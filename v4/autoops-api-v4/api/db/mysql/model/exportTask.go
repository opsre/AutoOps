package model

import (
	"dodevops-api/common/util"
	"time"
)

// ExportTaskStatus 导出任务状态
type ExportTaskStatus string

const (
	ExportTaskStatusPending    ExportTaskStatus = "pending"    // 等待中
	ExportTaskStatusProcessing ExportTaskStatus = "processing" // 处理中
	ExportTaskStatusCompleted  ExportTaskStatus = "completed"  // 已完成
	ExportTaskStatusFailed     ExportTaskStatus = "failed"     // 失败
)

// ExportTask 导出任务
type ExportTask struct {
	ID           uint             `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	TaskId       string           `gorm:"column:task_id;type:varchar(36);comment:'任务ID';uniqueIndex;NOT NULL" json:"taskId"`
	DbId         uint             `gorm:"column:db_id;comment:'数据库ID';NOT NULL" json:"dbId"`
	DbName       string           `gorm:"column:db_name;type:varchar(100);comment:'数据库名';NOT NULL" json:"dbName"`
	ExportType   string           `gorm:"column:export_type;type:varchar(20);comment:'导出类型:structure/full';NOT NULL" json:"exportType"`
	Status       ExportTaskStatus `gorm:"column:status;type:varchar(20);comment:'状态';NOT NULL" json:"status"`
	FilePath     string           `gorm:"column:file_path;type:varchar(500);comment:'文件路径'" json:"filePath"`
	FileSize     int64            `gorm:"column:file_size;comment:'文件大小(字节)'" json:"fileSize"`
	ErrorMessage string           `gorm:"column:error_message;type:text;comment:'错误信息'" json:"errorMessage"`
	StartTime    *time.Time       `gorm:"column:start_time;comment:'开始时间'" json:"startTime"`
	EndTime      *time.Time       `gorm:"column:end_time;comment:'结束时间'" json:"endTime"`
	CreateTime   util.HTime       `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime   util.HTime       `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
	Creator      string           `gorm:"column:creator;type:varchar(64);comment:'创建人'" json:"creator"`
	CreatorId    uint             `gorm:"column:creator_id;comment:'创建人ID'" json:"creatorId"`
}

func (ExportTask) TableName() string {
	return "db_export_task"
}

// ExportTaskVo 导出任务VO
type ExportTaskVo struct {
	TaskId       string           `json:"taskId"`
	DbName       string           `json:"dbName"`
	ExportType   string           `json:"exportType"`
	Status       ExportTaskStatus `json:"status"`
	FileSize     int64            `json:"fileSize"`
	ErrorMessage string           `json:"errorMessage,omitempty"`
	StartTime    *time.Time       `json:"startTime,omitempty"`
	EndTime      *time.Time       `json:"endTime,omitempty"`
	CreateTime   util.HTime       `json:"createTime"`
}
