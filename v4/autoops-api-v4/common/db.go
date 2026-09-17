package common

import (
	"dodevops-api/pkg/db"

	"gorm.io/gorm"
)

// GetDB 获取数据库连接（使用pkg/db中的连接）
func GetDB() *gorm.DB {
	if db.Db == nil {
		panic("Database connection is not initialized")
	}

	return db.Db
}
