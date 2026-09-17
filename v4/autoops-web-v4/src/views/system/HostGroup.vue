<template>
  <div class="group-management">
    <el-card shadow="hover" class="group-card">
      <template #header>
        <div class="card-header">
          <span class="title">堡垒机用户组管理</span>
        </div>
      </template>

      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :inline="true" :model="queryForm" class="search-form">
          <el-form-item label="用户组名称">
            <el-input
              v-model="queryForm.name"
              placeholder="请输入用户组名称"
              size="small"
              clearable
              class="modern-input"
            />
          </el-form-item>
          <el-form-item label="状态" style="width: 200px">
            <el-select
              v-model="queryForm.status"
              placeholder="请选择"
              clearable
              size="small"
              class="modern-select"
              style="width: 100%"
            >
              <el-option label="全部" value="" />
              <el-option label="启用" :value="1" />
              <el-option label="禁用" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              size="small"
              icon="Search"
              @click="handleQuery"
              class="modern-btn primary-btn"
            >
              查询
            </el-button>
            <el-button
              size="small"
              icon="Refresh"
              @click="handleReset"
              class="modern-btn reset-btn"
            >
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 操作按钮区域 -->
      <div class="action-section">
        <el-button
          type="success"
          size="small"
          icon="Plus"
          @click="handleCreate"
          class="modern-btn success-btn"
        >
          新建用户组
        </el-button>
      </div>

      <!-- 表格 -->
      <div class="table-section">
        <el-table
          v-loading="loading"
          :data="groupList"
          stripe
          class="modern-table"
        >
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="name" label="用户组名称" min-width="120" />
          <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
          <el-table-column prop="memberCount" label="成员" width="70" align="center">
            <template #default="{ row }">
              <el-tag size="small">{{ row.memberCount || 0 }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="hostCount" label="主机" width="70" align="center">
            <template #default="{ row }">
              <el-tag type="success" size="small">{{ row.hostCount || 0 }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="创建时间" width="160" />
          <el-table-column label="操作" width="280" fixed="right">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-button type="primary" size="small" @click="handleEdit(row)">
                  编辑
                </el-button>
                <el-button type="info" size="small" @click="handleMembers(row)">
                  成员
                </el-button>
                <el-button type="warning" size="small" @click="handleAuthorization(row)">
                  授权
                </el-button>
                <el-button type="danger" size="small" @click="handleDelete(row)">
                  删除
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <el-pagination
          :current-page="pagination.page"
          :page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
          class="modern-pagination"
        />
      </div>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="600px"
      class="modern-dialog"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        class="dialog-form"
      >
        <el-form-item label="用户组名称" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="请输入用户组名称"
            class="modern-input"
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述"
            class="modern-input"
          />
        </el-form-item>
        <el-form-item label="成员" prop="memberIds">
          <el-select
            v-model="formData.memberIds"
            multiple
            filterable
            remote
            placeholder="搜索并选择用户"
            :remote-method="searchUsers"
            :loading="userLoading"
            style="width: 100%;"
            class="modern-select"
          >
            <el-option
              v-for="user in userOptions"
              :key="user.id"
              :label="`${user.nickname || user.username} (${user.username})`"
              :value="user.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status" class="modern-radio-group">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false" class="modern-btn secondary-btn">取消</el-button>
          <el-button type="primary" :loading="submitLoading" @click="handleSubmit" class="modern-btn primary-btn">
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 成员管理对话框 -->
    <el-dialog
      title="成员管理"
      v-model="memberDialogVisible"
      width="800px"
      class="modern-dialog"
    >
      <div style="margin-bottom: 16px;">
        <el-button type="primary" size="small" @click="handleAddMembers" class="modern-btn primary-btn">
          <el-icon><Plus /></el-icon>
          添加成员
        </el-button>
        <el-button
          type="danger"
          size="small"
          :disabled="selectedMembers.length === 0"
          @click="handleRemoveMembers"
          class="modern-btn danger-btn"
        >
          <el-icon><Delete /></el-icon>
          移除成员 ({{ selectedMembers.length }})
        </el-button>
      </div>
      <el-table
        :data="memberList"
        border
        @selection-change="handleMemberSelectionChange"
        v-loading="memberLoading"
        class="modern-table"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="userId" label="用户ID" width="80" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="120" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column prop="phone" label="手机号" min-width="120" />
      </el-table>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="memberDialogVisible = false" class="modern-btn secondary-btn">关闭</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 添加成员对话框 -->
    <el-dialog
      title="添加成员"
      v-model="addMemberDialogVisible"
      width="600px"
      class="modern-dialog"
    >
      <el-form label-width="100px" class="dialog-form">
        <el-form-item label="选择用户">
          <el-select
            v-model="selectedUserIds"
            multiple
            filterable
            remote
            placeholder="搜索并选择用户"
            :remote-method="searchUsers"
            :loading="userLoading"
            style="width: 100%;"
            class="modern-select"
          >
            <el-option
              v-for="user in userOptions"
              :key="user.id"
              :label="`${user.nickname || user.username} (${user.username})`"
              :value="user.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="addMemberDialogVisible = false" class="modern-btn secondary-btn">取消</el-button>
          <el-button type="primary" @click="handleConfirmAddMembers" class="modern-btn primary-btn">
            确定添加 ({{ selectedUserIds.length }})
          </el-button>
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
  Delete
} from '@element-plus/icons-vue'

export default {
  name: 'UserGroupManagement',
  components: {
    Search,
    Refresh,
    Plus,
    Delete
  },
  data() {
    return {
      loading: false,
      submitLoading: false,
      groupList: [],
      queryForm: {
        name: '',
        status: ''
      },
      pagination: {
        page: 1,
        pageSize: 10,
        total: 0
      },
      dialogVisible: false,
      dialogTitle: '',
      dialogMode: 'create',
      formData: {
        id: null,
        name: '',
        description: '',
        memberIds: [],
        status: 1
      },
      formRules: {
        name: [
          { required: true, message: '请输入用户组名称', trigger: 'blur' }
        ],
        status: [
          { required: true, message: '请选择状态', trigger: 'change' }
        ]
      },
      userLoading: false,
      userOptions: [],
      memberDialogVisible: false,
      memberLoading: false,
      memberList: [],
      currentGroupId: null,
      selectedMembers: [],
      addMemberDialogVisible: false,
      selectedUserIds: []
    }
  },
  mounted() {
    this.fetchList()
  },
  methods: {
    async fetchList() {
      this.loading = true
      try {
        const { data: res } = await this.$api.getUserGroupList({
          ...this.queryForm,
          page: this.pagination.page,
          pageSize: this.pagination.pageSize
        })
        if (res.code === 200) {
          this.groupList = res.data.list || []
          this.pagination.total = res.data.total || 0
        } else {
          this.$message.error(res.message || '获取用户组列表失败')
        }
      } catch (error) {
        this.$message.error('获取用户组列表失败')
      } finally {
        this.loading = false
      }
    },

    handleQuery() {
      this.pagination.page = 1
      this.fetchList()
    },

    handleReset() {
      this.queryForm = {
        name: '',
        status: ''
      }
      this.handleQuery()
    },

    handleSizeChange(size) {
      this.pagination.pageSize = size
      this.fetchList()
    },

    handlePageChange(page) {
      this.pagination.page = page
      this.fetchList()
    },

    handleCreate() {
      this.dialogMode = 'create'
      this.dialogTitle = '新建用户组'
      this.dialogVisible = true
      this.loadUserOptions()
    },

    handleEdit(row) {
      this.dialogMode = 'edit'
      this.dialogTitle = '编辑用户组'
      this.formData = {
        id: row.id,
        name: row.name,
        description: row.description,
        memberIds: row.memberIds || [],
        status: row.status
      }
      this.dialogVisible = true
      this.loadUserOptions()
    },

    async handleSubmit() {
      this.$refs.formRef.validate(async valid => {
        if (!valid) return

        this.submitLoading = true
        try {
          let res
          if (this.dialogMode === 'create') {
            res = await this.$api.createUserGroup(this.formData)
          } else {
            res = await this.$api.updateUserGroup(this.formData)
          }

          if (res.data.code === 200) {
            this.$message.success(this.dialogMode === 'create' ? '创建成功' : '更新成功')
            this.dialogVisible = false
            this.fetchList()
          } else {
            this.$message.error(res.data.message || '操作失败')
          }
        } catch (error) {
          this.$message.error(error.message || '操作失败')
        } finally {
          this.submitLoading = false
        }
      })
    },

    handleDialogClose() {
      this.$refs.formRef.resetFields()
      this.formData = {
        id: null,
        name: '',
        description: '',
        memberIds: [],
        status: 1
      }
      this.userOptions = []
    },

    async handleMembers(row) {
      this.currentGroupId = row.id
      this.memberDialogVisible = true
      await this.fetchGroupMembers(row.id)
    },

    async fetchGroupMembers(groupId) {
      this.memberLoading = true
      try {
        const { data: res } = await this.$api.getUserGroupDetail(groupId)
        if (res.code === 200) {
          this.memberList = res.data.members || []
        } else {
          this.$message.error(res.message || '获取成员列表失败')
        }
      } catch (error) {
        this.$message.error('获取成员列表失败')
      } finally {
        this.memberLoading = false
      }
    },

    handleMemberSelectionChange(selection) {
      this.selectedMembers = selection
    },

    handleAddMembers() {
      this.selectedUserIds = []
      this.addMemberDialogVisible = true
      this.loadUserOptions()
    },

    async handleConfirmAddMembers() {
      if (this.selectedUserIds.length === 0) {
        this.$message.warning('请选择要添加的用户')
        return
      }

      try {
        const { data: res } = await this.$api.addGroupMembers({
          groupId: this.currentGroupId,
          userIds: this.selectedUserIds
        })
        if (res.code === 200) {
          this.$message.success('添加成员成功')
          this.addMemberDialogVisible = false
          await this.fetchGroupMembers(this.currentGroupId)
          await this.fetchList()
        } else {
          this.$message.error(res.message || '添加成员失败')
        }
      } catch (error) {
        this.$message.error('添加成员失败')
      }
    },

    async handleRemoveMembers() {
      if (this.selectedMembers.length === 0) {
        this.$message.warning('请选择要移除的成员')
        return
      }

      try {
        await this.$confirm(
          `确定要移除选中的 ${this.selectedMembers.length} 个成员吗？`,
          '提示',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        const userIds = this.selectedMembers.map(m => m.userId)
        const { data: res } = await this.$api.removeGroupMembers({
          groupId: this.currentGroupId,
          userIds
        })

        if (res.code === 200) {
          this.$message.success('移除成员成功')
          await this.fetchGroupMembers(this.currentGroupId)
          await this.fetchList()
        } else {
          this.$message.error(res.message || '移除成员失败')
        }
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('移除成员失败')
        }
      }
    },

    handleAuthorization(row) {
      this.$router.push({
        name: 'BastionAuthorization',
        query: {
          subjectType: 'group',
          subjectId: row.id
        }
      })
    },

    async handleDelete(row) {
      try {
        await this.$confirm(
          `确定要删除用户组"${row.name}"吗？删除后该用户组的所有授权关系也将被删除。`,
          '警告',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        const { data: res } = await this.$api.deleteUserGroup({ id: row.id })
        if (res.code === 200) {
          this.$message.success('删除成功')
          this.fetchList()
        } else {
          this.$message.error(res.message || '删除失败')
        }
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
        }
      }
    },

    async loadUserOptions() {
      this.userLoading = true
      try {
        const { data: res } = await this.$api.getUserList({
          page: 1,
          pageSize: 100
        })
        if (res.code === 200) {
          this.userOptions = res.data.list || []
        }
      } catch (error) {
        console.error('加载用户列表失败:', error)
      } finally {
        this.userLoading = false
      }
    },

    async searchUsers(query) {
      if (!query) {
        await this.loadUserOptions()
        return
      }

      this.userLoading = true
      try {
        const { data: res } = await this.$api.getUserList({
          keyword: query,
          page: 1,
          pageSize: 20
        })
        if (res.code === 200) {
          this.userOptions = res.data.list || []
        }
      } catch (error) {
        console.error('搜索用户失败:', error)
      } finally {
        this.userLoading = false
      }
    }
  }
}
</script>

<style lang="less" scoped>
.group-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.group-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;

  :deep(.el-card__header) {
    border-bottom: 1px solid rgba(103, 126, 234, 0.1);
    padding: 20px 24px;
  }

  :deep(.el-card__body) {
    padding: 20px 24px;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.search-section {
  margin-bottom: 16px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);

  .search-form {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    width: 100%;

    :deep(.el-form-item) {
      margin-bottom: 0;
    }

    :deep(.el-form-item__label) {
      color: #606266;
      font-weight: 500;
    }
  }
}

.action-section {
  margin-bottom: 16px;
  padding: 12px 0;

  .el-button {
    margin-right: 12px;
  }
}

.table-section {
  margin-bottom: 20px;
  flex: 1;
}

.modern-btn {
  border-radius: 8px;
  padding: 8px 20px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
  }
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

.secondary-btn {
  background: #fff;
  color: #606266;
  border: 1px solid #dcdfe6;

  &:hover {
    border-color: #667eea;
    background: rgba(103, 126, 234, 0.05);
  }
}

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
}

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
        padding: 12px !important;
      }

      &:hover {
        background: rgba(103, 126, 234, 0.05) !important;
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      }
    }
  }
}

.operation-buttons {
  display: flex;
  justify-content: center;
  gap: 4px;
  flex-wrap: wrap;

  .el-button {
    padding: 4px 12px;
    font-size: 12px;
    margin: 2px;
  }
}

.modern-pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;

  :deep(.btn-next),
  :deep(.btn-prev),
  :deep(.el-pager li) {
    border-radius: 6px;
    margin: 0 2px;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
  }
}

.modern-dialog {
  :deep(.el-dialog) {
    background: #ffffff;
    border-radius: 16px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);

    .el-dialog__header {
      background: linear-gradient(135deg, #667eea, #764ba2);
      border-radius: 16px 16px 0 0;
      padding: 20px 24px;

      .el-dialog__title {
        color: #ffffff;
        font-weight: 600;
        font-size: 16px;
      }

      .el-dialog__headerbtn {
        .el-dialog__close {
          color: #ffffff;
          font-size: 20px;

          &:hover {
            color: #f0f0f0;
          }
        }
      }
    }

    .el-dialog__body {
      padding: 24px;
    }

    .el-dialog__footer {
      border-top: 1px solid #f0f0f0;
      padding: 16px 24px;
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
      border-radius: 8px;
      box-shadow: none;

      &:hover {
        border-color: #c0c4cc;
      }

      &.is-focus {
        border-color: #409eff;
        box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
      }
    }

    .el-input__inner,
    .el-textarea__inner {
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
  :deep(.el-radio) {
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
  padding-top: 12px;
  display: flex;
  align-items: center;
  justify-content: flex-end;

  .modern-btn {
    margin-left: 8px;
  }
}

.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

@media (max-width: 768px) {
  .group-management {
    padding: 12px;
  }

  .group-card {
    :deep(.el-card__header) {
      padding: 16px;
    }

    :deep(.el-card__body) {
      padding: 16px;
    }
  }

  .modern-dialog :deep(.el-dialog) {
    width: 90% !important;
    margin: 5vh auto;
  }

  .search-section {
    padding: 12px;

    .search-form {
      gap: 8px;
    }
  }
}
</style>
