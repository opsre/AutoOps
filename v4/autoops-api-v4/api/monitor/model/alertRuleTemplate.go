package model

import "dodevops-api/common/util"

// AlertRuleTemplate 告警规则模板(仅定义指标信息)
type AlertRuleTemplate struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(200);not null;comment:指标名称,如CPU使用率"`
	Description string    `json:"description" gorm:"type:varchar(500);comment:指标描述"`

	// 数据类型/分类
	Category    string    `json:"category" gorm:"type:varchar(50);not null;index;comment:数据类型(cpu/memory/load/disk/io/network/connection/service)"`

	// 指标项
	MetricKey   string    `json:"metricKey" gorm:"type:varchar(100);not null;uniqueIndex;comment:指标项,如system_cpu_usage_percent"`

	// 单位
	Unit        string    `json:"unit" gorm:"type:varchar(20);comment:单位(percentage/string/number/bit/GB/MB/byte)"`

	// 值类型
	ValueType   string    `json:"valueType" gorm:"type:varchar(20);default:number;comment:值类型(number/string/boolean)"`

	// 模板标签
	Tags        string    `json:"tags" gorm:"type:text;comment:标签(JSON数组,如[\"基础监控\",\"资源告警\"])"`

	// 状态
	Status      int       `json:"status" gorm:"type:tinyint;default:1;comment:状态(1:启用,0:禁用)"`
	IsBuiltin   bool      `json:"isBuiltin" gorm:"type:tinyint;default:0;comment:是否内置模板"`

	CreatedBy   string    `json:"createdBy" gorm:"type:varchar(100);comment:创建人"`
	UpdatedBy   string    `json:"updatedBy" gorm:"type:varchar(100);comment:更新人"`
	CreatedAt   util.HTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   util.HTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (AlertRuleTemplate) TableName() string {
	return "monitor_alert_rule_template"
}

// AlertPolicy 告警策略(配置告警规则和通知)
type AlertPolicy struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(200);not null;uniqueIndex:idx_name;comment:策略名称"`
	Description string    `json:"description" gorm:"type:varchar(500);comment:策略描述"`

	// 关联业务组
	GroupID     uint      `json:"groupId" gorm:"type:int;not null;index;comment:业务组ID"`
	GroupName   string    `json:"groupName" gorm:"type:varchar(100);comment:业务组名称"`

	// 通知设置
	RobotIDs    string    `json:"robotIds" gorm:"type:text;comment:通知机器人ID列表(JSON数组)"`

	// 告警抑制
	SilencePeriod  int    `json:"silencePeriod" gorm:"type:int;default:300;comment:静默时间(秒)"`
	RepeatInterval int    `json:"repeatInterval" gorm:"type:int;default:600;comment:重复通知间隔(秒),0表示不重复通知"`
	MaxNotifyCount int    `json:"maxNotifyCount" gorm:"type:int;default:10;comment:最大通知次数"`

	// 恢复通知
	SendRecovery bool     `json:"sendRecovery" gorm:"type:tinyint;default:1;comment:是否发送恢复通知"`

	// 策略状态
	Status      int       `json:"status" gorm:"type:tinyint;default:0;index;comment:状态(1:启用,0:禁用)"`

	// 实时状态
	HealthStatus string   `json:"healthStatus" gorm:"type:varchar(20);default:healthy;comment:实时健康状态(healthy:健康,alerting:告警中)"`

	// 附加标签
	AppendTags  string    `json:"appendTags" gorm:"type:text;comment:附加标签(JSON数组)"`

	CreatedBy   string    `json:"createdBy" gorm:"type:varchar(100);comment:创建人"`
	UpdatedBy   string    `json:"updatedBy" gorm:"type:varchar(100);comment:更新人"`
	CreatedAt   util.HTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   util.HTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (AlertPolicy) TableName() string {
	return "monitor_alert_policy"
}

// AlertPolicyRule 告警策略规则(策略下的具体告警规则配置)
type AlertPolicyRule struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PolicyID    uint      `json:"policyId" gorm:"type:int;not null;index;comment:策略ID"`
	TemplateID  uint      `json:"templateId" gorm:"type:int;not null;index;comment:规则模板ID"`

	// 规则基本信息(从模板继承)
	RuleName    string    `json:"ruleName" gorm:"type:varchar(200);not null;comment:规则名称"`
	Description string    `json:"description" gorm:"type:varchar(500);comment:规则描述"`
	Category    string    `json:"category" gorm:"type:varchar(50);comment:数据类型"`
	MetricKey   string    `json:"metricKey" gorm:"type:varchar(100);comment:指标项"`
	Unit        string    `json:"unit" gorm:"type:varchar(20);comment:单位"`

	// 告警条件配置(在策略中配置)
	Operator    string    `json:"operator" gorm:"type:varchar(10);not null;comment:比较运算符(>/</>=/<=/==/!=)"`
	Threshold   float64   `json:"threshold" gorm:"type:decimal(10,2);not null;comment:阈值"`

	// PromQL(自动生成或手动编写)
	PromQL      string    `json:"promql" gorm:"type:text;not null;comment:Prometheus查询表达式"`

	// 告警级别
	Severity    int       `json:"severity" gorm:"type:tinyint;default:2;comment:告警级别(1:critical,2:warning,3:info)"`

	// 持续时间配置
	Duration    int       `json:"duration" gorm:"type:int;default:60;comment:持续时间(秒),持续多久才告警"`
	ContinuousTimes int   `json:"continuousTimes" gorm:"type:int;default:1;comment:持续次数,连续触发多少次才告警"`

	// 恢复配置
	RecoverDuration int   `json:"recoverDuration" gorm:"type:int;default:60;comment:恢复判定时间(秒)"`

	// 规则状态
	Status      int       `json:"status" gorm:"type:tinyint;default:1;comment:规则是否生效(1:生效,0:不生效)"`

	// 实时状态
	HealthStatus string   `json:"healthStatus" gorm:"type:varchar(20);default:healthy;comment:实时健康状态(healthy:健康,alerting:告警中)"`

	CreatedAt   util.HTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   util.HTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (AlertPolicyRule) TableName() string {
	return "monitor_alert_policy_rule"
}

// DTO - 创建规则模板
type CreateRuleTemplateDto struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Category    string `json:"category" validate:"required"`
	MetricKey   string `json:"metricKey" validate:"required"`
	Unit        string `json:"unit"`
	ValueType   string `json:"valueType"`
	Tags        string `json:"tags"`
	Status      int    `json:"status"`
}

// DTO - 更新规则模板
type UpdateRuleTemplateDto struct {
	ID          uint   `json:"id" validate:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	MetricKey   string `json:"metricKey"`
	Unit        string `json:"unit"`
	ValueType   string `json:"valueType"`
	Tags        string `json:"tags"`
	Status      int    `json:"status"`
}

// DTO - 创建告警策略
type CreateAlertPolicyDto struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description"`
	GroupID         uint   `json:"groupId" validate:"required"`
	GroupName       string `json:"groupName"`
	RobotIDs        string `json:"robotIds"`
	SilencePeriod   int    `json:"silencePeriod"`
	RepeatInterval  int    `json:"repeatInterval"`
	MaxNotifyCount  int    `json:"maxNotifyCount"`
	SendRecovery    bool   `json:"sendRecovery"`
	AppendTags      string `json:"appendTags"`
	Status          int    `json:"status"`
}

// DTO - 更新告警策略
type UpdateAlertPolicyDto struct {
	ID              uint   `json:"id" validate:"required"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	GroupID         uint   `json:"groupId"`
	GroupName       string `json:"groupName"`
	RobotIDs        string `json:"robotIds"`
	SilencePeriod   int    `json:"silencePeriod"`
	RepeatInterval  int    `json:"repeatInterval"`
	MaxNotifyCount  int    `json:"maxNotifyCount"`
	SendRecovery    bool   `json:"sendRecovery"`
	AppendTags      string `json:"appendTags"`
	Status          int    `json:"status"`
}

// DTO - 创建策略规则(从模板导入)
type CreatePolicyRuleDto struct {
	PolicyID        uint    `json:"policyId" validate:"required"`
	TemplateID      uint    `json:"templateId" validate:"required"`
	Description     string  `json:"description"`        // 规则描述
	Operator        string  `json:"operator" validate:"required"`
	Threshold       float64 `json:"threshold" validate:"required"`
	Severity        int     `json:"severity"`
	Duration        int     `json:"duration"`
	ContinuousTimes int     `json:"continuousTimes"`
	RecoverDuration int     `json:"recoverDuration"`
	Status          int     `json:"status"`
}

// DTO - 更新策略规则
type UpdatePolicyRuleDto struct {
	ID              uint    `json:"id" validate:"required"`
	RuleName        string  `json:"ruleName"`           // 规则名称
	Description     string  `json:"description"`        // 规则描述
	Operator        string  `json:"operator"`
	Threshold       float64 `json:"threshold"`
	Severity        int     `json:"severity"`
	Duration        int     `json:"duration"`
	ContinuousTimes int     `json:"continuousTimes"`
	RecoverDuration int     `json:"recoverDuration"`
	Status          int     `json:"status"`
}

// DTO - Prometheus查询预览
type PrometheusQueryPreviewDto struct {
	PromQL    string `json:"promql" validate:"required"`
	Timestamp int64  `json:"timestamp"` // Unix时间戳，0表示当前时间
}

// DTO - 策略规则预览响应
type PolicyRulePreviewResponse struct {
	RuleID      uint                       `json:"ruleId"`
	RuleName    string                     `json:"ruleName"`
	PromQL      string                     `json:"promql"`
	Threshold   float64                    `json:"threshold"`
	Operator    string                     `json:"operator"`
	Unit        string                     `json:"unit"`
	QueryTime   string                     `json:"queryTime"`
	ResultCount int                        `json:"resultCount"`
	Results     []PolicyRulePreviewResult  `json:"results"`
}

// DTO - 策略规则预览结果项
type PolicyRulePreviewResult struct {
	Hostname     string             `json:"hostname"`
	Instance     string             `json:"instance"`
	Value        float64            `json:"value"`
	IsAnomalous  bool               `json:"isAnomalous"`  // 是否异常（超过阈值）
	Labels       map[string]string  `json:"labels"`
}
