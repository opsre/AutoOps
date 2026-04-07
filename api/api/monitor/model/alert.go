package model

import "time"

type PrometheusAlertDB struct {
	Id                 int
	Tpltype            string //发送类型如钉钉、企业微信、飞书等
	Tpluse             string //接受目标如Prometheus、WebHook、graylog
	Tplname            string `gorm:"index"`
	Tpl                string `gorm:"type:text"`
	WebhookContentType string // webhook 请求的 contentType 如 application/json, application/x-www-form-urlencoded 等
	Created            time.Time
}

func (PrometheusAlertDB) TableName() string {
	return "monitor_prometheus_alert"
}

type AlertRecord struct {
	Id          int64
	Alertname   string
	AlertLevel  string
	Labels      string
	Instance    string
	StartsAt    string
	EndsAt      string
	Summary     string
	Description string
	AlertStatus string
	CreatedTime time.Time
	UpdatedBy   string
	UpdatedTime time.Time
}

func (AlertRecord) TableName() string {
	return "monitor_alert_record"
}

type AlertRouter struct {
	Id           int `gorm:"primaryKey;autoIncrement"`
	Name         string
	Tpl          *PrometheusAlertDB `gorm:"foreignKey:TplId"`
	TplId        int
	Rules        string
	UrlOrPhone   string
	AtSomeOne    string
	AtSomeOneRR  bool
	SendResolved bool
	Created      time.Time
}

type AlertRouterQuery struct {
	Name    string
	Webhook string
}

func (AlertRouter) TableName() string {
	return "monitor_alert_router"
}

type PrometheusAlertMsg struct {
        Tpl                string
        Type               string
        Ddurl              string
        Wxurl              string
        Fsurl              string
        Phone              string
        WebHookUrl         string
        ToUser             string
        Email              string
        EmailTitle         string
        ToParty            string
        ToTag              string
        GroupId            string
        AtSomeOne          string
        RoundRobin         string
        Split              string
        WebhookContentType string
}

type LabelMap struct {
Name  string
Value string
Regex bool
}



type AlertConfig struct {
        Id        int    `gorm:"primaryKey;autoIncrement"`
        ConfKey   string `gorm:"unique;column:conf_key;type:varchar(255)"`
        ConfValue string `gorm:"column:conf_value;type:text"`
}

func (AlertConfig) TableName() string {
        return "monitor_alert_config"
}
