package snmp

import (
	"dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	"dodevops-api/api/monitor/collector"
	"dodevops-api/common/config"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

// Manager SNMP采集管理器(参考agent模式,自动化管理)
type Manager struct {
	deviceDao      *dao.CmdbSnmpDeviceDao
	pushgatewayURL string
	devices        map[uint]*DeviceCollector // 设备ID -> 采集器
	mu             sync.RWMutex
	stopChan       chan struct{}
}

var (
	instance *Manager
	once     sync.Once
)

// GetManager 获取SNMP管理器单例
func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{
			deviceDao:      dao.NewCmdbSnmpDeviceDao(),
			pushgatewayURL: config.Config.Monitor.Pushgateway.URL,
			devices:        make(map[uint]*DeviceCollector),
			stopChan:       make(chan struct{}),
		}
	})
	return instance
}

// Start 启动SNMP采集管理器
func (m *Manager) Start() {
	log.Println("[SNMP Manager] 启动中...")

	// 1. 从数据库加载所有启用监控的设备
	m.loadDevices()

	// 2. 启动定时重载设备列表(每1分钟检查一次)
	go m.reloadLoop()

	log.Printf("[SNMP Manager] 已启动,当前监控 %d 个设备\n", len(m.devices))
}

// Reload 立即重载设备列表(API调用)
func (m *Manager) Reload() {
	log.Println("[SNMP Manager] 立即重载设备列表...")
	m.loadDevices()
}

// Stop 停止SNMP采集管理器
func (m *Manager) Stop() {
	log.Println("[SNMP Manager] 停止中...")
	close(m.stopChan)

	m.mu.Lock()
	for _, dc := range m.devices {
		dc.Stop()
	}
	m.devices = make(map[uint]*DeviceCollector)
	m.mu.Unlock()

	log.Println("[SNMP Manager] 已停止")
}

// loadDevices 从数据库加载设备并启动采集
func (m *Manager) loadDevices() {
	devices, err := m.deviceDao.GetMonitorEnabledDevices()
	if err != nil {
		log.Printf("[SNMP Manager] 加载设备失败: %v\n", err)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 停止已删除的设备
	for id := range m.devices {
		exists := false
		for _, d := range devices {
			if d.ID == id {
				exists = true
				break
			}
		}
		if !exists {
			deviceName := m.devices[id].device.DeviceName
			m.devices[id].Stop()
			delete(m.devices, id)
			log.Printf("[SNMP Manager] 停止设备: %s (ID:%d)\n", deviceName, id)
		}
	}

	// 启动新设备或更新现有设备
	for _, device := range devices {
		if dc, exists := m.devices[device.ID]; exists {
			// 已存在,检查是否需要更新
			if dc.device.IPAddress != device.IPAddress || dc.device.SnmpCommunity != device.SnmpCommunity {
				// 配置变化,重启采集器
				dc.Stop()
				delete(m.devices, device.ID)
			} else {
				continue // 无变化,跳过
			}
		}

		// 创建并启动新的采集器
		dc := NewDeviceCollector(&device, m.pushgatewayURL)
		dc.Start()
		m.devices[device.ID] = dc
		log.Printf("[SNMP Manager] 启动设备: %s (ID:%d, %s)\n", device.DeviceName, device.ID, device.IPAddress)
	}
}

// reloadLoop 定时重载设备列表
func (m *Manager) reloadLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.loadDevices()
		case <-m.stopChan:
			return
		}
	}
}

// CollectNow 立即采集指定设备(API调用)
func (m *Manager) CollectNow(deviceID uint) error {
	m.mu.RLock()
	dc, exists := m.devices[deviceID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("设备未启用监控或不存在")
	}

	return dc.CollectAndPush()
}

// ============= DeviceCollector 单个设备采集器 =============

type DeviceCollector struct {
	device         *cmdbModel.CmdbSnmpDevice
	pushgatewayURL string
	collector      *collector.SNMPCollector
	deviceDao      *dao.CmdbSnmpDeviceDao // 添加: 用于更新数据库状态
	vendor         string                 // 设备厂商(H3C/Cisco/Huawei等)

	// Prometheus指标 - 基础信息
	sysDescr      *prometheus.GaugeVec // 系统描述(作为label)
	cpuGauge      prometheus.Gauge
	cpuCoreGauge  *prometheus.GaugeVec // CPU使用率(每核)
	memoryGauge   prometheus.Gauge
	uptimeGauge   prometheus.Gauge
	statusGauge   prometheus.Gauge
	tempGauge     prometheus.Gauge     // 设备温度(交换机)

	// 带宽利用率指标
	portBandwidthInUtil  *prometheus.GaugeVec // 端口入向带宽利用率(%)
	portBandwidthOutUtil *prometheus.GaugeVec // 端口出向带宽利用率(%)

	// 磁盘指标
	diskUsage     *prometheus.GaugeVec // 磁盘使用率
	diskReadBytes *prometheus.CounterVec // 磁盘读字节
	diskWriteBytes *prometheus.CounterVec // 磁盘写字节

	// 端口指标
	portInBytes   *prometheus.CounterVec
	portOutBytes  *prometheus.CounterVec
	portStatus    *prometheus.GaugeVec
	portSpeed     *prometheus.GaugeVec // 端口速率

	// 端口流量速率(实时,易读)
	portInBytesRate  *prometheus.GaugeVec // 入向流量速率(Mbps)
	portOutBytesRate *prometheus.GaugeVec // 出向流量速率(Mbps)

	// 端口错误和丢包指标
	portInErrorsRate   *prometheus.GaugeVec // 入向错误包速率(errors/sec)
	portOutErrorsRate  *prometheus.GaugeVec // 出向错误包速率(errors/sec)
	portInDiscardsRate *prometheus.GaugeVec // 入向丢包速率(discards/sec)
	portOutDiscardsRate *prometheus.GaugeVec // 出向丢包速率(discards/sec)

	// 带宽计算辅助数据(用于计算增量)
	lastPortInBytes    map[int]uint64 // 端口索引 -> 上次入向字节数
	lastPortOutBytes   map[int]uint64 // 端口索引 -> 上次出向字节数
	lastPortInErrors   map[int]uint64 // 端口索引 -> 上次入向错误包数
	lastPortOutErrors  map[int]uint64 // 端口索引 -> 上次出向错误包数
	lastPortInDiscards map[int]uint64 // 端口索引 -> 上次入向丢包数
	lastPortOutDiscards map[int]uint64 // 端口索引 -> 上次出向丢包数
	lastCollectTime    time.Time      // 上次采集时间

	stopChan chan struct{}
}

func NewDeviceCollector(device *cmdbModel.CmdbSnmpDevice, pushgatewayURL string) *DeviceCollector {
	// 解析IP和端口
	ip, port := parseIPPort(device.IPAddress)

	dc := &DeviceCollector{
		device:              device,
		pushgatewayURL:      pushgatewayURL,
		deviceDao:           dao.NewCmdbSnmpDeviceDao(), // 初始化DAO
		lastPortInBytes:     make(map[int]uint64),
		lastPortOutBytes:    make(map[int]uint64),
		lastPortInErrors:    make(map[int]uint64),
		lastPortOutErrors:   make(map[int]uint64),
		lastPortInDiscards:  make(map[int]uint64),
		lastPortOutDiscards: make(map[int]uint64),
		lastCollectTime:     time.Now(),
		stopChan:            make(chan struct{}),
	}

	// 创建SNMP采集器
	dc.collector = collector.NewSNMPCollector(
		ip,
		uint16(port),
		device.SnmpCommunity,
		device.SnmpVersion,
		5, // 超时5秒
		3, // 重试3次
	)

	// 初始化Prometheus指标
	dc.initMetrics()

	return dc
}

// initMetrics 初始化Prometheus指标
func (dc *DeviceCollector) initMetrics() {
	labels := prometheus.Labels{
		"device_id":   fmt.Sprintf("%d", dc.device.ID),
		"device_name": dc.device.DeviceName,
		"device_type": dc.device.DeviceType,
		"ip":          dc.device.IPAddress,
	}

	// 系统描述(包含设备型号/OS信息)
	dc.sysDescr = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "snmp_device_info",
			Help:        "设备信息(sys_descr作为label)",
			ConstLabels: labels,
		},
		[]string{"sys_descr", "sys_name"},
	)

	// CPU整体使用率
	dc.cpuGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "snmp_device_cpu_percent",
		Help:        "CPU平均使用率(%)",
		ConstLabels: labels,
	})

	// CPU每核使用率
	dc.cpuCoreGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "snmp_device_cpu_core_percent",
			Help:        "CPU每核使用率(%)",
			ConstLabels: labels,
		},
		[]string{"core_index"},
	)

	// 内存使用率
	dc.memoryGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "snmp_device_memory_percent",
		Help:        "内存使用率(%)",
		ConstLabels: labels,
	})

	// 运行时间
	dc.uptimeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "snmp_device_uptime_days",
		Help:        "运行时间(天)",
		ConstLabels: labels,
	})

	// 设备状态
	dc.statusGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "snmp_device_status",
		Help:        "设备状态(1=在线,0=离线)",
		ConstLabels: labels,
	})

	// 设备温度(交换机)
	dc.tempGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "snmp_device_temperature_celsius",
		Help:        "设备温度(摄氏度)",
		ConstLabels: labels,
	})

	// 磁盘使用率
	dc.diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "snmp_device_disk_usage_percent",
			Help:        "磁盘使用率(%)",
			ConstLabels: labels,
		},
		[]string{"disk_descr", "disk_index"},
	)

	// 磁盘读字节数
	dc.diskReadBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "snmp_device_disk_read_bytes",
			Help:        "磁盘读取字节数",
			ConstLabels: labels,
		},
		[]string{"disk_descr", "disk_index"},
	)

	// 磁盘写字节数
	dc.diskWriteBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "snmp_device_disk_write_bytes",
			Help:        "磁盘写入字节数",
			ConstLabels: labels,
		},
		[]string{"disk_descr", "disk_index"},
	)

	// 端口入站字节
	dc.portInBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "snmp_port_in_bytes",
			Help: "端口入站字节数",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)

	// 端口出站字节
	dc.portOutBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "snmp_port_out_bytes",
			Help: "端口出站字节数",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)

	// 端口状态
	dc.portStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "snmp_port_status",
			Help: "端口状态(1=up,0=down)",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)

	// 端口速率
	dc.portSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "snmp_port_speed_mbps",
			Help: "端口速率(Mbps)",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)

	// 端口流量速率(实时)
	dc.portInBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "snmp_port_in_rate_mbps",
			Help: "端口入向流量速率(Mbps)",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)

	dc.portOutBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "snmp_port_out_rate_mbps",
			Help: "端口出向流量速率(Mbps)",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)

	// 带宽利用率
	dc.portBandwidthInUtil = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "snmp_port_bandwidth_in_utilization",
			Help: "端口入向带宽利用率(%)",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)

	dc.portBandwidthOutUtil = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "snmp_port_bandwidth_out_utilization",
			Help: "端口出向带宽利用率(%)",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)

	// 端口错误包速率
	dc.portInErrorsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "snmp_port_in_errors_per_second",
			Help: "端口入向错误包速率(errors/sec)",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)

	dc.portOutErrorsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "snmp_port_out_errors_per_second",
			Help: "端口出向错误包速率(errors/sec)",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)

	// 端口丢包速率
	dc.portInDiscardsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "snmp_port_in_discards_per_second",
			Help: "端口入向丢包速率(discards/sec)",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)

	dc.portOutDiscardsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "snmp_port_out_discards_per_second",
			Help: "端口出向丢包速率(discards/sec)",
		},
		[]string{"device_id", "device_name", "port_index", "port_name"},
	)
}

// Start 启动定时采集
func (dc *DeviceCollector) Start() {
	// 立即采集一次
	go dc.CollectAndPush()

	// 启动定时采集(默认60秒)
	go dc.collectLoop()
}

// Stop 停止采集
func (dc *DeviceCollector) Stop() {
	close(dc.stopChan)
}

// collectLoop 定时采集循环
func (dc *DeviceCollector) collectLoop() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = dc.CollectAndPush()
		case <-dc.stopChan:
			return
		}
	}
}

// CollectAndPush 采集并推送到Pushgateway
func (dc *DeviceCollector) CollectAndPush() error {
	// 1. 测试连接
	err := dc.collector.TestConnection()
	if err != nil {
		dc.statusGauge.Set(0) // 离线

		// 同步数据库状态为离线
		_ = dc.deviceDao.UpdateDeviceStatus(dc.device.ID, "offline")

		_ = dc.pushMetrics()
		log.Printf("[%s] 离线: %v\n", dc.device.DeviceName, err)
		return err
	}

	// 2. 设备在线
	dc.statusGauge.Set(1)

	// 同步数据库状态为在线,并更新最后在线时间
	_ = dc.deviceDao.UpdateDeviceStatus(dc.device.ID, "online")
	_ = dc.deviceDao.UpdateLastSeen(dc.device.ID)

	// 3. 识别设备厂商(首次采集时，或之前识别失败时)
	if dc.vendor == "" || dc.vendor == "Unknown" {
		dc.detectVendor()
	}

	// 4. 根据设备类型采集不同指标
	switch dc.device.DeviceType {
	case "switch", "router":
		// 交换机/路由器: 采集厂商特定指标 + 端口指标
		dc.collectSwitchMetrics()
	case "windows", "linux":
		// 服务器: 采集 Host Resources MIB 指标
		dc.collectServerMetrics()
	default:
		// 未知类型: 只采集基础信息
		dc.collectBasicInfo()
	}

	// 5. 推送到Pushgateway
	err = dc.pushMetrics()
	if err != nil {
		log.Printf("[%s] 推送失败: %v\n", dc.device.DeviceName, err)
		return err
	}

	log.Printf("[%s] 采集成功 (厂商:%s)\n", dc.device.DeviceName, dc.vendor)
	return nil
}

// collectBasicMetrics 采集基础指标(增强版)
func (dc *DeviceCollector) collectBasicMetrics() {
	// 1. 采集系统基础信息
	sysOids := []string{
		collector.OID_SYS_DESCR,   // 系统描述
		collector.OID_SYS_NAME,    // 系统名称
		collector.OID_SYS_UPTIME,  // 运行时间
	}

	sysData, err := dc.collector.Get(sysOids)
	if err != nil {
		return
	}

	// 系统描述(包含设备型号/OS信息)
	sysDescr := collector.ToString(collector.GetOIDValue(sysData, collector.OID_SYS_DESCR))
	sysName := collector.ToString(collector.GetOIDValue(sysData, collector.OID_SYS_NAME))
	if sysDescr != "" {
		dc.sysDescr.WithLabelValues(sysDescr, sysName).Set(1)
	}

	// 运行时间(已由parseValue转换为秒,直接除以86400转换为天)
	uptime := collector.ToInt64(collector.GetOIDValue(sysData, collector.OID_SYS_UPTIME))
	dc.uptimeGauge.Set(float64(uptime) / 86400.0)

	// 2. CPU使用率采集
	cpuData, _ := dc.collector.Walk(collector.OID_HR_PROCESSOR_LOAD)
	if len(cpuData) > 0 {
		totalCPU := 0.0
		for i, pdu := range cpuData {
			cpuUsage := collector.ToFloat64(pdu.Value)
			// 记录每核CPU使用率
			dc.cpuCoreGauge.WithLabelValues(fmt.Sprintf("%d", i)).Set(cpuUsage)
			totalCPU += cpuUsage
		}
		// 记录平均CPU使用率
		avgCPU := totalCPU / float64(len(cpuData))
		dc.cpuGauge.Set(avgCPU)
	}

	// 3. 内存使用率采集
	dc.collectMemoryMetrics()

	// 4. 磁盘使用率采集
	dc.collectDiskMetrics()
}

// collectMemoryMetrics 采集内存指标
func (dc *DeviceCollector) collectMemoryMetrics() {
	// 获取存储信息
	storageTypeData, _ := dc.collector.BulkWalk(collector.OID_HR_STORAGE_TYPE)
	storageSizeData, _ := dc.collector.BulkWalk(collector.OID_HR_STORAGE_SIZE)
	storageUsedData, _ := dc.collector.BulkWalk(collector.OID_HR_STORAGE_USED)

	if len(storageTypeData) == 0 {
		return
	}

	// 查找物理内存(RAM)
	var totalMemory, usedMemory float64
	for i, typePdu := range storageTypeData {
		// hrStorageRam OID: 1.3.6.1.2.1.25.2.1.2
		typeOID := collector.ToString(typePdu.Value)
		if strings.Contains(typeOID, "1.3.6.1.2.1.25.2.1.2") {
			// 找到物理内存
			if i < len(storageSizeData) && i < len(storageUsedData) {
				size := collector.ToFloat64(storageSizeData[i].Value)
				used := collector.ToFloat64(storageUsedData[i].Value)
				totalMemory += size
				usedMemory += used
			}
		}
	}

	if totalMemory > 0 {
		memoryPercent := (usedMemory / totalMemory) * 100.0
		dc.memoryGauge.Set(memoryPercent)
	}
}

// collectDiskMetrics 采集磁盘指标
func (dc *DeviceCollector) collectDiskMetrics() {
	// 获取存储设备信息
	storageDescrData, _ := dc.collector.BulkWalk(collector.OID_HR_STORAGE_DESCR)
	storageTypeData, _ := dc.collector.BulkWalk(collector.OID_HR_STORAGE_TYPE)
	storageSizeData, _ := dc.collector.BulkWalk(collector.OID_HR_STORAGE_SIZE)
	storageUsedData, _ := dc.collector.BulkWalk(collector.OID_HR_STORAGE_USED)

	if len(storageDescrData) == 0 {
		return
	}

	for i, descrPdu := range storageDescrData {
		if i >= len(storageTypeData) || i >= len(storageSizeData) || i >= len(storageUsedData) {
			continue
		}

		// 获取存储类型
		typeOID := collector.ToString(storageTypeData[i].Value)
		// 只采集固定磁盘: hrStorageFixedDisk OID: 1.3.6.1.2.1.25.2.1.4
		if !strings.Contains(typeOID, "1.3.6.1.2.1.25.2.1.4") {
			continue
		}

		diskDescr := collector.ToString(descrPdu.Value)
		diskIndex := extractPortIndex(descrPdu.Name)

		size := collector.ToFloat64(storageSizeData[i].Value)
		used := collector.ToFloat64(storageUsedData[i].Value)

		if size > 0 {
			usagePercent := (used / size) * 100.0
			dc.diskUsage.WithLabelValues(diskDescr, fmt.Sprintf("%d", diskIndex)).Set(usagePercent)
		}
	}

	// 采集磁盘IO(如果设备支持)
	dc.collectDiskIO()
}

// collectDiskIO 采集磁盘IO指标
func (dc *DeviceCollector) collectDiskIO() {
	// 使用 hrDiskStorageEntry 的读写计数器(部分设备支持)
	// OID: 1.3.6.1.2.1.25.3.6.1.x (HOST-RESOURCES-MIB)
	// 注意: 不是所有SNMP设备都支持磁盘IO统计

	diskReadOID := "1.3.6.1.2.1.25.3.6.1.12" // hrDiskStorageReads
	diskWriteOID := "1.3.6.1.2.1.25.3.6.1.13" // hrDiskStorageWrites

	readData, err1 := dc.collector.BulkWalk(diskReadOID)
	writeData, err2 := dc.collector.BulkWalk(diskWriteOID)

	if err1 != nil || err2 != nil {
		// 设备不支持磁盘IO统计,跳过
		return
	}

	for i, readPdu := range readData {
		if i >= len(writeData) {
			break
		}

		diskIndex := extractPortIndex(readPdu.Name)
		diskDescr := fmt.Sprintf("Disk%d", diskIndex)

		readBytes := collector.ToFloat64(readPdu.Value)
		writeBytes := collector.ToFloat64(writeData[i].Value)

		dc.diskReadBytes.WithLabelValues(diskDescr, fmt.Sprintf("%d", diskIndex)).Add(readBytes)
		dc.diskWriteBytes.WithLabelValues(diskDescr, fmt.Sprintf("%d", diskIndex)).Add(writeBytes)
	}
}

// collectPortMetrics 采集端口指标(增强版,包含带宽利用率、错误包和丢包)
func (dc *DeviceCollector) collectPortMetrics() {
	// 获取端口信息
	inOctetsPDUs, err := dc.collector.BulkWalk(collector.OID_IF_IN_OCTETS)
	if err != nil {
		return
	}

	outOctetsPDUs, _ := dc.collector.BulkWalk(collector.OID_IF_OUT_OCTETS)
	operStatusPDUs, _ := dc.collector.BulkWalk(collector.OID_IF_OPER_STA)
	ifDescrPDUs, _ := dc.collector.BulkWalk(collector.OID_IF_DESCR)
	ifSpeedPDUs, _ := dc.collector.BulkWalk(collector.OID_IF_SPEED)         // 端口速率(bps)
	ifHighSpeedPDUs, _ := dc.collector.BulkWalk(collector.OID_IF_HIGH_SPEED) // 高速端口速率(Mbps)

	// 获取错误包和丢包数据
	inErrorsPDUs, _ := dc.collector.BulkWalk(collector.OID_IF_IN_ERRORS)   // 入向错误包
	outErrorsPDUs, _ := dc.collector.BulkWalk(collector.OID_IF_OUT_ERRORS) // 出向错误包
	inDiscardsPDUs, _ := dc.collector.BulkWalk(collector.OID_IF_IN_DISCARD) // 入向丢包
	outDiscardsPDUs, _ := dc.collector.BulkWalk(collector.OID_IF_OUT_DISCA) // 出向丢包

	// 构建端口名称映射
	portNames := make(map[int]string)
	for _, pdu := range ifDescrPDUs {
		portIndex := extractPortIndex(pdu.Name)
		portNames[portIndex] = collector.ToString(pdu.Value)
	}

	// 构建端口速率映射(优先使用ifHighSpeed)
	portSpeeds := make(map[int]uint64)

	// 先从ifHighSpeed获取(单位:Mbps,需要转换为bps)
	highSpeedMap := make(map[int]uint64)
	for _, pdu := range ifHighSpeedPDUs {
		portIndex := extractPortIndex(pdu.Name)
		highSpeed := uint64(collector.ToFloat64(pdu.Value))
		if highSpeed > 0 {
			highSpeedMap[portIndex] = highSpeed * 1000000 // Mbps → bps
		}
	}

	// 再从ifSpeed获取(单位:bps)
	for _, pdu := range ifSpeedPDUs {
		portIndex := extractPortIndex(pdu.Name)
		speed := uint64(collector.ToFloat64(pdu.Value))

		// 优先级: ifHighSpeed > ifSpeed (如果ifSpeed不是溢出值)
		if highSpeed, hasHighSpeed := highSpeedMap[portIndex]; hasHighSpeed {
			// 优先使用ifHighSpeed
			portSpeeds[portIndex] = highSpeed
		} else if speed > 0 && speed != 4294967295 {
			// ifSpeed有效且未溢出,使用ifSpeed
			portSpeeds[portIndex] = speed
		}
		// 如果ifSpeed溢出(4294967295)且没有ifHighSpeed,则不记录速率
	}

	// 构建错误包映射
	portInErrors := make(map[int]uint64)
	for _, pdu := range inErrorsPDUs {
		portIndex := extractPortIndex(pdu.Name)
		portInErrors[portIndex] = uint64(collector.ToFloat64(pdu.Value))
	}

	portOutErrors := make(map[int]uint64)
	for _, pdu := range outErrorsPDUs {
		portIndex := extractPortIndex(pdu.Name)
		portOutErrors[portIndex] = uint64(collector.ToFloat64(pdu.Value))
	}

	// 构建丢包映射
	portInDiscards := make(map[int]uint64)
	for _, pdu := range inDiscardsPDUs {
		portIndex := extractPortIndex(pdu.Name)
		portInDiscards[portIndex] = uint64(collector.ToFloat64(pdu.Value))
	}

	portOutDiscards := make(map[int]uint64)
	for _, pdu := range outDiscardsPDUs {
		portIndex := extractPortIndex(pdu.Name)
		portOutDiscards[portIndex] = uint64(collector.ToFloat64(pdu.Value))
	}

	// 计算时间间隔
	currentTime := time.Now()
	intervalSeconds := currentTime.Sub(dc.lastCollectTime).Seconds()

	// 如果是第一次采集,或者间隔太短/太长,跳过带宽利用率计算
	skipBandwidthCalc := intervalSeconds < 10 || intervalSeconds > 300

	// 处理入站流量和带宽利用率
	for _, pdu := range inOctetsPDUs {
		portIndex := extractPortIndex(pdu.Name)
		portName := portNames[portIndex]
		if portName == "" {
			portName = fmt.Sprintf("Port%d", portIndex)
		}

		labels := prometheus.Labels{
			"device_id":   fmt.Sprintf("%d", dc.device.ID),
			"device_name": dc.device.DeviceName,
			"port_index":  fmt.Sprintf("%d", portIndex),
			"port_name":   portName,
		}

		currentInBytes := uint64(collector.ToFloat64(pdu.Value))

		// 计算入向流量和速率
		if !skipBandwidthCalc {
			if lastInBytes, exists := dc.lastPortInBytes[portIndex]; exists {
				// 计算字节增量(处理计数器溢出)
				var deltaBytes uint64
				if currentInBytes >= lastInBytes {
					deltaBytes = currentInBytes - lastInBytes
				} else {
					// 32位计数器溢出: 4294967295 → 0
					deltaBytes = (4294967295 - lastInBytes) + currentInBytes
				}

				// 添加增量到Counter(不是累计值!)
				dc.portInBytes.With(labels).Add(float64(deltaBytes))

				// 计算实时流量速率(Mbps)
				rateMbps := (float64(deltaBytes) * 8) / (intervalSeconds * 1000000.0)
				dc.portInBytesRate.With(labels).Set(rateMbps)

				// 计算带宽利用率(%)
				if speed, ok := portSpeeds[portIndex]; ok && speed > 0 && speed != 4294967295 {
					// 带宽利用率公式: (ΔOctets × 8) ÷ (interval_seconds × ifSpeed) × 100
					utilization := (float64(deltaBytes) * 8) / (intervalSeconds * float64(speed)) * 100.0

					// 限制范围: 0-100%
					if utilization < 0 {
						utilization = 0
					} else if utilization > 100 {
						utilization = 100
					}

					dc.portBandwidthInUtil.With(labels).Set(utilization)
				}
			}
		}

		// 保存当前值供下次计算
		dc.lastPortInBytes[portIndex] = currentInBytes

		// 计算入向错误包速率
		if currentInErrors, ok := portInErrors[portIndex]; ok {
			if !skipBandwidthCalc {
				if lastInErrors, exists := dc.lastPortInErrors[portIndex]; exists {
					// 计算错误包增量(处理计数器溢出)
					var deltaErrors uint64
					if currentInErrors >= lastInErrors {
						deltaErrors = currentInErrors - lastInErrors
					} else {
						// 32位计数器溢出
						deltaErrors = (4294967295 - lastInErrors) + currentInErrors
					}

					// 计算每秒错误包速率
					errorRate := float64(deltaErrors) / intervalSeconds
					dc.portInErrorsRate.With(labels).Set(errorRate)
				}
			}
			dc.lastPortInErrors[portIndex] = currentInErrors
		}

		// 计算入向丢包速率
		if currentInDiscards, ok := portInDiscards[portIndex]; ok {
			if !skipBandwidthCalc {
				if lastInDiscards, exists := dc.lastPortInDiscards[portIndex]; exists {
					// 计算丢包增量(处理计数器溢出)
					var deltaDiscards uint64
					if currentInDiscards >= lastInDiscards {
						deltaDiscards = currentInDiscards - lastInDiscards
					} else {
						// 32位计数器溢出
						deltaDiscards = (4294967295 - lastInDiscards) + currentInDiscards
					}

					// 计算每秒丢包速率
					discardRate := float64(deltaDiscards) / intervalSeconds
					dc.portInDiscardsRate.With(labels).Set(discardRate)
				}
			}
			dc.lastPortInDiscards[portIndex] = currentInDiscards
		}
	}

	// 处理出站流量和带宽利用率
	for _, pdu := range outOctetsPDUs {
		portIndex := extractPortIndex(pdu.Name)
		portName := portNames[portIndex]
		if portName == "" {
			portName = fmt.Sprintf("Port%d", portIndex)
		}

		labels := prometheus.Labels{
			"device_id":   fmt.Sprintf("%d", dc.device.ID),
			"device_name": dc.device.DeviceName,
			"port_index":  fmt.Sprintf("%d", portIndex),
			"port_name":   portName,
		}

		currentOutBytes := uint64(collector.ToFloat64(pdu.Value))

		// 计算出向流量和速率
		if !skipBandwidthCalc {
			if lastOutBytes, exists := dc.lastPortOutBytes[portIndex]; exists {
				// 计算字节增量(处理计数器溢出)
				var deltaBytes uint64
				if currentOutBytes >= lastOutBytes {
					deltaBytes = currentOutBytes - lastOutBytes
				} else {
					// 32位计数器溢出: 4294967295 → 0
					deltaBytes = (4294967295 - lastOutBytes) + currentOutBytes
				}

				// 添加增量到Counter(不是累计值!)
				dc.portOutBytes.With(labels).Add(float64(deltaBytes))

				// 计算实时流量速率(Mbps)
				rateMbps := (float64(deltaBytes) * 8) / (intervalSeconds * 1000000.0)
				dc.portOutBytesRate.With(labels).Set(rateMbps)

				// 计算带宽利用率(%)
				if speed, ok := portSpeeds[portIndex]; ok && speed > 0 && speed != 4294967295 {
					// 带宽利用率公式: (ΔOctets × 8) ÷ (interval_seconds × ifSpeed) × 100
					utilization := (float64(deltaBytes) * 8) / (intervalSeconds * float64(speed)) * 100.0

					// 限制范围: 0-100%
					if utilization < 0 {
						utilization = 0
					} else if utilization > 100 {
						utilization = 100
					}

					dc.portBandwidthOutUtil.With(labels).Set(utilization)
				}
			}
		}

		// 保存当前值供下次计算
		dc.lastPortOutBytes[portIndex] = currentOutBytes

		// 计算出向错误包速率
		if currentOutErrors, ok := portOutErrors[portIndex]; ok {
			if !skipBandwidthCalc {
				if lastOutErrors, exists := dc.lastPortOutErrors[portIndex]; exists {
					// 计算错误包增量(处理计数器溢出)
					var deltaErrors uint64
					if currentOutErrors >= lastOutErrors {
						deltaErrors = currentOutErrors - lastOutErrors
					} else {
						// 32位计数器溢出
						deltaErrors = (4294967295 - lastOutErrors) + currentOutErrors
					}

					// 计算每秒错误包速率
					errorRate := float64(deltaErrors) / intervalSeconds
					dc.portOutErrorsRate.With(labels).Set(errorRate)
				}
			}
			dc.lastPortOutErrors[portIndex] = currentOutErrors
		}

		// 计算出向丢包速率
		if currentOutDiscards, ok := portOutDiscards[portIndex]; ok {
			if !skipBandwidthCalc {
				if lastOutDiscards, exists := dc.lastPortOutDiscards[portIndex]; exists {
					// 计算丢包增量(处理计数器溢出)
					var deltaDiscards uint64
					if currentOutDiscards >= lastOutDiscards {
						deltaDiscards = currentOutDiscards - lastOutDiscards
					} else {
						// 32位计数器溢出
						deltaDiscards = (4294967295 - lastOutDiscards) + currentOutDiscards
					}

					// 计算每秒丢包速率
					discardRate := float64(deltaDiscards) / intervalSeconds
					dc.portOutDiscardsRate.With(labels).Set(discardRate)
				}
			}
			dc.lastPortOutDiscards[portIndex] = currentOutDiscards
		}
	}

	// 处理端口状态和速率
	for _, pdu := range operStatusPDUs {
		portIndex := extractPortIndex(pdu.Name)
		portName := portNames[portIndex]
		if portName == "" {
			portName = fmt.Sprintf("Port%d", portIndex)
		}

		labels := prometheus.Labels{
			"device_id":   fmt.Sprintf("%d", dc.device.ID),
			"device_name": dc.device.DeviceName,
			"port_index":  fmt.Sprintf("%d", portIndex),
			"port_name":   portName,
		}

		// 端口状态
		status := collector.ToInt64(pdu.Value)
		if status == 1 {
			dc.portStatus.With(labels).Set(1)
		} else {
			dc.portStatus.With(labels).Set(0)
		}

		// 端口速率(Mbps) - 将bps转换为Mbps,未使用的端口设为0
		if speed, ok := portSpeeds[portIndex]; ok && speed > 0 && speed != 4294967295 {
			// 将bps转换为Mbps
			speedMbps := float64(speed) / 1000000.0
			dc.portSpeed.With(labels).Set(speedMbps)
		} else {
			// 未使用或无效速率的端口设为0
			dc.portSpeed.With(labels).Set(0)
		}
	}

	// 更新最后采集时间
	dc.lastCollectTime = currentTime
}

// pushMetrics 推送指标到Pushgateway
func (dc *DeviceCollector) pushMetrics() error {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		// 基础指标
		dc.sysDescr,
		dc.cpuGauge,
		dc.cpuCoreGauge,
		dc.memoryGauge,
		dc.uptimeGauge,
		dc.statusGauge,
		dc.tempGauge, // 温度指标
		// 磁盘指标
		dc.diskUsage,
		dc.diskReadBytes,
		dc.diskWriteBytes,
		// 端口指标
		dc.portInBytes,
		dc.portOutBytes,
		dc.portStatus,
		dc.portSpeed,
		dc.portInBytesRate,  // 端口流量速率(Mbps)
		dc.portOutBytesRate, // 端口流量速率(Mbps)
		dc.portBandwidthInUtil,
		dc.portBandwidthOutUtil,
		// 端口错误和丢包指标
		dc.portInErrorsRate,
		dc.portOutErrorsRate,
		dc.portInDiscardsRate,
		dc.portOutDiscardsRate,
	)

	// 使用设备名称作为job名称
	jobName := fmt.Sprintf("snmp_%s", dc.device.DeviceName)

	pusher := push.New(dc.pushgatewayURL, jobName).Gatherer(registry)

	return pusher.Push()
}

// ============= 辅助函数 =============

// parseIPPort 解析IP:端口
func parseIPPort(ipPort string) (string, int) {
	parts := strings.Split(ipPort, ":")
	if len(parts) == 2 {
		var port int
		fmt.Sscanf(parts[1], "%d", &port)
		return parts[0], port
	}
	return ipPort, 161 // 默认161端口
}

// extractPortIndex 从OID提取端口索引
func extractPortIndex(oid string) int {
	parts := strings.Split(oid, ".")
	if len(parts) > 0 {
		var index int
		fmt.Sscanf(parts[len(parts)-1], "%d", &index)
		return index
	}
	return 0
}

// ============= 新增: 厂商识别和分类采集方法 =============

// detectVendor 识别设备厂商(从sysDescr)
func (dc *DeviceCollector) detectVendor() {
	sysData, err := dc.collector.Get([]string{collector.OID_SYS_DESCR})
	if err != nil {
		log.Printf("[%s] 厂商识别失败: 无法获取sysDescr - %v\n", dc.device.DeviceName, err)
		dc.vendor = "Unknown"
		return
	}

	rawValue := collector.GetOIDValue(sysData, collector.OID_SYS_DESCR)
	sysDescr := strings.ToLower(collector.ToString(rawValue))

	// 识别厂商(更宽松的匹配规则)
	if strings.Contains(sysDescr, "h3c") || strings.Contains(sysDescr, "3com") ||
	   strings.Contains(sysDescr, "comware") || strings.Contains(sysDescr, "s5130") {
		dc.vendor = "H3C"
	} else if strings.Contains(sysDescr, "cisco") {
		dc.vendor = "Cisco"
	} else if strings.Contains(sysDescr, "huawei") || strings.Contains(sysDescr, "vrp") {
		dc.vendor = "Huawei"
	} else if strings.Contains(sysDescr, "hp ") || strings.Contains(sysDescr, "procurve") {
		dc.vendor = "HP"
	} else if strings.Contains(sysDescr, "dell") || strings.Contains(sysDescr, "force10") {
		dc.vendor = "Dell"
	} else if strings.Contains(sysDescr, "juniper") || strings.Contains(sysDescr, "junos") {
		dc.vendor = "Juniper"
	} else if strings.Contains(sysDescr, "linux") {
		dc.vendor = "Linux"
	} else if strings.Contains(sysDescr, "windows") {
		dc.vendor = "Windows"
	} else {
		dc.vendor = "Unknown"
		log.Printf("[%s] 无法识别厂商,sysDescr: %s\n", dc.device.DeviceName, sysDescr)
	}

	log.Printf("[%s] 识别厂商: %s\n", dc.device.DeviceName, dc.vendor)
}

// collectBasicInfo 采集基础系统信息(所有设备通用)
func (dc *DeviceCollector) collectBasicInfo() {
	sysOids := []string{
		collector.OID_SYS_DESCR,
		collector.OID_SYS_NAME,
		collector.OID_SYS_UPTIME,
	}

	sysData, err := dc.collector.Get(sysOids)
	if err != nil {
		return
	}

	// 系统描述和名称
	sysDescr := collector.ToString(collector.GetOIDValue(sysData, collector.OID_SYS_DESCR))
	sysName := collector.ToString(collector.GetOIDValue(sysData, collector.OID_SYS_NAME))
	if sysDescr != "" {
		dc.sysDescr.WithLabelValues(sysDescr, sysName).Set(1)
	}

	// 运行时间(已由parseValue转换为秒,直接除以86400转换为天)
	uptime := collector.ToInt64(collector.GetOIDValue(sysData, collector.OID_SYS_UPTIME))
	dc.uptimeGauge.Set(float64(uptime) / 86400.0)
}

// collectSwitchMetrics 采集交换机指标(通用策略:尝试所有厂商OID)
func (dc *DeviceCollector) collectSwitchMetrics() {
	// 1. 采集基础信息
	dc.collectBasicInfo()

	// 2. 采集CPU/内存/温度(尝试所有常见厂商的OID)
	dc.collectSwitchCPU()
	dc.collectSwitchMemory()
	dc.collectSwitchTemperature()

	// 3. 采集端口流量(所有交换机通用)
	dc.collectPortMetrics()
}

// collectH3CMetrics 采集H3C交换机专有指标
func (dc *DeviceCollector) collectH3CMetrics() {
	// H3C CPU使用率(Walk方式,支持多CPU)
	cpuData, err := dc.collector.BulkWalk(collector.OID_H3C_CPU_USAGE)
	if err == nil && len(cpuData) > 0 {
		totalCPU := 0.0
		for i, pdu := range cpuData {
			cpuUsage := collector.ToFloat64(pdu.Value)
			// 记录每个CPU的使用率
			dc.cpuCoreGauge.WithLabelValues(fmt.Sprintf("%d", i+1)).Set(cpuUsage)
			totalCPU += cpuUsage
		}
		// 记录平均CPU使用率
		avgCPU := totalCPU / float64(len(cpuData))
		dc.cpuGauge.Set(avgCPU)
		log.Printf("[%s] H3C CPU: %.2f%% (共%d个CPU)\n", dc.device.DeviceName, avgCPU, len(cpuData))
	} else {
		log.Printf("[%s] H3C CPU采集失败: %v\n", dc.device.DeviceName, err)
	}

	// H3C 内存使用率
	memData, err := dc.collector.BulkWalk(collector.OID_H3C_MEM_USAGE)
	if err == nil && len(memData) > 0 {
		// 取第一个内存模块的使用率
		memUsage := collector.ToFloat64(memData[0].Value)
		dc.memoryGauge.Set(memUsage)
		log.Printf("[%s] H3C Memory: %.2f%%\n", dc.device.DeviceName, memUsage)
	} else {
		log.Printf("[%s] H3C Memory采集失败: %v\n", dc.device.DeviceName, err)
	}

	// H3C 设备温度
	tempData, err := dc.collector.BulkWalk(collector.OID_H3C_DEVICE_TEMP)
	if err == nil && len(tempData) > 0 {
		temp := collector.ToFloat64(tempData[0].Value)
		dc.tempGauge.Set(temp)
		log.Printf("[%s] H3C Temperature: %.1f°C\n", dc.device.DeviceName, temp)
	}
}

// collectCiscoMetrics 采集Cisco交换机专有指标
func (dc *DeviceCollector) collectCiscoMetrics() {
	// Cisco CPU使用率(5秒)
	cpuData, err := dc.collector.Get([]string{collector.OID_CISCO_CPU_5SEC})
	if err == nil {
		cpuUsage := collector.ToFloat64(collector.GetOIDValue(cpuData, collector.OID_CISCO_CPU_5SEC))
		dc.cpuGauge.Set(cpuUsage)
		log.Printf("[%s] Cisco CPU: %.2f%%\n", dc.device.DeviceName, cpuUsage)
	}

	// Cisco内存使用率需要计算 (已用/总量)
	// 这里简化处理,实际需要Walk内存池
	log.Printf("[%s] Cisco Memory采集需要额外实现\n", dc.device.DeviceName)
}

// collectHuaweiMetrics 采集华为交换机专有指标
func (dc *DeviceCollector) collectHuaweiMetrics() {
	// 华为CPU使用率
	cpuData, err := dc.collector.BulkWalk(collector.OID_HUAWEI_CPU)
	if err == nil && len(cpuData) > 0 {
		cpuUsage := collector.ToFloat64(cpuData[0].Value)
		dc.cpuGauge.Set(cpuUsage)
		log.Printf("[%s] Huawei CPU: %.2f%%\n", dc.device.DeviceName, cpuUsage)
	}

	// 华为内存使用率
	memData, err := dc.collector.BulkWalk(collector.OID_HUAWEI_MEM_USED)
	if err == nil && len(memData) > 0 {
		memUsage := collector.ToFloat64(memData[0].Value)
		dc.memoryGauge.Set(memUsage)
		log.Printf("[%s] Huawei Memory: %.2f%%\n", dc.device.DeviceName, memUsage)
	}
}

// collectServerMetrics 采集服务器指标(Windows/Linux,使用Host Resources MIB)
func (dc *DeviceCollector) collectServerMetrics() {
	// 1. 采集基础信息
	dc.collectBasicInfo()

	// 2. 采集服务器专有指标(保留原有逻辑)
	dc.collectBasicMetrics()
}

// ============= 通用采集方法(尝试所有厂商OID) =============

// collectSwitchCPU 采集交换机CPU(尝试所有厂商OID,哪个成功用哪个)
func (dc *DeviceCollector) collectSwitchCPU() {
	collected := false
	var cpuValue float64

	// 特殊处理H3C: 需要先获取实体索引
	if dc.vendor == "H3C" {
		cpuData, err := dc.collector.BulkWalk(collector.OID_H3C_ENTITY_CPU_BASE)
		if err == nil && len(cpuData) > 0 {
			// 找到第一个有效的CPU实例
			for _, pdu := range cpuData {
				cpuUsage := collector.ToFloat64(pdu.Value)
				if cpuUsage > 0 && cpuUsage <= 100 {
					cpuValue = cpuUsage
					collected = true
					break
				}
			}
		}
	}

	// 如果H3C采集失败,或者不是H3C设备,尝试其他厂商
	if !collected {
		// 定义其他厂商的CPU OID
		cpuOIDs := []struct {
			vendor string
			oid    string
			isWalk bool
		}{
			{"Cisco", collector.OID_CISCO_CPU_5SEC, false},
			{"Huawei", collector.OID_HUAWEI_CPU, true},
			{"HP", collector.OID_HP_CPU, false},
			{"Dell", collector.OID_DELL_CPU_5SEC, true},
			{"Juniper", collector.OID_JUNIPER_CPU, true},
		}

		for _, oidInfo := range cpuOIDs {
			if oidInfo.isWalk {
				cpuData, err := dc.collector.BulkWalk(oidInfo.oid)
				if err == nil && len(cpuData) > 0 {
					totalCPU := 0.0
					validCount := 0

					for _, pdu := range cpuData {
						cpuUsage := collector.ToFloat64(pdu.Value)
						if cpuUsage > 0 && cpuUsage < 1 {
							cpuUsage = cpuUsage * 100
						}
						if cpuUsage > 0 && cpuUsage <= 100 {
							totalCPU += cpuUsage
							validCount++
						}
					}
					if validCount > 0 {
						cpuValue = totalCPU / float64(validCount)
						collected = true
						break
					}
				}
			} else {
				cpuData, err := dc.collector.Get([]string{oidInfo.oid})
				if err == nil {
					usage := collector.ToFloat64(collector.GetOIDValue(cpuData, oidInfo.oid))
					if usage > 0 && usage < 1 {
						usage = usage * 100
					}
					if usage > 0 && usage <= 100 {
						cpuValue = usage
						collected = true
						break
					}
				}
			}
		}
	}

	if collected {
		dc.cpuGauge.Set(cpuValue)
	} else {
		log.Printf("[%s] ✗ CPU采集失败: 所有厂商OID均无响应\n", dc.device.DeviceName)
	}
}

// collectSwitchMemory 采集交换机内存(尝试所有厂商OID)
func (dc *DeviceCollector) collectSwitchMemory() {
	collected := false
	var memValue float64

	// 特殊处理H3C: 需要先获取实体索引
	if dc.vendor == "H3C" {
		memData, err := dc.collector.BulkWalk(collector.OID_H3C_ENTITY_MEM_BASE)
		if err == nil && len(memData) > 0 {
			// 找到第一个有效的内存实例
			for _, pdu := range memData {
				memUsage := collector.ToFloat64(pdu.Value)
				if memUsage > 0 && memUsage <= 100 {
					memValue = memUsage
					collected = true
					break
				}
			}
		}
	}

	// 如果H3C采集失败,或者不是H3C设备,尝试其他厂商
	if !collected {
		// 定义其他厂商的内存OID
		memOIDs := []struct {
			vendor string
			oid    string
			isWalk bool
		}{
			{"Huawei", collector.OID_HUAWEI_MEM_USED, true},
			{"HP", collector.OID_HP_MEM, true},
			{"Dell", collector.OID_DELL_MEM_USED, true},
			{"Juniper", collector.OID_JUNIPER_MEM, true},
		}

		for _, oidInfo := range memOIDs {
			if oidInfo.isWalk {
				memData, err := dc.collector.BulkWalk(oidInfo.oid)
				if err == nil && len(memData) > 0 {
					usage := collector.ToFloat64(memData[0].Value)
					if usage > 0 && usage < 1 {
						usage = usage * 100
					}
					if usage > 0 && usage <= 100 {
						memValue = usage
						collected = true
						break
					}
				}
			} else {
				memData, err := dc.collector.Get([]string{oidInfo.oid})
				if err == nil {
					usage := collector.ToFloat64(collector.GetOIDValue(memData, oidInfo.oid))
					if usage > 0 && usage < 1 {
						usage = usage * 100
					}
					if usage > 0 && usage <= 100 {
						memValue = usage
						collected = true
						break
					}
				}
			}
		}

		// 尝试Cisco内存(需要计算使用率)
		if !collected {
			memUsed, errUsed := dc.collector.BulkWalk(collector.OID_CISCO_MEM_POOL_USED)
			memFree, errFree := dc.collector.BulkWalk(collector.OID_CISCO_MEM_POOL_FREE)
			if errUsed == nil && errFree == nil && len(memUsed) > 0 && len(memFree) > 0 {
				used := collector.ToFloat64(memUsed[0].Value)
				free := collector.ToFloat64(memFree[0].Value)
				if used+free > 0 {
					memValue = (used / (used + free)) * 100.0
					collected = true
				}
			}
		}
	}

	if collected {
		dc.memoryGauge.Set(memValue)
	} else {
		log.Printf("[%s] ✗ Memory采集失败: 所有厂商OID均无响应\n", dc.device.DeviceName)
	}
}

// collectSwitchTemperature 采集交换机温度(尝试所有厂商OID)
func (dc *DeviceCollector) collectSwitchTemperature() {
	collected := false
	var tempValue float64

	// 特殊处理H3C: 需要先获取实体索引
	if dc.vendor == "H3C" {
		tempData, err := dc.collector.BulkWalk(collector.OID_H3C_ENTITY_TEMP_BASE)
		if err == nil && len(tempData) > 0 {
			// 找到第一个有效的温度实例
			for _, pdu := range tempData {
				temp := collector.ToFloat64(pdu.Value)
				// 温度通常在 -50 到 150 摄氏度之间
				if temp > -50 && temp < 150 {
					tempValue = temp
					collected = true
					break
				}
			}
		}
	}

	// 如果H3C采集失败,或者不是H3C设备,尝试其他厂商
	if !collected {
		tempOIDs := []struct {
			vendor string
			oid    string
		}{
			{"Huawei", collector.OID_HUAWEI_TEMP},
			{"HP", collector.OID_HP_TEMP},
			{"Dell", collector.OID_DELL_TEMP},
			{"Juniper", collector.OID_JUNIPER_TEMP},
		}

		for _, oidInfo := range tempOIDs {
			tempData, err := dc.collector.BulkWalk(oidInfo.oid)
			if err == nil && len(tempData) > 0 {
				temp := collector.ToFloat64(tempData[0].Value)
				if temp > -50 && temp < 150 {
					tempValue = temp
					collected = true
					break
				}
			}
		}
	}

	if collected {
		dc.tempGauge.Set(tempValue)
	}
	// 温度不是必需指标,失败不打印错误日志
}
