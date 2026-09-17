package dao

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

// NotifyRobotDao 通知机器人数据访问对象接口
type NotifyRobotDao interface {
	Create(robot *model.NotifyRobot) error
	GetByID(id uint) (*model.NotifyRobot, error)
	Update(robot *model.NotifyRobot) error
	Delete(id uint) error
	GetList(robotType string, status *int, page, pageSize int) ([]model.NotifyRobot, int64, error)
	GetByType(robotType string) ([]model.NotifyRobot, error)
	GetEnabledRobots() ([]model.NotifyRobot, error)
}

// NotifyRobotDaoImpl 通知机器人数据访问对象实现
type NotifyRobotDaoImpl struct {
	db *gorm.DB
}

// NewNotifyRobotDao 创建NotifyRobot DAO实例
func NewNotifyRobotDao() NotifyRobotDao {
	return &NotifyRobotDaoImpl{
		db: common.GetDB(),
	}
}

// Create 创建通知机器人记录
func (d *NotifyRobotDaoImpl) Create(robot *model.NotifyRobot) error {
	return d.db.Create(robot).Error
}

// GetByID 根据ID获取通知机器人记录
func (d *NotifyRobotDaoImpl) GetByID(id uint) (*model.NotifyRobot, error) {
	var robot model.NotifyRobot
	err := d.db.Where("id = ?", id).First(&robot).Error
	if err != nil {
		return nil, err
	}
	return &robot, nil
}

// Update 更新通知机器人记录
func (d *NotifyRobotDaoImpl) Update(robot *model.NotifyRobot) error {
	return d.db.Save(robot).Error
}

// Delete 删除通知机器人记录
func (d *NotifyRobotDaoImpl) Delete(id uint) error {
	return d.db.Delete(&model.NotifyRobot{}, id).Error
}

// GetList 分页获取通知机器人列表
func (d *NotifyRobotDaoImpl) GetList(robotType string, status *int, page, pageSize int) ([]model.NotifyRobot, int64, error) {
	var robots []model.NotifyRobot
	var total int64

	query := d.db.Model(&model.NotifyRobot{})

	// 按类型筛选
	if robotType != "" {
		query = query.Where("type = ?", robotType)
	}

	// 按状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&robots).Error
	if err != nil {
		return nil, 0, err
	}

	return robots, total, nil
}

// GetByType 根据类型获取所有机器人
func (d *NotifyRobotDaoImpl) GetByType(robotType string) ([]model.NotifyRobot, error) {
	var robots []model.NotifyRobot
	err := d.db.Where("type = ?", robotType).Find(&robots).Error
	return robots, err
}

// GetEnabledRobots 获取所有启用的机器人
func (d *NotifyRobotDaoImpl) GetEnabledRobots() ([]model.NotifyRobot, error) {
	var robots []model.NotifyRobot
	err := d.db.Where("status = ?", 1).Find(&robots).Error
	return robots, err
}
