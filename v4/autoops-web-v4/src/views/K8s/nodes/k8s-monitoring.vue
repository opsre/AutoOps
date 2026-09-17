<template>
  <div class="k8s-monitoring-dashboard">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-title">
        <h2>K8s 监控仪表板</h2>
        <p class="header-desc">实时监控集群资源使用情况</p>
      </div>
      <div class="header-controls">
        <el-select v-model="selectedClusterId" placeholder="选择集群" @change="handleClusterChange">
          <el-option 
            v-for="cluster in clusterList" 
            :key="cluster.id" 
            :label="cluster.name" 
            :value="cluster.id"
          />
        </el-select>
        <el-button :icon="Refresh" @click="refreshAllData" :loading="loading">刷新数据</el-button>
      </div>
    </div>

    <!-- 监控标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <!-- 节点监控 -->
      <el-tab-pane label="节点监控" name="nodes">
        <NodesMonitoring 
          :selected-cluster-id="selectedClusterId"
          ref="nodesMonitoringRef"
        />
      </el-tab-pane>

      <!-- 命名空间监控 -->
      <el-tab-pane label="命名空间监控" name="namespaces">
        <NamespaceMonitoring 
          :selected-cluster-id="selectedClusterId"
          ref="namespaceMonitoringRef"
        />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  View,
  Monitor
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import NodesMonitoring from './NodesMonitoring.vue'
import NamespaceMonitoring from './NamespaceMonitoring.vue'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const activeTab = ref('nodes')
const selectedClusterId = ref('')
const clusterList = ref([])

// 组件引用
const nodesMonitoringRef = ref(null)
const namespaceMonitoringRef = ref(null)

// 工具函数
const formatTimestamp = (timestamp) => {
  if (!timestamp) return '-'
  
  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    return timestamp
  }
}

const getProgressColor = (percentage) => {
  if (percentage >= 90) return '#f56c6c'
  if (percentage >= 75) return '#e6a23c'
  if (percentage >= 60) return '#409eff'
  return '#67c23a'
}

// API调用函数
const fetchClusterList = async () => {
  try {
    const response = await k8sApi.getClusterList()
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      const clusters = responseData.data?.list || responseData.data || []
      clusterList.value = clusters.map(cluster => ({
        id: cluster.id,
        name: cluster.name,
        status: cluster.status
      }))
      
      if (clusterList.value.length > 0 && !selectedClusterId.value) {
        const onlineCluster = clusterList.value.find(cluster => cluster.status === 2)
        selectedClusterId.value = onlineCluster ? onlineCluster.id : clusterList.value[0].id
      }
    } else {
      ElMessage.error(responseData.message || '获取集群列表失败')
    }
  } catch (error) {
    console.error('获取集群列表失败:', error)
    ElMessage.warning('无法获取集群列表，请检查后端服务')
  }
}


// 事件处理函数
const handleClusterChange = async () => {
  if (selectedClusterId.value) {
    await loadTabData()
    // 通知组件集群变更
    if (nodesMonitoringRef.value) {
      nodesMonitoringRef.value.handleClusterChange()
    }
    if (namespaceMonitoringRef.value) {
      namespaceMonitoringRef.value.handleClusterChange()
    }
  }
}

const handleTabChange = async (tabName) => {
  activeTab.value = tabName
  await loadTabData()
}

const loadTabData = async () => {
  if (!selectedClusterId.value) return
  
  if (activeTab.value === 'nodes') {
    // 节点监控由 NodesMonitoring 组件处理
    if (nodesMonitoringRef.value) {
      nodesMonitoringRef.value.refreshAllNodes()
    }
  } else if (activeTab.value === 'namespaces') {
    // 命名空间监控由 NamespaceMonitoring 组件处理
    if (namespaceMonitoringRef.value) {
      namespaceMonitoringRef.value.refreshAllData()
    }
  }
}

const refreshAllData = async () => {
  loading.value = true
  try {
    if (activeTab.value === 'nodes' && nodesMonitoringRef.value) {
      await nodesMonitoringRef.value.refreshAllNodes()
    } else if (activeTab.value === 'namespaces' && namespaceMonitoringRef.value) {
      await namespaceMonitoringRef.value.refreshAllData()
    } else {
      await loadTabData()
    }
    ElMessage.success('数据刷新成功')
  } catch (error) {
    ElMessage.error('数据刷新失败')
  } finally {
    loading.value = false
  }
}


// 组件挂载时初始化
onMounted(async () => {
  try {
    await fetchClusterList()
    if (selectedClusterId.value) {
      await loadTabData()
    }
  } catch (error) {
    console.error('页面初始化失败:', error)
  }
})
</script>

<style scoped>
.k8s-monitoring-dashboard {
  padding: 24px;
  background: #f5f5f5;
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  background: white;
  padding: 20px 24px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-title h2 {
  margin: 0 0 4px 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.header-desc {
  margin: 0;
  color: #606266;
  font-size: 14px;
}

.header-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.header-controls .el-select {
  width: 200px;
}

</style>