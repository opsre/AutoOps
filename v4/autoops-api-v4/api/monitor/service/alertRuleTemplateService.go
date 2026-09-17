package service

import (
	"fmt"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type AlertRuleTemplateService interface {
	CreateTemplate(c *gin.Context, dto *model.CreateRuleTemplateDto)
	UpdateTemplate(c *gin.Context, dto *model.UpdateRuleTemplateDto)
	DeleteTemplate(c *gin.Context, id uint)
	GetTemplate(c *gin.Context, id uint)
	GetTemplateList(c *gin.Context, category string, status *int, keyword string, page, pageSize int)
	UpdateTemplateStatus(c *gin.Context, id uint, status int)
	GetAllCategories(c *gin.Context)
}

type AlertRuleTemplateServiceImpl struct {
	templateDao dao.AlertRuleTemplateDao
}

func NewAlertRuleTemplateService() AlertRuleTemplateService {
	return &AlertRuleTemplateServiceImpl{
		templateDao: dao.NewAlertRuleTemplateDao(),
	}
}

// CreateTemplate 创建规则模板
func (s *AlertRuleTemplateServiceImpl) CreateTemplate(c *gin.Context, dto *model.CreateRuleTemplateDto) {
	template := &model.AlertRuleTemplate{
		Name:        dto.Name,
		Description: dto.Description,
		Category:    dto.Category,
		MetricKey:   dto.MetricKey,
		Unit:        dto.Unit,
		ValueType:   dto.ValueType,
		Tags:        dto.Tags,
		Status:      dto.Status,
	}

	// 设置默认值
	if template.ValueType == "" {
		template.ValueType = "number"
	}
	if template.Status == 0 {
		template.Status = 1
	}

	if err := s.templateDao.Create(template); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建规则模板失败: %v", err))
		return
	}

	result.Success(c, template)
}

// UpdateTemplate 更新规则模板
func (s *AlertRuleTemplateServiceImpl) UpdateTemplate(c *gin.Context, dto *model.UpdateRuleTemplateDto) {
	// 检查模板是否存在
	existing, err := s.templateDao.GetByID(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("规则模板不存在: %v", err))
		return
	}

	// 不允许修改内置模板
	if existing.IsBuiltin {
		result.Failed(c, int(result.ApiCode.FAILED), "内置模板不允许修改")
		return
	}

	template := &model.AlertRuleTemplate{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Category:    dto.Category,
		MetricKey:   dto.MetricKey,
		Unit:        dto.Unit,
		ValueType:   dto.ValueType,
		Tags:        dto.Tags,
		Status:      dto.Status,
	}

	if err := s.templateDao.Update(template); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新规则模板失败: %v", err))
		return
	}

	result.Success(c, template)
}

// DeleteTemplate 删除规则模板
func (s *AlertRuleTemplateServiceImpl) DeleteTemplate(c *gin.Context, id uint) {
	// 检查模板是否存在
	existing, err := s.templateDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("规则模板不存在: %v", err))
		return
	}

	// 不允许删除内置模板
	if existing.IsBuiltin {
		result.Failed(c, int(result.ApiCode.FAILED), "内置模板不允许删除")
		return
	}

	if err := s.templateDao.Delete(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("删除规则模板失败: %v", err))
		return
	}

	result.Success(c, "删除成功")
}

// GetTemplate 获取规则模板详情
func (s *AlertRuleTemplateServiceImpl) GetTemplate(c *gin.Context, id uint) {
	template, err := s.templateDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询规则模板失败: %v", err))
		return
	}

	result.Success(c, template)
}

// GetTemplateList 获取规则模板列表
func (s *AlertRuleTemplateServiceImpl) GetTemplateList(c *gin.Context, category string, status *int, keyword string, page, pageSize int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	templates, total, err := s.templateDao.GetList(category, status, keyword, page, pageSize)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询规则模板列表失败: %v", err))
		return
	}

	result.Success(c, gin.H{
		"list":     templates,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// UpdateTemplateStatus 更新规则模板状态
func (s *AlertRuleTemplateServiceImpl) UpdateTemplateStatus(c *gin.Context, id uint, status int) {
	// 检查模板是否存在
	_, err := s.templateDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("规则模板不存在: %v", err))
		return
	}

	if err := s.templateDao.UpdateStatus(id, status); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新规则模板状态失败: %v", err))
		return
	}

	statusText := "禁用"
	if status == 1 {
		statusText = "启用"
	}

	result.Success(c, fmt.Sprintf("规则模板已%s", statusText))
}

// GetAllCategories 获取所有分类
func (s *AlertRuleTemplateServiceImpl) GetAllCategories(c *gin.Context) {
	categories := []gin.H{
		{"value": "cpu", "label": "CPU"},
		{"value": "memory", "label": "内存"},
		{"value": "load", "label": "负载"},
		{"value": "disk", "label": "磁盘"},
		{"value": "io", "label": "IO"},
		{"value": "network", "label": "网络"},
		{"value": "connection", "label": "连接数"},
		{"value": "service", "label": "服务"},
		{"value": "custom", "label": "自定义"},
	}

	result.Success(c, categories)
}
