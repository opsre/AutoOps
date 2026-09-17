<template>
  <!-- 节点列表表格 -->
  <div class="table-section">
    <el-table
      :data="tableData"
      v-loading="loading"
      stripe
      style="width: 100%"
      class="nodes-table"
    >
      <el-table-column prop="nodeName" label="节点名称" min-width="180">
        <template #default="{ row }">
          <div class="node-name-container">
            <img src="@/assets/image/k8s.svg" alt="k8s" class="k8s-icon" />
            <el-button
              type="primary"
              link
              @click="navigateToNodeDetail(row)"
              class="node-name-link"
            >
              {{ row.nodeName }}
            </el-button>
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
      
      <el-table-column prop="role" label="角色" min-width="100">
        <template #default="{ row }">
          <el-tag
            :type="getRoleTag(row.role)"
            size="small"
            effect="plain"
          >
            {{ getRoleText(row.role) }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column prop="version" label="Kubelet版本" min-width="120">
        <template #default="{ row }">
          <span class="version-text">{{ row.version }}</span>
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
                @click="viewNodeLabels(row)"
                class="label-icon-button"
              >
                <img src="@/assets/image/标签.svg" alt="标签" width="14" height="14" />
              </el-button>
            </el-badge>
          </div>
        </template>
      </el-table-column>
      
      <el-table-column prop="age" label="运行时间" min-width="100">
        <template #default="{ row }">
          <el-tag size="small" type="info">{{ row.age }}</el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="CPU/内存" min-width="140">
        <template #default="{ row }">
          <div class="resource-info">
            <div class="resource-item">
              <span class="resource-label">CPU:</span>
              <span class="resource-value">{{ row.cpu }}</span>
            </div>
            <div class="resource-item">
              <span class="resource-label">Memory:</span>
              <span class="resource-value">{{ row.memory }}</span>
            </div>
          </div>
        </template>
      </el-table-column>
      
      <el-table-column prop="pods" label="Pod数量" min-width="100">
        <template #default="{ row }">
          <el-tag size="small" type="warning">{{ row.pods }}</el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="调度状态" min-width="100">
        <template #default="{ row }">
          <el-tag
            :type="row.schedulable ? 'success' : 'danger'"
            size="small"
            :icon="row.schedulable ? CircleCheck : Lock"
          >
            {{ row.schedulable ? '可调度' : '已封锁' }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="污点数量" min-width="90" align="center">
        <template #default="{ row }">
          <div class="taint-container">
            <el-badge :value="row.taints.length" :max="99" class="taint-badge">
              <el-icon size="18" class="taint-icon"><Warning /></el-icon>
            </el-badge>
          </div>
        </template>
      </el-table-column>
      
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <div class="operation-buttons">
            <el-button
              type="primary"
              :icon="View"
              size="small"
              v-authority="['k8s:node:details']"
              @click="navigateToNodeDetail(row)"
            >详情</el-button>

            <el-dropdown @command="(cmd) => handleMoreCommand(row, cmd)" trigger="click">
              <el-button type="warning" size="small" v-authority="['k8s:node:actions']" >
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="yaml" :icon="Document">yaml编辑</el-dropdown-item>
                  <el-dropdown-item command="drain" :icon="SwitchButton" :disabled="!row.schedulable">节点排水</el-dropdown-item>
                  <el-dropdown-item command="labels" :icon="Edit">标签管理</el-dropdown-item>
                  <el-dropdown-item command="taints" :icon="Warning">污点管理</el-dropdown-item>
                  <el-dropdown-item command="schedule" :icon="row.schedulable ? Lock : Unlock" divided>
                    {{ row.schedulable ? '封锁节点' : '解封节点' }}
                  </el-dropdown-item>
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
import { useRouter } from 'vue-router'
import {
  CircleCheck,
  Lock,
  Unlock,
  Warning,
  Edit,
  SwitchButton,
  View,
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
  },
  selectedClusterId: {
    type: [String, Number],
    required: true
  }
})

// Emits
const emit = defineEmits([
  'manageTaints',
  'manageLabels',
  'viewNodeLabels',
  'toggleCordon',
  'drainNode',
  'editYaml'
])

const router = useRouter()

// 导航到节点详情页面
const navigateToNodeDetail = (row) => {
  router.push(`/k8s/cluster/${props.selectedClusterId}/node/${row.nodeName}`)
}

// 获取状态标签样式
const getStatusTag = (status) => {
  switch (status) {
    case 'Ready':
      return 'success'
    case 'NotReady':
      return 'danger'
    case 'SchedulingDisabled':
      return 'warning'
    default:
      return 'info'
  }
}

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 'Ready':
      return '就绪'
    case 'NotReady':
      return '未就绪'
    case 'SchedulingDisabled':
      return '调度禁用'
    default:
      return '未知'
  }
}

// 获取角色标签样式
const getRoleTag = (role) => {
  switch (role) {
    case 'master':
    case 'control-plane':
      return 'danger'
    case 'worker':
      return 'success'
    default:
      return 'info'
  }
}

// 获取角色文本
const getRoleText = (role) => {
  switch (role) {
    case 'master':
    case 'control-plane':
      return 'master'
    case 'worker':
      return 'worker'
    default:
      return '未知'
  }
}

// 获取标签数量
const getLabelCount = (labels) => {
  return labels ? Object.keys(labels).length : 0
}

// 事件处理函数
const manageTaints = (row) => {
  emit('manageTaints', row)
}

const manageLabels = (row) => {
  emit('manageLabels', row)
}

const viewNodeLabels = (row) => {
  emit('viewNodeLabels', row)
}

const toggleCordon = (row) => {
  emit('toggleCordon', row)
}

const drainNode = (row) => {
  emit('drainNode', row)
}

const handleMoreCommand = (row, command) => {
  switch (command) {
    case 'yaml':     emit('editYaml', row); break
    case 'drain':    drainNode(row); break
    case 'labels':   manageLabels(row); break
    case 'taints':   manageTaints(row); break
    case 'schedule': toggleCordon(row); break
  }
}
</script>

<style scoped>
.table-section {
  margin-top: 20px;
}

.nodes-table {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  backdrop-filter: blur(10px);
  overflow: hidden;
}

.nodes-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #000000 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: none;
}

.nodes-table :deep(.el-table__header th .cell) {
  color: #000000 !important;
  font-weight: 700 !important;
}

.node-name-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.k8s-icon {
  width: 18px;
  height: 18px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.node-name-link {
  font-weight: 600;
  text-decoration: none;
}

.node-name-link:hover {
  color: #409EFF;
}

.version-text {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #606266;
}

.label-container,
.taint-container {
  display: flex;
  justify-content: center;
  align-items: center;
}

.label-badge,
.taint-badge {
  cursor: pointer;
}

.label-icon-button {
  padding: 8px;
  border-radius: 50%;
  transition: all 0.3s ease;
}

.label-icon-button:hover {
  background: rgba(64, 158, 255, 0.1);
  transform: scale(1.1);
}

.taint-icon {
  color: #E6A23C;
  cursor: pointer;
  transition: all 0.3s ease;
}

.taint-icon:hover {
  color: #F56C6C;
  transform: scale(1.1);
}

.resource-info {
  font-size: 12px;
}

.resource-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 2px;
}

.resource-label {
  color: #909399;
  margin-right: 8px;
}

.resource-value {
  font-weight: 600;
  color: #303133;
}

.operation-buttons {
  display: flex;
  gap: 6px;
  flex-wrap: nowrap;
  align-items: center;
}

.operation-buttons .el-button {
  transition: all 0.3s ease;
}

.operation-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}
</style>