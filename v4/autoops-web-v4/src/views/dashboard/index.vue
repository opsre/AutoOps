<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboardStats } from '@/api/dashboard'
import { GetAllTools, CreateTool, UpdateTool, DeleteTool as DeleteToolAPI, UploadIcon } from '@/api/tool'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as echarts from 'echarts'
import dashboardApi from '@/api/dashboard'
import monitorApi from '@/api/monitor'

// 导入子组件
import AlertStats from './components/AlertStats.vue'
import ZabbixAlertStats from './components/ZabbixAlertStats.vue'
import BillingStats from './components/BillingStats.vue'
import DeploymentStats from './components/DeploymentStats.vue'

const router = useRouter()

// 响应式数据
const loading = ref(true)
const editDialogVisible = ref(false)
const editingTool = ref(null)
const editingIndex = ref(-1)

// 告警源相关
const alertSourceList = ref([])
const selectedAlertSource = ref(null)

// 统计数据
const stats = reactive({
  assets: {
    title: '资产详情',
    items: [
      { label: '主机总数', value: 0 },
      { label: '数据库总数', value: 0 },
      { label: 'K8s集群数量', value: 0 }
    ]
  },
  services: {
    title: '服务详情',
    items: [
      { label: '应用总数', value: 0 },
      { label: '业务线总数', value: 0 }
    ]
  },
  platform: {
    title: '平台详情',
    items: [
      { label: 'CPU', value: 0, unit: '%', color: '#667eea' },
      { label: '内存', value: 0, unit: '%', color: '#f093fb' },
      { label: '负载', value: 0, unit: '', color: '#4facfe' },
      { label: '磁盘', value: 0, unit: '%', color: '#fa709a' }
    ]
  },
  deployment: {
    title: '发布详情',
    items: [
      { label: '应用发布', value: 0 },
      { label: '任务执行', value: 0 },
      { label: '成功率', value: 0, unit: '%' }
    ]
  },
  monitor: {
    title: '用户信息',
    items: [
      { label: '今日登录次数', value: 0 },
      { label: '今日活跃用户', value: 0 },
      { label: '系统用户总数', value: 0 }
    ]
  }
})

// 快捷导航工具数据
const quickTools = reactive([])

// 编辑工具表单
const toolForm = reactive({
  title: '',
  icon: '',
  link: '',
  sort: 0
})

// 资源使用率类型
const resourceType = ref('memory') // cpu, memory, disk - 默认展示内存

// TOP 使用率主机数据
const topUsageData = ref({
  topCpu: [],
  topMemory: [],
  topDisk: [],
  period: '15天',
  startTime: '',
  endTime: ''
})

// 当前选中的天数（用于按钮高亮）
const selectedDays = ref(15)

// 图表实例
let heatChart = null

// 打开编辑弹窗
const openEditDialog = (tool, index) => {
  editingIndex.value = index
  editingTool.value = tool
  Object.assign(toolForm, {
    title: tool.title,
    icon: tool.icon,
    link: tool.link,
    sort: tool.sort || 0
  })
  editDialogVisible.value = true
}

// 添加新工具
const addNewTool = () => {
  editingIndex.value = -1
  editingTool.value = null
  Object.assign(toolForm, {
    title: '',
    icon: '',
    link: '',
    sort: 0
  })
  editDialogVisible.value = true
}

// 保存编辑
const saveToolEdit = async () => {
  if (!toolForm.title.trim()) {
    ElMessage.warning('请输入导航标题')
    return
  }

  if (!toolForm.icon) {
    ElMessage.warning('请上传导航图标')
    return
  }

  if (!toolForm.link.trim()) {
    ElMessage.warning('请输入链接地址')
    return
  }

  // 校验链接地址必须包含 http:// 或 https://
  const link = toolForm.link.trim()
  if (!link.startsWith('http://') && !link.startsWith('https://')) {
    ElMessage.warning('链接地址必须以 http:// 或 https:// 开头')
    return
  }

  try {
    if (editingIndex.value >= 0) {
      // 编辑现有工具
      await UpdateTool({
        id: editingTool.value.id,
        title: toolForm.title,
        icon: toolForm.icon,
        link: toolForm.link,
        sort: toolForm.sort
      })
      ElMessage.success('更新成功')
    } else {
      // 添加新工具
      await CreateTool({
        title: toolForm.title,
        icon: toolForm.icon,
        link: toolForm.link,
        sort: toolForm.sort
      })
      ElMessage.success('添加成功')
    }

    editDialogVisible.value = false
    // 重新加载导航工具列表
    await loadTools()
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败，请稍后重试')
  }
}

// 删除工具
const deleteTool = (index) => {
  const tool = quickTools[index]
  ElMessageBox.confirm('确定要删除这个导航吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await DeleteToolAPI(tool.id)
      ElMessage.success('删除成功')
      // 重新加载导航工具列表
      await loadTools()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败，请稍后重试')
    }
  }).catch(() => {})
}

// 上传图标
const handleIconUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请上传图片文件')
    return
  }

  // 验证文件大小（限制为2MB）
  if (file.size > 2 * 1024 * 1024) {
    ElMessage.error('图片大小不能超过2MB')
    return
  }

  try {
    const formData = new FormData()
    formData.append('file', file)

    const response = await UploadIcon(formData)
    if (response.data && response.data.code === 200) {
      toolForm.icon = response.data.data
      ElMessage.success('图标上传成功')
    } else {
      ElMessage.error(response.data?.message || '图标上传失败')
    }
  } catch (error) {
    console.error('上传图标失败:', error)
    ElMessage.error('图标上传失败，请稍后重试')
  }
}

// 触发文件选择
const triggerIconUpload = () => {
  document.getElementById('iconUpload').click()
}

// 点击导航项
const handleToolClick = (tool) => {
  if (!tool.link) return

  // 判断是外部链接还是内部路由
  if (tool.link.startsWith('http://') || tool.link.startsWith('https://')) {
    // 外部链接，新窗口打开
    window.open(tool.link, '_blank')
  } else {
    // 内部路由
    router.push(tool.link)
  }
}

// 获取资源使用率数据
const getResourceData = (type) => {
  const dataMap = {
    memory: topUsageData.value.topMemory,
    cpu: topUsageData.value.topCpu,
    disk: topUsageData.value.topDisk
  }

  const hosts = dataMap[type] || []

  // 去重处理：按 hostId + ip 组合去重，保留第一个
  const uniqueHosts = []
  const seenKeys = new Set()

  for (const host of hosts) {
    const key = `${host.hostId}-${host.ip}`
    if (!seenKeys.has(key)) {
      seenKeys.add(key)
      uniqueHosts.push(host)
    }
  }

  // 按平均使用率降序排序（从大到小），取前5个
  const sortedHosts = uniqueHosts
    .sort((a, b) => (b.avgUsage || 0) - (a.avgUsage || 0))
    .slice(0, 5)

  console.log(`${type} 数据排序结果:`, sortedHosts.map(h => ({ name: h.hostName, avgUsage: h.avgUsage })))

  return sortedHosts.map(host => ({
    name: host.hostName || host.ip || `主机-${host.hostId}`,
    value: parseFloat((host.avgUsage || 0).toFixed(2)),
    maxValue: parseFloat((host.maxUsage || 0).toFixed(2)),
    currentValue: parseFloat((host.currentUsage || 0).toFixed(2)),
    ip: host.ip,
    hostId: host.hostId
  }))
}

// 加载 TOP 使用率主机数据
const loadTopUsageHosts = async (days = 15, limit = 5) => {
  try {
    const response = await monitorApi.hostMonitor.getTopUsageHosts({ days, limit })
    if (response.data && response.data.code === 200) {
      topUsageData.value = response.data.data

      // 调试信息
      console.log('API 返回数据:', {
        topCpu: topUsageData.value.topCpu?.length || 0,
        topMemory: topUsageData.value.topMemory?.length || 0,
        topDisk: topUsageData.value.topDisk?.length || 0,
        period: topUsageData.value.period
      })

      updateResourceChart()
    }
  } catch (error) {
    console.error('加载TOP使用率主机数据失败:', error)
  }
}

// 初始化资源使用率图表
const initHeatChart = () => {
  const chartDom = document.getElementById('heatChart')
  if (!chartDom) return

  heatChart = echarts.init(chartDom)
  updateResourceChart()
}

// 更新资源使用率图表
const updateResourceChart = () => {
  if (!heatChart) return

  const data = getResourceData(resourceType.value)
  const titles = {
    memory: '内存使用率 TOP5',
    cpu: 'CPU使用率 TOP5',
    disk: '磁盘占用 TOP5'
  }

  const gradientColors = {
    memory: [
      { color: ['#4facfe', '#00f2fe'], shadowColor: 'rgba(79, 172, 254, 0.3)' },
      { color: ['#4facfe', '#00f2fe'], shadowColor: 'rgba(79, 172, 254, 0.25)' },
      { color: ['#4facfe', '#00f2fe'], shadowColor: 'rgba(79, 172, 254, 0.2)' },
      { color: ['#4facfe', '#00f2fe'], shadowColor: 'rgba(79, 172, 254, 0.15)' },
      { color: ['#4facfe', '#00f2fe'], shadowColor: 'rgba(79, 172, 254, 0.1)' }
    ],
    cpu: [
      { color: ['#667eea', '#764ba2'], shadowColor: 'rgba(102, 126, 234, 0.3)' },
      { color: ['#667eea', '#764ba2'], shadowColor: 'rgba(102, 126, 234, 0.25)' },
      { color: ['#667eea', '#764ba2'], shadowColor: 'rgba(102, 126, 234, 0.2)' },
      { color: ['#667eea', '#764ba2'], shadowColor: 'rgba(102, 126, 234, 0.15)' },
      { color: ['#667eea', '#764ba2'], shadowColor: 'rgba(102, 126, 234, 0.1)' }
    ],
    disk: [
      { color: ['#fa709a', '#fee140'], shadowColor: 'rgba(250, 112, 154, 0.3)' },
      { color: ['#fa709a', '#fee140'], shadowColor: 'rgba(250, 112, 154, 0.25)' },
      { color: ['#fa709a', '#fee140'], shadowColor: 'rgba(250, 112, 154, 0.2)' },
      { color: ['#fa709a', '#fee140'], shadowColor: 'rgba(250, 112, 154, 0.15)' },
      { color: ['#fa709a', '#fee140'], shadowColor: 'rgba(250, 112, 154, 0.1)' }
    ]
  }

  const option = {
    title: {
      text: titles[resourceType.value],
      left: 24,
      top: 16,
      textStyle: {
        fontSize: 18,
        fontWeight: 600,
        color: '#1a1a1a'
      },
      subtext: topUsageData.value.period ? `统计周期：${topUsageData.value.period}` : '',
      subtextStyle: {
        fontSize: 13,
        color: '#8c8c8c',
        fontWeight: 400
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
        shadowStyle: {
          color: 'rgba(0, 0, 0, 0.05)'
        }
      },
      backgroundColor: 'rgba(255, 255, 255, 0.98)',
      borderWidth: 0,
      padding: [16, 20],
      textStyle: {
        color: '#333',
        fontSize: 13
      },
      boxShadow: '0 8px 24px rgba(0, 0, 0, 0.12)',
      formatter: (params) => {
        const item = params[0]
        const dataItem = data[item.dataIndex]
        let tip = `<div style="font-weight: 600; font-size: 14px; color: #1a1a1a; margin-bottom: 10px;">${item.name}</div>`
        if (dataItem.ip) {
          tip += `<div style="color: #999; font-size: 12px; margin-bottom: 12px; padding-bottom: 10px; border-bottom: 1px solid #f0f0f0;">IP: ${dataItem.ip}</div>`
        }
        tip += `<div style="margin: 8px 0; display: flex; align-items: center; justify-content: space-between;">
          <span style="color: #666;">平均使用率</span>
          <span style="font-weight: 600; font-size: 15px; color: #1a1a1a;">${item.value}%</span>
        </div>`
        if (dataItem.maxValue !== undefined) {
          tip += `<div style="margin: 8px 0; display: flex; align-items: center; justify-content: space-between;">
            <span style="color: #666;">峰值使用率</span>
            <span style="font-weight: 600; font-size: 15px; color: #ff6b6b;">${dataItem.maxValue}%</span>
          </div>`
        }
        if (dataItem.currentValue !== undefined) {
          tip += `<div style="margin: 8px 0; display: flex; align-items: center; justify-content: space-between;">
            <span style="color: #666;">当前使用率</span>
            <span style="font-weight: 600; font-size: 15px; color: #52c41a;">${dataItem.currentValue}%</span>
          </div>`
        }
        return tip
      }
    },
    grid: {
      left: 150,
      right: 70,
      top: 80,
      bottom: 50
    },
    xAxis: {
      type: 'value',
      max: 100,
      axisLabel: {
        formatter: '{value}%',
        color: '#8c8c8c',
        fontSize: 12,
        fontWeight: 400
      },
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      splitLine: {
        lineStyle: {
          color: '#f0f0f0',
          type: 'dashed'
        }
      }
    },
    yAxis: {
      type: 'category',
      data: data.map(item => item.name),
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: '#262626',
        fontSize: 13,
        fontWeight: 500,
        margin: 16
      }
    },
    series: [
      {
        type: 'bar',
        data: data.map((item, index) => ({
          value: item.value,
          itemStyle: {
            borderRadius: [0, 8, 8, 0],
            color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
              { offset: 0, color: gradientColors[resourceType.value][index].color[0] },
              { offset: 1, color: gradientColors[resourceType.value][index].color[1] }
            ]),
            shadowColor: gradientColors[resourceType.value][index].shadowColor,
            shadowBlur: 12,
            shadowOffsetX: 2,
            shadowOffsetY: 2
          }
        })),
        barWidth: 18,
        label: {
          show: true,
          position: 'right',
          formatter: '{c}%',
          color: '#262626',
          fontSize: 13,
          fontWeight: 600,
          distance: 10
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 20,
            shadowOffsetX: 0,
            shadowOffsetY: 0
          }
        }
      }
    ]
  }
  heatChart.setOption(option)
}

// 切换资源类型
const changeResourceType = (type) => {
  resourceType.value = type
  updateResourceChart()
}

// 切换统计周期
const changePeriod = async (days) => {
  selectedDays.value = days
  await loadTopUsageHosts(days, 5)
}

// 加载导航工具列表
const loadTools = async () => {
  try {
    const response = await GetAllTools()
    if (response.data && response.data.code === 200) {
      // 清空现有数据
      quickTools.splice(0, quickTools.length)
      // 添加新数据
      quickTools.push(...response.data.data)
    }
  } catch (error) {
    console.error('加载导航工具失败:', error)
  }
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const response = await getDashboardStats()
    if (response.data && response.data.code === 200) {
      const data = response.data.data

      // 更新资产详情
      stats.assets.items[0].value = data.hostStats?.total || 0
      stats.assets.items[1].value = data.databaseStats?.total || 0
      stats.assets.items[2].value = data.k8sClusterStats?.total || 0

      // 更新服务详情
      stats.services.items[0].value = data.serviceStats?.total || 0
      stats.services.items[1].value = data.serviceStats?.businessLines || 0

      // 更新发布详情
      stats.deployment.items[0].value = data.deploymentStats?.total || 0
      stats.deployment.items[1].value = data.taskStats?.total || 0
      stats.deployment.items[2].value = data.deploymentStats?.successRate || 0

      // 加载用户统计信息
      await loadUserStats()
    }
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取告警源列表（支持FlashDuty和Zabbix数据源）
const getAlertSources = async () => {
  try {
    const res = await dashboardApi.getAlertSourceList({ pageNum: 1, pageSize: 100, status: 1 })
    if (res.data && (res.data.code === 0 || res.data.code === 200)) {
      const allSources = res.data.data.list || []

      // 保留 FlashDuty (type=1) 和 Zabbix (type=6) 类型的告警源
      alertSourceList.value = allSources.filter(source => source.type === 1 || source.type === 6)

      // 优先选择 Zabbix 数据源,其次选择 FlashDuty
      const zabbixSource = alertSourceList.value.find(source => source.type === 6)
      const flashDutySource = alertSourceList.value.find(source => source.type === 1)

      if (zabbixSource) {
        selectedAlertSource.value = zabbixSource
      } else if (flashDutySource) {
        selectedAlertSource.value = flashDutySource
      } else if (alertSourceList.value.length > 0) {
        selectedAlertSource.value = alertSourceList.value[0]
      }

      // 查找主机类型的告警源 (type=5)，用于加载平台详情
      const platformSource = allSources.find(source => source.type === 5)
      if (platformSource && platformSource.hostId) {
        loadPlatformMetrics(platformSource.hostId)
      }
    }
  } catch (error) {
    console.error('获取告警源列表失败:', error)
  }
}

// 加载平台详情监控数据
const loadPlatformMetrics = async (hostId) => {
  try {
    const response = await dashboardApi.getHostMetrics(hostId)
    if (response.data && response.data.code === 200) {
      const metrics = response.data.data

      // 更新平台详情数据
      stats.platform.items[0].value = parseFloat(metrics.cpuUsage?.toFixed(1)) || 0  // CPU
      stats.platform.items[1].value = parseFloat(metrics.memoryUsage?.toFixed(1)) || 0  // 内存
      stats.platform.items[2].value = parseFloat(metrics.load5?.toFixed(2)) || 0  // 负载（5分钟）
      stats.platform.items[3].value = parseFloat(metrics.diskUsage?.toFixed(1)) || 0  // 磁盘

      // 注意：负载不是百分比，可能需要调整单位显示
      stats.platform.items[2].unit = ''  // 负载没有单位
    }
  } catch (error) {
    console.error('加载平台监控数据失败:', error)
  }
}

// 加载用户统计信息
const loadUserStats = async () => {
  try {
    const response = await monitorApi.getUserStats()
    if (response.data && response.data.code === 200) {
      const data = response.data.data

      // 更新用户信息
      stats.monitor.items[0].value = data.todayLoginCount || 0  // 今日登录次数
      stats.monitor.items[1].value = data.todayActiveUserCount || 0  // 今日活跃用户
      stats.monitor.items[2].value = data.totalUserCount || 0  // 系统用户总数
    }
  } catch (error) {
    console.error('加载用户统计信息失败:', error)
  }
}

// 窗口大小改变时重绘图表
const handleResize = () => {
  heatChart?.resize()
}

// 生命周期
onMounted(async () => {
  // 并行加载所有数据,提升性能
  await Promise.all([
    loadData(),
    loadTools(),
    getAlertSources(),
    loadTopUsageHosts(15, 5)
  ])

  setTimeout(() => {
    initHeatChart()
  }, 100)

  window.addEventListener('resize', handleResize)
})

import { onUnmounted } from 'vue'
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  heatChart?.dispose()
})
</script>

<template>
  <div class="dashboard">
    <!-- 顶部统计卡片 -->
    <div class="stats-cards">
      <!-- 资产详情 -->
      <div class="stat-card asset-card">
        <div class="stat-header">
          <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
            <img src="@/assets/image/云主机服务器.svg" alt="资产" />
          </div>
          <div class="stat-title">{{ stats.assets.title }}</div>
        </div>
        <div class="stat-items">
          <div class="stat-item" v-for="(item, index) in stats.assets.items" :key="item.label">
            <div class="item-left">
              <img
                :src="index === 0 ? require('@/assets/image/linux.svg') :
                      index === 1 ? require('@/assets/image/数据库.svg') :
                      require('@/assets/image/k8s.svg')"
                class="item-icon"
                :alt="item.label"
              />
              <span class="item-label">{{ item.label }}</span>
            </div>
            <span class="item-value">{{ item.value }}</span>
          </div>
        </div>
      </div>

      <!-- 服务详情 -->
      <div class="stat-card service-card">
        <div class="stat-header">
          <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
            <img src="@/assets/image/应用.svg" alt="服务" />
          </div>
          <div class="stat-title">{{ stats.services.title }}</div>
        </div>
        <div class="stat-items">
          <div class="stat-item" v-for="(item, index) in stats.services.items" :key="item.label">
            <div class="item-left">
              <img
                :src="index === 0 ? require('@/assets/image/统计管理.svg') :
                      require('@/assets/image/分组.svg')"
                class="item-icon"
                :alt="item.label"
              />
              <span class="item-label">{{ item.label }}</span>
            </div>
            <span class="item-value">{{ item.value }}</span>
          </div>
        </div>
      </div>

      <!-- 平台详情 -->
      <div class="stat-card platform-card">
        <div class="stat-header">
          <div class="stat-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
            <img src="@/assets/image/电脑.svg" alt="平台" />
          </div>
          <div class="stat-title">{{ stats.platform.title }}</div>
        </div>
        <div class="gauge-items">
          <div class="gauge-item" v-for="item in stats.platform.items" :key="item.label">
            <div class="gauge-wrapper">
              <svg class="gauge-svg" viewBox="0 0 100 60">
                <path
                  class="gauge-bg"
                  d="M 10 50 A 40 40 0 0 1 90 50"
                  fill="none"
                  stroke="#f0f0f0"
                  stroke-width="8"
                  stroke-linecap="round"
                />
                <path
                  class="gauge-progress"
                  d="M 10 50 A 40 40 0 0 1 90 50"
                  fill="none"
                  :stroke="item.color"
                  stroke-width="8"
                  stroke-linecap="round"
                  :stroke-dasharray="`${item.value * 1.26} 126`"
                />
              </svg>
              <div class="gauge-value">{{ item.value }}{{ item.unit }}</div>
            </div>
            <div class="gauge-label">{{ item.label }}</div>
          </div>
        </div>
      </div>

      <!-- 发布详情 -->
      <div class="stat-card deploy-card">
        <div class="stat-header">
          <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
            <img src="@/assets/image/版本.svg" alt="发布" />
          </div>
          <div class="stat-title">{{ stats.deployment.title }}</div>
        </div>
        <div class="stat-items">
          <div class="stat-item" v-for="(item, index) in stats.deployment.items" :key="item.label">
            <div class="item-left">
              <img
                :src="index === 0 ? require('@/assets/image/fluentd.svg') :
                      index === 1 ? require('@/assets/image/renwu.svg') :
                      require('@/assets/image/Prometheus.svg')"
                class="item-icon"
                :alt="item.label"
              />
              <span class="item-label">{{ item.label }}</span>
            </div>
            <span class="item-value">{{ item.value }}{{ item.unit || '' }}</span>
          </div>
        </div>
      </div>

      <!-- 监控告警 -->
      <div class="stat-card monitor-card">
        <div class="stat-header">
          <div class="stat-icon" style="background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);">
            <img src="@/assets/image/公告.svg" alt="告警" />
          </div>
          <div class="stat-title">{{ stats.monitor.title }}</div>
        </div>
        <div class="stat-items">
          <div class="stat-item" v-for="item in stats.monitor.items" :key="item.label">
            <span class="item-label">{{ item.label }}</span>
            <span class="item-value">{{ item.value }}{{ item.unit || '' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 告警统计三卡片区域 - 根据数据源类型显示 -->
    <!-- Zabbix告警统计 -->
    <ZabbixAlertStats
      v-if="!selectedAlertSource || selectedAlertSource.type === 6"
      :accountId="selectedAlertSource?.keyId"
    />

    <!-- FlashDuty告警统计 -->
    <AlertStats
      v-else-if="selectedAlertSource.type === 1"
      :selected-alert-source="selectedAlertSource.id"
    />

    <!-- 图表区域 - 发布统计 + 业务分布 -->
    <DeploymentStats />

    <!-- 云账单查询区域 -->
    <BillingStats />

    <!-- 底部区域 -->
    <div class="bottom-row">
      <!-- 快捷导航工具 -->
      <div class="tools-card">
        <div class="tools-header">
          <div class="tools-title">快捷导航工具</div>
          <button class="add-tool-btn" @click="addNewTool">
            <svg viewBox="0 0 24 24" fill="currentColor" style="width: 16px; height: 16px;">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
            </svg>
            添加
          </button>
        </div>
        <div class="tools-grid">
          <div
            class="tool-item"
            v-for="(tool, index) in quickTools"
            :key="tool.id"
          >
            <div class="tool-actions">
              <button class="action-btn edit-btn" @click.stop="openEditDialog(tool, index)">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"/>
                </svg>
              </button>
              <button class="action-btn delete-btn" @click.stop="deleteTool(index)">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/>
                </svg>
              </button>
            </div>
            <div class="tool-content" @click="handleToolClick(tool)">
              <div class="tool-icon">
                <img v-if="tool.icon" :src="tool.icon" :alt="tool.title" />
                <div v-else class="icon-placeholder">?</div>
              </div>
              <div class="tool-info">
                <div class="tool-name">{{ tool.title }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 编辑弹窗 -->
      <el-dialog
        v-model="editDialogVisible"
        :title="editingIndex >= 0 ? '编辑导航' : '添加导航'"
        width="550px"
        :close-on-click-modal="false"
      >
        <div class="edit-form">
          <div class="form-item">
            <label><span class="required">*</span> 导航标题</label>
            <input
              v-model="toolForm.title"
              type="text"
              placeholder="例如：百度"
              class="form-input"
            />
          </div>

          <div class="form-item">
            <label><span class="required">*</span> 导航图标</label>
            <div class="icon-upload-wrapper">
              <div class="icon-preview-box">
                <img v-if="toolForm.icon" :src="toolForm.icon" alt="图标预览" />
                <div v-else class="empty-icon">
                  <svg viewBox="0 0 24 24" fill="currentColor">
                    <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"/>
                  </svg>
                  <span>暂无图标</span>
                </div>
              </div>
              <button class="upload-btn" @click="triggerIconUpload">
                <svg viewBox="0 0 24 24" fill="currentColor" style="width: 16px; height: 16px;">
                  <path d="M9 16h6v-6h4l-7-7-7 7h4zm-4 2h14v2H5z"/>
                </svg>
                选择图标
              </button>
              <input
                id="iconUpload"
                type="file"
                accept="image/*"
                @change="handleIconUpload"
                style="display: none;"
              />
            </div>
            <div class="form-tip">支持 PNG、JPG、SVG 格式，大小不超过 2MB</div>
          </div>

          <div class="form-item">
            <label><span class="required">*</span> 链接地址</label>
            <input
              v-model="toolForm.link"
              type="text"
              placeholder="例如：https://www.baidu.com/"
              class="form-input"
            />
            <div class="form-tip">必须以 http:// 或 https:// 开头</div>
          </div>
        </div>

        <template #footer>
          <div class="dialog-footer">
            <button class="btn-cancel" @click="editDialogVisible = false">取消</button>
            <button class="btn-confirm" @click="saveToolEdit">保存</button>
          </div>
        </template>
      </el-dialog>

      <!-- 资源使用率 - 现代化设计 -->
      <div class="chart-card-premium">
        <div class="chart-control-bar">
          <div class="control-left">
            <div class="resource-type-selector">
              <button
                :class="['type-btn', { active: resourceType === 'memory' }]"
                @click="changeResourceType('memory')"
              >
                <span class="type-icon memory-icon"></span>
                <span class="type-label">内存</span>
              </button>
              <button
                :class="['type-btn', { active: resourceType === 'cpu' }]"
                @click="changeResourceType('cpu')"
              >
                <span class="type-icon cpu-icon"></span>
                <span class="type-label">CPU</span>
              </button>
              <button
                :class="['type-btn', { active: resourceType === 'disk' }]"
                @click="changeResourceType('disk')"
              >
                <span class="type-icon disk-icon"></span>
                <span class="type-label">磁盘</span>
              </button>
            </div>
          </div>
          <div class="control-right">
            <div class="period-selector">
              <button
                :class="['period-chip', { selected: selectedDays === 3 }]"
                @click="changePeriod(3)"
              >
                3天
              </button>
              <button
                :class="['period-chip', { selected: selectedDays === 7 }]"
                @click="changePeriod(7)"
              >
                7天
              </button>
              <button
                :class="['period-chip', { selected: selectedDays === 15 }]"
                @click="changePeriod(15)"
              >
                15天
              </button>
            </div>
          </div>
        </div>
        <div id="heatChart" class="premium-chart-container"></div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.dashboard {
  padding: 0;
  background: transparent;
  min-height: calc(100vh - 60px);
}

// 顶部统计卡片
.stats-cards {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.stat-card {
  position: relative;
  background: white;
  border-radius: 8px;
  padding: 10px 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);

    &::before {
      opacity: 1;
    }
  }
}

.asset-card::before {
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
}

.service-card::before {
  background: linear-gradient(90deg, #f093fb 0%, #f5576c 100%);
}

.deploy-card::before {
  background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
}

.monitor-card::before {
  background: linear-gradient(90deg, #fa709a 0%, #fee140 100%);
}

.platform-card::before {
  background: linear-gradient(90deg, #43e97b 0%, #38f9d7 100%);
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.stat-title {
  font-size: 13px;
  font-weight: 600;
  color: #333;
  flex: 1;
}

.stat-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);

  img {
    width: 18px;
    height: 18px;
    filter: brightness(0) invert(1);
  }

  svg {
    width: 18px;
    height: 18px;
  }
}

.stat-items {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
  transition: all 0.2s ease;

  &:hover {
    padding-left: 4px;
  }
}

.item-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.item-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  opacity: 0.75;
}

.item-label {
  font-size: 12px;
  color: #666;
  font-weight: 400;
}

.item-value {
  font-size: 15px;
  font-weight: 700;
  color: #333;
  letter-spacing: -0.3px;
}

// 仪表盘样式
.gauge-items {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 4px 6px;
}

.gauge-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.gauge-wrapper {
  position: relative;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.gauge-svg {
  width: 65px;
  height: 40px;
}

.gauge-progress {
  transition: stroke-dasharray 0.6s ease;
}

.gauge-value {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -20%);
  font-size: 11px;
  font-weight: 700;
  color: #333;
}

.gauge-label {
  font-size: 10px;
  color: #666;
  font-weight: 500;
  text-align: center;
}

// 底部行
.bottom-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-top: 20px;
}

// 快捷工具
.tools-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.tools-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.tools-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.add-tool-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: #2196f3;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    background: #1976d2;
    transform: translateY(-1px);
  }
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 15px;
}

.tool-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 6px;
  border-radius: 6px;
  transition: all 0.3s ease;

  &:hover {
    background: #f8f9fa;

    .tool-actions {
      opacity: 1;
    }
  }
}

.tool-actions {
  position: absolute;
  top: 4px;
  right: 4px;
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.action-btn {
  width: 24px;
  height: 24px;
  border: none;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;

  svg {
    width: 14px;
    height: 14px;
  }

  &.edit-btn {
    background: #e3f2fd;
    color: #2196f3;

    &:hover {
      background: #2196f3;
      color: white;
    }
  }

  &.delete-btn {
    background: #ffebee;
    color: #f44336;

    &:hover {
      background: #f44336;
      color: white;
    }
  }
}

.tool-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  width: fit-content;
}

.tool-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: transparent;
  border: none;

  img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  .icon-placeholder {
    font-size: 20px;
    color: #c0c4cc;
  }
}

.tool-info {
  text-align: center;
}

.tool-name {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}

// 编辑弹窗样式
.edit-form {
  padding: 10px 0;
}

.form-item {
  margin-bottom: 24px;

  label {
    display: block;
    font-size: 14px;
    color: #333;
    margin-bottom: 10px;
    font-weight: 500;

    .required {
      color: #f56c6c;
      margin-right: 4px;
    }
  }
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s ease;
  box-sizing: border-box;

  &:focus {
    outline: none;
    border-color: #2196f3;
  }

  &::placeholder {
    color: #c0c4cc;
  }
}

.form-tip {
  margin-top: 6px;
  font-size: 12px;
  color: #909399;
}

.icon-upload-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;
}

.icon-preview-box {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px dashed #dcdfe6;
  overflow: hidden;
  background: #fafafa;

  img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  .empty-icon {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    color: #c0c4cc;

    svg {
      width: 32px;
      height: 32px;
    }

    span {
      font-size: 11px;
    }
  }
}

.upload-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  background: #f5f7fa;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  color: #606266;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    background: #2196f3;
    border-color: #2196f3;
    color: white;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 10px;
}

.btn-cancel,
.btn-confirm {
  padding: 8px 20px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-cancel {
  background: #f5f7fa;
  color: #606266;

  &:hover {
    background: #e0e3e8;
  }
}

.btn-confirm {
  background: #2196f3;
  color: white;

  &:hover {
    background: #1976d2;
  }
}

// 资源使用率 - Premium 设计
.chart-card-premium {
  background: #ffffff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  height: 460px;
  border: 1px solid #f0f0f0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  &:hover {
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
    transform: translateY(-2px);
  }
}

.chart-control-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.control-left {
  flex: 1;
}

.control-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

// 资源类型选择器 - 左侧大按钮
.resource-type-selector {
  display: flex;
  gap: 12px;
}

.type-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  border: 2px solid transparent;
  background: #fafafa;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
    transition: left 0.5s ease;
  }

  &:hover {
    background: #f5f5f5;
    transform: translateY(-2px);

    &::before {
      left: 100%;
    }
  }

  &.active {
    border-color: currentColor;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.08), rgba(118, 75, 162, 0.08));
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);

    .type-label {
      font-weight: 600;
    }
  }

  // 内存按钮（第1个）
  &.active.type-btn:nth-child(1) {
    border-color: #4facfe;
    background: linear-gradient(135deg, rgba(79, 172, 254, 0.1), rgba(0, 242, 254, 0.1));

    .type-icon {
      background: linear-gradient(135deg, #4facfe, #00f2fe);
    }

    .type-label {
      background: linear-gradient(135deg, #4facfe, #00f2fe);
      -webkit-background-clip: text;
      background-clip: text;
      -webkit-text-fill-color: transparent;
    }
  }

  // CPU按钮（第2个）
  &.active.type-btn:nth-child(2) {
    border-color: #667eea;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));

    .type-icon {
      background: linear-gradient(135deg, #667eea, #764ba2);
    }

    .type-label {
      background: linear-gradient(135deg, #667eea, #764ba2);
      -webkit-background-clip: text;
      background-clip: text;
      -webkit-text-fill-color: transparent;
    }
  }

  // 磁盘按钮（第3个）
  &.active.type-btn:nth-child(3) {
    border-color: #fa709a;
    background: linear-gradient(135deg, rgba(250, 112, 154, 0.1), rgba(254, 225, 64, 0.1));

    .type-icon {
      background: linear-gradient(135deg, #fa709a, #fee140);
    }

    .type-label {
      background: linear-gradient(135deg, #fa709a, #fee140);
      -webkit-background-clip: text;
      background-clip: text;
      -webkit-text-fill-color: transparent;
    }
  }
}

.type-icon {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #d9d9d9;
  transition: all 0.3s ease;
}

.type-label {
  font-size: 14px;
  color: #595959;
  font-weight: 500;
  transition: all 0.3s ease;
}

// 周期选择器 - 右侧小标签
.period-selector {
  display: flex;
  gap: 6px;
  padding: 4px;
  background: #fafafa;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
}

.period-chip {
  padding: 6px 14px;
  border: none;
  background: transparent;
  color: #8c8c8c;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;

  &:hover {
    color: #262626;
    background: rgba(0, 0, 0, 0.02);
  }

  &.selected {
    background: #ffffff;
    color: #1a1a1a;
    font-weight: 600;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08),
                0 1px 2px rgba(0, 0, 0, 0.04);
  }
}

// 图表容器
.premium-chart-container {
  width: 100%;
  height: calc(100% - 90px);
}

// 响应式设计
@media (max-width: 1400px) {
  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .bottom-row {
    grid-template-columns: 1fr;
  }

  .tools-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .dashboard {
    padding: 12px;
  }

  .stats-cards {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .tools-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
