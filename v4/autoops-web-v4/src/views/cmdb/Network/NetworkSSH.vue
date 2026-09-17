<template>
  <div class="terminal-container" ref="terminalContainer">
    <!-- 终端工具栏 -->
    <div class="terminal-toolbar">
      <div class="toolbar-left">
        <el-button
          size="small"
          type="success"
          :icon="Connection"
          @click="connect"
          :disabled="isConnecting || !deviceInfo || isConnected"
          :loading="isConnecting"
        >
          {{ isConnecting ? '连接中' : '连接' }}
        </el-button>
        <el-button
          size="small"
          type="warning"
          :icon="SwitchButton"
          @click="reconnect"
          :disabled="!isConnected && !isConnecting"
        >
          重连
        </el-button>
        <el-button
          size="small"
          type="danger"
          :icon="Close"
          @click="disconnect"
          :disabled="!isConnected || isConnecting"
        >
          断开
        </el-button>
        <el-button
          size="small"
          type="info"
          :icon="Delete"
          @click="clear"
          :disabled="!term"
        >
          清屏
        </el-button>
        <el-button
          size="small"
          type="primary"
          :icon="Back"
          @click="goBack"
        >
          返回设备列表
        </el-button>
      </div>

      <div class="toolbar-right">
        <span class="connection-status" :class="statusClass">
          <span class="status-dot"></span>
          {{ statusText }}
        </span>
        <span v-if="deviceInfo" class="host-info">
          <el-tag type="success" size="small">{{ deviceTypeText }}</el-tag>
          {{ deviceInfo.deviceName }} ({{ deviceInfo.ipAddress }})
        </span>
      </div>
    </div>

    <!-- xterm 终端区域 -->
    <div id="xterm-wrapper" class="xterm-wrapper">
      <div id="network-xterm" class="xterm-terminal" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Connection, SwitchButton, Close, Delete, Back } from '@element-plus/icons-vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import storage from '@/utils/storage'
import API from '@/api/cmdb'

// 路由
const route = useRoute()
const router = useRouter()

// 状态管理
const terminalContainer = ref(null)
const term = ref(null)
const fitAddon = ref(null)
const socket = ref(null)
const isConnecting = ref(false)
const isConnected = ref(false)
const _termInitialized = ref(false)
const _inputHandlerBound = ref(false)

// 设备信息
const deviceInfo = ref(null)
const deviceId = ref(null)

// 设备类型文本映射
const deviceTypeMap = {
  'switch': '交换机',
  'router': '路由器',
  'windows': 'Windows'
}

const deviceTypeText = computed(() => {
  if (!deviceInfo.value) return ''
  return deviceTypeMap[deviceInfo.value.deviceType] || deviceInfo.value.deviceType
})

// 连接状态计算属性
const statusClass = computed(() => {
  if (isConnected.value) return 'status-connected'
  if (isConnecting.value) return 'status-connecting'
  return 'status-disconnected'
})

const statusText = computed(() => {
  if (isConnected.value) return '已连接'
  if (isConnecting.value) return '连接中...'
  return '未连接'
})

// 初始化终端
const initTerm = async () => {
  try {
    if (_termInitialized.value) {
      return
    }

    // 清理旧的终端实例
    if (term.value) {
      try {
        if (fitAddon.value && !fitAddon.value._disposed) {
          fitAddon.value.dispose()
        }
        term.value.dispose()
      } catch (e) {
        console.warn('Terminal dispose error:', e)
      }
    }

    _termInitialized.value = true

    // 创建新的终端实例
    term.value = new Terminal({
      fontSize: 14,
      fontFamily: 'Consolas, "Courier New", monospace',
      cursorBlink: true,
      cursorStyle: 'block',
      theme: {
        background: '#1e1e1e',
        foreground: '#d4d4d4',
        cursor: '#ffffff',
        selection: 'rgba(255, 255, 255, 0.3)',
        black: '#000000',
        red: '#cd3131',
        green: '#0dbc79',
        yellow: '#e5e510',
        blue: '#2472c8',
        magenta: '#bc3fbc',
        cyan: '#11a8cd',
        white: '#e5e5e5',
        brightBlack: '#666666',
        brightRed: '#f14c4c',
        brightGreen: '#23d18b',
        brightYellow: '#f5f543',
        brightBlue: '#3b8eea',
        brightMagenta: '#d670d6',
        brightCyan: '#29b8db',
        brightWhite: '#ffffff'
      },
      cols: 120,
      rows: 30,
      scrollback: 1000,
      convertEol: false,
      allowTransparency: false,
      lineHeight: 1.2,
      letterSpacing: 0,
      bellStyle: 'none',
      rightClickSelectsWord: true,
      allowProposedApi: true,
      disableStdin: false,
      windowsMode: false
    })

    // 创建并加载 FitAddon
    fitAddon.value = new FitAddon()
    term.value.loadAddon(fitAddon.value)

    // 挂载到 DOM
    const terminalElement = document.getElementById('network-xterm')
    if (terminalElement) {
      term.value.open(terminalElement)

      // 自适应大小
      await nextTick()
      fitAddon.value.fit()
      term.value.focus()

      // 绑定输入事件（只绑定一次）
      if (!_inputHandlerBound.value) {
        _inputHandlerBound.value = true
        term.value.onData(data => {
          if (socket.value && socket.value.readyState === WebSocket.OPEN) {
            socket.value.send(data)
          }
        })

        // 监听终端大小变化
        term.value.onResize(({ cols, rows }) => {
          if (socket.value && socket.value.readyState === WebSocket.OPEN) {
            socket.value.send(JSON.stringify({
              type: 'resize',
              cols: cols,
              rows: rows
            }))
          }
        })
      }

      // 显示欢迎信息
      term.value.clear()
      term.value.writeln('\x1b[1;32m═══════════════════════════════════════════════════════════\x1b[0m')
      term.value.writeln('\x1b[1;36m          欢迎使用网络设备 SSH 终端\x1b[0m')
      term.value.writeln('\x1b[1;32m═══════════════════════════════════════════════════════════\x1b[0m')
      term.value.writeln('')
      if (deviceInfo.value) {
        term.value.writeln(`\x1b[1;33m设备名称:\x1b[0m ${deviceInfo.value.deviceName}`)
        term.value.writeln(`\x1b[1;33m设备类型:\x1b[0m ${deviceTypeText.value}`)
        term.value.writeln(`\x1b[1;33m设备地址:\x1b[0m ${deviceInfo.value.ipAddress}`)
      }
      term.value.writeln('')
      term.value.writeln('\x1b[1;34m点击"连接"按钮开始 SSH 会话\x1b[0m')
      term.value.writeln('')
    }
  } catch (error) {
    console.error('Terminal initialization failed:', error)
    ElMessage.error('终端初始化失败: ' + error.message)
  }
}

// 建立 WebSocket 连接
const connectTerminal = async () => {
  if (!deviceInfo.value || isConnecting.value) return

  // 检查SSH配置
  if (!deviceInfo.value.remoteUsername) {
    ElMessage.error('该设备未配置SSH用户名，无法连接。请先编辑设备并填写远程连接信息。')
    term.value?.writeln('\r\n\x1b[31m✖ 该设备未配置SSH信息\x1b[0m')
    term.value?.writeln('\x1b[33m请先在设备管理中配置SSH用户名和密码\x1b[0m\r\n')
    return
  }

  // 如果已连接，先断开
  if (socket.value) {
    disconnectTerminal()
    await new Promise(resolve => setTimeout(resolve, 500))
  }

  isConnecting.value = true

  try {
    const token = storage.getItem('token')
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const wsUrl = `${protocol}//${host}/api/v1/cmdb/snmp-devices/${deviceId.value}/ssh/connect?token=${encodeURIComponent(token)}`

    socket.value = new WebSocket(wsUrl)

    // 连接超时处理
    const timeout = setTimeout(() => {
      if (socket.value.readyState === WebSocket.CONNECTING) {
        socket.value.close()
        term.value.writeln('\r\n\x1b[31m✖ 连接超时\x1b[0m')
        term.value.writeln('\x1b[33m请检查:\x1b[0m')
        term.value.writeln('  1. 后端 WebSocket 服务是否运行')
        term.value.writeln('  2. 设备SSH服务是否正常')
        term.value.writeln('  3. 网络连接是否正常')
        term.value.writeln('  4. SSH用户名密码是否正确')
        isConnecting.value = false
      }
    }, 10000)

    // WebSocket 连接成功
    socket.value.onopen = () => {
      clearTimeout(timeout)
      isConnected.value = true
      isConnecting.value = false
      term.value.clear()
      term.value.writeln('\x1b[32m✔ SSH 连接已建立\x1b[0m\r\n')

      // 发送终端初始化信息
      const { cols, rows } = term.value
      socket.value.send(JSON.stringify({
        type: 'init',
        cols: cols,
        rows: rows,
        term: 'xterm-256color'
      }))

      fitTerminal()
      ElMessage.success(`已连接到 ${deviceInfo.value.deviceName}`)
    }

    // 接收服务器消息
    socket.value.onmessage = (e) => {
      try {
        let data = e.data

        // 尝试解析 JSON（处理 resize 消息）
        try {
          const jsonData = JSON.parse(data)
          if (jsonData && jsonData.resize) {
            return // 忽略 resize 消息
          }
          if (jsonData && jsonData.data) {
            data = jsonData.data
          }
        } catch {
          // 不是 JSON 格式，继续处理原始数据
        }

        term.value.write(data)
      } catch (error) {
        console.error('Terminal write error:', error)
      }
    }

    // 连接关闭
    socket.value.onclose = (e) => {
      clearTimeout(timeout)
      isConnected.value = false
      isConnecting.value = false

      let message = `\r\n\r\n\x1b[31m✖ 连接已关闭\x1b[0m`
      if (e.code !== 1000) {
        message += ` (code: ${e.code}`
        if (e.reason) message += `, reason: ${e.reason}`
        message += ')'
      }
      message += '\r\n'

      setTimeout(() => {
        if (term.value) {
          term.value.write(message)
        }
      }, 200)

      ElMessage.warning('SSH 连接已断开')
    }

    // 连接错误
    socket.value.onerror = (error) => {
      clearTimeout(timeout)
      isConnected.value = false
      isConnecting.value = false

      term.value.writeln('\r\n\x1b[31m✖ WebSocket 连接错误\x1b[0m')
      term.value.writeln(`错误信息: ${error.message || '未知错误'}`)
      term.value.writeln('\x1b[33m请检查网络连接和服务器状态\x1b[0m')

      ElMessage.error('连接失败: ' + (error.message || '未知错误'))
    }
  } catch (error) {
    isConnecting.value = false
    term.value.writeln(`\r\n\x1b[31m✖ 连接失败: ${error.message}\x1b[0m\r\n`)
    ElMessage.error('连接失败: ' + error.message)
  }
}

// 自适应终端大小
const fitTerminal = () => {
  if (fitAddon.value && term.value) {
    try {
      fitAddon.value.fit()
      term.value.focus()
    } catch (error) {
      console.warn('Fit terminal error:', error)
    }
  }
}

// 断开连接
const disconnectTerminal = () => {
  if (socket.value) {
    socket.value.close()
    socket.value = null
  }
  isConnected.value = false
  isConnecting.value = false

  if (term.value) {
    term.value.writeln('\r\n\x1b[33m✔ SSH 连接已手动断开\x1b[0m')
    term.value.writeln('可以点击"连接"按钮重新连接\r\n')
  }
}

// 清屏
const clearTerminal = () => {
  if (term.value) {
    term.value.clear()
    term.value.writeln('\x1b[32m✔ 终端已清屏\x1b[0m\r\n')
  }
}

// 重新连接
const reconnect = () => {
  if (isConnected.value || isConnecting.value) {
    disconnectTerminal()
    setTimeout(() => {
      connect()
    }, 500)
  }
}

// 返回设备列表
const goBack = async () => {
  if (isConnected.value) {
    try {
      await ElMessageBox.confirm('当前SSH连接仍在进行中，确认返回吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      disconnectTerminal()
      router.push('/cmdb/switch')
    } catch {
      // 取消返回
    }
  } else {
    router.push('/cmdb/switch')
  }
}

// 对外暴露的方法
const connect = () => {
  if (!_termInitialized.value) {
    initTerm().then(() => {
      connectTerminal()
    })
  } else {
    connectTerminal()
  }
}

const disconnect = () => {
  disconnectTerminal()
}

const clear = () => {
  clearTerminal()
}

// 窗口大小变化处理
const handleResize = () => {
  fitTerminal()
}

// 获取设备信息
const loadDeviceInfo = async () => {
  try {
    deviceId.value = route.query.deviceId
    if (!deviceId.value) {
      ElMessage.error('缺少设备ID参数')
      router.push('/cmdb/switch')
      return
    }

    const { data: res } = await API.getSnmpDeviceById(deviceId.value)
    if (res.code === 200) {
      // 映射后端字段到前端统一字段
      deviceInfo.value = {
        ...res.data,
        // 根据设备类型映射远程连接字段
        remotePort: res.data.deviceType === 'windows' ? res.data.rdpPort : res.data.sshPort,
        remoteUsername: res.data.deviceType === 'windows' ? res.data.rdpUsername : res.data.sshUsername,
        remotePassword: res.data.deviceType === 'windows' ? res.data.rdpPassword : res.data.sshPassword,
        remoteDomain: res.data.rdpDomain || ''
      }
    } else {
      ElMessage.error(res.message || '获取设备信息失败')
      router.push('/cmdb/switch')
    }
  } catch (error) {
    console.error('加载设备信息失败:', error)
    ElMessage.error('加载设备信息失败: ' + error.message)
    router.push('/cmdb/switch')
  }
}

// 生命周期钩子
onMounted(async () => {
  window.addEventListener('resize', handleResize)
  await loadDeviceInfo()
  if (deviceInfo.value) {
    initTerm().then(() => {
      // 自动连接
      nextTick(() => {
        connect()
      })
    })
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)

  // 清理连接
  if (socket.value) {
    socket.value.close()
    socket.value = null
  }

  // 清理终端实例
  try {
    if (fitAddon.value && !fitAddon.value._disposed) {
      fitAddon.value.dispose()
    }
    if (term.value) {
      term.value.dispose()
    }
  } catch (error) {
    console.warn('Terminal cleanup error:', error)
  }

  isConnected.value = false
  isConnecting.value = false
  _termInitialized.value = false
  _inputHandlerBound.value = false
})

// 暴露方法给父组件
defineExpose({
  connect,
  disconnect,
  clear,
  reconnect,
  fitTerminal
})
</script>

<style scoped>
.terminal-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 100px);
  background: #1e1e1e;
  margin: 20px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

/* ==================== 工具栏样式 ==================== */
.terminal-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #2d2d30;
  border-bottom: 1px solid #3e3e42;
}

.toolbar-left {
  display: flex;
  gap: 6px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
  padding: 4px 12px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.05);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s ease-in-out infinite;
}

.status-connected .status-dot {
  background: #67c23a;
  box-shadow: 0 0 8px #67c23a;
}

.status-connecting .status-dot {
  background: #e6a23c;
  box-shadow: 0 0 8px #e6a23c;
}

.status-disconnected .status-dot {
  background: #909399;
  animation: none;
}

.status-connected {
  color: #67c23a;
}

.status-connecting {
  color: #e6a23c;
}

.status-disconnected {
  color: #909399;
}

.host-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #d4d4d4;
  font-family: 'Consolas', 'Courier New', monospace;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* ==================== 终端区域 ==================== */
.xterm-wrapper {
  flex: 1;
  overflow: hidden;
  padding: 10px;
  background: #1e1e1e;
}

.xterm-terminal {
  height: 100%;
  width: 100%;
}

/* ==================== 按钮样式优化 ==================== */
.terminal-toolbar .el-button {
  min-width: 50px;
  height: 26px;
  font-size: 12px;
  padding: 0 10px;
  border-radius: 4px;
}

.terminal-toolbar .el-button--small {
  padding: 0 10px;
  line-height: 24px;
}

.terminal-toolbar .el-button .el-icon {
  font-size: 12px;
  margin-right: 4px;
}

/* ==================== xterm 样式覆盖 ==================== */
:deep(.xterm) {
  height: 100% !important;
  padding: 0;
}

:deep(.xterm-viewport) {
  background-color: #1e1e1e !important;
  overflow-y: auto !important;
}

:deep(.xterm-screen) {
  width: 100% !important;
}

:deep(.xterm-viewport::-webkit-scrollbar) {
  width: 8px;
}

:deep(.xterm-viewport::-webkit-scrollbar-track) {
  background: #252526;
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb) {
  background: #424245;
  border-radius: 4px;
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb:hover) {
  background: #4e4e4e;
}
</style>
