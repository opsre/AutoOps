<script setup>
import { ref, onMounted, onUnmounted, reactive, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Edit,
  Delete,
  View,
  DeleteFilled,
  Back,
  DataAnalysis
} from '@element-plus/icons-vue'
import { policyRuleApi, templateApi } from '@/api/alert'
import healthIcon from '@/assets/image/健康.svg'
import fireIcon from '@/assets/image/起火.svg'
import severityIcon from '@/assets/image/严重.svg'
import warningIcon from '@/assets/image/警告.svg'
import tipsIcon from '@/assets/image/提示.svg'

const route = useRoute()
const router = useRouter()

// 获取路由参数
const policyId = ref(route.query.policyId)
const policyName = ref(route.query.policyName)

// 策略信息（从 API 获取）
const policyInfo = ref({
  id: null,
  name: '',
  groupId: null,
  groupName: '',
  ruleCount: 0,
  status: 1
})

// 左侧分类选中状态
const selectedCategory = ref('all')

// 搜索表单
const searchForm = reactive({
  name: '',
  status: -1,
  category: '',
  pageNum: 1,
  pageSize: 10
})

// 列表数据
const tableData = ref([])
const total = ref(0)
const loading = ref(false)
const selectedRows = ref([])

// 告警模板列表
const templateList = ref([])

// 规则分类列表（对应后端API的category参数）
const categoryList = ref([
  { label: '全部规则', value: 'all', count: 0 },
  { label: 'CPU', value: 'cpu', count: 0 },
  { label: '内存', value: 'memory', count: 0 },
  { label: '负载', value: 'load', count: 0 },
  { label: '磁盘', value: 'disk', count: 0 },
  { label: 'IO', value: 'io', count: 0 },
  { label: '网络', value: 'network', count: 0 },
  { label: '连接数', value: 'connection', count: 0 },
  { label: '服务', value: 'service', count: 0 },
  { label: '自定义', value: 'custom', count: 0 }
])

// 获取各分类的规则数量
const fetchCategoryCounts = async () => {
  try {
    // 并发请求所有分类的数量
    const categories = ['all', 'cpu', 'memory', 'load', 'disk', 'io', 'network', 'connection', 'service', 'custom']
    const promises = categories.map(async (cat) => {
      const params = {
        policyId: policyId.value,
        pageNum: 1,
        pageSize: 1 // 只获取总数，不需要数据
      }
      if (cat !== 'all') {
        params.category = cat
      }
      const { data: res } = await policyRuleApi.getList(params)
      return {
        category: cat,
        count: res.data?.totalCount || 0
      }
    })

    const results = await Promise.all(promises)
    results.forEach(result => {
      const cat = categoryList.value.find(c => c.value === result.category)
      if (cat) {
        cat.count = result.count
      }
    })
  } catch (error) {
    console.error('获取分类统计失败:', error)
  }
}

// 自动刷新定时器
const autoRefreshTimer = ref(null)

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('新增规则')
const formLoading = ref(false)

// 预览对话框
const previewDialogVisible = ref(false)
const previewLoading = ref(false)
const previewData = ref(null)

// 表单数据
const ruleForm = reactive({
  id: null,
  policyId: Number(policyId.value),
  name: '',
  description: '',
  templateId: null,
  templateName: '',
  operator: '>=',
  threshold: null,
  severity: 2,
  duration: 30,
  continuousTimes: 1,
  recoverDuration: 60,
  status: 1
})

const formRef = ref(null)

const statusOptions = [
  { label: '全部', value: -1 },
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

const comparisonOptions = [
  { label: '>', value: '>' },
  { label: '<', value: '<' },
  { label: '>=', value: '>=' },
  { label: '<=', value: '<=' },
  { label: '==', value: '==' },
  { label: '!=', value: '!=' }
]

const severityOptions = [
  { label: '严重', value: 1 },
  { label: '警告', value: 2 },
  { label: '提示', value: 3 }
]

// 表单验证规则
const formRules = {
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  templateId: [{ required: true, message: '请选择告警模板', trigger: 'change' }],
  operator: [{ required: true, message: '请选择比较操作符', trigger: 'change' }],
  threshold: [{ required: true, message: '请输入阈值', trigger: 'blur' }],
  severity: [{ required: true, message: '请选择严重级别', trigger: 'change' }]
}

// 获取规则列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = {
      ...searchForm,
      policyId: policyId.value
    }
    if (params.status === -1) delete params.status
    // 如果选中了分类，添加category参数
    if (selectedCategory.value !== 'all') {
      params.category = selectedCategory.value
    }

    const { data: res } = await policyRuleApi.getList(params)
    if (res.code === 200) {
      // 保存策略信息（包含最新的业务组信息）
      if (res.data.policyInfo) {
        policyInfo.value = res.data.policyInfo
      }

      // API 返回的是 rules 数组，而不是 list
      const rules = res.data.rules || []
      // 映射字段名：ruleName -> name，去掉策略名称前缀
      tableData.value = rules.map(rule => {
        let displayName = rule.ruleName || rule.name
        // 去掉 "策略名称-" 前缀，只保留规则名称
        if (displayName && displayName.includes('-')) {
          const parts = displayName.split('-')
          displayName = parts.slice(1).join('-') || parts[0]
        }
        return {
          ...rule,
          name: displayName,
          templateName: rule.metricKey || '',
          // 初始化实时状态，默认为 healthy
          healthStatus: rule.healthStatus || 'healthy'
        }
      })
      total.value = res.data.totalCount || rules.length

      // 更新分类统计
      await fetchCategoryCounts()
    }
  } catch (error) {
    ElMessage.error('获取规则列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 选择分类
const handleCategorySelect = (category) => {
  selectedCategory.value = category
  searchForm.pageNum = 1 // 重置到第一页
  fetchList() // 重新获取数据
}

// 获取告警模板列表（支持按分类过滤）
const fetchTemplates = async (category = '') => {
  try {
    const params = {
      pageNum: 1,
      pageSize: 100,
      status: 1
    }
    // 如果有分类且不是"全部"，则添加category参数
    if (category && category !== 'all') {
      params.category = category
    }

    const { data: res } = await templateApi.getList(params)
    if (res.code === 200) {
      templateList.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取告警模板失败:', error)
  }
}

// 选择模板时的回调
const handleTemplateChange = (templateId) => {
  const selected = templateList.value.find(t => t.id === templateId)
  if (selected) {
    ruleForm.templateName = selected.name
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
  dialogTitle.value = '新增规则'
  ruleForm.id = null
  ruleForm.policyId = Number(policyId.value)
  ruleForm.name = ''
  ruleForm.description = ''
  ruleForm.templateId = null
  ruleForm.templateName = ''
  ruleForm.threshold = null
  ruleForm.operator = '>='
  ruleForm.duration = 30
  ruleForm.continuousTimes = 1
  ruleForm.recoverDuration = 60
  ruleForm.severity = 2
  ruleForm.status = 1

  // 根据当前选中的分类加载对应的告警模板
  await fetchTemplates(selectedCategory.value)

  dialogVisible.value = true
}

// 编辑
const handleEdit = async (row) => {
  try {
    const { data: res } = await policyRuleApi.getDetail(row.id)
    if (res.code === 200) {
      dialogTitle.value = '编辑规则'
      const detail = res.data
      ruleForm.id = detail.id
      ruleForm.policyId = Number(detail.policyId)
      ruleForm.name = detail.ruleName || ''
      ruleForm.description = detail.description || ''
      ruleForm.templateId = Number(detail.templateId)
      ruleForm.templateName = detail.metricKey || ''
      ruleForm.threshold = Number(detail.threshold)
      ruleForm.operator = detail.operator || '>='
      ruleForm.duration = Number(detail.duration) || 60
      ruleForm.continuousTimes = Number(detail.continuousTimes) || 1
      ruleForm.recoverDuration = Number(detail.recoverDuration) || 60
      ruleForm.severity = Number(detail.severity) || 2
      ruleForm.status = Number(detail.status)

      // 编辑时加载全部模板（不根据分类过滤）
      await fetchTemplates()

      dialogVisible.value = true
    } else {
      ElMessage.error(res.message || '获取规则详情失败')
    }
  } catch (error) {
    ElMessage.error('获取规则详情失败')
    console.error(error)
  }
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该规则吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const { data: res } = await policyRuleApi.delete(row.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      fetchList()
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
    ElMessage.warning('请先选择要删除的规则')
    return
  }

  try {
    await ElMessageBox.confirm(`确定要删除 ${selectedRows.value.length} 个规则吗?`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const ids = selectedRows.value.map(row => row.id)
    const { data: res } = await policyRuleApi.batchDelete(ids)
    if (res.code === 200) {
      ElMessage.success('批量删除成功')
      selectedRows.value = []
      fetchList()
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

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    formLoading.value = true
    try {
      // 构建提交数据，确保所有数字字段都是 number 类型
      const data = {
        ruleName: ruleForm.name,
        description: ruleForm.description,
        policyId: Number(ruleForm.policyId),
        templateId: Number(ruleForm.templateId),
        operator: ruleForm.operator,
        threshold: Number(ruleForm.threshold),
        severity: Number(ruleForm.severity),
        duration: Number(ruleForm.duration),
        continuousTimes: Number(ruleForm.continuousTimes),
        recoverDuration: Number(ruleForm.recoverDuration),
        status: Number(ruleForm.status)
      }

      if (ruleForm.id) {
        data.id = ruleForm.id
      }

      let result
      if (ruleForm.id) {
        result = await policyRuleApi.update(data)
      } else {
        result = await policyRuleApi.create(data)
      }

      const { data: res } = result
      if (res.code === 200) {
        ElMessage.success(ruleForm.id ? '更新成功' : '创建成功')
        dialogVisible.value = false
        fetchList()
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

// 行选择
const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

// 返回
const handleBack = () => {
  router.back()
}

// 处理状态变更
const handleStatusChange = async (row) => {
  try {
    const { data: res } = await policyRuleApi.updateStatus(row.id, row.status)
    if (res.code === 200) {
      ElMessage.success(row.status === 1 ? '规则已启用' : '规则已禁用')
    } else {
      row.status = row.status === 1 ? 0 : 1 // 恢复状态
      ElMessage.error(res.message || '状态更新失败')
    }
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1 // 恢复状态
    ElMessage.error('状态更新失败')
    console.error(error)
  }
}

// 预览规则数据
const handlePreview = async (row) => {
  previewLoading.value = true
  try {
    const { data: res } = await policyRuleApi.preview(row.id)
    if (res.code === 200) {
      previewData.value = res.data
      previewDialogVisible.value = true
    } else {
      ElMessage.error(res.message || '获取预览数据失败')
    }
  } catch (error) {
    ElMessage.error('获取预览数据失败')
    console.error(error)
  } finally {
    previewLoading.value = false
  }
}

onMounted(async () => {
  await fetchTemplates()
  await fetchList()
  await fetchCategoryCounts()

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
  <div class="alert-policy-rule-container">
    <!-- 面包屑 -->
    <div class="breadcrumb-section">
      <el-button
        type="text"
        size="small"
        @click="handleBack"
        style="margin-right: 10px; padding: 0"
      >
        <el-icon><Back /></el-icon>
        返回
      </el-button>
      <span class="policy-info">
        > 告警策略: <strong>{{ policyInfo.name }}</strong> (ID: {{ policyInfo.id }})
        <span style="margin-left: 10px; color: #909399;">| 业务组: <strong style="color: #667eea;">{{ policyInfo.groupName }}</strong></span>
      </span>
    </div>

    <!-- 主内容区域：左右布局 -->
    <div class="main-content">
      <!-- 左侧分类栏 -->
      <div class="category-sidebar">
        <div class="category-header">规则分类</div>
        <div class="category-list">
          <div
            v-for="category in categoryList"
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
            <el-form-item label="规则名称">
              <el-input
                v-model="searchForm.name"
                placeholder="请输入规则名称"
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
            v-authority="['monitor:alert:policy-rule:add']"
            class="modern-btn success-btn"
          >
            <el-icon><Plus /></el-icon>
            新增规则
          </el-button>

          <el-button
            type="danger"
            size="small"
            @click="handleBatchDelete"
            :disabled="selectedRows.length === 0"
            v-authority="['monitor:alert:policy-rule:delete']"
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

            <el-table-column prop="name" label="名称" min-width="100">
              <template #default="{ row }">
                <span style="font-weight: 500">{{ row.name }}</span>
              </template>
            </el-table-column>

            <el-table-column label="实时状态" width="100" align="center">
              <template #default="{ row }">
                <div style="display: flex; align-items: center; justify-content: center; gap: 6px">
                  <img
                    :src="row.healthStatus === 'healthy' ? healthIcon : fireIcon"
                    :alt="row.healthStatus"
                    style="width: 18px; height: 18px"
                  />
                  <span style="font-size: 12px; font-weight: 500">
                    {{ row.healthStatus === 'healthy' ? '健康' : '告警中' }}
                  </span>
                </div>
              </template>
            </el-table-column>

            <el-table-column prop="templateName" label="告警模板" min-width="120">
              <template #default="{ row }">
                <span>{{ row.templateName }}</span>
              </template>
            </el-table-column>

            <el-table-column prop="operator" label="操作符" width="80" align="center">
              <template #default="{ row }">
                <el-tag type="info" size="small">{{ row.operator }}</el-tag>
              </template>
            </el-table-column>

            <el-table-column prop="threshold" label="阈值" width="100" align="center">
              <template #default="{ row }">
                <span>{{ row.threshold }}</span>
              </template>
            </el-table-column>

            <el-table-column prop="duration" label="时长(秒)" width="100" align="center">
              <template #default="{ row }">
                <span>{{ row.duration }}</span>
              </template>
            </el-table-column>

            <el-table-column prop="severity" label="严重级别" width="100" align="center">
              <template #default="{ row }">
                <div style="display: flex; align-items: center; justify-content: center; gap: 6px">
                  <img
                    :src="row.severity === 1 ? severityIcon : row.severity === 2 ? warningIcon : tipsIcon"
                    :alt="row.severity === 1 ? '严重' : row.severity === 2 ? '警告' : '提示'"
                    style="width: 18px; height: 18px"
                  />
                  <span style="font-size: 12px; font-weight: 500">
                    {{ row.severity === 1 ? '严重' : row.severity === 2 ? '警告' : '提示' }}
                  </span>
                </div>
              </template>
            </el-table-column>

            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-switch
                  v-model="row.status"
                  :active-value="1"
                  :inactive-value="0"
                  @change="handleStatusChange(row)"
                  v-authority="['monitor:alert:policy-rule:status']"
                />
              </template>
            </el-table-column>

            <el-table-column label="操作" fixed="right" width="230" align="center">
              <template #default="{ row }">
                <div class="operation-buttons">
                  <el-tooltip content="测试" placement="top">
                    <el-button
                      type="primary"
                      :icon="DataAnalysis"
                      size="small"
                      @click="handlePreview(row)"
                    >测试</el-button>
                  </el-tooltip>
                  <el-tooltip content="编辑" placement="top">
                    <el-button
                      type="warning"
                      :icon="Edit"
                      size="small"
                      @click="handleEdit(row)"
                      v-authority="['monitor:alert:policy-rule:edit']"
                    >编辑</el-button>
                  </el-tooltip>
                  <el-tooltip content="删除" placement="top">
                    <el-button
                      type="danger"
                      :icon="Delete"
                      size="small"
                      @click="handleDelete(row)"
                      v-authority="['monitor:alert:policy-rule:delete']"
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
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="ruleForm"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item label="规则名称" prop="name">
          <el-input v-model="ruleForm.name" placeholder="请输入规则名称" />
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="ruleForm.description" type="textarea" :rows="2" placeholder="请输入规则描述" />
        </el-form-item>

        <el-form-item label="告警模板" prop="templateId">
          <el-select
            v-model="ruleForm.templateId"
            placeholder="请选择告警模板"
            style="width: 100%"
            @change="handleTemplateChange"
          >
            <el-option
              v-for="item in templateList"
              :key="item.id"
              :label="`${item.name} (${item.metricKey})`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="比较操作符" prop="operator">
          <el-select v-model="ruleForm.operator" placeholder="请选择比较操作符" style="width: 100%">
            <el-option v-for="item in comparisonOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="阈值" prop="threshold">
          <el-input v-model.number="ruleForm.threshold" type="number" placeholder="请输入阈值" />
        </el-form-item>

        <el-form-item label="时长(秒)">
          <el-input v-model.number="ruleForm.duration" type="number" placeholder="默认30秒" disabled />
        </el-form-item>

        <el-form-item label="持续次数">
          <el-input v-model.number="ruleForm.continuousTimes" type="number" placeholder="默认1次" />
        </el-form-item>

        <el-form-item label="恢复判定时间(秒)">
          <el-input v-model.number="ruleForm.recoverDuration" type="number" placeholder="默认60秒" />
        </el-form-item>

        <el-form-item label="严重级别" prop="severity">
          <el-select v-model="ruleForm.severity" placeholder="请选择严重级别" style="width: 100%">
            <el-option v-for="item in severityOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="状态">
          <el-radio-group v-model="ruleForm.status">
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

    <!-- 预览对话框 -->
    <el-dialog
      v-model="previewDialogVisible"
      title="规则数据预览"
      width="800px"
      :close-on-click-modal="false"
    >
      <div v-loading="previewLoading" class="preview-content">
        <div v-if="previewData" class="preview-info">
          <!-- 规则信息 -->
          <div class="info-section">
            <div class="info-item">
              <span class="label">规则名称:</span>
              <span class="value">{{ previewData.ruleName }}</span>
            </div>
            <div class="info-item">
              <span class="label">PromQL:</span>
              <span class="value">{{ previewData.promql }}</span>
            </div>
            <div class="info-item">
              <span class="label">阈值:</span>
              <span class="value">{{ previewData.threshold }} {{ previewData.unit }}</span>
            </div>
            <div class="info-item">
              <span class="label">操作符:</span>
              <span class="value">{{ previewData.operator }}</span>
            </div>
            <div class="info-item">
              <span class="label">查询时间:</span>
              <span class="value">{{ new Date(previewData.queryTime * 1000).toLocaleString() }}</span>
            </div>
            <div class="info-item">
              <span class="label">结果数:</span>
              <span class="value">{{ previewData.resultCount }}</span>
            </div>
          </div>

          <!-- 查询结果表格 -->
          <div class="results-section">
            <h4>查询结果</h4>
            <el-table :data="previewData.results || []" stripe border style="width: 100%; font-size: 12px;">
              <el-table-column prop="hostname" label="主机名" width="100" />
              <el-table-column prop="instance" label="IP地址" width="130" />
              <el-table-column prop="value" label="当前值" width="100">
                <template #default="{ row }">
                  <span>{{ row.value }} {{ previewData.unit }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="isAnomalous" label="是否异常" width="90" align="center">
                <template #default="{ row }">
                  <el-tag :type="row.isAnomalous ? 'danger' : 'success'" size="small">
                    {{ row.isAnomalous ? '异常' : '正常' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="标签" min-width="200">
                <template #default="{ row }">
                  <div style="font-size: 11px; color: #909399;">
                    <div v-for="(value, key) in row.labels" :key="key" style="line-height: 1.5;">
                      <span>{{ key }}: {{ value }}</span>
                    </div>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="previewDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.alert-policy-rule-container {
  padding: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.breadcrumb-section {
  margin-bottom: 16px;
  padding: 12px 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 8px;
  border-left: 3px solid #667eea;
  display: flex;
  align-items: center;
  font-size: 14px;
}

.policy-info {
  color: #606266;
}

.policy-info strong {
  color: #667eea;
  font-weight: 600;
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

.preview-content {
  min-height: 200px;
}

.preview-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-section {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 16px;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.info-item {
  display: flex;
  gap: 8px;
  font-size: 13px;
}

.info-item .label {
  font-weight: 600;
  color: #2c3e50;
  min-width: 80px;
}

.info-item .value {
  color: #606266;
  word-break: break-all;
  flex: 1;
}

.results-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.results-section h4 {
  margin: 0;
  color: #2c3e50;
  font-size: 14px;
  font-weight: 600;
}
</style>
