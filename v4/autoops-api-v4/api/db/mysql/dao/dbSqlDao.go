package dao

import (
	"dodevops-api/api/db/mysql/model"
	. "dodevops-api/pkg/db"
)

// CreateDbSql 创建SQL模板
func CreateDbSql(sql *model.DbSql) error {
	return Db.Create(sql).Error
}

// UpdateDbSql 更新SQL模板
func UpdateDbSql(sql *model.DbSql) error {
	return Db.Model(&model.DbSql{}).Where("id = ?", sql.ID).Updates(sql).Error
}

// DeleteDbSql 删除SQL模板
func DeleteDbSql(id uint) error {
	return Db.Delete(&model.DbSql{}, id).Error
}

// GetDbSqlById 根据ID查询SQL模板
func GetDbSqlById(id uint) (sql model.DbSql, err error) {
	err = Db.Where("id = ?", id).First(&sql).Error
	return
}

// GetDbSqlByName 根据名称查询SQL模板
func GetDbSqlByName(dbId uint, name string) (sql model.DbSql, err error) {
	err = Db.Where("db_id = ? AND name = ?", dbId, name).First(&sql).Error
	return
}

// GetDbSqlList 获取SQL模板列表
func GetDbSqlList(query *model.DbSqlQuery) (sqls []model.DbSql, total int64, err error) {
	db := Db.Model(&model.DbSql{})

	// 构建查询条件
	if query.DbId > 0 {
		db = db.Where("db_id = ?", query.DbId)
	}
	if query.Db != "" {
		db = db.Where("db = ?", query.Db)
	}
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
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

	err = db.Order("create_time DESC").Find(&sqls).Error
	return
}

// GetDbSqlNames 获取SQL模板名称列表
func GetDbSqlNames(dbId uint, db string) (names []string, err error) {
	err = Db.Model(&model.DbSql{}).
		Where("db_id = ? AND db = ?", dbId, db).
		Pluck("name", &names).Error
	return
}

// CheckDbSqlExists 检查SQL模板是否存在
func CheckDbSqlExists(dbId uint, name string, excludeId uint) bool {
	var count int64
	db := Db.Model(&model.DbSql{}).Where("db_id = ? AND name = ?", dbId, name)
	if excludeId > 0 {
		db = db.Where("id != ?", excludeId)
	}
	db.Count(&count)
	return count > 0
}

// DeleteDbSqlsByDbId 根据数据库ID删除SQL模板
func DeleteDbSqlsByDbId(dbId uint) error {
	return Db.Where("db_id = ?", dbId).Delete(&model.DbSql{}).Error
}
