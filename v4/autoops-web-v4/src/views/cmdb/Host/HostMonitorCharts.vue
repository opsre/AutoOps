<template>
  <div class="monitor-charts-container">
    <!-- 时间选择器 -->
    <div class="time-selector">
      <el-radio-group v-model="timeRangeType" size="small" @change="handleTimeRangeChange">
        <el-radio-button label="preset">预设时间</el-radio-button>
        <el-radio-button label="custom">自定义时间</el-radio-button>
      </el-radio-group>

      <div v-if="timeRangeType === 'preset'" class="preset-buttons">
        <el-button
          v-for="duration in presetDurations"
          :key="duration.value"
          size="small"
          :type="activeDuration === duration.value ? 'primary' : ''"
          @click="selectPresetDuration(duration.value)"
        >
          {{ duration.label }}
        </el-button>
      </div>

      <div v-else class="custom-time">
        <el-date-picker
          v-model="customTimeRange"
          type="datetimerange"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          @change="handleCustomTimeChange"
        />
        <el-button
          size="small"
          type="primary"
          style="margin-left: 10px"
          @click="handleQueryClick"
          :disabled="!customTimeRange || customTimeRange.length !== 2"
        >
          查询
        </el-button>
      </div>
    </div>

    <!-- 监控卡片网格 -->
    <div class="charts-grid">
      <!-- 第一排 -->
      <el-row :gutter="20" class="chart-row">
        <el-col :span="6" v-for="card in metricsConfig[0]" :key="card.key">
          <el-card class="metric-card">
            <template #header>
              <div class="card-header">
                <span class="card-title">{{ card.title }}</span>
                <el-tag size="small" :type="getMetricTagType(card)">
                  {{ getCurrentValue(card) }}
                </el-tag>
              </div>
            </template>
            <div class="chart-container" :ref="el => setChartRef(card.key, el)"></div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 第二排 -->
      <el-row :gutter="20" class="chart-row">
        <el-col :span="6" v-for="card in metricsConfig[1]" :key="card.key">
          <el-card class="metric-card">
            <template #header>
              <div class="card-header">
                <span class="card-title">{{ card.title }}</span>
                <div v-if="card.series" class="multi-tags">
                  <el-tag
                    v-for="(seriesName, idx) in card.series"
                    :key="seriesName"
                    size="small"
                    :type="getSeriesTagType(idx)"
                    style="margin-left: 5px;"
                  >
                    {{ seriesName }}: {{ getCurrentSeriesValue(card, idx) }}
                  </el-tag>
                </div>
                <el-tag v-else size="small" type="info">
                  {{ getCurrentValue(card) }}
                </el-tag>
              </div>
            </template>
            <div class="chart-container" :ref="el => setChartRef(card.key, el)"></div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 第三排 -->
      <el-row :gutter="20" class="chart-row">
        <el-col :span="6" v-for="card in metricsConfig[2]" :key="card.key">
          <el-card class="metric-card">
            <template #header>
              <div class="card-header">
                <span class="card-title">{{ card.title }}</span>
                <div v-if="card.series" class="multi-tags">
                  <el-tag
                    v-for="(seriesName, idx) in card.series"
                    :key="seriesName"
                    size="small"
                    :type="getSeriesTagType(idx)"
                    style="margin-left: 5px;"
                  >
                    {{ seriesName }}: {{ getCurrentSeriesValue(card, idx) }}
                  </el-tag>
                </div>
                <el-tag v-else size="small" type="info">
                  {{ getCurrentValue(card) }}
                </el-tag>
              </div>
            </template>
            <div class="chart-container" :ref="el => setChartRef(card.key, el)"></div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import * as echarts from 'echarts'
import api from '@/api'
import { ElMessage } from 'element-plus'

const props = defineProps({
  hostId: {
    type: String,
    required: true
  }
})

// 时间选择器
const timeRangeType = ref('preset')
const activeDuration = ref('1h')
const presetDurations = [
  { label: '30分钟', value: '30m' },
  { label: '1小时', value: '1h' },
  { label: '3小时', value: '3h' },
  { label: '6小时', value: '6h' },
  { label: '12小时', value: '12h' },
  { label: '24小时', value: '24h' }
]
const customTimeRange = ref([])

// 图表实例管理
const chartInstances = {}
const chartRefs = {}

// 监控指标配置
const metricsConfig = [
  // 第一排
  [
    { key: 'cpu', title: 'CPU使用率', unit: '%', apiField: 'cpu', color: '#00d8ff' },
    { key: 'memory', title: '内存使用率', unit: '%', apiField: 'memory', color: '#67c23a' },
    { key: 'disk', title: '磁盘使用率', unit: '%', apiField: 'disk', color: '#e6a23c' },
    {
      key: 'load',
      title: '系统负载',
      unit: '',
      apiFields: ['load1min', 'load5min', 'load15min'],
      series: ['1分钟', '5分钟', '15分钟'],
      colors: ['#00d8ff', '#ff9800', '#f56c6c']
    }
  ],
  // 第二排
  [
    {
      key: 'diskIO',
      title: '磁盘I/O',
      unit: 'KB/s',
      apiFields: ['diskReadKB', 'diskWriteKB'],
      series: ['读取', '写入'],
      colors: ['#00d8ff', '#ff9800']
    },
    {
      key: 'network',
      title: '网络带宽',
      unit: 'KB/s',
      apiFields: ['networkReceive', 'networkSend'],
      series: ['接收', '发送'],
      colors: ['#17a2b8', '#6610f2']
    },
    { key: 'processes', title: '进程总数', unit: '个', apiField: 'totalProcesses', color: '#409eff' },
    {
      key: 'connections',
      title: 'TCP/UDP连接',
      unit: '个',
      apiFields: ['tcpConnections', 'udpConnections'],
      series: ['TCP', 'UDP'],
      colors: ['#f56c6c', '#e6a23c']
    }
  ],
  // 第三排
  [
    {
      key: 'memoryDetail',
      title: '内存详情',
      unit: 'GB',
      apiFields: ['memoryFree', 'memoryUsed'],
      series: ['可用', '已用'],
      colors: ['#67c23a', '#e6a23c']
    },
    {
      key: 'cpuDetail',
      title: 'CPU详情',
      unit: '%',
      apiFields: ['cpuSystem', 'cpuUser'],
      series: ['内核态', '用户态'],
      chartType: 'stack',
      colors: ['#409eff', '#67c23a']
    },
    {
      key: 'diskDetail',
      title: '磁盘空间',
      unit: 'GB',
      apiFields: ['diskTotal', 'diskUsed'],
      series: ['总空间', '已用'],
      colors: ['#909399', '#e6a23c']
    },
    { key: 'openFiles', title: '打开文件数', unit: '个', apiField: 'openFiles', color: '#67c23a' }
  ]
]

// 当前数据存储
const metricsData = ref({})
const refreshInterval = ref(null)

// 设置图表ref
const setChartRef = (key, el) => {
  if (el) {
    chartRefs[key] = el
  }
}

// 获取当前值（单指标）
const getCurrentValue = (card) => {
  if (!metricsData.value || !card.apiField) return '-'
  const data = metricsData.value[card.apiField]
  if (!data || data.length === 0) return '-'
  const latestValue = data[data.length - 1]?.value
  return latestValue !== undefined ? `${latestValue.toFixed(2)} ${card.unit}` : '-'
}

// 获取当前值（多指标）
const getCurrentSeriesValue = (card, idx) => {
  if (!metricsData.value || !card.apiFields) return '-'
  const field = card.apiFields[idx]
  const data = metricsData.value[field]
  if (!data || data.length === 0) return '-'
  const latestValue = data[data.length - 1]?.value
  return latestValue !== undefined ? `${latestValue.toFixed(2)} ${card.unit}` : '-'
}

// 获取标签类型
const getMetricTagType = (card) => {
  if (!card.apiField) return 'info'
  const data = metricsData.value[card.apiField]
  if (!data || data.length === 0) return 'info'
  const value = data[data.length - 1]?.value
  if (value === undefined) return 'info'

  // CPU、内存、磁盘使用率判断
  if (card.apiField === 'cpu' || card.apiField === 'memory' || card.apiField === 'disk') {
    if (value < 60) return 'success'
    if (value < 80) return 'warning'
    return 'danger'
  }

  return 'info'
}

// 获取系列标签类型
const getSeriesTagType = (idx) => {
  const types = ['success', 'warning', 'info', 'danger']
  return types[idx % types.length]
}

// 时间格式化
const formatDateTime = (date) => {
  if (!date) return ''
  const d = date instanceof Date ? date : new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day}+${hours}:${minutes}:${seconds}`
}

// 获取监控数据
const fetchMetricsData = async () => {
  try {
    console.log('=== 开始获取监控数据 ===')
    console.log('hostId:', props.hostId)

    const params = { id: props.hostId }

    if (timeRangeType.value === 'preset') {
      params.duration = activeDuration.value
      console.log('使用预设时间:', params.duration)
    } else if (customTimeRange.value && customTimeRange.value.length === 2) {
      params.start = formatDateTime(customTimeRange.value[0])
      params.end = formatDateTime(customTimeRange.value[1])
      console.log('使用自定义时间:', params.start, '至', params.end)
    }

    console.log('发送API请求，参数:', params)
    const { data: res } = await api.getHostMonitorData(params)
    console.log('API响应成功，数据:', res)

    if (res.code === 200) {
      metricsData.value = res.data
      console.log('开始更新图表...')
      updateAllCharts()
    } else {
      console.error('API返回错误:', res.message)
      ElMessage.error(res.message || '获取监控数据失败')
    }
  } catch (error) {
    console.error('获取监控数据失败 - 详细错误:', error)
    console.error('错误堆栈:', error.stack)
    ElMessage.error('获取监控数据失败: ' + (error.message || '未知错误'))
  }
}

// 更新所有图表
const updateAllCharts = () => {
  // 使用nextTick确保DOM已完全更新
  setTimeout(() => {
    metricsConfig.forEach(row => {
      row.forEach(card => {
        renderChart(card)
      })
    })
  }, 150)
}

// 渲染单个图表
const renderChart = (config) => {
  const chartDom = chartRefs[config.key]
  if (!chartDom) return

  // 确保容器有高度
  if (chartDom.offsetHeight === 0) {
    console.warn(`Chart container ${config.key} has no height`)
    return
  }

  let chartInstance = chartInstances[config.key]
  if (!chartInstance) {
    chartInstance = echarts.init(chartDom, null, { renderer: 'canvas' })
    chartInstances[config.key] = chartInstance
  }

  // 准备数据
  let xAxisData = []
  let seriesData = []

  if (config.apiField) {
    // 单指标
    const data = metricsData.value[config.apiField] || []
    xAxisData = data.map(point => {
      const date = new Date(point.timestamp * 1000)
      return `${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
    })
    seriesData = [{
      name: config.title,
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: data.map(point => point.value.toFixed(2)),
      lineStyle: {
        color: '#409EFF',
        width: 2,
        type: 'solid'
      },
      itemStyle: { color: '#409EFF' }
    }]
  } else if (config.apiFields) {
    // 多指标
    const firstFieldData = metricsData.value[config.apiFields[0]] || []
    xAxisData = firstFieldData.map(point => {
      const date = new Date(point.timestamp * 1000)
      return `${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
    })

    seriesData = config.apiFields.map((field, idx) => {
      const data = metricsData.value[field] || []
      const colors = ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399']
      return {
        name: config.series[idx],
        type: 'line',
        smooth: false,
        symbol: 'none',
        data: data.map(point => point.value.toFixed(2)),
        lineStyle: {
          color: colors[idx % colors.length],
          width: 2,
          type: 'solid'
        },
        itemStyle: { color: colors[idx % colors.length] },
        stack: config.chartType === 'stack' ? 'total' : undefined
      }
    })
  }

  const option = {
    backgroundColor: '#1f1f1f',
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: '#409EFF',
      borderWidth: 1,
      textStyle: { color: '#606266' },
      axisPointer: {
        type: 'line',
        lineStyle: { color: '#409EFF', width: 1, opacity: 0.5 }
      }
    },
    legend: config.series ? {
      data: config.series,
      textStyle: { color: '#ffffff', fontSize: 12 },
      top: 5,
      left: 'center'
    } : undefined,
    grid: {
      left: '8%',
      right: '4%',
      bottom: '8%',
      top: config.series ? '18%' : '8%',
      containLabel: false
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: xAxisData,
      axisLine: { lineStyle: { color: '#4a4a4a' } },
      axisLabel: {
        color: '#cccccc',
        fontSize: 11,
        interval: 'auto'
      }
    },
    yAxis: {
      type: 'value',
      name: config.unit,
      nameTextStyle: {
        color: '#cccccc',
        fontSize: 11
      },
      axisLine: { lineStyle: { color: '#4a4a4a' } },
      axisLabel: {
        color: '#cccccc',
        fontSize: 11
      },
      splitLine: { lineStyle: { color: '#3a3a3a' } },
      // 如果是百分比类型，设置最大值为100
      max: config.unit === '%' ? 100 : undefined,
      min: 0
    },
    series: seriesData
  }

  chartInstance.setOption(option, true)
  // 强制resize确保图表充满容器
  setTimeout(() => {
    chartInstance.resize()
  }, 100)
}

// 时间选择器事件
const handleTimeRangeChange = () => {
  if (timeRangeType.value === 'preset') {
    fetchMetricsData()
  }
}

const selectPresetDuration = (duration) => {
  activeDuration.value = duration
  fetchMetricsData()
}

const handleCustomTimeChange = () => {
  // 时间变化时不自动查询，等待用户点击查询按钮
}

const handleQueryClick = () => {
  if (customTimeRange.value && customTimeRange.value.length === 2) {
    fetchMetricsData()
  }
}

// 自动刷新
const startRefresh = () => {
  stopRefresh()
  fetchMetricsData()
  refreshInterval.value = setInterval(() => {
    fetchMetricsData()
  }, 10000)
}

const stopRefresh = () => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
    refreshInterval.value = null
  }
}

// 窗口大小调整
const handleResize = () => {
  Object.values(chartInstances).forEach(chart => {
    if (chart) chart.resize()
  })
}

// 生命周期
onMounted(() => {
  window.addEventListener('resize', handleResize)
  startRefresh()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  stopRefresh()
  Object.values(chartInstances).forEach(chart => {
    if (chart) chart.dispose()
  })
})

// 监听hostId变化
watch(() => props.hostId, (newVal) => {
  if (newVal) {
    startRefresh()
  }
})
</script>

<style scoped>
.monitor-charts-container {
  padding: 20px;
  background: #f5f7fa;
  min-height: 100%;
  width: 100%;
  box-sizing: border-box;
}

/* 时间选择器 */
.time-selector {
  margin-bottom: 20px;
  padding: 16px;
  background: #ffffff;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.preset-buttons {
  margin-top: 12px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.custom-time {
  margin-top: 12px;
  display: flex;
  align-items: center;
}

/* 图表网格 */
.charts-grid {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
}

.chart-row {
  margin-bottom: 0 !important;
  width: 100%;
}

/* 指标卡片 */
.metric-card {
  background: #1f1f1f;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  transition: all 0.3s ease;
  height: 380px;
  width: 100%;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.metric-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
  border-color: #409EFF;
}

.metric-card :deep(.el-card__header) {
  padding: 12px 16px;
  background: #2a2a2a;
  border-bottom: 1px solid #3a3a3a;
}

.metric-card :deep(.el-card__body) {
  padding: 12px;
  height: calc(100% - 53px);
  display: flex;
  flex-direction: column;
  background: #1f1f1f;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.card-title {
  font-size: 14px;
  font-weight: 600;
  color: #ffffff;
}

.multi-tags {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
  align-items: center;
}

/* 图表容器 */
.chart-container {
  width: 100%;
  flex: 1;
  min-height: 0;
  height: 100%;
}

/* 响应式布局 */
.chart-row :deep(.el-col) {
  padding-left: 10px;
  padding-right: 10px;
  margin-bottom: 0;
}

@media (min-width: 1920px) {
  .metric-card {
    height: 420px;
  }
}

@media (max-width: 1600px) {
  .metric-card {
    height: 360px;
  }
}

@media (max-width: 1400px) {
  .chart-row :deep(.el-col-6) {
    max-width: 50%;
    flex: 0 0 50%;
  }

  .metric-card {
    height: 380px;
  }
}

@media (max-width: 768px) {
  .monitor-charts-container {
    padding: 10px;
  }

  .chart-row :deep(.el-col-6) {
    max-width: 100%;
    flex: 0 0 100%;
  }

  .metric-card {
    margin-bottom: 15px;
    height: 320px;
  }
}

/* Element Plus 样式覆盖 */
:deep(.el-radio-button__inner) {
  background-color: #ffffff;
  color: #606266;
  border-color: #dcdfe6;
}

:deep(.el-radio-button__orig-radio:checked + .el-radio-button__inner) {
  background-color: #409EFF;
  color: #ffffff;
  border-color: #409EFF;
}

:deep(.el-button) {
  background-color: #ffffff;
  color: #606266;
  border-color: #dcdfe6;
}

:deep(.el-button:hover) {
  background-color: #ecf5ff;
  color: #409EFF;
  border-color: #c6e2ff;
}

:deep(.el-button--primary) {
  background-color: #409EFF;
  color: #ffffff;
  border-color: #409EFF;
}

:deep(.el-button--primary:hover) {
  background-color: #66b1ff;
  color: #ffffff;
  border-color: #66b1ff;
}

:deep(.el-date-editor) {
  background-color: #ffffff;
  border-color: #dcdfe6;
}

:deep(.el-date-editor .el-range-input) {
  background-color: #ffffff;
  color: #606266;
}

:deep(.el-date-editor .el-range-separator) {
  color: #606266;
}

:deep(.el-tag) {
  border-radius: 4px;
}
</style>
