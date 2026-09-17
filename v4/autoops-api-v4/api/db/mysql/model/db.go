package model

import (
	"dodevops-api/common/util"
)

// Db 数据库(库)模型
type Db struct {
	ID           uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Code         string     `gorm:"column:code;type:varchar(36);comment:'数据库编码';uniqueIndex;NOT NULL" json:"code"`
	Name         string     `gorm:"column:name;type:varchar(100);comment:'数据库名称';NOT NULL" json:"name"`
	Database     string     `gorm:"column:database;type:varchar(500);comment:'数据库名(多个用空格分隔)';NOT NULL" json:"database"`
	Remark       string     `gorm:"column:remark;type:varchar(500);comment:'备注'" json:"remark"`
	InstanceId   uint       `gorm:"column:instance_id;comment:'实例ID';index;NOT NULL" json:"instanceId"`
	InstanceCode string     `gorm:"column:instance_code;type:varchar(36);comment:'实例编码'" json:"instanceCode"`
	Status       int        `gorm:"column:status;comment:'状态:1->启用,2->禁用';default:1" json:"status"`
	CreateTime   util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime   util.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
	Creator      string     `gorm:"column:creator;type:varchar(64);comment:'创建人'" json:"creator"`
	CreatorId    uint       `gorm:"column:creator_id;comment:'创建人ID'" json:"creatorId"`
	Modifier     string     `gorm:"column:modifier;type:varchar(64);comment:'修改人'" json:"modifier"`
	ModifierId   uint       `gorm:"column:modifier_id;comment:'修改人ID'" json:"modifierId"`
}

func (Db) TableName() string {
	return "db"
}

// DbQuery 查询条件
type DbQuery struct {
	Name       string `json:"name"`
	Database   string `json:"database"`
	InstanceId uint   `json:"instanceId"`
	Status     int    `json:"status"`

	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

// DbDto 数据库DTO
type DbDto struct {
	ID         uint   `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name" validate:"required"`
	Database   string `json:"database" validate:"required"`
	Remark     string `json:"remark"`
	InstanceId uint   `json:"instanceId" validate:"required"`
}

// DbListVo 数据库列表VO
type DbListVo struct {
	ID           uint       `json:"id"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	Database     string     `json:"database"`
	Remark       string     `json:"remark"`
	InstanceId   uint       `json:"instanceId"`
	InstanceCode string     `json:"instanceCode"`
	InstanceName string     `json:"instanceName"`
	InstanceType string     `json:"instanceType"`
	Host         string     `json:"host"`
	Port         int        `json:"port"`
	Status       int        `json:"status"`
	CreateTime   util.HTime `json:"createTime"`
	Creator      string     `json:"creator"`
}
