<template>
  <el-dialog
    :model-value="visible"
    :title="`Pod列表 - ${workload.name || ''}`"
    width="1200px"
    class="pod-list-dialog"
    @update:model-value="handleClose"
  >
    <el-table
      :data="workload.pods || []"
      stripe
      style="width: 100%"
      size="small"
    >
      <el-table-column prop="name" label="Pod名称" min-width="200">
        <template #default="{ row }">
          <div class="pod-name-container">
            <img src="@/assets/image/k8s.svg" alt="k8s" class="k8s-icon" />
            <span class="pod-name">{{ row.name }}</span>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="status" label="状态" min-width="100">
        <template #default="{ row }">
          <el-tag
            :type="getPodStatusTag(row.status)"
            size="small"
            effect="dark"
          >
            {{ getPodStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="restartCount" label="重启次数" min-width="100">
        <template #default="{ row }">
          <el-tag size="small" :type="row.restartCount > 0 ? 'warning' : 'success'">
            {{ row.restartCount }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="nodeName" label="节点" min-width="150">
        <template #default="{ row }">
          <el-tag size="small" type="info">{{ row.nodeName }}</el-tag>
        </template>
      </el-table-column>

      <el-table-column label="IP地址" min-width="140">
        <template #default="{ row }">
          <div class="ip-info">
            <div class="ip-item">
              <span class="ip-label">Pod:</span>
              <span class="ip-value">{{ row.podIP }}</span>
            </div>
            <div class="ip-item">
              <span class="ip-label">Host:</span>
              <span class="ip-value">{{ row.hostIP }}</span>
            </div>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="资源" min-width="140">
        <template #default="{ row }">
          <div class="resource-info">
            <div class="resource-item">
              <span class="resource-label">CPU:</span>
              <span class="resource-value">{{ formatCpu(row.resources?.requests?.cpu) }}</span>
            </div>
            <div class="resource-item">
              <span class="resource-label">Memory:</span>
              <span class="resource-value">{{ formatMemory(row.resources?.requests?.memory) }}</span>
            </div>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="age" label="运行时间" min-width="120">
        <template #default="{ row }">
          <div class="time-info">
            <el-tag size="small" type="info">{{ row.age }}</el-tag>
            <div v-if="row.runningTime" class="running-time">
              运行: {{ formatRunningTime(row.runningTime) }}
            </div>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <div class="operation-buttons">
            <el-tooltip content="日志" placement="top">
              <el-button
                type="primary"
                :icon="Document"
                size="small"
                circle
                @click="$emit('view-logs', row)"
              />
            </el-tooltip>

            <el-tooltip content="YAML" placement="top">
              <el-button
                type="success"
                :icon="View"
                size="small"
                circle
                @click="$emit('view-yaml', row)"
              />
            </el-tooltip>

            <el-tooltip content="重构" placement="top">
              <el-button
                type="warning"
                size="small"
                circle
                v-authority="['k8s:workload:restart']"
                @click="$emit('rebuild-pod', row)"
              >
                <img src="@/assets/image/重构.svg" alt="重构" width="14" height="14" style="filter: brightness(0) invert(1);" />
              </el-button>
            </el-tooltip>

            <el-tooltip content="事件" placement="top">
              <el-button
                type="danger"
                size="small"
                circle
                @click="$emit('view-events', row)"
              >
                <img src="@/assets/image/事件.svg" alt="事件" width="14" height="14" />
              </el-button>
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </el-dialog>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'
import { Document, View } from '@element-plus/icons-vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  workload: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:visible', 'close', 'view-logs', 'view-yaml', 'rebuild-pod', 'view-events'])

const handleClose = () => {
  emit('update:visible', false)
  emit('close')
}

// 工具函数
const formatCpu = (cpuStr) => {
  if (!cpuStr || cpuStr === '0' || cpuStr === '') return '-'
  return cpuStr
}

const formatMemory = (memoryStr) => {
  if (!memoryStr || memoryStr === '0' || memoryStr === '') return '-'

  if (memoryStr.endsWith('Ki')) {
    const kb = parseInt(memoryStr.replace('Ki', ''))
    if (kb < 1024) return memoryStr
    const mb = (kb / 1024).toFixed(1)
    return `${mb}Mi`
  }

  if (memoryStr.endsWith('Mi')) {
    const mi = parseInt(memoryStr.replace('Mi', ''))
    if (mi < 1024) return memoryStr
    const gi = (mi / 1024).toFixed(1)
    return `${gi}Gi`
  }

  if (memoryStr.endsWith('Gi')) {
    return memoryStr
  }

  return memoryStr
}

const formatRunningTime = (runningTimeStr) => {
  if (!runningTimeStr) return '-'

  // 解析 "30h22m1.563771s" 格式
  const match = runningTimeStr.match(/^(\d+h)?(\d+m)?(\d+(?:\.\d+)?s)?$/)
  if (!match) return runningTimeStr

  const [, hours, minutes, seconds] = match
  let result = ''

  if (hours) result += hours
  if (minutes) result += minutes
  if (seconds && !hours && !minutes) {
    const sec = parseFloat(seconds.replace('s', ''))
    result += `${Math.round(sec)}s`
  }

  return result || runningTimeStr
}

const getPodStatusTag = (status) => {
  const tagMap = {
    'Running': 'success',
    'Pending': 'warning',
    'Failed': 'danger',
    'Succeeded': 'success',
    'Unknown': 'info'
  }
  return tagMap[status] || 'info'
}

const getPodStatusText = (status) => {
  const textMap = {
    'Running': '运行中',
    'Pending': '等待中',
    'Failed': '失败',
    'Succeeded': '成功',
    'Unknown': '未知'
  }
  return textMap[status] || status
}
</script>

<style scoped>
.pod-name-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.k8s-icon {
  width: 16px;
  height: 16px;
}

.pod-name {
  font-weight: 500;
  color: #6b9bd1;
}

.ip-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
}

.ip-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.ip-label {
  font-weight: 500;
  color: #666;
  min-width: 30px;
}

.ip-value {
  color: #409eff;
  font-family: monospace;
}

.resource-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
}

.resource-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.resource-label {
  font-weight: 500;
  color: #666;
  min-width: 45px;
}

.resource-value {
  color: #67c23a;
  font-family: monospace;
}

.time-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.running-time {
  font-size: 11px;
  color: #909399;
}

.operation-buttons {
  display: flex;
  gap: 4px;
  justify-content: center;
}

/* 对话框样式 */
.pod-list-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.pod-list-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.pod-list-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}
</style>