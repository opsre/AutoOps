package controller

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Tags 故障管理
// @Summary 创建故障记录
// @Description 创建故障记录
// @Produce json
// @Param data body model.CreateIncidentDto true "故障信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/incident [post]
// @Security ApiKeyAuth
func CreateIncident(c *gin.Context) {
	var dto model.CreateIncidentDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	service.GetIncidentService().CreateIncident(c, &dto)
}

// @Tags 故障管理
// @Summary 更新故障记录
// @Description 更新故障记录
// @Produce json
// @Param data body model.UpdateIncidentDto true "故障信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/incident [put]
// @Security ApiKeyAuth
func UpdateIncident(c *gin.Context) {
	var dto model.UpdateIncidentDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	service.GetIncidentService().UpdateIncident(c, &dto)
}

// @Tags 故障管理
// @Summary 删除故障记录
// @Description 删除故障记录
// @Produce json
// @Param id path int true "故障记录ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/incident/{id} [delete]
// @Security ApiKeyAuth
func DeleteIncident(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.GetIncidentService().DeleteIncident(c, uint(id))
}

// @Tags 故障管理
// @Summary 根据ID获取故障记录
// @Description 根据ID获取故障记录
// @Produce json
// @Param id path int true "故障记录ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/incident/{id} [get]
// @Security ApiKeyAuth
func GetIncidentById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.GetIncidentService().GetIncidentById(c, uint(id))
}

// @Tags 故障管理
// @Summary 获取故障列表
// @Description 获取故障列表
// @Produce json
// @Param businessLine query string false "业务线"
// @Param alertLevel query string false "告警级别"
// @Param status query int false "处理状态"
// @Param department query string false "所属部门"
// @Param handler query string false "处理人"
// @Param startTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/incident/list [get]
// @Security ApiKeyAuth
func GetIncidentList(c *gin.Context) {
	query := &model.IncidentQueryDto{}

	if v := c.Query("businessLineId"); v != "" {
		if businessLineId, _ := strconv.Atoi(v); businessLineId > 0 {
			query.BusinessLineID = uint(businessLineId)
		}
	}
	if v := c.Query("alertLevel"); v != "" {
		query.AlertLevel = v
	}
	if v := c.Query("status"); v != "" {
		if status, _ := strconv.Atoi(v); status > 0 {
			query.Status = status
		}
	}
	if v := c.Query("department"); v != "" {
		query.Department = v
	}
	if v := c.Query("handler"); v != "" {
		query.Handler = v
	}
	if v := c.Query("startTime"); v != "" {
		query.StartTime = v
	}
	if v := c.Query("endTime"); v != "" {
		query.EndTime = v
	}
	if v := c.Query("pageNum"); v != "" {
		if pageNum, _ := strconv.Atoi(v); pageNum > 0 {
			query.PageNum = pageNum
		}
	}
	if v := c.Query("pageSize"); v != "" {
		if pageSize, _ := strconv.Atoi(v); pageSize > 0 {
			query.PageSize = pageSize
		}
	}

	service.GetIncidentService().GetIncidentList(c, query)
}

// @Tags 故障管理
// @Summary 批量删除故障记录
// @Description 批量删除故障记录
// @Produce json
// @Param data body []uint true "故障记录ID列表"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/incident/batch-delete [post]
// @Security ApiKeyAuth
func BatchDeleteIncidents(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	service.GetIncidentService().BatchDeleteIncidents(c, ids)
}

// @Tags 故障管理
// @Summary 更新故障状态
// @Description 更新故障状态
// @Produce json
// @Param id path int true "故障记录ID"
// @Param status query int true "处理状态(1:未处理,2:处理中,3:已完成)"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/incident/{id}/status [put]
// @Security ApiKeyAuth
func UpdateIncidentStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	status, _ := strconv.Atoi(c.Query("status"))
	service.GetIncidentService().UpdateIncidentStatus(c, uint(id), status)
}

// @Tags 故障管理
// @Summary 获取故障统计信息
// @Description 获取故障统计信息
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/incident/statistics [get]
// @Security ApiKeyAuth
func GetIncidentStatistics(c *gin.Context) {
	service.GetIncidentService().GetIncidentStatistics(c)
}
