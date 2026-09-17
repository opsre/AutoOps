<template>
  <div class="create-pod">
    <!-- 创建Pod对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="通过YAML创建Pod"
      width="50%"
      class="create-pod-dialog"
      :close-on-click-modal="false"
    >
      <div class="create-pod-content">
        <!-- 操作按钮 -->
        <div class="pod-controls">
          <div class="control-group">
            <el-button 
              :icon="View" 
              type="warning" 
              @click="validatePodYaml" 
              :loading="loading"
            >
              校验YAML
            </el-button>
            <el-button 
              :icon="Search" 
              type="info" 
              @click="previewPodCreation" 
              :loading="loading"
            >
              DryRun预览
            </el-button>
            <el-button 
              :icon="CopyDocument" 
              @click="copyToClipboard(createForm.yamlContent, 'YAML已复制')"
            >
              复制YAML
            </el-button>
          </div>
          <div class="status-indicators">
            <el-tag 
              v-if="validationResult" 
              :type="validationResult.valid ? 'success' : 'danger'"
              size="small"
            >
              {{ validationResult.valid ? '✓ 格式正确' : '✗ 格式错误' }}
            </el-tag>
            <el-tag 
              v-if="dryRunResult" 
              :type="dryRunResult.success ? 'success' : 'danger'"
              size="small"
            >
              {{ dryRunResult.success ? '✓ 可以创建' : '✗ 创建失败' }}
            </el-tag>
          </div>
        </div>

        <!-- YAML编辑器 -->
        <div class="yaml-editor-container">
          <CodeEditor 
            ref="yamlEditor"
            v-model="createForm.yamlContent" 
            language="yaml" 
            height="500px"
            fontSize="14px"
            :readonly="false"
            placeholder="请输入Pod的YAML配置..."
          />
        </div>

        <!-- 验证结果显示 -->
        <div v-if="validationResult && !validationResult.valid" class="validation-result error">
          <el-alert
            :title="validationResult.message"
            type="error"
            :description="validationResult.details"
            show-icon
            :closable="false"
          />
        </div>

        <div v-if="validationResult && validationResult.valid" class="validation-result success">
          <el-alert
            :title="validationResult.message"
            type="success"
            show-icon
            :closable="false"
          />
        </div>

        <!-- DryRun结果显示 -->
        <div v-if="dryRunResult" class="dryrun-result">
          <el-alert
            :title="dryRunResult.message"
            :type="dryRunResult.success ? 'success' : 'error'"
            :description="dryRunResult.details"
            show-icon
            :closable="false"
          />
        </div>

        <!-- Pod创建选项 -->
        <div class="create-options">
          <el-card class="options-card">
            <template #header>
              <span>创建选项</span>
            </template>
            <el-descriptions :column="2" size="small">
              <el-descriptions-item label="目标集群">
                {{ clusterName || '未选择' }}
              </el-descriptions-item>
              <el-descriptions-item label="目标命名空间">
                {{ namespace || '未选择' }}
              </el-descriptions-item>
            </el-descriptions>
          </el-card>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="createDialogVisible = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="createPod" 
            :loading="loading"
            :disabled="validationResult && !validationResult.valid"
          >
            创建Pod
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { View, Search, CopyDocument } from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import CodeEditor from '@/components/CodeEditor.vue'

// Props
const props = defineProps({
  clusterId: {
    type: String,
    required: true
  },
  namespace: {
    type: String,
    required: true
  },
  clusterName: {
    type: String,
    default: ''
  }
})

// Emits
const emit = defineEmits(['created', 'refresh'])

// 创建Pod相关状态
const createDialogVisible = ref(false)
const loading = ref(false)
const yamlEditor = ref(null)

// 创建表单
const createForm = reactive({
  yamlContent: `apiVersion: v1
kind: Pod
metadata:
  name: example-pod
  namespace: default
  labels:
    app: example
    created-by: k8s-web-ui
spec:
  restartPolicy: Always
  containers:
  - name: example-container
    image: nginx:latest
    ports:
    - containerPort: 80
      name: http
      protocol: TCP
    env:
    - name: ENVIRONMENT
      value: "production"
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
    livenessProbe:
      httpGet:
        path: /
        port: 80
      initialDelaySeconds: 30
      periodSeconds: 10
    readinessProbe:
      httpGet:
        path: /
        port: 80
      initialDelaySeconds: 5
      periodSeconds: 5
  # 可选：添加持久化存储卷
  # volumes:
  # - name: app-storage
  #   emptyDir: {}`,
  dryRun: false,
  validateOnly: false
})

const validationResult = ref(null)
const dryRunResult = ref(null)

// 复制到剪贴板
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
    } catch (err) {
      ElMessage.error('复制失败')
    }
    document.body.removeChild(textArea)
  }
}

// 显示创建对话框
const showCreateDialog = () => {
  if (!props.clusterId) {
    ElMessage.warning('请先选择集群')
    return
  }
  if (!props.namespace) {
    ElMessage.warning('请先选择命名空间')
    return
  }
  
  // 重置表单和结果
  validationResult.value = null
  dryRunResult.value = null
  
  // 生成智能的Pod名称
  const timestamp = new Date().toISOString().replace(/[:.-]/g, '').toLowerCase().substr(0, 15)
  const podName = `custom-pod-${timestamp}`
  
  // 更新YAML内容，包括命名空间和Pod名称
  let updatedYaml = createForm.yamlContent.replace(
    /namespace:\s*\w+/,
    `namespace: ${props.namespace}`
  ).replace(
    /name:\s*[\w-]+/,
    `name: ${podName}`
  )
  
  createForm.yamlContent = updatedYaml
  
  createDialogVisible.value = true
  
  // 对话框打开后聚焦编辑器
  nextTick(() => {
    if (yamlEditor.value && yamlEditor.value.focus) {
      yamlEditor.value.focus()
    }
  })
}

// YAML校验函数
const validatePodYaml = async () => {
  if (!createForm.yamlContent.trim()) {
    ElMessage.warning('YAML内容不能为空')
    return
  }
  
  try {
    loading.value = true
    validationResult.value = null
    
    const response = await k8sApi.validateYaml(props.clusterId, {
      yamlContent: createForm.yamlContent,
      resourceType: 'pod'
    })
    
    const responseData = response.data || response
    
    if (responseData.code === 200) {
      validationResult.value = {
        valid: true,
        message: responseData.message || 'YAML格式验证通过',
        details: responseData.data
      }
      ElMessage.success('YAML格式验证通过')
    } else {
      validationResult.value = {
        valid: false,
        message: responseData.message || 'YAML格式验证失败',
        details: responseData.data
      }
      ElMessage.error(responseData.message || 'YAML格式验证失败')
    }
  } catch (error) {
    console.error('YAML校验失败:', error)
    validationResult.value = {
      valid: false,
      message: error.message || 'YAML格式校验失败',
      details: null
    }
    ElMessage.error('YAML格式校验失败: ' + (error.message || '网络错误'))
  } finally {
    loading.value = false
  }
}

// DryRun预览函数
const previewPodCreation = async () => {
  if (!createForm.yamlContent.trim()) {
    ElMessage.warning('YAML内容不能为空')
    return
  }
  
  try {
    loading.value = true
    dryRunResult.value = null
    
    const response = await k8sApi.createPodFromYaml(props.clusterId, props.namespace, {
      yamlContent: createForm.yamlContent,
      dryRun: true,
      validateOnly: false
    })
    
    const responseData = response.data || response
    
    if (responseData.code === 200) {
      dryRunResult.value = {
        success: true,
        message: responseData.message || 'Pod创建预览成功',
        details: responseData.data
      }
      ElMessage.success('Pod创建预览成功，可以正常创建')
    } else {
      dryRunResult.value = {
        success: false,
        message: responseData.message || 'Pod创建预览失败',
        details: responseData.data
      }
      ElMessage.error(responseData.message || 'Pod创建预览失败')
    }
  } catch (error) {
    console.error('DryRun预览失败:', error)
    dryRunResult.value = {
      success: false,
      message: error.message || 'Pod创建预览失败',
      details: null
    }
    ElMessage.error('Pod创建预览失败: ' + (error.message || '网络错误'))
  } finally {
    loading.value = false
  }
}

// 创建Pod函数
const createPod = async () => {
  if (!createForm.yamlContent.trim()) {
    ElMessage.warning('YAML内容不能为空')
    return
  }
  
  try {
    loading.value = true
    
    const response = await k8sApi.createPodFromYaml(props.clusterId, props.namespace, {
      yamlContent: createForm.yamlContent,
      dryRun: false,
      validateOnly: false
    })
    
    const responseData = response.data || response
    
    if (responseData.code === 200) {
      ElMessage.success('Pod创建成功!')
      createDialogVisible.value = false
      
      // 通知父组件刷新列表
      emit('created')
      emit('refresh')
    } else {
      ElMessage.error(responseData.message || 'Pod创建失败')
    }
  } catch (error) {
    console.error('Pod创建失败:', error)
    ElMessage.error('Pod创建失败: ' + (error.message || '网络错误'))
  } finally {
    loading.value = false
  }
}

// 监听对话框打开状态，自动聚焦编辑器
watch(createDialogVisible, (newVal) => {
  if (newVal) {
    nextTick(() => {
      if (yamlEditor.value && yamlEditor.value.focus) {
        yamlEditor.value.focus()
      }
    })
  }
})

// 暴露方法给父组件
defineExpose({
  showCreateDialog
})
</script>

<style scoped>
/* 创建Pod对话框样式 */
.create-pod-dialog .el-dialog__body {
  padding: 20px;
}

.create-pod-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.pod-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f8f9fa;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.control-group {
  display: flex;
  gap: 8px;
}

.status-indicators {
  display: flex;
  gap: 8px;
  align-items: center;
}

.yaml-editor-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

.validation-result,
.dryrun-result {
  margin-top: 12px;
}

.validation-result.error .el-alert,
.dryrun-result .el-alert--error {
  border-left: 4px solid #f56c6c;
}

.validation-result.success .el-alert,
.dryrun-result .el-alert--success {
  border-left: 4px solid #67c23a;
}

.create-options {
  margin-top: 16px;
}

.options-card {
  border: 1px solid #e4e7ed;
}

.options-card .el-card__header {
  padding: 12px 16px;
  background: #fafafa;
  border-bottom: 1px solid #e4e7ed;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
