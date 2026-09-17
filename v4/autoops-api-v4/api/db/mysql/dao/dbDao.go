package dao

import (
	"dodevops-api/api/db/mysql/model"
	. "dodevops-api/pkg/db"
)

// CreateDb 创建数据库
func CreateDb(db *model.Db) error {
	return Db.Create(db).Error
}

// UpdateDb 更新数据库
func UpdateDb(db *model.Db) error {
	return Db.Model(&model.Db{}).Where("id = ?", db.ID).Updates(db).Error
}

// DeleteDb 删除数据库
func DeleteDb(id uint) error {
	return Db.Delete(&model.Db{}, id).Error
}

// GetDbById 根据ID查询数据库
func GetDbById(id uint) (db model.Db, err error) {
	err = Db.Where("id = ?", id).First(&db).Error
	return
}

// GetDbByCode 根据编码查询数据库
func GetDbByCode(code string) (db model.Db, err error) {
	err = Db.Where("code = ?", code).First(&db).Error
	return
}

// GetDbList 获取数据库列表
func GetDbList(query *model.DbQuery) (dbs []model.Db, total int64, err error) {
	db := Db.Model(&model.Db{})

	// 构建查询条件
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Database != "" {
		db = db.Where("database LIKE ?", "%"+query.Database+"%")
	}
	if query.InstanceId > 0 {
		db = db.Where("instance_id = ?", query.InstanceId)
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

	err = db.Order("create_time DESC").Find(&dbs).Error
	return
}

// GetDbListWithInstance 获取数据库列表(包含实例信息)
func GetDbListWithInstance(query *model.DbQuery) (dbs []model.DbListVo, total int64, err error) {
	db := Db.Model(&model.Db{}).Select(`
		db.id, db.code, db.name, db.database, db.remark,
		db.instance_id, db.instance_code, db.status,
		db.create_time, db.creator,
		di.name as instance_name, di.type as instance_type,
		di.host, di.port
	`).Joins("LEFT JOIN db_instance di ON db.instance_id = di.id")

	// 构建查询条件
	if query.Name != "" {
		db = db.Where("db.name LIKE ?", "%"+query.Name+"%")
	}
	if query.Database != "" {
		db = db.Where("db.database LIKE ?", "%"+query.Database+"%")
	}
	if query.InstanceId > 0 {
		db = db.Where("db.instance_id = ?", query.InstanceId)
	}
	if query.Status > 0 {
		db = db.Where("db.status = ?", query.Status)
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

	err = db.Order("db.create_time DESC").Find(&dbs).Error
	return
}

// CheckDbExists 检查数据库是否存在
func CheckDbExists(code string, excludeId uint) bool {
	var count int64
	db := Db.Model(&model.Db{}).Where("code = ?", code)
	if excludeId > 0 {
		db = db.Where("id != ?", excludeId)
	}
	db.Count(&count)
	return count > 0
}

// GetDbsByInstanceId 根据实例ID获取数据库列表
func GetDbsByInstanceId(instanceId uint) (dbs []model.Db, err error) {
	err = Db.Where("instance_id = ? AND status = ?", instanceId, 1).Find(&dbs).Error
	return
}

// UpdateDbStatus 更新数据库状态
func UpdateDbStatus(id uint, status int) error {
	return Db.Model(&model.Db{}).Where("id = ?", id).Update("status", status).Error
}
