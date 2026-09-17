// Prompt模板管理服务
package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"dodevops-api/api/agent/dao"
	"dodevops-api/api/agent/model"

	"github.com/gin-gonic/gin"
)

// PromptTemplateService Prompt模板服务接口
type PromptTemplateService interface {
	CreateTemplate(c *gin.Context, dto *model.PromptTemplateDto) error
	GetTemplateByID(c *gin.Context, id uint) (*model.PromptTemplate, error)
	UpdateTemplate(c *gin.Context, dto *model.PromptTemplateDto) error
	DeleteTemplate(c *gin.Context, id uint) error
	GetTemplateList(c *gin.Context, query *model.PromptTemplateQuery) ([]model.PromptTemplate, int64, error)
	GetTemplate(category string, scene string) (*model.PromptTemplate, error)
	BuildPrompt(templateName string, vars map[string]interface{}) (string, error)
}

// PromptTemplateServiceImpl Prompt模板服务实现
type PromptTemplateServiceImpl struct {
	templateDao dao.PromptTemplateDao
}

// NewPromptTemplateService 创建Prompt模板服务实例
func NewPromptTemplateService() PromptTemplateService {
	return &PromptTemplateServiceImpl{
		templateDao: dao.NewPromptTemplateDao(),
	}
}

// CreateTemplate 创建模板
func (s *PromptTemplateServiceImpl) CreateTemplate(c *gin.Context, dto *model.PromptTemplateDto) error {
	// 检查名称是否已存在
	existing, _ := s.templateDao.GetByName(dto.Name)
	if existing != nil {
		return fmt.Errorf("模板名称已存在")
	}

	// 编码变量列表
	variablesJSON, err := json.Marshal(dto.Variables)
	if err != nil {
		return fmt.Errorf("变量列表格式错误: %v", err)
	}

	template := &model.PromptTemplate{
		Name:         dto.Name,
		Category:     dto.Category,
		Scene:        dto.Scene,
		Template:     dto.Template,
		Variables:    string(variablesJSON),
		SystemPrompt: dto.SystemPrompt,
		Temperature:  dto.Temperature,
		MaxTokens:    dto.MaxTokens,
		ModelID:      dto.ModelID,
		Enabled:      dto.Enabled,
	}

	return s.templateDao.Create(template)
}

// GetTemplateByID 根据ID获取模板
func (s *PromptTemplateServiceImpl) GetTemplateByID(c *gin.Context, id uint) (*model.PromptTemplate, error) {
	return s.templateDao.GetByID(id)
}

// UpdateTemplate 更新模板
func (s *PromptTemplateServiceImpl) UpdateTemplate(c *gin.Context, dto *model.PromptTemplateDto) error {
	// 检查模板是否存在
	template, err := s.templateDao.GetByID(dto.ID)
	if err != nil {
		return fmt.Errorf("模板不存在")
	}

	// 如果修改了名称，检查新名称是否已存在
	if dto.Name != template.Name {
		existing, _ := s.templateDao.GetByName(dto.Name)
		if existing != nil && existing.ID != dto.ID {
			return fmt.Errorf("模板名称已存在")
		}
	}

	// 编码变量列表
	variablesJSON, err := json.Marshal(dto.Variables)
	if err != nil {
		return fmt.Errorf("变量列表格式错误: %v", err)
	}

	// 更新字段
	template.Name = dto.Name
	template.Category = dto.Category
	template.Scene = dto.Scene
	template.Template = dto.Template
	template.Variables = string(variablesJSON)
	template.SystemPrompt = dto.SystemPrompt
	template.Temperature = dto.Temperature
	template.MaxTokens = dto.MaxTokens
	template.ModelID = dto.ModelID
	template.Enabled = dto.Enabled

	return s.templateDao.Update(template)
}

// DeleteTemplate 删除模板
func (s *PromptTemplateServiceImpl) DeleteTemplate(c *gin.Context, id uint) error {
	// 检查是否存在
	_, err := s.templateDao.GetByID(id)
	if err != nil {
		return fmt.Errorf("模板不存在")
	}

	return s.templateDao.Delete(id)
}

// GetTemplateList 获取模板列表
func (s *PromptTemplateServiceImpl) GetTemplateList(c *gin.Context, query *model.PromptTemplateQuery) ([]model.PromptTemplate, int64, error) {
	return s.templateDao.GetList(query)
}

// GetTemplate 根据分类和场景获取模板
func (s *PromptTemplateServiceImpl) GetTemplate(category string, scene string) (*model.PromptTemplate, error) {
	return s.templateDao.GetByCategoryAndScene(category, scene)
}

// BuildPrompt 构建Prompt
func (s *PromptTemplateServiceImpl) BuildPrompt(templateName string, vars map[string]interface{}) (string, error) {
	// 获取模板
	template, err := s.templateDao.GetByName(templateName)
	if err != nil {
		return "", fmt.Errorf("模板不存在: %s", templateName)
	}

	if template.Enabled != 1 {
		return "", fmt.Errorf("模板未启用: %s", templateName)
	}

	// 构建Prompt
	prompt := template.Template

	// 变量替换
	for key, value := range vars {
		placeholder := fmt.Sprintf("{{%s}}", key)
		prompt = strings.ReplaceAll(prompt, placeholder, fmt.Sprintf("%v", value))
	}

	return prompt, nil
}

// PromptBuilder Prompt构建器
type PromptBuilder struct {
	template *model.PromptTemplate
	vars     map[string]interface{}
}

// NewPromptBuilder 创建Prompt构建器
func NewPromptBuilder(template *model.PromptTemplate) *PromptBuilder {
	return &PromptBuilder{
		template: template,
		vars:     make(map[string]interface{}),
	}
}

// SetVar 设置变量
func (pb *PromptBuilder) SetVar(key string, value interface{}) *PromptBuilder {
	pb.vars[key] = value
	return pb
}

// SetVars 批量设置变量
func (pb *PromptBuilder) SetVars(vars map[string]interface{}) *PromptBuilder {
	for k, v := range vars {
		pb.vars[k] = v
	}
	return pb
}

// Build 构建Prompt
func (pb *PromptBuilder) Build() string {
	prompt := pb.template.Template

	// 变量替换
	for key, value := range pb.vars {
		placeholder := fmt.Sprintf("{{%s}}", key)
		prompt = strings.ReplaceAll(prompt, placeholder, fmt.Sprintf("%v", value))
	}

	return prompt
}

// GetSystemPrompt 获取系统提示词
func (pb *PromptBuilder) GetSystemPrompt() string {
	return pb.template.SystemPrompt
}

// GetTemperature 获取温度参数
func (pb *PromptBuilder) GetTemperature() float64 {
	return pb.template.Temperature
}

// GetMaxTokens 获取最大Token数
func (pb *PromptBuilder) GetMaxTokens() int {
	return pb.template.MaxTokens
}

// GetModelID 获取模型ID
func (pb *PromptBuilder) GetModelID() uint {
	return pb.template.ModelID
}

// InitDefaultTemplates 初始化默认模板
func InitDefaultTemplates() error {
	templateDao := dao.NewPromptTemplateDao()

	// 巡检总结模板
	inspectionSummaryTemplate := &model.PromptTemplate{
		Name:     "inspection_summary",
		Category: "inspection",
		Scene:    "summary",
		Template: `你是一个专业的运维工程师，请根据以下巡检结果生成总结报告：

## 巡检信息
- 主机: {{hostName}} ({{hostIp}})
- 巡检时间: {{inspectionTime}}
- 总检查项: {{totalChecks}}
- 通过项: {{passedChecks}}
- 失败项: {{failedChecks}}
- 警告项: {{warnChecks}}

## 检查结果
{{checkResults}}

请生成一份简洁的总结，包括：
1. 整体健康状态评估
2. 主要问题汇总
3. 风险等级评估
4. 关键建议（不超过5条）

要求：
- 简洁专业，重点突出
- 按优先级排序
- 给出可执行的建议`,
		Variables:    `["hostName","hostIp","inspectionTime","totalChecks","passedChecks","failedChecks","warnChecks","checkResults"]`,
		SystemPrompt: "你是一个专业的运维工程师，擅长分析系统健康状况和提供运维建议。",
		Temperature:  0.7,
		MaxTokens:    1000,
		Enabled:      1,
	}

	// 故障根因分析模板
	incidentAnalysisTemplate := &model.PromptTemplate{
		Name:     "incident_analysis",
		Category: "incident",
		Scene:    "analysis",
		Template: `你是一个专业的运维故障分析专家，请根据以下信息进行根因分析：

## 故障基本信息
- 主机: {{hostName}} ({{hostIp}})
- 故障类别: {{category}}
- 检测时间: {{detectTime}}
- 症状: {{symptoms}}

## 证据数据
{{evidenceData}}

## 分析要求
1. 分析故障的根本原因
2. 列出支持该结论的关键证据
3. 评估故障的影响范围
4. 给出详细的解决方案
5. 提供可执行的修复脚本（如适用）
6. 评估自动修复的风险等级（0-3）

请按照以下JSON格式返回分析结果：
{
    "rootCause": "根因描述",
    "keyEvidence": ["证据1", "证据2"],
    "impact": "影响范围描述",
    "solution": "解决方案描述",
    "fixScript": "#!/bin/bash\n修复脚本",
    "riskLevel": 0-3,
    "riskReason": "风险评估说明"
}`,
		Variables:    `["hostName","hostIp","category","detectTime","symptoms","evidenceData"]`,
		SystemPrompt: "你是一个专业的运维故障分析专家，擅长快速定位问题根因并提供解决方案。",
		Temperature:  0.5,
		MaxTokens:    2000,
		Enabled:      1,
	}

	// 容量预测建议模板
	predictionSuggestionTemplate := &model.PromptTemplate{
		Name:     "prediction_suggestion",
		Category: "prediction",
		Scene:    "suggestion",
		Template: `你是一个专业的运维容量规划专家，请根据以下信息给出优化建议：

## 资源信息
- 目标: {{targetName}}
- 指标: {{metricName}} ({{metricType}})
- 当前值: {{currentValue}}%
- 阈值: {{threshold}}%

## 趋势分析
- 趋势类型: {{trendType}}
- 变化率: {{changeRate}}%/天
- 预测剩余天数: {{daysLeft}}天
- 置信度: {{confidence}}%

## 预测数据
未来7天预测值: {{predictions}}

## 要求
请给出以下建议：
1. 当前资源状态评估
2. 风险分析和影响评估
3. 具体的优化措施（优先级排序）
4. 扩容建议（如需要）
5. 成本优化建议
6. 长期优化策略

请用简洁、专业的语言回答，重点突出可执行的建议。`,
		Variables:    `["targetName","metricName","metricType","currentValue","threshold","trendType","changeRate","daysLeft","confidence","predictions"]`,
		SystemPrompt: "你是一个专业的运维容量规划专家，擅长资源趋势分析和容量规划。",
		Temperature:  0.7,
		MaxTokens:    1500,
		Enabled:      1,
	}

	// 检查并创建模板
	templates := []*model.PromptTemplate{
		inspectionSummaryTemplate,
		incidentAnalysisTemplate,
		predictionSuggestionTemplate,
	}

	for _, tpl := range templates {
		existing, _ := templateDao.GetByName(tpl.Name)
		if existing == nil {
			if err := templateDao.Create(tpl); err != nil {
				return fmt.Errorf("创建模板失败 %s: %v", tpl.Name, err)
			}
		}
	}

	return nil
}
