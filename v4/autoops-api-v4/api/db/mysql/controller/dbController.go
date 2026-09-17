package controller

import (
    "dodevops-api/api/db/mysql/model"
    "dodevops-api/api/db/mysql/service"
    "dodevops-api/api/db/mysql/dao"
    "dodevops-api/pkg/utils"
    "github.com/gin-gonic/gin"
    "strconv"
    "io"
    "strings"
)

// @Tags MySQL数据库管理
// @Summary 获取数据库列表
// @Produce json
// @Description 获取数据库列表
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param name query string false "数据库名称"
// @Param database query string false "数据库名"
// @Param instanceId query int false "实例ID"
// @Param status query int false "状态"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases [get]
// @Security ApiKeyAuth
func GetDbList(c *gin.Context) {
	query := &model.DbQuery{
		Name:     c.Query("name"),
		Database: c.Query("database"),
		PageNum:  1,
		PageSize: 10,
	}

	if pageNum, err := strconv.Atoi(c.Query("pageNum")); err == nil && pageNum > 0 {
		query.PageNum = pageNum
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}
	if instanceId, err := strconv.ParseUint(c.Query("instanceId"), 10, 64); err == nil {
		query.InstanceId = uint(instanceId)
	}
	if status, err := strconv.Atoi(c.Query("status")); err == nil {
		query.Status = status
	}

	service.DbService().GetDbList(c, query)
}

// @Tags MySQL数据库管理
// @Summary 根据ID获取数据库
// @Produce json
// @Description 根据ID获取数据库
// @Param id path int true "数据库ID"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id} [get]
// @Security ApiKeyAuth
func GetDbById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.DbService().GetDbById(c, uint(id))
}

// @Tags MySQL数据库管理
// @Summary 创建数据库
// @Produce json
// @Description 创建数据库
// @Param data body model.DbDto true "数据库信息"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases [post]
// @Security ApiKeyAuth
func CreateDb(c *gin.Context) {
	var dto model.DbDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.DbService().CreateDb(c, &dto)
}

// @Tags MySQL数据库管理
// @Summary 更新数据库
// @Produce json
// @Description 更新数据库
// @Param data body model.DbDto true "数据库信息"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases [put]
// @Security ApiKeyAuth
func UpdateDb(c *gin.Context) {
	var dto model.DbDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.DbService().UpdateDb(c, &dto)
}

// @Tags MySQL数据库管理
// @Summary 删除数据库
// @Produce json
// @Description 删除数据库
// @Param id path int true "数据库ID"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id} [delete]
// @Security ApiKeyAuth
func DeleteDb(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.DbService().DeleteDb(c, uint(id))
}

// @Tags MySQL数据库管理
// @Summary 获取表列表
// @Produce json
// @Description 获取数据库的表列表
// @Param id path int true "数据库ID"
// @Param db query string true "数据库名"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/tables [get]
// @Security ApiKeyAuth
func GetTables(c *gin.Context) {
    // 保持与其余接口一致，路径中的 :id 为数据库表 db 的主键ID
    dbId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("dbName")
    if dbName == "" {
        dbName = c.Query("db")
    }
    // 通过数据库ID获取连接，避免与 Columns/Indexes/DDL 等接口语义不一致
    service.DbService().GetTables(c, uint(dbId), dbName)
}

// @Tags MySQL数据库管理
// @Summary 获取表列信息
// @Produce json
// @Description 获取表的列信息
// @Param id path int true "数据库ID"
// @Param db query string true "数据库名"
// @Param table query string true "表名"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/columns [get]
// @Security ApiKeyAuth
func GetTableColumns(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	dbName := c.Query("db")
	tableName := c.Query("table")
	service.DbService().GetTableColumns(c, uint(id), dbName, tableName)
}

// @Tags MySQL数据库管理
// @Summary 获取表索引信息
// @Produce json
// @Description 获取表的索引信息
// @Param id path int true "数据库ID"
// @Param db query string true "数据库名"
// @Param table query string true "表名"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/indexes [get]
// @Security ApiKeyAuth
func GetTableIndexes(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	dbName := c.Query("db")
	tableName := c.Query("table")
	service.DbService().GetTableIndexes(c, uint(id), dbName, tableName)
}

// @Tags MySQL数据库管理
// @Summary 获取表DDL
// @Produce json
// @Description 获取表的创建DDL语句
// @Param id path int true "数据库ID"
// @Param db query string true "数据库名"
// @Param table query string true "表名"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/ddl [get]
// @Security ApiKeyAuth
func GetTableDDL(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	dbName := c.Query("db")
	tableName := c.Query("table")
	service.DbService().GetTableDDL(c, uint(id), dbName, tableName)
}

// @Tags MySQL数据库管理
// @Summary 复制表
// @Produce json
// @Description 复制表结构和数据
// @Param data body model.CopyTableDto true "复制参数"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/copy-table [post]
// @Security ApiKeyAuth
func CopyTable(c *gin.Context) {
	var dto model.CopyTableDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.DbService().CopyTable(c, &dto)
}

// @Tags MySQL数据库管理
// @Summary 执行SQL
// @Produce json
// @Description 执行SQL语句
// @Param data body model.DbSqlExecDto true "SQL信息"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/exec-sql [post]
// @Security ApiKeyAuth
func ExecSql(c *gin.Context) {
	var dto model.DbSqlExecDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.DbService().ExecSql(c, &dto)
}

// @Tags MySQL数据库管理
// @Summary 执行SQL文件
// @Produce json
// @Description 上传SQL脚本并顺序执行（基础分号分割，复杂脚本建议拆分）
// @Accept multipart/form-data
// @Param id path int true "数据库ID或实例ID(兼容)"
// @Param db formData string true "数据库名"
// @Param file formData file true "SQL文件(.sql)"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/exec-sql-file [post]
// @Security ApiKeyAuth
func ExecSqlFile(c *gin.Context) {
    // 获取参数
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.PostForm("db")
    if dbName == "" {
        dbName = c.Query("db")
    }
    fh, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"code": 400, "message": "缺少文件: " + err.Error()})
        return
    }
    // 读取文件内容（限制大小 ~ 10MB）
    const maxSize = 10 << 20
    f, err := fh.Open()
    if err != nil { c.JSON(400, gin.H{"code":400, "message": err.Error()}); return }
    defer f.Close()
    var r io.Reader = io.LimitReader(f, maxSize)
    bs, err := io.ReadAll(r)
    if err != nil { c.JSON(400, gin.H{"code":400, "message": err.Error()}); return }
    content := string(bs)

    // 获取连接（dbId 优先，失败回退为 instanceId）
    instanceId := uint(id)
    if db, e := dao.GetDbById(uint(id)); e == nil && db.ID > 0 {
        instanceId = db.InstanceId
    }
    conn, err := service.DbInstanceService().GetDbConn(c.Request.Context(), instanceId, dbName)
    if err != nil { c.JSON(500, gin.H{"code":500, "message": err.Error()}); return }

    // 分割SQL语句（简单分号分割，忽略注释与字符串内的分号）
    stmts := splitSQLStatements(content)
    var total, success int
    var errs []string

    for _, s := range stmts {
        sql := strings.TrimSpace(s)
        if sql == "" { continue }
        up := strings.ToUpper(strings.TrimSpace(sql))
        // 仅执行非查询语句（常见导入：DDL/DML）。查询语句可根据需求扩展。
        if strings.HasPrefix(up, "SELECT") || strings.HasPrefix(up, "SHOW") || strings.HasPrefix(up, "DESC") || strings.HasPrefix(up, "EXPLAIN") {
            continue
        }
        total++
        if _, e := conn.ExecContext(c.Request.Context(), sql); e != nil {
            errs = append(errs, e.Error())
        } else {
            success++
        }
    }

    c.JSON(200, gin.H{
        "code": 200,
        "message": "执行完成",
        "data": gin.H{
            "total": total,
            "success": success,
            "failed": total - success,
            "errors": errs,
        },
    })
}

// splitSQLStatements 简易SQL分割：考虑 -- 行注释、/* */块注释、字符串引号，按分号分割
func splitSQLStatements(sql string) []string {
    var res []string
    var sb strings.Builder
    inSingle, inDouble := false, false
    inLineComment, inBlockComment := false, false
    runes := []rune(sql)
    for i := 0; i < len(runes); i++ {
        ch := runes[i]
        next := func() rune { if i+1 < len(runes) { return runes[i+1] }; return 0 }

        if inLineComment {
            if ch == '\n' { inLineComment = false; sb.WriteRune(ch) }
            continue
        }
        if inBlockComment {
            if ch == '*' && next() == '/' { inBlockComment = false; i++ }
            continue
        }
        if !inSingle && !inDouble {
            // 注释开始
            if ch == '-' && next() == '-' { inLineComment = true; i++; continue }
            if ch == '/' && next() == '*' { inBlockComment = true; i++; continue }
            // 语句结束
            if ch == ';' {
                res = append(res, sb.String())
                sb.Reset()
                continue
            }
        }
        // 引号处理
        if ch == '\'' && !inDouble {
            inSingle = !inSingle
            sb.WriteRune(ch)
            continue
        }
        if ch == '"' && !inSingle {
            inDouble = !inDouble
            sb.WriteRune(ch)
            continue
        }
        sb.WriteRune(ch)
    }
    last := strings.TrimSpace(sb.String())
    if last != "" { res = append(res, last) }
    return res
}

// @Tags MySQL数据库管理
// @Summary 获取提示表
// @Produce json
// @Description 获取SQL编辑器自动完成用的表和列信息
// @Param id path int true "数据库ID"
// @Param db query string true "数据库名"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/hint-tables [get]
// @Security ApiKeyAuth
func GetHintTables(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("db")
    service.DbService().GetHintTables(c, uint(id), dbName)
}

// @Tags MySQL数据库管理
// @Summary 快速获取表行数(近似值)
// @Produce json
// @Description 使用information_schema快速返回表的近似行数
// @Param id path int true "数据库ID或实例ID(兼容)"
// @Param db query string true "数据库名"
// @Param table query string true "表名"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/approx-count [get]
// @Security ApiKeyAuth
func GetApproxRowCount(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("db")
    if dbName == "" {
        dbName = c.Query("dbName")
    }
    table := c.Query("table")
    if table == "" {
        table = c.Query("tableName")
    }
    service.DbService().GetApproxRowCount(c, uint(id), dbName, table)
}

// @Tags MySQL数据库管理
// @Summary 获取表状态(快速统计)
// @Produce json
// @Description 使用 SHOW TABLE STATUS 返回近似行数与大小
// @Param id path int true "数据库ID或实例ID(兼容)"
// @Param db query string true "数据库名"
// @Param table query string false "表名(可选，指定则仅返回该表)"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/table-status [get]
// @Security ApiKeyAuth
func GetTableStatus(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("db")
    if dbName == "" {
        dbName = c.Query("dbName")
    }
    table := c.Query("table")
    if table == "" {
        table = c.Query("tableName")
    }
    service.DbService().GetTableStatus(c, uint(id), dbName, table)
}

// @Tags MySQL数据库管理
// @Summary 批量快速行数(近似)
// @Produce json
// @Description 返回多表或整个库的近似行数
// @Param id path int true "数据库ID或实例ID(兼容)"
// @Param db query string true "数据库名"
// @Param tables query string false "表名列表,逗号分隔; 为空则统计库内全部表"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/approx-counts [get]
// @Security ApiKeyAuth
func GetApproxRowCounts(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("db")
    if dbName == "" {
        dbName = c.Query("dbName")
    }
    tables := c.Query("tables")
    service.DbService().GetApproxRowCounts(c, uint(id), dbName, tables)
}

// @Tags MySQL数据库管理
// @Summary 获取单行数据（按主键）
// @Produce json
// @Param id path int true "数据库ID或实例ID(兼容)"
// @Param db query string true "数据库名"
// @Param table query string true "表名"
// @Param pk query string true "主键列名，多个用,分隔"
// @Param pkv query string true "主键值，多个与pk一一对应，用,分隔"
// @Router /api/v1/db/databases/{id}/row [get]
// @Security ApiKeyAuth
func GetRow(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("db")
    table := c.Query("table")
    pk := c.Query("pk")
    pkv := c.Query("pkv")
    service.DbService().GetRow(c, uint(id), dbName, table, pk, pkv)
}

// @Tags MySQL数据库管理
// @Summary 删除单行数据（按主键）
// @Produce json
// @Param id path int true "数据库ID或实例ID(兼容)"
// @Param db query string true "数据库名"
// @Param table query string true "表名"
// @Param pk query string true "主键列名，多个用,分隔"
// @Param pkv query string true "主键值，多个与pk一一对应，用,分隔"
// @Router /api/v1/db/databases/{id}/row [delete]
// @Security ApiKeyAuth
func DeleteRow(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("db")
    table := c.Query("table")
    pk := c.Query("pk")
    pkv := c.Query("pkv")
    service.DbService().DeleteRow(c, uint(id), dbName, table, pk, pkv)
}

// @Tags MySQL数据库管理
// @Summary 导出CSV（基于查询或表）
// @Produce text/csv
// @Param id path int true "数据库ID或实例ID(兼容)"
// @Param db query string true "数据库名"
// @Param sql query string false "SELECT语句（优先）"
// @Param table query string false "表名（当未提供sql时使用）"
// @Param where query string false "可选过滤（仅在table模式下使用）"
// @Router /api/v1/db/databases/{id}/export-csv [get]
// @Security ApiKeyAuth
func ExportCSV(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("db")
    sql := c.Query("sql")
    table := c.Query("table")
    where := c.Query("where")
    service.DbService().ExportCSV(c, uint(id), dbName, sql, table, where)
}

// ExportCSVPost 以POST JSON方式导出CSV，便于前端携带Authorization在当前页下载
type exportReq struct {
    Db    string `json:"db" binding:"required"`
    Sql   string `json:"sql"`
    Table string `json:"table"`
    Where string `json:"where"`
}

func ExportCSVPost(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    var req exportReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"code": 400, "message": err.Error()})
        return
    }
    service.DbService().ExportCSV(c, uint(id), req.Db, req.Sql, req.Table, req.Where)
}

// @Tags MySQL数据库管理
// @Summary 导出SQL脚本（INSERT语句）
// @Produce text/plain
// @Param id path int true "数据库ID或实例ID(兼容)"
// @Param db query string true "数据库名"
// @Param sql query string false "SELECT语句（优先）"
// @Param table query string false "表名（当未提供sql时使用）"
// @Param where query string false "可选过滤（仅在table模式下使用）"
// @Router /api/v1/db/databases/{id}/export-sql [get]
// @Security ApiKeyAuth
func ExportSQL(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("db")
    sql := c.Query("sql")
    table := c.Query("table")
    where := c.Query("where")
    service.DbService().ExportSQL(c, uint(id), dbName, sql, table, where)
}

// ExportSQLPost 以POST JSON方式导出SQL脚本
func ExportSQLPost(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    var req exportReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"code": 400, "message": err.Error()})
        return
    }
    service.DbService().ExportSQL(c, uint(id), req.Db, req.Sql, req.Table, req.Where)
}

// @Tags MySQL数据库管理
// @Summary 导出整个数据库结构
// @Produce text/plain
// @Description 导出整个数据库的所有表结构（DDL）
// @Param id path int true "数据库ID或实例ID"
// @Param db query string true "数据库名"
// @Success 200 {string} string "SQL文件"
// @Router /api/v1/db/databases/{id}/export-database-structure [get]
// @Security ApiKeyAuth
func ExportDatabaseStructure(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("db")
    if dbName == "" {
        c.JSON(400, gin.H{"code": 400, "message": "db 参数必填"})
        return
    }
    service.DbService().ExportDatabaseStructure(c, uint(id), dbName)
}

// @Tags MySQL数据库管理
// @Summary 导出整个数据库结构和数据（同步）
// @Produce text/plain
// @Description 导出整个数据库的所有表结构和数据（小数据库使用，大数据库建议使用异步接口）
// @Param id path int true "数据库ID或实例ID"
// @Param db query string true "数据库名"
// @Success 200 {string} string "SQL文件"
// @Router /api/v1/db/databases/{id}/export-database-full [get]
// @Security ApiKeyAuth
func ExportDatabaseFull(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("db")
    if dbName == "" {
        c.JSON(400, gin.H{"code": 400, "message": "db 参数必填"})
        return
    }
    service.DbService().ExportDatabaseStructureAndData(c, uint(id), dbName)
}

// @Tags MySQL数据库管理
// @Summary 异步导出整个数据库结构和数据
// @Produce json
// @Description 创建异步导出任务，适用于大数据库导出
// @Param id path int true "数据库ID或实例ID"
// @Param db query string true "数据库名"
// @Param exportType query string false "导出类型: structure(仅结构)/full(结构+数据)，默认full"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/export-database-full-async [post]
// @Security ApiKeyAuth
func ExportDatabaseFullAsync(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    dbName := c.Query("db")
    if dbName == "" {
        c.JSON(400, gin.H{"code": 400, "message": "db 参数必填"})
        return
    }
    exportType := c.Query("exportType")
    if exportType == "" {
        exportType = "full"
    }
    if exportType != "structure" && exportType != "full" {
        c.JSON(400, gin.H{"code": 400, "message": "exportType 必须是 structure 或 full"})
        return
    }
    service.DbService().ExportDatabaseAsync(c, uint(id), dbName, exportType)
}

// @Tags MySQL数据库管理
// @Summary 查询导出任务状态
// @Produce json
// @Description 查询导出任务的状态和进度
// @Param taskId path string true "任务ID"
// @Success 200 {object} result.Result
// @Router /api/v1/db/export-tasks/{taskId} [get]
// @Security ApiKeyAuth
func GetExportTaskStatus(c *gin.Context) {
    taskId := c.Param("taskId")
    if taskId == "" {
        c.JSON(400, gin.H{"code": 400, "message": "taskId 参数必填"})
        return
    }
    service.DbService().GetExportTaskStatus(c, taskId)
}

// @Tags MySQL数据库管理
// @Summary 下载导出文件
// @Produce application/octet-stream
// @Description 下载已完成的导出文件
// @Param taskId path string true "任务ID"
// @Success 200 {file} file "SQL文件"
// @Router /api/v1/db/export-tasks/{taskId}/download [get]
// @Security ApiKeyAuth
func DownloadExportFile(c *gin.Context) {
    taskId := c.Param("taskId")
    if taskId == "" {
        c.JSON(400, gin.H{"code": 400, "message": "taskId 参数必填"})
        return
    }
    service.DbService().DownloadExportFile(c, taskId)
}

// ---------- Create Database & Create Table ----------

// 使用 model 中的请求结构体

// @Tags MySQL数据库管理
// @Summary 创建数据库
// @Produce json
// @Description 在实例上创建数据库
// @Param id path int true "数据库ID或实例ID(兼容)"
// @Param data body CreateDatabaseReq true "创建库参数"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/create-database [post]
// @Security ApiKeyAuth
func CreateDatabase(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    var req model.CreateDatabaseReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"code": 400, "message": err.Error()})
        return
    }
    service.DbService().CreateDatabase(c, uint(id), &req)
}


// @Tags MySQL数据库管理
// @Summary 创建表
// @Produce json
// @Description 在指定数据库下创建表
// @Param id path int true "数据库ID或实例ID(兼容)"
// @Param data body CreateTableReq true "创建表参数"
// @Success 200 {object} result.Result
// @Router /api/v1/db/databases/{id}/create-table [post]
// @Security ApiKeyAuth
func CreateTable(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    var req model.CreateTableReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"code": 400, "message": err.Error()})
        return
    }
    service.DbService().CreateTable(c, uint(id), &req)
}

// @Tags MySQL数据库管理
// @Summary 格式化SQL语句
// @Accept json
// @Produce json
// @Description 将一长串SQL语句格式化，添加换行和缩进
// @Param body body object true "SQL内容"
// @Success 200 {object} result.Result
// @Router /api/v1/db/format-sql [post]
// @Security ApiKeyAuth
func FormatSQL(c *gin.Context) {
	var req struct {
		SQL string `json:"sql" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "请求参数错误: " + err.Error()})
		return
	}

	formattedSQL := utils.FormatSQL(req.SQL)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "格式化成功",
		"data": gin.H{
			"sql":          req.SQL,
			"formattedSql": formattedSQL,
		},
	})
}

// @Tags MySQL数据库管理
// @Summary SQL预检查
// @Accept json
// @Produce json
// @Description 预检查SQL，返回预估影响的行数（不实际执行）
// @Param body body object true "SQL信息"
// @Success 200 {object} result.Result
// @Router /api/v1/db/pre-check-sql [post]
// @Security ApiKeyAuth
func PreCheckSQL(c *gin.Context) {
	var dto model.DbSqlExecDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "请求参数错误: " + err.Error()})
		return
	}

	service.DbService().PreCheckSQL(c, &dto)
}
