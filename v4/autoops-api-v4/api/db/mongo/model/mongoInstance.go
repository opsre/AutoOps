package model

import "dodevops-api/common/util"

// MongoInstance MongoDB 实例模型
type MongoInstance struct {
	ID                 uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Code               string     `gorm:"column:code;type:varchar(36);comment:'实例编码';uniqueIndex;NOT NULL" json:"code"`
	Name               string     `gorm:"column:name;type:varchar(100);comment:'实例名称';NOT NULL" json:"name"`
	Type               string     `gorm:"column:type;type:varchar(20);comment:'MongoDB类型:mongodb,mongodb-atlas等';default:'mongodb'" json:"type"`
	Uri                string     `gorm:"column:uri;type:varchar(500);comment:'连接URI';NOT NULL" json:"uri"`
	SshTunnelMachineId uint       `gorm:"column:ssh_tunnel_machine_id;comment:'SSH隧道机器ID';default:0" json:"sshTunnelMachineId"`
	Remark             string     `gorm:"column:remark;type:varchar(500);comment:'备注'" json:"remark"`
	Status             int        `gorm:"column:status;comment:'状态:1->启用,2->禁用';default:1" json:"status"`
	CreateTime         util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime         util.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
	Creator            string     `gorm:"column:creator;type:varchar(64);comment:'创建人'" json:"creator"`
	CreatorId          uint       `gorm:"column:creator_id;comment:'创建人ID'" json:"creatorId"`
	Modifier           string     `gorm:"column:modifier;type:varchar(64);comment:'修改人'" json:"modifier"`
	ModifierId         uint       `gorm:"column:modifier_id;comment:'修改人ID'" json:"modifierId"`
}

func (MongoInstance) TableName() string { return "db_mongo_instance" }

// MongoInstanceVO 返回给前端的视图对象
type MongoInstanceVO struct {
	ID                 uint       `json:"id"`
	Code               string     `json:"code"`
	Name               string     `json:"name"`
	Type               string     `json:"type"`
	Host               string     `json:"host"`
	Port               int        `json:"port"`
	AuthDatabase       string     `json:"authDatabase"`
	Username           string     `json:"username"`
	Password           string     `json:"password"` // 脱敏后的密码
	SshTunnelMachineId uint       `json:"sshTunnelMachineId"`
	Remark             string     `json:"remark"`
	Status             int        `json:"status"`
	CreateTime         util.HTime `json:"createTime"`
	UpdateTime         util.HTime `json:"updateTime"`
	Creator            string     `json:"creator"`
	CreatorId          uint       `json:"creatorId"`
	Modifier           string     `json:"modifier"`
	ModifierId         uint       `json:"modifierId"`
}

// MongoInstanceQuery 查询条件
type MongoInstanceQuery struct {
	Name   string `json:"name"`
	Status int    `json:"status"`

	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

// MongoInstanceDto 创建/更新 DTO
type MongoInstanceDto struct {
	ID                 uint   `json:"id"`
	Code               string `json:"code"`
	Name               string `json:"name"`
	Type               string `json:"type"`
	Host               string `json:"host"`
	Port               int    `json:"port"`
	AuthDatabase       string `json:"authDatabase"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	Uri                string `json:"uri"`
	SshTunnelMachineId uint   `json:"sshTunnelMachineId"`
	Remark             string `json:"remark"`
}

// MongoDB 操作表单
type MongoFindForm struct {
	Database   string         `json:"database" binding:"required"`
	Collection string         `json:"collection" binding:"required"`
	Filter     map[string]any `json:"filter"`
	Sort       map[string]any `json:"sort"`
	Skip       int64          `json:"skip"`
	Limit      int64          `json:"limit"`
}

type MongoUpdateByIdForm struct {
	Database   string         `json:"database" binding:"required"`
	Collection string         `json:"collection" binding:"required"`
	Filter     map[string]any `json:"filter"` // 从 filter._id 获取
	DocId      any            `json:"docId"`  // 或直接传 docId
	ID         any            `json:"_id"`    // 或直接传 _id
	Update     map[string]any `json:"update" binding:"required"`
}

type MongoDeleteByIdForm struct {
	Database   string         `json:"database" binding:"required"`
	Collection string         `json:"collection" binding:"required"`
	Filter     map[string]any `json:"filter"` // 从 filter._id 获取
	DocId      any            `json:"docId"`  // 或直接传 docId
	ID         any            `json:"_id"`    // 或直接传 _id
}

type MongoInsertForm struct {
	Database   string         `json:"database" binding:"required"`
	Collection string         `json:"collection" binding:"required"`
	Doc        map[string]any `json:"doc"`      // 支持 doc 字段
	Document   map[string]any `json:"document"` // 支持 document 字段
}

type MongoRunCommandForm struct {
	Database string           `json:"database" binding:"required"`
	Command  []map[string]any `json:"command" binding:"required"`
}
