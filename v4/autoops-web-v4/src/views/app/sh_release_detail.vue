<script setup>
import { ref, reactive, onMounted, computed, defineOptions, onUnmounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'

defineOptions({
  name: 'ShReleaseDetail'
})
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft,
  CircleCheck,
  CircleClose,
  Loading,
  Document,
  User,
  Clock,
  Warning,
  SuccessFilled,
  WarningFilled,
  InfoFilled,
  Refresh,
  Download,
  Folder
} from '@element-plus/icons-vue'
import appApi from '@/api/app'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const releaseId = computed(() => route.params.id)

// 详情数据
const releaseData = ref({})

// 轮询定时器（用于实时获取执行状态和日志）
const pollTimer = ref(null)
let logStart = 0

// 审批对话框
const approveDialogVisible = ref(false)
const approveForm = reactive({
  action: 2,
  remark: 'ok'
})

// 执行Jenkins任务对话框
const executeDialogVisible = ref(false)

// 执行脚本对话框
const executeScriptDialogVisible = ref(false)

// 下载文件相关
const downloadDialogVisible = ref(false)
const downloading = ref(false)
const downloadingFile = ref(null)
const fileList = ref([])

// 是否合并文件下的开关
const isMerged = ref(true)

// Tab 选择
const activeTab = ref('detail')

// Jenkins执行日志内容
const executionLog = ref('')

// 脚本执行输出
const scriptOutput = ref('')

// 获取详情
const fetchDetail = async () => {
  loading.value = true
  try {
    const response = await appApi.getShReleaseDetail(releaseId.value)
    if (response.data?.code === 200) {
      releaseData.value = response.data.data || {}

      // 如果有 server_host_id 但没有 server_ip，从文件列表API获取
      if (releaseData.value.server_host_id && !releaseData.value.server_ip) {
        try {
          const filesResponse = await appApi.listShReleaseFiles(releaseId.value)
          if (filesResponse.data?.code === 200) {
            const filesData = filesResponse.data.data
            if (filesData?.host_ip) {
              releaseData.value.server_ip = filesData.host_ip
            }
          }
        } catch (err) {
          console.error('获取服务器IP失败:', err)
        }
      }
    } else {
      ElMessage.error(response.data?.msg || '获取详情失败')
    }
  } catch (error) {
    console.error('获取详情失败:', error)
    ElMessage.error('获取详情失败')
  } finally {
    loading.value = false
  }
}

// 启动轮询监控（执行中时）
const startPolling = () => {
  // 清除旧的定时器
  if (pollTimer.value) {
    clearInterval(pollTimer.value)
    pollTimer.value = null
  }

  console.log('启动轮询监控...')
  pollTimer.value = setInterval(async () => {
    try {
      // 获取最新详情
      const detailResponse = await appApi.getShReleaseDetail(releaseId.value)
      if (detailResponse.data?.code === 200) {
        const newData = detailResponse.data.data || {}
        releaseData.value = newData

        // 如果是Jenkins执行中(status=2)，获取Jenkins日志
        if (newData.execute_status === 2) {
          try {
            const logResponse = await appApi.getShReleaseLog(releaseId.value, logStart)
            if (logResponse.data?.code === 200) {
              const logData = logResponse.data.data
              if (logData.log_content) {
                executionLog.value += logData.log_content
              }
              logStart = parseInt(logData.next_start) || logStart
            }
          } catch (logError) {
            // 日志文件可能还未生成，忽略404错误
            if (logError.response?.status !== 404) {
              console.error('获取Jenkins日志失败:', logError)
            }
          }
        }

        // 如果是脚本执行中(status=5)，获取脚本输出
        if (newData.execute_status === 5) {
          try {
            await loadScriptOutput()
          } catch (scriptError) {
            // 脚本输出可能还未生成，忽略404错误
            if (scriptError.response?.status !== 404) {
              console.error('获取脚本输出失败:', scriptError)
            }
          }
        }

        // 如果不再是执行中状态，停止轮询
        if (newData.execute_status !== 2 && newData.execute_status !== 5) {
          console.log('执行状态已变化，停止轮询')
          clearInterval(pollTimer.value)
          pollTimer.value = null

          // 根据结果提示
          if (newData.execute_status === 3) {
            ElMessage.success('代码拉取成功，可以执行脚本！')
          } else if (newData.execute_status === 4) {
            ElMessage.error('Jenkins执行失败，请查看日志')
          } else if (newData.execute_status === 6) {
            ElMessage.success('脚本执行成功！')
            // 最后加载一次脚本输出
            await loadScriptOutput()
          } else if (newData.execute_status === 7) {
            ElMessage.error('脚本执行失败，请查看输出')
            // 最后加载一次脚本输出
            await loadScriptOutput()
          }
        }
      }
    } catch (error) {
      // 只记录非404的错误
      if (error.response?.status !== 404) {
        console.error('轮询获取详情失败:', error)
      }
    }
  }, 2000) // 每2秒轮询一次
}

// 停止轮询
const stopPolling = () => {
  if (pollTimer.value) {
    clearInterval(pollTimer.value)
    pollTimer.value = null
    console.log('停止轮询监控')
  }
}

// 返回列表
const goBack = () => {
  router.push('/app/sh-release')
}

// 审核状态映射
const getApprovalStatusInfo = (status) => {
  const map = {
    1: { text: '未审核', type: 'info', icon: InfoFilled },
    2: { text: '已审核', type: 'success', icon: SuccessFilled },
    3: { text: '已驳回', type: 'danger', icon: WarningFilled }
  }
  return map[status] || { text: '未知', type: 'info', icon: InfoFilled }
}

// 执行状态映射
const getExecuteStatusInfo = (status) => {
  const map = {
    1: { text: '未执行', type: 'info', icon: InfoFilled },
    2: { text: '代码拉取中', type: 'primary', icon: Loading },
    3: { text: '代码拉取成功', type: 'success', icon: SuccessFilled },
    4: { text: 'Jenkins执行失败', type: 'danger', icon: WarningFilled },
    5: { text: '脚本执行中', type: 'primary', icon: Loading },
    6: { text: '脚本执行成功', type: 'success', icon: SuccessFilled },
    7: { text: '脚本执行失败', type: 'danger', icon: WarningFilled }
  }
  return map[status] || { text: '未知', type: 'info', icon: InfoFilled }
}

// 整体状态映射
const getStatusInfo = (status) => {
  const map = {
    1: { text: '流程进行中', type: 'primary' },
    2: { text: '执行成功', type: 'success' },
    3: { text: '执行失败', type: 'danger' },
    4: { text: '已驳回', type: 'warning' },
    5: { text: '已取消', type: 'info' }
  }
  return map[status] || { text: '未知', type: 'info' }
}

// 流程节点状态
const getNodeStatus = (nodeIndex) => {
  const d = releaseData.value
  if (!d.id) return 'wait'

  // 0: 创建工单
  if (nodeIndex === 0) return 'finish'

  // 1: 审核人
  if (nodeIndex === 1) {
    if (d.approval_status === 2) return 'finish'
    if (d.approval_status === 3) return 'error'
    if (d.approval_status === 1) return 'process'
    return 'wait'
  }

  // 2: 拉取代码
  if (nodeIndex === 2) {
    if (d.approval_status !== 2) return 'wait'
    if (d.execute_status >= 3) return 'finish'
    if (d.execute_status === 4) return 'error'
    if (d.execute_status === 2) return 'process'
    return 'wait'
  }

  // 3: 执行脚本
  if (nodeIndex === 3) {
    if (d.execute_status < 3) return 'wait'
    if (d.execute_status === 6) return 'finish'
    if (d.execute_status === 7) return 'error'
    if (d.execute_status === 5) return 'process'
    return 'wait'
  }

  // 4: 完成
  if (nodeIndex === 4) {
    if (d.approval_status === 3) return 'error'
    if (d.execute_status === 4 || d.execute_status === 7) return 'error'
    if (d.execute_status === 6) return 'finish'
    return 'wait'
  }

  return 'wait'
}

// 检查权限
const checkPermission = async (operationType) => {
  try {
    const response = await appApi.checkShReleasePermission(releaseId.value, operationType)
    if (response.data?.code === 200) {
      const data = response.data.data
      if (!data.has_permission) {
        ElMessage.error(data.message || '您没有权限执行此操作')
        return false
      }
      return true
    }
    return false
  } catch (error) {
    console.error('权限检查失败:', error)
    ElMessage.error('权限检查失败')
    return false
  }
}

// 打开审批对话框
const openApproveDialog = async () => {
  // 先检查权限
  const hasPermission = await checkPermission('approval')
  if (!hasPermission) {
    return
  }

  approveForm.action = 2
  approveForm.remark = 'ok'
  approveDialogVisible.value = true
}

// 提交审批
const handleApprove = async () => {
  if (!approveForm.remark) {
    ElMessage.warning('请填写审核意见')
    return
  }

  try {
    const response = await appApi.approveShRelease({
      release_id: parseInt(releaseId.value),
      action: approveForm.action,
      remark: approveForm.remark
    })
    if (response.data?.code === 200) {
      ElMessage.success(approveForm.action === 2 ? '审核通过' : '审核驳回')
      approveDialogVisible.value = false
      fetchDetail()
    } else {
      ElMessage.error(response.data?.msg || '审核失败')
    }
  } catch (error) {
    console.error('审核失败:', error)
    ElMessage.error('审核失败')
  }
}

// 打开执行Jenkins任务对话框
const openExecuteDialog = async () => {
  // 检查执行权限
  const hasPermission = await checkPermission('execute')
  if (!hasPermission) {
    return
  }

  executeDialogVisible.value = true
}

// 提交执行Jenkins任务
const handleExecute = async () => {
  try {
    const response = await appApi.executeShRelease({
      release_id: parseInt(releaseId.value)
    })
    if (response.data?.code === 200) {
      ElMessage.success('Jenkins任务执行已启动...')
      executeDialogVisible.value = false
      await fetchDetail()

      // 切换到部署日志 tab
      activeTab.value = 'log'

      // 重置日志起始位置
      logStart = 0
      executionLog.value = ''

      // 启动轮询监控
      startPolling()
    } else {
      ElMessage.error(response.data?.msg || '执行失败')
    }
  } catch (error) {
    console.error('执行失败:', error)
    ElMessage.error('执行失败')
  }
}

// 打开执行脚本对话框
const openExecuteScriptDialog = async () => {
  // 检查执行权限
  const hasPermission = await checkPermission('execute')
  if (!hasPermission) {
    return
  }

  executeScriptDialogVisible.value = true
}

// 提交执行脚本
const handleExecuteScript = async () => {
  try {
    const response = await appApi.executeShReleaseScript({
      release_id: parseInt(releaseId.value)
    })
    if (response.data?.code === 200) {
      ElMessage.success('脚本执行已启动...')
      executeScriptDialogVisible.value = false
      await fetchDetail()

      // 切换到脚本输出 tab
      activeTab.value = 'script'

      // 启动轮询监控
      startPolling()
    } else {
      ElMessage.error(response.data?.msg || '执行脚本失败')
    }
  } catch (error) {
    console.error('执行脚本失败:', error)
    ElMessage.error('执行脚本失败')
  }
}

// 加载Jenkins执行日志
const loadExecutionLog = async () => {
  try {
    const response = await appApi.getShReleaseLog(releaseId.value, 0)
    if (response.data?.code === 200) {
      const logData = response.data.data
      executionLog.value = logData.log_content || '暂无日志'
      logStart = parseInt(logData.next_start) || 0
    }
  } catch (error) {
    console.error('获取日志失败:', error)
    executionLog.value = '获取日志失败'
  }
}

// 加载脚本执行输出
const loadScriptOutput = async () => {
  try {
    const response = await appApi.getShReleaseScriptOutput(releaseId.value)
    if (response.data?.code === 200) {
      const data = response.data.data
      scriptOutput.value = data.script_output || '暂无输出'
    }
  } catch (error) {
    console.error('获取脚本输出失败:', error)
    scriptOutput.value = '获取脚本输出失败'
  }
}

// Tab切换事件处理
const handleTabChange = async (tabName) => {
  if (tabName === 'log') {
    // 切换到Jenkins日志tab时，加载Jenkins日志
    await loadExecutionLog()
  } else if (tabName === 'script') {
    // 切换到脚本输出tab时，加载脚本输出
    await loadScriptOutput()
  }
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString()
}

// 下载文件
const handleDownloadFile = async () => {
  try {
    downloading.value = true
    const response = await appApi.listShReleaseFiles(releaseId.value)
    if (response.data?.code === 200) {
      const data = response.data.data
      fileList.value = data?.files || []

      // 从返回的数据中获取服务器IP（如果还没有的话）
      if (data?.host_ip && !releaseData.value.server_ip) {
        releaseData.value.server_ip = data.host_ip
      }

      downloadDialogVisible.value = true
    } else {
      ElMessage.error(response.data?.msg || '获取文件列表失败')
    }
  } catch (error) {
    console.error('获取文件列表失败:', error)
    ElMessage.error('获取文件列表失败')
  } finally {
    downloading.value = false
  }
}

// 下载选中的文件
const handleDownloadSelectedFile = async (file) => {
  try {
    downloadingFile.value = file.path
    ElMessage.info(`正在下载 ${file.name}...`)

    const response = await appApi.downloadShReleaseFile(releaseId.value, file.path)

    // 创建下载链接
    const url = window.URL.createObjectURL(response.data)
    const link = document.createElement('a')
    link.href = url
    link.download = file.name
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    ElMessage.success(`${file.name} 下载成功`)
  } catch (error) {
    console.error('文件下载失败:', error)
    ElMessage.error('文件下载失败: ' + (error.message || '未知错误'))
  } finally {
    downloadingFile.value = null
  }
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

// 可以审批
const canApprove = computed(() => {
  const d = releaseData.value
  return d.id && d.status === 1 && d.approval_status === 1
})

// 可以执行Jenkins任务（已审核 且 未执行Jenkins）
const canExecute = computed(() => {
  const d = releaseData.value
  return d.id && d.status === 1 && d.approval_status === 2 && (d.execute_status === 1 || !d.execute_status)
})

// 可以执行脚本（Jenkins任务执行成功后，代码拉取成功，execute_status === 3）
const canExecuteScript = computed(() => {
  const d = releaseData.value
  // 代码拉取成功(execute_status === 3)，且脚本还未执行(script_execute_status为空或1)
  return d.id && d.status === 1 && d.approval_status === 2 && d.execute_status === 3
})

onMounted(async () => {
  await fetchDetail()

  // 如果当前是执行中状态（代码拉取中或脚本执行中），启动轮询
  if (releaseData.value.execute_status === 2 || releaseData.value.execute_status === 5) {
    activeTab.value = 'log'
    await loadExecutionLog()
    startPolling()
  }
})

onBeforeUnmount(() => {
  // 组件卸载时清除定时器
  stopPolling()
})
</script>

<template>
  <div class="sh-release-detail-page" v-loading="loading">
    <!-- 顶部操作栏 -->
    <div class="top-bar">
      <el-button type="primary" :icon="ArrowLeft" size="small" @click="goBack">返回列表</el-button>
      <el-button type="success" :icon="Refresh" size="small" @click="fetchDetail">刷新</el-button>
    </div>

    <!-- 基本信息卡片 -->
    <el-card shadow="hover" class="info-card">
      <template #header>
        <div class="card-header">
          <span class="title">标题: {{ releaseData.title }}</span>
          <div class="header-right">
            <el-button
              v-if="canApprove"
              type="success"
              size="small"
              @click="openApproveDialog">
              执行工单
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="goBack">
              删除
            </el-button>
          </div>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="申请人">
          <el-tag type="primary" size="small">{{ releaseData.applicant_name }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="审批人">
          <el-tag size="small">{{ releaseData.approver_name }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述">
          {{ releaseData.reason }}
        </el-descriptions-item>
        <el-descriptions-item label="执行人">
          <el-tag size="small">{{ releaseData.executor_name }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="运行服务器">
          {{ releaseData.server_ip || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="Git地址">
          {{ releaseData.application?.repo_url || releaseData.repo_url || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ formatTime(releaseData.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="下载文件">
          <div style="display: flex; align-items: center; gap: 8px;">
            <el-button
              type="primary"
              size="small"
              :icon="Download"
              :loading="downloading"
              @click="handleDownloadFile"
              style="flex-shrink: 0;">
              下载
            </el-button>
            <span style="color: #909399; font-size: 12px; flex: 1; min-width: 0;">{{ releaseData.execute_dir }}</span>
          </div>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 工单进度和审核日志 -->
    <el-card shadow="hover" class="flow-card">
      <template #header>
        <span class="title">工单进度</span>
      </template>

      <div class="workflow-container">
        <svg class="workflow-svg" viewBox="0 0 750 160" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMidYMid meet">
          <!-- 定义箭头 -->
          <defs>
            <marker id="arrowhead" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
              <polygon points="0 0, 10 3, 0 6" fill="#CDC5BF" />
            </marker>
            <marker id="arrowhead-green" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
              <polygon points="0 0, 10 3, 0 6" fill="#67c23a" />
            </marker>
            <marker id="arrowhead-red" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
              <polygon points="0 0, 10 3, 0 6" fill="#f56c6c" />
            </marker>
          </defs>

          <!-- 创建工单 -->
          <rect x="20" y="50" width="100" height="60" :fill="getNodeStatus(0) === 'finish' ? '#67c23a' : '#f0f0f0'" :stroke="getNodeStatus(0) === 'finish' ? '#67c23a' : '#ddd'" stroke-width="2" rx="5"/>
          <text x="70" y="85" text-anchor="middle" :fill="getNodeStatus(0) === 'finish' ? '#fff' : '#666'" font-size="16" font-weight="600">创建工单</text>

          <!-- 箭头: 创建工单 -> 审核人 -->
          <line x1="120" y1="80" x2="165" y2="80" :stroke="getNodeStatus(0) === 'finish' ? '#67c23a' : '#CDC5BF'" stroke-width="2" :marker-end="getNodeStatus(0) === 'finish' ? 'url(#arrowhead-green)' : 'url(#arrowhead)'"/>

          <!-- 审核人 -->
          <polygon :points="`165,80 200,50 235,80 200,110`" :fill="getNodeStatus(1) === 'finish' ? '#67c23a' : getNodeStatus(1) === 'process' ? '#409eff' : getNodeStatus(1) === 'error' ? '#f56c6c' : '#f0f0f0'" :stroke="getNodeStatus(1) === 'finish' ? '#67c23a' : getNodeStatus(1) === 'process' ? '#409eff' : getNodeStatus(1) === 'error' ? '#f56c6c' : '#ddd'" stroke-width="2"/>
          <text x="200" y="85" text-anchor="middle" :fill="getNodeStatus(1) === 'wait' ? '#666' : '#fff'" font-size="15" font-weight="600">审核人</text>

          <!-- 箭头: 审核人 -> 拉取代码 -->
          <line x1="235" y1="80" x2="280" y2="80" :stroke="getNodeStatus(1) === 'finish' ? '#67c23a' : '#CDC5BF'" stroke-width="2" :marker-end="getNodeStatus(1) === 'finish' ? 'url(#arrowhead-green)' : 'url(#arrowhead)'"/>

          <!-- 拉取代码 -->
          <rect x="280" y="50" width="100" height="60" :fill="getNodeStatus(2) === 'finish' ? '#67c23a' : getNodeStatus(2) === 'process' ? '#409eff' : getNodeStatus(2) === 'error' ? '#f56c6c' : '#f0f0f0'" :stroke="getNodeStatus(2) === 'finish' ? '#67c23a' : getNodeStatus(2) === 'process' ? '#409eff' : getNodeStatus(2) === 'error' ? '#f56c6c' : '#ddd'" stroke-width="2" rx="5"/>
          <text x="330" y="85" text-anchor="middle" :fill="getNodeStatus(2) === 'wait' ? '#666' : '#fff'" font-size="16" font-weight="600">拉取代码</text>

          <!-- 箭头: 拉取代码 -> 执行脚本 -->
          <line x1="380" y1="80" x2="425" y2="80" :stroke="getNodeStatus(2) === 'finish' ? '#67c23a' : '#CDC5BF'" stroke-width="2" :marker-end="getNodeStatus(2) === 'finish' ? 'url(#arrowhead-green)' : 'url(#arrowhead)'"/>

          <!-- 执行脚本 -->
          <rect x="425" y="50" width="100" height="60" :fill="getNodeStatus(3) === 'finish' ? '#67c23a' : getNodeStatus(3) === 'process' ? '#409eff' : getNodeStatus(3) === 'error' ? '#f56c6c' : '#f0f0f0'" :stroke="getNodeStatus(3) === 'finish' ? '#67c23a' : getNodeStatus(3) === 'process' ? '#409eff' : getNodeStatus(3) === 'error' ? '#f56c6c' : '#ddd'" stroke-width="2" rx="5"/>
          <text x="475" y="85" text-anchor="middle" :fill="getNodeStatus(3) === 'wait' ? '#666' : '#fff'" font-size="16" font-weight="600">执行脚本</text>

          <!-- 箭头: 执行脚本 -> 完成 -->
          <line x1="525" y1="80" x2="570" y2="80" :stroke="getNodeStatus(3) === 'finish' ? '#67c23a' : '#CDC5BF'" stroke-width="2" :marker-end="getNodeStatus(3) === 'finish' ? 'url(#arrowhead-green)' : 'url(#arrowhead)'"/>

          <!-- 完成 -->
          <rect x="570" y="50" width="100" height="60" :fill="getNodeStatus(4) === 'finish' ? '#67c23a' : getNodeStatus(4) === 'error' ? '#f56c6c' : '#f0f0f0'" :stroke="getNodeStatus(4) === 'finish' ? '#67c23a' : getNodeStatus(4) === 'error' ? '#f56c6c' : '#ddd'" stroke-width="2" rx="5"/>
          <text x="620" y="85" text-anchor="middle" :fill="getNodeStatus(4) === 'finish' || getNodeStatus(4) === 'error' ? '#fff' : '#666'" font-size="16" font-weight="600">{{ getNodeStatus(4) === 'error' ? '失败' : '完成' }}</text>
        </svg>
      </div>

      <!-- 审核日志 -->
      <div class="audit-section">
        <div class="audit-title">审核日志</div>
        <div class="audit-list">
          <!-- 审核人记录 -->
          <div class="audit-item" v-if="releaseData.approval_status >= 1">
            <div class="audit-header">
              <span class="audit-user">审核人: {{ releaseData.approver_name }}</span>
              <el-tag :type="getApprovalStatusInfo(releaseData.approval_status).type" size="small">
                {{ getApprovalStatusInfo(releaseData.approval_status).text }}
              </el-tag>
              <span class="audit-time">审核时间: {{ formatTime(releaseData.approval_time) }}</span>
            </div>
            <div class="audit-remark" v-if="releaseData.approval_remark">
              附加信息: {{ releaseData.approval_remark }}
            </div>
          </div>

          <!-- 执行人记录 -->
          <div class="audit-item" v-if="releaseData.execute_status >= 2">
            <div class="audit-header">
              <span class="audit-user">执行人: {{ releaseData.executor_name }}</span>
              <el-tag :type="getExecuteStatusInfo(releaseData.execute_status).type" size="small">
                {{ getExecuteStatusInfo(releaseData.execute_status).text }}
              </el-tag>
              <span class="audit-time">审核时间: {{ formatTime(releaseData.execute_time) }}</span>
            </div>
            <div class="audit-remark" v-if="releaseData.execute_remark">
              附加信息: {{ releaseData.execute_remark }}
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮区域 -->
      <div class="action-buttons">
        <el-button
          v-if="canApprove"
          type="primary"
          @click="openApproveDialog">
          审核工单
        </el-button>
        <el-button
          v-if="canExecute"
          type="success"
          @click="openExecuteDialog">
          执行Jenkins任务
        </el-button>
        <el-button
          v-if="canExecuteScript"
          type="warning"
          @click="openExecuteScriptDialog">
          执行脚本
        </el-button>
      </div>
    </el-card>

    <!-- Tab 切换：工单内容 / Jenkins日志 / 脚本输出 -->
    <el-card shadow="hover" class="content-card">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="工单内容" name="detail">
          <div class="script-content">
            <div class="script-line">
              <span class="line-number">1</span>
              <code class="script-code">{{ releaseData.script_content }}</code>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="Jenkins日志" name="log">
          <div class="log-container">
            <pre>{{ executionLog || '暂无日志' }}</pre>
          </div>
        </el-tab-pane>
        <el-tab-pane label="脚本输出" name="script">
          <div class="log-container">
            <pre>{{ scriptOutput || '暂无输出' }}</pre>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 审批对话框 -->
    <el-dialog v-model="approveDialogVisible" title="审批" width="500px">
      <el-form :model="approveForm" label-width="80px">
        <el-form-item label="审批结果">
          <el-radio-group v-model="approveForm.action">
            <el-radio :value="2">已审核</el-radio>
            <el-radio :value="3">未处理</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审核意见" required>
          <el-input v-model="approveForm.remark" type="textarea" :rows="3" placeholder="请输入审核意见" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approveDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleApprove">确认</el-button>
      </template>
    </el-dialog>

    <!-- 执行Jenkins任务对话框 -->
    <el-dialog v-model="executeDialogVisible" title="执行Jenkins任务" width="500px">
      <el-alert
        type="warning"
        title="确认执行Jenkins任务？执行后将拉取代码到服务器"
        :closable="false"
        show-icon
      />
      <template #footer>
        <el-button @click="executeDialogVisible = false">取消</el-button>
        <el-button type="success" @click="handleExecute">确认执行</el-button>
      </template>
    </el-dialog>

    <!-- 执行脚本对话框 -->
    <el-dialog v-model="executeScriptDialogVisible" title="执行脚本" width="500px">
      <el-alert
        type="warning"
        title="确认执行脚本？将通过SSH在服务器上执行脚本"
        :closable="false"
        show-icon
      />
      <template #footer>
        <el-button @click="executeScriptDialogVisible = false">取消</el-button>
        <el-button type="warning" @click="handleExecuteScript">确认执行</el-button>
      </template>
    </el-dialog>

    <!-- 下载文件对话框 -->
    <el-dialog v-model="downloadDialogVisible" title="下载文件" width="700px">
      <el-table :data="fileList" style="width: 100%">
        <el-table-column label="文件名" min-width="200">
          <template #default="{ row }">
            <div style="display: flex; align-items: center; gap: 8px;">
              <el-icon :size="18" :color="row.isDir ? '#409eff' : '#67c23a'">
                <Folder v-if="row.isDir" />
                <Document v-else />
              </el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="120">
          <template #default="{ row }">
            {{ row.isDir ? '-' : formatFileSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column prop="modTime" label="修改时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="!row.isDir"
              type="primary"
              size="small"
              link
              :loading="downloadingFile === row.path"
              @click.stop="handleDownloadSelectedFile(row)">
              下载
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="downloadDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.sh-release-detail-page {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

/* 确保descriptions内容区域不被截断 */
.sh-release-detail-page :deep(.el-descriptions__body .el-descriptions__cell) {
  overflow: visible;
}

.sh-release-detail-page :deep(.el-descriptions__content) {
  overflow: visible;
  word-break: break-all;
}

.top-bar {
  margin-bottom: 20px;
  display: flex;
  gap: 12px;
}

.info-card,
.flow-card,
.audit-card,
.content-card {
  margin-bottom: 20px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-right {
  display: flex;
  gap: 8px;
}

.title {
  font-size: 16px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* SVG流程图样式 */
.workflow-container {
  position: relative;
  padding: 15px 20px;
  min-height: 130px;
  width: 100%;
}

.workflow-svg {
  width: 100%;
  max-width: 800px;
  height: auto;
  display: block;
  margin: 0 auto;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 20px 0;
  margin-top: 20px;
  border-top: 1px solid #ebeef5;
}

/* 审核日志区域样式 */
.audit-section {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}

.audit-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #303133;
}

/* 审核日志样式 */
.audit-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.audit-item {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
  border: 1px solid #e4e7ed;
}

.audit-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.audit-user {
  font-weight: 600;
  color: #303133;
}

.audit-time {
  font-size: 12px;
  color: #909399;
  margin-left: auto;
}

.audit-remark {
  font-size: 14px;
  color: #606266;
  margin-top: 8px;
  padding-left: 12px;
  border-left: 3px solid #409eff;
}

/* 脚本内容样式 */
.script-content {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
}

.script-line {
  display: flex;
  gap: 16px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
}

.line-number {
  color: #858585;
  user-select: none;
  min-width: 20px;
  text-align: right;
}

.script-code {
  color: #569cd6;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 日志容器样式 */
.log-container {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 15px;
  border-radius: 8px;
  max-height: 500px;
  overflow-y: auto;
}

.log-container pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
