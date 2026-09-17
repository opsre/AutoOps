<script setup>
import { ref, reactive, onMounted, computed, defineOptions, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'

defineOptions({
  name: 'ShRelease'
})
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Delete,
  View,
  Search,
  Refresh,
  Document,
  CircleCheck,
  CircleClose,
  Warning,
  Folder,
  FolderOpened,
  Setting,
  Check
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

// 已配置 prod 环境的应用列表
const prodApps = ref([])

// 创建工单对话框
const createDialogVisible = ref(false)

// 表单数据
const releaseForm = ref({
  title: '',
  reason: '',
  business_group_id: null,
  app_id: null,
  app_name: '',
  repo_url: '', // 仓库地址
  server_host_id: null, // ⭐️ 执行服务器主机ID（单选）
  server_ip: '', // 服务器IP
  approver_id: null,
  executor_id: null,
  execute_dir: '/home/dingding/saas3-data-export/',
  script_content: '',
  jenkins_server_id: null,
  job_name: '',
  parameters: null
})

// 应用选择弹框相关
const appSelectionDialogVisible = ref(false)
const businessTreeData = ref([]) // 业务树数据
const currentBusinessGroup = ref(null) // 当前选中的业务组
const businessGroupApps = ref([]) // 当前业务组的应用列表
const loadingApps = ref(false)
const appTable = ref() // 应用表格引用

// Jenkins任务参数配置
const jenkinsParamsDialogVisible = ref(false)
const jenkinsParamsLoading = ref(false)
const jenkinsJobParameters = ref([])
const jenkinsParamValues = ref({})

// 服务器选择弹框相关
const serverSelectionDialogVisible = ref(false)
const cmdbGroupsWithHosts = ref([]) // CMDB分组和主机树数据
const currentGroupHosts = ref([]) // 当前分组的主机列表
const selectedServer = ref(null) // ⭐️ 已选择的服务器（单选）
const loadingServers = ref(false)
const serverPagination = reactive({
  pageNum: 1,
  pageSize: 10,
  total: 0
})
const currentGroupId = ref(null)

// 日志查看对话框
const logDialogVisible = ref(false)
const logContent = ref('')
const currentLogReleaseId = ref(null)
let logTimer = null

// 状态选项
const statusOptions = [
  { value: '', label: '全部' },
  { value: 1, label: '流程进行中' },
  { value: 2, label: '执行成功' },
  { value: 3, label: '执行失败' },
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
    2: '执行成功',
    3: '执行失败',
    4: '已驳回',
    5: '已取消'
  }
  return map[status] || '未知'
}

// 获取审核状态标签类型
const getApprovalStatusType = (status) => {
  const map = {
    1: 'info',
    2: 'success',
    3: 'danger'
  }
  return map[status] || 'info'
}

// 获取审核状态文本
const getApprovalStatusText = (status) => {
  const map = {
    1: '未审核',
    2: '审核通过',
    3: '审核驳回'
  }
  return map[status] || '未知'
}

// 获取执行状态类型
const getExecuteStatusType = (status) => {
  const map = {
    1: 'info',        // 未执行
    2: 'warning',     // 执行中（Jenkins拉取代码）
    3: 'success',     // 代码拉取成功
    4: 'danger',      // Jenkins执行失败
    5: 'warning',     // 脚本执行中
    6: 'success',     // 脚本执行成功
    7: 'danger'       // 脚本执行失败
  }
  return map[status] || 'info'
}

// 获取执行状态文本
const getExecuteStatusText = (status) => {
  const map = {
    1: '未执行',
    2: '代码拉取中',
    3: '代码拉取成功',
    4: 'Jenkins执行失败',
    5: '脚本执行中',
    6: '脚本执行成功',
    7: '脚本执行失败'
  }
  return map[status] || '未知'
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

    const response = await appApi.getShReleaseList(params)
    if (response.data?.code === 200) {
      const list = response.data.data?.list || []

      // 处理每条记录
      releaseList.value = list.map(item => {
        // 如果后端没有返回business_group_name，从本地业务组数据中查找
        if (!item.business_group_name && item.business_group_id) {
          item.business_group_name = getGroupName(item.business_group_id)
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

// 查看详情
const handleViewDetail = (row) => {
  router.push(`/app/sh-release/${row.id}`)
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

    const response = await appApi.deleteShRelease(row.id)
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
    reason: '',
    business_group_id: null,
    app_id: null,
    app_name: '',
    repo_url: '',
    server_host_id: null,
    server_ip: '',
    approver_id: null,
    executor_id: null,
    execute_dir: '/home/dingding/saas3-data-export/',
    script_content: '',
    jenkins_server_id: null,
    job_name: '',
    parameters: null
  }
  selectedServer.value = null
}

// 提交创建
const handleSubmit = async () => {
  // 验证
  if (!releaseForm.value.title) {
    ElMessage.warning('请填写上线标题')
    return
  }
  if (!releaseForm.value.reason) {
    ElMessage.warning('请填写上线原因')
    return
  }
  if (!releaseForm.value.app_id) {
    ElMessage.warning('请选择应用')
    return
  }
  if (!selectedServer.value) {
    ElMessage.warning('请选择服务器')
    return
  }
  if (!releaseForm.value.approver_id) {
    ElMessage.warning('请选择审核人')
    return
  }
  if (!releaseForm.value.executor_id) {
    ElMessage.warning('请选择执行人')
    return
  }
  if (!releaseForm.value.execute_dir) {
    ElMessage.warning('请填写执行目录')
    return
  }
  if (!releaseForm.value.script_content) {
    ElMessage.warning('请填写脚本内容')
    return
  }

  try {
    // 构建提交数据
    const submitData = {
      title: releaseForm.value.title,
      reason: releaseForm.value.reason,
      business_group_id: releaseForm.value.business_group_id,
      app_id: releaseForm.value.app_id,
      repo_url: releaseForm.value.repo_url,
      server_host_id: releaseForm.value.server_host_id,
      server_ip: releaseForm.value.server_ip,
      approver_id: releaseForm.value.approver_id,
      executor_id: releaseForm.value.executor_id,
      execute_dir: releaseForm.value.execute_dir,
      script_content: releaseForm.value.script_content,
      jenkins_server_id: releaseForm.value.jenkins_server_id,
      job_name: releaseForm.value.job_name,
      environment: 'prod' // 固定为 prod 环境
    }

    // 如果配置了参数,将参数转换为字符串格式
    if (releaseForm.value.parameters) {
      submitData.parameters = {}
      Object.keys(releaseForm.value.parameters).forEach(key => {
        const value = releaseForm.value.parameters[key]
        // 处理不同类型的参数
        if (Array.isArray(value)) {
          // 多选参数,用逗号连接
          submitData.parameters[key] = value.join(',')
        } else if (typeof value === 'boolean') {
          // 布尔参数转换为字符串
          submitData.parameters[key] = value ? 'true' : 'false'
        } else {
          submitData.parameters[key] = value !== null && value !== undefined ? String(value) : ''
        }
      })
    }

    const response = await appApi.createShRelease(submitData)
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

// 审核
const handleApprove = async (row, action) => {
  try {
    const result = await ElMessageBox.prompt(
      action === 2 ? '请输入审核意见' : '请输入驳回原因',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /.+/,
        inputErrorMessage: '请输入审核意见'
      }
    )

    const response = await appApi.approveShRelease({
      release_id: row.id,
      action: action,
      remark: result.value
    })

    if (response.data?.code === 200) {
      ElMessage.success(action === 2 ? '审核通过' : '已驳回')
      fetchList()
    } else {
      ElMessage.error(response.data?.msg || '操作失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('审核失败:', error)
      ElMessage.error('操作失败')
    }
  }
}

// 执行脚本
const handleExecute = async (row) => {
  try {
    await ElMessageBox.confirm('确认执行该脚本？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const response = await appApi.executeShRelease({
      release_id: row.id
    })

    if (response.data?.code === 200) {
      ElMessage.success('脚本执行已启动')
      fetchList()

      // 自动打开日志查看
      setTimeout(() => handleViewLog(row.id), 1000)
    } else {
      ElMessage.error(response.data?.msg || '执行失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('执行失败:', error)
      ElMessage.error('执行失败')
    }
  }
}

// 查看日志
const handleViewLog = async (releaseId) => {
  logDialogVisible.value = true
  logContent.value = '正在加载日志...'
  currentLogReleaseId.value = releaseId

  let start = 0

  // 清除之前的定时器
  if (logTimer) {
    clearInterval(logTimer)
  }

  // 轮询获取日志
  logTimer = setInterval(async () => {
    try {
      const response = await appApi.getShReleaseLog(releaseId, start)

      if (response.data?.code === 200) {
        const data = response.data.data

        if (start === 0) {
          logContent.value = data.log_content || '暂无日志'
        } else {
          logContent.value += data.log_content || ''
        }

        start = parseInt(data.next_start) || start

        // 如果执行完成，停止轮询
        if (data.status === 3 || data.status === 4) {
          clearInterval(logTimer)
          fetchList()
        }
      }
    } catch (error) {
      console.error('获取日志失败:', error)
      clearInterval(logTimer)
    }
  }, 2000)
}

// 关闭日志对话框
const closeLogDialog = () => {
  logDialogVisible.value = false
  if (logTimer) {
    clearInterval(logTimer)
    logTimer = null
  }
  currentLogReleaseId.value = null
  logContent.value = ''
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

// 显示应用选择弹框
const showAppSelectionDialog = () => {
  appSelectionDialogVisible.value = true
  currentBusinessGroup.value = null
  businessGroupApps.value = []
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

      // 提取应用列表，只显示已配置 prod 环境的应用
      const apps = []
      serviceTreeData.forEach(businessLine => {
        businessLine.services?.forEach(service => {
          const hasProdEnv = service.jenkins_envs?.some(env => env.env_name === 'prod')
          if (hasProdEnv) {
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
          }
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

// 选择应用（单选）
const handleAppSelect = async (app) => {
  // 查找 prod 环境的配置
  const prodEnv = app.jenkins_envs.find(env => env.env_name === 'prod')

  if (!prodEnv) {
    ElMessage.warning('该应用未配置 prod 环境')
    return
  }

  releaseForm.value.app_id = app.id
  releaseForm.value.app_name = app.name
  releaseForm.value.business_group_id = app.business_group_id
  releaseForm.value.jenkins_server_id = prodEnv.jenkins_server_id
  releaseForm.value.job_name = prodEnv.job_name
  releaseForm.value.repo_url = app.repo_url || '' // 保存仓库地址

  ElMessage.success(`已选择应用：${app.name} (prod环境)`)
  appSelectionDialogVisible.value = false

  // 自动获取 Jenkins 参数
  if (prodEnv.jenkins_server_id && prodEnv.job_name) {
    await loadJenkinsParams(prodEnv.jenkins_server_id, prodEnv.job_name)
  }
}

// 加载 Jenkins 任务参数
const loadJenkinsParams = async (serverId, jobName) => {
  try {
    jenkinsParamsLoading.value = true
    const response = await appApi.getJenkinsJobParameters(serverId, jobName)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      jenkinsJobParameters.value = responseData.data?.parameters || []

      if (jenkinsJobParameters.value.length > 0) {
        // 初始化参数值
        jenkinsParamValues.value = {}
        jenkinsJobParameters.value.forEach(param => {
          // PT_CHECKBOX类型需要特殊处理
          if (param.type === 'PT_CHECKBOX') {
            jenkinsParamValues.value[param.name] = param.defaultValue
              ? (typeof param.defaultValue === 'string'
                  ? param.defaultValue.split(',').map(v => v.trim()).filter(v => v)
                  : [])
              : []
          } else {
            jenkinsParamValues.value[param.name] = param.defaultValue || ''
          }
        })

        // 自动打开参数配置对话框
        jenkinsParamsDialogVisible.value = true
      } else {
        ElMessage.info('该任务无需配置参数')
      }
    } else {
      ElMessage.error(responseData.message || '获取任务参数失败')
    }
  } catch (error) {
    console.error('获取Jenkins任务参数失败:', error)
    ElMessage.error('获取任务参数失败')
  } finally {
    jenkinsParamsLoading.value = false
  }
}

// 保存 Jenkins 参数配置
const saveJenkinsParams = () => {
  releaseForm.value.parameters = { ...jenkinsParamValues.value }
  ElMessage.success('参数配置已保存')
  closeParamsDialog()
}

// 关闭参数对话框
const closeParamsDialog = () => {
  jenkinsParamsDialogVisible.value = false
}

// 手动配置 Jenkins 参数
const handleConfigParams = async () => {
  if (!releaseForm.value.jenkins_server_id || !releaseForm.value.job_name) {
    ElMessage.warning('请先选择应用')
    return
  }

  await loadJenkinsParams(releaseForm.value.jenkins_server_id, releaseForm.value.job_name)
}

// 显示服务器选择弹框
const showServerSelectionDialog = () => {
  serverSelectionDialogVisible.value = true
  fetchCmdbGroupsWithHosts()
  getAllServers()
}

// 获取CMDB分组和主机树
const fetchCmdbGroupsWithHosts = async () => {
  try {
    const response = await cmdbApi.getGroupListWithHosts()
    if (response.data?.code === 200) {
      cmdbGroupsWithHosts.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取分组和主机列表失败:', error)
  }
}

// 获取所有服务器(默认显示)
const getAllServers = async () => {
  loadingServers.value = true
  try {
    const response = await cmdbApi.getCmdbHostList({
      page: serverPagination.pageNum,
      pageSize: serverPagination.pageSize
    })

    if (response.data?.code === 200) {
      currentGroupHosts.value = (response.data.data?.list || []).map(host => ({
        id: host.id,
        name: host.hostName || host.name || host.hostname,
        ip: host.privateIp || host.ip,
        publicIp: host.publicIp,
        os: host.os,
        remark: host.remark
      }))
      serverPagination.total = response.data.data?.total || 0
    }
  } catch (error) {
    console.error('获取服务器列表失败:', error)
    currentGroupHosts.value = []
  } finally {
    loadingServers.value = false
  }
}

// 处理分组树节点点击
const handleGroupClick = async (data) => {
  currentGroupId.value = data.id
  loadingServers.value = true
  try {
    const response = await cmdbApi.getCmdbHostsByGroupId(data.id, {
      page: serverPagination.pageNum,
      pageSize: serverPagination.pageSize
    })

    if (response.data?.code === 200) {
      const hostList = response.data.data || []
      currentGroupHosts.value = hostList.map(host => ({
        id: host.id,
        name: host.hostName || host.name || host.hostname,
        ip: host.privateIp || host.ip,
        publicIp: host.publicIp,
        os: host.os,
        remark: host.remark
      }))
      serverPagination.total = hostList.length
    }
  } catch (error) {
    console.error('获取主机列表失败:', error)
    currentGroupHosts.value = []
  } finally {
    loadingServers.value = false
  }
}

// 处理服务器选择变化（单选）
const handleServerRowClick = (row) => {
  selectedServer.value = row
  releaseForm.value.server_host_id = row.id
  releaseForm.value.server_ip = row.ip || '' // 保存服务器IP
}

// 确认选择服务器
const confirmServerSelection = () => {
  if (!selectedServer.value) {
    ElMessage.warning('请选择一台服务器')
    return
  }

  ElMessage.success(`已选择服务器: ${selectedServer.value.name || selectedServer.value.ip}`)
  serverSelectionDialogVisible.value = false
}

// 移除服务器
const removeServer = () => {
  selectedServer.value = null
  releaseForm.value.server_host_id = null
  ElMessage.success('已移除服务器')
}

// 分页变化
const handleServerPageChange = (page) => {
  serverPagination.pageNum = page
  if (currentGroupId.value) {
    handleGroupClick({ id: currentGroupId.value })
  } else {
    getAllServers()
  }
}

const handleServerSizeChange = (size) => {
  serverPagination.pageSize = size
  if (currentGroupId.value) {
    handleGroupClick({ id: currentGroupId.value })
  } else {
    getAllServers()
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

// 加载已配置 prod 环境的应用
const fetchProdApps = async () => {
  try {
    const response = await appApi.getServiceTree({})
    if (response.data?.code === 200) {
      const serviceTreeData = response.data.data || []

      // 提取配置了 prod 环境的应用列表
      const apps = []
      serviceTreeData.forEach(businessLine => {
        businessLine.services?.forEach(service => {
          const hasProdEnv = service.jenkins_envs?.some(env => env.env_name === 'prod')
          if (hasProdEnv) {
            apps.push({
              id: service.id,
              name: service.name,
              business_group_id: businessLine.business_group_id,
              business_group_name: businessLine.business_group_name
            })
          }
        })
      })

      prodApps.value = apps
    }
  } catch (error) {
    console.error('获取应用列表失败:', error)
  }
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

// 根据业务组过滤应用
const filteredProdApps = computed(() => {
  if (!releaseForm.value.business_group_id) {
    return prodApps.value
  }
  return prodApps.value.filter(app => app.business_group_id === releaseForm.value.business_group_id)
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

// 获取用户名
const getUserName = (userId) => {
  const user = userOptions.value.find(u => u.id === userId)
  return user ? user.name : '-'
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
  fetchProdApps()
})

// 组件销毁时清除定时器
onBeforeUnmount(() => {
  if (logTimer) {
    clearInterval(logTimer)
  }
})
</script>

<template>
  <div class="sh-release-page">
    <el-card shadow="hover" class="main-card">
      <template #header>
        <div class="card-header">
          <span class="title">脚本上线</span>
          <div class="header-actions">
            <el-button type="primary" v-authority="['app:service_sh:add']" :icon="Plus" size="small" @click="handleCreate">
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
          <el-table-column prop="title" label="上线标题" min-width="180" show-overflow-tooltip />
          <el-table-column prop="app_name" label="服务名称" width="150" show-overflow-tooltip align="center" />
          <el-table-column prop="approver_name" label="审核人" width="100" align="center" />
          <el-table-column prop="executor_name" label="执行人" width="100" align="center" />
          <el-table-column label="审核状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getApprovalStatusType(row.approval_status)" size="small">
                {{ getApprovalStatusText(row.approval_status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="执行状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getExecuteStatusType(row.execute_status)" size="small">
                {{ getExecuteStatusText(row.execute_status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="170" align="center">
            <template #default="{ row }">
              {{ formatTime(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="360" fixed="right" align="center">
            <template #default="{ row }">
              <el-tooltip content="查看详情">
                <el-button :icon="View" v-authority="['app:service_sh:get']" size="small" type="primary" @click="handleViewDetail(row)">详情</el-button>
              </el-tooltip>
              <el-tooltip content="审核通过" v-if="row.approval_status === 1">
                <el-button :icon="CircleCheck" size="small" type="success" @click="handleApprove(row, 2)">通过</el-button>
              </el-tooltip>
              <el-tooltip content="审核驳回" v-if="row.approval_status === 1">
                <el-button :icon="CircleClose" size="small" type="danger" @click="handleApprove(row, 3)">驳回</el-button>
              </el-tooltip>
              <el-tooltip content="执行脚本" v-if="row.approval_status === 2 && row.execute_status === 1">
                <el-button :icon="Warning" size="small" type="warning" @click="handleExecute(row)">执行</el-button>
              </el-tooltip>
              <el-tooltip content="删除">
                <el-button
                  :icon="Delete"
                  size="small"
                  v-authority="['app:service_sh:delete']"
                  type="danger"
                  @click="handleDelete(row)"
                  :disabled="row.execute_status === 2"
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
      title="新建脚本上线工单"
      width="700px"
      class="create-dialog"
    >
      <el-form :model="releaseForm" label-width="100px" size="small">
        <el-form-item label="上线标题" required>
          <el-input v-model="releaseForm.title" placeholder="请输入上线标题，描述本次脚本上线的目的" />
        </el-form-item>

        <el-form-item label="上线原因" required>
          <el-input
            v-model="releaseForm.reason"
            type="textarea"
            :rows="3"
            placeholder="详细说明为什么需要执行此脚本"
          />
        </el-form-item>

        <el-form-item label="应用配置" required>
          <div class="applications-section">
            <div class="applications-header">
              <el-button type="primary" size="small" @click="showAppSelectionDialog" :icon="Plus">
                选择应用
              </el-button>
              <span class="app-selected" v-if="releaseForm.app_name">
                已选择：<el-tag type="success" size="small">{{ releaseForm.app_name }}</el-tag>
                <el-tag type="danger" size="small" style="margin-left: 8px;">prod</el-tag>
                <el-button
                  v-if="releaseForm.jenkins_server_id && releaseForm.job_name"
                  size="small"
                  type="primary"
                  plain
                  style="margin-left: 12px;"
                  @click="handleConfigParams">
                  <el-icon><Setting /></el-icon>
                  {{ releaseForm.parameters ? '已配置参数' : '配置参数' }}
                </el-button>
              </span>
            </div>
          </div>
        </el-form-item>

        <div class="form-row">
          <el-form-item label="审核人" required class="form-item-half">
            <el-select v-model="releaseForm.approver_id" placeholder="选择审核人" filterable style="width: 100%">
              <el-option v-for="user in userOptions" :key="user.id" :label="user.name" :value="user.id" />
            </el-select>
          </el-form-item>
          <el-form-item label="执行人" required class="form-item-half">
            <el-select v-model="releaseForm.executor_id" placeholder="选择执行人" filterable style="width: 100%">
              <el-option v-for="user in userOptions" :key="user.id" :label="user.name" :value="user.id" />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="选择服务器" required>
          <div class="server-section">
            <div class="server-header">
              <el-button type="primary" size="small" @click="showServerSelectionDialog" :icon="Plus">
                选择主机
              </el-button>
            </div>
            <div v-if="selectedServer" class="selected-server-card">
              <div class="server-info">
                <div class="server-name">
                  <el-icon><Document /></el-icon>
                  <span>{{ selectedServer.name }}</span>
                </div>
                <div class="server-ip">
                  <el-tag type="primary" size="small">{{ selectedServer.ip }}</el-tag>
                  <span v-if="selectedServer.os" class="server-os">{{ selectedServer.os }}</span>
                </div>
              </div>
              <el-button
                type="danger"
                :icon="Delete"
                size="small"
                plain
                @click="removeServer"
              >移除</el-button>
            </div>
            <el-empty v-else description="请选择执行服务器" :image-size="60" />
          </div>
        </el-form-item>

        <el-form-item label="执行目录" required>
          <el-input
            v-model="releaseForm.execute_dir"
            placeholder="如：/home/dingding/saas3-data-export/"
          />
        </el-form-item>

        <el-form-item label="脚本内容" required>
          <el-input
            v-model="releaseForm.script_content"
            type="textarea"
            :rows="8"
            placeholder="如：node handle_pg/data/export_home_gateway_offline_data.js pro xxxx"
            class="script-textarea"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="createDialogVisible = false" size="small">取消</el-button>
        <el-button type="primary" @click="handleSubmit" size="small">提交工单</el-button>
      </template>
    </el-dialog>

    <!-- 日志查看对话框 -->
    <el-dialog
      v-model="logDialogVisible"
      title="执行日志"
      width="80%"
      @close="closeLogDialog"
      class="log-dialog"
    >
      <div class="log-container">
        <pre>{{ logContent }}</pre>
      </div>
      <template #footer>
        <el-button @click="closeLogDialog">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 应用选择弹框 -->
    <el-dialog
      v-model="appSelectionDialogVisible"
      title="选择应用"
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
                  <component :is="node.expanded ? FolderOpened : Folder" />
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
            <h4>应用列表</h4>
            <span v-if="currentBusinessGroup" class="selected-group">{{ currentBusinessGroup.name }}</span>
          </div>

          <el-table
            :data="businessGroupApps"
            border
            style="width: 100%"
            ref="appTable"
            v-loading="loadingApps">
            <el-table-column prop="name" label="应用名称" min-width="150">
              <template #default="{ row }">
                <div class="app-name-cell">
                  <el-icon><Document /></el-icon>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="code" label="应用代码" width="120">
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
            <el-table-column label="操作" width="100" align="center">
              <template #default="{ row }">
                <el-button
                  size="small"
                  type="primary"
                  @click="handleAppSelect(row)"
                  :disabled="row.status !== 2">
                  选择
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="empty-state" v-if="!currentBusinessGroup">
            <el-empty :image-size="100" description="请选择左侧业务组查看应用">
              <template #image>
                <el-icon size="100" color="#ccc"><Document /></el-icon>
              </template>
            </el-empty>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="appSelectionDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Jenkins任务参数配置对话框 -->
    <el-dialog
      v-model="jenkinsParamsDialogVisible"
      :title="`配置参数: ${releaseForm.app_name} - prod`"
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

    <!-- 服务器选择弹框 -->
    <el-dialog
      v-model="serverSelectionDialogVisible"
      title="选择执行主机"
      width="70%"
      class="modern-dialog server-selection-dialog">
      <div class="host-selection-layout-horizontal">
        <!-- 左侧分组树 -->
        <div class="tree-section-left">
          <div class="tree-header">
            <h4>主机分组</h4>
          </div>
          <el-tree
            :data="cmdbGroupsWithHosts"
            :props="{
              label: 'name',
              children: 'children'
            }"
            node-key="id"
            highlight-current
            @node-click="handleGroupClick"
          >
            <template #default="{ node, data }">
              <span class="tree-node">
                <el-icon v-if="node.expanded"><FolderOpened /></el-icon>
                <el-icon v-else><Folder /></el-icon>
                {{ node.label }}
              </span>
            </template>
          </el-tree>
        </div>

        <!-- 右侧主机表格 -->
        <div class="table-section-right">
          <el-table
            :data="currentGroupHosts"
            border
            style="width: 100%"
            highlight-current-row
            @row-click="handleServerRowClick"
            v-loading="loadingServers">
            <el-table-column width="55" align="center">
              <template #default="{ row }">
                <el-radio :model-value="selectedServer?.id" :label="row.id" @click.stop="handleServerRowClick(row)">
                  <span></span>
                </el-radio>
              </template>
            </el-table-column>
            <el-table-column prop="name" label="主机名称" min-width="150">
              <template #default="{ row }">
                <div class="hostname-cell">
                  <el-icon><Document /></el-icon>
                  <span>{{ row.name || '-' }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="ip" label="内网IP" width="140">
              <template #default="{ row }">
                <el-tag type="primary" size="small">{{ row.ip }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="publicIp" label="外网IP" width="140">
              <template #default="{ row }">
                <span>{{ row.publicIp || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="os" label="操作系统" width="180" />
            <el-table-column prop="remark" label="备注" show-overflow-tooltip />
          </el-table>

          <div class="pagination">
            <el-pagination
              @current-change="handleServerPageChange"
              @size-change="handleServerSizeChange"
              :current-page="serverPagination.pageNum"
              :page-sizes="[10, 20, 50, 100]"
              :page-size="serverPagination.pageSize"
              layout="total, sizes, prev, pager, next, jumper"
              :total="serverPagination.total"
            />
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="serverSelectionDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="confirmServerSelection">
          确认选择
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.sh-release-page {
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

/* 应用配置区域样式 */
.applications-section {
  width: 100%;
}

.applications-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.app-selected {
  color: #606266;
  font-size: 14px;
}

/* 服务器配置区域样式 */
.server-section {
  width: 100%;
}

.server-header {
  margin-bottom: 10px;
}

.selected-servers-list {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 10px;
  background: #f5f7fa;
  margin-top: 10px;
}

.selected-servers-title {
  color: #606266;
  margin-bottom: 8px;
  font-size: 14px;
}

.selected-servers-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.selected-server-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 10px;
  background: white;
  border-radius: 4px;
}

.selected-server-item span {
  font-size: 14px;
  color: #303133;
}

/* 服务器选择弹框样式 - 左右布局 */
.server-selection-dialog :deep(.el-dialog__body) {
  padding: 20px;
}

.host-selection-layout-horizontal {
  display: flex;
  height: 500px;
  gap: 16px;
}

.tree-section-left {
  width: 250px;
  border-right: 1px solid #ebeef5;
  padding-right: 10px;
  overflow-y: auto;
}

.tree-section-left .tree-header {
  margin-bottom: 12px;
}

.tree-section-left .tree-header h4 {
  margin: 0;
  color: #303133;
  font-size: 16px;
}

.table-section-right {
  flex: 1;
  padding-left: 10px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.table-section-right .pagination {
  margin-top: 15px;
  text-align: right;
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

/* 日志对话框样式 */
.log-dialog :deep(.el-dialog) {
  border-radius: 16px;
}

.log-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.log-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.log-container {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 15px;
  border-radius: 8px;
  max-height: 500px;
  overflow-y: auto;
}

.log-container pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 服务器相关样式增强 */
.server-section {
  width: 100%;
}

.server-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.server-hint {
  font-size: 12px;
  color: #909399;
}

.selected-server-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: linear-gradient(135deg, #f0f9ff 0%, #e8f5e9 100%);
  border-radius: 8px;
  border: 1px solid #67c23a;
  transition: all 0.3s;
}

.selected-server-card:hover {
  box-shadow: 0 4px 12px rgba(103, 194, 58, 0.2);
}

.server-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.server-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.server-name .el-icon {
  color: #67c23a;
  font-size: 18px;
}

.server-ip {
  display: flex;
  align-items: center;
  gap: 12px;
}

.server-os {
  font-size: 12px;
  color: #909399;
}

.server-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 脚本输入框样式 */
.script-textarea :deep(textarea) {
  font-family: Consolas, Monaco, 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  background: #1e1e1e;
  color: #d4d4d4;
  border: 1px solid #409eff;
  border-radius: 4px;
}

.script-textarea :deep(textarea)::placeholder {
  color: #909399;
  font-family: Consolas, Monaco, 'Courier New', monospace;
}

.script-textarea :deep(textarea):focus {
  background: #1e1e1e;
  color: #d4d4d4;
  border-color: #67c23a;
}
</style>
