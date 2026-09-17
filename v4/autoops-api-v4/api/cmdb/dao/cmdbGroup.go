package dao

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common"
	"errors"

	"gorm.io/gorm"
)

type CmdbGroupDao struct {
	db *gorm.DB
}

func NewCmdbGroupDao() CmdbGroupDao {
	return CmdbGroupDao{
		db: common.GetDB(),
	}
}

func (d *CmdbGroupDao) GetCmdbGroupById(id uint) (model.CmdbGroup, error) {
	var group model.CmdbGroup
	err := d.db.Where("id = ?", id).First(&group).Error
	return group, err
}

func (d *CmdbGroupDao) GetCmdbGroupByName(name string) (model.CmdbGroup, error) {
	var group model.CmdbGroup
	err := d.db.Where("name LIKE ?", "%"+name+"%").First(&group).Error
	return group, err
}

// CheckNameExists 检查分组名称是否存在（精确匹配）
// 如果 parentID 为 0，则检查一级分组中是否存在同名；
// 如果 parentID > 0，则检查所有二级分组中是否存在同名（不限父分组）
func (d *CmdbGroupDao) CheckNameExists(name string, parentID uint) bool {
	var count int64
	if parentID == 0 {
		// 创建一级分组：只检查一级分组中是否存在同名
		d.db.Model(&model.CmdbGroup{}).Where("name = ? AND parent_id = 0", name).Count(&count)
	} else {
		// 创建二级分组：检查所有二级分组（parent_id > 0）中是否存在同名
		d.db.Model(&model.CmdbGroup{}).Where("name = ? AND parent_id > 0", name).Count(&count)
	}
	return count > 0
}

func (d *CmdbGroupDao) CreateCmdbGroup(group *model.CmdbGroup) error {
	return d.db.Create(group).Error
}

func (d *CmdbGroupDao) UpdateCmdbGroup(id uint, group *model.CmdbGroup) error {
	return d.db.Model(&model.CmdbGroup{}).Where("id = ?", id).Updates(group).Error
}

// 递归获取所有子分组ID（包括自身）
func (d *CmdbGroupDao) getAllChildGroupIds(groupId uint) []uint {
	var groupIds []uint
	groupIds = append(groupIds, groupId) // 包含自身

	var childGroups []model.CmdbGroup
	d.db.Where("parent_id = ?", groupId).Find(&childGroups)

	for _, child := range childGroups {
		childIds := d.getAllChildGroupIds(child.ID)
		groupIds = append(groupIds, childIds...)
	}

	return groupIds
}

func (d *CmdbGroupDao) DeleteCmdbGroup(id uint) error {
	// 获取当前分组及所有子分组的ID
	groupIds := d.getAllChildGroupIds(id)

	// 检查这些分组中是否有主机
	var hostCount int64
	d.db.Model(&model.CmdbHost{}).Where("group_id IN ?", groupIds).Count(&hostCount)
	if hostCount > 0 {
		return errors.New("该分组或子分组下存在主机，无法删除")
	}

	// 批量删除所有分组（包括父分组和所有子分组）
	return d.db.Where("id IN ?", groupIds).Delete(&model.CmdbGroup{}).Error
}

func (d *CmdbGroupDao) GetCmdbGroupList() []model.CmdbGroup {
	var list []model.CmdbGroup
	d.db.Find(&list)
	return list
}

// GetCmdbGroupByParentAndName 根据父分组名称和子分组名称获取分组
// 返回匹配的子分组,如果不存在则返回nil
func (d *CmdbGroupDao) GetCmdbGroupByParentAndName(parentName, childName string) (*model.CmdbGroup, error) {
	// 1. 先查找父分组
	var parentGroup model.CmdbGroup
	if err := d.db.Where("name = ? AND parent_id = 0", parentName).First(&parentGroup).Error; err != nil {
		return nil, err
	}

	// 2. 查找该父分组下的子分组
	var childGroup model.CmdbGroup
	if err := d.db.Where("name = ? AND parent_id = ?", childName, parentGroup.ID).First(&childGroup).Error; err != nil {
		return nil, err
	}

	return &childGroup, nil
}

// GetCmdbGroupByExactName 根据精确名称获取分组(不使用模糊匹配)
func (d *CmdbGroupDao) GetCmdbGroupByExactName(name string) (*model.CmdbGroup, error) {
	var group model.CmdbGroup
	err := d.db.Where("name = ?", name).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// IsSecondLevelGroup 检查是否为二级分组(ParentID != 0)
func (d *CmdbGroupDao) IsSecondLevelGroup(groupID uint) (bool, error) {
	var group model.CmdbGroup
	err := d.db.Select("parent_id").Where("id = ?", groupID).First(&group).Error
	if err != nil {
		return false, err
	}
	return group.ParentID != 0, nil
}

// GetAllChildGroupIds 获取当前分组及所有子分组的ID（用于外部调用）
func (d *CmdbGroupDao) GetAllChildGroupIds(groupId uint) []uint {
	return d.getAllChildGroupIds(groupId)
}
