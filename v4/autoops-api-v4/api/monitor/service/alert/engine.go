package alert

import (
	"log"

	"dodevops-api/api/monitor/service/alert/scheduler"
)

// AlertEngine 告警引擎全局实例
var AlertEngine *scheduler.AlertScheduler

// InitAlertEngine 初始化告警引擎
// prometheusAddr: Prometheus地址，如 http://localhost:9090
func InitAlertEngine(prometheusAddr string) error {
	log.Println("初始化告警引擎...")

	// 创建调度器
	engine, err := scheduler.NewAlertScheduler(prometheusAddr)
	if err != nil {
		return err
	}

	AlertEngine = engine

	// 启动调度器
	AlertEngine.Start()

	log.Println("告警引擎初始化完成")
	return nil
}

// StopAlertEngine 停止告警引擎
func StopAlertEngine() {
	if AlertEngine != nil {
		AlertEngine.Stop()
		log.Println("告警引擎已停止")
	}
}

// GetAlertEngineStatus 获取告警引擎状态
func GetAlertEngineStatus() map[string]interface{} {
	if AlertEngine != nil {
		return AlertEngine.GetStatus()
	}
	return map[string]interface{}{
		"status": "not_initialized",
	}
}
