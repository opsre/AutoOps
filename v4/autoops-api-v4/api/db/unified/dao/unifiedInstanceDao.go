package dao

import (
	"dodevops-api/api/db/unified/model"
	dbmodel "dodevops-api/api/db/mysql/model"
	redismodel "dodevops-api/api/db/redis/model"
	mongomodel "dodevops-api/api/db/mongo/model"
	esmodel "dodevops-api/api/db/elasticsearch/model"
	"dodevops-api/pkg/db"
	"errors"
)

// GetUnifiedInstanceList 获取统一实例列表
func GetUnifiedInstanceList(query *model.UnifiedInstanceQuery) ([]model.UnifiedInstance, int64, error) {
	var list []model.UnifiedInstance
	var total int64

	// 构建查询
	dbQuery := db.Db.Model(&model.UnifiedInstance{})

	// 条件过滤
	if query.DbType != "" {
		dbQuery = dbQuery.Where("db_type = ?", query.DbType)
	}
	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Status != nil {
		dbQuery = dbQuery.Where("status = ?", *query.Status)
	}

	// 查询总数
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (query.PageNum - 1) * query.PageSize
	if err := dbQuery.
		Order("create_time DESC").
		Limit(query.PageSize).
		Offset(offset).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetUnifiedInstanceById 根据类型和ID获取实例
func GetUnifiedInstanceById(dbType string, id uint) (*model.UnifiedInstance, error) {
	var instance model.UnifiedInstance
	err := db.Db.Where("db_type = ? AND id = ?", dbType, id).First(&instance).Error
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// GetUnifiedInstanceByCode 根据 code 获取实例
func GetUnifiedInstanceByCode(code string) (*model.UnifiedInstance, error) {
	var instance model.UnifiedInstance
	err := db.Db.Where("code = ?", code).First(&instance).Error
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// GetUnifiedInstanceStats 获取统计信息
func GetUnifiedInstanceStats() (map[string]int64, error) {
	type Stats struct {
		DbType string `gorm:"column:db_type"`
		Count  int64  `gorm:"column:count"`
	}

	var stats []Stats
	err := db.Db.Model(&model.UnifiedInstance{}).
		Select("db_type, COUNT(*) as count").
		Group("db_type").
		Find(&stats).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int64)
	for _, stat := range stats {
		result[stat.DbType] = stat.Count
	}

	return result, nil
}

// GetInstancePassword 根据类型和ID获取数据库实例的真实密码
func GetInstancePassword(dbType string, id uint) (string, error) {
	switch dbType {
	case "mysql", "postgres":
		var instance dbmodel.DbInstance
		err := db.Db.Where("id = ?", id).First(&instance).Error
		if err != nil {
			return "", err
		}
		return instance.Password, nil

	case "redis":
		var instance redismodel.RedisInstance
		err := db.Db.Where("id = ?", id).First(&instance).Error
		if err != nil {
			return "", err
		}
		return instance.Password, nil

	case "mongodb":
		var instance mongomodel.MongoInstance
		err := db.Db.Where("id = ?", id).First(&instance).Error
		if err != nil {
			return "", err
		}
		// MongoDB 密码在 URI 中，需要解析
		return instance.Uri, nil

	case "elasticsearch":
		var instance esmodel.EsInstance
		err := db.Db.Where("id = ?", id).First(&instance).Error
		if err != nil {
			return "", err
		}
		return instance.Password, nil

	default:
		return "", errors.New("不支持的数据库类型: " + dbType)
	}
}
