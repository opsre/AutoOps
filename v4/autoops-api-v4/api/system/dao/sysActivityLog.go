// 系统动态日志 数据层
package dao

import (
	"dodevops-api/api/system/model"
	. "dodevops-api/pkg/db"
)

type SysActivityLogDao struct{}

func NewSysActivityLogDao() *SysActivityLogDao {
	return &SysActivityLogDao{}
}

// Create 创建系统动态日志
func (d *SysActivityLogDao) Create(log *model.SysActivityLog) error {
	return Db.Create(log).Error
}

// GetByID 根据ID查询
func (d *SysActivityLogDao) GetByID(id uint) (*model.SysActivityLog, error) {
	var log model.SysActivityLog
	if err := Db.Where("id = ?", id).First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// List 分页查询动态日志列表
func (d *SysActivityLogDao) List(activityType int, status int, page, pageSize int) ([]model.SysActivityLog, int64, error) {
	var logs []model.SysActivityLog
	var total int64

	query := Db.Model(&model.SysActivityLog{})

	// 筛选条件
	if activityType > 0 {
		query = query.Where("activity_type = ?", activityType)
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetLatest 获取最新的N条动态日志
func (d *SysActivityLogDao) GetLatest(limit int) ([]model.SysActivityLog, error) {
	var logs []model.SysActivityLog
	if err := Db.Order("create_time DESC").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// GetByRelatedID 根据关联ID查询
func (d *SysActivityLogDao) GetByRelatedID(activityType int, relatedID uint) ([]model.SysActivityLog, error) {
	var logs []model.SysActivityLog
	query := Db.Where("related_id = ?", relatedID)
	if activityType > 0 {
		query = query.Where("activity_type = ?", activityType)
	}
	if err := query.Order("create_time DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// DeleteByID 根据ID删除
func (d *SysActivityLogDao) DeleteByID(id uint) error {
	return Db.Delete(&model.SysActivityLog{}, id).Error
}

// DeleteBatch 批量删除
func (d *SysActivityLogDao) DeleteBatch(ids []uint) error {
	return Db.Where("id IN ?", ids).Delete(&model.SysActivityLog{}).Error
}

// Clean 清空动态日志（保留最近N天）
func (d *SysActivityLogDao) Clean(days int) error {
	return Db.Exec("DELETE FROM sys_activity_log WHERE create_time < DATE_SUB(NOW(), INTERVAL ? DAY)", days).Error
}
