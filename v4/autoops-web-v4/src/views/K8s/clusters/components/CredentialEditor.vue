<template>
  <div class="credential-editor-container">
    <div class="credential-editor">
      <CodeEditor
        v-model="credentialContent"
        language="yaml"
        :height="height"
        class="editor"
      />
    </div>
    <div class="credential-actions">
      <input
        ref="fileInputRef"
        type="file"
        accept="*"
        style="display: none"
        @change="handleFileChange"
      />
      <el-button 
        size="small" 
        type="success" 
        @click="triggerFileUpload"
        :disabled="disabled"
        class="upload-btn"
      >
        <el-icon><Upload /></el-icon>
        上传配置文件
      </el-button>
      <el-button 
        size="small" 
        type="warning" 
        @click="clearCredential" 
        plain
        :disabled="disabled"
      >
        <el-icon><Delete /></el-icon>
        清空
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, Delete } from '@element-plus/icons-vue'
import CodeEditor from '@/components/CodeEditor.vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  height: {
    type: String,
    default: '400px'
  },
  placeholder: {
    type: String,
    default: '请输入或上传集群凭证配置'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  maxFileSize: {
    type: Number,
    default: 5 * 1024 * 1024 // 5MB
  },
  acceptedFileNames: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const fileInputRef = ref(null)

// 双向绑定
const credentialContent = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 触发文件选择
const triggerFileUpload = () => {
  if (props.disabled) {
    ElMessage.warning('当前状态不允许上传文件')
    return
  }
  
  if (fileInputRef.value) {
    fileInputRef.value.click()
  } else {
    ElMessage.error('文件选择器初始化失败，请刷新页面重试')
  }
}

// 处理文件选择变化
const handleFileChange = (event) => {
  const file = event.target.files[0]
  if (!file) return

  console.log('选择的文件:', file.name, '文件大小:', file.size)
  
  // 仅在显式配置限制时校验文件名
  const fileName = file.name
  if (props.acceptedFileNames.length > 0) {
    const isValidFileName = props.acceptedFileNames.some(name =>
      fileName === name || fileName.toLowerCase().endsWith(`.${name}`)
    )

    if (!isValidFileName) {
      const acceptedNames = props.acceptedFileNames.join('", "')
      ElMessage.error(`请上传名为 "${acceptedNames}" 的文件`)
      event.target.value = ''
      return
    }
  }
  
  // 检查文件大小
  if (file.size > props.maxFileSize) {
    const sizeMB = (props.maxFileSize / (1024 * 1024)).toFixed(1)
    ElMessage.error(`文件大小不能超过 ${sizeMB}MB`)
    event.target.value = ''
    return
  }
  
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const content = e.target.result
      if (content && content.trim()) {
        credentialContent.value = content.trim()
        ElMessage.success(`凭证文件 ${fileName} 上传成功`)
      } else {
        ElMessage.warning('文件内容为空')
      }
    } catch (error) {
      console.error('读取文件失败:', error)
      ElMessage.error('读取文件失败，请检查文件格式')
    }
    event.target.value = ''
  }
  
  reader.onerror = (error) => {
    console.error('文件读取错误:', error)
    ElMessage.error('文件读取失败')
    event.target.value = ''
  }
  
  reader.readAsText(file, 'UTF-8')
}

// 清空凭证内容
const clearCredential = () => {
  if (props.disabled) {
    ElMessage.warning('当前状态不允许清空')
    return
  }
  
  if (!credentialContent.value.trim()) {
    ElMessage.info('凭证内容已经为空')
    return
  }
  
  ElMessageBox.confirm(
    '确定要清空凭证内容吗？',
    '清空确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  )
  .then(() => {
    credentialContent.value = ''
    ElMessage.success('凭证内容已清空')
  })
  .catch(() => {
    // 用户取消，不执行任何操作
  })
}
</script>

<style scoped>
.credential-editor-container {
  width: 100%;
}

.credential-editor {
  margin-bottom: 12px;
}

.editor {
  border-radius: 8px;
  overflow: hidden;
}

.credential-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-start;
  align-items: center;
}

.credential-actions .el-button {
  border-radius: 6px;
  font-size: 12px;
  height: 32px;
  padding: 0 14px;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.credential-actions .el-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.credential-actions .el-button .el-icon {
  font-size: 14px;
}

/* 确保上传按钮(绿色)正常显示 */
.credential-actions .upload-btn {
  background-color: #67c23a !important;
  border-color: #67c23a !important;
  color: white !important;
  cursor: pointer !important;
  opacity: 1 !important;
}

.credential-actions .upload-btn:hover:not(:disabled) {
  background-color: #5daf34 !important;
  border-color: #5daf34 !important;
  color: white !important;
}

.credential-actions .upload-btn:disabled {
  background-color: #a0cfff !important;
  border-color: #a0cfff !important;
  color: #ffffff !important;
  opacity: 0.6 !important;
  cursor: not-allowed !important;
}

/* 强制覆盖 Element Plus 的成功按钮样式 */
.credential-actions .el-button--success {
  background-color: #67c23a !important;
  border-color: #67c23a !important;
  color: white !important;
}

.credential-actions .el-button--success:hover:not(:disabled),
.credential-actions .el-button--success:focus:not(:disabled) {
  background-color: #5daf34 !important;
  border-color: #5daf34 !important;
  color: white !important;
}

.credential-actions .el-button--success:disabled {
  background-color: #a0cfff !important;
  border-color: #a0cfff !important;
  color: #ffffff !important;
  opacity: 0.6 !important;
  cursor: not-allowed !important;
}
</style>
