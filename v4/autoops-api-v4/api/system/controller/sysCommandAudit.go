// SSH命令审计管理 控制器
package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/api/system/model"
	"dodevops-api/api/system/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type SysCommandAuditController struct {
	auditService *service.CommandAuditService
}

// NewSysCommandAuditController 创建命令审计控制器
func NewSysCommandAuditController() *SysCommandAuditController {
	return &SysCommandAuditController{
		auditService: service.NewCommandAuditService(),
	}
}

// GetCommandAuditList 获取命令审计列表
// @Summary 获取命令审计列表
// @Description 分页获取SSH会话命令审计列表
// @Tags 系统-命令审计
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(100)
// @Param recordingId query uint false "录制记录ID过滤"
// @Param sessionId query string false "会话ID过滤"
// @Param isSensitive query bool false "是否敏感过滤"
// @Param riskLevel query int false "风险等级过滤"
// @Param keyword query string false "命令关键词搜索"
// @Success 200 {object} result.Result
// @Router /api/v1/system/recordings/commands [get]
// @Security ApiKeyAuth
func (c *SysCommandAuditController) GetCommandAuditList(ctx *gin.Context) {
	var dto model.CommandAuditListDto

	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&dto); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if dto.Page <= 0 {
		dto.Page = 1
	}
	if dto.PageSize <= 0 {
		dto.PageSize = 100
	}

	// 查询列表
	commands, total, err := c.auditService.GetCommandList(&dto)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"total":    total,
		"page":     dto.Page,
		"pageSize": dto.PageSize,
		"list":     commands,
	})
}

// GetSensitiveCommandCount 获取敏感命令数量
// @Summary 获取敏感命令数量
// @Description 获取指定录制记录的敏感命令数量
// @Tags 系统-命令审计
// @Accept json
// @Produce json
// @Param recordingId path uint true "录制记录ID"
// @Success 200 {object} result.Result
// @Router /api/v1/system/recordings/{recordingId}/commands/sensitive-count [get]
// @Security ApiKeyAuth
func (c *SysCommandAuditController) GetSensitiveCommandCount(ctx *gin.Context) {
	recordingIDStr := ctx.Param("recordingId")
	recordingID, err := strconv.ParseUint(recordingIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "录制记录ID格式错误")
		return
	}

	count, err := c.auditService.GetSensitiveCommandCount(uint(recordingID))
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"count": count,
	})
}

// GetHighestRiskLevel 获取最高风险等级
// @Summary 获取最高风险等级
// @Description 获取指定录制记录的最高风险等级
// @Tags 系统-命令审计
// @Accept json
// @Produce json
// @Param recordingId path uint true "录制记录ID"
// @Success 200 {object} result.Result
// @Router /api/v1/system/recordings/{recordingId}/commands/highest-risk [get]
// @Security ApiKeyAuth
func (c *SysCommandAuditController) GetHighestRiskLevel(ctx *gin.Context) {
	recordingIDStr := ctx.Param("recordingId")
	recordingID, err := strconv.ParseUint(recordingIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "录制记录ID格式错误")
		return
	}

	riskLevel, err := c.auditService.GetHighestRiskLevel(uint(recordingID))
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"riskLevel": riskLevel,
	})
}
