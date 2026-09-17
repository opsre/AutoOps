// 容量预测和告警数据访问层
package dao

import (
	"time"

	"dodevops-api/api/agent/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"
)

// PredictionDao 容量预测DAO接口
type PredictionDao interface {
	Create(prediction *model.CapacityPrediction) error
	GetByID(id uint) (*model.CapacityPrediction, error)
	Update(prediction *model.CapacityPrediction) error
	Delete(id uint) error
	GetList(query *model.CapacityPredictionQuery) ([]model.CapacityPrediction, int64, error)
	GetEnabledPredictions() ([]model.CapacityPrediction, error)
	GetWarningPredictions() ([]model.CapacityPrediction, error)
}

// PredictionDaoImpl 容量预测DAO实现
type PredictionDaoImpl struct{}

// NewPredictionDao 创建容量预测DAO实例
func NewPredictionDao() PredictionDao {
	return &PredictionDaoImpl{}
}

// Create 创建容量预测
func (d *PredictionDaoImpl) Create(prediction *model.CapacityPrediction) error {
	now := util.HTime{Time: time.Now()}
	prediction.CreateTime = now
	prediction.UpdateTime = now
	return Db.Create(prediction).Error
}

// GetByID 根据ID获取容量预测
func (d *PredictionDaoImpl) GetByID(id uint) (*model.CapacityPrediction, error) {
	var prediction model.CapacityPrediction
	err := Db.Where("id = ?", id).First(&prediction).Error
	if err != nil {
		return nil, err
	}
	return &prediction, nil
}

// Update 更新容量预测
func (d *PredictionDaoImpl) Update(prediction *model.CapacityPrediction) error {
	prediction.UpdateTime = util.HTime{Time: time.Now()}
	return Db.Model(&model.CapacityPrediction{}).Where("id = ?", prediction.ID).Updates(prediction).Error
}

// Delete 删除容量预测
func (d *PredictionDaoImpl) Delete(id uint) error {
	// 同时删除关联的历史记录
	Db.Where("prediction_id = ?", id).Delete(&model.PredictionHistory{})
	return Db.Delete(&model.CapacityPrediction{}, id).Error
}

// GetList 获取容量预测列表
func (d *PredictionDaoImpl) GetList(query *model.CapacityPredictionQuery) ([]model.CapacityPrediction, int64, error) {
	var predictions []model.CapacityPrediction
	var total int64

	db := Db.Model(&model.CapacityPrediction{})

	// 条件筛选
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.TargetType != "" {
		db = db.Where("target_type = ?", query.TargetType)
	}
	if query.TargetID > 0 {
		db = db.Where("target_id = ?", query.TargetID)
	}
	if query.MetricType != "" {
		db = db.Where("metric_type = ?", query.MetricType)
	}
	if query.RiskLevel >= 0 {
		db = db.Where("risk_level = ?", query.RiskLevel)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.Enabled > 0 {
		db = db.Where("enabled = ?", query.Enabled)
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

	// 按风险等级倒序、剩余天数正序
	db = db.Order("risk_level DESC, days_left ASC")

	if err := db.Find(&predictions).Error; err != nil {
		return nil, 0, err
	}

	return predictions, total, nil
}

// GetEnabledPredictions 获取所有启用的容量预测
func (d *PredictionDaoImpl) GetEnabledPredictions() ([]model.CapacityPrediction, error) {
	var predictions []model.CapacityPrediction
	err := Db.Where("enabled = ?", 1).Find(&predictions).Error
	return predictions, err
}

// GetWarningPredictions 获取有预警的容量预测
func (d *PredictionDaoImpl) GetWarningPredictions() ([]model.CapacityPrediction, error) {
	var predictions []model.CapacityPrediction
	err := Db.Where("enabled = ? AND risk_level >= ?", 1, 2). // 风险等级>=预警
									Order("risk_level DESC, days_left ASC").
									Find(&predictions).Error
	return predictions, err
}

// PredictionHistoryDao 预测历史DAO接口
type PredictionHistoryDao interface {
	Create(history *model.PredictionHistory) error
	GetByPredictionID(predictionID uint, limit int) ([]model.PredictionHistory, error)
	Delete(id uint) error
	DeleteOldRecords(predictionID uint, days int) error
}

// PredictionHistoryDaoImpl 预测历史DAO实现
type PredictionHistoryDaoImpl struct{}

// NewPredictionHistoryDao 创建预测历史DAO实例
func NewPredictionHistoryDao() PredictionHistoryDao {
	return &PredictionHistoryDaoImpl{}
}

// Create 创建预测历史
func (d *PredictionHistoryDaoImpl) Create(history *model.PredictionHistory) error {
	history.CreateTime = util.HTime{Time: time.Now()}
	return Db.Create(history).Error
}

// GetByPredictionID 根据预测ID获取历史记录
func (d *PredictionHistoryDaoImpl) GetByPredictionID(predictionID uint, limit int) ([]model.PredictionHistory, error) {
	var histories []model.PredictionHistory
	db := Db.Where("prediction_id = ?", predictionID).
		Order("predict_time DESC")

	if limit > 0 {
		db = db.Limit(limit)
	}

	err := db.Find(&histories).Error
	return histories, err
}

// Delete 删除预测历史
func (d *PredictionHistoryDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.PredictionHistory{}, id).Error
}

// DeleteOldRecords 删除旧的历史记录
func (d *PredictionHistoryDaoImpl) DeleteOldRecords(predictionID uint, days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	return Db.Where("prediction_id = ? AND create_time < ?", predictionID, cutoffTime).
		Delete(&model.PredictionHistory{}).Error
}

// AlertRuleDao 告警规则DAO接口
type AlertRuleDao interface {
	Create(rule *model.AlertRule) error
	GetByID(id uint) (*model.AlertRule, error)
	Update(rule *model.AlertRule) error
	Delete(id uint) error
	GetList(query *model.AlertRuleQuery) ([]model.AlertRule, int64, error)
	GetEnabledRules() ([]model.AlertRule, error)
}

// AlertRuleDaoImpl 告警规则DAO实现
type AlertRuleDaoImpl struct{}

// NewAlertRuleDao 创建告警规则DAO实例
func NewAlertRuleDao() AlertRuleDao {
	return &AlertRuleDaoImpl{}
}

// Create 创建告警规则
func (d *AlertRuleDaoImpl) Create(rule *model.AlertRule) error {
	now := util.HTime{Time: time.Now()}
	rule.CreateTime = now
	rule.UpdateTime = now
	return Db.Create(rule).Error
}

// GetByID 根据ID获取告警规则
func (d *AlertRuleDaoImpl) GetByID(id uint) (*model.AlertRule, error) {
	var rule model.AlertRule
	err := Db.Where("id = ?", id).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// Update 更新告警规则
func (d *AlertRuleDaoImpl) Update(rule *model.AlertRule) error {
	rule.UpdateTime = util.HTime{Time: time.Now()}
	return Db.Model(&model.AlertRule{}).Where("id = ?", rule.ID).Updates(rule).Error
}

// Delete 删除告警规则
func (d *AlertRuleDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.AlertRule{}, id).Error
}

// GetList 获取告警规则列表
func (d *AlertRuleDaoImpl) GetList(query *model.AlertRuleQuery) ([]model.AlertRule, int64, error) {
	var rules []model.AlertRule
	var total int64

	db := Db.Model(&model.AlertRule{})

	// 条件筛选
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.TargetType != "" {
		db = db.Where("target_type = ?", query.TargetType)
	}
	if query.MetricType != "" {
		db = db.Where("metric_type = ?", query.MetricType)
	}
	if query.Severity > 0 {
		db = db.Where("severity = ?", query.Severity)
	}
	if query.Enabled > 0 {
		db = db.Where("enabled = ?", query.Enabled)
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

	if err := db.Find(&rules).Error; err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

// GetEnabledRules 获取所有启用的告警规则
func (d *AlertRuleDaoImpl) GetEnabledRules() ([]model.AlertRule, error) {
	var rules []model.AlertRule
	err := Db.Where("enabled = ?", 1).Find(&rules).Error
	return rules, err
}

// AlertEventDao 告警事件DAO接口
type AlertEventDao interface {
	Create(event *model.AlertEvent) error
	GetByID(id uint) (*model.AlertEvent, error)
	Update(event *model.AlertEvent) error
	Delete(id uint) error
	GetList(query *model.AlertEventQuery) ([]model.AlertEvent, int64, error)
	GetActiveEvents() ([]model.AlertEvent, error)
	GetStatistics() (map[string]interface{}, error)
}

// AlertEventDaoImpl 告警事件DAO实现
type AlertEventDaoImpl struct{}

// NewAlertEventDao 创建告警事件DAO实例
func NewAlertEventDao() AlertEventDao {
	return &AlertEventDaoImpl{}
}

// Create 创建告警事件
func (d *AlertEventDaoImpl) Create(event *model.AlertEvent) error {
	now := util.HTime{Time: time.Now()}
	event.CreateTime = now
	event.UpdateTime = now
	return Db.Create(event).Error
}

// GetByID 根据ID获取告警事件
func (d *AlertEventDaoImpl) GetByID(id uint) (*model.AlertEvent, error) {
	var event model.AlertEvent
	err := Db.Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Update 更新告警事件
func (d *AlertEventDaoImpl) Update(event *model.AlertEvent) error {
	event.UpdateTime = util.HTime{Time: time.Now()}
	return Db.Model(&model.AlertEvent{}).Where("id = ?", event.ID).Updates(event).Error
}

// Delete 删除告警事件
func (d *AlertEventDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.AlertEvent{}, id).Error
}

// GetList 获取告警事件列表
func (d *AlertEventDaoImpl) GetList(query *model.AlertEventQuery) ([]model.AlertEvent, int64, error) {
	var events []model.AlertEvent
	var total int64

	db := Db.Model(&model.AlertEvent{})

	// 条件筛选
	if query.RuleID > 0 {
		db = db.Where("rule_id = ?", query.RuleID)
	}
	if query.HostID > 0 {
		db = db.Where("host_id = ?", query.HostID)
	}
	if query.HostName != "" {
		db = db.Where("host_name LIKE ?", "%"+query.HostName+"%")
	}
	if query.MetricType != "" {
		db = db.Where("metric_type = ?", query.MetricType)
	}
	if query.Severity > 0 {
		db = db.Where("severity = ?", query.Severity)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.StartTime != "" {
		db = db.Where("trigger_time >= ?", query.StartTime)
	}
	if query.EndTime != "" {
		db = db.Where("trigger_time <= ?", query.EndTime)
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

	// 按触发时间倒序
	db = db.Order("trigger_time DESC")

	if err := db.Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetActiveEvents 获取活跃的告警事件
func (d *AlertEventDaoImpl) GetActiveEvents() ([]model.AlertEvent, error) {
	var events []model.AlertEvent
	err := Db.Where("status = ?", 1). // 状态=触发
						Order("severity DESC, trigger_time ASC").
						Find(&events).Error
	return events, err
}

// GetStatistics 获取告警统计
func (d *AlertEventDaoImpl) GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 按状态统计
	var statusStats []struct {
		Status int
		Count  int64
	}
	err := Db.Model(&model.AlertEvent{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusStats).Error
	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]int64)
	for _, stat := range statusStats {
		switch stat.Status {
		case 1:
			statusMap["triggered"] = stat.Count
		case 2:
			statusMap["recovered"] = stat.Count
		case 3:
			statusMap["confirmed"] = stat.Count
		case 4:
			statusMap["resolved"] = stat.Count
		}
	}
	stats["by_status"] = statusMap

	// 按严重程度统计
	var severityStats []struct {
		Severity int
		Count    int64
	}
	err = Db.Model(&model.AlertEvent{}).
		Select("severity, count(*) as count").
		Group("severity").
		Scan(&severityStats).Error
	if err != nil {
		return nil, err
	}

	severityMap := make(map[string]int64)
	for _, stat := range severityStats {
		switch stat.Severity {
		case 1:
			severityMap["low"] = stat.Count
		case 2:
			severityMap["medium"] = stat.Count
		case 3:
			severityMap["high"] = stat.Count
		case 4:
			severityMap["critical"] = stat.Count
		}
	}
	stats["by_severity"] = severityMap

	// 今日告警数
	var todayCount int64
	Db.Model(&model.AlertEvent{}).
		Where("DATE(trigger_time) = CURDATE()").
		Count(&todayCount)
	stats["today_count"] = todayCount

	// 活跃告警数
	var activeCount int64
	Db.Model(&model.AlertEvent{}).
		Where("status = ?", 1).
		Count(&activeCount)
	stats["active_count"] = activeCount

	return stats, nil
}
