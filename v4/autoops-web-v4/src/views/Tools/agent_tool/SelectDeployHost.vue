<template>
  <el-dialog 
    title="选择部署Agent的主机" 
    v-model="visible"
    width="70%"
    :modal="false"
  >
    <div style="display: flex; height: 400px">
      <!-- 左侧分组树 -->
      <div style="width: 250px; border-right: 1px solid #ebeef5; padding-right: 10px; overflow-y: auto">
        <div style="margin-bottom: 10px">
          <el-button 
            type="primary" 
            size="small" 
            @click="showAllHosts" 
            style="width: 100%"
            :class="{ 'active-group': !currentGroupId }"
          >
            <el-icon><House /></el-icon>
            所有主机
          </el-button>
        </div>
        <el-tree
          :data="groupsWithHosts"
          :props="{
            label: 'name',
            children: 'children'
          }"
          node-key="id"
          highlight-current
          @node-click="handleGroupClick"
        >
          <template v-slot="{ node, data }">
            <span>
              <el-icon v-if="node.expanded" style="margin-right: 5px"><FolderOpened /></el-icon>
              <el-icon v-else style="margin-right: 5px"><Folder /></el-icon>
              {{ node.label }}
            </span>
          </template>
        </el-tree>
      </div>
      
      <!-- 右侧主机表格 -->
      <div style="flex: 1; padding-left: 10px; overflow-y: auto">
        <!-- 搜索框 -->
        <div style="margin-bottom: 10px">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索主机名称或IP地址"
            clearable
            size="small"
            style="width: 300px"
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
        
        <el-table
          :data="filteredHosts"
          border
          style="width: 100%"
          ref="hostTable"
          v-loading="loading"
          @selection-change="handleSelectionChange"
        >
          <el-table-column
            type="selection"
            width="55"
            :selectable="row => !hasAgent(row)"
          />
          <el-table-column
            prop="name"
            label="主机名称"
            width="150"
          >
            <template #default="{row}">
              <div class="host-cell">
                <el-icon><Monitor /></el-icon>
                <span style="margin-left: 5px">{{ row.hostName || row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="privateIp"
            label="内网IP"
            width="140"
          >
            <template #default="{row}">
              <div class="host-cell">
                <el-icon><HomeFilled /></el-icon>
                <span style="margin-left: 5px">{{ row.privateIp }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="publicIp"
            label="外网IP"
            width="140"
          >
            <template #default="{row}">
              <div class="host-cell">
                <el-icon color="#409EFF"><Connection /></el-icon>
                <span style="margin-left: 5px">{{ row.publicIp || '-' }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="os"
            label="操作系统"
            width="180"
          />
          <el-table-column
            label="Agent状态"
            width="120"
          >
            <template #default="{row}">
              <el-tag 
                :type="getAgentStatusTagType(row)"
                size="small"
              >
                {{ getAgentStatusText(row) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            prop="remark"
            label="备注"
          />
        </el-table>
        <div class="pagination">
          <el-pagination
            @current-change="handlePageChange"
            @size-change="handleSizeChange"
            :current-page="pagination.pageNum"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="pagination.pageSize"
            layout="total, sizes, prev, pager, next, jumper"
            :total="pagination.total"
          />
        </div>
      </div>
    </div>
    <template #footer>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <div style="color: #606266; font-size: 14px">
          <span v-if="tempSelectedHosts.length > 0">
            已选择 {{ tempSelectedHosts.length }} 台主机
          </span>
          <span v-else>请选择主机</span>
        </div>
        <div>
          <el-button @click="close">取消</el-button>
          <el-button
            type="primary"
            @click="confirmSelection"
            :disabled="tempSelectedHosts.length === 0"
          >部署到选中主机</el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, onMounted, computed } from 'vue'
import { 
  Search, 
  Monitor, 
  House, 
  FolderOpened, 
  Folder, 
  HomeFilled, 
  Connection 
} from '@element-plus/icons-vue'
import cmdbAPI from '@/api/cmdb'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'hosts-selected'])

const visible = ref(false)
const loading = ref(false)
const groupsWithHosts = ref([])
const currentGroupHosts = ref([])
const tempSelectedHosts = ref([])
const searchKeyword = ref('')
const existingAgents = ref([])
const pagination = ref({
  pageNum: 1,
  pageSize: 10,
  total: 0
})
const currentGroupId = ref(null)

const filteredHosts = computed(() => {
  // 当有搜索关键词时，不使用本地过滤，而是使用API搜索
  // 这样computed只负责返回当前显示的数据
  return currentGroupHosts.value
})

const fetchGroupsWithHosts = async () => {
  try {
    const response = await cmdbAPI.getGroupListWithHosts()
    if (response.data?.code === 200) {
      groupsWithHosts.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取分组和主机列表失败:', error)
  }
}

const fetchExistingAgents = async () => {
  try {
    const response = await cmdbAPI.getAgentList({
      page: 1,
      pageSize: 1000
    })
    if (response.data?.code === 200) {
      existingAgents.value = response.data.data?.list || []
    }
  } catch (error) {
    console.error('获取Agent列表失败:', error)
  }
}

const showAllHosts = async () => {
  currentGroupId.value = null
  await getAllHosts()
}

const getAllHosts = async () => {
  loading.value = true
  try {
    const response = await cmdbAPI.getCmdbHostList({
      page: pagination.value.pageNum,
      pageSize: pagination.value.pageSize
    })
    
    if (response.data?.code === 200) {
      currentGroupHosts.value = (response.data.data?.list || []).map(host => ({
        id: host.id,
        name: host.hostName || host.name,
        hostName: host.hostName || host.name,
        privateIp: host.privateIp,
        publicIp: host.publicIp,
        os: host.os,
        remark: host.remark
      }))
      pagination.value.total = response.data.data?.total || 0
    }
  } catch (error) {
    console.error('获取主机列表失败:', error)
    currentGroupHosts.value = []
  } finally {
    loading.value = false
  }
}

const handleGroupClick = async (data) => {
  currentGroupId.value = data.id
  loading.value = true
  try {
    const response = await cmdbAPI.getCmdbHostsByGroupId(data.id, {
      page: pagination.value.pageNum,
      pageSize: pagination.value.pageSize
    })
    
    if (response.data?.code === 200) {
      // API返回的数据结构：当有主机时，data是数组；当无主机时，data是null
      const hostList = response.data.data || []
      currentGroupHosts.value = hostList.map(host => ({
        id: host.id,
        name: host.hostName || host.name,
        hostName: host.hostName || host.name,
        privateIp: host.privateIp,
        publicIp: host.publicIp,
        os: host.os,
        remark: host.remark
      }))
      // 分组查询通常不返回total，使用数组长度
      pagination.value.total = hostList.length
    }
  } catch (error) {
    console.error('获取主机列表失败:', error)
    currentGroupHosts.value = []
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.value.pageNum = page
  // 如果正在搜索，继续搜索
  if (searchKeyword.value && searchKeyword.value.trim()) {
    searchHostsByName()
  } else if (currentGroupId.value) {
    handleGroupClick({ id: currentGroupId.value })
  } else {
    getAllHosts()
  }
}

const handleSizeChange = (size) => {
  pagination.value.pageSize = size
  pagination.value.pageNum = 1 // 改变页面大小时重置到第一页
  // 如果正在搜索，继续搜索
  if (searchKeyword.value && searchKeyword.value.trim()) {
    searchHostsByName()
  } else if (currentGroupId.value) {
    handleGroupClick({ id: currentGroupId.value })
  } else {
    getAllHosts()
  }
}

// 处理多选变化
const handleSelectionChange = (selection) => {
  tempSelectedHosts.value = selection
}

const handleSearch = async () => {
  // 搜索时重置到第一页
  pagination.value.pageNum = 1

  // 如果有搜索关键词，使用API搜索
  if (searchKeyword.value && searchKeyword.value.trim()) {
    await searchHostsByName()
  } else {
    // 如果搜索框为空，恢复显示当前分组或所有主机
    if (currentGroupId.value) {
      await handleGroupClick({ id: currentGroupId.value })
    } else {
      await getAllHosts()
    }
  }
}

const searchHostsByName = async () => {
  loading.value = true
  try {
    const response = await cmdbAPI.GetCmdbHostsByHostNameLike(searchKeyword.value.trim(), {
      pageNum: pagination.value.pageNum,
      pageSize: pagination.value.pageSize
    })

    if (response.data?.code === 200) {
      // 后端返回的data直接是数组，不是{list, total}结构
      const hosts = Array.isArray(response.data.data) ? response.data.data : []
      currentGroupHosts.value = hosts.map(host => ({
        id: host.id,
        name: host.hostName || host.name,
        hostName: host.hostName || host.name,
        privateIp: host.privateIp,
        publicIp: host.publicIp,
        os: host.os,
        remark: host.remark
      }))
      // 由于返回的是数组，使用数组长度作为total
      pagination.value.total = hosts.length
    }
  } catch (error) {
    console.error('搜索主机失败:', error)
    currentGroupHosts.value = []
  } finally {
    loading.value = false
  }
}

const hasAgent = (host) => {
  return existingAgents.value.some(agent => 
    agent.hostId === host.id && 
    (agent.status === 1 || agent.status === 3) // 部署中或运行中的不能再次部署
  )
}

const getAgentStatusText = (host) => {
  const agent = existingAgents.value.find(agent => agent.hostId === host.id)
  if (!agent) return '未部署'
  
  const statusTexts = {
    1: '部署中',
    2: '部署失败', 
    3: '运行中',
    4: '启动异常'
  }
  return statusTexts[agent.status] || '未知'
}

const getAgentStatusTagType = (host) => {
  const agent = existingAgents.value.find(agent => agent.hostId === host.id)
  if (!agent) return 'info'
  
  const statusTypes = {
    1: 'warning', // 部署中
    2: 'danger',  // 部署失败
    3: 'success', // 运行中
    4: 'danger'   // 启动异常
  }
  return statusTypes[agent.status] || 'info'
}

const confirmSelection = () => {
  emit('hosts-selected', tempSelectedHosts.value)
  tempSelectedHosts.value = []
}

const close = () => {
  visible.value = false
  tempSelectedHosts.value = []
  searchKeyword.value = ''
}

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    // 对话框打开时刷新数据
    fetchExistingAgents()
    getAllHosts()
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

onMounted(() => {
  fetchGroupsWithHosts()
  fetchExistingAgents()
})
</script>

<style scoped>
.host-cell {
  display: flex;
  align-items: center;
  padding: 5px;
  background-color: #f8f9fa;
  font-weight: 500;
  border-radius: 4px;
}

.pagination {
  margin-top: 15px;
  text-align: right;
}

.active-group {
  background-color: #409EFF;
  border-color: #409EFF;
  color: white;
}

.active-group:hover {
  background-color: #337ecc;
  border-color: #337ecc;
}

.el-tree {
  background: transparent;
}

.el-tree-node__content:hover {
  background-color: #f5f7fa;
}

.el-tree-node.is-current > .el-tree-node__content {
  background-color: #f0f9ff;
}
</style>