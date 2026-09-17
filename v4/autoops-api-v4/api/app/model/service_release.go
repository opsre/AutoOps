package model

import (
	"time"
)

// ServiceRelease 服务上线工单主表
type ServiceRelease struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string    `gorm:"type:varchar(255);not null;comment:发布标题" json:"title"`
	BusinessGroupID uint      `gorm:"not null;comment:业务组ID" json:"business_group_id"`
	ImpactFeature   string    `gorm:"type:text;comment:影响功能描述" json:"impact_feature"`

	// 申请人信息
	ApplicantID     uint      `gorm:"not null;comment:申请人ID" json:"applicant_id"`
	ApplicantName   string    `gorm:"type:varchar(100);not null;comment:申请人姓名" json:"applicant_name"`

	// 审批人信息
	OwnerApproverID    *uint  `gorm:"comment:负责人审批人ID" json:"owner_approver_id"`
	OwnerApproverName  string `gorm:"type:varchar(100);comment:负责人审批人姓名" json:"owner_approver_name"`
	SecurityApproverID *uint  `gorm:"comment:安全审批人ID" json:"security_approver_id"`
	SecurityApproverName string `gorm:"type:varchar(100);comment:安全审批人姓名" json:"security_approver_name"`
	TestApproverID     *uint  `gorm:"comment:测试审批人ID" json:"test_approver_id"`
	TestApproverName   string `gorm:"type:varchar(100);comment:测试审批人姓名" json:"test_approver_name"`

	// 审批状态 (1=未审批，2=通过，3=驳回)
	OwnerApprovalStatus    int    `gorm:"type:int;default:1;comment:负责人审批状态" json:"owner_approval_status"`
	SecurityApprovalStatus int    `gorm:"type:int;default:1;comment:安全审批状态" json:"security_approval_status"`
	TestApprovalStatus     int    `gorm:"type:int;default:1;comment:测试审批状态" json:"test_approval_status"`

	// 审批时间
	OwnerApprovalTime    *time.Time `gorm:"comment:负责人审批时间" json:"owner_approval_time"`
	SecurityApprovalTime *time.Time `gorm:"comment:安全审批时间" json:"security_approval_time"`
	TestApprovalTime     *time.Time `gorm:"comment:测试审批时间" json:"test_approval_time"`

	// 审批意见
	OwnerApprovalRemark    string `gorm:"type:text;comment:负责人审批意见" json:"owner_approval_remark"`
	SecurityApprovalRemark string `gorm:"type:text;comment:安全审批意见" json:"security_approval_remark"`
	TestApprovalRemark     string `gorm:"type:text;comment:测试审批意见" json:"test_approval_remark"`

	// 运维发布状态 (1=未发布，2=发布中，3=发布成功，4=发布失败)
	DeployStatus int `gorm:"type:int;default:1;comment:运维发布状态" json:"deploy_status"`

	// 回归测试状态 (1=未测试，2=测试中，3=测试成功，4=测试失败)
	RegressionTestStatus int `gorm:"type:int;default:1;comment:回归测试状态" json:"regression_test_status"`

	// 整体流程状态 (1=流程进行中，2=发布成功，3=发布失败，4=已驳回，5=已取消)
	Status int `gorm:"type:int;default:1;comment:流程状态" json:"status"`

	// 时间信息
	StartTime   *time.Time `gorm:"comment:发布开始时间" json:"start_time"`
	EndTime     *time.Time `gorm:"comment:发布结束时间" json:"end_time"`
	Duration    int        `gorm:"type:int;default:0;comment:发布耗时(秒)" json:"duration"`

	// 服务数量统计
	ServiceCount int `gorm:"type:int;default:0;comment:关联服务数量" json:"service_count"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// 关联的服务列表
	Services []ServiceReleaseItem `gorm:"foreignKey:ReleaseID;constraint:OnDelete:CASCADE" json:"services,omitempty"`
}

// TableName 指定表名
func (ServiceRelease) TableName() string {
	return "app_service_release"
}

// ServiceReleaseItem 服务上线工单明细表（关联的服务信息）
type ServiceReleaseItem struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ReleaseID  uint   `gorm:"not null;index;comment:上线工单ID" json:"release_id"`
	AppID      uint   `gorm:"not null;comment:应用ID" json:"app_id"`
	AppName    string `gorm:"type:varchar(255);not null;comment:应用名称" json:"app_name"`
	AppCode    string `gorm:"type:varchar(100);not null;comment:应用编码" json:"app_code"`

	// 项目信息
	ProjectName string `gorm:"type:varchar(255);not null;comment:项目服务名称" json:"project_name"`
	RepoURL     string `gorm:"type:varchar(500);not null;comment:项目地址" json:"repo_url"`
	Branch      string `gorm:"type:varchar(100);default:master;comment:发布分支" json:"branch"`
	CommitID    string `gorm:"type:varchar(100);not null;comment:Commit ID" json:"commit_id"`

	// 变更信息
	ImpactFeature   string `gorm:"type:text;comment:影响功能" json:"impact_feature"`
	FunctionModule  string `gorm:"type:text;comment:功能模块" json:"function_module"`
	DBChange        string `gorm:"type:text;comment:数据库变更" json:"db_change"`
	ConfigChange    string `gorm:"type:text;comment:配置变更" json:"config_change"`
	Remark          string `gorm:"type:text;comment:备注信息" json:"remark"`

	// Jenkins发布信息
	JenkinsEnvID    *uint  `gorm:"comment:Jenkins环境配置ID" json:"jenkins_env_id"`
	JenkinsJobURL   string `gorm:"type:varchar(500);comment:Jenkins任务URL" json:"jenkins_job_url"`
	Parameters      string `gorm:"type:text;comment:Jenkins构建参数(JSON格式)" json:"parameters"`
	BuildNumber     int    `gorm:"type:int;default:0;comment:构建编号" json:"build_number"`
	LogURL          string `gorm:"type:varchar(500);comment:构建日志URL" json:"log_url"`

	// 发布状态 (1=未发布，2=发布中，3=发布成功，4=发布失败)
	Status       int        `gorm:"type:int;default:1;comment:发布状态" json:"status"`
	StartTime    *time.Time `gorm:"comment:发布开始时间" json:"start_time"`
	EndTime      *time.Time `gorm:"comment:发布结束时间" json:"end_time"`
	Duration     int        `gorm:"type:int;default:0;comment:发布耗时(秒)" json:"duration"`
	ErrorMessage string     `gorm:"type:text;comment:错误信息" json:"error_message"`

	// 执行顺序
	ExecuteOrder int `gorm:"type:int;default:0;comment:执行顺序" json:"execute_order"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	Application *Application `gorm:"foreignKey:AppID" json:"application,omitempty"`
	JenkinsEnv  *JenkinsEnv  `gorm:"foreignKey:JenkinsEnvID" json:"jenkins_env,omitempty"`
}

// TableName 指定表名
func (ServiceReleaseItem) TableName() string {
	return "app_service_release_item"
}

// CreateServiceReleaseRequest 创建服务上线请求
type CreateServiceReleaseRequest struct {
	Title              string                            `json:"title" binding:"required"`
	BusinessGroupID    uint                              `json:"business_group_id" binding:"required"`
	ImpactFeature      string                            `json:"impact_feature" binding:"required"`
	OwnerApproverID    *uint                             `json:"owner_approver_id" binding:"required"`
	SecurityApproverID *uint                             `json:"security_approver_id" binding:"required"`
	TestApproverID     *uint                             `json:"test_approver_id" binding:"required"`
	Services           []CreateServiceReleaseItemRequest `json:"services" binding:"required,min=1"`
}

// CreateServiceReleaseItemRequest 创建服务上线明细请求
type CreateServiceReleaseItemRequest struct {
	AppID          uint              `json:"app_id" binding:"required"`
	ProjectName    string            `json:"project_name" binding:"required"`
	RepoURL        string            `json:"repo_url" binding:"required"`
	Branch         string            `json:"branch"`
	CommitID       string            `json:"commit_id" binding:"required"`
	ImpactFeature  string            `json:"impact_feature"`
	FunctionModule string            `json:"function_module"`
	DBChange       string            `json:"db_change"`
	ConfigChange   string            `json:"config_change"`
	Remark         string            `json:"remark"`
	Environment    string            `json:"environment" binding:"required"` // 发布环境
	Parameters     map[string]string `json:"parameters"`                     // Jenkins构建参数
}

// ServiceReleaseListRequest 服务上线列表请求
type ServiceReleaseListRequest struct {
	Page            int    `form:"page"`
	PageSize        int    `form:"pageSize"`
	BusinessGroupID *uint  `form:"business_group_id"`
	Status          *int   `form:"status"`
	ApplicantID     *uint  `form:"applicant_id"`
	Keyword         string `form:"keyword"` // 搜索关键词（标题）
}

// ServiceReleaseListResponse 服务上线列表响应
type ServiceReleaseListResponse struct {
	Total int              `json:"total"`
	List  []ServiceRelease `json:"list"`
}

// ApprovalRequest 审批请求
type ApprovalRequest struct {
	ReleaseID    uint   `json:"release_id" binding:"required"`
	ApprovalType string `json:"approval_type" binding:"required"` // owner/security/test
	Action       int    `json:"action" binding:"required"`        // 2=通过，3=驳回
	Remark       string `json:"remark"`
}

// DeployReleaseRequest 发布服务上线请求
type DeployReleaseRequest struct {
	ReleaseID     uint `json:"release_id" binding:"required"`
	ExecutionMode int  `json:"execution_mode"` // 1=并行，2=串行
}

// RegressionTestRequest 回归测试请求
type RegressionTestRequest struct {
	ReleaseID uint   `json:"release_id" binding:"required"`
	Status    int    `json:"status" binding:"required"` // 3=测试成功，4=测试失败
	Remark    string `json:"remark"`
}

// ServiceReleaseStatusText 状态文本映射
var ServiceReleaseStatusText = map[int]string{
	1: "流程进行中",
	2: "发布成功",
	3: "发布失败",
	4: "已驳回",
	5: "已取消",
}

// ApprovalStatusText 审批状态文本映射
var ApprovalStatusText = map[int]string{
	1: "未审批",
	2: "通过",
	3: "驳回",
}

// DeployStatusText 发布状态文本映射
var DeployStatusText = map[int]string{
	1: "未发布",
	2: "发布中",
	3: "发布成功",
	4: "发布失败",
}

// RegressionTestStatusText 回归测试状态文本映射
var RegressionTestStatusText = map[int]string{
	1: "未测试",
	2: "测试中",
	3: "测试成功",
	4: "测试失败",
}
