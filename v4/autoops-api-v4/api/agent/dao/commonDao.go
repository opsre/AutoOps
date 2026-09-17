// 知识库和Prompt模板数据访问层
package dao

import (
	"time"

	"dodevops-api/api/agent/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"
)

// KnowledgeBaseDao 知识库DAO接口
type KnowledgeBaseDao interface {
	Create(knowledge *model.KnowledgeBase) error
	GetByID(id uint) (*model.KnowledgeBase, error)
	Update(knowledge *model.KnowledgeBase) error
	Delete(id uint) error
	GetList(query *model.KnowledgeBaseQuery) ([]model.KnowledgeBase, int64, error)
	Search(keyword string, typeFilter string, category string, limit int) ([]model.KnowledgeBase, error)
	IncrementUseCount(id uint) error
}

// KnowledgeBaseDaoImpl 知识库DAO实现
type KnowledgeBaseDaoImpl struct{}

// NewKnowledgeBaseDao 创建知识库DAO实例
func NewKnowledgeBaseDao() KnowledgeBaseDao {
	return &KnowledgeBaseDaoImpl{}
}

// Create 创建知识库
func (d *KnowledgeBaseDaoImpl) Create(knowledge *model.KnowledgeBase) error {
	now := util.HTime{Time: time.Now()}
	knowledge.CreateTime = now
	knowledge.UpdateTime = now
	return Db.Create(knowledge).Error
}

// GetByID 根据ID获取知识库
func (d *KnowledgeBaseDaoImpl) GetByID(id uint) (*model.KnowledgeBase, error) {
	var knowledge model.KnowledgeBase
	err := Db.Where("id = ?", id).First(&knowledge).Error
	if err != nil {
		return nil, err
	}
	return &knowledge, nil
}

// Update 更新知识库
func (d *KnowledgeBaseDaoImpl) Update(knowledge *model.KnowledgeBase) error {
	knowledge.UpdateTime = util.HTime{Time: time.Now()}
	return Db.Model(&model.KnowledgeBase{}).Where("id = ?", knowledge.ID).Updates(knowledge).Error
}

// Delete 删除知识库
func (d *KnowledgeBaseDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.KnowledgeBase{}, id).Error
}

// GetList 获取知识库列表
func (d *KnowledgeBaseDaoImpl) GetList(query *model.KnowledgeBaseQuery) ([]model.KnowledgeBase, int64, error) {
	var knowledges []model.KnowledgeBase
	var total int64

	db := Db.Model(&model.KnowledgeBase{})

	// 条件筛选
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if query.Title != "" {
		db = db.Where("title LIKE ?", "%"+query.Title+"%")
	}
	if query.Keyword != "" {
		db = db.Where("keywords LIKE ?", "%"+query.Keyword+"%")
	}
	if query.Tag != "" {
		db = db.Where("tags LIKE ?", "%"+query.Tag+"%")
	}
	if query.Enabled > 0 {
		db = db.Where("enabled = ?", query.Enabled)
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

	// 按使用次数和评分排序
	db = db.Order("use_count DESC, score DESC")

	if err := db.Find(&knowledges).Error; err != nil {
		return nil, 0, err
	}

	return knowledges, total, nil
}

// Search 搜索知识库
func (d *KnowledgeBaseDaoImpl) Search(keyword string, typeFilter string, category string, limit int) ([]model.KnowledgeBase, error) {
	var knowledges []model.KnowledgeBase

	db := Db.Where("enabled = ?", 1)

	if typeFilter != "" {
		db = db.Where("type = ?", typeFilter)
	}

	if category != "" {
		db = db.Where("category = ?", category)
	}

	if keyword != "" {
		db = db.Where("title LIKE ? OR content LIKE ? OR keywords LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := db.Order("use_count DESC, score DESC").
		Limit(limit).
		Find(&knowledges).Error

	return knowledges, err
}

// IncrementUseCount 增加使用次数
func (d *KnowledgeBaseDaoImpl) IncrementUseCount(id uint) error {
	return Db.Model(&model.KnowledgeBase{}).
		Where("id = ?", id).
		UpdateColumn("use_count", Db.Raw("use_count + 1")).Error
}

// PromptTemplateDao Prompt模板DAO接口
type PromptTemplateDao interface {
	Create(template *model.PromptTemplate) error
	GetByID(id uint) (*model.PromptTemplate, error)
	GetByName(name string) (*model.PromptTemplate, error)
	Update(template *model.PromptTemplate) error
	Delete(id uint) error
	GetList(query *model.PromptTemplateQuery) ([]model.PromptTemplate, int64, error)
	GetByCategoryAndScene(category string, scene string) (*model.PromptTemplate, error)
}

// PromptTemplateDaoImpl Prompt模板DAO实现
type PromptTemplateDaoImpl struct{}

// NewPromptTemplateDao 创建Prompt模板DAO实例
func NewPromptTemplateDao() PromptTemplateDao {
	return &PromptTemplateDaoImpl{}
}

// Create 创建Prompt模板
func (d *PromptTemplateDaoImpl) Create(template *model.PromptTemplate) error {
	now := util.HTime{Time: time.Now()}
	template.CreateTime = now
	template.UpdateTime = now
	return Db.Create(template).Error
}

// GetByID 根据ID获取Prompt模板
func (d *PromptTemplateDaoImpl) GetByID(id uint) (*model.PromptTemplate, error) {
	var template model.PromptTemplate
	err := Db.Where("id = ?", id).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByName 根据名称获取Prompt模板
func (d *PromptTemplateDaoImpl) GetByName(name string) (*model.PromptTemplate, error) {
	var template model.PromptTemplate
	err := Db.Where("name = ?", name).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Update 更新Prompt模板
func (d *PromptTemplateDaoImpl) Update(template *model.PromptTemplate) error {
	template.UpdateTime = util.HTime{Time: time.Now()}
	return Db.Model(&model.PromptTemplate{}).Where("id = ?", template.ID).Updates(template).Error
}

// Delete 删除Prompt模板
func (d *PromptTemplateDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.PromptTemplate{}, id).Error
}

// GetList 获取Prompt模板列表
func (d *PromptTemplateDaoImpl) GetList(query *model.PromptTemplateQuery) ([]model.PromptTemplate, int64, error) {
	var templates []model.PromptTemplate
	var total int64

	db := Db.Model(&model.PromptTemplate{})

	// 条件筛选
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if query.Scene != "" {
		db = db.Where("scene = ?", query.Scene)
	}
	if query.Enabled > 0 {
		db = db.Where("enabled = ?", query.Enabled)
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

	if err := db.Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// GetByCategoryAndScene 根据分类和场景获取模板
func (d *PromptTemplateDaoImpl) GetByCategoryAndScene(category string, scene string) (*model.PromptTemplate, error) {
	var template model.PromptTemplate
	err := Db.Where("category = ? AND scene = ? AND enabled = ?", category, scene, 1).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// AgentChatHistoryDao Agent对话历史DAO接口
type AgentChatHistoryDao interface {
	Create(history *model.AgentChatHistory) error
	GetBySessionID(sessionID string, limit int) ([]model.AgentChatHistory, error)
	DeleteBySessionID(sessionID string) error
	DeleteOldRecords(days int) error
}

// AgentChatHistoryDaoImpl Agent对话历史DAO实现
type AgentChatHistoryDaoImpl struct{}

// NewAgentChatHistoryDao 创建Agent对话历史DAO实例
func NewAgentChatHistoryDao() AgentChatHistoryDao {
	return &AgentChatHistoryDaoImpl{}
}

// Create 创建对话历史
func (d *AgentChatHistoryDaoImpl) Create(history *model.AgentChatHistory) error {
	history.CreateTime = util.HTime{Time: time.Now()}
	return Db.Create(history).Error
}

// GetBySessionID 根据会话ID获取对话历史
func (d *AgentChatHistoryDaoImpl) GetBySessionID(sessionID string, limit int) ([]model.AgentChatHistory, error) {
	var histories []model.AgentChatHistory

	db := Db.Where("session_id = ?", sessionID).
		Order("create_time ASC")

	if limit > 0 {
		// 获取最近的N条记录
		db = db.Limit(limit)
	}

	err := db.Find(&histories).Error
	return histories, err
}

// DeleteBySessionID 根据会话ID删除对话历史
func (d *AgentChatHistoryDaoImpl) DeleteBySessionID(sessionID string) error {
	return Db.Where("session_id = ?", sessionID).Delete(&model.AgentChatHistory{}).Error
}

// DeleteOldRecords 删除旧的对话记录
func (d *AgentChatHistoryDaoImpl) DeleteOldRecords(days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	return Db.Where("create_time < ?", cutoffTime).Delete(&model.AgentChatHistory{}).Error
}
