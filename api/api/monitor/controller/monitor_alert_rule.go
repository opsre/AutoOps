package controller

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MonitorAlertRuleController struct {
	ruleService service.MonitorAlertRuleService
}

func NewMonitorAlertRuleController() *MonitorAlertRuleController {
	return &MonitorAlertRuleController{
		ruleService: service.NewMonitorAlertRuleService(),
	}
}

// ==== Group Rule API ====

func (c *MonitorAlertRuleController) CreateGroup(ctx *gin.Context) {
	var group model.MonitorAlertGroupRule
	if err := ctx.ShouldBindJSON(&group); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "参数绑定失败: "+err.Error())
		return
	}
	if err := c.ruleService.CreateGroup(&group); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}
	result.Success(ctx, group)
}

func (c *MonitorAlertRuleController) UpdateGroup(ctx *gin.Context) {
	var group model.MonitorAlertGroupRule
	if err := ctx.ShouldBindJSON(&group); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "参数绑定失败: "+err.Error())
		return
	}
	if group.ID == 0 {
		result.Failed(ctx, int(result.ApiCode.FAILED), "群组ID不能为空")
		return
	}
	if err := c.ruleService.UpdateGroup(&group); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}
	result.Success(ctx, group)
}

func (c *MonitorAlertRuleController) DeleteGroup(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "无效的ID")
		return
	}
	if err := c.ruleService.DeleteGroup(uint(id)); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}
	result.Success(ctx, map[string]interface{}{"id": id})
}

func (c *MonitorAlertRuleController) GetGroupByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "无效的ID")
		return
	}
	data, err := c.ruleService.GetGroupByID(uint(id))
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "获取群组失败")
		return
	}
	result.Success(ctx, data)
}

func (c *MonitorAlertRuleController) GetGroupList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	list, total, err := c.ruleService.GetGroupList(page, pageSize)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "获取群组列表失败")
		return
	}
	result.SuccessWithPage(ctx, list, total, page, pageSize)
}

// ==== Monitor Alert Rule (Children) API ====

func (c *MonitorAlertRuleController) CreateRule(ctx *gin.Context) {
	var rule model.MonitorAlertRule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "参数绑定失败: "+err.Error())
		return
	}
	if rule.GroupID == 0 {
		result.Failed(ctx, int(result.ApiCode.FAILED), "请选择归属群组")
		return
	}
	if err := c.ruleService.CreateRule(&rule); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}
	result.Success(ctx, rule)
}

func (c *MonitorAlertRuleController) UpdateRule(ctx *gin.Context) {
	var rule model.MonitorAlertRule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "参数绑定失败: "+err.Error())
		return
	}
	if rule.ID == 0 {
		result.Failed(ctx, int(result.ApiCode.FAILED), "规则ID不能为空")
		return
	}
	if err := c.ruleService.UpdateRule(&rule); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}
	result.Success(ctx, rule)
}

func (c *MonitorAlertRuleController) DeleteRule(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "无效的ID")
		return
	}
	if err := c.ruleService.DeleteRule(uint(id)); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}
	result.Success(ctx, map[string]interface{}{"id": id})
}

func (c *MonitorAlertRuleController) GetRuleByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "无效的ID")
		return
	}
	data, err := c.ruleService.GetRuleByID(uint(id))
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "获取规则失败")
		return
	}
	result.Success(ctx, data)
}

func (c *MonitorAlertRuleController) GetRuleListByGroup(ctx *gin.Context) {
	groupIdStr := ctx.Param("groupId")
	groupId, _ := strconv.Atoi(groupIdStr)

	var req model.MonitorAlertRuleQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "查询参数错误")
		return
	}

	// 强制使用PathParam指定的GroupID
	req.GroupID = uint(groupId)

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	list, total, err := c.ruleService.GetRuleList(&req)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "获取规则列表失败")
		return
	}
	result.SuccessWithPage(ctx, list, total, req.Page, req.PageSize)
}

func (c *MonitorAlertRuleController) GetRuleList(ctx *gin.Context) {
	var req model.MonitorAlertRuleQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "查询参数错误")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	list, total, err := c.ruleService.GetRuleList(&req)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), "获取规则列表失败")
		return
	}
	result.SuccessWithPage(ctx, list, total, req.Page, req.PageSize)
}

// CheckRuleExpr 表达式健康检查
func (c *MonitorAlertRuleController) CheckRuleExpr(ctx *gin.Context) {
	var req model.MonitorAlertRuleCheckReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "无效的请求参数")
		return
	}

	data, err := c.ruleService.CheckRuleExpr(req.DataSourceID, req.Expr)
	if err != nil {
		result.Failed(ctx, 500, err.Error())
		return
	}

	result.Success(ctx, data)
}
