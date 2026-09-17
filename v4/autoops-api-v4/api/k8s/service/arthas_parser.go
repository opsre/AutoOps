package service

import (
	"dodevops-api/api/k8s/model"
	"regexp"
	"strconv"
	"strings"
)

// stripANSI 移除 ANSI 转义序列
func stripANSI(text string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(text, "")
}

// extractArthasOutput 提取 Arthas 命令输出（去除 banner 和前缀）
func extractArthasOutput(rawOutput string) string {
	cleaned := stripANSI(rawOutput)
	cleaned = strings.ReplaceAll(cleaned, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", "\n")
	
	// 找到 "| plaintext" 之后的内容
	if idx := strings.Index(cleaned, "| plaintext"); idx >= 0 {
		cleaned = cleaned[idx+len("| plaintext"):]
	}
	
	return strings.TrimSpace(cleaned)
}

// parseDashboard 解析 dashboard 输出
func parseDashboard(rawOutput string) *model.DashboardData {
	output := extractArthasOutput(rawOutput)
	data := &model.DashboardData{
		Threads: model.ThreadStats{},
		Memory:  model.MemoryStats{
			Heap:    model.MemoryUsage{Max: -1},
			NonHeap: model.MemoryUsage{Max: -1},
		},
		Runtime: model.ArthasRuntimeInfo{},
	}
	
	lines := strings.Split(output, "\n")
	
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// 解析内存部分 - 查找 "Memory" 开头的行
		if strings.HasPrefix(line, "Memory") {
			// 下一行应该是 heap 数据
			for j := i + 1; j < len(lines) && j < i+10; j++ {
				dataLine := strings.TrimSpace(lines[j])
				if dataLine == "" {
					continue
				}
				
				// 使用正则匹配 heap 行（处理多个空格）
				// 格式: heap             129M  396M  595M       gc.copy.count       129
				if strings.HasPrefix(dataLine, "heap") {
					re := regexp.MustCompile(`^heap\s+(\S+)\s+(\S+)\s+(\S+)`)
					matches := re.FindStringSubmatch(dataLine)
					if len(matches) >= 4 {
						data.Memory.Heap.Used = parseMemoryValue(matches[1])
						data.Memory.Heap.Committed = parseMemoryValue(matches[2])
						data.Memory.Heap.Max = parseMemoryValue(matches[3])
					}
					break
				}
			}
		}
		
		// 解析运行时信息
		if strings.HasPrefix(line, "Runtime") {
			// 后续行包含 os.name, os.version, java.version
			for j := i + 1; j < len(lines) && j < i+10; j++ {
				infoLine := strings.TrimSpace(lines[j])
				if infoLine == "" || strings.HasPrefix(infoLine, "Process ends") {
					break
				}
				
				if strings.HasPrefix(infoLine, "os.name") {
					parts := strings.Fields(infoLine)
					if len(parts) >= 2 {
						data.Runtime.OSName = parts[len(parts)-1]
					}
				} else if strings.HasPrefix(infoLine, "os.version") {
					parts := strings.Fields(infoLine)
					if len(parts) >= 2 {
						data.Runtime.OSVersion = parts[len(parts)-1]
					}
				} else if strings.HasPrefix(infoLine, "java.version") {
					parts := strings.Fields(infoLine)
					if len(parts) >= 2 {
						data.Runtime.JavaVersion = parts[len(parts)-1]
					}
				}
			}
		}
	}
	
	// 从原始输出中提取线程统计（更准确的方式）
	// 查找 JVM 输出中的 THREAD 部分
	if strings.Contains(rawOutput, "COUNT") && strings.Contains(rawOutput, "DAEMON-COUNT") {
		re := regexp.MustCompile(`(?m)^\s*COUNT\s+(\d+)`)
		matches := re.FindStringSubmatch(rawOutput)
		if len(matches) >= 2 {
			data.Threads.Live, _ = strconv.Atoi(matches[1])
		}
		
		re = regexp.MustCompile(`(?m)^\s*DAEMON-COUNT\s+(\d+)`)
		matches = re.FindStringSubmatch(rawOutput)
		if len(matches) >= 2 {
			data.Threads.Daemon, _ = strconv.Atoi(matches[1])
		}
		
		re = regexp.MustCompile(`(?m)^\s*PEAK-COUNT\s+(\d+)`)
		matches = re.FindStringSubmatch(rawOutput)
		if len(matches) >= 2 {
			data.Threads.Peak, _ = strconv.Atoi(matches[1])
		}
	} else {
		// 从 dashboard 输出中统计线程（计数线程表格行）
		inThreadTable := false
		threadCount := 0
		
		for _, line := range lines {
			line = strings.TrimSpace(line)
			
			// 线程表格开始
			if strings.HasPrefix(line, "ID") && strings.Contains(line, "NAME") && strings.Contains(line, "STATE") {
				inThreadTable = true
				continue
			}
			
			// 线程表格结束
			if inThreadTable && (strings.HasPrefix(line, "Memory") || strings.HasPrefix(line, "Runtime")) {
				break
			}
			
			// 统计线程行
			if inThreadTable && line != "" {
				fields := strings.Fields(line)
				// 有效的线程行至少有 ID, NAME, GROUP, PRIORITY, STATE 这5个字段
				if len(fields) >= 5 {
					threadCount++
				}
			}
		}
		
		if threadCount > 0 {
			data.Threads.Live = threadCount
		}
	}
	
	return data
}

// parseThreadList 解析 thread 列表输出
func parseThreadList(rawOutput string) *model.ThreadListData {
	output := extractArthasOutput(rawOutput)
	data := &model.ThreadListData{
		Threads: []model.ThreadInfo{},
	}
	
	lines := strings.Split(output, "\n")
	
	// 解析总数（第一行 "Threads Total: 143, NEW: 0..."）
	for _, line := range lines {
		if strings.HasPrefix(line, "Threads Total:") {
			re := regexp.MustCompile(`Threads Total:\s*(\d+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 2 {
				data.Total, _ = strconv.Atoi(matches[1])
			}
			break
		}
	}
	
	// 解析线程表格
	inTable := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// 表头行
		if strings.HasPrefix(line, "ID") && strings.Contains(line, "NAME") && strings.Contains(line, "STATE") {
			inTable = true
			continue
		}
		
		// 数据行
		if inTable && !strings.HasPrefix(line, "[arthas@") && !strings.HasPrefix(line, "Process ends") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				thread := model.ThreadInfo{
					ThreadID:   fields[0],
					ThreadName: fields[1],
					State:      fields[4],
				}
				if len(fields) >= 6 {
					thread.CPU = fields[5]
				}
				if len(fields) >= 8 {
					thread.Time = fields[7]
				}
				data.Threads = append(data.Threads, thread)
			}
		}
	}
	
	return data
}

// parseJvmInfo 解析 jvm 输出
func parseJvmInfo(rawOutput string) *model.JvmData {
	data := &model.JvmData{}
	
	// 清理输出
	cleaned := stripANSI(rawOutput)
	cleaned = strings.ReplaceAll(cleaned, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", " ")
	
	// 使用正则提取关键信息（多行模式）
	extractField := func(pattern string) string {
		re := regexp.MustCompile(`(?m)^\s*` + pattern + `\s+(.+)$`)
		matches := re.FindStringSubmatch(cleaned)
		if len(matches) >= 2 {
			return strings.TrimSpace(matches[1])
		}
		return ""
	}
	
	data.OSName = extractField("OS")
	data.OSVersion = extractField("VERSION")
	data.JavaVersion = extractField("SPEC-VERSION")
	if data.JavaVersion == "" {
		data.JavaVersion = extractField("VM-VERSION")
	}
	data.JavaVendor = extractField("SPEC-VENDOR")
	data.VMName = extractField("VM-NAME")
	data.VMVersion = extractField("VM-VERSION")
	data.ClassPath = extractField("CLASS-PATH")
	
	// 解析处理器数量
	procStr := extractField("PROCESSORS-COUNT")
	if procStr != "" {
		data.ProcessorCount, _ = strconv.Atoi(procStr)
	}
	
	return data
}

// parseMemoryInfo 解析 memory 输出
func parseMemoryInfo(rawOutput string) *model.MemoryData {
	output := extractArthasOutput(rawOutput)
	data := &model.MemoryData{
		HeapMemoryUsage:    model.MemoryUsage{Max: -1},
		NonHeapMemoryUsage: model.MemoryUsage{Max: -1},
	}
	
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// 跳过表头
		if strings.HasPrefix(line, "Memory") {
			continue
		}
		
		// 使用正则匹配（处理多个空格）
		if strings.HasPrefix(line, "heap") {
			re := regexp.MustCompile(`^heap\s+(\S+)\s+(\S+)\s+(\S+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 4 {
				data.HeapMemoryUsage.Used = parseMemoryValue(matches[1])
				data.HeapMemoryUsage.Committed = parseMemoryValue(matches[2])
				data.HeapMemoryUsage.Max = parseMemoryValue(matches[3])
			}
		} else if strings.HasPrefix(line, "nonheap") {
			re := regexp.MustCompile(`^nonheap\s+(\S+)\s+(\S+)\s+(\S+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 4 {
				data.NonHeapMemoryUsage.Used = parseMemoryValue(matches[1])
				data.NonHeapMemoryUsage.Committed = parseMemoryValue(matches[2])
				data.NonHeapMemoryUsage.Max = parseMemoryValue(matches[3])
			}
		}
	}
	
	return data
}

// parseSysprop 解析 sysprop 输出为 key-value map
func parseSysprop(rawOutput string) map[string]string {
	output := extractArthasOutput(rawOutput)
	result := make(map[string]string)
	
	lines := strings.Split(output, "\n")
	inTable := false
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// 表头行
		if strings.HasPrefix(line, "KEY") && strings.Contains(line, "VALUE") {
			inTable = true
			continue
		}
		
		// 分隔线
		if strings.HasPrefix(line, "---") {
			continue
		}
		
		// 数据行
		if inTable && !strings.HasPrefix(line, "[arthas@") && !strings.HasPrefix(line, "Process ends") {
			// 按空格分割，KEY 和 VALUE 之间可能有多个空格
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				key := parts[0]
				value := strings.Join(parts[1:], " ")
				result[key] = value
			}
		}
	}
	
	return result
}

// parseMemoryValue 解析内存值（如 "129M" -> bytes）
func parseMemoryValue(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" || s == "-1" {
		return -1
	}
	
	multiplier := int64(1)
	if strings.HasSuffix(s, "K") {
		multiplier = 1024
		s = strings.TrimSuffix(s, "K")
	} else if strings.HasSuffix(s, "M") {
		multiplier = 1024 * 1024
		s = strings.TrimSuffix(s, "M")
	} else if strings.HasSuffix(s, "G") {
		multiplier = 1024 * 1024 * 1024
		s = strings.TrimSuffix(s, "G")
	}
	
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	
	return int64(val * float64(multiplier))
}
