<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import yaml from 'js-yaml'
import {
  Search,
  Refresh,
  Plus,
  ArrowDown,
  Delete,
  View,
  Setting,
  Warning,
  Monitor,
  DocumentCopy,
  Document,
  Connection
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import {
  getStoredClusterId,
  saveStoredClusterId,
  resolvePreferredClusterId
} from '@/utils/k8sClusterSelection'
import { getYamlContentFromApiResponse } from './utils/yaml'

// Import modular components
import ClusterSelector from './pods/ClusterSelector.vue'
import NamespaceSelector from './pods/NamespaceSelector.vue'
import PodConfigDialog from './pods/PodConfigDialog.vue'
import CreatePodDialog from './pods/CreatePodDialog.vue'
import { openYamlDialog } from './services/yamlDialog'

const router = useRouter()
const route = useRoute()
const ALL_WORKLOAD_TYPE = 'all'
const STANDALONE_PODS_TYPE = 'pods'
const WORKLOAD_SCOPE_ALL = 'all'
const WORKLOAD_SCOPE_ABNORMAL = 'abnormal'
const WORKLOAD_TAB_QUERY_KEY = 'tab'
const WORKLOAD_TAB_STORAGE_KEY = 'k8s-workloads-active-tab'

const loading = ref(false)

// 工作负载类型选项
const workloadTypeOptions = [
  { label: '全部', value: ALL_WORKLOAD_TYPE },
  { label: 'Deployment', value: 'deployments' },
  { label: 'StatefulSet', value: 'statefulsets' },
  { label: 'DaemonSet', value: 'daemonsets' },
  { label: 'Job', value: 'jobs' },
  { label: 'CronJob', value: 'cronjobs' },
  { label: 'Pods', value: STANDALONE_PODS_TYPE }
]

const validWorkloadTypes = new Set(workloadTypeOptions.map(option => option.value))

const getRouteTab = () => {
  const routeTab = route.query[WORKLOAD_TAB_QUERY_KEY]
  return Array.isArray(routeTab) ? routeTab[0] : routeTab
}

const normalizeWorkloadType = (type) =>
  validWorkloadTypes.has(type) ? type : ALL_WORKLOAD_TYPE

const getStoredTab = () => {
  if (typeof window === 'undefined') return ''
  return window.sessionStorage.getItem(WORKLOAD_TAB_STORAGE_KEY) || ''
}

const getInitialWorkloadType = () =>
  normalizeWorkloadType(getRouteTab() || getStoredTab())

const initialWorkloadType = getInitialWorkloadType()
const activeTab = ref(initialWorkloadType)
const queryParams = reactive({
  name: '',
  type: initialWorkloadType,
  namespace: 'default',
  scope: WORKLOAD_SCOPE_ALL
})

const workloadScopeOptions = [
  { label: '全部工作负载', value: WORKLOAD_SCOPE_ALL },
  { label: '仅看异常', value: WORKLOAD_SCOPE_ABNORMAL }
]

const persistTab = (tabName) => {
  if (typeof window === 'undefined') return

  const normalizedTab = normalizeWorkloadType(tabName)

  if (normalizedTab === ALL_WORKLOAD_TYPE) {
    window.sessionStorage.removeItem(WORKLOAD_TAB_STORAGE_KEY)
    return
  }

  window.sessionStorage.setItem(WORKLOAD_TAB_STORAGE_KEY, normalizedTab)
}

const syncTabState = (tabName) => {
  const normalizedTab = normalizeWorkloadType(tabName)
  const currentRouteTab = getRouteTab()
  const nextQuery = { ...route.query }

  persistTab(normalizedTab)

  if (normalizedTab === ALL_WORKLOAD_TYPE) {
    if (!currentRouteTab) return
    delete nextQuery[WORKLOAD_TAB_QUERY_KEY]
  } else {
    if (currentRouteTab === normalizedTab) return
    nextQuery[WORKLOAD_TAB_QUERY_KEY] = normalizedTab
  }

  router.replace({
    path: route.path,
    query: nextQuery
  })
}

const restoreTabState = () => {
  const routeTab = getRouteTab()
  const restoredTab = normalizeWorkloadType(routeTab || getStoredTab())

  activeTab.value = restoredTab
  queryParams.type = restoredTab
  syncTabState(restoredTab)
}

const isAllWorkloadsTab = computed(() => getRequestWorkloadType() === ALL_WORKLOAD_TYPE)
const isStandalonePodsTab = computed(() => getRequestWorkloadType() === STANDALONE_PODS_TYPE)
const isAbnormalMode = computed(() => queryParams.scope === WORKLOAD_SCOPE_ABNORMAL)
const workloadNameLabel = computed(() => (isStandalonePodsTab.value ? 'Pod 名称' : '工作负载名称'))
const workloadNamePlaceholder = computed(() => (
  isStandalonePodsTab.value
    ? '请输入 Pod 名称，支持模糊搜索'
    : '请输入工作负载名称，支持模糊搜索'
))
const replicaColumnLabel = computed(() => (isStandalonePodsTab.value ? 'Pod 状态' : '容器组数量'))
const workloadTableEmptyText = computed(() => {
  if (isStandalonePodsTab.value) {
    if (isAbnormalMode.value) {
      return '当前条件下暂无异常 Pod'
    }

    return queryParams.name?.trim() ? '未匹配到 Pod' : '暂无独立 Pod 数据'
  }

  if (isAbnormalMode.value) {
    return '当前条件下暂无异常工作负载'
  }

  return queryParams.name?.trim() ? '未匹配到工作负载' : '暂无工作负载数据'
})

const tableData = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const selectedClusterId = ref(getStoredClusterId())
const clusterList = ref([])
// clusterList 已移至 ClusterSelector 组件
// namespaceList 和 namespaceLoading 已移至 NamespaceSelector 组件

// 对话框状态
const scaleDialogVisible = ref(false)
const workloadLabelsDialogVisible = ref(false)
const allImagesDialogVisible = ref(false)
const podConfigDialogVisible = ref(false)
const schedulingDialogVisible = ref(false)
const createPodDialogVisible = ref(false)

// 当前操作的工作负载或Pod
const currentWorkload = ref({})
const createPodDialogRef = ref(null)

// 扩容缩容表单
const scaleForm = reactive({
  replicas: 1
})

// 获取集群列表
const fetchClusterList = async () => {
  try {
    const response = await k8sApi.getClusterList()
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      const clusters = responseData.data?.list || responseData.data || []
      clusterList.value = clusters.map(cluster => ({
        id: cluster.id,
        name: cluster.name,
        status: cluster.status
      }))

      const resolvedClusterId = resolvePreferredClusterId(
        clusterList.value,
        selectedClusterId.value || getStoredClusterId()
      )

      if (resolvedClusterId) {
        selectedClusterId.value = resolvedClusterId
        saveStoredClusterId(resolvedClusterId)
      }
      
      console.log('集群列表加载成功:', clusterList.value)
    } else {
      ElMessage.error(responseData.message || '获取集群列表失败')
    }
  } catch (error) {
    console.error('获取集群列表失败:', error)
    ElMessage.warning('无法获取集群列表，请检查后端服务')
  }
}

// namespaceRequestPromise 已移至 NamespaceSelector 组件

// fetchNamespaceList 已移至 NamespaceSelector 组件

// 处理标签页切换
const handleTabChange = (tabName) => {
  const normalizedTab = normalizeWorkloadType(tabName)
  console.log('标签页切换到:', normalizedTab)
  activeTab.value = normalizedTab
  queryParams.type = normalizedTab
  syncTabState(normalizedTab)

  if (normalizedTab === ALL_WORKLOAD_TYPE) {
    queryParams.name = ''
    queryParams.scope = WORKLOAD_SCOPE_ALL
  }

  currentPage.value = 1
  handleQuery()
}

// 监听activeTab变化，同步到queryParams.type
watch(activeTab, (newType) => {
  const normalizedType = normalizeWorkloadType(newType)
  queryParams.type = normalizedType
  persistTab(normalizedType)
})

watch(
  () => route.query[WORKLOAD_TAB_QUERY_KEY],
  (newTab) => {
    const routeTab = Array.isArray(newTab) ? newTab[0] : newTab
    const normalizedTab = normalizeWorkloadType(routeTab)

    if (normalizedTab === activeTab.value) {
      return
    }

    activeTab.value = normalizedTab
    queryParams.type = normalizedTab
    currentPage.value = 1

    if (!selectedClusterId.value) {
      resetTableData()
      return
    }

    if (!queryParams.namespace && normalizedTab !== ALL_WORKLOAD_TYPE) {
      resetTableData()
      return
    }

    handleQuery()
  }
)

const resetTableData = () => {
  tableData.value = []
  total.value = 0
}

const getRequestWorkloadType = () =>
  !queryParams.type || queryParams.type === ALL_WORKLOAD_TYPE ? ALL_WORKLOAD_TYPE : queryParams.type

const getKeyword = () => (queryParams.name || '').trim()
const getResourceNamespace = (resource = {}) =>
  resource?.namespace || currentWorkload.value?.namespace || queryParams.namespace

const getStandalonePodStatus = (pod = {}) => pod.status || pod.phase || 'Unknown'

const isStandalonePodReady = (pod = {}) => {
  const status = getStandalonePodStatus(pod)
  return ['Running', 'Succeeded', 'Completed'].includes(status)
}

const isStandalonePodAbnormal = (pod = {}) => !isStandalonePodReady(pod)

const extractPodImages = (containers = []) => {
  const images = containers
    .map(container => {
      if (typeof container === 'string') return ''
      return container?.image || container?.imageName || ''
    })
    .filter(Boolean)

  return [...new Set(images)]
}

const extractPodResourcesForTable = (containers = []) => {
  const resources = {
    requests: { cpu: '0', memory: '0' },
    limits: { cpu: '0', memory: '0' }
  }
  let hasRequestCPU = false
  let hasRequestMemory = false
  let hasLimitCPU = false
  let hasLimitMemory = false

  containers.forEach(container => {
    const containerResources = container?.resources || {}
    const requests = containerResources.requests || {}
    const limits = containerResources.limits || {}

    if (requests.cpu) {
      resources.requests.cpu = addCpuResources(resources.requests.cpu, requests.cpu)
      hasRequestCPU = true
    }
    if (requests.memory) {
      resources.requests.memory = addMemoryResources(resources.requests.memory, requests.memory)
      hasRequestMemory = true
    }
    if (limits.cpu) {
      resources.limits.cpu = addCpuResources(resources.limits.cpu, limits.cpu)
      hasLimitCPU = true
    }
    if (limits.memory) {
      resources.limits.memory = addMemoryResources(resources.limits.memory, limits.memory)
      hasLimitMemory = true
    }
  })

  return {
    requests: {
      cpu: hasRequestCPU ? resources.requests.cpu : '',
      memory: hasRequestMemory ? resources.requests.memory : ''
    },
    limits: {
      cpu: hasLimitCPU ? resources.limits.cpu : '',
      memory: hasLimitMemory ? resources.limits.memory : ''
    }
  }
}

const mapStandalonePodToTableRow = (pod) => {
  const readyReplicas = isStandalonePodReady(pod) ? 1 : 0
  const podResources = extractPodResourcesForTable(pod.containers || [])
  const createdAt = pod.createdAt || pod.creationTimestamp || pod.startTime || ''

  return {
    id: pod.name,
    name: pod.name,
    type: 'pod',
    namespace: pod.namespace || queryParams.namespace,
    replicas: `${readyReplicas}/1`,
    readyReplicas,
    totalReplicas: 1,
    images: extractPodImages(pod.containers || []),
    labels: pod.labels || {},
    status: getStandalonePodStatus(pod),
    age: formatAge(createdAt),
    updateTime: createdAt,
    updatedAt: pod.updatedAt || createdAt,
    conditions: pod.conditions || [],
    resources: podResources,
    rawData: pod
  }
}

// 查询工作负载列表
const handleQuery = async () => {
  const queryStartTime = Date.now()
  
  try {
    if (!selectedClusterId.value) {
      ElMessage.warning('请选择一个集群')
      return
    }
    
    if (!queryParams.namespace && !isAllWorkloadsTab.value) {
      ElMessage.warning('请选择命名空间')
      return
    }
    
    console.log('🔍 开始查询工作负载:', {
      clusterId: selectedClusterId.value,
      namespace: queryParams.namespace,
      scope: queryParams.scope,
      type: getRequestWorkloadType(),
      keyword: getKeyword(),
      page: currentPage.value,
      pageSize: pageSize.value
    })
    
    loading.value = true
    
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    const keyword = getKeyword()
    let response

    if (isStandalonePodsTab.value) {
      response = await k8sApi.getStandalonePodList(
        selectedClusterId.value,
        queryParams.namespace,
        keyword ? { keyword } : {}
      )
    } else if (isAllWorkloadsTab.value) {
      if (isAbnormalMode.value) {
        response = await k8sApi.getClusterAbnormalWorkloadList(selectedClusterId.value, {
          ...params,
          ...(keyword ? { keyword } : {})
        })
      } else if (keyword) {
        response = await k8sApi.searchClusterWorkloadList(selectedClusterId.value, {
          ...params,
          keyword
        })
      } else {
        response = await k8sApi.getClusterWorkloadList(selectedClusterId.value, params)
      }
    } else if (isAbnormalMode.value) {
      response = await k8sApi.getAbnormalWorkloadList(
        selectedClusterId.value,
        queryParams.namespace,
        {
          ...params,
          type: getRequestWorkloadType(),
          ...(keyword ? { keyword } : {})
        }
      )
    } else if (keyword) {
      response = await k8sApi.searchWorkloadList(
        selectedClusterId.value,
        queryParams.namespace,
        {
          ...params,
          type: getRequestWorkloadType(),
          keyword
        }
      )
    } else {
      response = await k8sApi.getWorkloadList(selectedClusterId.value, queryParams.namespace, {
        ...params,
        type: getRequestWorkloadType()
      })
    }
    
    const responseData = response.data || response
    console.log('工作负载查询API响应:', responseData)
    
    if (responseData.code === 200 || responseData.success) {
      if (isStandalonePodsTab.value) {
        const podData = responseData.data || {}
        const pods = Array.isArray(podData) ? podData : (podData.pods || [])
        const standalonePods = Array.isArray(pods) ? pods : []
        let filteredPods = standalonePods

        if (keyword) {
          const loweredKeyword = keyword.toLowerCase()
          filteredPods = filteredPods.filter(pod => String(pod.name || '').toLowerCase().includes(loweredKeyword))
        }

        if (isAbnormalMode.value) {
          filteredPods = filteredPods.filter(isStandalonePodAbnormal)
        }

        const serverTotal = filteredPods.length
        const serverPageSize = pageSize.value
        const serverTotalPages = serverTotal > 0 ? Math.ceil(serverTotal / serverPageSize) : 0

        if (serverTotalPages > 0 && currentPage.value > serverTotalPages) {
          currentPage.value = serverTotalPages
          await handleQuery()
          return
        }

        const startIndex = (currentPage.value - 1) * serverPageSize
        const pagedPods = filteredPods.slice(startIndex, startIndex + serverPageSize)

        tableData.value = pagedPods.map(mapStandalonePodToTableRow)
        total.value = serverTotal
      } else {
        const workloadData = responseData.data || {}
        const workloads = Array.isArray(workloadData) ? workloadData : (workloadData.workloads || [])
        // 确保workloads是数组
        const workloadList = Array.isArray(workloads) ? workloads : []
        const serverTotal = Array.isArray(workloadData)
          ? workloadList.length
          : Number(workloadData.total ?? workloadList.length)
        const serverPage = Array.isArray(workloadData)
          ? currentPage.value
          : Number(workloadData.page || currentPage.value)
        const serverPageSize = Array.isArray(workloadData)
          ? pageSize.value
          : Number(workloadData.pageSize || pageSize.value)
        const serverTotalPages = Array.isArray(workloadData)
          ? (serverTotal > 0 ? Math.ceil(serverTotal / serverPageSize) : 0)
          : Number(workloadData.totalPages || (serverTotal > 0 ? Math.ceil(serverTotal / serverPageSize) : 0))

        if (serverTotalPages > 0 && serverPage > serverTotalPages) {
          currentPage.value = serverTotalPages
          await handleQuery()
          return
        }

        tableData.value = workloadList.map(workload => ({
          id: workload.name,
          name: workload.name,
          type: workload.type?.toLowerCase() || workload.kind?.toLowerCase(),
          namespace: workload.namespace,
          replicas: `${workload.readyReplicas || 0}/${workload.replicas || 0}`,
          readyReplicas: workload.readyReplicas || 0,
          totalReplicas: workload.replicas || 0,
          images: workload.images || [],
          labels: workload.labels || {},
          status: workload.status || getWorkloadStatus(workload),
          age: formatAge(workload.createdAt),
          updateTime: workload.createdAt,
          updatedAt: workload.updatedAt,
          conditions: workload.conditions || [],
          resources: workload.resources || {
            cpu: { requests: '0', limits: '0' },
            memory: { requests: '0', limits: '0' }
          },
          rawData: workload
        }))

        total.value = serverTotal
        currentPage.value = serverPage
        pageSize.value = serverPageSize
      }
      
      console.log('工作负载列表加载成功:', tableData.value)
    } else {
      const errorMsg = responseData.message || '获取工作负载列表失败'

      // 特殊处理资源不存在的错误
      if (errorMsg.includes('the server could not find the requested resource')) {
        if (queryParams.type === 'cronjobs') {
          ElMessage.warning('当前集群不支持CronJob资源，可能是Kubernetes版本过低')
        } else if (queryParams.type && queryParams.type !== ALL_WORKLOAD_TYPE) {
          ElMessage.warning(`当前集群不支持${queryParams.type}资源类型`)
        } else {
          ElMessage.warning('请求的资源不存在，请检查集群配置')
        }
      } else {
        ElMessage.error(errorMsg)
      }

      resetTableData()
    }
  } catch (error) {
    console.error('获取工作负载列表失败:', error)
    
    if (error.code === 'ERR_NETWORK' || 
        error.message?.includes('ERR_CONNECTION_REFUSED') ||
        error.message?.includes('Failed to fetch')) {
      ElMessage.warning('后端服务连接失败，请检查服务状态')
    } else if (error.response?.status === 401) {
      ElMessage.error('认证失败，请重新登录')
    } else if (error.response?.status === 403) {
      ElMessage.error('权限不足，请联系管理员')
    } else {
      console.warn('API调用异常，但可能数据已正确加载')
    }
    
    resetTableData()
  } finally {
    loading.value = false
    console.log('✅ 工作负载查询完成，耗时:', Date.now() - queryStartTime + 'ms')
  }
}

// 获取工作负载状态
const getWorkloadStatus = (workload) => {
  // 如果后端直接返回了状态，优先使用
  if (workload.status) return workload.status
  
  const replicas = workload.replicas || 0
  const readyReplicas = workload.readyReplicas || 0
  
  if (workload.type === 'job' || workload.kind === 'Job') {
    return workload.succeeded ? 'Completed' : 
           workload.failed ? 'Failed' : 'Running'
  }
  
  if (workload.type === 'cronjob' || workload.kind === 'CronJob') {
    return workload.lastScheduleTime ? 'Active' : 'Suspended'
  }
  
  if (replicas === 0) return 'Stopped'
  if (readyReplicas === 0) return 'Pending'
  if (readyReplicas < replicas) return 'Partial'
  return 'Running'
}

// CPU资源相加
const addCpuResources = (a, b) => {
  const parseMillicores = (cpu) => {
    if (!cpu || cpu === '0') return 0
    if (cpu.endsWith('m')) return parseInt(cpu.replace('m', ''))
    return parseFloat(cpu) * 1000
  }
  
  const total = parseMillicores(a) + parseMillicores(b)
  return total < 1000 ? `${total}m` : `${(total / 1000).toFixed(1)}`
}

// 内存资源相加
const addMemoryResources = (a, b) => {
  const parseBytes = (memory) => {
    if (!memory || memory === '0') return 0
    if (memory.endsWith('Ki')) return parseInt(memory.replace('Ki', '')) * 1024
    if (memory.endsWith('Mi')) return parseInt(memory.replace('Mi', '')) * 1024 * 1024
    if (memory.endsWith('Gi')) return parseInt(memory.replace('Gi', '')) * 1024 * 1024 * 1024
    return parseInt(memory)
  }
  
  const totalBytes = parseBytes(a) + parseBytes(b)
  if (totalBytes < 1024 * 1024) return `${Math.round(totalBytes / 1024)}Ki`
  if (totalBytes < 1024 * 1024 * 1024) return `${Math.round(totalBytes / (1024 * 1024))}Mi`
  return `${Math.round(totalBytes / (1024 * 1024 * 1024))}Gi`
}

const formatAge = (createdTimestamp) => {
  if (!createdTimestamp) return 'Unknown'
  
  const now = new Date()
  const created = new Date(createdTimestamp)
  const diff = Math.floor((now - created) / 1000)
  
  if (diff < 60) return `${diff}s`
  if (diff < 3600) return `${Math.floor(diff / 60)}m`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h`
  return `${Math.floor(diff / 86400)}d`
}

const formatDateTime = (timestamp) => {
  if (!timestamp) return '-'
  
  const date = new Date(timestamp)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

const resetQuery = () => {
  const currentTab = normalizeWorkloadType(activeTab.value)

  queryParams.name = ''
  queryParams.scope = WORKLOAD_SCOPE_ALL
  queryParams.type = currentTab
  activeTab.value = currentTab
  currentPage.value = 1
  syncTabState(currentTab)
  handleQuery()
}

const handleScopeChange = (scope) => {
  queryParams.scope = scope
  currentPage.value = 1
  handleQuery()
}

const handleSearch = () => {
  queryParams.name = getKeyword()
  currentPage.value = 1
  handleQuery()
}

const handlePageChange = (page) => {
  currentPage.value = page
  handleQuery()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  handleQuery()
}

// 导航到监控仪表板
const navigateToMonitoring = () => {
  router.push('/k8s/monitoring')
}

// 创建工作负载相关函数
const showCreatePodDialog = () => {
  if (!selectedClusterId.value) {
    ElMessage.warning('请先选择集群')
    return
  }
  if (!queryParams.namespace) {
    ElMessage.warning('请先选择命名空间')
    return
  }

  createPodDialogVisible.value = true
}

// 处理YAML创建校验
const handlePodPreview = async (data) => {
  try {
    createPodDialogRef.value?.setLoading(true)

    // 使用validateYaml来校验YAML格式
    const response = await k8sApi.validateYaml(selectedClusterId.value, data.yamlContent)
    const responseData = response.data || response

    const result = {
      success: responseData.code === 200,
      message: responseData.message || (responseData.code === 200 ? '可以创建工作负载' : '创建预览失败'),
      details: responseData.data
    }

    createPodDialogRef.value?.setDryRunResult(result)

    if (result.success) {
      ElMessage.success(result.message)
    } else {
      ElMessage.error(result.message)
    }
  } catch (error) {
    const result = {
      success: false,
      message: error.message || '工作负载创建预览失败',
      details: null
    }
    createPodDialogRef.value?.setDryRunResult(result)
    ElMessage.error('工作负载创建预览失败: ' + (error.message || '网络错误'))
  } finally {
    createPodDialogRef.value?.setLoading(false)
  }
}

// 处理工作负载创建
const handlePodCreate = async (data) => {
  try {
    createPodDialogRef.value?.setLoading(true)

    // 使用createPodFromYaml来创建工作负载（该API支持多种资源类型）
    const response = await k8sApi.createPodFromYaml(selectedClusterId.value, queryParams.namespace, data)
    const responseData = response.data || response

    if (responseData.code === 200) {
      ElMessage.success('工作负载创建成功!')
      createPodDialogVisible.value = false
      handleQuery() // 刷新工作负载列表
    } else {
      ElMessage.error(responseData.message || '工作负载创建失败')
    }
  } catch (error) {
    console.error('工作负载创建失败:', error)
    ElMessage.error('工作负载创建失败: ' + (error.message || '网络错误'))
  } finally {
    createPodDialogRef.value?.setLoading(false)
  }
}


// 集群选择变化处理
const handleClusterChange = async () => {
  // 清空数据，NamespaceSelector 组件会自动处理命名空间列表加载
  resetTableData()
  currentPage.value = 1

  if (selectedClusterId.value && queryParams.namespace) {
    saveStoredClusterId(selectedClusterId.value)
    handleQuery()
    return
  }

  saveStoredClusterId(selectedClusterId.value)
}

// 命名空间选择变化处理
const handleNamespaceChange = () => {
  if (selectedClusterId.value && queryParams.namespace) {
    currentPage.value = 1
    handleQuery()
  } else {
    resetTableData()
  }
}

// 导航到容器详情页面
const navigateToPodDetail = async (row) => {
  try {
    const targetNamespace = getResourceNamespace(row)
    if (row.type === 'pod') {
      router.push({
        path: `/k8s/pod/${selectedClusterId.value}/${targetNamespace}/${row.name}`
      })
      return
    }

    console.log('🔍 点击工作负载名称，跳转到Pod详情:', row)
    console.log('📊 工作负载详情:', {
      name: row.name,
      type: row.type,
      namespace: targetNamespace
    })

    // 使用新的专门API获取该工作负载下的Pod列表
    const response = await k8sApi.getWorkloadPods(
      selectedClusterId.value,
      targetNamespace,
      row.type.toLowerCase(),
      row.name
    )
    const responseData = response.data || response

    if (responseData.code === 200 && responseData.data && responseData.data.length > 0) {
      // 获取第一个Pod（如果有多个Pod，跳转到第一个）
      const firstPod = responseData.data[0]
      console.log('🎯 跳转到第一个Pod:', firstPod.name)

      router.push({
        path: `/k8s/pod/${selectedClusterId.value}/${targetNamespace}/${firstPod.name}`
      })
    } else {
      ElMessage.warning('该工作负载下暂无Pod或Pod信息获取失败')
    }
  } catch (error) {
    console.error('获取工作负载Pod信息失败:', error)
    ElMessage.error('获取工作负载Pod信息失败，请检查网络连接')
  }
}


const handleMoreCommand = (row, command) => {
  switch (command) {
    case 'scale':
      scaleWorkload(row)
      break
    case 'restart':
      restartWorkload(row)
      break
    case 'config':
      updatePodConfig(row)
      break
    case 'schedule':
      updateScheduling(row)
      break
    case 'yaml':
      if (row.type === 'pod') {
        viewYaml(row.rawData || row)
      } else {
        editWorkloadYaml(row)
      }
      break
    case 'delete':
      deleteWorkload(row)
      break
  }
}

// 查看YAML
const viewYaml = async (pod) => {
  const targetNamespace = getResourceNamespace(pod)

  openYamlDialog({
    clusterId: selectedClusterId.value,
    resourceType: 'Pod',
    resourceName: pod.name,
    editable: false,
    fetchYaml: (clusterId, resourceName) => k8sApi.getPodYaml(clusterId, targetNamespace, resourceName),
    fallbackYaml: `apiVersion: v1\nkind: Pod\nmetadata:\n  name: ${pod.name}`,
    loadErrorMessage: '获取Pod YAML失败'
  })
}

// 扩容缩容
const scaleWorkload = (row) => {
  if (!canScale(row)) {
    ElMessage.warning('该工作负载不支持扩缩容操作')
    return
  }
  
  currentWorkload.value = row
  scaleForm.replicas = row.totalReplicas || 1
  scaleDialogVisible.value = true
}

// 提交扩容缩容
const submitScale = async () => {
  try {
    const response = await k8sApi.scaleDeployment(
      selectedClusterId.value, 
      getResourceNamespace(currentWorkload.value), 
      currentWorkload.value.name, 
      { replicas: scaleForm.replicas }
    )
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`${currentWorkload.value.name} 扩缩容成功`)
      scaleDialogVisible.value = false
      handleQuery()
    } else {
      ElMessage.error(responseData.message || '扩缩容失败')
    }
  } catch (error) {
    console.error('扩缩容失败:', error)
    ElMessage.error('扩缩容失败，请检查网络连接')
  }
}

// 重启工作负载
const restartWorkload = async (row) => {
  if (!canRestart(row)) {
    ElMessage.warning('该工作负载不支持重启操作')
    return
  }
  
  try {
    const targetNamespace = getResourceNamespace(row)
    await ElMessageBox.confirm(
      `确定要重启 ${row.type} "${row.name}" 吗？`,
      '重启确认',
      {
        confirmButtonText: '确定重启',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const response = await k8sApi.restartDeployment(selectedClusterId.value, targetNamespace, row.name)
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`${row.name} 重启成功`)
      handleQuery()
    } else {
      ElMessage.error(responseData.message || '重启失败')
    }
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消重启操作')
    } else {
      console.error('重启失败:', error)
      ElMessage.error('重启失败，请检查网络连接')
    }
  }
}

const getWorkloadTypeName = (type) => {
  const nameMap = {
    'deployment': 'Deployment',
    'statefulset': 'StatefulSet',
    'daemonset': 'DaemonSet',
    'job': 'Job',
    'cronjob': 'CronJob',
    'pod': 'Pod'
  }
  return nameMap[type] || type
}

// 复制内容到剪贴板
const copyToClipboard = async (text, successMessage = '已复制到剪贴板') => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(successMessage)
  } catch (error) {
    console.error('复制失败:', error)
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()
    try {
      document.execCommand('copy')
      ElMessage.success(successMessage)
    } catch (fallbackError) {
      ElMessage.error('复制失败，请手动复制')
    }
    document.body.removeChild(textArea)
  }
}

// 判断是否为系统命名空间或系统工作负载
const isSystemWorkload = (workload) => {
  // 系统命名空间
  const systemNamespaces = ['kube-system', 'kube-public', 'kube-node-lease', 'calico-system', 'tigera-operator']
  if (systemNamespaces.includes(workload.namespace)) {
    return true
  }
  
  // 系统工作负载名称前缀
  const systemPrefixes = ['kube-', 'calico-', 'coredns', 'metrics-server', 'node-local-dns', 'kubernetes-dashboard']
  if (systemPrefixes.some(prefix => workload.name.startsWith(prefix))) {
    return true
  }
  
  return false
}

// 判断是否可以扩缩容
const canScale = (workload) => {
  // 只有 Deployment 和 StatefulSet 支持扩缩容
  if (!['deployment', 'statefulset'].includes(workload.type)) {
    return false
  }
  
  // 系统工作负载不允许扩缩容
  return !isSystemWorkload(workload)
}

// 判断是否可以重启
const canRestart = (workload) => {
  // 只有 Deployment 支持重启
  if (workload.type !== 'deployment') {
    return false
  }
  
  // 系统工作负载不允许重启
  return !isSystemWorkload(workload)
}

// 判断是否可以删除
const canDelete = (workload) => {
  // 系统工作负载不允许删除
  return !isSystemWorkload(workload)
}

// 判断是否可以更新Pod配置
const canUpdatePodConfig = (workload) => {
  // 只有 Deployment 和 StatefulSet 支持Pod配置更新
  if (!['deployment', 'statefulset'].includes(workload.type)) {
    return false
  }
  
  // 系统工作负载不允许配置更新
  return !isSystemWorkload(workload)
}

// 判断是否可以更新调度
const canUpdateScheduling = (workload) => {
  // 只有 Deployment、StatefulSet 支持调度更新
  if (!['deployment', 'statefulset'].includes(workload.type)) {
    return false
  }
  
  // 系统工作负载不允许调度更新
  return !isSystemWorkload(workload)
}

// 判断是否可以编辑YAML
const canEditYaml = (workload) => {
  if (workload.type === 'pod') {
    return true
  }
  // 系统工作负载不允许编辑YAML
  return !isSystemWorkload(workload)
}

// 获取可见标签数量（排除系统标签）
const getVisibleLabelCount = (labels) => {
  if (!labels) return 0
  
  const systemLabelPrefixes = [
    'kubernetes.io/',
    'beta.kubernetes.io/',
    'node-role.kubernetes.io/',
    'node.kubernetes.io/',
    'app.kubernetes.io/managed-by',
    'pod-template-hash'
  ]
  
  return Object.keys(labels).filter(key => 
    !systemLabelPrefixes.some(prefix => key.startsWith(prefix))
  ).length
}

// 根据副本数获取Pod状态标签类型
const getPodStatusTagByReplicas = (readyReplicas, totalReplicas) => {
  if (totalReplicas === 0) return 'info'
  if (readyReplicas === 0) return 'danger'
  if (readyReplicas < totalReplicas) return 'warning'
  return 'success'
}

// 获取副本状态文本
const getReplicaStatusText = (readyReplicas, totalReplicas) => {
  if (totalReplicas === 0) return '已停止'
  if (readyReplicas === 0) return '启动中'
  if (readyReplicas < totalReplicas) return '部分就绪'
  return '全部就绪'
}

// 获取副本状态样式类
const getReplicaStatusClass = (readyReplicas, totalReplicas) => {
  if (totalReplicas === 0) return 'status-stopped'
  if (readyReplicas === 0) return 'status-starting'
  if (readyReplicas < totalReplicas) return 'status-partial'
  return 'status-ready'
}

// 查看工作负载标签
const viewWorkloadLabels = (row) => {
  currentWorkload.value = row
  workloadLabelsDialogVisible.value = true
}

// 查看所有镜像
const viewAllImages = (row) => {
  currentWorkload.value = row
  allImagesDialogVisible.value = true
}

// 更新Pod配置
const updatePodConfig = (row) => {
  if (!canUpdatePodConfig(row)) {
    ElMessage.warning('该工作负载不支持Pod配置更新')
    return
  }
  
  currentWorkload.value = row
  podConfigDialogVisible.value = true
}

// 更新调度
const updateScheduling = (row) => {
  if (!canUpdateScheduling(row)) {
    ElMessage.warning('该工作负载不支持调度更新')
    return
  }
  
  currentWorkload.value = row
  schedulingDialogVisible.value = true
}

// 编辑工作负载YAML
const editWorkloadYaml = async (row) => {
  if (!canEditYaml(row)) {
    ElMessage.warning('系统工作负载不允许编辑YAML')
    return
  }

  currentWorkload.value = row
  const targetNamespace = getResourceNamespace(row)

  openYamlDialog({
    clusterId: selectedClusterId.value,
    resourceType: row.type || 'Workload',
    resourceName: row.name,
    editable: true,
    fetchYaml: async (clusterId, resourceName) => {
      console.log('🔍 开始获取工作负载YAML...', row)

      try {
        const response = await k8sApi.getWorkloadYaml(clusterId, targetNamespace, row.type, resourceName)
        const responseData = response.data || response

        if (responseData.code === 200 || responseData.success) {
          return response
        }
      } catch (workloadError) {
        console.log('工作负载YAML获取失败，尝试通过Pod获取:', workloadError)
      }

      console.log('🔍 通过获取工作负载Pod列表来获取Pod YAML...')
      const detailResponse = await k8sApi.getWorkloadPods(clusterId, targetNamespace, row.type.toLowerCase(), resourceName)
      const detailData = detailResponse.data || detailResponse

      if ((detailData.code === 200 || detailData.success) && detailData.data && detailData.data.length > 0) {
        const firstPod = detailData.data[0]
        console.log('🔍 获取Pod YAML:', firstPod.name)

        const podYamlResponse = await k8sApi.getPodYaml(clusterId, targetNamespace, firstPod.name)
        const podYamlData = podYamlResponse.data || podYamlResponse

        if (podYamlData.code === 200 || podYamlData.success) {
          return {
            code: 200,
            data: {
              yamlContent: getYamlContentFromApiResponse(
                podYamlData,
                `# Pod YAML for ${firstPod.name}\n# 工作负载: ${resourceName} (${row.type})\n`
              )
            }
          }
        }

        throw new Error(podYamlData.message || 'Pod YAML获取失败')
      }

      console.log('⚠️ 没有找到Pod，生成基础YAML模板')
      ElMessage.warning('未找到实际YAML内容，显示基础模板')
      return {
        code: 200,
        data: {
          yamlContent: generateWorkloadYamlTemplate(row)
        }
      }
    },
    updateYaml: async (clusterId, resourceName, yamlContent) => {
      const { workloadType, workloadName } = parseWorkloadFromYaml(yamlContent)
      return k8sApi.updateWorkloadYaml(
        clusterId,
        targetNamespace,
        workloadType,
        workloadName || resourceName,
        yamlContent
      )
    },
    fallbackYaml: generateWorkloadYamlTemplate(row),
    loadErrorMessage: '获取工作负载YAML失败',
    saveSuccessMessage: `${row.name} YAML配置保存成功`,
    saveErrorMessage: 'YAML配置保存失败',
    onSaved: handleQuery
  })
}

// 生成工作负载YAML模板
const generateWorkloadYamlTemplate = (workload) => {
  const kind = workload.type.charAt(0).toUpperCase() + workload.type.slice(1)
  const targetNamespace = getResourceNamespace(workload)
  return `apiVersion: apps/v1
kind: ${kind}
metadata:
  name: ${workload.name}
  namespace: ${targetNamespace}
  labels:
    app: ${workload.name}
spec:
  replicas: ${workload.totalReplicas || 1}
  selector:
    matchLabels:
      app: ${workload.name}
  template:
    metadata:
      labels:
        app: ${workload.name}
    spec:
      containers:
      - name: ${workload.name}
        image: nginx:latest
        ports:
        - containerPort: 80
---
# 注意: 这是一个基础模板，请根据实际需求修改
# 工作负载类型: ${workload.type}
# 当前状态: ${workload.status}
# 副本数: ${workload.replicas}`
}

onMounted(async () => {
  try {
    console.log('🚀 开始加载k8s工作负载页面')
    const startTime = Date.now()
    restoreTabState()
    
    // 加载集群列表
    console.log('📡 正在加载集群列表...')
    await fetchClusterList()
    console.log('✅ 集群列表加载完成，耗时:', Date.now() - startTime + 'ms')
    
    // 如果有选中的集群，立即开始加载数据
    if (selectedClusterId.value) {
      console.log('🔄 开始并行加载命名空间和工作负载数据')
      
      // 命名空间加载已移至 NamespaceSelector 组件
      
      // 如果有默认命名空间，立即开始查询工作负载
      if (queryParams.namespace) {
        console.log('📦 立即开始查询工作负载:', queryParams.namespace)
        // 不等待命名空间，直接查询工作负载
        handleQuery().catch(error => {
          console.error('工作负载初始查询失败:', error)
        })
      }
    }
    
    console.log('🎉 页面初始化完成，总耗时:', Date.now() - startTime + 'ms')
  } catch (error) {
    console.error('页面初始化失败:', error)
  }
})

// 辅助函数
const formatCpu = (cpuStr) => {
  if (!cpuStr || cpuStr === '0' || cpuStr === '') return '-'
  return cpuStr
}

const formatMemory = (memoryStr) => {
  if (!memoryStr || memoryStr === '0' || memoryStr === '') return '-'

  if (memoryStr.endsWith('Ki')) {
    const kb = parseInt(memoryStr.replace('Ki', ''))
    if (kb < 1024) return memoryStr
    const mb = (kb / 1024).toFixed(1)
    return `${mb}Mi`
  }

  if (memoryStr.endsWith('Mi')) {
    const mb = parseInt(memoryStr.replace('Mi', ''))
    if (mb < 1024) return memoryStr
    const gb = (mb / 1024).toFixed(1)
    return `${gb}Gi`
  }

  if (memoryStr.endsWith('Gi')) {
    return memoryStr
  }

  // 如果没有单位，假设是字节
  const bytes = parseInt(memoryStr)
  if (!isNaN(bytes) && bytes > 0) {
    if (bytes < 1024) return `${bytes}B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)}Ki`
    if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)}Mi`
    return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)}Gi`
  }

  return memoryStr || '-'
}

const getImageTag = (image) => {
  if (!image) return 'latest'
  const parts = image.split(':')
  return parts.length > 1 ? parts[parts.length - 1] : 'latest'
}

const getImageRegistry = (image) => {
  if (!image) return 'docker.io'
  const parts = image.split('/')
  if (parts.length === 1) return 'docker.io'
  if (parts[0].includes('.') || parts[0].includes(':')) {
    return parts[0]
  }
  return 'docker.io'
}

// 获取用户自定义标签
const getUserLabels = (labels) => {
  if (!labels) return {}

  const systemLabelPrefixes = [
    'kubernetes.io/',
    'beta.kubernetes.io/',
    'node-role.kubernetes.io/',
    'node.kubernetes.io/',
    'app.kubernetes.io/managed-by',
    'pod-template-hash'
  ]

  const userLabels = {}
  Object.entries(labels).forEach(([key, value]) => {
    if (!systemLabelPrefixes.some(prefix => key.startsWith(prefix))) {
      userLabels[key] = value
    }
  })

  return userLabels
}

// 获取系统标签
const getSystemLabels = (labels) => {
  if (!labels) return {}

  const systemLabelPrefixes = [
    'kubernetes.io/',
    'beta.kubernetes.io/',
    'node-role.kubernetes.io/',
    'node.kubernetes.io/',
    'app.kubernetes.io/managed-by',
    'pod-template-hash'
  ]

  const systemLabels = {}
  Object.entries(labels).forEach(([key, value]) => {
    if (systemLabelPrefixes.some(prefix => key.startsWith(prefix))) {
      systemLabels[key] = value
    }
  })

  return systemLabels
}

const toLabelTableData = (labels) =>
  Object.entries(labels || {}).map(([key, value]) => ({ key, value }))

// 处理Pod配置提交
const handlePodConfigSubmit = async (updateData) => {
  console.log('🔧 开始处理Pod配置更新:', {
    clusterId: selectedClusterId.value,
    namespace: getResourceNamespace(currentWorkload.value),
    workloadName: currentWorkload.value.name,
    updateData: updateData
  })

  try {
    const response = await k8sApi.updateDeployment(
      selectedClusterId.value,
      getResourceNamespace(currentWorkload.value),
      currentWorkload.value.name,
      updateData
    )

    const responseData = response.data || response
    console.log('📤 API响应:', responseData)

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`${currentWorkload.value.name} Pod配置更新成功`)
      podConfigDialogVisible.value = false
      handleQuery() // 刷新列表
    } else {
      console.warn('❌ API返回错误:', responseData)
      ElMessage.error(responseData.message || 'Pod配置更新失败')
    }
  } catch (error) {
    console.error('💥 Pod配置更新异常:', error)
    ElMessage.error(`Pod配置更新失败: ${error.message || '请检查网络连接'}`)
  }
}

// 提交调度配置
const submitScheduling = async () => {
  try {
    const updateData = {
      template: {
        nodeSelector: currentWorkload.value.nodeSelector || {},
        tolerations: currentWorkload.value.tolerations || []
      }
    }

    const response = await k8sApi.updateDeployment(
      selectedClusterId.value,
      getResourceNamespace(currentWorkload.value),
      currentWorkload.value.name,
      updateData
    )

    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`${currentWorkload.value.name} 调度配置更新成功`)
      schedulingDialogVisible.value = false
      handleQuery() // 刷新列表
    } else {
      ElMessage.error(responseData.message || '调度配置更新失败')
    }
  } catch (error) {
    console.error('调度配置更新失败:', error)
    ElMessage.error('调度配置更新失败，请检查网络连接')
  }
}

// 从YAML内容中解析工作负载信息
const parseWorkloadFromYaml = (yamlContent) => {
  try {
    const yamlObj = yaml.load(yamlContent)
    if (!yamlObj || typeof yamlObj !== 'object') {
      throw new Error('无效的YAML格式')
    }

    const workloadType = yamlObj.kind?.toLowerCase()
    const workloadName = yamlObj.metadata?.name

    if (!workloadType || !workloadName) {
      throw new Error('YAML中缺少kind或metadata.name字段')
    }

    return { workloadType, workloadName }
  } catch (error) {
    console.error('解析YAML失败:', error)
    throw error
  }
}

const deleteWorkload = async (row) => {
  if (!canDelete(row)) {
    ElMessage.warning('系统工作负载不允许删除')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除 ${getWorkloadTypeName(row.type)} "${row.name}" 吗？此操作不可恢复！`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )

    console.log('🗑️ 开始删除工作负载:', {
      clusterId: selectedClusterId.value,
      namespace: getResourceNamespace(row),
      workloadType: row.type.toLowerCase(),
      workloadName: row.name
    })

    let response

    // 根据类型选择不同的删除API
    switch (row.type.toLowerCase()) {
      case 'pod':
        // Pod使用专门的删除Pod API
        response = await k8sApi.deletePod(selectedClusterId.value, getResourceNamespace(row), row.name)
        break
      case 'deployment':
        // Deployment使用专门的删除Deployment API
        response = await k8sApi.deleteDeployment(selectedClusterId.value, getResourceNamespace(row), row.name)
        break
      default:
        // 其他工作负载使用通用工作负载API
        response = await k8sApi.deleteWorkload(
          selectedClusterId.value,
          getResourceNamespace(row),
          row.type.toLowerCase(),
          row.name
        )
        break
    }

    console.log('📤 删除API响应:', response)

    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`${getWorkloadTypeName(row.type)} "${row.name}" 删除成功`)
      // 刷新工作负载列表
      await handleQuery()
    } else {
      ElMessage.error(responseData.message || `删除 ${getWorkloadTypeName(row.type)} 失败`)
    }
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消删除操作')
    } else {
      console.error('删除工作负载失败:', error)
      ElMessage.error(`删除失败: ${error.message || '请检查网络连接'}`)
    }
  }
}
</script>

<template>
  <div class="k8s-workloads-management">
    <el-card shadow="hover" class="workloads-card">
      <template #header>
        <div class="card-header">
          <span class="title">K8s 工作负载管理</span>
          <div class="header-actions">
            <ClusterSelector
              v-model="selectedClusterId"
              width="240px"
              size="default"
              margin-right="0"
              theme="cluster"
              @change="handleClusterChange"
            />
            
            <NamespaceSelector
              v-model="queryParams.namespace"
              :cluster-id="selectedClusterId"
              filterable
              width="260px"
              size="default"
              theme="namespace"
              @change="handleNamespaceChange"
            />
          </div>
        </div>
      </template>
      
        <!-- 搜索表单 -->
        <div class="search-section">
          <el-form :inline="true" :model="queryParams" class="search-form">
          <el-form-item :label="workloadNameLabel">
            <el-input
              v-model="queryParams.name"
              :placeholder="workloadNamePlaceholder"
              clearable
              size="default"
              class="search-name-input"
              style="width: 250px"
              @keyup.enter="handleSearch"
              @clear="handleSearch"
            >
              <template #prefix>
                <el-icon class="search-name-icon"><Search /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="展示范围">
            <el-select
              v-model="queryParams.scope"
              placeholder="请选择展示范围"
              size="default"
              class="scope-select"
              style="width: 220px"
              @change="handleScopeChange"
            >
              <template #prefix>
                <el-icon class="scope-select-icon"><Warning /></el-icon>
              </template>
              <el-option
                v-for="option in workloadScopeOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :icon="Search" size="small" @click="handleSearch">
              搜索
            </el-button>
            <el-button type="warning" :icon="Refresh" size="small" @click="resetQuery">
              刷新
            </el-button>
            <el-button :icon="Monitor" type="success" size="small" @click="navigateToMonitoring">
              监控仪表板
            </el-button>
            <el-button
              :icon="Plus"
              v-authority="['k8s:workload:add']"
              type="primary"
              size="small"
              @click="showCreatePodDialog"
            >
              创建工作负载
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 工作负载类型标签页 -->
        <div class="workload-type-section">
          <el-tabs v-model="activeTab" @tab-change="handleTabChange">
            <el-tab-pane
              v-for="option in workloadTypeOptions"
              :key="option.value"
              :label="option.label"
              :name="option.value"
            />
          </el-tabs>
        </div>
      </div>

      <!-- 工作负载列表表格 -->
      <div class="table-section">
        <el-table
          :data="tableData"
          v-loading="loading"
          :empty-text="workloadTableEmptyText"
          stripe
          style="width: 100%"
          class="workloads-table"
        >
          <el-table-column prop="name" label="名称" min-width="200">
            <template #default="{ row }">
              <div class="workload-name-container">
                <img src="@/assets/image/k8s.svg" alt="k8s" class="k8s-icon" />
                <div class="workload-info">
                  <div 
                    class="workload-name"
                  >
                    {{ row.name }}
                  </div>
                  <span
                    class="workload-type-label"
                  >
                    {{ getWorkloadTypeName(row.type) }}
                  </span>
                  <span
                    v-if="isAllWorkloadsTab"
                    class="workload-namespace-label"
                  >
                    {{ row.namespace || '-' }}
                  </span>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="标签" min-width="100" align="center">
            <template #default="{ row }">
              <div class="label-container">
                <el-badge :value="getVisibleLabelCount(row.labels)" :max="99" class="label-badge">
                  <el-button
                    type="text"
                    size="small"
                    circle
                    @click="viewWorkloadLabels(row)"
                    class="label-icon-button"
                  >
                    <img src="@/assets/image/标签.svg" alt="标签" width="14" height="14" />
                  </el-button>
                </el-badge>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column :label="replicaColumnLabel" min-width="120" align="center">
            <template #default="{ row }">
              <div class="pod-status-container">
                <el-tag 
                  :type="getPodStatusTagByReplicas(row.readyReplicas, row.totalReplicas)" 
                  size="default"
                  class="pod-count-tag"
                >
                  <el-icon class="pod-icon"><Monitor /></el-icon>
                  {{ row.replicas }}
                </el-tag>
                <div class="pod-status-text">
                  <span :class="getReplicaStatusClass(row.readyReplicas, row.totalReplicas)">
                    {{ getReplicaStatusText(row.readyReplicas, row.totalReplicas) }}
                  </span>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="Request/Limits" min-width="160">
            <template #default="{ row }">
              <div class="resource-info">
                <div class="resource-row">
                  <span class="resource-type">CPU:</span>
                  <span class="resource-values">
                    <span class="request-value">{{ formatCpu(row.resources?.requests?.cpu) }}</span>
                    <span class="separator">/</span>
                    <span class="limit-value">{{ formatCpu(row.resources?.limits?.cpu) }}</span>
                  </span>
                </div>
                <div class="resource-row">
                  <span class="resource-type">Mem:</span>
                  <span class="resource-values">
                    <span class="request-value">{{ formatMemory(row.resources?.requests?.memory) }}</span>
                    <span class="separator">/</span>
                    <span class="limit-value">{{ formatMemory(row.resources?.limits?.memory) }}</span>
                  </span>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="镜像" min-width="250">
            <template #default="{ row }">
              <div class="images-list">
                <div
                  v-for="(image, index) in row.images.slice(0, 1)"
                  :key="index"
                  class="image-tag-wrapper"
                  @click="copyToClipboard(image, '镜像地址已复制')"
                >
                  <el-icon class="copy-icon"><DocumentCopy /></el-icon>
                  <span class="full-image-name">{{ image }}</span>
                </div>
                <el-button
                  v-if="row.images.length > 1"
                  type="text"
                  size="small"
                  class="more-images-btn"
                  @click="viewAllImages(row)"
                >
                  +{{ row.images.length - 1 }}个镜像
                </el-button>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="创建时间" min-width="150">
            <template #default="{ row }">
              <div class="time-info">
                <span class="datetime-text">{{ formatDateTime(row.updateTime) }}</span>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="更新时间" min-width="150">
            <template #default="{ row }">
              <div class="time-info">
                <span v-if="row.updatedAt" class="datetime-text">
                  {{ formatDateTime(row.updatedAt) }}
                </span>
                <span v-else class="no-update">-</span>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="190" fixed="right">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-button
                  type="primary"
                  :icon="View"
                  size="small"
                  @click="navigateToPodDetail(row)"
                >详情</el-button>

                <el-dropdown @command="(cmd) => handleMoreCommand(row, cmd)" trigger="click">
                  <el-button type="warning" size="small" v-authority="['k8s:workload:actions']">
                    更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item
                        command="scale"
                        :icon="Plus"
                        v-authority="['k8s:workload:expandable']"
                        :disabled="!canScale(row)"
                      >伸缩</el-dropdown-item>
                      <el-dropdown-item
                        command="restart"
                        :icon="Refresh"
                        v-authority="['k8s:workload:restart']"
                        :disabled="!canRestart(row)"
                      >重构</el-dropdown-item>
                      <el-dropdown-item
                        command="config"
                        :icon="Setting"
                        v-authority="['k8s:workload:resource']"
                        :disabled="!canUpdatePodConfig(row)"
                      >Pod配置</el-dropdown-item>
                      <el-dropdown-item
                        command="schedule"
                        :icon="Monitor"
                        v-authority="['k8s:workload:dispatch']"
                        :disabled="!canUpdateScheduling(row)"
                      >调度更新</el-dropdown-item>
                      <el-dropdown-item
                        command="yaml"
                        :icon="Document"
                        v-authority="['k8s:workload:edityaml']"
                        :disabled="!canEditYaml(row)"
                      >YAML编辑</el-dropdown-item>
                      <el-dropdown-item
                        command="delete"
                        :icon="Delete"
                        v-authority="['k8s:workload:delete']"
                        :disabled="!canDelete(row)"
                        divided
                      >删除</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="pagination-wrapper" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>

    <!-- 扩缩容对话框 -->
    <el-dialog
      v-model="scaleDialogVisible"
      :title="`扩缩容 - ${currentWorkload.name || ''}`"
      width="400px"
      class="scale-dialog"
    >
      <el-form :model="scaleForm" label-width="80px">
        <el-form-item label="副本数" required>
          <el-input-number
            v-model="scaleForm.replicas"
            :min="0"
            :max="100"
            style="width: 100%"
          />
          <div class="form-tip">当前副本数: {{ currentWorkload.totalReplicas }}</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="scaleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitScale">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 工作负载标签查看对话框 -->
    <el-dialog
      v-model="workloadLabelsDialogVisible"
      :title="`工作负载标签 - ${currentWorkload.name || ''}`"
      width="800px"
      class="workload-labels-view-dialog"
    >
      <div class="labels-view-content" v-if="currentWorkload.labels">
        <!-- 用户自定义标签 -->
        <div class="labels-section" v-if="Object.keys(getUserLabels(currentWorkload.labels)).length > 0">
          <h4 class="section-title">用户标签</h4>
          <div class="labels-table-wrapper">
            <el-table
              :data="toLabelTableData(getUserLabels(currentWorkload.labels))"
              border
              stripe
              size="small"
              max-height="300"
              class="labels-table"
            >
              <el-table-column prop="key" label="Key" width="250" show-overflow-tooltip />
              <el-table-column prop="value" label="Value" show-overflow-tooltip />
            </el-table>
          </div>
        </div>

        <!-- 系统标签 -->
        <div class="labels-section" v-if="Object.keys(getSystemLabels(currentWorkload.labels)).length > 0">
          <h4 class="section-title">系统标签</h4>
          <div class="labels-table-wrapper">
            <el-table
              :data="toLabelTableData(getSystemLabels(currentWorkload.labels))"
              border
              stripe
              size="small"
              max-height="300"
              class="labels-table"
            >
              <el-table-column prop="key" label="Key" width="250" show-overflow-tooltip />
              <el-table-column prop="value" label="Value" show-overflow-tooltip />
            </el-table>
          </div>
        </div>

        <!-- 没有标签的提示 -->
        <div v-if="!currentWorkload.labels || Object.keys(currentWorkload.labels).length === 0" class="no-labels">
          <el-empty description="该工作负载没有标签" :image-size="60" />
        </div>
      </div>
    </el-dialog>

    <!-- 所有镜像查看对话框 -->
    <el-dialog
      v-model="allImagesDialogVisible"
      :title="`镜像列表 - ${currentWorkload.name || ''}`"
      width="800px"
      class="all-images-dialog"
    >
      <div class="images-view-content" v-if="currentWorkload.images">
        <div class="images-section">
          <h4>容器镜像 ({{ currentWorkload.images?.length || 0 }}个)</h4>
          <div class="all-images-list">
            <el-card
              v-for="(image, index) in currentWorkload.images"
              :key="index"
              class="image-card"
              shadow="hover"
            >
              <div class="image-info">
                <div class="image-name">
                  <el-icon class="image-icon"><Connection /></el-icon>
                  <span class="full-image-name">{{ image }}</span>
                </div>
                <div class="image-actions">
                  <el-button
                    type="primary"
                    size="small"
                    :icon="DocumentCopy"
                    @click="copyToClipboard(image, '镜像地址已复制')"
                  >
                    复制
                  </el-button>
                </div>
              </div>
              <div class="image-details">
                <el-tag size="small" type="info">{{ getImageTag(image) }}</el-tag>
                <el-tag size="small" type="success">{{ getImageRegistry(image) }}</el-tag>
              </div>
            </el-card>
          </div>
        </div>
        
        <!-- 没有镜像的提示 -->
        <div v-if="!currentWorkload.images || currentWorkload.images.length === 0" class="no-images">
          <el-empty description="该工作负载没有镜像" :image-size="60" />
        </div>
      </div>
    </el-dialog>

    <!-- Pod配置更新对话框 -->
    <PodConfigDialog
      :visible="podConfigDialogVisible"
      :workload="currentWorkload"
      @update:visible="podConfigDialogVisible = $event"
      @close="podConfigDialogVisible = false"
      @submit="handlePodConfigSubmit"
    />

    <!-- 调度更新对话框 -->
    <el-dialog
      v-model="schedulingDialogVisible"
      :title="`更新调度 - ${currentWorkload.name || ''}`"
      width="500px"
      class="scheduling-dialog"
    >
      <el-form :model="currentWorkload" label-width="120px">
        <el-form-item label="节点选择器">
          <el-input
            v-model="currentWorkload.nodeSelector"
            placeholder="如: kubernetes.io/arch=amd64"
            style="width: 100%"
          />
          <div class="form-tip">格式: key=value，多个用逗号分隔</div>
        </el-form-item>
        <el-form-item label="节点亲和性">
          <el-select
            v-model="currentWorkload.nodeAffinity"
            placeholder="请选择节点亲和性"
            style="width: 100%"
          >
            <el-option label="无要求" value="none" />
            <el-option label="偏好调度" value="preferred" />
            <el-option label="必须调度" value="required" />
          </el-select>
        </el-form-item>
        <el-form-item label="Pod反亲和性">
          <el-switch
            v-model="currentWorkload.podAntiAffinity"
            active-text="启用"
            inactive-text="禁用"
          />
          <div class="form-tip">避免Pod调度到同一节点</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="schedulingDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitScheduling">更新调度</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 创建工作负载对话框 -->
    <CreatePodDialog
      ref="createPodDialogRef"
      :visible="createPodDialogVisible"
      :cluster-id="selectedClusterId"
      :cluster-name="clusterList.find(c => c.id === selectedClusterId)?.name"
      :namespace="queryParams.namespace"
      @update:visible="createPodDialogVisible = $event"
      @close="createPodDialogVisible = false"
      @preview="handlePodPreview"
      @create="handlePodCreate"
    />
  </div>
</template>


<style scoped>
.k8s-workloads-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.workloads-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #6b9bd1;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-section {
  margin-bottom: 24px;
  padding: 20px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

.search-name-icon {
  color: #409eff;
  font-size: 16px;
}

.search-name-input :deep(.el-input__wrapper) {
  min-height: 40px;
  padding: 0 12px;
  background: rgba(64, 158, 255, 0.06);
  border-color: rgba(64, 158, 255, 0.22);
}

.search-name-input :deep(.el-input__wrapper:hover) {
  border-color: #409eff;
}

.search-name-input :deep(.el-input__wrapper.is-focus) {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.15);
}

.search-name-input :deep(.el-input__wrapper.is-disabled) {
  background: rgba(245, 247, 250, 0.95);
  border-color: rgba(230, 162, 60, 0.18);
  box-shadow: none;
}

.scope-select-icon {
  color: #e6a23c;
}

.scope-select :deep(.el-input__wrapper) {
  min-height: 40px;
  padding: 0 12px;
  background: rgba(230, 162, 60, 0.08);
  border-color: rgba(230, 162, 60, 0.25);
}

.scope-select :deep(.el-input__wrapper:hover) {
  border-color: #e6a23c;
}

.scope-select :deep(.el-input__wrapper.is-focus) {
  border-color: #e6a23c;
  box-shadow: 0 0 0 2px rgba(230, 162, 60, 0.18);
}

.search-mode-hint {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(230, 162, 60, 0.12);
  border: 1px solid rgba(230, 162, 60, 0.2);
  color: #b88230;
  font-size: 12px;
  font-weight: 500;
}

/* 工作负载类型标签页样式 */
.workload-type-section {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid rgba(103, 126, 234, 0.1);
}

.table-section {
  margin-top: 20px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
  padding: 12px 0;
}

.workloads-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.workloads-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.workloads-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #000000 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: none;
}

.workloads-table :deep(.el-table__header th .cell) {
  color: #000000 !important;
  font-weight: 700 !important;
}

.workloads-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.workloads-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.workload-name-container {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: default;
}

.workload-name-container:hover {
  transform: none !important;
  background-color: transparent !important;
}

.workload-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.workload-name {
  font-weight: 600;
  color: #6b9bd1;
  font-size: 14px;
}

.workload-type-label {
  font-size: 12px;
  color: #E6A23C;
  font-weight: 500;
  pointer-events: none;
  user-select: none;
}

.workload-namespace-label {
  font-size: 12px;
  color: #909399;
  font-weight: 500;
}

.workload-type-tag {
  font-size: 11px;
  height: 18px;
  line-height: 16px;
  padding: 0 6px;
}

.pod-name-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.k8s-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.workload-name-link {
  font-weight: 600;
  color: #667eea;
  text-decoration: none;
  transition: all 0.3s ease;
}

.workload-name-link:hover {
  color: #764ba2;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.pod-name {
  font-weight: 500;
  color: #6b9bd1;
}

.resource-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.resource-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

.resource-type {
  font-size: 12px;
  color: #909399;
  min-width: 35px;
  font-weight: 500;
}

.resource-values {
  display: flex;
  align-items: center;
  gap: 2px;
}

.request-value {
  font-size: 12px;
  color: #67c23a;
  font-weight: 500;
}

.separator {
  font-size: 12px;
  color: #dcdfe6;
  margin: 0 2px;
}

.limit-value {
  font-size: 12px;
  color: #e6a23c;
  font-weight: 500;
}

.resource-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.resource-label {
  font-size: 12px;
  color: #909399;
  min-width: 55px;
}

.resource-value {
  font-size: 12px;
  color: #606266;
  font-weight: 500;
}

.ip-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.ip-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.ip-label {
  font-size: 12px;
  color: #909399;
  min-width: 35px;
}

.ip-value {
  font-size: 12px;
  color: #606266;
  font-weight: 500;
}

.images-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.image-tag {
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 4px;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.image-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.copy-icon {
  font-size: 10px;
}

.more-images {
  cursor: default;
}

/* Pod状态样式 */
.pod-status-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.pod-count-tag {
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  font-weight: 600;
}

.pod-icon {
  font-size: 14px;
}

.pod-status-text {
  font-size: 11px;
  line-height: 1.2;
}

.status-ready {
  color: #67c23a;
  font-weight: 500;
}

.status-partial {
  color: #e6a23c;
  font-weight: 500;
}

.status-starting {
  color: #f56c6c;
  font-weight: 500;
}

.status-stopped {
  color: #909399;
  font-weight: 500;
}

/* 标签容器样式 */
.label-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  padding: 4px 0;
}

.label-badge {
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.label-icon-button {
  background: transparent;
  border: none;
  color: #606266;
  transition: all 0.3s ease;
}

.label-icon-button:hover {
  background: transparent;
  color: #409eff;
  transform: scale(1.1);
}

/* 镜像显示优化 */
.image-tag-wrapper {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 4px;
}

.image-tag-wrapper:hover {
  transform: translateY(-1px);
}

.image-tag-wrapper .full-image-name {
  font-family: 'Monaco', 'Courier New', monospace;
  font-size: 11px;
  color: #6b9bd1;
  word-break: break-all;
  line-height: 1.4;
  white-space: normal;
}

.image-tag-wrapper .copy-icon {
  color: #666;
  font-size: 12px;
  flex-shrink: 0;
}

.more-images-btn {
  color: #409eff;
  font-size: 12px;
  padding: 2px 6px;
  margin-left: 4px;
}

.more-images-btn:hover {
  color: #66b1ff;
  background-color: rgba(64, 158, 255, 0.1);
}

/* 时间显示 */
.time-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.datetime-text {
  font-family: 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  color: #6b9bd1;
}

.no-update {
  color: #909399;
  font-size: 12px;
}

.running-time {
  font-size: 10px;
  color: #909399;
  line-height: 1.2;
}

.no-update {
  color: #c0c4cc;
  font-size: 12px;
}


.operation-buttons {
  display: flex;
  gap: 6px;
  align-items: center;
  justify-content: center;
  flex-wrap: nowrap;
}

.operation-buttons .el-button {
  transition: all 0.3s ease;
}

.operation-buttons :deep(.el-dropdown) {
  display: inline-flex;
}

.operation-buttons .el-button:hover:not(.is-disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.operation-buttons .el-button.is-disabled {
  cursor: not-allowed;
  opacity: 0.5;
  background-color: #f5f7fa !important;
  border-color: #e4e7ed !important;
  color: #c0c4cc !important;
}

.operation-buttons .el-button.is-disabled:hover {
  transform: none;
  box-shadow: none;
}

/* 对话框样式 */
.yaml-dialog :deep(.el-dialog),
.scale-dialog :deep(.el-dialog),
.workload-labels-view-dialog :deep(.el-dialog),
.all-images-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.yaml-dialog :deep(.el-dialog__header),
.scale-dialog :deep(.el-dialog__header),
.workload-labels-view-dialog :deep(.el-dialog__header),
.all-images-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.yaml-dialog :deep(.el-dialog__title),
.scale-dialog :deep(.el-dialog__title),
.workload-labels-view-dialog :deep(.el-dialog__title),
.all-images-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.log-controls,
.yaml-controls {
  margin-bottom: 16px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 8px;
}

.log-content,
.yaml-content {
  background: #6b9bd1;
  color: #ecf0f1;
  padding: 16px;
  border-radius: 8px;
  max-height: 400px;
  overflow: auto;
}

.log-content pre,
.yaml-content pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.4;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.4;
}

.dialog-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

/* 集群选择样式 */
.cluster-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.cluster-name {
  font-weight: 500;
  color: #6b9bd1;
}

.cluster-status-tag {
  margin-left: 8px;
}

/* 命名空间选择样式 */
.namespace-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.namespace-name {
  font-weight: 500;
  color: #6b9bd1;
}

.namespace-status-tag {
  margin-left: 8px;
}

/* 通用样式 */
.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.el-input :deep(.el-input__wrapper),
.el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__wrapper):hover,
.el-select :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.el-input :deep(.el-input__wrapper.is-focus),
.el-select :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .operation-buttons {
    gap: 4px;
  }
  
  .operation-buttons .el-button {
    margin: 1px;
  }
  
  .header-actions .el-select {
    min-width: 180px;
  }
}

@media (max-width: 768px) {
  .k8s-workloads-management {
    padding: 10px;
  }
  
  .search-form {
    flex-direction: column;
  }
  
  .search-form .el-form-item {
    margin-right: 0;
    margin-bottom: 12px;
  }
  
  .operation-buttons {
    flex-direction: column;
    gap: 4px;
  }
  
  .workloads-table :deep(.el-table__row:hover) {
    transform: none;
  }
}

/* 标签和镜像对话框样式 */
.labels-view-content,
.images-view-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.labels-view-content {
  gap: 24px;
}

.labels-section {
  margin-top: 0;
}

.labels-section h4,
.section-title {
  margin: 0 0 16px 0;
  color: #6b9bd1;
  font-weight: 600;
  font-size: 16px;
  padding-bottom: 8px;
  border-bottom: 2px solid #e4e7ed;
}

.labels-table-wrapper {
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #dcdfe6;
}

.labels-table {
  width: 100%;
}

.labels-table :deep(.el-table__header-wrapper) {
  background: rgba(103, 126, 234, 0.05);
}

.labels-table :deep(.el-table__header th) {
  background: rgba(103, 126, 234, 0.1) !important;
  color: #6b9bd1;
  font-weight: 600;
}

.labels-table :deep(.el-table__body tr:hover > td) {
  background: rgba(103, 126, 234, 0.08) !important;
}

.labels-table :deep(.el-table td) {
  padding: 10px 12px;
  word-break: break-all;
}

.labels-table :deep(.el-table__body-wrapper) {
  max-height: 300px;
  overflow-y: auto;
}

.no-labels,
.no-images {
  text-align: center;
  padding: 40px 20px;
  color: #909399;
}

/* 镜像对话框特殊样式 */
.all-images-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.image-card {
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  transition: all 0.3s ease;
}

.image-card:hover {
  border-color: #667eea;
  box-shadow: 0 4px 12px rgba(103, 126, 234, 0.15);
}

.image-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.image-name {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.image-icon {
  color: #667eea;
  font-size: 16px;
  flex-shrink: 0;
}

.full-image-name {
  font-family: 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  color: #6b9bd1;
  word-break: break-all;
  line-height: 1.4;
}

.image-actions {
  flex-shrink: 0;
}

.image-details {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

</style>
