<template>
  <el-dialog
    v-model="visible"
    :title="`容器日志 - ${podName}`"
    width="80%"
    align-center
    class="logs-dialog"
    destroy-on-close
    @opened="handleOpened"
    @closed="handleClose"
  >
    <div class="logs-tab-content">
      <div class="logs-controls">
        <el-row :gutter="12" style="margin-bottom: 16px;">
          <el-col :span="6">
            <el-select v-model="selectedContainer" placeholder="选择容器" style="width: 100%;" @change="fetchLogs">
              <el-option 
                v-for="container in parsedContainers" 
                :key="container.name"
                :label="container.name" 
                :value="container.name"
              />
            </el-select>
          </el-col>
          <el-col :span="4">
            <el-select v-model="logTailLines" placeholder="行数" style="width: 100%;" @change="fetchLogs">
              <el-option label="最近100行" :value="100" />
              <el-option label="最近300行" :value="300" />
              <el-option label="最近500行" :value="500" />
              <el-option label="最近1000行" :value="1000" />
            </el-select>
          </el-col>
          <el-col :span="6">
            <div style="display: flex; align-items: center; height: 32px; gap: 12px; white-space: nowrap;">
              <el-checkbox v-model="followLogs" @change="toggleFollow">实时跟踪</el-checkbox>
              <el-checkbox v-model="showPreviousLogs" @change="fetchLogs">上个容器退出日志</el-checkbox>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="logs-actions" style="display: flex; justify-content: flex-end; align-items: center;">
              <el-button size="small" type="primary" :icon="Refresh" @click="fetchLogs" :loading="logsLoading">刷新</el-button>
              <el-button size="small" type="success" :icon="Download" @click="handleDownloadLogs">下载</el-button>
            </div>
          </el-col>
        </el-row>
      </div>
      
      <div class="logs-display" v-loading="logsLoading">
        <div v-if="!selectedContainer" class="logs-placeholder">
          <el-empty description="请选择容器查看日志" />
        </div>
        <div v-else class="logs-content-display">
          <div class="logs-header" style="margin-bottom: 8px; font-size: 13px; color: #909399;">
            <span class="container-name">{{ selectedContainer }}</span> 
            <span class="logs-info" style="margin-left: 10px;">
              最近 {{ logTailLines || 1000 }} 行 
            </span>
            <span v-if="followLogs" style="margin-left: 10px; color: #67c23a;">
              <el-icon class="is-loading" style="margin-right:4px;"><Loading /></el-icon>实时接收中...
            </span>
          </div>
          
          <div class="logs-editor-container">
            <div ref="terminalElement" class="xterm-container"></div>
          </div>
        </div>
      </div>
    </div>
    
    <template #footer>
      <el-button @click="visible = false">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onBeforeUnmount, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Download, Loading } from '@element-plus/icons-vue'
import { Terminal as XTerm } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import k8sApi from '@/api/k8s'
import storage from '@/utils/storage'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  clusterId: { type: [String, Number] },
  namespace: { type: String },
  podName: { type: String },
  containers: { type: Array, default: () => [] }
})

const emit = defineEmits(['update:modelValue', 'close'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// UI States
const selectedContainer = ref('')
const logTailLines = ref(500)
const followLogs = ref(false)
const showPreviousLogs = ref(false)
const logsLoading = ref(false)
const currentLogs = ref('')

const wsInstance = ref(null)
const terminalElement = ref(null)
let terminal = null
let fitAddon = null

const parsedContainers = computed(() => {
  // containers 可能是数组或者对象数组
  return props.containers.map(c => {
    if (typeof c === 'string') return { name: c }
    return c
  })
})

const initTerminal = () => {
  if (terminal) return
  
  terminal = new XTerm({
    cursorBlink: false,
    fontSize: 13, /* 对标 Ansible 的 fontSize: 13 */
    fontFamily: "'Courier New', Consolas, monospace", /* 对标 Ansible 的英文字体 */
    lineHeight: 1.4, /* 对标 Ansible 的行高 */
    theme: {
      background: '#1e1e1e', /* 对标 Ansible 的主黑底 */
      foreground: '#d4d4d4', /* 对标 Ansible 的字体主色 */
      cursor: '#ffffff',
      selection: '#ffffff80',
      black: '#000000',
      red: '#ff5c5c',
      green: '#5af78e',
      yellow: '#f3f99d',
      blue: '#57c7ff',
      magenta: '#ff6ac1',
      cyan: '#9aedfe',
      white: '#f1f1f0',
      brightBlack: '#686868',
      brightRed: '#ff5c5c',
      brightGreen: '#5af78e',
      brightYellow: '#f3f99d',
      brightBlue: '#57c7ff',
      brightMagenta: '#ff6ac1',
      brightCyan: '#9aedfe',
      brightWhite: '#f1f1f0'
    },
    disableStdin: true,
    convertEol: true,
    scrollback: 10000,
    cols: 100,
    rows: 40,
    smoothScrollDuration: 200, // 开启平滑滚动，让滚动时的视觉更流式、更顺滑
  })
  
  fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)
  
  if (terminalElement.value) {
    terminal.open(terminalElement.value)
    fitAddon.fit()
    
    const resizeObserver = new ResizeObserver(() => {
      if (fitAddon) fitAddon.fit()
    })
    resizeObserver.observe(terminalElement.value)
  }
}

const handleOpened = () => {
  initTerminal()
  if (selectedContainer.value) {
    fetchLogs()
  }
}

// 统一使用 WebSocket 获取日志
const startFollow = () => {
  stopFollow()
  if (!selectedContainer.value) return
  
  currentLogs.value = '' // 清空日志以便下载
  if (terminal) {
    terminal.clear()
    terminal.writeln(`\x1B[1;36m>> 正在准备连接到容器 ${selectedContainer.value}...\x1B[0m`)
  }
  
  logsLoading.value = true
  
  const params = {
    container: selectedContainer.value,
    tailLines: logTailLines.value,
    previous: showPreviousLogs.value,
    follow: followLogs.value,
    token: storage.getItem('token')
  }
  
  const wsUrl = k8sApi.getPodLogsWsUrl(
    props.clusterId,
    props.namespace,
    props.podName,
    params
  )
  
  try {
    const ws = new WebSocket(wsUrl)
    wsInstance.value = ws
    
    ws.onopen = () => {
      if (wsInstance.value !== ws) return
      logsLoading.value = false
      if (terminal) {
        terminal.clear()
        terminal.writeln(`\x1B[1;32m✓ 成功连接到 ${selectedContainer.value} 日志服务...\x1B[0m\n`)
      }
      console.log('日志WebSocket连接成功')
    }
    
    ws.onmessage = (event) => {
      if (wsInstance.value !== ws) return
      currentLogs.value += event.data
      if (terminal) terminal.write(event.data)
    }
    
    ws.onerror = (err) => {
      if (wsInstance.value !== ws) return
      console.error('实时日志WS出错', err)
      if (terminal) terminal.writeln('\n\x1B[1;31m❌ 连接发生错误或中断\x1B[0m')
      logsLoading.value = false
    }
    
    ws.onclose = () => {
      if (wsInstance.value !== ws) return
      console.log('实时日志连接关闭')
      if (terminal) {
        if (followLogs.value) {
          terminal.writeln('\n\x1B[1;33m🔌 连接已结束\x1B[0m')
        } else {
          // 短连接结束后可以给个友好的EOF标记
          terminal.writeln('\n\x1B[1;30m(EOF)\x1B[0m')
        }
      }
      logsLoading.value = false
    }
  } catch (error) {
    console.error('建连失败', error)
    if (terminal) terminal.writeln('\n\x1B[1;31m❌ 建连失败\x1B[0m')
    logsLoading.value = false
  }
}

// 统一封装入口
const fetchLogs = () => {
  if (!selectedContainer.value) return
  startFollow()
}

const toggleFollow = () => {
  startFollow()
}

const initLogState = () => {
  if (parsedContainers.value.length > 0) {
    selectedContainer.value = parsedContainers.value[0].name
  }
}

// 监听打开状态
watch(() => visible.value, (newVal) => {
  if (newVal) {
    currentLogs.value = ''
    followLogs.value = false
    showPreviousLogs.value = false
    initLogState()
  } else {
    stopFollow()
  }
})

const stopFollow = () => {
  if (wsInstance.value) {
    wsInstance.value.close()
    wsInstance.value = null
  }
}

const handleClose = () => {
  stopFollow()
  currentLogs.value = ''
  if (terminal) {
    terminal.dispose()
    terminal = null
  }
  emit('close')
}

const handleDownloadLogs = () => {
  if (!currentLogs.value) {
    ElMessage.warning('没有日志内容可下载')
    return
  }
  
  const blob = new Blob([currentLogs.value], { type: 'text/plain;charset=utf-8' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${props.podName}-${selectedContainer.value}-logs.txt`
  a.style.display = 'none'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  window.URL.revokeObjectURL(url)
}

onBeforeUnmount(() => {
  stopFollow()
})
</script>

<style scoped>
.logs-dialog :deep(.el-dialog__body) {
  padding: 10px 20px;
}
.logs-editor-container {
  border: 1px solid #e4e7ed; /* 采用 AnsibleLogDialog 类似的浅色外边框统一风格 */
  border-radius: 6px;
  overflow: hidden;
  height: 600px; /* 提升高度对齐 AnsibleLogDialog 的高度 */
  background: #1e1e1e; /* 对标 AnsibleLogDialog 的背景黑 */
}
.xterm-container {
  height: 100%;
  width: 100%;
  box-sizing: border-box;
  background-color: #1e1e1e;
  padding: 10px 0 10px 10px; /* 彻底移除右侧 padding 让滚动条贴边 */
}

/* 美化 xterm 的内部滚动条 (模仿 AnsibleLogDialog 的基础滚动条表现) */
.xterm-container :deep(.xterm-viewport) {
  /* xterm.js 特性修复：使用 padding-right 确保画布排版被强制推离右侧，防止截断 */
  box-sizing: border-box !important;
  padding-right: 15px !important;
}

.xterm-container :deep(.xterm-viewport::-webkit-scrollbar) {
  width: 8px; /* 统一宽度为 AnsibleLogDialog 的 8px */
  height: 8px;
}

.xterm-container :deep(.xterm-viewport::-webkit-scrollbar-track) {
  background: #2d2d2d; /* 统一颜色为 AnsibleLogDialog 的 #2d2d2d */
  border-radius: 0;
}

.xterm-container :deep(.xterm-viewport::-webkit-scrollbar-thumb) {
  background: #555; /* 统一滑块为 #555 */
  border-radius: 4px;
}

.xterm-container :deep(.xterm-viewport::-webkit-scrollbar-thumb:hover) {
  background: #777;
}
</style>


