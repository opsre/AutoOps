// Agent模块路由注册
package agent

import (
	"dodevops-api/api/agent/controller"
	"github.com/gin-gonic/gin"
)

// RegisterAgentRoutes 注册Agent相关路由
func RegisterAgentRoutes(router *gin.RouterGroup) {
	// ============== 运维巡检 ==============
	router.POST("/agent/inspection/create", controller.CreateInspection)              // 创建巡检任务
	router.GET("/agent/inspection/info/:id", controller.GetInspectionInfo)            // 获取巡检任务详情
	router.PUT("/agent/inspection/update", controller.UpdateInspection)               // 更新巡检任务
	router.DELETE("/agent/inspection/delete/:id", controller.DeleteInspection)        // 删除巡检任务
	router.GET("/agent/inspection/list", controller.GetInspectionList)                // 获取巡检任务列表
	router.POST("/agent/inspection/execute/:id", controller.ExecuteInspection)        // 执行巡检任务

	// 巡检记录
	router.GET("/agent/inspection/record/:id", controller.GetInspectionRecord)        // 获取巡检记录详情
	router.GET("/agent/inspection/record/list", controller.GetInspectionRecordList)   // 获取巡检记录列表

	// ============== AI对话（自然语言交互）==============
	router.POST("/agent/chat", controller.AgentChat)                                  // AI对话接口
	router.GET("/agent/chat/history", controller.GetChatHistory)                      // 获取对话历史
	router.DELETE("/agent/chat/clear", controller.ClearChatHistory)                   // 清除对话历史

	// TODO: 后续添加其他模块路由
	// ============== 故障分析 ==============
	// ============== 压力预警 ==============
	// ============== 知识库 ==============
	// ============== Agent任务 ==============
}
