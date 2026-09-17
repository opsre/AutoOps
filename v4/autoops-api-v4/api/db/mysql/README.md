# MySQL数据库管理功能说明

## 概述

本项目已成功集成了完整的MySQL数据库管理功能,改编自开源项目Mayfly-Go的数据库管理模块。该功能提供了企业级的数据库运维管理能力。

## 功能特性

### 1. 数据库实例管理
- 多数据库类型支持 (MySQL, PostgreSQL, Oracle等)
- 实例连接信息管理
- 连接测试
- SSH隧道支持
- 连接池管理

### 2. 数据库管理
- 数据库(Schema)的增删改查
- 关联实例管理
- 状态管理

### 3. 表管理
- 查看表列表
- 查看表结构(列信息)
- 查看表索引
- 获取表DDL语句
- 复制表(结构+数据)

### 4. SQL执行
- SQL语句执行(查询/更新)
- 执行历史记录
- 执行结果查看
- 执行时长统计
- 回滚支持(记录旧值)

### 5. SQL模板管理
- SQL模板保存
- 模板分类
- 快速执行

### 6. SQL编辑器支持
- 表名自动补全
- 字段名自动补全

## 项目结构

```
api/db/mysql/
├── controller/              # 控制层
│   ├── dbInstanceController.go    # 实例管理控制器
│   ├── dbController.go            # 数据库管理控制器
│   └── dbSqlController.go         # SQL管理控制器
├── service/                # 服务层
│   ├── dbInstanceService.go       # 实例服务
│   ├── dbService.go               # 数据库服务
│   ├── dbSqlService.go            # SQL模板服务
│   └── dbSqlExecService.go        # SQL执行服务
├── dao/                    # 数据访问层
│   ├── dbInstanceDao.go           # 实例DAO
│   ├── dbDao.go                   # 数据库DAO
│   ├── dbSqlDao.go                # SQL模板DAO
│   └── dbSqlExecDao.go            # SQL执行DAO
└── model/                  # 模型层
    ├── dbInstance.go              # 实例模型
    ├── db.go                      # 数据库模型
    ├── dbSql.go                   # SQL模板模型
    ├── dbSqlExec.go               # SQL执行记录模型
    └── metadata.go                # 元数据模型

pkg/dbm/                    # 数据库核心抽象层
├── dbi/                    # 数据库接口定义
│   ├── metadata.go         # 元数据访问接口
│   ├── dialect.go          # 数据库方言接口
│   ├── conn.go             # 数据库连接
│   └── db_info.go          # 数据库信息
├── mysql/                  # MySQL实现
│   ├── meta.go             # MySQL元信息
│   ├── metadata.go         # MySQL元数据实现
│   ├── dialect.go          # MySQL方言实现
│   └── sqlgen.go           # MySQL SQL生成器
└── dbm.go                  # 连接池管理

pkg/pool/                   # 资源池
└── pool.go                 # 通用资源池实现

router/db/                  # 路由
└── db.go                   # 数据库管理路由
```

## API端点列表

### 数据库实例管理

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/api/v1/db/instances` | 获取实例列表 |
| GET | `/api/v1/db/instances/:id` | 获取实例详情 |
| POST | `/api/v1/db/instances` | 创建实例 |
| PUT | `/api/v1/db/instances` | 更新实例 |
| DELETE | `/api/v1/db/instances/:id` | 删除实例 |
| POST | `/api/v1/db/instances/test-conn` | 测试连接 |
| GET | `/api/v1/db/instances/:id/databases` | 获取数据库列表 |

### 数据库管理

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/api/v1/db/databases` | 获取数据库列表 |
| GET | `/api/v1/db/databases/:id` | 获取数据库详情 |
| POST | `/api/v1/db/databases` | 创建数据库 |
| PUT | `/api/v1/db/databases` | 更新数据库 |
| DELETE | `/api/v1/db/databases/:id` | 删除数据库 |
| GET | `/api/v1/db/databases/:id/tables` | 获取表列表 |
| GET | `/api/v1/db/databases/:id/columns` | 获取表列信息 |
| GET | `/api/v1/db/databases/:id/indexes` | 获取表索引 |
| GET | `/api/v1/db/databases/:id/ddl` | 获取表DDL |
| GET | `/api/v1/db/databases/:id/approx-count` | 快速行数(近似) |
| GET | `/api/v1/db/databases/:id/table-status` | 表状态(近似行数/大小) |
| GET | `/api/v1/db/databases/:id/approx-counts` | 批量快速行数 |
| POST | `/api/v1/db/databases/copy-table` | 复制表 |
| POST | `/api/v1/db/databases/exec-sql` | 执行SQL |
| GET | `/api/v1/db/databases/:id/hint-tables` | 获取提示表 |

### SQL模板管理

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/api/v1/db/sql-templates` | 获取SQL模板列表 |
| GET | `/api/v1/db/sql-templates/names` | 获取模板名称列表 |
| GET | `/api/v1/db/sql-templates/:id` | 获取模板详情 |
| POST | `/api/v1/db/sql-templates` | 创建SQL模板 |
| PUT | `/api/v1/db/sql-templates` | 更新SQL模板 |
| DELETE | `/api/v1/db/sql-templates/:id` | 删除SQL模板 |

### SQL执行历史

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/api/v1/db/sql-execs` | 获取执行历史列表 |
| GET | `/api/v1/db/sql-execs/:id` | 获取执行历史详情 |
| GET | `/api/v1/db/sql-execs/statistics` | 获取执行统计 |

## 使用示例

### 1. 创建数据库实例

```bash
POST /api/v1/db/instances
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "生产环境MySQL",
  "type": "mysql",
  "host": "192.168.1.100",
  "port": 3306,
  "username": "root",
  "password": "password",
  "network": "tcp",
  "params": "charset=utf8mb4&parseTime=True",
  "remark": "生产环境主库"
}
```

### 2. 测试连接

```bash
POST /api/v1/db/instances/test-conn
Content-Type: application/json
Authorization: Bearer {token}

{
  "type": "mysql",
  "host": "192.168.1.100",
  "port": 3306,
  "username": "root",
  "password": "password",
  "database": "test"
}
```

### 3. 创建数据库

```bash
POST /api/v1/db/databases
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "测试数据库",
  "database": "test_db",
  "instanceId": 1,
  "remark": "测试环境数据库"
}
```

### 4. 执行SQL

```bash
POST /api/v1/db/databases/exec-sql
Content-Type: application/json
Authorization: Bearer {token}

{
  "dbId": 1,
  "db": "test_db",
  "sql": "SELECT * FROM users",
  "remark": "查询用户列表",
  "useCache": true,           // 可选：结果缓存
  "cacheSeconds": 120,        // 可选：缓存过期(秒)
  "maxRows": 500,             // 可选：最大返回行数(0=不限制)
  "timeoutSeconds": 60        // 可选：超时秒数(0=默认60)
}
```

### 5. 获取表列表

```bash
GET /api/v1/db/databases/1/tables?db=test_db
Authorization: Bearer {token}
```

### 6. 获取表结构

```bash
GET /api/v1/db/databases/1/columns?db=test_db&table=users
Authorization: Bearer {token}
```

### 7. 快速行数与表状态

```bash
GET /api/v1/db/databases/1/approx-count?db=test_db&table=users
Authorization: Bearer {token}
```

响应示例:

```json
{ "code": 200, "data": { "table": "users", "rows": 123456 } }
```

```bash
GET /api/v1/db/databases/1/table-status?db=test_db
Authorization: Bearer {token}
```

响应示例(部分字段):

```json
{
  "code": 200,
  "data": [
    { "tableName": "users", "rows": 123456, "dataLength": 1048576, "indexLength": 524288, "comment": "用户表" }
  ]
}
```

```bash
GET /api/v1/db/databases/1/approx-counts?db=test_db&tables=users,orders
Authorization: Bearer {token}
```

响应示例:

```json
{ "code": 200, "data": { "users": 123456, "orders": 98765 } }
```

## 数据库表结构

执行项目启动时会自动创建以下表:

1. **db_instance** - 数据库实例表
2. **db** - 数据库表
3. **db_sql** - SQL模板表
4. **db_sql_exec** - SQL执行历史表

## 核心特性

### 1. 连接池管理
- 自动连接池复用
- 30分钟无活动自动释放
- 支持多实例多数据库并发访问

### 2. 多数据库支持
通过抽象的Dialect模式,理论上支持:
- MySQL/MariaDB
- PostgreSQL
- Oracle
- SQL Server
- 达梦数据库(DM)
- SQLite

当前已完整实现MySQL支持,其他数据库需要验证相关实现。

### 3. 安全特性
- 密码需加密存储(TODO: 当前未实现,建议使用AES加密)
- SQL执行记录审计
- UPDATE/DELETE操作记录旧值支持回滚
- JWT认证保护所有接口

### 4. 性能优化
- 连接池复用避免频繁创建连接
- 元数据查询可缓存
- 分页查询支持

## 待优化项

1. **密码加密**: 当前密码明文存储在数据库,建议实现AES加密
2. **SQL解析**: 可以集成完整的SQL解析器,提供更精确的SQL分析
3. **审批流程**: 对危险SQL(DELETE/DROP等)集成审批工作流
4. **查询限制**: 对大查询结果集实现限制,避免内存溢出
5. **慢查询分析**: 记录慢查询并提供分析功能
6. **数据导入导出**: 实现数据库备份/恢复功能
7. **权限控制**: 实现细粒度的数据库/表级别权限控制

## 依赖说明

### Go依赖
- `github.com/gin-gonic/gin` - Web框架
- `gorm.io/gorm` - ORM框架
- `github.com/google/uuid` - UUID生成

### 内部依赖
- `dodevops-api/pkg/db` - 数据库连接
- `dodevops-api/common/result` - 统一响应格式
- `dodevops-api/middleware` - 中间件

## 开发指南

### 添加新的数据库类型支持

1. 在 `pkg/dbm/` 下创建新目录,如 `postgres/`
2. 实现 `Meta` 接口
3. 实现 `Metadata` 接口
4. 实现 `Dialect` 接口
5. 在 `meta.go` 中注册: `dbi.Register(DbTypePostgres, meta)`

### 扩展功能

如需添加新功能,遵循现有的分层架构:
1. 在 `model/` 添加数据模型
2. 在 `dao/` 添加数据访问方法
3. 在 `service/` 添加业务逻辑
4. 在 `controller/` 添加API处理器
5. 在 `router/db/db.go` 注册路由

## 问题排查

### 连接失败
1. 检查数据库实例配置是否正确
2. 检查网络连通性
3. 检查数据库用户权限
4. 查看日志文件 `logs/app.log`

### SQL执行失败
1. 检查SQL语法是否正确
2. 检查数据库权限
3. 查看执行历史表获取详细错误信息

## 致谢

本功能改编自开源项目 [Mayfly-Go](https://github.com/dromara/mayfly-go),感谢原作者的优秀设计。

## 许可证

本项目遵循原项目许可证。
