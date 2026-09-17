package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service/alert/eval"
	"dodevops-api/api/monitor/service/alert/notify"
	"dodevops-api/api/monitor/service/alert/process"
	"dodevops-api/common/prometheus"
)

// RuleWorker 规则工作器
type RuleWorker struct {
	policy     *model.AlertPolicy
	rule       *model.AlertPolicyRule
	promClient *prometheus.Client
	evaluator  *eval.RuleEvaluator
	processor  *process.EventProcessor
	notifier   *notify.AlertNotifier
}

// AlertScheduler 告警调度器
type AlertScheduler struct {
	policyDao     dao.AlertPolicyDao
	policyRuleDao dao.AlertPolicyRuleDao
	promClient    *prometheus.Client
	evaluator     *eval.RuleEvaluator
	processor     *process.EventProcessor
	notifier      *notify.AlertNotifier

	// 规则工作器映射 (policyID -> RuleWorkers)
	workers map[uint][]*RuleWorker
	mu      sync.RWMutex

	// 同步定时器 (每30秒检查一次所有规则)
	syncTicker  *time.Ticker
	checkTicker *time.Ticker
	stopChan    chan struct{}
}

// NewAlertScheduler 创建告警调度器
func NewAlertScheduler(prometheusAddr string) (*AlertScheduler, error) {
	promClient, err := prometheus.NewClient(prometheusAddr)
	if err != nil {
		return nil, fmt.Errorf("创建Prometheus客户端失败: %v", err)
	}

	evaluator := eval.NewRuleEvaluator(promClient)
	processor := process.NewEventProcessor()
	notifier := notify.NewAlertNotifier()

	// 从数据库加载正在触发的告警
	if err := processor.LoadFiringEventsFromDB(); err != nil {
		log.Printf("加载触发中的告警失败: %v", err)
	}

	return &AlertScheduler{
		policyDao:     dao.NewAlertPolicyDao(),
		policyRuleDao: dao.NewAlertPolicyRuleDao(),
		promClient:    promClient,
		evaluator:     evaluator,
		processor:     processor,
		notifier:      notifier,
		workers:       make(map[uint][]*RuleWorker),
		stopChan:      make(chan struct{}),
	}, nil
}

// Start 启动调度器
func (s *AlertScheduler) Start() {
	log.Println("告警调度器启动中...")

	// 首次同步规则
	if err := s.syncRules(); err != nil {
		log.Printf("同步规则失败: %v", err)
	}

	// 启动定时同步(每30秒同步一次策略和规则)
	s.syncTicker = time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-s.syncTicker.C:
				if err := s.syncRules(); err != nil {
					log.Printf("同步规则失败: %v", err)
				}
			case <-s.stopChan:
				return
			}
		}
	}()

	// 启动定时检查(每30秒检查一次所有规则)
	s.checkTicker = time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-s.checkTicker.C:
				s.checkAllRules()
			case <-s.stopChan:
				return
			}
		}
	}()

	log.Println("告警调度器已启动")
}

// Stop 停止调度器
func (s *AlertScheduler) Stop() {
	log.Println("告警调度器停止中...")

	// 停止同步定时器
	if s.syncTicker != nil {
		s.syncTicker.Stop()
	}
	if s.checkTicker != nil {
		s.checkTicker.Stop()
	}
	close(s.stopChan)

	log.Println("告警调度器已停止")
}

// syncRules 同步规则
func (s *AlertScheduler) syncRules() error {
	// 获取所有启用的策略
	policies, err := s.policyDao.GetAllEnabled()
	if err != nil {
		return fmt.Errorf("获取启用策略失败: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 清空旧的workers
	oldWorkers := s.workers
	s.workers = make(map[uint][]*RuleWorker)

	// 收集所有启用的规则ID
	enabledRuleIDs := make(map[uint]bool)

	// 为每个策略创建规则工作器
	for _, policy := range policies {
		// 获取策略下的所有启用规则
		rules, err := s.policyRuleDao.GetByPolicyID(policy.ID)
		if err != nil {
			log.Printf("获取策略规则失败 [策略:%s]: %v", policy.Name, err)
			continue
		}

		var policyWorkers []*RuleWorker
		for _, rule := range rules {
			if rule.Status == 1 { // 只处理启用的规则
				enabledRuleIDs[rule.ID] = true
				worker := &RuleWorker{
					policy:     policy,
					rule:       rule,
					promClient: s.promClient,
					evaluator:  s.evaluator,
					processor:  s.processor,
					notifier:   s.notifier,
				}
				policyWorkers = append(policyWorkers, worker)
			}
		}

		if len(policyWorkers) > 0 {
			s.workers[policy.ID] = policyWorkers
			log.Printf("策略工作器已创建 [策略:%s, 规则数:%d]", policy.Name, len(policyWorkers))
		}
	}

	// 检查是否有规则被禁用，需要恢复对应的告警
	for _, workers := range oldWorkers {
		for _, oldWorker := range workers {
			// 如果旧规则不在新的启用列表中，说明被禁用了
			if !enabledRuleIDs[oldWorker.rule.ID] {
				log.Printf("🔕 规则已禁用，恢复相关告警 [策略:%s, 规则:%s]", oldWorker.policy.Name, oldWorker.rule.RuleName)
				// 强制恢复该规则的所有告警
				s.processor.ForceRecoverRuleAlerts(oldWorker.policy, oldWorker.rule)
			}
		}
	}

	log.Printf("规则同步完成: 当前运行%d个策略", len(s.workers))
	return nil
}

// checkAllRules 检查所有规则
func (s *AlertScheduler) checkAllRules() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for policyID, workers := range s.workers {
		for _, worker := range workers {
			go worker.Check() // 并发检查每个规则
		}
		log.Printf("策略检查完成 [策略ID:%d, 规则数:%d]", policyID, len(workers))
	}
}

// GetStatus 获取调度器状态
func (s *AlertScheduler) GetStatus() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalRules := 0
	for _, workers := range s.workers {
		totalRules += len(workers)
	}

	return map[string]any{
		"policies_count":     len(s.workers),
		"rules_count":        totalRules,
		"firing_events":      s.processor.GetFiringEvents(),
		"prometheus_address": s.promClient.GetAddress(),
	}
}

// Check 检查规则(RuleWorker方法)
func (w *RuleWorker) Check() {
	startTime := time.Now()

	// 1. 执行PromQL查询
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("开始执行PromQL查询 [策略:%s, 规则:%s, PromQL:%s]", w.policy.Name, w.rule.RuleName, w.rule.PromQL)
	results, err := w.promClient.Query(ctx, w.rule.PromQL)
	if err != nil {
		log.Printf("PromQL查询失败 [策略:%s, 规则:%s]: %v", w.policy.Name, w.rule.RuleName, err)
		return
	}

	log.Printf("PromQL查询成功 [策略:%s, 规则:%s, 结果数:%d]", w.policy.Name, w.rule.RuleName, len(results))

	// 2. 过滤满足阈值条件的异常点
	var anomalyResults []prometheus.QueryResult
	for _, result := range results {
		isAnomaly := false
		switch w.rule.Operator {
		case ">=":
			isAnomaly = result.Value >= w.rule.Threshold
		case ">":
			isAnomaly = result.Value > w.rule.Threshold
		case "<=":
			isAnomaly = result.Value <= w.rule.Threshold
		case "<":
			isAnomaly = result.Value < w.rule.Threshold
		case "==":
			isAnomaly = result.Value == w.rule.Threshold
		case "!=":
			isAnomaly = result.Value != w.rule.Threshold
		}

		if isAnomaly {
			anomalyResults = append(anomalyResults, result)
			log.Printf("检测到异常 [策略:%s, 规则:%s, 值:%.2f, 阈值:%.2f, 运算符:%s]",
				w.policy.Name, w.rule.RuleName, result.Value, w.rule.Threshold, w.rule.Operator)
		}
	}

	// 3. 判断是否有异常点
	hasAnomalies := len(anomalyResults) > 0

	if hasAnomalies {
		// 4. 处理告警事件
		for _, result := range anomalyResults {
			// 构建告警点
			point := &eval.AnomalyPoint{
				Metric:      result.Metric,
				Value:       result.Value,
				Threshold:   w.rule.Threshold,
				Severity:    w.rule.Severity,
				Labels:      w.buildLabels(result.Metric),
				Annotations: w.buildAnnotations(),
				Timestamp:   result.Time,
				IsRecovered: false,
			}

			// 处理异常点(去重、状态管理、写入数据库)
			if err := w.processor.ProcessPolicyRule(w.policy, w.rule, point); err != nil {
				log.Printf("处理异常点失败 [规则:%s]: %v", w.rule.RuleName, err)
				continue
			}

			// 注意：只有当告警真正创建或更新到数据库后，才发送通知
			// ProcessPolicyRule 会在达到连续触发次数后创建告警
			// 在未达到连续次数时，会返回 nil 但不创建记录
			// 通知逻辑会自动检查记录是否存在，不存在则跳过
			// 所以这里始终尝试发送通知，notifier内部会处理
			if err := w.notifier.NotifyPolicyRule(w.policy, w.rule, point); err != nil {
				// 如果是 "record not found" 错误，说明还在计数阶段，这是正常的
				if !strings.Contains(err.Error(), "record not found") {
					log.Printf("发送告警通知失败 [规则:%s]: %v", w.rule.RuleName, err)
				}
			}
		}
	} else {
		// 5. 没有异常点,检查是否需要发送恢复通知
		if w.policy.SendRecovery {
			log.Printf("未检测到异常,检查恢复 [策略:%s, 规则:%s]", w.policy.Name, w.rule.RuleName)
			recoveredRecords, err := w.processor.ProcessPolicyRuleRecovery(w.policy, w.rule)
			if err != nil {
				log.Printf("处理恢复失败 [规则:%s]: %v", w.rule.RuleName, err)
			} else if len(recoveredRecords) > 0 {
				// 发送恢复通知
				for _, record := range recoveredRecords {
					log.Printf("发送恢复通知 [策略:%s, 规则:%s, 记录ID:%d]", w.policy.Name, w.rule.RuleName, record.ID)
					if err := w.notifier.NotifyRecovery(w.policy, w.rule, record); err != nil {
						log.Printf("发送恢复通知失败 [规则:%s]: %v", w.rule.RuleName, err)
					}
				}
			}
		}
	}

	duration := time.Since(startTime)
	log.Printf("规则检查完成 [策略:%s, 规则:%s, 耗时:%v, 异常点数:%d]",
		w.policy.Name, w.rule.RuleName, duration, len(results))
}

// buildLabels 构建标签
func (w *RuleWorker) buildLabels(metric map[string]string) map[string]string {
	labels := make(map[string]string)

	// 添加指标标签
	for k, v := range metric {
		labels[k] = v
	}

	// 添加策略信息
	labels["policy_id"] = fmt.Sprintf("%d", w.policy.ID)
	labels["policy_name"] = w.policy.Name
	labels["group_id"] = fmt.Sprintf("%d", w.policy.GroupID)
	labels["group_name"] = w.policy.GroupName

	// 添加规则信息
	labels["rule_id"] = fmt.Sprintf("%d", w.rule.ID)
	labels["rule_name"] = w.rule.RuleName
	labels["category"] = w.rule.Category
	labels["severity"] = fmt.Sprintf("%d", w.rule.Severity)

	// 添加附加标签
	if w.policy.AppendTags != "" {
		var appendTags []string
		if err := json.Unmarshal([]byte(w.policy.AppendTags), &appendTags); err == nil {
			for i, tag := range appendTags {
				labels[fmt.Sprintf("tag_%d", i)] = tag
			}
		}
	}

	return labels
}

// buildAnnotations 构建注解
func (w *RuleWorker) buildAnnotations() map[string]string {
	return map[string]string{
		"description":       w.rule.Description,
		"metric_key":        w.rule.MetricKey,
		"operator":          w.rule.Operator,
		"threshold":         fmt.Sprintf("%.2f", w.rule.Threshold),
		"unit":              w.rule.Unit,
		"continuous_times":  fmt.Sprintf("%d", w.rule.ContinuousTimes),
		"duration":          fmt.Sprintf("%ds", w.rule.Duration),
		"recover_duration":  fmt.Sprintf("%ds", w.rule.RecoverDuration),
		"silence_period":    fmt.Sprintf("%ds", w.policy.SilencePeriod),
		"repeat_interval":   fmt.Sprintf("%ds", w.policy.RepeatInterval),
		"max_notify_count":  fmt.Sprintf("%d", w.policy.MaxNotifyCount),
	}
}
