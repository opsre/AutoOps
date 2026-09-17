package model

import (
	"dodevops-api/common/util"
)

// 堡垒机主机授权模型
// 支持对用户组授权（subject_type='group'），也支持对单个用户授权（subject_type='user'）
type BastionHostAuthorization struct {
	ID          uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	SubjectType string     `gorm:"column:subject_type;varchar(10);comment:'授权对象类型:user/group';NOT NULL;index:idx_subject_host,priority:1" json:"subjectType"`
	SubjectID   uint       `gorm:"column:subject_id;comment:'授权对象ID(用户ID或用户组ID)';NOT NULL;index:idx_subject_host,priority:2" json:"subjectId"`
	HostID      uint       `gorm:"column:host_id;comment:'主机ID';NOT NULL;index:idx_subject_host,priority:3" json:"hostId"`
	CreateTime  util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
}

func (BastionHostAuthorization) TableName() string {
	return "bastion_host_authorization"
}

// 授权主机给用户DTO
type AuthorizeHostsToUserDto struct {
	UserID  uint   `validate:"required" json:"userId"`  // 用户ID
	HostIDs []uint `validate:"required" json:"hostIds"` // 主机ID列表
}

// 授权主机给用户组DTO
type AuthorizeHostsToGroupDto struct {
	GroupID uint   `validate:"required" json:"groupId"` // 用户组ID
	HostIDs []uint `validate:"required" json:"hostIds"` // 主机ID列表
}

// 取消用户的主机授权DTO
type RevokeHostsFromUserDto struct {
	UserID  uint   `validate:"required" json:"userId"`  // 用户ID
	HostIDs []uint `validate:"required" json:"hostIds"` // 主机ID列表
}

// 取消用户组的主机授权DTO
type RevokeHostsFromGroupDto struct {
	GroupID uint   `validate:"required" json:"groupId"` // 用户组ID
	HostIDs []uint `validate:"required" json:"hostIds"` // 主机ID列表
}

// 查询用户/用户组已授权主机DTO
type GetAuthorizedHostsDto struct {
	SubjectType string `form:"subjectType" binding:"required,oneof=user group"` // 授权对象类型
	SubjectID   uint   `form:"subjectId" binding:"required"`                    // 授权对象ID
	Page        int    `form:"page" binding:"required,min=1"`
	PageSize    int    `form:"pageSize" binding:"required,min=1,max=100"`
}

// 检查用户是否有主机权限DTO
type CheckUserHostPermissionDto struct {
	UserID uint `form:"userId" validate:"required"` // 用户ID
	HostID uint `form:"hostId" validate:"required"` // 主机ID
}

// 授权VO
type BastionHostAuthorizationVo struct {
	ID          uint       `json:"id"`
	SubjectType string     `json:"subjectType"` // user/group
	SubjectID   uint       `json:"subjectId"`
	SubjectName string     `json:"subjectName"` // 用户名或用户组名
	HostID      uint       `json:"hostId"`
	HostName    string     `json:"hostName"`
	SSHIP       string     `json:"sshIp"`
	CreateTime  util.HTime `json:"createTime"`
}
