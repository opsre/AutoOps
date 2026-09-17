package controller

import (
	"net/http"
	"strconv"
	"time"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type AlertRuleController struct {
	alertRuleService service.AlertRuleService
}

func NewAlertRuleController() *AlertRuleController {
	return &AlertRuleController{
		alertRuleService: service.NewAlertRuleService(),
	}
}

// CreateRule 创建告警规则
// @Summary 创建告警规则
// @Description 创建新的告警规则
// @Tags 告警规则
// @Accept json
// @Produce json
// @Param rule body model.AlertRule true "告警规则"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/rules [post]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) CreateRule(c *gin.Context) {
	var rule model.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.alertRuleService.CreateRule(c, &rule)
}

// UpdateRule 更新告警规则
// @Summary 更新告警规则
// @Description 更新已有的告警规则
// @Tags 告警规则
// @Accept json
// @Produce json
// @Param id path uint true "规则ID"
// @Param rule body model.AlertRule true "告警规则"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/rules/{id} [put]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) UpdateRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	var rule model.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	rule.ID = uint(id)
	ctrl.alertRuleService.UpdateRule(c, &rule)
}

// DeleteRule 删除告警规则
// @Summary 删除告警规则
// @Description 删除指定的告警规则
// @Tags 告警规则
// @Accept json
// @Produce json
// @Param id path uint true "规则ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/rules/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) DeleteRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	ctrl.alertRuleService.DeleteRule(c, uint(id))
}

// GetRule 获取告警规则
// @Summary 获取告警规则
// @Description 获取指定ID的告警规则
// @Tags 告警规则
// @Accept json
// @Produce json
// @Param id path uint true "规则ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/rules/{id} [get]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) GetRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	ctrl.alertRuleService.GetRule(c, uint(id))
}

// GetRuleList 获取告警规则列表
// @Summary 获取告警规则列表
// @Description 获取告警规则列表，支持分页和筛选
// @Tags 告警规则
// @Accept json
// @Produce json
// @Param status query int false "状态(0:禁用,1:启用)"
// @Param groupId query uint false "业务组ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/rules [get]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) GetRuleList(c *gin.Context) {
	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	var groupID *uint
	if groupIDStr := c.Query("groupId"); groupIDStr != "" {
		g, _ := strconv.ParseUint(groupIDStr, 10, 32)
		gid := uint(g)
		groupID = &gid
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	ctrl.alertRuleService.GetRuleList(c, status, groupID, page, pageSize)
}

// UpdateRuleStatus 更新规则状态
// @Summary 更新规则状态
// @Description 启用或禁用告警规则
// @Tags 告警规则
// @Accept json
// @Produce json
// @Param id path uint true "规则ID"
// @Param status query int true "状态(0:禁用,1:启用)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/rules/{id}/status [put]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) UpdateRuleStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	statusStr := c.Query("status")
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的状态值")
		return
	}

	ctrl.alertRuleService.UpdateRuleStatus(c, uint(id), status)
}

// GetAlertRecords 获取告警记录列表
// @Summary 获取告警记录列表
// @Description 获取告警记录列表，支持多条件筛选
// @Tags 告警记录
// @Accept json
// @Produce json
// @Param status query string false "状态(firing/resolved)"
// @Param severity query string false "严重级别(1/2/3)"
// @Param ip query string false "IP地址(精确匹配)"
// @Param ruleId query uint false "规则ID"
// @Param groupId query uint false "业务组ID"
// @Param startTime query string false "开始时间(格式: 2006-01-02 15:04:05)"
// @Param endTime query string false "结束时间(格式: 2006-01-02 15:04:05)"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/records [get]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) GetAlertRecords(c *gin.Context) {
	status := c.Query("status")
	severity := c.Query("severity")
	ip := c.Query("ip")

	var ruleID *uint
	if ruleIDStr := c.Query("ruleId"); ruleIDStr != "" {
		r, _ := strconv.ParseUint(ruleIDStr, 10, 32)
		rid := uint(r)
		ruleID = &rid
	}

	var groupID *uint
	if groupIDStr := c.Query("groupId"); groupIDStr != "" {
		g, _ := strconv.ParseUint(groupIDStr, 10, 32)
		gid := uint(g)
		groupID = &gid
	}

	var startTime, endTime *time.Time
	if startTimeStr := c.Query("startTime"); startTimeStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTimeStr); err == nil {
			startTime = &t
		}
	}
	if endTimeStr := c.Query("endTime"); endTimeStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTimeStr); err == nil {
			endTime = &t
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	ctrl.alertRuleService.GetAlertRecords(c, status, severity, ip, ruleID, groupID, startTime, endTime, page, pageSize)
}

// GetAlertRecord 获取告警记录详情
// @Summary 获取告警记录详情
// @Description 获取指定ID的告警记录详情
// @Tags 告警记录
// @Accept json
// @Produce json
// @Param id path uint true "记录ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/records/{id} [get]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) GetAlertRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	ctrl.alertRuleService.GetAlertRecord(c, uint(id))
}

// GetAlertStats 获取告警统计
// @Summary 获取告警统计
// @Description 获取告警统计信息
// @Tags 告警记录
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/stats [get]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) GetAlertStats(c *gin.Context) {
	ctrl.alertRuleService.GetAlertStats(c)
}

// ProcessAlertRecord 统一处理告警记录
// @Summary 统一处理告警记录
// @Description 支持认领自己、转派给其他系统用户、记录处理步骤
// @Tags 告警记录
// @Accept json
// @Produce json
// @Param id path uint true "记录ID"
// @Param data body object true "处理信息 {\"action\":\"acknowledge|assign|handle\",\"assigneeId\":123,\"handleNote\":\"处理说明\"}"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/records/{id}/process [put]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) ProcessAlertRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	var req struct {
		Action     string `json:"action" binding:"required,oneof=acknowledge assign handle"`
		AssigneeID *uint  `json:"assigneeId"`
		HandleNote string `json:"handleNote"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.alertRuleService.ProcessAlertRecord(c, uint(id), req.Action, req.AssigneeID, req.HandleNote)
}

// HandleAlertRecord 处理告警记录
// @Summary 处理告警记录
// @Description 兼容旧接口，使用当前登录用户记录处理步骤
// @Tags 告警记录
// @Accept json
// @Produce json
// @Param id path uint true "记录ID"
// @Param data body object true "处理信息 {\"assigneeId\":123,\"handleNote\":\"处理备注\"}"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/records/{id}/handle [put]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) HandleAlertRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	var req struct {
		AssigneeID *uint  `json:"assigneeId"`
		HandleNote string `json:"handleNote"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.alertRuleService.ProcessAlertRecord(c, uint(id), "handle", req.AssigneeID, req.HandleNote)
}

// AcknowledgeAlert 认领告警
// @Summary 认领告警
// @Description 兼容旧接口，默认由当前登录用户认领
// @Tags 告警记录
// @Accept json
// @Produce json
// @Param id path uint true "记录ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/records/{id}/acknowledge [put]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) AcknowledgeAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	ctrl.alertRuleService.ProcessAlertRecord(c, uint(id), "acknowledge", nil, "")
}

// DeleteAlertRecord 删除单条告警记录
// @Summary 删除单条告警记录
// @Description 删除指定ID的告警记录
// @Tags 告警记录
// @Accept json
// @Produce json
// @Param id path uint true "记录ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/records/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) DeleteAlertRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	ctrl.alertRuleService.DeleteAlertRecord(c, uint(id))
}

// BatchDeleteAlertRecords 批量删除告警记录
// @Summary 批量删除告警记录
// @Description 批量删除指定ID列表的告警记录
// @Tags 告警记录
// @Accept json
// @Produce json
// @Param ids body object true "记录ID列表 {\"ids\":[1,2,3]}"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/records/batch-delete [post]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) BatchDeleteAlertRecords(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.alertRuleService.BatchDeleteAlertRecords(c, req.IDs)
}

// ClearAllAlertRecords 清空所有告警记录
// @Summary 清空所有告警记录
// @Description 清空所有告警历史记录（慎用）
// @Tags 告警记录
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/records/clear-all [delete]
// @Security ApiKeyAuth
func (ctrl *AlertRuleController) ClearAllAlertRecords(c *gin.Context) {
	ctrl.alertRuleService.ClearAllAlertRecords(c)
}
