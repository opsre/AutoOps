package controller

import (
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
	"strconv"
)

type WindowsExporterController struct {
	exporterService *service.WindowsExporterService
}

func NewWindowsExporterController() *WindowsExporterController {
	return &WindowsExporterController{
		exporterService: service.NewWindowsExporterService(),
	}
}

// InstallExporter 安装Windows Exporter
// @Summary 安装Windows Exporter
// @Description 在指定Windows设备上安装Windows Exporter
// @Tags CMDB-Windows-Exporter
// @Accept json
// @Produce json
// @Param request body service.InstallExporterRequest true "安装参数"
// @Success 200 {object} result.Result
// @Router /cmdb/windows-exporter/install [post]
func (c *WindowsExporterController) InstallExporter(ctx *gin.Context) {
	var req service.InstallExporterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, int(result.ApiCode.ValidationParameterError), "参数错误: "+err.Error())
		return
	}

	// 异步安装（避免超时）
	go func() {
		c.exporterService.InstallWindowsExporter(&req)
	}()

	result.Success(ctx, "安装任务已提交，请稍后查询安装状态")
}

// GetExporterStatus 获取Exporter状态
// @Summary 获取Exporter状态
// @Description 获取指定设备的Windows Exporter安装和运行状态
// @Tags CMDB-Windows-Exporter
// @Accept json
// @Produce json
// @Param deviceId path int true "设备ID"
// @Success 200 {object} result.Result{data=service.ExporterStatusResponse}
// @Router /cmdb/windows-exporter/status/{deviceId} [get]
func (c *WindowsExporterController) GetExporterStatus(ctx *gin.Context) {
	deviceIDStr := ctx.Param("deviceId")
	deviceID, err := strconv.ParseUint(deviceIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.ValidationParameterError), "设备ID格式错误")
		return
	}

	status, err := c.exporterService.GetExporterStatus(uint(deviceID))
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(ctx, status)
}

// UninstallExporter 卸载Exporter
// @Summary 卸载Windows Exporter
// @Description 从指定Windows设备卸载Windows Exporter
// @Tags CMDB-Windows-Exporter
// @Accept json
// @Produce json
// @Param deviceId path int true "设备ID"
// @Success 200 {object} result.Result
// @Router /cmdb/windows-exporter/uninstall/{deviceId} [delete]
func (c *WindowsExporterController) UninstallExporter(ctx *gin.Context) {
	deviceIDStr := ctx.Param("deviceId")
	deviceID, err := strconv.ParseUint(deviceIDStr, 10, 32)
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.ValidationParameterError), "设备ID格式错误")
		return
	}

	// 异步卸载
	go func() {
		c.exporterService.UninstallExporter(uint(deviceID))
	}()

	result.Success(ctx, "卸载任务已提交")
}

// BatchInstallExporter 批量安装Exporter
// @Summary 批量安装Windows Exporter
// @Description 在多个Windows设备上批量安装Windows Exporter
// @Tags CMDB-Windows-Exporter
// @Accept json
// @Produce json
// @Param request body BatchInstallRequest true "批量安装参数"
// @Success 200 {object} result.Result
// @Router /cmdb/windows-exporter/batch-install [post]
func (c *WindowsExporterController) BatchInstallExporter(ctx *gin.Context) {
	var req BatchInstallRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, int(result.ApiCode.ValidationParameterError), "参数错误: "+err.Error())
		return
	}

	// 异步批量安装
	go func() {
		installReq := &service.InstallExporterRequest{
			InstallerURL: req.InstallerURL,
			Version:      req.Version,
			ExporterPort: req.ExporterPort,
		}
		c.exporterService.BatchInstallExporter(installReq, req.DeviceIDs)
	}()

	result.Success(ctx, "批量安装任务已提交")
}

// GetAllWindowsDevicesStatus 获取所有Windows设备状态
// @Summary 获取所有Windows设备的Exporter状态
// @Description 获取所有Windows设备的Windows Exporter安装和运行状态
// @Tags CMDB-Windows-Exporter
// @Accept json
// @Produce json
// @Success 200 {object} result.Result{data=[]service.ExporterStatusResponse}
// @Router /cmdb/windows-exporter/status-all [get]
func (c *WindowsExporterController) GetAllWindowsDevicesStatus(ctx *gin.Context) {
	statuses, err := c.exporterService.GetAllWindowsDevicesStatus()
	if err != nil {
		result.Failed(ctx, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(ctx, statuses)
}

// BatchInstallRequest 批量安装请求
type BatchInstallRequest struct {
	DeviceIDs    []uint `json:"deviceIds" validate:"required"`
	InstallerURL string `json:"installerUrl" validate:"required,url"`
	Version      string `json:"version" validate:"required"`
	ExporterPort int    `json:"exporterPort"`
}
