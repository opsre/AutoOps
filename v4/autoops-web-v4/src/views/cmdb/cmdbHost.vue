<template>
  <div class="cmdb-host-management">
    <el-card shadow="hover" class="host-card">
    <!-- 左右布局容器 -->
    <div class="host-management-container">
      <!-- 左侧分组树 -->
      <CmdbGroup
        ref="cmdbGroup"
        :group-list="groupList"
        :expanded-keys="expandedKeys"
        @group-search="handleGroupSearch"
        @group-click="handleGroupClick"
        @node-expand="handleNodeExpand"
        @node-collapse="handleNodeCollapse"
        @collapse-all="handleCollapseAll"
        @expand-all="handleExpandAll"
        @refresh-groups="handleRefreshGroups"
        @create-group="handleCreateGroup"
        @update-group="handleUpdateGroup"
        @delete-group="handleDeleteGroup"
      />
      <!-- 右侧主机管理区域 -->
      <div class="host-table-section">
        <!-- 搜索表单 -->
        <div class="search-section">
          <el-form :inline="true" :model="queryParams" class="demo-form-inline">
            <el-form-item label="主机名称" prop="hostName">
              <el-input
                  size="small"
                  placeholder="请输入主机名称"
                  clearable
                  v-model="queryParams.hostName"
                  @keyup.enter="handleQuery"
                  style="width: 160px;"
              />
            </el-form-item>
            <el-form-item label="IP地址" prop="ip">
              <el-input
                  size="small"
                  placeholder="请输入IP地址"
                  clearable  
                  v-model="queryParams.ip"
                  @keyup.enter="handleQuery"
                  style="width: 120px;"
              />
            </el-form-item>
            <el-form-item label="主机状态" prop="status">
              <el-select size="small" placeholder="请选择状态" v-model="queryParams.status" style="width: 120px;">
                <el-option v-for="item in statusList" :key="item.value" :label="item.label" :value="item.value"></el-option>
              </el-select>
            </el-form-item>

            <!-- 主机统计信息 -->
            <el-form-item v-if="hostStatistics" style="margin-left: 20px;">
              <div class="inline-statistics">
                <div class="stat-item">
                  <el-icon :size="20" color="#409EFF"><Monitor /></el-icon>
                  <span class="stat-text">总计: {{ hostStatistics.total || 0 }}</span>
                </div>
                <div v-for="detail in hostStatistics.details" :key="detail.vendor" class="stat-item">
                  <el-icon v-if="detail.vendor === 1" :size="20" color="#409EFF"><OfficeBuilding /></el-icon>
                  <img v-else-if="detail.vendor === 2" src="@/assets/image/aliyun.png" class="stat-icon"/>
                  <img v-else-if="detail.vendor === 3" src="@/assets/image/tengxun.png" class="stat-icon"/>
                  <img v-else-if="detail.vendor === 4" src="@/assets/image/baidu.svg" class="stat-icon"/>
                  <img v-else-if="detail.vendor === 5" src="@/assets/image/华为.svg" class="stat-icon"/>
                  <img v-else-if="detail.vendor === 6" src="@/assets/image/aws.svg" class="stat-icon"/>
                  <span class="stat-text">{{ detail.count }}</span>
                </div>
              </div>
            </el-form-item>
            <!-- 操作按钮 -->
            <div class="action-section">
            <el-row :gutter="10" class="mb8" style="text-align: left">
              <el-col :span="24">
                <!-- 搜索按钮 - 蓝色 -->
                <el-button type="primary" size="small" @click="handleQuery" style="margin-right: 10px">
                  <el-icon><Search /></el-icon>
                  <span style="margin-left: 4px">搜索</span>
                </el-button>
                
                <!-- 重置按钮 - 黄色 -->
                <el-button type="warning" size="small" @click="resetQuery" style="margin-right: 10px">
                  <el-icon><Refresh /></el-icon>
                  <span style="margin-left: 4px">重置</span>
                </el-button>
                
                <!--新建主机 - 绿色-->
                <el-dropdown
                  ref="createDropdown"
                  @command="handleCreateCommand"
                  @visible-change="handleDropdownVisibleChange"
                  :hide-on-click="true"
                  trigger="click"
                  placement="bottom-start">
                  <el-button
                    type="success"
                    size="small"
                    style="margin-right: 10px"
                    v-authority="['cmdb:ecs:add']"
                    @click.stop="handleCreateClick">
                    <el-icon><Plus /></el-icon>
                    <span style="margin-left: 4px">新建</span>
                    <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="importHost"><el-icon color="#409EFC" :size="20"><Edit /></el-icon>新增主机</el-dropdown-item>
                      <el-dropdown-item command="excelImport"><el-icon color="#409EFC" :size="20"><Folder /></el-icon>Excel创建主机</el-dropdown-item>
                      <el-dropdown-item command="cloudHost"><el-icon color="#409EFC" :size="21"><MostlyCloudy /></el-icon>新增云主机</el-dropdown-item>
                      <el-dropdown-item command="jumpserverImport"><el-icon color="#67C23A" :size="20"><Connection /></el-icon>JumpServer导入</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>

                <!--终端按钮 - 深色-->
                <el-button
                  size="small"
                  v-authority="['cmdb:ecs:terminal']"
                  @click="handleHostSSH"
                  class="terminal-dark-btn"
                  style="margin-left: 10px"
                >
                  <el-icon><Monitor /></el-icon>
                  <span style="margin-left: 4px">终端</span>
                </el-button>

                <!--资产授权按钮 - 深蓝色-->
                <el-button
                  size="small"
                  type="primary"
                  @click="handleHostAuthorization"
                  style="margin-left: 10px; background: #003d99; border-color: #003d99;"
                >
                  <el-icon><Key /></el-icon>
                  <span style="margin-left: 4px">资产授权</span>
                </el-button>

                <!-- 批量操作按钮 -->
                <el-dropdown
                  @command="handleBatchCommand"
                  trigger="click"
                  placement="bottom-start"
                  style="margin-left: 10px">
                  <el-button
                    size="small"
                    :disabled="!hasSelection"
                    :type="hasSelection ? 'primary' : 'info'"
                    v-authority="['cmdb:ecs:batch']">
                    <el-icon><Operation /></el-icon>
                    <span style="margin-left: 4px">批量操作</span>
                    <span v-if="selectionCount > 0" style="margin-left: 4px">({{ selectionCount }})</span>
                    <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="updateGroup">
                        <el-icon color="#409EFC" :size="18"><FolderOpened /></el-icon>
                        修改分组
                      </el-dropdown-item>
                      <el-dropdown-item command="updateVendor">
                        <el-icon color="#67C23A" :size="18"><Edit /></el-icon>
                        修改类型
                      </el-dropdown-item>
                      <el-dropdown-item command="bindTag">
                        <el-icon color="#E6A23C" :size="18"><PriceTag /></el-icon>
                        绑定标签
                      </el-dropdown-item>
                      <el-dropdown-item command="batchDelete" divided>
                        <el-icon color="#F56C6C" :size="18"><Delete /></el-icon>
                        批量删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-col>
            </el-row>
            </div>
          </el-form>
        </div>

        <!-- 主机表格 -->
        <CmdbHostTable
          ref="hostTable"
          :key="$route.fullPath"
          :host-list="hostList"
          :loading="loading"
          @show-detail="showHostDetail"
          @edit-host="showEditHostDialog"
          @show-download="showDownloadDialog"
          @delete-host="handleHostDelete"
          @selection-change="handleHostSelectionChange"
          @refresh-list="getHostList"
        />

        <!-- 分页 -->
        <div class="pagination-section">
          <el-pagination
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              :current-page="queryParams.pageNum"
              :page-sizes="[10, 50, 100]"
              :page-size="queryParams.pageSize"
              layout="total, sizes, prev, pager, next, jumper"
              :total="total"
          ></el-pagination>
        </div>
      </div>
    </div>

    <!-- 新增主机对话框 -->
    <CreateHost
      :visible="addDialogVisible"
      :group-list="groupList"
      :auth-list="authList"
      @close="addDialogVisible = false"
      @submit="addHost"
      @refresh-auth-list="getAuthList"
    />

    <!-- 编辑主机对话框 -->
    <EditHost
      :visible="editDialogVisible"
      :host-info="hostInfo"
      :group-list="groupList"
      :auth-list="authList"
      @close="editDialogVisible = false"
      @submit="editHost"
    />

    <!-- 导入云主机对话框 -->
    <CreateCloud v-model="cloudDialogVisible" @success="handleCloudImportSuccess" />

    <!-- Excel导入对话框 -->
    <CreateExcel v-model="ExcelDialogVisible" @success="handleExcelImportSuccess" />

    <!-- JumpServer导入对话框 -->
    <CreateJumpServer v-model="jumpServerDialogVisible" @success="handleJumpServerImportSuccess" />

    <!-- SSH终端对话框 -->
    <HostSSH 
      v-if="sshDialogVisible"
      :visible="sshDialogVisible"
      :host-id="currentHostId"
      @update:visible="val => {
        sshDialogVisible = val
      }"
    />

    <!-- 主机详情抽屉 -->
    <el-drawer
      v-model="detailDrawer"
      title="主机详情"
      direction="rtl"
      size="40%"
      :before-close="handleDetailClose">

      <!-- 仪表盘部分 -->
      <div class="dashboard-section">
        <div class="gauge-container">
          <div ref="cpuGauge" class="gauge-item"></div>
          <div ref="memoryGauge" class="gauge-item"></div>
          <div ref="diskGauge" class="gauge-item"></div>
        </div>
      </div>

      <!-- 基本信息部分 -->
      <h3 style="margin: 5px 0 10px 0">基本信息</h3>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="主机名称">{{ hostDetail.hostName }}</el-descriptions-item>
        <el-descriptions-item label="主机分组">{{ getGroupName(hostDetail.groupId) }}</el-descriptions-item>
        <el-descriptions-item label="连接地址">
          {{ hostDetail.sshName }}@{{ hostDetail.sshIp }}:{{ hostDetail.sshPort }}
        </el-descriptions-item>
        <el-descriptions-item label="认证类型">
          <el-tag :type="getAuthTypeTagType(hostDetail.sshKeyId)" size="small">
            {{ getAuthTypeName(hostDetail.sshKeyId) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="认证状态">
          <el-tag :type="getStatusTagType(hostDetail.status)">
            {{ getStatusText(hostDetail.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述信息">{{ hostDetail.remark }}</el-descriptions-item>
      </el-descriptions>

      <!-- 扩展信息部分 -->
      <div style="margin: 20px 0 10px 0; display: flex; justify-content: space-between; align-items: center;">
        <h3 style="margin: 0">扩展信息</h3>
        <el-button 
          type="primary" 
          size="mini" 
          icon="Refresh"
          :loading="syncLoading"
          v-authority="['cmdb:ecs:rsync']"
          @click="handleHostSync"
        >
          {{ syncLoading ? '同步中...' : '同步' }}
        </el-button>
      </div>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="实例ID">{{ hostDetail.instanceId }}</el-descriptions-item>
        <el-descriptions-item label="实例名称">{{ hostDetail.name }}</el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ hostDetail.os }}</el-descriptions-item>
        <el-descriptions-item label="CPU">{{ hostDetail.cpu }}核</el-descriptions-item>
        <el-descriptions-item label="内存">{{ hostDetail.memory }}G</el-descriptions-item>
        <el-descriptions-item label="磁盘">{{ hostDetail.disk }}GB</el-descriptions-item>
        <el-descriptions-item label="内网IP">{{ hostDetail.privateIp }}</el-descriptions-item>
        <el-descriptions-item label="公网IP">{{ hostDetail.publicIp || '无' }}</el-descriptions-item>
        <el-descriptions-item label="实例计费方式">{{ hostDetail.billingType }}</el-descriptions-item>
        <el-descriptions-item label="网络计费方式">{{ hostDetail.networkBillingType || '按流量计费' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ hostDetail.createTime }}</el-descriptions-item>
        <el-descriptions-item label="到期时间">{{ hostDetail.expireTime || '无' }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ hostDetail.updateTime }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>

    <!-- 文件管理组件 -->
    <FileManager
      ref="fileManager"
      @upload-success="handleFileOperationSuccess"
      @download-success="handleFileOperationSuccess"
    />
  </el-card>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import configApi from '@/api/config'
import cmdbApi from '@/api/cmdb'
import CreateCloud from './Host/CreateCloud.vue'
import HostSSH from './Host/SSH.vue'
import CreateExcel from './Host/CreateExcel.vue'
import CreateJumpServer from './Host/CreateJumpServer.vue'
import CmdbGroup from './Host/CmdbGroup.vue'
import CmdbHostTable from './Host/CmdbHostTable.vue'
import CreateHost from './Host/CreateHost.vue'
import EditHost from './Host/EditHost.vue'
import FileManager from './Host/FileManager.vue'

export default {
  components: {
    CreateCloud,
    HostSSH,
    CreateExcel,
    CreateJumpServer,
    CmdbGroup,
    CmdbHostTable,
    CreateHost,
    EditHost,
    FileManager
  },
  data() {
    return {
      ExcelDialogVisible: false,
      jumpServerDialogVisible: false,
      expandedKeys: [], // 用于跟踪展开的节点key
      // 批量操作相关
      selectedHosts: [],
      statusList: [
        { value: 2, label: '未认证' },
        { value: 1, label: '认证成功' },
        { value: 3, label: '认证失败' }
      ],
      loading: false,
      queryParams: {
        pageNum: 1,
        pageSize: 10,
        hostName: '',
        ip: '',
        status: '',
        groupId: ''
      },
      hostList: [],
      total: 0,
      hostStatistics: null, // 主机统计数据
      addDialogVisible: false,
      editDialogVisible: false,
      cloudDialogVisible: false,
      groupList: [],
      defaultProps: {
        children: 'children',
        label: 'name'
      },
      currentGroupId: null,
      authList: [],
      addForm: {
        hostName: '',
        ip: '',
        port: 22,
        username: '',
        authId: '',
        groupId: '',
        remark: ''
      },
      addFormRules: {
        hostName: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
        ip: [{ required: true, message: '请输入IP地址', trigger: 'blur' }],
        port: [{ required: true, message: '请输入端口号', trigger: 'blur' }],
        username: [{ required: true, message: '请输入连接用户名', trigger: 'blur' }],
        authId: [{ required: true, message: '请选择认证凭据', trigger: 'change' }],
        groupId: [{ required: true, message: '请选择所属分组', trigger: 'change' }]
      },
      hostInfo: {},
      editFormRules: {
        hostName: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
        ip: [{ required: true, message: '请输入IP地址', trigger: 'blur' }],
        port: [{ required: true, message: '请输入端口号', trigger: 'blur' }],
        username: [{ required: true, message: '请输入连接用户名', trigger: 'blur' }],
        authId: [{ required: true, message: '请选择认证凭据', trigger: 'change' }],
        groupId: [{ required: true, message: '请选择所属分组', trigger: 'change' }]
      },
    // SSH终端对话框相关
    sshDialogVisible: false,
    currentHostId: null,

    // 主机详情相关
    detailDrawer: false,
    syncLoading: false, // 同步加载状态
    hostDetail: {
        hostName: '',
        groupId: '',
        privateIp: '',
        publicIp: '',
        sshIp: '',
        sshName: '',
        sshKeyId: '',
        sshPort: 22,
        remark: '',
        vendor: '',
        region: '',
        instanceId: '',
        os: '',
        status: 0,
        cpu: '',
        memory: '',
        disk: '',
        billingType: '',
        createTime: '',
        expireTime: '',
        updateTime: '',
        name: '',
        cpuUsage: 0,
        memoryUsage: 0,
        diskUsage: 0
      },
    // ECharts 实例
    cpuChart: null,
    memoryChart: null,
    diskChart: null
    }
  },
  computed: {
    // 是否有选中的主机
    hasSelection() {
      return this.selectedHosts.length > 0
    },
    // 选中的主机数量
    selectionCount() {
      return this.selectedHosts.length
    }
  },
  created() {
    this.getAllGroups()
    this.getAuthList()
    // 默认加载所有主机
    this.getHostList()
    // 加载主机统计信息
    this.fetchHostStatistics()
  },

    beforeRouteEnter(to, from, next) {
      next(vm => {
        // 立即获取主机列表
        vm.getHostList().then(() => {
          // 主机列表加载完成后立即触发监控数据加载
          vm.$refs.hostTable?.fetchMonitorData()
        })
      })
    },

    beforeRouteUpdate(to, from, next) {
      // 立即获取主机列表
      this.getHostList().then(() => {
        // 主机列表加载完成后立即触发监控数据加载
        this.$refs.hostTable?.fetchMonitorData()
        next()
      })
    },
  methods: {
    // 获取所有分组
    async getAllGroups() {
      const { data: res } = await this.$api.getAllCmdbGroups()
      if (res.code === 200) {
        this.groupList = res.data
        // 设置默认分组为业务组
        const businessGroup = this.groupList.find(group => group.name === '业务组')
        if (businessGroup) {
          this.addForm.groupId = businessGroup.id
        }
      }
    },

    // 获取主机统计信息
    async fetchHostStatistics() {
      try {
        const { data: res } = await this.$api.getHostCloudStatistics()
        if (res.code === 200) {
          this.hostStatistics = res.data
        }
      } catch (error) {
        console.error('获取主机统计信息失败:', error)
      }
    },

    // 处理分组搜索
    async handleGroupSearch(searchText) {
      this.groupSearchText = searchText
      if (!this.groupSearchText) {
        this.expandedKeys = []
        return
      }
      
      try {
        const { data: res } = await this.$api.getCmdbGroupByName(this.groupSearchText)
        
        if (res.code === 200 && res.data) {
          
          // 获取CmdbGroup组件的树引用
          const cmdbGroupRef = this.$refs.cmdbGroup
          const tree = cmdbGroupRef ? cmdbGroupRef.$refs.groupTree : null
          if (!tree) {
            console.error('树组件引用不存在')
            return
          }



          // 找到匹配的分组并展开其父级
          const findAndExpandParent = (groups, targetId, path = []) => {
            for (const group of groups) {

              if (group.id === targetId) {

                return [...path, group.id]
              }
              if (group.children && group.children.length > 0) {
                const foundPath = findAndExpandParent(group.children, targetId, [...path, group.id])
                if (foundPath) {
                  return foundPath
                }
              }
            }
            return null
          }
          
          // 获取展开路径
          const expandPath = findAndExpandParent(this.groupList, res.data.id)
          
          if (expandPath) {
            // 设置展开的keys
            this.expandedKeys = expandPath.slice(0, -1)
            
            // 强制更新树组件
            this.$nextTick(() => {
              tree.setCurrentKey(res.data.id)
              
              // 确保树组件已更新
              setTimeout(() => {

              }, 500)
            })
          } else {
            console.warn('未找到匹配分组的路径')
            this.$message.warning('未找到匹配的分组路径')
          }
        } else {
          console.warn('未找到匹配分组')
        }
      } catch (error) {
        console.error('搜索分组失败:', error)
        this.$message.error('搜索分组失败: ' + (error?.message || '未知错误'))
      }
    },

    // 获取认证凭据列表
    async getAuthList() {
      try {
        const response = await configApi.getEcsAuthList({
          page: 1,
          pageSize: 100  // 获取认证凭据，用于下拉选择
        })
        
        if (response && response.data) {
          const res = response.data
          
          if (res.code === 200) {
            this.authList = Array.isArray(res.data?.list) ? res.data.list : []
            if (this.authList.length === 0) {
              console.warn('认证凭据列表为空')
            }
          } else {
            console.error('获取认证凭据失败:', res.message || '未知错误')
            this.$message.error(`获取认证凭据失败: ${res.message || '未知错误'}`)
          }
        } else {
          console.error('无效的响应格式:', response)
          this.$message.error('获取认证凭据失败: 无效的响应格式')
        }
      } catch (error) {
        console.error('获取认证凭据异常:', error)
        this.$message.error(`获取认证凭据异常: ${error?.message || '未知错误'}`)
        // 临时添加模拟数据用于测试
        this.authList = [
          { id: 1, name: '默认凭据', username: 'root' },
          { id: 2, name: '测试凭据', username: 'test' }
        ]
        console.warn('使用模拟凭据数据:', this.authList)
      }
    },
    
    // 获取主机列表
    async getHostList() {
      this.loading = true
      try {
        let response
        const { hostName, ip, status, pageNum, pageSize } = this.queryParams

        // 构建分页参数
        const baseParams = {
          page: pageNum,
          pageSize: pageSize,
          _t: Date.now() // 添加时间戳防止缓存
        }

        // 优先级：分组查询 > 条件查询 > 全部查询
        // 如果当前有选中的分组，使用分组查询
        if (this.currentGroupId) {
          response = await this.$api.getCmdbHostsByGroupId(this.currentGroupId, baseParams)
        }
        // 根据查询条件选择API调用
        else if (hostName && !ip && !status) {
          response = await this.$api.GetCmdbHostsByHostNameLike(hostName, baseParams)
        } else if (ip && !hostName && !status) {
          response = await this.$api.GetCmdbHostsByIP(ip, baseParams)
        } else if (status && !hostName && !ip) {
          response = await this.$api.GetCmdbHostsByStatus(status, baseParams)
        } else {
          response = await this.$api.getCmdbHostList(baseParams)
        }


        // 处理axios响应结构
        const axiosResponse = response?.data ? response : { data: response }

        // 严格检查响应格式
        if (!axiosResponse || typeof axiosResponse !== 'object') {
          throw new Error('API返回无效响应格式')
        }

        // 检查响应数据
        const res = axiosResponse.data
        if (!res || typeof res !== 'object') {
          throw new Error('无效的响应数据结构')
        }

        // 检查响应码
        if (res.code === undefined || res.code !== 200) {
          throw new Error(res.message || '获取主机列表失败')
        }

        // 确保data存在，即使为空数组
        if (res.data === undefined) {
          throw new Error('响应缺少data字段')
        }

        // 处理响应数据 - 适配不同API返回格式
        if (Array.isArray(res.data)) {
          // 直接返回数组的情况（如GetCmdbHostsByIP）
          this.hostList = res.data
          this.total = res.data.length
        } else if (res.data?.list) {
          // 返回分页格式的情况（如getCmdbHostList）
          this.hostList = res.data.list
          this.total = res.data.total
          if (res.data.page) {
            this.queryParams.pageNum = res.data.page
          }
          if (res.data.pageSize) {
            this.queryParams.pageSize = res.data.pageSize
          }
        } else {
          // 其他情况
          this.hostList = []
          this.total = 0
        }

        // 主机列表加载完成后立即触发监控数据加载
        this.$nextTick(() => {
          this.$refs.hostTable?.fetchMonitorData()
        })

      } catch (error) {
        console.error('获取主机列表异常:', {
          error: error?.message || '未知错误',
          stack: error?.stack || '无堆栈信息',
          queryParams: this.queryParams,
          currentGroupId: this.currentGroupId
        })
        this.$message.error(`获取主机列表失败: ${error?.message || '未知错误'}`)
        this.hostList = []
        this.total = 0
      } finally {
        this.loading = false
      }
    },
    
    // 处理分组选择变化
    handleGroupChange(value) {
      if (value && value.length > 0) {
        // 取最后一级作为选中分组ID
        this.addForm.groupId = value[value.length - 1]
        this.hostInfo.groupId = value[value.length - 1]
      } else {
        // 如果没有选择分组，设置默认分组
        const defaultGroup = this.groupList.find(item => item.isDefault)
        if (defaultGroup) {
          this.addForm.groupId = defaultGroup.id
          this.hostInfo.groupId = defaultGroup.id
        }
      }
    },

    // 根据分组获取主机
    async getHostsByGroup(groupId) {
      this.loading = true
      this.queryParams.groupId = groupId
      try {
        const { data: res } = await this.$api.getCmdbHostsByGroupId(groupId, {
          page: this.queryParams.pageNum,
          pageSize: this.queryParams.pageSize
        })
        if (res.code === 200) {
          // 处理响应数据 - 适配分页格式
          if (Array.isArray(res.data)) {
            // 直接返回数组的情况（旧格式）
            this.hostList = res.data
            this.total = res.data.length
          } else if (res.data?.list) {
            // 返回分页格式的情况（新格式）
            this.hostList = res.data.list
            this.total = res.data.total
            if (res.data.page) {
              this.queryParams.pageNum = res.data.page
            }
            if (res.data.pageSize) {
              this.queryParams.pageSize = res.data.pageSize
            }
          } else {
            // 其他情况
            this.hostList = []
            this.total = 0
          }

          // 主机列表加载完成后立即触发监控数据加载
          this.$nextTick(() => {
            this.$refs.hostTable?.fetchMonitorData()
          })
        }
      } catch (error) {
        console.error('获取主机列表失败:', error)
        this.hostList = []
        this.total = 0
      } finally {
        this.loading = false
      }
    },
    
    // 点击分组节点
    handleGroupClick(node, element) {
      let groupId
      if (element && element.data && element.data.id) {
        groupId = element.data.id
      } else if (element && element.id) {
        groupId = element.id
      } else if (node && node.key) {
        groupId = node.key
      }

      if (!groupId) {
        this.$message.warning("无法获取分组ID")
        return
      }

      this.currentGroupId = groupId
      // 重置到第一页
      this.queryParams.pageNum = 1
      this.getHostsByGroup(groupId)
    },

    handleNodeExpand(data, node) {
      if (!this.expandedKeys.includes(node.key)) {
        this.expandedKeys.push(node.key)
      }
    },

    handleNodeCollapse(data, node) {
      this.expandedKeys = this.expandedKeys.filter(key => key !== node.key)
    },

    // 折叠所有节点
    handleCollapseAll() {
      this.expandedKeys = []
    },

    // 展开所有节点
    handleExpandAll() {
      const allKeys = []
      const collectKeys = (nodes) => {
        nodes.forEach(node => {
          allKeys.push(node.id)
          if (node.children && node.children.length > 0) {
            collectKeys(node.children)
          }
        })
      }
      collectKeys(this.groupList)
      this.expandedKeys = allKeys
    },

    // 刷新分组列表
    async handleRefreshGroups() {
      try {
        await this.getAllGroups()
        this.$message.success('刷新成功')
      } catch (error) {
        console.error('刷新分组失败:', error)
        this.$message.error('刷新分组失败')
      }
    },
    
    // 搜索按钮操作
    handleQuery() {
      this.queryParams.pageNum = 1
      // 搜索时清除分组选择，使用搜索条件
      this.currentGroupId = null
      this.getHostList()
    },
    
    // 重置按钮操作
    resetQuery() {
      this.queryParams = {
        pageNum: 1,
        pageSize: 10,
        hostName: '',
        ip: '',
        status: '',
        groupId: ''
      }
      this.currentGroupId = null
      this.getHostList()
    },
    
    // pageSize变化
    handleSizeChange(newSize) {
      this.queryParams.pageSize = newSize
      this.getHostList()
    },
    
    // pageNum变化
    handleCurrentChange(newPage) {
      this.queryParams.pageNum = newPage
      this.getHostList()
    },
    
    // 新增主机
    async addHost(requestData) {
      try {
        const { data: res } = await this.$api.createCmdbHost(requestData)

        if (res.code === 200) {
          this.$message.success('新增主机成功，正在同步主机信息...')
          this.addDialogVisible = false

          // 立即刷新一次列表显示新创建的主机
          await this.getHostList()

          // 等待3秒让后端同步主机状态信息，然后再次刷新
          setTimeout(async () => {
            await this.getHostList()
            this.$message.success('主机信息同步完成')
          }, 3000)
        } else if (res.code === 426) {
          this.$message.error(`认证凭据不存在(凭据ID: ${requestData.sshKeyId})，请检查后重试`)
          // 刷新凭据列表
          await this.getAuthList()
        } else {
          this.$message.error(res.message || '新增主机失败')
        }
      } catch (error) {
        console.error('新增主机失败:', error)
        this.$message.error('新增主机失败: ' + error.message)
      }
    },
    
    // 展示编辑主机对话框
    async showEditHostDialog(id) {
      const { data: res } = await this.$api.getCmdbHostById(id)
      if (res.code === 200) {
        this.hostInfo = {
          id: res.data.id,
          hostName: res.data.hostName,
          groupId: res.data.groupId,
          remark: res.data.remark,
          ip: res.data.sshIp,
          port: res.data.sshPort,
          username: res.data.sshName,
          authId: res.data.sshKeyId,
          vendor: res.data.vendor,
          tag: res.data.tag
        }
        this.editDialogVisible = true
      }
    },
    
    // 监听编辑主机对话框关闭
    editDialogClosed() {
      this.$refs.editFormRef.resetFields()
    },
    
    // 编辑主机信息
    async editHost(requestData) {
      try {
        // 验证凭据是否存在
        const authExists = this.authList.some(auth => auth.id === requestData.sshKeyId)
        if (!authExists) {
          this.$message.error('选择的认证凭据不存在，请刷新凭据列表后重试')
          return false
        }

        // 确保端口为数字
        const port = Number(requestData.sshPort)
        if (isNaN(port) || port < 1 || port > 65535) {
          this.$message.error('端口号必须为1-65535之间的数字')
          return false
        }
        const { data: res } = await this.$api.updateCmdbHost(requestData)
        if (res.code === 200) {
          this.$message.success('修改主机成功')
          this.editDialogVisible = false
          this.getHostList()
          return true
        } else if (res.code === 426) {
          this.$message.error(`认证凭据不存在(凭据ID: ${requestData.sshKeyId})，请检查后重试`)
          // 刷新凭据列表
          await this.getAuthList()
          return false
        } else {
          this.$message.error(res.message || '修改主机失败')
          return false
        }
      } catch (error) {
        console.error('修改主机失败:', error)
        this.$message.error('修改主机失败: ' + error.message)
        return false
      }
    },
    
    // 获取状态文本

    // 根据分组ID获取分组名称
    getGroupName(groupId) {
      if (!groupId) return '未分组'
      const findGroup = (groups, id) => {
        for (const group of groups) {
          if (group.id === id) return group.name
          if (group.children && group.children.length > 0) {
            const found = findGroup(group.children, id)
            if (found) return found
          }
        }
        return null
      }
      return findGroup(this.groupList, groupId) || '未知分组'
    },

    // 根据认证凭据ID获取认证类型名称
    getAuthTypeName(authId) {
      if (!authId) return '未设置'
      const auth = this.authList.find(item => item.id === authId)
      if (!auth) return '未知类型'

      switch (auth.type) {
        case 1:
          return '密码认证'
        case 2:
          return '密钥认证'
        case 3:
          return '公钥认证'
        case 4:
          return '跳板机认证'
        default:
          return '未知类型'
      }
    },

    // 获取认证类型标签颜色
    getAuthTypeTagType(authId) {
      if (!authId) return 'info'
      const auth = this.authList.find(item => item.id === authId)
      if (!auth) return 'warning'

      switch (auth.type) {
        case 1:
          return 'primary'  // 密码认证 - 蓝色
        case 2:
          return 'success'  // 密钥认证 - 绿色
        case 3:
          return 'warning'  // 公钥认证 - 橙色
        case 4:
          return 'danger'   // 堡垒机认证 - 红色
        default:
          return 'info'
      }
    },

    // 获取认证状态文本
    getStatusText(status) {
      const statusMap = {
        1: '认证成功',
        2: '未认证',
        3: '认证失败'
      }
      return statusMap[status] || '未知状态'
    },

    // 获取认证状态标签类型
    getStatusTagType(status) {
      const typeMap = {
        1: 'success',
        2: 'warning',
        3: 'danger'
      }
      return typeMap[status] || 'info'
    },

    // 显示主机详情
    async showHostDetail(row) {
      try {
        const { data: res } = await this.$api.getCmdbHostById(row.id)

        if (res.code === 200) {
          this.hostDetail = res.data
          this.detailDrawer = true

          // 获取监控数据
          await this.fetchHostMonitorData(row.id)

          // 等待 DOM 更新后初始化图表
          this.$nextTick(() => {
            this.initGaugeCharts()
          })
        } else {
          console.error('获取主机详情失败:', res.message)
          this.$message.error(res.message || '获取主机详情失败')
        }
      } catch (error) {
        console.error('获取主机详情失败:', error)
        this.$message.error('获取主机详情失败: ' + error.message)
      }
    },

    // 获取主机监控数据
    async fetchHostMonitorData(hostId) {
      try {
        const { data: res } = await this.$api.getHostsMonitorData(hostId)

        if (res.code === 200 && res.data) {
          const monitorData = res.data[hostId]
          if (monitorData) {
            // 更新监控数据，保留两位小数
            this.hostDetail.cpuUsage = parseFloat(monitorData.cpuUsage?.toFixed(2) || 0)
            this.hostDetail.memoryUsage = parseFloat(monitorData.memoryUsage?.toFixed(2) || 0)
            this.hostDetail.diskUsage = parseFloat(monitorData.diskUsage?.toFixed(2) || 0)
          }
        }
      } catch (error) {
        console.error('获取主机监控数据失败:', error)
        // 不显示错误提示，使用默认值
        this.hostDetail.cpuUsage = 0
        this.hostDetail.memoryUsage = 0
        this.hostDetail.diskUsage = 0
      }
    },

    // 关闭详情抽屉
    handleDetailClose() {
      this.detailDrawer = false
      // 销毁图表实例
      this.destroyGaugeCharts()
    },

    // 处理云主机导入成功
    handleCloudImportSuccess() {
      this.cloudDialogVisible = false
      this.getHostList()
      this.fetchHostStatistics()
    },

    // 处理Excel导入成功
    handleExcelImportSuccess() {
      this.getHostList()
      this.fetchHostStatistics()
    },

    // 处理JumpServer导入成功
    handleJumpServerImportSuccess() {
      this.$message.success('JumpServer同步任务已在后台执行,请稍后刷新列表查看结果')
      // 延迟5秒后自动刷新列表
      setTimeout(() => {
        this.getHostList()
        this.fetchHostStatistics()
      }, 5000)
    },

    // 连接SSH终端
    handleHostSSH() {
      if (!this.hostList.length) {
        this.$message.warning('请先选择主机')
        return
      }
      const selectedHost = this.hostList[0] // 默认选择第一个主机
      this.$router.push({
        path: '/cmdb/ssh',
        query: {
          hostId: selectedHost.id
        }
      })
    },

    // 跳转到资产授权页面
    handleHostAuthorization() {
      this.$router.push({
        path: '/system/post'
      })
    },

    // 跳转到主机详情页面
    handleHostDetail() {
      console.log('点击主机详情按钮')
      console.log('selectedHosts:', this.selectedHosts)
      console.log('hostList:', this.hostList)

      // 如果没有选中主机，但有主机列表，使用第一台主机
      if (!this.selectedHosts || this.selectedHosts.length === 0) {
        if (this.hostList && this.hostList.length > 0) {
          this.$router.push(`/cmdb/host/detail/${this.hostList[0].id}`)
          return
        } else {
          this.$message.warning('请先选择一台主机')
          return
        }
      }

      if (this.selectedHosts.length > 1) {
        this.$message.warning('请只选择一台主机')
        return
      }

      this.$router.push(`/cmdb/host/detail/${this.selectedHosts[0].id}`)
    },

    // 显示下载对话框
    showDownloadDialog(row) {
      this.$refs.fileManager.showDownloadDialog(row)
    },

    // 文件操作成功回调
    handleFileOperationSuccess() {
      // 可以在这里添加需要的刷新逻辑
      console.log('文件操作成功')
    },

    // 处理下拉框显示状态变化
    handleDropdownVisibleChange(visible) {
      console.log('下拉框显示状态变化:', visible)
      if (visible) {
        // 当下拉框即将显示时检查权限
        const hasPermission = this.checkPermission(['cmdb:ecs:add'])
        console.log('权限检查结果:', hasPermission)

        if (!hasPermission) {
          this.$message.warning('您没有新建主机的权限')
          // 阻止下拉框显示
          this.$nextTick(() => {
            this.$refs.createDropdown?.hide()
          })
          return false
        }
      }
    },

    // 处理新建按钮点击（现在只是一个占位符，真正的逻辑在visible-change中）
    handleCreateClick(event) {
      console.log('点击了新建按钮')
      // 权限检查会在 handleDropdownVisibleChange 中进行
    },

    // 处理下拉框选项点击
    handleCreateCommand(command) {
      console.log('选择了下拉框选项:', command)

      switch (command) {
        case 'importHost':
          this.addDialogVisible = true
          break
        case 'excelImport':
          this.ExcelDialogVisible = true
          break
        case 'jumpserverImport':
          this.jumpServerDialogVisible = true
          break
        case 'cloudHost':
          this.cloudDialogVisible = true
          break
      }
    },

    // 检查权限方法
    checkPermission(permissions) {
      console.log('checkPermission被调用，权限列表:', permissions)

      // 临时返回true用于测试
      console.log('权限检查通过')
      return true

      // TODO: 实现真正的权限检查逻辑
      // 假设您有全局的权限检查方法
      // if (this.$checkPermission) {
      //   return this.$checkPermission(permissions)
      // }

      // 或者检查store中的权限
      // if (this.$store && this.$store.getters.permissions) {
      //   const userPermissions = this.$store.getters.permissions
      //   return permissions.some(permission => userPermissions.includes(permission))
      // }
    },

    // 删除主机
    async handleHostDelete(row) {
      const confirmResult = await this.$confirm('是否确认删除主机"' + row.hostName + '"?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)

      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消删除')
      }

      const { data: res } = await this.$api.deleteCmdbHost(row.id)
      if (res.code === 200) {
        this.$message.success('删除成功')
        this.getHostList()
        this.fetchHostStatistics()
      }
    },

    // 处理表格选择变化
    handleHostSelectionChange(selection) {
      this.selectedHosts = selection
    },

    // 处理批量操作命令
    handleBatchCommand(command) {
      const hostTable = this.$refs.hostTable
      if (!this.selectedHosts || this.selectedHosts.length === 0) {
        this.$message.warning('请先选择要操作的主机')
        return
      }

      switch (command) {
        case 'updateGroup':
          // 调用子组件的批量修改分组方法
          hostTable.showBatchUpdateGroupDialog()
          break
        case 'updateVendor':
          this.$message.info('修改类型功能开发中...')
          // TODO: 实现批量修改类型功能
          break
        case 'bindTag':
          this.$message.info('绑定标签功能开发中...')
          // TODO: 实现批量绑定标签功能
          break
        case 'batchDelete':
          this.handleBatchDelete()
          break
        default:
          break
      }
    },

    // 批量删除主机
    async handleBatchDelete() {
      if (!this.selectedHosts || this.selectedHosts.length === 0) {
        this.$message.warning('请先选择要删除的主机')
        return
      }

      const hostNames = this.selectedHosts.map(h => h.hostName).join('、')
      const confirmResult = await this.$confirm(
        `是否确认删除以下 ${this.selectedHosts.length} 台主机？\n\n${hostNames}`,
        '批量删除确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
          dangerouslyUseHTMLString: false
        }
      ).catch(err => err)

      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消删除')
      }

      try {
        // 提取所有主机ID
        const hostIds = this.selectedHosts.map(h => h.id)

        // 调用批量删除接口
        const { data: res } = await this.$api.batchDeleteCmdbHost(hostIds)

        if (res.code === 200) {
          const deletedCount = res.data?.deletedCount || hostIds.length
          this.$message.success(res.data?.message || `成功删除 ${deletedCount} 台主机`)

          // 清除选择
          const hostTable = this.$refs.hostTable
          if (hostTable) {
            hostTable.clearSelection()
          }

          // 刷新列表
          this.getHostList()
          this.fetchHostStatistics()
        } else {
          this.$message.error(res.message || '批量删除失败')
        }
      } catch (error) {
        console.error('批量删除异常:', error)
        this.$message.error('批量删除操作失败')
      }
    },

    // 同步主机配置信息
    async handleHostSync() {
      try {
        // 检查是否有当前主机详情
        if (!this.hostDetail || !this.hostDetail.id) {
          this.$message.warning('请先选择要同步的主机')
          return
        }
        
        this.syncLoading = true
        const { data: res } = await this.$api.syncHostConfig(this.hostDetail.id)
        
        if (res.code === 200) {
          this.$message.success(res.data?.message || '开始同步主机信息，请稍后查看结果')
          // 同步成功后刷新主机列表，延迟3秒后刷新以获取最新数据
          setTimeout(() => {
            this.getHostList()
            // 如果详情面板打开，也刷新详情信息
            if (this.detailDrawer && this.hostDetail.id) {
              this.showHostDetail({ id: this.hostDetail.id })
            }
          }, 3000)
        } else {
          this.$message.error(res.message || '同步失败')
        }
      } catch (error) {
        console.error('同步主机配置失败:', error)
        this.$message.error('同步失败: ' + (error.message || '未知错误'))
      } finally {
        this.syncLoading = false
      }
    },

    // 分组管理 - 创建分组
    async handleCreateGroup(groupData) {
      try {
        console.log('创建分组数据:', groupData)
        const { data: res } = await cmdbApi.createCmdbGroup(groupData)
        
        if (res.code === 200) {
          this.$message.success('创建分组成功')
          // 刷新分组列表
          await this.getAllGroups()
        } else {
          this.$message.error(res.message || '创建分组失败')
        }
      } catch (error) {
        console.error('创建分组失败:', error)
        this.$message.error('创建分组失败: ' + (error.response?.data?.message || error.message))
      }
    },

    // 分组管理 - 更新分组
    async handleUpdateGroup(groupData) {
      try {
        console.log('更新分组数据:', groupData)
        const { data: res } = await cmdbApi.updateCmdbGroup(groupData)
        
        if (res.code === 200) {
          this.$message.success('更新分组成功')
          // 刷新分组列表
          await this.getAllGroups()
        } else {
          this.$message.error(res.message || '更新分组失败')
        }
      } catch (error) {
        console.error('更新分组失败:', error)
        this.$message.error('更新分组失败: ' + (error.response?.data?.message || error.message))
      }
    },

    // 分组管理 - 删除分组
    async handleDeleteGroup(groupId) {
      try {
        console.log('删除分组ID:', groupId)
        const { data: res } = await cmdbApi.deleteCmdbGroup(groupId)

        if (res.code === 200) {
          this.$message.success('删除分组成功')
          // 刷新分组列表
          await this.getAllGroups()
          // 如果删除的是当前选中的分组，重置选择
          if (this.currentGroupId === groupId) {
            this.currentGroupId = null
            this.getHostList()
          }
        } else {
          this.$message.error(res.message || '删除分组失败')
        }
      } catch (error) {
        console.error('删除分组失败:', error)
        this.$message.error('删除分组失败: ' + (error.response?.data?.message || error.message))
      }
    },

    // 初始化仪表盘图表
    initGaugeCharts() {
      this.$nextTick(() => {
        // 初始化 CPU 仪表盘
        if (this.$refs.cpuGauge) {
          this.cpuChart = echarts.init(this.$refs.cpuGauge)
          this.cpuChart.setOption(this.getGaugeOption(this.hostDetail.cpuUsage || 0, 'CPU'))
        }

        // 初始化内存仪表盘
        if (this.$refs.memoryGauge) {
          this.memoryChart = echarts.init(this.$refs.memoryGauge)
          this.memoryChart.setOption(this.getGaugeOption(this.hostDetail.memoryUsage || 0, '内存'))
        }

        // 初始化磁盘仪表盘
        if (this.$refs.diskGauge) {
          this.diskChart = echarts.init(this.$refs.diskGauge)
          this.diskChart.setOption(this.getGaugeOption(this.hostDetail.diskUsage || 0, '磁盘'))
        }
      })
    },

    // 获取仪表盘配置
    getGaugeOption(value, name) {
      return {
        series: [
          {
            type: 'gauge',
            startAngle: 225,
            endAngle: -45,
            min: 0,
            max: 100,
            radius: '85%',
            center: ['50%', '60%'],
            splitNumber: 10,
            axisLine: {
              lineStyle: {
                width: 10,
                color: [
                  [0.3, '#67e0e3'],
                  [0.7, '#37a2da'],
                  [1, '#fd666d']
                ]
              }
            },
            pointer: {
              length: '60%',
              width: 4,
              itemStyle: {
                color: '#4169E1'
              }
            },
            axisTick: {
              distance: -10,
              length: 4,
              lineStyle: {
                color: '#333',
                width: 1
              }
            },
            splitLine: {
              distance: -10,
              length: 8,
              lineStyle: {
                color: '#333',
                width: 2
              }
            },
            axisLabel: {
              color: '#666',
              distance: 20,
              fontSize: 12,
              formatter: function(value) {
                return value
              }
            },
            detail: {
              valueAnimation: true,
              formatter: '{value}%',
              color: '#ff0000',
              fontSize: 20,
              fontWeight: 'bold',
              offsetCenter: [0, '70%']
            },
            title: {
              offsetCenter: [0, '50%'],
              fontSize: 14,
              color: '#333',
              fontWeight: 'bold'
            },
            data: [
              {
                value: value,
                name: name
              }
            ]
          }
        ]
      }
    },

    // 销毁图表实例
    destroyGaugeCharts() {
      if (this.cpuChart) {
        this.cpuChart.dispose()
        this.cpuChart = null
      }
      if (this.memoryChart) {
        this.memoryChart.dispose()
        this.memoryChart = null
      }
      if (this.diskChart) {
        this.diskChart.dispose()
        this.diskChart = null
      }
    }

  }
}
</script>

<style scoped>
/* 🎨 现代化科技感设计风格 - 仿照cmdbDB.vue */

.cmdb-host-management {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.host-card {
  background: transparent !important;
  backdrop-filter: none !important;
  border-radius: 0 !important;
  box-shadow: none !important;
  border: none !important;
  padding: 0 !important;
}

:deep(.host-card) {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
  padding: 0 !important;
}

:deep(.el-card) {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
  padding: 0 !important;
}

:deep(.el-card__body) {
  padding: 0 !important;
}

.host-management-container {
  display: flex;
  height: 100vh;
}

.group-tree-section {
  width: 250px;
  margin-right: 10px;
}

.host-table-section {
  flex: 1;
  overflow-x: auto;
  overflow-y: visible;
  min-width: 0; /* 允许flex容器压缩 */
}

/* 🔍 搜索区域样式 */
.search-section {
  margin-bottom: 20px;
  padding: 20px;
  background: rgba(103, 126, 234, 0.15);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.2);
}

/* 📊 内联统计信息样式 */
.inline-statistics {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.stat-item {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.stat-icon {
  width: 20px;
  height: 20px;
  object-fit: contain;
  flex-shrink: 0;
}

.stat-text {
  font-weight: 600;
  font-size: 14px;
  color: #2c3e50;
}

/* 🎯 操作按钮区域 */
.action-section {
  margin-top: 15px;
  margin-bottom: 20px;
  padding-left: 0;
}

.table-section {
  margin-bottom: 15px;
}

.pagination-section {
  text-align: right;
  margin-top: 20px;
}

.font-weight-bold {
  font-weight: bold;
}

.table-operation {
  display: flex;
  justify-content: space-around;
}

/* 完全移除表单分割线 */
.el-dialog .el-form-item {
  border-bottom: none !important;
  margin-bottom: 12px;
  padding-bottom: 0;
}

/* 移除行和列之间的分割线 */
.el-row {
  border-bottom: none !important;
}

.el-col {
  border-right: none !important;
  padding-right: 0 !important;
  margin-right: 0 !important;
}

/* 移除最后一个列的右边距 */
.el-col:last-child {
  padding-right: 0 !important;
  margin-right: 0 !important;
}

/* 🎨 按钮样式优化 */
.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 📝 输入框和选择器样式 */
.el-input :deep(.el-input__wrapper),
.el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
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

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #2c3e50;
}

/* 🔧 表单项样式 */
.search-section .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-section .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

/* 🖥️ 终端按钮深色样式 */
.terminal-dark-btn {
  background: #2c3e50 !important;
  border-color: #2c3e50 !important;
  color: white !important;
}

.terminal-dark-btn:hover {
  background: #34495e !important;
  border-color: #34495e !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(44, 62, 80, 0.4);
}

/* 🎯 抽屉内容区域样式 - 减少顶部间距 */
.el-drawer :deep(.el-drawer__body) {
  padding-top: 10px;
}

/* 🎯 仪表盘样式 */
.dashboard-section {
  margin: 0;
  padding: 0;
}

.gauge-container {
  display: flex;
  justify-content: space-around;
  align-items: center;
  gap: 20px;
  margin-bottom: 5px;
}

.gauge-item {
  flex: 1;
  height: 180px;
  min-width: 0;
}

/* 📁 文件浏览器样式 */
.breadcrumb-nav {
  margin-bottom: 16px;
  padding: 12px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.breadcrumb-nav :deep(.el-breadcrumb__item) {
  cursor: pointer;
  transition: all 0.3s;
}

.breadcrumb-nav :deep(.el-breadcrumb__item:hover) {
  color: #409eff;
}

.current-path {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.current-path .el-input {
  flex: 1;
}

.current-path .path-action-btn {
  flex-shrink: 0;
}

.current-path :deep(.el-input-group__prepend) {
  background-color: #ecf5ff;
  border-color: #409eff;
  color: #409eff;
  padding: 0 15px;
}

.current-path :deep(.el-input__wrapper) {
  border-left-color: #409eff;
}

.current-path .el-text {
  color: #606266;
  font-weight: 500;
}
</style>
