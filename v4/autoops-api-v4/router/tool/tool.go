// router/tool/tool.go
package tool

import (
	"dodevops-api/api/tool/controller"
	"github.com/gin-gonic/gin"
)

// RegisterToolRoutes 注册导航工具相关路由
func RegisterToolRoutes(router *gin.RouterGroup) {
	// 导航工具管理
	router.POST("/tool", controller.CreateTool)        // 创建导航工具
	router.GET("/tool/:id", controller.GetToolByID)    // 获取导航工具详情
	router.PUT("/tool", controller.UpdateTool)         // 更新导航工具
	router.DELETE("/tool/:id", controller.DeleteTool)  // 删除导航工具
	router.GET("/tool/list", controller.GetToolList)   // 获取导航工具列表（分页）
	router.GET("/tool/all", controller.GetAllTools)    // 获取所有导航工具（不分页）

	// 服务部署管理
	router.GET("/tool/services", controller.GetServicesList)             // 获取可部署服务列表
	router.GET("/tool/services/:serviceId", controller.GetServiceDetail) // 获取服务详情
	router.POST("/tool/deploy", controller.CreateDeploy)                 // 创建部署任务
	router.GET("/tool/deploy/list", controller.GetDeployList)            // 获取部署历史列表
	router.GET("/tool/deploy/:id/status", controller.GetDeployStatus)    // 获取部署状态
	router.DELETE("/tool/deploy/:id", controller.DeleteDeploy)           // 卸载服务

	// 飞书知识库（方案A：纯代理，数据全在飞书，无本地数据库表）
	router.GET("/tool/feishu/apps", controller.GetFeishuApps)                         // 获取飞书应用列表
	router.GET("/tool/feishu/validate", controller.ValidateFeishu)                    // 诊断：验证 token + 原始响应
	router.GET("/tool/feishu/spaces", controller.GetFeishuWikiSpaces)                 // 获取 Wiki 空间列表（支持 openId 参数）
	router.GET("/tool/feishu/spaces/:spaceId/nodes", controller.GetFeishuWikiNodes)   // 获取 Wiki 节点树（支持 openId 参数）
	router.GET("/tool/feishu/docs/:docToken/content", controller.GetFeishuDocContent) // 获取文档内容（支持 openId 参数）
	router.GET("/tool/feishu/search", controller.SearchFeishuDocs)                    // 全文搜索（支持 openId 参数）

	// 飞书 OAuth（用于访问个人知识库）
	router.GET("/tool/feishu/oauth/url", controller.GetFeishuOAuthURL)              // 获取 OAuth 授权链接
	// callback 路由注册在公共路由组（无需 JWT），见 router/router.go
	router.GET("/tool/feishu/oauth/users", controller.ListFeishuUserTokens)         // 列出已授权用户
	router.DELETE("/tool/feishu/oauth/users", controller.RevokeFeishuUserToken)     // 撤销用户授权
}
