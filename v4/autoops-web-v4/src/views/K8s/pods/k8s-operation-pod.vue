<template>
  <div class="operation-pod">
    <!-- 操作按钮组 -->
    <div class="operation-buttons">
      <!-- 第一行：4个按钮 -->
      <div class="button-row">
        <el-tooltip :content="canScale ? '伸缩' : '不支持伸缩'" placement="top">
          <el-button
            type="primary"
            size="small"
            circle
            :disabled="!canScale"
            @click="scaleWorkload"
          >
            <img src="@/assets/image/伸缩.svg" alt="伸缩" width="14" height="14" style="filter: brightness(0) invert(1);" />
          </el-button>
        </el-tooltip>
        
        <el-tooltip :content="canRestart ? '重启' : '不支持重启'" placement="top">
          <el-button
            type="warning"
            :icon="RefreshRight"
            size="small"
            circle
            :disabled="!canRestart"
            @click="restartWorkload"
          />
        </el-tooltip>
        
        <el-tooltip content="回滚" placement="top">
          <el-button
            type="info"
            :icon="Document"
            size="small"
            circle
            @click="handleRollback"
          />
        </el-tooltip>
        
        <el-tooltip content="编辑YAML" placement="top">
          <el-button
            type="primary"
            :icon="Edit"
            size="small"
            circle
            @click="handleEditWorkloadYaml"
          />
        </el-tooltip>
      </div>
      
      <!-- 第二行：4个按钮 -->
      <div class="button-row">
        <el-tooltip content="更新Pod配置" placement="top">
          <el-button
            type="success"
            :icon="Setting"
            size="small"
            circle
            @click="handleUpdatePodConfig"
          />
        </el-tooltip>
        
        <el-tooltip content="更新调度" placement="top">
          <el-button
            type="info"
            :icon="Monitor"
            size="small"
            circle
            @click="handleUpdateScheduling"
          />
        </el-tooltip>
        
        <el-tooltip content="更新工作负载" placement="top">
          <el-button
            type="success"
            :icon="RefreshRight"
            size="small"
            circle
            @click="editWorkload"
          />
        </el-tooltip>
        
        <el-tooltip content="删除工作负载" placement="top">
          <el-button
            type="danger"
            :icon="Delete"
            size="small"
            circle
            @click="deleteWorkload"
          />
        </el-tooltip>
      </div>
    </div>

<!-- 日志对话框 -->
    <el-dialog
      v-model="logDialogVisible"
      :title="`Pod日志 - ${currentPod.name || ''}`"
      width="80%"
      class="log-dialog"
    >
      <div class="log-controls">
        <el-form inline>
          <el-form-item label="容器">
            <el-select v-model="logParams.container" style="width: 200px;">
              <el-option 
                v-for="container in currentPod.containers" 
                :key="container.name" 
                :label="container.name" 
                :value="container.name"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="行数">
            <el-input-number v-model="logParams.lines" :min="10" :max="1000" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="refreshPodLogs">刷新日志</el-button>
            <el-button @click="copyToClipboard(currentPodLogs, '日志已复制')">复制日志</el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="log-content">
        <pre v-loading="logLoading">{{ currentPodLogs || '暂无日志数据' }}</pre>
      </div>
    </el-dialog>

    <!-- YAML查看对话框 -->
    <el-dialog
      v-model="yamlDialogVisible"
      :title="`Pod YAML - ${currentPod.name || ''}`"
      width="80%"
      class="yaml-dialog"
    >
      <div class="yaml-controls">
        <el-button @click="copyToClipboard(currentPodYaml, 'YAML已复制')">复制YAML</el-button>
      </div>
      <div class="yaml-content">
        <CodeEditor 
          v-model="currentPodYaml" 
          language="yaml" 
          height="500px"
          :readonly="true"
        />
      </div>
    </el-dialog>

    <!-- 扩缩容对话框 -->
    <el-dialog
      v-model="scaleDialogVisible"
      :title="`扩缩容 - ${workload.name || ''}`"
      width="400px"
      class="scale-dialog"
    >
      <el-form>
        <el-form-item label="当前副本数">
          <el-input :value="workload.totalReplicas" readonly />
        </el-form-item>
        <el-form-item label="目标副本数">
          <el-input-number 
            v-model="scaleForm.replicas" 
            :min="0" 
            :max="100" 
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="scaleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitScale" :loading="scaleLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  View,
  Edit,
  Delete,
  Setting,
  RefreshRight,
  Monitor,
  Document
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import CodeEditor from '@/components/CodeEditor.vue'
import { getYamlContentFromApiResponse } from '../utils/yaml'

// Props
const props = defineProps({
  workload: {
    type: Object,
    required: true
  },
  clusterId: {
    type: String,
    required: true
  },
  namespace: {
    type: String,
    required: true
  }
})

// Emits
const emit = defineEmits([
  'refresh', 
  'show-monitoring',
  'scale-workload',
  'restart-workload', 
  'update-pod-config',
  'update-scheduling',
  'edit-workload-yaml',
  'delete-workload',
  'rollback-workload'
])

// 状态变量
const logDialogVisible = ref(false)
const yamlDialogVisible = ref(false)
const scaleDialogVisible = ref(false)

const logLoading = ref(false)
const scaleLoading = ref(false)

const currentPod = ref({})
const currentPodLogs = ref('')
const currentPodYaml = ref('')

// 表单数据
const logParams = reactive({
  container: '',
  lines: 100,
  follow: false
})

const scaleForm = reactive({
  replicas: 1
})

// 计算属性
const canScale = computed(() => {
  const scalableTypes = ['deployment', 'statefulset', 'replicaset']
  return scalableTypes.includes(props.workload.type?.toLowerCase())
})

const canRestart = computed(() => {
  const restartableTypes = ['deployment', 'statefulset', 'daemonset']
  return restartableTypes.includes(props.workload.type?.toLowerCase())
})


// 工具函数
const getPodStatusTag = (status) => {
  const statusMap = {
    'Running': 'success',
    'Pending': 'warning', 
    'Succeeded': 'success',
    'Failed': 'danger',
    'Unknown': 'info',
    'CrashLoopBackOff': 'danger',
    'Error': 'danger'
  }
  return statusMap[status] || 'info'
}

const getPodStatusText = (status) => {
  const textMap = {
    'Running': '运行中',
    'Pending': '等待中',
    'Succeeded': '成功',
    'Failed': '失败',
    'Unknown': '未知',
    'CrashLoopBackOff': '崩溃重启',
    'Error': '错误'
  }
  return textMap[status] || status
}

const copyToClipboard = async (text, successMessage = '已复制到剪贴板') => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(successMessage)
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败')
  }
}

// 业务方法
const getWorkloadLabelSelector = (workload) => {
  if (workload.labels && workload.labels.app) {
    return `app=${workload.labels.app}`
  }
  if (workload.labels && workload.labels['app.kubernetes.io/name']) {
    return `app.kubernetes.io/name=${workload.labels['app.kubernetes.io/name']}`
  }
  return `app=${workload.name}`
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

const viewPodLogs = async (pod) => {
  try {
    currentPod.value = pod
    logParams.container = pod.containers?.[0]?.name || ''
    
    logLoading.value = true
    const response = await k8sApi.getPodLogs(props.clusterId, props.namespace, pod.name, {
      container: logParams.container,
      tailLines: logParams.lines
    })
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      currentPodLogs.value = responseData.data || '暂无日志'
      logDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取Pod日志失败')
    }
  } catch (error) {
    console.error('获取Pod日志失败:', error)
    ElMessage.error('获取Pod日志失败，请检查网络连接')
  } finally {
    logLoading.value = false
  }
}

const refreshPodLogs = async () => {
  if (currentPod.value.name) {
    await viewPodLogs(currentPod.value)
  }
}

const viewPodYaml = async (pod) => {
  try {
    const response = await k8sApi.getPodYaml(props.clusterId, props.namespace, pod.name)
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      currentPodYaml.value = getYamlContentFromApiResponse(
        responseData,
        `apiVersion: v1\nkind: Pod\nmetadata:\n  name: ${pod.name}`
      )
      currentPod.value = pod
      yamlDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取Pod YAML失败')
    }
  } catch (error) {
    console.error('获取Pod YAML失败:', error)
    ElMessage.error('获取Pod YAML失败，请检查网络连接')
  }
}

const deletePod = async (pod) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除Pod "${pod.name}" 吗？`,
      '删除Pod确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const response = await k8sApi.deletePod(props.clusterId, props.namespace, pod.name)
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`Pod ${pod.name} 删除成功`)
      // 重新获取Pod列表
      await viewPodList()
    } else {
      ElMessage.error(responseData.message || '删除Pod失败')
    }
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消删除操作')
    } else {
      console.error('删除Pod失败:', error)
      ElMessage.error('删除Pod失败，请检查网络连接')
    }
  }
}

const editWorkload = () => {
  ElMessage.info('编辑功能开发中...')
}

const scaleWorkload = () => {
  if (!canScale.value) {
    ElMessage.warning('该工作负载不支持扩缩容操作')
    return
  }
  emit('scale-workload', props.workload)
  
  scaleForm.replicas = props.workload.totalReplicas || 1
  scaleDialogVisible.value = true
}

const submitScale = async () => {
  try {
    scaleLoading.value = true
    const response = await k8sApi.scaleDeployment(
      props.clusterId, 
      props.namespace, 
      props.workload.name, 
      { replicas: scaleForm.replicas }
    )
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`${props.workload.name} 扩缩容成功`)
      scaleDialogVisible.value = false
      emit('refresh')
    } else {
      ElMessage.error(responseData.message || '扩缩容失败')
    }
  } catch (error) {
    console.error('扩缩容失败:', error)
    ElMessage.error('扩缩容失败，请检查网络连接')
  } finally {
    scaleLoading.value = false
  }
}

const restartWorkload = async () => {
  if (!canRestart.value) {
    ElMessage.warning('该工作负载不支持重启操作')
    return
  }
  emit('restart-workload', props.workload)
  
  try {
    await ElMessageBox.confirm(
      `确定要重启 ${props.workload.type} "${props.workload.name}" 吗？`,
      '重启确认',
      {
        confirmButtonText: '确定重启',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const response = await k8sApi.restartDeployment(props.clusterId, props.namespace, props.workload.name)
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`${props.workload.name} 重启成功`)
      emit('refresh')
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

const deleteWorkload = async () => {
  emit('delete-workload', props.workload)
}

const handleUpdatePodConfig = () => {
  emit('update-pod-config', props.workload)
}

const handleUpdateScheduling = () => {
  emit('update-scheduling', props.workload)
}

const handleEditWorkloadYaml = () => {
  emit('edit-workload-yaml', props.workload)
}

const handleRollback = () => {
  emit('rollback-workload', props.workload)
}

const originalDeleteWorkload = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 ${props.workload.type} "${props.workload.name}" 吗？此操作不可撤销！`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    ElMessage.info('删除功能开发中...')
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消删除操作')
    }
  }
}

// 暴露方法给父组件
defineExpose({
  viewPodList
})
</script>

<style scoped>
.operation-buttons {
  display: flex;
  gap: 8px;
}

.pod-stats {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 6px;
  margin-bottom: 16px;
}

.empty-state {
  padding: 40px;
  text-align: center;
}

.log-controls {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 6px;
  margin-bottom: 16px;
}

.log-content {
  max-height: 500px;
  overflow-y: auto;
  background: #1e1e1e;
  color: #ffffff;
  padding: 16px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.4;
}

.log-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.yaml-controls {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 6px;
  margin-bottom: 16px;
}

.yaml-content {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

/* 按钮行布局样式 */
.operation-buttons {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.button-row {
  display: flex;
  gap: 6px;
  justify-content: center;
}

.button-row .el-button {
  width: 32px;
  height: 32px;
  transition: all 0.3s ease;
}

.button-row .el-button:hover:not(.is-disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.button-row .el-button.is-disabled {
  cursor: not-allowed;
  opacity: 0.5;
  background-color: #f5f7fa !important;
  border-color: #e4e7ed !important;
  color: #c0c4cc !important;
}

.button-row .el-button.is-disabled:hover {
  transform: none;
  box-shadow: none;
}
</style>
