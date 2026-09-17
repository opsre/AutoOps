// AI模块路由注册
package ai

import (
	"dodevops-api/api/ai/controller"
	"github.com/gin-gonic/gin"
)

// RegisterAiRoutes 注册AI相关路由
func RegisterAiRoutes(router *gin.RouterGroup) {
	// ============== AI模型管理 ==============
	router.POST("/ai/model/add", controller.CreateAiModel)                // 创建AI模型配置
	router.GET("/ai/model/info", controller.GetAiModelInfo)               // 获取模型详情
	router.PUT("/ai/model/update", controller.UpdateAiModel)              // 更新模型配置
	router.DELETE("/ai/model/delete", controller.DeleteAiModel)           // 删除模型
	router.GET("/ai/model/list", controller.GetAiModelList)               // 获取模型列表(支持分页、筛选)
	router.PUT("/ai/model/updateStatus", controller.UpdateAiModelStatus)  // 更新模型状态(启用/禁用)
	router.GET("/ai/model/enabled", controller.GetEnabledAiModels)        // 获取所有启用的模型

	// ============== AI对话 ==============
	router.POST("/ai/chat", controller.AiChat) // AI对话接口(支持阿里百炼等多种模型)
}
