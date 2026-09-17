package eval

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/prometheus"
)

// AnomalyPoint 异常点(告警点)
type AnomalyPoint struct {
	Metric       map[string]string `json:"metric"`        // 指标标签
	Value        float64           `json:"value"`         // 当前值
	Threshold    float64           `json:"threshold"`     // 阈值
	Severity     int               `json:"severity"`      // 告警级别
	Labels       map[string]string `json:"labels"`        // 额外标签
	Annotations  map[string]string `json:"annotations"`   // 注解
	Timestamp    time.Time         `json:"timestamp"`     // 时间戳
	IsRecovered  bool              `json:"is_recovered"`  // 是否恢复
}

// RuleEvaluator 规则评估器
type RuleEvaluator struct {
	promClient *prometheus.Client
	hostDao    dao.CmdbHostDao
}

// NewRuleEvaluator 创建规则评估器
func NewRuleEvaluator(promClient *prometheus.Client) *RuleEvaluator {
	return &RuleEvaluator{
		promClient: promClient,
		hostDao:    dao.NewCmdbHostDao(),
	}
}

// Eval 评估规则
func (e *RuleEvaluator) Eval(rule *model.AlertRule) ([]*AnomalyPoint, error) {
	// 执行PromQL查询
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results, err := e.promClient.Query(ctx, rule.PromQL)
	if err != nil {
		return nil, fmt.Errorf("执行PromQL查询失败: %v", err)
	}

	if len(results) == 0 {
		return nil, nil
	}

	// 解析附加标签
	var appendTags []string
	if rule.AppendTags != "" {
		json.Unmarshal([]byte(rule.AppendTags), &appendTags)
	}

	// 解析注解
	var annotations map[string]string
	if rule.Annotations != "" {
		json.Unmarshal([]byte(rule.Annotations), &annotations)
	}

	// 转换为异常点
	var points []*AnomalyPoint
	for _, result := range results {
		point := &AnomalyPoint{
			Metric:      result.Metric,
			Value:       result.Value,
			Threshold:   rule.Threshold,
			Severity:    rule.Severity,
			Labels:      e.buildLabels(result.Metric, appendTags),
			Annotations: annotations,
			Timestamp:   result.Time,
			IsRecovered: false,
		}

		points = append(points, point)
	}

	return points, nil
}

// buildLabels 构建标签
func (e *RuleEvaluator) buildLabels(metric map[string]string, appendTags []string) map[string]string {
	labels := make(map[string]string)

	// 复制原始标签
	for k, v := range metric {
		labels[k] = v
	}

	// 添加附加标签
	for _, tag := range appendTags {
		// 解析 key=value 格式
		var key, value string
		fmt.Sscanf(tag, "%s=%s", &key, &value)
		if key != "" && value != "" {
			labels[key] = value
		}
	}

	return labels
}

// GenerateFingerprint 生成事件指纹(用于去重)
// 只使用核心标识标签，忽略动态变化的标签
func GenerateFingerprint(ruleID uint, labels map[string]string) string {
	// 定义用于生成fingerprint的核心标签（这些标签唯一标识一个告警对象）
	coreKeys := []string{
		"instance",   // 实例IP/地址
		"hostname",   // 主机名
		"job",        // 任务名
		"device",     // 设备名（磁盘、网卡等）
		"mountpoint", // 挂载点
		"fstype",     // 文件系统类型
		"group_id",   // 业务组ID
	}

	// 提取核心标签
	data := fmt.Sprintf("rule:%d", ruleID)
	for _, key := range coreKeys {
		if value, exists := labels[key]; exists {
			data += fmt.Sprintf(",%s=%s", key, value)
		}
	}

	// 计算MD5
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// BuildAlertMessage 构建告警消息
func BuildAlertMessage(rule *model.AlertRule, point *AnomalyPoint, targetName string) string {
	severityNames := map[int]string{
		1: "严重",
		2: "警告",
		3: "提示",
	}

	severityEmojis := map[int]string{
		1: "🔴",
		2: "⚠️",
		3: "ℹ️",
	}

	// 基础消息
	message := fmt.Sprintf(`%s %s告警

📋 规则名称: %s
🖥  目标名称: %s
📊 指标信息: %s
📈 当前值: %.2f
⚠️  阈值: %.2f
📝 描述: %s
⏰ 触发时间: %s`,
		severityEmojis[rule.Severity],
		severityNames[rule.Severity],
		rule.Name,
		targetName,
		rule.MetricName,
		point.Value,
		rule.Threshold,
		rule.Description,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	// 添加标签信息
	if len(point.Labels) > 0 {
		message += "\n\n🏷  标签:"
		for k, v := range point.Labels {
			if k != "__name__" && k != "job" && k != "instance" {
				message += fmt.Sprintf("\n  - %s: %s", k, v)
			}
		}
	}

	return message
}

// BuildRecoveryMessage 构建恢复消息
func BuildRecoveryMessage(rule *model.AlertRule, point *AnomalyPoint, targetName string, duration int) string {
	message := fmt.Sprintf(`✅ 告警恢复

📋 规则名称: %s
🖥  目标名称: %s
📊 指标信息: %s
📈 当前值: %.2f
⏱  持续时间: %d秒
⏰ 恢复时间: %s`,
		rule.Name,
		targetName,
		rule.MetricName,
		point.Value,
		duration,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return message
}
