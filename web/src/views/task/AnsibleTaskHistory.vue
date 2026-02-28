<template>
  <div class="task-history-management">
    <el-card shadow="hover" class="task-history-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button 
              type="primary" 
              link 
              @click="$router.back()" 
              style="margin-right: 15px; font-weight: 600"
            >
              <el-icon style="margin-right: 4px"><ArrowLeft /></el-icon> 返回
            </el-button>
            <span class="title">任务执行历史 - {{ taskName }}</span>
          </div>
          <el-button type="primary" size="small" icon="Refresh" @click="fetchHistory" :loading="loading">刷新列表</el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-section">
        <el-form :inline="true" :model="queryParams" class="search-form">
          <el-form-item label="执行状态">
            <el-select 
              v-model="queryParams.status" 
              placeholder="请选择状态" 
              clearable 
              style="width: 200px"
              size="small"
              @change="handleSearch"
            >
              <el-option label="等待中" :value="1" />
              <el-option label="运行中" :value="2" />
              <el-option label="成功" :value="3" />
              <el-option label="失败" :value="4" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="Search" size="small" @click="handleSearch">搜索</el-button>
            <el-button type="warning" icon="Refresh" size="small" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 列表区域 -->
      <div class="table-section">
        <el-table
          v-loading="loading"
          :data="historyList"
          border
          stripe
          style="width: 100%"
          :header-cell-style="{ background: '#eef1f6', color: '#606266' }"
        >
          <el-table-column prop="ID" label="执行ID" width="100" align="center">
            <template #default="{ row }">
              <span class="id-tag">#{{ row.ID }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="status" label="状态" width="130" align="center">
            <template #default="{ row }">
              <div style="display: flex; align-items: center; justify-content: center; gap: 6px;">
                 <img 
                  src="@/assets/image/zhuangtai.svg" 
                  alt="状态"
                  style="width: 16px; height: 16px;"
                />
                <span :class="getStatusClass(row.status)">
                  {{ getStatusText(row.status) }}
                </span>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="OperatorName" label="执行人" width="150" align="center">
            <template #default="{ row }">
              <div style="display: flex; align-items: center; justify-content: center; gap: 6px;">
                 <img 
                  src="@/assets/image/ren.svg" 
                  alt="用户"
                  style="width: 16px; height: 16px;"
                />
                <!-- 后端 OperatorName 可能为空，可以给个默认值 -->
                <span>{{ row.OperatorName || 'System' }}</span>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="CreatedAt" label="开始时间" width="180" align="center">
            <template #default="{ row }">
              {{ formatTime(row.CreatedAt) }}
            </template>
          </el-table-column>
          
          <el-table-column prop="FinishedAt" label="结束时间" width="180" align="center">
            <template #default="{ row }">
              {{ formatTime(row.FinishedAt) }}
            </template>
          </el-table-column>
          
          <el-table-column prop="TotalDuration" label="耗时" width="120" align="center">
            <template #default="{row}">
              <div style="display: flex; align-items: center; justify-content: center; gap: 6px;">
                <img 
                  src="@/assets/image/定时关闭.svg" 
                  alt="时间"
                  style="width: 16px; height: 16px;"
                />
                <span>{{ row.TotalDuration }}s</span>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="150" align="center" fixed="right">
            <template #default="{ row }">
              <div class="operation-buttons">
                 <el-tooltip content="查看日志" placement="top">
                  <el-button
                    type="primary"
                    :icon="Document"
                    size="small"
                    circle
                    @click="viewLog(row)"
                  />
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页区域 -->
      <div class="pagination-section">
        <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="queryParams.page"
          :page-sizes="[10, 20, 50]"
          :page-size="queryParams.pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
        />
      </div>
    </el-card>

    <!-- 历史详情流程弹窗 -->
    <el-dialog
      v-model="logDialogVisible"
      title="任务执行详情"
      width="80%"
      top="5vh"
      destroy-on-close
      custom-class="ansible-flow-dialog"
    >
      <div class="flow-container" v-loading="detailLoading">
        <AnsibleFlowTemp
          v-if="historySteps.length > 0"
          :steps="historySteps"
          :task-id="taskId"
          :history-mode="true"
          :history-id="currentHistoryId"
        />
        <el-empty v-else description="暂无执行详情数据" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Search, Refresh, Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import AnsibleFlowTemp from './Job/AnsibleFlowTemp.vue'
import { GetAnsibleTaskHistory, GetAnsibleTaskDetail, GetAnsibleHistoryDetail } from '@/api/task'

const route = useRoute()
const router = useRouter()

// 状态变量
const taskId = ref(route.query.id)
const taskName = ref(route.query.name || '未知任务')
const loading = ref(false)
const historyList = ref([])
const total = ref(0)
const logDialogVisible = ref(false)
const currentHistoryId = ref(null)
const historySteps = ref([])

const queryParams = reactive({
  id: taskId.value,
  page: 1,
  pageSize: 10,
  status: '',
  taskId: taskId.value
})

// Mock数据，实际应调用API
const fetchHistory = async () => {
  loading.value = true
  try {
    const res = await GetAnsibleTaskHistory(queryParams)
    // 尝试做兼容处理
    const responseData = res.data || res
    
    // 取出内部的 data 和 total
    // 路径: responseData -> data -> data (list)
    if (responseData && responseData.data) {
        historyList.value = responseData.data.data || []
        total.value = responseData.data.total || 0
    } else {
        historyList.value = []
        total.value = 0
    }
  } catch (error) {
    console.error('获取历史记录失败', error)
    ElMessage.error('获取历史记录失败')
    historyList.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  queryParams.page = 1
  fetchHistory()
}

const resetQuery = () => {
  queryParams.status = ''
  handleSearch()
}

const handleSizeChange = (val) => {
  queryParams.pageSize = val
  fetchHistory()
}

const handleCurrentChange = (val) => {
  queryParams.page = val
  fetchHistory()
}

const detailLoading = ref(false)

const viewLog = async (row) => {
  logDialogVisible.value = true
  detailLoading.value = true
  currentHistoryId.value = row.ID
  historySteps.value = []

  try {
    // 尝试获取详细的历史执行信息 (包含子任务状态)
    const historyDetailRes = await GetAnsibleHistoryDetail({
      id: taskId.value,
      historyId: row.ID
    })

    let works = []
    
    // 检查是否有 WorkHistories 数据
    if (historyDetailRes.data && historyDetailRes.data.data && historyDetailRes.data.data.WorkHistories) {
       works = historyDetailRes.data.data.WorkHistories
       
       console.log('获取到详细历史子任务数据:', works)
       
       historySteps.value = works.map(work => ({
         // WorkHistories 中的字段映射
         task_id: taskId.value,
         work_id: work.WorkID,
         entry_file_name: work.HostName, // 后端返回的是 HostName 对应文件名
         status: work.Status, 
         duration: work.Duration,
         // 保留原始数据以备不时之需
         original_work: work 
       }))
    } 
    // 如果没有获取到详细历史，回退使用列表中的 Works 数据
    else if (row.Works && row.Works.length > 0) {
      console.log('使用列表中的 Works 数据')
      historySteps.value = row.Works.map(work => ({
        ...work,
        task_id: taskId.value,
        work_id: work.ID || work.work_id || work.WorkID, 
        entry_file_name: work.EntryFileName || work.entry_file_name,
        // 直接使用Work的Status，如果后端返回的是数值1~4
        status: work.Status !== undefined ? work.Status : (work.status !== undefined ? work.status : row.status), 
        duration: work.Duration || work.duration
      }))
    } else {
        // 如果都没有，尝试获取任务详情作为结构补充（这种情况下状态不可用，默认成功或失败）
        try {
           const res = await GetAnsibleTaskDetail(taskId.value)
           if (res.data && res.data.data && res.data.data.task_info && res.data.data.task_info.Works) {
             const staticWorks = res.data.data.task_info.Works
             console.log('使用任务定义的 Works 结构')
             historySteps.value = staticWorks.map(work => ({
               ...work,
               task_id: taskId.value,
               work_id: work.workid, 
               entry_file_name: work.EntryFileName,
               status: row.status, // 使用整体状态作为默认状态
               duration: 0
             }))
           }
        } catch (e) {
            console.error('获取任务详情失败', e)
        }
    }
    
    console.log('History Steps computed:', historySteps.value)
  } catch (error) {
    console.error('获取任务详情失败:', error)
    // 出错时也尝试使用 fallback
    if (row.Works && row.Works.length > 0) {
        historySteps.value = row.Works.map(work => ({
        ...work,
        task_id: taskId.value,
        work_id: work.ID || work.work_id || work.WorkID, 
        entry_file_name: work.EntryFileName || work.entry_file_name,
        status: work.Status !== undefined ? work.Status : (work.status !== undefined ? work.status : row.status), 
        duration: work.Duration || work.duration
      }))
    } else {
        ElMessage.error('获取任务详情失败')
    }
  } finally {
    detailLoading.value = false
  }
}

// 辅助函数
const getStatusClass = (status) => {
  switch(status) {
    case 1: return 'status-waiting'
    case 2: return 'status-running'
    case 3: return 'status-success'
    case 4: return 'status-error'
    default: return ''
  }
}

const getStatusText = (status) => {
  const map = { 1: '等待中', 2: '运行中', 3: '成功', 4: '失败' }
  return map[status] || '未知'
}

const formatTime = (timeStr) => {
  if (!timeStr) return '-'
  try {
    const date = new Date(timeStr)
    // 检查是否是一年后的时间或无效时间，如果是1年后的时间(Golang zero time通常很小，但这里可能是null)
    // 0001-01-01T00:00:00Z 是 Go 的零值时间
    if (timeStr.startsWith('0001-01-01')) return '-'
    
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    }).replace(/\//g, '-')
  } catch (e) {
    return timeStr
  }
}

onMounted(() => {
  if (!taskId.value) {
     // ElMessage.warning('未指定任务ID')
  }
  fetchHistory()
})
</script>

<style scoped>
.task-history-management {
  padding: 20px;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.task-history-card {
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

.header-left {
  display: flex;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-left: 10px;
}

.search-section {
  margin-bottom: 16px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

.table-section {
  margin-top: 0;
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

/* 统一表格样式 */
:deep(.el-table) {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

:deep(.el-table__header) {
  background: rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-table__header th) {
  background: rgba(102, 126, 234, 0.1) !important;
  color: #2c3e50 !important;
  font-weight: 600;
  border: none;
}

:deep(.el-table__body tr:hover > td) {
  background-color: rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-table td) {
  border: none;
}

:deep(.el-table::before) {
  display: none;
}

:deep(.el-table--border::after) {
  display: none;
}

/* 按钮样式复用 */
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

:deep(.el-button--primary) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

:deep(.el-button--primary:hover) {
  background: linear-gradient(135deg, #5a6fd8, #6a4190);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

:deep(.el-button--primary:active) {
  transform: translateY(0);
}

:deep(.el-button--warning) {
  background: linear-gradient(135deg, #f39c12, #e67e22);
  border: none;
  border-radius: 8px;
  font-weight: 500;
  color: white;
  transition: all 0.3s ease;
}

:deep(.el-button--warning:hover) {
  background: linear-gradient(135deg, #e67e22, #d35400);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(243, 156, 18, 0.4);
}

:deep(.el-button--warning:active) {
  transform: translateY(0);
}

/* 输入框样式 */
:deep(.el-form-item__label) {
  color: #2c3e50;
  font-weight: 500;
}

:deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(102, 126, 234, 0.3);
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.el-input__wrapper:hover) {
  border-color: rgba(102, 126, 234, 0.5);
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.2);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

:deep(.el-select .el-input.is-focus .el-input__wrapper) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

:deep(.el-pagination) {
  --el-pagination-bg-color: transparent;
}

:deep(.el-pager li) {
  background: rgba(255, 255, 255, 0.8);
  border-radius: 6px;
  margin: 0 2px;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

:deep(.el-pager li:hover) {
  background: rgba(102, 126, 234, 0.1);
}

:deep(.el-pager li.is-active) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
}

.id-tag {
  color: #909399; 
  font-family: monospace; 
  background: #f4f4f5; 
  padding: 2px 6px; 
  border-radius: 4px;
}

.result-text {
  display: inline-block;
  max-width: 250px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #606266;
}

.status-waiting {
  color: #909399;
  font-weight: 500;
}

.status-running {
  color: #E6A23C;
  font-weight: 500;
}

.status-success {
  color: #67C23A;
  font-weight: 500;
}


:deep(.ansible-flow-dialog) {
  background: #e6e6e6;
}

:deep(.ansible-flow-dialog .el-dialog__header) {
  background: #e6e6e6;
  border-bottom: 1px solid #d9d9d9;
}

:deep(.ansible-flow-dialog .el-dialog__body) {
  padding: 20px;
  background: #e6e6e6;
}

.flow-container {
  min-height: 400px;
  padding: 20px;
  background: #e8e8e8;
  border-radius: 6px;
}
</style>
