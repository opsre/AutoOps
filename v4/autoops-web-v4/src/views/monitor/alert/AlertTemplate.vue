<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Edit,
  Delete,
  View,
  DeleteFilled
} from '@element-plus/icons-vue'
import { templateApi } from '@/api/alert'

// 定义事件
const emit = defineEmits(['refresh-stats'])

// 左侧分类选中状态
const selectedCategory = ref('all')

// 搜索表单
const searchForm = reactive({
  name: '',
  status: -1,
  pageNum: 1,
  pageSize: 10
})

// 列表数据
const tableData = ref([])
const total = ref(0)
const loading = ref(false)
const selectedRows = ref([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('新增模板')
const formLoading = ref(false)

// 表单数据
const templateForm = reactive({
  id: null,
  name: '',
  category: '',
  description: '',
  metricKey: '',      // 指标键名 (如: system_cpu_usage_percent)
  unit: '',           // 单位 (如: percentage, GB, MB)
  valueType: '',      // 值类型 (number, string, boolean)
  status: 1
})

const formRef = ref(null)
const categoryList = ref([])

// 分类选项 (与后端一致)
const categoryOptions = [
  { label: '全部模板', value: 'all', count: 0 },
  { label: 'CPU', value: 'cpu', count: 0 },
  { label: '内存', value: 'memory', count: 0 },
  { label: '负载', value: 'load', count: 0 },
  { label: '磁盘', value: 'disk', count: 0 },
  { label: 'IO', value: 'io', count: 0 },
  { label: '网络', value: 'network', count: 0 },
  { label: '连接数', value: 'connection', count: 0 },
  { label: '服务', value: 'service', count: 0 },
  { label: '自定义', value: 'custom', count: 0 }
]

const statusOptions = [
  { label: '全部', value: -1 },
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

// 表单验证规则
const formRules = {
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  metricKey: [{ required: true, message: '请输入指标键', trigger: 'blur' }],
  unit: [{ required: true, message: '请输入单位', trigger: 'blur' }],
  valueType: [{ required: true, message: '请选择值类型', trigger: 'change' }]
}

// 获取各分类的模板数量
const fetchCategoryCounts = async () => {
  try {
    const categories = ['all', 'cpu', 'memory', 'load', 'disk', 'io', 'network', 'connection', 'service', 'custom']
    const promises = categories.map(async (cat) => {
      const params = {
        pageNum: 1,
        pageSize: 1,
        status: 1
      }
      if (cat !== 'all') {
        params.category = cat
      }
      const { data: res } = await templateApi.getList(params)
      return {
        category: cat,
        count: res.data?.total || 0
      }
    })

    const results = await Promise.all(promises)
    results.forEach(result => {
      const cat = categoryOptions.find(c => c.value === result.category)
      if (cat) {
        cat.count = result.count
      }
    })
  } catch (error) {
    console.error('获取分类统计失败:', error)
  }
}

// 获取分类列表 (从后端)
const fetchCategories = async () => {
  try {
    const { data: res } = await templateApi.getCategories()
    if (res.code === 200) {
      // 如果后端返回分类列表，使用后端数据，否则使用本地默认
      categoryList.value = res.data && res.data.length > 0 ? res.data : categoryOptions.map(c => c.value)
    }
  } catch (error) {
    // 如果后端错误，使用本地默认分类
    console.error('获取分类列表失败，使用默认分类:', error)
    categoryList.value = categoryOptions.map(c => c.value)
  }
}

// 获取模板列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = { ...searchForm }
    if (params.status === -1) delete params.status

    // 如果选中了分类，添加category参数
    if (selectedCategory.value !== 'all') {
      params.category = selectedCategory.value
    }

    const { data: res } = await templateApi.getList(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0

      // 更新分类统计
      await fetchCategoryCounts()
    }
  } catch (error) {
    ElMessage.error('获取模板列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 选择分类
const handleCategorySelect = (category) => {
  selectedCategory.value = category
  searchForm.pageNum = 1
  fetchList()
}

// 搜索
const handleSearch = () => {
  searchForm.pageNum = 1
  fetchList()
}

// 重置
const handleReset = () => {
  searchForm.name = ''
  searchForm.status = -1
  searchForm.pageNum = 1
  handleSearch()
}

// 分页
const handlePageChange = (page) => {
  searchForm.pageNum = page
  fetchList()
}

const handleSizeChange = (size) => {
  searchForm.pageSize = size
  searchForm.pageNum = 1
  fetchList()
}

// 新增
const handleAdd = () => {
  dialogTitle.value = '新增模板'
  templateForm.id = null
  templateForm.name = ''
  templateForm.category = ''
  templateForm.description = ''
  templateForm.metricKey = ''
  templateForm.unit = ''
  templateForm.valueType = ''
  templateForm.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = async (row) => {
  try {
    const { data: res } = await templateApi.getDetail(row.id)
    if (res.code === 200) {
      dialogTitle.value = '编辑模板'
      const detail = res.data
      templateForm.id = detail.id
      templateForm.name = detail.name
      templateForm.category = detail.category
      templateForm.description = detail.description || ''
      templateForm.metricKey = detail.metricKey || ''
      templateForm.unit = detail.unit || ''
      templateForm.valueType = detail.valueType || ''
      templateForm.status = detail.status
      dialogVisible.value = true
    } else {
      ElMessage.error(res.message || '获取模板详情失败')
    }
  } catch (error) {
    ElMessage.error('获取模板详情失败')
    console.error(error)
  }
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该模板吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const { data: res } = await templateApi.delete(row.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      fetchList()
      emit('refresh-stats')
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error(error)
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择要删除的模板')
    return
  }

  try {
    await ElMessageBox.confirm(`确定要删除 ${selectedRows.value.length} 个模板吗?`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const ids = selectedRows.value.map(row => row.id)
    const { data: res } = await templateApi.batchDelete(ids)
    if (res.code === 200) {
      ElMessage.success('批量删除成功')
      selectedRows.value = []
      fetchList()
      emit('refresh-stats')
    } else {
      ElMessage.error(res.message || '批量删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
      console.error(error)
    }
  }
}

// 状态切换
const handleStatusChange = async (row) => {
  try {
    const { data: res } = await templateApi.updateStatus(row.id, row.status)
    if (res.code === 200) {
      ElMessage.success(row.status === 1 ? '启用成功' : '禁用成功')
      emit('refresh-stats')
    } else {
      row.status = row.status === 1 ? 0 : 1 // 恢复原状态
      ElMessage.error(res.message || '操作失败')
    }
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1 // 恢复原状态
    ElMessage.error('操作失败')
    console.error(error)
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    formLoading.value = true
    try {
      const data = { ...templateForm }

      let result
      if (data.id) {
        result = await templateApi.update(data)
      } else {
        delete data.id
        result = await templateApi.create(data)
      }

      const { data: res } = result
      if (res.code === 200) {
        ElMessage.success(data.id ? '更新成功' : '创建成功')
        dialogVisible.value = false
        fetchList()
        emit('refresh-stats')
      } else {
        ElMessage.error(res.message || '操作失败')
      }
    } catch (error) {
      ElMessage.error('操作失败')
      console.error(error)
    } finally {
      formLoading.value = false
    }
  })
}

// 获取分类图标颜色
const getCategoryColor = (category) => {
  const colorMap = {
    'cpu': '#FF6B6B',
    'memory': '#4ECDC4',
    'disk': '#45B7D1',
    'io': '#96CEB4',
    'network': '#FFEAA7',
    'connection': '#DFE6E9',
    'service': '#A29BFE',
    'custom': '#74B9FF'
  }
  return colorMap[category] || '#909399'
}

// 行选择
const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

onMounted(() => {
  fetchCategories()
  fetchList()
  fetchCategoryCounts()
})
</script>

<template>
  <div class="alert-template-container">
    <!-- 主内容区域：左右布局 -->
    <div class="main-content">
      <!-- 左侧分类栏 -->
      <div class="category-sidebar">
        <div class="category-header">模板分类</div>
        <div class="category-list">
          <div
            v-for="category in categoryOptions"
            :key="category.value"
            :class="['category-item', { active: selectedCategory === category.value }]"
            @click="handleCategorySelect(category.value)"
          >
            <span class="category-label">{{ category.label }}</span>
            <span class="category-count">({{ category.count }})</span>
          </div>
        </div>
      </div>

      <!-- 右侧内容区域 -->
      <div class="content-area">
        <!-- 搜索区域 -->
        <div class="search-section">
          <el-form :inline="true" :model="searchForm" class="search-form">
            <el-form-item label="模板名称">
              <el-input
                v-model="searchForm.name"
                placeholder="请输入模板名称"
                size="small"
                clearable
                @keyup.enter="handleSearch"
              />
            </el-form-item>

            <el-form-item label="状态">
              <el-select v-model="searchForm.status" size="small" style="width: 120px">
                <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" size="small" @click="handleSearch" class="modern-btn primary-btn">
                <el-icon><Search /></el-icon>
                搜索
              </el-button>
              <el-button size="small" @click="handleReset" class="modern-btn reset-btn">
                <el-icon><Refresh /></el-icon>
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 操作按钮 -->
        <div class="action-section">
          <el-button
            type="success"
            size="small"
            @click="handleAdd"
            v-authority="['monitor:alert:template:add']"
            class="modern-btn success-btn"
          >
            <el-icon><Plus /></el-icon>
            新增模板
          </el-button>

          <el-button
            type="danger"
            size="small"
            @click="handleBatchDelete"
            :disabled="selectedRows.length === 0"
            v-authority="['monitor:alert:template:delete']"
            class="modern-btn danger-btn"
          >
            <el-icon><DeleteFilled /></el-icon>
            批量删除
          </el-button>
        </div>

        <!-- 表格 -->
        <div class="table-section">
          <el-table
            :data="tableData"
            v-loading="loading"
            stripe
            border
            style="width: 100%"
            :header-cell-style="{ background: '#eef1f6', color: '#606266' }"
            @selection-change="handleSelectionChange"
          >
        <el-table-column type="selection" width="50" align="center" />

        <el-table-column prop="id" label="ID" width="70" align="center" />

        <el-table-column prop="name" label="模板名称" min-width="150">
          <template #default="{ row }">
            <span style="font-weight: 500">{{ row.name }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="category" label="分类" width="120" align="center">
          <template #default="{ row }">
            <el-tag
              :style="{ backgroundColor: getCategoryColor(row.category), color: '#fff' }"
              size="small"
            >
              {{ row.category }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="指标键" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <span style="font-family: monospace; font-size: 12px">{{ row.metricKey }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="unit" label="单位" width="100" align="center">
          <template #default="{ row }">
            <span style="font-size: 12px">{{ row.unit }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="valueType" label="值类型" width="100" align="center">
          <template #default="{ row }">
            <span style="font-size: 12px">{{ row.valueType }}</span>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
              v-authority="['monitor:alert:template:status']"
            />
          </template>
        </el-table-column>

        <el-table-column prop="createdAt" label="创建时间" width="170" align="center">
          <template #default="{ row }">
            <span style="font-size: 12px">{{ row.createdAt }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" fixed="right" width="250" align="center">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-tooltip content="查看" placement="top">
                <el-button
                  type="primary"
                  :icon="View"
                  size="small"
                  @click="handleEdit(row)"
                  v-authority="['monitor:alert:template:edit']"
                >查看</el-button>
              </el-tooltip>
              <el-tooltip content="修改" placement="top">
                <el-button
                  type="warning"
                  :icon="Edit"
                  size="small"
                  @click="handleEdit(row)"
                  v-authority="['monitor:alert:template:edit']"
                >修改</el-button>
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button
                  type="danger"
                  :icon="Delete"
                  size="small"
                  @click="handleDelete(row)"
                  v-authority="['monitor:alert:template:delete']"
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
            v-model:current-page="searchForm.pageNum"
            v-model:page-size="searchForm.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
          />
        </div>
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="800px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="templateForm"
        :rules="formRules"
        label-width="120px"
      >
        <!-- 基础信息 -->
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="templateForm.name" placeholder="请输入模板名称" />
        </el-form-item>

        <el-form-item label="分类" prop="category">
          <el-select v-model="templateForm.category" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="item in categoryOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="templateForm.description" type="textarea" :rows="2" placeholder="请输入指标描述" />
        </el-form-item>

        <!-- 指标信息 -->
        <el-form-item label="指标键" prop="metricKey">
          <el-input
            v-model="templateForm.metricKey"
            placeholder="例如: system_cpu_usage_percent"
            hint="指标的唯一标识符"
          />
        </el-form-item>

        <el-form-item label="单位" prop="unit">
          <el-input v-model="templateForm.unit" placeholder="例如: percentage, GB, MB, ms" />
        </el-form-item>

        <el-form-item label="值类型" prop="valueType">
          <el-select v-model="templateForm.valueType" placeholder="请选择值类型" style="width: 100%">
            <el-option label="数值" value="number" />
            <el-option label="字符串" value="string" />
            <el-option label="布尔值" value="boolean" />
          </el-select>
        </el-form-item>

        <el-form-item label="状态">
          <el-radio-group v-model="templateForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="formLoading">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.alert-template-container {
  padding: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 主内容左右布局 */
.main-content {
  display: flex;
  gap: 16px;
  flex: 1;
  min-height: 0;
}

/* 左侧分类栏 */
.category-sidebar {
  width: 180px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  overflow: hidden;
  flex-shrink: 0;
}

.category-header {
  padding: 12px 16px;
  background: #f5f7fa;
  color: #303133;
  font-size: 14px;
  font-weight: 600;
  border-bottom: 1px solid #e4e7ed;
}

.category-list {
  padding: 8px 0;
}

.category-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 14px;
  color: #606266;
}

.category-item:hover {
  background: #f5f7fa;
}

.category-item.active {
  background: #ecf5ff;
  color: #409eff;
  font-weight: 500;
  border-right: 3px solid #409eff;
}

.category-label {
  flex: 1;
}

.category-count {
  font-size: 12px;
  color: #909399;
  margin-left: 8px;
}

.category-item.active .category-count {
  color: #409eff;
}

/* 右侧内容区域 */
.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.search-section {
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form :deep(.el-form-item) {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form :deep(.el-form-item__label) {
  color: #606266;
  font-weight: 500;
}

.action-section {
  display: flex;
  gap: 10px;
}

.table-section {
  flex: 1;
  min-height: 0;
  overflow: auto;
}

.pagination-section {
  display: flex;
  justify-content: center;
  padding-top: 12px;
}

:deep(.el-table) {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  font-size: 12px;
}

:deep(.el-table__header) {
  background: rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-table__header th) {
  background: rgba(102, 126, 234, 0.1) !important;
  color: #2c3e50 !important;
  font-weight: 600;
  border: none;
  padding: 8px 0;
  font-size: 13px;
}

:deep(.el-table__body tr:hover > td) {
  background-color: rgba(102, 126, 234, 0.08) !important;
}

:deep(.el-table td) {
  border: none;
  padding: 6px 0;
}

:deep(.el-table .cell) {
  padding: 0 8px;
}

.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.operation-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
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

:deep(.danger-btn) {
  background: linear-gradient(45deg, #F56C6C, #F78989);
  color: white;
}
</style>
