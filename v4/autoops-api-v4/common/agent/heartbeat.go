package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// HeartbeatData 心跳数据结构
type HeartbeatData struct {
	PID        int    `json:"pid"`
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	Token      string `json:"token"`
	Version    string `json:"version"`    // Agent版本号
	VersionTag string `json:"versionTag"` // Agent版本标签
}

// getLocalIP 获取本机IP地址
func getLocalIP() string {
	// 尝试连接外部地址来获取本机IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Printf("Error getting local IP: %v", err)
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// sendHeartbeat 发送心跳到服务端
func sendHeartbeat(serverURL, token string) {
	// 获取当前进程ID
	pid := os.Getpid()

	// 获取本机IP
	localIP := getLocalIP()
	if localIP == "" {
		log.Println("Warning: Could not determine local IP address")
		return
	}

	// 构造心跳数据
	heartbeat := HeartbeatData{
		PID:        pid,
		IP:         localIP,
		Port:       12580, // Agent默认端口
		Token:      token,
		Version:    AgentVersion,
		VersionTag: AgentVersionTag,
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(heartbeat)
	if err != nil {
		log.Printf("Error marshaling heartbeat data: %v", err)
		return
	}

	// 发送HTTP POST请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending heartbeat: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Printf("Heartbeat sent successfully - PID: %d, IP: %s, Port: %d, Token: %s", pid, localIP, 12580, token)
	} else {
		log.Printf("Heartbeat failed with status: %d", resp.StatusCode)
	}
}
