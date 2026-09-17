<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`查看集群凭证 - ${credential.clusterName}`"
    width="700px"
    class="credential-view-dialog"
    @closed="handleDialogClosed"
  >
    <div class="credential-content">
      <div class="credential-header">
        <el-button 
          size="small" 
          @click="copyCredential"
          :icon="DocumentCopy"
          :loading="copyLoading"
        >
          复制凭证
        </el-button>
        <el-button 
          size="small" 
          @click="downloadCredential"
          :icon="Download"
          :loading="downloadLoading"
        >
          下载文件
        </el-button>
      </div>
      
      <!-- 使用代码高亮显示凭证内容 -->
      <div class="code-container">
        <pre class="code-block yaml-code" v-html="highlightedCredential"></pre>
      </div>
      
      <div class="credential-info">
        <el-tag size="small" type="info">
          字符数：{{ credential.content?.length || 0 }}
        </el-tag>
        <el-tag size="small" type="success" style="margin-left: 8px;">
          格式：YAML/Config
        </el-tag>
        <el-tag size="small" type="warning" style="margin-left: 8px;" v-if="!credential.content?.trim()">
          内容为空
        </el-tag>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { DocumentCopy, Download } from '@element-plus/icons-vue'
import { highlight } from '@/utils/highlight'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  credential: {
    type: Object,
    default: () => ({
      clusterName: '',
      content: ''
    })
  }
})

const emit = defineEmits(['update:visible', 'closed'])

const copyLoading = ref(false)
const downloadLoading = ref(false)

// 双向绑定 visible
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 计算高亮后的凭证内容
const highlightedCredential = computed(() => {
  if (!props.credential.content) return '<span style="color: #909399; font-style: italic;">暂无凭证内容</span>'
  return highlight(props.credential.content, 'yaml')
})

// 复制凭证内容
const copyCredential = async () => {
  if (!props.credential.content?.trim()) {
    ElMessage.warning('凭证内容为空，无法复制')
    return
  }
  
  try {
    copyLoading.value = true
    await navigator.clipboard.writeText(props.credential.content)
    ElMessage.success('凭证内容已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    // 降级处理：使用传统方式复制
    try {
      const textArea = document.createElement('textarea')
      textArea.value = props.credential.content
      textArea.style.position = 'fixed'
      textArea.style.opacity = '0'
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      ElMessage.success('凭证内容已复制到剪贴板')
    } catch (fallbackError) {
      console.error('降级复制也失败:', fallbackError)
      ElMessage.error('复制失败，请手动复制')
    }
  } finally {
    copyLoading.value = false
  }
}

// 下载凭证文件
const downloadCredential = () => {
  if (!props.credential.content?.trim()) {
    ElMessage.warning('凭证内容为空，无法下载')
    return
  }
  
  if (!props.credential.clusterName?.trim()) {
    ElMessage.warning('集群名称为空，无法生成文件名')
    return
  }
  
  try {
    downloadLoading.value = true
    const blob = new Blob([props.credential.content], { type: 'text/plain;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${props.credential.clusterName}-config`
    link.style.display = 'none'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    ElMessage.success('凭证文件下载成功')
  } catch (error) {
    console.error('下载失败:', error)
    ElMessage.error('下载失败，请重试')
  } finally {
    downloadLoading.value = false
  }
}

// 对话框关闭处理
const handleDialogClosed = () => {
  emit('closed')
}

// 监听凭证变化，验证数据完整性
watch(() => props.credential, (newCredential) => {
  if (props.visible && newCredential) {
    if (!newCredential.clusterName?.trim()) {
      console.warn('凭证对话框：集群名称为空')
    }
    if (!newCredential.content?.trim()) {
      console.warn('凭证对话框：凭证内容为空')
    }
  }
}, { deep: true, immediate: true })
</script>

<style scoped>
/* 凭证查看对话框样式 */
.credential-view-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.credential-view-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.credential-view-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
  font-size: 16px;
}

.credential-view-dialog :deep(.el-dialog__headerbtn) {
  top: 20px;
  right: 20px;
}

.credential-view-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
  color: white;
  font-size: 18px;
}

.credential-view-dialog :deep(.el-dialog__body) {
  padding: 20px 24px;
}

.credential-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.credential-header {
  display: flex;
  gap: 12px;
  justify-content: flex-start;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #ebeef5;
}

.credential-header .el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.credential-header .el-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.code-container {
  position: relative;
  max-height: 400px;
  overflow: auto;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #fafafa;
}

.code-block {
  margin: 0;
  padding: 20px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.6;
  border-radius: 8px;
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
  min-height: 100px;
}

/* highlight.js 主题样式会自动应用，这里只需要基础样式 */
.yaml-code {
  background: #0d1117;  /* GitHub Dark 主题背景色 */
  color: #c9d1d9;       /* GitHub Dark 主题文字色 */
}

.credential-info {
  display: flex;
  gap: 8px;
  justify-content: flex-start;
  align-items: center;
  padding: 12px 0;
  border-top: 1px solid #ebeef5;
  flex-wrap: wrap;
}

/* 滚动条样式 */
.code-container::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.code-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.code-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.code-container::-webkit-scrollbar-thumb:hover {
  background: #a1a1a1;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .credential-view-dialog :deep(.el-dialog) {
    width: 95% !important;
    margin: 5vh auto !important;
  }
  
  .credential-header {
    flex-direction: column;
    gap: 8px;
    align-items: stretch;
  }
  
  .credential-header .el-button {
    width: 100%;
  }
  
  .code-container {
    max-height: 300px;
  }
  
  .code-block {
    padding: 15px;
    font-size: 12px;
  }
}
</style>