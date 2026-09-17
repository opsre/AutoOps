// 对话历史DAO
package dao

import (
	"dodevops-api/api/agent/model"
	"dodevops-api/common"
)

// ChatHistoryDao 对话历史数据访问接口
type ChatHistoryDao interface {
	Create(history *model.AgentChatHistory) error
	GetBySessionID(sessionID string, limit int) ([]model.AgentChatHistory, error)
	GetByUserID(userID uint, pageNum, pageSize int) ([]model.AgentChatHistory, int64, error)
	DeleteBySessionID(sessionID string) error
}

// ChatHistoryDaoImpl 对话历史DAO实现
type ChatHistoryDaoImpl struct{}

// NewChatHistoryDao 创建对话历史DAO实例
func NewChatHistoryDao() ChatHistoryDao {
	return &ChatHistoryDaoImpl{}
}

// Create 创建对话记录
func (d *ChatHistoryDaoImpl) Create(history *model.AgentChatHistory) error {
	return common.GetDB().Create(history).Error
}

// GetBySessionID 根据会话ID获取对话历史
func (d *ChatHistoryDaoImpl) GetBySessionID(sessionID string, limit int) ([]model.AgentChatHistory, error) {
	var histories []model.AgentChatHistory
	query := common.GetDB().Where("session_id = ?", sessionID).
		Order("create_time ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&histories).Error
	return histories, err
}

// GetByUserID 根据用户ID分页获取对话历史
func (d *ChatHistoryDaoImpl) GetByUserID(userID uint, pageNum, pageSize int) ([]model.AgentChatHistory, int64, error) {
	var histories []model.AgentChatHistory
	var total int64

	db := common.GetDB().Model(&model.AgentChatHistory{}).Where("user_id = ?", userID)

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (pageNum - 1) * pageSize
	err := db.Order("create_time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&histories).Error

	return histories, total, err
}

// DeleteBySessionID 删除指定会话的所有历史记录
func (d *ChatHistoryDaoImpl) DeleteBySessionID(sessionID string) error {
	return common.GetDB().Where("session_id = ?", sessionID).Delete(&model.AgentChatHistory{}).Error
}
