// SSH命令审计相关模型
// 用于记录SSH会话中执行的所有命令
package model

import "dodevops-api/common/util"

// SysCommandAudit SSH命令审计记录表
type SysCommandAudit struct {
	ID uint `gorm:"primaryKey;comment:'主键ID'" json:"id"`

	// 关联信息
	RecordingID uint   `gorm:"not null;index;comment:'录制记录ID'" json:"recordingId"`   // 关联的录制记录ID
	SessionID   string `gorm:"type:varchar(64);index;not null;comment:'会话ID'" json:"sessionId"` // 会话ID

	// 命令信息
	Command   string  `gorm:"type:text;not null;comment:'执行的命令'" json:"command"`       // 执行的命令
	Timestamp float64 `gorm:"not null;comment:'相对时间戳(秒)'" json:"timestamp"`           // 相对录制开始的时间戳
	Sequence  int     `gorm:"not null;comment:'命令序号'" json:"sequence"`                 // 命令执行序号

	// 审计信息
	IsSensitive bool   `gorm:"default:false;index;comment:'是否为敏感命令'" json:"isSensitive"` // 是否为敏感命令
	RiskLevel   int    `gorm:"default:0;index;comment:'风险等级:0-正常 1-可疑 2-高危'" json:"riskLevel"` // 风险等级
	RiskReason  string `gorm:"type:varchar(512);comment:'风险原因'" json:"riskReason"`      // 风险原因说明

	// 元数据
	ExecuteTime util.HTime `gorm:"not null;index;comment:'执行时间(绝对时间)'" json:"executeTime"` // 命令执行的绝对时间

	CreateTime util.HTime `gorm:"not null;comment:'创建时间'" json:"createTime"` // 记录创建时间
}

func (SysCommandAudit) TableName() string {
	return "sys_command_audit"
}

// CreateCommandAuditDto 创建命令审计记录DTO
type CreateCommandAuditDto struct {
	RecordingID uint    `json:"recordingId" binding:"required"` // 录制记录ID
	SessionID   string  `json:"sessionId" binding:"required"`   // 会话ID
	Command     string  `json:"command" binding:"required"`     // 执行的命令
	Timestamp   float64 `json:"timestamp" binding:"required"`   // 相对时间戳
	Sequence    int     `json:"sequence" binding:"required"`    // 命令序号
	IsSensitive bool    `json:"isSensitive"`                    // 是否敏感
	RiskLevel   int     `json:"riskLevel"`                      // 风险等级
	RiskReason  string  `json:"riskReason"`                     // 风险原因
	ExecuteTime string  `json:"executeTime"`                    // 执行时间
}

// CommandAuditListDto 命令审计列表查询DTO
type CommandAuditListDto struct {
	Page        int    `form:"page" binding:"min=1"`             // 页码
	PageSize    int    `form:"pageSize" binding:"min=1,max=200"` // 每页数量
	RecordingID uint   `form:"recordingId"`                      // 录制记录ID过滤
	SessionID   string `form:"sessionId"`                        // 会话ID过滤
	IsSensitive *bool  `form:"isSensitive"`                      // 是否敏感过滤
	RiskLevel   *int   `form:"riskLevel"`                        // 风险等级过滤
	Keyword     string `form:"keyword"`                          // 命令关键词搜索
}
