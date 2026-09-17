<script setup>
import { ref, onMounted, onUnmounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Edit,
  Delete,
  View,
  DeleteFilled,
  Eye,
  Setting
} from '@element-plus/icons-vue'
import { policyApi, policyRuleApi } from '@/api/alert'
import cmdbApi from '@/api/cmdb'
import { robotApi } from '@/api/monitor'
import healthIcon from '@/assets/image/健康.svg'
import fireIcon from '@/assets/image/起火.svg'
import dingtalkIcon from '@/assets/image/钉钉01.svg'
import feishuIcon from '@/assets/image/飞书.svg'
import wechatIcon from '@/assets/image/企业微信.svg'
import emailIcon from '@/assets/image/邮箱.svg'
import webhookIcon from '@/assets/image/url.svg'
import teamsIcon from '@/assets/image/通知.svg'

const router = useRouter()

// 定义事件
const emit = defineEmits(['refresh-stats'])

// 搜索表单
const searchForm = reactive({
  name: '',
  businessGroup: '',
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
const dialogTitle = ref('新增策略')
const formLoading = ref(false)

// 表单数据
const policyForm = reactive({
  id: null,
  name: '',
  description: '',
  groupId: '',
  groupName: '',
  robotIds: [],
  silencePeriod: 300,
  repeatInterval: 600,
  maxNotifyCount: 10,
  sendRecovery: true,
  appendTags: '[]',
  status: 1
})

const formRef = ref(null)
const statusOptions = [
  { label: '全部', value: -1 },
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

// 表单验证规则
const formRules = {
  name: [{ required: true, message: '请输入策略名称', trigger: 'blur' }],
  groupId: [{ required: true, message: '请选择业务组', trigger: 'change' }]
}

// 模板列表
const businessGroupList = ref([])
const businessGroupMap = ref({}) // 用于快速查询业务组名称
const robotList = ref([])
const robotMap = ref({}) // 用于快速查询机器人完整信息

// 自动刷新定时器
const autoRefreshTimer = ref(null)

// PromQL 预览
const previewDialogVisible = ref(false)
const previewLoading = ref(false)
const previewForm = reactive({
  promql: '',
  timestamp: Math.floor(Date.now() / 1000)
})
const previewResult = ref(null)

// 获取业务组列表：与 CMDB 分组树保持一致，取根分组下的业务子分组
const fetchBusinessGroups = async () => {
  try {
    const { data: res } = await cmdbApi.getAllCmdbGroups()
    if (res.code === 200 && res.data) {
      businessGroupMap.value = {}

      const groupList = Array.isArray(res.data) ? res.data : []
      const businessGroups = []

      const collectBusinessGroups = (groups, parentId) => {
        if (!Array.isArray(groups)) return

        groups.forEach(group => {
          const groupItem = {
            id: group.id,
            name: group.name,
            parentId: group.parentId || parentId
          }
          businessGroups.push(groupItem)
          businessGroupMap.value[group.id] = group.name

          if (group.children && group.children.length > 0) {
            collectBusinessGroups(group.children, group.id)
          }
        })
      }

      groupList.forEach(group => {
        collectBusinessGroups(group.children, group.id)
      })
      businessGroupList.value = businessGroups
    }
  } catch (error) {
    console.error('获取业务组列表失败:', error)
  }
}

// 获取机器人列表
const fetchRobots = async () => {
  try {
    const { data: res } = await robotApi.getList({
      page: 1,
      pageSize: 100,
      status: 1  // 只获取启用的机器人
    })
    if (res.code === 200) {
      const robots = res.data.list || []
      robotList.value = robots
      // 构建机器人 Map，存储完整的机器人对象（包含type）
      robots.forEach(robot => {
        robotMap.value[robot.id] = {
          id: robot.id,
          name: robot.name,
          type: robot.type || 'webhook'
        }
      })
    }
  } catch (error) {
    console.error('获取机器人列表失败:', error)
  }
}

// 根据 groupId 获取业务组名称
const getGroupName = (groupId) => {
  if (!groupId) return ''
  const group = businessGroupList.value.find(g => g.id === groupId)
  return group ? group.name : `分组${groupId}`
}

// 获取策略列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = { ...searchForm }
    if (params.status === -1) delete params.status
    if (!params.businessGroup) delete params.businessGroup

    const { data: res } = await policyApi.getList(params)
    if (res.code === 200) {
      // 根据 groupId 查询业务组名称
      const list = res.data.list || []
      list.forEach(item => {
        // 使用 getGroupName 函数根据 groupId 获取最新的业务组名称
        item.groupName = getGroupName(item.groupId)
        // 初始化规则数量为 0
        item.ruleCount = 0
        // 初始化实时状态，默认为 healthy
        item.healthStatus = item.healthStatus || 'healthy'

        // 解析 robotIds 并转换为机器人对象数组
        try {
          const robotIds = item.robotIds ? JSON.parse(item.robotIds) : []
          item.robots = robotIds.map(id => {
            const robot = robotMap.value[id]
            return {
              id: id,
              name: robot?.name || `机器人${id}`,
              type: robot?.type || 'webhook'
            }
          })
        } catch (e) {
          item.robots = []
        }
      })
      tableData.value = list
      total.value = res.data.total || 0

      // 为每个策略获取规则数量
      const policyRuleCounts = await Promise.all(
        list.map(policy =>
          policyRuleApi.getList({ policyId: policy.id, pageNum: 1, pageSize: 1 })
            .then(({ data: ruleRes }) => ({
              policyId: policy.id,
              count: ruleRes.data?.totalCount || ruleRes.data?.rules?.length || 0
            }))
            .catch(() => ({ policyId: policy.id, count: 0 }))
        )
      )

      // 更新每个策略的规则数量
      policyRuleCounts.forEach(({ policyId, count }) => {
        const policy = tableData.value.find(p => p.id === policyId)
        if (policy) {
          policy.ruleCount = count
        }
      })
    }
  } catch (error) {
    ElMessage.error('获取策略列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  searchForm.pageNum = 1
  fetchList()
}

// 重置
const handleReset = () => {
  searchForm.name = ''
  searchForm.businessGroup = ''
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
const handleAdd = async () => {
  // 打开新增对话框前，刷新业务组和机器人列表，保证最新
  await fetchBusinessGroups()
  await fetchRobots()

  dialogTitle.value = '新增策略'
  policyForm.id = null
  policyForm.name = ''
  policyForm.description = ''
  policyForm.groupId = ''
  policyForm.groupName = ''
  policyForm.robotIds = []
  policyForm.silencePeriod = 300
  policyForm.repeatInterval = 600
  policyForm.maxNotifyCount = 10
  policyForm.sendRecovery = true
  policyForm.appendTags = '[]'
  policyForm.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = async (row) => {
  try {
    const { data: res } = await policyApi.getDetail(row.id)
    if (res.code === 200) {
      dialogTitle.value = '编辑策略'
      const detail = res.data
      policyForm.id = detail.id
      policyForm.name = detail.name
      policyForm.description = detail.description || ''
      policyForm.groupId = detail.groupId || ''
      // 根据 groupId 查询业务组名称
      policyForm.groupName = getGroupName(detail.groupId)
      // 解析 robotIds JSON 字符串
      try {
        policyForm.robotIds = detail.robotIds ? JSON.parse(detail.robotIds) : []
      } catch (e) {
        policyForm.robotIds = []
      }
      policyForm.silencePeriod = detail.silencePeriod || 300
      policyForm.repeatInterval = detail.repeatInterval || 600
      policyForm.maxNotifyCount = detail.maxNotifyCount || 10
      policyForm.sendRecovery = detail.sendRecovery !== false
      policyForm.appendTags = detail.appendTags || '[]'
      policyForm.status = detail.status
      dialogVisible.value = true
    } else {
      ElMessage.error(res.message || '获取策略详情失败')
    }
  } catch (error) {
    ElMessage.error('获取策略详情失败')
    console.error(error)
  }
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该策略吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const { data: res } = await policyApi.delete(row.id)
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
    ElMessage.warning('请先选择要删除的策略')
    return
  }

  try {
    await ElMessageBox.confirm(`确定要删除 ${selectedRows.value.length} 个策略吗?`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const ids = selectedRows.value.map(row => row.id)
    const { data: res } = await policyApi.batchDelete(ids)
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

// 配置告警规则
const handleConfigureRule = (row) => {
  // 跳转到告警规则页面，传递策略ID
  router.push({
    path: '/monitor/alert-rule',
    query: {
      policyId: row.id,
      policyName: row.name
    }
  })
}

// 状态切换
const handleStatusChange = async (row) => {
  try {
    const { data: res } = await policyApi.updateStatus(row.id, row.status)
    if (res.code === 200) {
      ElMessage.success(row.status === 1 ? '启用成功' : '禁用成功')
      emit('refresh-stats')
    } else {
      row.status = row.status === 1 ? 0 : 1
      ElMessage.error(res.message || '操作失败')
    }
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('操作失败')
    console.error(error)
  }
}

// PromQL 预览
const handlePreview = async () => {
  if (!previewForm.promql) {
    ElMessage.warning('请输入PromQL表达式')
    return
  }

  previewLoading.value = true
  try {
    const { data: res } = await policyApi.preview({
      promql: previewForm.promql,
      timestamp: previewForm.timestamp
    })
    if (res.code === 200) {
      // 转换数据格式：提取value数组中的实际值，移除时间戳
      const formattedData = (res.data.results || []).map(item => ({
        metric: JSON.stringify(item.metric),
        value: item.value ? item.value[1] : 'N/A' // 提取值，移除时间戳
      }))
      previewResult.value = formattedData
    } else {
      ElMessage.error(res.message || 'PromQL预览失败')
    }
  } catch (error) {
    ElMessage.error('PromQL预览失败')
    console.error(error)
  } finally {
    previewLoading.value = false
  }
}

// 打开PromQL预览对话框
const openPreviewDialog = () => {
  previewForm.promql = ''
  previewForm.timestamp = Math.floor(Date.now() / 1000)
  previewResult.value = null
  previewDialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    formLoading.value = true
    try {
      const data = { ...policyForm }
      // 处理JSON字段
      if (Array.isArray(data.robotIds)) {
        data.robotIds = JSON.stringify(data.robotIds)
      }

      let result
      if (data.id) {
        result = await policyApi.update(data)
      } else {
        delete data.id
        result = await policyApi.create(data)
      }

      const { data: res } = result
      if (res.code === 200) {
        ElMessage.success(data.id ? '更新成功' : '创建成功')
        dialogVisible.value = false
        // 重新获取业务组和机器人列表，保证列表最新
        await fetchBusinessGroups()
        await fetchRobots()
        await fetchList()
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

// 获取模板名称
const getTemplateName = (id) => {
  const template = templateList.value.find(t => t.id === id)
  return template ? template.name : `模板${id}`
}

// 获取机器人类型标签
const getRobotTypeLabel = (type) => {
  const typeMap = {
    'dingtalk': '钉钉',
    'feishu': '飞书',
    'wechat': '企业微信',
    'email': '邮箱',
    'webhook': 'Webhook',
    'teams': 'Teams'
  }
  return typeMap[type] || '未知'
}

// 获取机器人类型标签样式
const getRobotTypeTagType = (type) => {
  const tagMap = {
    'dingtalk': 'primary',
    'feishu': 'success',
    'wechat': 'warning',
    'email': 'info',
    'webhook': 'danger',
    'teams': 'primary'
  }
  return tagMap[type] || 'info'
}

// 获取机器人图标
const getRobotIcon = (type) => {
  const iconMap = {
    'dingtalk': dingtalkIcon,
    'feishu': feishuIcon,
    'wechat': wechatIcon,
    'email': emailIcon,
    'webhook': webhookIcon,
    'teams': teamsIcon
  }
  return iconMap[type] || ''
}

// 行选择
const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

onMounted(async () => {
  // 先获取业务组列表，再获取机器人列表和策略列表
  await fetchBusinessGroups()
  await fetchRobots()
  await fetchList()

  // 启动自动刷新定时器（每60秒刷新一次状态）
  autoRefreshTimer.value = setInterval(() => {
    fetchList()
  }, 60000) // 60000ms = 60秒
})

// 组件卸载时清除定时器
onUnmounted(() => {
  if (autoRefreshTimer.value) {
    clearInterval(autoRefreshTimer.value)
  }
})
</script>

<template>
  <div class="alert-policy-container">
    <!-- 搜索区域 -->
    <div class="search-section">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="策略名称">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入策略名称"
            size="small"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>

        <el-form-item label="业务组">
          <el-select
            v-model="searchForm.businessGroup"
            placeholder="全部业务组"
            size="small"
            clearable
            style="width: 150px"
          >
            <el-option
              v-for="item in businessGroupList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
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
        v-authority="['monitor:alert:policy:add']"
        class="modern-btn success-btn"
      >
        <el-icon><Plus /></el-icon>
        新增策略
      </el-button>

      <el-button
        type="danger"
        size="small"
        @click="handleBatchDelete"
        :disabled="selectedRows.length === 0"
        v-authority="['monitor:alert:policy:delete']"
        class="modern-btn danger-btn"
      >
        <el-icon><DeleteFilled /></el-icon>
        批量删除
      </el-button>

      <el-button
        type="info"
        size="small"
        @click="openPreviewDialog"
        v-authority="['monitor:alert:policy:preview']"
        class="modern-btn info-btn"
      >
        <el-icon><Eye /></el-icon>
        PromQL预览
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

        <el-table-column prop="name" label="策略名称" min-width="120">
          <template #default="{ row }">
            <span style="font-weight: 400">{{ row.name }}</span>
          </template>
        </el-table-column>

        <el-table-column label="实时状态" width="120" align="center">
          <template #default="{ row }">
            <div style="display: flex; align-items: center; justify-content: center; gap: 8px">
              <img
                :src="row.healthStatus === 'healthy' ? healthIcon : fireIcon"
                :alt="row.healthStatus"
                style="width: 20px; height: 20px"
              />
              <span style="font-size: 12px; font-weight: 500">
                {{ row.healthStatus === 'healthy' ? '健康' : '告警中' }}
              </span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="groupName" label="业务组" width="120">
          <template #default="{ row }">
            <span>{{ row.groupName }}</span>
          </template>
        </el-table-column>

        <el-table-column label="关联规则数量" width="120" align="center">
          <template #default="{ row }">
            <el-tag type="info" size="small">
              {{ row.ruleCount || 0 }} 个
            </el-tag>
          </template>
        </el-table-column>



        <el-table-column label="通知机器人" min-width="200">
          <template #default="{ row }">
            <div v-if="row.robots && row.robots.length > 0" style="display: flex; gap: 8px; flex-wrap: wrap">
              <div v-for="robot in row.robots" :key="robot.id" style="display: flex; align-items: center; gap: 6px; padding: 4px 8px; background: #f5f7fa; border-radius: 4px;">
                <img :src="getRobotIcon(robot.type)" :alt="robot.type" style="width: 16px; height: 16px;" />
                <span style="font-size: 12px; color: #606266;">{{ robot.name }}</span>
              </div>
            </div>
            <span v-else style="color: #c0c4cc">-</span>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>

        <el-table-column prop="createdAt" label="创建时间" width="170" align="center">
          <template #default="{ row }">
            <span style="font-size: 12px">{{ row.createdAt }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" fixed="right" width="300" align="center">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-tooltip content="编辑" placement="top">
                <el-button
                  type="warning"
                  :icon="Edit"
                  size="small"
                  @click="handleEdit(row)"
                  v-authority="['monitor:alert:policy:edit']"
                >编辑</el-button>
              </el-tooltip>
              <el-tooltip content="配置告警规则" placement="top">
                <el-button
                  type="primary"
                  :icon="Setting"
                  size="small"
                  @click="handleConfigureRule(row)"
                  v-authority="['monitor:alert:policy:rule']"
                >配置规则</el-button>
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button
                  type="danger"
                  :icon="Delete"
                  size="small"
                  @click="handleDelete(row)"
                  v-authority="['monitor:alert:policy:delete']"
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

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="800px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="policyForm"
        :rules="formRules"
        label-width="140px"
      >
        <!-- 基础信息 -->
        <el-form-item label="策略名称" prop="name">
          <el-input v-model="policyForm.name" placeholder="请输入策略名称" />
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="policyForm.description" type="textarea" :rows="3" placeholder="请输入策略描述" />
        </el-form-item>

        <el-form-item label="业务组" prop="groupId">
          <el-select v-model="policyForm.groupId" placeholder="请选择业务组" style="width: 100%">
            <el-option
              v-for="item in businessGroupList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <!-- 通知配置 -->
        <el-form-item label="通知机器人">
          <el-select
            v-model="policyForm.robotIds"
            placeholder="请选择通知机器人"
            multiple
            style="width: 100%"
          >
            <el-option
              v-for="item in robotList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="静默周期（秒）">
          <el-input v-model.number="policyForm.silencePeriod" type="number" placeholder="默认300秒" />
        </el-form-item>

        <el-form-item label="重复通知间隔（秒）">
          <el-input v-model.number="policyForm.repeatInterval" type="number" placeholder="默认600秒" />
        </el-form-item>

        <el-form-item label="最大通知次数">
          <el-input v-model.number="policyForm.maxNotifyCount" type="number" placeholder="默认10次" />
        </el-form-item>

        <el-form-item label="发送恢复通知">
          <el-radio-group v-model="policyForm.sendRecovery">
            <el-radio :label="true">是</el-radio>
            <el-radio :label="false">否</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 高级配置 -->
        <el-form-item label="附加标签">
          <el-input v-model="policyForm.appendTags" type="textarea" :rows="3" placeholder='JSON格式，如: [{"key":"env","value":"prod"}]' />
        </el-form-item>

        <el-form-item label="状态">
          <el-radio-group v-model="policyForm.status">
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

    <!-- PromQL预览对话框 -->
    <el-dialog
      v-model="previewDialogVisible"
      title="PromQL预览"
      width="700px"
      :close-on-click-modal="false"
    >
      <el-form label-width="120px">
        <el-form-item label="PromQL表达式">
          <el-input
            v-model="previewForm.promql"
            type="textarea"
            :rows="4"
            placeholder="请输入PromQL表达式"
          />
        </el-form-item>
      </el-form>

      <div v-if="previewResult" style="margin-top: 20px">
        <el-divider />
        <h4>查询结果（共 {{ previewResult.length }} 条）</h4>
        <el-table :data="previewResult" size="small" style="width: 100%" max-height="400px">
          <el-table-column prop="metric" label="指标标签" min-width="300" show-overflow-tooltip />
          <el-table-column prop="value" label="值" width="120" align="right" />
        </el-table>
      </div>

      <template #footer>
        <el-button @click="previewDialogVisible = false">关闭</el-button>
        <el-button type="primary" :loading="previewLoading" @click="handlePreview">
          预览
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.alert-policy-container {
  padding: 0;
}

.search-section {
  margin-bottom: 16px;
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
  margin-bottom: 16px;
  display: flex;
  gap: 10px;
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

:deep(.info-btn) {
  background: linear-gradient(45deg, #909399, #A6ADB7);
  color: white;
}

/* 实时状态样式 */
:deep(.el-tag.success) {
  --el-tag-text-color: #67C23A;
  --el-tag-bg-color: #F0F9FF;
  --el-tag-border-color: #67C23A;
}

:deep(.el-tag.danger) {
  --el-tag-text-color: #F56C6C;
  --el-tag-bg-color: #FEF0F0;
  --el-tag-border-color: #F56C6C;
}
</style>
