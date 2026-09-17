package guacamole

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ProxyClient 智能代理客户端（处理初始握手）
type ProxyClient struct {
	wsConn    *websocket.Conn
	guacdConn net.Conn
	closeOnce sync.Once
	closed    chan struct{}
	done      chan struct{} // 用于通知连接完全结束
	rdpConfig RDPConfig
}

// NewProxyClient 创建智能代理客户端
func NewProxyClient(wsConn *websocket.Conn, guacdAddr string, rdpConfig RDPConfig) (*ProxyClient, error) {
	// 连接 guacd
	log.Printf("========== [Proxy层] 创建代理客户端 ==========")
	log.Printf("[Proxy] guacd地址: %s", guacdAddr)
	log.Printf("[Proxy] RDP目标: %s:%d", rdpConfig.Hostname, rdpConfig.Port)
	log.Printf("[Proxy] 正在连接 guacd...")

	guacdConn, err := net.Dial("tcp", guacdAddr)
	if err != nil {
		log.Printf("[Proxy] ❌ 连接 guacd 失败: %v", err)
		return nil, err
	}

	log.Printf("[Proxy] ✅ 成功连接到 guacd: %s", guacdAddr)
	return &ProxyClient{
		wsConn:    wsConn,
		guacdConn: guacdConn,
		closed:    make(chan struct{}),
		done:      make(chan struct{}),
		rdpConfig: rdpConfig,
	}, nil
}

// Start 启动主动握手代理（后端负责协议握手）
func (p *ProxyClient) Start() error {
	log.Println("[Proxy] ========== 开始主动握手模式 ==========")

	// 1. 发送 select 指令（选择 RDP 协议）
	selectCmd := "6.select,3.rdp;"
	log.Printf("[Proxy] → guacd: %s", selectCmd)
	if _, err := p.guacdConn.Write([]byte(selectCmd)); err != nil {
		log.Printf("[Proxy] ❌ 发送 select 失败: %v", err)
		return err
	}

	// 2. 读取 args 响应
	argsResp, err := p.readGuacdInstruction()
	if err != nil {
		log.Printf("[Proxy] ❌ 读取 args 失败: %v", err)
		return err
	}
	log.Printf("[Proxy] ← guacd: %s", argsResp)

	// 3. 解析 args 响应
	argNames, err := p.parseArgs(argsResp)
	if err != nil {
		log.Printf("[Proxy] ❌ 解析 args 失败: %v", err)
		return err
	}
	log.Printf("[Proxy] 解析到参数列表: %v", argNames)

	// 4. 构造并发送 connect 指令
	connectCmd := p.buildConnectCommand(argNames)
	log.Printf("[Proxy] → guacd: %s", connectCmd)
	if _, err := p.guacdConn.Write([]byte(connectCmd)); err != nil {
		log.Printf("[Proxy] ❌ 发送 connect 失败: %v", err)
		return err
	}

	// 5. 读取 ready 响应
	readyResp, err := p.readGuacdInstruction()
	if err != nil {
		log.Printf("[Proxy] ❌ 读取 ready 失败: %v", err)
		return err
	}
	log.Printf("[Proxy] ← guacd: %s", readyResp)

	if !strings.Contains(readyResp, "ready") {
		return fmt.Errorf("guacd 未返回 ready 指令: %s", readyResp)
	}

	log.Println("[Proxy] ✅ Guacamole 握手完成，开始双向透传")

	// 6. 启动双向透传（异步）
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		p.guacdToWebSocket()
	}()

	go func() {
		defer wg.Done()
		p.webSocketToGuacd()
	}()

	// 等待双向转发都完成后通知 done
	go func() {
		wg.Wait()
		close(p.done)
		log.Println("[Proxy] 双向转发已全部结束")
	}()

	log.Println("[Proxy] 双向转发已启动（后台运行）")
	return nil
}

// Wait 等待连接关闭
func (p *ProxyClient) Wait() {
	<-p.done
}

// guacdToWebSocket 从 guacd 读取数据并转发到 WebSocket
func (p *ProxyClient) guacdToWebSocket() {
	defer func() {
		log.Println("[Guac-Proxy] guacd->WebSocket 转发结束")
		p.Close()
	}()

	buffer := make([]byte, 8192)
	for {
		select {
		case <-p.closed:
			log.Println("[Guac-Proxy] guacd->WebSocket 收到关闭信号")
			return
		default:
		}

		// 设置读取超时，避免永久阻塞
		p.guacdConn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		n, err := p.guacdConn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 超时继续循环，检查 closed 信号
				continue
			}
			if err != io.EOF {
				log.Printf("[Guac-Proxy] 读取 guacd 数据失败: %v", err)
			} else {
				log.Println("[Guac-Proxy] guacd 连接关闭 (EOF)")
			}
			return
		}

		log.Printf("[Guac-Proxy] ← guacd: 收到 %d 字节", n)

		// 直接转发到 WebSocket（文本模式）
		if err := p.wsConn.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
			log.Printf("[Guac-Proxy] 发送 WebSocket 数据失败: %v", err)
			return
		}
		log.Printf("[Guac-Proxy] → WebSocket: 已发送 %d 字节", n)
	}
}

// webSocketToGuacd 从 WebSocket 读取数据并转发到 guacd
func (p *ProxyClient) webSocketToGuacd() {
	defer func() {
		log.Println("[Guac-Proxy] WebSocket->guacd 转发结束")
		p.Close()
	}()

	for {
		select {
		case <-p.closed:
			log.Println("[Guac-Proxy] WebSocket->guacd 收到关闭信号")
			return
		default:
		}

		_, message, err := p.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("[Guac-Proxy] WebSocket 意外关闭: %v", err)
			} else {
				log.Println("[Guac-Proxy] WebSocket 正常关闭")
			}
			return
		}

		log.Printf("[Guac-Proxy] ← WebSocket: 收到 %d 字节: %s", len(message), string(message[:min(100, len(message))]))

		// 直接转发到 guacd
		if _, err := p.guacdConn.Write(message); err != nil {
			log.Printf("[Guac-Proxy] 发送 guacd 数据失败: %v", err)
			return
		}
		log.Printf("[Guac-Proxy] → guacd: 已发送 %d 字节", len(message))
	}
}

// Close 关闭所有连接
func (p *ProxyClient) Close() {
	p.closeOnce.Do(func() {
		log.Println("[Guac-Proxy] 正在关闭所有连接...")
		close(p.closed) // 通知所有 goroutine 退出
		if p.guacdConn != nil {
			p.guacdConn.Close()
			log.Println("[Guac-Proxy] guacd 连接已关闭")
		}
		if p.wsConn != nil {
			p.wsConn.Close()
			log.Println("[Guac-Proxy] WebSocket 连接已关闭")
		}
	})
}

// readGuacdInstruction 从 guacd 读取一条完整指令
func (p *ProxyClient) readGuacdInstruction() (string, error) {
	buffer := make([]byte, 8192)
	result := strings.Builder{}

	for {
		n, err := p.guacdConn.Read(buffer)
		if err != nil {
			return "", err
		}

		result.Write(buffer[:n])

		// 检查是否包含指令结束符 ';'
		if strings.Contains(result.String(), ";") {
			return strings.Split(result.String(), ";")[0] + ";", nil
		}
	}
}

// parseArgs 解析 args 指令，提取参数名列表
func (p *ProxyClient) parseArgs(argsResp string) ([]string, error) {
	// 格式: 4.args,8.hostname,4.port,8.username,8.password,...;
	if !strings.HasPrefix(argsResp, "4.args,") {
		return nil, fmt.Errorf("无效的 args 响应: %s", argsResp)
	}

	// 移除 "4.args," 和末尾的 ";"
	content := strings.TrimPrefix(argsResp, "4.args,")
	content = strings.TrimSuffix(content, ";")

	// 解析参数名
	var argNames []string
	parts := strings.Split(content, ",")

	for i := 0; i < len(parts); i++ {
		// 格式: "长度.值"
		if strings.Contains(parts[i], ".") {
			argName := strings.SplitN(parts[i], ".", 2)[1]
			argNames = append(argNames, argName)
		}
	}

	return argNames, nil
}

// buildConnectCommand 根据参数列表构造 connect 指令
func (p *ProxyClient) buildConnectCommand(argNames []string) string {
	// 准备参数值映射（仅针对真实参数名）
	paramValues := map[string]string{
		"hostname":    p.rdpConfig.Hostname,
		"port":        fmt.Sprintf("%d", p.rdpConfig.Port),
		"username":    p.rdpConfig.Username,
		"password":    p.rdpConfig.Password,
		"domain":      p.rdpConfig.Domain,
		"width":       "1920",
		"height":      "1080",
		"dpi":         "96",
		"ignore-cert": "true",
		"security":    "any",
		"color-depth": "32",
	}

	var params []string

	// 第一个参数：version → 必须传空字符串
	params = append(params, "0.") // 表示长度为0的空字符串

	// 从 argNames[1] 开始，才是真正的参数名
	for i := 1; i < len(argNames); i++ {
		argName := argNames[i]
		value := paramValues[argName]
		// 如果未配置该参数，传空
		if value == "" {
			params = append(params, "0.")
		} else {
			params = append(params, fmt.Sprintf("%d.%s", len(value), value))
		}
	}

	// 构造完整的 connect 指令
	// 格式: 7.connect,参数1,参数2,...;
	return fmt.Sprintf("7.connect,%s;", strings.Join(params, ","))
}
