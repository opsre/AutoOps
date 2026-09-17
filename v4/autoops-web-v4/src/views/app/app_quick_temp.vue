<template>
  <div class="quick-deployment-container">
    <div class="glass-card main-card">
      <div class="deployment-header">
      <div class="header-info">
        <h2 class="deployment-title">{{ deploymentInfo.title || '发布任务' }}</h2>
        <div class="deployment-meta">
          <template v-if="orderedTasks.length > 0">
            <el-tag
              v-for="(task, index) in orderedTasks"
              :key="`env-${index}`"
              :type="getEnvironmentType(task.environment)"
              size="small"
              style="margin-right: 4px;">
              {{ task.app_name }}-{{ task.environment }}
            </el-tag>
          </template>
          <el-tag :type="getStatusType(deploymentInfo.status)" size="small">
            {{ getStatusText(deploymentInfo.status) }}
          </el-tag>
          <span class="meta-text">创建人: {{ deploymentInfo.creator_name || '-' }}</span>
          <span class="meta-text">创建时间: {{ deploymentInfo.created_at ? new Date(deploymentInfo.created_at).toLocaleString() : '-' }}</span>
        </div>
      </div>
      <div class="header-actions">
        <el-dropdown @command="handleExecuteDeployment" trigger="click" :disabled="executing">
          <el-button
            size="small"
            v-authority="['app:quick-release:jobstart']"
            class="modern-btn green-btn"
            :loading="executing">
            <el-icon><VideoPlay /></el-icon>
            启动jenkins任务
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item :command="2" class="execution-mode-item">
                <div class="mode-option">
                  <el-icon><Sort /></el-icon>
                  <div class="mode-info">
                    <div class="mode-title">串行执行</div>
                    <div class="mode-desc">按顺序执行，适合有依赖关系的任务</div>
                  </div>
                  <el-tag size="small" type="success">默认</el-tag>
                </div>
              </el-dropdown-item>
              <el-dropdown-item :command="1" class="execution-mode-item">
                <div class="mode-option">
                  <el-icon><Grid /></el-icon>
                  <div class="mode-info">
                    <div class="mode-title">并行执行</div>
                    <div class="mode-desc">同时执行，适合独立的部署任务</div>
                  </div>
                </div>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button @click="handleRefresh" size="small" class="modern-btn blue-btn">
          <el-icon><Refresh /></el-icon>
          刷新状态
        </el-button>

        <el-button @click="handleGoBack" size="small" class="modern-btn green-btn">
          <el-icon><ArrowLeft /></el-icon>
          返回列表
        </el-button>
      </div>
      </div>

      <div class="deployment-steps-section">
        <div class="deployment-steps-container">
      <svg class="connector-canvas" ref="connectorCanvas"></svg>

      <div
        v-for="(task, index) in orderedTasks"
        :key="task.id"
        class="deployment-step"
        :style="{'--step-order': task.execute_order}">

      <svg class="connector-defs">
        <defs>
          <marker :id="'arrowhead-' + index" markerWidth="6" markerHeight="4"
                  refX="5" refY="2" orient="auto">
            <path d="M0,0 L0,4 L6,2 Z" :fill="getTaskStatusColor(task.status)"/>
          </marker>
          <marker :id="'flow-arrow-' + index" markerWidth="6" markerHeight="4"
                  refX="5" refY="2" orient="auto">
            <path d="M0,0 L0,4 L6,2 Z" :fill="getTaskStatusColor(task.status)"/>
          </marker>
        </defs>
      </svg>

      <div :class="['deployment-card', getTaskClass(task)]">
        <div class="task-status-indicator" :class="'status-' + getTaskStatusName(task.status)">
          <el-icon v-if="task.status === 3" class="status-icon success"><SuccessFilled /></el-icon>
          <el-icon v-else-if="task.status === 4" class="status-icon error"><CircleCloseFilled /></el-icon>
          <el-icon v-else-if="task.status === 2" class="status-icon running"><Loading /></el-icon>
          <el-icon v-else class="status-icon pending"><Clock /></el-icon>
        </div>

        <div class="status-text" :class="'status-' + getTaskStatusName(task.status)">
          {{
            task.status === 3 ? '完成' :
            task.status === 2 ? '运行中' :
            task.status === 4 ? '异常' : '等待'
          }}
        </div>

        <h3 class="card-title">
          <svg class="title-icon" viewBox="0 0 1024 1024" width="20" height="20">
            <path fill="#1A237E" d="M512 1024C229.23 1024 0 794.77 0 512S229.23 0 512 0s512 229.23 512 512-229.23 512-512 512zm0-938.67C276.36 85.33 85.33 276.36 85.33 512S276.36 938.67 512 938.67 938.67 747.64 938.67 512 747.64 85.33 512 85.33z"/>
            <path fill="#1A237E" d="M341.33 426.67a42.67 42.67 0 100-85.34 42.67 42.67 0 000 85.34zM682.67 426.67a42.67 42.67 0 100-85.34 42.67 42.67 0 000 85.34z"/>
            <path fill="#1A237E" d="M512 768c-117.82 0-213.33-95.51-213.33-213.33h426.66C725.33 672.49 629.82 768 512 768z"/>
          </svg>
          {{ task.app_name }}
        </h3>

        <p class="card-content">应用: {{ task.app_code }}</p>
        <p class="card-content">环境: {{ task.environment }}</p>
        <div class="card-duration" v-if="task.duration > 0">
          执行时长: {{ formatDuration(task.duration) }}
        </div>
        <div class="card-buttons">
          <!-- 调试信息 -->
          <div style="font-size: 10px; color: #999; margin-bottom: 5px;">
            任务ID: {{ task.id }} | 状态: {{ task.status }} | 应用: {{ task.app_name }}<br>
            Jenkins URL: {{ task.jenkins_job_url || task.log_url || '无' }}
          </div>

          <el-tooltip effect="dark" content="查看日志" placement="top">
            <button @click.stop="showTaskDetail(task)" class="btn-icon">
              <img src="@/assets/image/日志.svg" width="20" height="20" alt="日志" />
            </button>
          </el-tooltip>
          <el-tooltip effect="dark" content="停止任务" placement="top">
            <button v-authority="['app:quick-release:jobstop']" @click.stop="handleStopTask(task)" class="btn-icon stop-btn">
              <img
                src="@/assets/image/停止.svg"
                width="20"
                height="20"
                alt="停止" />
            </button>
          </el-tooltip>
          <el-tooltip v-if="task.jenkins_job_url || task.log_url" effect="dark" content="Jenkins任务" placement="top">
            <button @click.stop="openJenkinsJob(task.jenkins_job_url || task.log_url)" class="btn-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" stroke="white" stroke-width="2" stroke-linecap="round"/>
                <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" stroke="white" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </button>
          </el-tooltip>
        </div>

        <!-- 错误信息显示 -->
        <div v-if="task.status === 4 && task.error_message" class="error-message">
          <el-alert
            :title="task.error_message"
            type="error"
            :closable="false"
            show-icon>
          </el-alert>
        </div>

        <!-- 构建信息 -->
        <div v-if="task.build_number" class="build-info">
          <el-tag size="small" type="info">
            构建 #{{ task.build_number }}
          </el-tag>
        </div>
      </div>
    </div>
        </div>
      </div>
    </div>

    <!-- 任务详情对话框 -->
    <el-dialog
      v-model="taskDetailVisible"
      :title="currentTask?.app_name || '任务详情'"
      width="600px"
      class="modern-dialog task-detail-dialog">
      <div v-if="currentTask" class="task-detail-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="应用名称">{{ currentTask.app_name }}</el-descriptions-item>
          <el-descriptions-item label="应用编码">{{ currentTask.app_code }}</el-descriptions-item>
          <el-descriptions-item label="环境">{{ currentTask.environment }}</el-descriptions-item>
          <el-descriptions-item label="执行顺序">{{ currentTask.execute_order }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getTaskStatusType(currentTask.status)">{{ getTaskStatusText(currentTask.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="Jenkins任务ID">{{ currentTask.jenkins_env_id || '-' }}</el-descriptions-item>
          <el-descriptions-item label="开始时间">
            {{ currentTask.start_time ? new Date(currentTask.start_time).toLocaleString() : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="结束时间">
            {{ currentTask.end_time ? new Date(currentTask.end_time).toLocaleString() : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="耗时">
            {{ currentTask.duration ? formatDuration(currentTask.duration) : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="构建号">
            {{ currentTask.build_number || '-' }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="currentTask.error_message" class="error-section">
          <h4>错误信息</h4>
          <el-alert
            :title="currentTask.error_message"
            type="error"
            :closable="false"
            show-icon>
          </el-alert>
        </div>

        <!-- 实时日志显示 -->
        <div class="log-section">
          <div class="log-header">
            <h4>实时日志 <small>（任务ID: {{ currentTask.id }}）</small></h4>
            <el-button size="small" @click="fetchTaskLogImmediately(currentTask.id)" type="primary">
              <el-icon><Refresh /></el-icon>
              刷新日志
            </el-button>
          </div>
          <div class="log-container">
            <CodeEditor
              :model-value="getTaskLog(currentTask.id) || '暂无日志数据...'"
              language="bash"
              height="300px"
              :readonly="true"
              fontSize="12px"
            />
          </div>
          <div class="log-info">
            <small>日志长度: {{ getTaskLog(currentTask.id).length }} 字符 |
            实时更新: {{ monitoringIntervals.has(`log_${currentTask.id}`) ? '开启' : '关闭' }}</small>
          </div>
        </div>

        <div class="action-links">
          <el-button
            v-if="currentTask.jenkins_job_url || currentTask.log_url"
            @click="openJenkinsJob(currentTask.jenkins_job_url || currentTask.log_url)"
            type="primary">
            <el-icon><Link /></el-icon>
            打开Jenkins任务
          </el-button>
          <el-button
            v-if="currentTask.log_url"
            @click="openLogUrl(currentTask.log_url)"
            type="success">
            <el-icon><Document /></el-icon>
            查看构建日志
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, nextTick, watch, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Refresh, VideoPlay, ArrowLeft, SuccessFilled, CircleCloseFilled,
  Loading, Clock, Link, Document, View, ArrowDown, Sort, Grid
} from '@element-plus/icons-vue'
import appApi from '@/api/app'
import CodeEditor from '@/components/CodeEditor.vue'

defineOptions({ name: 'AppQuickTemp' })

const route = useRoute()
const router = useRouter()

// 响应式数据
const loading = ref(false)
const executing = ref(false)
const deploymentInfo = ref({})
const tasks = ref([])
const taskDetailVisible = ref(false)
const currentTask = ref(null)
const connectorCanvas = ref(null)

// 获取部署ID
const deploymentId = computed(() => route.params.id)

// 按执行顺序排序的任务
const orderedTasks = computed(() => {
  return [...tasks.value].sort((a, b) => a.execute_order - b.execute_order)
})

// 工具函数
const getEnvironmentType = (env) => {
  const types = { test: 'warning', staging: 'info', prod: 'danger' }
  return types[env] || 'info'
}

const getEnvironmentText = (env) => {
  const texts = { test: '测试', staging: '预发', prod: '生产' }
  return texts[env] || env
}

const getStatusType = (status) => {
  const types = { 1: 'info', 2: 'warning', 3: 'success', 4: 'danger' }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = { 1: '待发布', 2: '发布中', 3: '发布成功', 4: '发布失败' }
  return texts[status] || '未知'
}

const getTaskStatusName = (status) => {
  const names = { 1: 'pending', 2: 'running', 3: 'success', 4: 'error' }
  return names[status] || 'pending'
}

const getTaskStatusText = (status) => {
  const texts = { 1: '等待部署', 2: '部署中', 3: '部署成功', 4: '部署失败' }
  return texts[status] || '未知'
}

const getTaskStatusType = (status) => {
  const types = { 1: 'info', 2: 'warning', 3: 'success', 4: 'danger' }
  return types[status] || 'info'
}

const getTaskStatusColor = (status) => {
  const colors = {
    1: 'rgba(128,128,128,0.7)',
    2: 'rgba(33,150,243,0.7)',
    3: 'rgba(76,175,80,0.7)',
    4: 'rgba(244,67,54,0.7)'
  }
  return colors[status] || 'rgba(128,128,128,0.7)'
}

const getTaskClass = (task) => {
  const statusName = getTaskStatusName(task.status)
  return `task-${statusName}`
}

const formatDuration = (duration) => {
  if (!duration) return '-'
  const minutes = Math.floor(duration / 60)
  const seconds = duration % 60
  return minutes > 0 ? `${minutes}分${seconds}秒` : `${seconds}秒`
}

// 数据加载函数
const fetchDeploymentDetail = async () => {
  loading.value = true
  try {
    console.log('获取部署详情，ID:', deploymentId.value)
    const response = await appApi.getDeploymentDetail(deploymentId.value)
    console.log('API响应:', response)

    if (response.data.code === 200) {
      const data = response.data.data || {}
      console.log('部署数据:', data)

      // 基础部署信息
      deploymentInfo.value = {
        id: data.id,
        title: data.title,
        description: data.description,
        status: data.status,
        creator_name: data.creator_name,
        created_at: data.created_at,
        environment: data.environment // 如果有单一环境字段
      }

      // 处理应用任务数据
      if (data.tasks && Array.isArray(data.tasks)) {
        // 使用 tasks 数据，这里包含正确的 task_id
        tasks.value = data.tasks.map((task, index) => ({
          id: task.id, // 使用正确的任务ID
          app_id: task.app_id,
          app_name: task.app_name || `应用${task.app_id}`,
          app_code: task.app_code || '',
          environment: task.environment,
          execute_order: task.execute_order || index + 1,
          status: task.status || 1,
          start_time: task.start_time,
          end_time: task.end_time,
          duration: task.duration,
          jenkins_job_url: task.jenkins_job_url,
          jenkins_env_id: task.jenkins_env_id,
          build_number: task.build_number,
          log_url: task.log_url,
          error_message: task.error_message
        }))
      } else if (data.applications && Array.isArray(data.applications)) {
        // 备用：多应用场景（但需要获取正确的task_id）
        tasks.value = data.applications.map((app, index) => ({
          id: app.task_id || app.id, // 优先使用 task_id
          app_id: app.app_id,
          app_name: app.app_name || `应用${app.app_id}`,
          app_code: app.app_code || '',
          environment: app.environment,
          execute_order: index + 1,
          status: app.status || 1,
          start_time: app.start_time,
          end_time: app.end_time,
          duration: app.duration,
          jenkins_job_url: app.jenkins_job_url,
          jenkins_env_id: app.jenkins_env_id,
          build_number: app.build_number,
          log_url: app.log_url,
          error_message: app.error_message
        }))
      } else {
        tasks.value = []
      }

      console.log('处理后的任务列表:', tasks.value)

      // 绘制连接线
      await nextTick()
      drawConnectors()

      // 开始监控任务
      startAllTasksMonitoring()
    } else {
      ElMessage.error(response.data.message || '获取部署详情失败')
    }
  } catch (error) {
    console.error('获取部署详情失败:', error)
    ElMessage.error('获取部署详情失败')
  } finally {
    loading.value = false
  }
}

// 绘制步骤间的连接线 - 3列布局
const drawConnectors = () => {
  if (!connectorCanvas.value || orderedTasks.value.length <= 1) return

  const canvas = connectorCanvas.value
  const steps = document.querySelectorAll('.deployment-step')

  // 清空之前的线条
  canvas.innerHTML = ''

  // 获取容器尺寸
  const rect = canvas.getBoundingClientRect()
  canvas.setAttribute('width', rect.width)
  canvas.setAttribute('height', rect.height)

  // 3列布局的连接线逻辑
  for (let i = 0; i < steps.length - 1; i++) {
    const currentStep = steps[i]
    const nextStep = steps[i + 1]

    if (!currentStep || !nextStep) continue

    const currentRect = currentStep.getBoundingClientRect()
    const nextRect = nextStep.getBoundingClientRect()
    const canvasRect = canvas.getBoundingClientRect()

    const task = orderedTasks.value[i]
    const color = getTaskStatusColor(task.status)

    // 计算当前位置在3列中的位置
    const currentCol = i % 3
    const nextCol = (i + 1) % 3
    const currentRow = Math.floor(i / 3)
    const nextRow = Math.floor((i + 1) / 3)

    let startX, startY, endX, endY

    if (currentRow === nextRow) {
      // 同一行：水平连接
      startX = currentRect.right - canvasRect.left
      startY = currentRect.top + currentRect.height / 2 - canvasRect.top
      endX = nextRect.left - canvasRect.left
      endY = nextRect.top + nextRect.height / 2 - canvasRect.top
    } else {
      // 不同行：垂直连接
      if (currentCol === 2) {
        // 从第3列到下一行第1列
        startX = currentRect.left + currentRect.width / 2 - canvasRect.left
        startY = currentRect.bottom - canvasRect.top
        endX = nextRect.left + nextRect.width / 2 - canvasRect.left
        endY = nextRect.top - canvasRect.top
      } else {
        // 普通情况
        startX = currentRect.right - canvasRect.left
        startY = currentRect.top + currentRect.height / 2 - canvasRect.top
        endX = nextRect.left - canvasRect.left
        endY = nextRect.top + nextRect.height / 2 - canvasRect.top
      }
    }

    // 创建路径
    const path = document.createElementNS('http://www.w3.org/2000/svg', 'path')
    let pathData

    if (currentRow === nextRow) {
      // 直线连接
      pathData = `M${startX},${startY} L${endX},${endY}`
    } else {
      // L形连接
      const midX = startX
      const midY = endY
      pathData = `M${startX},${startY} L${midX},${midY} L${endX},${endY}`
    }

    path.setAttribute('d', pathData)
    path.setAttribute('stroke', color)
    path.setAttribute('stroke-width', '2')
    path.setAttribute('fill', 'none')
    path.setAttribute('marker-end', `url(#arrowhead-${i})`)
    path.setAttribute('stroke-dasharray', '5,5')

    canvas.appendChild(path)
  }
}

// 事件处理函数
const handleRefresh = () => {
  fetchDeploymentDetail()
}

const handleExecuteDeployment = async (executionMode = 2) => {
  try {
    const modeText = executionMode === 1 ? '并行执行' : '串行执行'
    await ElMessageBox.confirm(
      `确定要${modeText}发布任务"${deploymentInfo.value.title}"吗？`,
      '确认执行',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    executing.value = true
    const response = await appApi.executeQuickDeployment({
      deployment_id: parseInt(deploymentId.value),
      execution_mode: executionMode
    })
    if (response.data.code === 200) {
      ElMessage.success(`发布任务已启动（${modeText}）`)
      // 刷新数据并开始监控
      setTimeout(() => {
        fetchDeploymentDetail()
      }, 1000)
    } else {
      ElMessage.error(response.data.message || '启动发布失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('执行发布失败:', error)
      ElMessage.error('执行发布失败')
    }
  } finally {
    executing.value = false
  }
}

const handleGoBack = () => {
  router.push('/app/quick-release')
}

const showTaskDetail = (task) => {
  console.log(`显示任务 ${task.id} 的详情和日志`)
  currentTask.value = task
  taskDetailVisible.value = true

  // 总是立即获取日志，无论任务状态如何
  console.log(`任务 ${task.id} 立即获取日志`)
  fetchTaskLogImmediately(task.id)

  // 只有运行中的任务才需要监控
  if (task.status === 2 && !monitoringIntervals.has(`log_${task.id}`)) {
    console.log(`任务 ${task.id} 正在运行，开始监控`)
    startTaskMonitoring(task.id)
  }
}

// 立即获取任务日志（不依赖轮询）
const fetchTaskLogImmediately = async (taskId) => {
  try {
    console.log(`立即获取任务 ${taskId} 的完整日志`)
    const response = await appApi.getTaskLog(taskId, 0) // 从头开始获取
    if (response.data.code === 200) {
      const logData = response.data.data
      if (logData.log) {
        console.log(`立即获取到任务 ${taskId} 日志，长度: ${logData.log.length}`)
        const cleanedLog = cleanAnsiEscapeSequences(logData.log)
        taskLogs.value.set(taskId, cleanedLog)
        // 更新日志起始位置
        if (logData.text_size) {
          logStarts.value.set(taskId, logData.text_size)
        }
      }
    }
  } catch (error) {
    console.error(`立即获取任务 ${taskId} 日志失败:`, error)
  }
}

const openJenkinsJob = (url) => {
  if (url) {
    window.open(url, '_blank')
  }
}

const openLogUrl = (url) => {
  if (url) {
    window.open(url, '_blank')
  }
}

const handleStopTask = async (task) => {
  // 如果任务已完成或失败，不允许停止
  if (task.status === 3 || task.status === 4) {
    ElMessage.warning('任务已完成或失败，无法停止')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要停止任务"${task.app_name}"吗？`,
      '确认停止',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // 调用停止任务的API
    const response = await appApi.stopTask(task.id)
    if (response.data.code === 200) {
      ElMessage.success('任务停止成功')
      // 停止对该任务的监控
      stopTaskMonitoring(task.id)
      // 刷新数据
      fetchDeploymentDetail()
    } else {
      ElMessage.error(response.data.message || '停止任务失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('停止任务失败:', error)
      ElMessage.error('停止任务失败')
    }
  }
}

// === 实时监控相关函数 ===

// 监控间隔Map
const monitoringIntervals = new Map()
const taskLogs = ref(new Map()) // 存储每个任务的日志
const logStarts = ref(new Map()) // 存储每个任务的日志起始位置

// 开始监控所有任务
const startAllTasksMonitoring = () => {
  if (orderedTasks.value.length === 0) return

  console.log('开始监控所有任务:', orderedTasks.value.map(t => ({ id: t.id, name: t.app_name })))

  orderedTasks.value.forEach(task => {
    startTaskMonitoring(task.id)
  })
}

// 开始监控单个任务
const startTaskMonitoring = (taskId) => {
  console.log(`开始监控任务 ${taskId}`)

  // 检查任务是否已经完成
  const task = tasks.value.find(t => t.id === taskId)
  if (task && (task.status === 3 || task.status === 4)) {
    console.log(`任务 ${taskId} 已完成，状态: ${task.status}，跳过监控`)
    return
  }

  // 如果已经在监控，先停止
  stopTaskMonitoring(taskId)

  // 初始化日志起始位置
  if (!logStarts.value.has(taskId)) {
    logStarts.value.set(taskId, 0)
  }

  // 状态监控
  const statusInterval = setInterval(async () => {
    try {
      console.log(`正在检查任务 ${taskId} 的状态...`)
      const response = await appApi.getTaskStatus(taskId)
      console.log(`任务 ${taskId} 状态响应:`, response.data)

      if (response.data.code === 200) {
        updateTaskStatus(taskId, response.data.data)

        // 如果任务完成（成功或失败），停止状态监控
        if (response.data.data.status === 3 || response.data.data.status === 4) {
          console.log(`任务 ${taskId} 已完成，状态: ${response.data.data.status_text}，停止状态监控`)
          const statusKey = `status_${taskId}`
          if (monitoringIntervals.has(statusKey)) {
            clearInterval(monitoringIntervals.get(statusKey))
            monitoringIntervals.delete(statusKey)
            console.log(`已停止任务 ${taskId} 的状态监控`)
          }
        }
      } else {
        console.error(`任务 ${taskId} 状态检查失败:`, response.data)
      }
    } catch (error) {
      console.error(`获取任务 ${taskId} 状态失败:`, error)
      console.error('错误详情:', error.response?.data || error.message)
    }
  }, 3000) // 每3秒检查状态

  // 日志监控 - 持续运行，不因为任务完成而停止
  const logInterval = setInterval(async () => {
    try {
      const start = logStarts.value.get(taskId) || 0
      console.log(`正在获取任务 ${taskId} 的日志，从位置 ${start} 开始...`)

      const response = await appApi.getTaskLog(taskId, start)
      console.log(`任务 ${taskId} 日志响应:`, response.data)

      if (response.data.code === 200) {
        const logData = response.data.data
        if (logData.log) {
          console.log(`任务 ${taskId} 获取到新日志，长度: ${logData.log.length}`)
          appendTaskLog(taskId, logData.log)
        }
        if (logData.has_more && logData.text_size) {
          logStarts.value.set(taskId, logData.text_size)
          console.log(`任务 ${taskId} 更新日志起始位置为: ${logData.text_size}`)
        } else if (!logData.has_more) {
          console.log(`任务 ${taskId} 日志已完整，停止日志监控`)
          const logKey = `log_${taskId}`
          if (monitoringIntervals.has(logKey)) {
            clearInterval(monitoringIntervals.get(logKey))
            monitoringIntervals.delete(logKey)
            console.log(`已停止任务 ${taskId} 的日志监控`)
          }
        }
      } else {
        console.error(`任务 ${taskId} 日志获取失败:`, response.data)
      }
    } catch (error) {
      console.error(`获取任务 ${taskId} 日志失败:`, error)
      console.error('日志错误详情:', error.response?.data || error.message)
    }
  }, 2000) // 每2秒检查日志

  // 保存间隔ID
  monitoringIntervals.set(`status_${taskId}`, statusInterval)
  monitoringIntervals.set(`log_${taskId}`, logInterval)
}

// 停止监控单个任务
const stopTaskMonitoring = (taskId) => {
  console.log(`停止监控任务 ${taskId}`)

  const statusKey = `status_${taskId}`
  const logKey = `log_${taskId}`

  if (monitoringIntervals.has(statusKey)) {
    clearInterval(monitoringIntervals.get(statusKey))
    monitoringIntervals.delete(statusKey)
  }

  if (monitoringIntervals.has(logKey)) {
    clearInterval(monitoringIntervals.get(logKey))
    monitoringIntervals.delete(logKey)
  }
}

// 停止所有监控
const stopAllMonitoring = () => {
  console.log('停止所有监控')
  monitoringIntervals.forEach((interval, key) => {
    clearInterval(interval)
  })
  monitoringIntervals.clear()
}

// 更新任务状态
const updateTaskStatus = (taskId, statusData) => {
  const taskIndex = tasks.value.findIndex(t => t.id === taskId)
  if (taskIndex !== -1) {
    // 更新任务状态，保持原有的jenkins_job_url不变
    tasks.value[taskIndex] = {
      ...tasks.value[taskIndex],
      status: statusData.status,
      build_number: statusData.build_number,
      start_time: statusData.start_time,
      end_time: statusData.end_time,
      duration: statusData.duration,
      error_message: statusData.error_message,
      // 保持原有的jenkins_job_url，如果statusData有jenkins_job_url则使用新的
      jenkins_job_url: statusData.jenkins_job_url || tasks.value[taskIndex].jenkins_job_url
    }

    console.log(`任务 ${taskId} 状态更新:`, statusData)
  }
}

// 清理ANSI转义序列
const cleanAnsiEscapeSequences = (text) => {
  // 移除ANSI转义序列
  return text.replace(/\u001b\[[0-9;]*[mGKH]/g, '')
}

// 添加任务日志
const appendTaskLog = (taskId, newLog) => {
  if (!taskLogs.value.has(taskId)) {
    taskLogs.value.set(taskId, '')
  }
  const currentLog = taskLogs.value.get(taskId)
  // 清理ANSI转义序列后再添加
  const cleanedLog = cleanAnsiEscapeSequences(newLog)
  taskLogs.value.set(taskId, currentLog + cleanedLog)

  // 如果当前正在查看该任务的日志，自动滚动到底部
  nextTick(() => {
    if (currentTask.value && currentTask.value.id === taskId && taskDetailVisible.value) {
      const logContainer = document.querySelector('.log-container')
      if (logContainer) {
        logContainer.scrollTop = logContainer.scrollHeight
      }
    }
  })
}

// 获取任务日志
const getTaskLog = (taskId) => {
  const log = taskLogs.value.get(taskId) || ''
  console.log(`获取任务 ${taskId} 的日志，长度: ${log.length}`)
  return log
}

// 生命周期
onMounted(() => {
  fetchDeploymentDetail()
})

// 监听数据变化，启动轮询
let pollingInterval = null

watch(() => deploymentInfo.value.status, (newStatus) => {
  // 清除之前的轮询
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }

  // 如果是发布中状态，开启轮询
  if (newStatus === 2) {
    pollingInterval = setInterval(() => {
      fetchDeploymentDetail().then(() => {
        // 如果状态不再是发布中，停止轮询
        if (deploymentInfo.value.status !== 2) {
          clearInterval(pollingInterval)
          pollingInterval = null
        }
      })
    }, 5000) // 每5秒刷新一次
  }
}, { immediate: true })

// 组件卸载时清除定时器
onUnmounted(() => {
  if (pollingInterval) {
    clearInterval(pollingInterval)
  }
  // 停止所有监控
  stopAllMonitoring()
})

// 监听窗口大小变化，重新绘制连接线
window.addEventListener('resize', () => {
  setTimeout(drawConnectors, 100)
})
</script>

<style scoped>
.quick-deployment-container {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
  position: relative;
  display: flex;
  flex-direction: column;
}

.glass-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.main-card {
  padding: 24px;
}

.deployment-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 2px solid #e3e8f0;
}

.header-info {
  flex: 1;
}

.deployment-title {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 12px 0;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.deployment-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.meta-text {
  color: #606266;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.connector-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 10;
}

/* 流程容器 - 类似 Ansible */
.deployment-steps-section {
  margin-top: 24px;
}

.deployment-steps-container {
  display: flex;
  flex-wrap: wrap;
  gap: 90px;
  padding: 50px;
  position: relative;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
  min-height: 400px;
  justify-content: flex-start;
  align-items: flex-start;
}

.deployment-step {
  width: calc(33.33% - 60px);
  flex: 0 0 calc(33.33% - 60px);
  order: var(--step-order);
  position: relative;
  z-index: 5;
  min-width: 280px;
}

.connector-defs {
  position: absolute;
  width: 0;
  height: 0;
  pointer-events: none;
}

.deployment-card {
  position: relative;
  padding: 16px 16px 12px 16px;
  width: 100%;
  height: 180px;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  color: #333;
  border: 1px solid rgba(255, 255, 255, 0.3);
  transform: translateZ(0);
  will-change: transform, box-shadow;
  overflow: hidden;
  z-index: 5;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  box-sizing: border-box;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

/* 任务状态样式 - 参考 Ansible */
.deployment-card.task-pending {
  background: linear-gradient(to bottom right, #FF9800 0%, #FFE0B2 100%) !important;
}

.deployment-card.task-running {
  background: linear-gradient(to bottom right, #2196F3 0%, #BBDEFB 100%) !important;
}

.deployment-card.task-success {
  background: linear-gradient(to bottom right, #4CAF50 0%, #C8E6C9 100%) !important;
}

.deployment-card.task-error {
  background: linear-gradient(to bottom right, #F44336 0%, #FFCDD2 100%) !important;
}

.deployment-card:hover {
  transform: translateY(-8px) scale(1.03);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  z-index: 10;
}

/* 状态文本样式 */
.status-text {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 10px;
  padding: 4px 8px;
  border-radius: 12px;
  white-space: nowrap;
  display: inline-block;
  min-width: auto;
  max-width: 60px;
  text-align: center;
  font-weight: 500;
}

.status-text.status-completed {
  background: #f0f9eb;
  color: #67c23a;
}

.status-text.status-running {
  background: #ecf5ff;
  color: #409eff;
}

.status-text.status-pending {
  background: #fdf6ec;
  color: #e6a23c;
}

.status-text.status-error {
  background: #fef0f0;
  color: #f56c6c;
}

/* 卡片标题样式 */
.card-title {
  display: flex;
  align-items: center;
  margin: 0 0 6px;
  font-size: 14px;
  color: #1A237E;
  font-weight: 500;
}

.title-icon {
  margin-right: 6px;
  width: 16px;
  height: 16px;
}

/* 卡片内容样式 */
.card-content {
  margin: 0 0 4px;
  font-size: 12px;
  color: #333;
  font-weight: 500;
}

.card-duration {
  margin: 6px 0;
  font-size: 11px;
  color: #666;
}

/* 按钮样式 - 类似 Ansible */
.card-buttons {
  display: flex;
  gap: 6px;
  justify-content: flex-end;
  align-items: center;
  min-height: 32px;
  position: absolute;
  bottom: 10px;
  right: 16px;
}

.btn-icon {
  background: linear-gradient(45deg, #409eff, #66b3ff);
  border: none;
  border-radius: 6px;
  padding: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
}

.btn-icon svg {
  width: 16px;
  height: 16px;
}

.btn-icon:hover {
  background: linear-gradient(45deg, #3a8ee6, #5ca7f0);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 任务状态指示器 */
.task-status-indicator {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  z-index: 10;
}

/* 错误和构建信息样式 */
.error-message {
  margin-top: 10px;
}

.build-info {
  margin-top: 8px;
}

/* 任务详情对话框样式 */
.task-detail-dialog {
  border-radius: 12px;
}

.task-detail-content {
  padding: 10px 0;
}

.error-section {
  margin-top: 20px;
}

.error-section h4 {
  margin-bottom: 10px;
  color: #f56c6c;
}

.action-links {
  margin-top: 20px;
  display: flex;
  gap: 12px;
}

.deployment-card.task-error {
  border-left: 4px solid #f56c6c;
}

.deployment-card.task-running {
  border-left: 4px solid #409eff;
  animation: pulse 2s infinite;
}

.deployment-card.task-pending {
  border-left: 4px solid #909399;
}

@keyframes pulse {
  0% { box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1); }
  50% { box-shadow: 0 8px 32px rgba(64, 158, 255, 0.3); }
  100% { box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1); }
}

.status-icon {
  font-size: 14px;
}

.status-icon.success {
  color: #67c23a;
}

.status-icon.error {
  color: #f56c6c;
}

.status-icon.running {
  color: #409eff;
  animation: spin 1s linear infinite;
}

.status-icon.pending {
  color: #909399;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.status-text {
  position: absolute;
  top: 8px;
  left: 20px;
  font-size: 12px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.9);
}

.status-text.status-success {
  color: #67c23a;
}

.status-text.status-error {
  color: #f56c6c;
}

.status-text.status-running {
  color: #409eff;
}

.status-text.status-pending {
  color: #909399;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  margin: 20px 0 12px 0;
  color: #303133;
}

.title-icon {
  flex-shrink: 0;
}

.card-content {
  margin-bottom: 16px;
}

.task-info {
  margin: 6px 0;
  font-size: 14px;
  color: #606266;
}

.task-timing {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.card-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.error-message {
  margin-top: 12px;
}

.build-info {
  margin-top: 8px;
}

.task-detail-dialog .task-detail-content {
  padding: 0;
}

.error-section {
  margin-top: 20px;
}

.error-section h4 {
  margin: 0 0 12px 0;
  color: #f56c6c;
  font-size: 16px;
}

.log-section {
  margin-top: 20px;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.log-section h4 {
  margin: 0;
  color: #303133;
  font-size: 16px;
}

.log-container {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #fafafa;
}

.log-content {
  margin: 0;
  padding: 12px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.4;
  color: #333;
  white-space: pre-wrap;
  word-break: break-all;
  background: transparent;
}

.log-info {
  padding: 8px 12px;
  background: #f5f5f5;
  border-top: 1px solid #e4e7ed;
  font-size: 11px;
  color: #666;
}

.action-links {
  margin-top: 20px;
  display: flex;
  gap: 12px;
}

/* 现代对话框样式 */
.modern-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  padding: 20px 24px;
}

.modern-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.modern-btn {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.modern-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.green-btn {
  background: linear-gradient(45deg, #67c23a, #85ce61);
  border: none;
  color: white;
}

.green-btn:hover {
  background: linear-gradient(45deg, #5daf34, #73c050);
}

.blue-btn {
  background: linear-gradient(45deg, #409eff, #66b3ff);
  border: none;
  color: white;
}

.blue-btn:hover {
  background: linear-gradient(45deg, #3a8ee6, #5ca7f0);
}

.stop-btn {
  background: linear-gradient(45deg, #409eff, #66b3ff);
}

.stop-btn:hover:not(:disabled) {
  background: linear-gradient(45deg, #3a8ee6, #5ca7f0);
}

.stop-btn:disabled {
  background: linear-gradient(45deg, #ccc, #ddd);
  cursor: not-allowed;
}

.stop-btn:disabled:hover {
  transform: none;
  box-shadow: none;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .deployment-header {
    flex-direction: column;
    gap: 16px;
  }

  .header-actions {
    width: 100%;
    justify-content: center;
  }

  .deployment-card {
    width: 100%;
    max-width: 350px;
  }

  .deployment-step:nth-child(even) {
    justify-content: flex-start;
  }
}

/* 执行模式选择样式 */
.execution-mode-item {
  padding: 0 !important;
}

.mode-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  width: 100%;
}

.mode-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mode-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.mode-desc {
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
}

.mode-option .el-icon {
  font-size: 18px;
  color: #409eff;
}

.execution-mode-item:hover .mode-option {
  background-color: #f5f7fa;
}
</style>
