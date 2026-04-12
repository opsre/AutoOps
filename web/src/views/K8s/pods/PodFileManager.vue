<template>
  <div class="pod-file-manager">
    <div class="file-toolbar">
      <div class="breadcrumb-container" :class="{ 'is-root': currentPath === '/' }">
        <el-icon class="home-icon" @click="navigateTo(0)"><Location /></el-icon>
        <el-breadcrumb separator="/" style="flex: 1; margin-left: 8px;">
          <el-breadcrumb-item 
            v-for="(part, index) in pathParts" 
            :key="index"
            @click="navigateTo(index)"
            class="breadcrumb-item"
          >
            {{ part || '根目录' }}
          </el-breadcrumb-item>
        </el-breadcrumb>
      </div>

      <div class="toolbar-actions">
        <el-button type="primary" plain size="small" :icon="Refresh" @click="refreshFiles" :loading="loading">刷新</el-button>
        <el-button type="primary" plain size="small" :icon="FolderAdd" @click="showCreateDirDialog">新建目录</el-button>
        <el-upload
          class="upload-demo"
          action="#"
          :http-request="customUpload"
          :show-file-list="false"
        >
          <el-button type="success" size="small" :icon="Upload">上传文件</el-button>
        </el-upload>
        <el-button type="warning" plain size="small" :icon="Cpu" @click="handleHotReload" :loading="reloadLoading">热加载</el-button>
      </div>
    </div>

    <el-table
      v-loading="loading"
      :data="fileList"
      style="width: 100%"
      height="500px"
      border
      stripe
      highlight-current-row
      @row-dblclick="handleRowDblClick"
      class="custom-file-table"
    >
      <el-table-column prop="name" label="名称" min-width="240" show-overflow-tooltip>
        <template #default="scope">
          <div style="display: flex; align-items: center; cursor: pointer;" class="file-name-cell">
            <el-icon class="file-icon" 
              :color="scope.row.isDir ? '#e6a23c' : '#409EFF'">
              <Folder v-if="scope.row.isDir" />
              <Document v-else />
            </el-icon>
            <span class="file-name-text" :class="{ 'is-directory': scope.row.isDir }">{{ scope.row.name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="size" label="大小" width="120" align="right">
        <template #default="scope">
          <span class="file-size-text">{{ scope.row.isDir ? '-' : formatSize(scope.row.size) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="modTime" label="修改时间" width="180" align="center">
        <template #default="scope">
          <span class="file-time-text">{{ scope.row.modTime || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="mode" label="权限" width="120" align="center">
        <template #default="scope">
          <el-tag size="small" type="info" effect="plain">{{ scope.row.mode || '-' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right" align="center">
        <template #default="scope">
          <div class="action-buttons">
            <el-tooltip v-if="!scope.row.isDir" content="查看/编辑" placement="top" :show-after="200">
              <el-button type="primary" circle plain size="small" :icon="View" @click="viewAndEditFile(scope.row)" />
            </el-tooltip>
            <el-tooltip v-if="!scope.row.isDir" content="下载" placement="top" :show-after="200">
              <el-button type="success" circle plain size="small" :icon="Download" @click="downloadFile(scope.row)" />
            </el-tooltip>
            <el-popconfirm title="确定要删除这个文件/目录吗？" confirm-button-type="danger" @confirm="deleteFile(scope.row)" width="200px">
              <template #reference>
                <div class="delete-btn-wrapper">
                  <el-tooltip content="删除" placement="top" :show-after="200">
                    <el-button type="danger" circle plain size="small" :icon="Delete" />
                  </el-tooltip>
                </div>
              </template>
            </el-popconfirm>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <!-- 文件查看/编辑弹窗 -->
    <el-dialog
      v-model="editorVisible"
      :title="`编辑文件: ${currentEditFileName}`"
      width="800px"
      top="5vh"
      append-to-body
      destroy-on-close
      class="beautiful-editor-dialog"
    >
      <div v-loading="editorLoading" class="editor-container">
        <div class="editor-header">
          <el-icon class="header-icon"><Document /></el-icon>
          <span class="file-path">{{ currentEditFilePath }}</span>
        </div>
        <el-input
          v-model="fileContent"
          type="textarea"
          :rows="22"
          class="code-textarea beautified"
          placeholder="请输入文件内容..."
          spellcheck="false"
        />
      </div>
      <template #footer>
        <div class="dialog-footer-custom">
          <el-button @click="editorVisible = false" plain>取 消</el-button>
          <el-button type="primary" @click="saveFileContent" :loading="savingFile">保 存</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 新建目录弹窗 -->
    <el-dialog
      v-model="createDirVisible"
      title="新建目录"
      width="400px"
      append-to-body
      destroy-on-close
    >
      <el-input 
        v-model="newDirName" 
        placeholder="请输入目录名称" 
        @keyup.enter="handleCreateDir" 
      />
      <template #footer>
        <span class="dialog-footer-custom">
          <el-button @click="createDirVisible = false" plain>取 消</el-button>
          <el-button type="primary" @click="handleCreateDir" :loading="creatingDir">确 认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Folder, Document, Refresh, Upload, Download, Delete, Location, View, Cpu, FolderAdd } from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import storage from '@/utils/storage'

const props = defineProps({
  clusterId: { type: [String, Number], required: true },
  namespace: { type: String, required: true },
  podName: { type: String, required: true },
  containerName: { type: String, required: true }
})

const currentPath = ref('/')
const fileList = ref([])
const loading = ref(false)
const reloadLoading = ref(false)

// Editor State
const editorVisible = ref(false)
const editorLoading = ref(false)
const savingFile = ref(false)
const fileContent = ref('')
const currentEditFileName = ref('')
const currentEditFilePath = ref('')

// Create Directory State
const createDirVisible = ref(false)
const newDirName = ref('')
const creatingDir = ref(false)

const pathParts = computed(() => {
  const parts = currentPath.value.split('/').filter(p => p)
  return [''].concat(parts)
})

const fetchFiles = async () => {
  if (!props.containerName) return
  loading.value = true
  try {
    const res = await k8sApi.getPodFileList(
      props.clusterId,
      props.namespace,
      props.podName,
      props.containerName,
      currentPath.value
    )
    if (res.data && res.data.code === 200) {
      // 假设后端返回的数据结构中有 isDir 标识和文件属性，如果是上一级目录可以在前端增加或者后端返回 '..'
      console.log('获取文件列表成功:', res.data.data)
      const files = res.data.data || []
      // 排序：目录在前，文件在后
      fileList.value = files.sort((a, b) => {
        if (a.isDir && !b.isDir) return -1
        if (!a.isDir && b.isDir) return 1
        return a.name.localeCompare(b.name)
      })
    } else {
      ElMessage.error(res.data?.message || '获取文件列表失败')
    }
  } catch (error) {
    ElMessage.error('获取文件列表异常')
  } finally {
    loading.value = false
  }
}

const refreshFiles = () => {
  fetchFiles()
}

const showCreateDirDialog = () => {
  newDirName.value = ''
  createDirVisible.value = true
}

const handleCreateDir = async () => {
  if (!newDirName.value.trim()) {
    ElMessage.warning('请输入目录名称')
    return
  }
  
  creatingDir.value = true
  const dirPath = getFullFilePath(newDirName.value.trim())

  try {
    const res = await k8sApi.createPodDirectory(
      props.clusterId,
      props.namespace,
      props.podName,
      props.containerName,
      dirPath
    )
    if (res.data && res.data.code === 200) {
      ElMessage.success('目录创建成功')
      createDirVisible.value = false
      fetchFiles()
    } else {
      ElMessage.error(res.data?.message || '创建失败')
    }
  } catch (error) {
    ElMessage.error('创建异常')
  } finally {
    creatingDir.value = false
  }
}

const navigateTo = (index) => {
  if (index === 0) {
    currentPath.value = '/'
  } else {
    currentPath.value = '/' + pathParts.value.slice(1, index + 1).join('/')
  }
  fetchFiles()
}

const handleRowDblClick = (row) => {
  if (row.isDir) {
    if (currentPath.value.endsWith('/')) {
      currentPath.value += row.name
    } else {
      currentPath.value += '/' + row.name
    }
    fetchFiles()
  }
}

const getFullFilePath = (fileName) => {
  return currentPath.value.endsWith('/') 
    ? `${currentPath.value}${fileName}` 
    : `${currentPath.value}/${fileName}`
}

const viewAndEditFile = async (row) => {
  const filePath = getFullFilePath(row.name)
  currentEditFileName.value = row.name
  currentEditFilePath.value = filePath
  editorVisible.value = true
  editorLoading.value = true
  fileContent.value = ''

  try {
    const res = await k8sApi.getPodFileContent(
      props.clusterId,
      props.namespace,
      props.podName,
      props.containerName,
      filePath
    )
    if (res.data && res.data.code === 200) {
      fileContent.value = res.data.data || ''
    } else {
      ElMessage.error(res.data?.message || '读取文件内容失败')
      editorVisible.value = false
    }
  } catch (error) {
    ElMessage.error('读取文件内容异常')
    editorVisible.value = false
  } finally {
    editorLoading.value = false
  }
}

const saveFileContent = async () => {
  savingFile.value = true
  try {
    const res = await k8sApi.updatePodFileContent(
      props.clusterId,
      props.namespace,
      props.podName,
      props.containerName,
      currentEditFilePath.value,
      fileContent.value
    )
    if (res.data && res.data.code === 200) {
      ElMessage.success('保存成功')
      editorVisible.value = false
      refreshFiles()
    } else {
      ElMessage.error(res.data?.message || '保存失败')
    }
  } catch (error) {
    ElMessage.error('保存异常: ' + error.message)
  } finally {
    savingFile.value = false
  }
}

const handleHotReload = async () => {
  ElMessageBox.confirm(
    '确认要在该容器上执行热加载操作吗？',
    '提示',
    {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    reloadLoading.value = true
    try {
      const res = await k8sApi.hotReloadPod(
        props.clusterId,
        props.namespace,
        props.podName,
        props.containerName
      )
      if (res.data && res.data.code === 200) {
        ElMessage.success('热加载指令执行成功')
      } else {
        ElMessage.error(res.data?.message || '热加载操作失败')
      }
    } catch (err) {
      ElMessage.error('热加载操作异常: ' + err.message)
    } finally {
      reloadLoading.value = false
    }
  }).catch(() => {})
}

const downloadFile = async (row) => {
  const filePath = getFullFilePath(row.name)

  try {
    const res = await k8sApi.downloadPodFile(
      props.clusterId,
      props.namespace,
      props.podName,
      props.containerName,
      filePath
    )
    
    // Convert blob to object URL and trigger download
    const blob = new Blob([res.data])
    const downloadUrl = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = downloadUrl
    a.download = row.name
    a.style.display = 'none'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(downloadUrl)
    
  } catch (error) {
    ElMessage.error('下载失败')
  }
}

const deleteFile = async (row) => {
  const filePath = getFullFilePath(row.name)
    
  try {
    const res = await k8sApi.deletePodFile(
      props.clusterId,
      props.namespace,
      props.podName,
      props.containerName,
      filePath
    )
    if (res.data && res.data.code === 200) {
      ElMessage.success('删除成功')
      fetchFiles()
    } else {
      ElMessage.error(res.data?.message || '删除失败')
    }
  } catch (err) {
    ElMessage.error('删除异常')
  }
}

const customUpload = async (options) => {
  const formData = new FormData()
  formData.append('file', options.file)
  formData.append('containerName', props.containerName)
  formData.append('path', getFullFilePath(options.file.name))

  try {
    const res = await k8sApi.uploadPodFile(
      props.clusterId,
      props.namespace,
      props.podName,
      formData
    )
    if (res.data && res.data.code === 200) {
      ElMessage.success(`${options.file.name} 上传成功`)
      fetchFiles()
      options.onSuccess(res.data, options.file)
    } else {
      ElMessage.error(`${options.file.name} 上传失败: ${res.data?.message || '未知错误'}`)
      options.onError(new Error(res.data?.message || '未知错误'))
    }
  } catch (error) {
    ElMessage.error(`${options.file.name} 上传异常`)
    options.onError(error)
  }
}

const formatSize = (size) => {
  if (size === undefined || size === null) return '-'
  if (size === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(size) / Math.log(k))
  return (size / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i]
}

// 监听容器名称变更自动重新加载
watch(() => props.containerName, (newVal) => {
  if (newVal) {
    currentPath.value = '/'
    fetchFiles()
  }
})

onMounted(() => {
  if (props.containerName) {
    fetchFiles()
  }
})
</script>

<style scoped>
.pod-file-manager {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  border-radius: 6px;
  overflow: hidden;
}

.file-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #ebeef5;
}

.breadcrumb-container {
  display: flex;
  align-items: center;
  flex: 1;
  background-color: #fff;
  padding: 6px 12px;
  border-radius: 4px;
  border: 1px solid #dcdfe6;
  margin-right: 16px;
}

.home-icon {
  font-size: 16px;
  color: #909399;
  cursor: pointer;
  transition: color 0.2s;
}

.home-icon:hover {
  color: #409eff;
}

.breadcrumb-container.is-root .home-icon {
  color: #409eff;
}

.toolbar-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.upload-demo {
  display: inline-block;
  margin-left: 0;
}

.breadcrumb-item {
  cursor: pointer;
}
.breadcrumb-item:hover :deep(.el-breadcrumb__inner) {
  color: #409eff;
  text-decoration: underline;
}

/* Table styling */
.custom-file-table {
  --el-table-header-bg-color: #f5f7fa;
  --el-table-header-text-color: #606266;
}

.file-name-cell {
  padding: 4px 0;
}

.file-name-cell:hover .file-name-text {
  color: #409eff;
}

.file-icon {
  font-size: 20px;
  margin-right: 12px;
}

.file-name-text {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
  transition: color 0.2s;
}

.file-name-text.is-directory {
  color: #303133;
}

.file-size-text, .file-time-text {
  color: #909399;
  font-size: 13px;
}

.action-buttons {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
}

.delete-btn-wrapper {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

:deep(.el-table__row) {
  cursor: pointer;
}

:deep(.el-table--striped .el-table__body tr.el-table__row--striped td.el-table__cell) {
  background-color: #fafbfc;
}

.beautiful-editor-dialog :deep(.el-dialog__body) {
  padding: 15px 23px 5px;
}

.editor-container {
  border: 1px solid #3d424b;
  border-radius: 6px;
  overflow: hidden;
  box-shadow: 0 4px 14px 0 rgba(0, 0, 0, 0.1);
  background-color: #282c34;
}

.editor-header {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  background-color: #21252b;
  border-bottom: 1px solid #181a1f;
  color: #abb2bf;
  font-size: 13px;
}

.header-icon {
  margin-right: 6px;
  font-size: 16px;
  color: #e06c75;
}

.file-path {
  font-family: Consolas, Monaco, "Courier New", monospace;
  color: #61afef;
  font-weight: 500;
}

.code-textarea.beautified :deep(.el-textarea__inner) {
  border: none;
  border-radius: 0;
  box-shadow: none;
  padding: 16px;
  background-color: #282c34;
  color: #abb2bf;
  font-family: 'Fira Code', Consolas, Monaco, "Courier New", monospace;
  font-size: 14.5px;
  line-height: 1.6;
  white-space: pre;
  overflow: auto;
}

.code-textarea.beautified :deep(.el-textarea__inner:focus) {
  box-shadow: none;
  outline: none;
}

.dialog-footer-custom {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 10px;
}
</style>
