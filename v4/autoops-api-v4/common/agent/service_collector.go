package agent

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// monitorServices 监控运行的服务(类似check_services.sh)
func monitorServices() {
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
			javaCmd := exec.Command("bash", "-c", "ps -ef | grep '[j]ava' | awk '{for(i=8;i<=NF;i++) printf \"%s \", $i; print \"\"}'")
			javaOut, err := javaCmd.Output()
			if err == nil {
				// 查找jar/war文件
				jarCmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | grep -oP '(\\S+\\.jar|\\S+\\.war)' | sort -u", string(javaOut)))
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
			pythonCmd := exec.Command("bash", "-c", "ps -ef | grep '[p]ython' | awk '{for(i=8;i<=NF;i++) printf \"%s \", $i; print \"\"}'")
			pythonOut, err := pythonCmd.Output()
			if err == nil {
				pyCmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | grep -oP '\\s(\\S+\\.py)\\b' | sed 's/^\\s*//' | sort -u", string(pythonOut)))
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
			nodeCmd := exec.Command("bash", "-c", "ps -ef | grep '[n]ode' | awk '{for(i=8;i<=NF;i++) printf \"%s \", $i; print \"\"}'")
			nodeOut, err := nodeCmd.Output()
			if err == nil {
				jsCmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | grep -oP '\\s(\\S+\\.js)\\b' | sed 's/^\\s*//' | sort -u", string(nodeOut)))
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
		cmd = exec.Command("bash", "-c", fmt.Sprintf("ps aux | grep 'java' | grep '%s' | grep -v grep | awk '{cpu+=$3; mem+=$4} END {printf \"%%f %%f\", cpu, mem}'", jarName))
	} else if strings.HasPrefix(serviceName, "python: ") {
		scriptName := strings.TrimPrefix(serviceName, "python: ")
		cmd = exec.Command("bash", "-c", fmt.Sprintf("ps aux | grep 'python' | grep '%s' | grep -v grep | awk '{cpu+=$3; mem+=$4} END {printf \"%%f %%f\", cpu, mem}'", scriptName))
	} else if strings.HasPrefix(serviceName, "node: ") {
		scriptName := strings.TrimPrefix(serviceName, "node: ")
		cmd = exec.Command("bash", "-c", fmt.Sprintf("ps aux | grep 'node' | grep '%s' | grep -v grep | awk '{cpu+=$3; mem+=$4} END {printf \"%%f %%f\", cpu, mem}'", scriptName))
	} else {
		// 参考check_services.sh的实现：匹配进程名在第11列
		cmd = exec.Command("bash", "-c", fmt.Sprintf("ps aux | awk -v svc='%s' '$11 ~ svc || $11 ~ \"/\"svc\"$\" {cpu+=$3; mem+=$4} END {if(cpu>0 || mem>0) printf \"%%f %%f\", cpu, mem; else print \"0 0\"}'", serviceName))
	}

	output, err := cmd.Output()
	if err != nil {
		log.Printf("[本地DEBUG] 服务[%s] 命令执行失败: %v", serviceName, err)
		return 0, 0
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return 0, 0
	}

	fmt.Sscanf(result, "%f %f", &cpu, &mem)
	return cpu, mem
}
