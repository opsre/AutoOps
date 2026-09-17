package service

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math"
	"net"
	"os"
	"strings"
	"time"

	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const licenseConfigKey = "license"

const builtInLicensePublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoR3HFiJo8Znk/gQTuCYM
tU/Ce45WT6RVn30qp8sgw9V1qwKLUhnNoRkqfDRuZY/TWhBxi9eipyqd6z+XozNE
6CB9N9w5PkT/YqVcHY3N9kgxabKknYLTbYSOqBIqsNcWDi7F9HE5pjjzDuxB+PEw
ahPK9ibxOP/DPXLNZEHZyDBRUJq2HO389cGHrbJmhIcXn0IwLVrMNOtxyxevxJzV
U7UpF+WGigbad0O1rnw+Y8IH7M9I3tEoXWJSlDTslhmx7oFqGyBQ0ePO3MvbO2Xz
hcrt47qFVCBWE2ewarTHtcPpnRuOELW1RmH+eNS6+qeSr9PP4lOP4OsjXR0DeWrm
OwIDAQAB
-----END PUBLIC KEY-----`

type LicenseAvailability struct {
	Activated     bool
	Expired       bool
	MachineCode   string
	CustomerName  string
	LicenseType   string
	IssuedAt      string
	ActivatedAt   string
	ExpireAt      *string
	DaysRemaining int
	ValidityText  string
}

type LicenseService interface {
	GetMachineCode(c *gin.Context)
	GetLicenseStatus(c *gin.Context)
	Activate(c *gin.Context, req *model.ActivateLicenseRequest)
	GetAvailability() (*LicenseAvailability, error)
}

type LicenseServiceImpl struct {
	configDao dao.SysConfigDao
}

func NewLicenseService() LicenseService {
	return &LicenseServiceImpl{
		configDao: dao.NewSysConfigDao(),
	}
}

func normalizePEM(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}
	return strings.ReplaceAll(value, "\\n", "\n")
}

func parseLicenseToken(licenseKey string) (*model.LicensePayload, []byte, []byte, error) {
	parts := strings.Split(strings.TrimSpace(licenseKey), ".")
	if len(parts) != 2 {
		return nil, nil, nil, fmt.Errorf("激活码格式无效")
	}

	payloadRaw, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("激活码载荷解析失败")
	}
	signature, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("激活码签名解析失败")
	}

	var payload model.LicensePayload
	if err := json.Unmarshal(payloadRaw, &payload); err != nil {
		return nil, nil, nil, fmt.Errorf("激活码内容解析失败")
	}
	return &payload, payloadRaw, signature, nil
}

func parseLicensePublicKey() (*rsa.PublicKey, error) {
	publicKeyPEM := normalizePEM(os.Getenv("LICENSE_PUBLIC_KEY"))
	if publicKeyPEM == "" {
		publicKeyPEM = builtInLicensePublicKey
	}

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("LICENSE_PUBLIC_KEY 不是有效的 PEM 格式")
	}

	if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("LICENSE_PUBLIC_KEY 不是 RSA 公钥")
		}
		return rsaPub, nil
	}

	rsaPub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("LICENSE_PUBLIC_KEY 解析失败")
	}
	return rsaPub, nil
}

func verifyLicenseSignature(payloadRaw, signature []byte) error {
	publicKey, err := parseLicensePublicKey()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(payloadRaw)
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature); err != nil {
		return fmt.Errorf("激活码签名校验失败")
	}
	return nil
}

func readMachineID() string {
	candidates := []string{
		"/etc/machine-id",
		"/var/lib/dbus/machine-id",
	}
	for _, path := range candidates {
		if data, err := os.ReadFile(path); err == nil {
			value := strings.TrimSpace(string(data))
			if value != "" {
				return value
			}
		}
	}
	return ""
}

func readPrimaryMAC() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		if hw := strings.TrimSpace(iface.HardwareAddr.String()); hw != "" {
			return hw
		}
	}
	return ""
}

func generateMachineCode() string {
	hostname, _ := os.Hostname()
	seed := os.Getenv("LICENSE_MACHINE_SEED")
	source := strings.Join([]string{
		readMachineID(),
		hostname,
		readPrimaryMAC(),
		seed,
	}, "|")

	sum := sha256.Sum256([]byte(source))
	return hex.EncodeToString(sum[:])
}

func parseLicenseExpireAt(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	if strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", *value, time.Local)
	if err != nil {
		return nil, fmt.Errorf("expireAt 时间格式错误，应为 2006-01-02 15:04:05")
	}
	return &t, nil
}

func validateLicensePayload(payload *model.LicensePayload, machineCode string) error {
	if payload.MachineCode == "" {
		return fmt.Errorf("激活码缺少机器码")
	}
	if payload.MachineCode != machineCode {
		return fmt.Errorf("激活码与当前部署机器不匹配")
	}
	switch payload.LicenseType {
	case "half_year", "annual", "permanent":
	default:
		return fmt.Errorf("激活码中的授权类型无效")
	}

	expireAt, err := parseLicenseExpireAt(payload.ExpireAt)
	if err != nil {
		return err
	}
	if payload.LicenseType != "permanent" && expireAt == nil {
		return fmt.Errorf("限时授权必须包含过期时间")
	}
	if expireAt != nil && time.Now().After(*expireAt) {
		return fmt.Errorf("激活码已过期")
	}
	return nil
}

func (s *LicenseServiceImpl) loadStoredLicense() (*model.LicenseStored, error) {
	config, err := s.configDao.GetByKey(licenseConfigKey)
	if err != nil {
		return nil, err
	}
	var stored model.LicenseStored
	if err := json.Unmarshal([]byte(config.ConfigData), &stored); err != nil {
		return nil, err
	}
	return &stored, nil
}

func (s *LicenseServiceImpl) GetAvailability() (*LicenseAvailability, error) {
	machineCode := generateMachineCode()
	stored, err := s.loadStoredLicense()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &LicenseAvailability{
				Activated:   false,
				Expired:     false,
				MachineCode: machineCode,
			}, nil
		}
		return nil, err
	}

	// 已存授权绑定的是旧机器时，当前部署应视为未激活。
	if strings.TrimSpace(stored.MachineCode) != "" && stored.MachineCode != machineCode {
		return &LicenseAvailability{
			Activated:    false,
			Expired:      false,
			MachineCode:  machineCode,
			ValidityText: "未激活",
		}, nil
	}

	expired := false
	var expireAtTime *time.Time
	if stored.ExpireAt != nil {
		expireAtTime, err = parseLicenseExpireAt(stored.ExpireAt)
		if err != nil {
			return nil, err
		}
		if expireAtTime != nil && time.Now().After(*expireAtTime) {
			expired = true
		}
	}

	daysRemaining := 0
	validityText := "永久激活"
	if expireAtTime != nil && !expired {
		diff := expireAtTime.Sub(time.Now())
		daysRemaining = int(math.Ceil(diff.Hours() / 24))
		if diff > 0 && daysRemaining == 0 {
			daysRemaining = 1
		}
		validityText = fmt.Sprintf("剩余%d天", daysRemaining)
	} else if expireAtTime != nil && expired {
		validityText = "已过期"
	}

	return &LicenseAvailability{
		Activated:     stored.Activated,
		Expired:       expired,
		MachineCode:   machineCode,
		CustomerName:  stored.CustomerName,
		LicenseType:   stored.LicenseType,
		IssuedAt:      stored.IssuedAt,
		ActivatedAt:   stored.ActivatedAt,
		ExpireAt:      stored.ExpireAt,
		DaysRemaining: daysRemaining,
		ValidityText:  validityText,
	}, nil
}

func (s *LicenseServiceImpl) GetMachineCode(c *gin.Context) {
	result.Success(c, gin.H{
		"machineCode": generateMachineCode(),
	})
}

func (s *LicenseServiceImpl) GetLicenseStatus(c *gin.Context) {
	availability, err := s.GetAvailability()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("读取激活信息失败: %v", err))
		return
	}

	result.Success(c, gin.H{
		"activated":     availability.Activated,
		"expired":       availability.Expired,
		"machineCode":   availability.MachineCode,
		"customerName":  availability.CustomerName,
		"licenseType":   availability.LicenseType,
		"issuedAt":      availability.IssuedAt,
		"activatedAt":   availability.ActivatedAt,
		"expireAt":      availability.ExpireAt,
		"daysRemaining": availability.DaysRemaining,
		"validityText":  availability.ValidityText,
	})
}

func (s *LicenseServiceImpl) Activate(c *gin.Context, req *model.ActivateLicenseRequest) {
	machineCode := generateMachineCode()
	payload, payloadRaw, signature, err := parseLicenseToken(req.LicenseKey)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	if err := verifyLicenseSignature(payloadRaw, signature); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	if err := validateLicensePayload(payload, machineCode); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	stored := &model.LicenseStored{
		Activated:    true,
		LicenseKey:   req.LicenseKey,
		CustomerName: payload.CustomerName,
		MachineCode:  payload.MachineCode,
		LicenseType:  payload.LicenseType,
		IssuedAt:     payload.IssuedAt,
		ActivatedAt:  time.Now().Format("2006-01-02 15:04:05"),
		ExpireAt:     payload.ExpireAt,
	}

	configData, err := json.Marshal(stored)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("保存激活信息失败: %v", err))
		return
	}

	existing, err := s.configDao.GetByKey(licenseConfigKey)
	if err != nil && err != gorm.ErrRecordNotFound {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("读取激活配置失败: %v", err))
		return
	}

	if existing != nil {
		existing.ConfigType = "license"
		existing.ConfigData = string(configData)
		existing.Status = 1
		existing.Remark = "系统激活信息"
		if err := s.configDao.Update(existing); err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新激活信息失败: %v", err))
			return
		}
	} else {
		config := &model.SysConfig{
			ConfigKey:  licenseConfigKey,
			ConfigType: "license",
			ConfigData: string(configData),
			Status:     1,
			Remark:     "系统激活信息",
		}
		if err := s.configDao.Create(config); err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建激活信息失败: %v", err))
			return
		}
	}

	result.Success(c, gin.H{
		"activated":    true,
		"machineCode":  machineCode,
		"customerName": stored.CustomerName,
		"licenseType":  stored.LicenseType,
		"expireAt":     stored.ExpireAt,
		"validityText": func() string {
			if stored.LicenseType == "permanent" {
				return "永久激活"
			}
			if stored.ExpireAt != nil {
				if expireAtTime, err := parseLicenseExpireAt(stored.ExpireAt); err == nil && expireAtTime != nil {
					diff := expireAtTime.Sub(time.Now())
					if diff > 0 {
						return fmt.Sprintf("剩余%d天", int(math.Ceil(diff.Hours()/24)))
					}
					return "已过期"
				}
			}
			return ""
		}(),
	})
}
