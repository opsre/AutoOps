package service

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	systemDao "dodevops-api/api/system/dao"
	systemModel "dodevops-api/api/system/model"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type BastionUserGroupServiceInterface interface {
	CreateUserGroup(c *gin.Context, dto *model.CreateBastionUserGroupDto)      // 创建用户组
	UpdateUserGroup(c *gin.Context, dto *model.UpdateBastionUserGroupDto)      // 更新用户组
	DeleteUserGroup(c *gin.Context, id uint)                                   // 删除用户组
	GetUserGroupById(c *gin.Context, id uint)                                  // 根据ID获取用户组
	GetUserGroupList(c *gin.Context, page, pageSize int)                       // 获取用户组列表(分页)
	AddGroupMembers(c *gin.Context, dto *model.AddGroupMembersDto)             // 添加用户组成员
	RemoveGroupMembers(c *gin.Context, dto *model.RemoveGroupMembersDto)       // 删除用户组成员
	GetGroupMembers(c *gin.Context, groupId uint)                              // 获取用户组成员列表
	GetUserGroups(c *gin.Context, userId uint)                                 // 获取用户所属的用户组
}

type BastionUserGroupServiceImpl struct {
	dao     dao.BastionUserGroupDao
	authDao dao.BastionHostAuthorizationDao
}

// 创建用户组
func (s *BastionUserGroupServiceImpl) CreateUserGroup(c *gin.Context, dto *model.CreateBastionUserGroupDto) {
	// 检查用户组名称是否已存在
	if s.dao.CheckNameExists(dto.Name) {
		result.Failed(c, constant.INVALID_PARAMS, "用户组名称已存在")
		return
	}

	// 验证成员ID是否存在
	if len(dto.MemberIds) > 0 {
		for _, userId := range dto.MemberIds {
			user := systemDao.GetSysAdminInfo(int(userId))
			if user.ID == 0 {
				result.Failed(c, constant.INVALID_PARAMS, fmt.Sprintf("用户ID %d 不存在", userId))
				return
			}
		}
	}

	// 创建用户组
	group := model.BastionUserGroup{
		Name:        dto.Name,
		Description: dto.Description,
		MemberIds:   model.UintArray(dto.MemberIds),
		Status:      dto.Status,
		CreateTime:  util.HTime{Time: time.Now()},
		UpdateTime:  util.HTime{Time: time.Now()},
	}

	if err := s.dao.CreateUserGroup(&group); err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "创建用户组失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{
		"id":      group.ID,
		"message": "用户组创建成功",
	})
}

// 更新用户组
func (s *BastionUserGroupServiceImpl) UpdateUserGroup(c *gin.Context, dto *model.UpdateBastionUserGroupDto) {
	// 检查用户组是否存在
	existGroup, err := s.dao.GetUserGroupById(dto.ID)
	if err != nil {
		result.Failed(c, constant.CMDB_HOST_NOT_FOUND, "用户组不存在")
		return
	}

	// 检查用户组名称是否重复(排除自身)
	if dto.Name != existGroup.Name && s.dao.CheckNameExistsExcludeId(dto.Name, dto.ID) {
		result.Failed(c, constant.INVALID_PARAMS, "用户组名称已存在")
		return
	}

	// 验证成员ID是否存在
	if len(dto.MemberIds) > 0 {
		for _, userId := range dto.MemberIds {
			user := systemDao.GetSysAdminInfo(int(userId))
			if user.ID == 0 {
				result.Failed(c, constant.INVALID_PARAMS, fmt.Sprintf("用户ID %d 不存在", userId))
				return
			}
		}
	}

	// 更新用户组
	group := model.BastionUserGroup{
		Name:        dto.Name,
		Description: dto.Description,
		MemberIds:   model.UintArray(dto.MemberIds),
		Status:      dto.Status,
		UpdateTime:  util.HTime{Time: time.Now()},
	}

	if err := s.dao.UpdateUserGroup(dto.ID, &group); err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "更新用户组失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{"message": "用户组更新成功"})
}

// 删除用户组
func (s *BastionUserGroupServiceImpl) DeleteUserGroup(c *gin.Context, id uint) {
	// 检查用户组是否存在
	_, err := s.dao.GetUserGroupById(id)
	if err != nil {
		result.Failed(c, constant.CMDB_HOST_NOT_FOUND, "用户组不存在")
		return
	}

	// 删除用户组的所有授权
	if err := s.authDao.RevokeAllSubjectAuthorization("group", id); err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "删除用户组授权失败: "+err.Error())
		return
	}

	// 删除用户组
	if err := s.dao.DeleteUserGroup(id); err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "删除用户组失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{"message": "用户组删除成功"})
}

// 根据ID获取用户组
func (s *BastionUserGroupServiceImpl) GetUserGroupById(c *gin.Context, id uint) {
	group, err := s.dao.GetUserGroupById(id)
	if err != nil {
		result.Failed(c, constant.CMDB_HOST_NOT_FOUND, "用户组不存在")
		return
	}

	// 获取授权主机数量
	hostCount, _ := s.authDao.CountAuthorizedHosts("group", id)

	// 获取成员详细信息
	var members []model.MemberInfo
	for _, userId := range group.MemberIds {
		user := systemDao.GetSysAdminInfo(int(userId))
		if user.ID > 0 {
			members = append(members, model.MemberInfo{
				UserID:   user.ID,
				Username: user.Username,
				Nickname: user.Nickname,
				Email:    user.Email,
				Phone:    user.Phone,
			})
		}
	}

	vo := model.BastionUserGroupDetailVo{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		MemberIds:   []uint(group.MemberIds),
		Members:     members,
		HostCount:   int(hostCount),
		Status:      group.Status,
		CreateTime:  group.CreateTime,
		UpdateTime:  group.UpdateTime,
	}

	result.Success(c, vo)
}

// 获取用户组列表(分页)
func (s *BastionUserGroupServiceImpl) GetUserGroupList(c *gin.Context, page, pageSize int) {
	list, total := s.dao.GetUserGroupListWithPage(page, pageSize)

	var vos []model.BastionUserGroupVo
	for _, group := range list {
		// 获取授权主机数量
		hostCount, _ := s.authDao.CountAuthorizedHosts("group", group.ID)

		vos = append(vos, model.BastionUserGroupVo{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			MemberIds:   []uint(group.MemberIds),
			MemberCount: len(group.MemberIds),
			HostCount:   int(hostCount),
			Status:      group.Status,
			CreateTime:  group.CreateTime,
			UpdateTime:  group.UpdateTime,
		})
	}

	result.SuccessWithPage(c, vos, total, page, pageSize)
}

// 添加用户组成员
func (s *BastionUserGroupServiceImpl) AddGroupMembers(c *gin.Context, dto *model.AddGroupMembersDto) {
	// 检查用户组是否存在
	_, err := s.dao.GetUserGroupById(dto.GroupID)
	if err != nil {
		result.Failed(c, constant.CMDB_HOST_NOT_FOUND, "用户组不存在")
		return
	}

	// 验证用户ID是否存在
	for _, userId := range dto.UserIDs {
		user := systemDao.GetSysAdminInfo(int(userId))
		if user.ID == 0 {
			result.Failed(c, constant.INVALID_PARAMS, fmt.Sprintf("用户ID %d 不存在", userId))
			return
		}
	}

	// 添加成员
	if err := s.dao.AddGroupMembers(dto.GroupID, dto.UserIDs); err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "添加用户组成员失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{"message": fmt.Sprintf("成功添加 %d 个用户到用户组", len(dto.UserIDs))})
}

// 删除用户组成员
func (s *BastionUserGroupServiceImpl) RemoveGroupMembers(c *gin.Context, dto *model.RemoveGroupMembersDto) {
	// 检查用户组是否存在
	_, err := s.dao.GetUserGroupById(dto.GroupID)
	if err != nil {
		result.Failed(c, constant.CMDB_HOST_NOT_FOUND, "用户组不存在")
		return
	}

	// 删除成员
	if err := s.dao.RemoveGroupMembers(dto.GroupID, dto.UserIDs); err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "删除用户组成员失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{"message": fmt.Sprintf("成功从用户组移除 %d 个用户", len(dto.UserIDs))})
}

// 获取用户组成员列表
func (s *BastionUserGroupServiceImpl) GetGroupMembers(c *gin.Context, groupId uint) {
	// 检查用户组是否存在
	group, err := s.dao.GetUserGroupById(groupId)
	if err != nil {
		result.Failed(c, constant.CMDB_HOST_NOT_FOUND, "用户组不存在")
		return
	}

	// 获取成员详细信息
	var members []systemModel.SysAdminInfo
	for _, userId := range group.MemberIds {
		user := systemDao.GetSysAdminInfo(int(userId))
		if user.ID > 0 {
			members = append(members, user)
		}
	}

	result.Success(c, gin.H{
		"groupId":     groupId,
		"groupName":   group.Name,
		"members":     members,
		"memberCount": len(members),
	})
}

// 获取用户所属的用户组
func (s *BastionUserGroupServiceImpl) GetUserGroups(c *gin.Context, userId uint) {
	// 验证用户是否存在
	user := systemDao.GetSysAdminInfo(int(userId))
	if user.ID == 0 {
		result.Failed(c, constant.INVALID_PARAMS, "用户不存在")
		return
	}

	// 获取用户所属的用户组
	groups, err := s.dao.GetUserGroups(userId)
	if err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "获取用户组失败: "+err.Error())
		return
	}

	var vos []model.BastionUserGroupVo
	for _, group := range groups {
		// 获取授权主机数量
		hostCount, _ := s.authDao.CountAuthorizedHosts("group", group.ID)

		vos = append(vos, model.BastionUserGroupVo{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			MemberIds:   []uint(group.MemberIds),
			MemberCount: len(group.MemberIds),
			HostCount:   int(hostCount),
			Status:      group.Status,
			CreateTime:  group.CreateTime,
			UpdateTime:  group.UpdateTime,
		})
	}

	result.Success(c, gin.H{
		"userId": userId,
		"groups": vos,
		"total":  len(vos),
	})
}

func GetBastionUserGroupService() BastionUserGroupServiceInterface {
	return &BastionUserGroupServiceImpl{
		dao:     dao.NewBastionUserGroupDao(),
		authDao: dao.NewBastionHostAuthorizationDao(),
	}
}
