package dao

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
)

type APIBusinessFunctionDao interface {
	Create(function *model.APIBusinessFunction) error
	Update(function *model.APIBusinessFunction) error
	Delete(id uint) error
	GetByID(id uint) (*model.APIBusinessFunction, error)
	GetAll() ([]*model.APIBusinessFunction, error)
	GetByGroupID(groupID uint) ([]*model.APIBusinessFunction, error)
	GetByGroupIDAndName(groupID uint, name string) (*model.APIBusinessFunction, error)
	GetList(groupID *uint, keyword string, page, pageSize int) ([]*model.APIBusinessFunction, int64, error)
	UpdateStatus(id uint, status int) error
}

type APIBusinessFunctionDaoImpl struct{}

func NewAPIBusinessFunctionDao() APIBusinessFunctionDao {
	return &APIBusinessFunctionDaoImpl{}
}

// Create 创建业务功能
func (d *APIBusinessFunctionDaoImpl) Create(function *model.APIBusinessFunction) error {
	return common.GetDB().Create(function).Error
}

// Update 更新业务功能
func (d *APIBusinessFunctionDaoImpl) Update(function *model.APIBusinessFunction) error {
	return common.GetDB().Model(&model.APIBusinessFunction{}).
		Where("id = ?", function.ID).
		Updates(function).Error
}

// Delete 删除业务功能
func (d *APIBusinessFunctionDaoImpl) Delete(id uint) error {
	return common.GetDB().Delete(&model.APIBusinessFunction{}, id).Error
}

// GetByID 根据ID获取业务功能
func (d *APIBusinessFunctionDaoImpl) GetByID(id uint) (*model.APIBusinessFunction, error) {
	var function model.APIBusinessFunction
	err := common.GetDB().Where("id = ?", id).First(&function).Error
	return &function, err
}

// GetAll 获取所有业务功能
func (d *APIBusinessFunctionDaoImpl) GetAll() ([]*model.APIBusinessFunction, error) {
	var functions []*model.APIBusinessFunction
	err := common.GetDB().Order("group_id, sort").Find(&functions).Error
	return functions, err
}

// GetByGroupID 根据业务组ID获取业务功能列表
func (d *APIBusinessFunctionDaoImpl) GetByGroupID(groupID uint) ([]*model.APIBusinessFunction, error) {
	var functions []*model.APIBusinessFunction
	err := common.GetDB().Where("group_id = ?", groupID).Order("sort").Find(&functions).Error
	return functions, err
}

// GetByGroupIDAndName 根据业务组ID和名称查询业务功能
func (d *APIBusinessFunctionDaoImpl) GetByGroupIDAndName(groupID uint, name string) (*model.APIBusinessFunction, error) {
	var function model.APIBusinessFunction
	err := common.GetDB().Where("group_id = ? AND name = ?", groupID, name).First(&function).Error
	if err != nil {
		return nil, err
	}
	return &function, nil
}

// GetList 分页查询业务功能列表
func (d *APIBusinessFunctionDaoImpl) GetList(groupID *uint, keyword string, page, pageSize int) ([]*model.APIBusinessFunction, int64, error) {
	var functions []*model.APIBusinessFunction
	var total int64

	db := common.GetDB().Model(&model.APIBusinessFunction{})

	// 业务组筛选
	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	}

	// 关键字搜索
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 统计总数
	db.Count(&total)

	// 分页
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}

	err := db.Order("group_id, sort").Find(&functions).Error
	return functions, total, err
}

// UpdateStatus 更新状态
func (d *APIBusinessFunctionDaoImpl) UpdateStatus(id uint, status int) error {
	return common.GetDB().Model(&model.APIBusinessFunction{}).
		Where("id = ?", id).
		Update("status", status).Error
}
