<script setup>
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  View,
  Edit,
  Delete,
  Search,
  Refresh,
  Plus,
  DocumentCopy,
  WarningFilled,
  CircleCheck,
  Clock,
  Select,
  Link
} from '@element-plus/icons-vue'
import { incidentApi } from '@/api/monitor'
import cmdbAPI from '@/api/cmdb'

// 统计数据
const statistics = ref({
  total: 0,
  pending: 0,
  processing: 0,
  completed: 0
})

// 业务线分组数据（二级分组）
const businessLineOptions = ref([])

// 搜索表单
const searchForm = reactive({
  businessLineId: null,  // 改为 businessLineId
  alertLevel: '',
  status: -1,
  department: '',
  startTime: '',
  endTime: '',
  pageNum: 1,
  pageSize: 10
})

// 列表数据
const tableData = ref([])
const total = ref(0)
const loading = ref(false)

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('添加故障记录')
const formLoading = ref(false)

// 表单数据
const incidentForm = reactive({
  id: null,
  alertTime: '',
  businessLineId: null,  // 改为 businessLineId，存储分组ID
  frequency: '偶发',
  alertDesc: '',
  alertLevel: 'P4',
  incidentCause: '',
  department: '运维部',
  solution: '',
  detailUrl: '',
  status: 1,
  remark: ''
})

// 表单验证规则
const formRules = {
  alertTime: [
    { required: true, message: '请选择告警时间', trigger: 'change' }
  ],
  businessLineId: [
    { required: true, message: '请选择业务线', trigger: 'change' }
  ],
  frequency: [
    { required: true, message: '请选择频率', trigger: 'change' }
  ],
  alertDesc: [
    { required: true, message: '请输入告警描述', trigger: 'blur' }
  ]
}

const formRef = ref(null)

// 选项配置
const alertLevelOptions = [
  { label: '全部', value: '' },
  { label: 'P1', value: 'P1' },
  { label: 'P2', value: 'P2' },
  { label: 'P3', value: 'P3' },
  { label: 'P4', value: 'P4' }
]

const statusOptions = [
  { label: '全部', value: -1 },
  { label: '未处理', value: 1 },
  { label: '处理中', value: 2 },
  { label: '已完成', value: 3 }
]

const frequencyOptions = [
  { label: '偶发', value: '偶发' },
  { label: '频繁', value: '频繁' }
]

const departmentOptions = [
  { label: '研发部', value: '研发部' },
  { label: '运维部', value: '运维部' },
  { label: '第三方', value: '第三方' }
]

// 获取业务线（二级分组）列表
const fetchBusinessLines = async () => {
  try {
    const { data: res } = await cmdbAPI.getAllCmdbGroups()
    console.log('获取分组数据响应:', res)

    if (res.code === 200) {
      // 提取所有二级分组（parentId !== 0 的分组）
      const secondLevelGroups = []
      const extractSecondLevel = (groups) => {
        if (!Array.isArray(groups)) {
          console.warn('groups 不是数组:', groups)
          return
        }

        groups.forEach(group => {
          if (group.parentId !== 0) {
            secondLevelGroups.push({
              label: group.name,
              value: group.id,  // value 改为 ID
              id: group.id
            })
          }
          if (group.children && group.children.length > 0) {
            extractSecondLevel(group.children)
          }
        })
      }

      extractSecondLevel(res.data || [])
      console.log('提取的二级分组:', secondLevelGroups)
      businessLineOptions.value = secondLevelGroups

      if (secondLevelGroups.length === 0) {
        console.warn('未找到任何二级分组')
        ElMessage.warning('未找到可用的业务线分组')
      }
    } else {
      console.error('获取分组数据失败:', res)
      ElMessage.error(res.message || '获取业务线列表失败')
    }
  } catch (error) {
    console.error('获取业务线列表失败:', error)
    ElMessage.error('获取业务线列表失败')
  }
}

// 获取统计信息
const fetchStatistics = async () => {
  try {
    const { data: res } = await incidentApi.getStatistics()
    if (res.code === 200) {
      statistics.value = res.data
    }
  } catch (error) {
    console.error('获取统计信息失败:', error)
  }
}

// 获取故障列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = { ...searchForm }
    // 处理 -1 值和空值，不传给后端
    if (params.status === -1) delete params.status
    if (!params.alertLevel) delete params.alertLevel
    if (!params.businessLineId) delete params.businessLineId

    const { data: res } = await incidentApi.getList(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    } else {
      ElMessage.error(res.msg || '获取故障列表失败')
    }
  } catch (error) {
    ElMessage.error('获取故障列表失败')
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
  searchForm.businessLineId = null
  searchForm.alertLevel = ''
  searchForm.status = -1
  searchForm.department = ''
  searchForm.startTime = ''
  searchForm.endTime = ''
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

// 添加故障记录
const handleAdd = () => {
  dialogTitle.value = '添加故障记录'
  incidentForm.id = null
  incidentForm.alertTime = ''
  incidentForm.businessLineId = null
  incidentForm.frequency = '偶发'
  incidentForm.alertDesc = ''
  incidentForm.alertLevel = 'P4'
  incidentForm.incidentCause = ''
  incidentForm.department = '运维部'
  incidentForm.solution = ''
  incidentForm.detailUrl = ''
  incidentForm.status = 1
  incidentForm.remark = ''
  dialogVisible.value = true
}

// 编辑故障记录
const handleEdit = (row) => {
  dialogTitle.value = '编辑故障记录'
  incidentForm.id = row.id
  incidentForm.alertTime = row.alertTime
  incidentForm.businessLineId = row.businessLineId || row.businessLine  // 兼容两种字段
  incidentForm.frequency = row.frequency
  incidentForm.alertDesc = row.alertDesc
  incidentForm.alertLevel = row.alertLevel
  incidentForm.incidentCause = row.incidentCause
  incidentForm.department = row.department
  incidentForm.solution = row.solution
  incidentForm.detailUrl = row.detailUrl
  incidentForm.status = row.status
  incidentForm.remark = row.remark
  dialogVisible.value = true
}

// 删除故障记录
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该故障记录吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const { data: res } = await incidentApi.delete(row.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      fetchList()
      fetchStatistics()
    } else {
      ElMessage.error(res.msg || '删除失败')
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
      const data = { ...incidentForm }
      let result
      if (data.id) {
        result = await incidentApi.update(data)
      } else {
        delete data.id
        result = await incidentApi.create(data)
      }

      const { data: res } = result
      if (res.code === 200) {
        ElMessage.success(data.id ? '更新成功' : '添加成功')
        dialogVisible.value = false
        fetchList()
        fetchStatistics()
      } else {
        ElMessage.error(res.msg || '操作失败')
      }
    } catch (error) {
      ElMessage.error('操作失败')
      console.error(error)
    } finally {
      formLoading.value = false
    }
  })
}

// 快速更新状态
const handleQuickUpdateStatus = async (row, status) => {
  try {
    const statusText = ['', '未处理', '处理中', '已完成'][status]
    await ElMessageBox.confirm(`确定将状态更新为"${statusText}"吗?`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })

    const { data: res } = await incidentApi.updateStatus(row.id, status)
    if (res.code === 200) {
      ElMessage.success('状态更新成功')
      fetchList()
      fetchStatistics()
    } else {
      ElMessage.error(res.msg || '状态更新失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('状态更新失败')
      console.error(error)
    }
  }
}

// 状态类型
const getStatusType = (status) => {
  const typeMap = {
    1: 'danger',
    2: 'warning',
    3: 'success'
  }
  return typeMap[status] || 'info'
}

// 状态文本
const getStatusText = (status) => {
  const textMap = {
    1: '未处理',
    2: '处理中',
    3: '已完成'
  }
  return textMap[status] || '-'
}

// 告警级别类型
const getAlertLevelType = (level) => {
  const typeMap = {
    'P1': 'danger',
    'P2': 'warning',
    'P3': 'info',
    'P4': 'success'
  }
  return typeMap[level] || 'info'
}

// 格式化日期时间
const formatDateTime = (dateTime) => {
  if (!dateTime) return '-'
  try {
    return new Date(dateTime).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (e) {
    return dateTime
  }
}

// 打开详情链接
const handleOpenUrl = (url) => {
  if (!url) return
  window.open(url, '_blank')
}

// 初始化
onMounted(() => {
  fetchBusinessLines()
  fetchStatistics()
  fetchList()
})
</script>

<template>
  <div class="incident-container">

    <!-- 统计卡片 -->
    <div class="statistics-cards">
      <div class="stat-card stat-total">
        <div class="stat-icon">
          <el-icon :size="32"><DocumentCopy /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">故障总数</div>
          <div class="stat-value">{{ statistics.total }}</div>
        </div>
      </div>

      <div class="stat-card stat-pending">
        <div class="stat-icon">
          <el-icon :size="32"><WarningFilled /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">未处理</div>
          <div class="stat-value">{{ statistics.pending }}</div>
        </div>
      </div>

      <div class="stat-card stat-processing">
        <div class="stat-icon">
          <el-icon :size="32"><Clock /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">处理中</div>
          <div class="stat-value">{{ statistics.processing }}</div>
        </div>
      </div>

      <div class="stat-card stat-completed">
        <div class="stat-icon">
          <el-icon :size="32"><CircleCheck /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">已完成</div>
          <div class="stat-value">{{ statistics.completed }}</div>
        </div>
      </div>
    </div>

    <!-- 主卡片 -->
    <el-card shadow="hover" class="incident-card">
      <template #header>
        <div class="card-header">
          <span class="title">故障管理</span>
        </div>
      </template>

      <!-- 搜索区域 -->
      <div class="search-section">
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="业务线">
            <el-select v-model="searchForm.businessLineId" size="small" placeholder="请选择" clearable style="width: 140px">
              <el-option v-for="item in businessLineOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>

          <el-form-item label="故障级别">
            <el-select v-model="searchForm.alertLevel" size="small" placeholder="请选择" clearable style="width: 120px">
              <el-option v-for="item in alertLevelOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>

          <el-form-item label="处理状态">
            <el-select v-model="searchForm.status" size="small" placeholder="请选择" style="width: 120px">
              <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>

          <el-form-item label="所属部门">
            <el-select v-model="searchForm.department" size="small" placeholder="请选择" clearable style="width: 120px">
              <el-option v-for="item in departmentOptions" :key="item.value" :label="item.label" :value="item.value" />
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
          v-authority="['monitor:incident:add']"
          class="modern-btn success-btn"
        >
          <el-icon><Plus /></el-icon>
          添加故障记录
        </el-button>
      </div>

      <!-- 故障列表 -->
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

          <el-table-column prop="alertTime" label="告警时间" width="150">
            <template #default="{ row }">
              <span style="font-size: 12px; color: #606266">{{ formatDateTime(row.alertTime) }}</span>
            </template>
          </el-table-column>


          <el-table-column prop="businessLine" label="业务线" width="120" align="center">
            <template #default="{ row }">
              <el-tag size="small" type="primary" style="font-size: 11px">
                {{ row.businessLineName || '-' }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="故障级别" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="getAlertLevelType(row.alertLevel)" size="small" style="font-size: 11px">
                {{ row.alertLevel }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="frequency" label="频率" width="80" align="center">
            <template #default="{ row }">
              <span style="font-size: 12px">{{ row.frequency }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="alertDesc" label="告警描述" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
              <span style="font-size: 12px; color: #606266">{{ row.alertDesc }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="detailUrl" label="故障详情" width="100" align="center">
            <template #default="{ row }">
              <el-button
                v-if="row.detailUrl"
                v-authority="['monitor:incident:url']"
                type="primary"
                :icon="Link"
                size="small"
                @click="handleOpenUrl(row.detailUrl)"
              >详情链接</el-button>
              <span v-else style="font-size: 12px; color: #909399">-</span>
            </template>
          </el-table-column>

          <el-table-column prop="department" label="所属部门" width="100" align="center">
            <template #default="{ row }">
              <span style="font-size: 12px">{{ row.department || '-' }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="handler" label="处理人" width="100" align="center">
            <template #default="{ row }">
              <span style="font-size: 12px">{{ row.handler || '-' }}</span>
            </template>
          </el-table-column>

          <el-table-column label="处理状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small" style="font-size: 11px">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="300" fixed="right" align="center">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-tooltip content="编辑" placement="top">
                  <el-button
                    type="warning"
                    :icon="Edit"
                    size="small"
                    v-authority="['monitor:incident:edit']"
                    @click="handleEdit(row)"
                  >编辑</el-button>
                </el-tooltip>
                <el-dropdown trigger="click" @command="(cmd) => handleQuickUpdateStatus(row, cmd)">
                  <el-button type="success" :icon="Select" v-authority="['monitor:incident:status']" size="small">更新状态</el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item :command="1">未处理</el-dropdown-item>
                      <el-dropdown-item :command="2">处理中</el-dropdown-item>
                      <el-dropdown-item :command="3">已完成</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
                <el-tooltip content="删除" placement="top">
                  <el-button
                    type="danger"
                    :icon="Delete"
                    size="small"
                    v-authority="['monitor:incident:delete']"
                    @click="handleDelete(row)"
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
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="700px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="incidentForm"
        :rules="formRules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="告警时间" prop="alertTime">
              <el-date-picker
                v-model="incidentForm.alertTime"
                type="datetime"
                placeholder="选择日期时间"
                style="width: 100%"
                value-format="YYYY-MM-DDTHH:mm:ssZ"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="业务线" prop="businessLineId">
              <el-select v-model="incidentForm.businessLineId" placeholder="请选择业务线（二级分组）" style="width: 100%" filterable>
                <el-option v-for="item in businessLineOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="故障级别" prop="alertLevel">
              <el-select v-model="incidentForm.alertLevel" placeholder="请选择故障级别" style="width: 100%">
                <el-option label="P1" value="P1" />
                <el-option label="P2" value="P2" />
                <el-option label="P3" value="P3" />
                <el-option label="P4" value="P4" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="频率" prop="frequency">
              <el-select v-model="incidentForm.frequency" placeholder="请选择频率" style="width: 100%">
                <el-option v-for="item in frequencyOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="所属部门">
              <el-select v-model="incidentForm.department" placeholder="请选择部门" style="width: 100%">
                <el-option v-for="item in departmentOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="告警描述" prop="alertDesc">
          <el-input
            v-model="incidentForm.alertDesc"
            type="textarea"
            :rows="3"
            placeholder="请输入告警描述"
          />
        </el-form-item>

        <el-form-item label="根因分析">
          <el-input
            v-model="incidentForm.incidentCause"
            type="textarea"
            :rows="2"
            placeholder="请分析故障根因"
          />
        </el-form-item>

        <el-form-item label="详情链接">
          <el-input v-model="incidentForm.detailUrl" placeholder="例如: https://monitor.example.com/alert/12345" />
        </el-form-item>

        <el-form-item label="处理状态">
          <el-radio-group v-model="incidentForm.status">
            <el-radio :label="1">未处理</el-radio>
            <el-radio :label="2">处理中</el-radio>
            <el-radio :label="3">已完成</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="备注">
          <el-input
            v-model="incidentForm.remark"
            type="textarea"
            :rows="2"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="formLoading">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.incident-container {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

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

.stat-card.stat-pending .stat-icon {
  background: linear-gradient(135deg, #f56c6c 0%, #f78989 100%);
  color: white;
}

.stat-card.stat-pending .stat-value {
  color: #f56c6c;
}

.stat-card.stat-processing .stat-icon {
  background: linear-gradient(135deg, #e6a23c 0%, #eebe77 100%);
  color: white;
}

.stat-card.stat-processing .stat-value {
  color: #e6a23c;
}

.stat-card.stat-completed .stat-icon {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  color: white;
}

.stat-card.stat-completed .stat-value {
  color: #67c23a;
}

/* 主卡片 */
.incident-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
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

/* 响应式调整 */
@media (max-width: 1400px) {
  .statistics-cards {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }
}

@media (max-width: 768px) {
  .incident-container {
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
</style>
