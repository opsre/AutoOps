package service

import (
	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
)

type MonitorDataSourceService interface {
	Create(data *model.MonitorDataSource) error
	Delete(id uint) error
	Update(data *model.MonitorDataSource) error
	GetByID(id uint) (*model.MonitorDataSource, error)
	GetList(page, pageSize int) ([]*model.MonitorDataSource, int64, error)
}

type monitorDataSourceService struct {
	dao          dao.MonitorDataSourceDao
	groupRuleDao dao.MonitorAlertGroupRuleDao
	ruleDao      dao.MonitorAlertRuleDao
}

func NewMonitorDataSourceService() MonitorDataSourceService {
	return &monitorDataSourceService{
		dao:          dao.NewMonitorDataSourceDao(),
		groupRuleDao: dao.NewMonitorAlertGroupRuleDao(),
		ruleDao:      dao.NewMonitorAlertRuleDao(),
	}
}

func (s *monitorDataSourceService) Create(data *model.MonitorDataSource) error {
	return s.dao.Create(data)
}

func (s *monitorDataSourceService) Delete(id uint) error {
	// First get all groups belonging to this data source
	groups, err := s.groupRuleDao.GetAll()
	if err == nil {
		for _, g := range groups {
			if g.DataSourceID == id {
				// Delete rules under this group
				s.ruleDao.DeleteByGroupID(g.ID)
				// Delete the group
				s.groupRuleDao.Delete(g.ID)
			}
		}
	}
	return s.dao.Delete(id)
}

func (s *monitorDataSourceService) Update(data *model.MonitorDataSource) error {
	return s.dao.Update(data)
}

func (s *monitorDataSourceService) GetByID(id uint) (*model.MonitorDataSource, error) {
	return s.dao.GetByID(id)
}

func (s *monitorDataSourceService) GetList(page, pageSize int) ([]*model.MonitorDataSource, int64, error) {
	return s.dao.GetList(page, pageSize)
}
