
<template>
  <div class="alarm-rules-management">
    <el-row :gutter="20" style="height: 100%">
      <!-- 左侧：业务线与群组列表 -->
      <el-col :span="6">
        <el-card class="modern-card menu-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="header-title">
                <el-icon><FolderOpened /></el-icon> 告警群组
              </span>
              <el-button type="primary" link icon="Plus" @click="openGroupDialog()">新建群组</el-button>
            </div>
          </template>
          
          <div class="group-select-wrap">
            <div class="ds-selector" style="margin-bottom: 10px;">
              <el-select v-model="activeDataSourceId" placeholder="选择数据源" style="width: 100%;" class="ds-select" @change="handleDataSourceChange">
                <template #prefix>
                  <el-icon v-if="!activeDataSourceId"><Monitor /></el-icon>
                  <img v-else :src="getDataSourceIcon(activeDataSourceId)" class="ds-icon-small" />
                </template>
                <el-option v-for="ds in dataSources" :key="ds.id || ds.ID" :label="ds.name" :value="ds.id || ds.ID" class="ds-option">
                  <div class="ds-option-content">
                    <img :src="getDsTypeIcon(ds.type)" class="ds-icon" />
                    <span class="ds-name">{{ ds.name }}</span>
                    <el-tag size="small" type="info" class="ds-deploy-tag" v-if="ds.deployMethod">
                      <img :src="getDeployIcon(ds.deployMethod)" class="deploy-icon" v-if="getDeployIcon(ds.deployMethod)" />
                      {{ ds.deployMethod }}
                    </el-tag>
                  </div>
                </el-option>
              </el-select>
            </div>
            <el-select v-model="filterCmdbGroup" placeholder="筛选业务线/CMDB集群" clearable style="width: 100%; margin-bottom: 10px;">
              <el-option v-for="c in cmdbGroups" :key="c.id" :label="c.name" :value="c.id" />
            </el-select>
          </div>

          <!-- 群组菜单 -->
          <el-menu 
            :default-active="activeGroupId" 
            @select="handleGroupSelect" 
            class="group-menu transparent-bg"
            v-loading="loadingGroups"
          >
            <el-menu-item index="all">
              <img src="https://img.icons8.com/3d-fluency/94/layers.png" class="group-colored-icon" />
              <template #title>
                <div class="menu-item-content">
                  <span class="group-name">全部规则</span>
                </div>
              </template>
            </el-menu-item>
            <el-menu-item v-for="g in filteredGroupList" :key="g.ID || g.id" :index="(g.ID || g.id) + ''">
              <img src="https://img.icons8.com/3d-fluency/94/folder-invoices.png" class="group-colored-icon" />
              <template #title>
                <div class="menu-item-content">
                  <span class="group-name" :title="g.group_name">
                    {{ g.group_name }}
                    <el-tag v-if="getClusterLabel(g.labels)" size="small" type="info" class="cluster-tag" disable-transitions>{{ getClusterLabel(g.labels) }}</el-tag>
                  </span>
                  <div class="group-actions">
                    <el-button link type="primary" icon="EditPen" @click.stop="openGroupDialog(g)"></el-button>
                    <el-button link type="danger" icon="Delete" @click.stop="deleteGroup(g.ID || g.id)"></el-button>
                  </div>
                </div>
              </template>
            </el-menu-item>
            <el-empty v-if="filteredGroupList.length === 0" description="暂无群组数据" :image-size="60"></el-empty>
          </el-menu>
        </el-card>
      </el-col>

      <!-- 右侧：规则列表 -->
      <el-col :span="18">
        <el-card class="modern-card main-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="header-title">
                <el-icon><Warning /></el-icon> 
                {{ activeGroupId === 'all' ? '所有群组告警规则' : activeGroupName ? `【${activeGroupName}】告警规则` : '告警规则' }}
              </span>
              <div class="right-actions">
                <el-button type="info" plain icon="Collection" @click="openStyleDialog">分类(Style)管理</el-button>
                <el-button type="primary" icon="Plus" :disabled="!activeGroupId" @click="openRuleDialog()">新增规则</el-button>
              </div>
            </div>
          </template>

          <!-- 搜索过滤 -->
          <div class="search-section">
            <el-form :inline="true" :model="query" class="search-form">
              <el-form-item label="名称">
                <el-input v-model="query.alert" placeholder="告警名称" clearable @keyup.enter="fetchRules" style="width: 150px" />
              </el-form-item>
              <el-form-item label="优先级">
                <el-select v-model="query.severity" clearable placeholder="优先级" @change="fetchRules" style="width: 120px">
                  <el-option label="Critical (严重)" value="critical" />
                  <el-option label="Warning (警告)" value="warning" />
                  <el-option label="Info (通知)" value="info" />
                </el-select>
              </el-form-item>
              <el-form-item label="状态">
                <el-select v-model="query.status" clearable placeholder="告警状态" @change="fetchRules" style="width: 120px">
                  <el-option label="触发中 (Firing)" value="firing" />
                  <el-option label="等待中 (Pending)" value="pending" />
                  <el-option label="未触发 (Inactive)" value="inactive" />
                </el-select>
              </el-form-item>
              <el-form-item label="规则分类">
                <el-select v-model="query.style" clearable placeholder="选择分类" @change="fetchRules" style="width: 150px">
                  <el-option v-for="s in styleList" :key="s.name" :label="s.name" :value="s.name" />
                </el-select>
              </el-form-item>
              <el-form-item label="启用状态">
                <el-select v-model="query.enabled" clearable placeholder="状态" @change="fetchRules" style="width: 120px;">
                  <el-option label="已启用" :value="1" />
                  <el-option label="已禁用" :value="0" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="Search" @click="fetchRules">查询</el-button>
                <el-button type="warning" icon="Refresh" @click="resetQuery">重置</el-button>
              </el-form-item>
            </el-form>
          </div>

          <!-- 规则表格 -->
          <el-table 
            :data="ruleList" 
            v-loading="loadingRules" 
            border 
            stripe 
            style="width: 100%"
            :header-cell-style="{ background: 'rgba(102, 126, 234, 0.1)', color: '#2c3e50', fontWeight: '600' }"
          >
            <el-table-column prop="alert" label="告警名称 (Alert)" min-width="150" show-overflow-tooltip />
            <el-table-column v-if="activeGroupId === 'all'" label="所属群组" min-width="140" show-overflow-tooltip>
              <template #default="{ row }">
                <el-tag size="small" type="warning" plain>{{ getGroupName(row.group_id) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="severity" label="优先级" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getSeverityType(row.severity)" effect="dark">{{ row.severity || '未知' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="style" label="分类" width="120" align="center">
              <template #default="{ row }">
                <el-tag type="info">{{ row.style || '暂无分类' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="expr" label="PromQL 表达式 (Expr)" min-width="250" show-overflow-tooltip>
              <template #default="{ row }">
                <code class="promql-code">{{ row.expr }}</code>
              </template>
            </el-table-column>
            <el-table-column prop="for_duration" label="持续时间(For)" width="120" align="center" />
            <el-table-column prop="status" label="告警状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="row.status === 'firing' ? 'danger' : row.status === 'pending' ? 'warning' : 'success'">
                  {{ row.status === 'firing' ? '触发中' : row.status === 'pending' ? '等待中' : '未触发' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="enabled" label="启用状态" width="100" align="center">
              <template #default="{ row }">
                <el-switch 
                  v-model="row.enabled" 
                  :active-value="1" 
                  :inactive-value="0" 
                  @change="toggleRuleStatus(row)"
                  style="--el-switch-on-color: #13ce66;"
                />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="220" align="center" fixed="right">
              <template #default="{ row }">
                <div class="operation-buttons">
                  <el-button link type="success" icon="CircleCheck" @click="checkRule(row)">检查</el-button>
                  <el-button link type="primary" icon="EditPen" @click="openRuleDialog(row)">编辑</el-button>
                  <el-button link type="danger" icon="Delete" @click="deleteRule(row.ID || row.id)">删除</el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>

          <!-- 分页 -->
          <div class="pagination-section">
            <el-pagination
              background
              layout="total, sizes, prev, pager, next, jumper"
              :total="totalRules"
              v-model:current-page="query.page"
              v-model:page-size="query.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              @current-change="fetchRules"
              @size-change="fetchRules"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 群组表单弹窗 -->
    <el-dialog :title="groupForm.id ? '编辑告警群组' : '新建告警群组'" v-model="groupDialogVisible" width="650px" class="modern-dialog" :append-to-body="true" :modal="false">
      <el-form :model="groupForm" :rules="groupRules" ref="groupFormRef" label-width="120px" class="rule-form">
        <!-- Section 1: 基本配置 -->
        <div class="form-section-title"><el-icon><Setting /></el-icon> 基本配置</div>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="群组名称" prop="group_name">
              <el-input v-model="groupForm.group_name" placeholder="例如：node-system-usage" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="数据源" prop="data_source_id">
              <el-select v-model="groupForm.data_source_id" placeholder="请选择数据源" style="width: 100%" :disabled="!!groupForm.id">
                <el-option v-for="ds in dataSources" :key="ds.id || ds.ID" :label="ds.name" :value="ds.id || ds.ID">
                  <div class="ds-option-content" style="margin-left: 0;">
                    <img :src="getDsTypeIcon(ds.type)" class="ds-icon" />
                    <span>{{ ds.name }}</span>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- Section 2: 群组定义 -->
        <div class="form-section-title"><el-icon><DataLine /></el-icon> 群组定义</div>
        <el-form-item label="业务线/Cluster" prop="cluster">
          <el-input v-model="groupForm.cluster" placeholder="例如: beijing-core" />
          <span class="help-text" style="display: block; margin-top: 8px;">结合 CMDB 模块下发的业务标识，自动转为实际的 {"cluster": "..."} 结构保存。</span>
        </el-form-item>

        <div class="form-section-title"><el-icon><DataLine /></el-icon> yaml定义</div>
        <div class="code-editor-container">
          <div class="code-editor-header">
            <span class="mac-dot close"></span>
            <span class="mac-dot minimize"></span>
            <span class="mac-dot maximize"></span>
            <span class="file-name">PrometheusRule.yaml</span>
          </div>
          <el-input v-model="groupForm.rule_content" type="textarea" class="promql-input custom-scrollbar" :rows="8" placeholder="（可选）PrometheusRule CRD 骨架..." />
        </div>
        
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="groupDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitGroup" :loading="submitting">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 规则表单弹窗 -->
    <el-dialog :title="ruleForm.id ? '编辑告警规则' : '新增告警规则'" v-model="ruleDialogVisible" width="760px" class="modern-dialog" :append-to-body="true" :modal="false">
      <el-form :model="ruleForm" :rules="ruleRules" ref="ruleFormRef" label-width="120px" class="rule-form">
        <!-- Section 1: 基本配置 -->
        <div class="form-section-title"><el-icon><Setting /></el-icon> 基本配置</div>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="所属群组" prop="group_id">
              <el-select v-model="ruleForm.group_id" placeholder="选择群组" style="width: 100%">
                <el-option v-for="g in groupList" :key="g.ID || g.id" :label="g.group_name" :value="parseInt(g.ID || g.id)" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="告警分类" prop="style">
              <el-select v-model="ruleForm.style" placeholder="选择分类" style="width: 100%">
                <el-option v-for="s in styleList" :key="s.name" :label="s.name" :value="s.name" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="告警名称(Alert)" prop="alert">
              <el-input v-model="ruleForm.alert" placeholder="例如: NodeCPUUsage" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="告警级别" prop="severity">
              <el-select v-model="ruleForm.severity" style="width: 100%;">
                <el-option label="Critical (严重)" value="critical" />
                <el-option label="Warning (警告)" value="warning" />
                <el-option label="Info (通知)" value="info" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- Section 2: 规则定义 -->
        <div class="form-section-title"><el-icon><DataLine /></el-icon> 规则定义</div>
        <el-form-item label="表达式(Expr)" prop="expr" class="expr-form-item">
          <el-input v-model="ruleForm.expr" type="textarea" :rows="4" class="promql-input" placeholder="PromQL, 比如: 100 - (avg by (instance) (rate(node_cpu_seconds_total{mode='idle'}[5m])) * 100) > 90" />
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="持续时间(For)" prop="for_duration">
              <el-input v-model="ruleForm.for_duration" placeholder="例如: 5m, 1h" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="是否启用">
              <el-switch v-model="ruleForm.enabled" :active-value="1" :inactive-value="0" />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- Section 3: 附加信息 -->
        <div class="form-section-title"><el-icon><Document /></el-icon> 附加信息</div>
        <el-form-item label="摘要(Summary)" prop="summary">
          <el-input v-model="ruleForm.summary" placeholder="告警摘要，短描述" />
        </el-form-item>
        
        <el-form-item label="描述(Desc)" prop="description">
          <el-input v-model="ruleForm.description" type="textarea" :rows="2" placeholder="详细的告警描述 (支持模板变量)" />
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="约束条件" prop="constraints">
              <el-input v-model="ruleForm.constraints" placeholder='JSON格式，如 {"env": "prod"}' />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="附加标签" prop="labels">
              <el-input v-model="ruleForm.labels" placeholder='JSON格式，如 {"team": "devops"}' />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="ruleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitRule" :loading="submitting">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 告警分类管理弹窗 -->
    <el-dialog title="管理告警分类 (Styles)" v-model="styleDialogVisible" width="650px" class="modern-dialog" :append-to-body="true" :modal="false">
      <div style="margin-bottom: 10px;">
        <el-button type="primary" size="small" icon="Plus" @click="openStyleForm()">新增分类</el-button>
      </div>
      <el-table :data="styleList" border size="small"
        :header-cell-style="{ background: 'rgba(102, 126, 234, 0.1)', color: '#2c3e50', fontWeight: '600' }">
        <el-table-column prop="name" label="分类名称" width="120" />
        <el-table-column prop="description" label="描述" min-width="150" />
        <el-table-column prop="CreatedAt" label="创建时间" width="160" align="center">
          <template #default="{ row }">
            {{ row.CreatedAt ? new Date(row.CreatedAt).toLocaleString() : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" align="center">
          <template #default="{ row }">
            <div class="operation-buttons">
               <el-button link type="primary" size="small" icon="EditPen" @click="openStyleForm(row)">编辑</el-button>
               <el-button link type="danger" size="small" icon="Delete" @click="deleteStyle(row.ID || row.id)">删除</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 新建/编辑分类内部弹窗 -->
    <el-dialog :title="styleForm.id ? '编辑分类' : '新增分类'" v-model="styleFormVisible" width="400px" class="modern-dialog" append-to-body>
      <el-form :model="styleForm" ref="styleFormRef" label-width="80px">
        <el-form-item label="分类名称" required>
          <el-input v-model="styleForm.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="styleForm.description" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="styleFormVisible = false">取消</el-button>
        <el-button type="primary" @click="submitStyle">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import yaml from 'js-yaml'
import {
  getDataSources,
  getAlertGroupList,
  getAlertGroupById,
  createAlertGroup,
  updateAlertGroup,
  deleteAlertGroup,
  getAlertRulesList,  getRuleListByGroup,  createAlertRule,
  updateAlertRule,
  deleteAlertRule,
  checkAlertRule,
  getAlertStyles,
  createAlertStyle,
  updateAlertStyle,
  deleteAlertStyle
} from '@/api/monitor'

// ---- 数据状态 ----
const dataSources = ref([])
const activeDataSourceId = ref(null)

const groupList = ref([])
const ruleList = ref([])
const styleList = ref([])
const filterCmdbGroup = ref('')

// 动态提取业务线(Cluster)作为筛选
const cmdbGroups = computed(() => {
  const clusters = new Set()
  groupList.value.forEach(g => {
    try {
      const l = JSON.parse(g.labels || '{}')
      if (l.cluster) clusters.add(l.cluster)
    } catch(e) {}
  })
  return Array.from(clusters).map(c => ({ id: c, name: c }))
})

const activeGroupId = ref('')
const activeGroupName = computed(() => {
  const g = groupList.value.find(item => (item.ID || item.id) + '' === activeGroupId.value)
  return g ? g.group_name : ''
})

const loadingGroups = ref(false)
const loadingRules = ref(false)
const submitting = ref(false)

const query = reactive({
  alert: '',
  severity: '',
  status: '',
  style: '',
  enabled: '',
  page: 1,
  pageSize: 10
})
const totalRules = ref(0) // 总数
    
// 根据下拉框筛选 CMDB 业务线和数据源 (群组级别)
const filteredGroupList = computed(() => {
  let list = groupList.value
  
  if (activeDataSourceId.value) {
    list = list.filter(g => g.data_source_id === activeDataSourceId.value)
  }

  if (!filterCmdbGroup.value) return list

  return list.filter(g => {
    try {
      const l = JSON.parse(g.labels || '{}')
      return l.cluster === filterCmdbGroup.value
    } catch(e) {
      return false
    }
  })
})

// 工具：取群组 cluster 名称用于渲染 Tag
const getClusterLabel = (labelsStr) => {
  try {
    const l = JSON.parse(labelsStr)
    return l.cluster || ''
  } catch(e) {
    return ''
  }
}

const parseId = (data) => data.ID || data.id

// ---- 初始化生命周期 ----
onMounted(() => {
  fetchDataSources()
  fetchStyles()
  fetchGroups()
})

// ---- HTTP/API 调用 ----
const fetchDataSources = async () => {
  try {
    const res = await getDataSources({ page: 1, pageSize: 100 })
    if (res.code === 200 || res.data) {
      console.log('数据源成功:')
      dataSources.value = res.data.data?.list || []
      if (dataSources.value.length > 0 && !activeDataSourceId.value) {
        activeDataSourceId.value = dataSources.value[0].id || dataSources.value[0].ID
      }
    }
  } catch (error) {
    console.error('获取数据源失败', error)
  }
}

const fetchStyles = async () => {
  try {
    const res = await getAlertStyles()
    if (res.code === 200 || res.data) {
      styleList.value = Array.isArray(res.data.data) ? res.data.data : []
    }
  } catch (error) {
    console.error('获取分类失败', error)
  }
}

const fetchGroups = async () => {
  loadingGroups.value = true
  try {
    const res = await getAlertGroupList({ page: 1, pageSize: 500 }) // 确保拿到完整群组列表以正常进行业务线分组与过滤
    if (res.code === 200 || res.data) {
      if (res.data.data?.list) {
        groupList.value = res.data.data.list
      } else {
        groupList.value = []
      }
      
      if (!activeGroupId.value) {
        activeGroupId.value = 'all'
        fetchRules()
      }
    }
  } catch (error) {
    console.error('获取群组失败', error)
  } finally {
    loadingGroups.value = false
  }
}

const fetchRules = async () => {
  loadingRules.value = true
  try {
    const params = { ...query }
    if (params.alert === '') delete params.alert
    if (params.severity === '') delete params.severity
    if (params.status === '') delete params.status
    if (params.style === '') delete params.style
    if (params.enabled === '') delete params.enabled
    
    let res = null
    if (activeGroupId.value && activeGroupId.value !== 'all') {
      res = await getRuleListByGroup(activeGroupId.value, params)
    } else {
      res = await getAlertRulesList(params)
    }

    if (res.code === 200 || res.data) {
      // 适配给出的接口返回格式 res.data.list
      if (res.data.data?.list) {
        let list = res.data.data.list
        
        // 当选择"全部规则"时，进一步根据所选的数据源在前端进行过滤（如果后端接口不支持按数据源查询）
        if (activeGroupId.value === 'all' && activeDataSourceId.value) {
          const dsGroupIds = new Set(groupList.value.filter(g => g.data_source_id === activeDataSourceId.value).map(g => (g.ID || g.id) + ''))
          list = list.filter(r => dsGroupIds.has((r.group_id || '') + ''))
        }
        
        ruleList.value = list
        totalRules.value = res.data.total || list.length
      } else {
        ruleList.value = []
        totalRules.value = 0
      }
    }
  } catch (error) {
    console.error('获取规则失败', error)
  } finally {
    loadingRules.value = false
  }
}

// 侧边栏菜单切换组
const handleGroupSelect = (index) => {
  activeGroupId.value = index
  resetQuery()
}

const resetQuery = () => {
  query.alert = ''
  query.severity = ''
  query.status = ''
  query.style = ''
  query.enabled = ''
  query.page = 1
  fetchRules()
}

// ---- 群组管理 ----
const groupDialogVisible = ref(false)
const groupFormRef = ref(null)
const groupForm = reactive({ id: null, group_name: '', data_source_id: 1, cluster: '', rule_content: '' })

const validateYaml = (rule, value, callback) => {
  if (!value) return callback()
  try {
    yaml.load(value)
    callback()
  } catch (e) {
    callback(new Error(`YAML 格式错误: ${e.message.split('\\n')[0]}`))
  }
}

const groupRules = {
  group_name: [{ required: true, message: '请输入群组名称', trigger: 'blur' }],
  cluster: [{ required: true, message: '请输入业务线/Cluster标识', trigger: 'blur' }],
  rule_content: [{ validator: validateYaml, trigger: 'blur' }]
}

const openGroupDialog = (row = null) => {
  if (row) {
    groupForm.id = parseId(row)
    groupForm.group_name = row.group_name
    groupForm.data_source_id = row.data_source_id || 1
    
    let clusterVal = ''
    try {
      const l = JSON.parse(row.labels || '{}')
      clusterVal = l.cluster || ''
    } catch(e) {}
    groupForm.cluster = clusterVal
    
    groupForm.rule_content = row.rule_content || ''
  } else {
    groupForm.id = null
    groupForm.group_name = ''
    groupForm.data_source_id = 1
    groupForm.cluster = 'default'
    groupForm.rule_content = 'apiVersion: monitoring.coreos.com/v1\nkind: PrometheusRule\nmetadata:\n  name: default-rules\nspec:\n  groups:\n    - name: default.rules\n      rules: []'
  }
  groupDialogVisible.value = true
}

const submitGroup = async () => {
  if (!groupFormRef.value) return
  await groupFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const action = groupForm.id ? updateAlertGroup : createAlertGroup
        const payload = { ...groupForm }
        payload.labels = JSON.stringify({ cluster: groupForm.cluster })
        delete payload.cluster
        if (groupForm.id) payload.id = groupForm.id
        
        await action(payload)
        ElMessage.success(groupForm.id ? '编辑群组成功' : '新增群组成功')
        groupDialogVisible.value = false
        fetchGroups()
      } catch (error) {
        console.error(error)
        ElMessage.error('保存群组失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

const deleteGroup = (id) => {
  ElMessageBox.confirm('确认删除该群组及其下属所有规则吗?', '警告', { type: 'warning' }).then(async () => {
    try {
      await deleteAlertGroup(id)
      ElMessage.success('删除组成功')
      if (activeGroupId.value == id) {
        activeGroupId.value = ''
        ruleList.value = []
      }
      fetchGroups()
    } catch(err) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

// ---- 规则管理 ----
const ruleDialogVisible = ref(false)
const ruleFormRef = ref(null)
const ruleForm = reactive({
  id: null, group_id: null, alert: '', expr: '', for_duration: '5m',
  severity: 'warning', summary: '', description: '', constraints: '{}', labels: '{}', style: '', enabled: 1
})

const validateExpr = (rule, value, callback) => {
  if (!value) return callback(new Error('请输入告警表达式'))
  const stack = []
  const pairs = { ')': '(', '}': '{', ']': '[' }
  let inString = false
  let stringChar = ''
  
  for (let i = 0; i < value.length; i++) {
    const char = value[i]
    if (!inString && (char === '"' || char === "'" || char === '`')) {
      inString = true
      stringChar = char
      continue
    } else if (inString && char === stringChar && value[i-1] !== '\\') {
      inString = false
      continue
    }
    if (inString) continue
    
    if (['(', '{', '['].includes(char)) { stack.push(char) } 
    else if ([')', '}', ']'].includes(char)) {
      if (stack.length === 0 || stack.pop() !== pairs[char]) return callback(new Error(`表达式格式有误: 括号 '${char}' 不匹配`))
    }
  }
  
  if (inString) return callback(new Error('表达式格式有误: 存在未闭合的引号'))
  if (stack.length > 0) return callback(new Error(`表达式格式有误: 存在未闭合的括号`))
  
  callback()
}

const ruleRules = {
  alert: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  group_id: [{ required: true, message: '请选择所属群组', trigger: 'change' }],
  expr: [{ required: true, validator: validateExpr, trigger: 'blur' }],
  style: [{ required: true, message: '请选择分类', trigger: 'change' }]
}

const openRuleDialog = (row = null) => {
  if (!activeGroupId.value) return ElMessage.warning('请先选择一个告警群组')
  if (row) {
    Object.keys(ruleForm).forEach(k => { ruleForm[k] = row[k] !== undefined ? row[k] : ruleForm[k] })
    ruleForm.id = parseId(row)
    ruleForm.group_id = parseInt(row.group_id || activeGroupId.value)
  } else {
    ruleForm.id = null
    ruleForm.group_id = activeGroupId.value === 'all' ? null : parseInt(activeGroupId.value)
    ruleForm.alert = ''
    ruleForm.expr = ''
    ruleForm.for_duration = '5m'
    ruleForm.severity = 'warning'
    ruleForm.summary = ''
    ruleForm.description = ''
    ruleForm.constraints = '{"env": "prod"}'
    ruleForm.labels = '{"team": "devops"}'
    ruleForm.style = ''
    ruleForm.enabled = 1
  }
  ruleDialogVisible.value = true
}

const submitRule = async () => {
  if (!ruleFormRef.value) return
  await ruleFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const action = ruleForm.id ? updateAlertRule : createAlertRule
        const payload = { ...ruleForm }
        if (ruleForm.id) payload.id = ruleForm.id
        
        await action(payload)
        ElMessage.success(ruleForm.id ? '编辑规则成功' : '新增规则成功')
        ruleDialogVisible.value = false
        fetchRules() // 刷新列表
      } catch (error) {
        ElMessage.error('保存规则失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

const toggleRuleStatus = async (row) => {
  try {
    const payload = { ...row }
    payload.id = parseId(row)
    await updateAlertRule(payload)
    ElMessage.success(`告警规则已${row.enabled ? '启用' : '禁用'}`)
  } catch (error) {
    row.enabled = row.enabled ? 0 : 1 // 恢复状态
    ElMessage.error('状态修改失败')
  }
}

const checkRule = async (row) => {
  if (!row.expr) return ElMessage.warning('规则表达式为空，无法检查')
  try {
    const group = groupList.value.find(item => parseInt(item.ID || item.id) === parseInt(row.group_id))
    const dataSourceId = group ? group.data_source_id : activeDataSourceId.value
    
    if (!dataSourceId) return ElMessage.warning('未找到对应的数据源，无法检查')

    const res = await checkAlertRule({
      data_source_id: parseInt(dataSourceId),
      expr: row.expr
    })
    
    if (res.code === 200 || (res.data && res.data.code === 200)) {
       ElMessage.success('检查通过，规则生效中')
    } else {
       ElMessage.warning('检查返回失败: ' + (res.data?.message || '请确认表达式'))
    }
  } catch (error) {
    ElMessage.error(error.message || '检查发生错误，可能是由于表达式不合法')
  }
}

const deleteRule = (id) => {
  ElMessageBox.confirm('确定删除此告警规则吗？', '警告', { type: 'warning' }).then(async () => {
    try {
      await deleteAlertRule(id)
      ElMessage.success('删除成功')
      fetchRules()
    } catch(err) {
      ElMessage.error('删除规则失败')
    }
  }).catch(() => {})
}

// ---- 分类管理 ----
const styleDialogVisible = ref(false)
const styleFormVisible = ref(false)
const styleFormRef = ref(null)
const styleForm = reactive({ id: null, name: '', description: '' })

const openStyleDialog = () => {
  styleDialogVisible.value = true
  fetchStyles()
}

const openStyleForm = (row = null) => {
  if (row) {
    styleForm.id = parseId(row)
    styleForm.name = row.name
    styleForm.description = row.description
  } else {
    styleForm.id = null
    styleForm.name = ''
    styleForm.description = ''
  }
  styleFormVisible.value = true
}

const submitStyle = async () => {
  if (!styleForm.name) return ElMessage.warning('分类名称不能为空')
  try {
    const action = styleForm.id ? updateAlertStyle : createAlertStyle
    const payload = { ...styleForm }
    if (styleForm.id) payload.id = styleForm.id
    
    await action(payload)
    ElMessage.success('保存分类成功')
    styleFormVisible.value = false
    fetchStyles()
  } catch (error) {
    ElMessage.error('保存分类失败')
  }
}

const deleteStyle = (id) => {
  ElMessageBox.confirm('确定删除该分类吗？', '警告', { type: 'warning' }).then(async () => {
    try {
      await deleteAlertStyle(id)
      ElMessage.success('删除成功')
      fetchStyles()
    } catch (err) { }
  }).catch(() => {})
}

const getDataSourceIcon = (id) => {
  const ds = dataSources.value.find(item => (item.id || item.ID) === id)
  return ds ? getDsTypeIcon(ds.type) : ''
}

const getDsTypeIcon = (type) => {
  const t = (type || '').toLowerCase()
  if (t.includes('prometheus')) return 'https://upload.wikimedia.org/wikipedia/commons/3/38/Prometheus_software_logo.svg'
  if (t.includes('Zabbix')) return 'https://upload.wikimedia.org/wikipedia/commons/3/33/Zabbix_logo.svg'
  if (t.includes('Loki')) return 'https://upload.wikimedia.org/wikipedia/commons/6/60/Grafana_logo.svg'
  return 'https://cdn-icons-png.flaticon.com/512/3168/3168610.png' // default database/monitor
}

const getDeployIcon = (method) => {
  const m = (method || '').toLowerCase()
  if (m.includes('Docker')) return 'https://upload.wikimedia.org/wikipedia/commons/4/4e/Docker_%28container_engine%29_logo.svg'
  if (m.includes('kubenetes') || m.includes('Kubernetes')) return 'https://upload.wikimedia.org/wikipedia/commons/3/39/Kubernetes_logo_without_workmark.svg'
  return ''
}

const handleDataSourceChange = () => {
  activeGroupId.value = 'all'
  fetchRules()
}

// ---- 辅助函数 ----
const getSeverityType = (severity) => {
  switch (severity?.toLowerCase()) {
    case 'critical': return 'danger'
    case 'warning': return 'warning'
    case 'info': return 'info'
    default: return ''
  }
}

const getGroupName = (id) => {
  const g = groupList.value.find(item => (item.ID || item.id) + '' === id + '')
  return g ? g.group_name : `ID: ${id}`
}
</script>

<style scoped>
.alarm-rules-management {
  padding: 20px;
  min-height: calc(100vh - 120px);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.modern-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  flex-direction: column;
}

.menu-card {
  height: 100%;
}

.main-card {
  height: 100%;
}

:deep(.modern-card .el-card__body) {
  flex: 1;
  overflow: auto;
  padding: 15px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.transparent-bg {
  background-color: transparent !important;
}

.group-menu {
  border-right: none;
}
.menu-item-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}
.group-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: flex;
  align-items: center;
}

.cluster-tag {
  margin-left: 6px;
  transform: scale(0.85);
  transform-origin: left center;
}

.group-actions {
  display: none;
  align-items: center;
  gap: 4px;
}
.el-menu-item:hover .group-actions {
  display: flex !important;
}

/* 搜索表单部分 */
.search-section {
  margin-bottom: 16px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form .el-form-item {
  margin-bottom: 12px;
  margin-right: 16px;
}

.search-form .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

/* 表格美化 */
:deep(.el-table) {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

:deep(.el-table__header th) {
  border: none;
}

:deep(.el-table__body tr:hover > td) {
  background-color: rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-table td) {
  border: none;
}

:deep(.el-table::before) {
  display: none;
}
:deep(.el-table--border::after) {
  display: none;
}

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
}

/* 按钮美化 */
:deep(.el-button--primary) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.el-button--primary:hover) {
  background: linear-gradient(135deg, #5a6fd8, #6a4190);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

:deep(.el-button--warning) {
  background: linear-gradient(135deg, #f39c12, #e67e22);
  border: none;
  border-radius: 8px;
  color: white;
  transition: all 0.3s ease;
}

:deep(.el-button--warning:hover) {
  background: linear-gradient(135deg, #e67e22, #d35400);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(243, 156, 18, 0.4);
}

.promql-code {
  background-color: rgba(103, 126, 234, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
  color: #c7254e;
  font-family: Menlo, Monaco, Consolas, "Courier New", monospace;
  word-break: break-all;
  white-space: normal;
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

/* 对话框美化 */
:deep(.modern-dialog .el-dialog) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(31, 38, 135, 0.37);
}

:deep(.modern-dialog .el-dialog__header) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
  border-radius: 16px 16px 0 0;
  padding: 20px 24px 16px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.2);
}

:deep(.modern-dialog .el-dialog__title) {
  color: #2c3e50;
  font-weight: 600;
  font-size: 18px;
}

:deep(.modern-dialog .el-dialog__body) {
  padding: 24px;
}

:deep(.modern-dialog .el-dialog__footer) {
  padding: 16px 24px 24px;
  background: rgba(248, 249, 250, 0.8);
  border-radius: 0 0 16px 16px;
}

.help-text {
  font-size: 12px;
  color: #999;
  line-height: 1.2;
  margin-top: 5px;
}

/* 告警表单美化 */
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

.promql-input :deep(textarea) {
  font-family: Menlo, Monaco, Consolas, "Courier New", monospace;
  background-color: #1e1e1f;
  color: #d4d4d4;
  border: 1px solid #333;
  padding: 12px;
  font-size: 14px;
  line-height: 1.6;
  border-radius: 8px;
  box-shadow: inset 0 2px 4px rgba(0,0,0,0.2);
  transition: border 0.3s;
}

.promql-input :deep(textarea):focus {
  border-color: #667eea;
  box-shadow: inset 0 2px 4px rgba(0,0,0,0.2), 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.expr-form-item :deep(.el-form-item__label) {
  font-weight: 600;
  color: #e74c3c;
}

/* 数据源下拉框样式 */
.ds-selector {
  margin-bottom: 15px !important;
  background: rgba(255, 255, 255, 0.4);
  border-radius: 10px;
  padding: 4px;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

.ds-select {
  --el-select-input-focus-border-color: transparent;
}

.ds-select :deep(.el-input__wrapper) {
  box-shadow: none !important;
  background-color: transparent;
}

.ds-icon-small {
  width: 20px;
  height: 20px;
  margin-right: 6px;
  object-fit: contain;
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.1));
}

.ds-option {
  height: auto !important;
  padding: 10px 12px;
}

.ds-option-content {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}

.ds-icon {
  width: 26px;
  height: 26px;
  object-fit: contain;
  filter: drop-shadow(0 2px 5px rgba(0,0,0,0.15));
}

.ds-name {
  font-weight: 600;
  color: #2c3e50;
  font-size: 14px;
  flex: 1;
}

.group-colored-icon {
  width: 22px;
  height: 22px;
  margin-right: 12px;
  object-fit: contain;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.15));
  transition: transform 0.3s ease;
}

.el-menu-item:hover .group-colored-icon {
  transform: scale(1.15) translateY(-2px);
}

.ds-deploy-tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  border-radius: 20px;
  padding: 0 10px;
  height: 24px;
  border: none;
  background: rgba(102, 126, 234, 0.1);
  color: #5a6fd8;
  font-weight: 500;
}

.deploy-icon {
  width: 14px;
  height: 14px;
  object-fit: contain;
}

/* 代码编辑器包装器 */
.code-editor-container {
  background: #1e1e1f;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #333;
  margin-top: 10px;
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
  font-family: Menlo, Monaco, Consolas, "Courier New", monospace;
  opacity: 0.8;
}

.code-editor-container .promql-input :deep(textarea) {
  border: none;
  border-radius: 0 0 8px 8px;
  resize: vertical;
  padding-top: 12px;
  background: #1e1e1f !important;
  box-shadow: none;
}
.code-editor-container .promql-input :deep(textarea):focus {
  box-shadow: none;
}
</style>
