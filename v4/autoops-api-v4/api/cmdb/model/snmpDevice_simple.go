package model

import "dodevops-api/common/util"

// CmdbSnmpDevice SNMP设备资产表(简化版,仅存储基本信息)
type CmdbSnmpDevice struct {
	ID uint `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	// 基础信息(必填)
	DeviceName string `gorm:"column:device_name;type:varchar(100);comment:'设备名称';NOT NULL" json:"deviceName"`
	DeviceType string `gorm:"column:device_type;type:varchar(20);comment:'设备类型(switch/router/windows/linux)';NOT NULL" json:"deviceType"`
	IPAddress  string `gorm:"column:ip_address;type:varchar(50);comment:'设备IP地址';NOT NULL;uniqueIndex" json:"ipAddress"`

	// SNMP监控配置(可选,有默认值)
	SnmpVersion   string `gorm:"column:snmp_version;type:varchar(10);comment:'SNMP版本';default:'v2c'" json:"snmpVersion"`
	SnmpCommunity string `gorm:"column:snmp_community;type:varchar(100);comment:'Community字符串';default:'public'" json:"snmpCommunity"`
	SnmpPort      int    `gorm:"column:snmp_port;comment:'SNMP端口';default:161" json:"snmpPort"`

	// 远程访问配置(SSH/RDP通用) - 根据deviceType自动判断协议
	RemotePort     int    `gorm:"column:remote_port;comment:'远程端口(SSH默认22,RDP默认3389)'" json:"remotePort"`
	RemoteUsername string `gorm:"column:remote_username;type:varchar(100);comment:'远程登录用户名'" json:"remoteUsername"`
	RemotePassword string `gorm:"column:remote_password;type:varchar(255);comment:'远程登录密码'" json:"remotePassword"`
	RemoteDomain   string `gorm:"column:remote_domain;type:varchar(100);comment:'域名(仅Windows RDP需要)'" json:"remoteDomain"`

	// Windows Exporter 配置
	ExporterInstalled int    `gorm:"column:exporter_installed;comment:'Exporter状态(0=未安装,1=已安装运行中,2=安装中,3=已安装未运行,4=安装失败,5=启动中)';default:0" json:"exporterInstalled"`
	ExporterPort      int    `gorm:"column:exporter_port;comment:'Exporter监听端口';default:9182" json:"exporterPort"`
	ExporterVersion   string `gorm:"column:exporter_version;type:varchar(50);comment:'Exporter版本号'" json:"exporterVersion"`

	// 监控和业务信息
	MonitorEnabled int        `gorm:"column:monitor_enabled;comment:'是否启用监控(0=否,1=是)';default:1" json:"monitorEnabled"`
	Location       string     `gorm:"column:location;type:varchar(200);comment:'物理位置'" json:"location"`
	Owner          string     `gorm:"column:owner;type:varchar(100);comment:'负责人'" json:"owner"`
	Remark         string     `gorm:"column:remark;type:text;comment:'备注'" json:"remark"`
	Status         string     `gorm:"column:status;type:varchar(20);comment:'设备状态(online/offline)';default:'offline'" json:"status"`
	LastSeen       util.HTime `gorm:"column:last_seen;comment:'最后在线时间'" json:"lastSeen"`

	// 时间戳
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime util.HTime `gorm:"column:update_time;comment:'更新时间';NOT NULL" json:"updateTime"`
}

func (CmdbSnmpDevice) TableName() string {
	return "cmdb_snmp_devices"
}

// ============= DTOs =============

// CreateSnmpDeviceDto 创建SNMP设备DTO(极简)
type CreateSnmpDeviceDto struct {
	DeviceName    string `json:"deviceName" validate:"required"`                                       // 必填
	DeviceType    string `json:"deviceType" validate:"required,oneof=switch router windows linux"`    // 必填
	IPAddress     string `json:"ipAddress" validate:"required"`                                        // 必填
	SnmpVersion   string `json:"snmpVersion"`                                                          // 可选,默认v2c
	SnmpCommunity string `json:"snmpCommunity"`                                                        // 可选,默认public
	SnmpPort      int    `json:"snmpPort"`                                                             // 可选,默认161
	RemotePort    int    `json:"remotePort"`                                                           // 可选,SSH默认22/RDP默认3389
	RemoteUsername string `json:"remoteUsername"`                                                      // 可选
	RemotePassword string `json:"remotePassword"`                                                      // 可选
	RemoteDomain   string `json:"remoteDomain"`                                                        // 可选,仅Windows需要
	Location       string `json:"location"`                                                            // 可选
	Owner          string `json:"owner"`                                                               // 可选
	Remark         string `json:"remark"`                                                              // 可选
}

// UpdateSnmpDeviceDto 更新SNMP设备DTO
type UpdateSnmpDeviceDto struct {
	ID             uint   `json:"id" validate:"required"`
	DeviceName     string `json:"deviceName"`
	DeviceType     string `json:"deviceType"`
	IPAddress      string `json:"ipAddress"`
	SnmpVersion    string `json:"snmpVersion"`
	SnmpCommunity  string `json:"snmpCommunity"`
	SnmpPort       int    `json:"snmpPort"`
	RemotePort     int    `json:"remotePort"`
	RemoteUsername string `json:"remoteUsername"`
	RemotePassword string `json:"remotePassword"`
	RemoteDomain   string `json:"remoteDomain"`
	MonitorEnabled int    `json:"monitorEnabled" validate:"omitempty,oneof=0 1"`
	Location       string `json:"location"`
	Owner          string `json:"owner"`
	Remark         string `json:"remark"`
}

// SnmpDeviceQueryDto SNMP设备查询DTO
type SnmpDeviceQueryDto struct {
	DeviceType     string `json:"deviceType" form:"deviceType"`
	Status         string `json:"status" form:"status"`
	MonitorEnabled int    `json:"monitorEnabled" form:"monitorEnabled"`
	Keyword        string `json:"keyword" form:"keyword"` // 支持名称、IP搜索
	Page           int    `json:"page" form:"page" validate:"min=1"`
	PageSize       int    `json:"pageSize" form:"pageSize" validate:"min=1,max=100"`
}

// SnmpDeviceListResult SNMP设备列表返回结果
type SnmpDeviceListResult struct {
	Items    []CmdbSnmpDevice `json:"items"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

// SnmpDeviceStatsDto SNMP设备统计信息DTO
type SnmpDeviceStatsDto struct {
	TotalCount       int `json:"totalCount"`       // 总数
	OnlineCount      int `json:"onlineCount"`      // 在线数
	OfflineCount     int `json:"offlineCount"`     // 离线数
	MaintenanceCount int `json:"maintenanceCount"` // 维护中数
	SwitchCount      int `json:"switchCount"`      // 交换机数量
	WindowsCount     int `json:"windowsCount"`     // Windows服务器数量
	RouterCount      int `json:"routerCount"`      // 路由器数量
	OtherCount       int `json:"otherCount"`       // 其他设备数量
}

// SnmpDeviceVO SNMP设备视图对象(用于API返回,兼容前端字段名)
type SnmpDeviceVO struct {
	ID             uint   `json:"id"`
	DeviceName     string `json:"deviceName"`
	DeviceType     string `json:"deviceType"`
	IPAddress      string `json:"ipAddress"`
	SnmpVersion    string `json:"snmpVersion"`
	SnmpCommunity  string `json:"snmpCommunity"`
	SshPort        int    `json:"sshPort"`        // SSH端口(DeviceType=linux时使用)
	SshUsername    string `json:"sshUsername"`    // SSH用户名
	SshPassword    string `json:"sshPassword"`    // SSH密码
	RdpPort           int    `json:"rdpPort"`           // RDP端口(DeviceType=windows时使用)
	RdpUsername       string `json:"rdpUsername"`       // RDP用户名
	RdpPassword       string `json:"rdpPassword"`       // RDP密码
	RdpDomain         string `json:"rdpDomain"`         // RDP域
	ExporterInstalled int    `json:"exporterInstalled"` // Exporter安装状态
	ExporterPort      int    `json:"exporterPort"`      // Exporter端口
	ExporterVersion   string `json:"exporterVersion"`   // Exporter版本
	MonitorEnabled    int    `json:"monitorEnabled"`
	Location       string `json:"location"`
	Owner          string `json:"owner"`
	Remark         string `json:"remark"`
	Status         string `json:"status"`
	LastSeen       string `json:"lastSeen"`
	CreateTime     string `json:"createTime"`
	UpdateTime     string `json:"updateTime"`
}

// ToVO 将Model转换为VO
func (d *CmdbSnmpDevice) ToVO() *SnmpDeviceVO {
	return &SnmpDeviceVO{
		ID:                d.ID,
		DeviceName:        d.DeviceName,
		DeviceType:        d.DeviceType,
		IPAddress:         d.IPAddress,
		SnmpVersion:       d.SnmpVersion,
		SnmpCommunity:     d.SnmpCommunity,
		SshPort:           d.RemotePort,
		SshUsername:       d.RemoteUsername,
		SshPassword:       d.RemotePassword,
		RdpPort:           d.RemotePort,
		RdpUsername:       d.RemoteUsername,
		RdpPassword:       d.RemotePassword,
		RdpDomain:         d.RemoteDomain,
		ExporterInstalled: d.ExporterInstalled,
		ExporterPort:      d.ExporterPort,
		ExporterVersion:   d.ExporterVersion,
		MonitorEnabled:    d.MonitorEnabled,
		Location:          d.Location,
		Owner:             d.Owner,
		Remark:            d.Remark,
		Status:            d.Status,
		LastSeen:          d.LastSeen.String(),
		CreateTime:        d.CreateTime.String(),
		UpdateTime:        d.UpdateTime.String(),
	}
}

// ============= 监控数据DTOs =============

// DeviceMonitorMetrics 设备完整监控指标
type DeviceMonitorMetrics struct {
	DeviceID     uint    `json:"deviceId"`     // 设备ID
	DeviceName   string  `json:"deviceName"`   // 设备名称
	DeviceType   string  `json:"deviceType"`   // 设备类型
	DeviceModel  string  `json:"deviceModel"`  // 设备型号(从sysDescr解析)
	IPAddress    string  `json:"ipAddress"`    // IP地址
	Status       string  `json:"status"`       // 设备状态(online/offline)

	// 基础监控指标
	CpuUsage     float64 `json:"cpuUsage"`     // CPU使用率(%)
	MemoryUsage  float64 `json:"memoryUsage"`  // 内存使用率(%)
	Temperature  float64 `json:"temperature"`  // 设备温度(摄氏度)
	UptimeDays   float64 `json:"uptimeDays"`   // 运行时间(天)

	// 端口流量统计
	PortInTop5    []PortTrafficStat `json:"portInTop5"`    // 入向流量Top5
	PortOutTop5   []PortTrafficStat `json:"portOutTop5"`   // 出向流量Top5
	IdlePortCount int               `json:"idlePortCount"` // 闲置端口数量(流量=0)
	TotalPortCount int              `json:"totalPortCount"` // 总端口数量

	// 历史趋势数据(最近24小时)
	CpuTrend         []TrendPoint `json:"cpuTrend"`         // CPU趋势
	MemoryTrend      []TrendPoint `json:"memoryTrend"`      // 内存趋势
	TemperatureTrend []TrendPoint `json:"temperatureTrend"` // 温度趋势

	// 时间戳
	LastUpdateTime string `json:"lastUpdateTime"` // 最后更新时间
}

// PortTrafficStat 端口流量统计
type PortTrafficStat struct {
	PortIndex      int     `json:"portIndex"`      // 端口索引
	PortName       string  `json:"portName"`       // 端口名称
	TrafficRateMbps float64 `json:"trafficRateMbps"` // 流量速率(Mbps)
	Utilization    float64 `json:"utilization"`    // 带宽利用率(%)
	PortSpeed      float64 `json:"portSpeed"`      // 端口速率(Mbps)
}

// TrendPoint 趋势数据点
type TrendPoint struct {
	Timestamp string  `json:"timestamp"` // 时间戳
	Value     float64 `json:"value"`     // 指标值
}

// DeviceResourceMetrics 设备资源使用率(轻量级API,仅CPU+内存)
type DeviceResourceMetrics struct {
	DeviceID    uint    `json:"deviceId"`    // 设备ID
	DeviceName  string  `json:"deviceName"`  // 设备名称
	IPAddress   string  `json:"ipAddress"`   // IP地址
	Status      string  `json:"status"`      // 设备状态(online/offline)
	CpuUsage    float64 `json:"cpuUsage"`    // CPU使用率(%)
	MemoryUsage float64 `json:"memoryUsage"` // 内存使用率(%)
	UpdateTime  string  `json:"updateTime"`  // 更新时间
}

// BatchResourceQueryDto 批量资源查询DTO
type BatchResourceQueryDto struct {
	DeviceIDs []uint `json:"deviceIds" validate:"required,min=1,max=100"` // 设备ID列表(最多100个)
}

// BatchResourceMetricsResult 批量资源查询结果
type BatchResourceMetricsResult struct {
	Items []DeviceResourceMetrics `json:"items"` // 设备列表
	Total int                     `json:"total"` // 总数
}
