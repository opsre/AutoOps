package model

import (
	"database/sql"
)

// DomainSchedule 域名监控定时任务配置
type DomainSchedule struct {
	ID              uint         `gorm:"primarykey" json:"id"`
	Enabled         bool         `gorm:"type:tinyint(1);default:0;comment:是否启用" json:"enabled"`
	CronExpr        string       `gorm:"type:varchar(100);comment:Cron表达式" json:"cronExpr"`
	NotifyEnabled   bool         `gorm:"type:tinyint(1);default:0;comment:是否开启通知" json:"notifyEnabled"`
	NotifyRobotId   *uint        `gorm:"type:int;comment:通知机器人ID" json:"notifyRobotId"`
	NextRunAt       sql.NullTime `gorm:"type:datetime;comment:下次执行时间" json:"nextRunAt"`
	LastRunAt       sql.NullTime `gorm:"type:datetime;comment:上次执行时间" json:"lastRunAt"`
	Status          string       `gorm:"type:varchar(50);comment:任务状态(running/paused/stopped)" json:"status"`
	CreateTime      sql.NullTime `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime      sql.NullTime `gorm:"autoUpdateTime" json:"updateTime"`
}

// TableName 指定表名
func (DomainSchedule) TableName() string {
	return "monitor_domain_schedule"
}

// DomainScheduleDto 域名监控定时任务配置DTO
type DomainScheduleDto struct {
	Enabled       bool   `json:"enabled" binding:"required"`   // 是否启用
	CronExpr      string `json:"cronExpr" binding:"required"`  // Cron表达式
	NotifyEnabled bool   `json:"notifyEnabled"`                // 是否开启通知
	NotifyRobotId *uint  `json:"notifyRobotId"`                // 通知机器人ID
}
