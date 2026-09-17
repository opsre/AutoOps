package controller

import (
	"dodevops-api/api/db/mongo/model"
	"dodevops-api/api/db/mongo/service"
	"dodevops-api/common/result"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Tags MongoDB管理
// @Summary 获取MongoDB实例列表
// @Produce json
// @Description 获取MongoDB实例列表
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param name query string false "实例名称"
// @Param status query int false "状态"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb [get]
// @Security ApiKeyAuth
func GetMongoInstanceList(c *gin.Context) {
	q := &model.MongoInstanceQuery{PageNum: 1, PageSize: 10}
	if v := c.Query("name"); v != "" {
		q.Name = v
	}
	if v := c.Query("status"); v != "" {
		if n, _ := strconv.Atoi(v); n >= 0 {
			q.Status = n
		}
	}
	if v := c.Query("pageNum"); v != "" {
		if n, _ := strconv.Atoi(v); n > 0 {
			q.PageNum = n
		}
	}
	if v := c.Query("pageSize"); v != "" {
		if n, _ := strconv.Atoi(v); n > 0 {
			q.PageSize = n
		}
	}
	service.MongoInstanceService().GetMongoInstanceList(c, q)
}

// @Tags MongoDB管理
// @Summary 根据ID获取MongoDB实例
// @Produce json
// @Description 根据ID获取MongoDB实例
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb/{id} [get]
// @Security ApiKeyAuth
func GetMongoInstanceById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.MongoInstanceService().GetMongoInstanceById(c, uint(id))
}

// @Tags MongoDB管理
// @Summary 创建MongoDB实例
// @Produce json
// @Description 创建MongoDB实例
// @Param data body model.MongoInstanceDto true "实例信息"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb [post]
// @Security ApiKeyAuth
func CreateMongoInstance(c *gin.Context) {
	var dto model.MongoInstanceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	service.MongoInstanceService().SaveMongoInstance(c, &dto)
}

// @Tags MongoDB管理
// @Summary 更新MongoDB实例
// @Produce json
// @Description 更新MongoDB实例
// @Param data body model.MongoInstanceDto true "实例信息"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb [put]
// @Security ApiKeyAuth
func UpdateMongoInstance(c *gin.Context) {
	var dto model.MongoInstanceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	service.MongoInstanceService().SaveMongoInstance(c, &dto)
}

// @Tags MongoDB管理
// @Summary 删除MongoDB实例
// @Produce json
// @Description 删除MongoDB实例
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb/{id} [delete]
// @Security ApiKeyAuth
func DeleteMongoInstance(c *gin.Context) {
	idsStr := c.Param("id")
	ids := strings.Split(idsStr, ",")

	if len(ids) == 1 {
		id, _ := strconv.ParseUint(ids[0], 10, 64)
		service.MongoInstanceService().DeleteMongoInstance(c, uint(id))
	} else {
		var idList []uint
		for _, idStr := range ids {
			if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
				idList = append(idList, uint(id))
			}
		}
		service.MongoInstanceService().BatchDeleteMongoInstances(c, idList)
	}
}

// @Tags MongoDB管理
// @Summary 测试MongoDB连接
// @Produce json
// @Description 测试MongoDB连接
// @Param data body model.MongoInstanceDto true "连接信息(URI)"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb/test-conn [post]
// @Security ApiKeyAuth
func TestConnection(c *gin.Context) {
	var dto model.MongoInstanceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	service.MongoService().TestConn(c, &dto)
}

// @Tags MongoDB管理
// @Summary 获取所有数据库
// @Produce json
// @Description 获取MongoDB实例下所有数据库
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb/{id}/databases [get]
// @Security ApiKeyAuth
func GetDatabases(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.MongoService().GetDatabases(c, uint(id))
}

// @Tags MongoDB管理
// @Summary 获取指定库的所有集合
// @Produce json
// @Description 获取MongoDB指定数据库的所有集合
// @Param id path int true "实例ID"
// @Param database query string true "数据库名称"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb/{id}/collections [get]
// @Security ApiKeyAuth
func GetCollections(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	database := c.Query("database")
	service.MongoService().GetCollections(c, uint(id), database)
}

// @Tags MongoDB操作
// @Summary 执行MongoDB命令
// @Produce json
// @Description 执行MongoDB RunCommand
// @Param id path int true "实例ID"
// @Param data body model.MongoRunCommandForm true "命令信息"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb/{id}/run-command [post]
// @Security ApiKeyAuth
func RunCommand(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var form model.MongoRunCommandForm
	if err := c.ShouldBindJSON(&form); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	service.MongoService().RunCommand(c, uint(id), &form)
}

// @Tags MongoDB操作
// @Summary 查询文档
// @Produce json
// @Description 执行MongoDB Find命令
// @Param id path int true "实例ID"
// @Param data body model.MongoFindForm true "查询条件"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb/{id}/command/find [post]
// @Security ApiKeyAuth
func Find(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var form model.MongoFindForm
	if err := c.ShouldBindJSON(&form); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	service.MongoService().Find(c, uint(id), &form)
}

// @Tags MongoDB操作
// @Summary 根据ID更新文档
// @Produce json
// @Description 执行MongoDB UpdateByID命令
// @Param id path int true "实例ID"
// @Param data body model.MongoUpdateByIdForm true "更新信息"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb/{id}/command/update-by-id [post]
// @Security ApiKeyAuth
func UpdateById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var form model.MongoUpdateByIdForm
	if err := c.ShouldBindJSON(&form); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	service.MongoService().UpdateById(c, uint(id), &form)
}

// @Tags MongoDB操作
// @Summary 根据ID删除文档
// @Produce json
// @Description 执行MongoDB DeleteOne命令
// @Param id path int true "实例ID"
// @Param data body model.MongoDeleteByIdForm true "删除信息"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb/{id}/command/delete-by-id [post]
// @Security ApiKeyAuth
func DeleteById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var form model.MongoDeleteByIdForm
	if err := c.ShouldBindJSON(&form); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	service.MongoService().DeleteById(c, uint(id), &form)
}

// @Tags MongoDB操作
// @Summary 插入文档
// @Produce json
// @Description 执行MongoDB InsertOne命令
// @Param id path int true "实例ID"
// @Param data body model.MongoInsertForm true "文档数据"
// @Success 200 {object} result.Result
// @Router /api/v1/mongodb/{id}/command/insert [post]
// @Security ApiKeyAuth
func InsertOne(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var form model.MongoInsertForm
	if err := c.ShouldBindJSON(&form); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	service.MongoService().InsertOne(c, uint(id), &form)
}
