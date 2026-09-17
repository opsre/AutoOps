package agent

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

func monitorPing(target string) {
	start := time.Now()

	// 使用系统ping命令
	cmd := exec.Command("ping", "-c", "1", target)
	err := cmd.Run()
	duration := time.Since(start).Milliseconds()

	status := 0
	if err == nil {
		status = 1
	}

	pingResponseTime.WithLabelValues(hostname, target).Set(float64(duration))
	pingStatus.WithLabelValues(hostname, target).Set(float64(status))
	log.Printf("Ping %s - Status: %d, Response: %dms", target, status, duration)
}

// monitorTCPPorts TCP端口监控
func monitorTCPPorts() {
	// 执行netstat命令获取TCP监听端口
	cmd := exec.Command("netstat", "-ntlp")
	output, err := cmd.Output()
	if err != nil {
		//log.Printf("Failed to execute netstat: %%v", err)
		return
	}

	// 解析netstat输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines[2:] { // 跳过标题行
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 7 && fields[5] == "LISTEN" {
			// 解析本地地址和端口
			localAddr := fields[3]
			if strings.Contains(localAddr, ":") {
				parts := strings.Split(localAddr, ":")
				if len(parts) >= 2 {
					port := parts[len(parts)-1]

					// 解析进程信息
					processInfo := fields[6]
					pid := ""
					service := ""

					if processInfo != "-" && strings.Contains(processInfo, "/") {
						processParts := strings.Split(processInfo, "/")
						if len(processParts) >= 2 {
							pid = processParts[0]
							service = processParts[1]
						}
					}

					// 如果无法获取服务名，使用端口号作为服务名
					if service == "" {
						service = "port-" + port
					}

					// 记录TCP端口状态
					tcpPortStatus.WithLabelValues(hostname, port, service, pid).Set(1)

					log.Printf("TCP Port Monitor: %s:%s - Service: %s, PID: %s", hostname, port, service, pid)
				}
			}
		}
	}
}
