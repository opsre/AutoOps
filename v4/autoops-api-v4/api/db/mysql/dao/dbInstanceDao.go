package dao

import (
	"dodevops-api/api/db/mysql/model"
	. "dodevops-api/pkg/db"
	"fmt"
)

// CreateDbInstance 创建数据库实例
func CreateDbInstance(instance *model.DbInstance) error {
	return Db.Create(instance).Error
}

// UpdateDbInstance 更新数据库实例
func UpdateDbInstance(instance *model.DbInstance) error {
	return Db.Model(&model.DbInstance{}).Where("id = ?", instance.ID).Updates(instance).Error
}

// DeleteDbInstance 删除数据库实例
func DeleteDbInstance(id uint) error {
	return Db.Delete(&model.DbInstance{}, id).Error
}

// GetDbInstanceById 根据ID查询实例
func GetDbInstanceById(id uint) (instance model.DbInstance, err error) {
	err = Db.Where("id = ?", id).First(&instance).Error
	return
}

// GetDbInstanceByCode 根据编码查询实例
func GetDbInstanceByCode(code string) (instance model.DbInstance, err error) {
	err = Db.Where("code = ?", code).First(&instance).Error
	return
}

// GetDbInstanceList 获取实例列表
func GetDbInstanceList(query *model.DbInstanceQuery) (instances []model.DbInstance, total int64, err error) {
	db := Db.Model(&model.DbInstance{})

	// 构建查询条件
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.Host != "" {
		db = db.Where("host LIKE ?", "%"+query.Host+"%")
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}

	// 统计总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 分页查询
	if query.PageNum > 0 && query.PageSize > 0 {
		offset := (query.PageNum - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	err = db.Order("create_time DESC").Find(&instances).Error
	return
}

// CheckInstanceExists 检查实例是否存在
func CheckInstanceExists(code string, excludeId uint) bool {
	var count int64
	db := Db.Model(&model.DbInstance{}).Where("code = ?", code)
	if excludeId > 0 {
		db = db.Where("id != ?", excludeId)
	}
	db.Count(&count)
	return count > 0
}

// UpdateInstanceStatus 更新实例状态
func UpdateInstanceStatus(id uint, status int) error {
	return Db.Model(&model.DbInstance{}).Where("id = ?", id).Update("status", status).Error
}

// GetAllActiveInstances 获取所有启用的实例
func GetAllActiveInstances() (instances []model.DbInstance, err error) {
	err = Db.Where("status = ?", 1).Order("create_time DESC").Find(&instances).Error
	return
}

// BatchDeleteDbInstances 批量删除实例
func BatchDeleteDbInstances(ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids cannot be empty")
	}
	return Db.Where("id IN ?", ids).Delete(&model.DbInstance{}).Error
}
