<template>
  <div class="loginlog-management">
    <el-card shadow="hover" class="loginlog-card">
      <template #header>
        <div class="card-header">
          <span class="title">登录日志管理</span>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :model="queryParams" :inline="true" class="search-form">
          <el-form-item prop="username" label="用户名称">
            <el-input v-model="queryParams.username" placeholder="请输入用户名称" clearable size="small"
                      style="width: 200px" @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item prop="loginStatus" label="登录状态">
            <el-select size="small" placeholder="请选择状态" v-model="queryParams.loginStatus" style="width: 150px" clearable>
              <el-option v-for="item in loginStatusList" :key="item.value" :label="item.label" :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item prop="beginTime" label="开始时间">
            <el-date-picker v-model="queryParams.beginTime" size="small" type="date"
                            style="width: 190px" value-format="yyyy-MM-dd" clearable placeholder="请选择开始时间"
                            @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item prop="endTime" label="结束时间">
            <el-date-picker v-model="queryParams.endTime" size="small" type="date"
                            style="width: 190px" value-format="yyyy-MM-dd" clearable placeholder="请选择结束时间"
                            @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="Search" size="small" @click="handleQuery">搜索</el-button>
            <el-button type="warning" icon="Refresh" size="small" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 操作按钮 -->
      <div class="action-section">
        <el-button type="danger" plain icon="Delete" size="small" :disabled="multiple"
                   @click="batchHandleDelete" v-authority="['monitor:loginLog:delete']">批量删除
        </el-button>
        <el-button type="danger" plain icon="Delete" size="small" @click="handleClean" 
                   v-authority="['monitor:loginLog:clean']">清空日志</el-button>
      </div>

      <!-- 列表区域 -->
      <div class="table-section">
        <el-table v-loading="Loading" :data="sysLoginInfoList" stripe style="width: 100%"
                  class="loginlog-table" @selection-change="handleSelectionChange">
          <el-table-column type="selection" />
          <el-table-column label="ID" prop="id" v-if="false" />
          <el-table-column label="用户账号" prop="username" min-width="120" />
          <el-table-column label="登录IP" prop="ipAddress" min-width="140" />
          <el-table-column label="登录地点" prop="loginLocation" min-width="150" show-overflow-tooltip />
          <el-table-column label="浏览器" prop="browser" min-width="120" show-overflow-tooltip />
          <el-table-column label="操作系统" prop="os" min-width="120" show-overflow-tooltip />
          <el-table-column label="登录状态" prop="loginStatus" min-width="100">
            <template #default="scope">
              <el-tag v-if="scope.row.loginStatus === 1" type="success" size="small" effect="dark">成功</el-tag>
              <el-tag v-else-if="scope.row.loginStatus === 2" type="danger" size="small" effect="dark">失败</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="提示消息" prop="message" min-width="150" show-overflow-tooltip />
          <el-table-column label="访问时间" prop="loginTime" min-width="180" />
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="scope">
              <div class="operation-buttons">
                <el-tooltip content="删除" placement="top">
                  <el-button
                    type="danger"
                    icon="Delete"
                    size="small"
                    @click="handleDelete(scope.row.id)"
                    v-authority="['monitor:loginLog:delete']"
                  >删除</el-button>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页区域 -->
      <div class="pagination-section">
        <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange"
                       :current-page="queryParams.pageNum" :page-sizes="[10, 50, 100, 500, 1000]" :page-size="queryParams.pageSize"
                       layout="total, sizes, prev, pager, next, jumper" :total="total" />
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  data() {
    return {
      queryParams: {},
      loginStatusList: [{
        value: '1',
        label: '成功'
      }, {
        value: '2',
        label: '失败'
      }],
      sysLoginInfoList: [],
      Loading: true,
      ids: [],
      single: true,
      multiple: true,
      total: 0,
    }
  },
  methods: {
    // 多选框选中数据
    handleSelectionChange(selection) {
      this.ids = selection.map(item => item.id);
      this.single = selection.length != 1;
      this.multiple = !selection.length;
    },
    // 查询列表
    async getSysLoginInfoList() {
      this.Loading = true
      const { data: res } = await this.$api.querySysLoginInfoList(this.queryParams)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.sysLoginInfoList = res.data.list
        this.total = res.data.total
        this.Loading = false
      }
    },
    // 搜索按钮操作
    handleQuery() {
      this.getSysLoginInfoList();
    },
    // 重置按钮操作
    resetQuery() {
      this.queryParams = {}
      this.getSysLoginInfoList();
      this.$message.success("重置成功")
    },
    // pageSize
    handleSizeChange(newSize) {
      this.queryParams.pageSize = newSize
      this.getSysLoginInfoList()
    },
    // pageNum
    handleCurrentChange(newPage) {
      this.queryParams.pageNum = newPage
      this.getSysLoginInfoList()
    },
    // 清空
    async handleClean() {
      const confirmResult = await this.$confirm('是否清空登录日志？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消')
      }
      const { data: res } = await this.$api.cleanSysLoginInfo()
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('清空成功')
        await this.getSysLoginInfoList()
      }
    },
    // 删除
    async handleDelete(id) {
      const confirmResult = await this.$confirm('是否确认删除登录日志编号为"' + id + '"的数据项？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消删除')
      }
      const { data: res } = await this.$api.deleteSysLoginInfo(id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.getSysLoginInfoList()
      }
    },
    // 批量删除
    async batchHandleDelete() {
      const loginInfoIds = this.ids;
      const confirmResult = await this.$confirm('是否确认删除登录日志编号为"' + loginInfoIds + '"的数据项？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消删除')
      }
      const { data: res } = await this.$api.batchDeleteSysLoginInfo(loginInfoIds)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.getSysLoginInfoList()
      }
    },
  },
  created() {
    this.getSysLoginInfoList()
  }
}
</script>

<style scoped>
.loginlog-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.loginlog-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
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
  margin-bottom: 24px;
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
  margin-bottom: 24px;
  padding: 12px 0;
}

.action-section .el-button {
  margin-right: 12px;
}

.table-section {
  margin-bottom: 24px;
}

.loginlog-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.loginlog-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.loginlog-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.loginlog-table :deep(.el-table__header th .cell) {
  color: #2c3e50 !important;
  font-weight: 700 !important;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.loginlog-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.loginlog-table :deep(.el-table__row:hover) {
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

.pagination-section {
  display: flex;
  justify-content: center;
  padding: 20px 0;
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

.el-date-picker :deep(.el-input__inner) {
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.2);
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__wrapper):hover,
.el-select :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.el-input :deep(.el-input__wrapper.is-focus),
.el-select :deep(.el-input__wrapper.is-focus),
.el-date-picker :deep(.el-input__inner):focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
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
