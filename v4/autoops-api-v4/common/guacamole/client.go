package guacamole

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"dodevops-api/common/config"

	"github.com/gorilla/websocket"
)

// GuacamoleClient Guacamole客户端
type GuacamoleClient struct {
	guacdConn net.Conn
	wsConn    *websocket.Conn
}

// RDPConfig RDP连接配置
type RDPConfig struct {
	Hostname string
	Port     int
	Username string
	Password string
	Domain   string
	Width    int
	Height   int
}

// NewGuacamoleClient 创建Guacamole客户端
func NewGuacamoleClient(wsConn *websocket.Conn) *GuacamoleClient {
	return &GuacamoleClient{
		wsConn: wsConn,
	}
}

// ConnectRDP 连接RDP
func (c *GuacamoleClient) ConnectRDP(rdpConfig *RDPConfig) error {
	// 1. 连接guacd守护进程
	guacdAddr := fmt.Sprintf("%s:%d", config.Config.Guacd.Host, config.Config.Guacd.Port)
	log.Printf("[Guacamole] 连接guacd: %s", guacdAddr)

	timeout := time.Duration(config.Config.Guacd.Timeout) * time.Second
	conn, err := net.DialTimeout("tcp", guacdAddr, timeout)
	if err != nil {
		return fmt.Errorf("连接guacd失败: %v", err)
	}
	c.guacdConn = conn
	log.Println("[Guacamole] 连接guacd成功")

	// 2. 发送Guacamole协议握手
	if err := c.sendHandshake(rdpConfig); err != nil {
		return fmt.Errorf("握手失败: %v", err)
	}

	// 3. 启动数据转发
	go c.guacdToWebSocket()
	go c.webSocketToGuacd()

	return nil
}

// sendHandshake 发送Guacamole协议握手
func (c *GuacamoleClient) sendHandshake(rdpConfig *RDPConfig) error {
	// 设置读取超时
	c.guacdConn.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer c.guacdConn.SetReadDeadline(time.Time{})

	response := make([]byte, 8192)

	// 1. 发送select指令
	selectCmd := "6.select,3.rdp;"
	log.Printf("[Guacamole] 发送select: %s", selectCmd)
	if _, err := c.guacdConn.Write([]byte(selectCmd)); err != nil {
		return fmt.Errorf("发送select失败: %v", err)
	}

	// 2. 读取args响应
	n, err := c.guacdConn.Read(response)
	if err != nil {
		return fmt.Errorf("读取args响应失败: %v", err)
	}
	argsResp := string(response[:n])
	log.Printf("[Guacamole] 收到args响应: %s", argsResp)

	// 3. 解析args，提取参数名（跳过 "4.args" 和 "VERSION_1_5_0"）
	argsParts := strings.Split(strings.TrimSuffix(argsResp, ";"), ",")
	if len(argsParts) < 3 || argsParts[0] != "4.args" {
		return fmt.Errorf("无效的args响应: %s", argsResp)
	}

	var paramNames []string
	for i := 2; i < len(argsParts); i++ { // 从索引2开始，跳过 "4.args" 和 "13.VERSION_1_5_0"
		p := argsParts[i]
		if idx := strings.Index(p, "."); idx != -1 {
			paramName := p[idx+1:]
			paramNames = append(paramNames, paramName)
		}
	}
	log.Printf("[Guacamole] 解析到 %d 个参数名", len(paramNames))

	// 4. 构建参数值映射表
	paramValues := map[string]string{
		"hostname":     rdpConfig.Hostname,
		"port":         fmt.Sprint(rdpConfig.Port),
		"domain":       rdpConfig.Domain,
		"username":     rdpConfig.Username,
		"password":     rdpConfig.Password,
		"width":        fmt.Sprint(rdpConfig.Width),
		"height":       fmt.Sprint(rdpConfig.Height),
		"dpi":          "96",
		"color-depth":  "16",
		"security":     "any",        // 关键：安全模式
		"ignore-cert":  "true",       // 关键：忽略证书
		"client-name":  "Go-Guacamole",
		// 其他参数留空
	}

	// 5. 按paramNames顺序构造connect参数
	var connectParams []string
	for _, name := range paramNames {
		value := paramValues[name]
		if value == "" {
			connectParams = append(connectParams, "0.")
		} else {
			connectParams = append(connectParams, fmt.Sprintf("%d.%s", len(value), value))
		}
	}

	connectCmd := "7.connect," + strings.Join(connectParams, ",") + ";"
	log.Printf("[Guacamole] 发送connect: 参数数量=%d", len(connectParams))
	if _, err := c.guacdConn.Write([]byte(connectCmd)); err != nil {
		return fmt.Errorf("发送connect失败: %v", err)
	}

	// 6. 读取ready响应
	n, err = c.guacdConn.Read(response)
	if err != nil {
		return fmt.Errorf("读取ready响应失败: %v", err)
	}
	readyResp := string(response[:n])
	log.Printf("[Guacamole] 收到ready响应: %s", readyResp)

	if !strings.HasPrefix(readyResp, "5.ready") {
		return fmt.Errorf("握手失败，期望ready，收到: %s", readyResp)
	}

	log.Println("[Guacamole] 握手完成，准备开始数据转发")
	return nil
}

// guacdToWebSocket 将guacd数据转发到WebSocket
func (c *GuacamoleClient) guacdToWebSocket() {
	defer func() {
		log.Println("[Guacamole] guacd->WebSocket转发结束")
		c.Close()
	}()

	buffer := make([]byte, 8192)
	firstRead := true
	for {
		n, err := c.guacdConn.Read(buffer)
		if err != nil {
			log.Printf("[Guacamole] 读取guacd数据失败: %v", err)
			return
		}

		// 打印第一次读取的数据用于调试
		if firstRead {
			log.Printf("[Guacamole] 首次从guacd读取 %d 字节: %s", n, string(buffer[:min(n, 200)]))
			firstRead = false
		}

		// 转发到WebSocket
		if err := c.wsConn.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
			log.Printf("[Guacamole] 发送WebSocket数据失败: %v", err)
			return
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// webSocketToGuacd 将WebSocket数据转发到guacd
func (c *GuacamoleClient) webSocketToGuacd() {
	defer func() {
		log.Println("[Guacamole] WebSocket->guacd转发结束")
		c.Close()
	}()

	for {
		_, message, err := c.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("[Guacamole] WebSocket读取错误: %v", err)
			}
			return
		}

		// 转发到guacd
		if _, err := c.guacdConn.Write(message); err != nil {
			log.Printf("[Guacamole] 发送guacd数据失败: %v", err)
			return
		}
	}
}

// Close 关闭连接
func (c *GuacamoleClient) Close() {
	if c.guacdConn != nil {
		c.guacdConn.Close()
	}
	if c.wsConn != nil {
		c.wsConn.Close()
	}
}
