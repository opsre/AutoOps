<template>
  <div class="rdp-container" ref="rdpContainer">
    <!-- RDP 工具栏 -->
    <div class="rdp-toolbar">
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
          <el-tag type="primary" size="small">Windows RDP (Guacamole)</el-tag>
          {{ deviceInfo.deviceName }} ({{ deviceInfo.ipAddress }})
        </span>
      </div>
    </div>

    <!-- Guacamole RDP 显示区域 -->
    <div class="rdp-display-wrapper">
      <!-- 欢迎屏幕 -->
      <div v-if="!isConnected && !isConnecting" class="welcome-screen">
        <div class="welcome-content">
          <el-icon class="welcome-icon" :size="80">
            <Monitor />
          </el-icon>
          <h2>欢迎使用 RDP 远程桌面</h2>
          <div v-if="deviceInfo" class="device-info">
            <p><strong>设备名称:</strong> {{ deviceInfo.deviceName }}</p>
            <p><strong>设备类型:</strong> Windows 服务器</p>
            <p><strong>设备地址:</strong> {{ deviceInfo.ipAddress }}</p>
            <p><strong>RDP端口:</strong> {{ deviceInfo.remotePort || 3389 }}</p>
          </div>
          <p class="welcome-tip">点击"连接"按钮开始 RDP 远程桌面会话</p>
        </div>
      </div>

      <!-- 连接中屏幕 -->
      <div v-if="isConnecting" class="connecting-screen">
        <el-icon class="loading-icon is-loading" :size="60">
          <Loading />
        </el-icon>
        <p>正在连接到远程桌面...</p>
        <p class="connect-info">{{ deviceInfo?.deviceName }} ({{ deviceInfo?.ipAddress }})</p>
      </div>

      <!-- Guacamole 显示容器 -->
      <div
        ref="guacDisplay"
        id="guac-display"
        class="guac-display"
        :class="{ 'display-hidden': !isConnected }"
      ></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Connection, SwitchButton, Close, Back, Monitor, Loading } from '@element-plus/icons-vue'
import storage from '@/utils/storage'
import API from '@/api/cmdb'
import * as GuacamoleLib from 'guacamole-common-js'

// 使用导入的 Guacamole 对象
const Guacamole = GuacamoleLib.default || GuacamoleLib

// 路由
const route = useRoute()
const router = useRouter()

// 状态管理
const rdpContainer = ref(null)
const guacDisplay = ref(null)
const isConnecting = ref(false)
const isConnected = ref(false)

// 设备信息
const deviceInfo = ref(null)
const deviceId = ref(null)

// Guacamole 客户端
const guacClient = ref(null)
const guacTunnel = ref(null)
const guacMouse = ref(null)
const guacKeyboard = ref(null)

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

// 建立 Guacamole RDP 连接
const connectRDP = async () => {
  if (!deviceInfo.value || isConnecting.value) return

  // 调试：检查 Guacamole 对象
  console.log('Guacamole 对象:', Guacamole)
  console.log('Guacamole.Client:', Guacamole.Client)
  console.log('Guacamole.Tunnel:', Guacamole.Tunnel)
  console.log('Guacamole 状态常量:', {
    'Client.State.IDLE': Guacamole.Client?.State?.IDLE,
    'Client.State.CONNECTING': Guacamole.Client?.State?.CONNECTING,
    'Client.State.CONNECTED': Guacamole.Client?.State?.CONNECTED,
    'Client.State.DISCONNECTED': Guacamole.Client?.State?.DISCONNECTED,
    'Tunnel.State.CONNECTING': Guacamole.Tunnel?.State?.CONNECTING,
    'Tunnel.State.OPEN': Guacamole.Tunnel?.State?.OPEN,
    'Tunnel.State.CLOSED': Guacamole.Tunnel?.State?.CLOSED
  })

  // 检查RDP配置 - 兼容 remoteUsername 和 rdpUsername 两种字段名
  const username = deviceInfo.value.rdpUsername || deviceInfo.value.remoteUsername
  if (!username) {
    ElMessage.error('该设备未配置RDP用户名，无法连接。请先编辑设备并填写远程连接信息。')
    return
  }

  if (deviceInfo.value.deviceType !== 'windows') {
    ElMessage.error('该设备不是Windows服务器，无法使用RDP连接')
    return
  }

  // 如果已连接，先断开
  if (guacClient.value) {
    disconnectRDP()
    await new Promise(resolve => setTimeout(resolve, 500))
  }

  isConnecting.value = true

  try {
    const token = storage.getItem('token')

    // ✅ 直接连接后端 8000 端口（绕过 Vue DevServer 代理）
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsHost = window.location.hostname + ':8000'
    const wsUrl = `${protocol}//${wsHost}/api/v1/cmdb/snmp-devices/${deviceId.value}/guac-rdp/connect?token=${encodeURIComponent(token)}`

    console.log('==========================================')
    console.log('🔌 正在连接 RDP WebSocket:', wsUrl)
    console.log('设备:', deviceInfo.value.ipAddress)
    console.log('后端主动握手模式：后端会主动发送 Guacamole 协议')
    console.log('==========================================')

    // ✅ 使用原生 WebSocket（不使用 Guacamole.WebSocketTunnel）
    const ws = new WebSocket(wsUrl)

    // 创建 Guacamole.Client（用于协议解析和渲染）
    // 使用自定义 Tunnel 来桥接原生 WebSocket 和 Guacamole.Client
    const customTunnel = new Guacamole.Tunnel()
    guacClient.value = new Guacamole.Client(customTunnel)

    // WebSocket 打开
    ws.onopen = () => {
      console.log('✅ WebSocket 连接成功，等待后端主动握手...')
      customTunnel.setState(Guacamole.Tunnel.State.OPEN)
    }

    // WebSocket 接收消息 → 转发给 Guacamole.Client
    ws.onmessage = (event) => {
      console.log('📩 收到后端数据:', event.data.substring(0, 100) + '...')

      try {
        // 将数据传递给 Guacamole.Client 处理
        if (customTunnel.oninstruction) {
          // 解析 Guacamole 指令格式: "opcode,arg1,arg2,...;"
          const instructions = event.data.split(';')
          instructions.forEach(inst => {
            if (!inst) return
            const parts = inst.split(',')
            if (parts.length > 0) {
              const opcode = parts[0].split('.').pop() // 去掉长度前缀
              const args = parts.slice(1).map(arg => {
                // 去掉长度前缀，如 "4.1024" → "1024"
                return arg.includes('.') ? arg.split('.').pop() : arg
              })
              console.log('📨 指令:', opcode, args)
              customTunnel.oninstruction(opcode, args)
            }
          })
        }
      } catch (err) {
        console.error('解析 Guacamole 指令失败:', err)
      }
    }

    // WebSocket 错误
    ws.onerror = (error) => {
      console.error('❌ WebSocket 错误:', error)
      customTunnel.setState(Guacamole.Tunnel.State.CLOSED)
      isConnecting.value = false
      isConnected.value = false
      ElMessage.error('RDP连接失败，请检查网络连接')
    }

    // WebSocket 关闭
    ws.onclose = (event) => {
      console.log('🔌 WebSocket 已关闭')
      console.log('Code:', event.code, 'Reason:', event.reason)
      customTunnel.setState(Guacamole.Tunnel.State.CLOSED)
      isConnecting.value = false
      isConnected.value = false
      if (event.code !== 1000) {
        ElMessage.warning('RDP连接已断开')
      }
    }

    // 保存 WebSocket 引用
    guacTunnel.value = { socket: ws, customTunnel }

    // 设置 Guacamole.Client 事件监听
    guacClient.value.onstatechange = (state) => {
      console.log('Guacamole 客户端状态:', state, {
        IDLE: Guacamole.Client.State.IDLE,
        CONNECTING: Guacamole.Client.State.CONNECTING,
        WAITING: Guacamole.Client.State.WAITING,
        CONNECTED: Guacamole.Client.State.CONNECTED,
        DISCONNECTING: Guacamole.Client.State.DISCONNECTING,
        DISCONNECTED: Guacamole.Client.State.DISCONNECTED
      })

      if (state === Guacamole.Client.State.CONNECTED) {
        isConnected.value = true
        isConnecting.value = false
        ElMessage.success(`已连接到 ${deviceInfo.value.deviceName}`)
        setupMouseEvents()
        setupKeyboardEvents()
      } else if (state === Guacamole.Client.State.DISCONNECTED) {
        isConnected.value = false
        isConnecting.value = false
      }
    }

    guacClient.value.onerror = (error) => {
      console.error('Guacamole 客户端错误:', error)
      isConnecting.value = false
      isConnected.value = false
      ElMessage.error('RDP连接失败: ' + (error.message || '未知错误'))
    }

    // 挂载显示
    guacDisplay.value.innerHTML = ''
    guacDisplay.value.appendChild(guacClient.value.getDisplay().getElement())

    // ✅ 不调用 connect()，等待后端主动发送数据
    console.log('✅ 等待后端主动发送 Guacamole 协议握手...')

    // 超时处理
    setTimeout(() => {
      if (isConnecting.value) {
        disconnectRDP()
        ElMessage.error('RDP连接超时，请检查设备配置和网络连接')
      }
    }, 15000)

  } catch (error) {
    isConnecting.value = false
    console.error('建立RDP连接失败:', error)
    ElMessage.error('RDP连接失败: ' + error.message)
  }
}

// 设置鼠标事件
const setupMouseEvents = () => {
  if (!guacClient.value || !guacDisplay.value) return

  const display = guacClient.value.getDisplay()
  guacMouse.value = new Guacamole.Mouse(display.getElement())

  guacMouse.value.onmousedown =
  guacMouse.value.onmouseup =
  guacMouse.value.onmousemove = (mouseState) => {
    if (guacClient.value && isConnected.value) {
      guacClient.value.sendMouseState(mouseState)
    }
  }
}

// 设置键盘事件
const setupKeyboardEvents = () => {
  if (!guacClient.value) return

  guacKeyboard.value = new Guacamole.Keyboard(document)

  guacKeyboard.value.onkeydown = (keysym) => {
    if (guacClient.value && isConnected.value) {
      guacClient.value.sendKeyEvent(1, keysym)
    }
  }

  guacKeyboard.value.onkeyup = (keysym) => {
    if (guacClient.value && isConnected.value) {
      guacClient.value.sendKeyEvent(0, keysym)
    }
  }
}

// 断开 RDP 连接
const disconnectRDP = () => {
  console.log('正在断开 RDP 连接...')

  if (guacClient.value) {
    try {
      guacClient.value.disconnect()
    } catch (error) {
      console.warn('断开Guacamole客户端时出错:', error)
    }
    guacClient.value = null
  }

  if (guacTunnel.value) {
    try {
      // 关闭原生 WebSocket
      if (guacTunnel.value.socket) {
        guacTunnel.value.socket.close()
      }
      // 关闭自定义 Tunnel
      if (guacTunnel.value.customTunnel) {
        guacTunnel.value.customTunnel.setState(Guacamole.Tunnel.State.CLOSED)
      }
    } catch (error) {
      console.warn('断开 WebSocket 时出错:', error)
    }
    guacTunnel.value = null
  }

  if (guacMouse.value) {
    guacMouse.value = null
  }

  if (guacKeyboard.value) {
    guacKeyboard.value = null
  }

  isConnected.value = false
  isConnecting.value = false
}

// 重新连接
const reconnect = () => {
  if (isConnected.value || isConnecting.value) {
    disconnectRDP()
    setTimeout(() => {
      connect()
    }, 500)
  }
}

// 返回设备列表
const goBack = async () => {
  if (isConnected.value) {
    try {
      await ElMessageBox.confirm('当前RDP连接仍在进行中，确认返回吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      disconnectRDP()
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
  connectRDP()
}

const disconnect = () => {
  disconnectRDP()
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
      deviceInfo.value = res.data

      // 调试：打印设备信息
      console.log('设备信息:', JSON.stringify(deviceInfo.value, null, 2))
      console.log('rdpUsername:', deviceInfo.value.rdpUsername)
      console.log('remoteUsername:', deviceInfo.value.remoteUsername)

      // 验证是否为Windows设备
      if (deviceInfo.value.deviceType !== 'windows') {
        ElMessage.error('该设备不是Windows服务器，无法使用RDP连接')
        router.push('/cmdb/switch')
        return
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
  await loadDeviceInfo()
})

onBeforeUnmount(() => {
  // 清理 RDP 连接
  disconnectRDP()
})

// 暴露方法给父组件
defineExpose({
  connect,
  disconnect,
  reconnect
})
</script>

<style scoped>
.rdp-container {
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
.rdp-toolbar {
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

/* ==================== RDP 显示区域 ==================== */
.rdp-display-wrapper {
  flex: 1;
  overflow: auto;
  background: #000000;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

/* ==================== Guacamole 显示容器 ==================== */
.guac-display {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #000000;
}

.display-hidden {
  display: none;
}

/* ==================== 欢迎屏幕 ==================== */
.welcome-screen {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background: #1e1e1e;
}

.welcome-content {
  text-align: center;
  color: #d4d4d4;
  padding: 40px;
}

.welcome-icon {
  color: #409eff;
  margin-bottom: 20px;
}

.welcome-content h2 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 30px;
  color: #409eff;
}

.device-info {
  margin: 30px 0;
  padding: 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  text-align: left;
  max-width: 400px;
  margin-left: auto;
  margin-right: auto;
}

.device-info p {
  margin: 10px 0;
  font-size: 14px;
  line-height: 1.8;
}

.device-info strong {
  color: #409eff;
  margin-right: 10px;
}

.welcome-tip {
  font-size: 14px;
  color: #909399;
  margin-top: 20px;
}

/* ==================== 连接中屏幕 ==================== */
.connecting-screen {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: #e6a23c;
  background: #1e1e1e;
}

.loading-icon {
  margin-bottom: 20px;
}

.connecting-screen p {
  font-size: 16px;
  font-weight: 500;
  margin: 5px 0;
}

.connect-info {
  font-size: 14px;
  color: #909399;
  margin-top: 10px;
}

/* ==================== 按钮样式优化 ==================== */
.rdp-toolbar .el-button {
  min-width: 50px;
  height: 26px;
  font-size: 12px;
  padding: 0 10px;
  border-radius: 4px;
}

.rdp-toolbar .el-button--small {
  padding: 0 10px;
  line-height: 24px;
}

.rdp-toolbar .el-button .el-icon {
  font-size: 12px;
  margin-right: 4px;
}

/* ==================== 响应式布局 ==================== */
@media (max-width: 768px) {
  .rdp-container {
    margin: 10px;
  }

  .welcome-content h2 {
    font-size: 20px;
  }

  .device-info {
    padding: 15px;
  }
}

</style>
