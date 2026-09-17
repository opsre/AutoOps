package controller

import (
    "net/http"
    "strconv"
    "strings"

    cmdbsvc "dodevops-api/api/cmdb/service"
    "dodevops-api/common/result"

    "github.com/gin-gonic/gin"
)

// WindowsMonitorController 提供Windows主机监控API
type WindowsMonitorController struct {
    svc *cmdbsvc.WindowsMonitorService
}

func NewWindowsMonitorController() *WindowsMonitorController {
    return &WindowsMonitorController{svc: cmdbsvc.NewWindowsMonitorService()}
}

// GetWindowsMonitorMetrics 获取Windows主机关键指标（CPU/内存/磁盘/IO/网络/进程/uptime），支持30m/1d/7d
// @Summary 获取Windows主机监控数据
// @Description 获取Windows主机监控数据
// @Tags CMDB-Windows监控
// @Produce json
// @Param id path uint true "设备ID(cmdb_snmp_devices表，DeviceType=windows)"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/windows-devices/{id}/monitor [get]
func (c *WindowsMonitorController) GetWindowsMonitorMetrics(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        result.Failed(ctx, http.StatusBadRequest, "无效的设备ID")
        return
    }

    data, err := c.svc.GetWindowsMonitorMetrics(uint(id))
    if err != nil {
        result.Failed(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    result.Success(ctx, data)
}

// GetWindowsRealtimeMetrics 获取Windows主机实时CPU和内存使用率
// @Summary 获取Windows主机实时资源使用率
// @Description 获取Windows主机实时资源使用率
// @Tags CMDB-Windows监控
// @Produce json
// @Param id path uint true "设备ID(cmdb_snmp_devices表，DeviceType=windows)"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/windows-devices/{id}/realtime [get]
func (c *WindowsMonitorController) GetWindowsRealtimeMetrics(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        result.Failed(ctx, http.StatusBadRequest, "无效的设备ID")
        return
    }

    data, err := c.svc.GetWindowsRealtimeMetrics(uint(id))
    if err != nil {
        result.Failed(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    result.Success(ctx, data)
}

// GetBatchWindowsRealtimeMetrics 批量获取Windows主机实时CPU和内存使用率
// @Summary 批量获取Windows主机实时资源使用率
// @Description 批量获取Windows主机实时资源使用率
// @Tags CMDB-Windows监控
// @Produce json
// @Param ids query string true "设备ID列表，逗号分隔，例如: 1,2,3"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/windows-devices/batch-realtime [get]
func (c *WindowsMonitorController) GetBatchWindowsRealtimeMetrics(ctx *gin.Context) {
    idsStr := ctx.Query("ids")
    if idsStr == "" {
        result.Failed(ctx, http.StatusBadRequest, "参数ids不能为空")
        return
    }

    // 解析ID列表
    idStrs := strings.Split(idsStr, ",")
    deviceIDs := make([]uint, 0, len(idStrs))
    for _, idStr := range idStrs {
        idStr = strings.TrimSpace(idStr)
        if idStr == "" {
            continue
        }
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
            result.Failed(ctx, http.StatusBadRequest, "无效的设备ID: "+idStr)
            return
        }
        deviceIDs = append(deviceIDs, uint(id))
    }

    if len(deviceIDs) == 0 {
        result.Failed(ctx, http.StatusBadRequest, "未提供有效的设备ID")
        return
    }

    data, err := c.svc.GetBatchWindowsRealtimeMetrics(deviceIDs)
    if err != nil {
        result.Failed(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    result.Success(ctx, data)
}

