<template>
  <el-dialog title="编辑主机" v-model="dialogVisible" width="50%" @close="handleClose">
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
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
                v-model="form.port" 
                placeholder="端口" 
                style="width: 100px" 
              />
            </div>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12">
          <el-form-item label="认证凭据" prop="authId">
            <el-select 
              v-model="form.authId" 
              placeholder="请选择认证凭据" 
              style="width: 100%"
              filterable
              clearable
            >
              <el-option
                v-for="item in authList"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="主机类型" prop="vendor">
            <el-select v-model="form.vendor" placeholder="请选择主机类型" style="width: 100%">
              <el-option :value="1" :label="'自建主机'">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <el-icon :size="18" color="#409EFF"><OfficeBuilding /></el-icon>
                  <span>自建主机</span>
                </div>
              </el-option>
              <el-option :value="2" :label="'阿里云'">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img src="@/assets/image/aliyun.png" style="width: 18px; height: 18px"/>
                  <span>阿里云</span>
                </div>
              </el-option>
              <el-option :value="3" :label="'腾讯云'">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img src="@/assets/image/tengxun.png" style="width: 18px; height: 18px"/>
                  <span>腾讯云</span>
                </div>
              </el-option>
              <el-option :value="4" :label="'百度云'">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img src="@/assets/image/baidu.svg" style="width: 18px; height: 18px"/>
                  <span>百度云</span>
                </div>
              </el-option>
              <el-option :value="5" :label="'华为云'">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img src="@/assets/image/华为.svg" style="width: 18px; height: 18px"/>
                  <span>华为云</span>
                </div>
              </el-option>
              <el-option :value="6" :label="'AWS云'">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img src="@/assets/image/aws.svg" style="width: 18px; height: 18px"/>
                  <span>AWS</span>
                </div>
              </el-option>
            </el-select>
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
              多个标签用逗号分隔，格式为 key=value
            </div>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24">
          <el-form-item label="备注" prop="remark">
            <el-input v-model="form.remark" type="textarea" placeholder="请输入备注信息" />
          </el-form-item>
        </el-col>
      </el-row>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitForm">确 定</el-button>
        <el-button @click="dialogVisible = false">取 消</el-button>
      </span>
    </el-form>
  </el-dialog>
</template>

<script>
export default {
  name: 'EditHost',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    hostInfo: {
      type: Object,
      required: true
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
      formReady: false,
      form: {
        id: '',
        hostName: '',
        groupId: '',
        remark: '',
        ip: '',
        port: 22,
        username: '',
        authId: '',
        vendor: 1,
        tag: ''
      },
      rules: {
        hostName: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
        ip: [
          { required: true, message: '请输入IP地址', trigger: 'blur' },
          { pattern: /^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$/, message: '请输入有效的IP地址', trigger: 'blur' }
        ],
        port: [
          { required: true, message: '请输入端口号', trigger: 'blur' },
          { type: 'number', message: '端口号必须为数字', trigger: 'blur' },
          { min: 1, max: 65535, message: '端口号范围1-65535', trigger: 'blur' }
        ],
        username: [
          { required: true, message: '请输入连接用户名', trigger: 'blur' },
          { min: 1, max: 32, message: '用户名长度1-32个字符', trigger: 'blur' }
        ],
        authId: [{ required: true, message: '请选择认证凭据', trigger: 'change' }],
        groupId: [{ required: true, message: '请选择所属分组', trigger: 'change' }],
        vendor: [{ required: true, message: '请选择主机类型', trigger: 'change' }]
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
      if (newVal) {
        this.$nextTick(() => {
          this.initForm()
        })
      }
    },
    hostInfo: {
      deep: true,
      immediate: true,
      handler(newVal) {
        if (newVal) {
          this.form = {
            id: newVal.id || '',
            hostName: newVal.hostName || '',
            groupId: newVal.groupId || (this.groupList.length > 0 ? this.groupList[0].id : ''),
            remark: newVal.remark || '',
            ip: newVal.sshIp || newVal.ip || '',
            port: newVal.sshPort || newVal.port || 22,
            username: newVal.sshName || newVal.username || 'root',
            authId: newVal.sshKeyId || newVal.authId || (this.authList.length > 0 ? this.authList[0].id : ''),
            vendor: newVal.vendor !== undefined ? Number(newVal.vendor) : 1,
            tag: newVal.tag || ''
          }
        }
      }
    }
  },
  mounted() {
    if (this.visible) {
      this.initForm()
    }
  },
  methods: {
    async initForm() {
      await this.$nextTick()
      this.formReady = !!this.$refs.formRef
      if (!this.formReady) {
        console.error('表单初始化失败', {
          refs: this.$refs,
          formRef: this.$refs.formRef
        })
      }
    },
    handleClose() {
      if (this.formReady) {
        this.$refs.formRef.resetFields()
      }
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
    async submitForm() {
      try {
        // 确保表单组件已挂载
        await this.$nextTick()
        
        // 直接检查表单引用
        if (!this.$refs.formRef || !this.$refs.formRef.validate) {
          console.error('表单引用无效:', this.$refs.formRef)
          return this.$message.error('表单组件未准备好，请稍后再试')
        }

        // 直接调用validate方法
        return new Promise((resolve) => {
          this.$refs.formRef.validate((valid) => {
            if (valid) {
              // 验证凭据是否存在
              const authExists = this.authList.some(auth => auth.id === this.form.authId)
              if (!authExists) {
                this.$message.error('选择的认证凭据不存在，请刷新凭据列表后重试')
                return resolve(false)
              }

              // 确保端口为数字
              const port = Number(this.form.port)
              if (isNaN(port) || port < 1 || port > 65535) {
                this.$message.error('端口号必须为1-65535之间的数字')
                return resolve(false)
              }

              // 转换数据格式为后端期望的结构
              const requestData = {
                id: this.form.id,
                groupId: this.form.groupId,
                hostName: this.form.hostName,
                remark: this.form.remark,
                sshIp: this.form.ip,
                sshKeyId: this.form.authId,
                sshName: this.form.username,
                sshPort: port,
                vendor: this.form.vendor,
                tag: this.form.tag || ''
              }

              this.$emit('submit', requestData)
              this.dialogVisible = false
              resolve(true)
            } else {
              this.$message.error('请填写完整的表单信息')
              resolve(false)
            }
          })
        })
      } catch (error) {
        console.error('表单提交失败:', error)
        this.$message.error(`表单提交失败: ${error.message}`)
        return false
      }
    }
  }
}
</script>

<style scoped>
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
