package model

// ArthasBaseQuery Arthas通用查询参数（绑定 query string）
type ArthasBaseQuery struct {
	Namespace string `form:"namespace" binding:"required"`
	Pod       string `form:"pod" binding:"required"`
	Container string `form:"container" binding:"required"`
	ProcessID string `form:"processId"`
}

// ArthasCheckQuery 检查和安装不需要 processId
type ArthasCheckQuery struct {
	Namespace string `form:"namespace" binding:"required"`
	Pod       string `form:"pod" binding:"required"`
	Container string `form:"container" binding:"required"`
}

// ArthasInstallRequest 安装Arthas请求
type ArthasInstallRequest struct {
	Namespace string `json:"namespace" binding:"required"`
	Pod       string `json:"pod" binding:"required"`
	Container string `json:"container" binding:"required"`
}

// ArthasCommandRequest 自定义命令请求
type ArthasCommandRequest struct {
	Namespace string `json:"namespace" binding:"required"`
	Pod       string `json:"pod" binding:"required"`
	Container string `json:"container" binding:"required"`
	ProcessID string `json:"processId" binding:"required"`
	Command   string `json:"command" binding:"required"`
}

// ProcessInfo Java进程信息
type ProcessInfo struct {
	ProcessID string `json:"processId"`
	MainClass string `json:"mainClass"`
}

// ArthasCheckResult Arthas安装检查结果
type ArthasCheckResult struct {
	Installed bool   `json:"installed"`
	JavaPath  string `json:"javaPath,omitempty"`
	Message   string `json:"message"`
}

// ThreadInfo 线程信息
type ThreadInfo struct {
	ThreadID   string `json:"threadId"`
	ThreadName string `json:"threadName"`
	State      string `json:"state"`
	CPU        string `json:"cpu,omitempty"`
	Time       string `json:"time,omitempty"`
}

// DashboardData 仪表板数据
type DashboardData struct {
	Threads ThreadStats       `json:"threads"`
	Memory  MemoryStats       `json:"memory"`
	Runtime ArthasRuntimeInfo `json:"runtime,omitempty"`
}

// ThreadStats 线程统计
type ThreadStats struct {
	Live   int `json:"live"`
	Peak   int `json:"peak"`
	Daemon int `json:"daemon"`
}

// MemoryStats 内存统计
type MemoryStats struct {
	Heap    MemoryUsage `json:"heap"`
	NonHeap MemoryUsage `json:"nonHeap,omitempty"`
}

// MemoryUsage 内存使用情况
type MemoryUsage struct {
	Init      int64   `json:"init"`
	Used      int64   `json:"used"`
	Committed int64   `json:"committed"`
	Max       int64   `json:"max"`
	Usage     float64 `json:"usage,omitempty"`
}

// ArthasRuntimeInfo Arthas 运行时信息
type ArthasRuntimeInfo struct {
	OSName      string `json:"osName,omitempty"`
	OSVersion   string `json:"osVersion,omitempty"`
	JavaVersion string `json:"javaVersion,omitempty"`
}

// ThreadListData 线程列表数据
type ThreadListData struct {
	Total   int          `json:"total"`
	Threads []ThreadInfo `json:"threads"`
}

// JvmData JVM详细信息
type JvmData struct {
	OSName         string `json:"osName"`
	OSVersion      string `json:"osVersion"`
	ProcessorCount int    `json:"processorCount"`
	JavaVersion    string `json:"javaVersion"`
	JavaVendor     string `json:"javaVendor"`
	ClassPath      string `json:"classPath"`
	VMName         string `json:"vmName,omitempty"`
	VMVersion      string `json:"vmVersion,omitempty"`
}

// MemoryData 内存详细信息
type MemoryData struct {
	HeapMemoryUsage    MemoryUsage `json:"heapMemoryUsage"`
	NonHeapMemoryUsage MemoryUsage `json:"nonHeapMemoryUsage"`
}

// JvmInfo JVM信息
type JvmInfo struct {
	ClassPath string `json:"classPath"`
	BootPath  string `json:"bootPath"`
}

// FlameGraphResult 火焰图结果
type FlameGraphResult struct {
	Data    *FlameGraphData `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}

// FlameGraphTaskStatus 火焰图任务状态
type FlameGraphTaskStatus string

const (
	FlameGraphTaskPending  FlameGraphTaskStatus = "pending"   // 等待中
	FlameGraphTaskRunning  FlameGraphTaskStatus = "running"   // 采样中
	FlameGraphTaskDone     FlameGraphTaskStatus = "done"      // 完成
	FlameGraphTaskFailed   FlameGraphTaskStatus = "failed"    // 失败
)

// FlameGraphTask 火焰图异步任务
type FlameGraphTask struct {
	TaskID    string               `json:"taskId"`
	Status    FlameGraphTaskStatus `json:"status"`
	Duration  int                  `json:"duration"`
	StartedAt int64                `json:"startedAt"`
	FinishAt  int64                `json:"finishAt,omitempty"`
	Message   string               `json:"message,omitempty"`
	Result    *FlameGraphResult    `json:"result,omitempty"`
}

// FlameGraphData 火焰图数据（用于前端渲染）
type FlameGraphData struct {
	Name     string            `json:"name"`           // 显示名称（最后一段函数名）
	FullName string            `json:"fullName"`        // 完整调用链（用于 tooltip）
	Value    int64             `json:"value"`
	Children []*FlameGraphData `json:"children,omitempty"`
}

// FlameGraphNode 火焰图节点（内部构建用）
type FlameGraphNode struct {
	Name     string
	Value    int64
	Children map[string]*FlameGraphNode
}
