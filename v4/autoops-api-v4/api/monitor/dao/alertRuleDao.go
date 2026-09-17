package dao

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

type AlertRuleDao interface {
	Create(rule *model.AlertRule) error
	Update(rule *model.AlertRule) error
	Delete(id uint) error
	GetByID(id uint) (*model.AlertRule, error)
	GetList(status *int, groupID *uint, page, pageSize int) ([]*model.AlertRule, int64, error)
	GetActiveRules() ([]*model.AlertRule, error)
	UpdateStatus(id uint, status int) error
}

type AlertRuleDaoImpl struct {
	db *gorm.DB
}

func NewAlertRuleDao() AlertRuleDao {
	return &AlertRuleDaoImpl{
		db: common.GetDB(),
	}
}

func (d *AlertRuleDaoImpl) Create(rule *model.AlertRule) error {
	return d.db.Create(rule).Error
}

func (d *AlertRuleDaoImpl) Update(rule *model.AlertRule) error {
	return d.db.Save(rule).Error
}

func (d *AlertRuleDaoImpl) Delete(id uint) error {
	return d.db.Delete(&model.AlertRule{}, id).Error
}

func (d *AlertRuleDaoImpl) GetByID(id uint) (*model.AlertRule, error) {
	var rule model.AlertRule
	err := d.db.First(&rule, id).Error
	return &rule, err
}

func (d *AlertRuleDaoImpl) GetList(status *int, groupID *uint, page, pageSize int) ([]*model.AlertRule, int64, error) {
	var rules []*model.AlertRule
	var total int64

	query := d.db.Model(&model.AlertRule{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if groupID != nil {
		query = query.Where("group_id = ?", *groupID)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&rules).Error

	return rules, total, err
}

func (d *AlertRuleDaoImpl) GetActiveRules() ([]*model.AlertRule, error) {
	var rules []*model.AlertRule
	err := d.db.Where("status = ?", 1).Find(&rules).Error
	return rules, err
}

func (d *AlertRuleDaoImpl) UpdateStatus(id uint, status int) error {
	return d.db.Model(&model.AlertRule{}).Where("id = ?", id).Update("status", status).Error
}
