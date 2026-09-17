package model

import (
	"dodevops-api/common/util"
)

// DbSql SQL模板模型
type DbSql struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	DbId       uint       `gorm:"column:db_id;comment:'数据库ID';index;NOT NULL" json:"dbId"`
	Db         string     `gorm:"column:db;type:varchar(100);comment:'数据库名';NOT NULL" json:"db"`
	Name       string     `gorm:"column:name;type:varchar(100);comment:'SQL名称';NOT NULL" json:"name"`
	Type       int        `gorm:"column:type;comment:'类型:1->查询,2->更新';default:1" json:"type"`
	Sql        string     `gorm:"column:sql;type:text;comment:'SQL语句';NOT NULL" json:"sql"`
	Remark     string     `gorm:"column:remark;type:varchar(500);comment:'备注'" json:"remark"`
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime util.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
	Creator    string     `gorm:"column:creator;type:varchar(64);comment:'创建人'" json:"creator"`
	CreatorId  uint       `gorm:"column:creator_id;comment:'创建人ID'" json:"creatorId"`
	Modifier   string     `gorm:"column:modifier;type:varchar(64);comment:'修改人'" json:"modifier"`
	ModifierId uint       `gorm:"column:modifier_id;comment:'修改人ID'" json:"modifierId"`
}

func (DbSql) TableName() string {
	return "db_sql"
}

// DbSqlDto SQL模板DTO
type DbSqlDto struct {
	ID     uint   `json:"id"`
	DbId   uint   `json:"dbId" validate:"required"`
	Db     string `json:"db" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Type   int    `json:"type"`
	Sql    string `json:"sql" validate:"required"`
	Remark string `json:"remark"`
}

// DbSqlQuery 查询条件
type DbSqlQuery struct {
	DbId uint   `json:"dbId"`
	Db   string `json:"db"`
	Name string `json:"name"`

	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}
