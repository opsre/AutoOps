package dao

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

// WebhookDao Webhook数据访问对象接口
type WebhookDao interface {
	CreateLog(log *model.WebhookLog) error
	CreateNotifyLog(log *model.WebhookNotifyLog) error
	GetLogByID(id uint) (*model.WebhookLog, error)
	GetLogList(source, status string, page, pageSize int) ([]model.WebhookLog, int64, error)
	GetNotifyLogsByWebhookID(webhookLogID uint) ([]model.WebhookNotifyLog, error)
}

// WebhookDaoImpl Webhook数据访问对象实现
type WebhookDaoImpl struct {
	db *gorm.DB
}

// NewWebhookDao 创建Webhook DAO实例
func NewWebhookDao() WebhookDao {
	return &WebhookDaoImpl{
		db: common.GetDB(),
	}
}

// CreateLog 创建Webhook日志记录
func (d *WebhookDaoImpl) CreateLog(log *model.WebhookLog) error {
	return d.db.Create(log).Error
}

// CreateNotifyLog 创建Webhook通知日志记录
func (d *WebhookDaoImpl) CreateNotifyLog(log *model.WebhookNotifyLog) error {
	return d.db.Create(log).Error
}

// GetLogByID 根据ID获取Webhook日志
func (d *WebhookDaoImpl) GetLogByID(id uint) (*model.WebhookLog, error) {
	var log model.WebhookLog
	err := d.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetLogList 分页获取Webhook日志列表
func (d *WebhookDaoImpl) GetLogList(source, status string, page, pageSize int) ([]model.WebhookLog, int64, error) {
	var logs []model.WebhookLog
	var total int64

	query := d.db.Model(&model.WebhookLog{})

	// 按来源筛选
	if source != "" {
		query = query.Where("source = ?", source)
	}

	// 按状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetNotifyLogsByWebhookID 根据WebhookLogID获取通知日志列表
func (d *WebhookDaoImpl) GetNotifyLogsByWebhookID(webhookLogID uint) ([]model.WebhookNotifyLog, error) {
	var logs []model.WebhookNotifyLog
	err := d.db.Where("webhook_log_id = ?", webhookLogID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}
