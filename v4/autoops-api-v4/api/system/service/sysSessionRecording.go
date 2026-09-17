// 会话录制 服务层
package service

import (
	"fmt"
	"io"
	"os"
	"time"

	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/util"
	"dodevops-api/common/util/recording"
)

// SessionRecordingService 会话录制服务
type SessionRecordingService struct {
	storage recording.Storage
}

// NewSessionRecordingService 创建会话录制服务
func NewSessionRecordingService(storagePath string) (*SessionRecordingService, error) {
	storage, err := recording.NewLocalStorage(storagePath)
	if err != nil {
		return nil, fmt.Errorf("初始化存储失败: %w", err)
	}

	return &SessionRecordingService{
		storage: storage,
	}, nil
}

// CreateRecorder 创建录制器
func (s *SessionRecordingService) CreateRecorder(metadata recording.RecordingMetadata) recording.Recorder {
	return recording.NewAsciicastRecorder(s.storage, metadata)
}

// StartRecording 开始录制会话（创建数据库记录）
func (s *SessionRecordingService) StartRecording(dto *model.CreateSessionRecordingDto) (*model.SysSessionRecording, error) {
	// 创建录制记录
	rec := &model.SysSessionRecording{
		SessionID:      dto.SessionID,
		AdminID:        dto.AdminID,
		Username:       dto.Username,
		HostID:         dto.HostID,
		HostName:       dto.HostName,
		HostIP:         dto.HostIP,
		SSHUser:        dto.SSHUser,
		TerminalWidth:  dto.TerminalWidth,
		TerminalHeight: dto.TerminalHeight,
		ClientIP:       dto.ClientIP,
		UserAgent:      dto.UserAgent,
		StartTime:      util.HTime{Time: time.Now()},
		Status:         model.SessionStatusRecording,
		StorageType:    model.StorageTypeLocal,
		FilePath:       "", // 稍后由录制器填充
		CreateTime:     util.HTime{Time: time.Now()},
	}

	// 保存到数据库
	if err := dao.CreateSessionRecording(rec); err != nil {
		return nil, fmt.Errorf("创建录制记录失败: %w", err)
	}

	return rec, nil
}

// FinishRecording 完成录制（更新数据库记录并保存命令审计）
func (s *SessionRecordingService) FinishRecording(sessionID string, info *recording.RecordingInfo, recorder recording.Recorder) error {
	// 获取录制记录
	rec, err := dao.GetSessionRecordingBySessionID(sessionID)
	if err != nil {
		return fmt.Errorf("获取录制记录失败: %w", err)
	}

	// 保存命令审计（如果录制器支持）
	if asciiRecorder, ok := recorder.(*recording.AsciicastRecorder); ok {
		commands := asciiRecorder.GetParsedCommands()
		if len(commands) > 0 {
			// 创建命令审计服务
			cmdAuditService := NewCommandAuditService()

			// 保存命令审计记录
			if err := cmdAuditService.SaveCommands(rec.ID, sessionID, rec.StartTime.Time, commands); err != nil {
				// 命令审计保存失败不影响录制记录更新
				fmt.Printf("保存命令审计失败: %v\n", err)
			}

			// 评估整体风险等级
			riskLevel, err := cmdAuditService.GetHighestRiskLevel(rec.ID)
			if err == nil && riskLevel > 0 {
				// 更新录制记录的风险等级
				info.RiskLevel = riskLevel
			}

			// 检查是否有敏感命令
			sensitiveCount, err := cmdAuditService.GetSensitiveCommandCount(rec.ID)
			if err == nil && sensitiveCount > 0 {
				info.HasSensitiveCmd = true
			}
		}
	}

	// 构建更新数据
	updates := map[string]interface{}{
		"file_path":     info.FilePath,
		"file_size":     info.FileSize,
		"duration":      info.Duration,
		"input_count":   info.InputCount,
		"output_count":  info.OutputCount,
		"resize_count":  info.ResizeCount,
		"command_count": info.CommandCount,
		"end_time":      util.HTime{Time: time.Now()},
		"status":        model.SessionStatusCompleted,
		"update_time":   util.HTime{Time: time.Now()},
	}

	// 如果有风险评估结果，也更新
	if info.RiskLevel > 0 {
		updates["risk_level"] = info.RiskLevel
	}
	if info.HasSensitiveCmd {
		updates["has_sensitive_cmd"] = true
	}

	// 更新数据库
	if err := dao.UpdateSessionRecordingBySessionID(sessionID, updates); err != nil {
		return fmt.Errorf("更新录制记录失败: %w", err)
	}

	return nil
}

// FailRecording 标记录制失败
func (s *SessionRecordingService) FailRecording(sessionID string, errMsg string) error {
	updates := map[string]interface{}{
		"status":      model.SessionStatusFailed,
		"error_msg":   errMsg,
		"end_time":    util.HTime{Time: time.Now()},
		"update_time": util.HTime{Time: time.Now()},
	}

	return dao.UpdateSessionRecordingBySessionID(sessionID, updates)
}

// GetRecordingByID 根据ID获取录制记录
func (s *SessionRecordingService) GetRecordingByID(id uint) (*model.SysSessionRecording, error) {
	return dao.GetSessionRecordingByID(id)
}

// GetRecordingList 分页获取录制列表
func (s *SessionRecordingService) GetRecordingList(dto model.SessionRecordingListDto) ([]model.SysSessionRecording, int64, error) {
	return dao.GetSessionRecordingList(dto)
}

// GetRecordingContent 获取录制文件内容
func (s *SessionRecordingService) GetRecordingContent(id uint) (io.ReadCloser, error) {
	// 获取录制记录
	rec, err := dao.GetSessionRecordingByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取录制记录失败: %w", err)
	}

	// 从存储加载文件
	if rec.FilePath != "" {
		file, err := os.Open(rec.FilePath)
		if err != nil {
			return nil, fmt.Errorf("open recording file failed: %w", err)
		}
		return file, nil
	}

	reader, err := s.storage.Load(rec.SessionID)
	if err != nil {
		return nil, fmt.Errorf("加载录制文件失败: %w", err)
	}

	return reader, nil
}

// DeleteRecording 删除录制记录（包括文件）
func (s *SessionRecordingService) DeleteRecording(id uint) error {
	// 获取录制记录
	rec, err := dao.GetSessionRecordingByID(id)
	if err != nil {
		return fmt.Errorf("获取录制记录失败: %w", err)
	}

	// 删除文件
	if rec.FilePath != "" {
		_ = os.Remove(rec.FilePath)
	} else if err := s.storage.Delete(rec.SessionID); err != nil {
		// 文件删除失败，记录日志但继续删除数据库记录
		fmt.Printf("删除录制文件失败: %v\n", err)
	}

	// 删除数据库记录
	if err := dao.DeleteSessionRecordingByID(id); err != nil {
		return fmt.Errorf("删除录制记录失败: %w", err)
	}

	return nil
}

// BatchDeleteRecording 批量删除录制记录
func (s *SessionRecordingService) BatchDeleteRecording(ids []uint) error {
	// 逐个删除（确保文件也被删除）
	for _, id := range ids {
		if err := s.DeleteRecording(id); err != nil {
			return err
		}
	}
	return nil
}

// DownloadRecording 下载录制文件
func (s *SessionRecordingService) DownloadRecording(id uint) (string, *os.File, error) {
	// 获取录制记录
	rec, err := dao.GetSessionRecordingByID(id)
	if err != nil {
		return "", nil, fmt.Errorf("获取录制记录失败: %w", err)
	}

	// 生成下载文件名
	filename := fmt.Sprintf("session_%s_%s@%s.cast",
		rec.StartTime.Format("20060102_150405"),
		rec.SSHUser,
		rec.HostIP,
	)

	// 打开文件
	file, err := os.Open(rec.FilePath)
	if err != nil {
		return "", nil, fmt.Errorf("打开录制文件失败: %w", err)
	}

	return filename, file, nil
}
