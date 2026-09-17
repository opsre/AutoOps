<template>
  <div class="ssh-container">
    <div class="ssh-layout">
      <!-- 左侧主机树区域 -->
      <div class="tree-section" :style="{ width: treeWidth + 'px' }">
        <el-card shadow="never" class="group-card">
          <template #header>
            <div class="card-header">
              <span class="header-title">主机列表</span>
            </div>
          </template>

          <!-- 搜索框和刷新按钮 -->
          <div class="search-container">
            <el-input
              v-model="groupSearchText"
              placeholder="输入主机名/IP..."
              clearable
              size="small"
              class="search-input"
              @input="handleGroupSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button
              size="small"
              :icon="Refresh"
              class="refresh-btn"
              @click="handleRefresh"
              title="刷新主机列表"
            />
          </div>

          <!-- 主机树 -->
          <el-tree
            ref="groupTree"
            :data="groupList"
            :props="defaultProps"
            node-key="id"
            :expanded-keys="expandedKeys"
            :highlight-current="true"
            @node-click="handleGroupClick"
            @node-expand="handleNodeExpand"
            @node-collapse="handleNodeCollapse"
            class="group-tree"
          >
            <template #default="{ node, data }">
              <span :class="{ 'business-group': !data.parentId, 'host-node': data.isHost }">
                <template v-if="data.isHost">
                  <!-- 主机在线图标 -->
                  <span class="tree-icon host-icon">
                    <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg">
                      <path d="M839.68 759.808H184.32c-22.528 0-40.96-18.432-40.96-40.96V184.32c0-22.528 18.432-40.96 40.96-40.96h655.36c22.528 0 40.96 18.432 40.96 40.96v534.528c0 22.528-18.432 40.96-40.96 40.96z" fill="#004FFF"></path>
                      <path d="M622.592 852.992H401.408l17.408-93.184h186.368z" fill="#FBB03B"></path>
                      <path d="M622.592 852.992H401.408l-36.864 16.384h294.912zM364.544 869.376h295.936V880.64H364.544z" fill="#004FFF"></path>
                      <path d="M171.008 172.032h683.008v424.96H171.008z" fill="#E9F3FB"></path>
                      <path d="M491.52 670.72a20.48 20.48 0 1 0 40.96 0 20.48 20.48 0 1 0-40.96 0Z" fill="#E9F3FB"></path>
                      <path d="M197.632 196.608h628.736v375.808H197.632z" fill="#FBB03B"></path>
                      <path d="M300.032 362.496V196.608h-102.4v375.808h628.736v-81.92H428.032c-70.656 0-128-57.344-128-128z" fill="#F79839"></path>
                    </svg>
                  </span>
                </template>
                <template v-else-if="data.parentId">
                  <!-- 子分组：根据展开状态显示打开/关闭文件夹 -->
                  <span class="tree-icon folder-icon">
                    <img v-if="node.expanded" :src="FolderOpenIcon" alt="打开文件夹" />
                    <img v-else :src="FolderCloseIcon" alt="关闭文件夹" />
                  </span>
                </template>
                <template v-else>
                  <!-- 业务组：根据展开状态显示打开/关闭文件夹 -->
                  <span class="tree-icon folder-icon">
                    <img v-if="node.expanded" :src="FolderOpenIcon" alt="打开文件夹" />
                    <img v-else :src="FolderCloseIcon" alt="关闭文件夹" />
                  </span>
                </template>
                <span class="node-label">{{ node.label }}</span>
                <span v-if="!data.isHost && data.hostCount !== undefined" class="host-count">({{ data.hostCount }})</span>
              </span>
            </template>
          </el-tree>
        </el-card>
      </div>

      <!-- 拖拽分隔条 -->
      <div
        class="resizer"
        @mousedown="startResize"
        @dblclick="resetTreeWidth"
        title="拖动调整宽度，双击恢复默认"
      >
        <div class="resizer-line"></div>
      </div>

      <!-- 右侧终端标签页区域 -->
      <div class="terminal-section">
        <!-- 空状态 - 当没有打开任何终端时显示 -->
        <div v-if="terminalTabs.length === 0" class="welcome-placeholder">
          <div class="welcome-icon">
            <svg viewBox="0 0 200 150" xmlns="http://www.w3.org/2000/svg">
              <!-- 显示器 -->
              <rect x="20" y="20" width="160" height="90" fill="none" stroke="#409eff" stroke-width="3" rx="5"/>
              <!-- 屏幕内容 -->
              <rect x="30" y="30" width="140" height="70" fill="#2c3e50"/>
              <!-- 终端提示符 -->
              <text x="40" y="55" font-family="monospace" font-size="12" fill="#67c23a">$</text>
              <line x1="50" y1="48" x2="70" y2="48" stroke="#67c23a" stroke-width="2"/>
              <!-- 底座 -->
              <line x1="80" y1="110" x2="120" y2="110" stroke="#409eff" stroke-width="3"/>
              <line x1="100" y1="110" x2="100" y2="125" stroke="#409eff" stroke-width="3"/>
              <rect x="70" y="125" width="60" height="5" fill="#409eff" rx="2"/>
            </svg>
          </div>
          <div class="welcome-text">
            <h1>AutoOps</h1>
            <p>双击左侧主机节点开启终端会话</p>
          </div>
        </div>

        <!-- 终端标签页 -->
        <el-tabs
          v-else
          v-model="activeTabName"
          type="card"
          closable
          @tab-remove="handleTabRemove"
          @tab-click="handleTabClick"
          class="terminal-tabs"
        >
          <el-tab-pane
            v-for="tab in terminalTabs"
            :key="tab.name"
            :label="tab.title"
            :name="tab.name"
          >
            <!-- 每个标签页对应一个终端实例 -->
            <terminal
              v-if="tab.visible"
              :ref="'terminal_' + tab.name"
              :current-host="tab.host"
              :tab-name="tab.name"
              @terminal-closed="handleTerminalClosed"
            />
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import Terminal from './Terminal.vue'
import cmdbApi from '@/api/cmdb'
import FolderOpenIcon from '@/assets/image/打开文件夹.svg'
import FolderCloseIcon from '@/assets/image/关闭文件夹.svg'

// 数据状态
const groupSearchText = ref('')
const expandedKeys = ref([])
const groupList = ref([])
const groupTree = ref(null)

const defaultProps = {
  children: 'children',
  label: 'name'
}

// 终端标签页管理
const terminalTabs = ref([])
const activeTabName = ref('')
let tabIndex = 0

// 拖拽调整宽度相关
const treeWidth = ref(280) // 默认宽度280px
const minTreeWidth = 200 // 最小宽度
const maxTreeWidth = 600 // 最大宽度
const isResizing = ref(false)

// 加载主机分组列表
const loadGroupList = async () => {
  try {
    const response = await cmdbApi.getGroupListWithHosts()
    if (response.data.code === 200) {
      groupList.value = response.data.data.map(group => ({
        ...group,
        children: group.children ? group.children.map(child => ({
          ...child,
          children: child.hosts ? child.hosts.map(host => ({
            id: `host_${host.id}`,
            name: host.hostName,
            isHost: true,
            hostData: host
          })) : []
        })) : []
      }))

      // 默认展开第一个分组
      if (groupList.value.length > 0) {
        expandedKeys.value = [groupList.value[0].id]
      }
    }
  } catch (error) {
    console.error('加载分组列表失败:', error)
    ElMessage.error('加载分组列表失败')
  }
}

// 处理主机节点点击（双击打开终端）
const handleGroupClick = (data, node) => {
  if (data.isHost && node.parent.parent) {
    // 双击检测逻辑
    const now = Date.now()
    const lastClickTime = node._lastClickTime || 0

    if (now - lastClickTime < 300) {
      // 双击 - 打开终端
      openTerminal(data.hostData)
    }

    node._lastClickTime = now
  }
}

// 打开新的终端标签页
const openTerminal = (host) => {
  // 检查是否已经打开该主机的终端
  const existingTab = terminalTabs.value.find(tab => tab.host.id === host.id)

  if (existingTab) {
    // 如果已存在，切换到该标签页
    activeTabName.value = existingTab.name
    ElMessage.info(`已切换到 ${host.hostName} 的终端`)
    return
  }

  // 创建新的标签页
  const newTabName = `tab_${++tabIndex}`
  const newTab = {
    name: newTabName,
    title: host.hostName,
    host: {
      id: host.id,
      hostName: host.hostName,
      sshName: host.sshName,
      sshIp: host.sshIp,
      sshPort: host.sshPort
    },
    visible: true
  }

  terminalTabs.value.push(newTab)
  activeTabName.value = newTabName

  // 等待 DOM 更新后自动连接
  nextTick(() => {
    ElMessage.success(`已打开 ${host.hostName} 的终端`)
  })
}

// 关闭标签页
const handleTabRemove = (targetName) => {
  const tabs = terminalTabs.value
  let nextActiveTab = activeTabName.value

  if (nextActiveTab === targetName) {
    tabs.forEach((tab, index) => {
      if (tab.name === targetName) {
        const nextTab = tabs[index + 1] || tabs[index - 1]
        if (nextTab) {
          nextActiveTab = nextTab.name
        }
      }
    })
  }

  activeTabName.value = nextActiveTab
  terminalTabs.value = tabs.filter(tab => tab.name !== targetName)
}

// 标签页点击切换
const handleTabClick = (tab) => {
  // 可以在这里添加标签切换的额外逻辑
}

// 终端关闭回调
const handleTerminalClosed = (tabName) => {
  handleTabRemove(tabName)
}

// 刷新主机列表
const handleRefresh = async () => {
  groupSearchText.value = ''
  expandedKeys.value = []
  await loadGroupList()
  ElMessage.success('主机列表已刷新')
}

// 搜索分组
const handleGroupSearch = () => {
  if (!groupSearchText.value) {
    expandedKeys.value = []
    return
  }

  const findMatchingNode = (nodes, searchText) => {
    const results = []
    const traverse = (nodeList, parentKeys = []) => {
      nodeList.forEach(node => {
        const currentKeys = [...parentKeys, node.id]
        const searchLower = searchText.toLowerCase()
        const nameMatch = node.name.toLowerCase().includes(searchLower)
        const ipMatch = node.isHost && node.hostData && node.hostData.sshIp && node.hostData.sshIp.toLowerCase().includes(searchLower)

        if (nameMatch || ipMatch) {
          results.push({ node, keys: currentKeys })
        }
        if (node.children && node.children.length > 0) {
          traverse(node.children, currentKeys)
        }
      })
    }
    traverse(nodes)
    return results
  }

  const matches = findMatchingNode(groupList.value, groupSearchText.value)
  if (matches.length > 0) {
    // 展开所有匹配节点的父路径
    const allKeys = new Set()
    matches.forEach(match => {
      match.keys.slice(0, -1).forEach(key => allKeys.add(key))
    })
    expandedKeys.value = Array.from(allKeys)

    // 高亮第一个匹配
    nextTick(() => {
      if (groupTree.value) {
        groupTree.value.setCurrentKey(matches[0].node.id)
      }
    })
  } else {
    ElMessage.warning('未找到匹配的主机')
  }
}

// 节点展开/折叠
const handleNodeExpand = (data, node) => {
  if (!expandedKeys.value.includes(node.key)) {
    expandedKeys.value.push(node.key)
  }
}

const handleNodeCollapse = (data, node) => {
  expandedKeys.value = expandedKeys.value.filter(key => key !== node.key)
}

// 开始拖拽调整宽度
const startResize = (e) => {
  e.preventDefault()
  isResizing.value = true
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'

  // 获取容器的左边界位置
  const container = e.target.closest('.ssh-layout')
  if (!container) return

  const containerRect = container.getBoundingClientRect()

  const handleMouseMove = (moveEvent) => {
    if (!isResizing.value) return

    // 计算鼠标相对于容器左边界的距离，这就是新的宽度
    let newWidth = moveEvent.clientX - containerRect.left

    // 限制在最小和最大宽度之间
    if (newWidth < minTreeWidth) {
      newWidth = minTreeWidth
    } else if (newWidth > maxTreeWidth) {
      newWidth = maxTreeWidth
    }

    treeWidth.value = newWidth
  }

  const handleMouseUp = () => {
    isResizing.value = false
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)

    // 保存宽度到本地存储
    localStorage.setItem('ssh-tree-width', treeWidth.value.toString())
  }

  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
}

// 重置树宽度
const resetTreeWidth = () => {
  treeWidth.value = 280
  localStorage.setItem('ssh-tree-width', '280')
  ElMessage.success('已恢复默认宽度')
}

// 从本地存储加载宽度
const loadTreeWidth = () => {
  const savedWidth = localStorage.getItem('ssh-tree-width')
  if (savedWidth) {
    const width = parseInt(savedWidth, 10)
    if (width >= minTreeWidth && width <= maxTreeWidth) {
      treeWidth.value = width
    }
  }
}

// 组件挂载
onMounted(() => {
  loadGroupList()
  loadTreeWidth()
})

// 组件卸载时清理
onUnmounted(() => {
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
})
</script>

<style scoped>
.ssh-container {
  height: calc(100vh - 60px);
  padding: 0;
  background-color: #f5f7fa;
}

.ssh-layout {
  display: flex;
  height: 100%;
  gap: 0;
}

/* ==================== 左侧主机树区域 ==================== */
.tree-section {
  flex-shrink: 0;
  transition: width 0.1s ease-out;
}

/* ==================== 拖拽分隔条 ==================== */
.resizer {
  width: 8px;
  background: transparent;
  cursor: col-resize;
  position: relative;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease;
}

.resizer:hover {
  background-color: rgba(64, 158, 255, 0.1);
}

.resizer:hover .resizer-line {
  background-color: #409eff;
  width: 3px;
}

.resizer-line {
  width: 1px;
  height: 40px;
  background-color: #dcdfe6;
  border-radius: 2px;
  transition: all 0.2s ease;
}

.resizer:active .resizer-line {
  background-color: #409eff;
  width: 3px;
}

.group-card {
  height: 100%;
  border-radius: 0;
  border: 1px solid #e4e7ed;
  border-top: none;
  border-left: none;
  border-bottom: none;
  overflow: hidden;
  margin: 0;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.search-container {
  display: flex;
  gap: 8px;
  margin-bottom: 15px;
}

.search-input {
  flex: 1;
}

:deep(.search-input .el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #003d99 inset;
}

:deep(.search-input .el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #003d99 inset !important;
}

:deep(.search-input .el-icon) {
  color: #003d99;
}

.refresh-btn {
  flex-shrink: 0;
}

.group-tree {
  max-height: calc(100% - 100px);
  overflow-y: auto;
}

.tree-icon {
  margin-right: 6px;
  font-size: 14px;
  display: inline-flex;
  align-items: center;
  vertical-align: middle;
}

.tree-icon svg,
.tree-icon img {
  width: 16px;
  height: 16px;
}

/* 文件夹图标 */
.folder-icon img {
  width: 14px;
  height: 14px;
}

.business-group {
  font-weight: 600;
  color: #303133;
}

.host-node {
  color: #303133;
  font-weight: 500;
}

.node-label {
  flex: 1;
}

.host-count {
  margin-left: 4px;
  color: #909399;
  font-size: 12px;
  font-weight: 400;
}

/* ==================== 右侧终端区域 ==================== */
.terminal-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
  border-radius: 0;
  border: none;
  overflow: hidden;
}

.terminal-tabs {
  height: 100%;
  display: flex;
  flex-direction: column;
}

:deep(.el-tabs__header) {
  margin: 0;
  border-bottom: 1px solid #e4e7ed;
  background: #f5f7fa;
}

:deep(.el-tabs__content) {
  flex: 1;
  overflow: hidden;
  padding: 0;
}

:deep(.el-tab-pane) {
  height: 100%;
}

:deep(.el-tabs__item) {
  height: 40px;
  line-height: 40px;
  padding: 0 20px;
  color: #606266;
  font-size: 14px;
}

:deep(.el-tabs__item.is-active) {
  background: #ffffff;
  color: #409eff;
  font-weight: 600;
}

:deep(.el-tabs__item:hover) {
  color: #409eff;
}

/* ==================== 欢迎占位页面 ==================== */
.welcome-placeholder {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #1e1e1e;
  gap: 30px;
}

.welcome-icon {
  width: 400px;
  height: 300px;
  opacity: 0.8;
  transition: all 0.3s ease;
}

.welcome-icon:hover {
  opacity: 1;
  transform: scale(1.05);
}

.welcome-icon svg {
  width: 100%;
  height: 100%;
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.1));
}

.welcome-text {
  text-align: center;
}

.welcome-text h1 {
  font-size: 48px;
  font-weight: 700;
  color: #ffffff;
  margin: 0 0 12px 0;
  letter-spacing: 2px;
  font-family: 'Arial', 'Microsoft YaHei', sans-serif;
}

.welcome-text p {
  font-size: 16px;
  color: #b0b0b0;
  margin: 0;
}

/* ==================== 滚动条样式 ==================== */
.group-tree::-webkit-scrollbar {
  width: 6px;
}

.group-tree::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.group-tree::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.group-tree::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* ==================== 树节点样式优化 ==================== */
:deep(.el-tree-node__content) {
  height: 32px;
  border-radius: 4px;
  margin-bottom: 2px;
}

:deep(.el-tree-node__content:hover) {
  background-color: #f5f7fa;
}

:deep(.el-tree-node.is-current > .el-tree-node__content) {
  background-color: #ecf5ff;
  color: #409eff;
}

/* ==================== el-card 样式覆盖 ==================== */
:deep(.el-card__header) {
  padding: 12px 16px;
  border-bottom: 1px solid #e4e7ed;
}

:deep(.el-card__body) {
  padding: 16px;
  height: calc(100% - 49px);
  overflow: hidden;
}
</style>
