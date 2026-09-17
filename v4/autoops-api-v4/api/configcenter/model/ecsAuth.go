// ECS认证凭证模型
// xiaoRui
package model

import (
	"dodevops-api/common/util"
)

// ECS认证凭证模型
type EcsAuth struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Name       string     `gorm:"column:name;varchar(64);comment:'凭证名称';NOT NULL" json:"name"`
	Type       int        `gorm:"column:type;comment:'认证类型:1->密码,2->私钥,3->公钥(免认证),4->跳板机';NOT NULL" json:"type"`
	Username   string     `gorm:"column:username;varchar(64);comment:'用户名'" json:"username"`
	Password   string     `gorm:"column:password;varchar(256);comment:'密码(type=1时使用)'" json:"password"`
	PublicKey  string     `gorm:"column:public_key;type:text;comment:'私钥内容(type=2时使用，字段名历史原因)'" json:"publicKey"` // 实际存储私钥
	Port       int        `gorm:"column:port;comment:'端口号';default:22" json:"port"`
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	Remark     string     `gorm:"column:remark;varchar(500);comment:'备注'" json:"remark"`

	// 跳板机相关字段 (type=4时使用)
	BastionHost       string `gorm:"column:bastion_host;varchar(64);comment:'跳板机IP'" json:"bastionHost"`
	BastionPort       int    `gorm:"column:bastion_port;comment:'跳板机端口';default:22" json:"bastionPort"`
	BastionUsername   string `gorm:"column:bastion_username;varchar(64);comment:'跳板机用户名'" json:"bastionUsername"`
	BastionAuthType   int    `gorm:"column:bastion_auth_type;comment:'跳板机认证类型:1->密码,2->私钥,3->免密'" json:"bastionAuthType"`
	BastionPassword   string `gorm:"column:bastion_password;varchar(256);comment:'跳板机密码'" json:"bastionPassword"`
	BastionPrivateKey string `gorm:"column:bastion_private_key;type:text;comment:'跳板机私钥'" json:"bastionPrivateKey"`
	TargetAuthType    int    `gorm:"column:target_auth_type;comment:'目标机器认证类型:1->密码,2->私钥,3->免密'" json:"targetAuthType"`
	TargetPassword    string `gorm:"column:target_password;varchar(256);comment:'目标机器密码'" json:"targetPassword"`
	TargetPrivateKey  string `gorm:"column:target_private_key;type:text;comment:'目标机器私钥'" json:"targetPrivateKey"`
}

func (EcsAuth) TableName() string {
	return "config_ecsauth"
}

// 创建ECS密码认证DTO
type CreateEcsPasswordAuthDto struct {
	Name      string `validate:"required" json:"name"` // 凭证名称
	Type      int    `validate:"required" json:"type"` // 认证类型:1->密码
	Username  string `validate:"required" json:"username"` // 用户名
	Password  string `json:"password"` // 密码
	PublicKey string `json:"publicKey"` // 公钥
	Port      int    `validate:"required" json:"port"` // 端口号
	Remark    string `json:"remark"` // 备注

	// 跳板机相关字段（type=4时使用）
	BastionHost       string `json:"bastionHost"`       // 跳板机IP
	BastionPort       int    `json:"bastionPort"`       // 跳板机端口
	BastionUsername   string `json:"bastionUsername"`   // 跳板机用户名
	BastionAuthType   int    `json:"bastionAuthType"`   // 跳板机认证类型
	BastionPassword   string `json:"bastionPassword"`   // 跳板机密码
	BastionPrivateKey string `json:"bastionPrivateKey"` // 跳板机私钥
	TargetAuthType    int    `json:"targetAuthType"`    // 目标机器认证类型
	TargetPassword    string `json:"targetPassword"`    // 目标机器密码
	TargetPrivateKey  string `json:"targetPrivateKey"`  // 目标机器私钥
}

// 创建ECS密钥认证DTO
type CreateEcsKeyAuthDto struct {
	Name      string `validate:"required"` // 凭证名称
	Type      int    `validate:"required"` // 认证类型:2->密钥
	PublicKey string `validate:"required"` // 私钥内容(字段名历史原因，实际传入私钥)
	Username  string `validate:"required"` // 用户名
	Port      int    `validate:"required"` // 端口号
	Remark    string // 备注
}

// 创建ECS公钥认证DTO (type=3, 免认证)
type CreateEcsPublicKeyAuthDto struct {
	Name     string `validate:"required"` // 凭证名称
	Type     int    `validate:"required"` // 认证类型:3->公钥免认证
	Username string `validate:"required"` // 用户名
	Port     int    `validate:"required"` // 端口号
	Remark   string // 备注
}

// 创建跳板机认证DTO (type=4)
type CreateEcsBastionAuthDto struct {
	Name              string `validate:"required" json:"name"`              // 凭证名称
	Type              int    `validate:"required" json:"type"`              // 认证类型:4->跳板机
	Username          string `validate:"required" json:"username"`          // 目标机器用户名
	Port              int    `validate:"required" json:"port"`              // 目标机器端口
	Remark            string `json:"remark"`                                // 备注

	// 跳板机配置
	BastionHost       string `validate:"required" json:"bastionHost"`       // 跳板机IP
	BastionPort       int    `validate:"required" json:"bastionPort"`       // 跳板机端口
	BastionUsername   string `validate:"required" json:"bastionUsername"`   // 跳板机用户名
	BastionAuthType   int    `validate:"required" json:"bastionAuthType"`   // 跳板机认证类型:1->密码,2->私钥,3->免密
	BastionPassword   string `json:"bastionPassword"`                       // 跳板机密码(bastionAuthType=1时必填)
	BastionPrivateKey string `json:"bastionPrivateKey"`                     // 跳板机私钥(bastionAuthType=2时必填)

	// 目标机器配置
	TargetAuthType    int    `validate:"required" json:"targetAuthType"`    // 目标机器认证类型:1->密码,2->私钥,3->免密
	TargetPassword    string `json:"targetPassword"`                        // 目标机器密码(targetAuthType=1时必填)
	TargetPrivateKey  string `json:"targetPrivateKey"`                      // 目标机器私钥(targetAuthType=2时必填)
}

// ID参数
type EcsAuthIdDto struct {
	Id uint `json:"id"` // ID
}

// 更新ECS认证DTO
type UpdateEcsAuthDto struct {
	EcsAuthIdDto
	CreateEcsPasswordAuthDto
}

// 认证凭证列表VO
type EcsAuthVo struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	Type       int        `json:"type"`
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	PublicKey  string     `json:"publicKey"`
	Port       int        `json:"port"`
	CreateTime util.HTime `json:"createTime"`
	Remark     string     `json:"remark"`

	// 跳板机相关字段
	BastionHost       string `json:"bastionHost"`
	BastionPort       int    `json:"bastionPort"`
	BastionUsername   string `json:"bastionUsername"`
	BastionAuthType   int    `json:"bastionAuthType"`
	BastionPassword   string `json:"bastionPassword"`
	BastionPrivateKey string `json:"bastionPrivateKey"`
	TargetAuthType    int    `json:"targetAuthType"`
	TargetPassword    string `json:"targetPassword"`
	TargetPrivateKey  string `json:"targetPrivateKey"`
}
