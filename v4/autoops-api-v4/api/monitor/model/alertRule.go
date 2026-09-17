package model

import "dodevops-api/common/util"

// AlertRule 告警规则模型
type AlertRule struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(200);not null;uniqueIndex:idx_name;comment:规则名称"`
	Description string    `json:"description" gorm:"type:varchar(500);comment:规则描述"`

	// 监控对象
	TargetType  string    `json:"targetType" gorm:"type:varchar(20);not null;index;comment:监控类型(host/service/custom)"`
	TargetIDs   string    `json:"targetIds" gorm:"type:text;comment:监控目标ID列表(JSON数组,如[1,2,3])"`

	// 告警条件 - PromQL方式
	PromQL      string    `json:"promql" gorm:"type:text;not null;comment:Prometheus查询表达式"`
	Duration    int       `json:"duration" gorm:"type:int;default:60;comment:持续时间(秒),触发后持续此时间才告警"`

	// 阈值条件(可选,用于简单规则)
	MetricName  string    `json:"metricName" gorm:"type:varchar(100);comment:指标名称(cpu/memory/disk等)"`
	Operator    string    `json:"operator" gorm:"type:varchar(10);comment:比较运算符(>/</>=/<=/==/!=)"`
	Threshold   float64   `json:"threshold" gorm:"type:decimal(10,2);comment:阈值"`

	// 告警级别
	Severity    int       `json:"severity" gorm:"type:tinyint;default:2;comment:告警级别(1:critical,2:warning,3:info)"`

	// 通知设置
	RobotIDs    string    `json:"robotIds" gorm:"type:text;comment:通知机器人ID列表(JSON数组)"`

	// 告警抑制
	SilencePeriod int     `json:"silencePeriod" gorm:"type:int;default:300;comment:静默期(秒),同一告警在此期间不重复发送"`
	RepeatInterval int    `json:"repeatInterval" gorm:"type:int;default:3600;comment:重复通知间隔(秒),0表示不重复通知"`
	MaxNotifyCount int    `json:"maxNotifyCount" gorm:"type:int;default:3;comment:最大通知次数,0表示不限制"`

	// 恢复通知
	SendRecovery bool     `json:"sendRecovery" gorm:"type:tinyint;default:1;comment:是否发送恢复通知"`
	RecoverDuration int   `json:"recoverDuration" gorm:"type:int;default:60;comment:恢复判定时间(秒),持续正常此时间后发送恢复通知"`

	// 调度配置
	CronPattern string    `json:"cronPattern" gorm:"type:varchar(100);default:@every 30s;comment:调度表达式(cron格式)"`

	// 标签与注解
	AppendTags  string    `json:"appendTags" gorm:"type:text;comment:附加标签(JSON数组,如[\"env=prod\",\"team=ops\"])"`
	Annotations string    `json:"annotations" gorm:"type:text;comment:注解信息(JSON对象,如{\"summary\":\"CPU告警\"})"`

	// 规则状态
	Status      int       `json:"status" gorm:"type:tinyint;default:0;index;comment:状态(1:启用,0:禁用)"`

	// 业务信息
	GroupID     uint      `json:"groupId" gorm:"type:int;default:0;index;comment:业务组ID"`
	GroupName   string    `json:"groupName" gorm:"type:varchar(100);comment:业务组名称"`

	CreatedBy   string    `json:"createdBy" gorm:"type:varchar(100);comment:创建人"`
	UpdatedBy   string    `json:"updatedBy" gorm:"type:varchar(100);comment:更新人"`
	CreatedAt   util.HTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   util.HTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (AlertRule) TableName() string {
	return "monitor_alert_rule"
}

// AlertRecord 告警记录模型
type AlertRecord struct {
	ID          uint      `json:"id" gorm:"primaryKey"`

	// 策略和规则信息
	PolicyID    uint      `json:"policyId" gorm:"type:int;index;comment:告警策略ID"`
	PolicyName  string    `json:"policyName" gorm:"type:varchar(200);comment:策略名称"`
	RuleID      uint      `json:"ruleId" gorm:"type:int;not null;index:idx_rule_target;comment:告警规则ID(模板实例ID)"`
	RuleName    string    `json:"ruleName" gorm:"type:varchar(200);comment:规则名称"`
	TemplateID  uint      `json:"templateId" gorm:"type:int;index;comment:规则模板ID"`

	// 告警对象
	TargetType  string    `json:"targetType" gorm:"type:varchar(20);index;comment:监控类型"`
	TargetID    uint      `json:"targetId" gorm:"type:int;index:idx_rule_target;comment:目标ID"`
	TargetName  string    `json:"targetName" gorm:"type:varchar(200);comment:目标名称"`

	// 告警信息
	MetricName  string    `json:"metricName" gorm:"type:varchar(100);index;comment:指标名称"`
	CurrentValue float64  `json:"currentValue" gorm:"type:decimal(10,2);comment:当前值"`
	Threshold   float64   `json:"threshold" gorm:"type:decimal(10,2);comment:阈值"`
	Severity    int       `json:"severity" gorm:"type:tinyint;index;comment:告警级别(1:critical,2:warning,3:info)"`
	Message     string    `json:"message" gorm:"type:text;comment:告警消息"`
	Summary     string    `json:"summary" gorm:"type:varchar(500);comment:告警摘要"`

	// 标签
	Labels      string    `json:"labels" gorm:"type:text;comment:标签信息(JSON对象)"`

	// 告警状态
	Status      string    `json:"status" gorm:"type:varchar(20);default:firing;index;comment:告警状态(firing:触发中,resolved:已恢复,acknowledged:已认领)"`
	StartTime   util.HTime `json:"startTime" gorm:"index;comment:开始时间"`
	EndTime     *util.HTime `json:"endTime" gorm:"comment:结束时间"`
	Duration    int       `json:"duration" gorm:"type:int;comment:持续时间(秒)"`

	// 处理信息
	Handler     string    `json:"handler" gorm:"type:varchar(100);comment:处理人"`
	HandleTime  *util.HTime `json:"handleTime" gorm:"comment:处理时间"`
	HandleNote  string    `json:"handleNote" gorm:"type:text;comment:处理备注"`

	// 通知状态
	NotifySent  bool      `json:"notifySent" gorm:"type:tinyint;default:0;comment:是否已发送通知"`
	NotifyTime  *util.HTime `json:"notifyTime" gorm:"comment:通知时间"`
	NotifyCount int       `json:"notifyCount" gorm:"type:int;default:0;comment:通知次数"`
	NotifyError string    `json:"notifyError" gorm:"type:text;comment:通知错误信息"`

	// 业务信息
	GroupID     uint      `json:"groupId" gorm:"type:int;default:0;index;comment:业务组ID"`
	GroupName   string    `json:"groupName" gorm:"type:varchar(100);comment:业务组名称"`

	// 事件唯一标识(用于去重)
	Fingerprint string    `json:"fingerprint" gorm:"type:varchar(64);uniqueIndex;comment:事件指纹(MD5)"`

	CreatedAt   util.HTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   util.HTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (AlertRecord) TableName() string {
	return "monitor_alert_record"
}

// AlertMute 告警静默配置
type AlertMute struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(200);not null;comment:静默名称"`

	// 静默范围
	RuleIDs     string    `json:"ruleIds" gorm:"type:text;comment:规则ID列表(JSON数组),为空表示全局"`
	TargetType  string    `json:"targetType" gorm:"type:varchar(20);comment:目标类型"`
	TargetIDs   string    `json:"targetIds" gorm:"type:text;comment:目标ID列表(JSON数组)"`

	// 标签匹配
	LabelMatchers string  `json:"labelMatchers" gorm:"type:text;comment:标签匹配器(JSON数组)"`

	// 时间范围
	StartTime   util.HTime `json:"startTime" gorm:"comment:开始时间"`
	EndTime     util.HTime `json:"endTime" gorm:"index;comment:结束时间"`

	// 周期性静默(可选)
	Periodic    bool      `json:"periodic" gorm:"type:tinyint;default:0;comment:是否周期性静默"`
	PeriodDays  string    `json:"periodDays" gorm:"type:varchar(50);comment:周期天数(如1,2,3表示周一到周三)"`
	PeriodStart string    `json:"periodStart" gorm:"type:varchar(10);comment:每天开始时间(如09:00)"`
	PeriodEnd   string    `json:"periodEnd" gorm:"type:varchar(10);comment:每天结束时间(如18:00)"`

	Status      int       `json:"status" gorm:"type:tinyint;default:1;comment:状态(1:启用,0:禁用)"`
	Remark      string    `json:"remark" gorm:"type:varchar(500);comment:备注"`

	CreatedBy   string    `json:"createdBy" gorm:"type:varchar(100);comment:创建人"`
	CreatedAt   util.HTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   util.HTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (AlertMute) TableName() string {
	return "monitor_alert_mute"
}
