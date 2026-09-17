package model

import (
	"dodevops-api/common/util"
)

// DbInstance 数据库实例模型
type DbInstance struct {
	ID                 uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Code               string     `gorm:"column:code;type:varchar(36);comment:'实例编码';uniqueIndex;NOT NULL" json:"code"`
	Name               string     `gorm:"column:name;type:varchar(100);comment:'实例名称';NOT NULL" json:"name"`
	Type               string     `gorm:"column:type;type:varchar(20);comment:'数据库类型:mysql,postgres,oracle等';NOT NULL" json:"type"`
	Host               string     `gorm:"column:host;type:varchar(100);comment:'主机地址';NOT NULL" json:"host"`
	Port               int        `gorm:"column:port;comment:'端口';NOT NULL" json:"port"`
	Network            string     `gorm:"column:network;type:varchar(20);comment:'网络类型:tcp,unix';default:'tcp'" json:"network"`
	Params             string     `gorm:"column:params;type:varchar(500);comment:'连接参数'" json:"params"`
	Username           string     `gorm:"column:username;type:varchar(100);comment:'用户名';NOT NULL" json:"username"`
	Password           string     `gorm:"column:password;type:varchar(500);comment:'密码(加密)';NOT NULL" json:"password"`
	Remark             string     `gorm:"column:remark;type:varchar(500);comment:'备注'" json:"remark"`
	SshTunnelMachineId uint       `gorm:"column:ssh_tunnel_machine_id;comment:'SSH隧道机器ID';default:0" json:"sshTunnelMachineId"`
	Status             int        `gorm:"column:status;comment:'状态:1->启用,2->禁用';default:1" json:"status"`
	CreateTime         util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime         util.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
	Creator            string     `gorm:"column:creator;type:varchar(64);comment:'创建人'" json:"creator"`
	CreatorId          uint       `gorm:"column:creator_id;comment:'创建人ID'" json:"creatorId"`
	Modifier           string     `gorm:"column:modifier;type:varchar(64);comment:'修改人'" json:"modifier"`
	ModifierId         uint       `gorm:"column:modifier_id;comment:'修改人ID'" json:"modifierId"`
}

func (DbInstance) TableName() string {
	return "db_instance"
}

// DbInstanceQuery 查询条件
type DbInstanceQuery struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Host   string `json:"host"`
	Status int    `json:"status"`

	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

// DbInstanceDto 实例DTO
type DbInstanceDto struct {
	ID                 uint   `json:"id"`
	Code               string `json:"code"`
	Name               string `json:"name" validate:"required"`
	Type               string `json:"type" validate:"required"`
	Host               string `json:"host" validate:"required"`
	Port               int    `json:"port" validate:"required"`
	Network            string `json:"network"`
	Params             string `json:"params"`
	Username           string `json:"username" validate:"required"`
	Password           string `json:"password"`
	Remark             string `json:"remark"`
	SshTunnelMachineId uint   `json:"sshTunnelMachineId"`
}

// TestConnDto 测试连接DTO
type TestConnDto struct {
    InstanceId uint   `json:"instanceId"` // 可选：若提供则使用已保存实例信息
    Type     string `json:"type" validate:"required"`
    Host     string `json:"host" validate:"required"`
    Port     int    `json:"port" validate:"required"`
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
    Database string `json:"database"`
    Network  string `json:"network"`
    Params   string `json:"params"`
}
