package agent

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

// StartAgent 启动代理服务
func StartAgent() error {
	return StartAgentWithConfig(&defaultConfig)
}

// StartAgentWithConfig 使用指定配置启动代理服务
func StartAgentWithConfig(config *Config) error {
	// 设置默认值（如果配置为空）
	if config.Pushgateway.URL == "" {
		config.Pushgateway.URL = "http://8.130.14.34:9091"
	}
	if config.Heartbeat.ServerURL == "" {
		config.Heartbeat.ServerURL = "http://127.0.0.1:8000/api/v1/monitor/agent/heartbeat"
	}

	// 支持通过环境变量覆盖配置
	if url := os.Getenv("PUSHGATEWAY_URL"); url != "" {
		config.Pushgateway.URL = url
	}
	if job := os.Getenv("PUSHGATEWAY_JOB"); job != "" {
		config.Pushgateway.Job = job
	}
	if interval := os.Getenv("PUSHGATEWAY_INTERVAL"); interval != "" {
		config.Pushgateway.PushInterval = interval
	}
	if collectInterval := os.Getenv("COLLECT_INTERVAL"); collectInterval != "" {
		config.CollectInterval = collectInterval
	}
	// 心跳相关环境变量
	if heartbeatURL := os.Getenv("HEARTBEAT_SERVER_URL"); heartbeatURL != "" {
		config.Heartbeat.ServerURL = heartbeatURL
	}
	if heartbeatInterval := os.Getenv("HEARTBEAT_INTERVAL"); heartbeatInterval != "" {
		config.Heartbeat.Interval = heartbeatInterval
	}
	if heartbeatToken := os.Getenv("HEARTBEAT_TOKEN"); heartbeatToken != "" {
		config.Heartbeat.Token = heartbeatToken
	}

	// 定时采集
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in monitor goroutine: %v", r)
			}
		}()

		// 解析采集间隔
		collectInterval, err := time.ParseDuration(config.CollectInterval)
		if err != nil {
			collectInterval = 15 * time.Second // 默认值
			log.Printf("Invalid collect interval, using default 15s")
		}
		log.Printf("Collect interval set to: %v", collectInterval)

		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic in monitor function: %v", r)
					}
				}()

				monitorSystemMetrics(collectInterval.Seconds())
				monitorProcesses(config)
				monitorEndpoints(config)
				monitorPing("8.8.8.8") // 默认ping Google DNS
				monitorTCPPorts()      // TCP端口监控
				monitorServices()      // 服务监控
			}()
			time.Sleep(collectInterval)
		}
	}()

	// 启动HTTP服务
	// 添加版本查询路由
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		version := map[string]string{
			"version":    AgentVersion,
			"versionTag": AgentVersionTag,
			"buildTime":  AgentBuildTime,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(version)
	})

	// 添加查询示例路由
	http.HandleFunc("/query-examples", func(w http.ResponseWriter, r *http.Request) {
		examples := `
Prometheus查询示例：

1. 查询所有核心的CPU使用率:
system_cpu_usage_percent{}

2. 查询特定核心(core0)的CPU使用率:
system_cpu_usage_percent{core="core0"}

3. 查询多个核心的CPU使用率:
system_cpu_usage_percent{core=~"core0|core1"}

4. 查询内存使用率大于70%的指标:
system_memory_usage_percent > 70

5. 查询特定设备的磁盘读取速度:
system_disk_read_kb_per_second{device="disk0"}

6. 查询网络接收速率:
system_network_receive_kb_per_second

7. 查询进程CPU使用率(按进程名):
process_cpu_percent{name="go"}

8. 查询端点响应时间大于10ms的:
endpoint_response_time_ms > 10

8. 查询状态异常的端点:
endpoint_status == 0

9. 5分钟内的CPU使用率变化:
rate(system_cpu_usage_percent[5m])

10. 各核心CPU使用率平均值:
avg by (core) (system_cpu_usage_percent)
`
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(examples))
	})

	// 创建自定义注册表，不包含Go运行时指标
	registry := prometheus.NewRegistry()

	// 注册所有自定义指标
	registry.MustRegister(
		cpuTotalUsage,
		cpuUserUsage,
		cpuSystemUsage,
		systemLoad,
		memUsage,
		memUsedGB,
		memTotalGB,
		memFreeGB,
		swapUsedGB,
		diskUsage,
		diskUsedBytes,
		diskTotalBytes,
		diskReadRate,
		diskWriteRate,
		diskReadIOPS,
		diskWriteIOPS,
		totalProcesses,
		netRecv,
		netSent,
		processCPU,
		processMem,
		endpointResponseTime,
		endpointStatus,
		pingResponseTime,
		pingStatus,
		tcpPortStatus,
		systemUptimeDays,
		tcpConnections,
		udpConnections,
		serviceCPUUsage,
		serviceMemoryUsage,
		agentInfo, // Agent版本信息
		openFileDescriptors,
	)

	// 注册推送状态指标
	registry.MustRegister(pushSuccess, pushFailures)

	// 使用自定义注册表
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	// 启动Pushgateway推送goroutine
	if config.Pushgateway.URL != "" {
		pushInterval, err := time.ParseDuration(config.Pushgateway.PushInterval)
		if err != nil {
			pushInterval = 30 * time.Second
		}

		instance := strings.ReplaceAll(config.Pushgateway.Instance, "$HOSTNAME", hostname)

		go func() {
			for {
				pusher := push.New(config.Pushgateway.URL, config.Pushgateway.Job)
				pusher = pusher.Grouping("instance", instance)
				err := pusher.Gatherer(registry).Push()

				if err != nil {
					log.Printf("Push to Pushgateway failed: %v", err)
					pushFailures.Inc()
				} else {
					pushSuccess.Inc()
				}

				time.Sleep(pushInterval)
			}
		}()
	}

	// 启动心跳goroutine
	if config.Heartbeat.ServerURL != "" {
		heartbeatInterval, err := time.ParseDuration(config.Heartbeat.Interval)
		if err != nil {
			heartbeatInterval = 30 * time.Second // 默认30秒
			log.Printf("Invalid heartbeat interval, using default 30s")
		}
		log.Printf("Heartbeat interval set to: %v, server: %s", heartbeatInterval, config.Heartbeat.ServerURL)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in heartbeat goroutine: %v", r)
				}
			}()

			// 立即发送一次心跳
			sendHeartbeat(config.Heartbeat.ServerURL, config.Heartbeat.Token)

			// 定时发送心跳
			ticker := time.NewTicker(heartbeatInterval)
			defer ticker.Stop()

			for range ticker.C {
				sendHeartbeat(config.Heartbeat.ServerURL, config.Heartbeat.Token)
			}
		}()
	} else {
		log.Println("Heartbeat disabled: no server URL configured")
	}

	// 持续运行HTTP服务
	server := &http.Server{
		Addr:    ":12580",
		Handler: nil,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	log.Println("Agent started successfully on port 12580")
	return nil
}
