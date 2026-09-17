package model

// AlertSource 告警源配置模型
type AlertSource struct {
	ID         uint   `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Name       string `gorm:"column:name;varchar(64);comment:'告警源名称';NOT NULL" json:"name"`
	Type       int    `gorm:"column:type;comment:'告警源类型(1=flashduty,2=zabbix,3=n9e)';NOT NULL" json:"type"`
	KeyID      uint   `gorm:"column:key_id;comment:'云账号ID';default:0" json:"keyId"`
	HostID     uint   `gorm:"column:host_id;comment:'主机ID';default:0" json:"hostId"`
	AppKey     string `gorm:"column:app_key;varchar(255);comment:'应用密钥';NOT NULL" json:"appKey"`
	ApiBaseURL string `gorm:"column:api_base_url;varchar(255);comment:'API基础URL'" json:"apiBaseUrl"`
	Status     int    `gorm:"column:status;comment:'状态(0=禁用,1=启用)';default:1" json:"status"`
	Remark     string `gorm:"column:remark;type:text;comment:'备注'" json:"remark"`
	CreateTime int64  `gorm:"column:create_time;comment:'创建时间';autoCreateTime" json:"createTime"`
	UpdateTime int64  `gorm:"column:update_time;comment:'更新时间';autoUpdateTime" json:"updateTime"`
}

func (AlertSource) TableName() string {
	return "monitor_alert_source"
}

// AlertSourceType 告警源类型常量
const (
	AlertSourceTypeFlashDuty = 1 // FlashDuty
	AlertSourceTypeZabbix    = 2 // Zabbix
	AlertSourceTypeN9e       = 3 // N9e
)

// AlertStatsQueryDto 告警统计查询DTO
type AlertStatsQueryDto struct {
	SourceID          uint     `json:"sourceId" validate:"required"`          // 告警源ID
	StartTime         int64    `json:"startTime" validate:"required"`         // 开始时间(Unix时间戳-秒)
	EndTime           int64    `json:"endTime" validate:"required"`           // 结束时间(Unix时间戳-秒)
	CompareStartTime  int64    `json:"compareStartTime"`                      // 对比开始时间(Unix时间戳-秒)
	CompareEndTime    int64    `json:"compareEndTime"`                        // 对比结束时间(Unix时间戳-秒)
	TeamIDs           []int64  `json:"teamIds"`                               // 团队ID列表
	IsMyTeam          bool     `json:"isMyTeam"`                              // 是否只查询我的团队
	ChannelIDs        []int64  `json:"channelIds"`                            // 频道ID列表
	Severities        []string `json:"severities"`                            // 严重程度列表
	TimeZone          string   `json:"timeZone"`                              // 时区,默认Asia/Shanghai
	Query             string   `json:"query"`                                 // 查询关键字
	Labels            map[string]string `json:"labels"`                       // 标签过滤
	Fields            map[string]string `json:"fields"`                       // 字段过滤
	SecondsToCloseFrom int64  `json:"secondsToCloseFrom"`                    // 关闭时间下限(秒)
	SecondsToCloseTo   int64  `json:"secondsToCloseTo"`                      // 关闭时间上限(秒)
	SecondsToAckFrom   int64  `json:"secondsToAckFrom"`                      // 确认时间下限(秒)
	SecondsToAckTo     int64  `json:"secondsToAckTo"`                        // 确认时间上限(秒)
	OrderBy           string   `json:"orderby"`                               // 排序字段
	Asc               bool     `json:"asc"`                                   // 是否升序
	K                 int      `json:"k"`                                     // TopK数量,默认100
	Label             string   `json:"label" validate:"required"`             // 统计维度标签名称(如:check)
}

// AlertStatsItem 告警统计项
type AlertStatsItem struct {
	Label               string  `json:"label"`               // 标签值
	TotalAlertCnt       int     `json:"totalAlertCnt"`       // 告警总数
	TotalAlertEventCnt  int     `json:"totalAlertEventCnt"`  // 告警事件总数
	Percentage          float64 `json:"percentage"`          // 占比(%)
}

// AlertStatsCompareItem 告警统计对比项(含同环比)
type AlertStatsCompareItem struct {
	Label                    string  `json:"label"`                    // 标签值
	TotalAlertCnt            int     `json:"totalAlertCnt"`            // 当前周期告警总数
	TotalAlertEventCnt       int     `json:"totalAlertEventCnt"`       // 当前周期告警事件总数
	Percentage               float64 `json:"percentage"`               // 当前周期占比(%)
	CompareTotalAlertCnt     int     `json:"compareTotalAlertCnt"`     // 对比周期告警总数
	CompareTotalAlertEventCnt int    `json:"compareTotalAlertEventCnt"`// 对比周期告警事件总数
	ComparePercentage        float64 `json:"comparePercentage"`        // 对比周期占比(%)
	AlertCntChange           int     `json:"alertCntChange"`           // 告警总数变化量
	AlertEventCntChange      int     `json:"alertEventCntChange"`      // 事件总数变化量
	AlertCntChangeRate       float64 `json:"alertCntChangeRate"`       // 告警总数变化率(%)
	AlertEventCntChangeRate  float64 `json:"alertEventCntChangeRate"`  // 事件总数变化率(%)
}

// AlertStatsResult 告警统计结果
type AlertStatsResult struct {
	Items      []AlertStatsItem `json:"items"`      // 统计项列表
	Total      int              `json:"total"`      // 总数
	StartTime  int64            `json:"startTime"`  // 开始时间
	EndTime    int64            `json:"endTime"`    // 结束时间
	QueryLabel string           `json:"queryLabel"` // 查询的标签名称
}

// AlertStatsCompareResult 告警统计对比结果
type AlertStatsCompareResult struct {
	Items               []AlertStatsCompareItem `json:"items"`               // 统计对比项列表
	Total               int                     `json:"total"`               // 当前周期总数
	CompareTotal        int                     `json:"compareTotal"`        // 对比周期总数
	TotalChange         int                     `json:"totalChange"`         // 总数变化量
	TotalChangeRate     float64                 `json:"totalChangeRate"`     // 总数变化率(%)
	StartTime           int64                   `json:"startTime"`           // 当前周期开始时间
	EndTime             int64                   `json:"endTime"`             // 当前周期结束时间
	CompareStartTime    int64                   `json:"compareStartTime"`    // 对比周期开始时间
	CompareEndTime      int64                   `json:"compareEndTime"`      // 对比周期结束时间
	QueryLabel          string                  `json:"queryLabel"`          // 查询的标签名称
}

// FlashDutyTopkRequest FlashDuty TopK请求结构
type FlashDutyTopkRequest struct {
	StartTime          int64             `json:"start_time"`
	EndTime            int64             `json:"end_time"`
	TeamIDs            []int64           `json:"team_ids,omitempty"`
	IsMyTeam           bool              `json:"is_my_team"`
	ChannelIDs         []int64           `json:"channel_ids,omitempty"`
	Severities         []string          `json:"severities,omitempty"`
	TimeZone           string            `json:"time_zone"`
	Query              string            `json:"query"`
	Labels             map[string]string `json:"labels"`
	Fields             map[string]string `json:"fields"`
	SecondsToCloseFrom int64             `json:"seconds_to_close_from"`
	SecondsToCloseTo   int64             `json:"seconds_to_close_to"`
	SecondsToAckFrom   int64             `json:"seconds_to_ack_from"`
	SecondsToAckTo     int64             `json:"seconds_to_ack_to"`
	OrderBy            string            `json:"orderby"`
	Asc                bool              `json:"asc"`
	K                  int               `json:"k"`
	Label              string            `json:"label"`
}

// FlashDutyTopkResponse FlashDuty TopK响应结构
type FlashDutyTopkResponse struct {
	RequestID string `json:"request_id"`
	Data      struct {
		Items []struct {
			Label              string `json:"label"`
			TotalAlertCnt      int    `json:"total_alert_cnt"`
			TotalAlertEventCnt int    `json:"total_alert_event_cnt"`
		} `json:"items"`
	} `json:"data"`
}

// CreateAlertSourceDto 创建告警源DTO
type CreateAlertSourceDto struct {
	Name       string `json:"name" validate:"required"`       // 告警源名称
	Type       int    `json:"type" validate:"required,min=1,max=3"` // 告警源类型
	KeyID      uint   `json:"keyId"`                          // 云账号ID
	HostID     uint   `json:"hostId"`                         // 主机ID
	AppKey     string `json:"appKey" validate:"required"`     // 应用密钥
	ApiBaseURL string `json:"apiBaseUrl"`                     // API基础URL
	Status     int    `json:"status"`                         // 状态
	Remark     string `json:"remark"`                         // 备注
}

// UpdateAlertSourceDto 更新告警源DTO
type UpdateAlertSourceDto struct {
	ID         uint   `json:"id" validate:"required"`
	Name       string `json:"name"`
	Type       int    `json:"type" validate:"min=1,max=3"`
	KeyID      uint   `json:"keyId"`
	HostID     uint   `json:"hostId"`
	AppKey     string `json:"appKey"`
	ApiBaseURL string `json:"apiBaseUrl"`
	Status     int    `json:"status"`
	Remark     string `json:"remark"`
}

// AlertSourceQueryDto 告警源查询DTO
type AlertSourceQueryDto struct {
	Type     int    `json:"type"`     // 告警源类型
	Status   *int   `json:"status"`   // 状态(使用指针，支持nil表示不过滤)
	Keyword  string `json:"keyword"`  // 关键字(名称)
	PageNum  int    `json:"pageNum"`  // 页码
	PageSize int    `json:"pageSize"` // 页大小
}

// AlertNotificationQueryDto 告警通知统计查询DTO
type AlertNotificationQueryDto struct {
	SourceID   uint    `json:"sourceId" validate:"required"`   // 告警源ID
	ChannelIDs []int64 `json:"channelIds"`                     // 协作空间ID列表(可选，为空代表查全部)
	Step       string  `json:"step" validate:"required"`       // 统计周期(day/week/month)
	StartTime  int64   `json:"startTime" validate:"required"`  // 开始时间戳(秒)
	EndTime    int64   `json:"endTime" validate:"required"`    // 结束时间戳(秒)
}

// AlertNotificationItem 告警通知统计项
type AlertNotificationItem struct {
	Timestamp  int64 `json:"timestamp"`  // 时间戳
	SmsCount   int   `json:"smsCount"`   // 短信数量
	EmailCount int   `json:"emailCount"` // 邮件数量
	VoiceCount int   `json:"voiceCount"` // 语音数量
	TotalCount int   `json:"totalCount"` // 总数
}

// AlertNotificationResult 告警通知统计结果
type AlertNotificationResult struct {
	Items         []AlertNotificationItem `json:"items"`         // 统计项列表
	TotalSms      int                     `json:"totalSms"`      // 总短信数
	TotalEmail    int                     `json:"totalEmail"`    // 总邮件数
	TotalVoice    int                     `json:"totalVoice"`    // 总语音数
	GrandTotal    int                     `json:"grandTotal"`    // 总计
	Step          string                  `json:"step"`          // 统计周期
	StartTime     int64                   `json:"startTime"`     // 开始时间
	EndTime       int64                   `json:"endTime"`       // 结束时间
}

// FlashDutyNotificationRequest FlashDuty通知统计请求结构
type FlashDutyNotificationRequest struct {
	ChannelIDs []int64 `json:"channel_ids"`       // 协作空间列表
	Step       string  `json:"step"`              // 统计周期(day/week/month)
	StartTime  int64   `json:"start_time"`        // 开始时间戳
	EndTime    int64   `json:"end_time"`          // 结束时间戳
}

// FlashDutyNotificationResponse FlashDuty通知统计响应结构
type FlashDutyNotificationResponse struct {
	RequestID string `json:"request_id"`
	Data      struct {
		Items []struct {
			Timestamp  int64 `json:"ts"`         // 时间戳
			SmsCount   int   `json:"sms_cnt"`    // 短信数量
			EmailCount int   `json:"email_cnt"`  // 邮件数量
			VoiceCount int   `json:"voice_cnt"`  // 语音数量
		} `json:"items"`
	} `json:"data"`
}

// ChannelStatsQueryDto 协作空间统计查询DTO
type ChannelStatsQueryDto struct {
	SourceID             uint              `json:"sourceId" validate:"required"`     // 告警源ID
	StartTime            int64             `json:"startTime" validate:"required"`    // 开始时间戳(秒)
	EndTime              int64             `json:"endTime" validate:"required"`      // 结束时间戳(秒)
	TeamIDs              []int64           `json:"teamIds"`                          // 团队ID列表
	IsMyTeam             bool              `json:"isMyTeam"`                         // 是否只查询我的团队
	ChannelIDs           []int64           `json:"channelIds"`                       // 协作空间ID列表
	Severities           []string          `json:"severities"`                       // 严重程度列表
	TimeZone             string            `json:"timeZone"`                         // 时区
	Query                string            `json:"query"`                            // 故障标题(模糊匹配)
	Labels               map[string]string `json:"labels"`                           // 故障标签
	Fields               map[string]string `json:"fields"`                           // 自定义字段
	SecondsToCloseFrom   int64             `json:"secondsToCloseFrom"`               // 最短存活时间
	SecondsToCloseTo     int64             `json:"secondsToCloseTo"`                 // 最长存活时间
	SecondsToAckFrom     int64             `json:"secondsToAckFrom"`                 // 最短认领耗时
	SecondsToAckTo       int64             `json:"secondsToAckTo"`                   // 最长认领耗时
	AggregateUnit        string            `json:"aggregateUnit"`                    // 聚合粒度(day/week/month)
	SplitHours           bool              `json:"splitHours"`                       // 是否区分时间段
	OrderBy              string            `json:"orderBy"`                          // 排序字段
	Asc                  bool              `json:"asc"`                              // 是否升序
}

// ChannelStatsItem 协作空间统计项
type ChannelStatsItem struct {
	Timestamp                     int64   `json:"timestamp"`                     // 时间戳
	ChannelID                     int64   `json:"channelId"`                     // 协作空间ID
	ChannelName                   string  `json:"channelName"`                   // 协作空间名称
	TotalIncidentCnt              int     `json:"totalIncidentCnt"`              // 故障总数
	TotalIncidentsAcknowledged    int     `json:"totalIncidentsAcknowledged"`    // 已认领故障数
	TotalIncidentsClosed          int     `json:"totalIncidentsClosed"`          // 已关闭故障数
	TotalIncidentsAutoClosed      int     `json:"totalIncidentsAutoClosed"`      // 自动关闭故障数
	TotalIncidentsManuallyClosed  int     `json:"totalIncidentsManuallyClosed"`  // 手动关闭故障数
	TotalIncidentsTimeoutClosed   int     `json:"totalIncidentsTimeoutClosed"`   // 超时关闭故障数
	TotalIncidentsEscalated       int     `json:"totalIncidentsEscalated"`       // 已升级故障数
	TotalIncidentsManuallyEscalated int   `json:"totalIncidentsManuallyEscalated"` // 手动升级故障数
	TotalIncidentsTimeoutEscalated  int   `json:"totalIncidentsTimeoutEscalated"`  // 超时升级故障数
	TotalIncidentsReassigned      int     `json:"totalIncidentsReassigned"`      // 已重新分配故障数
	TotalInterruptions            int     `json:"totalInterruptions"`            // 打断次数
	TotalNotifications            int     `json:"totalNotifications"`            // 通知总量
	TotalEngagedSeconds           int64   `json:"totalEngagedSeconds"`           // 总参与时间(秒)
	TotalSecondsToAck             int64   `json:"totalSecondsToAck"`             // 总认领耗时(秒)
	TotalSecondsToClose           int64   `json:"totalSecondsToClose"`           // 总关闭耗时(秒)
	MeanSecondsToAck              float64 `json:"meanSecondsToAck"`              // 平均认领耗时(秒)
	MeanSecondsToClose            float64 `json:"meanSecondsToClose"`            // 平均关闭耗时(秒)
	NoiseReductionPct             float64 `json:"noiseReductionPct"`             // 噪声抑制率(%)
	AcknowledgementPct            float64 `json:"acknowledgementPct"`            // 认领率(%)
	AutoClosePct                  float64 `json:"autoClosePct"`                  // 自动关闭率(%)
	TotalAlertCnt                 int     `json:"totalAlertCnt"`                 // 告警总数
	TotalAlertEventCnt            int     `json:"totalAlertEventCnt"`            // 告警事件总数
}

// ChannelStatsResult 协作空间统计结果
type ChannelStatsResult struct {
	Items         []ChannelStatsItem `json:"items"`         // 统计项列表
	Total         int                `json:"total"`         // 总协作空间数
	StartTime     int64              `json:"startTime"`     // 开始时间
	EndTime       int64              `json:"endTime"`       // 结束时间
	AggregateUnit string             `json:"aggregateUnit"` // 聚合粒度
}

// FlashDutyChannelRequest FlashDuty协作空间统计请求结构
type FlashDutyChannelRequest struct {
	StartTime            int64             `json:"start_time"`                      // 开始时间戳
	EndTime              int64             `json:"end_time"`                        // 结束时间戳
	TeamIDs              []int64           `json:"team_ids,omitempty"`              // 团队ID列表
	IsMyTeam             bool              `json:"is_my_team"`                      // 是否只查询我的团队
	ChannelIDs           []int64           `json:"channel_ids,omitempty"`           // 协作空间ID列表
	Severities           []string          `json:"severities,omitempty"`            // 严重程度列表
	TimeZone             string            `json:"time_zone"`                       // 时区
	Query                string            `json:"query"`                           // 故障标题
	Labels               map[string]string `json:"labels"`                          // 故障标签
	Fields               map[string]string `json:"fields"`                          // 自定义字段
	SecondsToCloseFrom   int64             `json:"seconds_to_close_from,omitempty"` // 最短存活时间
	SecondsToCloseTo     int64             `json:"seconds_to_close_to,omitempty"`   // 最长存活时间
	SecondsToAckFrom     int64             `json:"seconds_to_ack_from,omitempty"`   // 最短认领耗时
	SecondsToAckTo       int64             `json:"seconds_to_ack_to,omitempty"`     // 最长认领耗时
	AggregateUnit        string            `json:"aggregate_unit,omitempty"`        // 聚合粒度
	SplitHours           bool              `json:"split_hours,omitempty"`           // 是否区分时间段
	OrderBy              string            `json:"order_by,omitempty"`              // 排序字段
	Asc                  bool              `json:"asc,omitempty"`                   // 是否升序
}

// FlashDutyChannelResponse FlashDuty协作空间统计响应结构
type FlashDutyChannelResponse struct {
	RequestID string `json:"request_id"`
	Data      struct {
		Items []struct {
			Timestamp                       int64   `json:"ts"`
			ChannelID                       int64   `json:"channel_id"`
			ChannelName                     string  `json:"channel_name"`
			TotalIncidentCnt                int     `json:"total_incident_cnt"`
			TotalIncidentsAcknowledged      int     `json:"total_incidents_acknowledged"`
			TotalIncidentsClosed            int     `json:"total_incidents_closed"`
			TotalIncidentsAutoClosed        int     `json:"total_incidents_auto_closed"`
			TotalIncidentsManuallyClosed    int     `json:"total_incidents_manually_closed"`
			TotalIncidentsTimeoutClosed     int     `json:"total_incidents_timeout_closed"`
			TotalIncidentsEscalated         int     `json:"total_incidents_escalated"`
			TotalIncidentsManuallyEscalated int     `json:"total_incidents_manually_escalated"`
			TotalIncidentsTimeoutEscalated  int     `json:"total_incidents_timeout_escalated"`
			TotalIncidentsReassigned        int     `json:"total_incidents_reassigned"`
			TotalInterruptions              int     `json:"total_interruptions"`
			TotalNotifications              int     `json:"total_notifications"`
			TotalEngagedSeconds             int64   `json:"total_engaged_seconds"`
			TotalSecondsToAck               int64   `json:"total_seconds_to_ack"`
			TotalSecondsToClose             int64   `json:"total_seconds_to_close"`
			MeanSecondsToAck                float64 `json:"mean_seconds_to_ack"`
			MeanSecondsToClose              float64 `json:"mean_seconds_to_close"`
			NoiseReductionPct               float64 `json:"noise_reduction_pct"`
			AcknowledgementPct              float64 `json:"acknowledgement_pct"`
			TotalAlertCnt                   int     `json:"total_alert_cnt"`
			TotalAlertEventCnt              int     `json:"total_alert_event_cnt"`
		} `json:"items"`
	} `json:"data"`
}

// ====================  Zabbix 告警统计相关 ====================

// ZabbixAlertTopItem Zabbix告警TOP项
type ZabbixAlertTopItem struct {
	Name       string  `json:"name"`       // 告警名称或主机名或主机群组名
	Count      int     `json:"count"`      // 告警数量
	Percentage float64 `json:"percentage"` // 占比(%)
}

// ZabbixAlertTopResult Zabbix告警TOP结果
type ZabbixAlertTopResult struct {
	Items     []ZabbixAlertTopItem `json:"items"`     // TOP列表
	Total     int                  `json:"total"`     // 总告警数
	StartTime int64                `json:"startTime"` // 开始时间
	EndTime   int64                `json:"endTime"`   // 结束时间
	TimeRange string               `json:"timeRange"` // 时间范围描述(today/yesterday/last7days/last30days)
}

// ZabbixNotificationStats Zabbix告警通知统计
type ZabbixNotificationStats struct {
	Items      []ZabbixNotificationItem `json:"items"`     // 每日通知统计列表
	TotalCount int                      `json:"totalCount"` // 总通知数
	StartTime  int64                    `json:"startTime"`  // 开始时间
	EndTime    int64                    `json:"endTime"`    // 结束时间
	TimeRange  string                   `json:"timeRange"`  // 时间范围描述
}

// ZabbixNotificationItem Zabbix每日通知统计项
type ZabbixNotificationItem struct {
	Date       string `json:"date"`       // 日期 (YYYY-MM-DD)
	Timestamp  int64  `json:"timestamp"`  // 时间戳
	Count      int    `json:"count"`      // 当天通知数量
	Percentage float64 `json:"percentage"` // 占总数的百分比
}

// ZabbixEvent Zabbix事件
type ZabbixEvent struct {
	EventID  string `json:"eventid"`  // 事件ID
	Source   string `json:"source"`   // 源
	Object   string `json:"object"`   // 对象
	ObjectID string `json:"objectid"` // 对象ID
	Clock    string `json:"clock"`    // 时间戳
	Value    string `json:"value"`    // 值(0=OK, 1=Problem)
	Severity string `json:"severity"` // 严重程度
	Name     string `json:"name"`     // 名称
}

// ZabbixTrigger Zabbix触发器
type ZabbixTrigger struct {
	TriggerID   string `json:"triggerid"`   // 触发器ID
	Description string `json:"description"` // 描述
	Priority    string `json:"priority"`    // 优先级
}

// ZabbixHost Zabbix主机
type ZabbixHost struct {
	HostID string `json:"hostid"` // 主机ID
	Host   string `json:"host"`   // 主机名
	Name   string `json:"name"`   // 主机显示名
}

// ZabbixHostGroup Zabbix主机群组
type ZabbixHostGroup struct {
	GroupID string `json:"groupid"` // 群组ID
	Name    string `json:"name"`    // 群组名称
}

