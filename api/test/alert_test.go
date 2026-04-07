package test

import (
	"log"
	"testing"
	"time"

	"dodevops-api/api/monitor/model"
	"dodevops-api/pkg/db"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// 使用内存 sqlite 进行单元测试
	database, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(database)
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
	return database
}

func TestAlertDao(t *testing.T) {
	database := setupTestDB()

	// 初始化 Dao
	// 实际开发中可以通过注入将 test 的 db 传给 dao
	// 这里假设我们在测试时手动注册了

	// 为了不修改原有 common.GetDB()，如果在业务中紧耦合了，
	// 最好是使用我们独立的 GORM 方法直接进行测试：

	tpl := &model.PrometheusAlertDB{
		Tpltype: "dingtalk",
		Tpluse:  "Prometheus",
		Tplname: "TestAlert",
		Tpl:     "alert.data",
		Created: time.Now(),
	}

	// Create
	err := database.Create(tpl).Error
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	if tpl.Id == 0 {
		t.Fatalf("Template id is 0 after creation")
	}

	// Get
	var result model.PrometheusAlertDB
	err = database.First(&result, tpl.Id).Error
	if err != nil {
		t.Fatalf("Failed to fetch template: %v", err)
	}
	if result.Tplname != "TestAlert" {
		t.Fatalf("Expected Tplname TestAlert, got %s", result.Tplname)
	}

	// Update
	result.Tplname = "UpdatedAlert"
	err = database.Save(&result).Error
	if err != nil {
		t.Fatalf("Failed to update template: %v", err)
	}

	var updatedResult model.PrometheusAlertDB
	database.First(&updatedResult, tpl.Id)
	if updatedResult.Tplname != "UpdatedAlert" {
		t.Fatalf("Expected Tplname UpdatedAlert, got %s", updatedResult.Tplname)
	}

	// Delete
	err = database.Delete(&model.PrometheusAlertDB{}, tpl.Id).Error
	if err != nil {
		t.Fatalf("Failed to delete template: %v", err)
	}

	var count int64
	database.Model(&model.PrometheusAlertDB{}).Where("id = ?", tpl.Id).Count(&count)
	if count != 0 {
		t.Fatalf("Expected record count to be 0 after delete")
	}
}
