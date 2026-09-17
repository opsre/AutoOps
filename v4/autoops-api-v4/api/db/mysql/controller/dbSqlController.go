package controller

import (
	"dodevops-api/api/db/mysql/model"
	"dodevops-api/api/db/mysql/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Tags MySQL数据库管理
// @Summary 获取SQL模板列表
// @Produce json
// @Description 获取SQL模板列表
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param dbId query int false "数据库ID"
// @Param db query string false "数据库名"
// @Param name query string false "模板名称"
// @Success 200 {object} result.Result
// @Router /api/v1/db/sql-templates [get]
// @Security ApiKeyAuth
func GetDbSqlList(c *gin.Context) {
	query := &model.DbSqlQuery{
		Db:       c.Query("db"),
		Name:     c.Query("name"),
		PageNum:  1,
		PageSize: 10,
	}

	if pageNum, err := strconv.Atoi(c.Query("pageNum")); err == nil && pageNum > 0 {
		query.PageNum = pageNum
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}
	if dbId, err := strconv.ParseUint(c.Query("dbId"), 10, 64); err == nil {
		query.DbId = uint(dbId)
	}

	service.DbSqlService().GetDbSqlList(c, query)
}

// @Tags MySQL数据库管理
// @Summary 根据ID获取SQL模板
// @Produce json
// @Description 根据ID获取SQL模板
// @Param id path int true "模板ID"
// @Success 200 {object} result.Result
// @Router /api/v1/db/sql-templates/{id} [get]
// @Security ApiKeyAuth
func GetDbSqlById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.DbSqlService().GetDbSqlById(c, uint(id))
}

// @Tags MySQL数据库管理
// @Summary 获取SQL模板名称列表
// @Produce json
// @Description 获取指定数据库的SQL模板名称列表
// @Param dbId query int true "数据库ID"
// @Param db query string true "数据库名"
// @Success 200 {object} result.Result
// @Router /api/v1/db/sql-templates/names [get]
// @Security ApiKeyAuth
func GetDbSqlNames(c *gin.Context) {
	dbId, _ := strconv.ParseUint(c.Query("dbId"), 10, 64)
	db := c.Query("db")
	service.DbSqlService().GetDbSqlNames(c, uint(dbId), db)
}

// @Tags MySQL数据库管理
// @Summary 创建SQL模板
// @Produce json
// @Description 创建SQL模板
// @Param data body model.DbSqlDto true "模板信息"
// @Success 200 {object} result.Result
// @Router /api/v1/db/sql-templates [post]
// @Security ApiKeyAuth
func CreateDbSql(c *gin.Context) {
	var dto model.DbSqlDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.DbSqlService().CreateDbSql(c, &dto)
}

// @Tags MySQL数据库管理
// @Summary 更新SQL模板
// @Produce json
// @Description 更新SQL模板
// @Param data body model.DbSqlDto true "模板信息"
// @Success 200 {object} result.Result
// @Router /api/v1/db/sql-templates [put]
// @Security ApiKeyAuth
func UpdateDbSql(c *gin.Context) {
	var dto model.DbSqlDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.DbSqlService().UpdateDbSql(c, &dto)
}

// @Tags MySQL数据库管理
// @Summary 删除SQL模板
// @Produce json
// @Description 删除SQL模板
// @Param id path int true "模板ID"
// @Success 200 {object} result.Result
// @Router /api/v1/db/sql-templates/{id} [delete]
// @Security ApiKeyAuth
func DeleteDbSql(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.DbSqlService().DeleteDbSql(c, uint(id))
}

// @Tags MySQL数据库管理
// @Summary 获取SQL执行历史列表
// @Produce json
// @Description 获取SQL执行历史列表
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param dbId query int false "数据库ID"
// @Param dbName query string false "数据库名"
// @Param tableName query string false "表名"
// @Param type query int false "SQL类型"
// @Param status query int false "执行状态"
// @Param creator query string false "创建人"
// @Param startTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} result.Result
// @Router /api/v1/db/sql-execs [get]
// @Security ApiKeyAuth
func GetDbSqlExecList(c *gin.Context) {
	query := &model.DbSqlExecQuery{
		DbName:    c.Query("dbName"),
		TableName: c.Query("tableName"),
		Creator:   c.Query("creator"),
		StartTime: c.Query("startTime"),
		EndTime:   c.Query("endTime"),
		PageNum:   1,
		PageSize:  10,
	}

	if pageNum, err := strconv.Atoi(c.Query("pageNum")); err == nil && pageNum > 0 {
		query.PageNum = pageNum
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}
	if dbId, err := strconv.ParseUint(c.Query("dbId"), 10, 64); err == nil {
		query.DbId = uint(dbId)
	}
	if sqlType, err := strconv.Atoi(c.Query("type")); err == nil {
		query.Type = int8(sqlType)
	}
	if status, err := strconv.Atoi(c.Query("status")); err == nil {
		query.Status = int8(status)
	}

	service.DbSqlExecService().GetDbSqlExecList(c, query)
}

// @Tags MySQL数据库管理
// @Summary 根据ID获取SQL执行历史
// @Produce json
// @Description 根据ID获取SQL执行历史
// @Param id path int true "记录ID"
// @Success 200 {object} result.Result
// @Router /api/v1/db/sql-execs/{id} [get]
// @Security ApiKeyAuth
func GetDbSqlExecById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.DbSqlExecService().GetDbSqlExecById(c, uint(id))
}

// @Tags MySQL数据库管理
// @Summary 获取SQL执行统计
// @Produce json
// @Description 获取SQL执行统计信息
// @Param dbId query int true "数据库ID"
// @Param startTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} result.Result
// @Router /api/v1/db/sql-execs/statistics [get]
// @Security ApiKeyAuth
func GetDbSqlExecStatistics(c *gin.Context) {
	dbId, _ := strconv.ParseUint(c.Query("dbId"), 10, 64)
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	service.DbSqlExecService().GetDbSqlExecStatistics(c, uint(dbId), startTime, endTime)
}
