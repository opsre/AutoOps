// 运维巡检相关模型
package model

import (
	"dodevops-api/common/util"
)

// AgentInspection 巡检任务
type AgentInspection struct {
	ID           uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Name         string     `gorm:"column:name;type:varchar(200);comment:'巡检任务名称';NOT NULL" json:"name"`
	Type         string     `gorm:"column:type;type:varchar(50);comment:'巡检类型:system/service/security/performance/backup';NOT NULL" json:"type"`
	TargetType   string     `gorm:"column:target_type;type:varchar(20);comment:'目标类型:host/group/tag';NOT NULL" json:"targetType"`
	TargetIDs    string     `gorm:"column:target_ids;type:text;comment:'目标ID列表(JSON数组)';NOT NULL" json:"targetIds"`
	ScheduleType string     `gorm:"column:schedule_type;type:varchar(20);comment:'调度类型:manual/cron';NOT NULL;default:'manual'" json:"scheduleType"`
	CronExpr     string     `gorm:"column:cron_expr;type:varchar(100);comment:'Cron表达式'" json:"cronExpr"`
	Enabled      int        `gorm:"column:enabled;comment:'是否启用:1=启用,2=禁用';NOT NULL;default:1" json:"enabled"`
	Template     string     `gorm:"column:template;type:text;comment:'巡检模板(JSON配置)';NOT NULL" json:"template"`
	Remark       string     `gorm:"column:remark;type:varchar(500);comment:'备注'" json:"remark"`
	CreateTime   util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime   util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (AgentInspection) TableName() string {
	return "ai_agent_inspection"
}

// InspectionRecord 巡检记录
type InspectionRecord struct {
	ID           uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	InspectionID uint       `gorm:"column:inspection_id;comment:'关联巡检任务';NOT NULL;index" json:"inspectionId"`
	HostID       uint       `gorm:"column:host_id;comment:'主机ID';NOT NULL;index" json:"hostId"`
	HostName     string     `gorm:"column:host_name;type:varchar(200);comment:'主机名';NOT NULL" json:"hostName"`
	HostIP       string     `gorm:"column:host_ip;type:varchar(50);comment:'主机IP'" json:"hostIp"`
	Status       int        `gorm:"column:status;comment:'状态:1=执行中,2=成功,3=失败,4=部分成功';NOT NULL;index" json:"status"`
	StartTime    util.HTime `gorm:"column:start_time;comment:'开始时间';NOT NULL" json:"startTime"`
	EndTime      util.HTime `gorm:"column:end_time;comment:'结束时间'" json:"endTime"`
	Duration     int        `gorm:"column:duration;comment:'耗时(秒)'" json:"duration"`
	TotalChecks  int        `gorm:"column:total_checks;comment:'总检查项';default:0" json:"totalChecks"`
	PassedChecks int        `gorm:"column:passed_checks;comment:'通过数';default:0" json:"passedChecks"`
	FailedChecks int        `gorm:"column:failed_checks;comment:'失败数';default:0" json:"failedChecks"`
	WarnChecks   int        `gorm:"column:warn_checks;comment:'警告数';default:0" json:"warnChecks"`
	RiskLevel    int        `gorm:"column:risk_level;comment:'风险等级:0=正常,1=低危,2=中危,3=高危';default:0;index" json:"riskLevel"`
	Result       string     `gorm:"column:result;type:text;comment:'详细结果(JSON)'" json:"result"`
	Summary      string     `gorm:"column:summary;type:text;comment:'AI生成的总结'" json:"summary"`
	Suggestions  string     `gorm:"column:suggestions;type:text;comment:'AI生成的建议'" json:"suggestions"`
	ReportURL    string     `gorm:"column:report_url;type:varchar(500);comment:'报告文件URL'" json:"reportUrl"`
	ErrorMsg     string     `gorm:"column:error_msg;type:text;comment:'错误信息'" json:"errorMsg"`
	CreateTime   util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
}

func (InspectionRecord) TableName() string {
	return "inspection_record"
}

// InspectionItem 巡检项目
type InspectionItem struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	RecordID   uint       `gorm:"column:record_id;comment:'关联巡检记录';NOT NULL;index" json:"recordId"`
	Category   string     `gorm:"column:category;type:varchar(50);comment:'分类:system/service/security/performance';NOT NULL;index" json:"category"`
	Name       string     `gorm:"column:name;type:varchar(200);comment:'检查项名称';NOT NULL" json:"name"`
	Status     int        `gorm:"column:status;comment:'状态:1=通过,2=失败,3=警告';NOT NULL" json:"status"`
	Actual     string     `gorm:"column:actual;type:text;comment:'实际值'" json:"actual"`
	Expected   string     `gorm:"column:expected;type:varchar(500);comment:'期望值'" json:"expected"`
	Message    string     `gorm:"column:message;type:text;comment:'详细信息'" json:"message"`
	Suggestion string     `gorm:"column:suggestion;type:text;comment:'修复建议'" json:"suggestion"`
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
}

func (InspectionItem) TableName() string {
	return "inspection_item"
}

// InspectionDto 巡检任务DTO
type InspectionDto struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name" validate:"required"`
	Type         string   `json:"type" validate:"required,oneof=system service security performance backup"`
	TargetType   string   `json:"targetType" validate:"required,oneof=host group tag"`
	TargetIDs    []uint   `json:"targetIds" validate:"required"`
	ScheduleType string   `json:"scheduleType" validate:"required,oneof=manual cron"`
	CronExpr     string   `json:"cronExpr"`
	Enabled      int      `json:"enabled" validate:"oneof=1 2"`
	Template     string   `json:"template" validate:"required"`
	Remark       string   `json:"remark"`
}

// InspectionQuery 巡检任务查询条件
type InspectionQuery struct {
	Name         string `form:"name"`
	Type         string `form:"type"`
	TargetType   string `form:"targetType"`
	ScheduleType string `form:"scheduleType"`
	Enabled      int    `form:"enabled"`
	PageNum      int    `form:"pageNum"`
	PageSize     int    `form:"pageSize"`
}

// InspectionRecordQuery 巡检记录查询条件
type InspectionRecordQuery struct {
	InspectionID uint   `form:"inspectionId"`
	HostID       uint   `form:"hostId"`
	HostName     string `form:"hostName"`
	Status       int    `form:"status"`
	RiskLevel    int    `form:"riskLevel"`
	StartTime    string `form:"startTime"`
	EndTime      string `form:"endTime"`
	PageNum      int    `form:"pageNum"`
	PageSize     int    `form:"pageSize"`
}

// InspectionTemplate 巡检模板配置
type InspectionTemplate struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Checks      []InspectionCheck  `json:"checks"`
}

// InspectionCheck 巡检检查项配置
type InspectionCheck struct {
	Name       string            `json:"name"`
	Category   string            `json:"category"`
	Type       string            `json:"type"` // command/metric/file/process
	Command    string            `json:"command,omitempty"`
	Metric     string            `json:"metric,omitempty"`
	FilePath   string            `json:"filePath,omitempty"`
	Process    string            `json:"process,omitempty"`
	Expected   string            `json:"expected"`
	Operator   string            `json:"operator"` // eq/ne/gt/lt/gte/lte/contains/regex
	Severity   string            `json:"severity"` // pass/warn/fail
	Message    string            `json:"message"`
	Suggestion string            `json:"suggestion"`
	Params     map[string]string `json:"params,omitempty"`
}

// InspectionResultSummary 巡检结果汇总
type InspectionResultSummary struct {
	RecordID     uint                        `json:"recordId"`
	HostID       uint                        `json:"hostId"`
	HostName     string                      `json:"hostName"`
	HostIP       string                      `json:"hostIp"`
	Status       int                         `json:"status"`
	RiskLevel    int                         `json:"riskLevel"`
	TotalChecks  int                         `json:"totalChecks"`
	PassedChecks int                         `json:"passedChecks"`
	FailedChecks int                         `json:"failedChecks"`
	WarnChecks   int                         `json:"warnChecks"`
	Duration     int                         `json:"duration"`
	Items        []InspectionItem            `json:"items"`
	Summary      string                      `json:"summary"`
	Suggestions  string                      `json:"suggestions"`
	CategoryStats map[string]CategoryStat    `json:"categoryStats"`
}

// CategoryStat 分类统计
type CategoryStat struct {
	Category string `json:"category"`
	Total    int    `json:"total"`
	Passed   int    `json:"passed"`
	Failed   int    `json:"failed"`
	Warned   int    `json:"warned"`
}
