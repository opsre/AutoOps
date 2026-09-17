package agent

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	psnet "github.com/shirou/gopsutil/v3/net"
)

var (
	// 保存上次IO计数器值
	lastDiskIO  map[string]disk.IOCountersStat
	lastNetIO   psnet.IOCountersStat
	lastCPUTime cpu.TimesStat // 保存上次CPU时间用于计算增量

	// 新增指标
	diskReadRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_read_kb_per_second",
			Help: "磁盘读取速率(KB/s)",
		},
		[]string{"instance", "device"},
	)

	diskWriteRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_write_kb_per_second",
			Help: "磁盘写入速率(KB/s)",
		},
		[]string{"instance", "device"},
	)

	totalProcesses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_total_processes",
			Help: "系统进程总数",
		},
		[]string{"instance"},
	)

	// Pushgateway状态指标
	pushSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pushgateway_success_total",
			Help: "成功推送到Pushgateway的次数",
		},
	)

	pushFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pushgateway_failures_total",
			Help: "推送到Pushgateway失败的次数",
		},
	)

	// 系统指标
	cpuTotalUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_usage_percent",
			Help: "CPU总使用百分比",
		},
		[]string{"instance"},
	)

	systemLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_load_average",
			Help: "系统负载平均值",
		},
		[]string{"instance", "period"},
	)

	memUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_usage_percent",
			Help: "系统内存使用百分比",
		},
		[]string{"instance"},
	)

	diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_usage_percent",
			Help: "磁盘使用率百分比",
		},
		[]string{"instance", "device", "mountpoint"},
	)

	netRecv = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_network_receive_kb_per_second",
			Help: "网络接收速率(KB/s)",
		},
		[]string{"instance"},
	)

	netSent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_network_send_kb_per_second",
			Help: "网络发送速率(KB/s)",
		},
		[]string{"instance"},
	)

	// 进程CPU占用监控
	processCPU = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_cpu_percent",
			Help: "TOP进程CPU使用率百分比",
		},
		[]string{"host", "name", "pid"},
	)

	// 进程内存占用监控
	processMem = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_memory_percent",
			Help: "TOP进程内存使用率百分比",
		},
		[]string{"host", "name", "pid"},
	)

	// 端点响应时间监控
	endpointResponseTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "endpoint_response_time_ms",
			Help: "HTTP端点响应时间(毫秒)",
		},
		[]string{"host", "url", "method"},
	)

	// 端点存活监控
	endpointStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "endpoint_status",
			Help: "HTTP端点存活状态(1=正常,0=异常)",
		},
		[]string{"host", "url", "method"},
	)

	// Ping响应时间监控
	pingResponseTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ping_response_time_ms",
			Help: "Ping网络连通性响应时间(毫秒)",
		},
		[]string{"host", "target"},
	)

	// Ping存活监控
	pingStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ping_status",
			Help: "Ping网络连通性存活状态(1=成功,0=失败)",
		},
		[]string{"host", "target"},
	)

	// TCP端口服务监控
	tcpPortStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tcp_port_listening",
			Help: "TCP端口监听状态(1=监听中,0=未监听)",
		},
		[]string{"instance", "port", "service", "pid"},
	)

	// TCP连接数
	tcpConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_tcp_connections",
			Help: "TCP连接总数",
		},
		[]string{"instance"},
	)

	// UDP连接数
	udpConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_udp_connections",
			Help: "UDP连接总数",
		},
		[]string{"instance"},
	)

	// 磁盘IOPS - 读
	diskReadIOPS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_read_iops",
			Help: "磁盘读IOPS(次/秒)",
		},
		[]string{"instance", "device"},
	)

	// 磁盘IOPS - 写
	diskWriteIOPS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_write_iops",
			Help: "磁盘写IOPS(次/秒)",
		},
		[]string{"instance", "device"},
	)

	// CPU用户态使用率
	cpuUserUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_user_percent",
			Help: "CPU用户态使用率百分比",
		},
		[]string{"instance"},
	)

	// CPU内核态使用率
	cpuSystemUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_system_percent",
			Help: "CPU内核态使用率百分比",
		},
		[]string{"instance"},
	)

	// 已用内存(GB)
	memUsedGB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_used_gb",
			Help: "已使用内存容量(GB)",
		},
		[]string{"instance"},
	)

	// 总内存(GB)
	memTotalGB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_total_gb",
			Help: "系统总内存容量(GB)",
		},
		[]string{"instance"},
	)

	// 可用内存(GB)
	memFreeGB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_free_gb",
			Help: "可用内存容量(GB)",
		},
		[]string{"instance"},
	)

	// 交换内存(GB)
	swapUsedGB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_swap_used_gb",
			Help: "交换内存使用量(GB)",
		},
		[]string{"instance"},
	)

	// 磁盘已用空间(GB)
	diskUsedBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_used_bytes",
			Help: "磁盘已用空间(GB) - 直接存储GB值",
		},
		[]string{"instance", "device", "mountpoint"},
	)

	// 磁盘总空间(GB)
	diskTotalBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_total_bytes",
			Help: "磁盘总空间(GB) - 直接存储GB值",
		},
		[]string{"instance", "device", "mountpoint"},
	)

	// 服务CPU使用率监控
	serviceCPUUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_cpu_usage_percent",
			Help: "服务CPU使用率百分比",
		},
		[]string{"instance", "service_name"},
	)

	// 服务内存使用率监控
	serviceMemoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_memory_usage_percent",
			Help: "服务内存使用率百分比",
		},
		[]string{"instance", "service_name"},
	)

	// Agent版本信息
	agentInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "agent_info",
			Help: "Agent版本信息(值固定为1)",
		},
		[]string{"instance", "version", "version_tag"},
	)

	// 系统运行天数
	systemUptimeDays = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_uptime_days",
			Help: "系统运行天数",
		},
		[]string{"instance"},
	)

	// 文件打开数量
	openFileDescriptors = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_open_file_descriptors",
			Help: "系统当前打开的文件描述符数量",
		},
		[]string{"instance"},
	)
)
