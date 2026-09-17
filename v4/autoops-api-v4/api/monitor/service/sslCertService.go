package service

import (
	"archive/zip"
	"bytes"
	"dodevops-api/api/monitor/model"
	cmdbService "dodevops-api/api/cmdb/service"
	"dodevops-api/common"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/common/util/acme"
	"dodevops-api/common/util/aliyun"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// min 返回两个整数中的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// SslCertService SSL证书服务接口
type SslCertService interface {
	ApplySslCert(c *gin.Context, dto *model.CreateSslCertDto)
	QueryCertState(c *gin.Context, certID uint)
	DownloadCert(c *gin.Context, certID uint)
	ListSslCerts(c *gin.Context, dto *model.SslCertListDto)
	DeleteSslCert(c *gin.Context, certID uint)
	AddDnsValidation(c *gin.Context, certID uint)
	CheckDnsValidation(c *gin.Context, certID uint)
	DeploySslCert(c *gin.Context, req *model.DeploySslCertRequest)
	GetDeployLogs(c *gin.Context, req *model.GetDeployLogsRequest)
	UploadSslCert(c *gin.Context)
}

// SslCertServiceImpl SSL证书服务实现
type SslCertServiceImpl struct{}

// NewSslCertService 创建SSL证书服务实例
func NewSslCertService() SslCertService {
	return &SslCertServiceImpl{}
}

// ApplySslCert 申请SSL证书 - 支持阿里云CAS和ACME两种方式
func (s *SslCertServiceImpl) ApplySslCert(c *gin.Context, dto *model.CreateSslCertDto) {
	db := common.GetDB()

	// 1. 查询阿里云配置
	var aliyunConfig model.AliyunConfig
	if err := db.Where("id = ? AND status = 1", dto.AliyunConfigID).First(&aliyunConfig).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "阿里云配置不存在或已禁用")
			return
		}
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 2. 检查是否已存在相同域名的证书申请
	var existCert model.SslCert
	err := db.Where("domain = ? AND aliyun_config_id = ? AND status IN (?, ?)",
		dto.Domain, dto.AliyunConfigID, model.SslCertStatusApplying, model.SslCertStatusValidating).First(&existCert).Error
	if err == nil {
		result.Failed(c, int(result.ApiCode.FAILED), "该域名已有正在申请中的证书")
		return
	}

	// 3. 验证证书类型
	if dto.CertType != model.CertTypeACME && dto.CertType != model.CertTypeAliyunCAS {
		result.Failed(c, int(result.ApiCode.FAILED), "证书类型无效，必须为1(ACME)或2(阿里云CAS)")
		return
	}

	// 4. 根据证书类型设置证书来源
	var certSource string
	if dto.CertType == model.CertTypeACME {
		certSource = model.CertSourceACME
	} else {
		certSource = model.CertSourceAliyunCAS
	}

	// 5. 设置默认值
	validateType := dto.ValidateType
	if validateType == "" {
		validateType = "DNS"
	}

	algorithm := dto.Algorithm
	if algorithm == "" {
		algorithm = "RSA2048"
	}

	// 6. 创建初始证书记录（待处理状态）
	now := time.Now()
	sslCert := model.SslCert{
		Domain:         dto.Domain,
		CertSource:     certSource,
		CAProvider:     dto.CAProvider,
		AliyunConfigID: dto.AliyunConfigID,
		CertName:       fmt.Sprintf("cert-%s-%d", dto.Domain, now.Unix()),
		Algorithm:      algorithm,
		Status:         model.SslCertStatusApplying, // 申请中
		ValidateType:   validateType,
		CreateTime:     util.HTime{Time: now},
		UpdateTime:     util.HTime{Time: now},
	}

	// 7. 保存初始记录
	if err := db.Create(&sslCert).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 8. 立即返回结果给用户
	result.Success(c, gin.H{
		"message":        "证书申请已提交，正在后台处理",
		"certificateId":  sslCert.ID,
		"domain":         sslCert.Domain,
		"status":         sslCert.Status,
		"statusText":     "申请中",
		"certificateInfo": sslCert.ToVO(),
	})

	// 9. 异步处理证书申请
	go func() {
		// 使用新的数据库连接（协程安全）
		asyncDB := common.GetDB()

		if certSource == model.CertSourceACME {
			// ACME方式申请证书（支持通配符）
			s.applyAcmeCertAsync(asyncDB, sslCert.ID, dto, aliyunConfig, validateType, algorithm)
		} else {
			// 阿里云CAS方式申请证书（仅单域名）
			s.applyAliyunCasCertAsync(asyncDB, sslCert.ID, dto, aliyunConfig, validateType, algorithm)
		}
	}()
}

// applyAliyunCasCertAsync 异步通过阿里云CAS申请证书
func (s *SslCertServiceImpl) applyAliyunCasCertAsync(db *gorm.DB, certID uint, dto *model.CreateSslCertDto,
	aliyunConfig model.AliyunConfig, validateType, algorithm string) {

	fmt.Printf("[INFO] 开始异步申请阿里云CAS证书 - CertID: %d, Domain: %s\n", certID, dto.Domain)

	// 1. 创建CAS客户端
	casClient, err := aliyun.NewCasClient(aliyun.CasConfig{
		AccessKey:    aliyunConfig.AccessKey,
		AccessSecret: aliyunConfig.AccessSecret,
	})
	if err != nil {
		// 更新证书状态为失败
		db.Model(&model.SslCert{}).Where("id = ?", certID).Updates(map[string]interface{}{
			"status":      model.SslCertStatusApplying,
			"error_msg":   fmt.Sprintf("创建阿里云CAS客户端失败: %v", err),
			"update_time": time.Now(),
		})
		fmt.Printf("[ERROR] 创建阿里云CAS客户端失败 - CertID: %d, Error: %v\n", certID, err)
		return
	}

	// 2. 申请证书
	applyResp, err := casClient.ApplyFreeCert(aliyun.ApplyFreeCertRequest{
		Domain:       dto.Domain,
		Email:        aliyunConfig.Email,
		Phone:        aliyunConfig.Phone,
		Username:     aliyunConfig.Username,
		ValidateType: validateType,
		ProductCode:  "digicert-free-1-free",
	})
	if err != nil {
		// 更新证书状态为失败
		db.Model(&model.SslCert{}).Where("id = ?", certID).Updates(map[string]interface{}{
			"status":      model.SslCertStatusApplying,
			"error_msg":   fmt.Sprintf("申请证书失败: %v", err),
			"update_time": time.Now(),
		})
		fmt.Printf("[ERROR] 申请阿里云CAS证书失败 - CertID: %d, Error: %v\n", certID, err)
		return
	}

	// 3. 更新证书记录
	updates := map[string]interface{}{
		"order_id":     fmt.Sprintf("%d", applyResp.OrderID),
		"product_code": "digicert-free-1-free",
		"ca_provider":  model.CAProviderDigiCert,
		"status":       model.SslCertStatusApplying,
		"update_time":  time.Now(),
	}

	if err := db.Model(&model.SslCert{}).Where("id = ?", certID).Updates(updates).Error; err != nil {
		fmt.Printf("[ERROR] 更新证书记录失败 - CertID: %d, Error: %v\n", certID, err)
		return
	}

	fmt.Printf("[INFO] 阿里云CAS证书申请成功 - CertID: %d, OrderID: %d\n", certID, applyResp.OrderID)
}

// applyAcmeCertAsync 异步通过ACME协议申请证书
func (s *SslCertServiceImpl) applyAcmeCertAsync(db *gorm.DB, certID uint, dto *model.CreateSslCertDto,
	aliyunConfig model.AliyunConfig, validateType, algorithm string) {

	fmt.Printf("[INFO] 开始异步申请ACME证书 - CertID: %d, Domain: %s\n", certID, dto.Domain)

	caProvider := dto.CAProvider
	if caProvider == "" {
		caProvider = model.CAProviderLetsEncrypt // 默认使用Let's Encrypt
	}

	// 1. 创建ACME客户端
	acmeClient, err := acme.NewAcmeClient(acme.AcmeConfig{
		Email:      dto.Email,
		CAProvider: caProvider,
		Algorithm:  algorithm,
		EABKid:     aliyunConfig.EABKid,     // 传递EAB凭证
		EABHmacKey: aliyunConfig.EABHmacKey, // 传递EAB凭证
	})
	if err != nil {
		// 更新证书状态为失败
		db.Model(&model.SslCert{}).Where("id = ?", certID).Updates(map[string]interface{}{
			"status":      model.SslCertStatusApplying,
			"error_msg":   fmt.Sprintf("创建ACME客户端失败: %v", err),
			"update_time": time.Now(),
		})
		fmt.Printf("[ERROR] 创建ACME客户端失败 - CertID: %d, Error: %v\n", certID, err)
		return
	}

	// 2. 设置DNS提供商（使用阿里云DNS）
	dnsClient, err := aliyun.NewDnsClient(aliyun.DnsConfig{
		AccessKey:    aliyunConfig.AccessKey,
		AccessSecret: aliyunConfig.AccessSecret,
	})
	if err != nil {
		// 更新证书状态为失败
		db.Model(&model.SslCert{}).Where("id = ?", certID).Updates(map[string]interface{}{
			"status":      model.SslCertStatusApplying,
			"error_msg":   fmt.Sprintf("创建DNS客户端失败: %v", err),
			"update_time": time.Now(),
		})
		fmt.Printf("[ERROR] 创建DNS客户端失败 - CertID: %d, Error: %v\n", certID, err)
		return
	}

	// 设置DNS回调函数
	addFunc := func(fqdn, value string) error {
		fmt.Printf("[DEBUG] ACME添加DNS记录 - FQDN: %s, Value: %s\n", fqdn, value)
		_, err := dnsClient.AddDnsTxtRecord(aliyun.AddDnsTxtRecordRequest{
			Domain: fqdn,
			Value:  value,
			TTL:    600,
		})
		if err != nil {
			fmt.Printf("[ERROR] ACME添加DNS记录失败: %v\n", err)
		} else {
			fmt.Printf("[INFO] ACME添加DNS记录成功 - FQDN: %s\n", fqdn)
		}
		return err
	}

	cleanupFunc := func(fqdn, value string) error {
		fmt.Printf("[DEBUG] ACME清理DNS记录 - FQDN: %s\n", fqdn)
		err := dnsClient.CleanupCertValidationRecord(fqdn)
		if err != nil {
			fmt.Printf("[ERROR] ACME清理DNS记录失败: %v\n", err)
		} else {
			fmt.Printf("[INFO] ACME清理DNS记录成功 - FQDN: %s\n", fqdn)
		}
		return err
	}

	if err := acmeClient.SetDNSProvider(addFunc, cleanupFunc); err != nil {
		// 更新证书状态为失败
		db.Model(&model.SslCert{}).Where("id = ?", certID).Updates(map[string]interface{}{
			"status":      model.SslCertStatusApplying,
			"error_msg":   fmt.Sprintf("设置DNS提供商失败: %v", err),
			"update_time": time.Now(),
		})
		fmt.Printf("[ERROR] 设置DNS提供商失败 - CertID: %d, Error: %v\n", certID, err)
		return
	}

	// 3. 申请证书
	acmeResp, err := acmeClient.ApplyCert(acme.ApplyCertRequest{
		Domain: dto.Domain,
		Email:  dto.Email,
	})
	if err != nil {
		// 更新证书状态为失败
		db.Model(&model.SslCert{}).Where("id = ?", certID).Updates(map[string]interface{}{
			"status":      model.SslCertStatusApplying,
			"error_msg":   fmt.Sprintf("申请证书失败: %v", err),
			"update_time": time.Now(),
		})
		fmt.Printf("[ERROR] 申请ACME证书失败 - CertID: %d, Error: %v\n", certID, err)
		return
	}

	// 4. 解析有效期计算剩余天数
	issueTime, _ := time.Parse("2006-01-02", acmeResp.IssueDate)
	expireTime, _ := time.Parse("2006-01-02", acmeResp.ExpireDate)

	// 5. 更新证书记录为已签发
	updates := map[string]interface{}{
		"ca_provider": caProvider,
		"status":      model.SslCertStatusIssued, // ACME方式直接签发
		"cert":        acmeResp.Cert,
		"private_key": acmeResp.PrivateKey,
		"issuer_cert": acmeResp.IssuerCert,
		"issue_time":  issueTime,
		"expire_time": expireTime,
		"days_left":   acmeResp.DaysLeft,
		"error_msg":   "", // 清空错误信息
		"update_time": time.Now(),
	}

	if err := db.Model(&model.SslCert{}).Where("id = ?", certID).Updates(updates).Error; err != nil {
		fmt.Printf("[ERROR] 更新证书记录失败 - CertID: %d, Error: %v\n", certID, err)
		return
	}

	fmt.Printf("[INFO] ACME证书申请成功 - CertID: %d, Domain: %s, ExpiresIn: %d days\n", certID, dto.Domain, acmeResp.DaysLeft)
}

// QueryCertState 查询证书状态并更新
func (s *SslCertServiceImpl) QueryCertState(c *gin.Context, certID uint) {
	db := common.GetDB()

	// 1. 查询证书记录
	var sslCert model.SslCert
	if err := db.First(&sslCert, certID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "证书记录不存在")
			return
		}
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 2. 查询阿里云配置
	var aliyunConfig model.AliyunConfig
	if err := db.First(&aliyunConfig, sslCert.AliyunConfigID).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 3. 创建CAS客户端
	casClient, err := aliyun.NewCasClient(aliyun.CasConfig{
		AccessKey:    aliyunConfig.AccessKey,
		AccessSecret: aliyunConfig.AccessSecret,
	})
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建阿里云CAS客户端失败: %v", err))
		return
	}

	// 4. 解析订单ID
	var orderID int64
	fmt.Sscanf(sslCert.OrderID, "%d", &orderID)

	// 5. 查询证书状态
	stateResp, err := casClient.GetCertState(orderID)
	if err != nil {
		sslCert.ErrorMsg = err.Error()
		sslCert.UpdateTime = util.HTime{Time: time.Now()}
		db.Save(&sslCert)
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询证书状态失败: %v", err))
		return
	}

	// 6. 更新证书状态
	updates := map[string]interface{}{
		"validate_info": stateResp.ValidateInfo,
		"update_time":   time.Now(),
	}

	// 根据证书状态更新本地状态
	switch stateResp.Status {
	case "checking", "domain_verify":
		updates["status"] = model.SslCertStatusValidating
	case "certificate":
		updates["status"] = model.SslCertStatusIssued
	case "revoked", "expired":
		updates["status"] = model.SslCertStatusExpired
	}

	if err := db.Model(&sslCert).Updates(updates).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 重新查询更新后的记录
	db.First(&sslCert, certID)

	vo := sslCert.ToVO()
	vo.ValidateInfo = stateResp.ValidateInfo

	result.Success(c, vo)
}

// DownloadCert 下载证书（打包为zip文件）
func (s *SslCertServiceImpl) DownloadCert(c *gin.Context, certID uint) {
	db := common.GetDB()

	// 1. 查询证书记录
	var sslCert model.SslCert
	if err := db.First(&sslCert, certID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "证书记录不存在")
			return
		}
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 2. 如果已有证书内容,直接打包返回
	if sslCert.Cert != "" && sslCert.PrivateKey != "" {
		s.downloadCertAsZip(c, &sslCert)
		return
	}

	// 3. 查询阿里云配置
	var aliyunConfig model.AliyunConfig
	if err := db.First(&aliyunConfig, sslCert.AliyunConfigID).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	casClient, err := aliyun.NewCasClient(aliyun.CasConfig{
		AccessKey:    aliyunConfig.AccessKey,
		AccessSecret: aliyunConfig.AccessSecret,
	})
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建阿里云CAS客户端失败: %v", err))
		return
	}

	// 4. 解析OrderID并下载证书
	var orderID int64
	fmt.Sscanf(sslCert.OrderID, "%d", &orderID)

	// 先检查证书状态
	stateResp, err := casClient.GetCertState(orderID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询证书状态失败: %v", err))
		return
	}

	fmt.Printf("[DEBUG] 证书状态: %s, 域名: %s\n", stateResp.Status, stateResp.Domain)

	// 检查证书是否已签发
	if stateResp.Status != "certificate" {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("证书尚未签发，当前状态: %s", stateResp.Status))
		return
	}

	// 5. 通过OrderID下载证书
	var downloadResp *aliyun.DownloadCertResponse

	// 优先尝试使用OrderID下载
	downloadResp, err = casClient.DownloadCertByOrderID(orderID)
	if err != nil {
		// 如果OrderID下载失败，尝试使用CertID下载
		if sslCert.CertID != 0 {
			fmt.Printf("[DEBUG] OrderID下载失败: %v，尝试使用CertID: %d\n", err, sslCert.CertID)
			downloadResp, err = casClient.DownloadCert(sslCert.CertID)
			if err != nil {
				result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("下载证书失败: %v", err))
				return
			}
		} else {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("下载证书失败: %v", err))
			return
		}
	}

	// 6. 解析时间并更新证书记录
	issueTime, _ := time.Parse("2006-01-02", downloadResp.IssueDate)
	expireTime, _ := time.Parse("2006-01-02", downloadResp.ExpireDate)

	updates := map[string]interface{}{
		"cert_id":     sslCert.CertID,
		"cert":        downloadResp.Cert,
		"private_key": downloadResp.PrivateKey,
		"issuer_cert": downloadResp.CertChain, // 保存中间证书链
		"issue_time":  issueTime,
		"expire_time": expireTime,
		"days_left":   downloadResp.DaysLeft,
		"status":      model.SslCertStatusIssued,
		"update_time": time.Now(),
	}

	if err := db.Model(&sslCert).Updates(updates).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 重新查询更新后的记录
	db.First(&sslCert, certID)

	// 返回zip文件
	s.downloadCertAsZip(c, &sslCert)
}

// downloadCertAsZip 将证书和私钥打包为zip文件下载
func (s *SslCertServiceImpl) downloadCertAsZip(c *gin.Context, sslCert *model.SslCert) {
	// 创建zip缓冲区
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// 添加.pem文件（证书）
	pemFileName := sslCert.Domain + ".pem"
	pemFile, err := zipWriter.Create(pemFileName)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建pem文件失败: %v", err))
		return
	}
	if _, err := pemFile.Write([]byte(sslCert.Cert)); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("写入pem文件失败: %v", err))
		return
	}

	// 添加.key文件（私钥）
	keyFileName := sslCert.Domain + ".key"
	keyFile, err := zipWriter.Create(keyFileName)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建key文件失败: %v", err))
		return
	}
	if _, err := keyFile.Write([]byte(sslCert.PrivateKey)); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("写入key文件失败: %v", err))
		return
	}

	// 关闭zip写入器
	if err := zipWriter.Close(); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("关闭zip失败: %v", err))
		return
	}

	// 设置响应头
	zipFileName := sslCert.Domain + "-cert.zip"
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipFileName))
	c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	c.Writer.WriteHeader(200)

	// 返回zip文件
	c.Writer.Write(buf.Bytes())
}

// ListSslCerts 查询SSL证书列表
func (s *SslCertServiceImpl) ListSslCerts(c *gin.Context, dto *model.SslCertListDto) {
	db := common.GetDB()
	var sslCerts []model.SslCert
	var total int64

	query := db.Model(&model.SslCert{})

	// 条件过滤
	if dto.Domain != "" {
		query = query.Where("domain LIKE ?", "%"+dto.Domain+"%")
	}
	if dto.AliyunConfigID > 0 {
		query = query.Where("aliyun_config_id = ?", dto.AliyunConfigID)
	}
	if dto.Status > 0 {
		query = query.Where("status = ?", dto.Status)
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	page := dto.Page
	if page < 1 {
		page = 1
	}
	pageSize := dto.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	if err := query.Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&sslCerts).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 查询关联的阿里云配置名称
	configMap := make(map[uint]string)
	var configs []model.AliyunConfig
	db.Find(&configs)
	for _, config := range configs {
		configMap[config.ID] = config.Name
	}

	// 转换为VO
	vos := make([]*model.SslCertVO, 0, len(sslCerts))
	for i := range sslCerts {
		vo := sslCerts[i].ToVO()
		vo.ConfigName = configMap[sslCerts[i].AliyunConfigID]
		vos = append(vos, vo)
	}

	result.SuccessWithPage(c, vos, total, page, pageSize)
}

// DeleteSslCert 删除SSL证书
func (s *SslCertServiceImpl) DeleteSslCert(c *gin.Context, certID uint) {
	db := common.GetDB()

	var sslCert model.SslCert
	if err := db.First(&sslCert, certID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "证书记录不存在")
			return
		}
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 删除本地记录
	if err := db.Delete(&sslCert).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, "删除成功")
}

// AddDnsValidation 添加DNS验证记录
func (s *SslCertServiceImpl) AddDnsValidation(c *gin.Context, certID uint) {
	db := common.GetDB()

	// 1. 查询证书记录
	var sslCert model.SslCert
	if err := db.First(&sslCert, certID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "证书记录不存在")
			return
		}
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 2. 检查验证方式
	if sslCert.ValidateType != "DNS" {
		result.Failed(c, int(result.ApiCode.FAILED), "该证书不使用DNS验证方式")
		return
	}

	// 3. 查询阿里云配置
	var aliyunConfig model.AliyunConfig
	if err := db.First(&aliyunConfig, sslCert.AliyunConfigID).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 4. 先查询证书状态获取验证信息
	casClient, err := aliyun.NewCasClient(aliyun.CasConfig{
		AccessKey:    aliyunConfig.AccessKey,
		AccessSecret: aliyunConfig.AccessSecret,
	})
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建阿里云CAS客户端失败: %v", err))
		return
	}

	var orderID int64
	fmt.Sscanf(sslCert.OrderID, "%d", &orderID)

	// 自动重试机制: 阿里云生成验证信息可能需要几秒钟
	var stateResp *aliyun.CertStateResponse
	maxRetries := 3
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		resp, err := casClient.GetCertState(orderID)
		if err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询证书状态失败: %v", err))
			return
		}

		stateResp = resp

		// 记录调试信息
		fmt.Printf("[DEBUG] 第%d次查询证书状态 - OrderID: %d, Status: '%s', ValidateInfo: '%s'\n",
			i+1, stateResp.OrderID, stateResp.Status, stateResp.ValidateInfo)

		// 如果已经有验证信息,或者证书已签发,跳出循环
		if stateResp.ValidateInfo != "" || stateResp.Status == "certificate" {
			break
		}

		// 如果还没到最后一次重试,等待后继续
		if i < maxRetries-1 {
			fmt.Printf("[INFO] 验证信息为空,等待%v后重试...\n", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	// 检查证书状态，如果已签发则下载证书内容并更新数据库
	if stateResp.Status == "certificate" {
		fmt.Printf("[INFO] 证书已签发，开始下载证书内容 - OrderID: %d\n", orderID)

		// 下载证书内容（包含完整证书链和私钥）
		downloadResp, downloadErr := casClient.DownloadCertByOrderID(orderID)
		if downloadErr != nil {
			// 如果下载失败，至少更新状态
			fmt.Printf("[ERROR] 证书已签发但下载证书失败: %v\n", downloadErr)
			updates := map[string]interface{}{
				"status":      model.SslCertStatusIssued,
				"error_msg":   fmt.Sprintf("证书已签发，但下载失败: %v", downloadErr),
				"update_time": time.Now(),
			}
			db.Model(&sslCert).Updates(updates)

			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("证书已签发，但下载失败: %v", downloadErr))
			return
		}

		fmt.Printf("[INFO] 证书下载成功 - IssueDate: %s, ExpireDate: %s, DaysLeft: %d\n",
			downloadResp.IssueDate, downloadResp.ExpireDate, downloadResp.DaysLeft)

		// 解析时间并更新证书完整信息（包含证书内容、私钥、有效期）
		issueTime, _ := time.Parse("2006-01-02", downloadResp.IssueDate)
		expireTime, _ := time.Parse("2006-01-02", downloadResp.ExpireDate)

		updates := map[string]interface{}{
			"cert":        downloadResp.Cert,
			"private_key": downloadResp.PrivateKey,
			"issuer_cert": downloadResp.CertChain, // 保存中间证书链
			"issue_time":  issueTime,
			"expire_time": expireTime,
			"days_left":   downloadResp.DaysLeft,
			"status":      model.SslCertStatusIssued,
			"error_msg":   "", // 清空错误信息
			"update_time": time.Now(),
		}
		if err := db.Model(&sslCert).Updates(updates).Error; err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), err.Error())
			return
		}

		// 重新查询更新后的记录
		db.First(&sslCert, certID)

		// 返回完整的证书详细信息
		vo := sslCert.ToVO()
		result.Success(c, gin.H{
			"message":        "证书已签发并已下载",
			"status":         "issued",
			"certificateInfo": vo,
		})
		return
	}

	if stateResp.ValidateInfo == "" {
		result.Failed(c, int(result.ApiCode.FAILED),
			fmt.Sprintf("验证信息为空，证书状态: %s。阿里云证书申请后需要等待几秒钟才会生成验证信息，请稍后重试", stateResp.Status))
		return
	}

	// 5. 创建DNS客户端
	dnsClient, err := aliyun.NewDnsClient(aliyun.DnsConfig{
		AccessKey:    aliyunConfig.AccessKey,
		AccessSecret: aliyunConfig.AccessSecret,
	})
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建阿里云DNS客户端失败: %v", err))
		return
	}

	// 6. 添加DNS验证记录
	addResp, err := dnsClient.AddCertValidationRecord(sslCert.Domain, stateResp.ValidateInfo)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("添加DNS验证记录失败: %v", err))
		return
	}

	// 7. 更新证书记录
	updates := map[string]interface{}{
		"validate_info": stateResp.ValidateInfo,
		"status":        model.SslCertStatusValidating,
		"update_time":   time.Now(),
	}

	if err := db.Model(&sslCert).Updates(updates).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, gin.H{
		"recordId":     addResp.RecordID,
		"message":      addResp.Message,
		"validateInfo": stateResp.ValidateInfo,
	})
}

// CheckDnsValidation 检查DNS验证状态
func (s *SslCertServiceImpl) CheckDnsValidation(c *gin.Context, certID uint) {
	db := common.GetDB()

	// 1. 查询证书记录
	var sslCert model.SslCert
	if err := db.First(&sslCert, certID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "证书记录不存在")
			return
		}
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 2. 判断证书来源，ACME证书不需要调用阿里云API验证
	if sslCert.CertSource == "acme" {
		// ACME证书已经签发完成，直接返回当前状态
		response := map[string]interface{}{
			"status":      sslCert.Status,
			"statusText":  getStatusText(sslCert.Status),
			"certSource":  sslCert.CertSource,
			"message":     "ACME证书无需验证，当前状态已是最新",
			"issueTime":   sslCert.IssueTime,
			"expireTime":  sslCert.ExpireTime,
			"daysLeft":    sslCert.DaysLeft,
		}
		result.Success(c, response)
		return
	}

	// 3. 阿里云CAS证书需要查询验证状态
	// 查询阿里云配置
	var aliyunConfig model.AliyunConfig
	if err := db.First(&aliyunConfig, sslCert.AliyunConfigID).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 4. 创建CAS客户端查询证书状态
	casClient, err := aliyun.NewCasClient(aliyun.CasConfig{
		AccessKey:    aliyunConfig.AccessKey,
		AccessSecret: aliyunConfig.AccessSecret,
	})
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建阿里云CAS客户端失败: %v", err))
		return
	}

	var orderID int64
	fmt.Sscanf(sslCert.OrderID, "%d", &orderID)

	stateResp, err := casClient.GetCertState(orderID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询证书状态失败: %v", err))
		return
	}

	// 记录调试信息
	fmt.Printf("[DEBUG] CheckDnsValidation - OrderID: %d, Status: '%s', Type: %s\n",
		stateResp.OrderID, stateResp.Status, stateResp.Type)

	// 5. 检查证书状态，如果已签发则下载证书内容
	if stateResp.Status == "certificate" {
		fmt.Printf("[INFO] 证书已签发，开始下载证书内容 - OrderID: %d\n", orderID)

		// 下载证书内容（包含完整证书链和私钥）
		downloadResp, downloadErr := casClient.DownloadCertByOrderID(orderID)
		if downloadErr != nil {
			// 如果下载失败，至少更新状态
			fmt.Printf("[ERROR] 证书已签发但下载证书失败: %v\n", downloadErr)
			updates := map[string]interface{}{
				"status":      model.SslCertStatusIssued,
				"error_msg":   fmt.Sprintf("证书已签发，但下载失败: %v", downloadErr),
				"update_time": time.Now(),
			}
			db.Model(&sslCert).Updates(updates)
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("证书已签发，但下载失败: %v", downloadErr))
			return
		}

		fmt.Printf("[INFO] 证书下载成功 - IssueDate: %s, ExpireDate: %s, DaysLeft: %d\n",
			downloadResp.IssueDate, downloadResp.ExpireDate, downloadResp.DaysLeft)

		// 解析时间并更新证书完整信息（包含证书内容、私钥、有效期）
		issueTime, _ := time.Parse("2006-01-02", downloadResp.IssueDate)
		expireTime, _ := time.Parse("2006-01-02", downloadResp.ExpireDate)

		updates := map[string]interface{}{
			"cert":        downloadResp.Cert,
			"private_key": downloadResp.PrivateKey,
			"issuer_cert": downloadResp.CertChain, // 保存中间证书链
			"issue_time":  issueTime,
			"expire_time": expireTime,
			"days_left":   downloadResp.DaysLeft,
			"status":      model.SslCertStatusIssued,
			"error_msg":   "", // 清空错误信息
			"update_time": time.Now(),
		}

		if err := db.Model(&sslCert).Updates(updates).Error; err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), err.Error())
			return
		}

		// 重新查询更新后的记录
		db.First(&sslCert, certID)
		vo := sslCert.ToVO()

		result.Success(c, gin.H{
			"message":         "证书已签发并已下载",
			"status":          "issued",
			"certificateInfo": vo,
		})
		return
	}

	// 5. 其他状态的处理
	updates := map[string]interface{}{
		"update_time": time.Now(),
	}

	switch stateResp.Status {
	case "checking", "domain_verify":
		updates["status"] = model.SslCertStatusValidating
	case "revoked", "expired":
		updates["status"] = model.SslCertStatusExpired
	}

	if err := db.Model(&sslCert).Updates(updates).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, gin.H{
		"status":       stateResp.Status,
		"type":         stateResp.Type,
		"domain":       stateResp.Domain,
		"validateType": stateResp.ValidateType,
		"validateInfo": stateResp.ValidateInfo,
	})
}

// getStatusText 获取证书状态文本
func getStatusText(status int) string {
	statusMap := map[int]string{
		1: "待验证",
		2: "验证中",
		3: "已签发",
		4: "已过期",
		5: "已撤销",
	}
	if text, ok := statusMap[status]; ok {
		return text
	}
	return "未知状态"
}

// DeploySslCert 部署SSL证书到目标主机（异步）
func (s *SslCertServiceImpl) DeploySslCert(c *gin.Context, req *model.DeploySslCertRequest) {
	db := common.GetDB()

	// 1. 查询证书信息
	var sslCert model.SslCert
	if err := db.Where("id = ?", req.CertID).First(&sslCert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "证书不存在")
			return
		}
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 2. 检查证书是否已签发
	if sslCert.Status != model.SslCertStatusIssued {
		result.Failed(c, int(result.ApiCode.FAILED), "证书未签发，无法部署")
		return
	}

	// 3. 检查证书内容是否存在
	if sslCert.Cert == "" || sslCert.PrivateKey == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "证书内容或私钥不存在")
		return
	}

	// 4. 查询主机信息
	var host struct {
		ID       uint   `gorm:"column:id"`
		HostName string `gorm:"column:host_name"`
		SSHIP    string `gorm:"column:ssh_ip"`
		SSHPort  int    `gorm:"column:ssh_port"`
		SSHName  string `gorm:"column:ssh_name"`
		SSHKeyID uint   `gorm:"column:ssh_key_id"`
	}
	if err := db.Table("cmdb_host").Where("id = ?", req.HostID).First(&host).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "主机不存在")
			return
		}
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 5. 创建部署日志记录
	now := time.Now()
	deployLog := model.SslCertDeployLog{
		CertID:     sslCert.ID,
		Domain:     sslCert.Domain,
		HostID:     host.ID,
		HostName:   host.HostName,
		DeployPath: req.DeployPath,
		Status:     1, // 部署中
		CreateTime: util.HTime{Time: now},
		UpdateTime: util.HTime{Time: now},
	}

	if err := db.Create(&deployLog).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "创建部署日志失败: "+err.Error())
		return
	}

	// 6. 立即返回响应，告知前端已开始异步部署
	result.Success(c, &model.DeploySslCertResponse{
		Success:     true,
		Message:     "证书部署任务已创建，正在后台执行部署",
		DeployLogID: deployLog.ID,
		Logs:        []string{fmt.Sprintf("[%s] 部署任务已创建，开始后台部署", time.Now().Format("2006-01-02 15:04:05"))},
	})

	// 7. 在后台goroutine中执行实际的部署操作
	go s.executeDeployment(deployLog.ID, sslCert.ID, sslCert.Domain, sslCert.Cert, sslCert.PrivateKey, sslCert.IssuerCert, host.ID, host.HostName, host.SSHIP, req.DeployPath)
}

// executeDeployment 后台执行证书部署
func (s *SslCertServiceImpl) executeDeployment(deployLogID uint, certID uint, domain string, cert string, privateKey string, issuerCert string, hostID uint, hostName string, hostIP string, deployPath string) {
	db := common.GetDB()
	logs := []string{}
	backupFiles := []string{}
	deployFiles := []string{}
	var deployErr error

	// 创建一个新的gin.Context用于SSH操作
	// 注意：这里使用background context，因为原始请求已经结束
	c := &gin.Context{}

	// 执行部署
	{
		logs = append(logs, fmt.Sprintf("[%s] 开始部署证书到主机: %s (%s)", time.Now().Format("2006-01-02 15:04:05"), hostName, hostIP))

		// 创建SSH连接 (使用现有的SSH服务)
		sshService := &cmdbService.CmdbHostSSHServiceImpl{}

		// 1. 确保目标目录存在
		logs = append(logs, fmt.Sprintf("[%s] 检查部署目录: %s", time.Now().Format("2006-01-02 15:04:05"), deployPath))
		mkdirCmd := fmt.Sprintf("mkdir -p %s", deployPath)
		_, err := sshService.ExecuteCommand(c, hostID, mkdirCmd)
		if err != nil {
			deployErr = fmt.Errorf("创建目录失败: %v", err)
			goto DEPLOY_FAILED
		}
		logs = append(logs, fmt.Sprintf("[%s] 目录已就绪", time.Now().Format("2006-01-02 15:04:05")))

		// 2. 检查旧证书并备份
		timestamp := time.Now().Format("20060102150405")
		certFileName := fmt.Sprintf("%s.pem", domain)
		keyFileName := fmt.Sprintf("%s.key", domain)
		certFilePath := fmt.Sprintf("%s/%s", deployPath, certFileName)
		keyFilePath := fmt.Sprintf("%s/%s", deployPath, keyFileName)

		logs = append(logs, fmt.Sprintf("[%s] 检查旧证书文件", time.Now().Format("2006-01-02 15:04:05")))

		// 检查并备份证书文件
		checkCertCmd := fmt.Sprintf("[ -f %s ] && echo 'exists' || echo 'not_exists'", certFilePath)
		certCheckResp, err := sshService.ExecuteCommand(c, hostID, checkCertCmd)
		if err == nil && certCheckResp != nil && certCheckResp.Output == "exists" {
			backupCertFile := fmt.Sprintf("%s.%s.bak", certFilePath, timestamp)
			backupCmd := fmt.Sprintf("cp %s %s", certFilePath, backupCertFile)
			_, err := sshService.ExecuteCommand(c, hostID, backupCmd)
			if err != nil {
				logs = append(logs, fmt.Sprintf("[%s] 警告: 备份证书文件失败: %v", time.Now().Format("2006-01-02 15:04:05"), err))
			} else {
				backupFiles = append(backupFiles, backupCertFile)
				logs = append(logs, fmt.Sprintf("[%s] 已备份证书文件: %s", time.Now().Format("2006-01-02 15:04:05"), backupCertFile))
			}
		}

		// 检查并备份密钥文件
		checkKeyCmd := fmt.Sprintf("[ -f %s ] && echo 'exists' || echo 'not_exists'", keyFilePath)
		keyCheckResp, err := sshService.ExecuteCommand(c, hostID, checkKeyCmd)
		if err == nil && keyCheckResp != nil && keyCheckResp.Output == "exists" {
			backupKeyFile := fmt.Sprintf("%s.%s.bak", keyFilePath, timestamp)
			backupCmd := fmt.Sprintf("cp %s %s", keyFilePath, backupKeyFile)
			_, err := sshService.ExecuteCommand(c, hostID, backupCmd)
			if err != nil {
				logs = append(logs, fmt.Sprintf("[%s] 警告: 备份密钥文件失败: %v", time.Now().Format("2006-01-02 15:04:05"), err))
			} else {
				backupFiles = append(backupFiles, backupKeyFile)
				logs = append(logs, fmt.Sprintf("[%s] 已备份密钥文件: %s", time.Now().Format("2006-01-02 15:04:05"), backupKeyFile))
			}
		}

		// 3. 写入新证书文件
		logs = append(logs, fmt.Sprintf("[%s] 开始部署新证书", time.Now().Format("2006-01-02 15:04:05")))

		// 准备证书内容
		// 注意: cert字段可能已经包含完整证书链(ACME协议),或者只包含域名证书(阿里云CAS)
		fullChainCert := cert

		// 如果有单独的中间证书,并且cert中不包含它,则追加
		if issuerCert != "" && !strings.Contains(cert, strings.TrimSpace(issuerCert)) {
			// 确保证书之间有换行
			if !strings.HasSuffix(cert, "\n") {
				fullChainCert = cert + "\n"
			} else {
				fullChainCert = cert
			}
			fullChainCert = fullChainCert + issuerCert
			logs = append(logs, fmt.Sprintf("[%s] 已合并证书链(域名证书+中间证书)", time.Now().Format("2006-01-02 15:04:05")))
		} else if strings.Contains(cert, "-----END CERTIFICATE-----") && strings.Count(cert, "-----BEGIN CERTIFICATE-----") > 1 {
			logs = append(logs, fmt.Sprintf("[%s] 检测到证书已包含完整链", time.Now().Format("2006-01-02 15:04:05")))
		}

		// 写入完整证书链
		writeCertCmd := fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", certFilePath, fullChainCert)
		_, err = sshService.ExecuteCommand(c, hostID, writeCertCmd)
		if err != nil {
			deployErr = fmt.Errorf("写入证书文件失败: %v", err)
			goto DEPLOY_FAILED
		}
		deployFiles = append(deployFiles, certFilePath)
		logs = append(logs, fmt.Sprintf("[%s] 证书文件已部署: %s", time.Now().Format("2006-01-02 15:04:05"), certFilePath))

		// 写入私钥
		writeKeyCmd := fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", keyFilePath, privateKey)
		_, err = sshService.ExecuteCommand(c, hostID, writeKeyCmd)
		if err != nil {
			deployErr = fmt.Errorf("写入私钥文件失败: %v", err)
			goto DEPLOY_FAILED
		}
		deployFiles = append(deployFiles, keyFilePath)
		logs = append(logs, fmt.Sprintf("[%s] 私钥文件已部署: %s", time.Now().Format("2006-01-02 15:04:05"), keyFilePath))

		// 4. 设置文件权限
		chmodCmd := fmt.Sprintf("chmod 600 %s && chmod 600 %s", certFilePath, keyFilePath)
		_, err = sshService.ExecuteCommand(c, hostID, chmodCmd)
		if err != nil {
			logs = append(logs, fmt.Sprintf("[%s] 警告: 设置文件权限失败: %v", time.Now().Format("2006-01-02 15:04:05"), err))
		} else {
			logs = append(logs, fmt.Sprintf("[%s] 已设置文件权限为 600", time.Now().Format("2006-01-02 15:04:05")))
		}

		// 5. 重载nginx
		logs = append(logs, fmt.Sprintf("[%s] 开始重载Nginx", time.Now().Format("2006-01-02 15:04:05")))
		reloadCmd := "nginx -s reload"
		reloadResp, err := sshService.ExecuteCommand(c, hostID, reloadCmd)
		if err != nil {
			deployErr = fmt.Errorf("重载Nginx失败: %v", err)
			goto DEPLOY_FAILED
		}
		logs = append(logs, fmt.Sprintf("[%s] Nginx重载成功", time.Now().Format("2006-01-02 15:04:05")))
		if reloadResp != nil && reloadResp.Output != "" {
			logs = append(logs, fmt.Sprintf("[%s] Nginx输出: %s", time.Now().Format("2006-01-02 15:04:05"), reloadResp.Output))
		}

		logs = append(logs, fmt.Sprintf("[%s] 证书部署完成", time.Now().Format("2006-01-02 15:04:05")))
	}

	// 更新部署日志为成功
	{
		backupFilesJSON, _ := json.Marshal(backupFiles)
		deployFilesJSON, _ := json.Marshal(deployFiles)
		logsJSON, _ := json.Marshal(logs)

		updates := map[string]interface{}{
			"status":       2, // 成功
			"backup_files": string(backupFilesJSON),
			"deploy_files": string(deployFilesJSON),
			"logs":         string(logsJSON),
			"update_time":  time.Now(),
		}

		db.Model(&model.SslCertDeployLog{}).Where("id = ?", deployLogID).Updates(updates)
	}
	return

DEPLOY_FAILED:
	// 部署失败，更新日志
	{
		backupFilesJSON, _ := json.Marshal(backupFiles)
		deployFilesJSON, _ := json.Marshal(deployFiles)
		logsJSON, _ := json.Marshal(logs)

		logs = append(logs, fmt.Sprintf("[%s] 部署失败: %v", time.Now().Format("2006-01-02 15:04:05"), deployErr))
		logsJSON, _ = json.Marshal(logs)

		updates := map[string]interface{}{
			"status":       3, // 失败
			"backup_files": string(backupFilesJSON),
			"deploy_files": string(deployFilesJSON),
			"logs":         string(logsJSON),
			"error_msg":    deployErr.Error(),
			"update_time":  time.Now(),
		}

		db.Model(&model.SslCertDeployLog{}).Where("id = ?", deployLogID).Updates(updates)
	}
}

// GetDeployLogs 获取SSL证书部署日志列表
func (s *SslCertServiceImpl) GetDeployLogs(c *gin.Context, req *model.GetDeployLogsRequest) {
	db := common.GetDB()

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 构建查询条件
	query := db.Model(&model.SslCertDeployLog{})

	// 按证书ID筛选
	if req.CertID > 0 {
		query = query.Where("cert_id = ?", req.CertID)
	}

	// 按主机ID筛选
	if req.HostID > 0 {
		query = query.Where("host_id = ?", req.HostID)
	}

	// 按域名筛选(模糊查询)
	if req.Domain != "" {
		query = query.Where("domain LIKE ?", "%"+req.Domain+"%")
	}

	// 按状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询总数失败: "+err.Error())
		return
	}

	// 查询列表
	var logs []model.SslCertDeployLog
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("create_time DESC").Offset(offset).Limit(req.PageSize).Find(&logs).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询列表失败: "+err.Error())
		return
	}

	// 转换为VO
	list := make([]model.DeployLogVO, 0, len(logs))
	for _, log := range logs {
		vo := model.DeployLogVO{
			ID:         log.ID,
			CertID:     log.CertID,
			Domain:     log.Domain,
			HostID:     log.HostID,
			HostName:   log.HostName,
			DeployPath: log.DeployPath,
			Status:     log.Status,
			ErrorMsg:   log.ErrorMsg,
			CreateTime: log.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime: log.UpdateTime.Format("2006-01-02 15:04:05"),
		}

		// 状态文本
		statusMap := map[int]string{
			1: "部署中",
			2: "成功",
			3: "失败",
		}
		if text, ok := statusMap[log.Status]; ok {
			vo.StatusText = text
		} else {
			vo.StatusText = "未知"
		}

		// 解析JSON字段
		if log.BackupFiles != "" {
			var backupFiles []string
			if err := json.Unmarshal([]byte(log.BackupFiles), &backupFiles); err == nil {
				vo.BackupFiles = backupFiles
			}
		}

		if log.DeployFiles != "" {
			var deployFiles []string
			if err := json.Unmarshal([]byte(log.DeployFiles), &deployFiles); err == nil {
				vo.DeployFiles = deployFiles
			}
		}

		if log.Logs != "" {
			var logs []string
			if err := json.Unmarshal([]byte(log.Logs), &logs); err == nil {
				vo.Logs = logs
			}
		}

		list = append(list, vo)
	}

	// 返回结果
	result.Success(c, &model.GetDeployLogsResponse{
		Total: int(total),
		List:  list,
	})
}

// UploadSslCert 上传SSL证书(支持tar/zip压缩包)
func (s *SslCertServiceImpl) UploadSslCert(c *gin.Context) {
	db := common.GetDB()

	// 1. 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "文件上传失败: "+err.Error())
		return
	}

	// 2. 检查文件类型(仅支持zip)
	filename := file.Filename
	if !strings.HasSuffix(filename, ".zip") {
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的文件格式,仅支持 zip")
		return
	}

	// 3. 打开文件
	fileContent, err := file.Open()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "文件读取失败: "+err.Error())
		return
	}
	defer fileContent.Close()

	// 4. 解析ZIP压缩包
	var certContent, keyContent, issuerCert, domain string
	var certFilename, keyFilename string

	// 解析ZIP
	readerAt, err := file.Open()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "文件打开失败: "+err.Error())
		return
	}
	defer readerAt.Close()

	// 读取文件内容到内存
	fileBytes := make([]byte, file.Size)
	if _, err := readerAt.Read(fileBytes); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "文件读取失败: "+err.Error())
		return
	}

	zipReader, err := zip.NewReader(bytes.NewReader(fileBytes), file.Size)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ZIP解析失败: "+err.Error())
		return
	}

	for _, zipFile := range zipReader.File {
		fname := zipFile.Name
		fmt.Printf("[DEBUG] 解析ZIP文件: %s, Size: %d\n", fname, zipFile.UncompressedSize64)

		// 跳过目录
		if zipFile.FileInfo().IsDir() {
			fmt.Printf("[DEBUG] 跳过目录: %s\n", fname)
			continue
		}

		// 跳过隐藏文件和系统文件
		baseName := fname
		if idx := strings.LastIndex(fname, "/"); idx != -1 {
			baseName = fname[idx+1:]
		}
		if strings.HasPrefix(baseName, ".") || strings.HasPrefix(baseName, "__") {
			fmt.Printf("[DEBUG] 跳过隐藏/系统文件: %s\n", fname)
			continue
		}

		rc, err := zipFile.Open()
		if err != nil {
			fmt.Printf("[ERROR] 打开文件失败 %s: %v\n", fname, err)
			continue
		}

		// 使用io.ReadAll读取完整内容
		contentBytes, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			fmt.Printf("[ERROR] 读取文件失败 %s: %v\n", fname, err)
			continue
		}

		fileContent := string(contentBytes)
		fmt.Printf("[DEBUG] 文件 %s 内容长度: %d, 前50字符: %s\n", fname, len(fileContent),
			fileContent[:min(50, len(fileContent))])

		// 判断文件类型 - 先检查是否是私钥
		lowerFname := strings.ToLower(fname)
		if strings.HasSuffix(lowerFname, ".key") || strings.Contains(lowerFname, "private") || strings.Contains(lowerFname, "key") {
			// 私钥文件
			if strings.Contains(fileContent, "BEGIN") && strings.Contains(fileContent, "PRIVATE KEY") {
				keyContent = fileContent
				keyFilename = fname
				fmt.Printf("[INFO] 找到私钥文件: %s\n", fname)
			}
		} else if strings.HasSuffix(lowerFname, ".pem") || strings.HasSuffix(lowerFname, ".crt") ||
			strings.HasSuffix(lowerFname, ".cer") || strings.HasSuffix(lowerFname, ".cert") {
			// 证书文件
			if strings.Contains(fileContent, "BEGIN CERTIFICATE") {
				certContent = fileContent
				certFilename = fname
				fmt.Printf("[INFO] 找到证书文件: %s\n", fname)
			}
		}
	}

	// 5. 验证是否解析到证书和私钥
	if certContent == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "未找到有效的证书文件(.pem/.crt/.cer)")
		return
	}
	if keyContent == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "未找到有效的私钥文件(.key)")
		return
	}

	// 6. 从文件名提取域名
	// 优先从私钥文件名提取,其次从证书文件名提取
	if keyFilename != "" {
		domain = util.ExtractDomainFromFilename(keyFilename)
	}
	if domain == "" && certFilename != "" {
		domain = util.ExtractDomainFromFilename(certFilename)
	}
	if domain == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "无法从文件名提取域名,请确保文件名格式为 domain.key 或 domain.pem")
		return
	}

	fmt.Printf("[INFO] 从文件名提取域名: %s (key=%s, cert=%s)\n", domain, keyFilename, certFilename)

	// 7. 解析证书信息
	certInfo, err := util.ParseCertificate(certContent)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "证书解析失败: "+err.Error())
		return
	}

	// 8. 检查证书链,分离域名证书和中间证书
	certBlocks := strings.Count(certContent, "-----BEGIN CERTIFICATE-----")
	if certBlocks > 1 {
		// 包含多个证书,分离域名证书和中间证书
		certs := strings.Split(certContent, "-----END CERTIFICATE-----")
		if len(certs) >= 2 {
			// 第一个是域名证书
			domainCert := certs[0] + "-----END CERTIFICATE-----\n"
			// 其余是中间证书
			var issuerCerts []string
			for i := 1; i < len(certs)-1; i++ {
				if strings.Contains(certs[i], "-----BEGIN CERTIFICATE-----") {
					issuerCerts = append(issuerCerts, strings.TrimSpace(certs[i])+"-----END CERTIFICATE-----\n")
				}
			}
			if len(issuerCerts) > 0 {
				certContent = domainCert
				issuerCert = strings.Join(issuerCerts, "")
				fmt.Printf("[INFO] 检测到证书链: 域名证书 + %d 个中间证书\n", len(issuerCerts))
			}
		}
	}

	// 9. 创建证书记录
	sslCert := model.SslCert{
		Domain:         domain,
		CertSource:     "manual_upload", // 手动上传
		CAProvider:     certInfo.Issuer,
		Cert:           certContent,
		PrivateKey:     keyContent,
		IssuerCert:     issuerCert,
		Algorithm:      certInfo.Algorithm,
		Status:         model.SslCertStatusIssued, // 已签发
		ValidateType:   "manual",
		IssueTime:      util.HTime{Time: certInfo.NotBefore},
		ExpireTime:     util.HTime{Time: certInfo.NotAfter},
		DaysLeft:       certInfo.DaysLeft,
		CreateTime:     util.HTime{Time: time.Now()},
		UpdateTime:     util.HTime{Time: time.Now()},
	}

	if err := db.Create(&sslCert).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "证书保存失败: "+err.Error())
		return
	}

	fmt.Printf("[INFO] 手动上传证书成功 - CertID: %d, Domain: %s, ExpiresIn: %d days\n",
		sslCert.ID, sslCert.Domain, sslCert.DaysLeft)

	result.Success(c, &model.UploadSslCertResponse{
		Success: true,
		Message: "证书上传成功",
		CertID:  sslCert.ID,
		Domain:  sslCert.Domain,
	})
}


