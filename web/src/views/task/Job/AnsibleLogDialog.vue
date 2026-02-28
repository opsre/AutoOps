<template>
  <el-dialog
    title="Ansible任务日志"
    v-model="dialogVisible"
    width="80%"
    top="5vh"
    append-to-body
    :before-close="handleClose"
  >
    <div class="log-header" v-if="logInfo">
      <el-descriptions :column="3" border>
        <el-descriptions-item label="任务ID">{{ logInfo.taskId }}</el-descriptions-item>
        <el-descriptions-item label="Work ID">{{ logInfo.workId }}</el-descriptions-item>
        <el-descriptions-item label="文件名">{{ logInfo.fileName }}</el-descriptions-item>
      </el-descriptions>
    </div>

    <div class="log-container">
      <div class="log-controls">
        <el-button type="primary" icon="Refresh" size="small" @click="refreshLog" :loading="refreshing">
          刷新日志
        </el-button>
        <el-button type="info" icon="Download" size="small" @click="downloadLog">
          下载日志
        </el-button>
        <el-tag v-if="wsConnected" type="success" size="small">
          🔄 实时连接
        </el-tag>
        <el-tag v-else-if="isCompleted" type="info" size="small">
          ✅ 任务完成
        </el-tag>
        <el-tag v-else type="warning" size="small">
          ⚪ 未连接
        </el-tag>
      </div>
      
      <div class="log-content-wrapper">
        <!-- 日志统计信息 -->
        <div v-if="logs.length > 0" class="log-stats">
          📊 日志行数: {{ logs.length }} 行 
          <span v-if="isCompleted" style="color: #67c23a;">| ✅ 任务已完成</span>
          <span v-else-if="wsConnected" style="color: #409eff;">| 🔄 任务运行中</span>
          <span v-else style="color: #909399;">| 空闲</span>
        </div>
        
        <div 
          ref="logContainer" 
          class="log-content" 
          :class="{ 'loading': loading, 'error': isError }"
        >
          <!-- 日志内容 -->
          <div v-if="loading" class="log-loading">
            <el-icon class="is-loading"><Loading /></el-icon>
            <div class="loading-text">
              <div>加载历史日志中...</div>
              <div style="font-size: 12px; color: #909399; margin-top: 5px;">
                后端可能需要2-3分钟处理时间，请耐心等待
              </div>
            </div>
          </div>
          
          <div v-else-if="isError" class="log-error">
            <el-icon><WarningFilled /></el-icon>
            {{ errorMessage }}
          </div>
          
          <div v-else-if="logs.length === 0" class="log-empty">
            <el-icon><Document /></el-icon>
            暂无日志数据
          </div>
          
          <div v-else class="log-lines">
            <div 
              v-for="(log, index) in logs" 
              :key="`${log.line_num || index}`"
              class="log-line"
              :class="getLogLineClass(log)"
            >
              <span class="line-number">{{ log.line_num || index + 1 }}</span>
              <span class="line-content" v-html="formatLogContent(log.content)"></span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch, nextTick, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Loading, WarningFilled, Document } from '@element-plus/icons-vue'
import { GetAnsibleTaskLog, GetAnsibleTaskLogByHistory } from '@/api/task'

// Props
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  historyMode: {
    type: Boolean,
    default: false
  },
  historyId: {
    type: [String, Number],
    default: null
  }
})

// Emits
const emit = defineEmits(['update:visible'])

// 响应式数据
const dialogVisible = ref(false)
const loading = ref(false)
const refreshing = ref(false)
const isError = ref(false)
const errorMessage = ref('')
const isCompleted = ref(false)
const wsConnected = ref(false)

const logInfo = ref(null)
const logs = ref([])
const logContainer = ref(null)

// 分别存储历史日志和实时日志
const historyLogs = ref([])
const realtimeLogs = ref([])
const lastWebSocketLineNum = ref(0)

let websocket = null
let reconnectTimer = null
let reconnectAttempts = 0
const maxReconnectAttempts = 5

// HTTP轮询相关变量
let pollingTimer = null
const pollingInterval = 3000 // 3秒轮询一次

// 监听visible变化
watch(() => props.visible, (newVal) => {
  dialogVisible.value = newVal
})

watch(dialogVisible, (newVal) => {
  emit('update:visible', newVal)
  
  if (newVal) {
    nextTick(() => {
      scrollToBottom()
    })
  } else {
    // 关闭时断开WebSocket连接和HTTP轮询
    disconnectWebSocket()
    stopHttpPolling()
  }
})

// 显示日志对话框
const show = (info) => {
  logInfo.value = { ...info }
  logs.value = []
  historyLogs.value = []
  realtimeLogs.value = []
  lastWebSocketLineNum.value = 0
  isError.value = false
  isCompleted.value = false
  loading.value = true

  dialogVisible.value = true

  // 根据模式选择获取方式
  if (props.historyMode && props.historyId) {
    refreshLog()
  } else {
    // 默认模式：先尝试获取历史日志，然后连接WebSocket
    refreshLog()
    connectWebSocket()
  }
}

// WebSocket连接
const connectWebSocket = () => {
  if (!logInfo.value) return
  
  const { taskId, workId } = logInfo.value
  const wsUrl = `ws://localhost:8080/api/v1/ws/task/ansible/${taskId}/log/${workId}`
  
  console.log('🔌 连接WebSocket:', wsUrl)
  console.log('📊 WebSocket连接参数:', { taskId, workId, wsUrl })
  
  try {
    websocket = new WebSocket(wsUrl)
    
    websocket.onopen = () => {
      console.log('✅ WebSocket连接成功')
      wsConnected.value = true
      reconnectAttempts = 0
      loading.value = false
    }
    
    websocket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        handleWebSocketMessage(message)
      } catch (error) {
        console.error('解析WebSocket消息失败:', error)
      }
    }
    
    websocket.onclose = (event) => {
      console.log('🔌 WebSocket连接关闭:', event.code, event.reason)
      wsConnected.value = false
      
      // 如果不是手动关闭且对话框仍然打开，尝试重连
      if (event.code !== 1000 && dialogVisible.value && reconnectAttempts < maxReconnectAttempts) {
        reconnectWebSocket()
      }
    }
    
    websocket.onerror = (error) => {
      console.error('❌ WebSocket连接错误:', error)
      wsConnected.value = false
      
      // WebSocket连接失败时，自动切换到HTTP轮询模式
      console.log('🔄 WebSocket连接失败，切换到HTTP轮询模式')
      startHttpPolling()
    }
    
  } catch (error) {
    console.error('❌ 创建WebSocket连接失败:', error)
    isError.value = true
    errorMessage.value = '无法创建WebSocket连接'
  }
}

// 处理WebSocket消息
const handleWebSocketMessage = (message) => {
  console.log('📨 收到WebSocket消息:', message)

  switch (message.type) {
    case 'log':
      // 检查是否是新的日志行（避免重复）
      if (message.line_num > lastWebSocketLineNum.value) {
        lastWebSocketLineNum.value = message.line_num

        // 添加到实时日志
        realtimeLogs.value.push({
          line_num: message.line_num,
          content: message.content,
          timestamp: message.timestamp,
          source: 'websocket'
        })

        // 合并日志并更新显示
        mergeLogs()

        // 自动滚动到底部
        nextTick(() => {
          scrollToBottom()
        })
      }
      break

    case 'complete':
      // 任务完成
      isCompleted.value = true
      wsConnected.value = false
      console.log(`✅ 任务完成，共${message.line_num}行日志`)
      ElMessage.success(`任务完成，共${logs.value.length}行日志`)

      // 关闭WebSocket连接
      if (websocket) {
        websocket.close(1000, 'Task completed')
      }
      break

    case 'error':
      // 处理错误
      isError.value = true
      errorMessage.value = message.content || '任务执行错误'
      ElMessage.error(errorMessage.value)
      break

    default:
      console.warn('未知的WebSocket消息类型:', message.type)
  }
}

// WebSocket重连
const reconnectWebSocket = () => {
  if (reconnectAttempts >= maxReconnectAttempts) {
    console.log('❌ WebSocket重连次数超限，停止重连')
    return
  }
  
  reconnectAttempts++
  const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 10000) // 指数退避，最大10秒
  
  console.log(`🔄 WebSocket重连 (${reconnectAttempts}/${maxReconnectAttempts})，${delay/1000}秒后重试`)
  
  reconnectTimer = setTimeout(() => {
    if (dialogVisible.value) {
      connectWebSocket()
    }
  }, delay)
}

// 断开WebSocket连接
const disconnectWebSocket = () => {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  
  if (websocket) {
    websocket.close(1000, 'Dialog closed')
    websocket = null
  }
  
  wsConnected.value = false
  reconnectAttempts = 0
  
  // 同时停止HTTP轮询
  stopHttpPolling()
}

// 启动HTTP轮询
const startHttpPolling = () => {
  console.log('🔄 启动HTTP轮询模式')
  
  // 先停止已有的轮询
  stopHttpPolling()
  
  // 立即获取一次日志
  refreshLog()
  
  // 开始定时轮询
  pollingTimer = setInterval(async () => {
    try {
      await refreshLog()
      console.log('📥 HTTP轮询获取日志成功')
    } catch (error) {
      console.error('❌ HTTP轮询获取日志失败:', error)
    }
  }, pollingInterval)
  
  console.log(`✅ HTTP轮询已启动，间隔${pollingInterval/1000}秒`)
}

// 停止HTTP轮询
const stopHttpPolling = () => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
    console.log('🛑 HTTP轮询已停止')
  }
}

// 刷新日志（获取历史日志）
const refreshLog = async () => {
  if (!logInfo.value || refreshing.value) return

  refreshing.value = true
  isError.value = false

  try {
    console.log('🔄 获取历史日志...')

    // 显示友好的加载提示
    if (logs.value.length === 0) {
      ElMessage.info('正在获取历史日志，后端可能需要处理时间...')
    }

    let response
    
    if (props.historyMode && props.historyId) {
      console.log('📚 获取归档历史日志:', {
        taskId: logInfo.value.taskId,
        workId: logInfo.value.workId,
        historyId: props.historyId
      })
      
      response = await GetAnsibleTaskLogByHistory({
        id: logInfo.value.taskId,
        workId: logInfo.value.workId,
        historyId: props.historyId
      })
      
      // 历史日志API可能返回 {data: {content: "..."}} 或 {data: "..."}
      // 需要做兼容处理
      let responseData = response.data || response
      
      if (responseData) {
        let content = ''
        
        content = responseData.data

        const parsedHistoryLogs = parseHistoryLogs(content)
        historyLogs.value = parsedHistoryLogs
        
        // 历史模式下不需要合并实时日志，直接使用
        logs.value = parsedHistoryLogs
        
        console.log(`✅ 获取到 ${parsedHistoryLogs.length} 行归档日志`)
        isError.value = false
        errorMessage.value = ''
        isCompleted.value = true // 历史日志视为已完成
        
        if (parsedHistoryLogs.length > 0) {
          ElMessage.success(`成功获取 ${parsedHistoryLogs.length} 行归档日志`)
        } else {
             ElMessage.info('该任务没有产生日志或日志为空')
        }
        
        // 提前返回，不再执行后续逻辑
        return
      }
    } else {
      response = await GetAnsibleTaskLog(logInfo.value.taskId, logInfo.value.workId)
    }

    if (response && response.data) {
      // 解析历史日志
      const parsedHistoryLogs = parseHistoryLogs(response.data)
      historyLogs.value = parsedHistoryLogs

      // 合并历史日志和实时日志
      mergeLogs()

      console.log(`✅ 获取到 ${parsedHistoryLogs.length} 行历史日志`)

      // 清除之前的错误状态
      isError.value = false
      errorMessage.value = ''

      if (parsedHistoryLogs.length > 0) {
        ElMessage.success(`成功获取 ${parsedHistoryLogs.length} 行历史日志`)
      }
    }

  } catch (error) {
    console.error('❌ 获取历史日志失败:', error)

    // 区分不同类型的错误
    if (error.code === 'ECONNABORTED') {
      // 超时错误 - 使用更友好的提示
      const contextMsg = error.contextMessage || '后端处理时间较长'
      console.warn('⚠️ 历史日志获取超时，但WebSocket实时日志仍可用')

      ElMessage({
        type: 'warning',
        message: `${contextMsg}，实时日志正常运行中`,
        duration: 5000,
        showClose: true
      })

      // 不设置isError，避免整个界面显示错误状态
      if (logs.value.length === 0) {
        ElMessage({
          type: 'info',
          message: '暂无历史日志，等待实时日志数据...',
          duration: 3000
        })
      }
    } else {
      // 其他错误类型
      isError.value = true
      errorMessage.value = '获取历史日志失败: ' + (error.message || '未知错误')
      ElMessage.error(errorMessage.value)
    }
  } finally {
    refreshing.value = false
    loading.value = false
  }
}

// 解析历史日志
const parseHistoryLogs = (logData) => {
  if (typeof logData !== 'string') return []

  const lines = logData.split('\n').filter(line => line.trim())
  return lines.map((line, index) => ({
    line_num: index + 1,
    content: line,
    timestamp: null,
    source: 'history'
  }))
}

// 合并历史日志和实时日志
const mergeLogs = () => {
  // 创建一个映射，用于存储每一行的内容（以内容为key避免重复）
  const logMap = new Map()

  // 首先添加历史日志
  historyLogs.value.forEach(log => {
    const key = log.content.trim()
    if (key && !logMap.has(key)) {
      logMap.set(key, {
        ...log,
        finalLineNum: log.line_num
      })
    }
  })

  // 然后添加实时日志（跳过已存在的内容）
  realtimeLogs.value.forEach(log => {
    const key = log.content.trim()
    if (key && !logMap.has(key)) {
      // 实时日志的行号应该继续历史日志的行号
      const maxHistoryLineNum = historyLogs.value.length
      logMap.set(key, {
        ...log,
        finalLineNum: maxHistoryLineNum + realtimeLogs.value.indexOf(log) + 1
      })
    }
  })

  // 转换为数组并按行号排序
  const mergedLogs = Array.from(logMap.values()).sort((a, b) => a.finalLineNum - b.finalLineNum)

  // 重新分配连续的行号
  logs.value = mergedLogs.map((log, index) => ({
    ...log,
    line_num: index + 1
  }))

  console.log(`🔄 日志合并完成: 历史${historyLogs.value.length}行 + 实时${realtimeLogs.value.length}行 = 最终${logs.value.length}行`)
}

// 格式化日志内容
const formatLogContent = (content) => {
  if (!content) return ''
  
  // ANSI颜色代码转换为HTML样式
  return content
    .replace(/\u001b\[0;32m/g, '<span style="color: #67c23a;">')  // 绿色
    .replace(/\u001b\[0;31m/g, '<span style="color: #f56c6c;">')  // 红色
    .replace(/\u001b\[0;33m/g, '<span style="color: #e6a23c;">')  // 黄色
    .replace(/\u001b\[0;34m/g, '<span style="color: #409eff;">')  // 蓝色
    .replace(/\u001b\[0m/g, '</span>')                              // 结束颜色
    .replace(/\u001b\[[\d;]*m/g, '')                                // 清除其他ANSI代码
}

// 获取日志行样式
const getLogLineClass = (log) => {
  const content = log.content.toLowerCase()
  if (content.includes('error') || content.includes('failed')) {
    return 'log-error-line'
  } else if (content.includes('warn')) {
    return 'log-warn-line'
  } else if (content.includes('success') || content.includes('ok')) {
    return 'log-success-line'
  }
  return ''
}

// 滚动到底部
const scrollToBottom = () => {
  if (logContainer.value) {
    logContainer.value.scrollTop = logContainer.value.scrollHeight
  }
}

// 下载日志
const downloadLog = () => {
  try {
    const logContent = logs.value.map(log => log.content).join('\n')
    const blob = new Blob([logContent], { type: 'text/plain;charset=utf-8' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `ansible-task-${logInfo.value?.taskId}-work-${logInfo.value?.workId}-logs.txt`
    a.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('日志下载成功')
  } catch (error) {
    console.error('下载日志失败:', error)
    ElMessage.error('下载日志失败')
  }
}

// 处理关闭
const handleClose = () => {
  disconnectWebSocket()
  dialogVisible.value = false
}

// 组件卸载时清理
onUnmounted(() => {
  disconnectWebSocket()
  stopHttpPolling()
})

// 导出方法供父组件使用
defineExpose({
  show
})
</script>

<style scoped>
.log-header {
  margin-bottom: 20px;
}

.log-container {
  height: 600px;
  display: flex;
  flex-direction: column;
}

.log-controls {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 0;
  border-bottom: 1px solid #e4e7ed;
}

.log-content-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.log-stats {
  background: #f8f9fa;
  padding: 5px 10px;
  margin-bottom: 10px;
  font-size: 11px;
  color: #666;
  border-left: 3px solid #28a745;
}

.log-content {
  flex: 1;
  overflow: auto;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  background: #1e1e1e;
  color: #d4d4d4;
  font-family: 'Courier New', Consolas, monospace;
  font-size: 13px;
  line-height: 1.4;
  position: relative;
}

.log-loading, .log-error, .log-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  flex-direction: column;
  gap: 10px;
  color: #909399;
  font-size: 14px;
}

.loading-text {
  text-align: center;
  line-height: 1.5;
}

.log-error {
  color: #f56c6c;
}

.log-lines {
  padding: 10px;
}

.log-line {
  display: flex;
  margin-bottom: 2px;
  min-height: 18px;
  align-items: flex-start;
}

.line-number {
  width: 60px;
  text-align: right;
  color: #606266;
  margin-right: 15px;
  font-size: 12px;
  flex-shrink: 0;
  padding-top: 1px;
}

.line-content {
  flex: 1;
  word-break: break-all;
  white-space: pre-wrap;
}

.log-error-line .line-content {
  color: #f56c6c;
}

.log-warn-line .line-content {
  color: #e6a23c;
}

.log-success-line .line-content {
  color: #67c23a;
}

/* 滚动条样式 */
.log-content::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.log-content::-webkit-scrollbar-track {
  background: #2d2d2d;
}

.log-content::-webkit-scrollbar-thumb {
  background: #555;
  border-radius: 4px;
}

.log-content::-webkit-scrollbar-thumb:hover {
  background: #777;
}
</style>