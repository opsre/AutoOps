package service

import (
	"fmt"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// DomainScheduleService 域名监控定时任务服务接口
type DomainScheduleService interface {
	GetScheduleConfig(c *gin.Context)
	UpdateScheduleConfig(c *gin.Context, dto *model.DomainScheduleDto)
	PauseSchedule(c *gin.Context)
	ResumeSchedule(c *gin.Context)
	GetScheduleStatus(c *gin.Context)
}

// DomainScheduleServiceImpl 域名监控定时任务服务实现
type DomainScheduleServiceImpl struct {
	scheduleDao dao.DomainScheduleDao
	scheduler   *DomainScheduler
}

// NewDomainScheduleService 创建定时任务服务实例
func NewDomainScheduleService() DomainScheduleService {
	return &DomainScheduleServiceImpl{
		scheduleDao: dao.NewDomainScheduleDao(),
		scheduler:   GetDomainScheduler(),
	}
}

// GetScheduleConfig 获取定时任务配置
func (s *DomainScheduleServiceImpl) GetScheduleConfig(c *gin.Context) {
	schedule, err := s.scheduleDao.GetOrCreate()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取配置失败: %v", err))
		return
	}

	result.Success(c, schedule)
}

// UpdateScheduleConfig 更新定时任务配置
func (s *DomainScheduleServiceImpl) UpdateScheduleConfig(c *gin.Context, dto *model.DomainScheduleDto) {
	// 参数验证
	if dto.CronExpr == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "Cron表达式不能为空")
		return
	}

	// 如果开启了通知，验证机器人ID
	if dto.NotifyEnabled {
		if dto.NotifyRobotId == nil {
			result.Failed(c, int(result.ApiCode.FAILED), "开启通知时必须选择通知机器人")
			return
		}

		// 验证机器人是否存在且启用
		robotDao := dao.NewNotifyRobotDao()
		robot, err := robotDao.GetByID(*dto.NotifyRobotId)
		if err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("通知机器人不存在: %v", err))
			return
		}

		if robot.Status != 1 {
			result.Failed(c, int(result.ApiCode.FAILED), "选择的通知机器人未启用")
			return
		}
	}

	// 获取当前配置
	schedule, err := s.scheduleDao.GetOrCreate()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取配置失败: %v", err))
		return
	}

	// 更新配置
	schedule.Enabled = dto.Enabled
	schedule.CronExpr = dto.CronExpr
	schedule.NotifyEnabled = dto.NotifyEnabled
	schedule.NotifyRobotId = dto.NotifyRobotId

	if err := s.scheduleDao.Update(schedule); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新配置失败: %v", err))
		return
	}

	// 根据启用状态控制调度器
	if dto.Enabled {
		// 如果调度器正在运行，先停止
		if s.scheduler.IsRunning() {
			s.scheduler.Stop()
		}
		// 启动调度器
		if err := s.scheduler.Start(dto.CronExpr); err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("启动定时任务失败: %v", err))
			return
		}
	} else {
		// 停止调度器
		if s.scheduler.IsRunning() {
			s.scheduler.Stop()
		}
	}

	result.Success(c, schedule)
}

// PauseSchedule 暂停定时任务
func (s *DomainScheduleServiceImpl) PauseSchedule(c *gin.Context) {
	if err := s.scheduler.Pause(); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("暂停失败: %v", err))
		return
	}

	// 更新配置
	schedule, _ := s.scheduleDao.GetOrCreate()
	schedule.Enabled = false
	s.scheduleDao.Update(schedule)

	result.Success(c, map[string]interface{}{
		"message": "定时任务已暂停",
		"status":  "paused",
	})
}

// ResumeSchedule 恢复定时任务
func (s *DomainScheduleServiceImpl) ResumeSchedule(c *gin.Context) {
	// 获取配置
	schedule, err := s.scheduleDao.GetOrCreate()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取配置失败: %v", err))
		return
	}

	if err := s.scheduler.Resume(schedule.CronExpr); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("恢复失败: %v", err))
		return
	}

	// 更新配置
	schedule.Enabled = true
	s.scheduleDao.Update(schedule)

	result.Success(c, map[string]interface{}{
		"message": "定时任务已恢复",
		"status":  "running",
	})
}

// GetScheduleStatus 获取定时任务状态
func (s *DomainScheduleServiceImpl) GetScheduleStatus(c *gin.Context) {
	status := s.scheduler.GetStatus()
	result.Success(c, status)
}
