package model

import (
	"time"
)

// ShRelease 脚本上线工单主表
type ShRelease struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string    `gorm:"type:varchar(255);not null;comment:上线标题" json:"title"`
	Reason          string    `gorm:"type:text;not null;comment:上线原因" json:"reason"`
	BusinessGroupID uint      `gorm:"not null;comment:业务线ID" json:"business_group_id"`

	// 服务信息
	AppID           uint      `gorm:"not null;comment:服务ID" json:"app_id"`
	AppName         string    `gorm:"type:varchar(255);not null;comment:服务名称" json:"app_name"`
	AppCode         string    `gorm:"type:varchar(100);not null;comment:服务编码" json:"app_code"`

	// 申请人信息
	ApplicantID     uint      `gorm:"not null;comment:申请人ID" json:"applicant_id"`
	ApplicantName   string    `gorm:"type:varchar(100);not null;comment:申请人姓名" json:"applicant_name"`

	// 审核人和执行人信息
	ApproverID      *uint     `gorm:"comment:审核人ID" json:"approver_id"`
	ApproverName    string    `gorm:"type:varchar(100);comment:审核人姓名" json:"approver_name"`
	ExecutorID      *uint     `gorm:"comment:执行人ID" json:"executor_id"`
	ExecutorName    string    `gorm:"type:varchar(100);comment:执行人姓名" json:"executor_name"`

	// 脚本执行信息
	ServerHostID    uint      `gorm:"not null;comment:执行服务器主机ID(关联cmdb_host)" json:"server_host_id"` // 执行服务器
	ExecuteDir      string    `gorm:"type:varchar(500);not null;comment:执行目录" json:"execute_dir"`
	ScriptContent   string    `gorm:"type:text;not null;comment:脚本内容" json:"script_content"`
	Parameters      string    `gorm:"type:text;comment:Jenkins构建参数(JSON格式)" json:"parameters"` // 额外的Jenkins参数

	// 审核状态 (1=未审核，2=审核通过，3=审核驳回)
	ApprovalStatus  int       `gorm:"type:int;default:1;comment:审核状态" json:"approval_status"`
	ApprovalTime    *time.Time `gorm:"comment:审核时间" json:"approval_time"`
	ApprovalRemark  string    `gorm:"type:text;comment:审核意见" json:"approval_remark"`

	// 执行状态 (1=未执行，2=拉取代码中，3=代码拉取成功，4=代码拉取失败，5=脚本执行中，6=脚本执行成功，7=脚本执行失败)
	ExecuteStatus   int       `gorm:"type:int;default:1;comment:执行状态" json:"execute_status"`

	// 整体流程状态 (1=流程进行中，2=执行成功，3=执行失败，4=已驳回，5=已取消)
	Status          int       `gorm:"type:int;default:1;comment:流程状态" json:"status"`

	// 时间信息
	PullCodeStartTime   *time.Time `gorm:"comment:拉取代码开始时间" json:"pull_code_start_time"`
	PullCodeEndTime     *time.Time `gorm:"comment:拉取代码结束时间" json:"pull_code_end_time"`
	StartTime           *time.Time `gorm:"comment:脚本执行开始时间" json:"start_time"`
	EndTime             *time.Time `gorm:"comment:脚本执行结束时间" json:"end_time"`
	Duration            int        `gorm:"type:int;default:0;comment:执行耗时(秒)" json:"duration"`

	// Jenkins执行信息（用于拉取代码）
	JenkinsEnvID    *uint     `gorm:"comment:Jenkins环境配置ID" json:"jenkins_env_id"`
	BuildNumber     int       `gorm:"type:int;default:0;comment:构建编号" json:"build_number"`
	LogURL          string    `gorm:"type:varchar(500);comment:构建日志URL" json:"log_url"`

	// 脚本执行信息
	ScriptOutput    string    `gorm:"type:longtext;comment:脚本执行输出" json:"script_output"`
	ErrorMessage    string    `gorm:"type:text;comment:错误信息" json:"error_message"`

	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	Application     *Application `gorm:"foreignKey:AppID" json:"application,omitempty"`
	JenkinsEnv      *JenkinsEnv  `gorm:"foreignKey:JenkinsEnvID" json:"jenkins_env,omitempty"`
}

// TableName 指定表名
func (ShRelease) TableName() string {
	return "app_sh_release"
}

// CreateShReleaseRequest 创建脚本上线请求
type CreateShReleaseRequest struct {
	Title           string            `json:"title" binding:"required"`             // 上线标题
	Reason          string            `json:"reason" binding:"required"`            // 上线原因
	BusinessGroupID uint              `json:"business_group_id" binding:"required"` // 业务线ID
	AppID           uint              `json:"app_id" binding:"required"`            // 服务ID
	Environment     string            `json:"environment" binding:"required"`       // 环境（prod/test/dev）
	ApproverID      *uint             `json:"approver_id" binding:"required"`       // 审核人ID
	ExecutorID      *uint             `json:"executor_id" binding:"required"`       // 执行人ID
	ServerHostID    uint              `json:"server_host_id" binding:"required"`    // 执行服务器主机ID
	ExecuteDir      string            `json:"execute_dir" binding:"required"`       // 执行目录
	ScriptContent   string            `json:"script_content" binding:"required"`    // 脚本内容
	Parameters      map[string]string `json:"parameters"`                           // Jenkins构建参数（可选，如commit_id等）
}

// ShReleaseListRequest 脚本上线列表请求
type ShReleaseListRequest struct {
	Page            int    `form:"page"`
	PageSize        int    `form:"pageSize"`
	BusinessGroupID *uint  `form:"business_group_id"`
	Status          *int   `form:"status"`
	ApplicantID     *uint  `form:"applicant_id"`
	Keyword         string `form:"keyword"` // 搜索关键词（标题）
}

// ShReleaseListResponse 脚本上线列表响应
type ShReleaseListResponse struct {
	Total int         `json:"total"`
	List  []ShRelease `json:"list"`
}

// ShApprovalRequest 脚本审核请求
type ShApprovalRequest struct {
	ReleaseID uint   `json:"release_id" binding:"required"`
	Action    int    `json:"action" binding:"required"` // 2=通过，3=驳回
	Remark    string `json:"remark"`
}

// ShExecuteRequest 脚本执行请求（拉取代码）
type ShExecuteRequest struct {
	ReleaseID uint `json:"release_id" binding:"required"`
}

// ShExecuteScriptRequest 执行脚本请求（在服务器上执行）
type ShExecuteScriptRequest struct {
	ReleaseID uint `json:"release_id" binding:"required"`
}

// ShReleaseStatusText 状态文本映射
var ShReleaseStatusText = map[int]string{
	1: "流程进行中",
	2: "执行成功",
	3: "执行失败",
	4: "已驳回",
	5: "已取消",
}

// ShApprovalStatusText 审核状态文本映射
var ShApprovalStatusText = map[int]string{
	1: "未审核",
	2: "审核通过",
	3: "审核驳回",
}

// ShExecuteStatusText 执行状态文本映射
var ShExecuteStatusText = map[int]string{
	1: "未执行",
	2: "拉取代码中",
	3: "代码拉取成功",
	4: "代码拉取失败",
	5: "脚本执行中",
	6: "脚本执行成功",
	7: "脚本执行失败",
}
