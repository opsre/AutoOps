package dao

import (
	"fmt"
	"time"

	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
	"dodevops-api/common/util"

	"gorm.io/gorm"
)

type AlertRecordDao interface {
	Create(record *model.AlertRecord) error
	Update(record *model.AlertRecord) error
	UpdateAlertMetrics(id uint, currentValue float64, duration int, targetName string) error
	GetByID(id uint) (*model.AlertRecord, error)
	GetByFingerprint(fingerprint string) (*model.AlertRecord, error)
	GetActiveRecord(ruleID, targetID uint) (*model.AlertRecord, error)
	GetList(status, severity, ip string, ruleID *uint, groupID *uint, startTime, endTime *time.Time, page, pageSize int) ([]*model.AlertRecord, int64, error)
	GetLastNotifyTime(ruleID, targetID uint) (time.Time, error)
	UpdateNotifyInfo(id uint, notifyTime time.Time, notifyCount int, notifyError string) error
	UpdateHandler(id uint, handler string, handleNote string) error
	UpdateStatus(id uint, status string) error
	ClearFingerprintForResolved(fingerprint string) error
	Delete(id uint) error
	BatchDelete(ids []uint) error
	DeleteAll() error
}

type AlertRecordDaoImpl struct {
	db *gorm.DB
}

func NewAlertRecordDao() AlertRecordDao {
	return &AlertRecordDaoImpl{
		db: common.GetDB(),
	}
}

func (d *AlertRecordDaoImpl) Create(record *model.AlertRecord) error {
	return d.db.Create(record).Error
}

func (d *AlertRecordDaoImpl) Update(record *model.AlertRecord) error {
	return d.db.Save(record).Error
}

func (d *AlertRecordDaoImpl) GetByID(id uint) (*model.AlertRecord, error) {
	var record model.AlertRecord
	err := d.db.First(&record, id).Error
	return &record, err
}

func (d *AlertRecordDaoImpl) GetByFingerprint(fingerprint string) (*model.AlertRecord, error) {
	var record model.AlertRecord
	err := d.db.Where("fingerprint = ? AND status = ?", fingerprint, "firing").First(&record).Error
	return &record, err
}

func (d *AlertRecordDaoImpl) GetActiveRecord(ruleID, targetID uint) (*model.AlertRecord, error) {
	var record model.AlertRecord
	err := d.db.Where("rule_id = ? AND target_id = ? AND status = ?", ruleID, targetID, "firing").
		Order("id DESC").First(&record).Error
	return &record, err
}

func (d *AlertRecordDaoImpl) GetList(status, severity, ip string, ruleID *uint, groupID *uint, startTime, endTime *time.Time, page, pageSize int) ([]*model.AlertRecord, int64, error) {
	var records []*model.AlertRecord
	var total int64

	query := d.db.Model(&model.AlertRecord{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if ruleID != nil {
		query = query.Where("rule_id = ?", *ruleID)
	}
	if groupID != nil {
		query = query.Where("group_id = ?", *groupID)
	}
	if ip != "" {
		query = query.Where(
			"(target_name = ? OR target_name LIKE ? OR labels LIKE ? OR labels LIKE ?)",
			ip,
			"%("+ip+")%",
			"%\"instance\":\""+ip+"\"%",
			"%\"job\":\""+ip+"\"%",
		)
	}
	if startTime != nil {
		query = query.Where("start_time >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("start_time <= ?", *endTime)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&records).Error

	return records, total, err
}

func (d *AlertRecordDaoImpl) GetLastNotifyTime(ruleID, targetID uint) (time.Time, error) {
	var record model.AlertRecord
	err := d.db.Where("rule_id = ? AND target_id = ? AND notify_sent = ?", ruleID, targetID, true).
		Order("notify_time DESC").First(&record).Error

	if err != nil {
		return time.Time{}, err
	}

	if record.NotifyTime == nil {
		return time.Time{}, nil
	}

	return record.NotifyTime.Time, nil
}

func (d *AlertRecordDaoImpl) UpdateNotifyInfo(id uint, notifyTime time.Time, notifyCount int, notifyError string) error {
	updates := map[string]interface{}{
		"notify_sent":  true,
		"notify_time":  util.HTime{Time: notifyTime},
		"notify_count": notifyCount,
	}
	if notifyError != "" {
		updates["notify_error"] = notifyError
	}
	return d.db.Model(&model.AlertRecord{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateAlertMetrics 只更新告警指标相关字段（不影响通知状态）
func (d *AlertRecordDaoImpl) UpdateAlertMetrics(id uint, currentValue float64, duration int, targetName string) error {
	updates := map[string]interface{}{
		"current_value": currentValue,
		"duration":      duration,
		"target_name":   targetName,
	}
	return d.db.Model(&model.AlertRecord{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateHandler 更新告警记录处理人
func (d *AlertRecordDaoImpl) UpdateHandler(id uint, handler string, handleNote string) error {
	now := util.HTime{Time: time.Now()}
	updates := map[string]interface{}{
		"handler":     handler,
		"handle_time": now,
		"handle_note": handleNote,
		"status":      "acknowledged",
	}
	return d.db.Model(&model.AlertRecord{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateStatus 更新告警记录状态
func (d *AlertRecordDaoImpl) UpdateStatus(id uint, status string) error {
	return d.db.Model(&model.AlertRecord{}).Where("id = ?", id).Update("status", status).Error
}

// ClearFingerprintForResolved 清空已恢复告警的fingerprint
// 当已恢复的告警重新触发时，需要清空旧记录的fingerprint以便创建新记录
// 注意: fingerprint有唯一索引，不能设置为空字符串，而是添加时间戳后缀使其唯一
func (d *AlertRecordDaoImpl) ClearFingerprintForResolved(fingerprint string) error {
	// 使用时间戳纳秒作为后缀，确保每次更新都是唯一的
	uniqueSuffix := fmt.Sprintf("_resolved_%d", time.Now().UnixNano())
	return d.db.Model(&model.AlertRecord{}).
		Where("fingerprint = ? AND status != ?", fingerprint, "firing").
		Update("fingerprint", fingerprint+uniqueSuffix).Error
}

// Delete 删除单条告警记录
func (d *AlertRecordDaoImpl) Delete(id uint) error {
	return d.db.Delete(&model.AlertRecord{}, id).Error
}

// BatchDelete 批量删除告警记录
func (d *AlertRecordDaoImpl) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.Delete(&model.AlertRecord{}, ids).Error
}

// DeleteAll 清空所有告警记录
func (d *AlertRecordDaoImpl) DeleteAll() error {
	return d.db.Where("1 = 1").Delete(&model.AlertRecord{}).Error
}
