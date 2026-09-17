// 容量预测和告警相关模型
package model

import (
	"dodevops-api/common/util"
)

// CapacityPrediction 容量预测任务
type CapacityPrediction struct {
	ID           uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Name         string     `gorm:"column:name;type:varchar(200);comment:'预测任务名称';NOT NULL" json:"name"`
	TargetType   string     `gorm:"column:target_type;type:varchar(20);comment:'目标类型:host/group/cluster';NOT NULL" json:"targetType"`
	TargetID     uint       `gorm:"column:target_id;comment:'目标ID';NOT NULL;index" json:"targetId"`
	TargetName   string     `gorm:"column:target_name;type:varchar(200);comment:'目标名称';NOT NULL" json:"targetName"`
	MetricType   string     `gorm:"column:metric_type;type:varchar(50);comment:'指标类型:cpu/memory/disk/network/custom';NOT NULL;index" json:"metricType"`
	MetricName   string     `gorm:"column:metric_name;type:varchar(100);comment:'指标名称';NOT NULL" json:"metricName"`
	CurrentValue float64    `gorm:"column:current_value;type:decimal(10,2);comment:'当前值'" json:"currentValue"`
	Threshold    float64    `gorm:"column:threshold;type:decimal(10,2);comment:'阈值(百分比或绝对值)';NOT NULL" json:"threshold"`
	PredictDays  int        `gorm:"column:predict_days;comment:'预测天数';NOT NULL;default:30" json:"predictDays"`
	Algorithm    string     `gorm:"column:algorithm;type:varchar(20);comment:'算法:linear/polynomial/arima/lstm';default:'linear'" json:"algorithm"`
	Confidence   float64    `gorm:"column:confidence;type:decimal(3,2);comment:'置信度(0-1)'" json:"confidence"`
	PredictTime  util.HTime `gorm:"column:predict_time;comment:'预测到达阈值时间'" json:"predictTime"`
	DaysLeft     int        `gorm:"column:days_left;comment:'剩余天数'" json:"daysLeft"`
	TrendType    string     `gorm:"column:trend_type;type:varchar(20);comment:'趋势类型:increasing/decreasing/stable/fluctuating'" json:"trendType"`
	RiskLevel    int        `gorm:"column:risk_level;comment:'风险等级:0=安全,1=关注,2=预警,3=紧急';default:0;index" json:"riskLevel"`
	Suggestion   string     `gorm:"column:suggestion;type:text;comment:'AI建议'" json:"suggestion"`
	Status       int        `gorm:"column:status;comment:'状态:1=正常,2=预警,3=告警';default:1;index" json:"status"`
	Enabled      int        `gorm:"column:enabled;comment:'是否启用:1=启用,2=禁用';default:1" json:"enabled"`
	LastCalcTime util.HTime `gorm:"column:last_calc_time;comment:'最后计算时间'" json:"lastCalcTime"`
	CreateTime   util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime   util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (CapacityPrediction) TableName() string {
	return "capacity_prediction"
}

// PredictionHistory 预测历史记录
type PredictionHistory struct {
	ID           uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	PredictionID uint       `gorm:"column:prediction_id;comment:'预测任务ID';NOT NULL;index" json:"predictionId"`
	PredictTime  util.HTime `gorm:"column:predict_time;comment:'预测时间点';NOT NULL" json:"predictTime"`
	PredictValue float64    `gorm:"column:predict_value;type:decimal(10,2);comment:'预测值'" json:"predictValue"`
	ActualValue  float64    `gorm:"column:actual_value;type:decimal(10,2);comment:'实际值'" json:"actualValue"`
	Accuracy     float64    `gorm:"column:accuracy;type:decimal(5,2);comment:'准确度'" json:"accuracy"`
	CreateTime   util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
}

func (PredictionHistory) TableName() string {
	return "prediction_history"
}

// AlertRule 智能告警规则
type AlertRule struct {
	ID             uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Name           string     `gorm:"column:name;type:varchar(200);comment:'规则名称';NOT NULL" json:"name"`
	Type           string     `gorm:"column:type;type:varchar(50);comment:'类型:threshold/trend/anomaly/prediction';NOT NULL;index" json:"type"`
	TargetType     string     `gorm:"column:target_type;type:varchar(20);comment:'目标类型:host/group/cluster';NOT NULL" json:"targetType"`
	TargetIDs      string     `gorm:"column:target_ids;type:text;comment:'目标ID列表(JSON)';NOT NULL" json:"targetIds"`
	MetricType     string     `gorm:"column:metric_type;type:varchar(50);comment:'指标类型';NOT NULL;index" json:"metricType"`
	Condition      string     `gorm:"column:condition;type:text;comment:'告警条件(JSON配置)';NOT NULL" json:"condition"`
	Severity       int        `gorm:"column:severity;comment:'严重程度:1=低,2=中,3=高,4=紧急';NOT NULL" json:"severity"`
	NotifyChannels string     `gorm:"column:notify_channels;type:varchar(500);comment:'通知渠道(JSON数组)'" json:"notifyChannels"`
	NotifyUsers    string     `gorm:"column:notify_users;type:text;comment:'通知用户(JSON数组)'" json:"notifyUsers"`
	SilencePeriod  int        `gorm:"column:silence_period;comment:'静默期(分钟)';default:60" json:"silencePeriod"`
	Enabled        int        `gorm:"column:enabled;comment:'是否启用:1=启用,2=禁用';default:1;index" json:"enabled"`
	CreateTime     util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime     util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (AlertRule) TableName() string {
	return "alert_rule"
}

// AlertEvent 告警事件
type AlertEvent struct {
	ID           uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	RuleID       uint       `gorm:"column:rule_id;comment:'关联告警规则';NOT NULL;index" json:"ruleId"`
	RuleName     string     `gorm:"column:rule_name;type:varchar(200);comment:'规则名称';NOT NULL" json:"ruleName"`
	HostID       uint       `gorm:"column:host_id;comment:'主机ID';NOT NULL;index" json:"hostId"`
	HostName     string     `gorm:"column:host_name;type:varchar(200);comment:'主机名';NOT NULL" json:"hostName"`
	HostIP       string     `gorm:"column:host_ip;type:varchar(50);comment:'主机IP'" json:"hostIp"`
	MetricType   string     `gorm:"column:metric_type;type:varchar(50);comment:'指标类型';NOT NULL;index" json:"metricType"`
	MetricName   string     `gorm:"column:metric_name;type:varchar(100);comment:'指标名称';NOT NULL" json:"metricName"`
	CurrentValue float64    `gorm:"column:current_value;type:decimal(10,2);comment:'当前值'" json:"currentValue"`
	Threshold    float64    `gorm:"column:threshold;type:decimal(10,2);comment:'阈值'" json:"threshold"`
	Severity     int        `gorm:"column:severity;comment:'严重程度:1=低,2=中,3=高,4=紧急';NOT NULL;index" json:"severity"`
	Status       int        `gorm:"column:status;comment:'状态:1=触发,2=恢复,3=已确认,4=已解决';NOT NULL;index" json:"status"`
	TriggerTime  util.HTime `gorm:"column:trigger_time;comment:'触发时间';NOT NULL" json:"triggerTime"`
	RecoverTime  util.HTime `gorm:"column:recover_time;comment:'恢复时间'" json:"recoverTime"`
	Duration     int        `gorm:"column:duration;comment:'持续时间(秒)'" json:"duration"`
	AISummary    string     `gorm:"column:ai_summary;type:text;comment:'AI分析摘要'" json:"aiSummary"`
	AISuggestion string     `gorm:"column:ai_suggestion;type:text;comment:'AI处理建议'" json:"aiSuggestion"`
	NotifyStatus int        `gorm:"column:notify_status;comment:'通知状态:1=待发送,2=已发送,3=发送失败';default:1" json:"notifyStatus"`
	NotifyError  string     `gorm:"column:notify_error;type:text;comment:'通知错误信息'" json:"notifyError"`
	CreateTime   util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime   util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (AlertEvent) TableName() string {
	return "ai_alert_event"
}

// CapacityPredictionDto 容量预测DTO
type CapacityPredictionDto struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name" validate:"required"`
	TargetType  string  `json:"targetType" validate:"required,oneof=host group cluster"`
	TargetID    uint    `json:"targetId" validate:"required"`
	MetricType  string  `json:"metricType" validate:"required,oneof=cpu memory disk network custom"`
	MetricName  string  `json:"metricName" validate:"required"`
	Threshold   float64 `json:"threshold" validate:"required"`
	PredictDays int     `json:"predictDays" validate:"required,min=1,max=365"`
	Algorithm   string  `json:"algorithm" validate:"oneof=linear polynomial arima lstm"`
	Enabled     int     `json:"enabled" validate:"oneof=1 2"`
}

// CapacityPredictionQuery 容量预测查询条件
type CapacityPredictionQuery struct {
	Name       string `form:"name"`
	TargetType string `form:"targetType"`
	TargetID   uint   `form:"targetId"`
	MetricType string `form:"metricType"`
	RiskLevel  int    `form:"riskLevel"`
	Status     int    `form:"status"`
	Enabled    int    `form:"enabled"`
	PageNum    int    `form:"pageNum"`
	PageSize   int    `form:"pageSize"`
}

// AlertRuleDto 告警规则DTO
type AlertRuleDto struct {
	ID             uint     `json:"id"`
	Name           string   `json:"name" validate:"required"`
	Type           string   `json:"type" validate:"required,oneof=threshold trend anomaly prediction"`
	TargetType     string   `json:"targetType" validate:"required,oneof=host group cluster"`
	TargetIDs      []uint   `json:"targetIds" validate:"required"`
	MetricType     string   `json:"metricType" validate:"required"`
	Condition      string   `json:"condition" validate:"required"`
	Severity       int      `json:"severity" validate:"required,min=1,max=4"`
	NotifyChannels []string `json:"notifyChannels"`
	NotifyUsers    []uint   `json:"notifyUsers"`
	SilencePeriod  int      `json:"silencePeriod"`
	Enabled        int      `json:"enabled" validate:"oneof=1 2"`
}

// AlertRuleQuery 告警规则查询条件
type AlertRuleQuery struct {
	Name       string `form:"name"`
	Type       string `form:"type"`
	TargetType string `form:"targetType"`
	MetricType string `form:"metricType"`
	Severity   int    `form:"severity"`
	Enabled    int    `form:"enabled"`
	PageNum    int    `form:"pageNum"`
	PageSize   int    `form:"pageSize"`
}

// AlertEventQuery 告警事件查询条件
type AlertEventQuery struct {
	RuleID     uint   `form:"ruleId"`
	HostID     uint   `form:"hostId"`
	HostName   string `form:"hostName"`
	MetricType string `form:"metricType"`
	Severity   int    `form:"severity"`
	Status     int    `form:"status"`
	StartTime  string `form:"startTime"`
	EndTime    string `form:"endTime"`
	PageNum    int    `form:"pageNum"`
	PageSize   int    `form:"pageSize"`
}

// TrendAnalysis 趋势分析结果
type TrendAnalysis struct {
	Type        string    `json:"type"` // increasing/decreasing/stable/fluctuating
	Slope       float64   `json:"slope"`
	ChangeRate  float64   `json:"changeRate"`
	Periodicity string    `json:"periodicity"` // daily/weekly/monthly/none
	Anomalies   []float64 `json:"anomalies"`
}

// PredictResult 预测结果
type PredictResult struct {
	Predictions   []float64  `json:"predictions"`
	ThresholdTime util.HTime `json:"thresholdTime"`
	DaysLeft      int        `json:"daysLeft"`
	Confidence    float64    `json:"confidence"`
}
