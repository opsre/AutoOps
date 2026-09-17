<template>
  <div class="modern-admin-container">
    <!-- 主内容区域 -->
    <div class="glass-card main-card">
      <template v-if="false"></template>
      <div class="card-header">
        <span class="main-title">用户管理系统</span>
      </div>
      <!-- 左右布局容器 -->
      <div class="user-management-container">
        <!-- 左侧部门树 -->
        <div class="dept-tree-section">
          <div class="dept-card">
            <!-- 科技感标题区域 -->
            <div class="card-header">
              <div class="title-wrapper">
                <div class="title-icon">
                  <div class="pulse-ring"></div>
                  <el-icon class="main-icon"><OfficeBuilding /></el-icon>
                </div>
                <div class="title-content">
                  <h3 class="title">组织架构</h3>
                  <span class="subtitle">Department</span>
                </div>
              </div>
              <div class="stats-indicator">
                <span class="stats-text">{{ deptList.length }}</span>
                <span class="stats-label">Depts</span>
                <div class="action-buttons-wrapper">
                  <div class="toggle-btn" @click="toggleAll" :title="isExpanded ? '折叠全部' : '展开全部'">
                    <img
                      :src="require('@/assets/image/折叠.svg')"
                      class="toggle-icon"
                      :class="{ 'expanded': isExpanded }"
                      alt="折叠展开"
                    />
                  </div>
                  <div class="refresh-btn" @click="handleRefresh" title="刷新部门">
                    <el-icon class="refresh-icon"><Refresh /></el-icon>
                  </div>
                </div>
              </div>
            </div>

            <!-- 高级搜索区域 -->
            <div class="search-container">
              <div class="search-wrapper">
                <div class="search-input-wrapper">
                  <el-input
                    v-model="deptSearchText"
                    placeholder="搜索部门..."
                    clearable
                    size="small"
                    @input="handleDeptSearch"
                    class="tech-input"
                  >
                    <template #prefix>
                      <el-icon class="search-icon"><Search /></el-icon>
                    </template>
                  </el-input>
                  <div class="search-glow"></div>
                </div>
              </div>
            </div>

            <!-- 树形结构区域 -->
            <div class="tree-container">
              <el-tree
                ref="deptTree"
                :data="deptList"
                :props="defaultProps"
                node-key="id"
                :expanded-keys="expandedKeys"
                :highlight-current="true"
                :default-expanded-keys="expandedKeys"
                @node-click="handleDeptClick"
                @node-expand="handleNodeExpand"
                @node-collapse="handleNodeCollapse"
                class="tech-tree"
              >
                <template v-slot="{ node, data }">
                  <div
                    class="tree-node"
                    :class="{ 'parent-node': !data.parentId }"
                    @contextmenu.prevent="showContextMenu($event, node, data)"
                  >
                    <div class="node-content">
                      <div class="node-icon-wrapper">
                        <!-- 公司（deptType=1）：根据展开状态切换文件夹图标 -->
                        <template v-if="data.deptType === 1">
                          <div class="parent-icon" :class="{ 'expanded': expandedKeys.includes(node.key) }">
                            <div class="icon-bg"></div>
                            <img
                              v-if="expandedKeys.includes(node.key)"
                              :src="require('@/assets/image/打开文件夹.svg')"
                              class="parent-icon-img"
                              alt="打开文件夹"
                            />
                            <img
                              v-else
                              :src="require('@/assets/image/关闭文件夹.svg')"
                              class="parent-icon-img"
                              alt="关闭文件夹"
                            />
                          </div>
                        </template>
                        <!-- 中心部门（deptType=2）：组织图标 -->
                        <template v-else-if="data.deptType === 2">
                          <div class="center-icon">
                            <div class="dot-indicator"></div>
                            <img
                              :src="require('@/assets/image/组织.svg')"
                              class="center-icon-img"
                              alt="中心部门"
                            />
                          </div>
                        </template>
                        <!-- 子部门（deptType=3）：人员分组图标 -->
                        <template v-else>
                          <div class="child-icon">
                            <div class="dot-indicator"></div>
                            <img
                              :src="require('@/assets/image/人员分组.svg')"
                              class="child-icon-img"
                              alt="子部门"
                            />
                          </div>
                        </template>
                      </div>
                      <div class="node-text">
                        <span class="node-label">
                          {{ node.label }}
                        </span>
                        <div class="connection-line" v-if="data.deptType === 1"></div>
                      </div>
                    </div>
                    <div class="hover-effect"></div>
                  </div>
                </template>
              </el-tree>
            </div>
          </div>
        </div>

        <!-- 右键菜单 -->
        <div
          v-if="contextMenuVisible"
          class="context-menu"
          :style="{ left: contextMenuPosition.x + 'px', top: contextMenuPosition.y + 'px' }"
          @click.stop
        >
          <div class="context-menu-item" @click="createDept">
            <el-icon><Plus /></el-icon>
            <span>创建部门</span>
          </div>
          <div
            class="context-menu-item"
            @click="renameDept"
            v-if="contextMenuData.data && contextMenuData.data.id !== 1"
          >
            <el-icon><EditPen /></el-icon>
            <span>重命名</span>
          </div>
          <div
            class="context-menu-item"
            @click="editDept"
            v-if="contextMenuData.data && contextMenuData.data.id !== 1"
          >
            <el-icon><Edit /></el-icon>
            <span>修改部门</span>
          </div>
          <div
            class="context-menu-item danger"
            @click="deleteDept"
            v-if="contextMenuData.data && contextMenuData.data.id !== 1"
          >
            <el-icon><Delete /></el-icon>
            <span>删除部门</span>
          </div>
          <!-- 一级部门的重命名按钮 -->
          <div
            class="context-menu-item"
            @click="renameDept"
            v-if="contextMenuData.data && contextMenuData.data.id === 1"
          >
            <el-icon><EditPen /></el-icon>
            <span>重命名</span>
          </div>
        </div>

        <!-- 创建/编辑部门对话框 -->
        <el-dialog
          :title="getDialogTitle()"
          v-model="deptDialogVisible"
          width="500px"
          @close="handleDialogClose"
          class="modern-dialog"
        >
          <el-form
            :model="deptForm"
            :rules="deptRules"
            ref="deptFormRef"
            label-width="100px"
            class="modern-dialog-form"
          >
            <el-form-item label="部门类型" prop="deptType" v-if="dialogMode === 'create' || dialogMode === 'edit'">
              <el-radio-group v-model="deptForm.deptType" class="modern-radio-group" @change="handleDeptTypeChange">
                <el-radio :label="2" class="modern-radio">中心部门</el-radio>
                <el-radio :label="3" class="modern-radio">子部门</el-radio>
              </el-radio-group>
              <div style="margin-top: 8px; color: #909399; font-size: 12px;">
                <div v-if="deptForm.deptType === 2">中心部门将自动创建在公司下</div>
                <div v-else-if="deptForm.deptType === 3 && contextMenuData.data">
                  子部门将创建在 <span style="color: #667eea; font-weight: 500;">{{ contextMenuData.data.deptName }}</span> 下
                </div>
              </div>
            </el-form-item>
            <el-form-item label="部门名称" prop="deptName">
              <el-input
                v-model="deptForm.deptName"
                placeholder="请输入部门名称"
                class="modern-input"
              />
            </el-form-item>
            <el-form-item
              label="父级部门"
              prop="parentId"
              v-if="false"
            >
              <treeselect
                v-model="deptForm.parentId"
                :options="deptSelectOptions"
                :show-count="true"
                placeholder="请选择父级部门"
                class="modern-treeselect"
              />
            </el-form-item>
            <el-form-item label="排序" prop="sort">
              <el-input-number
                v-model="deptForm.sort"
                :min="0"
                :max="9999"
                controls-position="right"
                class="modern-input"
              />
            </el-form-item>
            <el-form-item label="部门状态" prop="status">
              <el-radio-group v-model="deptForm.status" class="modern-radio-group">
                <el-radio :label="1" class="modern-radio">启用</el-radio>
                <el-radio :label="0" class="modern-radio">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-form>
          <template #footer>
            <el-button @click="deptDialogVisible = false" class="modern-btn secondary-btn">取消</el-button>
            <el-button type="primary" @click="submitDeptForm" class="modern-btn primary-btn">确定</el-button>
          </template>
        </el-dialog>

        <!-- 右侧用户管理区域 -->
        <div class="user-table-section">
          <!-- 搜索表单 -->
          <div class="glass-card search-card">
            <el-form :inline="true" :model="queryParams" class="modern-form">
              <el-form-item label="用户账号" prop="username">
                <el-input
                    size="small"
                    placeholder="请输入用户账号"
                    clearable
                    v-model="queryParams.username"
                    @keyup.enter="handleQuery"
                    class="modern-input"
                />
              </el-form-item>
              <el-form-item label="账号状态" prop="status"  style="width: 150px;" >
                <el-select size="small" placeholder="账号状态" v-model="queryParams.status" class="modern-select">
                  <el-option v-for="item in statusList" :key="item.value" :label="item.label" :value="item.value"></el-option>
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
                <el-button type="success" size="small" @click="addDialogVisible = true" v-authority="['base:admin:add']" class="modern-btn success-btn">
                  <el-icon><Plus /></el-icon>
                  新增
                </el-button>
              </el-form-item>
            </el-form>
          </div>

          <!-- 用户表格 -->
          <div class="glass-card table-card">
            <div class="card-header">
              <h3 class="section-title">用户列表</h3>
            </div>
            <el-table
                v-loading="Loading"
                :data="adminList"
                class="modern-table"
                :row-class-name="tableRowClassName"
            >
              <el-table-column label="ID" prop="id" v-if="false" />

              <el-table-column label="用户账号" prop="username"  />
              <el-table-column label="用户头像" prop="icon" >
                <template v-slot="scope">
                  <el-avatar shape="circle" :src="scope.row.icon" class="modern-avatar"></el-avatar>
                </template>
              </el-table-column>
              <el-table-column label="用户昵称" prop="nickname" />
              <el-table-column label="角色名称" prop="roleName" >
                <template v-slot="scope">
                  <el-tag type="success" class="modern-tag">{{ scope.row.roleName }}</el-tag>
                </template>
              </el-table-column>

              <el-table-column label="岗位" >
                <template v-slot="scope">
                  <div class="modern-text">{{ scope.row.postName }}</div>
                </template>
              </el-table-column>
              <el-table-column label="账号状态" >
                <template v-slot="scope">
                  <el-switch
                      v-model="scope.row.status"
                      :active-value="1"
                      :inactive-value="2"
                      active-color="#67C23A"
                      inactive-color="#F56C6C"
                      active-text="启用"
                      inactive-text="停用"
                      @change="adminUpdateStatus(scope.row)"
                      class="modern-switch"
                  ></el-switch>
                </template>
              </el-table-column>

              <el-table-column label="操作" fixed="right" width="280">
                <template v-slot="scope">
                  <div class="operation-buttons">
                    <el-tooltip content="编辑" placement="top">
                      <el-button
                        type="warning"
                        size="small"
                        @click="showEditAdminDialog(scope.row.id)"
                        v-authority="['base:admin:edit']"
                      >
                        <el-icon><Edit /></el-icon>编辑
                      </el-button>
                    </el-tooltip>
                    <el-tooltip content="删除" placement="top">
                      <el-button
                        type="danger"
                        size="small"
                        @click="handleAdminDelete(scope.row)"
                        v-authority="['base:admin:delete']"
                      >
                        <el-icon><Delete /></el-icon>删除
                      </el-button>
                    </el-tooltip>
                    <el-tooltip content="重置密码" placement="top">
                      <el-button
                        type="info"
                        size="small"
                        @click="handleResetPwd(scope.row)"
                        v-authority="['base:admin:reset']"
                      >
                        <el-icon><Key /></el-icon>重置密码
                      </el-button>
                    </el-tooltip>
                  </div>
                </template>
              </el-table-column>
            </el-table>

            <!-- 分页 -->
            <div class="pagination-wrapper">
              <el-pagination
                  @size-change="handleSizeChange"
                  @current-change="handleCurrentChange"
                  :current-page="queryParams.pageNum"
                  :page-sizes="[10, 50, 100, 500, 1000]"
                  :page-size="queryParams.pageSize"
                  layout="total, sizes, prev, pager, next, jumper"
                  :total="total"
                  class="modern-pagination"
              ></el-pagination>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 新增用户对话框 -->
    <el-dialog title="新增用户" v-model="addDialogVisible" width="40%" @close="addDialogClosed" class="modern-dialog">
      <el-form :model="addForm" :rules="addFormRules" ref="addFormRefForm" label-width="80px" class="modern-dialog-form">
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户名称" prop="username">
              <el-input v-model="addForm.username" placeholder="请输入英文名称" maxlength="30" class="modern-input" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户密码" prop="password">
              <el-input v-model="addForm.password" placeholder="请输入用户密码" type="password" maxlength="20" show-password class="modern-input" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户昵称" prop="nickname">
              <el-input v-model="addForm.nickname" placeholder="请输入中文昵称" maxlength="30" class="modern-input" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item size="small" label="归属部门" prop="deptId">
              <treeselect v-model="addForm.deptId" :options="deptSelectOptions" :show-count="true"
                          placeholder="请选择归属部门" class="modern-treeselect" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="手机号码" prop="phone">
              <el-input v-model="addForm.phone" placeholder="请输入手机号码" maxlength="11" class="modern-input" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户邮箱" prop="email">
              <el-input v-model="addForm.email" placeholder="请输入邮箱" maxlength="50" class="modern-input" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户状态" prop="status">
              <el-radio-group v-model="addForm.status" class="modern-radio-group">
                <el-radio :label="1" class="modern-radio">正常</el-radio>
                <el-radio :label="2" class="modern-radio">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户岗位" prop="postId">
              <el-select placeholder="请选择岗位" v-model="addForm.postId" class="modern-select">
                <el-option v-for="item in postList" :key="item.id" :label="item.postName" :value="item.id"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户角色" prop="roleId">
              <el-select placeholder="请选择角色" v-model="addForm.roleId" class="modern-select">
                <el-option v-for="item in roleList" :key="item.id" :label="item.roleName" :value="item.id"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <el-form-item label="个人简介" prop="note">
              <el-input v-model="addForm.note" type="textarea" placeholder="请输入内容" class="modern-textarea"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <span slot="footer" class="dialog-footer">
          <el-button type="primary" @click="addAdmin" class="modern-btn primary-btn">确 定</el-button>
          <el-button type="info" @click="addDialogVisible = false" class="modern-btn secondary-btn">取 消</el-button>
        </span>
      </el-form>
    </el-dialog>

    <!-- 编辑用户对话框 -->
    <el-dialog title="修改用户" v-model="editDialogVisible" width="40%" @close="editDialogClosed" class="modern-dialog">
      <el-form :model="adminInfo" :rules="editFormRules" ref="editFormRefForm" label-width="80px" class="modern-dialog-form">
        <el-row>
          <el-col :span="24">
            <el-form-item label="用户头像" prop="icon">
              <el-upload
                :headers="uploadHeaders"
                class="avatar-uploader"
                :action="uploadIconUrl"
                :show-file-list="false"
                :on-success="handleEditAvatarSuccess"
              >
                <img v-if="editIcon" :src="editIcon" :key="editAvatarKey" class="avatar" title="点击更换头像">
                <img v-else :src="adminInfo.icon" :key="editAvatarKey" class="avatar" title="点击更换头像">
              </el-upload>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户昵称" prop="nickname">
              <el-input v-model="adminInfo.nickname" placeholder="请输入用户昵称" maxlength="30" class="modern-input" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item size="small" label="归属部门" prop="deptId">
              <treeselect v-model="adminInfo.deptId" :options="deptSelectOptions" :show-count="true"
                          placeholder="请选择归属部门" class="modern-treeselect" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="手机号码" prop="phone">
              <el-input v-model="adminInfo.phone" placeholder="请输入手机号码" maxlength="11" class="modern-input" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户邮箱" prop="email">
              <el-input v-model="adminInfo.email" placeholder="请输入邮箱" maxlength="50" class="modern-input" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户名称" prop="username">
              <el-input v-model="adminInfo.username" placeholder="请输入用户名称" maxlength="30" class="modern-input" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户状态" prop="status">
              <el-radio-group v-model="adminInfo.status" class="modern-radio-group">
                <el-radio :label="1" class="modern-radio">正常</el-radio>
                <el-radio :label="2" class="modern-radio">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户岗位" prop="postId">
              <el-select placeholder="请选择岗位" v-model="adminInfo.postId" class="modern-select">
                <el-option v-for="item in postList" :key="item.id" :label="item.postName" :value="item.id"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户角色" prop="roleId">
              <el-select placeholder="请选择角色" v-model="adminInfo.roleId" class="modern-select">
                <el-option v-for="item in roleList" :key="item.id" :label="item.roleName" :value="item.id"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <el-form-item label="个人简介" prop="note">
              <el-input v-model="adminInfo.note" type="textarea" placeholder="请输入内容" class="modern-textarea"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <span slot="footer" class="dialog-footer">
          <el-button type="primary" @click="editAdminInfo" class="modern-btn primary-btn">确 定</el-button>
          <el-button type="info" @click="editDialogVisible = false" class="modern-btn secondary-btn">取 消</el-button>
        </span>
      </el-form>
    </el-dialog>
  </div>
</template>



<script>
import Treeselect from 'vue3-treeselect'
import 'vue3-treeselect/dist/vue3-treeselect.css'
import storage from '@/utils/storage'
import {
  Search,
  Refresh,
  Plus,
  Edit,
  Delete,
  Key,
  OfficeBuilding,
  DataBoard,
  FolderAdd,
  EditPen
} from '@element-plus/icons-vue'

export default {
  components: {
    Treeselect
  },
  data() {
    return {
      uploadIconUrl: '/api/v1/upload',
      uploadHeaders: {
        Authorization: "Bearer " + storage.getItem("token"),
      },
      editIcon: '',
      editAvatarKey: Date.now(),
      expandedKeys: [], // 用于跟踪展开的节点key
      deptSearchText: '', // 部门搜索文本
      isExpanded: true, // 默认展开状态
      // 右键菜单相关
      contextMenuVisible: false,
      contextMenuPosition: { x: 0, y: 0 },
      contextMenuData: {},
      // 部门管理对话框
      deptDialogVisible: false,
      dialogMode: 'create', // create, edit, rename
      deptForm: {
        id: null,
        deptName: '',
        parentId: null,
        sort: 0,
        status: 1,
        deptType: 2 // 默认为中心部门
      },
      deptRules: {
        deptName: [{ required: true, message: '请输入部门名称', trigger: 'blur' }],
        deptType: [{ required: true, message: '请选择部门类型', trigger: 'change' }],
        parentId: [
          {
            validator: (rule, value, callback) => {
              // 如果是子部门(deptType=3)，必须选择父级部门
              if (this.deptForm.deptType === 3 && !value) {
                callback(new Error('子部门必须选择父级部门'))
              } else {
                callback()
              }
            },
            trigger: 'change'
          }
        ]
      },
      deptSelectOptions: [],
      statusList: [{
        value: '1',
        label: '启用',
      }, {
        value: '2',
        label: '停用',
      }],
      Loading: false,
      queryParams: {},
      adminList: [],
      total: 0,
      addDialogVisible: false,
      deptList: [],// 确保初始化为空数组
      roleList: [],
      postList: [],
      addForm: {
        username: '',
        password: '',
        deptId: undefined,
        postId: undefined,
        roleId: undefined,
        email: '',
        nickname: '',
        status: 1,
        phone: '',
        note: ''
      },
      addFormRules: {
        deptId: [{required: true, message: '请选择部门', trigger: 'blur'}],
        postId: [{required: true, message: '请选择岗位', trigger: 'blur'}],
        roleId: [{required: true, message: '请选择角色', trigger: 'blur'}],
        username: [{required: true, message: '请输入用户账号', trigger: 'blur'}],
        password: [{required: true, message: '请输入用户密码', trigger: 'blur'}],
        status: [{required: true, message: '请选择状态', trigger: 'blur'}],
        email: [{required: true, message: '请输入用户邮箱', trigger: 'blur'}],
        nickname: [{required: true, message: '请输入用户昵称', trigger: 'blur'}],
        phone: [{required: true, message: '请输入用户手机', trigger: 'blur'}]
      },
      editDialogVisible: false,
      adminInfo: {},
      editFormRules: {
        deptId: [{required: true, message: '请选择部门', trigger: 'blur'}],
        postId: [{required: true, message: '请选择岗位', trigger: 'blur'}],
        roleId: [{required: true, message: '请选择角色', trigger: 'blur'}],
        username: [{required: true, message: '请输入用户账号', trigger: 'blur'}],
        status: [{required: true, message: '请选择状态', trigger: 'blur'}],
        email: [{required: true, message: '请输入用户邮箱', trigger: 'blur'}],
        nickname: [{required: true, message: '请输入用户昵称', trigger: 'blur'}],
        phone: [{required: true, message: '请输入用户手机', trigger: 'blur'}]
      },
      defaultProps: {
        children: 'children',
        label: 'deptName'
      },
      currentDeptId: null
    }
  },
  created() {
    this.getDeptVoList()
    this.getRoleVoList()
    this.getPostVoList()
    this.getAdminList()
  },
  mounted() {
    // 点击其他地方隐藏右键菜单
    document.addEventListener('click', this.hideContextMenu)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.hideContextMenu)
  },
  watch: {
    expandedKeys: {
      handler(newKeys) {
        this.$nextTick(() => {
          const tree = this.$refs.deptTree
          if (tree) {
            // 如果是空数组，先获取所有展开的节点然后逐个折叠
            if (newKeys.length === 0) {
              const currentExpanded = tree.store.nodesMap
              Object.keys(currentExpanded).forEach(key => {
                const node = currentExpanded[key]
                if (node.expanded) {
                  node.collapse()
                }
              })
            } else {
              tree.setExpandedKeys(newKeys)
            }
          }
        })
      },
      deep: true,
      immediate: false
    },
    deptList: {
      handler(newList) {
        // 更新部门选择器选项
        this.deptSelectOptions = this.formatDeptOptions(newList)
      },
      deep: true,
      immediate: true
    }
  },
  methods: {
    // 表格行样式
    tableRowClassName({row, rowIndex}) {
      return 'modern-table-row';
    },
    
    // 查询列表
    async getAdminList() {
      this.Loading = true
      const {data: res} = await this.$api.queryAdminList(this.queryParams)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.adminList = res.data.list
        this.total = res.data.total
        this.Loading = false
      }
    },

    // 根据部门查询用户
    async loadAdminsByDept(deptId) {
      console.log('开始加载部门 ID:', deptId);
      this.Loading = true;
      const { data: res } = await this.$api.deptUsers(deptId);
      console.log('接口返回结果:', res);

      if (res.code !== 200) {
        this.$message.error(res.message);
      } else {
        if (!res.data || res.data.length === 0) {
          // 🔥 显式处理无数据的情况
          this.adminList = [];
          this.total = 0;
          this.queryParams.pageNum = 1;
          this.Loading = false;
          return;
        }

        // 添加字段映射逻辑
        const users = res.data.map(item => ({
          id: item.id,
          username: item.username,
          nickname: item.nickname,
          status: item.status,
          icon: item.icon,
          email: item.email,
          phone: item.phone,
          note: item.note,
          createTime: item.create_time || item.createTime,
          deptName: item.dept_name || item.deptName,
          postName: item.post_name || item.postName,
          roleName: item.role_name || item.roleName
        }));
        this.adminList = users;
        this.total = users.length;
        this.queryParams.pageNum = 1; // 回到第一页
      }
      this.Loading = false;
    }

    ,

    // 搜索按钮操作
    handleQuery() {
      this.getAdminList();
    },

    // 重置按钮操作
    resetQuery() {
      this.queryParams = {}
      this.getAdminList();
      this.$message.success("重置成功")
    },

    // pageSize
    handleSizeChange(newSize) {
      this.queryParams.pageSize = newSize
      this.getAdminList()
    },

    // pageNum
    handleCurrentChange(newPage) {
      this.queryParams.pageNum = newPage
      this.getAdminList()
    },

    // 修改用户状态
    async adminUpdateStatus(row) {
      let text = row.status === 2 ? "停用" : "启用";
      const confirmResult = await this.$confirm('确认要"' + text + '""' + row.username + '"用户吗?', "警告", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).catch(err => err);

      if (confirmResult !== 'confirm') {
        await this.getAdminList();
        return this.$message.info('已取消删除');
      }

      // 更新状态
      const { data: res } = await this.$api.updateAdminStatus(row.id, row.status);
      if (res.code === 200) {
        this.$message.success(text + "成功");

        // 手动更新当前行状态，避免重新加载整个列表
        const index = this.adminList.findIndex(item => item.id === row.id);
        if (index > -1) {
          this.$set(this.adminList, index, {...row}); // 强制触发 Vue 响应式更新
        }
      } else {
        this.$message.error(res.message);
      }
    },

    // 部门下拉列表
    async getDeptVoList() {
      const { data: res } = await this.$api.queryDeptList()
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        console.log('原始部门数据:', res.data) // 调试：查看原始数据
        const treeData = this.$handleTree.handleTree(res.data, "id")
        console.log('处理后的树数据:', treeData) // 调试：查看处理后的数据
        this.deptList = treeData

        // 默认展开所有节点
        this.$nextTick(() => {
          if (this.isExpanded && this.deptList && this.deptList.length > 0) {
            const allKeys = this.getAllNodeKeys(this.deptList)
            this.expandedKeys = [...allKeys]

            // 确保树组件已渲染后再设置展开状态
            const tree = this.$refs.deptTree
            if (tree) {
              tree.setExpandedKeys(allKeys)
            }
          }
        })
      }
    }
    ,

    // 角色下拉列表
    async getRoleVoList() {
      const {data: res} = await this.$api.querySysRoleVoList()
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.roleList = res.data
      }
    },

    // 岗位下拉列表
    async getPostVoList() {
      const {data: res} = await this.$api.querySysPostVoList()
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.postList = res.data
      }
    },

    // 监听添加用户对话框关闭
    addDialogClosed() {
      this.$refs.addFormRefForm.resetFields()
    },

    // 新增用户
    addAdmin() {
      this.$refs.addFormRefForm.validate(async valid => {
        if (!valid) return
        const {data: res} = await this.$api.addAdmin(this.addForm);
        if (res.code !== 200) {
          this.$message.error(res.message)
        } else {
          this.$message.success('新增用户成功')
          this.addDialogVisible = false
          await this.getAdminList()
        }
      })
    },

    // 展示编辑用户的对话框
    async showEditAdminDialog(id) {
      const {data: res} = await this.$api.adminInfo(id) // [adminInfo](file:///Users/apple/Desktop/zhangfan/2025-06-30/gin-web/src/api/index.js#L333-L342)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.adminInfo = res.data
        
        // 如果用户有部门ID，获取部门信息以确保treeselect正确显示
        if (this.adminInfo.deptId) {
          try {
            const {data: deptRes} = await this.$api.deptInfo(this.adminInfo.deptId)
            if (deptRes.code === 200) {
              // 确保部门信息在treeselect选项中可用
              console.log('Department info loaded:', deptRes.data)
            }
          } catch (error) {
            console.warn('Failed to load department info:', error)
          }
        }
        
        this.editDialogVisible = true
      }
    },

    // 监听修改用户对话框的关闭事件
    editDialogClosed() {
      this.$refs.editFormRefForm.resetFields()
      this.editIcon = ''
      this.editAvatarKey = Date.now()
    },

    // 编辑头像上传成功回调
    handleEditAvatarSuccess(res) {
      if (res.code === 200) {
        this.editIcon = res.data
        this.editAvatarKey = Date.now()
      } else {
        this.$message.error(res.message || '头像上传失败')
      }
    },

    // 修改用户信息并提交
    editAdminInfo() {
      this.$refs.editFormRefForm.validate(async valid => {
        if (!valid) return
        const updateData = {
          ...this.adminInfo,
          icon: this.editIcon || this.adminInfo.icon
        }
        const {data: res} = await this.$api.adminUpdate(updateData)
        if (res.code !== 200) {
          this.$message.error(res.message)
        } else {
          this.editDialogVisible = false
          this.editIcon = ''
          this.editAvatarKey = Date.now()
          await this.getAdminList()
          this.$message.success('修改用户成功')
        }
      })
    },

    // 删除
    async handleAdminDelete(row) {
      const confirmResult = await this.$confirm('是否确认删除用户为"' + row.username + '"的数据项？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消删除')
      }
      const {data: res} = await this.$api.deleteAdmin(row.id) // [deleteAdmin](file:///Users/apple/Desktop/zhangfan/2025-06-30/gin-web/src/api/index.js#L373-L382)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.getAdminList()
      }
    },

    // 重置密码
    handleResetPwd(row) {
      this.$prompt('请输入"' + row.username + '"的新密码', "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        closeOnClickModal: false,
        inputPattern: /^.{5,20}$/,
        inputErrorMessage: "用户密码长度必须介于 5 和 20 之间"
      }).then(({value}) => {
        this.$api.resetPassword(row.id, value).then(() => { // [resetPassword](file:///Users/apple/Desktop/zhangfan/2025-06-30/gin-web/src/api/index.js#L363-L371)
          this.$message.success("修改成功，新密码是：" + value);
        });
      }).catch(() => {
      });
    },

// 点击部门节点
    handleNodeExpand(data, node) {
      if (!this.expandedKeys.includes(node.key)) {
        this.expandedKeys.push(node.key)
      }
    },

    handleNodeCollapse(data, node) {
      const collapsedKeys = this.getNodeAndDescendantKeys(data)
      this.expandedKeys = this.expandedKeys.filter(key => !collapsedKeys.includes(key))
    },

    handleDeptClick(node, element) {
      let deptId;

      // 多种方式尝试获取 deptId
      if (element && element.data && element.data.id) {
        deptId = element.data.id;
      } else if (element && element.id) {
        deptId = element.id;
      } else if (node && node.key) {
        deptId = node.key;
      }
      console.log('最终获取到的 deptId:', deptId);
      if (!deptId) {
        this.$message.warning("无法获取部门ID");
        return;
      }
      this.currentDeptId = deptId;
      this.loadAdminsByDept(deptId);
    },

    // 切换全部展开/折叠
    toggleAll() {
      this.isExpanded = !this.isExpanded
      if (this.isExpanded) {
        this.expandAllNodes()
      } else {
        this.collapseAllNodes()
      }
    },

    // 展开所有节点
    expandAllNodes() {
      this.$nextTick(() => {
        const tree = this.$refs.deptTree
        if (tree && this.deptList && this.deptList.length > 0) {
          const allKeys = this.getAllNodeKeys(this.deptList)
          this.expandedKeys = [...allKeys]
          // 强制更新树的展开状态
          tree.setExpandedKeys(allKeys)
        }
      })
    },

    // 折叠所有节点
    collapseAllNodes() {
      this.$nextTick(() => {
        const tree = this.$refs.deptTree
        if (tree) {
          this.expandedKeys = []
          // 手动折叠所有已展开的节点
          const currentExpanded = tree.store.nodesMap
          Object.keys(currentExpanded).forEach(key => {
            const node = currentExpanded[key]
            if (node.expanded) {
              node.collapse()
            }
          })
        }
      })
    },

    // 递归获取所有节点的key
    getAllNodeKeys(nodes) {
      let keys = []
      nodes.forEach(node => {
        keys.push(node.id)
        if (node.children && node.children.length > 0) {
          keys = keys.concat(this.getAllNodeKeys(node.children))
        }
      })
      return keys
    },

    // 获取当前节点及其所有子孙节点 key
    getNodeAndDescendantKeys(node) {
      if (!node) return []

      let keys = [node.id]
      if (node.children && node.children.length > 0) {
        keys = keys.concat(this.getAllNodeKeys(node.children))
      }
      return keys
    },

    // 刷新部门树
    handleRefresh() {
      this.getDeptVoList()
      this.$message.success('刷新成功')
    },

    // 搜索部门
    handleDeptSearch() {
      // 可以在这里实现搜索逻辑
      console.log('搜索部门:', this.deptSearchText)
    },

    // 显示右键菜单
    showContextMenu(event, node, data) {
      event.preventDefault()
      event.stopPropagation()

      // 获取部门卡片容器的位置信息
      const deptCard = this.$el.querySelector('.dept-card')
      const cardRect = deptCard.getBoundingClientRect()

      // 计算相对于部门卡片的位置
      this.contextMenuPosition = {
        x: event.clientX - cardRect.left + 10,
        y: event.clientY - cardRect.top + 5
      }

      this.contextMenuData = { node, data }
      this.contextMenuVisible = true
    },

    // 隐藏右键菜单
    hideContextMenu() {
      this.contextMenuVisible = false
    },

    // 部门类型变化时的处理
    handleDeptTypeChange(newType) {
      // 如果改为中心部门，设置父级为公司
      if (newType === 2) {
        const companyNode = this.deptList.find(dept => dept.deptType === 1)
        if (companyNode) {
          this.deptForm.parentId = companyNode.id
        } else {
          this.deptForm.parentId = null
        }
      } else if (newType === 3) {
        // 如果改为子部门，保持原来点击的部门作为父级
        if (this.contextMenuData.data) {
          // 如果点击的是中心部门或子部门，使用它作为父级
          this.deptForm.parentId = this.contextMenuData.data.id
        } else {
          // 如果没有点击数据，清空父级并提示
          this.deptForm.parentId = null
          this.$message.warning('子部门必须有父级部门')
        }
      }
    },

    // 创建部门
    createDept() {
      this.dialogMode = 'create'

      // 根据点击的位置智能判断部门类型和父级
      if (!this.contextMenuData.data) {
        // 点击空白区域，创建中心部门
        this.deptForm = {
          id: null,
          deptName: '',
          parentId: null,
          sort: 0,
          status: 1,
          deptType: 2 // 中心部门
        }
      } else {
        const clickedDept = this.contextMenuData.data
        // 点击某个部门，根据其类型决定
        if (clickedDept.deptType === 1) {
          // 点击公司，创建中心部门
          this.deptForm = {
            id: null,
            deptName: '',
            parentId: clickedDept.id,
            sort: 0,
            status: 1,
            deptType: 2
          }
        } else if (clickedDept.deptType === 2) {
          // 点击中心，创建子部门
          this.deptForm = {
            id: null,
            deptName: '',
            parentId: clickedDept.id,
            sort: 0,
            status: 1,
            deptType: 3
          }
        } else {
          // 点击子部门，创建同级子部门
          this.deptForm = {
            id: null,
            deptName: '',
            parentId: clickedDept.parentId,
            sort: 0,
            status: 1,
            deptType: 3
          }
        }
      }

      this.deptDialogVisible = true
      this.hideContextMenu()
    },

    // 重命名部门
    renameDept() {
      this.dialogMode = 'rename'
      this.deptForm = {
        id: this.contextMenuData.data.id,
        deptName: this.contextMenuData.data.deptName,
        parentId: this.contextMenuData.data.parentId || null,
        sort: this.contextMenuData.data.sort || 0,
        status: this.contextMenuData.data.deptStatus || 1,
        deptType: this.contextMenuData.data.deptType || 3
      }
      this.deptDialogVisible = true
      this.hideContextMenu()
    },

    // 编辑部门
    async editDept() {
      this.dialogMode = 'edit'
      const deptId = this.contextMenuData.data.id

      try {
        const { data: res } = await this.$api.deptInfo(deptId)
        if (res.code === 200) {
          this.deptForm = {
            id: res.data.id,
            deptName: res.data.deptName,
            parentId: res.data.parentId || null,
            sort: res.data.sort || 0,
            status: res.data.deptStatus || 1,
            deptType: res.data.deptType || 3
          }
          this.deptDialogVisible = true
        } else {
          this.$message.error(res.message || '获取部门信息失败')
        }
      } catch (error) {
        this.$message.error('获取部门信息失败')
        console.error(error)
      }

      this.hideContextMenu()
    },

    // 删除部门
    async deleteDept() {
      const deptData = this.contextMenuData.data

      try {
        await this.$confirm(
          `确定要删除部门"${deptData.deptName}"吗？如果该部门下有子部门或用户，将无法删除。`,
          '确认删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        const { data: res } = await this.$api.deleteDept(deptData.id)
        if (res.code === 200) {
          this.$message.success('删除成功')
          await this.getDeptVoList()
        } else {
          this.$message.error(res.message || '删除失败')
        }
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
          console.error(error)
        }
      }

      this.hideContextMenu()
    },

    // 获取对话框标题
    getDialogTitle() {
      switch (this.dialogMode) {
        case 'create':
          return '创建部门'
        case 'rename':
          return '重命名部门'
        case 'edit':
          return '编辑部门'
        default:
          return '部门管理'
      }
    },

    // 提交部门表单
    async submitDeptForm() {
      try {
        await this.$refs.deptFormRef.validate()

        // 查找公司节点
        const companyNode = this.deptList.find(dept => dept.deptType === 1)
        const companyId = companyNode ? companyNode.id : 1

        const formData = {
          ...this.deptForm,
          deptStatus: this.deptForm.status,
          // 如果是中心部门(deptType=2)，父级必须设为公司ID
          // 如果是子部门(deptType=3)，必须有父级部门
          parentId: this.deptForm.deptType === 2 ? companyId : this.deptForm.parentId
        }

        // 验证子部门必须有父级
        if (formData.deptType === 3 && !formData.parentId) {
          this.$message.error('子部门必须选择父级部门')
          return
        }

        // 验证中心部门不能挂在非公司节点下
        if (formData.deptType === 2) {
          formData.parentId = companyId
          console.log('中心部门自动设置 parentId 为公司:', companyId)
        }

        console.log('提交的部门数据:', formData)

        let res
        if (this.dialogMode === 'create') {
          const { data } = await this.$api.addDept(formData)
          res = data
        } else if (this.dialogMode === 'edit' || this.dialogMode === 'rename') {
          const { data } = await this.$api.deptUpdate(formData)
          res = data
        }

        if (res.code === 200) {
          this.$message.success(
            this.dialogMode === 'create' ? '创建成功' : '更新成功'
          )
          this.deptDialogVisible = false
          await this.getDeptVoList()
        } else {
          this.$message.error(res.message || '操作失败')
        }
      } catch (error) {
        if (error !== false) {
          console.error('表单验证失败:', error)
        }
      }
    },

    // 对话框关闭处理
    handleDialogClose() {
      this.$refs.deptFormRef?.resetFields()
    },

    // 格式化部门选项供选择器使用
    formatDeptOptions(depts) {
      return depts.map(dept => ({
        id: dept.id,
        label: dept.deptName,
        deptName: dept.deptName,
        deptType: dept.deptType,
        parentId: dept.parentId,
        sort: dept.sort,
        status: dept.status,
        children:
          dept.children && dept.children.length > 0
            ? this.formatDeptOptions(dept.children)
            : []
      }))
    }
  }
}
</script>

<style scoped>
/* 主容器 - 渐变背景 */
.modern-admin-container {
  min-height: 100vh;
  background: #eef1f6;
  padding: 5px;
  position: relative;
}

/* 主标题样式 */
.main-title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
   background-clip: text;
  -webkit-text-fill-color: transparent;
}

/* 毛玻璃卡片基础样式 */
.glass-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-sizing: border-box;
  width: 100%;
}

.glass-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.main-card {
  padding: 24px;
}

/* 布局容器 */
.user-management-container {
  display: flex;
  gap: 20px;
  min-height: 600px;
  overflow: hidden;
}

/* ================ 部门树样式开始 ================ */
.dept-tree-section {
  width: 280px;
  flex-shrink: 0;
  min-width: 200px;
  position: relative;
}

/* 🎨 高科技卡片设计 */
.dept-card {
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.95), rgba(255, 255, 255, 0.85));
  backdrop-filter: blur(20px);
  border-radius: 20px;
  box-shadow:
    0 20px 40px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
  height: 100%;
  overflow: hidden;
  position: relative;
}

.dept-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.5), transparent);
  z-index: 1;
}

/* 🎯 科技感标题区域 */
.dept-card .card-header {
  padding: 20px 20px 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  position: relative;
  margin-bottom: 0;
}

.title-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pulse-ring {
  position: absolute;
  width: 35px;
  height: 35px;
  border: 2px solid rgba(102, 126, 234, 0.3);
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.main-icon {
  font-size: 20px;
  color: #667eea;
  z-index: 2;
  background: rgba(102, 126, 234, 0.1);
  padding: 8px;
  border-radius: 8px;
}

@keyframes pulse {
  0% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.2); opacity: 0.7; }
  100% { transform: scale(1.4); opacity: 0; }
}

.title-content .title {
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 2px 0;
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  white-space: nowrap;
}

.subtitle {
  font-size: 11px;
  color: rgba(102, 126, 234, 0.6);
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 500;
}

.stats-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: rgba(102, 126, 234, 0.05);
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.1);
  position: relative;
}

.stats-text {
  font-size: 16px;
  font-weight: 700;
  color: #667eea;
  line-height: 1;
}

.stats-label {
  font-size: 10px;
  color: rgba(102, 126, 234, 0.6);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* 🔍 高级搜索区域 */
.search-container {
  padding: 15px 20px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.08);
}

.search-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.search-input-wrapper {
  position: relative;
}

.tech-input :deep(.el-input__wrapper) {
  background: rgba(102, 126, 234, 0.03);
  border: 1px solid rgba(102, 126, 234, 0.15);
  border-radius: 12px;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.02);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.tech-input :deep(.el-input__wrapper):hover {
  border-color: rgba(102, 126, 234, 0.3);
  background: rgba(102, 126, 234, 0.05);
}

.tech-input :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.search-icon {
  color: rgba(102, 126, 234, 0.6);
  transition: all 0.3s ease;
}

.tech-input:focus-within .search-icon {
  color: #667eea;
}

/* 操作按钮容器 */
.action-buttons-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  position: absolute;
  bottom: -21px;
  right: 55px;
}

.toggle-btn {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 3px;
  transition: all 0.3s ease;
  background: transparent;
  border: none;
  z-index: 10;
}

.toggle-btn:hover {
  transform: scale(1.1);
}

.toggle-icon {
  width: 20px;
  height: 20px;
  transition: all 0.3s ease;
  filter: hue-rotate(30deg) saturate(1.2) brightness(1.1);
}

.toggle-icon.expanded {
  transform: rotate(180deg);
  filter: hue-rotate(60deg) saturate(1.4) brightness(1.2);
}

/* 刷新按钮样式 */
.refresh-btn {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 3px;
  transition: all 0.3s ease;
  background: transparent;
  border: none;
  z-index: 10;
}

.refresh-btn:hover {
  transform: scale(1.1);
  background: rgba(102, 126, 234, 0.1);
}

.refresh-btn:active {
  transform: scale(0.95);
}

.refresh-icon {
  font-size: 16px;
  color: #667eea;
  transition: all 0.3s ease;
}

.refresh-btn:hover .refresh-icon {
  transform: rotate(180deg);
  color: #764ba2;
}

/* 🌳 高科技树形结构 */
.tree-container {
  padding: 15px 0;
  max-height: calc(100vh - 300px);
  overflow-y: auto;
  position: relative;
}

.tech-tree {
  border: none;
  background: transparent;
}

.tech-tree :deep(.el-tree-node) {
  margin: 0;
}

.tech-tree :deep(.el-tree-node__content) {
  height: auto;
  padding: 0;
  background: transparent;
  border-radius: 0;
  position: relative;
  pointer-events: auto;
}

/* 调整树节点展开按钮位置和旋转效果 - 放在图标前面 */
.tech-tree :deep(.el-tree-node__expand-icon) {
  position: absolute;
  left: 25px !important;
  top: 50%;
  transform: translateY(-50%);
  z-index: 5;
  color: #667eea;
  font-size: 12px;
  transition: transform 0.3s ease;
}

.tech-tree :deep(.el-tree-node__expand-icon.expanded) {
  transform: translateY(-50%) rotate(90deg);
}

.tech-tree :deep(.el-tree-node__expand-icon.is-leaf) {
  display: none;
}

.tree-node {
  position: relative;
  margin: 0px 15px 0px 15px;
  border-radius: 8px;
  overflow: visible;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.tree-node:hover {
  transform: none;
}

/* 只有直接悬停父节点时才显示阴影效果 */
.tree-node.parent-node:hover > .hover-effect {
  opacity: 1;
  transform: scale(1);
}

/* 确保子节点不受父节点阴影影响 */
.tree-node:not(.parent-node) .hover-effect {
  display: none;
}

.hover-effect {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.08), rgba(118, 75, 162, 0.08));
  border-radius: 12px;
  opacity: 0;
  transform: scale(1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1;
}

.node-content {
  display: flex;
  align-items: center;
  padding: 2px 12px;
  position: relative;
  z-index: 2;
}

.parent-node .node-content {
  padding: 3px 12px;
}

.node-icon-wrapper {
  margin-left: 32px;
  margin-right: 2px;
  display: flex;
  align-items: center;
}

.parent-icon {
  position: relative;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.parent-icon .icon-bg {
  display: none;
}

.parent-icon.expanded .icon-bg {
  display: none;
}

.parent-icon-img {
  width: 16px;
  height: 16px;
  z-index: 1;
  transition: all 0.3s ease;
}

.child-icon {
  position: relative;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dot-indicator {
  position: absolute;
  width: 6px;
  height: 6px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 50%;
  top: 50%;
  left: -8px;
  transform: translateY(-50%);
}

.child-icon-img {
  width: 18px;
  height: 18px;
}

/* 中心部门图标样式 */
.center-icon {
  position: relative;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.center-icon-img {
  width: 18px;
  height: 18px;
}

.node-text {
  flex: 1;
  position: relative;
}

.node-label {
  font-size: 14px;
  font-weight: 400;
  color: #2c3e50;
  transition: all 0.3s ease;
}

.parent-node .node-label {
  font-weight: 400;
  font-size: 15px;
}

.connection-line {
  position: absolute;
  bottom: -6px;
  left: 0;
  height: 1px;
  width: 0;
  background: linear-gradient(90deg, #667eea, transparent);
  transition: width 0.3s ease;
}

.parent-node:hover .connection-line {
  width: 0;
}

/* 🎯 选中状态 - 只针对被选中的节点本身 */
.tech-tree :deep(.el-tree-node.is-current > .el-tree-node__content) .tree-node {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.08), rgba(118, 75, 162, 0.08));
  border: none;
  transform: none;
}

.tech-tree :deep(.el-tree-node.is-current > .el-tree-node__content) .hover-effect {
  opacity: 0 !important;
}

.tech-tree :deep(.el-tree-node.is-current > .el-tree-node__content) .node-label {
  color: #2c3e50;
  font-weight: 500;
}

/* 📱 滚动条样式 */
.tree-container::-webkit-scrollbar {
  width: 4px;
}

.tree-container::-webkit-scrollbar-track {
  background: rgba(102, 126, 234, 0.05);
  border-radius: 2px;
}

.tree-container::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.3), rgba(118, 75, 162, 0.3));
  border-radius: 2px;
}

.tree-container::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.5), rgba(118, 75, 162, 0.5));
}

/* 🚀 数据加载动画 */
@keyframes dataLoad {
  0% { opacity: 0; transform: translateY(20px); }
  100% { opacity: 1; transform: translateY(0); }
}

.tree-node {
  animation: dataLoad 0.4s ease-out;
}

.tree-node:nth-child(1) { animation-delay: 0.1s; }
.tree-node:nth-child(2) { animation-delay: 0.2s; }
.tree-node:nth-child(3) { animation-delay: 0.3s; }
.tree-node:nth-child(4) { animation-delay: 0.4s; }
.tree-node:nth-child(5) { animation-delay: 0.5s; }

/* ================ 部门树样式结束 ================ */

/* 右键菜单样式 */
.context-menu {
  position: absolute;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 4px;
  z-index: 9999;
  width: 140px;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s ease;
  font-size: 14px;
  color: #2c3e50;
}

.context-menu-item:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.context-menu-item.danger {
  color: #f56c6c;
}

.context-menu-item.danger:hover {
  background: rgba(245, 108, 108, 0.1);
  color: #f56c6c;
}

.context-menu-item .el-icon {
  font-size: 14px;
}


.user-table-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
  min-width: 0;
  overflow: hidden;
}

/* 卡片头部 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(103, 126, 234, 0.1);
}

.section-title {
  color: #2c3e50;
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  background: linear-gradient(45deg, #667eea, #764ba2);
   background-clip: text;
  -webkit-text-fill-color: transparent;
}

/* 搜索表单样式 */
.search-card {
  padding: 20px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
  overflow: hidden;
}

.modern-form {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 16px;
  width: 100%;
  box-sizing: border-box;
}

.modern-form .el-form-item {
  margin-bottom: 0;
  flex-shrink: 0;
}

.modern-form .el-form-item__label {
  color: #606266 !important;
  font-weight: 500;
}

.modern-form .el-form-item:last-child {
  display: flex;
  gap: 12px;
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

.secondary-btn {
  background: linear-gradient(45deg, #909399, #B1B3B8);
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

/* 输入框样式 */
.modern-input :deep(.el-input__wrapper),
.modern-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.modern-input :deep(.el-input__wrapper):hover,
.modern-select :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.modern-input :deep(.el-input__wrapper.is-focus),
.modern-select :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.modern-input :deep(.el-input__inner),
.modern-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #2c3e50;
}

.modern-input :deep(.el-input__inner::placeholder),
.modern-select :deep(.el-input__inner::placeholder) {
  color: rgba(44, 62, 80, 0.6);
}

.modern-select :deep(.el-select-dropdown) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.modern-date-picker {
  width: 190px;
  min-width: 150px;
}

/* 确保输入框不会超出容器 */
.modern-input,
.modern-select,
.modern-date-picker {
  max-width: 100%;
  box-sizing: border-box;
}

/* 表格样式 */
.table-card {
  padding: 20px;
  flex: 1;
  overflow: hidden;
  min-width: 0;
}

.modern-table {
  background: transparent;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  width: 100%;
}

.modern-table :deep(.el-table) {
  overflow-x: auto;
}

.modern-table :deep(.el-table__body-wrapper) {
  overflow-x: auto;
}

.modern-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.modern-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
  padding: 8px 12px !important;
  height: 40px;
}

.modern-table :deep(.el-table__header th .cell) {
  color: #2c3e50 !important;
  font-weight: 700 !important;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.modern-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.modern-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.modern-table :deep(.el-table__row td) {
  border: none;
  color: #2c3e50;
  padding: 8px 12px !important;
  height: 40px;
}

.modern-table :deep(.el-table__cell) {
  padding: 8px 12px !important;
}

.modern-table :deep(.el-table__empty-block) {
  background: rgba(255, 255, 255, 0.05);
}

.modern-table :deep(.el-table__empty-text) {
  color: rgba(44, 62, 80, 0.6);
}

/* 表格内组件 */
.modern-tag {
  border-radius: 6px;
  font-weight: 500;
}

.modern-avatar {
  border: 2px solid rgba(255, 255, 255, 0.3);
  transition: transform 0.3s ease;
}

.modern-avatar:hover {
  transform: scale(1.1);
}

.modern-switch :deep(.el-switch__core) {
  border-radius: 12px;
}

.modern-text {
  color: #2c3e50;
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

/* 分页样式 */
.pagination-wrapper {
  padding: 20px;
  text-align: center;
  border-top: 1px solid rgba(102, 126, 234, 0.1);
  margin-top: 20px;
}

.modern-pagination :deep(.el-pagination) {
  justify-content: center;
}

.modern-pagination :deep(.el-pagination .el-pager li),
.modern-pagination :deep(.el-pagination .btn-prev),
.modern-pagination :deep(.el-pagination .btn-next),
.modern-pagination :deep(.el-pagination .el-pagination__sizes),
.modern-pagination :deep(.el-pagination .el-pagination__jump) {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #ffffff;
  border-radius: 6px;
  margin: 0 2px;
}

.modern-pagination :deep(.el-pagination .el-pager li:hover),
.modern-pagination :deep(.el-pagination .btn-prev:hover),
.modern-pagination :deep(.el-pagination .btn-next:hover) {
  background: rgba(255, 255, 255, 0.2);
  color: #ffffff;
}

.modern-pagination :deep(.el-pagination .el-pager li.active) {
  background: linear-gradient(45deg, #409EFF, #66B3FF);
  color: #ffffff;
}

/* 对话框样式 */
.modern-dialog :deep(.el-dialog) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.2);
}

.modern-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px 16px 0 0;
  padding: 16px 20px;
}

.modern-dialog :deep(.el-dialog__title) {
  color: #ffffff;
  font-weight: 600;
  font-size: 16px;
}

.modern-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
  color: #ffffff;
  font-size: 16px;
}

.modern-dialog-form {
  padding: 20px;
}

.modern-dialog-form .el-form-item__label {
  color: #606266 !important;
  font-weight: 500;
}

.modern-dialog-form .modern-input :deep(.el-input__wrapper),
.modern-dialog-form .modern-select :deep(.el-input__wrapper) {
  background: rgba(0, 0, 0, 0.02);
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.modern-dialog-form .modern-input :deep(.el-input__inner),
.modern-dialog-form .modern-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #333333;
}

.modern-dialog-form .modern-textarea :deep(.el-textarea__inner) {
  background: rgba(0, 0, 0, 0.02);
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  color: #333333;
  transition: all 0.3s ease;
}

.modern-dialog-form .modern-input :deep(.el-input__wrapper):hover,
.modern-dialog-form .modern-select :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.modern-dialog-form .modern-input :deep(.el-input__wrapper.is-focus),
.modern-dialog-form .modern-select :deep(.el-input__wrapper.is-focus),
.modern-dialog-form .modern-textarea :deep(.el-textarea__inner):focus {
  border-color: #409EFF;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
  background: rgba(255, 255, 255, 1);
}

.modern-dialog-form .modern-select {
  width: 100%;
}

.modern-dialog-form .avatar-uploader {
  border: 2px dashed rgba(103, 126, 234, 0.3);
  border-radius: 12px;
  cursor: pointer;
  width: 80px;
  height: 80px;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  background: rgba(103, 126, 234, 0.05);
}

.modern-dialog-form .avatar-uploader:hover {
  border-color: #667eea;
  background: rgba(103, 126, 234, 0.1);
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(103, 126, 234, 0.2);
}

.modern-dialog-form .avatar-uploader .avatar {
  width: 80px;
  height: 80px;
  display: block;
  border-radius: 8px;
}

.modern-radio-group .modern-radio {
  margin-right: 20px;
}

.modern-dialog-form .modern-radio :deep(.el-radio__label) {
  color: #606266 !important;
}

.modern-dialog-form .modern-radio :deep(.el-radio__input.is-checked .el-radio__inner) {
  border-color: #409EFF;
  background: #409EFF;
}

.modern-treeselect :deep(.vue-treeselect__control) {
  background: rgba(0, 0, 0, 0.02);
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.modern-treeselect :deep(.vue-treeselect__control):hover {
  border-color: #c0c4cc;
}

.modern-treeselect :deep(.vue-treeselect__control--is-focused) {
  border-color: #409EFF;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
  background: rgba(255, 255, 255, 1);
}

.modern-treeselect :deep(.vue-treeselect__placeholder),
.modern-treeselect :deep(.vue-treeselect__single-value),
.modern-treeselect :deep(.vue-treeselect__option) {
  color: #606266 !important;
}

.dialog-footer {
  text-align: center;
  padding: 20px;
  border-top: 1px solid #f0f0f0;
}

/* 响应式布局 */
@media (max-width: 1200px) {
  .user-management-container {
    flex-direction: column;
    gap: 16px;
  }
  
  .dept-tree-section {
    width: 100%;
    order: 2;
  }
  
  .user-table-section {
    order: 1;
  }
  
  .modern-form {
    flex-direction: column;
    align-items: stretch;
  }
  
  .modern-form .el-form-item {
    width: 100%;
  }
  
  .button-group {
    flex-direction: column;
    width: 100%;
  }
  
  .button-group .modern-btn {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .modern-admin-container {
    padding: 10px;
  }
  
  .main-card {
    padding: 16px;
  }
  
  .search-card, .table-card, .tree-card {
    padding: 12px;
  }
  
  .modern-form {
    gap: 8px;
  }
  
  .modern-input :deep(.el-input__inner),
  .modern-select :deep(.el-input__inner),
  .modern-date-picker :deep(.el-input__inner) {
    font-size: 14px;
  }
  
  .modern-table :deep(.el-table) {
    font-size: 12px;
  }
  
  .main-title {
    font-size: 18px;
  }
  
  .section-title {
    font-size: 14px;
  }
}

@media (max-width: 480px) {
  .modern-admin-container {
    padding: 8px;
  }
  
  .main-card {
    padding: 12px;
  }
  
  .search-card, .table-card, .tree-card {
    padding: 8px;
  }
  
  .user-management-container {
    gap: 12px;
  }
  
  .main-title {
    font-size: 16px;
  }
  
  .modern-form .el-form-item__label {
    font-size: 12px !important;
  }
  
  .modern-table :deep(.el-table) {
    font-size: 11px;
  }
}
</style>
