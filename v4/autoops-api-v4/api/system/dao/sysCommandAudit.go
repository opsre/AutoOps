// 命令审计 数据访问层
package dao

import (
	"dodevops-api/api/system/model"
	. "dodevops-api/pkg/db"
)

// BatchCreateCommandAudit 批量创建命令审计记录
func BatchCreateCommandAudit(audits []*model.SysCommandAudit) error {
	if len(audits) == 0 {
		return nil
	}

	return Db.Create(&audits).Error
}

// GetCommandAuditList 获取命令审计列表
func GetCommandAuditList(dto *model.CommandAuditListDto) ([]*model.SysCommandAudit, int64, error) {
	var list []*model.SysCommandAudit
	var total int64

	query := Db.Model(&model.SysCommandAudit{})

	// 条件过滤
	if dto.RecordingID > 0 {
		query = query.Where("recording_id = ?", dto.RecordingID)
	}

	if dto.SessionID != "" {
		query = query.Where("session_id = ?", dto.SessionID)
	}

	if dto.IsSensitive != nil {
		query = query.Where("is_sensitive = ?", *dto.IsSensitive)
	}

	if dto.RiskLevel != nil {
		query = query.Where("risk_level = ?", *dto.RiskLevel)
	}

	if dto.Keyword != "" {
		query = query.Where("command LIKE ?", "%"+dto.Keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (dto.Page - 1) * dto.PageSize
	if err := query.Order("sequence ASC").
		Offset(offset).
		Limit(dto.PageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// DeleteCommandAuditByRecordingID 删除指定录制记录的所有命令审计
func DeleteCommandAuditByRecordingID(recordingID uint) error {
	return Db.Where("recording_id = ?", recordingID).Delete(&model.SysCommandAudit{}).Error
}

// DeleteCommandAuditByRecordingIDs 批量删除命令审计记录
func DeleteCommandAuditByRecordingIDs(recordingIDs []uint) error {
	if len(recordingIDs) == 0 {
		return nil
	}

	return Db.Where("recording_id IN ?", recordingIDs).Delete(&model.SysCommandAudit{}).Error
}

// GetSensitiveCommandCount 获取敏感命令数量
func GetSensitiveCommandCount(recordingID uint) (int64, error) {
	var count int64
	err := Db.Model(&model.SysCommandAudit{}).
		Where("recording_id = ? AND is_sensitive = ?", recordingID, true).
		Count(&count).Error
	return count, err
}

// GetHighestRiskLevel 获取最高风险等级
func GetHighestRiskLevel(recordingID uint) (int, error) {
	var result struct {
		MaxRiskLevel int
	}

	err := Db.Model(&model.SysCommandAudit{}).
		Select("COALESCE(MAX(risk_level), 0) as max_risk_level").
		Where("recording_id = ?", recordingID).
		Scan(&result).Error

	return result.MaxRiskLevel, err
}
