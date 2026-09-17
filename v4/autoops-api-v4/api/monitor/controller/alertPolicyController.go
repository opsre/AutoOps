package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type AlertPolicyController struct {
	policyService service.AlertPolicyService
}

func NewAlertPolicyController() *AlertPolicyController {
	return &AlertPolicyController{
		policyService: service.NewAlertPolicyService(),
	}
}

// CreatePolicy 创建告警策略
// @Summary 创建告警策略
// @Description 创建新的告警策略
// @Tags 告警策略
// @Accept json
// @Produce json
// @Param policy body model.CreateAlertPolicyDto true "告警策略"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policies [post]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) CreatePolicy(c *gin.Context) {
	var dto model.CreateAlertPolicyDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.policyService.CreatePolicy(c, &dto)
}

// UpdatePolicy 更新告警策略
// @Summary 更新告警策略
// @Description 更新已有的告警策略
// @Tags 告警策略
// @Accept json
// @Produce json
// @Param id path uint true "策略ID"
// @Param policy body model.UpdateAlertPolicyDto true "告警策略"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policies/{id} [put]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) UpdatePolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的策略ID")
		return
	}

	var dto model.UpdateAlertPolicyDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	dto.ID = uint(id)
	ctrl.policyService.UpdatePolicy(c, &dto)
}

// DeletePolicy 删除告警策略
// @Summary 删除告警策略
// @Description 删除指定的告警策略
// @Tags 告警策略
// @Accept json
// @Produce json
// @Param id path uint true "策略ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policies/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) DeletePolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的策略ID")
		return
	}

	ctrl.policyService.DeletePolicy(c, uint(id))
}

// GetPolicy 获取告警策略
// @Summary 获取告警策略
// @Description 获取指定ID的告警策略
// @Tags 告警策略
// @Accept json
// @Produce json
// @Param id path uint true "策略ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policies/{id} [get]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) GetPolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的策略ID")
		return
	}

	ctrl.policyService.GetPolicy(c, uint(id))
}

// GetPolicyList 获取告警策略列表
// @Summary 获取告警策略列表
// @Description 获取告警策略列表，支持分页和筛选
// @Tags 告警策略
// @Accept json
// @Produce json
// @Param groupId query uint false "业务组ID"
// @Param status query int false "状态(0:禁用,1:启用)"
// @Param keyword query string false "关键字(名称或描述)"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policies [get]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) GetPolicyList(c *gin.Context) {
	keyword := c.Query("keyword")

	var groupID *uint
	if groupIDStr := c.Query("groupId"); groupIDStr != "" {
		g, _ := strconv.ParseUint(groupIDStr, 10, 32)
		gid := uint(g)
		groupID = &gid
	}

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	ctrl.policyService.GetPolicyList(c, groupID, status, keyword, page, pageSize)
}

// UpdatePolicyStatus 更新策略状态
// @Summary 更新策略状态
// @Description 启用或禁用告警策略
// @Tags 告警策略
// @Accept json
// @Produce json
// @Param id path uint true "策略ID"
// @Param status query int true "状态(0:禁用,1:启用)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policies/{id}/status [put]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) UpdatePolicyStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的策略ID")
		return
	}

	statusStr := c.Query("status")
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的状态值")
		return
	}

	ctrl.policyService.UpdatePolicyStatus(c, uint(id), status)
}

// PreviewPrometheusQuery 预览Prometheus查询
// @Summary 预览Prometheus查询
// @Description 动态查询Prometheus数据预览
// @Tags 告警策略
// @Accept json
// @Produce json
// @Param query body model.PrometheusQueryPreviewDto true "查询参数"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policies/preview [post]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) PreviewPrometheusQuery(c *gin.Context) {
	var dto model.PrometheusQueryPreviewDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.policyService.PreviewPrometheusQuery(c, &dto)
}

// CreatePolicyRule 为策略添加告警规则(从模板导入)
// @Summary 为策略添加告警规则
// @Description 从模板导入告警规则到策略中
// @Tags 告警策略规则
// @Accept json
// @Produce json
// @Param rule body model.CreatePolicyRuleDto true "策略规则"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policy-rules [post]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) CreatePolicyRule(c *gin.Context) {
	var dto model.CreatePolicyRuleDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.policyService.CreatePolicyRule(c, &dto)
}

// GetPolicyRule 获取单个策略规则详情
// @Summary 获取单个策略规则详情
// @Description 根据规则ID获取策略规则详情
// @Tags 告警策略规则
// @Accept json
// @Produce json
// @Param id path uint true "规则ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policy-rules/{id} [get]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) GetPolicyRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	ctrl.policyService.GetPolicyRule(c, uint(id))
}

// UpdatePolicyRule 更新策略规则
// @Summary 更新策略规则
// @Description 更新策略下的告警规则配置
// @Tags 告警策略规则
// @Accept json
// @Produce json
// @Param id path uint true "规则ID"
// @Param rule body model.UpdatePolicyRuleDto true "策略规则"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policy-rules/{id} [put]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) UpdatePolicyRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	var dto model.UpdatePolicyRuleDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	dto.ID = uint(id)
	ctrl.policyService.UpdatePolicyRule(c, &dto)
}

// DeletePolicyRule 删除策略规则
// @Summary 删除策略规则
// @Description 删除策略下的告警规则
// @Tags 告警策略规则
// @Accept json
// @Produce json
// @Param id path uint true "规则ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policy-rules/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) DeletePolicyRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	ctrl.policyService.DeletePolicyRule(c, uint(id))
}

// GetPolicyRules 获取策略下的规则列表
// @Summary 获取策略下的规则列表
// @Description 获取指定策略下的所有告警规则
// @Tags 告警策略规则
// @Accept json
// @Produce json
// @Param policyId query uint true "策略ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policy-rules [get]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) GetPolicyRules(c *gin.Context) {
	policyIDStr := c.Query("policyId")
	policyID, err := strconv.ParseUint(policyIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的策略ID")
		return
	}

	category := c.Query("category") // 可选的类型过滤参数
	ctrl.policyService.GetPolicyRules(c, uint(policyID), category)
}

// UpdatePolicyRuleStatus 更新策略规则状态
// @Summary 更新策略规则状态
// @Description 启用或禁用策略规则
// @Tags 告警策略规则
// @Accept json
// @Produce json
// @Param id path uint true "规则ID"
// @Param status query int true "状态(0:禁用,1:启用)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/alert/policy-rules/{id}/status [put]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) UpdatePolicyRuleStatus(c *gin.Context) {
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

	ctrl.policyService.UpdatePolicyRuleStatus(c, uint(id), status)
}

// PreviewPolicyRule 预览策略规则查询结果
// @Summary 预览策略规则查询结果
// @Description 预览指定策略规则的PromQL查询结果，显示当前值和是否超过阈值
// @Tags 告警策略规则
// @Accept json
// @Produce json
// @Param id path uint true "规则ID"
// @Success 200 {object} result.Result{data=model.PolicyRulePreviewResponse}
// @Router /api/v1/monitor/alert/policy-rules/{id}/preview [get]
// @Security ApiKeyAuth
func (ctrl *AlertPolicyController) PreviewPolicyRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	ctrl.policyService.PreviewPolicyRule(c, uint(id))
}

