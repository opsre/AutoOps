<template>
  <div class="taskjob-management">
    <el-card shadow="hover" class="taskjob-card">
      <template #header>
        <div class="card-header">
          <span class="title">任务管理</span>
        </div>
      </template>
      
      <!-- 搜索区域 -->
      <div class="search-section">
        <el-form :inline="true" :model="queryParams" class="search-form">
          <el-form-item label="任务名称">
            <el-input
              placeholder="请输入任务名称"
              size="small"
              clearable
              v-model="queryParams.name"
              @keyup.enter="searchByName(queryParams.name)"
            />
          </el-form-item>
          <el-form-item label="任务类型" style="width: 200px">
            <el-select
              size="small"
              placeholder="请选择任务类型"
              v-model="queryParams.type"
              @change="searchByType(queryParams.type)"
              style="width: 100%"
            >
              <el-option label="普通任务" :value="1" />
              <el-option label="定时任务" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item label="任务状态" style="width: 200px">
            <el-select
              size="small"
              placeholder="请选择任务状态"
              v-model="queryParams.status"
              @change="searchByStatus(queryParams.status)"
              style="width: 100%"
            >
              <el-option label="等待中" :value="1" />
              <el-option label="运行中" :value="2" />
              <el-option label="成功" :value="3" />
              <el-option label="异常" :value="4" />
              <el-option label="已暂停" :value="5" />
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
          v-authority="['task:job:add']"
          class="modern-btn success-btn"
        >
          <el-icon><Plus /></el-icon>
          新增任务
        </el-button>
      </div>

      <!-- 列表区域 -->
      <div class="table-section">
        <el-table
          v-loading="loading"
          :data="tasks"
          border
          stripe
          style="width: 100%"
          :header-cell-style="{ background: '#eef1f6', color: '#606266' }"
        >
          <el-table-column prop="name" label="任务名称" >
            <template #default="{row}">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img 
                  src="@/assets/image/renwu.svg" 
                  alt="任务"
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
          <el-table-column prop="status" label="状态"  >
            <template #default="{row}">
              <span :class="getStatusClass(row.status)">
                {{ getStatusText(row.status) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="cronExpr" label="定时表达式" >
            <template #default="{row}">
              <span>{{ row.type === 2 ? (row.cronExpr || '-') : '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="duration" label="执行耗时(秒)" >
            <template #default="{row}">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img
                  src="@/assets/image/shijian.svg"
                  alt="时间"
                  style="width: 18px; height: 18px; object-fit: contain; flex-shrink: 0;"
                />
                <span>{{ row.duration || '0' }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="executeCount" label="执行次数" >
            <template #default="{row}">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img
                  src="@/assets/image/统计管理.svg"
                  alt="次数"
                  style="width: 18px; height: 18px; object-fit: contain; flex-shrink: 0;"
                />
                <span class="execute-count">{{ row.executeCount || '0' }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="任务开关"  align="center">
            <template #default="{row}">
              <el-switch
                v-if="row.type === 2"
                v-model="row.isActive"
                :disabled="row.status === 1"
                @change="handleToggleTask(row)"
                active-text="启动"
                inactive-text="暂停"
              />
              <span v-else style="color: #C0C4CC;">-</span>
            </template>
          </el-table-column>
          <el-table-column prop="nextRunTime" label="下次执行时间" width="220" >
            <template #default="{row}">
              <div style="display: flex; align-items: center; gap: 8px;" v-if="row.type === 2 && row.nextRunTime">
                <img
                  src="@/assets/image/dingshirenwu.svg"
                  alt="定时"
                  style="width: 18px; height: 18px; object-fit: contain; flex-shrink: 0;"
                />
                <span class="next-run-time">{{ formatTime(row.nextRunTime) }}</span>
              </div>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注信息" />

          <el-table-column label="操作" >
            <template #default="{row}">
              <div class="operation-buttons">
                <el-tooltip content="启动任务" placement="top">
                  <el-button
                    type="success"
                    :icon="VideoPlay"
                    size="small"
                    v-authority="['task:job:start']"
                    @click="async () => {
                      console.log('Starting task with ID:', row.id);
                      if (!row.id) {
                        ElMessage.error('任务ID无效');
                        return;
                      }
                      const apiTaskId = await jobFlowRef?.showFlow?.(row.id);
                      console.log('API Task ID:', apiTaskId);
                    }"
                  >启动</el-button>
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button
                    type="danger"
                    :icon="Delete"
                    size="small"
                    v-authority="['task:job:delete']"
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

    <!--任务表单对话框-->
    <el-dialog
      title="新增任务"
      v-model="formVisible"
      width="35%"
      @close="formDialogClosed"
    >
      <el-form
        ref="taskFormRef"
        label-width="100px"
        :model="currentTask"
        :rules="taskRules"
      >
        <el-form-item label="任务名称" prop="name" label-width="120px">
          <el-input v-model="currentTask.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="任务类型" prop="type" label-width="120px">
          <el-select v-model="currentTask.type" placeholder="请选择任务类型">
            <el-option label="普通任务" :value="1" />
            <el-option label="定时任务" :value="2" />
          </el-select>
        </el-form-item>
                <el-form-item label="定时表达式" prop="cronExpr" v-if="currentTask.type === 2">
          <el-input
            v-model="currentTask.cronExpr"
            placeholder="请输入cron表达式，如: * * * * *"
            @change="calculateNextExecution"
          />
          <span style="margin-left: 10px; color: #67C23A; font-weight: 500;" v-if="nextExecutionTime">
            下次执行时间: {{ nextExecutionTime }}
          </span>
        </el-form-item>
        <el-form-item label="任务模版" prop="shell" style="margin-bottom: 15px">
          <div style="display: flex; align-items: center; margin-bottom: 10px">
            <el-button
              type="primary"
              icon="Plus"
              size="small"
              @click="handleAddTemplate"
              style="margin-left: 10px"
            >添加模板</el-button>
          </div>
          <div v-if="selectedTemplates.length > 0" style="border: 1px solid #ebeef5; border-radius: 4px; padding: 10px; background: #f5f7fa">
            <div style="color: #606266; margin-bottom: 8px">已选择模板:</div>
            <div style="display: flex; flex-direction: column; gap: 8px">
              <div
                v-for="templateId in selectedTemplates"
                :key="templateId"
                style="display: flex; justify-content: space-between; align-items: center; padding: 6px 10px; background: white; border-radius: 4px"
              >
                <span>{{ getTemplateName(templateId) }}</span>
                <el-button
                  type="danger"
                  icon="Close"
                  size="small"
                  plain
                  @click="removeTemplate(templateId)"
                >移除</el-button>
              </div>
            </div>
          </div>
        </el-form-item>

        <el-dialog
          title="选择任务模板"
          v-model="showTemplateDialog"
          width="35%"
        >
          <el-select
            v-model="selectedTemplates"
            multiple
            placeholder="请选择任务模板"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="template in templates"
              :key="template.id"
              :label="template.name"
              :value="template.id"
            />
          </el-select>
          <template #footer>
            <el-button @click="showTemplateDialog = false">取消</el-button>
            <el-button
              type="primary"
              @click="handleTemplateSelected(selectedTemplates)"
            >确定</el-button>
          </template>
        </el-dialog>

        <el-form-item label="执行主机" prop="hostIDs">
          <div style="display: flex; align-items: center; margin-bottom: 10px">
        <el-button
          type="primary"
          icon="Plus"
          size="small"
          @click="handleAddHost"
          style="margin-left: 10px"
        >添加主机</el-button>
          </div>
          <div v-if="selectedHosts.length > 0" style="border: 1px solid #ebeef5; border-radius: 4px; padding: 10px; background: #f5f7fa">
            <div style="color: #606266; margin-bottom: 8px">已选择主机:</div>
            <div style="display: flex; flex-direction: column; gap: 8px">
              <div
                v-for="host in selectedHosts"
                :key="host.id"
                style="display: flex; justify-content: space-between; align-items: center; padding: 6px 10px; background: white; border-radius: 4px"
              >
                <span>{{ host.name }} ({{ host.ip }})</span>
                <el-button
                  type="danger"
                  icon="Close"
                  size="small"
                  plain
                  @click="removeHost(host.id)"
                >移除</el-button>
              </div>
            </div>
          </div>
        </el-form-item>

        <HostSelectorDialog
          v-model="showHostDialog"
          title="选择主机"
          @hosts-selected="handleHostSelected"
        />
        <el-form-item label="备注信息">
          <el-input
            v-model="currentTask.remark"
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

    <!--JobFlow组件-->
    <JobFlow ref="jobFlowRef" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, shallowRef, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, VideoPlay, Delete } from '@element-plus/icons-vue'
import {
  CreateTask,
  UpdateTask,
  DeleteTask,
  ListTasks,
  GetTasksByName,
  GetTasksByType,
  GetTasksByStatus,
  GetNextExecutionTime,
  GetAllTemplates,
  PauseScheduledTask,
  ResumeScheduledTask
} from '@/api/task'
import JobFlow from './Job/JobFlow.vue'
import cmdbAPI from '@/api/cmdb'
import HostSelectorDialog from '@/components/HostSelectorDialog.vue'

const tasks = ref([])
const loading = ref(false)
const formVisible = ref(false)
const taskFormRef = shallowRef(null)
const total = ref(0)
const jobFlowRef = shallowRef(null)
const nextExecutionTime = ref('')
const templates = ref([])
const groupsWithHosts = ref([])
const selectedTemplates = ref([])
const selectedHosts = ref([])
const showTemplateDialog = ref(false)
const showHostDialog = ref(false)

const queryParams = reactive({
  page: 1,
  pageSize: 10,
  name: '',
  type: null,
  status: null
})

const currentTask = ref({
  name: '',
  type: 1,
  shell: '',
  hostIDs: '',
  cronExpr: '',
  remark: ''
})

const taskRules = reactive({
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  shell: [{
    required: true,
    validator: (rule, value, callback) => {
      if (selectedTemplates.value.length > 0) {
        callback()
      } else {
        callback(new Error('请选择任务模板'))
      }
    },
    trigger: 'change'
  }],
  cronExpr: [{
    required: false,
    validator: (rule, value, callback) => {
      if (currentTask.value.type === 2 && !value) {
        callback(new Error('请输入定时表达式'))
      } else {
        callback()
      }
    },
    trigger: 'blur'
  }],
  hostIDs: [{ required: true, message: '请选择执行主机', trigger: 'blur' }]
})

const fetchTasks = async () => {
  loading.value = true
  try {
    console.log('请求参数:', JSON.stringify(queryParams, null, 2))
    const response = await ListTasks({
      page: queryParams.page,
      pageSize: queryParams.pageSize,
      name: queryParams.name || undefined,
      type: queryParams.type || undefined,
      status: queryParams.status || undefined
    })
    console.log('API响应:', JSON.stringify(response, null, 2))
    const responseData = response?.data || {}
    if (responseData.code !== 200) {
      throw new Error(responseData.message || '获取任务列表失败')
    }
    tasks.value = Array.isArray(responseData.data?.list)
      ? responseData.data.list.map(item => ({
          id: item.id,
          name: item.name,
          type: item.type,
          shell: item.shell,
          hostIDs: item.host_ids,
          cronExpr: item.cron_expr,
          status: item.status,
          duration: item.duration,
          executeCount: item.execute_count || 0,
          nextRunTime: item.next_run_time,
          remark: item.remark,
          startTime: item.start_time ? formatTime(item.start_time) : '',
          endTime: item.end_time ? formatTime(item.end_time) : '',
          createdAt: item.created_at ? formatTime(item.created_at) : '',
          isActive: item.status === 2 // status=2运行中为true(右边启动)，status=5已暂停为false(左边暂停)
        }))
      : []
    total.value = responseData.data?.total || tasks.value.length
    console.log('处理后的任务列表数据:', tasks.value)
  } catch (error) {
    console.error('获取任务列表异常:', error)
    tasks.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  currentTask.value = {
    name: '',
    type: 1,
    shell: '',
    hostIDs: '',
    cronExpr: '',
    remark: ''
  }
  selectedTemplates.value = []
  selectedHosts.value = []
  nextTick(() => {
    formVisible.value = true
  })
}


const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除该任务吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await DeleteTask({ id })
    ElMessage.success('删除成功')
    fetchTasks()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除任务失败:', error)
      ElMessage.error(`删除失败: ${error.response?.data?.message || error.message || '未知错误'}`)
    }
  }
}


const handleSubmit = async () => {
  try {
    await taskFormRef.value?.validate()

    // 确保shell内容不为空
    if (!currentTask.value.shell && selectedTemplates.value.length > 0) {
      currentTask.value.shell = selectedTemplates.value.join(',')
    }

    const requestData = {
      name: currentTask.value.name.trim(),
      type: currentTask.value.type,
      shell: currentTask.value.shell || '', // 确保shell字段存在
      host_ids: currentTask.value.hostIDs,
      remark: currentTask.value.remark
    }

    // 如果是定时任务，添加cron表达式
    if (currentTask.value.type === 2) {
      requestData.cron_expr = currentTask.value.cronExpr
    }

    console.log('提交的任务数据:', JSON.stringify(requestData, null, 2))
    
    const response = await CreateTask(requestData)
    console.log('API响应:', JSON.stringify(response, null, 2))
    
    if (response.data?.code === 200) {
      ElMessage.success('任务创建成功')
      formVisible.value = false
      await nextTick()
      fetchTasks()
    } else {
      throw new Error(response.data?.message || '创建任务失败')
    }
  } catch (error) {
    if (error.name !== 'ValidationError') {
      console.error('创建任务失败:', error)
      ElMessage.error(`创建任务失败: ${error.response?.data?.message || error.message || '未知错误'}`)
    }
  }
}

const searchByName = async (name) => {
  try {
    loading.value = true
    queryParams.name = name
    queryParams.pageNum = 1
    const response = await GetTasksByName({
      name: name || undefined,
      pageNum: queryParams.pageNum,
      pageSize: queryParams.pageSize
    })
    tasks.value = response.data?.data || []
    total.value = response.data?.total || tasks.value.length
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
    const response = await GetTasksByType({
      type: type || undefined,
      pageNum: queryParams.pageNum,
      pageSize: queryParams.pageSize
    })
    tasks.value = response.data?.data || []
    total.value = response.data?.total || tasks.value.length
  } catch (error) {
    console.error('按类型搜索失败:', error)
    ElMessage.error(`按类型搜索失败: ${error.message || '未知错误'}`)
  } finally {
    loading.value = false
  }
}

const searchByStatus = async (status) => {
  try {
    loading.value = true
    queryParams.status = status
    queryParams.pageNum = 1
    const response = await GetTasksByStatus({
      status: status || undefined,
      pageNum: queryParams.pageNum,
      pageSize: queryParams.pageSize
    })
    tasks.value = response.data?.data || []
    total.value = response.data?.total || tasks.value.length
  } catch (error) {
    console.error('按状态搜索失败:', error)
    ElMessage.error(`按状态搜索失败: ${error.message || '未知错误'}`)
  } finally {
    loading.value = false
  }
}

const resetQuery = () => {
  queryParams.name = ''
  queryParams.type = null
  queryParams.status = null
  queryParams.pageNum = 1
  fetchTasks()
}

const formDialogClosed = () => {
  taskFormRef.value?.resetFields()
}

const handleSizeChange = (val) => {
  queryParams.pageSize = val
  queryParams.page = 1
  fetchTasks()
}

const handleCurrentChange = (val) => {
  queryParams.page = val
  fetchTasks()
}

const getTypeClass = (type) => {
  switch(type) {
    case 1: return 'type-normal'
    case 2: return 'type-cron'
    case 3: return 'type-ansible'
    case 4: return 'type-job'
    default: return ''
  }
}

const getStatusClass = (status) => {
  switch(status) {
    case 1: return 'status-waiting'
    case 2: return 'status-running'
    case 3: return 'status-success'
    case 4: return 'status-error'
    case 5: return 'status-paused'
    default: return ''
  }
}

const getStatusText = (status) => {
  switch(status) {
    case 1: return '等待中'
    case 2: return '运行中'
    case 3: return '成功'
    case 4: return '异常'
    case 5: return '已暂停'
    default: return '未知'
  }
}

const getTypeIcon = (type) => {
  switch(type) {
    case 1: return require('@/assets/image/putongrenwu.svg')        // 普通任务
    case 2: return require('@/assets/image/dingshirenwu.svg')    // 定时任务
    case 3: return require('@/assets/image/ansible.svg')    // Ansible任务
    case 4: return require('@/assets/image/renwu.svg')      // 工作作业
    default: return require('@/assets/image/普通.svg')
  }
}

const getTypeName = (type) => {
  switch(type) {
    case 1: return '普通任务'
    case 2: return '定时任务'
    case 3: return 'Ansible'
    case 4: return '工作作业'
    default: return '普通任务'
  }
}

const calculateNextExecution = async () => {
  if (!currentTask.value.cronExpr || currentTask.value.type !== 2) {
    nextExecutionTime.value = ''
    return
  }

  try {
    const response = await GetNextExecutionTime({ cron: currentTask.value.cronExpr })
    if (response.data?.data?.next_execution_time) {
      nextExecutionTime.value = formatTime(response.data.data.next_execution_time)
    } else {
      nextExecutionTime.value = '无法计算'
    }
  } catch (error) {
    console.error('计算下次执行时间失败:', error)
    nextExecutionTime.value = '计算失败'
  }
}

const fetchTemplates = async () => {
  try {
    const response = await GetAllTemplates()
    if (response.data?.code === 200) {
      templates.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取模板列表失败:', error)
  }
}

const fetchGroupsWithHosts = async () => {
  try {
    const response = await cmdbAPI.getGroupListWithHosts()
    if (response.data?.code === 200) {
      groupsWithHosts.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取分组和主机列表失败:', error)
  }
}

const handleAddTemplate = () => {
  showTemplateDialog.value = true
}

const getTemplateName = (templateId) => {
  const template = templates.value.find(t => t.id === templateId)
  return template ? template.name : '未知模板'
}

const removeTemplate = (templateId) => {
  selectedTemplates.value = selectedTemplates.value.filter(id => id !== templateId)
  updateTaskShellContent()
}

const handleTemplateSelected = (templateIds) => {
  // 确保templateIds是数组
  const ids = Array.isArray(templateIds) ? templateIds : []
  selectedTemplates.value = [...ids] // 创建新数组确保响应性
  currentTask.value.shell = ids.join(',')
  showTemplateDialog.value = false
}

// 移除getTemplatePlaceholder方法

const updateTaskShellContent = () => {
  currentTask.value.shell = selectedTemplates.value.join(',')
}

const handleAddHost = () => {
  showHostDialog.value = true
}

const handleHostSelected = (newHosts) => {
  // 确保newHosts是数组
  const hosts = Array.isArray(newHosts) ? newHosts : []
  const existingIds = new Set(selectedHosts.value.map(h => h.id))

  selectedHosts.value = [
    ...selectedHosts.value,
    ...hosts.filter(host => host?.id && !existingIds.has(host.id))
  ]

  currentTask.value.hostIDs = selectedHosts.value
    .filter(h => h?.id)
    .map(h => h.id)
    .join(',')
  showHostDialog.value = false
}

const removeHost = (hostId) => {
  selectedHosts.value = selectedHosts.value.filter(h => h?.id !== hostId)
  currentTask.value.hostIDs = selectedHosts.value
    .filter(h => h?.id)
    .map(h => h.id)
    .join(',')
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).replace(/\//g, '-')
}

// 切换任务开关（暂停/恢复）
const handleToggleTask = async (row) => {
  const isResuming = row.isActive // true表示要启动(恢复)，false表示要暂停
  const action = isResuming ? '启动' : '暂停'
  const message = isResuming
    ? '确定要启动该定时任务吗？启动后任务将重新开始自动执行。'
    : '确定要暂停该定时任务吗？暂停后任务将不再自动执行。'

  try {
    await ElMessageBox.confirm(message, `${action}任务`, {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: isResuming ? 'info' : 'warning'
    })

    loading.value = true
    const apiCall = isResuming ? ResumeScheduledTask : PauseScheduledTask
    const response = await apiCall(row.id)

    if (response.data?.code === 200) {
      ElMessage.success(`任务${action}成功`)
      // 更新本地状态
      row.status = isResuming ? 2 : 5
      fetchTasks()
    } else {
      throw new Error(response.data?.message || `${action}任务失败`)
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error(`${action}任务失败:`, error)
      ElMessage.error(`${action}失败: ${error.response?.data?.message || error.message || '未知错误'}`)
      // 恢复开关状态
      row.isActive = !row.isActive
    } else {
      // 用户取消，恢复开关状态
      row.isActive = !row.isActive
    }
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    await Promise.all([
      fetchTasks(),
      fetchTemplates(),
      fetchGroupsWithHosts()
    ])
  } catch (error) {
    console.error('初始化失败:', error)
  }
})
</script>

<style scoped>
.taskjob-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.taskjob-card {
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

.type-normal {
  color: #409EFF;
  font-weight: 500;
}

.type-cron {
  color: #E6A23C;
  font-weight: 500;
}

.type-ansible {
  color: #67C23A;
  font-weight: 500;
}

.type-job {
  color: #909399;
  font-weight: 500;
}

.status-waiting {
  color: #909399;
  font-weight: 500;
}

.status-running {
  color: #E6A23C;
  font-weight: 500;
}

.status-success {
  color: #67C23A;
  font-weight: 500;
}

.status-error {
  color: #F56C6C;
  font-weight: 500;
}

.status-paused {
  color: #909399;
  font-weight: 500;
}

.execute-count {
  color: #409EFF;
  font-weight: 600;
}

.next-run-time {
  color: #67C23A;
  font-weight: 500;
  font-size: 13px;
}
</style>
