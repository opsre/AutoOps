// Agent模块初始化
package agent

import (
	"fmt"
	"dodevops-api/api/agent/service"
)

// InitAgentModule 初始化Agent模块
func InitAgentModule() error {
	fmt.Println("正在初始化智能Agent模块...")

	// 1. 初始化默认Prompt模板
	if err := service.InitDefaultTemplates(); err != nil {
		return fmt.Errorf("初始化Prompt模板失败: %v", err)
	}
	fmt.Println("✓ Prompt模板初始化完成")

	// 2. 获取调度器单例并启动
	scheduler := service.GetAgentScheduler()

	// 3. 注册处理器
	scheduler.RegisterHandler(service.NewInspectionHandler())
	fmt.Println("✓ 巡检处理器已注册")

	// 4. 启动调度器
	scheduler.Start()
	fmt.Println("✓ Agent调度器已启动")

	// TODO: 后续添加其他处理器
	// scheduler.RegisterHandler(service.NewIncidentHandler())
	// scheduler.RegisterHandler(service.NewPredictionHandler())

	fmt.Println("智能Agent模块初始化完成！")
	return nil
}

// StopAgentModule 停止Agent模块
func StopAgentModule() {
	scheduler := service.GetAgentScheduler()
	scheduler.Stop()
	fmt.Println("Agent调度器已停止")
}
