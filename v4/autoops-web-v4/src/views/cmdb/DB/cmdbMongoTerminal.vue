<template>
  <div class="mongo-terminal-page">
    <el-container class="terminal-container">
      <!-- 左侧树形结构 -->
      <el-aside :width="asideWidth + 'px'" class="tree-aside">
        <el-card shadow="never" class="tree-card">
          <template #header>
            <div class="tree-header">
              <span>MongoDB 资源</span>
              <el-button text size="small" @click="refreshTree" :loading="loadingTree">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
          </template>

          <el-input v-model="treeFilter" size="small" placeholder="过滤..." clearable class="mb10">
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>

          <el-scrollbar height="calc(100vh - 200px)">
            <el-tree
              ref="treeRef"
              :data="treeData"
              :props="treeProps"
              :filter-node-method="filterNode"
              node-key="id"
              highlight-current
              :current-node-key="selectedNode && selectedNode.id"
              :expand-on-click-node="false"
              @node-click="handleNodeClick"
              @node-expand="handleNodeExpand"
              @node-contextmenu="handleNodeRightClick"
            >
              <template #default="{ node, data }">
                <span class="tree-node">
                  <el-icon v-if="data.type === 'instance'">
                    <img :src="mongoIcon" class="node-icon" />
                  </el-icon>
                  <el-icon v-else-if="data.type === 'database'"><img src="@/assets/image/数据库.svg" class="node-icon" /></el-icon>
                  <el-icon v-else-if="data.type === 'collection'"><img src="@/assets/image/数据表.svg" class="node-icon" /></el-icon>
                  <el-icon v-else><Database /></el-icon>
                  <span class="node-label">{{ node.label }}</span>
                  <span v-if="data.docCount !== undefined" class="node-count">({{ data.docCount }})</span>
                </span>
              </template>
            </el-tree>
          </el-scrollbar>

          <!-- 右键菜单 -->
          <div
            v-show="contextMenuVisible"
            ref="contextMenuRef"
            class="context-menu"
            :style="{ left: contextMenuX + 'px', top: contextMenuY + 'px' }"
            @click.stop
          >
            <!-- 实例右键菜单 -->
            <template v-if="contextMenuNode && contextMenuNode.type === 'instance'">
              <div class="context-menu-item" @click="handleRefreshDatabases">
                <el-icon><Refresh /></el-icon>
                <span>刷新数据库</span>
              </div>
              <div class="context-menu-item" @click="handleAddDatabase">
                <el-icon><Plus /></el-icon>
                <span>创建数据库</span>
              </div>
            </template>

            <!-- 数据库右键菜单 -->
            <template v-if="contextMenuNode && contextMenuNode.type === 'database'">
              <div class="context-menu-item" @click="handleRefreshCollections">
                <el-icon><Refresh /></el-icon>
                <span>刷新集合</span>
              </div>
              <div class="context-menu-item" @click="handleAddCollection">
                <el-icon><Plus /></el-icon>
                <span>创建集合</span>
              </div>
              <div class="context-menu-item" @click="handleDropDatabase">
                <el-icon><Delete /></el-icon>
                <span>删除数据库</span>
              </div>
            </template>

            <!-- 集合右键菜单 -->
            <template v-if="contextMenuNode && contextMenuNode.type === 'collection'">
              <div class="context-menu-item" @click="handleCollectionStats">
                <el-icon><View /></el-icon>
                <span>集合统计</span>
              </div>
              <div class="context-menu-item" @click="handleDropCollection">
                <el-icon><Delete /></el-icon>
                <span>删除集合</span>
              </div>
              <div class="context-menu-item" @click="handleCopyCollectionName">
                <el-icon><CopyDocument /></el-icon>
                <span>复制名字</span>
              </div>
            </template>
          </div>
        </el-card>
      </el-aside>

      <!-- Draggable resizer -->
      <div class="aside-resizer" @mousedown="startAsideResize">
        <el-icon class="resizer-icon"><DArrowRight /></el-icon>
      </div>

      <!-- 右侧内容区域 -->
      <el-main class="main-content" ref="mainContent">
        <el-card class="main-card" shadow="never">
          <div v-if="openTabs.length > 0">
            <el-tabs v-model="activeTab" type="card" closable @tab-remove="handleTabRemove">
              <el-tab-pane v-for="tab in openTabs" :key="tab.name" :name="tab.name">
                <template #label>
                  <el-tooltip :content="tab.title" placement="top">
                    <span class="tab-label-text">{{ getTabShortLabel(tab) }}</span>
                  </el-tooltip>
                </template>

                <div class="content-header">
                  <div class="header-info">
                    <el-icon><img src="@/assets/image/数据表.svg" class="node-icon" /></el-icon>
                    <span class="header-title">{{ tab.title }}</span>
                  </div>
                  <div v-if="tab.connection" class="header-connection-info">
                    <span style="margin-right: 15px;"><strong>实例:</strong> <span class="info-value">{{ tab.connection.instanceName }}</span></span>
                    <span style="margin-right: 15px;"><strong>数据库:</strong> <span class="info-value">{{ tab.connection.database }}</span></span>
                    <span style="margin-right: 15px;"><strong>集合:</strong> <span class="info-value">{{ tab.connection.collection }}</span></span>
                    <span style="margin-right: 15px;" v-if="tab.meta"><strong>文档数:</strong> <span class="info-value">{{ tab.meta.docCount }}</span></span>
                  </div>
                  <div class="header-actions">
                    <el-button type="primary" size="small" @click="handleAddDoc(tab)">
                      <el-icon><Plus /></el-icon>
                      新增文档
                    </el-button>
                    <el-button type="success" size="small" @click="reloadCollectionData(tab)" :loading="tab.loading">
                      <el-icon><Refresh /></el-icon>
                      刷新
                    </el-button>
                  </div>
                </div>

                <!-- 查询条件 -->
                <div class="query-section">
                  <el-form :inline="true" size="small">
                    <el-form-item label="查询条件">
                      <el-input v-model="tab.queryText" placeholder='{"field": "value"}' style="width: 350px" clearable>
                        <template #prepend>Filter</template>
                      </el-input>
                    </el-form-item>
                    <el-form-item label="排序">
                      <el-input v-model="tab.sortText" placeholder='{"_id": -1}' style="width: 200px" clearable>
                        <template #prepend>Sort</template>
                      </el-input>
                    </el-form-item>
                    <el-form-item label="限制">
                      <el-input-number v-model="tab.limit" :min="1" :max="1000" style="width: 120px" />
                    </el-form-item>
                    <el-form-item>
                      <el-button type="primary" @click="handleQuery(tab)">查询</el-button>
                      <el-button @click="resetQuery(tab)">重置</el-button>
                    </el-form-item>
                  </el-form>
                </div>

                <!-- 数据表格 -->
                <div>
                  <el-table
                    v-if="tab.rows && tab.rows.length"
                    :data="tab.rows"
                    border
                    stripe
                    max-height="450"
                    size="small"
                    :cell-style="{ padding: '4px 0' }"
                    :header-cell-style="{ padding: '8px 0' }"
                  >
                    <el-table-column v-for="col in tab.columns" :key="col.prop" :prop="col.prop" :label="col.label" :min-width="col.minWidth || 120" show-overflow-tooltip>
                      <template #default="{ row }">
                        <span>{{ formatValue(row[col.prop]) }}</span>
                      </template>
                    </el-table-column>
                    <el-table-column fixed="right" label="操作" width="200">
                      <template #default="{ row }">
                        <el-button size="small" type="primary" text @click="handleViewDoc(row, tab)">查看</el-button>
                        <el-button size="small" type="warning" text @click="handleEditDoc(row, tab)">编辑</el-button>
                        <el-button size="small" type="danger" text @click="handleDeleteDoc(row, tab)">删除</el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                  <el-empty v-else description="暂无数据" />
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
          <el-empty v-else description="请从左侧选择一个集合" />
        </el-card>
      </el-main>
    </el-container>

    <!-- 查看/编辑文档对话框 -->
    <el-dialog v-model="docDialogVisible" :title="docDialogTitle" width="60%" destroy-on-close>
      <div class="code-editor-wrapper">
        <CodeEditor v-model="docContent" language="json" :readonly="docReadonly" />
      </div>
      <template #footer>
        <el-button @click="docDialogVisible = false">取消</el-button>
        <el-button v-if="!docReadonly" type="primary" @click="saveDoc">保存</el-button>
      </template>
    </el-dialog>

    <!-- 创建数据库对话框 -->
    <el-dialog v-model="addDatabaseDialogVisible" title="创建数据库" width="40%" destroy-on-close>
      <el-form :model="databaseForm" :rules="databaseRules" ref="databaseFormRef" label-width="100px">
        <el-form-item label="数据库名" prop="dbName">
          <el-input v-model="databaseForm.dbName" placeholder="请输入数据库名称" />
        </el-form-item>
        <el-form-item label="集合名" prop="collectionName">
          <el-input v-model="databaseForm.collectionName" placeholder="请输入初始集合名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDatabaseDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAddDatabase">确定</el-button>
      </template>
    </el-dialog>

    <!-- 创建集合对话框 -->
    <el-dialog v-model="addCollectionDialogVisible" title="创建集合" width="40%" destroy-on-close>
      <el-form :model="collectionForm" :rules="collectionRules" ref="collectionFormRef" label-width="100px">
        <el-form-item label="集合名" prop="collectionName">
          <el-input v-model="collectionForm.collectionName" placeholder="请输入集合名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addCollectionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAddCollection">确定</el-button>
      </template>
    </el-dialog>

    <!-- 集合统计对话框 -->
    <el-dialog v-model="statsDialogVisible" title="集合统计" width="60%" destroy-on-close>
      <div class="code-editor-wrapper">
        <CodeEditor v-model="statsContent" language="json" :readonly="true" />
      </div>
      <template #footer>
        <el-button type="primary" @click="statsDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { Refresh, Search, DArrowRight, Database, Plus, View, Delete, CopyDocument } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { mongoAPI } from '@/api/database'
import CodeEditor from '@/components/CodeEditor.vue'

export default {
  components: { Refresh, Search, DArrowRight, Database, CodeEditor, Plus, View, Delete, CopyDocument },
  data() {
    return {
      mongoIcon: require('@/assets/image/mongodb.svg'),
      asideWidth: 300,
      resizing: false,
      startX: 0,
      startWidth: 300,
      loadingTree: false,
      treeFilter: '',
      treeData: [],
      treeProps: { children: 'children', label: 'label' },
      selectedNode: null,
      activeTab: '',
      openTabs: [],
      docDialogVisible: false,
      docDialogTitle: '',
      docContent: '',
      docReadonly: false,
      currentDoc: null,
      currentTab: null,
      // 右键菜单
      contextMenuVisible: false,
      contextMenuX: 0,
      contextMenuY: 0,
      contextMenuNode: null,
      // 创建数据库
      addDatabaseDialogVisible: false,
      databaseForm: { dbName: '', collectionName: '' },
      databaseRules: {
        dbName: [{ required: true, message: '请输入数据库名称', trigger: 'blur' }],
        collectionName: [{ required: true, message: '请输入集合名称', trigger: 'blur' }]
      },
      // 创建集合
      addCollectionDialogVisible: false,
      collectionForm: { collectionName: '' },
      collectionRules: {
        collectionName: [{ required: true, message: '请输入集合名称', trigger: 'blur' }]
      },
      // 集合统计
      statsDialogVisible: false,
      statsContent: ''
    }
  },
  watch: {
    treeFilter(val) {
      this.$refs.treeRef && this.$refs.treeRef.filter(val)
    }
  },
  mounted() {
    document.addEventListener('click', this.closeContextMenu)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.closeContextMenu)
  },
  created() {
    this.loadInstances()
  },
  methods: {
    async refreshTree() {
      await this.loadInstances()
    },
    filterNode(value, data) {
      if (!value) return true
      return String(data.label || '').toLowerCase().includes(value.toLowerCase())
    },
    startAsideResize(e) {
      this.resizing = true
      this.startX = e.clientX
      this.startWidth = this.asideWidth
      document.addEventListener('mousemove', this.onAsideResize)
      document.addEventListener('mouseup', this.stopAsideResize)
    },
    onAsideResize(e) {
      if (!this.resizing) return
      const delta = e.clientX - this.startX
      this.asideWidth = Math.min(600, Math.max(220, this.startWidth + delta))
    },
    stopAsideResize() {
      this.resizing = false
      document.removeEventListener('mousemove', this.onAsideResize)
      document.removeEventListener('mouseup', this.stopAsideResize)
    },
    async loadInstances() {
      this.loadingTree = true
      try {
        const res = await mongoAPI.getInstanceList({ pageNum: 1, pageSize: 100 })
        if (res.data?.code === 0 || res.data?.code === 200) {
          let instances = res.data?.data
          if (!Array.isArray(instances)) instances = instances?.list || instances?.data || []
          this.treeData = instances.map(inst => ({
            id: `inst-${inst.id}`,
            label: inst.name,
            type: 'instance',
            instanceId: inst.id,
            children: [],
            isLeaf: false
          }))
        } else {
          ElMessage.error(res.data?.message || '加载实例列表失败')
        }
      } catch (e) {
        ElMessage.error('加载实例列表失败')
      } finally {
        this.loadingTree = false
      }
    },
    async handleNodeExpand(data) {
      if (data.type === 'instance' && data.children.length === 0) {
        await this.loadDatabases(data)
      } else if (data.type === 'database' && data.children.length === 0) {
        await this.loadCollections(data)
      }
    },
    async handleNodeClick(data) {
      this.selectedNode = data
      if (data.type === 'instance') {
        if (data.children.length === 0) await this.loadDatabases(data)
      } else if (data.type === 'database') {
        if (data.children.length === 0) await this.loadCollections(data)
      } else if (data.type === 'collection') {
        this.openCollectionTab(data)
      }
    },
    async loadDatabases(instNode) {
      try {
        const res = await mongoAPI.getDatabases(instNode.instanceId)
        console.log('loadDatabases response:', res)

        let databases = []
        if (res.data?.code === 0 || res.data?.code === 200) {
          const dbData = res.data.data?.Databases || res.data.data?.databases || []
          databases = dbData
        }

        instNode.children = databases.map(db => ({
          id: `db-${instNode.instanceId}-${db.Name || db.name}`,
          label: db.Name || db.name,
          type: 'database',
          instanceId: instNode.instanceId,
          dbName: db.Name || db.name,
          children: [],
          isLeaf: false
        }))
      } catch (e) {
        console.error('loadDatabases error:', e)
        ElMessage.error('加载数据库失败: ' + (e.message || ''))
      }
    },
    async loadCollections(dbNode) {
      try {
        const res = await mongoAPI.getCollections(dbNode.instanceId, {
          database: dbNode.dbName
        })
        console.log('loadCollections response:', res)

        let collections = []
        if (res.data?.code === 0 || res.data?.code === 200) {
          collections = res.data.data || []
        }

        dbNode.children = collections.map(collName => ({
          id: `coll-${dbNode.instanceId}-${dbNode.dbName}-${collName}`,
          label: collName,
          type: 'collection',
          instanceId: dbNode.instanceId,
          dbName: dbNode.dbName,
          collectionName: collName,
          isLeaf: true
        }))
      } catch (e) {
        console.error('loadCollections error:', e)
        ElMessage.error('加载集合失败: ' + (e.message || ''))
      }
    },
    getTabShortLabel(tab) {
      const label = tab?.title || '未命名'
      if (label.length > 24) {
        return label.slice(0, 12) + '…' + label.slice(-8)
      }
      return label
    },
    handleTabRemove(targetName) {
      const idx = this.openTabs.findIndex(t => t.name === targetName)
      if (idx >= 0) {
        this.openTabs.splice(idx, 1)
        if (this.activeTab === targetName) {
          this.activeTab = this.openTabs.length ? this.openTabs[0].name : ''
        }
      }
    },
    async openCollectionTab(node) {
      const tabName = `coll-${node.instanceId}-${node.dbName}-${node.collectionName}`
      const exist = this.openTabs.find(t => t.name === tabName)
      if (exist) {
        this.activeTab = tabName
        if (!exist.rows || exist.rows.length === 0) {
          await this.reloadCollectionData(exist)
        }
        return
      }

      let instanceName = this.treeData.find(i => i.instanceId === node.instanceId)?.label || `#${node.instanceId}`
      const tab = {
        name: tabName,
        title: `${node.dbName}.${node.collectionName}`,
        node,
        connection: { instanceName, database: node.dbName, collection: node.collectionName },
        columns: [],
        rows: [],
        meta: null,
        loading: false,
        queryText: '',
        sortText: '',
        limit: 20
      }
      this.openTabs.push(tab)
      this.activeTab = tabName
      const reactiveTab = this.openTabs[this.openTabs.length - 1]
      await this.reloadCollectionData(reactiveTab)
    },
    async reloadCollectionData(tab) {
      await this.handleQuery(tab)
    },
    async handleQuery(tab) {
      const { instanceId, dbName, collectionName } = tab.node
      tab.loading = true
      try {
        let filter = {}
        if (tab.queryText) {
          try {
            filter = JSON.parse(tab.queryText)
          } catch {
            ElMessage.error('查询条件格式错误')
            return
          }
        }

        let sort = {}
        if (tab.sortText) {
          try {
            sort = JSON.parse(tab.sortText)
          } catch {
            ElMessage.error('排序条件格式错误')
            return
          }
        }

        const requestData = {
          database: dbName,
          collection: collectionName,
          filter,
          limit: tab.limit || 20,
          skip: 0
        }

        if (Object.keys(sort).length > 0) {
          requestData.sort = sort
        }

        const res = await mongoAPI.findDocuments(instanceId, requestData)
        console.log('handleQuery response:', res)

        let documents = []
        if (res.data?.code === 0 || res.data?.code === 200) {
          documents = res.data.data || []
        }

        // 构建列
        const allFields = new Set(['_id'])
        documents.forEach(doc => {
          Object.keys(doc).forEach(k => allFields.add(k))
        })

        tab.columns = Array.from(allFields).map(field => ({
          prop: field,
          label: field,
          minWidth: field === '_id' ? 200 : 150
        }))

        tab.rows = documents.map(doc => ({ ...doc, _rawDoc: doc }))
        tab.meta = { docCount: documents.length }
      } catch (e) {
        console.error('handleQuery error:', e)
        ElMessage.error('查询失败: ' + (e.message || ''))
        tab.rows = []
      } finally {
        tab.loading = false
      }
    },
    resetQuery(tab) {
      tab.queryText = ''
      tab.sortText = ''
      tab.limit = 20
      this.handleQuery(tab)
    },
    formatValue(val) {
      if (val === null) return 'null'
      if (val === undefined) return 'undefined'
      if (typeof val === 'object') return JSON.stringify(val)
      return String(val)
    },
    handleViewDoc(row, tab) {
      this.docDialogTitle = `查看文档 - ${row._id}`
      this.docContent = JSON.stringify(row._rawDoc, null, 2)
      this.docReadonly = true
      this.currentDoc = row
      this.currentTab = tab
      this.docDialogVisible = true
    },
    handleEditDoc(row, tab) {
      this.docDialogTitle = `编辑文档 - ${row._id}`
      this.docContent = JSON.stringify(row._rawDoc, null, 2)
      this.docReadonly = false
      this.currentDoc = row
      this.currentTab = tab
      this.docDialogVisible = true
    },
    handleAddDoc(tab) {
      this.docDialogTitle = '新增文档'
      this.docContent = JSON.stringify({}, null, 2)
      this.docReadonly = false
      this.currentDoc = null
      this.currentTab = tab
      this.docDialogVisible = true
    },
    async saveDoc() {
      try {
        const newDoc = JSON.parse(this.docContent)
        const { instanceId, dbName, collectionName } = this.currentTab.node

        // 如果 currentDoc 为 null，说明是新增文档
        if (!this.currentDoc) {
          const res = await mongoAPI.insertDocument(instanceId, {
            database: dbName,
            collection: collectionName,
            document: newDoc
          })
          console.log('insertDoc response:', res)

          if (res.data?.code === 0 || res.data?.code === 200) {
            ElMessage.success('新增成功')
            this.docDialogVisible = false
            await this.reloadCollectionData(this.currentTab)
          } else {
            ElMessage.error(res.data?.message || '新增失败')
          }
        } else {
          // 编辑已有文档
          const requestData = {
            database: dbName,
            collection: collectionName,
            filter: { _id: this.currentDoc._id },
            update: newDoc
          }

          const res = await mongoAPI.updateDocument(instanceId, requestData)
          console.log('saveDoc response:', res)

          if (res.data?.code === 0 || res.data?.code === 200) {
            ElMessage.success('保存成功')
            this.docDialogVisible = false
            await this.reloadCollectionData(this.currentTab)
          } else {
            ElMessage.error(res.data?.message || '保存失败')
          }
        }
      } catch (e) {
        console.error('saveDoc error:', e)
        ElMessage.error('保存失败: ' + (e.message || ''))
      }
    },
    async handleDeleteDoc(row, tab) {
      try {
        await ElMessageBox.confirm(`确定删除文档 ${row._id} 吗？`, '删除确认', { type: 'warning' })
        const { instanceId, dbName, collectionName } = tab.node

        const requestData = {
          database: dbName,
          collection: collectionName,
          filter: { _id: row._id }
        }

        const res = await mongoAPI.deleteDocument(instanceId, requestData)
        console.log('handleDeleteDoc response:', res)

        if (res.data?.code === 0 || res.data?.code === 200) {
          ElMessage.success('删除成功')
          await this.reloadCollectionData(tab)
        } else {
          ElMessage.error(res.data?.message || '删除失败')
        }
      } catch (e) {
        console.error('handleDeleteDoc error:', e)
        if (e !== 'cancel') ElMessage.error('删除失败: ' + (e.message || ''))
      }
    },
    // 右键菜单
    handleNodeRightClick(event, data, node, element) {
      if (data.type !== 'instance' && data.type !== 'database' && data.type !== 'collection') return

      event.preventDefault()
      event.stopPropagation()

      this.contextMenuNode = data
      this.contextMenuX = event.clientX
      this.contextMenuY = event.clientY
      this.contextMenuVisible = true
    },
    closeContextMenu() {
      this.contextMenuVisible = false
    },
    async handleRefreshDatabases() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false
      await this.loadDatabases(this.contextMenuNode)
      ElMessage.success('刷新成功')
    },
    handleAddDatabase() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false
      this.databaseForm = { dbName: '', collectionName: '' }
      this.addDatabaseDialogVisible = true
    },
    async submitAddDatabase() {
      try {
        await this.$refs.databaseFormRef.validate()

        const { instanceId } = this.contextMenuNode
        const { dbName, collectionName } = this.databaseForm

        // MongoDB 创建数据库需要先创建集合
        const requestData = {
          database: dbName,
          collection: collectionName
        }

        const res = await mongoAPI.insertDocument(instanceId, {
          database: dbName,
          collection: collectionName,
          document: { _initialized: true }
        })
        console.log('submitAddDatabase response:', res)

        if (res.data?.code === 0 || res.data?.code === 200) {
          ElMessage.success('数据库创建成功')
          this.addDatabaseDialogVisible = false
          await this.loadDatabases(this.contextMenuNode)
        } else {
          ElMessage.error(res.data?.message || '创建失败')
        }
      } catch (e) {
        console.error('submitAddDatabase error:', e)
        if (e !== false) {
          ElMessage.error('创建失败: ' + (e.message || ''))
        }
      }
    },
    async handleRefreshCollections() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false
      await this.loadCollections(this.contextMenuNode)
      ElMessage.success('刷新成功')
    },
    handleAddCollection() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false
      this.collectionForm = { collectionName: '' }
      this.addCollectionDialogVisible = true
    },
    async submitAddCollection() {
      try {
        await this.$refs.collectionFormRef.validate()

        const { instanceId, dbName } = this.contextMenuNode
        const { collectionName } = this.collectionForm

        // 通过插入一个初始化文档来创建集合
        const res = await mongoAPI.insertDocument(instanceId, {
          database: dbName,
          collection: collectionName,
          document: { _initialized: true }
        })
        console.log('submitAddCollection response:', res)

        if (res.data?.code === 0 || res.data?.code === 200) {
          ElMessage.success('集合创建成功')
          this.addCollectionDialogVisible = false
          await this.loadCollections(this.contextMenuNode)
        } else {
          ElMessage.error(res.data?.message || '创建失败')
        }
      } catch (e) {
        console.error('submitAddCollection error:', e)
        if (e !== false) {
          ElMessage.error('创建失败: ' + (e.message || ''))
        }
      }
    },
    async handleDropDatabase() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false

      try {
        const { instanceId, dbName } = this.contextMenuNode
        await ElMessageBox.confirm(`确定删除数据库 "${dbName}" 吗？此操作不可恢复！`, '删除确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        const res = await mongoAPI.dropDatabase(instanceId, {
          database: dbName
        })
        console.log('handleDropDatabase response:', res)

        if (res.data?.code === 0 || res.data?.code === 200) {
          ElMessage.success('数据库删除成功')
          const parentNode = this.treeData.find(inst => inst.instanceId === instanceId)
          if (parentNode) {
            await this.loadDatabases(parentNode)
          }
        } else {
          ElMessage.error(res.data?.message || '删除失败')
        }
      } catch (e) {
        console.error('handleDropDatabase error:', e)
        if (e !== 'cancel') {
          ElMessage.error('删除失败: ' + (e.message || ''))
        }
      }
    },
    async handleCollectionStats() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false

      try {
        const { instanceId, dbName, collectionName } = this.contextMenuNode

        const res = await mongoAPI.collectionStats(instanceId, {
          database: dbName,
          collection: collectionName
        })
        console.log('handleCollectionStats response:', res)

        let stats = null
        if (res.data?.code === 0 || res.data?.code === 200) {
          stats = res.data.data
        }

        this.statsContent = JSON.stringify(stats, null, 2)
        this.statsDialogVisible = true
      } catch (e) {
        console.error('handleCollectionStats error:', e)
        ElMessage.error('获取统计信息失败: ' + (e.message || ''))
      }
    },
    async handleDropCollection() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false

      try {
        const { instanceId, dbName, collectionName } = this.contextMenuNode
        await ElMessageBox.confirm(`确定删除集合 "${collectionName}" 吗？此操作不可恢复！`, '删除确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        const res = await mongoAPI.dropCollection(instanceId, {
          database: dbName,
          collection: collectionName
        })
        console.log('handleDropCollection response:', res)

        if (res.data?.code === 0 || res.data?.code === 200) {
          ElMessage.success('集合删除成功')
          const parentNode = this.treeData.find(inst => inst.instanceId === instanceId)
          if (parentNode) {
            const dbNode = parentNode.children.find(db => db.dbName === dbName)
            if (dbNode) {
              await this.loadCollections(dbNode)
            }
          }
        } else {
          ElMessage.error(res.data?.message || '删除失败')
        }
      } catch (e) {
        console.error('handleDropCollection error:', e)
        if (e !== 'cancel') {
          ElMessage.error('删除失败: ' + (e.message || ''))
        }
      }
    },
    handleCopyCollectionName() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false

      const collectionName = this.contextMenuNode.collectionName
      if (navigator.clipboard) {
        navigator.clipboard.writeText(collectionName).then(() => {
          ElMessage.success(`已复制: ${collectionName}`)
        }).catch(() => {
          ElMessage.error('复制失败')
        })
      } else {
        const textarea = document.createElement('textarea')
        textarea.value = collectionName
        textarea.style.position = 'fixed'
        textarea.style.opacity = '0'
        document.body.appendChild(textarea)
        textarea.select()
        try {
          document.execCommand('copy')
          ElMessage.success(`已复制: ${collectionName}`)
        } catch (e) {
          ElMessage.error('复制失败')
        }
        document.body.removeChild(textarea)
      }
    }
  }
}
</script>

<style scoped>
.mongo-terminal-page { padding: 20px; min-height: 100vh; background: #eef1f6; }
.terminal-container { background: transparent; }
.tree-aside { background: transparent; padding: 0 10px 0 0; position: relative; }
.tree-card { border-radius: 8px; height: calc(100vh - 40px); }
.tree-header { display: flex; justify-content: space-between; align-items: center; }
.node-icon { width: 18px; height: 18px; }
.tree-node { display: inline-flex; align-items: center; gap: 8px; }
.tree-card :deep(.el-tree--highlight-current .el-tree-node.is-current > .el-tree-node__content) {
  background-color: #eef5ff !important;
  color: #409eff;
  font-weight: 600;
  border-radius: 6px;
}
.tree-card :deep(.el-tree--highlight-current .el-tree-node.is-current > .el-tree-node__content .node-label) {
  color: #409eff;
}
.node-count { font-size: 12px; color: #909399; }
.aside-resizer { width: 6px; cursor: col-resize; background: transparent; display: flex; align-items: center; }
.resizer-icon { color: #909399; }
.main-content { padding-left: 12px; padding-right: 0; }
.main-card { border-radius: 8px; height: calc(100vh - 40px); overflow: hidden; display: flex; flex-direction: column; }
.main-card :deep(.el-card__body) { display: flex; flex-direction: column; height: 100%; }
.content-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.header-info { display: flex; align-items: center; gap: 10px; }
.header-title { font-size: 18px; font-weight: 600; color: #2c3e50; }
.header-connection-info { margin-left: 20px; font-size: 14px; color: #606266; display: flex; align-items: center; }
.info-value { color: #667eea; font-weight: 500; }
.mb10 { margin-bottom: 10px; }
.tab-label-text { display: inline-block; max-width: 220px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; vertical-align: middle; }
.query-section { margin-bottom: 15px; padding: 15px; background: rgba(103, 126, 234, 0.05); border-radius: 8px; }
.code-editor-wrapper { width: 100%; height: 500px; border: 1px solid #dcdfe6; border-radius: 4px; overflow: hidden; }

/* 右键菜单样式 */
.context-menu {
  position: fixed;
  z-index: 9999;
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 4px 0;
  min-width: 140px;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  cursor: pointer;
  font-size: 14px;
  color: #606266;
  transition: all 0.3s;
}

.context-menu-item:hover {
  background: #f5f7fa;
  color: #409eff;
}

.context-menu-item .el-icon {
  font-size: 16px;
}
</style>
