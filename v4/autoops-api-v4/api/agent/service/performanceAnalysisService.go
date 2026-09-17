package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	agentModel "dodevops-api/api/agent/model"
	aiDao "dodevops-api/api/ai/dao"
	aiModel "dodevops-api/api/ai/model"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	monitorService "dodevops-api/api/monitor/service"
	"dodevops-api/common/util"
	"dodevops-api/pkg/log"
)

// PerformanceAnalysisService 性能分析服务
type PerformanceAnalysisService struct {
	hostDao    cmdbDao.CmdbHostDao
	ecsAuthDao monitorService.EcsAuthDao
	sshUtil    *util.SSHUtil
	aiModelDao aiDao.AiModelDao
}

// NewPerformanceAnalysisService 创建性能分析服务实例
func NewPerformanceAnalysisService() *PerformanceAnalysisService {
	return &PerformanceAnalysisService{
		hostDao:    cmdbDao.NewCmdbHostDao(),
		ecsAuthDao: monitorService.NewEcsAuthDao(),
		sshUtil:    util.NewSSHUtil(),
		aiModelDao: aiDao.NewAiModelDao(),
	}
}

// AnalyzeServerPressureByIP 通过IP分析服务器压力
func (s *PerformanceAnalysisService) AnalyzeServerPressureByIP(ip string) (*agentModel.PressureAnalysisResult, error) {
	// 1. 通过IP查询主机
	host, err := s.hostDao.GetCmdbHostByIP(ip)
	if err != nil {
		return nil, fmt.Errorf("未找到IP为 %s 的服务器: %v", ip, err)
	}

	return s.analyzeServerPressure(&host)
}

// AnalyzeServerPressureByName 通过名称分析服务器压力
func (s *PerformanceAnalysisService) AnalyzeServerPressureByName(name string) (*agentModel.PressureAnalysisResult, error) {
	// 1. 通过名称查询主机
	host, err := s.hostDao.GetCmdbHostByName(name)
	if err != nil {
		return nil, fmt.Errorf("未找到名称为 %s 的服务器: %v", name, err)
	}

	return s.analyzeServerPressure(&host)
}

// analyzeServerPressure 分析服务器压力（内部方法）
func (s *PerformanceAnalysisService) analyzeServerPressure(host *cmdbModel.CmdbHost) (*agentModel.PressureAnalysisResult, error) {
	result := &agentModel.PressureAnalysisResult{
		HostID:      host.ID,
		HostName:    host.HostName,
		HostIP:      host.SSHIP,
		CollectTime: time.Now(),
		Errors:      []string{},
	}

	// 2. 获取SSH认证信息
	sshConfig, err := s.getSSHConfig(host)
	if err != nil {
		return nil, fmt.Errorf("获取SSH认证信息失败: %v", err)
	}

	// 3. 收集性能数据
	s.collectCPUInfo(sshConfig, result)
	s.collectMemoryInfo(sshConfig, result)
	s.collectDiskInfo(sshConfig, result)
	s.collectLoadInfo(sshConfig, result)
	s.collectTopProcesses(sshConfig, result)
	s.collectNetworkInfo(sshConfig, result)

	// 4. 评估压力等级
	result.PressureLevel = s.evaluatePressureLevel(result)

	return result, nil
}

// getSSHConfig 获取SSH配置
func (s *PerformanceAnalysisService) getSSHConfig(host *cmdbModel.CmdbHost) (*util.SSHConfig, error) {
	// 获取SSH认证信息
	ecsAuth, err := s.ecsAuthDao.GetEcsAuthById(host.SSHKeyID)
	if err != nil {
		return nil, fmt.Errorf("获取SSH认证信息失败: %v", err)
	}

	// 构建SSH配置
	sshConfig := &util.SSHConfig{
		IP:       host.SSHIP,
		Port:     host.SSHPort,
		Username: host.SSHName,
		Type:     ecsAuth.Type,
	}

	// 根据认证类型设置认证信息
	switch ecsAuth.Type {
	case 1: // 密码认证
		sshConfig.Password = ecsAuth.Password
	case 2: // 私钥认证
		sshConfig.PublicKey = ecsAuth.PublicKey
	case 4: // 跳板机认证
		sshConfig.BastionHost = ecsAuth.BastionHost
		sshConfig.BastionPort = ecsAuth.BastionPort
		sshConfig.BastionUsername = ecsAuth.BastionUsername
		sshConfig.BastionAuthType = ecsAuth.BastionAuthType
		sshConfig.BastionPassword = ecsAuth.BastionPassword
		sshConfig.BastionPrivateKey = ecsAuth.BastionPrivateKey
		sshConfig.TargetAuthType = ecsAuth.TargetAuthType
		sshConfig.TargetPassword = ecsAuth.TargetPassword
		sshConfig.TargetPrivateKey = ecsAuth.TargetPrivateKey
	}

	return sshConfig, nil
}

// collectCPUInfo 收集CPU信息
func (s *PerformanceAnalysisService) collectCPUInfo(sshConfig *util.SSHConfig, result *agentModel.PressureAnalysisResult) {
	cpuInfo := make(map[string]interface{})

	// CPU核心数
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "nproc")
	if err == nil {
		cpuInfo["cores"] = strings.TrimSpace(output)
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("获取CPU核心数失败: %v", err))
	}

	// CPU使用率
	output, err = s.sshUtil.ExecuteRemoteCommand(sshConfig, "top -bn1 | grep 'Cpu(s)' | awk '{print $2}' | cut -d'%' -f1")
	if err == nil {
		cpuInfo["usage"] = strings.TrimSpace(output)
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("获取CPU使用率失败: %v", err))
	}

	// CPU详细信息
	output, err = s.sshUtil.ExecuteRemoteCommand(sshConfig, "lscpu | grep -E 'Model name|CPU MHz|CPU max MHz'")
	if err == nil {
		cpuInfo["details"] = strings.TrimSpace(output)
	}

	result.CPUInfo = cpuInfo
}

// collectMemoryInfo 收集内存信息
func (s *PerformanceAnalysisService) collectMemoryInfo(sshConfig *util.SSHConfig, result *agentModel.PressureAnalysisResult) {
	memInfo := make(map[string]interface{})

	// 内存使用情况
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "free -m | grep Mem | awk '{print $2,$3,$4,$7}'")
	if err == nil {
		parts := strings.Fields(strings.TrimSpace(output))
		if len(parts) >= 4 {
			total, _ := strconv.ParseFloat(parts[0], 64)
			used, _ := strconv.ParseFloat(parts[1], 64)
			free, _ := strconv.ParseFloat(parts[2], 64)
			available, _ := strconv.ParseFloat(parts[3], 64)

			memInfo["total_mb"] = total
			memInfo["used_mb"] = used
			memInfo["free_mb"] = free
			memInfo["available_mb"] = available
			if total > 0 {
				memInfo["usage_percent"] = fmt.Sprintf("%.2f", (used/total)*100)
			}
		}
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("获取内存信息失败: %v", err))
	}

	result.MemoryInfo = memInfo
}

// collectDiskInfo 收集磁盘信息
func (s *PerformanceAnalysisService) collectDiskInfo(sshConfig *util.SSHConfig, result *agentModel.PressureAnalysisResult) {
	diskInfo := make(map[string]interface{})

	// 磁盘使用情况
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "df -h | grep -vE '^Filesystem|tmpfs|cdrom'")
	if err == nil {
		diskInfo["usage"] = strings.TrimSpace(output)
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("获取磁盘信息失败: %v", err))
	}

	// 根分区使用率
	output, err = s.sshUtil.ExecuteRemoteCommand(sshConfig, "df -h / | tail -1 | awk '{print $5}' | cut -d'%' -f1")
	if err == nil {
		diskInfo["root_usage_percent"] = strings.TrimSpace(output)
	}

	result.DiskInfo = diskInfo
}

// collectLoadInfo 收集负载信息
func (s *PerformanceAnalysisService) collectLoadInfo(sshConfig *util.SSHConfig, result *agentModel.PressureAnalysisResult) {
	loadInfo := make(map[string]interface{})

	// 系统负载
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "uptime | awk -F'load average:' '{print $2}'")
	if err == nil {
		loadInfo["load_average"] = strings.TrimSpace(output)
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("获取系统负载失败: %v", err))
	}

	result.LoadInfo = loadInfo
}

// collectTopProcesses 收集TOP进程信息
func (s *PerformanceAnalysisService) collectTopProcesses(sshConfig *util.SSHConfig, result *agentModel.PressureAnalysisResult) {
	processes := []map[string]interface{}{}

	// 获取CPU占用最高的10个进程
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "ps aux --sort=-%cpu | head -11 | tail -10")
	if err == nil {
		lines := strings.Split(strings.TrimSpace(output), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 11 {
				process := map[string]interface{}{
					"user":    fields[0],
					"pid":     fields[1],
					"cpu":     fields[2],
					"mem":     fields[3],
					"command": strings.Join(fields[10:], " "),
				}
				processes = append(processes, process)
			}
		}
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("获取TOP进程失败: %v", err))
	}

	result.TopProcesses = processes
}

// collectNetworkInfo 收集网络信息
func (s *PerformanceAnalysisService) collectNetworkInfo(sshConfig *util.SSHConfig, result *agentModel.PressureAnalysisResult) {
	netInfo := make(map[string]interface{})

	// 网络连接数
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "ss -s | grep 'TCP:' | awk '{print $2}'")
	if err == nil {
		netInfo["tcp_connections"] = strings.TrimSpace(output)
	}

	// 网络接口信息
	output, err = s.sshUtil.ExecuteRemoteCommand(sshConfig, "ip -s link | grep -E '^[0-9]+:|RX:|TX:' | head -20")
	if err == nil {
		netInfo["interfaces"] = strings.TrimSpace(output)
	}

	result.NetworkInfo = netInfo
}

// evaluatePressureLevel 评估压力等级
func (s *PerformanceAnalysisService) evaluatePressureLevel(result *agentModel.PressureAnalysisResult) string {
	score := 0

	// CPU使用率评分
	if cpuUsage, ok := result.CPUInfo["usage"].(string); ok {
		if usage, err := strconv.ParseFloat(cpuUsage, 64); err == nil {
			if usage > 80 {
				score += 3
			} else if usage > 60 {
				score += 2
			} else if usage > 40 {
				score += 1
			}
		}
	}

	// 内存使用率评分
	if memUsage, ok := result.MemoryInfo["usage_percent"].(string); ok {
		if usage, err := strconv.ParseFloat(memUsage, 64); err == nil {
			if usage > 85 {
				score += 3
			} else if usage > 70 {
				score += 2
			} else if usage > 50 {
				score += 1
			}
		}
	}

	// 磁盘使用率评分
	if diskUsage, ok := result.DiskInfo["root_usage_percent"].(string); ok {
		if usage, err := strconv.ParseFloat(diskUsage, 64); err == nil {
			if usage > 90 {
				score += 3
			} else if usage > 80 {
				score += 2
			} else if usage > 70 {
				score += 1
			}
		}
	}

	// 根据总分判断压力等级
	if score >= 6 {
		return "高压力"
	} else if score >= 3 {
		return "中等压力"
	}
	return "正常"
}

// GeneratePressureReport 生成压力分析报告
func (s *PerformanceAnalysisService) GeneratePressureReport(result *agentModel.PressureAnalysisResult) string {
	report := fmt.Sprintf("# 服务器压力巡检报告\n\n")
	report += fmt.Sprintf("**主机信息**\n")
	report += fmt.Sprintf("- 主机名称: %s\n", result.HostName)
	report += fmt.Sprintf("- 主机IP: %s\n", result.HostIP)
	report += fmt.Sprintf("- 检查时间: %s\n", result.CollectTime.Format("2006-01-02 15:04:05"))
	report += fmt.Sprintf("- 压力等级: **%s**\n\n", result.PressureLevel)

	// CPU信息
	report += fmt.Sprintf("## CPU信息\n")
	if cores, ok := result.CPUInfo["cores"]; ok {
		report += fmt.Sprintf("- CPU核心数: %v\n", cores)
	}
	if usage, ok := result.CPUInfo["usage"]; ok {
		report += fmt.Sprintf("- CPU使用率: %v%%\n", usage)
	}
	report += "\n"

	// 内存信息
	report += fmt.Sprintf("## 内存信息\n")
	if total, ok := result.MemoryInfo["total_mb"]; ok {
		report += fmt.Sprintf("- 总内存: %.0f MB\n", total)
	}
	if used, ok := result.MemoryInfo["used_mb"]; ok {
		report += fmt.Sprintf("- 已使用: %.0f MB\n", used)
	}
	if usage, ok := result.MemoryInfo["usage_percent"]; ok {
		report += fmt.Sprintf("- 使用率: %v%%\n", usage)
	}
	report += "\n"

	// 磁盘信息
	report += fmt.Sprintf("## 磁盘信息\n")
	if usage, ok := result.DiskInfo["root_usage_percent"]; ok {
		report += fmt.Sprintf("- 根分区使用率: %v%%\n", usage)
	}
	if diskUsage, ok := result.DiskInfo["usage"]; ok {
		report += fmt.Sprintf("```\n%v\n```\n", diskUsage)
	}
	report += "\n"

	// 系统负载
	report += fmt.Sprintf("## 系统负载\n")
	if load, ok := result.LoadInfo["load_average"]; ok {
		report += fmt.Sprintf("- 负载平均值: %v\n", load)
	}
	report += "\n"

	// TOP进程
	report += fmt.Sprintf("## TOP 10 进程（按CPU使用率）\n")
	report += fmt.Sprintf("| 用户 | PID | CPU%% | MEM%% | 命令 |\n")
	report += fmt.Sprintf("|------|-----|------|------|------|\n")
	for _, proc := range result.TopProcesses {
		report += fmt.Sprintf("| %v | %v | %v | %v | %v |\n",
			proc["user"], proc["pid"], proc["cpu"], proc["mem"], proc["command"])
	}
	report += "\n"

	// AI智能分析
	report += fmt.Sprintf("## 🤖 AI智能分析\n\n")
	aiAnalysis := s.GenerateAIAnalysisForPressure(result)
	report += aiAnalysis + "\n\n"

	// 错误信息
	if len(result.Errors) > 0 {
		report += fmt.Sprintf("## 错误信息\n")
		for _, err := range result.Errors {
			report += fmt.Sprintf("- %s\n", err)
		}
	}

	return report
}

// ==================== 内存问题分析 ====================

// AnalyzeMemoryIssueByName 通过名称分析内存问题
func (s *PerformanceAnalysisService) AnalyzeMemoryIssueByName(name string) (*agentModel.MemoryAnalysisResult, error) {
	// 1. 通过名称查询主机
	host, err := s.hostDao.GetCmdbHostByName(name)
	if err != nil {
		return nil, fmt.Errorf("未找到名称为 %s 的服务器: %v", name, err)
	}

	return s.analyzeMemoryIssue(&host)
}

// AnalyzeMemoryIssueByIP 通过IP分析内存问题
func (s *PerformanceAnalysisService) AnalyzeMemoryIssueByIP(ip string) (*agentModel.MemoryAnalysisResult, error) {
	// 1. 通过IP查询主机
	host, err := s.hostDao.GetCmdbHostByIP(ip)
	if err != nil {
		return nil, fmt.Errorf("未找到IP为 %s 的服务器: %v", ip, err)
	}

	return s.analyzeMemoryIssue(&host)
}

// analyzeMemoryIssue 分析内存问题（内部方法）
func (s *PerformanceAnalysisService) analyzeMemoryIssue(host *cmdbModel.CmdbHost) (*agentModel.MemoryAnalysisResult, error) {
	result := &agentModel.MemoryAnalysisResult{
		HostID:      host.ID,
		HostName:    host.HostName,
		HostIP:      host.SSHIP,
		CollectTime: time.Now(),
		Issues:      []string{},
	}

	// 2. 获取SSH认证信息
	sshConfig, err := s.getSSHConfig(host)
	if err != nil {
		return nil, fmt.Errorf("获取SSH认证信息失败: %v", err)
	}

	// 3. 收集详细的内存信息
	s.collectDetailedMemoryInfo(sshConfig, result)

	// 4. 收集内存占用最高的进程
	s.collectMemoryTopProcesses(sshConfig, result)

	// 5. 分析内存使用模式
	result.MemoryPattern = s.analyzeMemoryPattern(result)

	// 6. 识别内存问题
	result.Issues = s.identifyMemoryIssues(result)

	return result, nil
}

// collectDetailedMemoryInfo 收集详细的内存信息
func (s *PerformanceAnalysisService) collectDetailedMemoryInfo(sshConfig *util.SSHConfig, result *agentModel.MemoryAnalysisResult) {
	memInfo := make(map[string]interface{})

	// 基本内存信息
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "free -m | grep Mem | awk '{print $2,$3,$4,$5,$6,$7}'")
	if err == nil {
		parts := strings.Fields(strings.TrimSpace(output))
		if len(parts) >= 6 {
			total, _ := strconv.ParseFloat(parts[0], 64)
			used, _ := strconv.ParseFloat(parts[1], 64)
			free, _ := strconv.ParseFloat(parts[2], 64)
			shared, _ := strconv.ParseFloat(parts[3], 64)
			bufferCache, _ := strconv.ParseFloat(parts[4], 64)
			available, _ := strconv.ParseFloat(parts[5], 64)

			memInfo["total_mb"] = total
			memInfo["used_mb"] = used
			memInfo["free_mb"] = free
			memInfo["shared_mb"] = shared
			memInfo["buffer_cache_mb"] = bufferCache
			memInfo["available_mb"] = available
			if total > 0 {
				memInfo["usage_percent"] = fmt.Sprintf("%.2f", (used/total)*100)
				memInfo["available_percent"] = fmt.Sprintf("%.2f", (available/total)*100)
			}
		}
	}

	// Swap信息
	output, err = s.sshUtil.ExecuteRemoteCommand(sshConfig, "free -m | grep Swap | awk '{print $2,$3,$4}'")
	if err == nil {
		parts := strings.Fields(strings.TrimSpace(output))
		if len(parts) >= 3 {
			swapTotal, _ := strconv.ParseFloat(parts[0], 64)
			swapUsed, _ := strconv.ParseFloat(parts[1], 64)
			swapFree, _ := strconv.ParseFloat(parts[2], 64)

			memInfo["swap_total_mb"] = swapTotal
			memInfo["swap_used_mb"] = swapUsed
			memInfo["swap_free_mb"] = swapFree
			if swapTotal > 0 {
				memInfo["swap_usage_percent"] = fmt.Sprintf("%.2f", (swapUsed/swapTotal)*100)
			}
		}
	}

	// 内存详细统计
	output, err = s.sshUtil.ExecuteRemoteCommand(sshConfig, "cat /proc/meminfo | grep -E 'MemTotal|MemFree|MemAvailable|Buffers|Cached|SwapCached|Active|Inactive|Slab'")
	if err == nil {
		memInfo["meminfo_details"] = strings.TrimSpace(output)
	}

	result.MemoryInfo = memInfo
}

// collectMemoryTopProcesses 收集内存占用最高的进程
func (s *PerformanceAnalysisService) collectMemoryTopProcesses(sshConfig *util.SSHConfig, result *agentModel.MemoryAnalysisResult) {
	processes := []map[string]interface{}{}

	// 获取内存占用最高的15个进程
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "ps aux --sort=-%mem | head -16 | tail -15")
	if err == nil {
		lines := strings.Split(strings.TrimSpace(output), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 11 {
				process := map[string]interface{}{
					"user":    fields[0],
					"pid":     fields[1],
					"cpu":     fields[2],
					"mem":     fields[3],
					"vsz":     fields[4],
					"rss":     fields[5],
					"command": strings.Join(fields[10:], " "),
				}
				processes = append(processes, process)
			}
		}
	}

	result.TopProcesses = processes
}

// analyzeMemoryPattern 分析内存使用模式
func (s *PerformanceAnalysisService) analyzeMemoryPattern(result *agentModel.MemoryAnalysisResult) string {
	var patterns []string

	// 分析内存使用率
	if usageStr, ok := result.MemoryInfo["usage_percent"].(string); ok {
		if usage, err := strconv.ParseFloat(usageStr, 64); err == nil {
			if usage > 90 {
				patterns = append(patterns, "内存使用率极高(>90%)")
			} else if usage > 80 {
				patterns = append(patterns, "内存使用率较高(>80%)")
			} else if usage > 70 {
				patterns = append(patterns, "内存使用率中等(>70%)")
			} else {
				patterns = append(patterns, "内存使用率正常")
			}
		}
	}

	// 分析Swap使用情况
	if swapUsageStr, ok := result.MemoryInfo["swap_usage_percent"].(string); ok {
		if swapUsage, err := strconv.ParseFloat(swapUsageStr, 64); err == nil {
			if swapUsage > 50 {
				patterns = append(patterns, "Swap使用率高(>50%),可能存在内存不足")
			} else if swapUsage > 20 {
				patterns = append(patterns, "Swap使用率中等(>20%)")
			}
		}
	}

	// 分析可用内存
	if availableStr, ok := result.MemoryInfo["available_percent"].(string); ok {
		if available, err := strconv.ParseFloat(availableStr, 64); err == nil {
			if available < 10 {
				patterns = append(patterns, "可用内存极低(<10%)")
			} else if available < 20 {
				patterns = append(patterns, "可用内存较低(<20%)")
			}
		}
	}

	if len(patterns) == 0 {
		return "内存使用正常"
	}

	return strings.Join(patterns, "; ")
}

// identifyMemoryIssues 识别内存问题
func (s *PerformanceAnalysisService) identifyMemoryIssues(result *agentModel.MemoryAnalysisResult) []string {
	issues := []string{}

	// 检查内存使用率
	if usageStr, ok := result.MemoryInfo["usage_percent"].(string); ok {
		if usage, err := strconv.ParseFloat(usageStr, 64); err == nil {
			if usage > 90 {
				issues = append(issues, fmt.Sprintf("内存使用率过高: %.2f%%", usage))
			}
		}
	}

	// 检查可用内存
	if availableMB, ok := result.MemoryInfo["available_mb"].(float64); ok {
		if availableMB < 500 {
			issues = append(issues, fmt.Sprintf("可用内存不足: %.0f MB", availableMB))
		}
	}

	// 检查Swap使用
	if swapUsedMB, ok := result.MemoryInfo["swap_used_mb"].(float64); ok {
		if swapUsedMB > 1000 {
			issues = append(issues, fmt.Sprintf("Swap使用过多: %.0f MB,可能导致性能下降", swapUsedMB))
		}
	}

	// 分析TOP进程
	if len(result.TopProcesses) > 0 {
		topProcess := result.TopProcesses[0]
		if memStr, ok := topProcess["mem"].(string); ok {
			if mem, err := strconv.ParseFloat(memStr, 64); err == nil && mem > 30 {
				issues = append(issues, fmt.Sprintf("发现高内存占用进程: %v (PID: %v, 内存: %.1f%%)",
					topProcess["command"], topProcess["pid"], mem))
			}
		}
	}

	// 检查是否有多个高内存进程
	highMemProcessCount := 0
	for _, proc := range result.TopProcesses {
		if memStr, ok := proc["mem"].(string); ok {
			if mem, err := strconv.ParseFloat(memStr, 64); err == nil && mem > 10 {
				highMemProcessCount++
			}
		}
	}
	if highMemProcessCount > 5 {
		issues = append(issues, fmt.Sprintf("存在%d个高内存占用进程(>10%%)", highMemProcessCount))
	}

	return issues
}

// GenerateMemoryAnalysisReport 生成内存分析报告
func (s *PerformanceAnalysisService) GenerateMemoryAnalysisReport(result *agentModel.MemoryAnalysisResult) string {
	report := fmt.Sprintf("# 服务器内存分析报告\n\n")
	report += fmt.Sprintf("**主机信息**\n")
	report += fmt.Sprintf("- 主机名称: %s\n", result.HostName)
	report += fmt.Sprintf("- 主机IP: %s\n", result.HostIP)
	report += fmt.Sprintf("- 分析时间: %s\n\n", result.CollectTime.Format("2006-01-02 15:04:05"))

	// 内存使用模式
	report += fmt.Sprintf("## 内存使用模式\n")
	report += fmt.Sprintf("%s\n\n", result.MemoryPattern)

	// 内存详细信息
	report += fmt.Sprintf("## 内存详细信息\n")
	if total, ok := result.MemoryInfo["total_mb"]; ok {
		report += fmt.Sprintf("- 总内存: %.0f MB\n", total)
	}
	if used, ok := result.MemoryInfo["used_mb"]; ok {
		report += fmt.Sprintf("- 已使用: %.0f MB\n", used)
	}
	if available, ok := result.MemoryInfo["available_mb"]; ok {
		report += fmt.Sprintf("- 可用内存: %.0f MB\n", available)
	}
	if usage, ok := result.MemoryInfo["usage_percent"]; ok {
		report += fmt.Sprintf("- 使用率: %v%%\n", usage)
	}
	report += "\n"

	// Swap信息
	report += fmt.Sprintf("## Swap信息\n")
	if swapTotal, ok := result.MemoryInfo["swap_total_mb"]; ok {
		report += fmt.Sprintf("- Swap总量: %.0f MB\n", swapTotal)
	}
	if swapUsed, ok := result.MemoryInfo["swap_used_mb"]; ok {
		report += fmt.Sprintf("- Swap已使用: %.0f MB\n", swapUsed)
	}
	if swapUsage, ok := result.MemoryInfo["swap_usage_percent"]; ok {
		report += fmt.Sprintf("- Swap使用率: %v%%\n", swapUsage)
	}
	report += "\n"

	// TOP进程
	report += fmt.Sprintf("## TOP 15 进程（按内存使用率）\n")
	report += fmt.Sprintf("| 用户 | PID | CPU%% | MEM%% | VSZ | RSS | 命令 |\n")
	report += fmt.Sprintf("|------|-----|------|------|-----|-----|------|\n")
	for _, proc := range result.TopProcesses {
		report += fmt.Sprintf("| %v | %v | %v | %v | %v | %v | %v |\n",
			proc["user"], proc["pid"], proc["cpu"], proc["mem"], proc["vsz"], proc["rss"], proc["command"])
	}
	report += "\n"

	// 发现的问题
	if len(result.Issues) > 0 {
		report += fmt.Sprintf("## 发现的问题\n")
		for i, issue := range result.Issues {
			report += fmt.Sprintf("%d. %s\n", i+1, issue)
		}
		report += "\n"
	}

	// AI智能分析
	report += fmt.Sprintf("## 🤖 AI智能分析与优化建议\n\n")
	aiAnalysis := s.GenerateAIAnalysisForMemory(result)
	report += aiAnalysis + "\n"

	return report
}

// generateMemorySuggestions 生成内存优化建议
func (s *PerformanceAnalysisService) generateMemorySuggestions(result *agentModel.MemoryAnalysisResult) string {
	suggestions := ""

	// 根据问题类型给出建议
	for _, issue := range result.Issues {
		if strings.Contains(issue, "内存使用率过高") {
			suggestions += "- **内存使用率过高**: 建议检查是否有内存泄漏的进程，考虑增加物理内存或优化应用程序\n"
		}
		if strings.Contains(issue, "可用内存不足") {
			suggestions += "- **可用内存不足**: 建议释放不必要的缓存，关闭不需要的服务，或增加物理内存\n"
		}
		if strings.Contains(issue, "Swap使用过多") {
			suggestions += "- **Swap使用过多**: Swap使用会严重影响性能，建议增加物理内存或优化内存使用\n"
		}
		if strings.Contains(issue, "高内存占用进程") {
			suggestions += "- **高内存占用进程**: 检查该进程是否正常，考虑优化程序配置或限制其内存使用\n"
		}
		if strings.Contains(issue, "高内存占用进程(>10%)") {
			suggestions += "- **多个高内存进程**: 系统运行了多个高内存占用的进程，建议评估是否都是必需的\n"
		}
	}

	// 通用建议
	if suggestions == "" {
		suggestions = "系统内存使用正常，无需特殊优化。\n"
	} else {
		suggestions += "\n**通用优化建议**:\n"
		suggestions += "1. 定期监控内存使用情况，及时发现异常\n"
		suggestions += "2. 使用 `free -h` 和 `top` 命令实时查看内存状态\n"
		suggestions += "3. 检查系统日志，查找可能的内存相关错误\n"
		suggestions += "4. 考虑使用内存分析工具（如 valgrind）检查内存泄漏\n"
	}

	return suggestions
}

// ==================== AI智能分析 ====================

// callAIForAnalysis 调用AI大模型进行智能分析
func (s *PerformanceAnalysisService) callAIForAnalysis(systemPrompt, userPrompt string) (string, error) {
	// 1. 获取启用的AI模型（优先使用第一个启用的模型）
	models, err := s.aiModelDao.GetEnabledModels()
	if err != nil {
		log.Log().Errorf("获取AI模型失败: %v", err)
		return "", fmt.Errorf("获取AI模型失败: %v", err)
	}

	if len(models) == 0 {
		log.Log().Warnf("未找到可用的AI模型，请在系统中配置AI模型")
		return "", fmt.Errorf("未配置AI模型，请先在【AI模型管理】中添加并启用AI模型")
	}

	aiModelConfig := &models[0] // 使用第一个启用的模型（取指针）
	log.Log().Infof("使用AI模型: %s (类型: %s, 模型: %s)", aiModelConfig.Name, aiModelConfig.Type, aiModelConfig.Model)

	// 2. 构建消息
	messages := []aiModel.AiChatMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}

	// 3. 根据模型类型调用不同的API
	switch aiModelConfig.Type {
	case "阿里百炼":
		return s.callAliyunAI(aiModelConfig, messages)
	case "DeepSeek", "OpenAI":
		return s.callOpenAICompatibleAI(aiModelConfig, messages)
	default:
		return "", fmt.Errorf("不支持的模型类型: %s", aiModelConfig.Type)
	}
}

// callAliyunAI 调用阿里百炼API
func (s *PerformanceAnalysisService) callAliyunAI(model *aiModel.AiModel, messages []aiModel.AiChatMessage) (string, error) {
	// 转换消息格式
	aliMessages := make([]aiModel.AliChatMessage, len(messages))
	for i, msg := range messages {
		aliMessages[i] = aiModel.AliChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 构建请求
	req := aiModel.AliyunChatRequest{
		Model: model.Model,
		Input: aiModel.AliyunInput{
			Messages: aliMessages,
		},
		Parameters: map[string]interface{}{
			"temperature": 0.7,
		},
	}

	reqBody, _ := json.Marshal(req)

	// API URL
	apiURL := model.URL
	if apiURL == "" || apiURL == "https://bailian.cn-beijing.aliyuncs.com/api/v1/chat/completions" {
		apiURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Authorization", "Bearer "+model.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求 - 增加超时时间到120秒，因为AI分析需要较长时间
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI API错误(状态码 %d): %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var aliyunResp aiModel.AliyunChatResponse
	if err := json.Unmarshal(body, &aliyunResp); err != nil {
		return "", err
	}

	return aliyunResp.Output.Text, nil
}

// callOpenAICompatibleAI 调用OpenAI兼容的API (DeepSeek/OpenAI)
func (s *PerformanceAnalysisService) callOpenAICompatibleAI(model *aiModel.AiModel, messages []aiModel.AiChatMessage) (string, error) {
	// 转换消息格式
	openaiMessages := make([]aiModel.OpenAIChatMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = aiModel.OpenAIChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 构建请求
	req := aiModel.OpenAIChatRequest{
		Model:       model.Model,
		Messages:    openaiMessages,
		Temperature: 0.7,
	}

	reqBody, _ := json.Marshal(req)

	// 记录请求数据（用于调试）
	log.Log().Infof("[AI请求] API URL: %s, Model: %s, 消息数量: %d", model.URL, model.Model, len(messages))
	log.Log().Debugf("[AI请求] 完整请求体长度: %d bytes", len(reqBody))

	// 构建完整的API URL
	apiURL := model.URL
	if apiURL[len(apiURL)-1:] != "/" {
		apiURL += "/"
	}
	apiURL += "chat/completions"

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Authorization", "Bearer "+model.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求 - 增加超时时间到120秒，因为AI分析需要较长时间
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Log().Errorf("[AI响应] HTTP状态码: %d, 响应体: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("AI API错误(状态码 %d): %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var openaiResp aiModel.OpenAIChatResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return "", err
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("AI未返回内容")
	}

	log.Log().Infof("[AI响应] 成功获取AI分析结果")
	return openaiResp.Choices[0].Message.Content, nil
}

// GenerateAIAnalysisForPressure 为压力分析生成AI建议
func (s *PerformanceAnalysisService) GenerateAIAnalysisForPressure(result *agentModel.PressureAnalysisResult) string {
	systemPrompt := `你是一位资深的Linux系统运维专家，擅长服务器性能分析和优化。请根据用户提供的服务器监控数据，进行深度分析并给出专业的优化建议。

要求：
1. 分析要深入、专业，找出潜在问题的根本原因
2. 建议要具体、可操作，包含详细的命令或配置
3. 优先级排序，先解决最紧急的问题
4. 考虑服务器的实际使用场景
5. 输出格式清晰，使用Markdown格式`

	// 构建用户提示词，包含服务器详细数据（简化版，避免触发API限制）
	userPrompt := fmt.Sprintf(`请分析以下服务器的性能数据：

服务器信息：
- 主机名: %s
- IP地址: %s
- 压力等级: %s
- 采集时间: %s

CPU信息：
- CPU核心数: %v
- CPU使用率: %v%%

内存信息：
- 总内存: %.0f MB
- 已使用: %.0f MB
- 使用率: %v%%

磁盘信息：
- 根分区使用率: %v%%

系统负载：
- 负载平均值: %v

TOP 10进程（按CPU排序）：
`,
		result.HostName,
		result.HostIP,
		result.PressureLevel,
		result.CollectTime.Format("2006-01-02 15:04:05"),
		result.CPUInfo["cores"],
		result.CPUInfo["usage"],
		result.MemoryInfo["total_mb"],
		result.MemoryInfo["used_mb"],
		result.MemoryInfo["usage_percent"],
		result.DiskInfo["root_usage_percent"],
		result.LoadInfo["load_average"],
	)

	// 添加简化的进程信息（只保留关键数据，去除完整命令行）
	for i, proc := range result.TopProcesses {
		if i >= 10 {
			break
		}
		// 简化命令名称，只取基本名称
		command := fmt.Sprintf("%v", proc["command"])
		if len(command) > 50 {
			command = command[:50] + "..."
		}
		userPrompt += fmt.Sprintf("%d. PID:%v CPU:%v%% MEM:%v%% 命令:%s\n",
			i+1, proc["pid"], proc["cpu"], proc["mem"], command)
	}

	userPrompt += "\n请提供：\n1. 性能瓶颈分析 - 指出当前最主要的性能问题\n2. 根因分析 - 分析问题的根本原因\n3. 优化建议 - 给出3到5条具体可行的优化措施\n4. 风险提示 - 如果有潜在风险，请提前告知"

	// 调用AI分析
	aiResponse, err := s.callAIForAnalysis(systemPrompt, userPrompt)
	if err != nil {
		log.Log().Warnf("AI分析失败: %v", err)
		return fmt.Sprintf("⚠️ %s\n\n**温馨提示**：AI分析功能需要配置正确的API权限。", err.Error())
	}

	return aiResponse
}

// GenerateAIAnalysisForMemory 为内存分析生成AI建议
func (s *PerformanceAnalysisService) GenerateAIAnalysisForMemory(result *agentModel.MemoryAnalysisResult) string {
	systemPrompt := `你是一位资深的Linux内存管理专家，擅长诊断和解决内存相关问题。请根据用户提供的内存监控数据，进行深度分析并给出专业的优化建议。

要求：
1. 准确诊断内存问题类型（内存泄漏、配置不合理、缓存过多等）
2. 提供详细的排查步骤和命令
3. 给出具体的优化方案和配置建议
4. 评估是否需要扩容内存
5. 输出格式清晰，使用Markdown格式`

	// 构建简化的用户提示词（避免触发API内容检测）
	userPrompt := fmt.Sprintf(`请分析以下服务器的内存使用情况：

服务器信息：
- 主机名: %s
- IP地址: %s
- 分析时间: %s

内存使用模式：
%s

内存详细信息：
- 总内存: %.0f MB
- 已使用: %.0f MB
- 可用内存: %.0f MB
- 使用率: %v%%

Swap信息：
- Swap总量: %.0f MB
- Swap已使用: %.0f MB
- Swap使用率: %v%%

TOP 15进程（按内存排序）：
`,
		result.HostName,
		result.HostIP,
		result.CollectTime.Format("2006-01-02 15:04:05"),
		result.MemoryPattern,
		result.MemoryInfo["total_mb"],
		result.MemoryInfo["used_mb"],
		result.MemoryInfo["available_mb"],
		result.MemoryInfo["usage_percent"],
		result.MemoryInfo["swap_total_mb"],
		result.MemoryInfo["swap_used_mb"],
		result.MemoryInfo["swap_usage_percent"],
	)

	// 添加简化的进程信息（只保留关键数据，去除完整命令行）
	for i, proc := range result.TopProcesses {
		if i >= 15 {
			break
		}
		// 简化命令名称
		command := fmt.Sprintf("%v", proc["command"])
		if len(command) > 50 {
			command = command[:50] + "..."
		}
		userPrompt += fmt.Sprintf("%d. PID:%v MEM:%v%% RSS:%v KB 命令:%s\n",
			i+1, proc["pid"], proc["mem"], proc["rss"], command)
	}

	// 添加发现的问题
	if len(result.Issues) > 0 {
		userPrompt += "\n发现的问题：\n"
		for i, issue := range result.Issues {
			userPrompt += fmt.Sprintf("%d. %s\n", i+1, issue)
		}
	}

	userPrompt += "\n请提供：\n1. 内存问题诊断 - 判断是否存在内存泄漏、内存不足等问题\n2. 根因分析 - 分析内存占用高的原因\n3. 优化建议 - 给出具体的优化措施（包括命令和配置）\n4. 是否需要扩容 - 评估是否需要增加物理内存"

	// 调用AI分析
	aiResponse, err := s.callAIForAnalysis(systemPrompt, userPrompt)
	if err != nil {
		log.Log().Warnf("AI分析失败: %v", err)
		// 降级到基础建议
		return s.generateBasicMemoryAdvice(result)
	}

	return aiResponse
}

// generateBasicMemoryAdvice 生成基础内存建议（无需AI）
func (s *PerformanceAnalysisService) generateBasicMemoryAdvice(result *agentModel.MemoryAnalysisResult) string {
	advice := "### 💡 基础优化建议\n\n"

	// 根据内存使用率给建议
	if usage, ok := result.MemoryInfo["usage_percent"].(string); ok {
		if usageFloat, err := strconv.ParseFloat(usage, 64); err == nil {
			if usageFloat > 90 {
				advice += "#### 🔴 内存使用率极高 (>90%)\n\n"
				advice += "**立即处理建议：**\n"
				advice += "1. 检查MySQL进程内存占用，优化配置\n"
				advice += "```bash\n"
				advice += "# 优化MySQL内存配置\n"
				advice += "vi /etc/mysql/my.cnf\n"
				advice += "# 添加或修改：\n"
				advice += "innodb_buffer_pool_size = 512M  # 降低缓冲池大小\n"
				advice += "max_connections = 50            # 限制最大连接数\n"
				advice += "```\n\n"

				advice += "2. 重启占用内存高的服务\n"
				advice += "```bash\n"
				advice += "systemctl restart mysql\n"
				advice += "```\n\n"

				advice += "3. 清理系统缓存\n"
				advice += "```bash\n"
				advice += "sync && echo 3 > /proc/sys/vm/drop_caches\n"
				advice += "```\n\n"

				advice += "**长期优化：**\n"
				advice += "- 考虑升级内存到至少4GB\n"
				advice += "- 将部分服务迁移到其他服务器\n"
			} else if usageFloat > 80 {
				advice += "#### 🟡 内存使用率较高 (>80%)\n\n"
				advice += "建议监控内存趋势，考虑优化配置或扩容\n"
			}
		}
	}

	return advice
}

// ==================== 磁盘目录分析 ====================

// AnalyzeDiskDirectoriesByIP 通过IP分析磁盘目录占用
func (s *PerformanceAnalysisService) AnalyzeDiskDirectoriesByIP(ip string) (*agentModel.DiskAnalysisResult, error) {
	// 1. 通过IP查询主机
	host, err := s.hostDao.GetCmdbHostByIP(ip)
	if err != nil {
		return nil, fmt.Errorf("未找到IP为 %s 的服务器: %v", ip, err)
	}

	return s.analyzeDiskDirectories(&host)
}

// AnalyzeDiskDirectoriesByName 通过名称分析磁盘目录占用
func (s *PerformanceAnalysisService) AnalyzeDiskDirectoriesByName(name string) (*agentModel.DiskAnalysisResult, error) {
	// 1. 通过名称查询主机
	host, err := s.hostDao.GetCmdbHostByName(name)
	if err != nil {
		return nil, fmt.Errorf("未找到名称为 %s 的服务器: %v", name, err)
	}

	return s.analyzeDiskDirectories(&host)
}

// AnalyzeDirectoryDeepByIP 通过IP深度分析指定目录
func (s *PerformanceAnalysisService) AnalyzeDirectoryDeepByIP(ip, directoryPath string) (*agentModel.DiskAnalysisResult, error) {
	host, err := s.hostDao.GetCmdbHostByIP(ip)
	if err != nil {
		return nil, fmt.Errorf("未找到IP为 %s 的服务器: %v", ip, err)
	}

	return s.analyzeDirectoryDeep(&host, directoryPath)
}

// analyzeDiskDirectories 分析磁盘目录占用（内部方法）
func (s *PerformanceAnalysisService) analyzeDiskDirectories(host *cmdbModel.CmdbHost) (*agentModel.DiskAnalysisResult, error) {
	result := &agentModel.DiskAnalysisResult{
		HostID:      host.ID,
		HostName:    host.HostName,
		HostIP:      host.SSHIP,
		CollectTime: time.Now(),
		Directories: []agentModel.DirectoryInfo{},
	}

	// 2. 获取SSH认证信息
	sshConfig, err := s.getSSHConfig(host)
	if err != nil {
		return nil, fmt.Errorf("获取SSH认证信息失败: %v", err)
	}

	// 3. 收集磁盘使用情况
	s.collectDiskUsage(sshConfig, result)

	// 4. 收集目录占用情况
	s.collectTopDirectories(sshConfig, result)

	return result, nil
}

func (s *PerformanceAnalysisService) analyzeDirectoryDeep(host *cmdbModel.CmdbHost, directoryPath string) (*agentModel.DiskAnalysisResult, error) {
	cleanPath, err := sanitizeDirectoryPath(directoryPath)
	if err != nil {
		return nil, err
	}

	result := &agentModel.DiskAnalysisResult{
		HostID:       host.ID,
		HostName:     host.HostName,
		HostIP:       host.SSHIP,
		CollectTime:  time.Now(),
		AnalyzedPath: cleanPath,
		Directories:  []agentModel.DirectoryInfo{},
		TopFiles:     []agentModel.DirectoryInfo{},
		CleanupHints: []string{},
	}

	sshConfig, err := s.getSSHConfig(host)
	if err != nil {
		return nil, fmt.Errorf("获取SSH认证信息失败: %v", err)
	}

	existsOutput, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, fmt.Sprintf("test -d %s && echo OK || echo NOT_FOUND", cleanPath))
	if err != nil {
		return nil, fmt.Errorf("检查目录失败: %v", err)
	}
	if strings.TrimSpace(existsOutput) != "OK" {
		return nil, fmt.Errorf("目录 %s 不存在或不是目录", cleanPath)
	}

	s.collectDiskUsage(sshConfig, result)
	s.collectDirectoryChildren(sshConfig, cleanPath, result)
	s.collectDirectoryTopFiles(sshConfig, cleanPath, result)
	result.CleanupHints = s.buildDirectoryCleanupHints(cleanPath)

	return result, nil
}

// collectDiskUsage 收集磁盘使用情况
func (s *PerformanceAnalysisService) collectDiskUsage(sshConfig *util.SSHConfig, result *agentModel.DiskAnalysisResult) {
	// 获取磁盘分区使用情况
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "df -h | grep -vE '^Filesystem|tmpfs|cdrom'")
	if err == nil {
		result.DiskUsage = strings.TrimSpace(output)
	} else {
		log.Log().Warnf("获取磁盘使用情况失败: %v", err)
	}

	// 获取根分区使用率
	output, err = s.sshUtil.ExecuteRemoteCommand(sshConfig, "df -h / | tail -1 | awk '{print $5}' | cut -d'%' -f1")
	if err == nil {
		if usage, err := strconv.ParseFloat(strings.TrimSpace(output), 64); err == nil {
			result.RootUsagePercent = usage
		}
	}
}

// collectTopDirectories 收集占用空间最大的目录
func (s *PerformanceAnalysisService) collectTopDirectories(sshConfig *util.SSHConfig, result *agentModel.DiskAnalysisResult) {
	// 分析根目录下各目录的占用情况
	// 使用 du 命令，但限制深度为1，避免扫描时间过长
	log.Log().Infof("开始分析磁盘目录占用，主机: %s", result.HostIP)

	// 优化的命令：超时30秒，只扫描第一层目录
	command := "timeout 30 du -sh /* 2>/dev/null | sort -hr | head -5"
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, command)

	if err != nil {
		log.Log().Warnf("获取目录占用失败: %v", err)
		return
	}

	// 解析输出
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		// 格式：大小 \t 路径
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			size := parts[0]
			path := parts[1]

			dirInfo := agentModel.DirectoryInfo{
				Path: path,
				Size: size,
			}

			// 解析大小为MB（用于排序和比较）
			sizeInMB := s.parseSizeToMB(size)
			dirInfo.SizeInMB = sizeInMB

			result.Directories = append(result.Directories, dirInfo)
		}
	}

	log.Log().Infof("分析完成，找到 %d 个目录", len(result.Directories))
}

func (s *PerformanceAnalysisService) collectDirectoryChildren(sshConfig *util.SSHConfig, directoryPath string, result *agentModel.DiskAnalysisResult) {
	command := fmt.Sprintf("timeout 30 sh -c 'du -sh %s/* %s/.[!.]* %s/..?* 2>/dev/null | sort -hr | head -5'", directoryPath, directoryPath, directoryPath)
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, command)
	if err != nil {
		log.Log().Warnf("获取目录子项占用失败: %v", err)
		return
	}

	s.parseDirectoryEntries(output, &result.Directories)
}

func (s *PerformanceAnalysisService) collectDirectoryTopFiles(sshConfig *util.SSHConfig, directoryPath string, result *agentModel.DiskAnalysisResult) {
	command := fmt.Sprintf("timeout 45 sh -c 'find %s -xdev -type f -exec du -h {} + 2>/dev/null | sort -hr | head -5'", directoryPath)
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, command)
	if err != nil {
		log.Log().Warnf("获取目录大文件失败: %v", err)
		return
	}

	s.parseDirectoryEntries(output, &result.TopFiles)
}

func (s *PerformanceAnalysisService) parseDirectoryEntries(output string, target *[]agentModel.DirectoryInfo) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		size := parts[0]
		path := parts[1]
		*target = append(*target, agentModel.DirectoryInfo{
			Path:     path,
			Size:     size,
			SizeInMB: s.parseSizeToMB(size),
		})
	}
}

// parseSizeToMB 将大小字符串转换为MB
func (s *PerformanceAnalysisService) parseSizeToMB(size string) float64 {
	size = strings.TrimSpace(size)
	if size == "" {
		return 0
	}

	// 提取数字部分和单位
	var value float64
	var unit string

	_, err := fmt.Sscanf(size, "%f%s", &value, &unit)
	if err != nil {
		return 0
	}

	// 转换为MB
	unit = strings.ToUpper(unit)
	switch unit {
	case "K", "KB":
		return value / 1024
	case "M", "MB":
		return value
	case "G", "GB":
		return value * 1024
	case "T", "TB":
		return value * 1024 * 1024
	default:
		return value / (1024 * 1024) // 假设是字节
	}
}

// GenerateDiskAnalysisReport 生成磁盘目录分析报告
func (s *PerformanceAnalysisService) GenerateDiskAnalysisReport(result *agentModel.DiskAnalysisResult) string {
	var report strings.Builder

	report.WriteString("# 磁盘空间分析报告\n\n")
	report.WriteString(fmt.Sprintf("**主机信息**\n"))
	report.WriteString(fmt.Sprintf("- 主机名称: %s\n", result.HostName))
	report.WriteString(fmt.Sprintf("- 主机IP: %s\n", result.HostIP))
	report.WriteString(fmt.Sprintf("- 分析时间: %s\n\n", result.CollectTime.Format("2006-01-02 15:04:05")))

	// 磁盘使用情况
	report.WriteString("## 磁盘使用情况\n\n")
	if result.DiskUsage != "" {
		report.WriteString("```\n")
		report.WriteString(result.DiskUsage)
		report.WriteString("\n```\n\n")
	}

	if result.RootUsagePercent > 0 {
		report.WriteString(fmt.Sprintf("**根分区使用率**: %.1f%%", result.RootUsagePercent))
		if result.RootUsagePercent > 90 {
			report.WriteString(" 🔴 **告警：磁盘空间严重不足！**\n\n")
		} else if result.RootUsagePercent > 80 {
			report.WriteString(" 🟡 **警告：磁盘空间偏高**\n\n")
		} else {
			report.WriteString(" ✅ **正常**\n\n")
		}
	}

	// 占用空间最大的目录
	if len(result.Directories) > 0 {
		report.WriteString("## 占用空间最大的目录（TOP 5）\n\n")
		report.WriteString("| 排名 | 目录路径 | 占用大小 |\n")
		report.WriteString("|------|----------|----------|\n")

		for i, dir := range result.Directories {
			if i >= 5 {
				break
			}
			report.WriteString(fmt.Sprintf("| %d | %s | **%s** |\n", i+1, dir.Path, dir.Size))
		}
		report.WriteString("\n")

		// 分析建议
		report.WriteString("## 💡 优化建议\n\n")
		if len(result.Directories) > 0 {
			topDir := result.Directories[0]
			report.WriteString(fmt.Sprintf("**最大占用目录**: %s (%s)\n\n", topDir.Path, topDir.Size))

			// 根据目录给出建议
			if strings.Contains(topDir.Path, "/var/log") {
				report.WriteString("**建议**: 日志目录占用过大，可以清理旧日志\n")
				report.WriteString("```bash\n")
				report.WriteString("# 查找并删除7天前的日志文件\n")
				report.WriteString("find /var/log -name '*.log' -mtime +7 -delete\n")
				report.WriteString("```\n\n")
			} else if strings.Contains(topDir.Path, "/tmp") {
				report.WriteString("**建议**: 临时目录占用过大，可以安全清理\n")
				report.WriteString("```bash\n")
				report.WriteString("rm -rf /tmp/*\n")
				report.WriteString("```\n\n")
			} else if strings.Contains(topDir.Path, "/home") || strings.Contains(topDir.Path, "/root") {
				report.WriteString("**建议**: 用户目录占用大，检查是否有不必要的大文件\n")
				report.WriteString("```bash\n")
				report.WriteString(fmt.Sprintf("# 查找%s下最大的10个文件\n", topDir.Path))
				report.WriteString(fmt.Sprintf("find %s -type f -exec du -h {} + 2>/dev/null | sort -rh | head -10\n", topDir.Path))
				report.WriteString("```\n\n")
			} else if strings.Contains(topDir.Path, "/var/lib/docker") {
				report.WriteString("**建议**: Docker目录占用大，清理无用的镜像和容器\n")
				report.WriteString("```bash\n")
				report.WriteString("# 清理未使用的Docker资源\n")
				report.WriteString("docker system prune -a --volumes\n")
				report.WriteString("```\n\n")
			}
		}
	} else {
		report.WriteString("## 占用空间最大的目录\n\n")
		report.WriteString("⚠️ 未能获取目录占用信息（可能是权限问题或超时）\n\n")
	}

	return report.String()
}

// GenerateDirectoryDeepAnalysisReport 生成目录深度分析报告
func (s *PerformanceAnalysisService) GenerateDirectoryDeepAnalysisReport(result *agentModel.DiskAnalysisResult) string {
	var report strings.Builder

	report.WriteString("# 目录深度分析报告\n\n")
	report.WriteString("**分析目录**: ")
	report.WriteString(result.AnalyzedPath)
	report.WriteString("\n\n")

	if len(result.Directories) > 0 {
		report.WriteString("## 子项占用 TOP 5\n\n")
		report.WriteString("| 排名 | 路径 | 占用大小 |\n")
		report.WriteString("|------|------|----------|\n")
		for i, dir := range result.Directories {
			report.WriteString(fmt.Sprintf("| %d | %s | **%s** |\n", i+1, dir.Path, dir.Size))
		}
		report.WriteString("\n")
	}

	if len(result.TopFiles) > 0 {
		report.WriteString("## 最大文件 TOP 5\n\n")
		report.WriteString("| 排名 | 文件路径 | 大小 | 删除建议 |\n")
		report.WriteString("|------|----------|------|----------|\n")
		for i, file := range result.TopFiles {
			report.WriteString(fmt.Sprintf("| %d | %s | **%s** | %s |\n", i+1, file.Path, file.Size, classifyFileDeletionPolicy(file.Path)))
		}
		report.WriteString("\n")
	}

	if len(result.CleanupHints) > 0 {
		report.WriteString("## 可清理建议\n\n")
		for _, hint := range result.CleanupHints {
			report.WriteString(fmt.Sprintf("- %s\n", hint))
		}
		report.WriteString("\n")
	}

	report.WriteString("## 删除规则\n\n")
	report.WriteString("- 可删除: `*.log`、`*.txt`、`*.gz`、`*.zip`\n")
	report.WriteString("- 不可自动删除: 除上述后缀外的其他文件，一律只允许用户手动确认处理\n")

	return report.String()
}

func sanitizeDirectoryPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "`*\"'，。,:：;；()[]{}<>")
	if path == "" || !strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("目录路径必须是绝对路径")
	}
	for _, r := range path {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '/', r == '-', r == '_', r == '.', r == '@':
		default:
			return "", fmt.Errorf("目录路径包含非法字符: %s", path)
		}
	}
	return path, nil
}

func (s *PerformanceAnalysisService) buildDirectoryCleanupHints(path string) []string {
	switch {
	case strings.HasPrefix(path, "/var/log"):
		return []string{
			"/var/log 主要是系统日志、应用日志、journal 日志，优先清理历史日志和已轮转压缩日志",
			"可以先看是否是某个应用日志暴涨，再考虑 logrotate、journalctl 清理，而不是直接删当前正在写入的日志文件",
			"可执行 `journalctl --vacuum-time=7d` 或清理 7 天前的归档日志",
		}
	case strings.HasPrefix(path, "/var/lib/docker"):
		return []string{
			"/var/lib/docker 主要是镜像、容器层、卷和容器日志",
			"优先清理无用镜像、停止的容器、无用 volume，不要手动删除 overlay2 子目录",
			"可执行 `docker system df` 先确认占用，再用 `docker system prune` 做清理",
		}
	case strings.HasPrefix(path, "/tmp"):
		return []string{
			"/tmp 主要是临时文件，通常可清理，但要确认没有正在运行的程序依赖这些文件",
			"优先删除较早的临时文件，而不是直接全量删除",
		}
	case strings.HasPrefix(path, "/var/cache"):
		return []string{
			"/var/cache 主要是包管理、应用缓存，一般比业务数据更适合清理",
			"优先按应用分类清缓存，例如 apt/yum 缓存，而不是直接删除整个 /var",
		}
	case strings.HasPrefix(path, "/var"):
		return []string{
			"/var 常见内容包括日志(/var/log)、缓存(/var/cache)、包数据(/var/lib)、临时文件(/var/tmp)、队列和运行时数据",
			"通常优先排查 /var/log、/var/lib/docker、/var/cache、/var/tmp，这几类最容易膨胀",
			"不要直接删除 /var/lib 下的数据库、容器、服务状态目录，先确认具体子目录用途",
		}
	default:
		return []string{
			"先确认该目录属于日志、缓存、临时文件，还是业务数据目录",
			"能删的一般是日志归档、缓存、临时文件；数据库数据、上传文件、应用运行目录通常不能直接删",
		}
	}
}

func (s *PerformanceAnalysisService) describeDirectoryPurpose(path string) string {
	switch {
	case path == "/var":
		return "/var 主要存放系统运行期经常变化的数据，包括日志、缓存、应用状态、包管理数据、容器数据、临时文件。常见大户是 /var/log、/var/lib/docker、/var/cache、/var/tmp。"
	case strings.HasPrefix(path, "/var/log"):
		return "/var/log 主要存放系统日志、应用日志、审计日志、journal 日志和轮转后的压缩日志。通常可以清理旧日志，但不要直接删当前正在写入的活跃日志。"
	case strings.HasPrefix(path, "/var/lib/docker"):
		return "/var/lib/docker 主要存放 Docker 镜像层、容器可写层、volume 数据和容器日志。这里经常占满磁盘，但不能手动乱删目录。"
	case strings.HasPrefix(path, "/var/cache"):
		return "/var/cache 主要存放包管理器和应用缓存，一般比业务数据更适合清理。"
	case strings.HasPrefix(path, "/tmp"):
		return "/tmp 主要存放临时文件，很多场景下可清理，但需要避开正在使用的文件。"
	default:
		return "这个目录的具体内容取决于部署的软件。先看子目录和大文件，再判断哪些属于日志、缓存、临时文件，哪些属于业务数据。"
	}
}

func classifyFileDeletionPolicy(path string) string {
	lower := strings.ToLower(strings.TrimSpace(path))
	switch {
	case strings.HasSuffix(lower, ".log"),
		strings.HasSuffix(lower, ".txt"),
		strings.HasSuffix(lower, ".gz"),
		strings.HasSuffix(lower, ".zip"):
		return "可删除"
	default:
		return "不可自动删除，需用户手动处理"
	}
}

// QueryRunningServicesByIP 通过IP查询运行中的服务
func (s *PerformanceAnalysisService) QueryRunningServicesByIP(ip string) ([]string, error) {
	// 1. 通过IP查询主机
	host, err := s.hostDao.GetCmdbHostByIP(ip)
	if err != nil {
		return nil, fmt.Errorf("未找到IP为 %s 的服务器: %v", ip, err)
	}

	// 2. 获取SSH认证信息
	sshConfig, err := s.getSSHConfig(&host)
	if err != nil {
		return nil, fmt.Errorf("获取SSH认证信息失败: %v", err)
	}

	// 3. 执行 netstat 或 ss 命令获取监听端口的服务
	cmd := `netstat -nutpl 2>/dev/null || ss -tulpn 2>/dev/null`
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, cmd)
	if err != nil {
		return nil, fmt.Errorf("执行命令失败: %v", err)
	}

	// 4. 解析输出，提取进程名
	services := make(map[string]bool)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Active") || strings.HasPrefix(line, "Proto") {
			continue
		}

		// 提取最后一列的 PID/Program name
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		// 最后一列格式: "1115/mysqld" 或 "931/nginx: master p"
		lastField := fields[len(fields)-1]
		if strings.Contains(lastField, "/") {
			parts := strings.SplitN(lastField, "/", 2)
			if len(parts) == 2 {
				programName := parts[1]
				// 去掉可能的参数和后缀
				programName = strings.TrimSpace(programName)
				if idx := strings.Index(programName, ":"); idx > 0 {
					programName = programName[:idx]
				}
				if programName != "" && programName != "-" {
					services[programName] = true
				}
			}
		}
	}

	// 5. 转换为列表并返回
	result := make([]string, 0, len(services))
	for service := range services {
		result = append(result, service)
	}

	return result, nil
}
