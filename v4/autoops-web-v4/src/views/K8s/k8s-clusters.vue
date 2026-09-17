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
  Connection,
  Setting
} from '@element-plus/icons-vue'
import cmdbApi from '@/api/cmdb'
import k8sApi from '@/api/k8s'
import HostSelector from '@/components/HostSelector.vue'
import ClusterTable from './clusters/ClusterTable.vue'
import CredentialEditor from './clusters/components/CredentialEditor.vue'
import CredentialDialog from './clusters/dialogs/CredentialDialog.vue'
import RegisterDialog from './clusters/dialogs/RegisterDialog.vue'
import EditDialog from './clusters/dialogs/EditDialog.vue'

const router = useRouter()

const loading = ref(false)
const clusterTableRef = ref(null)
const refreshTrigger = ref(0)

// 移除了 statusOptions，因为编辑对话框不再需要状态选择
// 状态相关逻辑已移至 ClusterTable 组件中处理

const registerDialogVisible = ref(false)
const createDialogVisible = ref(false)
const editDialogVisible = ref(false)
const credentialDialogVisible = ref(false)

// 编辑时的集群数据
const currentEditCluster = reactive({
  id: '',
  clusterName: '',
  version: '',
  credential: '',
  description: ''
})

const createForm = reactive({
  name: '',
  version: '1.30.4',
  description: '',
  nodeConfig: {
    masterHostIds: [],
    workerHostIds: [],
    etcdHostIds: []
  },
  autoDeploy: false,
  deploymentMode: 1,
  networkPlugin: 'calico', // 默认选择 Calico
  enabledComponents: ['coredns'], // 默认包含必选组件
  usePrivateRegistry: false, // 是否使用私有仓库
  registryConfig: {
    privateRegistry: 'crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s',
    registryUsername: 'zhangfan5391621',
    registryPassword: 'zf5391621'
  }
})

const hostSelectorVisible = ref(false)
const currentNodeType = ref('Master')
const selectedMasterHosts = ref([])
const selectedWorkerHosts = ref([])
const selectedEtcdHosts = ref([])

const currentCredential = reactive({
  clusterName: '',
  content: ''
})

// 网络插件配置（二选一）
const networkPluginOptions = [
  { label: 'Calico', value: 'calico', description: '高性能网络插件，支持网络策略' },
  { label: 'Flannel', value: 'flannel', description: '轻量级网络插件，配置简单' }
]

// 必选插件
const requiredComponents = [
  { label: 'CoreDNS', value: 'coredns', description: 'Kubernetes集群DNS服务', required: true }
]

// 可选插件
const optionalComponents = [
  { label: 'Metrics Server', value: 'metrics-server', description: '集群资源监控' },
  { label: 'Ingress Nginx', value: 'ingress-nginx', description: 'HTTP/HTTPS负载均衡' },
  { label: 'Dashboard', value: 'dashboard', description: 'Kubernetes Web管理界面' },
  { label: 'Prometheus', value: 'prometheus', description: '监控和告警系统' },
  { label: 'Grafana', value: 'grafana', description: '监控数据可视化' }
]

const k8sVersionOptions = [
  { label: 'v1.30.4', value: '1.30.4' },
  { label: 'v1.31.0', value: '1.31.0' },
  { label: 'v1.32.0', value: '1.32.0' },
]

// 集群表格相关事件处理
const handleClusterSync = (row) => {
  refreshTrigger.value++
}

const handleClusterDelete = (row) => {
  refreshTrigger.value++
}

const handleClusterEdit = async (row) => {
  console.log('handleClusterEdit 被调用，参数:', row)
  try {
    loading.value = true
    console.log('调用后端API获取集群详情，ID:', row.id)
    
    // 调用后端 API 获取完整的集群详情
    const response = await k8sApi.getClusterById(row.id)
    
    // response可能是axios响应对象，需要获取data部分
    const responseData = response.data || response
    console.log('获取集群详情API响应:', responseData)
    
    if (responseData.code === 200 || responseData.success) {
      // 兼容不同的数据结构：有些API返回data.cluster，有些直接返回data
      const clusterData = responseData.data.cluster || responseData.data
      
      console.log('原始集群数据:', clusterData)
      
      // 填充编辑集群数据
      currentEditCluster.id = clusterData.id || row.id
      currentEditCluster.clusterName = clusterData.name || row.clusterName
      currentEditCluster.version = clusterData.version || row.version || ''
      currentEditCluster.credential = clusterData.credential || clusterData.kubeconfig || row.credential || ''
      currentEditCluster.description = clusterData.description || clusterData.remark || row.description || ''
      
      console.log('编辑集群数据已填充:', {
        id: currentEditCluster.id,
        clusterName: currentEditCluster.clusterName,
        credentialLength: currentEditCluster.credential?.length || 0
      })
      
      // 打开编辑对话框
      editDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取集群详情失败')
    }
  } catch (error) {
    console.error('获取集群详情失败:', error)
    ElMessage.error('获取集群详情失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 查看集群凭证
const handleViewCredential = (row) => {
  currentCredential.clusterName = row.clusterName
  currentCredential.content = row.credential || ''
  credentialDialogVisible.value = true
}



const handleRegister = () => {
  registerDialogVisible.value = true
}

const handleCreate = () => {
  createDialogVisible.value = true
  // 检查后端连接状态
  checkBackendConnection()
}

const checkBackendConnection = async () => {
  try {
    await k8sApi.getClusterList()
  } catch (error) {
    if (error.code === 'ERR_NETWORK' || error.message.includes('ERR_CONNECTION_REFUSED')) {
      ElMessage.warning('后端服务未启动，主机选择功能可能不可用')
    }
  }
}

const openHostSelector = (nodeType) => {
  currentNodeType.value = nodeType
  hostSelectorVisible.value = true
}

const handleHostsSelected = (hosts) => {
  const hostIds = hosts.map(host => host.id)
  
  if (currentNodeType.value === 'Master') {
    // 单节点模式只能选择一台Master
    if (createForm.deploymentMode === 1 && hosts.length > 1) {
      ElMessage.warning('单节点模式只能选择一台Master节点主机')
      return
    }
    
    selectedMasterHosts.value = hosts
    createForm.nodeConfig.masterHostIds = hostIds
    // 如果没有设置ETCD节点，默认使用Master节点
    if (selectedEtcdHosts.value.length === 0) {
      selectedEtcdHosts.value = [...hosts]
      createForm.nodeConfig.etcdHostIds = [...hostIds]
    }
  } else if (currentNodeType.value === 'Worker') {
    selectedWorkerHosts.value = hosts
    createForm.nodeConfig.workerHostIds = hostIds
  } else if (currentNodeType.value === 'ETCD') {
    selectedEtcdHosts.value = hosts
    createForm.nodeConfig.etcdHostIds = hostIds
  }
  
  ElMessage.success(`已选择 ${hosts.length} 台${currentNodeType.value}节点主机`)
}

const removeHost = (hostId, nodeType) => {
  if (nodeType === 'Master') {
    selectedMasterHosts.value = selectedMasterHosts.value.filter(host => host.id !== hostId)
    createForm.nodeConfig.masterHostIds = createForm.nodeConfig.masterHostIds.filter(id => id !== hostId)
  } else if (nodeType === 'Worker') {
    selectedWorkerHosts.value = selectedWorkerHosts.value.filter(host => host.id !== hostId)
    createForm.nodeConfig.workerHostIds = createForm.nodeConfig.workerHostIds.filter(id => id !== hostId)
  } else if (nodeType === 'ETCD') {
    selectedEtcdHosts.value = selectedEtcdHosts.value.filter(host => host.id !== hostId)
    createForm.nodeConfig.etcdHostIds = createForm.nodeConfig.etcdHostIds.filter(id => id !== hostId)
  }
}

const getExcludeHosts = () => {
  // 获取已选择的所有主机ID，避免重复选择
  const allSelectedIds = [
    ...createForm.nodeConfig.masterHostIds,
    ...createForm.nodeConfig.workerHostIds,
    ...createForm.nodeConfig.etcdHostIds
  ]
  return [...new Set(allSelectedIds)]
}

// 部署模式切换处理
const handleDeploymentModeChange = () => {
  // 如果从集群模式切换到单节点模式，且已选择多个Master节点
  if (createForm.deploymentMode === 1 && selectedMasterHosts.value.length > 1) {
    ElMessageBox.confirm(
      '切换到单节点模式将只保留第一个Master节点，其他Master节点将被移除。是否继续？',
      '模式切换确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    ).then(() => {
      // 只保留第一个Master节点
      const firstMaster = selectedMasterHosts.value[0]
      selectedMasterHosts.value = [firstMaster]
      createForm.nodeConfig.masterHostIds = [firstMaster.id]
      
      // 同步更新ETCD节点
      if (selectedEtcdHosts.value.length === 0 || 
          selectedEtcdHosts.value.some(host => host.id === firstMaster.id)) {
        selectedEtcdHosts.value = [firstMaster]
        createForm.nodeConfig.etcdHostIds = [firstMaster.id]
      }
      
      ElMessage.success('已切换到单节点模式，只保留一个Master节点')
    }).catch(() => {
      // 用户取消，恢复到集群模式
      createForm.deploymentMode = 2
    })
  }
}

// 注册成功事件处理
const handleRegisterSuccess = () => {
  console.log('注册成功，触发立即刷新')
  refreshTrigger.value++
  
  // 如果可以访问子组件，直接调用刷新方法
  if (clusterTableRef.value) {
    clusterTableRef.value.handleQuery()
  }
  
  // 延时刷新，等待后端自动同步完成
  setTimeout(() => {
    console.log('延时刷新，等待后端同步完成')
    refreshTrigger.value++
    if (clusterTableRef.value) {
      clusterTableRef.value.handleQuery()
    }
  }, 3000)
}

// 注册对话框关闭事件处理
const handleRegisterClosed = () => {
  console.log('注册对话框已关闭')
}


// 编辑成功事件处理
const handleEditSuccess = () => {
  console.log('编辑成功，触发表格刷新')
  refreshTrigger.value++
}

// 编辑对话框关闭事件处理
const handleEditClosed = () => {
  console.log('编辑对话框已关闭')
}


const submitCreate = async () => {
  try {
    if (!createForm.name) {
      ElMessage.warning('请输入集群名称')
      return
    }
    if (!createForm.nodeConfig.masterHostIds.length) {
      ElMessage.warning('请选择Master节点')
      return
    }
    
    // 单节点模式验证
    if (createForm.deploymentMode === 1) {
      if (createForm.nodeConfig.masterHostIds.length > 1) {
        ElMessage.warning('单节点模式只能选择1个Master节点')
        return
      }
    }
    
    // 集群模式验证
    if (createForm.deploymentMode === 2) {
      if (createForm.nodeConfig.masterHostIds.length < 3) {
        ElMessage.warning('集群模式至少需要选择3个Master节点以确保高可用')
        return
      }
    }
    
    if (!createForm.nodeConfig.etcdHostIds.length) {
      createForm.nodeConfig.etcdHostIds = [...createForm.nodeConfig.masterHostIds]
    }

    // 确保网络插件和CoreDNS包含在组件列表中
    if (!createForm.enabledComponents.includes('coredns')) {
      createForm.enabledComponents.push('coredns')
    }
    if (!createForm.enabledComponents.includes(createForm.networkPlugin)) {
      createForm.enabledComponents.push(createForm.networkPlugin)
    }

    const response = await k8sApi.createCluster(createForm)
    
    // response可能是axios响应对象，需要获取data部分
    const responseData = response.data || response
    console.log('创建集群API响应:', responseData)

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('K8s集群创建任务已提交成功！')
      createDialogVisible.value = false
      resetCreateForm()
      refreshTrigger.value++  // 触发表格刷新
    } else {
      ElMessage.error(responseData.message || '创建集群失败')
    }
  } catch (error) {
    console.error('创建集群失败:', error)
    ElMessage.error(error.response?.data?.message || '创建集群失败，请检查网络连接')
  }
}

const resetCreateForm = () => {
  createForm.name = ''
  createForm.version = '1.30.4'
  createForm.description = ''
  createForm.nodeConfig.masterHostIds = []
  createForm.nodeConfig.workerHostIds = []
  createForm.nodeConfig.etcdHostIds = []
  createForm.autoDeploy = false
  createForm.deploymentMode = 1
  createForm.networkPlugin = 'calico'
  createForm.enabledComponents = ['coredns']
  createForm.usePrivateRegistry = false
  createForm.registryConfig.privateRegistry = 'crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s'
  createForm.registryConfig.registryUsername = 'zhangfan5391621'
  createForm.registryConfig.registryPassword = 'zf5391621'
  
  // 清空选中的主机
  selectedMasterHosts.value = []
  selectedWorkerHosts.value = []
  selectedEtcdHosts.value = []
}

// 处理可选组件的选择
const toggleOptionalComponent = (componentValue) => {
  const index = createForm.enabledComponents.indexOf(componentValue)
  if (index > -1) {
    // 如果已存在，则移除（但不能移除必选组件）
    if (componentValue !== 'coredns') {
      createForm.enabledComponents.splice(index, 1)
    }
  } else {
    // 如果不存在，则添加
    createForm.enabledComponents.push(componentValue)
  }
}

onMounted(() => {
  // 组件挂载时不需要手动调用，ClusterTable 会自动加载数据
})
</script>

<template>
  <div class="k8s-cluster-management">
    <el-card shadow="hover" class="cluster-card">
      <template #header>
        <div class="card-header">
          <span class="title">K8s 集群管理</span>
          <div class="header-actions">
            <el-button type="primary" :icon="Plus" v-authority="['cloud:k8s:register']" size="small" @click="handleRegister">
              注册集群
            </el-button>
            <el-button type="success" :icon="Plus" v-authority="['cloud:k8s:add']" size="small" @click="handleCreate">
              创建集群
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 集群表格组件 -->
      <ClusterTable
        ref="clusterTableRef"
        :refresh-trigger="refreshTrigger"
        @cluster-sync="handleClusterSync"
        @cluster-edit="handleClusterEdit"
        @cluster-delete="handleClusterDelete"
        @view-credential="handleViewCredential"
      />
    </el-card>

    <!-- 注册集群对话框 -->
    <RegisterDialog
      v-model:visible="registerDialogVisible"
      @success="handleRegisterSuccess"
      @closed="handleRegisterClosed"
    />

    <!-- 编辑集群对话框 -->
    <EditDialog
      v-model:visible="editDialogVisible"
      :cluster="currentEditCluster"
      @success="handleEditSuccess"
      @closed="handleEditClosed"
    />

    <!-- 创建集群对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="创建 K8s 集群"
      width="700px"
      class="create-dialog compact-dialog"
    >
      <el-form :model="createForm" label-width="90px" class="compact-form">
        <!-- 基础配置行 -->
        <div class="form-row">
          <el-form-item label="集群名称" required class="form-item-half">
            <el-input
              v-model="createForm.name"
              placeholder="请输入集群名称"
              size="small"
            />
          </el-form-item>
          <el-form-item label="K8s版本" required class="form-item-half">
            <el-select v-model="createForm.version" style="width: 100%" size="small">
              <el-option 
                v-for="option in k8sVersionOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value" 
              />
            </el-select>
          </el-form-item>
        </div>
        
        <!-- 部署模式和自动部署行 -->
        <div class="form-row">
          <el-form-item label="部署模式" required class="form-item-half">
            <el-radio-group v-model="createForm.deploymentMode" @change="handleDeploymentModeChange" size="small">
              <el-radio :label="1">单节点</el-radio>
              <el-radio :label="2">集群</el-radio>
            </el-radio-group>
            <div class="mode-tip">
              {{ createForm.deploymentMode === 1 ? '适合测试' : '适合生产' }}
            </div>
          </el-form-item>
          <el-form-item label="自动部署" class="form-item-half">
            <el-switch
              v-model="createForm.autoDeploy"
              active-text="立即部署"
              inactive-text="稍后部署"
              size="small"
            />
            <div class="auto-deploy-tip">
              开启后自动执行部署脚本
            </div>
          </el-form-item>
        </div>
        
        <!-- 节点选择区域 -->
        <el-form-item label="节点配置" required class="nodes-config">
          <div class="nodes-selection">
            <!-- Master节点 -->
            <div class="node-section">
              <div class="node-header">
                <el-button 
                  type="primary" 
                  size="small"
                  @click="openHostSelector('Master')"
                >
                  Master{{ createForm.deploymentMode === 1 ? '(1台)' : '(≥3台)' }}
                </el-button>
                <span class="node-count">{{ selectedMasterHosts.length }}台</span>
              </div>
              <div v-if="selectedMasterHosts.length > 0" class="selected-tags">
                <el-tag
                  v-for="host in selectedMasterHosts"
                  :key="host.id"
                  closable
                  @close="removeHost(host.id, 'Master')"
                  type="success"
                  size="small"
                >
                  {{ host.name }}
                </el-tag>
              </div>
            </div>

            <!-- Worker节点 -->
            <div class="node-section">
              <div class="node-header">
                <el-button 
                  type="info" 
                  size="small"
                  @click="openHostSelector('Worker')"
                >
                  Worker(可选)
                </el-button>
                <span class="node-count">{{ selectedWorkerHosts.length }}台</span>
              </div>
              <div v-if="selectedWorkerHosts.length > 0" class="selected-tags">
                <el-tag
                  v-for="host in selectedWorkerHosts"
                  :key="host.id"
                  closable
                  @close="removeHost(host.id, 'Worker')"
                  type="primary"
                  size="small"
                >
                  {{ host.name }}
                </el-tag>
              </div>
            </div>

            <!-- ETCD节点 -->
            <div class="node-section">
              <div class="node-header">
                <el-button 
                  type="warning" 
                  size="small"
                  @click="openHostSelector('ETCD')"
                >
                  ETCD(可选)
                </el-button>
                <span class="node-count">{{ selectedEtcdHosts.length }}台</span>
              </div>
              <div v-if="selectedEtcdHosts.length > 0" class="selected-tags">
                <el-tag
                  v-for="host in selectedEtcdHosts"
                  :key="host.id"
                  closable
                  @close="removeHost(host.id, 'ETCD')"
                  type="warning"
                  size="small"
                >
                  {{ host.name }}
                </el-tag>
              </div>
              <div v-else class="etcd-tip">默认使用Master节点</div>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item label="网络插件" required>
          <div class="plugin-single-line">
            <div class="plugin-group">
              <div class="plugin-buttons">
                <el-button
                  v-for="plugin in networkPluginOptions"
                  :key="plugin.value"
                  :type="createForm.networkPlugin === plugin.value ? 'primary' : ''"
                  :plain="createForm.networkPlugin !== plugin.value"
                  size="small"
                  @click="createForm.networkPlugin = plugin.value"
                >
                  {{ plugin.label }}
                </el-button>
              </div>
            </div>
            
            <div class="plugin-group">
              <span class="plugin-label">必选组件:</span>
              <div class="plugin-buttons">
                <el-button
                  v-for="component in requiredComponents"
                  :key="component.value"
                  type="success"
                  size="small"
                  disabled
                >
                  {{ component.label }}
                  <el-tag type="success" size="small" style="margin-left: 4px;">必选</el-tag>
                </el-button>
              </div>
            </div>
          </div>
        </el-form-item>

        <el-form-item label="可选组件">
          <div class="plugin-single-line">
            <div class="plugin-group full-width">
              <div class="plugin-buttons">
                <el-button
                  v-for="component in optionalComponents"
                  :key="component.value"
                  :type="createForm.enabledComponents.includes(component.value) ? 'warning' : ''"
                  :plain="!createForm.enabledComponents.includes(component.value)"
                  size="small"
                  @click="toggleOptionalComponent(component.value)"
                >
                  {{ component.label }}
                </el-button>
              </div>
            </div>
          </div>
        </el-form-item>
        
        <!-- 镜像仓库配置（可选） -->
        <el-form-item label="镜像仓库" class="registry-config">
          <div class="form-row">
            <el-form-item label="" class="form-item-half">
              <el-switch
                v-model="createForm.usePrivateRegistry"
                active-text="使用私有仓库"
                inactive-text="使用默认仓库"
                size="small"
              />
            </el-form-item>
          </div>
          
          <div v-if="createForm.usePrivateRegistry" class="private-registry-config">
            <div class="form-row">
              <el-form-item label="" class="form-item-full">
                <el-input
                  v-model="createForm.registryConfig.privateRegistry"
                  placeholder="私有镜像仓库地址，如：registry.example.com"
                  size="small"
                >
                  <template #prepend>仓库地址</template>
                </el-input>
              </el-form-item>
            </div>
            <div class="form-row">
              <el-form-item label="" class="form-item-half">
                <el-input
                  v-model="createForm.registryConfig.registryUsername"
                  placeholder="仓库用户名"
                  size="small"
                >
                  <template #prepend>用户名</template>
                </el-input>
              </el-form-item>
              <el-form-item label="" class="form-item-half">
                <el-input
                  v-model="createForm.registryConfig.registryPassword"
                  type="password"
                  placeholder="仓库密码"
                  size="small"
                  show-password
                >
                  <template #prepend>密码</template>
                </el-input>
              </el-form-item>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item label="描述信息">
          <el-input
            v-model="createForm.description"
            type="textarea"
            :rows="2"
            placeholder="请输入集群描述"
            size="small"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createDialogVisible = false">取消</el-button>
          <el-button type="success" @click="submitCreate">创建</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 凭证查看对话框 -->
    <CredentialDialog
      v-model:visible="credentialDialogVisible"
      :credential="currentCredential"
    />

    <!-- 主机选择器 -->
    <HostSelector
      v-model="hostSelectorVisible"
      :node-type="currentNodeType"
      :exclude-hosts="getExcludeHosts()"
      :multiple="!(currentNodeType === 'Master' && createForm.deploymentMode === 1)"
      @hosts-selected="handleHostsSelected"
    />
  </div>
</template>

<style scoped>
.k8s-cluster-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.cluster-card {
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


.register-dialog :deep(.el-dialog),
.edit-dialog :deep(.el-dialog),
.create-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.register-dialog :deep(.el-dialog__header),
.edit-dialog :deep(.el-dialog__header),
.create-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.register-dialog :deep(.el-dialog__title),
.edit-dialog :deep(.el-dialog__title),
.create-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.register-dialog :deep(.el-dialog__body),
.edit-dialog :deep(.el-dialog__body),
.create-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

.dialog-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

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

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #6b9bd1;
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

.el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

.host-selection-container {
  width: 100%;
}

.selected-hosts {
  background: rgba(103, 126, 234, 0.05);
  padding: 12px;
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.2);
}

.hosts-info {
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
}

.host-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.no-hosts {
  color: #909399;
  font-size: 13px;
  font-style: italic;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 6px;
  border: 1px dashed #dcdfe6;
}

.el-tag {
  max-width: 200px;
}

.el-tag span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 紧凑表单样式 */
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

.form-item-full {
  width: 100%;
  margin-bottom: 0 !important;
}

.mode-tip, .auto-deploy-tip {
  font-size: 10px;
  color: #909399;
  margin-top: 2px;
  line-height: 1.2;
}

/* 节点配置区域 */
.nodes-config :deep(.el-form-item__content) {
  line-height: normal;
}

.nodes-selection {
  display: flex;
  gap: 16px;
  width: 100%;
}

.node-section {
  flex: 1;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  padding: 8px;
  background: #fafafa;
}

.node-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.node-count {
  font-size: 11px;
  color: #606266;
  font-weight: 500;
}

.selected-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 3px;
  margin-top: 4px;
}

.selected-tags .el-tag {
  font-size: 10px;
  height: 20px;
  line-height: 18px;
  max-width: 80px;
}

.selected-tags .el-tag span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 60px;
  display: inline-block;
}

.etcd-tip {
  font-size: 10px;
  color: #909399;
  text-align: center;
  margin-top: 4px;
  font-style: italic;
}

/* 插件按钮选择样式 */
.plugin-single-line {
  display: flex;
  align-items: center;
  width: 100%;
  gap: 20px;
  flex-wrap: nowrap;
}

.plugin-group {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.plugin-group.full-width {
  flex: 1;
  min-width: 0;
}

.plugin-label {
  font-size: 13px;
  font-weight: 600;
  color: #606266;
  white-space: nowrap;
  flex-shrink: 0;
}

.plugin-buttons {
  display: flex;
  gap: 6px;
  flex-wrap: nowrap;
}

/* 按钮样式调整 */
.plugin-single-line .el-button {
  border-radius: 12px;
  font-size: 10px;
  padding: 3px 6px;
  height: 22px;
  line-height: 1;
  white-space: nowrap;
  flex-shrink: 0;
}

.plugin-single-line .el-button.is-disabled {
  opacity: 0.8;
}

.plugin-single-line .el-tag {
  height: 14px;
  line-height: 12px;
  font-size: 9px;
  padding: 0 3px;
}

/* 响应式设计 */
@media (max-width: 1000px) {
  .plugin-single-line {
    flex-wrap: wrap;
    gap: 12px;
  }
  
  .plugin-group.full-width {
    width: 100%;
  }
  
  .plugin-buttons {
    flex-wrap: wrap;
  }
}

@media (max-width: 768px) {
  .compact-dialog {
    width: 95% !important;
  }
  
  .form-row {
    flex-direction: column;
    gap: 8px;
  }
  
  .nodes-selection {
    flex-direction: column;
    gap: 8px;
  }
  
  .node-section {
    padding: 6px;
  }
  
  .plugin-single-line {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .plugin-group {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
    width: 100%;
  }
  
  .plugin-buttons {
    width: 100%;
    justify-content: flex-start;
    flex-wrap: wrap;
  }
}

@media (max-width: 480px) {
  .plugin-single-line .el-button {
    font-size: 10px;
    padding: 3px 6px;
    height: 22px;
  }
  
  .plugin-single-line .el-tag {
    font-size: 8px;
    height: 12px;
    line-height: 10px;
    padding: 0 2px;
  }
}

/* 版本提示样式 */
.version-tip {
  margin-top: 4px;
  font-style: italic;
}

/* 私有仓库配置样式 */
.registry-config :deep(.el-form-item__label) {
  margin-bottom: 8px;
}

.private-registry-config {
  margin-top: 12px;
  padding: 12px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.private-registry-config .el-input :deep(.el-input-group__prepend) {
  background: rgba(103, 126, 234, 0.1);
  color: #606266;
  font-weight: 500;
  font-size: 12px;
  border-color: rgba(103, 126, 234, 0.2);
}

.private-registry-config .form-row {
  margin-bottom: 8px;
}

.private-registry-config .form-row:last-child {
  margin-bottom: 0;
}

</style>
