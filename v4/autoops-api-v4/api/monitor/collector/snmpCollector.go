package collector

import (
	"fmt"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
)

// SNMPCollector SNMP采集器
type SNMPCollector struct {
	Target    string
	Port      uint16
	Community string
	Version   gosnmp.SnmpVersion
	Timeout   time.Duration
	Retries   int
}

// OID定义 - 通用MIB-II标准
const (
	// System组 (1.3.6.1.2.1.1)
	OID_SYS_DESCR    = "1.3.6.1.2.1.1.1.0"    // 系统描述
	OID_SYS_OBJECTID = "1.3.6.1.2.1.1.2.0"    // 系统OID
	OID_SYS_UPTIME   = "1.3.6.1.2.1.1.3.0"    // 系统运行时间
	OID_SYS_CONTACT  = "1.3.6.1.2.1.1.4.0"    // 系统联系人
	OID_SYS_NAME     = "1.3.6.1.2.1.1.5.0"    // 系统名称
	OID_SYS_LOCATION = "1.3.6.1.2.1.1.6.0"    // 系统位置

	// Interface组 (1.3.6.1.2.1.2)
	OID_IF_NUMBER     = "1.3.6.1.2.1.2.1.0"      // 接口数量
	OID_IF_TABLE      = "1.3.6.1.2.1.2.2"        // 接口表
	OID_IF_INDEX      = "1.3.6.1.2.1.2.2.1.1"    // 接口索引
	OID_IF_DESCR      = "1.3.6.1.2.1.2.2.1.2"    // 接口描述
	OID_IF_TYPE       = "1.3.6.1.2.1.2.2.1.3"    // 接口类型
	OID_IF_MTU        = "1.3.6.1.2.1.2.2.1.4"    // 接口MTU
	OID_IF_SPEED      = "1.3.6.1.2.1.2.2.1.5"    // 接口速度
	OID_IF_PHYS_ADDR  = "1.3.6.1.2.1.2.2.1.6"    // 物理地址(MAC)
	OID_IF_ADMIN_STA  = "1.3.6.1.2.1.2.2.1.7"    // 管理状态
	OID_IF_OPER_STA   = "1.3.6.1.2.1.2.2.1.8"    // 操作状态
	OID_IF_IN_OCTETS  = "1.3.6.1.2.1.2.2.1.10"   // 入站字节数
	OID_IF_IN_PKTS    = "1.3.6.1.2.1.2.2.1.11"   // 入站包数
	OID_IF_IN_ERRORS  = "1.3.6.1.2.1.2.2.1.14"   // 入站错误
	OID_IF_IN_DISCARD = "1.3.6.1.2.1.2.2.1.13"   // 入站丢弃
	OID_IF_OUT_OCTETS = "1.3.6.1.2.1.2.2.1.16"   // 出站字节数
	OID_IF_OUT_PKTS   = "1.3.6.1.2.1.2.2.1.17"   // 出站包数
	OID_IF_OUT_ERRORS = "1.3.6.1.2.1.2.2.1.20"   // 出站错误
	OID_IF_OUT_DISCA  = "1.3.6.1.2.1.2.2.1.19"   // 出站丢弃
	OID_IF_ALIAS      = "1.3.6.1.2.1.31.1.1.1.18" // 接口别名(ifXTable)
	OID_IF_HIGH_SPEED = "1.3.6.1.2.1.31.1.1.1.15" // 高速端口速率(Mbps,用于10G+端口)

	// IP组 (1.3.6.1.2.1.4)
	OID_IP_FORWARDING = "1.3.6.1.2.1.4.1.0" // IP转发

	// Host Resources MIB - CPU和内存
	OID_HR_PROCESSOR_LOAD = "1.3.6.1.2.1.25.3.3.1.2"    // CPU负载
	OID_HR_STORAGE_TABLE  = "1.3.6.1.2.1.25.2.3"        // 存储表
	OID_HR_STORAGE_INDEX  = "1.3.6.1.2.1.25.2.3.1.1"    // 存储索引
	OID_HR_STORAGE_TYPE   = "1.3.6.1.2.1.25.2.3.1.2"    // 存储类型
	OID_HR_STORAGE_DESCR  = "1.3.6.1.2.1.25.2.3.1.3"    // 存储描述
	OID_HR_STORAGE_UNIT   = "1.3.6.1.2.1.25.2.3.1.4"    // 分配单元大小
	OID_HR_STORAGE_SIZE   = "1.3.6.1.2.1.25.2.3.1.5"    // 总大小
	OID_HR_STORAGE_USED   = "1.3.6.1.2.1.25.2.3.1.6"    // 已使用
	OID_HR_SW_RUN_INDEX   = "1.3.6.1.2.1.25.4.2.1.1"    // 运行进程索引

	// UCD-SNMP-MIB (NET-SNMP) - Linux系统常用
	OID_LOAD_1MIN  = "1.3.6.1.4.1.2021.10.1.3.1"  // 1分钟负载
	OID_LOAD_5MIN  = "1.3.6.1.4.1.2021.10.1.3.2"  // 5分钟负载
	OID_LOAD_15MIN = "1.3.6.1.4.1.2021.10.1.3.3"  // 15分钟负载
	OID_MEM_TOTAL  = "1.3.6.1.4.1.2021.4.5.0"     // 总内存
	OID_MEM_AVAIL  = "1.3.6.1.4.1.2021.4.6.0"     // 可用内存
	OID_MEM_BUFFER = "1.3.6.1.4.1.2021.4.14.0"    // 缓冲区
	OID_MEM_CACHED = "1.3.6.1.4.1.2021.4.15.0"    // 缓存

	// Cisco私有MIB
	OID_CISCO_CPU_5SEC     = "1.3.6.1.4.1.9.2.1.56.0"      // Cisco 5秒CPU使用率
	OID_CISCO_CPU_1MIN     = "1.3.6.1.4.1.9.2.1.57.0"      // Cisco 1分钟CPU使用率
	OID_CISCO_MEM_POOL_USED = "1.3.6.1.4.1.9.9.48.1.1.1.5" // Cisco 已用内存
	OID_CISCO_MEM_POOL_FREE = "1.3.6.1.4.1.9.9.48.1.1.1.6" // Cisco 可用内存

	// H3C私有MIB (HH3C-ENTITY-EXT-MIB)
	// 注意: H3C的OID需要加上实体索引,通常是.65537
	OID_H3C_CPU_USAGE       = "1.3.6.1.4.1.25506.2.6.1.1.1.1.6.1"  // H3C CPU使用率 (hh3cEntityCpuUsage) - 废弃,改用Walk
	OID_H3C_MEM_USAGE       = "1.3.6.1.4.1.25506.2.6.1.1.1.1.8.1"  // H3C 内存使用率 (hh3cEntityMemUsage) - 废弃,改用Walk
	OID_H3C_MEM_SIZE        = "1.3.6.1.4.1.25506.2.6.1.1.1.1.10.1" // H3C 内存总大小
	OID_H3C_MEM_FREE        = "1.3.6.1.4.1.25506.2.6.1.1.1.1.11.1" // H3C 内存可用大小
	OID_H3C_DEVICE_TEMP     = "1.3.6.1.4.1.25506.2.6.1.1.1.1.12.1" // H3C 设备温度
	OID_H3C_ENTITY_CPU_BASE = "1.3.6.1.4.1.25506.2.6.1.1.1.1.6"    // H3C 实体CPU基础OID (用于Walk)
	OID_H3C_ENTITY_MEM_BASE = "1.3.6.1.4.1.25506.2.6.1.1.1.1.8"    // H3C 实体内存基础OID (用于Walk)
	OID_H3C_ENTITY_TEMP_BASE = "1.3.6.1.4.1.25506.2.6.1.1.1.1.12"  // H3C 实体温度基础OID (用于Walk)

	// Huawei私有MIB
	OID_HUAWEI_CPU      = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.5" // Huawei CPU使用率
	OID_HUAWEI_MEM_USED = "1.3.6.1.4.1.2011.6.3.5.1.1.2"       // Huawei 已用内存
	OID_HUAWEI_MEM_FREE = "1.3.6.1.4.1.2011.6.3.5.1.1.3"       // Huawei 可用内存
	OID_HUAWEI_TEMP     = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.11" // Huawei 设备温度

	// HP/HPE ProCurve私有MIB
	OID_HP_CPU    = "1.3.6.1.4.1.11.2.14.11.5.1.9.6.1.0"    // HP CPU使用率
	OID_HP_MEM    = "1.3.6.1.4.1.11.2.14.11.5.1.1.2.2.1.1.7" // HP 内存使用率
	OID_HP_TEMP   = "1.3.6.1.4.1.11.2.14.11.1.2.6.1.4"      // HP 设备温度

	// Dell/Force10私有MIB
	OID_DELL_CPU_5SEC = "1.3.6.1.4.1.6027.3.10.1.2.9.1.3"  // Dell 5秒CPU使用率
	OID_DELL_MEM_USED = "1.3.6.1.4.1.6027.3.10.1.2.2.1.4"  // Dell 已用内存
	OID_DELL_TEMP     = "1.3.6.1.4.1.6027.3.10.1.2.3.1.4"  // Dell 设备温度

	// Juniper私有MIB
	OID_JUNIPER_CPU = "1.3.6.1.4.1.2636.3.1.13.1.8"    // Juniper CPU使用率
	OID_JUNIPER_MEM = "1.3.6.1.4.1.2636.3.1.13.1.11"   // Juniper 内存使用率
	OID_JUNIPER_TEMP = "1.3.6.1.4.1.2636.3.1.13.1.7"   // Juniper 设备温度
)

// NewSNMPCollector 创建SNMP采集器
func NewSNMPCollector(target string, port uint16, community string, version string, timeout int, retries int) *SNMPCollector {
	var snmpVersion gosnmp.SnmpVersion
	switch version {
	case "v1":
		snmpVersion = gosnmp.Version1
	case "v2c":
		snmpVersion = gosnmp.Version2c
	case "v3":
		snmpVersion = gosnmp.Version3
	default:
		snmpVersion = gosnmp.Version2c
	}

	return &SNMPCollector{
		Target:    target,
		Port:      port,
		Community: community,
		Version:   snmpVersion,
		Timeout:   time.Duration(timeout) * time.Second,
		Retries:   retries,
	}
}

// Connect 连接到SNMP设备
func (c *SNMPCollector) Connect() (*gosnmp.GoSNMP, error) {
	snmp := &gosnmp.GoSNMP{
		Target:    c.Target,
		Port:      c.Port,
		Community: c.Community,
		Version:   c.Version,
		Timeout:   c.Timeout,
		Retries:   c.Retries,
	}

	err := snmp.Connect()
	if err != nil {
		return nil, fmt.Errorf("SNMP连接失败: %v", err)
	}

	return snmp, nil
}

// Get 获取单个或多个OID的值
func (c *SNMPCollector) Get(oids []string) (map[string]interface{}, error) {
	snmp, err := c.Connect()
	if err != nil {
		return nil, err
	}
	defer snmp.Conn.Close()

	result, err := snmp.Get(oids)
	if err != nil {
		return nil, fmt.Errorf("SNMP GET失败: %v", err)
	}

	metrics := make(map[string]interface{})
	for _, variable := range result.Variables {
		metrics[variable.Name] = c.parseValue(variable)
	}

	return metrics, nil
}

// Walk 遍历OID树
func (c *SNMPCollector) Walk(oid string) ([]gosnmp.SnmpPDU, error) {
	snmp, err := c.Connect()
	if err != nil {
		return nil, err
	}
	defer snmp.Conn.Close()

	var pdus []gosnmp.SnmpPDU
	err = snmp.Walk(oid, func(pdu gosnmp.SnmpPDU) error {
		pdus = append(pdus, pdu)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("SNMP WALK失败: %v", err)
	}

	return pdus, nil
}

// BulkWalk 批量遍历OID树(更高效,仅v2c和v3支持)
func (c *SNMPCollector) BulkWalk(oid string) ([]gosnmp.SnmpPDU, error) {
	if c.Version == gosnmp.Version1 {
		// v1不支持BulkWalk,降级为普通Walk
		return c.Walk(oid)
	}

	snmp, err := c.Connect()
	if err != nil {
		return nil, err
	}
	defer snmp.Conn.Close()

	var pdus []gosnmp.SnmpPDU
	err = snmp.BulkWalk(oid, func(pdu gosnmp.SnmpPDU) error {
		pdus = append(pdus, pdu)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("SNMP BULKWALK失败: %v", err)
	}

	return pdus, nil
}

// TestConnection 测试SNMP连接
func (c *SNMPCollector) TestConnection() error {
	snmp, err := c.Connect()
	if err != nil {
		return err
	}
	defer snmp.Conn.Close()

	// 尝试获取系统名称
	result, err := snmp.Get([]string{OID_SYS_NAME})
	if err != nil {
		return fmt.Errorf("SNMP测试失败: %v", err)
	}

	if len(result.Variables) == 0 {
		return fmt.Errorf("SNMP响应为空")
	}

	return nil
}

// GetSystemInfo 获取系统基本信息
func (c *SNMPCollector) GetSystemInfo() (map[string]string, error) {
	oids := []string{
		OID_SYS_DESCR,
		OID_SYS_OBJECTID,
		OID_SYS_UPTIME,
		OID_SYS_CONTACT,
		OID_SYS_NAME,
		OID_SYS_LOCATION,
	}

	rawData, err := c.Get(oids)
	if err != nil {
		return nil, err
	}

	info := make(map[string]string)
	for oid, value := range rawData {
		info[oid] = fmt.Sprintf("%v", value)
	}

	return info, nil
}

// parseValue 解析SNMP返回值
func (c *SNMPCollector) parseValue(variable gosnmp.SnmpPDU) interface{} {
	switch variable.Type {
	case gosnmp.OctetString:
		// 尝试转换为字符串
		bytes := variable.Value.([]byte)
		return string(bytes)
	case gosnmp.Integer, gosnmp.Gauge32, gosnmp.Counter32, gosnmp.Counter64, gosnmp.Uinteger32:
		return variable.Value
	case gosnmp.TimeTicks:
		// 转换为秒
		ticks := variable.Value.(uint32)
		return ticks / 100 // TimeTicks是百分之一秒
	case gosnmp.IPAddress:
		return variable.Value
	default:
		return variable.Value
	}
}

// GetOIDValue 辅助方法:从map中安全获取OID值
func GetOIDValue(data map[string]interface{}, oid string) interface{} {
	// 先尝试直接匹配
	if val, ok := data[oid]; ok {
		return val
	}

	// 尝试添加前导点匹配
	if val, ok := data["."+oid]; ok {
		return val
	}

	// 尝试移除前导点匹配
	if strings.HasPrefix(oid, ".") {
		if val, ok := data[oid[1:]]; ok {
			return val
		}
	}

	return nil
}

// ToInt64 辅助方法:将interface{}转换为int64
func ToInt64(value interface{}) int64 {
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	default:
		return 0
	}
}

// ToFloat64 辅助方法:将interface{}转换为float64
func ToFloat64(value interface{}) float64 {
	if value == nil {
		return 0.0
	}

	switch v := value.(type) {
	case float32:
		return float64(v)
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	default:
		return 0.0
	}
}

// ToString 辅助方法:将interface{}转换为string
func ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	// 特殊处理字节数组,转换为字符串
	if bytes, ok := value.([]byte); ok {
		return string(bytes)
	}
	return fmt.Sprintf("%v", value)
}
