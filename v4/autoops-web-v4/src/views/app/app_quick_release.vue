<template>
  <div class="quick-release-management">
    <div class="glass-card main-card">
      <!-- 卡片标题 -->
      <div class="card-header">
        <h1 class="gradient-title">快速发布</h1>
        <el-button type="primary" size="small" v-authority="['app:quick-release:add']"  @click="handleCreateRelease" class="modern-btn">
          <el-icon><Plus /></el-icon>
          新建发布计划
        </el-button>
      </div>

      <!-- 搜索区域 -->
      <div class="search-section">
        <el-form :model="queryParams" :inline="true" class="search-form">
          <el-form-item label="发布标题" prop="title">
            <el-input
              placeholder="请输入发布标题"
              clearable
              size="small"
              class="modern-input"
              v-model="queryParams.title">
            </el-input>
          </el-form-item>
          <el-form-item label="环境" prop="environment">
            <el-select
              v-model="queryParams.environment"
              placeholder="请选择环境"
              size="small"
              style="width: 150px;"
              class="modern-select"
              clearable>
              <el-option label="测试环境" value="test" />
              <el-option label="预发环境" value="staging" />
              <el-option label="生产环境" value="prod" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select
              v-model="queryParams.status"
              placeholder="发布状态"
              size="small"
              style="width: 150px;"
              class="modern-select"
              clearable>
              <el-option
                v-for="item in statusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"/>
            </el-select>
          </el-form-item>
          <el-form-item label="业务组" prop="business_group_id">
            <el-select
              v-model="queryParams.business_group_id"
              placeholder="请选择业务组"
              size="small"
              style="width: 150px;"
              class="modern-select"
              clearable>
              <el-option
                v-for="group in businessGroups"
                :key="group.id"
                :label="group.name"
                :value="group.id"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              size="small"
              @click="handleQuery"
              class="modern-btn primary-btn">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="resetQuery"
              class="modern-btn reset-btn">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 发布列表 -->
      <div class="table-section">
        <el-table
          :data="releaseList"
          v-loading="loading"
          class="modern-table"
          stripe>
          <el-table-column prop="title" label="发布标题" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="table-cell-with-icon">
                <img src="@/assets/image/renwu.svg" width="16" height="16" alt="任务" />
                <span>{{ row.title }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="applications" label="发布服务">
            <template #default="{ row }">
              <div class="services-cell">
                <template v-if="row.tasks && row.tasks.length > 0">
                  <div
                    v-for="(task, index) in row.tasks"
                    :key="`${task.app_id}-${index}`"
                    class="service-item">
                    <span class="service-name">{{ task.app_name }}</span>
                    <div class="environment-wrapper">
                      <img src="@/assets/image/环境.svg" width="12" height="12" alt="环境" />
                      <span class="environment-tag" :class="getEnvironmentClass(task.environment)">
                        {{ task.environment }}
                      </span>
                    </div>
                  </div>
                </template>
                <template v-else-if="row.applications && row.applications.length > 0">
                  <div
                    v-for="(app, index) in row.applications"
                    :key="`${app.app_id}-${index}`"
                    class="service-item">
                    <span class="service-name">{{ app.app_name || `服务${app.app_id}` }}</span>
                    <div class="environment-wrapper" v-if="app.environment">
                      <img src="@/assets/image/环境.svg" width="12" height="12" alt="环境" />
                      <span class="environment-tag" :class="getEnvironmentClass(app.environment)">
                        {{ app.environment }}
                      </span>
                    </div>
                  </div>
                </template>
                <span v-else class="no-service">-</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="duration" label="耗时">
            <template #default="{ row }">
              <div class="table-cell-with-icon">
                <img src="@/assets/image/shijian.svg" width="16" height="16" alt="时间" />
                <span>{{ formatDuration(row.duration) }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="creator_name" label="创建人">
            <template #default="{ row }">
              <span>{{ row.creator_username || row.creator_name || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间">
            <template #default="{ row }">
              <span>{{ formatDateTime(row.created_at) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="操作" fixed="right" width="180">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-tooltip v-if="!executedReleases.has(row.id)" content="服务上线" placement="top">
                  <el-button :icon="VideoPlay" v-authority="['app:quick-release:start']" size="small" type="primary" @click="handleExecute(row)">服务上线</el-button>
                </el-tooltip>
                <el-tooltip :content="row.status === 2 ? '发布中的任务不能删除' : '删除任务'" placement="top">
                  <el-button
                    :icon="Delete"
                    size="small"
                    type="danger"
                    v-authority="['app:quick-release:delete']"
                    :disabled="row.status === 2"
                    @click="handleDelete(row)">删除</el-button>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页组件 -->
      <div class="pagination-section">
        <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="pagination.page"
          :page-sizes="[10, 20, 50, 100]"
          :page-size="pagination.pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          class="modern-pagination"
        />
      </div>
    </div>

    <!-- 新建发布对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="新建快速发布"
      width="800px"
      class="modern-dialog">
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="120px" class="modern-form">
        <el-form-item label="发布标题" prop="title" required>
          <el-input
            v-model="createForm.title"
            placeholder="请输入发布标题"
            clearable
          />
        </el-form-item>
        <div class="form-row">
          <el-form-item label="业务组" prop="business_group_id" required class="form-item-half">
            <el-cascader
              v-model="createForm.business_group_id"
              :options="cascaderBusinessGroups"
              placeholder="请选择业务组"
              style="width: 100%"
              @change="handleGroupChange"
              :props="{
                expandTrigger: 'hover',
                checkStrictly: true,
                emitPath: false,
                leaf: 'leaf'
              }"
            />
          </el-form-item>
          <el-form-item label="业务部门" prop="business_dept_id" required class="form-item-half">
            <el-select
              v-model="createForm.business_dept_id"
              placeholder="请选择业务部门"
              style="width: 100%"
              size="small"
              filterable
              clearable
              @change="handleDeptChange">
              <el-option
                v-for="dept in businessDepts"
                :key="dept.id"
                :label="dept.deptName"
                :value="dept.id" />
            </el-select>
          </el-form-item>
        </div>
        <el-form-item label="应用配置" required>
          <div class="applications-section">
            <div class="applications-header">
              <el-button type="primary" size="small" @click="showAppSelectionDialog" icon="Plus">
                添加应用
              </el-button>
              <span class="app-count" v-if="createForm.applications.length > 0">
                已选择 {{ createForm.applications.length }} 个应用
              </span>
            </div>

            <!-- 已选择的应用列表 -->
            <div class="selected-apps-list" v-if="createForm.applications.length > 0">
              <div
                v-for="(app, index) in createForm.applications"
                :key="`${app.app_id}-${index}`"
                class="selected-app-item">
                <div class="app-info">
                  <span class="app-name">{{ app.app_name || `应用${app.app_id}` }}</span>
                  <el-tag size="small" type="info">{{ app.business_group_name }}</el-tag>
                </div>
                <div class="app-environment">
                  <el-select
                    v-model="app.environment"
                    placeholder="选择环境"
                    size="small"
                    style="width: 120px;"
                    @change="handleEnvironmentChange(app)">
                    <el-option
                      v-for="env in app.available_environments || []"
                      :key="env.value"
                      :label="env.label"
                      :value="env.value" />
                  </el-select>
                </div>
                <div class="app-params">
                  <el-button
                    size="small"
                    :type="isAppParametersConfigured(app) ? 'success' : 'warning'"
                    plain
                    :disabled="!app.environment"
                    @click="handleConfigParams(app, index)">
                    <el-icon><Setting /></el-icon>
                    {{ isAppParametersConfigured(app) ? '已配置参数' : '配置参数' }}
                  </el-button>
                </div>
                <div class="app-actions">
                  <el-button
                    size="small"
                    type="danger"
                    icon="Delete"
                    @click="removeApplication(index)">
                    移除
                  </el-button>
                </div>
              </div>
            </div>

            <!-- 空状态提示 -->
            <div class="empty-apps" v-if="createForm.applications.length === 0">
              <el-empty :image-size="100" description="请添加要发布的应用">
                <template #image>
                  <el-icon size="100" color="#ccc"><Document /></el-icon>
                </template>
              </el-empty>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="发布描述" prop="description" required>
          <el-input v-model="createForm.description" type="textarea" rows="3" placeholder="请输入发布描述" />
        </el-form-item>

        <el-form-item label="是否通知">
          <el-switch v-model="createForm.notice" />
        </el-form-item>

        <el-form-item v-if="createForm.notice" label="通知机器人" prop="robot_id">
          <el-select
            v-model="createForm.robot_id"
            placeholder="请选择通知机器人"
            style="width: 100%"
            :loading="notifyRobotLoading"
            clearable
          >
            <el-option
              v-for="robot in notifyRobotList"
              :key="robot.id"
              :label="`${robot.name}${robot.status === 1 ? '(启用)' : '(禁用)'}`"
              :value="robot.id"
              :disabled="robot.status === 0"
            >
              <div style="display: flex; justify-content: space-between; align-items: center; width: 100%;">
                <span>{{ robot.name }}</span>
                <el-tag :type="robot.status === 1 ? 'success' : 'info'" size="small" style="font-size: 12px;">
                  {{ robot.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </div>
            </el-option>
          </el-select>
          <div v-if="notifyRobotList.length === 0 && !notifyRobotLoading" style="margin-top: 8px; color: #F56C6C; font-size: 12px;">
            暂无可用的通知机器人，请先在系统设置中配置
          </div>
          <div v-else-if="notifyRobotList.filter(r => r.status === 1).length === 0 && !notifyRobotLoading" style="margin-top: 8px; color: #E6A23C; font-size: 12px;">
            没有启用的通知机器人，请先启用相关机器人
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateSubmit" :loading="creating">创建发布</el-button>
      </template>
    </el-dialog>

    <!-- 应用选择弹框 -->
    <el-dialog
      v-model="appSelectionDialogVisible"
      title="选择应用"
      width="1000px"
      class="modern-dialog app-selection-dialog">
      <div class="host-selection-layout">
        <!-- 左侧业务树 -->
        <div class="tree-section">
          <div class="tree-header">
            <h4>业务树</h4>
          </div>
          <el-tree
            :data="businessTreeData"
            :props="{
              label: 'name',
              children: 'children'
            }"
            node-key="id"
            highlight-current
            @node-click="handleBusinessTreeClick"
          >
            <template #default="{ node, data }">
              <span class="tree-node">
                <el-icon v-if="data.children?.length > 0">
                  <component :is="node.expanded ? 'FolderOpened' : 'Folder'" />
                </el-icon>
                <el-icon v-else><Document /></el-icon>
                {{ node.label }}
              </span>
            </template>
          </el-tree>
        </div>

        <!-- 右侧应用表格 -->
        <div class="services-section">
          <div class="services-header">
            <h4>应用列表</h4>
            <span v-if="currentBusinessGroup" class="selected-group">{{ currentBusinessGroup.name }}</span>
          </div>

          <el-table
            :data="businessGroupApps"
            border
            style="width: 100%"
            @selection-change="handleAppSelectionChange"
            ref="appTable"
            v-loading="loadingApps">
            <el-table-column
              type="selection"
              width="55"
              :selectable="(row) => row.status === 2"
            />
            <el-table-column prop="name" label="应用名称" min-width="150">
              <template #default="{ row }">
                <div class="app-name-cell">
                  <el-icon><Document /></el-icon>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="code" label="应用代码" width="120">
              <template #default="{ row }">
                <el-tag type="info" size="small">{{ row.code }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 2 ? 'success' : 'danger'" size="small">
                  {{ row.status === 2 ? '已激活' : '未激活' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="jenkins_envs" label="支持环境" min-width="200">
              <template #default="{ row }">
                <div class="environments-cell">
                  <el-tag
                    v-for="env in row.jenkins_envs"
                    :key="env.env_name"
                    size="small"
                    type="warning"
                    style="margin-right: 4px; margin-bottom: 2px;">
                    {{ getEnvironmentText(env.env_name) }}
                  </el-tag>
                  <span v-if="!row.jenkins_envs || row.jenkins_envs.length === 0" class="no-env">
                    暂无环境配置
                  </span>
                </div>
              </template>
            </el-table-column>
          </el-table>

          <div class="empty-state" v-if="!currentBusinessGroup">
            <el-empty :image-size="100" description="请选择左侧业务组查看应用">
              <template #image>
                <el-icon size="100" color="#ccc"><Document /></el-icon>
              </template>
            </el-empty>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="appSelectionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAppSelection" :disabled="selectedAppsTemp.length === 0">
          确认选择 ({{ selectedAppsTemp.length }})
        </el-button>
      </template>
    </el-dialog>

    <!-- Jenkins任务参数配置对话框 -->
    <el-dialog
      v-model="jenkinsParamsDialogVisible"
      :title="`配置参数: ${currentAppConfig.app_name} - ${currentAppConfig.environment}`"
      width="800px"
      class="modern-dialog jenkins-params-dialog"
      @close="closeParamsDialog"
    >
      <el-card shadow="hover" v-loading="jenkinsParamsLoading">
        <div v-if="jenkinsJobParameters.length === 0" class="no-params">
          <el-empty description="该任务无需参数，可直接发布" />
        </div>

        <el-form v-else label-width="150px" class="params-form">
          <div v-for="param in jenkinsJobParameters" :key="param.name" class="param-item">
            <!-- String 类型 -->
            <el-form-item
              v-if="param.type === 'String'"
              :label="param.name"
              :required="param.required"
            >
              <el-input
                v-model="jenkinsParamValues[param.name]"
                :placeholder="param.description || `请输入${param.name}`"
                clearable
              />
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>

            <!-- Choice 类型(单选) -->
            <el-form-item
              v-else-if="param.type === 'Choice' && param.choices && param.choices.length > 0"
              :label="param.name"
              :required="param.required"
            >
              <el-select
                v-model="jenkinsParamValues[param.name]"
                :placeholder="param.description || `请选择${param.name}`"
                clearable
                style="width: 100%"
              >
                <el-option
                  v-for="choice in param.choices"
                  :key="choice"
                  :label="choice"
                  :value="choice"
                />
              </el-select>
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>

            <!-- Boolean 类型 -->
            <el-form-item
              v-else-if="param.type === 'Boolean'"
              :label="param.name"
            >
              <el-switch
                v-model="jenkinsParamValues[param.name]"
                :active-text="param.description || '是'"
                inactive-text="否"
              />
            </el-form-item>

            <!-- Text 类型 -->
            <el-form-item
              v-else-if="param.type === 'Text'"
              :label="param.name"
              :required="param.required"
            >
              <el-input
                v-model="jenkinsParamValues[param.name]"
                type="textarea"
                :rows="5"
                :placeholder="param.description || `请输入${param.name}`"
              />
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>

            <!-- Password 类型 -->
            <el-form-item
              v-else-if="param.type === 'Password'"
              :label="param.name"
              :required="param.required"
            >
              <el-input
                v-model="jenkinsParamValues[param.name]"
                type="password"
                :placeholder="param.description || `请输入${param.name}`"
                show-password
                clearable
              />
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>

            <!-- PT_CHECKBOX 类型(多选框) -->
            <el-form-item
              v-else-if="param.type === 'PT_CHECKBOX' && param.choices && param.choices.length > 0"
              :label="param.name"
              :required="param.required"
            >
              <el-checkbox-group v-model="jenkinsParamValues[param.name]">
                <el-checkbox
                  v-for="choice in param.choices"
                  :key="choice"
                  :label="choice"
                  :value="choice"
                >
                  {{ choice }}
                </el-checkbox>
              </el-checkbox-group>
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>

            <!-- 默认类型(Text输入) -->
            <el-form-item
              v-else
              :label="param.name"
              :required="param.required"
            >
              <el-input
                v-model="jenkinsParamValues[param.name]"
                :placeholder="param.description || `请输入${param.name}`"
                clearable
              />
              <span v-if="param.description" class="param-description">{{ param.description }}</span>
            </el-form-item>
          </div>
        </el-form>
      </el-card>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeParamsDialog">取消</el-button>
          <el-button type="primary" @click="saveJenkinsParams">
            <el-icon><Check /></el-icon>
            保存配置
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Search, Refresh, View, Document, VideoPlay, Delete, Folder, FolderOpened, Setting, Check
} from '@element-plus/icons-vue'
import appApi from '@/api/app'
import cmdbApi from '@/api/cmdb'
import systemApi from '@/api/system'
import monitorApi from '@/api/monitor'

defineOptions({ name: 'AppQuickRelease' })

const router = useRouter()

// 响应式数据
const loading = ref(false)
const releaseList = ref([])
const businessGroups = ref([]) // 业务组列表
const businessDepts = ref([]) // 业务部门列表
const serviceTree = ref([]) // 服务树数据
const loadingApps = ref(false)
const executedReleases = ref(new Set()) // 已点击发布的任务ID集合

// 应用选择弹框相关
const appSelectionDialogVisible = ref(false)
const businessTreeData = ref([]) // 业务树数据
const currentBusinessGroup = ref(null) // 当前选中的业务组
const businessGroupApps = ref([]) // 当前业务组的应用列表
const selectedAppsTemp = ref([]) // 临时选中的应用
const appTable = ref() // 应用表格引用
const creating = ref(false)

// Jenkins任务参数配置
const jenkinsParamsDialogVisible = ref(false)
const jenkinsParamsLoading = ref(false)
const jenkinsJobParameters = ref([])
const jenkinsParamValues = ref({})
const currentAppConfig = ref({}) // 当前配置参数的应用信息
const currentAppIndex = ref(-1) // 当前配置参数的应用索引

// 查询参数
const queryParams = reactive({
  title: '',
  environment: '',
  status: null,
  business_group_id: null,
  page: 1,
  pageSize: 10
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 状态选项
const statusOptions = [
  { label: '待发布', value: 1 },
  { label: '发布中', value: 2 },
  { label: '发布成功', value: 3 },
  { label: '发布失败', value: 4 }
]

// 新建发布表单
const createDialogVisible = ref(false)
const createFormRef = ref()
const createForm = reactive({
  title: '', // 发布标题
  business_group_id: null, // 业务组ID
  business_dept_id: null,  // 业务部门ID
  description: '', // 发布描述
  applications: [], // 应用数组 [{app_id, app_name, business_group_name, environment, available_environments}]
  notice: false, // 是否启用通知
  robot_id: null // 通知机器人ID
})

// 通知机器人相关
const notifyRobotList = ref([]) // 通知机器人列表
const notifyRobotLoading = ref(false) // 加载状态

const createRules = {
  title: [{ required: true, message: '请输入发布标题', trigger: 'blur' }],
  business_group_id: [{ required: true, message: '请选择业务组', trigger: 'change' }],
  business_dept_id: [{ required: true, message: '请选择业务部门', trigger: 'change' }],
  description: [{ required: true, message: '请输入发布描述', trigger: 'blur' }],
  applications: [
    {
      validator: (rule, value, callback) => {
        if (!value || value.length === 0) {
          callback(new Error('请至少添加一个应用'))
        } else if (value.some(app => !app.environment)) {
          callback(new Error('请为所有应用选择环境'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  robot_id: [
    {
      validator: (rule, value, callback) => {
        if (createForm.notice && !value) {
          callback(new Error('启用通知时必须选择通知机器人'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ]
}

// 工具函数
const getEnvironmentType = (env) => {
  const types = { test: 'warning', staging: 'info', prod: 'danger' }
  return types[env] || 'info'
}

const getEnvironmentText = (env) => {
  const texts = { test: 'test', dev: 'dev', prod: 'prod' }
  return texts[env] || env
}

// 获取环境样式类
const getEnvironmentClass = (env) => {
  const classes = {
    test: 'env-test',
    staging: 'env-staging',
    prod: 'env-prod',
    dev: 'env-dev'
  }
  return classes[env] || 'env-default'
}

const getStatusType = (status) => {
  const types = { 1: 'info', 2: 'warning', 3: 'success', 4: 'danger' }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = { 1: '待发布', 2: '发布中', 3: '发布成功', 4: '发布失败' }
  return texts[status] || '未知'
}

// 根据部门ID获取部门名称
const getDeptNameById = (deptId) => {
  if (!deptId) return '-'
  const dept = businessDepts.value.find(d => d.id === deptId)
  return dept ? dept.deptName : `部门ID: ${deptId}`
}

const formatDuration = (duration) => {
  if (!duration) return '-'
  const minutes = Math.floor(duration / 60)
  const seconds = duration % 60
  return minutes > 0 ? `${minutes}分${seconds}秒` : `${seconds}秒`
}

// 格式化日期时间为 YYYY-MM-DD HH:mm:ss 格式
const formatDateTime = (dateTime) => {
  if (!dateTime) return '-'
  const date = new Date(dateTime)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 数据加载函数
const fetchReleaseList = async () => {
  loading.value = true
  try {
    const params = {
      ...queryParams,
      page: pagination.page,
      pageSize: pagination.pageSize
    }
    Object.keys(params).forEach(key => {
      if (params[key] === '' || params[key] === null) {
        delete params[key]
      }
    })

    const response = await appApi.getDeploymentList(params)
    if (response.data.code === 200) {
      releaseList.value = response.data.data.list || []
      pagination.total = response.data.data.total || 0
    } else {
      ElMessage.error(response.data.message || '获取发布列表失败')
    }
  } catch (error) {
    console.error('获取发布列表失败:', error)
    ElMessage.error('获取发布列表失败')
  } finally {
    loading.value = false
  }
}

// 加载业务组列表
const loadBusinessGroups = async () => {
  try {
    const response = await cmdbApi.getAllCmdbGroups()
    if (response.data?.code === 200) {
      businessGroups.value = response.data.data || []

      // 同时构建业务树数据（用于应用选择弹框）
      businessTreeData.value = businessGroups.value.map(group => ({
        id: group.id,
        name: group.name,
        type: 'group',
        children: group.children?.map(child => ({
          id: child.id,
          name: child.name,
          type: 'subgroup',
          parentId: group.id,
          parentName: group.name
        })) || []
      }))
    }
  } catch (error) {
    console.error('获取业务组列表失败:', error)
  }
}

// 生成级联选择器的业务组数据：只能选择二级分组
const cascaderBusinessGroups = computed(() => {
  return businessGroups.value.map(group => ({
    value: group.id,
    label: group.name,
    leaf: false, // 一级分组不是叶子节点，不可选中
    children: group.children?.filter(child => !child.children || child.children.length === 0).map(child => ({
      value: child.id,
      label: child.name,
      leaf: true // 二级分组是叶子节点，可以选中
    })) || []
  }))
})

const loadBusinessDepts = async () => {
  try {
    // 使用 queryDeptList 获取部门数据（与application.vue保持一致）
    const response = await systemApi.queryDeptList()
    if (response.data?.code === 200) {
      businessDepts.value = (response.data.data || []).filter(dept => dept.deptName)
    } else {
      // 如果第一个API失败，尝试第二个
      try {
        const response2 = await systemApi.querySysDeptVoList()
        if (response2.data?.code === 200) {
          businessDepts.value = (response2.data.data || []).filter(dept => dept.deptName)
        }
      } catch (error2) {
        console.error('获取业务部门列表失败:', error2)
      }
    }
  } catch (error) {
    console.error('获取业务部门列表失败:', error)
  }
}

// 获取通知机器人列表
const fetchNotifyRobots = async () => {
  try {
    notifyRobotLoading.value = true
    const response = await monitorApi.domain.getNotifyRobots()
    if (response.data?.code === 200) {
      notifyRobotList.value = response.data.data?.list || []
    } else {
      console.error('获取通知机器人列表失败:', response.data?.message)
    }
  } catch (error) {
    console.error('获取通知机器人列表失败:', error)
  } finally {
    notifyRobotLoading.value = false
  }
}


// 事件处理函数
const handleQuery = () => {
  pagination.page = 1
  fetchReleaseList()
}

const resetQuery = () => {
  Object.assign(queryParams, {
    title: '',
    environment: '',
    status: null,
    business_group_id: null
  })
  pagination.page = 1
  fetchReleaseList()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchReleaseList()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  fetchReleaseList()
}

const handleCreateRelease = () => {
  createDialogVisible.value = true
  resetCreateForm()
}

// 业务组变化处理
const handleGroupChange = () => {
  // 业务组变化时清空应用列表
  createForm.applications = []
}

// 业务部门变化处理
const handleDeptChange = () => {
  // 部门变化时不影响应用列表，只是记录部门ID
}

// 显示应用选择弹框
const showAppSelectionDialog = () => {
  appSelectionDialogVisible.value = true
  currentBusinessGroup.value = null
  businessGroupApps.value = []
  selectedAppsTemp.value = []
}

// 处理业务树节点点击
const handleBusinessTreeClick = async (data) => {
  // 只处理二级节点（子分组）点击
  if (data.type === 'subgroup') {
    currentBusinessGroup.value = {
      id: data.id,
      name: data.name,
      parentId: data.parentId,
      parentName: data.parentName
    }
    await loadBusinessGroupApps(data.id)
  }
}

// 加载业务组下的应用
const loadBusinessGroupApps = async (businessGroupId) => {
  if (!businessGroupId) {
    businessGroupApps.value = []
    return
  }

  loadingApps.value = true
  try {
    const params = {
      business_group_ids: [businessGroupId]
    }
    const response = await appApi.getServiceTree(params)
    if (response.data.code === 200) {
      const serviceTreeData = response.data.data || []

      // 提取应用列表
      const apps = []
      serviceTreeData.forEach(businessLine => {
        businessLine.services?.forEach(service => {
          apps.push({
            id: service.id,
            name: service.name,
            code: service.code,
            status: service.status,
            business_group_id: businessLine.business_group_id,
            business_group_name: businessLine.business_group_name,
            jenkins_envs: service.jenkins_envs || []
          })
        })
      })

      businessGroupApps.value = apps
    }
  } catch (error) {
    console.error('获取业务组应用失败:', error)
    ElMessage.error('获取业务组应用失败')
  } finally {
    loadingApps.value = false
  }
}

// 处理应用表格选择变化
const handleAppSelectionChange = (selection) => {
  selectedAppsTemp.value = selection
}

// 确认应用选择
const confirmAppSelection = () => {
  selectedAppsTemp.value.forEach(app => {
    // 检查是否已存在
    const exists = createForm.applications.some(existing => existing.app_id === app.id)
    if (!exists) {
      // 添加应用到列表，包含可用环境
      const availableEnvironments = app.jenkins_envs.map(env => ({
        value: env.env_name,
        label: getEnvironmentText(env.env_name)
      }))

      createForm.applications.push({
        app_id: app.id,
        app_name: app.name,
        business_group_name: app.business_group_name,
        environment: '', // 初始为空，用户需要选择
        available_environments: availableEnvironments,
        jenkins_envs: app.jenkins_envs // 保存完整的 jenkins_envs 配置
      })
    }
  })

  appSelectionDialogVisible.value = false
}

// 移除应用
const removeApplication = (index) => {
  createForm.applications.splice(index, 1)
}

// 环境变化处理
const handleEnvironmentChange = async (app) => {
  // 环境变化时清空已配置的参数
  app.parameters = null
  app.jenkins_server_id = null
  app.job_name = null
}

const isAppParametersConfigured = (app) =>
  app?.parameters !== null && app?.parameters !== undefined

// 配置应用的Jenkins参数
const handleConfigParams = async (app, index) => {
  if (!app.environment) {
    ElMessage.warning('请先选择环境')
    return
  }

  try {
    jenkinsParamsLoading.value = true
    currentAppConfig.value = app
    currentAppIndex.value = index

    // 根据应用和环境查找对应的 Jenkins 配置
    const jenkinsEnv = app.jenkins_envs?.find(env => env.env_name === app.environment)

    if (!jenkinsEnv || !jenkinsEnv.jenkins_server_id || !jenkinsEnv.job_name) {
      ElMessage.warning('该应用在当前环境未配置 Jenkins 任务')
      return
    }

    // 保存 Jenkins 配置信息到应用对象
    app.jenkins_server_id = jenkinsEnv.jenkins_server_id
    app.job_name = jenkinsEnv.job_name

    // 获取任务参数
    const response = await appApi.getJenkinsJobParameters(jenkinsEnv.jenkins_server_id, jenkinsEnv.job_name)
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      jenkinsJobParameters.value = responseData.data?.parameters || []

      // 初始化参数值(使用已保存的值或默认值)
      jenkinsParamValues.value = {}
      jenkinsJobParameters.value.forEach(param => {
        // PT_CHECKBOX 类型需要特殊处理
        if (param.type === 'PT_CHECKBOX') {
          // 如果已保存的值是字符串(逗号分隔),转换为数组
          if (app.parameters?.[param.name]) {
            const savedValue = app.parameters[param.name]
            jenkinsParamValues.value[param.name] = typeof savedValue === 'string'
              ? savedValue.split(',').map(v => v.trim()).filter(v => v)
              : (Array.isArray(savedValue) ? savedValue : [])
          } else {
            // 默认值也可能是逗号分隔的字符串
            jenkinsParamValues.value[param.name] = param.defaultValue
              ? (typeof param.defaultValue === 'string'
                  ? param.defaultValue.split(',').map(v => v.trim()).filter(v => v)
                  : [])
              : []
          }
        } else {
          jenkinsParamValues.value[param.name] = app.parameters?.[param.name] ?? param.defaultValue
        }
      })

      jenkinsParamsDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取任务参数失败')
    }
  } catch (error) {
    console.error('获取Jenkins任务参数失败:', error)
    ElMessage.error(error.response?.data?.message || '获取任务参数失败')
  } finally {
    jenkinsParamsLoading.value = false
  }
}

// 保存Jenkins参数配置
const saveJenkinsParams = () => {
  if (currentAppIndex.value >= 0 && createForm.applications[currentAppIndex.value]) {
    // 保存参数到应用对象
    createForm.applications[currentAppIndex.value].parameters = { ...jenkinsParamValues.value }
    ElMessage.success('参数配置已保存')
    closeParamsDialog()
  }
}

// 关闭参数对话框
const closeParamsDialog = () => {
  jenkinsParamsDialogVisible.value = false
  jenkinsJobParameters.value = []
  jenkinsParamValues.value = {}
  currentAppConfig.value = {}
  currentAppIndex.value = -1
}

// 验证应用环境配置
const validateAppEnvironment = async (app) => {
  if (!app.environment) return

  try {
    const response = await appApi.getEnvironment({
      app_id: app.app_id,
      environment: app.environment
    })

    if (response.data.code === 200) {
      const config = response.data.data
      if (!config.is_configured) {
        ElMessage.warning(`${app.app_name} 的 ${app.environment} 环境配置不完整`)
      }
    }
  } catch (error) {
    console.error('验证环境配置失败:', error)
  }
}


const handleCreateSubmit = async () => {
  try {
    await createFormRef.value.validate()

    if (createForm.applications.length === 0) {
      ElMessage.warning('请至少添加一个应用')
      return
    }

    // 检查所有应用是否都选择了环境
    const missingEnv = createForm.applications.find(app => !app.environment)
    if (missingEnv) {
      ElMessage.warning(`请为应用"${missingEnv.app_name}"选择环境`)
      return
    }

    const missingParamsConfig = createForm.applications.find(app => !isAppParametersConfigured(app))
    if (missingParamsConfig) {
      ElMessage.warning(`应用"${missingParamsConfig.app_name}"尚未配置参数，无法创建发布计划`)
      return
    }

    // 检查通知配置
    if (createForm.notice && !createForm.robot_id) {
      ElMessage.warning('启用通知时必须选择通知机器人')
      return
    }

    creating.value = true

    // 构建新的API格式
    const submitData = {
      title: createForm.title,
      business_group_id: createForm.business_group_id,
      business_dept_id: createForm.business_dept_id,
      description: createForm.description,
      applications: createForm.applications.map(app => {
        const appData = {
          app_id: app.app_id,
          environment: app.environment
        }

        // 如果配置了参数,将参数转换为字符串格式
        if (app.parameters) {
          appData.parameters = {}
          Object.keys(app.parameters).forEach(key => {
            const value = app.parameters[key]
            // 处理不同类型的参数
            if (Array.isArray(value)) {
              // 多选参数,用逗号连接
              appData.parameters[key] = value.join(',')
            } else if (typeof value === 'boolean') {
              // 布尔参数转换为字符串
              appData.parameters[key] = value ? 'true' : 'false'
            } else {
              appData.parameters[key] = value !== null && value !== undefined ? String(value) : ''
            }
          })
        }

        return appData
      })
    }

    // 添加通知配置
    if (createForm.notice && createForm.robot_id) {
      submitData.notice = true
      submitData.robot_id = createForm.robot_id
    }

    const response = await appApi.createQuickDeployment(submitData)
    if (response.data.code === 200) {
      ElMessage.success('发布创建成功')
      createDialogVisible.value = false
      resetCreateForm()
      fetchReleaseList()
    } else {
      ElMessage.error(response.data.message || '创建发布失败')
    }
  } catch (error) {
    console.error('创建发布失败:', error)
    ElMessage.error('创建发布失败')
  } finally {
    creating.value = false
  }
}

const resetCreateForm = () => {
  Object.assign(createForm, {
    title: '',
    business_group_id: null,
    business_dept_id: null,
    description: '',
    applications: [],
    notice: false,
    robot_id: null
  })
  serviceTree.value = []
}

const handleExecute = (row) => {
  // 将任务ID添加到已执行集合中，隐藏发布按钮
  executedReleases.value.add(row.id)
  // 跳转到发布详情页
  router.push(`/app/quick-temp/${row.id}`)
}

const handleDelete = async (row) => {
  try {
    // 检查发布状态，发布中不允许删除
    if (row.status === 2) {
      ElMessage.warning('发布中的任务不允许删除')
      return
    }

    await ElMessageBox.confirm(
      `确定要删除发布任务"${row.title}"吗？此操作将同时删除所有相关的子任务记录，且无法恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )

    // 调用删除API
    const response = await appApi.deleteDeployment(row.id)
    if (response.data.code === 200) {
      ElMessage.success('删除成功')
      // 刷新列表
      fetchReleaseList()
    } else {
      ElMessage.error(response.data.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

// 生命周期
onMounted(() => {
  fetchReleaseList()
  loadBusinessGroups() // 业务组列表
  loadBusinessDepts() // 业务部门列表
  fetchNotifyRobots() // 通知机器人列表
})
</script>

<style scoped>
.quick-release-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.glass-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 2px solid #e3e8f0;
}

.gradient-title {
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.search-section {
  margin-bottom: 24px;
  padding: 20px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.modern-btn {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.modern-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.table-section {
  margin-bottom: 24px;
}

.modern-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.modern-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.modern-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
}

.id-badge {
  background: linear-gradient(45deg, #667eea, #764ba2);
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.operation-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.pagination-section {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.modern-form .form-row {
  display: flex;
  gap: 20px;
}

.modern-form .form-item-half {
  flex: 1;
}

/* 应用配置区域样式 */
.applications-section {
  width: 100%;
}

.applications-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.app-count {
  color: #606266;
  font-size: 14px;
}

.selected-apps-list {
  space-y: 12px;
}

.selected-app-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #fafafa;
  margin-bottom: 8px;
}

.app-info {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.app-name {
  font-weight: 500;
  margin-bottom: 4px;
}

.app-environment {
  margin: 0 16px;
}

.empty-apps {
  text-align: center;
  padding: 40px 0;
}

/* 应用选择弹框样式 - 左右布局 */
.app-selection-dialog .el-dialog__body {
  padding: 20px;
  height: 500px;
}

.host-selection-layout {
  display: flex;
  height: 100%;
  gap: 16px;
}

/* 左侧业务树 */
.tree-section {
  width: 280px;
  border-right: 1px solid #ebeef5;
  padding-right: 16px;
}

.tree-header {
  margin-bottom: 12px;
}

.tree-header h4 {
  margin: 0;
  color: #303133;
  font-size: 16px;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 6px;
}

.tree-node .el-icon {
  color: #409eff;
}

/* 右侧应用表格 */
.services-section {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.services-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.services-header h4 {
  margin: 0;
  color: #303133;
  font-size: 16px;
}

.selected-group {
  color: #409eff;
  font-size: 14px;
  font-weight: 500;
}

.app-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-name-cell .el-icon {
  color: #409eff;
}

.environments-cell {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
}

.services-cell {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
}

.no-service {
  color: #909399;
  font-size: 12px;
  font-style: italic;
}

.no-env {
  color: #909399;
  font-size: 12px;
  font-style: italic;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.apps-list {
  margin-top: 16px;
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #e3e8f0;
  border-radius: 8px;
  padding: 12px;
}

.app-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.app-item:last-child {
  border-bottom: none;
}

.app-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.app-name {
  font-weight: 500;
  color: #303133;
}

.app-reason {
  font-size: 12px;
  color: #909399;
}

.modern-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  padding: 20px 24px;
}

.modern-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

/* 表格单元格图标样式 */
.table-cell-with-icon {
  display: flex;
  align-items: center;
  gap: 6px;
}

.table-cell-with-icon img {
  flex-shrink: 0;
}

.table-cell-with-icon span {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 发布服务样式 */
.services-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.service-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 8px;
  border-radius: 6px;
  background: rgba(64, 158, 255, 0.1);
}

.service-name {
  font-weight: 500;
  color: #303133;
  font-size: 13px;
}

.environment-wrapper {
  display: flex;
  align-items: center;
  gap: 3px;
  padding: 2px 6px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.8);
}

.environment-tag {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

/* 不同环境的颜色 */
.env-test {
  color: #e6a23c;
}

.env-dev {
  color: #67c23a;
}

.env-staging {
  color: #909399;
}

.env-prod {
  color: #f56c6c;
}

.env-default {
  color: #409eff;
}

/* Jenkins参数对话框样式 */
.jenkins-params-dialog :deep(.el-dialog__body) {
  padding: 20px;
}

.jenkins-params-dialog .no-params {
  padding: 40px;
  text-align: center;
}

.jenkins-params-dialog .params-form {
  max-height: 500px;
  overflow-y: auto;
  padding: 20px;
}

.jenkins-params-dialog .param-item {
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f2f5;
}

.jenkins-params-dialog .param-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.jenkins-params-dialog .param-description {
  display: block;
  font-size: 12px;
  color: #909399;
  margin-top: 6px;
  line-height: 1.4;
}

.jenkins-params-dialog :deep(.el-form-item__label) {
  font-weight: 600;
  color: #303133;
}

.jenkins-params-dialog :deep(.el-switch) {
  --el-switch-on-color: #67c23a;
}

.jenkins-params-dialog :deep(.el-checkbox-group) {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.jenkins-params-dialog :deep(.el-checkbox) {
  margin-right: 0;
  height: auto;
}

.app-params {
  display: flex;
  align-items: center;
  margin: 0 8px;
}

.selected-app-item {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
