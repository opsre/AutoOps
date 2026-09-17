package controller

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util"

	"github.com/gin-gonic/gin"
)

type BastionHostAuthorizationController struct {
	service service.BastionHostAuthorizationServiceInterface
}

func NewBastionHostAuthorizationController() *BastionHostAuthorizationController {
	return &BastionHostAuthorizationController{
		service: service.GetBastionHostAuthorizationService(),
	}
}

// 授权主机给用户
// @Summary 授权主机给用户
// @Description 批量授权主机给指定用户
// @Tags 堡垒机授权管理
// @Accept json
// @Produce json
// @Param data body model.AuthorizeHostsToUserDto true "授权参数"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/auth/user/authorize [post]
// @Security ApiKeyAuth
func (c *BastionHostAuthorizationController) AuthorizeHostsToUser(ctx *gin.Context) {
	var dto model.AuthorizeHostsToUserDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.AuthorizeHostsToUser(ctx, &dto)
}

// 授权主机给用户组
// @Summary 授权主机给用户组
// @Description 批量授权主机给指定用户组
// @Tags 堡垒机授权管理
// @Accept json
// @Produce json
// @Param data body model.AuthorizeHostsToGroupDto true "授权参数"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/auth/group/authorize [post]
// @Security ApiKeyAuth
func (c *BastionHostAuthorizationController) AuthorizeHostsToGroup(ctx *gin.Context) {
	var dto model.AuthorizeHostsToGroupDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.AuthorizeHostsToGroup(ctx, &dto)
}

// 取消用户的主机授权
// @Summary 取消用户的主机授权
// @Description 批量取消用户的主机授权
// @Tags 堡垒机授权管理
// @Accept json
// @Produce json
// @Param data body model.RevokeHostsFromUserDto true "取消授权参数"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/auth/user/revoke [post]
// @Security ApiKeyAuth
func (c *BastionHostAuthorizationController) RevokeHostsFromUser(ctx *gin.Context) {
	var dto model.RevokeHostsFromUserDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.RevokeHostsFromUser(ctx, &dto)
}

// 取消用户组的主机授权
// @Summary 取消用户组的主机授权
// @Description 批量取消用户组的主机授权
// @Tags 堡垒机授权管理
// @Accept json
// @Produce json
// @Param data body model.RevokeHostsFromGroupDto true "取消授权参数"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/auth/group/revoke [post]
// @Security ApiKeyAuth
func (c *BastionHostAuthorizationController) RevokeHostsFromGroup(ctx *gin.Context) {
	var dto model.RevokeHostsFromGroupDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.RevokeHostsFromGroup(ctx, &dto)
}

// 获取授权主机列表
// @Summary 获取授权主机列表
// @Description 获取用户/用户组已授权的主机列表(分页)
// @Tags 堡垒机授权管理
// @Accept json
// @Produce json
// @Param subjectType query string true "授权对象类型(user/group)"
// @Param subjectId query uint true "授权对象ID"
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/auth/hosts [get]
// @Security ApiKeyAuth
func (c *BastionHostAuthorizationController) GetAuthorizedHosts(ctx *gin.Context) {
	var dto model.GetAuthorizedHostsDto
	if err := ctx.ShouldBindQuery(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.GetAuthorizedHosts(ctx, &dto)
}

// 获取用户有权限的主机列表
// @Summary 获取用户有权限的主机列表
// @Description 获取用户有权限的主机列表(包括直接授权和通过用户组授权，自动去重)
// @Tags 堡垒机授权管理
// @Accept json
// @Produce json
// @Param userId query uint true "用户ID"
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/user/authorized-hosts [get]
// @Security ApiKeyAuth
func (c *BastionHostAuthorizationController) GetUserAuthorizedHosts(ctx *gin.Context) {
	userId := util.StringToUint(ctx.Query("userId"))
	page := util.StringToInt(ctx.Query("page"))
	pageSize := util.StringToInt(ctx.Query("pageSize"))

	if userId == 0 {
		result.Failed(ctx, constant.INVALID_PARAMS, "用户ID不能为空")
		return
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	c.service.GetUserAuthorizedHosts(ctx, userId, page, pageSize)
}

// 检查用户是否有主机权限
// @Summary 检查用户是否有主机权限
// @Description 检查指定用户是否有指定主机的访问权限
// @Tags 堡垒机授权管理
// @Accept json
// @Produce json
// @Param userId query uint true "用户ID"
// @Param hostId query uint true "主机ID"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/user/check-permission [get]
// @Security ApiKeyAuth
func (c *BastionHostAuthorizationController) CheckUserHostPermission(ctx *gin.Context) {
	userId := util.StringToUint(ctx.Query("userId"))
	hostId := util.StringToUint(ctx.Query("hostId"))

	if userId == 0 || hostId == 0 {
		result.Failed(ctx, constant.INVALID_PARAMS, "用户ID和主机ID不能为空")
		return
	}

	c.service.CheckUserHostPermission(ctx, userId, hostId)
}
