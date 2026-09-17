<template>
  <div class="notify-robot-manage">
    <div class="form-title">
      <el-icon><Bell /></el-icon>
      <span>报警推送机器人配置</span>
    </div>

    <!-- 机器人列表操作按钮 -->
    <div class="robot-actions" style="margin-bottom: 15px;">
      <el-button type="primary" size="small" @click="handleAddRobot">
        <el-icon><Plus /></el-icon>
        新增机器人
      </el-button>
      <el-button size="small" @click="handleRefreshRobots">
        <el-icon><RefreshLeft /></el-icon>
        刷新
      </el-button>
      <el-button type="info" size="small" @click="handleShowCallback">
        <el-icon><Document /></el-icon>
        回调示例
      </el-button>
    </div>

    <!-- 机器人列表 -->
    <el-table
      :data="robotList"
      v-loading="robotLoading"
      border
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="机器人名称" min-width="150" />
      <el-table-column prop="type" label="类型" width="150">
        <template #default="{ row }">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img :src="getRobotIcon(row.type)" :alt="getTypeLabel(row.type)" style="width: 20px; height: 20px;" />
            <el-tag :type="getTypeTagType(row.type)">
              {{ getTypeLabel(row.type) }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="webhook" label="Webhook地址" min-width="250" show-overflow-tooltip />
      <el-table-column prop="secret" label="加签密钥" min-width="150" show-overflow-tooltip>
        <template #default="{ row }">
          <span v-if="row.secret">{{ row.secret.substring(0, 10) }}...</span>
          <span v-else>-</span>
        </template>
      </el-table-column>
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
      <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
      <el-table-column prop="createdAt" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDateTime(row.createdAt) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="240" fixed="right">
        <template #default="{ row }">
          <el-button type="success" link size="small" @click="handleTestRobot(row)">
            测试
          </el-button>
          <el-button type="primary" link size="small" @click="handleEditRobot(row)">
            编辑
          </el-button>
          <el-button type="danger" link size="small" @click="handleDeleteRobot(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      style="margin-top: 20px; justify-content: flex-end;"
      @size-change="handlePageChange"
      @current-change="handlePageChange"
    />

    <!-- 回调示例对话框 -->
    <el-dialog
      v-model="callbackDialogVisible"
      title="Webhook 回调配置示例"
      width="850px"
      :close-on-click-modal="false"
    >
      <div class="callback-guide">
        <div class="guide-intro">
          <p>本系统支持通过 Webhook 接收外部告警通知，以下是配置示例：</p>
        </div>

        <el-tabs v-model="activeCallbackTab">
          <el-tab-pane label="X-Webhook-Token 方式" name="header">
            <div class="callback-section">
              <h4>1. 接口信息</h4>
              <div class="info-box">
                <p><strong>接口地址：</strong><code>http://10.7.16.100:8080/api/v1/monitor/webhook/receive</code></p>
                <p><strong>请求方式：</strong><code>POST</code></p>
                <p><strong>认证方式：</strong>在 Header 中添加 <code>X-Webhook-Token</code></p>
              </div>

              <h4>2. 请求示例</h4>
              <el-input
                v-model="callbackExample1"
                type="textarea"
                :rows="12"
                readonly
                class="callback-code"
              />
              <el-button
                type="primary"
                size="small"
                style="margin-top: 10px"
                @click="copyToClipboard(callbackExample1)"
              >
                复制 curl 命令
              </el-button>

              <h4>3. 请求参数说明</h4>
              <el-table :data="parameterData" border size="small" style="margin-top: 10px">
                <el-table-column prop="field" label="字段" width="120" />
                <el-table-column prop="type" label="类型" width="80" />
                <el-table-column prop="required" label="必填" width="80" />
                <el-table-column prop="description" label="说明" />
              </el-table>

              <h4>4. 响应示例</h4>
              <el-input
                v-model="responseExample"
                type="textarea"
                :rows="8"
                readonly
                class="callback-code"
              />
            </div>
          </el-tab-pane>

          <el-tab-pane label="Authorization Bearer 方式" name="bearer">
            <div class="callback-section">
              <h4>1. 接口信息</h4>
              <div class="info-box">
                <p><strong>接口地址：</strong><code>http://10.7.16.100:8080/api/v1/monitor/webhook/receive</code></p>
                <p><strong>请求方式：</strong><code>POST</code></p>
                <p><strong>认证方式：</strong>在 Header 中添加 <code>Authorization: Bearer {token}</code></p>
              </div>

              <h4>2. 请求示例</h4>
              <el-input
                v-model="callbackExample2"
                type="textarea"
                :rows="12"
                readonly
                class="callback-code"
              />
              <el-button
                type="primary"
                size="small"
                style="margin-top: 10px"
                @click="copyToClipboard(callbackExample2)"
              >
                复制 curl 命令
              </el-button>

              <h4>3. 请求参数说明</h4>
              <el-table :data="parameterData" border size="small" style="margin-top: 10px">
                <el-table-column prop="field" label="字段" width="120" />
                <el-table-column prop="type" label="类型" width="80" />
                <el-table-column prop="required" label="必填" width="80" />
                <el-table-column prop="description" label="说明" />
              </el-table>

              <h4>4. 响应示例</h4>
              <el-input
                v-model="responseExample"
                type="textarea"
                :rows="8"
                readonly
                class="callback-code"
              />
            </div>
          </el-tab-pane>
        </el-tabs>

        <div class="guide-tips">
          <h4>📌 使用说明</h4>
          <ul>
            <li>两种认证方式任选其一即可，效果相同</li>
            <li>Token 需要根据实际系统配置进行替换（示例中为 <code>webhook-notify-token-2024</code>）</li>
            <li>告警级别 <code>level</code> 支持：<code>info</code>、<code>warning</code>、<code>error</code>、<code>critical</code></li>
            <li>来源 <code>source</code> 支持任意自定义值，如：zabbix、prometheus、grafana 等</li>
            <li>系统会将收到的告警通知推送到所有已启用的机器人</li>
          </ul>
        </div>
      </div>
    </el-dialog>

    <!-- 新增/编辑机器人对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="robotForm.id ? '编辑机器人' : '新增机器人'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="robotFormRef"
        :model="robotForm"
        :rules="robotRules"
        label-width="120px"
      >
        <el-form-item label="机器人名称" prop="name">
          <el-input v-model="robotForm.name" placeholder="例如：钉钉告警机器人" clearable />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select
            v-model="robotForm.type"
            placeholder="请选择机器人类型"
            style="width: 100%"
            @change="handleTypeChange"
          >
            <el-option value="dingtalk" label="钉钉机器人">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img src="@/assets/image/钉钉01.svg" alt="钉钉" style="width: 20px; height: 20px;" />
                <span>钉钉机器人</span>
              </div>
            </el-option>
            <el-option value="feishu" label="飞书机器人">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img src="@/assets/image/飞书.svg" alt="飞书" style="width: 20px; height: 20px;" />
                <span>飞书机器人</span>
              </div>
            </el-option>
            <el-option value="wechat" label="企业微信机器人">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img src="@/assets/image/企业微信.svg" alt="企业微信" style="width: 20px; height: 20px;" />
                <span>企业微信机器人</span>
              </div>
            </el-option>
            <el-option value="email" label="邮箱通知">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img src="@/assets/image/邮箱.svg" alt="邮箱" style="width: 20px; height: 20px;" />
                <span>邮箱通知</span>
              </div>
            </el-option>
            <el-option value="webhook" label="Webhook通知">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img src="@/assets/image/url.svg" alt="Webhook" style="width: 20px; height: 20px;" />
                <span>Webhook通知</span>
              </div>
            </el-option>
            <el-option value="teams" label="Microsoft Teams">
              <div style="display: flex; align-items: center; gap: 8px;">
                <img src="@/assets/image/通知.svg" alt="Teams" style="width: 20px; height: 20px;" />
                <span>Microsoft Teams</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-if="robotForm.type !== 'email'" label="Webhook地址" prop="webhook">
          <el-input
            v-model="robotForm.webhook"
            type="textarea"
            :rows="3"
            placeholder="请输入Webhook地址"
            clearable
          />
          <div style="margin-top: 8px; color: #909399; font-size: 12px;">
            <template v-if="robotForm.type === 'dingtalk'">
              钉钉：https://oapi.dingtalk.com/robot/send?access_token=xxx
            </template>
            <template v-else-if="robotForm.type === 'feishu'">
              飞书：https://open.feishu.cn/open-apis/bot/v2/hook/xxx
            </template>
            <template v-else-if="robotForm.type === 'wechat'">
              企业微信：https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx
            </template>
            <template v-else-if="robotForm.type === 'teams'">
              Teams：https://outlook.office.com/webhook/xxx
            </template>
            <template v-else-if="robotForm.type === 'webhook'">
              自定义Webhook地址
            </template>
          </div>
        </el-form-item>
        <el-form-item
          v-if="robotForm.type === 'dingtalk'"
          label="加签密钥"
          prop="secret"
        >
          <el-input
            v-model="robotForm.secret"
            type="textarea"
            :rows="2"
            placeholder="请输入加签密钥（可选，钉钉机器人安全设置中的加签）"
            clearable
          />
          <div style="margin-top: 8px; color: #909399; font-size: 12px;">
            加签密钥格式：SECxxx...（可选，若钉钉机器人启用了加签安全设置，则需填写）
          </div>
        </el-form-item>
        <!-- 邮箱配置字段 -->
        <template v-if="robotForm.type === 'email'">
          <el-form-item label="SMTP服务器" prop="server">
            <el-input
              v-model="robotForm.server"
              placeholder="例如：smtp.qq.com"
              clearable
            />
          </el-form-item>
          <el-form-item label="SMTP端口" prop="port">
            <el-input-number
              v-model="robotForm.port"
              :min="1"
              :max="65535"
              placeholder="465"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="发件人用户名" prop="username">
            <el-input
              v-model="robotForm.username"
              placeholder="请输入发件人-用户名"
              clearable
            />
          </el-form-item>
          <el-form-item label="授权码/密码" prop="password">
            <el-input
              v-model="robotForm.password"
              type="password"
              placeholder="请输入邮箱授权码或密码"
              show-password
              clearable
            />
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              QQ邮箱等需要使用授权码，而非登录密码
            </div>
          </el-form-item>
          <el-form-item label="发件人昵称" prop="nickname">
            <el-input
              v-model="robotForm.nickname"
              placeholder="例如：DoDevOps告警"
              clearable
            />
          </el-form-item>
        </template>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="robotForm.status" :active-value="1" :inactive-value="0" />
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">
            {{ robotForm.status === 1 ? '启用' : '禁用' }}
          </span>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="robotForm.remark"
            type="textarea"
            :rows="2"
            placeholder="请输入备注说明"
            clearable
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveRobot" :loading="saving">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { Bell, Plus, RefreshLeft, Document } from '@element-plus/icons-vue'
import monitorApi from '@/api/monitor'
import { ElMessage, ElMessageBox } from 'element-plus'
import dingtalkIcon from '@/assets/image/钉钉01.svg'
import feishuIcon from '@/assets/image/飞书.svg'
import wechatIcon from '@/assets/image/企业微信.svg'
import emailIcon from '@/assets/image/邮箱.svg'
import webhookIcon from '@/assets/image/url.svg'
import teamsIcon from '@/assets/image/通知.svg'

export default {
  name: 'NotifyRobotManage',
  components: {
    Bell,
    Plus,
    RefreshLeft,
    Document
  },
  data() {
    return {
      robotList: [],
      robotLoading: false,
      pagination: {
        page: 1,
        pageSize: 10,
        total: 0
      },
      dialogVisible: false,
      saving: false,
      callbackDialogVisible: false,
      activeCallbackTab: 'header',
      currentRobot: null,
      callbackExample1: '',
      callbackExample2: '',
      responseExample: '',
      parameterData: [
        { field: 'source', type: 'string', required: '是', description: '告警来源，如：zabbix、prometheus、grafana 等' },
        { field: 'title', type: 'string', required: '是', description: '告警标题' },
        { field: 'content', type: 'string', required: '是', description: '告警内容详情' },
        { field: 'level', type: 'string', required: '是', description: '告警级别：info、warning、error、critical' }
      ],
      robotForm: {
        id: null,
        name: '',
        type: 'feishu',
        webhook: '',
        secret: '',
        // 邮箱专用字段
        server: '',
        port: '',
        username: '',
        password: '',
        nickname: '',
        status: 1,
        remark: ''
      },
      robotRules: {
        name: [
          { required: true, message: '请输入机器人名称', trigger: 'blur' },
          { max: 100, message: '机器人名称最多100个字符', trigger: 'blur' }
        ],
        type: [
          { required: true, message: '请选择机器人类型', trigger: 'change' }
        ],
        webhook: [
          {
            validator: (rule, value, callback) => {
              // 邮箱类型不需要webhook
              if (this.robotForm.type === 'email') {
                callback()
                return
              }
              // 其他类型必填
              if (!value) {
                callback(new Error('请输入Webhook地址'))
                return
              }
              // 验证URL格式
              if (!value.startsWith('http')) {
                callback(new Error('Webhook地址必须以http或https开头'))
              } else {
                callback()
              }
            },
            trigger: 'blur'
          },
          { max: 500, message: 'Webhook地址最多500个字符', trigger: 'blur' }
        ],
        secret: [
          { max: 200, message: '加签密钥最多200个字符', trigger: 'blur' }
        ],
        server: [
          {
            validator: (rule, value, callback) => {
              if (this.robotForm.type === 'email' && !value) {
                callback(new Error('请输入SMTP服务器'))
              } else {
                callback()
              }
            },
            trigger: 'blur'
          }
        ],
        port: [
          {
            validator: (rule, value, callback) => {
              if (this.robotForm.type === 'email' && !value) {
                callback(new Error('请输入SMTP端口'))
              } else {
                callback()
              }
            },
            trigger: 'blur'
          }
        ],
        username: [
          {
            validator: (rule, value, callback) => {
              if (this.robotForm.type === 'email' && !value) {
                callback(new Error('请输入发件人邮箱'))
              } else if (value) {
                const emailReg = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
                if (!emailReg.test(value)) {
                  callback(new Error('请输入有效的邮箱地址'))
                } else {
                  callback()
                }
              } else {
                callback()
              }
            },
            trigger: 'blur'
          }
        ],
        password: [
          {
            validator: (rule, value, callback) => {
              if (this.robotForm.type === 'email' && !value) {
                callback(new Error('请输入授权码或密码'))
              } else {
                callback()
              }
            },
            trigger: 'blur'
          }
        ],
        template: [
          { max: 5000, message: '模板内容最多5000个字符', trigger: 'blur' }
        ],
        remark: [
          { max: 500, message: '备注最多500个字符', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
  },
  methods: {
    // 获取机器人列表
    async getRobotList() {
      try {
        this.robotLoading = true
        const res = await monitorApi.robot.getList({
          page: this.pagination.page,
          pageSize: this.pagination.pageSize
        })

        if (res.data && res.data.code === 200) {
          // 获取列表后，删除每个机器人的模板字段（由后端管理，前端不需要）
          const list = res.data.data.list || []
          list.forEach(item => {
            delete item.template
          })
          this.robotList = list
          this.pagination.total = res.data.data.total || 0
        } else {
          ElMessage.error(res.data?.msg || res.data?.message || '获取机器人列表失败')
        }
      } catch (error) {
        console.error('获取机器人列表失败:', error)
        ElMessage.error('获取机器人列表失败')
      } finally {
        this.robotLoading = false
      }
    },

    // 刷新机器人列表
    handleRefreshRobots() {
      this.pagination.page = 1
      this.getRobotList()
    },

    // 分页变化
    handlePageChange() {
      this.getRobotList()
    },

    // 新增机器人
    handleAddRobot() {
      this.robotForm = {
        id: null,
        name: '',
        type: 'feishu',
        webhook: '',
        secret: '',
        server: '',
        port: '',
        username: '',
        password: '',
        nickname: '',
        status: 1,
        remark: ''
      }
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs.robotFormRef?.clearValidate()
      })
    },

    // 编辑机器人
    handleEditRobot(row) {
      this.robotForm = {
        id: row.id,
        name: row.name,
        type: row.type,
        webhook: row.webhook || '',
        secret: row.secret || '',
        server: row.server || '',
        port: row.port || '',
        username: row.username || '',
        password: row.password || '',
        nickname: row.nickname || '',
        status: row.status,
        remark: row.remark || ''
      }
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs.robotFormRef?.clearValidate()
      })
    },

    // 保存机器人
    async handleSaveRobot() {
      if (!this.$refs.robotFormRef) return

      try {
        await this.$refs.robotFormRef.validate()
        this.saving = true

        // 构建只包含必要字段的请求数据，排除模板字段（由后端管理）
        const submitData = {
          id: this.robotForm.id,
          name: this.robotForm.name,
          type: this.robotForm.type,
          webhook: this.robotForm.webhook || '',
          secret: this.robotForm.secret || '',
          status: this.robotForm.status,
          remark: this.robotForm.remark || ''
        }

        // 邮箱配置字段
        if (this.robotForm.type === 'email') {
          submitData.server = this.robotForm.server || ''
          submitData.port = this.robotForm.port || ''
          submitData.username = this.robotForm.username || ''
          submitData.password = this.robotForm.password || ''
          submitData.nickname = this.robotForm.nickname || ''
        }

        const apiMethod = this.robotForm.id
          ? monitorApi.robot.update
          : monitorApi.robot.create

        const res = await apiMethod(submitData)

        if (res.data && res.data.code === 200) {
          ElMessage.success(this.robotForm.id ? '更新成功' : '创建成功')
          this.dialogVisible = false
          this.getRobotList()
        } else {
          ElMessage.error(res.data?.msg || res.data?.message || '保存失败')
        }
      } catch (error) {
        if (error !== false) {
          ElMessage.error('保存失败: ' + (error.message || '未知错误'))
        }
      } finally {
        this.saving = false
      }
    },

    // 测试机器人
    async handleTestRobot(row) {
      try {
        ElMessage.info('正在发送测试消息...')

        const res = await monitorApi.robot.test(row.id)

        if (res.data && res.data.code === 200) {
          ElMessage.success('测试消息发送成功！请检查机器人接收情况')
        } else {
          ElMessage.error(res.data?.msg || res.data?.message || '测试消息发送失败')
        }
      } catch (error) {
        console.error('测试消息发送失败:', error)
        ElMessage.error('测试消息发送失败: ' + (error.message || '未知错误'))
      }
    },

    // 删除机器人
    async handleDeleteRobot(row) {
      try {
        await ElMessageBox.confirm('确定删除该机器人吗？', '提示', {
          type: 'warning'
        })

        const res = await monitorApi.robot.delete(row.id)

        if (res.data && res.data.code === 200) {
          ElMessage.success('删除成功')
          this.getRobotList()
        } else {
          ElMessage.error(res.data?.msg || res.data?.message || '删除失败')
        }
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除失败:', error)
          ElMessage.error('删除失败')
        }
      }
    },

    // 更新机器人状态
    async handleStatusChange(row) {
      try {
        const res = await monitorApi.robot.update({
          id: row.id,
          name: row.name,
          type: row.type,
          webhook: row.webhook,
          secret: row.secret || '',
          status: row.status,
          remark: row.remark || ''
        })

        if (res.data && res.data.code === 200) {
          ElMessage.success('状态更新成功')
        } else {
          ElMessage.error(res.data?.msg || res.data?.message || '状态更新失败')
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

    // 类型变化处理
    handleTypeChange() {
      // 清空webhook和secret及邮箱字段
      this.robotForm.webhook = ''
      this.robotForm.secret = ''
      this.robotForm.server = ''
      this.robotForm.port = ''
      this.robotForm.username = ''
      this.robotForm.password = ''
      this.robotForm.nickname = ''

      // 清除验证
      this.$nextTick(() => {
        this.$refs.robotFormRef?.clearValidate()
      })
    },

    // 模板选择处理
    handleTemplateSelect(template) {
      if (template && template.template) {
        this.robotForm.template = template.template
      }
    },

    // 获取类型标签
    getTypeLabel(type) {
      const typeMap = {
        'dingtalk': '钉钉',
        'feishu': '飞书',
        'wechat': '企业微信',
        'email': '邮箱',
        'webhook': 'Webhook',
        'teams': 'Teams'
      }
      return typeMap[type] || '未知'
    },

    // 获取类型标签样式
    getTypeTagType(type) {
      const tagMap = {
        'dingtalk': 'primary',
        'feishu': 'success',
        'wechat': 'warning',
        'email': 'info',
        'webhook': 'danger',
        'teams': 'primary'
      }
      return tagMap[type] || 'info'
    },

    // 获取机器人图标
    getRobotIcon(type) {
      const iconMap = {
        'dingtalk': dingtalkIcon,
        'feishu': feishuIcon,
        'wechat': wechatIcon,
        'email': emailIcon,
        'webhook': webhookIcon,
        'teams': teamsIcon
      }
      return iconMap[type] || ''
    },

    // 格式化日期时间
    formatDateTime(dateStr) {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    },

    // 显示回调示例弹框
    handleShowCallback() {
      this.generateCallbackExamples()
      this.callbackDialogVisible = true
    },

    // 生成回调示例
    generateCallbackExamples() {
      // 示例1: 使用 X-Webhook-Token Header
      this.callbackExample1 = `curl -X POST http://10.7.16.100:8080/api/v1/monitor/webhook/receive \\
  -H "Content-Type: application/json" \\
  -H "X-Webhook-Token: webhook-notify-token-2024" \\
  -d '{
    "source": "zabbix",
    "title": "CPU告警",
    "content": "CPU使用率过高",
    "level": "critical"
  }'`

      // 示例2: 使用 Authorization Bearer
      this.callbackExample2 = `curl -X POST http://10.7.16.100:8080/api/v1/monitor/webhook/receive \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer webhook-notify-token-2024" \\
  -d '{
    "source": "prometheus",
    "title": "服务宕机",
    "content": "API服务无响应",
    "level": "critical"
  }'`

      // 响应示例
      this.responseExample = `{
  "code": 200,
  "message": "成功",
  "data": {
    "status": "success",
    "notifyCount": 1,
    "successCount": 1,
    "failedCount": 0,
    "webhookLogId": 1
  }
}`
    },

    // 复制到剪贴板
    copyToClipboard(text) {
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text).then(() => {
          ElMessage.success('已复制到剪贴板')
        }).catch(() => {
          this.fallbackCopy(text)
        })
      } else {
        this.fallbackCopy(text)
      }
    },

    // 备用复制方法
    fallbackCopy(text) {
      const textarea = document.createElement('textarea')
      textarea.value = text
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      try {
        document.execCommand('copy')
        ElMessage.success('已复制到剪贴板')
      } catch (err) {
        ElMessage.error('复制失败,请手动复制')
      }
      document.body.removeChild(textarea)
    }
  },
  mounted() {
    this.getRobotList()
  }
}
</script>

<style scoped>
.notify-robot-manage {
  padding: 0;
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

.robot-actions {
  display: flex;
  gap: 10px;
}

.callback-section {
  padding: 10px 0;
}

.callback-section h4 {
  margin: 0 0 15px 0;
  font-size: 14px;
  color: #606266;
  font-weight: 600;
}

.callback-code :deep(textarea) {
  font-family: 'Courier New', Consolas, monospace;
  font-size: 13px;
  line-height: 1.6;
  background-color: #f5f7fa;
}

.callback-guide {
  max-height: 70vh;
  overflow-y: auto;
}

.guide-intro {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f0f9ff;
  border-left: 4px solid #409eff;
  border-radius: 4px;
}

.guide-intro p {
  margin: 0;
  color: #606266;
  font-size: 14px;
}

.info-box {
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 15px;
}

.info-box p {
  margin: 8px 0;
  font-size: 14px;
  color: #606266;
}

.info-box code {
  padding: 2px 6px;
  background-color: #e6e8eb;
  border-radius: 3px;
  color: #e74c3c;
  font-family: 'Courier New', Consolas, monospace;
  font-size: 13px;
}

.guide-tips {
  margin-top: 25px;
  padding: 20px;
  background-color: #fffbf0;
  border-radius: 4px;
  border: 1px solid #ffe58f;
}

.guide-tips h4 {
  margin: 0 0 15px 0;
  font-size: 15px;
  color: #606266;
  font-weight: 600;
}

.guide-tips ul {
  margin: 0;
  padding-left: 20px;
}

.guide-tips li {
  margin: 8px 0;
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
}

.guide-tips code {
  padding: 2px 6px;
  background-color: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 3px;
  color: #e74c3c;
  font-family: 'Courier New', Consolas, monospace;
  font-size: 13px;
}
</style>
