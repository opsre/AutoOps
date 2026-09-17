package dao

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

type AlertStatsDao struct {
	db *gorm.DB
}

func NewAlertStatsDao() AlertStatsDao {
	return AlertStatsDao{
		db: common.GetDB(),
	}
}

// CreateAlertSource 创建告警源
func (d *AlertStatsDao) CreateAlertSource(source *model.AlertSource) error {
	return d.db.Create(source).Error
}

// UpdateAlertSource 更新告警源
func (d *AlertStatsDao) UpdateAlertSource(id uint, source *model.AlertSource) error {
	return d.db.Model(&model.AlertSource{}).Where("id = ?", id).Updates(source).Error
}

// DeleteAlertSource 删除告警源
func (d *AlertStatsDao) DeleteAlertSource(id uint) error {
	return d.db.Delete(&model.AlertSource{}, id).Error
}

// GetAlertSourceById 根据ID获取告警源
func (d *AlertStatsDao) GetAlertSourceById(id uint) (model.AlertSource, error) {
	var source model.AlertSource
	err := d.db.Where("id = ?", id).First(&source).Error
	return source, err
}

// GetAlertSourceList 获取告警源列表（带分页和条件查询）
func (d *AlertStatsDao) GetAlertSourceList(query *model.AlertSourceQueryDto) (list []model.AlertSource, total int64, err error) {
	db := d.db.Model(&model.AlertSource{})

	// 条件过滤
	if query.Type > 0 {
		db = db.Where("type = ?", query.Type)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
	}

	// 统计总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	println("DEBUG: Count total=", total)

	// 分页
	if query.PageNum > 0 && query.PageSize > 0 {
		offset := (query.PageNum - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 按创建时间倒序排列
	err = db.Order("create_time DESC").Find(&list).Error
	println("DEBUG: After Find, list length=", len(list))

	// 调试日志
	if err == nil {
		for i, item := range list {
			println("DEBUG: item", i, "id=", item.ID, "name=", item.Name, "type=", item.Type)
		}
	}

	return
}

// GetAllAlertSources 获取所有启用的告警源
func (d *AlertStatsDao) GetAllAlertSources() ([]model.AlertSource, error) {
	var list []model.AlertSource
	err := d.db.Where("status = ?", 1).Order("create_time DESC").Find(&list).Error
	return list, err
}

// BatchDeleteAlertSources 批量删除告警源
func (d *AlertStatsDao) BatchDeleteAlertSources(ids []uint) error {
	return d.db.Where("id IN ?", ids).Delete(&model.AlertSource{}).Error
}

// UpdateAlertSourceStatus 更新告警源状态
func (d *AlertStatsDao) UpdateAlertSourceStatus(id uint, status int) error {
	return d.db.Model(&model.AlertSource{}).Where("id = ?", id).Update("status", status).Error
}

// CheckAlertSourceExists 检查告警源名称是否已存在
func (d *AlertStatsDao) CheckAlertSourceExists(name string, excludeID uint) (bool, error) {
	var count int64
	query := d.db.Model(&model.AlertSource{}).Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
