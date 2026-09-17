package controller

import (
	"strconv"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// WebhookController Webhook控制器
type WebhookController struct {
	webhookService service.WebhookService
}

// NewWebhookController 创建Webhook控制器实例
func NewWebhookController() *WebhookController {
	return &WebhookController{
		webhookService: service.NewWebhookService(),
	}
}

// @Tags Webhook
// @Summary 接收Webhook通知
// @Description 接收Webhook通知
// @Accept json
// @Produce json
// @Description 接收外部系统的Webhook通知并转发到配置的通知机器人
// @Param data body model.WebhookRequest true "Webhook请求数据"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/webhook/receive [post]
func (ctrl *WebhookController) ReceiveWebhook(c *gin.Context) {
	var req model.WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	ctrl.webhookService.ReceiveWebhook(c, &req)
}

// @Tags Webhook
// @Summary 获取Webhook日志详情
// @Description 获取Webhook日志详情
// @Produce json
// @Description 根据ID获取Webhook日志详情
// @Param id path int true "日志ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/webhook/log/{id} [get]
// @Security ApiKeyAuth
func (ctrl *WebhookController) GetWebhookLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.webhookService.GetLogByID(c, uint(id))
}

// @Tags Webhook
// @Summary 获取Webhook日志列表
// @Description 获取Webhook日志列表
// @Produce json
// @Description 分页获取Webhook日志列表
// @Param source query string false "来源标识"
// @Param status query string false "状态(success/failed/partial)"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/webhook/log/list [get]
// @Security ApiKeyAuth
func (ctrl *WebhookController) GetWebhookLogList(c *gin.Context) {
	source := c.Query("source")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	ctrl.webhookService.GetLogList(c, source, status, page, pageSize)
}

// @Tags Webhook
// @Summary 获取Webhook通知日志
// @Description 获取Webhook通知日志
// @Produce json
// @Description 获取指定Webhook日志的通知详情列表
// @Param webhookLogId path int true "Webhook日志ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/webhook/log/{webhookLogId}/notify-logs [get]
// @Security ApiKeyAuth
func (ctrl *WebhookController) GetNotifyLogs(c *gin.Context) {
	webhookLogID, err := strconv.ParseUint(c.Param("webhookLogId"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.webhookService.GetNotifyLogs(c, uint(webhookLogID))
}
