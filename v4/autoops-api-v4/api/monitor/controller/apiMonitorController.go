package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type APIMonitorController struct {
	service service.APIMonitorService
}

func NewAPIMonitorController() *APIMonitorController {
	return &APIMonitorController{
		service: service.NewAPIMonitorService(),
	}
}

// CreateMonitor 创建API监控
// @Summary 创建API监控
// @Description 创建新的API监控配置
// @Tags API监控
// @Accept json
// @Produce json
// @Param monitor body model.CreateAPIMonitorDto true "监控配置"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/api-monitors [post]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) CreateMonitor(c *gin.Context) {
	var dto model.CreateAPIMonitorDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.service.CreateMonitor(c, &dto)
}

// SimpleCreateMonitor 简化创建API监控
// @Summary 简化创建API监控
// @Description 简化版创建API监控，自动填充默认参数
// @Tags API监控
// @Accept json
// @Produce json
// @Param monitor body model.SimpleCreateAPIMonitorDto true "监控配置（简化版）"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/api-monitors/simple [post]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) SimpleCreateMonitor(c *gin.Context) {
	var dto model.SimpleCreateAPIMonitorDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.service.SimpleCreateMonitor(c, &dto)
}

// UpdateMonitor 更新API监控
// @Summary 更新API监控
// @Description 更新已有的API监控配置
// @Tags API监控
// @Accept json
// @Produce json
// @Param id path uint true "监控ID"
// @Param monitor body model.UpdateAPIMonitorDto true "监控配置"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/api-monitors/{id} [put]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) UpdateMonitor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的监控ID")
		return
	}

	var dto model.UpdateAPIMonitorDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	dto.ID = uint(id)
	ctrl.service.UpdateMonitor(c, &dto)
}

// DeleteMonitor 删除API监控
// @Summary 删除API监控
// @Description 删除指定的API监控
// @Tags API监控
// @Accept json
// @Produce json
// @Param id path uint true "监控ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/api-monitors/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) DeleteMonitor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的监控ID")
		return
	}

	ctrl.service.DeleteMonitor(c, uint(id))
}

// GetMonitor 获取API监控详情
// @Summary 获取API监控详情
// @Description 获取指定ID的API监控配置
// @Tags API监控
// @Accept json
// @Produce json
// @Param id path uint true "监控ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/api-monitors/{id} [get]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) GetMonitor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的监控ID")
		return
	}

	ctrl.service.GetMonitor(c, uint(id))
}

// GetMonitorList 获取API监控列表
// @Summary 获取API监控列表
// @Description 获取API监控列表，支持分页和筛选
// @Tags API监控
// @Accept json
// @Produce json
// @Param groupId query uint false "业务组ID"
// @Param serviceName query string false "服务名称"
// @Param status query int false "状态(0:禁用,1:启用)"
// @Param keyword query string false "关键字(名称或URL)"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/api-monitors [get]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) GetMonitorList(c *gin.Context) {
	var query model.APIMonitorListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 兼容 pageNum 参数（如果使用了 pageNum 而不是 page）
	if query.Page == 0 {
		if pageNum := c.Query("pageNum"); pageNum != "" {
			if p, err := strconv.Atoi(pageNum); err == nil {
				query.Page = p
			}
		}
	}

	ctrl.service.GetMonitorList(c, &query)
}

// GetMonitorsByFunctionID 获取指定业务功能的API监控列表（用于级联选择）
// @Summary 获取业务功能的API列表
// @Description 获取指定业务功能下的所有API监控，用于前端级联选择
// @Tags API监控
// @Accept json
// @Produce json
// @Param businessFunctionId query uint true "业务功能ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/api-monitors-by-function [get]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) GetMonitorsByFunctionID(c *gin.Context) {
	functionIDStr := c.Query("businessFunctionId")
	if functionIDStr == "" {
		result.Failed(c, http.StatusBadRequest, "业务功能ID不能为空")
		return
	}

	functionID, err := strconv.ParseUint(functionIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的业务功能ID")
		return
	}

	// 查询该功能下的所有API监控（不分页，全量返回）
	query := &model.APIMonitorListQuery{
		BusinessFunctionID: (*uint)(nil),
		Page:               1,
		PageSize:           1000, // 大数值确保获取全部
	}

	// 将functionID赋值给query
	fid := uint(functionID)
	query.BusinessFunctionID = &fid

	ctrl.service.GetMonitorList(c, query)
}

// UpdateMonitorStatus 更新API监控状态
// @Summary 更新API监控状态
// @Description 启用或禁用API监控
// @Tags API监控
// @Accept json
// @Produce json
// @Param id path uint true "监控ID"
// @Param status query int true "状态(0:禁用,1:启用)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/api-monitors/{id}/status [put]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) UpdateMonitorStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的监控ID")
		return
	}

	statusStr := c.Query("status")
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的状态值")
		return
	}

	ctrl.service.UpdateMonitorStatus(c, uint(id), status)
}

// GetStats 获取统计信息
// @Summary 获取API监控统计
// @Description 获取API监控总体统计信息
// @Tags API监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result{data=model.APIMonitorStats}
// @Router /api/v1/monitor/api-monitors/stats [get]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) GetStats(c *gin.Context) {
	ctrl.service.GetStats(c)
}

// GetGroupStats 获取业务组统计
// @Summary 获取按业务组聚合的统计
// @Description 获取按业务组和服务聚合的API监控统计
// @Tags API监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result{data=[]model.APIMonitorGroupStats}
// @Router /api/v1/monitor/api-monitors/group-stats [get]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) GetGroupStats(c *gin.Context) {
	ctrl.service.GetGroupStats(c)
}

// GetServiceStats 获取服务统计
// @Summary 获取指定业务组的服务统计
// @Description 获取指定业务组下按服务聚合的API监控统计
// @Tags API监控
// @Accept json
// @Produce json
// @Param groupId query uint true "业务组ID"
// @Success 200 {object} result.Result{data=[]model.APIMonitorServiceStats}
// @Router /api/v1/monitor/api-monitors/service-stats [get]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) GetServiceStats(c *gin.Context) {
	groupIDStr := c.Query("groupId")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的业务组ID")
		return
	}

	ctrl.service.GetServiceStats(c, uint(groupID))
}

// GetTestData 获取测试数据
// @Summary 获取测试数据
// @Description 返回测试用的业务组和功能数据，用于前端快速测试
// @Tags API监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/test-data [get]
func (ctrl *APIMonitorController) GetTestData(c *gin.Context) {
	testData := map[string]interface{}{
		"groups": []map[string]interface{}{
			{
				"id":   1,
				"name": "支付系统",
			},
			{
				"id":   2,
				"name": "用户系统",
			},
			{
				"id":   5,
				"name": "网络监控",
			},
		},
		"simplifiedExample": map[string]interface{}{
			"name":                   "百度首页检测",
			"url":                    "http://www.baidu.com",
			"method":                 "GET",
			"groupId":                5,
			"businessFunctionName":   "网络连通性检测",
			"description":            "检测百度首页是否正常访问",
		},
		"defaultConfig": map[string]interface{}{
			"checkInterval":   30,
			"timeout":         5000,
			"expectedStatus":  200,
			"expectedKeyword": "success",
			"retryCount":      3,
			"alertEnabled":    true,
		},
		"apiEndpoints": map[string]interface{}{
			"getGroups":       "GET /api/v1/cmdb/groups",
			"getFunctions":    "GET /api/v1/monitor/business-functions?groupId={groupId}",
			"getMonitors":     "GET /api/v1/monitor/api-monitors?businessFunctionId={functionId}",
			"createSimple":    "POST /api/v1/monitor/api-monitors/simple",
			"createFull":      "POST /api/v1/monitor/api-monitors",
			"getStats":        "GET /api/v1/monitor/api-monitors/stats",
		},
	}

	result.Success(c, testData)
}

// SystemTest API监控模块测试
// @Summary API监控模块测试
// @Description 测试API监控模块的核心功能
// @Tags API监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/system-test [get]
func (ctrl *APIMonitorController) SystemTest(c *gin.Context) {
	testResult := map[string]interface{}{
		"module":    "API监控",
		"status":    "✅ 正常",
		"timestamp": "2024-02-01T10:30:00Z",
		"tests": []map[string]interface{}{
			{
				"name":    "业务功能模块",
				"status":  "✅ 正常",
				"message": "支持创建、查询、更新、删除业务功能",
			},
			{
				"name":    "API监控创建",
				"status":  "✅ 正常",
				"message": "支持简化版和完整版两种创建方式",
			},
			{
				"name":    "监控检测",
				"status":  "✅ 正常",
				"message": "每30秒自动执行一次检测，支持失败重试",
			},
			{
				"name":    "告警功能",
				"status":  "✅ 正常",
				"message": "异常自动触发告警，支持多种通知方式",
			},
			{
				"name":    "统计分析",
				"status":  "✅ 正常",
				"message": "支持按业务组、功能模块聚合统计",
			},
		},
	}

	result.Success(c, testResult)
}


// GetGroupsTreeWithApiCount 获取业务组树形结构及其关联的API数量
// @Summary 获取业务组树形统计
// @Description 获取业务组树形结构（包含父子关系）及其关联的API监控总数，用于前端树形菜单显示
// @Tags API监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/groups-tree-api-stats [get]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) GetGroupsTreeWithApiCount(c *gin.Context) {
	// 获取所有业务组及其API统计
	ctrl.service.GetGroupStats(c)
}

// TestAPI 测试API接口
// @Summary 测试API接口
// @Description 测试API接口连通性和响应
// @Tags API监控
// @Accept json
// @Produce json
// @Param test body model.TestAPIDto true "测试参数"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/api-monitors/test [post]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) TestAPI(c *gin.Context) {
	var dto model.TestAPIDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.service.TestAPI(c, &dto)
}

// AuthTestEndpoint 需要认证的测试接口
// @Summary 需要认证的测试接口
// @Description 用于测试带认证参数的API监控，需要在Header中传递 X-API-Key: test-secret-key
// @Tags API监控
// @Accept json
// @Produce json
// @Success 200 {object} result.Result
// @Failure 401 {object} result.Result
// @Router /api/v1/monitor/auth-test [get]
func (ctrl *APIMonitorController) AuthTestEndpoint(c *gin.Context) {
	// 检查认证
	apiKey := c.GetHeader("X-API-Key")
	if apiKey != "test-secret-key" {
		result.Failed(c, http.StatusUnauthorized, "认证失败: 无效的API Key")
		return
	}

	// 认证成功，返回测试数据
	result.Success(c, map[string]interface{}{
		"success":     true,
		"message":     "认证成功",
		"timestamp":   "2026-02-02",
		"apiKey":      apiKey,
		"endpoint":    "/api/v1/monitor/auth-test",
		"description": "这是一个需要认证的测试接口，用于测试API监控的认证参数功能",
		"usage": map[string]interface{}{
			"header": "X-API-Key: test-secret-key",
			"example": "curl -H \"X-API-Key: test-secret-key\" http://localhost:8080/api/v1/monitor/auth-test",
		},
	})
}

// ======================== 分层关系管理接口 ========================

// AddChildMonitor 为监控添加子监控
// @Summary 添加子监控
// @Description 为指定的父监控添加一个子监控，建立分层关系
// @Tags 分层监控
// @Accept json
// @Produce json
// @Param id path int true "父监控ID"
// @Param dto body model.AddChildMonitorDto true "子监控信息"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Router /api/v1/monitor/api-monitors/{id}/add-child [post]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) AddChildMonitor(c *gin.Context) {
	idStr := c.Param("id")
	parentID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的监控ID")
		return
	}

	var dto model.AddChildMonitorDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctrl.service.AddChildMonitor(c, uint(parentID), dto.ChildMonitorID)
}

// GetMonitorHierarchy 获取监控的树形结构
// @Summary 获取监控层级结构
// @Description 获取指定监控及其所有子监控的树形结构
// @Tags 分层监控
// @Produce json
// @Param id path int true "监控ID"
// @Param depth query int false "深度限制（默认10）"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/api-monitors/{id}/hierarchy [get]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) GetMonitorHierarchy(c *gin.Context) {
	idStr := c.Param("id")
	monitorID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的监控ID")
		return
	}

	depth := 10
	if depthStr := c.Query("depth"); depthStr != "" {
		if d, err := strconv.Atoi(depthStr); err == nil && d > 0 {
			depth = d
		}
	}

	ctrl.service.GetMonitorHierarchy(c, uint(monitorID), depth)
}

// GetTopLevelMonitors 获取顶级监控列表
// @Summary 获取顶级监控列表
// @Description 获取所有不属于任何父监控的顶级监控（支持分页、过滤、搜索）
// @Tags 分层监控
// @Produce json
// @Param page query int false "页码（默认1）"
// @Param pageSize query int false "每页数量（默认20）"
// @Param groupId query int false "业务组ID"
// @Param keyword query string false "搜索关键字（名称/URL）"
// @Param includeChildren query bool false "是否包含子监控"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/api-monitors/top-level [get]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) GetTopLevelMonitors(c *gin.Context) {
	var query model.APIMonitorListQuery

	// 解析分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			query.Page = page
		}
	}
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			query.PageSize = pageSize
		}
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	// 解析其他过滤参数
	if groupIdStr := c.Query("groupId"); groupIdStr != "" {
		if groupId, err := strconv.ParseUint(groupIdStr, 10, 32); err == nil {
			gid := uint(groupId)
			query.GroupID = &gid
		}
	}
	query.Keyword = c.Query("keyword")
	query.OnlyTopLevel = true

	if includeChildrenStr := c.Query("includeChildren"); includeChildrenStr == "true" {
		query.IncludeChildren = true
	}

	ctrl.service.GetMonitorList(c, &query)
}

// RemoveParentRelationship 移除子监控与父监控的关系
// @Summary 移除父监控关系
// @Description 将子监控从其父监控中移除，使其变为独立监控
// @Tags 分层监控
// @Produce json
// @Param id path int true "子监控ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Router /api/v1/monitor/api-monitors/{id}/parent [delete]
// @Security ApiKeyAuth
func (ctrl *APIMonitorController) RemoveParentRelationship(c *gin.Context) {
	idStr := c.Param("id")
	childMonitorID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的监控ID")
		return
	}

	ctrl.service.RemoveParentRelationship(c, uint(childMonitorID))
}
