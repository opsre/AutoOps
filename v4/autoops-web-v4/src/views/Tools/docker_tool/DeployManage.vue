<template>
  <!-- 标准页面模式 -->
  <div v-if="!embedded" class="deploy-manage">
    <el-card shadow="hover" class="manage-card">
      <template #header>
        <div class="card-header">
          <span class="title">部署管理</span>
        </div>
      </template>

      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :model="queryParams" :inline="true" class="search-form">
          <el-form-item label="服务名称">
            <el-input
              v-model="queryParams.serviceName"
              placeholder="请输入服务名称"
              clearable
              size="small"
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="queryParams.status" clearable size="small" placeholder="请选择状态">
              <el-option label="全部" :value="-1" />
              <el-option label="部署中" :value="0" />
              <el-option label="运行中" :value="1" />
              <el-option label="已停止" :value="2" />
              <el-option label="部署失败" :value="3" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="Search" size="small" @click="handleQuery">查询</el-button>
            <el-button icon="Refresh" size="small" @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格 -->
      <div class="table-section">
        <el-table
          :data="tableData"
          v-loading="loading"
          stripe
          class="deploy-table"
        >
          <el-table-column prop="serviceName" label="服务名称" width="150">
            <template #default="{ row }">
              <div style="display: flex; align-items: center; gap: 8px;">
                <el-icon :color="getServiceColor(row.serviceName)">
                  <component :is="getServiceIcon(row.serviceName)" />
                </el-icon>
                <span>{{ row.serviceName }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="version" label="版本" width="150" />
          <el-table-column label="主机" width="200">
            <template #default="{ row }">
              <div>{{ row.hostName }}</div>
              <div style="color: #999; font-size: 12px">{{ row.hostIp }}</div>
            </template>
          </el-table-column>
          <el-table-column prop="installDir" label="安装目录" width="200" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">
                <el-icon><component :is="getStatusIcon(row.status)" /></el-icon>
                {{ row.statusText }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="部署时间" width="180">
            <template #default="{ row }">
              {{ formatTime(row.createTime) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" fixed="right" width="200">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                icon="View"
                @click="viewLog(row)"
              >
                查看日志
              </el-button>
              <el-popconfirm
                title="确定要卸载此服务吗？"
                @confirm="handleUninstall(row.id)"
              >
                <template #reference>
                  <el-button type="danger" size="small" icon="Delete">
                    卸载
                  </el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination-section">
        <el-pagination
          v-model:current-page="queryParams.pageNum"
          v-model:page-size="queryParams.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleQuery"
          @current-change="handleQuery"
        />
      </div>
    </el-card>

    <!-- 日志查看对话框 -->
    <DeployProgress
      v-model="logVisible"
      :deploy-id="currentDeployId"
    />
  </div>

  <!-- 嵌入模式（用于服务市场内部切换） -->
  <template v-else>
    <div class="deploy-embedded">
      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :model="queryParams" :inline="true" class="search-form">
          <el-form-item label="服务名称">
            <el-input
              v-model="queryParams.serviceName"
              placeholder="请输入服务名称"
              clearable
              size="small"
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="queryParams.status" clearable size="small" placeholder="请选择状态">
              <el-option label="全部" :value="-1" />
              <el-option label="部署中" :value="0" />
              <el-option label="运行中" :value="1" />
              <el-option label="已停止" :value="2" />
              <el-option label="部署失败" :value="3" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="Search" size="small" @click="handleQuery">查询</el-button>
            <el-button icon="Refresh" size="small" @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格 -->
      <div class="table-section">
        <el-table
          :data="tableData"
          v-loading="loading"
          stripe
          class="deploy-table"
        >
          <el-table-column prop="serviceName" label="服务名称" width="150">
            <template #default="{ row }">
              <div style="display: flex; align-items: center; gap: 8px;">
                <el-icon :color="getServiceColor(row.serviceName)">
                  <component :is="getServiceIcon(row.serviceName)" />
                </el-icon>
                <span>{{ row.serviceName }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="version" label="版本" width="150" />
          <el-table-column label="主机" width="200">
            <template #default="{ row }">
              <div>{{ row.hostName }}</div>
              <div style="color: #999; font-size: 12px">{{ row.hostIp }}</div>
            </template>
          </el-table-column>
          <el-table-column prop="installDir" label="安装目录" width="200" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">
                <el-icon><component :is="getStatusIcon(row.status)" /></el-icon>
                {{ row.statusText }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="部署时间" width="180">
            <template #default="{ row }">
              {{ formatTime(row.createTime) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" fixed="right" width="200">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                icon="View"
                @click="viewLog(row)"
              >
                查看日志
              </el-button>
              <el-popconfirm
                title="确定要卸载此服务吗？"
                @confirm="handleUninstall(row.id)"
              >
                <template #reference>
                  <el-button type="danger" size="small" icon="Delete">
                    卸载
                  </el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination-section">
        <el-pagination
          v-model:current-page="queryParams.pageNum"
          v-model:page-size="queryParams.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleQuery"
          @current-change="handleQuery"
        />
      </div>

      <!-- 日志查看对话框 -->
      <DeployProgress
        v-model="logVisible"
        :deploy-id="currentDeployId"
      />
    </div>
  </template>
</template>

<script setup>
const props = defineProps({
  embedded: { type: Boolean, default: false }
})
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { getDeployList, deleteDeploy } from '@/api/tool'
import { ElMessage } from 'element-plus'
import {
  DataAnalysis,
  Operation,
  Monitor,
  CircleCheck,
  CircleClose,
  Loading,
  Warning
} from '@element-plus/icons-vue'
import DeployProgress from './DeployProgress.vue'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const logVisible = ref(false)
const currentDeployId = ref(null)
let pollingTimer = null

const queryParams = reactive({
  serviceName: '',
  status: -1,
  pageNum: 1,
  pageSize: 10
})

// 服务图标映射
const serviceIconMap = {
  'mysql': DataAnalysis,
  'redis': DataAnalysis,
  'postgresql': DataAnalysis,
  'jenkins': Operation,
  'gitlab': Operation,
  'grafana': Monitor,
  'prometheus': Monitor
}

const getServiceIcon = (serviceName) => {
  const name = serviceName?.toLowerCase() || ''
  return serviceIconMap[name] || Operation
}

const getServiceColor = (serviceName) => {
  const name = serviceName?.toLowerCase() || ''
  if (name.includes('mysql') || name.includes('redis') || name.includes('postgresql')) {
    return '#409eff'
  } else if (name.includes('jenkins') || name.includes('gitlab')) {
    return '#67c23a'
  } else if (name.includes('grafana') || name.includes('prometheus')) {
    return '#e6a23c'
  }
  return '#909399'
}

// 状态类型映射
const getStatusType = (status) => {
  const types = {
    0: 'primary',
    1: 'success',
    2: 'info',
    3: 'danger'
  }
  return types[status] || 'info'
}

const getStatusIcon = (status) => {
  const icons = {
    0: 'Loading',
    1: 'CircleCheck',
    2: 'Warning',
    3: 'CircleClose'
  }
  return icons[status] || 'Warning'
}

// 查询列表
const handleQuery = async () => {
  loading.value = true
  try {
    const params = {
      pageNum: queryParams.pageNum,
      pageSize: queryParams.pageSize
    }
    if (queryParams.serviceName) {
      params.serviceName = queryParams.serviceName
    }
    if (queryParams.status !== -1) {
      params.status = queryParams.status
    }

    const res = await getDeployList(params)
    if (res.data?.code === 200) {
      tableData.value = res.data.data?.list || []
      total.value = res.data.data?.total || 0
    } else {
      throw new Error(res.data?.message || '查询失败')
    }
  } catch (error) {
    console.error('查询失败:', error)
    ElMessage.error(`查询失败: ${error.message}`)
  } finally {
    loading.value = false
  }
}

// 重置查询
const handleReset = () => {
  queryParams.serviceName = ''
  queryParams.status = -1
  queryParams.pageNum = 1
  handleQuery()
}

// 查看日志
const viewLog = (row) => {
  currentDeployId.value = row.id
  logVisible.value = true
}

// 卸载服务
const handleUninstall = async (id) => {
  try {
    const res = await deleteDeploy(id)
    if (res.data?.code === 200) {
      ElMessage.success('卸载成功')
      handleQuery()
    } else {
      ElMessage.error(res.data?.message || '卸载失败')
    }
  } catch (error) {
    console.error('卸载失败:', error)
    ElMessage.error('卸载失败')
  }
}

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr || timeStr === '0001-01-01 00:00:00') return '无'
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

// 开始轮询
const startPolling = () => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
  }
  pollingTimer = setInterval(() => {
    const hasDeploying = tableData.value.some(item => item.status === 0)
    if (hasDeploying) {
      handleQuery()
    }
  }, 3000)
}

// 停止轮询
const stopPolling = () => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
}

onMounted(() => {
  handleQuery()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.deploy-manage {
  padding: 20px;
  height: 80vh;
  background: #eef1f6;
}

.manage-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  height: calc(80vh - 40px);
  display: flex;
  flex-direction: column;
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
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.search-section {
  margin-bottom: 5px;
  padding: 20px;
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
  margin-bottom: 20px;
  flex: 1;
  overflow-y: auto;
}

.deploy-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.deploy-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.deploy-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.deploy-table :deep(.el-table__header th .cell) {
  color: #2c3e50 !important;
  font-weight: 700 !important;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.deploy-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.deploy-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.pagination-section {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.el-input :deep(.el-input__wrapper),
.el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__wrapper):hover,
.el-select :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.el-input :deep(.el-input__wrapper.is-focus),
.el-select :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.el-pagination {
  background: transparent;
}

.el-pagination :deep(.btn-next),
.el-pagination :deep(.btn-prev),
.el-pagination :deep(.el-pager li) {
  border-radius: 6px;
  margin: 0 2px;
  transition: all 0.3s ease;
}

.el-pagination :deep(.btn-next):hover,
.el-pagination :deep(.btn-prev):hover,
.el-pagination :deep(.el-pager li):hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}
</style>
