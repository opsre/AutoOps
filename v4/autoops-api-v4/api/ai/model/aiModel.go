// AI模型管理相关结构体
package model

import (
	"dodevops-api/common/util"
	"encoding/json"
)

// AiModel AI模型配置
type AiModel struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Name       string     `gorm:"column:name;type:varchar(100);comment:'模型名称';NOT NULL" json:"name"`
	Type       string     `gorm:"column:type;type:varchar(50);comment:'模型类型(阿里百炼,OpenAI等)';NOT NULL" json:"type"`
	URL        string     `gorm:"column:url;type:varchar(500);comment:'API地址';NOT NULL" json:"url"`
	APIKey     string     `gorm:"column:api_key;type:varchar(200);comment:'API密钥';NOT NULL" json:"apiKey"`
	Model      string     `gorm:"column:model;type:varchar(100);comment:'模型标识(qwen-plus,gpt-4等)';NOT NULL" json:"model"`
	Status     int        `gorm:"column:status;default:1;comment:'状态:1->启用,2->禁用';NOT NULL" json:"status"`
	Remark     string     `gorm:"column:remark;type:varchar(500);comment:'备注'" json:"remark"`
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (AiModel) TableName() string {
	return "ai_model"
}

// AiModelDto AI模型DTO
type AiModelDto struct {
	ID     uint   `json:"id"`
	Name   string `json:"name" validate:"required"`   // 模型名称
	Type   string `json:"type" validate:"required"`   // 模型类型
	URL    string `json:"url" validate:"required"`    // API地址
	APIKey string `json:"apiKey" validate:"required"` // API密钥
	Model  string `json:"model" validate:"required"`  // 模型标识
	Status int    `json:"status"`                     // 状态
	Remark string `json:"remark"`                     // 备注
}

// AiModelQuery AI模型查询条件
type AiModelQuery struct {
	Name     string `form:"name"`     // 模型名称
	Type     string `form:"type"`     // 模型类型
	Status   int    `form:"status"`   // 状态(-1:全部,1:启用,2:禁用)
	PageNum  int    `form:"pageNum"`  // 页码
	PageSize int    `form:"pageSize"` // 每页数量
}

// AiChatRequest AI对话请求
type AiChatRequest struct {
	ModelID    uint                   `json:"modelId" validate:"required"`  // 模型ID
	Messages   []AiChatMessage        `json:"messages" validate:"required"` // 消息列表
	Stream     bool                   `json:"stream"`                       // 是否流式返回
	Parameters map[string]interface{} `json:"parameters"`                   // 额外参数(温度、top_p等)
	SessionID  string                 `json:"sessionId,omitempty"`          // 会话ID（内存记忆，不落库）
}

// UnmarshalJSON 兼容 sessionId / sessionID 两种前端字段写法
func (r *AiChatRequest) UnmarshalJSON(data []byte) error {
	type Alias AiChatRequest
	aux := struct {
		Alias
		SessionIDAlt string `json:"sessionID"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*r = AiChatRequest(aux.Alias)
	if r.SessionID == "" {
		r.SessionID = aux.SessionIDAlt
	}
	return nil
}

// AiChatMessage AI对话消息
type AiChatMessage struct {
	Role    string `json:"role" validate:"required,oneof=user assistant system"` // 角色
	Content string `json:"content" validate:"required"`                          // 内容
}

// AliyunChatRequest 阿里百炼请求格式
type AliyunChatRequest struct {
	Model      string                 `json:"model"`
	Input      AliyunInput            `json:"input"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// AliyunInput 阿里百炼输入格式
type AliyunInput struct {
	Messages []AliChatMessage `json:"messages"`
}

// AliChatMessage 阿里百炼消息格式
type AliChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AliyunChatResponse 阿里百炼响应格式
type AliyunChatResponse struct {
	Output    AliyunOutput `json:"output"`
	Usage     AliyunUsage  `json:"usage"`
	RequestID string       `json:"request_id"`
}

// AliyunOutput 阿里百炼输出
type AliyunOutput struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

// AliyunUsage 阿里百炼使用情况
type AliyunUsage struct {
	OutputTokens int `json:"output_tokens"`
	InputTokens  int `json:"input_tokens"`
}

// AiChatResponse AI对话响应（结构化返回）
type AiChatResponse struct {
	Text         string         `json:"text"`                // 原始文本
	Content      []ContentBlock `json:"content"`             // 结构化内容块
	FinishReason string         `json:"finishReason"`        // 完成原因
	Usage        AliyunUsage    `json:"usage"`               // Token使用情况
	RequestID    string         `json:"requestId"`           // 请求ID
	SessionID    string         `json:"sessionId,omitempty"` // 会话ID，便于前端复用
}

// ContentBlock 内容块（区分文本、代码等）
type ContentBlock struct {
	Type     string `json:"type"`               // 类型: text/code/heading/list
	Content  string `json:"content"`            // 内容
	Language string `json:"language,omitempty"` // 代码语言（仅code类型）
	Level    int    `json:"level,omitempty"`    // 标题级别（仅heading类型）
}

// OpenAIChatRequest OpenAI/DeepSeek请求格式
type OpenAIChatRequest struct {
	Model       string              `json:"model"`
	Messages    []OpenAIChatMessage `json:"messages"`
	Stream      bool                `json:"stream,omitempty"`
	Temperature float64             `json:"temperature,omitempty"`
	TopP        float64             `json:"top_p,omitempty"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
}

// OpenAIChatMessage OpenAI/DeepSeek消息格式
type OpenAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIChatResponse OpenAI/DeepSeek响应格式
type OpenAIChatResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
	Usage   OpenAIUsage    `json:"usage"`
}

// OpenAIChoice OpenAI/DeepSeek选项
type OpenAIChoice struct {
	Index        int               `json:"index"`
	Message      OpenAIChatMessage `json:"message"`
	FinishReason string            `json:"finish_reason"`
}

// OpenAIUsage OpenAI/DeepSeek使用情况
type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
