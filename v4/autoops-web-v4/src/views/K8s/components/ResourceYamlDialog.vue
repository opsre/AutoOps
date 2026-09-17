<template>
  <el-dialog
    :model-value="visible"
    :title="dialogTitle"
    width="80%"
    class="resource-yaml-dialog"
    :close-on-click-modal="false"
    @update:model-value="handleClose"
  >
    <div class="yaml-content">
      <div class="yaml-controls">
        <div class="control-group">
          <el-tag type="info" size="small">{{ resourceType }}</el-tag>
          <el-tag type="warning" size="small">{{ resourceName }}</el-tag>
        </div>
        <div class="control-group">
          <el-button
            :icon="DocumentCopy"
            type="success"
            size="small"
            circle
            @click="copyToClipboard"
            :disabled="!yamlContent"
          />
        </div>
      </div>

      <div class="yaml-editor-wrapper" v-loading="fetching">
        <div class="yaml-editor-container">
          <div v-if="!editable" class="yaml-viewer">
            <pre class="yaml-text">{{ yamlContent }}</pre>
          </div>
          <CodeEditor
            v-else
            v-model="yamlContent"
            language="yaml"
            height="560px"
            fontSize="14px"
            :readonly="false"
            :placeholder="placeholderText"
          />
        </div>
      </div>

      <div class="action-buttons">
        <el-button @click="handleClose">{{ editable ? '取消' : '关闭' }}</el-button>
        <el-button v-if="editable" type="primary" @click="saveYaml" :loading="saving">保存</el-button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { computed, ref, watch, defineProps, defineEmits } from 'vue'
import { ElMessage } from 'element-plus'
import { DocumentCopy } from '@element-plus/icons-vue'
import CodeEditor from '@/components/CodeEditor.vue'
import { copyYamlToClipboard, getYamlContentFromApiResponse } from '../utils/yaml'

const props = defineProps({
  visible: { type: Boolean, default: false },
  clusterId: { type: [String, Number], required: true },
  resourceType: { type: String, required: true },
  resourceName: { type: String, required: true },
  fetchYaml: { type: Function, required: true },
  updateYaml: { type: Function, default: null },
  editable: { type: Boolean, default: true },
  fallbackYaml: { type: String, default: '' },
  loadErrorMessage: { type: String, default: '获取 YAML 失败' },
  saveSuccessMessage: { type: String, default: 'YAML 更新成功' },
  saveErrorMessage: { type: String, default: '更新 YAML 失败' }
})

const emit = defineEmits(['update:visible', 'saved'])

const yamlContent = ref('')
const fetching = ref(false)
const saving = ref(false)

const dialogTitle = computed(() => `${props.editable ? 'YAML 编辑' : `${props.resourceType} YAML`} - ${props.resourceName}`)
const placeholderText = computed(() => `加载${props.resourceType} YAML 中...`)

const loadYaml = async () => {
  if (!props.clusterId || !props.resourceName) return

  fetching.value = true
  yamlContent.value = ''

  try {
    const res = await props.fetchYaml(props.clusterId, props.resourceName)
    const data = res.data || res

    if (data.code === 200 || data.success) {
      yamlContent.value = getYamlContentFromApiResponse(data, props.fallbackYaml)
    } else {
      ElMessage.error(data.message || props.loadErrorMessage)
    }
  } catch (error) {
    console.error(props.loadErrorMessage, error)
    ElMessage.error(props.loadErrorMessage)
  } finally {
    fetching.value = false
  }
}

watch(() => props.visible, (val) => {
  if (val) loadYaml()
  else yamlContent.value = ''
}, { immediate: true })

const saveYaml = async () => {
  if (!props.editable) {
    return
  }

  if (!yamlContent.value.trim()) {
    ElMessage.warning('YAML 内容不能为空')
    return
  }

  saving.value = true

  try {
    const res = await props.updateYaml(props.clusterId, props.resourceName, yamlContent.value)
    const data = res.data || res

    if (data.code === 200 || data.success) {
      ElMessage.success(props.saveSuccessMessage)
      emit('saved')
      handleClose()
    } else {
      ElMessage.error(data.message || '更新失败')
    }
  } catch (error) {
    console.error(props.saveErrorMessage, error)
    ElMessage.error(props.saveErrorMessage)
  } finally {
    saving.value = false
  }
}

const copyToClipboard = async () => {
  const result = await copyYamlToClipboard(yamlContent.value)

  if (result.success) {
    ElMessage.success('已复制到剪贴板')
    return
  }

  if (result.message === '没有内容可以复制') {
    ElMessage.warning(result.message)
    return
  }

  ElMessage.error(result.message)
}

const handleClose = () => emit('update:visible', false)
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
  padding: 12px 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.control-group {
  display: flex;
  gap: 8px;
  align-items: center;
}

.yaml-editor-wrapper {
  width: 100%;
  height: 560px;
}

.yaml-editor-container {
  height: 100%;
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  overflow: hidden;
}

.yaml-viewer {
  height: 100%;
  background: #f8f9fa;
  overflow: auto;
}

.yaml-text {
  margin: 0;
  padding: 16px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #334155;
  background: transparent;
  white-space: pre-wrap;
  word-break: break-word;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

.resource-yaml-dialog :deep(.el-dialog) {
  border-radius: 16px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.resource-yaml-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.resource-yaml-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.resource-yaml-dialog :deep(.el-dialog__body) {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
  max-height: calc(90vh - 120px);
}
</style>
