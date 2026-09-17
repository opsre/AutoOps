<template>
  <el-dialog 
    title="新增主机" 
    v-model="dialogVisible" 
    width="40%" 
    @close="handleClose"
  >
    <el-form 
      :model="form" 
      :rules="rules" 
      ref="formRef" 
      label-width="100px"
    >
      <el-row>
        <el-col :span="12">
          <el-form-item label="主机名称" prop="hostName">
            <el-input v-model="form.hostName" placeholder="请输入主机名称" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="所属分组" prop="groupId">
            <el-select
              v-model="form.groupId"
              placeholder="请选择分组"
              style="width: 100%"
              clearable
              filterable
              @change="handleGroupChange"
            >
              <el-option
                v-for="group in selectableGroups"
                :key="group.id"
                :label="group.displayName"
                :value="group.id"
                :class="group.isBusinessGroup ? 'business-group-option' : 'sub-group-option'"
              >
                <div class="group-option-content">
                  <span :class="group.isBusinessGroup ? 'business-group-text' : 'sub-group-text'">
                    {{ group.displayName }}
                  </span>
                  <el-tag v-if="group.isBusinessGroup" type="success" size="small" class="business-tag">
                    默认
                  </el-tag>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24">
          <el-form-item label="SSH连接" prop="sshConnection">
            <div style="display: flex; align-items: center">
              <el-input 
                v-model="form.username" 
                placeholder="用户名" 
                style="width: 120px; margin-right: 8px" 
              />
              <span style="margin: 0 4px">@</span>
              <el-input 
                v-model="form.ip" 
                placeholder="IP地址" 
                style="width: 150px; margin-right: 8px" 
              />
              <span style="margin: 0 4px">-p</span>
              <el-input
                v-model.number="form.port"
                type="number"
                :min="1"
                :max="65535"
                placeholder="端口"
                style="width: 80px"
              />
            </div>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24">
          <el-form-item label="认证凭据" prop="authId">
            <div style="display: flex; align-items: center; gap: 8px;">
              <el-select 
                v-model="form.authId" 
                placeholder="请选择认证凭据" 
                style="flex: 1"
                filterable
              >
                <el-option
                  v-for="item in authList"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
              <el-button 
                type="warning" 
                icon="Plus" 
                size="default"
                class="create-credential-btn"
                @click="showCreateCredentialDialog"
              >
                创建凭据
              </el-button>
            </div>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24">
          <el-form-item label="标签" prop="tag">
            <el-input
              v-model="form.tag"
              placeholder="格式: key=value,key=value 例如: bline=s3,env=prod"
            />
            <div style="font-size: 12px; color: #909399; margin-top: 4px;">
              多个标签用逗号分隔，格式为 bline=s3,env=prod
            </div>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24">
          <el-form-item label="备注" prop="remark">
            <el-input v-model="form.remark" type="textarea" placeholder="请输入备注信息"></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitForm">确 定</el-button>
        <el-button @click="dialogVisible = false">取 消</el-button>
      </span>
    </el-form>
    
    <!-- 创建凭据对话框 -->
    <el-dialog
      title="创建凭据"
      v-model="credentialDialogVisible"
      width="50%"
      :modal="false"
      @close="handleCredentialDialogClose"
    >
      <el-form 
        :model="credentialForm" 
        :rules="credentialRules" 
        ref="credentialFormRef" 
        label-width="100px"
      >
        <el-form-item label="凭据名称" prop="name">
          <el-input v-model="credentialForm.name" placeholder="请输入凭据名称" />
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="credentialForm.username" placeholder="请输入用户名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="端口" prop="port">
              <el-input 
                v-model.number="credentialForm.port" 
                type="number" 
                :min="1" 
                :max="65535"
                placeholder="请输入端口"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="认证类型" prop="type">
          <el-radio-group v-model="credentialForm.type">
            <el-radio :label="1">密码认证</el-radio>
            <el-radio :label="2">密钥认证</el-radio>
            <el-radio :label="3">公钥认证</el-radio>
            <el-radio :label="4">跳板机认证</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 公钥认证提示信息 -->
        <el-form-item v-if="credentialForm.type === 3" label="配置说明">
          <el-alert
            title="公钥认证配置说明"
            type="info"
            :closable="false"
            show-icon>
            <template #default>
              <div style="line-height: 1.6; margin-top: 8px;">
                <ol style="margin: 0; padding-left: 20px;">
                  <li>复制DevOps服务器的公钥: <code style="background: #f5f5f5; padding: 2px 4px; border-radius: 3px;">cat ~/.ssh/id_rsa.pub</code></li>
                  <li>将公钥添加到目标主机: <code style="background: #f5f5f5; padding: 2px 4px; border-radius: 3px;">echo "公钥内容" >> /root/.ssh/authorized_keys</code></li>
                </ol>
                <p style="margin: 12px 0 0 0; color: #909399; font-size: 13px;">
                  💡 提示:公钥认证无需存储密码或密钥，系统会自动使用DevOps服务器的私钥进行认证。
                </p>
              </div>
            </template>
          </el-alert>
        </el-form-item>
        
        <el-form-item v-if="credentialForm.type === 1" label="密码" prop="password">
          <el-input 
            v-model="credentialForm.password" 
            show-password 
            placeholder="请输入密码"
          />
        </el-form-item>
        
        <!-- 密钥认证配置说明 -->
        <el-form-item v-if="credentialForm.type === 2" label="配置说明">
          <el-alert
            title="密钥认证配置说明"
            type="warning"
            :closable="false"
            show-icon>
            <template #default>
              <div style="line-height: 1.6; margin-top: 8px;">
                <p style="margin: 0 0 8px 0; font-weight: 600;">请按以下步骤配置密钥认证：</p>
                <ol style="margin: 0; padding-left: 20px;">
                  <li>在目标主机执行: <code style="background: #f5f5f5; padding: 2px 4px; border-radius: 3px;">cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys</code></li>
                  <li>复制目标主机的私钥内容到下方文本框</li>
                  <li>私钥格式应包含完整的BEGIN和END标记</li>
                </ol>
                <p style="margin: 12px 0 0 0; color: #909399; font-size: 13px;">
                  💡 提示：密钥认证需要预先在目标主机配置公钥，然后上传对应的私钥内容。
                </p>
              </div>
            </template>
          </el-alert>
        </el-form-item>

        <el-form-item v-if="credentialForm.type === 2" label="私钥内容" prop="publicKey">
          <el-input
            v-model="credentialForm.publicKey"
            type="textarea"
            :rows="8"
            placeholder="请输入SSH私钥内容，格式如下：
-----BEGIN OPENSSH PRIVATE KEY-----
xxxxxxxxxxx
-----END OPENSSH PRIVATE KEY-----"
          />
        </el-form-item>

        <!-- 跳板机认证配置 -->
        <template v-if="credentialForm.type === 4">
          <el-divider content-position="left">跳板机配置</el-divider>

          <el-row>
            <el-col :span="12">
              <el-form-item label="跳板机IP" prop="bastionHost">
                <el-input v-model="credentialForm.bastionHost" placeholder="例如：123.45.67.89"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="跳板机端口" prop="bastionPort">
                <el-input-number v-model="credentialForm.bastionPort" :min="1" :max="65535" style="width: 100%"></el-input-number>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="跳板机用户名" prop="bastionUsername">
            <el-input v-model="credentialForm.bastionUsername" placeholder="例如：admin"></el-input>
          </el-form-item>

          <el-form-item label="跳板机认证方式" prop="bastionAuthType">
            <el-select v-model="credentialForm.bastionAuthType" style="width: 100%">
              <el-option label="密码认证" :value="1"></el-option>
              <el-option label="私钥认证" :value="2"></el-option>
              <el-option label="免密认证" :value="3"></el-option>
            </el-select>
          </el-form-item>

          <el-form-item v-if="credentialForm.bastionAuthType === 1" label="跳板机密码" prop="bastionPassword">
            <el-input v-model="credentialForm.bastionPassword" type="password" show-password></el-input>
          </el-form-item>

          <el-form-item v-if="credentialForm.bastionAuthType === 2" label="跳板机私钥" prop="bastionPrivateKey">
            <el-input
              v-model="credentialForm.bastionPrivateKey"
              type="textarea"
              :rows="6"
              placeholder="粘贴私钥内容（PEM格式）"></el-input>
          </el-form-item>

          <el-divider content-position="left">目标主机认证配置</el-divider>

          <el-form-item label="目标主机认证方式" prop="targetAuthType">
            <el-select v-model="credentialForm.targetAuthType" style="width: 100%">
              <el-option label="密码认证" :value="1"></el-option>
              <el-option label="私钥认证" :value="2"></el-option>
              <el-option label="跳板机代理免密" :value="4"></el-option>
            </el-select>
          </el-form-item>

          <el-form-item v-if="credentialForm.targetAuthType === 1" label="目标主机密码" prop="targetPassword">
            <el-input v-model="credentialForm.targetPassword" type="password" show-password></el-input>
          </el-form-item>

          <el-form-item v-if="credentialForm.targetAuthType === 2" label="目标主机私钥" prop="targetPrivateKey">
            <el-input
              v-model="credentialForm.targetPrivateKey"
              type="textarea"
              :rows="6"
              placeholder="粘贴私钥内容（PEM格式）"></el-input>
          </el-form-item>
        </template>

        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="credentialForm.remark"
            type="textarea"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitCredentialForm">确 定</el-button>
        <el-button @click="credentialDialogVisible = false">取 消</el-button>
      </span>
    </el-dialog>
  </el-dialog>
</template>

<script>
import API from '@/api/config'

export default {
  name: 'CreateHost',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    groupList: {
      type: Array,
      required: true
    },
    authList: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      dialogVisible: this.visible,
      credentialDialogVisible: false,
      form: {
        hostName: '',
        ip: '',
        port: 22,
        username: '',
        authId: '',
        groupId: '',
        tag: '',
        remark: ''
      },
      credentialForm: {
        name: '',
        type: undefined,
        username: '',
        password: '',
        publicKey: '',
        port: '',
        remark: '',
        // 跳板机字段
        bastionHost: '',
        bastionPort: 22,
        bastionUsername: '',
        bastionAuthType: 1,
        bastionPassword: '',
        bastionPrivateKey: '',
        targetAuthType: 4,
        targetPassword: '',
        targetPrivateKey: ''
      },
      rules: {
        hostName: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
        ip: [{ required: true, message: '请输入IP地址', trigger: 'blur' }],
        port: [{ required: true, message: '请输入端口号', trigger: 'blur' }],
        username: [{ required: true, message: '请输入连接用户名', trigger: 'blur' }],
        authId: [{ required: true, message: '请选择认证凭据', trigger: 'change' }],
        groupId: [{ required: true, message: '请选择所属分组', trigger: 'change' }]
      },
      credentialRules: {
        name: [{ required: true, message: '请输入凭据名称', trigger: 'blur' }],
        type: [{ required: true, message: '请选择认证类型', trigger: 'change' }],
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [
          { 
            required: true, 
            message: '请输入密码', 
            trigger: 'blur',
            validator: (rule, value, callback) => {
              if (this.credentialForm.type === 1 && !value) {
                callback(new Error('请输入密码'))
              } else {
                callback()
              }
            }
          }
        ],
        publicKey: [
          { 
            required: true, 
            message: '请输入公钥', 
            trigger: 'blur',
            validator: (rule, value, callback) => {
              if (this.credentialForm.type === 2 && !value) {
                callback(new Error('请输入公钥'))
              } else {
                callback()
              }
            }
          }
        ],
        port: [
          { 
            required: true, 
            message: '请输入端口号', 
            trigger: 'blur',
            validator: (rule, value, callback) => {
              if (!value) {
                callback(new Error('请输入端口号'))
              } else if (isNaN(value) || value < 1 || value > 65535) {
                callback(new Error('请输入1-65535之间的有效端口号'))
              } else {
                callback()
              }
            }
          }
        ]
      }
    }
  },
  computed: {
    // 生成可选择的分组列表：业务组 + 所有二级分组
    selectableGroups() {
      const groups = []
      
      // 遍历分组树，提取业务组和二级分组
      const extractGroups = (groupList, parentName = '') => {
        groupList.forEach(group => {
          // 如果是业务组，直接添加
          if (group.name === '业务组' || group.isDefault) {
            groups.push({
              id: group.id,
              name: group.name,
              displayName: group.name,
              isBusinessGroup: true
            })
          }
          
          // 如果有子分组，遍历子分组（二级）
          if (group.children && group.children.length > 0) {
            group.children.forEach(subGroup => {
              // 只有叶子节点（没有子节点）才可选
              if (!subGroup.children || subGroup.children.length === 0) {
                groups.push({
                  id: subGroup.id,
                  name: subGroup.name,
                  displayName: subGroup.name,
                  parentName: group.name,
                  isBusinessGroup: false
                })
              }
            })
          }
        })
      }
      
      extractGroups(this.groupList)
      return groups
    }
  },
  watch: {
    visible(newVal) {
      this.dialogVisible = newVal
    }
  },
  methods: {
    handleClose() {
      this.$refs.formRef.resetFields()
      this.$emit('close')
    },
    handleGroupChange(value) {
      if (value) {
        this.form.groupId = value
      } else {
        // 清空时设置默认分组
        const defaultGroup = this.selectableGroups.find(group => group.isBusinessGroup)
        if (defaultGroup) {
          this.form.groupId = defaultGroup.id
        }
      }
    },
    showCreateCredentialDialog() {
      this.credentialForm = {
        name: '',
        type: undefined,
        username: '',
        password: '',
        publicKey: '',
        port: '',
        remark: '',
        bastionHost: '',
        bastionPort: 22,
        bastionUsername: '',
        bastionAuthType: 1,
        bastionPassword: '',
        bastionPrivateKey: '',
        targetAuthType: 4,
        targetPassword: '',
        targetPrivateKey: ''
      }
      this.credentialDialogVisible = true
    },
    
    handleCredentialDialogClose() {
      this.$refs.credentialFormRef?.resetFields()
    },
    
    async submitCredentialForm() {
      try {
        await this.$refs.credentialFormRef.validate()

        const typeValue = Number(this.credentialForm.type)
        if (typeValue !== 1 && typeValue !== 2 && typeValue !== 3 && typeValue !== 4) {
          throw new Error('请选择有效的认证类型')
        }

        const formData = {
          id: '',
          name: this.credentialForm.name,
          type: typeValue,
          username: this.credentialForm.username,
          password: this.credentialForm.password || '',
          publicKey: this.credentialForm.publicKey || '',
          port: this.credentialForm.port,
          remark: this.credentialForm.remark || ''
        }

        // 如果是跳板机认证，添加跳板机相关字段
        if (typeValue === 4) {
          formData.bastionHost = this.credentialForm.bastionHost
          formData.bastionPort = this.credentialForm.bastionPort
          formData.bastionUsername = this.credentialForm.bastionUsername
          formData.bastionAuthType = this.credentialForm.bastionAuthType
          formData.targetAuthType = this.credentialForm.targetAuthType

          // 根据跳板机认证方式添加对应字段
          if (this.credentialForm.bastionAuthType === 1) {
            formData.bastionPassword = this.credentialForm.bastionPassword
          } else if (this.credentialForm.bastionAuthType === 2) {
            formData.bastionPrivateKey = this.credentialForm.bastionPrivateKey
          }

          // 根据目标主机认证方式添加对应字段
          if (this.credentialForm.targetAuthType === 1) {
            formData.targetPassword = this.credentialForm.targetPassword
          } else if (this.credentialForm.targetAuthType === 2) {
            formData.targetPrivateKey = this.credentialForm.targetPrivateKey
          }
        }

        console.log('Creating credential with data:', formData)

        // 直接调用API创建凭据
        const res = await API.createEcsAuth(formData)
        
        if (res.data.code === 200) {
          this.$message.success('创建凭据成功')
          this.credentialDialogVisible = false
          
          // 通知父组件刷新认证列表
          this.$emit('refresh-auth-list')
          
          // 自动选中新创建的凭据
          const newCredential = res.data.data
          if (newCredential && newCredential.id) {
            this.form.authId = newCredential.id
          }
        } else {
          this.$message.error(res.data.message || '创建凭据失败')
        }
      } catch (error) {
        console.error('创建凭据失败:', error)
        this.$message.error('创建凭据失败: ' + (error.response?.data?.message || error.message))
      }
    },
    
    async submitForm() {
      try {
        await this.$refs.formRef.validate()

        const authExists = this.authList.some(auth => auth.id === this.form.authId)
        if (!authExists) {
          return this.$message.error('选择的认证凭据不存在，请刷新凭据列表后重试')
        }

        const requestData = {
          groupId: this.form.groupId,
          hostName: this.form.hostName,
          remark: this.form.remark,
          sshIp: this.form.ip,
          sshKeyId: this.form.authId,
          sshName: this.form.username,
          sshPort: parseInt(this.form.port) || 22,
          tag: this.form.tag || ''
        }

        // 等待父组件处理完成后再关闭对话框
        this.$emit('submit', requestData)
        // 注意：对话框的关闭现在由父组件通过 close 事件或直接设置 visible 来控制
      } catch (error) {
        console.error('表单验证失败:', error)
      }
    }
  }
}
</script>

<style scoped>
.create-credential-btn {
  background-color: #f39c12 !important;
  border-color: #f39c12 !important;
  color: #ffffff !important;
  border-radius: 12px !important;
  padding: 10px 16px !important;
  font-weight: 500 !important;
  transition: all 0.3s ease !important;
}

.create-credential-btn:hover {
  background-color: #e67e22 !important;
  border-color: #e67e22 !important;
  color: #ffffff !important;
  transform: translateY(-1px) !important;
  box-shadow: 0 4px 8px rgba(243, 156, 18, 0.3) !important;
}

.create-credential-btn:focus,
.create-credential-btn:active {
  background-color: #d68910 !important;
  border-color: #d68910 !important;
  color: #ffffff !important;
}

/* 分组选择器样式 */
.group-option-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.business-group-option {
  background-color: #f0f9ff !important;
}

.business-group-text {
  color: #059669 !important;
  font-weight: 600 !important;
}

.sub-group-text {
  color: #374151 !important;
}

.business-tag {
  margin-left: 8px !important;
}

.sub-group-option:hover {
  background-color: #f3f4f6 !important;
}

.business-group-option:hover {
  background-color: #e0f2fe !important;
}
</style>
