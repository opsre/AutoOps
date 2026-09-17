<template>
  <div class="role-container">
    <div class="glass-card main-card">
      <!-- 卡片标题 -->
      <div class="card-header">
        <h2 class="gradient-title">角色管理</h2>
      </div>
      
      <!-- 搜索区域 -->
      <div class="search-section">
        <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="角色名称" prop="roleName">
          <el-input 
            placeholder="请输入角色名称" 
            size="small" 
            clearable 
            v-model="queryParams.roleName"
            @keyup.enter="handleQuery"
            class="search-input">
          </el-input>
        </el-form-item>
        <el-form-item label="账号状态" prop="status" style="width: 150px;" >
          <el-select size="small" placeholder="账号状态" v-model="queryParams.status" class="search-select">
            <el-option v-for="item in statusList" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="small" @click="handleQuery" class="modern-btn primary-btn">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button size="small" @click="resetQuery" class="modern-btn reset-btn">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
          <el-button 
            type="success" 
            size="small"
            @click="addRoleDialogVisible = true" 
            v-authority="['base:role:add']"
            class="modern-btn success-btn">
            <el-icon><Plus /></el-icon>
            新增角色
          </el-button>
        </el-form-item>
        </el-form>
      </div>

      <!-- 表格区域 -->
      <div class="table-section">
      <el-table 
        v-loading="Loading" 
        :data="roleList" 
        class="modern-table"
        style="width: 100%">
        <el-table-column label="ID" prop="id" v-if="false" />
        <el-table-column label="角色名称" prop="roleName"  />
        <el-table-column label="角色标识" prop="roleKey"  />
        <el-table-column label="创建时间" prop="createTime"  />
        <el-table-column label="账号状态" >
          <template v-slot="scope">
            <el-switch 
              v-model="scope.row.status" 
              :active-value="1" 
              :inactive-value="2" 
              active-color="#667eea"
              inactive-color="#f56565" 
              active-text="启用" 
              inactive-text="停用" 
              @change="roleUpdateStatus(scope.row)"
              class="status-switch">
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column label="描述" prop="description" />
        <el-table-column label="操作" width="280">
          <template v-slot="scope">
            <div class="operation-buttons">
              <el-tooltip content="编辑" placement="top">
                <el-button
                  type="warning"
                  size="small"
                  @click="showEditRoleDialog(scope.row.id)"
                  v-authority="['base:role:edit']"
                >
                  <el-icon><Edit /></el-icon>编辑
                </el-button>
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button
                  type="danger"
                  size="small"
                  @click="handleRoleDelete(scope.row)"
                  v-authority="['base:role:delete']"
                >
                  <el-icon><Delete /></el-icon>删除
                </el-button>
              </el-tooltip>
              <el-tooltip content="分配权限" placement="top">
                <el-button
                  type="primary"
                  size="small"
                  @click="showSetMenuDialog(scope.row)"
                  v-authority="['base:role:assign']"
                >
                  <el-icon><Setting /></el-icon>分配权限
                </el-button>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
        </el-table>
        
        <!-- 分页区域 -->
        <el-pagination 
          @size-change="handleSizeChange" 
          @current-change="handleCurrentChange"
          :current-page="queryParams.pageNum" 
          :page-sizes="[10, 50, 100, 500, 1000]" 
          :page-size="queryParams.pageSize"
          layout="total, sizes, prev, pager, next, jumper" 
          :total="total"
          class="modern-pagination">
        </el-pagination>
      </div>
    </div>

    <!-- 新增角色对话框 -->
    <el-dialog 
      v-model="addRoleDialogVisible" 
      title="新增角色"
      width="35%" 
      @close="addRoleDialogClosed"
      :show-close="true"
      class="clean-dialog">
      <el-form ref="addRoleFormRefForm" label-width="80px" :model="addRoleForm" :rules="addRoleFormRules" class="dialog-form">
        <el-form-item label="角色名称" prop="roleName">
          <el-input placeholder="请输入角色名称" v-model="addRoleForm.roleName" class="form-input" />
        </el-form-item>
        <el-form-item label="角色标识" prop="roleKey">
          <el-input placeholder="请输入角色标识" v-model="addRoleForm.roleKey" class="form-input" />
        </el-form-item>
        <el-form-item label="角色状态" prop="status">
          <el-radio-group v-model="addRoleForm.status" class="radio-group">
            <el-radio :label="1" class="custom-radio">正常</el-radio>
            <el-radio :label="2" class="custom-radio">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="角色描述" prop="description">
          <el-input placeholder="请输入角色描述" type="textarea" v-model="addRoleForm.description" class="form-textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="addRole" class="dialog-btn confirm-btn">确 定</el-button>
          <el-button @click="addRoleDialogVisible = false" class="dialog-btn cancel-btn">取 消</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑角色对话框 -->
    <el-dialog 
      v-model="editRoleDialogVisible" 
      title="修改角色"
      width="35%" 
      @close="editRoleDialogClosed"
      :show-close="true"
      class="clean-dialog">
      <el-form ref="editRoleFormRefForm" label-width="80px" :model="roleInfo" :rules="editRoleFormRules" class="dialog-form">
        <el-form-item label="角色名称" prop="roleName">
          <el-input placeholder="请输入角色名称" v-model="roleInfo.roleName" class="form-input" />
        </el-form-item>
        <el-form-item label="角色标识" prop="roleKey">
          <el-input placeholder="请输入角色标识" v-model="roleInfo.roleKey" class="form-input" />
        </el-form-item>
        <el-form-item label="角色状态" prop="status">
          <el-radio-group v-model="roleInfo.status" class="radio-group">
            <el-radio :label="1" class="custom-radio">正常</el-radio>
            <el-radio :label="2" class="custom-radio">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="角色描述" prop="description">
          <el-input placeholder="请输入角色描述" type="textarea" v-model="roleInfo.description" class="form-textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="editRole" class="dialog-btn confirm-btn">确 定</el-button>
          <el-button @click="editRoleDialogVisible = false" class="dialog-btn cancel-btn">取 消</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 分配权限对话框 -->
    <el-dialog 
      v-model="setMenuDialogVisible" 
      title="分配权限"
      width="25%" 
      @close="setRightDialogClosed"
      :show-close="true"
      class="clean-dialog">
      <el-tree 
        :data="menuList" 
        :props="treeProps" 
        show-checkbox 
        node-key="id"
        :default-checked-keys="defKeys" 
        ref="treeRef" 
        @check="handleTreeCheck"
        class="permission-tree">
      </el-tree>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="allotMenus" class="dialog-btn confirm-btn">确 定</el-button>
          <el-button @click="setMenuDialogVisible = false" class="dialog-btn cancel-btn">取 消</el-button>
        </span>
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
  Edit,
  Delete,
  Setting
} from '@element-plus/icons-vue'
export default {
  components: { Treeselect },
  data() {
    return {
      statusList: [{
        value: '1',
        label: '启用'
      }, {
        value: '2',
        label: '停用'
      }],
      queryParams: {},
      Loading: false,
      roleList: [],
      total: 0,
      addRoleDialogVisible: false,
      addRoleForm: {
        roleName: '',
        roleKey: '',
        description: '',
        status: 1
      },
      addRoleFormRules: {
        roleName: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
        roleKey: [{ required: true, message: '请角色权限标识', trigger: 'blur' }],
        status: [{ required: true, message: '请输入角色状态', trigger: 'blur' }],
        description: [{ required: true, message: '请输入角色描述', trigger: 'blur' }],
      },
      editRoleDialogVisible: false,
      roleInfo: {},
      editRoleFormRules: {
        roleName: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
        roleKey: [{ required: true, message: '请输入角色权限标识', trigger: 'blur' }],
        status: [{ required: true, message: '请输入角色状态', trigger: 'blur' }],
        description: [{ required: true, message: '请输入角色描述', trigger: 'blur' }],
      },
      setMenuDialogVisible: false,
      menuList: [],
      treeProps: {
        label: 'label'
      },
      defKeys: [],
      id: '',
    }
  },
  methods: {
    // 递归过滤出有权限的菜单
    filterMenuTree(menuList, permissionIds) {
      return menuList.filter(node => {
        // 如果当前节点在权限列表中，保留它
        if (permissionIds.includes(node.id)) {
          return true;
        }

        // 如果有子节点，递归处理子节点
        if (node.children && node.children.length > 0) {
          node.children = this.filterMenuTree(node.children, permissionIds)
          return node.children.length > 0
        }

        return false
      })
    },

    // 过滤出叶子节点权限，避免父权限导致自动全选
    filterLeafPermissions(menuTree, rolePermissions) {
      const leafPermissions = [];

      // 递归收集所有叶子节点ID
      const collectLeafNodes = (nodes) => {
        nodes.forEach(node => {
          if (!node.children || node.children.length === 0) {
            // 叶子节点
            leafPermissions.push(node.id);
          } else {
            // 有子节点，继续递归
            collectLeafNodes(node.children);
          }
        });
      };

      collectLeafNodes(menuTree);

      // 只返回在角色权限列表中且是叶子节点的权限
      return rolePermissions.filter(id => leafPermissions.includes(id));
    },

    // 基于真实sort字段对菜单数据排序
    sortMenusByOrder(menuList, menuWithSort) {
      if (!menuList || !Array.isArray(menuList)) {
        return [];
      }

      // 创建id到sort值的映射
      const sortMap = {};
      if (menuWithSort && Array.isArray(menuWithSort)) {
        menuWithSort.forEach(menu => {
          sortMap[menu.id] = menu.sort || 0;
        });
      }

      // 对当前层级按真实sort字段排序
      const sortedList = [...menuList].sort((a, b) => {
        const sortA = sortMap[a.id] || 999;
        const sortB = sortMap[b.id] || 999;
        return sortA - sortB;
      });

      // 递归处理子菜单
      return sortedList.map(menu => {
        if (menu.children && menu.children.length > 0) {
          return {
            ...menu,
            children: this.sortMenusByOrder(menu.children, menuWithSort)
          };
        }
        return menu;
      });
    },
    // 列表
    async getRoleList() {
      this.Loading = true
      const { data: res } = await this.$api.queryRoleList(this.queryParams)
      if (res.code !== 200) {
        this.$message.error(res.message);
      } else {
        this.roleList = res.data.list
        this.total = res.data.total
        this.Loading = false
      }
    },
    // 搜索
    handleQuery() {
      this.getRoleList()
    },
    // 重置
    resetQuery() {
      this.queryParams = {}
      this.getRoleList()
      this.$message.success("重置成功")
    },
    // pageSize
    handleSizeChange(newSize) {
      this.queryParams.pageSize = newSize
      this.getRoleList()
    },
    // pageNum
    handleCurrentChange(newPage) {
      this.queryParams.pageNum = newPage
      this.getRoleList()
    },
    // 启用/停用
    async roleUpdateStatus(row) {
      let text = row.status === 2 ? "停用" : "启用"
      const confirmResult = await this.$confirm('确认要"' + text + '""' + row.roleName + '"角色吗?', "警告", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        await this.getRoleList()
        return this.$message.info('已取消删除')
      }
      await this.$api.updateRoleStatus(row.id, row.status)
      return this.$message.success(text + "成功")
      // eslint-disable-next-line no-unreachable
      await this.getRoleList()
    },
    // 监听添加角色对话框关闭
    addRoleDialogClosed() {
      this.$refs.addRoleFormRefForm.resetFields()
    },
    // 新增
    addRole() {
      this.$refs.addRoleFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.addRole(this.addRoleForm);
        if (res.code !== 200) {
          this.$message.error(res.message);
        } else {
          this.$message.success("新增角色成功")
          this.addRoleDialogVisible = false
          await this.getRoleList()
        }
      })
    },
    // 监听修改角色对话框的关闭事件
    editRoleDialogClosed() {
      this.$refs.editRoleFormRefForm.resetFields()
    },
    // 展示修改对话框
    async showEditRoleDialog(id) {
      const { data: res } = await this.$api.roleInfo(id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.roleInfo = res.data
        this.editRoleDialogVisible = true
      }
    },
    // 修改角色信息并提交
    editRole() {
      this.$refs.editRoleFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.roleUpdate(this.roleInfo)
        if (res.code !== 200) {
          this.$message.error(res.message)
        } else {
          this.editRoleDialogVisible = false
          await this.getRoleList()
          this.$message.success("修改角色成功")
        }
      })
    },
    // 删除角色
    async handleRoleDelete(row) {
      const confirmResult = await this.$confirm('是否确认删除角色为"' + row.roleName + '"的数据项？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消删除')
      }
      const { data: res } = await this.$api.deleteRole(row.id)
      console.log(row.id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.getRoleList()
      }
    },
// 展示分配菜单权限对话框
    async showSetMenuDialog(role) {
      this.id = role.id;
      this.setMenuDialogVisible = true;

      try {
        const [menuRes, roleMenuRes, menuWithSortRes] = await Promise.all([
          this.$api.querySysMenuVoList(),
          this.$api.QueryRoleMenuIdList(role.id),
          this.$api.queryMenuList({}) // 获取包含sort字段的菜单列表
        ]);

        const allMenus = menuRes.data.data; // 完整菜单树（用于权限分配）
        const rolePermissions = roleMenuRes.data.data; // 已分配权限ID列表
        const menuWithSort = menuWithSortRes.data.data; // 包含sort字段的菜单列表

        // 调试：查看菜单数据和sort字段
        console.log('权限菜单数据:', allMenus);
        console.log('带sort菜单数据:', menuWithSort);
        console.log('sort字段示例:', menuWithSort.slice(0, 5).map(item => ({ id: item.id, menuName: item.menuName, sort: item.sort })));

        // 对菜单数据按sort字段排序（使用真实的sort信息）
        const sortedMenus = this.sortMenusByOrder(allMenus, menuWithSort);
        console.log('排序后菜单数据:', sortedMenus);

        // 构建完整菜单树用于显示
        this.menuList = this.$handleTree.handleTree(sortedMenus, "id");
        console.log('构建后的菜单树:', this.menuList);

        // 过滤出叶子节点权限，避免父权限导致全选
        this.defKeys = this.filterLeafPermissions(this.menuList, rolePermissions);

        console.log('角色原始权限:', rolePermissions);
        console.log('过滤后的叶子权限:', this.defKeys);

      } catch (error) {
        console.error('加载权限失败:', error);
      }
    }
    ,
    // 监听对话框关闭事件
    setRightDialogClosed() {
      this.defKeys = []
    },

    // 处理树节点选中事件
    handleTreeCheck(checkedNode, checkedInfo) {
      // 移除自动选择父子节点的逻辑，让用户精确控制权限分配
      // 不做任何额外处理，让用户手动选择需要的权限
    },

    // 分配菜单权限
    async allotMenus() {
      // 获取完全选中的节点
      const checkedKeys = this.$refs.treeRef.getCheckedKeys()
      // 获取半选中的父节点（部分子节点被选中时，父节点会处于半选中状态）
      const halfCheckedKeys = this.$refs.treeRef.getHalfCheckedKeys()

      // 合并完全选中和半选中的权限ID，确保父权限也被保存
      const allPermissionIds = [...checkedKeys, ...halfCheckedKeys]

      console.log('完全选中的权限:', checkedKeys)
      console.log('半选中的父权限:', halfCheckedKeys)
      console.log('最终保存的权限ID:', allPermissionIds)

      const { data: res } = await this.$api.AssignPermissions(this.id, allPermissionIds)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('分配权限成功')
        await this.getRoleList()
        this.setMenuDialogVisible = false
      }
    }
  },
  created() {
    this.getRoleList()
  }
}
</script>



<style scoped>
.role-container {
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
  
  .table-section {
    background: transparent;
    border-radius: 12px;
    padding: 20px;
    overflow: hidden;
    min-width: 0;
  }
}

/* 输入框样式 */
:deep(.search-input .el-input__wrapper),
:deep(.search-select .el-input__wrapper),
:deep(.search-date .el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

:deep(.search-input .el-input__wrapper):hover,
:deep(.search-select .el-input__wrapper):hover,
:deep(.search-date .el-input__wrapper):hover {
  border-color: #c0c4cc;
}

:deep(.search-input .el-input__wrapper.is-focus),
:deep(.search-select .el-input__wrapper.is-focus),
:deep(.search-date .el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

:deep(.search-input .el-input__inner),
:deep(.search-select .el-input__inner),
:deep(.search-date .el-input__inner) {
  background: transparent;
  border: none;
  color: #2c3e50;
}

:deep(.search-input .el-input__inner::placeholder),
:deep(.search-select .el-input__inner::placeholder),
:deep(.search-date .el-input__inner::placeholder) {
  color: rgba(44, 62, 80, 0.6);
}

:deep(.search-date .el-input__suffix-inner),
:deep(.search-select .el-input__suffix-inner) {
  color: #606266;
}

.search-date {
  width: 190px;
}

.search-select {
  width: 120px;
}

/* 现代化按钮 */
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

/* 操作区域 */
.action-section {
  margin-bottom: 25px;
  display: flex;
  gap: 15px;
}

/* 表格区域 */
.table-card {
  padding: 24px;
}

/* 现代化表格样式 */
:deep(.modern-table) {
  background: transparent;
  border-radius: 12px;
  overflow: hidden;
}

:deep(.modern-table .el-table__header) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}

:deep(.modern-table .el-table__header th) {
  background: transparent;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border: none;
  padding: 8px 12px !important;
  height: 40px;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

:deep(.modern-table .el-table__row) {
  transition: all 0.3s ease;
}

:deep(.modern-table .el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

:deep(.modern-table .el-table__row td) {
  border: none;
  color: #2c3e50;
  padding: 8px 12px !important;
  height: 40px;
}

:deep(.modern-table .el-table__cell) {
  padding: 8px 12px !important;
}

:deep(.modern-table .el-table__body tr:not(:last-child) td) {
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
}

/* 操作按钮 */
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

/* 状态开关 */
:deep(.status-switch) {
  transform: scale(0.9);
}

/* 分页区域 */
.pagination-card {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

:deep(.modern-pagination) {
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 16px 24px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

:deep(.modern-pagination .el-pagination__total),
:deep(.modern-pagination .el-pagination__jump) {
  color: #ffffff;
}

:deep(.modern-pagination .btn-prev),
:deep(.modern-pagination .btn-next),
:deep(.modern-pagination .el-pager li) {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #ffffff;
  border-radius: 6px;
  margin: 0 4px;
  transition: all 0.3s ease;
}

:deep(.modern-pagination .btn-prev:hover),
:deep(.modern-pagination .btn-next:hover),
:deep(.modern-pagination .el-pager li:hover) {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-1px);
}

:deep(.modern-pagination .el-pager li.is-active) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-color: #667eea;
}

/* 对话框样式 */
:deep(.modern-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.modern-dialog .el-dialog) {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(240, 248, 255, 0.95) 100%);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.2);
}

:deep(.modern-dialog .el-dialog__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
  margin: 0;
}

:deep(.modern-dialog .el-dialog__title) {
  color: white;
  font-weight: 600;
  font-size: 18px;
}

:deep(.modern-dialog .el-dialog__headerbtn .el-dialog__close) {
  color: white;
  font-size: 18px;
}

:deep(.modern-dialog .el-dialog__body) {
  padding: 24px;
}

/* 简洁对话框样式 - 去掉标题和背景 */
:deep(.clean-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.clean-dialog .el-dialog) {
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(0, 0, 0, 0.1);
}

:deep(.clean-dialog .el-dialog__header) {
  padding: 16px 20px;
  margin: 0;
  border-bottom: 1px solid #ebeef5;
}

:deep(.clean-dialog .el-dialog__title) {
  color: #303133;
  font-weight: 600;
  font-size: 16px;
}

:deep(.clean-dialog .el-dialog__headerbtn .el-dialog__close) {
  color: #909399;
  font-size: 16px;
}

:deep(.clean-dialog .el-dialog__body) {
  padding: 24px;
}

:deep(.clean-dialog .el-dialog__headerbtn) {
  position: absolute;
  top: 16px;
  right: 20px;
  z-index: 10;
}

:deep(.clean-dialog .el-dialog__close:hover) {
  color: #303133;
}

/* 表单样式 */
.dialog-form .el-form-item__label {
  color: #333;
  font-weight: 500;
}

:deep(.form-input .el-input__wrapper),
:deep(.form-textarea .el-textarea__inner) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.form-input .el-input__wrapper):hover,
:deep(.form-textarea .el-textarea__inner):hover {
  border-color: #667eea;
  background: rgba(255, 255, 255, 0.9);
}

:deep(.form-input .el-input__wrapper.is-focus),
:deep(.form-textarea .el-textarea__inner:focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
  background: white;
}

/* 单选按钮组 */
.radio-group {
  display: flex;
  gap: 16px;
}

:deep(.custom-radio.el-radio) {
  margin-right: 0;
}

:deep(.custom-radio .el-radio__input.is-checked + .el-radio__label) {
  color: #667eea;
}

:deep(.custom-radio .el-radio__input.is-checked .el-radio__inner) {
  background-color: #667eea;
  border-color: #667eea;
}

:deep(.custom-radio .el-radio__inner:hover) {
  border-color: #667eea;
}

/* 对话框底部按钮 */
.dialog-footer {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.dialog-btn {
  padding: 10px 24px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
  border: none;
}

.confirm-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.confirm-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.cancel-btn {
  background: rgba(108, 117, 125, 0.1);
  color: #6c757d;
  border: 1px solid rgba(108, 117, 125, 0.2);
}

.cancel-btn:hover {
  background: rgba(108, 117, 125, 0.2);
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 权限树样式 */
.permission-tree {
  background: transparent;
  border-radius: 8px;
  padding: 4px;
}

:deep(.permission-tree .el-tree-node) {
  margin-bottom: 0;
}

:deep(.permission-tree .el-tree-node__content) {
  border-radius: 3px;
  transition: all 0.3s ease;
  padding: 2px 6px;
  height: 24px;
  line-height: 20px;
}

:deep(.permission-tree .el-tree-node__content:hover) {
  background: rgba(102, 126, 234, 0.1);
}

:deep(.permission-tree .el-tree-node__expand-icon) {
  padding: 1px;
  margin-right: 4px;
}

:deep(.permission-tree .el-checkbox) {
  margin-right: 4px;
}

:deep(.permission-tree .el-checkbox__input) {
  line-height: 1;
}

:deep(.permission-tree .el-tree-node__label) {
  font-size: 13px;
}

:deep(.permission-tree .el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: #667eea;
  border-color: #667eea;
}

:deep(.permission-tree .el-checkbox__inner:hover) {
  border-color: #667eea;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .role-container {
    padding: 12px;
  }
  
  .glass-card {
    margin-bottom: 16px;
    border-radius: 12px;
  }
  
  .search-card {
    padding: 16px;
  }
  
  .table-card {
    padding: 16px;
    overflow-x: auto;
  }
  
  .search-form {
    flex-direction: column;
  }
  
  .search-form .el-form-item {
    width: 100%;
    margin-bottom: 12px;
  }
  
  .search-date,
  .search-select {
    width: 100%;
  }
  
  .table-actions {
    flex-direction: column;
    gap: 4px;
  }
  
  :deep(.modern-dialog .el-dialog) {
    margin: 20px;
    width: calc(100vw - 40px) !important;
  }
}

/* 动画效果 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.glass-card {
  animation: fadeInUp 0.6s ease-out;
}

.glass-card:nth-child(1) { animation-delay: 0.1s; }
.glass-card:nth-child(2) { animation-delay: 0.2s; }
.glass-card:nth-child(3) { animation-delay: 0.3s; }

/* 加载动画优化 */
:deep(.el-loading-mask) {
  background: rgba(102, 126, 234, 0.1);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

:deep(.el-loading-spinner) {
  color: #667eea;
}
</style>
