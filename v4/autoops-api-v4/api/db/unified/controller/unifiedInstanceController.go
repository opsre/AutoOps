package controller

import (
	"dodevops-api/api/db/unified/model"
	"dodevops-api/api/db/unified/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Tags 数据库统一管理
// @Summary 获取所有类型数据库实例列表
// @Produce json
// @Description 统一查询 MySQL/PostgreSQL/Redis/MongoDB/Elasticsearch 实例
// @Param dbType query string false "数据库类型: mysql, postgres, redis, mongodb, elasticsearch"
// @Param name query string false "实例名称（模糊搜索）"
// @Param status query int false "状态: 1-启用, 2-禁用"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/db/instances/all [get]
// @Security ApiKeyAuth
func GetAllInstancesList(c *gin.Context) {
	query := &model.UnifiedInstanceQuery{
		DbType:   c.Query("dbType"),
		Name:     c.Query("name"),
		PageNum:  1,
		PageSize: 10,
	}

	// 解析分页参数
	if pageNum := c.Query("pageNum"); pageNum != "" {
		if num, err := strconv.Atoi(pageNum); err == nil && num > 0 {
			query.PageNum = num
		}
	}
	if pageSize := c.Query("pageSize"); pageSize != "" {
		if size, err := strconv.Atoi(pageSize); err == nil && size > 0 {
			query.PageSize = size
		}
	}

	// 解析状态参数
	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			query.Status = &status
		}
	}

	service.UnifiedInstanceService().GetUnifiedInstanceList(c, query)
}

// @Tags 数据库统一管理
// @Summary 获取实例详情
// @Produce json
// @Description 根据类型和ID获取实例详情
// @Param dbType path string true "数据库类型"
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/db/instances/{dbType}/{id} [get]
// @Security ApiKeyAuth
func GetInstanceById(c *gin.Context) {
	dbType := c.Param("dbType")
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	service.UnifiedInstanceService().GetUnifiedInstanceById(c, dbType, uint(id))
}

// @Tags 数据库统一管理
// @Summary 获取实例统计信息
// @Produce json
// @Description 获取各类型数据库实例数量统计
// @Success 200 {object} result.Result
// @Router /api/v1/db/instances/stats [get]
// @Security ApiKeyAuth
func GetInstancesStats(c *gin.Context) {
	service.UnifiedInstanceService().GetUnifiedInstanceStats(c)
}

// @Tags 数据库统一管理
// @Summary 获取数据库实例密码
// @Produce json
// @Description 根据数据库类型和ID获取实例的真实密码（敏感操作）
// @Param dbType path string true "数据库类型: mysql, postgres, redis, mongodb, elasticsearch"
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/db/unified/instances/{dbType}/{id}/password [get]
// @Security ApiKeyAuth
func GetInstancePassword(c *gin.Context) {
	dbType := c.Param("dbType")
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	service.UnifiedInstanceService().GetInstancePassword(c, dbType, uint(id))
}
