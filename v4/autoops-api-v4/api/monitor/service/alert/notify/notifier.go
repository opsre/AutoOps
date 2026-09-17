package notify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service/alert/eval"
)

// AlertNotifier 告警通知器
type AlertNotifier struct {
	robotDao  dao.NotifyRobotDao
	recordDao dao.AlertRecordDao
}

// NewAlertNotifier 创建告警通知器
func NewAlertNotifier() *AlertNotifier {
	return &AlertNotifier{
		robotDao:  dao.NewNotifyRobotDao(),
		recordDao: dao.NewAlertRecordDao(),
	}
}

// Notify 发送告警通知
func (n *AlertNotifier) Notify(rule *model.AlertRule, points []*eval.AnomalyPoint) error {
	// 解析机器人ID列表
	var robotIDs []uint
	if err := json.Unmarshal([]byte(rule.RobotIDs), &robotIDs); err != nil {
		return fmt.Errorf("解析机器人ID失败: %v", err)
	}

	if len(robotIDs) == 0 {
		return nil // 没有配置机器人,跳过通知
	}

	// 遍历每个异常点
	for _, point := range points {
		// 生成指纹
		fingerprint := eval.GenerateFingerprint(rule.ID, point.Labels)

		// 获取告警记录
		record, err := n.recordDao.GetByFingerprint(fingerprint)
		if err != nil {
			log.Printf("获取告警记录失败: %v", err)
			continue
		}

		// 检查静默期
		if n.isInSilencePeriod(rule, record) {
			log.Printf("告警在静默期内,跳过通知 [规则:%s, fingerprint:%s]", rule.Name, fingerprint)
			continue
		}

		// 检查重复通知限制
		if n.exceedsMaxNotifyCount(rule, record) {
			log.Printf("告警已达到最大通知次数,跳过通知 [规则:%s, fingerprint:%s]", rule.Name, fingerprint)
			continue
		}

		// 发送到每个机器人
		for _, robotID := range robotIDs {
			if err := n.sendToRobot(robotID, record); err != nil {
				log.Printf("发送告警到机器人失败 [机器人ID:%d]: %v", robotID, err)
			}
		}

		// 更新通知状态
		now := time.Now()
		n.recordDao.UpdateNotifyInfo(record.ID, now, record.NotifyCount+1, "")
	}

	return nil
}

// sendToRobot 发送到指定机器人
func (n *AlertNotifier) sendToRobot(robotID uint, record *model.AlertRecord) error {
	// 获取机器人配置
	robot, err := n.robotDao.GetByID(robotID)
	if err != nil {
		return fmt.Errorf("获取机器人配置失败: %v", err)
	}

	if robot.Status != 1 {
		return fmt.Errorf("机器人已禁用")
	}

	// 渲染消息内容（支持模板）
	message, err := n.renderMessage(robot, record)
	if err != nil {
		return fmt.Errorf("渲染消息失败: %v", err)
	}

	// 根据类型发送消息
	switch robot.Type {
	case "dingtalk":
		return n.sendDingTalk(robot, message)
	case "feishu":
		return n.sendFeishu(robot, message)
	case "wechat":
		return n.sendWechat(robot, message)
	case "webhook":
		return n.sendWebhook(robot, record)
	case "teams":
		return n.sendTeams(robot, record)
	default:
		return fmt.Errorf("不支持的机器人类型: %s", robot.Type)
	}
}

// renderMessage 渲染消息模板
func (n *AlertNotifier) renderMessage(robot *model.NotifyRobot, record *model.AlertRecord) (string, error) {
	// 如果没有配置模板，使用默认消息
	if robot.Template == "" {
		return record.Message, nil
	}

	// 解析标签
	var labels map[string]string
	if record.Labels != "" {
		if err := json.Unmarshal([]byte(record.Labels), &labels); err != nil {
			log.Printf("解析标签失败: %v", err)
			labels = make(map[string]string)
		}
	} else {
		labels = make(map[string]string)
	}

	// 标签字段映射（支持多种命名方式）
	n.normalizeLabels(labels)

	// 构建告警标题：策略名称-规则描述
	alertTitle := fmt.Sprintf("%s-%s", record.PolicyName, getShortDescription(record))

	// 构建模板数据
	data := map[string]interface{}{
		"IsRecovered":               record.Status == "resolved",
		"alert_title":               alertTitle,       // 告警标题
		"target_name":               record.TargetName, // 目标名称（含IP）
		"rule_name":                 record.RuleName,
		"policy_name":               record.PolicyName,
		"fingerprint":               record.Fingerprint,
		"severity":                  n.getSeverityName(record.Severity),
		"severity_level":            record.Severity,
		"metric_name":               record.MetricName,
		"metric_expression":         n.buildMetricExpression(record, labels), // 指标表达式（从labels获取operator）
		"current_value":             fmt.Sprintf("%.2f", record.CurrentValue),
		"threshold":                 fmt.Sprintf("%.2f", record.Threshold),
		"first_trigger_time":        record.StartTime,
		"first_trigger_time_format": record.StartTime.Format("2006-01-02 15:04:05"),
		"duration":                  record.Duration,
		"duration_format":           n.formatDuration(record.Duration),
		"group_name":                record.GroupName,
		"message":                   record.Message,
		"annotations":               record.Message,
		"labels":                    labels,
		"duty_user":                 "暂无", // TODO: 可以从其他系统获取值班人员
		"faultCenterId":             record.ID,
	}

	// 解析并执行模板
	tmpl, err := template.New("alert").Parse(robot.Template)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %v", err)
	}

	message := buf.String()

	// 替换 ${xxx} 格式的变量（兼容旧的模板格式）
	message = n.replaceVariables(message, data, labels)

	return message, nil
}

// normalizeLabels 标准化标签（兼容不同的命名方式）
func (n *AlertNotifier) normalizeLabels(labels map[string]string) {
	// 如果没有 node_name，尝试从 hostname 或 instance 映射
	if _, exists := labels["node_name"]; !exists {
		if hostname, ok := labels["hostname"]; ok {
			labels["node_name"] = hostname
		} else if instance, ok := labels["instance"]; ok {
			labels["node_name"] = instance
		}
	}

	// 如果没有 pod，尝试从其他字段映射
	if _, exists := labels["pod"]; !exists {
		if podName, ok := labels["pod_name"]; ok {
			labels["pod"] = podName
		} else {
			labels["pod"] = "-"
		}
	}

	// 如果没有 namespace，设置默认值
	if _, exists := labels["namespace"]; !exists {
		labels["namespace"] = "-"
	}
}

// replaceVariables 替换 ${xxx} 格式的变量
func (n *AlertNotifier) replaceVariables(message string, data map[string]interface{}, labels map[string]string) string {
	// 替换简单变量 ${key}
	for key, value := range data {
		if key == "labels" {
			continue
		}
		placeholder := fmt.Sprintf("${%s}", key)
		message = strings.ReplaceAll(message, placeholder, fmt.Sprintf("%v", value))
	}

	// 替换标签变量 ${labels.xxx}
	for key, value := range labels {
		placeholder := fmt.Sprintf("${labels.%s}", key)
		message = strings.ReplaceAll(message, placeholder, value)
	}

	return message
}

// buildMetricExpression 构建指标表达式（格式：指标名 运算符 阈值）
func (n *AlertNotifier) buildMetricExpression(record *model.AlertRecord, labels map[string]string) string {
	// 尝试从 labels 中获取 operator
	operator := labels["operator"]

	// 如果 labels 中没有，尝试从 Message 中提取
	if operator == "" {
		// 从消息中提取 "阈值: 运算符 值" 这一行
		lines := strings.Split(record.Message, "\n")
		for _, line := range lines {
			if strings.Contains(line, "阈值:") || strings.Contains(line, "阈值：") {
				// 提取运算符（>、<、>=、<=、==、!=）
				for _, op := range []string{">=", "<=", "==", "!=", ">", "<"} {
					if strings.Contains(line, op) {
						operator = op
						break
					}
				}
				break
			}
		}
	}

	// 如果还是没有找到，使用默认运算符
	if operator == "" {
		operator = ">"
	}

	// 格式化：指标名 运算符 阈值
	return fmt.Sprintf("%s %s %.2f", record.MetricName, operator, record.Threshold)
}

// getSeverityName 获取告警级别名称
func (n *AlertNotifier) getSeverityName(severity int) string {
	severityMap := map[int]string{
		1: "严重",
		2: "警告",
		3: "提示",
	}
	if name, ok := severityMap[severity]; ok {
		return name
	}
	return "未知"
}

// formatDuration 格式化持续时间
func (n *AlertNotifier) formatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%d秒", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%d分钟", seconds/60)
	} else if seconds < 86400 {
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		if minutes > 0 {
			return fmt.Sprintf("%d小时%d分钟", hours, minutes)
		}
		return fmt.Sprintf("%d小时", hours)
	} else {
		days := seconds / 86400
		hours := (seconds % 86400) / 3600
		if hours > 0 {
			return fmt.Sprintf("%d天%d小时", days, hours)
		}
		return fmt.Sprintf("%d天", days)
	}
}

// getShortDescription 从告警消息中提取简短描述
func getShortDescription(record *model.AlertRecord) string {
	// 尝试从消息中提取描述行（格式：描述: xxx）
	lines := strings.Split(record.Message, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "描述:") || strings.HasPrefix(line, "描述：") {
			desc := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "描述:"), "描述："))
			if desc != "" {
				return desc
			}
		}
	}

	// 如果没找到描述，返回规则名称
	return record.RuleName
}

// sendDingTalk 发送钉钉消息
func (n *AlertNotifier) sendDingTalk(robot *model.NotifyRobot, message string) error {
	timestamp := time.Now().UnixMilli()
	sign := ""

	// 如果配置了密钥,生成签名
	if robot.Secret != "" {
		sign = n.generateDingTalkSign(timestamp, robot.Secret)
	}

	// 构建请求URL
	url := robot.Webhook
	if sign != "" {
		url = fmt.Sprintf("%s&timestamp=%d&sign=%s", robot.Webhook, timestamp, sign)
	}

	// 构建消息体
	body := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
	}

	return n.sendHTTPRequest(url, body)
}

// generateDingTalkSign 生成钉钉签名
func (n *AlertNotifier) generateDingTalkSign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// sendFeishu 发送飞书消息
func (n *AlertNotifier) sendFeishu(robot *model.NotifyRobot, message string) error {
	body := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": message,
		},
	}

	return n.sendHTTPRequest(robot.Webhook, body)
}

// sendWechat 发送企业微信消息
func (n *AlertNotifier) sendWechat(robot *model.NotifyRobot, message string) error {
	body := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
	}

	return n.sendHTTPRequest(robot.Webhook, body)
}

// sendWebhook 发送通用Webhook
func (n *AlertNotifier) sendWebhook(robot *model.NotifyRobot, record *model.AlertRecord) error {
	// 构建消息体
	body := map[string]interface{}{
		"rule_id":       record.RuleID,
		"rule_name":     record.RuleName,
		"target_name":   record.TargetName,
		"metric_name":   record.MetricName,
		"current_value": record.CurrentValue,
		"threshold":     record.Threshold,
		"severity":      record.Severity,
		"status":        record.Status,
		"message":       record.Message,
		"start_time":    record.StartTime,
	}

	return n.sendHTTPRequest(robot.Webhook, body)
}

// sendTeams 发送Teams消息
func (n *AlertNotifier) sendTeams(robot *model.NotifyRobot, record *model.AlertRecord) error {
	severityColors := map[int]string{
		1: "FF0000", // 严重 - 红色
		2: "FFA500", // 警告 - 橙色
		3: "0078D7", // 提示 - 蓝色
	}

	body := map[string]interface{}{
		"@type":      "MessageCard",
		"@context":   "https://schema.org/extensions",
		"summary":    fmt.Sprintf("告警: %s", record.RuleName),
		"themeColor": severityColors[record.Severity],
		"title":      fmt.Sprintf("🔔 %s", record.RuleName),
		"sections": []map[string]interface{}{
			{
				"facts": []map[string]string{
					{"name": "目标", "value": record.TargetName},
					{"name": "指标", "value": record.MetricName},
					{"name": "当前值", "value": fmt.Sprintf("%.2f", record.CurrentValue)},
					{"name": "阈值", "value": fmt.Sprintf("%.2f", record.Threshold)},
					{"name": "时间", "value": record.StartTime.Format("2006-01-02 15:04:05")},
				},
			},
		},
	}

	return n.sendHTTPRequest(robot.Webhook, body)
}

// sendHTTPRequest 发送HTTP请求
func (n *AlertNotifier) sendHTTPRequest(url string, body map[string]interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP状态码异常: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 检查响应
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err == nil {
		if errcode, ok := response["errcode"].(float64); ok && errcode != 0 {
			return fmt.Errorf("发送失败，错误码: %.0f, 错误信息: %v", errcode, response["errmsg"])
		}
		if code, ok := response["code"].(float64); ok && code != 0 {
			return fmt.Errorf("发送失败，错误码: %.0f, 错误信息: %v", code, response["msg"])
		}
	}

	return nil
}

// isInSilencePeriod 检查是否在静默期内
func (n *AlertNotifier) isInSilencePeriod(rule *model.AlertRule, record *model.AlertRecord) bool {
	if record.NotifyTime == nil {
		return false
	}

	elapsed := time.Since(record.NotifyTime.Time).Seconds()
	return elapsed < float64(rule.SilencePeriod)
}

// exceedsMaxNotifyCount 检查是否超过最大通知次数
func (n *AlertNotifier) exceedsMaxNotifyCount(rule *model.AlertRule, record *model.AlertRecord) bool {
	if rule.MaxNotifyCount == 0 {
		return false // 0表示不限制
	}

	return record.NotifyCount >= rule.MaxNotifyCount
}

// NotifyPolicyRule 发送告警策略规则通知
func (n *AlertNotifier) NotifyPolicyRule(policy *model.AlertPolicy, rule *model.AlertPolicyRule, point *eval.AnomalyPoint) error {
	// 解析机器人ID列表
	var robotIDs []uint
	if err := json.Unmarshal([]byte(policy.RobotIDs), &robotIDs); err != nil {
		return fmt.Errorf("解析机器人ID失败: %v", err)
	}

	if len(robotIDs) == 0 {
		log.Printf("未配置通知机器人，跳过通知 [策略:%s, 规则:%s]", policy.Name, rule.RuleName)
		return nil // 没有配置机器人,跳过通知
	}

	// 生成指纹
	fingerprint := eval.GenerateFingerprint(rule.ID, point.Labels)

	// 获取告警记录
	record, err := n.recordDao.GetByFingerprint(fingerprint)
	if err != nil {
		return fmt.Errorf("获取告警记录失败: %v", err)
	}

	log.Printf("告警记录状态 [策略:%s, 规则:%s, 记录ID:%d, 上次通知时间:%v, 通知次数:%d, 静默期:%ds, 重复间隔:%ds, 最大通知次数:%d]",
		policy.Name, rule.RuleName, record.ID, record.NotifyTime, record.NotifyCount, policy.SilencePeriod, policy.RepeatInterval, policy.MaxNotifyCount)

	// 检查重复通知限制（先检查次数上限）
	if n.exceedsMaxNotifyCountPolicy(policy, record) {
		log.Printf("⛔ 告警已达到最大通知次数,跳过通知 [策略:%s, 规则:%s, fingerprint:%s, 最大次数:%d, 当前次数:%d]",
			policy.Name, rule.RuleName, fingerprint, policy.MaxNotifyCount, record.NotifyCount)
		return nil
	}

	// 检查静默期
	if n.isInSilencePeriodPolicy(policy, record) {
		log.Printf("⏰ 告警在静默期内,跳过通知 [策略:%s, 规则:%s, fingerprint:%s]",
			policy.Name, rule.RuleName, fingerprint)
		return nil
	}

	log.Printf("准备发送告警通知 [策略:%s, 规则:%s, 机器人数:%d, 通知次数:%d]",
		policy.Name, rule.RuleName, len(robotIDs), record.NotifyCount)

	// 发送到每个机器人
	successCount := 0
	for _, robotID := range robotIDs {
		log.Printf("🤖 正在发送到机器人 [机器人ID:%d]", robotID)
		if err := n.sendToRobot(robotID, record); err != nil {
			log.Printf("❌ 发送告警到机器人失败 [机器人ID:%d]: %v", robotID, err)
		} else {
			successCount++
			log.Printf("✅ 发送告警到机器人成功 [机器人ID:%d]", robotID)
		}
	}

	log.Printf("📊 发送结果统计 [策略:%s, 规则:%s, 成功:%d/%d]",
		policy.Name, rule.RuleName, successCount, len(robotIDs))

	if successCount > 0 {
		// 更新通知状态
		now := time.Now()
		err := n.recordDao.UpdateNotifyInfo(record.ID, now, record.NotifyCount+1, "")
		if err != nil {
			log.Printf("❌ 更新通知状态失败: %v", err)
		} else {
			log.Printf("✅ 告警通知已发送并更新状态 [策略:%s, 规则:%s, 成功:%d/%d, 总通知次数:%d]",
				policy.Name, rule.RuleName, successCount, len(robotIDs), record.NotifyCount+1)
		}
	} else {
		log.Printf("⚠️ 所有机器人发送失败，不更新通知状态")
	}

	return nil
}

// isInSilencePeriodPolicy 检查策略规则是否在静默期内
func (n *AlertNotifier) isInSilencePeriodPolicy(policy *model.AlertPolicy, record *model.AlertRecord) bool {
	if record.NotifyTime == nil {
		log.Printf("✅ 首次通知，不在静默期")
		return false // 首次通知，不在静默期
	}

	elapsed := time.Since(record.NotifyTime.Time).Seconds()
	log.Printf("⏱️ 距离上次通知时间: %.0f秒 (上次通知: %v)", elapsed, record.NotifyTime.Time.Format("2006-01-02 15:04:05"))

	// 如果配置了重复通知间隔
	if policy.RepeatInterval > 0 {
		log.Printf("📋 配置了重复通知间隔: %ds, 当前已过时间: %.0fs", policy.RepeatInterval, elapsed)

		// 距离上次通知时间小于重复间隔，则跳过
		if elapsed < float64(policy.RepeatInterval) {
			log.Printf("❌ 距离上次通知仅%.0f秒，小于重复间隔%ds，跳过通知", elapsed, policy.RepeatInterval)
			return true
		}
		log.Printf("✅ 已超过重复间隔，允许通知")
		return false
	}

	// 如果没有配置重复间隔，使用静默期
	log.Printf("📋 未配置重复间隔，使用静默期: %ds", policy.SilencePeriod)
	if elapsed < float64(policy.SilencePeriod) {
		log.Printf("❌ 在静默期内，跳过通知")
		return true
	}
	log.Printf("✅ 已超过静默期，允许通知")
	return false
}

// exceedsMaxNotifyCountPolicy 检查策略规则是否超过最大通知次数
func (n *AlertNotifier) exceedsMaxNotifyCountPolicy(policy *model.AlertPolicy, record *model.AlertRecord) bool {
	if policy.MaxNotifyCount == 0 {
		return false // 0表示不限制
	}

	return record.NotifyCount >= policy.MaxNotifyCount
}

// NotifyRecovery 发送告警恢复通知
func (n *AlertNotifier) NotifyRecovery(policy *model.AlertPolicy, rule *model.AlertPolicyRule, record *model.AlertRecord) error {
	// 解析机器人ID列表
	var robotIDs []uint
	if err := json.Unmarshal([]byte(policy.RobotIDs), &robotIDs); err != nil {
		return fmt.Errorf("解析机器人ID失败: %v", err)
	}

	if len(robotIDs) == 0 {
		log.Printf("未配置通知机器人，跳过恢复通知 [策略:%s, 规则:%s]", policy.Name, rule.RuleName)
		return nil
	}

	log.Printf("准备发送恢复通知 [策略:%s, 规则:%s, 机器人数:%d]", policy.Name, rule.RuleName, len(robotIDs))

	// 发送到每个机器人
	successCount := 0
	for _, robotID := range robotIDs {
		log.Printf("🤖 正在发送恢复通知到机器人 [机器人ID:%d]", robotID)
		if err := n.sendRecoveryToRobot(robotID, record); err != nil {
			log.Printf("❌ 发送恢复通知到机器人失败 [机器人ID:%d]: %v", robotID, err)
		} else {
			successCount++
			log.Printf("✅ 发送恢复通知到机器人成功 [机器人ID:%d]", robotID)
		}
	}

	log.Printf("恢复通知发送完成 [策略:%s, 规则:%s, 成功:%d/%d]", policy.Name, rule.RuleName, successCount, len(robotIDs))

	return nil
}

// sendRecoveryToRobot 发送恢复通知到指定机器人
func (n *AlertNotifier) sendRecoveryToRobot(robotID uint, record *model.AlertRecord) error {
	// 获取机器人配置
	robot, err := n.robotDao.GetByID(robotID)
	if err != nil {
		return fmt.Errorf("获取机器人配置失败: %v", err)
	}

	if robot.Status != 1 {
		return fmt.Errorf("机器人已禁用")
	}

	// 渲染恢复消息
	message, err := n.renderRecoveryMessage(robot, record)
	if err != nil {
		return fmt.Errorf("渲染恢复消息失败: %v", err)
	}

	// 根据机器人类型发送
	switch robot.Type {
	case "dingtalk":
		return n.sendDingTalk(robot, message)
	case "feishu":
		return n.sendFeishu(robot, message)
	case "wechat":
		return n.sendWechat(robot, message)
	default:
		return fmt.Errorf("不支持的机器人类型: %s", robot.Type)
	}
}

// renderRecoveryMessage 渲染恢复消息
func (n *AlertNotifier) renderRecoveryMessage(robot *model.NotifyRobot, record *model.AlertRecord) (string, error) {
	// 使用默认恢复消息模板
	duration := time.Since(record.StartTime.Time)
	severityText := n.getSeverityName(record.Severity)

	return fmt.Sprintf(`✅ 告警已恢复

【级别】: %s
【策略】: %s
【规则】: %s
【目标】: %s
【指标】: %s
【触发时间】: %s
【恢复时间】: %s
【持续时长】: %s
【描述】: %s`,
		severityText,
		record.PolicyName,
		record.RuleName,
		record.TargetName,
		record.MetricName,
		record.StartTime.Format("2006-01-02 15:04:05"),
		record.EndTime.Format("2006-01-02 15:04:05"),
		formatDuration(duration),
		record.Message,
	), nil
}

// formatDuration 格式化持续时长
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟%d秒", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分钟%d秒", minutes, seconds)
	} else {
		return fmt.Sprintf("%d秒", seconds)
	}
}
