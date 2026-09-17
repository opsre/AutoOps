<template>
  <div class="cluster-table-container">
    <!-- 搜索表单 -->
    <div class="search-section">
      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="集群名称">
          <el-input
            v-model="queryParams.clusterName"
            placeholder="请输入集群名称"
            clearable
            size="small"
            style="width: 200px"
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item label="集群状态">
          <el-select
            v-model="queryParams.status"
            placeholder="请选择状态"
            clearable
            size="small"
            style="width: 150px"
          >
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="集群版本">
          <el-input
            v-model="queryParams.version"
            placeholder="请输入版本"
            clearable
            size="small"
            style="width: 150px"
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" size="small" @click="handleQuery">
            搜索
          </el-button>
          <el-button type="warning" :icon="Refresh" size="small" @click="resetQuery">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 集群列表表格 -->
    <div class="table-section">
      <el-table
        :data="tableData"
        v-loading="loading"
        stripe
        style="width: 100%"
        class="cluster-table"
      >
        <el-table-column prop="clusterName" label="集群名称" min-width="180">
          <template #default="{ row }">
            <div class="cluster-name-container">
              <img src="@/assets/image/k8s.svg" alt="k8s" class="k8s-icon" />
              <el-button
                type="primary"
                link
                @click="handleClusterNameClick(row)"
                class="cluster-name-link"
              >
                {{ row.clusterName }}
              </el-button>
              <el-tooltip content="复制集群名称" placement="top">
                <el-button
                  type="info"
                  :icon="DocumentCopy"
                  size="small"
                  circle
                  @click.stop="copyClusterName(row)"
                  class="copy-name-btn"
                />
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="集群版本" min-width="140">
          <template #default="{ row }">
            <div class="version-container">
            <img src="@/assets/image/k8s-io.svg" alt="k8s-io" class="k8s-io-icon" />
            <span class="version-text">{{ row.version }}</span>
              
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="nodeCount" label="节点数量" min-width="110">
          <template #default="{ row }">
            <el-tag size="small" type="info">
              {{ row.nodeCount }} 个
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="集群状态" min-width="110">
          <template #default="{ row }">
            <el-tag
              :type="getStatusTag(row.status)"
              size="small"
              effect="dark"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="集群凭证" width="100" align="center">
          <template #default="{ row }">
            <el-tooltip 
              content="查看集群凭证" 
              placement="top"
              v-if="row.credential && row.credential.trim()"
            >
              <el-button
                type="primary"
                :icon="View"
                size="small"
                circle
                v-authority="['cloud:k8s:get']"
                @click="viewCredential(row)"
              />
            </el-tooltip>
            <el-tag v-else size="small" type="info">未设置</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="备注信息" min-width="160">
          <template #default="{ row }">
            <div class="description-container">
              <el-tooltip 
                :content="row.description || '暂无备注'"
                placement="top"
                :disabled="!row.description || row.description.length <= 20"
              >
                <span class="description-text">
                  {{ row.description ? (row.description.length > 20 ? row.description.substring(0, 20) + '...' : row.description) : '暂无备注' }}
                </span>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" min-width="180" />
        <el-table-column label="操作" width="230" fixed="right">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-button
                type="primary"
                :icon="Connection"
                size="small"
                v-authority="['cloud:k8s:rsync']"
                :disabled="!row.credential || !row.credential.trim()"
                @click="handleSync(row)"
              >同步</el-button>
              <el-button
                type="warning"
                :icon="Edit"
                size="small"
                v-authority="['cloud:k8s:edit']"
                @click="handleEdit(row)"
              >编辑</el-button>
              <el-button
                type="danger"
                :icon="Delete"
                size="small"
                v-authority="['cloud:k8s:delete']"
                @click="handleDelete(row)"
              >删除</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Edit,
  Delete,
  View,
  Connection,
  DocumentCopy
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'

const router = useRouter()

const props = defineProps({
  refreshTrigger: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['cluster-sync', 'cluster-edit', 'cluster-delete', 'view-credential'])

const loading = ref(false)
const queryParams = reactive({
  clusterName: '',
  status: '',
  version: ''
})

const statusOptions = [
  { label: '全部', value: '' },
  { label: '创建中', value: 1 },
  { label: '运行中', value: 2 },
  { label: '离线', value: 3 }
]

const tableData = ref([])

const handleQuery = async () => {
  try {
    loading.value = true
    
    // 只有当查询参数不为空时才传递参数
    const params = {}
    if (queryParams.clusterName) params.name = queryParams.clusterName
    if (queryParams.version) params.version = queryParams.version  
    if (queryParams.status) params.status = queryParams.status
    
    const response = await k8sApi.getClusterList(Object.keys(params).length > 0 ? params : undefined)
    
    // response可能是axios响应对象，需要获取data部分
    const responseData = response.data || response
    console.log('API响应数据:', responseData)
    
    if (responseData.code === 200 || responseData.success) {
      // 根据实际API响应结构处理数据
      const clusters = responseData.data?.list || responseData.data || []
      tableData.value = clusters.map(cluster => ({
        id: cluster.id,
        clusterName: cluster.name,
        version: cluster.version || '', // 直接使用后端返回的版本信息，不添加前缀
        nodeCount: cluster.nodeCount || 0,
        status: getClusterStatus(cluster.status),
        credential: cluster.credential || '',  // 保留原始凭证数据
        description: cluster.description || cluster.remark || '', // 备注信息
        createTime: cluster.createdAt ? new Date(cluster.createdAt).toLocaleString() : cluster.createTime
      }))
      
      console.log('集群列表加载成功:', tableData.value)
    } else {
      ElMessage.error(responseData.message || '获取集群列表失败')
      tableData.value = []
    }
  } catch (error) {
    console.error('获取集群列表失败:', error)
    
    // 检查是否是网络连接问题
    if (error.code === 'ERR_NETWORK' || 
        error.message?.includes('ERR_CONNECTION_REFUSED') ||
        error.message?.includes('Failed to fetch')) {
      ElMessage.warning('后端服务连接失败，请检查服务状态')
    } else if (error.response?.status === 401) {
      ElMessage.error('认证失败，请重新登录')
    } else if (error.response?.status === 403) {
      ElMessage.error('权限不足，请联系管理员')
    } else {
      // 只在真正异常时才显示错误
      console.warn('API调用异常，但可能数据已正确加载')
    }
    
    // 设置空数据，避免界面显示异常
    tableData.value = []
  } finally {
    loading.value = false
  }
}

const getClusterStatus = (statusCode) => {
  const statusMap = {
    1: 'creating',
    2: 'online', 
    3: 'offline'
  }
  return statusMap[statusCode] || 'unknown'
}

const resetQuery = () => {
  queryParams.clusterName = ''
  queryParams.status = ''
  queryParams.version = ''
  // 重置后自动刷新数据
  handleQuery()
}

const handleClusterNameClick = (row) => {
  // 跳转到集群详情页面
  router.push(`/k8s/cluster/${row.id}`)
}

// 复制集群名称
const copyClusterName = async (row) => {
  try {
    await navigator.clipboard.writeText(row.clusterName)
    ElMessage.success(`集群名称 "${row.clusterName}" 已复制到剪贴板`)
  } catch (error) {
    console.error('复制失败:', error)
    // 降级处理：使用传统方式复制
    try {
      const textArea = document.createElement('textarea')
      textArea.value = row.clusterName
      textArea.style.position = 'fixed'
      textArea.style.opacity = '0'
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      ElMessage.success(`集群名称 "${row.clusterName}" 已复制到剪贴板`)
    } catch (fallbackError) {
      console.error('降级复制也失败:', fallbackError)
      ElMessage.error('复制失败，请手动复制')
    }
  }
}

const handleSync = async (row) => {
  try {
    // 检查是否有集群凭证
    if (!row.credential || !row.credential.trim()) {
      ElMessage.warning('该集群没有配置凭证，无法进行同步操作')
      return
    }
    
    const response = await k8sApi.syncCluster(row.id)
    
    // response可能是axios响应对象，需要获取data部分
    const responseData = response.data || response
    console.log('同步集群API响应:', responseData)
    
    if (responseData.code === 200 || responseData.success || responseData.message === '成功') {
      ElMessage.success(`集群 ${row.clusterName} 同步成功`)
      // 触发父组件刷新
      emit('cluster-sync', row)
      // 立即刷新一次
      await handleQuery()
    } else if (responseData.code === 500) {
      // 处理超时或连接失败的情况
      const isTimeoutError = responseData.message && (
        responseData.message.includes('timeout') ||
        responseData.message.includes('Timeout exceeded') ||
        responseData.message.includes('connection') ||
        responseData.message.includes('已将集群状态设为离线')
      )
      
      if (isTimeoutError) {
        ElMessage.warning(`集群 ${row.clusterName} 连接超时，已设为离线状态`)
        console.log('检测到超时错误，准备延迟刷新状态')
        
        // 立即刷新一次，可能状态还没更新
        await handleQuery()
        
        // 延迟刷新，确保后端状态已更新
        setTimeout(async () => {
          console.log('延迟刷新集群状态')
          await handleQuery()
        }, 1000) // 1秒后再次刷新
      } else {
        ElMessage.error(responseData.message || '同步失败')
      }
      
      // 触发父组件刷新
      emit('cluster-sync', row)
    } else {
      ElMessage.error(responseData.message || '同步失败')
      // 即使失败也要刷新，可能状态有变化
      await handleQuery()
      emit('cluster-sync', row)
    }
  } catch (error) {
    console.error('同步集群失败:', error)
    
    // 检查错误信息是否包含超时相关内容
    const errorMessage = error.message || error.toString()
    const isNetworkTimeout = errorMessage.includes('timeout') || 
                            errorMessage.includes('Network Error') ||
                            errorMessage.includes('ERR_NETWORK')
    
    if (isNetworkTimeout) {
      ElMessage.warning(`集群 ${row.clusterName} 网络连接超时`)
      // 网络超时也要延迟刷新状态
      setTimeout(async () => {
        console.log('网络超时后延迟刷新')
        await handleQuery()
      }, 1000)
    } else {
      ElMessage.error('同步集群失败，请检查网络连接')
    }
    
    emit('cluster-sync', row)
  }
}

const handleEdit = async (row) => {
  emit('cluster-edit', row)
}

const viewCredential = (row) => {
  emit('view-credential', row)
}

const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除集群 ${row.clusterName} 吗？此操作不可逆！`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  )
  .then(async () => {
    try {
      const response = await k8sApi.deleteCluster(row.id)
      
      // response可能是axios响应对象，需要获取data部分
      const responseData = response.data || response
      console.log('删除集群API响应:', responseData)
      
      if (responseData.code === 200 || responseData.success || responseData.message === '成功') {
        ElMessage.success(`集群 ${row.clusterName} 已删除`)
        emit('cluster-delete', row)
        handleQuery() // 刷新列表
      } else {
        ElMessage.error(responseData.message || '删除失败')
      }
    } catch (error) {
      console.error('删除集群失败:', error)
      ElMessage.error('删除集群失败，请检查网络连接')
    }
  })
  .catch(() => {
    ElMessage.info('已取消删除')
  })
}

const getStatusTag = (status) => {
  const tagMap = {
    'creating': 'info',
    'online': 'success',
    'offline': 'warning'
  }
  return tagMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    'creating': '创建中',
    'online': '运行中',
    'offline': '离线'
  }
  return textMap[status] || '未知'
}

// 监听刷新触发器
watch(() => props.refreshTrigger, (newVal, oldVal) => {
  console.log(`表格刷新触发器变化: ${oldVal} -> ${newVal}`)
  handleQuery()
})

// 暴露方法给父组件调用
defineExpose({
  handleQuery,
  resetQuery
})

onMounted(() => {
  handleQuery()
})
</script>

<style scoped>
.cluster-table-container {
  width: 100%;
}

.search-section {
  margin-bottom: 24px;
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
  margin-top: 20px;
}

.cluster-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.cluster-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.cluster-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #000000 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: none;
}

.cluster-table :deep(.el-table__header th .cell) {
  color: #000000 !important;
  font-weight: 700 !important;
  text-shadow: none;
}

.cluster-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.cluster-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.cluster-name-container {
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
}

.k8s-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.version-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.version-text {
  color: #606266;
  font-weight: 500;
}

.k8s-io-icon {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
  filter: drop-shadow(0 1px 3px rgba(0, 0, 0, 0.12));
}

.cluster-name-link {
  font-weight: 600;
  color: #667eea;
  text-decoration: none;
  transition: all 0.3s ease;
}

.cluster-name-link:hover {
  color: #764ba2;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.copy-name-btn {
  opacity: 0;
  transition: all 0.3s ease;
  transform: scale(0.8);
}

.cluster-name-container:hover .copy-name-btn {
  opacity: 1;
  transform: scale(1);
}

.copy-name-btn:hover {
  background-color: #409eff !important;
  border-color: #409eff !important;
  color: white !important;
  transform: scale(1.1) !important;
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

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #6b9bd1;
}

.el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

/* 备注信息样式 */
.description-container {
  max-width: 160px;
}

.description-text {
  color: #606266;
  font-size: 13px;
  line-height: 1.4;
  display: inline-block;
  max-width: 100%;
  word-break: break-word;
  cursor: default;
}

.description-text:empty::before {
  content: '暂无备注';
  color: #c0c4cc;
  font-style: italic;
}
</style>