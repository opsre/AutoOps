package dao

import (
	"time"

	"dodevops-api/api/monitor/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

// DomainScheduleDao 域名监控定时任务配置数据访问对象接口
type DomainScheduleDao interface {
	GetOrCreate() (*model.DomainSchedule, error)
	Update(schedule *model.DomainSchedule) error
	UpdateStatus(status string) error
	UpdateNextRunTime(nextRunAt time.Time) error
}

// DomainScheduleDaoImpl 域名监控定时任务配置数据访问对象实现
type DomainScheduleDaoImpl struct {
	db *gorm.DB
}

// NewDomainScheduleDao 创建DomainSchedule DAO实例
func NewDomainScheduleDao() DomainScheduleDao {
	return &DomainScheduleDaoImpl{
		db: common.GetDB(),
	}
}

// GetOrCreate 获取或创建定时任务配置（单例）
func (d *DomainScheduleDaoImpl) GetOrCreate() (*model.DomainSchedule, error) {
	var schedule model.DomainSchedule

	// 尝试获取第一条记录
	err := d.db.First(&schedule).Error
	if err == gorm.ErrRecordNotFound {
		// 不存在则创建默认配置
		schedule = model.DomainSchedule{
			Enabled:  false,
			CronExpr: "0 */2 * * *", // 默认每2小时执行一次
			Status:   "stopped",
		}
		schedule.CreateTime.Time = time.Now()
		schedule.CreateTime.Valid = true
		schedule.UpdateTime.Time = time.Now()
		schedule.UpdateTime.Valid = true

		if err := d.db.Create(&schedule).Error; err != nil {
			return nil, err
		}
		return &schedule, nil
	}

	if err != nil {
		return nil, err
	}

	return &schedule, nil
}

// Update 更新定时任务配置
func (d *DomainScheduleDaoImpl) Update(schedule *model.DomainSchedule) error {
	schedule.UpdateTime.Time = time.Now()
	schedule.UpdateTime.Valid = true
	return d.db.Save(schedule).Error
}

// UpdateStatus 更新任务状态
func (d *DomainScheduleDaoImpl) UpdateStatus(status string) error {
	// 先查询第一条记录的ID
	var schedule model.DomainSchedule
	if err := d.db.First(&schedule).Error; err != nil {
		return err
	}

	return d.db.Model(&model.DomainSchedule{}).
		Where("id = ?", schedule.ID).
		Updates(map[string]interface{}{
			"status":      status,
			"update_time": time.Now(),
		}).Error
}

// UpdateNextRunTime 更新下次执行时间
func (d *DomainScheduleDaoImpl) UpdateNextRunTime(nextRunAt time.Time) error {
	// 先查询第一条记录的ID
	var schedule model.DomainSchedule
	if err := d.db.First(&schedule).Error; err != nil {
		return err
	}

	return d.db.Model(&model.DomainSchedule{}).
		Where("id = ?", schedule.ID).
		Updates(map[string]interface{}{
			"next_run_at": nextRunAt,
			"last_run_at": time.Now(),
			"update_time": time.Now(),
		}).Error
}
