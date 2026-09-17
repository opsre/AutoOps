// Asciicast格式处理模块
// 实现asciicast v2格式的编码和解析
package recording

import (
	"encoding/json"
	"fmt"
	"time"
)

// AsciicastHeader asciicast v2文件头部
type AsciicastHeader struct {
	Version       int               `json:"version"`                   // 格式版本，固定为2
	Width         int               `json:"width"`                     // 终端宽度(列数)
	Height        int               `json:"height"`                    // 终端高度(行数)
	Timestamp     int64             `json:"timestamp"`                 // Unix时间戳(会话开始时间)
	Title         string            `json:"title,omitempty"`           // 录制标题
	IdleTimeLimit float64           `json:"idle_time_limit,omitempty"` // 空闲时间限制(秒)
	Env           map[string]string `json:"env,omitempty"`             // 环境变量
}

// NewAsciicastHeader 创建asciicast头部
func NewAsciicastHeader(width, height int, title string) *AsciicastHeader {
	return &AsciicastHeader{
		Version:   2,
		Width:     width,
		Height:    height,
		Timestamp: time.Now().Unix(),
		Title:     title,
		Env: map[string]string{
			"TERM":  "xterm-256color",
			"SHELL": "/bin/bash",
		},
	}
}

// MarshalHeader 将头部序列化为JSON字符串(带换行符)
func (h *AsciicastHeader) MarshalHeader() ([]byte, error) {
	data, err := json.Marshal(h)
	if err != nil {
		return nil, fmt.Errorf("marshal header failed: %w", err)
	}
	// 添加换行符
	return append(data, '\n'), nil
}

// AsciicastEvent asciicast事件
// 格式: [时间戳, 类型, 数据]
type AsciicastEvent struct {
	Time      float64 // 相对时间戳(秒，从会话开始计算)
	EventType string  // 事件类型: "o"=输出, "i"=输入, "r"=调整大小
	Data      string  // 事件数据内容
}

// 事件类型常量
const (
	EventTypeOutput = "o" // 输出事件(SSH服务器返回的数据)
	EventTypeInput  = "i" // 输入事件(用户键入的数据)
	EventTypeResize = "r" // 窗口调整事件
)

// NewOutputEvent 创建输出事件
func NewOutputEvent(elapsed float64, data string) *AsciicastEvent {
	return &AsciicastEvent{
		Time:      elapsed,
		EventType: EventTypeOutput,
		Data:      data,
	}
}

// NewInputEvent 创建输入事件
func NewInputEvent(elapsed float64, data string) *AsciicastEvent {
	return &AsciicastEvent{
		Time:      elapsed,
		EventType: EventTypeInput,
		Data:      data,
	}
}

// NewResizeEvent 创建窗口调整事件
func NewResizeEvent(elapsed float64, width, height int) *AsciicastEvent {
	return &AsciicastEvent{
		Time:      elapsed,
		EventType: EventTypeResize,
		Data:      fmt.Sprintf("%dx%d", width, height),
	}
}

// MarshalEvent 将事件序列化为JSON数组字符串(带换行符)
// 格式: [时间戳, "类型", "数据"]\n
func (e *AsciicastEvent) MarshalEvent() ([]byte, error) {
	// 构造三元组数组
	arr := []interface{}{e.Time, e.EventType, e.Data}
	data, err := json.Marshal(arr)
	if err != nil {
		return nil, fmt.Errorf("marshal event failed: %w", err)
	}
	// 添加换行符
	return append(data, '\n'), nil
}

// UnmarshalEvent 从JSON数组字节反序列化事件
func UnmarshalEvent(data []byte) (*AsciicastEvent, error) {
	var arr []interface{}
	if err := json.Unmarshal(data, &arr); err != nil {
		return nil, fmt.Errorf("unmarshal event failed: %w", err)
	}

	if len(arr) != 3 {
		return nil, fmt.Errorf("invalid event format: expected 3 elements, got %d", len(arr))
	}

	// 解析时间戳
	timestamp, ok := arr[0].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid timestamp type: %T", arr[0])
	}

	// 解析事件类型
	eventType, ok := arr[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid event type: %T", arr[1])
	}

	// 解析数据
	eventData, ok := arr[2].(string)
	if !ok {
		return nil, fmt.Errorf("invalid event data type: %T", arr[2])
	}

	return &AsciicastEvent{
		Time:      timestamp,
		EventType: eventType,
		Data:      eventData,
	}, nil
}

// IsOutputEvent 判断是否为输出事件
func (e *AsciicastEvent) IsOutputEvent() bool {
	return e.EventType == EventTypeOutput
}

// IsInputEvent 判断是否为输入事件
func (e *AsciicastEvent) IsInputEvent() bool {
	return e.EventType == EventTypeInput
}

// IsResizeEvent 判断是否为窗口调整事件
func (e *AsciicastEvent) IsResizeEvent() bool {
	return e.EventType == EventTypeResize
}
