<script setup>
import { ref, reactive, onMounted, onUnmounted, watch, nextTick } from 'vue'
import dashboardApi from '@/api/dashboard'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'

// Props - 从父组件接收Zabbix账号ID
const props = defineProps({
  accountId: {
    type: Number,
    default: null
  }
})

// 告警名称TOP统计数据
const alertNameStatsData = reactive({
  items: [],
  total: 0,
  loading: false,
  timeRange: 'today' // today, yesterday, last7days, last30days
})

// 告警主机群组TOP统计数据
const alertHostgroupStatsData = reactive({
  items: [],
  total: 0,
  loading: false,
  timeRange: 'today'
})

// 告警通知统计数据
const notificationStatsData = reactive({
  items: [],
  totalCount: 0,
  loading: false,
  timeRange: 'last30days'
})

// 图表实例
let notificationChart = null
let hostgroupChart = null

// 测试数据
const getTestAlertNameStats = () => ({
  items: [
    { name: '磁盘空间不足', count: 245, percentage: 28.5 },
    { name: 'CPU使用率过高', count: 198, percentage: 21.2 },
    { name: '内存不足', count: 156, percentage: 17.1 },
    { name: '网络延迟', count: 134, percentage: 13.2 },
    { name: '服务不可用', count: 112, percentage: 11.2 }
  ],
  total: 845
})

const getTestAlertHostgroupStats = () => ({
  items: [
    { name: '生产环境', count: 156, percentage: 28.5 },
    { name: '测试环境', count: 134, percentage: 21.2 },
    { name: '开发环境', count: 112, percentage: 17.1 },
    { name: '预发布环境', count: 98, percentage: 13.2 },
    { name: '灾备环境', count: 87, percentage: 11.2 }
  ],
  total: 587
})

const getTestNotificationStats = () => ({
  items: [
    { timestamp: Math.floor(Date.now() / 1000) - 86400 * 6, count: 45, percentage: 12.3 },
    { timestamp: Math.floor(Date.now() / 1000) - 86400 * 5, count: 52, percentage: 14.2 },
    { timestamp: Math.floor(Date.now() / 1000) - 86400 * 4, count: 38, percentage: 10.4 },
    { timestamp: Math.floor(Date.now() / 1000) - 86400 * 3, count: 61, percentage: 16.7 },
    { timestamp: Math.floor(Date.now() / 1000) - 86400 * 2, count: 48, percentage: 13.1 },
    { timestamp: Math.floor(Date.now() / 1000) - 86400, count: 55, percentage: 15.0 },
    { timestamp: Math.floor(Date.now() / 1000), count: 67, percentage: 18.3 }
  ],
  totalCount: 366
})

// 加载告警名称TOP统计
const loadAlertNameStats = async () => {
  if (!props.accountId) {
    const testData = getTestAlertNameStats()
    alertNameStatsData.items = testData.items
    alertNameStatsData.total = testData.total
    return
  }

  try {
    alertNameStatsData.loading = true

    const res = await dashboardApi.getZabbixAlertTopByName(
      props.accountId,
      alertNameStatsData.timeRange,
      5
    )

    if (res.data && res.data.code === 200) {
      const data = res.data.data
      alertNameStatsData.items = data.items || []
      alertNameStatsData.total = data.total || 0
    } else {
      const testData = getTestAlertNameStats()
      alertNameStatsData.items = testData.items
      alertNameStatsData.total = testData.total
    }
  } catch (error) {
    console.error('加载告警名称统计失败:', error)
    const testData = getTestAlertNameStats()
    alertNameStatsData.items = testData.items
    alertNameStatsData.total = testData.total
  } finally {
    alertNameStatsData.loading = false
  }
}

// 加载告警主机群组TOP统计
const loadAlertHostgroupStats = async () => {
  if (!props.accountId) {
    const testData = getTestAlertHostgroupStats()
    alertHostgroupStatsData.items = testData.items
    alertHostgroupStatsData.total = testData.total
    await nextTick()
    updateHostgroupChart()
    return
  }

  try {
    alertHostgroupStatsData.loading = true

    console.log('[主机群组] 开始加载数据, accountId:', props.accountId, 'timeRange:', alertHostgroupStatsData.timeRange)

    const res = await dashboardApi.getZabbixAlertTopByHostgroup(
      props.accountId,
      alertHostgroupStatsData.timeRange,
      5
    )

    console.log('[主机群组] API响应:', res.data)

    if (res.data && res.data.code === 200) {
      const data = res.data.data
      alertHostgroupStatsData.items = data.items || []
      alertHostgroupStatsData.total = data.total || 0
    } else {
      const testData = getTestAlertHostgroupStats()
      alertHostgroupStatsData.items = testData.items
      alertHostgroupStatsData.total = testData.total
    }

    console.log('[主机群组] 数据已保存到state, items:', alertHostgroupStatsData.items, 'total:', alertHostgroupStatsData.total)

    // 更新图表
    await nextTick()
    console.log('[主机群组] 准备更新图表, hostgroupChart存在:', !!hostgroupChart, 'items长度:', alertHostgroupStatsData.items.length)
    updateHostgroupChart()
  } catch (error) {
    console.error('[主机群组] 加载告警主机群组统计失败:', error)
    const testData = getTestAlertHostgroupStats()
    alertHostgroupStatsData.items = testData.items
    alertHostgroupStatsData.total = testData.total

    await nextTick()
    updateHostgroupChart()
  } finally {
    alertHostgroupStatsData.loading = false
  }
}

// 加载告警通知统计
const loadNotificationStats = async () => {
  if (!props.accountId) {
    const testData = getTestNotificationStats()
    notificationStatsData.items = testData.items
    notificationStatsData.totalCount = testData.totalCount
    await nextTick()
    updateNotificationChart()
    return
  }

  try {
    notificationStatsData.loading = true

    console.log('[通知统计] 开始加载数据, accountId:', props.accountId, 'timeRange:', notificationStatsData.timeRange)

    const res = await dashboardApi.getZabbixNotificationStats(
      props.accountId,
      notificationStatsData.timeRange
    )

    console.log('[通知统计] API响应:', res.data)

    if (res.data && res.data.code === 200) {
      const data = res.data.data
      notificationStatsData.items = data.items || []
      notificationStatsData.totalCount = data.totalCount || 0
    } else {
      const testData = getTestNotificationStats()
      notificationStatsData.items = testData.items
      notificationStatsData.totalCount = testData.totalCount
    }

    console.log('[通知统计] 数据已保存到state, items:', notificationStatsData.items, 'totalCount:', notificationStatsData.totalCount)

    // 更新图表
    await nextTick()
    console.log('[通知统计] 准备更新图表, notificationChart存在:', !!notificationChart, 'items长度:', notificationStatsData.items.length)
    updateNotificationChart()
  } catch (error) {
    console.error('[通知统计] 加载告警通知统计失败:', error)
    const testData = getTestNotificationStats()
    notificationStatsData.items = testData.items
    notificationStatsData.totalCount = testData.totalCount

    await nextTick()
    updateNotificationChart()
  } finally {
    notificationStatsData.loading = false
  }
}

// 切换时间范围
const changeAlertNameTimeRange = (range) => {
  alertNameStatsData.timeRange = range
  loadAlertNameStats()
}

const changeAlertHostgroupTimeRange = (range) => {
  alertHostgroupStatsData.timeRange = range
  loadAlertHostgroupStats()
}

const changeNotificationTimeRange = (range) => {
  notificationStatsData.timeRange = range
  loadNotificationStats()
}

// 获取表格行的class
const getAlertRowClass = ({ rowIndex }) => {
  if (rowIndex === 0) return 'alert-row-rank-1'
  if (rowIndex === 1) return 'alert-row-rank-2'
  if (rowIndex === 2) return 'alert-row-rank-3'
  if (rowIndex === 3) return 'alert-row-rank-4'
  if (rowIndex === 4) return 'alert-row-rank-5'
  return ''
}

// 初始化告警通知图表
const initNotificationChart = () => {
  console.log('[通知统计] 开始初始化图表')
  const chartDom = document.getElementById('notificationChart')
  console.log('[通知统计] chartDom元素:', chartDom)

  if (!chartDom) {
    console.warn('[通知统计] 未找到图表DOM元素 #notificationChart')
    return
  }

  notificationChart = echarts.init(chartDom)
  console.log('[通知统计] 图表初始化完成, notificationChart:', !!notificationChart)
  // 不在初始化时更新,等数据加载后再更新
}

// 更新告警通知图表
const updateNotificationChart = () => {
  console.log('[通知统计] updateNotificationChart 被调用, notificationChart:', !!notificationChart, 'items:', notificationStatsData.items)

  if (!notificationChart || !notificationStatsData.items.length) {
    console.log('[通知统计] 提前返回 - chart存在:', !!notificationChart, 'items长度:', notificationStatsData.items.length)
    return
  }

  console.log('[通知统计] 开始更新图表配置')

  const dates = notificationStatsData.items.map(item => {
    const date = new Date(item.timestamp * 1000)
    return `${date.getMonth() + 1}/${date.getDate()}`
  })

  const counts = notificationStatsData.items.map(item => item.count)

  console.log('[通知统计] 图表数据处理完成, dates:', dates, 'counts:', counts)

  const option = {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'line',
        lineStyle: {
          color: 'rgba(230, 162, 60, 0.3)',
          width: 2
        }
      },
      backgroundColor: 'rgba(255, 255, 255, 0.98)',
      borderWidth: 0,
      padding: [12, 16],
      textStyle: {
        color: '#333',
        fontSize: 13
      },
      boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
      formatter: (params) => {
        const item = params[0]
        const dataItem = notificationStatsData.items[item.dataIndex]
        return `
          <div style="font-weight: 600; font-size: 14px; color: #1a1a1a; margin-bottom: 8px;">${item.name}</div>
          <div style="display: flex; align-items: center; justify-content: space-between; gap: 20px;">
            <span style="color: #666;">告警数量</span>
            <span style="font-weight: 600; font-size: 16px; color: #e6a23c;">${item.value}</span>
          </div>
          <div style="display: flex; align-items: center; justify-content: space-between; gap: 20px; margin-top: 6px;">
            <span style="color: #666;">占比</span>
            <span style="font-weight: 500; color: #999;">${(dataItem.percentage || 0).toFixed(2)}%</span>
          </div>
        `
      }
    },
    grid: {
      top: 40,
      bottom: 40,
      left: 50,
      right: 30
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: {
        color: '#666',
        fontSize: 12,
        rotate: 30
      },
      axisLine: {
        lineStyle: { color: '#e0e0e0' }
      },
      axisTick: {
        show: false
      }
    },
    yAxis: {
      type: 'value',
      name: '告警数量',
      nameTextStyle: {
        color: '#666',
        fontSize: 12,
        padding: [0, 0, 0, 0]
      },
      splitLine: {
        lineStyle: { color: '#f0f0f0', type: 'dashed' }
      },
      axisLabel: {
        color: '#666',
        fontSize: 12
      }
    },
    series: [
      {
        type: 'line',
        data: counts,
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        itemStyle: {
          color: '#e6a23c',
          borderWidth: 2,
          borderColor: '#fff'
        },
        lineStyle: {
          width: 3,
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: '#e6a23c' },
            { offset: 1, color: '#f5c542' }
          ])
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(230, 162, 60, 0.3)' },
            { offset: 1, color: 'rgba(230, 162, 60, 0.05)' }
          ])
        },
        emphasis: {
          focus: 'series',
          itemStyle: {
            color: '#e6a23c',
            borderWidth: 3,
            borderColor: '#fff',
            shadowBlur: 10,
            shadowColor: 'rgba(230, 162, 60, 0.5)'
          }
        }
      }
    ]
  }

  notificationChart.setOption(option)
  console.log('[通知统计] 图表配置已设置完成')
}

// 窗口大小改变时重绘图表
const handleResize = () => {
  notificationChart?.resize()
  hostgroupChart?.resize()
}

// 初始化主机群组图表
const initHostgroupChart = () => {
  console.log('[主机群组] 开始初始化图表')
  const chartDom = document.getElementById('hostgroupChart')
  console.log('[主机群组] chartDom元素:', chartDom)

  if (!chartDom) {
    console.warn('[主机群组] 未找到图表DOM元素 #hostgroupChart')
    return
  }

  hostgroupChart = echarts.init(chartDom)
  console.log('[主机群组] 图表初始化完成, hostgroupChart:', !!hostgroupChart)
  // 不在初始化时更新,等数据加载后再更新
}

// 更新主机群组图表
const updateHostgroupChart = () => {
  console.log('[主机群组] updateHostgroupChart 被调用, hostgroupChart:', !!hostgroupChart, 'items:', alertHostgroupStatsData.items)

  if (!hostgroupChart || !alertHostgroupStatsData.items.length) {
    console.log('[主机群组] 提前返回 - chart存在:', !!hostgroupChart, 'items长度:', alertHostgroupStatsData.items.length)
    return
  }

  console.log('[主机群组] 开始更新图表配置')

  const groupNames = alertHostgroupStatsData.items.map(item => {
    const name = item.name || ''
    return name.length > 8 ? name.substring(0, 8) + '...' : name
  })
  const counts = alertHostgroupStatsData.items.map(item => item.count || 0)

  console.log('[主机群组] 图表数据处理完成, groupNames:', groupNames, 'counts:', counts)

  // 渐变色配置 - 使用黄色系 (5个柱状图)
  const colors = [
    ['#fadb14', '#ffd666'],  // 亮黄色
    ['#ffc53d', '#ffe58f'],  // 金黄色
    ['#faad14', '#ffd666'],  // 橙黄色
    ['#d4b106', '#fadb14'],  // 深黄色
    ['#ad8b00', '#d48806']   // 暗黄色
  ]

  const option = {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      backgroundColor: 'rgba(255, 255, 255, 0.98)',
      borderWidth: 0,
      padding: [10, 14],
      textStyle: {
        color: '#333',
        fontSize: 12
      },
      boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
      formatter: (params) => {
        const item = alertHostgroupStatsData.items[params[0].dataIndex]
        return `
          <div style="font-size: 12px;">
            <div style="font-weight: 600; margin-bottom: 6px;">${item.name}</div>
            <div style="display: flex; align-items: center; justify-content: space-between; gap: 20px;">
              <span style="color: #666;">告警数量</span>
              <span style="color: #faad14; font-weight: 600;">${item.count}</span>
            </div>
            <div style="display: flex; align-items: center; justify-content: space-between; gap: 20px; margin-top: 4px;">
              <span style="color: #666;">占比</span>
              <span style="color: #999; font-weight: 500;">${(item.percentage || 0).toFixed(2)}%</span>
            </div>
          </div>
        `
      }
    },
    grid: {
      top: 25,
      bottom: 50,
      left: 45,
      right: 15
    },
    xAxis: {
      type: 'category',
      data: groupNames,
      axisLabel: {
        color: '#666',
        fontSize: 10,
        interval: 0,
        rotate: 35
      },
      axisLine: {
        lineStyle: { color: '#e0e0e0' }
      },
      axisTick: {
        show: false
      }
    },
    yAxis: {
      type: 'value',
      name: '告警数量',
      nameTextStyle: {
        color: '#666',
        fontSize: 11,
        padding: [0, 0, 0, 0]
      },
      splitLine: {
        lineStyle: { color: '#f0f0f0', type: 'dashed' }
      },
      axisLabel: {
        color: '#666',
        fontSize: 10
      }
    },
    series: [{
      type: 'bar',
      data: counts.map((value, index) => ({
        value,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: colors[index % colors.length][0] },
            { offset: 1, color: colors[index % colors.length][1] }
          ]),
          borderRadius: [3, 3, 0, 0]
        }
      })),
      barWidth: '35%',
      label: {
        show: true,
        position: 'top',
        color: '#333',
        fontSize: 11,
        fontWeight: 600
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 8,
          shadowColor: 'rgba(250, 173, 20, 0.4)'
        }
      }
    }]
  }

  hostgroupChart.setOption(option)
  console.log('[主机群组] 图表配置已设置完成')
}

// 加载所有数据
const loadAllData = async () => {
  await Promise.all([
    loadAlertNameStats(),
    loadAlertHostgroupStats(),
    loadNotificationStats()
  ])
}

// 监听accountId变化
watch(() => props.accountId, () => {
  loadAllData()
})

// 生命周期
onMounted(async () => {
  // 延迟初始化图表,确保DOM已渲染
  await nextTick()
  setTimeout(async () => {
    initNotificationChart()
    initHostgroupChart()

    // 无论是否有accountId,都加载数据(无accountId时显示测试数据)
    await loadAllData()

    // 数据加载后,确保图表正确渲染尺寸
    await nextTick()
    setTimeout(() => {
      notificationChart?.resize()
      hostgroupChart?.resize()
    }, 200)
  }, 100)

  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  notificationChart?.dispose()
  hostgroupChart?.dispose()
})

// 暴露方法给父组件
defineExpose({
  loadAllData
})
</script>

<template>
  <div class="zabbix-alert-stats-section">
    <!-- 1. 告警名称TOP5 -->
    <div class="stat-card alert-card">
      <div class="card-header-with-tabs">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="#e91e63" style="width: 20px; height: 20px; margin-right: 8px;">
            <path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"/>
          </svg>
          告警名称 Top5
        </div>
        <div class="time-range-tabs">
          <button
            :class="['tab-btn', { active: alertNameStatsData.timeRange === 'today' }]"
            @click="changeAlertNameTimeRange('today')"
          >
            今天
          </button>
          <button
            :class="['tab-btn', { active: alertNameStatsData.timeRange === 'yesterday' }]"
            @click="changeAlertNameTimeRange('yesterday')"
          >
            昨天
          </button>
          <button
            :class="['tab-btn', { active: alertNameStatsData.timeRange === 'last7days' }]"
            @click="changeAlertNameTimeRange('last7days')"
          >
            7天
          </button>
          <button
            :class="['tab-btn', { active: alertNameStatsData.timeRange === 'last30days' }]"
            @click="changeAlertNameTimeRange('last30days')"
          >
            30天
          </button>
        </div>
      </div>
      <div class="alert-stats-table-wrapper" v-loading="alertNameStatsData.loading">
        <div v-if="alertNameStatsData.items.length === 0" class="no-data">暂无数据</div>
        <el-table v-else :data="alertNameStatsData.items" :row-class-name="getAlertRowClass" style="width: 100%">
          <el-table-column type="index" label="排名" width="80" align="center">
            <template #default="{ $index }">
              <div :class="['rank-badge', { 'top-five': $index < 5 }]">{{ $index + 1 }}</div>
            </template>
          </el-table-column>
          <el-table-column prop="name" label="告警名称" min-width="180" show-overflow-tooltip />
          <el-table-column prop="count" label="数量" width="80" align="center" />
          <el-table-column label="占比" width="90" align="center">
            <template #default="{ row }">
              <span>{{ (row.percentage || 0).toFixed(2) }}%</span>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- 2. 告警通知统计 -->
    <div class="stat-card notification-card">
      <div class="card-header-with-tabs">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="#e6a23c" style="width: 20px; height: 20px; margin-right: 8px;">
            <path d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.89 2 2 2zm6-6v-5c0-3.07-1.64-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.63 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2z"/>
          </svg>
          告警通知统计
        </div>
        <div class="time-range-tabs">
          <button
            :class="['tab-btn', { active: notificationStatsData.timeRange === 'today' }]"
            @click="changeNotificationTimeRange('today')"
          >
            今天
          </button>
          <button
            :class="['tab-btn', { active: notificationStatsData.timeRange === 'yesterday' }]"
            @click="changeNotificationTimeRange('yesterday')"
          >
            昨天
          </button>
          <button
            :class="['tab-btn', { active: notificationStatsData.timeRange === 'last7days' }]"
            @click="changeNotificationTimeRange('last7days')"
          >
            7天
          </button>
          <button
            :class="['tab-btn', { active: notificationStatsData.timeRange === 'last30days' }]"
            @click="changeNotificationTimeRange('last30days')"
          >
            30天
          </button>
        </div>
      </div>
      <div v-loading="notificationStatsData.loading" class="notification-chart-wrapper">
        <div v-show="notificationStatsData.items.length === 0" class="no-data">暂无数据</div>
        <div id="notificationChart" class="notification-chart" v-show="notificationStatsData.items.length > 0"></div>
      </div>
    </div>

    <!-- 3. 告警主机群组TOP5 -->
    <div class="stat-card alert-card">
      <div class="card-header-with-tabs">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="#67c23a" style="width: 20px; height: 20px; margin-right: 8px;">
            <path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/>
          </svg>
          告警主机群组 Top5
        </div>
        <div class="time-range-tabs">
          <button
            :class="['tab-btn', { active: alertHostgroupStatsData.timeRange === 'today' }]"
            @click="changeAlertHostgroupTimeRange('today')"
          >
            今天
          </button>
          <button
            :class="['tab-btn', { active: alertHostgroupStatsData.timeRange === 'yesterday' }]"
            @click="changeAlertHostgroupTimeRange('yesterday')"
          >
            昨天
          </button>
          <button
            :class="['tab-btn', { active: alertHostgroupStatsData.timeRange === 'last7days' }]"
            @click="changeAlertHostgroupTimeRange('last7days')"
          >
            7天
          </button>
          <button
            :class="['tab-btn', { active: alertHostgroupStatsData.timeRange === 'last30days' }]"
            @click="changeAlertHostgroupTimeRange('last30days')"
          >
            30天
          </button>
        </div>
      </div>
      <div class="hostgroup-chart-wrapper" v-loading="alertHostgroupStatsData.loading">
        <div v-show="alertHostgroupStatsData.items.length === 0" class="no-data">暂无数据</div>
        <div id="hostgroupChart" class="hostgroup-chart" v-show="alertHostgroupStatsData.items.length > 0"></div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.zabbix-alert-stats-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 20px;

  .stat-card {
    max-height: 400px !important;
    min-height: 400px !important;
    height: 400px !important;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background: white;
    border-radius: 8px;
    padding: 0;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    transition: all 0.3s ease;

    &:hover {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
      transform: translateY(-2px);
    }
  }

  .notification-card {
    grid-column: span 1;
  }
}

.card-header-with-tabs {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 18px;
  border-bottom: 1px solid #f0f0f0;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  display: flex;
  align-items: center;
}

.time-range-tabs {
  display: flex;
  gap: 8px;
  background: #f5f7fa;
  padding: 4px;
  border-radius: 6px;
}

.tab-btn {
  padding: 6px 16px;
  border: none;
  background: transparent;
  color: #666;
  font-size: 13px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s ease;

  &:hover {
    color: #333;
    background: rgba(255, 255, 255, 0.5);
  }

  &.active {
    background: white;
    color: #2196f3;
    font-weight: 500;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }
}

.no-data {
  text-align: center;
  padding: 40px 20px;
  color: #999;
  font-size: 14px;
}

.alert-stats-table-wrapper {
  flex: 1;
  overflow-y: auto;
  padding: 8px 12px;

  :deep(.el-table) {
    font-size: 13px;

    .el-table__header th {
      padding: 8px 0;
      font-size: 13px;
      background: #fafafa;
      font-weight: 600;
    }

    .el-table__row td {
      padding: 6px 0;
    }

    .el-table__body {
      font-size: 13px;
    }

    .cell {
      padding: 0 6px;
      line-height: 1.4;
    }
  }
}

.rank-badge {
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e0e0e0, #c0c0c0);
  color: #666;
  border-radius: 50%;
  font-size: 13px;
  font-weight: bold;
  margin: 0 auto;

  &.top-five {
    color: white;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
  }
}

:deep(.alert-row-rank-1) {
  background-color: #fff1f0 !important;

  .rank-badge {
    background: linear-gradient(135deg, #ff4d4f, #ff7875);
  }
}

:deep(.alert-row-rank-2) {
  background-color: #fff7e6 !important;

  .rank-badge {
    background: linear-gradient(135deg, #ff9c6e, #ffbb96);
  }
}

:deep(.alert-row-rank-3) {
  background-color: #fffbe6 !important;

  .rank-badge {
    background: linear-gradient(135deg, #ffc53d, #ffd666);
  }
}

:deep(.alert-row-rank-4) {
  background-color: #f6ffed !important;

  .rank-badge {
    background: linear-gradient(135deg, #95de64, #b7eb8f);
  }
}

:deep(.alert-row-rank-5) {
  background-color: #e6f7ff !important;

  .rank-badge {
    background: linear-gradient(135deg, #40a9ff, #69c0ff);
  }
}

.notification-chart-wrapper {
  flex: 1;
  overflow: hidden;
  padding: 10px 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  min-height: 0;
}

.notification-chart {
  width: calc(100% - 36px);
  height: calc(100% - 20px);
  position: absolute;
  top: 10px;
  left: 18px;
}

.hostgroup-chart-wrapper {
  flex: 1;
  overflow: hidden;
  padding: 10px 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  min-height: 0;
}

.hostgroup-chart {
  width: calc(100% - 36px);
  height: calc(100% - 20px);
  position: absolute;
  top: 10px;
  left: 18px;
}

@media (max-width: 1600px) {
  .zabbix-alert-stats-section {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 1200px) {
  .zabbix-alert-stats-section {
    grid-template-columns: 1fr;
  }
}
</style>
