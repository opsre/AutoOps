package controller

import (
    "dodevops-api/api/db/redis/dao"
    "dodevops-api/api/db/redis/model"
    "dodevops-api/api/db/redis/service"
    "dodevops-api/common/result"
    "github.com/gin-gonic/gin"
    "strconv"
)

// @Tags Redis管理
// @Summary 获取Redis实例列表
// @Produce json
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param name query string false "实例名称"
// @Param host query string false "主机地址"
// @Param mode query string false "模式(standalone|cluster|sentinel)"
// @Param status query int false "状态"
// @Success 200 {object} result.Result
// @Router /api/v1/redis [get]
// @Security ApiKeyAuth
func GetRedisInstanceList(c *gin.Context) {
    q := &model.RedisInstanceQuery{PageNum: 1, PageSize: 10}
    if v := c.Query("name"); v != "" { q.Name = v }
    if v := c.Query("host"); v != "" { q.Host = v }
    if v := c.Query("mode"); v != "" { q.Mode = v }
    if v := c.Query("status"); v != "" { if n, _ := strconv.Atoi(v); n >= 0 { q.Status = n } }
    if v := c.Query("pageNum"); v != "" { if n, _ := strconv.Atoi(v); n > 0 { q.PageNum = n } }
    if v := c.Query("pageSize"); v != "" { if n, _ := strconv.Atoi(v); n > 0 { q.PageSize = n } }
    service.RedisInstanceService().GetRedisInstanceList(c, q)
}

// @Tags Redis管理
// @Summary 根据ID获取Redis实例
// @Produce json
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/{id} [get]
// @Security ApiKeyAuth
func GetRedisInstanceById(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    service.RedisInstanceService().GetRedisInstanceById(c, uint(id))
}

// @Tags Redis管理
// @Summary 新增/更新Redis实例
// @Produce json
// @Param data body model.RedisInstanceDto true "实例信息"
// @Success 200 {object} result.Result
// @Router /api/v1/redis [post]
// @Router /api/v1/redis [put]
// @Security ApiKeyAuth
func SaveRedisInstance(c *gin.Context) {
    var dto model.RedisInstanceDto
    if err := c.ShouldBindJSON(&dto); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    service.RedisInstanceService().SaveRedisInstance(c, &dto)
}

// @Tags Redis管理
// @Summary 删除Redis实例
// @Produce json
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/{id} [delete]
// @Security ApiKeyAuth
func DeleteRedisInstance(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    service.RedisInstanceService().DeleteRedisInstance(c, uint(id))
}

// @Tags Redis管理
// @Summary 测试连接
// @Produce json
// @Param data body model.Conn true "连接信息"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/test-conn [post]
// @Security ApiKeyAuth
func TestConnection(c *gin.Context) {
    var conn model.Conn
    if err := c.ShouldBindJSON(&conn); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    service.RedisService().TestConn(c, &conn)
}

// ---------- 基于实例ID/DB的操作 ----------

func buildConnFromId(c *gin.Context) (*model.Conn, bool) {
    id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil || id64 == 0 { result.Failed(c, int(result.ApiCode.FAILED), "非法实例ID"); return nil, false }
    inst, err2 := dao.GetRedisInstanceById(uint(id64))
    if err2 != nil { result.Failed(c, int(result.ApiCode.FAILED), "实例不存在"); return nil, false }
    dbNum := 0
    if v := c.Param("db"); v != "" { if n, _ := strconv.Atoi(v); n >= 0 { dbNum = n } }
    conn := &model.Conn{ Host: inst.Host, Port: inst.Port, Password: inst.Password, DB: dbNum, Mode: inst.Mode, RedisNodePassword: inst.RedisNodePassword }
    return conn, true
}

// @Tags Redis管理
// @Summary 获取Redis INFO
// @Produce json
// @Param id path int true "实例ID"
// @Param format query string false "map 为解析后的结构"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/{id}/info [get]
// @Security ApiKeyAuth
func GetRedisInfo(c *gin.Context) {
    conn, ok := buildConnFromId(c); if !ok { return }
    service.RedisService().Info(c, conn)
}

// @Tags Redis管理
// @Summary 获取集群信息
// @Produce json
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/{id}/cluster-info [get]
// @Security ApiKeyAuth
func GetRedisClusterInfo(c *gin.Context) {
    conn, ok := buildConnFromId(c); if !ok { return }
    service.RedisService().ClusterInfo(c, conn)
}

// @Tags Redis管理
// @Summary 扫描Keys
// @Produce json
// @Param id path int true "实例ID"
// @Param db path int true "数据库编号"
// @Param data body model.ScanReq true "扫描参数"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/{id}/{db}/scan [post]
// @Security ApiKeyAuth
func ScanRedisKeys(c *gin.Context) {
    conn, ok := buildConnFromId(c); if !ok { return }
    // 仅绑定非连接字段，连接信息来自实例
    type scanReqBody struct {
        Match     string            `json:"match"`
        Pattern   string            `json:"pattern"`
        Count     int               `json:"count"`
        CursorMap map[string]uint64 `json:"cursor"`
        Cursor    uint64            `json:"cursorSingle"`
    }
    var body scanReqBody
    if err := c.ShouldBindJSON(&body); err != nil && err.Error() != "EOF" { // 空体也允许
        result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return
    }
    req := model.ScanReq{Conn: *conn, Match: body.Match, Pattern: body.Pattern, Count: body.Count, CursorMap: body.CursorMap, Cursor: body.Cursor}
    service.RedisService().Scan(c, &req)
}

// @Tags Redis管理
// @Summary 扫描全部 DB 的 Keys（无需选择DB）
// @Produce json
// @Param id path int true "实例ID"
// @Param data body model.ScanReq true "扫描参数（可传cursor为db:0等键）"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/{id}/scan [post]
// @Security ApiKeyAuth
func ScanRedisKeysAllDb(c *gin.Context) {
    conn, ok := buildConnFromId(c); if !ok { return }
    // 仅绑定非连接字段，连接信息来自实例
    type scanReqBody struct {
        Match     string            `json:"match"`
        Pattern   string            `json:"pattern"`
        Count     int               `json:"count"`
        CursorMap map[string]uint64 `json:"cursor"`
        Cursor    uint64            `json:"cursorSingle"`
    }
    var body scanReqBody
    if err := c.ShouldBindJSON(&body); err != nil && err.Error() != "EOF" { // 空体也允许
        result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return
    }
    req := model.ScanReq{Conn: *conn, Match: body.Match, Pattern: body.Pattern, Count: body.Count, CursorMap: body.CursorMap, Cursor: body.Cursor}
    service.RedisService().ScanAllDb(c, &req)
}

// @Tags Redis管理
// @Summary Key详情
// @Produce json
// @Param id path int true "实例ID"
// @Param db path int true "数据库编号"
// @Param key query string true "Key"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/{id}/{db}/key-info [get]
// @Security ApiKeyAuth
func GetRedisKeyInfo(c *gin.Context) {
    conn, ok := buildConnFromId(c); if !ok { return }
    key := c.Query("key"); if key == "" { result.Failed(c, int(result.ApiCode.FAILED), "key不能为空"); return }
    service.RedisService().KeyInfo(c, conn, key)
}

// @Tags Redis管理
// @Summary TTL信息
// @Produce json
// @Param id path int true "实例ID"
// @Param db path int true "数据库编号"
// @Param key query string true "Key"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/{id}/{db}/key-ttl [get]
// @Security ApiKeyAuth
func GetRedisKeyTtl(c *gin.Context) {
    conn, ok := buildConnFromId(c); if !ok { return }
    key := c.Query("key"); if key == "" { result.Failed(c, int(result.ApiCode.FAILED), "key不能为空"); return }
    service.RedisService().TtlKey(c, conn, key)
}

// @Tags Redis管理
// @Summary Key内存使用
// @Produce json
// @Param id path int true "实例ID"
// @Param db path int true "数据库编号"
// @Param key query string true "Key"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/{id}/{db}/key-memuse [get]
// @Security ApiKeyAuth
func GetRedisKeyMemuse(c *gin.Context) {
    conn, ok := buildConnFromId(c); if !ok { return }
    key := c.Query("key"); if key == "" { result.Failed(c, int(result.ApiCode.FAILED), "key不能为空"); return }
    service.RedisService().MemoryUsage(c, conn, key)
}

// @Tags Redis管理
// @Summary 执行命令
// @Produce json
// @Param id path int true "实例ID"
// @Param db path int true "数据库编号"
// @Param data body model.CmdReq true "命令参数"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/{id}/{db}/run-cmd [post]
// @Security ApiKeyAuth
func RunRedisCmd(c *gin.Context) {
    conn, ok := buildConnFromId(c); if !ok { return }
    type cmdReqBody struct {
        Command string   `json:"command" binding:"required"`
        Args    []string `json:"args"`
    }
    var body cmdReqBody
    if err := c.ShouldBindJSON(&body); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    req := model.CmdReq{Conn: *conn, Command: body.Command, Args: body.Args}
    service.RedisService().RunCmd(c, &req)
}

// @Tags Redis管理
// @Summary 列出实例的DB列表
// @Produce json
// @Param id path int true "实例ID"
// @Success 200 {object} result.Result
// @Router /api/v1/redis/{id}/dbs [get]
// @Security ApiKeyAuth
func GetRedisDatabases(c *gin.Context) {
    conn, ok := buildConnFromId(c); if !ok { return }
    service.RedisService().ListDbs(c, conn)
}
