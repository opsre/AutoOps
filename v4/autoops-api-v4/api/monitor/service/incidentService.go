package service

import (
	cmdbDao "dodevops-api/api/cmdb/dao"
	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/pkg/jwt"
	"time"

	"github.com/gin-gonic/gin"
)

type IncidentServiceInterface interface {
	CreateIncident(c *gin.Context, dto *model.CreateIncidentDto)
	UpdateIncident(c *gin.Context, dto *model.UpdateIncidentDto)
	DeleteIncident(c *gin.Context, id uint)
	GetIncidentById(c *gin.Context, id uint)
	GetIncidentList(c *gin.Context, query *model.IncidentQueryDto)
	BatchDeleteIncidents(c *gin.Context, ids []uint)
	UpdateIncidentStatus(c *gin.Context, id uint, status int)
	GetIncidentStatistics(c *gin.Context)
}

type IncidentServiceImpl struct{}

// CreateIncident 创建故障记录
func (s *IncidentServiceImpl) CreateIncident(c *gin.Context, dto *model.CreateIncidentDto) {
	incidentDao := dao.NewIncidentDao()
	groupDao := cmdbDao.NewCmdbGroupDao()

	// 验证业务线ID（必须是二级分组，即 parent_id > 0）
	group, err := groupDao.GetCmdbGroupById(dto.BusinessLineID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "业务线不存在")
		return
	}
	if group.ParentID == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "业务线必须选择二级分组")
		return
	}

	// 从 token 获取当前登录用户信息
	currentUser, err := jwt.GetAdmin(c)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "无法获取当前用户信息")
		return
	}

	// 设置默认值
	if dto.AlertLevel == "" {
		dto.AlertLevel = "P4" // 默认P4级别
	}
	if dto.Status == 0 {
		dto.Status = model.IncidentStatusPending // 默认未处理
	}

	incident := &model.Incident{
		AlertTime:      dto.AlertTime,
		BusinessLineID: dto.BusinessLineID,
		Frequency:      dto.Frequency,
		AlertDesc:      dto.AlertDesc,
		AlertLevel:     dto.AlertLevel,
		IncidentCause:  dto.IncidentCause,
		Department:     dto.Department,
		Solution:       dto.Solution,
		DetailURL:      dto.DetailURL,
		Handler:        currentUser.Username, // 从 token 自动获取
		HandlerID:      currentUser.ID,       // 从 token 自动获取
		Status:         dto.Status,
		Remark:         dto.Remark,
		CreateTime:     util.HTime{Time: time.Now()},
		UpdateTime:     util.HTime{Time: time.Now()},
	}

	err = incidentDao.CreateIncident(incident)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "创建故障记录失败")
		return
	}

	// 填充业务线名称
	incident.BusinessLineName = group.Name

	result.Success(c, incident.ToVO())
}

// UpdateIncident 更新故障记录
func (s *IncidentServiceImpl) UpdateIncident(c *gin.Context, dto *model.UpdateIncidentDto) {
	incidentDao := dao.NewIncidentDao()
	groupDao := cmdbDao.NewCmdbGroupDao()

	// 检查记录是否存在
	_, err := incidentDao.GetIncidentById(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "故障记录不存在")
		return
	}

	// 如果更新了业务线ID，需要验证
	if dto.BusinessLineID > 0 {
		group, err := groupDao.GetCmdbGroupById(dto.BusinessLineID)
		if err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), "业务线不存在")
			return
		}
		if group.ParentID == 0 {
			result.Failed(c, int(result.ApiCode.FAILED), "业务线必须选择二级分组")
			return
		}
	}

	incident := &model.Incident{
		ID:             dto.ID,
		AlertTime:      dto.AlertTime,
		BusinessLineID: dto.BusinessLineID,
		Frequency:      dto.Frequency,
		AlertDesc:      dto.AlertDesc,
		AlertLevel:     dto.AlertLevel,
		IncidentCause:  dto.IncidentCause,
		Department:     dto.Department,
		Solution:       dto.Solution,
		DetailURL:      dto.DetailURL,
		// 注意：Handler 和 HandlerID 不允许更新
		Status:     dto.Status,
		Remark:     dto.Remark,
		UpdateTime: util.HTime{Time: time.Now()},
	}

	err = incidentDao.UpdateIncident(dto.ID, incident)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "更新故障记录失败")
		return
	}

	result.Success(c, true)
}

// DeleteIncident 删除故障记录
func (s *IncidentServiceImpl) DeleteIncident(c *gin.Context, id uint) {
	incidentDao := dao.NewIncidentDao()

	err := incidentDao.DeleteIncident(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "删除故障记录失败")
		return
	}

	result.Success(c, true)
}

// GetIncidentById 根据ID获取故障记录
func (s *IncidentServiceImpl) GetIncidentById(c *gin.Context, id uint) {
	incidentDao := dao.NewIncidentDao()
	groupDao := cmdbDao.NewCmdbGroupDao()

	incident, err := incidentDao.GetIncidentById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "故障记录不存在")
		return
	}

	// 填充业务线名称
	if incident.BusinessLineID > 0 {
		group, err := groupDao.GetCmdbGroupById(incident.BusinessLineID)
		if err == nil {
			incident.BusinessLineName = group.Name
		}
	}

	result.Success(c, incident.ToVO())
}

// GetIncidentList 获取故障列表
func (s *IncidentServiceImpl) GetIncidentList(c *gin.Context, query *model.IncidentQueryDto) {
	incidentDao := dao.NewIncidentDao()
	groupDao := cmdbDao.NewCmdbGroupDao()

	// 设置默认分页参数
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	list, total, err := incidentDao.GetIncidentList(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询故障列表失败")
		return
	}

	// 填充业务线名称并转换为VO
	var voList []*model.IncidentVO
	for _, incident := range list {
		// 填充业务线名称
		if incident.BusinessLineID > 0 {
			group, err := groupDao.GetCmdbGroupById(incident.BusinessLineID)
			if err == nil {
				incident.BusinessLineName = group.Name
			}
		}
		voList = append(voList, incident.ToVO())
	}

	result.Success(c, gin.H{
		"list":     voList,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	})
}

// BatchDeleteIncidents 批量删除故障记录
func (s *IncidentServiceImpl) BatchDeleteIncidents(c *gin.Context, ids []uint) {
	if len(ids) == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "请选择要删除的记录")
		return
	}

	incidentDao := dao.NewIncidentDao()
	err := incidentDao.BatchDeleteIncidents(ids)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "批量删除失败")
		return
	}

	result.Success(c, true)
}

// UpdateIncidentStatus 更新故障状态
func (s *IncidentServiceImpl) UpdateIncidentStatus(c *gin.Context, id uint, status int) {
	// 验证状态值
	if status < model.IncidentStatusPending || status > model.IncidentStatusCompleted {
		result.Failed(c, int(result.ApiCode.FAILED), "无效的状态值")
		return
	}

	incidentDao := dao.NewIncidentDao()
	err := incidentDao.UpdateIncidentStatus(id, status)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "更新状态失败")
		return
	}

	result.Success(c, true)
}

// GetIncidentStatistics 获取故障统计信息
func (s *IncidentServiceImpl) GetIncidentStatistics(c *gin.Context) {
	incidentDao := dao.NewIncidentDao()

	stats, err := incidentDao.GetIncidentStatistics()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "获取统计信息失败")
		return
	}

	result.Success(c, stats)
}

// GetIncidentService 获取故障管理服务实例
func GetIncidentService() IncidentServiceInterface {
	return &IncidentServiceImpl{}
}
