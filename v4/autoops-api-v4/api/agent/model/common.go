// Agent任务、知识库、Prompt模板等通用模型
package model

import (
	"dodevops-api/common/util"
	"time"
)

// AgentTask Agent任务
type AgentTask struct {
	ID          uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Type        string     `gorm:"column:type;type:varchar(50);comment:'任务类型:inspection/incident/prediction/custom';NOT NULL;index" json:"type"`
	Name        string     `gorm:"column:name;type:varchar(200);comment:'任务名称';NOT NULL" json:"name"`
	Description string     `gorm:"column:description;type:varchar(500);comment:'任务描述'" json:"description"`
	Priority    int        `gorm:"column:priority;comment:'优先级:1-10';default:5" json:"priority"`
	Status      int        `gorm:"column:status;comment:'状态:1=待执行,2=执行中,3=成功,4=失败,5=取消';NOT NULL;index" json:"status"`
	Params      string     `gorm:"column:params;type:text;comment:'任务参数(JSON)'" json:"params"`
	Result      string     `gorm:"column:result;type:text;comment:'执行结果(JSON)'" json:"result"`
	ErrorMsg    string     `gorm:"column:error_msg;type:text;comment:'错误信息'" json:"errorMsg"`
	Retry       int        `gorm:"column:retry;comment:'重试次数';default:0" json:"retry"`
	MaxRetry    int        `gorm:"column:max_retry;comment:'最大重试次数';default:3" json:"maxRetry"`
	StartTime   util.HTime `gorm:"column:start_time;comment:'开始时间'" json:"startTime"`
	EndTime     util.HTime `gorm:"column:end_time;comment:'结束时间'" json:"endTime"`
	CreateTime  util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime  util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (AgentTask) TableName() string {
	return "ai_agent_task"
}

// KnowledgeBase 知识库
type KnowledgeBase struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Type       string     `gorm:"column:type;type:varchar(50);comment:'类型:command/solution/best-practice/faq';NOT NULL;index" json:"type"`
	Category   string     `gorm:"column:category;type:varchar(50);comment:'分类';NOT NULL;index" json:"category"`
	Title      string     `gorm:"column:title;type:varchar(500);comment:'标题';NOT NULL" json:"title"`
	Content    string     `gorm:"column:content;type:text;comment:'内容';NOT NULL" json:"content"`
	Keywords   string     `gorm:"column:keywords;type:varchar(1000);comment:'关键词(JSON数组)'" json:"keywords"`
	Tags       string     `gorm:"column:tags;type:varchar(500);comment:'标签(JSON数组)'" json:"tags"`
	Score      float64    `gorm:"column:score;type:decimal(3,2);comment:'评分';default:0.5" json:"score"`
	UseCount   int        `gorm:"column:use_count;comment:'使用次数';default:0" json:"useCount"`
	Enabled    int        `gorm:"column:enabled;comment:'是否启用:1=启用,2=禁用';default:1" json:"enabled"`
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (KnowledgeBase) TableName() string {
	return "knowledge_base"
}

// PromptTemplate Prompt模板
type PromptTemplate struct {
	ID           uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Name         string     `gorm:"column:name;type:varchar(200);comment:'模板名称';NOT NULL;uniqueIndex" json:"name"`
	Category     string     `gorm:"column:category;type:varchar(50);comment:'分类:inspection/incident/prediction/general';NOT NULL;index" json:"category"`
	Scene        string     `gorm:"column:scene;type:varchar(50);comment:'场景:summary/analysis/suggestion/fix';NOT NULL;index" json:"scene"`
	Template     string     `gorm:"column:template;type:text;comment:'模板内容(支持变量替换)';NOT NULL" json:"template"`
	Variables    string     `gorm:"column:variables;type:varchar(1000);comment:'变量列表(JSON数组)'" json:"variables"`
	SystemPrompt string     `gorm:"column:system_prompt;type:text;comment:'系统提示词'" json:"systemPrompt"`
	Temperature  float64    `gorm:"column:temperature;type:decimal(3,2);comment:'温度参数';default:0.7" json:"temperature"`
	MaxTokens    int        `gorm:"column:max_tokens;comment:'最大Token数';default:2000" json:"maxTokens"`
	ModelID      uint       `gorm:"column:model_id;comment:'默认模型ID'" json:"modelId"`
	Enabled      int        `gorm:"column:enabled;comment:'是否启用:1=启用,2=禁用';default:1" json:"enabled"`
	CreateTime   util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime   util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (PromptTemplate) TableName() string {
	return "prompt_template"
}

// AgentChatHistory Agent对话历史
type AgentChatHistory struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	SessionID  string     `gorm:"column:session_id;type:varchar(64);comment:'会话ID';NOT NULL;index" json:"sessionId"`
	UserID     uint       `gorm:"column:user_id;comment:'用户ID';NOT NULL;index" json:"userId"`
	Role       string     `gorm:"column:role;type:varchar(20);comment:'角色:user/assistant/system';NOT NULL" json:"role"`
	Message    string     `gorm:"column:message;type:text;comment:'消息内容';NOT NULL" json:"message"`
	Intent     string     `gorm:"column:intent;type:varchar(50);comment:'识别的意图类型'" json:"intent"`
	IntentConf float64    `gorm:"column:intent_conf;type:decimal(3,2);comment:'意图识别置信度'" json:"intentConf"`
	Entities   string     `gorm:"column:entities;type:json;comment:'提取的实体(JSON)'" json:"entities"`
	TaskID     uint       `gorm:"column:task_id;comment:'关联的任务ID';index" json:"taskId"`
	TaskType   string     `gorm:"column:task_type;type:varchar(50);comment:'任务类型'" json:"taskType"`
	Status     int        `gorm:"column:status;type:tinyint;comment:'状态:1=处理中,2=成功,3=失败';default:1" json:"status"`
	ErrorMsg   string     `gorm:"column:error_msg;type:text;comment:'错误信息'" json:"errorMsg"`
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL;autoCreateTime" json:"createTime"`
}

func (AgentChatHistory) TableName() string {
	return "ai_agent_chat_history"
}

// AgentTaskDto Agent任务DTO
type AgentTaskDto struct {
	ID          uint                   `json:"id"`
	Type        string                 `json:"type" validate:"required"`
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description"`
	Priority    int                    `json:"priority" validate:"min=1,max=10"`
	Params      map[string]interface{} `json:"params"`
	MaxRetry    int                    `json:"maxRetry"`
}

// AgentTaskQuery Agent任务查询条件
type AgentTaskQuery struct {
	Type      string `form:"type"`
	Status    int    `form:"status"`
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	PageNum   int    `form:"pageNum"`
	PageSize  int    `form:"pageSize"`
}

// KnowledgeBaseDto 知识库DTO
type KnowledgeBaseDto struct {
	ID       uint     `json:"id"`
	Type     string   `json:"type" validate:"required,oneof=command solution best-practice faq"`
	Category string   `json:"category" validate:"required"`
	Title    string   `json:"title" validate:"required"`
	Content  string   `json:"content" validate:"required"`
	Keywords []string `json:"keywords"`
	Tags     []string `json:"tags"`
	Score    float64  `json:"score" validate:"min=0,max=1"`
	Enabled  int      `json:"enabled" validate:"oneof=1 2"`
}

// KnowledgeBaseQuery 知识库查询条件
type KnowledgeBaseQuery struct {
	Type     string `form:"type"`
	Category string `form:"category"`
	Title    string `form:"title"`
	Keyword  string `form:"keyword"`
	Tag      string `form:"tag"`
	Enabled  int    `form:"enabled"`
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
}

// PromptTemplateDto Prompt模板DTO
type PromptTemplateDto struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name" validate:"required"`
	Category     string   `json:"category" validate:"required"`
	Scene        string   `json:"scene" validate:"required"`
	Template     string   `json:"template" validate:"required"`
	Variables    []string `json:"variables"`
	SystemPrompt string   `json:"systemPrompt"`
	Temperature  float64  `json:"temperature" validate:"min=0,max=2"`
	MaxTokens    int      `json:"maxTokens" validate:"min=1"`
	ModelID      uint     `json:"modelId"`
	Enabled      int      `json:"enabled" validate:"oneof=1 2"`
}

// PromptTemplateQuery Prompt模板查询条件
type PromptTemplateQuery struct {
	Name     string `form:"name"`
	Category string `form:"category"`
	Scene    string `form:"scene"`
	Enabled  int    `form:"enabled"`
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
}

// AgentChatRequest Agent对话请求
type AgentChatRequest struct {
	Message               string `json:"message" binding:"required"` // 用户消息
	SessionID             string `json:"sessionId"`                  // 会话ID（可选，用于多轮对话）
	ResolvedTargetType    string `json:"resolvedTargetType,omitempty"`
	ResolvedTargetName    string `json:"resolvedTargetName,omitempty"`
	ResolvedCheckType     string `json:"resolvedCheckType,omitempty"`
	ResolvedQueryType     string `json:"resolvedQueryType,omitempty"`
	ResolvedDirectoryPath string `json:"resolvedDirectoryPath,omitempty"`
	PreferResolvedTarget  bool   `json:"preferResolvedTarget,omitempty"`
}

// AgentChatResponse Agent对话响应
type AgentChatResponse struct {
	SessionID string                 `json:"sessionId"`          // 会话ID
	Message   string                 `json:"message"`            // AI回复
	Type      string                 `json:"type"`               // 响应类型：text/task/error
	Intent    *IntentResult          `json:"intent,omitempty"`   // 意图识别结果
	TaskInfo  map[string]interface{} `json:"taskInfo,omitempty"` // 任务信息
}

// IntentResult 意图识别结果
type IntentResult struct {
	Type       string            `json:"type"`       // 意图类型：inspection/incident/prediction/query/general
	Action     string            `json:"action"`     // 具体动作
	Confidence float64           `json:"confidence"` // 置信度（0-1）
	Entities   map[string]string `json:"entities"`   // 提取的实体
}

// ChatHistoryQuery 对话历史查询
type ChatHistoryQuery struct {
	SessionID string `json:"sessionId"` // 会话ID
	UserID    uint   `json:"userId"`    // 用户ID
	PageNum   int    `json:"pageNum"`   // 页码
	PageSize  int    `json:"pageSize"`  // 每页数量
}

// PressureAnalysisResult 压力分析结果
type PressureAnalysisResult struct {
	HostID        uint                     `json:"hostId"`
	HostName      string                   `json:"hostName"`
	HostIP        string                   `json:"hostIp"`
	CollectTime   time.Time                `json:"collectTime"`
	CPUInfo       map[string]interface{}   `json:"cpuInfo"`
	MemoryInfo    map[string]interface{}   `json:"memoryInfo"`
	DiskInfo      map[string]interface{}   `json:"diskInfo"`
	LoadInfo      map[string]interface{}   `json:"loadInfo"`
	TopProcesses  []map[string]interface{} `json:"topProcesses"`
	NetworkInfo   map[string]interface{}   `json:"networkInfo"`
	PressureLevel string                   `json:"pressureLevel"` // 正常/中等压力/高压力
	Errors        []string                 `json:"errors"`
}

// MemoryAnalysisResult 内存分析结果
type MemoryAnalysisResult struct {
	HostID        uint                     `json:"hostId"`
	HostName      string                   `json:"hostName"`
	HostIP        string                   `json:"hostIp"`
	CollectTime   time.Time                `json:"collectTime"`
	MemoryInfo    map[string]interface{}   `json:"memoryInfo"`
	TopProcesses  []map[string]interface{} `json:"topProcesses"`
	MemoryPattern string                   `json:"memoryPattern"` // 内存使用模式描述
	Issues        []string                 `json:"issues"`        // 发现的问题列表
}
