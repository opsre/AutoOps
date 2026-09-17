package model

import (
	"dodevops-api/common/util"
)

// Incident 故障管理模型
type Incident struct {
	ID               uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	AlertTime        util.HTime `gorm:"column:alert_time;comment:'告警时间';NOT NULL" json:"alertTime"`
	BusinessLineID   uint       `gorm:"column:business_line_id;comment:'业务线ID(关联cmdb_group二级分组)'" json:"businessLineId"`
	BusinessLineName string     `gorm:"-" json:"businessLineName"` // 业务线名称（虚拟字段，用于返回）
	Frequency        string     `gorm:"column:frequency;varchar(32);comment:'频率(偶发/频繁)'" json:"frequency"`
	AlertDesc        string     `gorm:"column:alert_desc;type:text;comment:'告警描述'" json:"alertDesc"`
	AlertLevel       string     `gorm:"column:alert_level;varchar(16);comment:'告警级别(P1/P2/P3/P4)';default:'P4'" json:"alertLevel"`
	IncidentCause    string     `gorm:"column:incident_cause;type:text;comment:'故障原因'" json:"incidentCause"`
	Department       string     `gorm:"column:department;varchar(64);comment:'所属部门(研发部/运维部)'" json:"department"`
	Solution         string     `gorm:"column:solution;type:text;comment:'解决方案'" json:"solution"`
	DetailURL        string     `gorm:"column:detail_url;varchar(512);comment:'故障详情URL链接'" json:"detailUrl"`
	Handler          string     `gorm:"column:handler;varchar(64);comment:'处理人(用户名)'" json:"handler"`
	HandlerID        uint       `gorm:"column:handler_id;comment:'处理人ID'" json:"handlerId"`
	Status           int        `gorm:"column:status;comment:'处理状态:1->未处理,2->处理中,3->已完成';default:1" json:"status"`
	Remark           string     `gorm:"column:remark;type:text;comment:'备注信息'" json:"remark"`
	CreateTime       util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime       util.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
}

func (Incident) TableName() string {
	return "monitor_incident"
}

// IncidentStatus 故障处理状态常量
const (
	IncidentStatusPending    = 1 // 未处理
	IncidentStatusProcessing = 2 // 处理中
	IncidentStatusCompleted  = 3 // 已完成
)

// GetStatusText 获取状态文本
func (i *Incident) GetStatusText() string {
	statusMap := map[int]string{
		IncidentStatusPending:    "未处理",
		IncidentStatusProcessing: "处理中",
		IncidentStatusCompleted:  "已完成",
	}
	if text, ok := statusMap[i.Status]; ok {
		return text
	}
	return "未知状态"
}

// CreateIncidentDto 创建故障DTO
type CreateIncidentDto struct {
	AlertTime      util.HTime `json:"alertTime" validate:"required"`     // 告警时间
	BusinessLineID uint       `json:"businessLineId" validate:"required"` // 业务线ID(必须是二级分组)
	Frequency      string     `json:"frequency" validate:"required"`      // 频率
	AlertDesc      string     `json:"alertDesc" validate:"required"`      // 告警描述
	AlertLevel     string     `json:"alertLevel"`                         // 告警级别，默认P4
	IncidentCause  string     `json:"incidentCause"`                      // 故障原因
	Department     string     `json:"department"`                         // 所属部门
	Solution       string     `json:"solution"`                           // 解决方案
	DetailURL      string     `json:"detailUrl"`                          // 故障详情URL
	Status         int        `json:"status"`                             // 处理状态
	Remark         string     `json:"remark"`                             // 备注
	// 注意：Handler 和 HandlerID 字段已移除，将从 token 自动获取
}

// UpdateIncidentDto 更新故障DTO
type UpdateIncidentDto struct {
	ID             uint       `json:"id" validate:"required"`
	AlertTime      util.HTime `json:"alertTime"`
	BusinessLineID uint       `json:"businessLineId"`
	Frequency      string     `json:"frequency"`
	AlertDesc      string     `json:"alertDesc"`
	AlertLevel     string     `json:"alertLevel"`
	IncidentCause  string     `json:"incidentCause"`
	Department     string     `json:"department"`
	Solution       string     `json:"solution"`
	DetailURL      string     `json:"detailUrl"`
	Status         int        `json:"status"`
	Remark         string     `json:"remark"`
	// 注意：Handler 和 HandlerID 字段已移除，更新时不允许修改处理人
}

// IncidentQueryDto 故障查询DTO
type IncidentQueryDto struct {
	BusinessLineID uint   `json:"businessLineId"` // 业务线ID
	AlertLevel     string `json:"alertLevel"`     // 告警级别
	Status         int    `json:"status"`         // 处理状态
	Department     string `json:"department"`     // 所属部门
	Handler        string `json:"handler"`        // 处理人
	StartTime      string `json:"startTime"`      // 开始时间
	EndTime        string `json:"endTime"`        // 结束时间
	PageNum        int    `json:"pageNum"`        // 页码
	PageSize       int    `json:"pageSize"`       // 页大小
}

// IncidentVO 故障视图对象
type IncidentVO struct {
	ID               uint       `json:"id"`
	AlertTime        util.HTime `json:"alertTime"`
	BusinessLineID   uint       `json:"businessLineId"`
	BusinessLineName string     `json:"businessLineName"` // 业务线名称
	Frequency        string     `json:"frequency"`
	AlertDesc        string     `json:"alertDesc"`
	AlertLevel       string     `json:"alertLevel"`
	IncidentCause    string     `json:"incidentCause"`
	Department       string     `json:"department"`
	Solution         string     `json:"solution"`
	DetailURL        string     `json:"detailUrl"`
	Handler          string     `json:"handler"`
	HandlerID        uint       `json:"handlerId"`
	Status           int        `json:"status"`
	StatusText       string     `json:"statusText"`
	Remark           string     `json:"remark"`
	CreateTime       util.HTime `json:"createTime"`
	UpdateTime       util.HTime `json:"updateTime"`
}

// ToVO 转换为VO
func (i *Incident) ToVO() *IncidentVO {
	return &IncidentVO{
		ID:               i.ID,
		AlertTime:        i.AlertTime,
		BusinessLineID:   i.BusinessLineID,
		BusinessLineName: i.BusinessLineName,
		Frequency:        i.Frequency,
		AlertDesc:        i.AlertDesc,
		AlertLevel:       i.AlertLevel,
		IncidentCause:    i.IncidentCause,
		Department:       i.Department,
		Solution:         i.Solution,
		DetailURL:        i.DetailURL,
		Handler:          i.Handler,
		HandlerID:        i.HandlerID,
		Status:           i.Status,
		StatusText:       i.GetStatusText(),
		Remark:           i.Remark,
		CreateTime:       i.CreateTime,
		UpdateTime:       i.UpdateTime,
	}
}
