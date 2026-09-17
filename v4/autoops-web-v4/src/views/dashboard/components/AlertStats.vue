<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import dashboardApi from '@/api/dashboard'

// Props - 从父组件接收告警源信息
const props = defineProps({
  selectedAlertSource: {
    type: Number,
    default: null
  }
})

// 告警统计数据标签类型
const alertStatsLabel = ref('check')

// 告警统计数据
const alertStatsData = reactive({
  items: [],
  total: 0,
  compareTotal: 0,
  totalChange: 0,
  totalChangeRate: 0,
  loading: false,
  timeRange: 'yesterday' // today, yesterday, last7days, last30days, today-vs-yesterday, thisweek-vs-lastweek - 默认显示昨天
})

// 告警通知统计数据
const notificationStatsData = reactive({
  items: [],
  totalSms: 0,
  totalEmail: 0,
  totalVoice: 0,
  grandTotal: 0,
  loading: false,
  timeRange: 'last30days' // today, yesterday, last7days, last30days - 默认显示30天
})

// 协作空间统计数据
const channelStatsData = reactive({
  items: [],
  loading: false,
  timeRange: 'last7days' // today, yesterday, last7days, last30days - 默认显示7天
})

// 图表实例
let notificationChart = null
let channelChart = null

// 测试数据
const getTestAlertStatsData = () => ({
  items: [
    { label: '数据库连接超时', totalAlertCnt: 245, totalAlertEventCnt: 89, percentage: 28.5, alertEventCntChangeRate: 12.3 },
    { label: 'CPU使用率过高', totalAlertCnt: 198, totalAlertEventCnt: 72, percentage: 21.2, alertEventCntChangeRate: -5.2 },
    { label: '内存不足', totalAlertCnt: 156, totalAlertEventCnt: 58, percentage: 17.1, alertEventCntChangeRate: 8.9 },
    { label: '磁盘空间告急', totalAlertCnt: 134, totalAlertEventCnt: 45, percentage: 13.2, alertEventCntChangeRate: 3.1 },
    { label: '网络延迟', totalAlertCnt: 112, totalAlertEventCnt: 38, percentage: 11.2, alertEventCntChangeRate: -2.4 }
  ],
  total: 845,
  compareTotal: 750,
  totalChange: 95,
  totalChangeRate: 12.67
})

const getTestNotificationStatsData = () => ({
  items: [
    { timestamp: Math.floor(Date.now() / 1000) - 86400 * 6, smsCount: 45, emailCount: 120, voiceCount: 23 },
    { timestamp: Math.floor(Date.now() / 1000) - 86400 * 5, smsCount: 52, emailCount: 135, voiceCount: 28 },
    { timestamp: Math.floor(Date.now() / 1000) - 86400 * 4, smsCount: 38, emailCount: 98, voiceCount: 19 },
    { timestamp: Math.floor(Date.now() / 1000) - 86400 * 3, smsCount: 61, emailCount: 156, voiceCount: 35 },
    { timestamp: Math.floor(Date.now() / 1000) - 86400 * 2, smsCount: 48, emailCount: 128, voiceCount: 26 },
    { timestamp: Math.floor(Date.now() / 1000) - 86400, smsCount: 55, emailCount: 142, voiceCount: 31 },
    { timestamp: Math.floor(Date.now() / 1000), smsCount: 67, emailCount: 178, voiceCount: 42 }
  ],
  totalSms: 366,
  totalEmail: 957,
  totalVoice: 204,
  grandTotal: 1527
})

const getTestChannelStatsData = () => ({
  items: [
    { channelName: '支付系统', totalIncidentCnt: 156 },
    { channelName: '用户中心', totalIncidentCnt: 134 },
    { channelName: '订单系统', totalIncidentCnt: 112 },
    { channelName: '库存管理', totalIncidentCnt: 98 },
    { channelName: '推荐引擎', totalIncidentCnt: 87 },
    { channelName: '搜索服务', totalIncidentCnt: 76 }
  ]
})

// 加载告警统计数据（根据时间范围）
const loadAlertStatsData = async () => {
  if (!props.selectedAlertSource) {
    const testData = getTestAlertStatsData()
    alertStatsData.items = testData.items
    alertStatsData.total = testData.total
    alertStatsData.compareTotal = testData.compareTotal
    alertStatsData.totalChange = testData.totalChange
    alertStatsData.totalChangeRate = testData.totalChangeRate
    return
  }

  try {
    alertStatsData.loading = true

    // 根据时间范围选择API
    const apiMap = {
      'today': dashboardApi.getAlertStatsToday,
      'yesterday': dashboardApi.getAlertStatsYesterday,
      'last7days': dashboardApi.getAlertStatsLast7Days,
      'last30days': dashboardApi.getAlertStatsLast30Days,
      'today-vs-yesterday': dashboardApi.getAlertStatsTodayVsYesterday,
      'thisweek-vs-lastweek': dashboardApi.getAlertStatsThisWeekVsLastWeek
    }

    const apiFunc = apiMap[alertStatsData.timeRange]
    if (!apiFunc) return

    const res = await apiFunc(props.selectedAlertSource, alertStatsLabel.value)

    if (res.data && (res.data.code === 0 || res.data.code === 200)) {
      const data = res.data.data
      alertStatsData.items = data.items || []
      alertStatsData.total = data.total || 0
      alertStatsData.compareTotal = data.compareTotal || 0
      alertStatsData.totalChange = data.totalChange || 0
      alertStatsData.totalChangeRate = data.totalChangeRate || 0
    } else {
      // 使用测试数据
      const testData = getTestAlertStatsData()
      alertStatsData.items = testData.items
      alertStatsData.total = testData.total
      alertStatsData.compareTotal = testData.compareTotal
      alertStatsData.totalChange = testData.totalChange
      alertStatsData.totalChangeRate = testData.totalChangeRate
    }
  } catch (error) {
    console.error('加载告警统计失败:', error)
    // 使用测试数据
    const testData = getTestAlertStatsData()
    alertStatsData.items = testData.items
    alertStatsData.total = testData.total
    alertStatsData.compareTotal = testData.compareTotal
    alertStatsData.totalChange = testData.totalChange
    alertStatsData.totalChangeRate = testData.totalChangeRate
  } finally {
    alertStatsData.loading = false
  }
}

// 加载告警通知统计数据（根据时间范围）
const loadNotificationStatsData = async () => {
  if (!props.selectedAlertSource) {
    const testData = getTestNotificationStatsData()
    notificationStatsData.items = testData.items
    notificationStatsData.totalSms = testData.totalSms
    notificationStatsData.totalEmail = testData.totalEmail
    notificationStatsData.totalVoice = testData.totalVoice
    notificationStatsData.grandTotal = testData.grandTotal
    nextTick(() => { updateNotificationChart() })
    return
  }

  try {
    notificationStatsData.loading = true

    // 根据时间范围选择API
    const apiMap = {
      'today': dashboardApi.getAlertNotificationsToday,
      'yesterday': dashboardApi.getAlertNotificationsYesterday,
      'last7days': dashboardApi.getAlertNotificationsLast7Days,
      'last30days': dashboardApi.getAlertNotificationsLast30Days
    }

    const apiFunc = apiMap[notificationStatsData.timeRange]
    if (!apiFunc) return

    const res = await apiFunc(props.selectedAlertSource)

    if (res.data && (res.data.code === 0 || res.data.code === 200)) {
      const data = res.data.data
      notificationStatsData.items = data.items || []
      notificationStatsData.totalSms = data.totalSms || 0
      notificationStatsData.totalEmail = data.totalEmail || 0
      notificationStatsData.totalVoice = data.totalVoice || 0
      notificationStatsData.grandTotal = data.grandTotal || 0

      // 更新图表
      nextTick(() => {
        updateNotificationChart()
      })
    } else {
      // 使用测试数据
      const testData = getTestNotificationStatsData()
      notificationStatsData.items = testData.items
      notificationStatsData.totalSms = testData.totalSms
      notificationStatsData.totalEmail = testData.totalEmail
      notificationStatsData.totalVoice = testData.totalVoice
      notificationStatsData.grandTotal = testData.grandTotal

      nextTick(() => {
        updateNotificationChart()
      })
    }
  } catch (error) {
    console.error('加载告警通知统计失败:', error)
    // 使用测试数据
    const testData = getTestNotificationStatsData()
    notificationStatsData.items = testData.items
    notificationStatsData.totalSms = testData.totalSms
    notificationStatsData.totalEmail = testData.totalEmail
    notificationStatsData.totalVoice = testData.totalVoice
    notificationStatsData.grandTotal = testData.grandTotal

    nextTick(() => {
      updateNotificationChart()
    })
  } finally {
    notificationStatsData.loading = false
  }
}

// 加载协作空间统计数据（根据时间范围）
const loadChannelStatsData = async () => {
  if (!props.selectedAlertSource) {
    channelStatsData.items = getTestChannelStatsData().items
    await nextTick()
    setTimeout(() => { !channelChart ? initChannelChart() : updateChannelChart() }, 50)
    return
  }

  try {
    channelStatsData.loading = true

    // 根据时间范围选择API
    const apiMap = {
      'today': dashboardApi.getChannelStatsToday,
      'yesterday': dashboardApi.getChannelStatsYesterday,
      'last7days': dashboardApi.getChannelStatsLast7Days,
      'last30days': dashboardApi.getChannelStatsLast30Days
    }

    const apiFunc = apiMap[channelStatsData.timeRange]
    if (!apiFunc) {
      console.error('业务线告警统计: 未找到对应的API方法')
      return
    }

    const res = await apiFunc(props.selectedAlertSource)

    let sortedItems = []
    if (res.data && (res.data.code === 0 || res.data.code === 200)) {
      const channelData = res.data.data
      // 按照故障总数从多到少排序，取前6条
      sortedItems = (channelData.items || [])
        .sort((a, b) => (b.totalIncidentCnt || 0) - (a.totalIncidentCnt || 0))
        .slice(0, 6)
    } else {
      // 使用测试数据
      sortedItems = getTestChannelStatsData().items
    }

    channelStatsData.items = sortedItems
    console.log('业务线告警统计数据加载成功:', sortedItems.length, '条')

    // 确保DOM已渲染后再更新图表
    await nextTick()

    // 如果图表未初始化,先初始化;否则直接更新
    if (!channelChart) {
      setTimeout(() => {
        initChannelChart()
      }, 50) // 减少延迟
    } else {
      updateChannelChart()
    }
  } catch (error) {
    console.error('加载业务线告警统计失败:', error)
    // 使用测试数据
    const testData = getTestChannelStatsData()
    channelStatsData.items = testData.items
    console.log('业务线告警统计使用测试数据:', testData.items.length, '条')

    await nextTick()
    if (!channelChart) {
      setTimeout(() => {
        initChannelChart()
      }, 50)
    } else {
      updateChannelChart()
    }
  } finally {
    channelStatsData.loading = false
  }
}

// 切换告警统计时间范围
const changeAlertStatsTimeRange = (range) => {
  alertStatsData.timeRange = range
  loadAlertStatsData()
}

// 切换告警通知时间范围
const changeNotificationTimeRange = (range) => {
  notificationStatsData.timeRange = range
  loadNotificationStatsData()
}

// 切换协作空间时间范围
const changeChannelTimeRange = (range) => {
  channelStatsData.timeRange = range
  loadChannelStatsData()
}

// 获取告警表格行的class（为前5名添加不同颜色）
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
  const chartDom = document.getElementById('notificationChart')
  if (!chartDom) return

  notificationChart = echarts.init(chartDom)
  notificationChart.setOption({
    backgroundColor: 'transparent',
    tooltip: { trigger: 'axis' },
    legend: {
      data: ['短信', '邮件', '语音'],
      textStyle: { color: '#666', fontSize: 11 },
      top: 5,
      itemWidth: 14,
      itemHeight: 10
    },
    grid: { top: 35, bottom: 25, left: 45, right: 15 },
    xAxis: {
      type: 'category',
      data: [],
      axisLabel: { color: '#666', fontSize: 11 },
      axisLine: { lineStyle: { color: '#e0e0e0' } }
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: '#f0f0f0' } },
      axisLabel: { color: '#666', fontSize: 11 }
    },
    series: [
      {
        name: '短信',
        type: 'line',
        smooth: true,
        data: [],
        itemStyle: { color: '#409eff' },
        areaStyle: { color: 'rgba(64, 158, 255, 0.2)' }
      },
      {
        name: '邮件',
        type: 'line',
        smooth: true,
        data: [],
        itemStyle: { color: '#67c23a' },
        areaStyle: { color: 'rgba(103, 194, 58, 0.2)' }
      },
      {
        name: '语音',
        type: 'line',
        smooth: true,
        data: [],
        itemStyle: { color: '#e6a23c' },
        areaStyle: { color: 'rgba(230, 162, 60, 0.2)' }
      }
    ]
  })
}

// 更新告警通知图表
const updateNotificationChart = () => {
  if (!notificationChart || notificationStatsData.items.length === 0) return

  const dates = notificationStatsData.items.map(item => {
    const date = new Date(item.timestamp * 1000)
    return `${date.getMonth() + 1}/${date.getDate()}`
  })

  const smsData = notificationStatsData.items.map(item => item.smsCount)
  const emailData = notificationStatsData.items.map(item => item.emailCount)
  const voiceData = notificationStatsData.items.map(item => item.voiceCount)

  notificationChart.setOption({
    xAxis: { data: dates },
    series: [
      { data: smsData },
      { data: emailData },
      { data: voiceData }
    ]
  })
}

// 初始化业务线告警柱状图
const initChannelChart = () => {
  const chartDom = document.getElementById('channelChart')
  if (!chartDom) {
    console.warn('业务线告警图表DOM未找到')
    return
  }

  // 如果图表已存在,先销毁
  if (channelChart) {
    channelChart.dispose()
    channelChart = null
  }

  channelChart = echarts.init(chartDom)
  updateChannelChart()
}

// 更新业务线告警柱状图
const updateChannelChart = () => {
  if (!channelChart) {
    console.warn('业务线告警图表实例不存在,尝试初始化')
    initChannelChart()
    return
  }

  // 如果没有数据,显示空状态
  if (!channelStatsData.items || channelStatsData.items.length === 0) {
    channelChart.clear()
    return
  }

  const channelNames = channelStatsData.items.map(item => {
    // 截取长名称
    const name = item.channelName || ''
    return name.length > 10 ? name.substring(0, 10) + '...' : name
  })
  const incidentCounts = channelStatsData.items.map(item => item.totalIncidentCnt || 0)

  // 渐变色配置
  const colors = [
    ['#ff6b6b', '#ee5a52'],
    ['#f093fb', '#f5576c'],
    ['#4facfe', '#00f2fe'],
    ['#43e97b', '#38f9d7'],
    ['#fa709a', '#fee140'],
    ['#667eea', '#764ba2']
  ]

  channelChart.setOption({
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      formatter: (params) => {
        const item = channelStatsData.items[params[0].dataIndex]
        return `
          <div style="font-size: 12px;">
            <div style="font-weight: 600; margin-bottom: 6px;">${item.channelName}</div>
            <div style="color: #666;">故障总数: <span style="color: #333; font-weight: 600;">${item.totalIncidentCnt}</span></div>
          </div>
        `
      }
    },
    grid: {
      top: 30,
      bottom: 60,
      left: 50,
      right: 20
    },
    xAxis: {
      type: 'category',
      data: channelNames,
      axisLabel: {
        color: '#666',
        fontSize: 11,
        interval: 0,
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
      name: '故障总数',
      nameTextStyle: {
        color: '#666',
        fontSize: 12,
        padding: [0, 0, 0, 0]
      },
      splitLine: {
        lineStyle: { color: '#f0f0f0' }
      },
      axisLabel: {
        color: '#666',
        fontSize: 11
      }
    },
    series: [{
      type: 'bar',
      data: incidentCounts.map((value, index) => ({
        value,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: colors[index][0] },
            { offset: 1, color: colors[index][1] }
          ]),
          borderRadius: [4, 4, 0, 0]
        }
      })),
      barWidth: '50%',
      label: {
        show: true,
        position: 'top',
        color: '#333',
        fontSize: 12,
        fontWeight: 600
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowColor: 'rgba(0, 0, 0, 0.3)'
        }
      }
    }]
  })
}

// 格式化时间（短格式）
const formatSecondsShort = (seconds) => {
  if (seconds < 60) return `${Math.floor(seconds)}秒`
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes}分钟`
  const hours = Math.floor(minutes / 60)
  return `${hours}小时${minutes % 60}分钟`
}

// 窗口大小改变时重绘图表
const handleResize = () => {
  notificationChart?.resize()
  channelChart?.resize()
}

// 加载所有数据 - 并行加载提升性能
const loadAllData = async () => {
  console.log('开始并行加载告警统计数据...')
  const startTime = Date.now()

  await Promise.all([
    loadAlertStatsData(),
    loadNotificationStatsData(),
    loadChannelStatsData()
  ])

  const endTime = Date.now()
  console.log(`告警统计数据加载完成,耗时: ${endTime - startTime}ms`)
}

// 监听 selectedAlertSource 变化
import { watch } from 'vue'
watch(() => props.selectedAlertSource, () => {
  loadAllData()
})

// 监听业务线告警数据变化,确保图表更新
watch(() => channelStatsData.items, async (newItems) => {
  if (newItems && newItems.length > 0) {
    console.log('业务线告警数据变化,触发图表更新:', newItems.length, '条')
    await nextTick()

    // 延迟确保DOM已完全渲染
    setTimeout(() => {
      if (!channelChart) {
        initChannelChart()
      } else {
        updateChannelChart()
      }
    }, 50)
  }
}, { deep: true })

// 生命周期
onMounted(async () => {
  // 先初始化图表DOM
  await nextTick()

  // 立即初始化图表和加载数据,并行执行
  setTimeout(() => {
    initNotificationChart()
    initChannelChart()

    // 无论是否有告警源,都加载数据(无告警源时显示测试数据)
    loadAllData()
  }, 50) // 减少延迟时间

  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  notificationChart?.dispose()
  channelChart?.dispose()
})

// 暴露方法给父组件
defineExpose({
  loadAllData
})
</script>

<template>
  <div class="alert-stats-section">
    <!-- 1. 告警统计 Top5 -->
    <div class="stat-card alert-card">
      <div class="card-header-with-tabs">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="#e91e63" style="width: 20px; height: 20px; margin-right: 8px;">
            <path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"/>
          </svg>
          告警统计 Top5
        </div>
        <div class="time-range-tabs">
          <button
            :class="['tab-btn', { active: alertStatsData.timeRange === 'today' }]"
            @click="changeAlertStatsTimeRange('today')"
          >
            今天
          </button>
          <button
            :class="['tab-btn', { active: alertStatsData.timeRange === 'yesterday' }]"
            @click="changeAlertStatsTimeRange('yesterday')"
          >
            昨天
          </button>
          <button
            :class="['tab-btn', { active: alertStatsData.timeRange === 'last7days' }]"
            @click="changeAlertStatsTimeRange('last7days')"
          >
            7天
          </button>
          <button
            :class="['tab-btn', { active: alertStatsData.timeRange === 'last30days' }]"
            @click="changeAlertStatsTimeRange('last30days')"
          >
            30天
          </button>
          <button
            :class="['tab-btn', { active: alertStatsData.timeRange === 'today-vs-yesterday' }]"
            @click="changeAlertStatsTimeRange('today-vs-yesterday')"
          >
            环比
          </button>
          <button
            :class="['tab-btn', { active: alertStatsData.timeRange === 'thisweek-vs-lastweek' }]"
            @click="changeAlertStatsTimeRange('thisweek-vs-lastweek')"
          >
            周比
          </button>
        </div>
      </div>
      <div class="alert-stats-table-wrapper" v-loading="alertStatsData.loading">
        <div v-if="alertStatsData.items.length === 0" class="no-data">暂无数据</div>
        <el-table v-else :data="alertStatsData.items" :row-class-name="getAlertRowClass" style="width: 100%">
          <el-table-column type="index" label="排名" width="80" align="center">
            <template #default="{ $index }">
              <div :class="['rank-badge', { 'top-five': $index < 5 }]">{{ $index + 1 }}</div>
            </template>
          </el-table-column>
          <el-table-column prop="label" label="告警标签" min-width="150" show-overflow-tooltip />
          <el-table-column prop="totalAlertCnt" label="告警数" width="70" align="center" />
          <el-table-column prop="totalAlertEventCnt" label="事件数" width="70" align="center" />
          <el-table-column label="占比" width="80" align="center">
            <template #default="{ row }">
              <span>{{ (row.percentage || 0).toFixed(2) }}%</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="alertStatsData.timeRange === 'today-vs-yesterday' || alertStatsData.timeRange === 'thisweek-vs-lastweek'"
            label="变化率"
            width="120"
            align="center"
          >
            <template #default="{ row }">
              <span :class="['change-rate', row.alertEventCntChangeRate >= 0 ? 'positive' : 'negative']">
                {{ (row.alertEventCntChangeRate || 0) >= 0 ? '+' : '' }}{{ (row.alertEventCntChangeRate || 0).toFixed(2) }}%
              </span>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- 2. 告警通知统计 -->
    <div class="stat-card alert-card notification-card">
      <div class="card-header-with-tabs">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="#409eff" style="width: 20px; height: 20px; margin-right: 8px;">
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
      <div v-loading="notificationStatsData.loading" style="flex: 1; display: flex; flex-direction: column; overflow: hidden;">
        <div class="notification-summary">
          <div class="summary-item">
            <div class="summary-label">短信</div>
            <div class="summary-value" style="color: #409eff;">{{ notificationStatsData.totalSms }}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">邮件</div>
            <div class="summary-value" style="color: #67c23a;">{{ notificationStatsData.totalEmail }}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">语音</div>
            <div class="summary-value" style="color: #e6a23c;">{{ notificationStatsData.totalVoice }}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">总计</div>
            <div class="summary-value" style="color: #333;">{{ notificationStatsData.grandTotal }}</div>
          </div>
        </div>
        <div id="notificationChart" style="flex: 1; width: 100%; min-height: 0;"></div>
      </div>
    </div>

    <!-- 3. 业务线告警统计 -->
    <div class="stat-card alert-card channel-card">
      <div class="card-header-with-tabs">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="#67c23a" style="width: 20px; height: 20px; margin-right: 8px;">
            <path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/>
          </svg>
          业务线告警统计
        </div>
        <div class="time-range-tabs">
          <button
            :class="['tab-btn', { active: channelStatsData.timeRange === 'today' }]"
            @click="changeChannelTimeRange('today')"
          >
            今天
          </button>
          <button
            :class="['tab-btn', { active: channelStatsData.timeRange === 'yesterday' }]"
            @click="changeChannelTimeRange('yesterday')"
          >
            昨天
          </button>
          <button
            :class="['tab-btn', { active: channelStatsData.timeRange === 'last7days' }]"
            @click="changeChannelTimeRange('last7days')"
          >
            7天
          </button>
          <button
            :class="['tab-btn', { active: channelStatsData.timeRange === 'last30days' }]"
            @click="changeChannelTimeRange('last30days')"
          >
            30天
          </button>
        </div>
      </div>
      <div class="channel-chart-wrapper" v-loading="channelStatsData.loading">
        <div v-if="channelStatsData.items.length === 0" class="no-data">暂无数据</div>
        <div v-else id="channelChart" style="width: 100%; height: 100%;"></div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.alert-stats-section {
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

// 告警统计表格
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

// 前5名告警行的背景颜色
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

.change-rate {
  font-weight: 600;

  &.positive {
    color: #f56c6c;
  }

  &.negative {
    color: #67c23a;
  }
}

// 告警通知统计
.notification-card {
  .notification-summary {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 12px;
    padding: 10px 18px;
    border-bottom: 1px solid #f0f0f0;
    flex-shrink: 0;
  }

  .summary-item {
    text-align: center;
  }

  .summary-label {
    font-size: 11px;
    color: #999;
    margin-bottom: 4px;
  }

  .summary-value {
    font-size: 20px;
    font-weight: bold;
  }

  #notificationChart {
    padding: 5px 10px;
  }
}

// 协作空间统计
.channel-card {
  .channel-chart-wrapper {
    flex: 1;
    overflow: hidden;
    padding: 10px 18px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

// 响应式布局
@media (max-width: 1200px) {
  .alert-stats-section {
    grid-template-columns: 1fr;
  }
}
</style>
