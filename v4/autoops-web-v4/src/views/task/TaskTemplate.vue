<template>
  <div class="tasktemplate-management">
    <el-card shadow="hover" class="tasktemplate-card">
      <template #header>
        <div class="card-header">
          <span class="title">任务模板管理</span>
        </div>
      </template>
      
      <!-- 搜索区域 -->
      <div class="search-section">
        <el-form :inline="true" :model="queryParams" class="search-form">
          <el-form-item label="模板名称">
            <el-input 
              placeholder="请输入模板名称" 
              size="small" 
              clearable 
              v-model="queryParams.name"
              @keyup.enter="searchByName(queryParams.name)"
            />
          </el-form-item>
          <el-form-item label="模板类型" style="width: 200px">
            <el-select 
              size="small" 
              placeholder="请选择模板类型" 
              v-model="queryParams.type"
              @change="searchByType(queryParams.type)"
              style="width: 100%"
            >
              <el-option :value="1">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img 
                    src="@/assets/image/shell.svg" 
                    alt="Shell"
                    style="width: 14px; height: 14px; object-fit: contain; flex-shrink: 0;"
                  />
                  <span>Shell</span>
                </div>
              </el-option>
              <el-option :value="2">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img 
                    src="@/assets/image/Python.svg" 
                    alt="Python"
                    style="width: 14px; height: 14px; object-fit: contain; flex-shrink: 0;"
                  />
                  <span>Python</span>
                </div>
              </el-option>
              <el-option :value="3">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img 
                    src="@/assets/image/ansible.svg" 
                    alt="Ansible"
                    style="width: 14px; height: 14px; object-fit: contain; flex-shrink: 0;"
                  />
                  <span>Ansible</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="small" @click="searchByName(queryParams.name)" class="modern-btn primary-btn">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button size="small" @click="resetQuery" class="modern-btn reset-btn">
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
          @click="handleCreate"
          v-authority="['task:template:add']"
          class="modern-btn success-btn"
        >
          <el-icon><Plus /></el-icon>
          新增模板
        </el-button>
      </div>

      <!-- 列表区域 -->
      <div class="table-section">
        <el-table 
          v-loading="loading" 
          :data="templates" 
          border 
          stripe 
          style="width: 100%"
          :header-cell-style="{ background: '#eef1f6', color: '#606266' }"
        >
          <el-table-column prop="name" label="模板名称" width="270">
            <template #default="{row}">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img 
                  src="@/assets/image/moban.svg" 
                  alt="模板"
                  style="width: 20px; height: 20px; object-fit: contain; flex-shrink: 0;"
                />
                <span>{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="type" label="类型" width="140">
            <template #default="{row}">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img 
                  :src="getTypeIcon(row.type)" 
                  :alt="getTypeName(row.type)"
                  style="width: 18px; height: 18px; object-fit: contain; flex-shrink: 0;"
                />
                <span :class="getTypeClass(row.type)">
                  {{ getTypeName(row.type) }}
                </span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="脚本内容" width="100">
            <template #default="{row}">
              <el-tooltip content="查看脚本" placement="top">
                <el-button
                  type="primary"
                  :icon="View"
                  size="small"
                  @click="handleViewContent(row.id)"
                >脚本</el-button>
              </el-tooltip>
            </template>
          </el-table-column>
          <el-table-column prop="createdBy" label="创建人" width="140">
            <template #default="{row}">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img 
                  src="@/assets/image/ren.svg" 
                  alt="创建人"
                  style="width: 18px; height: 18px; object-fit: contain; flex-shrink: 0;"
                />
                <span>{{ row.createdBy }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="180" />
          <el-table-column prop="remark" label="备注信息" >
            <template #default="{row}">
              {{ row.remark }}
            </template>
          </el-table-column>

          <el-table-column label="操作" width="180">
            <template #default="{row}">
              <div class="operation-buttons">
                <el-tooltip content="编辑" placement="top">
                  <el-button
                    type="warning"
                    :icon="Edit"
                    size="small"
                    v-authority="['task:template:edit']"
                    @click="handleEdit(row)"
                  >编辑</el-button>
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button
                    type="danger"
                    :icon="Delete"
                    size="small"
                    v-authority="['task:template:delete']"
                    @click="handleDelete(row.id)"
                  >删除</el-button>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页区域 -->
      <div class="pagination-section">
        <el-pagination 
          @size-change="handleSizeChange" 
          @current-change="handleCurrentChange"
          :current-page="queryParams.pageNum" 
          :page-sizes="[10, 50, 100]" 
          :page-size="queryParams.pageSize"
          layout="total, sizes, prev, pager, next, jumper" 
          :total="total"
        />
      </div>
    </el-card>

    <!--模板表单对话框-->
    <el-dialog 
      :title="currentTemplate.id ? '编辑模板' : '新增模板'" 
      v-model="formVisible" 
      width="40%"
      @close="formDialogClosed"
    >
      <el-form 
        ref="templateFormRef" 
        label-width="100px" 
        :model="currentTemplate" 
        :rules="templateRules"
      >
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="currentTemplate.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="模板类型" prop="type">
          <el-select v-model="currentTemplate.type" placeholder="请选择模板类型">
            <el-option :value="1">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img 
                  src="@/assets/image/shell.svg" 
                  alt="Shell"
                  style="width: 16px; height: 16px; object-fit: contain; flex-shrink: 0;"
                />
                <span>Shell</span>
              </div>
            </el-option>
            <el-option :value="2">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img 
                  src="@/assets/image/Python.svg" 
                  alt="Python"
                  style="width: 16px; height: 16px; object-fit: contain; flex-shrink: 0;"
                />
                <span>Python</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="模板内容" prop="content">
          <div class="code-editor-wrapper">
            <CodeEditor
              v-model="currentTemplate.content"
              :language="getLanguage(currentTemplate.type)"
              placeholder="请输入模板内容"
            />
          </div>
        </el-form-item>
        <el-form-item label="备注信息">
          <el-input 
            v-model="currentTemplate.remark" 
            type="textarea" 
            :rows="2"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!--脚本内容查看对话框-->
    <el-dialog 
      title="脚本内容" 
      v-model="contentDialogVisible" 
      width="60%"
    >
      <pre 
        style="white-space: pre-wrap; background: #1e1e1e; padding: 10px; border-radius: 4px;"
        v-html="scriptContent"
      ></pre>
      <template #footer>
        <el-button @click="contentDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import CodeEditor from '@/components/CodeEditor.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Edit, Delete, Search, Refresh, Plus } from '@element-plus/icons-vue'
import { highlight } from '@/utils/highlight'
import { 
  CreateTemplate, 
  GetAllTemplates,
  UpdateTemplate,
  DeleteTemplate,
  GetTemplateByID,
  GetTemplatesByName,
  GetTemplatesByType,
  GetTemplateContent
} from '@/api/task'

const templates = ref([])
const loading = ref(false)
const formVisible = ref(false)
const templateFormRef = ref(null)
const total = ref(0)

const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  name: '',
  type: null
})

const currentTemplate = ref({
  name: '',
  type: 1,
  content: '',
  remark: ''
})

const contentDialogVisible = ref(false)
const scriptContent = ref('')

const templateRules = reactive({
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择模板类型', trigger: 'change' }],
  content: [{ required: true, message: '请输入模板内容', trigger: 'blur' }]
})

const fetchTemplates = async () => {
  loading.value = true
  try {
    const response = await GetAllTemplates(queryParams)
    const responseData = response?.data || {}
    templates.value = Array.isArray(responseData.data) 
      ? responseData.data.map(item => ({
          id: item.id,
          name: item.name,
          type: item.type,
          content: item.content,
          remark: item.remark || '', // 确保remark字段存在
          createdBy: item.createdBy,
          updatedBy: item.updatedBy,
          createdAt: item.createdAt,
          updatedAt: item.updatedAt
        })) 
      : []
    console.log('模板数据:', JSON.stringify(templates.value, null, 2)) // 详细日志
    total.value = templates.value.length
  } catch (error) {
    console.error('获取模板列表异常:', error)
    templates.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  currentTemplate.value = {
    name: '',
    type: 1,
    content: ''
  }
  formVisible.value = true
}

const handleEdit = (template) => {
  currentTemplate.value = { ...template }
  formVisible.value = true
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除该模板吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await DeleteTemplate({ id })
    ElMessage.success('删除成功')
    fetchTemplates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleViewContent = async (id) => {
  try {
    const response = await GetTemplateContent({ id })
    const template = templates.value.find(t => t.id === id)
    const language = template?.type === 1 ? 'bash' : 
                    template?.type === 2 ? 'python' : 'yaml'
    scriptContent.value = highlight(response.data, language)
    contentDialogVisible.value = true
  } catch (error) {
    console.error('获取脚本内容失败:', error)
    scriptContent.value = '获取脚本内容失败'
  }
}

const handleSubmit = async () => {
  await templateFormRef.value.validate()
  if (currentTemplate.value.id) {
    const updateData = {
      id: currentTemplate.value.id,
      name: currentTemplate.value.name,
      type: currentTemplate.value.type,
      content: currentTemplate.value.content,
      remark: currentTemplate.value.remark || ''
    }
    await UpdateTemplate(updateData)
  } else {
    await CreateTemplate(currentTemplate.value)
  }
  formVisible.value = false
  fetchTemplates()
}

const searchByName = async (name) => {
  try {
    loading.value = true
    queryParams.name = name
    queryParams.pageNum = 1
    const params = {
      name: name || undefined,
      pageNum: queryParams.pageNum,
      pageSize: queryParams.pageSize
    }
    const res = await GetTemplatesByName(params)
    console.log('按名称搜索结果:', res)
    
    // 处理不同格式的响应
    if (res.data && Array.isArray(res.data)) {
      templates.value = res.data
      total.value = res.data.length
    } else if (res.data && res.data.data) {
      templates.value = res.data.data
      total.value = res.data.total || res.data.data.length
    } else {
      templates.value = []
      total.value = 0
      ElMessage.warning('未获取到有效数据')
    }
  } catch (error) {
    console.error('按名称搜索失败:', error)
    ElMessage.error(`按名称搜索失败: ${error.message || '未知错误'}`)
  } finally {
    loading.value = false
  }
}

const searchByType = async (type) => {
  try {
    loading.value = true
    queryParams.type = type
    queryParams.pageNum = 1
    const params = {
      type: type || undefined,
      pageNum: queryParams.pageNum,
      pageSize: queryParams.pageSize
    }
    const res = await GetTemplatesByType(params)
    console.log('按类型搜索结果:', res)
    
    // 处理不同格式的响应
    if (res.data && Array.isArray(res.data)) {
      templates.value = res.data
      total.value = res.data.length
    } else if (res.data && res.data.data) {
      templates.value = res.data.data
      total.value = res.data.total || res.data.data.length
    } else {
      templates.value = []
      total.value = 0
      ElMessage.warning('未获取到有效数据')
    }
  } catch (error) {
    console.error('按类型搜索失败:', error)
    ElMessage.error(`按类型搜索失败: ${error.message || '未知错误'}`)
  } finally {
    loading.value = false
  }
}

const resetQuery = () => {
  queryParams.name = ''
  queryParams.type = null
  queryParams.pageNum = 1
  fetchTemplates()
}

const formDialogClosed = () => {
  templateFormRef.value?.resetFields()
}

const handleSizeChange = (val) => {
  queryParams.pageSize = val
  fetchTemplates()
}

const handleCurrentChange = (val) => {
  queryParams.pageNum = val
  fetchTemplates()
}

const getTypeClass = (type) => {
  switch(type) {
    case 1: return 'type-shell'
    case 2: return 'type-python'
    case 3: return 'type-ansible'
    default: return ''
  }
}

const getLanguage = (type) => {
  switch(type) {
    case 1: return 'bash'
    case 2: return 'python'
    case 3: return 'yaml'
    default: return null
  }
}

const getTypeIcon = (type) => {
  switch(type) {
    case 1: return require('@/assets/image/shell.svg')
    case 2: return require('@/assets/image/Python.svg')
    case 3: return require('@/assets/image/ansible.svg')
    default: return require('@/assets/image/shell.svg')
  }
}

const getTypeName = (type) => {
  switch(type) {
    case 1: return 'Shell'
    case 2: return 'Python'
    case 3: return 'Ansible'
    default: return 'Shell'
  }
}

onMounted(() => {
  fetchTemplates()
})
</script>

<style scoped>
.tasktemplate-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.tasktemplate-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.search-section {
  margin-bottom: 16px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

.action-section {
  margin-bottom: 16px;
}

.table-section {
  margin-top: 0;
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

:deep(.el-table) {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

:deep(.el-table__header) {
  background: rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-table__header th) {
  background: rgba(102, 126, 234, 0.1) !important;
  color: #2c3e50 !important;
  font-weight: 600;
  border: none;
}

:deep(.el-table__body tr:hover > td) {
  background-color: rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-table td) {
  border: none;
}

:deep(.el-table::before) {
  display: none;
}

:deep(.el-table--border::after) {
  display: none;
}

.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.operation-buttons .el-button {
  transition: all 0.3s ease;
}

.operation-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.type-shell {
  color: #67C23A;
  font-weight: 500;
}

.type-python {
  color: #409EFF;
  font-weight: 500;
}

.type-ansible {
  color: #E6A23C;
  font-weight: 500;
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

:deep(.reset-btn) {
  background: linear-gradient(45deg, #E6A23C, #EEBE77);
  color: white;
}

:deep(.success-btn) {
  background: linear-gradient(45deg, #67C23A, #85CE61);
  color: white;
}

.code-editor-wrapper {
  width: 100%;
  height: 300px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}
</style>
