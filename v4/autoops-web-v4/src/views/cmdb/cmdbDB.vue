<template>
  <div class="cmdb-db-management">
    <el-card shadow="hover" class="db-card">
      <template #header>
        <div class="card-header">
          <span class="title">数据库实例管理</span>
        </div>
      </template>

      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :inline="true" :model="queryParams" ref="queryForm" class="search-form">
          <el-form-item label="实例名称">
            <el-input
              v-model="queryParams.name"
              placeholder="请输入实例名称"
              clearable
              size="small"
              style="width: 200px"
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item label="数据库类型">
            <el-select
              v-model="queryParams.dbType"
              placeholder="请选择数据库类型"
              clearable
              size="small"
              style="width: 180px"
            >
              <el-option
                v-for="item in dbTypeOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img
                    :src="getDbIcon(item.value)"
                    :alt="item.label"
                    style="width: 16px; height: 16px; object-fit: contain; flex-shrink: 0;"
                  />
                  <span>{{ item.label }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="主机地址">
            <el-input
              v-model="queryParams.host"
              placeholder="请输入主机地址"
              clearable
              size="small"
              style="width: 200px"
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select
              v-model="queryParams.status"
              placeholder="请选择状态"
              clearable
              size="small"
              style="width: 120px"
            >
              <el-option label="启用" :value="1" />
              <el-option label="禁用" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="small" @click="handleQuery">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button type="warning" size="small" @click="resetQuery">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 操作按钮 -->
      <div class="action-section">
        <el-button type="success" size="small" v-authority="['cmdb:db:add']" @click="showAddDialog">
          <el-icon><Plus /></el-icon>
          创建实例
        </el-button>
        <el-button type="primary" size="small" v-authority="['cmdb:db:dbms']" @click="goToTerminalPage">
          <el-icon><Monitor /></el-icon>
          DBMS终端
        </el-button>
        <el-button type="primary" size="small" v-authority="['cmdb:db:redis']" @click="goToRedisTerminalPage">
          <el-icon><Monitor /></el-icon>
          Redis终端
        </el-button>
        <el-button type="primary" size="small" v-authority="['cmdb:db:es']" @click="goToESTerminalPage">
          <el-icon><Monitor /></el-icon>
          ES终端
        </el-button>
        <el-button type="primary" size="small" v-authority="['cmdb:db:mongo']" @click="goToMongoTerminalPage">
          <el-icon><Monitor /></el-icon>
          Mongo终端
        </el-button>
      </div>

      <!-- 数据库实例列表表格 -->
      <div class="table-section">
        <el-table
          :data="dbList"
          v-loading="loading"
          stripe
          style="width: 100%"
          class="db-table"
        >
          <el-table-column prop="name" label="实例名称" min-width="200">
            <template #default="{ row }">
              <div class="db-name-container">
                <img :src="getDbIcon(row.dbType || row.type)" :alt="getDbName(row.dbType || row.type)" class="db-icon" />
                <span class="db-name-text">{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="dbType" label="数据库类型" min-width="140">
            <template #default="{ row }">
              <el-tag
                :type="getDbTagType(row.dbType || row.type)"
                size="small"
                effect="dark"
              >
                {{ getDbName(row.dbType || row.type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="连接信息" min-width="220">
            <template #default="{ row }">
              <span>
                <!-- 优先显示 MongoDB 的连接 URI，其次展示 host:port -->
                {{ (row.connectionConfig && row.connectionConfig.uri) ? row.connectionConfig.uri : `${row.host || ''}${row.port ? `:${row.port}` : ''}` }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="username" label="用户名" min-width="120" />
          <el-table-column prop="status" label="状态" min-width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="创建时间" min-width="180" />
          <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-tooltip content="测试连接" placement="top">
                  <el-button
                    type="primary"
                    size="small"
                    v-authority="['cmdb:db:test']"
                    circle
                    @click="testConnection(row)"
                  >
                    <el-icon><Link /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="查看密码" placement="top">
                  <el-button
                    type="info"
                    size="small"
                    v-authority="['cmdb:db:passwd']"
                    circle
                    @click="viewPassword(row)"
                  >
                    <el-icon><View /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="修改" placement="top">
                  <el-button
                    type="warning"
                    size="small"
                    circle
                    v-authority="['cmdb:db:edit']"
                    @click="showEditDialog(row)"
                  >
                    <el-icon><Edit /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button
                    type="danger"
                    size="small"
                    circle
                    v-authority="['cmdb:db:delete']"
                    @click="handleDelete(row)"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination-section">
        <el-pagination
          v-model:current-page="queryParams.pageNum"
          v-model:page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    
    <!--新增/编辑对话框-->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="45%">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="110px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="实例名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入实例名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="数据库类型" prop="type">
              <el-select v-model="formData.type" placeholder="请选择数据库类型" style="width: 100%">
                <el-option
                  v-for="item in dbTypeOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                >
                  <div style="display: flex; align-items: center; gap: 8px;">
                    <img
                      :src="getDbIcon(item.value)"
                      :alt="item.label"
                      style="width: 16px; height: 16px; object-fit: contain;"
                    />
                    <span>{{ item.label }}</span>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="主机地址" prop="host">
              <el-input v-model="formData.host" placeholder="请输入主机地址" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="端口" prop="port">
              <el-input-number
                v-model="formData.port"
                :min="1"
                :max="65535"
                placeholder="请输入端口"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="formData.username" placeholder="可选：用户名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码" prop="password">
              <el-input
                v-model="formData.password"
                type="password"
                placeholder="可选：密码"
                show-password
              />
            </el-form-item>
          </el-col>
        </el-row>
        <!-- Redis 默认DB（可选） -->
        <el-row v-if="formData.type === 'redis'" :gutter="20">
          <el-col :span="12">
            <el-form-item label="默认DB" prop="db">
              <el-input-number v-model="formData.db" :min="0" :max="1024" placeholder="默认 0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <!-- 仅 MySQL/PostgreSQL 显示网络与参数 -->
        <el-row v-if="formData.type === 'mysql' || formData.type === 'postgres'" :gutter="20">
          <el-col :span="12">
            <el-form-item label="网络类型" prop="network">
              <el-input v-model="formData.network" placeholder="默认 tcp" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="连接参数" prop="params">
              <el-input v-model="formData.params" placeholder="如: charset=utf8mb4" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="备注" prop="remark">
              <el-input
                v-model="formData.remark"
                type="textarea"
                :rows="3"
                placeholder="请输入备注信息"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="submitForm">确 定</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 查看密码对话框 -->
    <el-dialog :title="passwordDialog.title" v-model="passwordDialog.visible" width="40%">
      <div>
        <el-alert
          v-if="passwordDialog.isMongo"
          type="info"
          show-icon
          class="mb10"
          description="MongoDB 返回完整连接 URI，包含用户名与密码。"
        />
        <el-input
          type="textarea"
          :rows="passwordDialog.isMongo ? 4 : 2"
          readonly
          v-model="passwordDialog.value"
        />
        <div style="margin-top: 10px; display: flex; gap: 8px; justify-content: flex-end;">
          <el-button @click="passwordDialog.visible = false">关 闭</el-button>
          <el-button type="primary" @click="copyToClipboard(passwordDialog.value)">
            <el-icon><CopyDocument /></el-icon>
            复制
          </el-button>
        </div>
      </div>
    </el-dialog>
    </el-card>
  </div>
</template>

<script>
import { dbAPI, redisAPI, esAPI, mongoAPI } from '@/api/database'
import { Edit, Delete, Search, Refresh, Plus, Link, Monitor, View, CopyDocument } from '@element-plus/icons-vue'

export default {
  components: {
    Edit,
    Delete,
    Search,
    Refresh,
    Plus,
    Link,
    Monitor,
    View,
    CopyDocument
  },
  data() {
    return {
      queryParams: {
        name: '',
        dbType: '',
        host: '',
        status: undefined,
        pageNum: 1,
        pageSize: 10
      },
      dbTypeOptions: [
        { label: 'MySQL', value: 'mysql' },
        { label: 'PostgreSQL', value: 'postgres' },
        { label: 'MongoDB', value: 'mongodb' },
        { label: 'Redis', value: 'redis' },
        { label: 'Elasticsearch', value: 'elasticsearch' }
      ],
      loading: false,
      dbList: [],
      total: 0,
      dialogVisible: false,
      dialogTitle: '',
      passwordDialog: { visible: false, title: '查看密码', value: '', isMongo: false },
      formData: {
        id: '',
        name: '',
        type: 'mysql',
        host: '',
        port: 3306,
        // mysql/postgres 专用
        network: 'tcp',
        params: '',
        // 通用
        username: '',
        password: '',
        remark: '',
        // redis 专用（最小化参数集）
        db: 0,
        // es 专用（协议省略则默认 http）
        protocol: 'http'
      },
      formRules: {
        name: [{ required: true, message: '请输入实例名称', trigger: 'blur' }],
        type: [{ required: true, message: '请选择数据库类型', trigger: 'change' }],
        host: [{ required: true, message: '请输入主机地址', trigger: 'blur' }],
        port: [{ required: true, message: '请输入端口', trigger: 'blur' }]
      }
    }
  },
  methods: {
    // 获取数据库图标
    getDbIcon(type) {
      const iconMap = {
        'mysql': require('@/assets/image/mysqldb.svg'),
        'postgres': require('@/assets/image/PostgreSQL.svg'),
        'mongodb': require('@/assets/image/mongodb.svg'),
        'redis': require('@/assets/image/redis.svg'),
        'elasticsearch': require('@/assets/image/Elasticsearch.svg')
      }
      return iconMap[type] || require('@/assets/image/mysqldb.svg')
    },
    // 查看密码/凭据
    async viewPassword(row) {
      try {
        const t = (row.dbType || row.type || '').toLowerCase()
        if (!t) {
          this.$message.error('无法识别数据库类型')
          return
        }
        const res = await dbAPI.getInstancePassword(t, row.id)
        if (res.data?.code === 200) {
          const pwd = res.data?.data?.password || ''
          const isMongo = t === 'mongodb'
          this.passwordDialog = {
            visible: true,
            title: `${row.name} 的${isMongo ? '连接URI' : '密码'}`,
            value: pwd,
            isMongo
          }
        } else {
          this.$message.error(res.data?.message || '获取密码失败')
        }
      } catch (e) {
        this.$message.error('获取密码失败: ' + (e.message || ''))
      }
    },
    copyToClipboard(text) {
      if (!text) return
      if (navigator?.clipboard?.writeText) {
        navigator.clipboard.writeText(text).then(() => this.$message.success('已复制'))
          .catch(() => this.$message.error('复制失败'))
      } else {
        const ta = document.createElement('textarea')
        ta.value = text
        document.body.appendChild(ta)
        ta.select()
        try { document.execCommand('copy'); this.$message.success('已复制') } catch { this.$message.error('复制失败') }
        document.body.removeChild(ta)
      }
    },

    // 获取数据库名称
    getDbName(type) {
      const nameMap = {
        'mysql': 'MySQL',
        'postgres': 'PostgreSQL',
        'mongodb': 'MongoDB',
        'redis': 'Redis',
        'elasticsearch': 'Elasticsearch'
      }
      return nameMap[type] || 'MySQL'
    },

    // 获取数据库标签类型
    getDbTagType(type) {
      const tagTypeMap = {
        'mysql': 'success',
        'postgres': 'warning',
        'mongodb': 'success',
        'redis': 'danger',
        'elasticsearch': 'primary'
      }
      return tagTypeMap[type] || 'success'
    },

    // 跳转到数据库终端页面
    goToTerminalPage() {
      this.$router.push({
        path: '/cmdb/dbterminal'
      })
    },
    // 跳转到 Redis 终端
    goToRedisTerminalPage() {
      this.$router.push({ path: '/cmdb/redis-terminal' })
    },
    // 跳转到 ES 终端
    goToESTerminalPage() {
      this.$router.push({ path: '/cmdb/es-terminal' })
    },
    // 跳转到 Mongo 终端
    goToMongoTerminalPage() {
      this.$router.push({ path: '/cmdb/mongo-terminal' })
    },

    // 获取数据库实例列表
    async getList() {
      this.loading = true
      try {
        const res = await dbAPI.getInstanceList(this.queryParams)
        console.log('API响应:', res)

        if (res.data?.code === 200) {
          // 处理不同的后端响应格式
          let dataList = res.data?.data

          // 如果data是数组，直接使用
          if (Array.isArray(dataList)) {
            this.dbList = dataList
            this.total = res.data?.total || dataList.length
          }
          // 如果data是对象且包含list属性
          else if (dataList && Array.isArray(dataList.list)) {
            this.dbList = dataList.list
            this.total = dataList.total || dataList.list.length
          }
          // 如果data是对象且包含data属性
          else if (dataList && Array.isArray(dataList.data)) {
            this.dbList = dataList.data
            this.total = dataList.total || dataList.data.length
          }
          // 默认为空数组
          else {
            this.dbList = []
            this.total = 0
            console.warn('未知的数据格式:', dataList)
          }
        } else {
          this.dbList = []
          this.total = 0
          this.$message.error(res.data?.message || '获取实例列表失败')
        }
      } catch (error) {
        console.error('获取实例列表失败:', error)
        this.dbList = []
        this.total = 0
        this.$message.error('获取实例列表失败: ' + (error.message || ''))
      } finally {
        this.loading = false
      }
    },

    // 搜索
    handleQuery() {
      this.queryParams.pageNum = 1
      this.getList()
    },

    // 重置搜索
    resetQuery() {
      this.queryParams = {
        name: '',
        dbType: '',
        host: '',
        status: undefined,
        pageNum: 1,
        pageSize: 10
      }
      this.getList()
    },

    // 分页相关
    handleSizeChange(val) {
      this.queryParams.pageSize = val
      this.getList()
    },

    handlePageChange(val) {
      this.queryParams.pageNum = val
      this.getList()
    },

    // 显示新增对话框（根据类型展示最小必填项）
    showAddDialog() {
      this.dialogTitle = '创建数据库实例'
      this.$nextTick(() => {
        if (this.$refs.formRef) {
          this.$refs.formRef.resetFields()
        }
        this.formData = {
          name: '',
          type: 'mysql',
          host: '',
          port: 3306,
          network: 'tcp',
          params: '',
          username: '',
          password: '',
          remark: '',
          db: 0,
          protocol: 'http'
        }
        this.dialogVisible = true
      })
    },

    // 显示编辑对话框
    showEditDialog(row) {
      this.dialogTitle = '修改数据库实例'
      this.$nextTick(() => {
        const t = (row.dbType || row.type || '').toLowerCase()
        const cfg = row.connectionConfig || {}
        const base = {
          id: row.id,
          name: row.name,
          type: t || 'mysql',
          host: row.host || '',
          port: row.port || (t === 'mysql' ? 3306 : t === 'postgres' ? 5432 : t === 'redis' ? 6379 : t === 'elasticsearch' ? 9200 : t === 'mongodb' ? 27017 : ''),
          username: row.username || '',
          password: '',
          remark: row.remark || ''
        }
        // 额外字段按类型映射
        if (t === 'mysql' || t === 'postgres') {
          base.network = cfg.network || 'tcp'
          base.params = cfg.params || ''
        } else if (t === 'redis') {
          base.db = typeof cfg.db === 'number' ? cfg.db : 0
        } else if (t === 'elasticsearch') {
          base.protocol = cfg.protocol || 'http'
        } else if (t === 'mongodb') {
          // 从 uri 尝试解析 host:port（容错）
          const uri = cfg.uri || ''
          try {
            if (uri && uri.includes('://')) {
              const u = new URL(uri)
              base.host = u.hostname || base.host
              base.port = Number(u.port) || base.port
              base.username = row.username || u.username || ''
            } else if (uri && uri.includes(':')) {
              const [h, p] = uri.split(':')
              base.host = base.host || h
              base.port = base.port || Number(p)
            }
          } catch (_) {}
        }
        this.formData = base
        this.dialogVisible = true
      })
    },

    // 测试连接：按类型调用“列库/索引/DBs”接口验证
    async testConnection(row) {
      const loading = this.$loading({
        lock: true,
        text: '正在测试连接...',
        background: 'rgba(0, 0, 0, 0.7)'
      })

      try {
        const t = (row.dbType || row.type || '').toLowerCase()
        let ok = false
        let count = 0
        if (t === 'redis') {
          const res = await redisAPI.getDBs(row.id)
          if (res.data?.code === 200 || res.data?.code === 0) {
            const d = res.data?.data || {}
            const dbs = Array.isArray(d.dbs) ? d.dbs : []
            count = dbs.length || 0
            ok = true
          }
        } else if (t === 'elasticsearch') {
          const res = await esAPI.proxy(row.id, '_cat/indices?format=json', 'get')
          let indices = []
          if (res.data?.code === 200 || res.data?.code === 0) {
            indices = res.data?.data || []
          } else if (Array.isArray(res.data)) {
            indices = res.data
          }
          count = indices.length || 0
          ok = true
        } else if (t === 'mongodb') {
          const res = await mongoAPI.getDatabases(row.id)
          if (res.data?.code === 200 || res.data?.code === 0) {
            const dbs = res.data?.data?.Databases || res.data?.data?.databases || []
            count = Array.isArray(dbs) ? dbs.length : 0
            ok = true
          }
        } else {
          // mysql / postgres
          const res = await dbAPI.getInstanceDatabases(row.id)
          if (res.data?.code === 200) {
            const dbs = res.data?.data || []
            count = Array.isArray(dbs) ? dbs.length : 0
            ok = true
          }
        }

        if (ok) {
          this.$message.success(`连接正常！可用资源数量: ${count}`)
        } else {
          this.$message.error('连接失败或无权限访问资源')
        }
      } catch (error) {
        console.error('测试连接失败:', error)
        this.$message.error('测试连接失败: ' + (error.message || '未知错误'))
      } finally {
        loading.close()
      }
    },

    // 提交表单
    async submitForm() {
      try {
        await this.$refs.formRef.validate()
        const t = (this.formData.type || '').toLowerCase()
        let res
        if (t === 'redis') {
          const payload = {
            id: this.formData.id || undefined,
            name: this.formData.name,
            host: this.formData.host,
            port: this.formData.port,
            db: typeof this.formData.db === 'number' ? this.formData.db : 0,
            username: this.formData.username || '',
            remark: this.formData.remark || ''
          }
          if (this.formData.password) payload.password = this.formData.password
          // 省略 mode，默认 standalone 由后端兜底或无需
          res = payload.id ? await redisAPI.updateInstance(payload) : await redisAPI.createInstance(payload)
        } else if (t === 'elasticsearch') {
          const payload = {
            id: this.formData.id || undefined,
            name: this.formData.name,
            protocol: this.formData.protocol || 'http',
            host: this.formData.host,
            port: this.formData.port,
            username: this.formData.username || '',
            remark: this.formData.remark || ''
          }
          if (this.formData.password && this.formData.password !== '******') payload.password = this.formData.password
          res = payload.id ? await esAPI.updateInstance(payload) : await esAPI.createInstance(payload)
        } else if (t === 'mongodb') {
          const payload = {
            id: this.formData.id || undefined,
            name: this.formData.name,
            host: this.formData.host,
            port: this.formData.port,
            username: this.formData.username || '',
            remark: this.formData.remark || ''
          }
          if (this.formData.password && this.formData.password !== '******') payload.password = this.formData.password
          res = payload.id ? await mongoAPI.updateInstance(payload) : await mongoAPI.createInstance(payload)
        } else {
          // 默认 mysql（或 postgres）
          const payload = {
            id: this.formData.id || undefined,
            name: this.formData.name,
            type: t || 'mysql',
            host: this.formData.host,
            port: this.formData.port,
            network: this.formData.network || 'tcp',
            params: this.formData.params || '',
            username: this.formData.username,
            remark: this.formData.remark || ''
          }
          if (this.formData.password) payload.password = this.formData.password
          res = payload.id ? await dbAPI.updateInstance(payload) : await dbAPI.createInstance(payload)
        }

        if (res.data?.code === 200 || res.data?.code === 0) {
          this.$message.success(this.formData.id ? '修改成功' : '创建成功')
          this.dialogVisible = false
          this.getList()
        } else {
          this.$message.error(res.data?.message || (this.formData.id ? '修改失败' : '创建失败'))
        }
      } catch (error) {
        if (error !== false) { // 排除表单验证失败的情况
          console.error('操作失败:', error)
          this.$message.error('操作失败: ' + (error.message || '未知错误'))
        }
      }
    },

    // 删除数据库实例（根据类型选择不同 API）
    async handleDelete(row) {
      try {
        await this.$confirm(`确定删除数据库实例"${row.name}"?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        const t = (row.dbType || row.type || '').toLowerCase()
        let res
        if (t === 'redis') res = await redisAPI.deleteInstance(row.id)
        else if (t === 'elasticsearch') res = await esAPI.deleteInstance(row.id)
        else if (t === 'mongodb') res = await mongoAPI.deleteInstance(row.id)
        else res = await dbAPI.deleteInstance(row.id)

        if (res.data?.code === 200 || res.data?.code === 0) {
          this.$message.success('删除成功')
          this.getList()
        } else {
          this.$message.error(res.data?.message || '删除失败')
        }
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除失败:', error)
          this.$message.error('删除失败: ' + (error.message || '未知错误'))
        }
      }
    }
  },
  created() {
    this.getList()
  },

  watch: {
    'formData.type'(newType) {
      // 根据数据库类型自动设置默认端口和参数
      if (!this.formData.id) { // 只在新增时自动设置
        if (newType === 'mysql') {
          this.formData.port = 3306
          this.formData.params = 'charset=utf8mb4&parseTime=true'
          this.formData.network = 'tcp'
        } else if (newType === 'postgres') {
          this.formData.port = 5432
          this.formData.params = 'sslmode=disable'
        } else if (newType === 'mongodb') {
          this.formData.port = 27017
          this.formData.params = ''
        } else if (newType === 'redis') {
          this.formData.port = 6379
          this.formData.db = 0
          this.formData.params = ''
        } else if (newType === 'elasticsearch') {
          this.formData.port = 9200
          this.formData.protocol = 'http'
          this.formData.params = ''
        }
      }
    }
  }
}
</script>

<style scoped>
/* 现代化科技感设计风格 */

.cmdb-db-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.db-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  min-height: 800px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

/* 搜索区域样式 */
.search-section {
  margin-bottom: 20px;
  padding: 20px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

/* 操作按钮区域 */
.action-section {
  margin-bottom: 20px;
  padding-left: 0;
}

/* 表格区域样式 */
.table-section {
  margin-top: 0;
}

.db-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.db-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.db-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.db-table :deep(.el-table__header th .cell) {
  color: #2c3e50 !important;
  font-weight: 700 !important;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.db-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.db-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* 数据库名称容器样式 */
.db-name-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.db-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.db-name-text {
  font-weight: 600;
  color: #2c3e50;
}

/* 操作按钮样式 */
.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.operation-buttons .el-button {
  transition: all 0.3s ease;
}

.operation-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

/* 分页样式 */
.pagination-section {
  margin-top: 20px;
  padding: 20px;
  display: flex;
  justify-content: flex-end;
  background: transparent;
  border-radius: 0;
}

.pagination-section :deep(.el-pagination) {
  gap: 8px;
}

.pagination-section :deep(.el-pagination button),
.pagination-section :deep(.el-pager li) {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 8px;
  transition: all 0.3s ease;
}

.pagination-section :deep(.el-pagination button:hover),
.pagination-section :deep(.el-pager li:hover) {
  background: #667eea;
  color: #fff;
  transform: translateY(-2px);
}

.pagination-section :deep(.el-pager li.is-active) {
  background: linear-gradient(45deg, #667eea, #764ba2);
  color: #fff;
}

/* 标签样式优化 */
.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

/* 按钮样式优化 */
.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 输入框和选择器样式 */
.el-input :deep(.el-input__wrapper),
.el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__wrapper):hover,
.el-select :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.el-input :deep(.el-input__wrapper.is-focus),
.el-select :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #2c3e50;
}

/* 加载遮罩样式 */
.el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

/* 对话框样式 */
.el-dialog {
  border-radius: 12px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
