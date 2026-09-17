package model

import (
	"database/sql"
)

// Domain 域名监控实体
type Domain struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Domain      string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"domain"` // 域名
	Tags        string         `gorm:"type:varchar(500)" json:"tags"`                        // 标签(逗号分隔)
	Remark      string         `gorm:"type:text" json:"remark"`                              // 备注信息
	Status      int            `gorm:"type:int;default:1;comment:状态(1:启用,0:禁用)" json:"status"`
	IsAlive     int            `gorm:"type:int;default:0;comment:存活状态(1:正常,0:异常)" json:"isAlive"`
	StatusCode  int            `gorm:"type:int;comment:HTTP状态码" json:"statusCode"`
	ResponseTime int           `gorm:"type:int;comment:响应时间(ms)" json:"responseTime"`
	SSLExpireAt sql.NullTime   `gorm:"type:datetime;comment:SSL证书过期时间" json:"sslExpireAt"`
	SSLDaysLeft int            `gorm:"type:int;comment:SSL证书剩余天数" json:"sslDaysLeft"`
	SSLIssuer   string         `gorm:"type:varchar(255);comment:SSL证书颁发者" json:"sslIssuer"`
	LastCheckAt sql.NullTime   `gorm:"type:datetime;comment:最后检查时间" json:"lastCheckAt"`
	ErrorMsg    string         `gorm:"type:text;comment:错误信息" json:"errorMsg"`
	CreateTime  sql.NullTime   `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime  sql.NullTime   `gorm:"autoUpdateTime" json:"updateTime"`
}

// TableName 指定表名
func (Domain) TableName() string {
	return "monitor_domain"
}

// DomainDto 域名监控DTO
type DomainDto struct {
	ID     uint   `json:"id"`
	Domain string `json:"domain" binding:"required"` // 域名
	Tags   string `json:"tags"`                      // 标签
	Remark string `json:"remark"`                    // 备注
	Status int    `json:"status"`                    // 状态
}

// DomainQuery 域名监控查询条件
type DomainQuery struct {
	Domain   string `form:"domain"`   // 域名(模糊查询)
	Tags     string `form:"tags"`     // 标签(模糊查询)
	Status   int    `form:"status"`   // 状态(-1:全部,1:启用,0:禁用)
	IsAlive  int    `form:"isAlive"`  // 存活状态(-1:全部,1:正常,0:异常)
	PageNum  int    `form:"pageNum"`  // 页码
	PageSize int    `form:"pageSize"` // 每页数量
}

// DomainCheckResult 域名检查结果
type DomainCheckResult struct {
	DomainID     uint   `json:"domainId"`
	IsAlive      int    `json:"isAlive"`      // 存活状态(1:正常,0:异常)
	StatusCode   int    `json:"statusCode"`   // HTTP状态码
	ResponseTime int    `json:"responseTime"` // 响应时间(ms)
	SSLExpireAt  string `json:"sslExpireAt"`  // SSL证书过期时间
	SSLDaysLeft  int    `json:"sslDaysLeft"`  // SSL证书剩余天数
	SSLIssuer    string `json:"sslIssuer"`    // SSL证书颁发者
	ErrorMsg     string `json:"errorMsg"`     // 错误信息
}

// DomainStatistics 域名监控统计
type DomainStatistics struct {
	Total      int64 `json:"total"`      // 总数
	Alive      int64 `json:"alive"`      // 正常数量
	Dead       int64 `json:"dead"`       // 异常数量
	SSLExpired int64 `json:"sslExpired"` // SSL证书已过期
	SSLSoon    int64 `json:"sslSoon"`    // SSL证书即将过期(30天内)
}

// ImportFromAliyunDNSDto 从阿里云解析导入域名DTO
type ImportFromAliyunDNSDto struct {
	KeyID            uint   `json:"keyId" binding:"required"`         // 密钥管理ID
	Tags             string `json:"tags"`                             // 统一标签(可选)
	Remark           string `json:"remark"`                           // 统一备注(可选)
	IncludeSubdomain bool   `json:"includeSubdomain"`                 // 是否包含子域名(默认false,只导入主域名)
}

// AliyunDNSImportResult 阿里云DNS导入结果
type AliyunDNSImportResult struct {
	Total       int      `json:"total"`       // 总域名数
	Imported    int      `json:"imported"`    // 成功导入数
	Skipped     int      `json:"skipped"`     // 跳过数(已存在)
	Failed      int      `json:"failed"`      // 失败数
	FailedList  []string `json:"failedList"`  // 失败的域名列表
}
