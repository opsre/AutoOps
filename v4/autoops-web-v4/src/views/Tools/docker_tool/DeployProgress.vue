<template>
  <el-dialog
    v-model="visible"
    title="部署进度"
    width="800px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <!-- 状态卡片 -->
    <el-card class="status-card">
      <div class="status-info">
        <div class="status-icon">
          <el-icon :size="60" :color="statusColor">
            <Loading v-if="deployInfo.status === 0" class="rotating" />
            <CircleCheck v-else-if="deployInfo.status === 1" />
            <CircleClose v-else-if="deployInfo.status === 3" />
            <Remove v-else />
          </el-icon>
        </div>
        <div class="status-text">
          <h3>{{ deployInfo.serviceName }} {{ deployInfo.version }}</h3>
          <el-tag :type="statusTagType" size="large">
            {{ deployInfo.statusText }}
          </el-tag>
          <p class="host-info">
            主机: {{ deployInfo.hostName }} ({{ deployInfo.hostIp }})
          </p>
          <p class="install-dir">
            安装目录: {{ deployInfo.installDir }}
          </p>
        </div>
      </div>
    </el-card>

    <!-- 部署日志 -->
    <div class="deploy-log">
      <div class="log-header">
        <h4>部署日志</h4>
        <el-button
          size="small"
          @click="copyLog"
          icon="CopyDocument"
        >
          复制日志
        </el-button>
      </div>
      <div class="log-content" ref="logRef">
        <pre>{{ deployInfo.deployLog || '正在获取日志...' }}</pre>
      </div>
    </div>

    <template #footer>
      <el-button
        v-if="deployInfo.status === 0"
        type="info"
        disabled
      >
        部署中，请稍候...
      </el-button>
      <el-button
        v-else
        type="primary"
        @click="visible = false"
      >
        关闭
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onUnmounted, nextTick } from 'vue'
import { getDeployStatus } from '@/api/tool'
import { ElMessage } from 'element-plus'
import {
  Loading,
  CircleCheck,
  CircleClose,
  Remove
} from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: Boolean,
  deployId: [Number, String]
})

const emit = defineEmits(['update:modelValue', 'close'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const deployInfo = ref({
  status: 0,
  statusText: '部署中',
  serviceName: '',
  version: '',
  hostName: '',
  hostIp: '',
  installDir: '',
  deployLog: ''
})

const logRef = ref()
let pollingTimer = null
let hasNotifiedFinal = false

// 状态颜色
const statusColor = computed(() => {
  const colors = {
    0: '#409eff', // 部署中
    1: '#67c23a', // 运行中
    2: '#909399', // 已停止
    3: '#f56c6c'  // 部署失败
  }
  return colors[deployInfo.value.status] || '#909399'
})

// 状态标签类型
const statusTagType = computed(() => {
  const types = {
    0: 'primary',
    1: 'success',
    2: 'info',
    3: 'danger'
  }
  return types[deployInfo.value.status] || 'info'
})

// 获取部署状态
const fetchDeployStatus = async () => {
  if (!props.deployId) return

  try {
    const res = await getDeployStatus(props.deployId)
    if (res.data?.code === 200) {
      deployInfo.value = res.data.data || deployInfo.value

      // 自动滚动到日志底部
      nextTick(() => {
        if (logRef.value) {
          logRef.value.scrollTop = logRef.value.scrollHeight
        }
      })

      // 如果部署完成或失败，停止轮询并仅提示一次
      if (res.data.data?.status === 1 || res.data.data?.status === 3) {
        stopPolling()
        if (!hasNotifiedFinal) {
          hasNotifiedFinal = true
          if (res.data.data?.status === 1) {
            ElMessage.success('部署成功！')
          } else {
            ElMessage.error('部署失败，请查看日志')
          }
        }
      }
    }
  } catch (error) {
    console.error('获取部署状态失败:', error)
  }
}

// 开始轮询
const startPolling = () => {
  // 避免重复启动轮询
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
  hasNotifiedFinal = false
  fetchDeployStatus() // 立即获取一次
  pollingTimer = setInterval(fetchDeployStatus, 3000) // 每3秒轮询一次
}

// 停止轮询
const stopPolling = () => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
}

// 复制日志
const copyLog = () => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(deployInfo.value.deployLog)
    ElMessage.success('日志已复制到剪贴板')
  }
}

const handleClose = () => {
  stopPolling()
  emit('close')
}

// 监听 deployId 变化
watch(() => props.deployId, (newVal) => {
  if (newVal && props.modelValue) {
    startPolling()
  }
}, { immediate: true })

// 监听对话框显示状态
watch(() => props.modelValue, (newVal) => {
  if (newVal && props.deployId) {
    startPolling()
  } else {
    stopPolling()
  }
})

// 组件卸载时清理
onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.status-card {
  margin-bottom: 20px;
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.2);
}

.status-info {
  display: flex;
  align-items: center;
  gap: 20px;
}

.status-icon {
  flex-shrink: 0;
}

.rotating {
  animation: rotate 1.5s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.status-text {
  flex: 1;
}

.status-text h3 {
  margin: 0 0 10px 0;
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.host-info, .install-dir {
  margin: 5px 0;
  color: #666;
  font-size: 14px;
}

.deploy-log {
  border: 1px solid #dcdfe6;
  border-radius: 12px;
  overflow: hidden;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-bottom: 1px solid #dcdfe6;
}

.log-header h4 {
  margin: 0;
  color: #fff;
  font-weight: 600;
}

.log-content {
  height: 400px;
  overflow-y: auto;
  padding: 15px;
  background: #1e1e1e;
  color: #d4d4d4;
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.log-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

:deep(.el-dialog) {
  border-radius: 12px;
}

:deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  padding: 20px;
  border-radius: 12px 12px 0 0;
}

:deep(.el-dialog__title) {
  color: #fff;
  font-weight: 600;
}

.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
  margin-bottom: 10px;
}

.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
