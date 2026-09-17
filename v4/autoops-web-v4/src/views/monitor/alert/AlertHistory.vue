<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Download,
  Check,
  Reading,
  Delete,
  DeleteFilled
} from '@element-plus/icons-vue'
import { recordApi } from '@/api/alert'
import systemApi from '@/api/system'
import dayjs from 'dayjs'

// 定义事件
const emit = defineEmits(['refresh-stats'])

// 搜索表单
const searchForm = reactive({
  status: '',
  severity: '',
  ruleId: '',
  businessGroup: '',
  startTime: null,
  endTime: null,
  keyword: '',
  pageNum: 1,
  pageSize: 10
})

// 列表数据
const tableData = ref([])
const total = ref(0)
const loading = ref(false)

// 行选择
const selectedRows = ref([])

// 详情抽屉
const drawerVisible = ref(false)
const detailData = ref(null)

// 统一处理对话框
const processDialogVisible = ref(false)
const currentRecordId = ref(null)
const currentRecord = ref(null)
const processLoading = ref(false)
const userOptions = ref([])
const processForm = reactive({
  action: '',
  assigneeId: undefined,
  handleNote: ''
})

const actionOptions = [
  { label: '认领', value: 'acknowledge' },
  { label: '转派', value: 'assign' },
  { label: '处理', value: 'handle' }
]

// 下拉选项
const statusOptions = [
  { label: '全部', value: '' },
  { label: '告警中', value: 'firing' },
  { label: '已恢复', value: 'resolved' },
  { label: '已认领', value: 'acknowledged' }
]

const severityOptions = [
  { label: '全部', value: '' },
  { label: '严重', value: 1 },
  { label: '警告', value: 2 },
  { label: '信息', value: 3 }
]

const ruleList = ref([])
const businessGroupList = ref([])

const normalizeStatus = (status) => {
  if (status === null || status === undefined) return ''
  return String(status).trim().toLowerCase()
}

const isResolvedStatus = (status) => {
  const normalized = normalizeStatus(status)
  return normalized === 'resolved' || normalized === '已恢复'
}

const isAcknowledgedStatus = (status) => {
  const normalized = normalizeStatus(status)
  return normalized === 'acknowledged' || normalized === '已认领'
}

const canProcessRecord = (row) => {
  return !isResolvedStatus(row?.status)
}

const isAssignAction = computed(() => processForm.action === 'assign')
const isHandleAction = computed(() => processForm.action === 'handle')
const shouldShowAssignee = computed(() => isAssignAction.value)

// 获取告警记录列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = { ...searchForm, pageNum: searchForm.pageNum, pageSize: searchForm.pageSize }

    // 处理时间范围
    if (searchForm.startTime) {
      params.startTime = dayjs(searchForm.startTime).format('YYYY-MM-DD HH:mm:ss')
    }
    if (searchForm.endTime) {
      params.endTime = dayjs(searchForm.endTime).format('YYYY-MM-DD HH:mm:ss')
    }

    const { data: res } = await recordApi.getList(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error) {
    ElMessage.error('获取告警记录列表失败')
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
  searchForm.status = ''
  searchForm.severity = ''
  searchForm.ruleId = ''
  searchForm.businessGroup = ''
  searchForm.startTime = null
  searchForm.endTime = null
  searchForm.keyword = ''
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

// 查看详情
const handleViewDetail = async (row) => {
  try {
    const { data: res } = await recordApi.getDetail(row.id)
    if (res.code === 200) {
      detailData.value = {
        ...row,
        ...(res.data || {})
      }
      drawerVisible.value = true
    } else {
      ElMessage.error(res.message || '获取详情失败')
    }
  } catch (error) {
    ElMessage.error('获取详情失败')
    console.error(error)
  }
}

const getAlertValue = (record) => {
  return record?.currentValue ?? record?.value ?? '-'
}

const getHandleTime = (record) => {
  return record?.handleTime || record?.handledAt || '-'
}

const getHandleNote = (record) => {
  return record?.handleNote || record?.remark || '-'
}

const parseLabels = (labels) => {
  if (!labels) return {}
  if (typeof labels === 'object') return labels

  try {
    return JSON.parse(labels)
  } catch (error) {
    console.error('解析 labels 失败:', error)
    return {}
  }
}

const fetchUserOptions = async () => {
  try {
    const response = await systemApi.queryAdminList({
      pageNum: 1,
      pageSize: 1000,
      status: '1'
    })
    if (response.data?.code === 200) {
      userOptions.value = (response.data.data?.list || []).map(user => ({
        id: user.id,
        name: user.nickname || user.username
      })).filter(user => user.id && user.name)
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  }
}

const resetProcessForm = () => {
  processForm.action = ''
  processForm.assigneeId = undefined
  processForm.handleNote = ''
}

const openProcessDialog = async (row) => {
  currentRecordId.value = row.id
  currentRecord.value = row
  resetProcessForm()
  processDialogVisible.value = true

  if (userOptions.value.length === 0) {
    await fetchUserOptions()
  }
}

const buildProcessPayload = () => {
  const payload = {
    action: processForm.action
  }

  if (processForm.assigneeId !== undefined && processForm.assigneeId !== null && processForm.assigneeId !== '') {
    payload.assigneeId = processForm.assigneeId
  }

  if (processForm.handleNote.trim()) {
    payload.handleNote = processForm.handleNote.trim()
  }

  return payload
}

// 统一处理
const handleProcess = async () => {
  if (!processForm.action) {
    ElMessage.warning('请选择操作类型')
    return
  }

  if (isAssignAction.value && !processForm.assigneeId) {
    ElMessage.warning('请选择用户')
    return
  }

  if (isHandleAction.value && !processForm.handleNote.trim()) {
    ElMessage.warning('请输入处理步骤')
    return
  }

  processLoading.value = true
  try {
    const { data: res } = await recordApi.process(currentRecordId.value, buildProcessPayload())
    if (res.code === 200) {
      ElMessage.success('操作成功')
      processDialogVisible.value = false
      resetProcessForm()
      fetchList()
      emit('refresh-stats')
    } else {
      ElMessage.error(res.message || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败')
    console.error(error)
  } finally {
    processLoading.value = false
  }
}

// 导出
const handleExport = () => {
  ElMessage.info('导出功能开发中...')
}

// 刷新
const handleRefresh = () => {
  fetchList()
}

// 行选择
const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

// 删除单条告警记录
const handleDelete = async (row) => {
  try {
    const { data: res } = await recordApi.delete(row.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      fetchList()
      emit('refresh-stats')
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (error) {
    ElMessage.error('删除失败')
    console.error(error)
  }
}

// 批量删除告警记录
const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择要删除的告警记录')
    return
  }

  try {
    const ids = selectedRows.value.map(row => row.id)
    const { data: res } = await recordApi.batchDelete(ids)
    if (res.code === 200) {
      ElMessage.success('批量删除成功')
      selectedRows.value = []
      fetchList()
      emit('refresh-stats')
    } else {
      ElMessage.error(res.message || '批量删除失败')
    }
  } catch (error) {
    ElMessage.error('批量删除失败')
    console.error(error)
  }
}

// 清空所有告警记录
const handleClearAll = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有告警记录吗？此操作不可恢复', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const { data: res } = await recordApi.clearAll()
    if (res.code === 200) {
      ElMessage.success('清空成功')
      selectedRows.value = []
      fetchList()
      emit('refresh-stats')
    } else {
      ElMessage.error(res.message || '清空失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清空失败')
      console.error(error)
    }
  }
}

// 获取状态标签类型
const getStatusType = (status) => {
  const normalized = normalizeStatus(status)
  const typeMap = {
    'firing': 'danger',      // 告警中 - 红色
    'resolved': 'success',   // 已恢复 - 绿色
    'acknowledged': 'warning' // 已认领 - 黄色
  }
  return typeMap[normalized] || 'info'
}

// 获取严重级别标签类型
const getSeverityType = (severity) => {
  const typeMap = {
    1: 'danger',    // critical - 红色
    2: 'warning',   // warning - 黄色
    3: 'info'       // info - 蓝色
  }
  return typeMap[severity] || 'info'
}

// 获取状态标签文本
const getStatusLabel = (status) => {
  const normalized = normalizeStatus(status)
  const labelMap = {
    'firing': '告警中',
    'resolved': '已恢复',
    'acknowledged': '已认领'
  }
  return labelMap[normalized] || status
}

// 获取严重级别标签文本
const getSeverityLabel = (severity) => {
  const labelMap = {
    1: '严重',
    2: '警告',
    3: '信息'
  }
  return labelMap[severity] || '未知'
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

// 计算持续时长
const getDuration = (startTime, endTime) => {
  if (!startTime) return '-'
  const end = endTime ? dayjs(endTime) : dayjs()
  const start = dayjs(startTime)
  const diff = end.diff(start, 'second')

  if (diff < 60) return `${diff}秒`
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时`
  return `${Math.floor(diff / 86400)}天`
}

onMounted(() => {
  // 初始化规则列表和业务组列表
  ruleList.value = [
    { id: 1, name: 'CPU使用率' },
    { id: 2, name: '内存使用率' },
    { id: 3, name: '磁盘使用率' }
  ]
  businessGroupList.value = [
    { id: 1, name: '数据库' },
    { id: 2, name: '应用' },
    { id: 3, name: '基础设施' }
  ]
  fetchList()
})
</script>

<template>
  <div class="alert-history-container">
    <!-- 搜索区域 -->
    <div class="search-section">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="告警状态">
          <el-select v-model="searchForm.status" placeholder="全部" size="small" style="width: 120px" clearable>
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="严重级别">
          <el-select v-model="searchForm.severity" placeholder="全部" size="small" style="width: 120px" clearable>
            <el-option v-for="item in severityOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="规则">
          <el-select v-model="searchForm.ruleId" placeholder="全部规则" size="small" style="width: 120px" clearable>
            <el-option
              v-for="item in ruleList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="业务组">
          <el-select v-model="searchForm.businessGroup" placeholder="全部业务组" size="small" style="width: 120px" clearable>
            <el-option
              v-for="item in businessGroupList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.startTime"
            type="datetime"
            placeholder="开始时间"
            size="small"
            style="width: 150px"
          />
          <span style="margin: 0 8px">~</span>
          <el-date-picker
            v-model="searchForm.endTime"
            type="datetime"
            placeholder="结束时间"
            size="small"
            style="width: 150px"
          />
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
        type="primary"
        size="small"
        @click="handleRefresh"
        class="modern-btn primary-btn"
      >
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>

      <el-button
        type="info"
        size="small"
        @click="handleExport"
        class="modern-btn info-btn"
      >
        <el-icon><Download /></el-icon>
        导出
      </el-button>

      <el-button
        type="danger"
        size="small"
        @click="handleBatchDelete"
        :disabled="selectedRows.length === 0"
        v-authority="['monitor:alert:record:delete']"
        class="modern-btn danger-btn"
      >
        <el-icon><DeleteFilled /></el-icon>
        批量删除
      </el-button>

      <el-button
        type="danger"
        size="small"
        @click="handleClearAll"
        v-authority="['monitor:alert:record:clear']"
        class="modern-btn danger-btn"
      >
        <el-icon><Delete /></el-icon>
        清空所有
      </el-button>
    </div>

    <!-- 表格 -->
    <div class="table-section">
      <el-table
        :data="tableData"
        v-loading="loading"
        stripe
        border
        style="width: 100%"
        :header-cell-style="{ background: '#eef1f6', color: '#606266' }"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" align="center" />

        <el-table-column prop="id" label="ID" width="70" align="center" />

        <el-table-column label="告警状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="严重级别" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getSeverityType(row.severity)" size="small">
              {{ getSeverityLabel(row.severity) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="policyName" label="策略名称" min-width="150">
          <template #default="{ row }">
            <span style="font-weight: 500">{{ row.policyName }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="targetName" label="目标名称" min-width="120" show-overflow-tooltip>
          <template #default="{ row }">
            <span>{{ row.targetName || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="message" label="告警信息" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <span style="font-size: 12px; color: #606266">{{ row.message }}</span>
          </template>
        </el-table-column>

        <el-table-column label="当前值 / 阈值" width="120" align="center">
          <template #default="{ row }">
            <span style="font-size: 12px; font-family: monospace">
              {{ row.value }} / {{ row.threshold }}
            </span>
          </template>
        </el-table-column>

        <el-table-column prop="startTime" label="开始时间" width="170" align="center">
          <template #default="{ row }">
            <span style="font-size: 12px">{{ formatTime(row.startTime) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="持续时长" width="100" align="center">
          <template #default="{ row }">
            <span style="font-size: 12px">{{ getDuration(row.startTime, row.endTime) }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="handler" label="处理人" width="100" align="center">
          <template #default="{ row }">
            <span style="font-size: 12px">{{ row.handler || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" fixed="right" width="240" align="center">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-tooltip content="查看详情" placement="top">
                <el-button
                  type="primary"
                  :icon="Reading"
                  size="small"
                  @click="handleViewDetail(row)"
                >详情</el-button>
              </el-tooltip>
              <el-tooltip content="处理" placement="top">
                <el-button
                  type="warning"
                  size="small"
                  @click="openProcessDialog(row)"
                  :disabled="!canProcessRecord(row)"
                >
                  <el-icon><Check /></el-icon>处理
                </el-button>
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button
                  type="danger"
                  :icon="Delete"
                  size="small"
                  @click="handleDelete(row)"
                >删除</el-button>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 分页 -->
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

    <!-- 详情抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      title="告警详情"
      size="40%"
      :close-on-click-modal="false"
    >
      <div v-if="detailData" class="detail-content">
        <!-- 基本信息 -->
        <div class="detail-section">
          <h4 class="section-title">基本信息</h4>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">告警状态：</span>
                <el-tag :type="getStatusType(detailData.status)">
                  {{ getStatusLabel(detailData.status) }}
                </el-tag>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">严重级别：</span>
                <el-tag :type="getSeverityType(detailData.severity)">
                  {{ getSeverityLabel(detailData.severity) }}
                </el-tag>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">策略名称：</span>
                <span>{{ detailData.policyName || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">规则名称：</span>
                <span>{{ detailData.ruleName || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">目标类型：</span>
                <span>{{ detailData.targetType || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">目标名称：</span>
                <span>{{ detailData.targetName || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">指标名称：</span>
                <span>{{ detailData.metricName || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">指纹：</span>
                <span>{{ detailData.fingerprint || '-' }}</span>
              </div>
            </el-col>
          </el-row>
        </div>

        <!-- 时间信息 -->
        <div class="detail-section">
          <h4 class="section-title">时间信息</h4>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">开始时间：</span>
                <span>{{ formatTime(detailData.startTime) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">结束时间：</span>
                <span>{{ formatTime(detailData.endTime) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">持续时长：</span>
                <span>{{ getDuration(detailData.startTime, detailData.endTime) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">通知时间：</span>
                <span>{{ formatTime(detailData.notifyTime) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">通知次数：</span>
                <span>{{ detailData.notifyCount ?? '-' }}</span>
              </div>
            </el-col>
          </el-row>
        </div>

        <!-- 告警信息 -->
        <div class="detail-section">
          <h4 class="section-title">告警信息</h4>
          <el-row :gutter="20">
            <el-col :span="24">
              <div class="detail-item">
                <span class="label">告警信息：</span>
                <span>{{ detailData.message }}</span>
              </div>
            </el-col>
            <el-col :span="24" v-if="detailData.summary">
              <div class="detail-item">
                <span class="label">摘要：</span>
                <span>{{ detailData.summary }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">当前值：</span>
                <span style="font-family: monospace; font-weight: bold">{{ getAlertValue(detailData) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">阈值：</span>
                <span style="font-family: monospace; font-weight: bold">{{ detailData.threshold }}</span>
              </div>
            </el-col>
          </el-row>
        </div>

        <div v-if="Object.keys(parseLabels(detailData.labels)).length > 0" class="detail-section">
          <h4 class="section-title">标签信息</h4>
          <div class="labels-wrap">
            <el-tag
              v-for="(value, key) in parseLabels(detailData.labels)"
              :key="key"
              size="small"
              class="label-tag"
            >
              {{ key }}={{ value }}
            </el-tag>
          </div>
        </div>

        <!-- 处理信息 -->
        <div v-if="detailData.handler" class="detail-section">
          <h4 class="section-title">处理信息</h4>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">处理人：</span>
                <span>{{ detailData.handler }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">处理时间：</span>
                <span>{{ formatTime(getHandleTime(detailData)) }}</span>
              </div>
            </el-col>
            <el-col :span="24">
              <div class="detail-item">
                <span class="label">处理备注：</span>
                <span class="multiline-text">{{ getHandleNote(detailData) }}</span>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-drawer>

    <!-- 统一处理对话框 -->
    <el-dialog
      v-model="processDialogVisible"
      title="处理告警"
      width="520px"
      :close-on-click-modal="false"
    >
      <el-form label-width="90px">
        <el-form-item label="操作类型" required>
          <el-radio-group v-model="processForm.action">
            <el-radio
              v-for="item in actionOptions"
              :key="item.value"
              :label="item.value"
            >
              {{ item.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="shouldShowAssignee" label="选择用户" required>
          <el-select
            v-model="processForm.assigneeId"
            placeholder="请选择用户"
            filterable
            clearable
            style="width: 100%"
          >
            <el-option
              v-for="user in userOptions"
              :key="user.id"
              :label="user.name"
              :value="user.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item :label="isHandleAction ? '处理步骤' : '处理说明'">
          <el-input
            v-model="processForm.handleNote"
            type="textarea"
            :rows="4"
            :placeholder="isHandleAction ? '请输入处理步骤' : '请输入备注说明'"
          />
        </el-form-item>

        <div v-if="currentRecord" class="process-dialog-tip">
          当前记录：#{{ currentRecord.id }} / {{ currentRecord.policyName || '-' }}
        </div>
      </el-form>

      <template #footer>
        <el-button @click="processDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="processLoading" @click="handleProcess">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.alert-history-container {
  padding: 0;
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

.process-dialog-tip {
  margin-top: 8px;
  color: #909399;
  font-size: 12px;
}

.labels-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.label-tag {
  margin: 0;
}

.multiline-text {
  white-space: pre-wrap;
  word-break: break-word;
}

.action-section {
  margin-bottom: 16px;
  display: flex;
  gap: 10px;
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
}

.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.operation-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
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

:deep(.info-btn) {
  background: linear-gradient(45deg, #909399, #A6ADB7);
  color: white;
}

:deep(.danger-btn) {
  background: linear-gradient(45deg, #F56C6C, #F78989);
  color: white;
}

/* 详情抽屉样式 */
.detail-content {
  padding: 0 20px;
}

.detail-section {
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.detail-section:last-child {
  border-bottom: none;
}

.section-title {
  margin: 0 0 16px 0;
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
}

.detail-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
  font-size: 13px;
}

.detail-item .label {
  color: #909399;
  min-width: 80px;
  font-weight: 500;
  margin-right: 8px;
}

.detail-item span:not(.label) {
  color: #606266;
  flex: 1;
  word-break: break-all;
}
</style>
