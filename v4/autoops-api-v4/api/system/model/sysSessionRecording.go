// SSH会话录制相关模型
// 用于记录和管理SSH会话的录制信息
package model

import "dodevops-api/common/util"

// SysSessionRecording SSH会话录制记录表
type SysSessionRecording struct {
	ID        uint   `gorm:"primaryKey;comment:'主键ID'" json:"id"`
	SessionID string `gorm:"type:varchar(64);uniqueIndex;not null;comment:'会话唯一标识'" json:"sessionId"` // 会话唯一标识

	// 用户信息
	AdminID  uint   `gorm:"not null;comment:'操作用户ID'" json:"adminId"`               // 操作用户ID
	Username string `gorm:"type:varchar(64);not null;comment:'用户名'" json:"username"` // 用户名

	// 主机信息
	HostID   uint   `gorm:"not null;index;comment:'目标主机ID'" json:"hostId"`          // 目标主机ID
	HostName string `gorm:"type:varchar(128);not null;comment:'主机名称'" json:"hostName"` // 主机名称
	HostIP   string `gorm:"type:varchar(64);not null;index;comment:'主机IP'" json:"hostIp"` // 主机IP
	SSHUser  string `gorm:"type:varchar(64);not null;comment:'SSH登录用户'" json:"sshUser"` // SSH登录用户

	// 会话元数据
	StartTime      util.HTime `gorm:"not null;index;comment:'会话开始时间'" json:"startTime"`   // 会话开始时间
	EndTime        util.HTime `gorm:"comment:'会话结束时间'" json:"endTime"`                     // 会话结束时间
	Duration       int        `gorm:"comment:'会话时长(秒)'" json:"duration"`                   // 会话时长(秒)
	TerminalWidth  int        `gorm:"default:80;comment:'终端宽度'" json:"terminalWidth"`      // 终端宽度
	TerminalHeight int        `gorm:"default:24;comment:'终端高度'" json:"terminalHeight"`     // 终端高度

	// 录制文件信息
	FilePath    string `gorm:"type:varchar(512);not null;comment:'录制文件路径'" json:"filePath"` // 录制文件路径
	FileSize    int64  `gorm:"comment:'文件大小(字节)'" json:"fileSize"`                           // 文件大小(字节)
	StorageType int    `gorm:"default:1;comment:'存储类型:1-本地 2-OSS'" json:"storageType"`    // 存储类型
	OssKey      string `gorm:"type:varchar(512);comment:'OSS对象key'" json:"ossKey"`          // OSS对象key

	// 统计信息
	InputCount   int `gorm:"default:0;comment:'输入事件数量'" json:"inputCount"`       // 输入事件数量
	OutputCount  int `gorm:"default:0;comment:'输出事件数量'" json:"outputCount"`      // 输出事件数量
	ResizeCount  int `gorm:"default:0;comment:'窗口调整次数'" json:"resizeCount"`      // 窗口调整次数
	CommandCount int `gorm:"default:0;comment:'命令执行次数(估算)'" json:"commandCount"` // 命令执行次数(估算)

	// 审计信息
	ClientIP        string `gorm:"type:varchar(64);comment:'客户端IP'" json:"clientIp"`                 // 客户端IP
	UserAgent       string `gorm:"type:varchar(512);comment:'用户代理'" json:"userAgent"`                // 用户代理
	RiskLevel       int    `gorm:"default:0;index;comment:'风险等级:0-正常 1-可疑 2-高危'" json:"riskLevel"` // 风险等级
	HasSensitiveCmd bool   `gorm:"default:false;comment:'是否包含敏感命令'" json:"hasSensitiveCmd"`         // 是否包含敏感命令

	// 状态管理
	Status   int    `gorm:"default:1;index;comment:'状态:1-录制中 2-已完成 3-异常终止'" json:"status"` // 状态
	ErrorMsg string `gorm:"type:varchar(512);comment:'错误信息'" json:"errorMsg"`                 // 错误信息

	CreateTime util.HTime `gorm:"not null;comment:'创建时间'" json:"createTime"` // 创建时间
	UpdateTime util.HTime `gorm:"comment:'更新时间'" json:"updateTime"`          // 更新时间
	DeleteTime util.HTime `gorm:"comment:'删除时间'" json:"deleteTime"`          // 删除时间
}

// 状态常量
const (
	SessionStatusRecording = 1 // 录制中
	SessionStatusCompleted = 2 // 已完成
	SessionStatusFailed    = 3 // 异常终止
)

// 风险等级常量
const (
	RiskLevelNormal   = 0 // 正常
	RiskLevelSuspect  = 1 // 可疑
	RiskLevelCritical = 2 // 高危
)

// 存储类型常量
const (
	StorageTypeLocal = 1 // 本地存储
	StorageTypeOSS   = 2 // 对象存储
)

func (SysSessionRecording) TableName() string {
	return "sys_session_recording"
}

// CreateSessionRecordingDto 创建会话录制DTO
type CreateSessionRecordingDto struct {
	SessionID      string `json:"sessionId" binding:"required"`      // 会话唯一标识
	AdminID        uint   `json:"adminId" binding:"required"`        // 操作用户ID
	Username       string `json:"username" binding:"required"`       // 用户名
	HostID         uint   `json:"hostId" binding:"required"`         // 目标主机ID
	HostName       string `json:"hostName" binding:"required"`       // 主机名称
	HostIP         string `json:"hostIp" binding:"required"`         // 主机IP
	SSHUser        string `json:"sshUser" binding:"required"`        // SSH登录用户
	TerminalWidth  int    `json:"terminalWidth"`                     // 终端宽度
	TerminalHeight int    `json:"terminalHeight"`                    // 终端高度
	ClientIP       string `json:"clientIp"`                          // 客户端IP
	UserAgent      string `json:"userAgent"`                         // 用户代理
}

// UpdateSessionRecordingDto 更新会话录制DTO
type UpdateSessionRecordingDto struct {
	SessionID       string `json:"sessionId" binding:"required"` // 会话唯一标识
	FilePath        string `json:"filePath"`                     // 录制文件路径
	FileSize        int64  `json:"fileSize"`                     // 文件大小
	Duration        int    `json:"duration"`                     // 会话时长
	InputCount      int    `json:"inputCount"`                   // 输入事件数量
	OutputCount     int    `json:"outputCount"`                  // 输出事件数量
	ResizeCount     int    `json:"resizeCount"`                  // 窗口调整次数
	CommandCount    int    `json:"commandCount"`                 // 命令执行次数
	Status          int    `json:"status"`                       // 状态
	RiskLevel       int    `json:"riskLevel"`                    // 风险等级
	HasSensitiveCmd bool   `json:"hasSensitiveCmd"`              // 是否包含敏感命令
	ErrorMsg        string `json:"errorMsg"`                     // 错误信息
}

// SessionRecordingListDto 会话录制列表查询DTO
type SessionRecordingListDto struct {
	Page      int    `form:"page" binding:"min=1"`             // 页码
	PageSize  int    `form:"pageSize" binding:"min=1,max=100"` // 每页数量
	AdminID   uint   `form:"adminId"`                          // 用户ID过滤
	HostID    uint   `form:"hostId"`                           // 主机ID过滤
	Status    int    `form:"status"`                           // 状态过滤
	RiskLevel int    `form:"riskLevel"`                        // 风险等级过滤
	Keyword   string `form:"keyword"`                          // 关键词搜索(主机名、IP、用户名)
	StartTime string `form:"startTime"`                        // 开始时间过滤
	EndTime   string `form:"endTime"`                          // 结束时间过滤
}

// SessionRecordingIdDto ID参数
type SessionRecordingIdDto struct {
	ID uint `json:"id" binding:"required"` // 录制记录ID
}

// SessionRecordingSessionIdDto SessionID参数
type SessionRecordingSessionIdDto struct {
	SessionID string `json:"sessionId" binding:"required"` // 会话ID
}

// BatchDeleteSessionRecordingDto 批量删除DTO
type BatchDeleteSessionRecordingDto struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // ID列表
}
