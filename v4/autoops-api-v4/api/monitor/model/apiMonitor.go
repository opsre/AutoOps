package model

import "dodevops-api/common/util"

// APIMonitor API监控配置
type APIMonitor struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"type:varchar(200);not null;comment:监控名称"`
	Description string     `json:"description" gorm:"type:varchar(500);comment:描述"`

	// 业务分类（三层结构）
	GroupID            uint       `json:"groupId" gorm:"type:int;index;comment:业务组ID"`
	GroupName          string     `json:"groupName" gorm:"type:varchar(100);comment:业务组名称"`
	BusinessFunctionID uint       `json:"businessFunctionId" gorm:"type:int;index;comment:业务功能ID"`
	BusinessFunctionName string   `json:"businessFunctionName" gorm:"type:varchar(200);comment:业务功能名称"`
	ServiceName        string     `json:"serviceName" gorm:"type:varchar(100);index;comment:服务名称(可选,用于进一步分类)"`

	// 请求配置
	URL         string     `json:"url" gorm:"type:varchar(1000);not null;comment:接口地址"`
	Method      string     `json:"method" gorm:"type:varchar(10);default:GET;comment:HTTP方法"`
	Headers     string     `json:"headers" gorm:"type:text;comment:请求头(JSON)"`
	Body        string     `json:"body" gorm:"type:text;comment:请求体"`
	Timeout     int        `json:"timeout" gorm:"type:int;default:5000;comment:超时时间(毫秒)"`

	// 校验规则
	ExpectedStatus   int    `json:"expectedStatus" gorm:"type:int;default:200;comment:期望状态码"`
	ExpectedKeyword  string `json:"expectedKeyword" gorm:"type:varchar(500);comment:期望响应关键字"`
	ValidateJSON     string `json:"validateJson" gorm:"type:text;comment:JSON校验规则"`

	// 监控配置
	CheckInterval int        `json:"checkInterval" gorm:"type:int;default:60;comment:检测间隔(秒)"`
	RetryCount    int        `json:"retryCount" gorm:"type:int;default:3;comment:失败重试次数"`

	// 告警配置
	AlertEnabled  bool       `json:"alertEnabled" gorm:"type:tinyint;default:1;comment:是否启用告警"`
	PolicyID      *uint      `json:"policyId" gorm:"type:int;comment:关联告警策略ID"`

	// 状态
	Status        int        `json:"status" gorm:"type:tinyint;default:1;index;comment:启用状态(1:启用,0:禁用)"`
	HealthStatus  string     `json:"healthStatus" gorm:"type:varchar(20);default:unknown;comment:健康状态(healthy/unhealthy/unknown)"`

	// 最近检测结果
	LastCheckTime      *util.HTime `json:"lastCheckTime" gorm:"comment:最近检测时间"`
	LastResponseTime   int         `json:"lastResponseTime" gorm:"type:int;comment:最近响应时间(毫秒)"`
	LastStatusCode     int         `json:"lastStatusCode" gorm:"type:int;comment:最近状态码"`
	LastErrorMessage   string      `json:"lastErrorMessage" gorm:"type:varchar(1000);comment:最近错误信息"`
	ContinuousFailures int         `json:"continuousFailures" gorm:"type:int;default:0;comment:连续失败次数"`
	HealthRate         *float64    `json:"healthRate" gorm:"-"` // 健康率(24小时内),不存储到数据库

	// 标签
	Tags        string     `json:"tags" gorm:"type:text;comment:标签(JSON数组)"`

	// 分层关系字段
	ParentID       *uint      `json:"parentId" gorm:"type:int;index;comment:父监控ID,用于实现分层关系"`
	HierarchyLevel int        `json:"hierarchyLevel" gorm:"type:int;default:0;comment:层级深度: 0=独立监控, 1=一级子监控, 2=二级子监控"`
	IsAggregate    bool       `json:"isAggregate" gorm:"type:tinyint;default:0;comment:是否为聚合监控(基于子监控结果计算)"`

	// 关系字段（用于查询时加载）
	ParentMonitor  *APIMonitor  `json:"parentMonitor,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	ChildMonitors  []*APIMonitor `json:"childMonitors,omitempty" gorm:"foreignKey:ParentID;references:ID"`

	// 计算字段（来自service，不存储到数据库）
	ChildCount     int        `json:"childCount,omitempty" gorm:"-"`        // 子监控数量
	HealthyChildren int       `json:"healthyChildren,omitempty" gorm:"-"`   // 健康的子监控数
	ParentName     string     `json:"parentName,omitempty" gorm:"-"`        // 父监控名称（冗余字段）

	CreatedBy   string     `json:"createdBy" gorm:"type:varchar(100);comment:创建人"`
	UpdatedBy   string     `json:"updatedBy" gorm:"type:varchar(100);comment:更新人"`
	CreatedAt   util.HTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   util.HTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (APIMonitor) TableName() string {
	return "monitor_api_monitor"
}

// APIMonitorResult API监控检测结果
type APIMonitorResult struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	MonitorID     uint       `json:"monitorId" gorm:"type:int;not null;index;comment:监控配置ID"`

	// 检测结果
	CheckTime     util.HTime `json:"checkTime" gorm:"not null;index;comment:检测时间"`
	Success       bool       `json:"success" gorm:"type:tinyint;not null;comment:是否成功"`
	StatusCode    int        `json:"statusCode" gorm:"type:int;comment:HTTP状态码"`
	ResponseTime  int        `json:"responseTime" gorm:"type:int;comment:响应时间(毫秒)"`
	ErrorMessage  string     `json:"errorMessage" gorm:"type:varchar(1000);comment:错误信息"`

	// 响应内容(可选存储，用于排查问题)
	ResponseBody  string     `json:"responseBody" gorm:"type:text;comment:响应内容"`
}

func (APIMonitorResult) TableName() string {
	return "monitor_api_monitor_result"
}

// DTO - 创建API监控
type CreateAPIMonitorDto struct {
	Name                 string `json:"name" binding:"required"`
	Description          string `json:"description"`
	GroupID              uint   `json:"groupId" binding:"required"`
	GroupName            string `json:"groupName"`
	BusinessFunctionID   uint   `json:"businessFunctionId" binding:"required"`
	BusinessFunctionName string `json:"businessFunctionName"`
	ServiceName          string `json:"serviceName"`
	URL                  string `json:"url" binding:"required"`
	Method               string `json:"method"`
	Headers              string `json:"headers"`
	Body                 string `json:"body"`
	Timeout              int    `json:"timeout"`
	ExpectedStatus       int    `json:"expectedStatus"`
	ExpectedKeyword      string `json:"expectedKeyword"`
	ValidateJSON         string `json:"validateJson"`
	CheckInterval        int    `json:"checkInterval"`
	RetryCount           int    `json:"retryCount"`
	AlertEnabled         bool   `json:"alertEnabled"`
	PolicyID             *uint  `json:"policyId"`
	Status               int    `json:"status"`
	Tags                 string `json:"tags"`
	ParentID             *uint  `json:"parentId"`             // 父监控ID（可选，用于创建子监控）
	HierarchyLevel       int    `json:"hierarchyLevel"`       // 层级深度（自动计算）
}

// DTO - 简化版创建API监控
type SimpleCreateAPIMonitorDto struct {
	Name                  string `json:"name" binding:"required"`                          // 监控名称
	URL                   string `json:"url" binding:"required"`                           // 接口地址
	Method                string `json:"method"`                                           // HTTP方法（默认GET）
	GroupID               uint   `json:"groupId" binding:"required"`                       // 业务组ID
	BusinessFunctionName  string `json:"businessFunctionName" binding:"required"`          // 业务功能名称
	Description           string `json:"description"`                                      // 描述（可选）
	ParentID              *uint  `json:"parentId"`                                         // 父监控ID（可选，用于创建子监控）
}

// DTO - 更新API监控
type UpdateAPIMonitorDto struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	GroupID              uint   `json:"groupId"`
	GroupName            string `json:"groupName"`
	BusinessFunctionID   uint   `json:"businessFunctionId"`
	BusinessFunctionName string `json:"businessFunctionName"`
	ServiceName          string `json:"serviceName"`
	URL                  string `json:"url"`
	Method               string `json:"method"`
	Headers              string `json:"headers"`
	Body                 string `json:"body"`
	Timeout              int    `json:"timeout"`
	ExpectedStatus       int    `json:"expectedStatus"`
	ExpectedKeyword      string `json:"expectedKeyword"`
	ValidateJSON         string `json:"validateJson"`
	CheckInterval        int    `json:"checkInterval"`
	RetryCount           int    `json:"retryCount"`
	AlertEnabled         bool   `json:"alertEnabled"`
	PolicyID             *uint  `json:"policyId"`
	Status               int    `json:"status"`
	Tags                 string `json:"tags"`
}

// DTO - 测试API接口
type TestAPIDto struct {
	URL            string `json:"url" binding:"required"`
	Method         string `json:"method"`
	Headers        string `json:"headers"`
	Body           string `json:"body"`
	Timeout        int    `json:"timeout"`
	ExpectedStatus int    `json:"expectedStatus"`
}

// DTO - API监控统计
type APIMonitorStats struct {
	TotalCount      int     `json:"totalCount"`      // 总数
	EnabledCount    int     `json:"enabledCount"`    // 启用数
	HealthyCount    int     `json:"healthyCount"`    // 健康数
	UnhealthyCount  int     `json:"unhealthyCount"`  // 异常数
	AvgResponseTime float64 `json:"avgResponseTime"` // 平均响应时间
}

// DTO - 按业务组聚合的API监控统计
type APIMonitorGroupStats struct {
	GroupID             uint                              `json:"id"`                  // 业务组ID (改为id以匹配期望输出)
	GroupName           string                            `json:"name"`                // 业务组名称 (改为name以匹配期望输出)
	ParentID            uint                              `json:"-"`                   // 父级分组ID (用于构建树形结构，不输出)
	TotalAPIs           int                               `json:"apiCount"`            // API总数 (改为apiCount以匹配期望输出)
	HealthyAPIs         int                               `json:"healthy"`             // 健康API数 (改为healthy以匹配期望输出)
	UnhealthyAPIs       int                               `json:"unhealthy"`           // 异常API数 (改为unhealthy以匹配期望输出)
	AvgResponseTime     float64                           `json:"avgResponseTime"`     // 平均响应时间
	BusinessFunctions   []APIMonitorBusinessFunctionStats `json:"businessFunctions,omitempty"` // 业务功能列表
	Children            []*APIMonitorGroupStats           `json:"children,omitempty"`  // 子业务组列表
}

// DTO - 按业务功能聚合的API监控统计
type APIMonitorBusinessFunctionStats struct {
	BusinessFunctionID   uint         `json:"businessFunctionId"`   // 业务功能ID
	BusinessFunctionName string       `json:"businessFunctionName"` // 业务功能名称
	TotalAPIs            int          `json:"totalApis"`            // API总数
	HealthyAPIs          int          `json:"healthyApis"`          // 健康API数
	UnhealthyAPIs        int          `json:"unhealthyApis"`        // 异常API数
	AvgResponseTime      float64      `json:"avgResponseTime"`      // 平均响应时间
	APIs                 []APIMonitor `json:"apis"`                 // API列表
}

// BuildAPIMonitorGroupTree 构建业务组树形结构
func BuildAPIMonitorGroupTree(groups []*APIMonitorGroupStats) []*APIMonitorGroupStats {
	// 创建分组映射
	groupMap := make(map[uint]*APIMonitorGroupStats)
	for _, group := range groups {
		groupMap[group.GroupID] = group
	}

	// 构建树形结构
	var tree []*APIMonitorGroupStats
	for _, group := range groups {
		if group.ParentID == 0 {
			// 一级业务组（父组）
			tree = append(tree, buildAPIMonitorSubTree(group, groupMap))
		}
	}

	return tree
}

// buildAPIMonitorSubTree 递归构建子树
func buildAPIMonitorSubTree(group *APIMonitorGroupStats, groupMap map[uint]*APIMonitorGroupStats) *APIMonitorGroupStats {
	// 查找当前组的所有子组
	for _, child := range groupMap {
		if child.ParentID == group.GroupID {
			// 递归构建子组的子树
			childTree := buildAPIMonitorSubTree(child, groupMap)
			group.Children = append(group.Children, childTree)

			// 累加子组的统计数据到父组
			group.TotalAPIs += child.TotalAPIs
			group.HealthyAPIs += child.HealthyAPIs
			group.UnhealthyAPIs += child.UnhealthyAPIs

			// 响应时间需要重新计算平均值
			// 这里简化处理，如果需要精确的平均值，应该传递总数和总和
		}
	}

	return group
}


// DTO - API监控列表查询参数
type APIMonitorListQuery struct {
	GroupID            *uint  `form:"groupId"`            // 业务组ID
	BusinessFunctionID *uint  `form:"businessFunctionId"` // 业务功能ID
	ServiceName        string `form:"serviceName"`        // 服务名称
	Status             *int   `form:"status"`             // 状态
	Keyword            string `form:"keyword"`            // 关键字（名称/URL）
	Page               int    `form:"page"`               // 页码
	PageSize           int    `form:"pageSize"`           // 每页数量

	// 分层查询参数
	ParentID        *uint `form:"parentId"`        // 按父监控过滤
	OnlyTopLevel    bool  `form:"onlyTopLevel"`    // 只查询顶级监控（parentId=null）
	IncludeChildren bool  `form:"includeChildren"` // 是否加载子监控
}

// DTO - 添加子监控
type AddChildMonitorDto struct {
	ChildMonitorID uint `json:"childMonitorId" binding:"required"` // 子监控ID
}

// DTO - 监控树形结构（用于展示层级关系）
type HierarchyDto struct {
	ID              uint          `json:"id"`
	Name            string        `json:"name"`
	URL             string        `json:"url"`
	HealthStatus    string        `json:"healthStatus"`
	HierarchyLevel  int           `json:"hierarchyLevel"`
	ChildCount      int           `json:"childCount"`
	HealthyChildren int           `json:"healthyChildren"`
	Children        []HierarchyDto `json:"children,omitempty"`
}

// DTO - 子监控统计信息
type ChildrenStats struct {
	ParentID        uint   `json:"parentId"`
	ChildCount      int    `json:"childCount"`
	HealthyChildren int    `json:"healthyChildren"`
	HealthStatus    string `json:"healthStatus"` // 聚合后的父监控健康状态
}
