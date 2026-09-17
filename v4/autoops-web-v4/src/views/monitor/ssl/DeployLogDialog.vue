<template>
  <el-dialog
    v-model="isVisible"
    title="部署日志"
    width="900px"
    @close="handleClose"
  >
    <!-- 证书信息 -->
    <div class="cert-info" v-if="cert">
      <span>证书域名: {{ cert.domain }} (ID: {{ cert.id }})</span>
    </div>

    <!-- 日志列表 -->
    <el-table
      :data="logs"
      border
      stripe
      :loading="loading"
      max-height="500px"
    >
      <el-table-column prop="id" label="ID" width="80" align="center" />
      <el-table-column prop="hostName" label="主机" width="150" />
      <el-table-column prop="deployPath" label="部署路径" width="180" show-overflow-tooltip />
      <el-table-column label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">{{ row.statusText }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="部署时间" width="160" align="center" />
      <el-table-column label="操作" width="120" fixed="right" align="center">
        <template #default="{ row }">
          <el-tooltip content="查看日志" placement="top">
            <el-button
              type="primary"
              :icon="View"
              size="small"
              @click="handleShowDetail(row)"
            >日志详情</el-button>
          </el-tooltip>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="currentPageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </div>

    <!-- 日志详情对话框 -->
    <el-dialog v-model="detailVisible" title="部署详情" width="800px" append-to-body>
      <div v-if="selectedLog">
        <!-- 部署信息 -->
        <el-descriptions :column="2" border>
          <el-descriptions-item label="主机">{{ selectedLog.hostName }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(selectedLog.status)" size="small">{{ selectedLog.statusText }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="部署路径" :span="2">{{ selectedLog.deployPath }}</el-descriptions-item>
          <el-descriptions-item label="部署时间" :span="2">{{ selectedLog.createTime }}</el-descriptions-item>
        </el-descriptions>

        <!-- 备份文件 -->
        <div v-if="selectedLog.backupFiles && selectedLog.backupFiles.length > 0" style="margin-top: 20px;">
          <h4>备份文件</h4>
          <el-tag v-for="(file, index) in selectedLog.backupFiles" :key="index" style="margin: 4px;">
            {{ file }}
          </el-tag>
        </div>

        <!-- 部署文件 -->
        <div v-if="selectedLog.deployFiles && selectedLog.deployFiles.length > 0" style="margin-top: 20px;">
          <h4>部署文件</h4>
          <el-tag v-for="(file, index) in selectedLog.deployFiles" :key="index" type="success" style="margin: 4px;">
            {{ file }}
          </el-tag>
        </div>

        <!-- 错误信息 -->
        <div v-if="selectedLog.errorMsg" style="margin-top: 20px;">
          <h4>错误信息</h4>
          <el-alert :title="selectedLog.errorMsg" type="error" :closable="false" />
        </div>

        <!-- 部署日志 -->
        <div style="margin-top: 20px;">
          <h4>部署日志</h4>
          <div class="log-content">
            <div v-for="(log, index) in selectedLog.logs" :key="index" class="log-line">
              {{ log }}
            </div>
          </div>
        </div>
      </div>
    </el-dialog>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { View } from '@element-plus/icons-vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  cert: {
    type: Object,
    default: null
  },
  logs: {
    type: Array,
    default: () => []
  },
  total: {
    type: Number,
    default: 0
  },
  page: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'close', 'page-change', 'size-change'])

const isVisible = ref(false)
const currentPage = ref(1)
const currentPageSize = ref(10)
const detailVisible = ref(false)
const selectedLog = ref(null)

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

// 监听 page 和 pageSize prop
watch(() => props.page, (newVal) => {
  currentPage.value = newVal
})

watch(() => props.pageSize, (newVal) => {
  currentPageSize.value = newVal
})

// 获取状态类型
const getStatusType = (status) => {
  // 1=部署中, 2=成功, 3=失败
  switch (status) {
    case 1:
      return 'warning'  // 部署中
    case 2:
      return 'success'  // 成功
    case 3:
      return 'danger'   // 失败
    default:
      return 'info'
  }
}

// 显示日志详情
const handleShowDetail = (log) => {
  selectedLog.value = log
  detailVisible.value = true
}

// 分页变化
const handlePageChange = (page) => {
  currentPage.value = page
  emit('page-change', page)
}

const handleSizeChange = (size) => {
  currentPageSize.value = size
  emit('size-change', size)
}

// 关闭对话框
const handleClose = () => {
  isVisible.value = false
  emit('close')
}
</script>

<style scoped>
.cert-info {
  margin-bottom: 20px;
  padding: 12px 16px;
  background: linear-gradient(135deg, rgba(103, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border-left: 4px solid #667eea;
  border-radius: 8px;
  font-weight: 500;
  color: #2c3e50;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.log-content {
  max-height: 400px;
  overflow-y: auto;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  font-family: 'Courier New', Consolas, monospace;
  font-size: 13px;
}

.log-line {
  padding: 4px 0;
  line-height: 1.6;
  color: #606266;
  white-space: pre-wrap;
  word-break: break-all;
}

.log-line:hover {
  background: rgba(64, 158, 255, 0.1);
  padding-left: 4px;
  transition: all 0.2s ease;
}

h4 {
  margin: 0 0 12px 0;
  padding-bottom: 8px;
  border-bottom: 2px solid #e4e7ed;
  color: #2c3e50;
  font-size: 15px;
  font-weight: 600;
}

:deep(.el-descriptions__label) {
  font-weight: 600;
  color: #606266;
}

:deep(.el-descriptions__content) {
  color: #303133;
}
</style>
