package db

import (
    "dodevops-api/api/db/mysql/controller"
    unified "dodevops-api/api/db/unified/controller"
    esController "dodevops-api/api/db/elasticsearch/controller"
    mongoController "dodevops-api/api/db/mongo/controller"
    redisController "dodevops-api/api/db/redis/controller"
    "dodevops-api/middleware"
    "github.com/gin-gonic/gin"
)

// InitDbRouter 初始化数据库管理路由
func InitDbRouter(router *gin.Engine) {
	// 数据库管理路由组
	dbGroup := router.Group("/api/v1/db")
	dbGroup.Use(middleware.AuthMiddleware()) // 使用JWT认证中间件
	dbGroup.Use(middleware.LogMiddleware())  // 使用操作日志中间件

	{
		// 数据库实例管理
		instances := dbGroup.Group("/instances")
		{
			instances.GET("", controller.GetDbInstanceList)           // 获取实例列表
			instances.GET("/:id", controller.GetDbInstanceById)       // 获取实例详情
			instances.POST("", controller.CreateDbInstance)           // 创建实例
			instances.PUT("", controller.UpdateDbInstance)            // 更新实例
			instances.DELETE("/:id", controller.DeleteDbInstance)     // 删除实例
			instances.POST("/test-conn", controller.TestConnection)   // 测试连接
			instances.GET("/:id/databases", controller.GetDatabases)  // 获取数据库列表
		}



		// 数据库管理
		databases := dbGroup.Group("/databases")
		{
			databases.GET("", controller.GetDbList)                   // 获取数据库列表
			databases.GET("/:id", controller.GetDbById)               // 获取数据库详情
			databases.POST("", controller.CreateDb)                   // 创建数据库
			databases.PUT("", controller.UpdateDb)                    // 更新数据库
			databases.DELETE("/:id", controller.DeleteDb)             // 删除数据库

			// 表管理
			databases.GET("/:id/tables", controller.GetTables)        // 获取表列表
			databases.GET("/:id/columns", controller.GetTableColumns) // 获取表列信息
			databases.GET("/:id/indexes", controller.GetTableIndexes) // 获取表索引
			databases.GET("/:id/ddl", controller.GetTableDDL)         // 获取表DDL
			databases.GET("/:id/approx-count", controller.GetApproxRowCount) // 快速行数
			databases.GET("/:id/table-status", controller.GetTableStatus)     // 表状态
			databases.GET("/:id/approx-counts", controller.GetApproxRowCounts) // 批量快速行数
			databases.GET("/:id/row", controller.GetRow)                    // 获取单行
			databases.DELETE("/:id/row", controller.DeleteRow)              // 删除单行
			databases.GET("/:id/export-csv", controller.ExportCSV)          // 导出CSV
			databases.POST("/:id/export-csv", controller.ExportCSVPost)     // 导出CSV(POST)
			databases.GET("/:id/export-sql", controller.ExportSQL)          // 导出SQL脚本
			databases.POST("/:id/export-sql", controller.ExportSQLPost)     // 导出SQL脚本(POST)
			databases.GET("/:id/export-database-structure", controller.ExportDatabaseStructure) // 导出整个数据库结构
			databases.GET("/:id/export-database-full", controller.ExportDatabaseFull)         // 导出整个数据库结构+数据(同步)
			databases.POST("/:id/export-database-full-async", controller.ExportDatabaseFullAsync) // 异步导出数据库
			// 创建库与表
			databases.POST("/:id/create-database", controller.CreateDatabase)
			databases.POST("/:id/create-table", controller.CreateTable)
			databases.POST("/copy-table", controller.CopyTable)       // 复制表

			// SQL执行
			databases.POST("/exec-sql", controller.ExecSql)           // 执行SQL
			databases.POST("/:id/exec-sql-file", controller.ExecSqlFile) // 执行SQL文件
			databases.GET("/:id/hint-tables", controller.GetHintTables) // 获取提示表
		}

		// 导出任务管理
		exportTasks := dbGroup.Group("/export-tasks")
		{
			exportTasks.GET("/:taskId", controller.GetExportTaskStatus)         // 查询导出任务状态
			exportTasks.GET("/:taskId/download", controller.DownloadExportFile) // 下载导出文件
		}

		// SQL模板管理
		sqlTemplates := dbGroup.Group("/sql-templates")
		{
			sqlTemplates.GET("", controller.GetDbSqlList)             // 获取SQL模板列表
			sqlTemplates.GET("/names", controller.GetDbSqlNames)      // 获取SQL模板名称列表
			sqlTemplates.GET("/:id", controller.GetDbSqlById)         // 获取SQL模板详情
			sqlTemplates.POST("", controller.CreateDbSql)             // 创建SQL模板
			sqlTemplates.PUT("", controller.UpdateDbSql)              // 更新SQL模板
			sqlTemplates.DELETE("/:id", controller.DeleteDbSql)       // 删除SQL模板
		}

		// SQL执行历史
		sqlExecs := dbGroup.Group("/sql-execs")
		{
			sqlExecs.GET("", controller.GetDbSqlExecList)                 // 获取执行历史列表
			sqlExecs.GET("/:id", controller.GetDbSqlExecById)             // 获取执行历史详情
			sqlExecs.GET("/statistics", controller.GetDbSqlExecStatistics) // 获取执行统计
		}

		// SQL工具
		dbGroup.POST("/format-sql", controller.FormatSQL)    // SQL格式化
		dbGroup.POST("/pre-check-sql", controller.PreCheckSQL) // SQL预检查

		// 统一实例管理（新增 - 使用独立路径避免冲突）
		unifiedGroup := dbGroup.Group("/unified")
		{
			unifiedGroup.GET("/instances", unified.GetAllInstancesList)                        // 统一查询所有类型实例
			unifiedGroup.GET("/instances/:dbType/:id", unified.GetInstanceById)                // 根据类型和ID查询实例
			unifiedGroup.GET("/instances/stats", unified.GetInstancesStats)                    // 实例统计信息
			unifiedGroup.GET("/instances/:dbType/:id/password", unified.GetInstancePassword)   // 获取实例密码（敏感操作）
		}
	}

	// Elasticsearch 管理路由
	esGroup := router.Group("/api/v1/elasticsearch")
	esGroup.Use(middleware.AuthMiddleware())
	esGroup.Use(middleware.LogMiddleware())
	{
		// 实例管理
		esGroup.GET("", esController.GetEsInstanceList)
		esGroup.GET(":id", esController.GetEsInstanceById)
		esGroup.POST("", esController.CreateEsInstance)
		esGroup.PUT("", esController.UpdateEsInstance)
		esGroup.DELETE(":id", esController.DeleteEsInstance)
		esGroup.POST("/test-conn", esController.TestConnection)

		// 反向代理 ES API
		esGroup.Any(":id/proxy/*path", esController.Proxy)
	}

	// MongoDB 管理路由
	mongoGroup := router.Group("/api/v1/mongodb")
	mongoGroup.Use(middleware.AuthMiddleware())
	mongoGroup.Use(middleware.LogMiddleware())
	{
		// 实例管理
		mongoGroup.GET("", mongoController.GetMongoInstanceList)
		mongoGroup.GET(":id", mongoController.GetMongoInstanceById)
		mongoGroup.POST("", mongoController.CreateMongoInstance)
		mongoGroup.PUT("", mongoController.UpdateMongoInstance)
		mongoGroup.DELETE(":id", mongoController.DeleteMongoInstance)
		mongoGroup.POST("/test-conn", mongoController.TestConnection)

		// MongoDB 操作
		mongoGroup.GET(":id/databases", mongoController.GetDatabases)
		mongoGroup.GET(":id/collections", mongoController.GetCollections)
		mongoGroup.POST(":id/run-command", mongoController.RunCommand)
		mongoGroup.POST(":id/command/find", mongoController.Find)
		mongoGroup.POST(":id/command/update-by-id", mongoController.UpdateById)
		mongoGroup.POST(":id/command/delete-by-id", mongoController.DeleteById)
		mongoGroup.POST(":id/command/insert", mongoController.InsertOne)
	}

	// Redis 管理路由
	redisGroup := router.Group("/api/v1/redis")
	redisGroup.Use(middleware.AuthMiddleware())
	redisGroup.Use(middleware.LogMiddleware())
	{
		// 实例管理
		redisGroup.GET("", redisController.GetRedisInstanceList)
		redisGroup.GET(":id", redisController.GetRedisInstanceById)
		redisGroup.POST("", redisController.SaveRedisInstance)
		redisGroup.PUT("", redisController.SaveRedisInstance)
		redisGroup.DELETE(":id", redisController.DeleteRedisInstance)
		redisGroup.POST("/test-conn", redisController.TestConnection)

		// 基于实例ID/DB的业务接口
		redisGroup.POST(":id/scan", redisController.ScanRedisKeysAllDb)
		redisGroup.GET(":id/dbs", redisController.GetRedisDatabases)
		redisGroup.GET(":id/info", redisController.GetRedisInfo)
		redisGroup.GET(":id/cluster-info", redisController.GetRedisClusterInfo)
		redisGroup.POST(":id/:db/scan", redisController.ScanRedisKeys)
		redisGroup.GET(":id/:db/key-info", redisController.GetRedisKeyInfo)
		redisGroup.GET(":id/:db/key-ttl", redisController.GetRedisKeyTtl)
		redisGroup.GET(":id/:db/key-memuse", redisController.GetRedisKeyMemuse)
		redisGroup.POST(":id/:db/run-cmd", redisController.RunRedisCmd)
	}
}
