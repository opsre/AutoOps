<template>
  <div class="k8s-cluster-details">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <el-button 
          :icon="ArrowLeft" 
          size="small" 
          @click="goBack"
          class="back-btn"
        >
          返回
        </el-button>
        <div class="cluster-info">
          <h1 class="cluster-title">
            <img src="@/assets/image/k8s.svg" alt="k8s" class="k8s-icon" />
            {{ clusterDetails.clusterName || '集群详情' }}
          </h1>
          <div class="cluster-meta">
            <el-tag 
              :type="getStatusTag(clusterDetails.status)" 
              size="small" 
              effect="dark"
            >
              {{ getStatusText(clusterDetails.status) }}
            </el-tag>
            <span class="cluster-version">{{ clusterDetails.version }}</span>
            <span class="cluster-type">{{ getClusterTypeText(clusterDetails.clusterType) }}</span>
          </div>
        </div>
      </div>
      <div class="header-actions">
        <el-button 
          :icon="Refresh" 
          size="small" 
          @click="refreshData"
          :loading="loading"
        >
          刷新
        </el-button>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="details-content" v-loading="loading">
      <!-- 集群概览卡片 -->
      <div class="overview-cards">
        <el-row :gutter="16">
          <el-col :span="6">
            <el-card class="overview-card nodes-card">
              <div class="card-content">
                <div class="card-icon">
                  <el-icon size="24"><Monitor /></el-icon>
                </div>
                <div class="card-info">
                  <div class="card-value">{{ clusterDetails.nodeStats?.totalNodes || 0 }}</div>
                  <div class="card-label">节点总数</div>
                  <div class="card-sub">就绪: {{ clusterDetails.nodeStats?.readyNodes || 0 }}</div>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="overview-card workloads-card">
              <div class="card-content">
                <div class="card-icon">
                  <el-icon size="24"><Box /></el-icon>
                </div>
                <div class="card-info">
                  <div class="card-value">{{ getTotalWorkloads() }}</div>
                  <div class="card-label">工作负载</div>
                  <div class="card-sub">Pod: {{ clusterDetails.workloadStats?.totalPods || 0 }}</div>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="overview-card cpu-card">
              <div class="card-content">
                <div class="card-icon">
                  <img src="@/assets/image/cpu.svg" alt="CPU" class="svg-icon" />
                </div>
                <div class="card-info">
                  <div class="card-value" :class="getCpuUsageClass()">{{ clusterDetails.monitoring?.cpu?.usagePercent || 0 }}%</div>
                  <div class="card-label">CPU使用率</div>
                  <div class="card-sub">{{ clusterDetails.monitoring?.cpu?.availableCores || 0 }}核可用 / {{ clusterDetails.monitoring?.cpu?.totalCores || 0 }}核总量</div>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="overview-card memory-card">
              <div class="card-content">
                <div class="card-icon">
                  <img src="@/assets/image/内存.svg" alt="内存" class="svg-icon" />
                </div>
                <div class="card-info">
                  <div class="card-value" :class="getMemoryUsageClass()">{{ clusterDetails.monitoring?.memory?.usagePercent || 0 }}%</div>
                  <div class="card-label">内存使用率</div>
                  <div class="card-sub">{{ clusterDetails.monitoring?.memory?.availableMi || 0 }}Mi可用 / {{ clusterDetails.monitoring?.memory?.totalMi || 0 }}Mi总量</div>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 详细信息区域 -->
      <div class="details-sections">
        <el-row :gutter="16">
          <!-- 左侧列 -->
          <el-col :span="12">
            <!-- 基础信息 -->
            <el-card class="info-card" shadow="hover">
              <template #header>
                <div class="card-header">
                  <el-icon><InfoFilled /></el-icon>
                  <span>基础信息</span>
                </div>
              </template>
              <div class="info-content">
                <div class="info-row">
                  <span class="label">集群名称:</span>
                  <span class="value">{{ clusterDetails.clusterName }}</span>
                </div>
                <div class="info-row">
                  <span class="label">Kubernetes版本:</span>
                  <span class="value">{{ clusterDetails.version }}</span>
                </div>
                <div class="info-row">
                  <span class="label">集群类型:</span>
                  <span class="value">{{ getClusterTypeText(clusterDetails.clusterType) }}</span>
                </div>
                <div class="info-row">
                  <span class="label">创建时间:</span>
                  <span class="value">{{ formatTime(clusterDetails.createdAt) }}</span>
                </div>
                <div class="info-row">
                  <span class="label">最后同步:</span>
                  <span class="value">{{ formatTime(clusterDetails.lastSyncTime) }}</span>
                </div>
                <div class="info-row">
                  <span class="label">运行时间:</span>
                  <span class="value">{{ clusterDetails.runtime?.uptime || '未知' }}</span>
                </div>
              </div>
            </el-card>

            <!-- 网络配置 -->
            <el-card class="info-card" shadow="hover">
              <template #header>
                <div class="card-header">
                  <el-icon><Connection /></el-icon>
                  <span>网络配置</span>
                </div>
              </template>
              <div class="info-content">
                <div class="info-row">
                  <span class="label">Service CIDR:</span>
                  <span class="value">{{ clusterDetails.networkConfig?.serviceCIDR || '未配置' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">Pod CIDR:</span>
                  <span class="value">{{ clusterDetails.networkConfig?.podCIDR || '未配置' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">API Server:</span>
                  <span class="value">{{ clusterDetails.networkConfig?.apiServerEndpoint || '未配置' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">网络插件:</span>
                  <span class="value">{{ clusterDetails.networkConfig?.networkPlugin || '未知' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">服务转发:</span>
                  <span class="value">{{ clusterDetails.networkConfig?.proxyMode || '未知' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">DNS服务:</span>
                  <span class="value">{{ clusterDetails.networkConfig?.dnsService || '未知' }}</span>
                </div>
              </div>
            </el-card>

            <!-- 存储信息 -->
            <el-card class="info-card" shadow="hover">
              <template #header>
                <div class="card-header">
                  <el-icon><FolderOpened /></el-icon>
                  <span>存储信息</span>
                </div>
              </template>
              <div class="info-content">
                <div class="info-row">
                  <span class="label">PV总数:</span>
                  <span class="value">{{ clusterDetails.monitoring?.storage?.totalPVs || 0 }}</span>
                </div>
                <div class="info-row">
                  <span class="label">PVC总数:</span>
                  <span class="value">{{ clusterDetails.monitoring?.storage?.totalPVCs || 0 }}</span>
                </div>
                <div class="info-row">
                  <span class="label">存储类:</span>
                  <span class="value">{{ clusterDetails.monitoring?.storage?.storageClasses || 0 }}</span>
                </div>
              </div>
            </el-card>
          </el-col>

          <!-- 右侧列 -->
          <el-col :span="12">
            <!-- 工作负载统计 -->
            <el-card class="info-card" shadow="hover">
              <template #header>
                <div class="card-header">
                  <el-icon><Box /></el-icon>
                  <span>工作负载统计</span>
                </div>
              </template>
              <div class="workload-stats">
                <div class="stat-item">
                  <div class="stat-icon deployments">
                    <el-icon><Grid /></el-icon>
                  </div>
                  <div class="stat-info">
                    <div class="stat-value">{{ clusterDetails.workloadStats?.deployments || 0 }}</div>
                    <div class="stat-label">Deployments</div>
                  </div>
                </div>
                <div class="stat-item">
                  <div class="stat-icon statefulsets">
                    <el-icon><Files /></el-icon>
                  </div>
                  <div class="stat-info">
                    <div class="stat-value">{{ clusterDetails.workloadStats?.statefulSets || 0 }}</div>
                    <div class="stat-label">StatefulSets</div>
                  </div>
                </div>
                <div class="stat-item">
                  <div class="stat-icon daemonsets">
                    <el-icon><Service /></el-icon>
                  </div>
                  <div class="stat-info">
                    <div class="stat-value">{{ clusterDetails.workloadStats?.daemonSets || 0 }}</div>
                    <div class="stat-label">DaemonSets</div>
                  </div>
                </div>
                <div class="stat-item">
                  <div class="stat-icon jobs">
                    <el-icon><Timer /></el-icon>
                  </div>
                  <div class="stat-info">
                    <div class="stat-value">{{ clusterDetails.workloadStats?.jobs || 0 }}</div>
                    <div class="stat-label">Jobs</div>
                  </div>
                </div>
                <div class="stat-item">
                  <div class="stat-icon pods">
                    <el-icon><Box /></el-icon>
                  </div>
                  <div class="stat-info">
                    <div class="stat-value">{{ clusterDetails.workloadStats?.totalPods || 0 }}</div>
                    <div class="stat-label">总Pod数</div>
                  </div>
                </div>
                <div class="stat-item">
                  <div class="stat-icon running-pods">
                    <el-icon><CircleCheck /></el-icon>
                  </div>
                  <div class="stat-info">
                    <div class="stat-value">{{ clusterDetails.workloadStats?.runningPods || 0 }}</div>
                    <div class="stat-label">运行中Pod</div>
                  </div>
                </div>
              </div>
            </el-card>

            <!-- 安装的组件 -->
            <el-card class="info-card" shadow="hover">
              <template #header>
                <div class="card-header">
                  <el-icon><Setting /></el-icon>
                  <span>安装的组件</span>
                </div>
              </template>
              <div class="components-grid">
                <div 
                  v-for="component in clusterDetails.installedComponents || []" 
                  :key="component.name"
                  class="component-item"
                >
                  <div class="component-status" :class="component.status">
                    <el-icon v-if="component.status === 'Running'"><CircleCheck /></el-icon>
                    <el-icon v-else><Warning /></el-icon>
                  </div>
                  <div class="component-info">
                    <div class="component-name">{{ component.name }}</div>
                    <div class="component-version">{{ component.version || '未知版本' }}</div>
                  </div>
                </div>
              </div>
            </el-card>

            <!-- 运行时信息 -->
            <el-card class="info-card" shadow="hover">
              <template #header>
                <div class="card-header">
                  <el-icon><Cpu /></el-icon>
                  <span>运行时信息</span>
                </div>
              </template>
              <div class="info-content">
                <div class="info-row">
                  <span class="label">容器运行时:</span>
                  <span class="value">{{ clusterDetails.runtime?.containerRuntimeVersion || '未知' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">API Server版本:</span>
                  <span class="value">{{ clusterDetails.runtime?.apiServerVersion || '未知' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">ETCD版本:</span>
                  <span class="value">{{ clusterDetails.runtime?.etcdVersion || '未知' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">CoreDNS版本:</span>
                  <span class="value">{{ clusterDetails.runtime?.coreDNSVersion || '未知' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">Kube-proxy版本:</span>
                  <span class="value">{{ clusterDetails.runtime?.kubeProxyVersion || '未知' }}</span>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 节点信息 -->
        <el-card class="nodes-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Monitor /></el-icon>
              <span>节点信息</span>
              <el-tag size="small" type="info">{{ clusterDetails.nodeList?.length || 0 }}个节点</el-tag>
            </div>
          </template>
          <el-table 
            :data="clusterDetails.nodeList || []" 
            size="small"
            class="nodes-table"
            max-height="300"
          >
            <el-table-column prop="name" label="节点名称" width="200" />
            <el-table-column prop="role" label="角色" width="100">
              <template #default="{ row }">
                <el-tag 
                  :type="row.role === 'master' ? 'danger' : 'primary'" 
                  size="small"
                >
                  {{ row.role === 'master' ? 'Master' : 'Worker' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag 
                  :type="row.status === 'Ready' ? 'success' : 'warning'" 
                  size="small"
                >
                  {{ row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="internalIP" label="内部IP" width="150">
              <template #default="{ row }">
                <el-button 
                  type="text" 
                  size="small" 
                  @click="goToNodeDetail(row.name, row.internalIP)"
                  class="ip-link"
                >
                  {{ row.internalIP }}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column prop="externalIP" label="外部IP" width="150" />
            <el-table-column prop="version" label="K8s版本" width="120" />
            <el-table-column prop="osImage" label="操作系统" width="200" show-overflow-tooltip />
            <el-table-column prop="createdAt" label="创建时间" width="160">
              <template #default="{ row }">
                {{ formatTime(row.createdAt) }}
              </template>
            </el-table-column>
          </el-table>
          <div v-if="!clusterDetails.nodeList || clusterDetails.nodeList.length === 0" class="no-data">
            暂无节点信息
          </div>
        </el-card>

        <!-- 集群事件 -->
        <el-card class="events-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><DocumentCopy /></el-icon>
              <span>最近事件</span>
              <el-tag size="small" type="info">最近50条</el-tag>
            </div>
          </template>
          <el-table 
            :data="clusterDetails.events || []" 
            size="small"
            class="events-table"
            max-height="300"
          >
            <el-table-column prop="type" label="类型" width="80">
              <template #default="{ row }">
                <el-tag 
                  :type="row.type === 'Warning' ? 'warning' : 'success'" 
                  size="small"
                >
                  {{ row.type }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="reason" label="原因" width="120" />
            <el-table-column prop="message" label="消息" min-width="200" show-overflow-tooltip />
            <el-table-column prop="source" label="来源" width="150" show-overflow-tooltip />
            <el-table-column prop="count" label="次数" width="60" />
            <el-table-column prop="lastTimestamp" label="最后发生时间" width="160">
              <template #default="{ row }">
                {{ formatTime(row.lastTimestamp) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  Refresh,
  Monitor,
  Box,
  Cpu,
  Connection,
  Histogram,
  InfoFilled,
  FolderOpened,
  Setting,
  Grid,
  Files,
  Service,
  Timer,
  CircleCheck,
  Warning,
  DocumentCopy
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const clusterId = ref(route.params.clusterId)

// 集群详情数据
const clusterDetails = reactive({
  clusterName: '',
  version: '',
  status: '',
  clusterType: '',
  createdAt: '',
  lastSyncTime: '',
  nodeStats: {
    totalNodes: 0,
    readyNodes: 0,
    masterNodes: 0,
    workerNodes: 0
  },
  workloadStats: {
    deployments: 0,
    statefulSets: 0,
    daemonSets: 0,
    jobs: 0,
    cronJobs: 0,
    totalPods: 0,
    runningPods: 0
  },
  installedComponents: [],
  networkConfig: {
    serviceCIDR: '',
    podCIDR: '',
    apiServerEndpoint: '',
    networkPlugin: '',
    proxyMode: '',
    dnsService: ''
  },
  monitoring: {
    cpu: {
      totalCores: 0,
      availableCores: 0,
      usagePercent: 0
    },
    memory: {
      totalMi: 0,
      availableMi: 0,
      usagePercent: 0
    },
    network: {
      inboundTraffic: '',
      outboundTraffic: ''
    },
    storage: {
      totalPVs: 0,
      totalPVCs: 0,
      storageClasses: 0
    }
  },
  runtime: {
    kubernetesVersion: '',
    containerRuntimeVersion: '',
    apiServerVersion: '',
    etcdVersion: '',
    coreDNSVersion: '',
    kubeProxyVersion: '',
    uptime: ''
  },
  events: [],
  nodeList: []
})

// 获取集群详情
const fetchClusterDetails = async () => {
  try {
    loading.value = true
    console.log('获取集群详情，ID:', clusterId.value)
    
    const response = await k8sApi.getClusterDetail(clusterId.value)
    const responseData = response.data || response
    
    console.log('集群详情API响应:', responseData)
    
    if (responseData.code === 200 || responseData.success) {
      const data = responseData.data
      const cluster = data.cluster || {}
      
      console.log('解析集群数据:', cluster)
      
      // 填充基础信息
      clusterDetails.clusterName = cluster.name || ''
      clusterDetails.version = cluster.version || ''
      clusterDetails.status = cluster.status || ''
      clusterDetails.clusterType = cluster.clusterType || ''
      clusterDetails.createdAt = cluster.createdAt || ''
      clusterDetails.lastSyncTime = cluster.lastSyncAt || ''
      
      // 填充节点统计 - 使用summary数据
      if (data.summary) {
        clusterDetails.nodeStats.totalNodes = data.summary.totalNodes || 0
        clusterDetails.nodeStats.readyNodes = data.summary.readyNodes || 0
        clusterDetails.nodeStats.masterNodes = data.summary.masterNodes || 0
        clusterDetails.nodeStats.workerNodes = data.summary.workerNodes || 0
      }
      
      // 填充工作负载统计 - 使用workloads数据
      if (data.workloads) {
        clusterDetails.workloadStats.deployments = data.workloads.totalDeployments || 0
        clusterDetails.workloadStats.statefulSets = data.workloads.totalStatefulSets || 0
        clusterDetails.workloadStats.daemonSets = data.workloads.totalDaemonSets || 0
        clusterDetails.workloadStats.jobs = data.workloads.totalJobs || 0
        clusterDetails.workloadStats.cronJobs = data.workloads.totalCronJobs || 0
        clusterDetails.workloadStats.totalPods = data.workloads.totalPods || 0
        clusterDetails.workloadStats.runningPods = data.workloads.runningPods || 0
      }
      
      // 填充安装的组件 - 处理可能的版本字段映射
      if (data.components) {
        clusterDetails.installedComponents = data.components.map(component => {
          let version = component.version || 
                       component.componentVersion || 
                       component.currentVersion ||
                       component.serverVersion ||
                       ''
          
          // 处理 "unknown" 字符串
          if (!version || version === 'unknown' || version === 'Unknown') {
            version = '未知版本'
          }
          
          return {
            name: component.name || component.componentName || '',
            version: version,
            status: component.status || component.state || 'Unknown'
          }
        })
      } else {
        clusterDetails.installedComponents = []
      }
      
      // 填充网络配置
      if (data.network) {
        clusterDetails.networkConfig.serviceCIDR = data.network.serviceCIDR || ''
        clusterDetails.networkConfig.podCIDR = data.network.podCIDR || ''
        clusterDetails.networkConfig.apiServerEndpoint = data.network.apiServerEndpoint || ''
        clusterDetails.networkConfig.networkPlugin = data.network.networkPlugin || ''
        clusterDetails.networkConfig.proxyMode = data.network.proxyMode || ''
        clusterDetails.networkConfig.dnsService = data.network.dnsService || ''
      }
      
      // 填充监控信息
      if (data.monitoring) {
        // CPU信息
        if (data.monitoring.cpu) {
          const totalCores = parseFloat(data.monitoring.cpu.total?.replace(' cores', '')) || 0
          const availableCores = parseFloat(data.monitoring.cpu.available?.replace(' cores', '')) || 0
          const usedCores = parseFloat(data.monitoring.cpu.used?.replace(' cores', '')) || 0
          
          clusterDetails.monitoring.cpu.totalCores = totalCores
          clusterDetails.monitoring.cpu.availableCores = availableCores
          
          // 计算CPU使用率：如果API提供了usageRate使用API值，否则计算 (total - available) / total * 100
          let cpuUsage = data.monitoring.cpu.usageRate || 0
          if (cpuUsage === 0 && totalCores > 0 && availableCores > 0) {
            cpuUsage = Math.round(((totalCores - availableCores) / totalCores) * 100)
          }
          clusterDetails.monitoring.cpu.usagePercent = cpuUsage
        }
        
        // 内存信息
        if (data.monitoring.memory) {
          const totalMi = parseFloat(data.monitoring.memory.total?.replace(' Mi', '')) || 0
          const availableMi = parseFloat(data.monitoring.memory.available?.replace(' Mi', '')) || 0
          const usedMi = parseFloat(data.monitoring.memory.used?.replace(' Mi', '')) || 0
          
          clusterDetails.monitoring.memory.totalMi = totalMi
          clusterDetails.monitoring.memory.availableMi = availableMi
          
          // 计算内存使用率：如果API提供了usageRate使用API值，否则计算 (total - available) / total * 100
          let memoryUsage = data.monitoring.memory.usageRate || 0
          if (memoryUsage === 0 && totalMi > 0 && availableMi > 0) {
            memoryUsage = Math.round(((totalMi - availableMi) / totalMi) * 100)
          }
          clusterDetails.monitoring.memory.usagePercent = memoryUsage
        }
        
        // 网络信息
        if (data.monitoring.network) {
          clusterDetails.monitoring.network.inboundTraffic = data.monitoring.network.inboundTraffic || ''
          clusterDetails.monitoring.network.outboundTraffic = data.monitoring.network.outboundTraffic || ''
        }
        
        // 存储信息
        if (data.monitoring.storage) {
          clusterDetails.monitoring.storage.totalPVs = data.monitoring.storage.totalPVs || 0
          clusterDetails.monitoring.storage.totalPVCs = data.monitoring.storage.totalPVCs || 0
          clusterDetails.monitoring.storage.storageClasses = data.monitoring.storage.storageClasses?.length || 0
        }
      }
      
      // 填充运行时信息
      if (data.runtime) {
        clusterDetails.runtime.kubernetesVersion = data.runtime.kubernetesVersion || data.runtime.k8sVersion || ''
        clusterDetails.runtime.containerRuntimeVersion = data.runtime.containerRuntime || data.runtime.containerRuntimeVersion || ''
        clusterDetails.runtime.apiServerVersion = data.runtime.apiServerVersion || data.runtime.apiServer || ''
        // 尝试多种可能的etcd版本字段名，处理空字符串
        const etcdVersion = data.runtime.etcdVersion || 
                          data.runtime.etcd?.version || 
                          data.runtime.etcd?.serverVersion ||
                          data.runtime.etcdServerVersion ||
                          ''
        clusterDetails.runtime.etcdVersion = etcdVersion || '未知'
        clusterDetails.runtime.coreDNSVersion = data.runtime.coreDNSVersion || data.runtime.coreDns?.version || ''
        clusterDetails.runtime.kubeProxyVersion = data.runtime.kubeProxyVersion || data.runtime.kubeProxy?.version || ''
        clusterDetails.runtime.uptime = data.runtime.upTime || data.runtime.uptime || ''
      }
      
      // 填充事件信息
      clusterDetails.events = data.events || []
      
      // 填充节点列表 - 映射正确的字段名
      if (data.nodes) {
        clusterDetails.nodeList = data.nodes.map(node => ({
          name: node.name || '',
          role: node.role || 'unknown',
          status: node.status || 'Unknown',
          internalIP: node.internalIP || '',
          externalIP: node.externalIP || '无',  // 处理空的外部IP
          version: node.version || '',
          osImage: node.os || node.osImage || '未知',  // 使用 os 字段
          createdAt: node.createdAt || node.creationTimestamp || '未知',  // 尝试多个时间字段
          // 保持其他可能需要的原始字段
          capacity: node.capacity,
          allocatable: node.allocatable,
          conditions: node.conditions
        }))
      } else {
        clusterDetails.nodeList = []
      }
      
    } else {
      ElMessage.error(responseData.message || '获取集群详情失败')
    }
  } catch (error) {
    console.error('获取集群详情失败:', error)
    
    if (error.response?.status === 404) {
      ElMessage.error('集群不存在或已被删除')
      router.push('/k8s/list')
    } else if (error.response?.status === 401) {
      ElMessage.error('权限不足，无法访问集群详情')
    } else if (error.code === 'ERR_NETWORK') {
      ElMessage.error('网络连接失败，请检查网络状态')
    } else {
      ElMessage.error('获取集群详情失败，请重试')
    }
  } finally {
    loading.value = false
  }
}

// 返回上一页
const goBack = () => {
  router.push('/k8s/list')
}

// 刷新数据
const refreshData = () => {
  fetchClusterDetails()
}

// 获取状态标签类型
const getStatusTag = (status) => {
  const tagMap = {
    1: 'info',      // 创建中
    2: 'success',   // 运行中
    3: 'warning',   // 离线
    'creating': 'info',
    'online': 'success', 
    'offline': 'warning'
  }
  return tagMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const textMap = {
    1: '创建中',
    2: '运行中', 
    3: '离线',
    'creating': '创建中',
    'online': '运行中',
    'offline': '离线'
  }
  return textMap[status] || '未知'
}

// 获取集群类型文本
const getClusterTypeText = (type) => {
  const typeMap = {
    1: '自建集群',
    2: '导入集群'
  }
  return typeMap[type] || '未知类型'
}

// 获取工作负载总数
const getTotalWorkloads = () => {
  const stats = clusterDetails.workloadStats
  return (stats.deployments || 0) + 
         (stats.statefulSets || 0) + 
         (stats.daemonSets || 0) + 
         (stats.jobs || 0) + 
         (stats.cronJobs || 0)
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '未知'
  return new Date(time).toLocaleString('zh-CN')
}

// 获取CPU使用率样式类
const getCpuUsageClass = () => {
  const usage = clusterDetails.monitoring?.cpu?.usagePercent || 0
  if (usage >= 80) return 'usage-high'
  if (usage >= 60) return 'usage-medium'
  return 'usage-normal'
}

// 获取内存使用率样式类
const getMemoryUsageClass = () => {
  const usage = clusterDetails.monitoring?.memory?.usagePercent || 0
  if (usage >= 80) return 'usage-high'
  if (usage >= 60) return 'usage-medium'
  return 'usage-normal'
}

// 跳转到节点详情
const goToNodeDetail = (nodeName, nodeIP) => {
  console.log('跳转到节点详情:', nodeName, nodeIP)
  router.push(`/k8s/cluster/${clusterId.value}/node/${nodeName}`)
}

onMounted(() => {
  fetchClusterDetails()
})
</script>

<style scoped>
.k8s-cluster-details {
  padding: 20px;
  min-height: 100vh;
  background: #eef1f6;
  position: relative;
}

.k8s-cluster-details::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.05) 100%);
  pointer-events: none;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  position: relative;
  z-index: 1;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-btn {
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(102, 126, 234, 0.3);
  color: #667eea;
  backdrop-filter: blur(5px);
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: rgba(102, 126, 234, 0.1);
  border-color: #667eea;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.cluster-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cluster-title {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #6b9bd1;
}

.k8s-icon {
  width: 32px;
  height: 32px;
}

.cluster-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.cluster-version {
  color: #606266;
  font-size: 14px;
  background: #f0f2f5;
  padding: 4px 8px;
  border-radius: 4px;
}

.cluster-type {
  color: #909399;
  font-size: 12px;
}

/* 概览卡片 */
.overview-cards {
  margin-bottom: 24px;
  position: relative;
  z-index: 1;
}

.overview-card {
  border-radius: 12px;
  border: none;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.overview-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(102, 126, 234, 0.3);
  background: rgba(255, 255, 255, 0.95);
}

.card-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.nodes-card .card-icon {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.workloads-card .card-icon {
  background: linear-gradient(135deg, #f093fb, #f5576c);
}

.cpu-card .card-icon {
  background: linear-gradient(135deg, #4facfe, #00f2fe);
}

.memory-card .card-icon {
  background: linear-gradient(135deg, #43e97b, #38f9d7);
}

/* SVG图标样式 */
.svg-icon {
  width: 24px;
  height: 24px;
  filter: brightness(0) invert(1);
  transition: all 0.3s ease;
}

.overview-card:hover .svg-icon {
  transform: scale(1.1) rotate(5deg);
  filter: brightness(0) invert(1) drop-shadow(0 0 8px rgba(255, 255, 255, 0.5));
}

.card-info {
  flex: 1;
}

.card-value {
  font-size: 28px;
  font-weight: 700;
  color: #6b9bd1;
  line-height: 1;
}

.card-label {
  font-size: 14px;
  color: #606266;
  margin: 4px 0;
}

.card-sub {
  font-size: 12px;
  color: #909399;
}

/* 使用率颜色指示 */
.usage-normal {
  color: #67c23a !important;
}

.usage-medium {
  color: #e6a23c !important;
}

.usage-high {
  color: #f56c6c !important;
}

/* 详细信息区域 */
.details-sections {
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: relative;
  z-index: 1;
}

.info-card {
  border-radius: 12px;
  border: none;
  margin-bottom: 16px;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.info-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.2);
  background: rgba(255, 255, 255, 0.95);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #6b9bd1;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f2f5;
}

.info-row:last-child {
  border-bottom: none;
}

.info-row .label {
  color: #606266;
  font-weight: 500;
  min-width: 120px;
}

.info-row .value {
  color: #6b9bd1;
  text-align: right;
  word-break: break-all;
}

/* 工作负载统计 */
.workload-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8f9fc;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.stat-item:hover {
  background: #e6f3ff;
}

.stat-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
}

.stat-icon.deployments {
  background: #409eff;
}

.stat-icon.statefulsets {
  background: #67c23a;
}

.stat-icon.daemonsets {
  background: #e6a23c;
}

.stat-icon.jobs {
  background: #f56c6c;
}

.stat-icon.pods {
  background: #909399;
}

.stat-icon.running-pods {
  background: #67c23a;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #6b9bd1;
  line-height: 1;
}

.stat-label {
  font-size: 12px;
  color: #606266;
  margin-top: 2px;
}

/* 组件网格 */
.components-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.component-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8f9fc;
  border-radius: 8px;
  border-left: 4px solid #e4e7ed;
}

.component-item .component-status.Running {
  color: #67c23a;
}

.component-item .component-status:not(.Running) {
  color: #f56c6c;
}

.component-status {
  font-size: 16px;
}

.component-name {
  font-weight: 500;
  color: #6b9bd1;
  font-size: 14px;
}

.component-version {
  font-size: 12px;
  color: #909399;
}

/* 节点表格 */
.nodes-card {
  margin-bottom: 16px;
}

.nodes-table {
  border-radius: 8px;
  overflow: hidden;
}

.nodes-table :deep(.el-table__header) {
  background: #f8f9fc;
}

.nodes-table :deep(.el-table__row:hover) {
  background: #f0f9ff;
}

.ip-link {
  color: #667eea !important;
  font-weight: 500;
  transition: all 0.3s ease;
}

.ip-link:hover {
  color: #409eff !important;
  text-decoration: underline;
}

.no-data {
  text-align: center;
  color: #909399;
  padding: 40px;
  font-size: 14px;
}

/* 事件表格 */
.events-card {
  margin-top: 16px;
}

.events-table {
  border-radius: 8px;
  overflow: hidden;
}

.events-table :deep(.el-table__header) {
  background: #f8f9fc;
}

.events-table :deep(.el-table__row:hover) {
  background: #f0f9ff;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .overview-cards .el-col {
    margin-bottom: 16px;
  }
  
  .workload-stats {
    grid-template-columns: 1fr;
  }
  
  .components-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .k8s-cluster-details {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .cluster-title {
    font-size: 20px;
  }
  
  .cluster-meta {
    flex-wrap: wrap;
  }
  
  .details-sections .el-row {
    flex-direction: column;
  }
  
  .info-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .info-row .value {
    text-align: left;
  }
}
</style>
