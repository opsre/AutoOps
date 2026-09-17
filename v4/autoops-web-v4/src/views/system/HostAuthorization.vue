<template>
  <div class="authorization-management">
    <el-card shadow="hover" class="auth-card">
      <template #header>
        <div class="card-header">
          <span class="title">堡垒机主机授权</span>
        </div>
      </template>

      <!-- 授权对象选择 -->
      <div class="subject-selector">
        <div class="selector-content">
          <el-select
            v-model="selectedSubjectId"
            filterable
            placeholder="请选择用户组"
            style="width: 300px"
            class="modern-select"
            @change="handleSubjectChange"
          >
            <el-option
              v-for="group in groupOptions"
              :key="group.id"
              :label="group.name"
              :value="group.id"
            >
              <span>{{ group.name }}</span>
              <span style="color: #8492a6; font-size: 13px; margin-left: 8px;">
                ({{ group.memberCount || 0 }}人)
              </span>
            </el-option>
          </el-select>

          <el-button
            type="primary"
            size="small"
            icon="Plus"
            style="margin-left: 20px"
            :disabled="!selectedSubjectId"
            @click="handleAddAuthorization"
            class="modern-btn primary-btn"
          >
            新增授权
          </el-button>

          <el-button
            type="default"
            size="small"
            icon="Refresh"
            style="margin-left: 10px"
            :disabled="!selectedSubjectId"
            :loading="tableLoading"
            @click="handleRefresh"
            class="modern-btn secondary-btn"
          >
            刷新
          </el-button>
        </div>
      </div>

      <!-- 已授权主机列表 -->
      <div v-if="selectedSubjectId" class="host-list-section">
        <div class="section-header">
          <h3 class="section-title">已授权主机列表</h3>
          <el-button
            type="danger"
            size="small"
            :disabled="selectedHosts.length === 0"
            @click="handleBatchRevoke"
            class="modern-btn danger-btn"
          >
            <el-icon><Delete /></el-icon>
            批量取消授权 ({{ selectedHosts.length }})
          </el-button>
        </div>

        <el-table
          v-loading="tableLoading"
          :data="hostList"
          border
          stripe
          class="modern-table"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="hostName" label="主机名称" min-width="150" />
          <el-table-column prop="sshIp" label="SSH地址" min-width="150" />
          <el-table-column prop="privateIp" label="内网IP" min-width="150" />
          <el-table-column prop="groupName" label="业务组" min-width="120" />
          <el-table-column prop="status" label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag
                :type="getStatusType(row.status)"
                size="small"
              >
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button
                type="danger"
                size="small"
                @click="handleRevokeSingle(row)"
              >
                取消授权
              </el-button>
            </template>
          </el-table-column>
        </el-table>

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

    <!-- 主机选择对话框 -->
    <HostSelectorDialog
      v-model="selectorVisible"
      title="选择要授权的主机"
      @hosts-selected="handleHostsSelected"
    />
  </div>
</template>

<script>
import {
  Search,
  Plus,
  Delete,
  Refresh
} from '@element-plus/icons-vue'
import HostSelectorDialog from '@/components/HostSelectorDialog.vue'

export default {
  name: 'HostAuthorization',
  components: {
    Search,
    Plus,
    Delete,
    Refresh,
    HostSelectorDialog
  },
  data() {
    return {
      selectedSubjectId: null,
      groupOptions: [],
      tableLoading: false,
      hostList: [],
      selectedHosts: [],
      pagination: {
        page: 1,
        pageSize: 10,
        total: 0
      },
      selectorVisible: false
    }
  },
  mounted() {
    // 从路由参数中获取初始值
    const { subjectId } = this.$route.query
    if (subjectId) {
      this.selectedSubjectId = Number(subjectId)
      this.fetchHostList()
    }

    this.loadGroupOptions()
  },
  methods: {
    // 加载用户组选项
    async loadGroupOptions() {
      try {
        const { data: res } = await this.$api.getUserGroupList({
          page: 1,
          pageSize: 100
        })
        if (res.code === 200) {
          this.groupOptions = res.data.list || []
        }
      } catch (error) {
        console.error('加载用户组失败:', error)
      }
    },

    // 授权对象变化
    handleSubjectChange() {
      this.pagination.page = 1
      this.fetchHostList()
    },

    // 刷新列表
    handleRefresh() {
      this.fetchHostList()
    },

    // 获取已授权主机列表
    async fetchHostList() {
      if (!this.selectedSubjectId) return

      this.tableLoading = true
      try {
        const { data: res } = await this.$api.getAuthorizedHosts({
          subjectType: 'group',
          subjectId: this.selectedSubjectId,
          page: this.pagination.page,
          pageSize: this.pagination.pageSize
        })
        if (res.code === 200) {
          this.hostList = res.data.list || []
          this.pagination.total = res.data.total || 0
        } else {
          this.$message.error(res.message || '获取授权主机列表失败')
        }
      } catch (error) {
        this.$message.error('获取授权主机列表失败')
      } finally {
        this.tableLoading = false
      }
    },

    // 分页变化
    handleSizeChange(size) {
      this.pagination.pageSize = size
      this.fetchHostList()
    },

    handlePageChange(page) {
      this.pagination.page = page
      this.fetchHostList()
    },

    // 表格选择变化
    handleSelectionChange(selection) {
      this.selectedHosts = selection
    },

    // 新增授权
    handleAddAuthorization() {
      this.selectorVisible = true
    },

    // 主机选择完成
    async handleHostsSelected(hosts) {
      if (!hosts || hosts.length === 0) {
        this.$message.warning('请选择要授权的主机')
        return
      }

      try {
        const hostIds = hosts.map(h => h.id)
        const data = {
          hostIds,
          groupId: this.selectedSubjectId
        }

        const res = await this.$api.authorizeHostsToGroup(data)

        if (res.data.code === 200) {
          this.$message.success('授权成功')
          this.selectorVisible = false
          this.fetchHostList()
        } else {
          this.$message.error(res.data.message || '授权失败')
        }
      } catch (error) {
        this.$message.error(error.message || '授权失败')
      }
    },

    // 取消单个授权
    handleRevokeSingle(row) {
      this.$confirm('确定要取消该主机的授权吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.revokeAuthorization([row.id])
      }).catch(() => {})
    },

    // 批量取消授权
    handleBatchRevoke() {
      this.$confirm(
        `确定要取消选中的 ${this.selectedHosts.length} 台主机的授权吗？`,
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(() => {
        const hostIds = this.selectedHosts.map(h => h.id)
        this.revokeAuthorization(hostIds)
      }).catch(() => {})
    },

    // 取消授权
    async revokeAuthorization(hostIds) {
      try {
        const data = {
          hostIds,
          groupId: this.selectedSubjectId
        }

        const res = await this.$api.revokeHostsFromGroup(data)

        if (res.data.code === 200) {
          this.$message.success('取消授权成功')
          this.fetchHostList()
        } else {
          this.$message.error(res.data.message || '取消授权失败')
        }
      } catch (error) {
        this.$message.error('取消授权失败')
      }
    },

    // 获取状态文本
    getStatusText(status) {
      const map = {
        1: '认证成功',
        2: '未认证',
        3: '认证失败'
      }
      return map[status] || '未知'
    },

    // 获取状态类型
    getStatusType(status) {
      const map = {
        1: 'success',
        2: 'info',
        3: 'danger'
      }
      return map[status] || 'info'
    }
  }
}
</script>

<style lang="less" scoped>
.authorization-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.auth-card {
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

.subject-selector {
  margin-bottom: 24px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);

  .selector-content {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 16px;
  }
}

.host-list-section {
  margin-bottom: 20px;

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    padding: 16px;
    background: rgba(103, 126, 234, 0.05);
    border-radius: 12px;

    .section-title {
      color: #2c3e50;
      font-size: 16px;
      font-weight: 600;
      margin: 0;
    }
  }
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
  .authorization-management {
    padding: 12px;
  }

  .auth-card {
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

  .subject-selector {
    padding: 12px;
  }

  .selector-content {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
