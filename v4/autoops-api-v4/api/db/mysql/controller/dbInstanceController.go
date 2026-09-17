package controller

import (
	"dodevops-api/api/db/mysql/model"
	"dodevops-api/api/db/mysql/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Tags MySQL数据库管理
// @Summary 获取数据库实例列表
// @Produce json
// @Description 获取数据库实例列表
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param name query string false "实例名称"
// @Param type query string false "数据库类型"
// @Param host query string false "主机地址"
// @Param status query int false "状态"
// @Success 200 {object} result.Result
// @Router /api/v1/db/instances [get]
// @Security ApiKeyAuth
func GetDbInstanceList(c *gin.Context) {
	query := &model.DbInstanceQuery{
		Name:     c.Query("name"),
		Type:     c.Query("type"),
		Host:     c.Query("host"),
		PageNum:  1,
		PageSize: 10,
	}

	if pageNum, err := strconv.Atoi(c.Query("pageNum")); err == nil && pageNum > 0 {
		query.PageNum = pageNum
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}
	if status, err := strconv.Atoi(c.Query("status")); err == nil {
		query.Status = status
	}

	service.DbInstanceService().GetDbInstanceList(c, query)
}

// @Tags MySQL数据库管理
// @Summary 根据ID获取数据库实例
// @Produce json
// @Description 根据ID获取数据库实例
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/db/instances/{id} [get]
// @Security ApiKeyAuth
func GetDbInstanceById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.DbInstanceService().GetDbInstanceById(c, uint(id))
}

// @Tags MySQL数据库管理
// @Summary 创建数据库实例
// @Produce json
// @Description 创建数据库实例
// @Param data body model.DbInstanceDto true "实例信息"
// @Success 200 {object} result.Result
// @Router /api/v1/db/instances [post]
// @Security ApiKeyAuth
func CreateDbInstance(c *gin.Context) {
	var dto model.DbInstanceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.DbInstanceService().CreateDbInstance(c, &dto)
}

// @Tags MySQL数据库管理
// @Summary 更新数据库实例
// @Produce json
// @Description 更新数据库实例
// @Param data body model.DbInstanceDto true "实例信息"
// @Success 200 {object} result.Result
// @Router /api/v1/db/instances [put]
// @Security ApiKeyAuth
func UpdateDbInstance(c *gin.Context) {
	var dto model.DbInstanceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.DbInstanceService().UpdateDbInstance(c, &dto)
}

// @Tags MySQL数据库管理
// @Summary 删除数据库实例
// @Produce json
// @Description 删除数据库实例
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/db/instances/{id} [delete]
// @Security ApiKeyAuth
func DeleteDbInstance(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.DbInstanceService().DeleteDbInstance(c, uint(id))
}

// @Tags MySQL数据库管理
// @Summary 测试数据库连接
// @Produce json
// @Description 测试数据库连接
// @Param data body model.TestConnDto true "连接信息"
// @Success 200 {object} result.Result
// @Router /api/v1/db/instances/test-conn [post]
// @Security ApiKeyAuth
func TestConnection(c *gin.Context) {
	var dto model.TestConnDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.DbInstanceService().TestConnection(c, &dto)
}

// @Tags MySQL数据库管理
// @Summary 获取实例的数据库列表
// @Produce json
// @Description 获取实例的所有数据库
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/db/instances/{id}/databases [get]
// @Security ApiKeyAuth
func GetDatabases(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.DbInstanceService().GetDatabases(c, uint(id))
}
