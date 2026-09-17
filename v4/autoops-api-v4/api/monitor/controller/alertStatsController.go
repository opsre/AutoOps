package controller

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AlertStatsController struct {
	alertStatsService service.AlertStatsServiceInterface
}

func NewAlertStatsController() *AlertStatsController {
	return &AlertStatsController{
		alertStatsService: service.GetAlertStatsService(),
	}
}

// CreateAlertSource 创建告警源
// @Summary 创建告警源
// @Description 创建新的告警源配置
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param dto body model.CreateAlertSourceDto true "告警源信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-sources [post]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) CreateAlertSource(c *gin.Context) {
	var dto model.CreateAlertSourceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, err.Error())
		return
	}
	ctrl.alertStatsService.CreateAlertSource(c, &dto)
}

// UpdateAlertSource 更新告警源
// @Summary 更新告警源
// @Description 更新告警源配置
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param dto body model.UpdateAlertSourceDto true "告警源信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-sources [put]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) UpdateAlertSource(c *gin.Context) {
	var dto model.UpdateAlertSourceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, err.Error())
		return
	}
	ctrl.alertStatsService.UpdateAlertSource(c, &dto)
}

// DeleteAlertSource 删除告警源
// @Summary 删除告警源
// @Description 删除指定的告警源
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param id path uint true "告警源ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-sources/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) DeleteAlertSource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的ID")
		return
	}
	ctrl.alertStatsService.DeleteAlertSource(c, uint(id))
}

// GetAlertSourceById 根据ID获取告警源
// @Summary 根据ID获取告警源
// @Description 获取指定ID的告警源详情
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param id path uint true "告警源ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-sources/{id} [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertSourceById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的ID")
		return
	}
	ctrl.alertStatsService.GetAlertSourceById(c, uint(id))
}

// GetAlertSourceList 获取告警源列表
// @Summary 获取告警源列表
// @Description 获取告警源列表（支持分页和条件查询）
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param type query int false "告警源类型"
// @Param status query int false "状态"
// @Param keyword query string false "关键字"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-sources/list [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertSourceList(c *gin.Context) {
	var query model.AlertSourceQueryDto

	// 从query参数获取数据
	if typeStr := c.Query("type"); typeStr != "" {
		if typeVal, err := strconv.Atoi(typeStr); err == nil {
			query.Type = typeVal
		}
	}

	// 只有明确传递status参数时才设置
	if statusStr := c.Query("status"); statusStr != "" {
		if statusVal, err := strconv.Atoi(statusStr); err == nil {
			query.Status = &statusVal
		}
	}

	query.Keyword = c.Query("keyword")

	if pageNumStr := c.Query("pageNum"); pageNumStr != "" {
		if pageNum, err := strconv.Atoi(pageNumStr); err == nil {
			query.PageNum = pageNum
		}
	}

	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			query.PageSize = pageSize
		}
	}

	ctrl.alertStatsService.GetAlertSourceList(c, &query)
}

// BatchDeleteAlertSources 批量删除告警源
// @Summary 批量删除告警源
// @Description 批量删除指定的告警源
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param ids body []uint true "告警源ID列表"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-sources/batch-delete [post]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) BatchDeleteAlertSources(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		result.Failed(c, http.StatusBadRequest, err.Error())
		return
	}
	ctrl.alertStatsService.BatchDeleteAlertSources(c, ids)
}

// UpdateAlertSourceStatus 更新告警源状态
// @Summary 更新告警源状态
// @Description 启用或禁用告警源
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param id path uint true "告警源ID"
// @Param status body int true "状态(0=禁用,1=启用)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-sources/{id}/status [put]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) UpdateAlertSourceStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var statusDto struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&statusDto); err != nil {
		result.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	ctrl.alertStatsService.UpdateAlertSourceStatus(c, uint(id), statusDto.Status)
}

// GetAlertStatsToday 获取今天的告警统计
// @Summary 获取今天的告警统计
// @Description 获取今天的告警统计Top20
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Param label query string true "统计维度标签(如:check)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-stats/today [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertStatsToday(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	label := c.Query("label")
	if label == "" {
		result.Failed(c, http.StatusBadRequest, "请提供统计维度标签")
		return
	}

	ctrl.alertStatsService.GetAlertStatsToday(c, uint(sourceID), label)
}

// GetAlertStatsYesterday 获取昨天的告警统计
// @Summary 获取昨天的告警统计
// @Description 获取昨天的告警统计Top20
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Param label query string true "统计维度标签(如:check)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-stats/yesterday [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertStatsYesterday(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	label := c.Query("label")
	if label == "" {
		result.Failed(c, http.StatusBadRequest, "请提供统计维度标签")
		return
	}

	ctrl.alertStatsService.GetAlertStatsYesterday(c, uint(sourceID), label)
}

// GetAlertStatsLast7Days 获取最近7天的告警统计
// @Summary 获取最近7天的告警统计
// @Description 获取最近7天的告警统计Top20
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Param label query string true "统计维度标签(如:check)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-stats/last7days [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertStatsLast7Days(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	label := c.Query("label")
	if label == "" {
		result.Failed(c, http.StatusBadRequest, "请提供统计维度标签")
		return
	}

	ctrl.alertStatsService.GetAlertStatsLast7Days(c, uint(sourceID), label)
}

// GetAlertStatsLast30Days 获取最近30天的告警统计
// @Summary 获取最近30天的告警统计
// @Description 获取最近30天的告警统计Top20
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Param label query string true "统计维度标签(如:check)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-stats/last30days [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertStatsLast30Days(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	label := c.Query("label")
	if label == "" {
		result.Failed(c, http.StatusBadRequest, "请提供统计维度标签")
		return
	}

	ctrl.alertStatsService.GetAlertStatsLast30Days(c, uint(sourceID), label)
}

// GetAlertStatsTodayVsYesterday 今天vs昨天同环比
// @Summary 今天vs昨天同环比
// @Description 获取今天与昨天的告警统计对比
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Param label query string true "统计维度标签(如:check)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-stats/today-vs-yesterday [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertStatsTodayVsYesterday(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	label := c.Query("label")
	if label == "" {
		result.Failed(c, http.StatusBadRequest, "请提供统计维度标签")
		return
	}

	ctrl.alertStatsService.GetAlertStatsTodayVsYesterday(c, uint(sourceID), label)
}

// GetAlertStatsThisWeekVsLastWeek 本周vs上周同环比
// @Summary 本周vs上周同环比
// @Description 获取本周与上周的告警统计对比
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Param label query string true "统计维度标签(如:check)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-stats/thisweek-vs-lastweek [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertStatsThisWeekVsLastWeek(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	label := c.Query("label")
	if label == "" {
		result.Failed(c, http.StatusBadRequest, "请提供统计维度标签")
		return
	}

	ctrl.alertStatsService.GetAlertStatsThisWeekVsLastWeek(c, uint(sourceID), label)
}

// GetAlertNotificationsToday 获取今天的告警通知统计
// @Summary 获取今天的告警通知统计
// @Description 获取今天的告警通知统计(短信、邮件、语音)
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-notifications/today [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertNotificationsToday(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	ctrl.alertStatsService.GetAlertNotificationsToday(c, uint(sourceID))
}

// GetAlertNotificationsYesterday 获取昨天的告警通知统计
// @Summary 获取昨天的告警通知统计
// @Description 获取昨天的告警通知统计(短信、邮件、语音)
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-notifications/yesterday [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertNotificationsYesterday(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	ctrl.alertStatsService.GetAlertNotificationsYesterday(c, uint(sourceID))
}

// GetAlertNotificationsLast7Days 获取最近7天的告警通知统计
// @Summary 获取最近7天的告警通知统计
// @Description 获取最近7天的告警通知统计(短信、邮件、语音)
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-notifications/last7days [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertNotificationsLast7Days(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	ctrl.alertStatsService.GetAlertNotificationsLast7Days(c, uint(sourceID))
}

// GetAlertNotificationsLast30Days 获取最近30天的告警通知统计
// @Summary 获取最近30天的告警通知统计
// @Description 获取最近30天的告警通知统计(短信、邮件、语音)
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert-notifications/last30days [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetAlertNotificationsLast30Days(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	ctrl.alertStatsService.GetAlertNotificationsLast30Days(c, uint(sourceID))
}

// GetChannelStatsToday 获取今天的协作空间统计
// @Summary 获取今天的协作空间统计
// @Description 获取今天的协作空间统计(故障数量、自动关闭率、认领率等)
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/channel-stats/today [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetChannelStatsToday(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	ctrl.alertStatsService.GetChannelStatsToday(c, uint(sourceID))
}

// GetChannelStatsYesterday 获取昨天的协作空间统计
// @Summary 获取昨天的协作空间统计
// @Description 获取昨天的协作空间统计(故障数量、自动关闭率、认领率等)
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/channel-stats/yesterday [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetChannelStatsYesterday(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	ctrl.alertStatsService.GetChannelStatsYesterday(c, uint(sourceID))
}

// GetChannelStatsLast7Days 获取最近7天的协作空间统计
// @Summary 获取最近7天的协作空间统计
// @Description 获取最近7天的协作空间统计(故障数量、自动关闭率、认领率等)
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/channel-stats/last7days [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetChannelStatsLast7Days(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	ctrl.alertStatsService.GetChannelStatsLast7Days(c, uint(sourceID))
}

// GetChannelStatsLast30Days 获取最近30天的协作空间统计
// @Summary 获取最近30天的协作空间统计
// @Description 获取最近30天的协作空间统计(故障数量、自动关闭率、认领率等)
// @Tags 告警统计
// @Accept json
// @Produce json
// @Param sourceId query uint true "告警源ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/channel-stats/last30days [get]
// @Security ApiKeyAuth
func (ctrl *AlertStatsController) GetChannelStatsLast30Days(c *gin.Context) {
	sourceIDStr := c.Query("sourceId")
	sourceID, err := strconv.ParseUint(sourceIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的告警源ID")
		return
	}

	ctrl.alertStatsService.GetChannelStatsLast30Days(c, uint(sourceID))
}


