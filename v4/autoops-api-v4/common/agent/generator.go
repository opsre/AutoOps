package agent

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Agent版本信息在这里定义已经在其他文件中定义了，这里不需要重复

// findProjectRoot 查找包含go.mod的项目根目录
func findProjectRoot() string {
	// 获取当前执行文件的路径
	execPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "/tmp" // 如果都失败了，返回临时目录
	}

	// 从执行文件所在目录向上查找go.mod
	dir := execPath
	for i := 0; i < 10; i++ { // 最多向上查找10级目录
		dir = dir[:strings.LastIndex(dir, "/")]
		if _, err := os.Stat(dir + "/go.mod"); err == nil {
			return dir
		}
		if dir == "/" || dir == "" {
			break
		}
	}

	// 如果还是找不到，返回当前目录
	if pwd, err := os.Getwd(); err == nil {
		return pwd
	}
	return "/tmp"
}

// CreateAgentMainFile 创建独立可执行的agent main文件（完全独立，不依赖任何外部包）
// 新增 groupID 和 groupName 参数，用于在所有指标中注入分组标签
func CreateAgentMainFile(heartbeatURL, heartbeatToken, pushgatewayURL string, groupID uint, groupName string) string {
	log.Printf("CreateAgentMainFile: 开始生成agent代码，heartbeatURL=%s, groupID=%d, groupName=%s", heartbeatURL, groupID, groupName)

	// 设置默认值
	if heartbeatURL == "" {
		heartbeatURL = "http://127.0.0.1:8000/api/v1/monitor/agent/heartbeat"
	}
	if heartbeatToken == "" {
		heartbeatToken = "agent-heartbeat-token-2024"
	}
	if pushgatewayURL == "" {
		pushgatewayURL = "http://8.130.14.34:9091"
	}
	if groupName == "" {
		groupName = "default"
	}

	// 直接构建完整的agent代码
	mainCode := fmt.Sprintf(`package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// Agent版本信息
const (
	AgentVersion    = "2.0.0" // Agent版本号
	AgentVersionTag = "v2"    // Agent版本标签
	AgentBuildTime  = ""      // 构建时间(编译时可注入)
)

// 主机分组信息 - 编译时注入
const (
	HostGroupID   = %[4]d     // 主机分组ID
	HostGroupName = "%[5]s"   // 主机分组名称
)

// Config 结构体定义
type Config struct {
	Processes []struct {
		Name string
	}
	Endpoints []struct {
		URL    string
		Method string
	}
	Pushgateway struct {
		URL          string
		PushInterval string
		Job          string
		Instance     string
	}
	Heartbeat struct {
		ServerURL string
		Interval  string
		Token     string
	}
	CollectInterval string
}

// 声明全局变量
var (
	// CPU使用率监控
	cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_usage_percent",
			Help: "CPU使用率百分比",
		},
		[]string{"instance", "cpu"},
	)

	// 内存使用率监控
	memUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_usage_percent",
			Help: "内存使用率百分比",
		},
		[]string{"instance"},
	)

	// 总内存容量监控(GB)
	memTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_total_gb",
			Help: "系统总内存容量(GB)",
		},
		[]string{"instance"},
	)

	// 已用内存监控(GB)
	memUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_used_gb",
			Help: "已使用内存容量(GB)",
		},
		[]string{"instance"},
	)

	// 可用内存监控(GB)
	memFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_free_gb",
			Help: "可用内存容量(GB)",
		},
		[]string{"instance"},
	)

	// 磁盘使用率监控
	diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_usage_percent",
			Help: "磁盘使用率百分比",
		},
		[]string{"instance", "device", "mountpoint"},
	)

	// 磁盘inode使用率监控
	diskInodeUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_inode_usage_percent",
			Help: "磁盘inode使用率百分比",
		},
		[]string{"instance", "device", "mountpoint"},
	)

	// 磁盘inode总数
	diskInodeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_inode_total",
			Help: "磁盘inode总数",
		},
		[]string{"instance", "device", "mountpoint"},
	)

	// 磁盘inode已用数
	diskInodeUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_inode_used",
			Help: "磁盘inode已用数",
		},
		[]string{"instance", "device", "mountpoint"},
	)

	// 磁盘inode可用数
	diskInodeFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_inode_free",
			Help: "磁盘inode可用数",
		},
		[]string{"instance", "device", "mountpoint"},
	)

	// 系统负载监控 - 1分钟
	systemLoad1min = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_load_average_1min",
			Help: "系统负载平均值-1分钟",
		},
		[]string{"instance"},
	)

	// 系统负载监控 - 5分钟
	systemLoad5min = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_load_average_5min",
			Help: "系统负载平均值-5分钟",
		},
		[]string{"instance"},
	)

	// 系统负载监控 - 15分钟
	systemLoad15min = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_load_average_15min",
			Help: "系统负载平均值-15分钟",
		},
		[]string{"instance"},
	)

	// 系统负载比 - 1分钟 (负载/CPU核数)
	systemLoadRatio1min = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_load_ratio_1min",
			Help: "系统负载比-1分钟(负载/CPU核数),>1表示过载",
		},
		[]string{"instance"},
	)

	// 系统负载比 - 5分钟
	systemLoadRatio5min = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_load_ratio_5min",
			Help: "系统负载比-5分钟(负载/CPU核数),>1表示过载",
		},
		[]string{"instance"},
	)

	// 系统负载比 - 15分钟
	systemLoadRatio15min = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_load_ratio_15min",
			Help: "系统负载比-15分钟(负载/CPU核数),>1表示过载",
		},
		[]string{"instance"},
	)

	// CPU核数(用于负载比计算)
	cpuCoreCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_core_count",
			Help: "系统CPU核心数",
		},
		[]string{"instance"},
	)

	// 磁盘读取速度监控
	diskRead = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_read_kb_per_second",
			Help: "磁盘读取速度(KB/秒)",
		},
		[]string{"instance", "device"},
	)

	// 磁盘写入速度监控
	diskWrite = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_write_kb_per_second",
			Help: "磁盘写入速度(KB/秒)",
		},
		[]string{"instance", "device"},
	)

	// 进程总数监控
	totalProcesses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_total_processes",
			Help: "系统进程总数",
		},
		[]string{"instance"},
	)

	// 网络接收流量监控
	netRecv = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_network_receive_kb_per_second",
			Help: "网络接收流量速度(KB/秒)",
		},
		[]string{"instance"},
	)

	// 网络发送流量监控
	netSent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_network_send_kb_per_second",
			Help: "网络发送流量速度(KB/秒)",
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

	// TCP端口服务监控
	tcpPortStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tcp_port_listening",
			Help: "TCP端口监听状态(1=监听中,0=未监听)",
		},
		[]string{"instance", "port", "service", "pid"},
	)

	// 进程CPU使用率监控(按PID)
	processCPUUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_cpu_usage_percent",
			Help: "进程CPU使用率百分比(按PID)",
		},
		[]string{"instance", "pid", "name"},
	)

	// 进程内存使用率监控(按PID)
	processMemoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_memory_usage_percent",
			Help: "进程内存使用率百分比(按PID)",
		},
		[]string{"instance", "pid", "name"},
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


	// TCP连接数
	tcpConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_tcp_connections",
			Help: "TCP连接总数",
		},
		[]string{"instance"},
	)

	// V2新增指标 - UDP连接数
	udpConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_udp_connections",
			Help: "UDP连接总数",
		},
		[]string{"instance"},
	)

	// V2新增指标 - 磁盘读IOPS
	diskReadIOPS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_read_iops",
			Help: "磁盘读IOPS(次/秒)",
		},
		[]string{"instance", "device"},
	)

	// V2新增指标 - 磁盘写IOPS
	diskWriteIOPS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_write_iops",
			Help: "磁盘写IOPS(次/秒)",
		},
		[]string{"instance", "device"},
	)

	// V2新增指标 - CPU用户态使用率
	cpuUserUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_user_percent",
			Help: "CPU用户态使用率百分比",
		},
		[]string{"instance"},
	)

	// V2新增指标 - CPU内核态使用率
	cpuSystemUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cpu_system_percent",
			Help: "CPU内核态使用率百分比",
		},
		[]string{"instance"},
	)

	// 交换内存使用量(GB)
	swapUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_swap_used_gb",
			Help: "交换内存使用量(GB)",
		},
		[]string{"instance"},
	)

	// V2新增指标 - 磁盘已用空间(GB)
	diskUsedBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_used_bytes",
			Help: "磁盘已用空间(GB) - 直接存储GB值",
		},
		[]string{"instance", "device", "mountpoint"},
	)

	// V2新增指标 - 磁盘总空间(GB)
	diskTotalBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_disk_total_bytes",
			Help: "磁盘总空间(GB) - 直接存储GB值",
		},
		[]string{"instance", "device", "mountpoint"},
	)

	// V2新增指标 - 服务CPU使用率
	serviceCPUUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_cpu_usage_percent",
			Help: "服务CPU使用率百分比",
		},
		[]string{"instance", "service_name"},
	)

	// V2新增指标 - 服务内存使用率
	serviceMemoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_memory_usage_percent",
			Help: "服务内存使用率百分比",
		},
		[]string{"instance", "service_name"},
	)

	// V2新增指标 - Agent版本信息
	agentInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "agent_info",
			Help: "Agent版本信息(值固定为1)",
		},
		[]string{"instance", "version", "version_tag"},
	)

	// V2.1新增指标 - 已用内存(MB)
	memUsedMB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_used_mb",
			Help: "已使用内存容量(MB)",
		},
		[]string{"instance"},
	)

	// V2.1新增指标 - 总内存(MB)
	memTotalMB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_total_mb",
			Help: "系统总内存容量(MB)",
		},
		[]string{"instance"},
	)

	// V2.1新增指标 - 可用内存(MB)
	memFreeMB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_memory_free_mb",
			Help: "可用内存容量(MB)",
		},
		[]string{"instance"},
	)

	// V2.1新增指标 - 交换内存(MB)
	swapUsedMB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_swap_used_mb",
			Help: "交换内存使用量(MB)",
		},
		[]string{"instance"},
	)

	// V2.1新增指标 - 系统运行天数
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

	// 上一次读取的磁盘IO数据，用于计算速率
	lastDiskIO  disk.IOCountersStat
	lastNetIO   psnet.IOCountersStat
	lastCPUTime cpu.TimesStat // 保存上次CPU时间用于计算增量
)

// HeartbeatData 心跳数据结构
type HeartbeatData struct {
	PID      int    `+"`json:\"pid\"`"+`
	IP       string `+"`json:\"ip\"`"+`
	Hostname string `+"`json:\"hostname\"`"+`
	Port     int    `+"`json:\"port\"`"+`
	Token    string `+"`json:\"token\"`"+`
	Version    string `+"`json:\"version\"`"+`    // Agent版本号
	VersionTag string `+"`json:\"versionTag\"`"+` // Agent版本标签
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Printf("Error getting local IP: %%v", err)
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func sendHeartbeat(serverURL, token string) {
	pid := os.Getpid()
	localIP := getLocalIP()
	if localIP == "" {
		log.Println("Warning: Could not determine local IP address")
		return
	}

	// 获取hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	heartbeat := HeartbeatData{
		PID:        pid,
		IP:         localIP,
		Hostname:   hostname,
		Port:       12580,
		Token:      token,
		Version:    AgentVersion,
		VersionTag: AgentVersionTag,
	}

	jsonData, err := json.Marshal(heartbeat)
	if err != nil {
		log.Printf("Error marshaling heartbeat data: %%v", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending heartbeat: %%v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Printf("Heartbeat sent successfully - PID: %%d, IP: %%s, Port: %%d, Token: %%s", pid, localIP, 12580, token)
	} else {
		log.Printf("Heartbeat failed with status: %%d", resp.StatusCode)
	}
}

func collectSystemMetrics() {
	// 使用本地IP作为instance标识（用于Prometheus查询区分主机）
	instanceID := getLocalIP()
	if instanceID == "" {
		// 如果获取IP失败，使用hostname作为备选
		instanceID, err := os.Hostname()
		if err != nil || instanceID == "" {
			instanceID = "unknown"
		}
	}

	// 采集CPU指标（总体使用率）
	if cpuPercents, err := cpu.Percent(0, false); err == nil {
		if len(cpuPercents) > 0 {
			cpuUsage.WithLabelValues(instanceID, "total").Set(cpuPercents[0])
		}
	}

	// 采集内存指标
	if memStats, err := mem.VirtualMemory(); err == nil {
		memUsage.WithLabelValues(instanceID).Set(memStats.UsedPercent)
		// GB单位的内存指标（使用Total-Available计算真实已用内存，与free命令一致）
		realUsed := float64(memStats.Total - memStats.Available)
		memTotal.WithLabelValues(instanceID).Set(float64(memStats.Total) / 1024 / 1024 / 1024)
		memUsed.WithLabelValues(instanceID).Set(realUsed / 1024 / 1024 / 1024)
		memFree.WithLabelValues(instanceID).Set(float64(memStats.Available) / 1024 / 1024 / 1024)
	}

	// 采集磁盘指标
	if diskStats, err := disk.Partitions(false); err == nil {
		for _, partition := range diskStats {
			// 只采集根分区
			if partition.Mountpoint == "/" {
				if usage, err := disk.Usage(partition.Mountpoint); err == nil {
					// 磁盘空间使用率
					diskUsage.WithLabelValues(instanceID, partition.Device, partition.Mountpoint).Set(usage.UsedPercent)
					diskUsedBytes.WithLabelValues(instanceID, partition.Device, partition.Mountpoint).Set(float64(usage.Used) / 1024 / 1024 / 1024)
					diskTotalBytes.WithLabelValues(instanceID, partition.Device, partition.Mountpoint).Set(float64(usage.Total) / 1024 / 1024 / 1024)

					// inode使用率
					if usage.InodesTotal > 0 {
						inodeUsedPercent := float64(usage.InodesUsed) / float64(usage.InodesTotal) * 100
						diskInodeUsage.WithLabelValues(instanceID, partition.Device, partition.Mountpoint).Set(inodeUsedPercent)
						diskInodeTotal.WithLabelValues(instanceID, partition.Device, partition.Mountpoint).Set(float64(usage.InodesTotal))
						diskInodeUsed.WithLabelValues(instanceID, partition.Device, partition.Mountpoint).Set(float64(usage.InodesUsed))
						diskInodeFree.WithLabelValues(instanceID, partition.Device, partition.Mountpoint).Set(float64(usage.InodesFree))
					}
				}
			}
		}
	}

	// 采集负载平均值和负载比
	if loadStats, err := load.Avg(); err == nil {
		// 原始负载值
		systemLoad1min.WithLabelValues(instanceID).Set(loadStats.Load1)
		systemLoad5min.WithLabelValues(instanceID).Set(loadStats.Load5)
		systemLoad15min.WithLabelValues(instanceID).Set(loadStats.Load15)

		// 获取CPU核数
		cpuCount, err := cpu.Counts(true)
		if err != nil {
			log.Printf("获取CPU核数失败: %v", err)
		} else if cpuCount <= 0 {
			log.Printf("CPU核数异常: %d", cpuCount)
		} else {
			// CPU核数指标
			cpuCoreCount.WithLabelValues(instanceID).Set(float64(cpuCount))

			// 计算负载比 (负载/核数)
			ratio1min := loadStats.Load1 / float64(cpuCount)
			ratio5min := loadStats.Load5 / float64(cpuCount)
			ratio15min := loadStats.Load15 / float64(cpuCount)

			systemLoadRatio1min.WithLabelValues(instanceID).Set(ratio1min)
			systemLoadRatio5min.WithLabelValues(instanceID).Set(ratio5min)
			systemLoadRatio15min.WithLabelValues(instanceID).Set(ratio15min)

			log.Printf("负载比计算: Load1=%.2f, Load5=%.2f, Load15=%.2f, CPU核数=%d, Ratio1min=%.4f, Ratio5min=%.4f, Ratio15min=%.4f",
				loadStats.Load1, loadStats.Load5, loadStats.Load15, cpuCount, ratio1min, ratio5min, ratio15min)
		}
	} else {
		log.Printf("获取负载失败: %v", err)
	}

	// 采集磁盘IO指标
	const intervalSeconds = 10.0
	if diskIOStats, err := disk.IOCounters(); err == nil {
		var totalReadBytes, totalWriteBytes uint64
		var totalReadCount, totalWriteCount uint64
		for _, diskIO := range diskIOStats {
			totalReadBytes += diskIO.ReadBytes
			totalWriteBytes += diskIO.WriteBytes
			totalReadCount += diskIO.ReadCount
			totalWriteCount += diskIO.WriteCount
		}
		diskIO := disk.IOCountersStat{
			ReadBytes:  totalReadBytes,
			WriteBytes: totalWriteBytes,
			ReadCount:  totalReadCount,
			WriteCount: totalWriteCount,
		}

		if lastDiskIO.ReadBytes != 0 {
			// 计算吞吐量(KB/s)
			readRate := float64(diskIO.ReadBytes-lastDiskIO.ReadBytes) / 1024 / intervalSeconds
			writeRate := float64(diskIO.WriteBytes-lastDiskIO.WriteBytes) / 1024 / intervalSeconds
			diskRead.WithLabelValues(instanceID, "vda").Set(readRate)
			diskWrite.WithLabelValues(instanceID, "vda").Set(writeRate)

			// 计算IOPS(次/秒)
			readIOPS := float64(diskIO.ReadCount-lastDiskIO.ReadCount) / intervalSeconds
			writeIOPS := float64(diskIO.WriteCount-lastDiskIO.WriteCount) / intervalSeconds
			diskReadIOPS.WithLabelValues(instanceID, "vda").Set(readIOPS)
			diskWriteIOPS.WithLabelValues(instanceID, "vda").Set(writeIOPS)
		}
		lastDiskIO = diskIO
	}

	// 采集进程总数
	if processes, err := process.Processes(); err == nil {
		totalProcesses.WithLabelValues(instanceID).Set(float64(len(processes)))
	}

	// 采集网络指标
	if netStats, err := psnet.IOCounters(true); err == nil {
		if len(netStats) > 0 {
			var mainInterface *psnet.IOCountersStat
			for i := range netStats {
				name := netStats[i].Name
				if strings.HasPrefix(name, "docker") ||
					strings.HasPrefix(name, "br-") ||
					strings.HasPrefix(name, "veth") ||
					strings.HasPrefix(name, "lo") {
					continue
				}
				if netStats[i].BytesRecv > 0 || netStats[i].BytesSent > 0 {
					mainInterface = &netStats[i]
					break
				}
			}
			if mainInterface == nil {
				for i := range netStats {
					name := netStats[i].Name
					if !strings.HasPrefix(name, "docker") &&
						!strings.HasPrefix(name, "br-") &&
						!strings.HasPrefix(name, "veth") &&
						!strings.HasPrefix(name, "lo") {
						mainInterface = &netStats[i]
						break
					}
				}
			}
			if mainInterface == nil && len(netStats) > 0 {
				mainInterface = &netStats[0]
			}

			if lastNetIO.Name != "" {
				const intervalSeconds = 10.0
				recvRate := float64(mainInterface.BytesRecv-lastNetIO.BytesRecv) / 1024 / intervalSeconds
				sentRate := float64(mainInterface.BytesSent-lastNetIO.BytesSent) / 1024 / intervalSeconds
				netRecv.WithLabelValues(instanceID).Set(recvRate)
				netSent.WithLabelValues(instanceID).Set(sentRate)
			}
			lastNetIO = *mainInterface
		}
	}

	// 采集CPU用户态/内核态（使用增量计算）
	if times, err := cpu.Times(false); err == nil && len(times) > 0 {
		t := times[0]
		// 如果有上次的数据，计算增量
		if lastCPUTime.User > 0 {
			deltaUser := t.User - lastCPUTime.User
			deltaSystem := t.System - lastCPUTime.System
			deltaIdle := t.Idle - lastCPUTime.Idle
			deltaTotal := deltaUser + deltaSystem + deltaIdle +
				(t.Nice - lastCPUTime.Nice) +
				(t.Iowait - lastCPUTime.Iowait) +
				(t.Irq - lastCPUTime.Irq) +
				(t.Softirq - lastCPUTime.Softirq) +
				(t.Steal - lastCPUTime.Steal)

			if deltaTotal > 0 {
				userPercent := (deltaUser / deltaTotal) * 100
				systemPercent := (deltaSystem / deltaTotal) * 100
				cpuUserUsage.WithLabelValues(instanceID).Set(userPercent)
				cpuSystemUsage.WithLabelValues(instanceID).Set(systemPercent)
			}
		}
		// 保存当前值用于下次计算
		lastCPUTime = t
	}

	// 采集交换内存(GB)
	if swapInfo, err := mem.SwapMemory(); err == nil {
		swapUsed.WithLabelValues(instanceID).Set(float64(swapInfo.Used) / 1024 / 1024 / 1024)
	}

	// 采集系统运行时间（天）
	if content, err := exec.Command("cat", "/proc/uptime").Output(); err == nil {
		if fields := strings.Fields(string(content)); len(fields) > 0 {
			if uptime, err := strconv.ParseFloat(fields[0], 64); err == nil {
				// 只保留天数
				systemUptimeDays.WithLabelValues(instanceID).Set(uptime / 86400)
			}
		}
	}

	// V2新增 - 采集TCP/UDP连接数（使用ss命令统计，更准确）
	// 统计TCP连接数（排除表头）
	if tcpOut, err := exec.Command("bash", "-c", "ss -tan 2>/dev/null | tail -n +2 | wc -l").Output(); err == nil {
		if tcpCount, err := strconv.Atoi(strings.TrimSpace(string(tcpOut))); err == nil {
			tcpConnections.WithLabelValues(instanceID).Set(float64(tcpCount))
		}
	}

	// 统计UDP连接数（排除表头）
	if udpOut, err := exec.Command("bash", "-c", "ss -uan 2>/dev/null | tail -n +2 | wc -l").Output(); err == nil {
		if udpCount, err := strconv.Atoi(strings.TrimSpace(string(udpOut))); err == nil {
			udpConnections.WithLabelValues(instanceID).Set(float64(udpCount))
		}
	}

	// 采集文件打开数量
	if fileNrOut, err := exec.Command("bash", "-c", "cat /proc/sys/fs/file-nr | awk '{print $1}'").Output(); err == nil {
		if openFiles, err := strconv.Atoi(strings.TrimSpace(string(fileNrOut))); err == nil {
			openFileDescriptors.WithLabelValues(instanceID).Set(float64(openFiles))
		}
	}
}

func monitorProcesses(config *Config) {
	// 使用本地IP作为host标识
	hostIdentifier := getLocalIP()
	if hostIdentifier == "" {
		// 如果获取IP失败，使用hostname作为备选
		hostname, err := os.Hostname()
		if err != nil {
			hostIdentifier = "unknown"
		} else {
			hostIdentifier = hostname
		}
	}
	processes, _ := process.Processes()

	type procStats struct {
		name string
		pid  int32
		cpu  float64
		mem  float32
	}
	var allProcs []procStats

	for _, p := range processes {
		name, _ := p.Name()
		pid := p.Pid
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()

		allProcs = append(allProcs, procStats{
			name: name,
			pid:  pid,
			cpu:  cpuPercent,
			mem:  memPercent,
		})
	}

	// 按CPU使用率排序
	sort.Slice(allProcs, func(i, j int) bool {
		return allProcs[i].cpu > allProcs[j].cpu
	})

	processCPU.Reset()

	// 记录CPU使用率最高的5个进程
	for i := 0; i < 5 && i < len(allProcs); i++ {
		p := allProcs[i]
		processCPU.WithLabelValues(hostIdentifier, p.name, fmt.Sprint(p.pid)).Set(p.cpu)
	}

	// 按内存使用率排序
	sort.Slice(allProcs, func(i, j int) bool {
		return allProcs[i].mem > allProcs[j].mem
	})

	processMem.Reset()

	// 记录内存使用率最高的5个进程
	for i := 0; i < 5 && i < len(allProcs); i++ {
		p := allProcs[i]
		processMem.WithLabelValues(hostIdentifier, p.name, fmt.Sprint(p.pid)).Set(float64(p.mem))
	}
}

func monitorEndpoints(config *Config) {
	// 使用本地IP作为host标识
	hostIdentifier := getLocalIP()
	if hostIdentifier == "" {
		// 如果获取IP失败，使用hostname作为备选
		hostname, err := os.Hostname()
		if err != nil {
			hostIdentifier = "unknown"
		} else {
			hostIdentifier = hostname
		}
	}
	for _, endpoint := range config.Endpoints {
		start := time.Now()
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(endpoint.URL)
		duration := time.Since(start).Milliseconds()

		status := 0
		if err == nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				status = 1
			}
			resp.Body.Close()
		}

		endpointResponseTime.WithLabelValues(hostIdentifier, endpoint.URL, endpoint.Method).Set(float64(duration))
		endpointStatus.WithLabelValues(hostIdentifier, endpoint.URL, endpoint.Method).Set(float64(status))
	}
}

func monitorPing(target string) {
	// 使用本地IP作为host标识
	hostIdentifier := getLocalIP()
	if hostIdentifier == "" {
		// 如果获取IP失败，使用hostname作为备选
		hostname, err := os.Hostname()
		if err != nil {
			hostIdentifier = "unknown"
		} else {
			hostIdentifier = hostname
		}
	}
	start := time.Now()

	cmd := exec.Command("ping", "-c", "1", target)
	err := cmd.Run()
	duration := time.Since(start).Milliseconds()

	status := 0
	if err == nil {
		status = 1
	}

	pingResponseTime.WithLabelValues(hostIdentifier, target).Set(float64(duration))
	pingStatus.WithLabelValues(hostIdentifier, target).Set(float64(status))
}

func monitorTCPPorts() {
	// 使用本地IP作为instance标识
	hostIdentifier := getLocalIP()
	if hostIdentifier == "" {
		// 如果获取IP失败，使用hostname作为备选
		hostname, err := os.Hostname()
		if err != nil {
			hostIdentifier = "unknown"
		} else {
			hostIdentifier = hostname
		}
	}

	// 重置端口监控指标，清除已停止监听的端口
	tcpPortStatus.Reset()
	// 重置进程资源监控指标，清除已停止的进程
	processCPUUsage.Reset()
	processMemoryUsage.Reset()

	// 优先尝试 ss 命令
	cmd := exec.Command("ss", "-ntlp")
	output, err := cmd.Output()
	usingSS := true
	if err != nil {
		// ss 失败，尝试 netstat
		cmd = exec.Command("netstat", "-ntlp")
		output, err = cmd.Output()
		usingSS = false
		if err != nil {
			// 两个命令都失败，直接返回
			return
		}
	}

	lines := strings.Split(string(output), "\n")
	startLine := 1
	if !usingSS {
		startLine = 2 // netstat 需要跳过两行标题
	}

	for i := startLine; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)

		var localAddr, processInfo string
		var isListening bool

		if usingSS && len(fields) >= 6 && fields[0] == "LISTEN" {
			// ss 格式
			isListening = true
			localAddr = fields[3]
			processInfo = fields[5]
		} else if !usingSS && len(fields) >= 7 && fields[5] == "LISTEN" {
			// netstat 格式
			isListening = true
			localAddr = fields[3]
			processInfo = fields[6]
		}

		if !isListening {
			continue
		}

		if strings.Contains(localAddr, ":") {
			parts := strings.Split(localAddr, ":")
			if len(parts) >= 2 {
				port := parts[len(parts)-1]
				pid := ""
				service := ""

				// 解析进程信息
				if processInfo != "" && processInfo != "-" {
					if strings.Contains(processInfo, "pid=") {
						// ss 格式: users:(("nginx",pid=2642106,fd=6))
						// 提取服务名: 查找第一个引号对
						startQuote := strings.Index(processInfo, "\"")
						if startQuote != -1 {
							endQuote := strings.Index(processInfo[startQuote+1:], "\"")
							if endQuote != -1 {
								service = processInfo[startQuote+1 : startQuote+1+endQuote]
							}
						}
						// 提取 PID
						pidIdx := strings.Index(processInfo, "pid=")
						if pidIdx != -1 {
							pidStr := processInfo[pidIdx+4:]
							commaIdx := strings.Index(pidStr, ",")
							if commaIdx != -1 {
								pid = pidStr[:commaIdx]
							}
						}
					} else if strings.Contains(processInfo, "/") {
						// netstat 格式: 123/nginx
						parts := strings.Split(processInfo, "/")
						if len(parts) >= 2 {
							pid = parts[0]
							service = parts[1]
						}
					}
				}

				if service == "" {
					service = "port-" + port
				}

				tcpPortStatus.WithLabelValues(hostIdentifier, port, service, pid).Set(1)

				// 收集该进程的CPU和内存使用率
				if pid != "" {
					if pidInt, err := strconv.Atoi(pid); err == nil {
						if p, err := process.NewProcess(int32(pidInt)); err == nil {
							if cpuPercent, err := p.CPUPercent(); err == nil {
								processCPUUsage.WithLabelValues(hostIdentifier, pid, service).Set(cpuPercent)
							}
							if memPercent, err := p.MemoryPercent(); err == nil {
								processMemoryUsage.WithLabelValues(hostIdentifier, pid, service).Set(float64(memPercent))
							}
						}
					}
				}
			}
		}
	}
}

func monitorServices() {
	// 获取hostname作为instance标识
	hostname := getLocalIP()
	if hostname == "" {
		// 如果获取IP失败，使用hostname作为备选
		if h, err := os.Hostname(); err == nil {
			hostname = h
		} else {
			hostname = "unknown"
		}
	}

	// 排除的系统服务模式
	excludePattern := "systemd|systemd-resolve|systemd-network|sshd|chronyd|pmproxy|pmcd|pmlogger|pmie|docker-proxy"

	// 存储找到的服务
	servicesMap := make(map[string]bool)

	// 1. 从监听端口获取服务
	cmd := exec.Command("bash", "-c", "ss -ntpl 2>/dev/null || netstat -ntlp 2>/dev/null")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" || strings.Contains(line, "State") {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) < 6 {
				continue
			}

			// 提取进程信息(ss和netstat格式略有不同)
			var processInfo string
			for _, field := range fields {
				if strings.Contains(field, "users:") || strings.Contains(field, "/") {
					processInfo = field
					break
				}
			}

			if processInfo == "" {
				continue
			}

			// 解析服务名
			var serviceName string

			// ss格式: users:(("nginx",pid=123,fd=6))
			if strings.Contains(processInfo, "users:") {
				start := strings.Index(processInfo, "\"")
				if start != -1 {
					end := strings.Index(processInfo[start+1:], "\"")
					if end != -1 {
						serviceName = processInfo[start+1 : start+1+end]
					}
				}
			} else if strings.Contains(processInfo, "/") {
				// netstat格式: 123/nginx
				parts := strings.Split(processInfo, "/")
				if len(parts) >= 2 {
					serviceName = parts[1]
				}
			}

			// 跳过系统服务
			if serviceName != "" && !strings.Contains(excludePattern, serviceName) {
				// 特殊处理java/python/node
				if serviceName == "java" {
					servicesMap["java"] = true
				} else if strings.HasPrefix(serviceName, "python") {
					servicesMap["python"] = true
				} else if serviceName == "node" || serviceName == "nodejs" {
					servicesMap["node"] = true
				} else {
					servicesMap[serviceName] = true
				}
			}
		}
	}

	// 2. 检查docker/containerd进程
	dockerCmd := exec.Command("bash", "-c", "ps -ef | grep -E '[d]ockerd|[c]ontainerd' | grep -v 'containerd-shim' | awk '{print $8}'")
	dockerOutput, err := dockerCmd.Output()
	if err == nil {
		lines := strings.Split(string(dockerOutput), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				serviceName := strings.TrimPrefix(line, "/usr/bin/")
				serviceName = strings.TrimPrefix(serviceName, "/usr/local/bin/")
				if serviceName != "" {
					servicesMap[serviceName] = true
				}
			}
		}
	}

	// 清空旧的服务指标
	serviceCPUUsage.Reset()
	serviceMemoryUsage.Reset()

	// 3. 获取每个服务的资源使用情况
	for serviceName := range servicesMap {
		var cpuUsage, memUsage float64

		switch serviceName {
		case "java":
			// 查找所有java进程及其jar/war文件
			javaCmd := exec.Command("bash", "-c", "ps -ef | grep '[j]ava' | awk '{for(i=8;i<=NF;i++) printf \"%%s \", $i; print \"\"}'")
			javaOut, err := javaCmd.Output()
			if err == nil {
				// 查找jar/war文件
				jarCmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%%s' | grep -oP '(\\S+\\.jar|\\S+\\.war)' | sort -u", string(javaOut)))
				jarOut, err := jarCmd.Output()
				if err == nil {
					jars := strings.Split(strings.TrimSpace(string(jarOut)), "\n")
					for _, jar := range jars {
						if jar != "" {
							jarBase := jar[strings.LastIndex(jar, "/")+1:]
							cpu, mem := getServiceResourceUsage("java: " + jarBase)
							serviceCPUUsage.WithLabelValues(hostname, "java: "+jarBase).Set(cpu)
							serviceMemoryUsage.WithLabelValues(hostname, "java: "+jarBase).Set(mem)
						}
					}
				}
			}

		case "python":
			// 查找所有python进程及其.py文件
			pythonCmd := exec.Command("bash", "-c", "ps -ef | grep '[p]ython' | awk '{for(i=8;i<=NF;i++) printf \"%%s \", $i; print \"\"}'")
			pythonOut, err := pythonCmd.Output()
			if err == nil {
				pyCmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%%s' | grep -oP '\\s(\\S+\\.py)\\b' | sed 's/^\\s*//' | sort -u", string(pythonOut)))
				pyOut, err := pyCmd.Output()
				if err == nil {
					scripts := strings.Split(strings.TrimSpace(string(pyOut)), "\n")
					for _, script := range scripts {
						if script != "" {
							scriptBase := script[strings.LastIndex(script, "/")+1:]
							cpu, mem := getServiceResourceUsage("python: " + scriptBase)
							serviceCPUUsage.WithLabelValues(hostname, "python: "+scriptBase).Set(cpu)
							serviceMemoryUsage.WithLabelValues(hostname, "python: "+scriptBase).Set(mem)
						}
					}
				}
			}

		case "node":
			// 查找所有node进程及其.js文件
			nodeCmd := exec.Command("bash", "-c", "ps -ef | grep '[n]ode' | awk '{for(i=8;i<=NF;i++) printf \"%%s \", $i; print \"\"}'")
			nodeOut, err := nodeCmd.Output()
			if err == nil {
				jsCmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%%s' | grep -oP '\\s(\\S+\\.js)\\b' | sed 's/^\\s*//' | sort -u", string(nodeOut)))
				jsOut, err := jsCmd.Output()
				if err == nil {
					scripts := strings.Split(strings.TrimSpace(string(jsOut)), "\n")
					for _, script := range scripts {
						if script != "" {
							scriptBase := script[strings.LastIndex(script, "/")+1:]
							cpu, mem := getServiceResourceUsage("node: " + scriptBase)
							serviceCPUUsage.WithLabelValues(hostname, "node: "+scriptBase).Set(cpu)
							serviceMemoryUsage.WithLabelValues(hostname, "node: "+scriptBase).Set(mem)
						}
					}
				}
			}

		default:
			// 其他服务直接查询
			cpuUsage, memUsage = getServiceResourceUsage(serviceName)
			serviceCPUUsage.WithLabelValues(hostname, serviceName).Set(cpuUsage)
			serviceMemoryUsage.WithLabelValues(hostname, serviceName).Set(memUsage)
		}
	}
}

// getServiceResourceUsage 获取服务的CPU和内存使用率
func getServiceResourceUsage(serviceName string) (cpu float64, mem float64) {
	var cmd *exec.Cmd

	if strings.HasPrefix(serviceName, "java: ") {
		jarName := strings.TrimPrefix(serviceName, "java: ")
		cmdStr := fmt.Sprintf("ps aux | grep 'java' | grep '%%s' | grep -v grep | awk '{cpu+=$3; mem+=$4} END {printf \"%%%%f %%%%f\", cpu, mem}'", jarName)
		log.Printf("[DEBUG] Java服务[%%s] 执行命令: %%s", serviceName, cmdStr)
		cmd = exec.Command("bash", "-c", cmdStr)
	} else if strings.HasPrefix(serviceName, "python: ") {
		scriptName := strings.TrimPrefix(serviceName, "python: ")
		cmdStr := fmt.Sprintf("ps aux | grep 'python' | grep '%%s' | grep -v grep | awk '{cpu+=$3; mem+=$4} END {printf \"%%%%f %%%%f\", cpu, mem}'", scriptName)
		log.Printf("[DEBUG] Python服务[%%s] 执行命令: %%s", serviceName, cmdStr)
		cmd = exec.Command("bash", "-c", cmdStr)
	} else if strings.HasPrefix(serviceName, "node: ") {
		scriptName := strings.TrimPrefix(serviceName, "node: ")
		cmdStr := fmt.Sprintf("ps aux | grep 'node' | grep '%%s' | grep -v grep | awk '{cpu+=$3; mem+=$4} END {printf \"%%%%f %%%%f\", cpu, mem}'", scriptName)
		log.Printf("[DEBUG] Node服务[%%s] 执行命令: %%s", serviceName, cmdStr)
		cmd = exec.Command("bash", "-c", cmdStr)
	} else {
		// 参考check_services.sh的实现：匹配进程名在第11列（可能是完整路径或进程名）
		// $11 ~ svc 匹配包含服务名的进程
		// $11 ~ "/"svc"$" 匹配路径末尾是服务名的进程（如 /usr/bin/containerd）
		cmdStr := fmt.Sprintf("ps aux | awk -v svc='%%s' '$11 ~ svc || $11 ~ \"/\"svc\"$\" {cpu+=$3; mem+=$4} END {if(cpu>0 || mem>0) printf \"%%%%f %%%%f\", cpu, mem; else print \"0 0\"}'", serviceName)
		log.Printf("[DEBUG] 普通服务[%%s] 执行命令: %%s", serviceName, cmdStr)
		cmd = exec.Command("bash", "-c", cmdStr)
	}

	output, err := cmd.Output()
	if err != nil {
		log.Printf("[DEBUG] 服务[%%s] 命令执行失败: %%v", serviceName, err)
		return 0, 0
	}

	result := strings.TrimSpace(string(output))
	log.Printf("[DEBUG] 服务[%%s] 命令输出: [%%s]", serviceName, result)

	if result == "" {
		log.Printf("[DEBUG] 服务[%%s] 输出为空", serviceName)
		return 0, 0
	}

	n, err := fmt.Sscanf(result, "%%f %%f", &cpu, &mem)
	if err != nil || n != 2 {
		log.Printf("[DEBUG] 服务[%%s] 解析失败: scanned=%%d, err=%%v, result=[%%s]", serviceName, n, err, result)
		return 0, 0
	}

	log.Printf("[DEBUG] 服务[%%s] 最终结果: CPU=%%.2f%%%%, MEM=%%.2f%%%%", serviceName, cpu, mem)
	return cpu, mem
}
func registerMetrics(registry *prometheus.Registry) {
	registry.MustRegister(cpuUsage)
	registry.MustRegister(memUsage)
	registry.MustRegister(memTotal)
	registry.MustRegister(memUsed)
	registry.MustRegister(memFree)
	registry.MustRegister(diskUsage)
	registry.MustRegister(diskInodeUsage)
	registry.MustRegister(diskInodeTotal)
	registry.MustRegister(diskInodeUsed)
	registry.MustRegister(diskInodeFree)
	registry.MustRegister(systemLoad1min)
	registry.MustRegister(systemLoad5min)
	registry.MustRegister(systemLoad15min)
	registry.MustRegister(systemLoadRatio1min)
	registry.MustRegister(systemLoadRatio5min)
	registry.MustRegister(systemLoadRatio15min)
	registry.MustRegister(cpuCoreCount)
	registry.MustRegister(diskRead)
	registry.MustRegister(diskWrite)
	registry.MustRegister(totalProcesses)
	registry.MustRegister(netRecv)
	registry.MustRegister(netSent)
	registry.MustRegister(processCPU)
	registry.MustRegister(processMem)
	registry.MustRegister(tcpPortStatus)
	registry.MustRegister(processCPUUsage)
	registry.MustRegister(processMemoryUsage)
	registry.MustRegister(endpointResponseTime)
	registry.MustRegister(endpointStatus)
	registry.MustRegister(pingResponseTime)
	registry.MustRegister(pingStatus)
	registry.MustRegister(tcpConnections)
	registry.MustRegister(udpConnections)
	registry.MustRegister(diskReadIOPS)
	registry.MustRegister(diskWriteIOPS)
	registry.MustRegister(cpuUserUsage)
	registry.MustRegister(cpuSystemUsage)
	registry.MustRegister(swapUsed)
	registry.MustRegister(diskUsedBytes)
	registry.MustRegister(diskTotalBytes)
	registry.MustRegister(serviceCPUUsage)
	registry.MustRegister(serviceMemoryUsage)
	registry.MustRegister(agentInfo)
	registry.MustRegister(systemUptimeDays)
	registry.MustRegister(openFileDescriptors)
}
func getEnvVar(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func StartAgentWithConfig(config *Config) error {
	// 创建自定义注册表，不包含Go运行时指标
	registry := prometheus.NewRegistry()
	registerMetrics(registry)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	go func() {
		log.Println("Agent starting on :12580")
		if err := http.ListenAndServe(":12580", nil); err != nil {
			log.Printf("HTTP server failed: %%v", err)
		}
	}()

	// 定期采集系统指标
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	// 定期推送到Pushgateway
	pushTicker := time.NewTicker(30 * time.Second)
	defer pushTicker.Stop()

	// 心跳
	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-ticker.C:
			collectSystemMetrics()
			monitorProcesses(config)
			monitorEndpoints(config)
			monitorPing("8.8.8.8")
			monitorTCPPorts()
			monitorServices()
		case <-pushTicker.C:
			if config.Pushgateway.URL != "" {
				pusher := push.New(config.Pushgateway.URL, config.Pushgateway.Job)
				// 使用hostname作为instance分组标签
				pusher = pusher.Grouping("hostname", config.Pushgateway.Instance)
				// 添加分组标签
				pusher = pusher.Grouping("group_id", fmt.Sprintf("%%d", HostGroupID))
				pusher = pusher.Grouping("group_name", HostGroupName)

				log.Printf("[分组标签] 推送指标 - group_id: %%d, group_name: %%s", HostGroupID, HostGroupName)

				// 使用registry推送所有指标（不包含Go运行时指标）
				if err := pusher.Gatherer(registry).Push(); err != nil {
					log.Printf("Could not push to gateway: %%v", err)
				} else {
					log.Printf("Successfully pushed metrics to gateway: %%s (group: %%s)", config.Pushgateway.URL, HostGroupName)
				}
			}
		case <-heartbeatTicker.C:
			if config.Heartbeat.ServerURL != "" {
				sendHeartbeat(config.Heartbeat.ServerURL, config.Heartbeat.Token)
			}
		}
	}
}

func main() {
	log.Println("Starting dodevops-agent...")
	log.Printf("[分组信息] HostGroupID: %%d, HostGroupName: %%s", HostGroupID, HostGroupName)

	// 获取本地IP作为Job名称
	localIP := getLocalIP()
	if localIP == "" {
		localIP = "unknown-ip"
		log.Printf("Warning: could not get local IP")
	}
	log.Printf("Agent local IP: %%s", localIP)

	// 获取hostname作为Instance标识
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
		log.Printf("Warning: could not get hostname: %%v", err)
	}
	log.Printf("Agent hostname: %%s", hostname)

	config := &Config{
		Processes: []struct {
			Name string
		}{
			{Name: "nginx"},
			{Name: "mysql"},
			{Name: "redis"},
		},
		Endpoints: []struct {
			URL    string
			Method string
		}{
			{URL: "http://localhost:12580/metrics", Method: "GET"},
		},
		Pushgateway: struct {
			URL          string
			PushInterval string
			Job          string
			Instance     string
		}{
			URL:          getEnvVar("PUSHGATEWAY_URL", "%[1]s"),
			PushInterval: "30s",
			Job:          localIP,      // Job使用本地IP
			Instance:     hostname,     // Instance使用hostname
		},
		Heartbeat: struct {
			ServerURL string
			Interval  string
			Token     string
		}{
			ServerURL: getEnvVar("HEARTBEAT_SERVER_URL", "%[2]s"),
			Interval:  "30s",
			Token:     getEnvVar("HEARTBEAT_TOKEN", "%[3]s"),
		},
		CollectInterval: "15s",
	}

	if err := StartAgentWithConfig(config); err != nil {
		log.Fatalf("Failed to start agent: %%v", err)
	}
	select {}
}
`, pushgatewayURL, heartbeatURL, heartbeatToken, groupID, groupName)

	log.Printf("CreateAgentMainFile: 生成的代码长度: %d", len(mainCode))
	log.Printf("CreateAgentMainFile: 代码生成完成，heartbeat URL=%s, groupID=%d, groupName=%s", heartbeatURL, groupID, groupName)

	return mainCode
}

// ValidateAgentCode 验证生成的agent代码是否有效
func ValidateAgentCode(code string) error {
	// 检查基本结构
	if !strings.Contains(code, "package main") {
		return fmt.Errorf("缺少package main声明")
	}

	if !strings.Contains(code, "func main()") {
		return fmt.Errorf("缺少main函数")
	}

	// 简化验证：只检查明显的语法错误，不做复杂的解析
	lines := strings.Split(code, "\n")

	// 检查是否有明显的语法问题
	openBraces := 0
	openParens := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		openBraces += strings.Count(trimmed, "{")
		openBraces -= strings.Count(trimmed, "}")
		openParens += strings.Count(trimmed, "(")
		openParens -= strings.Count(trimmed, ")")
	}

	// 检查括号匹配
	if openBraces != 0 {
		return fmt.Errorf("大括号不匹配")
	}
	if openParens != 0 {
		return fmt.Errorf("圆括号不匹配")
	}

	return nil
}
