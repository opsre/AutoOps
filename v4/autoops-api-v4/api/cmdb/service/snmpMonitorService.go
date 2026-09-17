package service

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SnmpMonitorService struct {
	deviceDao      *dao.CmdbSnmpDeviceDao
	pushgatewayURL string
}

func NewSnmpMonitorService() *SnmpMonitorService {
	pushgatewayURL := "http://8.130.14.34:9091"
	if config.Config != nil && config.Config.Monitor.Pushgateway.URL != "" {
		pushgatewayURL = config.Config.Monitor.Pushgateway.URL
	}

	return &SnmpMonitorService{
		deviceDao:      dao.NewCmdbSnmpDeviceDao(),
		pushgatewayURL: pushgatewayURL,
	}
}

// GetDeviceMonitorMetrics 获取设备基础监控指标
func (s *SnmpMonitorService) GetDeviceMonitorMetrics(deviceID uint) (*model.DeviceMonitorMetrics, error) {
	// 1. 从数据库获取设备信息
	device, err := s.deviceDao.GetSnmpDeviceById(deviceID)
	if err != nil {
		return nil, fmt.Errorf("设备不存在: %v", err)
	}

	// 2. 构建返回结果
	metrics := &model.DeviceMonitorMetrics{
		DeviceID:       device.ID,
		DeviceName:     device.DeviceName,
		DeviceType:     device.DeviceType,
		IPAddress:      device.IPAddress,
		Status:         device.Status,
		LastUpdateTime: device.LastSeen.String(),
	}

	// 3. 如果设备离线,直接返回基础信息
	if device.Status == "offline" {
		return metrics, nil
	}

	// 4. 从Pushgateway获取监控数据
	jobName := fmt.Sprintf("snmp_%s", device.DeviceName)
	metricsData, err := s.fetchMetricsFromPushgateway(jobName)
	if err != nil {
		log.Printf("[GetDeviceMonitorMetrics] 从Pushgateway获取数据失败: %v", err)
		return metrics, nil // 返回基础信息,不报错
	}

	// 5. 解析监控指标
	metrics.CpuUsage = s.parseMetric(metricsData, "snmp_device_cpu_percent")
	metrics.MemoryUsage = s.parseMetric(metricsData, "snmp_device_memory_percent")
	metrics.Temperature = s.parseMetric(metricsData, "snmp_device_temperature_celsius")
	metrics.UptimeDays = s.parseMetric(metricsData, "snmp_device_uptime_days")

	// 6. 从sysDescr解析设备型号
	metrics.DeviceModel = s.parseDeviceModel(metricsData)

	// 7. 解析端口流量统计
	portInTop5, portOutTop5, idleCount, totalCount := s.parsePortTraffic(metricsData)
	metrics.PortInTop5 = portInTop5
	metrics.PortOutTop5 = portOutTop5
	metrics.IdlePortCount = idleCount
	metrics.TotalPortCount = totalCount

	// 8. 获取历史趋势数据(最近24小时)
	prometheusURL := "http://8.130.14.34:9090"
	if config.Config != nil && config.Config.Monitor.Prometheus.URL != "" {
		prometheusURL = config.Config.Monitor.Prometheus.URL
	}

	metrics.CpuTrend = s.queryTrendData(prometheusURL, fmt.Sprintf(`snmp_device_cpu_percent{device_name="%s"}`, device.DeviceName))
	metrics.MemoryTrend = s.queryTrendData(prometheusURL, fmt.Sprintf(`snmp_device_memory_percent{device_name="%s"}`, device.DeviceName))
	metrics.TemperatureTrend = s.queryTrendData(prometheusURL, fmt.Sprintf(`snmp_device_temperature_celsius{device_name="%s"}`, device.DeviceName))

	return metrics, nil
}

// fetchMetricsFromPushgateway 从Pushgateway获取原始指标数据
func (s *SnmpMonitorService) fetchMetricsFromPushgateway(jobName string) (string, error) {
	// Pushgateway API: GET /metrics
	url := fmt.Sprintf("%s/metrics", s.pushgatewayURL)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求Pushgateway失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Pushgateway返回错误状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 过滤出当前job的指标
	return s.filterMetricsByJob(string(body), jobName), nil
}

// filterMetricsByJob 过滤出指定job的指标
func (s *SnmpMonitorService) filterMetricsByJob(allMetrics, jobName string) string {
	var result strings.Builder
	lines := strings.Split(allMetrics, "\n")

	for _, line := range lines {
		// 跳过注释和空行
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// 检查是否包含当前job
		if strings.Contains(line, fmt.Sprintf(`job="%s"`, jobName)) {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	return result.String()
}

// parseMetric 从Prometheus格式的指标文本中解析指定指标的值
func (s *SnmpMonitorService) parseMetric(metricsData, metricName string) float64 {
	// 示例行: snmp_device_cpu_percent{device_id="1",device_name="H3C-Switch",device_type="switch",ip="10.7.80.101"} 15.5
	lines := strings.Split(metricsData, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, metricName) {
			// 提取数值部分
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				value, err := strconv.ParseFloat(parts[len(parts)-1], 64)
				if err == nil {
					return value
				}
			}
		}
	}

	return 0
}

// parseDeviceModel 从sysDescr中解析设备型号
func (s *SnmpMonitorService) parseDeviceModel(metricsData string) string {
	// 示例: snmp_device_info{device_id="1",device_name="H3C-Switch",device_type="switch",ip="10.7.80.101",sys_descr="H3C Comware Platform Software, Software Version 7.1.070, Release 1234P05",sys_name="H3C-Switch"} 1

	lines := strings.Split(metricsData, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "snmp_device_info") {
			// 使用正则提取sys_descr
			re := regexp.MustCompile(`sys_descr="([^"]*)"`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				sysDescr := matches[1]
				// 尝试从sysDescr中提取型号
				return extractModelFromSysDescr(sysDescr)
			}
		}
	}

	return "Unknown"
}

// extractModelFromSysDescr 从sysDescr中提取设备型号
func extractModelFromSysDescr(sysDescr string) string {
	// H3C设备: 匹配 S5110V2-28P-HPWR 这种格式
	re := regexp.MustCompile(`(S\d+[A-Z0-9-]+)`)
	matches := re.FindStringSubmatch(sysDescr)
	if len(matches) > 1 {
		return matches[1]
	}

	// Cisco设备: 匹配 Cisco Catalyst 2960
	if strings.Contains(strings.ToLower(sysDescr), "cisco") {
		re := regexp.MustCompile(`Cisco\s+(\w+\s+\d+\w*)`)
		matches := re.FindStringSubmatch(sysDescr)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// 如果无法提取型号,返回前50个字符
	if len(sysDescr) > 50 {
		return sysDescr[:50] + "..."
	}
	return sysDescr
}

// GetDeviceMonitorMetricsJSON 获取设备监控指标(从Prometheus返回JSON格式)
// 备用方案: 如果有Prometheus查询接口,可以用这个方法
func (s *SnmpMonitorService) GetDeviceMonitorMetricsJSON(deviceID uint) (map[string]interface{}, error) {
	device, err := s.deviceDao.GetSnmpDeviceById(deviceID)
	if err != nil {
		return nil, fmt.Errorf("设备不存在: %v", err)
	}

	// 构建Prometheus查询
	prometheusURL := "http://8.130.14.34:9090" // Prometheus地址
	queries := map[string]string{
		"cpu":         fmt.Sprintf(`snmp_device_cpu_percent{device_name="%s"}`, device.DeviceName),
		"memory":      fmt.Sprintf(`snmp_device_memory_percent{device_name="%s"}`, device.DeviceName),
		"temperature": fmt.Sprintf(`snmp_device_temperature_celsius{device_name="%s"}`, device.DeviceName),
		"uptime":      fmt.Sprintf(`snmp_device_uptime_days{device_name="%s"}`, device.DeviceName),
	}

	result := make(map[string]interface{})
	result["deviceId"] = device.ID
	result["deviceName"] = device.DeviceName
	result["deviceType"] = device.DeviceType
	result["ipAddress"] = device.IPAddress
	result["status"] = device.Status

	// 查询每个指标
	for key, query := range queries {
		value, err := s.queryPrometheus(prometheusURL, query)
		if err != nil {
			log.Printf("[查询Prometheus] %s 失败: %v", key, err)
			result[key] = 0
		} else {
			result[key] = value
		}
	}

	return result, nil
}

// queryPrometheus 查询Prometheus获取实时指标值
func (s *SnmpMonitorService) queryPrometheus(prometheusURL, query string) (float64, error) {
	url := fmt.Sprintf("%s/api/v1/query?query=%s&time=%d", prometheusURL, query, time.Now().Unix())

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// 解析Prometheus JSON响应
	var promResp struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Value  []interface{}     `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &promResp); err != nil {
		return 0, err
	}

	if promResp.Status != "success" || len(promResp.Data.Result) == 0 {
		return 0, fmt.Errorf("无数据")
	}

	// 提取数值
	if len(promResp.Data.Result[0].Value) >= 2 {
		valueStr := promResp.Data.Result[0].Value[1].(string)
		return strconv.ParseFloat(valueStr, 64)
	}

	return 0, fmt.Errorf("解析失败")
}

// parsePortTraffic 解析端口流量数据,返回入向Top5、出向Top5、闲置端口数、总端口数
func (s *SnmpMonitorService) parsePortTraffic(metricsData string) ([]model.PortTrafficStat, []model.PortTrafficStat, int, int) {
	lines := strings.Split(metricsData, "\n")

	// 存储所有端口的入向和出向流量
	portInMap := make(map[string]*model.PortTrafficStat)  // key: portIndex_portName
	portOutMap := make(map[string]*model.PortTrafficStat) // key: portIndex_portName
	portSpeedMap := make(map[string]float64)              // key: portIndex_portName

	// 正则表达式提取标签
	labelRe := regexp.MustCompile(`port_index="(\d+)".*port_name="([^"]*)"`)
	valueRe := regexp.MustCompile(`}\s+([\d.]+)`)

	// 1. 先解析端口速率
	for _, line := range lines {
		if strings.HasPrefix(line, "snmp_port_speed_mbps") {
			labelMatches := labelRe.FindStringSubmatch(line)
			valueMatches := valueRe.FindStringSubmatch(line)
			if len(labelMatches) > 2 && len(valueMatches) > 1 {
				portIndex := labelMatches[1]
				portName := labelMatches[2]
				key := portIndex + "_" + portName
				speed, _ := strconv.ParseFloat(valueMatches[1], 64)
				portSpeedMap[key] = speed
			}
		}
	}

	// 2. 解析入向流量速率(Mbps)
	for _, line := range lines {
		if strings.HasPrefix(line, "snmp_port_in_rate_mbps") {
			labelMatches := labelRe.FindStringSubmatch(line)
			valueMatches := valueRe.FindStringSubmatch(line)
			if len(labelMatches) > 2 && len(valueMatches) > 1 {
				portIndex, _ := strconv.Atoi(labelMatches[1])
				portName := labelMatches[2]
				trafficRate, _ := strconv.ParseFloat(valueMatches[1], 64)
				key := labelMatches[1] + "_" + portName

				portSpeed := portSpeedMap[key]
				utilization := 0.0
				if portSpeed > 0 {
					utilization = (trafficRate / portSpeed) * 100.0
				}

				portInMap[key] = &model.PortTrafficStat{
					PortIndex:       portIndex,
					PortName:        portName,
					TrafficRateMbps: trafficRate,
					Utilization:     utilization,
					PortSpeed:       portSpeed,
				}
			}
		}
	}

	// 3. 解析出向流量速率(Mbps)
	for _, line := range lines {
		if strings.HasPrefix(line, "snmp_port_out_rate_mbps") {
			labelMatches := labelRe.FindStringSubmatch(line)
			valueMatches := valueRe.FindStringSubmatch(line)
			if len(labelMatches) > 2 && len(valueMatches) > 1 {
				portIndex, _ := strconv.Atoi(labelMatches[1])
				portName := labelMatches[2]
				trafficRate, _ := strconv.ParseFloat(valueMatches[1], 64)
				key := labelMatches[1] + "_" + portName

				portSpeed := portSpeedMap[key]
				utilization := 0.0
				if portSpeed > 0 {
					utilization = (trafficRate / portSpeed) * 100.0
				}

				portOutMap[key] = &model.PortTrafficStat{
					PortIndex:       portIndex,
					PortName:        portName,
					TrafficRateMbps: trafficRate,
					Utilization:     utilization,
					PortSpeed:       portSpeed,
				}
			}
		}
	}

	// 4. 统计总端口数和闲置端口数
	totalPortCount := len(portInMap)
	idlePortCount := 0
	for _, port := range portInMap {
		if port.TrafficRateMbps == 0 {
			idlePortCount++
		}
	}

	// 5. 排序获取Top5
	portInTop5 := s.getTopPorts(portInMap, 5)
	portOutTop5 := s.getTopPorts(portOutMap, 5)

	return portInTop5, portOutTop5, idlePortCount, totalPortCount
}

// getTopPorts 从端口Map中获取流量Top N
func (s *SnmpMonitorService) getTopPorts(portMap map[string]*model.PortTrafficStat, topN int) []model.PortTrafficStat {
	// 转换为切片
	ports := make([]model.PortTrafficStat, 0, len(portMap))
	for _, port := range portMap {
		ports = append(ports, *port)
	}

	// 按流量速率降序排序(冒泡排序,简单实现)
	for i := 0; i < len(ports)-1; i++ {
		for j := 0; j < len(ports)-i-1; j++ {
			if ports[j].TrafficRateMbps < ports[j+1].TrafficRateMbps {
				ports[j], ports[j+1] = ports[j+1], ports[j]
			}
		}
	}

	// 取前N个
	if len(ports) > topN {
		return ports[:topN]
	}
	return ports
}

// queryTrendData 查询Prometheus获取最近24小时的趋势数据
func (s *SnmpMonitorService) queryTrendData(prometheusURL, query string) []model.TrendPoint {
	// 构建Range Query: 最近24小时,每5分钟一个点
	endTime := time.Now().Unix()
	startTime := endTime - 86400 // 24小时前
	step := 300                  // 5分钟

	url := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%d&end=%d&step=%d",
		prometheusURL, query, startTime, endTime, step)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[queryTrendData] 请求Prometheus失败: %v", err)
		return []model.TrendPoint{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[queryTrendData] 读取响应失败: %v", err)
		return []model.TrendPoint{}
	}

	// 解析Prometheus Range Query响应
	var promResp struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Values [][]interface{}   `json:"values"` // [[timestamp, value], [timestamp, value], ...]
			} `json:"result"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &promResp); err != nil {
		log.Printf("[queryTrendData] 解析JSON失败: %v", err)
		return []model.TrendPoint{}
	}

	if promResp.Status != "success" || len(promResp.Data.Result) == 0 {
		log.Printf("[queryTrendData] 无数据")
		return []model.TrendPoint{}
	}

	// 转换为TrendPoint
	values := promResp.Data.Result[0].Values
	trendPoints := make([]model.TrendPoint, 0, len(values))

	for _, v := range values {
		if len(v) >= 2 {
			timestamp := int64(v[0].(float64))
			valueStr := v[1].(string)
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				continue
			}

			// 格式化时间戳为ISO 8601
			timeStr := time.Unix(timestamp, 0).Format(time.RFC3339)
			trendPoints = append(trendPoints, model.TrendPoint{
				Timestamp: timeStr,
				Value:     value,
			})
		}
	}

	return trendPoints
}

// GetDeviceResourceMetrics 获取设备资源使用率(轻量级API,仅返回CPU和内存)
func (s *SnmpMonitorService) GetDeviceResourceMetrics(deviceID uint) (*model.DeviceResourceMetrics, error) {
	// 1. 从数据库获取设备信息
	device, err := s.deviceDao.GetSnmpDeviceById(deviceID)
	if err != nil {
		return nil, fmt.Errorf("设备不存在: %v", err)
	}

	// 2. 构建返回结果
	metrics := &model.DeviceResourceMetrics{
		DeviceID:   device.ID,
		DeviceName: device.DeviceName,
		IPAddress:  device.IPAddress,
		Status:     device.Status,
		UpdateTime: device.LastSeen.String(),
	}

	// 3. 如果设备离线,直接返回基础信息
	if device.Status == "offline" {
		return metrics, nil
	}

	// 4. 从Pushgateway获取监控数据
	jobName := fmt.Sprintf("snmp_%s", device.DeviceName)
	metricsData, err := s.fetchMetricsFromPushgateway(jobName)
	if err != nil {
		log.Printf("[GetDeviceResourceMetrics] 从Pushgateway获取数据失败: %v", err)
		return metrics, nil // 返回基础信息,不报错
	}

	// 5. 解析CPU和内存指标
	metrics.CpuUsage = s.parseMetric(metricsData, "snmp_device_cpu_percent")
	metrics.MemoryUsage = s.parseMetric(metricsData, "snmp_device_memory_percent")

	return metrics, nil
}

// GetBatchDeviceResourceMetrics 批量获取设备资源使用率
func (s *SnmpMonitorService) GetBatchDeviceResourceMetrics(deviceIDs []uint) (*model.BatchResourceMetricsResult, error) {
	if len(deviceIDs) == 0 {
		return &model.BatchResourceMetricsResult{
			Items: []model.DeviceResourceMetrics{},
			Total: 0,
		}, nil
	}

	// 使用channel收集结果
	type result struct {
		metrics *model.DeviceResourceMetrics
		err     error
	}

	resultChan := make(chan result, len(deviceIDs))

	// 并发查询所有设备
	for _, deviceID := range deviceIDs {
		go func(id uint) {
			metrics, err := s.GetDeviceResourceMetrics(id)
			resultChan <- result{metrics: metrics, err: err}
		}(deviceID)
	}

	// 收集结果
	items := make([]model.DeviceResourceMetrics, 0, len(deviceIDs))
	for i := 0; i < len(deviceIDs); i++ {
		res := <-resultChan
		if res.err == nil && res.metrics != nil {
			items = append(items, *res.metrics)
		} else if res.err != nil {
			log.Printf("[GetBatchDeviceResourceMetrics] 查询设备失败: %v", res.err)
		}
	}

	return &model.BatchResourceMetricsResult{
		Items: items,
		Total: len(items),
	}, nil
}
