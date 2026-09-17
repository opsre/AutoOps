package agent

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

var hostname string

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// getSystemUptime 获取系统运行时间(秒)
func getSystemUptime() (float64, error) {
	// 读取/proc/uptime文件
	content, err := exec.Command("cat", "/proc/uptime").Output()
	if err != nil {
		return 0, err
	}

	// /proc/uptime第一个值是系统运行时间(秒)
	fields := strings.Fields(string(content))
	if len(fields) > 0 {
		uptimeSeconds, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return 0, err
		}
		return uptimeSeconds, nil
	}

	return 0, fmt.Errorf("无法解析uptime")
}

func monitorSystemMetrics(intervalSeconds float64) {
	// 采集CPU指标(总使用率)
	if percents, err := cpu.Percent(0, false); err == nil { // false表示获取总CPU使用率
		if len(percents) > 0 {
			cpuTotalUsage.WithLabelValues(hostname).Set(percents[0])
		}
	}

	// 采集CPU分类使用率(用户态/内核态) - 使用增量计算
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
				cpuUserUsage.WithLabelValues(hostname).Set(userPercent)
				cpuSystemUsage.WithLabelValues(hostname).Set(systemPercent)
			}
		}
		// 保存当前值用于下次计算
		lastCPUTime = t
	}

	// 采集系统负载
	if load, err := load.Avg(); err == nil {
		systemLoad.WithLabelValues(hostname, "1min").Set(load.Load1)
		systemLoad.WithLabelValues(hostname, "5min").Set(load.Load5)
		systemLoad.WithLabelValues(hostname, "15min").Set(load.Load15)
	}

	// 采集内存指标
	if memInfo, err := mem.VirtualMemory(); err == nil {
		memUsage.WithLabelValues(hostname).Set(memInfo.UsedPercent)
		// GB单位的内存指标（使用Total-Available计算真实已用内存，与free命令一致）
		realUsed := float64(memInfo.Total - memInfo.Available)
		memUsedGB.WithLabelValues(hostname).Set(realUsed / 1024 / 1024 / 1024)
		memTotalGB.WithLabelValues(hostname).Set(float64(memInfo.Total) / 1024 / 1024 / 1024)
		memFreeGB.WithLabelValues(hostname).Set(float64(memInfo.Available) / 1024 / 1024 / 1024)
	}

	// 采集交换内存
	if swapInfo, err := mem.SwapMemory(); err == nil {
		// GB单位的交换内存
		swapUsedGB.WithLabelValues(hostname).Set(float64(swapInfo.Used) / 1024 / 1024 / 1024)
	}

	// 采集系统运行时间（天）
	if uptimeSeconds, err := getSystemUptime(); err == nil {
		// 只保留天数
		systemUptimeDays.WithLabelValues(hostname).Set(uptimeSeconds / 86400)
	}

	// 采集连接数（使用ss命令统计，更准确）
	// 统计TCP连接数（排除表头）
	if tcpOut, err := exec.Command("bash", "-c", "ss -tan 2>/dev/null | tail -n +2 | wc -l").Output(); err == nil {
		if tcpCount, err := strconv.Atoi(strings.TrimSpace(string(tcpOut))); err == nil {
			tcpConnections.WithLabelValues(hostname).Set(float64(tcpCount))
		}
	}

	// 统计UDP连接数（排除表头）
	if udpOut, err := exec.Command("bash", "-c", "ss -uan 2>/dev/null | tail -n +2 | wc -l").Output(); err == nil {
		if udpCount, err := strconv.Atoi(strings.TrimSpace(string(udpOut))); err == nil {
			udpConnections.WithLabelValues(hostname).Set(float64(udpCount))
		}
	}

	// 采集文件打开数量
	if fileNrOut, err := exec.Command("bash", "-c", "cat /proc/sys/fs/file-nr | awk '{print $1}'").Output(); err == nil {
		if openFiles, err := strconv.Atoi(strings.TrimSpace(string(fileNrOut))); err == nil {
			openFileDescriptors.WithLabelValues(hostname).Set(float64(openFiles))
		}
	}

	// 采集磁盘使用率(只采集根分区)
	if partitions, err := disk.Partitions(false); err == nil {
		for _, part := range partitions {
			if part.Mountpoint == "/" { // 只处理根分区
				if usage, err := disk.Usage(part.Mountpoint); err == nil {
					diskUsage.WithLabelValues(
						hostname,
						part.Device,
						part.Mountpoint,
					).Set(usage.UsedPercent)
					diskUsedBytes.WithLabelValues(
						hostname,
						part.Device,
						part.Mountpoint,
					).Set(float64(usage.Used) / 1024 / 1024 / 1024) // 转换为GB
					diskTotalBytes.WithLabelValues(
						hostname,
						part.Device,
						part.Mountpoint,
					).Set(float64(usage.Total) / 1024 / 1024 / 1024) // 转换为GB
				}
			}
		}
	}

	// 采集磁盘IO指标(吞吐量和IOPS)
	if diskIO, err := disk.IOCounters(); err == nil {
		if lastDiskIO != nil {
			for device, stats := range diskIO {
				if lastStats, ok := lastDiskIO[device]; ok {
					// 计算间隔内的吞吐速率(KB/s)
					readRate := float64(stats.ReadBytes-lastStats.ReadBytes) / 1024 / intervalSeconds
					writeRate := float64(stats.WriteBytes-lastStats.WriteBytes) / 1024 / intervalSeconds
					diskReadRate.WithLabelValues(hostname, device).Set(readRate)
					diskWriteRate.WithLabelValues(hostname, device).Set(writeRate)

					// 计算IOPS (次/秒)
					readIOPS := float64(stats.ReadCount-lastStats.ReadCount) / intervalSeconds
					writeIOPS := float64(stats.WriteCount-lastStats.WriteCount) / intervalSeconds
					diskReadIOPS.WithLabelValues(hostname, device).Set(readIOPS)
					diskWriteIOPS.WithLabelValues(hostname, device).Set(writeIOPS)
				}
			}
		}
		lastDiskIO = diskIO
	}

	// 采集进程总数
	if processes, err := process.Processes(); err == nil {
		totalProcesses.WithLabelValues(hostname).Set(float64(len(processes)))
	}

	// 采集网络指标
	if netStats, err := psnet.IOCounters(true); err == nil { // true表示获取所有接口
		if len(netStats) > 0 {
			// 找到活动的物理网络接口(排除docker、虚拟网桥等)
			var mainInterface *psnet.IOCountersStat
			for i := range netStats {
				name := netStats[i].Name
				// 跳过虚拟接口和docker网桥
				if strings.HasPrefix(name, "docker") ||
					strings.HasPrefix(name, "br-") ||
					strings.HasPrefix(name, "veth") ||
					strings.HasPrefix(name, "lo") {
					continue
				}
				// 优先选择有流量的接口
				if netStats[i].BytesRecv > 0 || netStats[i].BytesSent > 0 {
					mainInterface = &netStats[i]
					break
				}
			}
			// 如果没有找到活动接口，使用第一个非虚拟接口
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
			// 如果还是没有找到，使用第一个接口
			if mainInterface == nil && len(netStats) > 0 {
				mainInterface = &netStats[0]
			}

			// 记录详细统计
			log.Printf("Selected interface %s - BytesRecv: %d, BytesSent: %d",
				mainInterface.Name, mainInterface.BytesRecv, mainInterface.BytesSent)

			if lastNetIO.Name != "" {
				recvRate := float64(mainInterface.BytesRecv-lastNetIO.BytesRecv) / 1024 / intervalSeconds
				sentRate := float64(mainInterface.BytesSent-lastNetIO.BytesSent) / 1024 / intervalSeconds
				log.Printf("Network rates - Receive: %.2f KB/s, Send: %.2f KB/s", recvRate, sentRate)
				netRecv.WithLabelValues(hostname).Set(recvRate)
				netSent.WithLabelValues(hostname).Set(sentRate)
			}
			lastNetIO = *mainInterface
		} else {
			log.Println("No network interfaces found")
		}
	} else {
		log.Printf("Failed to get network stats: %v", err)
	}
}

func init() {
	hostname = getHostname()
	if hostname == "" {
		hostname = "unknown"
	}

	// 设置agent版本信息指标
	agentInfo.WithLabelValues(hostname, AgentVersion, AgentVersionTag).Set(1)
}
