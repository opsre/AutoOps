package dao

import (
	"dodevops-api/api/db/mysql/model"
	. "dodevops-api/pkg/db"
	"time"
)

// CreateDbSqlExec 创建SQL执行记录
func CreateDbSqlExec(exec *model.DbSqlExec) error {
	return Db.Create(exec).Error
}

// GetDbSqlExecById 根据ID查询SQL执行记录
func GetDbSqlExecById(id uint) (exec model.DbSqlExec, err error) {
	err = Db.Where("id = ?", id).First(&exec).Error
	return
}

// GetDbSqlExecList 获取SQL执行记录列表
func GetDbSqlExecList(query *model.DbSqlExecQuery) (execs []model.DbSqlExec, total int64, err error) {
	db := Db.Model(&model.DbSqlExec{})

	// 构建查询条件
	if query.DbId > 0 {
		db = db.Where("db_id = ?", query.DbId)
	}
	if query.DbName != "" {
		db = db.Where("db_name = ?", query.DbName)
	}
	if query.TableName != "" {
		db = db.Where("table_name LIKE ?", "%"+query.TableName+"%")
	}
	if query.Type > 0 {
		db = db.Where("type = ?", query.Type)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.Creator != "" {
		db = db.Where("creator LIKE ?", "%"+query.Creator+"%")
	}
	if query.StartTime != "" {
		db = db.Where("create_time >= ?", query.StartTime)
	}
	if query.EndTime != "" {
		db = db.Where("create_time <= ?", query.EndTime)
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

	err = db.Order("create_time DESC").Find(&execs).Error
	return
}

// DeleteDbSqlExecsByDbId 根据数据库ID删除SQL执行记录
func DeleteDbSqlExecsByDbId(dbId uint) error {
	return Db.Where("db_id = ?", dbId).Delete(&model.DbSqlExec{}).Error
}

// DeleteOldDbSqlExecs 删除旧的SQL执行记录
func DeleteOldDbSqlExecs(days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	return Db.Where("create_time < ?", cutoffTime).Delete(&model.DbSqlExec{}).Error
}

// GetDbSqlExecStatistics 获取SQL执行统计
func GetDbSqlExecStatistics(dbId uint, startTime, endTime string) (stats map[string]interface{}, err error) {
	stats = make(map[string]interface{})

	db := Db.Model(&model.DbSqlExec{}).Where("db_id = ?", dbId)

	if startTime != "" {
		db = db.Where("create_time >= ?", startTime)
	}
	if endTime != "" {
		db = db.Where("create_time <= ?", endTime)
	}

	// 总执行次数
	var totalCount int64
	db.Count(&totalCount)
	stats["totalCount"] = totalCount

	// 成功次数
	var successCount int64
	db.Where("status = ?", model.ExecStatusSuccess).Count(&successCount)
	stats["successCount"] = successCount

	// 失败次数
	stats["failCount"] = totalCount - successCount

	// 各类型SQL执行次数
	var typeStats []struct {
		Type  int8  `json:"type"`
		Count int64 `json:"count"`
	}
	db.Select("type, COUNT(*) as count").Group("type").Scan(&typeStats)
	stats["typeStats"] = typeStats

	return
}
