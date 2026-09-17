<!-- eslint-disable vue/multi-word-component-names -->
<script setup>
import { ref, reactive, onMounted, computed, defineOptions } from 'vue'
import { useRouter } from 'vue-router'

defineOptions({
  name: 'AppApplication'
})
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Edit,
  Delete,
  View,
  Search,
  Refresh,
  Monitor,
  Setting,
  SuccessFilled,
  WarningFilled,
  Close,
  Promotion,
  ArrowUp,
  ArrowDown
} from '@element-plus/icons-vue'
import appApi from '@/api/app'
import cmdbApi from '@/api/cmdb'
import systemApi from '@/api/system'
import databaseApi from '@/api/database'

// 基础状态
const loading = ref(false)
const searchKeyword = ref('')
const router = useRouter()

// 应用列表数据
const applicationList = ref([])
const currentApplication = ref({})

// 对话框状态
const createDialogVisible = ref(false)
const editDialogVisible = ref(false)
const detailDialogVisible = ref(false)

// 分页参数
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 表单数据
const applicationForm = ref({
  // 必填参数
  name: '',                    // 应用名称
  business_group_id: null,     // 业务组ID
  business_dept_id: null,      // 业务部门ID

  // 基本信息（可选）
  code: '',                    // 应用编码（可选，不填则自动生成）
  description: '',             // 应用介绍
  repo_url: '',               // 仓库地址

  // 负责人信息（可选）
  dev_owners: [],             // 研发负责人ID数组
  test_owners: [],            // 测试负责人ID数组
  ops_owners: [],             // 运维负责人ID数组

  // 技术信息（可选）
  programming_lang: '',       // 编程语言
  start_command: '',          // 启动命令
  stop_command: '',           // 停止命令
  health_api: '',             // 健康检查接口

  // 关联资源（可选）
  domains: [],                // 关联域名数组
  hosts: [],                  // 关联主机ID数组
  databases: [],              // 关联数据库ID数组
  other_res: {                // 其他资源配置
    rabbitmq: [],
    zookeeper: [],
    kafka: [],
    other: []
  },

  // Jenkins环境配置（可选）
  jenkins_envs: []            // Jenkins环境配置数组
})

// 更多配置展开状态
const showMoreConfig = ref(false)

// 编程语言选项
const programmingLanguages = [
  'Java', 'Go', 'Python', 'Node.js', 'C++', 'C#', '.NET', 'PHP', 'Ruby', 'Rust'
]


// 业务组选项
const businessGroups = ref([])

// 业务部门选项
const businessDepts = ref([])

// 用户选项（负责人）
const userOptions = ref([])

// 主机选项
const hostOptions = ref([])

// 数据库选项
const databaseOptions = ref([])

// Jenkins环境配置
const jenkinsEnvDialogVisible = ref(false)
const jenkinsEnvs = ref([])
const jenkinsEnvForm = ref({
  env_name: '',
  jenkins_server_id: null,
  job_name: ''
})
const editingJenkinsEnv = ref(null)
const showJenkinsEnvForm = ref(false)
const currentAppId = ref(null)
const currentAppName = ref('')

// Jenkins服务器和任务选项
const jenkinsServers = ref([])
const jenkinsJobs = ref([])
const jobSearchLoading = ref(false)
const jobValidationLoading = ref(false)
const jobValidationStatus = ref(null) // null, 'success', 'error'

// 过滤后的应用列表
const filteredApplicationList = computed(() => {
  const list = Array.isArray(applicationList.value) ? applicationList.value : []
  if (!searchKeyword.value) return list
  return list.filter(item =>
    item.name?.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    item.code?.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    item.programming_lang?.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

// 获取应用列表
const fetchApplicationList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize
    }

    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }

    const response = await appApi.getApplicationList(params)

    if (response.data?.code === 200) {
      applicationList.value = response.data.data?.list || []
      pagination.total = response.data.data?.total || 0
    } else {
      ElMessage.error(response.data?.message || '获取应用列表失败')
    }
  } catch (error) {
    console.error('获取应用列表失败:', error)
    ElMessage.error('获取应用列表失败')
  } finally {
    loading.value = false
  }
}

// 创建应用
const handleCreateApp = () => {
  resetForm()
  createDialogVisible.value = true
}

// 编辑应用
const handleEditApp = (row) => {
  currentApplication.value = { ...row }

  // 复制基本数据
  Object.assign(applicationForm.value, row)

  // 处理资源关联字段 - 将数组转换为逗号分隔的字符串
  if (row.domains && Array.isArray(row.domains)) {
    applicationForm.value.domains = row.domains.join(', ')
  }

  if (row.hosts && Array.isArray(row.hosts)) {
    applicationForm.value.hosts = row.hosts
  }

  // 处理 databases - 需要将ID转换为组合key以匹配选项
  if (row.databases && Array.isArray(row.databases)) {
    // 将后端返回的数据库ID转换为组合key
    applicationForm.value.databases = row.databases.map(dbId => {
      // 在 databaseOptions 中查找匹配的数据库
      const matchedDb = databaseOptions.value.find(db => db.originalId === dbId)
      if (matchedDb) {
        return matchedDb.id // 返回组合key
      }
      // 如果找不到匹配的，尝试在所有选项中查找包含该ID的
      const foundDb = databaseOptions.value.find(db =>
        db.id.endsWith(`_${dbId}`)
      )
      return foundDb ? foundDb.id : dbId
    })
  }

  // 处理 other_res 对象
  if (row.other_res) {
    if (row.other_res.rabbitmq && Array.isArray(row.other_res.rabbitmq)) {
      applicationForm.value.other_res.rabbitmq = row.other_res.rabbitmq.join(', ')
    }
    if (row.other_res.zookeeper && Array.isArray(row.other_res.zookeeper)) {
      applicationForm.value.other_res.zookeeper = row.other_res.zookeeper.join(', ')
    }
    if (row.other_res.kafka && Array.isArray(row.other_res.kafka)) {
      applicationForm.value.other_res.kafka = row.other_res.kafka.join(', ')
    }
    if (row.other_res.other && Array.isArray(row.other_res.other)) {
      applicationForm.value.other_res.other = row.other_res.other.join(', ')
    }
  }

  editDialogVisible.value = true
}

// 查看应用详情
const handleViewApp = (row) => {
  currentApplication.value = { ...row }
  detailDialogVisible.value = true
}


// 处理Jenkins环境配置
const handleJenkinsEnvConfig = async (row) => {
  currentAppId.value = row.id
  currentAppName.value = row.name
  jenkinsEnvDialogVisible.value = true
  await Promise.all([
    loadJenkinsEnvs(),
    loadJenkinsServers()
  ])
}

// 加载Jenkins环境列表
const loadJenkinsEnvs = async () => {
  try {
    const response = await appApi.getAppJenkinsEnvs(currentAppId.value)
    if (response.data?.code === 200) {
      jenkinsEnvs.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取Jenkins环境列表失败:', error)
    ElMessage.error('获取Jenkins环境列表失败')
  }
}

// 添加Jenkins环境
const handleAddJenkinsEnv = () => {
  editingJenkinsEnv.value = null
  jenkinsEnvForm.value = {
    env_name: '',
    jenkins_server_id: null,
    job_name: ''
  }
  jenkinsJobs.value = []
  jobValidationStatus.value = null
  showJenkinsEnvForm.value = true
}

// 编辑Jenkins环境
const handleEditJenkinsEnv = (env) => {
  editingJenkinsEnv.value = env
  showJenkinsEnvForm.value = true
  jenkinsEnvForm.value = {
    env_name: env.env_name,
    jenkins_server_id: env.jenkins_server_id,
    job_name: env.job_name
  }
  jobValidationStatus.value = null
}

// 保存Jenkins环境
const handleSaveJenkinsEnv = async () => {
  try {
    // 验证必填字段
    if (!jenkinsEnvForm.value.env_name?.trim()) {
      ElMessage.warning('请填写环境名称')
      return
    }

    // 如果选择了Jenkins服务器和任务名称，先验证任务是否存在
    const hasServerAndJob = jenkinsEnvForm.value.jenkins_server_id &&
                           jenkinsEnvForm.value.job_name?.trim()

    if (hasServerAndJob && (!jobValidationStatus.value || !jobValidationStatus.value.exists)) {
      try {
        await ElMessageBox.confirm(
          `Jenkins任务 "${jenkinsEnvForm.value.job_name.trim()}" 在所选服务器中可能不存在，确定要保存吗？`,
          '任务验证提醒',
          {
            confirmButtonText: '确定保存',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
      } catch (cancelError) {
        return // 用户取消操作
      }
    }

    // 准备提交数据，确保符合API规范
    const data = {
      env_name: jenkinsEnvForm.value.env_name.trim(),
      jenkins_server_id: jenkinsEnvForm.value.jenkins_server_id || null,
      job_name: jenkinsEnvForm.value.job_name?.trim() || ''
    }

    // 确保如果没有选择服务器，任务名称也为空
    if (!data.jenkins_server_id) {
      data.job_name = ''
    }

    if (editingJenkinsEnv.value) {
      // 更新环境
      await appApi.updateAppJenkinsEnv(currentAppId.value, editingJenkinsEnv.value.id, data)
      ElMessage.success('Jenkins环境更新成功')
    } else {
      // 新增环境
      await appApi.addAppJenkinsEnv(currentAppId.value, data)
      ElMessage.success('Jenkins环境创建成功')
    }

    resetJenkinsEnvForm()
    await loadJenkinsEnvs()
  } catch (error) {
    console.error('保存Jenkins环境失败:', error)
    const errorMessage = error.response?.data?.message || '保存Jenkins环境失败'
    ElMessage.error(errorMessage)
  }
}

// 删除Jenkins环境
const handleDeleteJenkinsEnv = async (env) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除环境 "${env.env_name}" 吗？`,
      '删除确认',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await appApi.deleteAppJenkinsEnv(currentAppId.value, env.id)
    ElMessage.success('Jenkins环境删除成功')
    await loadJenkinsEnvs()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除Jenkins环境失败:', error)
      ElMessage.error(error.response?.data?.message || '删除Jenkins环境失败')
    }
  }
}

// 重置Jenkins环境表单
const resetJenkinsEnvForm = () => {
  editingJenkinsEnv.value = null
  showJenkinsEnvForm.value = false
  jenkinsEnvForm.value = {
    env_name: '',
    jenkins_server_id: null,
    job_name: ''
  }
  jenkinsJobs.value = []
  jobValidationStatus.value = null
}

// 加载Jenkins服务器列表
const loadJenkinsServers = async () => {
  try {
    const response = await appApi.getJenkinsServers()
    if (response.data?.code === 200) {
      jenkinsServers.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取Jenkins服务器列表失败:', error)
    ElMessage.error('获取Jenkins服务器列表失败')
  }
}

// Jenkins服务器选择变化时，重置任务名称和搜索任务
const handleServerChange = async (serverId) => {
  jenkinsEnvForm.value.job_name = ''
  jobValidationStatus.value = null
  if (serverId) {
    await searchJenkinsJobs('')
  } else {
    jenkinsJobs.value = []
  }
}

// 搜索Jenkins任务
const searchJenkinsJobs = async (searchKey = '') => {
  if (!jenkinsEnvForm.value.jenkins_server_id) {
    jenkinsJobs.value = []
    return
  }

  jobSearchLoading.value = true
  try {
    const response = await appApi.searchJenkinsJobs(
      jenkinsEnvForm.value.jenkins_server_id,
      searchKey
    )
    if (response.data?.code === 200) {
      jenkinsJobs.value = response.data.data || []
    }
  } catch (error) {
    console.error('搜索Jenkins任务失败:', error)
    ElMessage.error('搜索Jenkins任务失败')
  } finally {
    jobSearchLoading.value = false
  }
}

// 任务名称输入变化时的防抖搜索
let searchTimeout = null
const handleJobNameInput = (value) => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    if (jenkinsEnvForm.value.jenkins_server_id) {
      searchJenkinsJobs(value)
    }
  }, 500)
}

// 验证Jenkins任务是否存在
const validateJenkinsJob = async () => {
  if (!jenkinsEnvForm.value.jenkins_server_id || !jenkinsEnvForm.value.job_name) {
    jobValidationStatus.value = null
    return
  }

  jobValidationLoading.value = true
  try {
    const response = await appApi.validateJenkinsJob({
      jenkins_server_id: jenkinsEnvForm.value.jenkins_server_id,
      job_name: jenkinsEnvForm.value.job_name
    })

    if (response.data?.code === 200) {
      const result = response.data.data
      jobValidationStatus.value = {
        exists: result.exists,
        message: result.message,
        job_url: result.job_url
      }

      if (!result.exists) {
        ElMessage.warning(`任务验证失败: ${result.message}`)
      }
      // 成功时不显示消息提示，只显示图标
    }
  } catch (error) {
    console.error('验证Jenkins任务失败:', error)
    ElMessage.error('验证Jenkins任务失败')
    jobValidationStatus.value = { exists: false, message: '验证失败' }
  } finally {
    jobValidationLoading.value = false
  }
}

// 删除应用
const handleDeleteApp = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认删除应用 "${row.name}"？删除后将无法恢复。`,
      '删除确认',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    const response = await appApi.deleteApplication(row.id)

    if (response.data?.code === 200) {
      ElMessage.success('应用删除成功')
      fetchApplicationList()
    } else {
      ElMessage.error(response.data?.message || '删除应用失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除应用失败:', error)
      ElMessage.error(error.response?.data?.message || '删除应用失败')
    }
  }
}

// 保存应用
const handleSaveApp = async () => {
  try {
    // 验证必填字段
    if (!applicationForm.value.name) {
      ElMessage.warning('请填写应用名称')
      return
    }
    if (!applicationForm.value.business_group_id) {
      ElMessage.warning('请选择业务组')
      return
    }
    if (!applicationForm.value.business_dept_id) {
      ElMessage.warning('请选择业务部门')
      return
    }

    // 准备提交数据，过滤空值
    const submitData = {}

    // 必填字段
    submitData.name = applicationForm.value.name
    submitData.business_group_id = applicationForm.value.business_group_id
    submitData.business_dept_id = applicationForm.value.business_dept_id

    // 可选字段，只包含非空值
    if (applicationForm.value.code) submitData.code = applicationForm.value.code
    if (applicationForm.value.description) submitData.description = applicationForm.value.description
    if (applicationForm.value.repo_url) submitData.repo_url = applicationForm.value.repo_url
    if (applicationForm.value.programming_lang) submitData.programming_lang = applicationForm.value.programming_lang
    if (applicationForm.value.start_command) submitData.start_command = applicationForm.value.start_command
    if (applicationForm.value.stop_command) submitData.stop_command = applicationForm.value.stop_command
    if (applicationForm.value.health_api) submitData.health_api = applicationForm.value.health_api

    // 负责人数组，只包含非空数组
    if (applicationForm.value.dev_owners?.length > 0) submitData.dev_owners = applicationForm.value.dev_owners
    if (applicationForm.value.test_owners?.length > 0) submitData.test_owners = applicationForm.value.test_owners
    if (applicationForm.value.ops_owners?.length > 0) submitData.ops_owners = applicationForm.value.ops_owners

    // 关联资源，只包含非空数组
    // 处理 domains - 字符串数组
    if (applicationForm.value.domains) {
      if (typeof applicationForm.value.domains === 'string') {
        const domainArray = applicationForm.value.domains.split(',').map(d => d.trim()).filter(d => d)
        if (domainArray.length > 0) submitData.domains = domainArray
      } else if (Array.isArray(applicationForm.value.domains) && applicationForm.value.domains.length > 0) {
        submitData.domains = applicationForm.value.domains
      }
    }

    // 处理 hosts - uint数组
    if (applicationForm.value.hosts && Array.isArray(applicationForm.value.hosts) && applicationForm.value.hosts.length > 0) {
      submitData.hosts = applicationForm.value.hosts
    }

    // 处理 databases - 从组合key提取原始ID
    if (applicationForm.value.databases && Array.isArray(applicationForm.value.databases) && applicationForm.value.databases.length > 0) {
      // 如果是组合key（字符串），需要解析出原始ID
      submitData.databases = applicationForm.value.databases.map(dbKey => {
        if (typeof dbKey === 'string' && dbKey.includes('_')) {
          // 从 "dbType_sourceTable_id" 中提取 id
          const parts = dbKey.split('_')
          return parseInt(parts[parts.length - 1])
        }
        return dbKey
      })
    }

    // 其他资源，只包含非空配置
    const otherRes = {}

    // 处理 rabbitmq
    if (applicationForm.value.other_res.rabbitmq) {
      if (typeof applicationForm.value.other_res.rabbitmq === 'string') {
        const arr = applicationForm.value.other_res.rabbitmq.split(',').map(r => r.trim()).filter(r => r)
        if (arr.length > 0) otherRes.rabbitmq = arr
      } else if (Array.isArray(applicationForm.value.other_res.rabbitmq) && applicationForm.value.other_res.rabbitmq.length > 0) {
        otherRes.rabbitmq = applicationForm.value.other_res.rabbitmq
      }
    }

    // 处理 zookeeper
    if (applicationForm.value.other_res.zookeeper) {
      if (typeof applicationForm.value.other_res.zookeeper === 'string') {
        const arr = applicationForm.value.other_res.zookeeper.split(',').map(z => z.trim()).filter(z => z)
        if (arr.length > 0) otherRes.zookeeper = arr
      } else if (Array.isArray(applicationForm.value.other_res.zookeeper) && applicationForm.value.other_res.zookeeper.length > 0) {
        otherRes.zookeeper = applicationForm.value.other_res.zookeeper
      }
    }

    // 处理 kafka
    if (applicationForm.value.other_res.kafka) {
      if (typeof applicationForm.value.other_res.kafka === 'string') {
        const arr = applicationForm.value.other_res.kafka.split(',').map(k => k.trim()).filter(k => k)
        if (arr.length > 0) otherRes.kafka = arr
      } else if (Array.isArray(applicationForm.value.other_res.kafka) && applicationForm.value.other_res.kafka.length > 0) {
        otherRes.kafka = applicationForm.value.other_res.kafka
      }
    }

    // 处理 other
    if (applicationForm.value.other_res.other) {
      if (typeof applicationForm.value.other_res.other === 'string') {
        const arr = applicationForm.value.other_res.other.split(',').map(o => o.trim()).filter(o => o)
        if (arr.length > 0) otherRes.other = arr
      } else if (Array.isArray(applicationForm.value.other_res.other) && applicationForm.value.other_res.other.length > 0) {
        otherRes.other = applicationForm.value.other_res.other
      }
    }

    if (Object.keys(otherRes).length > 0) submitData.other_res = otherRes

    // Jenkins环境配置，只包含非空数组
    if (applicationForm.value.jenkins_envs?.length > 0) submitData.jenkins_envs = applicationForm.value.jenkins_envs

    let response
    if (editDialogVisible.value) {
      // 更新应用
      response = await appApi.updateApplication(currentApplication.value.id, submitData)
    } else {
      // 创建应用
      response = await appApi.createApplication(submitData)
    }

    if (response.data?.code === 200) {
      ElMessage.success(editDialogVisible.value ? '应用更新成功' : '应用创建成功')
      createDialogVisible.value = false
      editDialogVisible.value = false
      fetchApplicationList()
    } else {
      ElMessage.error(response.data?.message || '操作失败')
    }
  } catch (error) {
    console.error('保存应用失败:', error)
    ElMessage.error(error.response?.data?.message || '保存应用失败')
  }
}


// 重置表单
const resetForm = () => {
  applicationForm.value = {
    // 必填参数
    name: '',
    business_group_id: null,
    business_dept_id: null,

    // 基本信息（可选）
    code: '',
    description: '',
    repo_url: '',

    // 负责人信息（可选）
    dev_owners: [],
    test_owners: [],
    ops_owners: [],

    // 技术信息（可选）
    programming_lang: '',
    start_command: '',
    stop_command: '',
    health_api: '',

    // 关联资源（可选）
    domains: [],
    hosts: [],
    databases: [],
    other_res: {
      rabbitmq: [],
      zookeeper: [],
      kafka: [],
      other: []
    },

    // Jenkins环境配置（可选）
    jenkins_envs: []
  }
  // 重置展开状态
  showMoreConfig.value = false
}

// 刷新数据
const handleRefresh = () => {
  fetchApplicationList()
}

// 重置搜索
const resetSearch = () => {
  searchKeyword.value = ''
  fetchApplicationList()
}

// 快速发布
const handleQuickRelease = () => {
  router.push('/app/quick-release')
}

// 分页变化
const handlePageChange = (page) => {
  pagination.page = page
  fetchApplicationList()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchApplicationList()
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const statusMap = {
    1: 'warning',
    2: 'success',
    3: 'danger',
    0: 'info'
  }
  return statusMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    1: '未激活',
    2: '已激活',
    3: '异常',
    0: '未知'
  }
  return statusMap[status] || '未知'
}

// 获取业务组列表
const fetchBusinessGroups = async () => {
  try {
    const response = await cmdbApi.getAllCmdbGroups()
    if (response.data?.code === 200) {
      businessGroups.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取业务组列表失败:', error)
  }
}

// 生成级联选择器的业务组数据：只能选择二级分组
const cascaderBusinessGroups = computed(() => {
  return businessGroups.value.map(group => ({
    value: group.id,
    label: group.name,
    leaf: false, // 一级分组不是叶子节点，不可选中
    children: group.children?.filter(child => !child.children || child.children.length === 0).map(child => ({
      value: child.id,
      label: child.name,
      leaf: true // 二级分组是叶子节点，可以选中
    })) || []
  }))
})

// 获取业务部门列表
const fetchBusinessDepts = async () => {
  try {
    // 使用 queryDeptList 获取部门数据
    const response = await systemApi.queryDeptList()
    if (response.data?.code === 200) {
      businessDepts.value = (response.data.data || []).filter(dept => dept.deptName)
    } else {
      // 如果第一个API失败，尝试第二个
      try {
        const response2 = await systemApi.querySysDeptVoList()
        if (response2.data?.code === 200) {
          businessDepts.value = (response2.data.data || []).filter(dept => dept.deptName)
        }
      } catch (error2) {
        console.error('获取业务部门列表失败:', error2)
      }
    }
  } catch (error) {
    console.error('获取业务部门列表失败:', error)
  }
}

// 获取用户列表（负责人选项）
const fetchUserOptions = async () => {
  try {
    const response = await systemApi.queryAdminList({
      pageNum: 1,
      pageSize: 1000, // 获取所有用户作为负责人选项
      status: '1' // 只获取启用状态的用户 (1启用，2停用)
    })
    if (response.data?.code === 200) {
      userOptions.value = response.data.data?.list?.map(user => ({
        id: user.id,
        name: user.nickname || user.username
      })).filter(user => user.name) || [] // 只保留有名称的用户
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
  }
}

// 获取主机列表
const fetchHostOptions = async () => {
  try {
    const response = await cmdbApi.getCmdbHostList({
      page: 1,
      pageSize: 50
    })
    console.log('主机列表响应:', response)
    if (response.data?.code === 200) {
      const hostList = response.data.data?.list || []
      console.log('主机列表数据:', hostList)
      hostOptions.value = hostList.map(host => ({
        id: host.id,
        name: host.name || host.hostName || `主机${host.id}`,
        ip: host.privateIp || host.publicIp || host.sshIp || '-'
      }))
      console.log('处理后的主机选项:', hostOptions.value)
    } else {
      console.error('获取主机列表失败:', response.data?.message)
      ElMessage.warning(`获取主机列表失败: ${response.data?.message}`)
    }
  } catch (error) {
    console.error('获取主机列表失败:', error)
    ElMessage.error('获取主机列表失败')
  }
}

// 远程搜索主机
const remoteSearchHosts = async (query) => {
  if (!query || !query.trim()) {
    // 如果搜索框为空，加载默认主机列表
    await fetchHostOptions()
    return
  }

  try {
    const response = await cmdbApi.GetCmdbHostsByHostNameLike(query.trim(), {
      pageNum: 1,
      pageSize: 50
    })

    if (response.data?.code === 200) {
      const hosts = Array.isArray(response.data.data) ? response.data.data : []
      hostOptions.value = hosts.map(host => ({
        id: host.id,
        name: host.hostName || host.name || `主机${host.id}`,
        ip: host.privateIp || host.publicIp || host.sshIp || '-'
      }))
    }
  } catch (error) {
    console.error('搜索主机失败:', error)
  }
}

// 获取数据库列表
const fetchDatabaseOptions = async () => {
  try {
    const response = await databaseApi.getInstanceList({
      pageNum: 1,
      pageSize: 50
    })
    console.log('数据库列表响应:', response)
    if (response.data?.code === 200) {
      const dbList = response.data.data?.list || []
      console.log('数据库列表数据:', dbList)
      databaseOptions.value = dbList.map(db => ({
        // 使用组合key确保唯一性: dbType_sourceTable_id
        id: `${db.dbType}_${db.sourceTable}_${db.id}`,
        name: db.name || `数据库${db.id}`,
        type: db.dbType || db.subType || 'Unknown',
        originalId: db.id, // 保留原始 ID
        dbType: db.dbType,
        sourceTable: db.sourceTable
      }))
      console.log('处理后的数据库选项:', databaseOptions.value)
    } else {
      console.error('获取数据库列表失败:', response.data?.message)
      ElMessage.warning(`获取数据库列表失败: ${response.data?.message}`)
    }
  } catch (error) {
    console.error('获取数据库列表失败:', error)
    ElMessage.error('获取数据库列表失败')
  }
}

// ID映射计算属性
const getDeptNameById = (deptId) => {
  if (!deptId) return '-'
  const dept = businessDepts.value.find(d => d.id === deptId)
  return dept ? dept.deptName : `部门ID: ${deptId}`
}

const getGroupNameById = (groupId) => {
  if (!groupId) return '-'

  // 在业务组层次结构中查找指定ID的分组
  for (const group of businessGroups.value) {
    if (group.children && group.children.length > 0) {
      const subGroup = group.children.find(sg => sg.id === groupId)
      if (subGroup) {
        return subGroup.name
      }
    }
  }

  return `业务组ID: ${groupId}`
}

const getUserNamesByIds = (userIds) => {
  if (!userIds || !Array.isArray(userIds) || userIds.length === 0) return '-'
  const userNames = userIds.map(id => {
    const user = userOptions.value.find(u => u.id === id)
    return user ? user.name : `用户ID: ${id}`
  })
  return userNames.join(', ')
}

const getServerNameById = (serverId) => {
  if (!serverId) return '未配置'
  const server = jenkinsServers.value.find(s => s.id === serverId)
  return server ? server.name : `服务器${serverId}`
}

// 页面初始化
onMounted(() => {
  fetchApplicationList()
  fetchBusinessGroups()
  fetchBusinessDepts()
  fetchUserOptions()
  fetchHostOptions()
  fetchDatabaseOptions()
})
</script>

<template>
  <div class="application-management">
    <el-card shadow="hover" class="main-card">
      <template #header>
        <div class="card-header">
          <span class="title">应用管理</span>
          <div class="header-actions">
            <el-button type="primary" :icon="Plus"  v-authority="['app:application:add']" size="small" @click="handleCreateApp">
              创建应用
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索卡片 -->
      <el-card shadow="hover" class="search-card">
        <div class="search-toolbar">
          <el-form :inline="true" class="search-form">
            <el-form-item>
              <el-input
                v-model="searchKeyword"
                placeholder="搜索应用名称、编码或技术栈"
                clearable
                size="small"
                style="width: 300px"
                :prefix-icon="Search"
                @keyup.enter="resetSearch"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :icon="Search" size="small" @click="resetSearch">
                搜索
              </el-button>
              <el-button type="success" :icon="Refresh" size="small" @click="handleRefresh">
                刷新
              </el-button>
              <el-button type="warning" :icon="Promotion" size="small" @click="handleQuickRelease">
                快速发布
              </el-button>
            </el-form-item>
          </el-form>
          <div class="toolbar-right">
            <span class="resource-count">共 {{ filteredApplicationList.length }} 个应用</span>
          </div>
        </div>
      </el-card>

      <!-- 应用列表卡片 -->
      <el-card shadow="never" class="table-card">
        <el-table
          :data="filteredApplicationList"
          v-loading="loading"
          element-loading-text="加载中..."
          class="app-table"
          empty-text="暂无应用数据"
          stripe
        >
          <el-table-column prop="name" label="应用名称" min-width="250" fixed="left">
            <template #default="{ row }">
              <div class="app-name">
                <el-icon class="app-icon"><Monitor /></el-icon>
                <span class="app-name-link" @click="handleViewApp(row)">{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>


          <el-table-column prop="business_group_id" label="业务线" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <el-tag size="small" type="primary">
                {{ getGroupNameById(row.business_group_id) }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="business_dept_id" label="业务部门" width="120">
            <template #default="{ row }">
              <el-tag size="small" type="info">
                {{ getDeptNameById(row.business_dept_id) }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="dev_owners" label="研发负责人" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <span style="font-size: 12px;">{{ getUserNamesByIds(row.dev_owners) }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="test_owners" label="测试负责人" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <span style="font-size: 12px;">{{ getUserNamesByIds(row.test_owners) }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="ops_owners" label="运维负责人" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <span style="font-size: 12px;">{{ getUserNamesByIds(row.ops_owners) }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="programming_lang" label="技术栈" width="100">
            <template #default="{ row }">
              <el-tag :type="row.programming_lang === 'Java' ? 'warning' : 'primary'" size="small">
                {{ row.programming_lang || '-' }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusTagType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="created_at" label="创建时间" width="180">
            <template #default="{ row }">
              <span>{{ row.created_at ? new Date(row.created_at).toLocaleString() : '-' }}</span>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="300" fixed="right">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-tooltip content="编辑" placement="top">
                  <el-button :icon="Edit" v-authority="['app:application:edit']" size="small" type="primary" @click="handleEditApp(row)">编辑</el-button>
                </el-tooltip>
                <el-tooltip content="环境配置" placement="top">
                  <el-button :icon="Setting" v-authority="['app:application:env']" size="small" type="warning" @click="handleJenkinsEnvConfig(row)">环境配置</el-button>
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button :icon="Delete" v-authority="['app:application:delete']" size="small" type="danger" @click="handleDeleteApp(row)">删除</el-button>
                </el-tooltip>
              </div>
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

    <!-- 创建应用对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="创建应用"
      width="700px"
      class="create-dialog compact-dialog"
    >
      <el-form :model="applicationForm" label-width="90px" class="compact-form">
        <!-- 必填字段 -->
        <el-form-item label="应用名称" required class="form-item-full">
          <el-input
            v-model="applicationForm.name"
            placeholder="请输入应用名称"
            size="small"
          />
        </el-form-item>

        <!-- 业务信息行 -->
        <div class="form-row">
          <el-form-item label="业务组" required class="form-item-half">
            <el-cascader
              v-model="applicationForm.business_group_id"
              :options="cascaderBusinessGroups"
              placeholder="请选择业务组"
              style="width: 100%"
              size="small"
              :props="{
                expandTrigger: 'hover',
                checkStrictly: true,
                emitPath: false,
                leaf: 'leaf'
              }"
            />
          </el-form-item>
          <el-form-item label="业务部门" required class="form-item-half">
            <el-select
              v-model="applicationForm.business_dept_id"
              placeholder="请选择业务部门"
              style="width: 100%"
              size="small"
              filterable
              clearable
            >
              <el-option
                v-for="dept in businessDepts"
                :key="dept.id"
                :label="dept.deptName"
                :value="dept.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="仓库地址" required class="form-item-full">
          <el-input
            v-model="applicationForm.repo_url"
            placeholder="https://git.company.com/project.git"
            size="small"
          />
        </el-form-item>

        <!-- 可选字段分隔线 -->
        <div class="optional-divider">
          <div class="divider-line"></div>
          <el-button
            text
            size="small"
            class="toggle-more-btn"
            @click="showMoreConfig = !showMoreConfig"
          >
            <el-icon>
              <ArrowUp v-if="showMoreConfig" />
              <ArrowDown v-else />
            </el-icon>
            {{ showMoreConfig ? '收起更多配置' : '展开更多配置' }}
          </el-button>
          <div class="divider-line"></div>
        </div>

        <!-- 可选字段（折叠区域） -->
        <el-collapse-transition>
          <div v-show="showMoreConfig" class="more-config-section">
            <!-- 应用描述 -->
            <el-form-item label="应用描述" class="form-item-full">
              <el-input
                v-model="applicationForm.description"
                type="textarea"
                :rows="2"
                placeholder="请输入应用描述"
                size="small"
              />
            </el-form-item>

            <!-- 技术配置行 -->
            <div class="form-row">
              <el-form-item label="编程语言" class="form-item-half">
                <el-select
                  v-model="applicationForm.programming_lang"
                  placeholder="请选择编程语言"
                  style="width: 100%"
                  size="small"
                >
                  <el-option
                    v-for="lang in programmingLanguages"
                    :key="lang"
                    :label="lang"
                    :value="lang"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="健康检查" class="form-item-half">
                <el-input
                  v-model="applicationForm.health_api"
                  placeholder="/api/v1/health"
                  size="small"
                />
              </el-form-item>
            </div>

            <!-- 命令配置行 -->
            <div class="form-row">
              <el-form-item label="启动命令" class="form-item-half">
                <el-input
                  v-model="applicationForm.start_command"
                  placeholder="./app --config=config.yaml"
                  size="small"
                />
              </el-form-item>
              <el-form-item label="停止命令" class="form-item-half">
                <el-input
                  v-model="applicationForm.stop_command"
                  placeholder="kill -TERM $PID"
                  size="small"
                />
              </el-form-item>
            </div>

            <!-- 负责人配置 -->
            <div class="form-row">
              <el-form-item label="研发负责人" class="form-item-half">
                <el-select
                  v-model="applicationForm.dev_owners"
                  placeholder="请选择研发负责人"
                  style="width: 100%"
                  size="small"
                  multiple
                  collapse-tags
                  collapse-tags-tooltip
                  filterable
                  clearable
                >
                  <el-option
                    v-for="user in userOptions"
                    :key="user.id"
                    :label="user.name"
                    :value="user.id"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="测试负责人" class="form-item-half">
                <el-select
                  v-model="applicationForm.test_owners"
                  placeholder="请选择测试负责人"
                  style="width: 100%"
                  size="small"
                  multiple
                  collapse-tags
                  collapse-tags-tooltip
                  filterable
                  clearable
                >
                  <el-option
                    v-for="user in userOptions"
                    :key="user.id"
                    :label="user.name"
                    :value="user.id"
                  />
                </el-select>
              </el-form-item>
            </div>

            <el-form-item label="运维负责人" class="form-item-full">
              <el-select
                v-model="applicationForm.ops_owners"
                placeholder="请选择运维负责人"
                style="width: 100%"
                size="small"
                multiple
                collapse-tags
                collapse-tags-tooltip
                filterable
                clearable
              >
                <el-option
                  v-for="user in userOptions"
                  :key="user.id"
                  :label="user.name || `用户${user.id}`"
                  :value="user.id"
                />
              </el-select>
            </el-form-item>

            <!-- 关联资源配置 -->
            <div class="resource-section-divider">
              <span class="section-title-small">关联资源</span>
            </div>

            <el-form-item label="关联域名" class="form-item-full">
              <el-input
                v-model="applicationForm.domains"
                placeholder="多个域名用逗号分隔"
                size="small"
              />
            </el-form-item>

            <el-form-item label="关联主机" class="form-item-full">
              <el-select
                v-model="applicationForm.hosts"
                placeholder="请选择关联主机"
                style="width: 100%"
                size="small"
                multiple
                collapse-tags
                collapse-tags-tooltip
                filterable
                remote
                :remote-method="remoteSearchHosts"
                clearable
              >
                <el-option
                  v-for="host in hostOptions"
                  :key="host.id"
                  :label="`${host.name} (${host.ip})`"
                  :value="host.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="关联数据库" class="form-item-full">
              <el-select
                v-model="applicationForm.databases"
                placeholder="请选择关联数据库"
                style="width: 100%"
                size="small"
                multiple
                collapse-tags
                collapse-tags-tooltip
                filterable
                clearable
              >
                <el-option
                  v-for="db in databaseOptions"
                  :key="db.id"
                  :label="`${db.name} (${db.type})`"
                  :value="db.id"
                />
              </el-select>
            </el-form-item>

            <!-- 其他资源配置 -->
            <div class="resource-section-divider">
              <span class="section-title-small">其他资源</span>
            </div>

            <div class="form-row">
              <el-form-item label="Kafka" class="form-item-half">
                <el-input
                  v-model="applicationForm.other_res.kafka"
                  placeholder="Kafka实例名,用逗号分隔"
                  size="small"
                />
              </el-form-item>
              <el-form-item label="RabbitMQ" class="form-item-half">
                <el-input
                  v-model="applicationForm.other_res.rabbitmq"
                  placeholder="RabbitMQ实例名,用逗号分隔"
                  size="small"
                />
              </el-form-item>
            </div>

            <div class="form-row">
              <el-form-item label="kubernetes" class="form-item-half">
                <el-input
                  v-model="applicationForm.other_res.zookeeper"
                  placeholder="kubernetes,用逗号分隔"
                  size="small"
                />
              </el-form-item>
              <el-form-item label="其他资源" class="form-item-half">
                <el-input
                  v-model="applicationForm.other_res.other"
                  placeholder="其他资源,用逗号分隔"
                  size="small"
                />
              </el-form-item>
            </div>
          </div>
        </el-collapse-transition>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="createDialogVisible = false" size="small">取消</el-button>
          <el-button type="primary" @click="handleSaveApp" size="small">确认创建</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 编辑应用对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑应用"
      width="700px"
      class="edit-dialog compact-dialog"
    >
      <el-form :model="applicationForm" label-width="90px" class="compact-form">
        <!-- 必填字段 -->
        <el-form-item label="应用名称" required class="form-item-full">
          <el-input
            v-model="applicationForm.name"
            placeholder="请输入应用名称"
            size="small"
          />
        </el-form-item>

        <!-- 业务信息行 -->
        <div class="form-row">
          <el-form-item label="业务组" required class="form-item-half">
            <el-cascader
              v-model="applicationForm.business_group_id"
              :options="cascaderBusinessGroups"
              placeholder="请选择业务组"
              style="width: 100%"
              size="small"
              :props="{
                expandTrigger: 'hover',
                checkStrictly: true,
                emitPath: false,
                leaf: 'leaf'
              }"
            />
          </el-form-item>
          <el-form-item label="业务部门" required class="form-item-half">
            <el-select
              v-model="applicationForm.business_dept_id"
              placeholder="请选择业务部门"
              style="width: 100%"
              size="small"
              filterable
              clearable
            >
              <el-option
                v-for="dept in businessDepts"
                :key="dept.id"
                :label="dept.deptName"
                :value="dept.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="仓库地址" required class="form-item-full">
          <el-input
            v-model="applicationForm.repo_url"
            placeholder="https://git.company.com/project.git"
            size="small"
          />
        </el-form-item>

        <!-- 可选字段分隔线 -->
        <div class="optional-divider">
          <div class="divider-line"></div>
          <el-button
            text
            size="small"
            class="toggle-more-btn"
            @click="showMoreConfig = !showMoreConfig"
          >
            <el-icon>
              <ArrowUp v-if="showMoreConfig" />
              <ArrowDown v-else />
            </el-icon>
            {{ showMoreConfig ? '收起更多配置' : '展开更多配置' }}
          </el-button>
          <div class="divider-line"></div>
        </div>

        <!-- 可选字段（折叠区域） -->
        <el-collapse-transition>
          <div v-show="showMoreConfig" class="more-config-section">
            <!-- 应用描述 -->
            <el-form-item label="应用描述" class="form-item-full">
              <el-input
                v-model="applicationForm.description"
                type="textarea"
                :rows="2"
                placeholder="请输入应用描述"
                size="small"
              />
            </el-form-item>

            <!-- 技术配置行 -->
            <div class="form-row">
              <el-form-item label="编程语言" class="form-item-half">
                <el-select
                  v-model="applicationForm.programming_lang"
                  placeholder="请选择编程语言"
                  style="width: 100%"
                  size="small"
                >
                  <el-option
                    v-for="lang in programmingLanguages"
                    :key="lang"
                    :label="lang"
                    :value="lang"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="健康检查" class="form-item-half">
                <el-input
                  v-model="applicationForm.health_api"
                  placeholder="/api/v1/health"
                  size="small"
                />
              </el-form-item>
            </div>

            <!-- 命令配置行 -->
            <div class="form-row">
              <el-form-item label="启动命令" class="form-item-half">
                <el-input
                  v-model="applicationForm.start_command"
                  placeholder="./app --config=config.yaml"
                  size="small"
                />
              </el-form-item>
              <el-form-item label="停止命令" class="form-item-half">
                <el-input
                  v-model="applicationForm.stop_command"
                  placeholder="kill -TERM $PID"
                  size="small"
                />
              </el-form-item>
            </div>

            <!-- 负责人配置 -->
            <div class="form-row">
              <el-form-item label="研发负责人" class="form-item-half">
                <el-select
                  v-model="applicationForm.dev_owners"
                  placeholder="请选择研发负责人"
                  style="width: 100%"
                  size="small"
                  multiple
                  collapse-tags
                  collapse-tags-tooltip
                  filterable
                  clearable
                >
                  <el-option
                    v-for="user in userOptions"
                    :key="user.id"
                    :label="user.name"
                    :value="user.id"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="测试负责人" class="form-item-half">
                <el-select
                  v-model="applicationForm.test_owners"
                  placeholder="请选择测试负责人"
                  style="width: 100%"
                  size="small"
                  multiple
                  collapse-tags
                  collapse-tags-tooltip
                  filterable
                  clearable
                >
                  <el-option
                    v-for="user in userOptions"
                    :key="user.id"
                    :label="user.name"
                    :value="user.id"
                  />
                </el-select>
              </el-form-item>
            </div>

            <el-form-item label="运维负责人" class="form-item-full">
              <el-select
                v-model="applicationForm.ops_owners"
                placeholder="请选择运维负责人"
                style="width: 100%"
                size="small"
                multiple
                collapse-tags
                collapse-tags-tooltip
                filterable
                clearable
              >
                <el-option
                  v-for="user in userOptions"
                  :key="user.id"
                  :label="user.name"
                  :value="user.id"
                />
              </el-select>
            </el-form-item>

            <!-- 关联资源配置 -->
            <div class="resource-section-divider">
              <span class="section-title-small">关联资源</span>
            </div>

            <el-form-item label="关联域名" class="form-item-full">
              <el-input
                v-model="applicationForm.domains"
                placeholder="多个域名用逗号分隔"
                size="small"
              />
            </el-form-item>

            <el-form-item label="关联主机" class="form-item-full">
              <el-select
                v-model="applicationForm.hosts"
                placeholder="请选择关联主机"
                style="width: 100%"
                size="small"
                multiple
                collapse-tags
                collapse-tags-tooltip
                filterable
                remote
                :remote-method="remoteSearchHosts"
                clearable
              >
                <el-option
                  v-for="host in hostOptions"
                  :key="host.id"
                  :label="`${host.name} (${host.ip})`"
                  :value="host.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="关联数据库" class="form-item-full">
              <el-select
                v-model="applicationForm.databases"
                placeholder="请选择关联数据库"
                style="width: 100%"
                size="small"
                multiple
                collapse-tags
                collapse-tags-tooltip
                filterable
                clearable
              >
                <el-option
                  v-for="db in databaseOptions"
                  :key="db.id"
                  :label="`${db.name} (${db.type})`"
                  :value="db.id"
                />
              </el-select>
            </el-form-item>

            <!-- 其他资源配置 -->
            <div class="resource-section-divider">
              <span class="section-title-small">其他资源</span>
            </div>

            <div class="form-row">
              <el-form-item label="Kafka" class="form-item-half">
                <el-input
                  v-model="applicationForm.other_res.kafka"
                  placeholder="Kafka实例名,用逗号分隔"
                  size="small"
                />
              </el-form-item>
              <el-form-item label="RabbitMQ" class="form-item-half">
                <el-input
                  v-model="applicationForm.other_res.rabbitmq"
                  placeholder="RabbitMQ实例名,用逗号分隔"
                  size="small"
                />
              </el-form-item>
            </div>

            <div class="form-row">
              <el-form-item label="kubernetes" class="form-item-half">
                <el-input
                  v-model="applicationForm.other_res.zookeeper"
                  placeholder="kubernetes,用逗号分隔"
                  size="small"
                />
              </el-form-item>
              <el-form-item label="其他资源" class="form-item-half">
                <el-input
                  v-model="applicationForm.other_res.other"
                  placeholder="其他资源,用逗号分隔"
                  size="small"
                />
              </el-form-item>
            </div>
          </div>
        </el-collapse-transition>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="editDialogVisible = false" size="small">取消</el-button>
          <el-button type="primary" @click="handleSaveApp" size="small">确认更新</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 应用详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="应用详情"
      width="600px"
      class="detail-dialog compact-dialog"
    >
      <div class="detail-content">
        <div class="detail-item">
          <span class="detail-label">应用名称:</span>
          <span class="detail-value">{{ currentApplication.name }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">应用编码:</span>
          <span class="detail-value">{{ currentApplication.code }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">应用描述:</span>
          <span class="detail-value">{{ currentApplication.description || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">技术栈:</span>
          <span class="detail-value">{{ currentApplication.programming_lang || '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">状态:</span>
          <el-tag :type="getStatusTagType(currentApplication.status)" size="small">
            {{ getStatusText(currentApplication.status) }}
          </el-tag>
        </div>
        <div class="detail-item">
          <span class="detail-label">创建时间:</span>
          <span class="detail-value">
            {{ currentApplication.created_at ? new Date(currentApplication.created_at).toLocaleString() : '-' }}
          </span>
        </div>
      </div>
    </el-dialog>


    <!-- Jenkins环境配置对话框 -->
    <el-dialog
      v-model="jenkinsEnvDialogVisible"
      :title="`${currentAppName} - Jenkins环境配置`"
      width="900px"
      class="jenkins-env-dialog modern-dialog"
    >
      <el-card shadow="hover" class="env-main-card">
        <!-- 环境列表 -->
        <div class="env-list-section">
          <div class="section-header">
            <div class="header-left">
              <el-icon class="section-icon"><Monitor /></el-icon>
              <span class="section-title">环境列表</span>
            </div>
            <el-button v-authority="['app:application:envadd']" type="primary" size="small" @click="handleAddJenkinsEnv" class="add-btn">
              <el-icon><Plus /></el-icon>
              新增环境
            </el-button>
          </div>

          <el-table :data="jenkinsEnvs" size="small" class="env-table">
            <el-table-column prop="env_name" label="环境名称" width="120">
              <template #default="{ row }">
                <el-tag type="info" size="small" class="env-tag">{{ row.env_name }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="jenkins_server_id" label="Jenkins服务器" width="150">
              <template #default="{ row }">
                <span class="server-name">{{ getServerNameById(row.jenkins_server_id) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="job_name" label="任务名称" min-width="150">
              <template #default="{ row }">
                <span class="job-name">{{ row.job_name || '未配置' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="配置状态" width="100">
              <template #default="{ row }">
                <el-tag
                  :type="(row.jenkins_server_id && row.job_name) ? 'success' : 'warning'"
                  size="small"
                  class="status-tag"
                >
                  <el-icon>
                    <component :is="(row.jenkins_server_id && row.job_name) ? SuccessFilled : WarningFilled" />
                  </el-icon>
                  {{ (row.jenkins_server_id && row.job_name) ? '已配置' : '待配置' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="创建时间" width="160">
              <template #default="{ row }">
                <span class="create-time">{{ row.created_at ? new Date(row.created_at).toLocaleString() : '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180">
              <template #default="{ row }">
                <div class="table-actions">
                  <el-tooltip content="编辑环境" placement="top">
                    <el-button v-authority="['app:application:envedit']" size="small" type="primary" @click="handleEditJenkinsEnv(row)">
                      <el-icon><Edit /></el-icon>编辑
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="删除环境" placement="top">
                    <el-button v-authority="['app:application:envdelete']"
 size="small" type="danger" @click="handleDeleteJenkinsEnv(row)">
                      <el-icon><Delete /></el-icon>删除
                    </el-button>
                  </el-tooltip>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 环境表单 -->
        <div class="env-form-section" v-if="showJenkinsEnvForm || editingJenkinsEnv !== null">
          <el-card shadow="hover" class="form-card">
            <div class="section-header">
              <div class="header-left">
                <el-icon class="section-icon"><Setting /></el-icon>
                <span class="section-title">
                  {{ editingJenkinsEnv ? '编辑环境' : '新增环境' }}
                </span>
              </div>
            </div>

            <el-form :model="jenkinsEnvForm" label-width="120px" size="small" class="modern-form">
              <div class="form-row">
                <el-form-item label="环境名称" required class="form-item-half">
                  <el-input
                    v-model="jenkinsEnvForm.env_name"
                    placeholder="请输入环境名称（如：prod、test、dev）"
                  />
                </el-form-item>
                <el-form-item label="Jenkins服务器" class="form-item-half">
                  <el-select
                    v-model="jenkinsEnvForm.jenkins_server_id"
                    placeholder="请选择Jenkins服务器"
                    style="width: 100%"
                    filterable
                    clearable
                    @change="handleServerChange"
                  >
                    <el-option
                      v-for="server in jenkinsServers"
                      :key="server.id"
                      :label="server.name"
                      :value="server.id"
                    />
                  </el-select>
                </el-form-item>
              </div>
              <el-form-item label="任务名称" class="form-item-full">
                <div class="job-name-input-group">
                  <el-select
                    v-model="jenkinsEnvForm.job_name"
                    placeholder="请选择或输入Jenkins任务名称"
                    style="width: 100%"
                    filterable
                    allow-create
                    clearable
                    :loading="jobSearchLoading"
                    @input="handleJobNameInput"
                    :disabled="!jenkinsEnvForm.jenkins_server_id"
                  >
                    <el-option
                      v-for="job in jenkinsJobs"
                      :key="job"
                      :label="job"
                      :value="job"
                    />
                  </el-select>
                  <el-button
                    type="primary"
                    :icon="Search"
                    :loading="jobValidationLoading"
                    @click="validateJenkinsJob"
                    :disabled="!jenkinsEnvForm.jenkins_server_id || !jenkinsEnvForm.job_name"
                    class="validate-btn"
                  >
                    验证任务
                  </el-button>
                </div>

                <!-- 验证状态显示 -->
                <div v-if="jobValidationStatus" class="validation-status">
                  <div class="validation-result">
                    <el-icon v-if="jobValidationStatus.exists" class="success-icon">
                      <SuccessFilled />
                    </el-icon>
                    <el-icon v-else class="error-icon">
                      <WarningFilled />
                    </el-icon>
                    <span v-if="!jobValidationStatus.exists" class="error-message">
                      {{ jobValidationStatus.message }}
                    </span>
                  </div>
                </div>
              </el-form-item>
              <el-form-item class="form-actions">
                <el-button type="primary" @click="handleSaveJenkinsEnv">
                  <el-icon><component :is="editingJenkinsEnv ? Edit : Plus" /></el-icon>
                  {{ editingJenkinsEnv ? '更新环境' : '创建环境' }}
                </el-button>
                <el-button @click="resetJenkinsEnvForm">取消</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </div>
      </el-card>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="jenkinsEnvDialogVisible = false" size="small">
            <el-icon><Close /></el-icon>
            关闭
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.application-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

/* 主卡片样式 */
.main-card {
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
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-actions {
  display: flex;
  gap: 12px;
}

/* 搜索卡片样式 */
.search-card {
  background: rgba(103, 126, 234, 0.05);
  backdrop-filter: blur(5px);
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(103, 126, 234, 0.1);
  margin-bottom: 20px;
}

.search-card :deep(.el-card__body) {
  padding: 20px;
}

/* 表格卡片样式 */
.table-card {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(5px);
  border-radius: 12px;
  border: none;
  box-shadow: none;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

/* 搜索工具栏样式 */
.search-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  display: flex;
  align-items: center;
  gap: 12px;
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

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.resource-count {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
  white-space: nowrap;
}

/* 表格样式 */
.app-table {
  border-radius: 8px;
  overflow: hidden;
  border: none;
}

.app-table :deep(.el-table) {
  border: none;
}

.app-table :deep(.el-table__header) {
  background: #f8f9fa;
  border: none;
}

.app-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

.app-table :deep(.el-table__border) {
  border: none;
}

.app-table :deep(.el-table th) {
  border: none;
}

.app-table :deep(.el-table td) {
  border: none;
}

.app-table :deep(.el-table__body-wrapper) {
  border: none;
}

.app-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-icon {
  color: #409EFF;
  font-size: 16px;
}

.app-name-link {
  color: #409EFF;
  cursor: pointer;
  font-weight: 500;
  transition: color 0.3s;
}

.app-name-link:hover {
  color: #66b1ff;
  text-decoration: underline;
}

.operation-buttons {
  display: flex;
  gap: 4px;
  justify-content: center;
  align-items: center;
  flex-wrap: nowrap;
}

.operation-buttons .el-button {
  padding: 6px;
  min-width: 32px;
  height: 32px;
}

/* 分页样式 */
.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  padding: 20px;
  background: rgba(248, 249, 250, 0.8);
  border-radius: 0 0 12px 12px;
}

/* 对话框样式 */
.compact-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.compact-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.compact-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.compact-dialog :deep(.el-dialog__body) {
  padding: 15px 20px;
}

.compact-form {
  margin: 0;
}

.compact-form :deep(.el-form-item) {
  margin-bottom: 12px;
}

.form-row {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.form-item-half {
  flex: 1;
  margin-bottom: 0 !important;
}

.form-item-full {
  width: 100%;
  margin-bottom: 0 !important;
}

.dialog-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

/* 详情对话框样式 */
.detail-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.detail-label {
  font-weight: 600;
  color: #606266;
  min-width: 80px;
  flex-shrink: 0;
}

.detail-value {
  color: #303133;
  word-break: break-all;
  flex: 1;
}

/* 通用样式 */
.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.el-input :deep(.el-input__wrapper),
.el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__wrapper):hover,
.el-select :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.el-input :deep(.el-input__wrapper.is-focus),
.el-select :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #2c3e50;
}

.el-textarea :deep(.el-textarea__inner) {
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.2);
  transition: all 0.3s ease;
}

.el-textarea :deep(.el-textarea__inner):focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
}

.el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

/* 可选字段分隔线样式 */
.optional-divider {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 20px 0;
}

.divider-line {
  flex: 1;
  height: 1px;
  background: repeating-linear-gradient(
    to right,
    #dcdfe6 0,
    #dcdfe6 4px,
    transparent 4px,
    transparent 8px
  );
}

.toggle-more-btn {
  color: #409EFF;
  font-weight: 500;
  padding: 8px 16px;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.toggle-more-btn:hover {
  background: rgba(64, 158, 255, 0.1);
  border-radius: 8px;
}

.toggle-more-btn .el-icon {
  transition: transform 0.3s ease;
}

/* 更多配置区域样式 */
.more-config-section {
  padding: 16px 0 0 0;
  border-radius: 8px;
}

/* 资源配置区域分隔线样式 */
.resource-section-divider {
  display: flex;
  align-items: center;
  margin: 24px 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 1px solid #e4e7ed;
}

.section-title-small {
  font-size: 14px;
  font-weight: 600;
  color: #409eff;
  padding-right: 12px;
  position: relative;
}

.section-title-small::before {
  content: '';
  position: absolute;
  left: -8px;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 14px;
  background: linear-gradient(135deg, #409eff, #66b1ff);
  border-radius: 2px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .application-management {
    padding: 12px;
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .header-right {
    justify-content: flex-end;
  }

  .search-toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .toolbar-right {
    justify-content: center;
  }

  .search-form {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }

  .compact-dialog {
    width: 95% !important;
  }

  .form-row {
    flex-direction: column;
    gap: 8px;
  }

  .detail-content {
    grid-template-columns: 1fr;
  }

  .operation-buttons {
    flex-direction: column;
    gap: 4px;
  }

  .app-table :deep(.el-table) {
    font-size: 12px;
  }

  .pagination-wrapper {
    padding: 16px 12px;
  }
}

@media (max-width: 480px) {
  .app-table :deep(.el-table__cell) {
    padding: 8px 4px;
  }

  .operation-buttons .el-button {
    padding: 4px;
    width: 28px;
    height: 28px;
  }
}

/* Jenkins环境配置对话框样式 */
.jenkins-env-dialog.modern-dialog :deep(.el-dialog__body) {
  padding: 0;
}

.jenkins-env-dialog .env-main-card {
  border: none;
  box-shadow: none;
  background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
}

.jenkins-env-dialog .env-list-section {
  margin-bottom: 24px;
}

.jenkins-env-dialog .env-form-section {
  margin-top: 24px;
}

.jenkins-env-dialog .form-card {
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border: 1px solid #e3e8f0;
}

.jenkins-env-dialog .section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 2px solid #e3e8f0;
}

.jenkins-env-dialog .header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.jenkins-env-dialog .section-icon {
  color: #409EFF;
  font-size: 18px;
}

.jenkins-env-dialog .section-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.jenkins-env-dialog .add-btn {
  border-radius: 20px;
  padding: 8px 16px;
  font-weight: 500;
}

.jenkins-env-dialog .env-table {
  border: 1px solid #e3e8f0;
  border-radius: 12px;
  overflow: hidden;
}

.jenkins-env-dialog .env-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #409EFF 0%, #66b1ff 100%);
}

.jenkins-env-dialog .env-table :deep(.el-table__header th) {
  background: transparent;
  color: #2c3e50;
  font-weight: 600;
  border: none;
}

.jenkins-env-dialog .env-table :deep(.el-table__body tr:hover) {
  background-color: rgba(64, 158, 255, 0.1);
}

.jenkins-env-dialog .env-tag {
  background: linear-gradient(135deg, #409EFF 0%, #66b1ff 100%);
  color: white;
  border: none;
  font-weight: 500;
}

.jenkins-env-dialog .status-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
}

.jenkins-env-dialog .server-id,
.jenkins-env-dialog .job-name,
.jenkins-env-dialog .create-time {
  color: #606266;
  font-size: 13px;
}

.jenkins-env-dialog .table-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.jenkins-env-dialog .modern-form .form-row {
  display: flex;
  gap: 20px;
}

.jenkins-env-dialog .modern-form .form-item-half {
  flex: 1;
}

.jenkins-env-dialog .modern-form .form-item-full {
  width: 100%;
}

.jenkins-env-dialog .modern-form .form-actions {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #e3e8f0;
}

.jenkins-env-dialog .modern-form .el-form-item {
  margin-bottom: 20px;
}

.jenkins-env-dialog .dialog-footer {
  text-align: center;
  padding-top: 16px;
  border-top: 1px solid #e3e8f0;
}

/* 任务名称输入组样式 */
.jenkins-env-dialog .job-name-input-group {
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.jenkins-env-dialog .job-name-input-group .el-select {
  flex: 1;
}

.jenkins-env-dialog .validate-btn {
  flex-shrink: 0;
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 500;
}

.jenkins-env-dialog .validation-status {
  margin-top: 12px;
}

.jenkins-env-dialog .validation-result {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  background-color: rgba(0, 0, 0, 0.02);
}

.jenkins-env-dialog .success-icon {
  color: #67c23a;
  font-size: 18px;
}

.jenkins-env-dialog .error-icon {
  color: #e6a23c;
  font-size: 18px;
}

.jenkins-env-dialog .error-message {
  color: #e6a23c;
  font-size: 14px;
  font-weight: 500;
}

.jenkins-env-dialog .server-name {
  color: #409EFF;
  font-weight: 500;
}
</style>
