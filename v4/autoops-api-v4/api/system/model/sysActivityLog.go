// 系统动态日志相关模型
package model

import "dodevops-api/common/util"

// SysActivityLog 系统动态日志表
type SysActivityLog struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ActivityType int        `gorm:"not null;comment:'动态类型：1=密钥同步，2=域名检查，3=服务器巡检，4=定时任务，5=其他'" json:"activityType"` // 动态类型
	Title        string     `gorm:"type:varchar(200);not null;comment:'动态标题'" json:"title"`                              // 动态标题
	Content      string     `gorm:"type:text;comment:'动态详细内容'" json:"content"`                                           // 动态详细内容
	Status       int        `gorm:"not null;default:1;comment:'状态：1=成功，2=失败，3=部分成功'" json:"status"`                     // 状态
	RelatedID    uint       `gorm:"comment:'关联ID（如同步任务ID、域名ID等）'" json:"relatedId"`                                    // 关联ID
	Summary      string     `gorm:"type:varchar(500);comment:'摘要信息'" json:"summary"`                                     // 摘要信息
	Duration     int64      `gorm:"comment:'执行耗时（秒）'" json:"duration"`                                                  // 执行耗时（秒）
	CreateTime   util.HTime `gorm:"not null;comment:'创建时间'" json:"createTime"`                                          // 创建时间
}

// 动态类型常量
const (
	ActivityTypeSyncKey      = 1 // 密钥同步
	ActivityTypeCheckDomain  = 2 // 域名检查
	ActivityTypeCheckServer  = 3 // 服务器巡检
	ActivityTypeScheduleTask = 4 // 定时任务
	ActivityTypeOther        = 5 // 其他
)

func (SysActivityLog) TableName() string {
	return "sys_activity_log"
}

// CreateActivityLogDto 创建动态日志DTO
type CreateActivityLogDto struct {
	ActivityType int    `json:"activityType" binding:"required"` // 动态类型
	Title        string `json:"title" binding:"required"`        // 动态标题
	Content      string `json:"content"`                         // 动态详细内容
	Status       int    `json:"status"`                          // 状态
	RelatedID    uint   `json:"relatedId"`                       // 关联ID
	Summary      string `json:"summary"`                         // 摘要信息
	Duration     int64  `json:"duration"`                        // 执行耗时（秒）
}
