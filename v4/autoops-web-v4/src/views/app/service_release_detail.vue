<script setup>
import { ref, reactive, onMounted, computed, defineOptions, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

defineOptions({
  name: 'ServiceReleaseDetail'
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
  Link
} from '@element-plus/icons-vue'
import appApi from '@/api/app'
import CodeEditor from '@/components/CodeEditor.vue'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const releaseId = computed(() => route.params.id)

// 详情数据
const releaseData = ref({})

// 轮询定时器
const pollTimer = ref(null)

// 日志对话框
const logDialogVisible = ref(false)
const currentService = ref(null)
const serviceLog = ref('')
const logLoading = ref(false)

// 审批对话框
const approveDialogVisible = ref(false)
const approveForm = reactive({
  approval_type: '',
  action: 2,
  remark: ''
})

// 发布对话框
const deployDialogVisible = ref(false)
const deployForm = reactive({
  execution_mode: 1
})

// 回归测试对话框
const regressionDialogVisible = ref(false)
const regressionForm = reactive({
  status: 3,
  remark: ''
})

// 获取详情
const fetchDetail = async () => {
  loading.value = true
  try {
    const response = await appApi.getServiceReleaseDetail(releaseId.value)
    if (response.data?.code === 200) {
      releaseData.value = response.data.data || {}
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

// 启动轮询监控（发布中时）
const startPolling = () => {
  // 清除旧的定时器
  if (pollTimer.value) {
    clearInterval(pollTimer.value)
    pollTimer.value = null
  }

  console.log('启动轮询监控...')
  pollTimer.value = setInterval(async () => {
    try {
      const response = await appApi.getServiceReleaseDetail(releaseId.value)
      if (response.data?.code === 200) {
        const newData = response.data.data || {}
        releaseData.value = newData

        // 如果不再是发布中状态，停止轮询
        if (newData.deploy_status !== 2) {
          console.log('发布状态已变化，停止轮询')
          clearInterval(pollTimer.value)
          pollTimer.value = null

          // 根据结果提示
          if (newData.deploy_status === 3) {
            ElMessage.success('发布成功！')
          } else if (newData.deploy_status === 4) {
            ElMessage.error('发布失败，请查看详情')
          }
        }
      }
    } catch (error) {
      console.error('轮询获取详情失败:', error)
    }
  }, 3000) // 每3秒轮询一次
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
  router.push('/app/service_work')
}

// 审批状态映射
const getApprovalStatusInfo = (status) => {
  const map = {
    1: { text: '未审批', type: 'info', icon: InfoFilled },
    2: { text: '通过', type: 'success', icon: SuccessFilled },
    3: { text: '驳回', type: 'danger', icon: WarningFilled }
  }
  return map[status] || { text: '未知', type: 'info', icon: InfoFilled }
}

// 发布状态映射
const getDeployStatusInfo = (status) => {
  const map = {
    1: { text: '未发布', type: 'info', icon: InfoFilled },
    2: { text: '发布中', type: 'primary', icon: Loading },
    3: { text: '发布成功', type: 'success', icon: SuccessFilled },
    4: { text: '发布失败', type: 'danger', icon: WarningFilled }
  }
  return map[status] || { text: '未知', type: 'info', icon: InfoFilled }
}

// 回归测试状态映射
const getRegressionStatusInfo = (status) => {
  const map = {
    1: { text: '未测试', type: 'info', icon: InfoFilled },
    2: { text: '测试中', type: 'primary', icon: Loading },
    3: { text: '测试成功', type: 'success', icon: SuccessFilled },
    4: { text: '测试失败', type: 'danger', icon: WarningFilled }
  }
  return map[status] || { text: '未知', type: 'info', icon: InfoFilled }
}

// 整体状态映射
const getStatusInfo = (status) => {
  const map = {
    1: { text: '流程进行中', type: 'primary' },
    2: { text: '发布成功', type: 'success' },
    3: { text: '发布失败', type: 'danger' },
    4: { text: '已驳回', type: 'warning' },
    5: { text: '已取消', type: 'info' }
  }
  return map[status] || { text: '未知', type: 'info' }
}

// 流程节点状态
const getNodeStatus = (nodeIndex) => {
  const d = releaseData.value
  if (!d.id) return 'wait'

  // 0: 申请
  if (nodeIndex === 0) return 'finish'

  // 1: 负责人审批
  if (nodeIndex === 1) {
    if (d.owner_approval_status === 2) return 'finish'
    if (d.owner_approval_status === 3) return 'error'
    if (d.owner_approval_status === 1) return 'process'
    return 'wait'
  }

  // 2: 安全审批
  if (nodeIndex === 2) {
    if (d.owner_approval_status !== 2) return 'wait'
    if (d.security_approval_status === 2) return 'finish'
    if (d.security_approval_status === 3) return 'error'
    if (d.security_approval_status === 1) return 'process'
    return 'wait'
  }

  // 3: 测试审批
  if (nodeIndex === 3) {
    if (d.security_approval_status !== 2) return 'wait'
    if (d.test_approval_status === 2) return 'finish'
    if (d.test_approval_status === 3) return 'error'
    if (d.test_approval_status === 1) return 'process'
    return 'wait'
  }

  // 4: 运维发布
  if (nodeIndex === 4) {
    if (d.test_approval_status !== 2) return 'wait'
    if (d.deploy_status === 3) return 'finish'
    if (d.deploy_status === 4) return 'error'
    if (d.deploy_status === 2) return 'process'
    if (d.deploy_status === 1) return 'process'
    return 'wait'
  }

  // 5: 回归测试
  if (nodeIndex === 5) {
    if (d.deploy_status !== 3) return 'wait'
    if (d.regression_test_status === 3) return 'finish'
    if (d.regression_test_status === 4) return 'error'
    if (d.regression_test_status === 2) return 'process'
    if (d.regression_test_status === 1) return 'process'
    return 'wait'
  }

  // 6: 结束
  if (nodeIndex === 6) {
    // 任何审批驳回，显示为异常
    if (d.owner_approval_status === 3 || d.security_approval_status === 3 || d.test_approval_status === 3) {
      return 'error'
    }
    // 发布失败，显示为异常
    if (d.deploy_status === 4) {
      return 'error'
    }
    // 回归测试通过，流程完成
    if (d.regression_test_status === 3) return 'finish'
    // 其他状态为等待
    return 'wait'
  }

  return 'wait'
}

// 检查权限
const checkPermission = async (operationType) => {
  try {
    const response = await appApi.checkServiceReleasePermission(releaseId.value, operationType)
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
const openApproveDialog = async (type) => {
  // 映射审批类型到操作类型
  const operationMap = {
    'owner': 'owner_approval',
    'security': 'security_approval',
    'test': 'test_approval'
  }

  // 先检查权限
  const hasPermission = await checkPermission(operationMap[type])
  if (!hasPermission) {
    return
  }

  approveForm.approval_type = type
  approveForm.action = 2
  approveForm.remark = ''
  approveDialogVisible.value = true
}

// 提交审批
const handleApprove = async () => {
  try {
    const response = await appApi.approveServiceRelease({
      release_id: parseInt(releaseId.value),
      approval_type: approveForm.approval_type,
      action: approveForm.action,
      remark: approveForm.remark
    })
    if (response.data?.code === 200) {
      ElMessage.success('审批成功')
      approveDialogVisible.value = false
      fetchDetail()
    } else {
      ElMessage.error(response.data?.msg || '审批失败')
    }
  } catch (error) {
    console.error('审批失败:', error)
    ElMessage.error('审批失败')
  }
}

// 打开发布对话框
const openDeployDialog = async () => {
  // 检查发布权限
  const hasPermission = await checkPermission('deploy')
  if (!hasPermission) {
    return
  }

  deployForm.execution_mode = 1
  deployDialogVisible.value = true
}

// 提交发布
const handleDeploy = async () => {
  try {
    const response = await appApi.deployServiceRelease({
      release_id: parseInt(releaseId.value),
      execution_mode: deployForm.execution_mode
    })
    if (response.data?.code === 200) {
      ElMessage.success('发布已启动，正在执行Jenkins任务...')
      deployDialogVisible.value = false
      await fetchDetail()

      // 启动轮询监控
      startPolling()
    } else {
      ElMessage.error(response.data?.msg || '发布失败')
    }
  } catch (error) {
    console.error('发布失败:', error)
    ElMessage.error('发布失败')
  }
}

// 打开回归测试对话框
const openRegressionDialog = async () => {
  // 检查回归测试权限
  const hasPermission = await checkPermission('regression_test')
  if (!hasPermission) {
    return
  }

  regressionForm.status = 3
  regressionForm.remark = ''
  regressionDialogVisible.value = true
}

// 提交回归测试
const handleRegression = async () => {
  try {
    const response = await appApi.regressionTestServiceRelease({
      release_id: parseInt(releaseId.value),
      status: regressionForm.status,
      remark: regressionForm.remark
    })
    if (response.data?.code === 200) {
      ElMessage.success('回归测试结果已提交')
      regressionDialogVisible.value = false
      fetchDetail()
    } else {
      ElMessage.error(response.data?.msg || '提交失败')
    }
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('提交失败')
  }
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString()
}

// 可以审批
const canApprove = (type) => {
  const d = releaseData.value
  if (!d.id || d.status !== 1) return false
  if (type === 'owner') return d.owner_approval_status === 1
  if (type === 'security') return d.owner_approval_status === 2 && d.security_approval_status === 1
  if (type === 'test') return d.security_approval_status === 2 && d.test_approval_status === 1
  return false
}

// 可以发布
const canDeploy = computed(() => {
  const d = releaseData.value
  return d.id && d.status === 1 &&
    d.owner_approval_status === 2 &&
    d.security_approval_status === 2 &&
    d.test_approval_status === 2 &&
    (d.deploy_status === 1 || !d.deploy_status)
})

// 可以回归测试
const canRegression = computed(() => {
  const d = releaseData.value
  // 发布成功后，回归测试未完成（未测试或测试中）时可以提交结果
  return d.id && d.status === 1 && d.deploy_status === 3 && (d.regression_test_status === 1 || d.regression_test_status === 2 || !d.regression_test_status)
})

// 查看服务日志
const showServiceLog = async (service) => {
  currentService.value = service
  logDialogVisible.value = true
  serviceLog.value = '正在加载日志...'
  logLoading.value = true

  try {
    // 获取服务日志
    const response = await appApi.getServiceReleaseLog(service.id, 0)
    if (response.data?.code === 200) {
      const logData = response.data.data
      // 后端返回的字段是 log_content
      if (logData.log_content) {
        // 清理 ANSI 转义序列
        serviceLog.value = cleanAnsiEscapeSequences(logData.log_content)
      } else if (logData.log) {
        // 兼容 log 字段
        serviceLog.value = cleanAnsiEscapeSequences(logData.log)
      } else {
        serviceLog.value = '暂无日志数据...'
      }
    } else {
      serviceLog.value = response.data?.msg || '获取日志失败'
    }
  } catch (error) {
    console.error('获取日志失败:', error)
    serviceLog.value = '获取日志失败'
  } finally {
    logLoading.value = false
  }
}

// 清理 ANSI 转义序列
const cleanAnsiEscapeSequences = (text) => {
  if (!text) return ''
  // 移除 ANSI 颜色代码和其他控制字符
  return text.replace(/\x1B\[[0-9;]*[a-zA-Z]/g, '')
            .replace(/\r\n/g, '\n')
            .replace(/\r/g, '\n')
}

// 打开Jenkins日志
const openLogUrl = (url) => {
  if (url) {
    window.open(url, '_blank')
  }
}

onMounted(async () => {
  await fetchDetail()

  // 如果当前是发布中状态，启动轮询
  if (releaseData.value.deploy_status === 2) {
    startPolling()
  }
})

onUnmounted(() => {
  // 组件卸载时清除定时器
  stopPolling()
})
</script>

<template>
  <div class="release-detail-page" v-loading="loading">
    <!-- 顶部操作栏 -->
    <div class="top-bar">
      <el-button type="primary" :icon="ArrowLeft" size="small" @click="goBack">返回列表</el-button>
      <el-button type="success" :icon="Refresh" size="small" @click="fetchDetail">刷新</el-button>
    </div>

    <!-- 基本信息卡片 -->
    <el-card shadow="hover" class="info-card">
      <template #header>
        <div class="card-header">
          <span class="title">
            <el-icon><Document /></el-icon>
            工单详情 #{{ releaseData.id }}
          </span>
          <el-tag :type="getStatusInfo(releaseData.status).type" size="large">
            {{ getStatusInfo(releaseData.status).text }}
          </el-tag>
        </div>
      </template>

      <el-descriptions :column="3" border>
        <el-descriptions-item label="发布标题" :span="3">
          {{ releaseData.title }}
        </el-descriptions-item>
        <el-descriptions-item label="申请人">
          <el-icon><User /></el-icon>
          {{ releaseData.applicant_name }}
        </el-descriptions-item>
        <el-descriptions-item label="服务数量">
          {{ releaseData.service_count || 0 }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          <el-icon><Clock /></el-icon>
          {{ formatTime(releaseData.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="影响功能" :span="3">
          {{ releaseData.impact_feature }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 流程图卡片 -->
    <el-card shadow="hover" class="flow-card">
      <template #header>
        <span class="title">审批流程</span>
      </template>

      <div class="workflow-container">
        <svg class="workflow-svg" viewBox="0 0 1100 320" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMidYMid meet" style="max-width: 1200px;">
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

          <!-- 申请 (矩形) -->
          <rect x="30" y="180" width="70" height="50" :fill="getNodeStatus(0) === 'finish' ? '#67c23a' : '#f0f0f0'" stroke="#67c23a" stroke-width="2" rx="4"/>
          <text x="65" y="210" text-anchor="middle" fill="#fff" font-size="14" font-weight="600">申请</text>

          <!-- 箭头: 申请 -> 负责人审批 -->
          <line x1="100" y1="205" x2="160" y2="205" :stroke="getNodeStatus(0) === 'finish' ? '#67c23a' : '#CDC5BF'" stroke-width="2" :marker-end="getNodeStatus(0) === 'finish' ? 'url(#arrowhead-green)' : 'url(#arrowhead)'"/>
          <text x="130" y="195" text-anchor="middle" fill="#666" font-size="12">提交</text>

          <!-- 负责人审批 (菱形) -->
          <polygon :points="`160,205 190,175 220,205 190,235`" :fill="getNodeStatus(1) === 'finish' ? '#67c23a' : getNodeStatus(1) === 'process' ? '#409eff' : getNodeStatus(1) === 'error' ? '#f56c6c' : '#f0f0f0'" :stroke="getNodeStatus(1) === 'finish' ? '#67c23a' : getNodeStatus(1) === 'process' ? '#409eff' : getNodeStatus(1) === 'error' ? '#f56c6c' : '#ddd'" stroke-width="2"/>
          <text x="190" y="202" text-anchor="middle" :fill="getNodeStatus(1) === 'wait' ? '#666' : '#fff'" font-size="11">负责人</text>
          <text x="190" y="214" text-anchor="middle" :fill="getNodeStatus(1) === 'wait' ? '#666' : '#fff'" font-size="11">审批</text>

          <!-- 箭头: 负责人审批 -> 安全审批 -->
          <line x1="220" y1="205" x2="280" y2="205" :stroke="getNodeStatus(1) === 'finish' ? '#67c23a' : '#CDC5BF'" stroke-width="2" :marker-end="getNodeStatus(1) === 'finish' ? 'url(#arrowhead-green)' : 'url(#arrowhead)'"/>
          <text x="250" y="195" text-anchor="middle" fill="#666" font-size="12">通过</text>

          <!-- 安全审批 (菱形) -->
          <polygon :points="`280,205 310,175 340,205 310,235`" :fill="getNodeStatus(2) === 'finish' ? '#67c23a' : getNodeStatus(2) === 'process' ? '#409eff' : getNodeStatus(2) === 'error' ? '#f56c6c' : '#f0f0f0'" :stroke="getNodeStatus(2) === 'finish' ? '#67c23a' : getNodeStatus(2) === 'process' ? '#409eff' : getNodeStatus(2) === 'error' ? '#f56c6c' : '#ddd'" stroke-width="2"/>
          <text x="310" y="202" text-anchor="middle" :fill="getNodeStatus(2) === 'wait' ? '#666' : '#fff'" font-size="11">安全</text>
          <text x="310" y="214" text-anchor="middle" :fill="getNodeStatus(2) === 'wait' ? '#666' : '#fff'" font-size="11">审批</text>

          <!-- 箭头: 安全审批 -> 测试审批 -->
          <line x1="340" y1="205" x2="400" y2="205" :stroke="getNodeStatus(2) === 'finish' ? '#67c23a' : '#CDC5BF'" stroke-width="2" :marker-end="getNodeStatus(2) === 'finish' ? 'url(#arrowhead-green)' : 'url(#arrowhead)'"/>
          <text x="370" y="195" text-anchor="middle" fill="#666" font-size="12">通过</text>

          <!-- 测试审批 (菱形) -->
          <polygon :points="`400,205 430,175 460,205 430,235`" :fill="getNodeStatus(3) === 'finish' ? '#67c23a' : getNodeStatus(3) === 'process' ? '#409eff' : getNodeStatus(3) === 'error' ? '#f56c6c' : '#f0f0f0'" :stroke="getNodeStatus(3) === 'finish' ? '#67c23a' : getNodeStatus(3) === 'process' ? '#409eff' : getNodeStatus(3) === 'error' ? '#f56c6c' : '#ddd'" stroke-width="2"/>
          <text x="430" y="202" text-anchor="middle" :fill="getNodeStatus(3) === 'wait' ? '#666' : '#fff'" font-size="11">测试</text>
          <text x="430" y="214" text-anchor="middle" :fill="getNodeStatus(3) === 'wait' ? '#666' : '#fff'" font-size="11">审批</text>

          <!-- 箭头: 测试审批 -> 运维发布 -->
          <line x1="460" y1="205" x2="520" y2="205" :stroke="getNodeStatus(3) === 'finish' ? '#67c23a' : '#CDC5BF'" stroke-width="2" :marker-end="getNodeStatus(3) === 'finish' ? 'url(#arrowhead-green)' : 'url(#arrowhead)'"/>
          <text x="490" y="195" text-anchor="middle" fill="#666" font-size="12">通过</text>

          <!-- 运维发布 (矩形) -->
          <rect x="520" y="180" width="80" height="50" :fill="getNodeStatus(4) === 'finish' ? '#67c23a' : getNodeStatus(4) === 'process' ? '#409eff' : getNodeStatus(4) === 'error' ? '#f56c6c' : '#f0f0f0'" :stroke="getNodeStatus(4) === 'finish' ? '#67c23a' : getNodeStatus(4) === 'process' ? '#409eff' : getNodeStatus(4) === 'error' ? '#f56c6c' : '#ddd'" stroke-width="2" rx="4"/>
          <text x="560" y="210" text-anchor="middle" :fill="getNodeStatus(4) === 'wait' ? '#666' : '#fff'" font-size="14" font-weight="600">运维发布</text>

          <!-- 箭头: 运维发布 -> 回归测试 -->
          <line x1="600" y1="205" x2="680" y2="205" :stroke="getNodeStatus(4) === 'finish' ? '#67c23a' : '#CDC5BF'" stroke-width="2" :marker-end="getNodeStatus(4) === 'finish' ? 'url(#arrowhead-green)' : 'url(#arrowhead)'"/>
          <text x="640" y="195" text-anchor="middle" fill="#666" font-size="12">发布成功</text>

          <!-- 回归测试 (菱形) -->
          <polygon :points="`680,205 710,175 740,205 710,235`" :fill="getNodeStatus(5) === 'finish' ? '#67c23a' : getNodeStatus(5) === 'process' ? '#409eff' : getNodeStatus(5) === 'error' ? '#f56c6c' : '#f0f0f0'" :stroke="getNodeStatus(5) === 'finish' ? '#67c23a' : getNodeStatus(5) === 'process' ? '#409eff' : getNodeStatus(5) === 'error' ? '#f56c6c' : '#ddd'" stroke-width="2"/>
          <text x="710" y="202" text-anchor="middle" :fill="getNodeStatus(5) === 'wait' ? '#666' : '#fff'" font-size="11">回归</text>
          <text x="710" y="214" text-anchor="middle" :fill="getNodeStatus(5) === 'wait' ? '#666' : '#fff'" font-size="11">测试</text>

          <!-- 箭头: 回归测试 -> 结束 (通过) -->
          <line x1="740" y1="205" x2="800" y2="205" :stroke="getNodeStatus(5) === 'finish' ? '#67c23a' : '#CDC5BF'" stroke-width="2" :marker-end="getNodeStatus(5) === 'finish' ? 'url(#arrowhead-green)' : 'url(#arrowhead)'"/>
          <text x="770" y="195" text-anchor="middle" fill="#666" font-size="12">通过</text>

          <!-- 结束 (矩形) -->
          <rect x="800" y="180" width="70" height="50" :fill="getNodeStatus(6) === 'finish' ? '#67c23a' : getNodeStatus(6) === 'error' ? '#f56c6c' : '#f0f0f0'" :stroke="getNodeStatus(6) === 'finish' ? '#67c23a' : getNodeStatus(6) === 'error' ? '#f56c6c' : '#ddd'" stroke-width="2" rx="4"/>
          <text x="835" y="210" text-anchor="middle" :fill="getNodeStatus(6) === 'finish' || getNodeStatus(6) === 'error' ? '#fff' : '#666'" font-size="14" font-weight="600">{{ getNodeStatus(6) === 'error' ? '异常' : '结束' }}</text>

          <!-- 运维回滚 (矩形) - 在结束上方 -->
          <rect x="800" y="60" width="70" height="50" fill="#f0f0f0" stroke="#ddd" stroke-width="2" rx="4"/>
          <text x="835" y="90" text-anchor="middle" fill="#666" font-size="14" font-weight="600">运维回滚</text>

          <!-- 箭头: 运维回滚 -> 结束 (向下) -->
          <line x1="835" y1="110" x2="835" y2="180" stroke="#CDC5BF" stroke-width="2" marker-end="url(#arrowhead)"/>

          <!-- 折线: 运维发布 -> 运维回滚 (向上) -->
          <path d="M 560 180 L 560 85 L 800 85" stroke="#CDC5BF" stroke-width="2" fill="none" marker-end="url(#arrowhead)"/>
          <text x="560" y="130" text-anchor="start" fill="#666" font-size="12">回滚</text>

          <!-- 折线: 回归测试 -> 运维回滚 (不通过, 向上) -->
          <path d="M 710 175 L 710 85 L 800 85" stroke="#CDC5BF" stroke-width="2" fill="none" marker-end="url(#arrowhead)"/>
          <text x="745" y="75" text-anchor="middle" fill="#666" font-size="12">不通过</text>

          <!-- 驳回路径: 负责人/安全/测试审批 -> 异常结束 (向下弯曲) -->
          <path v-if="releaseData.owner_approval_status === 3 || releaseData.security_approval_status === 3 || releaseData.test_approval_status === 3" d="M 190 235 L 190 280 L 835 280 L 835 230" stroke="#f56c6c" stroke-width="2" fill="none" marker-end="url(#arrowhead-red)" stroke-dasharray="5,5"/>
          <text v-if="releaseData.owner_approval_status === 3 || releaseData.security_approval_status === 3 || releaseData.test_approval_status === 3" x="500" y="295" text-anchor="middle" fill="#f56c6c" font-size="12">驳回</text>
        </svg>
      </div>

      <!-- 操作按钮区域 -->
      <div class="action-buttons">
        <el-button
          v-if="canApprove('owner')"
          type="primary"
          @click="openApproveDialog('owner')">
          负责人审批
        </el-button>
        <el-button
          v-if="canApprove('security')"
          type="primary"
          @click="openApproveDialog('security')">
          安全审批
        </el-button>
        <el-button
          v-if="canApprove('test')"
          type="primary"
          @click="openApproveDialog('test')">
          测试审批
        </el-button>
        <el-button
          v-if="canDeploy"
          type="success"
          @click="openDeployDialog">
          启动任务
        </el-button>
        <el-button
          v-if="canRegression"
          type="warning"
          @click="openRegressionDialog">
          提交结果
        </el-button>
      </div>

      <!-- 驳回/失败分支说明 -->
      <div class="flow-notes" v-if="releaseData.status === 4 || releaseData.status === 3">
        <el-alert
          :title="releaseData.status === 4 ? '工单已被驳回' : '发布失败'"
          :type="releaseData.status === 4 ? 'warning' : 'error'"
          show-icon
          :closable="false"
        />
      </div>
    </el-card>

    <!-- 服务列表卡片 -->
    <el-card shadow="hover" class="services-card">
      <template #header>
        <span class="title">服务列表</span>
      </template>

      <!-- 发布中进度提示 -->
      <el-alert
        v-if="releaseData.deploy_status === 2"
        title="Jenkins任务执行中，请稍候..."
        type="info"
        :closable="false"
        style="margin-bottom: 16px">
        <template #default>
          <div style="display: flex; align-items: center; gap: 8px;">
            <el-icon class="is-loading"><Loading /></el-icon>
            <span>正在自动刷新状态（每3秒）</span>
          </div>
        </template>
      </el-alert>

      <div class="service-grid">
        <div
          v-for="(service, index) in (releaseData.services || [])"
          :key="service.id || index"
          class="service-card-item"
        >
          <div class="service-header">
            <div class="service-tags">
              <el-tag type="primary" size="small">{{ service.app_name || service.project_name }}</el-tag>
              <el-tag v-if="service.environment" size="small">{{ service.environment }}</el-tag>
              <el-tag v-if="service.branch" type="info" size="small">{{ service.branch }}</el-tag>
            </div>
            <el-tooltip effect="dark" content="查看日志" placement="top">
              <el-button
                @click="showServiceLog(service)"
                type="primary"
                :icon="Document"
                size="small">查看日志</el-button>
            </el-tooltip>
          </div>

          <el-descriptions :column="1" border size="small" class="service-descriptions">
            <el-descriptions-item label="Git仓库">
              {{ service.repo_url }}
            </el-descriptions-item>
            <el-descriptions-item label="Commit">
              <code>{{ service.commit_id }}</code>
            </el-descriptions-item>
            <el-descriptions-item label="影响功能">
              {{ service.impact_feature || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="db变更操作语句">
              {{ service.db_change || '无' }}
            </el-descriptions-item>
            <el-descriptions-item label="配置文件变更">
              {{ service.config_change || '无' }}
            </el-descriptions-item>
            <el-descriptions-item label="备注信息">
              {{ service.remark || '-' }}
            </el-descriptions-item>
          </el-descriptions>

          <!-- 部署状态 -->
          <div class="deploy-status-section" v-if="releaseData.deploy_status >= 2">
            <el-divider content-position="left">部署状态</el-divider>

            <!-- Jenkins任务调用状态 -->
            <div class="status-row">
              <span>Jenkins任务:</span>
              <template v-if="service.build_number > 0">
                <el-tag type="success" size="small">✅ 已调用</el-tag>
                <a v-if="service.log_url" :href="service.log_url" target="_blank" style="margin-left: 8px; font-size: 12px;">
                  #{{ service.build_number }} 查看日志
                </a>
                <span v-else style="margin-left: 8px; font-size: 12px;">#{{ service.build_number }}</span>
              </template>
              <el-tag v-else type="info" size="small">⏸️ 待执行</el-tag>
            </div>

            <!-- 部署结果状态 -->
            <div class="status-row">
              <span>部署结果:</span>
              <el-tag v-if="service.status === 1" type="info" size="small">⏸️ 待执行</el-tag>
              <el-tag v-else-if="service.status === 2" type="primary" size="small">
                <el-icon class="is-loading" style="margin-right: 4px;"><Loading /></el-icon>
                执行中
              </el-tag>
              <el-tag v-else-if="service.status === 3" type="success" size="small">✅ 部署成功</el-tag>
              <el-tag v-else-if="service.status === 4" type="danger" size="small">❌ 部署失败</el-tag>
              <el-tag v-else type="info" size="small">未知</el-tag>
            </div>

            <!-- 执行耗时 -->
            <div class="status-row" v-if="service.duration > 0">
              <span>执行耗时:</span>
              <span style="font-size: 12px; color: #666;">{{ service.duration }}秒</span>
            </div>

            <!-- 错误信息 -->
            <div class="status-row" v-if="service.status === 4 && service.error_message">
              <span>错误信息:</span>
              <el-text type="danger" size="small" style="word-break: break-all;">
                {{ service.error_message }}
              </el-text>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 审批对话框 -->
    <el-dialog v-model="approveDialogVisible" title="审批" width="500px">
      <el-form :model="approveForm" label-width="80px">
        <el-form-item label="审批结果">
          <el-radio-group v-model="approveForm.action">
            <el-radio :value="2">通过</el-radio>
            <el-radio :value="3">驳回</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审批意见">
          <el-input v-model="approveForm.remark" type="textarea" :rows="3" placeholder="请输入审批意见" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approveDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleApprove">确认</el-button>
      </template>
    </el-dialog>

    <!-- 发布对话框 -->
    <el-dialog v-model="deployDialogVisible" title="执行发布" width="500px">
      <el-form :model="deployForm" label-width="100px">
        <el-form-item label="执行模式">
          <el-radio-group v-model="deployForm.execution_mode">
            <el-radio :value="1">并行发布 (快速)</el-radio>
            <el-radio :value="2">串行发布 (安全)</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-alert
          type="warning"
          title="请确认所有服务配置正确，发布后将自动执行Jenkins任务"
          :closable="false"
          show-icon
        />
      </el-form>
      <template #footer>
        <el-button @click="deployDialogVisible = false">取消</el-button>
        <el-button type="success" @click="handleDeploy">确认发布</el-button>
      </template>
    </el-dialog>

    <!-- 回归测试对话框 -->
    <el-dialog v-model="regressionDialogVisible" title="回归测试" width="500px">
      <el-form :model="regressionForm" label-width="100px">
        <el-form-item label="测试结果">
          <el-radio-group v-model="regressionForm.status">
            <el-radio :value="3">测试成功</el-radio>
            <el-radio :value="4">测试失败</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="regressionForm.remark" type="textarea" :rows="3" placeholder="请输入测试备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="regressionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRegression">提交结果</el-button>
      </template>
    </el-dialog>

    <!-- 服务日志对话框 -->
    <el-dialog
      v-model="logDialogVisible"
      :title="currentService?.app_name || '服务部署日志'"
      width="800px">
      <div v-if="currentService" class="log-detail-content">
        <el-descriptions :column="2" border size="small" style="margin-bottom: 16px;">
          <el-descriptions-item label="应用名称">{{ currentService.app_name }}</el-descriptions-item>
          <el-descriptions-item label="环境">{{ currentService.environment }}</el-descriptions-item>
          <el-descriptions-item label="分支">{{ currentService.branch || '-' }}</el-descriptions-item>
          <el-descriptions-item label="Commit">
            <code style="font-size: 12px;">{{ currentService.commit_id?.substring(0, 8) }}</code>
          </el-descriptions-item>
          <el-descriptions-item label="构建号">
            {{ currentService.build_number || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag v-if="currentService.status === 1" type="info" size="small">待执行</el-tag>
            <el-tag v-else-if="currentService.status === 2" type="primary" size="small">执行中</el-tag>
            <el-tag v-else-if="currentService.status === 3" type="success" size="small">成功</el-tag>
            <el-tag v-else-if="currentService.status === 4" type="danger" size="small">失败</el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="currentService.error_message" style="margin-bottom: 16px;">
          <el-alert
            :title="currentService.error_message"
            type="error"
            :closable="false"
            show-icon>
          </el-alert>
        </div>

        <div class="log-section">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
            <h4 style="margin: 0;">日志内容</h4>
            <el-button
              v-if="currentService.log_url"
              @click="openLogUrl(currentService.log_url)"
              type="primary"
              size="small"
              :icon="Link">
              打开Jenkins日志
            </el-button>
          </div>
          <div v-loading="logLoading" class="log-content-wrapper">
            <CodeEditor
              :model-value="serviceLog"
              language="bash"
              height="400px"
              :readonly="true"
              fontSize="12px"
            />
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.release-detail-page {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.top-bar {
  margin-bottom: 20px;
  display: flex;
  gap: 12px;
}

.info-card,
.flow-card,
.services-card,
.approvers-card {
  margin-bottom: 20px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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
  padding: 20px;
  min-height: 350px;
  width: 100%;
}

.workflow-svg {
  width: 100%;
  height: auto;
  display: block;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 20px 0;
  margin-top: 20px;
  border-top: 1px solid #ebeef5;
}

.flow-notes {
  margin-top: 20px;
}

/* 服务列表网格布局 - 每行3个 */
.service-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

@media (max-width: 1400px) {
  .service-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 900px) {
  .service-grid {
    grid-template-columns: 1fr;
  }
}

.service-card-item {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
  border: 1px solid #e4e7ed;
  transition: all 0.3s;
}

.service-card-item:hover {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.service-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.service-descriptions {
  margin-top: 12px;
}

.service-descriptions :deep(.el-descriptions__label) {
  width: 120px;
  font-size: 12px;
}

.service-descriptions :deep(.el-descriptions__content) {
  font-size: 12px;
  word-break: break-all;
}

.deploy-status-section {
  margin-top: 12px;
}

.status-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 12px;
}

/* 日志内容容器 */
.log-content-wrapper {
  max-height: 400px;
  overflow-y: auto;
}

.log-content-wrapper :deep(.code-editor-container) {
  height: auto !important;
  min-height: 400px;
}
</style>
