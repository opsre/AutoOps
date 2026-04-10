<template>
  <div class="alarm-notify-management">
    <div class="integrated-layout custom-scrollbar">
      <!-- 左侧：菜单导航 -->
      <div class="sidebar-section">
        <div class="sidebar-header">
          <img src="https://img.icons8.com/3d-fluency/94/bell.png" class="header-icon" />
          <span class="header-title">通知配置</span>
        </div>

        <el-menu 
          :default-active="activeTab" 
          @select="handleMenuSelect" 
          class="group-menu transparent-bg"
        >
          <el-menu-item index="templates">
            <img src="https://img.icons8.com/3d-fluency/94/document.png" class="menu-colored-icon" />
            <template #title><span class="menu-name">告警模版</span></template>
          </el-menu-item>
          <el-menu-item index="routers">
            <img src="https://img.icons8.com/3d-fluency/94/network-cable.png" class="menu-colored-icon" />
            <template #title><span class="menu-name">告警路由</span></template>
          </el-menu-item>
          <el-menu-item index="records">
            <img src="https://img.icons8.com/3d-fluency/94/clock.png" class="menu-colored-icon" />
            <template #title><span class="menu-name">通知历史</span></template>
          </el-menu-item>
        </el-menu>
      </div>

      <!-- 右侧：主内容区 -->
      <div class="main-section">
        <div class="main-header">
          <span class="header-title">
            <el-icon style="margin-right: 8px"><component :is="activeIcon" /></el-icon> 
            {{ activeTitle }}
          </span>
          <div class="right-actions">
            <el-button v-if="activeTab === 'routers'" type="warning" plain icon="Refresh" @click="reloadRouters">重载路由配置</el-button>
            <el-button v-if="activeTab === 'templates'" type="primary" icon="Plus" @click="openTemplateDialog()">新建模版</el-button>
            <el-button v-if="activeTab === 'routers'" type="primary" icon="Plus" @click="openRouterDialog()">新建路由</el-button>
            <el-button v-if="activeTab === 'records'" type="danger" icon="Delete" @click="cleanRecords()">清空历史</el-button>
          </div>
        </div>

        <div class="content-body custom-scrollbar">
          <!-- 搜索区域 -->
          <div class="search-section">
            <!-- 模版搜索 -->
            <el-form v-if="activeTab === 'templates'" :inline="true" :model="templateQuery" class="search-form">
              <el-form-item label="名称">
                <el-input v-model="templateQuery.tplname" placeholder="请输入模版名称" clearable @keyup.enter="fetchData" style="width: 180px" />
              </el-form-item>
              <el-form-item label="类型">
                <el-select v-model="templateQuery.tpltype" clearable placeholder="类型" @change="fetchData" style="width: 120px">
                  <el-option label="微信" value="wx">
                    <span style="display: flex; align-items: center; gap: 6px;"><el-icon><ChatDotRound /></el-icon> 微信</span>
                  </el-option>
                  <el-option label="钉钉" value="dd">
                    <span style="display: flex; align-items: center; gap: 6px;"><el-icon><Bell /></el-icon> 钉钉</span>
                  </el-option>
                  <el-option label="飞书" value="fs">
                    <span style="display: flex; align-items: center; gap: 6px;"><el-icon><Promotion /></el-icon> 飞书</span>
                  </el-option>
                  <el-option label="邮箱" value="email">
                    <span style="display: flex; align-items: center; gap: 6px;"><el-icon><Message /></el-icon> 邮箱</span>
                  </el-option>
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="Search" @click="fetchData">查询</el-button>
                <el-button type="warning" icon="Refresh" @click="resetQuery">重置</el-button>
              </el-form-item>
            </el-form>

            <!-- 路由搜索 -->
            <el-form v-if="activeTab === 'routers'" :inline="true" :model="routerQuery" class="search-form">
              <el-form-item label="名称">
                <el-input v-model="routerQuery.name" placeholder="请输入路由名" clearable @keyup.enter="fetchData" style="width: 180px" />
              </el-form-item>
              <el-form-item label="推送地址">
                <el-input v-model="routerQuery.urlOrPhone" placeholder="URL 或 手机号" clearable @keyup.enter="fetchData" style="width: 180px" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="Search" @click="fetchData">查询</el-button>
                <el-button type="warning" icon="Refresh" @click="resetQuery">重置</el-button>
              </el-form-item>
            </el-form>

            <!-- 历史搜索 -->
            <el-form v-if="activeTab === 'records'" :inline="true" :model="recordQuery" class="search-form">
              <el-form-item label="告警名称">
                <el-input v-model="recordQuery.alertname" placeholder="模糊匹配" clearable @keyup.enter="fetchData" style="width: 150px" />
              </el-form-item>
              <el-form-item label="级别">
                <el-select v-model="recordQuery.alertLevel" clearable placeholder="级别" @change="fetchData" style="width: 110px">
                  <el-option label="Critical" value="critical" />
                  <el-option label="Warning" value="warning" />
                  <el-option label="Info" value="info" />
                </el-select>
              </el-form-item>
              <el-form-item label="故障实例">
                <el-input v-model="recordQuery.instance" placeholder="模糊匹配" clearable @keyup.enter="fetchData" style="width: 150px" />
              </el-form-item>
              <el-form-item label="状态">
                <el-select v-model="recordQuery.alertStatus" clearable placeholder="状态" @change="fetchData" style="width: 110px">
                  <el-option label="触发" value="firing" />
                  <el-option label="恢复" value="resolved" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="Search" @click="fetchData">查询</el-button>
                <el-button type="warning" icon="Refresh" @click="resetQuery">重置</el-button>
              </el-form-item>
            </el-form>
          </div>

          <!-- 模版管理 -->
          <div v-if="activeTab === 'templates'">
            <el-table :data="templateList" v-loading="loading" border stripe class="custom-table" :header-cell-style="{ background: 'rgba(102, 126, 234, 0.1)', color: '#2c3e50', fontWeight: '600' }">
              <el-table-column prop="Tplname" label="模版名称" min-width="150" />
              <el-table-column prop="Tpltype" label="类型" width="120" align="center">
                <template #default="{ row }">
                  <el-tag :type="getTypeColor(row.Tpltype)" style="display: inline-flex; align-items: center; gap: 4px;">
                    <el-icon><component :is="getTypeIcon(row.Tpltype)" /></el-icon>
                    {{ formatType(row.Tpltype) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="Tpluse" label="关联模块" width="120" align="center" />
              <el-table-column label="操作" width="180" align="center" fixed="right">
                <template #default="{ row }">
                  <div class="operation-buttons">
                    <el-button link type="primary" icon="EditPen" @click="openTemplateDialog(row)">编辑</el-button>
                    <el-button link type="danger" icon="Delete" @click="deleteTemplate(row.Id)">删除</el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
            <div class="pagination-section">
              <el-pagination
                background
                layout="total, sizes, prev, pager, next, jumper"
                :total="templateTotal"
                v-model:current-page="templateQuery.page"
                v-model:page-size="templateQuery.pageSize"
                :page-sizes="[10, 15, 30, 50]"
                @current-change="fetchTemplates"
                @size-change="fetchTemplates"
              />
            </div>
          </div>

          <!-- 路由管理 -->
          <div v-if="activeTab === 'routers'">
            <el-table :data="routerList" v-loading="loading" border stripe class="custom-table" :header-cell-style="{ background: 'rgba(102, 126, 234, 0.1)', color: '#2c3e50', fontWeight: '600' }">
              <el-table-column prop="Name" label="路由名称" min-width="150" />
              <el-table-column label="匹配规则 (Rules)" min-width="250" show-overflow-tooltip>
                <template #default="{ row }">
                  <code class="promql-code rule-truncate" v-if="row.Rules && formatRules(row.Rules) !== '[]'">{{ formatRules(row.Rules) }}</code>
                  <span v-else class="empty-text">无</span>
                </template>
              </el-table-column>
              <el-table-column label="发送目标" min-width="200" show-overflow-tooltip>
                <template #default="{ row }">
                  <span class="url-text">{{ row.UrlOrPhone }}</span>
                </template>
              </el-table-column>
              <el-table-column label="@某人" prop="AtSomeOne" width="120" show-overflow-tooltip />
              <el-table-column label="At恢复" width="80" align="center">
                <template #default="{ row }">
                  <el-tag :type="row.AtSomeOneRR ? 'success' : 'info'">{{ row.AtSomeOneRR ? '是' : '否' }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="发恢复" width="80" align="center">
                <template #default="{ row }">
                  <el-tag :type="row.SendResolved ? 'success' : 'info'">{{ row.SendResolved ? '是' : '否' }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="180" align="center" fixed="right">
                <template #default="{ row }">
                  <div class="operation-buttons">
                    <el-button link type="primary" icon="EditPen" @click="openRouterDialog(row)">编辑</el-button>
                    <el-button link type="danger" icon="Delete" @click="deleteRouter(row.Id)">删除</el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
            <div class="pagination-section">
              <el-pagination
                background
                layout="total, sizes, prev, pager, next, jumper"
                :total="routerTotal"
                v-model:current-page="routerQuery.page"
                v-model:page-size="routerQuery.pageSize"
                :page-sizes="[10, 15, 30, 50]"
                @current-change="fetchRouters"
                @size-change="fetchRouters"
              />
            </div>
          </div>

          <!-- 告警历史 -->
          <div v-if="activeTab === 'records'">
            <el-table :data="recordList" v-loading="loading" border stripe class="custom-table" :header-cell-style="{ background: 'rgba(102, 126, 234, 0.1)', color: '#2c3e50', fontWeight: '600' }">
              <el-table-column prop="Alertname" label="告警名称" min-width="150" show-overflow-tooltip />
              <el-table-column prop="Instance" label="故障实例" min-width="150" show-overflow-tooltip />
              <el-table-column label="级别" width="80" align="center">
                <template #default="{ row }">
                  <el-tag :type="getSeverityColor(row.AlertLevel || tryGetSeverity(row.Labels))" effect="dark">
                    {{ row.AlertLevel || tryGetSeverity(row.Labels) || '未知' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="AlertStatus" label="状态" width="80" align="center">
                <template #default="{ row }">
                  <el-tag :type="row.AlertStatus === 'firing' ? 'danger' : 'success'">
                    {{ row.AlertStatus === 'firing' ? '触发' : '恢复' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="Summary" label="摘要" min-width="150" show-overflow-tooltip />
              <el-table-column prop="StartsAt" label="告警发生时间" width="160" align="center">
                <template #default="{ row }">
                  {{ formatDate(row.StartsAt) }}
                </template>
              </el-table-column>
              <el-table-column prop="CreatedTime" label="发送时间" width="160" align="center">
                <template #default="{ row }">
                  {{ formatDate(row.CreatedTime) }}
                </template>
              </el-table-column>
            </el-table>
            <div class="pagination-section">
              <el-pagination
                background
                layout="total, sizes, prev, pager, next, jumper"
                :total="recordTotal"
                v-model:current-page="recordQuery.page"
                v-model:page-size="recordQuery.pageSize"
                :page-sizes="[10, 15, 30, 50]"
                @current-change="fetchRecords"
                @size-change="fetchRecords"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 模版弹窗 -->
    <el-dialog :title="tplForm.Id ? '编辑告警模版' : '新建告警模版'" v-model="tplDialogVisible" width="700px" class="modern-dialog">
      <el-form :model="tplForm" ref="tplFormRef" label-width="120px" class="rule-form">
        <div class="form-section-title"><el-icon><Setting /></el-icon> 模版基础信息</div>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="模版名称" required>
              <el-input v-model="tplForm.Tplname" placeholder="如: Webhook-DingTalk" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通知类型" required>
              <el-select v-model="tplForm.Tpltype" style="width: 100%">
                <el-option label="微信" value="wx">
                  <span style="display: flex; align-items: center; gap: 6px;"><el-icon><ChatDotRound /></el-icon> 微信</span>
                </el-option>
                <el-option label="钉钉" value="dd">
                  <span style="display: flex; align-items: center; gap: 6px;"><el-icon><Bell /></el-icon> 钉钉</span>
                </el-option>
                <el-option label="飞书" value="fs">
                  <span style="display: flex; align-items: center; gap: 6px;"><el-icon><Promotion /></el-icon> 飞书</span>
                </el-option>
                <el-option label="邮箱" value="email">
                  <span style="display: flex; align-items: center; gap: 6px;"><el-icon><Message /></el-icon> 邮箱</span>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="关联模块">
              <el-input v-model="tplForm.Tpluse" placeholder="如: Prometheus" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Content-Type">
              <el-input v-model="tplForm.WebhookContentType" placeholder="如: application/json" />
            </el-form-item>
          </el-col>
        </el-row>

        <div class="form-section-title"><el-icon><Document /></el-icon> 模版内容定义 (Go Template)</div>
        <div class="code-editor-container">
          <div class="code-editor-header">
            <span class="mac-dot close"></span>
            <span class="mac-dot minimize"></span>
            <span class="mac-dot maximize"></span>
            <span class="file-name">template.tmpl</span>
          </div>
          <el-input v-model="tplForm.Tpl" type="textarea" class="code-input custom-scrollbar" :rows="8" placeholder="填写Go Template格式内容... 比如 {{.Alerts}}" />
        </div>
      </el-form>
      <template #footer>
        <el-button @click="tplDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitTemplate">保存</el-button>
      </template>
    </el-dialog>

    <!-- 路由弹窗 -->
    <el-dialog :title="routerForm.Id ? '编辑告警路由' : '新建告警路由'" v-model="routerDialogVisible" width="750px" class="modern-dialog">
      <el-form :model="routerForm" ref="routerFormRef" label-width="120px" class="rule-form">
        <div class="form-section-title"><el-icon><Position /></el-icon> 路由基础配置</div>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="路由名称" required>
              <el-input v-model="routerForm.RouterName" placeholder="如: DBA-Alert-Route" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="绑定模版" required>
              <el-select v-model="routerForm.RouterTplId" placeholder="选择渲染模版" style="width: 100%">
                <el-option v-for="t in templateList" :key="t.Id" :label="t.Tplname" :value="t.Id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="推送地址" required>
          <el-input v-model="routerForm.RouterPurl" placeholder="Webhook URL 或者邮箱地址" />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="@某人 (At)">
              <el-input v-model="routerForm.RouterPat" placeholder="如: @all, 或手机号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="At恢复(RR)">
              <el-switch v-model="routerForm.RouterPatRR" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="发送恢复提醒">
              <el-switch v-model="routerForm.RouterSendResolved" />
            </el-form-item>
          </el-col>
        </el-row>

        <div class="form-section-title"><el-icon><Filter /></el-icon> 标签匹配规则 (Rules)</div>
        <div v-for="(rule, index) in routerRules" :key="index" style="display: flex; gap: 10px; margin-bottom: 10px; align-items: center; padding-left: 50px;">
          <el-input v-model="rule.key" placeholder="标签键 (如 severity)" style="width: 150px;" />
          <el-select v-model="rule.type" placeholder="关系" style="width: 100px;">
            <el-option label="等于" value="等于" />
            <el-option label="不等于" value="不等于" />
            <el-option label="正则" value="正则" />
          </el-select>
          <el-input v-model="rule.value" placeholder="匹配值" style="width: 180px;" />
          <el-button type="danger" icon="Delete" circle plain @click="removeRule(index)" />
        </div>
        <el-button type="primary" plain icon="Plus" style="margin-left: 50px; margin-bottom: 15px;" @click="addRule">添加匹配规则</el-button>

      </el-form>
      <template #footer>
        <el-button @click="routerDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitRouter">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ChatDotRound, Bell, Promotion, Message, Collection, Document, Connection, Timer } from '@element-plus/icons-vue'
import {
  getAlertTemplates, createAlertTemplate, updateAlertTemplate, deleteAlertTemplate,
  getAlertRouters, createAlertRouter, updateAlertRouter, deleteAlertRouter, reloadAlertRouters,
  getAlertRecords, cleanAlertRecords
} from '@/api/monitor'

// 激活状态
const activeTab = ref('templates')
const loading = ref(false)

const activeTitle = computed(() => {
  if (activeTab.value === 'templates') return '告警模版管理'
  if (activeTab.value === 'routers') return '告警路由策略'
  return '告警通知历史'
})

const activeIcon = computed(() => {
  if (activeTab.value === 'templates') return Document
  if (activeTab.value === 'routers') return Connection
  return Timer
})

const handleMenuSelect = (index) => {
  activeTab.value = index
  fetchData()
}

// 数据列表
const templateList = ref([])
const templateTotal = ref(0)
const templateQuery = reactive({ page: 1, pageSize: 10, tplname: '', tpltype: '' })

const routerList = ref([])
const routerTotal = ref(0)
const routerQuery = reactive({ page: 1, pageSize: 10, name: '', urlOrPhone: '' })

const recordList = ref([])
const recordTotal = ref(0)
const recordQuery = reactive({ page: 1, pageSize: 10, alertname: '', alertStatus: '', alertLevel: '', instance: '' })

const resetQuery = () => {
  if (activeTab.value === 'templates') { templateQuery.tplname = ''; templateQuery.tpltype = ''; templateQuery.page = 1; }
  else if (activeTab.value === 'routers') { routerQuery.name = ''; routerQuery.urlOrPhone = ''; routerQuery.page = 1; }
  else { recordQuery.alertname = ''; recordQuery.alertStatus = ''; recordQuery.alertLevel = ''; recordQuery.instance = ''; recordQuery.page = 1; }
  fetchData()
}

const fetchData = () => {
  if (activeTab.value === 'templates') fetchTemplates()
  else if (activeTab.value === 'routers') fetchRouters()
  else fetchRecords()
}

onMounted(() => {
  fetchData()
})

// --- API ---
const fetchTemplates = async () => {
  loading.value = true
  try {
    const res = await getAlertTemplates(templateQuery)
    const data = res.data?.data || res.data || {}
    templateList.value = data.list || (Array.isArray(data) ? data : [])
    templateTotal.value = data.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchRouters = async () => {
  loading.value = true
  try {
    const res = await getAlertRouters(routerQuery)
    const data = res.data?.data || res.data || {}
    routerList.value = data.list || (Array.isArray(data) ? data : [])
    routerTotal.value = data.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchRecords = async () => {
  loading.value = true
  try {
    const res = await getAlertRecords(recordQuery)
    const data = res.data?.data || res.data || {}
    recordList.value = data.list || (Array.isArray(data) ? data : [])
    recordTotal.value = data.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

// --- Templates 管理 ---
const tplDialogVisible = ref(false)
const tplForm = reactive({ Id: null, Tplname: '', Tpltype: 'wx', Tpluse: 'Prometheus', Tpl: '', WebhookContentType: 'application/json' })

const openTemplateDialog = (row = null) => {
  if (row) {
    tplForm.Id = row.Id
    tplForm.Tplname = row.Tplname
    tplForm.Tpltype = row.Tpltype
    tplForm.Tpluse = row.Tpluse
    tplForm.Tpl = row.Tpl
    tplForm.WebhookContentType = row.WebhookContentType || 'application/json'
  } else {
    tplForm.Id = null
    tplForm.Tplname = ''
    tplForm.Tpltype = 'wx'
    tplForm.Tpluse = 'Prometheus'
    tplForm.Tpl = '## [告警提醒]\n\n告警详情: {{.Alerts}}'
    tplForm.WebhookContentType = 'application/json'
  }
  tplDialogVisible.value = true
}

const submitTemplate = async () => {
  if (!tplForm.Tplname || !tplForm.Tpl) return ElMessage.warning('名称与内容必填项')
  try {
    if (tplForm.Id) {
      await updateAlertTemplate(tplForm.Id, tplForm)
    } else {
      await createAlertTemplate(tplForm)
    }
    ElMessage.success('模版保存成功')
    tplDialogVisible.value = false
    fetchTemplates()
  } catch (e) {
    ElMessage.error('保存失败')
  }
}

const deleteTemplate = (id) => {
  ElMessageBox.confirm('确认删除该模版?', '警告', { type: 'warning' }).then(async () => {
    await deleteAlertTemplate(id)
    ElMessage.success('删除成功')
    fetchTemplates()
  }).catch(() => {})
}

// --- Routers 管理 ---
const routerDialogVisible = ref(false)
const routerForm = reactive({ Id: null, RouterName: '', RouterTplId: '', RouterPurl: '', RouterPat: '', RouterPatRR: false, RouterSendResolved: true })
const routerRules = ref([])

const openRouterDialog = (row = null) => {
  if (row && (row.Id || row.id)) {
    routerForm.Id = row.Id || row.id
    routerForm.RouterName = row.Name || row.RouterName || ''
    routerForm.RouterTplId = row.TplId || row.RouterTplId || ''
    routerForm.RouterPurl = row.UrlOrPhone || row.RouterPurl || ''
    routerForm.RouterPat = row.AtSomeOne || row.RouterPat || ''
    routerForm.RouterPatRR = row.AtSomeOneRR === true || row.RouterPatRR === true
    routerForm.RouterSendResolved = row.SendResolved === true || row.RouterSendResolved === true
    try {
      let parsed = []
      if (typeof row.Rules === 'string') {
        const cleanedStr = row.Rules.replace(/critical}/g, 'critical\\"}')
        try {
          parsed = JSON.parse(cleanedStr)
        } catch(e) { parsed = [] }
      } else if (Array.isArray(row.Rules)) {
        parsed = row.Rules
      }
      
      if (Array.isArray(parsed) && parsed.length > 0) {
        routerRules.value = parsed.map(p => ({
            key: p.key || '',
            type: p.type || '等于',
            value: p.value || ''
        }))
      } else if (typeof parsed === 'object' && parsed !== null) {
        const arr = []
        for (const [k, v] of Object.entries(parsed)) {
          arr.push({ key: k, type: '等于', value: v })
        }
        routerRules.value = arr.length > 0 ? arr : [{ key: '', type: '等于', value: '' }]
      } else {
        routerRules.value = [{ key: '', type: '等于', value: '' }]
      }
    } catch(e) { routerRules.value = [{ key: '', type: '等于', value: '' }] }
  } else {
    routerForm.Id = null
    routerForm.RouterName = ''
    routerForm.RouterTplId = ''
    routerForm.RouterPurl = ''
    routerForm.RouterPat = '@all'
    routerForm.RouterPatRR = false
    routerForm.RouterSendResolved = true
    routerRules.value = [{ key: '', type: '等于', value: '' }]
    if (templateList.value.length === 0) fetchTemplates()
  }
  routerDialogVisible.value = true
}

const addRule = () => routerRules.value.push({ key: '', type: '等于', value: '' })
const removeRule = (index) => routerRules.value.splice(index, 1)

const submitRouter = async () => {
  if (!routerForm.RouterName || !routerForm.RouterTplId || !routerForm.RouterPurl) return ElMessage.warning('必填项未填完整')
  try {
    const activeRules = routerRules.value.filter(r => r.key && r.value)
    
    const payload = { 
      ...routerForm, 
      RouterTplId: Number(routerForm.RouterTplId),
      Rules: activeRules
    }
    
    if (payload.Id) {
      console.log('Updating router with payload:', payload)
      await updateAlertRouter(payload.Id, payload) 
    } else {
      await createAlertRouter(payload)
    }
    ElMessage.success('路由保存成功')
    routerDialogVisible.value = false
    fetchRouters()
  } catch (e) {
    ElMessage.error('路由保存失败')
  }
}

const deleteRouter = (id) => {
  ElMessageBox.confirm('确认删除该路由策略?', '警告', { type: 'warning' }).then(async () => {
    await deleteAlertRouter(id)
    ElMessage.success('删除成功')
    fetchRouters()
  }).catch(() => {})
}

const reloadRouters = async () => {
  try {
    await reloadAlertRouters()
    ElMessage.success('路由策略重载成功')
  } catch (e) {
    ElMessage.error('重载失败')
  }
}

// --- Records 管理 ---
const cleanRecords = () => {
  ElMessageBox.confirm('即将清空所有通知发送历史，操作不可逆，继续吗?', '高危操作', { type: 'error' }).then(async () => {
    await cleanAlertRecords()
    ElMessage.success('历史记录已清空')
    recordQuery.page = 1
    fetchRecords()
  }).catch(() => {})
}

// UI Helper
const getTypeIcon = (type) => {
  const map = { 'wx': ChatDotRound, 'dd': Bell, 'fs': Promotion, 'email': Message }
  return map[type] || Collection
}

const getTypeColor = (type) => {
  const map = { 'wx': 'success', 'dd': 'primary', 'fs': 'warning', 'email': 'info' }
  return map[type] || 'info'
}
const formatType = (type) => {
  const map = { 'wx': '企业微信', 'dd': '钉钉', 'fs': '飞书', 'email': '邮件' }
  return map[type] || type
}

const getSeverityColor = (sev) => {
  if (!sev) return 'info'
  const low = sev.toLowerCase()
  if (low.includes('crit') || low.includes('fatal') || low.includes('严重')) return 'danger'
  if (low.includes('warn') || low.includes('警告')) return 'warning'
  return 'info'
}

const tryGetSeverity = (labelsStr) => {
  try {
    const l = JSON.parse(labelsStr || '{}')
    return l.severity || ''
  } catch (e) {
    return ''
  }
}

const formatRules = (rules) => {
  if (!rules) return ''
  if (typeof rules === 'string') return rules
  try {
    return JSON.stringify(rules)
  } catch(e) { return String(rules) }
}

const formatDate = (dateStr) => {
  if (!dateStr || dateStr.startsWith('0001-01-01')) return '-'
  return new Date(dateStr).toLocaleString()
}

</script>

<style scoped>
.alarm-notify-management {
  padding: 20px;
  min-height: calc(100vh - 120px);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.integrated-layout {
  display: flex;
  height: calc(100vh - 160px);
  min-height: 600px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
}

.sidebar-section {
  width: 250px;
  background: rgba(102, 126, 234, 0.03);
  border-right: 1px solid rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.main-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.main-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  background: rgba(255, 255, 255, 0.5);
}

.content-body {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-icon {
  width: 24px;
  height: 24px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
}

.transparent-bg {
  background-color: transparent !important;
  border-right: none;
}

.menu-colored-icon {
  width: 22px;
  height: 22px;
  margin-right: 12px;
  object-fit: contain;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.15));
  transition: transform 0.3s ease;
}

.el-menu-item:hover .menu-colored-icon {
  transform: scale(1.15) translateY(-2px);
}

.menu-name {
  font-weight: 600;
  color: #2c3e50;
}

/* 按钮 & 表格 */
:deep(.el-table) {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  margin-top: 10px;
}
:deep(.el-table__header th) { border: none; }
:deep(.el-table__body tr:hover > td) { background-color: rgba(102, 126, 234, 0.1) !important; }

.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.promql-code {
  background-color: rgba(103, 126, 234, 0.1);
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 13px;
  color: #d35400;
  font-family: Menlo, Monaco, Consolas, "Courier New", monospace;
}

.rule-truncate {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: bottom;
}

.url-text {
  font-size: 12px;
  color: #555;
  word-break: break-all;
}

.empty-text {
  color: #999;
  font-style: italic;
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

/* 弹窗及表单 */
:deep(.modern-dialog .el-dialog) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}
:deep(.modern-dialog .el-dialog__header) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
  padding: 20px 24px 16px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.2);
}
:deep(.modern-dialog .el-dialog__title) { color: #2c3e50; font-weight: 600; }
:deep(.modern-dialog .el-dialog__footer) { background: rgba(248, 249, 250, 0.8); border-radius: 0 0 16px 16px; }

.form-section-title {
  font-size: 15px;
  font-weight: 600;
  color: #5a6fd8;
  margin: 10px 0 20px 0;
  display: flex;
  align-items: center;
  gap: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.2);
}

/* 代码编辑器外观 */
.code-editor-container {
  background: #1e1e1f;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #333;
}

.code-editor-header {
  height: 30px;
  background: #252526;
  border-bottom: 1px solid #333;
  display: flex;
  align-items: center;
  padding: 0 12px;
  position: relative;
}

.mac-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 6px;
  display: inline-block;
}
.mac-dot.close { background: #ff5f56; }
.mac-dot.minimize { background: #ffbd2e; }
.mac-dot.maximize { background: #27c93f; }

.file-name {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  color: #9cdcfe;
  font-size: 13px;
  font-family: Menlo, Monaco, Consolas, monospace;
  opacity: 0.8;
}

.code-input :deep(textarea) {
  background-color: #1e1e1f;
  color: #d4d4d4;
  border: none;
  font-family: Menlo, Monaco, Consolas, monospace;
  padding: 12px;
  font-size: 14px;
  border-radius: 0 0 8px 8px;
  box-shadow: none;
}
.code-input :deep(textarea):focus { box-shadow: none; }

/* 滚动条美化 */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(102, 126, 234, 0.4);
  border-radius: 3px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
</style>
