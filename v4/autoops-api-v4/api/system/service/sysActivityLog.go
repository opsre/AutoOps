// 系统动态日志 服务层
package service

import (
	"time"

	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/util"
)

type SysActivityLogService struct {
	dao *dao.SysActivityLogDao
}

func NewSysActivityLogService() *SysActivityLogService {
	return &SysActivityLogService{
		dao: dao.NewSysActivityLogDao(),
	}
}

// Create 创建系统动态日志
func (s *SysActivityLogService) Create(log *model.SysActivityLog) error {
	log.CreateTime = util.HTime{Time: time.Now()}
	return s.dao.Create(log)
}

// CreateWithDto 使用DTO创建动态日志
func (s *SysActivityLogService) CreateWithDto(dto *model.CreateActivityLogDto) error {
	log := &model.SysActivityLog{
		ActivityType: dto.ActivityType,
		Title:        dto.Title,
		Content:      dto.Content,
		Status:       dto.Status,
		RelatedID:    dto.RelatedID,
		Summary:      dto.Summary,
		Duration:     dto.Duration,
		CreateTime:   util.HTime{Time: time.Now()},
	}
	return s.dao.Create(log)
}

// GetByID 根据ID查询
func (s *SysActivityLogService) GetByID(id uint) (*model.SysActivityLog, error) {
	return s.dao.GetByID(id)
}

// List 分页查询动态日志列表
func (s *SysActivityLogService) List(activityType int, status int, page, pageSize int) ([]model.SysActivityLog, int64, error) {
	return s.dao.List(activityType, status, page, pageSize)
}

// GetLatest 获取最新的N条动态日志
func (s *SysActivityLogService) GetLatest(limit int) ([]model.SysActivityLog, error) {
	return s.dao.GetLatest(limit)
}

// GetByRelatedID 根据关联ID查询
func (s *SysActivityLogService) GetByRelatedID(activityType int, relatedID uint) ([]model.SysActivityLog, error) {
	return s.dao.GetByRelatedID(activityType, relatedID)
}

// DeleteByID 根据ID删除
func (s *SysActivityLogService) DeleteByID(id uint) error {
	return s.dao.DeleteByID(id)
}

// DeleteBatch 批量删除
func (s *SysActivityLogService) DeleteBatch(ids []uint) error {
	return s.dao.DeleteBatch(ids)
}

// Clean 清空动态日志（保留最近N天）
func (s *SysActivityLogService) Clean(days int) error {
	return s.dao.Clean(days)
}
