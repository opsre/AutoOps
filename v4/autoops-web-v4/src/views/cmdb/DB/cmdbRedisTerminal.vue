<template>
  <div class="db-terminal-page">
    <el-container class="terminal-container">
      <!-- 左侧树形结构 -->
      <el-aside :width="asideWidth + 'px'" class="tree-aside">
        <el-card shadow="never" class="tree-card">
          <template #header>
            <div class="tree-header">
              <span>Redis 资源</span>
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
            >
              <template #default="{ node, data }">
                <span class="tree-node">
                  <el-icon v-if="data.type === 'instance'">
                    <img :src="redisIcon" class="node-icon" />
                  </el-icon>
                  <el-icon v-else-if="data.type === 'database'"><img src="@/assets/image/数据库.svg" class="node-icon" /></el-icon>
                  <el-icon v-else-if="data.type === 'folder'"><img src="@/assets/image/打开文件夹.svg" class="node-icon" /></el-icon>
                  <el-icon v-else-if="data.type === 'table'"><img src="@/assets/image/数据表.svg" class="node-icon" /></el-icon>
                  <el-icon v-else-if="data.type === 'loadmore'"><img src="@/assets/image/打开文件夹.svg" class="node-icon" /></el-icon>
                  <span class="node-label">{{ node.label }}</span>
                </span>
              </template>
            </el-tree>
          </el-scrollbar>
        </el-card>
      </el-aside>
      <!-- Draggable resizer between aside and main -->
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
                    <span style="margin-right: 15px;"><strong>DB:</strong> <span class="info-value">{{ tab.connection.db }}</span></span>
                    <span style="margin-right: 15px;" v-if="tab.meta"><strong>类型:</strong> <span class="info-value">{{ tab.meta.type }}</span></span>
                    <span style="margin-right: 15px;" v-if="tab.meta"><strong>TTL:</strong> <span class="info-value">{{ tab.meta.ttl }}</span></span>
                    <span style="margin-right: 15px;" v-if="tab.meta"><strong>内存:</strong> <span class="info-value">{{ tab.meta.memuse }}</span></span>
                  </div>
                  <div class="header-actions">
                    <el-button type="success" circle @click="reloadKeyData(tab)" :loading="tab.loading" title="刷新">
                      <el-icon><Refresh /></el-icon>
                    </el-button>
                  </div>
                </div>
                <div>
                  <!-- 贴近 MySQL 样式：在表格上方展示可选 JSON 区域与工具按钮 -->
                  <div class="tools-row" v-if="(tab.meta?.type || '').toLowerCase() === 'string'">
                    <el-switch v-if="tab.isJson" v-model="tab.showJson" inline-prompt active-text="JSON" inactive-text="文本" />
                  </div>
                  <div v-if="(tab.meta?.type || '').toLowerCase() === 'string' && tab.isJson && tab.showJson" class="mb10">
                    <pre style="white-space: pre-wrap; background: #0b0e14; color:#eaeaea; padding: 12px; border-radius: 6px;">
                      <code v-html="tab.jsonHtml"></code>
                    </pre>
                  </div>
                  <el-table v-if="tab.rows && tab.rows.length" :data="tab.rows" border stripe max-height="520">
                    <el-table-column v-for="col in tab.columns" :key="col.prop" :prop="col.prop" :label="col.label" :min-width="col.minWidth || 120" />
                    <el-table-column fixed="right" label="操作" width="180">
                      <template #default="{ row }">
                        <el-button size="small" type="warning" text @click="handleEditRow(row, tab)">修改</el-button>
                        <el-button size="small" type="danger" text @click="handleDeleteRow(row, tab)">删除</el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                  <el-empty v-else description="暂无数据" />
               </div>
              </el-tab-pane>
            </el-tabs>
          </div>
          <el-empty v-else description="请从左侧选择一个 key" />
        </el-card>
      </el-main>
    </el-container>

    <!-- 无弹框，所有信息在右侧标签内展示 -->
  </div>
</template>

<script>
import { Refresh, Search, DArrowRight, View } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { redisAPI } from '@/api/database'
import { highlight } from '@/utils/highlight'

export default {
  components: { Refresh, Search, DArrowRight, View },
  data() {
    return {
      redisIcon: require('@/assets/image/redis.svg'),
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
      currentInstance: null,
      singleInstanceId: null,
      keySeparator: ':'
    }
  },
  watch: {
    treeFilter(val) {
      this.$refs.treeRef && this.$refs.treeRef.filter(val)
    }
  },
  created() {
    const idParam = this.$route?.params?.id || this.$route?.query?.id
    const idNum = idParam ? Number(idParam) : null
    if (idNum) {
      this.singleInstanceId = idNum
      this.loadSingleInstance(idNum)
    } else {
      this.loadInstances()
    }
  },
  methods: {
    async refreshTree() {
      if (this.singleInstanceId) {
        await this.loadSingleInstance(this.singleInstanceId)
      } else {
        await this.loadInstances()
      }
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
        const res = await redisAPI.getInstanceList({ pageNum: 1, pageSize: 100 })
        if (res.data?.code === 0 || res.data?.code === 200) {
          let instances = res.data?.data
          if (!Array.isArray(instances)) instances = instances?.list || instances?.data || []
          this.treeData = instances.map(inst => ({
            id: `inst-${inst.id}`,
            label: inst.name,
            type: 'instance',
            instanceId: inst.id,
            mode: inst.mode,
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
    getShortKeyLabel(key) {
      if (!key) return ''
      let label = String(key)
      label = label.includes(':') ? label.split(':').pop() : label
      if (label.length > 24) label = label.slice(0, 12) + '…' + label.slice(-8)
      return label
    },
    getTabShortLabel(tab) {
      const key = (tab && tab.title) ? String(tab.title) : ''
      // 优先显示最后一个冒号分段（更有业务意义）
      let label = key.includes(':') ? key.split(':').pop() : key
      // 再做长度控制
      if (label.length > 24) {
        label = label.slice(0, 12) + '…' + label.slice(-8)
      }
      return label || '未命名'
    },
    async loadSingleInstance(id) {
      this.loadingTree = true
      try {
        const res = await redisAPI.getInstanceById(id)
        if (res.data?.code === 0 || res.data?.code === 200) {
          const inst = res.data.data
          if (!inst) throw new Error('实例不存在')
          const node = {
            id: `inst-${inst.id}`,
            label: inst.name,
            type: 'instance',
            instanceId: inst.id,
            mode: inst.mode,
            children: [],
            isLeaf: false
          }
          this.treeData = [node]
          this.currentInstance = node
          await this.loadDatabases(node)
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
      if (data.type === 'instance') {
        const hasDbChildren = (data.children || []).some(c => c.type === 'database')
        if (!hasDbChildren) await this.loadDatabases(data)
      } else if (data.type === 'database' && data.children.length === 0) {
        await this.loadKeys(data)
      } else if (data.type === 'folder' && data.children.length === 0) {
        await this.loadFolderChildren(data)
      }
    },
    async handleNodeClick(data) {
      this.selectedNode = data
      if (data.type === 'instance') {
        this.currentInstance = data
        const hasDbChildren = (data.children || []).some(c => c.type === 'database')
        if (!hasDbChildren) await this.loadDatabases(data)
      } else if (data.type === 'database') {
        if (data.children.length === 0) await this.loadKeys(data)
      } else if (data.type === 'folder') {
        if (data.children.length === 0) await this.loadFolderChildren(data)
      } else if (data.type === 'table') {
        this.openKeyTab(data)
      } else if (data.type === 'loadmore') {
        await this.loadKeys(data.parent)
      }
    },
    async loadDatabases(instNode) {
      // 尝试新接口 /redis/:id/dbs 获取可用 DB 列表
      let dbs = []
      let dbNames = []
      try {
        const res = await redisAPI.getDBs(instNode.instanceId)
        if (res.data?.code === 0 || res.data?.code === 200) {
          const d = res.data?.data || {}
          dbs = Array.isArray(d.dbs) ? d.dbs : []
          dbNames = Array.isArray(d.dbNames) ? d.dbNames : []
        }
      } catch {}
      // 兼容兜底：INFO Keyspace 或 0..15
      if (!dbs.length) {
        try {
          const res = await redisAPI.getInfo(instNode.instanceId, { format: 'map' })
          if (res.data?.code === 0 || res.data?.code === 200) {
            const info = res.data.data || {}
            const keyspace = info.Keyspace || info.keyspace || {}
            const dbKeys = Object.keys(keyspace)
            if (dbKeys.length) {
              dbs = dbKeys.map(k => Number(k.replace('db', ''))).sort((a,b)=>a-b)
            }
          }
        } catch {}
      }
      if (!dbs.length) {
        dbs = Array.from({ length: 16 }, (_, i) => i)
        dbNames = dbs.map(i => `db${i}`)
      } else if (!dbNames.length) {
        dbNames = dbs.map(i => `db${i}`)
      }

      // 保留已存在的“全部DB”节点（仅非 cluster）
      const existingAll = (instNode.children || []).filter(c => c.type === 'all')
      const dbNodes = dbs.map((dbNo, idx) => ({
        id: `db-${instNode.instanceId}-${dbNo}`,
        label: dbNames[idx] || `db${dbNo}`,
        type: 'database',
        instanceId: instNode.instanceId,
        mode: instNode.mode,
        db: dbNo,
        children: [],
        isLeaf: false,
        cursorSingle: 0,
        cursor: null
      }))
      instNode.children = [...existingAll, ...dbNodes]
    },
    async loadKeys(dbNode) {
      const body = { match: '*', count: 200 }
      if (dbNode.mode === 'cluster') {
        body.cursor = dbNode.cursor || {}
      } else {
        body.cursorSingle = dbNode.cursorSingle || 0
      }
      try {
        const res = await redisAPI.scanKeys(dbNode.instanceId, dbNode.db, body)
        if (res.data?.code === 0 || res.data?.code === 200) {
          const data = res.data.data || {}
          const keys = data.keys || []
          // Merge loaded keys for grouping
          if (!Array.isArray(dbNode.loadedKeys)) dbNode.loadedKeys = []
          for (const k of keys) if (!dbNode.loadedKeys.includes(k)) dbNode.loadedKeys.push(k)
          // handle cursor
          if (dbNode.mode === 'cluster') {
            dbNode.cursor = data.cursor || {}
          } else {
            dbNode.cursorSingle = typeof data.cursorSingle === 'number' ? data.cursorSingle : 0
          }
          // add loadmore node when cursor indicates more
          const more = (dbNode.mode === 'cluster') ? Object.values(dbNode.cursor || {}).some(v => Number(v) !== 0) : (dbNode.cursorSingle !== 0)
          let children = this.buildGroupedChildren(dbNode, '')
          if (more) children = children.concat([{ id: `more-${dbNode.instanceId}-${dbNode.db}-${children.length}`, label: '加载更多...', type: 'loadmore', parent: dbNode, isLeaf: true }])
          dbNode.children = children
        } else {
          ElMessage.error(res.data?.message || '加载 keys 失败')
        }
      } catch (e) {
        ElMessage.error('加载 keys 失败')
      }
    },
    buildGroupedChildren(dbNode, prefix) {
      const sep = this.keySeparator || ':'
      const keys = Array.isArray(dbNode.loadedKeys) ? dbNode.loadedKeys.filter(k => (prefix ? (k === prefix || k.startsWith(prefix + sep)) : true)) : []
      const groups = new Map()
      const leaves = []
      for (const full of keys) {
        let rest
        if (!prefix) rest = full
        else if (full === prefix) rest = ''
        else if (full.startsWith(prefix + sep)) rest = full.slice(prefix.length + 1)
        else continue

        if (rest === '') {
          const tail = prefix.includes(sep) ? prefix.split(sep).pop() : prefix
          leaves.push(this.makeLeafNode(dbNode, full, tail))
          continue
        }
        const idx = rest.indexOf(sep)
        if (idx >= 0) {
          const seg = rest.slice(0, idx)
          if (!groups.has(seg)) groups.set(seg, true)
        } else {
          const tail = rest
          leaves.push(this.makeLeafNode(dbNode, full, tail))
        }
      }
      const folders = Array.from(groups.keys()).map(seg => ({
        id: `folder-${dbNode.instanceId}-${dbNode.db}-${prefix}-${seg}`,
        label: seg,
        type: 'folder',
        prefix: prefix ? `${prefix}${sep}${seg}` : seg,
        db: dbNode.db,
        instanceId: dbNode.instanceId,
        parent: dbNode,
        children: [],
        isLeaf: false
      }))
      folders.sort((a,b)=>a.label.localeCompare(b.label))
      leaves.sort((a,b)=>a.label.localeCompare(b.label))
      return folders.concat(leaves)
    },
    makeLeafNode(dbNode, fullKey, tailLabel) {
      return {
        id: `key-${dbNode.instanceId}-${dbNode.db}-${fullKey}`,
        label: tailLabel,
        type: 'table',
        key: fullKey,
        db: dbNode.db,
        instanceId: dbNode.instanceId,
        parent: dbNode,
        isLeaf: true
      }
    },
    async loadFolderChildren(folderNode) {
      const dbNode = folderNode.parent?.type === 'database' ? folderNode.parent : (folderNode.parent?.parent || null)
      if (!dbNode) return
      folderNode.children = this.buildGroupedChildren(dbNode, folderNode.prefix)
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
    async openKeyTab(node) {
      const tabName = `key-${node.instanceId}-${node.db}-${node.key}`
      const exist = this.openTabs.find(t => t.name === tabName)
      if (exist) {
        this.activeTab = tabName
        // 若数据尚未加载，立即触发加载
        if (!exist.rows || exist.rows.length === 0) {
          await this.reloadKeyData(exist)
        }
        return
      }

      // resolve instance name
      let instanceName = this.treeData.find(i => i.instanceId === node.instanceId)?.label || `#${node.instanceId}`
      const tab = {
        name: tabName,
        title: node.key,
        node,
        connection: { instanceName, db: node.db, key: node.key },
        columns: [],
        rows: [],
        meta: null,
        loading: false,
        isJson: false,
        showJson: true,
        jsonHtml: '',
        jsonText: '',
        rawText: ''
      }
      this.openTabs.push(tab)
      this.activeTab = tabName
      // 获取响应式代理后再加载，避免界面不更新
      const reactiveTab = this.openTabs[this.openTabs.length - 1]
      await this.reloadKeyData(reactiveTab)
    },
    async reloadKeyData(tab) {
      const { instanceId, db, key } = tab.node
      tab.loading = true
      try {
        const [info, ttl, mem] = await Promise.all([
          redisAPI.getKeyInfo(instanceId, db, key),
          redisAPI.getKeyTTL(instanceId, db, key),
          redisAPI.getKeyMemUse(instanceId, db, key)
        ])
        const type = info.data?.data?.type || 'string'
        const ttlv = ttl.data?.data
        const memv = mem.data?.data
        tab.meta = { type, ttl: ttlv, memuse: memv }

        // fetch value by type, normalize to id/field/value(/score)
        let rows = []
        const extractResult = (resp) => {
          // 后端可能返回在 data 或 data.result 下
          const d = resp?.data?.data
          if (d && Object.prototype.hasOwnProperty.call(d, 'result')) return d.result
          return d
        }
        // 处理 binary 类型数据：将 base64 解码为可读字符串
        const decodeBinary = (val) => {
          if (!val || typeof val !== 'object') return val
          if (val._type === 'binary' && val._base64) {
            try {
              // 尝试将 base64 解码为 UTF-8 字符串
              const decoded = atob(val._base64)
              // 检查是否为可打印字符，如果不是则返回十六进制表示
              if (/^[\x20-\x7E\s]*$/.test(decoded)) {
                return decoded
              } else {
                // 返回十六进制表示
                return Array.from(decoded).map(c => c.charCodeAt(0).toString(16).padStart(2, '0')).join(' ')
              }
            } catch (e) {
              return `[Binary: ${val._base64.substring(0, 20)}...]`
            }
          }
          return val
        }
        if (type === 'string') {
          const res = await redisAPI.runCmd(instanceId, db, { command: 'GET', args: [key] })
          let v = extractResult(res)
          if (Array.isArray(v)) v = v[0]
          v = decodeBinary(v)
          rows = [ { id: 1, value: v ?? '' } ]
        } else if (type === 'hash') {
          const res = await redisAPI.runCmd(instanceId, db, { command: 'HGETALL', args: [key] })
          const data = extractResult(res) || []
          const kv = Array.isArray(data)
            ? Array.from({ length: Math.ceil(data.length/2) }, (_,i)=>({ field: decodeBinary(data[i*2]), value: decodeBinary(data[i*2+1]) }))
            : Object.entries(data).map(([k,v])=>({ field: decodeBinary(k), value: decodeBinary(v) }))
          rows = kv.map((r,idx)=>({ id: idx+1, field: r.field, value: r.value }))
        } else if (type === 'list') {
          const res = await redisAPI.runCmd(instanceId, db, { command: 'LRANGE', args: [key, '0', '-1'] })
          const arr = extractResult(res) || []
          rows = (Array.isArray(arr) ? arr : []).map((v,i)=>({ id: i+1, value: decodeBinary(v) }))
        } else if (type === 'set') {
          const res = await redisAPI.runCmd(instanceId, db, { command: 'SMEMBERS', args: [key] })
          const arr = extractResult(res) || []
          rows = (Array.isArray(arr) ? arr : []).map((v,i)=>({ id: i+1, value: decodeBinary(v) }))
        } else if (type === 'zset') {
          const res = await redisAPI.runCmd(instanceId, db, { command: 'ZRANGE', args: [key, '0', '-1', 'WITHSCORES'] })
          const arr = extractResult(res) || []
          const rowsData = []
          if (Array.isArray(arr)) {
            for (let i=0;i<arr.length;i+=2) rowsData.push({ value: decodeBinary(arr[i]), score: Number(arr[i+1]) })
          }
          rows = rowsData.map((r,i)=>({ id: i+1, value: r.value, score: r.score }))
        } else {
          rows = [ { id: 1, value: `暂不支持的类型: ${type}` } ]
        }

        // build columns in required order
        const idLabel = `ID (Total: ${rows.length})`
        const hasField = rows.some(r => Object.prototype.hasOwnProperty.call(r, 'field'))
        const hasScore = rows.some(r => Object.prototype.hasOwnProperty.call(r, 'score'))
        const columns = [ { prop: 'id', label: idLabel, minWidth: 140 } ]
        if (hasField) columns.push({ prop: 'field', label: 'field', minWidth: 160 })
        columns.push({ prop: 'value', label: 'value', minWidth: 240 })
        if (hasScore) columns.push({ prop: 'score', label: 'score', minWidth: 120 })
        tab.columns = columns
        tab.rows = rows

        // Build JSON helpers for string type and general raw value
        const typeLower = (type || '').toLowerCase()
        let raw
        if (typeLower === 'string') {
          raw = (rows?.[0]?.value ?? '')
        } else if (typeLower === 'hash') {
          raw = {}
          for (const r of rows || []) raw[r.field] = r.value
        } else if (typeLower === 'list') {
          raw = (rows || []).map(r => r.value)
        } else if (typeLower === 'set') {
          raw = (rows || []).map(r => r.value)
        } else if (typeLower === 'zset') {
          raw = (rows || []).map(r => ({ value: r.value, score: r.score }))
        } else {
          raw = rows
        }
        let isJson = false
        let jsonText = ''
        try {
          if (typeof raw === 'string') {
            const s = raw.trim()
            if (s.startsWith('{') || s.startsWith('[')) {
              jsonText = JSON.stringify(JSON.parse(s), null, 2)
              isJson = true
            }
          } else {
            jsonText = JSON.stringify(raw, null, 2)
            isJson = true
          }
        } catch (_) { isJson = false }
        tab.isJson = isJson
        tab.jsonText = jsonText
        tab.rawText = typeof raw === 'string' ? raw : (isJson ? jsonText : String(raw))
        tab.jsonHtml = isJson ? highlight(jsonText, 'json') : ''
        if (!tab.viewTab) tab.viewTab = 'data'
      } catch (e) {
        ElMessage.error('加载Key数据失败')
        tab.columns = [ { prop: 'error', label: 'Error' } ]
        tab.rows = [ { error: e.message || '加载失败' } ]
      } finally {
        tab.loading = false
      }
    },
    
    async handleEditRow(row, tab) {
      const type = (tab.meta?.type || '').toLowerCase()
      const key = tab.connection?.key
      const db = tab.node.db
      const inst = tab.node.instanceId
      try {
        if (type === 'string') {
          const { value } = await ElMessageBox.prompt('请输入新值', '修改', { inputValue: String(row.value ?? ''), inputType: 'textarea' })
          const res = await redisAPI.runCmd(inst, db, { command: 'SET', args: [key, value] })
          if (res.data?.code === 0 || res.data?.code === 200) ElMessage.success('修改成功'); else throw new Error(res.data?.message || '修改失败')
        } else if (type === 'hash') {
          const { value } = await ElMessageBox.prompt(`字段 ${row.field} 新值`, '修改', { inputValue: String(row.value ?? ''), inputType: 'textarea' })
          const res = await redisAPI.runCmd(inst, db, { command: 'HSET', args: [key, row.field, value] })
          if (res.data?.code === 0 || res.data?.code === 200) ElMessage.success('修改成功'); else throw new Error(res.data?.message || '修改失败')
        } else if (type === 'list') {
          const { value } = await ElMessageBox.prompt(`索引 ${row.id - 1} 新值`, '修改', { inputValue: String(row.value ?? ''), inputType: 'textarea' })
          const res = await redisAPI.runCmd(inst, db, { command: 'LSET', args: [key, String(row.id - 1), value] })
          if (res.data?.code === 0 || res.data?.code === 200) ElMessage.success('修改成功'); else throw new Error(res.data?.message || '修改失败')
        } else if (type === 'set') {
          const { value } = await ElMessageBox.prompt('请输入新成员值', '修改', { inputValue: String(row.value ?? ''), inputType: 'textarea' })
          const res1 = await redisAPI.runCmd(inst, db, { command: 'SREM', args: [key, row.value] })
          const res2 = await redisAPI.runCmd(inst, db, { command: 'SADD', args: [key, value] })
          if ((res1.data?.code === 0 || res1.data?.code === 200) && (res2.data?.code === 0 || res2.data?.code === 200)) ElMessage.success('修改成功'); else throw new Error('修改失败')
        } else if (type === 'zset') {
          const { value } = await ElMessageBox.prompt(`成员 ${row.value} 新分数`, '修改', { inputValue: String(row.score ?? 0), inputPattern: /^-?\d+(\.\d+)?$/, inputErrorMessage: '请输入数字' })
          const res = await redisAPI.runCmd(inst, db, { command: 'ZADD', args: [key, String(value), row.value] })
          if (res.data?.code === 0 || res.data?.code === 200) ElMessage.success('修改成功'); else throw new Error(res.data?.message || '修改失败')
        } else {
          ElMessage.warning('该类型暂不支持修改')
          return
        }
        await this.reloadKeyData(tab)
      } catch (e) {
        if (e !== 'cancel') ElMessage.error(e.message || '修改失败')
      }
    },
    async handleDeleteRow(row, tab) {
      const type = (tab.meta?.type || '').toLowerCase()
      const key = tab.connection?.key
      const db = tab.node.db
      const inst = tab.node.instanceId
      try {
        await ElMessageBox.confirm('确定删除吗？', '删除', { type: 'warning' })
        let res
        if (type === 'string') {
          res = await redisAPI.runCmd(inst, db, { command: 'DEL', args: [key] })
        } else if (type === 'hash') {
          res = await redisAPI.runCmd(inst, db, { command: 'HDEL', args: [key, row.field] })
        } else if (type === 'list') {
          // 删除一个与当前值相等的成员
          res = await redisAPI.runCmd(inst, db, { command: 'LREM', args: [key, '1', row.value] })
        } else if (type === 'set') {
          res = await redisAPI.runCmd(inst, db, { command: 'SREM', args: [key, row.value] })
        } else if (type === 'zset') {
          res = await redisAPI.runCmd(inst, db, { command: 'ZREM', args: [key, row.value] })
        } else {
          ElMessage.warning('该类型暂不支持删除')
          return
        }
        if (res.data?.code === 0 || res.data?.code === 200) {
          ElMessage.success('删除成功')
          await this.reloadKeyData(tab)
        } else {
          throw new Error(res.data?.message || '删除失败')
        }
      } catch (e) {
        if (e !== 'cancel') ElMessage.error(e.message || '删除失败')
      }
    }
  }
}
</script>

<style scoped>
.db-terminal-page { padding: 20px; min-height: 100vh; background: #fff; }
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
.tools-row { display:flex; align-items:center; gap:10px; margin-bottom:8px; }
.inner-table-wrap { flex: 1; min-height: 0; }
.inner-table-wrap .el-table { height: 100%; }
.tab-label-text { display: inline-block; max-width: 220px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; vertical-align: middle; }
.key-value-text { display: inline-block; max-width: 360px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; vertical-align: middle; }
</style>

 
