package model

import (
	"time"

	"dodevops-api/common/util"
)

// IncidentAnalysis 故障分析任务
type IncidentAnalysis struct {
	ID          uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	IncidentID  string     `gorm:"column:incident_id;type:varchar(100);comment:'故障ID(唯一标识)';NOT NULL;uniqueIndex" json:"incidentId"`
	Title       string     `gorm:"column:title;type:varchar(500);comment:'故障标题';NOT NULL" json:"title"`
	HostID      uint       `gorm:"column:host_id;comment:'故障主机ID';NOT NULL;index" json:"hostId"`
	HostName    string     `gorm:"column:host_name;type:varchar(200);comment:'主机名';NOT NULL" json:"hostName"`
	HostIP      string     `gorm:"column:host_ip;type:varchar(50);comment:'主机IP'" json:"hostIp"`
	Category    string     `gorm:"column:category;type:varchar(50);comment:'故障类别:cpu/memory/disk/network/service/custom';NOT NULL;index" json:"category"`
	Severity    int        `gorm:"column:severity;comment:'严重程度:1=低,2=中,3=高,4=紧急';NOT NULL;index" json:"severity"`
	Status      int        `gorm:"column:status;comment:'状态:1=分析中,2=已定位,3=修复中,4=已修复,5=已关闭';NOT NULL;index" json:"status"`
	DetectTime  util.HTime `gorm:"column:detect_time;comment:'检测时间';NOT NULL" json:"detectTime"`
	StartTime   util.HTime `gorm:"column:start_time;comment:'分析开始时间'" json:"startTime"`
	EndTime     util.HTime `gorm:"column:end_time;comment:'分析结束时间'" json:"endTime"`
	Duration    int        `gorm:"column:duration;comment:'耗时(秒)'" json:"duration"`
	RootCause   string     `gorm:"column:root_cause;type:text;comment:'根因分析结果'" json:"rootCause"`
	Symptoms    string     `gorm:"column:symptoms;type:text;comment:'故障症状(JSON数组)'" json:"symptoms"`
	Evidence    string     `gorm:"column:evidence;type:text;comment:'证据数据(JSON)'" json:"evidence"`
	Solution    string     `gorm:"column:solution;type:text;comment:'解决方案'" json:"solution"`
	FixScript   string     `gorm:"column:fix_script;type:text;comment:'修复脚本'" json:"fixScript"`
	RiskLevel   int        `gorm:"column:risk_level;comment:'自动修复风险等级:0=安全,1=低危,2=中危,3=高危';default:0" json:"riskLevel"`
	IsAutoFixed bool       `gorm:"column:is_auto_fixed;comment:'是否自动修复';default:false" json:"isAutoFixed"`
	FixResult   string     `gorm:"column:fix_result;type:text;comment:'修复结果'" json:"fixResult"`
	ErrorMsg    string     `gorm:"column:error_msg;type:text;comment:'错误信息'" json:"errorMsg"`
	CreateTime  util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime  util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (IncidentAnalysis) TableName() string {
	return "ai_incident_analysis"
}

// IncidentEvidence 故障证据
type IncidentEvidence struct {
	ID          uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	IncidentID  string     `gorm:"column:incident_id;type:varchar(100);comment:'故障ID';NOT NULL;index" json:"incidentId"`
	Type        string     `gorm:"column:type;type:varchar(50);comment:'类型:log/metric/process/network/file';NOT NULL;index" json:"type"`
	Source      string     `gorm:"column:source;type:varchar(500);comment:'来源:系统日志路径、监控指标名等';NOT NULL" json:"source"`
	Content     string     `gorm:"column:content;type:text;comment:'证据内容';NOT NULL" json:"content"`
	CollectTime util.HTime `gorm:"column:collect_time;comment:'采集时间';NOT NULL" json:"collectTime"`
	Relevance   float64    `gorm:"column:relevance;type:decimal(3,2);comment:'相关度评分(0-1)';default:0.5" json:"relevance"`
	CreateTime  util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
}

func (IncidentEvidence) TableName() string {
	return "ai_incident_evidence"
}

// IncidentKnowledge 故障知识库
type IncidentKnowledge struct {
	ID          uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Title       string     `gorm:"column:title;type:varchar(500);comment:'故障标题';NOT NULL" json:"title"`
	Category    string     `gorm:"column:category;type:varchar(50);comment:'故障类别';NOT NULL;index" json:"category"`
	Symptoms    string     `gorm:"column:symptoms;type:text;comment:'症状描述(JSON数组)';NOT NULL" json:"symptoms"`
	RootCause   string     `gorm:"column:root_cause;type:text;comment:'根因';NOT NULL" json:"rootCause"`
	Solution    string     `gorm:"column:solution;type:text;comment:'解决方案';NOT NULL" json:"solution"`
	FixScript   string     `gorm:"column:fix_script;type:text;comment:'修复脚本'" json:"fixScript"`
	Keywords    string     `gorm:"column:keywords;type:varchar(1000);comment:'关键词(用于匹配)'" json:"keywords"`
	MatchCount  int        `gorm:"column:match_count;comment:'匹配次数(用于排序)';default:0" json:"matchCount"`
	SuccessRate float64    `gorm:"column:success_rate;type:decimal(5,2);comment:'成功率';default:0" json:"successRate"`
	Enabled     int        `gorm:"column:enabled;comment:'是否启用:1=启用,2=禁用';default:1" json:"enabled"`
	CreateTime  util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime  util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (IncidentKnowledge) TableName() string {
	return "incident_knowledge"
}

// IncidentDto 故障分析DTO
type IncidentDto struct {
	ID         uint     `json:"id"`
	IncidentID string   `json:"incidentId"`
	Title      string   `json:"title" validate:"required"`
	HostID     uint     `json:"hostId" validate:"required"`
	Category   string   `json:"category" validate:"required,oneof=cpu memory disk network service custom"`
	Severity   int      `json:"severity" validate:"required,min=1,max=4"`
	Symptoms   []string `json:"symptoms"`
}

// IncidentQuery 故障分析查询条件
type IncidentQuery struct {
	IncidentID string `form:"incidentId"`
	HostID     uint   `form:"hostId"`
	HostName   string `form:"hostName"`
	Category   string `form:"category"`
	Severity   int    `form:"severity"`
	Status     int    `form:"status"`
	StartTime  string `form:"startTime"`
	EndTime    string `form:"endTime"`
	PageNum    int    `form:"pageNum"`
	PageSize   int    `form:"pageSize"`
}

// IncidentAnalysisResult AI分析结果
type IncidentAnalysisResult struct {
	RootCause   string   `json:"rootCause"`
	KeyEvidence []string `json:"keyEvidence"`
	Impact      string   `json:"impact"`
	Solution    string   `json:"solution"`
	FixScript   string   `json:"fixScript"`
	RiskLevel   int      `json:"riskLevel"`
	RiskReason  string   `json:"riskReason"`
}

// IncidentKnowledgeDto 故障知识库DTO
type IncidentKnowledgeDto struct {
	ID        uint     `json:"id"`
	Title     string   `json:"title" validate:"required"`
	Category  string   `json:"category" validate:"required"`
	Symptoms  []string `json:"symptoms" validate:"required"`
	RootCause string   `json:"rootCause" validate:"required"`
	Solution  string   `json:"solution" validate:"required"`
	FixScript string   `json:"fixScript"`
	Keywords  []string `json:"keywords"`
}

// IncidentKnowledgeQuery 故障知识库查询条件
type IncidentKnowledgeQuery struct {
	Title    string `form:"title"`
	Category string `form:"category"`
	Keyword  string `form:"keyword"`
	Enabled  int    `form:"enabled"`
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
}

// ==================== 磁盘分析相关模型 ====================

// DiskAnalysisResult 磁盘分析结果
type DiskAnalysisResult struct {
	HostID           uint            `json:"hostId"`
	HostName         string          `json:"hostName"`
	HostIP           string          `json:"hostIp"`
	CollectTime      time.Time       `json:"collectTime"`
	DiskUsage        string          `json:"diskUsage"`        // df -h 输出
	RootUsagePercent float64         `json:"rootUsagePercent"` // 根分区使用率
	AnalyzedPath     string          `json:"analyzedPath"`     // 深度分析目录
	Directories      []DirectoryInfo `json:"directories"`      // 目录占用信息
	TopFiles         []DirectoryInfo `json:"topFiles"`         // 大文件信息
	CleanupHints     []string        `json:"cleanupHints"`     // 清理建议
}

// DirectoryInfo 目录信息
type DirectoryInfo struct {
	Path     string  `json:"path"`     // 目录路径
	Size     string  `json:"size"`     // 大小（带单位，如 "1.2G"）
	SizeInMB float64 `json:"sizeInMB"` // 大小（MB为单位，用于排序）
}
