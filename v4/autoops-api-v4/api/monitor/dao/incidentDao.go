package dao

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

type IncidentDao struct {
	db *gorm.DB
}

func NewIncidentDao() IncidentDao {
	return IncidentDao{
		db: common.GetDB(),
	}
}

// CreateIncident 创建故障记录
func (d *IncidentDao) CreateIncident(incident *model.Incident) error {
	return d.db.Create(incident).Error
}

// UpdateIncident 更新故障记录
func (d *IncidentDao) UpdateIncident(id uint, incident *model.Incident) error {
	return d.db.Model(&model.Incident{}).Where("id = ?", id).Updates(incident).Error
}

// DeleteIncident 删除故障记录
func (d *IncidentDao) DeleteIncident(id uint) error {
	return d.db.Delete(&model.Incident{}, id).Error
}

// GetIncidentById 根据ID获取故障记录
func (d *IncidentDao) GetIncidentById(id uint) (model.Incident, error) {
	var incident model.Incident
	err := d.db.Where("id = ?", id).First(&incident).Error
	return incident, err
}

// GetIncidentList 获取故障列表（带分页和条件查询）
func (d *IncidentDao) GetIncidentList(query *model.IncidentQueryDto) (list []model.Incident, total int64, err error) {
	db := d.db.Model(&model.Incident{})

	// 条件过滤
	if query.BusinessLineID > 0 {
		db = db.Where("business_line_id = ?", query.BusinessLineID)
	}
	if query.AlertLevel != "" {
		db = db.Where("alert_level = ?", query.AlertLevel)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.Department != "" {
		db = db.Where("department = ?", query.Department)
	}
	if query.Handler != "" {
		db = db.Where("handler LIKE ?", "%"+query.Handler+"%")
	}
	if query.StartTime != "" {
		db = db.Where("alert_time >= ?", query.StartTime)
	}
	if query.EndTime != "" {
		db = db.Where("alert_time <= ?", query.EndTime)
	}

	// 统计总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 分页
	if query.PageNum > 0 && query.PageSize > 0 {
		offset := (query.PageNum - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 按告警时间倒序排列
	err = db.Order("alert_time DESC").Find(&list).Error
	return
}

// BatchDeleteIncidents 批量删除故障记录
func (d *IncidentDao) BatchDeleteIncidents(ids []uint) error {
	return d.db.Where("id IN ?", ids).Delete(&model.Incident{}).Error
}

// UpdateIncidentStatus 更新故障状态
func (d *IncidentDao) UpdateIncidentStatus(id uint, status int) error {
	return d.db.Model(&model.Incident{}).Where("id = ?", id).Update("status", status).Error
}

// GetIncidentStatistics 获取故障统计信息
func (d *IncidentDao) GetIncidentStatistics() (map[string]interface{}, error) {
	var total int64
	var pendingCount int64
	var processingCount int64
	var completedCount int64

	// 总数
	d.db.Model(&model.Incident{}).Count(&total)

	// 各状态数量
	d.db.Model(&model.Incident{}).Where("status = ?", model.IncidentStatusPending).Count(&pendingCount)
	d.db.Model(&model.Incident{}).Where("status = ?", model.IncidentStatusProcessing).Count(&processingCount)
	d.db.Model(&model.Incident{}).Where("status = ?", model.IncidentStatusCompleted).Count(&completedCount)

	return map[string]interface{}{
		"total":      total,
		"pending":    pendingCount,
		"processing": processingCount,
		"completed":  completedCount,
	}, nil
}
