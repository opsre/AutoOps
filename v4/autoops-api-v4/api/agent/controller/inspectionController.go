// 巡检控制器
package controller

import (
	"encoding/json"
	"strconv"

	"dodevops-api/api/agent/model"
	"dodevops-api/api/agent/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// CreateInspection 创建巡检任务
func CreateInspection(c *gin.Context) {
	var dto model.InspectionDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	inspectionService := service.NewInspectionService()
	if err := inspectionService.CreateInspection(c, &dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, nil)
}

// GetInspectionInfo 获取巡检任务详情
func GetInspectionInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	inspectionService := service.NewInspectionService()
	inspection, err := inspectionService.GetInspectionByID(c, uint(id))
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, inspection)
}

// UpdateInspection 更新巡检任务
func UpdateInspection(c *gin.Context) {
	var dto model.InspectionDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	if dto.ID == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "ID不能为空")
		return
	}

	inspectionService := service.NewInspectionService()
	if err := inspectionService.UpdateInspection(c, &dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, nil)
}

// DeleteInspection 删除巡检任务
func DeleteInspection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	inspectionService := service.NewInspectionService()
	if err := inspectionService.DeleteInspection(c, uint(id)); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, nil)
}

// GetInspectionList 获取巡检任务列表
func GetInspectionList(c *gin.Context) {
	query := &model.InspectionQuery{
		Name:         c.Query("name"),
		Type:         c.Query("type"),
		TargetType:   c.Query("targetType"),
		ScheduleType: c.Query("scheduleType"),
		PageNum:      1,
		PageSize:     10,
	}

	if enabled, err := strconv.Atoi(c.Query("enabled")); err == nil {
		query.Enabled = enabled
	}
	if pageNum, err := strconv.Atoi(c.Query("pageNum")); err == nil && pageNum > 0 {
		query.PageNum = pageNum
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}

	inspectionService := service.NewInspectionService()
	inspections, total, err := inspectionService.GetInspectionList(c, query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	data := map[string]interface{}{
		"list":     inspections,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	}

	result.Success(c, data)
}

// ExecuteInspection 执行巡检任务
func ExecuteInspection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	// 创建Agent任务
	params := map[string]interface{}{
		"inspectionId": uint(id),
	}
	paramsJSON, _ := json.Marshal(params)

	task := &model.AgentTask{
		Type:        "inspection",
		Name:        "手动执行巡检",
		Description: "手动触发的巡检任务",
		Priority:    5,
		Params:      string(paramsJSON),
		MaxRetry:    3,
	}

	// 提交到调度器
	scheduler := service.GetAgentScheduler()
	if err := scheduler.SubmitTask(task); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "提交任务失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"taskId": task.ID,
		"message": "巡检任务已提交",
	})
}

// GetInspectionRecord 获取巡检记录详情
func GetInspectionRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	inspectionService := service.NewInspectionService()
	record, err := inspectionService.GetRecordByID(c, uint(id))
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, record)
}

// GetInspectionRecordList 获取巡检记录列表
func GetInspectionRecordList(c *gin.Context) {
	query := &model.InspectionRecordQuery{
		HostName:  c.Query("hostName"),
		StartTime: c.Query("startTime"),
		EndTime:   c.Query("endTime"),
		PageNum:   1,
		PageSize:  10,
	}

	if inspectionID, err := strconv.ParseUint(c.Query("inspectionId"), 10, 64); err == nil {
		query.InspectionID = uint(inspectionID)
	}
	if hostID, err := strconv.ParseUint(c.Query("hostId"), 10, 64); err == nil {
		query.HostID = uint(hostID)
	}
	if status, err := strconv.Atoi(c.Query("status")); err == nil {
		query.Status = status
	}
	if riskLevel, err := strconv.Atoi(c.Query("riskLevel")); err == nil {
		query.RiskLevel = riskLevel
	}
	if pageNum, err := strconv.Atoi(c.Query("pageNum")); err == nil && pageNum > 0 {
		query.PageNum = pageNum
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}

	inspectionService := service.NewInspectionService()
	records, total, err := inspectionService.GetRecordList(c, query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	data := map[string]interface{}{
		"list":     records,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	}

	result.Success(c, data)
}
