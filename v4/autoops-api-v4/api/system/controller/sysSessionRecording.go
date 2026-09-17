// SSH会话录制管理 控制器
package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"dodevops-api/api/system/model"
	"dodevops-api/api/system/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type SysSessionRecordingController struct {
	recordingService *service.SessionRecordingService
}

// NewSysSessionRecordingController 创建录制管理控制器
func NewSysSessionRecordingController(storagePath string) (*SysSessionRecordingController, error) {
	svc, err := service.NewSessionRecordingService(storagePath)
	if err != nil {
		return nil, err
	}

	return &SysSessionRecordingController{
		recordingService: svc,
	}, nil
}

// GetRecordingList 获取录制列表
// @Summary 获取会话录制列表
// @Description 分页获取SSH会话录制列表
// @Tags 系统-会话录制
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param adminId query uint false "用户ID过滤"
// @Param hostId query uint false "主机ID过滤"
// @Param status query int false "状态过滤"
// @Param riskLevel query int false "风险等级过滤"
// @Param keyword query string false "关键词搜索"
// @Param startTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} result.Result
// @Router /api/v1/system/recordings [get]
// @Security ApiKeyAuth
func (c *SysSessionRecordingController) GetRecordingList(ctx *gin.Context) {
	var dto model.SessionRecordingListDto

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
		dto.PageSize = 10
	}

	// 查询列表
	recordings, total, err := c.recordingService.GetRecordingList(dto)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"total": total,
		"page":  dto.Page,
		"pageSize": dto.PageSize,
		"list":  recordings,
	})
}

// GetRecordingDetail 获取录制详情
// @Summary 获取会话录制详情
// @Description 根据ID获取会话录制详细信息
// @Tags 系统-会话录制
// @Accept json
// @Produce json
// @Param id path uint true "录制记录ID"
// @Success 200 {object} result.Result
// @Router /api/v1/system/recordings/{id} [get]
// @Security ApiKeyAuth
func (c *SysSessionRecordingController) GetRecordingDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	recording, err := c.recordingService.GetRecordingByID(uint(id))
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, recording)
}

// GetRecordingCast 获取录制cast文件内容
// @Summary 获取录制cast文件内容
// @Description 获取asciicast格式的录制文件内容，用于播放器回放
// @Tags 系统-会话录制
// @Accept json
// @Produce text/plain
// @Param id path uint true "录制记录ID"
// @Success 200 {string} string "asciicast文件内容"
// @Router /api/v1/system/recordings/{id}/cast [get]
// @Security ApiKeyAuth
func (c *SysSessionRecordingController) GetRecordingCast(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	reader, err := c.recordingService.GetRecordingContent(uint(id))
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "获取录制文件失败: "+err.Error())
		return
	}
	defer reader.Close()

	// 设置响应头
	ctx.Header("Content-Type", "text/plain; charset=utf-8")
	ctx.Header("Cache-Control", "public, max-age=3600")

	// 流式返回文件内容
	ctx.DataFromReader(http.StatusOK, -1, "text/plain", reader, nil)
}

// DownloadRecording 下载录制文件
// @Summary 下载会话录制文件
// @Description 下载asciicast格式的录制文件
// @Tags 系统-会话录制
// @Accept json
// @Produce application/octet-stream
// @Param id path uint true "录制记录ID"
// @Success 200 {file} file "录制文件"
// @Router /api/v1/system/recordings/{id}/download [get]
// @Security ApiKeyAuth
func (c *SysSessionRecordingController) DownloadRecording(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	filename, file, err := c.recordingService.DownloadRecording(uint(id))
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "下载失败: "+err.Error())
		return
	}
	defer file.Close()

	// 设置响应头
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", "application/octet-stream")

	// 流式返回文件
	ctx.DataFromReader(http.StatusOK, -1, "application/octet-stream", file, nil)
}

// DeleteRecording 删除录制记录
// @Summary 删除会话录制记录
// @Description 删除指定的会话录制记录（包括文件）
// @Tags 系统-会话录制
// @Accept json
// @Produce json
// @Param id path uint true "录制记录ID"
// @Success 200 {object} result.Result
// @Router /api/v1/system/recordings/{id} [delete]
// @Security ApiKeyAuth
func (c *SysSessionRecordingController) DeleteRecording(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.recordingService.DeleteRecording(uint(id)); err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	result.Success(ctx, nil)
}

// BatchDeleteRecording 批量删除录制记录
// @Summary 批量删除会话录制记录
// @Description 批量删除指定的会话录制记录（包括文件）
// @Tags 系统-会话录制
// @Accept json
// @Produce json
// @Param dto body model.BatchDeleteSessionRecordingDto true "批量删除参数"
// @Success 200 {object} result.Result
// @Router /api/v1/system/recordings/batch [delete]
// @Security ApiKeyAuth
func (c *SysSessionRecordingController) BatchDeleteRecording(ctx *gin.Context) {
	var dto model.BatchDeleteSessionRecordingDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := c.recordingService.BatchDeleteRecording(dto.IDs); err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "批量删除失败: "+err.Error())
		return
	}

	result.Success(ctx, nil)
}
