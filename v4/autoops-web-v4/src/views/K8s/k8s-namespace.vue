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
  Monitor
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import {
  getStoredClusterId,
  saveStoredClusterId,
  resolvePreferredClusterId
} from '@/utils/k8sClusterSelection'
import NamespacesTable from './namespaces/NamespacesTable.vue'
import CreateNamespaceDialog from './namespaces/CreateNamespaceDialog.vue'
import NamespaceLabelsDialog from './namespaces/NamespaceLabelsDialog.vue'
import { openYamlDialog } from './services/yamlDialog'

const router = useRouter()

const loading = ref(false)
const queryParams = reactive({
  namespaceName: '',
  status: ''
})

const statusOptions = [
  { label: '全部', value: '' },
  { label: '活跃', value: 'Active' },
  { label: '终止中', value: 'Terminating' },
  { label: '未知', value: 'Unknown' }
]

const tableData = ref([])
const selectedClusterId = ref(getStoredClusterId())
const clusterList = ref([])

// 分页
const currentPage = ref(1)
const pageSize = ref(10)
const total = computed(() => tableData.value.length)
const pagedTableData = computed(() =>
  tableData.value.slice((currentPage.value - 1) * pageSize.value, currentPage.value * pageSize.value)
)
const handlePageChange = (page) => { currentPage.value = page }
const handleSizeChange = (size) => { pageSize.value = size; currentPage.value = 1 }

// 对话框状态
const createNamespaceDialogVisible = ref(false)
const namespaceDetailDialogVisible = ref(false)
const resourceQuotaDialogVisible = ref(false)
const limitRangeDialogVisible = ref(false)
const createQuotaDialogVisible = ref(false)
const createLimitRangeDialogVisible = ref(false)
const labelsDialogVisible = ref(false)

// 当前操作的命名空间
const currentNamespace = ref({})

// 创建对话框引用
const createDialogRef = ref(null)

// 资源配额数据 - 确保初始值为空数组
const resourceQuotas = ref([])
const limitRanges = ref([])

// ResourceQuota表单
const quotaForm = reactive({
  name: '',
  cpuLimit: '',
  memoryLimit: '',
  storageLimit: '',
  podLimit: '',
  serviceLimit: ''
})

// LimitRange表单
const limitRangeForm = reactive({
  name: '',
  containerCpuMin: '',
  containerCpuMax: '',
  containerMemoryMin: '',
  containerMemoryMax: '',
  containerCpuDefault: '',
  containerMemoryDefault: ''
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

const handleQuery = async () => {
  try {
    if (!selectedClusterId.value) {
      ElMessage.warning('请选择一个集群')
      return
    }

    loading.value = true
    currentPage.value = 1

    let response
    if (queryParams.namespaceName) {
      // 有名称时走模糊查询 API
      response = await k8sApi.searchNamespaces(selectedClusterId.value, queryParams.namespaceName)
    } else {
      // 无名称时走完整列表 API（可带状态过滤）
      const params = {}
      if (queryParams.status) params.status = queryParams.status
      response = await k8sApi.getNamespaces(selectedClusterId.value, Object.keys(params).length > 0 ? params : undefined)
    }

    const responseData = response.data || response
    console.log('命名空间列表API响应:', responseData)

    if (responseData.code === 200 || responseData.success) {
      const namespaces = responseData.data?.namespaces || responseData.data || []
      const namespacesArray = Array.isArray(namespaces) ? namespaces : []
      tableData.value = namespacesArray.map(ns => ({
        name: ns.name,
        status: ns.status || 'Active',
        age: formatAge(ns.createdAt),
        createdAt: ns.createdAt,
        labels: ns.labels || {},
        annotations: ns.annotations || {},
        resourceQuotas: ns.resourceQuotas || [],
        limitRanges: ns.limitRanges || [],
        resourceCount: {
          pods: ns.resourceCount?.podCount || 0,
          services: ns.resourceCount?.serviceCount || 0,
          secrets: ns.resourceCount?.secretCount || 0,
          configMaps: ns.resourceCount?.configMapCount || 0
        }
      }))

      console.log('命名空间列表加载成功:', tableData.value)
    } else {
      ElMessage.error(responseData.message || '获取命名空间列表失败')
      tableData.value = []
    }
  } catch (error) {
    console.error('获取命名空间列表失败:', error)

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

const resetQuery = () => {
  queryParams.namespaceName = ''
  queryParams.status = ''
  handleQuery()
}

// 集群选择变化处理
const handleClusterChange = () => {
  if (selectedClusterId.value) {
    saveStoredClusterId(selectedClusterId.value)
    handleQuery()
  } else {
    saveStoredClusterId('')
    tableData.value = []
  }
}

// 创建命名空间
const showCreateNamespaceDialog = () => {
  createNamespaceDialogVisible.value = true
}

// 处理创建命名空间成功
const handleCreateNamespace = async (data) => {
  try {
    if (createDialogRef.value) {
      createDialogRef.value.setSubmitting(true)
    }

    const response = await k8sApi.createNamespace(selectedClusterId.value, data)
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('命名空间创建成功')
      createNamespaceDialogVisible.value = false
      handleQuery()
    } else {
      ElMessage.error(responseData.message || '创建命名空间失败')
    }
  } catch (error) {
    console.error('创建命名空间失败:', error)
    ElMessage.error('创建命名空间失败，请检查网络连接')
  } finally {
    if (createDialogRef.value) {
      createDialogRef.value.setSubmitting(false)
    }
  }
}

// 查看命名空间详情
const viewNamespaceDetail = async (row) => {
  try {
    loading.value = true
    const response = await k8sApi.getNamespaceDetail(selectedClusterId.value, row.name)
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      // 如果返回的data中有单个命名空间对象，直接使用；否则使用整个data
      currentNamespace.value = responseData.data?.name ? responseData.data : row
      namespaceDetailDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取命名空间详情失败')
    }
  } catch (error) {
    console.error('获取命名空间详情失败:', error)
    ElMessage.error('获取命名空间详情失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 管理ResourceQuota
const manageResourceQuotas = async (row) => {
  try {
    currentNamespace.value = row
    
    const response = await k8sApi.getResourceQuotas(selectedClusterId.value, row.name)
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      // 适配后端返回的数据结构
      const quotaData = responseData.data?.resourceQuotas || responseData.data || []
      resourceQuotas.value = Array.isArray(quotaData) ? quotaData : []
      resourceQuotaDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取ResourceQuota失败')
    }
  } catch (error) {
    console.error('获取ResourceQuota失败:', error)
    ElMessage.error('获取ResourceQuota失败，请检查网络连接')
  }
}

// 管理LimitRange
const manageLimitRanges = async (row) => {
  try {
    currentNamespace.value = row
    
    const response = await k8sApi.getLimitRanges(selectedClusterId.value, row.name)
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      // 适配后端返回的数据结构：后端返回的是 limits 直接在对象下，需要转换为前端期望的 spec.limits 结构
      const limitData = responseData.data?.limitRanges || responseData.data || []
      const mappedData = Array.isArray(limitData) ? limitData.map(item => ({
        name: item.name,
        createdAt: item.createdAt,
        spec: {
          limits: item.limits || []  // 将后端的 limits 映射到 spec.limits
        }
      })) : []
      
      limitRanges.value = mappedData
      limitRangeDialogVisible.value = true
      
      console.log('LimitRange数据加载成功:', limitRanges.value)
    } else {
      ElMessage.error(responseData.message || '获取LimitRange失败')
    }
  } catch (error) {
    console.error('获取LimitRange失败:', error)
    ElMessage.error('获取LimitRange失败，请检查网络连接')
  }
}

// 删除命名空间
const deleteNamespace = async (row) => {
  try {
    // 额外检查系统命名空间
    if (isSystemNamespace(row.name)) {
      ElMessage.warning('系统命名空间不可删除')
      return
    }
    
    await ElMessageBox.confirm(
      `确定要删除命名空间 ${row.name} 吗？删除后将无法恢复，且会删除该命名空间下的所有资源！`,
      '删除命名空间确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'error',
      }
    )
    
    const response = await k8sApi.deleteNamespace(selectedClusterId.value, row.name)
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`命名空间 ${row.name} 删除成功`)
      handleQuery()
    } else {
      ElMessage.error(responseData.message || '删除命名空间失败')
    }
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消删除操作')
    } else {
      console.error('删除命名空间失败:', error)
      ElMessage.error('删除命名空间失败，请检查网络连接')
    }
  }
}

const getStatusTag = (status) => {
  const tagMap = {
    'Active': 'success',
    'Terminating': 'warning',
    'Unknown': 'info'
  }
  return tagMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    'Active': '活跃',
    'Terminating': '终止中',
    'Unknown': '未知'
  }
  return textMap[status] || '未知'
}

// 集群状态相关方法
const getClusterStatusTag = (status) => {
  const tagMap = {
    1: 'info',
    2: 'success',
    3: 'warning',
    4: 'danger',
    5: 'info'
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

// 获取标签数量（过滤系统标签）
const getLabelCount = (labels) => {
  if (!labels) return 0
  
  const systemLabelPrefixes = [
    'kubernetes.io/',
    'beta.kubernetes.io/',
    'app.kubernetes.io/'
  ]
  
  return Object.keys(labels).filter(key => 
    !systemLabelPrefixes.some(prefix => key.startsWith(prefix))
  ).length
}

// 获取用户自定义标签
const getUserLabels = (labels) => {
  if (!labels) return {}
  
  const systemLabelPrefixes = [
    'kubernetes.io/',
    'beta.kubernetes.io/',
    'app.kubernetes.io/'
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
    'app.kubernetes.io/'
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

// 查看命名空间标签
const viewNamespaceLabels = (row) => {
  viewNamespaceDetail(row)
}

// yaml编辑
const handleEditYaml = (row) => {
  currentNamespace.value = row
  openYamlDialog({
    clusterId: selectedClusterId.value,
    resourceType: 'Namespace',
    resourceName: row.name,
    fetchYaml: k8sApi.getNamespaceYaml,
    updateYaml: k8sApi.updateNamespaceYaml,
    fallbackYaml: `# Namespace ${row.name} YAML\napiVersion: v1\nkind: Namespace\nmetadata:\n  name: ${row.name}`,
    loadErrorMessage: '获取命名空间 YAML 失败',
    saveSuccessMessage: '命名空间 YAML 更新成功',
    saveErrorMessage: '更新命名空间 YAML 失败',
    onSaved: handleQuery
  })
}

// 修改标签
const handleManageLabels = (row) => {
  currentNamespace.value = row
  labelsDialogVisible.value = true
}

// 判断是否为系统命名空间
const isSystemNamespace = (namespaceName) => {
  const systemNamespaces = [
    'default',
    'kube-system', 
    'kube-public',
    'kube-node-lease'
  ]
  return systemNamespaces.includes(namespaceName)
}

// 获取注解数量
const getAnnotationCount = (annotations) => {
  if (!annotations || typeof annotations !== 'object') return 0
  return Object.keys(annotations).length
}

// 显示创建ResourceQuota对话框
const showCreateQuotaDialog = () => {
  quotaForm.name = ''
  quotaForm.cpuLimit = ''
  quotaForm.memoryLimit = ''
  quotaForm.storageLimit = ''
  quotaForm.podLimit = ''
  quotaForm.serviceLimit = ''
  quotaForm.isEditing = false
  quotaForm.originalName = ''
  createQuotaDialogVisible.value = true
}

// 显示创建LimitRange对话框
const showCreateLimitRangeDialog = () => {
  limitRangeForm.name = ''
  limitRangeForm.containerCpuMin = ''
  limitRangeForm.containerCpuMax = ''
  limitRangeForm.containerMemoryMin = ''
  limitRangeForm.containerMemoryMax = ''
  limitRangeForm.containerCpuDefault = ''
  limitRangeForm.containerMemoryDefault = ''
  limitRangeForm.isEditing = false
  limitRangeForm.originalName = ''
  createLimitRangeDialogVisible.value = true
}

// 提交创建ResourceQuota
const submitCreateQuota = async () => {
  try {
    if (!quotaForm.name) {
      ElMessage.warning('请输入配额名称')
      return
    }
    
    const data = {
      name: quotaForm.name,
      hard: {}
    }
    
    // 自动添加单位
    if (quotaForm.cpuLimit) {
      // CPU: 如果是纯数字，自动添加单位 (cores)，如果已有单位则保持
      const cpuValue = quotaForm.cpuLimit.trim()
      data.hard['limits.cpu'] = /^\d+$/.test(cpuValue) ? `${cpuValue}` : cpuValue
    }
    
    if (quotaForm.memoryLimit) {
      // 内存: 如果是纯数字，自动添加Gi单位
      const memoryValue = quotaForm.memoryLimit.trim()
      data.hard['limits.memory'] = /^\d+$/.test(memoryValue) ? `${memoryValue}Gi` : memoryValue
    }
    
    if (quotaForm.storageLimit) {
      // 存储: 如果是纯数字，自动添加Gi单位
      const storageValue = quotaForm.storageLimit.trim()
      data.hard['requests.storage'] = /^\d+$/.test(storageValue) ? `${storageValue}Gi` : storageValue
    }
    
    // 数量类型不需要单位
    if (quotaForm.podLimit) data.hard['count/pods'] = quotaForm.podLimit
    if (quotaForm.serviceLimit) data.hard['count/services'] = quotaForm.serviceLimit
    
    let response, responseData
    
    if (quotaForm.isEditing) {
      // 编辑模式：更新现有配额
      response = await k8sApi.updateResourceQuota(selectedClusterId.value, currentNamespace.value.name, quotaForm.originalName, data)
    } else {
      // 创建模式：创建新配额
      response = await k8sApi.createResourceQuota(selectedClusterId.value, currentNamespace.value.name, data)
    }
    
    responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(quotaForm.isEditing ? 'ResourceQuota更新成功' : 'ResourceQuota创建成功')
      createQuotaDialogVisible.value = false
      // 重新获取ResourceQuota列表
      manageResourceQuotas(currentNamespace.value)
    } else {
      ElMessage.error(responseData.message || (quotaForm.isEditing ? '更新ResourceQuota失败' : '创建ResourceQuota失败'))
    }
  } catch (error) {
    console.error((quotaForm.isEditing ? '更新' : '创建') + 'ResourceQuota失败:', error)
    ElMessage.error((quotaForm.isEditing ? '更新' : '创建') + 'ResourceQuota失败，请检查网络连接')
  }
}

// 提交创建LimitRange
const submitCreateLimitRange = async () => {
  try {
    if (!limitRangeForm.name) {
      ElMessage.warning('请输入限制名称')
      return
    }
    
    const data = {
      name: limitRangeForm.name,
      spec: {
        limits: [
          {
            type: 'Container',
            min: {},
            max: {},
            default: {},
            defaultRequest: {}
          }
        ]
      }
    }
    
    const limit = data.spec.limits[0]
    
    // 自动添加单位 - CPU
    if (limitRangeForm.containerCpuMin) {
      const value = limitRangeForm.containerCpuMin.trim()
      limit.min.cpu = /^\d+$/.test(value) ? `${value}m` : value // 纯数字默认为毫核心
    }
    if (limitRangeForm.containerCpuMax) {
      const value = limitRangeForm.containerCpuMax.trim()
      limit.max.cpu = /^\d+$/.test(value) ? `${value}` : value // 纯数字默认为核心
    }
    if (limitRangeForm.containerCpuDefault) {
      const value = limitRangeForm.containerCpuDefault.trim()
      limit.default.cpu = /^\d+$/.test(value) ? `${value}m` : value // 纯数字默认为毫核心
    }
    
    // 自动添加单位 - Memory
    if (limitRangeForm.containerMemoryMin) {
      const value = limitRangeForm.containerMemoryMin.trim()
      limit.min.memory = /^\d+$/.test(value) ? `${value}Mi` : value // 纯数字默认为Mi
    }
    if (limitRangeForm.containerMemoryMax) {
      const value = limitRangeForm.containerMemoryMax.trim()
      limit.max.memory = /^\d+$/.test(value) ? `${value}Gi` : value // 纯数字默认为Gi
    }
    if (limitRangeForm.containerMemoryDefault) {
      const value = limitRangeForm.containerMemoryDefault.trim()
      limit.default.memory = /^\d+$/.test(value) ? `${value}Mi` : value // 纯数字默认为Mi
    }
    
    let response, responseData
    
    if (limitRangeForm.isEditing) {
      // 编辑模式：更新现有限制
      response = await k8sApi.updateLimitRange(selectedClusterId.value, currentNamespace.value.name, limitRangeForm.originalName, data)
    } else {
      // 创建模式：创建新限制
      response = await k8sApi.createLimitRange(selectedClusterId.value, currentNamespace.value.name, data)
    }
    
    responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(limitRangeForm.isEditing ? 'LimitRange更新成功' : 'LimitRange创建成功')
      createLimitRangeDialogVisible.value = false
      // 重新获取LimitRange列表
      manageLimitRanges(currentNamespace.value)
    } else {
      ElMessage.error(responseData.message || (limitRangeForm.isEditing ? '更新LimitRange失败' : '创建LimitRange失败'))
    }
  } catch (error) {
    console.error((limitRangeForm.isEditing ? '更新' : '创建') + 'LimitRange失败:', error)
    ElMessage.error((limitRangeForm.isEditing ? '更新' : '创建') + 'LimitRange失败，请检查网络连接')
  }
}

// 格式化资源值（添加单位显示）
const formatResourceValue = (value) => {
  if (!value || value === '-') return '-'
  
  // 如果已经有单位，直接返回
  if (typeof value === 'string' && (value.includes('m') || value.includes('Gi') || value.includes('Mi') || value.includes('Ki'))) {
    return value
  }
  
  // 纯数字，尝试添加合适的单位
  if (/^\d+$/.test(value)) {
    // CPU: 纯数字认为是核心数
    return `${value} cores`
  }
  
  return value
}

// 格式化数量值（添加单位显示）
const formatCountValue = (value) => {
  if (!value || value === '-') return '-'
  return `${value} 个`
}

// 编辑ResourceQuota
const editResourceQuota = (quota) => {
  // 填充表单数据
  quotaForm.name = quota.name
  quotaForm.cpuLimit = quota.hard?.['limits.cpu']?.replace(/cores?/i, '').trim() || ''
  quotaForm.memoryLimit = quota.hard?.['limits.memory']?.replace(/Gi?/i, '').trim() || ''
  quotaForm.storageLimit = quota.hard?.['requests.storage']?.replace(/Gi?/i, '').trim() || ''
  quotaForm.podLimit = quota.hard?.['count/pods'] || ''
  quotaForm.serviceLimit = quota.hard?.['count/services'] || ''
  
  // 标记为编辑模式
  quotaForm.isEditing = true
  quotaForm.originalName = quota.name
  
  createQuotaDialogVisible.value = true
}

// 删除ResourceQuota
const deleteResourceQuota = async (quota) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除资源配额 "${quota.name}" 吗？删除后将无法恢复！`,
      '删除资源配额确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'error',
      }
    )
    
    const response = await k8sApi.deleteResourceQuota(selectedClusterId.value, currentNamespace.value.name, quota.name)
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`资源配额 "${quota.name}" 删除成功`)
      // 重新获取ResourceQuota列表
      manageResourceQuotas(currentNamespace.value)
    } else {
      ElMessage.error(responseData.message || '删除资源配额失败')
    }
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消删除操作')
    } else {
      console.error('删除资源配额失败:', error)
      ElMessage.error('删除资源配额失败，请检查网络连接')
    }
  }
}

// 编辑LimitRange
const editLimitRange = (limitRange) => {
  // 填充表单数据
  limitRangeForm.name = limitRange.name
  const limit = limitRange.spec?.limits?.[0] || {}
  
  // 处理CPU值 - 移除单位后缀
  const extractCpuValue = (value) => {
    if (!value) return ''
    return value.toString().replace(/m$/i, '').replace(/cores?$/i, '').trim()
  }
  
  // 处理内存值 - 移除单位后缀
  const extractMemoryValue = (value) => {
    if (!value) return ''
    return value.toString().replace(/Mi$/i, '').replace(/Gi$/i, '').replace(/Ki$/i, '').trim()
  }
  
  limitRangeForm.containerCpuMin = extractCpuValue(limit.min?.cpu)
  limitRangeForm.containerCpuMax = extractCpuValue(limit.max?.cpu)
  limitRangeForm.containerMemoryMin = extractMemoryValue(limit.min?.memory)
  limitRangeForm.containerMemoryMax = extractMemoryValue(limit.max?.memory)
  limitRangeForm.containerCpuDefault = extractCpuValue(limit.default?.cpu)
  limitRangeForm.containerMemoryDefault = extractMemoryValue(limit.default?.memory)
  
  // 标记为编辑模式
  limitRangeForm.isEditing = true
  limitRangeForm.originalName = limitRange.name
  
  console.log('编辑LimitRange表单数据:', limitRangeForm)
  createLimitRangeDialogVisible.value = true
}

// 删除LimitRange
const deleteLimitRange = async (limitRange) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除资源限制 "${limitRange.name}" 吗？删除后将无法恢复！`,
      '删除资源限制确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'error',
      }
    )
    
    const response = await k8sApi.deleteLimitRange(selectedClusterId.value, currentNamespace.value.name, limitRange.name)
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`资源限制 "${limitRange.name}" 删除成功`)
      // 重新获取LimitRange列表
      manageLimitRanges(currentNamespace.value)
    } else {
      ElMessage.error(responseData.message || '删除资源限制失败')
    }
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消删除操作')
    } else {
      console.error('删除资源限制失败:', error)
      ElMessage.error('删除资源限制失败，请检查网络连接')
    }
  }
}

onMounted(async () => {
  await fetchClusterList()
  if (selectedClusterId.value) {
    handleQuery()
  }
})
</script>

<template>
  <div class="k8s-namespace-management">
    <el-card shadow="hover" class="namespace-card">
      <template #header>
        <div class="card-header">
          <span class="title">K8s 命名空间管理</span>
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
          <el-form-item label="命名空间名称">
            <el-input
              v-model="queryParams.namespaceName"
              placeholder="请输入命名空间名称"
              clearable
              size="small"
              style="width: 200px"
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select
              v-model="queryParams.status"
              placeholder="请选择状态"
              clearable
              size="small"
              style="width: 150px"
            >
              <el-option
                v-for="item in statusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :icon="Search" size="small" @click="handleQuery">
              搜索
            </el-button>
            <el-button type="warning" :icon="Refresh" size="small" @click="resetQuery">
              重置
            </el-button>
            <el-button type="success" :icon="Plus" v-authority="['k8s:namespace:add']" size="small" @click="showCreateNamespaceDialog">
              创建命名空间
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 命名空间列表表格 -->
      <NamespacesTable
        :table-data="pagedTableData"
        :loading="loading"
        @viewNamespaceDetail="viewNamespaceDetail"
        @viewNamespaceLabels="viewNamespaceLabels"
        @manageResourceQuotas="manageResourceQuotas"
        @manageLimitRanges="manageLimitRanges"
        @deleteNamespace="deleteNamespace"
        @editYaml="handleEditYaml"
        @manageLabels="handleManageLabels"
      />

      <!-- 分页 -->
      <div class="pagination-wrapper">
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

    <!-- 创建命名空间对话框 -->
    <CreateNamespaceDialog
      v-model:visible="createNamespaceDialogVisible"
      :selected-cluster-id="selectedClusterId"
      @success="handleCreateNamespace"
      ref="createDialogRef"
    />

    <!-- 命名空间标签修改对话框 -->
    <NamespaceLabelsDialog
      v-model:visible="labelsDialogVisible"
      :cluster-id="selectedClusterId"
      :namespace-name="currentNamespace.name || ''"
      :current-labels="currentNamespace.labels || {}"
      @saved="handleQuery"
    />

    <!-- 命名空间详情对话框 -->
    <el-dialog
      v-model="namespaceDetailDialogVisible"
      :title="`命名空间详情 - ${currentNamespace.name || ''}`"
      width="800px"
      class="namespace-detail-dialog"
    >
      <div class="namespace-detail-content" v-if="currentNamespace">
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="名称">{{ currentNamespace.name }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusTag(currentNamespace.status)" size="small">
              {{ getStatusText(currentNamespace.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatAge(currentNamespace.createdAt) }}</el-descriptions-item>
          <el-descriptions-item label="注解数量">{{ getAnnotationCount(currentNamespace.annotations) }}</el-descriptions-item>
        </el-descriptions>

        <div class="labels-view-content">
          <!-- 用户自定义标签 -->
          <div class="labels-section" v-if="Object.keys(getUserLabels(currentNamespace.labels)).length > 0">
            <h4 class="section-title">用户标签</h4>
            <div class="labels-table-wrapper">
              <el-table
                :data="toLabelTableData(getUserLabels(currentNamespace.labels))"
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
          <div class="labels-section" v-if="Object.keys(getSystemLabels(currentNamespace.labels)).length > 0">
            <h4 class="section-title">系统标签</h4>
            <div class="labels-table-wrapper">
              <el-table
                :data="toLabelTableData(getSystemLabels(currentNamespace.labels))"
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
          <div v-if="!currentNamespace.labels || Object.keys(currentNamespace.labels).length === 0" class="no-labels">
            <el-empty description="该命名空间没有标签" :image-size="60" />
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- ResourceQuota管理对话框 -->
    <el-dialog
      v-model="resourceQuotaDialogVisible"
      :title="`资源配额管理 - ${currentNamespace.name || ''}`"
      width="900px"
      class="resource-quota-dialog"
    >
      <div class="quota-content">
        <!-- ResourceQuota列表 -->
        <div class="quota-list-section">
          <div class="section-header">
            <h4>当前资源配额</h4>
            <el-button type="primary" size="small" :icon="Plus" v-authority="['k8s:namespace:setupadd']" @click="showCreateQuotaDialog">添加配额</el-button>
          </div>
          
          <el-table :data="resourceQuotas" size="small" style="width: 100%">
            <el-table-column prop="name" label="配额名称" min-width="150" />
            <el-table-column label="CPU限制" min-width="100">
              <template #default="{ row }">
                <span>{{ formatResourceValue(row.hard?.['limits.cpu']) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="内存限制" min-width="100">
              <template #default="{ row }">
                <span>{{ formatResourceValue(row.hard?.['limits.memory']) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="Pod数量限制" min-width="120">
              <template #default="{ row }">
                <span>{{ formatCountValue(row.hard?.['count/pods']) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <div class="quota-operation-buttons">
                  <el-button type="primary" size="small" :icon="Edit" @click="editResourceQuota(row)">编辑</el-button>
                  <el-button type="danger" size="small" :icon="Delete" @click="deleteResourceQuota(row)">删除</el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
          
          <div v-if="resourceQuotas.length === 0" class="empty-state">
            <el-empty description="暂无资源配额" :image-size="60" />
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- LimitRange管理对话框 -->
    <el-dialog
      v-model="limitRangeDialogVisible"
      :title="`资源限制管理 - ${currentNamespace.name || ''}`"
      width="900px"
      class="limit-range-dialog"
    >
      <div class="limit-range-content">
        <!-- LimitRange列表 -->
        <div class="limit-range-list-section">
          <div class="section-header">
            <h4>当前资源限制</h4>
            <el-button type="primary" size="small" :icon="Plus" v-authority="['k8s:namespace:restrictionadd']" @click="showCreateLimitRangeDialog">添加限制</el-button>
          </div>
          
          <el-table :data="limitRanges" size="small" style="width: 100%">
            <el-table-column prop="name" label="限制名称" min-width="150" />
            <el-table-column label="容器CPU最小值" min-width="120">
              <template #default="{ row }">
                <el-tag size="small" type="info">
                  {{ row.spec?.limits?.[0]?.min?.cpu || '-' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="容器CPU最大值" min-width="120">
              <template #default="{ row }">
                <el-tag size="small" type="warning">
                  {{ row.spec?.limits?.[0]?.max?.cpu || '-' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="容器内存最小值" min-width="130">
              <template #default="{ row }">
                <el-tag size="small" type="info">
                  {{ row.spec?.limits?.[0]?.min?.memory || '-' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="容器内存最大值" min-width="130">
              <template #default="{ row }">
                <el-tag size="small" type="warning">
                  {{ row.spec?.limits?.[0]?.max?.memory || '-' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="创建时间" min-width="120">
              <template #default="{ row }">
                <el-tag size="small">
                  {{ formatAge(row.createdAt) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <div class="limit-range-operation-buttons">
                  <el-button type="primary" size="small" :icon="Edit" @click="editLimitRange(row)">编辑</el-button>
                  <el-button type="danger" size="small" :icon="Delete" @click="deleteLimitRange(row)">删除</el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
          
          <div v-if="limitRanges.length === 0" class="empty-state">
            <el-empty description="暂无资源限制" :image-size="60" />
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 创建ResourceQuota对话框 -->
    <el-dialog
      v-model="createQuotaDialogVisible"
      :title="quotaForm.isEditing ? '编辑资源配额' : '创建资源配额'"
      width="600px"
      class="create-quota-dialog"
    >
      <el-form :model="quotaForm" label-width="120px">
        <el-form-item label="配额名称" required>
          <el-input
            v-model="quotaForm.name"
            placeholder="请输入配额名称，如：compute-quota"
          />
        </el-form-item>
        <el-form-item label="CPU限制">
          <el-input
            v-model="quotaForm.cpuLimit"
            placeholder="如：4（自动添加核心单位）或 4000m（毫核心）"
          >
            <template #append>cores</template>
          </el-input>
        </el-form-item>
        <el-form-item label="内存限制">
          <el-input
            v-model="quotaForm.memoryLimit"
            placeholder="如：8（自动添加Gi单位）或 8192Mi"
          >
            <template #append>Gi</template>
          </el-input>
        </el-form-item>
        <el-form-item label="存储限制">
          <el-input
            v-model="quotaForm.storageLimit"
            placeholder="如：100（自动添加Gi单位）或 100Gi"
          >
            <template #append>Gi</template>
          </el-input>
        </el-form-item>
        <el-form-item label="Pod数量限制">
          <el-input
            v-model="quotaForm.podLimit"
            placeholder="如：50（最多50个Pod）"
          />
        </el-form-item>
        <el-form-item label="Service数量限制">
          <el-input
            v-model="quotaForm.serviceLimit"
            placeholder="如：20（最多20个Service）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createQuotaDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitCreateQuota">{{ quotaForm.isEditing ? '更新' : '创建' }}</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 创建LimitRange对话框 -->
    <el-dialog
      v-model="createLimitRangeDialogVisible"
      :title="limitRangeForm.isEditing ? '编辑资源限制' : '创建资源限制'"
      width="700px"
      class="create-limit-range-dialog"
    >
      <el-form :model="limitRangeForm" label-width="150px">
        <el-form-item label="限制名称" required>
          <el-input
            v-model="limitRangeForm.name"
            placeholder="请输入限制名称，如：container-limits"
          />
        </el-form-item>
        <el-divider>容器资源限制</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="CPU最小值">
              <el-input
                v-model="limitRangeForm.containerCpuMin"
                placeholder="如：100（自动添加m单位）"
              >
                <template #append>m</template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="CPU最大值">
              <el-input
                v-model="limitRangeForm.containerCpuMax"
                placeholder="如：2（核心数）"
              >
                <template #append>cores</template>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="内存最小值">
              <el-input
                v-model="limitRangeForm.containerMemoryMin"
                placeholder="如：128（自动添加Mi单位）"
              >
                <template #append>Mi</template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="内存最大值">
              <el-input
                v-model="limitRangeForm.containerMemoryMax"
                placeholder="如：4（自动添加Gi单位）"
              >
                <template #append>Gi</template>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-divider>默认资源配置</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="默认CPU">
              <el-input
                v-model="limitRangeForm.containerCpuDefault"
                placeholder="如：500（自动添加m单位）"
              >
                <template #append>m</template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="默认内存">
              <el-input
                v-model="limitRangeForm.containerMemoryDefault"
                placeholder="如：1（自动添加Mi单位）"
              >
                <template #append>Mi</template>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createLimitRangeDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitCreateLimitRange">{{ limitRangeForm.isEditing ? '更新' : '创建' }}</el-button>
        </span>
      </template>
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

export { getProgressColor }
</script>

<style scoped>
.k8s-namespace-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.namespace-card {
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
.namespace-detail-dialog :deep(.el-dialog),
.resource-quota-dialog :deep(.el-dialog),
.limit-range-dialog :deep(.el-dialog),
.create-quota-dialog :deep(.el-dialog),
.create-limit-range-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.namespace-detail-dialog :deep(.el-dialog__header),
.resource-quota-dialog :deep(.el-dialog__header),
.limit-range-dialog :deep(.el-dialog__header),
.create-quota-dialog :deep(.el-dialog__header),
.create-limit-range-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.create-namespace-dialog :deep(.el-dialog__title),
.namespace-detail-dialog :deep(.el-dialog__title),
.resource-quota-dialog :deep(.el-dialog__title),
.limit-range-dialog :deep(.el-dialog__title),
.create-quota-dialog :deep(.el-dialog__title),
.create-limit-range-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.create-namespace-dialog :deep(.el-dialog__body),
.namespace-detail-dialog :deep(.el-dialog__body),
.resource-quota-dialog :deep(.el-dialog__body),
.limit-range-dialog :deep(.el-dialog__body),
.create-quota-dialog :deep(.el-dialog__body),
.create-limit-range-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

.namespace-detail-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.labels-view-content {
  display: flex;
  flex-direction: column;
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

.no-labels {
  text-align: center;
  padding: 40px 20px;
  color: #909399;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h4 {
  margin: 0;
  color: #6b9bd1;
  font-weight: 600;
  font-size: 16px;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
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

/* ResourceQuota和LimitRange操作按钮样式 */
.quota-operation-buttons,
.limit-range-operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
  align-items: center;
  flex-wrap: nowrap;
  white-space: nowrap;
}

.quota-operation-buttons .el-button,
.limit-range-operation-buttons .el-button {
  transition: all 0.3s ease;
  flex-shrink: 0;
  min-width: auto;
}

.quota-operation-buttons .el-button:hover,
.limit-range-operation-buttons .el-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* 表单提示样式 */
.form-hint {
  margin-top: 4px;
  padding: 4px 8px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 4px;
  border-left: 3px solid rgba(103, 126, 234, 0.3);
}

.form-hint .el-text {
  line-height: 1.4;
}

.create-namespace-dialog :deep(.el-divider__text) {
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 600;
}

.create-namespace-dialog :deep(.el-divider) {
  margin: 20px 0 16px 0;
}

.create-namespace-dialog :deep(.el-textarea__inner) {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .operation-buttons {
    gap: 4px;
  }
  
  .operation-buttons .el-button {
    margin: 1px;
  }
  
  .quota-operation-buttons,
  .limit-range-operation-buttons {
    gap: 4px;
  }
  
  .quota-operation-buttons .el-button,
  .limit-range-operation-buttons .el-button {
    font-size: 12px;
    padding: 4px 8px;
  }
  
  .header-actions .el-select {
    min-width: 200px;
  }
}

@media (max-width: 768px) {
  .k8s-namespace-management {
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
  
  .namespace-table :deep(.el-table__row:hover) {
    transform: none;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
  padding: 12px 0;
}
</style>
