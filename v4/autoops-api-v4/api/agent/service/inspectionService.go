// 巡检服务层
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"dodevops-api/api/agent/dao"
	"dodevops-api/api/agent/model"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	"dodevops-api/common/util"

	"github.com/gin-gonic/gin"
)

// InspectionService 巡检服务接口
type InspectionService interface {
	// 巡检任务管理
	CreateInspection(c *gin.Context, dto *model.InspectionDto) error
	GetInspectionByID(c *gin.Context, id uint) (*model.AgentInspection, error)
	UpdateInspection(c *gin.Context, dto *model.InspectionDto) error
	DeleteInspection(c *gin.Context, id uint) error
	GetInspectionList(c *gin.Context, query *model.InspectionQuery) ([]model.AgentInspection, int64, error)

	// 巡检执行
	ExecuteInspection(ctx context.Context, inspectionID uint) error
	InspectHost(ctx context.Context, inspection *model.AgentInspection, host *cmdbModel.CmdbHost) (*model.InspectionRecord, error)

	// 巡检记录
	GetRecordByID(c *gin.Context, id uint) (*model.InspectionResultSummary, error)
	GetRecordList(c *gin.Context, query *model.InspectionRecordQuery) ([]model.InspectionRecord, int64, error)
}

// InspectionServiceImpl 巡检服务实现
type InspectionServiceImpl struct {
	inspectionDao dao.InspectionDao
	recordDao     dao.InspectionRecordDao
	itemDao       dao.InspectionItemDao
	hostDao       cmdbDao.CmdbHostDao
	aiEngine      *AIDecisionEngine
}

// NewInspectionService 创建巡检服务实例
func NewInspectionService() InspectionService {
	return &InspectionServiceImpl{
		inspectionDao: dao.NewInspectionDao(),
		recordDao:     dao.NewInspectionRecordDao(),
		itemDao:       dao.NewInspectionItemDao(),
		hostDao:       cmdbDao.NewCmdbHostDao(),
		aiEngine:      GetAIDecisionEngine(),
	}
}

// CreateInspection 创建巡检任务
func (s *InspectionServiceImpl) CreateInspection(c *gin.Context, dto *model.InspectionDto) error {
	// 编码目标ID列表
	targetIDsJSON, err := json.Marshal(dto.TargetIDs)
	if err != nil {
		return fmt.Errorf("目标ID列表格式错误: %v", err)
	}

	inspection := &model.AgentInspection{
		Name:         dto.Name,
		Type:         dto.Type,
		TargetType:   dto.TargetType,
		TargetIDs:    string(targetIDsJSON),
		ScheduleType: dto.ScheduleType,
		CronExpr:     dto.CronExpr,
		Enabled:      dto.Enabled,
		Template:     dto.Template,
		Remark:       dto.Remark,
	}

	return s.inspectionDao.Create(inspection)
}

// GetInspectionByID 根据ID获取巡检任务
func (s *InspectionServiceImpl) GetInspectionByID(c *gin.Context, id uint) (*model.AgentInspection, error) {
	return s.inspectionDao.GetByID(id)
}

// UpdateInspection 更新巡检任务
func (s *InspectionServiceImpl) UpdateInspection(c *gin.Context, dto *model.InspectionDto) error {
	inspection, err := s.inspectionDao.GetByID(dto.ID)
	if err != nil {
		return fmt.Errorf("巡检任务不存在")
	}

	// 编码目标ID列表
	targetIDsJSON, err := json.Marshal(dto.TargetIDs)
	if err != nil {
		return fmt.Errorf("目标ID列表格式错误: %v", err)
	}

	inspection.Name = dto.Name
	inspection.Type = dto.Type
	inspection.TargetType = dto.TargetType
	inspection.TargetIDs = string(targetIDsJSON)
	inspection.ScheduleType = dto.ScheduleType
	inspection.CronExpr = dto.CronExpr
	inspection.Enabled = dto.Enabled
	inspection.Template = dto.Template
	inspection.Remark = dto.Remark

	return s.inspectionDao.Update(inspection)
}

// DeleteInspection 删除巡检任务
func (s *InspectionServiceImpl) DeleteInspection(c *gin.Context, id uint) error {
	_, err := s.inspectionDao.GetByID(id)
	if err != nil {
		return fmt.Errorf("巡检任务不存在")
	}

	return s.inspectionDao.Delete(id)
}

// GetInspectionList 获取巡检任务列表
func (s *InspectionServiceImpl) GetInspectionList(c *gin.Context, query *model.InspectionQuery) ([]model.AgentInspection, int64, error) {
	return s.inspectionDao.GetList(query)
}

// ExecuteInspection 执行巡检任务
func (s *InspectionServiceImpl) ExecuteInspection(ctx context.Context, inspectionID uint) error {
	// 1. 获取巡检配置
	inspection, err := s.inspectionDao.GetByID(inspectionID)
	if err != nil {
		return fmt.Errorf("巡检任务不存在: %v", err)
	}

	if inspection.Enabled != 1 {
		return fmt.Errorf("巡检任务未启用")
	}

	// 2. 解析目标主机列表
	hosts, err := s.resolveTargets(inspection.TargetType, inspection.TargetIDs)
	if err != nil {
		return fmt.Errorf("解析目标主机失败: %v", err)
	}

	if len(hosts) == 0 {
		return fmt.Errorf("未找到目标主机")
	}

	// 3. 并发执行巡检
	for _, host := range hosts {
		// 这里可以使用goroutine并发执行，暂时串行
		_, err := s.InspectHost(ctx, inspection, &host)
		if err != nil {
			// 记录错误但继续执行其他主机
			fmt.Printf("巡检主机失败 %s: %v\n", host.HostName, err)
		}
	}

	return nil
}

// InspectHost 巡检单个主机
func (s *InspectionServiceImpl) InspectHost(ctx context.Context, inspection *model.AgentInspection, host *cmdbModel.CmdbHost) (*model.InspectionRecord, error) {
	// 1. 创建巡检记录
	record := &model.InspectionRecord{
		InspectionID: inspection.ID,
		HostID:       host.ID,
		HostName:     host.HostName,
		HostIP:       host.PrivateIP,
		Status:       1, // 执行中
		StartTime:    util.HTime{Time: time.Now()},
	}
	if err := s.recordDao.Create(record); err != nil {
		return nil, fmt.Errorf("创建巡检记录失败: %v", err)
	}

	// 2. 加载巡检模板
	template, err := s.loadTemplate(inspection.Template)
	if err != nil {
		record.Status = 3 // 失败
		record.ErrorMsg = fmt.Sprintf("加载巡检模板失败: %v", err)
		s.recordDao.Update(record)
		return record, err
	}

	// 3. 执行巡检项
	items := []model.InspectionItem{}
	for _, check := range template.Checks {
		item := s.executeCheck(host, check)
		item.RecordID = record.ID
		items = append(items, item)
	}

	// 4. 批量保存巡检项
	if err := s.itemDao.BatchCreate(items); err != nil {
		record.Status = 3 // 失败
		record.ErrorMsg = fmt.Sprintf("保存巡检项失败: %v", err)
		s.recordDao.Update(record)
		return record, err
	}

	// 5. 统计结果
	s.updateRecordStatistics(record, items)

	// 6. AI分析和建议
	summary, suggestions := s.analyzeWithAI(record, items)
	record.Summary = summary
	record.Suggestions = suggestions

	// 7. 完成巡检
	record.EndTime = util.HTime{Time: time.Now()}
	record.Duration = int(record.EndTime.Sub(record.StartTime.Time).Seconds())

	if record.FailedChecks == 0 && record.WarnChecks == 0 {
		record.Status = 2 // 成功
	} else if record.FailedChecks == 0 {
		record.Status = 4 // 部分成功
	} else {
		record.Status = 4 // 部分成功
	}

	s.recordDao.Update(record)

	return record, nil
}

// loadTemplate 加载巡检模板
func (s *InspectionServiceImpl) loadTemplate(templateJSON string) (*model.InspectionTemplate, error) {
	var template model.InspectionTemplate
	if err := json.Unmarshal([]byte(templateJSON), &template); err != nil {
		return nil, err
	}
	return &template, nil
}

// executeCheck 执行单个检查项
func (s *InspectionServiceImpl) executeCheck(host *cmdbModel.CmdbHost, check model.InspectionCheck) model.InspectionItem {
	item := model.InspectionItem{
		Category: check.Category,
		Name:     check.Name,
		Expected: check.Expected,
		Message:  check.Message,
	}

	// 根据检查类型执行不同的检查逻辑
	switch check.Type {
	case "command":
		item = s.executeCommandCheck(host, check)
	case "metric":
		item = s.executeMetricCheck(host, check)
	case "file":
		item = s.executeFileCheck(host, check)
	case "process":
		item = s.executeProcessCheck(host, check)
	default:
		item.Status = 2 // 失败
		item.Message = "不支持的检查类型: " + check.Type
	}

	return item
}

// executeCommandCheck 执行命令检查
func (s *InspectionServiceImpl) executeCommandCheck(host *cmdbModel.CmdbHost, check model.InspectionCheck) model.InspectionItem {
	item := model.InspectionItem{
		Category:   check.Category,
		Name:       check.Name,
		Expected:   check.Expected,
		Suggestion: check.Suggestion,
	}

	// 执行SSH命令
	// TODO: 需要实现SSH命令执行功能
	output := ""
	err := fmt.Errorf("SSH命令执行功能未实现")
	_ = output
	if err != nil {
		item.Status = 2 // 失败
		item.Actual = ""
		item.Message = fmt.Sprintf("命令执行失败: %v", err)
		return item
	}

	item.Actual = strings.TrimSpace(output)

	// 判断结果
	if s.checkResult(item.Actual, check.Expected, check.Operator) {
		item.Status = 1 // 通过
		item.Message = "检查通过"
	} else {
		if check.Severity == "fail" {
			item.Status = 2 // 失败
		} else {
			item.Status = 3 // 警告
		}
		item.Message = check.Message
	}

	return item
}

// executeMetricCheck 执行指标检查
func (s *InspectionServiceImpl) executeMetricCheck(host *cmdbModel.CmdbHost, check model.InspectionCheck) model.InspectionItem {
	item := model.InspectionItem{
		Category:   check.Category,
		Name:       check.Name,
		Expected:   check.Expected,
		Suggestion: check.Suggestion,
	}

	// TODO: 从监控系统获取指标
	// 这里暂时返回模拟数据
	item.Status = 1 // 通过
	item.Actual = "正常"
	item.Message = "指标检查通过（暂未实现实际检查）"

	return item
}

// executeFileCheck 执行文件检查
func (s *InspectionServiceImpl) executeFileCheck(host *cmdbModel.CmdbHost, check model.InspectionCheck) model.InspectionItem {
	item := model.InspectionItem{
		Category:   check.Category,
		Name:       check.Name,
		Expected:   check.Expected,
		Suggestion: check.Suggestion,
	}

	// 检查文件是否存在
	command := fmt.Sprintf("test -f %s && echo 'exists' || echo 'not_exists'", check.FilePath)
	// TODO: 需要实现SSH命令执行功能
	output := ""
	err := fmt.Errorf("SSH命令执行功能未实现")
	_ = command
	_ = output
	if err != nil {
		item.Status = 2 // 失败
		item.Message = fmt.Sprintf("文件检查失败: %v", err)
		return item
	}

	item.Actual = strings.TrimSpace(output)

	if item.Actual == "exists" {
		item.Status = 1 // 通过
		item.Message = "文件存在"
	} else {
		item.Status = 2 // 失败
		item.Message = "文件不存在"
	}

	return item
}

// executeProcessCheck 执行进程检查
func (s *InspectionServiceImpl) executeProcessCheck(host *cmdbModel.CmdbHost, check model.InspectionCheck) model.InspectionItem {
	item := model.InspectionItem{
		Category:   check.Category,
		Name:       check.Name,
		Expected:   check.Expected,
		Suggestion: check.Suggestion,
	}

	// 检查进程是否运行
	command := fmt.Sprintf("ps aux | grep -v grep | grep '%s' | wc -l", check.Process)
	// TODO: 需要实现SSH命令执行功能
	output := ""
	err := fmt.Errorf("SSH命令执行功能未实现")
	_ = command
	_ = output
	if err != nil {
		item.Status = 2 // 失败
		item.Message = fmt.Sprintf("进程检查失败: %v", err)
		return item
	}

	item.Actual = strings.TrimSpace(output)

	if item.Actual != "0" {
		item.Status = 1 // 通过
		item.Message = "进程运行中"
	} else {
		item.Status = 2 // 失败
		item.Message = "进程未运行"
	}

	return item
}

// checkResult 检查结果
func (s *InspectionServiceImpl) checkResult(actual string, expected string, operator string) bool {
	switch operator {
	case "eq":
		return actual == expected
	case "ne":
		return actual != expected
	case "contains":
		return strings.Contains(actual, expected)
	case "not_contains":
		return !strings.Contains(actual, expected)
	default:
		return actual == expected
	}
}

// updateRecordStatistics 更新巡检记录统计
func (s *InspectionServiceImpl) updateRecordStatistics(record *model.InspectionRecord, items []model.InspectionItem) {
	record.TotalChecks = len(items)
	record.PassedChecks = 0
	record.FailedChecks = 0
	record.WarnChecks = 0
	maxRiskLevel := 0

	for _, item := range items {
		switch item.Status {
		case 1:
			record.PassedChecks++
		case 2:
			record.FailedChecks++
			if maxRiskLevel < 2 {
				maxRiskLevel = 2
			}
		case 3:
			record.WarnChecks++
			if maxRiskLevel < 1 {
				maxRiskLevel = 1
			}
		}
	}

	// 根据失败和警告的比例确定风险等级
	if record.FailedChecks > record.TotalChecks/2 {
		record.RiskLevel = 3 // 高危
	} else if record.FailedChecks > 0 {
		record.RiskLevel = 2 // 中危
	} else if record.WarnChecks > 0 {
		record.RiskLevel = 1 // 低危
	} else {
		record.RiskLevel = 0 // 正常
	}
}

// analyzeWithAI AI分析和建议
func (s *InspectionServiceImpl) analyzeWithAI(record *model.InspectionRecord, items []model.InspectionItem) (string, string) {
	// 构建检查结果摘要
	checkResults := ""
	for _, item := range items {
		statusStr := "通过"
		if item.Status == 2 {
			statusStr = "失败"
		} else if item.Status == 3 {
			statusStr = "警告"
		}
		checkResults += fmt.Sprintf("- [%s] %s: %s\n", statusStr, item.Name, item.Message)
	}

	// 调用AI分析
	data := map[string]interface{}{
		"hostName":        record.HostName,
		"hostIp":          record.HostIP,
		"inspectionTime":  record.StartTime.Format("2006-01-02 15:04:05"),
		"totalChecks":     record.TotalChecks,
		"passedChecks":    record.PassedChecks,
		"failedChecks":    record.FailedChecks,
		"warnChecks":      record.WarnChecks,
		"checkResults":    checkResults,
	}

	response, err := s.aiEngine.AnalyzeWithTemplate("inspection_summary", data)
	if err != nil {
		// AI分析失败，返回默认摘要
		return s.generateDefaultSummary(record), s.generateDefaultSuggestions(items)
	}

	// 简单分割summary和suggestions
	parts := strings.Split(response, "关键建议")
	summary := strings.TrimSpace(parts[0])
	suggestions := ""
	if len(parts) > 1 {
		suggestions = strings.TrimSpace(parts[1])
	}

	return summary, suggestions
}

// generateDefaultSummary 生成默认摘要
func (s *InspectionServiceImpl) generateDefaultSummary(record *model.InspectionRecord) string {
	return fmt.Sprintf("主机 %s 巡检完成，总检查项 %d，通过 %d，失败 %d，警告 %d",
		record.HostName, record.TotalChecks, record.PassedChecks, record.FailedChecks, record.WarnChecks)
}

// generateDefaultSuggestions 生成默认建议
func (s *InspectionServiceImpl) generateDefaultSuggestions(items []model.InspectionItem) string {
	suggestions := ""
	for _, item := range items {
		if item.Status != 1 && item.Suggestion != "" {
			suggestions += fmt.Sprintf("- %s: %s\n", item.Name, item.Suggestion)
		}
	}
	if suggestions == "" {
		suggestions = "系统运行正常，暂无特殊建议"
	}
	return suggestions
}

// resolveTargets 解析目标主机列表
func (s *InspectionServiceImpl) resolveTargets(targetType string, targetIDsJSON string) ([]cmdbModel.CmdbHost, error) {
	var targetIDs []uint
	if err := json.Unmarshal([]byte(targetIDsJSON), &targetIDs); err != nil {
		return nil, err
	}

	var hosts []cmdbModel.CmdbHost

	switch targetType {
	case "host":
		// 直接根据主机ID列表获取
		for _, id := range targetIDs {
			host, err := s.hostDao.GetCmdbHostById(id)
			if err == nil {
				hosts = append(hosts, host)
			}
		}
	case "group":
		// 根据分组ID获取主机列表
		for _, groupID := range targetIDs {
			// TODO: 实现根据分组ID获取主机列表
			_ = groupID
		}
	default:
		return nil, fmt.Errorf("不支持的目标类型: %s", targetType)
	}

	return hosts, nil
}

// GetRecordByID 根据ID获取巡检记录详情
func (s *InspectionServiceImpl) GetRecordByID(c *gin.Context, id uint) (*model.InspectionResultSummary, error) {
	record, err := s.recordDao.GetByID(id)
	if err != nil {
		return nil, err
	}

	items, err := s.itemDao.GetByRecordID(id)
	if err != nil {
		return nil, err
	}

	// 构建分类统计
	categoryStats := make(map[string]model.CategoryStat)
	for _, item := range items {
		stat, exists := categoryStats[item.Category]
		if !exists {
			stat = model.CategoryStat{
				Category: item.Category,
			}
		}
		stat.Total++
		switch item.Status {
		case 1:
			stat.Passed++
		case 2:
			stat.Failed++
		case 3:
			stat.Warned++
		}
		categoryStats[item.Category] = stat
	}

	summary := &model.InspectionResultSummary{
		RecordID:      record.ID,
		HostID:        record.HostID,
		HostName:      record.HostName,
		HostIP:        record.HostIP,
		Status:        record.Status,
		RiskLevel:     record.RiskLevel,
		TotalChecks:   record.TotalChecks,
		PassedChecks:  record.PassedChecks,
		FailedChecks:  record.FailedChecks,
		WarnChecks:    record.WarnChecks,
		Duration:      record.Duration,
		Items:         items,
		Summary:       record.Summary,
		Suggestions:   record.Suggestions,
		CategoryStats: categoryStats,
	}

	return summary, nil
}

// GetRecordList 获取巡检记录列表
func (s *InspectionServiceImpl) GetRecordList(c *gin.Context, query *model.InspectionRecordQuery) ([]model.InspectionRecord, int64, error) {
	return s.recordDao.GetList(query)
}
