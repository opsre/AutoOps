<template>
  <div class="group-tree-section">
    <div class="group-card">
      <!-- 科技感标题区域 -->
      <div class="card-header">
        <div class="title-wrapper">
          <div class="title-icon">
            <div class="pulse-ring"></div>
            <el-icon class="main-icon"><DataBoard /></el-icon>
          </div>
          <div class="title-content">
            <h3 class="title">资产分组</h3>
            <span class="subtitle">Asset Groups</span>
          </div>
        </div>
        <div class="stats-indicator">
          <span class="stats-text">{{ groupList.length }}</span>
          <span class="stats-label">Groups</span>
          <div class="action-buttons-wrapper">
            <div class="toggle-btn" @click="toggleAll" :title="isExpanded ? '折叠全部' : '展开全部'">
              <img
                :src="require('@/assets/image/折叠.svg')"
                class="toggle-icon"
                :class="{ 'expanded': isExpanded }"
                alt="折叠展开"
              />
            </div>
            <div class="refresh-btn" @click="handleRefresh" title="刷新分组">
              <el-icon class="refresh-icon"><Refresh /></el-icon>
            </div>
            <div class="create-btn" @click="handleCreateGroup" title="创建分组">
              <el-icon class="create-icon"><Plus /></el-icon>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 高级搜索区域 -->
      <div class="search-container">
        <div class="search-wrapper">
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
            <div class="search-glow"></div>
          </div>
        </div>
      </div>

      <!-- 树形结构区域 -->
      <div class="tree-container">
        <el-tree
            ref="groupTree"
            :data="groupList"
            :props="defaultProps"
            node-key="id"
            :expanded-keys="expandedKeys"
            :highlight-current="true"
            :default-expanded-keys="expandedKeys"
            @node-click="handleGroupClick"
            @node-expand="handleNodeExpand"
            @node-collapse="handleNodeCollapse"
            class="tech-tree"
        >
          <template v-slot="{ node, data }">
            <div 
              class="tree-node" 
              :class="{ 'parent-node': !data.parentId }"
              @contextmenu.prevent="showContextMenu($event, node, data)"
            >
              <div class="node-content">
                <div class="node-icon-wrapper">
                  <!-- 一级分组：科技感图标 -->
                  <template v-if="!data.parentId">
                    <div class="parent-icon" :class="{ 'expanded': expandedKeys.includes(node.key) }">
                      <div class="icon-bg"></div>
                      <img 
                        v-if="expandedKeys.includes(node.key)"
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
                  <!-- 二级分组：子节点图标 -->
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
                  <div class="connection-line" v-if="!data.parentId"></div>
                </div>
              </div>
              <div class="hover-effect"></div>
            </div>
          </template>
        </el-tree>
      </div>
    </div>
    
    <!-- 右键菜单 -->
    <div 
      v-if="contextMenuVisible" 
      class="context-menu"
      :style="{ left: contextMenuPosition.x + 'px', top: contextMenuPosition.y + 'px' }"
      @click.stop
    >
      <div class="context-menu-item" @click="createRootGroup">
        <el-icon><Plus /></el-icon>
        <span>创建根分组</span>
      </div>
      <div class="context-menu-item" @click="createSubGroup" v-if="contextMenuData.data && contextMenuData.data.id !== 1 && !contextMenuData.data.parentId">
        <el-icon><FolderAdd /></el-icon>
        <span>创建子分组</span>
      </div>
      <div 
        class="context-menu-item" 
        @click="renameGroup" 
        v-if="contextMenuData.data && contextMenuData.data.id !== 1"
      >
        <el-icon><EditPen /></el-icon>
        <span>重命名</span>
      </div>
      <div 
        class="context-menu-item" 
        @click="editGroup" 
        v-if="contextMenuData.data && contextMenuData.data.id !== 1"
      >
        <el-icon><Edit /></el-icon>
        <span>修改分组</span>
      </div>
      <div 
        class="context-menu-item danger" 
        @click="deleteGroup" 
        v-if="contextMenuData.data && contextMenuData.data.id !== 1"
      >
        <el-icon><Delete /></el-icon>
        <span>删除分组</span>
      </div>
    </div>
    
    <!-- 创建/编辑分组对话框 -->
    <el-dialog 
      :title="getDialogTitle()" 
      v-model="groupDialogVisible" 
      width="400px"
      @close="handleDialogClose"
    >
      <el-form 
        :model="groupForm" 
        :rules="groupRules" 
        ref="groupFormRef" 
        label-width="80px"
      >
        <el-form-item label="分组名称" prop="name">
          <el-input 
            v-model="groupForm.name" 
            placeholder="请输入分组名称"
            :disabled="groupForm.id === 1"
          />
        </el-form-item>
        <el-form-item 
          label="父级分组" 
          prop="parentId" 
          v-if="dialogMode === 'create' || dialogMode === 'edit'"
        >
          <el-tree-select
            v-model="groupForm.parentId"
            :data="groupSelectOptions"
            :props="{
              label: 'name',
              value: 'id',
              disabled: (data) => data.id === groupForm.id || data.id === 1 || data.parentId
            }"
            check-strictly
            placeholder="请选择父级分组（不选择为根分组）"
            clearable
            :disabled="groupForm.id === 1"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="groupDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitGroupForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'CmdbGroup',
  props: {
    groupList: {
      type: Array,
      required: true
    },
    expandedKeys: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      groupSearchText: '',
      isExpanded: false, // 默认折叠状态
      defaultProps: {
        children: 'children',
        label: 'name'
      },
      // 右键菜单相关
      contextMenuVisible: false,
      contextMenuPosition: { x: 0, y: 0 },
      contextMenuData: {},
      // 分组管理对话框
      groupDialogVisible: false,
      dialogMode: 'create', // create, edit
      groupForm: {
        id: null,
        name: '',
        parentId: null
      },
      groupRules: {
        name: [{ required: true, message: '请输入分组名称', trigger: 'blur' }]
      },
      groupSelectOptions: []
    }
  },
  watch: {
    expandedKeys: {
      handler(newKeys) {
        this.$nextTick(() => {
          const tree = this.$refs.groupTree
          if (tree) {
            // 如果是空数组，先获取所有展开的节点然后逐个折叠
            if (newKeys.length === 0) {
              const currentExpanded = tree.store.nodesMap
              Object.keys(currentExpanded).forEach(key => {
                const node = currentExpanded[key]
                if (node.expanded) {
                  node.collapse()
                }
              })
            } else {
              tree.setExpandedKeys(newKeys)
            }
          }
        })
      },
      deep: true,
      immediate: true
    },
    groupList: {
      handler(newList) {
        this.groupSelectOptions = this.formatGroupOptions(newList)
      },
      deep: true,
      immediate: true
    }
  },
  computed: {
    // 格式化分组选项供选择器使用
    formatGroupOptions() {
      return (groups) => {
        return groups.map(group => ({
          ...group,
          children: group.children && group.children.length > 0 ? this.formatGroupOptions(group.children) : []
        }))
      }
    }
  },
  mounted() {
    // 点击其他地方隐藏右键菜单
    document.addEventListener('click', this.hideContextMenu)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.hideContextMenu)
  },
  methods: {
    handleGroupSearch() {
      this.$emit('group-search', this.groupSearchText)
    },
    handleGroupClick(node, element) {
      this.$emit('group-click', node, element)
    },
    handleNodeExpand(data, node) {
      this.$emit('node-expand', data, node)
    },
    handleNodeCollapse(data, node) {
      this.$emit('node-collapse', data, node)
    },
    toggleAll() {
      this.isExpanded = !this.isExpanded
      if (this.isExpanded) {
        this.$emit('expand-all')
      } else {
        this.$emit('collapse-all')
      }
    },

    // 刷新分组
    handleRefresh() {
      this.$emit('refresh-groups')
    },

    // 创建分组
    handleCreateGroup() {
      this.createRootGroup()
    },
    
    // 显示右键菜单
    showContextMenu(event, node, data) {
      event.preventDefault()
      event.stopPropagation()
      
      // 获取分组卡片容器的位置信息
      const groupCard = this.$el.querySelector('.group-card')
      const cardRect = groupCard.getBoundingClientRect()
      
      // 计算相对于分组卡片的位置
      this.contextMenuPosition = {
        x: event.clientX - cardRect.left + 10, // 相对于卡片的位置 + 10px偏移
        y: event.clientY - cardRect.top + 5    // 相对于卡片的位置 + 5px偏移
      }
      
      console.log('菜单位置:', this.contextMenuPosition)
      console.log('卡片位置:', cardRect)
      console.log('鼠标位置:', event.clientX, event.clientY)
      
      this.contextMenuData = { node, data }
      this.contextMenuVisible = true
    },
    
    // 隐藏右键菜单
    hideContextMenu() {
      this.contextMenuVisible = false
    },
    
    // 创建根分组
    createRootGroup() {
      this.dialogMode = 'create'
      this.groupForm = {
        id: null,
        name: '',
        parentId: null
      }
      this.groupDialogVisible = true
      this.hideContextMenu()
    },
    
    // 创建子分组
    createSubGroup() {
      this.dialogMode = 'create'
      this.groupForm = {
        id: null,
        name: '',
        parentId: this.contextMenuData.data.id
      }
      this.groupDialogVisible = true
      this.hideContextMenu()
    },
    
    // 重命名分组
    renameGroup() {
      this.dialogMode = 'rename'
      this.groupForm = {
        id: this.contextMenuData.data.id,
        name: this.contextMenuData.data.name,
        parentId: this.contextMenuData.data.parentId || null
      }
      this.groupDialogVisible = true
      this.hideContextMenu()
    },
    
    // 编辑分组
    editGroup() {
      this.dialogMode = 'edit'
      this.groupForm = {
        id: this.contextMenuData.data.id,
        name: this.contextMenuData.data.name,
        parentId: this.contextMenuData.data.parentId || null
      }
      this.groupDialogVisible = true
      this.hideContextMenu()
    },
    
    // 删除分组
    async deleteGroup() {
      const groupData = this.contextMenuData.data
      
      try {
        await this.$confirm(
          `确定要删除分组"${groupData.name}"吗？如果该分组下有主机，将无法删除。`,
          '确认删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        this.$emit('delete-group', groupData.id)
        this.hideContextMenu()
      } catch {
        // 用户取消删除
      }
    },
    
    // 获取对话框标题
    getDialogTitle() {
      switch (this.dialogMode) {
        case 'create':
          return '创建分组'
        case 'rename':
          return '重命名分组'
        case 'edit':
          return '编辑分组'
        default:
          return '分组管理'
      }
    },
    
    // 提交分组表单
    async submitGroupForm() {
      try {
        await this.$refs.groupFormRef.validate()
        
        const formData = { 
          ...this.groupForm,
          // 如果parentId为null或undefined，设为0表示根分组
          parentId: this.groupForm.parentId || 0
        }
        
        console.log('提交分组数据:', formData)
        
        if (this.dialogMode === 'create') {
          this.$emit('create-group', formData)
        } else if (this.dialogMode === 'edit' || this.dialogMode === 'rename') {
          this.$emit('update-group', formData)
        }
        
        this.groupDialogVisible = false
      } catch (error) {
        console.error('表单验证失败:', error)
      }
    },
    
    // 对话框关闭处理
    handleDialogClose() {
      this.$refs.groupFormRef?.resetFields()
    }
  }
}
</script>

<style scoped>
/* 🚀 现代化科技感分组树样式 */

.group-tree-section {
  width: 280px;
  margin-right: 20px;
  position: relative;
}

/* 🎨 高科技卡片设计 */
.group-card {
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.95), rgba(255, 255, 255, 0.85));
  backdrop-filter: blur(20px);
  border-radius: 20px;
  box-shadow: 
    0 20px 40px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
  height: 100%;
  overflow: hidden;
  position: relative;
}

.group-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.5), transparent);
  z-index: 1;
}

/* 🎯 科技感标题区域 */
.card-header {
  padding: 20px 20px 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  position: relative;
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

.main-icon {
  font-size: 20px;
  color: #667eea;
  z-index: 2;
  background: rgba(102, 126, 234, 0.1);
  padding: 8px;
  border-radius: 8px;
}

@keyframes pulse {
  0% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.2); opacity: 0.7; }
  100% { transform: scale(1.4); opacity: 0; }
}

.title-content .title {
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 2px 0;
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  white-space: nowrap;
}

.subtitle {
  font-size: 11px;
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
  position: relative;
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

/* 🔍 高级搜索区域 */
.search-container {
  padding: 15px 20px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.08);
}

.search-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.search-input-wrapper {
  position: relative;
}

.tech-input :deep(.el-input__wrapper) {
  background: rgba(102, 126, 234, 0.03);
  border: 1px solid rgba(102, 126, 234, 0.15);
  border-radius: 12px;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.02);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
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
  transition: all 0.3s ease;
}

.tech-input:focus-within .search-icon {
  color: #667eea;
}

/* 操作按钮容器 */
.action-buttons-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  position: absolute;
  bottom: -21px;
  right: 25px;
}

.toggle-btn {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 3px;
  transition: all 0.3s ease;
  background: transparent;
  border: none;
  z-index: 10;
}

.toggle-btn:hover {
  transform: scale(1.1);
}

.toggle-icon {
  width: 20px;
  height: 20px;
  transition: all 0.3s ease;
  filter: hue-rotate(30deg) saturate(1.2) brightness(1.1);
}

.toggle-icon.expanded {
  transform: rotate(180deg);
  filter: hue-rotate(60deg) saturate(1.4) brightness(1.2);
}

/* 刷新按钮样式 */
.refresh-btn {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.3s ease;
  background: rgba(102, 126, 234, 0.15);
  border: 1px solid rgba(102, 126, 234, 0.2);
  z-index: 10;
}

.refresh-btn:hover {
  transform: scale(1.1);
  background: rgba(102, 126, 234, 0.25);
  border-color: rgba(102, 126, 234, 0.35);
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

/* 创建按钮样式 */
.create-btn {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.3s ease;
  background: rgba(103, 194, 58, 0.15);
  border: 1px solid rgba(103, 194, 58, 0.2);
  z-index: 10;
}

.create-btn:hover {
  transform: scale(1.1);
  background: rgba(103, 194, 58, 0.25);
  border-color: rgba(103, 194, 58, 0.35);
}

.create-btn:active {
  transform: scale(0.95);
}

.create-icon {
  font-size: 16px;
  color: #67C23A;
  transition: all 0.3s ease;
}

.create-btn:hover .create-icon {
  color: #85CE61;
  transform: rotate(90deg);
}

/* 🌳 高科技树形结构 */
.tree-container {
  padding: 15px 0;
  max-height: calc(100vh - 300px);
  overflow-y: auto;
  position: relative;
}

/* 移除装饰性虚线 */

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

/* 调整树节点展开按钮位置和旋转效果 - 放在图标前面 */
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
}

.tree-node:hover {
  transform: none;
}

/* 只有直接悬停父节点时才显示阴影效果 */
.tree-node.parent-node:hover > .hover-effect {
  opacity: 1;
  transform: scale(1);
}

/* 确保子节点不受父节点阴影影响 */
.tree-node:not(.parent-node) .hover-effect {
  display: none;
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

.parent-icon.expanded .icon-bg {
  display: none;
}

.parent-icon-img {
  width: 16px;
  height: 16px;
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

.child-icon-img {
  width: 18px;
  height: 18px;
}

.node-text {
  flex: 1;
  position: relative;
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

.host-count {
  font-size: 12px;
  font-weight: 400;
  color: rgba(102, 126, 234, 0.7);
  margin-left: 4px;
}

.connection-line {
  position: absolute;
  bottom: -6px;
  left: 0;
  height: 1px;
  width: 0;
  background: linear-gradient(90deg, #667eea, transparent);
  transition: width 0.3s ease;
}

.parent-node:hover .connection-line {
  width: 0;
}

.node-status {
  display: flex;
  align-items: center;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  position: relative;
}

.status-dot.online {
  background: #67C23A;
  box-shadow: 0 0 6px rgba(103, 194, 58, 0.5);
}

.status-dot.online::before {
  content: '';
  position: absolute;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #67C23A;
  animation: ping 2s infinite;
}

@keyframes ping {
  75%, 100% {
    transform: scale(2);
    opacity: 0;
  }
}

/* 🎯 选中状态 - 只针对被选中的节点本身 */
.tech-tree :deep(.el-tree-node.is-current > .el-tree-node__content) .tree-node {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.08), rgba(118, 75, 162, 0.08));
  border: none;
  transform: none;
}

.tech-tree :deep(.el-tree-node.is-current > .el-tree-node__content) .hover-effect {
  opacity: 0 !important;
}

.tech-tree :deep(.el-tree-node.is-current > .el-tree-node__content) .node-label {
  color: #2c3e50;
  font-weight: 500;
}

/* 📱 滚动条样式 */
.tree-container::-webkit-scrollbar {
  width: 4px;
}

.tree-container::-webkit-scrollbar-track {
  background: rgba(102, 126, 234, 0.05);
  border-radius: 2px;
}

.tree-container::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.3), rgba(118, 75, 162, 0.3));
  border-radius: 2px;
}

.tree-container::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.5), rgba(118, 75, 162, 0.5));
}

/* 🎮 交互反馈增强 */
.group-card {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.group-card:hover {
  transform: translateY(-2px);
  box-shadow: 
    0 25px 50px rgba(0, 0, 0, 0.15),
    0 0 0 1px rgba(255, 255, 255, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.stats-indicator {
  transition: all 0.3s ease;
}

.stats-indicator:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

/* 🚀 数据加载动画 */
@keyframes dataLoad {
  0% { opacity: 0; transform: translateY(20px); }
  100% { opacity: 1; transform: translateY(0); }
}

.tree-node {
  animation: dataLoad 0.4s ease-out;
}

.tree-node:nth-child(1) { animation-delay: 0.1s; }
.tree-node:nth-child(2) { animation-delay: 0.2s; }
.tree-node:nth-child(3) { animation-delay: 0.3s; }
.tree-node:nth-child(4) { animation-delay: 0.4s; }
.tree-node:nth-child(5) { animation-delay: 0.5s; }

/* 操作按钮区域 */
.action-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}

.manage-btn {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 3px;
  transition: all 0.3s ease;
  background: transparent;
  border: none;
  z-index: 10;
}

.manage-btn:hover {
  transform: scale(1.1);
  background: rgba(102, 126, 234, 0.1);
}

.manage-icon {
  font-size: 16px;
  color: #667eea;
}

/* 右键菜单样式 */
.context-menu {
  position: absolute;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 4px;
  z-index: 9999;
  width: 130px;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s ease;
  font-size: 14px;
  color: #2c3e50;
}

.context-menu-item:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.context-menu-item.danger {
  color: #f56c6c;
}

.context-menu-item.danger:hover {
  background: rgba(245, 108, 108, 0.1);
  color: #f56c6c;
}

.context-menu-item .el-icon {
  font-size: 14px;
}
</style>
