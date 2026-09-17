package controller

import (
	"strconv"

	"dodevops-api/api/app/model"
	"dodevops-api/api/app/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ServiceReleaseController 服务上线控制器
type ServiceReleaseController struct {
	service service.IServiceReleaseService
}

// NewServiceReleaseController 创建服务上线控制器
func NewServiceReleaseController(db *gorm.DB) *ServiceReleaseController {
	return &ServiceReleaseController{
		service: service.NewServiceReleaseService(db),
	}
}

// CreateServiceRelease 创建服务上线工单
// @Summary 创建服务上线工单
// @Description 创建服务上线工单，支持多个服务的上线发布流程
// @Tags ServiceRelease
// @Accept json
// @Produce json
// @Param request body model.CreateServiceReleaseRequest true "创建服务上线请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.ServiceRelease}
// @Router /apps/service-release [post]
func (ctrl *ServiceReleaseController) CreateServiceRelease(c *gin.Context) {
	var req model.CreateServiceReleaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.CreateServiceRelease(c, &req)
}

// GetServiceReleaseList 获取服务上线工单列表
// @Summary 获取服务上线工单列表
// @Description 获取服务上线工单列表，支持分页和筛选
// @Tags ServiceRelease
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param business_group_id query int false "业务组ID"
// @Param status query int false "流程状态"
// @Param applicant_id query int false "申请人ID"
// @Param keyword query string false "搜索关键词（标题）"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.ServiceReleaseListResponse}
// @Router /apps/service-release/list [get]
func (ctrl *ServiceReleaseController) GetServiceReleaseList(c *gin.Context) {
	var req model.ServiceReleaseListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.GetServiceReleaseList(c, &req)
}

// GetServiceReleaseDetail 获取服务上线工单详情
// @Summary 获取服务上线工单详情
// @Description 获取服务上线工单详情，包含所有服务明细信息
// @Tags ServiceRelease
// @Accept json
// @Produce json
// @Param id path int true "工单ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.ServiceRelease}
// @Router /apps/service-release/{id} [get]
func (ctrl *ServiceReleaseController) GetServiceReleaseDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的工单ID")
		return
	}

	ctrl.service.GetServiceReleaseDetail(c, uint(id))
}

// DeleteServiceRelease 删除服务上线工单
// @Summary 删除服务上线工单
// @Description 删除服务上线工单（级联删除所有服务明细）
// @Tags ServiceRelease
// @Accept json
// @Produce json
// @Param id path int true "工单ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=string}
// @Router /apps/service-release/{id} [delete]
func (ctrl *ServiceReleaseController) DeleteServiceRelease(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的工单ID")
		return
	}

	ctrl.service.DeleteServiceRelease(c, uint(id))
}

// ApproveServiceRelease 审批服务上线工单
// @Summary 审批服务上线工单
// @Description 负责人、安全、测试三方审批服务上线工单
// @Tags ServiceRelease
// @Accept json
// @Produce json
// @Param request body model.ApprovalRequest true "审批请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=string}
// @Router /apps/service-release/approve [post]
func (ctrl *ServiceReleaseController) ApproveServiceRelease(c *gin.Context) {
	var req model.ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.ApproveRelease(c, &req)
}

// DeployServiceRelease 运维发布服务上线
// @Summary 运维发布服务上线
// @Description 三方审批通过后，运维执行发布操作
// @Tags ServiceRelease
// @Accept json
// @Produce json
// @Param request body model.DeployReleaseRequest true "发布请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=string}
// @Router /apps/service-release/deploy [post]
func (ctrl *ServiceReleaseController) DeployServiceRelease(c *gin.Context) {
	var req model.DeployReleaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.DeployRelease(c, &req)
}

// RegressionTest 回归测试
// @Summary 回归测试
// @Description 发布成功后进行回归测试验证
// @Tags ServiceRelease
// @Accept json
// @Produce json
// @Param request body model.RegressionTestRequest true "回归测试请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=string}
// @Router /apps/service-release/regression-test [post]
func (ctrl *ServiceReleaseController) RegressionTest(c *gin.Context) {
	var req model.RegressionTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.RegressionTest(c, &req)
}

// GetServiceItemLog 获取单个服务的Jenkins日志
// @Summary 获取单个服务的Jenkins日志
// @Description 获取服务上线工单中单个服务的Jenkins构建日志
// @Tags ServiceRelease
// @Accept json
// @Produce json
// @Param item_id path int true "服务明细ID"
// @Param start query int false "日志起始行数" default(0)
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=object}
// @Router /apps/service-release/services/:item_id/log [get]
func (ctrl *ServiceReleaseController) GetServiceItemLog(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的服务ID")
		return
	}

	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))

	ctrl.service.GetServiceItemLog(c, uint(itemID), start)
}

// CheckUserPermission 检查用户操作权限
// @Summary 检查用户操作权限
// @Description 检查当前用户是否有权限执行指定的操作
// @Tags ServiceRelease
// @Accept json
// @Produce json
// @Param id path int true "工单ID"
// @Param operation_type query string true "操作类型" Enums(owner_approval, security_approval, test_approval, deploy, regression_test)
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=object}
// @Router /apps/service-release/:id/permission [get]
func (ctrl *ServiceReleaseController) CheckUserPermission(c *gin.Context) {
	releaseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的工单ID")
		return
	}

	operationType := c.Query("operation_type")
	if operationType == "" {
		result.Failed(c, 400, "请指定操作类型")
		return
	}

	ctrl.service.CheckUserPermission(c, uint(releaseID), operationType)
}
