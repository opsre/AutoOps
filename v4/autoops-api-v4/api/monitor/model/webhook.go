package model

import "dodevops-api/common/util"

// WebhookRequest Webhook接收请求结构
type WebhookRequest struct {
	Source      string                 `json:"source" binding:"required"`      // 来源标识(如: zabbix/prometheus/alertmanager/custom)
	Title       string                 `json:"title" binding:"required"`       // 告警标题
	Content     string                 `json:"content" binding:"required"`     // 告警内容
	Level       string                 `json:"level"`                          // 告警级别(info/warning/error/critical)
	Tags        []string               `json:"tags"`                           // 标签
	Extra       map[string]interface{} `json:"extra"`                          // 额外数据
	NotifyRobotIDs []uint              `json:"notifyRobotIds"`                 // 指定通知的机器人ID列表(可选,不指定则通知所有启用的机器人)
}

// WebhookLog Webhook接收日志
type WebhookLog struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Source         string    `json:"source" gorm:"type:varchar(50);not null;comment:来源标识"`
	Title          string    `json:"title" gorm:"type:varchar(200);not null;comment:告警标题"`
	Content        string    `json:"content" gorm:"type:text;not null;comment:告警内容"`
	Level          string    `json:"level" gorm:"type:varchar(20);comment:告警级别"`
	Tags           string    `json:"tags" gorm:"type:varchar(500);comment:标签(JSON数组字符串)"`
	Extra          string    `json:"extra" gorm:"type:text;comment:额外数据(JSON字符串)"`
	NotifyRobotIDs string    `json:"notifyRobotIds" gorm:"type:varchar(200);comment:通知的机器人ID列表(逗号分隔)"`
	Status         string    `json:"status" gorm:"type:varchar(20);default:success;comment:处理状态(success/failed/partial)"`
	ErrorMsg       string    `json:"errorMsg" gorm:"type:text;comment:错误信息"`
	NotifyCount    int       `json:"notifyCount" gorm:"type:int;default:0;comment:通知数量"`
	SuccessCount   int       `json:"successCount" gorm:"type:int;default:0;comment:成功数量"`
	FailedCount    int       `json:"failedCount" gorm:"type:int;default:0;comment:失败数量"`
	CreatedAt      util.HTime `json:"createdAt" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (WebhookLog) TableName() string {
	return "monitor_webhook_log"
}

// WebhookNotifyLog Webhook通知日志详情
type WebhookNotifyLog struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	WebhookLogID uint      `json:"webhookLogId" gorm:"type:int;not null;index;comment:Webhook日志ID"`
	RobotID      uint      `json:"robotId" gorm:"type:int;not null;comment:机器人ID"`
	RobotName    string    `json:"robotName" gorm:"type:varchar(100);comment:机器人名称"`
	RobotType    string    `json:"robotType" gorm:"type:varchar(20);comment:机器人类型"`
	Status       string    `json:"status" gorm:"type:varchar(20);not null;comment:发送状态(success/failed)"`
	ErrorMsg     string    `json:"errorMsg" gorm:"type:text;comment:错误信息"`
	CreatedAt    util.HTime `json:"createdAt" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (WebhookNotifyLog) TableName() string {
	return "monitor_webhook_notify_log"
}
