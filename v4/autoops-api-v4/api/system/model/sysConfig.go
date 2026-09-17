// 系统配置相关结构体
package model

import (
	"dodevops-api/common/util"
)

// SysConfig 系统配置模型
type SysConfig struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	ConfigKey  string     `gorm:"column:config_key;type:varchar(100);uniqueIndex;comment:'配置键';NOT NULL" json:"configKey"`
	ConfigType string     `gorm:"column:config_type;type:varchar(50);comment:'配置类型(ldap,email,sms等)';NOT NULL" json:"configType"`
	ConfigData string     `gorm:"column:config_data;type:text;comment:'配置数据(JSON格式)';NOT NULL" json:"configData"`
	Status     int        `gorm:"column:status;default:1;comment:'状态:1->启用,2->禁用';NOT NULL" json:"status"`
	Remark     string     `gorm:"column:remark;type:varchar(500);comment:'备注'" json:"remark"`
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}

// LDAPConfig LDAP配置结构
type LDAPConfig struct {
	Enable          bool           `json:"enable"`          // 是否启用
	Host            string         `json:"host"`            // LDAP服务器地址
	Port            int            `json:"port"`            // LDAP端口
	BaseDn          string         `json:"baseDn"`          // BaseDN
	BindUser        string         `json:"bindUser"`        // 绑定用户
	BindPass        string         `json:"bindPass"`        // 绑定密码
	AuthFilter      string         `json:"authFilter"`      // 认证过滤器
	CoverAttributes bool           `json:"coverAttributes"` // 是否覆盖属性
	TLS             bool           `json:"tls"`             // 是否使用TLS
	StartTLS        bool           `json:"startTLS"`        // 是否使用StartTLS
	DefaultRoles    []string       `json:"defaultRoles"`    // 默认角色（已废弃，使用DefaultRoleID）
	DefaultRoleID   uint           `json:"defaultRoleId"`   // 默认角色ID
	Attributes      LDAPAttributes `json:"attributes"`      // 属性映射
}

// LDAPAttributes LDAP属性映射
type LDAPAttributes struct {
	Nickname string `json:"nickname"` // 昵称字段映射
	Phone    string `json:"phone"`    // 电话字段映射
	Email    string `json:"email"`    // 邮箱字段映射
}

// SysConfigDto 系统配置DTO
type SysConfigDto struct {
	ID         uint   `json:"id"`
	ConfigKey  string `json:"configKey" validate:"required"`  // 配置键
	ConfigType string `json:"configType" validate:"required"` // 配置类型
	ConfigData string `json:"configData" validate:"required"` // 配置数据(JSON)
	Status     int    `json:"status"`                         // 状态
	Remark     string `json:"remark"`                         // 备注
}

// LDAPConfigDto LDAP配置DTO
type LDAPConfigDto struct {
	Enable          bool           `json:"enable"`                                   // 是否启用
	Host            string         `json:"host" validate:"required"`                 // LDAP服务器地址
	Port            int            `json:"port" validate:"required,min=1,max=65535"` // LDAP端口
	BaseDn          string         `json:"baseDn" validate:"required"`               // BaseDN
	BindUser        string         `json:"bindUser" validate:"required"`             // 绑定用户
	BindPass        string         `json:"bindPass" validate:"required"`             // 绑定密码
	AuthFilter      string         `json:"authFilter" validate:"required"`           // 认证过滤器
	CoverAttributes bool           `json:"coverAttributes"`                          // 是否覆盖属性
	TLS             bool           `json:"tls"`                                      // 是否使用TLS
	StartTLS        bool           `json:"startTLS"`                                 // 是否使用StartTLS
	DefaultRoles    []string       `json:"defaultRoles"`                             // 默认角色（已废弃）
	DefaultRoleID   uint           `json:"defaultRoleId"`                            // 默认角色ID
	Attributes      LDAPAttributes `json:"attributes"`                               // 属性映射
	Remark          string         `json:"remark"`                                   // 备注
}

// SysConfigQuery 系统配置查询条件
type SysConfigQuery struct {
	ConfigKey  string `form:"configKey"`  // 配置键
	ConfigType string `form:"configType"` // 配置类型
	Status     int    `form:"status"`     // 状态(-1:全部,1:启用,2:禁用)
	PageNum    int    `form:"pageNum"`    // 页码
	PageSize   int    `form:"pageSize"`   // 每页数量
}
