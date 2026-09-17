package service

import (
	"fmt"
	"time"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	systemDao "dodevops-api/api/system/dao"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type AlertRuleService interface {
	// 规则管理
	CreateRule(c *gin.Context, rule *model.AlertRule)
	UpdateRule(c *gin.Context, rule *model.AlertRule)
	DeleteRule(c *gin.Context, id uint)
	GetRule(c *gin.Context, id uint)
	GetRuleList(c *gin.Context, status *int, groupID *uint, page, pageSize int)
	UpdateRuleStatus(c *gin.Context, id uint, status int)

	// 告警记录
	GetAlertRecords(c *gin.Context, status, severity, ip string, ruleID, groupID *uint, startTime, endTime *time.Time, page, pageSize int)
	GetAlertRecord(c *gin.Context, id uint)
	GetAlertStats(c *gin.Context)
	ProcessAlertRecord(c *gin.Context, id uint, action string, assigneeID *uint, handleNote string)
	HandleAlertRecord(c *gin.Context, id uint, handler, handleNote string)
	AcknowledgeAlert(c *gin.Context, id uint, handler string)
	DeleteAlertRecord(c *gin.Context, id uint)
	BatchDeleteAlertRecords(c *gin.Context, ids []uint)
	ClearAllAlertRecords(c *gin.Context)
}

type AlertRuleServiceImpl struct {
	ruleDao   dao.AlertRuleDao
	recordDao dao.AlertRecordDao
}

func NewAlertRuleService() AlertRuleService {
	return &AlertRuleServiceImpl{
		ruleDao:   dao.NewAlertRuleDao(),
		recordDao: dao.NewAlertRecordDao(),
	}
}

// CreateRule 创建告警规则
func (s *AlertRuleServiceImpl) CreateRule(c *gin.Context, rule *model.AlertRule) {
	// 验证必填字段
	if rule.Name == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "规则名称不能为空")
		return
	}

	if rule.PromQL == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "PromQL表达式不能为空")
		return
	}

	// 设置默认值
	if rule.CronPattern == "" {
		rule.CronPattern = "@every 30s"
	}

	if rule.Duration == 0 {
		rule.Duration = 60
	}

	if rule.Severity == 0 {
		rule.Severity = 2 // 默认警告级别
	}

	if rule.SilencePeriod == 0 {
		rule.SilencePeriod = 300 // 默认5分钟静默期
	}

	if rule.RecoverDuration == 0 {
		rule.RecoverDuration = 60 // 默认1分钟恢复判定时间
	}

	// 创建规则
	if err := s.ruleDao.Create(rule); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建规则失败: %v", err))
		return
	}

	result.Success(c, rule)
}

// UpdateRule 更新告警规则
func (s *AlertRuleServiceImpl) UpdateRule(c *gin.Context, rule *model.AlertRule) {
	// 检查规则是否存在
	existing, err := s.ruleDao.GetByID(rule.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("规则不存在: %v", err))
		return
	}

	// 验证必填字段
	if rule.Name == "" {
		rule.Name = existing.Name
	}

	if rule.PromQL == "" {
		rule.PromQL = existing.PromQL
	}

	// 更新规则
	if err := s.ruleDao.Update(rule); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新规则失败: %v", err))
		return
	}

	result.Success(c, rule)
}

// DeleteRule 删除告警规则
func (s *AlertRuleServiceImpl) DeleteRule(c *gin.Context, id uint) {
	// 检查规则是否存在
	_, err := s.ruleDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("规则不存在: %v", err))
		return
	}

	if err := s.ruleDao.Delete(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("删除规则失败: %v", err))
		return
	}

	result.Success(c, "删除成功")
}

// GetRule 获取告警规则
func (s *AlertRuleServiceImpl) GetRule(c *gin.Context, id uint) {
	rule, err := s.ruleDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询规则失败: %v", err))
		return
	}

	result.Success(c, rule)
}

// GetRuleList 获取告警规则列表
func (s *AlertRuleServiceImpl) GetRuleList(c *gin.Context, status *int, groupID *uint, page, pageSize int) {
	// 设置默认分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	rules, total, err := s.ruleDao.GetList(status, groupID, page, pageSize)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询规则列表失败: %v", err))
		return
	}

	result.Success(c, gin.H{
		"list":     rules,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// UpdateRuleStatus 更新规则状态
func (s *AlertRuleServiceImpl) UpdateRuleStatus(c *gin.Context, id uint, status int) {
	// 检查规则是否存在
	_, err := s.ruleDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("规则不存在: %v", err))
		return
	}

	if err := s.ruleDao.UpdateStatus(id, status); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新规则状态失败: %v", err))
		return
	}

	statusText := "禁用"
	if status == 1 {
		statusText = "启用"
	}

	result.Success(c, fmt.Sprintf("规则已%s", statusText))
}

// GetAlertRecords 获取告警记录列表
func (s *AlertRuleServiceImpl) GetAlertRecords(c *gin.Context, status, severity, ip string, ruleID, groupID *uint, startTime, endTime *time.Time, page, pageSize int) {
	// 设置默认分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	records, total, err := s.recordDao.GetList(status, severity, ip, ruleID, groupID, startTime, endTime, page, pageSize)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询告警记录失败: %v", err))
		return
	}

	result.Success(c, gin.H{
		"list":     records,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetAlertRecord 获取告警记录详情
func (s *AlertRuleServiceImpl) GetAlertRecord(c *gin.Context, id uint) {
	record, err := s.recordDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询告警记录失败: %v", err))
		return
	}

	result.Success(c, record)
}

// GetAlertStats 获取告警统计
func (s *AlertRuleServiceImpl) GetAlertStats(c *gin.Context) {
	// 统计触发中的告警
	firingRecords, _, _ := s.recordDao.GetList("firing", "", "", nil, nil, nil, nil, 1, 10000)

	// 按严重级别统计
	criticalCount := 0
	warningCount := 0
	infoCount := 0

	for _, record := range firingRecords {
		switch record.Severity {
		case 1:
			criticalCount++
		case 2:
			warningCount++
		case 3:
			infoCount++
		}
	}

	// 统计今日告警数
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayRecords, todayTotal, _ := s.recordDao.GetList("", "", "", nil, nil, &startOfDay, nil, 1, 10000)

	// 统计今日恢复数
	todayResolvedCount := 0
	for _, record := range todayRecords {
		if record.Status == "resolved" {
			todayResolvedCount++
		}
	}

	result.Success(c, gin.H{
		"firing_total":   len(firingRecords),
		"critical_count": criticalCount,
		"warning_count":  warningCount,
		"info_count":     infoCount,
		"today_total":    todayTotal,
		"today_resolved": todayResolvedCount,
		"today_firing":   int(todayTotal) - todayResolvedCount,
	})
}

func appendAlertHandleNote(existing, entry string) string {
	if existing == "" {
		return entry
	}
	return existing + "\n" + entry
}

func formatAlertActionLog(operatorName, actionLabel, detail string) string {
	ts := time.Now().Format("2006-01-02 15:04:05")
	if detail == "" {
		return fmt.Sprintf("[%s] %s%s", ts, operatorName, actionLabel)
	}
	return fmt.Sprintf("[%s] %s%s: %s", ts, operatorName, actionLabel, detail)
}

func (s *AlertRuleServiceImpl) processAlertRecord(c *gin.Context, id uint, action string, assigneeID *uint, handleNote string, fallbackHandler string) {
	record, err := s.recordDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("告警记录不存在: %v", err))
		return
	}

	if record.Status != "firing" && record.Status != "acknowledged" {
		result.Failed(c, int(result.ApiCode.FAILED), "只能操作触发中或已认领的告警")
		return
	}

	currentUser, err := jwt.GetAdmin(c)
	if err != nil || currentUser == nil {
		result.Failed(c, int(result.ApiCode.FAILED), "无法获取当前登录用户")
		return
	}

	targetHandler := currentUser.Username
	actionLog := ""
	successMessage := "告警处理成功"

	switch action {
	case "acknowledge":
		if record.Status != "firing" {
			result.Failed(c, int(result.ApiCode.FAILED), "只能认领触发中的告警")
			return
		}
		actionLog = formatAlertActionLog(currentUser.Username, "认领了告警", handleNote)
		successMessage = "告警已认领"
	case "assign":
		if assigneeID == nil || *assigneeID == 0 {
			result.Failed(c, int(result.ApiCode.FAILED), "转派时必须选择目标用户")
			return
		}
		assignee := systemDao.GetSysAdminInfo(int(*assigneeID))
		if assignee.ID == 0 {
			result.Failed(c, int(result.ApiCode.FAILED), "目标用户不存在")
			return
		}
		if assignee.Status != 1 {
			result.Failed(c, int(result.ApiCode.FAILED), "目标用户未启用，无法转派")
			return
		}
		targetHandler = assignee.Username
		actionLog = formatAlertActionLog(currentUser.Username, "转派给了"+assignee.Username, handleNote)
		successMessage = "告警已转派"
	case "handle":
		if handleNote == "" {
			result.Failed(c, int(result.ApiCode.FAILED), "处理时请填写处理步骤")
			return
		}
		if assigneeID != nil && *assigneeID > 0 {
			assignee := systemDao.GetSysAdminInfo(int(*assigneeID))
			if assignee.ID == 0 {
				result.Failed(c, int(result.ApiCode.FAILED), "目标用户不存在")
				return
			}
			if assignee.Status != 1 {
				result.Failed(c, int(result.ApiCode.FAILED), "目标用户未启用，无法处理")
				return
			}
			targetHandler = assignee.Username
			actionLog = formatAlertActionLog(currentUser.Username, "更新了处理记录并指派给"+assignee.Username, handleNote)
		} else if record.Handler != "" {
			targetHandler = record.Handler
			actionLog = formatAlertActionLog(currentUser.Username, "更新了处理记录", handleNote)
		} else if fallbackHandler != "" {
			targetHandler = fallbackHandler
			actionLog = formatAlertActionLog(currentUser.Username, "更新了处理记录", handleNote)
		} else {
			targetHandler = currentUser.Username
			actionLog = formatAlertActionLog(currentUser.Username, "认领并更新了处理记录", handleNote)
		}
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的处理动作")
		return
	}

	newNote := appendAlertHandleNote(record.HandleNote, actionLog)
	if err := s.recordDao.UpdateHandler(id, targetHandler, newNote); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新告警记录失败: %v", err))
		return
	}

	result.Success(c, gin.H{
		"id":         id,
		"action":     action,
		"handler":    targetHandler,
		"handleNote": newNote,
		"message":    successMessage,
	})
}

// HandleAlertRecord 处理告警记录
func (s *AlertRuleServiceImpl) HandleAlertRecord(c *gin.Context, id uint, handler, handleNote string) {
	s.processAlertRecord(c, id, "handle", nil, handleNote, handler)
}

// AcknowledgeAlert 认领告警
func (s *AlertRuleServiceImpl) AcknowledgeAlert(c *gin.Context, id uint, handler string) {
	s.processAlertRecord(c, id, "acknowledge", nil, "", handler)
}

func (s *AlertRuleServiceImpl) ProcessAlertRecord(c *gin.Context, id uint, action string, assigneeID *uint, handleNote string) {
	s.processAlertRecord(c, id, action, assigneeID, handleNote, "")
}

// DeleteAlertRecord 删除单条告警记录
func (s *AlertRuleServiceImpl) DeleteAlertRecord(c *gin.Context, id uint) {
	// 检查记录是否存在
	_, err := s.recordDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("告警记录不存在: %v", err))
		return
	}

	// 删除记录
	if err := s.recordDao.Delete(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("删除告警记录失败: %v", err))
		return
	}

	result.Success(c, "删除成功")
}

// BatchDeleteAlertRecords 批量删除告警记录
func (s *AlertRuleServiceImpl) BatchDeleteAlertRecords(c *gin.Context, ids []uint) {
	if len(ids) == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "请选择要删除的告警记录")
		return
	}

	// 批量删除
	if err := s.recordDao.BatchDelete(ids); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("批量删除告警记录失败: %v", err))
		return
	}

	result.Success(c, fmt.Sprintf("成功删除 %d 条告警记录", len(ids)))
}

// ClearAllAlertRecords 清空所有告警记录
func (s *AlertRuleServiceImpl) ClearAllAlertRecords(c *gin.Context) {
	// 清空所有记录
	if err := s.recordDao.DeleteAll(); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("清空告警记录失败: %v", err))
		return
	}

	result.Success(c, "已清空所有告警记录")
}
