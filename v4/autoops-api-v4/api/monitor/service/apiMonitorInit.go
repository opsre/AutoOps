package service

import "log"

var apiMonitorScheduler *APIMonitorScheduler

// InitAPIMonitorScheduler 初始化API监控调度器
func InitAPIMonitorScheduler() error {
	log.Println("🚀 初始化API监控调度器...")

	apiMonitorScheduler = NewAPIMonitorScheduler()
	apiMonitorScheduler.Start()

	log.Println("✅ API监控调度器初始化完成")
	return nil
}

// GetAPIMonitorScheduler 获取API监控调度器实例
func GetAPIMonitorScheduler() *APIMonitorScheduler {
	return apiMonitorScheduler
}

// StopAPIMonitorScheduler 停止API监控调度器
func StopAPIMonitorScheduler() {
	if apiMonitorScheduler != nil {
		log.Println("⏹  停止API监控调度器...")
		apiMonitorScheduler.Stop()
		log.Println("✅ API监控调度器已停止")
	}
}
