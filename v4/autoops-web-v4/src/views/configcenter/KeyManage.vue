<template>
  <div class="key-management">
    <el-card shadow="hover" class="key-card">
      <template #header>
        <div class="card-header">
          <span class="title">云厂商密钥管理</span>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :model="queryParams" :inline="true" class="search-form" ref="queryForm">
          <el-form-item label="云厂商类型" prop="keyType">
            <el-select placeholder="请选择云厂商类型" clearable size="small" v-model="queryParams.keyType" style="width: 150px">
              <el-option v-for="item in keyTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="Search" size="small" @click="handleQuery">搜索</el-button>
            <el-button type="warning" icon="Refresh" size="small" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    
      <!-- 标签页导航 -->
      <el-tabs v-model="activeTab" class="config-tabs" @tab-change="handleTabChange">
        <el-tab-pane label="密钥管理" name="keys">
          <div class="tab-content">
            <div class="content-header">
              <span class="resource-count">共 {{ keyList.length }} 个密钥配置</span>
              <el-button type="primary" icon="Plus" size="small" v-authority="['config:keymanage:create']" @click="showAddDialog">
                创建密钥
              </el-button>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="定时同步配置" name="schedule">
          <div class="tab-content">
            <div class="content-header">
              <span class="resource-count">共 {{ scheduleList.length }} 个同步任务</span>
              <el-button type="primary" icon="Plus" size="small" @click="showSyncScheduleDialog">
                创建定时同步
              </el-button>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    
      <!-- 密钥列表区域 -->
      <div class="table-section" v-show="activeTab === 'keys'">
        <el-table
          :data="keyList"
          v-loading="loading"
          element-loading-text="加载中..."
          class="resource-table"
          empty-text="暂无密钥配置"
        >
          <el-table-column label="云厂商类型" prop="keyType" width="180" min-width="150">
            <template #default="{ row }">
              <div class="resource-name">
                <img :src="getKeyTypeIcon(row.keyType)" class="cloud-icon"/>
                <el-tag :type="getKeyTypeTag(row.keyType).type">
                  {{ getKeyTypeTag(row.keyType).label }}
                </el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="密钥名称" prop="keyName" min-width="150">
            <template #default="{ row }">
              <span class="key-name-text">{{ row.keyName || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="密钥ID" prop="keyId" min-width="180">
            <template #default="{ row }">
              <div class="key-info">
                <el-icon class="key-icon"><Key /></el-icon>
                <span class="key-value">{{ maskKey(row.keyId) }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="密钥Secret" prop="keySecret" min-width="180">
            <template #default="{ row }">
              <div class="key-info">
                <el-icon class="secret-icon"><Lock /></el-icon>
                <span class="key-value">{{ maskKey(row.keySecret) }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="createdAt" width="180">
            <template #default="{ row }">
              <span class="time-text">{{ row.createdAt }}</span>
            </template>
          </el-table-column>
          <el-table-column label="备注" prop="remark" min-width="150">
            <template #default="{ row }">
              <span class="remark-text">{{ row.remark || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="380" fixed="right">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-tooltip content="查看密钥" placement="top">
                  <el-button size="small" type="warning" v-authority="['config:keymanage:decrypt']" icon="View" @click="handleDecryptKey(row)">查看密钥</el-button>
                </el-tooltip>
                <el-tooltip content="修改" placement="top">
                  <el-button size="small" type="primary" v-authority="['config:keymanage:edit']" icon="Edit" @click="showEditDialog(row)">修改</el-button>
                </el-tooltip>
                <el-tooltip content="同步主机" placement="top">
                  <el-button size="small" type="success" v-authority="['config:keymanage:rsync']" icon="Refresh" @click="handleSyncHosts(row)">同步主机</el-button>
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button size="small" type="danger" v-authority="['config:keymanage:delete']" icon="Delete" @click="handleDelete(row)">删除</el-button>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 定时同步配置列表区域 -->
      <div class="table-section" v-show="activeTab === 'schedule'">
        <el-table
          :data="scheduleList"
          v-loading="scheduleLoading"
          element-loading-text="加载中..."
          class="resource-table"
          empty-text="暂无定时同步配置"
        >
          <el-table-column label="配置名称" prop="name" min-width="160">
            <template #default="{ row }">
              <div class="resource-name">
                <el-icon class="schedule-icon"><Setting /></el-icon>
                <span class="schedule-name">{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="Cron表达式" prop="cronExpr" min-width="150">
            <template #default="{ row }">
              <div class="cron-info">
                <el-icon class="cron-icon"><Timer /></el-icon>
                <code class="cron-text">{{ row.cronExpr }}</code>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="云账号" prop="keyTypes" min-width="200">
            <template #default="{ row }">
              <div class="cloud-accounts">
                <el-tooltip
                  v-for="keyId in parseKeyTypes(row.keyTypes)"
                  :key="keyId"
                  :content="getKeyTagByKeyId(keyId).label"
                  placement="top"
                >
                  <img :src="getKeyTagByKeyId(keyId).icon" class="cloud-account-icon"/>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="状态" prop="status" width="100">
            <template #default="{ row }">
              <el-switch
                v-model="row.status"
                :active-value="1"
                :inactive-value="0"
                active-color="#13ce66"
                inactive-color="#ff4949"
                @change="handleToggleStatus(row)"
              />
            </template>
          </el-table-column>
          <el-table-column label="下次同步时间" prop="nextRunTime" width="180">
            <template #default="{ row }">
              <div class="next-run-time">
                <el-icon class="time-icon"><Timer /></el-icon>
                <span class="time-text">{{ row.nextRunTime || '-' }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="备注" prop="remark" min-width="120">
            <template #default="{ row }">
              <span class="remark-text">{{ row.remark || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="createdAt" width="180">
            <template #default="{ row }">
              <span class="time-text">{{ row.createdAt }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="380" fixed="right">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-tooltip content="立即同步" placement="top">
                  <el-button size="small" type="success" icon="Refresh" @click="handleTriggerSync(row)">立即同步</el-button>
                </el-tooltip>
                <el-tooltip content="同步日志" placement="top">
                  <el-button size="small" type="info" icon="Document" @click="showSyncLog(row)">同步日志</el-button>
                </el-tooltip>
                <el-tooltip content="编辑" placement="top">
                  <el-button size="small" type="primary" icon="Edit" @click="showEditSyncScheduleDialog(row)">编辑</el-button>
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button size="small" type="danger" icon="Delete" @click="handleDeleteSchedule(row)">删除</el-button>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页区域 -->
      <div class="pagination-section" v-show="activeTab === 'keys'">
        <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="queryParams.pageNum"
          :page-sizes="[10, 20, 50, 100]"
          :page-size="queryParams.pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
        />
      </div>

      <!-- 定时同步分页区域 -->
      <div class="pagination-section" v-show="activeTab === 'schedule'">
        <el-pagination
          @size-change="handleScheduleSizeChange"
          @current-change="handleScheduleCurrentChange"
          :current-page="scheduleQueryParams.pageNum"
          :page-sizes="[10, 20, 50, 100]"
          :page-size="scheduleQueryParams.pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="scheduleTotal"
        />
      </div>
    
      <!--新增/编辑对话框-->
      <el-dialog :title="dialogTitle" v-model="dialogVisible" width="35%" :modal="false">
        <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
          <el-form-item label="云厂商类型" prop="keyType">
            <el-select v-model="formData.keyType" placeholder="请选择云厂商类型" style="width: 100%">
              <el-option v-for="item in keyTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="密钥名称" prop="keyName">
            <el-input v-model="formData.keyName" placeholder="请输入密钥名称"></el-input>
          </el-form-item>
          <el-form-item label="密钥ID" prop="keyId">
            <el-input v-model="formData.keyId" placeholder="请输入密钥ID" show-password></el-input>
          </el-form-item>
          <el-form-item label="密钥Secret" prop="keySecret">
            <el-input v-model="formData.keySecret" placeholder="请输入密钥Secret" show-password></el-input>
          </el-form-item>
          <el-form-item label="备注信息" prop="remark">
            <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注信息"></el-input>
          </el-form-item>
        </el-form>
        <span slot="footer" class="dialog-footer">
          <el-button type="primary" @click="submitForm">确 定</el-button>
          <el-button @click="dialogVisible = false">取 消</el-button>
        </span>
      </el-dialog>

      <!--查看密钥对话框-->
      <el-dialog title="查看密钥" v-model="decryptDialogVisible" width="40%" :modal="false">
        <div v-loading="decryptLoading" class="decrypt-container">
          <el-alert
            title="安全提示"
            type="warning"
            :closable="false"
            style="margin-bottom: 16px;"
          >
            请妥善保管密钥信息，不要泄露给他人
          </el-alert>
          <div v-if="decryptedKey" class="decrypt-info">
            <el-form label-width="140px">
              <el-form-item label="云厂商类型">
                <el-tag :type="getKeyTypeTag(decryptedKey.keyType).type">
                  {{ getKeyTypeTag(decryptedKey.keyType).label }}
                </el-tag>
              </el-form-item>
              <el-form-item label="密钥名称">
                <span>{{ decryptedKey.keyName || '-' }}</span>
              </el-form-item>
              <el-form-item label="AccessKey ID">
                <div class="key-display">
                  <el-input
                    v-model="decryptedKey.keyId"
                    readonly
                  >
                    <template #append>
                      <el-button @click="copyToClipboard(decryptedKey.keyId, 'AccessKey ID')" icon="DocumentCopy">
                        复制
                      </el-button>
                    </template>
                  </el-input>
                </div>
              </el-form-item>
              <el-form-item label="AccessKey Secret">
                <div class="key-display">
                  <el-input
                    v-model="decryptedKey.keySecret"
                    readonly
                  >
                    <template #append>
                      <el-button @click="copyToClipboard(decryptedKey.keySecret, 'AccessKey Secret')" icon="DocumentCopy">
                        复制
                      </el-button>
                    </template>
                  </el-input>
                </div>
              </el-form-item>
            </el-form>
          </div>
        </div>
        <template #footer>
          <el-button type="primary" @click="decryptDialogVisible = false">关闭</el-button>
        </template>
      </el-dialog>

      <!--同步日志对话框-->
      <el-dialog title="同步日志" v-model="syncLogDialogVisible" width="50%" :modal="false">
        <div v-loading="syncLogLoading" class="sync-log-container">
          <div v-if="currentSyncLog" class="log-info">
            <div class="log-header">
              <div class="log-header-item">
                <span class="log-label">配置名称:</span>
                <span class="log-value">{{ currentSyncLog.name }}</span>
              </div>
              <div class="log-header-item">
                <span class="log-label">上次执行时间:</span>
                <span class="log-value">{{ currentSyncLog.lastRunTime || '-' }}</span>
              </div>
              <div class="log-header-item">
                <span class="log-label">下次执行时间:</span>
                <span class="log-value">{{ currentSyncLog.nextRunTime || '-' }}</span>
              </div>
            </div>
            <div class="log-content">
              <div class="log-title">同步日志:</div>
              <pre class="log-text">{{ currentSyncLog.syncLog || '暂无日志' }}</pre>
            </div>
          </div>
          <el-empty v-else description="暂无日志数据" />
        </div>
        <template #footer>
          <el-button @click="syncLogDialogVisible = false">关闭</el-button>
        </template>
      </el-dialog>

      <!--定时同步配置对话框-->
      <el-dialog :title="scheduleDialogTitle" v-model="scheduleDialogVisible" width="40%" :modal="false">
        <el-form :model="scheduleFormData" :rules="scheduleFormRules" ref="scheduleFormRef" label-width="120px">
          <el-form-item label="配置名称" prop="name">
            <el-input v-model="scheduleFormData.name" placeholder="请输入配置名称"></el-input>
          </el-form-item>
          <el-form-item label="Cron表达式" prop="cronExpr">
            <el-input
              v-model="scheduleFormData.cronExpr"
              placeholder="请输入Cron表达式"
              @input="handleCronExprChange"
            >
              <template #append>
                <el-popover placement="bottom" width="300" trigger="click">
                  <template #reference>
                    <el-button>模板</el-button>
                  </template>
                  <div class="cron-templates">
                    <div
                      v-for="template in cronTemplates"
                      :key="template.value"
                      class="cron-template-item"
                      @click="selectCronTemplate(template.value)"
                    >
                      <strong>{{ template.label }}</strong>: {{ template.value }}
                    </div>
                  </div>
                </el-popover>
              </template>
            </el-input>

            <!-- Cron表达式预览 - 参考TaskJob.vue的简洁设计 -->
            <div v-if="scheduleFormData.cronExpr && nextExecutionTimes.length > 0" style="margin-top: 8px;">
              <span style="margin-left: 10px; color: #67C23A; font-weight: 500;">
                下次执行时间: {{ nextExecutionTimes[0] }}
              </span>
            </div>
            <div v-else-if="scheduleFormData.cronExpr && cronPreviewError" style="margin-top: 8px;">
              <span style="margin-left: 10px; color: #F56C6C; font-weight: 500;">
                {{ cronPreviewError }}
              </span>
            </div>
          </el-form-item>
          <el-form-item label="云账号选择" prop="keyTypes">
            <el-select
              v-model="scheduleFormData.keyTypes"
              placeholder="请选择云账号"
              multiple
              collapse-tags
              collapse-tags-tooltip
              style="width: 100%"
            >
              <el-option
                v-for="key in keyList"
                :key="key.id"
                :label="`${getKeyTypeTag(key.keyType).label} - ${key.keyName || key.keyId}`"
                :value="key.id"
              >
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img :src="getKeyTypeIcon(key.keyType)" style="width: 16px; height: 16px;"/>
                  <el-tag :type="getKeyTypeTag(key.keyType).type" size="small">
                    {{ getKeyTypeTag(key.keyType).label }}
                  </el-tag>
                  <span>{{ key.keyName || maskKey(key.keyId) }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="启用状态" prop="status">
            <el-switch
              v-model="scheduleFormData.status"
              :active-value="1"
              :inactive-value="0"
              active-text="启用"
              inactive-text="禁用"
            />
          </el-form-item>
          <el-form-item label="备注信息" prop="remark">
            <el-input
              v-model="scheduleFormData.remark"
              type="textarea"
              :rows="3"
              placeholder="请输入备注信息"
            ></el-input>
          </el-form-item>
        </el-form>
        <span slot="footer" class="dialog-footer">
          <el-button type="primary" @click="submitScheduleForm">确 定</el-button>
          <el-button @click="scheduleDialogVisible = false">取 消</el-button>
        </span>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import { Plus, Edit, Delete, Refresh, Setting, Key, Lock, Timer, Loading, Warning, Document } from '@element-plus/icons-vue'
import API from '@/api/config'

export default {
  data() {
    return {
      activeTab: 'keys',
      queryParams: {
        keyType: '',
        pageNum: 1,
        pageSize: 10
      },
      loading: false,
      keyList: [],
      total: 0,
      dialogVisible: false,
      dialogTitle: '',
      scheduleQueryParams: {
        pageNum: 1,
        pageSize: 10
      },
      scheduleLoading: false,
      scheduleList: [],
      scheduleTotal: 0,
      scheduleDialogVisible: false,
      scheduleDialogTitle: '',
      syncLogDialogVisible: false,
      syncLogLoading: false,
      currentSyncLog: null,
      decryptDialogVisible: false,
      decryptLoading: false,
      decryptedKey: null,
      keyTypeOptions: [
        { value: 1, label: '阿里云' },
        { value: 2, label: '腾讯云' },
        { value: 3, label: '百度云' },
        { value: 4, label: '华为云' },
        { value: 5, label: 'AWS云' },
        { value: 6, label: '火山云' },
        { value: 7, label: '飞书应用' }
      ],
      cronTemplates: [
        { label: '每天凌晨2点', value: '0 0 2 * * ?' },
        { label: '每6小时', value: '0 0 */6 * * ?' },
        { label: '每天上午9点', value: '0 0 9 * * ?' },
        { label: '工作日下午6点', value: '0 0 18 * * MON-FRI' },
        { label: '每周日凌晨3点', value: '0 0 3 * * SUN' }
      ],
      formData: {
        id: '',
        keyType: '',
        keyName: '',
        keyId: '',
        keySecret: '',
        remark: ''
      },
      formRules: {
        keyType: [{ required: true, message: '请选择云厂商类型', trigger: 'change' }],
        keyId: [{ required: true, message: '请输入密钥ID', trigger: 'blur' }],
        keySecret: [{ required: true, message: '请输入密钥Secret', trigger: 'blur' }]
      },
      scheduleFormData: {
        id: '',
        name: '',
        cronExpr: '',
        keyTypes: [],
        status: 1,
        remark: ''
      },
      scheduleFormRules: {
        name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
        cronExpr: [{ required: true, message: '请输入Cron表达式', trigger: 'blur' }],
        keyTypes: [{ required: true, message: '请选择云账号', trigger: 'change' }]
      },
      cronPreviewLoading: false,
      cronPreviewError: '',
      nextExecutionTimes: [],
      cronPreviewTimer: null
    }
  },
  methods: {
    // 获取密钥列表
    async getList() {
      this.loading = true
      try {
        const { data: res } = await API.getKeyManageList({
          page: this.queryParams.pageNum,
          pageSize: this.queryParams.pageSize,
          keyType: this.queryParams.keyType || undefined
        })
        if (res.code === 200) {
          this.keyList = res.data?.list || []
          this.total = res.data?.total || 0
        } else {
          this.$message.error(res.message || '获取密钥列表失败')
        }
      } catch (error) {
        console.error('获取密钥列表失败:', error)
        this.$message.error('获取密钥列表失败')
      } finally {
        this.loading = false
      }
    },
    
    // 搜索
    async handleQuery() {
      this.queryParams.pageNum = 1
      this.getList()
    },
    
    // 重置搜索
    resetQuery() {
      this.queryParams = {
        keyType: '',
        pageNum: 1,
        pageSize: 10
      }
      this.getList()
    },
    
    // 分页处理
    handleSizeChange(val) {
      this.queryParams.pageSize = val
      this.queryParams.pageNum = 1
      this.getList()
    },
    
    handleCurrentChange(val) {
      this.queryParams.pageNum = val
      this.getList()
    },
    
    // 显示新增对话框
    showAddDialog() {
      this.dialogTitle = '创建密钥'
      this.$nextTick(() => {
        this.formData = {
          id: '',
          keyType: '',
          keyName: '',
          keyId: '',
          keySecret: '',
          remark: ''
        }
        this.dialogVisible = true
      })
    },
    
    // 显示编辑对话框
    showEditDialog(row) {
      this.dialogTitle = '修改密钥'
      this.$nextTick(() => {
        this.formData = {
          id: row.id,
          keyType: row.keyType,
          keyName: row.keyName || '',
          keyId: row.keyId,
          keySecret: row.keySecret,
          remark: row.remark || ''
        }
        this.dialogVisible = true
      })
    },
    
    // 提交表单
    async submitForm() {
      try {
        await this.$refs.formRef.validate()
        
        const formData = { ...this.formData }
        
        let res
        if (formData.id) {
          // 更新
          res = await API.updateKeyManage(formData)
        } else {
          // 新增
          res = await API.createKeyManage(formData)
        }

        if (res.data.code === 200) {
          this.$message.success(formData.id ? '修改成功' : '创建成功')
          this.dialogVisible = false
          await this.getList()
        } else {
          this.$message.error(res.data.message || (formData.id ? '修改失败' : '创建失败'))
        }
      } catch (error) {
        console.error('操作失败:', error)
        this.$message.error('操作失败: ' + error.message)
      }
    },
    
    // 删除密钥
    async handleDelete(row) {
      try {
        await this.$confirm(`确定删除该密钥配置?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        const { data: res } = await API.deleteKeyManage(row.id)
        if (res.code === 200) {
          this.$message.success('删除成功')
          this.getList()
        } else {
          this.$message.error(res.message || '删除失败')
        }
      } catch (error) {
        console.error('删除失败:', error)
      }
    },

    // 查看解密后的密钥
    async handleDecryptKey(row) {
      this.decryptDialogVisible = true
      this.decryptLoading = true
      this.decryptedKey = null

      try {
        const { data: res } = await API.decryptKeys(row.id)
        if (res.code === 200 && res.data) {
          this.decryptedKey = {
            keyType: row.keyType,
            keyName: row.keyName,
            keyId: res.data.accessKeyId || res.data.keyId,
            keySecret: res.data.accessKeySecret || res.data.keySecret
          }
        } else {
          this.$message.error(res.message || '解密失败')
          this.decryptDialogVisible = false
        }
      } catch (error) {
        console.error('解密失败:', error)
        this.$message.error('解密失败')
        this.decryptDialogVisible = false
      } finally {
        this.decryptLoading = false
      }
    },

    // 复制到剪贴板
    copyToClipboard(text, label) {
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text).then(() => {
          this.$message.success(`${label} 已复制到剪贴板`)
        }).catch(err => {
          console.error('复制失败:', err)
          this.fallbackCopy(text, label)
        })
      } else {
        this.fallbackCopy(text, label)
      }
    },

    // 降级复制方法
    fallbackCopy(text, label) {
      const textArea = document.createElement('textarea')
      textArea.value = text
      textArea.style.position = 'fixed'
      textArea.style.left = '-999999px'
      document.body.appendChild(textArea)
      textArea.select()
      try {
        document.execCommand('copy')
        this.$message.success(`${label} 已复制到剪贴板`)
      } catch (err) {
        console.error('复制失败:', err)
        this.$message.error('复制失败，请手动复制')
      }
      document.body.removeChild(textArea)
    },

    // 同步主机
    async handleSyncHosts(row) {
      try {
        await this.$confirm(`确定同步${this.getKeyTypeTag(row.keyType).label}主机?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info'
        })

        const syncData = {
          keyId: row.id,
          keyType: row.keyType
        }

        const { data: res } = await API.syncCloudHosts(syncData)

        if (res.code === 200) {
          this.$message.success('同步成功')
        } else {
          this.$message.error(res.message || '同步失败')
        }
      } catch (error) {
        console.error('同步失败:', error)
        this.$message.error('同步失败')
      }
    },
    
    // 获取云厂商标签样式
    getKeyTypeTag(keyType) {
      const typeMap = {
        1: { label: '阿里云', type: 'primary' },
        2: { label: '腾讯云', type: 'success' },
        3: { label: '百度云', type: 'info' },
        4: { label: '华为云', type: 'warning' },
        5: { label: 'AWS云', type: 'danger' },
        6: { label: '火山云', type: 'danger' },
        7: { label: '飞书应用', type: 'primary' }
      }
      return typeMap[keyType] || { label: '未知', type: 'info' }
    },
    
    // 获取云厂商图标
    getKeyTypeIcon(keyType) {
      const iconMap = {
        1: require('@/assets/image/aliyun.png'),      // 阿里云
        2: require('@/assets/image/tengxun.png'),     // 腾讯云
        3: require('@/assets/image/baidu.svg'),       // 百度云
        4: require('@/assets/image/huaweiyun.svg'),   // 华为云
        5: require('@/assets/image/aws.svg'),         // AWS云
        6: require('@/assets/image/火山引擎.svg'),      // 火山云
        7: require('@/assets/image/飞书.svg')          // 飞书应用
      }
      return iconMap[keyType] || require('@/assets/image/版本.svg')
    },
    
    // 密钥掩码显示
    maskKey(key) {
      if (!key) return ''
      if (key.length <= 8) return key
      return key.substring(0, 4) + '****' + key.substring(key.length - 4)
    },

    // 标签页切换
    handleTabChange(tabName) {
      if (tabName === 'schedule') {
        this.getScheduleList()
      }
    },

    // 获取定时同步配置列表
    async getScheduleList() {
      this.scheduleLoading = true
      try {
        const { data: res } = await API.getSyncScheduleList({
          page: this.scheduleQueryParams.pageNum,
          pageSize: this.scheduleQueryParams.pageSize
        })
        if (res.code === 200) {
          this.scheduleList = res.data?.list || []
          this.scheduleTotal = res.data?.total || 0
        } else {
          this.$message.error(res.message || '获取定时同步配置列表失败')
        }
      } catch (error) {
        console.error('获取定时同步配置列表失败:', error)
        this.$message.error('获取定时同步配置列表失败')
      } finally {
        this.scheduleLoading = false
      }
    },

    // 定时同步分页处理
    handleScheduleSizeChange(val) {
      this.scheduleQueryParams.pageSize = val
      this.scheduleQueryParams.pageNum = 1
      this.getScheduleList()
    },

    handleScheduleCurrentChange(val) {
      this.scheduleQueryParams.pageNum = val
      this.getScheduleList()
    },

    // 显示新增定时同步对话框
    showSyncScheduleDialog() {
      this.scheduleDialogTitle = '创建定时同步配置'
      // 确保有密钥列表
      if (this.keyList.length === 0) {
        this.getList()
      }
      this.$nextTick(() => {
        this.scheduleFormData = {
          id: '',
          name: '',
          cronExpr: '',
          keyTypes: [],
          status: 1,
          remark: ''
        }
        // 重置预览状态
        this.cronPreviewError = ''
        this.nextExecutionTimes = []
        this.cronPreviewLoading = false
        this.scheduleDialogVisible = true
      })
    },

    // 显示编辑定时同步对话框
    showEditSyncScheduleDialog(row) {
      this.scheduleDialogTitle = '编辑定时同步配置'
      // 确保有密钥列表
      if (this.keyList.length === 0) {
        this.getList()
      }
      this.$nextTick(() => {
        this.scheduleFormData = {
          id: row.id,
          name: row.name,
          cronExpr: row.cronExpr,
          keyTypes: this.parseKeyTypes(row.keyTypes),
          status: row.status,
          remark: row.remark || ''
        }
        // 重置预览状态
        this.cronPreviewError = ''
        this.nextExecutionTimes = []
        this.cronPreviewLoading = false
        this.scheduleDialogVisible = true

        // 如果有现有的Cron表达式，立即验证
        if (row.cronExpr) {
          this.handleCronExprChange(row.cronExpr)
        }
      })
    },

    // 选择Cron模板
    selectCronTemplate(cronExpr) {
      this.scheduleFormData.cronExpr = cronExpr
      this.handleCronExprChange(cronExpr)
    },

    // 处理Cron表达式变化
    handleCronExprChange(cronExpr) {
      // 清除之前的定时器
      if (this.cronPreviewTimer) {
        clearTimeout(this.cronPreviewTimer)
      }

      // 重置状态
      this.cronPreviewError = ''
      this.nextExecutionTimes = []

      // 如果表达式为空，不进行验证
      if (!cronExpr || cronExpr.trim() === '') {
        return
      }

      // 防抖处理，500ms 后才执行验证
      this.cronPreviewTimer = setTimeout(() => {
        this.validateCronExpression(cronExpr.trim())
      }, 500)
    },

    // 验证Cron表达式并获取执行时间
    async validateCronExpression(cronExpr) {
      this.cronPreviewLoading = true
      this.cronPreviewError = ''

      try {
        const response = await API.getNextExecutionTime(cronExpr)
        const responseData = response?.data || {}

        if (responseData.code === 200) {
          // 参考 TaskJob.vue 的处理方式
          if (responseData.data?.next_execution_time) {
            // 单个时间的处理
            this.nextExecutionTimes = [this.formatTime(responseData.data.next_execution_time)]
          } else if (Array.isArray(responseData.data)) {
            // 时间数组的处理
            this.nextExecutionTimes = responseData.data.map(time => this.formatTime(time))
          } else {
            this.nextExecutionTimes = []
          }

          if (this.nextExecutionTimes.length === 0) {
            this.cronPreviewError = 'Cron表达式有效，但无法计算执行时间'
          }
        } else {
          this.cronPreviewError = responseData.message || 'Cron表达式格式错误'
        }
      } catch (error) {
        console.error('验证Cron表达式失败:', error)
        this.cronPreviewError = 'Cron表达式格式错误，请检查语法'
      } finally {
        this.cronPreviewLoading = false
      }
    },

    // 格式化时间
    formatTime(timeStr) {
      if (!timeStr) return ''
      const date = new Date(timeStr)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      }).replace(/\//g, '-')
    },

    // 提交定时同步表单
    async submitScheduleForm() {
      try {
        await this.$refs.scheduleFormRef.validate()

        const formData = {
          ...this.scheduleFormData,
          keyTypes: JSON.stringify(this.scheduleFormData.keyTypes)
        }

        let res
        if (formData.id) {
          res = await API.updateSyncSchedule(formData)
        } else {
          res = await API.createSyncSchedule(formData)
        }

        if (res.data.code === 200) {
          this.$message.success(formData.id ? '更新成功' : '创建成功')
          this.scheduleDialogVisible = false
          await this.getScheduleList()
        } else {
          this.$message.error(res.data.message || (formData.id ? '更新失败' : '创建失败'))
        }
      } catch (error) {
        console.error('操作失败:', error)
        this.$message.error('操作失败: ' + error.message)
      }
    },

    // 切换定时同步状态
    async handleToggleStatus(row) {
      try {
        const { data: res } = await API.toggleSyncScheduleStatus({
          id: row.id,
          status: row.status
        })

        if (res.code === 200) {
          this.$message.success(`${row.status === 1 ? '启用' : '禁用'}成功`)
        } else {
          row.status = row.status === 1 ? 0 : 1
          this.$message.error(res.message || '状态切换失败')
        }
      } catch (error) {
        row.status = row.status === 1 ? 0 : 1
        console.error('状态切换失败:', error)
        this.$message.error('状态切换失败')
      }
    },

    // 删除定时同步配置
    async handleDeleteSchedule(row) {
      try {
        await this.$confirm(`确定删除定时同步配置"${row.name}"?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        const { data: res } = await API.deleteSyncSchedule(row.id)
        if (res.code === 200) {
          this.$message.success('删除成功')
          this.getScheduleList()
        } else {
          this.$message.error(res.message || '删除失败')
        }
      } catch (error) {
        console.error('删除失败:', error)
      }
    },

    // 立即触发同步任务
    async handleTriggerSync(row) {
      try {
        await this.$confirm(`确定立即执行同步任务"${row.name}"?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info'
        })

        const { data: res } = await API.triggerSyncSchedule(row.id)
        if (res.code === 200) {
          this.$message.success(res.message || '同步任务已触发，正在后台执行')
        } else {
          this.$message.error(res.message || '触发同步任务失败')
        }
      } catch (error) {
        if (error !== 'cancel') {
          console.error('触发同步任务失败:', error)
          this.$message.error('触发同步任务失败')
        }
      }
    },

    // 解析keyTypes字符串(实际存储的是keyId数组)
    parseKeyTypes(keyTypesStr) {
      if (!keyTypesStr) return []
      try {
        return JSON.parse(keyTypesStr)
      } catch (error) {
        console.error('解析keyTypes失败:', error)
        return []
      }
    },

    // 根据keyId获取密钥信息标签
    getKeyTagByKeyId(keyId) {
      const key = this.keyList.find(k => k.id === keyId)
      if (!key) {
        return { label: '未知账号', type: 'info', icon: require('@/assets/image/版本.svg') }
      }
      return {
        label: key.keyName || this.maskKey(key.keyId),
        type: this.getKeyTypeTag(key.keyType).type,
        icon: this.getKeyTypeIcon(key.keyType)
      }
    },

    // 显示同步日志
    async showSyncLog(row) {
      this.syncLogDialogVisible = true
      this.syncLogLoading = true
      this.currentSyncLog = null

      try {
        const { data: res } = await API.getSyncScheduleLog(row.id)
        if (res.code === 200 && res.data) {
          this.currentSyncLog = {
            id: res.data.id,
            name: res.data.name,
            syncLog: res.data.syncLog,
            lastRunTime: res.data.lastRunTime ? this.formatTime(res.data.lastRunTime) : '-',
            nextRunTime: res.data.nextRunTime ? this.formatTime(res.data.nextRunTime) : '-'
          }
        } else {
          this.$message.warning(res.message || '暂无日志数据')
        }
      } catch (error) {
        console.error('获取同步日志失败:', error)
        this.$message.error('获取同步日志失败')
      } finally {
        this.syncLogLoading = false
      }
    }
  },
  created() {
    this.getList()
  },
  beforeUnmount() {
    // 清理定时器
    if (this.cronPreviewTimer) {
      clearTimeout(this.cronPreviewTimer)
    }
  }
}
</script>

<style scoped>
.key-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.key-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  height: 800px;
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
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.search-section {
  margin-bottom: 20px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form {
  margin: 0;
}

.search-form .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

/* 按钮样式 - 与k8s-config.vue保持一致 */
.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 输入框样式 - 与k8s-config.vue保持一致 */
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

/* 标签样式 */
.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

.config-tabs {
  margin-top: 20px;
}

.config-tabs :deep(.el-tabs__header) {
  margin-bottom: 20px;
}

.config-tabs :deep(.el-tabs__item) {
  font-weight: 500;
  color: #606266;
}

.config-tabs :deep(.el-tabs__item.is-active) {
  color: #409EFF;
  font-weight: 600;
}

.tab-content {
  padding: 0;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px 0;
}

.resource-count {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.resource-table {
  border-radius: 8px;
  overflow: hidden;
}

.resource-table :deep(.el-table__header) {
  background: #f8f9fa;
}

.resource-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

.resource-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cloud-icon {
  width: 16px;
  height: 16px;
}

.key-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.key-icon {
  color: #409EFF;
  font-size: 16px;
}

.secret-icon {
  color: #E6A23C;
  font-size: 16px;
}

.key-value {
  color: #606266;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
}

.time-text {
  color: #909399;
  font-size: 13px;
}

.remark-text {
  color: #606266;
  font-size: 13px;
}

.schedule-icon {
  color: #409EFF;
  font-size: 16px;
}

.schedule-name {
  color: #409EFF;
  cursor: pointer;
  font-weight: 500;
  transition: color 0.3s;
}

.cron-info {
  display: flex;
  align-items: center;
  gap: 6px;
}

.cron-icon {
  color: #67C23A;
  font-size: 14px;
}

.cron-text {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  color: #E6A23C;
  border: 1px solid #e9ecef;
}

.cloud-types,
.cloud-accounts {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.cloud-account-icon {
  width: 24px;
  height: 24px;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.cloud-account-icon:hover {
  transform: scale(1.2);
}

.cloud-tag {
  font-size: 11px;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.next-run-time {
  display: flex;
  align-items: center;
  gap: 6px;
}

.time-icon {
  color: #67C23A;
  font-size: 14px;
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

.el-pagination {
  background: transparent;
}

.el-pagination :deep(.btn-next),
.el-pagination :deep(.btn-prev),
.el-pagination :deep(.el-pager li) {
  border-radius: 6px;
  margin: 0 2px;
  transition: all 0.3s ease;
}

.el-pagination :deep(.btn-next):hover,
.el-pagination :deep(.btn-prev):hover,
.el-pagination :deep(.el-pager li):hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.cron-templates {
  max-height: 200px;
  overflow-y: auto;
}

.cron-template-item {
  padding: 8px 12px;
  margin-bottom: 8px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.cron-template-item:hover {
  background: rgba(103, 126, 234, 0.1);
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.cron-template-item:last-child {
  margin-bottom: 0;
}

.cron-template-item strong {
  color: #667eea;
  display: block;
  margin-bottom: 2px;
}

.el-checkbox-group .el-checkbox {
  margin-right: 20px;
  margin-bottom: 10px;
}

.el-checkbox-group .el-checkbox:deep(.el-checkbox__label) {
  font-weight: 500;
}


/* 加载动画样式 */
.el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

/* 解密密钥对话框样式 */
.decrypt-container {
  min-height: 200px;
}

.decrypt-info {
  padding: 16px 0;
}

.key-display {
  width: 100%;
}

.key-display .el-input {
  font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
}

.key-display .el-input :deep(.el-input__inner) {
  font-size: 13px;
  font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
}

.key-display :deep(.el-input-group__append) {
  padding: 0;
}

.key-display :deep(.el-input-group__append .el-button) {
  margin: 0;
  border-left: 1px solid #dcdfe6;
}

/* 同步日志对话框样式 */
.sync-log-container {
  min-height: 200px;
}

.log-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.log-header {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.log-header-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.log-label {
  font-weight: 600;
  color: #606266;
  min-width: 100px;
}

.log-value {
  color: #409EFF;
  font-weight: 500;
}

.log-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.log-title {
  font-weight: 600;
  color: #606266;
  font-size: 14px;
}

.log-text {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #2c3e50;
  max-height: 400px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
  margin: 0;
}

.log-text::-webkit-scrollbar {
  width: 8px;
}

.log-text::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.log-text::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

.log-text::-webkit-scrollbar-thumb:hover {
  background: #555;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .key-management {
    padding: 12px;
  }

  .card-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .content-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .resource-table :deep(.el-table) {
    font-size: 12px;
  }

  .operation-buttons {
    flex-direction: column;
    gap: 4px;
  }

  .search-section {
    padding: 12px;
  }

  .search-form {
    flex-direction: column;
    align-items: stretch;
  }

  .search-form :deep(.el-form-item) {
    margin-bottom: 8px;
  }

  .cloud-types {
    gap: 2px;
  }

  .cloud-tag {
    max-width: 60px;
  }

  .log-header {
    padding: 12px;
  }

  .log-label {
    min-width: 80px;
    font-size: 13px;
  }

  .log-text {
    font-size: 12px;
    padding: 12px;
  }
}
</style>
