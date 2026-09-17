<template>
  <div class="post-management">
    <div class="glass-card main-card">
      <!-- 卡片标题 -->
      <div class="card-header">
        <h1 class="gradient-title">岗位管理</h1>
      </div>
      
      <!-- 搜索区域 -->
      <div class="search-section">
        <el-form :model="queryParams" :inline="true" class="search-form">
        <el-form-item label="岗位名称" prop="postName">
          <el-input 
            placeholder="请输入岗位名称" 
            clearable 
            size="small"
            class="modern-input"
            v-model="queryParams.postName">
          </el-input>
        </el-form-item>
        <el-form-item label="岗位状态" prop="postStatus">
          <el-select 
            v-model="queryParams.postStatus" 
            placeholder="岗位状态" 
            size="small"
            style="width: 150px;" 
            class="modern-select">
            <el-option
                v-for="item in postStatusList" 
                :key="item.value" 
                :label="item.label" 
                :value="item.value"/>
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
            size="small"
            @click="resetQuery"
            class="modern-btn reset-btn">
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
            @click="handleAddButtonClick" 
            v-authority="['base:post:add']"
            class="modern-btn success-btn">
            <el-icon><Plus /></el-icon>
            新增
          </el-button>
          <el-button 
            type="danger" 
            size="small"
            :disabled="multiple"
            @click="batchHandleDelete" 
            v-authority="['base:post:delete']"
            class="modern-btn danger-btn">
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </div>
      </div>

      <!-- 表格区域 -->
      <div class="table-section">
        <el-table 
        class="modern-table"
        v-loading="loading" 
        :data="postList" 
        @selection-change="handleSelectionChange">
        <el-table-column type="selection" />
        <el-table-column label="ID" v-if="false" prop="id" />
        <el-table-column label="岗位名称" prop="postName" />
        <el-table-column label="岗位编码" prop="postCode" />
        <el-table-column label="岗位状态" prop="postStatus">
          <template v-slot="scope">
            <el-switch 
              v-model="scope.row.postStatus" 
              :active-value="1" 
              :inactive-value="2" 
              active-color="#667eea"
              inactive-color="#F5222D" 
              active-text="启用" 
              inactive-text="停用" 
              class="modern-switch"
              @change="postUpdateStatus(scope.row)">
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="createTime" />
        <el-table-column label="描述" prop="remark" />
        <el-table-column label="操作" width="180" fixed="right">
          <template v-slot="scope">
            <div class="operation-buttons">
              <el-tooltip content="编辑" placement="top">
                <el-button
                  type="warning"
                  size="small"
                  @click="handleUpdate(scope.row.id)"
                  v-authority="['base:post:edit']"
                >
                  <el-icon><Edit /></el-icon>编辑
                </el-button>
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button
                  type="danger"
                  size="small"
                  @click="handleDelete(scope.row.id)"
                  v-authority="['base:post:delete']"
                >
                  <el-icon><Delete /></el-icon>删除
                </el-button>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!--分页-->
      <el-pagination 
        class="modern-pagination"
        @size-change="handleSizeChange" 
        @current-change="handleCurrentChange"
        :current-page="queryParams.pageNum" 
        :page-sizes="[10, 50, 100, 500, 1000]" 
        :page-size="queryParams.pageSize"
        layout="total, sizes, prev, pager, next, jumper" 
        :total="total">
      </el-pagination>
      </div>
    </div>

    <!--新增对话框-->
    <el-dialog 
      title="新增岗位" 
      v-model="addPostDialogVisible" 
      width="30%" 
      class="modern-dialog"
      @close="addPostDialogClosed">
      <el-form 
        label-width="80px" 
        ref="addPostFormRefForm" 
        :rules="addPostFormRules" 
        :model="addPostForm"
        class="dialog-form">
        <el-form-item label="岗位名称" prop="postName">
          <el-input 
            placeholder="请输入岗位名称" 
            class="modern-input"
            v-model="addPostForm.postName" />
        </el-form-item>
        <el-form-item label="岗位编码" prop="postCode">
          <el-input 
            placeholder="请输入岗位编码" 
            class="modern-input"
            v-model="addPostForm.postCode" />
        </el-form-item>
        <el-form-item label="岗位状态" prop="postStatus">
          <el-radio-group v-model="addPostForm.postStatus" class="modern-radio-group">
            <el-radio :label="1" class="modern-radio">启用</el-radio>
            <el-radio :label="2" class="modern-radio">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="岗位描述" prop="remark">
          <el-input 
            placeholder="请输入岗位描述" 
            class="modern-input"
            v-model="addPostForm.remark" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button 
            type="primary" 
            class="modern-btn primary-btn"
            @click="addPost">确定</el-button>
          <el-button 
            class="modern-btn secondary-btn"
            @click="addPostDialogVisible = false">取消</el-button>
        </div>
      </template>
    </el-dialog>

    <!--编辑对话框-->
    <el-dialog 
      title="编辑岗位" 
      v-model="editPostDialogVisible" 
      width="30%" 
      class="modern-dialog"
      @close="editPostDialogClosed">
      <el-form 
        label-width="80px" 
        ref="editPostFormRefForm" 
        :rules="editPostFormRules" 
        :model="editPostForm"
        class="dialog-form">
        <el-form-item label="岗位名称" prop="postName">
          <el-input 
            placeholder="请输入岗位名称" 
            class="modern-input"
            v-model="editPostForm.postName" />
        </el-form-item>
        <el-form-item label="岗位编码" prop="postCode">
          <el-input 
            placeholder="请输入岗位编码" 
            class="modern-input"
            v-model="editPostForm.postCode" />
        </el-form-item>
        <el-form-item label="岗位状态" prop="postStatus">
          <el-radio-group v-model="editPostForm.postStatus" class="modern-radio-group">
            <el-radio :label="1" class="modern-radio">启用</el-radio>
            <el-radio :label="2" class="modern-radio">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="岗位描述" prop="remark">
          <el-input 
            placeholder="请输入岗位描述" 
            class="modern-input"
            v-model="editPostForm.remark" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button 
            type="primary" 
            class="modern-btn primary-btn"
            @click="editPost">确定</el-button>
          <el-button 
            class="modern-btn secondary-btn"
            @click="editPostDialogVisible = false">取消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {
  Search,
  Refresh,
  Plus,
  Edit,
  Delete
} from '@element-plus/icons-vue'

export default {
  components: {
    Search,
    Refresh,
    Plus,
    Edit,
    Delete
  },
  data() {
    return {
      queryParams: {},
      postStatusList: [{
        value: '1',
        label: '启用'
      }, {
        value: '2',
        label: '停用'
      }],
      loading: true,
      postList: [],
      total: 0,
      addPostDialogVisible: false,
      addPostFormRules: {
        postName: [{ required: true, message: '请输入岗位名称', trigger: 'blur' }],
        postCode: [{ required: true, message: '请输入岗位标识', trigger: 'blur' }],
        postStatus: [{ required: true, message: '请输入岗位状态', trigger: 'blur' }]
      },
      addPostForm: {
        postName: '',
        postCode: '',
        postStatus: 1,
        remark: ''
      },
      editPostDialogVisible: false,
      editPostForm: {},
      editPostFormRules: {
        postName: [{ required: true, message: '请输入岗位名称', trigger: 'blur' }],
        postCode: [{ required: true, message: '请输入岗位标识', trigger: 'blur' }],
        postStatus: [{ required: true, message: '请输入岗位状态', trigger: 'blur' }]
      },
      ids: [],
      single: true,
      multiple: true
    }
  },
  methods: {
    // 新增岗位方法
    handleAddButtonClick() {
      console.log('新增岗位按钮被点击');
      this.addPostDialogVisible = true;
    },
    // 获取列表
    async getPostList() {
      this.loading = true
      const { data: res } = await this.$api.queryPostList(this.queryParams)  // 调用api
      // console.log("res数据:", res)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.postList = res.data.list
        this.total = res.data.total
        this.loading = false
      }
    },
    // 搜索
    handleQuery() {
      this.getPostList()
    },
    // 重置
    resetQuery() {
      this.queryParams = {}
      this.getPostList()
      this.$message.success("重置成功")
    },
    // pageSize
    handleSizeChange(newSize) {
      this.queryParams.pageSize = newSize
      this.getPostList()
    },
    // pageNum
    handleCurrentChange(newPage) {
      this.queryParams.pageNum = newPage
      this.getPostList()
    },
    // 岗位状态修改
    async postUpdateStatus(row) {
      let text = row.postStatus === 2 ? "停用" : "启用"
      const confirmResult = await this.$confirm('确认要"' + text + '""' + row.postName + '"岗位吗?', "警告", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).catch(err => err)
      if (confirmResult != 'confirm') {
        await this.getPostList()
        return this.$message.info('已取消修改')
      }
      await this.$api.updatePostStatus(row.id, row.postStatus)
      return this.$message.success(text + "成功")
      // eslint-disable-next-line no-unreachable
      await this.getPostList()
    },
    // 监听对话框的关闭
    addPostDialogClosed() {
     // console.log('Add dialog closed');
      this.$refs.addPostFormRefForm.resetFields()
    },
    // 新增操作
    addPost() {
      this.$refs.addPostFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.addPost(this.addPostForm)
        if (res.code !== 200) {
          this.$message.error(res.message)
        } else {
          this.$message.success("新增岗位成功")
          this.addPostDialogVisible = false
          await this.getPostList()
        }
      })
    },
    // 显示编辑对话框
    async handleUpdate(id) {
      const { data: res } = await this.$api.postInfo(id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.editPostForm = res.data
        this.editPostDialogVisible = true
      }
    },
    // 监听编辑岗位对话框
    editPostDialogClosed() {
      this.$refs.editPostFormRefForm.resetFields()
    },
    // 修改岗位
    editPost() {
      this.$refs.editPostFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.updatePost(this.editPostForm)
        if (res.code !== 200) {
          this.$message.error(res.message)
        } else {
          this.$message.success("修改岗位成功")
          this.editPostDialogVisible = false
          await this.getPostList()
        }
      })
    },
    // 删除岗位
    async handleDelete(id) {
      const confirmResult = await this.$confirm('是否确认删除岗位编号为"' + id + '"的数据项？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消删除')
      }
      const { data: res } = await this.$api.deleteSysPost(id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.getPostList()
      }
    },
    // 多选框选中数据
    handleSelectionChange(selection) {
      this.ids = selection.map(item => item.id);
      this.single = selection.length != 1;
      this.multiple = !selection.length;
    },
    // 批量删除
    async batchHandleDelete() {
      const postIds = this.ids
      const confirmResult = await this.$confirm('是否确认删除岗位编号为"' + postIds + '"的数据项？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消删除')
      }
      const { data: res } = await this.$api.batchDeleteSysPost(postIds)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.getPostList()
      }
    }
  },
  created() {
    this.getPostList()
  },
}
</script>

<style lang="less" scoped>
.post-management {
  min-height: 100vh;
  background: #eef1f6;
  padding: 20px;
  
  .main-card {
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

  .danger-btn {
    background: linear-gradient(45deg, #FF4D4F, #FF7875);
    color: white;
    
    &:hover:not(:disabled) {
      background: linear-gradient(45deg, #FF2629, #FF5E61);
    }
    
    &:disabled {
      background: linear-gradient(45deg, #ccc 0%, #bbb 100%);
      cursor: not-allowed;
      transform: none;
      box-shadow: none;
    }
  }

  // 现代化输入框样式
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

  // 现代化选择器样式
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

  // 现代化日期选择器样式
  .modern-date-picker {
    width: 190px;
    
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

  // 现代化表格样式
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

  // 动作按钮样式
  .action-btn {
    border-radius: 20px;
    padding: 4px 12px;
    margin: 0 4px;
    transition: all 0.3s ease;
    border: 1px solid transparent;
    
    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }
    
    &.edit-btn {
      color: #667eea;
      background: rgba(102, 126, 234, 0.1);
      border-color: rgba(102, 126, 234, 0.2);
      
      &:hover {
        background: rgba(102, 126, 234, 0.2);
        color: #5a6fd8;
      }
    }
    
    &.delete-btn {
      color: #ff6b6b;
      background: rgba(255, 107, 107, 0.1);
      border-color: rgba(255, 107, 107, 0.2);
      
      &:hover {
        background: rgba(255, 107, 107, 0.2);
        color: #ff5252;
      }
    }
  }

  // 现代化开关样式
  .modern-switch {
    /deep/ .el-switch__core {
      border-radius: 20px;
      height: 24px;
      
      &::after {
        border-radius: 50%;
        width: 20px;
        height: 20px;
      }
    }
  }

  // 现代化分页样式
  .modern-pagination {
    margin-top: 24px;
    display: flex;
    justify-content: center;
    
    /deep/ .el-pagination__total,
    /deep/ .el-pagination__jump,
    /deep/ .el-select .el-input__inner,
    /deep/ .el-pagination__editor.el-input__inner {
      color: #ffffff;
    }
    
    /deep/ .btn-prev,
    /deep/ .btn-next {
      background: rgba(255, 255, 255, 0.1);
      border: 1px solid rgba(255, 255, 255, 0.2);
      border-radius: 8px;
      color: #ffffff;
      transition: all 0.3s ease;
      
      &:hover {
        background: rgba(255, 255, 255, 0.2);
        transform: translateY(-1px);
      }
      
      &:disabled {
        background: rgba(255, 255, 255, 0.05);
        color: rgba(255, 255, 255, 0.3);
      }
    }
    
    /deep/ .el-pager li {
      background: rgba(255, 255, 255, 0.1);
      border: 1px solid rgba(255, 255, 255, 0.2);
      border-radius: 8px;
      color: #ffffff;
      margin: 0 4px;
      transition: all 0.3s ease;
      
      &:hover {
        background: rgba(255, 255, 255, 0.2);
        transform: translateY(-1px);
      }
      
      &.active {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: #ffffff;
        border-color: transparent;
      }
    }
  }

  // 现代化对话框样式
  .modern-dialog {
    :deep(.el-dialog) {
      background: #ffffff;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      
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
  }

  .dialog-form {
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
        
        &::placeholder {
          color: #c0c4cc;
        }
      }
    }
  }
  

  .modern-radio-group {
    /deep/ .el-radio {
      margin-right: 30px;
      
      .el-radio__label {
        color: #333;
        font-weight: 500;
      }
      
      .el-radio__input.is-checked .el-radio__inner {
        background-color: #667eea;
        border-color: #667eea;
      }
      
      &:hover .el-radio__inner {
        border-color: #667eea;
      }
    }
  }

  .dialog-footer {
    text-align: center;
    padding-top: 20px;
    
    .modern-btn {
      margin: 0 8px;
    }
  }

}

// 响应式设计
@media (max-width: 768px) {
  .post-management {
    padding: 12px;
    
    .page-header .gradient-title {
      font-size: 28px;
    }
    
    .glass-card {
      border-radius: 12px;
    }
    
    .search-card,
    .content-card {
      padding: 16px;
    }
    
    .modern-dialog /deep/ .el-dialog {
      width: 90% !important;
      margin: 5vh auto;
    }
  }
}
</style>
