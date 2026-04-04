<template>
  <div class="knowledge-container">
    <div class="card">
      <!-- 列表页 -->
      <div v-if="viewMode === 'list'" class="fade-in list-view">
        <!-- 左侧分类面板 -->
        <div class="category-panel">
          <div class="panel-header">
            <span>文档分类</span>
            <el-icon class="add-icon" @click="openCategoryDialog()"><Plus /></el-icon>
          </div>
          <div class="category-list">
            <div 
              class="category-item" 
              :class="{ active: currentCategory === '' }"
              @click="selectCategory('')"
            >
              <span><el-icon><Folder /></el-icon> 全部</span>
              <span class="doc-count">{{ totalDocs }}</span>
            </div>
            <div 
              v-for="cat in categoryList" 
              :key="cat.id"
              class="category-item" 
              :class="{ active: currentCategory === cat.name }"
              @click="selectCategory(cat.name)"
            >
              <span><el-icon><Folder /></el-icon> {{ cat.name }}</span>
              <div class="category-actions">
                <span class="doc-count">{{ cat.docCount }}</span>
                <el-icon class="action-icon" @click.stop="openCategoryDialog(cat)"><Edit /></el-icon>
                <el-icon class="action-icon delete" @click.stop="handleDeleteCategory(cat.id)"><Delete /></el-icon>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧表格区域 -->
        <div class="table-panel">
          <div class="search-row">
            <el-input v-model="searchForm.title" placeholder="输入关键字搜索..." style="width: 320px" clearable @keyup.enter="handleSearch" @clear="handleSearch" />
            <el-select v-model="searchForm.status" placeholder="选择状态" clearable style="width: 120px" @change="handleSearch">
              <el-option label="已发布" :value="1" />
              <el-option label="草稿" :value="2" />
            </el-select>
            <el-button type="primary" @click="handleSearch">
              <el-icon><Search /></el-icon>查询
            </el-button>
            <el-button @click="handleReset">
              <el-icon><Refresh /></el-icon>重置
            </el-button>
            <div style="flex: 1"></div>
            <el-button type="success" @click="goEdit(null)">
              <el-icon><Plus /></el-icon>新增知识
            </el-button>
          </div>

          <el-table :data="tableData" v-loading="loading" stripe style="width: 100%">
            <el-table-column type="index" label="序号" width="80" />
            <el-table-column prop="title" label="知识库标题" min-width="300">
              <template #default="{ row }">
                <el-link type="primary" @click="showDetail(row.id)" style="font-weight: 500">{{ row.title }}</el-link>
              </template>
            </el-table-column>
            <el-table-column prop="category" label="所属分类" width="120">
              <template #default="{ row }">
                <span class="category-tag">{{ row.category }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="author" label="作者" width="100" />
            <el-table-column prop="createTime" label="创建时间" width="160" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.statusText }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180" align="center">
              <template #default="{ row }">
                <el-button type="primary" link @click="goEdit(row.id)">编辑</el-button>
                <el-button type="danger" link @click="handleDelete(row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-wrapper">
            <el-pagination
              v-model:current-page="pagination.pageNum"
              v-model:page-size="pagination.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="pagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>
      </div>

      <!-- 编辑/新增页 -->
      <div v-if="viewMode === 'edit'" class="fade-in edit-view">
        <div class="edit-layout">
          <div class="edit-header">
            <el-input
              v-model="editForm.title"
              class="title-input"
              placeholder="请输入知识标题..."
            />
            <div style="display: flex; gap: 12px; margin-left: 40px">
              <el-button
                :type="isPreviewOpen ? 'primary' : 'default'"
                @click="togglePreview"
              >
                <el-icon><View /></el-icon> {{ isPreviewOpen ? '关闭预览' : '实时预览' }}
              </el-button>
              <el-button @click="goList">取消</el-button>
              <el-button type="primary" @click="handleSave(1)">保存并发布</el-button>
            </div>
          </div>

          <div class="category-select-row">
            <span class="label">选择分类：</span>
            <el-select v-model="editForm.category" placeholder="选择分类" style="width: 180px">
              <el-option v-for="cat in categoryList" :key="cat.id" :label="cat.name" :value="cat.name" />
            </el-select>
          </div>

          <div class="editor-container">
            <div 
              class="monaco-editor-pane" 
              :style="{ width: isPreviewOpen ? editorWidth + '%' : '100%' }"
            >
              <vue-monaco-editor
                v-model:value="editForm.content"
                language="markdown"
                :options="monacoOptions"
                @change="handleEditorChange"
              />
            </div>
            
            <div 
              v-if="isPreviewOpen" 
              class="resize-divider"
              @mousedown="startResize"
            >
              <div class="resize-line"></div>
            </div>
            
            <div 
              v-if="isPreviewOpen" 
              class="md-preview-pane"
              :style="{ width: (100 - editorWidth) + '%' }"
            >
              <div class="preview-content" v-html="previewHtml"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- 详情展示页 -->
      <div v-if="viewMode === 'detail'" class="fade-in">
        <div style="margin-bottom: 20px">
          <el-button @click="goList">
            <el-icon><ArrowLeft /></el-icon>返回
          </el-button>
          <el-button type="primary" style="float: right" @click="goEdit(detailData.id)">编辑文档</el-button>
        </div>
        <div class="detail-view">
          <h1 class="detail-title">{{ detailData.title }}</h1>
          <div class="detail-meta">
            发布于 <span>{{ detailData.createTime }}</span> | 分类：<b>{{ detailData.category }}</b>
          </div>
          <div class="detail-content" v-html="detailHtml" />
        </div>
      </div>
    </div>

    <!-- 分类编辑对话框 -->
    <el-dialog 
      v-model="categoryDialogVisible" 
      :title="categoryForm.id ? '编辑分类' : '新增分类'"
      width="450px"
    >
      <el-form :model="categoryForm" label-width="80px">
        <el-form-item label="分类名称" required>
          <el-input v-model="categoryForm.name" placeholder="请输入分类名称" maxlength="50" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="categoryForm.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="categoryForm.description" type="textarea" :rows="2" placeholder="请输入分类描述" maxlength="200" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveCategory">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, View, ArrowLeft, Folder, Edit, Delete } from '@element-plus/icons-vue'
import { VueMonacoEditor } from '@guolao/vue-monaco-editor'
import { 
  getKnowledgeList, 
  getKnowledgeById, 
  createKnowledge, 
  updateKnowledge, 
  deleteKnowledge,
  getCategoryList,
  createCategory,
  updateCategory,
  deleteCategory
} from '@/api/tool'
import { marked } from 'marked'

const viewMode = ref('list')
const loading = ref(false)
const tableData = ref([])
const isPreviewOpen = ref(true)
const editorWidth = ref(50)
const isResizing = ref(false)
const currentCategory = ref('')
const categoryList = ref([])
const categoryDialogVisible = ref(false)
const categoryForm = reactive({
  id: null,
  name: '',
  sort: 0,
  description: ''
})

const totalDocs = computed(() => {
  return categoryList.value.reduce((sum, cat) => sum + cat.docCount, 0)
})

const monacoOptions = {
  theme: 'vs',
  fontSize: 14,
  fontFamily: "'Fira Code', 'Consolas', monospace",
  automaticLayout: true,
  minimap: { enabled: true },
  wordWrap: 'on',
  lineHeight: 22,
  padding: { top: 20 },
  scrollBeyondLastLine: false,
  renderLineHighlight: 'all',
  cursorBlinking: 'smooth',
  tabSize: 2
}

const searchForm = reactive({
  title: '',
  status: null
})

const pagination = reactive({
  pageNum: 1,
  pageSize: 10,
  total: 0
})

const editForm = reactive({
  id: null,
  title: '',
  category: '',
  content: '',
  tags: '',
  status: 1
})

const detailData = ref({})
const previewHtml = ref('')
const detailHtml = ref('')

const fetchCategoryList = async () => {
  try {
    const res = await getCategoryList()
    if (res.data.code === 200) {
      categoryList.value = res.data.data || []
    }
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const params = {
      pageNum: pagination.pageNum,
      pageSize: pagination.pageSize,
      title: searchForm.title,
      category: currentCategory.value,
      status: searchForm.status
    }
    const res = await getKnowledgeList(params)
    if (res.data.code === 200) {
      tableData.value = res.data.data.list || []
      pagination.total = res.data.data.total || 0
    }
  } catch (error) {
    console.error('获取列表失败:', error)
  } finally {
    loading.value = false
  }
}

const selectCategory = (category) => {
  currentCategory.value = category
  pagination.pageNum = 1
  fetchList()
}

const openCategoryDialog = (category = null) => {
  if (category) {
    categoryForm.id = category.id
    categoryForm.name = category.name
    categoryForm.sort = category.sort
    categoryForm.description = category.description || ''
  } else {
    categoryForm.id = null
    categoryForm.name = ''
    categoryForm.sort = 0
    categoryForm.description = ''
  }
  categoryDialogVisible.value = true
}

const handleSaveCategory = async () => {
  if (!categoryForm.name) {
    ElMessage.warning('请输入分类名称')
    return
  }

  try {
    let res
    if (categoryForm.id) {
      res = await updateCategory(categoryForm)
    } else {
      res = await createCategory(categoryForm)
    }

    if (res.data.code === 200) {
      ElMessage.success(categoryForm.id ? '更新成功' : '创建成功')
      categoryDialogVisible.value = false
      fetchCategoryList()
    } else {
      ElMessage.error(res.data.message || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDeleteCategory = async (id) => {
  try {
    await ElMessageBox.confirm('确定删除该分类吗？该分类下的文档将归类为"未分类"', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await deleteCategory(id)
    if (res.data.code === 200) {
      ElMessage.success('删除成功')
      fetchCategoryList()
      if (currentCategory.value) {
        const deletedCategory = categoryList.value.find(c => c.id === id)
        if (deletedCategory && currentCategory.value === deletedCategory.name) {
          currentCategory.value = ''
          fetchList()
        }
      }
    } else {
      ElMessage.error(res.data.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSearch = () => {
  pagination.pageNum = 1
  fetchList()
}

const handleReset = () => {
  searchForm.title = ''
  searchForm.status = null
  currentCategory.value = ''
  pagination.pageNum = 1
  fetchList()
}

const handleSizeChange = (val) => {
  pagination.pageSize = val
  fetchList()
}

const handleCurrentChange = (val) => {
  pagination.pageNum = val
  fetchList()
}

const goList = () => {
  viewMode.value = 'list'
  fetchCategoryList()
  fetchList()
}

const goEdit = async (id) => {
  if (id) {
    try {
      const res = await getKnowledgeById(id)
      if (res.data.code === 200) {
        editForm.id = res.data.data.id
        editForm.title = res.data.data.title
        editForm.category = res.data.data.category
        editForm.content = res.data.data.content || ''
        editForm.tags = res.data.data.tags
        editForm.status = res.data.data.status
        updatePreview()
      }
    } catch (error) {
      ElMessage.error('获取详情失败')
    }
  } else {
    editForm.id = null
    editForm.title = ''
    editForm.category = currentCategory.value || (categoryList.value[0]?.name || '')
    editForm.content = '# 在此开始编写文档...\n\n'
    editForm.tags = ''
    editForm.status = 1
    previewHtml.value = ''
  }
  isPreviewOpen.value = true
  editorWidth.value = 50
  viewMode.value = 'edit'
}

const togglePreview = () => {
  isPreviewOpen.value = !isPreviewOpen.value
  if (isPreviewOpen.value) {
    updatePreview()
  }
}

const handleEditorChange = (value) => {
  editForm.content = value
  if (isPreviewOpen.value) {
    updatePreview()
  }
}

const updatePreview = () => {
  if (editForm.content) {
    previewHtml.value = marked(editForm.content)
  } else {
    previewHtml.value = ''
  }
}

const startResize = (e) => {
  isResizing.value = true
  document.addEventListener('mousemove', handleResize)
  document.addEventListener('mouseup', stopResize)
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
}

const handleResize = (e) => {
  if (!isResizing.value) return
  
  const container = document.querySelector('.editor-container')
  if (!container) return
  
  const rect = container.getBoundingClientRect()
  const newWidth = ((e.clientX - rect.left) / rect.width) * 100
  
  if (newWidth >= 20 && newWidth <= 80) {
    editorWidth.value = newWidth
  }
}

const stopResize = () => {
  isResizing.value = false
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
}

const handleSave = async (status) => {
  if (!editForm.title) {
    ElMessage.warning('请填写标题')
    return
  }

  const data = {
    title: editForm.title,
    category: editForm.category || '其他',
    content: editForm.content,
    tags: editForm.tags,
    status: status
  }

  try {
    let res
    if (editForm.id) {
      data.id = editForm.id
      res = await updateKnowledge(data)
    } else {
      res = await createKnowledge(data)
    }

    if (res.data.code === 200) {
      ElMessage.success(status === 1 ? '发布成功' : '保存草稿成功')
      goList()
    } else {
      ElMessage.error(res.data.message || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const showDetail = async (id) => {
  try {
    const res = await getKnowledgeById(id)
    if (res.data.code === 200) {
      detailData.value = res.data.data
      detailHtml.value = marked(res.data.data.content || '')
      viewMode.value = 'detail'
    }
  } catch (error) {
    ElMessage.error('获取详情失败')
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定删除该知识吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await deleteKnowledge(id)
    if (res.data.code === 200) {
      ElMessage.success('删除成功')
      fetchCategoryList()
      fetchList()
    } else {
      ElMessage.error(res.data.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  fetchCategoryList()
  fetchList()
})

onUnmounted(() => {
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
})
</script>

<style scoped>
.knowledge-container {
  height: 100%;
}

.card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  padding: 24px;
  min-height: 100%;
}

/* 列表页布局 */
.list-view {
  display: flex;
  gap: 15px;
  height: calc(100vh - 140px);
}

/* 左侧分类面板 */
.category-panel {
  width: 220px;
  background: white;
  border-radius: 8px;
  border: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.panel-header {
  padding: 15px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 14px;
}

.add-icon {
  cursor: pointer;
  color: #1890ff;
  font-size: 16px;
}

.add-icon:hover {
  color: #40a9ff;
}

.category-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px 0;
}

.category-item {
  padding: 10px 15px;
  cursor: pointer;
  font-size: 13px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  transition: 0.2s;
}

.category-item:hover {
  background: #f5f7fa;
}

.category-item.active {
  color: #1890ff;
  background: #e6f7ff;
  font-weight: 600;
}

.category-item .el-icon {
  margin-right: 8px;
}

.doc-count {
  color: #999;
  font-size: 12px;
}

.category-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-icon {
  display: none;
  font-size: 14px;
  color: #666;
  cursor: pointer;
}

.action-icon:hover {
  color: #1890ff;
}

.action-icon.delete:hover {
  color: #ff4d4f;
}

.category-item:hover .action-icon {
  display: inline-flex;
}

.category-item:hover .doc-count {
  display: none;
}

/* 右侧表格区域 */
.table-panel {
  flex: 1;
  background: white;
  border-radius: 8px;
  border: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  padding: 20px;
  overflow: hidden;
}

.search-row {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  align-items: center;
  flex-shrink: 0;
}

.category-tag {
  color: #1890ff;
}

.pagination-wrapper {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
  flex-shrink: 0;
}

/* 编辑页面布局 */
.edit-view {
  height: calc(100vh - 140px);
  display: flex;
  flex-direction: column;
}

.edit-layout {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.edit-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  flex-shrink: 0;
}

.title-input :deep(.el-input__wrapper) {
  font-size: 26px;
  font-weight: 700;
  border: none;
  border-bottom: 2px solid transparent;
  border-radius: 0;
  box-shadow: none;
  padding: 5px 0;
}

.title-input :deep(.el-input__wrapper:focus-within) {
  border-bottom-color: #1890ff;
}

.category-select-row {
  margin-bottom: 15px;
  flex-shrink: 0;
}

.category-select-row .label {
  font-size: 13px;
  color: #666;
}

/* 编辑器容器 */
.editor-container {
  display: flex;
  flex: 1;
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
  min-height: 500px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

/* Monaco编辑器 */
.monaco-editor-pane {
  height: 100%;
  position: relative;
  flex-shrink: 0;
}

/* 可拖拽分隔条 */
.resize-divider {
  width: 6px;
  background: #f0f0f0;
  cursor: col-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background 0.2s;
  z-index: 10;
}

.resize-divider:hover {
  background: #1890ff;
}

.resize-divider:hover .resize-line {
  background: white;
}

.resize-line {
  width: 2px;
  height: 30px;
  background: #ccc;
  border-radius: 1px;
  transition: background 0.2s;
}

/* 预览区 */
.md-preview-pane {
  height: 100%;
  background: white;
  border-left: 1px solid #eee;
  overflow-y: auto;
  padding: 25px;
  flex-shrink: 0;
}

.preview-content {
  line-height: 1.6;
}

.preview-content :deep(h1) {
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
  margin-bottom: 15px;
}

.preview-content :deep(h2) {
  border-left: 4px solid #1890ff;
  padding-left: 15px;
  margin: 20px 0 10px;
}

.preview-content :deep(pre) {
  background: #f6f8fa;
  padding: 16px;
  border-radius: 6px;
  overflow: auto;
  margin: 15px 0;
}

.preview-content :deep(code) {
  font-family: monospace;
  background: rgba(175, 184, 193, 0.2);
  padding: 0.2em 0.4em;
  border-radius: 6px;
}

.preview-content :deep(pre code) {
  background: transparent;
  padding: 0;
}

/* 详情页 */
.detail-view {
  max-width: 850px;
  margin: 0 auto;
  background: white;
  padding: 50px 60px;
  min-height: 100%;
  border-radius: 4px;
}

.detail-title {
  font-size: 32px;
  margin-bottom: 10px;
}

.detail-meta {
  color: #999;
  font-size: 14px;
  margin-bottom: 30px;
  border-bottom: 1px solid #eee;
  padding-bottom: 20px;
}

.detail-content {
  line-height: 1.8;
  font-size: 16px;
}

.detail-content :deep(h1),
.detail-content :deep(h2),
.detail-content :deep(h3) {
  margin: 25px 0 15px;
}

.detail-content :deep(h2) {
  border-left: 4px solid #1890ff;
  padding-left: 15px;
}

.detail-content :deep(pre) {
  background: #282c34;
  color: #abb2bf;
  padding: 15px;
  border-radius: 6px;
  overflow-x: auto;
  margin: 15px 0;
}

.detail-content :deep(code) {
  background: #f0f2f5;
  color: #e74c3c;
  padding: 2px 4px;
  border-radius: 4px;
  font-family: monospace;
}

.detail-content :deep(pre code) {
  background: transparent;
  color: inherit;
  padding: 0;
}

.fade-in {
  animation: fadeIn 0.3s;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>
