package dao

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

type BastionHostAuthorizationDao struct {
	db *gorm.DB
}

func NewBastionHostAuthorizationDao() BastionHostAuthorizationDao {
	return BastionHostAuthorizationDao{
		db: common.GetDB(),
	}
}

// 授权主机
func (d *BastionHostAuthorizationDao) AuthorizeHost(auth *model.BastionHostAuthorization) error {
	return d.db.Create(auth).Error
}

// 批量授权主机
func (d *BastionHostAuthorizationDao) BatchAuthorizeHosts(auths []model.BastionHostAuthorization) error {
	return d.db.Create(&auths).Error
}

// 取消授权
func (d *BastionHostAuthorizationDao) RevokeHostAuthorization(subjectType string, subjectId, hostId uint) error {
	return d.db.Where("subject_type = ? AND subject_id = ? AND host_id = ?", subjectType, subjectId, hostId).
		Delete(&model.BastionHostAuthorization{}).Error
}

// 批量取消授权
func (d *BastionHostAuthorizationDao) BatchRevokeHostAuthorization(subjectType string, subjectId uint, hostIds []uint) error {
	return d.db.Where("subject_type = ? AND subject_id = ? AND host_id IN ?", subjectType, subjectId, hostIds).
		Delete(&model.BastionHostAuthorization{}).Error
}

// 删除主体的所有授权(用户或用户组)
func (d *BastionHostAuthorizationDao) RevokeAllSubjectAuthorization(subjectType string, subjectId uint) error {
	return d.db.Where("subject_type = ? AND subject_id = ?", subjectType, subjectId).
		Delete(&model.BastionHostAuthorization{}).Error
}

// 删除主机的所有授权
func (d *BastionHostAuthorizationDao) RevokeAllHostAuthorization(hostId uint) error {
	return d.db.Where("host_id = ?", hostId).Delete(&model.BastionHostAuthorization{}).Error
}

// 检查是否已授权该主机
func (d *BastionHostAuthorizationDao) CheckAuthorization(subjectType string, subjectId, hostId uint) bool {
	var count int64
	d.db.Model(&model.BastionHostAuthorization{}).
		Where("subject_type = ? AND subject_id = ? AND host_id = ?", subjectType, subjectId, hostId).
		Count(&count)
	return count > 0
}

// 获取已授权的主机ID列表
func (d *BastionHostAuthorizationDao) GetAuthorizedHostIds(subjectType string, subjectId uint) ([]uint, error) {
	var hostIds []uint
	err := d.db.Model(&model.BastionHostAuthorization{}).
		Where("subject_type = ? AND subject_id = ?", subjectType, subjectId).
		Pluck("host_id", &hostIds).Error
	return hostIds, err
}

// 获取已授权的主机列表(分页)
func (d *BastionHostAuthorizationDao) GetAuthorizedHostsWithPage(subjectType string, subjectId uint, page, pageSize int) ([]model.BastionHostAuthorization, int64, error) {
	var auths []model.BastionHostAuthorization
	var total int64

	query := d.db.Model(&model.BastionHostAuthorization{}).
		Where("subject_type = ? AND subject_id = ?", subjectType, subjectId)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&auths).Error
	return auths, total, err
}

// 获取用户有权限的主机ID列表(包括直接授权和通过用户组授权，去重)
func (d *BastionHostAuthorizationDao) GetUserAuthorizedHostIds(userId uint) ([]uint, error) {
	// 获取用户所属的用户组
	groupDao := NewBastionUserGroupDao()
	groups, err := groupDao.GetUserGroups(userId)
	if err != nil {
		return nil, err
	}

	// 收集用户组ID
	var groupIds []uint
	for _, group := range groups {
		groupIds = append(groupIds, group.ID)
	}

	// 构建查询条件：直接授权给用户 OR 授权给用户所在的用户组
	var hostIds []uint
	query := d.db.Model(&model.BastionHostAuthorization{})

	if len(groupIds) > 0 {
		// 用户有用户组：查询直接授权 + 用户组授权
		query = query.Where("(subject_type = ? AND subject_id = ?) OR (subject_type = ? AND subject_id IN ?)",
			"user", userId, "group", groupIds)
	} else {
		// 用户没有用户组：只查询直接授权
		query = query.Where("subject_type = ? AND subject_id = ?", "user", userId)
	}

	err = query.Distinct("host_id").Pluck("host_id", &hostIds).Error
	return hostIds, err
}

// 获取用户有权限的主机列表(分页+去重)
func (d *BastionHostAuthorizationDao) GetUserAuthorizedHostsWithPage(userId uint, page, pageSize int) ([]uint, int64, error) {
	// 获取用户所属的用户组
	groupDao := NewBastionUserGroupDao()
	groups, err := groupDao.GetUserGroups(userId)
	if err != nil {
		return nil, 0, err
	}

	// 收集用户组ID
	var groupIds []uint
	for _, group := range groups {
		groupIds = append(groupIds, group.ID)
	}

	// 构建查询条件
	query := d.db.Model(&model.BastionHostAuthorization{})

	if len(groupIds) > 0 {
		query = query.Where("(subject_type = ? AND subject_id = ?) OR (subject_type = ? AND subject_id IN ?)",
			"user", userId, "group", groupIds)
	} else {
		query = query.Where("subject_type = ? AND subject_id = ?", "user", userId)
	}

	// 统计总数(去重)
	var total int64
	err = query.Distinct("host_id").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询主机ID(去重)
	var hostIds []uint
	err = query.Distinct("host_id").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Pluck("host_id", &hostIds).Error

	return hostIds, total, err
}

// 检查用户是否有主机权限
func (d *BastionHostAuthorizationDao) CheckUserHostPermission(userId, hostId uint) (bool, error) {
	// 首先检查直接授权
	var count int64
	err := d.db.Model(&model.BastionHostAuthorization{}).
		Where("subject_type = ? AND subject_id = ? AND host_id = ?", "user", userId, hostId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	// 检查用户组授权
	groupDao := NewBastionUserGroupDao()
	groups, err := groupDao.GetUserGroups(userId)
	if err != nil {
		return false, err
	}

	if len(groups) == 0 {
		return false, nil
	}

	// 收集用户组ID
	var groupIds []uint
	for _, group := range groups {
		groupIds = append(groupIds, group.ID)
	}

	// 检查用户组授权
	err = d.db.Model(&model.BastionHostAuthorization{}).
		Where("subject_type = ? AND subject_id IN ? AND host_id = ?", "group", groupIds, hostId).
		Count(&count).Error

	return count > 0, err
}

// 统计授权主机数量
func (d *BastionHostAuthorizationDao) CountAuthorizedHosts(subjectType string, subjectId uint) (int64, error) {
	var count int64
	err := d.db.Model(&model.BastionHostAuthorization{}).
		Where("subject_type = ? AND subject_id = ?", subjectType, subjectId).
		Count(&count).Error
	return count, err
}

// 统计用户有权限的主机数量(去重)
func (d *BastionHostAuthorizationDao) CountUserAuthorizedHosts(userId uint) (int64, error) {
	// 获取用户所属的用户组
	groupDao := NewBastionUserGroupDao()
	groups, err := groupDao.GetUserGroups(userId)
	if err != nil {
		return 0, err
	}

	// 收集用户组ID
	var groupIds []uint
	for _, group := range groups {
		groupIds = append(groupIds, group.ID)
	}

	// 构建查询条件
	query := d.db.Model(&model.BastionHostAuthorization{})

	if len(groupIds) > 0 {
		query = query.Where("(subject_type = ? AND subject_id = ?) OR (subject_type = ? AND subject_id IN ?)",
			"user", userId, "group", groupIds)
	} else {
		query = query.Where("subject_type = ? AND subject_id = ?", "user", userId)
	}

	// 统计去重后的主机数量
	var count int64
	err = query.Distinct("host_id").Count(&count).Error
	return count, err
}
