<template>
  <div class="sysconfig-management">
    <el-card shadow="hover" class="config-card">
      <template #header>
        <div class="card-header">
          <span class="title">系统配置管理</span>
        </div>
      </template>

      <!-- 配置类型选择区域 -->
      <div class="config-type-section">
        <div class="type-label">配置类型：</div>
        <div class="type-buttons">
          <el-button
            :type="activeMenu === 'ldap' ? 'primary' : ''"
            size="small"
            @click="handleMenuClick('ldap')"
          >
            LDAP设置
          </el-button>
          <el-button
            :type="activeMenu === 'alertSource' ? 'primary' : ''"
            size="small"
            @click="handleMenuClick('alertSource')"
          >
            大盘数据配置
          </el-button>
          <el-button
            :type="activeMenu === 'aiModel' ? 'primary' : ''"
            size="small"
            @click="handleMenuClick('aiModel')"
          >
            Ai大模型设置
          </el-button>
          <el-button
            :type="activeMenu === 'notifyRobot' ? 'primary' : ''"
            size="small"
            @click="handleMenuClick('notifyRobot')"
          >
            报警推送设置
          </el-button>
          <el-button size="small" disabled>开放服务设置</el-button>
        </div>
      </div>

      <!-- LDAP配置表单 -->
      <div v-if="activeMenu === 'ldap'" class="config-form-section">
        <div class="form-title">
          <el-icon><Setting /></el-icon>
          <span>LDAP设置</span>
        </div>

        <el-form
          ref="ldapFormRef"
          :model="ldapForm"
          :rules="ldapRules"
          label-width="140px"
          class="ldap-form"
          v-loading="loading"
        >
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="LDAP服务器" prop="host">
                <el-input
                  v-model="ldapForm.host"
                  placeholder="例如：ldap-yd-in.dding.net"
                  clearable
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="端口" prop="port">
                <el-input-number
                  v-model="ldapForm.port"
                  :min="1"
                  :max="65535"
                  placeholder="389"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="Base DN" prop="baseDn">
                <el-input
                  v-model="ldapForm.baseDn"
                  placeholder="例如：dc=dding,dc=cn"
                  clearable
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="绑定用户" prop="bindUser">
                <el-input
                  v-model="ldapForm.bindUser"
                  placeholder="例如：cn=reader,dc=dding,dc=cn"
                  clearable
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="绑定密码" prop="bindPass">
                <el-input
                  v-model="ldapForm.bindPass"
                  type="password"
                  placeholder="请输入绑定密码"
                  show-password
                  clearable
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="认证过滤器" prop="authFilter">
                <el-input
                  v-model="ldapForm.authFilter"
                  placeholder="例如：(&(cn=%s))"
                  clearable
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="默认角色" prop="defaultRoleId">
                <el-select
                  v-model="ldapForm.defaultRoleId"
                  placeholder="请选择默认角色"
                  clearable
                  filterable
                  style="width: 100%"
                >
                  <el-option
                    v-for="role in enabledRoleList"
                    :key="role.id"
                    :label="role.roleName"
                    :value="role.id"
                  >
                    <div style="display: flex; justify-content: space-between; align-items: center;">
                      <span>{{ role.roleName }}</span>
                      <el-tag size="small" type="info" style="margin-left: 10px;">{{ role.roleKey }}</el-tag>
                    </div>
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="启用LDAP" prop="enable">
                <el-switch v-model="ldapForm.enable" />
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item>
            <el-button
              type="success"
              size="small"
              @click="testLdapConnection"
              :loading="testLoading"
            >
              <el-icon><Link /></el-icon>
              测试连接
            </el-button>
            <el-button
              type="primary"
              size="small"
              @click="saveLdapConfig"
              :loading="loading"
            >
              <el-icon><Check /></el-icon>
              保存设置
            </el-button>
            <el-button
              size="small"
              @click="resetForm"
            >
              <el-icon><RefreshLeft /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 告警源统计配置 -->
      <div v-if="activeMenu === 'alertSource'" class="config-form-section">
        <div class="form-title">
          <el-icon><DataLine /></el-icon>
          <span>告警源统计配置</span>
        </div>

        <!-- 告警源列表 -->
        <div class="alert-source-actions" style="margin-bottom: 15px;">
          <el-button type="primary" size="small" @click="handleAddAlertSource">
            <el-icon><Plus /></el-icon>
            新增数据源
          </el-button>
          <el-button size="small" @click="handleRefreshAlertSource">
            <el-icon><RefreshLeft /></el-icon>
            刷新
          </el-button>
        </div>

        <el-table
          :data="alertSourceList"
          v-loading="alertSourceLoading"
          border
          style="width: 100%"
        >
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="告警源名称" min-width="150" />
          <el-table-column prop="type" label="类型" width="120">
            <template #default="{ row }">
              <el-tag :type="getTypeTagType(row.type)">
                {{ getTypeLabel(row.type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="hostName" label="关联主机" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-if="row.type === 5">{{ row.hostName || '-' }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="accountAlias" label="Zabbix账号" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-if="row.type === 6">{{ row.accountAlias || '-' }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="appKey" label="App Key" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-if="row.appKey">********</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="apiBaseUrl" label="API Base URL" min-width="180" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-switch
                v-model="row.status"
                :active-value="1"
                :inactive-value="0"
                @change="handleStatusChange(row)"
              />
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatTimestamp(row.createTime) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="handleEditAlertSource(row)">
                编辑
              </el-button>
              <el-button type="danger" link size="small" @click="handleDeleteAlertSource(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <el-pagination
          v-model:current-page="alertSourcePage.pageNum"
          v-model:page-size="alertSourcePage.pageSize"
          :total="alertSourcePage.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          style="margin-top: 20px; justify-content: flex-end;"
          @size-change="handleAlertSourcePageChange"
          @current-change="handleAlertSourcePageChange"
        />
      </div>

      <!-- 报警推送设置 -->
      <div v-if="activeMenu === 'notifyRobot'" class="config-form-section">
        <NotifyRobotManage />
      </div>

      <!-- AI大模型配置 -->
      <div v-if="activeMenu === 'aiModel'" class="config-form-section">
        <div class="form-title">
          <el-icon><Cpu /></el-icon>
          <span>AI大模型配置</span>
        </div>

        <!-- AI模型列表 -->
        <div class="ai-model-actions" style="margin-bottom: 15px;">
          <el-button type="primary" size="small" @click="handleAddAiModel">
            <el-icon><Plus /></el-icon>
            新增模型
          </el-button>
          <el-button size="small" @click="handleRefreshAiModel">
            <el-icon><RefreshLeft /></el-icon>
            刷新
          </el-button>
        </div>

        <el-table
          :data="aiModelList"
          v-loading="aiModelLoading"
          border
          style="width: 100%"
        >
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="模型名称" min-width="150" />
          <el-table-column prop="type" label="类型" width="120">
            <template #default="{ row }">
              <el-tag>{{ row.type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="model" label="模型标识" min-width="150" show-overflow-tooltip />
          <el-table-column prop="url" label="API URL" min-width="250" show-overflow-tooltip />
          <el-table-column prop="apiKey" label="API Key" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-if="row.apiKey">{{ row.apiKey.substring(0, 10) }}...</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-switch
                v-model="row.status"
                :active-value="1"
                :inactive-value="0"
                @change="handleAiModelStatusChange(row)"
              />
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatTimestamp(row.createTime) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="handleEditAiModel(row)">
                编辑
              </el-button>
              <el-button type="danger" link size="small" @click="handleDeleteAiModel(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <el-pagination
          v-model:current-page="aiModelPage.pageNum"
          v-model:page-size="aiModelPage.pageSize"
          :total="aiModelPage.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          style="margin-top: 20px; justify-content: flex-end;"
          @size-change="handleAiModelPageChange"
          @current-change="handleAiModelPageChange"
        />
      </div>

      <!-- 告警源新增/编辑对话框 -->
      <el-dialog
        v-model="alertSourceDialogVisible"
        :title="alertSourceForm.id ? '编辑告警源' : '新增告警源'"
        width="600px"
        :close-on-click-modal="false"
      >
        <el-form
          ref="alertSourceFormRef"
          :model="alertSourceForm"
          :rules="alertSourceRules"
          label-width="120px"
        >
          <el-form-item label="告警源名称" prop="name">
            <el-input v-model="alertSourceForm.name" placeholder="请输入告警源名称" clearable />
          </el-form-item>
          <el-form-item label="类型" prop="type">
            <el-select
              v-model="alertSourceForm.type"
              placeholder="请选择类型"
              style="width: 100%"
              @change="handleTypeChange"
            >
              <el-option :value="1" label="FlashDuty" />
              <el-option :value="2" label="阿里云" />
              <el-option :value="3" label="腾讯云" />
              <el-option :value="4" label="AWS云" />
              <el-option :value="5" label="主机" />
              <el-option :value="6" label="Zabbix" />
            </el-select>
          </el-form-item>
          <el-form-item
            v-if="isCloudType"
            label="云账号密钥"
            prop="keyId"
          >
            <el-select
              v-model="alertSourceForm.keyId"
              placeholder="请选择云账号密钥"
              style="width: 100%"
              filterable
            >
              <el-option
                v-for="key in filteredCloudKeys"
                :key="key.id"
                :label="key.keyName || `${key.remark}`"
                :value="key.id"
              >
                <div style="display: flex; justify-content: space-between; align-items: center;">
                  <span>{{ key.keyName || key.remark }}</span>
                  <el-tag size="small" type="info" style="margin-left: 10px;">
                    {{ getCloudKeyTypeLabel(key.keyType) }}
                  </el-tag>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item
            v-if="alertSourceForm.type === 5"
            label="选择主机"
            prop="hostId"
          >
            <el-select
              v-model="alertSourceForm.hostId"
              placeholder="请输入主机名搜索"
              style="width: 100%"
              filterable
              remote
              :remote-method="searchHosts"
              :loading="hostSearchLoading"
              clearable
            >
              <el-option
                v-for="host in hostList"
                :key="host.id"
                :label="`${host.hostName} (${host.publicIp || host.privateIp || host.sshIp})`"
                :value="host.id"
              >
                <div style="display: flex; justify-content: space-between; align-items: center;">
                  <span>{{ host.hostName }}</span>
                  <el-tag size="small" type="info" style="margin-left: 10px;">
                    {{ host.publicIp || host.privateIp || host.sshIp }}
                  </el-tag>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item
            v-if="alertSourceForm.type === 6"
            label="Zabbix账号"
            prop="accountId"
          >
            <el-select
              v-model="alertSourceForm.accountId"
              placeholder="请选择Zabbix账号"
              style="width: 100%"
              filterable
            >
              <el-option
                v-for="account in zabbixAccountList"
                :key="account.id"
                :label="account.alias"
                :value="account.id"
              >
                <div style="display: flex; justify-content: space-between; align-items: center;">
                  <span>{{ account.alias }}</span>
                  <el-tag size="small" type="info" style="margin-left: 10px;">
                    {{ account.host }}
                  </el-tag>
                </div>
              </el-option>
            </el-select>
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              在 <router-link to="/configcenter/accountauth" style="color: #409EFF;">配置中心 > 账号认证管理</router-link> 中管理Zabbix账号
            </div>
          </el-form-item>
          <el-form-item
            v-if="alertSourceForm.type === 1"
            label="App Key"
            prop="appKey"
          >
            <el-input
              v-model="alertSourceForm.appKey"
              type="textarea"
              :rows="3"
              placeholder="请输入App Key"
              clearable
            />
          </el-form-item>
          <el-form-item
            v-if="alertSourceForm.type === 1"
            label="API Base URL"
            prop="apiBaseUrl"
          >
            <el-input
              v-model="alertSourceForm.apiBaseUrl"
              placeholder="留空将自动设置默认URL"
              clearable
            />
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-switch v-model="alertSourceForm.status" :active-value="1" :inactive-value="0" />
          </el-form-item>
          <el-form-item label="备注" prop="remark">
            <el-input
              v-model="alertSourceForm.remark"
              type="textarea"
              :rows="2"
              placeholder="请输入备注"
              clearable
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="alertSourceDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSaveAlertSource" :loading="alertSourceSaving">
            保存
          </el-button>
        </template>
      </el-dialog>

      <!-- AI模型新增/编辑对话框 -->
      <el-dialog
        v-model="aiModelDialogVisible"
        :title="aiModelForm.id ? '编辑AI模型' : '新增AI模型'"
        width="600px"
        :close-on-click-modal="false"
      >
        <el-form
          ref="aiModelFormRef"
          :model="aiModelForm"
          :rules="aiModelRules"
          label-width="120px"
        >
          <el-form-item label="模型名称" prop="name">
            <el-input v-model="aiModelForm.name" placeholder="例如：私人大模型" clearable />
          </el-form-item>
          <el-form-item label="类型" prop="type">
            <el-select
              v-model="aiModelForm.type"
              placeholder="请选择类型"
              style="width: 100%"
              @change="handleAiTypeChange"
            >
              <el-option value="阿里百炼" label="阿里百炼" />
              <el-option value="DeepSeek" label="DeepSeek" />
              <el-option value="OpenAI" label="OpenAI" />
              <el-option value="智谱AI" label="智谱AI" />

            </el-select>
          </el-form-item>
          <el-form-item label="模型标识" prop="model">
            <el-select
              v-model="aiModelForm.model"
              placeholder="请选择或输入模型"
              style="width: 100%"
              filterable
              allow-create
              default-first-option
              clearable
            >
              <el-option
                v-for="model in availableModels"
                :key="model.value"
                :label="model.label"
                :value="model.value"
              >
                <div style="display: flex; justify-content: space-between; align-items: center; width: 100%;">
                  <span>{{ model.value }}</span>
                  <div style="display: flex; align-items: center; gap: 8px;">
                    <el-tag size="small" :type="model.tagType">{{ model.desc }}</el-tag>
                    <span style="color: #909399; font-size: 12px;">{{ model.tokens }}</span>
                  </div>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="API URL" prop="url">
            <el-input
              v-model="aiModelForm.url"
              placeholder="例如：https://dashscope.aliyuncs.com/api/v1/..."
              clearable
            />
          </el-form-item>
          <el-form-item label="API Key" prop="apiKey">
            <el-input
              v-model="aiModelForm.apiKey"
              type="textarea"
              :rows="2"
              placeholder="请输入API Key"
              clearable
            />
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-switch v-model="aiModelForm.status" :active-value="1" :inactive-value="0" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="aiModelDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSaveAiModel" :loading="aiModelSaving">
            保存
          </el-button>
        </template>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import { Setting, Link, Check, RefreshLeft, DataLine, Plus, Cpu } from '@element-plus/icons-vue'
import systemApi from '@/api/system'
import dashboardApi from '@/api/dashboard'
import configApi from '@/api/config'
import cmdbApi from '@/api/cmdb'
import aiApi from '@/api/ai'
import { ElMessage, ElMessageBox } from 'element-plus'
import NotifyRobotManage from './notify/NotifyRobotManage.vue'

export default {
  components: {
    Setting,
    Link,
    Check,
    RefreshLeft,
    DataLine,
    Plus,
    Cpu,
    NotifyRobotManage
  },
  data() {
    return {
      activeMenu: 'ldap',
      loading: false,
      testLoading: false,
      roleList: [],
      ldapForm: {
        enable: false,
        host: '',
        port: 389,
        baseDn: '',
        bindUser: '',
        bindPass: '',
        authFilter: '(&(cn=%s))',
        coverAttributes: true,
        tls: false,
        startTLS: false,
        defaultRoleId: null,
        attributes: {
          nickname: 'cn',
          phone: 'mobile',
          email: 'mail'
        },
        remark: ''
      },
      ldapRules: {
        host: [
          { required: true, message: '请输入LDAP服务器地址', trigger: 'blur' }
        ],
        port: [
          { required: true, message: '请输入端口', trigger: 'blur' }
        ],
        baseDn: [
          { required: true, message: '请输入Base DN', trigger: 'blur' }
        ],
        bindUser: [
          { required: true, message: '请输入绑定用户', trigger: 'blur' }
        ],
        bindPass: [
          { required: true, message: '请输入绑定密码', trigger: 'blur' }
        ],
        authFilter: [
          { required: true, message: '请输入认证过滤器', trigger: 'blur' }
        ]
      },
      // 告警源相关
      alertSourceList: [],
      alertSourceLoading: false,
      alertSourcePage: {
        pageNum: 1,
        pageSize: 10,
        total: 0
      },
      alertSourceDialogVisible: false,
      alertSourceSaving: false,
      alertSourceForm: {
        id: null,
        name: '',
        type: 1,
        keyId: null,
        hostId: null,
        accountId: null,
        appKey: '',
        apiBaseUrl: '',
        status: 1,
        remark: ''
      },
      cloudKeyList: [],
      hostList: [],
      zabbixAccountList: [],
      hostSearchLoading: false,
      alertSourceRules: {
        name: [
          { required: true, message: '请输入告警源名称', trigger: 'blur' }
        ],
        type: [
          { required: true, message: '请选择类型', trigger: 'change' }
        ],
        keyId: [
          {
            validator: (rule, value, callback) => {
              if (this.isCloudType && this.alertSourceForm.type !== 5 && !value) {
                callback(new Error('请选择云账号密钥'))
              } else {
                callback()
              }
            },
            trigger: 'change'
          }
        ],
        hostId: [
          {
            validator: (rule, value, callback) => {
              if (this.alertSourceForm.type === 5 && !value) {
                callback(new Error('请选择主机'))
              } else {
                callback()
              }
            },
            trigger: 'change'
          }
        ],
        accountId: [
          {
            validator: (rule, value, callback) => {
              if (this.alertSourceForm.type === 6 && !value) {
                callback(new Error('请选择Zabbix账号'))
              } else {
                callback()
              }
            },
            trigger: 'change'
          }
        ],
        appKey: [
          {
            validator: (rule, value, callback) => {
              if (this.alertSourceForm.type === 1 && !value) {
                callback(new Error('请输入App Key'))
              } else {
                callback()
              }
            },
            trigger: 'blur'
          }
        ]
      },
      // AI模型相关
      aiModelList: [],
      aiModelLoading: false,
      aiModelPage: {
        pageNum: 1,
        pageSize: 10,
        total: 0
      },
      aiModelDialogVisible: false,
      aiModelSaving: false,
      aiModelForm: {
        id: null,
        name: '',
        type: '阿里百炼',
        url: '',
        apiKey: '',
        model: 'qwen-plus',
        status: 1
      },
      // AI模型配置字典
      modelConfigs: {
        '阿里百炼': [
          { value: 'qwen-turbo', label: 'qwen-turbo', desc: '⚡ 速度快、价格低', tokens: '32K tokens', tagType: 'success', recommend: '高频/低成本场景' },
          { value: 'qwen-plus', label: 'qwen-plus', desc: '🧠 能力均衡、性价比高', tokens: '131K tokens', tagType: 'primary', recommend: '最推荐通用场景' },
          { value: 'qwen-max', label: 'qwen-max', desc: '🏆 最强性能、最强推理', tokens: '32K tokens', tagType: 'danger', recommend: '高质量输出' },
          { value: 'qwen-vl-plus', label: 'qwen-vl-plus', desc: '🖼️ 多模态（图文理解）', tokens: '-', tagType: 'warning', recommend: '图片理解' }
        ],
        'OpenAI': [
          { value: 'gpt-3.5-turbo', label: 'gpt-3.5-turbo', desc: '⚡ 速度快、成本低', tokens: '16K tokens', tagType: 'success', recommend: '日常对话' },
          { value: 'gpt-4', label: 'gpt-4', desc: '🧠 强大推理能力', tokens: '8K tokens', tagType: 'primary', recommend: '复杂任务' },
          { value: 'gpt-4-turbo', label: 'gpt-4-turbo', desc: '🚀 更快的GPT-4', tokens: '128K tokens', tagType: 'danger', recommend: '大上下文' }
        ],
        '智谱AI': [
          { value: 'glm-4', label: 'glm-4', desc: '🧠 通用对话模型', tokens: '128K tokens', tagType: 'primary', recommend: '通用场景' },
          { value: 'glm-4v', label: 'glm-4v', desc: '🖼️ 多模态理解', tokens: '2K tokens', tagType: 'warning', recommend: '图文理解' },
          { value: 'glm-3-turbo', label: 'glm-3-turbo', desc: '⚡ 高性价比', tokens: '128K tokens', tagType: 'success', recommend: '快速响应' }
        ],
        'DeepSeek': [
          { value: 'deepseek-v3-250324', label: 'deepseek-v3-250324', desc: '🏆 最强版本', tokens: '8K tokens', tagType: 'danger', recommend: '高精度' },
        ]
      },
      aiModelRules: {
        name: [
          { required: true, message: '请输入模型名称', trigger: 'blur' }
        ],
        type: [
          { required: true, message: '请选择类型', trigger: 'change' }
        ],
        url: [
          { required: true, message: '请输入API URL', trigger: 'blur' }
        ],
        apiKey: [
          { required: true, message: '请输入API Key', trigger: 'blur' }
        ],
        model: [
          { required: true, message: '请输入模型标识', trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    // 切换菜单
    handleMenuClick(key) {
      this.activeMenu = key
      if (key === 'alertSource') {
        this.getAlertSourceList()
      } else if (key === 'aiModel') {
        this.getAiModelList()
      }
    },

    // 获取LDAP配置
    async getLdapConfig() {
      try {
        this.loading = true
        const res = await systemApi.getLdapConfig()

        // 注意：响应数据结构是 res.data.data
        if (res.data && res.data.code === 200 && res.data.data) {
          const data = res.data.data

          // 直接赋值，处理默认值
          this.ldapForm.enable = !!data.enable
          this.ldapForm.host = data.host || ''
          this.ldapForm.port = data.port || 389
          this.ldapForm.baseDn = data.baseDn || ''
          this.ldapForm.bindUser = data.bindUser || ''
          this.ldapForm.bindPass = data.bindPass || ''
          this.ldapForm.authFilter = data.authFilter || '(&(cn=%s))'
          this.ldapForm.coverAttributes = data.coverAttributes !== false
          this.ldapForm.tls = !!data.tls
          this.ldapForm.startTLS = !!data.startTLS

          // defaultRoleId 为 0 时设置为 null（0表示未设置）
          this.ldapForm.defaultRoleId = data.defaultRoleId && data.defaultRoleId > 0 ? data.defaultRoleId : null

          // 处理属性映射
          if (data.attributes) {
            this.ldapForm.attributes.nickname = data.attributes.nickname || 'cn'
            this.ldapForm.attributes.phone = data.attributes.phone || 'mobile'
            this.ldapForm.attributes.email = data.attributes.email || 'mail'
          }

          this.ldapForm.remark = data.remark || ''
        }
      } catch (error) {
        console.error('获取LDAP配置失败:', error)
        // 首次加载失败不显示错误提示（可能还没有配置）
      } finally {
        // 延迟关闭loading，确保数据已渲染
        setTimeout(() => {
          this.loading = false
        }, 100)
      }
    },

    // 获取角色列表
    async getRoleList() {
      try {
        const res = await systemApi.querySysRoleVoList()
        if (res.data && res.data.code === 200) {
          this.roleList = res.data.data || []
        }
      } catch (error) {
        console.error('获取角色列表失败:', error)
      }
    },

    // 测试LDAP连接
    async testLdapConnection() {
      if (!this.$refs.ldapFormRef) return

      try {
        await this.$refs.ldapFormRef.validate()
        this.testLoading = true

        // 测试连接不需要传参数，会使用已保存的配置
        const res = await systemApi.testLdapConnection()

        if (res.data && res.data.code === 200) {
          this.$message.success(res.data.data?.message || res.data.message || 'LDAP连接测试成功')
        } else {
          this.$message.error(res.data?.message || 'LDAP连接测试失败')
        }
      } catch (error) {
        if (error !== false) {
          this.$message.error('LDAP连接测试失败: ' + (error.message || '未知错误'))
        }
      } finally {
        this.testLoading = false
      }
    },

    // 保存LDAP配置
    async saveLdapConfig() {
      if (!this.$refs.ldapFormRef) return

      try {
        await this.$refs.ldapFormRef.validate()
        this.loading = true

        // 构造请求数据，使用后端要求的字段结构
        const requestData = {
          enable: this.ldapForm.enable,
          host: this.ldapForm.host,
          port: this.ldapForm.port,
          baseDn: this.ldapForm.baseDn,
          bindUser: this.ldapForm.bindUser,
          bindPass: this.ldapForm.bindPass,
          authFilter: this.ldapForm.authFilter,
          coverAttributes: this.ldapForm.coverAttributes,
          tls: this.ldapForm.tls,
          startTLS: this.ldapForm.startTLS,
          defaultRoleId: this.ldapForm.defaultRoleId,
          attributes: {
            nickname: this.ldapForm.attributes.nickname,
            phone: this.ldapForm.attributes.phone,
            email: this.ldapForm.attributes.email
          },
          remark: this.ldapForm.remark
        }

        const res = await systemApi.saveLdapConfig(requestData)

        if (res.data && res.data.code === 200) {
          this.$message.success('LDAP配置保存成功')
          this.getLdapConfig()
        } else {
          this.$message.error(res.data?.message || 'LDAP配置保存失败')
        }
      } catch (error) {
        if (error !== false) {
          this.$message.error('LDAP配置保存失败: ' + (error.message || '未知错误'))
        }
      } finally {
        this.loading = false
      }
    },

    // 重置表单
    resetForm() {
      this.$refs.ldapFormRef.resetFields()
      this.getLdapConfig()
    },

    // ========== 告警源管理方法 ==========
    // 获取告警源列表
    async getAlertSourceList() {
      try {
        this.alertSourceLoading = true
        const res = await dashboardApi.getAlertSourceList({
          pageNum: this.alertSourcePage.pageNum,
          pageSize: this.alertSourcePage.pageSize
        })

        if (res.data && (res.data.code === 0 || res.data.code === 200)) {
          this.alertSourceList = res.data.data.list || []
          this.alertSourcePage.total = res.data.data.total || 0
        } else {
          ElMessage.error(res.data?.message || res.data?.msg || '获取告警源列表失败')
        }
      } catch (error) {
        console.error('获取告警源列表失败:', error)
        ElMessage.error('获取告警源列表失败')
      } finally {
        this.alertSourceLoading = false
      }
    },

    // 刷新告警源列表
    handleRefreshAlertSource() {
      this.alertSourcePage.pageNum = 1
      this.getAlertSourceList()
    },

    // 分页变化
    handleAlertSourcePageChange() {
      this.getAlertSourceList()
    },

    // 新增告警源
    handleAddAlertSource() {
      this.alertSourceForm = {
        id: null,
        name: '',
        type: 1,
        keyId: null,
        hostId: null,
        accountId: null,
        appKey: '',
        apiBaseUrl: '',
        status: 1,
        remark: ''
      }
      this.hostList = []
      this.getCloudKeyList()
      this.getZabbixAccountList()
      this.alertSourceDialogVisible = true
      this.$nextTick(() => {
        this.$refs.alertSourceFormRef?.clearValidate()
      })
    },

    // 编辑告警源
    async handleEditAlertSource(row) {
      this.alertSourceForm = {
        ...row,
        keyId: row.keyId || null,
        hostId: row.hostId || null,
        // 对于Zabbix类型,keyId实际存储的是accountId
        accountId: row.type === 6 ? (row.keyId || null) : null
      }
      this.getCloudKeyList()
      this.getZabbixAccountList()

      // 如果是主机类型且有主机ID，加载主机信息
      if (row.type === 5 && row.hostId) {
        await this.loadHostById(row.hostId)
      }

      this.alertSourceDialogVisible = true
      this.$nextTick(() => {
        this.$refs.alertSourceFormRef?.clearValidate()
      })
    },

    // 保存告警源
    async handleSaveAlertSource() {
      if (!this.$refs.alertSourceFormRef) return

      try {
        await this.$refs.alertSourceFormRef.validate()
        this.alertSourceSaving = true

        const apiMethod = this.alertSourceForm.id
          ? dashboardApi.updateAlertSource
          : dashboardApi.createAlertSource

        // 构造提交数据,对于Zabbix类型,将accountId保存到keyId字段
        const submitData = { ...this.alertSourceForm }
        if (submitData.type === 6) {
          // Zabbix类型使用accountId作为keyId
          submitData.keyId = submitData.accountId
        }

        const res = await apiMethod(submitData)

        if (res.data && (res.data.code === 0 || res.data.code === 200)) {
          ElMessage.success(this.alertSourceForm.id ? '更新成功' : '创建成功')
          this.alertSourceDialogVisible = false
          this.getAlertSourceList()
        } else {
          ElMessage.error(res.data?.message || res.data?.msg || '保存失败')
        }
      } catch (error) {
        if (error !== false) {
          ElMessage.error('保存失败: ' + (error.message || '未知错误'))
        }
      } finally {
        this.alertSourceSaving = false
      }
    },

    // 删除告警源
    async handleDeleteAlertSource(row) {
      try {
        await ElMessageBox.confirm('确定删除该告警源吗？', '提示', {
          type: 'warning'
        })

        const res = await dashboardApi.deleteAlertSource(row.id)

        if (res.data && (res.data.code === 0 || res.data.code === 200)) {
          ElMessage.success('删除成功')
          this.getAlertSourceList()
        } else {
          ElMessage.error(res.data?.message || res.data?.msg || '删除失败')
        }
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除失败:', error)
          ElMessage.error('删除失败')
        }
      }
    },

    // 更新告警源状态
    async handleStatusChange(row) {
      try {
        const res = await dashboardApi.updateAlertSourceStatus(row.id, row.status)

        if (res.data && (res.data.code === 0 || res.data.code === 200)) {
          ElMessage.success('状态更新成功')
        } else {
          ElMessage.error(res.data?.message || res.data?.msg || '状态更新失败')
          // 恢复原状态
          row.status = row.status === 1 ? 0 : 1
        }
      } catch (error) {
        console.error('状态更新失败:', error)
        ElMessage.error('状态更新失败')
        // 恢复原状态
        row.status = row.status === 1 ? 0 : 1
      }
    },

    // 格式化时间戳
    formatTimestamp(timestamp) {
      if (!timestamp) return '-'
      const date = new Date(timestamp * 1000)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    },

    // 获取云密钥列表
    async getCloudKeyList() {
      try {
        const res = await configApi.getKeyManageList({
          page: 1,
          pageSize: 100
        })

        if (res.data && res.data.code === 200) {
          this.cloudKeyList = res.data.data?.list || []
        }
      } catch (error) {
        console.error('获取云密钥列表失败:', error)
      }
    },

    // 获取Zabbix账号列表
    async getZabbixAccountList() {
      try {
        const res = await configApi.listAccountAuth({
          type: 5, // 5=Zabbix
          page: 1,
          pageSize: 100
        })

        if (res.data && res.data.code === 200) {
          this.zabbixAccountList = res.data.data?.list || []
        }
      } catch (error) {
        console.error('获取Zabbix账号列表失败:', error)
      }
    },

    // 类型变化处理
    handleTypeChange() {
      // 重置密钥ID、主机ID、账号ID和appKey
      this.alertSourceForm.keyId = null
      this.alertSourceForm.hostId = null
      this.alertSourceForm.accountId = null
      this.hostList = []

      if (this.isCloudType && this.alertSourceForm.type !== 5 && this.alertSourceForm.type !== 6) {
        this.alertSourceForm.appKey = ''
        this.alertSourceForm.apiBaseUrl = ''
      }

      // 清除验证
      this.$nextTick(() => {
        this.$refs.alertSourceFormRef?.clearValidate()
      })
    },

    // 搜索主机
    async searchHosts(query) {
      if (!query || query.trim() === '') {
        this.hostList = []
        return
      }

      try {
        this.hostSearchLoading = true
        const res = await cmdbApi.GetCmdbHostsByHostNameLike(query, {
          pageNum: 1,
          pageSize: 20
        })

        if (res.data && res.data.code === 200) {
          // 后端返回的data直接是数组，不是{list: []}结构
          const data = res.data.data
          this.hostList = Array.isArray(data) ? data : (data?.list || [])
        } else {
          this.hostList = []
        }
      } catch (error) {
        console.error('搜索主机失败:', error)
        this.hostList = []
      } finally {
        this.hostSearchLoading = false
      }
    },

    // 根据主机ID加载主机信息
    async loadHostById(hostId) {
      try {
        const res = await cmdbApi.getHostList({
          pageNum: 1,
          pageSize: 1,
          id: hostId
        })

        if (res.data && res.data.code === 200 && res.data.data?.list?.length > 0) {
          const host = res.data.data.list[0]
          this.hostList = [host]
        }
      } catch (error) {
        console.error('加载主机信息失败:', error)
      }
    },

    // 获取类型标签
    getTypeLabel(type) {
      const typeMap = {
        1: 'FlashDuty',
        2: '阿里云',
        3: '腾讯云',
        4: 'AWS云',
        5: '主机',
        6: 'Zabbix'
      }
      return typeMap[type] || '未知'
    },

    // 获取类型标签样式
    getTypeTagType(type) {
      const tagMap = {
        1: 'primary',
        2: 'success',
        3: 'warning',
        4: 'danger',
        5: 'info',
        6: ''
      }
      return tagMap[type] || 'info'
    },

    // 获取云密钥类型标签（keyType映射）
    getCloudKeyTypeLabel(keyType) {
      const typeMap = {
        1: '阿里云',
        2: '腾讯云',
        5: 'AWS云'
      }
      return typeMap[keyType] || '未知'
    },

    // ========== AI模型管理方法 ==========
    // 获取AI模型列表
    async getAiModelList() {
      try {
        this.aiModelLoading = true
        const res = await aiApi.getAiModelList({
          pageNum: this.aiModelPage.pageNum,
          pageSize: this.aiModelPage.pageSize
        })

        if (res.data && (res.data.code === 0 || res.data.code === 200)) {
          this.aiModelList = res.data.data.list || []
          this.aiModelPage.total = res.data.data.total || 0
        } else {
          ElMessage.error(res.data?.message || res.data?.msg || '获取AI模型列表失败')
        }
      } catch (error) {
        console.error('获取AI模型列表失败:', error)
        ElMessage.error('获取AI模型列表失败')
      } finally {
        this.aiModelLoading = false
      }
    },

    // 刷新AI模型列表
    handleRefreshAiModel() {
      this.aiModelPage.pageNum = 1
      this.getAiModelList()
    },

    // 分页变化
    handleAiModelPageChange() {
      this.getAiModelList()
    },

    // 新增AI模型
    handleAddAiModel() {
      this.aiModelForm = {
        id: null,
        name: '',
        type: '阿里百炼',
        url: '',
        apiKey: '',
        model: 'qwen-plus',
        status: 1
      }
      this.aiModelDialogVisible = true
      this.$nextTick(() => {
        this.$refs.aiModelFormRef?.clearValidate()
      })
    },

    // AI模型类型变化处理
    handleAiTypeChange() {
      // 如果当前是预置模型，则跟随类型切换到默认模型；手动输入的自定义值则保留。
      const defaultModels = {
        '阿里百炼': 'qwen-plus',
        'OpenAI': 'gpt-3.5-turbo',
        '智谱AI': 'glm-4',
        '文心一言': 'ernie-3.5'
      }
      const presetModels = Object.values(this.modelConfigs)
        .reduce((models, configList) => models.concat(configList), [])
        .map(model => model.value)

      if (!this.aiModelForm.model || presetModels.includes(this.aiModelForm.model)) {
        this.aiModelForm.model = defaultModels[this.aiModelForm.type] || ''
      }

      // 清除验证
      this.$nextTick(() => {
        this.$refs.aiModelFormRef?.clearValidate('model')
      })
    },

    // 编辑AI模型
    handleEditAiModel(row) {
      this.aiModelForm = { ...row }
      this.aiModelDialogVisible = true
      this.$nextTick(() => {
        this.$refs.aiModelFormRef?.clearValidate()
      })
    },

    // 保存AI模型
    async handleSaveAiModel() {
      if (!this.$refs.aiModelFormRef) return

      try {
        await this.$refs.aiModelFormRef.validate()
        this.aiModelSaving = true

        const apiMethod = this.aiModelForm.id
          ? aiApi.updateAiModel
          : aiApi.createAiModel

        const res = await apiMethod(this.aiModelForm)

        if (res.data && (res.data.code === 0 || res.data.code === 200)) {
          ElMessage.success(this.aiModelForm.id ? '更新成功' : '创建成功')
          this.aiModelDialogVisible = false
          this.getAiModelList()
        } else {
          ElMessage.error(res.data?.message || res.data?.msg || '保存失败')
        }
      } catch (error) {
        if (error !== false) {
          ElMessage.error('保存失败: ' + (error.message || '未知错误'))
        }
      } finally {
        this.aiModelSaving = false
      }
    },

    // 删除AI模型
    async handleDeleteAiModel(row) {
      try {
        await ElMessageBox.confirm('确定删除该AI模型吗？', '提示', {
          type: 'warning'
        })

        const res = await aiApi.deleteAiModel(row.id)

        if (res.data && (res.data.code === 0 || res.data.code === 200)) {
          ElMessage.success('删除成功')
          this.getAiModelList()
        } else {
          ElMessage.error(res.data?.message || res.data?.msg || '删除失败')
        }
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除失败:', error)
          ElMessage.error('删除失败')
        }
      }
    },

    // 更新AI模型状态
    async handleAiModelStatusChange(row) {
      try {
        const res = await aiApi.updateAiModelStatus(row.id, row.status)

        if (res.data && (res.data.code === 0 || res.data.code === 200)) {
          ElMessage.success('状态更新成功')
        } else {
          ElMessage.error(res.data?.message || res.data?.msg || '状态更新失败')
          // 恢复原状态
          row.status = row.status === 1 ? 0 : 1
        }
      } catch (error) {
        console.error('状态更新失败:', error)
        ElMessage.error('状态更新失败')
        // 恢复原状态
        row.status = row.status === 1 ? 0 : 1
      }
    }
  },
  mounted() {
    this.getLdapConfig()
    this.getRoleList()
  },
  computed: {
    // 获取角色列表（vo接口返回的数据没有status字段，直接返回所有角色）
    enabledRoleList() {
      return this.roleList || []
    },
    // 判断是否为云账号类型
    isCloudType() {
      return this.alertSourceForm.type >= 2 && this.alertSourceForm.type <= 4
    },
    // 根据选择的云厂商类型过滤密钥列表
    filteredCloudKeys() {
      if (!this.alertSourceForm.type || !this.isCloudType) {
        return []
      }
      // type 2=阿里云(keyType=1), 3=腾讯云(keyType=2), 4=AWS云(keyType=5)
      const typeToKeyTypeMap = {
        2: 1, // 阿里云
        3: 2, // 腾讯云
        4: 5  // AWS云
      }
      const keyType = typeToKeyTypeMap[this.alertSourceForm.type]
      return this.cloudKeyList.filter(key => key.keyType === keyType)
    },
    // 根据选择的AI模型类型返回可用的模型列表
    availableModels() {
      return this.modelConfigs[this.aiModelForm.type] || []
    }
  }
}
</script>

<style scoped>
/* 现代化科技感设计风格 */

.sysconfig-management {
  padding: 20px;
  min-height: 100vh;
  background: #eef1f6;
}

.config-card {
  height: 800px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
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
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

/* 配置类型选择区域 */
.config-type-section {
  margin-bottom: 30px;
  padding: 20px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
  display: flex;
  align-items: center;
  gap: 16px;
}

.type-label {
  font-size: 14px;
  font-weight: 600;
  color: #606266;
  white-space: nowrap;
}

.type-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

/* 表单区域 */
.config-form-section {
  padding: 20px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.form-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 2px solid rgba(103, 126, 234, 0.2);
}

.form-title .el-icon {
  font-size: 18px;
  color: #667eea;
}

/* 表单样式 */
.ldap-form {
  margin-top: 20px;
}

.ldap-form :deep(.el-form-item__label) {
  font-weight: 500;
  color: #606266;
}

.ldap-form :deep(.el-form-item) {
  margin-bottom: 18px;
}

/* 输入框和选择器样式 */
.ldap-form :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.ldap-form :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.ldap-form :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.ldap-form :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #2c3e50;
}

/* 按钮样式优化 */
.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover:not(.is-disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.el-button.is-disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 加载遮罩样式 */
:deep(.el-loading-mask) {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

/* 响应式布局 */
@media (max-width: 768px) {
  .config-type-section {
    flex-direction: column;
    align-items: flex-start;
  }

  .type-buttons {
    width: 100%;
  }

  .type-buttons .el-button {
    flex: 1;
    min-width: 120px;
  }

  .ldap-form :deep(.el-form-item__label) {
    text-align: left;
    width: 100% !important;
  }

  .ldap-form :deep(.el-form-item__content) {
    margin-left: 0 !important;
  }
}
</style>
