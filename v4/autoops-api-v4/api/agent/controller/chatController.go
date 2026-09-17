// AI对话控制器
package controller

import (
	"strconv"

	"dodevops-api/api/agent/model"
	"dodevops-api/api/agent/service"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AgentChat AI对话接口
// @Summary AI对话接口
// @Description 通过自然语言与Agent交互，执行运维任务
// @Tags Agent对话
// @Accept json
// @Produce json
// @Param request body model.AgentChatRequest true "对话请求"
// @Success 200 {object} result.Result{data=model.AgentChatResponse}
// @Router /api/v1/agent/chat [post]
func AgentChat(c *gin.Context) {
	var req model.AgentChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, err := jwt.GetAdminId(c)
	if err != nil {
		result.Failed(c, int(result.ApiCode.NOAUTH), "未授权")
		return
	}

	chatService := service.NewAgentChatService()
	response, err := chatService.ProcessMessage(c.Request.Context(), userID, &req)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, response)
}

// GetChatHistory 获取对话历史
// @Summary 获取对话历史
// @Description 获取指定会话的对话历史记录
// @Tags Agent对话
// @Accept json
// @Produce json
// @Param sessionId query string false "会话ID"
// @Param limit query int false "限制数量（默认50）"
// @Success 200 {object} result.Result{data=[]model.AgentChatHistory}
// @Router /api/v1/agent/chat/history [get]
func GetChatHistory(c *gin.Context) {
	sessionID := c.Query("sessionId")
	if sessionID == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "会话ID不能为空")
		return
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 获取当前用户ID
	userID, err := jwt.GetAdminId(c)
	if err != nil {
		result.Failed(c, int(result.ApiCode.NOAUTH), "未授权")
		return
	}

	chatService := service.NewAgentChatService()
	history, err := chatService.GetChatHistory(c.Request.Context(), userID, sessionID, limit)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, history)
}

// ClearChatHistory 清除对话历史
// @Summary 清除对话历史
// @Description 清除指定会话的所有对话历史
// @Tags Agent对话
// @Accept json
// @Produce json
// @Param sessionId query string true "会话ID"
// @Success 200 {object} result.Result
// @Router /api/v1/agent/chat/clear [delete]
func ClearChatHistory(c *gin.Context) {
	sessionID := c.Query("sessionId")
	if sessionID == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "会话ID不能为空")
		return
	}

	chatService := service.NewAgentChatService()
	if err := chatService.ClearChatHistory(c.Request.Context(), sessionID); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, nil)
}
