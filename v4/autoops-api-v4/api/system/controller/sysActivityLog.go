// 系统动态日志 控制层
package controller

import (
	"strconv"

	"dodevops-api/api/system/model"
	"dodevops-api/api/system/service"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

type SysActivityLogController struct {
	service *service.SysActivityLogService
}

func NewSysActivityLogController() *SysActivityLogController {
	return &SysActivityLogController{
		service: service.NewSysActivityLogService(),
	}
}

// Create 创建系统动态日志
// @Summary 创建系统动态日志
// @Description 创建系统动态日志
// @Tags System系统管理
// @Param log body model.CreateActivityLogDto true "动态日志信息"
// @Success 200 {object} result.Result
// @Router /api/v1/system/activity-log [post]
// @Security ApiKeyAuth
func (c *SysActivityLogController) Create(ctx *gin.Context) {
	var dto model.CreateActivityLogDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	if err := c.service.CreateWithDto(&dto); err != nil {
		result.Failed(ctx, 500, "创建失败: "+err.Error())
		return
	}

	result.Success(ctx, nil)
}

// GetByID 根据ID查询
// @Summary 根据ID查询动态日志
// @Description 根据ID查询动态日志
// @Tags System系统管理
// @Param id query uint true "日志ID"
// @Success 200 {object} result.Result{data=model.SysActivityLog}
// @Router /api/v1/system/activity-log/detail [get]
// @Security ApiKeyAuth
func (c *SysActivityLogController) GetByID(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	log, err := c.service.GetByID(req.ID)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, log)
}

// List 分页查询动态日志列表
// @Summary 分页查询动态日志列表
// @Description 分页查询动态日志列表
// @Tags System系统管理
// @Param activityType query int false "动态类型：1=密钥同步，2=域名检查，3=服务器巡检，4=定时任务，5=其他"
// @Param status query int false "状态"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} result.Result{data=result.PageResult{list=[]model.SysActivityLog}}
// @Router /api/v1/system/activity-log/list [get]
// @Security ApiKeyAuth
func (c *SysActivityLogController) List(ctx *gin.Context) {
	activityType, _ := strconv.Atoi(ctx.Query("activityType"))
	status, _ := strconv.Atoi(ctx.Query("status"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	logs, total, err := c.service.List(activityType, status, page, pageSize)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	pageResult := result.PageResult{
		List:     logs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	result.Success(ctx, pageResult)
}

// GetLatest 获取最新动态
// @Summary 获取最新N条动态日志
// @Description 获取最新N条动态日志
// @Tags System系统管理
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} result.Result{data=[]model.SysActivityLog}
// @Router /api/v1/system/activity-log/latest [get]
// @Security ApiKeyAuth
func (c *SysActivityLogController) GetLatest(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	logs, err := c.service.GetLatest(limit)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, logs)
}

// GetByRelatedID 根据关联ID查询
// @Summary 根据关联ID查询动态日志
// @Description 根据关联ID查询动态日志
// @Tags System系统管理
// @Param activityType query int false "动态类型：1=密钥同步，2=域名检查，3=服务器巡检，4=定时任务，5=其他"
// @Param relatedId query uint true "关联ID"
// @Success 200 {object} result.Result{data=[]model.SysActivityLog}
// @Router /api/v1/system/activity-log/related [get]
// @Security ApiKeyAuth
func (c *SysActivityLogController) GetByRelatedID(ctx *gin.Context) {
	var req struct {
		ActivityType int  `form:"activityType"`
		RelatedID    uint `form:"relatedId" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	logs, err := c.service.GetByRelatedID(req.ActivityType, req.RelatedID)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, logs)
}

// DeleteByID 根据ID删除
// @Summary 根据ID删除动态日志
// @Description 根据ID删除动态日志
// @Tags System系统管理
// @Param id query uint true "日志ID"
// @Success 200 {object} result.Result
// @Router /api/v1/system/activity-log [delete]
// @Security ApiKeyAuth
func (c *SysActivityLogController) DeleteByID(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	if err := c.service.DeleteByID(req.ID); err != nil {
		result.Failed(ctx, 500, "删除失败: "+err.Error())
		return
	}

	result.Success(ctx, nil)
}

// DeleteBatch 批量删除
// @Summary 批量删除动态日志
// @Description 批量删除动态日志
// @Tags System系统管理
// @Param ids body []uint true "日志ID列表"
// @Success 200 {object} result.Result
// @Router /api/v1/system/activity-log/batch [delete]
// @Security ApiKeyAuth
func (c *SysActivityLogController) DeleteBatch(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	if err := c.service.DeleteBatch(req.IDs); err != nil {
		result.Failed(ctx, 500, "删除失败: "+err.Error())
		return
	}

	result.Success(ctx, nil)
}

// Clean 清空动态日志
// @Summary 清空动态日志（保留最近N天）
// @Description 清空动态日志（保留最近N天）
// @Tags System系统管理
// @Param days query int false "保留天数" default(30)
// @Success 200 {object} result.Result
// @Router /api/v1/system/activity-log/clean [delete]
// @Security ApiKeyAuth
func (c *SysActivityLogController) Clean(ctx *gin.Context) {
	days, _ := strconv.Atoi(ctx.DefaultQuery("days", "30"))

	if err := c.service.Clean(days); err != nil {
		result.Failed(ctx, 500, "清空失败: "+err.Error())
		return
	}

	result.Success(ctx, nil)
}
