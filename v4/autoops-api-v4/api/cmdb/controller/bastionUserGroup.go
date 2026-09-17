package controller

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util"

	"github.com/gin-gonic/gin"
)

type BastionUserGroupController struct {
	service service.BastionUserGroupServiceInterface
}

func NewBastionUserGroupController() *BastionUserGroupController {
	return &BastionUserGroupController{
		service: service.GetBastionUserGroupService(),
	}
}

// 创建用户组
// @Summary 创建用户组
// @Description 创建堡垒机用户组
// @Tags 堡垒机用户组管理
// @Accept json
// @Produce json
// @Param data body model.CreateBastionUserGroupDto true "用户组信息"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/group/create [post]
// @Security ApiKeyAuth
func (c *BastionUserGroupController) CreateUserGroup(ctx *gin.Context) {
	var dto model.CreateBastionUserGroupDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.CreateUserGroup(ctx, &dto)
}

// 更新用户组
// @Summary 更新用户组
// @Description 更新堡垒机用户组
// @Tags 堡垒机用户组管理
// @Accept json
// @Produce json
// @Param data body model.UpdateBastionUserGroupDto true "用户组信息"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/group/update [put]
// @Security ApiKeyAuth
func (c *BastionUserGroupController) UpdateUserGroup(ctx *gin.Context) {
	var dto model.UpdateBastionUserGroupDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.UpdateUserGroup(ctx, &dto)
}

// 删除用户组
// @Summary 删除用户组
// @Description 删除堡垒机用户组
// @Tags 堡垒机用户组管理
// @Accept json
// @Produce json
// @Param data body model.BastionUserGroupIdDto true "用户组ID"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/group/delete [delete]
// @Security ApiKeyAuth
func (c *BastionUserGroupController) DeleteUserGroup(ctx *gin.Context) {
	var dto model.BastionUserGroupIdDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.DeleteUserGroup(ctx, dto.ID)
}

// 根据ID获取用户组
// @Summary 根据ID获取用户组
// @Description 获取堡垒机用户组详情
// @Tags 堡垒机用户组管理
// @Accept json
// @Produce json
// @Param id query uint true "用户组ID"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/group/info [get]
// @Security ApiKeyAuth
func (c *BastionUserGroupController) GetUserGroupById(ctx *gin.Context) {
	id := util.StringToUint(ctx.Query("id"))
	if id == 0 {
		result.Failed(ctx, constant.INVALID_PARAMS, "用户组ID不能为空")
		return
	}
	c.service.GetUserGroupById(ctx, id)
}

// 获取用户组列表(分页)
// @Summary 获取用户组列表(分页)
// @Description 获取堡垒机用户组列表(分页)
// @Tags 堡垒机用户组管理
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/group/list [get]
// @Security ApiKeyAuth
func (c *BastionUserGroupController) GetUserGroupList(ctx *gin.Context) {
	page := util.StringToInt(ctx.Query("page"))
	pageSize := util.StringToInt(ctx.Query("pageSize"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	c.service.GetUserGroupList(ctx, page, pageSize)
}

// 添加用户组成员
// @Summary 添加用户组成员
// @Description 批量添加用户到用户组
// @Tags 堡垒机用户组管理
// @Accept json
// @Produce json
// @Param data body model.AddGroupMembersDto true "添加成员参数"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/group/member/add [post]
// @Security ApiKeyAuth
func (c *BastionUserGroupController) AddGroupMembers(ctx *gin.Context) {
	var dto model.AddGroupMembersDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.AddGroupMembers(ctx, &dto)
}

// 删除用户组成员
// @Summary 删除用户组成员
// @Description 批量从用户组移除用户
// @Tags 堡垒机用户组管理
// @Accept json
// @Produce json
// @Param data body model.RemoveGroupMembersDto true "删除成员参数"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/group/member/remove [post]
// @Security ApiKeyAuth
func (c *BastionUserGroupController) RemoveGroupMembers(ctx *gin.Context) {
	var dto model.RemoveGroupMembersDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	c.service.RemoveGroupMembers(ctx, &dto)
}

// 获取用户组成员列表
// @Summary 获取用户组成员列表
// @Description 获取指定用户组的所有成员
// @Tags 堡垒机用户组管理
// @Accept json
// @Produce json
// @Param groupId query uint true "用户组ID"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/group/member/list [get]
// @Security ApiKeyAuth
func (c *BastionUserGroupController) GetGroupMembers(ctx *gin.Context) {
	groupId := util.StringToUint(ctx.Query("groupId"))
	if groupId == 0 {
		result.Failed(ctx, constant.INVALID_PARAMS, "用户组ID不能为空")
		return
	}
	c.service.GetGroupMembers(ctx, groupId)
}

// 获取用户所属的用户组
// @Summary 获取用户所属的用户组
// @Description 获取指定用户所属的所有用户组
// @Tags 堡垒机用户组管理
// @Accept json
// @Produce json
// @Param userId query uint true "用户ID"
// @Success 200 {object} result.Result
// @Router /api/v1/bastion/user/groups [get]
// @Security ApiKeyAuth
func (c *BastionUserGroupController) GetUserGroups(ctx *gin.Context) {
	userId := util.StringToUint(ctx.Query("userId"))
	if userId == 0 {
		result.Failed(ctx, constant.INVALID_PARAMS, "用户ID不能为空")
		return
	}
	c.service.GetUserGroups(ctx, userId)
}
