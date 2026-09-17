package monitor

import (
	"dodevops-api/api/monitor/controller"

	"github.com/gin-gonic/gin"
)

func InitMonitorRouter(r *gin.RouterGroup) {
	monitorController := controller.NewMonitorController()
	agentController := controller.NewAgentController()
	domainController := controller.NewDomainController()
	domainScheduleController := controller.NewDomainScheduleController()
	sslCertController := controller.NewSslCertController()
	aliyunConfigController := controller.NewAliyunConfigController()
	alertStatsController := controller.NewAlertStatsController()
	cloudBillingController := controller.NewCloudBillingController()
	zabbixStatsController := controller.NewZabbixStatsController()
	notifyRobotController := controller.NewNotifyRobotController()
	webhookController := controller.NewWebhookController()
	userStatsController := controller.NewUserStatsController()
	alertRuleController := controller.NewAlertRuleController()                     // 告警规则控制器
	alertTemplateController := controller.NewAlertRuleTemplateController()         // 告警规则模板控制器
	alertPolicyController := controller.NewAlertPolicyController()                 // 告警策略控制器
	apiMonitorController := controller.NewAPIMonitorController()                   // API监控控制器
	apiBusinessFunctionController := controller.NewAPIBusinessFunctionController() // 业务功能控制器
	systemHealthController := controller.NewSystemHealthController()               // 系统健康检查控制器

	monitorGroup := r.Group("/monitor")
	// 不要重复使用 AuthMiddleware，因为 jwtGroup 已经使用了

	// 主机监控
	monitorGroup.GET("/host/:id", monitorController.GetHostMetrics)                        // 获取主机监控数据
	monitorGroup.GET("/hosts", monitorController.BatchGetHostMetrics)                      // 批量获取主机监控数据(包含在线状态)
	monitorGroup.GET("/hosts/top-usage", monitorController.GetTopUsageHosts)               // 获取使用率TOP主机
	monitorGroup.GET("/hosts/:id/history", monitorController.GetHostMetricHistory)         // 获取主机指定指标的历史数据
	monitorGroup.GET("/hosts/:id/all-metrics", monitorController.GetHostAllMetricsHistory) // 获取主机所有指标的历史数据
	monitorGroup.GET("/hosts/:id/top-processes", monitorController.GetTopProcesses)        // 获取主机TOP进程使用率
	monitorGroup.GET("/hosts/:id/ports", monitorController.GetHostPorts)                   // 获取主机端口信息
	monitorGroup.GET("/hosts/:id/complete", monitorController.GetCompleteHostMetrics)      // 获取主机完整监控数据

	// Agent管理
	monitorGroup.POST("/agent/deploy", agentController.DeployAgent)           // 部署agent到指定主机(支持单个或多个)
	monitorGroup.DELETE("/agent/uninstall", agentController.UninstallAgent)   // 卸载指定主机的agent(支持单个或多个)
	monitorGroup.GET("/agent/status/:id", agentController.GetAgentStatus)     // 根据主机id获取agent状态
	monitorGroup.POST("/agent/restart/:id", agentController.RestartAgent)     // 重启agent
	monitorGroup.GET("/agent/list", agentController.GetAgentList)             // 获取agent列表
	monitorGroup.GET("/agent/statistics", agentController.GetAgentStatistics) // 获取统计信息
	monitorGroup.DELETE("/agent/delete/:id", agentController.DeleteAgent)     // 删除agent数据(用于离线服务器)

	// 域名监控
	monitorGroup.POST("/domain", domainController.CreateDomain)                          // 创建域名监控
	monitorGroup.POST("/domain/batch-upload", domainController.BatchUploadDomains)       // 批量上传域名文件
	monitorGroup.POST("/domain/import-aliyun-dns", domainController.ImportFromAliyunDNS) // 从阿里云解析导入域名
	monitorGroup.GET("/domain/:id", domainController.GetDomain)                          // 获取域名监控详情
	monitorGroup.PUT("/domain", domainController.UpdateDomain)                           // 更新域名监控
	monitorGroup.DELETE("/domain/:id", domainController.DeleteDomain)                    // 删除域名监控
	monitorGroup.GET("/domain/list", domainController.GetDomainList)                     // 获取域名监控列表
	monitorGroup.POST("/domain/check/:id", domainController.CheckDomain)                 // 检查单个域名
	monitorGroup.POST("/domain/check-all", domainController.CheckAllDomains)             // 检查所有域名
	monitorGroup.GET("/domain/ssl-expiring", domainController.GetSSLExpiringSoon)        // 获取SSL证书即将过期的域名
	monitorGroup.GET("/domain/statistics", domainController.GetStatistics)               // 获取统计信息

	// 域名定时巡检配置
	monitorGroup.GET("/domain/schedule/config", domainScheduleController.GetScheduleConfig)     // 获取定时巡检配置
	monitorGroup.POST("/domain/schedule/config", domainScheduleController.UpdateScheduleConfig) // 更新定时巡检配置
	monitorGroup.POST("/domain/schedule/pause", domainScheduleController.PauseSchedule)         // 暂停定时巡检
	monitorGroup.POST("/domain/schedule/resume", domainScheduleController.ResumeSchedule)       // 恢复定时巡检
	monitorGroup.GET("/domain/schedule/status", domainScheduleController.GetScheduleStatus)     // 获取定时巡检状态

	// 故障管理
	monitorGroup.POST("/incident", controller.CreateIncident)                    // 创建故障记录
	monitorGroup.PUT("/incident", controller.UpdateIncident)                     // 更新故障记录
	monitorGroup.DELETE("/incident/:id", controller.DeleteIncident)              // 删除故障记录
	monitorGroup.GET("/incident/:id", controller.GetIncidentById)                // 获取故障详情
	monitorGroup.GET("/incident/list", controller.GetIncidentList)               // 获取故障列表
	monitorGroup.POST("/incident/batch-delete", controller.BatchDeleteIncidents) // 批量删除故障记录
	monitorGroup.PUT("/incident/:id/status", controller.UpdateIncidentStatus)    // 更新故障状态
	monitorGroup.GET("/incident/statistics", controller.GetIncidentStatistics)   // 获取故障统计信息

	// SSL证书管理
	monitorGroup.POST("/ssl-cert/apply", sslCertController.ApplySslCert)                     // 申请SSL证书
	monitorGroup.POST("/ssl-cert/upload", sslCertController.UploadSslCert)                   // 上传SSL证书
	monitorGroup.GET("/ssl-cert/:id/state", sslCertController.QueryCertState)                // 查询证书状态
	monitorGroup.GET("/ssl-cert/:id/download", sslCertController.DownloadCert)               // 下载证书
	monitorGroup.GET("/ssl-cert/list", sslCertController.ListSslCerts)                       // 获取SSL证书列表
	monitorGroup.DELETE("/ssl-cert/:id", sslCertController.DeleteSslCert)                    // 删除SSL证书
	monitorGroup.POST("/ssl-cert/:id/dns-validation", sslCertController.AddDnsValidation)    // 添加DNS验证记录
	monitorGroup.GET("/ssl-cert/:id/check-validation", sslCertController.CheckDnsValidation) // 检查DNS验证状态
	monitorGroup.POST("/ssl-cert/deploy", sslCertController.DeploySslCert)                   // 部署SSL证书到主机
	monitorGroup.GET("/ssl-cert/deploy-logs", sslCertController.GetDeployLogs)               // 获取SSL证书部署日志

	// 阿里云配置管理
	monitorGroup.POST("/aliyun-config", aliyunConfigController.CreateAliyunConfig)       // 创建阿里云配置
	monitorGroup.PUT("/aliyun-config", aliyunConfigController.UpdateAliyunConfig)        // 更新阿里云配置
	monitorGroup.GET("/aliyun-config/:id", aliyunConfigController.GetAliyunConfig)       // 获取阿里云配置详情
	monitorGroup.GET("/aliyun-config/list", aliyunConfigController.ListAliyunConfigs)    // 获取阿里云配置列表
	monitorGroup.DELETE("/aliyun-config/:id", aliyunConfigController.DeleteAliyunConfig) // 删除阿里云配置

	// 告警源管理
	monitorGroup.POST("/alert-sources", alertStatsController.CreateAlertSource)                    // 创建告警源
	monitorGroup.PUT("/alert-sources", alertStatsController.UpdateAlertSource)                     // 更新告警源
	monitorGroup.GET("/alert-sources/list", alertStatsController.GetAlertSourceList)               // 获取告警源列表（必须在:id之前）
	monitorGroup.POST("/alert-sources/batch-delete", alertStatsController.BatchDeleteAlertSources) // 批量删除告警源
	monitorGroup.GET("/alert-sources/:id", alertStatsController.GetAlertSourceById)                // 获取告警源详情
	monitorGroup.DELETE("/alert-sources/:id", alertStatsController.DeleteAlertSource)              // 删除告警源
	monitorGroup.PUT("/alert-sources/:id/status", alertStatsController.UpdateAlertSourceStatus)    // 更新告警源状态

	// 告警统计
	monitorGroup.GET("/alert-stats/today", alertStatsController.GetAlertStatsToday)                             // 获取今天的告警统计
	monitorGroup.GET("/alert-stats/yesterday", alertStatsController.GetAlertStatsYesterday)                     // 获取昨天的告警统计
	monitorGroup.GET("/alert-stats/last7days", alertStatsController.GetAlertStatsLast7Days)                     // 获取最近7天的告警统计
	monitorGroup.GET("/alert-stats/last30days", alertStatsController.GetAlertStatsLast30Days)                   // 获取最近30天的告警统计
	monitorGroup.GET("/alert-stats/today-vs-yesterday", alertStatsController.GetAlertStatsTodayVsYesterday)     // 今天vs昨天同环比
	monitorGroup.GET("/alert-stats/thisweek-vs-lastweek", alertStatsController.GetAlertStatsThisWeekVsLastWeek) // 本周vs上周同环比

	// 告警通知统计
	monitorGroup.GET("/alert-notifications/today", alertStatsController.GetAlertNotificationsToday)           // 获取今天的告警通知统计
	monitorGroup.GET("/alert-notifications/yesterday", alertStatsController.GetAlertNotificationsYesterday)   // 获取昨天的告警通知统计
	monitorGroup.GET("/alert-notifications/last7days", alertStatsController.GetAlertNotificationsLast7Days)   // 获取最近7天的告警通知统计
	monitorGroup.GET("/alert-notifications/last30days", alertStatsController.GetAlertNotificationsLast30Days) // 获取最近30天的告警通知统计

	// 协作空间故障统计
	monitorGroup.GET("/channel-stats/today", alertStatsController.GetChannelStatsToday)           // 获取今天的协作空间统计
	monitorGroup.GET("/channel-stats/yesterday", alertStatsController.GetChannelStatsYesterday)   // 获取昨天的协作空间统计
	monitorGroup.GET("/channel-stats/last7days", alertStatsController.GetChannelStatsLast7Days)   // 获取最近7天的协作空间统计
	monitorGroup.GET("/channel-stats/last30days", alertStatsController.GetChannelStatsLast30Days) // 获取最近30天的协作空间统计

	// 云账单管理
	monitorGroup.GET("/billing/aliyun-product", cloudBillingController.GetAliyunProductBilling)    // 获取阿里云账单（按产品）
	monitorGroup.GET("/billing/tencent", cloudBillingController.GetTencentBilling)                 // 获取腾讯云账单
	monitorGroup.GET("/billing/aliyun-tag", cloudBillingController.GetAliyunTagBilling)            // 获取阿里云账单（按标签）
	monitorGroup.GET("/billing/aws", cloudBillingController.GetAWSBilling)                         // 获取AWS账单
	monitorGroup.GET("/billing/last-month-summary", cloudBillingController.GetLastMonthSummary)    // 获取上月费用汇总及同环比
	monitorGroup.GET("/billing/last-12-months-trend", cloudBillingController.GetLast12MonthsTrend) // 获取最近12个月费用趋势

	// Zabbix告警统计
	monitorGroup.GET("/zabbix/alert-top-by-name", zabbixStatsController.GetZabbixAlertTopByName)           // 获取告警TOP(按告警名称)
	monitorGroup.GET("/zabbix/alert-top-by-host", zabbixStatsController.GetZabbixAlertTopByHost)           // 获取告警TOP(按主机)
	monitorGroup.GET("/zabbix/alert-top-by-hostgroup", zabbixStatsController.GetZabbixAlertTopByHostGroup) // 获取告警TOP(按主机群组)
	monitorGroup.GET("/zabbix/notification-stats", zabbixStatsController.GetZabbixNotificationStats)       // 获取每日告警数量统计(折线图)

	// 通知机器人管理（统一接口，支持钉钉/飞书/企业微信/邮箱/通用Webhook/Teams）
	monitorGroup.POST("/robot", notifyRobotController.CreateRobot)              // 创建通知机器人
	monitorGroup.GET("/robot/list", notifyRobotController.GetRobotList)         // 获取通知机器人列表（支持按type筛选）- 必须在:id之前
	monitorGroup.POST("/robot/test/:id", notifyRobotController.TestSendMessage) // 测试发送消息 - 必须在:id之前
	monitorGroup.GET("/robot/:id", notifyRobotController.GetRobot)              // 获取通知机器人详情
	monitorGroup.PUT("/robot", notifyRobotController.UpdateRobot)               // 更新通知机器人
	monitorGroup.DELETE("/robot/:id", notifyRobotController.DeleteRobot)        // 删除通知机器人

	// Webhook日志查询接口（需要认证）
	// 注意: /webhook/receive 接口在 router.go 的公开路由中注册（供外部系统回调）
	monitorGroup.GET("/webhook/log/list", webhookController.GetWebhookLogList)        // 获取Webhook日志列表 - 必须在:id之前
	monitorGroup.GET("/webhook/log/:id/notify-logs", webhookController.GetNotifyLogs) // 获取通知日志详情 - 必须在:id之前
	monitorGroup.GET("/webhook/log/:id", webhookController.GetWebhookLog)             // 获取Webhook日志详情
	// 用户统计
	monitorGroup.GET("/user-stats/login", userStatsController.GetUserLoginStats) // 获取用户登录统计

	// API监控管理
	// Note: /test-data endpoint is registered in router.go as a public endpoint (no JWT required)
	monitorGroup.POST("/business-functions", apiBusinessFunctionController.Create)                        // 创建业务功能
	monitorGroup.PUT("/business-functions/:id", apiBusinessFunctionController.Update)                     // 更新业务功能
	monitorGroup.DELETE("/business-functions/:id", apiBusinessFunctionController.Delete)                  // 删除业务功能
	monitorGroup.GET("/business-functions", apiBusinessFunctionController.GetList)                        // 获取业务功能列表
	monitorGroup.GET("/business-functions-by-group", apiBusinessFunctionController.GetFunctionsByGroupID) // 获取指定业务组的所有功能
	monitorGroup.GET("/business-functions/:id", apiBusinessFunctionController.GetByID)                    // 获取业务功能详情
	monitorGroup.POST("/api-monitors/simple", apiMonitorController.SimpleCreateMonitor)                   // 简化版创建API监控
	monitorGroup.POST("/api-monitors", apiMonitorController.CreateMonitor)                                // 创建API监控
	monitorGroup.PUT("/api-monitors/:id", apiMonitorController.UpdateMonitor)                             // 更新API监控
	monitorGroup.DELETE("/api-monitors/:id", apiMonitorController.DeleteMonitor)                          // 删除API监控
	monitorGroup.GET("/api-monitors/:id", apiMonitorController.GetMonitor)                                // 获取API监控详情
	monitorGroup.GET("/api-monitors", apiMonitorController.GetMonitorList)                                // 获取API监控列表
	monitorGroup.GET("/api-monitors-by-function", apiMonitorController.GetMonitorsByFunctionID)           // 获取指定业务功能的API列表
	monitorGroup.GET("/groups-tree-api-stats", apiMonitorController.GetGroupsTreeWithApiCount)            // 获取业务组树形结构及API统计（用于菜单显示）
	monitorGroup.PUT("/api-monitors/:id/status", apiMonitorController.UpdateMonitorStatus)                // 更新API监控状态
	monitorGroup.GET("/api-monitors/stats", apiMonitorController.GetStats)                                // 获取统计信息
	monitorGroup.GET("/api-monitors/group-stats", apiMonitorController.GetGroupStats)                     // 获取业务组统计
	monitorGroup.GET("/api-monitors/service-stats", apiMonitorController.GetServiceStats)                 // 获取服务统计
	monitorGroup.POST("/api-monitors/test", apiMonitorController.TestAPI)                                 // 测试API接口

	// API监控分层管理（必须在通用路由之后，避免路由冲突）
	monitorGroup.POST("/api-monitors/:id/add-child", apiMonitorController.AddChildMonitor)         // 添加子监控
	monitorGroup.GET("/api-monitors/:id/hierarchy", apiMonitorController.GetMonitorHierarchy)      // 获取监控树形结构
	monitorGroup.GET("/api-monitors/top-level", apiMonitorController.GetTopLevelMonitors)          // 获取顶级监控列表
	monitorGroup.DELETE("/api-monitors/:id/parent", apiMonitorController.RemoveParentRelationship) // 移除父监控关系

	monitorGroup.GET("/auth-test", apiMonitorController.AuthTestEndpoint) // 需要认证的测试接口

	// 系统测试 - 各模块健康检查接口
	// Note: /system-test/all and /system-test/monitor endpoints are registered in router.go as public endpoints (no JWT required)
	monitorGroup.GET("/system-test/cmdb", systemHealthController.CMDBTest)                 // CMDB模块测试
	monitorGroup.GET("/system-test/k8s", systemHealthController.K8STest)                   // K8S模块测试
	monitorGroup.GET("/system-test/db", systemHealthController.DBTest)                     // 数据库模块测试
	monitorGroup.GET("/system-test/app", systemHealthController.AppTest)                   // 应用发布模块测试
	monitorGroup.GET("/system-test/task", systemHealthController.TaskTest)                 // 智能任务模块测试
	monitorGroup.GET("/system-test/system", systemHealthController.SystemTest)             // 系统管理模块测试
	monitorGroup.GET("/system-test/configcenter", systemHealthController.ConfigCenterTest) // 配置中心模块测试

	// 告警规则管理
	alertGroup := monitorGroup.Group("/alert")
	{
		// 告警规则模板
		alertGroup.POST("/templates", alertTemplateController.CreateTemplate)                 // 创建规则模板
		alertGroup.PUT("/templates/:id", alertTemplateController.UpdateTemplate)              // 更新规则模板
		alertGroup.DELETE("/templates/:id", alertTemplateController.DeleteTemplate)           // 删除规则模板
		alertGroup.GET("/templates/:id", alertTemplateController.GetTemplate)                 // 获取规则模板详情
		alertGroup.GET("/templates", alertTemplateController.GetTemplateList)                 // 获取规则模板列表
		alertGroup.PUT("/templates/:id/status", alertTemplateController.UpdateTemplateStatus) // 更新模板状态
		alertGroup.GET("/templates/categories", alertTemplateController.GetCategories)        // 获取所有分类

		// 告警策略
		alertGroup.POST("/policies", alertPolicyController.CreatePolicy)                   // 创建告警策略
		alertGroup.PUT("/policies/:id", alertPolicyController.UpdatePolicy)                // 更新告警策略
		alertGroup.DELETE("/policies/:id", alertPolicyController.DeletePolicy)             // 删除告警策略
		alertGroup.GET("/policies/:id", alertPolicyController.GetPolicy)                   // 获取告警策略详情
		alertGroup.GET("/policies", alertPolicyController.GetPolicyList)                   // 获取告警策略列表
		alertGroup.PUT("/policies/:id/status", alertPolicyController.UpdatePolicyStatus)   // 更新策略状态
		alertGroup.POST("/policies/preview", alertPolicyController.PreviewPrometheusQuery) // Prometheus查询预览

		// 告警策略规则
		alertGroup.POST("/policy-rules", alertPolicyController.CreatePolicyRule)                 // 创建策略规则(从模板导入)
		alertGroup.GET("/policy-rules/:id", alertPolicyController.GetPolicyRule)                 // 获取单个策略规则详情
		alertGroup.PUT("/policy-rules/:id", alertPolicyController.UpdatePolicyRule)              // 更新策略规则
		alertGroup.PUT("/policy-rules/:id/status", alertPolicyController.UpdatePolicyRuleStatus) // 更新策略规则状态(启用/禁用)
		alertGroup.GET("/policy-rules/:id/preview", alertPolicyController.PreviewPolicyRule)     // 预览策略规则查询结果
		alertGroup.DELETE("/policy-rules/:id", alertPolicyController.DeletePolicyRule)           // 删除策略规则
		alertGroup.GET("/policy-rules", alertPolicyController.GetPolicyRules)                    // 获取策略下的规则列表

		// 告警规则（旧版，保持兼容）
		alertGroup.POST("/rules", alertRuleController.CreateRule)                 // 创建告警规则
		alertGroup.PUT("/rules/:id", alertRuleController.UpdateRule)              // 更新告警规则
		alertGroup.DELETE("/rules/:id", alertRuleController.DeleteRule)           // 删除告警规则
		alertGroup.GET("/rules/:id", alertRuleController.GetRule)                 // 获取告警规则详情
		alertGroup.GET("/rules", alertRuleController.GetRuleList)                 // 获取告警规则列表
		alertGroup.PUT("/rules/:id/status", alertRuleController.UpdateRuleStatus) // 更新规则状态(启用/禁用)

		// 告警记录
		alertGroup.GET("/records", alertRuleController.GetAlertRecords)                       // 获取告警记录列表
		alertGroup.GET("/records/:id", alertRuleController.GetAlertRecord)                    // 获取告警记录详情
		alertGroup.PUT("/records/:id/process", alertRuleController.ProcessAlertRecord)        // 统一处理告警记录
		alertGroup.PUT("/records/:id/handle", alertRuleController.HandleAlertRecord)          // 处理告警记录
		alertGroup.PUT("/records/:id/acknowledge", alertRuleController.AcknowledgeAlert)      // 认领告警
		alertGroup.DELETE("/records/:id", alertRuleController.DeleteAlertRecord)              // 删除单条告警记录
		alertGroup.POST("/records/batch-delete", alertRuleController.BatchDeleteAlertRecords) // 批量删除告警记录
		alertGroup.DELETE("/records/clear-all", alertRuleController.ClearAllAlertRecords)     // 清空所有告警记录

		// 告警统计
		alertGroup.GET("/stats", alertRuleController.GetAlertStats) // 获取告警统计
	}

}
