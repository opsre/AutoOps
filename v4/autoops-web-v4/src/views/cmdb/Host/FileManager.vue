<template>
  <div>
    <!-- 文件上传对话框 -->
    <el-dialog title="文件上传" v-model="uploadDialogVisible" width="25%" @close="handleUploadDialogClose">
      <el-form :model="uploadForm" :rules="uploadRules" ref="uploadFormRef" label-width="100px">
        <el-form-item label="目标主机">
          <el-input v-model="currentUploadHost.hostName" disabled />
        </el-form-item>
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
                提示：请上传小于100MB的文件
              </div>
            </template>
          </el-upload>
        </el-form-item>
        <div v-if="isUploading" class="upload-progress-section">
          <div class="upload-progress-item">
            <div class="upload-progress-header">
              <span>上传到目标主机中...</span>
            </div>
            <el-progress
              :percentage="remoteUploadProgress"
              :status="uploadTaskStatus === 'completed' ? 'success' : uploadTaskStatus === 'failed' ? 'exception' : ''"
            >
              <span
                class="upload-progress-value"
                :class="{
                  'is-success': uploadTaskStatus === 'completed',
                  'is-failed': uploadTaskStatus === 'failed'
                }"
              >
                {{ remoteUploadProgress }}%
              </span>
            </el-progress>
          </div>
          <div class="upload-progress-tip">
            {{ uploadStatusText }}
          </div>
        </div>
      </el-form>
      <template #footer>
        <el-button @click="uploadDialogVisible = false" :disabled="isUploading">取消</el-button>
        <el-button
          type="primary"
          @click="handleUpload"
          :loading="isUploading"
          :disabled="!uploadForm.file || !uploadForm.hostId"
        >
          开始上传
        </el-button>
      </template>
    </el-dialog>

    <!-- 文件下载对话框 -->
    <el-dialog title="文件管理" v-model="downloadDialogVisible" width="50%">
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
          <img src="@/assets/image/刷新.svg" class="btn-icon" />
        </el-button>
        <el-button
          size="small"
          @click="navigateToParentPath"
          title="返回上一级"
          class="path-action-btn btn-back"
          type="warning"
          :disabled="currentPath === '/'"
        >
          <img src="@/assets/image/返回.svg" class="btn-icon" />
        </el-button>
        <el-button
          size="small"
          @click="handleUploadInCurrentPath"
          title="在当前目录上传文件"
          class="path-action-btn btn-upload"
          type="primary"
          plain
        >
          <img src="@/assets/image/上传.svg" class="btn-icon" />
        </el-button>
      </div>

      <!-- 文件列表 -->
      <el-table
        v-loading="isLoadingFiles"
        :data="fileList"
        style="width: 100%"
        @row-click="handleFileRowClick"
        highlight-current-row
        max-height="400"
      >
        <el-table-column label="名称" min-width="250">
          <template #default="{ row }">
            <img v-if="row.isDir" src="@/assets/image/关闭文件夹.svg" style="width: 16px; height: 16px; margin-right: 8px; vertical-align: middle;" />
            <el-icon v-else style="margin-right: 8px; color: #67c23a">
              <Document />
            </el-icon>
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
      <div v-if="downloadProgress.show" class="download-progress" style="margin-top: 16px;">
        <div style="margin-bottom: 8px; display: flex; justify-content: space-between; align-items: center;">
          <span style="font-size: 14px; color: #606266;">
            正在下载: {{ downloadProgress.fileName }}
          </span>
          <span style="font-size: 12px; color: #909399;">
            {{ formatFileSize(downloadProgress.received) }} / {{ formatFileSize(downloadProgress.total) }}
          </span>
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
        <el-button @click="downloadDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'FileManager',
  data() {
    return {
      // 上传相关
      uploadDialogVisible: false,
      uploadForm: {
        hostId: null,
        file: null,
        targetPath: '/tmp'
      },
      currentUploadHost: {
        hostName: ''
      },
      uploadRules: {
        file: [{ required: true, message: '请选择上传文件', trigger: 'change' }],
        targetPath: [{
          required: true,
          message: '请输入目标路径',
          trigger: ['blur', 'change'],
          validator: (rule, value, callback) => {
            if (value === '/tmp' || (value && value.trim() !== '')) {
              callback()
            } else {
              callback(new Error('请输入目标路径'))
            }
          }
        }]
      },
      isUploading: false,
      remoteUploadProgress: 0,
      uploadTaskId: '',
      uploadTaskStatus: '',
      uploadStatusText: '',
      uploadPollingTimer: null,

      // 下载相关
      downloadDialogVisible: false,
      downloadForm: {
        hostId: null,
        remotePath: ''
      },
      currentDownloadHost: {
        hostName: ''
      },
      downloadRules: {
        remotePath: [
          { required: true, message: '请输入远程文件路径', trigger: 'blur' }
        ]
      },
      isDownloading: false,
      // 文件浏览器相关
      currentPath: '/',
      fileList: [],
      isLoadingFiles: false,
      downloadingFile: null,
      // 下载进度
      downloadProgress: {
        show: false,
        fileName: '',
        received: 0,
        total: 0,
        percentage: 0,
        status: ''
      }
    }
  },
  computed: {
    // 计算面包屑路径
    pathParts() {
      if (this.currentPath === '/') return []
      return this.currentPath.split('/').filter(Boolean)
    }
  },
  methods: {
    // ==================== 上传相关方法 ====================

    // 显示上传对话框
    showUploadDialog(host) {
      this.currentUploadHost = host
      this.uploadForm = {
        hostId: host.id,
        file: null,
        targetPath: '/tmp'
      }
      this.$nextTick(() => {
        this.$refs.uploadFormRef?.clearValidate('targetPath')
      })
      this.uploadDialogVisible = true
    },

    // 文件选择处理
    handleFileChange(file) {
      this.uploadForm.file = file.raw
      this.$refs.uploadFormRef.validateField('file')
    },

    // 文件删除处理
    handleFileRemove() {
      this.uploadForm.file = null
      this.$refs.uploadFormRef.validateField('file')
    },

    extractUploadTaskInfo(responseData) {
      if (!responseData || typeof responseData !== 'object') {
        return { taskId: '', task: null }
      }

      const payload = responseData.data && typeof responseData.data === 'object'
        ? responseData.data
        : responseData
      const nestedTask = payload.task && typeof payload.task === 'object'
        ? payload.task
        : null

      return {
        taskId: nestedTask?.taskId || payload.taskId || payload.jobId || responseData.taskId || responseData.jobId || '',
        task: nestedTask || payload
      }
    },

    normalizeUploadTask(task) {
      const source = task && typeof task === 'object' && task.task && typeof task.task === 'object'
        ? task.task
        : task

      if (!source || typeof source !== 'object') {
        return null
      }

      const rawProgress = Number(source.progress ?? source.percentage ?? 0)
      return {
        status: source.status || source.state || '',
        progress: Math.min(100, Math.max(0, Math.round(rawProgress))),
        message: source.message || source.msg || '',
        taskId: source.taskId || ''
      }
    },

    clearUploadTaskPolling() {
      if (this.uploadPollingTimer) {
        clearInterval(this.uploadPollingTimer)
        this.uploadPollingTimer = null
      }
    },

    async pollUploadTaskStatus() {
      if (!this.uploadTaskId) {
        return
      }

      try {
        const { data: res } = await this.$api.getUploadFileTaskStatus(this.uploadTaskId)
        const task = this.normalizeUploadTask(res?.data || res)

        if (!task) {
          return
        }

        this.uploadTaskStatus = task.status
        this.remoteUploadProgress = task.progress
        this.uploadStatusText = task.message || (
          task.status === 'completed'
            ? '上传完成'
            : task.status === 'failed'
              ? '上传失败'
              : '传输到远程主机中...'
        )

        if (task.status === 'completed') {
          this.clearUploadTaskPolling()
          this.remoteUploadProgress = 100
          this.$message.success('文件上传成功')
          this.uploadDialogVisible = false
          this.resetUploadForm()
          this.$emit('upload-success')

          if (this.downloadDialogVisible) {
            this.loadRemoteFiles(this.currentPath)
          }
        } else if (task.status === 'failed') {
          this.clearUploadTaskPolling()
          this.isUploading = false
          this.$message.error(task.message || '文件上传失败')
        }
      } catch (error) {
        console.error('查询上传任务进度失败:', error)
        this.clearUploadTaskPolling()
        this.isUploading = false
        this.$message.error('查询上传进度失败: ' + (error.message || '未知错误'))
      }
    },

    startUploadTaskPolling(taskId) {
      this.uploadTaskId = taskId
      this.uploadTaskStatus = this.uploadTaskStatus || 'uploading'
      this.remoteUploadProgress = this.remoteUploadProgress || 0
      this.uploadStatusText = this.uploadStatusText || '传输到远程主机中...'
      this.clearUploadTaskPolling()
      this.pollUploadTaskStatus()
      this.uploadPollingTimer = setInterval(() => {
        this.pollUploadTaskStatus()
      }, 2000)
    },

    // 文件上传处理
    async handleUpload() {
      try {
        // 验证表单
        await this.$refs.uploadFormRef.validate()

        if (!this.uploadForm.file) {
          return this.$message.warning('请选择上传文件')
        }

        // 检查文件大小 (100MB限制)
        if (this.uploadForm.file.size > 100 * 1024 * 1024) {
          return this.$message.warning('文件大小不能超过100MB')
        }

        // 检查是否已有上传在进行
        if (this.isUploading) {
          return this.$message.warning('已有文件正在上传，请等待完成')
        }

        this.isUploading = true
        this.remoteUploadProgress = 0
        this.uploadTaskId = ''
        this.uploadTaskStatus = ''
        this.uploadStatusText = '上传到平台中...'
        this.clearUploadTaskPolling()

        // 确保目标路径有值，使用默认路径'/tmp'如果为空
        const destPath = this.uploadForm.targetPath || '/tmp'

        const formData = new FormData()
        formData.append('file', this.uploadForm.file)
        formData.append('destPath', destPath)

        const config = {
          headers: {
            'Content-Type': 'multipart/form-data'
          },
          timeout: 60000, // 60秒超时
          onUploadProgress: progressEvent => {
            if (!progressEvent.total) {
              return
            }
            if (progressEvent.loaded < progressEvent.total) {
              this.uploadStatusText = '上传到平台中...'
            }
          }
        }

        // 主机ID作为路径参数，其他参数通过formData传递
        const { data: res } = await this.$api.uploadFileToHost(
          this.uploadForm.hostId,
          formData,
          config
        )

        if (res.code === 200) {
          const { taskId, task } = this.extractUploadTaskInfo(res)

          if (taskId) {
            const normalizedTask = this.normalizeUploadTask(task)
            if (normalizedTask) {
              this.uploadTaskStatus = normalizedTask.status || 'uploading'
              this.remoteUploadProgress = normalizedTask.progress
              this.uploadStatusText = normalizedTask.message || '传输到远程主机中...'
            }
            this.startUploadTaskPolling(taskId)
          } else {
            this.uploadTaskStatus = 'failed'
            this.uploadStatusText = '平台已接收文件，但未返回上传任务ID'
            this.$message.error('平台已接收文件，但未返回上传任务ID，无法继续跟踪远端上传进度')
          }
        } else {
          this.uploadTaskStatus = 'failed'
          this.uploadStatusText = res.message || '文件上传失败'
          this.$message.error(res.message || '文件上传失败')
        }
      } catch (error) {
        console.error('文件上传失败:', error)
        this.$message.error('文件上传失败: ' + error.message)
      } finally {
        if (!this.uploadTaskId) {
          this.isUploading = false
        }
      }
    },

    // 重置上传表单
    resetUploadForm() {
      this.clearUploadTaskPolling()
      this.uploadForm = {
        hostId: this.currentUploadHost?.id || null,
        file: null,
        targetPath: '/tmp'
      }
      this.remoteUploadProgress = 0
      this.uploadTaskId = ''
      this.uploadTaskStatus = ''
      this.uploadStatusText = ''
      this.isUploading = false
      this.$nextTick(() => {
        this.$refs.uploadFormRef?.clearValidate('targetPath')
      })
    },

    handleUploadDialogClose() {
      if (!this.isUploading) {
        this.resetUploadForm()
      }
    },

    // ==================== 下载相关方法 ====================

    // 显示下载对话框
    showDownloadDialog(host) {
      this.currentDownloadHost = host
      this.downloadForm = {
        hostId: host.id,
        remotePath: ''
      }
      this.currentPath = '/'
      this.fileList = []
      this.downloadDialogVisible = true
      // 加载根目录文件
      this.loadRemoteFiles('/')
    },

    // 处理路径输入回车
    handlePathInputEnter() {
      if (!this.currentPath) {
        this.$message.warning('请输入路径')
        return
      }
      // 确保路径格式正确
      let path = this.currentPath.trim()
      if (!path.startsWith('/')) {
        path = '/' + path
      }
      this.loadRemoteFiles(path)
    },

    // 返回上一级目录
    navigateToParentPath() {
      if (this.currentPath === '/') {
        this.$message.info('已经在根目录')
        return
      }
      const pathParts = this.currentPath.split('/').filter(Boolean)
      if (pathParts.length === 0) {
        this.navigateToPath('/')
      } else {
        pathParts.pop() // 移除最后一级
        const parentPath = pathParts.length === 0 ? '/' : '/' + pathParts.join('/')
        this.navigateToPath(parentPath)
      }
    },

    // 在当前目录上传文件
    handleUploadInCurrentPath() {
      // 设置上传表单的目标路径为当前路径
      this.uploadForm.targetPath = this.currentPath
      this.uploadForm.hostId = this.currentDownloadHost.id
      this.currentUploadHost = this.currentDownloadHost
      this.uploadDialogVisible = true
    },

    // 加载远程文件列表
    async loadRemoteFiles(path) {
      this.isLoadingFiles = true
      try {
        const { data: res } = await this.$api.listRemoteFiles(
          this.currentDownloadHost.id,
          path
        )

        if (res.code === 200) {
          this.fileList = res.data.files || []
          this.currentPath = path
        } else {
          this.$message.error(res.message || '加载文件列表失败')
        }
      } catch (error) {
        console.error('加载文件列表失败:', error)
        this.$message.error('加载文件列表失败: ' + (error.message || '未知错误'))
      } finally {
        this.isLoadingFiles = false
      }
    },

    // 导航到指定路径
    navigateToPath(path) {
      this.loadRemoteFiles(path)
    },

    // 获取指定索引之前的路径
    getPathUpTo(index) {
      if (index < 0) return '/'
      const parts = this.pathParts.slice(0, index + 1)
      return '/' + parts.join('/')
    },

    // 处理文件行点击
    handleFileRowClick(row) {
      if (row.isDir) {
        this.loadRemoteFiles(row.path)
      }
    },

    // 下载选中的文件
    async handleDownloadSelectedFile(file) {
      try {
        this.downloadingFile = file.path

        // 显示进度条
        this.downloadProgress = {
          show: true,
          fileName: file.name,
          received: 0,
          total: 0,
          percentage: 0,
          status: ''
        }

        // 从 storage 获取 token（与 request.js 保持一致）
        const storage = await import('@/utils/storage')
        const token = storage.default.getItem("token") || {}
        // 如果token是对象，获取其access_token属性；如果是字符串直接使用
        const tokenValue = typeof token === 'object' ? token.access_token || token.token : token

        // 使用 XMLHttpRequest 下载文件以支持进度监听
        const url = `/api/v1/cmdb/hostssh/download/${this.currentDownloadHost.id}?remotePath=${encodeURIComponent(file.path)}`

        console.log('开始下载文件:', url)
        console.log('Token:', tokenValue ? '存在' : '不存在')

        // 使用 Promise 包装 XMLHttpRequest
        await new Promise((resolve, reject) => {
          const xhr = new XMLHttpRequest()

          console.log('1. 创建 XMLHttpRequest 对象')

          // 监听下载进度事件（必须在 open 之前设置）
          xhr.addEventListener('progress', (event) => {
            console.log('✓ progress 事件触发!', {
              lengthComputable: event.lengthComputable,
              loaded: event.loaded,
              total: event.total
            })

            if (event.lengthComputable) {
              const received = event.loaded
              const total = event.total
              const percentage = Math.round((received / total) * 100)

              this.downloadProgress.total = total
              this.downloadProgress.received = received
              this.downloadProgress.percentage = percentage

              console.log(`📊 下载进度: ${percentage}% (${(received / 1024 / 1024).toFixed(2)} MB / ${(total / 1024 / 1024).toFixed(2)} MB)`)

              // 强制更新视图
              this.$forceUpdate()
            } else {
              console.warn('⚠️ lengthComputable = false, 无法计算总进度')
              // 如果没有 Content-Length，只显示已接收大小
              this.downloadProgress.received = event.loaded
              console.log(`已接收: ${(event.loaded / 1024 / 1024).toFixed(2)} MB`)
            }
          })

          console.log('2. progress 事件监听器已绑定')

          xhr.open('GET', url, true)
          console.log('3. XMLHttpRequest 已 open')

          // 关键修复：使用 arraybuffer 而不是 blob，这样才能实时监听进度
          xhr.responseType = 'arraybuffer'
          console.log('4. responseType 设置为 arraybuffer')

          // 设置请求头
          if (tokenValue) {
            xhr.setRequestHeader('Authorization', `Bearer ${tokenValue}`)
            console.log('5. Authorization 头已设置')
          }

          // 请求完成
          xhr.addEventListener('load', () => {
            console.log('6. load 事件触发，status:', xhr.status)

            if (xhr.status === 200) {
              console.log('✓ 下载完成，创建文件...')

              // 从 ArrayBuffer 创建 Blob
              const blob = new Blob([xhr.response])
              console.log('7. Blob 已创建，大小:', blob.size, '字节')

              // 创建下载链接
              const downloadUrl = window.URL.createObjectURL(blob)
              const link = document.createElement('a')
              link.href = downloadUrl
              link.download = file.name
              document.body.appendChild(link)
              link.click()

              // 清理
              setTimeout(() => {
                document.body.removeChild(link)
                window.URL.revokeObjectURL(downloadUrl)
              }, 100)

              console.log('8. 下载链接已创建并触发')

              // 下载成功
              this.downloadProgress.percentage = 100
              this.downloadProgress.status = 'success'
              this.$message.success(`${file.name} 下载成功`)
              this.$emit('download-success')

              // 2秒后隐藏进度条
              setTimeout(() => {
                this.downloadProgress.show = false
              }, 2000)

              resolve()
            } else {
              console.error('✖ 下载失败，HTTP状态:', xhr.status)
              reject(new Error(`HTTP error! status: ${xhr.status}`))
            }
          })

          // 请求错误
          xhr.addEventListener('error', (e) => {
            console.error('✖ error 事件触发:', e)
            reject(new Error('网络请求失败'))
          })

          // 请求超时
          xhr.addEventListener('timeout', () => {
            console.error('✖ timeout 事件触发')
            reject(new Error('请求超时'))
          })

          // loadstart 事件 - 请求开始时触发
          xhr.addEventListener('loadstart', () => {
            console.log('✓ loadstart 事件触发 - 请求已开始')
          })

          // loadend 事件 - 请求结束时触发（无论成功或失败）
          xhr.addEventListener('loadend', () => {
            console.log('✓ loadend 事件触发 - 请求已结束')
          })

          // 设置超时时间（10分钟）
          xhr.timeout = 600000
          console.log('9. timeout 设置为 600000ms')

          console.log('10. 准备发送请求...')
          xhr.send()
          console.log('11. xhr.send() 已调用')
        })
      } catch (error) {
        console.error('文件下载失败:', error)
        console.error('错误堆栈:', error.stack)
        this.downloadProgress.status = 'exception'

        if (error.name === 'AbortError') {
          this.$message.error('文件下载已取消')
        } else if (error.message.includes('timeout') || error.message.includes('超时')) {
          this.$message.error('文件下载超时，请稍后重试或下载较小的文件')
        } else {
          this.$message.error('文件下载失败: ' + (error.message || '未知错误'))
        }

        // 3秒后隐藏进度条
        setTimeout(() => {
          this.downloadProgress.show = false
        }, 3000)
      } finally {
        this.downloadingFile = null
      }
    },

    // 格式化文件大小
    formatFileSize(bytes) {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
    }
  },
  beforeUnmount() {
    this.clearUploadTaskPolling()
  }
}
</script>

<style scoped>
/* 📁 文件浏览器样式 */
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

.current-path .btn-icon {
  width: 16px;
  height: 16px;
  display: block;
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

.current-path .el-text {
  color: #606266;
  font-weight: 500;
}

.upload-progress-section {
  margin-top: 16px;
}

.upload-progress-item + .upload-progress-item {
  margin-top: 16px;
}

.upload-progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  color: #606266;
  font-size: 14px;
}

.upload-progress-tip {
  margin-top: 12px;
  color: #909399;
  font-size: 13px;
}

.upload-progress-value {
  color: #409eff;
  font-weight: 600;
}

.upload-progress-value.is-success {
  color: #67c23a;
}

.upload-progress-value.is-failed {
  color: #f56c6c;
}
</style>
