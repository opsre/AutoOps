// 带录制功能的WebSSH包装器
// 使用装饰器模式在原有WebSSH基础上添加录制能力
package websocket

import (
	"context"
	"dodevops-api/common/util/recording"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

// RecordingWebSSH 带录制功能的WebSSH包装器
type RecordingWebSSH struct {
	*WebSSH                     // 嵌入原有WebSSH
	recorder recording.Recorder // 录制器
	enabled  bool               // 是否启用录制
	done     chan struct{}      // 用于等待连接结束
	closeOnce sync.Once         // 确保关闭只执行一次
}

// NewRecordingWebSSH 创建带录制的WebSSH实例
// webssh: 原始WebSSH实例
// recorder: 录制器实例，如果为nil则不启用录制
func NewRecordingWebSSH(webssh *WebSSH, recorder recording.Recorder) *RecordingWebSSH {
	return &RecordingWebSSH{
		WebSSH:   webssh,
		recorder: recorder,
		enabled:  recorder != nil,
		done:     make(chan struct{}),
	}
}

func (r *RecordingWebSSH) shutdown() {
	r.closeOnce.Do(func() {
		if r.cancel != nil {
			r.cancel()
		}
		if r.conn != nil {
			_ = r.conn.Close()
		}
		if r.session != nil {
			_ = r.session.Close()
		}
		if r.client != nil {
			_ = r.client.Close()
		}
	})
}

// Connect 建立SSH会话并启动录制
// 重写原有Connect方法，在其基础上添加录制功能
func (r *RecordingWebSSH) Connect(wsConn *websocket.Conn) error {
	// 如果session已存在，使用现有session（跳板机模式）
	if r.session == nil {
		session, err := r.client.NewSession()
		if err != nil {
			return fmt.Errorf("failed to create session: %v", err)
		}
		r.session = session
	}
	r.conn = wsConn

	// 获取输入管道（如果还没有）
	if r.stdinPipe == nil {
		var err error
		r.stdinPipe, err = r.session.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdin pipe: %v", err)
		}
	}

	// 设置输出
	stdout, err := r.session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	// 默认终端尺寸
	cols := 80
	rows := 24

	// 尝试从第一个WebSocket消息中获取终端尺寸
	r.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, data, err := r.conn.ReadMessage()
	r.conn.SetReadDeadline(time.Time{}) // 清除deadline

	if err == nil {
		var msg map[string]interface{}
		if json.Unmarshal(data, &msg) == nil {
			if msg["type"] == "init" {
				if c, ok := msg["cols"].(float64); ok {
					cols = int(c)
				}
				if r, ok := msg["rows"].(float64); ok {
					rows = int(r)
				}
				log.Printf("收到前端终端尺寸: %dx%d", cols, rows)
			}
		}
	}

	// 启动录制（在请求PTY之前，记录终端尺寸）
	if r.enabled {
		// 更新录制元数据中的终端尺寸
		metadata := r.recorder.GetMetadata()
		metadata.TerminalWidth = cols
		metadata.TerminalHeight = rows

		if err := r.recorder.Start(metadata); err != nil {
			log.Printf("启动录制失败: %v (继续执行SSH连接)", err)
			// 录制失败不影响SSH连接
		} else {
			log.Printf("会话录制已启动: %s", r.recorder.GetSessionID())
		}
	}

	// 请求PTY和启动Shell
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
		ssh.ONLCR:         1,
		ssh.OCRNL:         0,
		ssh.INLCR:         0,
		ssh.IGNCR:         0,
		ssh.ICRNL:         1,
		ssh.OPOST:         1,
		ssh.ONLRET:        0,
		ssh.ONOCR:         0,
	}
	log.Printf("请求PTY: xterm-256color, 尺寸=%dx%d", cols, rows)
	if err := r.session.RequestPty("xterm-256color", rows, cols, modes); err != nil {
		return fmt.Errorf("request pty failed: %v", err)
	}

	if err := r.session.Shell(); err != nil {
		return fmt.Errorf("start shell failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	// 使用WaitGroup等待两个handler goroutine
	var wg sync.WaitGroup
	wg.Add(2)

	// 处理WebSocket输入（带录制）
	go func() {
		defer wg.Done()
		r.handleInputWithRecording(ctx)
	}()

	// 处理SSH输出（带录制）
	go func() {
		defer wg.Done()
		r.handleOutputWithRecording(ctx, stdout)
	}()

	// 等待所有handler结束，然后关闭done channel
	go func() {
		wg.Wait()
		close(r.done)
	}()

	return nil
}

// handleInputWithRecording 处理WebSocket输入并录制
func (r *RecordingWebSSH) handleInputWithRecording(ctx context.Context) {
	defer func() {
		// handler退出时取消context，通知其他goroutine也退出
		if r.cancel != nil {
			r.cancel()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, data, err := r.conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				// WS断开时主动关闭SSH/WS，避免另一个goroutine阻塞导致Wait不返回
				r.shutdown()
				return
			}

			// 处理窗口调整
			var msg map[string]interface{}
			if err := json.Unmarshal(data, &msg); err == nil {
				if msg["type"] == "resize" {
					cols, _ := strconv.Atoi(fmt.Sprintf("%v", msg["cols"]))
					rows, _ := strconv.Atoi(fmt.Sprintf("%v", msg["rows"]))
					if cols > 0 && rows > 0 {
						r.session.WindowChange(rows, cols)

						// 录制窗口调整事件
						if r.enabled {
							if err := r.recorder.RecordResize(cols, rows); err != nil {
								log.Printf("录制窗口调整事件失败: %v", err)
							}
						}
					}
					continue
				}
			}

			// 录制用户输入
			if r.enabled {
				if err := r.recorder.RecordInput(data); err != nil {
					log.Printf("录制输入事件失败: %v", err)
					// 录制失败不影响正常使用
				}
			}

			// 发送输入到SSH
			if _, err := r.stdinPipe.Write(data); err != nil {
				log.Printf("SSH write error: %v", err)
				r.shutdown()
				return
			}
		}
	}
}

// handleOutputWithRecording 处理SSH输出并录制
func (r *RecordingWebSSH) handleOutputWithRecording(ctx context.Context, stdout io.Reader) {
	defer func() {
		// handler退出时取消context，通知其他goroutine也退出
		if r.cancel != nil {
			r.cancel()
		}
	}()

	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := stdout.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("SSH read error: %v", err)
				}
				r.shutdown()
				return
			}
			if n > 0 {
				// 录制SSH输出
				if r.enabled {
					if err := r.recorder.RecordOutput(buf[:n]); err != nil {
						log.Printf("录制输出事件失败: %v", err)
						// 录制失败不影响正常使用
					}
				}

				// 发送到WebSocket
				if err := r.conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					log.Printf("WebSocket write error: %v", err)
					r.shutdown()
					return
				}
			}
		}
	}
}

// Close 关闭连接并停止录制
func (r *RecordingWebSSH) Close() error {
	// 停止录制
	if r.enabled && r.recorder != nil {
		if info, err := r.recorder.Stop(); err != nil {
			log.Printf("停止录制失败: %v", err)
		} else {
			log.Printf("录制已完成: 文件=%s, 大小=%d字节, 时长=%d秒, 输入=%d, 输出=%d",
				info.FilePath, info.FileSize, info.Duration, info.InputCount, info.OutputCount)
		}
	}

	// 调用原有Close方法
	r.shutdown()
	return nil
}

// GetRecorder 获取录制器实例
func (r *RecordingWebSSH) GetRecorder() recording.Recorder {
	return r.recorder
}

// FinalizeRecording 停止录制并返回录制信息（会将录制标记为禁用，避免Close重复Stop）
func (r *RecordingWebSSH) FinalizeRecording() (*recording.RecordingInfo, error) {
	if !r.enabled || r.recorder == nil {
		return nil, nil
	}

	info, err := r.recorder.Stop()
	r.enabled = false
	r.recorder = nil
	return info, err
}

// IsRecordingEnabled 检查是否启用录制
func (r *RecordingWebSSH) IsRecordingEnabled() bool {
	return r.enabled
}

// Wait 等待连接结束
func (r *RecordingWebSSH) Wait() {
	<-r.done
}
