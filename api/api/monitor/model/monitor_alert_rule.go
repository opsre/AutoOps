package model

import "gorm.io/gorm"

type MonitorAlertRule struct {
	gorm.Model
	GroupID     uint   `json:"group_id" gorm:"column:group_id;comment:组ID"`
	Alert       string `json:"alert" gorm:"column:alert;type:varchar(255);comment:告警名称"`
	Expr        string `json:"expr" gorm:"column:expr;type:text;comment:告警表达式"`
	ForDuration string `json:"for_duration" gorm:"column:for_duration;type:varchar(64);comment:持续时间"`
	Labels      string `json:"labels" gorm:"column:labels;type:text;comment:规则私有label (JSON格式)"`
	Constraints string `json:"constraints" gorm:"column:constraints;type:text;comment:规则约束条件 (JSON格式)"`
	Severity    string `json:"severity" gorm:"column:severity;type:varchar(64);comment:严重程度"`
	Summary     string `json:"summary" gorm:"column:summary;type:varchar(255);comment:摘要"`
	Description string `json:"description" gorm:"column:description;type:text;comment:描述"`
	RuleContent string `json:"rule_content" gorm:"column:rule_content;type:text;comment:单个规则的完整yaml"`
	Status      string `json:"status" gorm:"column:status;type:varchar(64);default:inactive;comment:运行状态"`
	Style       string `json:"style" gorm:"column:style;type:varchar(64);comment:规则分类(如CPU,Memory)"`
	Enabled     *int   `json:"enabled" gorm:"column:enabled;type:tinyint(1);default:1;comment:是否启用(1启用,0禁用)"`
}

func (MonitorAlertRule) TableName() string {
	return "monitor_alert_rule"
}

// MonitorAlertRuleQuery 查询参数结构体
type MonitorAlertRuleQuery struct {
	Page        int    `form:"page" json:"page"`
	PageSize    int    `form:"pageSize" json:"pageSize"`
	GroupID     uint   `form:"groupId" json:"groupId"`
	Alert       string `form:"alert" json:"alert"`
	Expr        string `form:"expr" json:"expr"`
	ForDuration string `form:"for_duration" json:"for_duration"`
	Labels      string `form:"labels" json:"labels"`
	Constraints string `form:"constraints" json:"constraints"`
	Style       string `form:"style" json:"style"`
	Severity    string `form:"severity" json:"severity"`
	Summary     string `form:"summary" json:"summary"`
	Description string `form:"description" json:"description"`
	Status      string `form:"status" json:"status"`
	Enabled     *int   `form:"enabled" json:"enabled"`
}

// MonitorAlertRuleCheckReq 表达式健康检查请求参数
type MonitorAlertRuleCheckReq struct {
	DataSourceID uint   `json:"data_source_id" binding:"required"`
	Expr         string `json:"expr" binding:"required"`
}
