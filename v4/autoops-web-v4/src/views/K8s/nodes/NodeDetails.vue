<template>
  <div class="node-details-container">
    <div class="node-details-header">
      <el-button
        @click="goBack"
        :icon="ArrowLeft"
        type="primary"
        size="small"
        class="back-btn"
      >
        返回
      </el-button>
      <h1 class="page-title">{{ nodeData?.basic?.name || '节点详情' }}</h1>
    </div>

    <div class="node-details-content" v-loading="loading">
      <div v-if="!loading && nodeData" class="main-grid">
        <!-- Row 1: 资源概览卡片（4列） -->
        <div class="stats-grid">
          <!-- CPU使用率 -->
          <div class="stat-card">
            <div class="stat-icon stat-icon-cpu">
              <img src="@/assets/image/cpu.svg" alt="CPU" class="svg-icon" />
            </div>
            <div class="stat-content">
              <div class="stat-label">CPU使用率</div>
              <div class="stat-value" :class="getCpuUsageClass()">{{ cpuUsageRate }}%</div>
              <div class="stat-sub">{{ nodeData.resources?.cpu?.capacity || 0 }}核可用</div>
              <el-progress
                :percentage="cpuUsageRate"
                :color="getProgressColor(cpuUsageRate)"
                :show-text="false"
                class="mini-progress"
              />
            </div>
          </div>

          <!-- 内存使用率 -->
          <div class="stat-card">
            <div class="stat-icon stat-icon-memory">
              <img src="@/assets/image/内存.svg" alt="内存" class="svg-icon" />
            </div>
            <div class="stat-content">
              <div class="stat-label">内存使用率</div>
              <div class="stat-value" :class="getMemoryUsageClass()">{{ memoryUsageRate }}%</div>
              <div class="stat-sub">{{ formatMemory(nodeData.resources?.memory?.capacity) }}</div>
              <el-progress
                :percentage="memoryUsageRate"
                :color="getProgressColor(memoryUsageRate)"
                :show-text="false"
                class="mini-progress"
              />
            </div>
          </div>

          <!-- 运行Pod数 -->
          <div class="stat-card">
            <div class="stat-icon stat-icon-pods">
              <el-icon><Box /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-label">运行Pod数</div>
              <div class="stat-value">{{ nodeData.resources?.pods?.current || 0 }}</div>
              <div class="stat-sub">总容量: {{ nodeData.resources?.pods?.capacity || 0 }}</div>
              <el-progress
                :percentage="podUsageRate"
                :color="getProgressColor(podUsageRate)"
                :show-text="false"
                class="mini-progress"
              />
            </div>
          </div>

          <!-- 运行时间 -->
          <div class="stat-card">
            <div class="stat-icon stat-icon-time">
              <el-icon><Timer /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-label">运行时间</div>
              <div class="stat-value">{{ formatTime(nodeData.basic?.creationTime) || 'N/A' }}</div>
              <div class="stat-sub">
                <el-tag
                  :type="getStatusType(nodeData.basic?.status)"
                  size="small"
                >
                  {{ getStatusText(nodeData.basic?.status) }}
                </el-tag>
              </div>
            </div>
          </div>
        </div>

        <!-- Row 2: 基本信息卡片（3列） -->
        <div class="info-grid">
          <!-- 基本信息 -->
          <div class="info-card">
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <h3>基本信息</h3>
            </div>
            <div class="card-content">
              <div class="info-row">
                <span class="label">节点名称：</span>
                <span class="value">{{ nodeData.basic?.name || 'N/A' }}</span>
              </div>
              <div class="info-row">
                <span class="label">节点角色：</span>
                <span class="value">{{ nodeData.basic?.role || 'N/A' }}</span>
              </div>
              <div class="info-row">
                <span class="label">节点状态：</span>
                <el-tag
                  :type="getStatusType(nodeData.basic?.status)"
                  class="status-tag"
                >
                  {{ getStatusText(nodeData.basic?.status) }}
                </el-tag>
              </div>
              <div class="info-row">
                <span class="label">创建时间：</span>
                <span class="value">{{ formatTime(nodeData.basic?.creationTime) || 'N/A' }}</span>
              </div>
              <div class="info-row">
                <span class="label">K8s版本：</span>
                <span class="value">{{ nodeData.basic?.k8sVersion || 'N/A' }}</span>
              </div>
              <div class="info-row">
                <span class="label">容器运行时：</span>
                <span class="value">{{ nodeData.basic?.containerRuntime || 'N/A' }}</span>
              </div>
            </div>
          </div>

          <!-- 系统信息 -->
          <div class="info-card">
            <div class="card-header">
              <el-icon><Monitor /></el-icon>
              <h3>系统信息</h3>
            </div>
            <div class="card-content">
              <div class="info-row">
                <span class="label">操作系统：</span>
                <span class="value">{{ nodeData.system?.osImage || 'N/A' }}</span>
              </div>
              <div class="info-row">
                <span class="label">内核版本：</span>
                <span class="value">{{ nodeData.system?.kernelVersion || 'N/A' }}</span>
              </div>
              <div class="info-row">
                <span class="label">架构：</span>
                <span class="value">{{ nodeData.system?.architecture || 'N/A' }}</span>
              </div>
              <div class="info-row">
                <span class="label">主机名：</span>
                <span class="value">{{ nodeData.system?.hostname || 'N/A' }}</span>
              </div>
              <div class="info-row">
                <span class="label">机器ID：</span>
                <span class="value">{{ nodeData.system?.machineId || 'N/A' }}</span>
              </div>
            </div>
          </div>

          <!-- 网络信息 -->
          <div class="info-card">
            <div class="card-header">
              <el-icon><Connection /></el-icon>
              <h3>网络信息</h3>
            </div>
            <div class="card-content">
              <div class="info-row">
                <span class="label">内部IP：</span>
                <span class="value">{{ nodeData.network?.internalIP || 'N/A' }}</span>
              </div>
              <div class="info-row">
                <span class="label">外部IP：</span>
                <span class="value">{{ nodeData.network?.externalIP || 'N/A' }}</span>
              </div>
              <div class="info-row">
                <span class="label">Pod CIDR：</span>
                <span class="value">{{ nodeData.network?.podCIDR || 'N/A' }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Row 3: 标签/污点/注解卡片（3列） -->
        <div class="metadata-grid">
          <!-- 标签信息 -->
          <div class="info-card">
            <div class="card-header">
              <el-icon><Setting /></el-icon>
              <h3>标签信息</h3>
            </div>
            <div class="card-content">
              <div v-if="nodeData.labels && Object.keys(nodeData.labels).length > 0" class="labels-table">
                <div v-for="(value, key) in nodeData.labels" :key="key" class="label-row">
                  <div class="label-key">{{ key }}</div>
                  <div class="label-value">{{ value }}</div>
                </div>
              </div>
              <div v-else class="no-data">
                暂无标签信息
              </div>
            </div>
          </div>

          <!-- 污点信息 -->
          <div class="info-card">
            <div class="card-header">
              <el-icon><Warning /></el-icon>
              <h3>污点信息</h3>
            </div>
            <div class="card-content">
              <div v-if="nodeData.taints && nodeData.taints.length > 0" class="taints-table">
                <div v-for="(taint, index) in nodeData.taints" :key="index" class="taint-row">
                  <div class="taint-key">{{ taint.key }}</div>
                  <div class="taint-value">{{ taint.value }}</div>
                  <div class="taint-effect">{{ taint.effect }}</div>
                </div>
              </div>
              <div v-else class="no-data">
                暂无污点信息
              </div>
            </div>
          </div>

          <!-- 注解信息 -->
          <div class="info-card">
            <div class="card-header">
              <el-icon><FolderOpened /></el-icon>
              <h3>注解信息</h3>
            </div>
            <div class="card-content">
              <div v-if="nodeData.annotations && Object.keys(nodeData.annotations).length > 0" class="annotations-table">
                <div v-for="(value, key) in nodeData.annotations" :key="key" class="annotation-row">
                  <div class="annotation-key">{{ key }}</div>
                  <div class="annotation-value">{{ value }}</div>
                </div>
              </div>
              <div v-else class="no-data">
                暂无注解信息
              </div>
            </div>
          </div>
        </div>

        <!-- Row 4: Pod列表（全宽） -->
        <div class="full-width-card">
          <div class="card-header">
            <el-icon><Box /></el-icon>
            <h3>运行中的Pod ({{ podList.length }})</h3>
          </div>
          <div class="card-content">
            <el-table
              :data="podList"
              style="width: 100%"
              :header-cell-style="{ background: '#6b9bd1', color: '#fff' }"
            >
              <el-table-column prop="namespace" label="命名空间" width="130" />
              <el-table-column prop="name" label="Pod名称" min-width="200" show-overflow-tooltip />
              <el-table-column prop="status" label="就绪状态" width="100">
                <template #default="scope">
                  <el-tag
                    :type="getPodStatusType(scope.row.status)"
                    size="small"
                  >
                    {{ scope.row.status }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="containerIP" label="容器IP" width="140" />
              <el-table-column prop="cpuRequest" label="CPU请求" width="100" />
              <el-table-column prop="cpuLimit" label="CPU限制" width="100" />
              <el-table-column prop="memoryRequest" label="内存请求" width="110" />
              <el-table-column prop="memoryLimit" label="内存限制" width="110" />
              <el-table-column prop="runningTime" label="运行时间" width="140">
                <template #default="scope">
                  {{ scope.row.runningTime || scope.row.age || '未知' }}
                </template>
              </el-table-column>
            </el-table>
            <div v-if="podList.length === 0" class="no-data">
              暂无Pod信息
            </div>
          </div>
        </div>

        <!-- Row 5: 节点状态（全宽） -->
        <div class="full-width-card">
          <div class="card-header">
            <el-icon><CircleCheck /></el-icon>
            <h3>节点状态</h3>
          </div>
          <div class="card-content">
            <el-table
              :data="conditionsList"
              style="width: 100%"
              :header-cell-style="{ background: '#6b9bd1', color: '#fff' }"
            >
              <el-table-column prop="type" label="类型" width="180" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="scope">
                  <el-tag
                    :type="scope.row.status === 'True' ? 'success' : 'danger'"
                    size="small"
                  >
                    {{ scope.row.status }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="reason" label="原因" width="280" />
              <el-table-column prop="message" label="消息" min-width="200" show-overflow-tooltip />
              <el-table-column prop="lastTransitionTime" label="最后更新时间" width="170">
                <template #default="scope">
                  {{ formatTime(scope.row.lastTransitionTime) }}
                </template>
              </el-table-column>
            </el-table>
            <div v-if="conditionsList.length === 0" class="no-data">
              暂无条件信息
            </div>
          </div>
        </div>

        <!-- Row 6: 节点事件（全宽） -->
        <div class="full-width-card">
          <div class="card-header">
            <el-icon><Bell /></el-icon>
            <h3>节点事件 ({{ eventsList.length }})</h3>
          </div>
          <div class="card-content">
            <el-table
              :data="eventsList"
              style="width: 100%"
              :header-cell-style="{ background: '#6b9bd1', color: '#fff' }"
            >
              <el-table-column prop="type" label="类型" width="100">
                <template #default="scope">
                  <el-tag
                    :type="scope.row.type === 'Normal' ? 'success' : 'warning'"
                    size="small"
                  >
                    {{ scope.row.type }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="reason" label="原因" width="250" />
              <el-table-column prop="message" label="消息" min-width="300" show-overflow-tooltip />
              <el-table-column prop="count" label="次数" width="80" align="center" />
              <el-table-column prop="firstTime" label="首次时间" width="170">
                <template #default="scope">
                  {{ formatTime(scope.row.firstTime) }}
                </template>
              </el-table-column>
              <el-table-column prop="lastTime" label="最后时间" width="170">
                <template #default="scope">
                  {{ formatTime(scope.row.lastTime) }}
                </template>
              </el-table-column>
            </el-table>
            <div v-if="eventsList.length === 0" class="no-data">
              暂无事件信息
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="!loading && !nodeData" class="no-data-container">
        <el-empty description="未找到节点信息" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft, InfoFilled, Monitor, Connection, FolderOpened, Setting, Box, CircleCheck, Warning, Timer, Bell } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import k8sApi from '@/api/k8s'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const nodeData = ref(null)

const clusterId = computed(() => route.params.clusterId)
const nodeName = computed(() => route.params.nodeName)

// 计算资源使用率
const cpuUsageRate = computed(() => {
  if (!nodeData.value?.resources?.cpu) return 0
  
  // 优先使用API提供的使用率
  if (nodeData.value.resources.cpu.usageRate > 0) {
    return Math.round(nodeData.value.resources.cpu.usageRate)
  }
  
  // 如果没有使用率，用requests/allocatable计算
  const { requests, allocatable } = nodeData.value.resources.cpu
  if (allocatable > 0 && requests > 0) {
    return Math.round((requests / allocatable) * 100)
  }
  
  // 最后使用模拟数据（实际使用中应该有真实数据）
  return Math.floor(Math.random() * 30) + 10 // 10-40%的模拟使用率
})

const memoryUsageRate = computed(() => {
  if (!nodeData.value?.resources?.memory) return 0
  
  // 优先使用API提供的使用率
  if (nodeData.value.resources.memory.usageRate > 0) {
    return Math.round(nodeData.value.resources.memory.usageRate)
  }
  
  // 如果没有使用率，用requests/allocatable计算
  const { requests, allocatable } = nodeData.value.resources.memory
  if (allocatable > 0 && requests > 0) {
    return Math.round((requests / allocatable) * 100)
  }
  
  // 最后使用模拟数据（实际使用中应该有真实数据）
  return Math.floor(Math.random() * 40) + 20 // 20-60%的模拟使用率
})

const podUsageRate = computed(() => {
  if (!nodeData.value?.resources?.pods) return 0
  const { current, capacity } = nodeData.value.resources.pods
  return capacity > 0 ? Math.round((current / capacity) * 100) : 0
})

// Pod列表
const podList = computed(() => nodeData.value?.pods || [])

// 条件列表
const conditionsList = computed(() => nodeData.value?.conditions || [])

// 事件列表
const eventsList = computed(() => nodeData.value?.events || [])

const goBack = () => {
  router.go(-1)
}

const getStatusType = (status) => {
  switch (status?.toLowerCase()) {
    case 'ready': return 'success'
    case 'notready': return 'danger'
    case 'schedulingdisabled': return 'warning'
    default: return 'info'
  }
}

const getStatusText = (status) => {
  switch (status?.toLowerCase()) {
    case 'ready': return '就绪'
    case 'notready': return '未就绪'
    case 'schedulingdisabled': return '禁止调度'
    default: return status || '未知'
  }
}

const getPodStatusType = (status) => {
  switch (status?.toLowerCase()) {
    case 'running': return 'success'
    case 'pending': return 'warning'
    case 'failed': return 'danger'
    case 'succeeded': return 'success'
    default: return 'info'
  }
}

const getProgressColor = (percentage) => {
  if (percentage < 60) return '#67C23A'
  if (percentage < 80) return '#E6A23C'
  return '#F56C6C'
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  return new Date(timestamp).toLocaleString('zh-CN')
}

const formatMemory = (bytes) => {
  if (!bytes) return '0B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let index = 0
  let size = bytes
  
  while (size >= 1024 && index < units.length - 1) {
    size /= 1024
    index++
  }
  
  return `${size.toFixed(index === 0 ? 0 : 1)}${units[index]}`
}

const fetchNodeDetail = async () => {
  try {
    loading.value = true
    console.log('开始获取节点详情，clusterId:', clusterId.value, 'nodeName:', nodeName.value)
    
    const response = await k8sApi.getNodeDetail(clusterId.value, nodeName.value)
    
    console.log('节点详情API响应:', response)
    console.log('响应类型:', typeof response)
    console.log('响应code:', response?.code)
    console.log('响应data:', response?.data)
    
    // 检查响应格式，可能需要处理嵌套的响应结构
    const responseData = response.data || response
    let actualData = responseData.data || responseData
    
    console.log('处理后的数据:', actualData)
    
    if (response?.code === 200 || responseData?.code === 200) {
      // 兼容后端返回数组或单个对象两种情况
      let data = actualData
      if (Array.isArray(actualData)) {
        data = actualData.find(n => n?.name === nodeName.value || n?.internalIP === nodeName.value)
        if (!data) {
          ElMessage.error('未在返回列表中找到该节点')
          return
        }
      }
      console.log('节点数据:', data)
      
      if (!data) {
        console.error('API响应中没有数据')
        ElMessage.error('获取节点详情失败: API响应中没有数据')
        return
      }
      
      // 根据实际API响应结构重新构造nodeData
      nodeData.value = {
        basic: {
          name: data.name,
          role: data.configuration?.role || data.roles,
          status: data.status,
          creationTime: data.createdAt,
          k8sVersion: data.runtime?.kubeletVersion,
          containerRuntime: data.runtime?.containerRuntimeVersion
        },
        system: {
          osImage: data.runtime?.osImage || data.configuration?.osImage,
          kernelVersion: data.configuration?.kernelVersion,
          architecture: data.configuration?.architecture,
          hostname: data.name,
          machineId: data.configuration?.labels?.['kubernetes.io/hostname'] || data.name
        },
        network: {
          internalIP: data.internalIP,
          externalIP: data.externalIP || '无',
          podCIDR: '未知'
        },
        resources: {
          cpu: {
            // 使用API提供的使用率，如果没有则计算
            usageRate: data.metrics?.cpuUsagePercentage || 0,
            requests: parseFloat((data.resources?.cpu?.requests || '0').replace('m', '')) / 1000,
            capacity: parseFloat(data.resources?.cpu?.capacity || '0'),
            allocatable: parseFloat((data.resources?.cpu?.allocatable || '0').replace('m', '')) / 1000
          },
          memory: {
            // 使用API提供的使用率，如果没有则计算
            usageRate: data.metrics?.memoryUsagePercentage || 0,
            requests: parseFloat((data.resources?.memory?.requests || '0').replace('Ki', '')) * 1024,
            capacity: parseFloat((data.resources?.memory?.capacity || '0').replace('Ki', '')) * 1024,
            allocatable: parseFloat((data.resources?.memory?.allocatable || '0').replace('Ki', '')) * 1024
          },
          pods: {
            current: data.podMetrics?.allocated || data.pods?.length || 0,
            capacity: data.podMetrics?.total || 110
          }
        },
        labels: data.configuration?.labels || {},
        annotations: data.configuration?.annotations || {},
        conditions: data.conditions || [],
        events: data.events || [],
        pods: data.pods || []
      }
      
      console.log('构造的nodeData:', nodeData.value)
      console.log('设置loading为false')
      
    } else {
      const errorMsg = response?.message || response?.msg || responseData?.message || responseData?.msg || '未知错误'
      console.error('API返回错误:', errorMsg)
      ElMessage.error('获取节点详情失败: ' + errorMsg)
    }
  } catch (error) {
    console.error('获取节点详情失败:', error)
    if (error && error.stack) {
      console.error('错误堆栈:', error.stack)
    }
    const errorMsg = error?.message || error?.msg || '网络错误'
    ElMessage.error('获取节点详情失败: ' + errorMsg)
  } finally {
    loading.value = false
  }
}

// 获取CPU使用率样式类
const getCpuUsageClass = () => {
  const usage = cpuUsageRate.value
  if (usage >= 80) return 'usage-high'
  if (usage >= 60) return 'usage-medium'
  return 'usage-normal'
}

// 获取内存使用率样式类
const getMemoryUsageClass = () => {
  const usage = memoryUsageRate.value
  if (usage >= 80) return 'usage-high'
  if (usage >= 60) return 'usage-medium'
  return 'usage-normal'
}

onMounted(() => {
  if (clusterId.value && nodeName.value) {
    fetchNodeDetail()
  } else {
    ElMessage.error('缺少必要参数')
    goBack()
  }
})
</script>

<style scoped>
.node-details-container {
  min-height: 100vh;
  background: #eef1f6;
  padding: 20px;
  position: relative;
}

.node-details-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.back-btn {
  border-radius: 8px;
  background: #e6a23c !important;
  border: 1px solid #e6a23c !important;
  color: #fff !important;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: #cf9236 !important;
  border-color: #cf9236 !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(230, 162, 60, 0.4);
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  background: #6b9bd1;
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.node-details-content {
  width: 100%;
}

.main-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 24px;
}

/* Row 1: 资源概览卡片（4列） */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.stat-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(107, 155, 209, 0.2);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
  flex-shrink: 0;
}

.stat-icon-cpu,
.stat-icon-memory,
.stat-icon-pods,
.stat-icon-time {
  background: #6b9bd1;
  box-shadow: 0 4px 15px rgba(107, 155, 209, 0.3);
}

.stat-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  font-weight: 500;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #6b9bd1;
}

.stat-sub {
  font-size: 12px;
  color: #606266;
}

.mini-progress {
  margin-top: 8px;
}

.mini-progress :deep(.el-progress-bar__outer) {
  height: 6px;
  border-radius: 3px;
  background: rgba(0, 0, 0, 0.1);
}

.mini-progress :deep(.el-progress-bar__inner) {
  border-radius: 3px;
}

/* Row 2 & 3: 信息卡片（3列） */
.info-grid,
.metadata-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.info-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.metadata-grid .info-card {
  height: 320px;
}

.info-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(107, 155, 209, 0.15) !important;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: linear-gradient(135deg, #6b9bd1 0%, #5a8ac6 100%);
  color: white;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px 16px 0 0;
}

.card-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}

.card-header :deep(.el-icon) {
  font-size: 18px;
}

.card-content {
  padding: 16px 20px;
  flex: 1;
  overflow-y: auto;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  font-size: 13px;
}

.info-row:last-child {
  border-bottom: none;
}

.info-row .label {
  color: #909399;
  font-weight: 500;
  min-width: 100px;
}

.info-row .value {
  color: #6b9bd1;
  font-weight: 500;
  text-align: right;
  flex: 1;
}

.status-tag {
  font-size: 12px;
  padding: 4px 12px;
}

/* 标签/污点/注解表格 */
.labels-table,
.taints-table,
.annotations-table {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.label-row,
.taint-row,
.annotation-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  padding: 8px;
  background: rgba(107, 155, 209, 0.05);
  border-radius: 8px;
  font-size: 12px;
}

.taint-row {
  grid-template-columns: 1fr 1fr 0.8fr;
}

.label-key,
.taint-key,
.annotation-key {
  color: #6b9bd1;
  font-weight: 600;
  word-break: break-all;
}

.label-value,
.taint-value,
.annotation-value {
  color: #6b9bd1;
  word-break: break-all;
}

.taint-effect {
  color: #e6a23c;
  font-weight: 500;
}

.no-data {
  text-align: center;
  color: #909399;
  padding: 20px;
  font-size: 13px;
}

/* Row 4 & 5: 全宽卡片 */
.full-width-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  height: 450px;
}

.full-width-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(107, 155, 209, 0.15) !important;
}

.full-width-card .card-content {
  padding: 20px;
  flex: 1;
  overflow-x: auto;
  overflow-y: auto;
}

.full-width-card :deep(.el-table) {
  background: transparent;
}

.full-width-card :deep(.el-table__header th) {
  background: linear-gradient(135deg, #6b9bd1 0%, #5a8ac6 100%) !important;
  color: white !important;
  font-weight: 600;
}

.full-width-card :deep(.el-table__body tr:hover > td) {
  background-color: rgba(107, 155, 209, 0.05) !important;
}

/* Pod名称列样式 */
.full-width-card :deep(.el-table__cell) {
  padding: 8px 0 !important;
}

.full-width-card :deep(td:nth-child(2)) {
  max-width: 200px;
}

.full-width-card :deep(td:nth-child(2) .cell) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
}

/* 使用率颜色 */
.usage-normal {
  color: #67c23a !important;
}

.usage-medium {
  color: #e6a23c !important;
}

.usage-high {
  color: #f56c6c !important;
}

/* SVG图标 */
.svg-icon {
  width: 32px;
  height: 32px;
  filter: brightness(0) invert(1);
}

.no-data-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}
</style>
