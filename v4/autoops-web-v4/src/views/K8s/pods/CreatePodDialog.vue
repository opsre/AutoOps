<template>
  <el-dialog
    :model-value="visible"
    title="通过YAML创建工作负载"
    width="80%"
    class="create-pod-dialog"
    :close-on-click-modal="false"
    @update:model-value="handleClose"
  >
    <div class="create-pod-content">
      <!-- 操作按钮和目标信息 -->
      <div class="pod-controls">
        <div class="control-group">
          <el-tooltip content="YAML创建校验" placement="top">
            <el-button
              :icon="Search"
              type="primary"
              size="small"
              circle
              @click="previewPodCreation"
              :loading="loading"
              class="validate-btn"
            />
          </el-tooltip>
          <el-tooltip content="复制YAML" placement="top">
            <el-button
              :icon="DocumentCopy"
              type="success"
              size="small"
              circle
              @click="copyToClipboard"
              class="copy-btn"
            />
          </el-tooltip>
          <!-- 目标信息 -->
          <div class="target-info">
            <div class="target-item">
              <span class="target-label">目标集群:</span>
              <el-tag type="info" size="small">{{ clusterName || '未选择' }}</el-tag>
            </div>
            <div class="target-item">
              <span class="target-label">目标命名空间:</span>
              <el-tag type="warning" size="small">{{ namespace || '未选择' }}</el-tag>
            </div>
          </div>
        </div>
        <div class="status-indicators">
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
          v-model="yamlContent"
          language="yaml"
          height="500px"
          fontSize="14px"
          :readonly="false"
          placeholder="请输入工作负载的YAML配置..."
        />
      </div>

      <!-- YAML创建校验结果显示 -->
      <div v-if="dryRunResult" class="dryrun-result">
        <el-alert
          :title="dryRunResult.message"
          :type="dryRunResult.success ? 'success' : 'error'"
          :description="dryRunResult.details"
          show-icon
          :closable="false"
        />
      </div>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          @click="handleSubmit"
          :loading="loading"
        >
          创建工作负载
        </el-button>
      </div>
    </div>

  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch, nextTick, defineProps, defineEmits } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, DocumentCopy } from '@element-plus/icons-vue'
import CodeEditor from '@/components/CodeEditor.vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  clusterId: {
    type: String,
    default: ''
  },
  clusterName: {
    type: String,
    default: ''
  },
  namespace: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:visible', 'close', 'preview', 'create'])

const yamlEditor = ref(null)
const loading = ref(false)
const dryRunResult = ref(null)

// 默认的Deployment YAML模板
const defaultPodYaml = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: custom-deployment-\${timestamp}
  namespace: \${namespace}
  labels:
    app: custom-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: custom-deployment
  template:
    metadata:
      labels:
        app: custom-deployment
    spec:
      containers:
      - name: main
        image: nginx:latest
        ports:
        - containerPort: 80
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi`

const yamlContent = ref('')

// 初始化YAML内容
const initializeYaml = () => {
  const timestamp = new Date().getTime()

  let initialYaml = defaultPodYaml
    .replace('${timestamp}', timestamp)
    .replace('${namespace}', props.namespace || 'default')

  yamlContent.value = initialYaml
}

// 监听对话框打开状态
watch(() => props.visible, (newVal) => {
  if (newVal) {
    initializeYaml()
    // 清空之前的校验结果
    dryRunResult.value = null

    // 聚焦编辑器
    nextTick(() => {
      if (yamlEditor.value && yamlEditor.value.focus) {
        yamlEditor.value.focus()
      }
    })
  }
})

const handleClose = () => {
  emit('update:visible', false)
  emit('close')
}

// YAML创建校验
const previewPodCreation = async () => {
  if (!yamlContent.value.trim()) {
    ElMessage.warning('YAML内容不能为空')
    return
  }

  emit('preview', {
    yamlContent: yamlContent.value,
    dryRun: true
  })
}

// 复制到剪贴板
const copyToClipboard = async () => {
  console.log('Copy button clicked')
  try {
    console.log('YAML content to copy:', yamlContent.value ? yamlContent.value.substring(0, 100) + '...' : 'empty')

    if (!yamlContent.value || typeof yamlContent.value !== 'string' || yamlContent.value.trim() === '') {
      ElMessage.warning('没有内容可以复制')
      return
    }

    await navigator.clipboard.writeText(yamlContent.value)
    ElMessage.success('YAML已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    // 备用方案：使用旧方法复制
    try {
      const textArea = document.createElement('textarea')
      textArea.value = yamlContent.value
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      ElMessage.success('YAML已复制到剪贴板')
    } catch (fallbackError) {
      console.error('备用复制方案也失败:', fallbackError)
      ElMessage.error('复制失败，请手动复制')
    }
  }
}

// 提交创建
const handleSubmit = async () => {
  if (!yamlContent.value.trim()) {
    ElMessage.warning('YAML内容不能为空')
    return
  }

  loading.value = true

  try {
    emit('create', {
      yamlContent: yamlContent.value,
      dryRun: false
    })
  } finally {
    loading.value = false
  }
}

// 暴露方法供父组件调用
defineExpose({
  setDryRunResult: (result) => {
    dryRunResult.value = result
  },
  setLoading: (state) => {
    loading.value = state
  }
})
</script>

<style scoped>
.create-pod-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.pod-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.control-group {
  display: flex;
  gap: 20px;
  align-items: center;
  justify-content: flex-start;
  flex-wrap: wrap;
}

.target-info {
  display: flex;
  gap: 16px;
  margin-left: 12px;
}

.target-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.target-label {
  font-size: 13px;
  color: #6b7280;
  font-weight: 500;
  white-space: nowrap;
}

.status-indicators {
  display: flex;
  gap: 8px;
}

/* 自定义按钮样式 */
.validate-btn {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  border: none;
  color: white;
  box-shadow: 0 2px 8px rgba(79, 172, 254, 0.3);
  transition: all 0.3s ease;
}

.validate-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.4);
  background: linear-gradient(135deg, #00f2fe 0%, #4facfe 100%);
}

.copy-btn {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  border: none;
  color: #6b9bd1;
  box-shadow: 0 2px 8px rgba(168, 237, 234, 0.3);
  transition: all 0.3s ease;
}

.copy-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(168, 237, 234, 0.4);
  background: linear-gradient(135deg, #fed6e3 0%, #a8edea 100%);
}

/* Loading状态样式 */
.validate-btn.is-loading {
  background: linear-gradient(135deg, #94a3b8 0%, #64748b 100%);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.7; }
  100% { opacity: 1; }
}

.yaml-editor-container {
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  overflow: hidden;
}

.dryrun-result {
  margin-top: 16px;
}

.action-buttons {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

/* 创建Pod对话框样式 */
.create-pod-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.create-pod-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.create-pod-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.create-pod-dialog :deep(.el-dialog__body) {
  padding: 20px;
}
</style>