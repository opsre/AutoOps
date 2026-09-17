<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  Edit,
  Delete,
  View,
  Setting,
  Warning,
  CircleCheck,
  Monitor,
  Lock,
  Unlock,
  SwitchButton,
  DocumentCopy
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import NodesTable from './nodes/NodesTable.vue'
import { openYamlDialog } from './services/yamlDialog'

const router = useRouter()

const loading = ref(false)
const queryParams = reactive({
  nodeName: '',
  status: '',
  role: ''
})

const statusOptions = [
  { label: '全部', value: '' },
  { label: '就绪', value: 'Ready' },
  { label: '未就绪', value: 'NotReady' },
  { label: '不可调度', value: 'SchedulingDisabled' },
  { label: '未知', value: 'Unknown' }
]

const roleOptions = [
  { label: '全部', value: '' },
  { label: 'Master', value: 'master' },
  { label: 'Worker', value: 'worker' },
  { label: 'Control-plane', value: 'control-plane' }
]

const tableData = ref([])
const selectedClusterId = ref('') // 默认集群ID，初始为空
const clusterList = ref([]) // 集群列表

// 对话框状态
const taintDialogVisible = ref(false)
const labelDialogVisible = ref(false)
const resourcesDialogVisible = ref(false)
const drainDialogVisible = ref(false)
const nodeLabelsDialogVisible = ref(false)
// 当前操作的节点
const currentNode = ref({})

// 污点管理表单
const taintForm = reactive({
  operation: 'add', // add | remove
  key: '',
  value: '',
  effect: 'NoSchedule', // NoSchedule | PreferNoSchedule | NoExecute
  selectedTaint: null // 选中的要移除的污点
})

// 标签管理表单
const labelForm = reactive({
  operation: 'add', // add | remove
  key: '',
  value: '',
  selectedLabel: null // 选中的要移除的标签
})

// 驱逐配置表单
const drainForm = reactive({
  force: false,
  deleteLocalData: false,
  ignoreDaemonSets: true,
  gracePeriodSeconds: 30
})

// 节点资源详情
const nodeResources = ref({})

const effectOptions = [
  { 
    label: 'NoSchedule - 禁止调度', 
    value: 'NoSchedule',
    description: '不允许新Pod调度到此节点，但不影响已运行的Pod'
  },
  { 
    label: 'PreferNoSchedule - 尽量避免调度', 
    value: 'PreferNoSchedule',
    description: '尽量避免新Pod调度到此节点，但在资源不足时仍可调度'
  },
  { 
    label: 'NoExecute - 禁止执行', 
    value: 'NoExecute',
    description: '不允许新Pod调度，且会驱逐不能容忍此污点的已运行Pod'
  }
]

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
      
      // 如果有集群且当前没有选中集群，默认选择第一个在线集群
      if (clusterList.value.length > 0 && !selectedClusterId.value) {
        const onlineCluster = clusterList.value.find(cluster => cluster.status === 2) // 2表示在线
        selectedClusterId.value = onlineCluster ? onlineCluster.id : clusterList.value[0].id
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

const handleQuery = async () => {
  try {
    // 检查是否选择了集群
    if (!selectedClusterId.value) {
      ElMessage.warning('请选择一个集群')
      return
    }

    loading.value = true

    let response

    // 根据不同的查询条件调用不同的API
    if (queryParams.nodeName) {
      // 根据节点名称搜索
      response = await k8sApi.searchNodesByName(selectedClusterId.value, queryParams.nodeName)
    } else if (queryParams.role) {
      // 根据节点角色搜索
      response = await k8sApi.searchNodesByRole(selectedClusterId.value, queryParams.role)
    } else if (queryParams.status) {
      // 根据节点状态搜索
      response = await k8sApi.searchNodesByStatus(selectedClusterId.value, queryParams.status)
    } else {
      // 都没有条件，获取全部节点
      response = await k8sApi.getClusterNodes(selectedClusterId.value)
    }

    const responseData = response.data || response
    console.log('节点列表API响应:', responseData)

    if (responseData.code === 200 || responseData.code === 0 || responseData.success) {
      const nodes = responseData.data || []
      tableData.value = nodes.map(node => ({
        id: node.name,
        nodeName: node.name,
        status: node.status || 'Unknown',
        role: node.roles === 'control-plane,master' ? 'master' : (node.roles === 'control-plane' ? 'control-plane' : (node.roles === 'master' ? 'master' : 'worker')),
        version: node.runtime?.kubeletVersion || 'Unknown',
        age: formatAge(node.createdAt),
        cpu: `${node.resources?.cpu?.allocatable || 'Unknown'}`,
        memory: formatMemory(node.resources?.memory?.allocatable || '0'),
        pods: `${node.podMetrics?.allocated || 0}/${node.podMetrics?.total || 0}`,
        schedulable: !node.scheduling?.unschedulable,
        taints: node.scheduling?.taints || [],
        labels: node.configuration?.labels || {},
        conditions: node.conditions || [],
        nodeInfo: node.runtime || {},
        allocatedResources: node.resources || {},
        internalIP: node.internalIP,
        externalIP: node.externalIP,
        configuration: node.configuration,
        runtime: node.runtime
      }))

      console.log('节点列表加载成功:', tableData.value)
    } else {
      ElMessage.error(responseData.message || '获取节点列表失败')
      tableData.value = []
    }
  } catch (error) {
    console.error('获取节点列表失败:', error)
    
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
    
    tableData.value = []
  } finally {
    loading.value = false
  }
}

const getNodeStatus = (status) => {
  // 如果直接是字符串状态，直接返回
  if (typeof status === 'string') return status
  
  // 如果是包含conditions的对象
  if (!status || !status.conditions) return 'Unknown'
  
  const readyCondition = status.conditions.find(condition => condition.type === 'Ready')
  if (readyCondition) {
    return readyCondition.status === 'True' ? 'Ready' : 'NotReady'
  }
  
  return 'Unknown'
}

const getNodeRole = (labels) => {
  if (!labels) return 'worker'
  
  if (labels['node-role.kubernetes.io/master'] || 
      labels['node-role.kubernetes.io/control-plane']) {
    return 'master'
  }
  
  return 'worker'
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

const formatMemory = (memoryStr) => {
  if (!memoryStr || memoryStr === '0') return 'Unknown'
  
  // 处理Ki单位
  if (memoryStr.endsWith('Ki')) {
    const kb = parseInt(memoryStr.replace('Ki', ''))
    const gb = (kb / 1024 / 1024).toFixed(1)
    return `${gb}GB`
  }
  
  // 处理Mi单位
  if (memoryStr.endsWith('Mi')) {
    const mb = parseInt(memoryStr.replace('Mi', ''))
    const gb = (mb / 1024).toFixed(1)
    return `${gb}GB`
  }
  
  // 处理Gi单位
  if (memoryStr.endsWith('Gi')) {
    const gb = parseInt(memoryStr.replace('Gi', ''))
    return `${gb}GB`
  }
  
  return memoryStr
}

const resetQuery = () => {
  queryParams.nodeName = ''
  queryParams.status = ''
  queryParams.role = ''
  handleQuery()
}

// 导航到监控仪表板
const navigateToMonitoring = () => {
  router.push('/k8s/monitoring')
}

// 集群选择变化处理
const handleClusterChange = () => {
  if (selectedClusterId.value) {
    handleQuery()
  } else {
    tableData.value = []
  }
}


// 查看资源分配详情
const viewResources = async (row) => {
  try {
    loading.value = true
    const response = await k8sApi.getNodeResources(selectedClusterId.value, row.nodeName)
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      const rawData = responseData.data
      
      // 适配新的API响应格式
      nodeResources.value = {
        cpu: {
          capacity: rawData.capacity?.cpu || '0',
          allocatable: rawData.allocatable?.cpu || '0',
          allocated: rawData.allocated?.cpu || '0'
        },
        memory: {
          capacity: rawData.capacity?.memory || '0',
          allocatable: rawData.allocatable?.memory || '0',
          allocated: formatBytesToKi(rawData.allocated?.memory || '0')
        },
        pods: {
          capacity: parseInt(rawData.capacity?.pods || '0'),
          allocatable: parseInt(rawData.allocatable?.pods || '0'),
          allocated: rawData.podList?.length || 0
        },
        podList: rawData.podList || []
      }
      
      currentNode.value = row
      resourcesDialogVisible.value = true
      console.log('资源详情数据:', nodeResources.value)
    } else {
      ElMessage.error(responseData.message || '获取资源详情失败')
    }
  } catch (error) {
    console.error('获取资源详情失败:', error)
    ElMessage.error('获取资源详情失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 管理污点
const manageTaints = (row) => {
  currentNode.value = row
  taintForm.operation = 'add'
  taintForm.key = ''
  taintForm.value = ''
  taintForm.effect = 'NoSchedule'
  taintForm.selectedTaint = null // 选中的要移除的污点
  taintDialogVisible.value = true
}

// 管理标签
const manageLabels = (row) => {
  currentNode.value = row
  labelForm.operation = 'add'
  labelForm.key = ''
  labelForm.value = ''
  labelForm.selectedLabel = null // 选中的要移除的标签
  labelDialogVisible.value = true
}

// 封锁/解封节点
const toggleCordon = async (row) => {
  try {
    const action = row.schedulable ? '封锁' : '解封'
    await ElMessageBox.confirm(
      `确定要${action}节点 ${row.nodeName} 吗？`,
      `${action}节点确认`,
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const response = await k8sApi.cordonNode(selectedClusterId.value, row.nodeName, {
      unschedulable: row.schedulable,
      reason: `手动${action}节点`
    })
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`节点 ${row.nodeName} ${action}成功`)
      handleQuery()
    } else {
      ElMessage.error(responseData.message || `${action}节点失败`)
    }
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消操作')
    } else {
      console.error('节点封锁操作失败:', error)
      ElMessage.error('操作失败，请检查网络连接')
    }
  }
}

// yaml编辑
const handleEditYaml = (row) => {
  openYamlDialog({
    clusterId: selectedClusterId.value,
    resourceType: 'Node',
    resourceName: row.nodeName,
    fetchYaml: k8sApi.getNodeYaml,
    updateYaml: k8sApi.updateNodeYaml,
    fallbackYaml: `# Node ${row.nodeName} YAML\napiVersion: v1\nkind: Node\nmetadata:\n  name: ${row.nodeName}`,
    loadErrorMessage: '获取节点 YAML 失败',
    saveSuccessMessage: '节点 YAML 更新成功',
    saveErrorMessage: '更新节点 YAML 失败',
    onSaved: handleQuery
  })
}

// 驱逐节点
const drainNode = (row) => {
  currentNode.value = row
  drainForm.force = false
  drainForm.deleteLocalData = false
  drainForm.ignoreDaemonSets = true
  drainForm.gracePeriodSeconds = 30
  drainDialogVisible.value = true
}

// 提交污点操作
const submitTaintOperation = async () => {
  try {
    let params = {}
    
    if (taintForm.operation === 'add') {
      if (!taintForm.key) {
        ElMessage.warning('请输入污点键名')
        return
      }
      params = {
        key: taintForm.key,
        value: taintForm.value,
        effect: taintForm.effect
      }
    } else {
      // 移除操作
      if (!taintForm.selectedTaint) {
        ElMessage.warning('请选择要移除的污点')
        return
      }
      params = {
        key: taintForm.selectedTaint.key,
        effect: taintForm.selectedTaint.effect
      }
      if (taintForm.selectedTaint.value) {
        params.value = taintForm.selectedTaint.value
      }
    }
    
    const apiCall = taintForm.operation === 'add' 
      ? k8sApi.addNodeTaint(selectedClusterId.value, currentNode.value.nodeName, params)
      : k8sApi.removeNodeTaint(selectedClusterId.value, currentNode.value.nodeName, params)
    
    const response = await apiCall
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`污点${taintForm.operation === 'add' ? '添加' : '移除'}成功`)
      taintDialogVisible.value = false
      handleQuery()
    } else {
      ElMessage.error(responseData.message || '操作失败')
    }
  } catch (error) {
    console.error('污点操作失败:', error)
    ElMessage.error('操作失败，请检查网络连接')
  }
}

// 提交标签操作
const submitLabelOperation = async () => {
  try {
    let params = {}
    
    if (labelForm.operation === 'add') {
      if (!labelForm.key) {
        ElMessage.warning('请输入标签键名')
        return
      }
      params = {
        key: labelForm.key,
        value: labelForm.value
      }
    } else {
      // 移除操作
      if (!labelForm.selectedLabel) {
        ElMessage.warning('请选择要移除的标签')
        return
      }
      params = {
        key: labelForm.selectedLabel.key
      }
    }
    
    const apiCall = labelForm.operation === 'add' 
      ? k8sApi.addNodeLabel(selectedClusterId.value, currentNode.value.nodeName, params)
      : k8sApi.removeNodeLabel(selectedClusterId.value, currentNode.value.nodeName, params)
    
    const response = await apiCall
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`标签${labelForm.operation === 'add' ? '添加' : '移除'}成功`)
      labelDialogVisible.value = false
      handleQuery()
    } else {
      ElMessage.error(responseData.message || '操作失败')
    }
  } catch (error) {
    console.error('标签操作失败:', error)
    ElMessage.error('操作失败，请检查网络连接')
  }
}

// 提交驱逐操作
const submitDrainOperation = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要驱逐节点 ${currentNode.value.nodeName} 上的所有Pod吗？这将会影响运行中的服务！`,
      '驱逐节点确认',
      {
        confirmButtonText: '确定驱逐',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const response = await k8sApi.drainNode(selectedClusterId.value, currentNode.value.nodeName, drainForm)
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`节点 ${currentNode.value.nodeName} 驱逐成功`)
      drainDialogVisible.value = false
      handleQuery()
    } else {
      ElMessage.error(responseData.message || '驱逐失败')
    }
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消驱逐操作')
    } else {
      console.error('节点驱逐失败:', error)
      ElMessage.error('驱逐失败，请检查网络连接')
    }
  }
}


// 获取当前节点的污点选项（用于移除操作）
const currentNodeTaints = computed(() => {
  if (!currentNode.value || !currentNode.value.taints) return []
  
  return currentNode.value.taints.map(taint => ({
    label: `${taint.key}=${taint.value || ''}:${taint.effect}`,
    value: taint,
    taint: taint
  }))
})

// 获取当前节点的标签选项（用于移除操作）
const currentNodeLabels = computed(() => {
  if (!currentNode.value || !currentNode.value.labels) return []
  
  // 过滤掉系统标签，只显示用户可以操作的标签
  const systemLabelPrefixes = [
    'kubernetes.io/',
    'beta.kubernetes.io/',
    'node-role.kubernetes.io/',
    'node.kubernetes.io/'
  ]
  
  return Object.entries(currentNode.value.labels)
    .filter(([key]) => !systemLabelPrefixes.some(prefix => key.startsWith(prefix)))
    .map(([key, value]) => ({
      label: `${key}=${value}`,
      value: { key, value },
      labelData: { key, value }
    }))
})

// 集群状态相关方法
const getClusterStatusTag = (status) => {
  const tagMap = {
    1: 'info',      // 创建中
    2: 'success',   // 运行中
    3: 'warning',   // 已停止
    4: 'danger',    // 异常
    5: 'info'       // 已删除
  }
  return tagMap[status] || 'info'
}

const getClusterStatusText = (status) => {
  const textMap = {
    1: '创建中',
    2: '运行中',
    3: '已停止',
    4: '异常',
    5: '已删除'
  }
  return textMap[status] || '未知'
}

// 复制内容到剪贴板
const copyToClipboard = async (text, successMessage = '已复制到剪贴板') => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(successMessage)
  } catch (error) {
    console.error('复制失败:', error)
    // 降级处理
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

// 字节转换为Ki单位
const formatBytesToKi = (bytes) => {
  if (!bytes || bytes === '0') return '0Ki'
  const bytesNum = parseInt(bytes)
  if (isNaN(bytesNum)) return '0Ki'
  const ki = Math.round(bytesNum / 1024)
  return `${ki}Ki`
}

// 计算CPU使用百分比
const getCpuPercentage = (cpu) => {
  if (!cpu || !cpu.allocatable || !cpu.allocated) return 0
  
  // 统一转换为millicores进行计算
  const parseCpuToMillicores = (cpuStr) => {
    if (!cpuStr) return 0
    if (cpuStr.endsWith('m')) {
      return parseFloat(cpuStr.replace('m', ''))
    }
    // 如果是纯数字，表示cores，需要转换为millicores
    return parseFloat(cpuStr) * 1000
  }
  
  const allocatable = parseCpuToMillicores(cpu.allocatable)
  const allocated = parseCpuToMillicores(cpu.allocated)
  
  console.log('CPU计算:', {
    allocatable: cpu.allocatable,
    allocated: cpu.allocated,
    allocatableMillicores: allocatable,
    allocatedMillicores: allocated,
    percentage: allocatable > 0 ? Math.round((allocated / allocatable) * 100) : 0
  })
  
  if (allocatable === 0) return 0
  return Math.round((allocated / allocatable) * 100)
}

// 计算内存使用百分比
const getMemoryPercentage = (memory) => {
  if (!memory || !memory.allocatable || !memory.allocated) return 0
  
  const getMemoryKi = (memStr) => {
    if (!memStr) return 0
    if (memStr.endsWith('Ki')) return parseInt(memStr.replace('Ki', ''))
    if (memStr.endsWith('Mi')) return parseInt(memStr.replace('Mi', '')) * 1024
    if (memStr.endsWith('Gi')) return parseInt(memStr.replace('Gi', '')) * 1024 * 1024
    return parseInt(memStr) || 0
  }
  
  const allocatable = getMemoryKi(memory.allocatable)
  const allocated = getMemoryKi(memory.allocated)
  
  if (allocatable === 0) return 0
  return Math.round((allocated / allocatable) * 100)
}

// 计算Pod使用百分比 - 适配新的资源数据结构
const getPodPercentage = (pods) => {
  if (!pods || !pods.allocatable || typeof pods.allocated === 'undefined') return 0
  
  const total = parseInt(pods.allocatable) || 0
  const allocated = parseInt(pods.allocated) || 0
  
  if (total === 0) return 0
  return Math.round((allocated / total) * 100)
}


// 查看节点标签
const viewNodeLabels = (row) => {
  currentNode.value = row
  nodeLabelsDialogVisible.value = true
}

// 获取用户自定义标签
const getUserLabels = (labels) => {
  if (!labels) return {}
  
  const systemLabelPrefixes = [
    'kubernetes.io/',
    'beta.kubernetes.io/',
    'node-role.kubernetes.io/',
    'node.kubernetes.io/'
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
    'node.kubernetes.io/'
  ]
  
  const systemLabels = {}
  Object.entries(labels).forEach(([key, value]) => {
    if (systemLabelPrefixes.some(prefix => key.startsWith(prefix))) {
      systemLabels[key] = value
    }
  })
  
  return systemLabels
}

onMounted(async () => {
  // 先获取集群列表，然后根据选择的集群获取节点列表
  await fetchClusterList()
  // fetchClusterList会自动选择第一个可用集群并触发节点查询
  if (selectedClusterId.value) {
    handleQuery()
  }
})
</script>

<template>
  <div class="k8s-nodes-management">
    <el-card shadow="hover" class="nodes-card">
      <template #header>
        <div class="card-header">
          <span class="title">K8s 节点管理</span>
          <div class="header-actions">
            <el-select 
              v-model="selectedClusterId" 
              placeholder="请选择集群"
              size="small" 
              style="width: 250px" 
              @change="handleClusterChange"
              :loading="!clusterList.length"
            >
              <el-option
                v-for="cluster in clusterList"
                :key="cluster.id"
                :label="`${cluster.name} ${getClusterStatusText(cluster.status)}`"
                :value="cluster.id"
                :disabled="cluster.status !== 2"
              >
                <div class="cluster-option">
                  <span class="cluster-name">{{ cluster.name }}</span>
                  <el-tag 
                    :type="getClusterStatusTag(cluster.status)" 
                    size="small"
                    class="cluster-status-tag"
                  >
                    {{ getClusterStatusText(cluster.status) }}
                  </el-tag>
                </div>
              </el-option>
            </el-select>
          </div>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :inline="true" :model="queryParams" class="search-form">
          <el-form-item label="节点名称">
            <el-input
              v-model="queryParams.nodeName"
              placeholder="请输入节点名称或IP"
              clearable
              size="small"
              style="width: 200px"
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item label="节点角色">
            <el-select
              v-model="queryParams.role"
              placeholder="请选择角色"
              clearable
              size="small"
              style="width: 150px"
            >
              <el-option label="全部" value="" />
              <el-option label="Master" value="master" />
              <el-option label="Worker" value="worker" />
              <el-option label="Control-plane" value="control-plane" />
            </el-select>
          </el-form-item>
          <el-form-item label="节点状态">
            <el-select
              v-model="queryParams.status"
              placeholder="请选择状态"
              clearable
              size="small"
              style="width: 150px"
            >
              <el-option label="全部" value="" />
              <el-option label="就绪" value="Ready" />
              <el-option label="未就绪" value="NotReady" />
              <el-option label="未知" value="Unknown" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :icon="Search" size="small" @click="handleQuery">
              搜索
            </el-button>
            <el-button type="warning" :icon="Refresh" size="small" @click="resetQuery">
              重置
            </el-button>
            <el-button :icon="Monitor" type="success" size="small" @click="navigateToMonitoring">
              监控仪表板
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 节点列表表格 -->
      <NodesTable
        :table-data="tableData"
        :loading="loading"
        :selected-cluster-id="selectedClusterId"
        @view-resources="viewResources"
        @manage-taints="manageTaints"
        @manage-labels="manageLabels"
        @view-node-labels="viewNodeLabels"
        @toggle-cordon="toggleCordon"
        @drain-node="drainNode"
        @edit-yaml="handleEditYaml"
      />
    </el-card>


    <!-- 资源详情对话框 -->
    <el-dialog
      v-model="resourcesDialogVisible"
      :title="`资源详情 - ${currentNode.nodeName || ''}`"
      width="900px"
      class="resources-dialog"
    >
      <div class="resources-content" v-if="nodeResources">
        <!-- 资源概览 -->
        <div class="resources-overview">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-card class="resource-card">
                <div class="resource-header">
                  <span class="resource-title">CPU</span>
                  <span class="resource-unit">cores</span>
                </div>
                <div class="resource-progress">
                  <el-progress
                    :percentage="getCpuPercentage(nodeResources.cpu)"
                    :color="getProgressColor(getCpuPercentage(nodeResources.cpu))"
                  />
                </div>
                <div class="resource-details">
                  <span>已分配: {{ nodeResources.cpu?.allocated || '0' }}</span>
                  <span>总容量: {{ nodeResources.cpu?.allocatable || '0' }}</span>
                </div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card class="resource-card">
                <div class="resource-header">
                  <span class="resource-title">内存</span>
                  <span class="resource-unit">GB</span>
                </div>
                <div class="resource-progress">
                  <el-progress
                    :percentage="getMemoryPercentage(nodeResources.memory)"
                    :color="getProgressColor(getMemoryPercentage(nodeResources.memory))"
                  />
                </div>
                <div class="resource-details">
                  <span>已分配: {{ formatMemory(nodeResources.memory?.allocated || '0') }}</span>
                  <span>总容量: {{ formatMemory(nodeResources.memory?.allocatable || '0') }}</span>
                </div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card class="resource-card">
                <div class="resource-header">
                  <span class="resource-title">Pod</span>
                  <span class="resource-unit">个</span>
                </div>
                <div class="resource-progress">
                  <el-progress
                    :percentage="getPodPercentage(nodeResources.pods)"
                    :color="getProgressColor(getPodPercentage(nodeResources.pods))"
                  />
                </div>
                <div class="resource-details">
                  <span>运行中: {{ nodeResources.pods?.allocated || 0 }}</span>
                  <span>最大容量: {{ nodeResources.pods?.allocatable || 0 }}</span>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>
        
        <!-- Pod列表 -->
        <div class="pods-section" v-if="nodeResources.podList && nodeResources.podList.length > 0">
          <h4>Pod列表</h4>
          <el-table :data="nodeResources.podList" size="small" style="width: 100%">
            <el-table-column prop="name" label="Pod名称" min-width="200" />
            <el-table-column prop="namespace" label="命名空间" min-width="120" />
            <el-table-column label="CPU请求" min-width="100">
              <template #default="{ row }">
                <span>{{ row.requests?.cpu || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="内存请求" min-width="100">
              <template #default="{ row }">
                <span>{{ row.requests?.memory || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="CPU限制" min-width="100">
              <template #default="{ row }">
                <span>{{ row.limits?.cpu || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="内存限制" min-width="100">
              <template #default="{ row }">
                <span>{{ row.limits?.memory || '-' }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-dialog>

    <!-- 污点管理对话框 -->
    <el-dialog
      v-model="taintDialogVisible"
      :title="`管理污点 - ${currentNode.nodeName || ''}`"
      width="500px"
      class="taint-dialog"
    >
      <el-form :model="taintForm" label-width="100px">
        <el-form-item label="操作类型" required>
          <el-radio-group v-model="taintForm.operation">
            <el-radio label="add">添加污点</el-radio>
            <el-radio label="remove">移除污点</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <!-- 添加污点的表单 -->
        <template v-if="taintForm.operation === 'add'">
          <el-form-item label="污点键名" required>
            <el-input
              v-model="taintForm.key"
              placeholder="请输入污点键名，如：maintenance"
            />
          </el-form-item>
          <el-form-item label="污点值">
            <el-input
              v-model="taintForm.value"
              placeholder="请输入污点值，如：scheduled（可选）"
            />
          </el-form-item>
          <el-form-item label="污点效果" required>
            <el-select v-model="taintForm.effect" style="width: 100%" placeholder="请选择污点效果">
              <el-option
                v-for="option in effectOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              >
                <div class="effect-option">
                  <div class="effect-label">{{ option.label }}</div>
                  <div class="effect-description">{{ option.description }}</div>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </template>
        
        <!-- 移除污点的表单 -->
        <template v-else>
          <el-form-item label="选择污点" required v-if="currentNodeTaints.length > 0">
            <el-select 
              v-model="taintForm.selectedTaint" 
              style="width: 100%" 
              placeholder="请选择要移除的污点"
              value-key="key"
            >
              <el-option
                v-for="option in currentNodeTaints"
                :key="option.label"
                :label="option.label"
                :value="option.taint"
              >
                <div class="taint-option">
                  <el-tag type="warning" size="small">{{ option.label }}</el-tag>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-else>
            <el-alert
              title="该节点没有污点"
              type="info"
              description="当前节点没有污点，无需移除操作"
              show-icon
              :closable="false"
            />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="taintDialogVisible = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="submitTaintOperation"
            :disabled="taintForm.operation === 'remove' && currentNodeTaints.length === 0"
          >
            {{ taintForm.operation === 'add' ? '添加污点' : '移除污点' }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 标签管理对话框 -->
    <el-dialog
      v-model="labelDialogVisible"
      :title="`管理标签 - ${currentNode.nodeName || ''}`"
      width="500px"
      class="label-dialog"
    >
      <el-form :model="labelForm" label-width="100px">
        <el-form-item label="操作类型" required>
          <el-radio-group v-model="labelForm.operation">
            <el-radio label="add">添加标签</el-radio>
            <el-radio label="remove">移除标签</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <!-- 添加标签的表单 -->
        <template v-if="labelForm.operation === 'add'">
          <el-form-item label="标签键名" required>
            <el-input
              v-model="labelForm.key"
              placeholder="请输入标签键名，如：environment"
            />
          </el-form-item>
          <el-form-item label="标签值" required>
            <el-input
              v-model="labelForm.value"
              placeholder="请输入标签值，如：production"
            />
          </el-form-item>
        </template>
        
        <!-- 移除标签的表单 -->
        <template v-else>
          <el-form-item label="选择标签" required v-if="currentNodeLabels.length > 0">
            <el-select 
              v-model="labelForm.selectedLabel" 
              style="width: 100%" 
              placeholder="请选择要移除的标签"
              value-key="key"
            >
              <el-option
                v-for="option in currentNodeLabels"
                :key="option.label"
                :label="option.label"
                :value="option.labelData"
              >
                <div class="label-option">
                  <el-tag type="info" size="small">{{ option.label }}</el-tag>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-else>
            <el-alert
              title="该节点没有可移除的标签"
              type="info"
              description="当前节点没有用户自定义标签，只有系统标签无法移除"
              show-icon
              :closable="false"
            />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="labelDialogVisible = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="submitLabelOperation"
            :disabled="labelForm.operation === 'remove' && currentNodeLabels.length === 0"
          >
            {{ labelForm.operation === 'add' ? '添加标签' : '移除标签' }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 驱逐节点对话框 -->
    <el-dialog
      v-model="drainDialogVisible"
      :title="`驱逐节点 - ${currentNode.nodeName || ''}`"
      width="500px"
      class="drain-dialog"
    >
      <el-form :model="drainForm" label-width="120px">
        <el-form-item label="强制驱逐">
          <el-switch
            v-model="drainForm.force"
            active-text="是"
            inactive-text="否"
          />
          <div class="form-tip">强制驱逐会删除不受控制器管理的Pod</div>
        </el-form-item>
        <el-form-item label="删除本地数据">
          <el-switch
            v-model="drainForm.deleteLocalData"
            active-text="是"
            inactive-text="否"
          />
          <div class="form-tip">删除使用emptyDir的Pod</div>
        </el-form-item>
        <el-form-item label="忽略DaemonSet">
          <el-switch
            v-model="drainForm.ignoreDaemonSets"
            active-text="是"
            inactive-text="否"
          />
          <div class="form-tip">忽略DaemonSet管理的Pod</div>
        </el-form-item>
        <el-form-item label="优雅终止时间">
          <el-input-number
            v-model="drainForm.gracePeriodSeconds"
            :min="0"
            :max="300"
            style="width: 100%"
          />
          <div class="form-tip">Pod优雅终止等待时间（秒）</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="drainDialogVisible = false">取消</el-button>
          <el-button type="danger" @click="submitDrainOperation">驱逐</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 节点标签查看对话框 -->
    <el-dialog
      v-model="nodeLabelsDialogVisible"
      :title="`节点标签 - ${currentNode.nodeName || ''}`"
      width="800px"
      class="node-labels-view-dialog"
    >
      <div class="labels-view-content" v-if="currentNode.labels">
        <!-- 用户自定义标签 -->
        <div class="labels-section" v-if="Object.keys(getUserLabels(currentNode.labels)).length > 0">
          <h4 class="section-title">用户标签</h4>
          <div class="labels-table-wrapper">
            <el-table
              :data="Object.entries(getUserLabels(currentNode.labels)).map(([k, v]) => ({ key: k, value: v }))"
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
        <div class="labels-section" v-if="Object.keys(getSystemLabels(currentNode.labels)).length > 0">
          <h4 class="section-title">系统标签</h4>
          <div class="labels-table-wrapper">
            <el-table
              :data="Object.entries(getSystemLabels(currentNode.labels)).map(([k, v]) => ({ key: k, value: v }))"
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
        <div v-if="!currentNode.labels || Object.keys(currentNode.labels).length === 0" class="no-labels">
          <el-empty description="该节点没有标签" :image-size="60" />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
// 辅助函数
const getProgressColor = (percentage) => {
  if (percentage < 60) return '#67c23a'
  if (percentage < 80) return '#e6a23c'
  return '#f56c6c'
}

// 导出辅助函数以便在模板中使用
export { getProgressColor }
</script>

<style scoped>
.k8s-nodes-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.nodes-card {
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


/* 对话框样式 */
.resources-dialog :deep(.el-dialog),
.taint-dialog :deep(.el-dialog),
.label-dialog :deep(.el-dialog),
.drain-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.resources-dialog :deep(.el-dialog__header),
.taint-dialog :deep(.el-dialog__header),
.label-dialog :deep(.el-dialog__header),
.drain-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.resources-dialog :deep(.el-dialog__title),
.taint-dialog :deep(.el-dialog__title),
.label-dialog :deep(.el-dialog__title),
.drain-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.resources-dialog :deep(.el-dialog__body),
.taint-dialog :deep(.el-dialog__body),
.label-dialog :deep(.el-dialog__body),
.drain-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

.node-detail-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.taints-section,
.labels-section {
  margin-top: 20px;
}

.taints-section h4,
.labels-section h4,
.section-title {
  margin: 0 0 12px 0;
  color: #6b9bd1;
  font-weight: 600;
}

.taints-list,
.labels-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.taint-tag,
.label-tag {
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 4px;
}

.taint-tag:hover,
.label-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* 标签表格样式 */
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

/* 资源详情样式 */
.resources-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.resources-overview {
  margin-bottom: 20px;
}

.resource-card {
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.2);
}

.resource-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.resource-title {
  font-weight: 600;
  color: #6b9bd1;
  font-size: 16px;
}

.resource-unit {
  font-size: 12px;
  color: #909399;
  background: rgba(103, 126, 234, 0.1);
  padding: 2px 8px;
  border-radius: 12px;
}

.resource-progress {
  margin-bottom: 12px;
}

.resource-details {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #606266;
}

.pods-section {
  margin-top: 20px;
}

.pods-section h4 {
  margin: 0 0 12px 0;
  color: #6b9bd1;
  font-weight: 600;
}

/* 表单样式 */
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

/* 通用元素样式 */
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

.header-actions .el-select {
  min-width: 250px;
}

.header-actions .el-select .el-input__wrapper {
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(103, 126, 234, 0.3);
}

.header-actions .el-select .el-input__wrapper.is-focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
}

/* 污点效果选项样式 */
.effect-option {
  padding: 4px 0;
}

.effect-label {
  font-weight: 500;
  color: #6b9bd1;
  font-size: 14px;
  margin-bottom: 2px;
}

.effect-description {
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
  margin-left: 8px;
}

/* 污点选项样式 */
.taint-option {
  padding: 4px 0;
  display: flex;
  align-items: center;
}

/* 标签选项样式 */
.label-option {
  padding: 4px 0;
  display: flex;
  align-items: center;
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

.label-icon {
  color: #409eff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.label-icon:hover {
  color: #66b1ff;
  transform: scale(1.1);
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

/* 标签查看对话框样式 */
.node-labels-view-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.node-labels-view-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.node-labels-view-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.node-labels-view-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

.labels-view-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.labels-section h4 {
  margin: 0 0 16px 0;
  color: #6b9bd1;
  font-weight: 600;
  font-size: 16px;
  padding-bottom: 8px;
  border-bottom: 2px solid #e4e7ed;
}

.labels-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.label-tag {
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  font-size: 13px;
  line-height: 1.4;
  border-radius: 8px;
  max-width: 300px;
  word-break: break-all;
}

.label-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.label-tag.system-label {
  opacity: 0.8;
}

.tag-icon {
  font-size: 12px;
  flex-shrink: 0;
}

.no-labels {
  text-align: center;
  padding: 40px 20px;
  color: #909399;
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
    min-width: 200px;
  }
}

@media (max-width: 768px) {
  .k8s-nodes-management {
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
  
}
</style>
