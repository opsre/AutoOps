<template>
  <el-dialog 
    :title="dialogTitle" 
    v-model="visible"
    width="70%"
    :before-close="handleClose"
  >
    <div style="display: flex; height: 500px">
      <!-- 左侧分组树 -->
      <div style="width: 250px; border-right: 1px solid #ebeef5; padding-right: 10px; overflow-y: auto">
        <div class="group-header">
          <el-button 
            type="primary" 
            size="small" 
            @click="loadAllHosts"
            :loading="loading"
          >
            显示全部主机
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
          :loading="groupsLoading"
        >
          <template v-slot="{ node, data }">
            <span class="tree-node">
              <el-icon v-if="node.expanded" style="margin-right: 5px"><FolderOpened /></el-icon>
              <el-icon v-else style="margin-right: 5px"><Folder /></el-icon>
              {{ node.label }}
              <span class="host-count" v-if="data.hostCount !== undefined">({{ data.hostCount }})</span>
            </span>
          </template>
        </el-tree>
      </div>
      
      <!-- 右侧主机表格 -->
      <div style="flex: 1; padding-left: 10px; overflow-y: auto">
        <div class="table-header">
          <div class="selection-info">
            <span v-if="tempSelectedHosts.length > 0">
              已选择 <strong>{{ tempSelectedHosts.length }}</strong> 台主机
            </span>
            <span v-else class="text-muted">请从下方表格中选择主机</span>
          </div>
          <div class="table-actions">
            <el-button 
              size="small" 
              @click="clearSelection"
              :disabled="tempSelectedHosts.length === 0"
            >
              清空选择
            </el-button>
            <el-button 
              size="small" 
              type="primary" 
              @click="selectAll"
              :disabled="currentGroupHosts.length === 0"
            >
              全选当前页
            </el-button>
          </div>
        </div>
        
        <el-table
          :data="currentGroupHosts"
          v-loading="hostsLoading"
          border
          stripe
          style="width: 100%"
          @selection-change="handleHostSelectionChange"
          ref="hostTableRef"
        >
          <el-table-column
            type="selection"
            width="55"
            :selectable="isHostSelectable"
          />
          <el-table-column
            prop="name"
            label="主机名称"
            width="150"
            show-overflow-tooltip
          >
            <template #default="{row}">
              <div class="host-cell">
                <el-icon class="host-icon"><Monitor /></el-icon>
                <span>{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="ip"
            label="IP地址"
            width="120"
            show-overflow-tooltip
          >
            <template #default="{row}">
              <div class="host-cell">
                <el-icon class="ip-icon"><Connection /></el-icon>
                <span>{{ row.ip || row.privateIp }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="publicIp"
            label="公网IP"
            width="120"
            show-overflow-tooltip
          >
            <template #default="{row}">
              <div class="host-cell" v-if="row.publicIp">
                <el-icon class="public-ip-icon"><Globe /></el-icon>
                <span>{{ row.publicIp }}</span>
              </div>
              <span v-else class="text-muted">-</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="os"
            label="操作系统"
            width="120"
            show-overflow-tooltip
          />
          <el-table-column
            prop="status"
            label="状态"
            width="80"
          >
            <template #default="{row}">
              <el-tag 
                :type="getStatusType(row.status)" 
                size="small"
              >
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            prop="remark"
            label="备注"
            show-overflow-tooltip
          />
        </el-table>
        
        <div class="pagination" v-if="pagination.total > 0">
          <el-pagination
            @current-change="handlePageChange"
            @size-change="handleSizeChange"
            :current-page="pagination.pageNum"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="pagination.pageSize"
            layout="total, sizes, prev, pager, next, jumper"
            :total="pagination.total"
            small
          />
        </div>
      </div>
    </div>
    
    <template #footer>
      <div class="dialog-footer">
        <div class="footer-info">
          <span v-if="tempSelectedHosts.length > 0">
            将添加 <strong>{{ tempSelectedHosts.length }}</strong> 台主机到{{ nodeType }}节点
          </span>
        </div>
        <div class="footer-actions">
          <el-button @click="handleClose">取消</el-button>
          <el-button 
            type="primary" 
            @click="confirmSelection"
            :disabled="tempSelectedHosts.length === 0"
          >
            确认选择 ({{ tempSelectedHosts.length }})
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  FolderOpened, 
  Folder, 
  Monitor, 
  Connection, 
  Globe 
} from '@element-plus/icons-vue'
import cmdbApi from '@/api/cmdb'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  selectedHosts: {
    type: Array,
    default: () => []
  },
  nodeType: {
    type: String,
    default: 'Master'
  },
  multiple: {
    type: Boolean,
    default: true
  },
  excludeHosts: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'hosts-selected'])

const visible = ref(false)
const loading = ref(false)
const groupsLoading = ref(false)
const hostsLoading = ref(false)
const groupsWithHosts = ref([])
const currentGroupHosts = ref([])
const tempSelectedHosts = ref([])
const hostTableRef = ref()

const pagination = ref({
  pageNum: 1,
  pageSize: 20,
  total: 0
})

const currentGroupId = ref(null)

const dialogTitle = ref('')

// 计算对话框标题
watch(() => props.nodeType, (nodeType) => {
  dialogTitle.value = `选择${nodeType}节点主机`
}, { immediate: true })

// 获取分组列表
const fetchGroupsWithHosts = async () => {
  try {
    groupsLoading.value = true
    const response = await cmdbApi.getGroupListWithHosts()
    if (response.data?.code === 200) {
      groupsWithHosts.value = response.data.data?.map(group => ({
        id: group.id,
        name: group.name,
        hostCount: group.hostCount || 0
      })) || []
    }
  } catch (error) {
    console.error('获取分组列表失败:', error)
    ElMessage.error('获取分组列表失败')
  } finally {
    groupsLoading.value = false
  }
}

// 获取全部主机
const loadAllHosts = async () => {
  currentGroupId.value = null
  await getAllHosts()
}

const getAllHosts = async () => {
  try {
    hostsLoading.value = true
    const response = await cmdbApi.getCmdbHostList({
      page: pagination.value.pageNum,
      pageSize: pagination.value.pageSize
    })
    
    if (response.data?.code === 200) {
      const hostList = response.data.data?.list || []
      currentGroupHosts.value = hostList.map(host => ({
        id: host.id,
        name: host.hostName || host.name,
        ip: host.ip || host.privateIp,
        privateIp: host.privateIp,
        publicIp: host.publicIp,
        os: host.os || host.osType,
        status: host.status,
        remark: host.remark || host.description
      }))
      pagination.value.total = response.data.data?.total || 0
    }
  } catch (error) {
    console.error('获取主机列表失败:', error)
    ElMessage.error('获取主机列表失败')
    currentGroupHosts.value = []
  } finally {
    hostsLoading.value = false
  }
}

// 按分组获取主机
const handleGroupClick = async (data) => {
  currentGroupId.value = data.id
  pagination.value.pageNum = 1
  
  try {
    hostsLoading.value = true
    const response = await cmdbApi.getCmdbHostsByGroupId(data.id, {
      page: pagination.value.pageNum,
      pageSize: pagination.value.pageSize
    })
    
    if (response.data?.code === 200) {
      const hostList = response.data.data || []
      currentGroupHosts.value = hostList.map(host => ({
        id: host.id,
        name: host.hostName || host.name,
        ip: host.ip || host.privateIp,
        privateIp: host.privateIp,
        publicIp: host.publicIp,
        os: host.os || host.osType,
        status: host.status,
        remark: host.remark || host.description
      }))
      pagination.value.total = hostList.length
    }
  } catch (error) {
    console.error('获取分组主机失败:', error)
    ElMessage.error('获取分组主机失败')
    currentGroupHosts.value = []
  } finally {
    hostsLoading.value = false
  }
}

// 分页处理
const handlePageChange = (page) => {
  pagination.value.pageNum = page
  if (currentGroupId.value) {
    handleGroupClick({ id: currentGroupId.value })
  } else {
    getAllHosts()
  }
}

const handleSizeChange = (size) => {
  pagination.value.pageSize = size
  pagination.value.pageNum = 1
  if (currentGroupId.value) {
    handleGroupClick({ id: currentGroupId.value })
  } else {
    getAllHosts()
  }
}

// 主机选择处理
const handleHostSelectionChange = (selection) => {
  tempSelectedHosts.value = selection
}

const isHostSelectable = (row) => {
  // 检查是否在排除列表中
  return !props.excludeHosts.includes(row.id)
}

const clearSelection = () => {
  hostTableRef.value?.clearSelection()
  tempSelectedHosts.value = []
}

const selectAll = () => {
  nextTick(() => {
    currentGroupHosts.value.forEach(host => {
      if (isHostSelectable(host)) {
        hostTableRef.value?.toggleRowSelection(host, true)
      }
    })
  })
}

// 状态相关
const getStatusType = (status) => {
  const statusMap = {
    1: 'success',  // 在线
    0: 'danger',   // 离线
    2: 'warning'   // 维护中
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    1: '在线',
    0: '离线',
    2: '维护中'
  }
  return statusMap[status] || '未知'
}

// 确认选择
const confirmSelection = () => {
  if (tempSelectedHosts.value.length === 0) {
    ElMessage.warning('请选择至少一台主机')
    return
  }
  
  emit('hosts-selected', tempSelectedHosts.value)
  handleClose()
}

// 关闭对话框
const handleClose = () => {
  tempSelectedHosts.value = []
  clearSelection()
  visible.value = false
}

// 监听器
watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    // 对话框打开时重新加载数据
    fetchGroupsWithHosts()
    getAllHosts()
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

onMounted(() => {
  fetchGroupsWithHosts()
})
</script>

<style scoped>
.group-header {
  margin-bottom: 10px;
}

.tree-node {
  display: flex;
  align-items: center;
  width: 100%;
}

.host-count {
  color: #909399;
  font-size: 12px;
  margin-left: 5px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding: 8px 0;
  border-bottom: 1px solid #ebeef5;
}

.selection-info {
  font-size: 14px;
}

.text-muted {
  color: #909399;
}

.table-actions {
  display: flex;
  gap: 8px;
}

.host-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

.host-icon {
  color: #409eff;
}

.ip-icon {
  color: #67c23a;
}

.public-ip-icon {
  color: #e6a23c;
}

.pagination {
  margin-top: 15px;
  text-align: right;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.footer-info {
  flex: 1;
  font-size: 14px;
  color: #606266;
}

.footer-actions {
  display: flex;
  gap: 8px;
}

.el-table {
  border-radius: 8px;
}

:deep(.el-table__header-wrapper) {
  border-radius: 8px 8px 0 0;
}

:deep(.el-table__header) {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

:deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 600;
  border-bottom: 2px solid #e4e7ed;
}

:deep(.el-table__row:hover) {
  background-color: #f0f9ff !important;
}

:deep(.el-tree-node__content) {
  height: 32px;
  border-radius: 4px;
  margin: 2px 0;
}

:deep(.el-tree-node__content:hover) {
  background-color: #f0f9ff;
}

:deep(.el-tree-node.is-current .el-tree-node__content) {
  background-color: #e1f3d8;
  color: #67c23a;
  font-weight: 600;
}
</style>