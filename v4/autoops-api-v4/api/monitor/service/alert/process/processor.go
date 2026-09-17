package process

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"dodevops-api/api/cmdb/dao"
	alertDao "dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service/alert/eval"
	"dodevops-api/common/util"
)

// EventProcessor 事件处理器
type EventProcessor struct {
	ruleDao       alertDao.AlertRuleDao
	recordDao     alertDao.AlertRecordDao
	muteDao       alertDao.AlertMuteDao
	hostDao       dao.CmdbHostDao
	policyDao     alertDao.AlertPolicyDao
	policyRuleDao alertDao.AlertPolicyRuleDao

	// 内存中的告警事件表(用于去重和节流)
	firingEvents  map[string]*model.AlertRecord // 正在触发的告警 fingerprint -> record
	pendingEvents map[string]time.Time          // 等待恢复的告警 fingerprint -> lastCheckTime

	// 连续触发计数器(用于实现连续N次才告警)
	continuousCounter map[string]*ContinuousTrigger // fingerprint -> 触发信息

	// 恢复检测计数器(用于实现持续正常N秒才恢复)
	recoveryTracker map[string]time.Time // fingerprint -> 首次检测到正常的时间
	mu              sync.RWMutex
}

// ContinuousTrigger 连续触发信息
type ContinuousTrigger struct {
	Count          int       // 当前连续触发次数
	FirstTriggerAt time.Time // 首次触发时间
	LastTriggerAt  time.Time // 最后一次触发时间
}

// NewEventProcessor 创建事件处理器
func NewEventProcessor() *EventProcessor {
	return &EventProcessor{
		ruleDao:           alertDao.NewAlertRuleDao(),
		recordDao:         alertDao.NewAlertRecordDao(),
		muteDao:           alertDao.NewAlertMuteDao(),
		hostDao:           dao.NewCmdbHostDao(),
		policyDao:         alertDao.NewAlertPolicyDao(),
		policyRuleDao:     alertDao.NewAlertPolicyRuleDao(),
		firingEvents:      make(map[string]*model.AlertRecord),
		pendingEvents:     make(map[string]time.Time),
		continuousCounter: make(map[string]*ContinuousTrigger),
		recoveryTracker:   make(map[string]time.Time),
	}
}

// Process 处理异常点
func (p *EventProcessor) Process(rule *model.AlertRule, points []*eval.AnomalyPoint) error {
	for _, point := range points {
		if err := p.processPoint(rule, point); err != nil {
			log.Printf("处理异常点失败: %v", err)
		}
	}
	return nil
}

// processPoint 处理单个异常点
func (p *EventProcessor) processPoint(rule *model.AlertRule, point *eval.AnomalyPoint) error {
	// 生成事件指纹
	fingerprint := eval.GenerateFingerprint(rule.ID, point.Labels)

	// 检查静默
	if p.isMuted(rule, point) {
		log.Printf("告警被静默 [规则:%s, fingerprint:%s]", rule.Name, fingerprint)
		return nil
	}

	// 获取目标名称
	targetName := p.getTargetName(rule, point)

	p.mu.Lock()
	defer p.mu.Unlock()

	// 检查是否已存在触发中的告警
	existingRecord, exists := p.firingEvents[fingerprint]

	if !exists {
		// 新告警 - 创建记录
		return p.createNewAlert(rule, point, fingerprint, targetName)
	}

	// 已存在的告警 - 更新
	return p.updateExistingAlert(rule, point, existingRecord)
}

// createNewAlert 创建新告警
func (p *EventProcessor) createNewAlert(rule *model.AlertRule, point *eval.AnomalyPoint, fingerprint, targetName string) error {
	// 序列化标签
	labelsJSON, _ := json.Marshal(point.Labels)

	// 创建告警记录
	record := &model.AlertRecord{
		RuleID:       rule.ID,
		RuleName:     rule.Name,
		TargetType:   rule.TargetType,
		TargetID:     p.getTargetID(point),
		TargetName:   targetName,
		MetricName:   rule.MetricName,
		CurrentValue: point.Value,
		Threshold:    rule.Threshold,
		Severity:     rule.Severity,
		Status:       "firing",
		StartTime:    util.HTime{Time: time.Now()},
		Labels:       string(labelsJSON),
		GroupID:      rule.GroupID,
		GroupName:    rule.GroupName,
		Fingerprint:  fingerprint,
		Message:      eval.BuildAlertMessage(rule, point, targetName),
	}

	// 保存到数据库
	if err := p.recordDao.Create(record); err != nil {
		return fmt.Errorf("创建告警记录失败: %v", err)
	}

	// 加入内存映射
	p.firingEvents[fingerprint] = record

	log.Printf("新告警创建 [规则:%s, 目标:%s, fingerprint:%s]", rule.Name, targetName, fingerprint)

	return nil
}

// updateExistingAlert 更新已存在的告警
func (p *EventProcessor) updateExistingAlert(rule *model.AlertRule, point *eval.AnomalyPoint, record *model.AlertRecord) error {
	// 获取目标名称
	targetName := p.getTargetName(rule, point)

	// 计算当前持续时间
	duration := int(time.Since(record.StartTime.Time).Seconds())

	// 只更新告警指标相关字段，不影响通知状态（notify_time、notify_count）
	if err := p.recordDao.UpdateAlertMetrics(record.ID, point.Value, duration, targetName); err != nil {
		return fmt.Errorf("更新告警记录失败: %v", err)
	}

	// 同步更新内存中的记录
	record.CurrentValue = point.Value
	record.Duration = duration
	record.TargetName = targetName
	record.UpdatedAt = util.HTime{Time: time.Now()}

	return nil
}

// ProcessRecovery 处理恢复
func (p *EventProcessor) ProcessRecovery(rule *model.AlertRule) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 检查所有正在触发的告警
	for fingerprint, record := range p.firingEvents {
		// 检查是否应该恢复
		if p.shouldRecover(rule, record) {
			if err := p.recoverAlert(rule, record, fingerprint); err != nil {
				log.Printf("恢复告警失败: %v", err)
			}
		}
	}

	return nil
}

// shouldRecover 判断是否应该恢复
func (p *EventProcessor) shouldRecover(rule *model.AlertRule, record *model.AlertRecord) bool {
	// 检查恢复时间
	if time.Since(record.UpdatedAt.Time) < time.Duration(rule.RecoverDuration)*time.Second {
		return false
	}

	return true
}

// recoverAlert 恢复告警
func (p *EventProcessor) recoverAlert(rule *model.AlertRule, record *model.AlertRecord, fingerprint string) error {
	now := util.HTime{Time: time.Now()}
	record.Status = "resolved"
	record.EndTime = &now
	record.Duration = int(time.Since(record.StartTime.Time).Seconds())

	// 保存到数据库
	if err := p.recordDao.Update(record); err != nil {
		return fmt.Errorf("更新告警记录失败: %v", err)
	}

	// 从内存映射中移除
	delete(p.firingEvents, fingerprint)

	log.Printf("告警恢复 [规则:%s, fingerprint:%s, 持续:%ds]", rule.Name, fingerprint, record.Duration)

	return nil
}

// isMuted 检查是否被静默
func (p *EventProcessor) isMuted(rule *model.AlertRule, point *eval.AnomalyPoint) bool {
	mutes, err := p.muteDao.GetActiveMutes()
	if err != nil {
		return false
	}

	for _, mute := range mutes {
		if p.matchMute(mute, rule, point) {
			return true
		}
	}

	return false
}

// matchMute 匹配静默规则
func (p *EventProcessor) matchMute(mute *model.AlertMute, rule *model.AlertRule, point *eval.AnomalyPoint) bool {
	// 检查规则ID匹配
	if mute.RuleIDs != "" {
		var ruleIDs []uint
		json.Unmarshal([]byte(mute.RuleIDs), &ruleIDs)

		matched := false
		for _, id := range ruleIDs {
			if id == rule.ID {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查目标类型匹配
	if mute.TargetType != "" && mute.TargetType != rule.TargetType {
		return false
	}

	// 检查目标ID匹配
	if mute.TargetIDs != "" {
		var targetIDs []uint
		json.Unmarshal([]byte(mute.TargetIDs), &targetIDs)

		targetID := p.getTargetID(point)
		matched := false
		for _, id := range targetIDs {
			if id == targetID {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查标签匹配(简单实现,后续可以扩展为正则匹配)
	if mute.LabelMatchers != "" {
		// TODO: 实现标签匹配逻辑
	}

	return true
}

// getTargetID 从标签中提取目标ID
func (p *EventProcessor) getTargetID(point *eval.AnomalyPoint) uint {
	// 从instance标签获取IP,再查询主机ID
	instance, ok := point.Labels["instance"]
	if !ok {
		return 0
	}

	host, err := p.hostDao.GetCmdbHostByIP(instance)
	if err != nil {
		return 0
	}

	return host.ID
}

// getTargetName 获取目标名称
func (p *EventProcessor) getTargetName(rule *model.AlertRule, point *eval.AnomalyPoint) string {
	if rule.TargetType == "host" {
		targetID := p.getTargetID(point)
		if targetID > 0 {
			host, err := p.hostDao.GetCmdbHostById(targetID)
			if err == nil {
				return host.HostName
			}
		}
	}

	// 默认使用instance标签
	if instance, ok := point.Labels["instance"]; ok {
		return instance
	}

	return "unknown"
}

// GetFiringEvents 获取正在触发的告警数量
func (p *EventProcessor) GetFiringEvents() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.firingEvents)
}

// LoadFiringEventsFromDB 从数据库加载正在触发的告警
func (p *EventProcessor) LoadFiringEventsFromDB() error {
	records, _, err := p.recordDao.GetList("firing", "", "", nil, nil, nil, nil, 1, 10000)
	if err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// 统计每个策略和规则的告警数量，用于同步健康状态
	policyAlerts := make(map[uint]int)
	ruleAlerts := make(map[uint]int)

	for _, record := range records {
		p.firingEvents[record.Fingerprint] = record

		// 统计告警数量
		if record.PolicyID > 0 {
			policyAlerts[record.PolicyID]++
		}
		if record.RuleID > 0 {
			ruleAlerts[record.RuleID]++
		}
	}

	log.Printf("从数据库加载了 %d 条触发中的告警", len(records))

	// 同步健康状态：有告警的策略和规则应该是 alerting 状态
	for policyID := range policyAlerts {
		if err := p.policyDao.UpdateHealthStatus(policyID, "alerting"); err != nil {
			log.Printf("同步策略健康状态失败 [策略ID:%d]: %v", policyID, err)
		}
	}

	for ruleID := range ruleAlerts {
		if err := p.policyRuleDao.UpdateHealthStatus(ruleID, "alerting"); err != nil {
			log.Printf("同步规则健康状态失败 [规则ID:%d]: %v", ruleID, err)
		}
	}

	log.Printf("已同步健康状态 [策略:%d个, 规则:%d个]", len(policyAlerts), len(ruleAlerts))

	return nil
}

// ProcessPolicyRule 处理告警策略规则异常点
func (p *EventProcessor) ProcessPolicyRule(policy *model.AlertPolicy, rule *model.AlertPolicyRule, point *eval.AnomalyPoint) error {
	// 生成事件指纹
	fingerprint := eval.GenerateFingerprint(rule.ID, point.Labels)

	log.Printf("处理告警 [策略:%s, 规则:%s, fingerprint:%s]", policy.Name, rule.RuleName, fingerprint)

	// 检查静默 (可选：如果策略有静默规则)
	// if p.isMuted(rule, point) {
	// 	log.Printf("告警被静默 [规则:%s, fingerprint:%s]", rule.RuleName, fingerprint)
	// 	return nil
	// }

	// 获取目标名称
	targetName := p.getTargetNameFromLabels(point.Labels)
	targetID := p.getTargetIDFromLabels(point.Labels)

	p.mu.Lock()
	defer p.mu.Unlock()

	// 检查是否已存在触发中的告警
	existingRecord, exists := p.firingEvents[fingerprint]

	if !exists {
		// 新告警 - 检查连续触发次数
		if rule.ContinuousTimes > 1 {
			// 需要连续触发多次才告警
			trigger, hasTrigger := p.continuousCounter[fingerprint]
			now := time.Now()

			if !hasTrigger {
				// 首次触发，初始化计数器
				p.continuousCounter[fingerprint] = &ContinuousTrigger{
					Count:          1,
					FirstTriggerAt: now,
					LastTriggerAt:  now,
				}
				log.Printf("🔔 [第1次触发] 策略:%s | 规则:%s | 目标:%s", policy.Name, rule.RuleName, targetName)
				log.Printf("   ├─ 连续触发要求: %d次 (检查间隔约30秒)", rule.ContinuousTimes)
				log.Printf("   ├─ 当前进度: 1/%d 次 ⏳", rule.ContinuousTimes)
				log.Printf("   ├─ 时间窗口: %ds", rule.Duration*rule.ContinuousTimes)
				log.Printf("   └─ 首次触发时间: %s", now.Format("15:04:05"))
				return nil
			}

			// 检查是否在持续时间窗口内（Duration * ContinuousTimes）
			maxDuration := time.Duration(rule.Duration*rule.ContinuousTimes) * time.Second
			elapsedTime := now.Sub(trigger.FirstTriggerAt)
			timeSinceLastTrigger := now.Sub(trigger.LastTriggerAt)

			if elapsedTime > maxDuration {
				// 超过时间窗口，重置计数器
				log.Printf("⏰ [触发超时-重置] 策略:%s | 规则:%s", policy.Name, rule.RuleName)
				log.Printf("   ├─ 已过时间: %.0fs (超过窗口: %ds)", elapsedTime.Seconds(), rule.Duration*rule.ContinuousTimes)
				log.Printf("   ├─ 之前进度: %d/%d 次", trigger.Count, rule.ContinuousTimes)
				log.Printf("   └─ 重新开始: 1/%d 次 ⏳", rule.ContinuousTimes)
				p.continuousCounter[fingerprint] = &ContinuousTrigger{
					Count:          1,
					FirstTriggerAt: now,
					LastTriggerAt:  now,
				}
				return nil
			}

			// 递增计数器
			trigger.Count++
			trigger.LastTriggerAt = now

			log.Printf("🔔 [第%d次触发] 策略:%s | 规则:%s | 目标:%s", trigger.Count, policy.Name, rule.RuleName, targetName)
			log.Printf("   ├─ 当前进度: %d/%d 次", trigger.Count, rule.ContinuousTimes)
			log.Printf("   ├─ 距首次触发: %.0fs (窗口: %ds)", elapsedTime.Seconds(), rule.Duration*rule.ContinuousTimes)
			log.Printf("   ├─ 距上次检查: %.0fs", timeSinceLastTrigger.Seconds())
			log.Printf("   └─ 首次触发: %s | 上次触发: %s", trigger.FirstTriggerAt.Format("15:04:05"), trigger.LastTriggerAt.Format("15:04:05"))

			// 检查是否达到连续触发次数
			if trigger.Count < rule.ContinuousTimes {
				log.Printf("   ⏳ 还需 %d 次才会创建告警并发送通知", rule.ContinuousTimes-trigger.Count)
				return nil
			}

			// 达到连续触发次数，清除计数器，创建告警
			delete(p.continuousCounter, fingerprint)
			log.Printf("✅ [达到触发条件] 策略:%s | 规则:%s", policy.Name, rule.RuleName)
			log.Printf("   ├─ 已连续触发: %d 次 ✓", trigger.Count)
			log.Printf("   ├─ 总耗时: %.0fs (平均间隔: %.0fs)", elapsedTime.Seconds(), elapsedTime.Seconds()/float64(trigger.Count-1))
			log.Printf("   └─ 准备创建告警记录并发送通知 🚨")
		}

		// 创建新告警记录
		log.Printf("📝 [创建告警] 策略:%s | 规则:%s | 目标:%s", policy.Name, rule.RuleName, targetName)
		return p.createNewPolicyAlert(policy, rule, point, fingerprint, targetName, targetID)
	}

	// 已存在的告警 - 更新（已经触发过的告警，直接更新）
	// 清除连续计数器（如果存在）
	delete(p.continuousCounter, fingerprint)
	// 清除恢复跟踪器（告警重新触发，不再正常）
	delete(p.recoveryTracker, fingerprint)

	log.Printf("更新已存在告警 [策略:%s, 规则:%s, fingerprint:%s, 已存在ID:%d]", policy.Name, rule.RuleName, fingerprint, existingRecord.ID)
	return p.updateExistingPolicyAlert(policy, rule, point, existingRecord)
}

// createNewPolicyAlert 创建新的策略告警
func (p *EventProcessor) createNewPolicyAlert(policy *model.AlertPolicy, rule *model.AlertPolicyRule, point *eval.AnomalyPoint, fingerprint, targetName string, targetID uint) error {
	// 将 operator 添加到 labels 中（用于后续消息渲染）
	if point.Labels == nil {
		point.Labels = make(map[string]string)
	}
	point.Labels["operator"] = rule.Operator

	// 序列化标签
	labelsJSON, _ := json.Marshal(point.Labels)

	// 构建告警消息
	message := p.buildPolicyAlertMessage(policy, rule, point, targetName)

	// 创建告警记录
	record := &model.AlertRecord{
		PolicyID:     policy.ID,
		PolicyName:   policy.Name,
		RuleID:       rule.ID,
		RuleName:     rule.RuleName,
		TargetType:   "host", // 从标签中判断
		TargetID:     targetID,
		TargetName:   targetName,
		MetricName:   rule.MetricKey,
		CurrentValue: point.Value,
		Threshold:    rule.Threshold,
		Severity:     rule.Severity,
		Status:       "firing",
		StartTime:    util.HTime{Time: time.Now()},
		Labels:       string(labelsJSON),
		GroupID:      policy.GroupID,
		GroupName:    policy.GroupName,
		Fingerprint:  fingerprint,
		Message:      message,
	}

	// 保存到数据库
	if err := p.recordDao.Create(record); err != nil {
		// 如果是重复键错误，说明记录已存在
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "idx_monitor_alert_record_fingerprint") {
			log.Printf("告警记录已存在，尝试加载并更新 [fingerprint:%s]", fingerprint)

			// 从数据库查询已存在的记录（只查firing状态）
			existingRecord, dbErr := p.recordDao.GetByFingerprint(fingerprint)
			if dbErr != nil {
				// 如果查不到firing状态的记录，说明之前的告警已恢复
				// 这是一个新的告警触发，需要先将旧记录的fingerprint清空，再创建新记录
				log.Printf("未找到firing状态的告警，可能是已恢复的告警重新触发 [fingerprint:%s]", fingerprint)

				// 清空旧记录的fingerprint（已恢复的记录不需要保持fingerprint唯一性）
				if clearErr := p.recordDao.ClearFingerprintForResolved(fingerprint); clearErr != nil {
					log.Printf("清空已恢复告警的fingerprint失败: %v", clearErr)
				}

				// 重新尝试创建
				if retryErr := p.recordDao.Create(record); retryErr != nil {
					return fmt.Errorf("重新创建告警记录失败: %v", retryErr)
				}

				// 加入内存映射
				p.firingEvents[fingerprint] = record
				log.Printf("已恢复告警重新触发，创建新记录 [策略:%s, 规则:%s, fingerprint:%s]", policy.Name, rule.RuleName, fingerprint)
				return nil
			}

			// 找到了firing状态的记录，加入内存并更新
			p.firingEvents[fingerprint] = existingRecord
			return p.updateExistingPolicyAlert(policy, rule, point, existingRecord)
		}
		return fmt.Errorf("创建告警记录失败: %v", err)
	}

	// 加入内存映射
	p.firingEvents[fingerprint] = record

	// 更新策略和规则的健康状态为告警中
	if err := p.policyDao.UpdateHealthStatus(policy.ID, "alerting"); err != nil {
		log.Printf("更新策略健康状态失败: %v", err)
	}
	if err := p.policyRuleDao.UpdateHealthStatus(rule.ID, "alerting"); err != nil {
		log.Printf("更新规则健康状态失败: %v", err)
	}

	log.Printf("新告警创建 [策略:%s, 规则:%s, 目标:%s, fingerprint:%s]", policy.Name, rule.RuleName, targetName, fingerprint)

	return nil
}

// updateExistingPolicyAlert 更新已存在的策略告警
func (p *EventProcessor) updateExistingPolicyAlert(policy *model.AlertPolicy, rule *model.AlertPolicyRule, point *eval.AnomalyPoint, record *model.AlertRecord) error {
	// 获取目标名称
	targetName := p.getTargetNameFromLabels(point.Labels)

	// 计算当前持续时间
	duration := int(time.Since(record.StartTime.Time).Seconds())

	// 只更新告警指标相关字段，不影响通知状态（notify_time、notify_count）
	if err := p.recordDao.UpdateAlertMetrics(record.ID, point.Value, duration, targetName); err != nil {
		return fmt.Errorf("更新告警记录失败: %v", err)
	}

	// 确保策略和规则的健康状态为告警中（防止状态不一致）
	if err := p.policyDao.UpdateHealthStatus(policy.ID, "alerting"); err != nil {
		log.Printf("更新策略健康状态失败: %v", err)
	}
	if err := p.policyRuleDao.UpdateHealthStatus(rule.ID, "alerting"); err != nil {
		log.Printf("更新规则健康状态失败: %v", err)
	}

	// 同步更新内存中的记录
	record.CurrentValue = point.Value
	record.Duration = duration
	record.TargetName = targetName
	record.UpdatedAt = util.HTime{Time: time.Now()}

	log.Printf("告警更新 [策略:%s, 规则:%s, 持续:%ds, 当前值:%.2f]", policy.Name, rule.RuleName, duration, point.Value)

	return nil
}

// ProcessPolicyRuleRecovery 处理策略规则恢复
func (p *EventProcessor) ProcessPolicyRuleRecovery(policy *model.AlertPolicy, rule *model.AlertPolicyRule) ([]*model.AlertRecord, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var recoveredRecords []*model.AlertRecord
	now := time.Now()

	// 检查所有正在触发的告警
	for fingerprint, record := range p.firingEvents {
		// 检查是否属于该规则
		if record.RuleID == rule.ID && record.PolicyID == policy.ID {
			// 检查恢复时间跟踪器
			firstNormalTime, exists := p.recoveryTracker[fingerprint]
			if !exists {
				// 首次检测到正常，记录时间
				p.recoveryTracker[fingerprint] = now
				log.Printf("🟢 [首次检测正常] 策略:%s | 规则:%s | 目标:%s", policy.Name, rule.RuleName, record.TargetName)
				log.Printf("   ├─ 需要持续正常: %ds", rule.RecoverDuration)
				log.Printf("   └─ 首次正常时间: %s", now.Format("15:04:05"))
				continue
			}

			// 计算已经持续正常的时间
			normalDuration := now.Sub(firstNormalTime)
			log.Printf("🟢 [持续正常] 策略:%s | 规则:%s | 目标:%s", policy.Name, rule.RuleName, record.TargetName)
			log.Printf("   ├─ 已持续正常: %.0fs / %ds", normalDuration.Seconds(), rule.RecoverDuration)
			log.Printf("   └─ 首次正常: %s | 当前: %s", firstNormalTime.Format("15:04:05"), now.Format("15:04:05"))

			// 检查是否达到恢复时间
			if normalDuration < time.Duration(rule.RecoverDuration)*time.Second {
				log.Printf("   ⏳ 还需 %.0fs 才会标记为已恢复", float64(rule.RecoverDuration)-normalDuration.Seconds())
				continue
			}

			// 达到恢复时间，执行恢复
			log.Printf("✅ [达到恢复条件] 策略:%s | 规则:%s", policy.Name, rule.RuleName)
			log.Printf("   ├─ 已持续正常: %.0fs", normalDuration.Seconds())
			log.Printf("   └─ 准备标记为已恢复并发送通知 📢")

			if err := p.recoverPolicyAlert(policy, rule, record, fingerprint); err != nil {
				log.Printf("恢复告警失败: %v", err)
			} else {
				// 添加到已恢复列表
				recoveredRecords = append(recoveredRecords, record)
				// 清除恢复跟踪器
				delete(p.recoveryTracker, fingerprint)
			}
		}
	}

	// 只在有告警恢复时才清理连续触发计数器
	// 因为恢复说明指标已正常，相关的计数器应该重置
	if len(recoveredRecords) > 0 {
		// 清理所有连续触发计数器（告警恢复说明指标已正常）
		for fp := range p.continuousCounter {
			delete(p.continuousCounter, fp)
		}
	}

	return recoveredRecords, nil
}

// shouldRecoverPolicy 判断策略规则是否应该恢复
func (p *EventProcessor) shouldRecoverPolicy(rule *model.AlertPolicyRule, record *model.AlertRecord) bool {
	// 检查恢复时间：告警最后更新时间距今超过RecoverDuration，说明已持续正常
	// 这样可以避免指标短暂恢复后又立即触发的情况
	if time.Since(record.UpdatedAt.Time) < time.Duration(rule.RecoverDuration)*time.Second {
		log.Printf("告警未达到恢复时间 [规则:%s, 上次更新:%v, 需要持续:%ds]",
			rule.RuleName, time.Since(record.UpdatedAt.Time).Seconds(), rule.RecoverDuration)
		return false
	}

	return true
}

// ForceRecoverRuleAlerts 强制恢复指定规则的所有告警（用于规则禁用时）
func (p *EventProcessor) ForceRecoverRuleAlerts(policy *model.AlertPolicy, rule *model.AlertPolicyRule) {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf("🔕 强制恢复规则告警 [策略:%s, 规则:%s]", policy.Name, rule.RuleName)

	// 检查所有正在触发的告警
	for fingerprint, record := range p.firingEvents {
		// 检查是否属于该规则
		if record.RuleID == rule.ID && record.PolicyID == policy.ID {
			log.Printf("   └─ 恢复告警 [目标:%s, 记录ID:%d]", record.TargetName, record.ID)

			// 直接恢复，不需要等待RecoverDuration
			if err := p.recoverPolicyAlert(policy, rule, record, fingerprint); err != nil {
				log.Printf("      ❌ 恢复失败: %v", err)
			} else {
				log.Printf("      ✅ 已恢复")
			}
		}
	}

	// 清理该规则的连续触发计数器
	// 暂时清理所有（简化实现）
	for fp := range p.continuousCounter {
		delete(p.continuousCounter, fp)
	}

	// 清理该规则的恢复跟踪器
	for fp := range p.recoveryTracker {
		delete(p.recoveryTracker, fp)
	}
}

// recoverPolicyAlert 恢复策略告警
func (p *EventProcessor) recoverPolicyAlert(policy *model.AlertPolicy, rule *model.AlertPolicyRule, record *model.AlertRecord, fingerprint string) error {
	now := util.HTime{Time: time.Now()}
	record.Status = "resolved"
	record.EndTime = &now
	record.Duration = int(time.Since(record.StartTime.Time).Seconds())

	// 保存到数据库
	if err := p.recordDao.Update(record); err != nil {
		return fmt.Errorf("更新告警记录失败: %v", err)
	}

	// 清空已恢复告警的 fingerprint，避免下次触发时冲突
	if err := p.recordDao.ClearFingerprintForResolved(fingerprint); err != nil {
		log.Printf("清空已恢复告警的fingerprint失败: %v", err)
	}

	// 从内存映射中移除
	delete(p.firingEvents, fingerprint)

	// 检查规则下是否还有其他告警
	hasOtherAlerts := false
	for _, r := range p.firingEvents {
		if r.RuleID == rule.ID && r.PolicyID == policy.ID {
			hasOtherAlerts = true
			break
		}
	}

	// 如果该规则下没有其他告警，更新规则健康状态为健康
	if !hasOtherAlerts {
		if err := p.policyRuleDao.UpdateHealthStatus(rule.ID, "healthy"); err != nil {
			log.Printf("更新规则健康状态失败: %v", err)
		}
	}

	// 检查策略下是否还有其他告警
	hasPolicyAlerts := false
	for _, r := range p.firingEvents {
		if r.PolicyID == policy.ID {
			hasPolicyAlerts = true
			break
		}
	}

	// 如果该策略下没有其他告警，更新策略健康状态为健康
	if !hasPolicyAlerts {
		if err := p.policyDao.UpdateHealthStatus(policy.ID, "healthy"); err != nil {
			log.Printf("更新策略健康状态失败: %v", err)
		}
	}

	log.Printf("告警恢复 [策略:%s, 规则:%s, fingerprint:%s, 持续:%ds]", policy.Name, rule.RuleName, fingerprint, record.Duration)

	return nil
}

// getTargetNameFromLabels 从标签中获取目标名称（格式：主机名 (IP)）
func (p *EventProcessor) getTargetNameFromLabels(labels map[string]string) string {
	hostname := ""
	instance := ""

	// 获取hostname
	if h, ok := labels["hostname"]; ok && h != "" {
		hostname = h
	}

	// 获取instance (IP地址)
	if i, ok := labels["instance"]; ok && i != "" {
		instance = i
	}

	// 组合格式：主机名 (IP)
	if hostname != "" && instance != "" {
		return fmt.Sprintf("%s (%s)", hostname, instance)
	} else if hostname != "" {
		return hostname
	} else if instance != "" {
		return instance
	}

	return "unknown"
}

// getTargetIDFromLabels 从标签中获取目标ID
func (p *EventProcessor) getTargetIDFromLabels(labels map[string]string) uint {
	// 从instance标签获取IP，再查询主机ID
	instance, ok := labels["instance"]
	if !ok {
		return 0
	}

	host, err := p.hostDao.GetCmdbHostByIP(instance)
	if err != nil {
		return 0
	}

	return host.ID
}

// buildPolicyAlertMessage 构建策略告警消息
func (p *EventProcessor) buildPolicyAlertMessage(policy *model.AlertPolicy, rule *model.AlertPolicyRule, point *eval.AnomalyPoint, targetName string) string {
	severityMap := map[int]string{
		1: "严重",
		2: "警告",
		3: "提示",
	}

	// 级别对应的火焰emoji
	severityEmoji := map[int]string{
		1: "🔥🔥🔥",
		2: "🔥🔥",
		3: "🔥",
	}

	severity := severityMap[rule.Severity]
	if severity == "" {
		severity = "未知"
	}

	emoji := severityEmoji[rule.Severity]
	if emoji == "" {
		emoji = "🔥"
	}

	return fmt.Sprintf(
		"【%s】%s %s\n"+
			"【策略】: %s\n"+
			"【规则】: %s\n"+
			"【目标】: %s\n"+
			"【指标】: %s\n"+
			"【运算符】: %s\n"+
			"【阈值】: %.2f%s\n"+
			"【当前值】: %.2f%s\n"+
			"【时间】: %s\n"+
			"【描述】: %s",
		severity,
		rule.RuleName,
		emoji,
		policy.Name,
		rule.RuleName,
		targetName,
		rule.MetricKey,
		rule.Operator,
		rule.Threshold,
		rule.Unit,
		point.Value,
		rule.Unit,
		time.Now().Format("2006-01-02 15:04:05"),
		rule.Description,
	)
}
