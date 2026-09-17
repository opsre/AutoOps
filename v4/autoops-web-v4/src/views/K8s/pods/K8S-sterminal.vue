<template>
  <div class="terminal-container">
    <el-card>
      <template #header>
        <div class="terminal-header">
          <div class="terminal-info">
            <el-icon class="terminal-icon">
              <Monitor />
            </el-icon>
            <span>终端 - {{ displayPodName }}</span>
            <el-tag v-if="currentContainer" type="info" size="small">
              {{ currentContainerLabel }}
            </el-tag>
            <el-tag v-if="currentShell" type="success" size="small">
              {{ currentShell }}
            </el-tag>
            <el-tag v-if="isStandaloneDebugSession" type="warning" size="small">
              调试 Pod
            </el-tag>
          </div>
          <div class="terminal-controls">
            <el-select
              v-model="selectedContainer"
              placeholder="选择容器"
              size="small"
              style="width: 150px;"
              @change="handleContainerChange"
            >
              <el-option
                v-for="container in containers"
                :key="container.value"
                :label="container.label"
                :value="container.value"
              />
            </el-select>
            <el-select
              v-model="selectedShell"
              placeholder="选择Shell"
              size="small"
              style="width: 110px;"
              @change="handleShellChange"
            >
              <el-option
                v-for="shell in shellOptions"
                :key="shell.value"
                :label="shell.label"
                :value="shell.value"
              />
            </el-select>
            <el-button
              type="success"
              size="small"
              :icon="FolderOpened"
              @click="openFileManager"
            >
              文件管理
            </el-button>
            <el-button
              type="primary"
              size="small"
              :loading="connecting"
              @click="connect"
              :disabled="!selectedContainer"
            >
              {{ isConnected ? '重连' : '连接' }}
            </el-button>
            <el-button
              type="warning"
              size="small"
              :loading="creatingDebugContainer"
              :disabled="!canCreateDebugContainer"
              :title="debugActionTitle"
              @click="handleCreateDebugContainer"
            >
              一键调试
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="disconnect"
              :disabled="!isConnected"
            >
              断开
            </el-button>
          </div>
        </div>
      </template>
      
      <div class="terminal-content">
        <div ref="terminalElement" class="xterm-container" @click="focusTerminal"></div>
        <div v-if="!isConnected" class="terminal-placeholder">
          <div class="placeholder-content">
            <el-icon class="placeholder-icon">
              <Monitor />
            </el-icon>
            <div class="placeholder-text">
              <h3>Web 终端</h3>
              <p>选择容器并点击连接按钮开始使用终端</p>
              <p v-if="!selectedContainer" class="tip">请先选择一个容器</p>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <PodFileManagerDialog
      v-model:visible="fileManagerVisible"
      :cluster-id="clusterId"
      :namespace-name="namespace"
      :pod-name="displayPodName"
      :container-name="activeFileContainer"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Monitor, FolderOpened } from '@element-plus/icons-vue'
import { Terminal as XTerm } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import k8sApi from '@/api/k8s'
import storage from '@/utils/storage'
import PodFileManagerDialog from './PodFileManagerDialog.vue'

const route = useRoute()
const DEFAULT_DEBUG_IMAGE = 'crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/busybox:latest'

// 从路由参数获取Pod信息
const clusterId = computed(() => route.params.clusterId)
const namespace = computed(() => route.params.namespace)
const podName = computed(() => route.params.podName)
const containerFromQuery = computed(() => route.query.container)
const shellFromQuery = computed(() => route.query.shell)
const commandFromQuery = computed(() => route.query.command)

// 终端状态
const connecting = ref(false)
const isConnected = ref(false)
const selectedContainer = ref('')
const selectedShell = ref('bash')
const currentContainer = ref('')
const currentShell = ref('')
const activePodName = ref('')
const currentDebugMode = ref('')
const containers = ref([])
const creatingDebugContainer = ref(false)
const fileManagerVisible = ref(false)
const shellOptions = [
  { label: 'bash', value: 'bash' },
  { label: 'sh', value: 'sh' }
]

const activeFileContainer = computed(() => selectedContainer.value || currentContainer.value || '')
const displayPodName = computed(() => activePodName.value || podName.value)
const isStandaloneDebugSession = computed(() => currentDebugMode.value === 'standalone-pod')
const selectedContainerOption = computed(() => {
  return containers.value.find(container => container.value === selectedContainer.value) || null
})
const currentContainerLabel = computed(() => {
  const current = containers.value.find(container => container.value === currentContainer.value)
  return current?.label || currentContainer.value
})
const canCreateDebugContainer = computed(() => {
  return selectedContainerOption.value?.kind === 'normal' && !isStandaloneDebugSession.value
})
const debugActionTitle = computed(() => {
  if (!selectedContainer.value) return '请先选择业务容器'
  if (isStandaloneDebugSession.value) return '当前已切换到独立调试 Pod，不支持继续一键调试'
  if (selectedContainerOption.value?.kind === 'ephemeral') return '当前已是临时调试容器，请切换到业务容器后再创建'
  if (selectedContainerOption.value?.kind === 'init') return '初始化容器不支持创建调试容器'
  return '为当前业务容器注入临时调试容器并自动连接终端'
})

// 终端相关
const terminalElement = ref(null)
let terminal = null
let fitAddon = null
let websocket = null
let heartbeatInterval = null

const normalizeQueryValue = (value) => {
  if (Array.isArray(value)) return value[0] || ''
  return value || ''
}

const normalizeContainerValue = (value) => String(value || '').replace(/\s+\((init|ephemeral)\)$/i, '')

const getContainerKind = (name) => {
  const label = String(name || '')
  if (/\(ephemeral\)$/i.test(label)) return 'ephemeral'
  if (/\(init\)$/i.test(label)) return 'init'
  return 'normal'
}

const createContainerOption = (container) => {
  const label = container?.name || container
  return {
    label,
    value: normalizeContainerValue(label),
    kind: getContainerKind(label)
  }
}

const sanitizeDebugContainerName = (name) => {
  const normalized = String(name || '')
    .trim()
    .toLowerCase()
    .replaceAll('_', '-')
    .replace(/[^a-z0-9-]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 63)
    .replace(/-+$/g, '')

  return normalized || 'debugger'
}

const resolveTokenValue = (token) => {
  if (!token) return ''
  return typeof token === 'object' ? token.access_token || token.token || '' : token
}

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms))

const setSelectedContainerValue = (preferredContainer = '') => {
  const allowRouteQuery = activePodName.value === '' || activePodName.value === podName.value
  const preferredValue = normalizeContainerValue(preferredContainer)
  if (preferredValue && containers.value.some(container => container.value === preferredValue)) {
    selectedContainer.value = preferredValue
    return
  }

  const queryContainer = normalizeContainerValue(normalizeQueryValue(containerFromQuery.value))
  if (allowRouteQuery && queryContainer && containers.value.some(container => container.value === queryContainer)) {
    selectedContainer.value = queryContainer
    return
  }

  if (selectedContainer.value && containers.value.some(container => container.value === selectedContainer.value)) {
    return
  }

  selectedContainer.value = containers.value[0]?.value || ''
}

const resolveInitialShell = () => {
  if (shellFromQuery.value === 'sh') {
    return 'sh'
  }

  if (typeof commandFromQuery.value === 'string' && commandFromQuery.value.trim() === '/bin/sh') {
    return 'sh'
  }

  return 'bash'
}

selectedShell.value = resolveInitialShell()
activePodName.value = podName.value

// 初始化终端
const initTerminal = () => {
  if (!terminalElement.value) return
  
  terminal = new XTerm({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
    theme: {
      background: '#060101',
      foreground: '#ffffff',
      cursor: '#ffffff',
      selection: '#ffffff80'
    },
    // 确保终端可以接收输入
    disableStdin: false,
    convertEol: true,
    // 其他输入相关配置
    allowTransparency: false,
    rightClickSelectsWord: true,
    scrollback: 1000,
    cols: 80,
    rows: 24,
    windowsMode: false,
    lineHeight: 1.0,
    letterSpacing: 0,
    tabStopWidth: 8,
    bellStyle: 'none',
    drawBoldTextInBrightColors: true,
    cursorStyle: 'underline'
  })
  
  fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)
  
  terminal.open(terminalElement.value)
  fitAddon.fit()
  
  // 设置终端焦点，确保可以接收输入
  setTimeout(() => {
    terminal.focus()
  }, 100)
  
  // 监听窗口大小变化 - 发送正确的JSON格式resize消息
  const resizeObserver = new ResizeObserver(() => {
    if (fitAddon) {
      fitAddon.fit()
      if (websocket && websocket.readyState === WebSocket.OPEN) {
        const resizeMessage = {
          operation: 'resize',
          data: {
            cols: terminal.cols,
            rows: terminal.rows
          }
        }
        websocket.send(JSON.stringify(resizeMessage))
      }
    }
  })
  resizeObserver.observe(terminalElement.value)
  
  // 监听键盘输入 - 发送JSON格式的消息
  terminal.onData((data) => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      const inputMessage = {
        operation: 'stdin',
        data: data
      }
      websocket.send(JSON.stringify(inputMessage))
    }
  })
}

// 获取容器列表
const getContainers = async (targetPodName = displayPodName.value, preferredContainer = '') => {
  try {
    const response = await k8sApi.getPodContainers(clusterId.value, namespace.value, targetPodName)
    const responseData = response.data || response
    
    if (responseData.code === 200 && responseData.data) {
      containers.value = responseData.data.map(createContainerOption)
    } else {
      console.warn('获取容器列表响应异常:', responseData)
      // 如果API失败，尝试从Pod详情获取
      await getContainersFromPodDetail(targetPodName, preferredContainer)
    }
    
    setSelectedContainerValue(preferredContainer)
  } catch (error) {
    console.error('获取容器列表失败:', error)
    ElMessage.error('获取容器列表失败')
    // 降级处理：从Pod详情获取容器信息
    await getContainersFromPodDetail(targetPodName, preferredContainer)
  }
}

// 从工作负载获取容器列表（降级方案）
const getContainersFromPodDetail = async (targetPodName = podName.value, preferredContainer = '') => {
  try {
    if (targetPodName !== podName.value) {
      containers.value = [createContainerOption(preferredContainer || 'main')]
      setSelectedContainerValue(preferredContainer || 'main')
      return
    }

    // 从Pod名称推断Deployment名称
    const podNameParts = podName.value.split('-')
    let deploymentName = null
    if (podNameParts.length >= 3) {
      deploymentName = podNameParts.slice(0, -2).join('-')
      console.log('🔍 从Pod名称推断Deployment:', deploymentName)
    }

    if (deploymentName) {
      const response = await k8sApi.getWorkloadPods(
        clusterId.value,
        namespace.value,
        'deployment',
        deploymentName
      )
      const responseData = response.data || response

      if (responseData.code === 200 && responseData.data) {
        // 从Pod列表中找到当前Pod
        const currentPod = responseData.data.find(pod => pod.name === podName.value)
        if (currentPod) {
          if (currentPod.spec && currentPod.spec.containers) {
            containers.value = currentPod.spec.containers.map(container => createContainerOption(container.name))
          } else if (currentPod.containers) {
            containers.value = currentPod.containers.map(container => createContainerOption(container.name))
          }
          setSelectedContainerValue()
          return
        }
      }
    }

    // 如果无法通过工作负载API获取，使用默认容器名
    containers.value = [createContainerOption('main')]
    setSelectedContainerValue()
  } catch (error) {
    console.error('从工作负载获取容器列表失败:', error)
    // 最后的降级：使用默认容器名
    containers.value = [createContainerOption('main')]
    setSelectedContainerValue()
  }
}

const buildTerminalWsUrl = (token, options = {}) => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'

  if (options.terminalPath) {
    const url = new URL(options.terminalPath, window.location.origin)
    url.searchParams.set('token', token)
    return `${protocol}//${url.host || window.location.host}${url.pathname}?${url.searchParams.toString()}`
  }

  const targetPodName = options.podName || displayPodName.value
  const targetContainerName = normalizeContainerValue(options.containerName || selectedContainer.value)
  const targetShell = options.shell || selectedShell.value
  const params = new URLSearchParams()

  params.set('containerName', targetContainerName)

  if (targetShell === 'sh') {
    if (typeof commandFromQuery.value === 'string' && commandFromQuery.value.trim() === '/bin/sh' && !shellFromQuery.value) {
      params.set('command', '/bin/sh')
    } else {
      params.set('shell', 'sh')
    }
  } else {
    params.set('shell', 'bash')
  }

  params.set('token', token)

  return `${protocol}//${window.location.host}/api/v1/k8s/cluster/${clusterId.value}/namespaces/${namespace.value}/pods/${encodeURIComponent(targetPodName)}/terminal?${params.toString()}`
}

// 连接终端
const connect = async (options = {}) => {
  const targetPodName = options.podName || displayPodName.value
  const targetContainerName = normalizeContainerValue(options.containerName || selectedContainer.value)
  const targetShell = options.shell || selectedShell.value

  if (!targetContainerName) {
    ElMessage.error('请选择容器')
    return
  }

  if (websocket) {
    websocket.close()
  }

  connecting.value = true
  
  try {
    const token = storage.getItem('token')
    const tokenValue = resolveTokenValue(token)
    if (!tokenValue) {
      ElMessage.error('未找到token，请先登录')
      connecting.value = false
      return
    }

    const wsUrl = buildTerminalWsUrl(tokenValue, {
      podName: targetPodName,
      containerName: targetContainerName,
      shell: targetShell,
      terminalPath: options.terminalPath
    })
    
    console.log('🔌 WebSocket连接URL:', wsUrl)
    console.log('🔑 使用token:', tokenValue ? '***' + tokenValue.slice(-4) : 'null')
    
    websocket = new WebSocket(wsUrl)
    
    // 添加连接超时处理
    const connectionTimeout = setTimeout(() => {
      if (websocket.readyState === WebSocket.CONNECTING) {
        console.warn('⏰ WebSocket连接超时，主动关闭')
        websocket.close()
        connecting.value = false
        ElMessage.error('连接超时，请检查网络或服务器状态')
      }
    }, 10000) // 10秒超时
    
    websocket.onopen = () => {
      clearTimeout(connectionTimeout)
      connecting.value = false
      isConnected.value = true
      activePodName.value = targetPodName
      currentContainer.value = targetContainerName
      currentShell.value = targetShell
      ElMessage.success('终端连接成功')
      
      // 清空终端并显示欢迎信息
      if (terminal) {
        terminal.clear()
        terminal.writeln('\x1B[1;32m🖥️  K8S Pod 终端连接成功\x1B[0m')
        terminal.writeln('\x1B[1;34m📦 Pod: ' + targetPodName + '\x1B[0m')
        terminal.writeln('\x1B[1;33m🐳 容器: ' + targetContainerName + '\x1B[0m')
        terminal.writeln('\x1B[1;35m🐚 Shell: ' + targetShell + '\x1B[0m')
        terminal.writeln('\x1B[1;36m⚡ 等待服务器终端初始化...\x1B[0m')
        terminal.writeln('')
        
        // 发送初始终端大小 - 使用正确的JSON格式
        const resizeMessage = {
          operation: 'resize',
          data: {
            cols: terminal.cols,
            rows: terminal.rows
          }
        }
        console.log('📐 连接成功，发送初始终端大小:', resizeMessage)
        websocket.send(JSON.stringify(resizeMessage))
        
        // 暂时禁用心跳检测，避免干扰
        // startHeartbeat()
        console.log('🚫 已禁用心跳检测，避免干扰K8S终端')
        
        // 确保终端获得焦点
        setTimeout(() => {
          terminal.focus()
          console.log('🖥️ WebSocket连接成功，终端重新获得焦点')
        }, 200)
      }
    }
    
    websocket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        console.log('📨 收到消息:', message.operation)
        
        if ((message.operation === 'stdout' || message.operation === 'stderr') && message.data) {
          console.log('✍️ 写入终端')
          terminal.write(message.data)
        }
      } catch (error) {
        // 非JSON消息直接写入
        console.log('📝 非JSON消息，直接写入')
        terminal.write(event.data)
      }
    }
    
    websocket.onerror = (error) => {
      console.error('❌ WebSocket连接错误:', error)
      console.error('🔍 错误详情:', {
        type: error.type,
        target: error.target,
        readyState: error.target?.readyState,
        url: error.target?.url
      })
      ElMessage.error('终端连接出错，请检查网络或后端服务')
      connecting.value = false
      isConnected.value = false
    }
    
    websocket.onclose = (event) => {
      console.log('🔌 WebSocket连接关闭:', {
        code: event.code,
        reason: event.reason,
        wasClean: event.wasClean
      })
      isConnected.value = false
      connecting.value = false
      currentShell.value = ''
      
      // 停止心跳检测
      stopHeartbeat()
      
      // 提供更友好的错误信息
      const getCloseReason = (code, reason) => {
        switch (code) {
          case 1000:
            return '正常关闭'
          case 1001:
            return '服务器关闭或重启'
          case 1002:
            return '协议错误'
          case 1003:
            return '不支持的数据类型'
          case 1005:
            return '连接异常关闭'
          case 1006:
            return '连接异常断开，可能是网络问题'
          case 1007:
            return '数据格式错误'
          case 1008:
            return '违反协议策略'
          case 1009:
            return '数据过大'
          case 1010:
            return '缺少扩展支持'
          case 1011:
            return '服务器内部错误'
          case 1015:
            return 'TLS握手失败'
          default:
            return reason || '连接中断'
        }
      }
      
      if (event.code === 1000) {
        ElMessage.info('终端连接已正常断开')
      } else if (event.code === 1006) {
        ElMessage.error('终端连接异常断开，请检查网络连接')
      } else {
        const reasonText = getCloseReason(event.code, event.reason)
        ElMessage.warning(`终端连接关闭: ${reasonText}`)
      }
    }
    
  } catch (error) {
    console.error('❌ 创建WebSocket失败:', error)
    console.error('🔍 创建失败详情:', {
      name: error.name,
      message: error.message,
      stack: error.stack
    })
    ElMessage.error(`连接终端失败: ${error.message}`)
    connecting.value = false
    isConnected.value = false
  }
}

const stopHeartbeat = () => {
  if (heartbeatInterval) {
    clearInterval(heartbeatInterval)
    heartbeatInterval = null
    console.log('💔 心跳检测已停止')
  }
}

// 断开连接
const disconnect = () => {
  stopHeartbeat()
  if (websocket) {
    websocket.close()
  }
  currentShell.value = ''
  if (terminal) {
    terminal.clear()
  }
}

// 切换容器
const handleContainerChange = () => {
  if (isConnected.value) {
    ElMessage.info('切换容器需要重新连接')
    disconnect()
  }
}

const handleShellChange = () => {
  if (isConnected.value) {
    ElMessage.info('切换 Shell 需要重新连接')
    disconnect()
  }
}

const openFileManager = () => {
  fileManagerVisible.value = true
}

const waitForContainerToAppear = async (targetPodName, containerName, maxAttempts = 8, intervalMs = 1000) => {
  const normalizedName = normalizeContainerValue(containerName)
  if (!normalizedName) return false

  for (let attempt = 0; attempt < maxAttempts; attempt += 1) {
    try {
      const response = await k8sApi.getPodContainers(clusterId.value, namespace.value, targetPodName)
      const responseData = response.data || response

      if (responseData.code === 200 && Array.isArray(responseData.data)) {
        containers.value = responseData.data.map(createContainerOption)
      }
    } catch (error) {
      console.error('轮询调试容器列表失败:', error)
    }

    if (containers.value.some(container => container.value === normalizedName)) {
      selectedContainer.value = normalizedName
      return true
    }

    if (attempt < maxAttempts - 1) {
      await sleep(intervalMs)
    }
  }

  return false
}

const handleCreateDebugContainer = async () => {
  if (!canCreateDebugContainer.value) {
    ElMessage.warning('请先选择业务容器')
    return
  }

  creatingDebugContainer.value = true

  try {
    const response = await k8sApi.createPodDebugContainer(
      clusterId.value,
      namespace.value,
      podName.value,
      {
        targetContainerName: selectedContainer.value,
        containerName: sanitizeDebugContainerName(`debugger-${selectedContainer.value}-test`),
        debugImage: DEFAULT_DEBUG_IMAGE,
        privileged: true
      }
    )

    const responseData = response.data || response
    if (responseData.code !== 200 || !responseData.data) {
      ElMessage.error(responseData.message || '启动调试容器失败')
      return
    }

    const debugMode = responseData.data.mode || ''
    const targetPodName = responseData.data.podName || responseData.data.debugPodName || podName.value
    const debugContainerName = normalizeContainerValue(responseData.data.debugContainerName)
    if (!debugContainerName) {
      ElMessage.error('调试容器名称为空，无法连接')
      return
    }

    if (isConnected.value) {
      disconnect()
      await sleep(300)
    }

    currentDebugMode.value = debugMode
    activePodName.value = targetPodName
    containers.value = [createContainerOption(debugContainerName)]
    selectedContainer.value = debugContainerName
    selectedShell.value = 'bash'

    if (debugMode === 'standalone-pod' && responseData.data.terminalPath) {
      ElMessage.success(responseData.data.message || '独立调试 Pod 已就绪，正在连接')
      await connect({
        podName: targetPodName,
        containerName: debugContainerName,
        shell: 'bash',
        terminalPath: responseData.data.terminalPath
      })
      return
    }

    const found = await waitForContainerToAppear(targetPodName, debugContainerName)
    if (!found) {
      ElMessage.warning('调试容器已创建，但容器列表尚未刷新完成，请稍后手动选择后连接')
      return
    }

    ElMessage.success(responseData.data.message || '调试容器已创建，正在连接')
    await sleep(800)
    await connect({
      podName: targetPodName,
      containerName: debugContainerName,
      shell: 'bash'
    })
  } catch (error) {
    console.error('启动调试容器失败:', error)
    ElMessage.error(error.message || '启动调试容器失败')
  } finally {
    creatingDebugContainer.value = false
  }
}

// 手动设置终端焦点
const focusTerminal = () => {
  if (terminal) {
    terminal.focus()
    console.log('🖥️ 手动设置终端焦点')
  }
}

onMounted(() => {
  initTerminal()
  getContainers(displayPodName.value)
})

onUnmounted(() => {
  disconnect()
  stopHeartbeat()
  if (terminal) {
    terminal.dispose()
  }
})
</script>

<style scoped>
.terminal-container {
  height: 100vh;
  padding: 20px;
  background: #f5f5f5;
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.terminal-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.terminal-icon {
  color: #409eff;
}

.terminal-controls {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

.terminal-content {
  height: calc(100vh - 200px);
  position: relative;
}

.xterm-container {
  height: 100%;
  width: 100%;
  background-color: #060101;
  padding: 10px;
  border-radius: 4px;
  overflow: hidden;
}

.terminal-placeholder {
  position: absolute;
  inset: 0;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(6, 1, 1, 0.9);
  border: 2px dashed rgba(255, 255, 255, 0.12);
  border-radius: 4px;
}

.placeholder-content {
  text-align: center;
  color: #c0c4cc;
}

.placeholder-icon {
  font-size: 48px;
  margin-bottom: 16px;
  color: #409eff;
}

.placeholder-text h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: #ffffff;
}

.placeholder-text p {
  margin: 4px 0;
  color: #c0c4cc;
}

.placeholder-text .tip {
  color: #f56c6c;
  font-weight: bold;
}
</style>
