<script setup>
import { ref, reactive, onMounted, computed, defineOptions } from 'vue'
import { useRouter } from 'vue-router'

defineOptions({
  name: 'AppUpgrade'
})
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Delete,
  View,
  Search,
  Refresh,
  Document,
  Timer,
  User,
  CircleCheck,
  CircleClose,
  Loading,
  Warning,
  Setting,
  Check,
  Folder,
  FolderOpened
} from '@element-plus/icons-vue'
import appApi from '@/api/app'
import systemApi from '@/api/system'
import cmdbApi from '@/api/cmdb'

const router = useRouter()
const loading = ref(false)

// 列表数据
const releaseList = ref([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 搜索条件
const searchForm = reactive({
  keyword: '',
  status: '',
  business_group_id: ''
})

// 业务组选项
const businessGroups = ref([])

// 用户选项
const userOptions = ref([])

// 创建工单对话框
const createDialogVisible = ref(false)

// 表单数据
const releaseForm = ref({
  title: '',
  business_group_id: null,
  impact_feature: '',
  owner_approver_id: null,
  security_approver_id: null,
  test_approver_id: null,
  services: []
})

// 应用选择弹框相关
const appSelectionDialogVisible = ref(false)
const businessTreeData = ref([]) // 业务树数据
const currentBusinessGroup = ref(null) // 当前选中的业务组
const businessGroupApps = ref([]) // 当前业务组的应用列表
const selectedAppsTemp = ref([]) // 临时选中的应用
const appTable = ref() // 应用表格引用
const loadingApps = ref(false)
const confirmingSelection = ref(false) // 确认选择loading状态

// Jenkins任务参数配置
const jenkinsParamsDialogVisible = ref(false)
const jenkinsParamsLoading = ref(false)
const jenkinsJobParameters = ref([])
const jenkinsParamValues = ref({})
const currentAppConfig = ref({}) // 当前配置参数的应用信息
const currentServiceIndex = ref(-1) // 当前配置参数的服务索引

// 状态选项
const statusOptions = [
  { value: '', label: '全部' },
  { value: 1, label: '流程进行中' },
  { value: 2, label: '发布成功' },
  { value: 3, label: '发布失败' },
  { value: 4, label: '已驳回' },
  { value: 5, label: '已取消' }
]

// 获取状态标签类型
const getStatusType = (status) => {
  const map = {
    1: 'primary',
    2: 'success',
    3: 'danger',
    4: 'warning',
    5: 'info'
  }
  return map[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    1: '流程进行中',
    2: '发布成功',
    3: '发布失败',
    4: '已驳回',
    5: '已取消'
  }
  return map[status] || '未知'
}

// 获取审批状态标签类型
const getApprovalStatusType = (status) => {
  const map = {
    1: 'info',
    2: 'success',
    3: 'danger'
  }
  return map[status] || 'info'
}

// 获取审批状态文本
const getApprovalStatusText = (status) => {
  const map = {
    1: '未审批',
    2: '通过',
    3: '驳回'
  }
  return map[status] || '未知'
}



// 获取发布状态类型
const getDeployStatusType = (status) => {
  const map = {
    1: 'info',
    2: 'primary',
    3: 'success',
    4: 'danger'
  }
  return map[status] || 'info'
}

// 获取回归测试状态文本
const getRegressionStatusText = (status) => {
  const map = {
    1: '未测试',
    2: '测试中',
    3: '测试成功',
    4: '测试失败'
  }
  return map[status] || '未知'
}

// 获取回归测试状态类型
const getRegressionStatusType = (status) => {
  const map = {
    1: 'info',
    2: 'primary',
    3: 'success',
    4: 'danger'
  }
  return map[status] || 'info'
}

// 获取列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...searchForm
    }
    // 移除空值
    Object.keys(params).forEach(key => {
      if (params[key] === '' || params[key] === null || params[key] === undefined) {
        delete params[key]
      }
    })

    const response = await appApi.getServiceReleaseList(params)
    if (response.data?.code === 200) {
      const list = response.data.data?.list || []

      // 处理每条记录，确保业务线名称和服务名称正确显示
      releaseList.value = list.map(item => {
        // 如果后端没有返回business_group_name，从本地业务组数据中查找
        if (!item.business_group_name && item.business_group_id) {
          item.business_group_name = getGroupName(item.business_group_id)
        }

        // 如果后端返回的是services数组，提取服务名称
        if (item.services && Array.isArray(item.services) && item.services.length > 0) {
          item.service_names = item.services.map(s => s.app_name || s.project_name).filter(Boolean).join(', ')
        }

        return item
      })

      pagination.total = response.data.data?.total || 0
    } else {
      ElMessage.error(response.data?.msg || '获取列表失败')
    }
  } catch (error) {
    console.error('获取列表失败:', error)
    ElMessage.error('获取列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

// 重置搜索
const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  searchForm.business_group_id = ''
  pagination.page = 1
  fetchList()
}

// 分页变化
const handlePageChange = (page) => {
  pagination.page = page
  fetchList()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchList()
}

// 查看详情 - 跳转到工单模版详情页
const handleViewDetail = (row) => {
  router.push(`/app/service-release/${row.id}`)
}

// 删除工单
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除工单 "${row.title}" 吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const response = await appApi.deleteServiceRelease(row.id)
    if (response.data?.code === 200) {
      ElMessage.success('删除成功')
      fetchList()
    } else {
      ElMessage.error(response.data?.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 打开创建对话框
const handleCreate = () => {
  resetForm()
  createDialogVisible.value = true
}

// 重置表单
const resetForm = () => {
  releaseForm.value = {
    title: '',
    business_group_id: null,
    impact_feature: '',
    owner_approver_id: null,
    security_approver_id: null,
    test_approver_id: null,
    services: []
  }
}

// 删除服务
const removeService = (index) => {
  releaseForm.value.services.splice(index, 1)
}

// 提交创建
const handleSubmit = async () => {
  // 验证
  if (!releaseForm.value.title) {
    ElMessage.warning('请填写发布标题')
    return
  }
  if (!releaseForm.value.business_group_id) {
    ElMessage.warning('请选择业务组')
    return
  }
  if (!releaseForm.value.owner_approver_id) {
    ElMessage.warning('请选择负责人审批人')
    return
  }
  if (!releaseForm.value.security_approver_id) {
    ElMessage.warning('请选择安全审批人')
    return
  }
  if (!releaseForm.value.test_approver_id) {
    ElMessage.warning('请选择测试审批人')
    return
  }
  if (releaseForm.value.services.length === 0) {
    ElMessage.warning('请至少添加一个服务')
    return
  }

  // 验证服务
  for (let i = 0; i < releaseForm.value.services.length; i++) {
    const s = releaseForm.value.services[i]
    if (!s.app_id) {
      ElMessage.warning(`请选择第${i + 1}个服务的应用`)
      return
    }
  }

  try {
    // 构建提交数据
    const submitData = {
      title: releaseForm.value.title,
      business_group_id: releaseForm.value.business_group_id,
      impact_feature: releaseForm.value.impact_feature,
      owner_approver_id: releaseForm.value.owner_approver_id,
      security_approver_id: releaseForm.value.security_approver_id,
      test_approver_id: releaseForm.value.test_approver_id,
      services: releaseForm.value.services.map(s => {
        const serviceData = {
          app_id: s.app_id,
          project_name: s.project_name,
          repo_url: s.repo_url,
          branch: s.branch,
          impact_feature: s.impact_feature,
          db_change: s.db_change,
          config_change: s.config_change,
          remark: s.remark,
          environment: 'prod'
        }

        // 如果配置了参数,将参数转换为字符串格式
        if (s.parameters) {
          serviceData.parameters = {}
          Object.keys(s.parameters).forEach(key => {
            const value = s.parameters[key]
            // 处理不同类型的参数
            if (Array.isArray(value)) {
              // 多选参数,用逗号连接
              serviceData.parameters[key] = value.join(',')
            } else if (typeof value === 'boolean') {
              // 布尔参数转换为字符串
              serviceData.parameters[key] = value ? 'true' : 'false'
            } else {
              serviceData.parameters[key] = value !== null && value !== undefined ? String(value) : ''
            }
          })
        }

        return serviceData
      })
    }

    const response = await appApi.createServiceRelease(submitData)
    if (response.data?.code === 200) {
      ElMessage.success('工单创建成功')
      createDialogVisible.value = false
      fetchList()
    } else {
      ElMessage.error(response.data?.msg || '创建失败')
    }
  } catch (error) {
    console.error('创建失败:', error)
    ElMessage.error('创建失败')
  }
}

// 获取业务组
const fetchBusinessGroups = async () => {
  try {
    const response = await cmdbApi.getAllCmdbGroups()
    if (response.data?.code === 200) {
      businessGroups.value = response.data.data || []

      // 同时构建业务树数据(用于应用选择弹框)
      businessTreeData.value = businessGroups.value.map(group => ({
        id: group.id,
        name: group.name,
        type: 'group',
        children: group.children?.map(child => ({
          id: child.id,
          name: child.name,
          type: 'subgroup',
          parentId: group.id,
          parentName: group.name
        })) || []
      }))
    }
  } catch (error) {
    console.error('获取业务组失败:', error)
  }
}

// 获取用户列表
const fetchUserOptions = async () => {
  try {
    const response = await systemApi.queryAdminList({
      pageNum: 1,
      pageSize: 1000,
      status: '1'
    })
    if (response.data?.code === 200) {
      userOptions.value = response.data.data?.list?.map(user => ({
        id: user.id,
        name: user.nickname || user.username
      })).filter(user => user.name) || []
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
  }
}

// 显示应用选择弹框
const showAppSelectionDialog = () => {
  appSelectionDialogVisible.value = true
  currentBusinessGroup.value = null
  businessGroupApps.value = []
  selectedAppsTemp.value = []
}

// 处理业务树节点点击
const handleBusinessTreeClick = async (data) => {
  // 只处理二级节点(子分组)点击
  if (data.type === 'subgroup') {
    currentBusinessGroup.value = {
      id: data.id,
      name: data.name,
      parentId: data.parentId,
      parentName: data.parentName
    }
    await loadBusinessGroupApps(data.id)
  }
}

// 加载业务组下的应用
const loadBusinessGroupApps = async (businessGroupId) => {
  if (!businessGroupId) {
    businessGroupApps.value = []
    return
  }

  loadingApps.value = true
  try {
    const params = {
      business_group_ids: [businessGroupId]
    }
    const response = await appApi.getServiceTree(params)
    if (response.data.code === 200) {
      const serviceTreeData = response.data.data || []

      // 提取应用列表
      const apps = []
      serviceTreeData.forEach(businessLine => {
        businessLine.services?.forEach(service => {
          apps.push({
            id: service.id,
            name: service.name,
            code: service.code,
            status: service.status,
            repo_url: service.repo_url,
            business_group_id: businessLine.business_group_id,
            business_group_name: businessLine.business_group_name,
            jenkins_envs: service.jenkins_envs || []
          })
        })
      })

      businessGroupApps.value = apps
    }
  } catch (error) {
    console.error('获取业务组应用失败:', error)
    ElMessage.error('获取业务组应用失败')
  } finally {
    loadingApps.value = false
  }
}

// 处理应用表格选择变化
const handleAppSelectionChange = (selection) => {
  selectedAppsTemp.value = selection
}

// 确认应用选择
const confirmAppSelection = async () => {
  confirmingSelection.value = true
  try {
    // 为每个选中的应用获取详细信息
    for (const app of selectedAppsTemp.value) {
      // 检查是否已存在
      const exists = releaseForm.value.services.some(existing => existing.app_id === app.id)
      if (!exists) {
        try {
          // 通过应用ID获取应用详情，以获取repo_url
          const response = await appApi.getApplicationDetail(app.id)
          const appDetail = response.data?.data || {}

          // 查找prod环境的配置
          const prodEnv = app.jenkins_envs.find(env => env.env_name === 'prod')

          releaseForm.value.services.push({
            app_id: app.id,
            app_name: app.name,
            business_group_name: app.business_group_name,
            project_name: app.name,
            repo_url: appDetail.repo_url || '', // 从详情中获取repo_url
            branch: 'master',
            impact_feature: '',
            db_change: '',
            config_change: '',
            remark: '',
            environment: 'prod', // 固定为prod
            jenkins_server_id: prodEnv?.jenkins_server_id,
            job_name: prodEnv?.job_name,
            jenkins_envs: app.jenkins_envs,
            parameters: null
          })
        } catch (error) {
          console.error(`获取应用${app.name}详情失败:`, error)
          ElMessage.warning(`获取应用${app.name}的仓库地址失败`)
        }
      }
    }

    appSelectionDialogVisible.value = false
  } finally {
    confirmingSelection.value = false
  }
}

// 配置服务的Jenkins参数
const handleConfigParams = async (service, index) => {
  if (!service.jenkins_server_id || !service.job_name) {
    ElMessage.warning('该应用在prod环境未配置Jenkins任务')
    return
  }

  try {
    jenkinsParamsLoading.value = true
    currentAppConfig.value = service
    currentServiceIndex.value = index

    // 获取任务参数
    const response = await appApi.getJenkinsJobParameters(service.jenkins_server_id, service.job_name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      jenkinsJobParameters.value = responseData.data?.parameters || []

      // 初始化参数值(使用已保存的值或默认值)
      jenkinsParamValues.value = {}
      jenkinsJobParameters.value.forEach(param => {
        // PT_CHECKBOX类型需要特殊处理
        if (param.type === 'PT_CHECKBOX') {
          if (service.parameters?.[param.name]) {
            const savedValue = service.parameters[param.name]
            jenkinsParamValues.value[param.name] = typeof savedValue === 'string'
              ? savedValue.split(',').map(v => v.trim()).filter(v => v)
              : (Array.isArray(savedValue) ? savedValue : [])
          } else {
            jenkinsParamValues.value[param.name] = param.defaultValue
              ? (typeof param.defaultValue === 'string'
                  ? param.defaultValue.split(',').map(v => v.trim()).filter(v => v)
                  : [])
              : []
          }
        } else {
          jenkinsParamValues.value[param.name] = service.parameters?.[param.name] ?? param.defaultValue
        }
      })

      jenkinsParamsDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取任务参数失败')
    }
  } catch (error) {
    console.error('获取Jenkins任务参数失败:', error)
    ElMessage.error(error.response?.data?.message || '获取任务参数失败')
  } finally {
    jenkinsParamsLoading.value = false
  }
}

// 保存Jenkins参数配置
const saveJenkinsParams = () => {
  if (currentServiceIndex.value >= 0 && releaseForm.value.services[currentServiceIndex.value]) {
    // 保存参数到服务对象
    releaseForm.value.services[currentServiceIndex.value].parameters = { ...jenkinsParamValues.value }
    ElMessage.success('参数配置已保存')
    closeParamsDialog()
  }
}

// 关闭参数对话框
const closeParamsDialog = () => {
  jenkinsParamsDialogVisible.value = false
  jenkinsJobParameters.value = []
  jenkinsParamValues.value = {}
  currentAppConfig.value = {}
  currentServiceIndex.value = -1
}

// 获取环境文本
const getEnvironmentText = (env) => {
  const texts = { test: 'test', dev: 'dev', prod: 'prod' }
  return texts[env] || env
}

// 获取用户名
const getUserName = (userId) => {
  const user = userOptions.value.find(u => u.id === userId)
  return user ? user.name : '-'
}

// 业务组级联选择器数据
const cascaderBusinessGroups = computed(() => {
  return businessGroups.value.map(group => ({
    value: group.id,
    label: group.name,
    children: group.children?.map(child => ({
      value: child.id,
      label: child.name
    })) || []
  }))
})

// 获取业务组名称
const getGroupName = (groupId) => {
  for (const group of businessGroups.value) {
    if (group.id === groupId) return group.name
    if (group.children) {
      const child = group.children.find(c => c.id === groupId)
      if (child) return child.name
    }
  }
  return '-'
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString()
}

onMounted(() => {
  fetchList()
  fetchBusinessGroups()
  fetchUserOptions()
})
</script>

<template>
  <div class="service-release-page">
    <el-card shadow="hover" class="main-card">
      <template #header>
        <div class="card-header">
          <span class="title">服务上线</span>
          <div class="header-actions">
            <el-button type="primary" :icon="Plus" v-authority="['app:service_work:add']" size="small" @click="handleCreate">
              新建工单
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索区域 -->
      <el-card shadow="hover" class="search-card">
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="关键词">
            <el-input
              v-model="searchForm.keyword"
              placeholder="搜索标题"
              clearable
              size="small"
              style="width: 200px"
              @keyup.enter="handleSearch"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" size="small" clearable style="width: 140px">
              <el-option
                v-for="item in statusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :icon="Search" size="small" @click="handleSearch">搜索</el-button>
            <el-button type="warning" :icon="Refresh" size="small" @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 列表 -->
      <el-card shadow="never" class="table-card">
        <el-table :data="releaseList" v-loading="loading" stripe>
          <el-table-column prop="id" label="工单ID" width="80" align="center">
            <template #default="{ row }">
              <el-link type="primary" @click="handleViewDetail(row)">
                #{{ row.id }}
              </el-link>
            </template>
          </el-table-column>
          <el-table-column prop="applicant_name" label="申请人" width="100" align="center" />
          <el-table-column prop="title" label="发布标题" min-width="180" show-overflow-tooltip />
          <el-table-column prop="created_at" label="创建时间" width="170" align="center">
            <template #default="{ row }">
              {{ formatTime(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="流程状态" width="120" align="center">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="business_group_name" label="业务线" width="150" show-overflow-tooltip align="center" />
          <el-table-column prop="service_names" label="服务名称" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-if="row.service_names">{{ row.service_names }}</span>
              <span v-else class="no-data">-</span>
            </template>
          </el-table-column>
          <el-table-column prop="impact_feature" label="影响功能" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-if="row.impact_feature">{{ row.impact_feature }}</span>
              <span v-else class="no-data">-</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right" align="center">
            <template #default="{ row }">
              <el-tooltip content="查看详情">
                <el-button :icon="View" v-authority="['app:service_work:get']" size="small" type="primary" @click="handleViewDetail(row)">详情</el-button>
              </el-tooltip>
              <el-tooltip content="删除">
                <el-button
                  :icon="Delete"
                  size="small"
                  v-authority="['app:service_work:delete']"
                  type="danger"
                  @click="handleDelete(row)"
                  :disabled="row.deploy_status === 2"
                >删除</el-button>
              </el-tooltip>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
          />
        </div>
      </el-card>
    </el-card>

    <!-- 创建工单对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="新建服务上线工单"
      width="900px"
      class="create-dialog"
    >
      <el-form :model="releaseForm" label-width="100px" size="small">
        <el-form-item label="发布标题" required>
          <el-input v-model="releaseForm.title" placeholder="请输入发布标题" />
        </el-form-item>

        <div class="form-row">
          <el-form-item label="业务组" required class="form-item-half">
            <el-cascader
              v-model="releaseForm.business_group_id"
              :options="cascaderBusinessGroups"
              placeholder="请选择业务组"
              style="width: 100%"
              :props="{ expandTrigger: 'hover', checkStrictly: true, emitPath: false }"
            />
          </el-form-item>
        </div>

        <el-form-item label="影响功能" required>
          <el-input
            v-model="releaseForm.impact_feature"
            type="textarea"
            :rows="2"
            placeholder="描述本次发布影响的功能"
          />
        </el-form-item>

        <div class="form-row">
          <el-form-item label="负责人审批" required class="form-item-third">
            <el-select v-model="releaseForm.owner_approver_id" placeholder="选择" filterable style="width: 100%">
              <el-option v-for="user in userOptions" :key="user.id" :label="user.name" :value="user.id" />
            </el-select>
          </el-form-item>
          <el-form-item label="安全审批" required class="form-item-third">
            <el-select v-model="releaseForm.security_approver_id" placeholder="选择" filterable style="width: 100%">
              <el-option v-for="user in userOptions" :key="user.id" :label="user.name" :value="user.id" />
            </el-select>
          </el-form-item>
          <el-form-item label="测试审批" required class="form-item-third">
            <el-select v-model="releaseForm.test_approver_id" placeholder="选择" filterable style="width: 100%">
              <el-option v-for="user in userOptions" :key="user.id" :label="user.name" :value="user.id" />
            </el-select>
          </el-form-item>
        </div>

        <!-- 服务列表 -->
        <el-form-item label="服务配置" required>
          <div class="applications-section">
            <div class="applications-header">
              <el-button type="primary" size="small" @click="showAppSelectionDialog" :icon="Plus">
                添加服务
              </el-button>
              <span class="app-count" v-if="releaseForm.services.length > 0">
                已选择 {{ releaseForm.services.length }} 个服务
              </span>
            </div>

            <!-- 已选择的服务列表 -->
            <div class="selected-apps-list" v-if="releaseForm.services.length > 0">
              <div
                v-for="(service, index) in releaseForm.services"
                :key="`${service.app_id}-${index}`"
                class="selected-service-card">
                <!-- 服务头部信息 -->
                <div class="service-card-header">
                  <div class="service-title-row">
                    <span class="service-title">{{ service.app_name || service.project_name }} - 详细信息</span>
                    <div class="service-tags">
                      <el-tag size="small" type="danger">prod</el-tag>
                    </div>
                  </div>
                  <div class="service-actions-row">
                    <el-button
                      size="small"
                      type="primary"
                      plain
                      :disabled="!service.jenkins_server_id || !service.job_name"
                      @click="handleConfigParams(service, index)">
                      <el-icon><Setting /></el-icon>
                      {{ service.parameters ? '已配置参数' : '配置参数' }}
                    </el-button>
                    <el-button
                      size="small"
                      type="danger"
                      :icon="Delete"
                      @click="removeService(index)">
                      移除
                    </el-button>
                  </div>
                </div>

                <!-- 服务详细信息表单 -->
                <el-form label-width="100px" size="small" class="service-form">
                  <div class="form-row">
                    <el-form-item label="Git仓库" class="form-item-full">
                      <el-input v-model="service.repo_url" placeholder="Git仓库地址" disabled class="repo-url-input" />
                    </el-form-item>
                  </div>
                  <el-form-item label="分支">
                    <el-input v-model="service.branch" placeholder="master" disabled class="repo-url-input" />
                  </el-form-item>
                  <el-form-item label="影响功能">
                    <el-input v-model="service.impact_feature" type="textarea" :rows="2" placeholder="描述本次变更影响的功能" />
                  </el-form-item>
                  <el-form-item label="SQL变更">
                    <el-input v-model="service.db_change" type="textarea" :rows="2" placeholder="描述数据库变更内容,如:新增表、修改字段等" />
                  </el-form-item>
                  <el-form-item label="配置变更">
                    <el-input v-model="service.config_change" type="textarea" :rows="2" placeholder="描述配置文件变更内容" />
                  </el-form-item>
                  <el-form-item label="备注信息">
                    <el-input v-model="service.remark" type="textarea" :rows="2" placeholder="其他需要说明的信息" />
                  </el-form-item>
                </el-form>
              </div>
            </div>

            <!-- 空状态提示 -->
            <div class="empty-apps" v-if="releaseForm.services.length === 0">
              <el-empty :image-size="100" description="请添加要发布的服务">
                <template #image>
                  <el-icon size="100" color="#ccc"><Document /></el-icon>
                </template>
              </el-empty>
            </div>
          </div>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="createDialogVisible = false" size="small">取消</el-button>
        <el-button type="primary" @click="handleSubmit" size="small">提交工单</el-button>
      </template>
    </el-dialog>

    <!-- 应用选择弹框 -->
    <el-dialog
      v-model="appSelectionDialogVisible"
      title="选择服务"
      width="1000px"
      class="modern-dialog app-selection-dialog">
      <div class="host-selection-layout">
        <!-- 左侧业务树 -->
        <div class="tree-section">
          <div class="tree-header">
            <h4>业务树</h4>
          </div>
          <el-tree
            :data="businessTreeData"
            :props="{
              label: 'name',
              children: 'children'
            }"
            node-key="id"
            highlight-current
            @node-click="handleBusinessTreeClick"
          >
            <template #default="{ node, data }">
              <span class="tree-node">
                <el-icon v-if="data.children?.length > 0">
                  <component :is="node.expanded ? 'FolderOpened' : 'Folder'" />
                </el-icon>
                <el-icon v-else><Document /></el-icon>
                {{ node.label }}
              </span>
            </template>
          </el-tree>
        </div>

        <!-- 右侧应用表格 -->
        <div class="services-section">
          <div class="services-header">
            <h4>服务列表</h4>
            <span v-if="currentBusinessGroup" class="selected-group">{{ currentBusinessGroup.name }}</span>
          </div>

          <el-table
            :data="businessGroupApps"
            border
            style="width: 100%"
            @selection-change="handleAppSelectionChange"
            ref="appTable"
            v-loading="loadingApps">
            <el-table-column
              type="selection"
              width="55"
              :selectable="(row) => row.status === 2"
            />
            <el-table-column prop="name" label="服务名称" min-width="150">
              <template #default="{ row }">
                <div class="app-name-cell">
                  <el-icon><Document /></el-icon>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="code" label="服务代码" width="120">
              <template #default="{ row }">
                <el-tag type="info" size="small">{{ row.code }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 2 ? 'success' : 'danger'" size="small">
                  {{ row.status === 2 ? '已激活' : '未激活' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="jenkins_envs" label="Prod环境" min-width="120">
              <template #default="{ row }">
                <el-tag
                  v-if="row.jenkins_envs.find(env => env.env_name === 'prod')"
                  size="small"
                  type="danger">
                  prod
                </el-tag>
                <span v-else class="no-env">未配置</span>
              </template>
            </el-table-column>
          </el-table>

          <div class="empty-state" v-if="!currentBusinessGroup">
            <el-empty :image-size="100" description="请选择左侧业务组查看服务">
              <template #image>
                <el-icon size="100" color="#ccc"><Document /></el-icon>
              </template>
            </el-empty>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="appSelectionDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click="confirmAppSelection"
          :disabled="selectedAppsTemp.length === 0"
          :loading="confirmingSelection">
          确认选择 ({{ selectedAppsTemp.length }})
        </el-button>
      </template>
    </el-dialog>

    <!-- Jenkins任务参数配置对话框 -->
    <el-dialog
      v-model="jenkinsParamsDialogVisible"
      :title="`配置参数: ${currentAppConfig.app_name || currentAppConfig.project_name}`"
      width="800px"
      class="modern-dialog jenkins-params-dialog"
      @close="closeParamsDialog"
    >
      <el-card shadow="hover" v-loading="jenkinsParamsLoading">
        <div v-if="jenkinsJobParameters.length === 0" class="no-params">
          <el-empty description="该任务无需参数，可直接发布" />
        </div>

        <el-form v-else label-width="150px" class="params-form">
          <div v-for="param in jenkinsJobParameters" :key="param.name" class="param-item">
            <!-- String 类型 -->
            <el-form-item
              v-if="param.type === 'String'"
              :label="param.name"
              :required="param.required"
            >
              <el-input
                v-model="jenkinsParamValues[param.name]"
                :placeholder="param.description || `请输入${param.name}`"
                clearable
              />
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>

            <!-- Choice 类型(单选) -->
            <el-form-item
              v-else-if="param.type === 'Choice' && param.choices && param.choices.length > 0"
              :label="param.name"
              :required="param.required"
            >
              <el-select
                v-model="jenkinsParamValues[param.name]"
                :placeholder="param.description || `请选择${param.name}`"
                clearable
                style="width: 100%"
              >
                <el-option
                  v-for="choice in param.choices"
                  :key="choice"
                  :label="choice"
                  :value="choice"
                />
              </el-select>
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>

            <!-- Boolean 类型 -->
            <el-form-item
              v-else-if="param.type === 'Boolean'"
              :label="param.name"
            >
              <el-switch
                v-model="jenkinsParamValues[param.name]"
                :active-text="param.description || '是'"
                inactive-text="否"
              />
            </el-form-item>

            <!-- Text 类型 -->
            <el-form-item
              v-else-if="param.type === 'Text'"
              :label="param.name"
              :required="param.required"
            >
              <el-input
                v-model="jenkinsParamValues[param.name]"
                type="textarea"
                :rows="5"
                :placeholder="param.description || `请输入${param.name}`"
              />
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>

            <!-- Password 类型 -->
            <el-form-item
              v-else-if="param.type === 'Password'"
              :label="param.name"
              :required="param.required"
            >
              <el-input
                v-model="jenkinsParamValues[param.name]"
                type="password"
                :placeholder="param.description || `请输入${param.name}`"
                show-password
                clearable
              />
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>

            <!-- PT_CHECKBOX 类型(多选框) -->
            <el-form-item
              v-else-if="param.type === 'PT_CHECKBOX' && param.choices && param.choices.length > 0"
              :label="param.name"
              :required="param.required"
            >
              <el-checkbox-group v-model="jenkinsParamValues[param.name]">
                <el-checkbox
                  v-for="choice in param.choices"
                  :key="choice"
                  :label="choice"
                  :value="choice"
                >
                  {{ choice }}
                </el-checkbox>
              </el-checkbox-group>
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>

            <!-- 默认类型(Text输入) -->
            <el-form-item
              v-else
              :label="param.name"
              :required="param.required"
            >
              <el-input
                v-model="jenkinsParamValues[param.name]"
                :placeholder="param.description || `请输入${param.name}`"
                clearable
              />
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>
          </div>
        </el-form>
      </el-card>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeParamsDialog">取消</el-button>
          <el-button type="primary" @click="saveJenkinsParams">
            <el-icon><Check /></el-icon>
            保存配置
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.service-release-page {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.main-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(45deg, #667eea, #764ba2);
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.search-card {
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  margin-bottom: 20px;
}

.search-form {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.search-form :deep(.el-form-item) {
  margin-bottom: 0;
}

.table-card {
  background: rgba(255, 255, 255, 0.8);
  border-radius: 12px;
  border: none;
  box-shadow: none;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

.approval-status {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  padding: 20px;
  background: rgba(248, 249, 250, 0.8);
  border-radius: 0 0 12px 12px;
}

/* 创建对话框样式 */
.create-dialog :deep(.el-dialog) {
  border-radius: 16px;
}

.create-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.create-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.create-dialog :deep(.el-dialog__body) {
  max-height: 60vh;
  overflow-y: auto;
}

.form-row {
  display: flex;
  gap: 16px;
}

.form-item-half {
  flex: 1;
}

.form-item-third {
  flex: 1;
}

.form-item-full {
  flex: 1;
  width: 100%;
}

/* 应用配置区域样式 */
.applications-section {
  width: 100%;
}

.applications-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.app-count {
  color: #606266;
  font-size: 14px;
}

.selected-apps-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.selected-service-card {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #fff;
  overflow: hidden;
}

.service-card-header {
  background: linear-gradient(135deg, #f5f7fa, #e8edf3);
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;
}

.service-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.service-title {
  font-weight: 600;
  color: #409eff;
  font-size: 15px;
}

.service-tags {
  display: flex;
  gap: 8px;
}

.service-actions-row {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.service-form {
  padding: 16px;
}

/* Git仓库和分支文字黄色高亮样式 */
.repo-url-input :deep(.el-input__wrapper) {
  --el-input-text-color: #faad14 !important;
  --el-input-disabled-text-color: #faad14 !important;
}

.repo-url-input :deep(.el-input__inner) {
  color: #faad14 !important;
  font-weight: 600 !important;
  -webkit-text-fill-color: #faad14 !important;
}

.repo-url-input :deep(input.el-input__inner) {
  color: #faad14 !important;
  font-weight: 600 !important;
}

.empty-apps {
  text-align: center;
  padding: 40px 0;
}

/* 应用选择弹框样式 - 左右布局 */
.app-selection-dialog :deep(.el-dialog__body) {
  padding: 20px;
  height: 500px;
}

.host-selection-layout {
  display: flex;
  height: 100%;
  gap: 16px;
}

/* 左侧业务树 */
.tree-section {
  width: 280px;
  border-right: 1px solid #ebeef5;
  padding-right: 16px;
  overflow-y: auto;
}

.tree-header {
  margin-bottom: 12px;
}

.tree-header h4 {
  margin: 0;
  color: #303133;
  font-size: 16px;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 6px;
}

.tree-node .el-icon {
  color: #409eff;
}

/* 右侧应用表格 */
.services-section {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.services-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.services-header h4 {
  margin: 0;
  color: #303133;
  font-size: 16px;
}

.selected-group {
  color: #409eff;
  font-size: 14px;
  font-weight: 500;
}

.app-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-name-cell .el-icon {
  color: #409eff;
}

.no-env {
  color: #909399;
  font-size: 12px;
  font-style: italic;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modern-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  padding: 20px 24px;
}

.modern-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

/* Jenkins参数对话框样式 */
.jenkins-params-dialog :deep(.el-dialog__body) {
  padding: 20px;
}

.jenkins-params-dialog .no-params {
  padding: 40px;
  text-align: center;
}

.jenkins-params-dialog .params-form {
  max-height: 500px;
  overflow-y: auto;
  padding: 20px;
}

.jenkins-params-dialog .param-item {
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f2f5;
}

.jenkins-params-dialog .param-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.jenkins-params-dialog .param-description {
  display: block;
  font-size: 12px;
  color: #909399;
  margin-top: 6px;
  line-height: 1.4;
}

.jenkins-params-dialog :deep(.el-form-item__label) {
  font-weight: 600;
  color: #303133;
}

.jenkins-params-dialog :deep(.el-switch) {
  --el-switch-on-color: #67c23a;
}

.jenkins-params-dialog :deep(.el-checkbox-group) {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.jenkins-params-dialog :deep(.el-checkbox) {
  margin-right: 0;
  height: auto;
}

.service-item {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-weight: 600;
  color: #409EFF;
}

.no-data {
  color: #909399;
  font-style: italic;
}
</style>
