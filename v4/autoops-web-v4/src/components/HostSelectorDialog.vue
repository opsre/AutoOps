<template>
  <el-dialog
    :title="title"
    v-model="visible"
    width="75%"
    :modal="true"
    class="host-selector-dialog"
    @close="handleClose"
  >
    <div class="dialog-content">
      <!-- 左侧分组树 -->
      <div class="left-panel">
        <div class="group-card">
          <!-- 标题区域 -->
          <div class="card-header">
            <div class="title-wrapper">
              <div class="title-icon">
                <div class="pulse-ring"></div>
                <el-icon class="main-icon"><FolderOpened /></el-icon>
              </div>
              <div class="title-content">
                <h3 class="title">业务分组</h3>
                <span class="subtitle">Business Groups</span>
              </div>
            </div>
            <div class="stats-indicator">
              <span class="stats-text">{{ groupList.length }}</span>
              <span class="stats-label">Groups</span>
            </div>
          </div>

          <!-- 搜索区域 -->
          <div class="search-container">
            <div class="search-input-wrapper">
              <el-input
                v-model="groupSearchText"
                placeholder="搜索分组..."
                clearable
                size="small"
                @input="handleGroupSearch"
                class="tech-input"
              >
                <template #prefix>
                  <el-icon class="search-icon"><Search /></el-icon>
                </template>
              </el-input>
            </div>
          </div>

          <!-- 所有主机按钮 -->
          <div class="all-hosts-btn-wrapper">
            <el-button
              type="primary"
              size="small"
              @click="showAllHosts"
              :class="{ 'active-group': !currentGroupId }"
              class="all-hosts-btn"
            >
              <el-icon><House /></el-icon>
              所有主机
              <span class="host-count-badge" v-if="allHostsCount > 0">({{ allHostsCount }})</span>
            </el-button>
          </div>

          <!-- 树形结构 -->
          <div class="tree-container">
            <el-tree
              ref="groupTree"
              :data="filteredGroupList"
              :props="{ children: 'children', label: 'name' }"
              node-key="id"
              :highlight-current="true"
              @node-click="handleGroupClick"
              class="tech-tree"
            >
              <template v-slot="{ node, data }">
                <div class="tree-node" :class="{ 'parent-node': !data.parentId, 'active': currentGroupId === data.id }">
                  <div class="node-content">
                    <div class="node-icon-wrapper">
                      <!-- 一级分组 -->
                      <template v-if="!data.parentId">
                        <div class="parent-icon">
                          <img
                            v-if="node.expanded"
                            :src="require('@/assets/image/打开文件夹.svg')"
                            class="parent-icon-img"
                            alt="打开文件夹"
                          />
                          <img
                            v-else
                            :src="require('@/assets/image/关闭文件夹.svg')"
                            class="parent-icon-img"
                            alt="关闭文件夹"
                          />
                        </div>
                      </template>
                      <!-- 二级分组 -->
                      <template v-else>
                        <div class="child-icon">
                          <div class="dot-indicator"></div>
                          <img
                            :src="require('@/assets/image/分组.svg')"
                            class="child-icon-img"
                            alt="子分组"
                          />
                        </div>
                      </template>
                    </div>
                    <div class="node-text">
                      <span class="node-label">
                        {{ node.label }}
                        <span class="host-count" v-if="data.hostCount !== undefined">({{ data.hostCount }})</span>
                      </span>
                    </div>
                  </div>
                  <div class="hover-effect"></div>
                </div>
              </template>
            </el-tree>
          </div>
        </div>
      </div>

      <!-- 右侧主机表格 -->
      <div class="right-panel">
        <div class="host-table-card">
          <!-- 搜索和统计 -->
          <div class="table-header">
            <div class="search-row">
              <el-input
                v-model="hostSearchKeyword"
                placeholder="搜索主机名称或IP地址"
                clearable
                size="small"
                style="width: 350px"
                @input="handleHostSearch"
                class="tech-input"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <div class="selection-info">
                <span class="selection-text">已选择</span>
                <span class="selection-count">{{ tempSelectedHosts.length }}</span>
                <span class="selection-text">台主机</span>
              </div>
            </div>
          </div>

          <!-- 主机表格 -->
          <el-table
            :data="displayHosts"
            border
            v-loading="loading"
            ref="hostTable"
            @selection-change="handleSelectionChange"
            class="modern-table"
            max-height="450"
            :header-cell-style="{
              background: '#ffffff',
              color: '#2c3e50',
              fontWeight: '600',
              borderBottom: '1px solid #e5e7eb'
            }"
          >
            <el-table-column type="selection" width="55" fixed="left" />
            <el-table-column prop="id" label="ID" width="70" />
            <el-table-column prop="hostName" label="主机名称" min-width="150">
              <template #default="{ row }">
                <div class="host-cell">
                  <el-icon color="#409EFF"><Monitor /></el-icon>
                  <span>{{ row.hostName || row.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="privateIp" label="内网IP" width="140">
              <template #default="{ row }">
                <div class="host-cell">
                  <el-icon color="#67C23A"><HomeFilled /></el-icon>
                  <span>{{ row.privateIp }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="sshIp" label="SSH地址" width="140">
              <template #default="{ row }">
                <div class="host-cell">
                  <el-icon color="#E6A23C"><Connection /></el-icon>
                  <span>{{ row.sshIp || row.publicIp || '-' }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="groupName" label="业务组" width="130" />
            <el-table-column prop="os" label="操作系统" width="140" show-overflow-tooltip />
            <el-table-column prop="remark" label="备注" min-width="120" show-overflow-tooltip />
          </el-table>

          <!-- 分页 -->
          <div class="pagination-wrapper">
            <el-pagination
              @current-change="handlePageChange"
              @size-change="handleSizeChange"
              :current-page="pagination.pageNum"
              :page-sizes="[10, 20, 50, 100]"
              :page-size="pagination.pageSize"
              layout="total, sizes, prev, pager, next, jumper"
              :total="pagination.total"
              class="modern-pagination"
            />
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <div class="footer-info">
          <span v-if="tempSelectedHosts.length > 0">
            已选择 <strong>{{ tempSelectedHosts.length }}</strong> 台主机
          </span>
          <span v-else class="empty-text">请选择主机</span>
        </div>
        <div class="footer-actions">
          <el-button @click="handleClose">取消</el-button>
          <el-button
            type="primary"
            @click="confirmSelection"
            :disabled="tempSelectedHosts.length === 0"
          >
            确认选择
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, onMounted, computed } from 'vue'
import {
  Search,
  Monitor,
  House,
  FolderOpened,
  HomeFilled,
  Connection
} from '@element-plus/icons-vue'
import cmdbAPI from '@/api/cmdb'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '选择主机'
  }
})

const emit = defineEmits(['update:modelValue', 'hosts-selected'])

const visible = ref(false)
const loading = ref(false)
const groupList = ref([])
const filteredGroupList = ref([])
const groupSearchText = ref('')
const currentGroupId = ref(null)
const allHostsCount = ref(0)

const hostList = ref([])
const displayHosts = ref([])
const hostSearchKeyword = ref('')
const tempSelectedHosts = ref([])

const pagination = ref({
  pageNum: 1,
  pageSize: 10,
  total: 0
})

const groupTree = ref(null)

const hostTable = ref(null)

// 加载分组列表
const fetchGroupList = async () => {
  try {
    const response = await cmdbAPI.getAllCmdbGroups()
    if (response.data?.code === 200) {
      groupList.value = response.data.data || []
      filteredGroupList.value = groupList.value
    }
  } catch (error) {
    console.error('获取分组列表失败:', error)
  }
}

// 分组搜索
const handleGroupSearch = () => {
  if (!groupSearchText.value.trim()) {
    filteredGroupList.value = groupList.value
    return
  }

  const keyword = groupSearchText.value.toLowerCase()
  const filterGroups = (groups) => {
    return groups.filter(group => {
      const match = group.name.toLowerCase().includes(keyword)
      if (group.children && group.children.length > 0) {
        group.children = filterGroups(group.children)
        return match || group.children.length > 0
      }
      return match
    })
  }

  filteredGroupList.value = filterGroups(JSON.parse(JSON.stringify(groupList.value)))
}

// 显示所有主机
const showAllHosts = async () => {
  currentGroupId.value = null
  hostSearchKeyword.value = ''
  pagination.value.pageNum = 1
  await fetchAllHosts()
}

// 获取所有主机
const fetchAllHosts = async () => {
  loading.value = true
  try {
    const response = await cmdbAPI.getCmdbHostList({
      page: pagination.value.pageNum,
      pageSize: pagination.value.pageSize
    })

    if (response.data?.code === 200) {
      const data = response.data.data
      hostList.value = data.list || []
      displayHosts.value = hostList.value
      pagination.value.total = data.total || 0
      allHostsCount.value = data.total || 0
    }
  } catch (error) {
    console.error('获取主机列表失败:', error)
    hostList.value = []
    displayHosts.value = []
  } finally {
    loading.value = false
  }
}

// 点击分组节点
const handleGroupClick = async (data) => {
  currentGroupId.value = data.id
  hostSearchKeyword.value = ''
  pagination.value.pageNum = 1
  await fetchHostsByGroup(data.id)
}

// 根据分组获取主机
const fetchHostsByGroup = async (groupId) => {
  loading.value = true
  try {
    const response = await cmdbAPI.getCmdbHostsByGroupId(groupId, {
      page: pagination.value.pageNum,
      pageSize: pagination.value.pageSize
    })

    if (response.data?.code === 200) {
      const data = response.data.data
      // 后端返回的是 {total, list, page, pageSize} 结构
      if (data && data.list) {
        hostList.value = data.list || []
        displayHosts.value = data.list || []
        pagination.value.total = data.total || 0
      } else if (Array.isArray(data)) {
        // 兼容直接返回数组的情况
        hostList.value = data
        displayHosts.value = data
        pagination.value.total = data.length
      } else {
        hostList.value = []
        displayHosts.value = []
        pagination.value.total = 0
      }
    }
  } catch (error) {
    console.error('获取分组主机失败:', error)
    hostList.value = []
    displayHosts.value = []
  } finally {
    loading.value = false
  }
}

// 主机搜索
const handleHostSearch = async () => {
  pagination.value.pageNum = 1

  if (!hostSearchKeyword.value.trim()) {
    // 如果搜索框为空，恢复当前视图
    if (currentGroupId.value) {
      await fetchHostsByGroup(currentGroupId.value)
    } else {
      await fetchAllHosts()
    }
    return
  }

  // 执行搜索
  loading.value = true
  try {
    const response = await cmdbAPI.GetCmdbHostsByHostNameLike(hostSearchKeyword.value.trim(), {
      page: pagination.value.pageNum,
      pageSize: pagination.value.pageSize
    })

    if (response.data?.code === 200) {
      const data = response.data.data
      if (data && data.list) {
        displayHosts.value = data.list || []
        pagination.value.total = data.total || 0
      } else if (Array.isArray(data)) {
        displayHosts.value = data
        pagination.value.total = data.length
      } else {
        displayHosts.value = []
        pagination.value.total = 0
      }
    }
  } catch (error) {
    console.error('搜索主机失败:', error)
    displayHosts.value = []
  } finally {
    loading.value = false
  }
}

// 分页变化
const handlePageChange = (page) => {
  pagination.value.pageNum = page
  if (hostSearchKeyword.value.trim()) {
    handleHostSearch()
  } else if (currentGroupId.value) {
    fetchHostsByGroup(currentGroupId.value)
  } else {
    fetchAllHosts()
  }
}

const handleSizeChange = (size) => {
  pagination.value.pageSize = size
  pagination.value.pageNum = 1
  if (hostSearchKeyword.value.trim()) {
    handleHostSearch()
  } else if (currentGroupId.value) {
    fetchHostsByGroup(currentGroupId.value)
  } else {
    fetchAllHosts()
  }
}

// 选择变化
const handleSelectionChange = (selection) => {
  tempSelectedHosts.value = selection
}

// 确认选择
const confirmSelection = () => {
  emit('hosts-selected', tempSelectedHosts.value)
  handleClose()
}

// 关闭对话框
const handleClose = () => {
  visible.value = false
  tempSelectedHosts.value = []
  hostSearchKeyword.value = ''
  groupSearchText.value = ''
}

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    // 打开对话框时加载数据
    fetchGroupList()
    fetchAllHosts()
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})
</script>

<style scoped>
/* 表格表头颜色修复 - 黑色字体 */
:deep(.modern-table .el-table__header-wrapper .el-table__header thead th) {
  color: #2c3e50 !important;
  background-color: transparent !important;
}

:deep(.modern-table .el-table__header-wrapper .el-table__header thead th .cell) {
  color: #2c3e50 !important;
}

.host-selector-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: #f8f9fa;
}

.host-selector-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  padding: 20px 24px;
  border-radius: 16px 16px 0 0;
}

.host-selector-dialog :deep(.el-dialog__title) {
  color: #ffffff;
  font-weight: 600;
  font-size: 18px;
}

.host-selector-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
  color: #ffffff;
  font-size: 20px;
}

.host-selector-dialog :deep(.el-dialog__body) {
  padding: 0;
  height: 550px;
}

.dialog-content {
  display: flex;
  height: 100%;
  gap: 0;
}

/* 左侧面板 */
.left-panel {
  width: 280px;
  background: #f8f9fa;
  border-right: 1px solid #e5e7eb;
  padding: 15px;
}

.group-card {
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.95), rgba(255, 255, 255, 0.85));
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* 标题区域 */
.card-header {
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

.title-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
}

.title-icon {
  position: relative;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pulse-ring {
  position: absolute;
  width: 32px;
  height: 32px;
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
  font-size: 18px;
  color: #667eea;
  z-index: 2;
  background: rgba(102, 126, 234, 0.1);
  padding: 6px;
  border-radius: 6px;
}

.title-content .title {
  font-size: 16px;
  font-weight: 700;
  margin: 0 0 2px 0;
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.subtitle {
  font-size: 10px;
  color: rgba(102, 126, 234, 0.6);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.stats-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: rgba(102, 126, 234, 0.05);
  padding: 6px 10px;
  border-radius: 6px;
  border: 1px solid rgba(102, 126, 234, 0.1);
}

.stats-text {
  font-size: 16px;
  font-weight: 700;
  color: #667eea;
  line-height: 1;
}

.stats-label {
  font-size: 9px;
  color: rgba(102, 126, 234, 0.6);
  text-transform: uppercase;
}

/* 搜索区域 */
.search-container {
  padding: 12px 16px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.08);
}

.search-input-wrapper {
  position: relative;
}

.tech-input :deep(.el-input__wrapper) {
  background: rgba(102, 126, 234, 0.03);
  border: 1px solid rgba(102, 126, 234, 0.15);
  border-radius: 10px;
  transition: all 0.3s ease;
}

.tech-input :deep(.el-input__wrapper):hover {
  border-color: rgba(102, 126, 234, 0.3);
  background: rgba(102, 126, 234, 0.05);
}

.tech-input :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.search-icon {
  color: rgba(102, 126, 234, 0.6);
}

/* 所有主机按钮 */
.all-hosts-btn-wrapper {
  padding: 8px 16px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.08);
}

.all-hosts-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.all-hosts-btn.active-group {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-color: transparent;
  color: white;
}

.host-count-badge {
  font-size: 12px;
  margin-left: 4px;
}

/* 树形结构 */
.tree-container {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.tech-tree {
  background: transparent;
  border: none;
}

.tech-tree :deep(.el-tree-node__content) {
  height: auto;
  padding: 0;
  background: transparent;
}

.tree-node {
  position: relative;
  margin: 2px 12px;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.tree-node.active {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
}

.tree-node:hover > .hover-effect {
  opacity: 1;
}

.hover-effect {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05), rgba(118, 75, 162, 0.05));
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 1;
}

.node-content {
  display: flex;
  align-items: center;
  padding: 6px 10px;
  position: relative;
  z-index: 2;
}

.node-icon-wrapper {
  margin-right: 8px;
  display: flex;
  align-items: center;
}

.parent-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.parent-icon-img {
  width: 18px;
  height: 18px;
}

.child-icon {
  position: relative;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dot-indicator {
  position: absolute;
  width: 5px;
  height: 5px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 50%;
  left: -6px;
}

.child-icon-img {
  width: 16px;
  height: 16px;
}

.node-text {
  flex: 1;
}

.node-label {
  font-size: 13px;
  color: #2c3e50;
}

.parent-node .node-label {
  font-weight: 500;
  font-size: 14px;
}

.host-count {
  font-size: 11px;
  color: rgba(102, 126, 234, 0.7);
  margin-left: 4px;
}

/* 右侧面板 */
.right-panel {
  flex: 1;
  background: white;
  padding: 15px;
}

.host-table-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.table-header {
  margin-bottom: 15px;
}

.search-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.selection-info {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.1);
}

.selection-text {
  font-size: 13px;
  color: #666;
}

.selection-count {
  font-size: 18px;
  font-weight: 700;
  color: #667eea;
}

/* 表格样式 */
.modern-table {
  border-radius: 12px;
  overflow: hidden;
  flex: 1;
}

/* 表头背景和文字颜色 - 白色背景，黑色字体 */
.modern-table :deep(.el-table__header-wrapper) {
  background: #ffffff !important;
}

.modern-table :deep(.el-table__header) {
  background: #ffffff !important;
}

.modern-table :deep(.el-table__header th) {
  background: #ffffff !important;
  color: #2c3e50 !important;
  font-weight: 600 !important;
  border-bottom: 1px solid #e5e7eb !important;
}

.modern-table :deep(.el-table__header th .cell) {
  color: #2c3e50 !important;
  font-weight: 600 !important;
}

.modern-table :deep(.el-table th.el-table__cell) {
  background-color: #ffffff !important;
  color: #2c3e50 !important;
}

.modern-table :deep(.el-table__body tr:hover > td) {
  background: rgba(102, 126, 234, 0.05);
}

.host-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 分页 */
.pagination-wrapper {
  margin-top: 15px;
  display: flex;
  justify-content: center;
}

.modern-pagination :deep(.el-pager li.is-active) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
}

/* 底部 */
.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-top: 1px solid #e5e7eb;
  background: white;
}

.footer-info {
  font-size: 14px;
  color: #606266;
}

.footer-info strong {
  color: #667eea;
  font-size: 16px;
}

.empty-text {
  color: #999;
}

.footer-actions {
  display: flex;
  gap: 10px;
}

/* 滚动条 */
.tree-container::-webkit-scrollbar {
  width: 4px;
}

.tree-container::-webkit-scrollbar-track {
  background: rgba(102, 126, 234, 0.05);
}

.tree-container::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.3), rgba(118, 75, 162, 0.3));
  border-radius: 2px;
}
</style>
