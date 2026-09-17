<template>
  <div class="pod-monitor">
    <!-- 监控对话框 -->
    <el-dialog
      v-model="monitoringVisible"
      :title="`Pod监控 - ${selectedPod?.name || ''}`"
      width="60%"
      class="monitoring-dialog"
    >
      <div class="monitoring-container" v-loading="monitoringLoading">
        <el-row :gutter="16" style="margin-bottom: 16px;">
          <el-col :span="12">
            <el-card class="metric-card">
              <div class="metric-header">
                <span class="metric-title">CPU使用率</span>
                <span class="metric-value">{{ monitoringData.cpu.percentage }}%</span>
              </div>
              <el-progress 
                :percentage="monitoringData.cpu.percentage" 
                :color="monitoringData.cpu.percentage > 80 ? '#F56C6C' : monitoringData.cpu.percentage > 60 ? '#E6A23C' : '#67C23A'"
                :stroke-width="8"
              />
              <div class="metric-details">
                <span>已用: {{ monitoringData.cpu.used }}</span>
                <span>限制: {{ monitoringData.cpu.limit }}</span>
              </div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card class="metric-card">
              <div class="metric-header">
                <span class="metric-title">内存使用率</span>
                <span class="metric-value">{{ monitoringData.memory.percentage }}%</span>
              </div>
              <el-progress 
                :percentage="monitoringData.memory.percentage" 
                :color="monitoringData.memory.percentage > 80 ? '#F56C6C' : monitoringData.memory.percentage > 60 ? '#E6A23C' : '#67C23A'"
                :stroke-width="8"
              />
              <div class="metric-details">
                <span>已用: {{ monitoringData.memory.used }}</span>
                <span>限制: {{ monitoringData.memory.limit }}</span>
              </div>
            </el-card>
          </el-col>
        </el-row>
        
        <el-row :gutter="16">
          <el-col :span="12">
            <el-card class="metric-card">
              <div class="metric-header">
                <span class="metric-title">网络流量</span>
              </div>
              <div class="network-metrics">
                <div class="network-item">
                  <span class="network-label">入站:</span>
                  <span class="network-value">{{ monitoringData.network.rx }}</span>
                </div>
                <div class="network-item">
                  <span class="network-label">出站:</span>
                  <span class="network-value">{{ monitoringData.network.tx }}</span>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card class="metric-card">
              <div class="metric-header">
                <span class="metric-title">磁盘使用率</span>
                <span class="metric-value">{{ monitoringData.disk.percentage }}%</span>
              </div>
              <el-progress 
                :percentage="monitoringData.disk.percentage" 
                :color="monitoringData.disk.percentage > 80 ? '#F56C6C' : monitoringData.disk.percentage > 60 ? '#E6A23C' : '#67C23A'"
                :stroke-width="8"
              />
              <div class="metric-details">
                <span>已用: {{ monitoringData.disk.used }}</span>
                <span>限制: {{ monitoringData.disk.limit }}</span>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 容器级别监控数据 -->
        <el-card v-if="monitoringData.containers && monitoringData.containers.length > 0" style="margin-top: 16px;">
          <template #header>
            <div class="card-header">
              <span>容器监控详情</span>
              <el-tag v-if="monitoringData.timestamp" type="info" size="small">
                更新时间: {{ formatTimestamp(monitoringData.timestamp) }}
              </el-tag>
            </div>
          </template>
          
          <el-table :data="monitoringData.containers" size="small" stripe>
            <el-table-column prop="name" label="容器名称" min-width="120" />
            <el-table-column label="CPU使用" min-width="100">
              <template #default="{ row }">
                <div class="container-metric">
                  <span class="metric-text">{{ row.usage?.cpu || '0m' }}</span>
                  <el-progress 
                    v-if="row.usageRate?.cpuRate" 
                    :percentage="Math.round(row.usageRate.cpuRate)" 
                    :stroke-width="6"
                    :show-text="false"
                    :color="row.usageRate.cpuRate > 80 ? '#F56C6C' : row.usageRate.cpuRate > 60 ? '#E6A23C' : '#67C23A'"
                  />
                </div>
              </template>
            </el-table-column>
            <el-table-column label="内存使用" min-width="100">
              <template #default="{ row }">
                <div class="container-metric">
                  <span class="metric-text">{{ row.usage?.memory || '0Mi' }}</span>
                  <el-progress 
                    v-if="row.usageRate?.memoryRate" 
                    :percentage="Math.round(row.usageRate.memoryRate)" 
                    :stroke-width="6"
                    :show-text="false"
                    :color="row.usageRate.memoryRate > 80 ? '#F56C6C' : row.usageRate.memoryRate > 60 ? '#E6A23C' : '#67C23A'"
                  />
                </div>
              </template>
            </el-table-column>
            <el-table-column label="CPU使用率" min-width="80">
              <template #default="{ row }">
                <el-tag 
                  :type="row.usageRate?.cpuRate > 80 ? 'danger' : row.usageRate?.cpuRate > 60 ? 'warning' : 'success'"
                  size="small"
                >
                  {{ Math.round(row.usageRate?.cpuRate || 0) }}%
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="内存使用率" min-width="80">
              <template #default="{ row }">
                <el-tag 
                  :type="row.usageRate?.memoryRate > 80 ? 'danger' : row.usageRate?.memoryRate > 60 ? 'warning' : 'success'"
                  size="small"
                >
                  {{ Math.round(row.usageRate?.memoryRate || 0) }}%
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" min-width="80">
              <template #default="{ row }">
                <el-tag 
                  :type="row.state === 'Running' ? 'success' : row.state === 'Waiting' ? 'warning' : 'danger'"
                  size="small"
                >
                  {{ row.state || 'Unknown' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
      
      <template #footer>
        <el-button @click="monitoringVisible = false">关闭</el-button>
        <el-button type="primary" @click="refreshMonitoring" :loading="monitoringLoading">刷新数据</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import k8sApi from '@/api/k8s'

// Props
const props = defineProps({
  clusterId: {
    type: String,
    required: true
  },
  namespace: {
    type: String,
    required: true
  }
})

// 监控相关状态
const monitoringVisible = ref(false)
const monitoringLoading = ref(false)
const selectedPod = ref(null)
const monitoringData = ref({
  cpu: { used: '0', limit: '1000m', percentage: 0 },
  memory: { used: '0Mi', limit: '512Mi', percentage: 0 },
  network: { rx: '0KB/s', tx: '0KB/s' },
  disk: { used: '0GB', limit: '10GB', percentage: 0 },
  containers: [],
  timestamp: null
})

// 格式化时间戳
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

// 显示监控弹框
const showMonitoring = async (pod) => {
  selectedPod.value = pod
  monitoringVisible.value = true
  
  // 加载监控数据
  await loadMonitoringData(pod)
}

// 加载监控数据
const loadMonitoringData = async (pod) => {
  monitoringLoading.value = true
  
  try {
    console.log('🔍 加载Pod监控数据:', pod.name)
    
    // 调用真实的监控API
    const response = await k8sApi.getPodMetrics(props.clusterId, props.namespace, pod.name)
    const responseData = response.data || response
    
    if (responseData.code === 200 && responseData.data) {
      const metrics = responseData.data
      
      // 转换API响应数据格式
      monitoringData.value = {
        cpu: {
          used: metrics.totalUsage?.cpu || '0m',
          limit: metrics.resourceQuota?.cpu || '1000m',
          percentage: Math.round(metrics.usageRate?.cpuRate || 0)
        },
        memory: {
          used: metrics.totalUsage?.memory || '0Mi',
          limit: metrics.resourceQuota?.memory || '512Mi',
          percentage: Math.round(metrics.usageRate?.memoryRate || 0)
        },
        network: {
          rx: '0KB/s',
          tx: '0KB/s'
        },
        disk: {
          used: '0GB',
          limit: '10GB',
          percentage: 0
        },
        containers: metrics.containers || [],
        timestamp: metrics.timestamp
      }
      
      console.log('✅ 监控数据加载完成:', monitoringData.value)
    } else {
      console.warn('⚠️ 监控API返回异常:', responseData.message)
      ElMessage.warning('监控数据获取失败: ' + (responseData.message || '未知错误'))
    }

  } catch (error) {
    console.error('⚠️ 监控API调用失败:', error)
    ElMessage.error('监控数据获取失败，请检查网络连接或稍后重试')
  } finally {
    monitoringLoading.value = false
  }
}


// 刷新监控数据
const refreshMonitoring = async () => {
  if (selectedPod.value) {
    await loadMonitoringData(selectedPod.value)
  }
}

// 暴露方法给父组件
defineExpose({
  showMonitoring
})
</script>

<style scoped>
/* 监控弹框样式 */
.monitoring-dialog .el-dialog__body {
  padding: 20px;
}

.monitoring-container {
  min-height: 200px;
}

.metric-card {
  height: 150px;
  display: flex;
  flex-direction: column;
}

.metric-card .el-card__body {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16px;
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.metric-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.metric-value {
  font-size: 18px;
  font-weight: bold;
  color: #409eff;
}

.metric-details {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #666;
  margin-top: 8px;
}

.network-metrics {
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 60px;
}

.network-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.network-label {
  font-size: 14px;
  color: #666;
}

.network-value {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

/* 容器监控表格样式 */
.container-metric {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.metric-text {
  font-size: 12px;
  color: #606266;
  margin-bottom: 2px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  font-weight: 500;
}
</style>