package model

type HostMetrics struct {
	CPUUsage     float64 `json:"cpuUsage"`     // CPU使用率百分比
	MemoryUsage  float64 `json:"memoryUsage"`  // 内存使用率百分比
	DiskUsage    float64 `json:"diskUsage"`    // 磁盘使用率百分比
	Load5        float64 `json:"load5"`        // 5分钟系统负载
	OnlineStatus int     `json:"onlineStatus"` // 在线状态(0:在线,1:离线)
}

type PrometheusQueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"` // 用于matrix类型
			Value  []interface{}     `json:"value"`  // 用于vector类型
		} `json:"result"`
	} `json:"data"`
}

// HostMetricHistory 主机指标历史数据
type HostMetricHistory struct {
	HostID   uint              `json:"hostId"`   // 主机ID
	Metric   string            `json:"metric"`   // 指标类型(cpu/memory/disk)
	TimeData []MetricDataPoint `json:"timeData"` // 时间序列数据
}

// MetricDataPoint 指标数据点
type MetricDataPoint struct {
	Timestamp int64   `json:"timestamp"` // 时间戳(秒)
	Value     float64 `json:"value"`     // 指标值
}

// AllMetricsHistory 主机所有指标历史数据
type AllMetricsHistory struct {
	HostID         uint              `json:"hostId"`         // 主机ID
	CPU            []MetricDataPoint `json:"cpu"`            // CPU使用率历史数据(%)
	CPUUser        []MetricDataPoint `json:"cpuUser"`        // CPU用户态使用率历史数据(%)
	CPUSystem      []MetricDataPoint `json:"cpuSystem"`      // CPU内核态使用率历史数据(%)
	Memory         []MetricDataPoint `json:"memory"`         // 内存使用率历史数据(%)
	MemoryUsed     []MetricDataPoint `json:"memoryUsed"`     // 已用内存历史数据(GB)
	MemoryFree     []MetricDataPoint `json:"memoryFree"`     // 可用内存历史数据(GB)
	Disk           []MetricDataPoint `json:"disk"`           // 磁盘使用率历史数据(%)
	DiskUsed       []MetricDataPoint `json:"diskUsed"`       // 磁盘已用空间历史数据(GB)
	DiskTotal      []MetricDataPoint `json:"diskTotal"`      // 磁盘总空间历史数据(GB)
	DiskReadKB     []MetricDataPoint `json:"diskReadKB"`     // 磁盘读取速率(KB/s)
	DiskWriteKB    []MetricDataPoint `json:"diskWriteKB"`    // 磁盘写入速率(KB/s)
	NetworkReceive []MetricDataPoint `json:"networkReceive"` // 网络接收速率(KB/s)
	NetworkSend    []MetricDataPoint `json:"networkSend"`    // 网络发送速率(KB/s)
	Load1min       []MetricDataPoint `json:"load1min"`       // 1分钟系统负载
	Load5min       []MetricDataPoint `json:"load5min"`       // 5分钟系统负载
	Load15min      []MetricDataPoint `json:"load15min"`      // 15分钟系统负载
	TotalProcesses []MetricDataPoint `json:"totalProcesses"` // 系统进程总数
	TCPConnections []MetricDataPoint `json:"tcpConnections"` // TCP连接数历史数据
	UDPConnections []MetricDataPoint `json:"udpConnections"` // UDP连接数历史数据
	OpenFiles      []MetricDataPoint `json:"openFiles"`      // 系统打开文件数量历史数据
}

// ProcessMetrics 进程级别监控指标
type ProcessMetrics struct {
	PID        uint              `json:"pid"`        // 进程ID
	Name       string            `json:"name"`       // 进程名称
	CPUPercent []MetricDataPoint `json:"cpuPercent"` // CPU使用百分比历史数据
	MemPercent []MetricDataPoint `json:"memPercent"` // 内存使用百分比历史数据
	Host       string            `json:"host"`       // 主机名
}

// ProcessInfo 进程基本信息
type ProcessInfo struct {
	PID        uint    `json:"pid"`        // 进程ID
	Name       string  `json:"name"`       // 进程名称
	CPUPercent float64 `json:"cpuPercent"` // CPU使用百分比
	MemPercent float64 `json:"memPercent"` // 内存使用百分比
	Host       string  `json:"host"`       // 主机名
}

// TopProcessesResult 进程TOP排行结果
type TopProcessesResult struct {
	HostID     uint          `json:"hostId"`     // 主机ID
	HostName   string        `json:"hostName"`   // 主机名称
	TopCPU     []ProcessInfo `json:"topCPU"`     // CPU使用率前5的进程
	TopMemory  []ProcessInfo `json:"topMemory"`  // 内存使用率前5的进程
	UpdateTime int64         `json:"updateTime"` // 更新时间戳
}

// PortInfo TCP端口信息
type PortInfo struct {
	Port       string  `json:"port"`       // 端口号
	PID        string  `json:"pid"`        // 进程ID
	Service    string  `json:"service"`    // 服务名称
	Status     int     `json:"status"`     // 监听状态(1:监听中,0:未监听)
	CPUUsage   float64 `json:"cpuUsage"`   // 进程CPU使用率百分比
	MemUsage   float64 `json:"memUsage"`   // 进程内存使用率百分比
}

// HostPortsResult 主机端口信息结果
type HostPortsResult struct {
	HostID     uint       `json:"hostId"`     // 主机ID
	HostName   string     `json:"hostName"`   // 主机名称
	Ports      []PortInfo `json:"ports"`      // 端口列表
	Total      int        `json:"total"`      // 端口总数
	UpdateTime int64      `json:"updateTime"` // 更新时间戳
}

// HostUsageRank 主机使用率排名
type HostUsageRank struct {
	HostID       uint    `json:"hostId"`       // 主机ID
	HostName     string  `json:"hostName"`     // 主机名称
	IP           string  `json:"ip"`           // IP地址
	AvgUsage     float64 `json:"avgUsage"`     // 平均使用率
	MaxUsage     float64 `json:"maxUsage"`     // 最大使用率
	CurrentUsage float64 `json:"currentUsage"` // 当前使用率
}

// TopUsageHostsResult 使用率TOP主机结果
type TopUsageHostsResult struct {
	TopCPU    []HostUsageRank `json:"topCpu"`    // CPU使用率前N的主机
	TopMemory []HostUsageRank `json:"topMemory"` // 内存使用率前N的主机
	TopDisk   []HostUsageRank `json:"topDisk"`   // 磁盘使用率前N的主机
	Period    string          `json:"period"`    // 统计周期（如"15天"）
	StartTime string          `json:"startTime"` // 开始时间
	EndTime   string          `json:"endTime"`   // 结束时间
}

// CompleteHostMetrics 完整的主机监控数据
type CompleteHostMetrics struct {
	HostID   uint                  `json:"hostId"`   // 主机ID
	HostName string                `json:"hostName"` // 主机名称
	IP       string                `json:"ip"`       // IP地址

	// CPU相关
	CPU         CPUMetrics    `json:"cpu"`         // CPU指标

	// 内存相关
	Memory      MemoryMetrics `json:"memory"`      // 内存指标

	// 磁盘相关
	Disk        DiskMetrics   `json:"disk"`        // 磁盘指标

	// 负载
	Load        LoadMetrics   `json:"load"`        // 负载指标

	// 连接数
	Connections ConnMetrics   `json:"connections"` // 连接数指标

	// 网络
	Network     NetworkMetrics `json:"network"`    // 网络指标

	// 服务
	Services    []ServiceInfo  `json:"services"`   // 服务列表

	// 系统运行时间
	Uptime        float64        `json:"uptime"`        // 系统运行时间(天)
	UptimeStr     string         `json:"uptimeStr"`     // 系统运行时间(格式化字符串)

	// 系统进程总数
	ProcessCount  int            `json:"processCount"`  // 系统进程总数

	// 文件描述符数量
	OpenFiles     int            `json:"openFiles"`     // 打开的文件描述符数量

	UpdateTime  int64          `json:"updateTime"` // 更新时间戳
}

// CPUMetrics CPU指标
type CPUMetrics struct {
	Total      float64 `json:"total"`      // CPU总使用率(%)
	User       float64 `json:"user"`       // 用户态CPU使用率(%)
	System     float64 `json:"system"`     // 内核态CPU使用率(%)
}

// MemoryMetrics 内存指标
type MemoryMetrics struct {
	Used       uint64  `json:"used"`       // 已使用内存(bytes)
	UsedStr    string  `json:"usedStr"`    // 已使用内存(格式化字符串)
	Total      uint64  `json:"total"`      // 总内存(bytes)
	TotalStr   string  `json:"totalStr"`   // 总内存(格式化字符串)
	Free       uint64  `json:"free"`       // 可用内存(bytes)
	FreeStr    string  `json:"freeStr"`    // 可用内存(格式化字符串)
	UsageRate  float64 `json:"usageRate"`  // 内存使用率(%)
	SwapUsed   uint64  `json:"swapUsed"`   // 交换内存使用量(bytes)
	SwapUsedStr string `json:"swapUsedStr"` // 交换内存使用量(格式化字符串)
}

// LoadMetrics 负载指标
type LoadMetrics struct {
	Load1       float64 `json:"load1"`       // 1分钟负载
	Load5       float64 `json:"load5"`       // 5分钟负载
	Load15      float64 `json:"load15"`      // 15分钟负载
	LoadRatio1  float64 `json:"loadRatio1"`  // 1分钟负载比(负载/CPU核数)
	LoadRatio5  float64 `json:"loadRatio5"`  // 5分钟负载比
	LoadRatio15 float64 `json:"loadRatio15"` // 15分钟负载比
	CPUCores    int     `json:"cpuCores"`    // CPU核心数
}

// ConnMetrics 连接数指标
type ConnMetrics struct {
	TCP   int `json:"tcp"`   // TCP连接数
	UDP   int `json:"udp"`   // UDP连接数
	Total int `json:"total"` // 总连接数
}

// DiskMetrics 磁盘指标
type DiskMetrics struct {
	UsageRate     float64               `json:"usageRate"`     // 磁盘使用率(%)
	Used          uint64                `json:"used"`          // 磁盘使用量(bytes)
	UsedStr       string                `json:"usedStr"`       // 磁盘使用量(格式化字符串)
	Free          uint64                `json:"free"`          // 磁盘剩余量(bytes)
	FreeStr       string                `json:"freeStr"`       // 磁盘剩余量(格式化字符串)
	Total         uint64                `json:"total"`         // 磁盘总量(bytes)
	TotalStr      string                `json:"totalStr"`      // 磁盘总量(格式化字符串)
	DiskThroughput DiskThroughputMetrics `json:"diskThroughput"` // 磁盘吞吐量
	DiskIOPS      DiskIOPSMetrics       `json:"diskIops"`       // 磁盘IOPS
	InodeUsage    float64               `json:"inodeUsage"`    // Inode使用率(%)
	InodeTotal    uint64                `json:"inodeTotal"`    // Inode总数
	InodeUsed     uint64                `json:"inodeUsed"`     // Inode已用数
	InodeFree     uint64                `json:"inodeFree"`     // Inode可用数
}

// DiskThroughputMetrics 磁盘吞吐量指标
type DiskThroughputMetrics struct {
	ReadKBPerSec  float64 `json:"readKbPerSec"`  // 读取速率(KB/s)
	WriteKBPerSec float64 `json:"writeKbPerSec"` // 写入速率(KB/s)
}

// DiskIOPSMetrics 磁盘IOPS指标
type DiskIOPSMetrics struct {
	ReadPerSec  float64 `json:"readPerSec"`  // 读IOPS(次/s)
	WritePerSec float64 `json:"writePerSec"` // 写IOPS(次/s)
}

// NetworkMetrics 网络指标
type NetworkMetrics struct {
	ReceiveKBPerSec float64 `json:"receiveKbPerSec"` // 接收速率(KB/s)
	SendKBPerSec    float64 `json:"sendKbPerSec"`    // 发送速率(KB/s)
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name       string  `json:"name"`       // 服务名称
	CPU        float64 `json:"cpu"`        // CPU使用率(%)
	Memory     float64 `json:"memory"`     // 内存使用率(%)
}
