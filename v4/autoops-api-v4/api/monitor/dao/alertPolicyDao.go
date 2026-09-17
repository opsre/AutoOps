package dao

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
)

type AlertPolicyDao interface {
	Create(policy *model.AlertPolicy) error
	Update(policy *model.AlertPolicy) error
	Delete(id uint) error
	GetByID(id uint) (*model.AlertPolicy, error)
	GetList(groupID *uint, status *int, keyword string, page, pageSize int) ([]*model.AlertPolicy, int64, error)
	UpdateStatus(id uint, status int) error
	UpdateHealthStatus(id uint, healthStatus string) error
	GetAllEnabled() ([]*model.AlertPolicy, error)
	GetByGroupID(groupID uint) ([]*model.AlertPolicy, error)
}

type AlertPolicyDaoImpl struct{}

func NewAlertPolicyDao() AlertPolicyDao {
	return &AlertPolicyDaoImpl{}
}

// Create 创建告警策略
func (d *AlertPolicyDaoImpl) Create(policy *model.AlertPolicy) error {
	return common.GetDB().Create(policy).Error
}

// Update 更新告警策略
func (d *AlertPolicyDaoImpl) Update(policy *model.AlertPolicy) error {
	return common.GetDB().Model(&model.AlertPolicy{}).Where("id = ?", policy.ID).Updates(policy).Error
}

// Delete 删除告警策略
func (d *AlertPolicyDaoImpl) Delete(id uint) error {
	return common.GetDB().Delete(&model.AlertPolicy{}, id).Error
}

// GetByID 根据ID获取告警策略
func (d *AlertPolicyDaoImpl) GetByID(id uint) (*model.AlertPolicy, error) {
	var policy model.AlertPolicy
	if err := common.GetDB().First(&policy, id).Error; err != nil {
		return nil, err
	}
	return &policy, nil
}

// GetList 获取告警策略列表
func (d *AlertPolicyDaoImpl) GetList(groupID *uint, status *int, keyword string, page, pageSize int) ([]*model.AlertPolicy, int64, error) {
	var policies []*model.AlertPolicy
	var total int64

	query := common.GetDB().Model(&model.AlertPolicy{})

	// 按业务组筛选
	if groupID != nil {
		query = query.Where("group_id = ?", *groupID)
	}

	// 按状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 关键字搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&policies).Error; err != nil {
		return nil, 0, err
	}

	return policies, total, nil
}

// UpdateStatus 更新告警策略状态
func (d *AlertPolicyDaoImpl) UpdateStatus(id uint, status int) error {
	return common.GetDB().Model(&model.AlertPolicy{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateHealthStatus 更新告警策略健康状态
func (d *AlertPolicyDaoImpl) UpdateHealthStatus(id uint, healthStatus string) error {
	return common.GetDB().Model(&model.AlertPolicy{}).Where("id = ?", id).Update("health_status", healthStatus).Error
}

// GetAllEnabled 获取所有启用的告警策略
func (d *AlertPolicyDaoImpl) GetAllEnabled() ([]*model.AlertPolicy, error) {
	var policies []*model.AlertPolicy
	if err := common.GetDB().Where("status = ?", 1).Find(&policies).Error; err != nil {
		return nil, err
	}
	return policies, nil
}

// GetByGroupID 根据业务组ID获取告警策略
func (d *AlertPolicyDaoImpl) GetByGroupID(groupID uint) ([]*model.AlertPolicy, error) {
	var policies []*model.AlertPolicy
	if err := common.GetDB().Where("group_id = ?", groupID).Find(&policies).Error; err != nil {
		return nil, err
	}
	return policies, nil
}
