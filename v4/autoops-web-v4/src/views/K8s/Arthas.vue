<script setup>
import { ref, onMounted, reactive, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Cpu, Download, Upload, Refresh, Monitor, Search } from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import * as echarts from 'echarts'

// 当前活动的标签页
const activeTab = ref('dashboard')

// 集群ID
const clusterId = ref('')

// 参数选择
const params = reactive({
  namespace: 'default',
  workloadType: '',
  workloadName: '',
  pod: '',
  container: '',
  processId: ''
})

// 加载状态
const loading = ref(false)

// 集群、命名空间、工作负载、Pod、容器列表
const clusterList = ref([])
const namespaceList = ref([])
const workloadList = ref([])
const podList = ref([])
const containerList = ref([])
const processList = ref([])

// Arthas安装状态
const arthasStatus = reactive({
  installed: false,
  javaPath: '',
  message: '',
  checked: false
})

// 诊断数据
const dashboardData = ref(null)
const dashboardJsonStr = ref('')
const threadData = ref(null)
const jvmData = ref(null)
const memoryData = ref(null)
const sysenvData = ref(null)
const syspropData = ref(null)
const classSearchResults = ref([])
const decompileResult = ref('')
const flameGraphData = ref(null)
const flameGraphDuration = ref(30)
const flameGraphLoading = ref(false)
const flameGraphStatus = ref('')   // pending | running | done | failed
const flameGraphTaskId = ref('')
const flameGraphChartRef = ref(null)
const flameGraphChartHeight = ref(600)
let flameGraphChart = null
let flameGraphPollTimer = null

// 搜索参数
const searchParams = reactive({
  className: '',
  methodName: '',
  pattern: ''
})

// 加载集群列表
const loadClusters = async () => {
  try {
    const res = await k8sApi.getClusterList()
    const responseData = res.data || res
    if (responseData.code === 200 || responseData.code === 0 || responseData.success) {
      const clusters = responseData.data?.list || responseData.data || []
      clusterList.value = Array.isArray(clusters) ? clusters : []
      if (clusterList.value.length > 0 && !clusterId.value) {
        const onlineCluster = clusterList.value.find(c => c.status === 2)
        clusterId.value = onlineCluster ? onlineCluster.id : clusterList.value[0].id
      }
    }
  } catch (error) {
    console.error('加载集群失败:', error)
  }
}

// 加载命名空间列表
const loadNamespaces = async () => {
  if (!clusterId.value) return
  try {
    const res = await k8sApi.getNamespaces(clusterId.value)
    const responseData = res.data || res
    if (responseData.code === 200 || responseData.code === 0 || responseData.success) {
      const namespaces = responseData.data?.namespaces || responseData.data || []
      namespaceList.value = Array.isArray(namespaces)
        ? namespaces.map(ns => (typeof ns === 'string' ? ns : ns.name)).filter(Boolean)
        : []
    }
  } catch (error) {
    console.error('加载命名空间失败:', error)
  }
}

// 加载工作负载列表
const loadWorkloads = async () => {
  if (!clusterId.value || !params.namespace) return
  try {
    const res = await k8sApi.getWorkloadList(clusterId.value, params.namespace)
    const responseData = res.data || res
    if (responseData.code === 200 || responseData.code === 0 || responseData.success) {
      const workloads = responseData.data?.workloads || responseData.data || []
      workloadList.value = Array.isArray(workloads)
        ? workloads.map(w => ({
            name: w.name,
            type: (w.type || w.kind || 'deployment').toLowerCase()
          }))
        : []
    }
  } catch (error) {
    console.error('加载工作负载失败:', error)
  }
}

// 加载Pod列表
const loadPods = async () => {
  if (!clusterId.value || !params.namespace || !params.workloadName || !params.workloadType) return
  try {
    const res = await k8sApi.getWorkloadPods(
      clusterId.value,
      params.namespace,
      params.workloadType,
      params.workloadName
    )
    const responseData = res.data || res
    if (responseData.code === 200 || responseData.code === 0 || responseData.success) {
      const pods = responseData.data?.pods || responseData.data?.list || responseData.data || []
      podList.value = Array.isArray(pods)
        ? pods.map(p => (typeof p === 'string' ? p : p.name)).filter(Boolean)
        : []
    }
  } catch (error) {
    console.error('加载Pod失败:', error)
  }
}

// 加载容器列表
const loadContainers = async () => {
  if (!clusterId.value || !params.namespace || !params.pod) return
  try {
    const res = await k8sApi.getPodContainers(clusterId.value, params.namespace, params.pod)
    const responseData = res.data || res
    if (responseData.code === 200 || responseData.code === 0 || responseData.success) {
      const containers = responseData.data?.containers || responseData.data?.list || responseData.data || []
      containerList.value = Array.isArray(containers)
        ? containers.map(c => (typeof c === 'string' ? c : c.name)).filter(Boolean)
        : []
      if (containerList.value.length > 0) {
        params.container = containerList.value[0]
      }
    }
  } catch (error) {
    console.error('加载容器失败:', error)
  }
}

// 检查Arthas环境
const checkArthasInstallStatus = async () => {
  if (!validateParams()) return
  loading.value = true
  try {
    const res = await k8sApi.checkArthasInstall(clusterId.value, {
      namespace: params.namespace,
      pod: params.pod,
      container: params.container
    })
    const responseData = res.data || res
    if ((responseData.code === 200 || responseData.code === 0 || responseData.success) && responseData.data) {
      arthasStatus.installed = responseData.data.installed
      arthasStatus.javaPath = responseData.data.javaPath || ''
      arthasStatus.message = responseData.data.message || ''
      arthasStatus.checked = true
      if (responseData.data.installed) {
        ElMessage.success('环境检查完成，可以进行诊断')
      } else {
        ElMessage.warning(responseData.data.message)
      }
    }
  } catch (error) {
    ElMessage.error('检查失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 下载 Arthas
const downloadArthas = async () => {
  loading.value = true
  try {
    const res = await k8sApi.downloadArthas()
    const responseData = res.data || res
    if (responseData.code === 200 || responseData.code === 0 || responseData.success) {
      ElMessage.success('Arthas 下载成功')
    } else {
      ElMessage.error('下载失败: ' + responseData.message)
    }
  } catch (error) {
    ElMessage.error('下载失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 上传 Arthas
const uploadArthas = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.jar'
  input.onchange = async (e) => {
    const file = e.target.files[0]
    if (!file) return
    if (file.name !== 'arthas-boot.jar') {
      ElMessage.warning('请选择名为 arthas-boot.jar 的文件')
      return
    }
    loading.value = true
    try {
      const formData = new FormData()
      formData.append('file', file)
      const res = await k8sApi.uploadArthas(formData)
      const responseData = res.data || res
      if (responseData.code === 200 || responseData.code === 0 || responseData.success) {
        ElMessage.success(`Arthas 上传成功 (${(file.size / 1024 / 1024).toFixed(2)}MB)`)
      } else {
        ElMessage.error('上传失败: ' + responseData.message)
      }
    } catch (error) {
      ElMessage.error('上传失败: ' + error.message)
    } finally {
      loading.value = false
    }
  }
  input.click()
}

// 加载Java进程列表
const loadJavaProcesses = async () => {
  if (!validateParams()) return
  loading.value = true
  try {
    const res = await k8sApi.listJavaProcesses(clusterId.value, {
      namespace: params.namespace,
      pod: params.pod,
      container: params.container
    })
    const responseData = res.data || res
    if (responseData.code === 200 || responseData.code === 0 || responseData.success) {
      const processes = responseData.data?.list || responseData.data || []
      processList.value = Array.isArray(processes) ? processes : []
      ElMessage.success('共找到 ' + processList.value.length + ' 个Java进程')
    }
  } catch (error) {
    ElMessage.error('加载进程列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 参数验证
const validateParams = () => {
  if (!clusterId.value) { ElMessage.warning('请选择集群'); return false }
  if (!params.namespace) { ElMessage.warning('请选择命名空间'); return false }
  if (!params.workloadName) { ElMessage.warning('请选择工作负载'); return false }
  if (!params.pod) { ElMessage.warning('请选择Pod'); return false }
  if (!params.container) { ElMessage.warning('请选择容器'); return false }
  return true
}

// 获取JVM仪表板
const getDashboard = async () => {
  if (!params.processId) { ElMessage.warning('请先选择Java进程'); return }
  loading.value = true
  try {
    const res = await k8sApi.getArthasDashboard(clusterId.value, {
      namespace: params.namespace,
      pod: params.pod,
      container: params.container,
      processId: params.processId
    })
    const responseData = res.data || res
    if ((responseData.code === 200 || responseData.code === 0 || responseData.success) && responseData.data) {
      dashboardData.value = responseData.data
      dashboardJsonStr.value = JSON.stringify(responseData.data, null, 2)
      ElMessage.success('获取成功')
    }
  } catch (error) {
    ElMessage.error('获取失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 获取线程信息
const getThreadInfo = async () => {
  if (!params.processId) { ElMessage.warning('请先选择Java进程'); return }
  loading.value = true
  try {
    const res = await k8sApi.getArthasThreadInfo(clusterId.value, {
      namespace: params.namespace,
      pod: params.pod,
      container: params.container,
      processId: params.processId
    })
    const responseData = res.data || res
    if ((responseData.code === 200 || responseData.code === 0 || responseData.success) && responseData.data) {
      threadData.value = responseData.data
      ElMessage.success('获取成功')
    }
  } catch (error) {
    ElMessage.error('获取失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 获取JVM信息
const getJvmInfo = async () => {
  if (!params.processId) { ElMessage.warning('请先选择Java进程'); return }
  loading.value = true
  try {
    const res = await k8sApi.getArthasJvmInfo(clusterId.value, {
      namespace: params.namespace,
      pod: params.pod,
      container: params.container,
      processId: params.processId
    })
    const responseData = res.data || res
    if ((responseData.code === 200 || responseData.code === 0 || responseData.success) && responseData.data) {
      jvmData.value = responseData.data
      ElMessage.success('获取成功')
    }
  } catch (error) {
    ElMessage.error('获取失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 获取内存信息
const getMemoryInfo = async () => {
  if (!params.processId) { ElMessage.warning('请先选择Java进程'); return }
  loading.value = true
  try {
    const res = await k8sApi.getArthasMemoryInfo(clusterId.value, {
      namespace: params.namespace,
      pod: params.pod,
      container: params.container,
      processId: params.processId
    })
    const responseData = res.data || res
    if ((responseData.code === 200 || responseData.code === 0 || responseData.success) && responseData.data) {
      memoryData.value = responseData.data
      ElMessage.success('获取成功')
    }
  } catch (error) {
    ElMessage.error('获取失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 获取系统环境变量
const getSysenv = async () => {
  if (!params.processId) { ElMessage.warning('请先选择Java进程'); return }
  loading.value = true
  try {
    const res = await k8sApi.getArthasSysenv(clusterId.value, {
      namespace: params.namespace,
      pod: params.pod,
      container: params.container,
      processId: params.processId
    })
    const responseData = res.data || res
    if ((responseData.code === 200 || responseData.code === 0 || responseData.success) && responseData.data) {
      sysenvData.value = responseData.data
      ElMessage.success('获取成功')
    }
  } catch (error) {
    ElMessage.error('获取失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 获取系统属性
const getSysprop = async () => {
  if (!params.processId) { ElMessage.warning('请先选择Java进程'); return }
  loading.value = true
  try {
    const res = await k8sApi.getArthasSysprop(clusterId.value, {
      namespace: params.namespace,
      pod: params.pod,
      container: params.container,
      processId: params.processId
    })
    const responseData = res.data || res
    if ((responseData.code === 200 || responseData.code === 0 || responseData.success) && responseData.data) {
      syspropData.value = responseData.data
      ElMessage.success('获取成功')
    }
  } catch (error) {
    ElMessage.error('获取失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 搜索类
const searchClass = async () => {
  if (!params.processId) { ElMessage.warning('请先选择Java进程'); return }
  if (!searchParams.pattern) { ElMessage.warning('请输入搜索模式'); return }
  loading.value = true
  try {
    const res = await k8sApi.searchArthasClass(clusterId.value, {
      namespace: params.namespace,
      pod: params.pod,
      container: params.container,
      processId: params.processId,
      pattern: searchParams.pattern
    })
    const responseData = res.data || res
    if ((responseData.code === 200 || responseData.code === 0 || responseData.success) && responseData.data) {
      const results = responseData.data?.list || responseData.data || []
      classSearchResults.value = Array.isArray(results) ? results : []
      ElMessage.success('找到 ' + classSearchResults.value.length + ' 个类')
    }
  } catch (error) {
    ElMessage.error('搜索失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 反编译类
const decompileClass = async (className) => {
  if (!params.processId) { ElMessage.warning('请先选择Java进程'); return }
  loading.value = true
  try {
    const res = await k8sApi.getArthasJad(clusterId.value, {
      namespace: params.namespace,
      pod: params.pod,
      container: params.container,
      processId: params.processId,
      className: className
    })
    const responseData = res.data || res
    if ((responseData.code === 200 || responseData.code === 0 || responseData.success) && responseData.data) {
      decompileResult.value = responseData.data
      ElMessage.success('反编译成功')
    }
  } catch (error) {
    ElMessage.error('反编译失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 生成火焰图 - 第一步：启动采样
const generateFlameGraph = async () => {
  if (!params.processId) { ElMessage.warning('请先选择Java进程'); return }

  // 清除上次的轮询
  if (flameGraphPollTimer) {
    clearInterval(flameGraphPollTimer)
    flameGraphPollTimer = null
  }

  flameGraphLoading.value = true
  flameGraphData.value = null
  flameGraphStatus.value = 'pending'
  flameGraphTaskId.value = ''

  try {
    const res = await k8sApi.startArthasProfiler(clusterId.value, {
      namespace: params.namespace,
      pod: params.pod,
      container: params.container,
      processId: params.processId,
      duration: flameGraphDuration.value
    })
    const responseData = res.data || res
    if ((responseData.code === 200 || responseData.code === 0 || responseData.success) && responseData.data) {
      const task = responseData.data
      flameGraphTaskId.value = task.taskId
      flameGraphStatus.value = task.status || 'pending'
      ElMessage.info(`采样任务已提交，预计 ${flameGraphDuration.value} 秒后完成`)
      // 开始轮询
      pollFlameGraphTask(task.taskId)
    } else {
      ElMessage.error('启动失败: ' + (responseData.message || '未知错误'))
      flameGraphLoading.value = false
      flameGraphStatus.value = 'failed'
    }
  } catch (error) {
    ElMessage.error('启动失败: ' + error.message)
    flameGraphLoading.value = false
    flameGraphStatus.value = 'failed'
  }
}

// 生成火焰图 - 第二步：轮询任务状态
const pollFlameGraphTask = (taskId) => {
  flameGraphPollTimer = setInterval(async () => {
    try {
      const res = await k8sApi.getArthasProfilerTask(clusterId.value, taskId)
      const responseData = res.data || res
      if (!(responseData.code === 200 || responseData.code === 0 || responseData.success)) {
        clearInterval(flameGraphPollTimer)
        flameGraphPollTimer = null
        ElMessage.error('查询失败: ' + (responseData.message || '未知错误'))
        flameGraphLoading.value = false
        flameGraphStatus.value = 'failed'
        return
      }

      const task = responseData.data
      flameGraphStatus.value = task.status

      if (task.status === 'done') {
        clearInterval(flameGraphPollTimer)
        flameGraphPollTimer = null
        flameGraphLoading.value = false
        flameGraphData.value = task.result
        ElMessage.success(task.result?.message || '火焰图生成成功')
        await nextTick()
        const treeData = task.result?.data
        if (treeData) {
          renderFlameGraph(treeData)
        } else {
          ElMessage.error('火焰图数据为空')
        }
      } else if (task.status === 'failed') {
        clearInterval(flameGraphPollTimer)
        flameGraphPollTimer = null
        flameGraphLoading.value = false
        ElMessage.error('采样失败: ' + (task.message || '未知原因'))
      }
      // pending / running 继续等待
    } catch (error) {
      clearInterval(flameGraphPollTimer)
      flameGraphPollTimer = null
      flameGraphLoading.value = false
      flameGraphStatus.value = 'failed'
      ElMessage.error('轮询失败: ' + error.message)
    }
  }, 5000)
}

// 渲染火焰图
const renderFlameGraph = (data) => {
  if (!flameGraphChartRef.value || !data) {
    console.error('火焰图渲染失败: 容器或数据为空', { ref: flameGraphChartRef.value, data })
    return
  }

  console.log('开始渲染火焰图，根节点:', data)

  // 销毁旧图表
  if (flameGraphChart) {
    flameGraphChart.dispose()
  }

  // 将树形数据转换为火焰图数据
  const flameData = []
  let maxDepth = 0

  const traverse = (node, depth, start) => {
    if (!node) return start

    const item = {
      name: node.name,
      value: node.value,
      depth: depth,
      start: start,
      end: start + node.value
    }
    flameData.push(item)
    maxDepth = Math.max(maxDepth, depth)

    let currentStart = start
    if (node.children && node.children.length > 0) {
      node.children.forEach(child => {
        currentStart = traverse(child, depth + 1, currentStart)
      })
    }

    return start + node.value
  }

  traverse(data, 0, 0)

  console.log('火焰图数据转换完成:', {
    总节点数: flameData.length,
    最大深度: maxDepth,
    前5个节点: flameData.slice(0, 5)
  })

  // 动态计算图表高度
  const rowHeight = 20
  flameGraphChartHeight.value = Math.max(400, (maxDepth + 1) * rowHeight + 100)

  // 创建新图表
  flameGraphChart = echarts.init(flameGraphChartRef.value)

  const option = {
    title: {
      text: 'CPU 火焰图',
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 600,
        color: '#6b9bd1'
      }
    },
    tooltip: {
      formatter: (params) => {
        if (!params.data) return ''
        const percent = params.data.value ? ((params.data.value / data.value) * 100).toFixed(2) : 0
        return `<strong>${params.data.name}</strong><br/>样本数: ${params.data.value}<br/>占比: ${percent}%`
      }
    },
    grid: {
      left: 10,
      right: 10,
      top: 60,
      bottom: 10
    },
    xAxis: {
      type: 'value',
      max: data.value,
      show: false
    },
    yAxis: {
      type: 'value',
      max: maxDepth + 1,
      show: false,
      inverse: true
    },
    series: [
      {
        type: 'custom',
        renderItem: (params, api) => {
          const item = flameData[params.dataIndex]
          if (!item) return null

          const x = api.coord([item.start, item.depth])[0]
          const y = api.coord([item.start, item.depth])[1]
          const width = api.coord([item.end, item.depth])[0] - x
          const height = api.coord([item.start, item.depth + 1])[1] - y

          // 颜色根据深度变化
          const hue = (item.depth * 30) % 360
          const color = `hsl(${hue}, 70%, 60%)`

          return {
            type: 'rect',
            shape: {
              x: x,
              y: y,
              width: Math.max(width, 0),
              height: Math.max(height - 1, 0)
            },
            style: {
              fill: color,
              stroke: '#fff',
              lineWidth: 1
            },
            textContent: {
              style: {
                text: width > 50 ? item.name : '',
                fontSize: 11,
                fill: '#000',
                overflow: 'truncate',
                width: width - 4
              }
            },
            textConfig: {
              position: 'inside',
              distance: 2
            }
          }
        },
        data: flameData
      }
    ]
  }

  console.log('设置 ECharts 配置')
  flameGraphChart.setOption(option)
  console.log('火焰图渲染完成')

  // 响应式调整
  const resizeHandler = () => {
    flameGraphChart?.resize()
  }
  window.removeEventListener('resize', resizeHandler)
  window.addEventListener('resize', resizeHandler)
}

// 初始化
onMounted(() => {
  loadClusters()
})

// 监听集群选择
watch(() => clusterId.value, (newClusterId) => {
  if (newClusterId) {
    namespaceList.value = []
    workloadList.value = []
    podList.value = []
    containerList.value = []
    processList.value = []
    params.namespace = 'default'
    params.workloadType = ''
    params.workloadName = ''
    params.pod = ''
    params.container = ''
    params.processId = ''
    arthasStatus.checked = false
    loadNamespaces()
  }
})

// 监听命名空间选择
watch(() => params.namespace, (newNamespace) => {
  if (newNamespace && clusterId.value) {
    workloadList.value = []
    podList.value = []
    containerList.value = []
    params.workloadName = ''
    params.workloadType = ''
    params.pod = ''
    params.container = ''
    params.processId = ''
    arthasStatus.checked = false
    loadWorkloads()
  }
})

// 监听工作负载选择
watch(() => params.workloadName, (newWorkload) => {
  if (newWorkload && params.namespace && clusterId.value) {
    const found = workloadList.value.find(w => w.name === newWorkload)
    params.workloadType = found ? found.type : 'deployment'
    podList.value = []
    containerList.value = []
    params.pod = ''
    params.container = ''
    params.processId = ''
    arthasStatus.checked = false
    loadPods()
  }
})

// 监听Pod选择
watch(() => params.pod, (newPod) => {
  if (newPod && params.namespace && clusterId.value) {
    containerList.value = []
    params.container = ''
    params.processId = ''
    arthasStatus.checked = false
    loadContainers()
  }
})

// 监听容器变更
watch(() => params.container, () => {
  processList.value = []
  params.processId = ''
  arthasStatus.checked = false
})
</script>

<template>
  <div class="arthas-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="page-title-group">
        <div class="page-title-icon">
          <el-icon><Cpu /></el-icon>
        </div>
        <div>
          <h2 class="page-title">Arthas Java 诊断</h2>
          <p class="page-subtitle">在线Java应用诊断，包括JVM监控、线程分析、代码反编译等功能</p>
        </div>
      </div>
    </div>

    <!-- 选择器 + 操作栏 -->
    <div class="selector-bar">
      <div class="selector-inputs">
        <div class="selector-item">
          <span class="selector-label">集群</span>
          <el-select v-model="clusterId" placeholder="请选择集群" class="selector-select">
            <el-option
              v-for="cluster in clusterList"
              :key="cluster.id"
              :label="cluster.name || cluster.id"
              :value="cluster.id"
            />
          </el-select>
        </div>

        <div class="selector-item">
          <span class="selector-label">命名空间</span>
          <el-select
            v-model="params.namespace"
            placeholder="请选择命名空间"
            :disabled="!clusterId"
            class="selector-select"
          >
            <el-option v-for="ns in namespaceList" :key="ns" :label="ns" :value="ns" />
          </el-select>
        </div>

        <div class="selector-item">
          <span class="selector-label">工作负载</span>
          <el-select
            v-model="params.workloadName"
            placeholder="请选择工作负载"
            :disabled="!params.namespace"
            class="selector-select workload-select"
          >
            <el-option
              v-for="workload in workloadList"
              :key="workload.name"
              :label="`[${workload.type}] ${workload.name}`"
              :value="workload.name"
            />
          </el-select>
        </div>

        <div class="selector-item">
          <span class="selector-label">Pod</span>
          <el-select
            v-model="params.pod"
            placeholder="请选择Pod"
            :disabled="!params.workloadName"
            class="selector-select"
          >
            <el-option v-for="pod in podList" :key="pod" :label="pod" :value="pod" />
          </el-select>
        </div>

        <div class="selector-item">
          <span class="selector-label">容器</span>
          <el-select
            v-model="params.container"
            placeholder="请选择容器"
            :disabled="!params.pod"
            class="selector-select"
          >
            <el-option v-for="c in containerList" :key="c" :label="c" :value="c" />
          </el-select>
        </div>

        <div class="selector-item">
          <span class="selector-label">进程</span>
          <el-select
            v-model="params.processId"
            placeholder="请先刷新进程"
            :disabled="processList.length === 0"
            class="selector-select process-select"
          >
            <el-option
              v-for="process in processList"
              :key="process.processId"
              :label="`PID: ${process.processId} - ${process.mainClass}`"
              :value="process.processId"
            />
          </el-select>
        </div>
      </div>

      <div class="selector-actions">
        <el-tooltip content="检查Arthas环境是否可用" placement="top">
          <el-button
            :icon="Monitor"
            :loading="loading"
            @click="checkArthasInstallStatus"
            class="action-btn check-btn"
          >
            检查环境
          </el-button>
        </el-tooltip>
        <el-tooltip content="加载Java进程列表" placement="top">
          <el-button
            :icon="Refresh"
            :loading="loading"
            @click="loadJavaProcesses"
            class="action-btn refresh-btn"
          >
            刷新进程
          </el-button>
        </el-tooltip>
        <el-tooltip content="从网络下载Arthas" placement="top">
          <el-button
            :icon="Download"
            :loading="loading"
            @click="downloadArthas"
            class="action-btn download-btn"
          >
            下载
          </el-button>
        </el-tooltip>
        <el-tooltip content="手动上传 arthas-boot.jar" placement="top">
          <el-button
            :icon="Upload"
            :loading="loading"
            @click="uploadArthas"
            class="action-btn upload-btn"
          >
            上传
          </el-button>
        </el-tooltip>
      </div>
    </div>

    <!-- 环境检查状态 -->
    <div v-if="arthasStatus.checked" class="status-banner">
      <el-alert
        :title="`Arthas环境: ${arthasStatus.installed ? '✓ 就绪，可以开始诊断' : '✗ 不可用'}`"
        :type="arthasStatus.installed ? 'success' : 'warning'"
        :closable="false"
        show-icon
      >
        <template #default>
          <span v-if="arthasStatus.javaPath">Java路径: {{ arthasStatus.javaPath }}&nbsp;&nbsp;</span>
          <span v-if="arthasStatus.message">{{ arthasStatus.message }}</span>
        </template>
      </el-alert>
    </div>

    <!-- 诊断内容区 -->
    <div class="content-wrapper">
      <el-tabs v-if="params.processId" v-model="activeTab" class="diagnosis-tabs">

        <!-- JVM仪表板 -->
        <el-tab-pane name="dashboard" label="JVM仪表板">
          <div class="tab-content">
            <div class="tab-toolbar">
              <el-button :icon="Refresh" :loading="loading" @click="getDashboard" size="small" class="modern-btn primary-btn">
                刷新数据
              </el-button>
            </div>
            <div v-if="dashboardData">
              <el-row :gutter="16" class="stat-row">
                <el-col :span="6">
                  <div class="stat-card">
                    <div class="stat-label">活跃线程</div>
                    <div class="stat-value">{{ dashboardData.threads?.live ?? '-' }}</div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="stat-card">
                    <div class="stat-label">峰值线程</div>
                    <div class="stat-value">{{ dashboardData.threads?.peak ?? '-' }}</div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="stat-card">
                    <div class="stat-label">守护线程</div>
                    <div class="stat-value">{{ dashboardData.threads?.daemon ?? '-' }}</div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="stat-card">
                    <div class="stat-label">堆已用 (MB)</div>
                    <div class="stat-value stat-value--highlight">
                      {{ dashboardData.memory?.heap?.used
                        ? (dashboardData.memory.heap.used / 1024 / 1024).toFixed(1)
                        : '-' }}
                    </div>
                  </div>
                </el-col>
              </el-row>
              <el-collapse style="margin-top: 16px">
                <el-collapse-item title="堆内存详情" name="heap">
                  <el-descriptions :column="3" border>
                    <el-descriptions-item label="已使用 (MB)">
                      {{ dashboardData.memory?.heap?.used ? (dashboardData.memory.heap.used / 1024 / 1024).toFixed(2) : '-' }}
                    </el-descriptions-item>
                    <el-descriptions-item label="已提交 (MB)">
                      {{ dashboardData.memory?.heap?.committed ? (dashboardData.memory.heap.committed / 1024 / 1024).toFixed(2) : '-' }}
                    </el-descriptions-item>
                    <el-descriptions-item label="最大值 (MB)">
                      {{ dashboardData.memory?.heap?.max ? (dashboardData.memory.heap.max / 1024 / 1024).toFixed(2) : '-' }}
                    </el-descriptions-item>
                  </el-descriptions>
                </el-collapse-item>
                <el-collapse-item title="原始 JSON 数据" name="json">
                  <el-input v-model="dashboardJsonStr" type="textarea" :rows="10" readonly />
                </el-collapse-item>
              </el-collapse>
            </div>
            <el-empty v-else description="点击刷新获取数据" />
          </div>
        </el-tab-pane>

        <!-- 线程分析 -->
        <el-tab-pane name="thread" label="线程分析">
          <div class="tab-content">
            <div class="tab-toolbar">
              <el-button :icon="Refresh" :loading="loading" @click="getThreadInfo" size="small" class="modern-btn primary-btn">
                刷新数据
              </el-button>
            </div>
            <el-table
              v-if="threadData?.threads"
              :data="threadData.threads"
              stripe
              border
              :header-cell-style="{ background: 'rgba(103,126,234,0.08)', color: '#6b9bd1', fontWeight: '600' }"
            >
              <el-table-column prop="threadId" label="线程ID" width="100" />
              <el-table-column prop="threadName" label="线程名称" min-width="150" />
              <el-table-column prop="state" label="状态" width="120">
                <template #default="{ row }">
                  <el-tag
                    :type="row.state === 'RUNNABLE' ? 'success' : row.state === 'WAITING' ? 'warning' : 'info'"
                    size="small"
                  >
                    {{ row.state }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="cpu" label="CPU %" width="100" />
              <el-table-column prop="time" label="时间" width="120" />
            </el-table>
            <el-empty v-else description="点击刷新获取数据" />
          </div>
        </el-tab-pane>

        <!-- JVM信息 -->
        <el-tab-pane name="jvm" label="JVM信息">
          <div class="tab-content">
            <div class="tab-toolbar">
              <el-button :icon="Refresh" :loading="loading" @click="getJvmInfo" size="small" class="modern-btn primary-btn">
                刷新数据
              </el-button>
            </div>
            <el-descriptions v-if="jvmData" :column="2" border>
              <el-descriptions-item label="操作系统">{{ jvmData.osName }}</el-descriptions-item>
              <el-descriptions-item label="OS版本">{{ jvmData.osVersion }}</el-descriptions-item>
              <el-descriptions-item label="处理器数">{{ jvmData.processorCount }}</el-descriptions-item>
              <el-descriptions-item label="Java版本">{{ jvmData.javaVersion }}</el-descriptions-item>
              <el-descriptions-item label="Java厂商">{{ jvmData.javaVendor }}</el-descriptions-item>
              <el-descriptions-item label="Class Path" :span="2">
                <el-input v-model="jvmData.classPath" type="textarea" :rows="2" readonly />
              </el-descriptions-item>
            </el-descriptions>
            <el-empty v-else description="点击刷新获取数据" />
          </div>
        </el-tab-pane>

        <!-- 内存信息 -->
        <el-tab-pane name="memory" label="内存信息">
          <div class="tab-content">
            <div class="tab-toolbar">
              <el-button :icon="Refresh" :loading="loading" @click="getMemoryInfo" size="small" class="modern-btn primary-btn">
                刷新数据
              </el-button>
            </div>
            <div v-if="memoryData">
              <el-divider content-position="left">堆内存使用</el-divider>
              <el-row :gutter="16">
                <el-col :span="6">
                  <div class="stat-card">
                    <div class="stat-label">初始值 (MB)</div>
                    <div class="stat-value">
                      {{ memoryData.heapMemoryUsage?.init ? (memoryData.heapMemoryUsage.init / 1024 / 1024).toFixed(2) : '-' }}
                    </div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="stat-card">
                    <div class="stat-label">已使用 (MB)</div>
                    <div class="stat-value stat-value--highlight">
                      {{ memoryData.heapMemoryUsage?.used ? (memoryData.heapMemoryUsage.used / 1024 / 1024).toFixed(2) : '-' }}
                    </div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="stat-card">
                    <div class="stat-label">已提交 (MB)</div>
                    <div class="stat-value">
                      {{ memoryData.heapMemoryUsage?.committed ? (memoryData.heapMemoryUsage.committed / 1024 / 1024).toFixed(2) : '-' }}
                    </div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="stat-card">
                    <div class="stat-label">最大值 (MB)</div>
                    <div class="stat-value">
                      {{ memoryData.heapMemoryUsage?.max ? (memoryData.heapMemoryUsage.max / 1024 / 1024).toFixed(2) : '-' }}
                    </div>
                  </div>
                </el-col>
              </el-row>
            </div>
            <el-empty v-else description="点击刷新获取数据" />
          </div>
        </el-tab-pane>

        <!-- 环境变量 -->
        <el-tab-pane name="sysenv" label="环境变量">
          <div class="tab-content">
            <div class="tab-toolbar">
              <el-button :icon="Refresh" :loading="loading" @click="getSysenv" size="small" class="modern-btn primary-btn">
                刷新数据
              </el-button>
            </div>
            <el-table
              v-if="sysenvData"
              :data="Object.entries(sysenvData).map(([k, v]) => ({ key: k, value: v }))"
              stripe
              border
              :header-cell-style="{ background: 'rgba(103,126,234,0.08)', color: '#6b9bd1', fontWeight: '600' }"
            >
              <el-table-column prop="key" label="变量名" min-width="200" />
              <el-table-column prop="value" label="变量值" show-overflow-tooltip />
            </el-table>
            <el-empty v-else description="点击刷新获取数据" />
          </div>
        </el-tab-pane>

        <!-- 系统属性 -->
        <el-tab-pane name="sysprop" label="系统属性">
          <div class="tab-content">
            <div class="tab-toolbar">
              <el-button :icon="Refresh" :loading="loading" @click="getSysprop" size="small" class="modern-btn primary-btn">
                刷新数据
              </el-button>
            </div>
            <el-table
              v-if="syspropData"
              :data="Object.entries(syspropData).map(([k, v]) => ({ key: k, value: v }))"
              stripe
              border
              :header-cell-style="{ background: 'rgba(103,126,234,0.08)', color: '#6b9bd1', fontWeight: '600' }"
            >
              <el-table-column prop="key" label="属性名" min-width="200" />
              <el-table-column prop="value" label="属性值" show-overflow-tooltip />
            </el-table>
            <el-empty v-else description="点击刷新获取数据" />
          </div>
        </el-tab-pane>

        <!-- 代码分析 -->
        <el-tab-pane name="code" label="代码分析">
          <div class="tab-content">
            <div class="tab-toolbar">
              <el-input
                v-model="searchParams.pattern"
                placeholder="输入类名模式（支持通配符 *）"
                style="width: 320px"
                clearable
              />
              <el-button
                :icon="Search"
                :loading="loading"
                @click="searchClass"
                size="small"
                class="modern-btn primary-btn"
                style="margin-left: 8px"
              >
                搜索类
              </el-button>
            </div>
            <div v-if="classSearchResults.length > 0" style="margin-top: 16px">
              <el-table
                :data="classSearchResults.map(cls => ({ className: cls }))"
                stripe
                border
                :header-cell-style="{ background: 'rgba(103,126,234,0.08)', color: '#6b9bd1', fontWeight: '600' }"
              >
                <el-table-column prop="className" label="类名" show-overflow-tooltip />
                <el-table-column label="操作" width="100">
                  <template #default="{ row }">
                    <el-button size="small" type="primary" @click="decompileClass(row.className)">
                      反编译
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
            <div v-if="decompileResult" style="margin-top: 16px">
              <el-divider content-position="left">反编译结果</el-divider>
              <el-input v-model="decompileResult" type="textarea" :rows="15" readonly />
            </div>
            <el-empty
              v-if="classSearchResults.length === 0 && !decompileResult"
              description="输入类名模式开始搜索"
            />
          </div>
        </el-tab-pane>

        <!-- 火焰图 -->
        <el-tab-pane name="flamegraph" label="火焰图">
          <div class="tab-content">
            <div class="tab-toolbar">
              <span style="margin-right: 8px; color: #606266; font-weight: 500;">采样时长:</span>
              <el-radio-group v-model="flameGraphDuration" size="small">
                <el-radio-button :label="10">10秒</el-radio-button>
                <el-radio-button :label="30">30秒</el-radio-button>
                <el-radio-button :label="60">60秒</el-radio-button>
                <el-radio-button :label="120">120秒</el-radio-button>
              </el-radio-group>
              <el-button
                :icon="Cpu"
                :loading="flameGraphLoading"
                @click="generateFlameGraph"
                size="small"
                class="modern-btn primary-btn"
                style="margin-left: 12px"
              >
                生成火焰图
              </el-button>
            </div>
            <el-alert
              v-if="!flameGraphData && !flameGraphLoading"
              title="火焰图说明"
              type="info"
              :closable="false"
              style="margin-bottom: 16px"
            >
              <template #default>
                <div style="font-size: 13px; line-height: 1.6;">
                  <p style="margin: 0 0 8px 0;">火焰图用于分析 CPU 性能热点，帮助定位代码中的性能瓶颈。</p>
                  <p style="margin: 0;"><strong>推荐时长：</strong></p>
                  <ul style="margin: 4px 0 0 20px; padding: 0;">
                    <li>10秒 - 极高并发场景，快速采样</li>
                    <li>30-60秒 - 最佳推荐，覆盖完整请求链路</li>
                    <li>120秒 - 低并发或复杂业务场景</li>
                  </ul>
                </div>
              </template>
            </el-alert>
            <div v-if="flameGraphLoading" style="text-align: center; padding: 60px 20px;">
              <el-icon class="is-loading" style="font-size: 48px; color: #409EFF; margin-bottom: 16px;">
                <Refresh />
              </el-icon>
              <p style="font-size: 14px; color: #606266;">
                {{ flameGraphStatus === 'running' ? '正在采样中...' : '任务等待启动...' }}
              </p>
              <p style="font-size: 13px; color: #909399;">每 5 秒自动查询一次，采样时长 {{ flameGraphDuration }} 秒</p>
              <el-tag v-if="flameGraphTaskId" type="info" style="margin-top: 8px; font-size: 11px;">
                任务ID: {{ flameGraphTaskId }}
              </el-tag>
            </div>
            <div v-else-if="flameGraphData" style="margin-top: 16px;">
              <el-divider content-position="left">火焰图结果</el-divider>
              <div style="border: 1px solid #dcdfe6; border-radius: 8px; overflow: hidden; background: #fff; padding: 16px;">
                <div
                  ref="flameGraphChartRef"
                  :style="{ width: '100%', height: flameGraphChartHeight + 'px' }"
                ></div>
              </div>
              <el-descriptions :column="2" border style="margin-top: 16px;">
                <el-descriptions-item label="状态消息">{{ flameGraphData.message || '-' }}</el-descriptions-item>
                <el-descriptions-item label="采样时长">{{ flameGraphDuration }} 秒</el-descriptions-item>
              </el-descriptions>
              <el-alert
                title="使用提示"
                type="info"
                :closable="false"
                style="margin-top: 16px;"
              >
                <template #default>
                  <div style="font-size: 13px; line-height: 1.6;">
                    <p style="margin: 0;">• 横向宽度表示该方法占用 CPU 时间，越宽占用越多</p>
                    <p style="margin: 4px 0 0 0;">• 纵向深度表示调用栈层级，从上到下依次展开</p>
                    <p style="margin: 4px 0 0 0;">• 鼠标悬停可查看详细的样本数和占比</p>
                  </div>
                </template>
              </el-alert>
            </div>
            <el-empty
              v-else
              description="选择采样时长后点击生成火焰图"
            />
          </div>
        </el-tab-pane>

      </el-tabs>

      <!-- 未选择进程提示 -->
      <div v-else class="empty-state">
        <el-empty description="请依次选择集群、命名空间、工作负载、Pod、容器，点击「刷新进程」后选择进程开始诊断" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.arthas-container {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.page-title-group {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 22px;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  line-height: 1.3;
}

.page-subtitle {
  margin: 4px 0 0 0;
  font-size: 13px;
  color: #909399;
  line-height: 1.4;
}

/* 选择器栏 */
.selector-bar {
  margin-bottom: 12px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.selector-inputs {
  display: flex;
  gap: 8px;
  flex: 1;
  flex-wrap: nowrap;
  align-items: center;
  min-width: 0;
}

.selector-item {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
  flex-shrink: 0;
}

.selector-label {
  font-size: 12px;
  color: #606266;
  white-space: nowrap;
  font-weight: 500;
}

.selector-select {
  width: 110px;
}

.workload-select {
  width: 140px;
}

.process-select {
  width: 160px;
}

.selector-bar :deep(.el-select .el-input__wrapper) {
  border-radius: 8px;
  box-shadow: none;
  border: 1px solid #dcdfe6;
  transition: all 0.3s ease;
}

.selector-bar :deep(.el-select .el-input__wrapper:hover) {
  border-color: #667eea;
}

.selector-bar :deep(.el-select .el-input.is-focus .el-input__wrapper) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.15);
}

/* 操作按钮区 */
.selector-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.action-btn {
  border-radius: 8px;
  font-weight: 500;
  font-size: 13px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.18);
}

.check-btn {
  background: linear-gradient(45deg, #67C23A, #85CE61);
  color: #fff;
}

.refresh-btn {
  background: linear-gradient(45deg, #409EFF, #66B3FF);
  color: #fff;
}

.download-btn {
  background: linear-gradient(45deg, #667eea, #764ba2);
  color: #fff;
}

.upload-btn {
  background: linear-gradient(45deg, #E6A23C, #EEBE77);
  color: #fff;
}

/* 状态横幅 */
.status-banner {
  margin-bottom: 12px;
  border-radius: 12px;
  overflow: hidden;
}

/* 内容区域 */
.content-wrapper {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  min-height: 400px;
}

/* Tab 样式 */
:deep(.diagnosis-tabs .el-tabs__header) {
  margin: 0;
  padding: 0 20px;
  background: rgba(103, 126, 234, 0.05);
  border-bottom: 1px solid rgba(103, 126, 234, 0.15);
}

:deep(.diagnosis-tabs .el-tabs__nav-wrap::after) {
  display: none;
}

:deep(.diagnosis-tabs .el-tabs__item) {
  padding: 0 20px;
  height: 48px;
  line-height: 48px;
  font-size: 13px;
  color: #606266;
  font-weight: 500;
  transition: all 0.2s ease;
}

:deep(.diagnosis-tabs .el-tabs__item:hover) {
  color: #667eea;
}

:deep(.diagnosis-tabs .el-tabs__item.is-active) {
  color: #667eea;
  font-weight: 600;
}

:deep(.diagnosis-tabs .el-tabs__active-bar) {
  background: linear-gradient(45deg, #667eea, #764ba2);
  height: 3px;
  border-radius: 2px;
}

:deep(.diagnosis-tabs .el-tabs__content) {
  padding: 0;
}

/* Tab 内容 */
.tab-content {
  padding: 20px;
}

.tab-toolbar {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(103, 126, 234, 0.1);
}

/* 统计卡片 */
.stat-row {
  margin-bottom: 8px;
}

.stat-card {
  text-align: center;
  padding: 20px 16px;
  border-radius: 12px;
  background: rgba(103, 126, 234, 0.05);
  border: 1px solid rgba(103, 126, 234, 0.12);
  transition: all 0.3s ease;
}

.stat-card:hover {
  background: rgba(103, 126, 234, 0.1);
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.15);
}

.stat-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 8px;
  font-weight: 500;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #6b9bd1;
}

.stat-value--highlight {
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

/* 表格样式 */
:deep(.el-table) {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
}

:deep(.el-table__body tr:hover > td) {
  background-color: rgba(102, 126, 234, 0.06) !important;
}

/* 按钮样式 */
:deep(.modern-btn) {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

:deep(.modern-btn:hover) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.18);
}

:deep(.primary-btn) {
  background: linear-gradient(45deg, #409EFF, #66B3FF);
  color: white;
}

/* 空状态 */
.empty-state {
  padding: 60px 20px;
}

/* 响应式 */
@media (max-width: 1400px) {
  .selector-select {
    width: 120px;
  }
  .workload-select {
    width: 160px;
  }
  .process-select {
    width: 180px;
  }
}

@media (max-width: 1100px) {
  .selector-bar {
    flex-direction: column;
    align-items: flex-start;
  }
  .selector-actions {
    width: 100%;
    justify-content: flex-end;
  }
}

@media (max-width: 768px) {
  .selector-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  .selector-select,
  .workload-select,
  .process-select {
    width: 100%;
  }
  :deep(.diagnosis-tabs .el-tabs__item) {
    padding: 0 12px;
    font-size: 12px;
  }
}
</style>
