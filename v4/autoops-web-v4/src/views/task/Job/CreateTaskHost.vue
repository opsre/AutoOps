<template>
  <el-dialog 
    title="选择执行主机" 
    v-model="visible"
    width="70%"
  >
    <div style="display: flex; height: 500px">
      <!-- 左侧分组树 -->
      <div style="width: 250px; border-right: 1px solid #ebeef5; padding-right: 10px; overflow-y: auto">
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
        <el-table
          :data="currentGroupHosts"
          border
          style="width: 100%"
          @selection-change="handleHostSelectionChange"
          ref="hostTable"
        >
          <el-table-column
            type="selection"
            width="55"
            :selectable="() => true"
          />
          <el-table-column
            prop="name"
            label="主机名称"
            width="150"
          >
            <template #default="{row}">
              <div class="highlight-cell">
                <el-icon><User /></el-icon>
                <span style="margin-left: 5px">{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="privateIp"
            label="内网IP"
            width="120"
          >
            <template #default="{row}">
              <div class="highlight-cell">
                <el-icon><HomeFilled /></el-icon>
                <span style="margin-left: 5px">{{ row.privateIp }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="publicIp"
            label="外网IP"
            width="120"
          >
            <template #default="{row}">
              <div class="highlight-cell">
                <el-icon color="#409EFF"><Connection /></el-icon>
                <span style="margin-left: 5px">{{ row.publicIp }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="os"
            label="操作系统"
            width="180"
          />
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
      <el-button @click="close">关闭</el-button>
      <el-button 
        type="primary" 
        @click="confirmSelection"
      >添加选中主机</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import cmdbAPI from '@/api/cmdb'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  selectedHosts: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'hosts-selected'])

const visible = ref(false)
const groupsWithHosts = ref([])
const currentGroupHosts = ref([])
const tempSelectedHosts = ref([])
const pagination = ref({
  pageNum: 1,
  pageSize: 10,
  total: 0
})
const currentGroupId = ref(null)

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

const getAllHosts = async () => {
  try {
    const response = await cmdbAPI.getCmdbHostList({
      page: pagination.value.pageNum,
      pageSize: pagination.value.pageSize
    })
    
    if (response.data?.code === 200) {
      currentGroupHosts.value = (response.data.data?.list || []).map(host => ({
        id: host.id,
        name: host.hostName || host.name,
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
  }
}

const handleGroupClick = async (data) => {
  currentGroupId.value = data.id
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
  }
}

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
  if (currentGroupId.value) {
    handleGroupClick({ id: currentGroupId.value })
  } else {
    getAllHosts()
  }
}

const handleHostSelectionChange = (selection) => {
  tempSelectedHosts.value = selection
}

const confirmSelection = () => {
  emit('hosts-selected', tempSelectedHosts.value)
  tempSelectedHosts.value = []
}

const close = () => {
  visible.value = false
}

watch(() => props.modelValue, (val) => {
  visible.value = val
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

onMounted(() => {
  fetchGroupsWithHosts()
  getAllHosts()
})
</script>

<style scoped>
.highlight-cell {
  display: flex;
  align-items: center;
  padding: 5px;
  background-color: #fff8e6;
  font-weight: bold;
  border-radius: 4px;
}

.pagination {
  margin-top: 15px;
  text-align: right;
}
</style>
