package dao

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
	"time"
)

type APIMonitorResultDao interface {
	Create(result *model.APIMonitorResult) error
	GetLatest(monitorID uint) (*model.APIMonitorResult, error)
	GetHistory(monitorID uint, startTime, endTime time.Time) ([]*model.APIMonitorResult, error)
	DeleteOldRecords(days int) error
	GetHealthRate24h(monitorID uint) (float64, error)
	GetHealthRateBatch24h(monitorIDs []uint) (map[uint]float64, error)
}

type APIMonitorResultDaoImpl struct{}

func NewAPIMonitorResultDao() APIMonitorResultDao {
	return &APIMonitorResultDaoImpl{}
}

// Create 创建检测结果
func (d *APIMonitorResultDaoImpl) Create(result *model.APIMonitorResult) error {
	return common.GetDB().Create(result).Error
}

// GetLatest 获取最新检测结果（优化：只选择必要字段，不加载response_body）
func (d *APIMonitorResultDaoImpl) GetLatest(monitorID uint) (*model.APIMonitorResult, error) {
	var result model.APIMonitorResult
	err := common.GetDB().
		Select("id, monitor_id, check_time, success, status_code, response_time, error_message").
		Where("monitor_id = ?", monitorID).
		Order("check_time DESC").
		First(&result).Error
	return &result, err
}

// GetHistory 获取历史检测结果
func (d *APIMonitorResultDaoImpl) GetHistory(monitorID uint, startTime, endTime time.Time) ([]*model.APIMonitorResult, error) {
	var results []*model.APIMonitorResult
	err := common.GetDB().Where("monitor_id = ? AND check_time >= ? AND check_time <= ?",
		monitorID, startTime, endTime).
		Order("check_time DESC").
		Find(&results).Error
	return results, err
}

// DeleteOldRecords 删除旧记录（保留最近N天）
func (d *APIMonitorResultDaoImpl) DeleteOldRecords(days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	return common.GetDB().Where("check_time < ?", cutoffTime).
		Delete(&model.APIMonitorResult{}).Error
}

// GetHealthRate24h 获取过去24小时的健康率
func (d *APIMonitorResultDaoImpl) GetHealthRate24h(monitorID uint) (float64, error) {
	startTime := time.Now().Add(-24 * time.Hour)

	var total int64
	var successCount int64

	db := common.GetDB().Model(&model.APIMonitorResult{}).
		Where("monitor_id = ? AND check_time >= ?", monitorID, startTime)

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	// 如果没有检测记录，返回0
	if total == 0 {
		return 0, nil
	}

	// 统计成功次数
	if err := db.Where("success = ?", true).Count(&successCount).Error; err != nil {
		return 0, err
	}

	// 计算健康率
	healthRate := float64(successCount) / float64(total) * 100
	return healthRate, nil
}

// GetHealthRateBatch24h 批量获取过去24小时的健康率（优化版）
func (d *APIMonitorResultDaoImpl) GetHealthRateBatch24h(monitorIDs []uint) (map[uint]float64, error) {
	if len(monitorIDs) == 0 {
		return make(map[uint]float64), nil
	}

	startTime := time.Now().Add(-24 * time.Hour)

	type HealthRateResult struct {
		MonitorID    uint
		Total        int64
		SuccessCount int64
	}

	var results []HealthRateResult

	// 优化：
	// 1. 使用 SUM(success) 替代 SUM(CASE WHEN)，因为success是tinyint(0或1)
	// 2. 使用 FORCE INDEX 提示数据库使用正确的复合索引（需要先创建索引）
	// 3. 只选择必要的字段，避免全表扫描
	err := common.GetDB().Model(&model.APIMonitorResult{}).
		Select(`
			monitor_id,
			COUNT(*) as total,
			SUM(success) as success_count
		`).
		Where("monitor_id IN ? AND check_time >= ?", monitorIDs, startTime).
		Group("monitor_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 构建结果映射
	healthRates := make(map[uint]float64, len(monitorIDs))
	for _, result := range results {
		if result.Total > 0 {
			healthRates[result.MonitorID] = float64(result.SuccessCount) / float64(result.Total) * 100
		} else {
			healthRates[result.MonitorID] = 0
		}
	}

	// 对于没有检测记录的监控项，设置健康率为0
	for _, monitorID := range monitorIDs {
		if _, exists := healthRates[monitorID]; !exists {
			healthRates[monitorID] = 0
		}
	}

	return healthRates, nil
}
