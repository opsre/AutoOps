// 命令审计 服务层
package service

import (
	"fmt"
	"time"

	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/util"
	"dodevops-api/common/util/recording"
)

// CommandAuditService 命令审计服务
type CommandAuditService struct{}

// NewCommandAuditService 创建命令审计服务
func NewCommandAuditService() *CommandAuditService {
	return &CommandAuditService{}
}

// SaveCommands 批量保存命令审计记录
func (s *CommandAuditService) SaveCommands(recordingID uint, sessionID string, startTime time.Time, commands []recording.ParsedCommand) error {
	if len(commands) == 0 {
		return nil
	}

	// 转换为数据库模型
	audits := make([]*model.SysCommandAudit, 0, len(commands))

	for _, cmd := range commands {
		// 评估命令风险
		riskLevel, riskReason := recording.AssessCommandRisk(cmd.Command)
		isSensitive := recording.IsSensitiveCommand(cmd.Command)

		// 计算绝对执行时间
		executeTime := startTime.Add(time.Duration(cmd.Timestamp * float64(time.Second)))

		audit := &model.SysCommandAudit{
			RecordingID: recordingID,
			SessionID:   sessionID,
			Command:     cmd.Command,
			Timestamp:   cmd.Timestamp,
			Sequence:    cmd.Sequence,
			IsSensitive: isSensitive,
			RiskLevel:   riskLevel,
			RiskReason:  riskReason,
			ExecuteTime: util.HTime{Time: executeTime},
			CreateTime:  util.HTime{Time: time.Now()},
		}

		audits = append(audits, audit)
	}

	// 批量插入数据库
	if err := dao.BatchCreateCommandAudit(audits); err != nil {
		return fmt.Errorf("批量保存命令审计失败: %w", err)
	}

	return nil
}

// GetCommandList 获取命令审计列表
func (s *CommandAuditService) GetCommandList(dto *model.CommandAuditListDto) ([]*model.SysCommandAudit, int64, error) {
	return dao.GetCommandAuditList(dto)
}

// DeleteByRecordingID 删除指定录制记录的所有命令审计
func (s *CommandAuditService) DeleteByRecordingID(recordingID uint) error {
	return dao.DeleteCommandAuditByRecordingID(recordingID)
}

// DeleteByRecordingIDs 批量删除命令审计记录
func (s *CommandAuditService) DeleteByRecordingIDs(recordingIDs []uint) error {
	return dao.DeleteCommandAuditByRecordingIDs(recordingIDs)
}

// GetSensitiveCommandCount 获取敏感命令数量
func (s *CommandAuditService) GetSensitiveCommandCount(recordingID uint) (int64, error) {
	return dao.GetSensitiveCommandCount(recordingID)
}

// GetHighestRiskLevel 获取最高风险等级
func (s *CommandAuditService) GetHighestRiskLevel(recordingID uint) (int, error) {
	return dao.GetHighestRiskLevel(recordingID)
}
