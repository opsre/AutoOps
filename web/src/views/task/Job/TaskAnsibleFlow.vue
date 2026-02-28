<template>
  <el-dialog
    title="Ansible任务流程"
    v-model="visible"
    width="90%"
    top="5vh"
    :before-close="handleClose"
  >
    <div class="ansible-flow-container" v-if="taskData">
      <!-- 任务基本信息 -->
      <el-card class="task-info-card" shadow="never">
        <div class="task-header">
          <h3>{{ taskData.name }}</h3>
          <div class="task-meta">
            <el-tag :type="taskData.type === 1 ? 'primary' : 'success'">
              {{ taskData.type === 1 ? '手动任务' : '自动任务' }}
            </el-tag>
            <el-tag :type="getStatusType(taskData.status)" style="margin-left: 10px">
              {{ getStatusText(taskData.status) }}
            </el-tag>
          </div>
        </div>
        
        <!-- 操作按钮 -->
        <div class="task-actions">
          <el-button
            type="primary"
            icon="VideoPlay"
            @click="startTask"
            :loading="starting"
            :disabled="taskData.status === 2"
          >
            启动任务
          </el-button>
          <el-button
            type="info"
            icon="Refresh"
            @click="refreshTaskStatus"
            :loading="refreshing"
          >
            刷新状态
          </el-button>
          <el-button
            type="warning"
            icon="Document"
            @click="viewLogs"
          >
            查看日志
          </el-button>
        </div>
      </el-card>

      <!-- 主机分组信息 -->
      <el-card class="host-groups-card" shadow="never" v-if="hostGroups">
        <template #header>
          <span>主机分组配置</span>
        </template>
        <div class="host-groups">
          <div
            v-for="(hosts, groupName) in hostGroups"
            :key="groupName"
            class="host-group"
          >
            <div class="group-header">
              <el-icon><Server /></el-icon>
              <span class="group-name">{{ groupName }}</span>
              <el-tag size="small">{{ hosts.length }}台主机</el-tag>
            </div>
            <div class="hosts-list">
              <el-tag
                v-for="hostId in hosts"
                :key="hostId"
                size="small"
                type="info"
                style="margin: 2px;"
              >
                主机ID: {{ hostId }}
              </el-tag>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 任务流程图 -->
      <el-card class="flow-chart-card" shadow="never">
        <template #header>
          <span>任务执行流程</span>
        </template>
        
        <div class="flow-container">
          <!-- 动态生成任务步骤 -->
          <div
            v-for="(step, index) in taskSteps"
            :key="index"
            class="flow-step"
            :class="getStepClass(step.status)"
          >
            <div class="step-content">
              <div class="step-number">{{ index + 1 }}</div>
              <div class="step-info">
                <h4>{{ step.name }}</h4>
                <p>{{ step.description }}</p>
                <div class="step-meta">
                  <span class="step-status">{{ getStepStatusText(step.status) }}</span>
                  <span class="step-time" v-if="step.startTime">
                    {{ formatTime(step.startTime) }}
                  </span>
                </div>
              </div>
              
              <!-- 步骤操作按钮 -->
              <div class="step-actions">
                <el-tooltip content="查看日志" placement="top">
                  <el-button
                    size="small"
                    icon="Document"
                    circle
                    @click="viewStepLog(step)"
                    :disabled="step.status === 1"
                  />
                </el-tooltip>
                <el-tooltip content="重新执行" placement="top">
                  <el-button
                    size="small"
                    icon="Refresh"
                    circle
                    type="warning"
                    @click="retryStep(step)"
                    :disabled="step.status === 1 || step.status === 2"
                  />
                </el-tooltip>
              </div>
            </div>
            
            <!-- 连接线 -->
            <div v-if="index < taskSteps.length - 1" class="flow-connector">
              <div class="connector-line" :class="{ 'active': step.status === 3 }"></div>
              <div class="connector-arrow" :class="{ 'active': step.status === 3 }"></div>
            </div>
          </div>
        </div>
      </el-card>

      <!-- Git仓库信息 (自动任务) -->
      <el-card 
        class="git-info-card" 
        shadow="never" 
        v-if="taskData.type === 2 && taskData.gitRepo"
      >
        <template #header>
          <span>Git仓库信息</span>
        </template>
        <div class="git-info">
          <el-icon><FolderOpened /></el-icon>
          <span>{{ taskData.gitRepo }}</span>
        </div>
        
      </el-card>

      <!-- 全局变量 (手动任务) -->
      <el-card 
        class="variables-card" 
        shadow="never" 
        v-if="taskData.type === 1 && taskData.variables"
      >
        <template #header>
          <span>全局变量</span>
        </template>
        <pre class="variables-content">{{ formatVariables(taskData.variables) }}</pre>
      </el-card>
    </div>

    <!-- 日志查看对话框 -->
    <el-dialog
      title="任务日志"
      v-model="logDialogVisible"
      width="80%"
      append-to-body
    >
      <div class="log-container">
        <pre class="log-content">{{ currentLogs }}</pre>
      </div>
      <template #footer>
        <el-button @click="logDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="downloadLogs">下载日志</el-button>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Server, FolderOpened } from '@element-plus/icons-vue'
import {
  GetAnsibleTaskById,
  StartAnsibleTask,
  GetAnsibleTaskLog
} from '@/api/task'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  taskId: {
    type: [String, Number],
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const taskData = ref(null)
const starting = ref(false)
const refreshing = ref(false)
const logDialogVisible = ref(false)
const currentLogs = ref('')

// 解析主机分组
const hostGroups = computed(() => {
  if (!taskData.value?.hostGroups) return null
  try {
    return typeof taskData.value.hostGroups === 'string' 
      ? JSON.parse(taskData.value.hostGroups) 
      : taskData.value.hostGroups
  } catch {
    return null
  }
})

// 动态生成任务步骤
const taskSteps = computed(() => {
  if (!taskData.value) return []
  
  const steps = []
  
  // 基础步骤
  steps.push({
    name: '环境准备',
    description: '检查主机连接状态和环境配置',
    status: getStepStatus(1),
    startTime: taskData.value.createdAt
  })
  
  if (taskData.value.type === 2) {
    // 自动任务 - Git仓库步骤
    steps.push({
      name: 'Git代码拉取',
      description: `从 ${taskData.value.gitRepo} 拉取最新代码`,
      status: getStepStatus(2),
      startTime: taskData.value.startTime
    })
  } else {
    // 手动任务 - 文件上传步骤
    steps.push({
      name: '文件上传',
      description: '上传Playbook和Roles文件到执行环境',
      status: getStepStatus(2),
      startTime: taskData.value.startTime
    })
  }
  
  // 根据主机分组数量动态生成执行步骤
  if (hostGroups.value) {
    Object.keys(hostGroups.value).forEach((groupName, index) => {
      steps.push({
        name: `执行 ${groupName} 分组`,
        description: `在 ${groupName} 分组的 ${hostGroups.value[groupName].length} 台主机上执行Ansible任务`,
        status: getStepStatus(3 + index),
        startTime: taskData.value.startTime,
        groupName: groupName
      })
    })
  }
  
  // 最终步骤
  steps.push({
    name: '任务完成',
    description: '收集执行结果并生成报告',
    status: getStepStatus(steps.length + 1),
    startTime: taskData.value.endTime
  })
  
  return steps
})

// 根据任务状态计算步骤状态
const getStepStatus = (stepIndex) => {
  if (!taskData.value) return 1
  
  const taskStatus = taskData.value.status
  const totalSteps = Object.keys(hostGroups.value || {}).length + 3
  
  if (taskStatus === 1) return 1 // 等待中
  if (taskStatus === 4) return stepIndex <= 2 ? 4 : 1 // 失败
  if (taskStatus === 2) {
    // 运行中 - 根据进度计算
    const progress = Math.floor(totalSteps * 0.6) // 假设当前进度
    return stepIndex <= progress ? 3 : stepIndex === progress + 1 ? 2 : 1
  }
  if (taskStatus === 3) return 3 // 成功
  
  return 1
}

const fetchTaskData = async () => {
  if (!props.taskId) return
  
  try {
    const response = await GetAnsibleTaskById(props.taskId)
    if (response?.data?.code === 200) {
      taskData.value = response.data.data
    }
  } catch (error) {
    console.error('获取任务详情失败:', error)
    ElMessage.error('获取任务详情失败')
  }
}

const startTask = async () => {
  try {
    starting.value = true
    const response = await StartAnsibleTask(props.taskId)
    if (response?.data?.code === 200) {
      ElMessage.success('任务启动成功')
      await fetchTaskData()
    }
  } catch (error) {
    console.error('启动任务失败:', error)
    ElMessage.error('启动任务失败')
  } finally {
    starting.value = false
  }
}

const refreshTaskStatus = async () => {
  try {
    refreshing.value = true
    await fetchTaskData()
    ElMessage.success('状态已刷新')
  } catch (error) {
    console.error('刷新状态失败:', error)
    ElMessage.error('刷新状态失败')
  } finally {
    refreshing.value = false
  }
}

const viewLogs = async () => {
  try {
    // 这里需要根据实际API调整参数
    const response = await GetAnsibleTaskLog(props.taskId, 'main')
    if (response?.data) {
      currentLogs.value = response.data
      logDialogVisible.value = true
    }
  } catch (error) {
    console.error('获取日志失败:', error)
    ElMessage.error('获取日志失败')
  }
}

const viewStepLog = async (step) => {
  try {
    // 根据步骤获取对应日志
    const workId = step.groupName || 'main'
    const response = await GetAnsibleTaskLog(props.taskId, workId)
    if (response?.data) {
      currentLogs.value = response.data
      logDialogVisible.value = true
    }
  } catch (error) {
    console.error('获取步骤日志失败:', error)
    ElMessage.error('获取步骤日志失败')
  }
}

const retryStep = (step) => {
  ElMessage.info('重新执行功能暂未实现')
}

const downloadLogs = () => {
  const blob = new Blob([currentLogs.value], { type: 'text/plain' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `ansible-task-${props.taskId}-logs.txt`
  a.click()
  window.URL.revokeObjectURL(url)
}

const handleClose = () => {
  visible.value = false
}

const getStatusType = (status) => {
  switch(status) {
    case 1: return 'info'
    case 2: return 'warning'
    case 3: return 'success'
    case 4: return 'danger'
    default: return 'info'
  }
}

const getStatusText = (status) => {
  switch(status) {
    case 1: return '等待中'
    case 2: return '运行中'
    case 3: return '成功'
    case 4: return '失败'
    default: return '未知'
  }
}

const getStepClass = (status) => {
  switch(status) {
    case 1: return 'step-pending'
    case 2: return 'step-running'
    case 3: return 'step-success'
    case 4: return 'step-error'
    default: return 'step-pending'
  }
}

const getStepStatusText = (status) => {
  return getStatusText(status)
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).replace(/\//g, '-')
}

const formatVariables = (variables) => {
  try {
    const parsed = typeof variables === 'string' ? JSON.parse(variables) : variables
    return JSON.stringify(parsed, null, 2)
  } catch {
    return variables
  }
}

// 监听taskId变化
watch(() => props.taskId, (newId) => {
  if (newId && visible.value) {
    fetchTaskData()
  }
})

// 监听dialog打开
watch(() => visible.value, (newVisible) => {
  if (newVisible && props.taskId) {
    fetchTaskData()
  }
})
</script>

<style scoped>
.ansible-flow-container {
  max-height: 70vh;
  overflow-y: auto;
}

.task-info-card {
  margin-bottom: 20px;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.task-header h3 {
  margin: 0;
  color: #303133;
}

.task-meta {
  display: flex;
  align-items: center;
}

.task-actions {
  display: flex;
  gap: 10px;
}

.host-groups-card,
.flow-chart-card,
.git-info-card,
.variables-card {
  margin-bottom: 20px;
}

.host-groups {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.host-group {
  border: 1px solid #ebeef5;
  border-radius: 6px;
  padding: 15px;
  background: #fafafa;
}

.group-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
  font-weight: 500;
  color: #303133;
}

.group-name {
  font-size: 16px;
}

.hosts-list {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.flow-container {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.flow-step {
  position: relative;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 10px;
  transition: all 0.3s ease;
}

.step-pending {
  background: #f4f4f5;
  border-left: 4px solid #909399;
}

.step-running {
  background: #fdf6ec;
  border-left: 4px solid #e6a23c;
  animation: pulse 2s infinite;
}

.step-success {
  background: #f0f9eb;
  border-left: 4px solid #67c23a;
}

.step-error {
  background: #fef0f0;
  border-left: 4px solid #f56c6c;
}

.step-content {
  display: flex;
  align-items: flex-start;
  gap: 15px;
}

.step-number {
  flex-shrink: 0;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: #409eff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 14px;
}

.step-success .step-number {
  background: #67c23a;
}

.step-error .step-number {
  background: #f56c6c;
}

.step-running .step-number {
  background: #e6a23c;
}

.step-pending .step-number {
  background: #909399;
}

.step-info {
  flex: 1;
}

.step-info h4 {
  margin: 0 0 5px 0;
  color: #303133;
  font-size: 16px;
}

.step-info p {
  margin: 0 0 10px 0;
  color: #606266;
  font-size: 14px;
}

.step-meta {
  display: flex;
  gap: 15px;
  font-size: 12px;
  color: #909399;
}

.step-status {
  font-weight: 500;
}

.step-actions {
  display: flex;
  gap: 5px;
}

.flow-connector {
  position: relative;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 5px 0;
}

.connector-line {
  width: 2px;
  height: 15px;
  background: #dcdfe6;
  transition: background 0.3s ease;
}

.connector-line.active {
  background: #67c23a;
}

.connector-arrow {
  position: absolute;
  bottom: 0;
  width: 0;
  height: 0;
  border-left: 4px solid transparent;
  border-right: 4px solid transparent;
  border-top: 6px solid #dcdfe6;
  transition: border-top-color 0.3s ease;
}

.connector-arrow.active {
  border-top-color: #67c23a;
}

.git-info {
  display: flex;
  align-items: center;
  gap: 10px;
  font-family: 'Courier New', monospace;
  color: #606266;
}

.variables-content {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  padding: 15px;
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  color: #495057;
  overflow-x: auto;
}

.log-container {
  max-height: 500px;
  overflow-y: auto;
}

.log-content {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 15px;
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.4;
  border-radius: 4px;
  overflow-x: auto;
  white-space: pre-wrap;
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(230, 162, 60, 0.4);
  }
  70% {
    box-shadow: 0 0 0 10px rgba(230, 162, 60, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(230, 162, 60, 0);
  }
}
</style>