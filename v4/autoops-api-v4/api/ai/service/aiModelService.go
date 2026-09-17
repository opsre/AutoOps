// AI模型服务层
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"dodevops-api/api/ai/dao"
	"dodevops-api/api/ai/model"
	"dodevops-api/common/result"
	redisPkg "dodevops-api/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AiModelService AI模型服务接口
type AiModelService interface {
	CreateModel(c *gin.Context, dto *model.AiModelDto)
	GetModelByID(c *gin.Context, id uint)
	UpdateModel(c *gin.Context, dto *model.AiModelDto)
	DeleteModel(c *gin.Context, id uint)
	GetModelList(c *gin.Context, query *model.AiModelQuery)
	UpdateModelStatus(c *gin.Context, id uint, status int)
	GetEnabledModels(c *gin.Context)

	// AI对话功能
	Chat(c *gin.Context, req *model.AiChatRequest)
}

// AiModelServiceImpl AI模型服务实现
type AiModelServiceImpl struct {
	modelDao dao.AiModelDao
}

// NewAiModelService 创建AI模型服务实例
func NewAiModelService() AiModelService {
	return &AiModelServiceImpl{
		modelDao: dao.NewAiModelDao(),
	}
}

// ===== 内存会话存储（不落库） =====
type memoryStore struct {
	mu   sync.RWMutex
	data map[string][]model.AiChatMessage
}

type SessionStructuredContext struct {
	LastTargetType    string `json:"lastTargetType,omitempty"`
	LastTargetName    string `json:"lastTargetName,omitempty"`
	LastCheckType     string `json:"lastCheckType,omitempty"`
	LastQueryType     string `json:"lastQueryType,omitempty"`
	LastDirectoryPath string `json:"lastDirectoryPath,omitempty"`
	UpdatedAt         string `json:"updatedAt,omitempty"`
}

func newMemoryStore() *memoryStore {
	return &memoryStore{data: make(map[string][]model.AiChatMessage)}
}

func (s *memoryStore) Get(sessionID string) []model.AiChatMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if msgs, ok := s.data[sessionID]; ok {
		cp := make([]model.AiChatMessage, len(msgs))
		copy(cp, msgs)
		return cp
	}
	return nil
}

func (s *memoryStore) Set(sessionID string, msgs []model.AiChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]model.AiChatMessage, len(msgs))
	copy(cp, msgs)
	s.data[sessionID] = cp
}

var memStore = newMemoryStore()
var contextStore = struct {
	mu   sync.RWMutex
	data map[string]SessionStructuredContext
}{data: make(map[string]SessionStructuredContext)}

// AppendMemory 供控制器或其他模块将消息写入同一会话的内存上下文
func AppendMemory(sessionID string, msgs ...model.AiChatMessage) {
	if sessionID == "" || len(msgs) == 0 {
		return
	}
	current := memStore.Get(sessionID)
	current = append(current, msgs...)
	current = trimToLastRounds(current, 6)
	memStore.Set(sessionID, current)
}

// GetMemory 供读取当前会话的内存上下文（只读拷贝）
func GetMemory(sessionID string) []model.AiChatMessage {
	return memStore.Get(sessionID)
}

func GetSessionStructuredContext(sessionID string) SessionStructuredContext {
	if sessionID == "" {
		return SessionStructuredContext{}
	}
	if redisPkg.RedisDb != nil {
		ctx := context.Background()
		val, err := redisPkg.RedisDb.Get(ctx, "ai:ctx:"+sessionID).Result()
		if err == nil && val != "" {
			var stored SessionStructuredContext
			if json.Unmarshal([]byte(val), &stored) == nil {
				_ = redisPkg.RedisDb.Expire(ctx, "ai:ctx:"+sessionID, 30*time.Minute).Err()
				return stored
			}
		}
	}

	contextStore.mu.RLock()
	defer contextStore.mu.RUnlock()
	return contextStore.data[sessionID]
}

func SaveSessionStructuredContext(sessionID string, structured SessionStructuredContext) {
	if sessionID == "" {
		return
	}
	structured.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	if redisPkg.RedisDb != nil {
		ctx := context.Background()
		if b, err := json.Marshal(structured); err == nil {
			_ = redisPkg.RedisDb.Set(ctx, "ai:ctx:"+sessionID, b, 30*time.Minute).Err()
		}
	}
	contextStore.mu.Lock()
	defer contextStore.mu.Unlock()
	contextStore.data[sessionID] = structured
}

// ===== Redis 会话存储（可选，TTL自动过期） =====
type redisSessionStore struct{}

func (r *redisSessionStore) key(sessionID string) string { return "ai:session:" + sessionID }

func (r *redisSessionStore) Get(sessionID string) ([]model.AiChatMessage, bool) {
	if redisPkg.RedisDb == nil {
		return nil, false
	}
	ctx := context.Background()
	val, err := redisPkg.RedisDb.Get(ctx, r.key(sessionID)).Result()
	if err != nil {
		return nil, false
	}
	var msgs []model.AiChatMessage
	if err := json.Unmarshal([]byte(val), &msgs); err != nil {
		return nil, false
	}
	// 滑动刷新TTL
	_ = redisPkg.RedisDb.Expire(ctx, r.key(sessionID), 5*time.Minute).Err()
	return msgs, true
}

func (r *redisSessionStore) Set(sessionID string, msgs []model.AiChatMessage) error {
	if redisPkg.RedisDb == nil {
		return fmt.Errorf("redis not initialized")
	}
	ctx := context.Background()
	b, err := json.Marshal(msgs)
	if err != nil {
		return err
	}
	return redisPkg.RedisDb.Set(ctx, r.key(sessionID), b, 5*time.Minute).Err()
}

func useRedis() bool { return redisPkg.RedisDb != nil }

// CreateModel 创建AI模型
func (s *AiModelServiceImpl) CreateModel(c *gin.Context, dto *model.AiModelDto) {
	// 检查模型名称是否已存在
	existing, _ := s.modelDao.GetByName(dto.Name)
	if existing != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "模型名称已存在")
		return
	}

	aiModel := &model.AiModel{
		Name:   dto.Name,
		Type:   dto.Type,
		URL:    dto.URL,
		APIKey: dto.APIKey,
		Model:  dto.Model,
		Status: 1, // 默认启用
		Remark: dto.Remark,
	}

	if err := s.modelDao.Create(aiModel); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建失败: %v", err))
		return
	}

	result.Success(c, aiModel)
}

// GetModelByID 根据ID获取AI模型
func (s *AiModelServiceImpl) GetModelByID(c *gin.Context, id uint) {
	aiModel, err := s.modelDao.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "模型不存在")
		} else {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		}
		return
	}

	result.Success(c, aiModel)
}

// UpdateModel 更新AI模型
func (s *AiModelServiceImpl) UpdateModel(c *gin.Context, dto *model.AiModelDto) {
	// 检查模型是否存在
	aiModel, err := s.modelDao.GetByID(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "模型不存在")
		return
	}

	// 如果修改了模型名称，检查新名称是否已存在
	if dto.Name != aiModel.Name {
		existing, _ := s.modelDao.GetByName(dto.Name)
		if existing != nil && existing.ID != dto.ID {
			result.Failed(c, int(result.ApiCode.FAILED), "模型名称已存在")
			return
		}
	}

	// 更新字段
	aiModel.Name = dto.Name
	aiModel.Type = dto.Type
	aiModel.URL = dto.URL
	aiModel.APIKey = dto.APIKey
	aiModel.Model = dto.Model
	aiModel.Status = dto.Status
	aiModel.Remark = dto.Remark

	if err := s.modelDao.Update(aiModel); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新失败: %v", err))
		return
	}

	result.Success(c, aiModel)
}

// DeleteModel 删除AI模型
func (s *AiModelServiceImpl) DeleteModel(c *gin.Context, id uint) {
	// 检查是否存在
	_, err := s.modelDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "模型不存在")
		return
	}

	if err := s.modelDao.Delete(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("删除失败: %v", err))
		return
	}

	result.Success(c, nil)
}

// GetModelList 获取AI模型列表
func (s *AiModelServiceImpl) GetModelList(c *gin.Context, query *model.AiModelQuery) {
	aiModels, total, err := s.modelDao.GetList(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		return
	}

	data := map[string]interface{}{
		"list":     aiModels,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	}

	result.Success(c, data)
}

// UpdateModelStatus 更新模型状态
func (s *AiModelServiceImpl) UpdateModelStatus(c *gin.Context, id uint, status int) {
	// 检查是否存在
	_, err := s.modelDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "模型不存在")
		return
	}

	if err := s.modelDao.UpdateStatus(id, status); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新状态失败: %v", err))
		return
	}

	result.Success(c, nil)
}

// GetEnabledModels 获取所有启用的模型
func (s *AiModelServiceImpl) GetEnabledModels(c *gin.Context) {
	aiModels, err := s.modelDao.GetEnabledModels()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		return
	}

	result.Success(c, aiModels)
}

// Chat AI对话
func (s *AiModelServiceImpl) Chat(c *gin.Context, req *model.AiChatRequest) {
	// 1) 生成/获取 sessionId（不落库，内存会话）
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	// 2) 注意：历史合并已在 controller 层完成，这里直接使用 req.Messages
	//    统一裁剪：保留所有system + 最近6轮UA
	finalMsgs := trimToLastRounds(req.Messages, 6)
	req.Messages = finalMsgs
	req.SessionID = sessionID

	// 先写入会话存储（到用户发言为止，助手回覆稍后追加）
	if useRedis() {
		_ = (&redisSessionStore{}).Set(sessionID, finalMsgs)
	} else {
		memStore.Set(sessionID, finalMsgs)
	}
	// 获取模型配置
	aiModel, err := s.modelDao.GetByID(req.ModelID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "模型不存在")
		return
	}

	if aiModel.Status != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "模型未启用")
		return
	}

	// 根据模型类型调用不同的API
	switch aiModel.Type {
	case "阿里百炼":
		s.callAliyunAPI(c, aiModel, req, sessionID)
	case "DeepSeek", "OpenAI":
		s.callOpenAIAPI(c, aiModel, req, sessionID)
	default:
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("不支持的模型类型: %s", aiModel.Type))
	}
}

// trimToLastRounds 保留所有system消息 + 最近N轮(user/assistant共2N条)
func trimToLastRounds(msgs []model.AiChatMessage, rounds int) []model.AiChatMessage {
	if rounds <= 0 || len(msgs) == 0 {
		return msgs
	}

	var systems []model.AiChatMessage
	var ua []model.AiChatMessage
	for _, m := range msgs {
		if m.Role == "system" {
			systems = append(systems, m)
			continue
		}
		ua = append(ua, m)
	}

	keep := rounds * 2
	if keep > len(ua) {
		keep = len(ua)
	}
	// 只在需要截断时进行切片，避免切片越界
	var tail []model.AiChatMessage
	if keep > 0 {
		tail = ua[len(ua)-keep:]
	}
	trimmed := make([]model.AiChatMessage, 0, len(systems)+len(tail))
	trimmed = append(trimmed, systems...)
	trimmed = append(trimmed, tail...)
	return trimmed
}

// callAliyunAPI 调用阿里百炼API (DashScope)
func (s *AiModelServiceImpl) callAliyunAPI(c *gin.Context, aiModel *model.AiModel, req *model.AiChatRequest, sessionID string) {
	// 转换消息格式
	messages := make([]model.AliChatMessage, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = model.AliChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 构建阿里百炼请求 (DashScope格式)
	aliyunReq := model.AliyunChatRequest{
		Model: aiModel.Model,
		Input: model.AliyunInput{
			Messages: messages,
		},
		Parameters: req.Parameters,
	}

	// 序列化请求
	reqBody, err := json.Marshal(aliyunReq)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("请求序列化失败: %v", err))
		return
	}

	// 使用正确的DashScope API地址
	// 如果数据库中URL为空或不正确，使用默认地址
	apiURL := aiModel.URL
	if apiURL == "" || apiURL == "https://bailian.cn-beijing.aliyuncs.com/api/v1/chat/completions" {
		// 使用正确的DashScope API地址
		apiURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建请求失败: %v", err))
		return
	}

	// 设置请求头 (DashScope使用Bearer token)
	httpReq.Header.Set("Authorization", "Bearer "+aiModel.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("请求失败: %v", err))
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("读取响应失败: %v", err))
		return
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("API返回错误(状态码: %d, URL: %s): %s", resp.StatusCode, apiURL, string(body)))
		return
	}

	// 解析响应
	var aliyunResp model.AliyunChatResponse
	if err := json.Unmarshal(body, &aliyunResp); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("解析响应失败: %v, 响应内容: %s", err, string(body)))
		return
	}

	// 解析 Markdown 为结构化内容
	contentBlocks := ParseMarkdown(aliyunResp.Output.Text)

	// 返回结构化结果
	response := model.AiChatResponse{
		Text:         aliyunResp.Output.Text,
		Content:      contentBlocks,
		FinishReason: aliyunResp.Output.FinishReason,
		Usage:        aliyunResp.Usage,
		RequestID:    aliyunResp.RequestID,
		SessionID:    sessionID,
	}

	// 将助手回复追加到会话存储，并进行再次裁剪
	if response.Text != "" {
		var current []model.AiChatMessage
		if useRedis() {
			if h, ok := (&redisSessionStore{}).Get(sessionID); ok {
				current = h
			}
		} else {
			current = memStore.Get(sessionID)
		}
		current = append(current, model.AiChatMessage{Role: "assistant", Content: response.Text})
		current = trimToLastRounds(current, 6)
		if useRedis() {
			_ = (&redisSessionStore{}).Set(sessionID, current)
		} else {
			memStore.Set(sessionID, current)
		}
	}

	result.Success(c, response)
}

// callOpenAIAPI 调用OpenAI/DeepSeek API
func (s *AiModelServiceImpl) callOpenAIAPI(c *gin.Context, aiModel *model.AiModel, req *model.AiChatRequest, sessionID string) {
	// 转换消息格式
	messages := make([]model.OpenAIChatMessage, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = model.OpenAIChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 构建OpenAI/DeepSeek请求
	openaiReq := model.OpenAIChatRequest{
		Model:    aiModel.Model,
		Messages: messages,
		Stream:   req.Stream,
	}

	// 处理额外参数
	if req.Parameters != nil {
		if temp, ok := req.Parameters["temperature"].(float64); ok {
			openaiReq.Temperature = temp
		}
		if topP, ok := req.Parameters["top_p"].(float64); ok {
			openaiReq.TopP = topP
		}
		if maxTokens, ok := req.Parameters["max_tokens"].(float64); ok {
			openaiReq.MaxTokens = int(maxTokens)
		}
	}

	// 序列化请求
	reqBody, err := json.Marshal(openaiReq)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("请求序列化失败: %v", err))
		return
	}

	// 构建完整的API URL (OpenAI兼容API需要 /chat/completions 端点)
	apiURL := aiModel.URL
	if apiURL[len(apiURL)-1:] != "/" {
		apiURL += "/"
	}
	apiURL += "chat/completions"

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建请求失败: %v", err))
		return
	}

	// 设置请求头
	httpReq.Header.Set("Authorization", "Bearer "+aiModel.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("请求失败: %v", err))
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("读取响应失败: %v", err))
		return
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("API返回错误(状态码: %d, URL: %s): %s", resp.StatusCode, apiURL, string(body)))
		return
	}

	// 解析响应
	var openaiResp model.OpenAIChatResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("解析响应失败: %v", err))
		return
	}

	// 检查是否有返回内容
	if len(openaiResp.Choices) == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "API未返回任何内容")
		return
	}

	// 获取回复文本
	responseText := openaiResp.Choices[0].Message.Content

	// 解析 Markdown 为结构化内容
	contentBlocks := ParseMarkdown(responseText)

	// 转换Usage格式 (OpenAI格式 -> 统一格式)
	usage := model.AliyunUsage{
		OutputTokens: openaiResp.Usage.CompletionTokens,
		InputTokens:  openaiResp.Usage.PromptTokens,
	}

	// 返回结构化结果
	response := model.AiChatResponse{
		Text:         responseText,
		Content:      contentBlocks,
		FinishReason: openaiResp.Choices[0].FinishReason,
		Usage:        usage,
		RequestID:    openaiResp.ID,
		SessionID:    sessionID,
	}

	// 将助手回复追加到会话存储，并进行再次裁剪
	if response.Text != "" {
		var current []model.AiChatMessage
		if useRedis() {
			if h, ok := (&redisSessionStore{}).Get(sessionID); ok {
				current = h
			}
		} else {
			current = memStore.Get(sessionID)
		}
		current = append(current, model.AiChatMessage{Role: "assistant", Content: response.Text})
		current = trimToLastRounds(current, 6)
		if useRedis() {
			_ = (&redisSessionStore{}).Set(sessionID, current)
		} else {
			memStore.Set(sessionID, current)
		}
	}

	result.Success(c, response)
}
