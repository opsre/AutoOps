package main

import (
	"dodevops-api/common/config"
	"dodevops-api/common/snmp"
	"dodevops-api/pkg/db"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// SNMP 采集守护进程
// 功能:
// 1. 从数据库读取启用监控的SNMP设备
// 2. 定时采集SNMP指标
// 3. 推送到Pushgateway
// 4. Prometheus从Pushgateway拉取数据

func main() {
	log.Println("==========================================")
	log.Println("  SNMP Collector Daemon Starting...")
	log.Println("==========================================")

	// 1. 加载配置
	if err := config.LoadConfig(""); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 2. 初始化数据库
	if err := db.SetupDBLink(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 3. 获取配置
	pushgatewayURL := config.Config.Monitor.Pushgateway.URL

	log.Printf("Pushgateway URL: %s", pushgatewayURL)

	if pushgatewayURL == "" {
		log.Fatal("Pushgateway URL 未配置,请在 config.yaml 中配置 monitor.pushgateway.url")
	}

	// 4. 启动SNMP采集管理器
	snmp.GetManager().Start()

	// 5. 启动信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("==========================================")
	log.Println("  SNMP Collector Daemon Started")
	log.Println("  按 Ctrl+C 停止服务")
	log.Println("==========================================")

	// 6. 等待退出信号
	sig := <-sigChan
	log.Printf("接收到信号: %v, 正在关闭...", sig)

	// 7. 停止SNMP采集管理器
	snmp.GetManager().Stop()

	log.Println("SNMP Collector Daemon 已停止")
}
