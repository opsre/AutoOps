package model

import (
	"dodevops-api/common/util"
)

// SslCert SSL证书模型
type SslCert struct {
	ID             uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Domain         string     `gorm:"column:domain;varchar(255);comment:'域名';NOT NULL;index" json:"domain"`
	CertSource     string     `gorm:"column:cert_source;varchar(32);comment:'证书来源:aliyun_cas/acme';default:'aliyun_cas';index" json:"certSource"`
	CAProvider     string     `gorm:"column:ca_provider;varchar(64);comment:'CA提供商:DigiCert/LetsEncrypt/ZeroSSL';default:'DigiCert'" json:"caProvider"`
	AliyunConfigID uint       `gorm:"column:aliyun_config_id;comment:'阿里云配置ID';index" json:"aliyunConfigId"`
	OrderID        string     `gorm:"column:order_id;varchar(128);comment:'证书订单ID';index" json:"orderId"`
	CertID         int64      `gorm:"column:cert_id;comment:'证书ID'" json:"certId"`
	CertName       string     `gorm:"column:cert_name;varchar(255);comment:'证书名称'" json:"certName"`
	ProductCode    string     `gorm:"column:product_code;varchar(64);comment:'产品代码';default:'digicert-free-1-free'" json:"productCode"`
	Status         int        `gorm:"column:status;comment:'状态:1->申请中,2->验证中,3->已签发,4->已过期,5->申请失败'" json:"status"`
	ValidateType   string     `gorm:"column:validate_type;varchar(32);comment:'验证方式:DNS/FILE';default:'DNS'" json:"validateType"`
	ValidateInfo   string     `gorm:"column:validate_info;type:text;comment:'验证信息(DNS记录或FILE内容)'" json:"validateInfo"`
	Cert           string     `gorm:"column:cert;type:text;comment:'证书内容'" json:"cert"`
	PrivateKey     string     `gorm:"column:private_key;type:text;comment:'私钥内容'" json:"privateKey"`
	IssuerCert     string     `gorm:"column:issuer_cert;type:text;comment:'颁发者证书'" json:"issuerCert"`
	Algorithm      string     `gorm:"column:algorithm;varchar(32);comment:'加密算法:RSA2048/EC256';default:'RSA2048'" json:"algorithm"`
	IssueTime      util.HTime `gorm:"column:issue_time;comment:'签发时间'" json:"issueTime"`
	ExpireTime     util.HTime `gorm:"column:expire_time;comment:'过期时间'" json:"expireTime"`
	DaysLeft       int        `gorm:"column:days_left;comment:'剩余天数'" json:"daysLeft"`
	ErrorMsg       string     `gorm:"column:error_msg;type:text;comment:'错误信息'" json:"errorMsg"`
	CreateTime     util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime     util.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
}

func (SslCert) TableName() string {
	return "monitor_ssl_cert"
}

// SslCert状态常量
const (
	SslCertStatusApplying    = 1 // 申请中
	SslCertStatusValidating  = 2 // 验证中
	SslCertStatusIssued      = 3 // 已签发
	SslCertStatusExpired     = 4 // 已过期
	SslCertStatusApplyFailed = 5 // 申请失败
)

// 证书来源常量
const (
	CertSourceAliyunCAS = "aliyun_cas" // 阿里云CAS申请
	CertSourceACME      = "acme"       // ACME协议申请（Let's Encrypt等）
)

// CA提供商常量
const (
	CAProviderDigiCert     = "DigiCert"     // DigiCert（阿里云CAS）
	CAProviderLetsEncrypt  = "LetsEncrypt"  // Let's Encrypt
	CAProviderZeroSSL      = "ZeroSSL"      // ZeroSSL
	CAProviderGoogle       = "Google"       // Google Trust Services
	CAProviderSSLCOM       = "SSLCOM"       // SSL.COM
	CAProviderBuyPass      = "BuyPass"      // BuyPass
)

// GetStatusText 获取状态文本
func (s *SslCert) GetStatusText() string {
	statusMap := map[int]string{
		SslCertStatusApplying:    "申请中",
		SslCertStatusValidating:  "验证中",
		SslCertStatusIssued:      "已签发",
		SslCertStatusExpired:     "已过期",
		SslCertStatusApplyFailed: "申请失败",
	}
	if text, ok := statusMap[s.Status]; ok {
		return text
	}
	return "未知状态"
}

// IsValid 判断证书是否有效
func (s *SslCert) IsValid() bool {
	return s.Status == SslCertStatusIssued && s.DaysLeft > 0
}

// AliyunConfig 阿里云配置模型
type AliyunConfig struct {
	ID           uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Name         string     `gorm:"column:name;varchar(128);comment:'配置名称';NOT NULL" json:"name"`
	AccessKey    string     `gorm:"column:access_key;varchar(128);comment:'AccessKey';NOT NULL" json:"accessKey"`
	AccessSecret string     `gorm:"column:access_secret;varchar(128);comment:'AccessSecret';NOT NULL" json:"accessSecret"`
	Region       string     `gorm:"column:region;varchar(64);comment:'区域';default:'cn-hangzhou'" json:"region"`
	Email        string     `gorm:"column:email;varchar(128);comment:'联系邮箱'" json:"email"`
	Phone        string     `gorm:"column:phone;varchar(32);comment:'联系电话'" json:"phone"`
	Username     string     `gorm:"column:username;varchar(64);comment:'联系人姓名'" json:"username"`
	EABKid       string     `gorm:"column:eab_kid;varchar(256);comment:'ACME EAB Key ID (ZeroSSL等需要)'" json:"eabKid"`
	EABHmacKey   string     `gorm:"column:eab_hmac_key;varchar(256);comment:'ACME EAB HMAC Key (ZeroSSL等需要)'" json:"eabHmacKey"`
	Status       int        `gorm:"column:status;comment:'状态:1->启用,0->禁用';default:1" json:"status"`
	CreateTime   util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime   util.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
}

func (AliyunConfig) TableName() string {
	return "monitor_aliyun_config"
}

// 证书类型常量
const (
	CertTypeACME      = 1 // ACME协议证书
	CertTypeAliyunCAS = 2 // 阿里云CAS证书
)

// CreateSslCertDto 创建SSL证书DTO
type CreateSslCertDto struct {
	Domain         string `validate:"required" json:"domain"`         // 域名（支持通配符*.example.com）
	CertType       int    `validate:"required" json:"certType"`       // 证书类型: 1=ACME, 2=阿里云CAS
	CertSource     string `json:"certSource"`                         // 证书来源: aliyun_cas/acme（内部使用，根据CertType自动设置）
	CAProvider     string `json:"caProvider"`                         // CA提供商（ACME方式时需要）
	AliyunConfigID uint   `validate:"required" json:"aliyunConfigId"` // 阿里云配置ID
	ValidateType   string `json:"validateType"`                       // 验证方式(DNS/FILE)
	Algorithm      string `json:"algorithm"`                          // 加密算法（ACME方式时可选）
	Email          string `json:"email"`                              // 邮箱（ACME方式时需要）
}

// SslCertListDto SSL证书列表查询DTO
type SslCertListDto struct {
	Domain         string `json:"domain"`         // 域名(模糊查询)
	AliyunConfigID uint   `json:"aliyunConfigId"` // 阿里云配置ID
	Status         int    `json:"status"`         // 状态
	Page           int    `json:"page"`           // 页码
	PageSize       int    `json:"pageSize"`       // 页大小
}

// SslCertVO SSL证书视图对象
type SslCertVO struct {
	ID             uint       `json:"id"`
	Domain         string     `json:"domain"`
	CertSource     string     `json:"certSource"`     // 证书来源
	CertSourceText string     `json:"certSourceText"` // 证书来源文本
	CAProvider     string     `json:"caProvider"`     // CA提供商
	AliyunConfigID uint       `json:"aliyunConfigId"`
	ConfigName     string     `json:"configName"` // 配置名称
	OrderID        string     `json:"orderId"`
	CertID         int64      `json:"certId"`
	CertName       string     `json:"certName"`
	ProductCode    string     `json:"productCode"`
	Algorithm      string     `json:"algorithm"`
	Status         int        `json:"status"`
	StatusText     string     `json:"statusText"`
	ValidateType   string     `json:"validateType"`
	ValidateInfo   string     `json:"validateInfo"`
	IssueTime      util.HTime `json:"issueTime"`
	ExpireTime     util.HTime `json:"expireTime"`
	DaysLeft       int        `json:"daysLeft"`
	IsValid        bool       `json:"isValid"`
	ErrorMsg       string     `json:"errorMsg"`
	CreateTime     util.HTime `json:"createTime"`
	UpdateTime     util.HTime `json:"updateTime"`
}

// GetCertSourceText 获取证书来源文本
func (s *SslCert) GetCertSourceText() string {
	sourceMap := map[string]string{
		CertSourceAliyunCAS: "阿里云CAS",
		CertSourceACME:      "ACME协议",
	}
	if text, ok := sourceMap[s.CertSource]; ok {
		return text
	}
	return "未知"
}

// ToVO 转换为VO
func (s *SslCert) ToVO() *SslCertVO {
	return &SslCertVO{
		ID:             s.ID,
		Domain:         s.Domain,
		CertSource:     s.CertSource,
		CertSourceText: s.GetCertSourceText(),
		CAProvider:     s.CAProvider,
		AliyunConfigID: s.AliyunConfigID,
		OrderID:        s.OrderID,
		CertID:         s.CertID,
		CertName:       s.CertName,
		ProductCode:    s.ProductCode,
		Algorithm:      s.Algorithm,
		Status:         s.Status,
		StatusText:     s.GetStatusText(),
		ValidateType:   s.ValidateType,
		ValidateInfo:   s.ValidateInfo,
		IssueTime:      s.IssueTime,
		ExpireTime:     s.ExpireTime,
		DaysLeft:       s.DaysLeft,
		IsValid:        s.IsValid(),
		ErrorMsg:       s.ErrorMsg,
		CreateTime:     s.CreateTime,
		UpdateTime:     s.UpdateTime,
	}
}

// CreateAliyunConfigDto 创建阿里云配置DTO
type CreateAliyunConfigDto struct {
	Name         string `validate:"required" json:"name"`         // 配置名称
	AccessKey    string `validate:"required" json:"accessKey"`    // AccessKey
	AccessSecret string `validate:"required" json:"accessSecret"` // AccessSecret
	Region       string `json:"region"`                           // 区域
	Email        string `validate:"required,email" json:"email"`  // 联系邮箱
	Phone        string `json:"phone"`                            // 联系电话
	Username     string `json:"username"`                         // 联系人姓名
}

// UpdateAliyunConfigDto 更新阿里云配置DTO
type UpdateAliyunConfigDto struct {
	ID           uint   `validate:"required" json:"id"`           // ID
	Name         string `validate:"required" json:"name"`         // 配置名称
	AccessKey    string `validate:"required" json:"accessKey"`    // AccessKey
	AccessSecret string `validate:"required" json:"accessSecret"` // AccessSecret
	Region       string `json:"region"`                           // 区域
	Email        string `validate:"required,email" json:"email"`  // 联系邮箱
	Phone        string `json:"phone"`                            // 联系电话
	Username     string `json:"username"`                         // 联系人姓名
	Status       int    `json:"status"`                           // 状态
}

// AliyunConfigVO 阿里云配置视图对象
type AliyunConfigVO struct {
	ID           uint       `json:"id"`
	Name         string     `json:"name"`
	AccessKey    string     `json:"accessKey"`
	AccessSecret string     `json:"accessSecret"` // 敏感信息，可选择性脱敏
	Region       string     `json:"region"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone"`
	Username     string     `json:"username"`
	Status       int        `json:"status"`
	CreateTime   util.HTime `json:"createTime"`
	UpdateTime   util.HTime `json:"updateTime"`
}

// ToVO 转换为VO
func (a *AliyunConfig) ToVO(maskSecret bool) *AliyunConfigVO {
	vo := &AliyunConfigVO{
		ID:           a.ID,
		Name:         a.Name,
		AccessKey:    a.AccessKey,
		AccessSecret: a.AccessSecret,
		Region:       a.Region,
		Email:        a.Email,
		Phone:        a.Phone,
		Username:     a.Username,
		Status:       a.Status,
		CreateTime:   a.CreateTime,
		UpdateTime:   a.UpdateTime,
	}

	// 脱敏处理
	if maskSecret && len(a.AccessSecret) > 8 {
		vo.AccessSecret = a.AccessSecret[:4] + "****" + a.AccessSecret[len(a.AccessSecret)-4:]
	}

	return vo
}

// DeploySslCertRequest 部署SSL证书请求
type DeploySslCertRequest struct {
	CertID     uint   `json:"certId" validate:"required"`     // 证书ID
	HostID     uint   `json:"hostId" validate:"required"`     // 目标主机ID
	DeployPath string `json:"deployPath" validate:"required"` // 部署路径，如: /etc/nginx/ssl
}

// DeploySslCertResponse 部署SSL证书响应
type DeploySslCertResponse struct {
	Success      bool     `json:"success"`      // 是否成功
	Message      string   `json:"message"`      // 消息
	DeployLogID  uint     `json:"deployLogId"`  // 部署日志ID，前端可通过此ID查询部署进度
	BackupFiles  []string `json:"backupFiles"`  // 备份的文件列表
	DeployFiles  []string `json:"deployFiles"`  // 部署的文件列表
	Logs         []string `json:"logs"`         // 执行日志
}

// SslCertDeployLog SSL证书部署日志表
type SslCertDeployLog struct {
	ID          uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	CertID      uint       `gorm:"column:cert_id;comment:'证书ID';NOT NULL;index" json:"certId"`
	Domain      string     `gorm:"column:domain;varchar(255);comment:'域名'" json:"domain"`
	HostID      uint       `gorm:"column:host_id;comment:'目标主机ID';NOT NULL;index" json:"hostId"`
	HostName    string     `gorm:"column:host_name;varchar(255);comment:'主机名称'" json:"hostName"`
	DeployPath  string     `gorm:"column:deploy_path;varchar(500);comment:'部署路径'" json:"deployPath"`
	Status      int        `gorm:"column:status;comment:'部署状态:1->部署中,2->成功,3->失败';default:1" json:"status"`
	BackupFiles string     `gorm:"column:backup_files;type:text;comment:'备份的文件列表(JSON)'" json:"backupFiles"`
	DeployFiles string     `gorm:"column:deploy_files;type:text;comment:'部署的文件列表(JSON)'" json:"deployFiles"`
	Logs        string     `gorm:"column:logs;type:text;comment:'执行日志(JSON)'" json:"logs"`
	ErrorMsg    string     `gorm:"column:error_msg;type:text;comment:'错误信息'" json:"errorMsg"`
	CreateTime  util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime  util.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
}

func (SslCertDeployLog) TableName() string {
	return "monitor_ssl_cert_deploy_log"
}

// GetDeployLogsRequest 获取部署日志列表请求
type GetDeployLogsRequest struct {
	CertID   uint   `form:"certId"`   // 证书ID(可选)
	HostID   uint   `form:"hostId"`   // 主机ID(可选)
	Domain   string `form:"domain"`   // 域名(模糊查询,可选)
	Status   *int   `form:"status"`   // 部署状态(可选)
	Page     int    `form:"page"`     // 页码
	PageSize int    `form:"pageSize"` // 每页数量
}

// GetDeployLogsResponse 获取部署日志列表响应
type GetDeployLogsResponse struct {
	Total int                `json:"total"` // 总数
	List  []DeployLogVO      `json:"list"`  // 列表
}

// DeployLogVO 部署日志视图对象
type DeployLogVO struct {
	ID          uint     `json:"id"`          // ID
	CertID      uint     `json:"certId"`      // 证书ID
	Domain      string   `json:"domain"`      // 域名
	HostID      uint     `json:"hostId"`      // 主机ID
	HostName    string   `json:"hostName"`    // 主机名称
	DeployPath  string   `json:"deployPath"`  // 部署路径
	Status      int      `json:"status"`      // 部署状态
	StatusText  string   `json:"statusText"`  // 状态文本
	BackupFiles []string `json:"backupFiles"` // 备份文件列表
	DeployFiles []string `json:"deployFiles"` // 部署文件列表
	Logs        []string `json:"logs"`        // 执行日志
	ErrorMsg    string   `json:"errorMsg"`    // 错误信息
	CreateTime  string   `json:"createTime"`  // 创建时间
	UpdateTime  string   `json:"updateTime"`  // 更新时间
}

// UploadSslCertRequest 上传SSL证书请求
type UploadSslCertRequest struct {
	Domain   string `form:"domain" validate:"required"`   // 域名
	CertFile string `form:"certFile"`                     // 证书文件内容(可选,从压缩包解析)
	KeyFile  string `form:"keyFile"`                      // 私钥文件内容(可选,从压缩包解析)
}

// UploadSslCertResponse 上传SSL证书响应
type UploadSslCertResponse struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
	CertID  uint   `json:"certId"`  // 证书ID
	Domain  string `json:"domain"`  // 域名
}
