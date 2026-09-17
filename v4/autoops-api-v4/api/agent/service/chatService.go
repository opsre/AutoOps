package service

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"dodevops-api/api/agent/dao"
	"dodevops-api/api/agent/model"
	"dodevops-api/api/agent/skill"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	"dodevops-api/common/util"
	baseSkill "dodevops-api/pkg/agent/skill"
	"dodevops-api/pkg/log"

	"github.com/google/uuid"
)

// AgentChatService AI对话服务接口
type AgentChatService interface {
	ProcessMessage(ctx context.Context, userID uint, req *model.AgentChatRequest) (*model.AgentChatResponse, error)
	GetChatHistory(ctx context.Context, userID uint, sessionID string, limit int) ([]model.AgentChatHistory, error)
	ClearChatHistory(ctx context.Context, sessionID string) error
}

// AgentChatServiceImpl AI对话服务实现
type AgentChatServiceImpl struct {
	chatDao             dao.ChatHistoryDao
	hostDao             cmdbDao.CmdbHostDao
	inspectionSvc       InspectionService
	performanceAnalysis *PerformanceAnalysisService
	skillRegistry       *baseSkill.SkillRegistry
}

// NewAgentChatService 创建AI对话服务实例
func NewAgentChatService() AgentChatService {
	svc := &AgentChatServiceImpl{
		chatDao:             dao.NewChatHistoryDao(),
		hostDao:             cmdbDao.NewCmdbHostDao(),
		inspectionSvc:       NewInspectionService(),
		performanceAnalysis: NewPerformanceAnalysisService(),
		skillRegistry:       baseSkill.GetDefaultRegistry(),
	}

	// 注册所有技能
	svc.registerSkills()

	return svc
}

// registerSkills 注册所有技能
func (s *AgentChatServiceImpl) registerSkills() {
	log.Log().Info("🚀 [Skills 架构] 开始注册技能...")

	// 注册分析类技能
	s.skillRegistry.Register(skill.NewPressureAnalysisSkill())
	s.skillRegistry.Register(skill.NewMemoryAnalysisSkill())
	s.skillRegistry.Register(skill.NewDirectoryDeepAnalysisSkill())
	s.skillRegistry.Register(skill.NewDiskAnalysisSkill())
	s.skillRegistry.Register(skill.NewCPUAnalysisSkill())

	// 注册服务类统一技能
	s.skillRegistry.Register(skill.NewServiceManageSkill())

	// 注册操作类技能
	s.skillRegistry.Register(skill.NewCommandExecutionSkill())
	s.skillRegistry.Register(skill.NewScriptExecutionSkill())

	allSkills := s.skillRegistry.GetAll()
	log.Log().Infof("✅ [Skills 架构] 注册完成，共 %d 个技能", len(allSkills))
	for _, sk := range allSkills {
		log.Log().Infof("   - %s: %s (Category: %s)", sk.Name(), sk.Description(), sk.Category())
	}
}

// ProcessMessage 处理用户消息
func (s *AgentChatServiceImpl) ProcessMessage(ctx context.Context, userID uint, req *model.AgentChatRequest) (*model.AgentChatResponse, error) {
	// 1. 生成或使用现有会话ID
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	// 2. 意图识别（先识别意图，再决定是否记录）
	intent, err := s.recognizeIntent(ctx, req.Message)
	if err != nil {
		return s.createErrorResponse(sessionID, "意图识别失败: "+err.Error()), nil
	}
	if intent != nil {
		intent.Entities = mergeResolvedContextIntoIntent(intent.Entities, req)
	}

	// 3. 判断是否需要记录到数据库（只记录重要事件）
	shouldRecord := s.shouldRecordHistory(intent.Type)

	// 4. 如果需要记录，保存用户消息
	if shouldRecord {
		userHistory := &model.AgentChatHistory{
			SessionID:  sessionID,
			UserID:     userID,
			Role:       "user",
			Message:    req.Message,
			Intent:     intent.Type,
			IntentConf: intent.Confidence,
			Status:     1, // 处理中
			CreateTime: util.HTime{Time: time.Now()},
		}
		entitiesJSON, _ := json.Marshal(intent.Entities)
		userHistory.Entities = string(entitiesJSON)

		if err := s.chatDao.Create(userHistory); err != nil {
			log.Log().Errorf("保存用户消息失败: %v", err)
		}
	}

	// 5. 【Skills 路由】优先尝试使用 Skill 处理
	var response *model.AgentChatResponse

	skillIntent := &baseSkill.Intent{
		Type:       intent.Type,
		Action:     intent.Action,
		Confidence: intent.Confidence,
		Entities:   intent.Entities,
	}

	selectedSkill := s.skillRegistry.Find(skillIntent)
	if selectedSkill != nil {
		log.Log().Infof("✅ [Skills 架构] 找到匹配的 Skill: %s (Category: %s)", selectedSkill.Name(), selectedSkill.Category())

		// 解析目标列表（支持多主机）
		var targetList []string
		targetName := intent.Entities["targetName"]
		if targetName != "" {
			// 支持逗号分隔的多个目标（可能是英文逗号或中文逗号）
			targetName = strings.ReplaceAll(targetName, "，", ",")
			targetList = strings.Split(targetName, ",")
			// 去除空格
			for i := range targetList {
				targetList[i] = strings.TrimSpace(targetList[i])
			}
		}

		// 构建 Skill 请求
		skillReq := &baseSkill.SkillRequest{
			SessionID:  sessionID,
			UserID:     userID,
			Intent:     skillIntent,
			Message:    req.Message,
			TargetList: targetList, // 传递目标列表
		}

		// 验证参数
		if err := selectedSkill.Validate(skillReq); err != nil {
			response = s.createErrorResponse(sessionID, err.Error())
		} else {
			// 执行 Skill（获取目标信息）
			skillResp, err := selectedSkill.Execute(ctx, skillReq)
			if err != nil {
				log.Log().Warnf("❌ [Skills 架构] Skill 执行失败: %v", err)
				response = s.createErrorResponse(sessionID, err.Error())
			} else {
				// Skill 返回了目标信息，service 层执行实际分析
				response, err = s.executeSkillAnalysis(ctx, sessionID, selectedSkill.Name(), skillResp)
				if err != nil {
					log.Log().Warnf("⚠️ [Skills 架构] 分析失败，回退到旧逻辑: %v", err)
					response, err = s.processMessageLegacy(ctx, sessionID, userID, intent)
					if err != nil {
						response = s.createErrorResponse(sessionID, err.Error())
					}
				}
			}
		}
	} else {
		// 6. 【降级逻辑】使用原有的路由
		log.Log().Infof("⚠️ [旧逻辑] 未找到匹配的 Skill，使用旧逻辑处理")
		response, err = s.processMessageLegacy(ctx, sessionID, userID, intent)
		if err != nil {
			response = s.createErrorResponse(sessionID, err.Error())
		}
	}

	// 7. 保存 AI 回复（如果需要）
	if shouldRecord && response != nil {
		entitiesJSON, _ := json.Marshal(intent.Entities)
		assistantHistory := &model.AgentChatHistory{
			SessionID:  sessionID,
			UserID:     userID,
			Role:       "assistant",
			Message:    response.Message,
			Intent:     intent.Type,
			IntentConf: intent.Confidence,
			Entities:   string(entitiesJSON),
			Status:     2, // 成功
			CreateTime: util.HTime{Time: time.Now()},
		}
		s.chatDao.Create(assistantHistory)
	}
	if response != nil {
		response.Intent = intent
	}

	return response, nil
}

// processMessageLegacy 原有的消息处理逻辑（保留作为 fallback）
func (s *AgentChatServiceImpl) processMessageLegacy(ctx context.Context, sessionID string, userID uint, intent *model.IntentResult) (*model.AgentChatResponse, error) {
	var response *model.AgentChatResponse
	var err error

	// 根据意图路由到不同的处理器
	switch intent.Type {
	case "inspection":
		response, err = s.handleInspection(ctx, sessionID, userID, intent)
	case "query":
		response, err = s.handleQuery(ctx, sessionID, userID, intent)
	case "incident":
		response, err = s.handleIncident(ctx, sessionID, userID, intent)
	case "prediction":
		response = s.createTextResponse(sessionID, "容量预测功能正在开发中，敬请期待...")
	case "general":
		response, err = s.handleGeneral(ctx, sessionID, "")
	default:
		response = s.createTextResponse(sessionID, "抱歉，我还不太理解您的意思。您可以尝试：\n1. 检查服务器状态\n2. 执行系统巡检\n3. 查询主机信息")
	}

	return response, err
}

// executeSkillAnalysis 执行 Skill 的实际分析逻辑
func (s *AgentChatServiceImpl) executeSkillAnalysis(ctx context.Context, sessionID string, skillName string, skillResp *baseSkill.SkillResponse) (*model.AgentChatResponse, error) {
	log.Log().Infof("🔧 [Skills 架构] 执行 Skill 分析: %s", skillName)

	// 检查是否是批量处理
	if skillResp.Type == "batch" {
		return s.executeBatchAnalysis(ctx, sessionID, skillName, skillResp)
	}

	// 根据 Skill 名称执行对应的分析（单主机）
	switch skillName {
	case "pressure_analysis":
		return s.executePressureAnalysis(ctx, sessionID, skillResp)
	case "memory_analysis":
		return s.executeMemoryAnalysis(ctx, sessionID, skillResp)
	case "cpu_analysis":
		return s.executeCPUAnalysis(ctx, sessionID, skillResp)
	case "disk_analysis":
		return s.executeDiskAnalysis(ctx, sessionID, skillResp)
	case "directory_deep_analysis":
		return s.executeDirectoryDeepAnalysis(ctx, sessionID, skillResp)
	case "command_execution", "script_execution", "service_manage":
		// 命令执行和脚本执行已经在 Skill 中完成，直接返回结果
		response := &model.AgentChatResponse{
			SessionID: sessionID,
			Message:   skillResp.Message,
			Type:      skillResp.Type,
			TaskInfo:  skillResp.TaskInfo,
		}
		if skillName == "command_execution" {
			return s.enrichManagedFileOperationResponse(ctx, response)
		}
		return response, nil
	default:
		log.Log().Infof("ℹ️ [Skills 架构] %s 使用默认处理", skillName)
		// 默认直接返回 Skill 的响应
		return &model.AgentChatResponse{
			SessionID: sessionID,
			Message:   skillResp.Message,
			Type:      skillResp.Type,
			TaskInfo:  skillResp.TaskInfo,
		}, nil
	}
}

func (s *AgentChatServiceImpl) enrichManagedFileOperationResponse(ctx context.Context, response *model.AgentChatResponse) (*model.AgentChatResponse, error) {
	if response == nil || response.TaskInfo == nil {
		return response, nil
	}

	fileOp, _ := response.TaskInfo["fileOp"].(string)
	filePath, _ := response.TaskInfo["filePath"].(string)
	if fileOp == "" {
		return response, nil
	}

	directoryPath := managedFileFollowUpDirectory(filePath)
	if directoryPath == "" {
		return response, nil
	}

	hostIP, _ := response.TaskInfo["hostIP"].(string)
	if strings.TrimSpace(hostIP) == "" {
		return response, nil
	}

	switch fileOp {
	case "clear":
		message, err := s.buildManagedFileClearResponse(hostIP, filePath, directoryPath, response)
		if err != nil {
			response.Message = response.Message + fmt.Sprintf("\n\n---\n\n已完成文件清空，但后续信息采集失败: %v", err)
			return response, nil
		}
		response.Message = message
		response.TaskInfo["followUpTaskType"] = "managed_file_clear"
		response.TaskInfo["followUpDirectoryPath"] = directoryPath
	case "delete":
		message, err := s.buildManagedFileDeleteResponse(hostIP, filePath, directoryPath, response)
		if err != nil {
			response.Message = response.Message + fmt.Sprintf("\n\n---\n\n文件删除后续信息采集失败: %v", err)
			return response, nil
		}
		response.Message = message
		response.TaskInfo["followUpTaskType"] = "managed_file_delete"
		response.TaskInfo["followUpDirectoryPath"] = directoryPath
	}
	return response, nil
}

func managedFileFollowUpDirectory(filePath string) string {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" || !strings.HasPrefix(filePath, "/") {
		return ""
	}
	dir := path.Dir(filePath)
	if dir == "." || dir == "" {
		return ""
	}
	return dir
}

func mergeResolvedContextIntoIntent(entities map[string]string, req *model.AgentChatRequest) map[string]string {
	if entities == nil {
		entities = make(map[string]string)
	}

	if req == nil {
		return entities
	}

	if req.PreferResolvedTarget && req.ResolvedTargetName != "" {
		entities["targetName"] = req.ResolvedTargetName
		if req.ResolvedTargetType != "" {
			entities["targetType"] = req.ResolvedTargetType
		}
	}
	if entities["targetName"] == "" && req.ResolvedTargetName != "" {
		entities["targetName"] = req.ResolvedTargetName
		if req.ResolvedTargetType != "" {
			entities["targetType"] = req.ResolvedTargetType
		}
	}
	if entities["checkType"] == "" && req.ResolvedCheckType != "" {
		entities["checkType"] = req.ResolvedCheckType
	}

	// 文件操作不继承历史 queryType / directoryPath，避免把“清空日志”误污染为目录分析。
	if entities["fileOp"] == "" {
		if entities["queryType"] == "" && req.ResolvedQueryType != "" {
			entities["queryType"] = req.ResolvedQueryType
		}
		if entities["directoryPath"] == "" && req.ResolvedDirectoryPath != "" {
			entities["directoryPath"] = req.ResolvedDirectoryPath
		}
	}

	return entities
}

func (s *AgentChatServiceImpl) buildManagedFileClearResponse(hostIP, filePath, directoryPath string, response *model.AgentChatResponse) (string, error) {
	sizeOutput, err := s.runManagedFileFollowUpCommand(hostIP, fmt.Sprintf("if test -e %s; then du -sh %s; else printf 'NOT_FOUND\\t%s\\n'; fi", filePath, filePath, filePath))
	if err != nil {
		return "", err
	}

	result, err := s.performanceAnalysis.AnalyzeDirectoryDeepByIP(hostIP, directoryPath)
	if err != nil {
		return "", err
	}

	var report strings.Builder
	report.WriteString("# 文件清空结果\n\n")
	report.WriteString("- 文件: `")
	report.WriteString(filePath)
	report.WriteString("`\n")
	report.WriteString("- 结果: ")
	report.WriteString(formatManagedFileStatus(response, "clear"))
	report.WriteString("\n\n")
	report.WriteString("## 清空后文件大小\n\n")
	report.WriteString("```plaintext\n")
	report.WriteString(strings.TrimSpace(sizeOutput))
	report.WriteString("\n```\n\n")
	report.WriteString(buildTopFilesSection(result.TopFiles))

	return report.String(), nil
}

func (s *AgentChatServiceImpl) buildManagedFileDeleteResponse(hostIP, filePath, directoryPath string, response *model.AgentChatResponse) (string, error) {
	listOutput, err := s.runManagedFileFollowUpCommand(hostIP, fmt.Sprintf("ls -1A %s 2>/dev/null", directoryPath))
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(listOutput) == "" {
		listOutput = "(empty)"
	}

	var report strings.Builder
	report.WriteString("# 文件删除结果\n\n")
	report.WriteString("- 文件: `")
	report.WriteString(filePath)
	report.WriteString("`\n")
	report.WriteString("- 结果: ")
	report.WriteString(formatManagedFileStatus(response, "delete"))
	report.WriteString("\n\n")
	report.WriteString("## 当前目录剩余文件\n\n")
	report.WriteString("```plaintext\n")
	report.WriteString(strings.TrimSpace(listOutput))
	report.WriteString("\n```\n")

	return report.String(), nil
}

func (s *AgentChatServiceImpl) runManagedFileFollowUpCommand(hostIP, command string) (string, error) {
	host, err := s.performanceAnalysis.hostDao.GetCmdbHostByIP(hostIP)
	if err != nil {
		return "", err
	}

	sshConfig, err := s.performanceAnalysis.getSSHConfig(&host)
	if err != nil {
		return "", err
	}

	output, err := s.performanceAnalysis.sshUtil.ExecuteRemoteCommand(sshConfig, command)
	if err != nil {
		return "", err
	}

	return output, nil
}

func formatManagedFileStatus(response *model.AgentChatResponse, fileOp string) string {
	status := extractManagedFileStatus(response)
	switch fileOp {
	case "clear":
		switch status {
		case "cleared":
			return "清空成功"
		case "not_found":
			return "文件不存在"
		default:
			return "清空已执行"
		}
	case "delete":
		switch status {
		case "deleted":
			return "删除成功"
		case "not_found":
			return "文件不存在"
		default:
			return "删除已执行"
		}
	default:
		return "已执行"
	}
}

func extractManagedFileStatus(response *model.AgentChatResponse) string {
	if response == nil {
		return ""
	}
	for _, line := range strings.Split(response.Message, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "status=") {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "status="))
		}
	}
	if response.TaskInfo != nil {
		if status, ok := response.TaskInfo["status"].(string); ok {
			return strings.TrimSpace(status)
		}
	}
	return ""
}

func buildTopFilesSection(topFiles []model.DirectoryInfo) string {
	var report strings.Builder
	report.WriteString("## 最大文件 TOP 5\n\n")
	report.WriteString("| 排名 | 文件路径 | 大小 | 删除建议 |\n")
	report.WriteString("|------|----------|------|----------|\n")
	for i, file := range topFiles {
		report.WriteString(fmt.Sprintf("| %d | %s | **%s** | %s |\n", i+1, file.Path, file.Size, classifyFileDeletionPolicy(file.Path)))
	}
	if len(topFiles) == 0 {
		report.WriteString("| - | - | - | - |\n")
	}
	report.WriteString("\n")
	return report.String()
}

// executeBatchAnalysis 执行批量分析
func (s *AgentChatServiceImpl) executeBatchAnalysis(ctx context.Context, sessionID string, skillName string, skillResp *baseSkill.SkillResponse) (*model.AgentChatResponse, error) {
	if len(skillResp.BatchResult) == 0 {
		return nil, fmt.Errorf("批量结果为空")
	}

	var reports []string
	var successCount int

	// 对于命令执行和脚本执行，直接返回格式化的结果
	if skillName == "command_execution" || skillName == "script_execution" {
		for _, item := range skillResp.BatchResult {
			// 检查是否有错误
			if errMsg, ok := item["error"].(string); ok {
				if target, ok := item["target"].(string); ok {
					reports = append(reports, fmt.Sprintf("❌ **主机 %s**: %s", target, errMsg))
				} else {
					hostName := item["hostName"].(string)
					hostIP := item["hostIP"].(string)
					reports = append(reports, fmt.Sprintf("❌ **主机 %s (%s)**: %s", hostName, hostIP, errMsg))
				}
				continue
			}

			// 成功执行
			hostName := item["hostName"].(string)
			hostIP := item["hostIP"].(string)
			result := item["result"].(string)

			if skillName == "command_execution" {
				command := item["command"].(string)
				reports = append(reports, fmt.Sprintf("✅ **主机 %s (%s)**\n\n**命令**: `%s`\n\n**输出**:\n```\n%s\n```", hostName, hostIP, command, result))
			} else {
				reports = append(reports, fmt.Sprintf("✅ **主机 %s (%s)**\n\n**输出**:\n```\n%s\n```", hostName, hostIP, result))
			}
			successCount++
		}

		// 生成汇总报告
		summary := fmt.Sprintf("## 批量执行报告\n\n共执行 %d 台主机，成功 %d 台\n\n---\n\n%s",
			len(skillResp.BatchResult), successCount, strings.Join(reports, "\n\n---\n\n"))

		return &model.AgentChatResponse{
			SessionID: sessionID,
			Message:   summary,
			Type:      "text",
			TaskInfo: map[string]interface{}{
				"taskType":     "batch_" + skillName,
				"totalHosts":   len(skillResp.BatchResult),
				"successCount": successCount,
			},
		}, nil
	}

	// 分析类技能的批量处理
	for _, item := range skillResp.BatchResult {
		// 检查是否有错误
		if errMsg, ok := item["error"].(string); ok {
			target := item["target"].(string)
			reports = append(reports, fmt.Sprintf("❌ **主机 %s**: %s", target, errMsg))
			continue
		}

		// 根据 Skill 类型执行分析
		var report string
		var err error

		switch skillName {
		case "cpu_analysis":
			report, err = s.executeSingleCPUAnalysis(ctx, item)
		case "memory_analysis":
			report, err = s.executeSingleMemoryAnalysis(ctx, item)
		case "disk_analysis":
			report, err = s.executeSingleDiskAnalysis(ctx, item)
		default:
			err = fmt.Errorf("不支持的批量分析类型: %s", skillName)
		}

		if err != nil {
			hostName := item["hostName"].(string)
			reports = append(reports, fmt.Sprintf("❌ **主机 %s**: 分析失败 - %v", hostName, err))
		} else {
			reports = append(reports, report)
			successCount++
		}
	}

	// 生成汇总报告
	summary := fmt.Sprintf("## 批量分析报告\n\n共分析 %d 台主机，成功 %d 台\n\n---\n\n%s",
		len(skillResp.BatchResult), successCount, strings.Join(reports, "\n\n---\n\n"))

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   summary,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType":     "batch_" + skillName,
			"totalHosts":   len(skillResp.BatchResult),
			"successCount": successCount,
		},
	}, nil
}

// executeSingleCPUAnalysis 执行单个主机的CPU分析（用于批量处理）
func (s *AgentChatServiceImpl) executeSingleCPUAnalysis(ctx context.Context, taskInfo map[string]interface{}) (string, error) {
	hostIP, ok := taskInfo["hostIP"].(string)
	if !ok {
		return "", fmt.Errorf("无法获取主机IP")
	}

	hostName, _ := taskInfo["hostName"].(string)
	needAnalysis, _ := taskInfo["needAnalysis"].(bool)

	// 调用性能分析服务获取CPU数据
	result, err := s.performanceAnalysis.AnalyzeServerPressureByIP(hostIP)
	if err != nil {
		return "", err
	}

	// 生成CPU分析报告
	report := fmt.Sprintf("### %s (%s)\n\n%s", hostName, hostIP, s.generateCPUReport(result, needAnalysis))
	return report, nil
}

// executeSingleMemoryAnalysis 执行单个主机的内存分析（用于批量处理）
func (s *AgentChatServiceImpl) executeSingleMemoryAnalysis(ctx context.Context, taskInfo map[string]interface{}) (string, error) {
	hostIP, ok := taskInfo["hostIP"].(string)
	if !ok {
		return "", fmt.Errorf("无法获取主机IP")
	}

	hostName, _ := taskInfo["hostName"].(string)

	// 调用性能分析服务
	result, err := s.performanceAnalysis.AnalyzeMemoryIssueByIP(hostIP)
	if err != nil {
		return "", err
	}

	// 生成报告
	report := fmt.Sprintf("### %s (%s)\n\n%s", hostName, hostIP, s.performanceAnalysis.GenerateMemoryAnalysisReport(result))
	return report, nil
}

// executeSingleDiskAnalysis 执行单个主机的磁盘分析（用于批量处理）
func (s *AgentChatServiceImpl) executeSingleDiskAnalysis(ctx context.Context, taskInfo map[string]interface{}) (string, error) {
	hostIP, ok := taskInfo["hostIP"].(string)
	if !ok {
		return "", fmt.Errorf("无法获取主机IP")
	}

	hostName, _ := taskInfo["hostName"].(string)

	// 调用性能分析服务
	result, err := s.performanceAnalysis.AnalyzeServerPressureByIP(hostIP)
	if err != nil {
		return "", err
	}

	// 生成磁盘分析报告，默认不显示目录
	report := fmt.Sprintf("### %s (%s)\n\n%s", hostName, hostIP, s.generateDiskReport(result, false))
	return report, nil
}

// executePressureAnalysis 执行压力分析
func (s *AgentChatServiceImpl) executePressureAnalysis(ctx context.Context, sessionID string, skillResp *baseSkill.SkillResponse) (*model.AgentChatResponse, error) {
	hostIP, ok := skillResp.TaskInfo["hostIP"].(string)
	if !ok {
		return nil, fmt.Errorf("无法获取主机IP")
	}

	// 调用性能分析服务
	result, err := s.performanceAnalysis.AnalyzeServerPressureByIP(hostIP)
	if err != nil {
		return nil, err
	}

	// 生成报告
	report := s.performanceAnalysis.GeneratePressureReport(result)

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType":      "pressure_analysis",
			"hostId":        result.HostID,
			"hostName":      result.HostName,
			"pressureLevel": result.PressureLevel,
		},
	}, nil
}

// executeMemoryAnalysis 执行内存分析
func (s *AgentChatServiceImpl) executeMemoryAnalysis(ctx context.Context, sessionID string, skillResp *baseSkill.SkillResponse) (*model.AgentChatResponse, error) {
	hostIP, ok := skillResp.TaskInfo["hostIP"].(string)
	if !ok {
		return nil, fmt.Errorf("无法获取主机IP")
	}

	// 调用性能分析服务
	result, err := s.performanceAnalysis.AnalyzeMemoryIssueByIP(hostIP)
	if err != nil {
		return nil, err
	}

	// 生成报告
	report := s.performanceAnalysis.GenerateMemoryAnalysisReport(result)

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType":      "memory_analysis",
			"hostId":        result.HostID,
			"hostName":      result.HostName,
			"memoryPattern": result.MemoryPattern,
		},
	}, nil
}

// executeCPUAnalysis 执行CPU分析
func (s *AgentChatServiceImpl) executeCPUAnalysis(ctx context.Context, sessionID string, skillResp *baseSkill.SkillResponse) (*model.AgentChatResponse, error) {
	hostIP, ok := skillResp.TaskInfo["hostIP"].(string)
	if !ok {
		return nil, fmt.Errorf("无法获取主机IP")
	}

	needAnalysis, _ := skillResp.TaskInfo["needAnalysis"].(bool)

	// 调用性能分析服务获取CPU数据
	result, err := s.performanceAnalysis.AnalyzeServerPressureByIP(hostIP)
	if err != nil {
		return nil, err
	}

	// 生成CPU分析报告
	report := s.generateCPUReport(result, needAnalysis)

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType": "cpu_analysis",
			"hostId":   result.HostID,
			"hostName": result.HostName,
		},
	}, nil
}

// executeDiskAnalysis 执行磁盘分析
func (s *AgentChatServiceImpl) executeDiskAnalysis(ctx context.Context, sessionID string, skillResp *baseSkill.SkillResponse) (*model.AgentChatResponse, error) {
	hostIP, ok := skillResp.TaskInfo["hostIP"].(string)
	if !ok {
		return nil, fmt.Errorf("无法获取主机IP")
	}

	showDirectories, _ := skillResp.TaskInfo["showDirectories"].(bool)
	report := ""
	topDirectoryPath := ""
	hostName, _ := skillResp.TaskInfo["hostName"].(string)
	hostID, _ := skillResp.TaskInfo["hostId"].(uint)

	if showDirectories {
		diskResult, err := s.performanceAnalysis.AnalyzeDiskDirectoriesByIP(hostIP)
		if err != nil {
			return nil, err
		}
		report = s.performanceAnalysis.GenerateDiskAnalysisReport(diskResult)
		hostName = diskResult.HostName
		hostID = diskResult.HostID
		if len(diskResult.Directories) > 0 {
			topDirectoryPath = diskResult.Directories[0].Path
		}
	} else {
		result, err := s.performanceAnalysis.AnalyzeServerPressureByIP(hostIP)
		if err != nil {
			return nil, err
		}
		report = s.generateDiskReport(result, false)
		hostName = result.HostName
		hostID = result.HostID
	}

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType":         "disk_analysis",
			"hostId":           hostID,
			"hostName":         hostName,
			"topDirectoryPath": topDirectoryPath,
		},
	}, nil
}

func (s *AgentChatServiceImpl) executeDirectoryDeepAnalysis(ctx context.Context, sessionID string, skillResp *baseSkill.SkillResponse) (*model.AgentChatResponse, error) {
	hostIP, ok := skillResp.TaskInfo["hostIP"].(string)
	if !ok {
		return nil, fmt.Errorf("无法获取主机IP")
	}
	directoryPath, _ := skillResp.TaskInfo["directoryPath"].(string)
	if strings.TrimSpace(directoryPath) == "" {
		return nil, fmt.Errorf("无法获取目录路径")
	}
	queryType, _ := skillResp.TaskInfo["queryType"].(string)

	if queryType == "list_files" {
		report, err := s.buildDirectoryFileListResponse(hostIP, directoryPath)
		if err != nil {
			return nil, err
		}
		return &model.AgentChatResponse{
			SessionID: sessionID,
			Message:   report,
			Type:      "text",
			TaskInfo: map[string]interface{}{
				"taskType":      "directory_file_list",
				"directoryPath": directoryPath,
			},
		}, nil
	}

	result, err := s.performanceAnalysis.AnalyzeDirectoryDeepByIP(hostIP, directoryPath)
	if err != nil {
		return nil, err
	}

	report := s.performanceAnalysis.GenerateDirectoryDeepAnalysisReport(result)

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType":      "directory_deep_analysis",
			"hostId":        result.HostID,
			"hostName":      result.HostName,
			"directoryPath": result.AnalyzedPath,
		},
	}, nil
}

func (s *AgentChatServiceImpl) buildDirectoryFileListResponse(hostIP, directoryPath string) (string, error) {
	output, err := s.runManagedFileFollowUpCommand(hostIP, fmt.Sprintf("ls -lah %s 2>/dev/null", directoryPath))
	if err != nil {
		return "", err
	}

	var report strings.Builder
	report.WriteString("# 目录文件列表\n\n")
	report.WriteString("**目录**: ")
	report.WriteString(directoryPath)
	report.WriteString("\n\n")
	report.WriteString("```plaintext\n")
	report.WriteString(strings.TrimSpace(output))
	report.WriteString("\n```\n")

	return report.String(), nil
}

// executeServiceQuery 执行服务查询
func (s *AgentChatServiceImpl) executeServiceQuery(ctx context.Context, sessionID string, skillResp *baseSkill.SkillResponse) (*model.AgentChatResponse, error) {
	hostIP, ok := skillResp.TaskInfo["hostIP"].(string)
	if !ok {
		return nil, fmt.Errorf("无法获取主机IP")
	}

	hostName, _ := skillResp.TaskInfo["hostName"].(string)

	// 调用新的服务查询方法
	services, err := s.performanceAnalysis.QueryRunningServicesByIP(hostIP)
	if err != nil {
		return nil, err
	}

	// 生成服务报告
	report := s.generateServiceReportFromList(hostName, hostIP, services)

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType":     "service_query",
			"hostIP":       hostIP,
			"serviceCount": len(services),
		},
	}, nil
}

// recognizeIntent 意图识别
func (s *AgentChatServiceImpl) recognizeIntent(ctx context.Context, message string) (*model.IntentResult, error) {
	// TODO: 需要重构AI服务以支持非Gin Context调用
	// 暂时使用降级的关键词匹配方式
	return s.fallbackIntentRecognition(message), nil
}

// fallbackIntentRecognition 降级意图识别（关键词匹配）
func (s *AgentChatServiceImpl) fallbackIntentRecognition(message string) *model.IntentResult {
	messageLower := strings.ToLower(message)
	entities := s.extractEntities(message)
	isDirectoryInspection := entities["checkType"] == "disk" &&
		(entities["queryType"] == "directory_deep_analysis" || entities["directoryPath"] != "") &&
		(strings.Contains(messageLower, "检查") ||
			strings.Contains(messageLower, "查看") ||
			strings.Contains(messageLower, "分析") ||
			strings.Contains(messageLower, "有哪些") ||
			strings.Contains(messageLower, "文件") ||
			strings.Contains(messageLower, "目录"))

	// 添加调试日志
	log.Log().Infof("[意图识别] 原始消息: %s", message)
	log.Log().Infof("[意图识别] 转换后消息: %s", messageLower)

	if actionType := entities["actionType"]; actionType == "execute_command" || actionType == "run_command" ||
		actionType == "execute_script" || actionType == "run_script" ||
		actionType == "install_service" || actionType == "restart_service" || actionType == "query_service" {
		log.Log().Infof("[意图识别] 识别为命令/脚本执行，提取实体: %+v", entities)
		return &model.IntentResult{
			Type:       "incident",
			Action:     "execute",
			Confidence: 0.8,
			Entities:   entities,
		}
	}

	// 故障/压力分析相关关键词 - 优先级最高（因为更具体）
	if strings.Contains(messageLower, "故障") || strings.Contains(messageLower, "问题") ||
		strings.Contains(messageLower, "异常") || strings.Contains(messageLower, "错误") ||
		strings.Contains(messageLower, "压力") || strings.Contains(messageLower, "分析") ||
		strings.Contains(messageLower, "负载") || strings.Contains(messageLower, "性能") ||
		strings.Contains(messageLower, "卡顿") || strings.Contains(messageLower, "慢") ||
		strings.Contains(messageLower, "内存") || strings.Contains(messageLower, "cpu") ||
		strings.Contains(messageLower, "磁盘") || strings.Contains(messageLower, "端口") ||
		strings.Contains(messageLower, "进程") ||
		strings.Contains(messageLower, "为什么") || strings.Contains(messageLower, "怎么") ||
		isDirectoryInspection {

		log.Log().Infof("[意图识别] 识别为 incident 类型，提取实体: %+v", entities)

		return &model.IntentResult{
			Type:       "incident",
			Action:     "analyze",
			Confidence: 0.7,
			Entities:   entities,
		}
	}

	// 巡检相关关键词 - 必须有明确目标
	if strings.Contains(messageLower, "巡检") ||
		(strings.Contains(messageLower, "检查") && (strings.Contains(messageLower, "健康") || strings.Contains(messageLower, "状态"))) ||
		(strings.Contains(messageLower, "扫描") && strings.Contains(messageLower, "主机")) {

		// 必须有目标才执行巡检，否则可能是通用技术问题（如"写一个巡检脚本"）
		if entities["targetType"] != "" || entities["targetName"] != "" || entities["targetIP"] != "" {
			log.Log().Infof("[意图识别] 识别为 inspection 类型，提取实体: %+v", entities)
			return &model.IntentResult{
				Type:       "inspection",
				Action:     "execute",
				Confidence: 0.7,
				Entities:   entities,
			}
		}
		log.Log().Infof("[意图识别] 包含巡检关键词但无明确目标，判定为通用对话")
	}

	// 命令执行相关关键词（新增 - 优先级高）
	if strings.Contains(messageLower, "执行命令") || strings.Contains(messageLower, "运行命令") ||
		strings.Contains(messageLower, "执行脚本") || strings.Contains(messageLower, "运行脚本") ||
		(strings.Contains(messageLower, "批量") && (strings.Contains(messageLower, "执行") || strings.Contains(messageLower, "运行"))) {
		log.Log().Infof("[意图识别] 识别为命令/脚本执行，提取实体: %+v", entities)
		return &model.IntentResult{
			Type:       "incident",
			Action:     "execute",
			Confidence: 0.8,
			Entities:   entities,
		}
	}

	// 服务查询相关关键词（新增）
	if (strings.Contains(messageLower, "服务") || strings.Contains(messageLower, "应用")) &&
		(strings.Contains(messageLower, "查看") || strings.Contains(messageLower, "查询") ||
			strings.Contains(messageLower, "有哪些") || strings.Contains(messageLower, "部署") ||
			strings.Contains(messageLower, "运行") || strings.Contains(messageLower, "启动")) {
		return &model.IntentResult{
			Type:       "incident",
			Action:     "analyze",
			Confidence: 0.7,
			Entities:   entities,
		}
	}

	// 容量预测相关关键词
	if strings.Contains(messageLower, "预测") || strings.Contains(messageLower, "容量") ||
		strings.Contains(messageLower, "磁盘") && (strings.Contains(messageLower, "满") || strings.Contains(messageLower, "不足")) ||
		strings.Contains(messageLower, "资源") && strings.Contains(messageLower, "使用") {
		return &model.IntentResult{
			Type:       "prediction",
			Action:     "predict",
			Confidence: 0.7,
			Entities:   entities,
		}
	}

	// 默认为一般对话
	return &model.IntentResult{
		Type:       "general",
		Action:     "chat",
		Confidence: 0.5,
		Entities:   make(map[string]string),
	}
}

// extractEntities 提取实体（增强版，支持IP和主机名提取）
func (s *AgentChatServiceImpl) extractEntities(message string) map[string]string {
	entities := make(map[string]string)
	messageLower := strings.ToLower(message)
	directoryPath := extractDirectoryPathFromMessage(message)
	if directoryPath != "" {
		entities["directoryPath"] = directoryPath
	}

	// 0. 优先提取显式格式：name=xxx 或 ip=xxx，支持逗号分隔多个值
	// 但是要特别处理：如果消息中本身已经有多个IP，优先使用消息中的IP，而不是显式格式

	// 先尝试自动IP提取，看看消息中是否有多个IP
	autoIPMatches := findIPAddress(message)

	if strings.Contains(message, "name=") {
		parts := strings.Split(message, "name=")
		if len(parts) > 1 {
			targetName := extractExplicitTargetValue(parts[1])
			if isValidResolvedTarget(targetName) {
				entities["targetType"] = "host"
				entities["targetName"] = targetName
				if strings.Contains(targetName, ",") || strings.Contains(targetName, "，") {
					log.Log().Infof("提取到多个主机名(显式): %s", targetName)
				} else {
					log.Log().Infof("提取到主机名(显式): %s", targetName)
				}
			}
		}
	}

	// 处理 ip= 格式：如果自动提取到多个IP，优先使用自动提取的结果
	if entities["targetName"] == "" && strings.Contains(message, "ip=") {
		// 如果自动提取到了多个IP（2个以上），说明消息本身就有多个IP，使用自动提取结果
		if len(autoIPMatches) >= 2 {
			entities["targetType"] = "host"
			entities["targetName"] = strings.Join(autoIPMatches, ",")
			log.Log().Infof("检测到消息中有多个IP，使用自动提取结果: %s", entities["targetName"])
		} else {
			// 只有1个或0个IP，使用显式格式提取
			parts := strings.Split(message, "ip=")
			if len(parts) > 1 {
				targetIP := extractExplicitTargetValue(parts[1])
				if targetIP != "" {
					entities["targetType"] = "host"
					entities["targetName"] = targetIP
					if strings.Contains(targetIP, ",") || strings.Contains(targetIP, "，") {
						log.Log().Infof("提取到多个IP(显式): %s", targetIP)
					} else {
						log.Log().Infof("提取到IP(显式): %s", targetIP)
					}
				}
			}
		}
	}

	// 1. 提取IP地址 - 匹配常见的IP地址格式，支持多个IP
	if entities["targetName"] == "" {
		ipMatches := findIPAddress(message)
		if len(ipMatches) > 0 {
			entities["targetType"] = "host"
			if len(ipMatches) == 1 {
				// 单个IP
				entities["targetName"] = ipMatches[0]
				log.Log().Infof("提取到IP地址: %s", ipMatches[0])
			} else {
				// 多个IP，用逗号分隔
				entities["targetName"] = strings.Join(ipMatches, ",")
				log.Log().Infof("提取到多个IP地址: %s", entities["targetName"])
			}
		}
	}

	// 2. 提取主机名 - 匹配常见的主机名模式
	if entities["targetName"] == "" {
		words := strings.Fields(message)
		for _, word := range words {
			word = strings.Trim(word, ",.!?;:()[]{}<>*\"'`")
			// 检查是否像主机名（包含-或数字，或者包含点号）
			if (strings.Contains(word, "-") || strings.Contains(word, ".") || containsDigit(word)) &&
				!looksLikeSizeToken(word) &&
				!strings.Contains(word, "检查") &&
				!strings.Contains(word, "分析") &&
				!strings.Contains(word, "服务器") &&
				!strings.Contains(word, "主机") &&
				len(word) > 2 &&
				isValidResolvedTarget(word) {
				entities["targetType"] = "host"
				entities["targetName"] = word
				log.Log().Infof("提取到主机名: %s", word)
				break
			}
		}
	}

	// 3. 提取目标类型（如果没有具体主机，则提取类型）
	if entities["targetType"] == "" {
		if strings.Contains(messageLower, "所有") || strings.Contains(messageLower, "全部") {
			entities["targetType"] = "all"
		} else if strings.Contains(messageLower, "分组") || strings.Contains(messageLower, "组") {
			entities["targetType"] = "group"
		} else if strings.Contains(messageLower, "主机") || strings.Contains(messageLower, "服务器") {
			entities["targetType"] = "host"
		}
	}

	// 4. 提取检查类型 - 优先检查具体类型关键词
	// 注意：先检查更具体的关键词，避免被"服务器"误匹配
	if strings.Contains(messageLower, "内存") || strings.Contains(messageLower, "memory") {
		entities["checkType"] = "memory"
	} else if strings.Contains(messageLower, "cpu") || strings.Contains(messageLower, "处理器") {
		entities["checkType"] = "cpu"
	} else if strings.Contains(messageLower, "压力") || strings.Contains(messageLower, "负载") || strings.Contains(messageLower, "性能") {
		entities["checkType"] = "performance"
	} else if strings.Contains(messageLower, "磁盘") || strings.Contains(messageLower, "disk") || strings.Contains(messageLower, "空间") {
		entities["checkType"] = "disk"
	} else if directoryPath != "" || strings.Contains(messageLower, "目录") {
		entities["checkType"] = "disk"
	} else if strings.Contains(messageLower, "端口") || strings.Contains(messageLower, "port") {
		entities["checkType"] = "port"
	} else if strings.Contains(messageLower, "进程") || strings.Contains(messageLower, "process") {
		entities["checkType"] = "process"
	} else if strings.Contains(messageLower, "什么服务") || strings.Contains(messageLower, "哪些服务") ||
		strings.Contains(messageLower, "部署") && strings.Contains(messageLower, "服务") ||
		strings.Contains(messageLower, "运行") && strings.Contains(messageLower, "服务") ||
		strings.Contains(messageLower, "启动") && strings.Contains(messageLower, "服务") {
		// 检测询问服务相关的问题
		entities["checkType"] = "service"
	} else if strings.Contains(messageLower, "系统") {
		entities["checkType"] = "system"
	} else if strings.Contains(messageLower, "安全") {
		entities["checkType"] = "security"
	}

	// 5. 提取详细的分析原因（"为什么"）
	if strings.Contains(messageLower, "为什么") || strings.Contains(messageLower, "原因") {
		entities["needAnalysis"] = "true"
	}

	// 6. 提取操作类型（执行命令、执行脚本等）
	if serviceName := extractServiceName(messageLower); serviceName != "" {
		switch {
		case strings.Contains(messageLower, "安装") || strings.Contains(messageLower, "部署"):
			entities["actionType"] = "install_service"
			entities["checkType"] = "service"
			entities["serviceName"] = serviceName
		case strings.Contains(messageLower, "重启") || strings.Contains(messageLower, "restart"):
			entities["actionType"] = "restart_service"
			entities["checkType"] = "service"
			entities["serviceName"] = serviceName
		case strings.Contains(messageLower, "检查") || strings.Contains(messageLower, "查看") ||
			strings.Contains(messageLower, "状态") || strings.Contains(messageLower, "是否启动") ||
			strings.Contains(messageLower, "运行情况"):
			entities["actionType"] = "query_service"
			entities["checkType"] = "service"
			entities["serviceName"] = serviceName
		}
	} else if strings.Contains(messageLower, "执行命令") || strings.Contains(messageLower, "运行命令") {
		entities["actionType"] = "execute_command"
		// 提取命令
		if strings.Contains(message, "执行命令") {
			parts := strings.SplitN(message, "执行命令", 2)
			if len(parts) > 1 {
				cmd := strings.TrimSpace(parts[1])
				// 去除"在"、"上"等词
				cmd = strings.TrimPrefix(cmd, "在")
				cmd = strings.TrimSpace(cmd)
				// 去除末尾的 ip= 或 name= 参数（这些参数是自动追加的，不是命令的一部分）
				if idx := strings.Index(cmd, " ip="); idx != -1 {
					cmd = strings.TrimSpace(cmd[:idx])
				}
				if idx := strings.Index(cmd, " name="); idx != -1 {
					cmd = strings.TrimSpace(cmd[:idx])
				}
				if cmd != "" {
					entities["command"] = cmd
				}
			}
		}
	} else if clearPath := extractClearableFilePath(message); clearPath != "" &&
		(strings.Contains(messageLower, "清空") || strings.Contains(messageLower, "清除")) {
		entities["actionType"] = "execute_command"
		entities["fileOp"] = "clear"
		entities["filePath"] = clearPath
	} else if deletePath := extractAnyFilePath(message); deletePath != "" &&
		strings.Contains(messageLower, "删除") {
		entities["actionType"] = "execute_command"
		entities["fileOp"] = "delete"
		entities["filePath"] = deletePath
	} else if strings.Contains(messageLower, "执行脚本") || strings.Contains(messageLower, "运行脚本") {
		entities["actionType"] = "execute_script"
		// 提取脚本内容或路径
		if strings.Contains(message, "执行脚本") {
			parts := strings.SplitN(message, "执行脚本", 2)
			if len(parts) > 1 {
				script := strings.TrimSpace(parts[1])
				if script != "" {
					entities["script"] = script
					// 判断是文件路径还是内联脚本
					if strings.HasPrefix(script, "/") || strings.HasSuffix(script, ".sh") {
						entities["scriptType"] = "file"
					} else {
						entities["scriptType"] = "inline"
					}
				}
			}
		}
	}

	// 7. 提取特定查询（"哪些目录"、"高危端口"等）
	if shouldListDirectoryFiles(messageLower, directoryPath) {
		entities["queryType"] = "list_files"
	} else if shouldUseDirectoryDeepAnalysis(messageLower, directoryPath) {
		entities["queryType"] = "directory_deep_analysis"
	} else if strings.Contains(messageLower, "哪些目录") || strings.Contains(messageLower, "哪个目录") || strings.Contains(messageLower, "目录占用") {
		entities["queryType"] = "directories"
	}
	if strings.Contains(messageLower, "高危端口") || strings.Contains(messageLower, "危险端口") {
		entities["queryType"] = "risky_ports"
	}
	if strings.Contains(messageLower, "异常") || strings.Contains(messageLower, "僵尸") {
		entities["queryType"] = "abnormal"
	}

	return entities
}

func extractDirectoryPathFromMessage(message string) string {
	for i := 0; i < len(message); i++ {
		if message[i] != '/' {
			continue
		}
		j := i
		for j < len(message) {
			c := message[j]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
				c == '/' || c == '-' || c == '_' || c == '.' || c == '@' {
				j++
				continue
			}
			break
		}
		if j > i {
			field := strings.TrimSpace(message[i:j])
			field = strings.Trim(field, "`*\"'，。,:：;；()[]{}<>")
			if strings.HasPrefix(field, "/") && isSafeDirectoryPath(field) {
				return field
			}
		}
	}
	return ""
}

func extractClearableFilePath(message string) string {
	for i := 0; i < len(message); i++ {
		if message[i] != '/' {
			continue
		}
		j := i
		for j < len(message) {
			c := message[j]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
				c == '/' || c == '-' || c == '_' || c == '.' || c == '@' {
				j++
				continue
			}
			break
		}
		if j > i {
			field := strings.TrimSpace(message[i:j])
			field = strings.Trim(field, "`*\"'，。,:：;；()[]{}<>")
			if strings.HasPrefix(field, "/") && isSafeDirectoryPath(field) && isAllowedClearFile(field) {
				return field
			}
		}
	}
	return ""
}

func extractAnyFilePath(message string) string {
	for i := 0; i < len(message); i++ {
		if message[i] != '/' {
			continue
		}
		j := i
		for j < len(message) {
			c := message[j]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
				c == '/' || c == '-' || c == '_' || c == '.' || c == '@' {
				j++
				continue
			}
			break
		}
		if j > i {
			field := strings.TrimSpace(message[i:j])
			field = strings.Trim(field, "`*\"'，。,:：;；()[]{}<>")
			if strings.HasPrefix(field, "/") && isSafeDirectoryPath(field) {
				return field
			}
		}
	}
	return ""
}

func shouldListDirectoryFiles(messageLower, directoryPath string) bool {
	if directoryPath == "" {
		return false
	}
	return (strings.Contains(messageLower, "有哪些文件") ||
		strings.Contains(messageLower, "下面有哪些文件") ||
		strings.Contains(messageLower, "目录下有哪些文件") ||
		strings.Contains(messageLower, "列出文件") ||
		strings.Contains(messageLower, "看看文件")) &&
		!strings.Contains(messageLower, "删除") &&
		!strings.Contains(messageLower, "清空")
}

func extractServiceName(messageLower string) string {
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

func shouldUseDirectoryDeepAnalysis(messageLower, directoryPath string) bool {
	if isDirectFileOperationMessage(messageLower, directoryPath) {
		return false
	}
	if directoryPath != "" {
		if isAllowedClearFile(directoryPath) {
			return false
		}
		return true
	}
	return strings.Contains(messageLower, "这个目录") ||
		strings.Contains(messageLower, "该目录") ||
		strings.Contains(messageLower, "目录里面") ||
		(strings.Contains(messageLower, "里面") && (strings.Contains(messageLower, "文件") || strings.Contains(messageLower, "删除") || strings.Contains(messageLower, "清理")))
}

func isDirectFileOperationMessage(messageLower, path string) bool {
	if path == "" {
		return false
	}
	if !(strings.Contains(messageLower, "清空") || strings.Contains(messageLower, "清除") || strings.Contains(messageLower, "删除")) {
		return false
	}
	if strings.Contains(messageLower, "目录") || strings.Contains(messageLower, "里面") || strings.Contains(messageLower, "下面") {
		return false
	}
	return true
}

func isSafeDirectoryPath(path string) bool {
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

func isAllowedClearFile(path string) bool {
	lower := strings.ToLower(strings.TrimSpace(path))
	return strings.HasSuffix(lower, ".log") ||
		strings.HasSuffix(lower, ".txt") ||
		strings.HasSuffix(lower, ".gz") ||
		strings.HasSuffix(lower, ".zip")
}

func extractExplicitTargetValue(raw string) string {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		switch r {
		case ' ', '\n', '\r', '\t', '。', '、', '，', ',', ';', '；', ':', '：', ')', '(', ']', '[', '}', '{', '>', '<', '*', '"', '\'', '`':
			return true
		default:
			return false
		}
	})
	if len(fields) == 0 {
		return ""
	}
	return strings.TrimSpace(fields[0])
}

func isValidResolvedTarget(target string) bool {
	target = strings.TrimSpace(target)
	if target == "" {
		return false
	}
	if looksLikeSizeToken(target) {
		return false
	}
	if strings.Contains(target, "/") || strings.Contains(target, "\\") || strings.Contains(target, "**") {
		return false
	}
	lower := strings.ToLower(target)
	for _, fragment := range []string{"目录", "文件", "占用", "空间"} {
		if strings.Contains(lower, fragment) {
			return false
		}
	}
	if strings.HasPrefix(target, "(") || strings.HasSuffix(target, ")") {
		return false
	}
	return true
}

func looksLikeSizeToken(value string) bool {
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

// findIPAddress 查找IP地址，支持多种分隔符（空格、逗号、换行符）
func findIPAddress(text string) []string {
	var matches []string

	// 先统一分隔符，将换行符、中文逗号替换为空格
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.ReplaceAll(text, "，", " ")
	text = strings.ReplaceAll(text, ",", " ")

	parts := strings.Split(text, " ")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 简化的IP匹配：检查是否是点分十进制格式
		if strings.Count(part, ".") == 3 {
			segments := strings.Split(part, ".")
			if len(segments) == 4 {
				isIP := true
				for _, seg := range segments {
					// 移除可能的标点符号
					seg = strings.Trim(seg, ",.!?;:")
					if len(seg) == 0 || len(seg) > 3 {
						isIP = false
						break
					}
					// 检查是否全是数字
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
					// 清理标点符号
					cleanIP := strings.Trim(part, ",.!?;:")
					matches = append(matches, cleanIP)
				}
			}
		}
	}
	return matches
}

// containsDigit 检查字符串是否包含数字
func containsDigit(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}

// parseIntentResponse 解析AI意图识别响应
func (s *AgentChatServiceImpl) parseIntentResponse(response string) (*model.IntentResult, error) {
	// 尝试从响应中提取JSON
	response = strings.TrimSpace(response)

	// 查找JSON部分
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start == -1 || end == -1 || start >= end {
		// 如果没有找到JSON，使用降级策略
		return &model.IntentResult{
			Type:       "general",
			Confidence: 0.5,
			Entities:   make(map[string]string),
		}, nil
	}

	jsonStr := response[start : end+1]

	var intent model.IntentResult
	if err := json.Unmarshal([]byte(jsonStr), &intent); err != nil {
		log.Log().Warnf("解析意图JSON失败: %v, 响应: %s", err, jsonStr)
		return &model.IntentResult{
			Type:       "general",
			Confidence: 0.5,
			Entities:   make(map[string]string),
		}, nil
	}

	return &intent, nil
}

// handleInspection 处理巡检请求
func (s *AgentChatServiceImpl) handleInspection(ctx context.Context, sessionID string, userID uint, intent *model.IntentResult) (*model.AgentChatResponse, error) {
	// 1. 解析目标
	targets, err := s.resolveTargets(intent.Entities)
	if err != nil {
		return nil, fmt.Errorf("无法解析目标主机: %v", err)
	}

	if len(targets) == 0 {
		return s.createTextResponse(sessionID, "请指定要巡检的主机或分组，例如：\n- 检查生产环境的所有服务器\n- 巡检web分组\n- 检查主机192.168.1.100"), nil
	}

	// 2. 确定巡检类型
	checkType := intent.Entities["checkType"]
	if checkType == "" {
		checkType = "system" // 默认系统巡检
	}

	// 3. 创建临时巡检任务
	targetIDs := make([]uint, len(targets))
	for i, host := range targets {
		targetIDs[i] = host.ID
	}

	targetIDsJSON, _ := json.Marshal(targetIDs)
	inspection := &model.AgentInspection{
		Name:         fmt.Sprintf("临时巡检-%s", time.Now().Format("20060102150405")),
		Type:         checkType,
		TargetType:   "host",
		TargetIDs:    string(targetIDsJSON),
		ScheduleType: "manual",
		Enabled:      1,
		Template:     "{}",
		Remark:       fmt.Sprintf("AI对话触发，会话ID: %s", sessionID),
	}

	// TODO: 需要重构InspectionService以支持非Gin Context调用
	// 暂时跳过创建巡检任务
	_ = inspection

	// 5. 返回响应
	message := fmt.Sprintf("✅ 已为 %d 台主机创建巡检任务\n\n", len(targets))
	message += "目标主机：\n"
	for i, host := range targets {
		if i < 5 { // 最多显示5台
			message += fmt.Sprintf("- %s (%s)\n", host.HostName, host.SSHIP)
		}
	}
	if len(targets) > 5 {
		message += fmt.Sprintf("- ... 及其他 %d 台主机\n", len(targets)-5)
	}
	message += "\n巡检任务正在执行中，稍后可查看巡检报告..."

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   message,
		Type:      "task",
		Intent:    intent,
		TaskInfo: map[string]interface{}{
			"taskType":    "inspection",
			"targetCount": len(targets),
		},
	}, nil
}

// handleQuery 处理查询请求
func (s *AgentChatServiceImpl) handleQuery(ctx context.Context, sessionID string, userID uint, intent *model.IntentResult) (*model.AgentChatResponse, error) {
	// 解析目标
	targets, err := s.resolveTargets(intent.Entities)
	if err != nil || len(targets) == 0 {
		return s.createTextResponse(sessionID, "请指定要查询的主机，例如：查询web-server-01的状态"), nil
	}

	// 构建查询结果
	message := fmt.Sprintf("📊 查询结果（共 %d 台主机）\n\n", len(targets))
	for i, host := range targets {
		if i < 10 { // 最多显示10台
			status := "🟢 在线"
			if host.Status != 1 {
				status = "🔴 离线"
			}
			message += fmt.Sprintf("**%s** (%s)\n", host.HostName, host.SSHIP)
			message += fmt.Sprintf("  状态: %s\n", status)
			message += fmt.Sprintf("  系统: %s\n", host.OS)
			if host.Remark != "" {
				message += fmt.Sprintf("  备注: %s\n", host.Remark)
			}
			message += "\n"
		}
	}
	if len(targets) > 10 {
		message += fmt.Sprintf("... 及其他 %d 台主机\n", len(targets)-10)
	}

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   message,
		Type:      "text",
		Intent:    intent,
	}, nil
}

// handleGeneral 处理一般对话
func (s *AgentChatServiceImpl) handleGeneral(ctx context.Context, sessionID string, message string) (*model.AgentChatResponse, error) {
	// TODO: 需要重构AI服务以支持非Gin Context调用
	// 暂时返回默认回复
	defaultMessage := `你好，我是你的 **AutoOps 运维智能小助手**。

目前我可以帮你处理这些运维工作：

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

	return s.createTextResponse(sessionID, defaultMessage), nil
}

// handleIncident 处理故障/性能分析请求
func (s *AgentChatServiceImpl) handleIncident(ctx context.Context, sessionID string, userID uint, intent *model.IntentResult) (*model.AgentChatResponse, error) {
	// 1. 获取检查类型
	checkType := intent.Entities["checkType"]
	if checkType == "" {
		return s.createTextResponse(sessionID, "请指定分析类型，例如：\n- 检查服务器压力\n- 分析内存问题\n- 查看磁盘使用\n- 检查端口情况\n- 查看进程状态"), nil
	}

	// 2. 获取目标主机
	targetName := intent.Entities["targetName"]
	if targetName == "" {
		return s.createTextResponse(sessionID, "请指定要分析的主机IP或名称，例如：\n- 检查 ip=192.168.1.100 的服务器压力\n- 分析 name=web-server-01 的内存问题"), nil
	}

	// 3. 根据检查类型执行相应的分析
	switch checkType {
	case "performance":
		return s.handlePressureAnalysis(ctx, sessionID, targetName)
	case "memory":
		return s.handleMemoryAnalysis(ctx, sessionID, targetName)
	case "cpu":
		return s.handleCPUAnalysis(ctx, sessionID, targetName, intent)
	case "disk":
		return s.handleDiskAnalysis(ctx, sessionID, targetName, intent)
	case "port":
		return s.handlePortAnalysis(ctx, sessionID, targetName, intent)
	case "process":
		return s.handleProcessAnalysis(ctx, sessionID, targetName, intent)
	case "service":
		return s.handleServiceQuery(ctx, sessionID, targetName)
	default:
		return s.createTextResponse(sessionID, fmt.Sprintf("暂不支持 %s 类型的分析，目前支持：\n- 压力分析（performance）\n- 内存分析（memory）\n- CPU分析（cpu）\n- 磁盘分析（disk）\n- 端口检查（port）\n- 进程检查（process）\n- 服务查询（service）", checkType)), nil
	}
}

// handlePressureAnalysis 处理服务器压力分析
func (s *AgentChatServiceImpl) handlePressureAnalysis(ctx context.Context, sessionID string, target string) (*model.AgentChatResponse, error) {
	// 尝试通过IP分析
	result, err := s.performanceAnalysis.AnalyzeServerPressureByIP(target)
	if err != nil {
		// 如果IP查询失败，尝试通过名称分析
		result, err = s.performanceAnalysis.AnalyzeServerPressureByName(target)
		if err != nil {
			// 优化错误提示：如果是 record not found，提示配置密钥
			errMsg := err.Error()
			if strings.Contains(errMsg, "record not found") {
				return s.createErrorResponse(sessionID, fmt.Sprintf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", target)), nil
			}
			return s.createErrorResponse(sessionID, fmt.Sprintf("分析失败: %v", err)), nil
		}
	}

	// 生成报告
	report := s.performanceAnalysis.GeneratePressureReport(result)

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType":      "pressure_analysis",
			"hostId":        result.HostID,
			"hostName":      result.HostName,
			"pressureLevel": result.PressureLevel,
		},
	}, nil
}

// handleMemoryAnalysis 处理内存问题分析
func (s *AgentChatServiceImpl) handleMemoryAnalysis(ctx context.Context, sessionID string, target string) (*model.AgentChatResponse, error) {
	// 尝试通过名称分析
	result, err := s.performanceAnalysis.AnalyzeMemoryIssueByName(target)
	if err != nil {
		// 如果名称查询失败，尝试通过IP分析
		result, err = s.performanceAnalysis.AnalyzeMemoryIssueByIP(target)
		if err != nil {
			// 优化错误提示：如果是 record not found，提示配置密钥
			errMsg := err.Error()
			if strings.Contains(errMsg, "record not found") {
				return s.createErrorResponse(sessionID, fmt.Sprintf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", target)), nil
			}
			return s.createErrorResponse(sessionID, fmt.Sprintf("分析失败: %v", err)), nil
		}
	}

	// 生成报告
	report := s.performanceAnalysis.GenerateMemoryAnalysisReport(result)

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType":      "memory_analysis",
			"hostId":        result.HostID,
			"hostName":      result.HostName,
			"memoryPattern": result.MemoryPattern,
			"issueCount":    len(result.Issues),
		},
	}, nil
}

// resolveTargets 解析目标主机
func (s *AgentChatServiceImpl) resolveTargets(entities map[string]string) ([]*cmdbModel.CmdbHost, error) {
	targetType := entities["targetType"]
	targetName := entities["targetName"]

	var hosts []cmdbModel.CmdbHost
	var err error

	switch targetType {
	case "all":
		// 查询所有主机
		hosts = s.hostDao.GetCmdbHostList()
	case "group":
		// 根据分组名称查询
		if targetName != "" {
			// TODO: 实现根据分组名称查询主机
			hosts = s.hostDao.GetCmdbHostList()
		}
	case "host":
		// 根据主机名或IP查询
		if targetName != "" {
			// 尝试按IP查询
			host, err := s.hostDao.GetCmdbHostByIP(targetName)
			if err == nil {
				hosts = []cmdbModel.CmdbHost{host}
			} else {
				// 尝试按名称查询
				host, err := s.hostDao.GetCmdbHostByName(targetName)
				if err == nil {
					hosts = []cmdbModel.CmdbHost{host}
				}
			}
		}
	default:
		// 默认查询所有主机
		hosts = s.hostDao.GetCmdbHostList()
	}

	// 转换为指针类型
	hostPointers := make([]*cmdbModel.CmdbHost, len(hosts))
	for i := range hosts {
		hostPointers[i] = &hosts[i]
	}

	return hostPointers, err
}

// GetChatHistory 获取对话历史
func (s *AgentChatServiceImpl) GetChatHistory(ctx context.Context, userID uint, sessionID string, limit int) ([]model.AgentChatHistory, error) {
	return s.chatDao.GetBySessionID(sessionID, limit)
}

// ClearChatHistory 清除对话历史
func (s *AgentChatServiceImpl) ClearChatHistory(ctx context.Context, sessionID string) error {
	return s.chatDao.DeleteBySessionID(sessionID)
}

// 辅助方法

func (s *AgentChatServiceImpl) createTextResponse(sessionID, message string) *model.AgentChatResponse {
	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   message,
		Type:      "text",
	}
}

func (s *AgentChatServiceImpl) createErrorResponse(sessionID, errMsg string) *model.AgentChatResponse {
	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   "❌ " + errMsg,
		Type:      "error",
	}
}

func (s *AgentChatServiceImpl) getDefaultModelID() uint {
	// TODO: 从配置中获取默认模型ID
	return 1
}

// shouldRecordHistory 判断是否需要记录到数据库
// 只记录重要事件：inspection（巡检）、incident（故障/性能分析）、prediction（容量预测）
// 不记录简单查询：query（查询）、general（一般对话）
func (s *AgentChatServiceImpl) shouldRecordHistory(intentType string) bool {
	switch intentType {
	case "inspection", "incident", "prediction":
		return true // 重要事件，需要记录
	case "query", "general":
		return false // 简单查询，不需要记录
	default:
		return false // 未知类型，默认不记录
	}
}

// ============ 新增分析功能 ============

// handleCPUAnalysis 处理CPU分析
func (s *AgentChatServiceImpl) handleCPUAnalysis(ctx context.Context, sessionID string, target string, intent *model.IntentResult) (*model.AgentChatResponse, error) {
	// 尝试通过IP或名称查询主机
	host, err := s.getHostByTarget(target)
	if err != nil {
		// 优化错误提示
		errMsg := err.Error()
		if strings.Contains(errMsg, "未找到主机") {
			return s.createErrorResponse(sessionID, fmt.Sprintf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", target)), nil
		}
		return s.createErrorResponse(sessionID, fmt.Sprintf("未找到主机: %v", err)), nil
	}

	// 分析CPU使用率
	result, err := s.performanceAnalysis.AnalyzeServerPressureByIP(host.SSHIP)
	if err != nil {
		return s.createErrorResponse(sessionID, fmt.Sprintf("CPU分析失败: %v", err)), nil
	}

	// 生成CPU分析报告
	report := s.generateCPUReport(result, intent.Entities["needAnalysis"] == "true")

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType": "cpu_analysis",
			"hostId":   host.ID,
			"hostName": host.HostName,
		},
	}, nil
}

// handleDiskAnalysis 处理磁盘分析
func (s *AgentChatServiceImpl) handleDiskAnalysis(ctx context.Context, sessionID string, target string, intent *model.IntentResult) (*model.AgentChatResponse, error) {
	// 尝试通过IP或名称查询主机
	host, err := s.getHostByTarget(target)
	if err != nil {
		// 优化错误提示
		errMsg := err.Error()
		if strings.Contains(errMsg, "未找到主机") {
			return s.createErrorResponse(sessionID, fmt.Sprintf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", target)), nil
		}
		return s.createErrorResponse(sessionID, fmt.Sprintf("未找到主机: %v", err)), nil
	}

	// 获取磁盘信息
	result, err := s.performanceAnalysis.AnalyzeServerPressureByIP(host.SSHIP)
	if err != nil {
		return s.createErrorResponse(sessionID, fmt.Sprintf("磁盘分析失败: %v", err)), nil
	}

	// 生成磁盘分析报告
	report := s.generateDiskReport(result, intent.Entities["queryType"] == "directories")

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType": "disk_analysis",
			"hostId":   host.ID,
			"hostName": host.HostName,
		},
	}, nil
}

// handlePortAnalysis 处理端口分析
func (s *AgentChatServiceImpl) handlePortAnalysis(ctx context.Context, sessionID string, target string, intent *model.IntentResult) (*model.AgentChatResponse, error) {
	host, err := s.getHostByTarget(target)
	if err != nil {
		// 优化错误提示
		errMsg := err.Error()
		if strings.Contains(errMsg, "未找到主机") {
			return s.createErrorResponse(sessionID, fmt.Sprintf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", target)), nil
		}
		return s.createErrorResponse(sessionID, fmt.Sprintf("未找到主机: %v", err)), nil
	}

	// TODO: 需要实现端口扫描功能，通过 Agent 执行 ss -tuln 命令
	report := fmt.Sprintf("# 端口扫描报告\n\n**主机**: %s (%s)\n\n⚠️ 端口扫描功能开发中，敬请期待...\n\n建议手动执行：`ss -tuln | grep LISTEN`", host.HostName, host.SSHIP)

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType": "port_scan",
			"hostId":   host.ID,
		},
	}, nil
}

// handleProcessAnalysis 处理进程分析
func (s *AgentChatServiceImpl) handleProcessAnalysis(ctx context.Context, sessionID string, target string, intent *model.IntentResult) (*model.AgentChatResponse, error) {
	host, err := s.getHostByTarget(target)
	if err != nil {
		// 优化错误提示
		errMsg := err.Error()
		if strings.Contains(errMsg, "未找到主机") {
			return s.createErrorResponse(sessionID, fmt.Sprintf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", target)), nil
		}
		return s.createErrorResponse(sessionID, fmt.Sprintf("未找到主机: %v", err)), nil
	}

	result, err := s.performanceAnalysis.AnalyzeServerPressureByIP(host.SSHIP)
	if err != nil {
		return s.createErrorResponse(sessionID, fmt.Sprintf("进程分析失败: %v", err)), nil
	}

	// 生成进程分析报告
	report := s.generateProcessReport(result, intent.Entities["queryType"] == "abnormal")

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType": "process_analysis",
			"hostId":   host.ID,
		},
	}, nil
}

// handleServiceQuery 处理服务查询
func (s *AgentChatServiceImpl) handleServiceQuery(ctx context.Context, sessionID string, target string) (*model.AgentChatResponse, error) {
	host, err := s.getHostByTarget(target)
	if err != nil {
		// 优化错误提示
		errMsg := err.Error()
		if strings.Contains(errMsg, "未找到主机") {
			return s.createErrorResponse(sessionID, fmt.Sprintf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", target)), nil
		}
		return s.createErrorResponse(sessionID, fmt.Sprintf("未找到主机: %v", err)), nil
	}

	result, err := s.performanceAnalysis.AnalyzeServerPressureByIP(host.SSHIP)
	if err != nil {
		return s.createErrorResponse(sessionID, fmt.Sprintf("服务查询失败: %v", err)), nil
	}

	// 从进程列表中提取服务
	report := s.generateServiceReport(result)

	return &model.AgentChatResponse{
		SessionID: sessionID,
		Message:   report,
		Type:      "text",
		TaskInfo: map[string]interface{}{
			"taskType": "service_query",
			"hostId":   host.ID,
		},
	}, nil
}

// ============ 辅助函数 ============

// getHostByTarget 根据IP或名称获取主机
func (s *AgentChatServiceImpl) getHostByTarget(target string) (*cmdbModel.CmdbHost, error) {
	// 尝试按IP查询
	host, err := s.hostDao.GetCmdbHostByIP(target)
	if err == nil {
		return &host, nil
	}

	// 尝试按名称查询
	host, err = s.hostDao.GetCmdbHostByName(target)
	if err == nil {
		return &host, nil
	}

	return nil, fmt.Errorf("未找到主机 %s", target)
}

// generateCPUReport 生成CPU分析报告
func (s *AgentChatServiceImpl) generateCPUReport(result *model.PressureAnalysisResult, needAnalysis bool) string {
	var report strings.Builder

	report.WriteString("# CPU 使用率分析\n\n")
	report.WriteString(fmt.Sprintf("**主机**: %s (%s)\n\n", result.HostName, result.HostIP))

	// 从 CPUInfo 中获取数据
	if cpuCores, ok := result.CPUInfo["cores"].(float64); ok {
		report.WriteString(fmt.Sprintf("**CPU核心数**: %d\n", int(cpuCores)))
	}

	var cpuUsage float64
	if usage, ok := result.CPUInfo["usage"].(float64); ok {
		cpuUsage = usage
		report.WriteString(fmt.Sprintf("**CPU使用率**: %.1f%%\n\n", cpuUsage))
	}

	if needAnalysis && cpuUsage > 70 {
		report.WriteString("## 🔍 高CPU使用率分析\n\n")
		report.WriteString("CPU使用率较高，可能原因：\n\n")
		report.WriteString("### TOP进程占用分析\n\n")

		if len(result.TopProcesses) > 0 {
			report.WriteString("| 进程 | CPU% | 说明 |\n")
			report.WriteString("|------|------|------|\n")
			for i, proc := range result.TopProcesses {
				if i >= 5 {
					break
				}
				if command, ok := proc["command"].(string); ok {
					if cpuPercent, ok := proc["cpuPercent"].(float64); ok {
						report.WriteString(fmt.Sprintf("| %s | %.1f%% | - |\n", command, cpuPercent))
					}
				}
			}
		}

		report.WriteString("\n**建议**：\n")
		report.WriteString("1. 检查是否有异常进程占用大量CPU\n")
		report.WriteString("2. 考虑优化高CPU占用的应用程序\n")
		report.WriteString("3. 如果是正常业务高峰，考虑扩容\n")
	} else {
		report.WriteString("✅ CPU使用率正常\n")
	}

	return report.String()
}

// generateDiskReport 生成磁盘分析报告
func (s *AgentChatServiceImpl) generateDiskReport(result *model.PressureAnalysisResult, showDirectories bool) string {
	// 如果需要显示目录详情，调用完整的磁盘分析
	if showDirectories {
		log.Log().Infof("开始磁盘目录分析，主机: %s (%s)", result.HostName, result.HostIP)

		diskResult, err := s.performanceAnalysis.AnalyzeDiskDirectoriesByIP(result.HostIP)
		if err != nil {
			log.Log().Errorf("磁盘目录分析失败: %v", err)
			return fmt.Sprintf("# 磁盘目录分析失败\n\n错误信息: %v\n\n💡 请检查服务器SSH连接状态", err)
		}

		// 生成完整的磁盘分析报告
		return s.performanceAnalysis.GenerateDiskAnalysisReport(diskResult)
	}

	// 否则，返回基本的磁盘使用信息
	var report strings.Builder

	report.WriteString("# 磁盘空间分析\n\n")
	report.WriteString(fmt.Sprintf("**主机**: %s (%s)\n\n", result.HostName, result.HostIP))

	report.WriteString("## 磁盘使用情况\n\n")
	if diskUsage, ok := result.DiskInfo["usage"].(string); ok {
		report.WriteString("```\n")
		report.WriteString(diskUsage)
		report.WriteString("\n```\n\n")
	}

	if rootUsage, ok := result.DiskInfo["root_usage_percent"].(string); ok {
		report.WriteString(fmt.Sprintf("**根分区使用率**: %s%%\n\n", rootUsage))
	}

	report.WriteString("💡 提示：使用 \"哪个目录占用磁盘空间最大\" 查看详细目录分析\n")

	return report.String()
}

// generateProcessReport 生成进程分析报告
func (s *AgentChatServiceImpl) generateProcessReport(result *model.PressureAnalysisResult, checkAbnormal bool) string {
	var report strings.Builder

	report.WriteString("# 进程状态分析\n\n")
	report.WriteString(fmt.Sprintf("**主机**: %s (%s)\n\n", result.HostName, result.HostIP))

	if checkAbnormal {
		report.WriteString("## 异常进程检查\n\n")
		report.WriteString("⚠️ 异常进程检测功能开发中...\n\n")
		report.WriteString("💡 建议手动执行以下命令：\n\n")
		report.WriteString("```bash\n")
		report.WriteString("# 查找僵尸进程\n")
		report.WriteString("ps aux | awk '$8 == \"Z\" {print}'\n")
		report.WriteString("# 查找CPU占用超过50%的进程\n")
		report.WriteString("ps aux --sort=-%cpu | awk 'NR==1 || $3>50.0'\n")
		report.WriteString("```\n")
	} else {
		report.WriteString("## TOP 进程列表\n\n")
		if len(result.TopProcesses) > 0 {
			report.WriteString("| 用户 | PID | CPU% | MEM% | 命令 |\n")
			report.WriteString("|------|-----|------|------|------|\n")
			for i, proc := range result.TopProcesses {
				if i >= 10 {
					break
				}
				user, _ := proc["user"].(string)
				pid, _ := proc["pid"].(float64)
				cpuPercent, _ := proc["cpuPercent"].(float64)
				memPercent, _ := proc["memPercent"].(float64)
				command, _ := proc["command"].(string)

				report.WriteString(fmt.Sprintf("| %s | %d | %.1f | %.1f | %s |\n",
					user, int(pid), cpuPercent, memPercent, command))
			}
		}
	}

	return report.String()
}

// generateServiceReportFromList 根据服务列表生成报告
func (s *AgentChatServiceImpl) generateServiceReportFromList(hostName, hostIP string, serviceProcesses []string) string {
	var report strings.Builder

	report.WriteString("# 服务列表\n\n")
	report.WriteString(fmt.Sprintf("**主机**: %s (%s)\n\n", hostName, hostIP))

	if len(serviceProcesses) == 0 {
		report.WriteString("未检测到运行中的服务\n")
		return report.String()
	}

	// 服务名称映射（进程名 -> 友好名称）
	serviceMap := map[string]string{
		"mysqld":          "MySQL",
		"redis-server":    "Redis",
		"redis":           "Redis",
		"nginx":           "Nginx",
		"httpd":           "Apache",
		"apache2":         "Apache",
		"postgres":        "PostgreSQL",
		"mongod":          "MongoDB",
		"docker":          "Docker",
		"dockerd":         "Docker",
		"containerd":      "Containerd",
		"prometheus":      "Prometheus",
		"pushgateway":     "Prometheus Pushgateway",
		"grafana":         "Grafana",
		"influxd":         "InfluxDB",
		"node":            "Node.js",
		"java":            "Java",
		"python":          "Python",
		"python3":         "Python",
		"go":              "Go",
		"php-fpm":         "PHP-FPM",
		"tomcat":          "Tomcat",
		"etcd":            "Etcd",
		"consul":          "Consul",
		"kafka":           "Kafka",
		"zookeeper":       "Zookeeper",
		"elasticsearch":   "Elasticsearch",
		"kibana":          "Kibana",
		"logstash":        "Logstash",
		"rabbitmq":        "RabbitMQ",
		"memcached":       "Memcached",
		"sshd":            "SSH",
		"systemd":         "Systemd",
		"systemd-resolve": "Systemd DNS",
		"systemd-network": "Systemd Network",
		"crond":           "Cron",
		"chronyd":         "Chrony NTP",
		"rsyslogd":        "Rsyslog",
		"pmlogger":        "PCP Logger",
		"pmproxy":         "PCP Proxy",
		"pmcd":            "PCP Collector",
		"init":            "Init",
		"dodevops-agent":  "DoDevOps Agent",
		"dodevops-api":    "DoDevOps API",
	}

	// 检测到的服务（去重并映射友好名称）
	detectedServices := make(map[string]bool)

	for _, processName := range serviceProcesses {
		// 查找是否有对应的友好名称
		friendlyName := ""
		processLower := strings.ToLower(processName)

		for key, value := range serviceMap {
			if strings.Contains(processLower, key) || processLower == key {
				friendlyName = value
				break
			}
		}

		// 如果没有友好名称，使用原进程名
		if friendlyName == "" {
			friendlyName = processName
		}

		detectedServices[friendlyName] = true
	}

	if len(detectedServices) == 0 {
		report.WriteString("未检测到常见服务\n")
		return report.String()
	}

	report.WriteString("检测到以下服务：\n\n")

	// 按字母顺序排序输出
	serviceList := make([]string, 0, len(detectedServices))
	for service := range detectedServices {
		serviceList = append(serviceList, service)
	}

	// 简单排序
	for i := 0; i < len(serviceList); i++ {
		for j := i + 1; j < len(serviceList); j++ {
			if serviceList[i] > serviceList[j] {
				serviceList[i], serviceList[j] = serviceList[j], serviceList[i]
			}
		}
	}

	for _, service := range serviceList {
		report.WriteString(fmt.Sprintf("- %s\n", service))
	}

	report.WriteString(fmt.Sprintf("\n💡 共检测到 %d 个服务\n", len(detectedServices)))

	return report.String()
}

// generateServiceReport 生成服务报告（旧方法，保留用于其他调用）
func (s *AgentChatServiceImpl) generateServiceReport(result *model.PressureAnalysisResult) string {
	var report strings.Builder

	report.WriteString("# 服务列表\n\n")
	report.WriteString(fmt.Sprintf("**主机**: %s (%s)\n\n", result.HostName, result.HostIP))

	if len(result.TopProcesses) == 0 {
		report.WriteString("未检测到运行中的服务\n")
		return report.String()
	}

	// 服务识别映射
	serviceMap := map[string]string{
		"mysqld":        "MySQL",
		"redis-server":  "Redis",
		"redis":         "Redis",
		"nginx":         "Nginx",
		"httpd":         "Apache",
		"apache2":       "Apache",
		"postgres":      "PostgreSQL",
		"mongod":        "MongoDB",
		"docker":        "Docker",
		"dockerd":       "Docker",
		"containerd":    "Containerd",
		"prometheus":    "Prometheus",
		"grafana":       "Grafana",
		"node":          "Node.js",
		"java":          "Java",
		"python":        "Python",
		"python3":       "Python",
		"go":            "Go",
		"php-fpm":       "PHP-FPM",
		"tomcat":        "Tomcat",
		"etcd":          "Etcd",
		"consul":        "Consul",
		"kafka":         "Kafka",
		"zookeeper":     "Zookeeper",
		"elasticsearch": "Elasticsearch",
		"kibana":        "Kibana",
		"logstash":      "Logstash",
		"rabbitmq":      "RabbitMQ",
		"memcached":     "Memcached",
		"sshd":          "SSH",
		"systemd":       "Systemd",
		"crond":         "Cron",
		"rsyslogd":      "Rsyslog",
	}

	// 检测到的服务（去重）
	detectedServices := make(map[string]bool)

	for _, proc := range result.TopProcesses {
		command, ok := proc["command"].(string)
		if !ok {
			continue
		}
		cmdLower := strings.ToLower(command)

		for keyword, serviceName := range serviceMap {
			if strings.Contains(cmdLower, keyword) {
				detectedServices[serviceName] = true
				break
			}
		}
	}

	if len(detectedServices) == 0 {
		report.WriteString("未检测到常见服务\n\n")
		report.WriteString("💡 提示：可能服务使用了自定义名称，或者是容器化部署\n")
		return report.String()
	}

	report.WriteString("检测到以下服务：\n\n")

	// 按字母顺序排序输出
	serviceList := make([]string, 0, len(detectedServices))
	for service := range detectedServices {
		serviceList = append(serviceList, service)
	}

	// 简单排序
	for i := 0; i < len(serviceList); i++ {
		for j := i + 1; j < len(serviceList); j++ {
			if serviceList[i] > serviceList[j] {
				serviceList[i], serviceList[j] = serviceList[j], serviceList[i]
			}
		}
	}

	for _, service := range serviceList {
		report.WriteString(fmt.Sprintf("- %s\n", service))
	}

	report.WriteString(fmt.Sprintf("\n💡 共检测到 %d 个服务\n", len(detectedServices)))

	return report.String()
}
