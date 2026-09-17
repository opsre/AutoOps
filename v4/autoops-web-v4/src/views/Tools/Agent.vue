<template>
  <div class="agent-management">
    <el-card shadow="hover" class="agent-card">
      <template #header>
        <div class="card-header">
          <span class="title">Agent管理</span>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :model="queryParams" :inline="true" class="search-form">
      <el-form-item label="主机名称">
        <el-input
          placeholder="请输入主机名称"
          size="small"
          clearable
          v-model="queryParams.hostName"
          @keyup.enter="searchByHostName(queryParams.hostName)"
        />
      </el-form-item>
      <el-form-item label="Agent状态" style="width: 200px">
        <el-select
          size="small"
          placeholder="请选择Agent状态"
          v-model="queryParams.status"
          @change="searchByStatus(queryParams.status)"
          style="width: 100%"
        >
          <el-option label="部署中" :value="1" />
          <el-option label="部署失败" :value="2" />
          <el-option label="运行中" :value="3" />
          <el-option label="启动异常" :value="4" />
        </el-select>
      </el-form-item>
      <el-form-item label="版本">
        <el-input
          placeholder="请输入版本号"
          size="small"
          clearable
          v-model="queryParams.version"
          @keyup.enter="searchByVersion(queryParams.version)"
        />
      </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="Search" size="small" @click="fetchAgents">搜索</el-button>
            <el-button type="warning" icon="Refresh" size="small" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 操作按钮 -->
      <div class="action-section">
        <el-button
          type="success"
          plain
          icon="Plus"
          size="small"
          v-authority="['ops:agent:create']"
          @click="handleDeployToHosts"
        >部署Agent</el-button>
        <el-button
          type="danger"
          plain
          icon="Delete"
          v-authority="['ops:agent:deleteall']"
          size="small"
          @click="batchUninstall"
          :disabled="selectedAgents.length === 0"
        >批量卸载({{ selectedAgents.length }})</el-button>
      </div>

      <!-- 列表区域 -->
      <div class="table-section">
        <el-table
          v-loading="loading"
          :data="agents"
          stripe
          style="width: 100%"
          class="agent-table"
          @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="hostName" label="主机名称" width="150">
        <template #default="{row}">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img src="@/assets/image/主机在线.svg" style="width: 16px; height: 16px"/>
            <span>{{ row.hostName }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="sshIp" label="IP地址" width="140">
        <template #default="{row}">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img src="@/assets/image/ip.svg" style="width: 16px; height: 16px"/>
            <span>{{ row.sshIp }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="version" label="版本" width="120">
        <template #default="{row}">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img src="@/assets/image/版本.svg" style="width: 16px; height: 16px"/>
            <span>{{ row.version }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="120">
        <template #default="{row}">
          <span :class="getStatusClass(row.status)">
            <el-icon><component :is="getStatusIcon(row.status)" /></el-icon>
            {{ getStatusText(row.status) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="port" label="监听端口" width="100">
        <template #default="{row}">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img src="@/assets/image/端口.svg" style="width: 16px; height: 16px"/>
            <span>{{ row.port }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="安装进度" width="150">
        <template #default="{row}">
          <div v-if="row.status === 1">
            <el-progress 
              :percentage="row.installProgress || 0" 
              status="success"
              :stroke-width="8"
            />
            <div style="margin-top: 4px; font-size: 12px; color: #67C23A;">
              {{ row.installProgressText || '部署中...' }}
            </div>
          </div>
          <div v-else-if="row.status === 3 && row.installProgress === 100">
            <el-progress 
              :percentage="100" 
              status="success"
              :stroke-width="8"
            />
            <div style="margin-top: 4px; font-size: 12px; color: #67C23A;">
              {{ row.installProgressText || '部署完成' }}
            </div>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="健康状态" width="100">
        <template #default="{row}">
          <el-tag
            :type="getHealthStatusType(row)"
            size="small"
          >
            <el-icon><component :is="getHealthStatusIcon(row)" /></el-icon>
            {{ getHealthStatusText(row) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="lastHeartbeat" label="最后心跳" width="180">
        <template #default="{row}">
          <span v-if="row.lastHeartbeat && row.lastHeartbeat !== '0001-01-01 00:00:00'">
            {{ formatTime(row.lastHeartbeat) }}
          </span>
          <span v-else class="text-muted">未连接</span>
        </template>
      </el-table-column>
      <el-table-column prop="updateTime" label="更新时间" width="180">
        <template #default="{row}">
          {{ formatTime(row.updateTime) }}
        </template>
      </el-table-column>
      <el-table-column label="操作"  fixed="right">
        <template #default="{row}">
          <el-button 
            size="small" 
            type="success" 
            icon="Download" 
            @click="handleRedeploy(row)"
            :disabled="row.status === 1"
            v-if="row.status === 2"
          >重新部署</el-button>
          
          <el-button 
            size="small" 
            type="warning" 
            icon="Refresh" 
            @click="handleRestart(row)"
            :disabled="row.status === 1"
            v-if="row.status === 4"
          >重启</el-button>
          
          <el-tooltip content="卸载" placement="top">
            <el-button
              size="small"
              type="warning"
              v-authority="['ops:agent:delete']"
              @click="handleUninstall(row)"
            >
              <img src="@/assets/image/卸载.svg" style="width: 16px; height: 16px"/>卸载
            </el-button>
          </el-tooltip>

          <el-button
            size="small"
            type="danger"
            icon="Delete"
            v-authority="['ops:agent:delete']"
            @click="handleDelete(row)"
          >删除</el-button>

          <el-button
            size="small"
            type="success"
            icon="View"
            v-authority="['ops:agent:get']"
            @click="viewAgentDetails(row)"
          >详情</el-button>
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
          :page-sizes="[10, 50, 100]"
          :page-size="queryParams.pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
        />
      </div>

    <!--Agent详情对话框-->
    <el-dialog
      title="Agent详情"
      v-model="detailVisible"
      width="50%"
      :modal="false"
    >
      <el-descriptions 
        :column="2" 
        border 
        v-if="currentAgent"
      >
        <el-descriptions-item label="主机名称">{{ currentAgent.hostName }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentAgent.sshIp }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ currentAgent.version }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <span :class="getStatusClass(currentAgent.status)">
            <el-icon><component :is="getStatusIcon(currentAgent.status)" /></el-icon>
            {{ getStatusText(currentAgent.status) }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="安装路径">{{ currentAgent.installPath || '未设置' }}</el-descriptions-item>
        <el-descriptions-item label="监听端口">{{ currentAgent.port }}</el-descriptions-item>
        <el-descriptions-item label="进程ID">{{ currentAgent.pid || '无' }}</el-descriptions-item>
        <el-descriptions-item label="健康状态">
          <el-tag 
            :type="currentAgent.isHealthy ? 'success' : 'danger'"
            size="small"
          >
            {{ currentAgent.isHealthy ? '健康' : '异常' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="最后心跳" span="2">
          <span v-if="currentAgent.lastHeartbeat && currentAgent.lastHeartbeat !== '0001-01-01 00:00:00'">
            {{ formatTime(currentAgent.lastHeartbeat) }}
          </span>
          <span v-else class="text-muted">未连接</span>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" span="2">
          {{ formatTime(currentAgent.createTime) }}
        </el-descriptions-item>
        <el-descriptions-item label="更新时间" span="2">
          {{ formatTime(currentAgent.updateTime) }}
        </el-descriptions-item>
        <el-descriptions-item label="错误信息" span="2" v-if="currentAgent.errorMsg">
          <el-alert 
            :title="currentAgent.errorMsg" 
            type="error" 
            :closable="false"
            show-icon
          />
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!--选择主机部署Agent对话框-->
    <AgentHostSelectorDialog
      v-model="showDeployHostDialog"
      title="选择部署Agent的主机"
      @hosts-selected="handleHostsSelected"
    />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Search, 
  Refresh, 
  Download, 
  Delete, 
  View,
  CircleCheck,
  CircleClose,
  Finished,
  Loading,
  Warning,
  Close
} from '@element-plus/icons-vue'
import cmdbAPI from '@/api/cmdb'
import AgentHostSelectorDialog from '@/components/AgentHostSelectorDialog.vue'

const agents = ref([])
const loading = ref(false)
const detailVisible = ref(false)
const currentAgent = ref(null)
const selectedAgents = ref([])
const total = ref(0)
const showDeployHostDialog = ref(false)
let pollingTimer = null
let pollingCounter = 0
const MAX_POLLING_COUNT = 5 // 最多轮询5次(15秒)

const queryParams = reactive({
  page: 1,
  pageSize: 10,
  hostName: '',
  status: null,
  version: ''
})

const fetchAgents = async () => {
  loading.value = true
  try {
    const response = await cmdbAPI.getAgentList({
      page: queryParams.page,
      pageSize: queryParams.pageSize,
      hostName: queryParams.hostName || undefined,
      status: queryParams.status || undefined,
      version: queryParams.version || undefined
    })
    
    if (response.data?.code === 200) {
      // 映射API响应数据到组件所需的格式
      agents.value = (response.data.data?.list || []).map(item => ({
        id: item.id,
        hostId: item.hostId,
        hostName: item.hostName,
        sshIp: item.sshIp,
        version: item.version,
        status: item.status,
        statusText: item.statusText,
        installPath: item.installPath,
        port: item.port,
        pid: item.pid,
        lastHeartbeat: item.lastHeartbeat,
        updateTime: item.updateTime,
        createTime: item.createTime,
        isHealthy: item.isHealthy,
        errorMsg: item.errorMsg,
        installProgress: item.installProgress,
        installProgressText: item.installProgressText
      }))
      total.value = response.data.data?.total || 0
    } else {
      throw new Error(response.data?.message || '获取Agent列表失败')
    }
  } catch (error) {
    console.error('获取Agent列表失败:', error)
    ElMessage.error(`获取Agent列表失败: ${error.message || '未知错误'}`)
    agents.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleRedeploy = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要重新部署Agent到主机 ${row.hostName} 吗？`, '确认重新部署', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    await cmdbAPI.deployAgent([row.hostId], '1.0.0')
    ElMessage.success('Agent重新部署已启动，请等待部署完成')
    fetchAgents()
    
    startPolling()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重新部署Agent失败:', error)
      ElMessage.error(`重新部署Agent失败: ${error.response?.data?.message || error.message || '未知错误'}`)
    }
  }
}

const handleDeployToHosts = () => {
  showDeployHostDialog.value = true
}

const handleHostsSelected = async (selectedHosts) => {
  if (!selectedHosts || selectedHosts.length === 0) {
    ElMessage.warning('请选择要部署Agent的主机')
    return
  }

  try {
    const hostNames = selectedHosts.map(host => host.name).join(', ')
    await ElMessageBox.confirm(`确定要在以下主机上部署Agent吗？\n${hostNames}`, '确认部署', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    const hostIds = selectedHosts.map(host => host.id)
    await cmdbAPI.deployAgent(hostIds, '1.0.0')
    
    ElMessage.success(`已开始在 ${selectedHosts.length} 台主机上部署Agent，请等待部署完成`)
    showDeployHostDialog.value = false
    fetchAgents()
    startPolling()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('部署Agent失败:', error)
      ElMessage.error(`部署Agent失败: ${error.response?.data?.message || error.message || '未知错误'}`)
    }
  }
}

const handleUninstall = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要卸载主机 ${row.hostName} 的Agent吗？`, '确认卸载', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await cmdbAPI.uninstallAgent([row.hostId])
    ElMessage.success('Agent卸载已启动，请等待卸载完成')

    // 立即刷新一次数据
    await fetchAgents()

    // 启动轮询以监控卸载进度
    startPolling()

    // 额外设置一个3秒后的刷新，确保状态更新
    setTimeout(() => {
      fetchAgents()
    }, 3000)

  } catch (error) {
    if (error !== 'cancel') {
      console.error('卸载Agent失败:', error)
      ElMessage.error(`卸载Agent失败: ${error.response?.data?.message || error.message || '未知错误'}`)
    }
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除主机 ${row.hostName} 的Agent数据吗？此操作不可恢复！`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await cmdbAPI.deleteAgent(row.id)
    ElMessage.success('Agent数据删除成功')

    // 立即刷新数据
    await fetchAgents()

  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除Agent失败:', error)
      ElMessage.error(`删除Agent失败: ${error.response?.data?.message || error.message || '未知错误'}`)
    }
  }
}

const handleRestart = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要重启主机 ${row.hostName} 的Agent吗？`, '确认重启', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })

    await cmdbAPI.restartAgent(row.hostId)
    ElMessage.success('Agent重启指令已发送')

    // 立即刷新数据
    await fetchAgents()

    // 启动轮询以监控重启后的状态变化
    startPolling()

  } catch (error) {
    if (error !== 'cancel') {
      console.error('重启Agent失败:', error)
      ElMessage.error(`重启Agent失败: ${error.response?.data?.message || error.message || '未知错误'}`)
    }
  }
}


const batchUninstall = async () => {
  try {
    const hostNames = selectedAgents.value.map(agent => agent.hostName).join(', ')
    await ElMessageBox.confirm(`确定要批量卸载以下主机的Agent吗？\n${hostNames}`, '确认批量卸载', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const hostIds = selectedAgents.value.map(agent => agent.hostId)
    await cmdbAPI.uninstallAgent(hostIds)

    ElMessage.success('批量卸载已启动，请等待卸载完成')

    // 立即刷新一次数据
    await fetchAgents()

    // 启动轮询以监控卸载进度
    startPolling()

    // 额外设置一个3秒后的刷新，确保状态更新
    setTimeout(() => {
      fetchAgents()
    }, 3000)

  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量卸载失败:', error)
      ElMessage.error(`批量卸载失败: ${error.message || '未知错误'}`)
    }
  }
}

const viewAgentDetails = (row) => {
  currentAgent.value = row
  detailVisible.value = true
}

const handleSelectionChange = (selection) => {
  selectedAgents.value = selection
}

const searchByHostName = () => {
  queryParams.page = 1
  fetchAgents()
}

const searchByStatus = () => {
  queryParams.page = 1
  fetchAgents()
}

const searchByVersion = () => {
  queryParams.page = 1
  fetchAgents()
}

const resetQuery = () => {
  queryParams.hostName = ''
  queryParams.status = null
  queryParams.version = ''
  queryParams.page = 1
  fetchAgents()
}

const handleSizeChange = (val) => {
  queryParams.pageSize = val
  queryParams.page = 1
  fetchAgents()
}

const handleCurrentChange = (val) => {
  queryParams.page = val
  fetchAgents()
}

const getStatusClass = (status) => {
  const statusClasses = {
    1: 'status-deploying',
    2: 'status-deploy-failed', 
    3: 'status-running',
    4: 'status-exception'
  }
  return statusClasses[status] || ''
}

const getStatusIcon = (status) => {
  const statusIcons = {
    1: 'Loading',
    2: 'Close', 
    3: 'CircleCheck',
    4: 'Warning'
  }
  return statusIcons[status] || 'Warning'
}

const getStatusText = (status) => {
  const statusTexts = {
    1: '部署中',
    2: '部署失败', 
    3: '运行中',
    4: '启动异常'
  }
  return statusTexts[status] || '未知状态'
}

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

const getHealthStatusText = (row) => {
  // status: 1-部署中, 2-部署失败, 3-运行中, 4-启动异常
  switch (row.status) {
    case 1:
      return '部署中'
    case 2:
      return '异常'
    case 3:
      return '健康'
    case 4:
      return '异常'
    default:
      return '未知'
  }
}

const getHealthStatusType = (row) => {
  switch (row.status) {
    case 1:
      return 'warning'
    case 2:
      return 'danger'
    case 3:
      return 'success'
    case 4:
      return 'danger'
    default:
      return 'info'
  }
}

const getHealthStatusIcon = (row) => {
  switch (row.status) {
    case 1:
      return 'Loading'
    case 2:
      return 'CircleClose'
    case 3:
      return 'CircleCheck'
    case 4:
      return 'CircleClose'
    default:
      return 'Warning'
  }
}

const startPolling = () => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
  }

  // 重置计数器
  pollingCounter = 0

  pollingTimer = setInterval(() => {
    const hasActiveOperations = agents.value.some(agent =>
      agent.status === 1 // 部署中状态需要轮询
    )

    if (hasActiveOperations) {
      // 有活跃操作时重置计数器并继续轮询
      pollingCounter = 0
      fetchAgents()
    } else {
      // 没有活跃操作时增加计数器
      pollingCounter++

      // 继续轮询几次以确保状态同步（比如卸载操作完成后的状态更新）
      if (pollingCounter <= MAX_POLLING_COUNT) {
        fetchAgents()
      } else {
        // 超过最大轮询次数后停止
        stopPolling()
      }
    }
  }, 3000)
}

const stopPolling = () => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
}

onMounted(() => {
  fetchAgents()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.agent-management {
  padding: 5px;
  height: 80vh;
  background: #eef1f6;
}

.agent-card {
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

.action-section {
  margin-bottom: 5px;
  margin-top: 5px;
  padding: 12px 0;
}

.action-section .el-button {
  margin-right: 12px;
}

.table-section {
  margin-bottom: 20px;
  flex: 1;
  overflow-y: auto;
}

.agent-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.agent-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.agent-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.agent-table :deep(.el-table__header th .cell) {
  color: #2c3e50 !important;
  font-weight: 700 !important;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.agent-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.agent-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.pagination-section {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.text-muted {
  color: #909399;
}

.status-deploying {
  color: #E6A23C;
  font-weight: 500;
  animation: pulse 1.5s ease-in-out infinite alternate;
}

.status-deploy-failed {
  color: #F56C6C;
  font-weight: 500;
}

.status-running {
  color: #67C23A;
  font-weight: 500;
}

.status-exception {
  color: #F56C6C;
  font-weight: 500;
}

@keyframes pulse {
  from {
    opacity: 1;
  }
  to {
    opacity: 0.5;
  }
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

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #2c3e50;
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

.agent-table .el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
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

.el-progress {
  line-height: 1;
}
</style>
