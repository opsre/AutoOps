分层监控改造总结（中文）                                                  
                                                                               
  🎯 核心改造内容                                                              
                                                            
  已完成API监控系统分层监控改造，支持将多个独立的API监控组织成父子层级关系。   
                                                                               
  典型场景：
  父监控: https://api.dding.net/system/check （域名）
  ├─ 子监控1: http://192.168.7.6:7001/system/check （IP:Port）
  ├─ 子监控2: http://192.168.7.7:7001/system/check （IP:Port）
  └─ 子监控3: http://192.168.7.8:7001/system/check （IP:Port）

  ---
  💾 后端改造内容

  1️⃣  数据模型 (api/monitor/model/apiMonitor.go)

  新增字段：
  - ParentID: 父监控ID（可为空）
  - HierarchyLevel: 层级深度（0=独立, 1=一级子, 2=二级子）
  - IsAggregate: 是否为聚合监控
  - ChildMonitors: 子监控列表关系
  - ChildCount: 子监控数量（计算字段）
  - HealthyChildren: 健康的子监控数（计算字段）

  2️⃣  查询方法 (api/monitor/dao/apiMonitorDao.go)

  新增方法：
  ✓ GetMonitorTree() - 获取父监控及子监控树
  ✓ GetTopLevelMonitors() - 获取顶级监控（默认）
  ✓ GetChildMonitors() - 获取指定父监控的直接子监控
  ✓ GetAllDescendants() - 获取所有后代（递归）
  ✓ CheckCircularDependency() - 检查循环依赖
  ✓ GetChildrenStatsInBatch() - 批量获取子监控统计
  ✓ UpdateParentAggregateStats() - 更新父监控聚合数据

  增强的 GetList() 方法支持：
  - parentId 参数：查询某个父监控的子监控
  - onlyTopLevel 参数：只查询顶级监控（默认为true，向下兼容）
  - includeChildren 参数：返回结果中包含子监控列表

  3️⃣  业务服务 (api/monitor/service/apiMonitorService.go)

  新增方法：
  ✓ AddChildMonitor() - 为监控添加子监控
  ✓ ConvertToChildMonitor() - 将独立监控转换为子监控
  ✓ GetMonitorHierarchy() - 获取监控树形结构
  ✓ DeleteMonitorWithChildren() - 删除监控（级联/孤立/防止）
  ✓ CalculateParentHealthStatus() - 计算父监控聚合健康状态
  ✓ UpdateAllParentAggregateStats() - 定时更新所有父监控聚合数据
  ✓ RemoveParentRelationship() - 移除子监控与父监控的关系

  验证机制：
  - ❌ 禁止循环依赖（A→B→C→A 会被拒绝）
  - ❌ 最多3层级限制（顶级 → 一级子 → 二级子）
  - ✅ 自动合并子监控健康状态

  4️⃣  API 接口 (api/monitor/controller/apiMonitorController.go)

  新增4个接口：

  1️⃣  POST /api/v1/monitor/api-monitors/{id}/add-child
     功能：为父监控添加子监控
     请求体：{ "childMonitorId": 123 }

  2️⃣  GET /api/v1/monitor/api-monitors/{id}/hierarchy?depth=2
     功能：获取监控的树形结构
     参数：depth=深度限制（可选）

  3️⃣  GET /api/v1/monitor/api-monitors/top-level?page=1&pageSize=20
     功能：获取顶级监控列表（支持分页/过滤）
     参数：pageSize, groupId, keyword, includeChildren

  4️⃣  DELETE /api/v1/monitor/api-monitors/{id}/parent
     功能：移除子监控与父监控的关系，使其变为独立监控

  5️⃣  定时任务 (api/monitor/service/apiMonitorScheduler.go)

  - 每60秒自动更新所有父监控的聚合健康状态
  - 自动计算健康的子监控数量
  - 不阻塞主检测流程

  ---
  🎨 前端配置修改方案

  A. 监控列表展示

  原有展示（保持不变）：
  监控列表（平铺）
  ├─ api.dding.net/system/check
  ├─ 192.168.7.6:7001/system/check
  ├─ 192.168.7.7:7001/system/check
  └─ ...

  改造后展示（树形结构）：
  📋 监控列表（树形）

  🏢 api.dding.net/system/check [父] ❌ 1/3健康
     ├─ 🖥️   192.168.7.6:7001 ✅ healthy
     ├─ 🖥️   192.168.7.7:7001 ❌ unhealthy
     └─ 🖥️   192.168.7.8:7001 ✅ healthy

  🔌 redis.check.net/ping [独立监控] ✅ healthy

  ... 其他顶级监控 ...

  B. 创建监控流程

  Step 1: 创建父监控
  点击「新建监控」
  ├─ 监控名称：api.dding.net
  ├─ URL：https://api.dding.net/system/check
  ├─ 业务组：选择 xxx
  ├─ 其他配置：...
  └─ 保存 ✓

  Step 2: 创建子监控
  方式A - 从列表操作：
  监控列表 → 找到「api.dding.net」→ 右键菜单 → 「添加子监控」
    ├─ 预填充 URL 基础路径：/system/check
    ├─ 子监控地址：http://192.168.7.6:7001/system/check
    ├─ 其他配置：...
    └─ 保存 ✓

  方式B - 从监控详情页：
  监控详情 → 「子监控」标签页 → 「新增子监控」按钮
    └─ 同方式A

  C. 监控详情页修改

  原有详情页：
  ┌─────────────────────────┐
  │ 监控名称：API监控        │
  │ URL：https://...        │
  │ 状态：启用              │
  │ 业务组：xxx             │
  │ 业务功能：xxx           │
  │ ... 其他配置 ...        │
  └─────────────────────────┘

  改造后详情页：
  ┌──────────────────────────────────┐
  │ 监控名称：api.dding.net           │
  │ URL：https://...                 │
  │ 状态：启用                        │
  │ 业务组：xxx                       │
  │ 业务功能：xxx                     │
  │                                  │
  │ ⬇️  新增分层信息 ⬇️                 │
  │ 父监控关系：无 [移除] [选择父]    │
  │ 层级深度：0（独立监控）           │
  │ 子监控数：3                       │
  │ 健康子监控：2/3                   │
  │ 聚合健康状态：unhealthy           │
  │                                  │
  │ 📊 子监控列表 Tab：               │
  │ ├─ 192.168.7.6:7001  ✅ healthy │
  │ ├─ 192.168.7.7:7001  ❌ unhealthy│
  │ └─ 192.168.7.8:7001  ✅ healthy │
  │ [添加子监控]                     │
  │                                  │
  │ [编辑] [删除] [保存]              │
  └──────────────────────────────────┘

  D. 监控删除/管理

  删除有子监控的父监控：
  点击「删除」 → 出现提示框

  ⚠️  该监控有 3 个子监控，请选择处理方式：

  ○ 级联删除（同时删除3个子监控及检测结果）
  ○ 孤立处理（3个子监控变为独立监控）
  ○ 防止删除（禁止删除，请先处理子监控）

  [确定] [取消]

  E. 操作菜单/工具栏修改

  监控列表右键菜单 (或操作列按钮)：
  原有菜单：
  ├─ 编辑
  ├─ 删除
  ├─ 查看详情
  └─ 测试API

  改造后菜单（对于有子监控的）：
  ├─ 编辑
  ├─ 查看详情
  ├─ ➕ 添加子监控     ← 新增
  ├─ 📊 查看子监控     ← 新增
  ├─ 🔗 选择父监控     ← 新增（如果当前是子监控，改为"移除父关系"）
  ├─ 测试API
  └─ 删除

  F. 过滤和搜索

  监控列表工具栏：
  原有过滤：
  [业务组▼] [业务功能▼] [状态▼] [搜索框]

  改造后过滤：
  [业务组▼] [业务功能▼] [状态▼]
  [层级过滤: ○全部 ◉仅顶级 ○只看子监控▼]
  [搜索框]

  √ 搜索框支持搜索父监控名称来查看子监控

  G. 统计展示修改

  监控卡片/Dashboard（原有）：
  API监控概览
  ├─ 总数：45
  ├─ 启用：43
  ├─ 健康：40
  ├─ 异常：3
  └─ 平均响应时间：123ms

  改造后：
  API监控概览
  ├─ 总数：45
  │  ├─ 顶级监控：15
  │  └─ 子监控：30
  ├─ 启用：43
  ├─ 健康：40
  ├─ 异常：3
  ├─ 分层监控数：12（含子监控）
  └─ 平均响应时间：123ms

  监控卡片（列表中显示）：
  原有：
  ┌─────────────────────────┐
  │ 📡 api.dding.net        │
  │ 健康 ✅ / 最近响应: 45ms │
  │ 最后检测：2分钟前        │
  └─────────────────────────┘

  改造后：
  ┌─────────────────────────┐
  │ 📡 api.dding.net [父]   │
  │ 聚合状态 ❌ 2/3健康     │  ← 新增
  │ 最近响应: 45ms          │
  │ 最后检测：2分钟前        │
  │ [展开子监控▼]           │  ← 新增
  └─────────────────────────┘

  展开后：
  ├─ 🖥️   192.168.7.6:7001  ✅
  ├─ 🖥️   192.168.7.7:7001  ❌
  └─ 🖥️   192.168.7.8:7001  ✅

  ---
  📱 前端API调用示例

  获取顶级监控列表

  // 原有接口仍可用，新增参数支持
  GET /api/v1/monitor/api-monitors?page=1&pageSize=20

  // 改造后支持：
  GET /api/v1/monitor/api-monitors?
      page=1&
      pageSize=20&
      onlyTopLevel=true&           // 只查询顶级
      includeChildren=true         // 返回中包含子监控

  响应示例：
  {
    "code": 200,
    "data": {
      "list": [
        {
          "id": 1,
          "name": "api.dding.net",
          "url": "https://api.dding.net/system/check",
          "healthStatus": "unhealthy",
          "hierarchyLevel": 0,
          "parentId": null,
          "childCount": 3,           // ← 新增字段
          "healthyChildren": 2,      // ← 新增字段
          "childMonitors": [         // ← 新增字段（当includeChildren=true）
            {
              "id": 2,
              "name": "IP-192.168.7.6",
              "url": "http://192.168.7.6:7001/system/check",
              "healthStatus": "healthy",
              "hierarchyLevel": 1,
              "parentId": 1
            },
            ...
          ]
        }
      ],
      "total": 15,
      "page": 1,
      "pageSize": 20
    }
  }

  获取监控树形结构

  GET /api/v1/monitor/api-monitors/{id}/hierarchy?depth=2

  响应示例：
  {
    "code": 200,
    "data": {
      "id": 1,
      "name": "api.dding.net",
      "url": "https://api.dding.net/system/check",
      "healthStatus": "unhealthy",
      "hierarchyLevel": 0,
      "childCount": 3,
      "healthyChildren": 2,
      "children": [
        {
          "id": 2,
          "name": "IP-192.168.7.6",
          "url": "http://192.168.7.6:7001/system/check",
          "healthStatus": "healthy",
          "hierarchyLevel": 1,
          "childCount": 0,
          "children": []
        },
        ...
      ]
    }
  }

  添加子监控

  POST /api/v1/monitor/api-monitors/{parentId}/add-child
  Content-Type: application/json

  {
    "childMonitorId": 123
  }

  成功响应：
  {
    "code": 200,
    "data": {
      "message": "子监控添加成功",
      "childMonitorId": 123,
      "parentMonitorId": 1
    }
  }

  移除父监控关系

  DELETE /api/v1/monitor/api-monitors/{childId}/parent

  成功响应：
  {
    "code": 200,
    "data": {
      "message": "父监控关系已移除，监控现为独立监控",
      "monitorId": 123
    }
  }

  ---
  ⚠️  重要注意事项
  ┌──────────┬────────────────────────────────────────┐
  │   项目   │                  说明                  │
  ├──────────┼────────────────────────────────────────┤
  │ 向下兼容 │ ✅ 现有监控无需修改，自动作为顶级监控  │
  ├──────────┼────────────────────────────────────────┤
  │ 默认行为 │ 查询列表时默认只显示顶级监控，保持原样 │
  ├──────────┼────────────────────────────────────────┤
  │ 最大层级 │ ⚠️  最多3层（顶级 → 一级子 → 二级子）   │
  ├──────────┼────────────────────────────────────────┤
  │ 循环检查 │ ❌ 禁止 A→B→C→A 这样的循环依赖         │
  ├──────────┼────────────────────────────────────────┤
  │ 删除策略 │ 删除父监控前需选择处理子监控的方式     │
  ├──────────┼────────────────────────────────────────┤
  │ 自动更新 │ 后端每60秒自动更新父监控的聚合状态     │
  ├──────────┼────────────────────────────────────────┤
  │ 独立检测 │ 所有监控独立检测，不受父子关系影响     │
  └──────────┴────────────────────────────────────────┘
  ---
  🔄 前端实现建议

  优先级：
  1. ✅ 高优先级 - 树形列表展示、添加/移除子监控
  2. ✅ 中优先级 - 删除确认对话框、子监控数统计显示
  3. ⭕ 低优先级 - 分层过滤、批量操作

  技术方案：
  - 树形组件：可用 antd TreeSelect 或 Vue Tree
  - 列表展示：嵌套 Table 或手动缩进
  - 模态框：删除确认选择处理方式
  - 性能：使用虚拟滚动处理大量子监控

  改造已完成 ✅，现在支持灵活的分层监控来应对复杂的分布式服务场景！