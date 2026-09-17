<template>
  <!-- 命名空间列表表格 -->
  <div class="table-section">
    <el-table
      :data="tableData"
      v-loading="loading"
      stripe
      style="width: 100%"
      class="namespace-table"
    >
      <el-table-column prop="name" label="命名空间名称" min-width="180">
        <template #default="{ row }">
          <div class="namespace-name-container">
            <img src="@/assets/image/k8s.svg" alt="k8s" class="k8s-icon" />
            <el-button
              type="primary"
              link
              @click="viewNamespaceDetail(row)"
              class="namespace-name-link"
            >
              {{ row.name }}
            </el-button>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="资源统计" min-width="160">
        <template #default="{ row }">
          <div class="resource-stats">
            <div class="resource-item">
              <span class="resource-label pod-label">Pod</span>
              <span class="resource-value pod-value">{{ row.resourceCount.pods }}</span>
            </div>
            <div class="resource-item">
              <span class="resource-label svc-label">Service</span>
              <span class="resource-value svc-value">{{ row.resourceCount.services }}</span>
            </div>
            <div class="resource-item">
              <span class="resource-label secret-label">Secret</span>
              <span class="resource-value secret-value">{{ row.resourceCount.secrets }}</span>
            </div>
            <div class="resource-item">
              <span class="resource-label cm-label">ConfigMap</span>
              <span class="resource-value cm-value">{{ row.resourceCount.configMaps }}</span>
            </div>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="status" label="状态" min-width="100">
        <template #default="{ row }">
          <el-tag
            :type="getStatusTag(row.status)"
            size="small"
            effect="dark"
          >
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="age" label="创建时间" min-width="100">
        <template #default="{ row }">
          <el-tag size="small" type="info">{{ row.age }}</el-tag>
        </template>
      </el-table-column>

      <el-table-column label="标签" min-width="100" align="center">
        <template #default="{ row }">
          <div class="label-container">
            <el-badge :value="getLabelCount(row.labels)" :max="99" class="label-badge">
              <el-button
                type="text"
                size="small"
                circle
                @click="viewNamespaceLabels(row)"
                class="label-icon-button"
              >
                <img src="@/assets/image/标签.svg" alt="标签" width="14" height="14" />
              </el-button>
            </el-badge>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="配额/限制" min-width="120" align="center">
        <template #default="{ row }">
          <div class="quota-limit-info">
            <el-badge :value="row.resourceQuotas.length" :max="99" class="quota-badge">
              <el-tag size="small" type="warning">配额</el-tag>
            </el-badge>
            <el-badge :value="row.limitRanges.length" :max="99" class="limit-badge">
              <el-tag size="small" type="info">限制</el-tag>
            </el-badge>
          </div>
        </template>
      </el-table-column>
      
      <el-table-column label="操作" width="190" fixed="right">
        <template #default="{ row }">
          <div class="operation-buttons">
            <el-button
              type="primary"
              :icon="View"
              size="small"
              v-authority="['k8s:namespace:details']"
              @click="viewNamespaceDetail(row)"
            >详情</el-button>

            <el-dropdown @command="(cmd) => handleMoreCommand(row, cmd)" trigger="click">
              <el-button type="warning" size="small" v-authority="['k8s:namespace:actions']">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="yaml" :icon="Document">yaml编辑</el-dropdown-item>
                  <el-dropdown-item command="quota" :icon="Monitor">配额管理</el-dropdown-item>
                  <el-dropdown-item command="limit" :icon="Setting">资源限制</el-dropdown-item>
                  <el-dropdown-item command="labels" :icon="Edit">修改标签</el-dropdown-item>
                  <el-dropdown-item
                    command="delete"
                    :icon="Delete"
                    divided
                    :disabled="isSystemNamespace(row.name)"
                  >删除空间</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import {
  View,
  Delete,
  Monitor,
  Setting,
  Edit,
  Document,
  ArrowDown
} from '@element-plus/icons-vue'

// Props
const props = defineProps({
  tableData: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits([
  'viewNamespaceDetail',
  'viewNamespaceLabels',
  'manageResourceQuotas',
  'manageLimitRanges',
  'deleteNamespace',
  'editYaml',
  'manageLabels'
])

// 获取状态标签样式
const getStatusTag = (status) => {
  const tagMap = {
    'Active': 'success',
    'Terminating': 'warning',
    'Unknown': 'info'
  }
  return tagMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const textMap = {
    'Active': '活跃',
    'Terminating': '终止中',
    'Unknown': '未知'
  }
  return textMap[status] || '未知'
}

// 获取标签数量（过滤系统标签）
const getLabelCount = (labels) => {
  if (!labels) return 0
  
  const systemLabelPrefixes = [
    'kubernetes.io/',
    'beta.kubernetes.io/',
    'app.kubernetes.io/'
  ]
  
  return Object.keys(labels).filter(key => 
    !systemLabelPrefixes.some(prefix => key.startsWith(prefix))
  ).length
}

// 判断是否为系统命名空间
const isSystemNamespace = (namespaceName) => {
  const systemNamespaces = [
    'default',
    'kube-system', 
    'kube-public',
    'kube-node-lease'
  ]
  return systemNamespaces.includes(namespaceName)
}

// 事件处理函数
const viewNamespaceDetail = (row) => {
  emit('viewNamespaceDetail', row)
}

const viewNamespaceLabels = (row) => {
  emit('viewNamespaceLabels', row)
}

const manageResourceQuotas = (row) => {
  emit('manageResourceQuotas', row)
}

const manageLimitRanges = (row) => {
  emit('manageLimitRanges', row)
}

const deleteNamespace = (row) => {
  emit('deleteNamespace', row)
}

const handleMoreCommand = (row, command) => {
  switch (command) {
    case 'yaml':   emit('editYaml', row); break
    case 'quota':  manageResourceQuotas(row); break
    case 'limit':  manageLimitRanges(row); break
    case 'labels': emit('manageLabels', row); break
    case 'delete': deleteNamespace(row); break
  }
}
</script>

<style scoped>
.table-section {
  margin-top: 20px;
}

.namespace-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.namespace-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.namespace-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #000000 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: none;
}

.namespace-table :deep(.el-table__header th .cell) {
  color: #000000 !important;
  font-weight: 700 !important;
  text-shadow: none;
}

.namespace-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.namespace-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.namespace-name-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.k8s-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.namespace-name-link {
  font-weight: 600;
  color: #667eea;
  text-decoration: none;
  transition: all 0.3s ease;
}

.namespace-name-link:hover {
  color: #764ba2;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.resource-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 4px;
}

.resource-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.resource-label {
  font-size: 11px;
  font-weight: 700;
  min-width: 52px;
  padding: 1px 4px;
  border-radius: 3px;
}

.resource-value {
  font-size: 12px;
  font-weight: 700;
}

/* Pod - 蓝色 */
.pod-label   { color: #409eff; background: rgba(64,158,255,0.1); }
.pod-value   { color: #409eff; }

/* Service - 绿色 */
.svc-label   { color: #67c23a; background: rgba(103,194,58,0.1); }
.svc-value   { color: #67c23a; }

/* Secret - 橙色 */
.secret-label { color: #e6a23c; background: rgba(230,162,60,0.1); }
.secret-value { color: #e6a23c; }

/* ConfigMap - 紫色 */
.cm-label    { color: #764ba2; background: rgba(118,75,162,0.1); }
.cm-value    { color: #764ba2; }

.quota-limit-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.quota-badge, .limit-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

/* 标签容器样式 */
.label-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  padding: 4px 0;
}

.label-badge {
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.label-icon-button {
  background: transparent;
  border: none;
  color: #606266;
  transition: all 0.3s ease;
}

.label-icon-button:hover {
  background: transparent;
  color: #409eff;
  transform: scale(1.1);
}

.operation-buttons {
  display: flex;
  gap: 6px;
  align-items: center;
  flex-wrap: nowrap;
}

.operation-buttons .el-button {
  transition: all 0.3s ease;
  pointer-events: auto;
  position: relative;
  z-index: 1;
}

.operation-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.operation-buttons .el-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .operation-buttons {
    gap: 4px;
  }
  
  .operation-buttons .el-button {
    margin: 1px;
  }
}

@media (max-width: 768px) {
  .operation-buttons {
    flex-direction: column;
    gap: 4px;
  }
  
  .namespace-table :deep(.el-table__row:hover) {
    transform: none;
  }
}
</style>