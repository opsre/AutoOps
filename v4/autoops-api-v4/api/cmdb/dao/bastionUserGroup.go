package dao

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

type BastionUserGroupDao struct {
	db *gorm.DB
}

func NewBastionUserGroupDao() BastionUserGroupDao {
	return BastionUserGroupDao{
		db: common.GetDB(),
	}
}

// 创建用户组
func (d *BastionUserGroupDao) CreateUserGroup(group *model.BastionUserGroup) error {
	return d.db.Create(group).Error
}

// 更新用户组
func (d *BastionUserGroupDao) UpdateUserGroup(id uint, group *model.BastionUserGroup) error {
	return d.db.Model(&model.BastionUserGroup{}).Where("id = ?", id).Updates(group).Error
}

// 删除用户组
func (d *BastionUserGroupDao) DeleteUserGroup(id uint) error {
	return d.db.Delete(&model.BastionUserGroup{}, id).Error
}

// 根据ID获取用户组
func (d *BastionUserGroupDao) GetUserGroupById(id uint) (*model.BastionUserGroup, error) {
	var group model.BastionUserGroup
	err := d.db.Where("id = ?", id).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// 根据名称获取用户组
func (d *BastionUserGroupDao) GetUserGroupByName(name string) (*model.BastionUserGroup, error) {
	var group model.BastionUserGroup
	err := d.db.Where("name = ?", name).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// 检查用户组名称是否存在
func (d *BastionUserGroupDao) CheckNameExists(name string) bool {
	var count int64
	d.db.Model(&model.BastionUserGroup{}).Where("name = ?", name).Count(&count)
	return count > 0
}

// 检查用户组名称是否存在(排除指定ID)
func (d *BastionUserGroupDao) CheckNameExistsExcludeId(name string, id uint) bool {
	var count int64
	d.db.Model(&model.BastionUserGroup{}).Where("name = ? AND id != ?", name, id).Count(&count)
	return count > 0
}

// 获取用户组列表(分页)
func (d *BastionUserGroupDao) GetUserGroupListWithPage(page, pageSize int) ([]model.BastionUserGroup, int64) {
	var list []model.BastionUserGroup
	var total int64

	d.db.Model(&model.BastionUserGroup{}).Count(&total)
	d.db.Order("create_time DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)

	return list, total
}

// 获取用户组列表
func (d *BastionUserGroupDao) GetUserGroupList() ([]model.BastionUserGroup, error) {
	var list []model.BastionUserGroup
	err := d.db.Order("create_time DESC").Find(&list).Error
	return list, err
}

// 添加用户组成员
func (d *BastionUserGroupDao) AddGroupMembers(groupId uint, userIds []uint) error {
	// 获取当前用户组
	group, err := d.GetUserGroupById(groupId)
	if err != nil {
		return err
	}

	// 合并成员ID(去重)
	memberMap := make(map[uint]bool)
	for _, id := range group.MemberIds {
		memberMap[id] = true
	}
	for _, id := range userIds {
		memberMap[id] = true
	}

	// 转换回数组
	var newMemberIds []uint
	for id := range memberMap {
		newMemberIds = append(newMemberIds, id)
	}

	// 更新用户组
	return d.db.Model(&model.BastionUserGroup{}).
		Where("id = ?", groupId).
		Update("member_ids", model.UintArray(newMemberIds)).Error
}

// 删除用户组成员
func (d *BastionUserGroupDao) RemoveGroupMembers(groupId uint, userIds []uint) error {
	// 获取当前用户组
	group, err := d.GetUserGroupById(groupId)
	if err != nil {
		return err
	}

	// 创建要删除的用户ID map
	removeMap := make(map[uint]bool)
	for _, id := range userIds {
		removeMap[id] = true
	}

	// 过滤成员ID
	var newMemberIds []uint
	for _, id := range group.MemberIds {
		if !removeMap[id] {
			newMemberIds = append(newMemberIds, id)
		}
	}

	// 更新用户组
	return d.db.Model(&model.BastionUserGroup{}).
		Where("id = ?", groupId).
		Update("member_ids", model.UintArray(newMemberIds)).Error
}

// 获取用户所属的所有用户组
func (d *BastionUserGroupDao) GetUserGroups(userId uint) ([]model.BastionUserGroup, error) {
	var groups []model.BastionUserGroup
	// 使用JSON_CONTAINS查询，第二个参数需要是JSON格式的数字
	// 使用CAST转换为JSON，或直接传递字符串格式的JSON
	err := d.db.Where("JSON_CONTAINS(member_ids, CAST(? AS JSON))", userId).Find(&groups).Error
	return groups, err
}
