<template>
  <el-dialog
    v-model="isVisible"
    title="上传证书"
    width="600px"
    @close="handleClose"
  >
    <div class="upload-container">
      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :file-list="fileList"
        action="#"
        :before-upload="handleBeforeUpload"
        @change="handleChange"
        @remove="handleRemove"
        accept=".zip"
        drag
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          拖拽证书文件到这里 或 <em>点击选择</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            只能上传 .zip 格式的证书文件，大小不超过 10MB
          </div>
        </template>
      </el-upload>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          确定上传
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { UploadFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  fileList: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:visible', 'file-change', 'file-remove', 'submit'])

const isVisible = ref(false)
const uploadRef = ref(null)

// 监听 visible prop
watch(() => props.visible, (newVal) => {
  isVisible.value = newVal
})

// 监听 isVisible 变化，更新 visible prop
watch(isVisible, (newVal) => {
  if (!newVal) {
    emit('update:visible', false)
  }
})

// 上传前验证
const handleBeforeUpload = (file) => {
  const isZip = file.type === 'application/zip' ||
                file.type === 'application/x-zip-compressed' ||
                file.name.endsWith('.zip')

  if (!isZip) {
    ElMessage.error('只能上传 .zip 格式的证书文件')
    return false
  }

  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('证书文件大小不能超过 10MB')
    return false
  }

  return true
}

// 文件列表变化
const handleChange = (file, fileList) => {
  emit('file-change', file, fileList)
}

// 移除文件
const handleRemove = (file, fileList) => {
  emit('file-remove', file, fileList)
}

// 关闭对话框
const handleClose = () => {
  isVisible.value = false
}

// 提交上传
const handleSubmit = () => {
  if (props.fileList.length === 0) {
    ElMessage.warning('请选择要上传的证书文件')
    return
  }
  emit('submit')
}
</script>

<style scoped>
.upload-container {
  padding: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

:deep(.el-upload-dragger) {
  width: 100%;
}
</style>
