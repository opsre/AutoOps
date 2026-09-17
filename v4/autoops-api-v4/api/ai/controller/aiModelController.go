// AI模型控制器
package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	agentDao "dodevops-api/api/agent/dao"
	agentModel "dodevops-api/api/agent/model"
	agentService "dodevops-api/api/agent/service"
	"dodevops-api/api/ai/model"
	"dodevops-api/api/ai/service"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"
	"dodevops-api/pkg/log"
	redisPkg "dodevops-api/pkg/redis"

	"github.com/gin-gonic/gin"
)

// @Tags AI模型管理
// @Summary 创建AI模型
// @Description 创建AI模型
// @Accept json
// @Produce json
// @Description 创建新的AI模型配置
// @Param data body model.AiModelDto true "模型信息"
// @Success 200 {object} result.Result
// @Router /api/v1/ai/model/add [post]
// @Security ApiKeyAuth
func CreateAiModel(c *gin.Context) {
	var dto model.AiModelDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	modelService := service.NewAiModelService()
	modelService.CreateModel(c, &dto)
}

// @Tags AI模型管理
// @Summary 获取AI模型详情
// @Description 获取AI模型详情
// @Produce json
// @Description 根据ID获取AI模型详情
// @Param id query int true "模型ID"
// @Success 200 {object} result.Result
// @Router /api/v1/ai/model/info [get]
// @Security ApiKeyAuth
func GetAiModelInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	modelService := service.NewAiModelService()
	modelService.GetModelByID(c, uint(id))
}

// @Tags AI模型管理
// @Summary 更新AI模型
// @Description 更新AI模型
// @Accept json
// @Produce json
// @Description 更新AI模型信息
// @Param data body model.AiModelDto true "模型信息"
// @Success 200 {object} result.Result
// @Router /api/v1/ai/model/update [put]
// @Security ApiKeyAuth
func UpdateAiModel(c *gin.Context) {
	var dto model.AiModelDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	if dto.ID == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "ID不能为空")
		return
	}

	modelService := service.NewAiModelService()
	modelService.UpdateModel(c, &dto)
}

// @Tags AI模型管理
// @Summary 删除AI模型
// @Description 删除AI模型
// @Produce json
// @Description 根据ID删除AI模型
// @Param id query int true "模型ID"
// @Success 200 {object} result.Result
// @Router /api/v1/ai/model/delete [delete]
// @Security ApiKeyAuth
func DeleteAiModel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	modelService := service.NewAiModelService()
	modelService.DeleteModel(c, uint(id))
}

// @Tags AI模型管理
// @Summary 获取AI模型列表
// @Description 获取AI模型列表
// @Produce json
// @Description 获取AI模型列表（支持分页和筛选）
// @Param name query string false "模型名称(模糊查询)"
// @Param type query string false "模型类型"
// @Param status query int false "状态(-1:全部,1:启用,2:禁用)"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/ai/model/list [get]
// @Security ApiKeyAuth
func GetAiModelList(c *gin.Context) {
	query := &model.AiModelQuery{
		Name:     c.Query("name"),
		Type:     c.Query("type"),
		Status:   -1, // 默认查询全部
		PageNum:  1,
		PageSize: 10,
	}

	if pageNum, err := strconv.Atoi(c.Query("pageNum")); err == nil && pageNum > 0 {
		query.PageNum = pageNum
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}
	if status, err := strconv.Atoi(c.Query("status")); err == nil {
		query.Status = status
	}

	modelService := service.NewAiModelService()
	modelService.GetModelList(c, query)
}

// @Tags AI模型管理
// @Summary 更新模型状态
// @Description 更新模型状态
// @Accept json
// @Produce json
// @Description 更新AI模型状态
// @Param data body map[string]interface{} true "状态信息 {id:1, status:2}"
// @Success 200 {object} result.Result
// @Router /api/v1/ai/model/updateStatus [put]
// @Security ApiKeyAuth
func UpdateAiModelStatus(c *gin.Context) {
	var req struct {
		ID     uint `json:"id" validate:"required"`
		Status int  `json:"status" validate:"required,oneof=1 2"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	modelService := service.NewAiModelService()
	modelService.UpdateModelStatus(c, req.ID, req.Status)
}

// @Tags AI模型管理
// @Summary 获取启用的模型列表
// @Description 获取启用的模型列表
// @Produce json
// @Description 获取所有启用状态的AI模型
// @Success 200 {object} result.Result
// @Router /api/v1/ai/model/enabled [get]
// @Security ApiKeyAuth
func GetEnabledAiModels(c *gin.Context) {
	modelService := service.NewAiModelService()
	modelService.GetEnabledModels(c)
}

// ============== AI对话接口 ==============

// @Tags AI对话
// @Summary AI对话
// @Description AI对话
// @Accept json
// @Produce json
// @Description 使用AI模型进行对话（支持智能运维Agent）
// @Param data body model.AiChatRequest true "对话请求"
// @Success 200 {object} result.Result
// @Router /api/v1/ai/chat [post]
// @Security ApiKeyAuth
func AiChat(c *gin.Context) {
	var req model.AiChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	// 【核心修复】：无论走哪个流程，都先合并历史记录
	// 若带有sessionId，合并历史（优先Redis，其次进程内存）以便意图识别使用上下文
	if req.SessionID != "" {
		mergedOnce := false
		if redisPkg.RedisDb != nil {
			if hist := loadSessionHistory(req.SessionID); len(hist) > 0 {
				log.Log().Infof("[会话历史] 从 Redis 加载到 %d 条历史消息，会话ID: %s", len(hist), req.SessionID)
				// 仅当请求端只发了1条非system消息时合并历史
				nonSystem := 0
				for _, m := range req.Messages {
					if m.Role != "system" {
						nonSystem++
					}
				}
				if nonSystem <= 1 {
					merged := make([]model.AiChatMessage, 0, len(hist)+len(req.Messages))
					merged = append(merged, hist...)
					merged = append(merged, req.Messages...)
					req.Messages = merged
					mergedOnce = true
					log.Log().Infof("[会话历史] 已合并历史，总消息数: %d", len(req.Messages))
				}
			}
		}
		// 若未从Redis合并，尝试使用服务内存会话
		if !mergedOnce {
			if hist := service.GetMemory(req.SessionID); len(hist) > 0 {
				log.Log().Infof("[会话历史] 从内存加载到 %d 条历史消息，会话ID: %s", len(hist), req.SessionID)
				nonSystem := 0
				for _, m := range req.Messages {
					if m.Role != "system" {
						nonSystem++
					}
				}
				if nonSystem <= 1 {
					merged := make([]model.AiChatMessage, 0, len(hist)+len(req.Messages))
					merged = append(merged, hist...)
					merged = append(merged, req.Messages...)
					req.Messages = merged
					log.Log().Infof("[会话历史] 已合并历史，总消息数: %d", len(req.Messages))
				}
			} else {
				dbHist := loadSessionHistoryFromDB(req.SessionID)
				if len(dbHist) > 0 {
					log.Log().Infof("[会话历史] 从数据库加载到 %d 条历史消息，会话ID: %s", len(dbHist), req.SessionID)
					nonSystem := 0
					for _, m := range req.Messages {
						if m.Role != "system" {
							nonSystem++
						}
					}
					if nonSystem <= 1 {
						merged := make([]model.AiChatMessage, 0, len(dbHist)+len(req.Messages))
						merged = append(merged, dbHist...)
						merged = append(merged, req.Messages...)
						req.Messages = merged
						log.Log().Infof("[会话历史] 已合并数据库历史，总消息数: %d", len(req.Messages))
					}
				} else {
					log.Log().Warnf("[会话历史] 未找到会话历史，会话ID: %s", req.SessionID)
				}
			}
		}
	}

	// 尝试从消息中提取用户意图
	// 【优先判断】：如果是询问历史记录的问题（而非执行新操作），则走通用对话流程
	if isHistoryQueryIntent(req) {
		log.Log().Info("[意图识别] 检测到历史查询意图，走通用对话流程")
		// 走通用对话流程（会添加历史摘要辅助回答）
		handleGeneralChat(c, &req)
		return
	}

	// 如果是运维相关的意图（需要执行实际操作），转发到Agent服务处理
	if isAgentIntent(req) {
		log.Log().Info("[意图识别] 检测到 Agent 操作意图，转发到 Agent 服务")
		handleAgentRequest(c, &req)
		return
	}

	// 否则使用通用AI对话
	log.Log().Info("[意图识别] 通用对话")
	handleGeneralChat(c, &req)
}

// isHistoryQueryIntent 判断是否是查询历史记录的意图
func isHistoryQueryIntent(req model.AiChatRequest) bool {
	if len(req.Messages) == 0 {
		return false
	}

	lastMessage := req.Messages[len(req.Messages)-1]
	if lastMessage.Role != "user" {
		return false
	}

	message := lastMessage.Content
	messageLower := strings.ToLower(message)

	// 文件/目录 follow-up 属于继续执行或继续分析，不应归类为历史问答
	if isFileOrDirectoryFollowUp(req) {
		return false
	}

	// 1. 如果包含明确的IP地址或主机标识，一定不是历史查询
	if findIPAddress(message) != "" || strings.Contains(message, "ip=") || strings.Contains(message, "name=") {
		return false
	}

	// 2. 必须包含指代词（"刚才"、"那个"等），表示在引用历史
	hasReference := containsReferenceWords(message)
	if !hasReference {
		return false
	}

	// 3. 如果包含明确的查询动词，说明需要执行新操作（而非查询历史）
	hasQueryAction := strings.Contains(messageLower, "查询") ||
		strings.Contains(messageLower, "查看") ||
		strings.Contains(messageLower, "检查") ||
		strings.Contains(messageLower, "分析") ||
		strings.Contains(messageLower, "巡检") ||
		strings.Contains(messageLower, "扫描") ||
		strings.Contains(messageLower, "获取") ||
		strings.Contains(messageLower, "执行") ||
		strings.Contains(messageLower, "运行") ||
		strings.Contains(messageLower, "列出") ||
		strings.Contains(messageLower, "显示") ||
		strings.Contains(messageLower, "帮忙") ||
		strings.Contains(messageLower, "帮我")

	// 如果有查询动词，则不是历史查询，而是需要执行新操作
	if hasQueryAction {
		return false
	}

	// 4. 纯粹的问句（没有动词）才是历史查询
	// 例如："刚才检查的服务器叫什么名字" → 历史查询
	// 例如："那个服务器的IP是多少" → 历史查询
	hasQuestion := strings.Contains(messageLower, "什么") ||
		strings.Contains(messageLower, "哪") ||
		strings.Contains(messageLower, "怎么") ||
		strings.Contains(messageLower, "如何") ||
		strings.Contains(messageLower, "多少") ||
		strings.Contains(messageLower, "几") ||
		strings.Contains(messageLower, "是") ||
		strings.Contains(messageLower, "吗") ||
		strings.Contains(messageLower, "?") ||
		strings.Contains(messageLower, "？")

	return hasQuestion
}

// handleGeneralChat 处理通用对话（提取历史摘要增强上下文）
func handleGeneralChat(c *gin.Context, req *model.AiChatRequest) {
	if canned := buildGeneralAssistantIntro(req); canned != nil {
		result.Success(c, canned)
		return
	}

	modelService := service.NewAiModelService()

	// 【增强上下文理解】：如果消息中包含"刚才"、"之前"等指代词，且有历史记录
	// 则在消息前添加历史摘要，帮助 AI 理解上下文
	if req.SessionID != "" && len(req.Messages) > 1 {
		lastMsg := req.Messages[len(req.Messages)-1]
		if containsReferenceWords(lastMsg.Content) {
			// 提取历史中的关键信息（IP、主机名、操作类型等）
			summary := extractHistorySummary(req.Messages[:len(req.Messages)-1])
			if summary != "" {
				// 在当前消息前添加一个 system 消息，提供上下文
				systemMsg := model.AiChatMessage{
					Role:    "system",
					Content: fmt.Sprintf("上下文信息（用于回答用户问题）：%s", summary),
				}
				req.Messages = append([]model.AiChatMessage{systemMsg}, req.Messages...)
				log.Log().Infof("[上下文增强] 添加历史摘要: %s", summary)
			}
		}
	}

	modelService.Chat(c, req)
}

func buildGeneralAssistantIntro(req *model.AiChatRequest) *model.AiChatResponse {
	if len(req.Messages) == 0 {
		return nil
	}

	lastMsg := req.Messages[len(req.Messages)-1]
	if lastMsg.Role != "user" {
		return nil
	}

	content := strings.TrimSpace(lastMsg.Content)
	contentLower := strings.ToLower(content)
	normalized := strings.NewReplacer("，", "", "。", "", "！", "", "？", "", " ", "").Replace(contentLower)

	isGreeting := normalized == "你好" || normalized == "您好" || normalized == "hi" || normalized == "hello"
	isIdentityQuery := strings.Contains(content, "你是谁") ||
		strings.Contains(content, "你是干嘛的") ||
		strings.Contains(content, "你能做什么") ||
		strings.Contains(content, "你会什么") ||
		strings.Contains(content, "你可以做什么")

	if !isGreeting && !isIdentityQuery {
		return nil
	}

	text := `你好，我是你的 **AutoOps 运维智能小助手**。

目前我可以帮你做这些运维工作：

1. **服务器分析**
- CPU 分析
- 内存分析
- 磁盘分析
- 系统压力分析
- 目录深度分析

2. **远程操作**
- 在指定主机上执行命令
- 批量执行脚本或命令

3. **服务管理**
- 检查服务状态
- 安装服务
- 重启服务

你可以直接这样问我：
- 帮我检查 192.168.1.100 的 CPU 使用情况
- 帮我在 139.9.205.38 上执行命令 df -h
- 帮我检查 nginx 服务状态
- 帮我在 43.251.227.249 上安装 nginx`

	return &model.AiChatResponse{
		Text:         text,
		Content:      service.ParseMarkdown(text),
		FinishReason: "stop",
		Usage:        model.AliyunUsage{},
		SessionID:    req.SessionID,
	}
}

// containsReferenceWords 检查消息是否包含指代词
func containsReferenceWords(msg string) bool {
	referenceWords := []string{"刚才", "之前", "上面", "刚刚", "那个", "那台", "这个", "这台", "前面"}
	msgLower := strings.ToLower(msg)
	for _, word := range referenceWords {
		if strings.Contains(msgLower, word) {
			return true
		}
	}
	return false
}

func isFileOrDirectoryFollowUp(req model.AiChatRequest) bool {
	if len(req.Messages) == 0 {
		return false
	}

	lastMessage := req.Messages[len(req.Messages)-1]
	if lastMessage.Role != "user" {
		return false
	}

	message := lastMessage.Content
	messageLower := strings.ToLower(message)

	hasHistoryTarget := false
	if len(req.Messages) > 1 {
		if _, _, ok := findLastTargetInMessages(req.Messages[:len(req.Messages)-1]); ok {
			hasHistoryTarget = true
		}
	}

	hasPath := extractDirectoryPath(message) != "" || hasFilePath(message)
	hasReference := containsReferenceWords(message)
	hasFileOrDirNoun := strings.Contains(messageLower, "目录") ||
		strings.Contains(messageLower, "文件") ||
		strings.Contains(messageLower, "日志") ||
		strings.Contains(messageLower, "日志文件") ||
		strings.Contains(messageLower, "路径")
	hasFollowUpAction := strings.Contains(messageLower, "清空") ||
		strings.Contains(messageLower, "清除") ||
		strings.Contains(messageLower, "删除") ||
		strings.Contains(messageLower, "分析") ||
		strings.Contains(messageLower, "查看") ||
		strings.Contains(messageLower, "检查") ||
		strings.Contains(messageLower, "多大") ||
		strings.Contains(messageLower, "大小") ||
		strings.Contains(messageLower, "还能删") ||
		strings.Contains(messageLower, "哪些可以删") ||
		strings.Contains(messageLower, "是否清空") ||
		strings.Contains(messageLower, "占用")

	// 只要当前消息本身已经给出明确绝对路径，并且语义上是在看文件/目录内容，
	// 就应该允许进入 Agent，再由后续结构化上下文补齐目标主机。
	if hasPath && hasFileOrDirNoun && hasFollowUpAction {
		return true
	}

	return hasHistoryTarget && (hasPath || hasReference || hasFileOrDirNoun) && hasFollowUpAction
}

func shouldBackfillContextFromRecentHistory(message string) bool {
	messageLower := strings.ToLower(message)
	return extractDirectoryPath(message) != "" ||
		hasFilePath(message) ||
		strings.Contains(messageLower, "目录") ||
		strings.Contains(messageLower, "文件") ||
		strings.Contains(messageLower, "日志") ||
		strings.Contains(messageLower, "清空") ||
		strings.Contains(messageLower, "清除") ||
		strings.Contains(messageLower, "删除") ||
		strings.Contains(messageLower, "分析") ||
		strings.Contains(messageLower, "检查") ||
		strings.Contains(messageLower, "查看")
}

func loadRecentStructuredContextFromHistory(userID uint) (service.SessionStructuredContext, bool) {
	histories, _, err := agentDao.NewChatHistoryDao().GetByUserID(userID, 1, 20)
	if err != nil {
		log.Log().Warnf("[会话历史] 从最近历史回捞上下文失败: %v", err)
		return service.SessionStructuredContext{}, false
	}

	for _, history := range histories {
		if strings.TrimSpace(history.Entities) == "" {
			continue
		}

		entities := map[string]string{}
		if err := json.Unmarshal([]byte(history.Entities), &entities); err != nil {
			continue
		}

		targetName := strings.TrimSpace(entities["targetName"])
		if targetName == "" {
			continue
		}

		return service.SessionStructuredContext{
			LastTargetType:    strings.TrimSpace(entities["targetType"]),
			LastTargetName:    targetName,
			LastCheckType:     strings.TrimSpace(entities["checkType"]),
			LastQueryType:     strings.TrimSpace(entities["queryType"]),
			LastDirectoryPath: strings.TrimSpace(entities["directoryPath"]),
		}, true
	}

	return service.SessionStructuredContext{}, false
}

func mergeStructuredContext(current, fallback service.SessionStructuredContext) service.SessionStructuredContext {
	return service.SessionStructuredContext{
		LastTargetType:    firstNonEmpty(current.LastTargetType, fallback.LastTargetType),
		LastTargetName:    firstNonEmpty(current.LastTargetName, fallback.LastTargetName),
		LastCheckType:     firstNonEmpty(current.LastCheckType, fallback.LastCheckType),
		LastQueryType:     firstNonEmpty(current.LastQueryType, fallback.LastQueryType),
		LastDirectoryPath: firstNonEmpty(current.LastDirectoryPath, fallback.LastDirectoryPath),
	}
}

func hasFilePath(text string) bool {
	for _, field := range strings.Fields(text) {
		field = strings.TrimSpace(field)
		field = strings.Trim(field, "`*\"'，。,:：;；()[]{}<>")
		if strings.HasPrefix(field, "/") && strings.Contains(field, ".") && isLikelyLinuxPath(field) {
			return true
		}
	}
	return false
}

func isManagedFileOperationFollowUp(message string) bool {
	messageLower := strings.ToLower(message)
	if !(strings.Contains(messageLower, "清空") || strings.Contains(messageLower, "清除") || strings.Contains(messageLower, "删除")) {
		return false
	}

	for _, field := range strings.Fields(message) {
		field = strings.TrimSpace(field)
		field = strings.Trim(field, "`*\"'，。,:：;；()[]{}<>")
		if strings.HasPrefix(field, "/") && isLikelyLinuxPath(field) && isAllowedManagedFileSuffix(field) {
			return true
		}
	}

	return false
}

func isAllowedManagedFileSuffix(path string) bool {
	lower := strings.ToLower(strings.TrimSpace(path))
	return strings.HasSuffix(lower, ".log") ||
		strings.HasSuffix(lower, ".txt") ||
		strings.HasSuffix(lower, ".gz") ||
		strings.HasSuffix(lower, ".zip")
}

// extractHistorySummary 从历史对话中提取关键信息摘要
func extractHistorySummary(history []model.AiChatMessage) string {
	var summary strings.Builder

	// 提取最近的操作和目标
	var lastUserAction string
	var targets []string
	var hostNames []string
	var services []string

	log.Log().Infof("[历史摘要] 开始提取，历史消息数: %d", len(history))

	// 倒序遍历历史，提取最近的操作信息
	for i := len(history) - 1; i >= 0 && len(targets) < 3; i-- {
		msg := history[i]

		// 提取用户的操作意图
		if msg.Role == "user" {
			if lastUserAction == "" {
				if strings.Contains(msg.Content, "检查") || strings.Contains(msg.Content, "巡检") {
					lastUserAction = "检查服务器"
				} else if strings.Contains(msg.Content, "分析") {
					lastUserAction = "分析服务器"
				} else if strings.Contains(msg.Content, "查询") {
					lastUserAction = "查询信息"
				}
			}

			// 提取 IP 地址
			if ip := findIPAddress(msg.Content); ip != "" && !contains(targets, "IP:"+ip) {
				targets = append(targets, "IP:"+ip)
				log.Log().Infof("[历史摘要] 从用户消息提取到IP: %s", ip)
			}
		}

		// 从 Assistant 的回复中提取主机名、IP 和服务信息
		if msg.Role == "assistant" {
			lines := strings.Split(msg.Content, "\n")
			inProcessTable := false

			for _, line := range lines {
				line = strings.TrimSpace(line)

				// 提取主机名称（支持多种格式）
				if strings.Contains(line, "主机名称:") || strings.Contains(line, "主机名:") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) >= 2 {
						hostName := strings.TrimSpace(parts[1])
						hostName = strings.Trim(hostName, "`*_")
						if hostName != "" && !contains(hostNames, hostName) {
							hostNames = append(hostNames, hostName)
							log.Log().Infof("[历史摘要] 从助手回复提取到主机名: %s", hostName)
						}
					}
				}

				// 提取主机IP
				if strings.Contains(line, "主机IP:") || (strings.Contains(line, "IP:") && !strings.Contains(line, "主机名称")) {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) >= 2 {
						ipPart := strings.TrimSpace(parts[1])
						ipPart = strings.Trim(ipPart, "`*_")
						if ip := findIPAddress(ipPart); ip != "" && !contains(targets, "IP:"+ip) {
							targets = append(targets, "IP:"+ip)
							log.Log().Infof("[历史摘要] 从助手回复提取到IP: %s", ip)
						}
					}
				}

				// 提取进程表中的服务信息
				if strings.Contains(line, "TOP") && strings.Contains(line, "进程") {
					inProcessTable = true
					continue
				}

				// 从进程命令中提取服务名称
				if inProcessTable && strings.Contains(line, "|") {
					parts := strings.Split(line, "|")
					if len(parts) >= 5 {
						command := strings.TrimSpace(parts[4])
						if service := extractServiceFromCommand(command); service != "" && !contains(services, service) {
							services = append(services, service)
						}
					}
				}
			}
		}
	}

	// 构建摘要
	if lastUserAction != "" {
		summary.WriteString("最近操作: " + lastUserAction)
	}

	if len(hostNames) > 0 {
		if summary.Len() > 0 {
			summary.WriteString("; ")
		}
		summary.WriteString("涉及主机: " + strings.Join(hostNames, ", "))
	}

	if len(targets) > 0 {
		if summary.Len() > 0 {
			summary.WriteString("; ")
		}
		summary.WriteString("目标: " + strings.Join(targets, ", "))
	}

	if len(services) > 0 {
		if summary.Len() > 0 {
			summary.WriteString("; ")
		}
		summary.WriteString("检测到的服务: " + strings.Join(services, ", "))
	}

	result := summary.String()
	log.Log().Infof("[历史摘要] 提取完成: %s", result)
	return result
}

// extractServiceFromCommand 从进程命令中提取服务名称
func extractServiceFromCommand(command string) string {
	command = strings.ToLower(command)

	// 常见服务映射
	serviceMap := map[string]string{
		"mysqld":         "MySQL",
		"mysql":          "MySQL",
		"redis-server":   "Redis",
		"redis":          "Redis",
		"nginx":          "Nginx",
		"httpd":          "Apache",
		"postgres":       "PostgreSQL",
		"mongod":         "MongoDB",
		"docker":         "Docker",
		"prometheus":     "Prometheus",
		"grafana":        "Grafana",
		"containerd":     "Containerd",
		"systemd":        "", // 忽略系统进程
		"sshd":           "SSH",
		"node":           "Node.js",
		"java":           "Java应用",
		"python":         "Python应用",
		"dodevops-agent": "DoDevOps监控",
		"pushgateway":    "Pushgateway",
		"cloudmonitor":   "云监控",
	}

	for keyword, serviceName := range serviceMap {
		if strings.Contains(command, keyword) {
			return serviceName
		}
	}

	return ""
}

// contains 检查字符串切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isAgentIntent 判断是否是Agent意图（目标明确的运维查询）
func isAgentIntent(req model.AiChatRequest) bool {
	if len(req.Messages) == 0 {
		return false
	}

	// 获取最后一条用户消息
	lastMessage := req.Messages[len(req.Messages)-1]
	if lastMessage.Role != "user" {
		return false
	}

	message := lastMessage.Content
	messageLower := strings.ToLower(message)
	structuredCtx := service.GetSessionStructuredContext(req.SessionID)

	log.Log().Infof("[Agent意图判断] 原始消息: %s", message)

	if isManagedFileOperationFollowUp(message) {
		log.Log().Infof("[Agent意图判断] 检测到白名单文件操作请求，判定为Agent意图")
		return true
	}

	if isFileOrDirectoryFollowUp(req) {
		log.Log().Infof("[Agent意图判断] 检测到文件/目录 follow-up，判定为Agent意图")
		return true
	}

	if structuredCtx.LastTargetName != "" && isContinuationOperation(messageLower) {
		log.Log().Infof("[Agent意图判断] 检测到 continuation follow-up，复用会话目标 %s", structuredCtx.LastTargetName)
		return true
	}

	if structuredCtx.LastTargetName != "" && isServiceOperationFollowUp(messageLower) {
		log.Log().Infof("[Agent意图判断] 检测到服务操作 follow-up，复用会话目标 %s", structuredCtx.LastTargetName)
		return true
	}

	// -2. 最优先排除：历史查询意图（包含"刚才"、"之前"等指代词的问句）
	// 这类问题应该由通用AI基于历史上下文回答，而不是触发新的Agent操作
	if isHistoryQueryIntent(req) {
		log.Log().Infof("[Agent意图判断] 检测到历史查询意图（包含指代词），判定为非Agent意图")
		return false
	}

	// -1. 优先排除：通用技术问题（写脚本、教学类问题）
	if isGeneralTechnicalQuestion(messageLower) {
		log.Log().Infof("[Agent意图判断] 检测到通用技术问题，判定为非Agent意图")
		return false
	}

	// 0. 优先检查：如果包含IP地址或主机名 + 监控指标关键词，直接判定为Agent意图
	hasTarget := strings.Contains(message, "ip=") ||
		strings.Contains(message, "name=") ||
		findIPAddress(message) != "" ||
		hasHostNamePattern(message)

	hasMetrics := strings.Contains(messageLower, "内存") ||
		strings.Contains(messageLower, "cpu") ||
		strings.Contains(messageLower, "磁盘") ||
		strings.Contains(messageLower, "disk") ||
		strings.Contains(messageLower, "memory") ||
		strings.Contains(messageLower, "压力") ||
		strings.Contains(messageLower, "负载") ||
		strings.Contains(messageLower, "性能") ||
		strings.Contains(messageLower, "进程") ||
		strings.Contains(messageLower, "端口") ||
		strings.Contains(messageLower, "服务")

	// 如果同时包含目标和指标关键词，不管有没有动作词，都判定为Agent意图
	if hasTarget && hasMetrics {
		log.Log().Infof("[Agent意图判断] 检测到目标+指标关键词，判定为Agent意图")
		return true
	}

	// 1. 必须包含运维相关的动作关键词（扩展同义词）
	hasAction := strings.Contains(messageLower, "检查") ||
		strings.Contains(messageLower, "查询") ||
		strings.Contains(messageLower, "查看") ||
		strings.Contains(messageLower, "显示") ||
		strings.Contains(messageLower, "分析") ||
		strings.Contains(messageLower, "重启") ||
		strings.Contains(messageLower, "安装") ||
		strings.Contains(messageLower, "部署") ||
		strings.Contains(messageLower, "巡检") ||
		strings.Contains(messageLower, "扫描") ||
		strings.Contains(messageLower, "列出") ||
		strings.Contains(messageLower, "获取") ||
		strings.Contains(messageLower, "有哪些") ||
		strings.Contains(messageLower, "执行") ||
		strings.Contains(messageLower, "运行") ||
		strings.Contains(messageLower, "批量") ||
		strings.Contains(messageLower, "清空") ||
		strings.Contains(messageLower, "清除") ||
		strings.Contains(messageLower, "删除")

	if !hasAction {
		log.Log().Infof("[Agent意图判断] 未检测到动作关键词，判定为非Agent意图")
		return false
	}

	// 2. 有明确目标，或历史中存在目标（合并历史后可用）
	if hasTarget {
		log.Log().Infof("[Agent意图判断] 检测到明确目标，判定为Agent意图")
		return true
	}

	// 从历史消息中回溯寻找最近的目标
	if _, _, ok := findLastTargetInMessages(req.Messages); ok {
		log.Log().Infof("[Agent意图判断] 从历史中找到目标，判定为Agent意图")
		return true
	}

	log.Log().Infof("[Agent意图判断] 未找到目标，判定为非Agent意图")
	return false
}

func isContinuationOperation(messageLower string) bool {
	hasContinuation := strings.Contains(messageLower, "继续") ||
		strings.Contains(messageLower, "接着") ||
		strings.Contains(messageLower, "然后") ||
		strings.Contains(messageLower, "再") ||
		strings.Contains(messageLower, "还是这台") ||
		strings.Contains(messageLower, "这台机器继续") ||
		strings.Contains(messageLower, "这个主机继续")
	if !hasContinuation {
		return false
	}

	hasOperation := strings.Contains(messageLower, "检查") ||
		strings.Contains(messageLower, "查看") ||
		strings.Contains(messageLower, "分析") ||
		strings.Contains(messageLower, "重启") ||
		strings.Contains(messageLower, "安装") ||
		strings.Contains(messageLower, "部署") ||
		strings.Contains(messageLower, "执行命令") ||
		strings.Contains(messageLower, "运行命令") ||
		strings.Contains(messageLower, "执行脚本") ||
		strings.Contains(messageLower, "运行脚本") ||
		strings.Contains(messageLower, "列出") ||
		strings.Contains(messageLower, "获取") ||
		strings.Contains(messageLower, "查询") ||
		strings.Contains(messageLower, "清空") ||
		strings.Contains(messageLower, "清除") ||
		strings.Contains(messageLower, "删除")

	return hasOperation
}

func isServiceOperationFollowUp(messageLower string) bool {
	if strings.TrimSpace(messageLower) == "" {
		return false
	}

	hasAction := strings.Contains(messageLower, "安装") ||
		strings.Contains(messageLower, "部署") ||
		strings.Contains(messageLower, "重启") ||
		strings.Contains(messageLower, "restart") ||
		strings.Contains(messageLower, "检查") ||
		strings.Contains(messageLower, "查看") ||
		strings.Contains(messageLower, "状态") ||
		strings.Contains(messageLower, "是否启动") ||
		strings.Contains(messageLower, "运行情况")
	if !hasAction {
		return false
	}

	return extractServiceNameForIntent(messageLower) != "" ||
		strings.Contains(messageLower, "这个服务") ||
		strings.Contains(messageLower, "该服务")
}

// isGeneralTechnicalQuestion 判断是否是通用技术问题（非运维操作）
func isGeneralTechnicalQuestion(messageLower string) bool {
	// 1. 写脚本/写代码类问题
	if strings.Contains(messageLower, "写") && (strings.Contains(messageLower, "脚本") ||
		strings.Contains(messageLower, "代码") ||
		strings.Contains(messageLower, "程序") ||
		strings.Contains(messageLower, "示例")) {
		return true
	}

	// 2. 教学/咨询类问题
	if strings.Contains(messageLower, "如何") && (strings.Contains(messageLower, "编写") ||
		strings.Contains(messageLower, "实现") ||
		strings.Contains(messageLower, "配置") ||
		strings.Contains(messageLower, "安装") ||
		strings.Contains(messageLower, "部署") ||
		strings.Contains(messageLower, "使用")) {
		// 但如果同时包含具体IP或主机名，则可能是针对特定服务器的操作
		if findIPAddress(messageLower) == "" && !strings.Contains(messageLower, "ip=") && !strings.Contains(messageLower, "name=") {
			return true
		}
	}

	// 3. 学习/了解类问题
	if strings.Contains(messageLower, "怎么") && (strings.Contains(messageLower, "做") ||
		strings.Contains(messageLower, "实现") ||
		strings.Contains(messageLower, "编写")) {
		// 排除特定服务器的操作
		if findIPAddress(messageLower) == "" && !strings.Contains(messageLower, "ip=") && !strings.Contains(messageLower, "name=") {
			return true
		}
	}

	// 4. 帮我生成/创建类问题（不涉及具体服务器）
	if (strings.Contains(messageLower, "帮我") || strings.Contains(messageLower, "帮忙")) &&
		(strings.Contains(messageLower, "生成") ||
			strings.Contains(messageLower, "创建") ||
			strings.Contains(messageLower, "写一个") ||
			strings.Contains(messageLower, "给我")) &&
		findIPAddress(messageLower) == "" {
		return true
	}

	return false
}

// findIPAddress 查找IP地址
func findIPAddress(text string) string {
	parts := strings.Split(text, " ")
	for _, part := range parts {
		if strings.Count(part, ".") == 3 {
			segments := strings.Split(part, ".")
			if len(segments) == 4 {
				isIP := true
				for _, seg := range segments {
					seg = strings.Trim(seg, ",.!?;:")
					if len(seg) == 0 || len(seg) > 3 {
						isIP = false
						break
					}
					for _, c := range seg {
						if c < '0' || c > '9' {
							isIP = false
							break
						}
					}
					if !isIP {
						break
					}
				}
				if isIP {
					return strings.Trim(part, ",.!?;:")
				}
			}
		}
	}
	return ""
}

// hasHostNamePattern 检查是否有主机名模式
func hasHostNamePattern(message string) bool {
	words := strings.Fields(message)
	for _, word := range words {
		word = strings.Trim(word, ",.!?;:()[]{}<>*\"'`")
		if strings.HasPrefix(word, "/") || looksLikeSizeTokenForIntent(word) {
			continue
		}
		// 检查是否像主机名（包含-或数字，且不是动词）
		if (strings.Contains(word, "-") || containsDigit(word)) &&
			!strings.Contains(word, "检查") &&
			!strings.Contains(word, "分析") &&
			!strings.Contains(word, "查询") &&
			len(word) > 2 {
			return true
		}
	}
	return false
}

func looksLikeSizeTokenForIntent(value string) bool {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "()[]{}")
	if value == "" {
		return false
	}

	lower := strings.ToLower(value)
	units := []string{"kb", "mb", "gb", "tb", "pb", "k", "m", "g", "t", "p"}
	for _, unit := range units {
		if strings.HasSuffix(lower, unit) {
			numberPart := strings.TrimSpace(lower[:len(lower)-len(unit)])
			if numberPart == "" {
				return false
			}
			dotCount := 0
			for _, r := range numberPart {
				if r >= '0' && r <= '9' {
					continue
				}
				if r == '.' {
					dotCount++
					if dotCount <= 1 {
						continue
					}
				}
				return false
			}
			return true
		}
	}

	return false
}

// containsDigit 检查是否包含数字
func containsDigit(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}

// handleAgentRequest 处理Agent请求
func handleAgentRequest(c *gin.Context, req *model.AiChatRequest) {
	// 获取用户ID
	userID, err := jwt.GetAdminId(c)
	if err != nil {
		result.Failed(c, int(result.ApiCode.NOAUTH), "未授权")
		return
	}

	// 提取用户消息
	if len(req.Messages) == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "消息不能为空")
		return
	}

	lastMessage := req.Messages[len(req.Messages)-1]
	hasExplicitTarget := strings.Contains(lastMessage.Content, "ip=") || strings.Contains(lastMessage.Content, "name=") || len(findIPAddress(lastMessage.Content)) > 0
	resolvedTargetType := ""
	resolvedTargetName := ""
	resolvedCheckType := ""
	resolvedQueryType := ""
	resolvedDirectoryPath := extractDirectoryPath(lastMessage.Content)
	preferResolvedTarget := false

	structuredCtx := service.GetSessionStructuredContext(req.SessionID)
	if structuredCtx.LastTargetName == "" && shouldBackfillContextFromRecentHistory(lastMessage.Content) {
		if recentCtx, ok := loadRecentStructuredContextFromHistory(userID); ok {
			structuredCtx = mergeStructuredContext(structuredCtx, recentCtx)
		}
	}
	isFollowUp := shouldTreatAsFollowUp(lastMessage.Content, structuredCtx)
	if structuredCtx.LastTargetName != "" && !hasExplicitTarget && isFollowUp {
		resolvedTargetType = structuredCtx.LastTargetType
		resolvedTargetName = structuredCtx.LastTargetName
		resolvedCheckType = structuredCtx.LastCheckType
		resolvedQueryType = structuredCtx.LastQueryType
		preferResolvedTarget = true
	}
	if resolvedDirectoryPath == "" && shouldReuseLastDirectoryPath(lastMessage.Content) {
		resolvedDirectoryPath = structuredCtx.LastDirectoryPath
	}
	if resolvedTargetName == "" {
		if kind, target, ok := findLastTargetInMessages(req.Messages); ok {
			resolvedTargetType = kind
			resolvedTargetName = target
		}
	}

	// 将原始用户消息写入内存会话，避免污染原始文本
	if req.SessionID != "" {
		service.AppendMemory(req.SessionID, model.AiChatMessage{Role: "user", Content: lastMessage.Content})
	}

	// 构建Agent请求，将结构化上下文单独传递给Agent
	agentReq := &agentModel.AgentChatRequest{
		Message:               lastMessage.Content,
		SessionID:             req.SessionID,
		ResolvedTargetType:    resolvedTargetType,
		ResolvedTargetName:    resolvedTargetName,
		ResolvedCheckType:     resolvedCheckType,
		ResolvedQueryType:     resolvedQueryType,
		ResolvedDirectoryPath: resolvedDirectoryPath,
		PreferResolvedTarget:  preferResolvedTarget,
	}

	// 调用Agent服务
	agentSvc := agentService.NewAgentChatService()
	agentResponse, err := agentSvc.ProcessMessage(c.Request.Context(), userID, agentReq)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 把助手回复写入内存会话
	if req.SessionID != "" && agentResponse != nil && agentResponse.Message != "" {
		service.AppendMemory(req.SessionID, model.AiChatMessage{Role: "assistant", Content: agentResponse.Message})
		if agentResponse.Intent != nil && agentResponse.Intent.Entities != nil {
			entities := agentResponse.Intent.Entities
			targetType := entities["targetType"]
			targetName := entities["targetName"]
			directoryPath := firstNonEmpty(entities["directoryPath"], resolvedDirectoryPath)
			if directoryPath == "" && agentResponse.TaskInfo != nil {
				if taskDirectoryPath, ok := agentResponse.TaskInfo["directoryPath"].(string); ok {
					directoryPath = strings.TrimSpace(taskDirectoryPath)
				} else if topDirectoryPath, ok := agentResponse.TaskInfo["topDirectoryPath"].(string); ok {
					directoryPath = strings.TrimSpace(topDirectoryPath)
				}
			}
			if targetName == "" {
				targetType = resolvedTargetType
				targetName = resolvedTargetName
			}
			if targetName != "" {
				service.SaveSessionStructuredContext(req.SessionID, service.SessionStructuredContext{
					LastTargetType:    targetType,
					LastTargetName:    targetName,
					LastCheckType:     firstNonEmpty(entities["checkType"], resolvedCheckType),
					LastQueryType:     firstNonEmpty(entities["queryType"], resolvedQueryType),
					LastDirectoryPath: directoryPath,
				})
			}
		}
		log.Log().Infof("[会话历史] Agent 响应已保存到内存，会话ID: %s, 响应长度: %d 字符", req.SessionID, len(agentResponse.Message))
	}

	// 转换为统一的AI响应格式
	aiResponse := convertAgentResponseToAiResponse(agentResponse)

	// 返回结果
	result.Success(c, aiResponse)
}

// convertAgentResponseToAiResponse 将Agent响应转换为AI响应格式
func convertAgentResponseToAiResponse(agentResp *agentModel.AgentChatResponse) *model.AiChatResponse {
	// 使用ParseMarkdown解析消息为结构化内容
	contentBlocks := service.ParseMarkdown(agentResp.Message)

	// 构建统一的AI响应格式
	return &model.AiChatResponse{
		Text:         agentResp.Message,
		Content:      contentBlocks,
		FinishReason: "stop",
		Usage: model.AliyunUsage{
			OutputTokens: 0, // Agent服务不统计token，设为0
			InputTokens:  0,
		},
		RequestID: agentResp.SessionID,
		SessionID: agentResp.SessionID,
	}
}

// loadSessionHistory 从Redis加载历史messages
func loadSessionHistory(sessionID string) []model.AiChatMessage {
	ctx := context.Background()
	key := "ai:session:" + sessionID
	val, err := redisPkg.RedisDb.Get(ctx, key).Result()
	if err != nil || val == "" {
		return nil
	}
	var msgs []model.AiChatMessage
	if err := json.Unmarshal([]byte(val), &msgs); err != nil {
		return nil
	}
	return msgs
}

func loadSessionHistoryFromDB(sessionID string) []model.AiChatMessage {
	histories, err := agentDao.NewChatHistoryDao().GetBySessionID(sessionID, 20)
	if err != nil || len(histories) == 0 {
		return nil
	}

	return convertAgentHistoriesToAiMessages(histories)
}

func convertAgentHistoriesToAiMessages(histories []agentModel.AgentChatHistory) []model.AiChatMessage {
	msgs := make([]model.AiChatMessage, 0, len(histories))
	for _, history := range histories {
		role := strings.TrimSpace(history.Role)
		content := strings.TrimSpace(history.Message)
		if role == "" || content == "" {
			continue
		}
		msgs = append(msgs, model.AiChatMessage{
			Role:    role,
			Content: content,
		})
	}
	return msgs
}

// findLastTargetInMessages 回溯历史，寻找最近的目标（IP或主机名）
func findLastTargetInMessages(msgs []model.AiChatMessage) (kind string, target string, ok bool) {
	for i := len(msgs) - 1; i >= 0; i-- {
		m := msgs[i]
		if m.Role != "user" {
			continue
		}
		if strings.Contains(m.Content, "ip=") {
			parts := strings.Split(m.Content, "ip=")
			if len(parts) > 1 {
				rest := parts[1]
				t := strings.FieldsFunc(rest, func(r rune) bool { return r == ' ' || r == ',' || r == '，' || r == '。' })
				if len(t) > 0 {
					return "host", t[0], true
				}
			}
		}
		if strings.Contains(m.Content, "name=") {
			parts := strings.Split(m.Content, "name=")
			if len(parts) > 1 {
				rest := parts[1]
				t := strings.FieldsFunc(rest, func(r rune) bool { return r == ' ' || r == ',' || r == '，' || r == '。' })
				if len(t) > 0 {
					return "host", t[0], true
				}
			}
		}
		if ip := findIPAddress(m.Content); ip != "" {
			return "host", ip, true
		}
	}
	return "", "", false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func extractDirectoryPath(text string) string {
	for i := 0; i < len(text); i++ {
		if text[i] != '/' {
			continue
		}
		j := i
		for j < len(text) {
			c := text[j]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
				c == '/' || c == '-' || c == '_' || c == '.' || c == '@' {
				j++
				continue
			}
			break
		}
		if j > i {
			candidate := strings.TrimSpace(text[i:j])
			candidate = strings.Trim(candidate, "`*\"'，。,:：;；()[]{}<>")
			if strings.HasPrefix(candidate, "/") && isLikelyLinuxPath(candidate) {
				return candidate
			}
		}
	}
	return ""
}

func shouldReuseLastDirectoryPath(text string) bool {
	lower := strings.ToLower(text)
	if extractDirectoryPath(text) != "" {
		return false
	}
	return strings.Contains(lower, "这个目录") ||
		strings.Contains(lower, "该目录") ||
		strings.Contains(lower, "这个路径") ||
		strings.Contains(lower, "继续分析") ||
		(strings.Contains(lower, "里面") && (strings.Contains(lower, "删除") || strings.Contains(lower, "文件") || strings.Contains(lower, "清理")))
}

func shouldTreatAsFollowUp(text string, structuredCtx service.SessionStructuredContext) bool {
	if strings.TrimSpace(text) == "" {
		return false
	}
	if structuredCtx.LastTargetName == "" {
		return false
	}
	lower := strings.ToLower(text)
	if strings.Contains(text, "ip=") || strings.Contains(text, "name=") || findIPAddress(text) != "" {
		return false
	}
	if extractDirectoryPath(text) != "" && structuredCtx.LastTargetName != "" {
		return true
	}
	if isServiceOperationFollowUp(lower) {
		return true
	}
	followUpKeywords := []string{
		"这个目录", "该目录", "这个路径", "继续分析", "继续看", "接着分析", "上面的目录",
		"里面哪些", "哪些可以删除", "哪些文件", "哪些能清", "占用的空间比较大",
	}
	for _, keyword := range followUpKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

func extractServiceNameForIntent(messageLower string) string {
	serviceAliases := []struct {
		canonical string
		aliases   []string
	}{
		{canonical: "nginx", aliases: []string{"nginx", "ningx"}},
		{canonical: "apache2", aliases: []string{"apache2", "apache", "httpd"}},
		{canonical: "redis", aliases: []string{"redis", "redis-server"}},
		{canonical: "mysql", aliases: []string{"mysql", "mysqld"}},
		{canonical: "mariadb", aliases: []string{"mariadb"}},
		{canonical: "postgresql", aliases: []string{"postgresql", "postgres"}},
		{canonical: "mongodb", aliases: []string{"mongodb", "mongod"}},
		{canonical: "docker", aliases: []string{"docker"}},
		{canonical: "prometheus", aliases: []string{"prometheus"}},
		{canonical: "grafana", aliases: []string{"grafana"}},
	}
	for _, candidate := range serviceAliases {
		for _, alias := range candidate.aliases {
			if strings.Contains(messageLower, alias) {
				return candidate.canonical
			}
		}
	}
	return ""
}

func isLikelyLinuxPath(path string) bool {
	if path == "" || !strings.HasPrefix(path, "/") {
		return false
	}
	for _, r := range path {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '/', r == '-', r == '_', r == '.', r == '@':
		default:
			return false
		}
	}
	return true
}
