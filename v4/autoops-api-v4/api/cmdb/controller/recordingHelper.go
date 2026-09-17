// SSH会话录制集成辅助函数
// 提供WebSSH连接时的录制功能集成
package controller

import (
	"fmt"
	"log"

	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/system/model"
	sysService "dodevops-api/api/system/service"
	"dodevops-api/common/util/recording"
	"dodevops-api/common/util/websocket"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RecordingHelper 录制辅助工具
type RecordingHelper struct {
	recordingService *sysService.SessionRecordingService
	storagePath      string
}

// NewRecordingHelper 创建录制辅助工具
func NewRecordingHelper(storagePath string) (*RecordingHelper, error) {
	service, err := sysService.NewSessionRecordingService(storagePath)
	if err != nil {
		return nil, err
	}

	return &RecordingHelper{
		recordingService: service,
		storagePath:      storagePath,
	}, nil
}

// WrapWithRecording 将普通WebSSH包装为带录制的WebSSH
func (h *RecordingHelper) WrapWithRecording(
	ctx *gin.Context,
	webSSH *websocket.WebSSH,
	hostID uint,
) (*websocket.RecordingWebSSH, error) {

	// 获取主机信息
	host, err := dao.NewCmdbHostSSHDao().GetHostSSHInfo(hostID)
	if err != nil {
		log.Printf("获取主机信息失败: %v", err)
		return nil, fmt.Errorf("获取主机信息失败: %w", err)
	}

	// 从JWT获取用户信息
	claims, err := jwt.GetAdmin(ctx)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		// 使用默认用户信息
		claims = &model.JwtAdmin{
			ID:       0,
			Username: "unknown",
		}
	}

	// 生成会话ID
	sessionID := recording.GenerateSessionID()

	// 创建录制元数据
	metadata := recording.RecordingMetadata{
		SessionID:      sessionID,
		UserID:         claims.ID,
		Username:       claims.Username,
		HostID:         hostID,
		HostName:       host.HostName,
		HostIP:         host.SSHIP,
		SSHUser:        host.SSHName,
		TerminalWidth:  80, // 默认值，Connect时会更新
		TerminalHeight: 24, // 默认值，Connect时会更新
		ClientIP:       ctx.ClientIP(),
		UserAgent:      ctx.Request.UserAgent(),
	}

	// 创建录制器
	recorder := h.recordingService.CreateRecorder(metadata)

	// 创建数据库记录
	dto := &model.CreateSessionRecordingDto{
		SessionID:      sessionID,
		AdminID:        claims.ID,
		Username:       claims.Username,
		HostID:         hostID,
		HostName:       host.HostName,
		HostIP:         host.SSHIP,
		SSHUser:        host.SSHName,
		TerminalWidth:  80,
		TerminalHeight: 24,
		ClientIP:       ctx.ClientIP(),
		UserAgent:      ctx.Request.UserAgent(),
	}

	if _, err := h.recordingService.StartRecording(dto); err != nil {
		log.Printf("创建录制记录失败: %v (继续执行)", err)
		// 不影响SSH连接
	}

	// 包装WebSSH
	recordingWS := websocket.NewRecordingWebSSH(webSSH, recorder)

	return recordingWS, nil
}

// OnRecordingComplete 录制完成回调
func (h *RecordingHelper) OnRecordingComplete(sessionID string, info *recording.RecordingInfo, recorder recording.Recorder, err error) {
	if err != nil {
		log.Printf("录制失败: %v", err)
		if err := h.recordingService.FailRecording(sessionID, err.Error()); err != nil {
			log.Printf("更新录制失败状态失败: %v", err)
		}
		return
	}

	// 更新录制记录（并保存命令审计）
	if err := h.recordingService.FinishRecording(sessionID, info, recorder); err != nil {
		log.Printf("更新录制完成状态失败: %v", err)
	} else {
		log.Printf("录制完成: 会话=%s, 文件=%s, 大小=%d字节, 时长=%d秒, 命令=%d条",
			sessionID, info.FilePath, info.FileSize, info.Duration, info.CommandCount)
	}
}
