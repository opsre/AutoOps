package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/config"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type AlertPolicyService interface {
	CreatePolicy(c *gin.Context, dto *model.CreateAlertPolicyDto)
	UpdatePolicy(c *gin.Context, dto *model.UpdateAlertPolicyDto)
	DeletePolicy(c *gin.Context, id uint)
	GetPolicy(c *gin.Context, id uint)
	GetPolicyList(c *gin.Context, groupID *uint, status *int, keyword string, page, pageSize int)
	UpdatePolicyStatus(c *gin.Context, id uint, status int)
	PreviewPrometheusQuery(c *gin.Context, dto *model.PrometheusQueryPreviewDto)

	// 策略规则管理
	CreatePolicyRule(c *gin.Context, dto *model.CreatePolicyRuleDto)
	GetPolicyRule(c *gin.Context, id uint)
	UpdatePolicyRule(c *gin.Context, dto *model.UpdatePolicyRuleDto)
	DeletePolicyRule(c *gin.Context, id uint)
	GetPolicyRules(c *gin.Context, policyID uint, category string)
	UpdatePolicyRuleStatus(c *gin.Context, id uint, status int)
	PreviewPolicyRule(c *gin.Context, ruleID uint)
}

type AlertPolicyServiceImpl struct {
	policyDao     dao.AlertPolicyDao
	templateDao   dao.AlertRuleTemplateDao
	policyRuleDao dao.AlertPolicyRuleDao
}

func NewAlertPolicyService() AlertPolicyService {
	return &AlertPolicyServiceImpl{
		policyDao:     dao.NewAlertPolicyDao(),
		templateDao:   dao.NewAlertRuleTemplateDao(),
		policyRuleDao: dao.NewAlertPolicyRuleDao(),
	}
}

// CreatePolicy 创建告警策略
func (s *AlertPolicyServiceImpl) CreatePolicy(c *gin.Context, dto *model.CreateAlertPolicyDto) {
	// 创建策略
	policy := &model.AlertPolicy{
		Name:           dto.Name,
		Description:    dto.Description,
		GroupID:        dto.GroupID,
		GroupName:      dto.GroupName,
		RobotIDs:       dto.RobotIDs,
		SilencePeriod:  dto.SilencePeriod,
		RepeatInterval: dto.RepeatInterval,
		MaxNotifyCount: dto.MaxNotifyCount,
		SendRecovery:   dto.SendRecovery,
		AppendTags:     dto.AppendTags,
		Status:         dto.Status,
	}

	// 设置默认值
	if policy.SilencePeriod == 0 {
		policy.SilencePeriod = 300
	}
	if policy.RepeatInterval == 0 {
		policy.RepeatInterval = 600
	}
	if policy.MaxNotifyCount == 0 {
		policy.MaxNotifyCount = 10
	}

	// 创建策略
	if err := s.policyDao.Create(policy); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建告警策略失败: %v", err))
		return
	}

	result.Success(c, policy)
}

// UpdatePolicy 更新告警策略
func (s *AlertPolicyServiceImpl) UpdatePolicy(c *gin.Context, dto *model.UpdateAlertPolicyDto) {
	// 检查策略是否存在
	_, err := s.policyDao.GetByID(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("告警策略不存在: %v", err))
		return
	}

	// 更新策略
	policy := &model.AlertPolicy{
		ID:             dto.ID,
		Name:           dto.Name,
		Description:    dto.Description,
		GroupID:        dto.GroupID,
		GroupName:      dto.GroupName,
		RobotIDs:       dto.RobotIDs,
		SilencePeriod:  dto.SilencePeriod,
		RepeatInterval: dto.RepeatInterval,
		MaxNotifyCount: dto.MaxNotifyCount,
		SendRecovery:   dto.SendRecovery,
		AppendTags:     dto.AppendTags,
		Status:         dto.Status,
	}

	if err := s.policyDao.Update(policy); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新告警策略失败: %v", err))
		return
	}

	result.Success(c, policy)
}

// DeletePolicy 删除告警策略
func (s *AlertPolicyServiceImpl) DeletePolicy(c *gin.Context, id uint) {
	// 检查策略是否存在
	_, err := s.policyDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("告警策略不存在: %v", err))
		return
	}

	// 删除策略实例规则
	s.policyRuleDao.DeleteByPolicyID(id)

	// 删除策略
	if err := s.policyDao.Delete(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("删除告警策略失败: %v", err))
		return
	}

	result.Success(c, "删除成功")
}

// GetPolicy 获取告警策略详情
func (s *AlertPolicyServiceImpl) GetPolicy(c *gin.Context, id uint) {
	policy, err := s.policyDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询告警策略失败: %v", err))
		return
	}

	result.Success(c, policy)
}

// GetPolicyList 获取告警策略列表
func (s *AlertPolicyServiceImpl) GetPolicyList(c *gin.Context, groupID *uint, status *int, keyword string, page, pageSize int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	policies, total, err := s.policyDao.GetList(groupID, status, keyword, page, pageSize)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询告警策略列表失败: %v", err))
		return
	}

	result.Success(c, gin.H{
		"list":     policies,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// UpdatePolicyStatus 更新告警策略状态
func (s *AlertPolicyServiceImpl) UpdatePolicyStatus(c *gin.Context, id uint, status int) {
	// 检查策略是否存在
	_, err := s.policyDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("告警策略不存在: %v", err))
		return
	}

	if err := s.policyDao.UpdateStatus(id, status); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新告警策略状态失败: %v", err))
		return
	}

	statusText := "禁用"
	if status == 1 {
		statusText = "启用"
	}

	result.Success(c, fmt.Sprintf("告警策略已%s", statusText))
}

// PreviewPrometheusQuery 预览Prometheus查询
func (s *AlertPolicyServiceImpl) PreviewPrometheusQuery(c *gin.Context, dto *model.PrometheusQueryPreviewDto) {
	// 从配置中获取Prometheus地址
	prometheusURL := config.Config.Monitor.Prometheus.URL
	if prometheusURL == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "Prometheus地址未配置")
		return
	}

	// 构建查询URL,对PromQL进行URL编码
	queryURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, url.QueryEscape(dto.PromQL))
	if dto.Timestamp > 0 {
		queryURL = fmt.Sprintf("%s&time=%d", queryURL, dto.Timestamp)
	}

	// 发送HTTP请求
	resp, err := http.Get(queryURL)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询Prometheus失败: %v", err))
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("读取响应失败: %v", err))
		return
	}

	// 打印原始响应,便于调试
	bodyStr := string(body)

	// 解析响应
	var prometheusResp struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Value  []any             `json:"value"`
			} `json:"result"`
		} `json:"data"`
		Error     string `json:"error"`
		ErrorType string `json:"errorType"`
	}

	if err := json.Unmarshal(body, &prometheusResp); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("解析响应失败: %v, 原始响应: %s", err, bodyStr))
		return
	}

	if prometheusResp.Status != "success" {
		errMsg := prometheusResp.Error
		if errMsg == "" {
			errMsg = "未知错误"
		}
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败[%s]: %s", prometheusResp.ErrorType, errMsg))
		return
	}

	result.Success(c, gin.H{
		"resultType": prometheusResp.Data.ResultType,
		"results":    prometheusResp.Data.Result,
		"count":      len(prometheusResp.Data.Result),
	})
}

// CreatePolicyRule 创建策略规则(从模板导入)
func (s *AlertPolicyServiceImpl) CreatePolicyRule(c *gin.Context, dto *model.CreatePolicyRuleDto) {
	// 检查策略是否存在
	policy, err := s.policyDao.GetByID(dto.PolicyID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("告警策略不存在: %v", err))
		return
	}

	// 检查模板是否存在
	template, err := s.templateDao.GetByID(dto.TemplateID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("规则模板不存在: %v", err))
		return
	}

	// 构建PromQL (只查询指标数据，不包含比较条件，比较逻辑在代码中实现)
	promql := fmt.Sprintf("%s{group_id=\"%d\"}", template.MetricKey, policy.GroupID)

	// 生成规则名称
	ruleName := fmt.Sprintf("%s-%s", policy.Name, template.Name)

	// 使用用户提供的描述,如果没有则使用默认描述
	description := dto.Description
	if description == "" {
		description = fmt.Sprintf("当%s %s %.2f%s时触发告警", template.Name, dto.Operator, dto.Threshold, template.Unit)
	}

	// 创建策略规则
	rule := &model.AlertPolicyRule{
		PolicyID:        dto.PolicyID,
		TemplateID:      dto.TemplateID,
		RuleName:        ruleName,
		Description:     description,
		Category:        template.Category,
		MetricKey:       template.MetricKey,
		Unit:            template.Unit,
		Operator:        dto.Operator,
		Threshold:       dto.Threshold,
		PromQL:          promql,
		Severity:        dto.Severity,
		Duration:        dto.Duration,
		ContinuousTimes: dto.ContinuousTimes,
		RecoverDuration: dto.RecoverDuration,
		Status:          dto.Status,
	}

	// 设置默认值
	if rule.Severity == 0 {
		rule.Severity = 2 // warning
	}
	if rule.Duration == 0 {
		rule.Duration = 60
	}
	if rule.ContinuousTimes == 0 {
		rule.ContinuousTimes = 1
	}
	if rule.RecoverDuration == 0 {
		rule.RecoverDuration = 60
	}
	if rule.Status == 0 {
		rule.Status = 1
	}

	if err := s.policyRuleDao.Create(rule); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建策略规则失败: %v", err))
		return
	}

	result.Success(c, rule)
}

// UpdatePolicyRule 更新策略规则
func (s *AlertPolicyServiceImpl) UpdatePolicyRule(c *gin.Context, dto *model.UpdatePolicyRuleDto) {
	// 检查规则是否存在
	existing, err := s.policyRuleDao.GetByID(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("策略规则不存在: %v", err))
		return
	}

	// 获取策略信息(用于重新生成PromQL)
	policy, err := s.policyDao.GetByID(existing.PolicyID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("告警策略不存在: %v", err))
		return
	}

	// 更新PromQL (只查询指标数据，不包含比较条件)
	promql := existing.PromQL
	if dto.Operator != "" && dto.Threshold != 0 {
		promql = fmt.Sprintf("%s{group_id=\"%d\"}", existing.MetricKey, policy.GroupID)
	}

	rule := &model.AlertPolicyRule{
		ID:              dto.ID,
		RuleName:        dto.RuleName,
		Description:     dto.Description,
		Operator:        dto.Operator,
		Threshold:       dto.Threshold,
		PromQL:          promql,
		Severity:        dto.Severity,
		Duration:        dto.Duration,
		ContinuousTimes: dto.ContinuousTimes,
		RecoverDuration: dto.RecoverDuration,
		Status:          dto.Status,
	}

	if err := s.policyRuleDao.Update(rule); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新策略规则失败: %v", err))
		return
	}

	result.Success(c, rule)
}

// GetPolicyRule 获取单个策略规则详情
func (s *AlertPolicyServiceImpl) GetPolicyRule(c *gin.Context, id uint) {
	rule, err := s.policyRuleDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("策略规则不存在: %v", err))
		return
	}

	result.Success(c, rule)
}

// DeletePolicyRule 删除策略规则
func (s *AlertPolicyServiceImpl) DeletePolicyRule(c *gin.Context, id uint) {
	// 检查规则是否存在
	_, err := s.policyRuleDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("策略规则不存在: %v", err))
		return
	}

	if err := s.policyRuleDao.Delete(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("删除策略规则失败: %v", err))
		return
	}

	result.Success(c, "删除成功")
}

// GetPolicyRules 获取策略下的规则列表
func (s *AlertPolicyServiceImpl) GetPolicyRules(c *gin.Context, policyID uint, category string) {
	// 检查策略是否存在
	policy, err := s.policyDao.GetByID(policyID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("告警策略不存在: %v", err))
		return
	}

	rules, err := s.policyRuleDao.GetByPolicyID(policyID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询策略规则失败: %v", err))
		return
	}

	// 按类型分组
	groupByCategory := make(map[string][]*model.AlertPolicyRule)
	categoryLabels := map[string]string{
		"cpu":        "CPU",
		"memory":     "内存",
		"load":       "负载",
		"disk":       "磁盘",
		"io":         "IO",
		"network":    "网络",
		"connection": "连接数",
		"service":    "服务",
		"custom":     "自定义",
	}

	// 如果指定了category,则过滤规则
	filteredRules := rules
	if category != "" {
		filteredRules = make([]*model.AlertPolicyRule, 0)
		for _, rule := range rules {
			ruleCategory := rule.Category
			if ruleCategory == "" {
				ruleCategory = "custom"
			}
			if ruleCategory == category {
				filteredRules = append(filteredRules, rule)
			}
		}
	}

	// 对过滤后的规则进行分组
	for _, rule := range filteredRules {
		ruleCategory := rule.Category
		if ruleCategory == "" {
			ruleCategory = "custom"
		}
		groupByCategory[ruleCategory] = append(groupByCategory[ruleCategory], rule)
	}

	// 统计每个分类的规则数量
	categoryStats := make(map[string]int)
	for cat, ruleList := range groupByCategory {
		categoryStats[cat] = len(ruleList)
	}

	result.Success(c, gin.H{
		"policyInfo": gin.H{
			"id":        policy.ID,
			"name":      policy.Name,
			"groupId":   policy.GroupID,
			"groupName": policy.GroupName,
			"status":    policy.Status,
			"ruleCount": len(rules), // 总规则数不变
		},
		"rules":            filteredRules, // 返回过滤后的规则
		"groupByCategory":  groupByCategory,
		"categoryStats":    categoryStats,
		"categoryLabels":   categoryLabels,
		"totalCount":       len(filteredRules), // 过滤后的数量
		"filteredCategory": category,           // 当前过滤的类型
	})
}

// UpdatePolicyRuleStatus 更新策略规则状态
func (s *AlertPolicyServiceImpl) UpdatePolicyRuleStatus(c *gin.Context, id uint, status int) {
	// 检查规则是否存在
	_, err := s.policyRuleDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("策略规则不存在: %v", err))
		return
	}

	if err := s.policyRuleDao.UpdateStatus(id, status); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新策略规则状态失败: %v", err))
		return
	}

	statusText := "禁用"
	if status == 1 {
		statusText = "启用"
	}

	result.Success(c, fmt.Sprintf("策略规则已%s", statusText))
}

// PreviewPolicyRule 预览策略规则查询结果
func (s *AlertPolicyServiceImpl) PreviewPolicyRule(c *gin.Context, ruleID uint) {
	// 获取规则信息
	rule, err := s.policyRuleDao.GetByID(ruleID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("策略规则不存在: %v", err))
		return
	}

	// 从配置中获取Prometheus地址
	prometheusURL := config.Config.Monitor.Prometheus.URL
	if prometheusURL == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "Prometheus地址未配置")
		return
	}

	// 构建查询URL
	queryURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, url.QueryEscape(rule.PromQL))

	// 发送HTTP请求
	resp, err := http.Get(queryURL)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询Prometheus失败: %v", err))
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("读取响应失败: %v", err))
		return
	}

	// 解析响应
	var prometheusResp struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Value  []any             `json:"value"`
			} `json:"result"`
		} `json:"data"`
		Error     string `json:"error"`
		ErrorType string `json:"errorType"`
	}

	if err := json.Unmarshal(body, &prometheusResp); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("解析响应失败: %v", err))
		return
	}

	if prometheusResp.Status != "success" {
		errMsg := prometheusResp.Error
		if errMsg == "" {
			errMsg = "未知错误"
		}
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败[%s]: %s", prometheusResp.ErrorType, errMsg))
		return
	}

	// 构建预览响应
	results := make([]model.PolicyRulePreviewResult, 0)
	for _, item := range prometheusResp.Data.Result {
		// 解析值
		var value float64
		if len(item.Value) >= 2 {
			if v, ok := item.Value[1].(string); ok {
				fmt.Sscanf(v, "%f", &value)
			}
		}

		// 获取hostname和instance
		hostname := item.Metric["hostname"]
		instance := item.Metric["instance"]
		if hostname == "" {
			hostname = instance
		}

		// 判断是否异常（超过阈值）
		isAnomalous := false
		switch rule.Operator {
		case ">":
			isAnomalous = value > rule.Threshold
		case ">=":
			isAnomalous = value >= rule.Threshold
		case "<":
			isAnomalous = value < rule.Threshold
		case "<=":
			isAnomalous = value <= rule.Threshold
		case "==":
			isAnomalous = value == rule.Threshold
		case "!=":
			isAnomalous = value != rule.Threshold
		}

		results = append(results, model.PolicyRulePreviewResult{
			Hostname:    hostname,
			Instance:    instance,
			Value:       value,
			IsAnomalous: isAnomalous,
			Labels:      item.Metric,
		})
	}

	// 返回预览响应
	response := model.PolicyRulePreviewResponse{
		RuleID:      rule.ID,
		RuleName:    rule.RuleName,
		PromQL:      rule.PromQL,
		Threshold:   rule.Threshold,
		Operator:    rule.Operator,
		Unit:        rule.Unit,
		QueryTime:   fmt.Sprintf("%d", resp.Header.Get("Date")),
		ResultCount: len(results),
		Results:     results,
	}

	result.Success(c, response)
}
