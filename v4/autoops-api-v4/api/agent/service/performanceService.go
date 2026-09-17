// 服务器性能分析服务
package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/agent/model"
	cmdbModel "dodevops-api/api/cmdb/model"
	"dodevops-api/common/util"
	"dodevops-api/pkg/log"
)

// PerformanceService 性能分析服务接口
type PerformanceService interface {
	AnalyzeServerPressure(ctx context.Context, host *cmdbModel.CmdbHost) (*model.PressureAnalysisResult, error)
	AnalyzeMemoryUsage(ctx context.Context, host *cmdbModel.CmdbHost) (*model.MemoryAnalysisResult, error)
}

// PerformanceServiceImpl 性能分析服务实现
type PerformanceServiceImpl struct {
	sshUtil *util.SSHUtil
}

// NewPerformanceService 创建性能分析服务实例
func NewPerformanceService() PerformanceService {
	return &PerformanceServiceImpl{
		sshUtil: util.NewSSHUtil(),
	}
}

// AnalyzeServerPressure 分析服务器压力
func (s *PerformanceServiceImpl) AnalyzeServerPressure(ctx context.Context, host *cmdbModel.CmdbHost) (*model.PressureAnalysisResult, error) {
	result := &model.PressureAnalysisResult{
		HostID:      host.ID,
		HostName:    host.HostName,
		HostIP:      host.SSHIP,
		CollectTime: time.Now(),
	}

	// 获取SSH配置
	sshConfig, err := s.getSSHConfig(host)
	if err != nil {
		return nil, fmt.Errorf("获取SSH配置失败: %v", err)
	}

	// 1. 收集CPU信息
	cpuInfo, err := s.collectCPUInfo(sshConfig)
	if err != nil {
		log.Log().Warnf("收集CPU信息失败: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("CPU信息收集失败: %v", err))
	} else {
		result.CPUInfo = cpuInfo
	}

	// 2. 收集内存信息
	memInfo, err := s.collectMemoryInfo(sshConfig)
	if err != nil {
		log.Log().Warnf("收集内存信息失败: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("内存信息收集失败: %v", err))
	} else {
		result.MemoryInfo = memInfo
	}

	// 3. 收集磁盘信息
	diskInfo, err := s.collectDiskInfo(sshConfig)
	if err != nil {
		log.Log().Warnf("收集磁盘信息失败: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("磁盘信息收集失败: %v", err))
	} else {
		result.DiskInfo = diskInfo
	}

	// 4. 收集负载信息
	loadInfo, err := s.collectLoadInfo(sshConfig)
	if err != nil {
		log.Log().Warnf("收集负载信息失败: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("负载信息收集失败: %v", err))
	} else {
		result.LoadInfo = loadInfo
	}

	// 5. 收集TOP进程
	topProcesses, err := s.collectTopProcesses(sshConfig)
	if err != nil {
		log.Log().Warnf("收集TOP进程失败: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("进程信息收集失败: %v", err))
	} else {
		result.TopProcesses = topProcesses
	}

	// 6. 收集网络连接
	networkInfo, err := s.collectNetworkInfo(sshConfig)
	if err != nil {
		log.Log().Warnf("收集网络信息失败: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("网络信息收集失败: %v", err))
	} else {
		result.NetworkInfo = networkInfo
	}

	// 7. 评估压力等级
	result.PressureLevel = s.assessPressureLevel(result)

	return result, nil
}

// AnalyzeMemoryUsage 分析内存使用
func (s *PerformanceServiceImpl) AnalyzeMemoryUsage(ctx context.Context, host *cmdbModel.CmdbHost) (*model.MemoryAnalysisResult, error) {
	result := &model.MemoryAnalysisResult{
		HostID:      host.ID,
		HostName:    host.HostName,
		HostIP:      host.SSHIP,
		CollectTime: time.Now(),
	}

	// 获取SSH配置
	sshConfig, err := s.getSSHConfig(host)
	if err != nil {
		return nil, fmt.Errorf("获取SSH配置失败: %v", err)
	}

	// 1. 收集详细内存信息
	memInfo, err := s.collectDetailedMemoryInfo(sshConfig)
	if err != nil {
		return nil, fmt.Errorf("收集内存信息失败: %v", err)
	}
	result.MemoryInfo = memInfo

	// 2. 收集内存占用最高的进程
	topMemProcesses, err := s.collectTopMemoryProcesses(sshConfig, 10)
	if err != nil {
		log.Log().Warnf("收集内存占用进程失败: %v", err)
	} else {
		result.TopProcesses = topMemProcesses
	}

	// 3. 分析内存使用模式
	result.MemoryPattern = s.analyzeMemoryPattern(memInfo, topMemProcesses)

	// 4. 识别潜在问题
	result.Issues = s.identifyMemoryIssues(memInfo, topMemProcesses)

	return result, nil
}

// getSSHConfig 获取SSH配置
func (s *PerformanceServiceImpl) getSSHConfig(host *cmdbModel.CmdbHost) (*util.SSHConfig, error) {
	// TODO: 需要从数据库获取完整的SSH认证信息
	// 暂时返回基本配置
	config := &util.SSHConfig{
		IP:       host.SSHIP,
		Port:     host.SSHPort,
		Username: host.SSHName,
		Type:     1, // 默认密码认证，实际应从数据库读取
	}

	return config, nil
}

// collectCPUInfo 收集CPU信息
func (s *PerformanceServiceImpl) collectCPUInfo(sshConfig *util.SSHConfig) (map[string]interface{}, error) {
	// 执行命令获取CPU使用率
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "top -bn1 | grep 'Cpu(s)' | awk '{print $2}'")
	if err != nil {
		return nil, err
	}

	cpuUsage := strings.TrimSpace(strings.Replace(output, "%us,", "", -1))
	usage, _ := strconv.ParseFloat(cpuUsage, 64)

	return map[string]interface{}{
		"usage":   usage,
		"raw":     output,
		"command": "top -bn1 | grep 'Cpu(s)'",
	}, nil
}

// collectMemoryInfo 收集内存信息
func (s *PerformanceServiceImpl) collectMemoryInfo(sshConfig *util.SSHConfig) (map[string]interface{}, error) {
	// 使用free命令获取内存信息
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "free -m")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("内存信息格式错误")
	}

	// 解析内存行
	memLine := strings.Fields(lines[1])
	if len(memLine) < 7 {
		return nil, fmt.Errorf("内存信息解析失败")
	}

	total, _ := strconv.ParseFloat(memLine[1], 64)
	used, _ := strconv.ParseFloat(memLine[2], 64)
	free, _ := strconv.ParseFloat(memLine[3], 64)
	available, _ := strconv.ParseFloat(memLine[6], 64)

	usagePercent := 0.0
	if total > 0 {
		usagePercent = (used / total) * 100
	}

	return map[string]interface{}{
		"total":     total,
		"used":      used,
		"free":      free,
		"available": available,
		"usage":     usagePercent,
		"raw":       output,
	}, nil
}

// collectDetailedMemoryInfo 收集详细内存信息
func (s *PerformanceServiceImpl) collectDetailedMemoryInfo(sshConfig *util.SSHConfig) (map[string]interface{}, error) {
	return s.collectMemoryInfo(sshConfig)
}

// collectDiskInfo 收集磁盘信息
func (s *PerformanceServiceImpl) collectDiskInfo(sshConfig *util.SSHConfig) (map[string]interface{}, error) {
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "df -h / | tail -n 1 | awk '{print $5}'")
	if err != nil {
		return nil, err
	}

	usageStr := strings.TrimSpace(strings.Replace(output, "%", "", -1))
	usage, _ := strconv.ParseFloat(usageStr, 64)

	return map[string]interface{}{
		"usage": usage,
		"raw":   output,
	}, nil
}

// collectLoadInfo 收集负载信息
func (s *PerformanceServiceImpl) collectLoadInfo(sshConfig *util.SSHConfig) (map[string]interface{}, error) {
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "uptime")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"raw": output,
	}, nil
}

// collectTopProcesses 收集TOP进程
func (s *PerformanceServiceImpl) collectTopProcesses(sshConfig *util.SSHConfig) ([]map[string]interface{}, error) {
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "ps aux --sort=-%cpu | head -n 11")
	if err != nil {
		return nil, err
	}

	processes := []map[string]interface{}{}
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if i == 0 || line == "" {
			continue // 跳过表头
		}
		fields := strings.Fields(line)
		if len(fields) >= 11 {
			cpu, _ := strconv.ParseFloat(fields[2], 64)
			mem, _ := strconv.ParseFloat(fields[3], 64)
			processes = append(processes, map[string]interface{}{
				"user":    fields[0],
				"pid":     fields[1],
				"cpu":     cpu,
				"mem":     mem,
				"command": strings.Join(fields[10:], " "),
			})
		}
	}

	return processes, nil
}

// collectTopMemoryProcesses 收集内存占用最高的进程
func (s *PerformanceServiceImpl) collectTopMemoryProcesses(sshConfig *util.SSHConfig, limit int) ([]map[string]interface{}, error) {
	cmd := fmt.Sprintf("ps aux --sort=-%cmem | head -n %d", '%', limit+1)
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, cmd)
	if err != nil {
		return nil, err
	}

	processes := []map[string]interface{}{}
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 11 {
			cpu, _ := strconv.ParseFloat(fields[2], 64)
			mem, _ := strconv.ParseFloat(fields[3], 64)
			processes = append(processes, map[string]interface{}{
				"user":    fields[0],
				"pid":     fields[1],
				"cpu":     cpu,
				"mem":     mem,
				"vsz":     fields[4],
				"rss":     fields[5],
				"command": strings.Join(fields[10:], " "),
			})
		}
	}

	return processes, nil
}

// collectNetworkInfo 收集网络信息
func (s *PerformanceServiceImpl) collectNetworkInfo(sshConfig *util.SSHConfig) (map[string]interface{}, error) {
	output, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, "netstat -an | wc -l")
	if err != nil {
		return nil, err
	}

	connections, _ := strconv.Atoi(strings.TrimSpace(output))

	return map[string]interface{}{
		"connections": connections,
		"raw":         output,
	}, nil
}

// assessPressureLevel 评估压力等级
func (s *PerformanceServiceImpl) assessPressureLevel(result *model.PressureAnalysisResult) string {
	score := 0

	// CPU评分
	if cpuInfo, ok := result.CPUInfo["usage"].(float64); ok {
		if cpuInfo > 80 {
			score += 3
		} else if cpuInfo > 60 {
			score += 2
		} else if cpuInfo > 40 {
			score += 1
		}
	}

	// 内存评分
	if memInfo, ok := result.MemoryInfo["usage"].(float64); ok {
		if memInfo > 90 {
			score += 3
		} else if memInfo > 75 {
			score += 2
		} else if memInfo > 60 {
			score += 1
		}
	}

	// 磁盘评分
	if diskInfo, ok := result.DiskInfo["usage"].(float64); ok {
		if diskInfo > 90 {
			score += 2
		} else if diskInfo > 80 {
			score += 1
		}
	}

	if score >= 6 {
		return "高压力"
	} else if score >= 3 {
		return "中等压力"
	} else {
		return "正常"
	}
}

// analyzeMemoryPattern 分析内存使用模式
func (s *PerformanceServiceImpl) analyzeMemoryPattern(memInfo map[string]interface{}, processes []map[string]interface{}) string {
	usage, _ := memInfo["usage"].(float64)

	if usage > 90 {
		return "内存严重不足，可能存在内存泄漏或配置不足"
	} else if usage > 75 {
		return "内存使用偏高，需要关注"
	} else if usage > 50 {
		return "内存使用正常偏高"
	}
	return "内存使用正常"
}

// identifyMemoryIssues 识别内存问题
func (s *PerformanceServiceImpl) identifyMemoryIssues(memInfo map[string]interface{}, processes []map[string]interface{}) []string {
	issues := []string{}

	usage, _ := memInfo["usage"].(float64)
	if usage > 90 {
		issues = append(issues, "内存使用率超过90%，系统可能出现性能问题")
	}

	// 检查是否有进程占用过高
	for _, proc := range processes {
		if mem, ok := proc["mem"].(float64); ok && mem > 20 {
			command, _ := proc["command"].(string)
			issues = append(issues, fmt.Sprintf("进程 '%s' 占用内存过高: %.1f%%", command, mem))
		}
	}

	if len(issues) == 0 {
		issues = append(issues, "未发现明显内存问题")
	}

	return issues
}
