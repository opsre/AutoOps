<script setup>
import { reactive, onMounted, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import dashboardApi from '@/api/dashboard'

// 云账单数据
const cloudBillingData = reactive({
  activeProvider: 'tencent', // 当前激活的云厂商: tencent, aliyun, aws
  tencent: {
    loading: false,
    month: '',
    sourceId: null, // 告警源ID
    keyId: null, // 密钥ID（从告警源读取）
    data: {
      month: '',
      costs: [],
      total_cost: 0
    }
  },
  aliyun: {
    loading: false,
    month: '',
    sourceId: null,
    keyId: null,
    billingType: 'aliyun-tag', // aliyun-product or aliyun-tag，默认按标签
    forceRefresh: false, // 强制刷新标志
    data: {
      month: '',
      costs: [],
      total_cost: 0
    }
  },
  aws: {
    loading: false,
    month: '',
    sourceId: null,
    keyId: null,
    data: {
      month: '',
      costs: [],
      total_cost: 0,
      total_cost_usd: 0,
      exchange_rate: 0
    }
  },
  alertSources: [], // 云账号类型的告警源列表
  // 12个月趋势数据
  trends: {
    loading: false,
    data: []
  },
  // 上月费用汇总
  lastMonthSummary: {
    loading: false,
    currentMonth: '',
    previousMonth: '',
    summary: []
  }
})

// 图表实例
let trendBillingChart = null

// 获取云账单告警源列表
const getCloudBillingAlertSources = async () => {
  try {
    const res = await dashboardApi.getAlertSourceList({ pageNum: 1, pageSize: 100, status: 1 })
    if (res.data && (res.data.code === 0 || res.data.code === 200)) {
      const allSources = res.data.data.list || []
      // 只保留云账号类型的告警源 (type=2阿里云, type=3腾讯云, type=4AWS)
      cloudBillingData.alertSources = allSources.filter(source => [2, 3, 4].includes(source.type))

      // 为每个云厂商设置默认告警源
      const tencentSources = cloudBillingData.alertSources.filter(s => s.type === 3)
      const aliyunSources = cloudBillingData.alertSources.filter(s => s.type === 2)
      const awsSources = cloudBillingData.alertSources.filter(s => s.type === 4)

      if (tencentSources.length > 0) {
        cloudBillingData.tencent.sourceId = tencentSources[0].id
        cloudBillingData.tencent.keyId = tencentSources[0].keyId
      }
      if (aliyunSources.length > 0) {
        cloudBillingData.aliyun.sourceId = aliyunSources[0].id
        cloudBillingData.aliyun.keyId = aliyunSources[0].keyId
      }
      if (awsSources.length > 0) {
        cloudBillingData.aws.sourceId = awsSources[0].id
        cloudBillingData.aws.keyId = awsSources[0].keyId
      }
    }
  } catch (error) {
    console.error('获取云账单告警源列表失败:', error)
  }
}

// 查询腾讯云账单
const queryTencentBilling = async () => {
  if (!cloudBillingData.tencent.keyId) {
    ElMessage.warning('未找到腾讯云数据源，请先配置腾讯云告警源')
    return
  }

  cloudBillingData.tencent.loading = true
  try {
    const res = await dashboardApi.getTencentBilling({
      keyId: cloudBillingData.tencent.keyId,
      month: cloudBillingData.tencent.month || undefined
    })

    if (res.data && res.data.code === 200) {
      const responseData = res.data.data || {}
      cloudBillingData.tencent.data = {
        month: responseData.month || '',
        costs: responseData.costs || [],
        total_cost: responseData.total_cost || 0
      }
      ElMessage.success('查询成功')
    } else {
      ElMessage.error(res.data?.message || '查询失败')
    }
  } catch (error) {
    console.error('查询腾讯云账单失败:', error)
    ElMessage.error('查询失败: ' + (error.message || '未知错误'))
  } finally {
    cloudBillingData.tencent.loading = false
  }
}

// 查询阿里云账单
const queryAliyunBilling = async (forceRefresh = false) => {
  if (!cloudBillingData.aliyun.keyId) {
    ElMessage.warning('未找到阿里云数据源，请先配置阿里云告警源')
    return
  }

  cloudBillingData.aliyun.loading = true
  try {
    const params = {
      keyId: cloudBillingData.aliyun.keyId,
      month: cloudBillingData.aliyun.month || undefined,
      forceRefresh: forceRefresh
    }

    const apiFunc = cloudBillingData.aliyun.billingType === 'aliyun-product'
      ? dashboardApi.getAliyunProductBilling
      : dashboardApi.getAliyunTagBilling

    const res = await apiFunc(params)

    if (res.data && res.data.code === 200) {
      const responseData = res.data.data || {}
      cloudBillingData.aliyun.data = {
        month: responseData.month || '',
        costs: responseData.costs || [],
        total_cost: responseData.total_cost || 0
      }
      ElMessage.success('查询成功')
    } else {
      ElMessage.error(res.data?.message || '查询失败')
    }
  } catch (error) {
    console.error('查询阿里云账单失败:', error)
    ElMessage.error('查询失败: ' + (error.message || '未知错误'))
  } finally {
    cloudBillingData.aliyun.loading = false
  }
}

// 查询AWS账单
const queryAWSBilling = async () => {
  if (!cloudBillingData.aws.keyId) {
    ElMessage.warning('未找到AWS数据源，请先配置AWS告警源')
    return
  }

  cloudBillingData.aws.loading = true
  try {
    const res = await dashboardApi.getAWSBilling({
      keyId: cloudBillingData.aws.keyId,
      month: cloudBillingData.aws.month || undefined
    })

    if (res.data && res.data.code === 200) {
      const responseData = res.data.data || {}
      cloudBillingData.aws.data = {
        month: responseData.month || '',
        costs: responseData.costs || [],
        total_cost: responseData.total_cost || 0,
        total_cost_usd: responseData.total_cost_usd || 0,
        exchange_rate: responseData.exchange_rate || 0
      }
      ElMessage.success('查询成功')
    } else {
      ElMessage.error(res.data?.message || '查询失败')
    }
  } catch (error) {
    console.error('查询AWS账单失败:', error)
    ElMessage.error('查询失败: ' + (error.message || '未知错误'))
  } finally {
    cloudBillingData.aws.loading = false
  }
}

// 切换云厂商
const switchCloudProvider = (provider) => {
  cloudBillingData.activeProvider = provider
  // 切换时自动加载该厂商上月账单
  loadProviderLastMonthBilling(provider)
}

// 加载指定云厂商的上月账单
const loadProviderLastMonthBilling = async (provider) => {
  const providerMap = {
    tencent: {
      keyId: cloudBillingData.tencent.keyId,
      apiFunc: dashboardApi.getTencentBilling
    },
    aliyun: {
      keyId: cloudBillingData.aliyun.keyId,
      apiFunc: cloudBillingData.aliyun.billingType === 'aliyun-product'
        ? dashboardApi.getAliyunProductBilling
        : dashboardApi.getAliyunTagBilling
    },
    aws: {
      keyId: cloudBillingData.aws.keyId,
      apiFunc: dashboardApi.getAWSBilling
    }
  }

  const config = providerMap[provider]
  if (!config || !config.keyId) return

  cloudBillingData[provider].loading = true
  try {
    const res = await config.apiFunc({
      keyId: config.keyId,
      month: undefined // 不传month则查询上月
    })

    if (res.data && res.data.code === 200) {
      const responseData = res.data.data || {}
      if (provider === 'aws') {
        cloudBillingData[provider].data = {
          month: responseData.month || '',
          costs: responseData.costs || [],
          total_cost: responseData.total_cost || 0,
          total_cost_usd: responseData.total_cost_usd || 0,
          exchange_rate: responseData.exchange_rate || 0
        }
      } else {
        cloudBillingData[provider].data = {
          month: responseData.month || '',
          costs: responseData.costs || [],
          total_cost: responseData.total_cost || 0
        }
      }
    }
  } catch (error) {
    console.error(`查询${provider}账单失败:`, error)
  } finally {
    cloudBillingData[provider].loading = false
  }
}

// 加载12个月费用趋势
const load12MonthsTrend = async () => {
  if (!cloudBillingData.tencent.keyId && !cloudBillingData.aliyun.keyId && !cloudBillingData.aws.keyId) {
    return
  }

  cloudBillingData.trends.loading = true
  try {
    const res = await dashboardApi.getLast12MonthsTrend({
      tencentKeyId: cloudBillingData.tencent.keyId || undefined,
      aliyunKeyId: cloudBillingData.aliyun.keyId || undefined,
      awsKeyId: cloudBillingData.aws.keyId || undefined
    })

    if (res.data && res.data.code === 200) {
      const responseData = res.data.data || {}
      cloudBillingData.trends.data = responseData.data || []
      // 数据加载后，等待DOM更新，然后初始化并更新图表
      nextTick(() => {
        if (!trendBillingChart) {
          initTrendBillingChart()
        } else {
          updateTrendBillingChart()
        }
      })
    }
  } catch (error) {
    console.error('加载12个月费用趋势失败:', error)
  } finally {
    cloudBillingData.trends.loading = false
  }
}

// 加载上月费用汇总
const loadLastMonthSummary = async () => {
  if (!cloudBillingData.tencent.keyId && !cloudBillingData.aliyun.keyId && !cloudBillingData.aws.keyId) {
    return
  }

  cloudBillingData.lastMonthSummary.loading = true
  try {
    const res = await dashboardApi.getLastMonthSummary({
      tencentKeyId: cloudBillingData.tencent.keyId || undefined,
      aliyunKeyId: cloudBillingData.aliyun.keyId || undefined,
      awsKeyId: cloudBillingData.aws.keyId || undefined
    })

    if (res.data && res.data.code === 200) {
      const data = res.data.data
      cloudBillingData.lastMonthSummary.currentMonth = data.current_month || ''
      cloudBillingData.lastMonthSummary.previousMonth = data.previous_month || ''
      cloudBillingData.lastMonthSummary.summary = data.summary || []
    }
  } catch (error) {
    console.error('加载上月费用汇总失败:', error)
  } finally {
    cloudBillingData.lastMonthSummary.loading = false
  }
}

// 格式化金额
const formatMoney = (amount) => {
  if (!amount && amount !== 0) return '0.00'
  return Number(amount).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

// 初始化云账单12个月趋势图表
const initTrendBillingChart = () => {
  const chartDom = document.getElementById('trendBillingChart')
  if (!chartDom) return

  trendBillingChart = echarts.init(chartDom)
  updateTrendBillingChart()
}

// 更新云账单12个月趋势图表
const updateTrendBillingChart = () => {
  if (!trendBillingChart) return
  if (cloudBillingData.trends.data.length === 0) return

  const rawData = cloudBillingData.trends.data

  // 提取月份（从第一个云厂商的trends中获取）
  const months = rawData[0]?.trends?.map(item => item.month) || []

  // 提取各云厂商的数据
  const tencentProvider = rawData.find(item => item.cloud_provider === '腾讯云')
  const aliyunProvider = rawData.find(item => item.cloud_provider === '阿里云')
  const awsProvider = rawData.find(item => item.cloud_provider === 'AWS')

  const tencentData = tencentProvider?.trends?.map(item => item.amount || 0) || []
  const aliyunData = aliyunProvider?.trends?.map(item => item.amount || 0) || []
  const awsData = awsProvider?.trends?.map(item => item.amount || 0) || []

  const series = []

  if (cloudBillingData.tencent.sourceId && tencentData.length > 0) {
    series.push({
      name: '腾讯云',
      type: 'line',
      smooth: true,
      data: tencentData,
      itemStyle: { color: '#00a4ff' },
      lineStyle: { color: '#00a4ff', width: 2 },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(0, 164, 255, 0.3)' },
            { offset: 1, color: 'rgba(0, 164, 255, 0.05)' }
          ]
        }
      }
    })
  }

  if (cloudBillingData.aliyun.sourceId && aliyunData.length > 0) {
    series.push({
      name: '阿里云',
      type: 'line',
      smooth: true,
      data: aliyunData,
      itemStyle: { color: '#ff6a00' },
      lineStyle: { color: '#ff6a00', width: 2 },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(255, 106, 0, 0.3)' },
            { offset: 1, color: 'rgba(255, 106, 0, 0.05)' }
          ]
        }
      }
    })
  }

  if (cloudBillingData.aws.sourceId && awsData.length > 0) {
    series.push({
      name: 'AWS',
      type: 'line',
      smooth: true,
      data: awsData,
      itemStyle: { color: '#232f3e' },
      lineStyle: { color: '#232f3e', width: 2 },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(35, 47, 62, 0.3)' },
            { offset: 1, color: 'rgba(35, 47, 62, 0.05)' }
          ]
        }
      }
    })
  }

  const option = {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: { backgroundColor: '#6a7985' }
      },
      formatter: (params) => {
        let result = params[0].name + '<br/>'
        params.forEach(item => {
          result += `${item.marker} ${item.seriesName}: ¥${formatMoney(item.value)}<br/>`
        })
        return result
      }
    },
    legend: {
      data: series.map(s => s.name),
      textStyle: { color: '#666' },
      top: 10
    },
    grid: {
      left: 60,
      right: 30,
      top: 50,
      bottom: 30
    },
    xAxis: {
      type: 'category',
      data: months,
      axisLabel: { color: '#666', fontSize: 11 },
      axisLine: { lineStyle: { color: '#e0e0e0' } },
      axisTick: { show: false }
    },
    yAxis: {
      type: 'value',
      name: '费用(CNY)',
      nameTextStyle: { color: '#999', fontSize: 12 },
      splitLine: { lineStyle: { color: '#f5f5f5' } },
      axisLabel: {
        color: '#666',
        formatter: (value) => '¥' + (value / 1000).toFixed(0) + 'k'
      }
    },
    series: series
  }

  trendBillingChart.setOption(option)
}

// 窗口大小改变时重绘图表
const handleResize = () => {
  trendBillingChart?.resize()
}

// 初始化加载
const initLoad = async () => {
  await getCloudBillingAlertSources()

  if (cloudBillingData.tencent.keyId || cloudBillingData.aliyun.keyId || cloudBillingData.aws.keyId) {
    // 加载默认激活的腾讯云上月账单
    if (cloudBillingData.tencent.keyId) {
      loadProviderLastMonthBilling('tencent')
    } else if (cloudBillingData.aliyun.keyId) {
      cloudBillingData.activeProvider = 'aliyun'
      loadProviderLastMonthBilling('aliyun')
    } else if (cloudBillingData.aws.keyId) {
      cloudBillingData.activeProvider = 'aws'
      loadProviderLastMonthBilling('aws')
    }

    // 加载12个月趋势和上月汇总
    load12MonthsTrend()
    loadLastMonthSummary()
  }

  setTimeout(() => {
    initTrendBillingChart()
    nextTick(() => {
      if (cloudBillingData.trends.data.length > 0) {
        updateTrendBillingChart()
      }
    })
  }, 100)
}

onMounted(() => {
  initLoad()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendBillingChart?.dispose()
})
</script>

<template>
  <div v-if="cloudBillingData.alertSources.length > 0" class="cloud-billing-section-unified">
    <!-- 卡片1: 云账单明细 -->
    <div class="stat-card billing-card-unified">
      <div class="billing-header-oneline">
        <div class="billing-header-left">
          <img
            v-if="cloudBillingData.activeProvider === 'tencent'"
            src="@/assets/image/tengxun.png"
            alt="腾讯云"
            class="provider-icon"
          />
          <img
            v-else-if="cloudBillingData.activeProvider === 'aliyun'"
            src="@/assets/image/aliyun.png"
            alt="阿里云"
            class="provider-icon"
          />
          <img
            v-else
            src="@/assets/image/aws.svg"
            alt="AWS"
            class="provider-icon"
          />
          <span class="billing-title">{{ cloudBillingData.activeProvider === 'tencent' ? '腾讯云' : cloudBillingData.activeProvider === 'aliyun' ? '阿里云' : 'AWS' }}账单明细</span>
          <span class="billing-month">{{ cloudBillingData[cloudBillingData.activeProvider].data.month }}</span>
          <span class="billing-amount">¥ {{ formatMoney(cloudBillingData[cloudBillingData.activeProvider].data.total_cost) }}</span>
        </div>
        <div class="provider-switcher">
          <button
            v-if="cloudBillingData.tencent.sourceId"
            :class="['provider-btn', { active: cloudBillingData.activeProvider === 'tencent' }]"
            @click="switchCloudProvider('tencent')"
            title="腾讯云"
          >
            <img src="@/assets/image/tengxun.png" alt="腾讯云" style="width: 20px; height: 20px;" />
          </button>
          <button
            v-if="cloudBillingData.aliyun.sourceId"
            :class="['provider-btn', { active: cloudBillingData.activeProvider === 'aliyun' }]"
            @click="switchCloudProvider('aliyun')"
            title="阿里云"
          >
            <img src="@/assets/image/aliyun.png" alt="阿里云" style="width: 20px; height: 20px;" />
          </button>
          <button
            v-if="cloudBillingData.aws.sourceId"
            :class="['provider-btn', { active: cloudBillingData.activeProvider === 'aws' }]"
            @click="switchCloudProvider('aws')"
            title="AWS"
          >
            <img src="@/assets/image/aws.svg" alt="AWS" style="width: 20px; height: 20px;" />
          </button>
        </div>
      </div>

      <!-- 账单结果 -->
      <div
        v-loading="cloudBillingData[cloudBillingData.activeProvider].loading"
        class="billing-result-compact"
      >
        <div v-if="cloudBillingData[cloudBillingData.activeProvider].data.costs && cloudBillingData[cloudBillingData.activeProvider].data.costs.length > 0">
          <div class="billing-table-compact">
            <!-- 腾讯云和阿里云表格 -->
            <el-table
              v-if="cloudBillingData.activeProvider !== 'aws'"
              :data="cloudBillingData[cloudBillingData.activeProvider].data.costs"
              style="width: 100%"
              max-height="300"
              size="small"
            >
              <el-table-column type="index" label="#" width="50" align="center" />
              <el-table-column prop="business_line" :label="cloudBillingData.activeProvider === 'aliyun' ? '业务线/产品' : '业务线'" min-width="150" show-overflow-tooltip />
              <el-table-column label="费用(CNY)" width="130" align="right">
                <template #default="{ row }">
                  <span class="money-text">¥ {{ formatMoney(row.amount) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="占比" width="80" align="center">
                <template #default="{ row }">
                  <span class="percentage-text">{{ ((row.amount / cloudBillingData[cloudBillingData.activeProvider].data.total_cost) * 100).toFixed(1) }}%</span>
                </template>
              </el-table-column>
            </el-table>

            <!-- AWS表格 -->
            <el-table
              v-else
              :data="cloudBillingData.aws.data.costs"
              style="width: 100%"
              max-height="300"
              size="small"
            >
              <el-table-column type="index" label="#" width="50" align="center" />
              <el-table-column prop="business_line" label="服务" min-width="150" show-overflow-tooltip />
              <el-table-column label="费用(CNY)" width="110" align="right">
                <template #default="{ row }">
                  <span class="money-text">¥ {{ formatMoney(row.amount) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="费用(USD)" width="110" align="right">
                <template #default="{ row }">
                  <span class="money-text-usd">$ {{ formatMoney(row.original_amount) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="占比" width="70" align="center">
                <template #default="{ row }">
                  <span class="percentage-text">{{ ((row.amount / cloudBillingData.aws.data.total_cost) * 100).toFixed(1) }}%</span>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
        <div v-else-if="cloudBillingData[cloudBillingData.activeProvider].data.total_cost > 0" class="billing-summary-only">
          <div class="summary-title">账单汇总</div>
          <div class="summary-content">
            <div class="summary-row" v-if="cloudBillingData.activeProvider === 'aws'">
              <span class="summary-label">费用(CNY):</span>
              <span class="summary-amount">¥ {{ formatMoney(cloudBillingData.aws.data.total_cost) }}</span>
            </div>
            <div class="summary-row" v-if="cloudBillingData.activeProvider === 'aws'">
              <span class="summary-label">费用(USD):</span>
              <span class="summary-amount-usd">$ {{ formatMoney(cloudBillingData.aws.data.total_cost_usd) }}</span>
            </div>
            <div class="summary-row" v-if="cloudBillingData.activeProvider === 'aws'">
              <span class="summary-label">汇率:</span>
              <span class="summary-rate">{{ cloudBillingData.aws.data.exchange_rate }}</span>
            </div>
            <div class="summary-row" v-if="cloudBillingData.activeProvider !== 'aws'">
              <span class="summary-label">总费用:</span>
              <span class="summary-amount">¥ {{ formatMoney(cloudBillingData[cloudBillingData.activeProvider].data.total_cost) }}</span>
            </div>
          </div>
          <div class="no-detail-tip">暂无明细数据，仅显示总金额</div>
        </div>
        <div v-else class="billing-empty-compact">
          <div class="empty-text">暂无数据</div>
        </div>
      </div>
    </div>

    <!-- 卡片2: 12个月费用趋势 -->
    <div class="stat-card trend-card">
      <div class="card-header-with-tabs">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="#409eff" style="width: 20px; height: 20px; margin-right: 8px;">
            <path d="M3.5 18.49l6-6.01 4 4L22 6.92l-1.41-1.41-7.09 7.97-4-4L2 16.99z"/>
          </svg>
          近12个月费用趋势
        </div>
      </div>
      <div v-loading="cloudBillingData.trends.loading" style="padding: 8px 16px;">
        <div v-if="cloudBillingData.trends.data.length > 0" id="trendBillingChart" style="width: 100%; height: 180px;"></div>
        <div v-else class="billing-empty-compact">
          <div class="empty-text">暂无趋势数据</div>
        </div>
      </div>
    </div>

    <!-- 卡片3: 上月费用汇总 -->
    <div class="stat-card summary-card">
      <div class="card-header-with-tabs">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="#67c23a" style="width: 20px; height: 20px; margin-right: 8px;">
            <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM9 17H7v-7h2v7zm4 0h-2V7h2v10zm4 0h-2v-4h2v4z"/>
          </svg>
          上月费用汇总
        </div>
      </div>
      <div v-loading="cloudBillingData.lastMonthSummary.loading" style="padding: 8px 16px; flex: 1; overflow: auto;">
        <div v-if="cloudBillingData.lastMonthSummary.summary.length > 0">
          <div class="summary-header">
            <span class="summary-month-label">{{ cloudBillingData.lastMonthSummary.currentMonth }}</span>
            <span class="vs-label">vs</span>
            <span class="summary-month-label">{{ cloudBillingData.lastMonthSummary.previousMonth }}</span>
          </div>
          <div class="summary-list" style="display: flex !important; gap: 8px; flex-wrap: wrap;">
            <div
              v-for="item in cloudBillingData.lastMonthSummary.summary"
              :key="item.cloud_provider"
              class="summary-item-card"
              style="flex: 1; min-width: 0; padding: 6px; margin-bottom: 0; display: flex; flex-direction: column; gap: 4px;"
            >
              <div class="provider-header">
                <svg v-if="item.cloud_provider === '腾讯云'" viewBox="0 0 24 24" fill="#00a4ff" style="width: 16px; height: 16px;">
                  <circle cx="12" cy="12" r="8"/>
                </svg>
                <svg v-else-if="item.cloud_provider === '阿里云'" viewBox="0 0 24 24" fill="#ff6a00" style="width: 16px; height: 16px;">
                  <path d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="#ff9900" style="width: 16px; height: 16px;">
                  <path d="M19.35 10.04C18.67 6.59 15.64 4 12 4 9.11 4 6.6 5.64 5.35 8.04 2.34 8.36 0 10.91 0 14c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96z"/>
                </svg>
                <span class="provider-name">
                  {{ item.cloud_provider }}
                </span>
              </div>
              <div class="summary-cost">
                <div class="cost-row">
                  <span class="cost-label">当月</span>
                  <span class="cost-value">¥ {{ formatMoney(item.current_month) }}</span>
                </div>
                <div class="cost-row">
                  <span class="cost-label">上月</span>
                  <span class="cost-value">¥ {{ formatMoney(item.previous_month) }}</span>
                </div>
              </div>
              <div class="summary-change">
                <span :class="['change-badge', (item.month_on_month || 0) >= 0 ? 'increase' : 'decrease']">
                  {{ (item.month_on_month || 0) >= 0 ? '+' : '' }}{{ (item.month_on_month || 0).toFixed(2) }}%
                </span>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="billing-empty-compact">
          <div class="empty-text">暂无汇总数据</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.cloud-billing-section-unified {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 20px;

  .stat-card {
    height: 240px !important;
    min-height: 240px !important;
    max-height: 240px !important;
    background: white;
    border-radius: 8px;
    padding: 0;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    transition: all 0.3s ease;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    &:hover {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
      transform: translateY(-2px);
    }
  }
}

.billing-card-unified {
  .billing-header-oneline {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 16px;
    background: linear-gradient(135deg, #ffd166, #f9c74f);
    border-bottom: 1px solid #f0f0f0;
    flex-shrink: 0;
  }

  .billing-header-left {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: nowrap;
  }

  .provider-icon {
    width: 20px;
    height: 20px;
    object-fit: contain;
    flex-shrink: 0;
  }

  .billing-title {
    font-size: 12px;
    font-weight: 600;
    color: #333;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .billing-month {
    font-size: 11px;
    color: #666;
    padding: 2px 6px;
    background: rgba(255, 255, 255, 0.5);
    border-radius: 4px;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .billing-amount {
    font-size: 13px;
    font-weight: bold;
    color: #d63031;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .billing-result-compact {
    flex: 1;
    overflow: auto;
    padding: 8px 16px;
  }
}

.provider-switcher {
  display: flex;
  gap: 8px;
  background: #f5f7fa;
  padding: 4px;
  border-radius: 6px;
}

.provider-btn {
  padding: 8px 12px;
  border: none;
  background: transparent;
  color: #666;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;

  &:hover {
    color: #333;
    background: rgba(255, 255, 255, 0.5);
  }

  &.active {
    background: white;
    color: #2196f3;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }
}

.card-header-with-tabs {
  padding: 12px 16px;
  flex-shrink: 0;
  border-bottom: 1px solid #f0f0f0;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  display: flex;
  align-items: center;
}

.billing-table-compact {
  :deep(.el-table) {
    border-radius: 8px;
    overflow: hidden;
    font-size: 12px;

    .el-table__header th {
      background: #f8f9fa;
      font-weight: 600;
      color: #333;
      font-size: 12px;
      padding: 8px 0;
    }

    .el-table__body {
      font-size: 12px;
    }

    .el-table__row:hover {
      background-color: #f5f7fa;
    }
  }
}

.money-text {
  font-size: 13px;
  font-weight: 600;
  color: #f56c6c;
}

.money-text-usd {
  font-size: 12px;
  font-weight: 600;
  color: #67c23a;
}

.percentage-text {
  font-size: 12px;
  font-weight: 600;
  color: #409eff;
}

.billing-empty-compact {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.empty-text {
  font-size: 13px;
  color: #999;
}

.trend-card {
  #trendBillingChart {
    height: 180px !important;
  }
}

.summary-card {
  .summary-header {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 6px 0;
    margin-bottom: 6px;
    font-size: 12px;
  }

  .summary-month-label {
    font-size: 14px;
    font-weight: 600;
    color: #333;
  }

  .vs-label {
    font-size: 12px;
    color: #999;
    padding: 2px 8px;
    background: #f5f5f5;
    border-radius: 4px;
  }

  .summary-list {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .summary-item-card {
    flex: 1;
    min-width: 0;
    padding: 6px;
    margin-bottom: 0;
    display: flex;
    flex-direction: column;
    gap: 4px;
    background: #fafafa;
    border-radius: 8px;
    border-left: 3px solid #409eff;
    transition: all 0.3s ease;

    &:hover {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
      transform: translateX(2px);
    }
  }

  .provider-header {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-bottom: 4px;

    svg {
      width: 14px !important;
      height: 14px !important;
    }
  }

  .provider-name {
    font-size: 12px;
    font-weight: 600;
  }

  .summary-cost {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .cost-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .cost-label {
    font-size: 11px;
    color: #999;
  }

  .cost-value {
    font-size: 11px;
    font-weight: 600;
    color: #333;
  }

  .summary-change {
    text-align: center;
    margin-top: 2px;
  }

  .change-badge {
    font-size: 11px;
    padding: 2px 4px;
    display: inline-block;
    border-radius: 12px;
    font-weight: 600;

    &.increase {
      background: #fef0f0;
      color: #f56c6c;
    }

    &.decrease {
      background: #f0f9ff;
      color: #67c23a;
    }
  }
}

@media (max-width: 1400px) {
  .cloud-billing-section-unified {
    grid-template-columns: 1fr;
  }
}
</style>
