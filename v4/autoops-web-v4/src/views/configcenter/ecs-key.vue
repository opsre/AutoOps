<template>
  <div class="ecs-key-management">
    <el-card shadow="hover" class="ecs-key-card">
      <template #header>
        <div class="card-header">
          <span class="title">ECS密钥管理</span>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :model="queryParams" :inline="true" class="search-form" ref="queryForm">
      <el-form-item label="凭据名称" prop="name">
        <el-input placeholder="请输入凭据名称" clearable size="small" v-model="queryParams.name"
                  @keyup.enter.native="handleQuery" />
      </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="Search" size="small" @click="handleQuery">搜索</el-button>
            <el-button type="warning" icon="Refresh" size="small" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    
      <!-- 操作按钮 -->
      <div class="action-section">
        <el-button type="success" v-authority="['config:ecs:create']" plain icon="Plus" size="small"
                   @click="showAddDialog">创建凭据</el-button>
      </div>
    
      <!-- 列表区域 -->
      <div class="table-section">
        <el-table stripe v-loading="loading" :data="authList" class="ecs-key-table">
      <el-table-column label="凭据名称" prop="name">
        <template v-slot="scope">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img src="@/assets/image/凭据.svg" style="width: 16px; height: 16px"/>
            <span>{{ scope.row.name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="用户名" prop="username" v-if="authList.some(item => item.type === 1)">
        <template v-slot="scope">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img src="@/assets/image/ren.svg" style="width: 16px; height: 16px"/>
            <span>{{ scope.row.username }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="端口" prop="port" width="100">
        <template v-slot="scope">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img src="@/assets/image/端口.svg" style="width: 16px; height: 16px"/>
            <span>{{ scope.row.port }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="认证信息" width="120">
        <template v-slot="scope">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img
              :src="scope.row.type === 1 ? require('@/assets/image/密码.svg') : require('@/assets/image/密钥.svg')"
              style="width: 16px; height: 16px"
            />
            <el-tag :type="scope.row.type === 1 ? 'success' : (scope.row.type === 4 ? 'danger' : 'info')">
              {{
                scope.row.type === 1 ? '密码认证' :
                (scope.row.type === 3 ? '公钥认证' :
                (scope.row.type === 4 ? '跳板机认证' : '未知类型'))
              }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="createTime" />
      <el-table-column label="备注" prop="remark" />
      <el-table-column label="操作" class-name="small-padding fixed-width" width="220">
        <template v-slot="scope">
          <div class="operation-buttons">
            <el-tooltip content="修改" placement="top">
              <el-button size="small" v-authority="['config:ecs:edit']" type="primary" icon="Edit" @click="showEditDialog(scope.row)">修改</el-button>
            </el-tooltip>
            <el-tooltip content="删除" placement="top">
              <el-button size="small" v-authority="['config:ecs:delete']" type="danger" icon="Delete" @click="handleDelete(scope.row)">删除</el-button>
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
        </el-table>
      </div>

      <!-- 分页区域 -->
      <div class="pagination-section">
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
    
    <!--新增/编辑对话框-->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="45%" :modal="false">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="凭据名称" prop="name">
          <el-input v-model="formData.name"></el-input>
        </el-form-item>
        <el-form-item label="认证类型" prop="type">
          <el-radio-group v-model="formData.type">
            <el-radio :label="1">密码认证</el-radio>
            <el-radio :label="3">公钥认证</el-radio>
            <el-radio :label="4">跳板机认证</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 密码认证配置 -->
        <template v-if="formData.type === 1">
          <el-row>
            <el-col :span="12">
              <el-form-item label="用户名" prop="username">
                <el-input v-model="formData.username"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="端口" prop="port">
                <el-input v-model.number="formData.port" type="number" :min="1" :max="65535"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="密码" prop="password">
            <el-input v-model="formData.password" show-password></el-input>
          </el-form-item>
        </template>

        <!-- 公钥认证配置 -->
        <template v-if="formData.type === 3">
          <el-form-item label="配置说明">
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
        </template>

        <!-- 跳板机认证配置 -->
        <template v-if="formData.type === 4">
          <el-divider content-position="left">跳板机配置</el-divider>

          <el-row>
            <el-col :span="8">
              <el-form-item label="用户名" prop="bastionUsername">
                <el-input v-model="formData.bastionUsername" placeholder="例如：admin"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="跳板机IP" prop="bastionHost">
                <el-input v-model="formData.bastionHost" placeholder="例如：123.45.67.89"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="端口" prop="bastionPort">
                <el-input-number v-model="formData.bastionPort" :min="1" :max="65535" style="width: 100%"></el-input-number>
              </el-form-item>
            </el-col>

          </el-row>

          <el-form-item label="认证方式" prop="bastionAuthType">
            <el-select v-model="formData.bastionAuthType" style="width: 100%">
              <el-option label="密码认证" :value="1"></el-option>
              <el-option label="私钥认证" :value="2"></el-option>
              <el-option label="免密认证" :value="3"></el-option>
            </el-select>
          </el-form-item>

          <el-form-item v-if="formData.bastionAuthType === 1" label="跳板机密码" prop="bastionPassword">
            <el-input v-model="formData.bastionPassword" type="password" show-password></el-input>
          </el-form-item>

          <el-form-item v-if="formData.bastionAuthType === 2" label="跳板机私钥" prop="bastionPrivateKey">
            <el-input
              v-model="formData.bastionPrivateKey"
              type="textarea"
              :rows="6"
              placeholder="粘贴私钥内容（PEM格式）"></el-input>
          </el-form-item>

          <el-divider content-position="left">目标主机认证配置</el-divider>

          <el-form-item label="认证方式" prop="targetAuthType">
            <el-select v-model="formData.targetAuthType" style="width: 100%">
              <el-option label="密码认证" :value="1"></el-option>
              <el-option label="跳板机代理免密" :value="4"></el-option>
            </el-select>
          </el-form-item>

          <el-form-item v-if="formData.targetAuthType === 1" label="目标主机密码" prop="targetPassword">
            <el-input v-model="formData.targetPassword" type="password" show-password></el-input>
          </el-form-item>
        </template>

        <el-form-item label="备注" prop="remark">
          <el-input v-model="formData.remark" type="textarea"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitForm">确 定</el-button>
        <el-button @click="dialogVisible = false">取 消</el-button>
      </span>
    </el-dialog>
    </el-card>
  </div>
</template>

<script>
import API from '@/api/config'

export default {
  data() {
    return {
      queryParams: {
        name: '',
        pageNum: 1,
        pageSize: 10
      },
      loading: false,
      authList: [],
      total: 0,
      dialogVisible: false,
      dialogTitle: '',
      formData: {
        id: '',
        name: '',
        type: undefined, // 初始不设置默认值
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
      formRules: {
        name: [{ required: true, message: '请输入凭据名称', trigger: 'blur' }],
        type: [{ required: true, message: '请选择认证类型', trigger: 'change' }],
        username: [
          {
            validator: (rule, value, callback) => {
              if (this.formData.type === 1 && !value) {
                callback(new Error('请输入用户名'))
              } else {
                callback()
              }
            },
            trigger: 'blur'
          }
        ],
        password: [
          { 
            required: true, 
            message: '请输入密码', 
            trigger: 'blur',
            validator: (rule, value, callback) => {
              if (this.formData.type === 1 && !value) {
                callback(new Error('请输入密码'))
              } else {
                callback()
              }
            }
          }
        ],
        publicKey: [
          {
            required: false,
            message: '请输入公钥',
            trigger: 'blur'
          }
        ],
        port: [
          { 
            trigger: 'blur',
            validator: (rule, value, callback) => {
              if (this.formData.type !== 1) {
                callback()
              } else if (!value) {
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
  methods: {
    // 获取凭据列表
    async getList() {
      this.loading = true
      try {
        const { data: res } = await API.getEcsAuthList({
          page: this.queryParams.pageNum,
          pageSize: this.queryParams.pageSize,
          name: this.queryParams.name || undefined
        })
        if (res.code === 200) {
          this.authList = res.data?.list || []
          this.total = res.data?.total || 0
        } else {
          this.$message.error(res.message || '获取凭据列表失败')
        }
      } catch (error) {
        console.error('获取凭据列表失败:', error)
        this.$message.error('获取凭据列表失败')
      } finally {
        this.loading = false
      }
    },
    
    // 搜索
    async handleQuery() {
      this.queryParams.pageNum = 1  // 重置到第一页
      this.getList()
    },
    
    // 重置搜索
    resetQuery() {
      this.queryParams = {
        name: '',
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
      this.dialogTitle = '创建凭据'
      this.$nextTick(() => {
        this.formData = {
          id: '',
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
        this.dialogVisible = true
        console.log('Add dialog opened, formData:', this.formData)
      })
    },
    
    // 显示编辑对话框
    showEditDialog(row) {
      this.dialogTitle = '修改凭据'
      this.$nextTick(() => {
        this.formData = {
          id: row.id,
          name: row.name,
          type: Number(row.type),
          username: row.username,
          password: row.type === 1 ? row.password : '',
          publicKey: row.type === 2 ? row.publicKey : '',
          port: row.port || '',
          remark: row.remark || '',
          bastionHost: row.bastionHost || '',
          bastionPort: row.bastionPort || 22,
          bastionUsername: row.bastionUsername || '',
          bastionAuthType: row.bastionAuthType || 1,
          bastionPassword: row.bastionPassword || '',
          bastionPrivateKey: row.bastionPrivateKey || '',
          targetAuthType: row.targetAuthType || 4,
          targetPassword: row.targetPassword || '',
          targetPrivateKey: row.targetPrivateKey || ''
        }
        this.dialogVisible = true
        console.log('Edit dialog opened for row:', row, 'formData:', this.formData)
      })
    },
    
    // 提交表单
    async submitForm() {
      try {
        await this.$refs.formRef.validate()
        
        // 确保type是数字类型且有效
        const typeValue = Number(this.formData.type)
        if (typeValue !== 1 && typeValue !== 3 && typeValue !== 4) {
          throw new Error('请选择有效的认证类型')
        }

        const formData = {
          ...this.formData,
          type: typeValue
        }

        // 清理不需要的字段
        if (typeValue !== 4) {
          // 非跳板机认证：删除跳板机相关字段
          delete formData.bastionHost
          delete formData.bastionPort
          delete formData.bastionUsername
          delete formData.bastionAuthType
          delete formData.bastionPassword
          delete formData.bastionPrivateKey
          delete formData.targetAuthType
          delete formData.targetPassword
          delete formData.targetPrivateKey

          if (typeValue === 3) {
            delete formData.username
            delete formData.port
            delete formData.password
            delete formData.publicKey
          }
        } else {
          // 跳板机认证：删除传统字段（username, password, publicKey, port）
          // 因为目标主机信息来自前端传递（dto.SSHIP, dto.SSHPort, dto.SSHName）
          delete formData.username
          delete formData.password
          delete formData.publicKey
          delete formData.port

          // 清理跳板机不需要的认证字段
          if (formData.bastionAuthType !== 1) {
            delete formData.bastionPassword
          }
          if (formData.bastionAuthType !== 2) {
            delete formData.bastionPrivateKey
          }
          // 清理目标主机不需要的认证字段
          if (formData.targetAuthType !== 1) {
            delete formData.targetPassword
          }
          if (formData.targetAuthType !== 2) {
            delete formData.targetPrivateKey
          }
        }

        console.log('Submitting form with data:', formData)

        let res
        if (formData.id) {
          // 更新
          res = await API.updateEcsAuth(formData)
        } else {
          // 新增
          res = await API.createEcsAuth(formData)
        }

        if (res.data.code === 200) {
          this.$message.success(formData.id ? '修改成功' : '创建成功')
          this.dialogVisible = false
          // 强制刷新列表数据
          await this.getList()
          console.log('List after update:', this.authList)
        } else {
          this.$message.error(res.data.message || (formData.id ? '修改失败' : '创建失败'))
        }
      } catch (error) {
        console.error('操作失败:', error)
        this.$message.error('操作失败: ' + error.message)
      }
    },
    
    // 删除凭据
    async handleDelete(row) {
      try {
        await this.$confirm(`确定删除凭据"${row.name}"?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        const { data: res } = await API.deleteEcsAuth(row.id)
        if (res.code === 200) {
          this.$message.success('删除成功')
          this.getList()
        } else {
          this.$message.error(res.message || '删除失败')
        }
      } catch (error) {
        console.error('删除失败:', error)
      }
    }
  },
  created() {
    this.getList()
  }
}
</script>

<style scoped>
.ecs-key-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.ecs-key-card {
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
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.search-section {
  margin-bottom: 5px;
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

.action-section {
  margin-bottom: 5px;
  margin-top: 5px;
  padding: 12px 0;
}

.action-section .el-button {
  margin-right: 12px;
}

.table-section {
  margin-bottom: 40px;
}

.ecs-key-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.ecs-key-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.ecs-key-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.ecs-key-table :deep(.el-table__header th .cell) {
  color: #2c3e50 !important;
  font-weight: 700 !important;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.ecs-key-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.ecs-key-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
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
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.el-input :deep(.el-input__wrapper),
.el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #2c3e50;
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

.ecs-key-table .el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
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
</style>
