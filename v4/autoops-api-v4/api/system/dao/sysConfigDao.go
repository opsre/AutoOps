// 系统配置数据访问层
package dao

import (
	"time"

	"dodevops-api/api/system/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"
)

// SysConfigDao 系统配置DAO接口
type SysConfigDao interface {
	Create(config *model.SysConfig) error
	GetByID(id uint) (*model.SysConfig, error)
	GetByKey(configKey string) (*model.SysConfig, error)
	Update(config *model.SysConfig) error
	Delete(id uint) error
	GetList(query *model.SysConfigQuery) ([]model.SysConfig, int64, error)
	UpdateStatus(id uint, status int) error
}

// SysConfigDaoImpl 系统配置DAO实现
type SysConfigDaoImpl struct{}

// NewSysConfigDao 创建系统配置DAO实例
func NewSysConfigDao() SysConfigDao {
	return &SysConfigDaoImpl{}
}

// Create 创建系统配置
func (d *SysConfigDaoImpl) Create(config *model.SysConfig) error {
	now := util.HTime{Time: time.Now()}
	config.CreateTime = now
	config.UpdateTime = now
	return Db.Create(config).Error
}

// GetByID 根据ID获取系统配置
func (d *SysConfigDaoImpl) GetByID(id uint) (*model.SysConfig, error) {
	var config model.SysConfig
	err := Db.Where("id = ?", id).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetByKey 根据配置键获取系统配置
func (d *SysConfigDaoImpl) GetByKey(configKey string) (*model.SysConfig, error) {
	var config model.SysConfig
	err := Db.Where("config_key = ?", configKey).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Update 更新系统配置
func (d *SysConfigDaoImpl) Update(config *model.SysConfig) error {
	config.UpdateTime = util.HTime{Time: time.Now()}
	return Db.Model(&model.SysConfig{}).Where("id = ?", config.ID).Updates(config).Error
}

// Delete 删除系统配置
func (d *SysConfigDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.SysConfig{}, id).Error
}

// GetList 获取系统配置列表
func (d *SysConfigDaoImpl) GetList(query *model.SysConfigQuery) ([]model.SysConfig, int64, error) {
	var configs []model.SysConfig
	var total int64

	db := Db.Model(&model.SysConfig{})

	// 条件筛选
	if query.ConfigKey != "" {
		db = db.Where("config_key LIKE ?", "%"+query.ConfigKey+"%")
	}
	if query.ConfigType != "" {
		db = db.Where("config_type = ?", query.ConfigType)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if query.PageNum > 0 && query.PageSize > 0 {
		offset := (query.PageNum - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 按创建时间倒序
	db = db.Order("create_time DESC")

	if err := db.Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

// UpdateStatus 更新状态
func (d *SysConfigDaoImpl) UpdateStatus(id uint, status int) error {
	return Db.Model(&model.SysConfig{}).Where("id = ?", id).Update("status", status).Error
}
