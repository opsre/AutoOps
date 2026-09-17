<script setup>
import { ref, onMounted, onUnmounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Edit,
  Delete,
  Pointer,
  FolderOpened,
  Folder,
  DataBoard,
  MoreFilled,
  ArrowDown,
  DocumentAdd,
  ArrowRight,
  Warning
} from '@element-plus/icons-vue'
import { apiMonitorApi, businessGroupApi } from '@/api/monitor'
import healthIcon from '@/assets/image/健康.svg'
import fireIcon from '@/assets/image/起火.svg'

// ==================== 状态管理 ====================

// 左侧树形菜单数据
const treeData = ref([])
const selectedNodeData = ref(null)
const selectedFunctionId = ref(null)
const selectedGroupName = ref('')

// 二级分组列表（用于对话框选择）
const secondLevelGroups = ref([])

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: -1,
  pageNum: 1,
  pageSize: 10
})

// 列表数据
const tableData = ref([])
const total = ref(0)
const loading = ref(false)

// 新增/编辑对话框
const dialogVisible = ref(false)
const dialogTitle = ref('新增 API 监控')
const formLoading = ref(false)

// 表单数据（简化版 - 只有4个必填字段）
const monitorForm = reactive({
  id: null,
  name: '',
  url: '',
  method: 'GET',
  groupId: null,
  businessFunctionName: '',
  description: '',
  headers: '',
  // 子域名配置列表
  childUrls: []
})

// 子域名配置项结构
const createChildUrlItem = () => ({
  name: '',
  url: '',
  tempId: Date.now() + Math.random() // 临时ID用于v-for的key
})

const formRef = ref(null)

// 下拉选项
const statusOptions = [
  { label: '全部', value: -1 },
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' }
]

// 表单验证规则
const formRules = {
  name: [{ required: true, message: '请输入监控名称', trigger: 'blur' }],
  groupId: [{ required: true, message: '请选择业务组', trigger: 'change' }],
  url: [{ required: true, message: '请输入接口地址', trigger: 'blur' }],
  businessFunctionName: [{ required: true, message: '请输入业务功能名称', trigger: 'blur' }],
  headers: [
    {
      validator: (rule, value, callback) => {
        if (!value) {
          callback()
          return
        }
        try {
          JSON.parse(value)
          callback()
        } catch (e) {
          callback(new Error('Headers 必须是有效的 JSON 格式'))
        }
      },
      trigger: 'blur'
    }
  ]
}

// 测试按钮加载状态
const testLoading = ref(false)

// 自动刷新定时器
const autoRefreshTimer = ref(null)

// ==================== 分层监控相关状态 ====================

// 抽屉相关
const drawerVisible = ref(false)
const drawerMonitorDetail = ref(null) // 当前抽屉中显示的监控详情
const drawerLoading = ref(false)

// ==================== API 调用 ====================

// 获取业务组列表并构建树形结构（使用新的统计接口）
const fetchTreeData = async () => {
  try {
    const { data: res } = await businessGroupApi.getGroupTreeStats()
    if (res.code === 200) {
      const tree = []
      const groups = []

      // 遍历一级业务组
      for (const parentGroup of res.data) {
        // 一级分组节点
        const firstLevelNode = {
          id: `parent_${parentGroup.id}`,
          label: `${parentGroup.name} (${parentGroup.apiCount || 0})`,
          type: 'parent',
          parentId: parentGroup.id,
          apiCount: parentGroup.apiCount || 0,
          healthy: parentGroup.healthy || 0,
          unhealthy: parentGroup.unhealthy || 0,
          children: [],
          isLeaf: false
        }

        // 遍历二级业务组
        if (parentGroup.children && Array.isArray(parentGroup.children)) {
          for (const childGroup of parentGroup.children) {
            // 保存二级分组到列表（用于对话框）
            groups.push({
              id: childGroup.id,
              name: childGroup.name,
              parentName: parentGroup.name
            })

            // 二级分组节点
            const secondLevelNode = {
              id: `group_${childGroup.id}`,
              label: `${childGroup.name} (${childGroup.apiCount || 0})`,
              type: 'group',
              groupId: childGroup.id,
              parentId: parentGroup.id,
              apiCount: childGroup.apiCount || 0,
              healthy: childGroup.healthy || 0,
              unhealthy: childGroup.unhealthy || 0,
              businessFunctions: childGroup.businessFunctions || [],
              isLeaf: true
            }

            firstLevelNode.children.push(secondLevelNode)
          }
        }

        tree.push(firstLevelNode)
      }

      treeData.value = tree
      secondLevelGroups.value = groups
      console.log('树形数据构建完成:', tree)
    }
  } catch (error) {
    ElMessage.error('获取业务组树形统计失败')
    console.error(error)
  }
}

// 获取 API 监控列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = {
      ...searchForm
    }
    if (params.status === -1) delete params.status
    if (!params.keyword) delete params.keyword

    // 默认只显示顶级监控（父监控和独立监控）
    params.onlyTopLevel = true
    params.includeChildren = true // 包含子监控数据用于统计

    const { data: res } = await apiMonitorApi.getList(params)
    if (res.code === 200) {
      // 实际响应结构是 res.data.list 和 res.data.total
      tableData.value = Array.isArray(res.data.list) ? res.data.list : []
      total.value = res.data.total || 0
    }
  } catch (error) {
    ElMessage.error('获取 API 监控列表失败')
    tableData.value = []
    total.value = 0
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 获取指定业务功能的 API 列表
const fetchFunctionApiList = async (functionId) => {
  loading.value = true
  try {
    const { data: res } = await apiMonitorApi.getByBusinessFunction(functionId)
    if (res.code === 200) {
      let apiData = res.data
      if (apiData.list) {
        apiData = apiData.list
      }
      tableData.value = Array.isArray(apiData) ? apiData : []
      total.value = tableData.value.length
    }
  } catch (error) {
    ElMessage.error('获取 API 列表失败')
    tableData.value = []
    total.value = 0
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 获取指定业务组下所有API（直接按 groupId 查询）
const fetchGroupApiList = async (groupId) => {
  loading.value = true
  try {
    console.log('=== fetchGroupApiList 开始 ===')
    console.log('请求业务组ID:', groupId)

    // 直接按 groupId 查询该业务组下的所有API
    const params = {
      groupId: groupId,
      page: 1,
      pageSize: 1000, // 获取该业务组下的所有API
      onlyTopLevel: true, // 只显示顶级监控（父监控和独立监控）
      includeChildren: true // 包含子监控数据用于统计
    }

    const { data: res } = await apiMonitorApi.getList(params)
    console.log('API列表响应:', res)

    if (res.code === 200) {
      const apiList = res.data.list || []
      console.log('API数量:', apiList.length)
      console.log('API列表:', apiList)

      tableData.value = apiList
      total.value = res.data.total || apiList.length

      if (apiList.length === 0) {
        ElMessage.info('该业务组下暂无 API 监控配置')
      }
    } else {
      ElMessage.error(res.message || '获取API列表失败')
      tableData.value = []
      total.value = 0
    }
  } catch (error) {
    ElMessage.error('获取业务组 API 列表失败')
    tableData.value = []
    total.value = 0
    console.error('fetchGroupApiList 错误:', error)
  } finally {
    loading.value = false
    console.log('=== fetchGroupApiList 结束 ===')
  }
}

// ==================== 表单操作 ====================

// 搜索
const handleSearch = () => {
  if (selectedFunctionId.value) {
    // 如果选中了功能，使用功能的 API 列表进行过滤
    const filtered = tableData.value.filter(api =>
      api.name.toLowerCase().includes(searchForm.keyword.toLowerCase()) ||
      api.url.toLowerCase().includes(searchForm.keyword.toLowerCase())
    )
    tableData.value = filtered
  } else {
    // 否则获取全部 API 列表
    searchForm.pageNum = 1
    fetchList()
  }
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = -1
  searchForm.pageNum = 1
  if (selectedFunctionId.value) {
    fetchFunctionApiList(selectedFunctionId.value)
  } else {
    handleSearch()
  }
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
  dialogTitle.value = '新增 API 监控'
  monitorForm.id = null
  monitorForm.name = ''
  monitorForm.url = ''
  monitorForm.method = 'GET'
  monitorForm.groupId = null
  monitorForm.businessFunctionName = ''
  monitorForm.description = ''
  monitorForm.headers = ''
  monitorForm.childUrls = []

  dialogVisible.value = true
}

// 编辑
const handleEdit = async (row) => {
  try {
    const { data: res } = await apiMonitorApi.getDetail(row.id)
    if (res.code === 200) {
      dialogTitle.value = '编辑 API 监控'
      const detail = res.data
      monitorForm.id = detail.id
      monitorForm.name = detail.name || ''
      monitorForm.url = detail.url || ''
      monitorForm.method = detail.method || 'GET'
      monitorForm.groupId = Number(detail.groupId)
      monitorForm.businessFunctionName = detail.businessFunctionName || ''
      monitorForm.description = detail.description || ''
      monitorForm.headers = detail.headers || ''

      // 加载子监控配置
      if (detail.childMonitors && detail.childMonitors.length > 0) {
        monitorForm.childUrls = detail.childMonitors.map(child => ({
          id: child.id, // 保存已有子监控的ID
          name: child.name,
          url: child.url,
          tempId: Date.now() + Math.random()
        }))
      } else {
        monitorForm.childUrls = []
      }

      dialogVisible.value = true
    } else {
      ElMessage.error(res.message || '获取监控详情失败')
    }
  } catch (error) {
    ElMessage.error('获取监控详情失败')
    console.error(error)
  }
}

// 添加子域名配置项
const addChildUrl = () => {
  monitorForm.childUrls.push(createChildUrlItem())
}

// 删除子域名配置项
const removeChildUrl = async (index) => {
  const child = monitorForm.childUrls[index]

  // 如果是已保存的子监控（有id），需要先删除后端数据
  if (child.id) {
    try {
      await ElMessageBox.confirm(
        `确定要删除子监控"${child.name}"吗?`,
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )

      const { data: res } = await apiMonitorApi.delete(child.id)
      if (res.code === 200) {
        ElMessage.success('子监控删除成功')
        monitorForm.childUrls.splice(index, 1)
      } else {
        ElMessage.error(res.message || '删除失败')
      }
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error('删除失败')
        console.error(error)
      }
    }
  } else {
    // 新添加的子监控，直接从数组中删除
    monitorForm.childUrls.splice(index, 1)
  }
}

// 快速生成子域名URL
const quickGenerateChildUrl = (index, ipOrDomain) => {
  if (!ipOrDomain || !monitorForm.url) return

  try {
    const parentUrlObj = new URL(monitorForm.url)
    const parentPath = parentUrlObj.pathname + parentUrlObj.search

    // 判断是否为内网IP
    const isPrivateIP = /^(192\.168\.|10\.|172\.(1[6-9]|2[0-9]|3[01])\.)/.test(ipOrDomain)
    const protocol = isPrivateIP ? 'http' : parentUrlObj.protocol.replace(':', '')

    monitorForm.childUrls[index].url = `${protocol}://${ipOrDomain}${parentPath}`
  } catch (e) {
    ElMessage.error('URL 生成失败')
    console.error('URL解析失败', e)
  }
}

// 批量快速生成子域名
const batchGenerateChildUrls = async () => {
  try {
    const { value: ipsText } = await ElMessageBox.prompt(
      '请输入多个 IP 或域名（每行一个）',
      '批量生成子域名',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputType: 'textarea',
        inputPlaceholder: '例如:\n192.168.7.6:7001\n192.168.8.6:7001\n192.168.8.173:7001',
        inputValidator: (value) => {
          if (!value || !value.trim()) return '请输入至少一个 IP 或域名'
          return true
        }
      }
    )

    if (!ipsText || !monitorForm.url) return

    const ips = ipsText.trim().split('\n').filter(ip => ip.trim())

    ips.forEach((ip, idx) => {
      const trimmedIp = ip.trim()
      const childItem = createChildUrlItem()
      childItem.name = `${monitorForm.name}-${String(idx + 1).padStart(2, '0')}` // 自动生成名称

      try {
        const parentUrlObj = new URL(monitorForm.url)
        const parentPath = parentUrlObj.pathname + parentUrlObj.search
        const isPrivateIP = /^(192\.168\.|10\.|172\.(1[6-9]|2[0-9]|3[01])\.)/.test(trimmedIp)
        const protocol = isPrivateIP ? 'http' : parentUrlObj.protocol.replace(':', '')
        childItem.url = `${protocol}://${trimmedIp}${parentPath}`

        monitorForm.childUrls.push(childItem)
      } catch (e) {
        console.error('URL生成失败', e)
      }
    })

    ElMessage.success(`已批量生成 ${ips.length} 个子域名配置`)
  } catch (error) {
    // 用户取消
  }
}

// 删除
const handleDelete = async (row) => {
  try {
    // 检查是否有子监控
    if (row.childCount && row.childCount > 0) {
      // 有子监控,显示删除策略选择对话框
      const { value: strategy } = await ElMessageBox.prompt(
        `该监控有 ${row.childCount} 个子监控,请选择处理方式:`,
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          inputType: 'select',
          inputOptions: {
            cascade: '级联删除 (同时删除所有子监控)',
            orphan: '孤立处理 (子监控变为独立监控)',
            prevent: '防止删除 (保留所有监控)'
          },
          inputPlaceholder: '请选择删除策略',
          inputValidator: (value) => {
            if (!value) return '请选择删除策略'
            return true
          }
        }
      )

      if (strategy === 'prevent') {
        ElMessage.info('已取消删除操作')
        return
      }

      const { data: res } = await apiMonitorApi.deleteWithStrategy(row.id, strategy)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        if (selectedFunctionId.value) {
          fetchFunctionApiList(selectedFunctionId.value)
        } else {
          fetchList()
        }
        fetchTreeData()
      } else {
        ElMessage.error(res.message || '删除失败')
      }
    } else {
      // 没有子监控,直接删除
      await ElMessageBox.confirm('确定要删除该 API 监控吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      const { data: res } = await apiMonitorApi.delete(row.id)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        if (selectedFunctionId.value) {
          fetchFunctionApiList(selectedFunctionId.value)
        } else {
          fetchList()
        }
        fetchTreeData()
      } else {
        ElMessage.error(res.message || '删除失败')
      }
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error(error)
    }
  }
}

// 提交表单 - 使用简化版接口
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    formLoading.value = true
    try {
      const data = {
        name: monitorForm.name,
        url: monitorForm.url,
        method: monitorForm.method,
        groupId: Number(monitorForm.groupId),
        businessFunctionName: monitorForm.businessFunctionName,
        description: monitorForm.description
      }

      // 如果 headers 不为空，添加到请求数据中
      if (monitorForm.headers && monitorForm.headers.trim()) {
        data.headers = monitorForm.headers.trim()
      }

      let result
      let parentMonitorId

      if (monitorForm.id) {
        // 编辑时使用完整 API
        result = await apiMonitorApi.update(monitorForm.id, data)
        parentMonitorId = monitorForm.id
      } else {
        // 新增时使用简化版 API
        result = await apiMonitorApi.createSimple(data)
      }

      const { data: res } = result
      if (res.code === 200) {
        ElMessage.success(monitorForm.id ? '更新成功' : '创建成功')

        // 如果是新建，获取新创建的监控ID
        if (!monitorForm.id) {
          parentMonitorId = res.data.id || res.data
        }

        // 处理子域名配置
        if (monitorForm.childUrls && monitorForm.childUrls.length > 0) {
          await handleChildUrlsSubmit(parentMonitorId)
        }

        dialogVisible.value = false
        if (selectedFunctionId.value) {
          fetchFunctionApiList(selectedFunctionId.value)
        } else {
          fetchList()
        }
        fetchTreeData()
      } else {
        ElMessage.error(res.message || '操作失败')
      }
    } catch (error) {
      ElMessage.error('操作失败: ' + (error.message || ''))
      console.error(error)
    } finally {
      formLoading.value = false
    }
  })
}

// 处理子域名配置的提交
const handleChildUrlsSubmit = async (parentMonitorId) => {
  const newChildren = monitorForm.childUrls.filter(child => !child.id) // 新增的子监控
  const existingChildren = monitorForm.childUrls.filter(child => child.id) // 已有的子监控

  // 创建新的子监控
  for (const child of newChildren) {
    if (!child.name || !child.url) continue

    try {
      // 创建子监控，带上 parentId
      const createData = {
        name: child.name,
        url: child.url,
        method: monitorForm.method,
        groupId: monitorForm.groupId,
        businessFunctionName: monitorForm.businessFunctionName,
        description: `${monitorForm.name}的子监控`,
        parentId: parentMonitorId // 添加父监控ID
      }

      const { data: createRes } = await apiMonitorApi.createSimple(createData)
      if (createRes.code !== 200) {
        console.error('创建子监控失败:', createRes.message)
        ElMessage.warning(`子监控 ${child.name} 创建失败: ${createRes.message}`)
      }
    } catch (error) {
      console.error('创建子监控失败', error)
      ElMessage.warning(`子监控 ${child.name} 创建失败`)
    }
  }

  // 更新已有的子监控
  for (const child of existingChildren) {
    if (!child.name || !child.url || !child.id) continue

    try {
      const updateData = {
        name: child.name,
        url: child.url,
        method: monitorForm.method,
        groupId: monitorForm.groupId,
        businessFunctionName: monitorForm.businessFunctionName
      }

      await apiMonitorApi.update(child.id, updateData)
    } catch (error) {
      console.error('更新子监控失败', error)
    }
  }

  if (newChildren.length > 0) {
    ElMessage.success(`已批量创建 ${newChildren.length} 个子监控`)
  }
}

// 处理状态变更
const handleStatusChange = async (row) => {
  try {
    const { data: res } = await apiMonitorApi.updateStatus(row.id, row.status)
    if (res.code === 200) {
      ElMessage.success(row.status === 1 ? '监控已启用' : '监控已禁用')
    } else {
      row.status = row.status === 1 ? 0 : 1
      ElMessage.error(res.message || '状态更新失败')
    }
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('状态更新失败')
    console.error(error)
  }
}

// 测试 API 接口
const handleTest = async (row) => {
  testLoading.value = true
  try {
    const testData = {
      url: row.url,
      method: row.method || 'GET',
      timeout: 5000
    }

    const { data: res } = await apiMonitorApi.test(testData)
    if (res.code === 200) {
      const result = res.data
      ElMessage.success(`测试成功 - 状态码: ${result.statusCode}, 响应时间: ${result.responseTime}ms`)
    } else {
      ElMessage.error(res.message || '测试失败')
    }
  } catch (error) {
    ElMessage.error('测试失败: ' + error.message)
    console.error(error)
  } finally {
    testLoading.value = false
  }
}

// 处理下拉菜单命令
const handleCommand = (command, row) => {
  switch (command) {
    case 'test':
      handleTest(row)
      break
    case 'edit':
      handleEdit(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

// ==================== 分层监控操作方法 ====================

// 打开监控详情抽屉
const handleViewDetail = async (row) => {
  try {
    drawerLoading.value = true
    drawerVisible.value = true

    // 获取完整的监控详情（包含 childMonitors）
    const { data: res } = await apiMonitorApi.getDetail(row.id)
    if (res.code === 200) {
      const detail = res.data
      // 将 childMonitors 映射为 children 以统一命名
      if (detail.childMonitors && detail.childMonitors.length > 0) {
        detail.children = detail.childMonitors
      }
      drawerMonitorDetail.value = detail
    } else {
      ElMessage.error(res.message || '获取监控详情失败')
      drawerMonitorDetail.value = row
    }
  } catch (error) {
    ElMessage.error('获取监控详情失败')
    console.error(error)
    drawerMonitorDetail.value = row
  } finally {
    drawerLoading.value = false
  }
}

// 关闭抽屉
const closeDrawer = () => {
  drawerVisible.value = false
  drawerMonitorDetail.value = null
}

// 移除父监控关系
const handleRemoveParent = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要移除该监控与父监控"${row.parentName || '未知'}"的关系吗? 移除后该监控将成为独立监控。`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const { data: res } = await apiMonitorApi.removeParent(row.id)
    if (res.code === 200) {
      ElMessage.success('父监控关系已移除')
      fetchList()
      fetchTreeData()
      closeDrawer() // 关闭抽屉
    } else {
      ElMessage.error(res.message || '移除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('移除失败')
      console.error(error)
    }
  }
}

// 删除子监控（从抽屉中）
const handleDeleteChild = async (child) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除子监控"${child.name}"吗? 此操作将永久删除该监控及其所有检测记录。`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const { data: res } = await apiMonitorApi.delete(child.id)
    if (res.code === 200) {
      ElMessage.success('子监控删除成功')

      // 从抽屉详情中移除该子监控
      if (drawerMonitorDetail.value && drawerMonitorDetail.value.children) {
        const index = drawerMonitorDetail.value.children.findIndex(c => c.id === child.id)
        if (index > -1) {
          drawerMonitorDetail.value.children.splice(index, 1)
          // 更新子监控数量
          if (drawerMonitorDetail.value.childCount > 0) {
            drawerMonitorDetail.value.childCount--
          }
        }
      }

      // 刷新列表和树形数据
      fetchList()
      fetchTreeData()

      // 如果没有子监控了，可以关闭抽屉
      if (drawerMonitorDetail.value && (!drawerMonitorDetail.value.children || drawerMonitorDetail.value.children.length === 0)) {
        ElMessage.info('已删除所有子监控')
      }
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

// 树节点点击事件
const handleNodeClick = (data) => {
  console.log('树节点点击:', data)

  if (data.type === 'group') {
    // 点击二级分组 - 显示该分组下所有业务功能的所有API
    selectedFunctionId.value = null
    selectedNodeData.value = data
    selectedGroupName.value = data.label.replace(/\s*\(\d+\)$/, '') // 移除总数部分
    searchForm.keyword = ''
    searchForm.status = -1
    console.log('开始加载业务组API，groupId:', data.groupId)
    fetchGroupApiList(data.groupId)
  } else if (data.type === 'parent') {
    // 点击一级分组 - 不做任何操作，仅展开/折叠
    selectedFunctionId.value = null
    selectedNodeData.value = null
    selectedGroupName.value = ''
    tableData.value = []
    total.value = 0
  }
}

// ==================== 生命周期 ====================

onMounted(async () => {
  await fetchTreeData()
  await fetchList()

  // 启动自动刷新定时器（每60秒刷新一次）
  autoRefreshTimer.value = setInterval(() => {
    if (selectedFunctionId.value) {
      fetchFunctionApiList(selectedFunctionId.value)
    } else {
      fetchList()
    }
  }, 60000)
})

// 组件卸载时清除定时器
onUnmounted(() => {
  if (autoRefreshTimer.value) {
    clearInterval(autoRefreshTimer.value)
  }
})
</script>

<template>
  <div class="api-monitor-container">
    <!-- 左右布局 -->
    <div class="main-content">
      <!-- 左侧树形菜单 -->
      <div class="tree-sidebar">
        <div class="tree-header">
          <div class="title-wrapper">
            <div class="title-icon">
              <div class="pulse-ring"></div>
              <el-icon class="main-icon"><Pointer /></el-icon>
            </div>
            <div class="title-content">
              <h3 class="title">业务组/功能</h3>
              <span class="subtitle">Business Groups</span>
            </div>
          </div>
          <div class="header-actions">
            <div class="stats-indicator">
              <span class="stats-text">{{ treeData.length }}</span>
              <span class="stats-label">Groups</span>
            </div>
            <div class="refresh-btn" @click="fetchTreeData" title="刷新分组">
              <el-icon class="refresh-icon"><Refresh /></el-icon>
            </div>
          </div>
        </div>
        <div class="tree-content">
          <el-tree
            :data="treeData"
            node-key="id"
            :props="{ children: 'children', label: 'label' }"
            @node-click="handleNodeClick"
            :accordion="true"
            highlight-current
            class="tech-tree"
          >
            <template v-slot="{ node, data }">
              <div
                class="tree-node"
                :class="{ 'parent-node': data.type === 'parent' }"
              >
                <div class="node-content">
                  <div class="node-icon-wrapper">
                    <!-- 一级分组 -->
                    <template v-if="data.type === 'parent'">
                      <div class="parent-icon">
                        <div class="icon-bg"></div>
                        <el-icon class="parent-icon-el"><FolderOpened /></el-icon>
                      </div>
                    </template>
                    <!-- 二级分组 -->
                    <template v-else-if="data.type === 'group'">
                      <div class="child-icon">
                        <div class="dot-indicator"></div>
                        <el-icon class="child-icon-el"><Folder /></el-icon>
                      </div>
                    </template>
                  </div>
                  <div class="node-text">
                    <span class="node-label">
                      {{ node.label }}
                    </span>
                    <!-- 健康状态指示器 -->
                    <div class="health-indicators" v-if="data.type === 'group'">
                      <span class="health-badge healthy" v-if="data.healthy > 0">
                        <img :src="healthIcon" style="width: 12px; height: 12px" />
                        {{ data.healthy }}
                      </span>
                      <span class="health-badge unhealthy" v-if="data.unhealthy > 0">
                        <img :src="fireIcon" style="width: 12px; height: 12px" />
                        {{ data.unhealthy }}
                      </span>
                    </div>
                  </div>
                </div>
                <div class="hover-effect"></div>
              </div>
            </template>
          </el-tree>
        </div>
      </div>

      <!-- 右侧内容区 -->
      <div class="content-area">
        <!-- 顶部统计卡片 -->
        <div class="stats-card-section">
          <div class="stats-card total">
            <div class="stats-icon-wrapper">
              <el-icon class="stats-icon"><DataBoard /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-label">总 API 数</div>
              <div class="stats-value">{{ total }}</div>
            </div>
          </div>
          <div class="stats-card healthy">
            <div class="stats-icon-wrapper">
              <img :src="healthIcon" class="stats-icon-img" alt="健康" />
            </div>
            <div class="stats-info">
              <div class="stats-label">健康</div>
              <div class="stats-value">{{ tableData.filter(d => d.healthStatus === 'healthy').length }}</div>
            </div>
          </div>
          <div class="stats-card unhealthy">
            <div class="stats-icon-wrapper">
              <img :src="fireIcon" class="stats-icon-img" alt="异常" />
            </div>
            <div class="stats-info">
              <div class="stats-label">异常</div>
              <div class="stats-value">{{ tableData.filter(d => d.healthStatus === 'unhealthy').length }}</div>
            </div>
          </div>
        </div>

        <!-- 搜索和筛选区域 -->
        <div class="search-section">
          <el-form :inline="true" :model="searchForm" class="search-form">
            <el-form-item label="关键字">
              <el-input
                v-model="searchForm.keyword"
                placeholder="搜索 API 名称或地址"
                size="small"
                clearable
                @keyup.enter="handleSearch"
                style="width: 200px"
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
              <el-button
                type="success"
                size="small"
                @click="handleAdd"
                v-authority="['monitor:api:add']"
                class="modern-btn success-btn"
              >
                <el-icon><Plus /></el-icon>
                新增监控
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 卡片网格 -->
        <div class="cards-section" v-loading="loading">
          <div class="api-cards-grid">
            <div
              v-for="row in tableData"
              :key="row.id"
              class="api-card"
              :class="{ 'unhealthy': row.healthStatus === 'unhealthy', 'has-children': row.childCount > 0 }"
              @click="handleViewDetail(row)"
            >
              <!-- 卡片头部 -->
              <div class="card-header">
                <div class="card-title-section">
                  <div class="card-title-content">
                    <h4 class="card-title">
                      {{ row.name }}
                      <!-- 子监控数量标识 -->
                      <el-tag
                        v-if="row.childCount > 0"
                        type="warning"
                        size="small"
                        style="margin-left: 6px"
                      >
                        {{ row.childCount }} 个子监控
                      </el-tag>
                    </h4>
                    <span class="card-function">{{ row.businessFunctionName }}</span>
                  </div>
                </div>
                <div class="card-header-actions" @click.stop>
                  <el-switch
                    v-model="row.status"
                    :active-value="1"
                    :inactive-value="0"
                    @change="handleStatusChange(row)"
                    v-authority="['monitor:api:status']"
                    size="small"
                  />
                  <el-dropdown trigger="click" @command="(command) => handleCommand(command, row)">
                    <el-button text class="more-icon-btn">
                      <el-icon><MoreFilled /></el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item command="test" :icon="Pointer">
                          <span :class="{ 'loading-text': testLoading }">测试</span>
                        </el-dropdown-item>
                        <el-dropdown-item command="edit" :icon="Edit" v-authority="['monitor:api:edit']">
                          编辑
                        </el-dropdown-item>
                        <el-dropdown-item command="delete" :icon="Delete" v-authority="['monitor:api:delete']" divided>
                          删除
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </div>

              <!-- 卡片内容 -->
              <div class="card-body">
                <!-- API 方法和地址 -->
                <div class="card-info-row">
                  <div class="card-status-icon-wrapper">
                    <img
                      :src="row.healthStatus === 'healthy' ? healthIcon : fireIcon"
                      :alt="row.healthStatus"
                      class="status-icon"
                    />
                  </div>
                  <el-tag :type="row.method === 'GET' ? 'success' : 'primary'" size="small" class="method-tag">
                    {{ row.method || 'GET' }}
                  </el-tag>
                  <span class="api-url" :title="row.url">{{ row.url }}</span>
                </div>

                <!-- 父监控聚合状态显示 -->
                <div v-if="row.childCount > 0" class="child-stats-row">
                  <div class="child-stats-info">
                    <span class="stats-label">子监控状态:</span>
                    <span class="stats-value">
                      <span class="healthy-count">{{ row.healthyChildren || 0 }}</span>
                      <span class="divider">/</span>
                      <span class="total-count">{{ row.childCount }}</span>
                      <span class="stats-text">健康</span>
                    </span>
                  </div>
                  <el-icon class="view-detail-icon"><ArrowRight /></el-icon>
                </div>

                <!-- 监控数据 -->
                <div class="card-metrics">
                  <div class="metric-item">
                    <span class="metric-label">健康率</span>
                    <span class="metric-value health-rate" :class="{ 'rate-good': row.healthRate >= 95, 'rate-warning': row.healthRate < 95 && row.healthRate >= 80, 'rate-danger': row.healthRate < 80 }">
                      {{ row.healthRate !== undefined && row.healthRate !== null ? row.healthRate + '%' : '-' }}
                    </span>
                  </div>
                  <div class="metric-divider"></div>
                  <div class="metric-item">
                    <span class="metric-label">响应时间</span>
                    <span class="metric-value">{{ row.lastResponseTime || '-' }} ms</span>
                  </div>
                  <div class="metric-divider"></div>
                  <div class="metric-item">
                    <span class="metric-label">状态码</span>
                    <el-tag
                      :type="row.lastStatusCode === 200 ? 'success' : 'danger'"
                      size="small"
                      class="status-code-tag"
                    >
                      {{ row.lastStatusCode || '-' }}
                    </el-tag>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 无数据提示 -->
          <el-empty v-if="!loading && tableData.length === 0" description="暂无 API 监控数据" />
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

    <!-- 新增/编辑对话框 - 简化版表单 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="650px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="monitorForm"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item label="业务组" prop="groupId">
          <el-select v-model="monitorForm.groupId" placeholder="请选择业务组（二级）" style="width: 100%">
            <el-option
              v-for="group in secondLevelGroups"
              :key="group.id"
              :label="`${group.parentName} / ${group.name}`"
              :value="group.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="业务功能" prop="businessFunctionName">
          <el-input
            v-model="monitorForm.businessFunctionName"
            placeholder="输入现有或新的功能名称（自动创建）"
          />
        </el-form-item>

        <el-form-item label="监控名称" prop="name">
          <el-input v-model="monitorForm.name" placeholder="请输入监控名称，同组内唯一" />
        </el-form-item>

        <el-form-item label="接口地址" prop="url">
          <el-input v-model="monitorForm.url" placeholder="请输入完整的接口地址，如 http://example.com/api/login" />
        </el-form-item>

        <el-form-item label="HTTP 方法">
          <el-select v-model="monitorForm.method" style="width: 100%">
            <el-option v-for="item in methodOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="Headers" prop="headers">
          <el-input
            v-model="monitorForm.headers"
            type="textarea"
            :rows="3"
            placeholder='请输入 Headers（JSON 格式，非必填）&#10;例如：{"Authorization": "Bearer token", "Content-Type": "application/json"}'
          />
          <div style="font-size: 12px; color: #909399; margin-top: 4px;">
            提示：Headers 必须是有效的 JSON 格式字符串，如：{"key": "value"}
          </div>
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="monitorForm.description" type="textarea" :rows="3" placeholder="请输入监控描述" />
        </el-form-item>

        <!-- 子域名配置区域 -->
        <el-divider content-position="left">
          <span style="font-weight: 600; color: #409eff;">子域名配置（可选）</span>
        </el-divider>

        <div style="background: #f5f7fa; padding: 16px; border-radius: 8px; margin-bottom: 16px;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
            <span style="font-size: 13px; color: #606266; font-weight: 500;">
              配置子域名 ({{ monitorForm.childUrls.length }} 个)
            </span>
            <div style="display: flex; gap: 8px;">
              <el-button size="small" @click="batchGenerateChildUrls" :disabled="!monitorForm.url">
                <el-icon><DocumentAdd /></el-icon>
                批量生成
              </el-button>
              <el-button size="small" type="primary" @click="addChildUrl">
                <el-icon><Plus /></el-icon>
                添加一行
              </el-button>
            </div>
          </div>

          <!-- 子域名配置列表 -->
          <div v-if="monitorForm.childUrls.length > 0" style="max-height: 300px; overflow-y: auto;">
            <div
              v-for="(child, index) in monitorForm.childUrls"
              :key="child.tempId"
              style="display: flex; gap: 8px; margin-bottom: 8px; align-items: flex-start;"
            >
              <el-input
                v-model="child.name"
                placeholder="子监控名称"
                size="small"
                style="width: 180px;"
              />
              <el-input
                v-model="child.url"
                placeholder="子监控URL"
                size="small"
                style="flex: 1;"
              >
                <template #append>
                  <el-button
                    size="small"
                    @click="async () => {
                      const { value } = await ElMessageBox.prompt('请输入IP或域名', '快速生成', {
                        inputPlaceholder: '例如: 192.168.7.6:7001'
                      }).catch(() => ({ value: null }));
                      if (value) quickGenerateChildUrl(index, value);
                    }"
                  >
                    填充
                  </el-button>
                </template>
              </el-input>
              <el-button
                size="small"
                type="danger"
                :icon="Delete"
                @click="removeChildUrl(index)"
              >删除</el-button>
            </div>
          </div>

          <el-empty
            v-else
            description="暂无子域名配置"
            :image-size="60"
            style="padding: 20px 0;"
          />

          <div style="font-size: 12px; color: #909399; margin-top: 12px; line-height: 1.6;">
            💡 提示:<br>
            • 点击"批量生成"可一次性输入多个IP/域名(每行一个)<br>
            • 点击"填充"可基于父监控路径快速生成完整URL<br>
            • 子监控会自动继承父监控的业务组、业务功能和HTTP方法
          </div>
        </div>

        <!-- 自动配置提示 -->
        <div style="background: #f0f9ff; padding: 12px; border-radius: 8px; margin-bottom: 16px; border-left: 4px solid #409eff;">
          <div style="font-size: 13px; color: #333; font-weight: 500; margin-bottom: 8px;">🚀 后端自动配置：</div>
          <div style="font-size: 12px; color: #666; line-height: 1.8;">
            • 检测间隔：30秒 | 超时时间：5秒<br>
            • 期望状态：200 | 重试次数：3次<br>
            • 告警已启用 | 监控已启用
          </div>
        </div>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="formLoading">
          {{ monitorForm.id ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 监控详情抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      title="监控详情"
      :size="600"
      :close-on-click-modal="false"
      direction="rtl"
    >
      <div v-loading="drawerLoading" class="drawer-content">
        <template v-if="drawerMonitorDetail">
          <!-- 父监控卡片 -->
          <div class="drawer-section">
            <h3 class="drawer-section-title">
              <el-icon><DataBoard /></el-icon>
              父监控信息
            </h3>
            <div class="detail-card parent-detail-card">
              <div class="detail-card-header">
                <div class="detail-card-title">
                  <h4>{{ drawerMonitorDetail.name }}</h4>
                  <el-tag :type="drawerMonitorDetail.status === 1 ? 'success' : 'info'" size="small">
                    {{ drawerMonitorDetail.status === 1 ? '启用' : '禁用' }}
                  </el-tag>
                </div>
                <div class="detail-card-status">
                  <img
                    :src="drawerMonitorDetail.healthStatus === 'healthy' ? healthIcon : fireIcon"
                    class="status-icon-large"
                    :alt="drawerMonitorDetail.healthStatus"
                  />
                </div>
              </div>

              <div class="detail-card-body">
                <div class="detail-item">
                  <span class="detail-label">业务组:</span>
                  <span class="detail-value">{{ drawerMonitorDetail.groupName || '-' }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">业务功能:</span>
                  <span class="detail-value">{{ drawerMonitorDetail.businessFunctionName || '-' }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">请求方法:</span>
                  <el-tag :type="drawerMonitorDetail.method === 'GET' ? 'success' : 'primary'" size="small">
                    {{ drawerMonitorDetail.method || 'GET' }}
                  </el-tag>
                </div>
                <div class="detail-item">
                  <span class="detail-label">URL:</span>
                  <span class="detail-value url-value">{{ drawerMonitorDetail.url }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">健康率:</span>
                  <span
                    class="detail-value health-rate"
                    :class="{
                      'rate-good': drawerMonitorDetail.healthRate >= 95,
                      'rate-warning': drawerMonitorDetail.healthRate < 95 && drawerMonitorDetail.healthRate >= 80,
                      'rate-danger': drawerMonitorDetail.healthRate < 80
                    }"
                  >
                    {{ drawerMonitorDetail.healthRate !== undefined && drawerMonitorDetail.healthRate !== null ? drawerMonitorDetail.healthRate + '%' : '-' }}
                  </span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">最后响应时间:</span>
                  <span class="detail-value">{{ drawerMonitorDetail.lastResponseTime || '-' }} ms</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">最后状态码:</span>
                  <el-tag
                    :type="drawerMonitorDetail.lastStatusCode === 200 ? 'success' : 'danger'"
                    size="small"
                  >
                    {{ drawerMonitorDetail.lastStatusCode || '-' }}
                  </el-tag>
                </div>
                <div class="detail-item" v-if="drawerMonitorDetail.description">
                  <span class="detail-label">描述:</span>
                  <span class="detail-value">{{ drawerMonitorDetail.description }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 子监控列表 -->
          <div v-if="drawerMonitorDetail.children && drawerMonitorDetail.children.length > 0" class="drawer-section">
            <h3 class="drawer-section-title">
              <el-icon><FolderOpened /></el-icon>
              子监控列表 ({{ drawerMonitorDetail.children.length }})
            </h3>

            <div class="child-cards-list">
              <div
                v-for="child in drawerMonitorDetail.children"
                :key="child.id"
                class="child-detail-card"
              >
                <div class="child-card-header">
                  <div class="child-card-title">
                    <img
                      :src="child.healthStatus === 'healthy' ? healthIcon : fireIcon"
                      class="status-icon-small"
                      :alt="child.healthStatus"
                    />
                    <h5>{{ child.name }}</h5>
                  </div>
                  <div class="child-card-actions" @click.stop>
                    <el-button
                      text
                      type="primary"
                      size="small"
                      @click="handleTest(child)"
                    >
                      测试
                    </el-button>
                    <el-button
                      text
                      type="danger"
                      size="small"
                      @click="handleDeleteChild(child)"
                      v-authority="['monitor:api:delete']"
                    >
                      删除
                    </el-button>
                  </div>
                </div>

                <div class="child-card-body">
                  <div class="child-detail-row">
                    <el-tag :type="child.method === 'GET' ? 'success' : 'primary'" size="small">
                      {{ child.method || 'GET' }}
                    </el-tag>
                    <span class="child-url">{{ child.url }}</span>
                  </div>

                  <div class="child-metrics">
                    <div class="child-metric-item">
                      <span class="metric-label">健康率</span>
                      <span
                        class="metric-value"
                        :class="{
                          'rate-good': child.healthRate >= 95,
                          'rate-warning': child.healthRate < 95 && child.healthRate >= 80,
                          'rate-danger': child.healthRate < 80
                        }"
                      >
                        {{ child.healthRate !== undefined && child.healthRate !== null ? child.healthRate + '%' : '-' }}
                      </span>
                    </div>
                    <div class="child-metric-item">
                      <span class="metric-label">响应时间</span>
                      <span class="metric-value">{{ child.lastResponseTime || '-' }} ms</span>
                    </div>
                    <div class="child-metric-item">
                      <span class="metric-label">状态码</span>
                      <el-tag
                        :type="child.lastStatusCode === 200 ? 'success' : 'danger'"
                        size="small"
                      >
                        {{ child.lastStatusCode || '-' }}
                      </el-tag>
                    </div>
                  </div>

                  <div v-if="child.lastErrorMessage" class="child-error-message">
                    <el-icon color="#F56C6C"><Warning /></el-icon>
                    <span>{{ child.lastErrorMessage }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 无子监控提示 -->
          <div v-else class="drawer-section">
            <el-empty description="该监控暂无子监控" :image-size="100" />
          </div>
        </template>
      </div>

      <template #footer>
        <div class="drawer-footer">
          <el-button @click="closeDrawer">关闭</el-button>
          <el-button type="primary" @click="handleEdit(drawerMonitorDetail)" v-authority="['monitor:api:edit']">
            编辑监控
          </el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<style scoped>
.api-monitor-container {
  padding: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 左右布局 */
.main-content {
  display: flex;
  gap: 20px;
  flex: 1;
  min-height: 0;
}

/* 左侧树形菜单栏 - 参考 CmdbGroup.vue 的科技感设计 */
.tree-sidebar {
  width: 280px;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.95), rgba(255, 255, 255, 0.85));
  backdrop-filter: blur(20px);
  border-radius: 20px;
  box-shadow:
    0 20px 40px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
  overflow: hidden;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  position: relative;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.tree-sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.5), transparent);
  z-index: 1;
}

.tree-sidebar:hover {
  transform: translateY(-2px);
  box-shadow:
    0 25px 50px rgba(0, 0, 0, 0.15),
    0 0 0 1px rgba(255, 255, 255, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

/* 科技感标题区域 */
.tree-header {
  padding: 20px 20px 15px;
  background: transparent;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  flex-shrink: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pulse-ring {
  position: absolute;
  width: 35px;
  height: 35px;
  border: 2px solid rgba(102, 126, 234, 0.3);
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.2); opacity: 0.7; }
  100% { transform: scale(1.4); opacity: 0; }
}

.main-icon {
  font-size: 20px;
  color: #667eea;
  z-index: 2;
  background: rgba(102, 126, 234, 0.1);
  padding: 8px;
  border-radius: 8px;
}

.title-content .title {
  font-size: 16px;
  font-weight: 700;
  margin: 0 0 2px 0;
  background: linear-gradient(135deg, #667eea, #764ba2);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  white-space: nowrap;
}

.subtitle {
  font-size: 10px;
  color: rgba(102, 126, 234, 0.6);
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 500;
}

.stats-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: rgba(102, 126, 234, 0.05);
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.1);
  transition: all 0.3s ease;
}

.stats-indicator:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.stats-text {
  font-size: 16px;
  font-weight: 700;
  color: #667eea;
  line-height: 1;
}

.stats-label {
  font-size: 10px;
  color: rgba(102, 126, 234, 0.6);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* 树形内容区域 */
.tree-content {
  flex: 1;
  overflow-y: auto;
  padding: 15px 0;
}

/* 高科技树形结构 */
.tech-tree {
  border: none;
  background: transparent;
}

.tech-tree :deep(.el-tree-node) {
  margin: 0;
}

.tech-tree :deep(.el-tree-node__content) {
  height: auto;
  padding: 0;
  background: transparent;
  border-radius: 0;
  position: relative;
}

/* 树节点展开按钮样式 */
.tech-tree :deep(.el-tree-node__expand-icon) {
  position: absolute;
  left: 25px !important;
  top: 50%;
  transform: translateY(-50%);
  z-index: 5;
  color: #667eea;
  font-size: 12px;
  transition: transform 0.3s ease;
}

.tech-tree :deep(.el-tree-node__expand-icon.expanded) {
  transform: translateY(-50%) rotate(90deg);
}

.tech-tree :deep(.el-tree-node__expand-icon.is-leaf) {
  display: none;
}

.tree-node {
  position: relative;
  margin: 0px 15px 0px 15px;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  animation: dataLoad 0.4s ease-out;
}

@keyframes dataLoad {
  0% { opacity: 0; transform: translateY(20px); }
  100% { opacity: 1; transform: translateY(0); }
}

.hover-effect {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.08), rgba(118, 75, 162, 0.08));
  border-radius: 12px;
  opacity: 0;
  transform: scale(1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1;
}

.tree-node.parent-node:hover > .hover-effect {
  opacity: 1;
}

.tree-node:not(.parent-node) .hover-effect {
  display: none;
}

.node-content {
  display: flex;
  align-items: center;
  padding: 2px 12px;
  position: relative;
  z-index: 2;
}

.parent-node .node-content {
  padding: 3px 12px;
}

.node-icon-wrapper {
  margin-left: 16px;
  margin-right: 2px;
  display: flex;
  align-items: center;
}

.parent-icon {
  position: relative;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.parent-icon .icon-bg {
  display: none;
}

.parent-icon-el {
  font-size: 18px;
  color: #F7BA2A;
  z-index: 1;
  transition: all 0.3s ease;
}

.child-icon {
  position: relative;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dot-indicator {
  position: absolute;
  width: 6px;
  height: 6px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 50%;
  top: 50%;
  left: -8px;
  transform: translateY(-50%);
}

.child-icon-el {
  font-size: 16px;
  color: #F7BA2A;
}

.node-text {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.node-label {
  font-size: 14px;
  font-weight: 400;
  color: #2c3e50;
  transition: all 0.3s ease;
}

.parent-node .node-label {
  font-weight: 400;
  font-size: 15px;
}

.health-indicators {
  display: flex;
  align-items: center;
  gap: 4px;
}

.health-badge {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.health-badge.healthy {
  background: rgba(103, 194, 58, 0.1);
  color: #67C23A;
}

.health-badge.unhealthy {
  background: rgba(245, 108, 108, 0.1);
  color: #F56C6C;
}

/* 选中状态 */
.tech-tree :deep(.el-tree-node.is-current > .el-tree-node__content) .tree-node {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.08), rgba(118, 75, 162, 0.08));
  border: none;
}

.tech-tree :deep(.el-tree-node.is-current > .el-tree-node__content) .hover-effect {
  opacity: 0 !important;
}

.tech-tree :deep(.el-tree-node.is-current > .el-tree-node__content) .node-label {
  color: #2c3e50;
  font-weight: 500;
}

/* 滚动条美化 */
.tree-content::-webkit-scrollbar {
  width: 4px;
}

.tree-content::-webkit-scrollbar-track {
  background: rgba(102, 126, 234, 0.05);
  border-radius: 2px;
}

.tree-content::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.3), rgba(118, 75, 162, 0.3));
  border-radius: 2px;
}

.tree-content::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.5), rgba(118, 75, 162, 0.5));
}

/* 刷新按钮样式 */
.refresh-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.3s ease;
  background: rgba(102, 126, 234, 0.05);
  border: 1px solid rgba(102, 126, 234, 0.1);
}

.refresh-btn:hover {
  transform: scale(1.05);
  background: rgba(102, 126, 234, 0.15);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.refresh-btn:active {
  transform: scale(0.95);
}

.refresh-icon {
  font-size: 16px;
  color: #667eea;
  transition: all 0.3s ease;
}

.refresh-btn:hover .refresh-icon {
  transform: rotate(180deg);
  color: #764ba2;
}

.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

/* 统计卡片 - 重新设计更小巧专业 */
.stats-card-section {
  display: flex;
  gap: 12px;
}

.stats-card {
  flex: 1;
  background: white;
  border-radius: 10px;
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.stats-card::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  transition: all 0.3s ease;
}

.stats-card.total::before {
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
}

.stats-card.healthy::before {
  background: linear-gradient(180deg, #67C23A 0%, #529b2e 100%);
}

.stats-card.unhealthy::before {
  background: linear-gradient(180deg, #F56C6C 0%, #dd6161 100%);
}

.stats-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

.stats-icon-wrapper {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.3s ease;
}

.stats-card.total .stats-icon-wrapper {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
}

.stats-card.healthy .stats-icon-wrapper {
  background: linear-gradient(135deg, rgba(103, 194, 58, 0.1), rgba(82, 155, 46, 0.1));
}

.stats-card.unhealthy .stats-icon-wrapper {
  background: linear-gradient(135deg, rgba(245, 108, 108, 0.1), rgba(221, 97, 97, 0.1));
}

.stats-card:hover .stats-icon-wrapper {
  transform: scale(1.1);
}

.stats-icon {
  font-size: 22px;
}

.stats-card.total .stats-icon {
  color: #667eea;
}

.stats-icon-img {
  width: 22px;
  height: 22px;
}

.stats-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stats-label {
  font-size: 13px;
  font-weight: 500;
  color: #606266;
  line-height: 1;
}

.stats-value {
  font-size: 26px;
  font-weight: 700;
  line-height: 1.2;
}

.stats-card.total .stats-value {
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.stats-card.healthy .stats-value {
  color: #67C23A;
}

.stats-card.unhealthy .stats-value {
  color: #F56C6C;
}

/* 搜索区域 */
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

/* 操作按钮区域 */
.action-section {
  display: flex;
  gap: 10px;
}

/* 卡片区域 */
.cards-section {
  flex: 1;
  min-height: 0;
  overflow: auto;
}

.api-cards-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 4px;
}

@media (max-width: 1600px) {
  .api-cards-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 1200px) {
  .api-cards-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .api-cards-grid {
    grid-template-columns: 1fr;
  }
}

/* API 卡片样式 */
.api-card {
  background: white;
  border-radius: 10px;
  padding: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  gap: 8px;
  position: relative;
  overflow: hidden;
  cursor: pointer;
}

.api-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.api-card.has-children:hover {
  border-color: #409eff;
  box-shadow: 0 8px 24px rgba(64, 158, 255, 0.2);
}

/* 卡片头部 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
}

.card-title-section {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.card-title-content {
  flex: 1;
  min-width: 0;
}

.card-title {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
  line-height: 1.3;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-function {
  font-size: 11px;
  color: #909399;
  display: block;
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-header-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.more-icon-btn {
  padding: 3px;
  color: #909399;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.more-icon-btn:hover {
  color: #667eea;
  background: rgba(102, 126, 234, 0.1);
}

.more-icon-btn .el-icon {
  font-size: 16px;
}

/* 卡片内容 */
.card-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.card-info-row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 6px;
}

.card-status-icon-wrapper {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: rgba(103, 194, 58, 0.1);
  transition: all 0.3s ease;
  border: 2px solid rgba(103, 194, 58, 0.3);
}

.api-card.unhealthy .card-status-icon-wrapper {
  background: rgba(245, 108, 108, 0.1);
  border-color: rgba(245, 108, 108, 0.3);
}

.card-status-icon-wrapper .status-icon {
  width: 18px;
  height: 18px;
}

.method-tag {
  flex-shrink: 0;
}

.api-url {
  font-size: 11px;
  color: #606266;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

/* 监控指标 */
.card-metrics {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 0;
}

.metric-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.metric-label {
  font-size: 10px;
  color: #909399;
  font-weight: 500;
}

.metric-value {
  font-size: 14px;
  font-weight: 700;
  color: #2c3e50;
}

.health-rate.rate-good {
  color: #67C23A;
}

.health-rate.rate-warning {
  color: #E6A23C;
}

.health-rate.rate-danger {
  color: #F56C6C;
}

.status-code-tag {
  font-weight: 600;
}

.metric-divider {
  width: 1px;
  height: 24px;
  background: rgba(0, 0, 0, 0.1);
}

.loading-text {
  opacity: 0.6;
}

/* 分页区域 */
.pagination-section {
  display: flex;
  justify-content: center;
  padding-top: 12px;
}

:deep(.el-pagination) {
  justify-content: center;
}

:deep(.el-empty) {
  padding: 60px 0;
}

/* 下拉菜单样式优化 */
:deep(.el-dropdown-menu) {
  border-radius: 8px;
  padding: 6px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

:deep(.el-dropdown-menu__item) {
  border-radius: 4px;
  padding: 8px 16px;
  transition: all 0.2s ease;
}

:deep(.el-dropdown-menu__item:hover) {
  background: rgba(102, 126, 234, 0.08);
}

:deep(.el-dropdown-menu__item.is-divided) {
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  margin-top: 6px;
  padding-top: 8px;
}

:deep(.el-dropdown-menu__item:last-child:hover) {
  background: rgba(245, 108, 108, 0.1);
  color: #F56C6C;
}

/* 现代化按钮风格 */
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

/* 滚动条美化 */
.tree-content::-webkit-scrollbar {
  width: 6px;
}

.tree-content::-webkit-scrollbar-track {
  background: transparent;
}

.tree-content::-webkit-scrollbar-thumb {
  background: #d9d9d9;
  border-radius: 3px;
}

.tree-content::-webkit-scrollbar-thumb:hover {
  background: #bfbfbf;
}

/* ==================== 分层监控样式 ==================== */

/* 子监控统计行 */
.child-stats-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 10px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05), rgba(118, 75, 162, 0.05));
  border-radius: 6px;
  margin-top: 4px;
}

.child-stats-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
}

.child-stats-info .stats-label {
  color: #606266;
  font-weight: 500;
}

.child-stats-info .stats-value {
  font-weight: 600;
  color: #2c3e50;
  display: flex;
  align-items: center;
  gap: 3px;
}

.child-stats-info .healthy-count {
  color: #67C23A;
  font-weight: 700;
  font-size: 14px;
}

.child-stats-info .divider {
  color: #909399;
  font-size: 12px;
  font-weight: 400;
}

.child-stats-info .total-count {
  color: #909399;
  font-size: 14px;
  font-weight: 600;
}

.child-stats-info .stats-text {
  color: #606266;
  font-size: 11px;
  font-weight: 500;
  margin-left: 2px;
}

.view-detail-icon {
  font-size: 16px;
  color: #409eff;
  transition: transform 0.3s ease;
}

.api-card:hover .view-detail-icon {
  transform: translateX(4px);
}

/* ==================== 抽屉样式 ==================== */

.drawer-content {
  padding: 0;
  min-height: 400px;
}

.drawer-section {
  margin-bottom: 24px;
}

.drawer-section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 16px;
  padding: 0 0 12px 0;
  border-bottom: 2px solid rgba(102, 126, 234, 0.2);
}

.drawer-section-title .el-icon {
  color: #667eea;
}

/* 父监控详情卡片 */
.detail-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.detail-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

.detail-card-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.detail-card-title h4 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.detail-card-status {
  display: flex;
  align-items: center;
}

.status-icon-large {
  width: 32px;
  height: 32px;
}

.detail-card-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.detail-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.detail-label {
  min-width: 100px;
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}

.detail-value {
  flex: 1;
  font-size: 13px;
  color: #2c3e50;
  word-break: break-all;
}

.url-value {
  font-family: 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  background: rgba(0, 0, 0, 0.03);
  padding: 4px 8px;
  border-radius: 4px;
}

/* 子监控卡片列表 */
.child-cards-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.child-detail-card {
  background: white;
  border-radius: 10px;
  padding: 14px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
}

.child-detail-card:hover {
  transform: translateX(4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.child-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.child-card-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-icon-small {
  width: 20px;
  height: 20px;
}

.child-card-title h5 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
}

.child-card-actions {
  display: flex;
  gap: 4px;
}

.child-card-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.child-detail-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 6px;
}

.child-url {
  flex: 1;
  font-size: 12px;
  color: #606266;
  font-family: 'Monaco', 'Courier New', monospace;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.child-metrics {
  display: flex;
  gap: 16px;
  padding: 8px 0;
}

.child-metric-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.child-metric-item .metric-label {
  font-size: 11px;
  color: #909399;
  font-weight: 500;
}

.child-metric-item .metric-value {
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
}

.child-error-message {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  padding: 8px;
  background: rgba(245, 108, 108, 0.1);
  border-radius: 6px;
  border-left: 3px solid #F56C6C;
}

.child-error-message span {
  flex: 1;
  font-size: 12px;
  color: #F56C6C;
  line-height: 1.4;
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

</style>
