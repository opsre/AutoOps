package dao

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
)

type AlertRuleTemplateDao interface {
	Create(template *model.AlertRuleTemplate) error
	Update(template *model.AlertRuleTemplate) error
	Delete(id uint) error
	GetByID(id uint) (*model.AlertRuleTemplate, error)
	GetByIDs(ids []uint) ([]*model.AlertRuleTemplate, error)
	GetList(category string, status *int, keyword string, page, pageSize int) ([]*model.AlertRuleTemplate, int64, error)
	UpdateStatus(id uint, status int) error
	GetAllEnabled() ([]*model.AlertRuleTemplate, error)
}

type AlertRuleTemplateDaoImpl struct{}

func NewAlertRuleTemplateDao() AlertRuleTemplateDao {
	return &AlertRuleTemplateDaoImpl{}
}

// Create 创建规则模板
func (d *AlertRuleTemplateDaoImpl) Create(template *model.AlertRuleTemplate) error {
	return common.GetDB().Create(template).Error
}

// Update 更新规则模板
func (d *AlertRuleTemplateDaoImpl) Update(template *model.AlertRuleTemplate) error {
	return common.GetDB().Model(&model.AlertRuleTemplate{}).Where("id = ?", template.ID).Updates(template).Error
}

// Delete 删除规则模板
func (d *AlertRuleTemplateDaoImpl) Delete(id uint) error {
	return common.GetDB().Delete(&model.AlertRuleTemplate{}, id).Error
}

// GetByID 根据ID获取规则模板
func (d *AlertRuleTemplateDaoImpl) GetByID(id uint) (*model.AlertRuleTemplate, error) {
	var template model.AlertRuleTemplate
	if err := common.GetDB().First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByIDs 根据ID列表获取规则模板
func (d *AlertRuleTemplateDaoImpl) GetByIDs(ids []uint) ([]*model.AlertRuleTemplate, error) {
	var templates []*model.AlertRuleTemplate
	if err := common.GetDB().Where("id IN ?", ids).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// GetList 获取规则模板列表
func (d *AlertRuleTemplateDaoImpl) GetList(category string, status *int, keyword string, page, pageSize int) ([]*model.AlertRuleTemplate, int64, error) {
	var templates []*model.AlertRuleTemplate
	var total int64

	query := common.GetDB().Model(&model.AlertRuleTemplate{})

	// 按分类筛选
	if category != "" {
		query = query.Where("category = ?", category)
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
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// UpdateStatus 更新规则模板状态
func (d *AlertRuleTemplateDaoImpl) UpdateStatus(id uint, status int) error {
	return common.GetDB().Model(&model.AlertRuleTemplate{}).Where("id = ?", id).Update("status", status).Error
}

// GetAllEnabled 获取所有启用的规则模板
func (d *AlertRuleTemplateDaoImpl) GetAllEnabled() ([]*model.AlertRuleTemplate, error) {
	var templates []*model.AlertRuleTemplate
	if err := common.GetDB().Where("status = ?", 1).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}
