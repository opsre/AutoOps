package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	aiModel "dodevops-api/api/ai/model"
	aiService "dodevops-api/api/ai/service"
	"dodevops-api/api/agent/dao"
	"dodevops-api/api/agent/model"
)

// AIDecisionEngine AI决策引擎
type AIDecisionEngine struct {
	aiModelService aiService.AiModelService
	promptService  PromptTemplateService
	cache          *PromptCache
	mu             sync.RWMutex
}

var (
	aiEngineInstance *AIDecisionEngine
	aiEngineOnce     sync.Once
)

// GetAIDecisionEngine 获取AI决策引擎单例
func GetAIDecisionEngine() *AIDecisionEngine {
	aiEngineOnce.Do(func() {
		aiEngineInstance = &AIDecisionEngine{
			aiModelService: aiService.NewAiModelService(),
			promptService:  NewPromptTemplateService(),
			cache:          NewPromptCache(),
		}
	})
	return aiEngineInstance
}

// Analyze 智能分析
func (e *AIDecisionEngine) Analyze(category string, scene string, data map[string]interface{}) (string, error) {
	// 1. 加载Prompt模板
	template, err := e.promptService.GetTemplate(category, scene)
	if err != nil {
		return "", fmt.Errorf("加载模板失败: %v", err)
	}

	// 2. 构建Prompt
	builder := NewPromptBuilder(template)
	builder.SetVars(data)
	prompt := builder.Build()

	// 3. 检查缓存
	cacheKey := e.cache.GenerateKey(prompt)
	if cached := e.cache.Get(cacheKey); cached != "" {
		return cached, nil
	}

	// 4. 调用AI模型
	response, err := e.callAI(builder.GetSystemPrompt(), prompt, builder.GetTemperature(), builder.GetMaxTokens(), builder.GetModelID())
	if err != nil {
		return "", fmt.Errorf("AI调用失败: %v", err)
	}

	// 5. 缓存结果
	e.cache.Set(cacheKey, response, 1*time.Hour)

	return response, nil
}

// AnalyzeWithTemplate 使用指定模板进行分析
func (e *AIDecisionEngine) AnalyzeWithTemplate(templateName string, data map[string]interface{}) (string, error) {
	// 1. 根据模板名称获取模板
	templateDao := dao.NewPromptTemplateDao()
	template, err := templateDao.GetByName(templateName)
	if err != nil {
		return "", fmt.Errorf("模板不存在: %s", templateName)
	}

	// 2. 构建Prompt
	builder := NewPromptBuilder(template)
	builder.SetVars(data)
	prompt := builder.Build()

	// 3. 检查缓存
	cacheKey := e.cache.GenerateKey(prompt)
	if cached := e.cache.Get(cacheKey); cached != "" {
		return cached, nil
	}

	// 4. 调用AI模型
	response, err := e.callAI(builder.GetSystemPrompt(), prompt, builder.GetTemperature(), builder.GetMaxTokens(), builder.GetModelID())
	if err != nil {
		return "", fmt.Errorf("AI调用失败: %v", err)
	}

	// 5. 缓存结果
	e.cache.Set(cacheKey, response, 1*time.Hour)

	return response, nil
}

// AnalyzeJSON 智能分析并解析JSON结果
func (e *AIDecisionEngine) AnalyzeJSON(category string, scene string, data map[string]interface{}, result interface{}) error {
	response, err := e.Analyze(category, scene, data)
	if err != nil {
		return err
	}

	// 尝试提取JSON（可能包含在代码块中）
	jsonStr := extractJSON(response)
	if jsonStr == "" {
		jsonStr = response
	}

	// 解析JSON
	if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
		return fmt.Errorf("解析AI响应失败: %v, 响应内容: %s", err, response)
	}

	return nil
}

// callAI 调用AI模型
func (e *AIDecisionEngine) callAI(systemPrompt string, userPrompt string, temperature float64, maxTokens int, modelID uint) (string, error) {
	// 默认使用第一个启用的模型
	if modelID == 0 {
		// 直接从DAO层获取启用的模型，避免依赖Gin Context
		modelID = 1 // 默认使用ID为1的模型，后续可从配置获取
	}

	// 构建消息列表
	messages := []aiModel.AiChatMessage{}
	if systemPrompt != "" {
		messages = append(messages, aiModel.AiChatMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}
	messages = append(messages, aiModel.AiChatMessage{
		Role:    "user",
		Content: userPrompt,
	})

	// 构建请求
	req := &aiModel.AiChatRequest{
		ModelID:  modelID,
		Messages: messages,
		Parameters: map[string]interface{}{
			"temperature": temperature,
			"max_tokens":  maxTokens,
		},
	}

	// 由于AI服务的Chat方法需要gin.Context，这里需要直接调用底层实现
	// 或者创建一个内部方法来处理。暂时返回错误提示需要重构
	// TODO: 需要重构AI服务以支持非Gin Context调用
	_ = req
	return "", fmt.Errorf("AI调用暂不可用，需要重构AI服务接口")
}

// extractJSON 从文本中提取JSON
func extractJSON(text string) string {
	// 尝试提取 ```json ... ``` 中的内容
	start := -1
	end := -1

	// 查找 ```json
	if idx := findPattern(text, "```json"); idx >= 0 {
		start = idx + 7 // 跳过 ```json
		// 查找结束的 ```
		if idx := findPattern(text[start:], "```"); idx >= 0 {
			end = start + idx
		}
	}

	// 如果没找到，尝试查找 { ... }
	if start == -1 || end == -1 {
		for i := 0; i < len(text); i++ {
			if text[i] == '{' {
				start = i
				break
			}
		}
		if start >= 0 {
			for i := len(text) - 1; i >= 0; i-- {
				if text[i] == '}' {
					end = i + 1
					break
				}
			}
		}
	}

	if start >= 0 && end > start {
		return text[start:end]
	}

	return ""
}

// findPattern 查找子串
func findPattern(text string, pattern string) int {
	for i := 0; i <= len(text)-len(pattern); i++ {
		match := true
		for j := 0; j < len(pattern); j++ {
			if text[i+j] != pattern[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// PromptCache Prompt缓存
type PromptCache struct {
	cache map[string]*CacheItem
	mu    sync.RWMutex
}

// CacheItem 缓存项
type CacheItem struct {
	Value      string
	ExpireTime time.Time
}

// NewPromptCache 创建Prompt缓存
func NewPromptCache() *PromptCache {
	pc := &PromptCache{
		cache: make(map[string]*CacheItem),
	}

	// 启动清理协程
	go pc.cleanupLoop()

	return pc
}

// GenerateKey 生成缓存键
func (pc *PromptCache) GenerateKey(prompt string) string {
	// 简单的哈希函数
	hash := uint32(0)
	for i := 0; i < len(prompt); i++ {
		hash = hash*31 + uint32(prompt[i])
	}
	return fmt.Sprintf("%x", hash)
}

// Get 获取缓存
func (pc *PromptCache) Get(key string) string {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	item, exists := pc.cache[key]
	if !exists {
		return ""
	}

	if time.Now().After(item.ExpireTime) {
		return ""
	}

	return item.Value
}

// Set 设置缓存
func (pc *PromptCache) Set(key string, value string, duration time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.cache[key] = &CacheItem{
		Value:      value,
		ExpireTime: time.Now().Add(duration),
	}
}

// Delete 删除缓存
func (pc *PromptCache) Delete(key string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	delete(pc.cache, key)
}

// Clear 清空缓存
func (pc *PromptCache) Clear() {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.cache = make(map[string]*CacheItem)
}

// cleanupLoop 清理过期缓存
func (pc *PromptCache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		pc.cleanup()
	}
}

// cleanup 清理过期缓存
func (pc *PromptCache) cleanup() {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	now := time.Now()
	for key, item := range pc.cache {
		if now.After(item.ExpireTime) {
			delete(pc.cache, key)
		}
	}
}

// IntentRecognizer 意图识别器
type IntentRecognizer struct {
	aiEngine *AIDecisionEngine
}

// NewIntentRecognizer 创建意图识别器
func NewIntentRecognizer() *IntentRecognizer {
	return &IntentRecognizer{
		aiEngine: GetAIDecisionEngine(),
	}
}

// Recognize 识别意图
func (ir *IntentRecognizer) Recognize(message string) (*model.IntentResult, error) {
	prompt := fmt.Sprintf(`请识别以下用户消息的意图：
"%s"

可能的意图类型：
1. inspection - 运维巡检（例如：帮我检查一下服务器、巡检所有主机）
2. incident - 故障分析（例如：分析故障、服务器出问题了）
3. prediction - 容量预警（例如：磁盘什么时候会满、预测资源使用）
4. query - 信息查询（例如：查询主机状态、显示CPU使用率）
5. general - 一般对话

请返回JSON格式：
{
    "type": "意图类型",
    "action": "具体动作",
    "confidence": 0.0-1.0,
    "entities": {}
}`, message)

	response, err := ir.aiEngine.callAI(
		"你是一个智能运维助手，擅长理解用户的运维需求。",
		prompt,
		0.3, // 低温度，保证结果稳定
		500,
		0,
	)
	if err != nil {
		return nil, err
	}

	// 解析JSON
	var intent model.IntentResult
	jsonStr := extractJSON(response)
	if jsonStr == "" {
		jsonStr = response
	}

	if err := json.Unmarshal([]byte(jsonStr), &intent); err != nil {
		// 如果解析失败，返回默认意图
		return &model.IntentResult{
			Type:       "general",
			Action:     "chat",
			Confidence: 0.5,
			Entities:   make(map[string]string),
		}, nil
	}

	return &intent, nil
}
