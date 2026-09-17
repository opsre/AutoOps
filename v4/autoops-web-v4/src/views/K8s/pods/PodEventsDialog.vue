<template>
  <el-dialog
    :model-value="visible"
    :title="`Pod事件 - ${podName}`"
    width="1000px"
    class="pod-events-dialog"
    @update:model-value="handleClose"
  >
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="8" animated />
    </div>

    <div v-else-if="events.length === 0" class="empty-events">
      <el-empty description="暂无事件记录" />
    </div>

    <div v-else class="events-container">
      <el-table
        :data="events"
        stripe
        style="width: 100%"
        size="small"
        :default-sort="{ prop: 'lastTime', order: 'descending' }"
      >
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag
              :type="getEventTypeTag(row.type)"
              size="small"
              effect="dark"
            >
              {{ row.type }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="reason" label="原因" width="150">
          <template #default="{ row }">
            <span class="event-reason">{{ row.reason }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="message" label="消息" min-width="300">
          <template #default="{ row }">
            <div class="event-message">{{ row.message }}</div>
          </template>
        </el-table-column>

        <el-table-column prop="count" label="次数" width="80">
          <template #default="{ row }">
            <el-tag size="small" :type="row.count > 1 ? 'warning' : 'info'">
              {{ row.count }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="firstTime" label="首次发生" width="180" sortable>
          <template #default="{ row }">
            <div class="timestamp">{{ formatTimestamp(row.firstTime) }}</div>
          </template>
        </el-table-column>

        <el-table-column prop="lastTime" label="最后发生" width="180" sortable>
          <template #default="{ row }">
            <div class="timestamp">{{ formatTimestamp(row.lastTime) }}</div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button type="primary" @click="refreshEvents" :loading="loading">
          刷新
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, defineProps, defineEmits } from 'vue'
import { ElMessage } from 'element-plus'
import k8sApi from '@/api/k8s'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  clusterId: {
    type: String,
    default: ''
  },
  namespaceName: {
    type: String,
    default: ''
  },
  podName: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:visible', 'close'])

const loading = ref(false)
const events = ref([])

// 监听对话框打开状态
watch(() => props.visible, (newVal) => {
  if (newVal && props.clusterId && props.namespaceName && props.podName) {
    loadEvents()
  }
})

const handleClose = () => {
  emit('update:visible', false)
  emit('close')
}

const loadEvents = async () => {
  if (!props.clusterId || !props.namespaceName || !props.podName) {
    return
  }

  loading.value = true

  try {
    const response = await k8sApi.getPodEvents(props.clusterId, props.namespaceName, props.podName)
    const responseData = response.data || response

    if (responseData.code === 200 && responseData.data) {
      events.value = responseData.data.events || []
    } else {
      ElMessage.error(responseData.message || '获取Pod事件失败')
      events.value = []
    }
  } catch (error) {
    console.error('获取Pod事件失败:', error)
    ElMessage.error('获取Pod事件失败，请重试')
    events.value = []
  } finally {
    loading.value = false
  }
}

const refreshEvents = () => {
  loadEvents()
}

const getEventTypeTag = (type) => {
  const tagMap = {
    'Warning': 'danger',
    'Normal': 'success'
  }
  return tagMap[type] || 'info'
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return '-'

  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    return timestamp
  }
}
</script>

<style scoped>
.loading-container {
  padding: 20px;
}

.empty-events {
  padding: 40px 0;
}

.events-container {
  max-height: 600px;
  overflow-y: auto;
}

.event-reason {
  font-weight: 500;
  color: #6b9bd1;
}

.event-message {
  line-height: 1.4;
  word-break: break-word;
  color: #555;
}

.timestamp {
  font-family: monospace;
  font-size: 12px;
  color: #666;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 对话框样式 */
.pod-events-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.pod-events-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.pod-events-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.pod-events-dialog :deep(.el-dialog__body) {
  padding: 20px;
}

.pod-events-dialog :deep(.el-dialog__footer) {
  padding: 20px 24px;
  border-top: 1px solid #e5e7eb;
}

/* 事件类型样式 */
.pod-events-dialog :deep(.el-tag.el-tag--danger) {
  background: linear-gradient(135deg, #ff7875, #ff4d4f);
  border-color: transparent;
}

.pod-events-dialog :deep(.el-tag.el-tag--success) {
  background: linear-gradient(135deg, #73d13d, #52c41a);
  border-color: transparent;
}
</style>