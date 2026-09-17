<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { getBusinessDistribution, getDeploymentTrend } from '@/api/dashboard'
import { ElMessage } from 'element-plus'
import systemApi from '@/api/system'
import SpeakerIcon from '@/assets/image/喇叭.svg'
import NotificationIcon from '@/assets/image/通知.svg'

// 发布统计时间维度
const deployTimeRange = ref('week') // week, month, year

// 图表实例
let trendChart = null
let pieChart = null

// 加载状态
const loading = ref(false)

// 系统动态日志
const activityLogs = ref([])
const logsLoading = ref(false)

// 获取发布数据（根据时间维度）
const getDeploymentData = async (timeRange) => {
  try {
    loading.value = true
    const response = await getDeploymentTrend(timeRange)

    if (response.data && response.data.code === 200) {
      const result = response.data.data

      // 处理数据以适配图表格式
      const labels = []
      const serviceData = []
      const scriptData = []
      const quickDeployData = []

      // 处理日期标签格式
      result.data.forEach(item => {
        let label = ''
        if (timeRange === 'year') {
          // 年度视图：显示月份（如 "1月"）
          label = item.date.substring(5, 7) + '月'
          // 去掉月份前面的0
          label = label.replace(/^0/, '')
        } else {
          // 周/月视图：显示月-日（如 "12-18"）
          label = item.date.substring(5)
        }

        labels.push(label)
        serviceData.push(item.service_count)
        scriptData.push(item.script_count)
        quickDeployData.push(item.quick_deploy_count)
      })

      return {
        labels,
        service: serviceData,
        script: scriptData,
        quickDeploy: quickDeployData
      }
    } else {
      ElMessage.error('获取上线趋势数据失败')
      return null
    }
  } catch (error) {
    console.error('加载上线趋势数据失败:', error)
    ElMessage.error('加载上线趋势数据失败')
    return null
  } finally {
    loading.value = false
  }
}

// 初始化发布统计图
const initTrendChart = async () => {
  const chartDom = document.getElementById('deploymentTrendChart')
  if (!chartDom) return

  trendChart = echarts.init(chartDom)
  await updateTrendChart()
}

// 更新发布统计图
const updateTrendChart = async () => {
  if (!trendChart) return

  const data = await getDeploymentData(deployTimeRange.value)
  if (!data) return

  const option = {
    title: {
      text: '上线发布次数统计',
      left: 20,
      top: 10,
      textStyle: {
        fontSize: 16,
        fontWeight: 'normal',
        color: '#333'
      }
    },
    grid: {
      left: 60,
      right: 30,
      top: 60,
      bottom: 40
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985'
        }
      },
      formatter: (params) => {
        let result = params[0].name + '<br/>'
        params.forEach(item => {
          result += `${item.marker} ${item.seriesName}: ${item.value}次<br/>`
        })
        return result
      }
    },
    xAxis: {
      type: 'category',
      data: data.labels,
      axisLine: { lineStyle: { color: '#e0e0e0' } },
      axisTick: { show: false },
      axisLabel: { color: '#999' }
    },
    yAxis: {
      type: 'value',
      name: '发布次数',
      nameTextStyle: { color: '#999', fontSize: 12 },
      splitLine: { lineStyle: { color: '#f5f5f5' } },
      axisLabel: { color: '#999' }
    },
    legend: {
      data: ['生产上线', '脚本上线', '测试上线'],
      top: 15,
      right: 80
    },
    series: [
      {
        name: '生产上线',
        type: 'line',
        smooth: true,
        data: data.service,
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(255, 107, 107, 0.4)' },
              { offset: 1, color: 'rgba(255, 107, 107, 0.05)' }
            ]
          }
        },
        lineStyle: { color: '#ff6b6b', width: 2 },
        itemStyle: { color: '#ff6b6b' }
      },
      {
        name: '脚本上线',
        type: 'line',
        smooth: true,
        data: data.script,
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(78, 201, 176, 0.4)' },
              { offset: 1, color: 'rgba(78, 201, 176, 0.05)' }
            ]
          }
        },
        lineStyle: { color: '#4ecdc4', width: 2 },
        itemStyle: { color: '#4ecdc4' }
      },
      {
        name: '测试上线',
        type: 'line',
        smooth: true,
        data: data.quickDeploy,
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(255, 193, 7, 0.4)' },
              { offset: 1, color: 'rgba(255, 193, 7, 0.05)' }
            ]
          }
        },
        lineStyle: { color: '#ffc107', width: 2 },
        itemStyle: { color: '#ffc107' }
      }
    ]
  }
  trendChart.setOption(option)
}

// 切换时间维度
const changeTimeRange = async (range) => {
  deployTimeRange.value = range
  await updateTrendChart()
}

// 初始化环形图
const initPieChart = async () => {
  const chartDom = document.getElementById('businessPieChart')
  if (!chartDom) return

  pieChart = echarts.init(chartDom)

  // 加载业务分布数据
  let businessData = []
  try {
    const response = await getBusinessDistribution()
    if (response.data && response.data.code === 200) {
      const data = response.data.data
      const colors = ['#5dade2', '#f8b739', '#48c9b0', '#9b59b6', '#ec7063', '#ff6b6b', '#4ecdc4', '#45b7d1']
      businessData = data.businessLines.map((line, index) => ({
        value: line.serviceCount,
        name: line.name,
        itemStyle: { color: colors[index % colors.length] }
      }))
    }
  } catch (error) {
    console.error('加载业务分布数据失败:', error)
    // 使用默认数据
    businessData = [
      { value: 10, name: '暂无数据', itemStyle: { color: '#e0e0e0' } }
    ]
  }

  const option = {
    title: {
      text: '业务组应用分布',
      left: 20,
      top: 10,
      textStyle: {
        fontSize: 16,
        fontWeight: 'normal',
        color: '#333'
      }
    },
    tooltip: {
      trigger: 'item',
      formatter: '{b}<br/>应用数: {c}<br/>占比: {d}%'
    },
    legend: {
      orient: 'vertical',
      right: 20,
      bottom: 20,
      itemWidth: 10,
      itemHeight: 10,
      itemGap: 8,
      textStyle: { fontSize: 11, color: '#666' }
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '55%'],
        center: ['50%', '45%'],
        avoidLabelOverlap: false,
        label: { show: false },
        labelLine: { show: false },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.3)'
          }
        },
        data: businessData
      }
    ]
  }
  pieChart.setOption(option)
}

// 获取系统动态日志
const getActivityLogs = async () => {
  try {
    logsLoading.value = true
    const response = await systemApi.getActivityLogs({ limit: 10 })

    if (response.data && response.data.code === 200) {
      activityLogs.value = response.data.data || []
    } else {
      ElMessage.error('获取系统动态日志失败')
    }
  } catch (error) {
    console.error('加载系统动态日志失败:', error)
    ElMessage.error('加载系统动态日志失败')
  } finally {
    logsLoading.value = false
  }
}

// 格式化时间为 MM-DD HH:mm:ss
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  // 后端返回格式: "2025-12-27 10:34:54"
  // 提取 MM-DD HH:mm:ss
  const parts = timeStr.split(' ')
  if (parts.length !== 2) return timeStr

  const datePart = parts[0].split('-')
  if (datePart.length !== 3) return timeStr

  return `${datePart[1]}-${datePart[2]} ${parts[1]}`
}

// 获取活动类型标签
const getActivityTypeLabel = (type) => {
  const typeMap = {
    1: '密钥',
    2: '域名',
    3: 'ecs巡检',
    4: '定时任务',
    5: '其他'
  }
  return typeMap[type] || '其他'
}

// 获取活动类型颜色
const getActivityTypeColor = (type) => {
  const colorMap = {
    1: '#FF6B6B',      // 密钥 - 红色
    2: '#4ECDC4',      // 域名 - 青绿色
    3: '#45B7D1',      // ecs巡检 - 蓝色
    4: '#FFA502',      // 定时任务 - 橙色
    5: '#95A5A6'       // 其他 - 灰色
  }
  return colorMap[type] || '#95A5A6'
}

// 窗口大小改变时重绘图表
const handleResize = () => {
  trendChart?.resize()
  pieChart?.resize()
}

// 生命周期
onMounted(() => {
  setTimeout(async () => {
    initTrendChart()
    await initPieChart()
    await getActivityLogs()
  }, 100)
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  pieChart?.dispose()
})
</script>

<template>
  <div class="deployment-stats-section">
    <!-- 系统动态日志 -->
    <div class="chart-card activity-log-card">
      <div class="card-title">
        <img :src="SpeakerIcon" class="title-icon" alt="通知" />
        <span>系统动态</span>
      </div>
      <div v-loading="logsLoading" class="activity-log-content">
        <div v-if="activityLogs.length === 0" class="empty-log">
          暂无系统动态
        </div>
        <div v-else class="log-list">
          <div
            v-for="(log, index) in activityLogs"
            :key="index"
            class="log-item"
          >
            <div class="log-text">
              <img :src="NotificationIcon" class="log-icon" alt="通知" />
              <span class="log-time">{{ formatTime(log.createTime) }}</span>：<span class="log-type" :style="{ color: getActivityTypeColor(log.activityType) }">[{{ getActivityTypeLabel(log.activityType) }}]</span> {{ log.title }}{{ log.summary }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 发布统计图 -->
    <div class="chart-card large">
      <div class="chart-header">
        <div class="time-range-tabs">
          <button
            :class="['tab-btn', { active: deployTimeRange === 'week' }]"
            @click="changeTimeRange('week')"
          >
            周
          </button>
          <button
            :class="['tab-btn', { active: deployTimeRange === 'month' }]"
            @click="changeTimeRange('month')"
          >
            月
          </button>
          <button
            :class="['tab-btn', { active: deployTimeRange === 'year' }]"
            @click="changeTimeRange('year')"
          >
            年
          </button>
        </div>
      </div>
      <div id="deploymentTrendChart" style="width: 100%; height: calc(100% - 40px);"></div>
    </div>

    <!-- 环形图 -->
    <div class="chart-card">
      <div id="businessPieChart" style="width: 100%; height: 100%;"></div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.deployment-stats-section {
  display: grid;
  grid-template-columns: 0.8fr 2fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.chart-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  height: 400px;
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
    transform: translateY(-2px);
  }

  &.large {
    height: 400px;
  }
}

.chart-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
  padding-bottom: 10px;
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

// 系统动态日志样式
.activity-log-card {
  display: flex;
  flex-direction: column;

  .card-title {
    font-size: 16px;
    font-weight: 500;
    color: #333;
    margin-bottom: 15px;
    display: flex;
    align-items: center;
    gap: 8px;

    .title-icon {
      width: 20px;
      height: 20px;
    }
  }

  .activity-log-content {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .empty-log {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: #999;
    font-size: 14px;
  }

  .log-list {
    flex: 1;
    overflow-y: auto;
    padding-right: 5px;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-thumb {
      background: #ddd;
      border-radius: 3px;

      &:hover {
        background: #ccc;
      }
    }

    &::-webkit-scrollbar-track {
      background: transparent;
    }
  }

  .log-item {
    padding: 12px 0;
    border-bottom: 1px dashed #4ecdc4;

    &:last-child {
      border-bottom: none;
    }

    .log-text {
      font-size: 13px;
      color: #666;
      line-height: 1.6;
      word-break: break-all;
      display: flex;
      align-items: flex-start;
      gap: 0;

      .log-icon {
        width: 16px;
        height: 16px;
        margin-top: 2px;
        margin-right: 6px;
        flex-shrink: 0;
      }

      .log-time {
        color: #999;
        font-weight: 500;
      }

      .log-type {
        font-weight: 600;
        white-space: nowrap;
      }
    }
  }
}

// 响应式设计
@media (max-width: 1400px) {
  .deployment-stats-section {
    grid-template-columns: 1fr;
  }
}
</style>
