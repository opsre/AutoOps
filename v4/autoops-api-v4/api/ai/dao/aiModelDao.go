// AI模型数据访问层
package dao

import (
	"time"

	"dodevops-api/api/ai/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"
)

// AiModelDao AI模型DAO接口
type AiModelDao interface {
	Create(aiModel *model.AiModel) error
	GetByID(id uint) (*model.AiModel, error)
	GetByName(name string) (*model.AiModel, error)
	Update(aiModel *model.AiModel) error
	Delete(id uint) error
	GetList(query *model.AiModelQuery) ([]model.AiModel, int64, error)
	UpdateStatus(id uint, status int) error
	GetEnabledModels() ([]model.AiModel, error)
}

// AiModelDaoImpl AI模型DAO实现
type AiModelDaoImpl struct{}

// NewAiModelDao 创建AI模型DAO实例
func NewAiModelDao() AiModelDao {
	return &AiModelDaoImpl{}
}

// Create 创建AI模型
func (d *AiModelDaoImpl) Create(aiModel *model.AiModel) error {
	now := util.HTime{Time: time.Now()}
	aiModel.CreateTime = now
	aiModel.UpdateTime = now
	return Db.Create(aiModel).Error
}

// GetByID 根据ID获取AI模型
func (d *AiModelDaoImpl) GetByID(id uint) (*model.AiModel, error) {
	var aiModel model.AiModel
	err := Db.Where("id = ?", id).First(&aiModel).Error
	if err != nil {
		return nil, err
	}
	return &aiModel, nil
}

// GetByName 根据名称获取AI模型
func (d *AiModelDaoImpl) GetByName(name string) (*model.AiModel, error) {
	var aiModel model.AiModel
	err := Db.Where("name = ?", name).First(&aiModel).Error
	if err != nil {
		return nil, err
	}
	return &aiModel, nil
}

// Update 更新AI模型
func (d *AiModelDaoImpl) Update(aiModel *model.AiModel) error {
	aiModel.UpdateTime = util.HTime{Time: time.Now()}
	return Db.Model(&model.AiModel{}).Where("id = ?", aiModel.ID).Updates(aiModel).Error
}

// Delete 删除AI模型
func (d *AiModelDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.AiModel{}, id).Error
}

// GetList 获取AI模型列表
func (d *AiModelDaoImpl) GetList(query *model.AiModelQuery) ([]model.AiModel, int64, error) {
	var aiModels []model.AiModel
	var total int64

	db := Db.Model(&model.AiModel{})

	// 条件筛选
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
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

	if err := db.Find(&aiModels).Error; err != nil {
		return nil, 0, err
	}

	return aiModels, total, nil
}

// UpdateStatus 更新状态
func (d *AiModelDaoImpl) UpdateStatus(id uint, status int) error {
	return Db.Model(&model.AiModel{}).Where("id = ?", id).Update("status", status).Error
}

// GetEnabledModels 获取所有启用的模型
func (d *AiModelDaoImpl) GetEnabledModels() ([]model.AiModel, error) {
	var aiModels []model.AiModel
	err := Db.Where("status = ?", 1).Order("create_time DESC").Find(&aiModels).Error
	if err != nil {
		return nil, err
	}
	return aiModels, nil
}
