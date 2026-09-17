<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  Edit,
  Delete,
  View,
  ArrowDown,
  Connection,
  Setting,
  Monitor,
  Upload,
  Download,
  Document,
  CopyDocument,
  Link
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import PodYamlDialog from './pods/PodYamlDialog.vue'
import ClusterSelector from './pods/ClusterSelector.vue'
import NamespaceSelector from './pods/NamespaceSelector.vue'
import { getYamlContentFromApiResponse } from './utils/yaml'

const router = useRouter()
const route = useRoute()

// 集群和命名空间选择
const selectedClusterId = ref('')
const queryParams = reactive({
  namespace: 'default'
})

// 页面状态
const loading = ref(false)
const activeTab = ref('service')
const searchKeyword = ref('')

// Service 相关状态
const serviceList = ref([])
const serviceLoading = ref(false)

// 监听 serviceList 的变化，确保它始终是数组
watch(serviceList, (newValue) => {
  if (!Array.isArray(newValue)) {
    console.warn('serviceList 被设置为非数组值:', newValue, '重置为空数组')
    serviceList.value = []
  }
}, { immediate: true })
const serviceDialogVisible = ref(false)
const editServiceDialogVisible = ref(false)
const serviceDetailDialogVisible = ref(false)
const serviceYamlDialogVisible = ref(false)
const serviceEventsDialogVisible = ref(false)
const currentService = reactive({
  name: '',
  labels: {},
  type: 'ClusterIP',
  selector: {},
  ports: [{ name: '', port: 80, targetPort: 80, protocol: 'TCP' }],
  description: ''
})
const currentServiceForDetail = ref({})
const currentServiceYaml = ref('')
const currentServiceForEvents = ref({})
const serviceEvents = ref([])

// Ingress 相关状态
const ingressList = ref([])
const ingressLoading = ref(false)

// 监听 ingressList 的变化，确保它始终是数组
watch(ingressList, (newValue) => {
  if (!Array.isArray(newValue)) {
    console.warn('ingressList 被设置为非数组值:', newValue, '重置为空数组')
    ingressList.value = []
  }
}, { immediate: true })
const ingressDialogVisible = ref(false)
const editIngressDialogVisible = ref(false)
const ingressDetailDialogVisible = ref(false)
const ingressYamlDialogVisible = ref(false)
const ingressEventsDialogVisible = ref(false)
const currentIngress = reactive({
  name: '',
  annotations: {},
  rules: [{ host: '', paths: [{ path: '/', pathType: 'Prefix', serviceName: '', servicePort: 80 }] }],
  tls: [],
  description: ''
})
const currentIngressForDetail = ref({})
const currentIngressYaml = ref('')
const currentIngressForEvents = ref({})
const ingressEvents = ref([])

// Ingress 测试相关状态
const ingressTestDialogVisible = ref(false)
const ingressTestLoading = ref(false)
const ingressTestResult = ref(null)
const ingressTestForm = reactive({
  serviceName: '',
  servicePort: 80,
  path: '/',
  host: '',
  method: 'GET',
  timeout: 10
})

// Service 类型选项
const serviceTypeOptions = [
  { label: 'ClusterIP', value: 'ClusterIP', description: '集群内部访问' },
  { label: 'NodePort', value: 'NodePort', description: '节点端口访问' },
  { label: 'LoadBalancer', value: 'LoadBalancer', description: '负载均衡器访问' },
  { label: 'ExternalName', value: 'ExternalName', description: '外部名称映射' }
]

// 协议选项
const protocolOptions = [
  { label: 'TCP', value: 'TCP' },
  { label: 'UDP', value: 'UDP' },
  { label: 'SCTP', value: 'SCTP' }
]

// 路径类型选项
const pathTypeOptions = [
  { label: 'Prefix', value: 'Prefix', description: '前缀匹配' },
  { label: 'Exact', value: 'Exact', description: '精确匹配' },
  { label: 'ImplementationSpecific', value: 'ImplementationSpecific', description: '实现特定' }
]

// 过滤后的列表
const filteredServiceList = computed(() => {
  // 确保 serviceList.value 是数组
  const list = Array.isArray(serviceList.value) ? serviceList.value : []

  if (!searchKeyword.value) return list
  return list.filter(item =>
    item.name?.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    item.type?.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

const filteredIngressList = computed(() => {
  // 确保 ingressList.value 是数组
  const list = Array.isArray(ingressList.value) ? ingressList.value : []

  if (!searchKeyword.value) return list
  return list.filter(item =>
    item.name?.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    (item.rules && Array.isArray(item.rules) && item.rules.some(rule =>
      rule.host?.toLowerCase().includes(searchKeyword.value.toLowerCase())
    ))
  )
})

// 处理集群选择变化
const handleClusterChange = (clusterId) => {
  selectedClusterId.value = clusterId
  console.log('集群选择变化:', clusterId)
  if (clusterId && queryParams.namespace) {
    loadAllNetworkResources()
  }
}

// 处理命名空间选择变化
const handleNamespaceChange = (namespace) => {
  queryParams.namespace = namespace
  console.log('命名空间选择变化:', namespace)
  if (selectedClusterId.value && namespace) {
    loadAllNetworkResources()
  }
}

// 加载所有网络资源
const loadAllNetworkResources = () => {
  console.log('加载网络资源 - 集群:', selectedClusterId.value, '命名空间:', queryParams.namespace)

  // 清空现有数据
  serviceList.value = []
  ingressList.value = []

  // 重新加载数据
  if (activeTab.value === 'service') {
    fetchServiceList()
  } else {
    fetchIngressList()
  }
}

// 监听集群和命名空间变化，自动加载数据
watch(
  [selectedClusterId, () => queryParams.namespace],
  ([clusterId, namespace]) => {
    console.log('监听到变化 - 集群ID:', clusterId, '命名空间:', namespace)
    if (clusterId && namespace) {
      console.log('集群和命名空间都已选择，开始加载网络资源')
      loadAllNetworkResources()
    }
  },
  { immediate: true }
)

// 获取 Service 列表
const fetchServiceList = async () => {
  if (!selectedClusterId.value) {
    console.warn('集群ID为空，无法获取 Service 列表')
    return
  }

  try {
    serviceLoading.value = true
    console.log('正在获取 Service 列表，集群ID:', selectedClusterId.value, '命名空间:', queryParams.namespace)

    const response = await k8sApi.getServiceList(selectedClusterId.value, queryParams.namespace)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      // 后端返回的数据结构是 { data: { services: [...] } }
      // 处理 services 为 null 的情况
      const rawData = responseData.data?.services || responseData.data
      // 确保获取到的数据是数组格式，处理 null 和 undefined
      serviceList.value = Array.isArray(rawData) ? rawData : []
      console.log('获取到 Service 列表:', serviceList.value.length, '个')
    } else {
      ElMessage.error(responseData.message || '获取 Service 列表失败')
      serviceList.value = [] // 失败时确保是空数组
    }
  } catch (error) {
    console.error('获取 Service 列表失败:', error)
    ElMessage.error('获取 Service 列表失败，请检查网络连接')
    serviceList.value = [] // 异常时确保是空数组
  } finally {
    serviceLoading.value = false
  }
}

// 获取 Ingress 列表
const fetchIngressList = async () => {
  if (!selectedClusterId.value) {
    console.warn('集群ID为空，无法获取 Ingress 列表')
    return
  }

  try {
    ingressLoading.value = true
    console.log('正在获取 Ingress 列表，集群ID:', selectedClusterId.value, '命名空间:', queryParams.namespace)

    const response = await k8sApi.getIngressList(selectedClusterId.value, queryParams.namespace)
    console.log('Ingress API原始响应:', response)

    const responseData = response.data || response
    console.log('Ingress API响应数据:', responseData)

    if (responseData.code === 200 || responseData.success) {
      // 后端返回的数据结构可能是 { data: { ingresses: [...] } } 或类似结构
      // 处理 ingresses 为 null 的情况
      const rawData = responseData.data?.ingresses || responseData.data
      // 确保获取到的数据是数组格式，处理 null 和 undefined
      ingressList.value = Array.isArray(rawData) ? rawData : []
      console.log('获取到 Ingress 列表:', ingressList.value.length, '个')
    } else {
      console.error('Ingress API返回错误:', responseData)
      ElMessage.error(responseData.message || '获取 Ingress 列表失败')
      ingressList.value = [] // 失败时确保是空数组
    }
  } catch (error) {
    console.error('获取 Ingress 列表失败:', error)
    console.error('错误详情:', {
      message: error?.message,
      response: error?.response,
      request: error?.request
    })
    ElMessage.error(error?.response?.data?.message || '获取 Ingress 列表失败，请检查网络连接')
    ingressList.value = [] // 异常时确保是空数组
  } finally {
    ingressLoading.value = false
  }
}

// 刷新数据
const handleRefresh = () => {
  if (activeTab.value === 'service') {
    fetchServiceList()
  } else {
    fetchIngressList()
  }
}

// 标签页切换
const handleTabChange = (tab) => {
  activeTab.value = tab
  searchKeyword.value = ''
  if (tab === 'service') {
    fetchServiceList()
  } else {
    fetchIngressList()
  }
}

// Service 操作
const handleCreateService = () => {
  if (!selectedClusterId.value) {
    ElMessage.warning('请先选择集群')
    return
  }
  if (!queryParams.namespace) {
    ElMessage.warning('请先选择命名空间')
    return
  }

  // 设置默认的Service YAML模板
  const defaultYaml = `apiVersion: v1
kind: Service
metadata:
  name: new-service
  namespace: ${queryParams.namespace}
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
    name: http
  selector:
    app: my-app`

  currentServiceYaml.value = defaultYaml
  currentServiceForDetail.value = { name: 'new-service', namespace: queryParams.namespace }
  serviceYamlDialogVisible.value = true
}

// 查看 Service 详情
const handleViewService = async (row) => {
  try {
    const response = await k8sApi.getServiceDetail(selectedClusterId.value, queryParams.namespace, row.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      currentServiceForDetail.value = responseData.data || row
      serviceDetailDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取 Service 详情失败')
    }
  } catch (error) {
    console.error('获取 Service 详情失败:', error)
    ElMessage.error('获取 Service 详情失败')
  }
}

// 编辑 Service YAML
const handleEditServiceYaml = async (row) => {
  try {
    loading.value = true
    const response = await k8sApi.getServiceYaml(selectedClusterId.value, queryParams.namespace, row.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      currentServiceYaml.value = getYamlContentFromApiResponse(
        responseData,
        `# Service ${row.name} YAML\napiVersion: v1\nkind: Service\nmetadata:\n  name: ${row.name}`
      )
      currentServiceForDetail.value = row
      serviceYamlDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取 Service YAML 失败')
    }
  } catch (error) {
    console.error('获取 Service YAML 失败:', error)
    ElMessage.error('获取 Service YAML 失败')
  } finally {
    loading.value = false
  }
}

// 查看 Service 事件
const handleViewServiceEvents = async (row) => {
  try {
    currentServiceForEvents.value = row
    const response = await k8sApi.getServiceEvents(selectedClusterId.value, queryParams.namespace, row.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      serviceEvents.value = responseData.data || []
    } else {
      serviceEvents.value = []
      ElMessage.warning('暂无事件数据')
    }
    serviceEventsDialogVisible.value = true
  } catch (error) {
    console.error('获取 Service 事件失败:', error)
    serviceEvents.value = []
    serviceEventsDialogVisible.value = true
  }
}

const handleEditService = async (row) => {
  try {
    const response = await k8sApi.getServiceDetail(selectedClusterId.value, queryParams.namespace, row.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      const serviceData = responseData.data
      Object.assign(currentService, {
        name: serviceData.name,
        labels: serviceData.labels || {},
        type: serviceData.type,
        selector: serviceData.selector || {},
        ports: serviceData.ports || [{ name: '', port: 80, targetPort: 80, protocol: 'TCP' }],
        description: serviceData.description || ''
      })
      editServiceDialogVisible.value = true
    }
  } catch (error) {
    ElMessage.error('获取 Service 详情失败')
  }
}

const handleDeleteService = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 Service "${row.name}" 吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const response = await k8sApi.deleteService(selectedClusterId.value, queryParams.namespace, row.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('Service 删除成功')
      fetchServiceList()
    } else {
      ElMessage.error(responseData.message || 'Service 删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('Service 删除失败')
    }
  }
}

// Ingress 操作
const handleCreateIngress = () => {
  if (!selectedClusterId.value) {
    ElMessage.warning('请先选择集群')
    return
  }
  if (!queryParams.namespace) {
    ElMessage.warning('请先选择命名空间')
    return
  }

  // 设置默认的Ingress YAML模板（使用v1beta1版本以兼容旧版本K8s）
  const defaultYaml = `apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: new-ingress
  namespace: ${queryParams.namespace}
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          serviceName: my-service
          servicePort: 80`

  currentIngressYaml.value = defaultYaml
  currentIngressForDetail.value = { name: 'new-ingress', namespace: queryParams.namespace }
  ingressYamlDialogVisible.value = true
}

// 查看 Ingress 详情
const handleViewIngress = async (row) => {
  try {
    const response = await k8sApi.getIngressDetail(selectedClusterId.value, queryParams.namespace, row.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      currentIngressForDetail.value = responseData.data || row
      ingressDetailDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取 Ingress 详情失败')
    }
  } catch (error) {
    console.error('获取 Ingress 详情失败:', error)
    ElMessage.error('获取 Ingress 详情失败')
  }
}

// 编辑 Ingress YAML
const handleEditIngressYaml = async (row) => {
  try {
    loading.value = true
    const response = await k8sApi.getIngressYaml(selectedClusterId.value, queryParams.namespace, row.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      currentIngressYaml.value = getYamlContentFromApiResponse(
        responseData,
        `# Ingress ${row.name} YAML\napiVersion: networking.k8s.io/v1\nkind: Ingress\nmetadata:\n  name: ${row.name}`
      )
      currentIngressForDetail.value = row
      ingressYamlDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取 Ingress YAML 失败')
    }
  } catch (error) {
    console.error('获取 Ingress YAML 失败:', error)
    ElMessage.error('获取 Ingress YAML 失败')
  } finally {
    loading.value = false
  }
}

// 查看 Ingress 事件
const handleViewIngressEvents = async (row) => {
  try {
    currentIngressForEvents.value = row
    const response = await k8sApi.getIngressEvents(selectedClusterId.value, queryParams.namespace, row.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ingressEvents.value = responseData.data || []
    } else {
      ingressEvents.value = []
      ElMessage.warning('暂无事件数据')
    }
    ingressEventsDialogVisible.value = true
  } catch (error) {
    console.error('获取 Ingress 事件失败:', error)
    ingressEvents.value = []
    ingressEventsDialogVisible.value = true
  }
}

// 测试 Ingress 后端服务
const handleTestIngress = async (row) => {
  try {
    // 从 Ingress 数据中提取第一个规则的信息作为默认值
    if (row.rules && row.rules.length > 0) {
      const firstRule = row.rules[0]
      ingressTestForm.host = firstRule.host || ''

      if (firstRule.http?.paths && firstRule.http.paths.length > 0) {
        const firstPath = firstRule.http.paths[0]
        ingressTestForm.path = firstPath.path || '/'
        ingressTestForm.serviceName = firstPath.backend?.service?.name || firstPath.backend?.serviceName || ''
        ingressTestForm.servicePort = firstPath.backend?.service?.port?.number || firstPath.backend?.servicePort || 80
      }
    }

    ingressTestResult.value = null
    ingressTestDialogVisible.value = true
  } catch (error) {
    console.error('打开测试对话框失败:', error)
    ElMessage.error('打开测试对话框失败')
  }
}

// 执行 Ingress 后端测试
const executeIngressTest = async () => {
  if (!ingressTestForm.serviceName) {
    ElMessage.warning('请输入服务名称')
    return
  }

  if (!ingressTestForm.servicePort) {
    ElMessage.warning('请输入服务端口')
    return
  }

  try {
    ingressTestLoading.value = true
    ingressTestResult.value = null

    const response = await k8sApi.testIngressBackend(
      selectedClusterId.value,
      queryParams.namespace,
      {
        serviceName: ingressTestForm.serviceName,
        servicePort: ingressTestForm.servicePort,
        path: ingressTestForm.path || '/',
        host: ingressTestForm.host || '',
        method: ingressTestForm.method || 'GET',
        timeout: ingressTestForm.timeout || 10
      }
    )

    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ingressTestResult.value = responseData.data || responseData

      if (responseData.data?.success || responseData.success) {
        ElMessage.success('后端服务测试完成')
      } else {
        ElMessage.warning('后端服务测试发现问题,请查看详细信息')
      }
    } else {
      ElMessage.error(responseData.message || '测试失败')
      ingressTestResult.value = {
        success: false,
        status: '异常',
        message: responseData.message || '测试失败'
      }
    }
  } catch (error) {
    console.error('测试 Ingress 后端失败:', error)
    ElMessage.error(error.response?.data?.message || '测试失败,请检查网络连接')
    ingressTestResult.value = {
      success: false,
      status: '异常',
      message: error.response?.data?.message || error.message || '测试失败'
    }
  } finally {
    ingressTestLoading.value = false
  }
}

// 重置测试表单
const resetTestForm = () => {
  Object.assign(ingressTestForm, {
    serviceName: '',
    servicePort: 80,
    path: '/',
    host: '',
    method: 'GET',
    timeout: 10
  })
  ingressTestResult.value = null
}

// 保存 Service YAML
const handleServiceYamlSave = async (data) => {
  try {
    // 判断是否是创建新Service
    // 条件：资源名为 'new-service' 或者 currentServiceForDetail 没有 name 属性
    const isCreating = data.resourceName === 'new-service' || !currentServiceForDetail.value?.name

    let response
    if (isCreating) {
      // 创建新Service，使用createPodFromYaml API（支持多种资源类型）
      response = await k8sApi.createPodFromYaml(selectedClusterId.value, queryParams.namespace, { yamlContent: data.yamlContent })
    } else {
      // 更新现有Service，使用 PUT 接口
      response = await k8sApi.updateServiceYaml(selectedClusterId.value, queryParams.namespace, data.resourceName, data.yamlContent)
    }

    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(isCreating ? 'Service 创建成功' : 'Service YAML 更新成功')
      serviceYamlDialogVisible.value = false
      fetchServiceList() // 刷新列表
    } else {
      ElMessage.error(responseData.message || (isCreating ? 'Service 创建失败' : 'Service YAML 更新失败'))
    }
  } catch (error) {
    console.error('Service 操作失败:', error)
    ElMessage.error(error.response?.data?.message || 'Service 操作失败')
  }
}

// 保存 Ingress YAML
const handleIngressYamlSave = async (data) => {
  try {
    // 判断是否是创建新Ingress
    // 条件：资源名为 'new-ingress' 或者 currentIngressForDetail 没有 name 属性
    const isCreating = data.resourceName === 'new-ingress' || !currentIngressForDetail.value?.name

    let response
    if (isCreating) {
      // 创建新Ingress，使用createPodFromYaml API（支持多种资源类型）
      response = await k8sApi.createPodFromYaml(selectedClusterId.value, queryParams.namespace, { yamlContent: data.yamlContent })
    } else {
      // 更新现有Ingress，使用 PUT 接口
      response = await k8sApi.updateIngressYaml(selectedClusterId.value, queryParams.namespace, data.resourceName, data.yamlContent)
    }

    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(isCreating ? 'Ingress 创建成功' : 'Ingress YAML 更新成功')
      ingressYamlDialogVisible.value = false
      fetchIngressList() // 刷新列表
    } else {
      ElMessage.error(responseData.message || (isCreating ? 'Ingress 创建失败' : 'Ingress YAML 更新失败'))
    }
  } catch (error) {
    console.error('Ingress 操作失败:', error)
    ElMessage.error(error.response?.data?.message || 'Ingress 操作失败')
  }
}

const handleEditIngress = async (row) => {
  try {
    const response = await k8sApi.getIngressDetail(selectedClusterId.value, queryParams.namespace, row.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      const ingressData = responseData.data
      Object.assign(currentIngress, {
        name: ingressData.name,
        annotations: ingressData.annotations || {},
        rules: ingressData.rules || [{ host: '', paths: [{ path: '/', pathType: 'Prefix', serviceName: '', servicePort: 80 }] }],
        tls: ingressData.tls || [],
        description: ingressData.description || ''
      })
      editIngressDialogVisible.value = true
    }
  } catch (error) {
    ElMessage.error('获取 Ingress 详情失败')
  }
}

const handleDeleteIngress = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 Ingress "${row.name}" 吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const response = await k8sApi.deleteIngress(selectedClusterId.value, queryParams.namespace, row.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('Ingress 删除成功')
      fetchIngressList()
    } else {
      ElMessage.error(responseData.message || 'Ingress 删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('Ingress 删除失败')
    }
  }
}

// 表单重置
const resetServiceForm = () => {
  Object.assign(currentService, {
    name: '',
    labels: {},
    type: 'ClusterIP',
    selector: {},
    ports: [{ name: '', port: 80, targetPort: 80, protocol: 'TCP' }],
    description: ''
  })
}

const resetIngressForm = () => {
  Object.assign(currentIngress, {
    name: '',
    annotations: {},
    rules: [{ host: '', paths: [{ path: '/', pathType: 'Prefix', serviceName: '', servicePort: 80 }] }],
    tls: [],
    description: ''
  })
}

// Service 端口管理
const addServicePort = () => {
  currentService.ports.push({ name: '', port: 80, targetPort: 80, protocol: 'TCP' })
}

const removeServicePort = (index) => {
  if (currentService.ports.length > 1) {
    currentService.ports.splice(index, 1)
  }
}

// Ingress 规则管理
const addIngressRule = () => {
  currentIngress.rules.push({ host: '', paths: [{ path: '/', pathType: 'Prefix', serviceName: '', servicePort: 80 }] })
}

const removeIngressRule = (index) => {
  if (currentIngress.rules.length > 1) {
    currentIngress.rules.splice(index, 1)
  }
}

const addIngressPath = (ruleIndex) => {
  currentIngress.rules[ruleIndex].paths.push({ path: '/', pathType: 'Prefix', serviceName: '', servicePort: 80 })
}

const removeIngressPath = (ruleIndex, pathIndex) => {
  if (currentIngress.rules[ruleIndex].paths.length > 1) {
    currentIngress.rules[ruleIndex].paths.splice(pathIndex, 1)
  }
}

// 提交表单
const submitServiceForm = async () => {
  try {
    const response = await k8sApi.createService(selectedClusterId.value, queryParams.namespace, currentService)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('Service 创建成功')
      serviceDialogVisible.value = false
      fetchServiceList()
    } else {
      ElMessage.error(responseData.message || 'Service 创建失败')
    }
  } catch (error) {
    ElMessage.error('Service 创建失败')
  }
}

const submitEditServiceForm = async () => {
  try {
    const response = await k8sApi.updateService(selectedClusterId.value, queryParams.namespace, currentService.name, currentService)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('Service 更新成功')
      editServiceDialogVisible.value = false
      fetchServiceList()
    } else {
      ElMessage.error(responseData.message || 'Service 更新失败')
    }
  } catch (error) {
    ElMessage.error('Service 更新失败')
  }
}

const submitIngressForm = async () => {
  try {
    const response = await k8sApi.createIngress(selectedClusterId.value, queryParams.namespace, currentIngress)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('Ingress 创建成功')
      ingressDialogVisible.value = false
      fetchIngressList()
    } else {
      ElMessage.error(responseData.message || 'Ingress 创建失败')
    }
  } catch (error) {
    ElMessage.error('Ingress 创建失败')
  }
}

const submitEditIngressForm = async () => {
  try {
    const response = await k8sApi.updateIngress(selectedClusterId.value, queryParams.namespace, currentIngress.name, currentIngress)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('Ingress 更新成功')
      editIngressDialogVisible.value = false
      fetchIngressList()
    } else {
      ElMessage.error(responseData.message || 'Ingress 更新失败')
    }
  } catch (error) {
    ElMessage.error('Ingress 更新失败')
  }
}

// 复制名称到剪贴板
const copyToClipboard = async (text, resourceType = '资源') => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(`${resourceType}名称已复制到剪贴板`)
  } catch (error) {
    console.error('复制失败:', error)
    // 降级方案：使用传统的复制方法
    try {
      const textArea = document.createElement('textarea')
      textArea.value = text
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      ElMessage.success(`${resourceType}名称已复制到剪贴板`)
    } catch (fallbackError) {
      ElMessage.error('复制失败，请手动复制')
    }
  }
}

// 组件挂载时初始化
onMounted(async () => {
  console.log('K8s网络管理页面初始化完成')
})
</script>

<template>
  <div class="k8s-network-management">
    <el-card shadow="hover" class="network-card">
      <template #header>
        <div class="card-header">
          <span class="title">K8s 网络管理</span>
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

      <!-- 搜索表单和资源类型选择 -->
      <div class="search-section">
        <el-form :inline="true" class="search-form">
          <el-form-item label="网络资源名称">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索资源..."
              clearable
              size="small"
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :icon="Search" size="small" @click="handleRefresh">
              搜索
            </el-button>
            <el-button type="success" :icon="Refresh" size="small" @click="handleRefresh">
              刷新
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 资源类型选择器和创建按钮 -->
        <div class="resource-controls">
          <el-tabs v-model="activeTab" @tab-change="handleTabChange" class="network-tabs">
            <el-tab-pane label="Service" name="service">
              <template #label>
                <span class="tab-label">Service</span>
              </template>
            </el-tab-pane>
            <el-tab-pane label="Ingress" name="ingress">
              <template #label>
                <span class="tab-label">Ingress</span>
              </template>
            </el-tab-pane>
          </el-tabs>

          <div class="create-actions">
            <el-button
              v-if="activeTab === 'service'"
              type="primary"
              :icon="Plus"
              size="small"
              v-authority="['k8s:network:addservice']"
              @click="handleCreateService"
            >
              创建 Service
            </el-button>
            <el-button
              v-if="activeTab === 'ingress'"
              type="primary"
              :icon="Plus"
              size="small"
              v-authority="['k8s:network:addingress']"
              @click="handleCreateIngress"
            >
              创建 Ingress
            </el-button>
          </div>
        </div>
      </div>

      <!-- 网络资源表格 -->
      <div class="table-content">
        <!-- Service 表格 -->
        <div v-if="activeTab === 'service'" class="tab-content">

          <el-table
            :data="filteredServiceList"
            v-loading="serviceLoading"
            element-loading-text="加载中..."
            class="resource-table"
            empty-text="暂无 Service 资源"
          >
              <el-table-column prop="name" label="名称" min-width="180">
                <template #default="{ row }">
                  <div class="resource-name">
                    <el-icon class="resource-icon"><Connection /></el-icon>
                    <span class="resource-name-link" @click="handleViewService(row)">{{ row.name }}</span>
                    <el-tooltip content="复制名称" placement="top">
                      <el-button
                        :icon="CopyDocument"
                        size="small"
                        type="text"
                        class="copy-button"
                        @click.stop="copyToClipboard(row.name, 'Service')"
                      />
                    </el-tooltip>
                  </div>
                </template>
              </el-table-column>

              <el-table-column prop="type" label="类型" width="120">
                <template #default="{ row }">
                  <el-tag
                    :type="row.type === 'ClusterIP' ? 'success' : row.type === 'NodePort' ? 'warning' : 'primary'"
                    size="small"
                  >
                    {{ row.type }}
                  </el-tag>
                </template>
              </el-table-column>

              <el-table-column prop="clusterIP" label="集群 IP" width="140" />

              <el-table-column label="端口" min-width="200">
                <template #default="{ row }">
                  <div class="ports-info">
                    <el-tag
                      v-for="(port, index) in (row.ports || [])"
                      :key="index"
                      size="small"
                      class="port-tag"
                      :type="row.type === 'NodePort' && port.nodePort ? 'warning' : 'primary'"
                    >
                      {{ port.port }}{{ port.name ? ':' + port.name : '' }}/{{ port.protocol }}
                      <span v-if="row.type === 'NodePort' && port.nodePort" class="nodeport-info">
                        →{{ port.nodePort }}
                      </span>
                    </el-tag>
                  </div>
                </template>
              </el-table-column>

              <el-table-column label="选择器" min-width="150">
                <template #default="{ row }">
                  <div class="selector-info">
                    <el-tag
                      v-for="(value, key) in (row.selector || {})"
                      :key="key"
                      size="small"
                      type="info"
                      class="selector-tag"
                    >
                      {{ key }}={{ value }}
                    </el-tag>
                  </div>
                </template>
              </el-table-column>

              <el-table-column label="创建时间" width="160">
                <template #default="{ row }">
                  {{ row.createdAt ? new Date(row.createdAt).toLocaleString() : '-' }}
                </template>
              </el-table-column>

              <el-table-column label="操作" width="190" fixed="right">
                <template #default="{ row }">
                  <div class="operation-buttons">
                    <el-button
                      type="primary"
                      :icon="View"
                      size="small"
                      @click="handleViewService(row)"
                    >详情</el-button>
                    <el-dropdown trigger="click">
                      <el-button type="warning" size="small" v-authority="['k8s:network:actions']">
                        更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                      </el-button>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item
                            :icon="Edit"
                            v-authority="['k8s:network:editservice']"
                            @click="handleEditService(row)"
                          >编辑 Service</el-dropdown-item>
                          <el-dropdown-item
                            :icon="Document"
                            v-authority="['k8s:network:edit_service_yaml']"
                            @click="handleEditServiceYaml(row)"
                          >编辑 YAML</el-dropdown-item>
                          <el-dropdown-item
                            :icon="Monitor"
                            v-authority="['k8s:network:service_event']"
                            @click="handleViewServiceEvents(row)"
                          >查看事件</el-dropdown-item>
                          <el-dropdown-item
                            :icon="Delete"
                            v-authority="['k8s:network:deleteservice']"
                            @click="handleDeleteService(row)"
                          >删除</el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </template>
              </el-table-column>
            </el-table>
        </div>

        <!-- Ingress 表格 -->
        <div v-if="activeTab === 'ingress'" class="tab-content">

          <el-table
            :data="filteredIngressList"
            v-loading="ingressLoading"
            element-loading-text="加载中..."
            class="resource-table"
            empty-text="暂无 Ingress 资源"
          >
              <el-table-column prop="name" label="名称" min-width="200">
                <template #default="{ row }">
                  <div class="resource-name">
                    <el-icon class="resource-icon"><Monitor /></el-icon>
                    <span class="resource-name-link" @click="handleViewIngress(row)">{{ row.name }}</span>
                    <el-tooltip content="复制名称" placement="top">
                      <el-button
                        :icon="CopyDocument"
                        size="small"
                        type="text"
                        class="copy-button"
                        @click.stop="copyToClipboard(row.name, 'Ingress')"
                      />
                    </el-tooltip>
                  </div>
                </template>
              </el-table-column>

              <el-table-column label="类型" width="120">
                <template #default="{ row }">
                  <el-tag type="primary" size="small">
                    {{ row.type || '-' }}
                  </el-tag>
                </template>
              </el-table-column>

              <el-table-column label="访问 IP" width="140">
                <template #default="{ row }">
                  <span v-if="row.loadBalancer && row.loadBalancer.ingress && row.loadBalancer.ingress.length > 0">
                    {{ row.loadBalancer.ingress[0].ip || row.loadBalancer.ingress[0].hostname || '-' }}
                  </span>
                  <span v-else>-</span>
                </template>
              </el-table-column>

              <el-table-column label="Service:端口" min-width="200">
                <template #default="{ row }">
                  <div class="services-info">
                    <template v-for="(rule, ruleIndex) in (row.rules || [])" :key="ruleIndex">
                      <el-tag
                        v-for="(path, pathIndex) in (rule.http?.paths || [])"
                        :key="`${ruleIndex}-${pathIndex}`"
                        size="small"
                        type="success"
                        class="service-tag"
                      >
                        {{ path.backend?.service?.name || path.backend?.serviceName || '-' }}:{{ path.backend?.service?.port?.number || path.backend?.servicePort || '-' }}
                      </el-tag>
                    </template>
                  </div>
                </template>
              </el-table-column>

              <el-table-column label="Ingress 域名" min-width="250">
                <template #default="{ row }">
                  <div class="endpoints-info">
                    <template v-if="row.endpoints && row.endpoints.length > 0">
                      <div
                        v-for="(endpoint, index) in row.endpoints"
                        :key="index"
                        class="endpoint-item"
                      >
                        <a :href="endpoint" target="_blank" class="endpoint-link">
                          {{ endpoint }}
                        </a>
                        <el-tooltip content="复制链接" placement="top">
                          <el-button
                            :icon="CopyDocument"
                            size="small"
                            type="text"
                            class="copy-button"
                            @click.stop="copyToClipboard(endpoint, '域名')"
                          />
                        </el-tooltip>
                      </div>
                    </template>
                    <template v-else>
                      <el-tag
                        v-for="(rule, index) in (row.rules || [])"
                        :key="index"
                        size="small"
                        type="info"
                        class="host-tag"
                      >
                        {{ rule.host || '*' }}
                      </el-tag>
                    </template>
                  </div>
                </template>
              </el-table-column>

              <el-table-column label="TLS" width="80">
                <template #default="{ row }">
                  <el-tag
                    v-if="row.tls && row.tls.length > 0"
                    type="success"
                    size="small"
                  >
                    启用
                  </el-tag>
                  <el-tag
                    v-else
                    type="info"
                    size="small"
                  >
                    未启用
                  </el-tag>
                </template>
              </el-table-column>

              <el-table-column label="创建时间" width="160">
                <template #default="{ row }">
                  {{ row.createdAt ? new Date(row.createdAt).toLocaleString() : '-' }}
                </template>
              </el-table-column>

              <el-table-column label="操作" width="190" fixed="right">
                <template #default="{ row }">
                  <div class="operation-buttons">
                    <el-button
                      type="primary"
                      :icon="View"
                      size="small"
                      @click="handleViewIngress(row)"
                    >详情</el-button>
                    <el-dropdown trigger="click">
                      <el-button type="warning" size="small" v-authority="['k8s:network:actions']" >
                        更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                      </el-button>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item
                            :icon="Edit"
                            v-authority="['k8s:network:editingress']"
                            @click="handleEditIngress(row)"
                          >编辑 Ingress</el-dropdown-item>
                          <el-dropdown-item
                            :icon="Document"
                            v-authority="['k8s:network:edit_ingress_yaml']"
                            @click="handleEditIngressYaml(row)"
                          >编辑 YAML</el-dropdown-item>
                          <el-dropdown-item
                            :icon="Link"
                            @click="handleTestIngress(row)"
                          >测试后端服务</el-dropdown-item>
                          <el-dropdown-item
                            :icon="Delete"
                            v-authority="['k8s:network:delete_ingress']"
                            @click="handleDeleteIngress(row)"
                          >删除</el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </template>
              </el-table-column>
            </el-table>
        </div>
      </div>
    </el-card>

    <!-- Service 详情对话框 -->
    <el-dialog
      v-model="serviceDetailDialogVisible"
      title="Service 详情"
      width="800px"
      class="detail-dialog"
    >
      <div class="detail-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="名称">{{ currentServiceForDetail.name }}</el-descriptions-item>
          <el-descriptions-item label="命名空间">{{ currentServiceForDetail.namespace }}</el-descriptions-item>
          <el-descriptions-item label="类型">
            <el-tag :type="currentServiceForDetail.type === 'ClusterIP' ? 'success' : 'primary'">
              {{ currentServiceForDetail.type }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="集群IP">{{ currentServiceForDetail.clusterIP }}</el-descriptions-item>
          <el-descriptions-item label="创建时间" :span="2">
            {{ currentServiceForDetail.createdAt ? new Date(currentServiceForDetail.createdAt).toLocaleString() : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="标签" :span="2">
            <div class="labels-info">
              <el-tag
                v-for="(value, key) in (currentServiceForDetail.labels || {})"
                :key="key"
                size="small"
                type="info"
                style="margin-right: 4px;"
              >
                {{ key }}={{ value }}
              </el-tag>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="端口" :span="2">
            <div class="ports-info">
              <el-tag
                v-for="(port, index) in (currentServiceForDetail.ports || [])"
                :key="index"
                size="small"
                style="margin-right: 4px;"
                :type="currentServiceForDetail.type === 'NodePort' && port.nodePort ? 'warning' : 'primary'"
              >
                {{ port.port }}{{ port.name ? ':' + port.name : '' }}/{{ port.protocol }}
                {{ port.targetPort ? ' → ' + port.targetPort : '' }}
                <span v-if="currentServiceForDetail.type === 'NodePort' && port.nodePort" class="nodeport-info">
                  (NodePort: {{ port.nodePort }})
                </span>
              </el-tag>
            </div>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- Service YAML 编辑对话框 -->
    <PodYamlDialog
      v-model:visible="serviceYamlDialogVisible"
      :yaml-content="currentServiceYaml"
      :resource-name="currentServiceForDetail?.name || ''"
      resource-type="Service"
      :editable="true"
      @save="handleServiceYamlSave"
      @close="serviceYamlDialogVisible = false"
    />

    <!-- Service 事件对话框 -->
    <el-dialog
      v-model="serviceEventsDialogVisible"
      title="Service 事件"
      width="1000px"
      class="events-dialog"
    >
      <div class="events-content">
        <el-table :data="serviceEvents" stripe style="width: 100%">
          <el-table-column prop="type" label="类型" width="80">
            <template #default="{ row }">
              <el-tag :type="row.type === 'Warning' ? 'warning' : 'success'" size="small">
                {{ row.type }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="reason" label="原因" width="120" />
          <el-table-column prop="message" label="消息" min-width="300" />
          <el-table-column prop="firstTimestamp" label="首次时间" width="160">
            <template #default="{ row }">
              {{ row.firstTimestamp ? new Date(row.firstTimestamp).toLocaleString() : '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="count" label="次数" width="80" />
        </el-table>
      </div>
    </el-dialog>

    <!-- Ingress 详情对话框 -->
    <el-dialog
      v-model="ingressDetailDialogVisible"
      title="Ingress 详情"
      width="800px"
      class="detail-dialog"
    >
      <div class="detail-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="名称">{{ currentIngressForDetail.name }}</el-descriptions-item>
          <el-descriptions-item label="命名空间">{{ currentIngressForDetail.namespace }}</el-descriptions-item>
          <el-descriptions-item label="TLS" :span="2">
            <el-tag :type="currentIngressForDetail.tls && currentIngressForDetail.tls.length > 0 ? 'success' : 'info'" size="small">
              {{ currentIngressForDetail.tls && currentIngressForDetail.tls.length > 0 ? '启用' : '未启用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间" :span="2">
            {{ currentIngressForDetail.createdAt ? new Date(currentIngressForDetail.createdAt).toLocaleString() : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="规则" :span="2">
            <div class="rules-info">
              <div v-for="(rule, index) in (currentIngressForDetail.rules || [])" :key="index" style="margin-bottom: 8px;">
                <strong>{{ rule.host || '*' }}</strong>
                <div v-for="(path, pathIndex) in (rule.paths || [])" :key="pathIndex" style="margin-left: 16px;">
                  <el-tag size="small" type="info">
                    {{ path.path }} → {{ path.serviceName }}:{{ path.servicePort }}
                  </el-tag>
                </div>
              </div>
            </div>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- Ingress YAML 编辑对话框 -->
    <PodYamlDialog
      v-model:visible="ingressYamlDialogVisible"
      :yaml-content="currentIngressYaml"
      :resource-name="currentIngressForDetail?.name || ''"
      resource-type="Ingress"
      :editable="true"
      @save="handleIngressYamlSave"
      @close="ingressYamlDialogVisible = false"
    />

    <!-- Ingress 事件对话框 -->
    <el-dialog
      v-model="ingressEventsDialogVisible"
      title="Ingress 事件"
      width="1000px"
      class="events-dialog"
    >
      <div class="events-content">
        <el-table :data="ingressEvents" stripe style="width: 100%">
          <el-table-column prop="type" label="类型" width="80">
            <template #default="{ row }">
              <el-tag :type="row.type === 'Warning' ? 'warning' : 'success'" size="small">
                {{ row.type }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="reason" label="原因" width="120" />
          <el-table-column prop="message" label="消息" min-width="300" />
          <el-table-column prop="firstTimestamp" label="首次时间" width="160">
            <template #default="{ row }">
              {{ row.firstTimestamp ? new Date(row.firstTimestamp).toLocaleString() : '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="count" label="次数" width="80" />
        </el-table>
      </div>
    </el-dialog>

    <!-- 创建 Service 对话框 -->
    <el-dialog
      v-model="serviceDialogVisible"
      title="创建 Service"
      width="700px"
      class="create-dialog compact-dialog"
    >
      <el-form :model="currentService" label-width="90px" class="compact-form">
        <div class="form-row">
          <el-form-item label="名称" required class="form-item-half">
            <el-input
              v-model="currentService.name"
              placeholder="请输入 Service 名称"
              size="small"
            />
          </el-form-item>
          <el-form-item label="类型" required class="form-item-half">
            <el-select v-model="currentService.type" style="width: 100%" size="small">
              <el-option
                v-for="option in serviceTypeOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              >
                <div class="option-content">
                  <span>{{ option.label }}</span>
                  <span class="option-desc">{{ option.description }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="端口配置" required>
          <div class="ports-config">
            <div
              v-for="(port, index) in currentService.ports"
              :key="index"
              class="port-item"
            >
              <div class="port-fields">
                <el-input
                  v-model="port.name"
                  placeholder="端口d(可选)"
                  size="small"
                  style="width: 120px;"
                />
                <el-input-number
                  v-model="port.port"
                  :min="1"
                  :max="65535"
                  placeholder="端口"
                  size="small"
                  style="width: 100px;"
                />
                <el-input-number
                  v-model="port.targetPort"
                  :min="1"
                  :max="65535"
                  placeholder="目标端口"
                  size="small"
                  style="width: 100px;"
                />
                <el-select v-model="port.protocol" size="small" style="width: 80px;">
                  <el-option
                    v-for="protocol in protocolOptions"
                    :key="protocol.value"
                    :label="protocol.label"
                    :value="protocol.value"
                  />
                </el-select>
                <el-button
                  type="danger"
                  :icon="Delete"
                  size="small"
                  plain
                  @click="removeServicePort(index)"
                  :disabled="currentService.ports.length === 1"
                />
              </div>
            </div>
            <el-button type="primary" :icon="Plus" size="small" @click="addServicePort">
              添加端口
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="currentService.description"
            type="textarea"
            :rows="2"
            placeholder="请输入服务描述"
            size="small"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="serviceDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitServiceForm">创建</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑 Service 对话框 -->
    <el-dialog
      v-model="editServiceDialogVisible"
      title="编辑 Service"
      width="700px"
      class="edit-dialog compact-dialog"
    >
      <el-form :model="currentService" label-width="90px" class="compact-form">
        <div class="form-row">
          <el-form-item label="名称" class="form-item-half">
            <el-input
              v-model="currentService.name"
              disabled
              size="small"
            />
          </el-form-item>
          <el-form-item label="类型" required class="form-item-half">
            <el-select v-model="currentService.type" style="width: 100%" size="small">
              <el-option
                v-for="option in serviceTypeOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              >
                <div class="option-content">
                  <span>{{ option.label }}</span>
                  <span class="option-desc">{{ option.description }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="端口配置" required>
          <div class="ports-config">
            <div
              v-for="(port, index) in currentService.ports"
              :key="index"
              class="port-item"
            >
              <div class="port-fields">
                <el-input
                  v-model="port.name"
                  placeholder="端口名称(可选)"
                  size="small"
                  style="width: 120px;"
                />
                <el-input-number
                  v-model="port.port"
                  :min="1"
                  :max="65535"
                  placeholder="端口"
                  size="small"
                  style="width: 100px;"
                />
                <el-input-number
                  v-model="port.targetPort"
                  :min="1"
                  :max="65535"
                  placeholder="目标端口"
                  size="small"
                  style="width: 100px;"
                />
                <el-select v-model="port.protocol" size="small" style="width: 80px;">
                  <el-option
                    v-for="protocol in protocolOptions"
                    :key="protocol.value"
                    :label="protocol.label"
                    :value="protocol.value"
                  />
                </el-select>
                <el-button
                  type="danger"
                  :icon="Delete"
                  size="small"
                  plain
                  @click="removeServicePort(index)"
                  :disabled="currentService.ports.length === 1"
                />
              </div>
            </div>
            <el-button type="primary" :icon="Plus" size="small" @click="addServicePort">
              添加端口
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="currentService.description"
            type="textarea"
            :rows="2"
            placeholder="请输入服务描述"
            size="small"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editServiceDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitEditServiceForm">更新</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 创建 Ingress 对话框 -->
    <el-dialog
      v-model="ingressDialogVisible"
      title="创建 Ingress"
      width="800px"
      class="create-dialog compact-dialog"
    >
      <el-form :model="currentIngress" label-width="90px" class="compact-form">
        <el-form-item label="名称" required>
          <el-input
            v-model="currentIngress.name"
            placeholder="请输入 Ingress 名称"
            size="small"
            style="width: 300px;"
          />
        </el-form-item>

        <el-form-item label="规则配置" required>
          <div class="rules-config">
            <div
              v-for="(rule, ruleIndex) in currentIngress.rules"
              :key="ruleIndex"
              class="rule-item"
            >
              <div class="rule-header">
                <el-input
                  v-model="rule.host"
                  placeholder="主机名 (如: example.com)"
                  size="small"
                  style="width: 250px;"
                />
                <el-button
                  type="danger"
                  :icon="Delete"
                  size="small"
                  plain
                  @click="removeIngressRule(ruleIndex)"
                  :disabled="currentIngress.rules.length === 1"
                >
                  删除规则
                </el-button>
              </div>

              <div class="paths-config">
                <div
                  v-for="(path, pathIndex) in rule.paths"
                  :key="pathIndex"
                  class="path-item"
                >
                  <div class="path-fields">
                    <el-input
                      v-model="path.path"
                      placeholder="路径"
                      size="small"
                      style="width: 120px;"
                    />
                    <el-select v-model="path.pathType" size="small" style="width: 150px;">
                      <el-option
                        v-for="pathType in pathTypeOptions"
                        :key="pathType.value"
                        :label="pathType.label"
                        :value="pathType.value"
                      >
                        <div class="option-content">
                          <span>{{ pathType.label }}</span>
                          <span class="option-desc">{{ pathType.description }}</span>
                        </div>
                      </el-option>
                    </el-select>
                    <el-input
                      v-model="path.serviceName"
                      placeholder="服务名"
                      size="small"
                      style="width: 120px;"
                    />
                    <el-input-number
                      v-model="path.servicePort"
                      :min="1"
                      :max="65535"
                      placeholder="端口"
                      size="small"
                      style="width: 100px;"
                    />
                    <el-button
                      type="danger"
                      :icon="Delete"
                      size="small"
                      plain
                      @click="removeIngressPath(ruleIndex, pathIndex)"
                      :disabled="rule.paths.length === 1"
                    />
                  </div>
                </div>
                <el-button
                  type="primary"
                  :icon="Plus"
                  size="small"
                  @click="addIngressPath(ruleIndex)"
                >
                  添加路径
                </el-button>
              </div>
            </div>
            <el-button type="success" :icon="Plus" size="small" @click="addIngressRule">
              添加规则
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="currentIngress.description"
            type="textarea"
            :rows="2"
            placeholder="请输入 Ingress 描述"
            size="small"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="ingressDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitIngressForm">创建</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑 Ingress 对话框 -->
    <el-dialog
      v-model="editIngressDialogVisible"
      title="编辑 Ingress"
      width="800px"
      class="edit-dialog compact-dialog"
    >
      <el-form :model="currentIngress" label-width="90px" class="compact-form">
        <el-form-item label="名称">
          <el-input
            v-model="currentIngress.name"
            disabled
            size="small"
            style="width: 300px;"
          />
        </el-form-item>

        <el-form-item label="规则配置" required>
          <div class="rules-config">
            <div
              v-for="(rule, ruleIndex) in currentIngress.rules"
              :key="ruleIndex"
              class="rule-item"
            >
              <div class="rule-header">
                <el-input
                  v-model="rule.host"
                  placeholder="主机名 (如: example.com)"
                  size="small"
                  style="width: 250px;"
                />
                <el-button
                  type="danger"
                  :icon="Delete"
                  size="small"
                  plain
                  @click="removeIngressRule(ruleIndex)"
                  :disabled="currentIngress.rules.length === 1"
                >
                  删除规则
                </el-button>
              </div>

              <div class="paths-config">
                <div
                  v-for="(path, pathIndex) in rule.paths"
                  :key="pathIndex"
                  class="path-item"
                >
                  <div class="path-fields">
                    <el-input
                      v-model="path.path"
                      placeholder="路径"
                      size="small"
                      style="width: 120px;"
                    />
                    <el-select v-model="path.pathType" size="small" style="width: 150px;">
                      <el-option
                        v-for="pathType in pathTypeOptions"
                        :key="pathType.value"
                        :label="pathType.label"
                        :value="pathType.value"
                      >
                        <div class="option-content">
                          <span>{{ pathType.label }}</span>
                          <span class="option-desc">{{ pathType.description }}</span>
                        </div>
                      </el-option>
                    </el-select>
                    <el-input
                      v-model="path.serviceName"
                      placeholder="服务名"
                      size="small"
                      style="width: 120px;"
                    />
                    <el-input-number
                      v-model="path.servicePort"
                      :min="1"
                      :max="65535"
                      placeholder="端口"
                      size="small"
                      style="width: 100px;"
                    />
                    <el-button
                      type="danger"
                      :icon="Delete"
                      size="small"
                      plain
                      @click="removeIngressPath(ruleIndex, pathIndex)"
                      :disabled="rule.paths.length === 1"
                    />
                  </div>
                </div>
                <el-button
                  type="primary"
                  :icon="Plus"
                  size="small"
                  @click="addIngressPath(ruleIndex)"
                >
                  添加路径
                </el-button>
              </div>
            </div>
            <el-button type="success" :icon="Plus" size="small" @click="addIngressRule">
              添加规则
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="currentIngress.description"
            type="textarea"
            :rows="2"
            placeholder="请输入 Ingress 描述"
            size="small"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editIngressDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitEditIngressForm">更新</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Ingress 测试对话框 -->
    <el-dialog
      v-model="ingressTestDialogVisible"
      title="测试 Ingress 后端服务"
      width="900px"
      class="test-dialog"
      @close="resetTestForm"
    >
      <div class="test-content">
        <el-form :model="ingressTestForm" label-width="100px" class="test-form">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="服务名称" required>
                <el-input
                  v-model="ingressTestForm.serviceName"
                  placeholder="请输入Service名称"
                  size="small"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="服务端口" required>
                <el-input-number
                  v-model="ingressTestForm.servicePort"
                  :min="1"
                  :max="65535"
                  placeholder="请输入端口"
                  size="small"
                  style="width: 100%;"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="测试路径">
                <el-input
                  v-model="ingressTestForm.path"
                  placeholder="请输入测试路径,如:/api/health"
                  size="small"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="Host域名">
                <el-input
                  v-model="ingressTestForm.host"
                  placeholder="请输入Host域名(可选)"
                  size="small"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="请求方法">
                <el-select v-model="ingressTestForm.method" size="small" style="width: 100%;">
                  <el-option label="GET" value="GET" />
                  <el-option label="POST" value="POST" />
                  <el-option label="PUT" value="PUT" />
                  <el-option label="DELETE" value="DELETE" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="超时时间(秒)">
                <el-input-number
                  v-model="ingressTestForm.timeout"
                  :min="1"
                  :max="60"
                  size="small"
                  style="width: 100%;"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item>
            <el-button
              type="primary"
              :loading="ingressTestLoading"
              @click="executeIngressTest"
            >
              {{ ingressTestLoading ? '测试中...' : '开始测试' }}
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 测试结果显示区域 -->
        <div v-if="ingressTestResult" class="test-result">
          <el-divider content-position="left">测试结果</el-divider>

          <el-alert
            :title="ingressTestResult.status || '测试完成'"
            :type="ingressTestResult.success ? 'success' : 'error'"
            :closable="false"
            show-icon
          >
            <template #default>
              <div class="result-message">
                <p><strong>{{ ingressTestResult.message }}</strong></p>
                <p v-if="ingressTestResult.suggestion" class="suggestion">
                  💡 建议: {{ ingressTestResult.suggestion }}
                </p>
              </div>
            </template>
          </el-alert>

          <el-descriptions :column="2" border class="result-details" v-if="ingressTestResult.details">
            <el-descriptions-item label="服务名称">
              {{ ingressTestResult.details.serviceName || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="服务端口">
              {{ ingressTestResult.details.servicePort || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="测试路径">
              {{ ingressTestResult.details.testPath || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="Host域名">
              {{ ingressTestResult.details.testHost || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="请求方法">
              {{ ingressTestResult.details.method || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="HTTP状态码">
              <el-tag
                :type="ingressTestResult.statusCode >= 200 && ingressTestResult.statusCode < 300 ? 'success' : 'danger'"
                size="small"
              >
                {{ ingressTestResult.statusCode || '-' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="响应时间">
              <el-tag
                :type="ingressTestResult.responseTime < 1000 ? 'success' : 'warning'"
                size="small"
              >
                {{ ingressTestResult.responseTime }}ms
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="Endpoints状态">
              {{ ingressTestResult.details.endpointsReady || 0 }} / {{ ingressTestResult.details.endpointsTotal || 0 }} 就绪
            </el-descriptions-item>
          </el-descriptions>

          <!-- 显示已测试的Pod列表 -->
          <div v-if="ingressTestResult.details?.podsTested && ingressTestResult.details.podsTested.length > 0" class="pods-tested">
            <h4>已测试的Pod:</h4>
            <el-tag
              v-for="pod in ingressTestResult.details.podsTested"
              :key="pod"
              size="small"
              type="info"
              style="margin-right: 8px; margin-bottom: 8px;"
            >
              {{ pod }}
            </el-tag>
          </div>

          <!-- 显示响应头 -->
          <div v-if="ingressTestResult.details?.responseHeaders" class="response-headers">
            <el-divider content-position="left">响应头</el-divider>
            <el-scrollbar max-height="200px">
              <pre class="code-block">{{ JSON.stringify(ingressTestResult.details.responseHeaders, null, 2) }}</pre>
            </el-scrollbar>
          </div>

          <!-- 显示响应体 -->
          <div v-if="ingressTestResult.details?.responseBody" class="response-body">
            <el-divider content-position="left">响应体</el-divider>
            <el-scrollbar max-height="300px">
              <pre class="code-block">{{ ingressTestResult.details.responseBody }}</pre>
            </el-scrollbar>
          </div>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="ingressTestDialogVisible = false">关闭</el-button>
          <el-button type="primary" :loading="ingressTestLoading" @click="executeIngressTest">
            重新测试
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.k8s-network-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.network-card {
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
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.search-section {
  margin-bottom: 20px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form {
  margin: 0;
}

.resource-controls {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid rgba(103, 126, 234, 0.1);
}

.network-tabs {
  flex: 1;
  margin: 0;
}

.network-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
}

.network-tabs :deep(.el-tabs__item) {
  font-weight: 500;
  color: #606266;
}

.network-tabs :deep(.el-tabs__item.is-active) {
  color: #409EFF;
  font-weight: 600;
}

.create-actions {
  flex-shrink: 0;
  margin-top: 8px;
}

.cluster-option,
.namespace-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.cluster-name {
  font-weight: 500;
  color: #6b9bd1;
}

.cluster-status {
  margin-left: 8px;
}

.namespace-name {
  font-weight: 500;
  color: #6b9bd1;
}

.namespace-status-tag {
  margin-left: 8px;
}

.selector-section .el-select {
  border-radius: 8px;
}

.selector-section .el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(103, 126, 234, 0.3);
  border-radius: 8px;
  transition: all 0.3s ease;
}

.selector-section .el-select :deep(.el-input__wrapper):hover {
  border-color: #667eea;
}

.selector-section .el-select :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
}

.header-actions {
  display: flex;
  align-items: center;
}

.tab-content {
  padding: 0;
}


.resource-table {
  border-radius: 8px;
  overflow: hidden;
}

.resource-table :deep(.el-table__header) {
  background: rgba(103, 126, 234, 0.1);
}

.resource-table :deep(.el-table__header th) {
  background: rgba(103, 126, 234, 0.1);
  color: #6b9bd1;
  font-weight: 600;
  border: none;
}

.resource-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.resource-table :deep(.el-table__row:hover) {
  background: rgba(103, 126, 234, 0.05);
}

.resource-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.resource-icon {
  color: #667eea;
  font-size: 16px;
}

.resource-name-link {
  color: #667eea;
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.3s ease;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  flex: 1;
}

.resource-name-link:hover {
  color: #409eff;
  border-bottom-color: #409eff;
  transform: translateY(-1px);
}

.copy-button {
  color: #909399;
  padding: 2px 4px;
  margin-left: 4px;
  opacity: 0;
  transition: all 0.3s ease;
  min-width: auto;
  height: auto;
}

.resource-name:hover .copy-button {
  opacity: 1;
}

.copy-button:hover {
  color: #409eff;
  background-color: rgba(64, 158, 255, 0.1);
  transform: scale(1.1);
}

.ports-info,
.hosts-info,
.paths-info,
.selector-info,
.services-info {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.endpoints-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.port-tag,
.host-tag,
.path-tag,
.selector-tag,
.service-tag {
  font-size: 11px;
  border-radius: 4px;
  margin: 1px;
}

.endpoint-item {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 4px;
}

.endpoint-link {
  color: #409eff;
  text-decoration: none;
  font-size: 12px;
  transition: all 0.3s ease;
  border-bottom: 1px solid transparent;
}

.endpoint-link:hover {
  color: #66b1ff;
  border-bottom-color: #66b1ff;
}

.nodeport-info {
  font-weight: 600;
  color: #E6A23C;
  margin-left: 2px;
}

/* 操作按钮样式 */
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

/* 兼容旧的 action-buttons 样式 */
.action-buttons {
  display: flex;
  gap: 4px;
  align-items: center;
}

.action-buttons .el-button {
  padding: 4px 8px;
  font-size: 11px;
}

.action-buttons .el-dropdown {
  margin-right: 4px;
}

.action-buttons .el-dropdown .el-button {
  font-size: 11px;
  padding: 4px 8px;
}

/* 详情对话框样式 */
.detail-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.detail-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.detail-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.detail-content {
  padding: 20px;
}

/* YAML 对话框样式 */
.yaml-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.yaml-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.yaml-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.yaml-content {
  padding: 20px;
}

.yaml-content .el-textarea :deep(.el-textarea__inner) {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  font-family: 'Courier New', 'Consolas', monospace;
  line-height: 1.4;
}

/* 事件对话框样式 */
.events-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.events-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.events-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.events-content {
  padding: 20px;
}

/* 对话框样式 */
.create-dialog :deep(.el-dialog),
.edit-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.create-dialog :deep(.el-dialog__header),
.edit-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.create-dialog :deep(.el-dialog__title),
.edit-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.create-dialog :deep(.el-dialog__body),
.edit-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

.compact-dialog :deep(.el-dialog__body) {
  padding: 15px 20px;
}

.compact-form {
  margin: 0;
}

.compact-form :deep(.el-form-item) {
  margin-bottom: 12px;
}

.form-row {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.form-item-half {
  flex: 1;
  margin-bottom: 0 !important;
}

.option-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.option-desc {
  font-size: 11px;
  color: #909399;
  margin-left: 8px;
}

/* 端口配置样式 */
.ports-config {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 12px;
  background: rgba(249, 250, 251, 0.5);
}

.port-item {
  margin-bottom: 8px;
  padding: 8px;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  background: white;
}

.port-fields {
  display: flex;
  gap: 8px;
  align-items: center;
}

/* 规则配置样式 */
.rules-config {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 12px;
  background: rgba(249, 250, 251, 0.5);
}

.rule-item {
  margin-bottom: 16px;
  padding: 12px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: white;
}

.rule-header {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.paths-config {
  padding-left: 12px;
}

.path-item {
  margin-bottom: 8px;
  padding: 8px;
  border: 1px solid #f0f0f0;
  border-radius: 4px;
  background: rgba(249, 250, 251, 0.3);
}

.path-fields {
  display: flex;
  gap: 8px;
  align-items: center;
}

.dialog-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
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

.el-textarea :deep(.el-textarea__inner) {
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.2);
  transition: all 0.3s ease;
}

.el-textarea :deep(.el-textarea__inner):focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .k8s-network-management {
    padding: 10px;
  }

  .card-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .header-actions {
    justify-content: flex-end;
    flex-wrap: wrap;
  }

  .header-actions .el-select {
    width: 100% !important;
    margin-bottom: 8px;
  }

  .search-section {
    padding: 12px;
  }

  .search-form {
    flex-direction: column;
    align-items: stretch;
  }

  .search-form :deep(.el-form-item) {
    margin-bottom: 8px;
  }

  .resource-controls {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .network-tabs,
  .create-actions {
    width: 100%;
  }

  .create-actions {
    justify-content: center;
    margin-top: 12px;
  }


  .form-row {
    flex-direction: column;
    gap: 8px;
  }

  .port-fields,
  .path-fields {
    flex-direction: column;
    gap: 8px;
  }

  .rule-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .operation-buttons {
    gap: 4px;
  }

  .operation-buttons .el-button {
    width: 32px;
    height: 32px;
  }

  .action-buttons {
    flex-direction: column;
    gap: 2px;
  }

  .action-buttons .el-button {
    width: 100%;
    margin: 0;
  }

  .action-buttons .el-dropdown {
    width: 100%;
    margin-right: 0;
    margin-bottom: 2px;
  }

  .action-buttons .el-dropdown .el-button {
    width: 100%;
  }

  .resource-name-link {
    font-size: 13px;
    line-height: 1.2;
  }
}

/* 测试对话框样式 */
.test-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.test-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.test-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.test-content {
  padding: 20px;
  max-height: 70vh;
  overflow-y: auto;
}

.test-form {
  margin-bottom: 20px;
}

.test-result {
  margin-top: 20px;
  padding: 20px;
  background: rgba(249, 250, 251, 0.5);
  border-radius: 12px;
  border: 1px solid #e4e7ed;
}

.result-message {
  margin: 0;
}

.result-message p {
  margin: 8px 0;
}

.suggestion {
  color: #409eff;
  font-size: 14px;
  margin-top: 12px;
}

.result-details {
  margin-top: 16px;
}

.pods-tested {
  margin-top: 16px;
}

.pods-tested h4 {
  margin-bottom: 12px;
  font-size: 14px;
  color: #606266;
}

.response-headers,
.response-body {
  margin-top: 16px;
}

.code-block {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 12px;
  font-family: 'Courier New', 'Consolas', monospace;
  font-size: 12px;
  line-height: 1.5;
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: #333;
}

</style>
