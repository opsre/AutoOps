package service

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	configDao "dodevops-api/api/configcenter/dao"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type BastionHostAuthorizationServiceInterface interface {
	AuthorizeHostsToUser(c *gin.Context, dto *model.AuthorizeHostsToUserDto)      // 授权主机给用户
	AuthorizeHostsToGroup(c *gin.Context, dto *model.AuthorizeHostsToGroupDto)    // 授权主机给用户组
	RevokeHostsFromUser(c *gin.Context, dto *model.RevokeHostsFromUserDto)        // 取消用户的主机授权
	RevokeHostsFromGroup(c *gin.Context, dto *model.RevokeHostsFromGroupDto)      // 取消用户组的主机授权
	GetAuthorizedHosts(c *gin.Context, dto *model.GetAuthorizedHostsDto)          // 获取授权主机列表
	GetUserAuthorizedHosts(c *gin.Context, userId uint, page, pageSize int)       // 获取用户有权限的主机列表(分页)
	CheckUserHostPermission(c *gin.Context, userId, hostId uint)                  // 检查用户是否有主机权限
}

type BastionHostAuthorizationServiceImpl struct {
	authDao  dao.BastionHostAuthorizationDao
	hostDao  dao.CmdbHostDao
	groupDao dao.BastionUserGroupDao
	keyDao   *configDao.KeyManageDao
}

// 授权主机给用户
func (s *BastionHostAuthorizationServiceImpl) AuthorizeHostsToUser(c *gin.Context, dto *model.AuthorizeHostsToUserDto) {
	// 验证主机ID是否存在
	for _, hostId := range dto.HostIDs {
		_, err := s.hostDao.GetCmdbHostById(hostId)
		if err != nil {
			result.Failed(c, constant.INVALID_PARAMS, fmt.Sprintf("主机ID %d 不存在", hostId))
			return
		}
	}

	// 批量创建授权记录(过滤已存在的)
	var auths []model.BastionHostAuthorization
	var successCount, skipCount int

	for _, hostId := range dto.HostIDs {
		// 检查是否已授权
		if s.authDao.CheckAuthorization("user", dto.UserID, hostId) {
			skipCount++
			continue
		}

		auths = append(auths, model.BastionHostAuthorization{
			SubjectType: "user",
			SubjectID:   dto.UserID,
			HostID:      hostId,
			CreateTime:  util.HTime{Time: time.Now()},
		})
		successCount++
	}

	if len(auths) > 0 {
		if err := s.authDao.BatchAuthorizeHosts(auths); err != nil {
			result.Failed(c, constant.INVALID_PARAMS, "授权失败: "+err.Error())
			return
		}
	}

	result.Success(c, gin.H{
		"message":      fmt.Sprintf("授权成功: 新增%d条, 跳过%d条已存在", successCount, skipCount),
		"successCount": successCount,
		"skipCount":    skipCount,
	})
}

// 授权主机给用户组
func (s *BastionHostAuthorizationServiceImpl) AuthorizeHostsToGroup(c *gin.Context, dto *model.AuthorizeHostsToGroupDto) {
	// 验证用户组是否存在
	_, err := s.groupDao.GetUserGroupById(dto.GroupID)
	if err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "用户组不存在")
		return
	}

	// 验证主机ID是否存在
	for _, hostId := range dto.HostIDs {
		_, err := s.hostDao.GetCmdbHostById(hostId)
		if err != nil {
			result.Failed(c, constant.INVALID_PARAMS, fmt.Sprintf("主机ID %d 不存在", hostId))
			return
		}
	}

	// 批量创建授权记录(过滤已存在的)
	var auths []model.BastionHostAuthorization
	var successCount, skipCount int

	for _, hostId := range dto.HostIDs {
		// 检查是否已授权
		if s.authDao.CheckAuthorization("group", dto.GroupID, hostId) {
			skipCount++
			continue
		}

		auths = append(auths, model.BastionHostAuthorization{
			SubjectType: "group",
			SubjectID:   dto.GroupID,
			HostID:      hostId,
			CreateTime:  util.HTime{Time: time.Now()},
		})
		successCount++
	}

	if len(auths) > 0 {
		if err := s.authDao.BatchAuthorizeHosts(auths); err != nil {
			result.Failed(c, constant.INVALID_PARAMS, "授权失败: "+err.Error())
			return
		}
	}

	result.Success(c, gin.H{
		"message":      fmt.Sprintf("授权成功: 新增%d条, 跳过%d条已存在", successCount, skipCount),
		"successCount": successCount,
		"skipCount":    skipCount,
	})
}

// 取消用户的主机授权
func (s *BastionHostAuthorizationServiceImpl) RevokeHostsFromUser(c *gin.Context, dto *model.RevokeHostsFromUserDto) {
	// 批量取消授权
	if err := s.authDao.BatchRevokeHostAuthorization("user", dto.UserID, dto.HostIDs); err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "取消授权失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{
		"message": fmt.Sprintf("成功取消 %d 台主机的授权", len(dto.HostIDs)),
	})
}

// 取消用户组的主机授权
func (s *BastionHostAuthorizationServiceImpl) RevokeHostsFromGroup(c *gin.Context, dto *model.RevokeHostsFromGroupDto) {
	// 验证用户组是否存在
	_, err := s.groupDao.GetUserGroupById(dto.GroupID)
	if err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "用户组不存在")
		return
	}

	// 批量取消授权
	if err := s.authDao.BatchRevokeHostAuthorization("group", dto.GroupID, dto.HostIDs); err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "取消授权失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{
		"message": fmt.Sprintf("成功取消 %d 台主机的授权", len(dto.HostIDs)),
	})
}

// 获取授权主机列表
func (s *BastionHostAuthorizationServiceImpl) GetAuthorizedHosts(c *gin.Context, dto *model.GetAuthorizedHostsDto) {
	// 验证subject是否存在
	if dto.SubjectType == "group" {
		_, err := s.groupDao.GetUserGroupById(dto.SubjectID)
		if err != nil {
			result.Failed(c, constant.INVALID_PARAMS, "用户组不存在")
			return
		}
	}

	// 获取授权主机ID列表
	auths, total, err := s.authDao.GetAuthorizedHostsWithPage(dto.SubjectType, dto.SubjectID, dto.Page, dto.PageSize)
	if err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "获取授权主机失败: "+err.Error())
		return
	}

	// 查询主机详细信息
	var vos []model.CmdbHostVo
	for _, auth := range auths {
		host, err := s.hostDao.GetCmdbHostById(auth.HostID)
		if err != nil {
			continue
		}

		// 查询SSH密钥名称
		sshKeyName := ""
		if host.SSHKeyID > 0 {
			keyManage, err := s.keyDao.GetByID(host.SSHKeyID)
			if err == nil {
				sshKeyName = keyManage.KeyName
			}
		}

		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   host.Group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHKeyName:  sshKeyName,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Tag:         host.Tag,
			Vendor:      host.Vendor,
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}

	result.SuccessWithPage(c, vos, total, dto.Page, dto.PageSize)
}

// 获取用户有权限的主机列表(分页+去重)
func (s *BastionHostAuthorizationServiceImpl) GetUserAuthorizedHosts(c *gin.Context, userId uint, page, pageSize int) {
	// 获取用户有权限的主机ID列表(已去重)
	hostIds, total, err := s.authDao.GetUserAuthorizedHostsWithPage(userId, page, pageSize)
	if err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "获取授权主机失败: "+err.Error())
		return
	}

	if len(hostIds) == 0 {
		result.SuccessWithPage(c, []model.CmdbHostVo{}, 0, page, pageSize)
		return
	}

	// 批量查询主机详细信息
	hosts, err := s.hostDao.GetCmdbHostsByIds(hostIds)
	if err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "查询主机信息失败: "+err.Error())
		return
	}

	// 转换为VO
	var vos []model.CmdbHostVo
	for _, host := range hosts {
		// 查询SSH密钥名称
		sshKeyName := ""
		if host.SSHKeyID > 0 {
			keyManage, err := s.keyDao.GetByID(host.SSHKeyID)
			if err == nil {
				sshKeyName = keyManage.KeyName
			}
		}

		vos = append(vos, model.CmdbHostVo{
			ID:          host.ID,
			HostName:    host.HostName,
			Name:        host.Name,
			GroupID:     host.GroupID,
			GroupName:   host.Group.Name,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			SSHName:     host.SSHName,
			SSHKeyID:    host.SSHKeyID,
			SSHKeyName:  sshKeyName,
			SSHPort:     host.SSHPort,
			Remark:      host.Remark,
			Tag:         host.Tag,
			Vendor:      host.Vendor,
			Region:      host.Region,
			InstanceID:  host.InstanceID,
			OS:          host.OS,
			Status:      host.Status,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			BillingType: host.BillingType,
			CreateTime:  host.CreateTime,
			ExpireTime:  host.ExpireTime,
			UpdateTime:  host.UpdateTime,
		})
	}

	result.SuccessWithPage(c, vos, total, page, pageSize)
}

// 检查用户是否有主机权限
func (s *BastionHostAuthorizationServiceImpl) CheckUserHostPermission(c *gin.Context, userId, hostId uint) {
	hasPermission, err := s.authDao.CheckUserHostPermission(userId, hostId)
	if err != nil {
		result.Failed(c, constant.INVALID_PARAMS, "检查权限失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{
		"userId":        userId,
		"hostId":        hostId,
		"hasPermission": hasPermission,
	})
}

func GetBastionHostAuthorizationService() BastionHostAuthorizationServiceInterface {
	return &BastionHostAuthorizationServiceImpl{
		authDao:  dao.NewBastionHostAuthorizationDao(),
		hostDao:  dao.NewCmdbHostDao(),
		groupDao: dao.NewBastionUserGroupDao(),
		keyDao:   configDao.NewKeyManageDao(),
	}
}
