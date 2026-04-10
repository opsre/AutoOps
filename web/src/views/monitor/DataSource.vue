<template>
  <div class="datasource-management">
    <el-card class="modern-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="header-title">
            <el-icon><Monitor /></el-icon> 数据源管理
          </span>
          <div class="right-actions">
            <el-button type="primary" icon="Plus" @click="openDialog()">新增数据源</el-button>
          </div>
        </div>
      </template>

      <!-- 搜索 -->
      <div class="search-section">
        <el-form :inline="true" :model="query" class="search-form">
          <el-form-item label="名称">
            <el-input v-model="query.name" placeholder="数据源名称" clearable @keyup.enter="fetchData" style="width: 180px" />
          </el-form-item>
          <el-form-item label="类型">
            <el-select v-model="query.type" clearable placeholder="选择类型" @change="fetchData" style="width: 150px">
              <el-option label="Prometheus" value="Prometheus" />
              <el-option label="Zabbix" value="Zabbix" />
              <el-option label="Loki" value="Loki" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="Search" @click="fetchData">查询</el-button>
            <el-button type="warning" icon="Refresh" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格 -->
      <el-table
        :data="tableData"
        v-loading="loading"
        border
        stripe
        :header-cell-style="{ background: 'rgba(102, 126, 234, 0.1)', color: '#2c3e50', fontWeight: '600' }"
      >
        <el-table-column prop="name" label="数据源名称" min-width="150" />
        <el-table-column prop="type" label="类型" width="140" align="left">
          <template #default="{ row }">
             <div style="display: flex; align-items: center; gap: 8px;">
               <img :src="getDsTypeIcon(row.type)" class="ds-icon" v-if="getDsTypeIcon(row.type)" />
               <span>{{ row.type || '未知' }}</span>
             </div>
          </template>
        </el-table-column>
        <el-table-column prop="apiUrl" label="连接地址 (URL)" min-width="200" show-overflow-tooltip />
        <el-table-column prop="deployMethod" label="部署方式" width="140" align="center">
          <template #default="{ row }">
             <el-tag size="small" :type="row.deployMethod === 'Kubernetes' ? 'primary' : 'info'" style="border-radius: 12px;">
               {{ row.deployMethod || '未知' }}
             </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="config" label="配置" min-width="150" show-overflow-tooltip />
         <el-table-column label="创建时间" width="160" align="center">
          <template #default="{ row }">
            {{ formatDate(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-button link type="primary" icon="EditPen" @click="openDialog(row)">编辑</el-button>
              <el-button link type="danger" icon="Delete" @click="handleDelete(row.id || row.ID)">删除</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-section">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          v-model:current-page="query.page"
          v-model:page-size="query.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          @current-change="fetchData"
          @size-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 弹窗 -->
    <el-dialog :title="form.id ? '编辑数据源' : '新增数据源'" v-model="dialogVisible" width="600px" class="modern-dialog">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" class="rule-form">
        <div class="form-section-title"><el-icon><Setting /></el-icon> 基础配置</div>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="名称" prop="name">
              <el-input v-model="form.name" placeholder="请输入名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="类型" prop="type">
              <el-select v-model="form.type" placeholder="选择类型" style="width: 100%">
                <el-option label="Prometheus" value="Prometheus" />
                <el-option label="Zabbix" value="Zabbix" />
                <el-option label="Loki" value="Loki" />
                <el-option label="ElasticSearch" value="ElasticSearch" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="连接地址" prop="apiUrl">
          <el-input v-model="form.url" placeholder="http://ip:port (例如: http://prometheus:9090)" />
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="部署方式" prop="deployMethod">
              <el-select v-model="form.deployMethod" placeholder="请选择" style="width: 100%">
                <el-option label="Kubernetes" value="Kubernetes" />
                <el-option label="Docker" value="Docker" />
                <el-option label="Binary" value="Binary" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <!-- 预留其他配置位置 -->
          </el-col>
        </el-row>

        <el-form-item label="配置" prop="config">
          <el-input v-model="form.config" type="textarea" :rows="3" placeholder="写点描述..." />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Monitor, Setting } from '@element-plus/icons-vue'
import { getDataSources, createDataSource, updateDataSource, deleteDataSource } from '@/api/monitor'

const query = reactive({
  name: '',
  type: '',
  page: 1,
  pageSize: 10
})

const loading = ref(false)
const tableData = ref([])
const total = ref(0)

const formRef = ref(null)
const dialogVisible = ref(false)
const submitting = ref(false)

const form = reactive({
  id: null,
  name: '',
  type: 'Prometheus',
  url: '',
  deployMethod: 'Kubernetes',
  config: {},
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入数据源名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择数据源类型', trigger: 'change' }],
  url: [{ required: true, message: '请输入连接地址', trigger: 'blur' }]
}

onMounted(() => {
  fetchData()
})

const resetQuery = () => {
  query.name = ''
  query.type = ''
  query.page = 1
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getDataSources(query)
    const data = res.data?.data || res.data || {}
    tableData.value = data.list || (Array.isArray(data) ? data : [])
    total.value = data.total || tableData.value.length || 0
  } catch (error) {
    console.error('获取数据源列表失败', error)
  } finally {
    loading.value = false
  }
}

const openDialog = (row = null) => {
  if (formRef.value) formRef.value.resetFields()
  if (row) {
    form.id = row.id || row.ID
    form.name = row.name
    form.type = row.type
    form.url = row.apiUrl
    form.deployMethod = row.deployMethod
    form.config = row.config
  } else {
    form.id = null
    form.name = ''
    form.type = 'Prometheus'
    form.url = ''
    form.deployMethod = 'Kubernetes'
    form.config = {}
  }
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const action = form.id ? updateDataSource : createDataSource
        const payload = { ...form }
        if (payload.id) {
          payload.id = parseInt(payload.id)
        } else {
          delete payload.id
        }
        
        await action(payload)
        ElMessage.success(form.id ? '编辑数据源成功' : '新增数据源成功')
        dialogVisible.value = false
        fetchData()
      } catch (error) {
        ElMessage.error('保存数据源失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleDelete = (id) => {
  ElMessageBox.confirm('确定删除该数据源吗？相关依赖可能会受到影响。', '警告', { type: 'warning' }).then(async () => {
    try {
      await deleteDataSource(id)
      ElMessage.success('删除成功')
      if (tableData.value.length === 1 && query.page > 1) {
        query.page -= 1
      }
      fetchData()
    } catch(err) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

// GUI Tools
const getDsTypeIcon = (type) => {
  const t = (type || '').toLowerCase()
  if (t.includes('prometheus')) return 'https://upload.wikimedia.org/wikipedia/commons/3/38/Prometheus_software_logo.svg'
  if (t.includes('zabbix')) return 'https://upload.wikimedia.org/wikipedia/commons/3/33/Zabbix_logo.svg'
  if (t.includes('loki')) return 'https://upload.wikimedia.org/wikipedia/commons/6/60/Grafana_logo.svg'
  if (t.includes('elasticsearch')) return 'https://upload.wikimedia.org/wikipedia/commons/f/f4/Elasticsearch_logo.svg'
  return 'https://cdn-icons-png.flaticon.com/512/3168/3168610.png'
}

const formatDate = (dateStr) => {
  if (!dateStr || String(dateStr).startsWith('0001')) return '-'
  return new Date(dateStr).toLocaleString()
}
</script>

<style scoped>
.datasource-management {
  padding: 20px;
  min-height: calc(100vh - 120px);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.modern-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  min-height: 600px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 搜索表单部分 */
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

/* 表格美化 */
:deep(.el-table) {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
}

:deep(.el-table__header th) { border: none; }
:deep(.el-table__body tr:hover > td) { background-color: rgba(102, 126, 234, 0.1) !important; }
:deep(.el-table td) { border: none; }
:deep(.el-table::before) { display: none; }
:deep(.el-table--border::after) { display: none; }

.ds-icon {
  width: 20px;
  height: 20px;
  object-fit: contain;
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.1));
}

.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

/* 对话框美化 */
:deep(.modern-dialog .el-dialog) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(31, 38, 135, 0.37);
}

:deep(.modern-dialog .el-dialog__header) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
  border-radius: 16px 16px 0 0;
  padding: 20px 24px 16px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.2);
}

:deep(.modern-dialog .el-dialog__title) { color: #2c3e50; font-weight: 600; font-size: 18px; }
:deep(.modern-dialog .el-dialog__body) { padding: 24px; }
:deep(.modern-dialog .el-dialog__footer) {
  padding: 16px 24px 24px;
  background: rgba(248, 249, 250, 0.8);
  border-radius: 0 0 16px 16px;
}

.form-section-title {
  font-size: 15px;
  font-weight: 600;
  color: #5a6fd8;
  margin: 0 0 20px 0;
  display: flex;
  align-items: center;
  gap: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.2);
}
</style>
