package controller

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/result"
	"dodevops-api/common/snmp"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type CmdbSnmpDeviceController struct {
	service        *service.CmdbSnmpDeviceService
	monitorService *service.SnmpMonitorService
}

func NewCmdbSnmpDeviceController() *CmdbSnmpDeviceController {
	return &CmdbSnmpDeviceController{
		service:        service.NewCmdbSnmpDeviceService(),
		monitorService: service.NewSnmpMonitorService(),
	}
}

// CreateSnmpDevice 创建SNMP设备(极简API)
// @Summary 创建SNMP设备
// @Description 创建SNMP设备
// @Tags CMDB-SNMP设备
// @Accept json
// @Produce json
// @Param device body model.CreateSnmpDeviceDto true "设备信息"
// @Success 200 {object} result.Result{data=model.CmdbSnmpDevice}
// @Router /api/v1/cmdb/snmp-devices [post]
func (ctrl *CmdbSnmpDeviceController) CreateSnmpDevice(c *gin.Context) {
	var dto model.CreateSnmpDeviceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	username, _ := c.Get("username")
	usernameStr, ok := username.(string)
	if !ok {
		usernameStr = "system"
	}

	device, err := ctrl.service.CreateSnmpDevice(&dto, usernameStr)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建设备失败: "+err.Error())
		return
	}

	// 立即通知Manager重载设备列表
	snmp.GetManager().Reload()

	result.Success(c, device)
}

// UpdateSnmpDevice 更新SNMP设备
// @Summary 更新SNMP设备
// @Description 更新SNMP设备
// @Tags CMDB-SNMP设备
// @Accept json
// @Produce json
// @Param id path uint true "设备ID"
// @Param device body model.UpdateSnmpDeviceDto true "设备信息"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/snmp-devices/{id} [put]
func (ctrl *CmdbSnmpDeviceController) UpdateSnmpDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的设备ID")
		return
	}

	var dto model.UpdateSnmpDeviceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	dto.ID = uint(id)

	username, _ := c.Get("username")
	usernameStr, ok := username.(string)
	if !ok {
		usernameStr = "system"
	}

	err = ctrl.service.UpdateSnmpDevice(&dto, usernameStr)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新设备失败: "+err.Error())
		return
	}

	// 立即通知Manager重载设备列表
	snmp.GetManager().Reload()

	result.Success(c, nil)
}

// DeleteSnmpDevice 删除SNMP设备
// @Summary 删除SNMP设备
// @Description 删除SNMP设备
// @Tags CMDB-SNMP设备
// @Param id path uint true "设备ID"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/snmp-devices/{id} [delete]
func (ctrl *CmdbSnmpDeviceController) DeleteSnmpDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的设备ID")
		return
	}

	err = ctrl.service.DeleteSnmpDevice(uint(id))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除设备失败: "+err.Error())
		return
	}

	// 立即通知Manager重载设备列表
	snmp.GetManager().Reload()

	result.Success(c, nil)
}

// GetSnmpDeviceById 获取SNMP设备详情
// @Summary 获取SNMP设备详情
// @Description 获取SNMP设备详情
// @Tags CMDB-SNMP设备
// @Param id path uint true "设备ID"
// @Success 200 {object} result.Result{data=model.CmdbSnmpDevice}
// @Router /api/v1/cmdb/snmp-devices/{id} [get]
func (ctrl *CmdbSnmpDeviceController) GetSnmpDeviceById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的设备ID")
		return
	}

	device, err := ctrl.service.GetSnmpDeviceById(uint(id))
	if err != nil {
		result.Failed(c, http.StatusNotFound, "设备不存在")
		return
	}

	result.Success(c, device.ToVO())
}

// GetSnmpDeviceList 获取SNMP设备列表
// @Summary 获取SNMP设备列表
// @Description 获取SNMP设备列表
// @Tags CMDB-SNMP设备
// @Param deviceType query string false "设备类型"
// @Param status query string false "状态"
// @Param keyword query string false "关键词"
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Success 200 {object} result.Result{data=model.SnmpDeviceListResult}
// @Router /api/v1/cmdb/snmp-devices [get]
func (ctrl *CmdbSnmpDeviceController) GetSnmpDeviceList(c *gin.Context) {
	var query model.SnmpDeviceQueryDto

	query.DeviceType = c.Query("deviceType")
	query.Status = c.Query("status")
	query.Keyword = c.Query("keyword")

	monitorEnabledStr := c.Query("monitorEnabled")
	if monitorEnabledStr != "" {
		monitorEnabled, _ := strconv.Atoi(monitorEnabledStr)
		query.MonitorEnabled = monitorEnabled
	} else {
		query.MonitorEnabled = -1
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query.Page = page
	query.PageSize = pageSize

	listResult, err := ctrl.service.GetSnmpDeviceList(&query)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取设备列表失败: "+err.Error())
		return
	}

	// 转换为VO
	voItems := make([]*model.SnmpDeviceVO, 0, len(listResult.Items))
	for _, device := range listResult.Items {
		voItems = append(voItems, device.ToVO())
	}

	result.Success(c, map[string]interface{}{
		"items":    voItems,
		"total":    listResult.Total,
		"page":     listResult.Page,
		"pageSize": listResult.PageSize,
	})
}

// CollectNow 立即采集指定设备(手动触发)
// @Summary 立即采集指定设备
// @Description 立即采集指定设备
// @Tags CMDB-SNMP设备
// @Param id path uint true "设备ID"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/snmp-devices/{id}/collect [post]
func (ctrl *CmdbSnmpDeviceController) CollectNow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的设备ID")
		return
	}

	// 调用SNMP Manager立即采集
	err = snmp.GetManager().CollectNow(uint(id))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "采集失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{"message": "采集成功,数据已推送到Pushgateway"})
}

// EnableMonitor 启用设备监控(开启采集)
// @Summary 启用设备监控
// @Description 启用设备监控
// @Tags CMDB-SNMP设备
// @Param id path uint true "设备ID"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/snmp-devices/{id}/enable [post]
func (ctrl *CmdbSnmpDeviceController) EnableMonitor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的设备ID")
		return
	}

	// 更新数据库
	err = ctrl.service.SetMonitorEnabled(uint(id), 1)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "启用监控失败: "+err.Error())
		return
	}

	// 立即重载Manager,启动采集
	snmp.GetManager().Reload()

	result.Success(c, gin.H{"message": "已启用监控,采集器已启动"})
}

// DisableMonitor 停止设备监控(关闭采集)
// @Summary 停止设备监控
// @Description 停止设备监控
// @Tags CMDB-SNMP设备
// @Param id path uint true "设备ID"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/snmp-devices/{id}/disable [post]
func (ctrl *CmdbSnmpDeviceController) DisableMonitor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的设备ID")
		return
	}

	// 更新数据库
	err = ctrl.service.SetMonitorEnabled(uint(id), 0)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "停止监控失败: "+err.Error())
		return
	}

	// 立即重载Manager,停止采集
	snmp.GetManager().Reload()

	result.Success(c, gin.H{"message": "已停止监控,采集器已关闭"})
}

// GetDeviceMonitorMetrics 获取设备基础监控指标
// @Summary 获取设备基础监控指标
// @Description 获取设备的CPU使用率、内存使用率、设备温度、运行时间等实时监控数据
// @Tags CMDB-SNMP设备监控
// @Param id path uint true "设备ID"
// @Success 200 {object} result.Result{data=model.DeviceMonitorMetrics}
// @Router /api/v1/cmdb/snmp-devices/{id}/monitor [get]
func (ctrl *CmdbSnmpDeviceController) GetDeviceMonitorMetrics(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的设备ID")
		return
	}

	// 获取监控指标
	metrics, err := ctrl.monitorService.GetDeviceMonitorMetrics(uint(id))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取监控数据失败: "+err.Error())
		return
	}

	result.Success(c, metrics)
}

// GetDeviceResourceMetrics 获取设备资源使用率(轻量级API,仅返回CPU和内存)
// @Summary 获取设备资源使用率
// @Description 轻量级API,仅返回交换机的CPU使用率和内存使用率
// @Tags CMDB-SNMP设备监控
// @Param id path uint true "设备ID"
// @Success 200 {object} result.Result{data=model.DeviceResourceMetrics}
// @Router /api/v1/cmdb/snmp-devices/{id}/resource [get]
func (ctrl *CmdbSnmpDeviceController) GetDeviceResourceMetrics(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的设备ID")
		return
	}

	// 获取资源指标
	metrics, err := ctrl.monitorService.GetDeviceResourceMetrics(uint(id))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取资源数据失败: "+err.Error())
		return
	}

	result.Success(c, metrics)
}

// GetBatchDeviceResourceMetrics 批量获取设备资源使用率
// @Summary 批量获取设备资源使用率
// @Description 批量获取多个设备的CPU使用率和内存使用率,支持最多100个设备
// @Tags CMDB-SNMP设备监控
// @Param deviceIds query string true "设备ID列表,用逗号分隔,例如: 1,2,3,4,5"
// @Success 200 {object} result.Result{data=model.BatchResourceMetricsResult}
// @Router /api/v1/cmdb/snmp-devices/batch-resource [get]
func (ctrl *CmdbSnmpDeviceController) GetBatchDeviceResourceMetrics(c *gin.Context) {
	// 从查询参数获取设备ID列表
	deviceIdsStr := c.Query("deviceIds")
	if deviceIdsStr == "" {
		result.Failed(c, http.StatusBadRequest, "参数错误: deviceIds不能为空")
		return
	}

	// 解析设备ID列表(逗号分隔)
	deviceIdsStrArr := strings.Split(deviceIdsStr, ",")
	if len(deviceIdsStrArr) == 0 || len(deviceIdsStrArr) > 100 {
		result.Failed(c, http.StatusBadRequest, "参数错误: deviceIds数量必须在1-100之间")
		return
	}

	deviceIDs := make([]uint, 0, len(deviceIdsStrArr))
	for _, idStr := range deviceIdsStrArr {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
		if err != nil {
			result.Failed(c, http.StatusBadRequest, "参数错误: 无效的设备ID")
			return
		}
		deviceIDs = append(deviceIDs, uint(id))
	}

	// 获取批量资源指标
	metrics, err := ctrl.monitorService.GetBatchDeviceResourceMetrics(deviceIDs)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取资源数据失败: "+err.Error())
		return
	}

	result.Success(c, metrics)
}

