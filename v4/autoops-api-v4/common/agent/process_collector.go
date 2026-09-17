package agent

import (
	"fmt"
	"sort"

	"github.com/shirou/gopsutil/v3/process"
)

func monitorProcesses(config *Config) {
	processes, _ := process.Processes()

	// 收集所有进程的CPU和内存使用率
	type procStats struct {
		name string
		pid  int32
		cpu  float64
		mem  float32
	}
	var allProcs []procStats

	for _, p := range processes {
		name, _ := p.Name()
		pid := p.Pid
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()

		allProcs = append(allProcs, procStats{
			name: name,
			pid:  pid,
			cpu:  cpuPercent,
			mem:  memPercent,
		})
	}

	// 按CPU使用率排序
	sort.Slice(allProcs, func(i, j int) bool {
		return allProcs[i].cpu > allProcs[j].cpu
	})

	// 清理旧的CPU指标
	processCPU.Reset()

	// 记录CPU使用率最高的5个进程
	for i := 0; i < 5 && i < len(allProcs); i++ {
		p := allProcs[i]
		processCPU.WithLabelValues(hostname, p.name, fmt.Sprint(p.pid)).Set(p.cpu)
	}

	// 按内存使用率排序
	sort.Slice(allProcs, func(i, j int) bool {
		return allProcs[i].mem > allProcs[j].mem
	})

	// 清理旧的内存指标
	processMem.Reset()

	// 记录内存使用率最高的5个进程
	for i := 0; i < 5 && i < len(allProcs); i++ {
		p := allProcs[i]
		processMem.WithLabelValues(hostname, p.name, fmt.Sprint(p.pid)).Set(float64(p.mem))
	}
}
