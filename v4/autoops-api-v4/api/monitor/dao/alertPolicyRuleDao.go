package dao

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
)

type AlertPolicyRuleDao interface {
	Create(rule *model.AlertPolicyRule) error
	Update(rule *model.AlertPolicyRule) error
	Delete(id uint) error
	GetByID(id uint) (*model.AlertPolicyRule, error)
	GetByPolicyID(policyID uint) ([]*model.AlertPolicyRule, error)
	DeleteByPolicyID(policyID uint) error
	UpdateStatus(id uint, status int) error
	UpdateHealthStatus(id uint, healthStatus string) error
}

type AlertPolicyRuleDaoImpl struct{}

func NewAlertPolicyRuleDao() AlertPolicyRuleDao {
	return &AlertPolicyRuleDaoImpl{}
}

// Create 创建策略规则
func (d *AlertPolicyRuleDaoImpl) Create(rule *model.AlertPolicyRule) error {
	return common.GetDB().Create(rule).Error
}

// Update 更新策略规则
func (d *AlertPolicyRuleDaoImpl) Update(rule *model.AlertPolicyRule) error {
	return common.GetDB().Model(&model.AlertPolicyRule{}).Where("id = ?", rule.ID).Updates(rule).Error
}

// Delete 删除策略规则
func (d *AlertPolicyRuleDaoImpl) Delete(id uint) error {
	return common.GetDB().Delete(&model.AlertPolicyRule{}, id).Error
}

// GetByID 根据ID获取策略规则
func (d *AlertPolicyRuleDaoImpl) GetByID(id uint) (*model.AlertPolicyRule, error) {
	var rule model.AlertPolicyRule
	if err := common.GetDB().First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// GetByPolicyID 根据策略ID获取规则列表
func (d *AlertPolicyRuleDaoImpl) GetByPolicyID(policyID uint) ([]*model.AlertPolicyRule, error) {
	var rules []*model.AlertPolicyRule
	if err := common.GetDB().Where("policy_id = ?", policyID).Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// DeleteByPolicyID 删除策略下的所有规则
func (d *AlertPolicyRuleDaoImpl) DeleteByPolicyID(policyID uint) error {
	return common.GetDB().Where("policy_id = ?", policyID).Delete(&model.AlertPolicyRule{}).Error
}

// UpdateStatus 更新策略规则状态
func (d *AlertPolicyRuleDaoImpl) UpdateStatus(id uint, status int) error {
	return common.GetDB().Model(&model.AlertPolicyRule{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateHealthStatus 更新策略规则健康状态
func (d *AlertPolicyRuleDaoImpl) UpdateHealthStatus(id uint, healthStatus string) error {
	return common.GetDB().Model(&model.AlertPolicyRule{}).Where("id = ?", id).Update("health_status", healthStatus).Error
}
