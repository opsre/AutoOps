package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type AlertRuleTemplateController struct {
	templateService service.AlertRuleTemplateService
}

func NewAlertRuleTemplateController() *AlertRuleTemplateController {
	return &AlertRuleTemplateController{
		templateService: service.NewAlertRuleTemplateService(),
	}
}

// CreateTemplate 创建规则模板
// @Summary 创建规则模板
// @Description 创建新的告警规则模板
// @Tags 告警规则模板
// @Accept json
// @Produce json
// @Param template body model.CreateRuleTemplateDto true "规则模板"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/templates [post]
// @Security ApiKeyAuth
func (ctrl *AlertRuleTemplateController) CreateTemplate(c *gin.Context) {
	var dto model.CreateRuleTemplateDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.templateService.CreateTemplate(c, &dto)
}

// UpdateTemplate 更新规则模板
// @Summary 更新规则模板
// @Description 更新已有的告警规则模板
// @Tags 告警规则模板
// @Accept json
// @Produce json
// @Param id path uint true "模板ID"
// @Param template body model.UpdateRuleTemplateDto true "规则模板"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/templates/{id} [put]
// @Security ApiKeyAuth
func (ctrl *AlertRuleTemplateController) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	var dto model.UpdateRuleTemplateDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	dto.ID = uint(id)
	ctrl.templateService.UpdateTemplate(c, &dto)
}

// DeleteTemplate 删除规则模板
// @Summary 删除规则模板
// @Description 删除指定的告警规则模板
// @Tags 告警规则模板
// @Accept json
// @Produce json
// @Param id path uint true "模板ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/templates/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *AlertRuleTemplateController) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	ctrl.templateService.DeleteTemplate(c, uint(id))
}

// GetTemplate 获取规则模板
// @Summary 获取规则模板
// @Description 获取指定ID的告警规则模板
// @Tags 告警规则模板
// @Accept json
// @Produce json
// @Param id path uint true "模板ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/templates/{id} [get]
// @Security ApiKeyAuth
func (ctrl *AlertRuleTemplateController) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	ctrl.templateService.GetTemplate(c, uint(id))
}

// GetTemplateList 获取规则模板列表
// @Summary 获取规则模板列表
// @Description 获取告警规则模板列表，支持分页和筛选
// @Tags 告警规则模板
// @Accept json
// @Produce json
// @Param category query string false "分类(cpu/memory/load/disk/io/network/connection/service)"
// @Param status query int false "状态(0:禁用,1:启用)"
// @Param keyword query string false "关键字(名称或描述)"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/templates [get]
// @Security ApiKeyAuth
func (ctrl *AlertRuleTemplateController) GetTemplateList(c *gin.Context) {
	category := c.Query("category")
	keyword := c.Query("keyword")

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	ctrl.templateService.GetTemplateList(c, category, status, keyword, page, pageSize)
}

// UpdateTemplateStatus 更新规则模板状态
// @Summary 更新规则模板状态
// @Description 启用或禁用告警规则模板
// @Tags 告警规则模板
// @Accept json
// @Produce json
// @Param id path uint true "模板ID"
// @Param status query int true "状态(0:禁用,1:启用)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/templates/{id}/status [put]
// @Security ApiKeyAuth
func (ctrl *AlertRuleTemplateController) UpdateTemplateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	statusStr := c.Query("status")
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的状态值")
		return
	}

	ctrl.templateService.UpdateTemplateStatus(c, uint(id), status)
}

// GetCategories 获取所有分类
// @Summary 获取所有分类
// @Description 获取所有告警规则模板分类
// @Tags 告警规则模板
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/templates/categories [get]
// @Security ApiKeyAuth
func (ctrl *AlertRuleTemplateController) GetCategories(c *gin.Context) {
	ctrl.templateService.GetAllCategories(c)
}
