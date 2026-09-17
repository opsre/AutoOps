<script setup>
import { ref, onMounted, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import { businessGroupApi, apiBusinessFunctionListApi, apiMonitorApi } from '@/api/monitor'
import healthIcon from '@/assets/image/健康.svg'
import fireIcon from '@/assets/image/起火.svg'

// ==================== 状态管理 ====================

// 第一级：业务组列表
const groupList = ref([])
const selectedGroupId = ref(null)
const selectedGroupName = ref('')

// 第二级：业务功能列表
const functionList = ref([])
const selectedFunctionId = ref(null)
const selectedFunctionName = ref('')
const functionsLoading = ref(false)

// 第三级：API列表
const apiList = ref([])
const apisLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  keyword: ''
})

// ==================== API 调用 ====================

// 获取业务组列表
const fetchGroupList = async () => {
  try {
    const { data: res } = await businessGroupApi.getGroupList()
    if (res.code === 200) {
      // 提取二级业务组
      const secondLevelGroups = []
      res.data.forEach(parentGroup => {
        if (parentGroup.children && Array.isArray(parentGroup.children) && parentGroup.children.length > 0) {
          parentGroup.children.forEach(child => {
            secondLevelGroups.push({
              id: child.id,
              name: child.name,
              parentId: child.parentId,
              parentName: parentGroup.name
            })
          })
        }
      })
      groupList.value = secondLevelGroups
    }
  } catch (error) {
    ElMessage.error('获取业务组列表失败')
    console.error(error)
  }
}

// 获取业务功能列表
const fetchFunctionList = async (groupId) => {
  if (!groupId) {
    functionList.value = []
    apiList.value = []
    return
  }

  functionsLoading.value = true
  try {
    const { data: res } = await apiBusinessFunctionListApi.getByGroup(groupId)
    if (res.code === 200) {
      functionList.value = Array.isArray(res.data) ? res.data : []
      // 重置API列表
      apiList.value = []
      selectedFunctionId.value = null
    }
  } catch (error) {
    ElMessage.error('获取业务功能列表失败')
    functionList.value = []
    console.error(error)
  } finally {
    functionsLoading.value = false
  }
}

// 获取API列表
const fetchApiList = async (functionId) => {
  if (!functionId) {
    apiList.value = []
    return
  }

  apisLoading.value = true
  try {
    const { data: res } = await apiMonitorApi.getByBusinessFunction(functionId)
    if (res.code === 200) {
      let apiData = res.data
      // 支持分页响应
      if (apiData.list) {
        apiData = apiData.list
      }
      apiList.value = Array.isArray(apiData) ? apiData : []
    }
  } catch (error) {
    ElMessage.error('获取API列表失败')
    apiList.value = []
    console.error(error)
  } finally {
    apisLoading.value = false
  }
}

// ==================== 事件处理 ====================

// 选择业务组
const handleSelectGroup = (group) => {
  selectedGroupId.value = group.id
  selectedGroupName.value = `${group.parentName} / ${group.name}`
  selectedFunctionId.value = null
  selectedFunctionName.value = ''
  fetchFunctionList(group.id)
}

// 选择业务功能
const handleSelectFunction = (func) => {
  selectedFunctionId.value = func.id
  selectedFunctionName.value = func.name
  fetchApiList(func.id)
}

// 搜索API
const handleSearch = () => {
  // 可以在这里实现搜索逻辑
  // 当前版本简化处理，前端过滤
}

// 重置搜索
const handleReset = () => {
  searchForm.keyword = ''
}

// 获取健康状态显示
const getHealthStatus = (row) => {
  if (row.healthStatus === 'healthy') {
    return '健康'
  }
  return '异常'
}

// 过滤后的API列表
const filteredApiList = () => {
  if (!searchForm.keyword) {
    return apiList.value
  }
  const keyword = searchForm.keyword.toLowerCase()
  return apiList.value.filter(api =>
    api.name.toLowerCase().includes(keyword) ||
    api.url.toLowerCase().includes(keyword)
  )
}

// ==================== 生命周期 ====================

onMounted(() => {
  fetchGroupList()
})
</script>

<template>
  <div class="business-function-catalog">
    <!-- 主内容区域：左右布局 -->
    <div class="main-content">
      <!-- 左侧分类栏 -->
      <div class="category-sidebar">
        <!-- 业务组选择 -->
        <div class="category-section">
          <div class="section-header">业务组</div>
          <div class="category-list">
            <div
              v-for="group in groupList"
              :key="group.id"
              :class="['category-item', { active: selectedGroupId === group.id }]"
              @click="handleSelectGroup(group)"
              :title="`${group.parentName} / ${group.name}`"
            >
              <span class="category-label">{{ group.name }}</span>
            </div>
          </div>
        </div>

        <!-- 业务功能选择 -->
        <div class="category-section" v-if="selectedGroupId">
          <div class="section-header">业务功能</div>
          <div class="category-list">
            <div
              v-if="functionsLoading"
              style="padding: 12px 16px; text-align: center; color: #909399; font-size: 12px"
            >
              加载中...
            </div>
            <div
              v-else-if="functionList.length === 0"
              style="padding: 12px 16px; text-align: center; color: #909399; font-size: 12px"
            >
              暂无功能
            </div>
            <div
              v-for="func in functionList"
              v-else
              :key="func.id"
              :class="['category-item', { active: selectedFunctionId === func.id }]"
              @click="handleSelectFunction(func)"
              :title="func.name"
            >
              <span class="category-label">{{ func.name }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧内容区域 -->
      <div class="content-area">
        <!-- 面包屑导航 -->
        <div class="breadcrumb-section" v-if="selectedGroupId">
          <div class="breadcrumb">
            <span class="breadcrumb-item">{{ selectedGroupName }}</span>
            <span class="breadcrumb-separator">/</span>
            <span class="breadcrumb-item" v-if="selectedFunctionId">{{ selectedFunctionName }}</span>
            <span class="breadcrumb-item" v-else style="color: #bfbfbf">选择业务功能</span>
          </div>
        </div>

        <!-- 搜索区域 -->
        <div class="search-section" v-if="selectedFunctionId">
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

        <!-- 统计卡片 -->
        <div class="stats-section" v-if="selectedFunctionId">
          <div class="stats-card">
            <div class="stats-label">总 API 数</div>
            <div class="stats-value">{{ filteredApiList().length }}</div>
          </div>
          <div class="stats-card healthy">
            <div class="stats-label">健康</div>
            <div class="stats-value">
              {{ filteredApiList().filter(a => a.healthStatus === 'healthy').length }}
            </div>
          </div>
          <div class="stats-card unhealthy">
            <div class="stats-label">异常</div>
            <div class="stats-value">
              {{ filteredApiList().filter(a => a.healthStatus === 'unhealthy').length }}
            </div>
          </div>
        </div>

        <!-- 空状态提示 -->
        <div class="empty-state" v-if="!selectedFunctionId && selectedGroupId">
          <div class="empty-icon">📋</div>
          <div class="empty-text">请在左侧选择业务功能查看API详情</div>
        </div>

        <div class="empty-state" v-if="!selectedGroupId">
          <div class="empty-icon">🔍</div>
          <div class="empty-text">请在左侧选择业务组</div>
        </div>

        <!-- 表格 -->
        <div class="table-section" v-if="selectedFunctionId">
          <el-table
            :data="filteredApiList()"
            v-loading="apisLoading"
            stripe
            border
            style="width: 100%"
            :header-cell-style="{ background: '#eef1f6', color: '#606266' }"
          >
            <el-table-column prop="id" label="ID" width="70" align="center" />

            <el-table-column prop="name" label="监控名称" min-width="120">
              <template #default="{ row }">
                <span style="font-weight: 500">{{ row.name }}</span>
              </template>
            </el-table-column>

            <el-table-column prop="url" label="接口地址" min-width="200">
              <template #default="{ row }">
                <div style="display: flex; align-items: center; gap: 6px">
                  <el-tag size="small">{{ row.method || 'GET' }}</el-tag>
                  <span
                    style="color: #606266; font-size: 12px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; flex: 1"
                    :title="row.url"
                  >
                    {{ row.url }}
                  </span>
                </div>
              </template>
            </el-table-column>

            <el-table-column label="健康状态" width="110" align="center">
              <template #default="{ row }">
                <div style="display: flex; align-items: center; justify-content: center; gap: 6px">
                  <img
                    :src="row.healthStatus === 'healthy' ? healthIcon : fireIcon"
                    :alt="row.healthStatus"
                    style="width: 18px; height: 18px"
                  />
                  <span style="font-size: 12px; font-weight: 500">
                    {{ getHealthStatus(row) }}
                  </span>
                </div>
              </template>
            </el-table-column>

            <el-table-column prop="lastResponseTime" label="最近响应(ms)" width="120" align="center">
              <template #default="{ row }">
                <span>{{ row.lastResponseTime || '-' }}</span>
              </template>
            </el-table-column>

            <el-table-column prop="lastStatusCode" label="状态码" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="row.lastStatusCode === 200 ? 'success' : 'danger'" size="small">
                  {{ row.lastStatusCode || '-' }}
                </el-tag>
              </template>
            </el-table-column>

            <el-table-column label="状态" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>

            <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
          </el-table>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.business-function-catalog {
  padding: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
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
  width: 200px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  overflow: hidden;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  max-height: 100%;
}

.category-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  border-bottom: 1px solid #f0f0f0;
}

.category-section:last-child {
  border-bottom: none;
}

.section-header {
  padding: 12px 16px;
  background: #f5f7fa;
  color: #303133;
  font-size: 14px;
  font-weight: 600;
  border-bottom: 1px solid #e4e7ed;
  flex-shrink: 0;
}

.category-list {
  padding: 8px 0;
  flex: 1;
  overflow-y: auto;
}

.category-item {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 13px;
  color: #606266;
  border-left: 3px solid transparent;
}

.category-item:hover {
  background: #f5f7fa;
}

.category-item.active {
  background: #ecf5ff;
  color: #409eff;
  font-weight: 500;
  border-left-color: #409eff;
}

.category-label {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 右侧内容区域 */
.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

/* 面包屑导航 */
.breadcrumb-section {
  padding: 12px 16px;
  background: #fafbfc;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.breadcrumb {
  display: flex;
  align-items: center;
  font-size: 13px;
  color: #606266;
}

.breadcrumb-item {
  color: #409eff;
  font-weight: 500;
}

.breadcrumb-item:first-child {
  color: #303133;
  font-weight: 600;
}

.breadcrumb-separator {
  margin: 0 8px;
  color: #bfbfbf;
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

/* 统计卡片 */
.stats-section {
  display: flex;
  gap: 16px;
}

.stats-card {
  flex: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  padding: 16px;
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.2);
  transition: all 0.3s ease;
}

.stats-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.stats-card.healthy {
  background: linear-gradient(135deg, #67C23A 0%, #85CE61 100%);
  box-shadow: 0 4px 15px rgba(103, 194, 58, 0.2);
}

.stats-card.healthy:hover {
  box-shadow: 0 8px 25px rgba(103, 194, 58, 0.3);
}

.stats-card.unhealthy {
  background: linear-gradient(135deg, #F56C6C 0%, #F78989 100%);
  box-shadow: 0 4px 15px rgba(245, 108, 108, 0.2);
}

.stats-card.unhealthy:hover {
  box-shadow: 0 8px 25px rgba(245, 108, 108, 0.3);
}

.stats-label {
  font-size: 12px;
  opacity: 0.9;
  margin-bottom: 6px;
}

.stats-value {
  font-size: 24px;
  font-weight: bold;
}

/* 空状态 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #909399;
}

.empty-icon {
  font-size: 48px;
  opacity: 0.5;
}

.empty-text {
  font-size: 14px;
}

/* 表格区域 */
.table-section {
  flex: 1;
  min-height: 0;
  overflow: auto;
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

:deep(.reset-btn) {
  background: linear-gradient(45deg, #E6A23C, #EEBE77);
  color: white;
}

/* 滚动条美化 */
.category-list::-webkit-scrollbar {
  width: 6px;
}

.category-list::-webkit-scrollbar-track {
  background: transparent;
}

.category-list::-webkit-scrollbar-thumb {
  background: #d9d9d9;
  border-radius: 3px;
}

.category-list::-webkit-scrollbar-thumb:hover {
  background: #bfbfbf;
}
</style>
