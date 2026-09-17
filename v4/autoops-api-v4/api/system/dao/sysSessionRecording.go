// 会话录制 数据层
package dao

import (
	"dodevops-api/api/system/model"
	. "dodevops-api/pkg/db"
)

// CreateSessionRecording 创建会话录制记录
func CreateSessionRecording(recording *model.SysSessionRecording) error {
	return Db.Create(recording).Error
}

// UpdateSessionRecording 更新会话录制记录
func UpdateSessionRecording(recording *model.SysSessionRecording) error {
	return Db.Save(recording).Error
}

// UpdateSessionRecordingBySessionID 根据SessionID更新录制记录
func UpdateSessionRecordingBySessionID(sessionID string, updates map[string]interface{}) error {
	return Db.Model(&model.SysSessionRecording{}).
		Where("session_id = ?", sessionID).
		Updates(updates).Error
}

// GetSessionRecordingByID 根据ID获取录制记录
func GetSessionRecordingByID(id uint) (*model.SysSessionRecording, error) {
	var recording model.SysSessionRecording
	err := Db.Where("id = ?", id).First(&recording).Error
	if err != nil {
		return nil, err
	}
	return &recording, nil
}

// GetSessionRecordingBySessionID 根据SessionID获取录制记录
func GetSessionRecordingBySessionID(sessionID string) (*model.SysSessionRecording, error) {
	var recording model.SysSessionRecording
	err := Db.Where("session_id = ?", sessionID).First(&recording).Error
	if err != nil {
		return nil, err
	}
	return &recording, nil
}

// GetSessionRecordingList 分页查询会话录制列表
func GetSessionRecordingList(dto model.SessionRecordingListDto) ([]model.SysSessionRecording, int64, error) {
	var recordings []model.SysSessionRecording
	var count int64

	query := Db.Model(&model.SysSessionRecording{})

	// 条件过滤
	if dto.AdminID != 0 {
		query = query.Where("admin_id = ?", dto.AdminID)
	}
	if dto.HostID != 0 {
		query = query.Where("host_id = ?", dto.HostID)
	}
	if dto.Status != 0 {
		query = query.Where("status = ?", dto.Status)
	}
	if dto.RiskLevel != 0 {
		query = query.Where("risk_level = ?", dto.RiskLevel)
	}
	if dto.Keyword != "" {
		keyword := "%" + dto.Keyword + "%"
		query = query.Where("username LIKE ? OR host_name LIKE ? OR host_ip LIKE ?", keyword, keyword, keyword)
	}
	if dto.StartTime != "" {
		query = query.Where("start_time >= ?", dto.StartTime)
	}
	if dto.EndTime != "" {
		query = query.Where("start_time <= ?", dto.EndTime)
	}

	// 统计总数
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (dto.Page - 1) * dto.PageSize
	if err := query.Order("start_time DESC").
		Limit(dto.PageSize).
		Offset(offset).
		Find(&recordings).Error; err != nil {
		return nil, 0, err
	}

	return recordings, count, nil
}

// DeleteSessionRecordingByID 根据ID删除录制记录
func DeleteSessionRecordingByID(id uint) error {
	return Db.Delete(&model.SysSessionRecording{}, id).Error
}

// BatchDeleteSessionRecording 批量删除录制记录
func BatchDeleteSessionRecording(ids []uint) error {
	return Db.Where("id IN ?", ids).Delete(&model.SysSessionRecording{}).Error
}
