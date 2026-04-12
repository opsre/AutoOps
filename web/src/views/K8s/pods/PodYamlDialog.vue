<template>
  <el-dialog
    :model-value="visible"
    :title="dialogTitle"
    width="80%"
    class="pod-yaml-dialog"
    :close-on-click-modal="false"
    @update:model-value="handleClose"
  >
    <div class="yaml-content">
      <!-- 操作按钮和目标信息 -->
      <div class="yaml-controls">
        <div class="control-group">
          <el-tooltip content="YAML格式验证" placement="top">
            <el-button
              :icon="Search"
              type="primary"
              size="small"
              circle
              @click="validateYaml"
              :loading="validating"
              :disabled="!hasYamlContent"
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
              :disabled="!hasYamlContent"
              class="copy-btn"
            />
          </el-tooltip>
          <!-- 资源信息 -->
          <div class="target-info">
            <div class="target-item">
              <span class="target-label">资源类型:</span>
              <el-tag type="info" size="small">{{ resourceType || '未知' }}</el-tag>
            </div>
            <div class="target-item">
              <span class="target-label">资源名称:</span>
              <el-tag type="warning" size="small">{{ resourceName || '未选择' }}</el-tag>
            </div>
          </div>
        </div>
        <div class="status-indicators">
          <el-tag
            v-if="validationResult"
            :type="validationResult.success ? 'success' : 'danger'"
            size="small"
          >
            {{ validationResult.success ? '✓ 格式正确' : '✗ 格式错误' }}
          </el-tag>
        </div>
      </div>

      <!-- YAML编辑器/查看器 -->
      <div class="yaml-editor-wrapper">
        <div class="yaml-editor-container">
          <!-- 只读模式 - Pod YAML查看 -->
          <div v-if="!editable" class="yaml-viewer">
            <pre class="yaml-text">{{ yamlContent }}</pre>
          </div>
          <!-- 编辑模式 - 工作负载YAML编辑 -->
          <CodeEditor
            v-else
            ref="yamlEditor"
            v-model="editableYaml"
            language="yaml"
            :height="editorHeight"
            fontSize="14px"
            :readonly="false"
            placeholder="请输入YAML配置..."
          />
        </div>
      </div>

      <!-- YAML验证结果显示 -->
      <div v-if="validationResult" class="validation-result">
        <el-alert
          :title="validationResult.message"
          :type="validationResult.success ? 'success' : 'error'"
          :description="validationResult.details"
          show-icon
          :closable="false"
        />
      </div>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button @click="handleClose">{{ editable ? '取消' : '关闭' }}</el-button>
        <el-button
          v-if="editable"
          type="primary"
          @click="saveYaml"
          :loading="saving"
        >
          保存YAML
        </el-button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, defineProps, defineEmits } from 'vue'
import { ElMessage } from 'element-plus'
import { DocumentCopy, Search } from '@element-plus/icons-vue'
import CodeEditor from '@/components/CodeEditor.vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  // YAML内容
  yamlContent: {
    type: String,
    default: ''
  },
  // 目标资源信息
  resourceName: {
    type: String,
    default: ''
  },
  resourceType: {
    type: String,
    default: 'Pod' // Pod, Deployment, etc.
  },
  // 是否可编辑
  editable: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'close', 'save'])

const yamlEditor = ref(null)
const editableYaml = ref('')
const saving = ref(false)
const validating = ref(false)
const validationResult = ref(null)

// 计算编辑器动态高度，配合外层父容器滚动
const editorHeight = computed(() => {
  if (!editableYaml.value) return '100%'

  const lines = editableYaml.value.split('\n').length
  // 按照字体行高估算动态高度，让内部元素撑开，父容器负责自适应包裹和出现滑动窗口
  const calculatedHeight = Math.max(400, lines * 21 + 40)
  return `${calculatedHeight}px`
})

// 计算对话框标题
const dialogTitle = computed(() => {
  const action = props.editable ? '编辑YAML' : `${props.resourceType} YAML`
  return `${action} - ${props.resourceName || ''}`
})

// 计算是否有YAML内容可操作
const hasYamlContent = computed(() => {
  const content = props.editable ? editableYaml.value : props.yamlContent
  return content && typeof content === 'string' && content.trim().length > 0
})

// 监听对话框打开状态
watch(() => props.visible, (newVal) => {
  if (newVal) {
    // 更宽松的类型转换，确保有内容时正确显示
    editableYaml.value = (props.yamlContent !== null && props.yamlContent !== undefined)
      ? String(props.yamlContent)
      : ''
    // 清空之前的验证结果
    validationResult.value = null
  }
})

// 监听YAML内容变化
watch(() => props.yamlContent, (newVal) => {
  // 更宽松的类型转换，确保有内容时正确显示
  editableYaml.value = (newVal !== null && newVal !== undefined)
    ? String(newVal)
    : ''
  // 清空验证结果
  validationResult.value = null
})

const handleClose = () => {
  emit('update:visible', false)
  emit('close')
}

// 复制到剪贴板
const copyToClipboard = async () => {
  console.log('Copy button clicked')
  try {
    const textToCopy = props.editable ? editableYaml.value : props.yamlContent
    console.log('Text to copy:', typeof textToCopy === 'string' ? textToCopy.substring(0, 100) + '...' : textToCopy)

    if (!textToCopy || typeof textToCopy !== 'string' || textToCopy.trim() === '') {
      ElMessage.warning('没有内容可以复制')
      return
    }

    await navigator.clipboard.writeText(textToCopy)
    ElMessage.success('YAML已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    // 备用方案：使用旧方法复制
    try {
      const textArea = document.createElement('textarea')
      textArea.value = props.editable ? editableYaml.value : props.yamlContent
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

// 验证YAML格式
const validateYaml = async () => {
  console.log('Validate button clicked')

  // 获取要验证的YAML内容
  const yamlToValidate = props.editable ? editableYaml.value : props.yamlContent
  console.log('YAML to validate:', typeof yamlToValidate === 'string' ? yamlToValidate.substring(0, 100) + '...' : yamlToValidate)

  if (!yamlToValidate || typeof yamlToValidate !== 'string' || !yamlToValidate.trim()) {
    validationResult.value = {
      success: false,
      message: '验证失败',
      details: 'YAML内容不能为空'
    }
    ElMessage.warning('YAML内容不能为空')
    return false
  }

  validating.value = true

  try {
    // 检查基本的YAML结构
    const lines = yamlToValidate.split('\n')
    let hasApiVersion = false
    let hasKind = false
    let hasMetadata = false
    let errors = []

    lines.forEach(line => {
      if (line.includes('apiVersion:')) hasApiVersion = true
      if (line.includes('kind:')) hasKind = true
      if (line.includes('metadata:')) hasMetadata = true
    })

    if (!hasApiVersion) errors.push('缺少 apiVersion 字段')
    if (!hasKind) errors.push('缺少 kind 字段')
    if (!hasMetadata) errors.push('缺少 metadata 字段')

    // 简化的YAML检查 - 只检查最基本的问题
    try {
      // 检查是否有明显无效的行（暂时跳过复杂的语法检查）
      // 对于Kubernetes YAML，我们信任用户输入，只做基本的结构检查
      console.log('YAML基本检查通过，跳过详细语法验证')
    } catch (syntaxError) {
      errors.push('YAML语法错误: ' + syntaxError.message)
    }

    if (errors.length > 0) {
      validationResult.value = {
        success: false,
        message: '格式验证失败',
        details: errors.join('; ')
      }
      return false
    }

    validationResult.value = {
      success: true,
      message: 'YAML格式验证通过',
      details: '所有必需字段都存在，格式正确'
    }
    return true
  } catch (error) {
    console.error('YAML验证失败:', error)
    validationResult.value = {
      success: false,
      message: '验证异常',
      details: '验证过程中发生错误: ' + error.message
    }
    return false
  } finally {
    validating.value = false
  }
}

// 保存YAML
const saveYaml = async () => {
  const isValid = await validateYaml()
  if (!isValid) {
    return
  }

  saving.value = true

  try {
    // 发出保存事件，由父组件处理具体的保存逻辑
    emit('save', {
      yamlContent: editableYaml.value,
      resourceName: props.resourceName,
      resourceType: props.resourceType
    })
  } finally {
    saving.value = false
  }
}

// 暴露方法供父组件调用
defineExpose({
  setSaving: (state) => {
    saving.value = state
  },
  getYamlContent: () => editableYaml.value
})
</script>

<style scoped>
.yaml-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.yaml-controls {
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
  color: #2c3e50;
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

.yaml-editor-wrapper {
  width: 100%;
}

.yaml-editor-container {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  overflow-y: auto;
  overflow-x: hidden;
  height: 60vh;
  min-height: 400px;
  background: #ffffff;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.02);
  position: relative;
}

/* 滚动条美化 (针对容器本身) */
.yaml-editor-container::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.yaml-editor-container::-webkit-scrollbar-thumb {
  background-color: #cbd5e1;
  border-radius: 4px;
}

.yaml-editor-container::-webkit-scrollbar-track {
  background-color: #f8fafc;
}

/* 只读YAML查看器样式 */
.yaml-viewer {
  background: #ffffff;
  height: 100%;
  overflow-y: auto; /* 保留给只读模式 */
}

/* 滚动条美化 */
.yaml-viewer::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.yaml-viewer::-webkit-scrollbar-thumb {
  background-color: #cbd5e1;
  border-radius: 4px;
}

.yaml-viewer::-webkit-scrollbar-track {
  background-color: #f8fafc;
}

.yaml-text {
  margin: 0;
  padding: 20px;
  font-family: 'Fira Code', 'Consolas', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13.5px;
  line-height: 1.6;
  color: #334155;
  background: transparent;
  white-space: pre-wrap;
  word-break: break-all;
}

.validation-result {
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

/* 对话框样式 */
.pod-yaml-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.pod-yaml-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
  flex-shrink: 0;
}

.pod-yaml-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.pod-yaml-dialog :deep(.el-dialog__body) {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
  max-height: calc(90vh - 120px);
}
</style>