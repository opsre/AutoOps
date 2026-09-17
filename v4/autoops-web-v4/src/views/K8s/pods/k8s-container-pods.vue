<template>
  <div class="container-pods">
    <!-- 容器组信息显示 -->
    <div class="pod-info-card">
      <div class="pod-header">
        <div class="pod-basic-info">
          <img src="@/assets/image/k8s.svg" alt="k8s" class="k8s-icon" />
          <div class="workload-info">
            <div 
              class="workload-name clickable-name" 
              @click="navigateToPodDetail"
            >
              {{ workload.name }}
            </div>
            <span class="workload-type-label">
              {{ getWorkloadTypeName(workload.type) }}
            </span>
          </div>
        </div>
      </div>

      <div class="pod-details">
        <!-- 标签信息 -->
        <div class="label-section">
          <div class="label-container">
            <el-badge :value="getVisibleLabelCount(workload.labels)" :max="99" class="label-badge">
              <el-button
                type="text"
                size="small"
                circle
                @click="viewWorkloadLabels"
                class="label-icon-button"
                title="查看标签"
              >
                <img src="@/assets/image/标签.svg" alt="标签" width="14" height="14" />
              </el-button>
            </el-badge>
          </div>
        </div>

        <!-- 容器组数量 -->
        <div class="pod-status-section">
          <div class="pod-status-container">
            <el-tag 
              :type="getPodStatusTagByReplicas(workload.readyReplicas, workload.totalReplicas)" 
              size="default"
              class="pod-count-tag"
              @click="$emit('view-pod-list')"
            >
              <el-icon class="pod-icon"><Monitor /></el-icon>
              {{ workload.replicas }}
            </el-tag>
            <div class="pod-status-text">
              <span :class="getReplicaStatusClass(workload.readyReplicas, workload.totalReplicas)">
                {{ workload.readyReplicas }}/{{ workload.totalReplicas }} 就绪
              </span>
            </div>
          </div>
        </div>

        <!-- 镜像信息 -->
        <div class="images-section">
          <div class="images-list">
            <div
              v-for="(image, index) in workload.images.slice(0, 1)"
              :key="index"
              class="image-tag-wrapper"
              @click="copyToClipboard(image, '镜像地址已复制')"
              title="点击复制镜像地址"
            >
              <el-icon class="copy-icon"><DocumentCopy /></el-icon>
              <span class="full-image-name">{{ image }}</span>
            </div>
            <el-button
              v-if="workload.images.length > 1"
              type="text"
              size="small"
              class="more-images-btn"
              @click="viewAllImages"
            >
              +{{ workload.images.length - 1 }} 更多
            </el-button>
          </div>
        </div>

        <!-- 状态和更新时间 -->
        <div class="status-section">
          <el-tag 
            :type="getStatusTag(workload.status)" 
            size="default"
            class="status-tag"
          >
            {{ getStatusText(workload.status) }}
          </el-tag>
          <div class="update-time">
            <span class="update-label">更新时间:</span>
            <span class="update-value">{{ formatDateTime(workload.updatedAt) }}</span>
          </div>
        </div>

        <!-- 监控按钮 -->
        <div class="monitoring-section">
          <el-button 
            :icon="Monitor" 
            size="small" 
            circle 
            type="primary" 
            @click="$emit('show-monitoring', workload)"
            title="查看监控"
            class="monitor-btn"
          />
        </div>
      </div>
    </div>

    <!-- 标签详情对话框 -->
    <el-dialog
      v-model="labelsDialogVisible"
      :title="`标签详情 - ${workload.name || ''}`"
      width="600px"
      class="labels-dialog"
    >
      <div class="labels-content">
        <!-- 用户标签 -->
        <div class="labels-section" v-if="Object.keys(getUserLabels(workload.labels)).length > 0">
          <h4>应用标签</h4>
          <div class="labels-grid">
            <el-tag
              v-for="(value, key) in getUserLabels(workload.labels)"
              :key="key"
              type="primary"
              size="default"
              class="label-tag"
              @click="copyToClipboard(`${key}=${value}`, '标签信息已复制')"
            >
              <el-icon class="tag-icon"><DocumentCopy /></el-icon>
              {{ key }}={{ value }}
            </el-tag>
          </div>
        </div>
        
        <!-- 系统标签 -->
        <div class="labels-section" v-if="Object.keys(getSystemLabels(workload.labels)).length > 0">
          <h4>系统标签</h4>
          <div class="labels-grid">
            <el-tag
              v-for="(value, key) in getSystemLabels(workload.labels)"
              :key="key"
              type="info"
              size="default"
              class="label-tag system-label"
              @click="copyToClipboard(`${key}=${value}`, '标签信息已复制')"
            >
              <el-icon class="tag-icon"><DocumentCopy /></el-icon>
              {{ key }}={{ value }}
            </el-tag>
          </div>
        </div>
        
        <!-- 没有标签的提示 -->
        <div v-if="!workload.labels || Object.keys(workload.labels).length === 0" class="no-labels">
          <el-empty description="该工作负载没有标签" :image-size="60" />
        </div>
      </div>
    </el-dialog>

    <!-- 镜像详情对话框 -->
    <el-dialog
      v-model="imagesDialogVisible"
      :title="`镜像信息 - ${workload.name || ''}`"
      width="800px"
      class="images-dialog"
    >
      <div class="images-content">
        <div v-if="workload.images && workload.images.length > 0" class="images-grid">
          <el-card
            v-for="(image, index) in workload.images"
            :key="index"
            class="image-card"
          >
            <div class="image-info">
              <div class="image-name-section">
                <div class="full-image-name">{{ image }}</div>
                <div class="image-actions">
                  <el-button
                    type="primary"
                    size="small"
                    :icon="DocumentCopy"
                    @click="copyToClipboard(image, '镜像地址已复制')"
                  >
                    复制
                  </el-button>
                </div>
              </div>
              <div class="image-details">
                <el-tag size="small" type="info">{{ getImageTag(image) }}</el-tag>
                <el-tag size="small" type="success">{{ getImageRegistry(image) }}</el-tag>
              </div>
            </div>
          </el-card>
        </div>
        
        <div v-else class="no-images">
          <el-empty description="该工作负载没有镜像信息" :image-size="60" />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Monitor, DocumentCopy } from '@element-plus/icons-vue'

// Props
const props = defineProps({
  workload: {
    type: Object,
    required: true
  },
  clusterId: {
    type: String,
    required: true
  },
  namespace: {
    type: String,
    required: true
  }
})

// Emits
const emit = defineEmits(['show-monitoring', 'view-pod-list'])

const router = useRouter()

// 对话框状态
const labelsDialogVisible = ref(false)
const imagesDialogVisible = ref(false)

// 工具函数
const getWorkloadTypeName = (type) => {
  const nameMap = {
    'deployment': 'Deployment',
    'statefulset': 'StatefulSet',
    'daemonset': 'DaemonSet',
    'job': 'Job',
    'cronjob': 'CronJob'
  }
  return nameMap[type] || type?.toUpperCase() || 'Unknown'
}

const getStatusTag = (status) => {
  const tagMap = {
    'Running': 'success',
    'Pending': 'warning',
    'Partial': 'warning',
    'Stopped': 'info',
    'Failed': 'danger',
    'Completed': 'success',
    'Active': 'success',
    'Suspended': 'info',
    'Unknown': 'info'
  }
  return tagMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    'Running': '运行中',
    'Pending': '等待中',
    'Partial': '部分就绪',
    'Stopped': '已停止',
    'Failed': '失败',
    'Completed': '已完成',
    'Active': '活跃',
    'Suspended': '暂停',
    'Unknown': '未知'
  }
  return textMap[status] || '未知'
}

const getPodStatusTagByReplicas = (ready, total) => {
  if (total === 0) return 'info'
  if (ready === 0) return 'danger'
  if (ready < total) return 'warning'
  return 'success'
}

const getReplicaStatusClass = (ready, total) => {
  if (total === 0) return 'replica-stopped'
  if (ready === 0) return 'replica-failed'
  if (ready < total) return 'replica-partial'
  return 'replica-ready'
}

const getVisibleLabelCount = (labels) => {
  if (!labels || typeof labels !== 'object') return 0
  return Object.keys(labels).length
}

const formatDateTime = (timestamp) => {
  if (!timestamp) return '-'
  
  const date = new Date(timestamp)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

const getImageTag = (image) => {
  const parts = image.split(':')
  return parts.length > 1 ? parts[parts.length - 1] : 'latest'
}

const getImageRegistry = (image) => {
  const parts = image.split('/')
  if (parts.length > 2) {
    return parts[0]
  } else if (parts.length === 2 && parts[0].includes('.')) {
    return parts[0]
  }
  return 'Docker Hub'
}

const getUserLabels = (labels) => {
  if (!labels) return {}
  
  const systemPrefixes = [
    'kubernetes.io/',
    'k8s.io/',
    'app.kubernetes.io/',
    'controller-revision-hash',
    'pod-template-hash'
  ]
  
  const userLabels = {}
  Object.keys(labels).forEach(key => {
    const isSystemLabel = systemPrefixes.some(prefix => key.startsWith(prefix))
    if (!isSystemLabel) {
      userLabels[key] = labels[key]
    }
  })
  
  return userLabels
}

const getSystemLabels = (labels) => {
  if (!labels) return {}
  
  const systemPrefixes = [
    'kubernetes.io/',
    'k8s.io/',
    'app.kubernetes.io/',
    'controller-revision-hash',
    'pod-template-hash'
  ]
  
  const systemLabels = {}
  Object.keys(labels).forEach(key => {
    const isSystemLabel = systemPrefixes.some(prefix => key.startsWith(prefix))
    if (isSystemLabel) {
      systemLabels[key] = labels[key]
    }
  })
  
  return systemLabels
}

const copyToClipboard = async (text, successMessage = '已复制到剪贴板') => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(successMessage)
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败')
  }
}

// 事件处理
const navigateToPodDetail = () => {
  // 这里可以导航到Pod详情页面
  ElMessage.info('跳转到Pod详情页面功能开发中...')
}

const viewWorkloadLabels = () => {
  labelsDialogVisible.value = true
}

const viewAllImages = () => {
  imagesDialogVisible.value = true
}
</script>

<style scoped>
.pod-info-card {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
  background: white;
}

.pod-header {
  padding: 16px;
  background: #f8f9fa;
  border-bottom: 1px solid #e4e7ed;
}

.pod-basic-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.k8s-icon {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
}

.workload-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.workload-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.clickable-name {
  cursor: pointer;
  color: #409eff;
}

.clickable-name:hover {
  text-decoration: underline;
}

.workload-type-label {
  font-size: 12px;
  color: #f39c12;
  font-weight: 500;
}

.pod-details {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.label-section {
  display: flex;
  align-items: center;
}

.label-badge {
  cursor: pointer;
}

.label-icon-button {
  border: none;
  background: transparent;
}

.label-icon-button:hover {
  background: rgba(64, 158, 255, 0.1);
}

.pod-status-section {
  display: flex;
  align-items: center;
}

.pod-status-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.pod-count-tag {
  cursor: pointer;
  font-size: 14px;
  padding: 8px 16px;
}

.pod-count-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.pod-icon {
  margin-right: 4px;
}

.pod-status-text {
  font-size: 12px;
  color: #606266;
}

.replica-ready { color: #67c23a; }
.replica-partial { color: #e6a23c; }
.replica-failed { color: #f56c6c; }
.replica-stopped { color: #909399; }

.images-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.images-list {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.image-tag-wrapper {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: #f0f9ff;
  border: 1px solid #3b82f6;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  max-width: 200px;
}

.image-tag-wrapper:hover {
  background: #dbeafe;
  transform: translateY(-1px);
}

.copy-icon {
  width: 12px;
  height: 12px;
  color: #3b82f6;
}

.full-image-name {
  font-size: 11px;
  color: #1e40af;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.more-images-btn {
  font-size: 12px;
  color: #606266;
}

.status-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.status-tag {
  font-size: 12px;
}

.update-time {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  font-size: 11px;
}

.update-label {
  color: #909399;
}

.update-value {
  color: #606266;
  font-weight: 500;
}

.monitoring-section {
  display: flex;
  justify-content: center;
}

.monitor-btn {
  width: 36px;
  height: 36px;
}

/* 标签对话框样式 */
.labels-content {
  min-height: 200px;
}

.labels-section {
  margin-bottom: 24px;
}

.labels-section h4 {
  margin: 0 0 12px 0;
  color: #303133;
  font-size: 14px;
  font-weight: 600;
}

.labels-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.label-tag {
  cursor: pointer;
  transition: all 0.2s;
}

.label-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.tag-icon {
  margin-right: 4px;
  font-size: 12px;
}

.system-label {
  opacity: 0.8;
}

.no-labels {
  padding: 40px;
  text-align: center;
}

/* 镜像对话框样式 */
.images-content {
  min-height: 200px;
}

.images-grid {
  display: grid;
  gap: 12px;
}

.image-card {
  border: 1px solid #e4e7ed;
}

.image-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.image-name-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.image-actions {
  flex-shrink: 0;
}

.image-details {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.no-images {
  padding: 40px;
  text-align: center;
}
</style>