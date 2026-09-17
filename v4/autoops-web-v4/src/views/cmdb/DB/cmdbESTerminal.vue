<template>
  <div class="es-terminal-page">
    <el-container class="terminal-container">
      <!-- 左侧树形结构 -->
      <el-aside :width="asideWidth + 'px'" class="tree-aside">
        <el-card shadow="never" class="tree-card">
          <template #header>
            <div class="tree-header">
              <span>Elasticsearch 资源</span>
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
                    <img :src="esIcon" class="node-icon" />
                  </el-icon>
                  <el-icon v-else-if="data.type === 'index'"><img src="@/assets/image/数据库.svg" class="node-icon" /></el-icon>
                  <el-icon v-else-if="data.type === 'folder'"><img src="@/assets/image/打开文件夹.svg" class="node-icon" /></el-icon>
                  <el-icon v-else><DataAnalysis /></el-icon>
                  <span class="node-label">{{ node.label }}</span>
                  <span v-if="data.docCount !== undefined" class="node-count">({{ data.docCount }})</span>
                </span>
              </template>
            </el-tree>
          </el-scrollbar>

          <!-- 右键菜单 - 使用绝对定位 -->
          <div
            v-show="contextMenuVisible"
            ref="contextMenuRef"
            class="context-menu"
            :style="{ left: contextMenuX + 'px', top: contextMenuY + 'px' }"
            @click.stop
          >
            <!-- 实例右键菜单 -->
            <template v-if="contextMenuNode && contextMenuNode.type === 'instance'">
              <div class="context-menu-item" @click="handleRefreshIndices">
                <el-icon><Refresh /></el-icon>
                <span>刷新索引</span>
              </div>
              <div class="context-menu-item" @click="handleAddIndex">
                <el-icon><Plus /></el-icon>
                <span>添加索引</span>
              </div>
            </template>

            <!-- 索引右键菜单 -->
            <template v-if="contextMenuNode && contextMenuNode.type === 'index'">
              <div class="context-menu-item" @click="handleIndexDetails">
                <el-icon><View /></el-icon>
                <span>索引详情</span>
              </div>
              <div class="context-menu-item" @click="handleDeleteIndex">
                <el-icon><Delete /></el-icon>
                <span>删除索引</span>
              </div>
              <div class="context-menu-item" @click="handleCloseIndex">
                <el-icon><Lock /></el-icon>
                <span>关闭索引</span>
              </div>
              <div class="context-menu-item" @click="handleMigrateIndex">
                <el-icon><Promotion /></el-icon>
                <span>索引迁移</span>
              </div>
              <div class="context-menu-item" @click="handleCopyIndexName">
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
                    <span style="margin-right: 15px;"><strong>索引:</strong> <span class="info-value">{{ tab.connection.index }}</span></span>
                    <span style="margin-right: 15px;" v-if="tab.meta"><strong>文档数:</strong> <span class="info-value">{{ tab.meta.docCount }}</span></span>
                  </div>
                  <div class="header-actions">
                    <el-button type="success" size="small" @click="reloadIndexData(tab)" :loading="tab.loading">
                      <el-icon><Refresh /></el-icon>
                      刷新
                    </el-button>
                  </div>
                </div>

                <!-- 查询条件 -->
                <div class="query-section">
                  <el-form :inline="true" size="small">
                    <el-form-item label="查询条件">
                      <el-input v-model="tab.queryText" placeholder='{"match_all": {}}' style="width: 400px" clearable>
                        <template #prepend>Query DSL</template>
                      </el-input>
                    </el-form-item>
                    <el-form-item label="分页大小">
                      <el-input-number v-model="tab.pageSize" :min="1" :max="1000" style="width: 120px" />
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
                        <span v-if="col.prop === '_source'">{{ formatJSON(row[col.prop]) }}</span>
                        <span v-else>{{ row[col.prop] }}</span>
                      </template>
                    </el-table-column>
                    <el-table-column fixed="right" label="操作" width="180">
                      <template #default="{ row }">
                        <el-button size="small" type="primary" text @click="handleViewDoc(row, tab)">查看</el-button>
                        <el-button size="small" type="warning" text @click="handleEditDoc(row, tab)">编辑</el-button>
                        <el-button size="small" type="danger" text @click="handleDeleteDoc(row, tab)">删除</el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                  <el-empty v-else description="暂无数据" />

                  <!-- 分页 -->
                  <div v-if="tab.total > 0" class="pagination-section">
                    <el-pagination
                      v-model:current-page="tab.currentPage"
                      v-model:page-size="tab.pageSize"
                      :page-sizes="[10, 20, 50, 100]"
                      :total="tab.total"
                      layout="total, sizes, prev, pager, next"
                      @size-change="handlePageSizeChange(tab)"
                      @current-change="handlePageChange(tab)"
                    />
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
          <el-empty v-else description="请从左侧选择一个索引" />
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

    <!-- 添加索引对话框 -->
    <el-dialog v-model="addIndexDialogVisible" title="添加索引" width="40%" destroy-on-close>
      <el-form :model="indexForm" :rules="indexRules" ref="indexFormRef" label-width="100px">
        <el-form-item label="索引名称" prop="indexName">
          <el-input v-model="indexForm.indexName" placeholder="请输入索引名称（小写字母、数字、-、_）" />
        </el-form-item>
        <el-form-item label="分片数" prop="shards">
          <el-input-number v-model="indexForm.shards" :min="1" :max="10" />
        </el-form-item>
        <el-form-item label="副本数" prop="replicas">
          <el-input-number v-model="indexForm.replicas" :min="0" :max="5" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addIndexDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAddIndex">确定</el-button>
      </template>
    </el-dialog>

    <!-- 索引详情对话框 -->
    <el-dialog v-model="indexDetailsDialogVisible" title="索引详情" width="60%" destroy-on-close>
      <div class="code-editor-wrapper">
        <CodeEditor v-model="indexDetailsContent" language="json" :readonly="true" />
      </div>
      <template #footer>
        <el-button type="primary" @click="indexDetailsDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 索引迁移对话框 -->
    <el-dialog v-model="migrateIndexDialogVisible" title="索引迁移" width="50%" destroy-on-close>
      <el-form :model="migrateForm" :rules="migrateRules" ref="migrateFormRef" label-width="120px">
        <el-form-item label="源索引">
          <el-input v-model="migrateForm.sourceIndex" disabled />
        </el-form-item>
        <el-form-item label="目标索引" prop="targetIndex">
          <el-input v-model="migrateForm.targetIndex" placeholder="请输入目标索引名称" />
        </el-form-item>
        <el-form-item label="目标实例" prop="targetInstanceId">
          <el-select v-model="migrateForm.targetInstanceId" placeholder="选择目标实例" style="width: 100%">
            <el-option
              v-for="inst in treeData"
              :key="inst.instanceId"
              :label="inst.label"
              :value="inst.instanceId"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="迁移方式">
          <el-radio-group v-model="migrateForm.method">
            <el-radio value="reindex">Reindex（推荐）</el-radio>
            <el-radio value="snapshot">快照恢复</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="migrateIndexDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitMigrateIndex">开始迁移</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { Refresh, Search, DArrowRight, DataAnalysis, Plus, View, Delete, Lock, Promotion, CopyDocument } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { esAPI } from '@/api/database'
import CodeEditor from '@/components/CodeEditor.vue'

export default {
  components: { Refresh, Search, DArrowRight, DataAnalysis, CodeEditor, Plus, View, Delete, Lock, Promotion, CopyDocument },
  data() {
    return {
      esIcon: require('@/assets/image/Elasticsearch.svg'),
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
      // 右键菜单相关
      contextMenuVisible: false,
      contextMenuX: 0,
      contextMenuY: 0,
      contextMenuNode: null,
      // 添加索引对话框
      addIndexDialogVisible: false,
      indexForm: {
        indexName: '',
        shards: 1,
        replicas: 1
      },
      indexRules: {
        indexName: [
          { required: true, message: '请输入索引名称', trigger: 'blur' },
          { pattern: /^[a-z0-9_-]+$/, message: '只能包含小写字母、数字、-、_', trigger: 'blur' }
        ]
      },
      // 索引详情对话框
      indexDetailsDialogVisible: false,
      indexDetailsContent: '',
      // 索引迁移对话框
      migrateIndexDialogVisible: false,
      migrateForm: {
        sourceIndex: '',
        targetIndex: '',
        targetInstanceId: null,
        method: 'reindex'
      },
      migrateRules: {
        targetIndex: [
          { required: true, message: '请输入目标索引名称', trigger: 'blur' }
        ],
        targetInstanceId: [
          { required: true, message: '请选择目标实例', trigger: 'change' }
        ]
      }
    }
  },
  watch: {
    treeFilter(val) {
      this.$refs.treeRef && this.$refs.treeRef.filter(val)
    }
  },
  mounted() {
    // 监听点击事件，点击其他地方关闭右键菜单
    document.addEventListener('click', this.closeContextMenu)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.closeContextMenu)
  },
  created() {
    const idParam = this.$route?.params?.id || this.$route?.query?.id
    const idNum = idParam ? Number(idParam) : null
    if (idNum) {
      this.loadSingleInstance(idNum)
    } else {
      this.loadInstances()
    }
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
        const res = await esAPI.getInstanceList({ pageNum: 1, pageSize: 100 })
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
    async handleNodeClick(data) {
      this.selectedNode = data
      if (data.type === 'instance') {
        if (data.children.length === 0) await this.loadIndices(data)
      } else if (data.type === 'index') {
        this.openIndexTab(data)
      }
    },
    async loadSingleInstance(id) {
      this.loadingTree = true
      try {
        const res = await esAPI.getInstanceById(id)
        if (res.data?.code === 0 || res.data?.code === 200) {
          const inst = res.data.data
          if (!inst) throw new Error('实例不存在')
          const node = {
            id: `inst-${inst.id}`,
            label: inst.name,
            type: 'instance',
            instanceId: inst.id,
            children: [],
            isLeaf: false
          }
          this.treeData = [node]
          await this.loadIndices(node)
        } else {
          ElMessage.error(res.data?.message || '加载实例失败')
        }
      } catch (e) {
        ElMessage.error('加载实例失败')
      } finally {
        this.loadingTree = false
      }
    },
    async handleNodeExpand(data) {
      if (data.type === 'instance' && data.children.length === 0) {
        await this.loadIndices(data)
      }
    },
    
    async loadIndices(instNode) {
      try {
        const res = await esAPI.proxy(instNode.instanceId, '_cat/indices?format=json', 'get')
        console.log('loadIndices response:', res)

        // 后端返回的数据可能在 res.data.data 或直接在 res.data 中
        let indices = []
        if (res.data?.code === 0 || res.data?.code === 200) {
          indices = res.data.data || []
        } else if (Array.isArray(res.data)) {
          // 直接返回数组的情况
          indices = res.data
        } else {
          ElMessage.error('加载索引失败: ' + (res.data?.message || '未知错误'))
          return
        }

        console.log('parsed indices:', indices)

        instNode.children = indices.map(idx => ({
          id: `idx-${instNode.instanceId}-${idx.index}`,
          label: idx.index,
          type: 'index',
          instanceId: instNode.instanceId,
          indexName: idx.index,
          docCount: idx['docs.count'] || idx.docsCount || 0,
          isLeaf: true
        }))
      } catch (e) {
        console.error('loadIndices error:', e)
        ElMessage.error('加载索引失败: ' + (e.message || ''))
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
    async openIndexTab(node) {
      const tabName = `idx-${node.instanceId}-${node.indexName}`
      const exist = this.openTabs.find(t => t.name === tabName)
      if (exist) {
        this.activeTab = tabName
        if (!exist.rows || exist.rows.length === 0) {
          await this.reloadIndexData(exist)
        }
        return
      }

      let instanceName = this.treeData.find(i => i.instanceId === node.instanceId)?.label || `#${node.instanceId}`
      const tab = {
        name: tabName,
        title: node.indexName,
        node,
        connection: { instanceName, index: node.indexName },
        columns: [],
        rows: [],
        meta: { docCount: node.docCount },
        loading: false,
        queryText: '',
        pageSize: 20,
        currentPage: 1,
        total: 0
      }
      this.openTabs.push(tab)
      this.activeTab = tabName
      const reactiveTab = this.openTabs[this.openTabs.length - 1]
      await this.reloadIndexData(reactiveTab)
    },
    async reloadIndexData(tab) {
      await this.handleQuery(tab)
    },
    async handleQuery(tab) {
      const { instanceId, indexName } = tab.node
      tab.loading = true
      try {
        let query = {}
        if (tab.queryText) {
          try {
            query = JSON.parse(tab.queryText)
          } catch {
            query = { query_string: { query: tab.queryText } }
          }
        } else {
          query = { match_all: {} }
        }

        const from = (tab.currentPage - 1) * tab.pageSize
        const body = {
          query,
          from,
          size: tab.pageSize
        }

        const res = await esAPI.proxy(instanceId, `${indexName}/_search`, 'post', body)
        console.log('handleQuery response:', res)

        // 处理返回数据：可能直接在 res.data 或 res.data.data 中
        let data = null
        if (res.data?.code === 0 || res.data?.code === 200) {
          data = res.data.data || {}
        } else if (res.data?.hits) {
          // ES 直接返回的数据格式
          data = res.data
        } else {
          ElMessage.error(res.data?.message || '查询失败')
          tab.rows = []
          tab.total = 0
          return
        }

        console.log('parsed data:', data)

        const hits = data.hits?.hits || []
        tab.total = typeof data.hits?.total === 'number' ? data.hits.total : (data.hits?.total?.value || 0)

        // 构建列
        const allFields = new Set(['_id', '_index', '_score'])
        hits.forEach(hit => {
          if (hit._source) {
            Object.keys(hit._source).forEach(k => allFields.add(k))
          }
        })

        tab.columns = Array.from(allFields).map(field => ({
          prop: field,
          label: field,
          minWidth: field === '_id' ? 180 : field === '_source' ? 300 : 150
        }))

        // 构建行数据
        tab.rows = hits.map((hit, idx) => {
          const row = {
            _id: hit._id,
            _index: hit._index,
            _score: hit._score,
            _rawDoc: hit
          }
          if (hit._source) {
            Object.keys(hit._source).forEach(k => {
              row[k] = hit._source[k]
            })
          }
          return row
        })
      } catch (e) {
        console.error('handleQuery error:', e)
        ElMessage.error('查询失败: ' + (e.message || ''))
        tab.rows = []
        tab.total = 0
      } finally {
        tab.loading = false
      }
    },
    resetQuery(tab) {
      tab.queryText = ''
      tab.currentPage = 1
      tab.pageSize = 20
      this.handleQuery(tab)
    },
    handlePageSizeChange(tab) {
      tab.currentPage = 1
      this.handleQuery(tab)
    },
    handlePageChange(tab) {
      this.handleQuery(tab)
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
      this.docContent = JSON.stringify(row._rawDoc._source || {}, null, 2)
      this.docReadonly = false
      this.currentDoc = row
      this.currentTab = tab
      this.docDialogVisible = true
    },
    async saveDoc() {
      try {
        const newSource = JSON.parse(this.docContent)
        const { instanceId, indexName } = this.currentTab.node
        const docId = this.currentDoc._id

        const res = await esAPI.proxy(instanceId, `${indexName}/_doc/${docId}`, 'put', newSource)
        console.log('saveDoc response:', res)

        // 处理返回数据：ES 更新成功返回 result: "updated"
        let success = false
        if (res.data?.code === 0 || res.data?.code === 200) {
          success = true
        } else if (res.data?.result === 'updated' || res.data?.result === 'created') {
          // ES 直接返回更新结果
          success = true
        } else if (res.data?._shards?.successful > 0) {
          // 通过 shards 判断成功
          success = true
        }

        if (success) {
          ElMessage.success('保存成功')
          this.docDialogVisible = false
          await this.reloadIndexData(this.currentTab)
        } else {
          ElMessage.error(res.data?.message || '保存失败')
        }
      } catch (e) {
        console.error('saveDoc error:', e)
        ElMessage.error('保存失败: ' + (e.message || ''))
      }
    },
    async handleDeleteDoc(row, tab) {
      try {
        await ElMessageBox.confirm(`确定删除文档 ${row._id} 吗？`, '删除确认', { type: 'warning' })
        const { instanceId, indexName } = tab.node
        const res = await esAPI.proxy(instanceId, `${indexName}/_doc/${row._id}`, 'delete')
        console.log('handleDeleteDoc response:', res)

        // 处理返回数据：ES 删除成功返回 result: "deleted"
        let success = false
        if (res.data?.code === 0 || res.data?.code === 200) {
          success = true
        } else if (res.data?.result === 'deleted') {
          // ES 直接返回删除结果
          success = true
        } else if (res.data?._shards?.successful > 0) {
          // 通过 shards 判断成功
          success = true
        }

        if (success) {
          ElMessage.success('删除成功')
          await this.reloadIndexData(tab)
        } else {
          ElMessage.error(res.data?.message || '删除失败')
        }
      } catch (e) {
        console.error('handleDeleteDoc error:', e)
        if (e !== 'cancel') ElMessage.error('删除失败: ' + (e.message || ''))
      }
    },
    formatJSON(obj) {
      if (typeof obj === 'object') {
        return JSON.stringify(obj)
      }
      return String(obj)
    },
    // 右键菜单相关方法
    handleNodeRightClick(event, data, node, element) {
      // 只对实例和索引节点显示右键菜单
      if (data.type !== 'instance' && data.type !== 'index') return

      event.preventDefault()
      event.stopPropagation()

      this.contextMenuNode = data

      // 获取鼠标相对于视口的位置
      this.contextMenuX = event.clientX
      this.contextMenuY = event.clientY
      this.contextMenuVisible = true
    },
    closeContextMenu() {
      this.contextMenuVisible = false
    },
    async handleRefreshIndices() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false
      await this.loadIndices(this.contextMenuNode)
      ElMessage.success('刷新成功')
    },
    handleAddIndex() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false
      this.indexForm = {
        indexName: '',
        shards: 1,
        replicas: 1
      }
      this.addIndexDialogVisible = true
    },
    async submitAddIndex() {
      try {
        await this.$refs.indexFormRef.validate()

        const { instanceId } = this.contextMenuNode
        const { indexName, shards, replicas } = this.indexForm

        const body = {
          settings: {
            number_of_shards: shards,
            number_of_replicas: replicas
          }
        }

        const res = await esAPI.proxy(instanceId, indexName, 'put', body)
        console.log('submitAddIndex response:', res)

        // 处理返回数据
        let success = false
        if (res.data?.code === 0 || res.data?.code === 200) {
          success = true
        } else if (res.data?.acknowledged === true) {
          // ES 直接返回创建结果
          success = true
        }

        if (success) {
          ElMessage.success('索引创建成功')
          this.addIndexDialogVisible = false
          // 刷新索引列表
          await this.loadIndices(this.contextMenuNode)
        } else {
          ElMessage.error(res.data?.message || res.data?.error?.reason || '创建失败')
        }
      } catch (e) {
        console.error('submitAddIndex error:', e)
        if (e !== false) {
          ElMessage.error('创建失败: ' + (e.message || ''))
        }
      }
    },
    // 索引右键菜单方法
    async handleIndexDetails() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false

      try {
        const { instanceId, indexName } = this.contextMenuNode
        const res = await esAPI.proxy(instanceId, `${indexName}`, 'get')
        console.log('handleIndexDetails response:', res)

        let data = null
        if (res.data?.code === 0 || res.data?.code === 200) {
          data = res.data.data
        } else if (res.data && typeof res.data === 'object') {
          data = res.data
        }

        this.indexDetailsContent = JSON.stringify(data, null, 2)
        this.indexDetailsDialogVisible = true
      } catch (e) {
        console.error('handleIndexDetails error:', e)
        ElMessage.error('获取索引详情失败: ' + (e.message || ''))
      }
    },
    async handleDeleteIndex() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false

      try {
        const { instanceId, indexName } = this.contextMenuNode
        await ElMessageBox.confirm(`确定删除索引 "${indexName}" 吗？此操作不可恢复！`, '删除确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        const res = await esAPI.proxy(instanceId, indexName, 'delete')
        console.log('handleDeleteIndex response:', res)

        let success = false
        if (res.data?.code === 0 || res.data?.code === 200) {
          success = true
        } else if (res.data?.acknowledged === true) {
          success = true
        }

        if (success) {
          ElMessage.success('索引删除成功')
          // 刷新父节点的索引列表
          const parentNode = this.treeData.find(inst => inst.instanceId === instanceId)
          if (parentNode) {
            await this.loadIndices(parentNode)
          }
        } else {
          ElMessage.error(res.data?.message || res.data?.error?.reason || '删除失败')
        }
      } catch (e) {
        console.error('handleDeleteIndex error:', e)
        if (e !== 'cancel') {
          ElMessage.error('删除失败: ' + (e.message || ''))
        }
      }
    },
    async handleCloseIndex() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false

      try {
        const { instanceId, indexName } = this.contextMenuNode
        await ElMessageBox.confirm(`确定关闭索引 "${indexName}" 吗？`, '关闭确认', {
          type: 'warning'
        })

        const res = await esAPI.proxy(instanceId, `${indexName}/_close`, 'post')
        console.log('handleCloseIndex response:', res)

        let success = false
        if (res.data?.code === 0 || res.data?.code === 200) {
          success = true
        } else if (res.data?.acknowledged === true) {
          success = true
        }

        if (success) {
          ElMessage.success('索引关闭成功')
        } else {
          ElMessage.error(res.data?.message || res.data?.error?.reason || '关闭失败')
        }
      } catch (e) {
        console.error('handleCloseIndex error:', e)
        if (e !== 'cancel') {
          ElMessage.error('关闭失败: ' + (e.message || ''))
        }
      }
    },
    handleMigrateIndex() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false

      this.migrateForm = {
        sourceIndex: this.contextMenuNode.indexName,
        targetIndex: this.contextMenuNode.indexName + '_copy',
        targetInstanceId: this.contextMenuNode.instanceId,
        method: 'reindex'
      }
      this.migrateIndexDialogVisible = true
    },
    async submitMigrateIndex() {
      try {
        await this.$refs.migrateFormRef.validate()

        const { sourceIndex, targetIndex, targetInstanceId, method } = this.migrateForm
        const sourceInstanceId = this.contextMenuNode.instanceId

        if (method === 'reindex') {
          // 使用 Reindex API
          const body = {
            source: {
              index: sourceIndex
            },
            dest: {
              index: targetIndex
            }
          }

          const res = await esAPI.proxy(targetInstanceId, '_reindex', 'post', body)
          console.log('submitMigrateIndex response:', res)

          let success = false
          if (res.data?.code === 0 || res.data?.code === 200) {
            success = true
          } else if (res.data?.failures?.length === 0 || res.data?.total > 0) {
            success = true
          }

          if (success) {
            ElMessage.success('索引迁移成功')
            this.migrateIndexDialogVisible = false
            // 刷新目标实例的索引列表
            const targetNode = this.treeData.find(inst => inst.instanceId === targetInstanceId)
            if (targetNode) {
              await this.loadIndices(targetNode)
            }
          } else {
            ElMessage.error('迁移失败')
          }
        } else {
          ElMessage.warning('快照恢复功能暂未实现')
        }
      } catch (e) {
        console.error('submitMigrateIndex error:', e)
        if (e !== false) {
          ElMessage.error('迁移失败: ' + (e.message || ''))
        }
      }
    },
    handleCopyIndexName() {
      if (!this.contextMenuNode) return
      this.contextMenuVisible = false

      const indexName = this.contextMenuNode.indexName
      // 使用 Clipboard API 复制到剪贴板
      if (navigator.clipboard) {
        navigator.clipboard.writeText(indexName).then(() => {
          ElMessage.success(`已复制: ${indexName}`)
        }).catch(() => {
          ElMessage.error('复制失败')
        })
      } else {
        // 兼容旧浏览器
        const textarea = document.createElement('textarea')
        textarea.value = indexName
        textarea.style.position = 'fixed'
        textarea.style.opacity = '0'
        document.body.appendChild(textarea)
        textarea.select()
        try {
          document.execCommand('copy')
          ElMessage.success(`已复制: ${indexName}`)
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
.es-terminal-page { padding: 20px; min-height: 100vh; background: #eef1f6; }
.terminal-container { background: transparent; }
.tree-aside { background: transparent; padding: 0 10px 0 0; }
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
.pagination-section { margin-top: 15px; display: flex; justify-content: flex-end; }
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

.tree-card {
  position: relative;
}
</style>
