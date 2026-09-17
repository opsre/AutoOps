<template>
  <div class="accountauth-management">
    <el-card shadow="hover" class="accountauth-card">
      <template #header>
        <div class="card-header">
          <span class="title">账号认证管理</span>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :model="queryParams" :inline="true" class="search-form" ref="queryForm">
      <el-form-item label="账号别名" prop="alias">
        <el-input placeholder="请输入账号别名" clearable size="small" v-model="queryParams.alias"
                  @keyup.enter.native="handleQuery" />
      </el-form-item>
      <el-form-item label="账号类型" prop="type" style="width: 200px">
        <el-select v-model="queryParams.type" placeholder="请选择账号类型" size="small" clearable>
          <el-option label="Jenkins" :value="4" />
          <el-option label="通用账号" :value="5" />
        </el-select>
      </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="Search" size="small" @click="handleQuery">搜索</el-button>
            <el-button type="warning" icon="Refresh" size="small" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    
      <!-- 操作按钮 -->
      <div class="action-section">
        <el-button type="success" v-authority="['config:common:add']" plain icon="Plus" size="small"
                   @click="showAddDialog">新增账号</el-button>
                   
      </div>
    
      <!-- 列表区域 -->
      <div class="table-section">
        <el-table stripe v-loading="loading" :data="accountList" class="accountauth-table">
      <el-table-column label="账号别名" prop="alias">
        <template v-slot="scope">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img src="@/assets/image/账号.svg" style="width: 16px; height: 16px"/>
            <span>{{ scope.row.alias }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="服务地址" prop="host"  width="320" >
        <template v-slot="scope">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img src="@/assets/image/url.svg" style="width: 16px; height: 16px"/>
            <span>{{ scope.row.host }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="用户名" prop="name">
        <template v-slot="scope">
          <div style="display: flex; align-items: center; gap: 8px;">
            <img src="@/assets/image/ren.svg" style="width: 16px; height: 16px"/>
            <span>{{ scope.row.name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="账号类型" width="120">
        <template v-slot="scope">
          <el-tag :type="scope.row.type === 4 ? 'success' :
                         (scope.row.type === 5 ? 'warning' : 'info')">
            {{ scope.row.type === 4 ? 'Jenkins' :
               (scope.row.type === 5 ? 'Zabbix' : '通用账号') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="createdAt" />
      <el-table-column label="备注" prop="remark" />
      <el-table-column label="操作" class-name="small-padding fixed-width" width="380">
        <template v-slot="scope">
          <div class="operation-buttons">
            <el-tooltip content="测试连接" placement="top">
              <el-button size="small" v-authority="['config:common:test']" type="success" icon="Connection" @click="handleTest(scope.row)">测试连接</el-button>
            </el-tooltip>
            <el-tooltip content="修改" placement="top">
              <el-button size="small" v-authority="['config:common:edit']" type="primary" icon="Edit" @click="showEditDialog(scope.row)">修改</el-button>
            </el-tooltip>
            <el-tooltip content="删除" placement="top">
              <el-button size="small" v-authority="['config:common:delete']" type="danger" icon="Delete" @click="handleDelete(scope.row)">删除</el-button>
            </el-tooltip>
            <el-tooltip content="解密" placement="top">
              <el-button size="small" v-authority="['config:common:decrypt']" type="warning" icon="Key" @click="handleDecrypt(scope.row)">解密</el-button>
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
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="40%" :modal="false">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="账号别名" prop="alias">
          <el-input v-model="formData.alias" placeholder="请输入账号别名"></el-input>
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item label="服务地址" prop="host">
              <el-input v-model="formData.host" placeholder="192.168.1.1:3306 (无需协议)"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="服务类型" prop="type">
              <el-select v-model="formData.type" placeholder="请选择服务类型">
                <el-option label="Jenkins" :value="4" />
                <el-option label="Zabbix" :value="5" />
                <el-option label="通用账号" :value="6" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户名" prop="name">
              <el-input v-model="formData.name" placeholder="请输入用户名"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码" prop="password">
              <el-input v-model="formData.password" show-password placeholder="请输入密码"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="formData.remark" type="textarea" placeholder="请输入备注信息"></el-input>
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
        alias: '',
        type: undefined,
        pageNum: 1,
        pageSize: 10
      },
      loading: false,
      accountList: [],
      total: 0,
      dialogVisible: false,
      dialogTitle: '',
      formData: {
        id: '',
        alias: '',
        host: '',
        name: '',
        password: '',
        type: undefined,
        remark: ''
      },
      formRules: {
        alias: [{ required: true, message: '请输入账号别名', trigger: 'blur' }],
        host: [{ required: true, message: '请输入服务地址(格式: IP或域名:端口)', trigger: 'blur' }],
        name: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
        type: [{ required: true, message: '请选择服务类型', trigger: 'change' }]
      }
    }
  },
  methods: {
    // 获取账号列表
    async getList() {
      this.loading = true
      try {
        const { data: res } = await API.listAccountAuth({
          page: this.queryParams.pageNum,
          pageSize: this.queryParams.pageSize,
          alias: this.queryParams.alias || undefined,
          type: this.queryParams.type || undefined
        })
        if (res.code === 200) {
          this.accountList = res.data?.list || []
          this.total = res.data?.total || 0
        } else {
          this.$message.error(res.message || '获取账号列表失败')
        }
      } catch (error) {
        console.error('获取账号列表失败:', error)
        this.$message.error('获取账号列表失败')
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
        alias: '',
        type: undefined,
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
      this.dialogTitle = '创建账号'
      this.$nextTick(() => {
        this.formData = {
          id: '',
          alias: '',
          host: '',
          name: '',
          password: '',
          type: undefined,
          remark: ''
        }
        this.dialogVisible = true
      })
    },
    
    // 显示编辑对话框
    showEditDialog(row) {
      this.dialogTitle = '修改账号'
      this.$nextTick(() => {
        this.formData = {
          id: row.id,
          alias: row.alias,
          host: row.host,
          name: row.name,
          password: '',
          type: row.type,
          remark: row.remark,
          createdAt: row.createdAt,  // 保留创建时间
          updatedAt: row.updatedAt   // 保留更新时间
        }
        this.dialogVisible = true
      })
    },
    
    // 提交表单
    async submitForm() {
      try {
        await this.$refs.formRef.validate()

        const formData = {
          ...this.formData,
          type: Number(this.formData.type),
          id: this.formData.id ? Number(this.formData.id) : undefined,
          createdAt: this.formData.createdAt, // 保留创建时间
          updatedAt: this.formData.updatedAt  // 保留更新时间
        }

        let res
        if (formData.id) {
          // 更新 - 仅发送可修改字段
          const { createdAt, updatedAt, ...updateData } = formData
          res = await API.updateAccountAuth({
            id: updateData.id,
            alias: updateData.alias,
            host: updateData.host,
            name: updateData.name,
            password: updateData.password,
            type: updateData.type,
            remark: updateData.remark
          })
        } else {
          // 新增 - 确保不发送ID字段
          const { id, ...createData } = formData
          res = await API.createAccountAuth(createData)
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
    
    // 删除账号
    // 测试账号连接
    async handleTest(row) {
      const loading = this.$loading({
        lock: true,
        text: `正在测试 ${row.alias} 连接...`,
        background: 'rgba(0, 0, 0, 0.7)'
      })

      try {
        const res = await API.testAccountAuth(row.id)
        loading.close()

        if (res.data.code === 200 && res.data.data) {
          const { success, message } = res.data.data

          if (success) {
            this.$notify({
              title: '连接成功',
              message: `账号 "${row.alias}" ${message || '连接测试通过'}`,
              type: 'success',
              duration: 3000
            })
          } else {
            this.$notify({
              title: '连接失败',
              message: `账号 "${row.alias}" ${message || '连接测试失败'}`,
              type: 'error',
              duration: 5000
            })
          }
        } else {
          this.$message.error(res.data.message || '测试失败')
        }
      } catch (error) {
        loading.close()
        console.error('测试连接失败:', error)
        this.$notify({
          title: '测试失败',
          message: error.response?.data?.message || error.message || '网络请求失败',
          type: 'error',
          duration: 5000
        })
      }
    },

    // 解密密码
    async handleDecrypt(row) {
      try {
        const res = await API.decryptPassword({}, {
          params: { id: row.id }
        })
        if (res.data.code === 200) {
          this.$alert(`
            <div>
              <p>账号: ${row.alias}</p>
              <p>密码: <span style="color: #1890ff; font-weight: bold; font-size: 18px;">${res.data.data.password}</span></p>
            </div>
          `, '解密结果', {
            confirmButtonText: '确定',
            customClass: 'decrypt-result-alert',
            dangerouslyUseHTMLString: true
          })
        } else {
          this.$message.error(res.data.message || '解密失败')
        }
      } catch (error) {
        console.error('解密失败:', error)
        this.$message.error('解密失败: ' + (error.response?.data?.message || error.message))
      }
    },

    // 删除账号
    async handleDelete(row) {
      try {
        await this.$confirm(`确定删除账号"${row.alias}"?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        const res = await API.deleteAccountAuth(row.id)
        if (res.data.code === 200) {
          this.$message.success('删除成功')
          this.getList()
        } else {
          this.$message.error(res.data.message || '删除失败')
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
.accountauth-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.accountauth-card {
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

.accountauth-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.accountauth-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.accountauth-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.accountauth-table :deep(.el-table__header th .cell) {
  color: #2c3e50 !important;
  font-weight: 700 !important;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.accountauth-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.accountauth-table :deep(.el-table__row:hover) {
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

.accountauth-table .el-loading-mask {
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
