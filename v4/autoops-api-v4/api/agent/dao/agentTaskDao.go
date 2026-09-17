// Agent任务数据访问层
package dao

import (
	"time"

	"dodevops-api/api/agent/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"
)

// AgentTaskDao Agent任务DAO接口
type AgentTaskDao interface {
	Create(task *model.AgentTask) error
	GetByID(id uint) (*model.AgentTask, error)
	Update(task *model.AgentTask) error
	Delete(id uint) error
	GetList(query *model.AgentTaskQuery) ([]model.AgentTask, int64, error)
	GetPendingTasks(limit int) ([]model.AgentTask, error)
	GetStatistics() (map[string]interface{}, error)
}

// AgentTaskDaoImpl Agent任务DAO实现
type AgentTaskDaoImpl struct{}

// NewAgentTaskDao 创建Agent任务DAO实例
func NewAgentTaskDao() AgentTaskDao {
	return &AgentTaskDaoImpl{}
}

// Create 创建任务
func (d *AgentTaskDaoImpl) Create(task *model.AgentTask) error {
	now := util.HTime{Time: time.Now()}
	task.CreateTime = now
	task.UpdateTime = now
	return Db.Create(task).Error
}

// GetByID 根据ID获取任务
func (d *AgentTaskDaoImpl) GetByID(id uint) (*model.AgentTask, error) {
	var task model.AgentTask
	err := Db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Update 更新任务
func (d *AgentTaskDaoImpl) Update(task *model.AgentTask) error {
	task.UpdateTime = util.HTime{Time: time.Now()}
	return Db.Model(&model.AgentTask{}).Where("id = ?", task.ID).Updates(task).Error
}

// Delete 删除任务
func (d *AgentTaskDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.AgentTask{}, id).Error
}

// GetList 获取任务列表
func (d *AgentTaskDaoImpl) GetList(query *model.AgentTaskQuery) ([]model.AgentTask, int64, error) {
	var tasks []model.AgentTask
	var total int64

	db := Db.Model(&model.AgentTask{})

	// 条件筛选
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.StartTime != "" {
		db = db.Where("create_time >= ?", query.StartTime)
	}
	if query.EndTime != "" {
		db = db.Where("create_time <= ?", query.EndTime)
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

	if err := db.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetPendingTasks 获取待执行的任务
func (d *AgentTaskDaoImpl) GetPendingTasks(limit int) ([]model.AgentTask, error) {
	var tasks []model.AgentTask
	err := Db.Where("status = ?", 1). // 状态=待执行
						Order("priority DESC, create_time ASC"). // 优先级高的先执行
						Limit(limit).
						Find(&tasks).Error
	return tasks, err
}

// GetStatistics 获取任务统计
func (d *AgentTaskDaoImpl) GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 按状态统计
	var statusStats []struct {
		Status int
		Count  int64
	}
	err := Db.Model(&model.AgentTask{}).
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
			statusMap["pending"] = stat.Count
		case 2:
			statusMap["running"] = stat.Count
		case 3:
			statusMap["success"] = stat.Count
		case 4:
			statusMap["failed"] = stat.Count
		case 5:
			statusMap["cancelled"] = stat.Count
		}
	}
	stats["by_status"] = statusMap

	// 按类型统计
	var typeStats []struct {
		Type  string
		Count int64
	}
	err = Db.Model(&model.AgentTask{}).
		Select("type, count(*) as count").
		Group("type").
		Scan(&typeStats).Error
	if err != nil {
		return nil, err
	}

	typeMap := make(map[string]int64)
	for _, stat := range typeStats {
		typeMap[stat.Type] = stat.Count
	}
	stats["by_type"] = typeMap

	// 今日任务数
	var todayCount int64
	Db.Model(&model.AgentTask{}).
		Where("DATE(create_time) = CURDATE()").
		Count(&todayCount)
	stats["today_count"] = todayCount

	return stats, nil
}
