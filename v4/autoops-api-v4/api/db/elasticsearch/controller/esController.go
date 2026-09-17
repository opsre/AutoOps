package controller

import (
    "dodevops-api/api/db/elasticsearch/model"
    "dodevops-api/api/db/elasticsearch/service"
    "dodevops-api/common/result"
    "github.com/gin-gonic/gin"
    "strconv"
)

// @Tags Elasticsearch管理
// @Summary 获取Elasticsearch实例列表
// @Produce json
// @Description 获取Elasticsearch实例列表
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param name query string false "实例名称"
// @Param host query string false "主机地址"
// @Param status query int false "状态"
// @Success 200 {object} result.Result
// @Router /api/v1/elasticsearch [get]
// @Security ApiKeyAuth
func GetEsInstanceList(c *gin.Context) {
    q := &model.EsInstanceQuery{PageNum: 1, PageSize: 10}
    if v := c.Query("name"); v != "" { q.Name = v }
    if v := c.Query("host"); v != "" { q.Host = v }
    if v := c.Query("status"); v != "" { if n, _ := strconv.Atoi(v); n >= 0 { q.Status = n } }
    if v := c.Query("pageNum"); v != "" { if n, _ := strconv.Atoi(v); n > 0 { q.PageNum = n } }
    if v := c.Query("pageSize"); v != "" { if n, _ := strconv.Atoi(v); n > 0 { q.PageSize = n } }
    service.EsInstanceService().GetEsInstanceList(c, q)
}

// @Tags Elasticsearch管理
// @Summary 根据ID获取Elasticsearch实例
// @Produce json
// @Description 根据ID获取Elasticsearch实例
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/elasticsearch/{id} [get]
// @Security ApiKeyAuth
func GetEsInstanceById(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    service.EsInstanceService().GetEsInstanceById(c, uint(id))
}

// @Tags Elasticsearch管理
// @Summary 创建Elasticsearch实例
// @Produce json
// @Description 创建Elasticsearch实例
// @Param data body model.EsInstanceDto true "实例信息"
// @Success 200 {object} result.Result
// @Router /api/v1/elasticsearch [post]
// @Security ApiKeyAuth
func CreateEsInstance(c *gin.Context) {
    var dto model.EsInstanceDto
    if err := c.ShouldBindJSON(&dto); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    service.EsInstanceService().SaveEsInstance(c, &dto)
}

// @Tags Elasticsearch管理
// @Summary 更新Elasticsearch实例
// @Produce json
// @Description 更新Elasticsearch实例
// @Param data body model.EsInstanceDto true "实例信息"
// @Success 200 {object} result.Result
// @Router /api/v1/elasticsearch [put]
// @Security ApiKeyAuth
func UpdateEsInstance(c *gin.Context) {
    var dto model.EsInstanceDto
    if err := c.ShouldBindJSON(&dto); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    service.EsInstanceService().SaveEsInstance(c, &dto)
}

// @Tags Elasticsearch管理
// @Summary 删除Elasticsearch实例
// @Produce json
// @Description 删除Elasticsearch实例
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/elasticsearch/{id} [delete]
// @Security ApiKeyAuth
func DeleteEsInstance(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    service.EsInstanceService().DeleteEsInstance(c, uint(id))
}

// @Tags Elasticsearch管理
// @Summary 测试Elasticsearch连接
// @Produce json
// @Description 测试Elasticsearch连接
// @Param data body model.EsInstanceDto true "连接信息(协议/主机/端口/用户名/密码)"
// @Success 200 {object} result.Result
// @Router /api/v1/elasticsearch/test-conn [post]
// @Security ApiKeyAuth
func TestConnection(c *gin.Context) {
    var dto model.EsInstanceDto
    if err := c.ShouldBindJSON(&dto); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    inst := &model.EsInstance{Protocol: dto.Protocol, Host: dto.Host, Port: dto.Port, Username: dto.Username, Password: dto.Password}
    service.EsService().TestConn(c, inst)
}

// @Tags Elasticsearch管理
// @Summary 反向代理ES API
// @Produce json
// @Description 反向代理至目标Elasticsearch接口
// @Param id path int true "实例ID"
// @Param path path string false "目标路径(支持通配)"
// @Success 200 {object} result.Result
// @Router /api/v1/elasticsearch/{id}/proxy/{path} [get]
// @Security ApiKeyAuth
func Proxy(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    path := c.Param("path")
    service.EsService().ProxyById(c, uint(id), path)
}
