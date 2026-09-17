# UI Design Patterns - DoDevOps Web DB
# Source: k8s-clusters.vue + CmdbGroup.vue + CmdbHostTable.vue
# Purpose: AI Agent reference for generating consistent UI code

---

## 1. 颜色系统 (Color System)

### 1.1 品牌主色（Brand Colors）
```css
/* 主渐变 - 用于标题、对话框头部、表格头部、重要强调 */
--gradient-primary-135: linear-gradient(135deg, #667eea, #764ba2);
--gradient-primary-45:  linear-gradient(45deg,  #667eea, #764ba2);
--gradient-primary-90:  linear-gradient(90deg,  #667eea, #764ba2);
--gradient-primary-145: linear-gradient(145deg, rgba(255,255,255,0.95), rgba(255,255,255,0.85)); /* 卡片 */

/* 主色（单色用法）*/
--color-primary: #667eea;
--color-primary-alt: #764ba2;
--color-primary-rgb: 102, 126, 234;   /* 用于 rgba() */

/* 主色透明度变体 */
--color-primary-10:  rgba(102, 126, 234, 0.10);
--color-primary-08:  rgba(102, 126, 234, 0.08);
--color-primary-05:  rgba(102, 126, 234, 0.05);
--color-primary-03:  rgba(102, 126, 234, 0.03);
--color-primary-border: rgba(102, 126, 234, 0.15);
--color-primary-border-md: rgba(102, 126, 234, 0.2);
--color-primary-border-strong: rgba(102, 126, 234, 0.3);
```

### 1.2 语义色（Semantic Colors）
```css
/* 成功 / 在线 / 详情 */
--color-success: #67C23A;
--color-success-light: #85CE61;
--color-success-bg: rgba(103, 194, 58, 0.15);

/* 警告 / 文件管理 */
--color-warning: #E6A23C;
--color-warning-light: #EBB563;

/* 危险 / 删除 / 离线 */
--color-danger: #F56C6C;
--color-danger-light: #F78989;
--color-danger-bg: rgba(245, 108, 108, 0.1);

/* 信息 / 蓝色 / 编辑 */
--color-info: #409EFF;
--color-info-light: #66B1FF;

/* 监控（深蓝紫） */
--color-monitor-start: #1e3a8a;
--color-monitor-end: #3b82f6;
```

### 1.3 中性色（Neutral Colors）
```css
/* 文字 */
--text-primary:   #2c3e50;   /* 主文字、表格内容 */
--text-secondary: #606266;   /* 次要文字、标签 */
--text-hint:      #909399;   /* 提示文字、空状态 */
--text-disabled:  #c0c4cc;

/* 背景 */
--bg-page:    #eef1f6;                         /* 页面底色 */
--bg-card:    rgba(255, 255, 255, 0.95);       /* 普通卡片 */
--bg-section: #fafafa;                         /* 分区块底色 */
--bg-input:   rgba(255, 255, 255, 0.8);        /* 输入框 */
--bg-preview: #f5f7fa;                         /* 预览/标签列表底色 */

/* 边框 */
--border-default: #e4e7ed;
--border-light:   #dcdfe6;
```

### 1.4 资源使用率颜色（Usage Color Coding）
```javascript
// CPU/内存/磁盘使用率颜色逻辑
const getUsageTextColor = (usage) => {
  if (!usage) return '#909399'    // 无数据 - 灰色
  if (usage >= 85) return '#F56C6C'  // 危险 - 红色
  if (usage >= 50) return '#E6A23C'  // 警告 - 黄色
  return '#67C23A'                    // 正常 - 绿色
}
```

---

## 2. 布局结构 (Layout Patterns)

### 2.1 标准页面容器（Standard Page Container）
```vue
<template>
  <div class="page-management">
    <el-card shadow="hover" class="main-card">
      <template #header>
        <div class="card-header">
          <span class="title">页面标题</span>
          <div class="header-actions">
            <el-button type="primary" :icon="Plus" size="small">操作</el-button>
          </div>
        </div>
      </template>
      <XxxTable ref="tableRef" :refresh-trigger="refreshTrigger" />
    </el-card>
    <!-- Dialogs 放卡片外 -->
  </div>
</template>

<style scoped>
.page-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}
.main-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-actions { display: flex; gap: 12px; }
</style>
```

### 2.2 侧边栏 + 主内容布局（Sidebar + Content Layout）
> 来源：CmdbGroup.vue  ——  左树 + 右表的经典两栏

```vue
<template>
  <div class="host-page" style="display: flex; align-items: flex-start;">
    <!-- 左侧树形侧边栏 -->
    <div class="sidebar-section">
      <div class="sidebar-card">
        <!-- 科技感标题 -->
        <div class="card-header">
          <div class="title-wrapper">
            <div class="title-icon">
              <div class="pulse-ring"></div>
              <el-icon class="main-icon"><DataBoard /></el-icon>
            </div>
            <div class="title-content">
              <h3 class="title">分组标题</h3>
              <span class="subtitle">Asset Groups</span>
            </div>
          </div>
          <div class="stats-indicator">
            <span class="stats-text">{{ count }}</span>
            <span class="stats-label">Groups</span>
          </div>
        </div>
        <!-- 树形内容 -->
        <el-tree class="tech-tree" ... />
      </div>
    </div>
    <!-- 右侧主内容 -->
    <div class="main-content" style="flex: 1;">
      <el-card shadow="hover" class="main-card">...</el-card>
    </div>
  </div>
</template>

<style scoped>
.sidebar-section {
  width: 280px;
  margin-right: 20px;
  flex-shrink: 0;
}
.sidebar-card {
  /* 科技感增强版卡片 */
  background: linear-gradient(145deg, rgba(255,255,255,0.95), rgba(255,255,255,0.85));
  backdrop-filter: blur(20px);
  border-radius: 20px;
  box-shadow:
    0 20px 40px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
  overflow: hidden;
  position: relative;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
/* 顶部发光装饰线 */
.sidebar-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.5), transparent);
  z-index: 1;
}
.sidebar-card:hover {
  transform: translateY(-2px);
  box-shadow:
    0 25px 50px rgba(0, 0, 0, 0.15),
    0 0 0 1px rgba(255, 255, 255, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}
</style>
```

### 2.3 双列表单行（Two-Column Form Row）
```vue
<div class="form-row">
  <el-form-item label="字段A" required class="form-item-half">
    <el-input v-model="form.a" size="small" />
  </el-form-item>
  <el-form-item label="字段B" class="form-item-half">
    <el-select v-model="form.b" size="small" style="width:100%" />
  </el-form-item>
</div>

<style scoped>
.form-row       { display: flex; gap: 16px; margin-bottom: 12px; }
.form-item-half { flex: 1; margin-bottom: 0 !important; }
.form-item-full { width: 100%; margin-bottom: 0 !important; }
.compact-form :deep(.el-form-item) { margin-bottom: 12px; }
</style>
```

### 2.4 分组选择块（Group Selection Block）
> 来源：k8s-clusters.vue —— 节点选择区（Master/Worker/ETCD）

```vue
<div class="nodes-selection">
  <div class="node-section" v-for="group in groups" :key="group.type">
    <div class="node-header">
      <el-button :type="group.btnType" size="small" @click="openSelector(group.type)">
        {{ group.label }}
      </el-button>
      <span class="node-count">{{ group.selected.length }}台</span>
    </div>
    <div class="selected-tags">
      <el-tag
        v-for="item in group.selected" :key="item.id"
        closable @close="removeItem(item.id, group.type)"
        :type="group.tagType" size="small"
      >{{ item.name }}</el-tag>
    </div>
  </div>
</div>

<style scoped>
.nodes-selection { display: flex; gap: 16px; width: 100%; }
.node-section    { flex: 1; border: 1px solid #e4e7ed; border-radius: 6px; padding: 8px; background: #fafafa; }
.node-header     { display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px; }
.node-count      { font-size: 11px; color: #606266; font-weight: 500; }
.selected-tags   { display: flex; flex-wrap: wrap; gap: 3px; margin-top: 4px; }
.selected-tags .el-tag { font-size: 10px; height: 20px; max-width: 80px; }
</style>
```

### 2.5 展开配置块（Expandable Config Block）
```vue
<el-form-item label="高级配置">
  <el-switch v-model="form.useAdvanced" active-text="启用" inactive-text="关闭" size="small" />
  <div v-if="form.useAdvanced" class="advanced-config">
    <!-- 展开内容 -->
  </div>
</el-form-item>

<style scoped>
.advanced-config {
  margin-top: 12px;
  padding: 12px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}
</style>
```

---

## 3. 卡片设计 (Card Design)

### 3.1 标准卡片（Standard Card）
```css
/* 用于主要内容区域 */
.main-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}
```

### 3.2 科技感增强卡片（Tech-Style Card）
```css
/* 用于侧边栏、重点模块，有更强的视觉层次 */
.tech-card {
  background: linear-gradient(145deg, rgba(255,255,255,0.95), rgba(255,255,255,0.85));
  backdrop-filter: blur(20px);
  border-radius: 20px;
  box-shadow:
    0 20px 40px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
}
/* 顶部发光线 */
.tech-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0; height: 1px;
  background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.5), transparent);
}
/* hover 效果 */
.tech-card:hover {
  transform: translateY(-2px);
  box-shadow:
    0 25px 50px rgba(0, 0, 0, 0.15),
    0 0 0 1px rgba(255, 255, 255, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}
```

### 3.3 统计信息徽章（Stats Indicator）
```vue
<div class="stats-indicator">
  <span class="stats-text">{{ count }}</span>
  <span class="stats-label">Items</span>
</div>

<style scoped>
.stats-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: rgba(102, 126, 234, 0.05);
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.1);
  transition: all 0.3s ease;
}
.stats-indicator:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}
.stats-text  { font-size: 16px; font-weight: 700; color: #667eea; line-height: 1; }
.stats-label { font-size: 10px; color: rgba(102, 126, 234, 0.6); text-transform: uppercase; letter-spacing: 0.5px; }
</style>
```

---

## 4. 标题系统 (Typography)

### 4.1 渐变标题（Gradient Title）
```css
/* 页面/卡片主标题 */
.title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

/* 侧边栏标题（稍小）*/
.sidebar-title {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
```

### 4.2 英文副标题（Subtitle）
> 来源：CmdbGroup.vue —— 标题下方的小字英文说明

```css
.subtitle {
  font-size: 11px;
  color: rgba(102, 126, 234, 0.6);
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 500;
  display: block;
}
```

### 4.3 节点图标标题（Icon + Title Pattern）
```vue
<div class="title-wrapper">
  <div class="title-icon">
    <div class="pulse-ring"></div>
    <el-icon class="main-icon"><DataBoard /></el-icon>
  </div>
  <div class="title-content">
    <h3 class="title">中文标题</h3>
    <span class="subtitle">English Subtitle</span>
  </div>
</div>

<style scoped>
.title-wrapper { display: flex; align-items: center; gap: 12px; }
.title-icon    { position: relative; width: 40px; height: 40px; display: flex; align-items: center; justify-content: center; }
.main-icon     { font-size: 20px; color: #667eea; z-index: 2; background: rgba(102,126,234,0.1); padding: 8px; border-radius: 8px; }

/* 脉冲光环动画 */
.pulse-ring {
  position: absolute;
  width: 35px; height: 35px;
  border: 2px solid rgba(102, 126, 234, 0.3);
  border-radius: 50%;
  animation: pulse 2s infinite;
}
@keyframes pulse {
  0%   { transform: scale(1);   opacity: 1; }
  50%  { transform: scale(1.2); opacity: 0.7; }
  100% { transform: scale(1.4); opacity: 0; }
}
</style>
```

---

## 5. 表格设计 (Table Design)

### 5.1 表格基础样式
```css
/* 表格容器 */
.data-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

/* 行 hover 动效（重要：行悬停有上浮+阴影） */
.data-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}
.data-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
```

### 5.2 渐变表头（Gradient Table Header）
> 来源：CmdbHostTable.vue —— 最高等级表格视觉

```css
/* 渐变表头版（推荐用于重要表格）*/
.data-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}
.data-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}
.data-table :deep(.el-table__header th .cell) {
  color: #2c3e50 !important;
  font-weight: 700 !important;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}
```

### 5.3 轻量表头（Light Table Header）
```css
/* 普通版，用于次要表格 */
.data-table :deep(.el-table__header th) {
  background: rgba(102, 126, 234, 0.1);
  color: #2c3e50;
  font-weight: 600;
}
```

### 5.4 操作列按钮组（Table Operation Buttons）
> 来源：CmdbHostTable.vue —— 每个操作各有独立渐变色

```vue
<el-table-column label="操作" width="220" fixed="right" align="center">
  <template v-slot="scope">
    <div class="table-operation">
      <el-button-group class="operation-btn-group">
        <el-tooltip content="详情" placement="top">
          <el-button type="success" icon="Document" size="small"
            class="operation-btn operation-btn-detail" @click="handleDetail(scope.row)" />
        </el-tooltip>
        <el-tooltip content="编辑" placement="top">
          <el-button type="primary" icon="Edit" size="small"
            class="operation-btn operation-btn-edit"
            v-authority="['module:resource:edit']" @click="handleEdit(scope.row)" />
        </el-tooltip>
        <el-tooltip content="删除" placement="top">
          <el-button type="danger" icon="Delete" size="small"
            class="operation-btn operation-btn-delete"
            v-authority="['module:resource:delete']" @click="handleDelete(scope.row)" />
        </el-tooltip>
      </el-button-group>
    </div>
  </template>
</el-table-column>

<style scoped>
.table-operation {
  display: flex; justify-content: center; align-items: center;
  gap: 4px; padding: 4px 0;
}
.operation-btn-group { display: flex; gap: 4px; flex-wrap: nowrap; }

/* 统一按钮基础 */
.operation-btn {
  position: relative;
  width: 30px; height: 30px; padding: 0;
  border-radius: 6px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1.5px solid;
  overflow: hidden;
}
/* 水波纹点击效果 */
.operation-btn::before {
  content: '';
  position: absolute; top: 50%; left: 50%;
  width: 0; height: 0; border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  transform: translate(-50%, -50%);
  transition: width 0.4s ease, height 0.4s ease;
}
.operation-btn:hover::before { width: 80px; height: 80px; }

/* 详情 - 绿色 */
.operation-btn-detail {
  background: linear-gradient(135deg, #67C23A, #85CE61);
  border-color: #67C23A; color: white;
  box-shadow: 0 2px 6px rgba(103, 194, 58, 0.25);
}
.operation-btn-detail:hover {
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 4px 12px rgba(103, 194, 58, 0.4);
}

/* 编辑 - 蓝色 */
.operation-btn-edit {
  background: linear-gradient(135deg, #409EFF, #66B1FF);
  border-color: #409EFF; color: white;
  box-shadow: 0 2px 6px rgba(64, 158, 255, 0.25);
}
.operation-btn-edit:hover {
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
}

/* 文件/下载 - 橙色 */
.operation-btn-file {
  background: linear-gradient(135deg, #E6A23C, #EBB563);
  border-color: #E6A23C; color: white;
  box-shadow: 0 2px 6px rgba(230, 162, 60, 0.25);
}
.operation-btn-file:hover {
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 4px 12px rgba(230, 162, 60, 0.4);
}

/* 删除 - 红色 */
.operation-btn-delete {
  background: linear-gradient(135deg, #F56C6C, #F78989);
  border-color: #F56C6C; color: white;
  box-shadow: 0 2px 6px rgba(245, 108, 108, 0.25);
}
.operation-btn-delete:hover {
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 4px 12px rgba(245, 108, 108, 0.4);
}

/* 监控 - 深蓝紫 */
.operation-btn-monitor {
  background: linear-gradient(135deg, #1e3a8a, #3b82f6);
  border-color: #1e3a8a; color: white;
  box-shadow: 0 2px 6px rgba(30, 58, 138, 0.3);
}
.operation-btn-monitor:hover {
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.45);
}

/* 禁用状态 */
.operation-btn:disabled {
  opacity: 0.5; cursor: not-allowed;
  transform: none !important;
}
</style>
```

### 5.5 自定义滚动条（Custom Scrollbar）
```css
/* 细滚动条（侧边栏树形区域）*/
.tree-container::-webkit-scrollbar { width: 4px; }
.tree-container::-webkit-scrollbar-track {
  background: rgba(102, 126, 234, 0.05);
  border-radius: 2px;
}
.tree-container::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, rgba(102,126,234,0.3), rgba(118,75,162,0.3));
  border-radius: 2px;
}
.tree-container::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, rgba(102,126,234,0.5), rgba(118,75,162,0.5));
}
```

---

## 6. 按钮设计 (Button Design)

### 6.1 标准 el-button（在卡片头部/表单）
```css
.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}
.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
```

### 6.2 小型图标圆形按钮（Icon-Only Circle Buttons）
> 来源：CmdbGroup.vue —— 刷新、创建、折叠的图标按钮

```vue
<!-- 刷新按钮 -->
<div class="icon-btn refresh-btn" @click="handleRefresh" title="刷新">
  <el-icon class="refresh-icon"><Refresh /></el-icon>
</div>
<!-- 创建按钮 -->
<div class="icon-btn create-btn" @click="handleCreate" title="创建">
  <el-icon class="create-icon"><Plus /></el-icon>
</div>

<style scoped>
/* 通用图标按钮基础 */
.icon-btn {
  width: 20px; height: 20px;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; border-radius: 50%;
  transition: all 0.3s ease; border: 1px solid;
}
.icon-btn:hover  { transform: scale(1.1); }
.icon-btn:active { transform: scale(0.95); }

/* 刷新 - 蓝紫色 */
.refresh-btn { background: rgba(102,126,234,0.15); border-color: rgba(102,126,234,0.2); }
.refresh-btn:hover { background: rgba(102,126,234,0.25); border-color: rgba(102,126,234,0.35); }
.refresh-icon { font-size: 16px; color: #667eea; }
.refresh-btn:hover .refresh-icon { transform: rotate(180deg); color: #764ba2; }

/* 创建 - 绿色 */
.create-btn { background: rgba(103,194,58,0.15); border-color: rgba(103,194,58,0.2); }
.create-btn:hover { background: rgba(103,194,58,0.25); border-color: rgba(103,194,58,0.35); }
.create-icon { font-size: 16px; color: #67C23A; }
.create-btn:hover .create-icon { color: #85CE61; transform: rotate(90deg); }
</style>
```

### 6.3 胶囊切换按钮（Pill Toggle Buttons）
> 来源：k8s-clusters.vue —— 网络插件、可选组件选择

```css
.toggle-btn-group .el-button {
  border-radius: 12px;
  font-size: 10px;
  padding: 3px 6px;
  height: 22px;
  line-height: 1;
  white-space: nowrap;
}
```

### 6.4 带数量角标的按钮（Badge Button）
> 来源：CmdbHostTable.vue —— 标签数量按钮

```vue
<el-button type="primary" size="small" circle @click="showTags(row)">
  <el-icon><PriceTag /></el-icon>
  <span class="badge-count">{{ count }}</span>
</el-button>

<style scoped>
.badge-count {
  position: absolute;
  top: -6px; right: -6px;
  background: #f56c6c; color: white;
  font-size: 9px; font-weight: 600;
  padding: 1px 4px; border-radius: 8px;
  line-height: 1.2; min-width: 14px;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
}
</style>
```

---

## 7. 输入框与表单 (Form & Input)

### 7.1 科技感输入框（Tech Input）
```css
/* 标准输入框 */
.el-input :deep(.el-input__wrapper),
.el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}
.el-input :deep(.el-input__wrapper):hover { border-color: #c0c4cc; }
.el-input :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

/* 搜索框（更轻量，带内嵌阴影）*/
.tech-input :deep(.el-input__wrapper) {
  background: rgba(102, 126, 234, 0.03);
  border: 1px solid rgba(102, 126, 234, 0.15);
  border-radius: 12px;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.02);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.tech-input :deep(.el-input__wrapper):hover {
  border-color: rgba(102, 126, 234, 0.3);
  background: rgba(102, 126, 234, 0.05);
}
.tech-input :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}
```

### 7.2 文本域（Textarea）
```css
.el-textarea :deep(.el-textarea__inner) {
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.2);
  transition: all 0.3s ease;
}
.el-textarea :deep(.el-textarea__inner):focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
}
```

### 7.3 搜索图标颜色联动
```css
.search-icon { color: rgba(102, 126, 234, 0.6); transition: all 0.3s ease; }
.tech-input:focus-within .search-icon { color: #667eea; }
```

---

## 8. 对话框设计 (Dialog Design)

### 8.1 渐变头部对话框（Gradient Header Dialog）
```css
.xxx-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}
.xxx-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}
.xxx-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}
.xxx-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

/* 复杂表单对话框（紧凑版） */
.compact-dialog :deep(.el-dialog__body) { padding: 15px 20px; }
.compact-form :deep(.el-form-item)      { margin-bottom: 12px; }
```

### 8.2 对话框底部按钮
```vue
<template #footer>
  <span class="dialog-footer">
    <el-button @click="visible = false">取消</el-button>
    <el-button type="success" @click="handleSubmit">确认</el-button>
  </span>
</template>

<style scoped>
.dialog-footer { display: flex; gap: 12px; justify-content: flex-end; }
</style>
```

### 8.3 右键菜单（Context Menu）
> 来源：CmdbGroup.vue

```vue
<div
  v-if="contextMenuVisible"
  class="context-menu"
  :style="{ left: pos.x + 'px', top: pos.y + 'px' }"
>
  <div class="context-menu-item" @click="handleAction">
    <el-icon><Edit /></el-icon>
    <span>编辑</span>
  </div>
  <div class="context-menu-item danger" @click="handleDelete">
    <el-icon><Delete /></el-icon>
    <span>删除</span>
  </div>
</div>

<style scoped>
.context-menu {
  position: absolute;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 4px; z-index: 9999; min-width: 130px;
}
.context-menu-item {
  display: flex; align-items: center; gap: 8px;
  padding: 8px 12px; cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s ease;
  font-size: 14px; color: #2c3e50;
}
.context-menu-item:hover       { background: rgba(102, 126, 234, 0.1); color: #667eea; }
.context-menu-item.danger      { color: #f56c6c; }
.context-menu-item.danger:hover { background: rgba(245, 108, 108, 0.1); color: #f56c6c; }
</style>
```

---

## 9. 树形组件 (Tree Component)

### 9.1 科技感树形（Tech Tree）
```css
/* 树形基础 */
.tech-tree { border: none; background: transparent; }
.tech-tree :deep(.el-tree-node__content) {
  height: auto; padding: 0;
  background: transparent; border-radius: 0;
}

/* 展开图标 */
.tech-tree :deep(.el-tree-node__expand-icon) {
  color: #667eea;
  font-size: 12px;
  transition: transform 0.3s ease;
}
.tech-tree :deep(.el-tree-node__expand-icon.expanded) {
  transform: rotate(90deg);
}
.tech-tree :deep(.el-tree-node__expand-icon.is-leaf) { display: none; }

/* 树节点容器 */
.tree-node {
  margin: 0 15px;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  animation: dataLoad 0.4s ease-out;  /* 入场动画 */
}

/* 节点 hover 渐变蒙层 */
.hover-effect {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: linear-gradient(135deg, rgba(102,126,234,0.08), rgba(118,75,162,0.08));
  border-radius: 12px;
  opacity: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.tree-node.parent-node:hover > .hover-effect { opacity: 1; }

/* 选中状态 */
.tech-tree :deep(.el-tree-node.is-current > .el-tree-node__content) .tree-node {
  background: linear-gradient(135deg, rgba(102,126,234,0.08), rgba(118,75,162,0.08));
}

/* 子节点圆点指示器 */
.dot-indicator {
  position: absolute;
  width: 6px; height: 6px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 50%;
  top: 50%; left: -8px;
  transform: translateY(-50%);
}
```

### 9.2 节点入场动画（Staggered Entrance）
```css
@keyframes dataLoad {
  0%   { opacity: 0; transform: translateY(20px); }
  100% { opacity: 1; transform: translateY(0); }
}
.tree-node:nth-child(1) { animation-delay: 0.1s; }
.tree-node:nth-child(2) { animation-delay: 0.2s; }
.tree-node:nth-child(3) { animation-delay: 0.3s; }
.tree-node:nth-child(4) { animation-delay: 0.4s; }
.tree-node:nth-child(5) { animation-delay: 0.5s; }
```

---

## 10. 悬停复制模式 (Hover-to-Copy Pattern)

> 来源：CmdbHostTable.vue —— 鼠标悬停才显示复制图标

```vue
<div class="copyable-cell" @mouseenter="showCopy" @mouseleave="hideCopy">
  <span>{{ content }}</span>
  <el-icon class="copy-icon" @click="copyToClipboard(content)">
    <DocumentCopy />
  </el-icon>
</div>

<style scoped>
.copy-icon {
  opacity: 0;
  transition: all 0.3s ease;
  font-size: 14px;
  padding: 2px; border-radius: 4px;
  cursor: pointer; color: #409EFF;
}
.copy-icon:hover {
  background-color: rgba(64, 158, 255, 0.1);
  transform: scale(1.1);
}
.copyable-cell:hover .copy-icon { opacity: 1; }
.copy-icon.copied { color: #67C23A !important; transform: scale(1.2); }
</style>
```

---

## 11. 状态标识 (Status Indicators)

### 11.1 在线/离线状态
```vue
<div class="status-cell">
  <img :src="row.isAlive ? onlineSvg : offlineSvg" style="width:16px;height:16px;margin-right:6px" />
  <el-tag :type="row.isAlive ? 'success' : 'danger'" size="small">
    {{ row.isAlive ? '在线' : '离线' }}
  </el-tag>
</div>
```

### 11.2 脉冲呼吸灯（Online Pulse）
```css
.status-dot.online {
  background: #67C23A;
  box-shadow: 0 0 6px rgba(103, 194, 58, 0.5);
}
.status-dot.online::before {
  content: '';
  position: absolute;
  width: 8px; height: 8px; border-radius: 50%;
  background: #67C23A;
  animation: ping 2s infinite;
}
@keyframes ping {
  75%, 100% { transform: scale(2); opacity: 0; }
}
```

---

## 12. 动画规范 (Animation Standards)

```css
/* 标准过渡函数 */
--ease-standard: cubic-bezier(0.4, 0, 0.2, 1);  /* 大多数动画 */
--ease-bounce:   cubic-bezier(0.4, 0, 0.2, 1);   /* 弹性效果 */

/* 常用 transition 时长 */
--duration-fast:   0.3s;   /* 按钮、输入框、颜色变化 */
--duration-normal: 0.3s;   /* 卡片、面板 */
--duration-slow:   0.4s;   /* 入场动画 */

/* 按钮 hover 上浮 */
.hover-lift:hover  { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }

/* 操作按钮 hover */
.op-btn:hover      { transform: translateY(-2px) scale(1.05); }
.op-btn:active     { transform: translateY(-1px) scale(1.02); }

/* 图标旋转刷新 */
.refresh-icon { transition: all 0.3s ease; }
.refresh-btn:hover .refresh-icon { transform: rotate(180deg); }

/* 创建图标旋转 */
.create-btn:hover .create-icon { transform: rotate(90deg); }
```

---

## 13. 响应式断点 (Responsive Breakpoints)

```css
/* 平板端：按钮组换行 */
@media (max-width: 1000px) {
  .plugin-single-line { flex-wrap: wrap; gap: 12px; }
}

/* 移动端：布局变单列 */
@media (max-width: 768px) {
  .form-row           { flex-direction: column; gap: 8px; }
  .nodes-selection    { flex-direction: column; gap: 8px; }
  .compact-dialog     { width: 95% !important; }
  .plugin-single-line { flex-direction: column; align-items: flex-start; gap: 8px; }
}

/* 小屏：压缩尺寸 */
@media (max-width: 480px) {
  .toggle-btn-group .el-button { font-size: 10px; padding: 3px 6px; height: 22px; }
}
```

---

## 14. 圆角规范 (Border Radius Scale)

| 值 | 用途 |
|---|---|
| `4px` | 右键菜单项、微小元素 |
| `6px` | 操作按钮、分组选择块 |
| `8px` | 输入框、标签 Tag、普通按钮、图标背景 |
| `12px` | 搜索框、表格、胶囊按钮 |
| `16px` | 主卡片 `el-card`、对话框 |
| `20px` | 科技感侧边栏卡片 |
| `50%` | 圆形图标按钮（刷新、创建） |

---

## 15. 间距规范 (Spacing Scale)

```css
/* 组件间距 */
--gap-xs:   4px;   /* tag 之间、按钮图标内部 */
--gap-sm:   6px;   /* 小型按钮组 */
--gap-md:   8px;   /* 树节点、搜索容器内部 */
--gap-lg:   12px;  /* 标题图标、按钮组、右键菜单项 */
--gap-xl:   16px;  /* 表单行、分组选择块 */
--gap-2xl:  20px;  /* 侧边栏间距 */

/* 内边距 */
--padding-card:    20px 24px;  /* 卡片头部 */
--padding-dialog:  24px;       /* 对话框内容 */
--padding-compact: 15px 20px;  /* 紧凑表单对话框 */
--padding-section: 15px 20px;  /* 搜索区域 */
--padding-node:    8px;        /* 节点选择块 */
```

---

## 16. 反模式（Anti-Patterns）

```
❌ 硬编码颜色（在设计系统外写 #hex 值）
❌ 不同页面混用不同圆角值（统一用圆角规范）
❌ 使用 !important 覆盖样式（修正 CSS 特异性）
❌ inline style 写样式（统一用 scoped class）
❌ 缺少 hover/focus 交互反馈
❌ 表格不加行 hover 效果
❌ 操作按钮不用 el-button-group + el-tooltip
❌ 重要表格不加渐变表头
❌ 对话框缺少渐变头部（配置类对话框可不加，创建/注册等重要操作必须加）
❌ 不加 v-authority 指令到操作按钮
❌ 对话框宽度乱写（参考：简单表单 400-500px，复杂表单 700px）
```
