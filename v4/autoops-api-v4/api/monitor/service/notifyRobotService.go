package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// NotifyRobotService 通知机器人服务接口（通用接口，支持三种类型）
type NotifyRobotService interface {
	Create(c *gin.Context, robot *model.NotifyRobot, robotType string)
	GetByID(c *gin.Context, id uint, robotType string)
	Update(c *gin.Context, robot *model.NotifyRobot, robotType string)
	Delete(c *gin.Context, id uint, robotType string)
	GetList(c *gin.Context, robotType string, status *int, page, pageSize int)
	TestSend(c *gin.Context, id uint, tempTemplate string) // 测试发送消息（支持临时模板）
}

// NotifyRobotServiceImpl 通知机器人服务实现
type NotifyRobotServiceImpl struct {
	robotDao dao.NotifyRobotDao
}

// NewNotifyRobotService 创建通知机器人服务实例
func NewNotifyRobotService() NotifyRobotService {
	return &NotifyRobotServiceImpl{
		robotDao: dao.NewNotifyRobotDao(),
	}
}

// Create 创建通知机器人
func (s *NotifyRobotServiceImpl) Create(c *gin.Context, robot *model.NotifyRobot, robotType string) {
	// 设置机器人类型
	robot.Type = robotType

	// 验证必填字段
	if robot.Name == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "机器人名称不能为空")
		return
	}

	// 根据类型验证不同的字段
	if robotType == "email" {
		// 邮箱类型验证
		if robot.Server == "" || robot.Port == nil || *robot.Port == 0 || robot.Username == "" || robot.Password == "" {
			result.Failed(c, int(result.ApiCode.FAILED), "邮箱配置不完整，请检查 server/port/username/password 字段")
			return
		}
		// 如果没有填写接收邮箱，默认使用发件邮箱
		if robot.Webhook == "" {
			robot.Webhook = robot.Username
		}
	} else if robotType == "webhook" || robotType == "teams" {
		// Webhook/Teams类型验证
		if robot.Webhook == "" {
			result.Failed(c, int(result.ApiCode.FAILED), "Webhook地址不能为空")
			return
		}
		// 设置默认HTTP方法
		if robot.Method == "" {
			robot.Method = "POST"
		}
	} else {
		// 其他机器人类型验证（钉钉/飞书/企业微信）
		if robot.Webhook == "" {
			result.Failed(c, int(result.ApiCode.FAILED), "Webhook地址不能为空")
			return
		}
	}

	// 设置默认状态为启用
	if robot.Status == 0 {
		robot.Status = 1
	}

	if err := s.robotDao.Create(robot); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建失败: %v", err))
		return
	}

	result.Success(c, robot)
}

// GetByID 根据ID获取通知机器人
func (s *NotifyRobotServiceImpl) GetByID(c *gin.Context, id uint, robotType string) {
	robot, err := s.robotDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		return
	}

	// 验证机器人类型（如果指定了类型）
	if robotType != "" && robot.Type != robotType {
		result.Failed(c, int(result.ApiCode.FAILED), "机器人类型不匹配")
		return
	}

	result.Success(c, robot)
}

// Update 更新通知机器人
func (s *NotifyRobotServiceImpl) Update(c *gin.Context, robot *model.NotifyRobot, robotType string) {
	// 检查机器人是否存在
	existing, err := s.robotDao.GetByID(robot.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("机器人不存在: %v", err))
		return
	}

	// 验证机器人类型
	if existing.Type != robotType {
		result.Failed(c, int(result.ApiCode.FAILED), "机器人类型不匹配")
		return
	}

	// 合并数据：如果请求中没有提供某些字段，使用已有的值
	if robot.Name == "" {
		robot.Name = existing.Name
	}
	if robot.Server == "" {
		robot.Server = existing.Server
	}
	if robot.Port == nil {
		robot.Port = existing.Port
	}
	if robot.Username == "" {
		robot.Username = existing.Username
	}
	if robot.Password == "" {
		robot.Password = existing.Password
	}
	if robot.Webhook == "" {
		robot.Webhook = existing.Webhook
	}
	if robot.Method == "" {
		robot.Method = existing.Method
	}

	// 验证必填字段
	if robot.Name == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "机器人名称不能为空")
		return
	}

	// 根据类型验证不同的字段
	if robotType == "email" {
		// 邮箱类型验证
		if robot.Server == "" || robot.Port == nil || *robot.Port == 0 || robot.Username == "" || robot.Password == "" {
			result.Failed(c, int(result.ApiCode.FAILED), "邮箱配置不完整，请检查 server/port/username/password 字段")
			return
		}
		// 如果没有填写接收邮箱，默认使用发件邮箱
		if robot.Webhook == "" {
			robot.Webhook = robot.Username
		}
	} else if robotType == "webhook" || robotType == "teams" {
		// Webhook/Teams类型验证
		if robot.Webhook == "" {
			result.Failed(c, int(result.ApiCode.FAILED), "Webhook地址不能为空")
			return
		}
		// 设置默认HTTP方法
		if robot.Method == "" {
			robot.Method = "POST"
		}
	} else {
		// 其他机器人类型验证（钉钉/飞书/企业微信）
		if robot.Webhook == "" {
			result.Failed(c, int(result.ApiCode.FAILED), "Webhook地址不能为空")
			return
		}
	}

	// 保持类型不变
	robot.Type = robotType

	if err := s.robotDao.Update(robot); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新失败: %v", err))
		return
	}

	result.Success(c, robot)
}

// Delete 删除通知机器人
func (s *NotifyRobotServiceImpl) Delete(c *gin.Context, id uint, robotType string) {
	// 检查机器人是否存在
	existing, err := s.robotDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("机器人不存在: %v", err))
		return
	}

	// 验证机器人类型（如果指定了类型）
	if robotType != "" && existing.Type != robotType {
		result.Failed(c, int(result.ApiCode.FAILED), "机器人类型不匹配")
		return
	}

	if err := s.robotDao.Delete(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("删除失败: %v", err))
		return
	}

	result.Success(c, "删除成功")
}

// GetList 获取通知机器人列表
func (s *NotifyRobotServiceImpl) GetList(c *gin.Context, robotType string, status *int, page, pageSize int) {
	// 设置默认分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	robots, total, err := s.robotDao.GetList(robotType, status, page, pageSize)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		return
	}

	result.Success(c, gin.H{
		"list":     robots,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// TestSend 测试发送消息（支持临时模板）
func (s *NotifyRobotServiceImpl) TestSend(c *gin.Context, id uint, tempTemplate string) {
	// 获取机器人配置
	robot, err := s.robotDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("机器人不存在: %v", err))
		return
	}

	// 如果传入了临时模板，使用临时模板（不修改数据库）
	if tempTemplate != "" {
		// 创建一个副本，避免修改原对象
		testRobot := *robot
		testRobot.Template = tempTemplate
		robot = &testRobot
	}

	// 根据类型发送测试消息
	var sendErr error
	switch robot.Type {
	case "dingtalk":
		sendErr = s.sendDingTalkMessage(robot)
	case "feishu":
		sendErr = s.sendFeishuMessage(robot)
	case "wechat":
		sendErr = s.sendWechatMessage(robot)
	case "email":
		sendErr = s.sendEmailMessage(robot)
	case "webhook":
		sendErr = s.sendWebhookMessage(robot)
	case "teams":
		sendErr = s.sendTeamsMessage(robot)
	default:
		result.Failed(c, int(result.ApiCode.FAILED), "不支持的机器人类型")
		return
	}

	if sendErr != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("发送失败: %v", sendErr))
		return
	}

	result.Success(c, "测试消息发送成功")
}

// sendDingTalkMessage 发送钉钉消息
func (s *NotifyRobotServiceImpl) sendDingTalkMessage(robot *model.NotifyRobot) error {
	timestamp := time.Now().UnixMilli()
	sign := ""

	// 如果配置了密钥，生成签名
	if robot.Secret != "" {
		sign = s.generateDingTalkSign(timestamp, robot.Secret)
	}

	// 构建请求URL
	url := robot.Webhook
	if sign != "" {
		url = fmt.Sprintf("%s&timestamp=%d&sign=%s", robot.Webhook, timestamp, sign)
	}

	var messageText string
	if robot.Template != "" {
		// 使用模板渲染测试消息
		messageText = s.renderTestTemplate(robot.Template)
	} else {
		// 使用默认测试消息
		messageText = fmt.Sprintf("【测试消息】\n机器人名称：%s\n发送时间：%s\n这是一条来自DoDevOps平台的测试消息。",
			robot.Name, time.Now().Format("2006-01-02 15:04:05"))
	}

	// 构建钉钉消息体（支持 Markdown）
	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "测试消息",
			"text":  messageText,
		},
	}

	return s.sendHTTPRequest(url, message)
}

// sendFeishuMessage 发送飞书消息
func (s *NotifyRobotServiceImpl) sendFeishuMessage(robot *model.NotifyRobot) error {
	var messageText string

	if robot.Template != "" {
		// 使用模板渲染测试消息
		messageText = s.renderTestTemplate(robot.Template)
	} else {
		// 使用默认测试消息
		messageText = fmt.Sprintf("【测试消息】\n机器人名称：%s\n发送时间：%s\n这是一条来自DoDevOps平台的测试消息。",
			robot.Name, time.Now().Format("2006-01-02 15:04:05"))
	}

	// 构建飞书消息体
	message := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": messageText,
		},
	}

	return s.sendHTTPRequest(robot.Webhook, message)
}

// sendWechatMessage 发送企业微信消息
func (s *NotifyRobotServiceImpl) sendWechatMessage(robot *model.NotifyRobot) error {
	var messageText string
	if robot.Template != "" {
		// 使用模板渲染测试消息
		messageText = s.renderTestTemplate(robot.Template)
	} else {
		// 使用默认测试消息
		messageText = fmt.Sprintf("【测试消息】\n机器人名称：%s\n发送时间：%s\n这是一条来自DoDevOps平台的测试消息。",
			robot.Name, time.Now().Format("2006-01-02 15:04:05"))
	}

	// 构建企业微信消息体（支持 Markdown）
	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": messageText,
		},
	}

	return s.sendHTTPRequest(robot.Webhook, message)
}

// generateDingTalkSign 生成钉钉签名
func (s *NotifyRobotServiceImpl) generateDingTalkSign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// renderTestTemplate 渲染测试模板
func (s *NotifyRobotServiceImpl) renderTestTemplate(templateStr string) string {
	// 构建测试数据
	testData := map[string]interface{}{
		"IsRecovered":               false,
		"alert_title":               "测试告警-系统测试",
		"target_name":               "test-server (192.168.1.100)",
		"rule_name":                 "测试规则",
		"policy_name":               "测试策略",
		"fingerprint":               "test-fingerprint-12345",
		"severity":                  "测试",
		"severity_level":            2,
		"metric_name":               "test_metric",
		"current_value":             "85.5",
		"threshold":                 "80.0",
		"first_trigger_time_format": time.Now().Format("2006-01-02 15:04:05"),
		"duration":                  120,
		"duration_format":           "2分钟",
		"group_name":                "测试业务组",
		"message":                   "这是一条测试消息",
		"annotations":               "这是一条测试消息",
		"labels": map[string]string{
			"hostname":  "test-server",
			"instance":  "192.168.1.100",
			"pod":       "test-pod",
			"namespace": "test-namespace",
		},
		"duty_user":     "测试值班人员",
		"faultCenterId": 999,
	}

	// 解析并执行 Go template
	tmpl, err := template.New("test").Parse(templateStr)
	if err != nil {
		return fmt.Sprintf("模板解析失败: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, testData); err != nil {
		return fmt.Sprintf("模板执行失败: %v", err)
	}

	message := buf.String()

	// 替换 ${xxx} 格式的变量
	for key, value := range testData {
		if key == "labels" {
			// 处理 labels 子字段
			if labels, ok := value.(map[string]string); ok {
				for labelKey, labelValue := range labels {
					placeholder := fmt.Sprintf("${labels.%s}", labelKey)
					message = strings.ReplaceAll(message, placeholder, labelValue)
				}
			}
			continue
		}
		placeholder := fmt.Sprintf("${%s}", key)
		message = strings.ReplaceAll(message, placeholder, fmt.Sprintf("%v", value))
	}

	return message
}

// sendHTTPRequest 发送HTTP请求
func (s *NotifyRobotServiceImpl) sendHTTPRequest(url string, message map[string]interface{}) error {
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

// sendEmailMessage 发送邮件消息
func (s *NotifyRobotServiceImpl) sendEmailMessage(robot *model.NotifyRobot) error {
	// 验证必填字段
	if robot.Server == "" || robot.Port == nil || *robot.Port == 0 || robot.Username == "" || robot.Password == "" {
		return fmt.Errorf("邮箱配置不完整，请检查 server/port/username/password 字段")
	}

	// 如果没有接收邮箱，使用发件邮箱
	toEmail := robot.Webhook
	if toEmail == "" {
		toEmail = robot.Username
	}

	// 构建邮件主题和内容
	var subject, body string
	if robot.Template != "" {
		// 使用自定义模板
		templateStr := robot.Template
		templateStr = strings.ReplaceAll(templateStr, "{{name}}", robot.Name)
		templateStr = strings.ReplaceAll(templateStr, "{{time}}", time.Now().Format("2006-01-02 15:04:05"))
		templateStr = strings.ReplaceAll(templateStr, "{{content}}", "这是一条来自DoDevOps平台的测试消息。")

		// 尝试解析为JSON格式（包含subject和body字段）
		var emailTemplate map[string]string
		if err := json.Unmarshal([]byte(templateStr), &emailTemplate); err == nil {
			subject = emailTemplate["subject"]
			body = emailTemplate["body"]
		} else {
			// 如果不是JSON，整个模板作为邮件正文
			subject = "【DoDevOps告警通知】"
			body = templateStr
		}
	} else {
		// 使用默认模板
		subject = "【DoDevOps测试消息】"
		body = fmt.Sprintf("机器人名称：%s\n发送时间：%s\n这是一条来自DoDevOps平台的测试消息。",
			robot.Name, time.Now().Format("2006-01-02 15:04:05"))
	}

	// 获取发件人昵称（如果没有设置则使用用户名）
	fromName := robot.Nickname
	if fromName == "" {
		fromName = robot.Username
	}

	// 构建邮件
	from := robot.Username
	to := []string{toEmail}

	// 邮件头部
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, from)
	headers["To"] = toEmail
	headers["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=UTF-8"
	headers["Content-Transfer-Encoding"] = "base64"

	// 组装邮件内容
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// SMTP 服务器地址
	addr := fmt.Sprintf("%s:%d", robot.Server, *robot.Port)

	// 认证信息
	auth := smtp.PlainAuth("", robot.Username, robot.Password, robot.Server)

	// 根据端口自动判断加密方式
	// 465: SMTPS (直接SSL/TLS连接)
	// 587, 25: STARTTLS (先普通连接再升级) 或普通连接
	// 其他端口: 尝试STARTTLS,失败则使用普通连接
	port := *robot.Port

	if port == 465 {
		// 465端口：SMTPS，使用直接的SSL/TLS连接
		return s.sendEmailWithSSL(addr, auth, from, to, []byte(message), robot.Server)
	}

	// 587端口或其他端口：尝试STARTTLS
	if port == 587 || port > 25 {
		// 先尝试STARTTLS方式
		err := s.sendEmailWithSTARTTLS(addr, auth, from, to, []byte(message), robot.Server)
		if err == nil {
			return nil
		}
		// STARTTLS失败，如果是587端口则直接返回错误（587必须使用TLS）
		if port == 587 {
			return err
		}
		// 其他端口STARTTLS失败，继续尝试普通连接
	}

	// 25端口或其他端口普通方式发送
	return smtp.SendMail(addr, auth, from, to, []byte(message))
}

// sendEmailWithSSL 使用SSL/TLS发送邮件(465端口,直接SSL连接)
func (s *NotifyRobotServiceImpl) sendEmailWithSSL(addr string, auth smtp.Auth, from string, to []string, msg []byte, serverName string) error {
	// 建立 TLS 连接
	// 注意: InsecureSkipVerify设为true是为了兼容某些企业邮箱的自签名证书
	// 生产环境建议设为false并正确配置证书
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: true, // 跳过证书验证,解决某些企业邮箱证书问题
	})
	if err != nil {
		return fmt.Errorf("SSL连接失败(服务器: %s): %v", serverName, err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, serverName)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP认证失败(用户名: %s, 服务器: %s): %v", from, serverName, err)
		}
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("设置收件人失败: %v", err)
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("准备发送邮件失败: %v", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭邮件写入失败: %v", err)
	}

	return client.Quit()
}

// sendEmailWithSTARTTLS 使用STARTTLS发送邮件(587端口,先建立普通连接再升级)
func (s *NotifyRobotServiceImpl) sendEmailWithSTARTTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte, serverName string) error {
	// 先建立普通TCP连接
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %v", err)
	}
	defer client.Close()

	// 发送EHLO命令
	if err = client.Hello(serverName); err != nil {
		return fmt.Errorf("EHLO命令失败: %v", err)
	}

	// 升级到TLS
	if err = client.StartTLS(&tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: true, // 跳过证书验证,解决某些企业邮箱证书问题
	}); err != nil {
		return fmt.Errorf("STARTTLS升级失败: %v", err)
	}

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP认证失败(用户名: %s, 服务器: %s): %v", from, serverName, err)
		}
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("设置收件人失败: %v", err)
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("准备发送邮件失败: %v", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭邮件写入失败: %v", err)
	}

	return client.Quit()
}

// sendWebhookMessage 发送通用Webhook消息
func (s *NotifyRobotServiceImpl) sendWebhookMessage(robot *model.NotifyRobot) error {
	// 构建默认消息内容
	defaultContent := map[string]interface{}{
		"title":   "DoDevOps测试消息",
		"content": fmt.Sprintf("机器人名称：%s\n发送时间：%s\n这是一条来自DoDevOps平台的测试消息。", robot.Name, time.Now().Format("2006-01-02 15:04:05")),
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	}

	// 如果配置了自定义模板，使用模板渲染
	var message map[string]interface{}
	if robot.Template != "" {
		// 尝试解析模板为JSON
		if err := json.Unmarshal([]byte(robot.Template), &message); err != nil {
			// 如果不是JSON，作为纯文本处理
			message = map[string]interface{}{
				"text": robot.Template,
			}
		}
		// 模板变量替换（简单实现）
		messageJSON, _ := json.Marshal(message)
		messageStr := string(messageJSON)
		messageStr = strings.ReplaceAll(messageStr, "{{name}}", robot.Name)
		messageStr = strings.ReplaceAll(messageStr, "{{time}}", time.Now().Format("2006-01-02 15:04:05"))
		messageStr = strings.ReplaceAll(messageStr, "{{content}}", "这是一条来自DoDevOps平台的测试消息。")
		json.Unmarshal([]byte(messageStr), &message)
	} else {
		message = defaultContent
	}

	return s.sendCustomHTTPRequest(robot.Webhook, robot.Method, robot.Headers, message)
}

// sendTeamsMessage 发送Microsoft Teams消息
func (s *NotifyRobotServiceImpl) sendTeamsMessage(robot *model.NotifyRobot) error {
	var message map[string]interface{}

	if robot.Template != "" {
		// 使用自定义模板
		if err := json.Unmarshal([]byte(robot.Template), &message); err != nil {
			return fmt.Errorf("解析Teams模板失败: %v", err)
		}
		// 模板变量替换
		message = s.replaceTemplateVars(message, robot, "测试消息")
	} else {
		// 使用默认Teams Adaptive Card格式
		message = map[string]interface{}{
			"@type":      "MessageCard",
			"@context":   "https://schema.org/extensions",
			"summary":    "DoDevOps测试消息",
			"themeColor": "0078D7",
			"title":      "DoDevOps测试消息",
			"sections": []map[string]interface{}{
				{
					"activityTitle": fmt.Sprintf("机器人：%s", robot.Name),
					"facts": []map[string]string{
						{"name": "发送时间", "value": time.Now().Format("2006-01-02 15:04:05")},
						{"name": "消息内容", "value": "这是一条来自DoDevOps平台的测试消息。"},
					},
				},
			},
		}
	}

	return s.sendCustomHTTPRequest(robot.Webhook, "POST", robot.Headers, message)
}

// sendCustomHTTPRequest 发送自定义HTTP请求（支持自定义方法和Headers）
func (s *NotifyRobotServiceImpl) sendCustomHTTPRequest(url, method, headersJSON string, message map[string]interface{}) error {
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

// replaceTemplateVars 替换模板中的变量
// 支持的变量: {{name}}, {{time}}, {{content}}, {{title}}
func (s *NotifyRobotServiceImpl) replaceTemplateVars(message map[string]interface{}, robot *model.NotifyRobot, content string) map[string]interface{} {
	// 将消息转换为JSON字符串
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return message
	}

	messageStr := string(messageJSON)

	// 替换变量
	messageStr = strings.ReplaceAll(messageStr, "{{name}}", robot.Name)
	messageStr = strings.ReplaceAll(messageStr, "{{time}}", time.Now().Format("2006-01-02 15:04:05"))
	messageStr = strings.ReplaceAll(messageStr, "{{content}}", content)
	messageStr = strings.ReplaceAll(messageStr, "{{title}}", "DoDevOps告警通知")

	// 转换回map
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(messageStr), &result); err != nil {
		return message
	}

	return result
}
