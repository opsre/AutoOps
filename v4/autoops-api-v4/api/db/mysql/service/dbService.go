package service

import (
    "context"
    "dodevops-api/api/db/mysql/dao"
    "dodevops-api/api/db/mysql/model"
    "dodevops-api/common/result"
    "dodevops-api/common/util"
    "dodevops-api/pkg/redis"
    "dodevops-api/pkg/dbm/dbi"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "crypto/sha1"
    "encoding/hex"
    "encoding/json"
    "strconv"
    "strings"
    "time"
    "regexp"
    "os"
    "path/filepath"
)

// IDbService 数据库服务接口
type IDbService interface {
	GetDbList(c *gin.Context, query *model.DbQuery)
	GetDbById(c *gin.Context, id uint)
	CreateDb(c *gin.Context, dto *model.DbDto)
	UpdateDb(c *gin.Context, dto *model.DbDto)
	DeleteDb(c *gin.Context, id uint)
	GetTables(c *gin.Context, dbId uint, dbName string)
	GetTablesByInstanceId(c *gin.Context, instanceId uint, dbName string)
	GetTableColumns(c *gin.Context, dbId uint, dbName, tableName string)
	GetTableIndexes(c *gin.Context, dbId uint, dbName, tableName string)
	GetTableDDL(c *gin.Context, dbId uint, dbName, tableName string)
	CopyTable(c *gin.Context, dto *model.CopyTableDto)
    ExecSql(c *gin.Context, dto *model.DbSqlExecDto) *model.DbSqlExecResult
    GetHintTables(c *gin.Context, dbId uint, dbName string)
    GetApproxRowCount(c *gin.Context, dbId uint, dbName, table string)
    GetTableStatus(c *gin.Context, dbId uint, dbName, table string)
    GetApproxRowCounts(c *gin.Context, dbId uint, dbName, tables string)
    GetRow(c *gin.Context, dbId uint, dbName, table, pk, pkv string)
    DeleteRow(c *gin.Context, dbId uint, dbName, table, pk, pkv string)
    ExportCSV(c *gin.Context, dbId uint, dbName, sql, table, where string)
    ExportSQL(c *gin.Context, dbId uint, dbName, sql, table, where string)
    ExportDatabaseStructure(c *gin.Context, dbId uint, dbName string)
    ExportDatabaseStructureAndData(c *gin.Context, dbId uint, dbName string)
    ExportDatabaseAsync(c *gin.Context, dbId uint, dbName, exportType string)
    GetExportTaskStatus(c *gin.Context, taskId string)
    DownloadExportFile(c *gin.Context, taskId string)
    CreateDatabase(c *gin.Context, id uint, req *model.CreateDatabaseReq)
    PreCheckSQL(c *gin.Context, dto *model.DbSqlExecDto)
    CreateTable(c *gin.Context, id uint, req *model.CreateTableReq)
}

type DbServiceImpl struct{}

// 默认限制（可按需改为从配置加载）
const (
    defaultSqlTimeoutSeconds = 60
    defaultMaxQueryRows      = 1000
)

// paramFor returns a dialect-specific placeholder token for positional arg index i.
// For Postgres family, it returns $1, $2...; otherwise '?'.
func paramFor(conn *dbi.DbConn, i int) string {
    switch strings.ToLower(string(conn.Info.Type)) {
    case "postgres", "gauss", "kingbasees", "vastbase":
        return fmt.Sprintf("$%d", i+1)
    default:
        return "?"
    }
}

// GetDbList 获取数据库列表
func (s DbServiceImpl) GetDbList(c *gin.Context, query *model.DbQuery) {
	dbs, total, err := dao.GetDbListWithInstance(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	data := map[string]interface{}{
		"list":     dbs,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	}
	result.Success(c, data)
}

// GetDbById 根据ID获取数据库
func (s DbServiceImpl) GetDbById(c *gin.Context, id uint) {
	db, err := dao.GetDbById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, db)
}

// CreateDb 创建数据库
func (s DbServiceImpl) CreateDb(c *gin.Context, dto *model.DbDto) {
	// 检查实例是否存在
	instance, err := dao.GetDbInstanceById(dto.InstanceId)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}

	// 检查编码是否已存在
	if dto.Code != "" && dao.CheckDbExists(dto.Code, 0) {
		result.Failed(c, int(result.ApiCode.FAILED), "数据库编码已存在")
		return
	}

	// 生成编码
	if dto.Code == "" {
		dto.Code = uuid.New().String()
	}

	now := time.Now()
	db := &model.Db{
		Code:         dto.Code,
		Name:         dto.Name,
		Database:     dto.Database,
		Remark:       dto.Remark,
		InstanceId:   dto.InstanceId,
		InstanceCode: instance.Code,
		Status:       1,
		CreateTime:   util.HTime{Time: now},
		UpdateTime:   util.HTime{Time: now},
		Creator:      c.GetString("username"),
		CreatorId:    uint(c.GetInt("userId")),
	}

	if err := dao.CreateDb(db); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, db.ID)
}

// UpdateDb 更新数据库
func (s DbServiceImpl) UpdateDb(c *gin.Context, dto *model.DbDto) {
	// 检查数据库是否存在
	db, err := dao.GetDbById(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "数据库不存在")
		return
	}

	// 检查编码是否重复
	if dto.Code != "" && dto.Code != db.Code && dao.CheckDbExists(dto.Code, dto.ID) {
		result.Failed(c, int(result.ApiCode.FAILED), "数据库编码已存在")
		return
	}

	// 更新字段
	db.Name = dto.Name
	db.Database = dto.Database
	db.Remark = dto.Remark
	if dto.InstanceId > 0 && dto.InstanceId != db.InstanceId {
		instance, err := dao.GetDbInstanceById(dto.InstanceId)
		if err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
			return
		}
		db.InstanceId = dto.InstanceId
		db.InstanceCode = instance.Code
	}
	db.UpdateTime = util.HTime{Time: time.Now()}
	db.Modifier = c.GetString("username")
	db.ModifierId = uint(c.GetInt("userId"))

	if err := dao.UpdateDb(&db); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, nil)
}

// DeleteDb 删除数据库
func (s DbServiceImpl) DeleteDb(c *gin.Context, id uint) {
	// 删除相关的SQL模板
	dao.DeleteDbSqlsByDbId(id)
	// 删除相关的执行记录
	dao.DeleteDbSqlExecsByDbId(id)

	if err := dao.DeleteDb(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, nil)
}


// GetTablesByInstanceId 根据实例ID和数据库名获取表列表
func (s DbServiceImpl) GetTablesByInstanceId(c *gin.Context, instanceId uint, dbName string) {
	conn, err := DbInstanceService().GetDbConn(c.Request.Context(), instanceId, dbName)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	metadata := conn.GetMetadata()
	tables, err := metadata.GetTables()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取表列表失败: %v", err))
		return
	}

	// 转换为前端需要的格式
	var tableList []model.TableInfo
    for _, table := range tables {
        total := table.DataLength + table.IndexLength
        tableList = append(tableList, model.TableInfo{
            TableName:    table.TableName,
            TableComment: table.TableComment,
            TableType:    "",              // dbi.Table没有TableType字段，设为空
            CreateTime:   table.CreateTime,
            TableRows:    int64(table.TableRows),
            DataLength:   table.DataLength,
            IndexLength:  table.IndexLength,
            TotalLength:  total,
            Size:         humanSize(total),
        })
    }

	result.Success(c, tableList)
}
func (s DbServiceImpl) GetTables(c *gin.Context, dbId uint, dbName string) {
	conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	metadata := conn.GetMetadata()
	tables, err := metadata.GetTables()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取表列表失败: %v", err))
		return
	}

	// 转换为前端需要的格式
	var tableList []model.TableInfo
    for _, table := range tables {
        total := table.DataLength + table.IndexLength
        tableList = append(tableList, model.TableInfo{
            TableName:    table.TableName,
            TableComment: table.TableComment,
            TableType:    "",              // dbi.Table没有TableType字段，设为空
            CreateTime:   table.CreateTime,
            TableRows:    int64(table.TableRows),
            DataLength:   table.DataLength,
            IndexLength:  table.IndexLength,
            TotalLength:  total,
            Size:         humanSize(total),
        })
    }

	result.Success(c, tableList)
}

// GetTableColumns 获取表列信息
func (s DbServiceImpl) GetTableColumns(c *gin.Context, dbId uint, dbName, tableName string) {
	conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	metadata := conn.GetMetadata()
	columns, err := metadata.GetColumns(tableName)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取列信息失败: %v", err))
		return
	}

	// 转换为前端需要的格式
	var columnList []model.ColumnInfo
	for _, col := range columns {
		nullable := "NO"
		if col.Nullable {
			nullable = "YES"
		}
		columnList = append(columnList, model.ColumnInfo{
			ColumnName:    col.ColumnName,
			ColumnType:    col.GetColumnType(), // 调用方法获取完整类型
			ColumnComment: col.ColumnComment,
			DataType:      col.DataType,
			IsNullable:    nullable,
			ColumnKey:     "",        // dbi.Column没有此字段
			Extra:         "",        // 需要从Extra map转换
		})
	}

	result.Success(c, columnList)
}

// GetTableIndexes 获取表索引信息
func (s DbServiceImpl) GetTableIndexes(c *gin.Context, dbId uint, dbName, tableName string) {
	conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	metadata := conn.GetMetadata()
	indexes, err := metadata.GetTableIndex(tableName)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取索引信息失败: %v", err))
		return
	}

	// 转换为前端需要的格式
	var indexList []model.IndexInfo
	for _, idx := range indexes {
		nonUnique := 0
		if !idx.IsUnique {
			nonUnique = 1
		}
		indexList = append(indexList, model.IndexInfo{
			IndexName:  idx.IndexName,
			ColumnName: idx.ColumnName,
			IndexType:  idx.IndexType,
			SeqInIndex: idx.SeqInIndex,
			NonUnique:  nonUnique,
		})
	}

	result.Success(c, indexList)
}

// GetTableDDL 获取表DDL
func (s DbServiceImpl) GetTableDDL(c *gin.Context, dbId uint, dbName, tableName string) {
	conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	metadata := conn.GetMetadata()
	ddl, err := metadata.GetTableDDL(tableName, false)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取表DDL失败: %v", err))
		return
	}

	result.Success(c, map[string]interface{}{
		"tableName": tableName,
		"ddl":       ddl,
	})
}

// CopyTable 复制表
func (s DbServiceImpl) CopyTable(c *gin.Context, dto *model.CopyTableDto) {
	conn, err := s.getDbConn(c.Request.Context(), dto.DbId, dto.Db)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	dialect := conn.GetDialect()
	copyParam := &dbi.DbCopyTable{
		TableName: dto.TableName,
		// NewTable字段不存在，CopyTable方法会自动生成新表名
		CopyData:  dto.CopyData,
	}

	if err := dialect.CopyTable(copyParam); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("复制表失败: %v", err))
		return
	}

	result.Success(c, nil)
}

// ExecSql 执行SQL
func (s DbServiceImpl) ExecSql(c *gin.Context, dto *model.DbSqlExecDto) *model.DbSqlExecResult {
    startTime := time.Now()

    // 上下文超时控制（默认60秒，可被 dto.TimeoutSeconds 覆盖）
    ttl := time.Duration(defaultSqlTimeoutSeconds) * time.Second
    if dto.TimeoutSeconds > 0 {
        ttl = time.Duration(dto.TimeoutSeconds) * time.Second
    }
    ctx, cancel := context.WithTimeout(c.Request.Context(), ttl)
    defer cancel()

	var conn *dbi.DbConn
	var err error

	// 优先使用 InstanceId，兼容 DbId
    if dto.InstanceId > 0 {
        conn, err = DbInstanceService().GetDbConn(ctx, dto.InstanceId, dto.Db)
    } else {
        conn, err = s.getDbConn(ctx, dto.DbId, dto.Db)
    }

	if err != nil {
		s.saveExecRecord(c, dto, model.ExecStatusFail, 0, "", err.Error(), time.Since(startTime).Milliseconds())
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return nil
	}

	// 去除SQL首尾空格
	sql := strings.TrimSpace(dto.Sql)
	if sql == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "SQL不能为空")
		return nil
	}

    // 检查是否包含多条SQL语句（通过分号判断）
    statements := s.splitSQLStatements(sql)

    // 如果包含多条查询语句，执行多语句查询
    if len(statements) > 1 {
        execResult, execErr := s.execMultipleQueries(ctx, conn, statements, dto.MaxRows)
        execTime := time.Since(startTime).Milliseconds()
        if execErr != nil {
            s.saveExecRecord(c, dto, model.ExecStatusFail, model.SqlTypeQuery, "", execErr.Error(), execTime)
            result.Failed(c, int(result.ApiCode.FAILED), execErr.Error())
            return nil
        }
        if execResult != nil {
            execResult.ExecTime = execTime
        }
        s.saveExecRecord(c, dto, model.ExecStatusSuccess, model.SqlTypeQuery, "", "执行成功", execTime)
        result.Success(c, execResult)
        return execResult
    }

    // 单条SQL语句的原有逻辑
    // 判断SQL类型
    sqlType := s.getSqlType(sql)

    // 如果是查询语句，且开启缓存，尝试读取缓存
    if sqlType == model.SqlTypeQuery && dto.UseCache && dto.CacheSeconds > 0 && redis.RedisDb != nil {
        if cached := s.tryGetSqlCache(dto, sql); cached != nil {
            execTime := time.Since(startTime).Milliseconds()
            cached.ExecTime = execTime
            // 记录一次执行日志（标注缓存命中）
            s.saveExecRecord(c, dto, model.ExecStatusSuccess, sqlType, "", "缓存命中", execTime)
            result.Success(c, cached)
            return cached
        }
    }

	var execResult *model.DbSqlExecResult
	var execErr error

    if sqlType == model.SqlTypeQuery {
        // 默认最大返回行数
        if dto.MaxRows == 0 {
            dto.MaxRows = defaultMaxQueryRows
        }
        // 查询语句（支持最大行限制与超时）
        execResult, execErr = s.execQueryWithLimit(ctx, conn, sql, dto.MaxRows)
    } else {
        // 更新语句
        execResult, execErr = s.execUpdateCtx(ctx, conn, sql)
    }

    execTime := time.Since(startTime).Milliseconds()
    if execErr != nil {
        s.saveExecRecord(c, dto, model.ExecStatusFail, sqlType, "", execErr.Error(), execTime)
        result.Failed(c, int(result.ApiCode.FAILED), execErr.Error())
        return nil
    }
    if execResult != nil {
        execResult.ExecTime = execTime
    }

    s.saveExecRecord(c, dto, model.ExecStatusSuccess, sqlType, "", "执行成功", execTime)

    // 成功后写入缓存（仅查询类SQL）
    if sqlType == model.SqlTypeQuery && dto.UseCache && dto.CacheSeconds > 0 && redis.RedisDb != nil {
        _ = s.trySetSqlCache(dto, sql, execResult, time.Duration(dto.CacheSeconds)*time.Second)
    }
    result.Success(c, execResult)
    return execResult
}

// GetHintTables 获取提示表(用于SQL编辑器自动完成)
func (s DbServiceImpl) GetHintTables(c *gin.Context, dbId uint, dbName string) {
	conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	metadata := conn.GetMetadata()

	// 获取所有表
	tables, err := metadata.GetTables()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取表列表失败: %v", err))
		return
	}

	var hintTables []model.HintTable
	for _, table := range tables {
		// 获取表的列信息
		columns, err := metadata.GetColumns(table.TableName)
		if err != nil {
			continue
		}

		var hintColumns []model.HintColumn
		for _, col := range columns {
			hintColumns = append(hintColumns, model.HintColumn{
				ColumnName: col.ColumnName,
				ColumnType: col.GetColumnType(), // 调用方法
				ColumnComment: col.ColumnComment,
			})
		}

		hintTables = append(hintTables, model.HintTable{
			TableName: table.TableName,
			Columns:   hintColumns,
		})
	}

	result.Success(c, hintTables)
}

// 辅助方法

func (s DbServiceImpl) getDbConn(ctx context.Context, dbId uint, dbName string) (*dbi.DbConn, error) {
    // 先按数据库ID查找（标准行为）
    if db, err := dao.GetDbById(dbId); err == nil && db.ID > 0 {
        return DbInstanceService().GetDbConn(ctx, db.InstanceId, dbName)
    }

    // 兼容旧用法：如果按数据库ID未找到，尝试把 id 当作实例ID直接连接
    if conn, err := DbInstanceService().GetDbConn(ctx, dbId, dbName); err == nil {
        return conn, nil
    }

    return nil, fmt.Errorf("数据库不存在或实例ID无效")
}

func (s DbServiceImpl) getSqlType(sql string) int8 {
	upperSql := strings.ToUpper(strings.TrimSpace(sql))

	if strings.HasPrefix(upperSql, "SELECT") || strings.HasPrefix(upperSql, "SHOW") ||
	   strings.HasPrefix(upperSql, "DESC") || strings.HasPrefix(upperSql, "EXPLAIN") {
		return model.SqlTypeQuery
	} else if strings.HasPrefix(upperSql, "INSERT") {
		return model.SqlTypeInsert
	} else if strings.HasPrefix(upperSql, "UPDATE") {
		return model.SqlTypeUpdate
	} else if strings.HasPrefix(upperSql, "DELETE") {
		return model.SqlTypeDelete
	} else {
		return model.SqlTypeDDL
	}
}

func (s DbServiceImpl) execQueryWithLimit(ctx context.Context, conn *dbi.DbConn, sql string, maxRows int) (*model.DbSqlExecResult, error) {
    rows := make([]map[string]interface{}, 0, 16)
    current := 0
    cols, err := conn.WalkQueryRows(ctx, sql, func(row map[string]any, ccols []*dbi.QueryColumn) error {
        current++
        if maxRows != 0 && current > maxRows {
            return dbi.NewStopWalkQueryError(fmt.Sprintf("exceed the maximum number of query records %d", maxRows))
        }
        // 转换 map[string]any 到 map[string]interface{}
        m := make(map[string]interface{}, len(row))
        for k, v := range row {
            m[k] = v
        }
        rows = append(rows, m)
        return nil
    })
    if err != nil {
        return nil, err
    }

    // 转换列名
    var columnNames []string
    for _, col := range cols {
        columnNames = append(columnNames, col.Name)
    }

    // 构建列元信息：名称/类型/备注（备注需通过表元数据查询）
    metas := make([]model.ColumnMeta, 0, len(cols))
    comments := map[string]string{}
    // 简单提取表名，仅在简单查询时有效
    if tn := s.extractTableName(sql); tn != "" {
        // 用原始表名从元数据获取列备注
        md := conn.GetMetadata()
        if cs, e := md.GetColumns(tn); e == nil {
            for _, c := range cs {
                comments[c.ColumnName] = c.ColumnComment
            }
        }
    }
    for _, c := range cols {
        metas = append(metas, model.ColumnMeta{
            Name:    c.Name,
            Type:    c.Type,
            Comment: comments[c.Name],
        })
    }

    return &model.DbSqlExecResult{Columns: columnNames, Rows: rows, ColumnMetas: metas}, nil
}

func (s DbServiceImpl) execUpdateCtx(ctx context.Context, conn *dbi.DbConn, sql string) (*model.DbSqlExecResult, error) {
    affectedRows, err := conn.ExecContext(ctx, sql)
    if err != nil {
        return nil, err
    }

    result := &model.DbSqlExecResult{
        AffectedRows: affectedRows,
    }

    return result, nil
}

// execMultipleQueries 执行多条查询语句并合并结果
func (s DbServiceImpl) execMultipleQueries(ctx context.Context, conn *dbi.DbConn, statements []string, maxRows int) (*model.DbSqlExecResult, error) {
    if maxRows == 0 {
        maxRows = defaultMaxQueryRows
    }

    var allRows []map[string]interface{}
    var allColumns []string
    var allColumnMetas []model.ColumnMeta
    columnSet := make(map[string]bool)

    for _, stmt := range statements {
        stmt = strings.TrimSpace(stmt)
        if stmt == "" {
            continue
        }

        // 跳过非查询语句
        upperStmt := strings.ToUpper(stmt)
        if !strings.HasPrefix(upperStmt, "SELECT") &&
           !strings.HasPrefix(upperStmt, "SHOW") &&
           !strings.HasPrefix(upperStmt, "DESC") &&
           !strings.HasPrefix(upperStmt, "EXPLAIN") {
            continue
        }

        // 执行单条查询
        queryResult, err := s.execQueryWithLimit(ctx, conn, stmt, maxRows)
        if err != nil {
            return nil, fmt.Errorf("执行SQL失败 [%s]: %v", stmt, err)
        }

        // 合并列信息（使用set去重）
        for _, col := range queryResult.Columns {
            if !columnSet[col] {
                allColumns = append(allColumns, col)
                columnSet[col] = true
            }
        }

        // 合并列元信息
        for _, meta := range queryResult.ColumnMetas {
            if !columnSet[meta.Name] {
                allColumnMetas = append(allColumnMetas, meta)
            }
        }

        // 合并行数据
        allRows = append(allRows, queryResult.Rows...)

        // 检查是否超过最大行数限制
        if len(allRows) >= maxRows {
            allRows = allRows[:maxRows]
            break
        }
    }

    return &model.DbSqlExecResult{
        Columns:      allColumns,
        Rows:         allRows,
        ColumnMetas:  allColumnMetas,
    }, nil
}

// splitSQLStatements 分割SQL语句，支持注释和字符串
func (s DbServiceImpl) splitSQLStatements(sql string) []string {
    var res []string
    var sb strings.Builder
    inSingle, inDouble := false, false
    inLineComment, inBlockComment := false, false
    runes := []rune(sql)

    for i := 0; i < len(runes); i++ {
        ch := runes[i]
        next := func() rune {
            if i+1 < len(runes) {
                return runes[i+1]
            }
            return 0
        }

        // 处理行注释
        if inLineComment {
            if ch == '\n' {
                inLineComment = false
                sb.WriteRune(ch)
            }
            continue
        }

        // 处理块注释
        if inBlockComment {
            if ch == '*' && next() == '/' {
                inBlockComment = false
                i++
            }
            continue
        }

        // 不在字符串内时处理注释和语句分隔符
        if !inSingle && !inDouble {
            // 注释开始
            if ch == '-' && next() == '-' {
                inLineComment = true
                i++
                continue
            }
            if ch == '/' && next() == '*' {
                inBlockComment = true
                i++
                continue
            }
            // 语句结束
            if ch == ';' {
                stmt := strings.TrimSpace(sb.String())
                if stmt != "" {
                    res = append(res, stmt)
                }
                sb.Reset()
                continue
            }
        }

        // 处理单引号
        if ch == '\'' && !inDouble {
            inSingle = !inSingle
            sb.WriteRune(ch)
            continue
        }

        // 处理双引号
        if ch == '"' && !inSingle {
            inDouble = !inDouble
            sb.WriteRune(ch)
            continue
        }

        sb.WriteRune(ch)
    }

    // 处理最后一条语句
    last := strings.TrimSpace(sb.String())
    if last != "" {
        res = append(res, last)
    }

    return res
}

// humanSize 将字节数转为可读字符串
func humanSize(bytes int64) string {
    const (
        KB = 1024
        MB = 1024 * KB
        GB = 1024 * MB
        TB = 1024 * GB
    )
    b := float64(bytes)
    switch {
    case bytes >= TB:
        return fmt.Sprintf("%.0f TB", b/float64(TB))
    case bytes >= GB:
        return fmt.Sprintf("%.0f GB", b/float64(GB))
    case bytes >= MB:
        return fmt.Sprintf("%.0f MB", b/float64(MB))
    case bytes >= KB:
        return fmt.Sprintf("%.0f KB", b/float64(KB))
    default:
        return fmt.Sprintf("%d B", bytes)
    }
}

func (s DbServiceImpl) saveExecRecord(c *gin.Context, dto *model.DbSqlExecDto, status int8, sqlType int8, oldValue, res string, execTime int64) {
	// 提取表名(简单实现)
	tableName := s.extractTableName(dto.Sql)

	// 优先使用 InstanceId
	dbId := dto.DbId
	if dto.InstanceId > 0 {
		dbId = dto.InstanceId
	}

	record := &model.DbSqlExec{
		DbId:      dbId,
		DbName:    dto.Db,
		TblName:   tableName,
		Type:      sqlType,
		Sql:       dto.Sql,
		OldValue:  oldValue,
		Remark:    dto.Remark,
		Status:    status,
		Res:       res,
		ExecTime:  execTime,
		CreateTime: util.HTime{Time: time.Now()},
		Creator:   c.GetString("username"),
		CreatorId: uint(c.GetInt("userId")),
	}

	dao.CreateDbSqlExec(record)
}

func (s DbServiceImpl) extractTableName(sql string) string {
    // 简单的表名提取逻辑：用大写词匹配关键字，但返回原始SQL对应位置的词，保留表名大小写
    trimmed := strings.TrimSpace(sql)
    upperSql := strings.ToUpper(trimmed)
    upperWords := strings.Fields(upperSql)
    origWords := strings.Fields(trimmed)
    for i, word := range upperWords {
        if (word == "FROM" || word == "INTO" || word == "UPDATE") && i+1 < len(origWords) {
            tableName := origWords[i+1]
            // 移除可能的反引号与逗号
            tableName = strings.Trim(tableName, "` ,")
            return tableName
        }
    }
    return ""
}

// GetApproxRowCount 使用information_schema返回近似行数
func (s DbServiceImpl) GetApproxRowCount(c *gin.Context, dbId uint, dbName, table string) {
    if dbName == "" || table == "" {
        result.Failed(c, int(result.ApiCode.FAILED), "db和table必填")
        return
    }
    conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
    if err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), err.Error())
        return
    }
    sql := "SELECT TABLE_ROWS AS cnt FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? LIMIT 1"
    _, rows, err := conn.Query(sql, dbName, table)
    if err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), err.Error())
        return
    }
    var cnt int64
    if len(rows) > 0 {
        switch v := rows[0]["cnt"].(type) {
        case int64:
            cnt = v
        case int:
            cnt = int64(v)
        case uint64:
            cnt = int64(v)
        case []uint8:
            // 某些版本可能返回字节，需要转成字符串再解析
            s := string(v)
            if parsed, err := strconv.ParseInt(s, 10, 64); err == nil { cnt = parsed }
        }
    }
    result.Success(c, gin.H{"table": table, "rows": cnt})
}

// GetTableStatus 使用 SHOW TABLE STATUS 返回近似行数与大小
func (s DbServiceImpl) GetTableStatus(c *gin.Context, dbId uint, dbName, table string) {
    if dbName == "" {
        result.Failed(c, int(result.ApiCode.FAILED), "db必填")
        return
    }
    conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
    if err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), err.Error())
        return
    }
    var sql string
    var args []any
    if table != "" {
        sql = fmt.Sprintf("SHOW TABLE STATUS FROM `%s` LIKE ?", dbName)
        args = []any{table}
    } else {
        sql = fmt.Sprintf("SHOW TABLE STATUS FROM `%s`", dbName)
    }
    _, rows, err := conn.Query(sql, args...)
    if err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), err.Error())
        return
    }
    list := make([]model.TableStatus, 0, len(rows))
    for _, r := range rows {
        ts := model.TableStatus{}
        if v, ok := r["Name"].(string); ok { ts.TableName = v }
        switch rv := r["Rows"].(type) {
        case int64:
            ts.Rows = rv
        case int:
            ts.Rows = int64(rv)
        }
        switch dv := r["Data_length"].(type) {
        case int64:
            ts.DataLength = dv
        case int:
            ts.DataLength = int64(dv)
        }
        switch iv := r["Index_length"].(type) {
        case int64:
            ts.IndexLength = iv
        case int:
            ts.IndexLength = int64(iv)
        }
        if v, ok := r["Comment"].(string); ok { ts.Comment = v }
        if v, ok := r["Create_time"].(string); ok { ts.CreateTime = v }
        list = append(list, ts)
    }
    result.Success(c, list)
}

// GetApproxRowCounts 批量获取近似行数
func (s DbServiceImpl) GetApproxRowCounts(c *gin.Context, dbId uint, dbName, tables string) {
    if dbName == "" {
        result.Failed(c, int(result.ApiCode.FAILED), "db必填")
        return
    }
    conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
    if err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), err.Error())
        return
    }

    baseSql := "SELECT TABLE_NAME as name, TABLE_ROWS as cnt FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?"
    var args []any
    args = append(args, dbName)
    if strings.TrimSpace(tables) != "" {
        // 拆分并构造IN列表
        parts := strings.Split(tables, ",")
        // 去空白
        cleaned := make([]string, 0, len(parts))
        for _, p := range parts {
            p = strings.TrimSpace(p)
            if p != "" {
                cleaned = append(cleaned, p)
            }
        }
        if len(cleaned) > 0 {
            // 构造占位符
            placeholders := make([]string, len(cleaned))
            for i := range cleaned {
                placeholders[i] = "?"
                args = append(args, cleaned[i])
            }
            baseSql += " AND TABLE_NAME IN (" + strings.Join(placeholders, ",") + ")"
        }
    }

    _, rows, err := conn.Query(baseSql, args...)
    if err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), err.Error())
        return
    }
    // 返回 map[string]int64 结构
    res := map[string]int64{}
    for _, r := range rows {
        name, _ := r["name"].(string)
        var cnt int64
        switch v := r["cnt"].(type) {
        case int64:
            cnt = v
        case int:
            cnt = int64(v)
        case uint64:
            cnt = int64(v)
        case []uint8:
            if parsed, err := strconv.ParseInt(string(v), 10, 64); err == nil { cnt = parsed }
        }
        if name != "" {
            res[name] = cnt
        }
    }
    result.Success(c, res)
}

// GetRow 按主键获取单行
func (s DbServiceImpl) GetRow(c *gin.Context, dbId uint, dbName, table, pk, pkv string) {
    if dbName == "" || table == "" || pk == "" || pkv == "" {
        result.Failed(c, int(result.ApiCode.FAILED), "db/table/pk/pkv 不能为空")
        return
    }
    conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }

    pkCols := strings.Split(pk, ",")
    pkVals := strings.Split(pkv, ",")
    if len(pkCols) != len(pkVals) { result.Failed(c, int(result.ApiCode.FAILED), "pk 与 pkv 数量不一致"); return }

    conds := make([]string, 0, len(pkCols))
    args := make([]any, 0, len(pkCols))
    q := conn.GetDialect().Quoter()
    for i, col := range pkCols {
        conds = append(conds, fmt.Sprintf("%s = %s", q.Quote(strings.TrimSpace(col)), paramFor(conn, i)))
        args = append(args, strings.TrimSpace(pkVals[i]))
    }
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", q.Quote(table), strings.Join(conds, " AND "))
    _, rows, err := conn.Query(query, args...)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    if len(rows) == 0 { result.Success(c, nil); return }
    result.Success(c, rows[0])
}

// DeleteRow 按主键删除单行
func (s DbServiceImpl) DeleteRow(c *gin.Context, dbId uint, dbName, table, pk, pkv string) {
    if dbName == "" || table == "" || pk == "" || pkv == "" {
        result.Failed(c, int(result.ApiCode.FAILED), "db/table/pk/pkv 不能为空")
        return
    }
    conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }

    pkCols := strings.Split(pk, ",")
    pkVals := strings.Split(pkv, ",")
    if len(pkCols) != len(pkVals) { result.Failed(c, int(result.ApiCode.FAILED), "pk 与 pkv 数量不一致"); return }

    conds := make([]string, 0, len(pkCols))
    args := make([]any, 0, len(pkCols))
    q := conn.GetDialect().Quoter()
    for i, col := range pkCols {
        conds = append(conds, fmt.Sprintf("%s = %s", q.Quote(strings.TrimSpace(col)), paramFor(conn, i)))
        args = append(args, strings.TrimSpace(pkVals[i]))
    }
    del := fmt.Sprintf("DELETE FROM %s WHERE %s", q.Quote(table), strings.Join(conds, " AND "))
    
    _, err = conn.ExecContext(c.Request.Context(), del, args...)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    result.Success(c, nil)
}

// ExportCSV 导出CSV
func (s DbServiceImpl) ExportCSV(c *gin.Context, dbId uint, dbName, sql, table, where string) {
    conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    query := strings.TrimSpace(sql)
    if query == "" {
        if table == "" { result.Failed(c, int(result.ApiCode.FAILED), "sql 或 table 必填"); return }
        q := conn.GetDialect().Quoter()
        query = fmt.Sprintf("SELECT * FROM %s", q.Quote(table))
        if strings.TrimSpace(where) != "" { query += " WHERE " + where }
    }

    c.Header("Content-Type", "text/csv; charset=utf-8")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=export_%d.csv", time.Now().Unix()))

    w := c.Writer
    wroteHeader := false
    // 写入BOM以便Excel识别UTF-8
    _, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})
    // 行写入函数
    writeRow := func(values []string) {
        line := strings.Join(values, ",") + "\n"
        _, _ = w.Write([]byte(line))
    }

    _, err = conn.WalkQueryRows(c.Request.Context(), query, func(row map[string]any, cols []*dbi.QueryColumn) error {
        if !wroteHeader {
            headers := make([]string, 0, len(cols))
            for _, c := range cols { headers = append(headers, c.Name) }
            writeRow(headers)
            wroteHeader = true
        }
        vals := make([]string, 0, len(row))
        for _, c := range cols {
            v := row[c.Name]
            s := fmt.Sprintf("%v", v)
            s = strings.ReplaceAll(s, "\n", " ")
            s = strings.ReplaceAll(s, "\r", " ")
            s = strings.ReplaceAll(s, ",", "\t")
            // 简单转义双引号
            s = strings.ReplaceAll(s, "\"", "'")
            // 包裹以避免逗号问题
            vals = append(vals, fmt.Sprintf("\"%s\"", s))
        }
        writeRow(vals)
        return nil
    })
    if err != nil { /* 流式响应，错误不再写JSON */ }
}

// ExportSQL 导出INSERT脚本
func (s DbServiceImpl) ExportSQL(c *gin.Context, dbId uint, dbName, sql, table, where string) {
    conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    query := strings.TrimSpace(sql)
    if query == "" {
        if table == "" { result.Failed(c, int(result.ApiCode.FAILED), "sql 或 table 必填"); return }
        q := conn.GetDialect().Quoter()
        query = fmt.Sprintf("SELECT * FROM %s", q.Quote(table))
        if strings.TrimSpace(where) != "" { query += " WHERE " + where }
    }

    c.Header("Content-Type", "text/plain; charset=utf-8")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=export_%d.sql", time.Now().Unix()))
    w := c.Writer

    var colNames []string
    _, _ = conn.WalkQueryRows(c.Request.Context(), query, func(row map[string]any, cols []*dbi.QueryColumn) error {
        if colNames == nil {
            colNames = make([]string, len(cols))
            for i, c := range cols { colNames[i] = c.Name }
        }
        values := make([]string, len(colNames))
        for i, name := range colNames {
            v := row[name]
            // 使用 dbi 的数据类型转 SQL 字面量会更准确，这里简单处理
            switch vv := v.(type) {
            case nil:
                values[i] = "NULL"
            case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
                values[i] = fmt.Sprintf("%v", vv)
            case time.Time:
                values[i] = fmt.Sprintf("'%s'", vv.Format("2006-01-02 15:04:05"))
            default:
                s := fmt.Sprintf("%v", vv)
                s = strings.ReplaceAll(s, "'", "''")
                values[i] = fmt.Sprintf("'%s'", s)
            }
        }
        q := conn.GetDialect().Quoter()
        stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);\\n", q.Quote(tableNameForSql(table, query)), strings.Join(quoteColsWithQuoter(q, colNames), ","), strings.Join(values, ","))
        _, _ = w.Write([]byte(stmt))
        return nil
    })
}

func quoteColsWithQuoter(q dbi.Quoter, cols []string) []string {
    res := make([]string, len(cols))
    for i, c := range cols { res[i] = q.Quote(c) }
    return res
}

func tableNameForSql(table, sql string) string {
    if table != "" { return table }
    // 尝试从SQL中提取表名（与 extractTableName 逻辑一致的无接收者版本）
    trimmed := strings.TrimSpace(sql)
    upperSql := strings.ToUpper(trimmed)
    upperWords := strings.Fields(upperSql)
    origWords := strings.Fields(trimmed)
    for i, word := range upperWords {
        if (word == "FROM" || word == "INTO" || word == "UPDATE") && i+1 < len(origWords) {
            t := origWords[i+1]
            t = strings.Trim(t, "`\" ,")
            return t
        }
    }
    return ""
}

// ExportDatabaseStructure 导出整个数据库的结构（所有表的DDL）
func (s DbServiceImpl) ExportDatabaseStructure(c *gin.Context, dbId uint, dbName string) {
    conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
    if err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), err.Error())
        return
    }

    // 获取所有表
    metadata := conn.GetMetadata()
    tables, err := metadata.GetTables()
    if err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取表列表失败: %v", err))
        return
    }

    // 设置响应头
    c.Header("Content-Type", "text/plain; charset=utf-8")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s_structure_%d.sql", dbName, time.Now().Unix()))
    w := c.Writer

    // 写入数据库注释
    _, _ = w.Write([]byte(fmt.Sprintf("-- Database: %s\n", dbName)))
    _, _ = w.Write([]byte(fmt.Sprintf("-- Export Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))))
    _, _ = w.Write([]byte(fmt.Sprintf("-- Tables: %d\n\n", len(tables))))

    // 遍历所有表，导出DDL
    for _, table := range tables {
        ddl, err := metadata.GetTableDDL(table.TableName, false)
        if err != nil {
            _, _ = w.Write([]byte(fmt.Sprintf("-- Error getting DDL for table %s: %v\n\n", table.TableName, err)))
            continue
        }
        _, _ = w.Write([]byte(fmt.Sprintf("-- Table: %s\n", table.TableName)))
        _, _ = w.Write([]byte(ddl))
        _, _ = w.Write([]byte(";\n\n"))
    }
}

// ExportDatabaseStructureAndData 导出整个数据库的结构和数据
func (s DbServiceImpl) ExportDatabaseStructureAndData(c *gin.Context, dbId uint, dbName string) {
    conn, err := s.getDbConn(c.Request.Context(), dbId, dbName)
    if err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), err.Error())
        return
    }

    // 获取所有表
    metadata := conn.GetMetadata()
    tables, err := metadata.GetTables()
    if err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取表列表失败: %v", err))
        return
    }

    // 设置响应头
    c.Header("Content-Type", "text/plain; charset=utf-8")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s_full_%d.sql", dbName, time.Now().Unix()))
    w := c.Writer

    // 写入数据库注释
    _, _ = w.Write([]byte(fmt.Sprintf("-- Database: %s\n", dbName)))
    _, _ = w.Write([]byte(fmt.Sprintf("-- Export Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))))
    _, _ = w.Write([]byte(fmt.Sprintf("-- Tables: %d\n", len(tables))))
    _, _ = w.Write([]byte("-- Export Type: Structure + Data\n\n"))

    q := conn.GetDialect().Quoter()

    // 遍历所有表
    for _, table := range tables {
        // 导出表结构
        ddl, err := metadata.GetTableDDL(table.TableName, false)
        if err != nil {
            _, _ = w.Write([]byte(fmt.Sprintf("-- Error getting DDL for table %s: %v\n\n", table.TableName, err)))
            continue
        }
        _, _ = w.Write([]byte(fmt.Sprintf("-- Table structure for %s\n", table.TableName)))
        _, _ = w.Write([]byte(ddl))
        _, _ = w.Write([]byte(";\n\n"))

        // 导出表数据
        _, _ = w.Write([]byte(fmt.Sprintf("-- Data for table %s\n", table.TableName)))
        query := fmt.Sprintf("SELECT * FROM %s", q.Quote(table.TableName))

        var colNames []string
        rowCount := 0
        _, _ = conn.WalkQueryRows(c.Request.Context(), query, func(row map[string]any, cols []*dbi.QueryColumn) error {
            if colNames == nil {
                colNames = make([]string, len(cols))
                for i, c := range cols {
                    colNames[i] = c.Name
                }
            }
            values := make([]string, len(colNames))
            for i, name := range colNames {
                v := row[name]
                switch vv := v.(type) {
                case nil:
                    values[i] = "NULL"
                case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
                    values[i] = fmt.Sprintf("%v", vv)
                case time.Time:
                    values[i] = fmt.Sprintf("'%s'", vv.Format("2006-01-02 15:04:05"))
                default:
                    s := fmt.Sprintf("%v", vv)
                    s = strings.ReplaceAll(s, "'", "''")
                    values[i] = fmt.Sprintf("'%s'", s)
                }
            }
            stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);\n",
                q.Quote(table.TableName),
                strings.Join(quoteColsWithQuoter(q, colNames), ","),
                strings.Join(values, ","))
            _, _ = w.Write([]byte(stmt))
            rowCount++
            return nil
        })

        if rowCount == 0 {
            _, _ = w.Write([]byte("-- No data\n"))
        }
        _, _ = w.Write([]byte("\n"))
    }
}

// ---------- Create Database & Create Table ----------

func (s DbServiceImpl) CreateDatabase(c *gin.Context, id uint, req *model.CreateDatabaseReq) {
    // 解析实例ID（数据库ID -> 实例ID 兼容）
    instanceId := id
    if db, err := dao.GetDbById(id); err == nil && db.ID > 0 {
        instanceId = db.InstanceId
    }
    // 连接到实例（不指定数据库）
    conn, err := DbInstanceService().GetDbConn(c.Request.Context(), instanceId, "")
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }

    // 基本校验与构建DDL
    if !isValidIdentifier(req.DbName) {
        result.Failed(c, int(result.ApiCode.FAILED), "非法库名")
        return
    }
    parts := []string{"CREATE DATABASE"}
    if req.IfNotExists { parts = append(parts, "IF NOT EXISTS") }
    parts = append(parts, fmt.Sprintf("`%s`", req.DbName))
    if req.Charset != "" { parts = append(parts, fmt.Sprintf("DEFAULT CHARACTER SET %s", req.Charset)) }
    if req.Collate != "" { parts = append(parts, fmt.Sprintf("COLLATE %s", req.Collate)) }
    ddl := strings.Join(parts, " ")

    if _, err := conn.ExecContext(c.Request.Context(), ddl); err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return
    }
    result.Success(c, nil)
}

func (s DbServiceImpl) CreateTable(c *gin.Context, id uint, req *model.CreateTableReq) {
    // 解析实例ID
    instanceId := id
    if db, err := dao.GetDbById(id); err == nil && db.ID > 0 {
        instanceId = db.InstanceId
    }
    // 连接到指定数据库
    conn, err := DbInstanceService().GetDbConn(c.Request.Context(), instanceId, req.Db)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }

    if !isValidIdentifier(req.TableName) { result.Failed(c, int(result.ApiCode.FAILED), "非法表名"); return }
    if len(req.Columns) == 0 { result.Failed(c, int(result.ApiCode.FAILED), "列定义不能为空"); return }

    var colDefs []string
    var pkCols []string
    for _, col := range req.Columns {
        if !isValidIdentifier(col.Name) { result.Failed(c, int(result.ApiCode.FAILED), "非法列名: "+col.Name); return }
        def := fmt.Sprintf("`%s` %s", col.Name, col.Type)
        if col.AutoIncrement { def += " AUTO_INCREMENT" }
        if col.Nullable { def += " NULL" } else { def += " NOT NULL" }
        if col.Default != nil {
            def += " DEFAULT " + formatDefaultValue(*col.Default)
        }
        if col.Comment != "" {
            def += " COMMENT '" + escapeSingle(col.Comment) + "'"
        }
        colDefs = append(colDefs, def)
        if col.PrimaryKey { pkCols = append(pkCols, fmt.Sprintf("`%s`", col.Name)) }
    }
    if len(pkCols) > 0 {
        colDefs = append(colDefs, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pkCols, ",")))
    }

    // 索引定义
    for _, idx := range req.Indexes {
        if len(idx.Columns) == 0 { continue }
        var cols []string
        for _, cn := range idx.Columns {
            if !isValidIdentifier(cn) { result.Failed(c, int(result.ApiCode.FAILED), "非法索引列名: "+cn); return }
            cols = append(cols, fmt.Sprintf("`%s`", cn))
        }
        key := "KEY"
        if idx.Unique { key = "UNIQUE KEY" }
        name := idx.Name
        if name == "" { name = strings.Join(idx.Columns, "_") }
        if !isValidIdentifier(name) { result.Failed(c, int(result.ApiCode.FAILED), "非法索引名: "+name); return }
        colDefs = append(colDefs, fmt.Sprintf("%s `%s` (%s)", key, name, strings.Join(cols, ",")))
    }

    tableOptions := []string{}
    if req.Engine != "" { tableOptions = append(tableOptions, "ENGINE="+req.Engine) }
    if req.Charset != "" { tableOptions = append(tableOptions, "DEFAULT CHARSET="+req.Charset) }

    ddl := fmt.Sprintf("CREATE TABLE `%s` (\n  %s\n) %s;", req.TableName, strings.Join(colDefs, ",\n  "), strings.Join(tableOptions, " "))
    if _, err := conn.ExecContext(c.Request.Context(), ddl); err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return
    }
    result.Success(c, nil)
}

var identRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func isValidIdentifier(s string) bool { return identRe.MatchString(s) }

func escapeSingle(s string) string { return strings.ReplaceAll(s, "'", "''") }

func formatDefaultValue(v string) string {
    up := strings.ToUpper(strings.TrimSpace(v))
    // 常见函数/关键字不过引号
    switch up {
    case "NULL", "CURRENT_TIMESTAMP", "CURRENT_DATE", "CURRENT_TIME", "NOW()":
        return up
    }
    // 数字
    if numRe.MatchString(v) { return v }
    return "'" + escapeSingle(v) + "'"
}

var numRe = regexp.MustCompile(`^[-+]?[0-9]+(\.[0-9]+)?$`)

// 生成 SQL 结果缓存键
func (s DbServiceImpl) sqlCacheKey(dto *model.DbSqlExecDto, sql string) string {
    // 优先使用 InstanceId 作为命名空间；否则用 DbId
    ns := fmt.Sprintf("inst:%d", dto.InstanceId)
    if dto.InstanceId == 0 {
        ns = fmt.Sprintf("db:%d", dto.DbId)
    }
    // 哈希SQL，避免key过长
    h := sha1.Sum([]byte(strings.TrimSpace(sql)))
    return fmt.Sprintf("db:sql:result:%s:%s:%s", ns, dto.Db, hex.EncodeToString(h[:]))
}

// 读取缓存
func (s DbServiceImpl) tryGetSqlCache(dto *model.DbSqlExecDto, sql string) *model.DbSqlExecResult {
    if redis.RedisDb == nil { return nil }
    key := s.sqlCacheKey(dto, sql)
    bs, err := redis.RedisDb.Get(context.Background(), key).Bytes()
    if err != nil { return nil }
    var res model.DbSqlExecResult
    if err := json.Unmarshal(bs, &res); err != nil { return nil }
    return &res
}

// 写入缓存
func (s DbServiceImpl) trySetSqlCache(dto *model.DbSqlExecDto, sql string, val *model.DbSqlExecResult, ttl time.Duration) error {
    if redis.RedisDb == nil { return nil }
    key := s.sqlCacheKey(dto, sql)
    bs, err := json.Marshal(val)
    if err != nil { return err }
    return redis.RedisDb.Set(context.Background(), key, bs, ttl).Err()
}

// DbService 服务实例
func DbService() IDbService {
	return DbServiceImpl{}
}

// PreCheckSQL SQL预检查，返回预估影响行数（不实际执行）
func (s DbServiceImpl) PreCheckSQL(c *gin.Context, dto *model.DbSqlExecDto) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var conn *dbi.DbConn
	var err error

	// 优先使用 InstanceId，兼容 DbId
	if dto.InstanceId > 0 {
		conn, err = DbInstanceService().GetDbConn(ctx, dto.InstanceId, dto.Db)
	} else {
		conn, err = s.getDbConn(ctx, dto.DbId, dto.Db)
	}

	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "连接数据库失败: "+err.Error())
		return
	}

	// 去除SQL首尾空格
	sql := strings.TrimSpace(dto.Sql)
	if sql == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "SQL不能为空")
		return
	}

	// 判断SQL类型
	sqlType := s.getSqlType(sql)

	// 转换为预估查询
	estimateSQL, err := s.convertToEstimateSQL(sql, sqlType)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "预检查失败: "+err.Error())
		return
	}

	// 执行预估查询
	_, rows, err := conn.Query(estimateSQL)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "预检查失败: "+err.Error())
		return
	}

	var affectedRows int64 = 0
	if rows != nil && len(rows) > 0 {
		// 获取第一行第一列的值
		if count, ok := rows[0]["count"].(int64); ok {
			affectedRows = count
		} else if count, ok := rows[0]["count"].(string); ok {
			affectedRows, _ = strconv.ParseInt(count, 10, 64)
		}
	}

	result.Success(c, gin.H{
		"sqlType":      sqlType,
		"sql":          sql,
		"estimateSQL":  estimateSQL,
		"affectedRows": affectedRows,
		"sqlTypeText":  s.getSqlTypeText(sqlType),
	})
}

// convertToEstimateSQL 将SQL转换为预估查询SQL
func (s DbServiceImpl) convertToEstimateSQL(sql string, sqlType int8) (string, error) {
	sqlUpper := strings.ToUpper(sql)

	switch sqlType {
	case model.SqlTypeQuery:
		// SELECT 查询：提取 WHERE 条件和 FROM 子句
		return s.convertSelectToCount(sql)

	case model.SqlTypeUpdate:
		// UPDATE：提取 WHERE 条件
		return s.convertUpdateToCount(sql)

	case model.SqlTypeDelete:
		// DELETE：提取 WHERE 条件
		return s.convertDeleteToCount(sql)

	case model.SqlTypeInsert:
		// INSERT：返回固定值（插入语句通常影响1行或多行）
		if strings.Contains(sqlUpper, "SELECT") {
			// INSERT INTO ... SELECT ... 需要执行子查询
			return s.convertInsertSelectToCount(sql)
		}
		// 普通 INSERT，统计 VALUES 的数量
		return "", fmt.Errorf("INSERT 语句预检查：预估插入行数")

	default:
		return "", fmt.Errorf("不支持的SQL类型")
	}
}

// convertSelectToCount 将 SELECT 转换为 COUNT 查询
func (s DbServiceImpl) convertSelectToCount(sql string) (string, error) {
	// 简单处理：将 SELECT ... FROM 替换为 SELECT COUNT(*) FROM
	// 移除 ORDER BY、LIMIT 等子句

	sqlUpper := strings.ToUpper(sql)

	// 查找 FROM 位置
	fromIdx := strings.Index(sqlUpper, "FROM")
	if fromIdx == -1 {
		return "", fmt.Errorf("未找到 FROM 子句")
	}

	// 提取 FROM 之后的部分
	fromPart := sql[fromIdx:]

	// 移除 ORDER BY、LIMIT、OFFSET
	fromPart = s.removeOrderByLimitOffset(fromPart)

	// 构造 COUNT 查询
	return "SELECT COUNT(*) AS count " + fromPart, nil
}

// convertUpdateToCount 将 UPDATE 转换为 COUNT 查询
func (s DbServiceImpl) convertUpdateToCount(sql string) (string, error) {
	// UPDATE table SET ... WHERE ... -> SELECT COUNT(*) FROM table WHERE ...

	sqlUpper := strings.ToUpper(sql)

	// 提取表名
	updateIdx := strings.Index(sqlUpper, "UPDATE")
	if updateIdx == -1 {
		return "", fmt.Errorf("非法的 UPDATE 语句")
	}

	setIdx := strings.Index(sqlUpper, "SET")
	if setIdx == -1 {
		return "", fmt.Errorf("未找到 SET 子句")
	}

	// 提取表名（UPDATE 和 SET 之间）
	tablePart := strings.TrimSpace(sql[updateIdx+6 : setIdx])

	// 提取 WHERE 子句（如果有）
	whereIdx := strings.Index(sqlUpper, "WHERE")
	wherePart := ""
	if whereIdx != -1 {
		wherePart = sql[whereIdx:]
	}

	return fmt.Sprintf("SELECT COUNT(*) AS count FROM %s %s", tablePart, wherePart), nil
}

// convertDeleteToCount 将 DELETE 转换为 COUNT 查询
func (s DbServiceImpl) convertDeleteToCount(sql string) (string, error) {
	// DELETE FROM table WHERE ... -> SELECT COUNT(*) FROM table WHERE ...

	sqlUpper := strings.ToUpper(sql)

	// 查找 FROM 位置
	fromIdx := strings.Index(sqlUpper, "FROM")
	if fromIdx == -1 {
		return "", fmt.Errorf("未找到 FROM 子句")
	}

	// 提取 FROM 之后的部分
	fromPart := sql[fromIdx:]

	return "SELECT COUNT(*) AS count " + fromPart, nil
}

// convertInsertSelectToCount 将 INSERT INTO ... SELECT 转换为 COUNT 查询
func (s DbServiceImpl) convertInsertSelectToCount(sql string) (string, error) {
	// INSERT INTO table SELECT ... -> SELECT COUNT(*) FROM (SELECT ...) AS tmp

	sqlUpper := strings.ToUpper(sql)

	// 查找 SELECT 位置
	selectIdx := strings.Index(sqlUpper, "SELECT")
	if selectIdx == -1 {
		return "", fmt.Errorf("未找到 SELECT 子句")
	}

	// 提取 SELECT 部分
	selectPart := sql[selectIdx:]

	return fmt.Sprintf("SELECT COUNT(*) AS count FROM (%s) AS tmp", selectPart), nil
}

// removeOrderByLimitOffset 移除 ORDER BY、LIMIT、OFFSET 子句
func (s DbServiceImpl) removeOrderByLimitOffset(sql string) string {
	sqlUpper := strings.ToUpper(sql)

	// 查找 ORDER BY 位置
	orderByIdx := strings.Index(sqlUpper, "ORDER BY")
	if orderByIdx != -1 {
		sql = sql[:orderByIdx]
		sqlUpper = sqlUpper[:orderByIdx]
	}

	// 查找 LIMIT 位置
	limitIdx := strings.Index(sqlUpper, "LIMIT")
	if limitIdx != -1 {
		sql = sql[:limitIdx]
		sqlUpper = sqlUpper[:limitIdx]
	}

	// 查找 OFFSET 位置
	offsetIdx := strings.Index(sqlUpper, "OFFSET")
	if offsetIdx != -1 {
		sql = sql[:offsetIdx]
	}

	return strings.TrimSpace(sql)
}

// getSqlTypeText 获取SQL类型文本
func (s DbServiceImpl) getSqlTypeText(sqlType int8) string {
	switch sqlType {
	case model.SqlTypeQuery:
		return "查询"
	case model.SqlTypeUpdate:
		return "更新"
	case model.SqlTypeDelete:
		return "删除"
	case model.SqlTypeInsert:
		return "插入"
	case model.SqlTypeDDL:
		return "DDL"
	default:
		return "未知"
	}
}

// ExportDatabaseAsync 异步导出数据库
func (s DbServiceImpl) ExportDatabaseAsync(c *gin.Context, dbId uint, dbName, exportType string) {
	// 生成任务ID
	taskId := uuid.New().String()

	// 创建任务记录
	now := time.Now()
	task := &model.ExportTask{
		TaskId:     taskId,
		DbId:       dbId,
		DbName:     dbName,
		ExportType: exportType,
		Status:     model.ExportTaskStatusPending,
		CreateTime: util.HTime{Time: now},
		UpdateTime: util.HTime{Time: now},
		Creator:    c.GetString("username"),
		CreatorId:  uint(c.GetInt("userId")),
	}

	if err := dao.CreateExportTask(task); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "创建导出任务失败: "+err.Error())
		return
	}

	// 启动后台goroutine执行导出
	go s.executeExportTask(taskId, dbId, dbName, exportType)

	// 立即返回任务ID
	result.Success(c, gin.H{
		"taskId": taskId,
		"status": model.ExportTaskStatusPending,
		"message": "导出任务已创建，正在后台处理",
	})
}

// executeExportTask 在后台执行导出任务
func (s DbServiceImpl) executeExportTask(taskId string, dbId uint, dbName, exportType string) {
	// 获取任务
	task, err := dao.GetExportTaskByTaskId(taskId)
	if err != nil {
		return
	}

	// 更新状态为处理中
	startTime := time.Now()
	task.Status = model.ExportTaskStatusProcessing
	task.StartTime = &startTime
	task.UpdateTime = util.HTime{Time: startTime}
	dao.UpdateExportTask(task)

	// 创建临时文件目录
	exportDir := "./data/exports"
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		s.markExportTaskFailed(task, "创建导出目录失败: "+err.Error())
		return
	}

	// 生成文件路径
	fileName := fmt.Sprintf("%s_%s_%d.sql", dbName, exportType, time.Now().Unix())
	filePath := filepath.Join(exportDir, fileName)

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		s.markExportTaskFailed(task, "创建导出文件失败: "+err.Error())
		return
	}
	defer file.Close()

	// 执行导出
	ctx := context.Background()
	conn, err := s.getDbConn(ctx, dbId, dbName)
	if err != nil {
		s.markExportTaskFailed(task, "连接数据库失败: "+err.Error())
		return
	}

	// 获取所有表
	metadata := conn.GetMetadata()
	tables, err := metadata.GetTables()
	if err != nil {
		s.markExportTaskFailed(task, "获取表列表失败: "+err.Error())
		return
	}

	// 写入文件头
	file.WriteString(fmt.Sprintf("-- Database: %s\n", dbName))
	file.WriteString(fmt.Sprintf("-- Export Time: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	file.WriteString(fmt.Sprintf("-- Tables: %d\n", len(tables)))
	file.WriteString(fmt.Sprintf("-- Export Type: %s\n\n", exportType))

	q := conn.GetDialect().Quoter()

	// 遍历所有表
	for _, table := range tables {
		// 导出表结构
		ddl, err := metadata.GetTableDDL(table.TableName, false)
		if err != nil {
			file.WriteString(fmt.Sprintf("-- Error getting DDL for table %s: %v\n\n", table.TableName, err))
			continue
		}
		file.WriteString(fmt.Sprintf("-- Table structure for %s\n", table.TableName))
		file.WriteString(ddl)
		file.WriteString(";\n\n")

		// 如果是完整导出，还要导出数据
		if exportType == "full" {
			file.WriteString(fmt.Sprintf("-- Data for table %s\n", table.TableName))
			query := fmt.Sprintf("SELECT * FROM %s", q.Quote(table.TableName))

			var colNames []string
			rowCount := 0
			conn.WalkQueryRows(ctx, query, func(row map[string]any, cols []*dbi.QueryColumn) error {
				if colNames == nil {
					colNames = make([]string, len(cols))
					for i, c := range cols {
						colNames[i] = c.Name
					}
				}
				values := make([]string, len(colNames))
				for i, name := range colNames {
					v := row[name]
					switch vv := v.(type) {
					case nil:
						values[i] = "NULL"
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
						values[i] = fmt.Sprintf("%v", vv)
					case time.Time:
						values[i] = fmt.Sprintf("'%s'", vv.Format("2006-01-02 15:04:05"))
					default:
						s := fmt.Sprintf("%v", vv)
						s = strings.ReplaceAll(s, "'", "''")
						values[i] = fmt.Sprintf("'%s'", s)
					}
				}
				stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);\n",
					q.Quote(table.TableName),
					strings.Join(quoteColsWithQuoter(q, colNames), ","),
					strings.Join(values, ","))
				file.WriteString(stmt)
				rowCount++
				return nil
			})

			if rowCount == 0 {
				file.WriteString("-- No data\n")
			}
			file.WriteString("\n")
		}
	}

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		s.markExportTaskFailed(task, "获取文件信息失败: "+err.Error())
		return
	}

	// 更新任务状态为完成
	endTime := time.Now()
	task.Status = model.ExportTaskStatusCompleted
	task.FilePath = filePath
	task.FileSize = fileInfo.Size()
	task.EndTime = &endTime
	task.UpdateTime = util.HTime{Time: endTime}
	dao.UpdateExportTask(task)
}

// markExportTaskFailed 标记导出任务为失败
func (s DbServiceImpl) markExportTaskFailed(task *model.ExportTask, errMsg string) {
	endTime := time.Now()
	task.Status = model.ExportTaskStatusFailed
	task.ErrorMessage = errMsg
	task.EndTime = &endTime
	task.UpdateTime = util.HTime{Time: endTime}
	dao.UpdateExportTask(task)
}

// GetExportTaskStatus 查询导出任务状态
func (s DbServiceImpl) GetExportTaskStatus(c *gin.Context, taskId string) {
	task, err := dao.GetExportTaskByTaskId(taskId)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "任务不存在")
		return
	}

	vo := model.ExportTaskVo{
		TaskId:       task.TaskId,
		DbName:       task.DbName,
		ExportType:   task.ExportType,
		Status:       task.Status,
		FileSize:     task.FileSize,
		ErrorMessage: task.ErrorMessage,
		StartTime:    task.StartTime,
		EndTime:      task.EndTime,
		CreateTime:   task.CreateTime,
	}

	result.Success(c, vo)
}

// DownloadExportFile 下载导出文件
func (s DbServiceImpl) DownloadExportFile(c *gin.Context, taskId string) {
	task, err := dao.GetExportTaskByTaskId(taskId)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "任务不存在")
		return
	}

	if task.Status != model.ExportTaskStatusCompleted {
		result.Failed(c, int(result.ApiCode.FAILED), "任务未完成，无法下载")
		return
	}

	if task.FilePath == "" || !fileExists(task.FilePath) {
		result.Failed(c, int(result.ApiCode.FAILED), "导出文件不存在")
		return
	}

	// 设置响应头
	fileName := filepath.Base(task.FilePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Length", fmt.Sprintf("%d", task.FileSize))

	// 发送文件
	c.File(task.FilePath)
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
