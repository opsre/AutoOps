package model

import "dodevops-api/common/util"

// NotifyRobot 通知机器人模型
type NotifyRobot struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null;comment:机器人名称"`
	Type      string    `json:"type" gorm:"type:varchar(20);not null;comment:机器人类型(dingtalk/feishu/wechat/email/webhook/teams)"`
	Webhook   string    `json:"webhook" gorm:"type:varchar(500);comment:webhook地址或邮箱接收地址"`
	Secret    string    `json:"secret" gorm:"type:varchar(200);comment:加签密钥或邮箱密码"`
	// 邮箱配置字段
	Server    string    `json:"server" gorm:"type:varchar(100);comment:SMTP服务器地址"`
	Port      *int      `json:"port" gorm:"type:int;comment:SMTP端口"`                    // 使用指针类型,允许为nil
	Username  string    `json:"username" gorm:"type:varchar(100);comment:邮箱用户名"`
	Password  string    `json:"password" gorm:"type:varchar(200);comment:邮箱密码"`
	Nickname  string    `json:"nickname" gorm:"type:varchar(100);comment:发件人昵称"`
	// Webhook配置字段
	Headers   string    `json:"headers" gorm:"type:text;comment:自定义HTTP头(JSON格式)"`
	Method    string    `json:"method" gorm:"type:varchar(10);default:POST;comment:HTTP方法(GET/POST/PUT)"`
	// 告警模板字段（所有类型都支持）
	Template  string    `json:"template" gorm:"type:text;comment:告警消息模板(支持变量替换,不同机器人类型格式不同)"`
	Status    int       `json:"status" gorm:"type:tinyint;default:1;comment:状态(1:启用,0:禁用)"`
	Remark    string    `json:"remark" gorm:"type:varchar(500);comment:备注说明"`
	CreatedAt util.HTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt util.HTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (NotifyRobot) TableName() string {
	return "monitor_notify_robot"
}
