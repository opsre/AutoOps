package dao

import (
	"time"

	"dodevops-api/api/monitor/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

type AlertMuteDao interface {
	Create(mute *model.AlertMute) error
	Update(mute *model.AlertMute) error
	Delete(id uint) error
	GetByID(id uint) (*model.AlertMute, error)
	GetList(status *int, page, pageSize int) ([]*model.AlertMute, int64, error)
	GetActiveMutes() ([]*model.AlertMute, error)
}

type AlertMuteDaoImpl struct {
	db *gorm.DB
}

func NewAlertMuteDao() AlertMuteDao {
	return &AlertMuteDaoImpl{
		db: common.GetDB(),
	}
}

func (d *AlertMuteDaoImpl) Create(mute *model.AlertMute) error {
	return d.db.Create(mute).Error
}

func (d *AlertMuteDaoImpl) Update(mute *model.AlertMute) error {
	return d.db.Save(mute).Error
}

func (d *AlertMuteDaoImpl) Delete(id uint) error {
	return d.db.Delete(&model.AlertMute{}, id).Error
}

func (d *AlertMuteDaoImpl) GetByID(id uint) (*model.AlertMute, error) {
	var mute model.AlertMute
	err := d.db.First(&mute, id).Error
	return &mute, err
}

func (d *AlertMuteDaoImpl) GetList(status *int, page, pageSize int) ([]*model.AlertMute, int64, error) {
	var mutes []*model.AlertMute
	var total int64

	query := d.db.Model(&model.AlertMute{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&mutes).Error

	return mutes, total, err
}

func (d *AlertMuteDaoImpl) GetActiveMutes() ([]*model.AlertMute, error) {
	var mutes []*model.AlertMute
	now := time.Now()
	err := d.db.Where("status = ? AND start_time <= ? AND end_time >= ?", 1, now, now).Find(&mutes).Error
	return mutes, err
}
