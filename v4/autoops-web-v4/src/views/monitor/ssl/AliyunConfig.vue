<template>
  <div class="aliyun-config-container">
    <!-- 搜索和操作区域 -->
    <div class="search-section">
      <el-form :inline="true" :model="aliyunConfigSearchForm" class="search-form">
        <el-form-item label="配置名称">
          <el-input
            v-model="aliyunConfigSearchForm.name"
            placeholder="搜索配置名称"
            size="small"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>

        <el-form-item label="配置状态">
          <el-select
            v-model="aliyunConfigSearchForm.status"
            placeholder="配置状态"
            size="small"
            clearable
            style="width: 120px"
          >
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" size="small" @click="handleSearch" class="modern-btn primary-btn">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button size="small" @click="handleReset" class="modern-btn reset-btn">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 操作按钮区域 -->
    <div class="action-section">
      <el-button
        type="success"
        size="small"
        @click="handleAddAliyunConfig"
        class="modern-btn success-btn"
      >
        <el-icon><Plus /></el-icon>
        添加配置
      </el-button>

      <el-button
        size="small"
        @click="fetchAliyunConfigList"
        class="modern-btn info-btn"
      >
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <!-- 列表 -->
    <div class="table-section">
      <el-table
        :data="aliyunConfigList"
        border
        stripe
        :loading="aliyunConfigLoading"
        class="config-table"
        :header-cell-style="{ background: '#eef1f6', color: '#606266' }"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="配置名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="accessKey" label="AccessKey" min-width="200" show-overflow-tooltip />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="150" show-overflow-tooltip />
        <el-table-column prop="phone" label="电话" width="130" align="center" />
        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-tooltip content="编辑" placement="top">
                <el-button
                  type="primary"
                  :icon="Edit"
                  size="small"
                  @click="handleEditAliyunConfig(row)"
                >编辑</el-button>
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button
                  type="danger"
                  :icon="Delete"
                  size="small"
                  @click="handleDeleteAliyunConfig(row)"
                >删除</el-button>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 分页 -->
    <div class="pagination-section">
      <el-pagination
        v-model:current-page="aliyunConfigSearchForm.pageNum"
        v-model:page-size="aliyunConfigSearchForm.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="aliyunConfigTotal"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </div>

    <!-- 编辑对话框 -->
    <aliyun-config-dialog
      v-model:visible="aliyunConfigDialogVisible"
      :title="aliyunConfigDialogTitle"
      :form="aliyunConfigForm"
      :loading="aliyunConfigLoading"
      @submit="handleSubmitAliyunConfig"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh, Edit, Delete } from '@element-plus/icons-vue'
import { aliyunConfigApi } from '@/api/monitor'
import AliyunConfigDialog from './AliyunConfigDialog.vue'

// 阿里云配置搜索表单
const aliyunConfigSearchForm = reactive({
  name: '',
  status: -1,
  pageNum: 1,
  pageSize: 10
})

// 阿里云配置表单
const aliyunConfigForm = reactive({
  id: null,
  name: '',
  accessKey: '',
  accessSecret: '',
  email: '',
  phone: '',
  username: '',
  status: 1
})

// 阿里云配置列表状态
const aliyunConfigList = ref([])
const aliyunConfigTotal = ref(0)
const aliyunConfigLoading = ref(false)
const aliyunConfigDialogVisible = ref(false)
const aliyunConfigDialogTitle = ref('添加阿里云配置')

// 获取阿里云配置列表
const fetchAliyunConfigList = async () => {
  aliyunConfigLoading.value = true
  try {
    const { data: res } = await aliyunConfigApi.getList(aliyunConfigSearchForm)
    if (res.code === 200) {
      aliyunConfigList.value = Array.isArray(res.data) ? res.data : res.data?.list || []
      aliyunConfigTotal.value = res.data?.total || aliyunConfigList.value.length
    } else {
      ElMessage.error(res.message || '获取配置列表失败')
    }
  } catch (error) {
    console.error('获取配置列表失败:', error)
    ElMessage.error('获取配置列表失败')
  } finally {
    aliyunConfigLoading.value = false
  }
}

// 搜索
const handleSearch = () => {
  aliyunConfigSearchForm.pageNum = 1
  fetchAliyunConfigList()
}

// 重置
const handleReset = () => {
  aliyunConfigSearchForm.name = ''
  aliyunConfigSearchForm.status = -1
  aliyunConfigSearchForm.pageNum = 1
  fetchAliyunConfigList()
}

// 分页变化
const handlePageChange = (page) => {
  aliyunConfigSearchForm.pageNum = page
  fetchAliyunConfigList()
}

const handleSizeChange = (size) => {
  aliyunConfigSearchForm.pageSize = size
  fetchAliyunConfigList()
}

// 打开添加对话框
const handleAddAliyunConfig = () => {
  aliyunConfigDialogTitle.value = '添加阿里云配置'
  aliyunConfigForm.id = null
  aliyunConfigForm.name = ''
  aliyunConfigForm.accessKey = ''
  aliyunConfigForm.accessSecret = ''
  aliyunConfigForm.email = ''
  aliyunConfigForm.phone = ''
  aliyunConfigForm.username = ''
  aliyunConfigForm.status = 1
  aliyunConfigDialogVisible.value = true
}

// 打开编辑对话框
const handleEditAliyunConfig = (row) => {
  aliyunConfigDialogTitle.value = '编辑阿里云配置'
  aliyunConfigForm.id = row.id
  aliyunConfigForm.name = row.name
  aliyunConfigForm.accessKey = row.accessKey
  aliyunConfigForm.accessSecret = row.accessSecret || ''
  aliyunConfigForm.email = row.email
  aliyunConfigForm.phone = row.phone
  aliyunConfigForm.username = row.username
  aliyunConfigForm.status = row.status
  aliyunConfigDialogVisible.value = true
}

// 提交表单
const handleSubmitAliyunConfig = async () => {
  aliyunConfigLoading.value = true
  try {
    const data = { ...aliyunConfigForm }
    let result
    if (data.id) {
      result = await aliyunConfigApi.update(data)
    } else {
      delete data.id
      result = await aliyunConfigApi.create(data)
    }

    const { data: res } = result
    if (res.code === 200) {
      ElMessage.success(data.id ? '更新成功' : '添加成功')
      aliyunConfigDialogVisible.value = false
      fetchAliyunConfigList()
    } else {
      ElMessage.error(res.message || '操作失败')
    }
  } catch (error) {
    console.error('操作失败:', error)
    ElMessage.error('操作失败')
  } finally {
    aliyunConfigLoading.value = false
  }
}

// 删除
const handleDeleteAliyunConfig = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该配置吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const { data: res } = await aliyunConfigApi.delete(row.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      fetchAliyunConfigList()
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 工具函数
const getStatusType = (status) => {
  return status === 1 ? 'success' : 'info'
}

// 组件挂载
onMounted(() => {
  fetchAliyunConfigList()
})
</script>

<style scoped>
.aliyun-config-container {
  padding: 0;
}

.search-section {
  margin-bottom: 16px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form {
  margin: 0;
}

.search-form :deep(.el-form-item) {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form :deep(.el-form-item__label) {
  color: #606266;
  font-weight: 500;
}

.search-form :deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.search-form :deep(.el-select .el-input__wrapper) {
  border-radius: 8px;
}

.action-section {
  margin-bottom: 16px;
}

:deep(.modern-btn) {
  border-radius: 8px;
  padding: 8px 20px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

:deep(.modern-btn:hover) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
}

:deep(.primary-btn) {
  background: linear-gradient(45deg, #409EFF, #66B3FF);
  color: white;
}

:deep(.success-btn) {
  background: linear-gradient(45deg, #67C23A, #85CE61);
  color: white;
}

:deep(.reset-btn) {
  background: linear-gradient(45deg, #E6A23C, #EEBE77);
  color: white;
}

:deep(.info-btn) {
  background: linear-gradient(45deg, #909399, #b1b3b8);
  color: white;
}

.table-section {
  margin-top: 0;
}

.config-table {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  font-size: 12px;
}

.config-table :deep(.el-table__header) {
  background: rgba(102, 126, 234, 0.1) !important;
}

.config-table :deep(.el-table__header th) {
  background: rgba(102, 126, 234, 0.1) !important;
  color: #2c3e50 !important;
  font-weight: 600;
  border: none;
  padding: 8px 0;
  font-size: 13px;
}

.config-table :deep(.el-table__body tr:hover > td) {
  background-color: rgba(102, 126, 234, 0.08) !important;
}

.config-table :deep(.el-table td) {
  border: none;
  padding: 6px 0;
}

.config-table :deep(.el-table .cell) {
  padding: 0 8px;
  line-height: 1.4;
}

.config-table :deep(.el-table::before) {
  display: none;
}

.config-table :deep(.el-table--border::after) {
  display: none;
}

/* 确保操作列固定在最上层 */
.config-table :deep(.el-table__fixed-right) {
  z-index: 3 !important;
}

.config-table :deep(.el-table__fixed-right-patch) {
  z-index: 3 !important;
}

/* 操作按钮样式 */
.operation-buttons {
  display: flex;
  gap: 6px;
  justify-content: center;
  align-items: center;
}

.operation-buttons .el-button {
  transition: all 0.3s ease;
}

.operation-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.config-table :deep(.el-button.is-link) {
  font-weight: 500;
  transition: all 0.3s ease;
}

.config-table :deep(.el-button.is-link:hover) {
  transform: translateY(-1px);
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

:deep(.el-pagination) {
  padding: 0;
}

:deep(.el-pagination button),
:deep(.el-pagination .el-pager li) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.el-pagination button:hover),
:deep(.el-pagination .el-pager li:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

:deep(.el-pagination .el-pager li.is-active) {
  background: linear-gradient(45deg, #409EFF, #66B3FF);
  color: white;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}
</style>
