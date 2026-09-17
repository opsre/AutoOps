<template>
  <div>
    <!-- 文件上传对话框 -->
    <el-dialog title="文件上传" v-model="uploadDialogVisible" width="480px" :close-on-click-modal="!isUploading">
      <el-form :model="uploadForm" :rules="uploadRules" ref="uploadFormRef" label-width="100px">
        <el-form-item label="目标路径" prop="targetPath">
          <el-input v-model="uploadForm.targetPath" placeholder="请输入目标路径" />
        </el-form-item>
        <el-form-item label="上传文件" prop="file">
          <el-upload
            class="upload-demo"
            :auto-upload="false"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            :show-file-list="false"
          >
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip" v-if="uploadForm.file">
                已选择: {{ uploadForm.file.name }}
                <el-button
                  type="danger"
                  text
                  icon="Close"
                  circle
                  size="small"
                  @click.stop="handleFileRemove"
                  style="margin-left: 8px"
                />
              </div>
              <div class="el-upload__tip" style="color: #999; margin-top: 15px">
                提示：请上传小于 500MB 的文件
              </div>
            </template>
          </el-upload>
        </el-form-item>
        <el-progress
          v-if="isUploading"
          :percentage="uploadProgress"
          :status="uploadProgress === 100 ? 'success' : ''"
        />
      </el-form>
      <template #footer>
        <el-button @click="uploadDialogVisible = false" :disabled="isUploading">取消</el-button>
        <el-button
          type="primary"
          @click="handleUpload"
          :loading="isUploading"
          :disabled="!uploadForm.file"
        >
          开始上传
        </el-button>
      </template>
    </el-dialog>

    <!-- 文件管理主对话框 -->
    <el-dialog
      title="文件管理"
      v-model="dialogVisible"
      width="60%"
      :close-on-click-modal="false"
      :before-close="handleDialogClose"
    >
      <!-- Pod 信息概览 -->
      <div class="manager-overview">
        <span class="overview-item overview-pod">
          <span class="overview-label">Pod</span>
          <span class="overview-value">{{ podName }}</span>
        </span>
        <span class="overview-item overview-namespace">
          <span class="overview-label">命名空间</span>
          <span class="overview-value">{{ namespaceName }}</span>
        </span>
        <span class="overview-item overview-container">
          <span class="overview-label">容器</span>
          <span class="overview-value">{{ containerName || '默认容器' }}</span>
        </span>
      </div>

      <!-- 面包屑导航 -->
      <el-breadcrumb separator="/" class="breadcrumb-nav">
        <el-breadcrumb-item @click="navigateToPath('/')" style="cursor: pointer">
          <el-icon><HomeFilled /></el-icon>
          根目录
        </el-breadcrumb-item>
        <el-breadcrumb-item
          v-for="(part, index) in pathParts"
          :key="index"
          @click="navigateToPath(getPathUpTo(index))"
          style="cursor: pointer"
        >
          {{ part }}
        </el-breadcrumb-item>
      </el-breadcrumb>

      <!-- 当前路径 -->
      <div class="current-path">
        <el-input
          v-model="currentPath"
          placeholder="输入路径后按回车刷新"
          clearable
          size="small"
          @keyup.enter="handlePathInputEnter"
          @clear="navigateToPath('/')"
        >
          <template #prepend>
            <span style="font-weight: 500;">当前路径</span>
          </template>
        </el-input>
        <el-button
          size="small"
          @click="handlePathInputEnter"
          title="刷新当前路径"
          class="path-action-btn btn-refresh"
          type="success"
          plain
        >
          <el-icon><Refresh /></el-icon>
        </el-button>
        <el-button
          size="small"
          @click="navigateToParentPath"
          title="返回上一级"
          class="path-action-btn btn-back"
          type="warning"
          :disabled="currentPath === '/'"
        >
          <el-icon><Back /></el-icon>
        </el-button>
        <el-button
          size="small"
          @click="handleUploadInCurrentPath"
          title="在当前目录上传文件"
          class="path-action-btn btn-upload"
          type="primary"
          plain
        >
          <el-icon><Upload /></el-icon>
        </el-button>
      </div>

      <!-- 文件列表 -->
      <el-table
        v-loading="isLoadingFiles"
        :data="fileList"
        style="width: 100%"
        class="file-table"
        @row-click="handleFileRowClick"
        highlight-current-row
        max-height="400"
      >
        <el-table-column label="名称" min-width="250">
          <template #default="{ row }">
            <el-icon v-if="row.isDir" style="margin-right: 8px; color: #409eff; vertical-align: middle;"><Folder /></el-icon>
            <el-icon v-else style="margin-right: 8px; color: #67c23a; vertical-align: middle;"><Document /></el-icon>
            <span :style="{ cursor: row.isDir ? 'pointer' : 'default', color: row.isDir ? '#409eff' : '' }">
              {{ row.name }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="120">
          <template #default="{ row }">
            {{ row.isDir ? '-' : formatFileSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column prop="mode" label="权限" width="120" />
        <el-table-column prop="modTime" label="修改时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="!row.isDir"
              type="primary"
              size="small"
              link
              :loading="downloadingFile === row.path"
              @click.stop="handleDownloadSelectedFile(row)"
            >
              下载
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 下载进度条 -->
      <div v-if="downloadProgress.show" class="download-progress">
        <div class="progress-header">
          <span>正在下载: {{ downloadProgress.fileName }}</span>
          <span>{{ formatFileSize(downloadProgress.received) }} / {{ formatFileSize(downloadProgress.total) }}</span>
        </div>
        <el-progress
          :percentage="downloadProgress.percentage"
          :status="downloadProgress.status"
          :stroke-width="20"
        >
          <span style="font-size: 12px;">{{ downloadProgress.percentage }}%</span>
        </el-progress>
      </div>

      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { HomeFilled, Refresh, Back, Upload, Folder, Document } from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import storage from '@/utils/storage'

const props = defineProps({
  visible: { type: Boolean, default: false },
  clusterId: { type: [String, Number], required: true },
  namespaceName: { type: String, required: true },
  podName: { type: String, required: true },
  containerName: { type: String, default: '' }
})

const emit = defineEmits(['update:visible', 'upload-success', 'download-success'])

const dialogVisible = computed({
  get: () => props.visible,
  set: val => emit('update:visible', val)
})

// 文件浏览器状态
const currentPath = ref('/')
const fileList = ref([])
const isLoadingFiles = ref(false)
const downloadingFile = ref(null)
const downloadProgress = ref({
  show: false, fileName: '', received: 0, total: 0, percentage: 0, status: ''
})

// 上传状态
const uploadDialogVisible = ref(false)
const uploadFormRef = ref(null)
const isUploading = ref(false)
const uploadProgress = ref(0)
const uploadForm = ref({ targetPath: '/tmp', file: null })
const uploadRules = {
  targetPath: [{ required: true, message: '请输入目标路径', trigger: 'blur' }],
  file: [{ required: true, message: '请选择上传文件', trigger: 'change' }]
}

const pathParts = computed(() => {
  if (currentPath.value === '/') return []
  return currentPath.value.split('/').filter(Boolean)
})

watch(() => props.visible, val => {
  if (val) {
    currentPath.value = '/'
    fileList.value = []
    downloadProgress.value.show = false
    loadRemoteFiles('/')
  }
})

const handleDialogClose = done => {
  if (isUploading.value) {
    ElMessage.warning('文件上传进行中，请稍后再关闭')
    return
  }
  done()
}

// 面包屑路径计算
const getPathUpTo = index => {
  const parts = pathParts.value.slice(0, index + 1)
  return '/' + parts.join('/')
}

const navigateToPath = path => loadRemoteFiles(path)

const navigateToParentPath = () => {
  if (currentPath.value === '/') return
  const parts = currentPath.value.split('/').filter(Boolean)
  parts.pop()
  navigateToPath(parts.length === 0 ? '/' : '/' + parts.join('/'))
}

const handlePathInputEnter = () => {
  const path = (currentPath.value || '/').trim()
  loadRemoteFiles(path)
}

const handleFileRowClick = row => {
  if (row.isDir) loadRemoteFiles(row.path)
}

// 加载文件列表
const loadRemoteFiles = async path => {
  // 规范化路径，去除重复的 /
  const normalizedPath = ('/' + path).replace(/\/+/g, '/').replace(/(.+)\/$/, '$1') || '/'
  isLoadingFiles.value = true
  try {
    const params = { remotePath: normalizedPath }
    if (props.containerName) params.containerName = props.containerName

    const res = await k8sApi.listPodFiles(props.clusterId, props.namespaceName, props.podName, params)
    const data = res.data || res
    if (data.code === 200) {
      fileList.value = data.data?.files || data.data || []
      currentPath.value = normalizedPath
    } else {
      ElMessage.error(data.message || '加载文件列表失败')
    }
  } catch (e) {
    ElMessage.error('加载文件列表失败: ' + (e.message || '未知错误'))
  } finally {
    isLoadingFiles.value = false
  }
}

// 上传相关
const handleUploadInCurrentPath = () => {
  uploadForm.value = { targetPath: currentPath.value, file: null }
  uploadProgress.value = 0
  uploadDialogVisible.value = true
}

const handleFileChange = file => {
  uploadForm.value.file = file.raw || file
  uploadFormRef.value?.validateField('file')
}

const handleFileRemove = () => {
  uploadForm.value.file = null
  uploadFormRef.value?.validateField('file')
}

const handleUpload = async () => {
  try {
    await uploadFormRef.value.validate()
    if (!uploadForm.value.file) return ElMessage.warning('请选择上传文件')

    isUploading.value = true
    uploadProgress.value = 0

    const formData = new FormData()
    formData.append('file', uploadForm.value.file)
    formData.append('destPath', uploadForm.value.targetPath || '/tmp')
    if (props.containerName) formData.append('containerName', props.containerName)

    const res = await k8sApi.uploadPodFile(
      props.clusterId, props.namespaceName, props.podName, formData,
      {
        timeout: 600000,
        onUploadProgress: e => {
          const total = Number(e.total || uploadForm.value.file?.size || 0)
          if (total) uploadProgress.value = Math.min(99, Math.round((e.loaded / total) * 100))
        }
      }
    )

    const data = res.data || res
    if (data.code && data.code !== 200) throw new Error(data.message || '上传失败')

    uploadProgress.value = 100
    ElMessage.success('文件上传成功')
    uploadDialogVisible.value = false
    emit('upload-success')
    loadRemoteFiles(currentPath.value)
  } catch (e) {
    ElMessage.error('文件上传失败: ' + (e.message || '未知错误'))
  } finally {
    isUploading.value = false
  }
}

// 下载文件
const handleDownloadSelectedFile = async file => {
  downloadingFile.value = file.path
  downloadProgress.value = { show: true, fileName: file.name, received: 0, total: 0, percentage: 0, status: '' }

  try {
    const token = storage.getItem('token')
    const tokenValue = typeof token === 'object' ? token.access_token || token.token : token

    const params = new URLSearchParams({ remotePath: file.path })
    if (props.containerName) params.set('containerName', props.containerName)
    const url = `/api/v1/k8s/cluster/${props.clusterId}/namespaces/${props.namespaceName}/pods/${props.podName}/files/download?${params.toString()}`

    await new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest()
      xhr.responseType = 'arraybuffer'
      xhr.timeout = 600000

      xhr.addEventListener('progress', e => {
        downloadProgress.value.received = e.loaded
        if (e.lengthComputable) {
          downloadProgress.value.total = e.total
          downloadProgress.value.percentage = Math.round((e.loaded / e.total) * 100)
        }
      })

      xhr.addEventListener('load', () => {
        if (xhr.status === 200) {
          const blob = new Blob([xhr.response])
          const link = document.createElement('a')
          link.href = window.URL.createObjectURL(blob)
          link.download = file.name
          document.body.appendChild(link)
          link.click()
          setTimeout(() => {
            document.body.removeChild(link)
            window.URL.revokeObjectURL(link.href)
          }, 100)
          downloadProgress.value.percentage = 100
          downloadProgress.value.status = 'success'
          ElMessage.success(`${file.name} 下载成功`)
          emit('download-success')
          setTimeout(() => { downloadProgress.value.show = false }, 2000)
          resolve()
        } else {
          reject(new Error(`HTTP ${xhr.status}`))
        }
      })

      xhr.addEventListener('error', () => reject(new Error('网络请求失败')))
      xhr.addEventListener('timeout', () => reject(new Error('请求超时')))

      xhr.open('GET', url, true)
      if (tokenValue) xhr.setRequestHeader('Authorization', `Bearer ${tokenValue}`)
      xhr.send()
    })
  } catch (e) {
    downloadProgress.value.status = 'exception'
    ElMessage.error('文件下载失败: ' + (e.message || '未知错误'))
    setTimeout(() => { downloadProgress.value.show = false }, 3000)
  } finally {
    downloadingFile.value = null
  }
}

const formatFileSize = bytes => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
}
</script>

<style scoped>
/* Pod 信息概览 */
.manager-overview {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid #eef2f7;
}

.overview-item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.overview-label {
  font-size: 12px;
  font-weight: 600;
}

.overview-value {
  font-size: 14px;
  font-weight: 500;
  word-break: break-all;
}

.overview-pod { color: #2563eb; }
.overview-namespace { color: #ca8a04; }
.overview-container { color: #059669; }

/* 面包屑 */
.breadcrumb-nav {
  margin-bottom: 16px;
  padding: 12px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.breadcrumb-nav :deep(.el-breadcrumb__item) {
  cursor: pointer;
  transition: all 0.3s;
}

.breadcrumb-nav :deep(.el-breadcrumb__item:hover) {
  color: #409eff;
}

/* 路径栏 */
.current-path {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.current-path .el-input {
  flex: 1;
}

.current-path .path-action-btn {
  flex-shrink: 0;
  padding: 8px 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.current-path :deep(.el-input-group__prepend) {
  background-color: #ecf5ff;
  border-color: #409eff;
  color: #409eff;
  padding: 0 15px;
}

.current-path :deep(.el-input__wrapper) {
  border-left-color: #409eff;
}

/* 文件表格 */
.file-table {
  border-radius: 8px;
  overflow: hidden;
}

.file-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.file-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.file-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.file-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* 下载进度 */
.download-progress {
  margin-top: 16px;
  padding: 14px 16px;
  border-radius: 12px;
  background: #f8fafc;
  border: 1px solid #e5e7eb;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  color: #606266;
}
</style>
