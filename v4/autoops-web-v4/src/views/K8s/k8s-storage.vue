<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Edit,
  Delete,
  View,
  ArrowDown,
  Document,
  Refresh,
  Connection,
  DataLine,
  Cpu,
  Search,
  Monitor,
  Files
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import PodYamlDialog from './pods/PodYamlDialog.vue'
import ClusterSelector from './pods/ClusterSelector.vue'
import NamespaceSelector from './pods/NamespaceSelector.vue'
import { getYamlContentFromApiResponse } from './utils/yaml'

// 基础状态
const loading = ref(false)
const activeTab = ref('pvc')
const searchKeyword = ref('')

// 集群和命名空间状态
const selectedClusterId = ref('')
const queryParams = reactive({
  namespace: 'default'
})

// 存储资源数据状态
const pvcList = ref([])
const pvList = ref([])
const storageClassList = ref([])

// 对话框状态
const pvcYamlDialogVisible = ref(false)
const pvYamlDialogVisible = ref(false)
const storageClassYamlDialogVisible = ref(false)

// 详情对话框状态
const pvcDetailDialogVisible = ref(false)
const pvDetailDialogVisible = ref(false)
const storageClassDetailDialogVisible = ref(false)

// 当前操作的资源
const currentPVCYaml = ref('')
const currentPVYaml = ref('')
const currentStorageClassYaml = ref('')
const currentResourceName = ref('')
const currentResourceType = ref('')

// 当前查看的资源详情
const currentPVCForDetail = ref({})
const currentPVForDetail = ref({})
const currentStorageClassForDetail = ref({})

// 过滤后的列表
const filteredPVCList = computed(() => {
  const list = Array.isArray(pvcList.value) ? pvcList.value : []
  if (!searchKeyword.value) return list
  return list.filter(item =>
    item.name?.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    item.status?.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    item.storageClassName?.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

const filteredPVList = computed(() => {
  const list = Array.isArray(pvList.value) ? pvList.value : []
  if (!searchKeyword.value) return list
  return list.filter(item =>
    item.name?.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    item.status?.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    item.storageClassName?.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

const filteredStorageClassList = computed(() => {
  const list = Array.isArray(storageClassList.value) ? storageClassList.value : []
  if (!searchKeyword.value) return list
  return list.filter(item =>
    item.name?.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    item.provisioner?.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

// 处理集群选择变化
const handleClusterChange = (clusterId) => {
  selectedClusterId.value = clusterId
  console.log('集群选择变化:', clusterId)
  if (clusterId && queryParams.namespace) {
    loadAllStorageResources()
  }
}

// 处理命名空间选择变化
const handleNamespaceChange = (namespace) => {
  queryParams.namespace = namespace
  console.log('命名空间选择变化:', namespace)
  if (selectedClusterId.value && namespace) {
    loadAllStorageResources()
  }
}

// 获取命名空间列表
const fetchNamespaceList = async () => {
  if (!selectedClusterId.value) return

  try {
    const response = await k8sApi.getNamespaceList(selectedClusterId.value)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      namespaceList.value = responseData.data || []
      // 如果当前选中的命名空间不在新列表中，重置为default
      if (!namespaceList.value.some(ns => ns.name === queryParams.namespace)) {
        queryParams.namespace = namespaceList.value.length > 0 ? namespaceList.value[0].name : 'default'
      }
    } else {
      ElMessage.error(responseData.message || '获取命名空间列表失败')
    }
  } catch (error) {
    console.error('获取命名空间列表失败:', error)
    ElMessage.error('获取命名空间列表失败')
  }
}

// 加载所有存储资源
const loadAllStorageResources = async () => {
  if (!selectedClusterId.value) {
    console.warn('集群ID为空，无法加载存储资源')
    return
  }

  console.log('开始加载存储资源，集群ID:', selectedClusterId.value, '命名空间:', queryParams.namespace)

  loading.value = true
  try {
    // 并发加载所有存储资源
    await Promise.all([
      fetchPVCList(),
      fetchPVList(),
      fetchStorageClassList()
    ])
  } catch (error) {
    console.error('加载存储资源失败:', error)
  } finally {
    loading.value = false
  }
}

// 重置搜索
const resetSearch = () => {
  searchKeyword.value = ''
}


// 标签页切换处理
const handleTabChange = (tabName) => {
  console.log('标签页切换到:', tabName)
  activeTab.value = tabName
}

// 获取PVC列表
const fetchPVCList = async () => {
  if (!selectedClusterId.value || !queryParams.namespace) {
    console.warn('集群ID或命名空间为空，无法获取 PVC 列表')
    return
  }

  try {
    console.log('正在获取 PVC 列表，集群ID:', selectedClusterId.value, '命名空间:', queryParams.namespace)

    const response = await k8sApi.getPVCList(selectedClusterId.value, queryParams.namespace)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      const rawData = responseData.data?.pvcs || responseData.data
      pvcList.value = Array.isArray(rawData) ? rawData : []
      console.log('获取到 PVC 列表:', pvcList.value.length, '个')
    } else {
      ElMessage.error(responseData.message || '获取 PVC 列表失败')
      pvcList.value = []
    }
  } catch (error) {
    console.error('获取 PVC 列表失败:', error)
    ElMessage.error('获取 PVC 列表失败，请检查网络连接')
    pvcList.value = []
  }
}

// 获取PV列表
const fetchPVList = async () => {
  if (!selectedClusterId.value) {
    console.warn('集群ID为空，无法获取 PV 列表')
    return
  }

  try {
    console.log('正在获取 PV 列表，集群ID:', selectedClusterId.value)

    const response = await k8sApi.getPVList(selectedClusterId.value)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      const rawData = responseData.data?.pvs || responseData.data
      pvList.value = Array.isArray(rawData) ? rawData : []
      console.log('获取到 PV 列表:', pvList.value.length, '个')
    } else {
      ElMessage.error(responseData.message || '获取 PV 列表失败')
      pvList.value = []
    }
  } catch (error) {
    console.error('获取 PV 列表失败:', error)
    ElMessage.error('获取 PV 列表失败，请检查网络连接')
    pvList.value = []
  }
}

// 获取StorageClass列表
const fetchStorageClassList = async () => {
  if (!selectedClusterId.value) {
    console.warn('集群ID为空，无法获取 StorageClass 列表')
    return
  }

  try {
    console.log('正在获取 StorageClass 列表，集群ID:', selectedClusterId.value)

    const response = await k8sApi.getStorageClassList(selectedClusterId.value)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      const rawData = responseData.data?.storageClasses || responseData.data
      storageClassList.value = Array.isArray(rawData) ? rawData : []
      console.log('获取到 StorageClass 列表:', storageClassList.value.length, '个')
    } else {
      ElMessage.error(responseData.message || '获取 StorageClass 列表失败')
      storageClassList.value = []
    }
  } catch (error) {
    console.error('获取 StorageClass 列表失败:', error)
    ElMessage.error('获取 StorageClass 列表失败，请检查网络连接')
    storageClassList.value = []
  }
}

// 刷新数据
const handleRefresh = () => {
  loadAllStorageResources()
}

// PVC 操作
const handleCreatePVC = () => {
  if (!selectedClusterId.value) {
    ElMessage.warning('请先选择集群')
    return
  }
  if (!queryParams.namespace) {
    ElMessage.warning('请先选择命名空间')
    return
  }

  // 设置默认的PVC YAML模板
  const defaultYaml = `apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: new-pvc
  namespace: ${queryParams.namespace}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: standard`

  currentPVCYaml.value = defaultYaml
  currentPVCForDetail.value = { name: 'new-pvc', namespace: queryParams.namespace }
  pvcYamlDialogVisible.value = true
}

// 查看 PVC 详情
const handleViewPVC = async (row) => {
  try {
    console.log('🔍 获取PVC详情:', {
      clusterId: selectedClusterId.value,
      namespace: queryParams.namespace,
      pvcName: row?.name,
      row: row
    })

    const response = await k8sApi.getPVCDetail(selectedClusterId.value, queryParams.namespace, row?.name)
    const responseData = response.data || response

    console.log('📥 PVC详情API响应:', responseData)

    if (responseData.code === 200 || responseData.success) {
      currentPVCForDetail.value = responseData.data || row
      console.log('✅ PVC详情数据已设置，打开对话框')
      pvcDetailDialogVisible.value = true
    } else {
      console.error('❌ API返回失败:', responseData)
      ElMessage.error(responseData.message || '获取 PVC 详情失败')
    }
  } catch (error) {
    console.error('❌ 获取 PVC 详情失败:', error)
    ElMessage.error('获取 PVC 详情失败: ' + (error.message || '未知错误'))
  }
}

// 编辑 PVC YAML
const handleEditPVCYaml = async (row) => {
  try {
    loading.value = true
    const response = await k8sApi.getPVCYaml(selectedClusterId.value, queryParams.namespace, row?.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      currentPVCYaml.value = getYamlContentFromApiResponse(
        responseData,
        `# PVC ${row?.name} YAML\napiVersion: v1\nkind: PersistentVolumeClaim\nmetadata:\n  name: ${row?.name}`
      )
      currentPVCForDetail.value = row
      pvcYamlDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取 PVC YAML 失败')
    }
  } catch (error) {
    console.error('获取 PVC YAML 失败:', error)
    ElMessage.error('获取 PVC YAML 失败')
  } finally {
    loading.value = false
  }
}

// 保存 PVC YAML
const handlePVCYamlSave = async (data) => {
  try {
    // 检查是否是创建新PVC (new-pvc表示新建)
    const isCreating = data.resourceName === 'new-pvc' || !currentPVCForDetail.value?.creationTimestamp

    let response
    if (isCreating) {
      // 创建新PVC，使用createPodFromYaml API（支持多种资源类型）
      response = await k8sApi.createPodFromYaml(selectedClusterId.value, queryParams.namespace, { yamlContent: data.yamlContent })
    } else {
      // 更新现有PVC
      response = await k8sApi.updatePVCYaml(selectedClusterId.value, queryParams.namespace, data.resourceName, data.yamlContent)
    }

    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(isCreating ? 'PVC 创建成功' : 'PVC YAML 更新成功')
      pvcYamlDialogVisible.value = false
      fetchPVCList() // 刷新列表
    } else {
      ElMessage.error(responseData.message || (isCreating ? 'PVC 创建失败' : 'PVC YAML 更新失败'))
    }
  } catch (error) {
    console.error('PVC 操作失败:', error)
    ElMessage.error('PVC 操作失败')
  }
}

// 删除 PVC
const handleDeletePVC = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认删除 PVC "${row?.name}"？删除后数据将无法恢复。`,
      '删除确认',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    const response = await k8sApi.deletePVC(selectedClusterId.value, queryParams.namespace, row?.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('PVC 删除成功')
      fetchPVCList()
    } else {
      ElMessage.error(responseData.message || 'PVC 删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除 PVC 失败:', error)
      ElMessage.error('删除 PVC 失败')
    }
  }
}

// PV 操作
const handleCreatePV = () => {
  if (!selectedClusterId.value) {
    ElMessage.warning('请先选择集群')
    return
  }

  // 设置默认的PV YAML模板
  const defaultYaml = `apiVersion: v1
kind: PersistentVolume
metadata:
  name: new-pv
spec:
  capacity:
    storage: 20Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: standard
  hostPath:
    path: /data/storage`

  currentPVYaml.value = defaultYaml
  currentPVForDetail.value = { name: 'new-pv' }
  pvYamlDialogVisible.value = true
}

// 查看 PV 详情
const handleViewPV = async (row) => {
  try {
    console.log('🔍 获取PV详情:', {
      clusterId: selectedClusterId.value,
      pvName: row?.name,
      row: row
    })

    const response = await k8sApi.getPVDetail(selectedClusterId.value, row?.name)
    const responseData = response.data || response

    console.log('📥 PV详情API响应:', responseData)

    if (responseData.code === 200 || responseData.success) {
      currentPVForDetail.value = responseData.data || row
      console.log('✅ PV详情数据已设置，打开对话框')
      pvDetailDialogVisible.value = true
    } else {
      console.error('❌ API返回失败:', responseData)
      ElMessage.error(responseData.message || '获取 PV 详情失败')
    }
  } catch (error) {
    console.error('❌ 获取 PV 详情失败:', error)
    ElMessage.error('获取 PV 详情失败: ' + (error.message || '未知错误'))
  }
}

// 编辑 PV YAML
const handleEditPVYaml = async (row) => {
  try {
    loading.value = true
    const response = await k8sApi.getPVYaml(selectedClusterId.value, row?.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      currentPVYaml.value = getYamlContentFromApiResponse(
        responseData,
        `# PV ${row?.name} YAML\napiVersion: v1\nkind: PersistentVolume\nmetadata:\n  name: ${row?.name}`
      )
      currentPVForDetail.value = row
      pvYamlDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取 PV YAML 失败')
    }
  } catch (error) {
    console.error('获取 PV YAML 失败:', error)
    ElMessage.error('获取 PV YAML 失败')
  } finally {
    loading.value = false
  }
}

// 保存 PV YAML
const handlePVYamlSave = async (data) => {
  try {
    // 检查是否是创建新PV (new-pv表示新建)
    const isCreating = data.resourceName === 'new-pv' || !currentPVForDetail.value?.creationTimestamp

    let response
    if (isCreating) {
      // 创建新PV，使用createPodFromYaml API（支持多种资源类型）
      response = await k8sApi.createPodFromYaml(selectedClusterId.value, queryParams.namespace, { yamlContent: data.yamlContent })
    } else {
      // 更新现有PV
      response = await k8sApi.updatePVYaml(selectedClusterId.value, data.resourceName, data.yamlContent)
    }

    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(isCreating ? 'PV 创建成功' : 'PV YAML 更新成功')
      pvYamlDialogVisible.value = false
      fetchPVList() // 刷新列表
    } else {
      ElMessage.error(responseData.message || (isCreating ? 'PV 创建失败' : 'PV YAML 更新失败'))
    }
  } catch (error) {
    console.error('PV 操作失败:', error)
    ElMessage.error('PV 操作失败')
  }
}

// 删除 PV
const handleDeletePV = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认删除 PV "${row?.name}"？删除后数据将无法恢复。`,
      '删除确认',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    const response = await k8sApi.deletePV(selectedClusterId.value, row?.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('PV 删除成功')
      fetchPVList()
    } else {
      ElMessage.error(responseData.message || 'PV 删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除 PV 失败:', error)
      ElMessage.error('删除 PV 失败')
    }
  }
}

// StorageClass 操作
const handleCreateStorageClass = () => {
  if (!selectedClusterId.value) {
    ElMessage.warning('请先选择集群')
    return
  }

  // 设置默认的StorageClass YAML模板
  const defaultYaml = `apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: new-storageclass
provisioner: kubernetes.io/no-provisioner
reclaimPolicy: Delete
allowVolumeExpansion: true
volumeBindingMode: WaitForFirstConsumer`

  currentStorageClassYaml.value = defaultYaml
  currentStorageClassForDetail.value = { name: 'new-storageclass' }
  storageClassYamlDialogVisible.value = true
}

// 查看 StorageClass 详情
const handleViewStorageClass = async (row) => {
  try {
    console.log('🔍 获取StorageClass详情:', {
      clusterId: selectedClusterId.value,
      storageClassName: row?.name,
      row: row
    })

    const response = await k8sApi.getStorageClassDetail(selectedClusterId.value, row?.name)
    const responseData = response.data || response

    console.log('📥 StorageClass详情API响应:', responseData)

    if (responseData.code === 200 || responseData.success) {
      currentStorageClassForDetail.value = responseData.data || row
      console.log('✅ StorageClass详情数据已设置，打开对话框')
      storageClassDetailDialogVisible.value = true
    } else {
      console.error('❌ API返回失败:', responseData)
      ElMessage.error(responseData.message || '获取 StorageClass 详情失败')
    }
  } catch (error) {
    console.error('❌ 获取 StorageClass 详情失败:', error)
    ElMessage.error('获取 StorageClass 详情失败: ' + (error.message || '未知错误'))
  }
}

// 编辑 StorageClass YAML
const handleEditStorageClassYaml = async (row) => {
  try {
    loading.value = true
    const response = await k8sApi.getStorageClassYaml(selectedClusterId.value, row?.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      currentStorageClassYaml.value = getYamlContentFromApiResponse(
        responseData,
        `# StorageClass ${row?.name} YAML\napiVersion: storage.k8s.io/v1\nkind: StorageClass\nmetadata:\n  name: ${row?.name}`
      )
      currentStorageClassForDetail.value = row
      storageClassYamlDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取 StorageClass YAML 失败')
    }
  } catch (error) {
    console.error('获取 StorageClass YAML 失败:', error)
    ElMessage.error('获取 StorageClass YAML 失败')
  } finally {
    loading.value = false
  }
}

// 保存 StorageClass YAML
const handleStorageClassYamlSave = async (data) => {
  try {
    // 检查是否是创建新StorageClass (new-storageclass表示新建)
    const isCreating = data.resourceName === 'new-storageclass' || !currentStorageClassForDetail.value?.creationTimestamp

    let response
    if (isCreating) {
      // 创建新StorageClass，使用createPodFromYaml API（支持多种资源类型）
      response = await k8sApi.createPodFromYaml(selectedClusterId.value, queryParams.namespace, { yamlContent: data.yamlContent })
    } else {
      // 更新现有StorageClass
      response = await k8sApi.updateStorageClassYaml(selectedClusterId.value, data.resourceName, data.yamlContent)
    }

    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(isCreating ? 'StorageClass 创建成功' : 'StorageClass YAML 更新成功')
      storageClassYamlDialogVisible.value = false
      fetchStorageClassList() // 刷新列表
    } else {
      ElMessage.error(responseData.message || (isCreating ? 'StorageClass 创建失败' : 'StorageClass YAML 更新失败'))
    }
  } catch (error) {
    console.error('StorageClass 操作失败:', error)
    ElMessage.error('StorageClass 操作失败')
  }
}

// 删除 StorageClass
const handleDeleteStorageClass = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认删除 StorageClass "${row?.name}"？删除后可能影响使用该存储类的PVC。`,
      '删除确认',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    const response = await k8sApi.deleteStorageClass(selectedClusterId.value, row?.name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('StorageClass 删除成功')
      fetchStorageClassList()
    } else {
      ElMessage.error(responseData.message || 'StorageClass 删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除 StorageClass 失败:', error)
      ElMessage.error('删除 StorageClass 失败')
    }
  }
}

// 工具函数
const getClusterStatusText = (status) => {
  const statusMap = {
    1: '(运行中)',
    2: '(正常)',
    3: '(异常)',
    4: '(连接失败)'
  }
  return statusMap[status] || '(未知状态)'
}

const getClusterStatusTag = (status) => {
  const tagMap = {
    1: 'success',
    2: 'success',
    3: 'warning',
    4: 'danger'
  }
  return tagMap[status] || 'info'
}

const getStatusTagType = (status) => {
  const statusMap = {
    'Available': 'success',
    'Bound': 'success',
    'Released': 'warning',
    'Failed': 'danger',
    'Pending': 'info'
  }
  return statusMap[status] || 'info'
}

// 监听集群和命名空间变化，自动加载数据
watch(
  [selectedClusterId, () => queryParams.namespace],
  ([clusterId, namespace]) => {
    console.log('监听到变化 - 集群ID:', clusterId, '命名空间:', namespace)
    if (clusterId && namespace) {
      console.log('集群和命名空间都已选择，开始加载存储资源')
      loadAllStorageResources()
    }
  },
  { immediate: true }
)

// 页面初始化
onMounted(async () => {
  console.log('🚀 开始加载k8s存储管理页面')
  const startTime = Date.now()

  try {
    console.log('🎉 页面初始化完成，总耗时:', Date.now() - startTime + 'ms')
  } catch (error) {
    console.error('页面初始化失败:', error)
  }
})
</script>

<template>
  <div class="k8s-storage-management">
    <el-card shadow="hover" class="storage-card">
      <template #header>
        <div class="card-header">
          <span class="title">K8s 存储管理</span>
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
        <el-form :inline="true" class="search-form">
          <el-form-item label="存储资源名称">
            <el-input
              v-model="searchKeyword"
              placeholder="请输入名称"
              clearable
              size="small"
              style="width: 200px"
              @keyup.enter="resetSearch"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :icon="Search" size="small" @click="resetSearch">
              搜索
            </el-button>
            <el-button type="success" :icon="Refresh" size="small" @click="handleRefresh">
              刷新
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 存储资源表格 -->
      <el-tabs v-model="activeTab" @tab-change="handleTabChange" class="storage-tabs">
        <!-- PVC 标签页 -->
        <el-tab-pane label="PVC" name="pvc">
          <div class="tab-content">
            <div class="content-header">
              <span class="resource-count">共 {{ filteredPVCList.length }} 个 PVC</span>
              <el-button type="primary" :icon="Plus" size="small" @click="handleCreatePVC">
                创建 PVC
              </el-button>
            </div>

            <el-table
              :data="filteredPVCList"
              v-loading="loading"
              element-loading-text="加载中..."
              class="resource-table"
              empty-text="暂无 PVC 资源"
            >
              <el-table-column prop="name" label="名称" min-width="150">
                <template #default="{ row }">
                  <div class="resource-name">
                    <el-icon class="resource-icon"><DataLine /></el-icon>
                    <span class="resource-name-link" @click="handleViewPVC(row)">{{ row?.name || '-' }}</span>
                  </div>
                </template>
              </el-table-column>

              <el-table-column prop="status" label="状态" width="120">
                <template #default="{ row }">
                  <el-tag
                    :type="getStatusTagType(row?.status)"
                    size="small"
                  >
                    {{ row?.status || '未知' }}
                  </el-tag>
                </template>
              </el-table-column>

              <el-table-column prop="capacity" label="容量" width="100">
                <template #default="{ row }">
                  <span>{{ row?.capacity || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column prop="accessModes" label="访问模式" width="140">
                <template #default="{ row }">
                  <div class="access-modes">
                    <el-tag
                      v-for="mode in (row?.accessModes || [])"
                      :key="mode"
                      size="small"
                      type="info"
                      class="mode-tag"
                    >
                      {{ mode }}
                    </el-tag>
                  </div>
                </template>
              </el-table-column>

              <el-table-column prop="storageClassName" label="存储类" width="150">
                <template #default="{ row }">
                  <span>{{ row?.storageClassName || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column prop="creationTimestamp" label="创建时间" width="180">
                <template #default="{ row }">
                  <span>{{ row?.creationTimestamp || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column label="操作" width="190" fixed="right">
                <template #default="{ row }">
                  <div class="operation-buttons">
                    <el-button
                      type="primary"
                      :icon="View"
                      size="small"
                      @click="handleViewPVC(row)"
                      v-authority="['k8s:storage:get']"
                    >详情</el-button>
                    <el-dropdown trigger="click">
                      <el-button type="warning" size="small" v-authority="['k8s:storage:actions']">
                        更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                      </el-button>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item :icon="Edit" @click="handleEditPVCYaml(row)">编辑 YAML</el-dropdown-item>
                          <el-dropdown-item :icon="Delete" @click="handleDeletePVC(row)">删除</el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <!-- PV 标签页 -->
        <el-tab-pane label="PV" name="pv">
          <div class="tab-content">
            <div class="content-header">
              <span class="resource-count">共 {{ filteredPVList.length }} 个 PV</span>
              <el-button type="primary" :icon="Plus" size="small" @click="handleCreatePV">
                创建 PV
              </el-button>
            </div>

            <el-table
              :data="filteredPVList"
              v-loading="loading"
              element-loading-text="加载中..."
              class="resource-table"
              empty-text="暂无 PV 资源"
            >
              <el-table-column prop="name" label="名称" min-width="150">
                <template #default="{ row }">
                  <div class="resource-name">
                    <el-icon class="resource-icon"><Cpu /></el-icon>
                    <span class="resource-name-link" @click="handleViewPV(row)">{{ row?.name || '-' }}</span>
                  </div>
                </template>
              </el-table-column>

              <el-table-column prop="status" label="状态" width="120">
                <template #default="{ row }">
                  <el-tag
                    :type="getStatusTagType(row?.status)"
                    size="small"
                  >
                    {{ row?.status || '未知' }}
                  </el-tag>
                </template>
              </el-table-column>

              <el-table-column prop="capacity" label="容量" width="100">
                <template #default="{ row }">
                  <span>{{ row?.capacity || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column prop="accessModes" label="访问模式" width="140">
                <template #default="{ row }">
                  <div class="access-modes">
                    <el-tag
                      v-for="mode in (row?.accessModes || [])"
                      :key="mode"
                      size="small"
                      type="info"
                      class="mode-tag"
                    >
                      {{ mode }}
                    </el-tag>
                  </div>
                </template>
              </el-table-column>

              <el-table-column prop="reclaimPolicy" label="回收策略" width="120">
                <template #default="{ row }">
                  <span>{{ row?.reclaimPolicy || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column prop="storageClassName" label="存储类" width="150">
                <template #default="{ row }">
                  <span>{{ row?.storageClassName || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column prop="creationTimestamp" label="创建时间" width="180">
                <template #default="{ row }">
                  <span>{{ row?.creationTimestamp || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column label="操作" width="190" fixed="right">
                <template #default="{ row }">
                  <div class="operation-buttons">
                    <el-button
                      type="primary"
                      :icon="View"
                      size="small"
                      @click="handleViewPV(row)"
                      v-authority="['k8s:storage:get']"
                    >详情</el-button>
                    <el-dropdown trigger="click">
                      <el-button type="warning" size="small" v-authority="['k8s:storage:actions']">
                        更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                      </el-button>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item :icon="Edit" @click="handleEditPVYaml(row)">编辑 YAML</el-dropdown-item>
                          <el-dropdown-item :icon="Delete" @click="handleDeletePV(row)">删除</el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <!-- StorageClass 标签页 -->
        <el-tab-pane label="StorageClass" name="storageclass">
          <div class="tab-content">
            <div class="content-header">
              <span class="resource-count">共 {{ filteredStorageClassList.length }} 个 StorageClass</span>
              <el-button type="primary" :icon="Plus" size="small" @click="handleCreateStorageClass">
                创建 StorageClass
              </el-button>
            </div>

            <el-table
              :data="filteredStorageClassList"
              v-loading="loading"
              element-loading-text="加载中..."
              class="resource-table"
              empty-text="暂无 StorageClass 资源"
            >
              <el-table-column prop="name" label="名称" min-width="150">
                <template #default="{ row }">
                  <div class="resource-name">
                    <el-icon class="resource-icon"><Connection /></el-icon>
                    <span class="resource-name-link" @click="handleViewStorageClass(row)">{{ row?.name || '-' }}</span>
                  </div>
                </template>
              </el-table-column>

              <el-table-column prop="provisioner" label="供应商" min-width="200">
                <template #default="{ row }">
                  <span>{{ row?.provisioner || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column prop="reclaimPolicy" label="回收策略" width="120">
                <template #default="{ row }">
                  <span>{{ row?.reclaimPolicy || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column prop="allowVolumeExpansion" label="支持扩容" width="100">
                <template #default="{ row }">
                  <el-tag
                    :type="row?.allowVolumeExpansion ? 'success' : 'info'"
                    size="small"
                  >
                    {{ row?.allowVolumeExpansion ? '是' : '否' }}
                  </el-tag>
                </template>
              </el-table-column>

              <el-table-column prop="volumeBindingMode" label="绑定模式" width="150">
                <template #default="{ row }">
                  <span>{{ row?.volumeBindingMode || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column prop="creationTimestamp" label="创建时间" width="180">
                <template #default="{ row }">
                  <span>{{ row?.creationTimestamp || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column label="操作" width="190" fixed="right">
                <template #default="{ row }">
                  <div class="operation-buttons">
                    <el-button
                      type="primary"
                      :icon="View"
                      size="small"
                      @click="handleViewStorageClass(row)"
                      v-authority="['k8s:storage:get']"
                    >详情</el-button>
                    <el-dropdown trigger="click">
                      <el-button type="warning" size="small" v-authority="['k8s:storage:actions']">
                        更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                      </el-button>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item :icon="Edit" @click="handleEditStorageClassYaml(row)">编辑 YAML</el-dropdown-item>
                          <el-dropdown-item :icon="Delete" @click="handleDeleteStorageClass(row)">删除</el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- PVC 详情对话框 -->
    <el-dialog
      v-model="pvcDetailDialogVisible"
      title="PVC 详情"
      width="60%"
      class="detail-dialog"
    >
      <div class="detail-content">
        <div class="detail-item">
          <span class="detail-label">名称:</span>
          <span class="detail-value">{{ currentPVCForDetail.name || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">命名空间:</span>
          <span class="detail-value">{{ currentPVCForDetail.namespace || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">状态:</span>
          <el-tag :type="getStatusTagType(currentPVCForDetail.status)" size="small">
            {{ currentPVCForDetail.status || 'Unknown' }}
          </el-tag>
        </div>
        <div class="detail-item">
          <span class="detail-label">容量:</span>
          <span class="detail-value">{{ currentPVCForDetail.capacity || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">访问模式:</span>
          <div class="access-modes">
            <el-tag
              v-for="mode in (currentPVCForDetail.accessModes || [])"
              :key="mode"
              size="small"
              type="info"
              class="mode-tag"
            >
              {{ mode }}
            </el-tag>
          </div>
        </div>
        <div class="detail-item">
          <span class="detail-label">存储类:</span>
          <span class="detail-value">{{ currentPVCForDetail.storageClass || currentPVCForDetail.storageClassName || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">创建时间:</span>
          <span class="detail-value">{{ currentPVCForDetail.createdAt ? new Date(currentPVCForDetail.createdAt).toLocaleString() : (currentPVCForDetail.creationTimestamp ? new Date(currentPVCForDetail.creationTimestamp).toLocaleString() : '-') }}</span>
        </div>
        <div class="detail-item" v-if="currentPVCForDetail.volumeName">
          <span class="detail-label">绑定的PV:</span>
          <span class="detail-value">{{ currentPVCForDetail.volumeName }}</span>
        </div>
        <div class="detail-item" v-if="currentPVCForDetail.volumeMode">
          <span class="detail-label">卷模式:</span>
          <span class="detail-value">{{ currentPVCForDetail.volumeMode }}</span>
        </div>
      </div>
    </el-dialog>

    <!-- PV 详情对话框 -->
    <el-dialog
      v-model="pvDetailDialogVisible"
      title="PV 详情"
      width="60%"
      class="detail-dialog"
    >
      <div class="detail-content">
        <div class="detail-item">
          <span class="detail-label">名称:</span>
          <span class="detail-value">{{ currentPVForDetail.name || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">状态:</span>
          <el-tag :type="getStatusTagType(currentPVForDetail.status)" size="small">
            {{ currentPVForDetail.status || 'Unknown' }}
          </el-tag>
        </div>
        <div class="detail-item">
          <span class="detail-label">容量:</span>
          <span class="detail-value">{{ currentPVForDetail.capacity || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">访问模式:</span>
          <div class="access-modes">
            <el-tag
              v-for="mode in (currentPVForDetail.accessModes || [])"
              :key="mode"
              size="small"
              type="info"
              class="mode-tag"
            >
              {{ mode }}
            </el-tag>
          </div>
        </div>
        <div class="detail-item">
          <span class="detail-label">回收策略:</span>
          <span class="detail-value">{{ currentPVForDetail.reclaimPolicy || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">存储类:</span>
          <span class="detail-value">{{ currentPVForDetail.storageClass || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">创建时间:</span>
          <span class="detail-value">{{ currentPVForDetail.createdAt ? new Date(currentPVForDetail.createdAt).toLocaleString() : '-' }}</span>
        </div>
        <div class="detail-item" v-if="currentPVForDetail.volumeMode">
          <span class="detail-label">卷模式:</span>
          <span class="detail-value">{{ currentPVForDetail.volumeMode }}</span>
        </div>
        <div class="detail-item" v-if="currentPVForDetail.persistentVolumeSource">
          <span class="detail-label">存储源:</span>
          <div class="storage-source">
            <div v-if="currentPVForDetail.persistentVolumeSource.hostPath">
              <strong>HostPath:</strong> {{ currentPVForDetail.persistentVolumeSource.hostPath.path }}
              <br>
              <small>类型: {{ currentPVForDetail.persistentVolumeSource.hostPath.type }}</small>
            </div>
            <div v-else-if="currentPVForDetail.persistentVolumeSource.nfs">
              <strong>NFS:</strong> {{ currentPVForDetail.persistentVolumeSource.nfs.server }}:{{ currentPVForDetail.persistentVolumeSource.nfs.path }}
            </div>
            <div v-else>
              {{ Object.keys(currentPVForDetail.persistentVolumeSource)[0] }}
            </div>
          </div>
        </div>
        <div class="detail-item" v-if="currentPVForDetail.claimRef">
          <span class="detail-label">绑定的PVC:</span>
          <span class="detail-value">{{ currentPVForDetail.claimRef.namespace }}/{{ currentPVForDetail.claimRef.name }}</span>
        </div>
      </div>
    </el-dialog>

    <!-- StorageClass 详情对话框 -->
    <el-dialog
      v-model="storageClassDetailDialogVisible"
      title="StorageClass 详情"
      width="60%"
      class="detail-dialog"
    >
      <div class="detail-content">
        <div class="detail-item">
          <span class="detail-label">名称:</span>
          <span class="detail-value">{{ currentStorageClassForDetail.name }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">供应商:</span>
          <span class="detail-value">{{ currentStorageClassForDetail.provisioner }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">回收策略:</span>
          <span class="detail-value">{{ currentStorageClassForDetail.reclaimPolicy }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">支持扩容:</span>
          <el-tag
            :type="currentStorageClassForDetail.allowVolumeExpansion ? 'success' : 'info'"
            size="small"
          >
            {{ currentStorageClassForDetail.allowVolumeExpansion ? '是' : '否' }}
          </el-tag>
        </div>
        <div class="detail-item">
          <span class="detail-label">绑定模式:</span>
          <span class="detail-value">{{ currentStorageClassForDetail.volumeBindingMode }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">创建时间:</span>
          <span class="detail-value">{{ currentStorageClassForDetail.creationTimestamp }}</span>
        </div>
      </div>
    </el-dialog>

    <!-- PVC YAML 编辑对话框 -->
    <PodYamlDialog
      :visible="pvcYamlDialogVisible"
      :yaml-content="currentPVCYaml"
      :resource-name="currentPVCForDetail.name || 'new-pvc'"
      :resource-type="'PVC'"
      :editable="true"
      @update:visible="pvcYamlDialogVisible = $event"
      @close="pvcYamlDialogVisible = false"
      @save="handlePVCYamlSave"
    />

    <!-- PV YAML 编辑对话框 -->
    <PodYamlDialog
      :visible="pvYamlDialogVisible"
      :yaml-content="currentPVYaml"
      :resource-name="currentPVForDetail.name || 'new-pv'"
      :resource-type="'PV'"
      :editable="true"
      @update:visible="pvYamlDialogVisible = $event"
      @close="pvYamlDialogVisible = false"
      @save="handlePVYamlSave"
    />

    <!-- StorageClass YAML 编辑对话框 -->
    <PodYamlDialog
      :visible="storageClassYamlDialogVisible"
      :yaml-content="currentStorageClassYaml"
      :resource-name="currentStorageClassForDetail.name || 'new-storageclass'"
      :resource-type="'StorageClass'"
      :editable="true"
      @update:visible="storageClassYamlDialogVisible = $event"
      @close="storageClassYamlDialogVisible = false"
      @save="handleStorageClassYamlSave"
    />
  </div>
</template>

<style scoped>
.k8s-storage-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.storage-card {
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

/* 按钮样式 - 与k8s-clusters.vue保持一致 */
.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 输入框样式 - 与k8s-clusters.vue保持一致 */
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

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #6b9bd1;
}

/* 标签样式 */
.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

.storage-tabs {
  margin-top: 20px;
}

.storage-tabs :deep(.el-tabs__header) {
  margin-bottom: 20px;
}

.storage-tabs :deep(.el-tabs__item) {
  font-weight: 500;
  color: #606266;
}

.storage-tabs :deep(.el-tabs__item.is-active) {
  color: #409EFF;
  font-weight: 600;
}

.tab-content {
  padding: 0;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px 0;
}

.resource-count {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.resource-table {
  border-radius: 8px;
  overflow: hidden;
}

.resource-table :deep(.el-table__header) {
  background: #f8f9fa;
}

.resource-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

.resource-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.resource-icon {
  color: #409EFF;
  font-size: 16px;
}

.resource-name-link {
  color: #409EFF;
  cursor: pointer;
  font-weight: 500;
  transition: color 0.3s;
}

.resource-name-link:hover {
  color: #66b1ff;
  text-decoration: underline;
}

.access-modes {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.mode-tag {
  font-size: 11px;
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

/* 对话框样式 - 与k8s-clusters.vue保持一致 */
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

.detail-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

/* 加载动画样式 */
.el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

.detail-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.detail-label {
  font-weight: 600;
  color: #606266;
  min-width: 80px;
}

.detail-value {
  color: #303133;
  word-break: break-all;
}

.storage-source {
  background: #f8f9fa;
  padding: 8px 12px;
  border-radius: 6px;
  border-left: 3px solid #409eff;
  font-size: 13px;
  line-height: 1.6;
}

.storage-source strong {
  color: #409eff;
}

.storage-source small {
  color: #909399;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .k8s-storage-management {
    padding: 12px;
  }

  .card-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .header-actions {
    justify-content: flex-end;
  }

  .detail-content {
    grid-template-columns: 1fr;
  }

  .resource-table :deep(.el-table) {
    font-size: 12px;
  }

  .operation-buttons {
    flex-direction: column;
    gap: 4px;
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
}
</style>
