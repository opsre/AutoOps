package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	sysModel "dodevops-api/api/system/model"
	sysService "dodevops-api/api/system/service"
	"github.com/robfig/cron/v3"

	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/smtp"
	"strings"
)

// DomainScheduler 域名监控调度器
type DomainScheduler struct {
	domainDao     dao.DomainDao
	scheduleDao   dao.DomainScheduleDao
	domainService *DomainServiceImpl
	robotDao      dao.NotifyRobotDao
	cron          *cron.Cron
	entryID       cron.EntryID
	running       bool
	mutex         sync.RWMutex
}

var (
	domainScheduler *DomainScheduler
	schedulerOnce   sync.Once
)

// GetDomainScheduler 获取域名调度器实例（单例）
func GetDomainScheduler() *DomainScheduler {
	schedulerOnce.Do(func() {
		domainScheduler = &DomainScheduler{
			domainDao:     dao.NewDomainDao(),
			scheduleDao:   dao.NewDomainScheduleDao(),
			domainService: &DomainServiceImpl{domainDao: dao.NewDomainDao()},
			robotDao:      dao.NewNotifyRobotDao(),
			cron:          cron.New(),
			running:       false,
		}
	})
	return domainScheduler
}

// InitScheduler 初始化并启动定时任务（在应用启动时调用）
func InitScheduler() error {
	scheduler := GetDomainScheduler()

	// 获取配置
	schedule, err := scheduler.scheduleDao.GetOrCreate()
	if err != nil {
		return fmt.Errorf("获取调度配置失败: %v", err)
	}

	// 如果配置为启用状态，自动启动
	if schedule.Enabled && schedule.CronExpr != "" {
		if err := scheduler.Start(schedule.CronExpr); err != nil {
			return fmt.Errorf("启动定时任务失败: %v", err)
		}
		log.Printf("域名监控定时任务已自动启动，Cron表达式：%s", schedule.CronExpr)
	} else {
		log.Println("域名监控定时任务未启用")
	}

	return nil
}

// Start 启动定时任务
func (ds *DomainScheduler) Start(cronExpr string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	if ds.running {
		return fmt.Errorf("定时任务已在运行")
	}

	if cronExpr == "" {
		return fmt.Errorf("Cron表达式不能为空")
	}

	// 添加定时任务
	entryID, err := ds.cron.AddFunc(cronExpr, func() {
		ds.executeCheck()
	})
	if err != nil {
		return fmt.Errorf("无效的Cron表达式: %v", err)
	}

	ds.entryID = entryID
	ds.cron.Start()
	ds.running = true

	// 更新状态
	ds.scheduleDao.UpdateStatus("running")

	log.Printf("域名监控定时任务已启动，Cron表达式：%s", cronExpr)
	return nil
}

// Stop 停止定时任务
func (ds *DomainScheduler) Stop() error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	if !ds.running {
		return fmt.Errorf("定时任务未运行")
	}

	// 移除定时任务
	ds.cron.Remove(ds.entryID)
	ds.cron.Stop()
	ds.running = false

	// 更新状态
	ds.scheduleDao.UpdateStatus("stopped")

	log.Println("域名监控定时任务已停止")
	return nil
}

// Pause 暂停定时任务
func (ds *DomainScheduler) Pause() error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	if !ds.running {
		return fmt.Errorf("定时任务未运行")
	}

	// 移除定时任务但不停止cron
	ds.cron.Remove(ds.entryID)
	ds.running = false

	// 更新状态
	ds.scheduleDao.UpdateStatus("paused")

	log.Println("域名监控定时任务已暂停")
	return nil
}

// Resume 恢复定时任务
func (ds *DomainScheduler) Resume(cronExpr string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	if ds.running {
		return fmt.Errorf("定时任务已在运行")
	}

	if cronExpr == "" {
		return fmt.Errorf("Cron表达式不能为空")
	}

	// 重新添加定时任务
	entryID, err := ds.cron.AddFunc(cronExpr, func() {
		ds.executeCheck()
	})
	if err != nil {
		return fmt.Errorf("无效的Cron表达式: %v", err)
	}

	ds.entryID = entryID
	ds.running = true

	// 更新状态
	ds.scheduleDao.UpdateStatus("running")

	log.Printf("域名监控定时任务已恢复，Cron表达式：%s", cronExpr)
	return nil
}

// executeCheck 执行检查
func (ds *DomainScheduler) executeCheck() {
	startTime := time.Now()
	log.Println("开始执行域名监控定时检查...")

	// 获取所有启用的域名
	domains, err := ds.domainDao.GetEnabledDomains()
	if err != nil {
		log.Printf("获取域名列表失败: %v", err)
		// 记录失败日志
		ds.recordActivityLog(0, 0, time.Since(startTime).Seconds(), 2, fmt.Sprintf("获取域名列表失败: %v", err))
		return
	}

	totalCount := len(domains)
	log.Printf("共 %d 个域名需要检查", totalCount)

	// 使用并发控制检查域名
	maxConcurrent := 10
	semaphore := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	var mu sync.Mutex
	expiringCount := 0                     // 即将过期的域名数量（小于30天）
	abnormalCount := 0                     // 异常域名数量
	abnormalDomains := []string{}          // 异常域名列表
	expiringDomains := []string{}          // 即将过期域名列表

	for _, domain := range domains {
		wg.Add(1)

		go func(d model.Domain) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 执行检查
			checkResult := ds.domainService.performCheck(d.Domain)
			checkResult.DomainID = d.ID

			// 统计即将过期的域名（SSL证书剩余天数小于30天且大于0）
			if checkResult.SSLDaysLeft > 0 && checkResult.SSLDaysLeft < 30 {
				mu.Lock()
				expiringCount++
				expiringDomains = append(expiringDomains, d.Domain)
				mu.Unlock()
			}

			// 统计异常域名（IsAlive == 0）
			if checkResult.IsAlive == 0 {
				mu.Lock()
				abnormalCount++
				abnormalDomains = append(abnormalDomains, d.Domain)
				mu.Unlock()
			}

			// 更新检查结果
			if err := ds.domainDao.UpdateCheckResult(d.ID, checkResult); err != nil {
				log.Printf("更新域名 %s 检查结果失败: %v", d.Domain, err)
			} else {
				log.Printf("域名 %s 检查完成，状态：%d", d.Domain, checkResult.IsAlive)
			}
		}(domain)
	}

	wg.Wait()

	// 计算执行耗时
	duration := time.Since(startTime).Seconds()

	// 记录活动日志
	ds.recordActivityLog(totalCount, expiringCount, duration, 1, "")

	// 发送通知
	ds.sendNotification(totalCount, abnormalCount, expiringCount, duration, abnormalDomains, expiringDomains)

	log.Println("定时检查完成")
}

// recordActivityLog 记录域名巡检活动日志
func (ds *DomainScheduler) recordActivityLog(totalCount, expiringCount int, duration float64, status int, errorMsg string) {
	activityLogService := sysService.NewSysActivityLogService()

	var title, content, summary string
	if status == 1 {
		// 成功
		title = "域名定时巡检完成"
		summary = fmt.Sprintf("本次巡检 %d 个域名，%d 个域名即将过期", totalCount, expiringCount)
		content = fmt.Sprintf("巡检结果：\n- 巡检域名总数：%d\n- 即将过期域名数：%d（SSL证书剩余天数<30天）\n- 执行耗时：%.2f秒",
			totalCount, expiringCount, duration)
	} else {
		// 失败
		title = "域名定时巡检失败"
		summary = errorMsg
		content = fmt.Sprintf("巡检失败：%s\n执行耗时：%.2f秒", errorMsg, duration)
	}

	activityLog := &sysModel.SysActivityLog{
		ActivityType: sysModel.ActivityTypeCheckDomain, // 域名检查类型
		Title:        title,
		Content:      content,
		Status:       status,
		Summary:      summary,
		Duration:     int64(duration),
	}

	if err := activityLogService.Create(activityLog); err != nil {
		log.Printf("记录域名巡检活动日志失败: %v", err)
	}
}

// IsRunning 检查是否正在运行
func (ds *DomainScheduler) IsRunning() bool {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return ds.running
}

// GetStatus 获取调度器状态
func (ds *DomainScheduler) GetStatus() map[string]interface{} {
	ds.mutex.RLock()
	running := ds.running
	ds.mutex.RUnlock()

	schedule, _ := ds.scheduleDao.GetOrCreate()

	status := map[string]interface{}{
		"running": running,
		"config":  schedule,
	}

	return status
}

// sendNotification 发送巡检通知
func (ds *DomainScheduler) sendNotification(totalCount, abnormalCount, expiringCount int, duration float64, abnormalDomains, expiringDomains []string) {
	// 1. 获取调度配置
	schedule, err := ds.scheduleDao.GetOrCreate()
	if err != nil {
		log.Printf("获取调度配置失败，跳过通知: %v", err)
		return
	}

	// 2. 检查是否开启通知
	if !schedule.NotifyEnabled || schedule.NotifyRobotId == nil {
		log.Println("通知未开启或未配置机器人，跳过通知")
		return
	}

	log.Printf("开始发送通知: 异常数=%d, 过期数=%d, 耗时=%.2f秒", abnormalCount, expiringCount, duration)

	// 3. 获取机器人配置
	robot, err := ds.robotDao.GetByID(*schedule.NotifyRobotId)
	if err != nil {
		log.Printf("获取机器人配置失败: %v", err)
		return
	}

	log.Printf("获取到机器人: ID=%d, Name=%s, Type=%s, Status=%d", robot.ID, robot.Name, robot.Type, robot.Status)

	// 4. 检查机器人状态
	if robot.Status != 1 {
		log.Printf("机器人[%s]未启用，跳过通知", robot.Name)
		return
	}

	// 5. 根据机器人类型发送通知
	var sendErr error
	switch robot.Type {
	case "dingtalk":
		log.Println("发送钉钉通知...")
		sendErr = ds.sendDingTalkNotification(robot, totalCount, abnormalCount, expiringCount, duration, abnormalDomains, expiringDomains)
	case "feishu":
		log.Println("发送飞书通知...")
		sendErr = ds.sendFeishuNotification(robot, totalCount, abnormalCount, expiringCount, duration, abnormalDomains, expiringDomains)
	case "wechat":
		log.Println("发送企业微信通知...")
		sendErr = ds.sendWechatNotification(robot, totalCount, abnormalCount, expiringCount, duration, abnormalDomains, expiringDomains)
	case "email":
		log.Println("发送邮件通知...")
		sendErr = ds.sendEmailNotification(robot, totalCount, abnormalCount, expiringCount, duration, abnormalDomains, expiringDomains)
	case "webhook":
		log.Println("发送Webhook通知...")
		sendErr = ds.sendWebhookNotification(robot, totalCount, abnormalCount, expiringCount, duration, abnormalDomains, expiringDomains)
	case "teams":
		log.Println("发送Teams通知...")
		sendErr = ds.sendTeamsNotification(robot, totalCount, abnormalCount, expiringCount, duration, abnormalDomains, expiringDomains)
	default:
		log.Printf("不支持的机器人类型: %s", robot.Type)
		return
	}

	if sendErr != nil {
		log.Printf("❌ 发送通知失败: %v", sendErr)
	} else {
		log.Printf("✅ 巡检通知已成功发送到机器人[%s]", robot.Name)
	}
}

// determineAlertLevel 判断告警级别
func (ds *DomainScheduler) determineAlertLevel(totalCount, abnormalCount, expiringCount int) (level string, emoji string) {
	abnormalRate := 0.0
	if totalCount > 0 {
		abnormalRate = float64(abnormalCount) / float64(totalCount)
	}

	if abnormalRate >= 0.3 || abnormalCount >= 10 || expiringCount >= 5 {
		return "严重", "🔴"
	} else if abnormalRate >= 0.1 || abnormalCount >= 3 || expiringCount >= 2 {
		return "警告", "🟡"
	}
	return "正常", "🟢"
}

// sendDingTalkNotification 发送钉钉通知
func (ds *DomainScheduler) sendDingTalkNotification(robot *model.NotifyRobot, totalCount, abnormalCount, expiringCount int, duration float64, abnormalDomains, expiringDomains []string) error {
	level, emoji := ds.determineAlertLevel(totalCount, abnormalCount, expiringCount)

	timestamp := time.Now().UnixMilli()
	sign := ""
	if robot.Secret != "" {
		sign = ds.generateDingTalkSign(timestamp, robot.Secret)
	}

	url := robot.Webhook
	if sign != "" {
		url = fmt.Sprintf("%s&timestamp=%d&sign=%s", robot.Webhook, timestamp, sign)
	}

	// 构建异常域名列表字符串
	abnormalList := ""
	if abnormalCount > 0 {
		for _, domain := range abnormalDomains {
			abnormalList += "- " + domain + "\n"
		}
	}

	// 构建过期域名列表字符串
	expiringList := ""
	if expiringCount > 0 {
		for _, domain := range expiringDomains {
			expiringList += "- " + domain + "\n"
		}
	}

	content := fmt.Sprintf(
		"# %s 域名巡检完成通知 %s\n\n"+
			"**监控项目**: 域名巡检\n\n"+
			"**巡检时间**: %s\n\n"+
			"**执行耗时**: %.2f秒\n\n"+
			"**告警级别**: %s %s\n\n"+
			"**巡检统计**:\n"+
			"- 域名总数: %d\n"+
			"- 异常域名: %d\n"+
			"- SSL即将过期: %d（<30天）\n\n"+
			"%s%s"+
			"请及时处理异常域名和即将过期的SSL证书。",
		emoji, emoji,
		time.Now().Format("2006-01-02 15:04:05"),
		duration,
		level, emoji,
		totalCount, abnormalCount, expiringCount,
		func() string {
			if abnormalCount > 0 {
				return "**异常域名列表**:\n" + abnormalList + "\n"
			}
			return ""
		}(),
		func() string {
			if expiringCount > 0 {
				return "**SSL即将过期域名**:\n" + expiringList + "\n"
			}
			return ""
		}(),
	)

	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": fmt.Sprintf("%s 域名巡检完成通知 %s", emoji, emoji),
			"text":  content,
		},
	}

	return ds.sendHTTPRequest(url, message)
}

// sendFeishuNotification 发送飞书通知 - 使用简单可靠的Markdown文本格式
func (ds *DomainScheduler) sendFeishuNotification(robot *model.NotifyRobot, totalCount, abnormalCount, expiringCount int, duration float64, abnormalDomains, expiringDomains []string) error {
	level, emoji := ds.determineAlertLevel(totalCount, abnormalCount, expiringCount)

	// 构建异常域名列表字符串
	abnormalList := ""
	if abnormalCount > 0 {
		for _, domain := range abnormalDomains {
			abnormalList += "  • " + domain + "\n"
		}
	}

	// 构建过期域名列表字符串
	expiringList := ""
	if expiringCount > 0 {
		for _, domain := range expiringDomains {
			expiringList += "  • " + domain + "\n"
		}
	}

	// 构建Markdown格式的文本
	content := fmt.Sprintf(
		"%s 域名巡检完成通知 %s\n\n"+
			"**监控项目**: 域名巡检\n\n"+
			"**巡检时间**: %s\n\n"+
			"**执行耗时**: %.2f秒\n\n"+
			"**告警级别**: %s %s\n\n"+
			"**巡检统计**:\n"+
			"  • 域名总数: %d\n"+
			"  • 异常域名: %d\n"+
			"  • SSL即将过期: %d（<30天）\n\n"+
			"%s%s"+
			"请及时处理异常域名和即将过期的SSL证书。",
		emoji, emoji,
		time.Now().Format("2006-01-02 15:04:05"),
		duration,
		level, emoji,
		totalCount, abnormalCount, expiringCount,
		func() string {
			if abnormalCount > 0 {
				return "**异常域名列表**:\n" + abnormalList + "\n"
			}
			return ""
		}(),
		func() string {
			if expiringCount > 0 {
				return "**SSL即将过期域名**:\n" + expiringList + "\n"
			}
			return ""
		}(),
	)

	// 使用text消息类型发送Markdown格式内容
	message := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": content,
		},
	}

	return ds.sendHTTPRequest(robot.Webhook, message)
}

// sendWechatNotification 发送企业微信通知
func (ds *DomainScheduler) sendWechatNotification(robot *model.NotifyRobot, totalCount, abnormalCount, expiringCount int, duration float64, abnormalDomains, expiringDomains []string) error {
	level, emoji := ds.determineAlertLevel(totalCount, abnormalCount, expiringCount)

	// 构建异常域名列表字符串
	abnormalList := ""
	if abnormalCount > 0 {
		for _, domain := range abnormalDomains {
			abnormalList += "- " + domain + "\n"
		}
	}

	// 构建过期域名列表字符串
	expiringList := ""
	if expiringCount > 0 {
		for _, domain := range expiringDomains {
			expiringList += "- " + domain + "\n"
		}
	}

	content := fmt.Sprintf(
		"# %s 域名巡检完成通知 %s\n\n"+
			"**监控项目**: 域名巡检\n\n"+
			"**巡检时间**: %s\n\n"+
			"**执行耗时**: %.2f秒\n\n"+
			"**告警级别**: %s %s\n\n"+
			"**巡检统计**:\n"+
			"- 域名总数: %d\n"+
			"- 异常域名: %d\n"+
			"- SSL即将过期: %d（<30天）\n\n"+
			"%s%s"+
			"请及时处理异常域名和即将过期的SSL证书。",
		emoji, emoji,
		time.Now().Format("2006-01-02 15:04:05"),
		duration,
		level, emoji,
		totalCount, abnormalCount, expiringCount,
		func() string {
			if abnormalCount > 0 {
				return "**异常域名列表**:\n" + abnormalList + "\n"
			}
			return ""
		}(),
		func() string {
			if expiringCount > 0 {
				return "**SSL即将过期域名**:\n" + expiringList + "\n"
			}
			return ""
		}(),
	)

	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": content,
		},
	}

	return ds.sendHTTPRequest(robot.Webhook, message)
}

// sendEmailNotification 发送邮件通知
func (ds *DomainScheduler) sendEmailNotification(robot *model.NotifyRobot, totalCount, abnormalCount, expiringCount int, duration float64, abnormalDomains, expiringDomains []string) error {
	level, emoji := ds.determineAlertLevel(totalCount, abnormalCount, expiringCount)

	// 构建异常域名列表字符串
	abnormalList := ""
	if abnormalCount > 0 {
		for _, domain := range abnormalDomains {
			abnormalList += domain + "\n"
		}
	}

	// 构建过期域名列表字符串
	expiringList := ""
	if expiringCount > 0 {
		for _, domain := range expiringDomains {
			expiringList += domain + "\n"
		}
	}

	subject := fmt.Sprintf("%s 域名巡检完成通知 - %s", emoji, level)
	body := fmt.Sprintf(
		"域名巡检完成通知\n\n"+
			"监控项目：域名巡检\n"+
			"巡检时间：%s\n"+
			"执行耗时：%.2f秒\n"+
			"告警级别：%s %s\n\n"+
			"巡检统计：\n"+
			"- 域名总数：%d\n"+
			"- 异常域名：%d\n"+
			"- SSL即将过期：%d（<30天）\n\n"+
			"%s%s%s"+
			"请及时处理异常域名和即将过期的SSL证书。\n\n"+
			"---\n"+
			"此邮件由DoDevOps平台自动发送，请勿回复。",
		time.Now().Format("2006-01-02 15:04:05"),
		duration,
		level, emoji,
		totalCount, abnormalCount, expiringCount,
		func() string {
			if abnormalCount > 0 {
				return "异常域名列表:\n" + abnormalList + "\n"
			}
			return ""
		}(),
		func() string {
			if expiringCount > 0 {
				return "SSL即将过期域名:\n" + expiringList + "\n"
			}
			return ""
		}(),
		"",
	)

	return ds.sendEmail(robot, subject, body)
}

// sendWebhookNotification 发送Webhook通知
func (ds *DomainScheduler) sendWebhookNotification(robot *model.NotifyRobot, totalCount, abnormalCount, expiringCount int, duration float64, abnormalDomains, expiringDomains []string) error {
	level, emoji := ds.determineAlertLevel(totalCount, abnormalCount, expiringCount)

	message := map[string]interface{}{
		"title":             "域名巡检完成通知",
		"level":             level,
		"emoji":             emoji,
		"time":              time.Now().Format("2006-01-02 15:04:05"),
		"duration":          duration,
		"totalCount":        totalCount,
		"abnormalCount":     abnormalCount,
		"expiringCount":     expiringCount,
		"abnormalDomains":   abnormalDomains,
		"expiringDomains":   expiringDomains,
	}

	return ds.sendCustomHTTPRequest(robot.Webhook, robot.Method, robot.Headers, message)
}

// sendTeamsNotification 发送Teams通知
func (ds *DomainScheduler) sendTeamsNotification(robot *model.NotifyRobot, totalCount, abnormalCount, expiringCount int, duration float64, abnormalDomains, expiringDomains []string) error {
	level, emoji := ds.determineAlertLevel(totalCount, abnormalCount, expiringCount)

	// 构建异常域名字符串
	abnormalDomainsStr := ""
	if len(abnormalDomains) > 0 {
		abnormalDomainsStr = fmt.Sprintf("\n**异常域名列表**: %s", strings.Join(abnormalDomains, ", "))
	}

	// 构建过期域名字符串
	expiringDomainsStr := ""
	if len(expiringDomains) > 0 {
		expiringDomainsStr = fmt.Sprintf("\n**SSL即将过期域名**: %s", strings.Join(expiringDomains, ", "))
	}

	message := map[string]interface{}{
		"@type":      "MessageCard",
		"@context":   "https://schema.org/extensions",
		"summary":    "域名巡检完成通知",
		"themeColor": ds.getThemeColor(level),
		"title":      fmt.Sprintf("%s 域名巡检完成通知", emoji),
		"sections": []map[string]interface{}{
			{
				"activityTitle": fmt.Sprintf("告警级别：%s %s", level, emoji),
				"facts": []map[string]string{
					{"name": "监控项目", "value": "域名巡检"},
					{"name": "巡检时间", "value": time.Now().Format("2006-01-02 15:04:05")},
					{"name": "执行耗时", "value": fmt.Sprintf("%.2f秒", duration)},
					{"name": "域名总数", "value": fmt.Sprintf("%d", totalCount)},
					{"name": "异常域名", "value": fmt.Sprintf("%d", abnormalCount)},
					{"name": "SSL即将过期", "value": fmt.Sprintf("%d（<30天）", expiringCount)},
				},
				"text": fmt.Sprintf("请及时处理异常域名和即将过期的SSL证书。%s%s", abnormalDomainsStr, expiringDomainsStr),
			},
		},
	}

	return ds.sendCustomHTTPRequest(robot.Webhook, "POST", robot.Headers, message)
}

// getThemeColor 根据告警级别获取主题颜色
func (ds *DomainScheduler) getThemeColor(level string) string {
	switch level {
	case "严重":
		return "FF0000" // 红色
	case "警告":
		return "FFA500" // 橙色
	default:
		return "00FF00" // 绿色
	}
}

// generateDingTalkSign 生成钉钉签名
func (ds *DomainScheduler) generateDingTalkSign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// sendHTTPRequest 发送HTTP请求
func (ds *DomainScheduler) sendHTTPRequest(url string, message map[string]interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 打印请求体用于调试
	log.Printf("发送HTTP请求到: %s", url)
	log.Printf("请求体: %s", string(jsonData))

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

	// 打印响应状态码和响应体用于调试
	log.Printf("HTTP响应状态码: %d", resp.StatusCode)
	log.Printf("HTTP响应体: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP状态码异常: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 检查Feishu/飞书API返回的业务级错误码
	var apiResp map[string]interface{}
	if err := json.Unmarshal(body, &apiResp); err == nil {
		// 如果响应中包含code字段，检查是否为错误
		if code, ok := apiResp["code"]; ok {
			codeNum, isNum := code.(float64)
			if isNum && codeNum != 0 {
				msg := ""
				if msgVal, ok := apiResp["msg"]; ok {
					msg = fmt.Sprintf("%v", msgVal)
				}
				return fmt.Errorf("API返回错误: code=%v, msg=%s", codeNum, msg)
			}
		}
	}

	return nil
}

// sendCustomHTTPRequest 发送自定义HTTP请求（支持自定义方法和Headers）
func (ds *DomainScheduler) sendCustomHTTPRequest(url, method, headersJSON string, message map[string]interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	if method == "" {
		method = "POST"
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP状态码异常: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return nil
}

// sendEmail 发送邮件
func (ds *DomainScheduler) sendEmail(robot *model.NotifyRobot, subject, body string) error {
	// 验证必填字段
	if robot.Server == "" || robot.Port == nil || *robot.Port == 0 || robot.Username == "" || robot.Password == "" {
		return fmt.Errorf("邮箱配置不完整")
	}

	toEmail := robot.Webhook
	if toEmail == "" {
		toEmail = robot.Username
	}

	fromName := robot.Nickname
	if fromName == "" {
		fromName = robot.Username
	}

	from := robot.Username
	to := []string{toEmail}

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, from)
	headers["To"] = toEmail
	headers["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=UTF-8"
	headers["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	addr := fmt.Sprintf("%s:%d", robot.Server, *robot.Port)
	auth := smtp.PlainAuth("", robot.Username, robot.Password, robot.Server)

	port := *robot.Port
	if port == 465 {
		return ds.sendEmailWithSSL(addr, auth, from, to, []byte(message), robot.Server)
	}

	if port == 587 || port > 25 {
		err := ds.sendEmailWithSTARTTLS(addr, auth, from, to, []byte(message), robot.Server)
		if err == nil {
			return nil
		}
		if port == 587 {
			return err
		}
	}

	return smtp.SendMail(addr, auth, from, to, []byte(message))
}

// sendEmailWithSSL 使用SSL/TLS发送邮件(465端口,直接SSL连接)
func (ds *DomainScheduler) sendEmailWithSSL(addr string, auth smtp.Auth, from string, to []string, msg []byte, serverName string) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return fmt.Errorf("SSL连接失败: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, serverName)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP认证失败: %v", err)
		}
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("设置收件人失败: %v", err)
		}
	}

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
func (ds *DomainScheduler) sendEmailWithSTARTTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte, serverName string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %v", err)
	}
	defer client.Close()

	if err = client.Hello(serverName); err != nil {
		return fmt.Errorf("EHLO命令失败: %v", err)
	}

	if err = client.StartTLS(&tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: true,
	}); err != nil {
		return fmt.Errorf("STARTTLS升级失败: %v", err)
	}

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP认证失败: %v", err)
		}
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("设置收件人失败: %v", err)
		}
	}

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
