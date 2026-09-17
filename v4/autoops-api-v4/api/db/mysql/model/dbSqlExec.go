package model

import (
	"dodevops-api/common/util"
)

// DbSqlExec SQL执行历史模型
type DbSqlExec struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	DbId       uint       `gorm:"column:db_id;comment:'数据库ID';index;NOT NULL" json:"dbId"`
	DbName     string     `gorm:"column:db_name;type:varchar(100);comment:'数据库名';NOT NULL" json:"dbName"`
	TblName    string     `gorm:"column:table_name;type:varchar(100);comment:'表名'" json:"tableName"`
	Type       int8       `gorm:"column:type;comment:'类型:1->查询,2->插入,3->更新,4->删除,5->DDL';NOT NULL" json:"type"`
	Sql        string     `gorm:"column:sql;type:text;comment:'SQL语句';NOT NULL" json:"sql"`
	OldValue   string     `gorm:"column:old_value;type:longtext;comment:'旧值(用于回滚)'" json:"oldValue"`
	Remark     string     `gorm:"column:remark;type:varchar(500);comment:'备注'" json:"remark"`
	Status     int8       `gorm:"column:status;comment:'状态:1->成功,2->失败';NOT NULL" json:"status"`
	Res        string     `gorm:"column:res;type:text;comment:'执行结果'" json:"res"`
	ExecTime   int64      `gorm:"column:exec_time;comment:'执行时长(ms)'" json:"execTime"`
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';index;NOT NULL" json:"createTime"`
	Creator    string     `gorm:"column:creator;type:varchar(64);comment:'创建人'" json:"creator"`
	CreatorId  uint       `gorm:"column:creator_id;comment:'创建人ID'" json:"creatorId"`
}

func (DbSqlExec) TableName() string {
	return "db_sql_exec"
}

// SQL类型常量
const (
	SqlTypeQuery  int8 = 1 // 查询
	SqlTypeInsert int8 = 2 // 插入
	SqlTypeUpdate int8 = 3 // 更新
	SqlTypeDelete int8 = 4 // 删除
	SqlTypeDDL    int8 = 5 // DDL
)

// 执行状态常量
const (
	ExecStatusSuccess int8 = 1 // 成功
	ExecStatusFail    int8 = 2 // 失败
)

// DbSqlExecDto SQL执行DTO
type DbSqlExecDto struct {
    DbId       uint   `json:"dbId"`        // 兼容旧版，可选
    InstanceId uint   `json:"instanceId"`  // 数据库实例ID，优先使用
    Db         string `json:"db" validate:"required"`
    Sql        string `json:"sql" validate:"required"`
    Remark     string `json:"remark"`
    // 可选：查询结果缓存开关与过期时间（秒），仅对查询类SQL生效
    UseCache     bool `json:"useCache"`
    CacheSeconds int  `json:"cacheSeconds"`
    // 可选：查询最大返回行数（0=不限制），超时时间（秒，0=默认60s）
    MaxRows        int `json:"maxRows"`
    TimeoutSeconds int `json:"timeoutSeconds"`
}

// DbSqlExecQuery 查询条件
type DbSqlExecQuery struct {
	DbId       uint   `json:"dbId"`
	DbName     string `json:"dbName"`
	TableName  string `json:"tableName"`
	Type       int8   `json:"type"`
	Status     int8   `json:"status"`
	Creator    string `json:"creator"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`

	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

// DbSqlExecResult SQL执行结果
type DbSqlExecResult struct {
    Columns      []string                 `json:"columns"`      // 列名
    Rows         []map[string]interface{} `json:"rows"`         // 数据行
    AffectedRows int64                    `json:"affectedRows"` // 影响行数
    ExecTime     int64                    `json:"execTime"`     // 执行时长(ms)
    ColumnMetas  []ColumnMeta             `json:"columnMetas"`  // 列元信息(名称/类型/备注)，可选
}

// ColumnMeta 查询结果列的元信息
type ColumnMeta struct {
    Name    string `json:"name"`
    Type    string `json:"type"`
    Comment string `json:"comment"`
}
