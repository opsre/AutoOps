<script setup>
import { ref, reactive, onMounted, onUnmounted, computed, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Monitor,
  View,
  Edit,
  Terminal,
  Document,
  More,
  ArrowLeft,
  Delete,
  Copy,
  ArrowUp,
  ArrowDown
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import CodeEditor from '@/components/CodeEditor.vue'
import LogViewerDialog from './LogViewerDialog.vue'

const route = useRoute()
const router = useRouter()

// Pod状态相关函数
const getPodStatusTag = (status) => {
  const statusMap = {
    'Running': 'success',
    'Pending': 'warning', 
    'Succeeded': 'success',
    'Failed': 'danger',
    'Unknown': 'info'
  }
  return statusMap[status] || 'info'
}

const getPodStatusText = (status) => {
  const textMap = {
    'Running': '运行中',
    'Pending': '等待中',
    'Succeeded': '成功',
    'Failed': '失败',
    'Unknown': '未知'
  }
  return textMap[status] || status
}

// 格式化运行时间
const formatRunningTime = (runningTimeStr) => {
  if (!runningTimeStr) return '-'
  
  const match = runningTimeStr.match(/^(\d+h)?(\d+m)?(\d+(?:\.\d+)?s)?$/)
  if (!match) return runningTimeStr
  
  const [, hours, minutes] = match
  return `${hours || ''}${minutes || ''}`.replace(/^$/, '0m')
}

// 格式化时间戳
const formatTimestamp = (timestamp) => {
  if (!timestamp) return '-'
  
  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    return timestamp
  }
}


const loading = ref(false)
const activeTab = ref('containers')

// 路由参数
const routeParams = reactive({
  clusterId: route.params.clusterId,
  namespace: route.params.namespace,
  podName: route.params.podName
})

// 响应式数据
const podDetail = ref({})
const containers = ref([])
const allRelatedPods = ref([]) // 工作负载下的所有相关Pod
const events = ref([])
const historyVersions = ref([])
const logs = ref('')
const yamlContent = ref('')

// YAML搜索相关状态
const yamlSearchText = ref('')
const yamlSearchResults = ref([])
const yamlSearchCurrentIndex = ref(-1)

// 日志搜索相关状态
const logsSearchText = ref('')
const logsSearchResults = ref([])
const logsSearchCurrentIndex = ref(-1)

// 日志相关状态
const selectedContainerForLogs = ref('')
const logTailLines = ref(300)
const followLogs = ref(false)
const currentLogs = ref('')
const logsLoading = ref(false)
const lastLogRefreshTime = ref(null)
const showPreviousLogs = ref(false)

// YAML标签页相关状态
const yamlTabContent = ref('')
const yamlTabLoading = ref(false)

// 弹框控制
const dialogStates = reactive({
  logsVisible: false,
  yamlVisible: false,
  scaleVisible: false,
  terminalVisible: false,
  monitoringVisible: false,
  rollbackVisible: false
})

// 表单数据
const scaleForm = reactive({
  replicas: 1
})

// 回滚相关数据
const rollbackForm = reactive({
  targetVersion: null,
  versionInfo: {}
})
const rollbackLoading = ref(false)

// 监控相关状态
const selectedPodForMonitoring = ref(null)
const monitoringData = ref({
  cpu: { used: '0', limit: '1000m', percentage: 0 },
  memory: { used: '0Mi', limit: '512Mi', percentage: 0 },
  network: { rx: '0KB/s', tx: '0KB/s' },
  disk: { used: '0GB', limit: '10GB', percentage: 0 }
})
const monitoringLoading = ref(false)

// 当前选中的容器
const selectedContainer = ref(null)

// 标签展开状态
const labelsExpanded = ref(false)

// 计算属性
const podStatus = computed(() => {
  return podDetail.value.status || podDetail.value.phase || 'Unknown'
})

const podLabels = computed(() => {
  return podDetail.value.labels || {}
})

const podCreationTime = computed(() => {
  return podDetail.value.createdAt || podDetail.value.creationTimestamp || ''
})

// 期望pod数量
const expectedPodCount = computed(() => {
  // 如果有相关Pod数据，返回实际运行的Pod数量
  if (allRelatedPods.value && allRelatedPods.value.length > 0) {
    return allRelatedPods.value.length
  }

  // 否则尝试从pod详情中获取
  if (podDetail.value.spec?.ownerReferences) {
    const deployment = podDetail.value.spec.ownerReferences.find(ref =>
      ref.kind === 'ReplicaSet' || ref.kind === 'Deployment'
    )
    if (deployment) {
      return podDetail.value.expectedReplicas || 1
    }
  }
  return podDetail.value.expectedReplicas || 1
})

// 显示的标签数量限制
const maxVisibleLabels = 3

// 可见的标签
const visibleLabels = computed(() => {
  const labelEntries = Object.entries(podLabels.value)
  if (!labelsExpanded.value && labelEntries.length > maxVisibleLabels) {
    return Object.fromEntries(labelEntries.slice(0, maxVisibleLabels))
  }
  return podLabels.value
})

// 是否有更多标签
const hasMoreLabels = computed(() => {
  return Object.keys(podLabels.value).length > maxVisibleLabels
})

// 隐藏的标签数量
const hiddenLabelsCount = computed(() => {
  const totalLabels = Object.keys(podLabels.value).length
  return Math.max(0, totalLabels - maxVisibleLabels)
})

// 工作负载名称
const workloadName = computed(() => {
  // 先尝试从ownerReferences获取
  const ownerRefs = podDetail.value.metadata?.ownerReferences ||
                   podDetail.value.spec?.ownerReferences ||
                   podDetail.value.ownerReferences || []

  // 查找Deployment
  const deployment = ownerRefs.find(ref => ref.kind === 'Deployment')
  if (deployment) {
    return deployment.name
  }

  // 查找ReplicaSet，然后推断Deployment名称
  const replicaSet = ownerRefs.find(ref => ref.kind === 'ReplicaSet')
  if (replicaSet) {
    // ReplicaSet命名规则: deployment-name-pod-template-hash
    const parts = replicaSet.name.split('-')
    if (parts.length >= 2) {
      return parts.slice(0, -1).join('-')
    }
  }

  // 最后从Pod名称推断Deployment名称
  const podNameParts = routeParams.podName.split('-')
  if (podNameParts.length >= 3) {
    // Pod命名规律: deployment-name-replicaset-hash-pod-hash
    return podNameParts.slice(0, -2).join('-')
  }

  // 如果都无法推断，返回Pod名称
  return routeParams.podName
})

// 当前Pod信息（用于表格显示）
const currentPodInfo = computed(() => {
  if (!podDetail.value || Object.keys(podDetail.value).length === 0) {
    return {
      name: routeParams.podName,
      status: 'Loading...',
      restartCount: '-',
      nodeName: '-',
      podIP: '-',
      hostIP: '-',
      age: '-',
      runningTime: '-',
      containers: [],
      resources: { requests: { cpu: '', memory: '' }, limits: { cpu: '', memory: '' } },
      labels: {},
      conditions: []
    }
  }

  return {
    name: routeParams.podName,
    status: podDetail.value.status || podDetail.value.phase || 'Unknown',
    restartCount: podDetail.value.restartCount || 0,
    nodeName: podDetail.value.spec?.nodeName || podDetail.value.nodeName || 'Unknown',
    podIP: podDetail.value.podIP || podDetail.value.status?.podIP ||
           (podDetail.value.status === 'Pending' ? '等待分配' :
            podDetail.value.status === 'Terminating' ? '已释放' : 'Unknown'),
    hostIP: podDetail.value.hostIP || podDetail.value.status?.hostIP || 'Unknown',
    age: formatAge(podDetail.value.metadata?.creationTimestamp || podDetail.value.createdAt),
    runningTime: podDetail.value.runningTime || '',
    containers: containers.value || podDetail.value.spec?.containers || [],
    resources: podDetail.value.resources || { requests: { cpu: '', memory: '' }, limits: { cpu: '', memory: '' } },
    labels: podDetail.value.labels || {},
    conditions: podDetail.value.conditions || []
  }
})


// 切换标签展开状态
const toggleLabelsExpanded = () => {
  labelsExpanded.value = !labelsExpanded.value
}

// 监听标签页切换
const handleTabChange = (tabName) => {
  if (tabName === 'yaml' && !yamlTabContent.value) {
    loadYamlContent()
  }
  
  // 切换标签页时清空搜索状态
  if (tabName !== 'logs') {
    clearLogsSearch()
  }
  if (tabName !== 'yaml') {
    clearYamlSearch()
  }
}

// 加载YAML内容
const loadYamlContent = async () => {
  if (yamlTabLoading.value) return
  
  try {
    yamlTabLoading.value = true
    console.log('🔍 开始加载YAML内容...')
    
    const response = await k8sApi.getPodYaml(routeParams.clusterId, routeParams.namespace, routeParams.podName)
    const responseData = response.data || response
    
    if (responseData.code === 200) {
      const rawData = responseData.data || ''
      if (typeof rawData === 'string') {
        yamlTabContent.value = rawData
      } else if (typeof rawData === 'object') {
        yamlTabContent.value = JSON.stringify(rawData, null, 2)
      } else {
        yamlTabContent.value = String(rawData)
      }
      console.log('✅ YAML内容加载成功，长度:', yamlTabContent.value.length)
      
      if (!yamlTabContent.value.trim()) {
        console.log('⚠️ YAML内容为空，生成默认YAML')
        generateYamlFromDetail()
        return
      }
    } else {
      console.log('❌ YAML API返回错误，生成默认YAML:', responseData.message)
      generateYamlFromDetail()
    }
  } catch (error) {
    console.error('❌ 获取YAML内容失败，生成默认YAML:', error)
    generateYamlFromDetail()
  } finally {
    yamlTabLoading.value = false
  }
}

// 定时刷新器
let refreshTimer = null

// 页面初始化
onMounted(() => {
  if (routeParams.clusterId && routeParams.namespace && routeParams.podName) {
    handleQuery()
    // 移除自动刷新，仅支持手动刷新
    // startAutoRefresh()
    addVisibilityListener()
  }
})

// 页面卸载时清理
onUnmounted(() => {
  stopAutoRefresh()
  removeVisibilityListener()
})

// 监听路由参数变化
watch(() => [routeParams.clusterId, routeParams.namespace, routeParams.podName], () => {
  console.log('路由参数变化，重新加载数据')
  if (routeParams.clusterId && routeParams.namespace && routeParams.podName) {
    handleQuery()
  }
}, { immediate: false })

// 监听标签页切换 - 当切换到日志标签页时自动加载日志
watch(() => activeTab.value, (newTab) => {
  console.log('🔄 标签页切换到:', newTab)
  if (newTab === 'logs') {
    console.log('🔄 切换到日志标签页，检查是否需要加载日志')
    // 延迟一点确保DOM已更新
    nextTick(() => {
      if (selectedContainerForLogs.value && !currentLogs.value) {
        console.log('🔄 自动加载日志，因为当前无日志内容')
        handleRefreshLogs()
      } else if (!selectedContainerForLogs.value && allRelatedPods.value.length > 0) {
        console.log('🔄 自动选择Pod并加载日志')
        selectedContainerForLogs.value = allRelatedPods.value[0].name
        setTimeout(() => {
          handleRefreshLogs()
        }, 100)
      }
    })
  }
})

// 启动自动刷新 (每30秒刷新一次)
const startAutoRefresh = () => {
  refreshTimer = setInterval(() => {
    console.log('自动刷新Pod数据...')
    handleQuery()
  }, 30000) // 30秒
}

// 停止自动刷新
const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// 页面可见性变化处理
const handleVisibilityChange = () => {
  if (!document.hidden) {
    console.log('页面重新获得焦点，刷新数据')
    handleQuery()
  }
}

// 添加页面可见性监听器
const addVisibilityListener = () => {
  document.addEventListener('visibilitychange', handleVisibilityChange)
}

// 移除页面可见性监听器  
const removeVisibilityListener = () => {
  document.removeEventListener('visibilitychange', handleVisibilityChange)
}

// 获取Pod详情
const handleQuery = async (forceRefresh = false) => {
  try {
    loading.value = true
    console.log(forceRefresh ? '🔄 强制刷新Pod详情数据...' : '📥 获取Pod详情数据...')
    console.log('🔍 [Pod详情页面] 请求参数:', {
      clusterId: routeParams.clusterId,
      namespace: routeParams.namespace,
      podName: routeParams.podName
    })

    // 从Pod名称推断Deployment名称
    const podNameParts = routeParams.podName.split('-')
    let deploymentName = null
    if (podNameParts.length >= 3) {
      // Pod命名规律: deployment-name-replicaset-hash-pod-hash
      deploymentName = podNameParts.slice(0, -2).join('-')
      console.log('🔍 从Pod名称推断Deployment:', deploymentName)
    }

    let response, responseData
    if (deploymentName) {
      try {
        console.log('🔍 使用工作负载API获取Pod详情...')
        response = await k8sApi.getWorkloadPods(
          routeParams.clusterId,
          routeParams.namespace,
          'deployment',
          deploymentName
        )
        responseData = response.data || response
        console.log('工作负载Pod列表API响应:', responseData)

        if (responseData.code === 200 && responseData.data) {
          // 从Pod列表中找到当前Pod
          const currentPod = responseData.data.find(pod => pod.name === routeParams.podName)
          if (currentPod) {
            podDetail.value = currentPod
            containers.value = currentPod.containers || []
            console.log('✅ 从工作负载Pod列表中找到目标Pod:', currentPod.name)

            // 同时设置所有相关Pod
            allRelatedPods.value = responseData.data.map(pod => ({
              name: pod.name,
              status: pod.status || pod.phase || 'Unknown',
              restartCount: pod.restartCount || 0,
              nodeName: pod.nodeName || 'Unknown',
              podIP: pod.podIP || 'Unknown',
              hostIP: pod.hostIP || 'Unknown',
              age: pod.age || formatAge(pod.createdAt),
              runningTime: pod.runningTime || '',
              containers: pod.containers || [],
              resources: pod.resources || { requests: { cpu: '', memory: '' }, limits: { cpu: '', memory: '' } },
              labels: pod.labels || {},
              conditions: pod.conditions || [],
              rawData: pod
            }))
            console.log(`✅ 设置了 ${allRelatedPods.value.length} 个相关Pod`)

            // 自动选择当前Pod用于日志显示
            if (!selectedContainerForLogs.value && allRelatedPods.value.length > 0) {
              selectedContainerForLogs.value = routeParams.podName
              console.log('🔄 自动选择当前Pod用于日志:', selectedContainerForLogs.value)
            }
          } else {
            throw new Error('在工作负载Pod列表中未找到目标Pod')
          }
        } else {
          throw new Error(responseData.message || '工作负载API调用失败')
        }
      } catch (workloadError) {
        console.log('⚠️ 工作负载API失败，尝试其他方式:', workloadError.message)
        // 如果工作负载API失败，构造基本的Pod信息
        podDetail.value = {
          name: routeParams.podName,
          status: 'Unknown',
          metadata: { name: routeParams.podName, namespace: routeParams.namespace }
        }
        containers.value = []
      }
    } else {
      console.log('⚠️ 无法从Pod名称推断Deployment，构造基本信息')
      podDetail.value = {
        name: routeParams.podName,
        status: 'Unknown',
        metadata: { name: routeParams.podName, namespace: routeParams.namespace }
      }
      containers.value = []
    }

    if (podDetail.value && podDetail.value.name) {
      
      // 工作负载API不提供事件数据，直接获取事件
      console.log('🔄 获取Pod事件数据...')
      await getEvents()

      // 获取相关的所有Pod（如果还没有的话）
      if (allRelatedPods.value.length === 0) {
        console.log('🔄 [Pod详情页面] 开始获取相关Pod列表')
        await getRelatedPods()
      } else {
        console.log('✅ 已有相关Pod数据，跳过重复获取')
      }

      // 每次都尝试获取最新事件数据
      console.log('🔄 强制刷新事件数据...')
      await getEvents()
      
      // 获取历史版本数据
      await getHistoryVersions()
    } else {
      ElMessage.error('获取Pod详情失败')
    }
  } catch (error) {
    console.error('获取Pod详情失败:', error)
    ElMessage.error('获取Pod详情失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 获取工作负载下的所有相关Pod
const getRelatedPods = async () => {
  try {
    console.log('🔍 获取相关Pod列表...')
    console.log('📊 当前Pod详情数据:', {
      podName: routeParams.podName,
      metadata: podDetail.value.metadata,
      labels: podDetail.value.metadata?.labels || podDetail.value.labels,
      ownerReferences: podDetail.value.metadata?.ownerReferences || podDetail.value.ownerReferences,
      // 调试: 显示Pod详情的所有顶级字段
      allFields: Object.keys(podDetail.value)
    })

    console.log('📊 Pod详情完整数据结构:', podDetail.value)
    
    // 从当前Pod的ownerReferences获取工作负载信息
    const ownerRefs = podDetail.value.ownerReferences || podDetail.value.metadata?.ownerReferences || []
    console.log('🔍 [Pod详情页面] Pod的ownerReferences:', ownerRefs)
    let labelSelector = ''
    
    if (ownerRefs.length === 0) {
      console.log('⚠️ 当前Pod没有ownerReferences，尝试通过标签获取相关Pod')
      // 尝试通过app标签获取相关Pod
      const labels = podDetail.value.labels || podDetail.value.metadata?.labels || {}
      console.log('🔍 [Pod详情页面] 可用的标签:', labels)

      const appLabel = labels['app']
      const nameLabel = labels['app.kubernetes.io/name']
      
      if (appLabel) {
        labelSelector = `app=${appLabel}`
      } else if (nameLabel) {
        labelSelector = `app.kubernetes.io/name=${nameLabel}`
      }
      
      if (!labelSelector) {
        console.log('⚠️ 无法构建标签选择器，使用当前Pod')
        // 构造当前Pod数据结构
        allRelatedPods.value = [{
          name: routeParams.podName,
          status: podDetail.value.status || podDetail.value.phase || 'Unknown',
          restartCount: podDetail.value.restartCount || 0,
          nodeName: podDetail.value.spec?.nodeName || 'Unknown',
          podIP: podDetail.value.podIP || podDetail.value.status?.podIP || (podDetail.value.status === 'Pending' ? '等待分配' : podDetail.value.status === 'Terminating' ? '已释放' : 'Unknown'),
          hostIP: podDetail.value.hostIP || podDetail.value.status?.hostIP || 'Unknown',
          age: formatAge(podDetail.value.metadata?.creationTimestamp),
          runningTime: podDetail.value.runningTime || '',
          containers: containers.value || podDetail.value.spec?.containers || [],
          resources: { requests: { cpu: '', memory: '' }, limits: { cpu: '', memory: '' } },
          labels: podDetail.value.metadata?.labels || {},
          conditions: podDetail.value.status?.conditions || [],
          rawData: podDetail.value
        }]
        return
      }
    } else {
      // 查找ReplicaSet或直接的Deployment引用
      const replicaSetRef = ownerRefs.find(ref => ref.kind === 'ReplicaSet')
      const deploymentRef = ownerRefs.find(ref => ref.kind === 'Deployment')

      console.log('🔍 [Pod详情页面] 找到的owner引用:', {
        replicaSet: replicaSetRef?.name,
        deployment: deploymentRef?.name
      })
      
      if (replicaSetRef) {
        console.log('📋 通过ReplicaSet获取相关Pod:', replicaSetRef.name)
      } else if (deploymentRef) {
        console.log('📋 通过Deployment获取相关Pod:', deploymentRef.name)
      }
      
      // 优先使用app标签，这样可以获取所有相关的Pod（不仅仅是相同版本）
      const labels = podDetail.value.labels || podDetail.value.metadata?.labels || {}
      console.log('🔍 [Pod详情页面] 工作负载Pod的可用标签:', labels)

      const appLabel = labels['app']
      const nameLabel = labels['app.kubernetes.io/name']
      const instanceLabel = labels['app.kubernetes.io/instance']

      if (appLabel) {
        labelSelector = `app=${appLabel}`
        console.log('🔍 [Pod详情页面] 使用app标签选择器:', labelSelector)
      } else if (nameLabel) {
        labelSelector = `app.kubernetes.io/name=${nameLabel}`
        console.log('🔍 [Pod详情页面] 使用name标签选择器:', labelSelector)
      } else if (instanceLabel) {
        labelSelector = `app.kubernetes.io/instance=${instanceLabel}`
        console.log('🔍 [Pod详情页面] 使用instance标签选择器:', labelSelector)
      } else {
        // 最后选择：通过pod-template-hash（只会找到相同版本的Pod）
        const podTemplateHash = labels['pod-template-hash']
        if (podTemplateHash) {
          labelSelector = `pod-template-hash=${podTemplateHash}`
          console.log('🔍 [Pod详情页面] 使用pod-template-hash标签选择器:', labelSelector)
        }
      }
    }
    
    // 统一的标签选择器处理
    
    if (!labelSelector) {
      console.log('⚠️ 无法构建标签选择器，使用当前Pod')
      // 构造当前Pod数据结构
      allRelatedPods.value = [{
        name: routeParams.podName,
        status: podDetail.value.status || podDetail.value.phase || 'Unknown',
        restartCount: podDetail.value.restartCount || 0,
        nodeName: podDetail.value.spec?.nodeName || 'Unknown',
        podIP: podDetail.value.podIP || podDetail.value.status?.podIP || (podDetail.value.status === 'Pending' ? '等待分配' : podDetail.value.status === 'Terminating' ? '已释放' : 'Unknown'),
        hostIP: podDetail.value.hostIP || podDetail.value.status?.hostIP || 'Unknown',
        age: formatAge(podDetail.value.metadata?.creationTimestamp),
        runningTime: podDetail.value.runningTime || '',
        containers: containers.value || podDetail.value.spec?.containers || [],
        resources: {
          requests: { cpu: '', memory: '' },
          limits: { cpu: '', memory: '' }
        },
        labels: podDetail.value.metadata?.labels || {},
        conditions: podDetail.value.status?.conditions || [],
        rawData: podDetail.value
      }]
      return
    }
    
    // 尝试使用新的工作负载API获取Pod列表
    let response
    let useNewApi = false

    // 从Pod名称智能推断Deployment名称
    let inferredDeploymentName = null
    const podNameParts = routeParams.podName.split('-')
    if (podNameParts.length >= 3) {
      // Pod命名规律: deployment-name-replicaset-hash-pod-hash
      // 例如: zf-nginx-test-5d9f8cbcdb-8b94f
      // Deployment名称通常是去掉最后两个部分
      inferredDeploymentName = podNameParts.slice(0, -2).join('-')
      console.log('🔍 [Pod详情页面] 从Pod名称推断Deployment:', inferredDeploymentName)
    }

    // 检查是否可以推断出工作负载信息
    if (deploymentRef) {
      console.log('🔍 [Pod详情页面] 使用新API - Deployment:', deploymentRef.name)
      try {
        response = await k8sApi.getWorkloadPods(
          routeParams.clusterId,
          routeParams.namespace,
          'deployment',
          deploymentRef.name
        )
        useNewApi = true
      } catch (error) {
        console.log('⚠️ 新API调用失败，回退到标签选择器:', error.message)
      }
    } else if (replicaSetRef) {
      // 尝试从ReplicaSet名称推断Deployment名称
      const replicaSetName = replicaSetRef.name
      // ReplicaSet命名规则通常是: deployment-name-pod-template-hash
      const parts = replicaSetName.split('-')
      if (parts.length >= 2) {
        const deploymentName = parts.slice(0, -1).join('-')
        console.log('🔍 [Pod详情页面] 从ReplicaSet推断Deployment名称:', deploymentName)
        try {
          response = await k8sApi.getWorkloadPods(
            routeParams.clusterId,
            routeParams.namespace,
            'deployment',
            deploymentName
          )
          useNewApi = true
        } catch (error) {
          console.log('⚠️ 推断的Deployment名称无效，回退到标签选择器:', error.message)
        }
      }
    }

    // 如果新API没有成功，回退到旧的标签选择器方式
    if (!useNewApi) {
      console.log('🔍 [Pod详情页面] 使用旧API - 标签选择器:', labelSelector)
      response = await k8sApi.getPodList(routeParams.clusterId, routeParams.namespace, {
        labelSelector: labelSelector
      })
    }
    
    console.log('📡 Pod列表API响应:', response)
    const responseData = response.data || response
    console.log('📊 Pod列表响应数据:', responseData)

    if (responseData.code === 200 || responseData.success) {
      // 处理不同API的响应格式
      let pods = []
      if (useNewApi) {
        // 新API直接返回Pod数组
        pods = responseData.data || []
        console.log('📋 [新API] 获取到Pod列表:', pods.length, '个')
      } else {
        // 旧API返回 data.pods 或 data
        pods = responseData.data?.pods || responseData.data || []
        console.log('📋 [旧API] 获取到Pod列表:', pods.length, '个')
      }
      console.log('📋 原始Pod数据:', pods)
      console.log('📊 Pod数量:', pods.length)
      allRelatedPods.value = pods.map(pod => ({
        name: pod.name,
        status: pod.status || pod.phase || 'Unknown',
        restartCount: useNewApi ? (pod.restarts || pod.restartCount || 0) : (pod.restartCount || 0),
        nodeName: pod.nodeName || 'Unknown',
        podIP: pod.podIP || pod.status?.podIP || (pod.status === 'Pending' ? '等待分配' : pod.status === 'Terminating' ? '已释放' : 'Unknown'),
        hostIP: pod.hostIP || pod.status?.hostIP || 'Unknown',
        age: useNewApi ? (pod.age || formatAge(pod.createdAt)) : formatAge(pod.createdAt),
        runningTime: pod.runningTime || '',
        containers: pod.containers || [],
        resources: pod.resources || {
          requests: { cpu: '', memory: '' },
          limits: { cpu: '', memory: '' }
        },
        labels: pod.labels || {},
        conditions: pod.conditions || [],
        rawData: pod
      }))
      
      console.log(`✅ 成功获取 ${allRelatedPods.value.length} 个相关Pod`)
      
      // 自动选择第一个Pod用于日志显示
      if (allRelatedPods.value.length > 0 && !selectedContainerForLogs.value) {
        selectedContainerForLogs.value = allRelatedPods.value[0].name
        console.log('🔄 自动选择第一个Pod用于日志:', selectedContainerForLogs.value)
        // 自动加载第一个Pod的日志
        setTimeout(() => {
          handleRefreshLogs()
        }, 100)
      }
    } else {
      console.log('❌ 获取Pod列表失败:', responseData.message)
      // 构造当前Pod数据结构
      allRelatedPods.value = [{
        name: routeParams.podName,
        status: podDetail.value.status || podDetail.value.phase || 'Unknown',
        restartCount: podDetail.value.restartCount || 0,
        nodeName: podDetail.value.spec?.nodeName || 'Unknown',
        podIP: podDetail.value.podIP || podDetail.value.status?.podIP || (podDetail.value.status === 'Pending' ? '等待分配' : podDetail.value.status === 'Terminating' ? '已释放' : 'Unknown'),
        hostIP: podDetail.value.hostIP || podDetail.value.status?.hostIP || 'Unknown',
        age: formatAge(podDetail.value.metadata?.creationTimestamp),
        runningTime: podDetail.value.runningTime || '',
        containers: containers.value || podDetail.value.spec?.containers || [],
        resources: { requests: { cpu: '', memory: '' }, limits: { cpu: '', memory: '' } },
        labels: podDetail.value.metadata?.labels || {},
        conditions: podDetail.value.status?.conditions || [],
        rawData: podDetail.value
      }]
    }
  } catch (error) {
    console.error('❌ 获取相关Pod失败:', error)
    // 构造当前Pod数据结构
    allRelatedPods.value = [{
      name: routeParams.podName,
      status: podDetail.value.status || podDetail.value.phase || 'Unknown',
      restartCount: podDetail.value.restartCount || 0,
      nodeName: podDetail.value.spec?.nodeName || 'Unknown',
      podIP: podDetail.value.podIP || podDetail.value.status?.podIP || (podDetail.value.status === 'Pending' ? '等待分配' : podDetail.value.status === 'Terminating' ? '已释放' : 'Unknown'),
      hostIP: podDetail.value.hostIP || podDetail.value.status?.hostIP || 'Unknown',
      age: formatAge(podDetail.value.metadata?.creationTimestamp),
      runningTime: podDetail.value.runningTime || '',
      containers: podDetail.value.spec?.containers || [],
      resources: { requests: { cpu: '', memory: '' }, limits: { cpu: '', memory: '' } },
      labels: podDetail.value.metadata?.labels || {},
      conditions: podDetail.value.status?.conditions || [],
      rawData: podDetail.value
    }]
  }
  
  // 自动选择第一个Pod用于日志显示
  if (allRelatedPods.value.length > 0 && !selectedContainerForLogs.value) {
    selectedContainerForLogs.value = allRelatedPods.value[0].name
    console.log('🔄 自动选择第一个Pod用于日志:', selectedContainerForLogs.value)
    // 自动加载第一个Pod的日志
    setTimeout(() => {
      handleRefreshLogs()
    }, 100)
  }
}

// 格式化时间
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

// 获取相关的所有事件（Pod、ReplicaSet、Deployment）
const getEvents = async () => {
  try {
    console.log('开始获取事件数据...')
    console.log('请求参数:', {
      clusterId: routeParams.clusterId,
      namespace: routeParams.namespace,
      podName: routeParams.podName
    })
    
    // 获取最近1小时的时间戳
    const oneHourAgo = new Date(Date.now() - 60 * 60 * 1000)
    console.log('过滤时间: 最近1小时内的事件 (从', oneHourAgo.toISOString(), '开始)')
    
    // 根据后端API文档，构建事件查询参数
    console.log('构建事件查询参数...')
    
    // 从Pod名称中推断相关资源名称
    const podName = routeParams.podName
    let deploymentName = ''
    let replicaSetName = ''
    
    // 解析Pod名称格式：nginx-deployment-79df99db77-85qr5
    if (podName.includes('-')) {
      const parts = podName.split('-')
      if (parts.length >= 3) {
        // 去掉最后两个部分（ReplicaSet hash 和 Pod hash）
        deploymentName = parts.slice(0, -2).join('-')
        // 去掉最后一个部分（Pod hash）
        replicaSetName = parts.slice(0, -1).join('-')
      }
    }
    
    console.log('推断的资源名称:', {
      podName,
      replicaSetName,
      deploymentName
    })
    
    // 构建查询参数，获取Pod相关的所有事件
    const eventParams = {
      kind: 'Pod', // 首先获取Pod类型的事件
      name: podName, // 过滤特定Pod名称
      limit: 100 // 限制返回数量
    }
    
    console.log('事件API查询参数:', eventParams)
    const response = await k8sApi.getClusterEvents(routeParams.clusterId, eventParams)
    const responseData = response.data || response
    
    console.log('事件API完整响应结构:', {
      success: !!responseData,
      code: responseData.code,
      dataType: typeof responseData.data,
      hasEvents: !!responseData.data?.events,
      eventsLength: responseData.data?.events?.length || 0,
      total: responseData.data?.total,
      namespace: responseData.data?.namespace
    })
    
    if (responseData.code === 200) {
      // 根据新的API格式，事件数据在 data.events 中
      let podEvents = responseData.data?.events || []
      
      console.log('Pod事件数据:', {
        count: podEvents.length,
        total: responseData.data?.total,
        namespace: responseData.data?.namespace
      })
      
      // 如果没有Pod事件，尝试获取相关的Deployment和ReplicaSet事件
      if (podEvents.length === 0 && (deploymentName || replicaSetName)) {
        console.log('未找到Pod事件，尝试获取相关资源事件...')
        
        // 获取Deployment事件
        if (deploymentName) {
          try {
            const deploymentParams = { kind: 'Deployment', name: deploymentName, limit: 50 }
            const deploymentResponse = await k8sApi.getClusterEvents(routeParams.clusterId, deploymentParams)
            if (deploymentResponse.data?.code === 200) {
              const deploymentEvents = deploymentResponse.data?.data?.events || []
              podEvents = podEvents.concat(deploymentEvents)
              console.log(`添加了 ${deploymentEvents.length} 个Deployment事件`)
            }
          } catch (error) {
            console.log('获取Deployment事件失败:', error.message)
          }
        }
        
        // 获取ReplicaSet事件
        if (replicaSetName) {
          try {
            const replicaSetParams = { kind: 'ReplicaSet', name: replicaSetName, limit: 50 }
            const replicaSetResponse = await k8sApi.getClusterEvents(routeParams.clusterId, replicaSetParams)
            if (replicaSetResponse.data?.code === 200) {
              const replicaSetEvents = replicaSetResponse.data?.data?.events || []
              podEvents = podEvents.concat(replicaSetEvents)
              console.log(`添加了 ${replicaSetEvents.length} 个ReplicaSet事件`)
            }
          } catch (error) {
            console.log('获取ReplicaSet事件失败:', error.message)
          }
        }
      }
      
      // 检查事件数据结构
      console.log('处理后的事件数据结构检查:')
      if (podEvents.length > 0) {
        const sampleEvent = podEvents[0]
        console.log('第一个事件的字段:', Object.keys(sampleEvent))
        console.log('事件示例:', {
          type: sampleEvent.type,
          reason: sampleEvent.reason,
          message: sampleEvent.message,
          source: sampleEvent.source,
          count: sampleEvent.count,
          firstTime: sampleEvent.firstTime,
          lastTime: sampleEvent.lastTime
        })
      } else {
        console.log('没有获取到任何相关事件数据')
      }
      
      // 从Pod名称中推断相关资源名称
      const podNameParts = routeParams.podName.split('-')
      console.log('Pod名称分析:', {
        fullName: routeParams.podName,
        parts: podNameParts,
        partsCount: podNameParts.length
      })
      
      // 生成多种可能的名称匹配模式
      const matchPatterns = []
      
      // 精确匹配Pod名称
      matchPatterns.push({
        name: routeParams.podName,
        kind: 'Pod',
        pattern: 'exact_pod'
      })
      
      // 如果Pod名称有足够的段数，推断ReplicaSet名称 (去掉最后一段随机字符)
      if (podNameParts.length >= 3) {
        const replicaSetName = podNameParts.slice(0, -1).join('-')
        matchPatterns.push({
          name: replicaSetName,
          kind: 'ReplicaSet', 
          pattern: 'replicaset_inferred'
        })
        
        // 推断Deployment名称 (去掉最后两段：hash和随机字符)
        if (podNameParts.length >= 4) {
          const deploymentName = podNameParts.slice(0, -2).join('-')
          matchPatterns.push({
            name: deploymentName,
            kind: 'Deployment',
            pattern: 'deployment_inferred'
          })
        }
      }
      
      console.log('匹配模式:', matchPatterns)
      
      // 过滤事件（这部分代码已经不需要，因为新API已经在后端过滤）
      const filteredEvents = podEvents.filter(event => {
        // 检查基本数据结构
        if (!event.involvedObject || !event.involvedObject.name || !event.involvedObject.kind) {
          return false
        }
        
        const eventObject = event.involvedObject
        const eventTime = new Date(event.lastTimestamp || event.lastTime || event.firstTimestamp || 0)
        
        // 时间过滤：只显示最近1小时内的事件
        if (eventTime < oneHourAgo) {
          return false
        }
        
        // 命名空间过滤
        if (eventObject.namespace !== routeParams.namespace) {
          return false  
        }
        
        // 资源名称匹配
        const isMatched = matchPatterns.some(pattern => {
          // 精确匹配
          if (pattern.pattern === 'exact_pod') {
            return eventObject.name === pattern.name && eventObject.kind === 'Pod'
          }
          
          // ReplicaSet匹配
          if (pattern.pattern === 'replicaset_inferred') {
            return eventObject.name === pattern.name && eventObject.kind === 'ReplicaSet'
          }
          
          // Deployment匹配
          if (pattern.pattern === 'deployment_inferred') {
            return eventObject.name === pattern.name && eventObject.kind === 'Deployment'
          }
          
          // 前缀匹配 - 相关的资源
          if (eventObject.name.startsWith(pattern.name)) {
            return ['Pod', 'ReplicaSet', 'Deployment', 'Service'].includes(eventObject.kind)
          }
          
          return false
        })
        
        if (isMatched) {
          console.log('✅ 匹配事件:', {
            name: eventObject.name,
            kind: eventObject.kind,
            namespace: eventObject.namespace,
            type: event.type,
            reason: event.reason,
            time: eventTime.toLocaleString('zh-CN')
          })
        }
        
        return isMatched
      })
      
      console.log(`过滤结果: ${filteredEvents.length}/${podEvents.length} 个事件匹配`)
      
      // 按时间排序（最新的在前）
      filteredEvents.sort((a, b) => {
        const timeA = new Date(a.lastTimestamp || a.lastTime || a.firstTimestamp || 0)
        const timeB = new Date(b.lastTimestamp || b.lastTime || b.firstTimestamp || 0)
        return timeB - timeA
      })
      
      // 直接使用新API返回的已过滤事件
      events.value = podEvents
      
      // 结果处理
      if (events.value.length === 0) {
        console.log('❌ 未找到相关事件数据')
        // 移除警告提示，静默处理无事件的情况
      } else {
        console.log(`✅ 成功加载 ${events.value.length} 个相关事件`)
        // 静默加载，不显示成功提示
      }
    } else {
      console.warn('❌ 获取集群事件API失败:', responseData)
      ElMessage.error(responseData.message || '获取事件失败')
    }
  } catch (error) {
    console.error('❌ 获取事件异常:', error)
    
    // 404错误表示没有事件数据，这是正常情况
    if (error.response && error.response.status === 404) {
      console.log('📋 集群暂无事件数据，这是正常情况')
      events.value = []
      return
    }
    
    // 其他错误才显示错误消息
    ElMessage.error('获取事件失败: ' + (error.message || '网络连接异常'))
  }
}

// 直接显示Pod详情中的事件数据
const showPodDetailEvents = () => {
  console.log('尝试直接显示Pod详情中的事件数据')
  console.log('podDetail.value:', podDetail.value)
  
  // 检查Pod详情中是否有events字段
  if (podDetail.value.events && Array.isArray(podDetail.value.events)) {
    events.value = podDetail.value.events
    console.log('直接使用Pod详情中的事件数据:', events.value.length, '个事件')
    ElMessage.success(`已加载 ${events.value.length} 个事件`)
  } else {
    // 尝试其他可能的事件数据位置
    const possibleEventFields = ['events', 'eventList', 'podEvents', 'conditions']
    let foundEvents = false
    
    for (const field of possibleEventFields) {
      if (podDetail.value[field] && Array.isArray(podDetail.value[field])) {
        events.value = podDetail.value[field]
        console.log(`在 ${field} 字段找到事件数据:`, events.value.length, '个')
        ElMessage.success(`从 ${field} 字段加载了 ${events.value.length} 个事件`)
        foundEvents = true
        break
      }
    }
    
    if (!foundEvents) {
      console.log('Pod详情中未找到任何事件数据')
      console.log('Pod详情的所有字段:', Object.keys(podDetail.value))
      ElMessage.warning('Pod详情中未找到事件数据')
    }
  }
}

// 获取历史版本（通过workload获取）
const getHistoryVersions = async () => {
  try {
    console.log('📋 开始获取历史版本数据...')
    console.log('Pod详情数据:', podDetail.value)
    
    // 尝试从Pod的ownerReferences获取工作负载信息
    const ownerRefs = podDetail.value.metadata?.ownerReferences || 
                     podDetail.value.spec?.ownerReferences || 
                     podDetail.value.ownerReferences || []
    
    console.log('Pod ownerReferences:', ownerRefs)
    
    // 查找ReplicaSet或Deployment
    const replicaSet = ownerRefs.find(ref => ref.kind === 'ReplicaSet')
    const deployment = ownerRefs.find(ref => ref.kind === 'Deployment')
    
    console.log('识别的控制器:', {
      hasReplicaSet: !!replicaSet,
      hasDeployment: !!deployment,
      replicaSetName: replicaSet?.name,
      deploymentName: deployment?.name
    })
    
    if (deployment) {
      // 如果是Deployment管理的Pod，获取Deployment的历史版本
      console.log('🔍 [历史版本] 获取Deployment历史版本:', deployment.name)
      console.log('🔍 [历史版本] 调用getWorkloadDetail API')
      try {
        const response = await k8sApi.getWorkloadDetail(routeParams.clusterId, routeParams.namespace, 'deployment', deployment.name)
        const responseData = response.data || response
        
        console.log('Deployment详情API响应:', {
          success: responseData.code === 200,
          hasRevisionHistory: !!responseData.data?.revisionHistory,
          historyLength: responseData.data?.revisionHistory?.length || 0
        })
        
        if (responseData.code === 200 && responseData.data?.revisionHistory) {
          historyVersions.value = responseData.data.revisionHistory.map((version, index) => ({
            revision: version.revision || (index + 1),
            creationTime: version.creationTime || version.createdAt,
            replicas: version.replicas || responseData.data.replicas || 1,
            description: `Deployment: ${deployment.name} - Rev ${version.revision || (index + 1)}`,
            changeDescription: version.changeDescription || version.description,
            image: version.image,
            status: version.status || 'Active'
          }))
          console.log('✅ 获取到Deployment历史版本:', historyVersions.value.length, '个版本')
        } else {
          throw new Error('Deployment API返回失败或无历史版本数据')
        }
      } catch (deploymentError) {
        console.log('⚠️  获取Deployment详情失败，生成基础版本信息')
        historyVersions.value = [{
          revision: 1,
          creationTime: podDetail.value.createdAt || podDetail.value.creationTimestamp,
          replicas: 1,
          description: `Deployment: ${deployment.name} (详情获取失败)`,
          status: 'Current'
        }]
      }
    } else if (replicaSet) {
      // 如果是ReplicaSet管理的Pod
      console.log('📝 生成ReplicaSet版本信息:', replicaSet.name)
      
      // 尝试获取ReplicaSet的详细信息
      console.log('🔍 [历史版本] 调用getWorkloadDetail API - ReplicaSet')
      try {
        const response = await k8sApi.getWorkloadDetail(routeParams.clusterId, routeParams.namespace, 'replicaset', replicaSet.name)
        const responseData = response.data || response
        
        if (responseData.code === 200 && responseData.data) {
          historyVersions.value = [{
            revision: 1,
            creationTime: responseData.data.createdAt || responseData.data.creationTimestamp,
            replicas: responseData.data.replicas || 1,
            description: `ReplicaSet: ${replicaSet.name}`,
            image: responseData.data.containers?.[0]?.image,
            status: 'Active'
          }]
        } else {
          throw new Error('ReplicaSet API返回失败')
        }
      } catch (replicaSetError) {
        console.log('⚠️  获取ReplicaSet详情失败，使用基础信息')
        historyVersions.value = [{
          revision: 1,
          creationTime: podDetail.value.createdAt || podDetail.value.creationTimestamp,
          replicas: 1,
          description: `ReplicaSet: ${replicaSet.name} (详情获取失败)`,
          status: 'Current'
        }]
      }
    } else {
      // 如果没有控制器，显示Pod自身信息
      console.log('🔹 独立Pod，显示Pod自身信息')
      historyVersions.value = [{
        revision: 1,
        creationTime: podDetail.value.createdAt || podDetail.value.creationTimestamp,
        replicas: 1,
        description: `独立Pod: ${routeParams.podName}`,
        image: containers.value?.[0]?.image,
        status: podDetail.value.status || podDetail.value.phase || 'Unknown'
      }]
    }
    
    console.log('✅ 历史版本数据处理完成:', historyVersions.value)
    
  } catch (error) {
    console.error('❌ 获取历史版本异常:', error)
    // 如果获取失败，至少显示当前Pod信息
    historyVersions.value = [{
      revision: 1,
      creationTime: podDetail.value.createdAt || podDetail.value.creationTimestamp || new Date().toISOString(),
      replicas: 1,
      description: `当前Pod: ${routeParams.podName} (历史版本获取失败)`,
      status: 'Current',
      error: error.message
    }]
    
    ElMessage.warning('历史版本获取失败，显示当前Pod信息')
  }
}

// 打开容器日志对话框
const handleViewLogs = async (pod) => {
  if (allRelatedPods.value.length === 0) {
    await getRelatedPods()
  }
  
  if (pod && pod.name) {
    selectedContainer.value = pod
  } else if (allRelatedPods.value.length > 0) {
    selectedContainer.value = allRelatedPods.value[0]
  }
  
  dialogStates.logsVisible = true
}

// 处理容器选择变化
const handleContainerChange = (podName) => {
  console.log('🔄 选择Pod:', podName)
  selectedContainerForLogs.value = podName
  if (podName) {
    handleRefreshLogs()
  }
}

// 刷新日志
const handleRefreshLogs = async () => {
  if (!selectedContainerForLogs.value) {
    ElMessage.warning('请先选择Pod')
    return
  }
  
  console.log('🔄 开始刷新日志...')
  console.log('🔍 selectedContainerForLogs.value:', selectedContainerForLogs.value)
  console.log('🔍 allRelatedPods.value:', allRelatedPods.value.map(p => p.name))
  
  // 根据选中的Pod名称找到对应的Pod和容器
  const selectedPod = allRelatedPods.value.find(pod => pod.name === selectedContainerForLogs.value)
  if (!selectedPod) {
    console.error('❌ 未找到选中的Pod:', selectedContainerForLogs.value)
    console.error('❌ 可用的Pod列表:', allRelatedPods.value.map(p => p.name))
    ElMessage.error(`未找到选中的Pod: ${selectedContainerForLogs.value}`)
    return
  }
  
  console.log('✅ 找到选中的Pod:', selectedPod)
  
  // 获取Pod的第一个容器名称
  const firstContainer = selectedPod.containers?.[0]?.name
  if (!firstContainer) {
    console.error('❌ Pod中未找到容器:', selectedPod)
    console.error('❌ containers字段:', selectedPod.containers)
    ElMessage.error(`Pod "${selectedContainerForLogs.value}" 中未找到容器`)
    return
  }
  
  console.log('✅ 找到第一个容器:', firstContainer)
  
  try {
    logsLoading.value = true
    console.log('🔍 获取Pod日志:', {
      podName: selectedContainerForLogs.value,
      container: firstContainer,
      tailLines: logTailLines.value,
      previous: showPreviousLogs.value,
      follow: followLogs.value
    })
    
    // 构建日志查询参数
    const logParams = {
      container: firstContainer
    }
    
    // 添加行数限制（最大1000行）
    logParams.tailLines = logTailLines.value || 1000
    
    // 如果要显示上个容器的日志
    if (showPreviousLogs.value) {
      logParams.previous = true
    }
    
    // 如果要实时跟踪
    if (followLogs.value) {
      logParams.follow = true
    }
    
    const response = await k8sApi.getPodLogs(routeParams.clusterId, routeParams.namespace, selectedContainerForLogs.value, logParams)
    const responseData = response.data || response
    
    console.log('日志API响应:', {
      success: responseData.code === 200,
      dataLength: responseData.data?.length || 0,
      dataType: typeof responseData.data
    })
    
    if (responseData.code === 200) {
      // 处理嵌套的日志数据结构：data.logs 或直接 data
      const logData = responseData.data?.logs || responseData.data || ''
      console.log('📋 日志数据结构解析:', {
        hasDataLogs: !!responseData.data?.logs,
        hasDirectData: !!responseData.data && typeof responseData.data === 'string',
        logDataType: typeof logData,
        logDataLength: logData?.length || 0
      })
      
      currentLogs.value = logData
      lastLogRefreshTime.value = new Date()
      
      const logLineCount = currentLogs.value.split('\n').length
      console.log('✅ 日志获取成功:', logLineCount, '行')
      console.log('📄 日志内容预览:', currentLogs.value.substring(0, 200) + '...')
      
      if (currentLogs.value.trim()) {
        // 静默加载日志，不显示成功提示
        console.log(`✅ 日志加载成功，共 ${logLineCount} 行`)
      } else {
        ElMessage.info('该容器暂无日志输出')
        currentLogs.value = '# 该容器暂无日志输出\n# Container has no log output'
      }
    } else {
      console.error('❌ 日志API返回错误:', responseData.message)
      ElMessage.error(responseData.message || '获取容器日志失败')
      currentLogs.value = ''
    }
  } catch (error) {
    console.error('❌ 获取容器日志异常:', error)
    
    if (error.response?.status === 404) {
      ElMessage.error('容器不存在或日志不可用')
      currentLogs.value = '# 容器不存在或日志不可用\n# Container not found or logs not available'
    } else if (error.response?.status === 400) {
      ElMessage.error('日志参数错误，请检查容器状态')
      currentLogs.value = '# 日志参数错误\n# Invalid log parameters'
    } else {
      ElMessage.error('获取日志失败: ' + (error.message || '网络连接异常'))
      currentLogs.value = '# 日志获取失败\n# Failed to retrieve logs\n# Error: ' + (error.message || 'Network error')
    }
  } finally {
    logsLoading.value = false
  }
}

// 下载日志文件
const handleDownloadLogs = () => {
  if (!selectedContainerForLogs.value) {
    ElMessage.warning('请先选择Pod')
    return
  }
  
  if (!currentLogs.value) {
    ElMessage.warning('请先获取日志')
    return
  }
  
  try {
    // 创建日志文件内容
    const logContent = `# Container Logs
# Pod: ${routeParams.podName}
# Namespace: ${routeParams.namespace}
# Container: ${selectedContainerForLogs.value}
# Generated: ${new Date().toLocaleString('zh-CN')}
# Lines: Last ${logTailLines.value || 1000} lines
# Previous: ${showPreviousLogs.value ? 'Yes' : 'No'}

${currentLogs.value}`
    
    // 创建下载链接
    const blob = new Blob([logContent], { type: 'text/plain;charset=utf-8' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    
    // 生成文件名
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5)
    const filename = `${routeParams.podName}_${selectedContainerForLogs.value}_${timestamp}.log`
    
    link.href = url
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success(`日志文件已下载: ${filename}`)
    console.log('✅ 日志下载成功:', filename)
    
  } catch (error) {
    console.error('❌ 日志下载失败:', error)
    ElMessage.error('日志下载失败: ' + error.message)
  }
}

// 获取YAML内容
const handleViewYaml = async () => {
  try {
    loading.value = true
    console.log('🔍 开始获取YAML内容...')
    console.log('请求参数:', {
      clusterId: routeParams.clusterId,
      namespace: routeParams.namespace,
      podName: routeParams.podName
    })
    
    const response = await k8sApi.getPodYaml(routeParams.clusterId, routeParams.namespace, routeParams.podName)
    const responseData = response.data || response
    
    console.log('YAML API响应结构:', {
      hasResponse: !!responseData,
      code: responseData.code,
      hasData: !!responseData.data,
      dataType: typeof responseData.data,
      dataLength: responseData.data?.length || 0
    })
    
    if (responseData.code === 200) {
      const rawData = responseData.data || ''
      if (typeof rawData === 'string') {
        yamlContent.value = rawData
      } else if (typeof rawData === 'object') {
        yamlContent.value = JSON.stringify(rawData, null, 2)
      } else {
        yamlContent.value = String(rawData)
      }
      console.log('YAML内容检查:', {
        contentLength: yamlContent.value.length,
        firstLine: yamlContent.value.split('\n')[0],
        isString: typeof yamlContent.value === 'string',
        isEmpty: yamlContent.value.trim() === ''
      })
      
      if (yamlContent.value && yamlContent.value.trim()) {
        dialogStates.yamlVisible = true
        console.log('✅ YAML内容获取成功，打开对话框')
        ElMessage.success('YAML内容加载成功')
      } else {
        console.log('⚠️ YAML内容为空')
        ElMessage.warning('YAML内容为空，可能Pod资源不存在')
      }
    } else {
      console.log('❌ YAML API返回错误:', responseData)
      ElMessage.error(responseData.message || '获取YAML内容失败')
    }
  } catch (error) {
    console.error('❌ 获取YAML内容异常:', error)
    
    // 详细的错误处理
    if (error.response) {
      const status = error.response.status
      const statusText = error.response.statusText
      console.log('HTTP错误详情:', {
        status,
        statusText,
        data: error.response.data
      })
      
      switch (status) {
        case 404:
          ElMessage.error('Pod不存在或已被删除')
          break
        case 403:
          ElMessage.error('没有权限访问该Pod的YAML')
          break
        case 500:
          ElMessage.error('服务器内部错误，请稍后重试')
          break
        default:
          ElMessage.error(`获取YAML失败: ${status} ${statusText}`)
      }
    } else if (error.request) {
      console.log('网络请求错误:', error.request)
      ElMessage.error('网络连接失败，请检查网络连接')
    } else {
      console.log('其他错误:', error.message)
      ElMessage.error('获取YAML内容失败: ' + error.message)
    }
  } finally {
    loading.value = false
  }
}

// YAML编辑
const handleEditYaml = async (pod) => {
  try {
    loading.value = true
    console.log('🔍 开始获取Pod YAML内容用于编辑...', pod.name)
    
    // 使用pod参数中的信息或者从路由参数获取
    const clusterId = routeParams.clusterId
    const namespace = pod.namespace || routeParams.namespace
    const podName = pod.name
    
    console.log('请求参数:', { clusterId, namespace, podName })
    
    const response = await k8sApi.getPodYaml(clusterId, namespace, podName)
    const responseData = response.data || response
    
    if (responseData.code === 200) {
      const rawData = responseData.data || ''
      if (typeof rawData === 'string') {
        yamlContent.value = rawData
      } else if (typeof rawData === 'object') {
        yamlContent.value = JSON.stringify(rawData, null, 2)
      } else {
        yamlContent.value = String(rawData)
      }
      
      // 打开YAML编辑对话框
      dialogStates.yamlVisible = true
      console.log('✅ YAML内容加载成功，打开编辑对话框')
    } else {
      console.error('❌ 获取YAML内容失败:', responseData)
      ElMessage.error(responseData.message || '获取YAML内容失败')
    }
  } catch (error) {
    console.error('❌ 获取YAML内容异常:', error)
    ElMessage.error('获取YAML内容失败: ' + (error.message || '网络连接异常'))
  } finally {
    loading.value = false
  }
}

// 容器扩缩容
const handleScale = () => {
  scaleForm.replicas = podDetail.value.replicas || 1
  dialogStates.scaleVisible = true
}

// 处理回滚操作
const handleRollback = (version) => {
  rollbackForm.targetVersion = version.revision
  rollbackForm.versionInfo = version
  dialogStates.rollbackVisible = true
}

// 确认回滚
const handleConfirmRollback = async () => {
  if (!rollbackForm.targetVersion) {
    ElMessage.error('请选择回滚版本')
    return
  }

  // 获取部署名称
  const ownerRefs = podDetail.value.metadata?.ownerReferences ||
                   podDetail.value.spec?.ownerReferences ||
                   podDetail.value.ownerReferences || []

  console.log('🔍 [回滚] 检查ownerReferences:', ownerRefs)

  // 首先查找直接的Deployment引用
  let deployment = ownerRefs.find(ref => ref.kind === 'Deployment')
  let deploymentName = null

  if (deployment) {
    deploymentName = deployment.name
    console.log('✅ [回滚] 找到直接的Deployment:', deploymentName)
  } else {
    // 查找ReplicaSet，然后推断Deployment名称
    const replicaSet = ownerRefs.find(ref => ref.kind === 'ReplicaSet')
    if (replicaSet) {
      // ReplicaSet命名规则: deployment-name-pod-template-hash
      const parts = replicaSet.name.split('-')
      if (parts.length >= 2) {
        deploymentName = parts.slice(0, -1).join('-')
        console.log('✅ [回滚] 从ReplicaSet推断Deployment:', deploymentName)
      }
    }
  }

  // 如果还没找到，尝试从Pod名称推断
  if (!deploymentName) {
    const podNameParts = routeParams.podName.split('-')
    if (podNameParts.length >= 3) {
      deploymentName = podNameParts.slice(0, -2).join('-')
      console.log('✅ [回滚] 从Pod名称推断Deployment:', deploymentName)
    }
  }

  if (!deploymentName) {
    ElMessage.error('无法找到关联的Deployment，无法执行回滚')
    return
  }

  try {
    rollbackLoading.value = true
    console.log(`🔄 开始回滚Deployment ${deploymentName} 到版本 ${rollbackForm.targetVersion}`)

    const rollbackData = {
      revision: rollbackForm.targetVersion
    }

    const response = await k8sApi.rollbackDeployment(
      routeParams.clusterId,
      routeParams.namespace,
      deploymentName,
      rollbackData
    )

    const responseData = response.data || response
    if (responseData.code === 200) {
      ElMessage.success(`成功回滚到版本 ${rollbackForm.targetVersion}`)
      dialogStates.rollbackVisible = false

      // 刷新数据
      await handleQuery(true)
    } else {
      throw new Error(responseData.message || '回滚失败')
    }
  } catch (error) {
    console.error('❌ 回滚失败:', error)
    ElMessage.error(error.message || '回滚操作失败')
  } finally {
    rollbackLoading.value = false
  }
}

// 显示监控弹框
const handleShowMonitoring = async (pod) => {
  selectedPodForMonitoring.value = pod
  dialogStates.monitoringVisible = true
  
  // 加载监控数据
  await loadMonitoringData(pod)
}

// 加载监控数据
const loadMonitoringData = async (pod) => {
  monitoringLoading.value = true
  
  try {
    console.log('🔍 加载Pod监控数据:', pod.name)
    
    // 调用真实的监控API
    const response = await k8sApi.getPodMetrics(routeParams.clusterId, routeParams.namespace, pod.name)
    const responseData = response.data || response
    
    if (responseData.code === 200 && responseData.data) {
      const metrics = responseData.data
      
      // 转换API响应数据格式
      monitoringData.value = {
        cpu: {
          used: metrics.totalUsage?.cpu || '0m',
          limit: metrics.resourceQuota?.cpu || '1000m',
          percentage: Math.round(metrics.usageRate?.cpuRate || 0)
        },
        memory: {
          used: metrics.totalUsage?.memory || '0Mi',
          limit: metrics.resourceQuota?.memory || '512Mi',
          percentage: Math.round(metrics.usageRate?.memoryRate || 0)
        },
        network: {
          rx: '0KB/s',
          tx: '0KB/s'
        },
        disk: {
          used: '0GB',
          limit: '10GB',
          percentage: 0
        },
        containers: metrics.containers || [],
        timestamp: metrics.timestamp
      }
      
      console.log('✅ 监控数据加载完成:', monitoringData.value)
    } else {
      console.warn('⚠️ 监控API返回异常，使用模拟数据:', responseData.message)
      await loadMockMonitoringData(pod)
    }
    
  } catch (error) {
    console.warn('⚠️ 监控API调用失败，使用模拟数据:', error)
    await loadMockMonitoringData(pod)
  } finally {
    monitoringLoading.value = false
  }
}

// 加载模拟监控数据（作为后备方案）
const loadMockMonitoringData = async (pod) => {
  const mockData = {
    cpu: { 
      used: Math.floor(Math.random() * 500) + 100 + 'm', 
      limit: '1000m', 
      percentage: Math.floor(Math.random() * 80) + 10 
    },
    memory: { 
      used: Math.floor(Math.random() * 200) + 50 + 'Mi', 
      limit: '512Mi', 
      percentage: Math.floor(Math.random() * 70) + 15 
    },
    network: { 
      rx: Math.floor(Math.random() * 100) + 10 + 'KB/s', 
      tx: Math.floor(Math.random() * 50) + 5 + 'KB/s' 
    },
    disk: { 
      used: (Math.random() * 5 + 1).toFixed(2) + 'GB', 
      limit: '10GB', 
      percentage: Math.floor(Math.random() * 60) + 10 
    },
    containers: pod.containers?.map(container => ({
      name: container.name,
      usage: {
        cpu: Math.floor(Math.random() * 300) + 50 + 'm',
        memory: Math.floor(Math.random() * 200) + 50 + 'Mi'
      },
      usageRate: {
        cpuRate: Math.floor(Math.random() * 60) + 15,
        memoryRate: Math.floor(Math.random() * 70) + 10
      },
      state: container.state || 'Running'
    })) || [],
    timestamp: new Date().toISOString()
  }
  
  monitoringData.value = mockData
  console.log('✅ 模拟监控数据加载完成:', mockData)
}

// 保存YAML
const handleSaveYaml = async () => {
  try {
    loading.value = true
    console.log('💾 开始保存YAML内容...')
    
    if (!yamlContent.value || !yamlContent.value.trim()) {
      ElMessage.warning('YAML内容不能为空')
      return
    }
    
    // 验证YAML格式
    try {
      // 简单的YAML格式验证
      if (!yamlContent.value.includes('apiVersion') || !yamlContent.value.includes('kind')) {
        ElMessage.warning('YAML格式不正确，缺少必要的apiVersion或kind字段')
        return
      }
    } catch (yamlError) {
      ElMessage.error('YAML格式错误: ' + yamlError.message)
      return
    }
    
    // 尝试使用kubectl apply方式更新
    console.log('使用kubectl apply方式更新资源...')
    const response = await k8sApi.applyYaml(routeParams.clusterId, yamlContent.value)
    const responseData = response.data || response
    
    if (responseData.code === 200) {
      ElMessage.success('YAML保存成功！')
      dialogStates.yamlVisible = false
      
      // 刷新页面数据
      console.log('⚡ 刷新页面数据...')
      await loadPodDetail()
      await getRelatedPods()
      
      console.log('✅ YAML保存并刷新完成')
    } else {
      console.error('❌ 保存YAML失败:', responseData)
      ElMessage.error(responseData.message || '保存YAML失败')
    }
  } catch (error) {
    console.error('❌ 保存YAML异常:', error)
    
    if (error.response && error.response.status === 404) {
      ElMessage.error('kubectl apply接口未实现，请联系管理员')
    } else {
      ElMessage.error('保存YAML失败: ' + (error.message || '网络连接异常'))
    }
  } finally {
    loading.value = false
  }
}

// 打开终端
const handleTerminal = (container) => {
  selectedContainer.value = container
  ElMessage.success('正在跳转到终端页面...')
  
  // 跳转到终端页面
  router.push({
    path: `/k8s/terminal/${routeParams.clusterId}/${routeParams.namespace}/${routeParams.podName}`,
    query: {
      container: container.name || container
    }
  })
}

// 删除Pod
const handleDelete = async (podRow) => {
  const podToDelete = podRow || { name: routeParams.podName }

  try {
    await ElMessageBox.confirm(`确定要删除Pod "${podToDelete.name}" 吗？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    loading.value = true
    const response = await k8sApi.deletePod(routeParams.clusterId, routeParams.namespace, podToDelete.name)
    const responseData = response.data || response

    if (responseData.code === 200) {
      ElMessage.success(`Pod "${podToDelete.name}" 删除成功`)

      // 如果删除的是当前页面的Pod，则返回上一页
      if (podToDelete.name === routeParams.podName) {
        handleGoBack()
      } else {
        // 否则刷新页面数据
        handleQuery(true)
      }
    } else {
      ElMessage.error(responseData.message || 'Pod删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Pod删除失败:', error)
      ElMessage.error('Pod删除失败，请检查网络连接')
    }
  } finally {
    loading.value = false
  }
}

// 获取状态标签类型
const getStatusType = (status) => {
  const statusMap = {
    'Running': 'success',
    'Pending': 'warning',
    'Failed': 'danger',
    'Succeeded': 'success',
    'Unknown': 'info'
  }
  return statusMap[status] || 'info'
}

// 获取资源类型颜色
const getResourceTypeColor = (resourceKind) => {
  const colorMap = {
    'Pod': 'success',
    'ReplicaSet': 'primary', 
    'Deployment': 'warning',
    'Service': 'info',
    'ConfigMap': '',
    'Secret': 'danger',
    'Ingress': 'success'
  }
  return colorMap[resourceKind] || ''
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  return new Date(timestamp).toLocaleString('zh-CN')
}

// 复制到剪贴板
const copyToClipboard = (text, message = '已复制到剪贴板') => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success(message)
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// YAML搜索功能
const searchYamlContent = () => {
  if (!yamlSearchText.value.trim()) {
    yamlSearchResults.value = []
    yamlSearchCurrentIndex.value = -1
    return
  }

  const searchText = yamlSearchText.value.toLowerCase()
  const content = yamlContent.value
  const lines = content.split('\n')
  
  yamlSearchResults.value = []
  
  lines.forEach((line, index) => {
    const lowerLine = line.toLowerCase()
    let startIndex = 0
    
    while (true) {
      const matchIndex = lowerLine.indexOf(searchText, startIndex)
      if (matchIndex === -1) break
      
      yamlSearchResults.value.push({
        lineIndex: index,
        lineNumber: index + 1,
        matchIndex: matchIndex,
        line: line,
        preview: line.trim()
      })
      
      startIndex = matchIndex + searchText.length
    }
  })
  
  yamlSearchCurrentIndex.value = yamlSearchResults.value.length > 0 ? 0 : -1
  
  if (yamlSearchResults.value.length > 0) {
    ElMessage.success(`找到 ${yamlSearchResults.value.length} 个匹配项`)
  } else {
    ElMessage.info('未找到匹配内容')
  }
}

// 跳转到上一个搜索结果
const goToPreviousSearchResult = () => {
  if (yamlSearchResults.value.length === 0) return
  
  yamlSearchCurrentIndex.value = yamlSearchCurrentIndex.value <= 0 
    ? yamlSearchResults.value.length - 1 
    : yamlSearchCurrentIndex.value - 1
    
  scrollToSearchResult()
}

// 跳转到下一个搜索结果
const goToNextSearchResult = () => {
  if (yamlSearchResults.value.length === 0) return
  
  yamlSearchCurrentIndex.value = yamlSearchCurrentIndex.value >= yamlSearchResults.value.length - 1
    ? 0 
    : yamlSearchCurrentIndex.value + 1
    
  scrollToSearchResult()
}

// 滚动到当前搜索结果
const scrollToSearchResult = () => {
  if (yamlSearchCurrentIndex.value === -1 || yamlSearchResults.value.length === 0) return
  
  const currentResult = yamlSearchResults.value[yamlSearchCurrentIndex.value]
  console.log('📍 跳转到搜索结果:', currentResult.lineNumber, '行')
  
  // 通过事件通知CodeEditor滚动到指定行
  nextTick(() => {
    const yamlEditor = document.querySelector('.yaml-dialog .code-editor__textarea')
    if (yamlEditor) {
      const lines = yamlContent.value.split('\n')
      const targetLine = currentResult.lineIndex
      const lineHeight = 21 // 估算行高
      const scrollTop = targetLine * lineHeight
      
      yamlEditor.scrollTop = Math.max(0, scrollTop - 100) // 留一些上边距
    }
  })
}

// 清空搜索
const clearYamlSearch = () => {
  yamlSearchText.value = ''
  yamlSearchResults.value = []
  yamlSearchCurrentIndex.value = -1
}

// 日志搜索功能
const searchLogsContent = () => {
  if (!logsSearchText.value.trim()) {
    logsSearchResults.value = []
    logsSearchCurrentIndex.value = -1
    return
  }

  const searchText = logsSearchText.value.toLowerCase()
  const content = currentLogs.value
  const lines = content.split('\n')
  
  logsSearchResults.value = []
  
  lines.forEach((line, index) => {
    const lowerLine = line.toLowerCase()
    let startIndex = 0
    
    while (true) {
      const matchIndex = lowerLine.indexOf(searchText, startIndex)
      if (matchIndex === -1) break
      
      logsSearchResults.value.push({
        lineIndex: index,
        lineNumber: index + 1,
        matchIndex: matchIndex,
        line: line,
        preview: line.trim()
      })
      
      startIndex = matchIndex + searchText.length
    }
  })
  
  logsSearchCurrentIndex.value = logsSearchResults.value.length > 0 ? 0 : -1
  
  if (logsSearchResults.value.length > 0) {
    ElMessage.success(`找到 ${logsSearchResults.value.length} 个匹配项`)
    scrollToLogsSearchResult()
  } else {
    ElMessage.info('未找到匹配内容')
  }
}

// 跳转到上一个日志搜索结果
const goToPreviousLogsSearchResult = () => {
  if (logsSearchResults.value.length === 0) return
  
  logsSearchCurrentIndex.value = logsSearchCurrentIndex.value <= 0 
    ? logsSearchResults.value.length - 1 
    : logsSearchCurrentIndex.value - 1
    
  scrollToLogsSearchResult()
}

// 跳转到下一个日志搜索结果
const goToNextLogsSearchResult = () => {
  if (logsSearchResults.value.length === 0) return
  
  logsSearchCurrentIndex.value = logsSearchCurrentIndex.value >= logsSearchResults.value.length - 1
    ? 0 
    : logsSearchCurrentIndex.value + 1
    
  scrollToLogsSearchResult()
}

// 滚动到当前日志搜索结果
const scrollToLogsSearchResult = () => {
  if (logsSearchCurrentIndex.value === -1 || logsSearchResults.value.length === 0) return
  
  const currentResult = logsSearchResults.value[logsSearchCurrentIndex.value]
  console.log('📍 跳转到日志搜索结果:', currentResult.lineNumber, '行')
  
  // 通过事件通知CodeEditor滚动到指定行
  nextTick(() => {
    const logsEditor = document.querySelector('.logs-tab-content .code-editor__textarea')
    if (logsEditor) {
      const targetLine = currentResult.lineIndex
      const lineHeight = 21 // 估算行高
      const scrollTop = targetLine * lineHeight
      
      logsEditor.scrollTop = Math.max(0, scrollTop - 100) // 留一些上边距
    }
  })
}

// 清空日志搜索
const clearLogsSearch = () => {
  logsSearchText.value = ''
  logsSearchResults.value = []
  logsSearchCurrentIndex.value = -1
}

// 从Pod详情生成YAML内容
const generateYamlFromDetail = () => {
  try {
    console.log('🔧 从Pod详情生成YAML...')
    console.log('Pod详情数据:', podDetail.value)
    
    if (!podDetail.value || Object.keys(podDetail.value).length === 0) {
      ElMessage.warning('Pod详情数据为空，无法生成YAML')
      return
    }
    
    // 构造基本的Pod YAML结构
    const yamlObject = {
      apiVersion: 'v1',
      kind: 'Pod',
      metadata: {
        name: routeParams.podName,
        namespace: routeParams.namespace,
        labels: podDetail.value.labels || {},
        annotations: podDetail.value.annotations || {}
      },
      spec: {
        containers: (containers.value || []).map(container => ({
          name: container.name,
          image: container.image,
          ports: container.ports || [],
          env: container.env || [],
          resources: container.resources || {},
          volumeMounts: container.volumeMounts || []
        })),
        volumes: podDetail.value.volumes || [],
        restartPolicy: podDetail.value.spec?.restartPolicy || 'Always',
        nodeName: podDetail.value.nodeName,
        serviceAccount: podDetail.value.spec?.serviceAccount,
        serviceAccountName: podDetail.value.spec?.serviceAccountName,
        nodeSelector: podDetail.value.spec?.nodeSelector || {},
        tolerations: podDetail.value.spec?.tolerations || [],
        affinity: podDetail.value.spec?.affinity
      },
      status: {
        phase: podDetail.value.status || podDetail.value.phase,
        podIP: podDetail.value.podIP,
        hostIP: podDetail.value.hostIP,
        startTime: podDetail.value.startTime || podDetail.value.createdAt,
        conditions: podDetail.value.conditions || []
      }
    }
    
    // 清理undefined值
    const cleanObject = (obj) => {
      if (Array.isArray(obj)) {
        return obj.map(cleanObject).filter(item => item !== undefined)
      } else if (obj !== null && typeof obj === 'object') {
        const cleaned = {}
        Object.keys(obj).forEach(key => {
          const value = cleanObject(obj[key])
          if (value !== undefined && value !== null && value !== '' && 
              !(Array.isArray(value) && value.length === 0) &&
              !(typeof value === 'object' && Object.keys(value).length === 0)) {
            cleaned[key] = value
          }
        })
        return Object.keys(cleaned).length > 0 ? cleaned : undefined
      }
      return obj
    }
    
    const cleanedYaml = cleanObject(yamlObject)
    
    // 转换为YAML字符串
    const yamlString = `# Pod YAML (从详情生成)
# 名称: ${routeParams.podName}
# 命名空间: ${routeParams.namespace}
# 生成时间: ${new Date().toLocaleString('zh-CN')}

apiVersion: v1
kind: Pod
metadata:
  name: ${routeParams.podName}
  namespace: ${routeParams.namespace}
  labels:${Object.keys(cleanedYaml.metadata.labels || {}).length > 0 ? 
    Object.entries(cleanedYaml.metadata.labels).map(([k, v]) => `\n    ${k}: ${v}`).join('') : 
    '\n    # No labels'}
spec:
  containers:${cleanedYaml.spec.containers.map((container, index) => `
  - name: ${container.name}
    image: ${container.image}${container.ports?.length > 0 ? `
    ports:${container.ports.map(port => `
    - containerPort: ${port.containerPort}
      protocol: ${port.protocol || 'TCP'}`).join('')}` : ''}${container.env?.length > 0 ? `
    env:${container.env.map(envVar => `
    - name: ${envVar.name}
      value: ${envVar.value || envVar.valueFrom ? JSON.stringify(envVar.valueFrom) : ''}`).join('')}` : ''}`).join('')}
  restartPolicy: ${cleanedYaml.spec.restartPolicy || 'Always'}${cleanedYaml.spec.nodeName ? `
  nodeName: ${cleanedYaml.spec.nodeName}` : ''}${Object.keys(cleanedYaml.spec.nodeSelector || {}).length > 0 ? `
  nodeSelector:${Object.entries(cleanedYaml.spec.nodeSelector).map(([k, v]) => `
    ${k}: ${v}`).join('')}` : ''}
status:
  phase: ${cleanedYaml.status.phase || 'Unknown'}${cleanedYaml.status.podIP ? `
  podIP: ${cleanedYaml.status.podIP}` : ''}${cleanedYaml.status.hostIP ? `
  hostIP: ${cleanedYaml.status.hostIP}` : ''}${cleanedYaml.status.startTime ? `
  startTime: ${cleanedYaml.status.startTime}` : ''}
`
    
    // 根据调用上下文设置不同的YAML内容
    if (yamlTabLoading.value) {
      yamlTabContent.value = yamlString
      console.log('✅ YAML标签页内容生成成功')
    } else {
      yamlContent.value = yamlString
      dialogStates.yamlVisible = true
      ElMessage.success('YAML已从Pod详情生成')
      console.log('✅ YAML对话框生成成功')
    }
    
  } catch (error) {
    console.error('❌ 生成YAML失败:', error)
    ElMessage.error('生成YAML失败: ' + error.message)
  }
}

// 返回上一页
const handleGoBack = () => {
  router.back()
}
</script>

<template>
  <div class="k8s-pod-management">
    <el-card shadow="hover" class="pod-card">
      <!-- 页面头部 -->
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button type="success" :icon="ArrowLeft" size="small" @click="handleGoBack">返回</el-button>
            <span class="page-title">容器详情 - {{ workloadName }}</span>
          </div>
          <div class="header-actions">
            <el-button type="success" :icon="Refresh" size="small" @click="() => handleQuery(true)">刷新</el-button>
            <el-button type="primary" :icon="View" size="small" @click="handleViewYaml">查看YAML</el-button>
          </div>
        </div>
      </template>

      <!-- 基本信息区域 -->
      <div class="info-section">
        <el-row :gutter="24">
          <!-- 左侧基本信息 -->
          <el-col :span="12">
            <el-card class="info-card" header="基本信息">
              <el-descriptions :column="1" border size="small">
                <el-descriptions-item label="名称">{{ workloadName }}</el-descriptions-item>
                <el-descriptions-item label="命名空间">{{ routeParams.namespace }}</el-descriptions-item>
                <el-descriptions-item label="节点选择器">
                  <div v-if="podDetail.spec?.nodeSelector && Object.keys(podDetail.spec.nodeSelector).length > 0">
                    <el-tag
                      v-for="(value, key) in podDetail.spec.nodeSelector"
                      :key="key"
                      size="small"
                      style="margin: 2px;"
                      @click="copyToClipboard(`${key}=${value}`, '节点选择器已复制')"
                    >
                      {{ key }}={{ value }}
                    </el-tag>
                  </div>
                  <span v-else>-</span>
                </el-descriptions-item>
                <el-descriptions-item label="标签">
                  <div v-if="Object.keys(podLabels).length > 0" class="labels-container">
                    <div class="labels-display">
                      <el-tag
                        v-for="(value, key) in visibleLabels"
                        :key="key"
                        size="small"
                        style="margin: 2px;"
                        @click="copyToClipboard(`${key}=${value}`, '标签已复制')"
                      >
                        {{ key }}={{ value }}
                      </el-tag>
                      
                      <el-button
                        v-if="hasMoreLabels"
                        type="text"
                        size="small"
                        @click="toggleLabelsExpanded"
                        style="margin-left: 8px; color: #409eff;"
                      >
                        {{ labelsExpanded ? '折叠' : `展开 (+${hiddenLabelsCount})` }}
                      </el-button>
                    </div>
                  </div>
                  <span v-else>-</span>
                </el-descriptions-item>
              </el-descriptions>
            </el-card>
          </el-col>

          <!-- 右侧状态信息 -->
          <el-col :span="12">
            <el-card class="info-card" header="状态信息">
              <el-descriptions :column="1" border size="small">
                <el-descriptions-item label="创建时间">{{ formatTime(podCreationTime) }}</el-descriptions-item>
                <el-descriptions-item label="重启策略">{{ podDetail.spec?.restartPolicy || '-' }}</el-descriptions-item>
                <el-descriptions-item label="运行时间">{{ formatRunningTime(podDetail.runningTime) }}</el-descriptions-item>
                <el-descriptions-item label="期望Pod数量">{{ expectedPodCount }}</el-descriptions-item>
                <el-descriptions-item label="状态">
                  <el-tag :type="getStatusType(podStatus)" effect="dark">{{ podStatus }}</el-tag>
                </el-descriptions-item>
              </el-descriptions>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 标签页内容 -->
      <div class="tabs-section">
        <el-tabs v-model="activeTab" type="border-card" @tab-change="handleTabChange">
          <!-- 容器组标签页 -->
          <el-tab-pane label="容器组" name="containers">
            <el-table
              :data="allRelatedPods"
              v-loading="loading"
              stripe
              style="width: 100%"
              size="small"
            >
              <el-table-column prop="name" label="Pod名称" min-width="200">
                <template #default="{ row }">
                  <div class="pod-name-container">
                    <img src="@/assets/image/k8s.svg" alt="k8s" class="k8s-icon" width="16" height="16" />
                    <span class="pod-name">{{ row.name }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" min-width="100">
                <template #default="{ row }">
                  <el-tag :type="getPodStatusTag(row.status)" size="small">
                    {{ getPodStatusText(row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="restartCount" label="重启次数" min-width="100">
                <template #default="{ row }">
                  <el-tag size="small" :type="row.restartCount > 0 ? 'warning' : 'success'">
                    {{ row.restartCount }}
                  </el-tag>
                </template>
              </el-table-column>
              
              <el-table-column label="监控" min-width="60" align="center">
                <template #default="{ row }">
                  <el-button 
                    :icon="Monitor" 
                    size="small" 
                    circle 
                    type="primary" 
                    @click="handleShowMonitoring(row)"
                    title="查看监控"
                  />
                </template>
              </el-table-column>
              
              <el-table-column prop="nodeName" label="节点" min-width="150"></el-table-column>
              
              <el-table-column prop="podIP" label="Pod IP" min-width="120"></el-table-column>
              
              <el-table-column prop="age" label="运行时间" min-width="120">
                <template #default="{ row }">
                  <div class="time-info">
                    <el-tag size="small" type="info">{{ row.age }}</el-tag>
                    <div v-if="row.runningTime" class="running-time">
                      运行: {{ formatRunningTime(row.runningTime) }}
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="操作" min-width="380" fixed="right">
                <template #default="{ row }">
                  <div class="operation-buttons-container">
                    <el-button size="small" type="warning" :icon="Edit" v-authority="['k8s:workload:edityaml']" @click="handleEditYaml(row)">YAML编辑</el-button>
                    <el-button size="small" type="primary" :icon="Document" v-authority="['k8s:workload:podlog']" @click="handleViewLogs(row)">日志</el-button>
                    <el-button size="small" type="success" v-authority="['k8s:workload:terminal']"  @click="handleTerminal(row)">
                      <img src="@/assets/image/终端.svg" alt="terminal" class="custom-icon" />
                      终端
                    </el-button>
                    <el-button size="small" type="danger" :icon="Delete" v-authority="['k8s:workload:poddelete']" @click="handleDelete(row)">删除</el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <!-- 历史版本标签页 -->
          <el-tab-pane label="历史版本" name="history">
            <el-table
              :data="historyVersions"
              stripe
              style="width: 100%"
              size="small"
              v-loading="loading"
            >
              <el-table-column prop="revision" label="版本" min-width="80" />
              <el-table-column prop="description" label="描述" min-width="250" />
              <el-table-column prop="status" label="状态" min-width="100">
                <template #default="{ row }">
                  <el-tag :type="getStatusType(row.status)" effect="dark" size="small">{{ row.status }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="image" label="镜像" min-width="200">
                <template #default="{ row }">
                  <div class="image-cell">
                    <span v-if="row.image" class="image-info">{{ row.image }}</span>
                    <span v-else>-</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="replicas" label="副本数" min-width="100" />
              <el-table-column prop="creationTime" label="创建时间" min-width="180">
                <template #default="{ row }">
                  {{ formatTime(row.creationTime) }}
                </template>
              </el-table-column>
              <el-table-column prop="changeDescription" label="变更描述" min-width="200">
                <template #default="{ row }">
                  <span v-if="row.changeDescription">{{ row.changeDescription }}</span>
                  <span v-else>-</span>
                </template>
              </el-table-column>
              <el-table-column label="操作" min-width="150" fixed="right">
                <template #default="{ row }">
                  <el-button size="small" v-authority="['k8s:workload:rollback_version']" type="primary" @click="handleRollback(row)">回滚到此版本</el-button>
                </template>
              </el-table-column>
            </el-table>
            <div v-if="!loading && historyVersions.length === 0" style="padding: 20px; text-align: center;">
              <el-empty description="暂无历史版本数据">
                <div style="margin-top: 16px;">
                  <el-button type="primary" @click="getHistoryVersions">重新获取历史版本</el-button>
                </div>
              </el-empty>
            </div>
          </el-tab-pane>

          <!-- 事件标签页 -->
          <el-tab-pane label="事件" name="events">
            <el-table
              :data="events"
              stripe
              style="width: 100%"
              size="small"
              v-loading="loading"
            >
              <el-table-column prop="type" label="事件类型" min-width="100">
                <template #default="{ row }">
                  <el-tag :type="row.type === 'Warning' ? 'warning' : 'success'" effect="dark" size="small">
                    {{ row.type }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="来源" min-width="150">
                <template #default="{ row }">
                  <span>{{ row.source || '未知' }}</span>
                </template>
              </el-table-column>
              <el-table-column label="次数" min-width="80">
                <template #default="{ row }">
                  <el-tag size="small" type="info">{{ row.count || 1 }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="reason" label="原因" min-width="150" />
              <el-table-column prop="message" label="消息" min-width="300" />
              <el-table-column label="最后时间" min-width="180">
                <template #default="{ row }">
                  {{ formatTime(row.lastTime || row.firstTime) }}
                </template>
              </el-table-column>
            </el-table>
            <div v-if="!loading && events.length === 0" style="padding: 20px; text-align: center;">
              <el-empty description="暂无事件数据">
                <div style="margin-top: 16px;">
                  <el-button type="primary" @click="getEvents">重新获取事件</el-button>
                  <el-button type="success" @click="showPodDetailEvents" style="margin-left: 8px;">显示Pod详情事件</el-button>
                </div>
              </el-empty>
            </div>
          </el-tab-pane>

          <!-- 容器伸缩标签页 -->
          <el-tab-pane label="容器伸缩" name="scale">
            <div style="padding: 20px;">
              <el-button type="primary" @click="handleScale">容器扩缩容</el-button>
            </div>
          </el-tab-pane>

          <!-- YAML标签页 -->
          <el-tab-pane label="yaml" name="yaml">
            <div v-loading="yamlTabLoading" style="padding: 20px;">
              <div style="margin-bottom: 10px; text-align: right;">
                <el-button type="primary" size="small" :icon="Copy" @click="copyToClipboard(yamlTabContent, 'YAML已复制')">复制YAML</el-button>
              </div>
              <CodeEditor 
                v-model="yamlTabContent" 
                language="yaml" 
                height="500px"
                :readonly="true"
                fontSize="12px"
              />
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-card>

    <!-- 日志查看对话框 -->
    <LogViewerDialog
      v-model="dialogStates.logsVisible"
      :cluster-id="routeParams.clusterId"
      :namespace="routeParams.namespace"
      :pod-name="selectedContainer?.name"
      :containers="selectedContainer?.containers || []"
      @close="dialogStates.logsVisible = false"
    />

    <!-- YAML查看/编辑对话框 -->
    <el-dialog
      v-model="dialogStates.yamlVisible"
      title="YAML内容"
      width="50%"
      class="yaml-dialog"
      @closed="clearYamlSearch"
    >
      <!-- 搜索工具栏 -->
      <div class="yaml-search-toolbar">
        <el-row :gutter="12" type="flex" align="middle">
          <el-col :span="8">
            <el-input
              v-model="yamlSearchText"
              placeholder="搜索YAML内容..."
              :prefix-icon="Search"
              clearable
              @keyup.enter="searchYamlContent"
              @clear="clearYamlSearch"
              size="small"
              style="width: 280px;"
            />
          </el-col>
          <el-col :span="8">
            <el-button-group size="small">
              <el-button type="warning" :icon="Search" @click="searchYamlContent">搜索</el-button>
              <el-button type="success" :icon="ArrowUp" @click="goToPreviousSearchResult" :disabled="yamlSearchResults.length === 0"></el-button>
              <el-button type="primary" :icon="ArrowDown" @click="goToNextSearchResult" :disabled="yamlSearchResults.length === 0"></el-button>
            </el-button-group>
          </el-col>
          <el-col :span="8">
            <el-text v-if="yamlSearchResults.length > 0" type="info" size="small">
              {{ yamlSearchCurrentIndex + 1 }} / {{ yamlSearchResults.length }}
            </el-text>
          </el-col>
        </el-row>
      </div>
      
      <div class="yaml-container">
        <CodeEditor 
          v-model="yamlContent" 
          language="yaml" 
          height="500px"
          fontSize="12px"
          :searchText="yamlSearchText"
          :searchResults="yamlSearchResults"
          :currentSearchIndex="yamlSearchCurrentIndex"
        />
      </div>
      <template #footer>
        <el-button @click="copyToClipboard(yamlContent, 'YAML已复制')">复制YAML</el-button>
        <el-button type="primary" @click="handleSaveYaml">保存</el-button>
      </template>
    </el-dialog>

    <!-- 容器扩缩容对话框 -->
    <el-dialog
      v-model="dialogStates.scaleVisible"
      title="容器扩缩容"
      width="400px"
    >
      <el-form :model="scaleForm" label-width="80px">
        <el-form-item label="副本数">
          <el-input-number
            v-model="scaleForm.replicas"
            :min="0"
            :max="100"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogStates.scaleVisible = false">取消</el-button>
        <el-button type="primary" @click="ElMessage.success('扩缩容功能开发中')">确定</el-button>
      </template>
    </el-dialog>

    <!-- 回滚确认对话框 -->
    <el-dialog
      v-model="dialogStates.rollbackVisible"
      title="回滚确认"
      width="500px"
    >
      <div class="rollback-confirmation">
        <el-alert
          title="回滚操作风险提示"
          type="warning"
          description="回滚操作将会替换当前运行的版本，请确认您要执行此操作。"
          show-icon
          :closable="false"
          style="margin-bottom: 20px;"
        />

        <el-descriptions title="回滚详情" :column="1" border>
          <el-descriptions-item label="目标版本">{{ rollbackForm.targetVersion }}</el-descriptions-item>
          <el-descriptions-item label="版本描述">{{ rollbackForm.versionInfo.description }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatTime(rollbackForm.versionInfo.creationTime) }}</el-descriptions-item>
          <el-descriptions-item label="镜像" v-if="rollbackForm.versionInfo.image">{{ rollbackForm.versionInfo.image }}</el-descriptions-item>
        </el-descriptions>
      </div>

      <template #footer>
        <el-button @click="dialogStates.rollbackVisible = false">取消</el-button>
        <el-button
          type="danger"
          :loading="rollbackLoading"
          @click="handleConfirmRollback"
        >
          确认回滚
        </el-button>
      </template>
    </el-dialog>

    <!-- 终端对话框 -->
    <el-dialog
      v-model="dialogStates.terminalVisible"
      :title="`终端 - ${selectedContainer?.name || ''}`"
      width="80%"
      class="terminal-dialog"
    >
      <div class="terminal-container">
        <div class="terminal-content">
          <div class="terminal-header">终端连接功能开发中...</div>
          <div class="terminal-body">
            <p>即将支持Web终端连接功能</p>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 监控弹框 -->
    <el-dialog
      v-model="dialogStates.monitoringVisible"
      :title="`Pod监控 - ${selectedPodForMonitoring?.name || ''}`"
      width="60%"
      class="monitoring-dialog"
    >
      <div class="monitoring-container" v-loading="monitoringLoading">
        <el-row :gutter="16" style="margin-bottom: 16px;">
          <el-col :span="12">
            <el-card class="metric-card">
              <div class="metric-header">
                <span class="metric-title">CPU使用率</span>
                <span class="metric-value">{{ monitoringData.cpu.percentage }}%</span>
              </div>
              <el-progress 
                :percentage="monitoringData.cpu.percentage" 
                :color="monitoringData.cpu.percentage > 80 ? '#F56C6C' : monitoringData.cpu.percentage > 60 ? '#E6A23C' : '#67C23A'"
                :stroke-width="8"
              />
              <div class="metric-details">
                <span>已用: {{ monitoringData.cpu.used }}</span>
                <span>限制: {{ monitoringData.cpu.limit }}</span>
              </div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card class="metric-card">
              <div class="metric-header">
                <span class="metric-title">内存使用率</span>
                <span class="metric-value">{{ monitoringData.memory.percentage }}%</span>
              </div>
              <el-progress 
                :percentage="monitoringData.memory.percentage" 
                :color="monitoringData.memory.percentage > 80 ? '#F56C6C' : monitoringData.memory.percentage > 60 ? '#E6A23C' : '#67C23A'"
                :stroke-width="8"
              />
              <div class="metric-details">
                <span>已用: {{ monitoringData.memory.used }}</span>
                <span>限制: {{ monitoringData.memory.limit }}</span>
              </div>
            </el-card>
          </el-col>
        </el-row>
        
        <el-row :gutter="16">
          <el-col :span="12">
            <el-card class="metric-card">
              <div class="metric-header">
                <span class="metric-title">网络流量</span>
              </div>
              <div class="network-metrics">
                <div class="network-item">
                  <span class="network-label">入站:</span>
                  <span class="network-value">{{ monitoringData.network.rx }}</span>
                </div>
                <div class="network-item">
                  <span class="network-label">出站:</span>
                  <span class="network-value">{{ monitoringData.network.tx }}</span>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card class="metric-card">
              <div class="metric-header">
                <span class="metric-title">磁盘使用率</span>
                <span class="metric-value">{{ monitoringData.disk.percentage }}%</span>
              </div>
              <el-progress 
                :percentage="monitoringData.disk.percentage" 
                :color="monitoringData.disk.percentage > 80 ? '#F56C6C' : monitoringData.disk.percentage > 60 ? '#E6A23C' : '#67C23A'"
                :stroke-width="8"
              />
              <div class="metric-details">
                <span>已用: {{ monitoringData.disk.used }}</span>
                <span>限制: {{ monitoringData.disk.limit }}</span>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 容器级别监控数据 -->
        <el-card v-if="monitoringData.containers && monitoringData.containers.length > 0" style="margin-top: 16px;">
          <template #header>
            <div class="card-header">
              <span>容器监控详情</span>
              <el-tag v-if="monitoringData.timestamp" type="info" size="small">
                更新时间: {{ formatTimestamp(monitoringData.timestamp) }}
              </el-tag>
            </div>
          </template>
          
          <el-table :data="monitoringData.containers" size="small" stripe>
            <el-table-column prop="name" label="容器名称" min-width="120" />
            <el-table-column label="CPU使用" min-width="100">
              <template #default="{ row }">
                <div class="container-metric">
                  <span class="metric-text">{{ row.usage?.cpu || '0m' }}</span>
                  <el-progress 
                    v-if="row.usageRate?.cpuRate" 
                    :percentage="Math.round(row.usageRate.cpuRate)" 
                    :stroke-width="6"
                    :show-text="false"
                    :color="row.usageRate.cpuRate > 80 ? '#F56C6C' : row.usageRate.cpuRate > 60 ? '#E6A23C' : '#67C23A'"
                  />
                </div>
              </template>
            </el-table-column>
            <el-table-column label="内存使用" min-width="100">
              <template #default="{ row }">
                <div class="container-metric">
                  <span class="metric-text">{{ row.usage?.memory || '0Mi' }}</span>
                  <el-progress 
                    v-if="row.usageRate?.memoryRate" 
                    :percentage="Math.round(row.usageRate.memoryRate)" 
                    :stroke-width="6"
                    :show-text="false"
                    :color="row.usageRate.memoryRate > 80 ? '#F56C6C' : row.usageRate.memoryRate > 60 ? '#E6A23C' : '#67C23A'"
                  />
                </div>
              </template>
            </el-table-column>
            <el-table-column label="CPU使用率" min-width="80">
              <template #default="{ row }">
                <el-tag 
                  :type="row.usageRate?.cpuRate > 80 ? 'danger' : row.usageRate?.cpuRate > 60 ? 'warning' : 'success'"
                  size="small"
                >
                  {{ Math.round(row.usageRate?.cpuRate || 0) }}%
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="内存使用率" min-width="80">
              <template #default="{ row }">
                <el-tag 
                  :type="row.usageRate?.memoryRate > 80 ? 'danger' : row.usageRate?.memoryRate > 60 ? 'warning' : 'success'"
                  size="small"
                >
                  {{ Math.round(row.usageRate?.memoryRate || 0) }}%
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" min-width="80">
              <template #default="{ row }">
                <el-tag 
                  :type="row.state === 'Running' ? 'success' : row.state === 'Waiting' ? 'warning' : 'danger'"
                  size="small"
                >
                  {{ row.state || 'Unknown' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
      
      <template #footer>
        <el-button @click="dialogStates.monitoringVisible = false">关闭</el-button>
        <el-button type="primary" @click="loadMonitoringData(selectedPodForMonitoring)" :loading="monitoringLoading">刷新数据</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.k8s-pod-management {
  padding: 20px;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.pod-card {
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

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.info-section {
  margin-bottom: 20px;
}

.info-card {
  height: 100%;
  border-radius: 6px;
}

.info-card :deep(.el-card__header) {
  background: rgba(103, 126, 234, 0.05);
  border-bottom: 1px solid rgba(103, 126, 234, 0.1);
  font-weight: 600;
  padding: 12px 16px;
}

.tabs-section {
  margin-top: 20px;
}

.labels-container {
  max-width: 100%;
}

/* Pod名称容器样式 */
.pod-name-container {
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
}

.pod-name-container:hover {
  color: #409eff;
  transform: translateY(-1px);
}

.k8s-icon {
  flex-shrink: 0;
}

.pod-name {
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.labels-display {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 2px;
}

.image-link {
  color: #409eff;
  padding: 0;
  font-size: 12px;
  text-decoration: none;
}

.image-link:hover {
  text-decoration: underline;
}

.logs-dialog .logs-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #000;
  color: #fff;
  height: 400px;
  overflow: auto;
}

.logs-content {
  padding: 12px;
  margin: 0;
  white-space: pre-wrap;
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.yaml-search-toolbar {
  padding: 12px;
  border-bottom: 1px solid #ebeef5;
  background-color: #fafafa;
  margin-bottom: 16px;
}

.yaml-search-toolbar .el-button-group .el-button {
  padding: 7px 12px;
}

.yaml-search-toolbar .el-input {
  font-size: 14px;
}

.yaml-dialog .yaml-container {
  margin-top: 0;
}

.logs-editor-container {
  background-color: #1e1e1e;
  border-radius: 4px;
  overflow: hidden;
}

.yaml-dialog .yaml-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.yaml-editor :deep(.el-textarea__inner) {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  background-color: #f8f9fa;
}

.terminal-dialog .terminal-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #000;
  color: #fff;
  height: 400px;
  padding: 12px;
}

.terminal-header {
  color: #00ff00;
  margin-bottom: 10px;
}

.terminal-body {
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

:deep(.el-descriptions-item__label) {
  font-weight: 600;
}

:deep(.el-tag) {
  cursor: pointer;
}

:deep(.el-tag):hover {
  transform: scale(1.05);
}

:deep(.el-card) {
  border-radius: 8px;
}

:deep(.el-table) {
  border-radius: 8px;
}

/* 日志标签页样式 */
.logs-tab-content {
  padding: 16px;
}

.logs-controls {
  background-color: #f8f9fa;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.logs-actions {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.logs-display {
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  background-color: #000;
  height: 680px;
  overflow: hidden;
}

.logs-placeholder {
  padding: 40px;
  text-align: center;
  color: #00ff00;
  background-color: #000;
  height: 100%;
  border-radius: 8px;
}

.logs-content-display {
  height: 100%;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: #1a1a1a;
  border-bottom: 1px solid #333;
  border-radius: 8px 8px 0 0;
}

.container-name {
  font-weight: 600;
  color: #00ff00;
  font-size: 14px;
}

.logs-info {
  color: #888;
  font-size: 12px;
}

.logs-text {
  padding: 16px;
  margin: 0;
  font-family: 'Courier New', Monaco, 'Lucida Console', monospace;
  font-size: 12px;
  line-height: 1.4;
  background-color: #000;
  color: #00ff00;
  white-space: pre-wrap;
  word-break: break-all;
  height: calc(500px - 48px);
  overflow-y: auto;
  border-radius: 0 0 8px 8px;
}

.logs-empty {
  padding: 40px;
  text-align: center;
  color: #00ff00;
  background-color: #000;
  height: 100%;
  border-radius: 8px;
}

.image-info {
  font-family: 'Courier New', Monaco, monospace;
  font-size: 12px;
  color: #666;
}

/* 镜像单元格样式 */
.image-cell {
  word-break: break-all;
  white-space: normal;
  line-height: 1.4;
  max-width: 300px;
}

.image-cell .image-link {
  word-break: break-all;
  white-space: normal;
  text-align: left;
  height: auto;
  line-height: 1.4;
  padding: 4px 8px;
}

/* 操作按钮布局样式 */
.operation-buttons-container {
  display: flex;
  flex-direction: row;
  gap: 4px;
  flex-wrap: nowrap;
  justify-content: flex-start;
}

/* 自定义图标样式 */
.custom-icon {
  width: 16px;
  height: 16px;
  margin-right: 4px;
  vertical-align: middle;
  filter: brightness(0) invert(1);
}

/* 监控弹框样式 */
.monitoring-dialog .el-dialog__body {
  padding: 20px;
}

.monitoring-container {
  min-height: 300px;
}

.metric-card {
  height: 140px;
  margin-bottom: 16px;
}

.metric-card :deep(.el-card__body) {
  padding: 16px;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.metric-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.metric-value {
  font-size: 16px;
  font-weight: bold;
  color: #409EFF;
}

.metric-details {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #666;
  margin-top: 8px;
}

.network-metrics {
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 60px;
}

.network-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.network-label {
  font-size: 14px;
  color: #666;
}

.network-value {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

/* 容器监控表格样式 */
.container-metric {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.metric-text {
  font-size: 12px;
  color: #606266;
  margin-bottom: 2px;
}刚才修改

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center
  ;
}

.card-header span {
  font-weight: 500;
}

/* 回滚确认对话框样式 */
.rollback-confirmation {
  padding: 10px 0;
}

.rollback-confirmation .el-descriptions {
  margin-top: 10px;
}

.rollback-confirmation .el-alert {
  border-radius: 8px;
}
</style>