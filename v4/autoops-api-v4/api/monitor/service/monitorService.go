package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	monitormodel "dodevops-api/api/monitor/model"
	"dodevops-api/common/config"
	"dodevops-api/common/result"
	"dodevops-api/pkg/redis"

	"github.com/gin-gonic/gin"
)

type hostResult struct {
	id      uint
	metrics *monitormodel.HostMetrics
}

// MonitorServiceInterface 定义监控服务接口
type MonitorServiceInterface interface {
	GetHostMetrics(c *gin.Context, hostID uint)                      // 获取主机监控指标
	BatchGetHostMetrics(c *gin.Context, hostIDs []uint)              // 批量获取主机监控指标
	GetHostMetricHistory(c *gin.Context, hostID uint, metric string, // 获取主机指标历史数据
		start, end string, step string)
	GetHostAllMetricsHistory(c *gin.Context, hostID uint, // 获取主机所有指标历史数据
		start, end string, step string)
	GetTopProcesses(c *gin.Context, hostID uint)          // 获取主机TOP进程使用率
	GetHostPorts(c *gin.Context, hostID uint)             // 获取主机端口信息
	GetTopUsageHosts(c *gin.Context, days, limit int)     // 获取使用率TOP主机
	GetCompleteHostMetrics(c *gin.Context, hostID uint)   // 获取主机完整监控数据
}

// MonitorServiceImpl 监控服务实现
type MonitorServiceImpl struct {
	hostDao       dao.CmdbHostDao
	prometheusURL string
	cacheExpire   time.Duration
}

// 缓存键名常量
const (
	hostCacheKey = "monitor:host:%d" // 主机信息缓存键
	cacheTimeout = 5 * time.Minute   // 缓存过期时间
)

// NewMonitorService 创建监控服务实例
func NewMonitorService() MonitorServiceInterface {
	return &MonitorServiceImpl{
		hostDao:       dao.NewCmdbHostDao(),
		prometheusURL: config.Config.Monitor.Prometheus.URL,
		cacheExpire:   cacheTimeout,
	}
}

// getInstanceIdentifier 获取主机的instance标识（优先使用内网IP，其次使用Name）
func (s *MonitorServiceImpl) getInstanceIdentifier(host *model.CmdbHost) string {
	if host.PrivateIP != "" {
		return host.PrivateIP
	}
	if host.Name != "" {
		return host.Name
	}
	return ""
}

// checkHostOnline 检查主机是否在线(通过TCP端口检测或云厂商状态)
func (s *MonitorServiceImpl) checkHostOnline(host model.CmdbHost) bool {
	if host.SSHIP == "" || host.SSHPort == 0 {
		return false
	}

	// 判断是否为云厂商主机（1->自建, 2->阿里云, 3->腾讯云, 4->百度云, 5->华为云, 6->AWS）
	isCloudVendor := host.Vendor >= 2 && host.Vendor <= 6

	if isCloudVendor {
		// 云厂商主机：优先信任云厂商的Status字段
		// Status: 1->认证成功(在线), 2->未认证(在线), 3->认证失败(离线)
		return host.Status == 1 || host.Status == 2
	}

	// 自建主机：使用TCP端口探测，超时500ms
	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", host.SSHIP, host.SSHPort),
		500*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// isPrivateIP 判断是否为私网IP
func (s *MonitorServiceImpl) isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// 检查是否为私网地址段
	// 10.0.0.0/8
	// 172.16.0.0/12
	// 192.168.0.0/16
	// 100.64.0.0/10 (Carrier-grade NAT)
	privateIPBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"100.64.0.0/10",
	}

	for _, block := range privateIPBlocks {
		_, subnet, _ := net.ParseCIDR(block)
		if subnet.Contains(ip) {
			return true
		}
	}

	return false
}

// getHostFromCache 从缓存获取主机信息
func (s *MonitorServiceImpl) getHostFromCache(hostID uint) (*model.CmdbHost, error) {
	cacheKey := fmt.Sprintf(hostCacheKey, hostID)

	// 从 Redis 获取缓存数据
	cachedData, err := redis.RedisDb.Get(context.Background(), cacheKey).Result()
	if err != nil {
		// 缓存未命中或出错，返回 nil
		return nil, err
	}

	// 反序列化主机信息
	var host model.CmdbHost
	if err := json.Unmarshal([]byte(cachedData), &host); err != nil {
		return nil, err
	}

	return &host, nil
}

// setHostToCache 将主机信息存入缓存
func (s *MonitorServiceImpl) setHostToCache(hostID uint, host *model.CmdbHost) error {
	cacheKey := fmt.Sprintf(hostCacheKey, hostID)

	// 序列化主机信息
	hostData, err := json.Marshal(host)
	if err != nil {
		return err
	}

	// 存入 Redis，设置过期时间
	return redis.RedisDb.Set(context.Background(), cacheKey, hostData, s.cacheExpire).Err()
}

// getHostWithCache 优先从缓存获取主机信息，缓存未命中则查询数据库
func (s *MonitorServiceImpl) getHostWithCache(hostID uint) (*model.CmdbHost, error) {
	// 1. 先尝试从缓存获取
	if host, err := s.getHostFromCache(hostID); err == nil && host != nil {
		return host, nil
	}

	// 2. 缓存未命中，查询数据库
	host, err := s.hostDao.GetCmdbHostById(hostID)
	if err != nil {
		return nil, err
	}

	// 3. 将查询结果存入缓存(异步执行，不影响主流程)
	go func() {
		if err := s.setHostToCache(hostID, &host); err != nil {
			fmt.Printf("缓存主机信息失败: %v\n", err)
		}
	}()

	return &host, nil
}

// GetHostMetrics 获取单个主机监控指标
func (s *MonitorServiceImpl) GetHostMetrics(c *gin.Context, hostID uint) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 优先从缓存获取主机信息
	host, err := s.getHostWithCache(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	// 使用内网IP作为instance标识
	instanceIdentifier := s.getInstanceIdentifier(host)
	if instanceIdentifier == "" {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机内网IP和名称均为空")
		return
	}

	metrics := &monitormodel.HostMetrics{}

	// 检查主机是否在线
	if !s.checkHostOnline(*host) {
		metrics.CPUUsage = 0
		metrics.MemoryUsage = 0
		metrics.DiskUsage = 0
		metrics.Load5 = 0
		result.Success(c, metrics)
		return
	}

	// ... existing code ...
	now := time.Now()
	startTime := now.Add(-5 * time.Minute) // 查询最近5分钟的数据
	endTime := now
	step := "30s" // 30秒间隔

	// 查询CPU使用率
	cpuQuery := fmt.Sprintf("system_cpu_usage_percent{instance=\"%s\"}", instanceIdentifier)
	cpu, err := s.queryPrometheus(cpuQuery, startTime, endTime, step)
	if err == nil && len(cpu.Data.Result) > 0 && len(cpu.Data.Result[0].Values) > 0 {
		lastPoint := cpu.Data.Result[0].Values[len(cpu.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				metrics.CPUUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
			}
		}
	} else {
		metrics.CPUUsage = 0
	}

	// 查询内存使用率
	memQuery := fmt.Sprintf("system_memory_usage_percent{instance=\"%s\"}", instanceIdentifier)
	mem, err := s.queryPrometheus(memQuery, startTime, endTime, step)
	if err == nil && len(mem.Data.Result) > 0 && len(mem.Data.Result[0].Values) > 0 {
		lastPoint := mem.Data.Result[0].Values[len(mem.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				metrics.MemoryUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
			}
		}
	} else {
		metrics.MemoryUsage = 0
	}

	// 查询磁盘使用率
	diskQuery := fmt.Sprintf("system_disk_usage_percent{instance=\"%s\",mountpoint=\"/\"}", instanceIdentifier)
	disk, err := s.queryPrometheus(diskQuery, startTime, endTime, step)
	if err == nil && len(disk.Data.Result) > 0 && len(disk.Data.Result[0].Values) > 0 {
		lastPoint := disk.Data.Result[0].Values[len(disk.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				metrics.DiskUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
			}
		}
	} else {
		metrics.DiskUsage = 0
	}

	// 查询5分钟系统负载
	loadQuery := fmt.Sprintf("system_load_average_5min{instance=\"%s\"}", instanceIdentifier)
	load, err := s.queryPrometheus(loadQuery, startTime, endTime, step)
	if err == nil && len(load.Data.Result) > 0 && len(load.Data.Result[0].Values) > 0 {
		lastPoint := load.Data.Result[0].Values[len(load.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if loadValue, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				metrics.Load5, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", loadValue), 64)
			}
		}
	} else {
		metrics.Load5 = 0
	}

	result.Success(c, metrics)
}

// BatchGetHostMetrics 批量获取主机监控指标
func (s *MonitorServiceImpl) BatchGetHostMetrics(c *gin.Context, hostIDs []uint) {
	// 参数校验
	if len(hostIDs) == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 限制并发数量，防止过多并发请求导致数据库压力
	maxConcurrent := 10
	if len(hostIDs) < maxConcurrent {
		maxConcurrent = len(hostIDs)
	}
	semaphore := make(chan struct{}, maxConcurrent)

	// 使用sync.WaitGroup等待所有goroutine完成
	var wg sync.WaitGroup
	resultMap := make(map[uint]*monitormodel.HostMetrics, len(hostIDs))
	var mu sync.Mutex // 用于保护resultMap的并发访问

	// 为每个主机ID启动一个goroutine
	for _, id := range hostIDs {
		wg.Add(1)
		go func(hostID uint) {
			defer wg.Done()

			// 获取信号量，限制并发数
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 优先从缓存获取主机信息
			host, err := s.getHostWithCache(hostID)
			if err != nil {
				// 单个主机失败不影响其他主机
				return
			}

			metrics := &monitormodel.HostMetrics{}

			// 检测主机在线状态
			metrics.OnlineStatus = 1 // 默认离线
			if host.SSHIP != "" && host.SSHPort != 0 {
				// 判断是否为云厂商主机（1->自建, 2->阿里云, 3->腾讯云, 4->百度云, 5->华为云, 6->AWS）
				isCloudVendor := host.Vendor >= 2 && host.Vendor <= 6

				if isCloudVendor {
					// 云厂商主机：优先信任云厂商的Status字段
					// Status: 1->认证成功(在线), 2->未认证(在线), 3->认证失败(离线)
					if host.Status == 1 || host.Status == 2 {
						metrics.OnlineStatus = 0 // 在线
					} else {
						metrics.OnlineStatus = 1 // 离线
					}
				} else {
					// 自建主机：使用TCP端口探测
					conn, err := net.DialTimeout("tcp",
						fmt.Sprintf("%s:%d", host.SSHIP, host.SSHPort),
						500*time.Millisecond)
					if err == nil {
						metrics.OnlineStatus = 0 // 在线
						conn.Close()
					}
				}
			}

			// 如果主机离线，设置监控数据为0
			if metrics.OnlineStatus == 1 {
				metrics.CPUUsage = 0
				metrics.MemoryUsage = 0
				metrics.DiskUsage = 0
				mu.Lock()
				resultMap[hostID] = metrics
				mu.Unlock()
				return
			}

			// 主机在线，获取监控数据（优化查询频率）
			// 获取当前时间
			now := time.Now()
			startTime := now.Add(-3 * time.Minute) // 缩短查询时间范围从5分钟到3分钟
			endTime := now
			step := "60s" // 增大步长从30秒到60秒

			// 使用内网IP作为instance标识
			instanceIdentifier := s.getInstanceIdentifier(host)

			// 优化查询：只获取最后一个数据点
			if instanceIdentifier != "" {
				// 并发查询三个指标
				var metricWg sync.WaitGroup
				var metricMu sync.Mutex

				// 查询CPU使用率
				metricWg.Add(1)
				go func() {
					defer metricWg.Done()
					cpuQuery := fmt.Sprintf("system_cpu_usage_percent{instance=\"%s\"}", instanceIdentifier)
					cpu, err := s.queryPrometheus(cpuQuery, startTime, endTime, step)
					if err == nil && len(cpu.Data.Result) > 0 && len(cpu.Data.Result[0].Values) > 0 {
						lastPoint := cpu.Data.Result[0].Values[len(cpu.Data.Result[0].Values)-1]
						if len(lastPoint) >= 2 {
							if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
								metricMu.Lock()
								metrics.CPUUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
								metricMu.Unlock()
							}
						}
					}
				}()

				// 查询内存使用率
				metricWg.Add(1)
				go func() {
					defer metricWg.Done()
					memQuery := fmt.Sprintf("system_memory_usage_percent{instance=\"%s\"}", instanceIdentifier)
					mem, err := s.queryPrometheus(memQuery, startTime, endTime, step)
					if err == nil && len(mem.Data.Result) > 0 && len(mem.Data.Result[0].Values) > 0 {
						lastPoint := mem.Data.Result[0].Values[len(mem.Data.Result[0].Values)-1]
						if len(lastPoint) >= 2 {
							if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
								metricMu.Lock()
								metrics.MemoryUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
								metricMu.Unlock()
							}
						}
					}
				}()

				// 查询磁盘使用率
				metricWg.Add(1)
				go func() {
					defer metricWg.Done()
					diskQuery := fmt.Sprintf("system_disk_usage_percent{instance=\"%s\",mountpoint=\"/\"}", instanceIdentifier)
					disk, err := s.queryPrometheus(diskQuery, startTime, endTime, step)
					if err == nil && len(disk.Data.Result) > 0 && len(disk.Data.Result[0].Values) > 0 {
						lastPoint := disk.Data.Result[0].Values[len(disk.Data.Result[0].Values)-1]
						if len(lastPoint) >= 2 {
							if usage, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
								metricMu.Lock()
								metrics.DiskUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", usage), 64)
								metricMu.Unlock()
							}
						}
					}
				}()

				// 等待所有指标查询完成
				metricWg.Wait()
			}

			// 安全地将结果写入map
			mu.Lock()
			resultMap[hostID] = metrics
			mu.Unlock()
		}(id)
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 返回结果，即使部分主机查询失败也返回成功查询的结果
	result.Success(c, resultMap)
}

// queryPrometheus 查询Prometheus指标数据
func (s *MonitorServiceImpl) queryPrometheus(query string, start, end time.Time, step string) (*monitormodel.PrometheusQueryResult, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	// 检查Prometheus URL配置
	if s.prometheusURL == "" {
		return nil, fmt.Errorf("prometheus URL not configured")
	}

	var req *http.Request
	var q url.Values
	var err error

	// 对于时间范围查询使用query_range端点
	if !start.IsZero() && !end.IsZero() {
		req, err = http.NewRequest("GET", strings.TrimSuffix(s.prometheusURL, "/")+"/api/v1/query_range", nil)
		if err != nil {
			return nil, err
		}

		q = req.URL.Query()
		q.Add("query", query)
		q.Add("start", strconv.FormatInt(start.Unix(), 10))
		q.Add("end", strconv.FormatInt(end.Unix(), 10))
		q.Add("step", step)
	} else {
		// 即时查询使用query端点
		req, err = http.NewRequest("GET", s.prometheusURL+"/api/v1/query", nil)
		if err != nil {
			return nil, err
		}
		q = req.URL.Query()
		q.Set("query", query)
	}

	// 强制不使用缓存并获取最新数据
	q.Add("nocache", "1")
	q.Add("time", fmt.Sprintf("%d", time.Now().Unix()))  // 添加时间戳防止缓存
	q.Add("_", fmt.Sprintf("%d", time.Now().UnixNano())) // 添加纳秒级时间戳确保唯一性
	req.URL.RawQuery = q.Encode()

	// 设置请求头确保不使用缓存
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Expires", "0")
	req.Header.Set("X-Request-ID", fmt.Sprintf("%d", time.Now().UnixNano())) // 添加唯一请求ID

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result monitormodel.PrometheusQueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Prometheus response: %v\nRaw response: %s", err, string(body))
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("prometheus query failed: %s", string(body))
	}

	// 兼容vector和matrix类型结果
	if result.Data.ResultType != "matrix" && result.Data.ResultType != "vector" {
		return nil, fmt.Errorf("unexpected result type: %s", result.Data.ResultType)
	}

	return &result, nil
}

// GetTopProcesses 获取主机TOP进程使用率(前5名)
func (s *MonitorServiceImpl) GetTopProcesses(c *gin.Context, hostID uint) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 优先从缓存获取主机信息
	host, err := s.getHostWithCache(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	// 使用内网IP作为instance标识
	instanceIdentifier := s.getInstanceIdentifier(host)
	if instanceIdentifier == "" {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机内网IP和名称均为空")
		return
	}

	// 检查主机是否在线
	if !s.checkHostOnline(*host) {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机当前不在线")
		return
	}

	// ... existing code ...
	var wg sync.WaitGroup
	var topCPU, topMemory []monitormodel.ProcessInfo
	var mu sync.Mutex

	// 查询CPU使用率TOP5进程
	wg.Add(1)
	go func() {
		defer wg.Done()
		cpuQuery := fmt.Sprintf("topk(5, process_cpu_percent{host=\"%s\"})", instanceIdentifier)
		cpuResult, err := s.queryPrometheus(cpuQuery, time.Time{}, time.Time{}, "")
		if err == nil && len(cpuResult.Data.Result) > 0 {
			mu.Lock()
			for _, item := range cpuResult.Data.Result {
				if len(item.Value) >= 2 {
					// 解析指标数据
					var processInfo monitormodel.ProcessInfo
					processInfo.Host = item.Metric["host"]
					processInfo.Name = item.Metric["name"]

					// 解析PID
					if pidStr, ok := item.Metric["pid"]; ok {
						if pid, err := strconv.ParseUint(pidStr, 10, 32); err == nil {
							processInfo.PID = uint(pid)
						}
					}

					// 解析CPU使用率
					if valStr, ok := item.Value[1].(string); ok {
						if cpuVal, err := strconv.ParseFloat(valStr, 64); err == nil {
							processInfo.CPUPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", cpuVal), 64)
						}
					}

					topCPU = append(topCPU, processInfo)
				}
			}
			mu.Unlock()
		}
	}()

	// 查询内存使用率TOP5进程
	wg.Add(1)
	go func() {
		defer wg.Done()
		memQuery := fmt.Sprintf("topk(5, process_memory_percent{host=\"%s\"})", instanceIdentifier)
		memResult, err := s.queryPrometheus(memQuery, time.Time{}, time.Time{}, "")
		if err == nil && len(memResult.Data.Result) > 0 {
			mu.Lock()
			for _, item := range memResult.Data.Result {
				if len(item.Value) >= 2 {
					// 解析指标数据
					var processInfo monitormodel.ProcessInfo
					processInfo.Host = item.Metric["host"]
					processInfo.Name = item.Metric["name"]

					// 解析PID
					if pidStr, ok := item.Metric["pid"]; ok {
						if pid, err := strconv.ParseUint(pidStr, 10, 32); err == nil {
							processInfo.PID = uint(pid)
						}
					}

					// 解析内存使用率
					if valStr, ok := item.Value[1].(string); ok {
						if memVal, err := strconv.ParseFloat(valStr, 64); err == nil {
							processInfo.MemPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", memVal), 64)
						}
					}

					topMemory = append(topMemory, processInfo)
				}
			}
			mu.Unlock()
		}
	}()

	// 等待所有查询完成
	wg.Wait()

	// 构建返回结果
	topResult := &monitormodel.TopProcessesResult{
		HostID:     hostID,
		HostName:   instanceIdentifier,
		TopCPU:     topCPU,
		TopMemory:  topMemory,
		UpdateTime: time.Now().Unix(),
	}

	result.Success(c, topResult)
}

// GetHostAllMetricsHistory 获取主机所有指标历史数据
func (s *MonitorServiceImpl) GetHostAllMetricsHistory(c *gin.Context, hostID uint, startTime, endTime string, step string) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 获取主机信息
	host, err := s.hostDao.GetCmdbHostById(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	// 使用内网IP作为instance标识
	instanceIdentifier := s.getInstanceIdentifier(&host)
	if instanceIdentifier == "" {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机内网IP和名称均为空")
		return
	}

	// 解析时间参数
	loc, _ := time.LoadLocation("Asia/Shanghai")
	start, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, loc)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "开始时间格式错误，请使用: 2006-01-02 15:04:05")
		return
	}
	end, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, loc)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "结束时间格式错误，请使用: 2006-01-02 15:04:05")
		return
	}

	// 自动计算step参数
	if step == "" {
		step = calculateStep(end.Sub(start))
	} else {
		step = strings.Trim(step, `"`) // 移除可能存在的额外引号
	}

	// 验证step格式
	if _, err := time.ParseDuration(step); err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError),
			"无效的步长格式，请使用如: 15s, 1m, 5m")
		return
	}

	// 使用WaitGroup并行查询所有指标
	var wg sync.WaitGroup
	var mu sync.Mutex
	allMetrics := &monitormodel.AllMetricsHistory{
		HostID: hostID,
	}

	// 查询CPU指标
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("rate(system_cpu_usage_percent{instance=\"%s\"}[1m])", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.CPU = append(allMetrics.CPU, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询内存指标
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_memory_usage_percent{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.Memory = append(allMetrics.Memory, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询磁盘指标
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_disk_usage_percent{instance=\"%s\",mountpoint=\"/\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.Disk = append(allMetrics.Disk, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询网络接收指标
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_network_receive_kb_per_second{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.NetworkReceive = append(allMetrics.NetworkReceive, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询网络发送指标
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_network_send_kb_per_second{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.NetworkSend = append(allMetrics.NetworkSend, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询1分钟系统负载
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_load_average_1min{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.Load1min = append(allMetrics.Load1min, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询5分钟系统负载
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_load_average_5min{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.Load5min = append(allMetrics.Load5min, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询15分钟系统负载
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_load_average_15min{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.Load15min = append(allMetrics.Load15min, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询磁盘读取速率
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_disk_read_kb_per_second{instance=\"%s\",device=\"vda\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.DiskReadKB = append(allMetrics.DiskReadKB, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询磁盘写入速率
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_disk_write_kb_per_second{instance=\"%s\",device=\"vda\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.DiskWriteKB = append(allMetrics.DiskWriteKB, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询系统进程总数
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_total_processes{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.TotalProcesses = append(allMetrics.TotalProcesses, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询CPU用户态使用率
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_cpu_user_percent{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.CPUUser = append(allMetrics.CPUUser, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询CPU内核态使用率
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_cpu_system_percent{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.CPUSystem = append(allMetrics.CPUSystem, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询已用内存
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_memory_used_gb{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.MemoryUsed = append(allMetrics.MemoryUsed, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询可用内存
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_memory_free_gb{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.MemoryFree = append(allMetrics.MemoryFree, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询磁盘已用空间
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_disk_used_bytes{instance=\"%s\",mountpoint=\"/\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								// 注意: system_disk_used_bytes 实际存储的就是GB值，不需要转换
								allMetrics.DiskUsed = append(allMetrics.DiskUsed, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询磁盘总空间
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_disk_total_bytes{instance=\"%s\",mountpoint=\"/\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								// 注意: system_disk_total_bytes 实际存储的就是GB值，不需要转换
								allMetrics.DiskTotal = append(allMetrics.DiskTotal, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询TCP连接数
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_tcp_connections{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.TCPConnections = append(allMetrics.TCPConnections, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询UDP连接数
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_udp_connections{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.UDPConnections = append(allMetrics.UDPConnections, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 查询系统打开文件数量
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := fmt.Sprintf("system_open_file_descriptors{instance=\"%s\"}", instanceIdentifier)
		promResult, err := s.queryPrometheus(query, start, end, step)
		if err == nil && len(promResult.Data.Result) > 0 {
			mu.Lock()
			for _, resultItem := range promResult.Data.Result {
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								allMetrics.OpenFiles = append(allMetrics.OpenFiles, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// 等待所有查询完成
	wg.Wait()

	result.Success(c, allMetrics)
}

// calculateStep 根据时间范围计算步长，固定返回60个数据点
func calculateStep(duration time.Duration) string {
	step := duration / 60
	return fmt.Sprintf("%.0fs", step.Seconds())
}

// GetHostMetricHistory 获取主机指标历史数据
func (s *MonitorServiceImpl) GetHostMetricHistory(c *gin.Context, hostID uint, metric string, startTime, endTime string, step string) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 获取主机信息
	host, err := s.hostDao.GetCmdbHostById(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	// 使用内网IP作为instance标识
	instanceIdentifier := s.getInstanceIdentifier(&host)
	if instanceIdentifier == "" {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机内网IP和名称均为空")
		return
	}

	// 验证指标名称
	switch metric {
	case "cpu", "memory", "disk":
		// 有效指标
	default:
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "不支持的指标类型")
		return
	}

	// 解析时间参数
	loc, _ := time.LoadLocation("Asia/Shanghai")
	start, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, loc)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "开始时间格式错误，请使用: 2006-01-02 15:04:05")
		return
	}
	end, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, loc)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "结束时间格式错误，请使用: 2006-01-02 15:04:05")
		return
	}

	// 自动计算step参数
	if step == "" {
		step = calculateStep(end.Sub(start))
	} else {
		step = strings.Trim(step, `"`) // 移除可能存在的额外引号
	}

	// 验证step格式
	if _, err := time.ParseDuration(step); err != nil {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError),
			"无效的步长格式，请使用如: 15s, 1m, 5m")
		return
	}

	// 构建Prometheus查询
	var query string
	switch metric {
	case "cpu":
		query = fmt.Sprintf("rate(system_cpu_usage_percent{instance=\"%s\"}[1m])", instanceIdentifier)
	case "memory":
		query = fmt.Sprintf("system_memory_usage_percent{instance=\"%s\"}", instanceIdentifier)
	case "disk":
		query = fmt.Sprintf("system_disk_usage_percent{instance=\"%s\",mountpoint=\"/\"}", instanceIdentifier)
	}

	// 执行Prometheus查询
	promResult, err := s.queryPrometheus(query, start, end, step)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED),
			fmt.Sprintf("Prometheus查询失败: %v", err))
		return
	}

	history := &monitormodel.HostMetricHistory{
		HostID: hostID,
		Metric: metric,
	}

	// 提取时间序列数据
	if len(promResult.Data.Result) > 0 {
		for _, resultItem := range promResult.Data.Result {
			switch promResult.Data.ResultType {
			case "matrix":
				// 处理matrix类型结果(时间范围查询)
				for _, value := range resultItem.Values {
					if len(value) >= 2 {
						if timestamp, ok := value[0].(float64); ok {
							if valStr, ok := value[1].(string); ok {
								val, _ := strconv.ParseFloat(valStr, 64)
								history.TimeData = append(history.TimeData, monitormodel.MetricDataPoint{
									Timestamp: int64(timestamp),
									Value:     val,
								})
							}
						}
					}
				}
			case "vector":
				// 处理vector类型结果(即时查询)
				if len(resultItem.Value) >= 2 {
					if timestamp, ok := resultItem.Value[0].(float64); ok {
						if valStr, ok := resultItem.Value[1].(string); ok {
							val, _ := strconv.ParseFloat(valStr, 64)
							history.TimeData = append(history.TimeData, monitormodel.MetricDataPoint{
								Timestamp: int64(timestamp),
								Value:     val,
							})
						}
					}
				}
			default:
				// 未知结果类型，不做处理
			}
		}
	}

	result.Success(c, history)
}

// GetHostPorts 获取主机端口信息
func (s *MonitorServiceImpl) GetHostPorts(c *gin.Context, hostID uint) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 首先检查Redis缓存
	cacheKey := fmt.Sprintf("monitor:ports:%d", hostID)
	if cachedData, err := redis.RedisDb.Get(context.Background(), cacheKey).Result(); err == nil {
		var cachedResult monitormodel.HostPortsResult
		if err := json.Unmarshal([]byte(cachedData), &cachedResult); err == nil {
			// 检查缓存是否过期（5分钟内的数据认为有效）
			if time.Now().Unix()-cachedResult.UpdateTime < 300 {
				result.Success(c, &cachedResult)
				return
			}
		}
	}

	// 优先从缓存获取主机信息
	host, err := s.getHostWithCache(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	// 使用内网IP作为instance标识
	instanceIdentifier := s.getInstanceIdentifier(host)
	if instanceIdentifier == "" {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机内网IP和名称均为空")
		return
	}

	// 检查主机是否在线
	if !s.checkHostOnline(*host) {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机当前不在线")
		return
	}

	// 查询TCP端口监听状态
	query := fmt.Sprintf("tcp_port_listening{instance=\"%s\"}", instanceIdentifier)
	now := time.Now()
	startTime := now.Add(-5 * time.Minute) // 查询最近5分钟的数据
	endTime := now
	step := "30s"

	promResult, err := s.queryPrometheus(query, startTime, endTime, step)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED),
			fmt.Sprintf("Prometheus查询失败: %v", err))
		return
	}

	// 批量查询所有进程的CPU和内存使用率（一次查询获取所有数据）
	cpuQuery := fmt.Sprintf("process_cpu_usage_percent{instance=\"%s\"}", instanceIdentifier)
	memQuery := fmt.Sprintf("process_memory_usage_percent{instance=\"%s\"}", instanceIdentifier)

	// 并行查询CPU和内存数据
	var cpuResult, memResult *monitormodel.PrometheusQueryResult
	var cpuErr, memErr error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		cpuResult, cpuErr = s.queryPrometheus(cpuQuery, startTime, endTime, step)
	}()
	go func() {
		defer wg.Done()
		memResult, memErr = s.queryPrometheus(memQuery, startTime, endTime, step)
	}()
	wg.Wait()

	// 构建PID到CPU/内存使用率的映射
	cpuMap := make(map[string]float64)
	memMap := make(map[string]float64)

	if cpuErr == nil && len(cpuResult.Data.Result) > 0 {
		for _, resultItem := range cpuResult.Data.Result {
			pid := resultItem.Metric["pid"]
			var cpuUsage float64
			if len(resultItem.Values) > 0 {
				lastPoint := resultItem.Values[len(resultItem.Values)-1]
				if len(lastPoint) >= 2 {
					if valStr, ok := lastPoint[1].(string); ok {
						cpuUsage, _ = strconv.ParseFloat(valStr, 64)
					}
				}
			} else if len(resultItem.Value) >= 2 {
				if valStr, ok := resultItem.Value[1].(string); ok {
					cpuUsage, _ = strconv.ParseFloat(valStr, 64)
				}
			}
			cpuMap[pid] = cpuUsage
		}
	}

	if memErr == nil && len(memResult.Data.Result) > 0 {
		for _, resultItem := range memResult.Data.Result {
			pid := resultItem.Metric["pid"]
			var memUsage float64
			if len(resultItem.Values) > 0 {
				lastPoint := resultItem.Values[len(resultItem.Values)-1]
				if len(lastPoint) >= 2 {
					if valStr, ok := lastPoint[1].(string); ok {
						memUsage, _ = strconv.ParseFloat(valStr, 64)
					}
				}
			} else if len(resultItem.Value) >= 2 {
				if valStr, ok := resultItem.Value[1].(string); ok {
					memUsage, _ = strconv.ParseFloat(valStr, 64)
				}
			}
			memMap[pid] = memUsage
		}
	}

	var ports []monitormodel.PortInfo

	// 解析端口数据，并从映射中获取CPU/内存使用率
	if len(promResult.Data.Result) > 0 {
		for _, resultItem := range promResult.Data.Result {
			port := resultItem.Metric["port"]
			pid := resultItem.Metric["pid"]
			service := resultItem.Metric["service"]

			// 获取最新的状态值
			var status int
			if len(resultItem.Values) > 0 {
				lastPoint := resultItem.Values[len(resultItem.Values)-1]
				if len(lastPoint) >= 2 {
					if valStr, ok := lastPoint[1].(string); ok {
						if val, err := strconv.ParseFloat(valStr, 64); err == nil {
							status = int(val)
						}
					}
				}
			} else if len(resultItem.Value) >= 2 {
				if valStr, ok := resultItem.Value[1].(string); ok {
					if val, err := strconv.ParseFloat(valStr, 64); err == nil {
						status = int(val)
					}
				}
			}

			// 从映射中获取CPU和内存使用率
			cpuUsage := cpuMap[pid]
			memUsage := memMap[pid]

			ports = append(ports, monitormodel.PortInfo{
				Port:     port,
				PID:      pid,
				Service:  service,
				Status:   status,
				CPUUsage: cpuUsage,
				MemUsage: memUsage,
			})
		}
	}

	// 构造结果
	portResult := &monitormodel.HostPortsResult{
		HostID:     hostID,
		HostName:   instanceIdentifier,
		Ports:      ports,
		Total:      len(ports),
		UpdateTime: time.Now().Unix(),
	}

	// 将结果存储到Redis缓存（5分钟过期）
	if cacheData, err := json.Marshal(portResult); err == nil {
		redis.RedisDb.Set(context.Background(), cacheKey, cacheData, 5*time.Minute)
	}

	result.Success(c, portResult)
}

// GetTopUsageHosts 获取使用率TOP主机
func (s *MonitorServiceImpl) GetTopUsageHosts(c *gin.Context, days, limit int) {
	// 参数验证
	if days <= 0 || days > 15 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "天数参数必须在1-15之间")
		return
	}
	if limit <= 0 || limit > 20 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "limit参数必须在1-20之间")
		return
	}

	// 计算时间范围
	loc, _ := time.LoadLocation("Asia/Shanghai")
	endTime := time.Now().In(loc)
	startTime := endTime.Add(-time.Duration(days) * 24 * time.Hour)

	// 使用较大的步长减少数据点,但确保有足够的数据来计算平均值
	step := "1h"

	// 并发查询三个指标
	var wg sync.WaitGroup
	var mu sync.Mutex

	var topCPU []monitormodel.HostUsageRank
	var topMemory []monitormodel.HostUsageRank
	var topDisk []monitormodel.HostUsageRank

	// 查询 CPU 使用率 TOP N - 使用范围查询获取时间序列数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 先查询指定时间范围内的所有主机CPU数据
		cpuQuery := "system_cpu_usage_percent"
		cpuResult, err := s.queryPrometheus(cpuQuery, startTime, endTime, step)
		if err != nil {
			fmt.Printf("CPU查询失败: %v, query: %s\n", err, cpuQuery)
			return
		}
		if cpuResult != nil && len(cpuResult.Data.Result) > 0 {
			ranks := s.parseTopUsageResult(cpuResult, "cpu")
			mu.Lock()
			topCPU = ranks
			mu.Unlock()
		}
	}()

	// 查询内存使用率 TOP N
	wg.Add(1)
	go func() {
		defer wg.Done()
		memQuery := "system_memory_usage_percent"
		memResult, err := s.queryPrometheus(memQuery, startTime, endTime, step)
		if err != nil {
			fmt.Printf("内存查询失败: %v, query: %s\n", err, memQuery)
			return
		}
		if memResult != nil && len(memResult.Data.Result) > 0 {
			ranks := s.parseTopUsageResult(memResult, "memory")
			mu.Lock()
			topMemory = ranks
			mu.Unlock()
		}
	}()

	// 查询磁盘使用率 TOP N
	wg.Add(1)
	go func() {
		defer wg.Done()
		diskQuery := "system_disk_usage_percent{mountpoint=\"/\"}"
		diskResult, err := s.queryPrometheus(diskQuery, startTime, endTime, step)
		if err != nil {
			fmt.Printf("磁盘查询失败: %v, query: %s\n", err, diskQuery)
			return
		}
		if diskResult != nil && len(diskResult.Data.Result) > 0 {
			ranks := s.parseTopUsageResult(diskResult, "disk")
			mu.Lock()
			topDisk = ranks
			mu.Unlock()
		}
	}()

	// 等待所有查询完成
	wg.Wait()

	// 对结果进行排序并取TOP N
	topCPU = sortAndLimitRanks(topCPU, limit)
	topMemory = sortAndLimitRanks(topMemory, limit)
	topDisk = sortAndLimitRanks(topDisk, limit)

	// 构造返回结果
	topUsageResult := &monitormodel.TopUsageHostsResult{
		TopCPU:    topCPU,
		TopMemory: topMemory,
		TopDisk:   topDisk,
		Period:    fmt.Sprintf("%d天", days),
		StartTime: startTime.Format("2006-01-02 15:04:05"),
		EndTime:   endTime.Format("2006-01-02 15:04:05"),
	}

	result.Success(c, topUsageResult)
}

// parseTopUsageResult 解析 Prometheus TOP 查询结果
func (s *MonitorServiceImpl) parseTopUsageResult(queryResult *monitormodel.PrometheusQueryResult, metricType string) []monitormodel.HostUsageRank {
	// 使用map去重,key为主机名
	hostMap := make(map[string]*monitormodel.HostUsageRank)

	for _, item := range queryResult.Data.Result {
		// 获取主机名
		hostName, ok := item.Metric["instance"]
		if !ok {
			continue
		}

		// 计算平均值、最大值和当前值
		var avgUsage, maxUsage, currentUsage float64
		var valueCount int
		var valueSum float64

		if len(item.Values) > 0 {
			for _, point := range item.Values {
				if len(point) >= 2 {
					if val, err := strconv.ParseFloat(point[1].(string), 64); err == nil {
						valueSum += val
						valueCount++
						if val > maxUsage {
							maxUsage = val
						}
					}
				}
			}
			// 最后一个值作为当前值
			if lastPoint := item.Values[len(item.Values)-1]; len(lastPoint) >= 2 {
				if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
					currentUsage = val
				}
			}
		} else if len(item.Value) >= 2 {
			// vector 类型结果
			if val, err := strconv.ParseFloat(item.Value[1].(string), 64); err == nil {
				avgUsage = val
				maxUsage = val
				currentUsage = val
				valueCount = 1
			}
		}

		if valueCount > 0 {
			avgUsage = valueSum / float64(valueCount)
		}

		// 检查该主机是否已经存在
		if existingRank, exists := hostMap[hostName]; exists {
			// 如果存在,合并数据(取更大的值)
			if avgUsage > existingRank.AvgUsage {
				existingRank.AvgUsage = parseFloat(avgUsage)
			}
			if maxUsage > existingRank.MaxUsage {
				existingRank.MaxUsage = parseFloat(maxUsage)
			}
			if currentUsage > existingRank.CurrentUsage {
				existingRank.CurrentUsage = parseFloat(currentUsage)
			}
			continue
		}

		// 尝试从数据库获取主机信息
		host, err := s.hostDao.GetCmdbHostByName(hostName)
		var hostID uint
		var ip string

		if err != nil {
			// 如果数据库中不存在,使用默认值(主机ID为0,IP使用主机名)
			hostID = 0
			ip = hostName // 使用主机名作为IP显示
		} else {
			hostID = host.ID
			ip = host.PrivateIP
			if ip == "" {
				ip = host.SSHIP
			}
			if ip == "" {
				ip = hostName
			}
		}

		rank := monitormodel.HostUsageRank{
			HostID:       hostID,
			HostName:     hostName,
			IP:           ip,
			AvgUsage:     parseFloat(avgUsage),
			MaxUsage:     parseFloat(maxUsage),
			CurrentUsage: parseFloat(currentUsage),
		}

		hostMap[hostName] = &rank
	}

	// 将map转换为数组
	var ranks []monitormodel.HostUsageRank
	for _, rank := range hostMap {
		ranks = append(ranks, *rank)
	}

	return ranks
}

// sortAndLimitRanks 对主机使用率排名进行排序并限制数量
func sortAndLimitRanks(ranks []monitormodel.HostUsageRank, limit int) []monitormodel.HostUsageRank {
	if len(ranks) == 0 {
		return ranks
	}

	// 按平均使用率从高到低排序
	for i := 0; i < len(ranks)-1; i++ {
		for j := i + 1; j < len(ranks); j++ {
			if ranks[j].AvgUsage > ranks[i].AvgUsage {
				ranks[i], ranks[j] = ranks[j], ranks[i]
			}
		}
	}

	// 限制返回数量
	if len(ranks) > limit {
		return ranks[:limit]
	}

	return ranks
}

// parseFloat 格式化浮点数保留2位小数
func parseFloat(value float64) float64 {
	result, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return result
}

// GetCompleteHostMetrics 获取主机完整监控数据
func (s *MonitorServiceImpl) GetCompleteHostMetrics(c *gin.Context, hostID uint) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 优先从缓存获取主机信息
	host, err := s.getHostWithCache(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	// 使用内网IP作为instance标识
	instanceIdentifier := s.getInstanceIdentifier(host)
	if instanceIdentifier == "" {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机内网IP和名称均为空")
		return
	}

	// 检查主机是否在线
	if !s.checkHostOnline(*host) {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机离线")
		return
	}

	// 查询最近5分钟的数据
	now := time.Now()
	startTime := now.Add(-5 * time.Minute)
	endTime := now
	step := "30s"

	metrics := &monitormodel.CompleteHostMetrics{
		HostID:   hostID,
		HostName: host.Name,
		IP:       host.PrivateIP,
		UpdateTime: time.Now().Unix(),
	}

	// 查询CPU相关指标
	metrics.CPU = s.queryCPUMetrics(instanceIdentifier, startTime, endTime, step)

	// 查询内存相关指标
	metrics.Memory = s.queryMemoryMetrics(instanceIdentifier, startTime, endTime, step)

	// 查询磁盘相关指标
	metrics.Disk = s.queryDiskMetrics(instanceIdentifier, startTime, endTime, step)

	// 查询负载指标
	metrics.Load = s.queryLoadMetrics(instanceIdentifier, startTime, endTime, step)

	// 查询连接数指标
	metrics.Connections = s.queryConnectionMetrics(instanceIdentifier, startTime, endTime, step)

	// 查询网络指标
	metrics.Network = s.queryNetworkMetrics(instanceIdentifier, startTime, endTime, step)

	// 查询系统运行时间
	metrics.Uptime, metrics.UptimeStr = s.queryUptimeMetrics(instanceIdentifier, startTime, endTime, step)

	// 查询系统进程总数
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_total_processes{instance=\"%s\"}", instanceIdentifier),
		startTime, endTime, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				metrics.ProcessCount = int(val)
			}
		}
	}

	// 查询文件打开数量
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_open_file_descriptors{instance=\"%s\"}", instanceIdentifier),
		startTime, endTime, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				metrics.OpenFiles = int(val)
			}
		}
	}

	// 查询服务列表
	metrics.Services = s.queryServiceMetrics(instanceIdentifier, startTime, endTime, step)

	result.Success(c, metrics)
}

// queryCPUMetrics 查询CPU指标
func (s *MonitorServiceImpl) queryCPUMetrics(instance string, start, end time.Time, step string) monitormodel.CPUMetrics {
	cpuMetrics := monitormodel.CPUMetrics{}

	// CPU总使用率
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_cpu_usage_percent{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				cpuMetrics.Total = parseFloat(val)
			}
		}
	}

	// CPU用户态使用率
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_cpu_user_percent{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				cpuMetrics.User = parseFloat(val)
			}
		}
	}

	// CPU内核态使用率
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_cpu_system_percent{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				cpuMetrics.System = parseFloat(val)
			}
		}
	}

	return cpuMetrics
}

// queryMemoryMetrics 查询内存指标
func (s *MonitorServiceImpl) queryMemoryMetrics(instance string, start, end time.Time, step string) monitormodel.MemoryMetrics {
	memMetrics := monitormodel.MemoryMetrics{}

	// 内存使用率
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_memory_usage_percent{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				memMetrics.UsageRate = parseFloat(val)
			}
		}
	}

	// 已使用内存
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_memory_used_gb{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				// 转换为字节
				memMetrics.Used = uint64(val * 1024 * 1024 * 1024)
				memMetrics.UsedStr = formatBytes(memMetrics.Used)
			}
		}
	}

	// 总内存
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_memory_total_gb{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				// 转换为字节
				memMetrics.Total = uint64(val * 1024 * 1024 * 1024)
				memMetrics.TotalStr = formatBytes(memMetrics.Total)
			}
		}
	}

	// 可用内存
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_memory_free_gb{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				// 转换为字节
				memMetrics.Free = uint64(val * 1024 * 1024 * 1024)
				memMetrics.FreeStr = formatBytes(memMetrics.Free)
			}
		}
	}

	// 交换内存使用量
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_swap_used_gb{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				// 转换为字节
				memMetrics.SwapUsed = uint64(val * 1024 * 1024 * 1024)
				memMetrics.SwapUsedStr = formatBytes(memMetrics.SwapUsed)
			}
		}
	}

	return memMetrics
}

// queryDiskMetrics 查询磁盘指标
func (s *MonitorServiceImpl) queryDiskMetrics(instance string, start, end time.Time, step string) monitormodel.DiskMetrics {
	diskMetrics := monitormodel.DiskMetrics{}

	// 磁盘使用率
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_disk_usage_percent{instance=\"%s\",mountpoint=\"/\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				diskMetrics.UsageRate = parseFloat(val)
			}
		}
	}

	// 磁盘已用空间（不指定device，自动匹配）
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_disk_used_bytes{instance=\"%s\",mountpoint=\"/\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 {
		// 如果有多个设备，取第一个（通常就是根分区所在设备）
		if len(data.Data.Result[0].Values) > 0 {
			lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
			if len(lastPoint) >= 2 {
				if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
					diskMetrics.Used = uint64(val)
					diskMetrics.UsedStr = formatBytes(uint64(val))
				}
			}
		}
	}

	// 磁盘总空间（不指定device，自动匹配）
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_disk_total_bytes{instance=\"%s\",mountpoint=\"/\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 {
		// 如果有多个设备，取第一个（通常就是根分区所在设备）
		if len(data.Data.Result[0].Values) > 0 {
			lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
			if len(lastPoint) >= 2 {
				if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
					diskMetrics.Total = uint64(val)
					diskMetrics.TotalStr = formatBytes(uint64(val))
				}
			}
		}
	}

	// 磁盘读取速率（聚合所有设备）
	if data, err := s.queryPrometheus(
		fmt.Sprintf("sum(system_disk_read_kb_per_second{instance=\"%s\"})", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				diskMetrics.DiskThroughput.ReadKBPerSec = parseFloat(val)
			}
		}
	}

	// 磁盘写入速率（聚合所有设备）
	if data, err := s.queryPrometheus(
		fmt.Sprintf("sum(system_disk_write_kb_per_second{instance=\"%s\"})", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				diskMetrics.DiskThroughput.WriteKBPerSec = parseFloat(val)
			}
		}
	}

	// 磁盘读IOPS（聚合所有设备）
	if data, err := s.queryPrometheus(
		fmt.Sprintf("sum(system_disk_read_iops{instance=\"%s\"})", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				diskMetrics.DiskIOPS.ReadPerSec = parseFloat(val)
			}
		}
	}

	// 磁盘写IOPS（聚合所有设备）
	if data, err := s.queryPrometheus(
		fmt.Sprintf("sum(system_disk_write_iops{instance=\"%s\"})", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				diskMetrics.DiskIOPS.WritePerSec = parseFloat(val)
			}
		}
	}

	// Inode使用率
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_disk_inode_usage_percent{instance=\"%s\",mountpoint=\"/\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				diskMetrics.InodeUsage = parseFloat(val)
			}
		}
	}

	// Inode总数
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_disk_inode_total{instance=\"%s\",mountpoint=\"/\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				diskMetrics.InodeTotal = uint64(val)
			}
		}
	}

	// Inode已用数
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_disk_inode_used{instance=\"%s\",mountpoint=\"/\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				diskMetrics.InodeUsed = uint64(val)
			}
		}
	}

	// Inode可用数
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_disk_inode_free{instance=\"%s\",mountpoint=\"/\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				diskMetrics.InodeFree = uint64(val)
			}
		}
	}

	// 计算剩余空间
	if diskMetrics.Total > 0 && diskMetrics.Used > 0 {
		diskMetrics.Free = diskMetrics.Total - diskMetrics.Used
		diskMetrics.FreeStr = formatBytes(diskMetrics.Free)
	}

	return diskMetrics
}

// queryLoadMetrics 查询负载指标
func (s *MonitorServiceImpl) queryLoadMetrics(instance string, start, end time.Time, step string) monitormodel.LoadMetrics {
	loadMetrics := monitormodel.LoadMetrics{}

	// 1分钟负载
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_load_average_1min{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				loadMetrics.Load1 = parseFloat(val)
			}
		}
	}

	// 5分钟负载
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_load_average_5min{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				loadMetrics.Load5 = parseFloat(val)
			}
		}
	}

	// 15分钟负载
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_load_average_15min{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				loadMetrics.Load15 = parseFloat(val)
			}
		}
	}

	// 1分钟负载比
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_load_ratio_1min{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				loadMetrics.LoadRatio1 = parseFloat(val)
			}
		}
	}

	// 5分钟负载比
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_load_ratio_5min{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				loadMetrics.LoadRatio5 = parseFloat(val)
			}
		}
	}

	// 15分钟负载比
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_load_ratio_15min{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				loadMetrics.LoadRatio15 = parseFloat(val)
			}
		}
	}

	// CPU核心数
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_cpu_core_count{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				loadMetrics.CPUCores = int(val)
			}
		}
	}

	return loadMetrics
}

// queryConnectionMetrics 查询连接数指标
func (s *MonitorServiceImpl) queryConnectionMetrics(instance string, start, end time.Time, step string) monitormodel.ConnMetrics {
	connMetrics := monitormodel.ConnMetrics{}

	// TCP连接数
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_tcp_connections{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				connMetrics.TCP = int(val)
			}
		}
	}

	// UDP连接数
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_udp_connections{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				connMetrics.UDP = int(val)
			}
		}
	}

	connMetrics.Total = connMetrics.TCP + connMetrics.UDP

	return connMetrics
}

// queryNetworkMetrics 查询网络指标
func (s *MonitorServiceImpl) queryNetworkMetrics(instance string, start, end time.Time, step string) monitormodel.NetworkMetrics {
	netMetrics := monitormodel.NetworkMetrics{}

	// 网络接收速率
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_network_receive_kb_per_second{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				netMetrics.ReceiveKBPerSec = parseFloat(val)
			}
		}
	}

	// 网络发送速率
	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_network_send_kb_per_second{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				netMetrics.SendKBPerSec = parseFloat(val)
			}
		}
	}

	return netMetrics
}

// queryUptimeMetrics 查询系统运行时间
func (s *MonitorServiceImpl) queryUptimeMetrics(instance string, start, end time.Time, step string) (float64, string) {
	var uptimeDays float64

	if data, err := s.queryPrometheus(
		fmt.Sprintf("system_uptime_days{instance=\"%s\"}", instance),
		start, end, step,
	); err == nil && len(data.Data.Result) > 0 && len(data.Data.Result[0].Values) > 0 {
		lastPoint := data.Data.Result[0].Values[len(data.Data.Result[0].Values)-1]
		if len(lastPoint) >= 2 {
			if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
				uptimeDays = val
			}
		}
	}

	return parseFloat(uptimeDays), formatUptimeFromDays(uptimeDays)
}

// queryServiceMetrics 查询服务列表
func (s *MonitorServiceImpl) queryServiceMetrics(instance string, start, end time.Time, step string) []monitormodel.ServiceInfo {
	services := []monitormodel.ServiceInfo{}

	// 查询所有服务的CPU使用率
	cpuData, err := s.queryPrometheus(
		fmt.Sprintf("service_cpu_usage_percent{instance=\"%s\"}", instance),
		start, end, step,
	)
	if err != nil {
		return services
	}

	// 遍历每个服务
	for _, result := range cpuData.Data.Result {
		serviceName := result.Metric["service_name"]
		if serviceName == "" {
			continue
		}

		var cpuUsage, memUsage float64

		// 获取CPU使用率
		if len(result.Values) > 0 {
			lastPoint := result.Values[len(result.Values)-1]
			if len(lastPoint) >= 2 {
				if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
					cpuUsage = parseFloat(val)
				}
			}
		}

		// 查询该服务的内存使用率
		memData, err := s.queryPrometheus(
			fmt.Sprintf("service_memory_usage_percent{instance=\"%s\",service_name=\"%s\"}", instance, serviceName),
			start, end, step,
		)
		if err == nil && len(memData.Data.Result) > 0 && len(memData.Data.Result[0].Values) > 0 {
			lastPoint := memData.Data.Result[0].Values[len(memData.Data.Result[0].Values)-1]
			if len(lastPoint) >= 2 {
				if val, err := strconv.ParseFloat(lastPoint[1].(string), 64); err == nil {
					memUsage = parseFloat(val)
				}
			}
		}

		services = append(services, monitormodel.ServiceInfo{
			Name:   serviceName,
			CPU:    cpuUsage,
			Memory: memUsage,
		})
	}

	return services
}

// formatBytes 格式化字节大小
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatUptime 格式化系统运行时间(从秒转换)
func formatUptime(seconds float64) string {
	days := int(seconds) / 86400
	hours := (int(seconds) % 86400) / 3600
	minutes := (int(seconds) % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟", hours, minutes)
	} else {
		return fmt.Sprintf("%d分钟", minutes)
	}
}

// formatUptimeFromDays 格式化系统运行时间(从天数转换)
func formatUptimeFromDays(days float64) string {
	totalDays := int(days)
	hours := int((days - float64(totalDays)) * 24)
	minutes := int(((days - float64(totalDays)) * 24 - float64(hours)) * 60)

	if totalDays > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", totalDays, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟", hours, minutes)
	} else {
		return fmt.Sprintf("%d分钟", minutes)
	}
}
