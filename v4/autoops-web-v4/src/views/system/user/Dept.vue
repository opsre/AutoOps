<template>
  <div class="dept-container">
    <div class="glass-card">
      <div class="card-header">
        <h2 class="gradient-title">部门管理</h2>
      </div>
      
      <!--搜索-->
      <div class="search-section">
        <el-form :inline="true" :model="queryParams" class="search-form">
          <el-form-item label="部门名称">
            <el-input 
              placeholder="请输入部门名称" 
              clearable 
              v-model="queryParams.deptName"
              @keyup.enter="handleQuery"
              size="small" 
              class="modern-input"
            />
          </el-form-item>
          <el-form-item label="部门状态">
            <el-select 
              placeholder="部门状态"  
              v-model="queryParams.deptStatus" 
              style="width: 150px;" 
              size="small"
              class="modern-select"
            >
              <el-option 
                v-for="item in deptStatusList" 
                :key="item.value" 
                :label="item.label" 
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button 
              type="primary" 
              size="small"
              @click="handleQuery"
              class="modern-btn primary-btn"
            >
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button 
              size="small"
              @click="resetQuery"
              class="modern-btn reset-btn"
            >
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 操作按钮区域 -->
      <div class="action-section">
        <div class="action-buttons">
          <el-button 
            type="success" 
            size="small"
            @click="addDeptDialogVisible = true" 
            v-authority="['base:dept:add']"
            class="modern-btn success-btn"
          >
            <el-icon><Plus /></el-icon>
            新增部门
          </el-button>
          <el-button 
            size="small"
            @click="toggleExpandAll"
            class="modern-btn secondary-btn"
          >
            <el-icon><Sort /></el-icon>
            展开/折叠
          </el-button>
        </div>
      </div>
      
      <!--列表-->
      <div class="table-section">
        <el-table 
          v-if="refreshTable"
          v-loading="loading" 
          :data="deptList" 
          row-key="id" 
          :default-expand-all="isExpandAll"
          :tree-props="{ children: 'children' }"
          class="modern-table"
          :header-cell-style="{ background: 'transparent', color: '#2c3e50', fontWeight: 'bold' }"
          :row-style="{ background: 'rgba(255, 255, 255, 0.05)' }"
        >
          <el-table-column label="部门名称" prop="deptName" />
          <el-table-column label="部门类型" prop="deptType">
            <template v-slot="scope">
              <el-tag 
                v-if="scope.row.deptType === 1" 
                class="modern-tag company-tag"
              >
                公司
              </el-tag>
              <el-tag 
                v-else-if="scope.row.deptType === 2" 
                class="modern-tag center-tag"
              >
                中心
              </el-tag>
              <el-tag 
                v-else-if="scope.row.deptType === 3" 
                class="modern-tag dept-tag"
              >
                部门
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="部门状态" prop="deptStatus">
            <template v-slot="scope">
              <el-tag 
                v-if="scope.row.deptStatus === 1" 
                class="modern-tag success-tag"
              >
                正常
              </el-tag>
              <el-tag 
                v-else-if="scope.row.deptStatus === 2" 
                class="modern-tag danger-tag"
              >
                停用
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="createTime" />
          <el-table-column label="操作" width="150" fixed="right">
            <template v-slot="scope">
              <div class="operation-buttons">
                <el-tooltip content="修改" placement="top">
                  <el-button
                    type="warning"
                    size="small"
                    @click="showEditDeptDialog(scope.row.id)"
                    v-authority="['base:dept:edit']"
                  >
                    <el-icon><Edit /></el-icon>修改
                  </el-button>
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button
                    type="danger"
                    size="small"
                    @click="handleDeptDelete(scope.row)"
                    :disabled="scope.row.deptType == '1' ? true : false"
                    v-authority="['base:dept:delete']"
                  >
                    <el-icon><Delete /></el-icon>删除
                  </el-button>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!--新增部门-->
    <el-dialog 
      title="新增部门" 
      v-model="addDeptDialogVisible" 
      width="35%" 
      @close="addDeptDialogClosed"
      class="modern-dialog"
    >
      <div class="dialog-content">
        <el-form 
          :model="addDeptForm" 
          :rules="addDeptFormRules" 
          ref="addDeptFormRefForm" 
          label-width="90px"
          class="modern-form"
        >
          <el-form-item label="部门类型" prop="deptType">
            <el-radio-group v-model="addDeptForm.deptType" class="modern-radio-group">
              <el-radio :label="1" class="modern-radio">公司</el-radio>
              <el-radio :label="2" class="modern-radio">中心</el-radio>
              <el-radio :label="3" class="modern-radio">部门</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="上级部门" prop="parentId" v-if="addDeptForm.deptType != 1">
            <treeselect 
              :options="optionsDeptList" 
              placeholder="请选择上级部门" 
              v-model="addDeptForm.parentId"
              class="modern-treeselect"
            />
          </el-form-item>
          <el-form-item label="部门名称" prop="deptName">
            <el-input v-model="addDeptForm.deptName" class="modern-input" />
          </el-form-item>
          <el-form-item label="部门状态" prop="deptStatus" >
            <el-radio-group v-model="addDeptForm.deptStatus" class="modern-radio-group">
              <el-radio :label="1" class="modern-radio">正常</el-radio>
              <el-radio :label="2" class="modern-radio">停用</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="addDept" class="modern-btn primary-btn">确 定</el-button>
          <el-button @click="addDeptDialogVisible = false" class="modern-btn secondary-btn">取 消</el-button>
        </div>
      </template>
    </el-dialog>

    <!--修改部门-->
    <el-dialog 
      title="编辑部门" 
      v-model="editDeptDialogVisible" 
      width="35%"
      class="modern-dialog"
    >
      <div class="dialog-content">
        <el-form 
          :model="deptInfo" 
          :rules="editDeptFormRules" 
          ref="editDeptFormRefForm" 
          label-width="90px"
          class="modern-form"
        >
          <el-form-item label="部门类型" prop="deptType">
            <el-radio-group v-model="deptInfo.deptType" class="modern-radio-group">
              <el-radio :label="1" class="modern-radio">公司</el-radio>
              <el-radio :label="2" class="modern-radio">中心</el-radio>
              <el-radio :label="3" class="modern-radio">部门</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="上级部门" prop="parentId" v-if="deptInfo.deptType != 1">
            <treeselect 
              :options="optionsDeptList" 
              placeholder="请选择上级部门" 
              v-model="deptInfo.parentId"
              class="modern-treeselect"
            />
          </el-form-item>
          <el-form-item label="部门名称" prop="deptName">
            <el-input v-model="deptInfo.deptName" class="modern-input" />
          </el-form-item>
          <el-form-item label="部门状态" prop="deptStatus">
            <el-radio-group v-model="deptInfo.deptStatus" class="modern-radio-group">
              <el-radio :label="1" class="modern-radio">正常</el-radio>
              <el-radio :label="2" class="modern-radio">停用</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="editDept" class="modern-btn primary-btn">确 定</el-button>
          <el-button @click="editDeptDialogVisible = false" class="modern-btn secondary-btn">取 消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import Treeselect from 'vue3-treeselect'
import 'vue3-treeselect/dist/vue3-treeselect.css'
import {
  Search,
  Refresh,
  Plus,
  Sort,
  Edit,
  Delete
} from '@element-plus/icons-vue'

export default {
  components: { 
    Treeselect,
    Search,
    Refresh,
    Plus,
    Sort,
    Edit,
    Delete
  },
  data() {
    return {
      deptStatusList: [{
        value: '2',
        label: '停用'
      }, {
        value: '1',
        label: '正常'
      }],
      queryParams: {},
      loading: true,
      deptList: [],
      refreshTable: true,
      isExpandAll: true,
      optionsDeptList: [],
      addDeptDialogVisible: false,
      addDeptFormRules: {
        deptType: [{ required: true, message: "请选择部门类型", trigger: "blur" }],
        deptName: [{ required: true, message: '请输入部门名称', trigger: 'blur' }],
      },
      addDeptForm: {
        deptStatus: 1
      },
      editDeptDialogVisible: false,
      deptInfo: {},
      editDeptFormRules: {
        deptType: [{ required: true, message: "请选择部门类型", trigger: "blur" }],
        deptName: [{ required: true, message: '请输入部门名称', trigger: 'blur' }],
      }
    }
  },
  methods: {
    // 列表
    async getList() {
      this.loading = true
      const { data: res } = await this.$api.queryDeptList(this.queryParams)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.deptList = this.$handleTree.handleTree(res.data, "id")
        this.loading = false
      }
    },
    // 搜索
    handleQuery() {
      this.getList()
    },
    // 重置搜索
    resetQuery() {
      this.queryParams = {}
      this.getList()
      this.$message.success("重置成功")
    },
    // 展开和折叠
    toggleExpandAll() {
      this.refreshTable = false
      this.isExpandAll = !this.isExpandAll
      this.$nextTick(() => {
        this.refreshTable = true
      })
    },
    // 部门下拉列表
    async getDeptVoList() {
      const { data: res } = await this.$api.querySysDeptVoList()
      // console.log(res)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.optionsDeptList = this.$handleTree.handleTree(res.data, "id")
      }
    },
    // 监听新增部门对话框
    addDeptDialogClosed() {
      this.$refs.addDeptFormRefForm.resetFields()
    },
    // 新增
    addDept() {
      this.$refs.addDeptFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.addDept(this.addDeptForm);
        if (res.code !== 200) {
          this.$message.error(res.message)
        } else {
          this.$message.success('新增部门成功')
          this.addDeptDialogVisible = false
          await this.getList()
          await this.getDeptVoList()
        }
      })
    },
    // 展示编辑对话框
    async showEditDeptDialog(id) {
      const { data: res } = await this.$api.deptInfo(id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.deptInfo = res.data
        this.editDeptDialogVisible = true
      }
    },
    // 监听编辑部门
    editDeptDialogClosed() {
      this.$refs.editDeptFormRefForm.resetFields()
    },
    // 修改部门信息并提交
    editDept() {
      this.$refs.editDeptFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.deptUpdate(this.deptInfo)
        if (res.code !== 200) {
          this.$message.error(res.message)
        } else {
          this.editDeptDialogVisible = false
          await this.getList()
          this.$message.success("修改部门成功")
        }
      })
    },
    // 删除部门
    async handleDeptDelete(row) {
      const confirmResult = await this.$confirm('是否确认删除部门为"' + row.deptName + '"的数据项？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消删除')
      }
      const { data: res } = await this.$api.deleteDept(row.id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.getList()
      }
    }
  },
  created() {
    this.getList()
    this.getDeptVoList()
  }
}
</script>

<style lang="less" scoped>
.dept-container {
  min-height: 100vh;
  background: #eef1f6;
  padding: 20px;
  
  .glass-card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 16px;
    border: 1px solid rgba(255, 255, 255, 0.2);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    padding: 24px;
    transition: all 0.3s ease;
    box-sizing: border-box;
    width: 100%;
    overflow: hidden;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
    }
  }
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 12px;
    border-bottom: 1px solid rgba(103, 126, 234, 0.1);
    
    .gradient-title {
      color: #2c3e50;
      background: linear-gradient(45deg, #667eea, #764ba2);
      background-clip: text;
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      font-size: 20px;
      font-weight: 600;
      margin: 0;
    }
  }
  
  .search-section {
    margin-bottom: 24px;
    padding: 20px;
    background: rgba(103, 126, 234, 0.05);
    border-radius: 12px;
    border: 1px solid rgba(103, 126, 234, 0.1);
    overflow: hidden;
    margin-left: 0;
    margin-right: 0;
    
    .search-form {
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      gap: 16px;
      width: 100%;
      box-sizing: border-box;
      
      :deep(.el-form-item) {
        margin-bottom: 0;
        flex-shrink: 0;
      }
      
      :deep(.el-form-item__label) {
        color: #606266;
        font-weight: 500;
      }

      :deep(.el-form-item:last-child) {
        display: flex;
        gap: 12px;
      }
    }
  }
  
  // 操作按钮区域样式
  .action-section {
    margin-bottom: 5px;
    margin-left: 0;
    margin-right: 0;
    padding-left: 20px;
    
    .action-buttons {
      display: flex;
      gap: 12px;
      align-items: center;
      
      .modern-btn {
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        
        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
        }
      }
    }
  }
  
  .table-section {
    background: transparent;
    border-radius: 12px;
    padding: 20px;
    overflow: hidden;
    min-width: 0;
    margin-left: 0;
    margin-right: 0;
  }
  
  // 现代化按钮样式
  .modern-btn {
    border-radius: 8px;
    padding: 8px 20px;
    font-weight: 500;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border: none;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  .modern-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
  }

  .primary-btn {
    background: linear-gradient(45deg, #409EFF, #66B3FF);
    color: white;
  }

  .reset-btn {
    background: linear-gradient(45deg, #E6A23C, #EEBE77);
    color: white;
  }

  .success-btn {
    background: linear-gradient(45deg, #67C23A, #85CE61);
    color: white;
  }

  .secondary-btn {
    background: linear-gradient(45deg, #909399, #B1B3B8);
    color: white;
  }
  
  // 输入框样式
  .modern-input {
    :deep(.el-input__wrapper) {
      background: rgba(255, 255, 255, 0.8);
      border: 1px solid rgba(103, 126, 234, 0.2);
      border-radius: 8px;
      box-shadow: none;
      transition: all 0.3s ease;
    }

    :deep(.el-input__wrapper):hover {
      border-color: #c0c4cc;
    }

    :deep(.el-input__wrapper.is-focus) {
      border-color: #667eea;
      box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
      background: rgba(255, 255, 255, 1);
    }

    :deep(.el-input__inner) {
      background: transparent;
      border: none;
      color: #2c3e50;

      &::placeholder {
        color: rgba(44, 62, 80, 0.6);
      }
    }
  }
  
  // 选择框样式
  .modern-select {
    :deep(.el-input__wrapper) {
      background: rgba(255, 255, 255, 0.8);
      border: 1px solid rgba(103, 126, 234, 0.2);
      border-radius: 8px;
      box-shadow: none;
      transition: all 0.3s ease;
    }

    :deep(.el-input__wrapper):hover {
      border-color: #c0c4cc;
    }

    :deep(.el-input__wrapper.is-focus) {
      border-color: #667eea;
      box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
      background: rgba(255, 255, 255, 1);
    }

    :deep(.el-input__inner) {
      background: transparent;
      border: none;
      color: #2c3e50;

      &::placeholder {
        color: rgba(44, 62, 80, 0.6);
      }
    }
    
    :deep(.el-input__suffix-inner) {
      color: #606266;
    }
  }
  
  // 表格样式
  .modern-table {
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
    
    :deep(.el-table) {
      background: transparent;
      
      &::before {
        display: none;
      }
      
      th {
        background: linear-gradient(135deg, #667eea, #764ba2) !important;
        border-bottom: none;
        color: #2c3e50 !important;
        font-weight: 700 !important;
        text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
        
        .cell {
          color: #2c3e50 !important;
          font-weight: 700 !important;
          text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
        }
      }
      
      td {
        background: rgba(255, 255, 255, 0.05) !important;
        border-bottom: 1px solid rgba(103, 126, 234, 0.1);
        color: #2c3e50;
      }
      
      .el-table__row {
        transition: all 0.3s ease;
        
        td {
          padding: 8px 12px !important;
          height: 40px !important;
        }
        
        &:hover {
          background: rgba(103, 126, 234, 0.05) !important;
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }
      }
    }
  }
  
  // 标签样式
  .modern-tag {
    border-radius: 20px;
    border: none;
    padding: 4px 12px;
    font-weight: 500;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    }
    
    &.company-tag {
      background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
      color: #8b4513;
    }
    
    &.center-tag {
      background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
      color: #2d5a87;
    }
    
    &.dept-tag {
      background: linear-gradient(135deg, #fad0c4 0%, #ffd1ff 100%);
      color: #8b2b85;
    }
    
    &.success-tag {
      background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%);
      color: #0f5132;
    }
    
    &.danger-tag {
      background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
      color: #721c24;
    }
  }
  
  // 操作按钮区域样式
  .operation-buttons {
    display: flex;
    justify-content: center;
    gap: 8px;
    
    .el-button {
      width: 32px;
      height: 32px;
      padding: 0;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: all 0.3s ease;
      
      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
      }
      
      .el-icon {
        font-size: 14px;
      }
    }
  }
  
  // 操作按钮样式
  .action-btn {
    border-radius: 6px;
    padding: 4px 8px;
    margin: 0 2px;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-1px);
    }
    
    &.edit-btn {
      color: #40a9ff;
      
      &:hover {
        background: rgba(64, 169, 255, 0.1);
        color: #1890ff;
      }
    }
    
    &.delete-btn {
      color: #ff7875;
      
      &:hover {
        background: rgba(255, 120, 117, 0.1);
        color: #f5222d;
      }
      
      &:disabled {
        color: rgba(255, 255, 255, 0.3);
      }
    }
  }
  
  // 对话框样式
  :deep(.modern-dialog) {
    .el-dialog {
      background: #ffffff;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }
    
    .el-dialog__header {
      background: #ffffff;
      border-bottom: 1px solid #f0f0f0;
      padding: 16px 20px;
      
      .el-dialog__title {
        color: #303133;
        font-weight: 600;
        font-size: 16px;
      }
      
      .el-dialog__headerbtn {
        .el-dialog__close {
          color: #909399;
          font-size: 16px;
          
          &:hover {
            color: #303133;
          }
        }
      }
    }
    
    .el-dialog__body {
      padding: 24px;
    }
    
    .el-dialog__footer {
      border-top: 1px solid #f0f0f0;
      padding: 20px 24px;
    }
  }
  
  .dialog-content {
    .modern-form {
      :deep(.el-form-item__label) {
        color: #606266;
        font-weight: 500;
      }
      
      :deep(.el-form-item__content) {
        .el-input__wrapper {
          background: #ffffff;
          border: 1px solid #dcdfe6;
          border-radius: 4px;
          box-shadow: none;
          
          &:hover {
            border-color: #c0c4cc;
          }
          
          &.is-focus {
            border-color: #409eff;
            box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
          }
        }
        
        .el-input__inner {
          background: transparent;
          border: none;
          color: #606266;
        }
      }
    }
  }
  
  // 单选框样式
  .modern-radio-group {
    :deep(.el-radio) {
      .el-radio__input {
        .el-radio__inner {
          background: #ffffff;
          border-color: #dcdfe6;
          
          &:hover {
            border-color: #409eff;
          }
        }
        
        &.is-checked {
          .el-radio__inner {
            background: #409eff;
            border-color: #409eff;
          }
        }
      }
      
      .el-radio__label {
        color: #606266;
      }
    }
  }
  
  // 树选择器样式
  .modern-treeselect {
    :deep(.vue-treeselect__control) {
      background: #ffffff;
      border: 1px solid #dcdfe6;
      border-radius: 4px;
      
      .vue-treeselect__placeholder,
      .vue-treeselect__single-value {
        color: #606266;
      }
      
      &:hover {
        border-color: #409eff;
      }
      
      &.vue-treeselect--focused {
        border-color: #409eff;
        box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
      }
    }
    
    :deep(.vue-treeselect__menu) {
      background: #ffffff;
      border: 1px solid #dcdfe6;
      border-radius: 4px;
      
      .vue-treeselect__option {
        color: #606266;
        
        &:hover {
          background: #f5f7fa;
        }
        
        &.vue-treeselect__option--selected {
          background: #409eff;
          color: #ffffff;
        }
      }
    }
  }
  
  .dialog-footer {
    display: flex;
    justify-content: center;
    gap: 15px;
  }
  
  // 响应式设计
  @media (max-width: 1200px) {
    .search-form {
      flex-direction: column;
      align-items: stretch;
      
      :deep(.el-form-item) {
        width: 100%;
      }
    }
    
    .action-section {
      flex-direction: column;
      gap: 10px;
      
      .modern-btn {
        width: 100%;
      }
    }
  }
  
  @media (max-width: 768px) {
    padding: 10px;
    
    .glass-card {
      padding: 16px;
      border-radius: 12px;
    }
    
    .card-header .gradient-title {
      font-size: 18px;
    }
    
    .search-section, .action-section, .table-section {
      padding: 12px;
    }
    
    .search-section {
      margin-bottom: 16px;
    }
    
    .action-section {
      margin-bottom: 16px;
    }
    
    .modern-table {
      :deep(.el-table) {
        font-size: 12px;
      }
    }
  }
  
  @media (max-width: 480px) {
    padding: 8px;
    
    .glass-card {
      padding: 12px;
    }
    
    .card-header .gradient-title {
      font-size: 16px;
    }
    
    .search-section, .action-section, .table-section {
      padding: 8px;
    }
    
    .modern-table {
      :deep(.el-table) {
        font-size: 11px;
      }
    }
    
    .modern-input,
    .modern-select {
      max-width: 100%;
      box-sizing: border-box;
    }
  }
}

// 全局对话框遮罩样式
:deep(.el-overlay) {
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(5px);
}
</style>
