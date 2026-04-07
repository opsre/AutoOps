package monitor

import (
	"dodevops-api/api/monitor/controller"
	"dodevops-api/middleware"

	"github.com/gin-gonic/gin"
)

func InitMonitorRouter(r *gin.RouterGroup) {
	monitorController := controller.NewMonitorController()
	agentController := controller.NewAgentController()
	alertController := controller.NewAlertController() // 添加告警控制器

	// Webhook入口
	// PrometheusAlert webhooks (no auth usually needed)
	publicMonitorGroup := r.Group("/monitor")
	publicMonitorGroup.POST("/alert/webhook/gitlab", alertController.ReceiveGitlabWebhook)
	publicMonitorGroup.POST("/alert/webhook/zabbix", alertController.ReceiveZabbixWebhook)
	publicMonitorGroup.POST("/alert/webhook/prometheus", alertController.ReceivePrometheusWebhook)

	monitorGroup := r.Group("/monitor")
	monitorGroup.Use(middleware.AuthMiddleware())

	// PrometheusAlert
	// Template CRUD
	monitorGroup.POST("/alert/template", alertController.CreateTemplate)
	monitorGroup.DELETE("/alert/template/:id", alertController.DeleteTemplate)
	monitorGroup.PUT("/alert/template/:id", alertController.UpdateTemplate)
	monitorGroup.GET("/alert/templates", alertController.GetTemplate)
	// Router CRUD
	// monitorGroup.POST("/alert/router", alertController.CreateRouter)
	// monitorGroup.DELETE("/alert/router/:id", alertController.DeleteRouter)
	// monitorGroup.PUT("/alert/router/:id", alertController.UpdateRouter)
	// monitorGroup.GET("/alert/routers", alertController.GetRouter)
	// Component Config Operations
	// monitorGroup.POST("/alert/router/reload", alertController.ReloadConfig)
	// monitorGroup.GET("/alert/router/health", alertController.HealthCheck)
	// Records Operations
	monitorGroup.GET("/alert/records", alertController.GetRecords)
	monitorGroup.DELETE("/alert/records/clean", alertController.CleanRecords)

	// 主机监控
	monitorGroup.GET("/host/:id", monitorController.GetHostMetrics)                        // 获取主机监控数据
	monitorGroup.GET("/hosts", monitorController.BatchGetHostMetrics)                      // 批量获取主机监控数据(包含在线状态)
	monitorGroup.GET("/hosts/:id/history", monitorController.GetHostMetricHistory)         // 获取主机指定指标的历史数据
	monitorGroup.GET("/hosts/:id/all-metrics", monitorController.GetHostAllMetricsHistory) // 获取主机所有指标的历史数据
	monitorGroup.GET("/hosts/:id/top-processes", monitorController.GetTopProcesses)        // 获取主机TOP进程使用率
	monitorGroup.GET("/hosts/:id/ports", monitorController.GetHostPorts)                   // 获取主机端口信息

	// Agent管理
	monitorGroup.POST("/agent/deploy", agentController.DeployAgent)           // 部署agent到指定主机(支持单个或多个)
	monitorGroup.DELETE("/agent/uninstall", agentController.UninstallAgent)   // 卸载指定主机的agent(支持单个或多个)
	monitorGroup.GET("/agent/status/:id", agentController.GetAgentStatus)     // 根据主机id获取agent状态
	monitorGroup.POST("/agent/restart/:id", agentController.RestartAgent)     // 重启agent
	monitorGroup.GET("/agent/list", agentController.GetAgentList)             // 获取agent列表
	monitorGroup.GET("/agent/statistics", agentController.GetAgentStatistics) // 获取统计信息
	monitorGroup.DELETE("/agent/delete/:id", agentController.DeleteAgent)     // 删除agent数据(用于离线服务器)
}
