package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// WebhookService Webhook服务接口
type WebhookService interface {
	ReceiveWebhook(c *gin.Context, req *model.WebhookRequest)
	GetLogByID(c *gin.Context, id uint)
	GetLogList(c *gin.Context, source, status string, page, pageSize int)
	GetNotifyLogs(c *gin.Context, webhookLogID uint)
}

// WebhookServiceImpl Webhook服务实现
type WebhookServiceImpl struct {
	webhookDao dao.WebhookDao
	robotDao   dao.NotifyRobotDao
}

// NewWebhookService 创建Webhook服务实例
func NewWebhookService() WebhookService {
	return &WebhookServiceImpl{
		webhookDao: dao.NewWebhookDao(),
		robotDao:   dao.NewNotifyRobotDao(),
	}
}

// ReceiveWebhook 接收Webhook请求并转发通知
func (s *WebhookServiceImpl) ReceiveWebhook(c *gin.Context, req *model.WebhookRequest) {
	// 创建Webhook日志
	webhookLog := &model.WebhookLog{
		Source:  req.Source,
		Title:   req.Title,
		Content: req.Content,
		Level:   req.Level,
		Status:  "success",
	}

	// 转换Tags为JSON字符串
	if len(req.Tags) > 0 {
		tagsJSON, _ := json.Marshal(req.Tags)
		webhookLog.Tags = string(tagsJSON)
	}

	// 转换Extra为JSON字符串
	if len(req.Extra) > 0 {
		extraJSON, _ := json.Marshal(req.Extra)
		webhookLog.Extra = string(extraJSON)
	}

	// 转换NotifyRobotIDs为字符串
	if len(req.NotifyRobotIDs) > 0 {
		robotIDStrs := make([]string, len(req.NotifyRobotIDs))
		for i, id := range req.NotifyRobotIDs {
			robotIDStrs[i] = fmt.Sprintf("%d", id)
		}
		webhookLog.NotifyRobotIDs = strings.Join(robotIDStrs, ",")
	}

	// 获取需要通知的机器人列表
	var robots []model.NotifyRobot
	var err error

	if len(req.NotifyRobotIDs) > 0 {
		// 如果指定了机器人ID,则只通知指定的机器人
		for _, robotID := range req.NotifyRobotIDs {
			robot, robotErr := s.robotDao.GetByID(robotID)
			if robotErr == nil && robot.Status == 1 {
				robots = append(robots, *robot)
			}
		}
	} else {
		// 否则通知所有启用的机器人
		robots, err = s.robotDao.GetEnabledRobots()
		if err != nil {
			webhookLog.Status = "failed"
			webhookLog.ErrorMsg = fmt.Sprintf("获取通知机器人失败: %v", err)
			s.webhookDao.CreateLog(webhookLog)
			result.Failed(c, int(result.ApiCode.FAILED), webhookLog.ErrorMsg)
			return
		}
	}

	if len(robots) == 0 {
		webhookLog.Status = "failed"
		webhookLog.ErrorMsg = "没有可用的通知机器人"
		s.webhookDao.CreateLog(webhookLog)
		result.Failed(c, int(result.ApiCode.FAILED), webhookLog.ErrorMsg)
		return
	}

	webhookLog.NotifyCount = len(robots)

	// 先保存Webhook日志以获取ID
	if err := s.webhookDao.CreateLog(webhookLog); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建Webhook日志失败: %v", err))
		return
	}

	// 向每个机器人发送通知
	var successCount, failedCount int
	var failedErrors []string

	for _, robot := range robots {
		notifyLog := &model.WebhookNotifyLog{
			WebhookLogID: webhookLog.ID,
			RobotID:      robot.ID,
			RobotName:    robot.Name,
			RobotType:    robot.Type,
			Status:       "success",
		}

		// 根据机器人类型发送消息
		var sendErr error
		switch robot.Type {
		case "dingtalk":
			sendErr = s.sendDingTalkNotify(&robot, req)
		case "feishu":
			sendErr = s.sendFeishuNotify(&robot, req)
		case "wechat":
			sendErr = s.sendWechatNotify(&robot, req)
		case "email":
			sendErr = s.sendEmailNotify(&robot, req)
		case "webhook":
			sendErr = s.sendWebhookNotify(&robot, req)
		case "teams":
			sendErr = s.sendTeamsNotify(&robot, req)
		default:
			sendErr = fmt.Errorf("不支持的机器人类型: %s", robot.Type)
		}

		if sendErr != nil {
			notifyLog.Status = "failed"
			notifyLog.ErrorMsg = sendErr.Error()
			failedCount++
			failedErrors = append(failedErrors, fmt.Sprintf("[%s]%s", robot.Name, sendErr.Error()))
		} else {
			successCount++
		}

		// 保存通知日志
		s.webhookDao.CreateNotifyLog(notifyLog)
	}

	// 更新Webhook日志统计
	webhookLog.SuccessCount = successCount
	webhookLog.FailedCount = failedCount
	if failedCount > 0 {
		if successCount > 0 {
			webhookLog.Status = "partial"
		} else {
			webhookLog.Status = "failed"
		}
		webhookLog.ErrorMsg = strings.Join(failedErrors, "; ")
	}

	// 更新日志状态
	s.webhookDao.CreateLog(webhookLog)

	result.Success(c, gin.H{
		"webhookLogId": webhookLog.ID,
		"notifyCount":  webhookLog.NotifyCount,
		"successCount": successCount,
		"failedCount":  failedCount,
		"status":       webhookLog.Status,
	})
}

// GetLogByID 根据ID获取Webhook日志
func (s *WebhookServiceImpl) GetLogByID(c *gin.Context, id uint) {
	log, err := s.webhookDao.GetLogByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		return
	}

	result.Success(c, log)
}

// GetLogList 获取Webhook日志列表
func (s *WebhookServiceImpl) GetLogList(c *gin.Context, source, status string, page, pageSize int) {
	// 设置默认分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	logs, total, err := s.webhookDao.GetLogList(source, status, page, pageSize)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		return
	}

	result.Success(c, gin.H{
		"list":     logs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetNotifyLogs 获取Webhook通知日志列表
func (s *WebhookServiceImpl) GetNotifyLogs(c *gin.Context, webhookLogID uint) {
	logs, err := s.webhookDao.GetNotifyLogsByWebhookID(webhookLogID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		return
	}

	result.Success(c, logs)
}

// sendDingTalkNotify 发送钉钉通知
func (s *WebhookServiceImpl) sendDingTalkNotify(robot *model.NotifyRobot, req *model.WebhookRequest) error {
	timestamp := time.Now().UnixMilli()
	sign := ""

	// 如果配置了密钥，生成签名
	if robot.Secret != "" {
		sign = generateDingTalkSign(timestamp, robot.Secret)
	}

	// 构建请求URL
	url := robot.Webhook
	if sign != "" {
		url = fmt.Sprintf("%s&timestamp=%d&sign=%s", robot.Webhook, timestamp, sign)
	}

	// 构建消息内容
	content := formatNotifyContent(req)

	// 构建消息体
	message := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
	}

	return sendHTTPRequest(url, message)
}

// sendFeishuNotify 发送飞书通知
func (s *WebhookServiceImpl) sendFeishuNotify(robot *model.NotifyRobot, req *model.WebhookRequest) error {
	// 构建消息内容
	content := formatNotifyContent(req)

	// 构建消息体
	message := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": content,
		},
	}

	return sendHTTPRequest(robot.Webhook, message)
}

// sendWechatNotify 发送企业微信通知
func (s *WebhookServiceImpl) sendWechatNotify(robot *model.NotifyRobot, req *model.WebhookRequest) error {
	// 构建消息内容
	content := formatNotifyContent(req)

	// 构建消息体
	message := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
	}

	return sendHTTPRequest(robot.Webhook, message)
}

// sendEmailNotify 发送邮件通知
func (s *WebhookServiceImpl) sendEmailNotify(robot *model.NotifyRobot, req *model.WebhookRequest) error {
	// 这里可以复用 NotifyRobotService 中的邮件发送逻辑
	// 为了简化,这里返回一个待实现的错误
	// TODO: 实现邮件发送
	return fmt.Errorf("邮件通知暂未实现")
}

// sendWebhookNotify 发送通用Webhook通知
func (s *WebhookServiceImpl) sendWebhookNotify(robot *model.NotifyRobot, req *model.WebhookRequest) error {
	// 构建默认消息内容
	defaultContent := formatNotifyContent(req)

	message := map[string]interface{}{
		"source":  req.Source,
		"title":   req.Title,
		"content": req.Content,
		"level":   req.Level,
		"tags":    req.Tags,
		"extra":   req.Extra,
		"time":    time.Now().Format("2006-01-02 15:04:05"),
		"text":    defaultContent, // 格式化后的文本内容
	}

	// 如果配置了自定义模板，使用模板渲染
	if robot.Template != "" {
		var customMessage map[string]interface{}
		if err := json.Unmarshal([]byte(robot.Template), &customMessage); err == nil {
			// 模板变量替换
			messageJSON, _ := json.Marshal(customMessage)
			messageStr := string(messageJSON)
			messageStr = strings.ReplaceAll(messageStr, "{{source}}", req.Source)
			messageStr = strings.ReplaceAll(messageStr, "{{title}}", req.Title)
			messageStr = strings.ReplaceAll(messageStr, "{{content}}", req.Content)
			messageStr = strings.ReplaceAll(messageStr, "{{level}}", req.Level)
			messageStr = strings.ReplaceAll(messageStr, "{{time}}", time.Now().Format("2006-01-02 15:04:05"))
			json.Unmarshal([]byte(messageStr), &message)
		}
	}

	return sendCustomHTTPRequest(robot.Webhook, robot.Method, robot.Headers, message)
}

// sendTeamsNotify 发送Microsoft Teams通知
func (s *WebhookServiceImpl) sendTeamsNotify(robot *model.NotifyRobot, req *model.WebhookRequest) error {
	// 根据告警级别设置颜色
	themeColor := "0078D7" // 默认蓝色
	switch req.Level {
	case "warning":
		themeColor = "FFA500" // 橙色
	case "error":
		themeColor = "FF0000" // 红色
	case "critical":
		themeColor = "8B0000" // 深红色
	case "info":
		themeColor = "0078D7" // 蓝色
	}

	// 构建facts
	facts := []map[string]string{
		{"name": "来源", "value": req.Source},
		{"name": "级别", "value": req.Level},
		{"name": "内容", "value": req.Content},
		{"name": "时间", "value": time.Now().Format("2006-01-02 15:04:05")},
	}

	if len(req.Tags) > 0 {
		facts = append(facts, map[string]string{
			"name":  "标签",
			"value": strings.Join(req.Tags, ", "),
		})
	}

	// Teams MessageCard格式
	message := map[string]interface{}{
		"@type":      "MessageCard",
		"@context":   "https://schema.org/extensions",
		"summary":    req.Title,
		"themeColor": themeColor,
		"title":      req.Title,
		"sections": []map[string]interface{}{
			{
				"activityTitle": fmt.Sprintf("来自 %s 的告警", req.Source),
				"facts":         facts,
			},
		},
	}

	// 如果配置了自定义模板，使用模板
	if robot.Template != "" {
		var customMessage map[string]interface{}
		if err := json.Unmarshal([]byte(robot.Template), &customMessage); err == nil {
			// 模板变量替换
			messageJSON, _ := json.Marshal(customMessage)
			messageStr := string(messageJSON)
			messageStr = strings.ReplaceAll(messageStr, "{{source}}", req.Source)
			messageStr = strings.ReplaceAll(messageStr, "{{title}}", req.Title)
			messageStr = strings.ReplaceAll(messageStr, "{{content}}", req.Content)
			messageStr = strings.ReplaceAll(messageStr, "{{level}}", req.Level)
			messageStr = strings.ReplaceAll(messageStr, "{{time}}", time.Now().Format("2006-01-02 15:04:05"))
			json.Unmarshal([]byte(messageStr), &message)
		}
	}

	return sendCustomHTTPRequest(robot.Webhook, "POST", robot.Headers, message)
}

// formatNotifyContent 格式化通知内容
func formatNotifyContent(req *model.WebhookRequest) string {
	var buf bytes.Buffer

	// 添加告警级别标识
	levelIcon := "ℹ️"
	switch req.Level {
	case "warning":
		levelIcon = "⚠️"
	case "error":
		levelIcon = "❌"
	case "critical":
		levelIcon = "🔥"
	}

	buf.WriteString(fmt.Sprintf("%s 【%s】\n", levelIcon, req.Title))
	buf.WriteString(fmt.Sprintf("来源: %s\n", req.Source))
	if req.Level != "" {
		buf.WriteString(fmt.Sprintf("级别: %s\n", req.Level))
	}
	buf.WriteString(fmt.Sprintf("内容: %s\n", req.Content))

	if len(req.Tags) > 0 {
		buf.WriteString(fmt.Sprintf("标签: %s\n", strings.Join(req.Tags, ", ")))
	}

	buf.WriteString(fmt.Sprintf("时间: %s", time.Now().Format("2006-01-02 15:04:05")))

	return buf.String()
}

// generateDingTalkSign 生成钉钉签名
func generateDingTalkSign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// sendHTTPRequest 发送HTTP请求
func sendHTTPRequest(url string, message map[string]interface{}) error {
	jsonData, err := json.Marshal(message)
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

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP状态码异常: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应检查是否成功
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查各平台的响应格式
	if errcode, ok := response["errcode"].(float64); ok && errcode != 0 {
		return fmt.Errorf("发送失败，错误码: %.0f, 错误信息: %v", errcode, response["errmsg"])
	}
	if code, ok := response["code"].(float64); ok && code != 0 {
		return fmt.Errorf("发送失败，错误码: %.0f, 错误信息: %v", code, response["msg"])
	}

	return nil
}

// sendCustomHTTPRequest 发送自定义HTTP请求（支持自定义方法和Headers）
func sendCustomHTTPRequest(url, method, headersJSON string, message map[string]interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 设置默认方法
	if method == "" {
		method = "POST"
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置默认Content-Type
	req.Header.Set("Content-Type", "application/json")

	// 解析并设置自定义Headers
	if headersJSON != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(headersJSON), &headers); err == nil {
			for key, value := range headers {
				req.Header.Set(key, value)
			}
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Teams和通用Webhook可能返回2xx系列状态码都算成功
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP状态码异常: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return nil
}
