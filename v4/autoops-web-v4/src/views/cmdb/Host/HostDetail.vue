<template>
  <div class="host-detail-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <el-button :icon="ArrowLeft" @click="goBack" class="back-btn">返回</el-button>
        <h2 class="page-title">{{ hostInfo.hostName || '主机详情' }}</h2>
        <el-tag :type="getStatusTagType(hostInfo.status)" size="large">
          {{ getStatusText(hostInfo.status) }}
        </el-tag>
      </div>
    </div>

    <!-- Tab 切换 -->
    <div class="tab-container">
      <el-tabs v-model="activeTab" class="detail-tabs">
        <el-tab-pane label="主机详情" name="details"></el-tab-pane>
        <el-tab-pane label="监控图表" name="monitor"></el-tab-pane>
        <el-tab-pane label="告警记录" name="alert"></el-tab-pane>
      </el-tabs>
    </div>

    <el-skeleton :loading="loading" animated>
      <template #template>
        <el-skeleton-item variant="h1" style="width: 50%" />
        <el-skeleton-item variant="text" style="margin-top: 20px" />
        <el-skeleton-item variant="text" />
        <el-skeleton-item variant="text" />
      </template>

      <template #default>
        <!-- 主机详情页面 -->
        <div v-show="activeTab === 'details'" class="details-content">
          <!-- 左侧: 主机基本信息 -->
          <div class="left-section">
            <el-card shadow="hover" class="info-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><Monitor /></el-icon>
                  <span class="card-title">主机信息</span>
                </div>
              </template>
              <div class="info-content">
                <div class="info-row">
                  <span class="label">实例名称:</span>
                  <span class="value">{{ hostInfo.name || '-' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">业务组:</span>
                  <span class="value">{{ hostInfo.groupName || '-' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">内网IP:</span>
                  <span class="value">{{ hostInfo.privateIp || '-' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">公网IP:</span>
                  <span class="value">{{ hostInfo.publicIp || '无' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">配置:</span>
                  <span class="value">{{ hostInfo.cpu }}核-{{ hostInfo.memory }}GB-{{ hostInfo.disk }}GB磁盘</span>
                </div>
                <div class="info-row">
                  <span class="label">操作系统:</span>
                  <span class="value">{{ hostInfo.os || '-' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">系统运行天数:</span>
                  <span class="value">{{ monitorData.uptimeStr || '加载中...' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">系统进程总数:</span>
                  <span class="value">{{ monitorData.processCount || 0 }}</span>
                </div>
                <div class="info-row">
                  <span class="label">文件打开数量:</span>
                  <span class="value">{{ monitorData.openFiles || 0 }}</span>
                </div>
              </div>
            </el-card>
          </div>

          <!-- 右侧: 4个监控卡片 (2x2布局) -->
          <div class="right-section">
            <!-- 卡片2: CPU详情 -->
            <el-card shadow="hover" class="metric-card cpu-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><Cpu /></el-icon>
                  <span class="card-title">CPU详情</span>
                  <el-button
                    size="small"
                    :icon="Refresh"
                    :loading="refreshing"
                    @click="refreshMonitorData"
                    circle
                  />
                </div>
              </template>
              <div class="metric-content">
                <div class="metric-main">
                  <div class="metric-value-large">{{ monitorData.cpu?.total?.toFixed(1) || '0.0' }}%</div>
                  <div class="metric-label">总CPU使用率</div>
                </div>
                <el-divider />
                <div class="metric-details">
                  <div class="detail-item">
                    <span class="detail-label">用户态使用率</span>
                    <span class="detail-value">{{ monitorData.cpu?.user?.toFixed(1) || '0.0' }}%</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">内核态使用率</span>
                    <span class="detail-value">{{ monitorData.cpu?.system?.toFixed(1) || '0.0' }}%</span>
                  </div>
                </div>
              </div>
            </el-card>

            <!-- 卡片3: 内存详情 -->
            <el-card shadow="hover" class="metric-card memory-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><Memo /></el-icon>
                  <span class="card-title">内存详情</span>
                </div>
              </template>
              <div class="metric-content">
                <div class="metric-main">
                  <div class="metric-value-large">{{ monitorData.memory?.usageRate?.toFixed(1) || '0.0' }}%</div>
                  <div class="metric-label">内存使用率</div>
                </div>
                <el-divider />
                <div class="metric-details">
                  <div class="detail-item">
                    <span class="detail-label">可用内存</span>
                    <span class="detail-value">{{ monitorData.memory?.freeStr || '0 MB' }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">已使用内存</span>
                    <span class="detail-value">{{ monitorData.memory?.usedStr || '0 MB' }}</span>
                  </div>
                </div>
              </div>
            </el-card>

            <!-- 卡片4: 磁盘详情 -->
            <el-card shadow="hover" class="metric-card disk-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><HardDrive /></el-icon>
                  <span class="card-title">磁盘详情</span>
                </div>
              </template>
              <div class="metric-content">
                <div class="metric-main">
                  <div class="metric-value-large">{{ monitorData.disk?.usageRate?.toFixed(1) || '0.0' }}%</div>
                  <div class="metric-label">磁盘使用率</div>
                </div>
                <el-divider />
                <div class="metric-details">
                  <div class="detail-item">
                    <span class="detail-label">磁盘读取速率</span>
                    <span class="detail-value">{{ monitorData.disk?.diskThroughput?.readKbPerSec?.toFixed(2) || '0.00' }} KB/s</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">磁盘写入速率</span>
                    <span class="detail-value">{{ monitorData.disk?.diskThroughput?.writeKbPerSec?.toFixed(2) || '0.00' }} KB/s</span>
                  </div>
                </div>
              </div>
            </el-card>

            <!-- 卡片5: 连接数 -->
            <el-card shadow="hover" class="metric-card connection-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><Connection /></el-icon>
                  <span class="card-title">连接数</span>
                </div>
              </template>
              <div class="metric-content">
                <div class="connection-grid">
                  <div class="connection-item">
                    <div class="connection-value">{{ monitorData.connections?.tcp || 0 }}</div>
                    <div class="connection-label">TCP连接数</div>
                  </div>
                  <div class="connection-item">
                    <div class="connection-value">{{ monitorData.connections?.udp || 0 }}</div>
                    <div class="connection-label">UDP连接数</div>
                  </div>
                  <div class="connection-item">
                    <div class="connection-value">{{ monitorData.connections?.total || 0 }}</div>
                    <div class="connection-label">总连接数</div>
                  </div>
                </div>
              </div>
            </el-card>

            <!-- 卡片6: 系统负载 -->
            <el-card shadow="hover" class="metric-card load-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><TrendCharts /></el-icon>
                  <span class="card-title">系统负载</span>
                </div>
              </template>
              <div class="metric-content">
                <div class="load-grid">
                  <div class="load-item">
                    <div class="load-value">{{ monitorData.load?.load1?.toFixed(2) || '0.00' }}</div>
                    <div class="load-label">1分钟负载</div>
                  </div>
                  <div class="load-item">
                    <div class="load-value">{{ monitorData.load?.load5?.toFixed(2) || '0.00' }}</div>
                    <div class="load-label">5分钟负载</div>
                  </div>
                  <div class="load-item">
                    <div class="load-value">{{ monitorData.load?.load15?.toFixed(2) || '0.00' }}</div>
                    <div class="load-label">15分钟负载</div>
                  </div>
                </div>
              </div>
            </el-card>

            <!-- 卡片7: 网络流量 -->
            <el-card shadow="hover" class="metric-card network-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><Link /></el-icon>
                  <span class="card-title">网络流量</span>
                </div>
              </template>
              <div class="metric-content">
                <div class="network-grid">
                  <div class="network-item">
                    <el-icon class="network-icon download"><Download /></el-icon>
                    <div class="network-info">
                      <div class="network-value">{{ monitorData.network?.receiveKbPerSec?.toFixed(2) || '0.00' }} KB/s</div>
                      <div class="network-label">网络接收流量</div>
                    </div>
                  </div>
                  <div class="network-item">
                    <el-icon class="network-icon upload"><Upload /></el-icon>
                    <div class="network-info">
                      <div class="network-value">{{ monitorData.network?.sendKbPerSec?.toFixed(2) || '0.00' }} KB/s</div>
                      <div class="network-label">网络发送流量</div>
                    </div>
                  </div>
                </div>
              </div>
            </el-card>
          </div>

          <!-- 底部: 其他卡片 -->
          <div class="bottom-section">

            <!-- 卡片8: 系统服务详情 -->
            <el-card shadow="hover" class="services-card" v-if="monitorData.services && monitorData.services.length > 0">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><List /></el-icon>
                  <span class="card-title">系统服务详情 (Top {{ monitorData.services.length }})</span>
                </div>
              </template>
              <el-table :data="monitorData.services" style="width: 100%" stripe>
                <el-table-column prop="name" label="服务名称" min-width="200" />
                <el-table-column label="CPU使用率" width="150" align="center">
                  <template #default="{ row }">
                    <el-tag :type="getCpuTagType(row.cpu)" size="small">{{ row.cpu }}%</el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="内存使用率" width="150" align="center">
                  <template #default="{ row }">
                    <el-tag :type="getMemoryTagType(row.memory)" size="small">{{ row.memory }}%</el-tag>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>

            <!-- 卡片9: CPU Top5进程 -->
            <el-card shadow="hover" class="process-card cpu-top-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><Cpu /></el-icon>
                  <span class="card-title">CPU Top5进程</span>
                </div>
              </template>
              <el-table :data="topProcesses.topCPU" style="width: 100%" stripe>
                <el-table-column prop="pid" label="PID" width="100" align="center">
                  <template #default="{ row }">
                    <el-tag type="warning" size="small">{{ row.pid }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="name" label="进程名称" min-width="180" />
                <el-table-column label="CPU使用率" width="120" align="center">
                  <template #default="{ row }">
                    <el-tag :type="getCpuTagType(row.cpuPercent)" size="small">{{ row.cpuPercent }}%</el-tag>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>

            <!-- 卡片10: 内存Top5进程 -->
            <el-card shadow="hover" class="process-card memory-top-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><Memo /></el-icon>
                  <span class="card-title">内存Top5进程</span>
                </div>
              </template>
              <el-table :data="topProcesses.topMemory" style="width: 100%" stripe>
                <el-table-column prop="pid" label="PID" width="100" align="center">
                  <template #default="{ row }">
                    <el-tag type="warning" size="small">{{ row.pid }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="name" label="进程名称" min-width="180" />
                <el-table-column label="内存使用率" width="120" align="center">
                  <template #default="{ row }">
                    <el-tag :type="getMemoryTagType(row.memPercent)" size="small">{{ row.memPercent }}%</el-tag>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>
          </div>
        </div>

        <!-- 监控图表页面 -->
        <div v-show="activeTab === 'monitor'" class="monitor-content">
          <HostMonitorCharts :hostId="hostId" />
        </div>

        <!-- 告警记录页面 -->
        <div v-show="activeTab === 'alert'" class="alert-content">
          <el-card shadow="hover" class="alert-card">
            <template #header>
              <div class="card-header alert-header">
                <div class="alert-header-left">
                  <el-icon class="card-icon"><List /></el-icon>
                  <span class="card-title">历史告警记录</span>
                  <el-tag size="small" type="info">
                    IP: {{ alertQueryIp || '未获取到IP' }}
                  </el-tag>
                </div>
                <el-button
                  type="warning"
                  size="small"
                  :icon="Refresh"
                  :loading="alertLoading"
                  @click="fetchAlertRecords"
                >
                  刷新
                </el-button>
              </div>
            </template>

            <div v-if="!alertQueryIp" class="alert-empty-wrap">
              <el-empty description="当前主机未获取到可查询的 IP" />
            </div>

            <template v-else>
              <el-table
                :data="alertRecords"
                v-loading="alertLoading"
                stripe
                border
                style="width: 100%"
                :header-cell-style="{ background: '#eef1f6', color: '#606266' }"
              >
                <el-table-column prop="id" label="ID" width="80" align="center" />

                <el-table-column label="状态" width="100" align="center">
                  <template #default="{ row }">
                    <el-tag :type="getAlertStatusType(row.status)" size="small">
                      {{ getAlertStatusLabel(row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>

                <el-table-column label="级别" width="90" align="center">
                  <template #default="{ row }">
                    <el-tag :type="getAlertSeverityType(row.severity)" size="small">
                      {{ getAlertSeverityLabel(row.severity) }}
                    </el-tag>
                  </template>
                </el-table-column>

                <el-table-column prop="policyName" label="策略" min-width="140" show-overflow-tooltip />
                <el-table-column prop="ruleName" label="规则" min-width="180" show-overflow-tooltip />
                <el-table-column prop="targetName" label="目标" min-width="180" show-overflow-tooltip />
                <el-table-column prop="metricName" label="指标" min-width="180" show-overflow-tooltip />

                <el-table-column label="当前值 / 阈值" width="140" align="center">
                  <template #default="{ row }">
                    <span class="alert-mono-text">
                      {{ getAlertValue(row) }} / {{ row.threshold ?? '-' }}
                    </span>
                  </template>
                </el-table-column>

                <el-table-column prop="startTime" label="开始时间" width="160" align="center" />

                <el-table-column label="持续时长" width="100" align="center">
                  <template #default="{ row }">
                    <span>{{ getAlertDuration(row.startTime, row.endTime) }}</span>
                  </template>
                </el-table-column>

                <el-table-column label="操作" width="100" fixed="right" align="center">
                  <template #default="{ row }">
                    <el-button
                      type="primary"
                      size="small"
                      text
                      @click="handleViewAlertDetail(row)"
                    >
                      详情
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>

              <div class="alert-pagination">
                <el-pagination
                  v-model:current-page="alertPagination.page"
                  v-model:page-size="alertPagination.pageSize"
                  :page-sizes="[10, 20, 50, 100]"
                  :total="alertPagination.total"
                  layout="total, sizes, prev, pager, next, jumper"
                  @current-change="handleAlertPageChange"
                  @size-change="handleAlertSizeChange"
                />
              </div>
            </template>
          </el-card>
        </div>
      </template>
    </el-skeleton>

    <el-drawer
      v-model="alertDrawerVisible"
      title="告警详情"
      size="42%"
      :close-on-click-modal="false"
    >
      <div v-if="alertDetail" class="detail-content">
        <div class="detail-section">
          <h4 class="section-title">基本信息</h4>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">告警状态：</span>
                <el-tag :type="getAlertStatusType(alertDetail.status)">
                  {{ getAlertStatusLabel(alertDetail.status) }}
                </el-tag>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">严重级别：</span>
                <el-tag :type="getAlertSeverityType(alertDetail.severity)">
                  {{ getAlertSeverityLabel(alertDetail.severity) }}
                </el-tag>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">策略名称：</span>
                <span>{{ alertDetail.policyName || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">规则名称：</span>
                <span>{{ alertDetail.ruleName || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">目标类型：</span>
                <span>{{ alertDetail.targetType || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">目标名称：</span>
                <span>{{ alertDetail.targetName || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">指标名称：</span>
                <span>{{ alertDetail.metricName || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">指纹：</span>
                <span>{{ alertDetail.fingerprint || '-' }}</span>
              </div>
            </el-col>
          </el-row>
        </div>

        <div class="detail-section">
          <h4 class="section-title">时间信息</h4>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">开始时间：</span>
                <span>{{ formatAlertTime(alertDetail.startTime) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">结束时间：</span>
                <span>{{ formatAlertTime(alertDetail.endTime) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">持续时长：</span>
                <span>{{ getAlertDuration(alertDetail.startTime, alertDetail.endTime) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">通知时间：</span>
                <span>{{ formatAlertTime(alertDetail.notifyTime) }}</span>
              </div>
            </el-col>
          </el-row>
        </div>

        <div class="detail-section">
          <h4 class="section-title">告警内容</h4>
          <el-row :gutter="20">
            <el-col :span="24">
              <div class="detail-item">
                <span class="label">告警信息：</span>
                <span class="multiline-text">{{ alertDetail.message || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">当前值：</span>
                <span class="alert-mono-text">{{ getAlertValue(alertDetail) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">阈值：</span>
                <span class="alert-mono-text">{{ alertDetail.threshold ?? '-' }}</span>
              </div>
            </el-col>
          </el-row>
        </div>

        <div v-if="Object.keys(parseAlertLabels(alertDetail.labels)).length > 0" class="detail-section">
          <h4 class="section-title">标签信息</h4>
          <div class="labels-wrap">
            <el-tag
              v-for="(value, key) in parseAlertLabels(alertDetail.labels)"
              :key="key"
              size="small"
              class="label-tag"
            >
              {{ key }}={{ value }}
            </el-tag>
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  Refresh,
  Cpu,
  Memo,
  HardDrive,
  Monitor,
  TrendCharts,
  Connection,
  Link,
  List,
  Download,
  Upload
} from '@element-plus/icons-vue'
import api from '@/api'
import { recordApi } from '@/api/alert'
import HostMonitorCharts from './HostMonitorCharts.vue'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const refreshing = ref(false)
const hostInfo = ref({})
const monitorData = ref({})
const topProcesses = ref({ topCPU: [], topMemory: [] })
const activeTab = ref('details')
const alertLoading = ref(false)
const alertRecords = ref([])
const alertDrawerVisible = ref(false)
const alertDetail = ref(null)
const alertPagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
})
// 使用 computed 从路由参数中获取主机ID，防止切换页面后丢失
const hostId = computed(() => route.params.id)
const alertQueryIp = computed(() => hostInfo.value.privateIp || hostInfo.value.publicIp || '')
let refreshTimer = null

onMounted(() => {
  // 检查主机ID是否存在
  if (!hostId.value) {
    ElMessage.error('缺少主机ID参数')
    return
  }

  initData()
  // 每30秒自动刷新监控数据
  refreshTimer = setInterval(() => {
    fetchMonitorData()
    fetchTopProcesses()
  }, 30000)
})

onBeforeUnmount(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})

watch(activeTab, (tab) => {
  if (tab === 'alert') {
    fetchAlertRecords()
  }
})

watch(alertQueryIp, (newIp, oldIp) => {
  if (newIp && newIp !== oldIp) {
    alertPagination.value.page = 1
    if (activeTab.value === 'alert') {
      fetchAlertRecords()
    }
  }
})

// 获取主机基本信息
const fetchHostInfo = async () => {
  try {
    if (!hostId.value) {
      ElMessage.error('缺少主机ID参数')
      return
    }

    const { data: res } = await api.getCmdbHostById(hostId.value)
    if (res.code === 200) {
      hostInfo.value = res.data
    } else {
      ElMessage.error(res.message || '获取主机信息失败')
    }
  } catch (error) {
    console.error('获取主机信息失败:', error)
    ElMessage.error('获取主机信息失败')
  }
}

// 获取监控数据
const fetchMonitorData = async () => {
  try {
    if (!hostId.value) return

    const { data: res } = await api.getHostCompleteMonitorData(hostId.value)
    if (res.code === 200) {
      monitorData.value = res.data
    } else {
      ElMessage.error(res.message || '获取监控数据失败')
    }
  } catch (error) {
    console.error('获取监控数据失败:', error)
  }
}

// 获取Top进程数据
const fetchTopProcesses = async () => {
  try {
    if (!hostId.value) return

    const { data: res } = await api.getHostTopProcesses(hostId.value)
    if (res.code === 200) {
      topProcesses.value = {
        topCPU: res.data.topCPU || [],
        topMemory: res.data.topMemory || []
      }
    } else {
      ElMessage.error(res.message || '获取进程数据失败')
    }
  } catch (error) {
    console.error('获取进程数据失败:', error)
  }
}

const fetchAlertRecords = async () => {
  if (!alertQueryIp.value) {
    alertRecords.value = []
    alertPagination.value.total = 0
    return
  }

  alertLoading.value = true
  try {
    const params = {
      ip: alertQueryIp.value,
      page: alertPagination.value.page,
      pageSize: alertPagination.value.pageSize
    }

    const { data: res } = await recordApi.getList(params)
    if (res.code === 200) {
      alertRecords.value = res.data?.list || []
      alertPagination.value.total = res.data?.total || 0
    } else {
      ElMessage.error(res.message || '获取告警记录失败')
    }
  } catch (error) {
    console.error('获取告警记录失败:', error)
    ElMessage.error('获取告警记录失败')
  } finally {
    alertLoading.value = false
  }
}

const handleAlertPageChange = (page) => {
  alertPagination.value.page = page
  fetchAlertRecords()
}

const handleAlertSizeChange = (size) => {
  alertPagination.value.pageSize = size
  alertPagination.value.page = 1
  fetchAlertRecords()
}

const handleViewAlertDetail = async (row) => {
  try {
    const { data: res } = await recordApi.getDetail(row.id)
    if (res.code === 200) {
      alertDetail.value = {
        ...row,
        ...(res.data || {})
      }
      alertDrawerVisible.value = true
    } else {
      ElMessage.error(res.message || '获取告警详情失败')
    }
  } catch (error) {
    console.error('获取告警详情失败:', error)
    ElMessage.error('获取告警详情失败')
  }
}

// 刷新监控数据
const refreshMonitorData = async () => {
  refreshing.value = true
  await Promise.all([fetchMonitorData(), fetchTopProcesses()])
  refreshing.value = false
  ElMessage.success('监控数据已刷新')
}

// 初始化数据
const initData = async () => {
  loading.value = true
  await Promise.all([fetchHostInfo(), fetchMonitorData(), fetchTopProcesses()])
  loading.value = false
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    1: '认证成功',
    2: '未认证',
    3: '认证失败'
  }
  return statusMap[status] || '未知状态'
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const typeMap = {
    1: 'success',
    2: 'warning',
    3: 'danger'
  }
  return typeMap[status] || 'info'
}

// 获取CPU标签类型
const getCpuTagType = (cpu) => {
  if (cpu < 50) return 'success'
  if (cpu < 80) return 'warning'
  return 'danger'
}

// 获取内存标签类型
const getMemoryTagType = (memory) => {
  if (memory < 60) return 'success'
  if (memory < 80) return 'warning'
  return 'danger'
}

const normalizeAlertStatus = (status) => {
  if (status === null || status === undefined) return ''
  return String(status).trim().toLowerCase()
}

const getAlertStatusType = (status) => {
  const typeMap = {
    firing: 'danger',
    resolved: 'success',
    acknowledged: 'warning'
  }
  return typeMap[normalizeAlertStatus(status)] || 'info'
}

const getAlertStatusLabel = (status) => {
  const labelMap = {
    firing: '告警中',
    resolved: '已恢复',
    acknowledged: '已认领'
  }
  return labelMap[normalizeAlertStatus(status)] || status || '-'
}

const getAlertSeverityType = (severity) => {
  const typeMap = {
    1: 'danger',
    2: 'warning',
    3: 'info'
  }
  return typeMap[severity] || 'info'
}

const getAlertSeverityLabel = (severity) => {
  const labelMap = {
    1: '严重',
    2: '警告',
    3: '信息'
  }
  return labelMap[severity] || '未知'
}

const getAlertValue = (record) => {
  return record?.currentValue ?? record?.value ?? '-'
}

const formatAlertTime = (time) => {
  if (!time) return '-'
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

const getAlertDuration = (startTime, endTime) => {
  if (!startTime) return '-'
  const end = endTime ? dayjs(endTime) : dayjs()
  const start = dayjs(startTime)
  const diff = end.diff(start, 'second')

  if (diff < 60) return `${diff}秒`
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时`
  return `${Math.floor(diff / 86400)}天`
}

const parseAlertLabels = (labels) => {
  if (!labels) return {}
  if (typeof labels === 'object') return labels

  try {
    return JSON.parse(labels)
  } catch (error) {
    console.error('解析告警 labels 失败:', error)
    return {}
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}
</script>

<style scoped>
/* 容器样式 */
.host-detail-container {
  padding: 0;
  background: #f5f7fa;
  min-height: 100vh;
  position: relative;
}

/* 页面头部 */
.page-header {
  background: #ffffff;
  padding: 16px 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 10;
  border-bottom: 2px solid transparent;
  background-clip: padding-box;
  animation: slideDown 0.5s ease-out;
}

.page-header::before {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, #667eea, #764ba2, #667eea);
  background-size: 200% 100%;
  animation: gradientMove 3s linear infinite;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes gradientMove {
  0% {
    background-position: 0% 0%;
  }
  100% {
    background-position: 200% 0%;
  }
}

.header-content {
  display: flex;
  align-items: center;
  gap: 16px;
  max-width: 100%;
  margin: 0 auto;
  padding: 0 4px;
  justify-content: center;
  position: relative;
}

.back-btn {
  font-size: 14px;
  padding: 8px 16px;
  position: absolute;
  left: 4px;
  background-color: #67c23a;
  color: #ffffff;
  border: none;
}

.back-btn:hover {
  background-color: #5daf34;
  color: #ffffff;
}

.page-title {
  flex: none;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

/* Tab 容器 */
.tab-container {
  background: #ffffff;
  padding: 0 24px;
  margin-bottom: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  position: relative;
  z-index: 9;
  animation: slideDown 0.6s ease-out 0.1s both;
}

.detail-tabs {
  max-width: 100%;
  margin: 0 auto;
}

/* 内容区域 */
.details-content {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 20px 20px;
  position: relative;
  z-index: 1;
}

.monitor-content,
.alert-content {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 20px 20px;
}

/* 顶部布局：左侧主机信息 + 右侧4个卡片(2x2) */
.details-content {
  display: grid;
  grid-template-columns: 400px 1fr;
  gap: 12px;
  grid-template-rows: auto auto;
}

.left-section {
  grid-row: 1 / 2;
  grid-column: 1;
}

.right-section {
  grid-row: 1 / 2;
  grid-column: 2;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  align-items: start;
}

/* 底部区域 */
.bottom-section {
  grid-row: 2;
  grid-column: 1 / -1;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 12px;
  margin-top: 0;
}

.monitor-content,
.alert-content {
  display: block;
  min-height: auto;
}

.alert-card {
  border-radius: 12px;
  border: 1px solid #e4e7ed;
}

.alert-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.alert-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.alert-empty-wrap {
  padding: 24px 0;
}

.alert-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.alert-mono-text {
  font-family: monospace;
  font-size: 12px;
}

.detail-content {
  padding-right: 8px;
}

.detail-section {
  margin-bottom: 24px;
}

.section-title {
  margin: 0 0 16px;
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.detail-item {
  margin-bottom: 14px;
  line-height: 1.7;
  word-break: break-all;
}

.detail-item .label {
  color: #909399;
}

.labels-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.label-tag {
  margin-right: 0;
}

.multiline-text {
  white-space: pre-wrap;
}

/* 卡片通用样式 */
.info-card,
.metric-card,
.services-card {
  border-radius: 12px;
  border: 1px solid #e4e7ed;
  background: #ffffff;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  animation: fadeInUp 0.6s ease-out both;
}

/* 卡片动画延迟 */
.info-card {
  animation-delay: 0.1s;
}

.metric-card:nth-child(1) {
  animation-delay: 0.2s;
}

.metric-card:nth-child(2) {
  animation-delay: 0.25s;
}

.metric-card:nth-child(3) {
  animation-delay: 0.3s;
}

.metric-card:nth-child(4) {
  animation-delay: 0.35s;
}

.metric-card:nth-child(5) {
  animation-delay: 0.4s;
}

.metric-card:nth-child(6) {
  animation-delay: 0.45s;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 卡片发光边框效果 */
.info-card::before,
.metric-card::before,
.services-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.3), transparent);
  transition: left 0.5s ease;
}

.metric-card :deep(.el-card__header) {
  padding: 12px 16px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

.metric-card :deep(.el-card__body) {
  padding: 12px 16px;
}

.info-card:hover,
.metric-card:hover,
.services-card:hover {
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.3);
  transform: translateY(-4px) scale(1.02);
  border-color: rgba(102, 126, 234, 0.5);
}

.info-card:hover::before,
.metric-card:hover::before,
.services-card:hover::before {
  left: 100%;
}

/* 卡片1: 主机信息 */
.info-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.info-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 10px;
  background: transparent;
  border-radius: 6px;
  border: none;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.info-row::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 3px;
  background: transparent;
  transform: scaleY(0);
  transition: transform 0.3s ease;
}

.info-row:hover {
  background: #f5f7fa;
  transform: translateX(5px);
  box-shadow: none;
}

.info-row:hover::before {
  transform: scaleY(1);
}

.info-row .label {
  font-size: 13px;
  color: #909399;
  white-space: nowrap;
}

.info-row .value {
  font-size: 13px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-align: right;
  word-break: break-all;
}

/* 卡片头部 */
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-icon {
  font-size: 16px;
  color: #667eea;
  animation: iconPulse 2s ease-in-out infinite;
  filter: drop-shadow(0 0 4px rgba(102, 126, 234, 0.5));
}

@keyframes iconPulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
}

.card-title {
  flex: 1;
  font-size: 15px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* CPU卡片 */
.cpu-card .card-icon {
  color: #409eff;
  filter: drop-shadow(0 0 8px rgba(64, 158, 255, 0.6));
}

.cpu-card {
  background: #ffffff;
}

/* 内存卡片 */
.memory-card .card-icon {
  color: #67c23a;
  filter: drop-shadow(0 0 8px rgba(103, 194, 58, 0.6));
}

.memory-card {
  background: #ffffff;
}

/* 负载卡片 */
.load-card .card-icon {
  color: #e6a23c;
  filter: drop-shadow(0 0 8px rgba(230, 162, 60, 0.6));
}

.load-card {
  background: #ffffff;
}

/* 连接数卡片 */
.connection-card .card-icon {
  color: #f56c6c;
  filter: drop-shadow(0 0 8px rgba(245, 108, 108, 0.6));
}

.connection-card {
  background: #ffffff;
}

/* 磁盘卡片 */
.disk-card .card-icon {
  color: #e6a23c;
  filter: drop-shadow(0 0 8px rgba(230, 162, 60, 0.6));
}

.disk-card {
  background: #ffffff;
}

/* 网络卡片 */
.network-card .card-icon {
  color: #409eff;
  filter: drop-shadow(0 0 8px rgba(64, 158, 255, 0.6));
}

.network-card {
  background: #ffffff;
}

/* 系统信息卡片 */
.system-card .card-icon {
  color: #909399;
  filter: drop-shadow(0 0 8px rgba(144, 147, 153, 0.6));
}

.system-card {
  background: #ffffff;
}

/* 服务列表卡片 */
.services-card {
  width: 100%;
}

.services-card .card-icon {
  color: #67c23a;
}

/* 进程列表卡片 */
.process-card {
  width: 100%;
}

.cpu-top-card .card-icon {
  color: #409eff;
  filter: drop-shadow(0 0 8px rgba(64, 158, 255, 0.6));
}

.memory-top-card .card-icon {
  color: #67c23a;
  filter: drop-shadow(0 0 8px rgba(103, 194, 58, 0.6));
}

/* 指标内容 */
.metric-content {
  padding: 2px 0;
}

.metric-main {
  text-align: center;
  padding: 8px 0;
}

.metric-value-large {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
  margin-bottom: 4px;
  animation: numberGlow 2s ease-in-out infinite;
  text-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
  position: relative;
}

@keyframes numberGlow {
  0%, 100% {
    filter: brightness(1);
  }
  50% {
    filter: brightness(1.2);
  }
}

.metric-label {
  font-size: 13px;
  color: #909399;
}

.metric-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-radius: 6px;
  border: 1px solid rgba(102, 126, 234, 0.1);
  transition: all 0.3s ease;
}

.detail-item:hover {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  transform: translateX(3px);
  box-shadow: 0 2px 6px rgba(102, 126, 234, 0.15);
}

.detail-label {
  font-size: 12px;
  color: #909399;
}

.detail-value {
  font-size: 13px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 负载网格 */
.load-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  padding: 8px 0;
}

.load-item {
  text-align: center;
  padding: 8px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.1);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.load-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}

.load-value {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 4px;
  animation: numberFloat 3s ease-in-out infinite;
}

@keyframes numberFloat {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-2px);
  }
}

.load-label {
  font-size: 12px;
  color: #909399;
}

/* 连接数网格 */
.connection-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  padding: 8px 0;
}

.connection-item {
  text-align: center;
  padding: 8px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.1);
  transition: all 0.3s ease;
}

.connection-item:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}

.connection-value {
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 4px;
  animation: numberPulse 2s ease-in-out infinite;
}

@keyframes numberPulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.connection-label {
  font-size: 12px;
  color: #909399;
}

/* 网络流量网格 */
.network-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
  padding: 8px 0;
}

.network-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.1);
  transition: all 0.3s ease;
}

.network-item:hover {
  transform: translateX(5px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}

.network-icon {
  font-size: 24px;
  animation: iconRotate 3s linear infinite;
}

@keyframes iconRotate {
  0%, 90%, 100% {
    transform: rotate(0deg);
  }
  95% {
    transform: rotate(10deg);
  }
}

.network-icon.download {
  color: #67c23a;
  filter: drop-shadow(0 0 8px rgba(103, 194, 58, 0.5));
  animation: iconBounce 2s ease-in-out infinite;
}

@keyframes iconBounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-3px);
  }
}

.network-icon.upload {
  color: #409eff;
  filter: drop-shadow(0 0 8px rgba(64, 158, 255, 0.5));
  animation: iconBounce 2s ease-in-out infinite 0.5s;
}

.network-info {
  flex: 1;
}

.network-value {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 2px;
}

.network-label {
  font-size: 12px;
  color: #909399;
}

/* 响应式 */
@media (max-width: 1200px) {
  .details-content {
    grid-template-columns: 1fr;
    grid-template-rows: auto auto auto;
  }

  .left-section {
    grid-row: 1;
    grid-column: 1;
  }

  .right-section {
    grid-row: 2;
    grid-column: 1;
  }

  .bottom-section {
    grid-row: 3;
    grid-column: 1;
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .host-detail-container {
    padding: 0;
  }

  .page-header {
    padding: 12px 16px;
  }

  .header-content {
    gap: 12px;
  }

  .page-title {
    font-size: 16px;
  }

  .tab-container {
    padding: 0 16px;
  }

  .details-content {
    padding: 0 16px 16px;
    grid-template-columns: 1fr;
  }

  .right-section {
    grid-template-columns: 1fr;
    grid-template-rows: auto;
  }

  .bottom-section {
    grid-template-columns: 1fr;
  }

  .metric-value-large {
    font-size: 28px;
  }

  .load-grid,
  .connection-grid {
    grid-template-columns: 1fr;
  }

  .network-grid {
    grid-template-columns: 1fr;
  }
}
</style>
