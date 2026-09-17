<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  View,
  Edit,
  Delete,
  Search,
  Refresh,
  Plus,
  Platform,
  SuccessFilled,
  CircleCloseFilled,
  WarningFilled,
  RemoveFilled,
  Upload,
  Document,
  Setting,
  Monitor,
  CloudUpload
} from '@element-plus/icons-vue'
import { domainApi } from '@/api/monitor'
import configApi from '@/api/config'
import SslCertApply from './ssl/SslCertApply.vue'
import AliyunConfig from './ssl/AliyunConfig.vue'
import dingtalkIcon from '@/assets/image/钉钉01.svg'
import feishuIcon from '@/assets/image/飞书.svg'
import wechatIcon from '@/assets/image/企业微信.svg'
import emailIcon from '@/assets/image/邮箱.svg'
import webhookIcon from '@/assets/image/url.svg'
import teamsIcon from '@/assets/image/通知.svg'

// 当前活动的标签页
const activeTab = ref('ssl-monitor')

// 统计数据
const statistics = ref({
  total: 0,
  alive: 0,
  dead: 0,
  sslExpired: 0,
  sslSoon: 0
})

// 定时巡检（精简显示与配置）
const scheduleConfig = reactive({ enabled: false, cronExpr: '', notifyEnabled: false, notifyRobotId: null })
const scheduleStatus = reactive({ is_running: false, next_run_time: '-' })
const scheduleDialogVisible = ref(false)
const scheduleForm = reactive({ enabled: false, cronExpr: '', notifyEnabled: false, notifyRobotId: null })
const scheduleFormLoading = ref(false)
const nextExecutionTimes = ref([])
const cronPreviewError = ref('')
const cronPreviewLoading = ref(false)
let cronPreviewTimer = null
const notifyRobotList = ref([]) // 通知机器人列表
const notifyRobotLoading = ref(false) // 加载状态

// 搜索表单
const searchForm = reactive({
  domain: '',
  tags: '',
  status: -1,
  isAlive: -1,
  pageNum: 1,
  pageSize: 10
})

// 列表数据
const tableData = ref([])
const total = ref(0)
const loading = ref(false)
const checkingAll = ref(false)

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('添加域名')
const formLoading = ref(false)

// 表单数据
const domainForm = reactive({
  id: null,
  domain: '',
  tags: '',
  remark: '',
  status: 1
})

// 表单验证规则
const formRules = {
  domain: [
    { required: true, message: '请输入域名', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:[0-9]{1,5})?$/, message: '请输入有效的域名格式', trigger: 'blur' }
  ]
}

const formRef = ref(null)

// 批量上传相关
const batchUploadDialogVisible = ref(false)
const uploadFileList = ref([])
const uploadLoading = ref(false)
const uploadRef = ref(null)
const showUploadTips = ref(false)

// 云导入相关
const cloudImportDialogVisible = ref(false)
const cloudImportLoading = ref(false)
const cloudImportForm = reactive({
  keyId: null,
  includeSubdomain: true,
  tags: '',
  remark: ''
})
const cloudImportFormRef = ref(null)
const aliyunKeyList = ref([]) // 阿里云密钥列表

// 状态选项
const statusOptions = [
  { label: '全部', value: -1 },
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

const aliveOptions = [
  { label: '全部', value: -1 },
  { label: '正常', value: 1 },
  { label: '异常', value: 0 }
]

// 获取统计信息
const fetchStatistics = async () => {
  try {
    const { data: res } = await domainApi.getStatistics()
    if (res.code === 200) {
      statistics.value = res.data
    }
  } catch (error) {
    console.error('获取统计信息失败:', error)
  }
}

// 获取域名列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = { ...searchForm }
    // 处理 -1 值，不传给后端
    if (params.status === -1) delete params.status
    if (params.isAlive === -1) delete params.isAlive

    const { data: res } = await domainApi.getList(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error) {
    ElMessage.error('获取域名列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  searchForm.pageNum = 1
  fetchList()
}

// 重置
const handleReset = () => {
  searchForm.domain = ''
  searchForm.tags = ''
  searchForm.status = -1
  searchForm.isAlive = -1
  searchForm.pageNum = 1
  handleSearch()
}

// 分页
const handlePageChange = (page) => {
  searchForm.pageNum = page
  fetchList()
}

const handleSizeChange = (size) => {
  searchForm.pageSize = size
  searchForm.pageNum = 1
  fetchList()
}

// 添加域名
const handleAdd = () => {
  dialogTitle.value = '添加域名'
  domainForm.id = null
  domainForm.domain = ''
  domainForm.tags = ''
  domainForm.remark = ''
  domainForm.status = 1
  dialogVisible.value = true
}

// 编辑域名
const handleEdit = (row) => {
  dialogTitle.value = '编辑域名'
  domainForm.id = row.id
  domainForm.domain = row.domain
  domainForm.tags = row.tags
  domainForm.remark = row.remark
  domainForm.status = row.status
  dialogVisible.value = true
}

// 删除域名
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该域名监控吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const { data: res } = await domainApi.delete(row.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      fetchList()
      fetchStatistics()
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error(error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    formLoading.value = true
    try {
      const data = { ...domainForm }
      let result
      if (data.id) {
        result = await domainApi.update(data)
      } else {
        delete data.id
        result = await domainApi.create(data)
      }

      const { data: res } = result
      if (res.code === 200) {
        ElMessage.success(data.id ? '更新成功' : '添加成功')
        dialogVisible.value = false
        fetchList()
        fetchStatistics()
      } else {
        ElMessage.error(res.message || '操作失败')
      }
    } catch (error) {
      ElMessage.error('操作失败')
      console.error(error)
    } finally {
      formLoading.value = false
    }
  })
}

// 检查单个域名
const handleCheck = async (row) => {
  const loadingMsg = ElMessage({
    message: '检查中...',
    type: 'info',
    duration: 0
  })

  try {
    const { data: res } = await domainApi.check(row.id)
    loadingMsg.close()

    if (res.code === 200) {
      ElMessage.success('检查完成')
      fetchList()
      fetchStatistics()
    } else {
      ElMessage.error(res.message || '检查失败')
    }
  } catch (error) {
    loadingMsg.close()
    ElMessage.error('检查失败')
    console.error(error)
  }
}

// 检查所有域名
const handleCheckAll = async () => {
  checkingAll.value = true
  const loadingMsg = ElMessage({
    message: '正在启动检查任务...',
    type: 'info',
    duration: 0
  })

  try {
    const { data: res } = await domainApi.checkAll()
    loadingMsg.close()

    if (res.code === 200) {
      const { total, message, status } = res.data

      // 显示后台任务启动成功的提示
      ElMessage.success(message || `检查任务已启动！共 ${total} 个域名正在后台检查...`)

      // 等待2秒后刷新列表，让用户看到初步结果
      setTimeout(() => {
        fetchList()
        fetchStatistics()
      }, 2000)

      // 再等待5秒后再次刷新，获取更多结果
      setTimeout(() => {
        fetchList()
        fetchStatistics()
      }, 5000)
    } else {
      ElMessage.error(res.message || '检查失败')
    }
  } catch (error) {
    loadingMsg.close()
    ElMessage.error('检查失败')
    console.error(error)
  } finally {
    checkingAll.value = false
  }
}

// 获取阿里云密钥列表
const fetchAliyunKeys = async () => {
  try {
    const { data: res } = await configApi.getKeyManageList({
      keyType: 1, // 1=阿里云
      page: 1,
      pageSize: 100
    })
    if (res.code === 200) {
      aliyunKeyList.value = res.data?.list || []
    }
  } catch (error) {
    console.error('获取阿里云密钥列表失败:', error)
  }
}

// 打开云导入对话框
const handleCloudImport = async () => {
  // 获取阿里云密钥列表
  await fetchAliyunKeys()

  if (aliyunKeyList.value.length === 0) {
    ElMessage.warning('请先在配置中心 > 密钥管理中添加阿里云密钥')
    return
  }

  // 重置表单
  cloudImportForm.keyId = null
  cloudImportForm.includeSubdomain = true
  cloudImportForm.tags = ''
  cloudImportForm.remark = ''

  cloudImportDialogVisible.value = true
}

// 提交云导入
const handleCloudImportSubmit = async () => {
  if (!cloudImportForm.keyId) {
    ElMessage.warning('请选择阿里云密钥')
    return
  }

  try {
    cloudImportLoading.value = true

    const { data: res } = await domainApi.importAliyunDns({
      keyId: cloudImportForm.keyId,
      includeSubdomain: cloudImportForm.includeSubdomain,
      tags: cloudImportForm.tags || undefined,
      remark: cloudImportForm.remark || undefined
    })

    if (res.code === 200) {
      const { total, message } = res.data
      ElMessage.success(message || `导入任务已启动，共 ${total} 个域名正在后台异步导入`)
      cloudImportDialogVisible.value = false

      // 等待3秒后刷新列表
      setTimeout(() => {
        fetchList()
        fetchStatistics()
      }, 3000)

      // 再等待8秒后再次刷新，获取更多结果
      setTimeout(() => {
        fetchList()
        fetchStatistics()
      }, 8000)
    } else {
      ElMessage.error(res.message || '导入失败')
    }
  } catch (error) {
    ElMessage.error('导入失败')
    console.error(error)
  } finally {
    cloudImportLoading.value = false
  }
}

// 获取定时巡检配置
const fetchScheduleConfig = async () => {
  try {
    const { data: res } = await domainApi.getScheduleConfig()
    if (res.code === 200 && res.data) {
      scheduleConfig.enabled = !!res.data.enabled
      scheduleConfig.cronExpr = res.data.cronExpr || ''
      scheduleConfig.notifyEnabled = !!res.data.notifyEnabled
      scheduleConfig.notifyRobotId = res.data.notifyRobotId || null
    }
  } catch (e) {
    console.error('获取定时巡检配置失败', e)
  } finally {
    // no-op
  }
}

// 获取定时巡检状态
const fetchScheduleStatus = async () => {
  try {
    const { data: res } = await domainApi.getScheduleStatus()
    if (res.code === 200 && res.data) {
      scheduleStatus.is_running = !!res.data.is_running
      scheduleStatus.next_run_time = res.data.next_run_time || '-'
    }
  } catch (e) {
    console.error('获取定时巡检状态失败', e)
  } finally {
    // no-op
  }
}

// 打开配置弹窗
const openScheduleDialog = async () => {
  scheduleForm.enabled = scheduleConfig.enabled
  scheduleForm.cronExpr = scheduleConfig.cronExpr
  scheduleForm.notifyEnabled = scheduleConfig.notifyEnabled
  scheduleForm.notifyRobotId = scheduleConfig.notifyRobotId
  nextExecutionTimes.value = []
  cronPreviewError.value = ''
  scheduleDialogVisible.value = true

  // 获取通知机器人列表
  await fetchNotifyRobots()

  // 如果有现有的Cron表达式，立即验证
  if (scheduleForm.cronExpr) {
    handleCronExprChange(scheduleForm.cronExpr)
  }
}

// 处理Cron表达式变化
const handleCronExprChange = (cronExpr) => {
  // 清除之前的定时器
  if (cronPreviewTimer) {
    clearTimeout(cronPreviewTimer)
  }

  // 清除之前的错误信息和预览时间
  cronPreviewError.value = ''
  nextExecutionTimes.value = []

  // 如果表达式为空，不进行验证
  if (!cronExpr || cronExpr.trim() === '') {
    return
  }

  // 防抖处理，500ms 后才执行验证
  cronPreviewTimer = setTimeout(() => {
    validateCronExpression(cronExpr.trim())
  }, 500)
}

// 验证Cron表达式并获取执行时间
const validateCronExpression = async (cronExpr) => {
  cronPreviewLoading.value = true
  cronPreviewError.value = ''

  try {
    // 使用 config API 获取下次执行时间
    const response = await domainApi.getNextExecutionTime(cronExpr)
    const responseData = response?.data || {}

    if (responseData.code === 200) {
      if (responseData.data?.next_execution_time) {
        nextExecutionTimes.value = [responseData.data.next_execution_time]
      } else {
        cronPreviewError.value = 'Cron表达式解析成功，但无法获取下次执行时间'
      }
    } else {
      cronPreviewError.value = responseData.message || 'Cron表达式格式错误'
    }
  } catch (error) {
    console.error('验证Cron表达式失败:', error)
    cronPreviewError.value = 'Cron表达式验证失败，请检查格式'
  } finally {
    cronPreviewLoading.value = false
  }
}

// 选择Cron模板
const selectCronTemplate = (cronExpr) => {
  scheduleForm.cronExpr = cronExpr
  handleCronExprChange(cronExpr)
}

// 获取机器人图标
const getRobotIcon = (type) => {
  const iconMap = {
    'dingtalk': dingtalkIcon,
    'feishu': feishuIcon,
    'wechat': wechatIcon,
    'email': emailIcon,
    'webhook': webhookIcon,
    'teams': teamsIcon
  }
  return iconMap[type] || ''
}

// 获取类型标签文本
const getTypeLabel = (type) => {
  const typeMap = {
    'dingtalk': '钉钉',
    'feishu': '飞书',
    'wechat': '企业微信',
    'email': '邮箱',
    'webhook': 'Webhook',
    'teams': 'Teams'
  }
  return typeMap[type] || '未知'
}

// 获取通知机器人列表
const fetchNotifyRobots = async () => {
  try {
    notifyRobotLoading.value = true
    const { data: res } = await domainApi.getNotifyRobots()
    if (res.code === 200) {
      notifyRobotList.value = res.data.list || []
    } else {
      console.error('获取通知机器人列表失败', res.message)
    }
  } catch (e) {
    console.error('获取通知机器人列表失败', e)
  } finally {
    notifyRobotLoading.value = false
  }
}

// 保存配置
const handleSaveSchedule = async () => {
  try {
    if (scheduleForm.enabled && (!scheduleForm.cronExpr || scheduleForm.cronExpr.trim() === '')) {
      ElMessage.warning('启用定时巡检时必须填写Cron表达式')
      return
    }
    if (scheduleForm.notifyEnabled && !scheduleForm.notifyRobotId) {
      ElMessage.warning('开启通知时必须选择通知机器人')
      return
    }
    scheduleFormLoading.value = true
    const { data: res } = await domainApi.updateScheduleConfig({
      enabled: scheduleForm.enabled,
      cronExpr: scheduleForm.cronExpr,
      notifyEnabled: scheduleForm.notifyEnabled,
      notifyRobotId: scheduleForm.notifyRobotId || null
    })
    if (res.code === 200) {
      ElMessage.success(res.message || '配置更新成功')
      scheduleDialogVisible.value = false
      await fetchScheduleConfig()
      await fetchScheduleStatus()
    } else {
      ElMessage.error(res.message || '配置更新失败')
    }
  } catch (e) {
    ElMessage.error('配置更新失败')
    console.error(e)
  } finally {
    scheduleFormLoading.value = false
  }
}

// SSL证书剩余天数样式
const getSslDaysClass = (daysLeft) => {
  if (daysLeft <= 0) return 'ssl-expired'
  if (daysLeft <= 7) return 'ssl-critical'
  if (daysLeft <= 30) return 'ssl-warning'
  return 'ssl-normal'
}

// 存活状态类型
const getAliveType = (isAlive) => {
  return isAlive === 1 ? 'success' : 'danger'
}

// 启用状态类型
const getStatusType = (status) => {
  return status === 1 ? 'success' : 'info'
}

// 检查是否有SSL信息
const hasSslInfo = (row) => {
  // 检查 sslExpireAt 是否有效
  if (row.sslExpireAt) {
    // 处理对象格式
    if (typeof row.sslExpireAt === 'object' && row.sslExpireAt.hasOwnProperty('Valid')) {
      if (row.sslExpireAt.Valid && row.sslExpireAt.Time && !row.sslExpireAt.Time.startsWith('0001-01-01')) {
        return true
      }
    }
    // 处理字符串格式
    if (typeof row.sslExpireAt === 'string' && row.sslExpireAt && !row.sslExpireAt.startsWith('0001-01-01')) {
      return true
    }
  }
  return false
}

// 格式化日期 - 处理后端返回的 {Time: "", Valid: false} 格式
const formatDate = (dateObj) => {
  if (!dateObj) return '-'

  // 处理后端返回的对象格式 {Time: "", Valid: false}
  if (typeof dateObj === 'object' && dateObj.hasOwnProperty('Valid')) {
    if (!dateObj.Valid || !dateObj.Time) return '-'
    const dateStr = dateObj.Time
    // 检查是否是无效日期
    if (dateStr.startsWith('0001-01-01')) return '-'
    return new Date(dateStr).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  // 处理字符串格式
  if (typeof dateObj === 'string') {
    if (!dateObj || dateObj.startsWith('0001-01-01')) return '-'
    return new Date(dateObj).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  return '-'
}

// 查询即将过期的SSL证书
const handleQueryExpiring = async (days) => {
  loading.value = true
  try {
    const { data: res } = await domainApi.getSslExpiring(days)

    if (res.code === 200) {
      const { total: totalCount, domains } = res.data
      tableData.value = domains || []
      total.value = totalCount || 0

      if (totalCount > 0) {
        ElMessage.success(`查询到 ${totalCount} 个域名的证书将在 ${days} 天内过期`)
      } else {
        ElMessage.info(`未查询到 ${days} 天内即将过期的证书`)
      }
    } else {
      ElMessage.error(res.message || '查询失败')
    }
  } catch (error) {
    ElMessage.error('查询失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 查询异常域名
const handleQueryAbnormal = async () => {
  searchForm.isAlive = 0 // 0表示异常
  searchForm.pageNum = 1
  await fetchList()

  if (tableData.value.length > 0) {
    ElMessage.success(`查询到 ${total.value} 个异常域名`)
  } else {
    ElMessage.info('未查询到异常域名')
  }
}

// 打开批量上传对话框
const openBatchUploadDialog = () => {
  uploadFileList.value = []
  showUploadTips.value = false
  batchUploadDialogVisible.value = true
}

// 文件上传前验证
const handleBeforeUpload = (file) => {
  const isText = file.type === 'text/plain' || file.name.endsWith('.txt')
  if (!isText) {
    ElMessage.error('只能上传 .txt 文本文件')
    return false
  }
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB')
    return false
  }
  return true
}

// 文件列表变化
const handleFileChange = (file, fileList) => {
  uploadFileList.value = fileList
}

// 移除文件
const handleFileRemove = (file, fileList) => {
  uploadFileList.value = fileList
}

// 提交批量上传
const handleBatchUpload = async () => {
  if (uploadFileList.value.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }

  uploadLoading.value = true
  try {
    const file = uploadFileList.value[0].raw
    const { data: res } = await domainApi.batchUpload(file)

    if (res.code === 200) {
      const { total, message } = res.data
      ElMessage.success(message || `批量上传成功，共解析 ${total} 个域名`)
      batchUploadDialogVisible.value = false
      uploadFileList.value = []

      // 等待2秒后刷新列表
      setTimeout(() => {
        fetchList()
        fetchStatistics()
      }, 2000)

      // 再等待5秒后再次刷新
      setTimeout(() => {
        fetchList()
        fetchStatistics()
      }, 5000)
    } else {
      ElMessage.error(res.message || '批量上传失败')
    }
  } catch (error) {
    ElMessage.error('批量上传失败')
    console.error(error)
  } finally {
    uploadLoading.value = false
  }
}

onMounted(() => {
  fetchStatistics()
  fetchList()
  fetchScheduleConfig()
  fetchScheduleStatus()
})
</script>

<template>
  <div class="ssl-monitor-container">

    <!-- 统计卡片 -->
    <div class="statistics-cards">
      <div class="stat-card stat-total">
        <div class="stat-icon">
          <el-icon :size="32"><Platform /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">域名总数</div>
          <div class="stat-value">{{ statistics.total }}</div>
        </div>
      </div>

      <div class="stat-card stat-alive">
        <div class="stat-icon">
          <el-icon :size="32"><SuccessFilled /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">正常域名</div>
          <div class="stat-value">{{ statistics.alive }}</div>
        </div>
      </div>

      <div class="stat-card stat-dead">
        <div class="stat-icon">
          <el-icon :size="32"><CircleCloseFilled /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">异常域名</div>
          <div class="stat-value">{{ statistics.dead }}</div>
        </div>
      </div>

      <div class="stat-card stat-soon">
        <div class="stat-icon">
          <el-icon :size="32"><WarningFilled /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">即将过期</div>
          <div class="stat-value">{{ statistics.sslSoon }}</div>
        </div>
      </div>

      <div class="stat-card stat-expired">
        <div class="stat-icon">
          <el-icon :size="32"><RemoveFilled /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">已过期</div>
          <div class="stat-value">{{ statistics.sslExpired }}</div>
        </div>
      </div>
    </div>

    <!-- 域名监控卡片 -->
    <el-card shadow="hover" class="ssl-monitor-card">
      <template #header>
        <div class="card-header">
          <span class="title">域名监控</span>
        </div>
      </template>

      <!-- 标签页 -->
      <el-tabs v-model="activeTab" @tab-change="handleTabChange" class="ssl-tabs">
        <!-- SSL证书监控标签页 -->
        <el-tab-pane name="ssl-monitor">
          <template #label>
            <span class="tab-label">
              <el-icon><Monitor /></el-icon>
              SSL证书监控
            </span>
          </template>

      <!-- 搜索区域 -->
      <div class="search-section">
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="域名">
            <el-input
              v-model="searchForm.domain"
              placeholder="请输入域名"
              size="small"
              clearable
              @keyup.enter="handleSearch"
            />
          </el-form-item>

          <el-form-item label="标签">
            <el-input
              v-model="searchForm.tags"
              placeholder="请输入标签"
              size="small"
              clearable
            />
          </el-form-item>

          <el-form-item label="状态">
            <el-select v-model="searchForm.status" size="small" placeholder="请选择" style="width: 120px">
              <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>

          <el-form-item label="存活状态">
            <el-select v-model="searchForm.isAlive" size="small" placeholder="请选择" style="width: 120px">
              <el-option v-for="item in aliveOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" size="small" @click="handleSearch" class="modern-btn primary-btn">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button size="small" @click="handleReset" class="modern-btn reset-btn">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 操作按钮 -->
      <div class="action-section">
        <el-button
          type="success"
          size="small"
          @click="handleAdd"
          v-authority="['monitor:domain:add']"
          class="modern-btn success-btn"
        >
          <el-icon><Plus /></el-icon>
          添加域名
        </el-button>

        <el-button
          type="primary"
          size="small"
          @click="openBatchUploadDialog"
          v-authority="['monitor:domain:add_all']"
          class="modern-btn primary-btn"
        >
          <el-icon><Upload /></el-icon>
          批量上传
        </el-button>
        <el-button
          type="primary"
          size="small"
           v-authority="['monitor:domain:cloud_import']"
          @click="handleCloudImport"
          :loading="cloudImportLoading"
          class="modern-btn primary-btn"
        >
          <el-icon><CloudUpload /></el-icon>
          云导入
        </el-button>


        <el-divider direction="vertical" style="height: 24px; margin: 0 12px" />

        <el-button
          type="primary"
          size="small"
          v-authority="['monitor:domain:auto_inspection']"
          class="modern-btn primary-btn"
          @click="openScheduleDialog"
        >
          自动巡检
        </el-button>
        <el-button
          type="warning"
          size="small"
          @click="handleQueryExpiring(30)"
          class="modern-btn warning-btn"
          plain
        >
          <el-icon><WarningFilled /></el-icon>
          30天
        </el-button>
        <el-button
          type="danger"
          size="small"
          @click="handleQueryAbnormal"
          class="modern-btn danger-btn"
          plain
        >
          <el-icon><CircleCloseFilled /></el-icon>
          异常域名
        </el-button>

      </div>

      <!-- 域名列表 -->
      <div class="table-section">
        <el-table
          :data="tableData"
          v-loading="loading"
          stripe
          border
          style="width: 100%"
          :header-cell-style="{ background: '#eef1f6', color: '#606266' }"
        >
          <el-table-column prop="id" label="ID" width="60" align="center" />

          <el-table-column prop="domain" label="域名" width="220">
            <template #default="{ row }">
              <div style="line-height: 1.4">
                <div style="font-weight: 500; font-size: 13px">{{ row.domain }}</div>
                <el-tag v-if="row.tags" size="small" type="info" style="margin-top: 4px; font-size: 11px">
                  {{ row.tags }}
                </el-tag>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="getAliveType(row.isAlive)" size="small" style="font-size: 11px">
                {{ row.isAlive === 1 ? '正常' : '异常' }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="statusCode" label="状态码" width="70" align="center">
            <template #default="{ row }">
              <span style="font-size: 12px">{{ row.statusCode || '-' }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="responseTime" label="响应" width="70" align="center">
            <template #default="{ row }">
              <span style="font-size: 11px">{{ row.responseTime > 0 ? row.responseTime + 'ms' : '-' }}</span>
            </template>
          </el-table-column>



          <el-table-column label="SSL证书" width="160">
            <template #default="{ row }">
              <div v-if="hasSslInfo(row)" style="font-size: 11px; line-height: 1.6">
                <div style="color: #606266">{{ formatDate(row.sslExpireAt) }}</div>
                <el-tag :class="getSslDaysClass(row.sslDaysLeft)" size="small" style="margin-top: 4px; font-size: 10px">
                  剩余 {{ row.sslDaysLeft }} 天
                </el-tag>
              </div>
              <span v-else style="color: #c0c4cc; font-size: 11px">-</span>
            </template>
          </el-table-column>

          <el-table-column prop="lastCheckAt" label="最后检查" width="135">
            <template #default="{ row }">
              <span style="font-size: 11px; color: #606266">{{ formatDate(row.lastCheckAt) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="启用" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small" style="font-size: 11px">
                {{ row.status === 1 ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="remark" label="备注" width="180" show-overflow-tooltip>
            <template #default="{ row }">
              <span style="font-size: 11px; color: #606266">{{ row.remark || '-' }}</span>
            </template>
          </el-table-column>

          <el-table-column label="证书提供商"  show-overflow-tooltip>
            <template #default="{ row }">
              <span v-if="row.sslIssuer" style="font-size: 11px; color: #606266">{{ row.sslIssuer }}</span>
              <span v-else style="color: #c0c4cc; font-size: 11px">-</span>
            </template>
          </el-table-column>

          <el-table-column label="操作"  fixed="right" align="center">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-tooltip content="检查" placement="top">
                  <el-button
                    type="primary"
                    :icon="View"
                    size="small"
                    @click="handleCheck(row)"
                    v-authority="['monitor:domain:ops_inspection']"
                  >检查</el-button>
                </el-tooltip>
                <el-tooltip content="编辑" placement="top">
                  <el-button
                    type="warning"
                    :icon="Edit"
                    size="small"
                    @click="handleEdit(row)"
                    v-authority="['monitor:domain:edit']"
                  >编辑</el-button>
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button
                    type="danger"
                    :icon="Delete"
                    size="small"
                    @click="handleDelete(row)"
                    v-authority="['monitor:domain:delete']"
                  >删除</el-button>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页区域 -->
      <div class="pagination-section">
        <el-pagination
          v-model:current-page="searchForm.pageNum"
          v-model:page-size="searchForm.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
        </el-tab-pane>

        <!-- SSL证书申请标签页 -->
        <el-tab-pane name="ssl-cert-apply">
          <template #label>
            <span class="tab-label">
              <el-icon><Document /></el-icon>
              SSL证书申请
            </span>
          </template>

          <div class="tab-content">
            <ssl-cert-apply />
          </div>
        </el-tab-pane>

        <!-- 云厂商配置标签页 -->
        <el-tab-pane name="aliyun-config">
          <template #label>
            <span class="tab-label">
              <el-icon><Setting /></el-icon>
              DNS云厂商配置
            </span>
          </template>

          <div class="tab-content">
            <aliyun-config />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="domainForm"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="域名" prop="domain">
          <el-input v-model="domainForm.domain" placeholder="例如: www.example.com" />
        </el-form-item>

        <el-form-item label="标签">
          <el-input v-model="domainForm.tags" placeholder="多个标签用逗号分隔" />
        </el-form-item>

        <el-form-item label="备注">
          <el-input
            v-model="domainForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
          />
        </el-form-item>

        <el-form-item label="状态">
          <el-radio-group v-model="domainForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="formLoading">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 定时巡检配置对话框 -->
    <el-dialog
      v-model="scheduleDialogVisible"
      title="配置定时巡检"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form label-width="120px">
        <el-form-item label="是否启用">
          <el-switch v-model="scheduleForm.enabled" />
        </el-form-item>

        <el-form-item label="Cron表达式" prop="cronExpr">
          <el-input
            v-model="scheduleForm.cronExpr"
            placeholder="请输入Cron表达式"
            @input="handleCronExprChange"
          >
            <template #append>
              <el-dropdown @command="selectCronTemplate">
                <el-button>
                  常用模板
                  <el-icon class="el-icon--right"><arrow-down /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="0 */2 * * *">每2小时</el-dropdown-item>
                    <el-dropdown-item command="0 2 * * *">每天凌晨2点</el-dropdown-item>
                    <el-dropdown-item command="0 0 */6 * *">每6小时</el-dropdown-item>
                    <el-dropdown-item command="30 1 * * *">每天凌晨1:30</el-dropdown-item>
                    <el-dropdown-item command="0 */4 * * *">每4小时</el-dropdown-item>
                    <el-dropdown-item command="0 0 * * *">每天午夜</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-input>

          <!-- Cron表达式预览 -->
          <div v-if="scheduleForm.cronExpr && nextExecutionTimes.length > 0" style="margin-top: 8px;">
            <span style="margin-left: 10px; color: #67C23A; font-weight: 500;">
              下次执行时间: {{ nextExecutionTimes[0] }}
            </span>
          </div>
          <div v-else-if="scheduleForm.cronExpr && cronPreviewError" style="margin-top: 8px;">
            <span style="margin-left: 10px; color: #F56C6C; font-weight: 500;">
              {{ cronPreviewError }}
            </span>
          </div>
          <div v-else style="margin-top: 8px;">
            <span style="margin-left: 10px; color: #909399; font-size: 12px;">
              格式: 分 时 日 月 周，例如: 0 2 * * * 表示每天凌晨2点执行
            </span>
          </div>
        </el-form-item>

        <el-form-item label="是否通知">
          <el-switch v-model="scheduleForm.notifyEnabled" />
        </el-form-item>

        <el-form-item v-if="scheduleForm.notifyEnabled" label="通知机器人" prop="notifyRobotId">
          <el-select
            v-model="scheduleForm.notifyRobotId"
            placeholder="请选择通知机器人"
            style="width: 100%"
            :loading="notifyRobotLoading"
          >
            <el-option
              v-for="robot in notifyRobotList"
              :key="robot.id"
              :label="`${robot.name}${robot.status === 1 ? '(启用)' : '(禁用)'}`"
              :value="robot.id"
              :disabled="robot.status === 0"
            >
              <div style="display: flex; justify-content: space-between; align-items: center; width: 100%;">
                <div style="display: flex; align-items: center; gap: 8px; flex: 1;">
                  <img :src="getRobotIcon(robot.type)" :alt="getTypeLabel(robot.type)" style="width: 20px; height: 20px;" />
                  <span>{{ robot.name }}</span>
                </div>
                <div style="display: flex; align-items: center; gap: 8px;">
                  <span style="color: #8492a6; font-size: 12px;">{{ getTypeLabel(robot.type) }}</span>
                  <el-tag :type="robot.status === 1 ? 'success' : 'info'" size="small" style="font-size: 12px;">
                    {{ robot.status === 1 ? '启用' : '禁用' }}
                  </el-tag>
                </div>
              </div>
            </el-option>
          </el-select>
          <div v-if="notifyRobotList.length === 0 && !notifyRobotLoading" style="margin-top: 8px; color: #F56C6C; font-size: 12px;">
            暂无可用的通知机器人，请先在系统设置中配置
          </div>
          <div v-else-if="notifyRobotList.filter(r => r.status === 1).length === 0 && !notifyRobotLoading" style="margin-top: 8px; color: #E6A23C; font-size: 12px;">
            没有启用的通知机器人，请先启用相关机器人
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="scheduleDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="scheduleFormLoading" @click="handleSaveSchedule">保存</el-button>
      </template>
    </el-dialog>

    <!-- 批量上传对话框 -->
    <el-dialog
      v-model="batchUploadDialogVisible"
      title="批量上传域名"
      width="600px"
      :close-on-click-modal="false"
    >
      <div style="margin-bottom: 20px">
        <div class="upload-tips-header">
          <el-icon
            class="tips-icon"
            :size="20"
            @click="showUploadTips = !showUploadTips"
            style="cursor: pointer;"
          >
            <WarningFilled />
          </el-icon>
          <span
            class="tips-text"
            @click="showUploadTips = !showUploadTips"
            style="cursor: pointer;"
          >
            文件格式说明
          </span>
        </div>

        <transition name="el-fade-in">
          <div v-show="showUploadTips" class="upload-tips-content">
            <p>1. 仅支持 .txt 文本文件，每行一个域名</p>
            <p>2. 自动去除 http://、https://、www. 前缀和路径</p>
            <p>3. 自动过滤重复域名和以 # 开头的注释行</p>
            <p>4. 已存在的域名会自动跳过</p>
            <p>5. 后台异步创建，最多同时创建 5 个</p>
          </div>
        </transition>
      </div>

      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :limit="1"
        :file-list="uploadFileList"
        :before-upload="handleBeforeUpload"
        :on-change="handleFileChange"
        :on-remove="handleFileRemove"
        drag
        accept=".txt"
      >
        <el-icon class="el-icon--upload"><Upload /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            仅支持 .txt 文件，文件大小不超过 10MB
          </div>
        </template>
      </el-upload>

      <template #footer>
        <el-button @click="batchUploadDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="uploadLoading"
          @click="handleBatchUpload"
        >
          开始上传
        </el-button>
      </template>
    </el-dialog>

    <!-- 云导入对话框 -->
    <el-dialog
      v-model="cloudImportDialogVisible"
      title="阿里云DNS域名导入"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="cloudImportFormRef"
        :model="cloudImportForm"
        label-width="120px"
      >
        <el-form-item label="阿里云密钥" required>
          <el-select
            v-model="cloudImportForm.keyId"
            placeholder="请选择阿里云密钥"
            style="width: 100%"
          >
            <el-option
              v-for="key in aliyunKeyList"
              :key="key.id"
              :label="key.keyName || '未命名'"
              :value="key.id"
            />
          </el-select>
          <div style="margin-top: 8px; color: #909399; font-size: 12px;">
            在 <router-link to="/configcenter/keymanage" style="color: #409EFF;">配置中心 > 密钥管理</router-link> 中管理密钥
          </div>
        </el-form-item>

        <el-form-item label="导入模式">
          <el-radio-group v-model="cloudImportForm.includeSubdomain">
            <el-radio :label="true">
              导入所有子域名 ⭐
              <span style="color: #909399; font-size: 12px; margin-left: 8px;">
                (如: dding.net, www.dding.net, api.dding.net 等)
              </span>
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="统一标签">
          <el-input
            v-model="cloudImportForm.tags"
            placeholder="为导入的域名添加统一标签，如：生产环境"
          />
        </el-form-item>

        <el-form-item label="统一备注">
          <el-input
            v-model="cloudImportForm.remark"
            type="textarea"
            :rows="3"
            placeholder="为导入的域名添加统一备注，如：从阿里云解析导入"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="cloudImportDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="cloudImportLoading"
          @click="handleCloudImportSubmit"
        >
          开始导入
        </el-button>
      </template>
    </el-dialog>

  </div>
</template>

<style scoped>
.ssl-monitor-container {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}
.schedule-inline-text { margin-left: 8px; color: #606266; }

/* 统计卡片区域 */
.statistics-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 12px;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 6px;
  font-weight: 500;
}

.stat-value {
  font-size: 26px;
  font-weight: bold;
  line-height: 1;
}

/* 不同类型统计卡片颜色 */
.stat-card.stat-total .stat-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-card.stat-total .stat-value {
  color: #667eea;
}

.stat-card.stat-alive .stat-icon {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  color: white;
}

.stat-card.stat-alive .stat-value {
  color: #67c23a;
}

.stat-card.stat-dead .stat-icon {
  background: linear-gradient(135deg, #f56c6c 0%, #f78989 100%);
  color: white;
}

.stat-card.stat-dead .stat-value {
  color: #f56c6c;
}

.stat-card.stat-soon .stat-icon {
  background: linear-gradient(135deg, #e6a23c 0%, #eebe77 100%);
  color: white;
}

.stat-card.stat-soon .stat-value {
  color: #e6a23c;
}

.stat-card.stat-expired .stat-icon {
  background: linear-gradient(135deg, #909399 0%, #b1b3b8 100%);
  color: white;
}

.stat-card.stat-expired .stat-value {
  color: #909399;
}

/* 主卡片 */
.ssl-monitor-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  min-height: calc(100vh - 100px);
}

.ssl-monitor-card :deep(.el-card__body) {
  min-height: 500px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text; /* WebKit 浏览器支持 */
  background-clip: text; /* 标准写法 */
  -webkit-text-fill-color: transparent; /* WebKit 浏览器支持 */
  color: transparent; /* 非 WebKit 浏览器支持 */
}

.search-section {
  margin-bottom: 16px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form :deep(.el-form-item) {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form :deep(.el-form-item__label) {
  color: #606266;
  font-weight: 500;
}

.action-section {
  margin-bottom: 16px;
}

.table-section {
  margin-top: 0;
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

:deep(.el-table) {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  font-size: 12px;
}

:deep(.el-table__header) {
  background: rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-table__header th) {
  background: rgba(102, 126, 234, 0.1) !important;
  color: #2c3e50 !important;
  font-weight: 600;
  border: none;
  padding: 8px 0;
  font-size: 13px;
}

:deep(.el-table__body tr:hover > td) {
  background-color: rgba(102, 126, 234, 0.08) !important;
}

:deep(.el-table td) {
  border: none;
  padding: 6px 0;
}

:deep(.el-table .cell) {
  padding: 0 8px;
  line-height: 1.4;
}

:deep(.el-table::before) {
  display: none;
}

:deep(.el-table--border::after) {
  display: none;
}

/* 确保操作列固定在最上层 */
:deep(.el-table__fixed-right) {
  z-index: 3 !important;
}

:deep(.el-table__fixed-right-patch) {
  z-index: 3 !important;
}

.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.operation-buttons .el-button {
  transition: all 0.3s ease;
}

.operation-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

/* SSL天数标签颜色 */
.ssl-normal {
  background-color: #f0f9ff;
  color: #67c23a;
  border-color: #c2e7b0;
}

.ssl-warning {
  background-color: #fdf6ec;
  color: #e6a23c;
  border-color: #f5dab1;
}

.ssl-critical {
  background-color: #fef0f0;
  color: #f56c6c;
  border-color: #fbc4c4;
}

.ssl-expired {
  background-color: #f0f0f0;
  color: #909399;
  border-color: #d3d4d6;
}

:deep(.modern-btn) {
  border-radius: 8px;
  padding: 8px 20px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

:deep(.modern-btn:hover) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
}

:deep(.primary-btn) {
  background: linear-gradient(45deg, #409EFF, #66B3FF);
  color: white;
}

:deep(.reset-btn) {
  background: linear-gradient(45deg, #E6A23C, #EEBE77);
  color: white;
}

:deep(.success-btn) {
  background: linear-gradient(45deg, #67C23A, #85CE61);
  color: white;
}

:deep(.warning-btn) {
  background: linear-gradient(45deg, #E6A23C, #EEBE77);
  color: white;
}

:deep(.warning-btn.is-plain) {
  background: rgba(230, 162, 60, 0.1);
  color: #E6A23C;
  border: 1px solid rgba(230, 162, 60, 0.3);
}

:deep(.warning-btn.is-plain:hover) {
  background: linear-gradient(45deg, #E6A23C, #EEBE77);
  color: white;
  border-color: transparent;
}

/* 批量上传提示样式 */
.upload-tips-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 0;
}

.tips-icon {
  color: #E6A23C;
  transition: transform 0.3s ease;
}

.tips-icon:hover {
  transform: scale(1.1);
}

.tips-text {
  color: #E6A23C;
  font-weight: 500;
  font-size: 14px;
  user-select: none;
}

.upload-tips-content {
  margin-top: 10px;
  padding: 15px;
  background: rgba(230, 162, 60, 0.1);
  border-left: 3px solid #E6A23C;
  border-radius: 4px;
}

.upload-tips-content p {
  margin: 8px 0;
  color: #E6A23C;
  font-size: 13px;
  line-height: 1.6;
}

.upload-tips-content p:first-child {
  margin-top: 0;
}

.upload-tips-content p:last-child {
  margin-bottom: 0;
}

/* 响应式调整 */
@media (max-width: 1400px) {
  .statistics-cards {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }
}

@media (max-width: 768px) {
  .ssl-monitor-container {
    padding: 10px;
  }

  .statistics-cards {
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
    margin-bottom: 15px;
  }

  .stat-card {
    padding: 12px;
    gap: 10px;
  }

  .stat-icon {
    width: 42px;
    height: 42px;
  }

  .stat-icon :deep(.el-icon) {
    font-size: 22px !important;
  }

  .stat-label {
    font-size: 12px;
  }

  .stat-value {
    font-size: 20px;
  }
}

/* 标签页样式 */
.ssl-tabs {
  margin-top: 0;
}

.ssl-tabs :deep(.el-tabs__header) {
  margin-bottom: 20px;
}

.ssl-tabs :deep(.el-tabs__item) {
  font-weight: 500;
  color: #606266;
  padding: 0 20px;
}

.ssl-tabs :deep(.el-tabs__item.is-active) {
  color: #409EFF;
  font-weight: 600;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 6px;
}

.tab-label .el-icon {
  font-size: 16px;
}

.tab-content {
  padding: 0;
}

/* SSL证书申请表单样式 */
.cert-list-card {
  margin-bottom: 20px;
  border: 1px solid rgba(103, 126, 234, 0.1);
  border-radius: 12px;
}

.cert-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cert-list-header .header-actions {
  display: flex;
  gap: 8px;
}

.form-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
}

.form-tip {
  margin-top: 6px;
  padding-left: 2px;
}

.form-tip .el-text {
  line-height: 1.4;
}

/* 部署日志样式 */
.deploy-logs {
  max-height: 300px;
  overflow-y: auto;
  background: #f5f7fa;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 12px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.log-item {
  padding: 4px 0;
  color: #606266;
  line-height: 1.6;
  word-break: break-all;
}

.log-item:not(:last-child) {
  border-bottom: 1px dashed #e4e7ed;
}
</style>
