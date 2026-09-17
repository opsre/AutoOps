// 录制引擎核心实现
// 负责管理会话录制的生命周期和事件记录
package recording

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

// RecordingMetadata 录制元数据
type RecordingMetadata struct {
	SessionID      string // 会话唯一标识
	UserID         uint   // 用户ID
	Username       string // 用户名
	HostID         uint   // 主机ID
	HostName       string // 主机名称
	HostIP         string // 主机IP
	SSHUser        string // SSH登录用户
	TerminalWidth  int    // 终端宽度
	TerminalHeight int    // 终端高度
	ClientIP       string // 客户端IP
	UserAgent      string // 用户代理
}

// RecordingInfo 录制完成信息
type RecordingInfo struct {
	SessionID       string  // 会话ID
	FilePath        string  // 文件路径
	FileSize        int64   // 文件大小
	Duration        int     // 持续时长(秒)
	InputCount      int     // 输入事件数量
	OutputCount     int     // 输出事件数量
	ResizeCount     int     // 窗口调整次数
	CommandCount    int     // 命令数量(估算)
	RiskLevel       int     // 风险等级(由命令审计填充)
	HasSensitiveCmd bool    // 是否包含敏感命令(由命令审计填充)
}

// Recorder 会话录制器接口
type Recorder interface {
	// Start 开始录制
	Start(metadata RecordingMetadata) error

	// RecordInput 记录用户输入
	RecordInput(data []byte) error

	// RecordOutput 记录SSH输出
	RecordOutput(data []byte) error

	// RecordResize 记录窗口调整
	RecordResize(width, height int) error

	// Stop 停止录制并保存
	Stop() (*RecordingInfo, error)

	// GetSessionID 获取会话ID
	GetSessionID() string

	// GetMetadata 获取录制元数据
	GetMetadata() RecordingMetadata
}

// AsciicastRecorder asciicast格式录制器实现
type AsciicastRecorder struct {
	metadata  RecordingMetadata // 录制元数据
	storage   Storage           // 存储后端
	startTime time.Time         // 录制开始时间

	// 临时文件写入
	tempFile   *os.File       // 临时文件
	tempWriter *bufio.Writer  // 缓冲写入器
	mu         sync.Mutex     // 并发保护

	// 统计信息
	inputCount   int // 输入事件计数
	outputCount  int // 输出事件计数
	resizeCount  int // 窗口调整计数
	commandCount int // 命令计数(估算)

	// 命令解析
	cmdParser *CommandParser // 命令解析器

	// 状态标记
	started bool // 是否已开始
	stopped bool // 是否已停止
}

// NewAsciicastRecorder 创建asciicast录制器
func NewAsciicastRecorder(storage Storage, metadata RecordingMetadata) *AsciicastRecorder {
	return &AsciicastRecorder{
		metadata:  metadata,
		storage:   storage,
		cmdParser: NewCommandParser(),
	}
}

// Start 开始录制(实现Recorder接口)
func (r *AsciicastRecorder) Start(metadata RecordingMetadata) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.started {
		return fmt.Errorf("recorder already started")
	}

	r.metadata = metadata
	r.startTime = time.Now()

	// 创建临时文件
	tempFile, err := os.CreateTemp("", fmt.Sprintf("recording_%s_*.cast", metadata.SessionID))
	if err != nil {
		return fmt.Errorf("create temp file failed: %w", err)
	}
	r.tempFile = tempFile
	r.tempWriter = bufio.NewWriter(tempFile)

	// 写入asciicast头部
	header := NewAsciicastHeader(
		metadata.TerminalWidth,
		metadata.TerminalHeight,
		fmt.Sprintf("SSH %s@%s", metadata.SSHUser, metadata.HostIP),
	)

	headerData, err := header.MarshalHeader()
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return fmt.Errorf("marshal header failed: %w", err)
	}

	if _, err := r.tempWriter.Write(headerData); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return fmt.Errorf("write header failed: %w", err)
	}

	// 刷新缓冲区
	if err := r.tempWriter.Flush(); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return fmt.Errorf("flush header failed: %w", err)
	}

	r.started = true
	return nil
}

// RecordInput 记录用户输入(实现Recorder接口)
func (r *AsciicastRecorder) RecordInput(data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.started || r.stopped {
		return fmt.Errorf("recorder not active")
	}

	// 计算相对时间戳
	elapsed := time.Since(r.startTime).Seconds()

	// 创建输入事件
	event := NewInputEvent(elapsed, string(data))
	eventData, err := event.MarshalEvent()
	if err != nil {
		return fmt.Errorf("marshal input event failed: %w", err)
	}

	// 写入事件
	if _, err := r.tempWriter.Write(eventData); err != nil {
		return fmt.Errorf("write input event failed: %w", err)
	}

	r.inputCount++

	// 从输入流解析命令
	if hasCommand, cmd := r.cmdParser.ParseInput(data, elapsed); hasCommand && cmd != nil {
		r.commandCount++
	}

	return nil
}

// RecordOutput 记录SSH输出(实现Recorder接口)
func (r *AsciicastRecorder) RecordOutput(data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.started || r.stopped {
		return fmt.Errorf("recorder not active")
	}

	// 计算相对时间戳
	elapsed := time.Since(r.startTime).Seconds()

	// 创建输出事件
	event := NewOutputEvent(elapsed, string(data))
	eventData, err := event.MarshalEvent()
	if err != nil {
		return fmt.Errorf("marshal output event failed: %w", err)
	}

	// 写入事件
	if _, err := r.tempWriter.Write(eventData); err != nil {
		return fmt.Errorf("write output event failed: %w", err)
	}

	r.outputCount++

	// 输出流解析作为未来增强(当前返回空)
	// 未来可用于检测Tab补全后的完整命令
	_ = r.cmdParser.ParseOutput(data, elapsed)

	return nil
}

// RecordResize 记录窗口调整(实现Recorder接口)
func (r *AsciicastRecorder) RecordResize(width, height int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.started || r.stopped {
		return fmt.Errorf("recorder not active")
	}

	// 计算相对时间戳
	elapsed := time.Since(r.startTime).Seconds()

	// 创建窗口调整事件
	event := NewResizeEvent(elapsed, width, height)
	eventData, err := event.MarshalEvent()
	if err != nil {
		return fmt.Errorf("marshal resize event failed: %w", err)
	}

	// 写入事件
	if _, err := r.tempWriter.Write(eventData); err != nil {
		return fmt.Errorf("write resize event failed: %w", err)
	}

	r.resizeCount++

	return nil
}

// Stop 停止录制并保存(实现Recorder接口)
func (r *AsciicastRecorder) Stop() (*RecordingInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.started {
		return nil, fmt.Errorf("recorder not started")
	}

	if r.stopped {
		return nil, fmt.Errorf("recorder already stopped")
	}

	// 刷新缓冲区
	if err := r.tempWriter.Flush(); err != nil {
		return nil, fmt.Errorf("flush buffer failed: %w", err)
	}

	// 关闭临时文件
	tempFilePath := r.tempFile.Name()
	if err := r.tempFile.Close(); err != nil {
		return nil, fmt.Errorf("close temp file failed: %w", err)
	}

	// 重新打开用于读取
	tempFile, err := os.Open(tempFilePath)
	if err != nil {
		return nil, fmt.Errorf("reopen temp file failed: %w", err)
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFilePath) // 清理临时文件
	}()

	// 保存到存储
	filePath, fileSize, err := r.storage.Save(r.metadata.SessionID, tempFile)
	if err != nil {
		return nil, fmt.Errorf("save to storage failed: %w", err)
	}

	// 计算持续时长
	duration := int(time.Since(r.startTime).Seconds())

	r.stopped = true

	return &RecordingInfo{
		SessionID:    r.metadata.SessionID,
		FilePath:     filePath,
		FileSize:     fileSize,
		Duration:     duration,
		InputCount:   r.inputCount,
		OutputCount:  r.outputCount,
		ResizeCount:  r.resizeCount,
		CommandCount: r.commandCount,
	}, nil
}

// GetSessionID 获取会话ID(实现Recorder接口)
func (r *AsciicastRecorder) GetSessionID() string {
	return r.metadata.SessionID
}

// GetMetadata 获取录制元数据(实现Recorder接口)
func (r *AsciicastRecorder) GetMetadata() RecordingMetadata {
	return r.metadata
}

// GetParsedCommands 获取解析后的命令列表
func (r *AsciicastRecorder) GetParsedCommands() []ParsedCommand {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cmdParser == nil {
		return []ParsedCommand{}
	}

	return r.cmdParser.GetCommands()
}
