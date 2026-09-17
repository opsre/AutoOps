<script setup>
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Monitor,
  Folder,
  Setting,
  Refresh
} from '@element-plus/icons-vue'
import appApi from '@/api/app'
import cmdbApi from '@/api/cmdb'
import databaseApi from '@/api/database'

// 导入自定义SVG图标
import parentGroupIcon from '@/assets/image/业务组.svg'
import childGroupIcon from '@/assets/image/子业务.svg'
import appIcon from '@/assets/image/应用.svg'

// 基础状态
const loading = ref(false)
const treeLoading = ref(false)
const detailLoading = ref(false)

// 业务树数据
const businessTreeData = ref([])
const defaultExpandedKeys = ref([])

// 当前选中的应用环境
const currentEnv = ref(null)
const currentApplication = ref(null)

// 详情页面切换
const activeDetailTab = ref('basic-info')

// 用于显示名称的选项数据（仍然需要用于详情面板的ID映射）
const hostOptions = ref([])
const databaseOptions = ref([])

// 主机列表详细数据
const hostListData = ref([])
const hostListLoading = ref(false)

// 发布历史数据
const releaseHistoryData = ref([])
const releaseHistoryLoading = ref(false)
const releaseHistoryTotal = ref(0)
const releaseHistoryPage = ref(1)
const releaseHistoryPageSize = ref(10)
const selectedEnvironment = ref('') // 环境筛选

// 获取业务树数据（从后端获取已构建好的树结构）
const fetchTreeData = async () => {
  treeLoading.value = true
  try {
    const response = await appApi.getApplicationViewTree()

    if (response.data?.code === 200) {
      let treeData = response.data.data || []
      // 按 group_id 排序，确保顺序固定
      treeData = treeData.sort((a, b) => a.group_id - b.group_id)
      // 转换后端返回的树结构为前端el-tree组件需要的格式
      businessTreeData.value = transformTreeData(treeData)

      // 默认展开第一层
      if (businessTreeData.value.length > 0) {
        defaultExpandedKeys.value = [businessTreeData.value[0].id]
      }
    } else {
      ElMessage.error(response.data?.message || '获取业务树失败')
    }
  } catch (error) {
    console.error('获取业务树失败:', error)
    ElMessage.error('获取业务树失败')
  } finally {
    treeLoading.value = false
  }
}

// 转换后端树数据为前端需要的格式
const transformTreeData = (backendData) => {
  return backendData.map(group => {
    const parentNode = {
      id: `parent-${group.group_id}`,
      label: group.group_name,
      type: 'parent-group',
      children: []
    }

    // 处理子业务组
    if (group.sub_groups && group.sub_groups.length > 0) {
      parentNode.children = group.sub_groups.map(subGroup => {
        const childNode = {
          id: `group-${subGroup.sub_group_id}`,
          label: subGroup.sub_group_name,
          type: 'child-group',
          groupId: subGroup.sub_group_id,
          children: []
        }

        // 处理应用
        if (subGroup.applications && subGroup.applications.length > 0) {
          childNode.children = subGroup.applications.map(app => {
            const appNode = {
              id: `app-${app.app_id}`,
              label: app.app_name,
              type: 'app',
              appId: app.app_id,
              appData: {
                id: app.app_id,
                name: app.app_name,
                code: app.code || '',
                description: app.description || '',
                repo_url: app.repo_url || '',
                business_group_id: subGroup.sub_group_id,
                business_dept_id: app.business_dept_id,
                programming_lang: app.programming_lang || '',
                start_command: app.start_command || '',
                stop_command: app.stop_command || '',
                health_api: app.health_api || '',
                dev_owners: app.dev_owners || [],
                test_owners: app.test_owners || [],
                ops_owners: app.ops_owners || [],
                domains: app.domains || [],
                hosts: app.hosts || [],
                databases: app.databases || [],
                other_res: app.other_res || {},
                status: app.status,
                jenkins_envs: app.environments || []
              },
              children: []
            }

            // 处理环境
            if (app.environments && app.environments.length > 0) {
              appNode.children = app.environments.map(env => ({
                id: `env-${app.app_id}-${env.env_id}`,
                label: env.env_name,
                type: 'env',
                appId: app.app_id,
                envId: env.env_id,
                appData: appNode.appData,
                envData: {
                  id: env.env_id,
                  env_name: env.env_name,
                  jenkins_server_id: env.jenkins_server_id,
                  job_name: env.job_name || '',
                  is_configured: env.is_configured || false
                }
              }))
            }

            return appNode
          })
        }

        return childNode
      })
    }

    return parentNode
  })
}

// 获取主机列表
const fetchHostOptions = async () => {
  try {
    const response = await cmdbApi.getCmdbHostList({
      page: 1,
      pageSize: 100 // 设置合理的分页大小
    })
    if (response.data?.code === 200) {
      const hostList = response.data.data?.list || []
      console.log('主机列表原始数据:', hostList)
      hostOptions.value = hostList.map(host => ({
        id: host.id,
        name: host.name || host.hostName || `主机${host.id}`,
        ip: host.privateIp || host.publicIp || host.sshIp || '-'
      }))
      console.log('处理后的主机选项:', hostOptions.value)
      console.log('主机ID映射:', hostOptions.value.reduce((acc, h) => {
        acc[h.id] = h.name
        return acc
      }, {}))
    }
  } catch (error) {
    console.error('获取主机列表失败:', error)
  }
}

// 获取数据库列表
const fetchDatabaseOptions = async () => {
  try {
    const response = await databaseApi.getInstanceList({
      pageNum: 1,
      pageSize: 100 // 设置合理的分页大小
    })
    if (response.data?.code === 200) {
      const dbList = response.data.data?.list || []
      console.log('数据库列表原始数据:', dbList)
      databaseOptions.value = dbList.map(db => ({
        id: db.id,
        name: db.name || `数据库${db.id}`,
        type: db.dbType || db.subType || 'Unknown'
      }))
      console.log('处理后的数据库选项:', databaseOptions.value)
    }
  } catch (error) {
    console.error('获取数据库列表失败:', error)
  }
}

// 注: fetchApplicationList 和 buildBusinessTree 已移除
// 现在直接使用 fetchTreeData() 从后端获取已构建好的树结构

// 获取主机详细列表
const fetchHostDetailList = async (hostIds) => {
  if (!hostIds || hostIds.length === 0) {
    hostListData.value = []
    return
  }

  hostListLoading.value = true
  try {
    // 批量获取主机详细信息
    const hostPromises = hostIds.map(id => cmdbApi.getCmdbHostById(id))
    const responses = await Promise.all(hostPromises)

    hostListData.value = responses
      .filter(res => res.data?.code === 200)
      .map(res => res.data.data)
      .filter(host => host) // 过滤掉空数据

    console.log('主机详细列表:', hostListData.value)
  } catch (error) {
    console.error('获取主机详细信息失败:', error)
    ElMessage.error('获取主机详细信息失败')
    hostListData.value = []
  } finally {
    hostListLoading.value = false
  }
}

// 获取发布历史
const fetchReleaseHistory = async () => {
  if (!currentApplication.value?.app_id) {
    releaseHistoryData.value = []
    return
  }

  releaseHistoryLoading.value = true
  try {
    const params = {
      app_id: currentApplication.value.app_id,
      page: releaseHistoryPage.value,
      pageSize: releaseHistoryPageSize.value
    }

    // 如果选择了环境，添加环境筛选
    if (selectedEnvironment.value) {
      params.environment = selectedEnvironment.value
    }

    const response = await appApi.getApplicationReleaseHistory(params)

    if (response.data?.code === 200) {
      releaseHistoryData.value = response.data.data?.list || []
      releaseHistoryTotal.value = response.data.data?.total || 0
      console.log('发布历史列表:', releaseHistoryData.value)
    } else {
      ElMessage.error(response.data?.message || '获取发布历史失败')
    }
  } catch (error) {
    console.error('获取发布历史失败:', error)
    ElMessage.error('获取发布历史失败')
    releaseHistoryData.value = []
  } finally {
    releaseHistoryLoading.value = false
  }
}

// 处理分页变化
const handleReleasePageChange = (page) => {
  releaseHistoryPage.value = page
  fetchReleaseHistory()
}

// 处理每页数量变化
const handleReleasePageSizeChange = (pageSize) => {
  releaseHistoryPageSize.value = pageSize
  releaseHistoryPage.value = 1
  fetchReleaseHistory()
}

// 处理环境筛选变化
const handleEnvironmentChange = () => {
  releaseHistoryPage.value = 1
  fetchReleaseHistory()
}

// 格式化发布时长
const formatDuration = (seconds) => {
  if (!seconds) return '-'
  const minutes = Math.floor(seconds / 60)
  const secs = seconds % 60
  return minutes > 0 ? `${minutes}分${secs}秒` : `${secs}秒`
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const statusMap = {
    1: 'info',      // 待审批
    2: 'warning',   // 审批通过
    3: 'success',   // 发布成功
    4: 'danger',    // 发布失败
    5: 'info'       // 已取消
  }
  return statusMap[status] || 'info'
}

// 树节点点击事件
const handleNodeClick = async (data) => {
  console.log('点击节点:', data)

  if (data.type === 'env' || data.type === 'app') {
    // 点击应用或环境节点，从后端获取完整的应用详情
    detailLoading.value = true
    try {
      const response = await appApi.getApplicationDetail(data.appId)

      if (response.data?.code === 200) {
        // 使用后端返回的完整详情数据
        currentApplication.value = response.data.data

        // 如果点击的是环境节点，设置当前环境信息
        if (data.type === 'env') {
          // 从完整详情中找到对应的环境配置
          const envDetail = response.data.data.environments?.find(env => env.env_id === data.envId)
          currentEnv.value = envDetail || data.envData
        } else {
          currentEnv.value = null
        }

        // 获取主机详细列表
        const hostIds = response.data.data.hosts || []
        await fetchHostDetailList(hostIds)

        // 获取发布历史
        releaseHistoryPage.value = 1
        selectedEnvironment.value = ''
        await fetchReleaseHistory()
      } else {
        ElMessage.error(response.data?.message || '获取应用详情失败')
      }
    } catch (error) {
      console.error('获取应用详情失败:', error)
      ElMessage.error('获取应用详情失败')
    } finally {
      detailLoading.value = false
    }
  } else {
    // 点击业务组节点,清空详情
    currentEnv.value = null
    currentApplication.value = null
    hostListData.value = []
    releaseHistoryData.value = []
  }
}

// 获取主机名称列表
const getHostNamesByIds = (hostIds) => {
  console.log('getHostNamesByIds 被调用，参数:', hostIds)
  console.log('当前 hostOptions.value:', hostOptions.value)

  if (!hostIds || !Array.isArray(hostIds) || hostIds.length === 0) return '-'

  const hostNames = hostIds.map(id => {
    console.log(`查找主机ID: ${id}`)
    const host = hostOptions.value.find(h => h.id === id)
    console.log(`找到的主机:`, host)
    return host ? `${host.name} (${host.ip})` : `主机ID: ${id}`
  })

  const result = hostNames.join(', ')
  console.log('最终返回结果:', result)
  return result
}

// 获取数据库名称列表
const getDatabaseNamesByIds = (dbIds) => {
  if (!dbIds || !Array.isArray(dbIds) || dbIds.length === 0) return '-'
  const dbNames = dbIds.map(id => {
    const db = databaseOptions.value.find(d => d.id === id)
    return db ? `${db.name} (${db.type})` : `数据库ID: ${id}`
  })
  return dbNames.join(', ')
}

// 格式化数组显示
const formatArrayDisplay = (arr) => {
  if (!arr || !Array.isArray(arr) || arr.length === 0) return '-'
  return arr.join(', ')
}

// 检查是否有其他资源需要显示
const hasOtherResources = computed(() => {
  if (!currentApplication.value?.other_res) return false
  const res = currentApplication.value.other_res
  return (res.kafka && res.kafka.length > 0) ||
         (res.rabbitmq && res.rabbitmq.length > 0) ||
         (res.zookeeper && res.zookeeper.length > 0) ||
         (res.other && res.other.length > 0)
})

// 刷新数据
const handleRefresh = async () => {
  loading.value = true
  try {
    await Promise.all([
      fetchTreeData(),
      fetchHostOptions(),
      fetchDatabaseOptions()
    ])
    ElMessage.success('刷新成功')
  } catch (error) {
    console.error('刷新失败:', error)
    ElMessage.error('刷新失败')
  } finally {
    loading.value = false
  }
}

// 页面初始化
onMounted(async () => {
  loading.value = true
  try {
    // 先加载辅助数据，用于详情面板的ID映射（主机和数据库）
    await Promise.all([
      fetchHostOptions(),
      fetchDatabaseOptions()
    ])
    // 然后加载树数据
    await fetchTreeData()
  } catch (error) {
    console.error('初始化失败:', error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="application-view">
    <el-card shadow="hover" class="main-card">
      <!-- 左右布局容器 -->
      <div class="view-container">
        <!-- 左侧业务树 -->
        <div class="tree-section">
          <div class="tree-header">
            <el-icon class="header-icon"><Folder /></el-icon>
            <span class="tree-title">业务树</span>
          </div>

          <el-tree
            :data="businessTreeData"
            :props="{ children: 'children', label: 'label' }"
            :default-expanded-keys="defaultExpandedKeys"
            node-key="id"
            @node-click="handleNodeClick"
            v-loading="treeLoading"
            class="business-tree"
          >
            <template #default="{ node, data }">
              <div class="tree-node">
                <!-- 一级分组 - 使用业务组.svg -->
                <img v-if="data.type === 'parent-group'" :src="parentGroupIcon" class="node-icon-img parent-group-icon" />

                <!-- 二级分组 - 使用子业务.svg -->
                <img v-else-if="data.type === 'child-group'" :src="childGroupIcon" class="node-icon-img child-group-icon" />

                <!-- 三级应用 - 使用应用.svg -->
                <img v-else-if="data.type === 'app'" :src="appIcon" class="node-icon-img app-icon" />

                <!-- 四级环境 - 使用 Element Plus 图标 -->
                <el-icon v-else-if="data.type === 'env'" class="node-icon env-icon">
                  <Setting />
                </el-icon>

                <span class="node-label">{{ node.label }}</span>
              </div>
            </template>
          </el-tree>
        </div>

        <!-- 右侧详情区域 -->
        <div class="detail-section">
          <!-- 操作栏 -->
          <div class="action-section">
            <span class="section-label">应用详情</span>
            <el-button type="success" size="small" :icon="Refresh" @click="handleRefresh" :loading="loading">
              刷新
            </el-button>
          </div>

          <!-- 未选择时的提示 -->
          <div v-if="!currentApplication" class="empty-state">
            <el-icon class="empty-icon"><Monitor /></el-icon>
            <p class="empty-text">请在左侧业务树中选择应用或环境查看详情</p>
          </div>

          <!-- 应用详情 -->
          <div v-else class="detail-content" v-loading="detailLoading">
            <!-- 标签页切换 -->
            <el-tabs v-model="activeDetailTab" class="detail-tabs">
              <!-- 基本信息标签页 -->
              <el-tab-pane label="基本信息" name="basic-info">
                <!-- 基本信息 -->
                <div class="info-section">
                  <h3 class="section-title">基本信息</h3>
                  <el-descriptions :column="2" border>
                    <el-descriptions-item label="应用名称">{{ currentApplication.app_name }}</el-descriptions-item>
                    <el-descriptions-item label="业务线">{{ currentApplication.business_group_name || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="业务部门">{{ currentApplication.business_dept_name || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="环境" v-if="currentEnv">
                      <el-tag type="success" size="small">{{ currentEnv.env_name }}</el-tag>
                    </el-descriptions-item>
                    <el-descriptions-item label="应用描述" :span="2" v-if="currentApplication.description">
                      {{ currentApplication.description }}
                    </el-descriptions-item>
                    <el-descriptions-item label="仓库地址" :span="2" v-if="currentApplication.repo_url">
                      <el-link :href="currentApplication.repo_url" target="_blank" type="primary">
                        {{ currentApplication.repo_url }}
                      </el-link>
                    </el-descriptions-item>
                  </el-descriptions>
                </div>

                <!-- Jenkins环境信息 -->
                <div class="info-section" v-if="currentEnv">
                  <h3 class="section-title">Jenkins环境配置</h3>
                  <el-descriptions :column="2" border>
                    <el-descriptions-item label="环境名称">{{ currentEnv.env_name }}</el-descriptions-item>
                    <el-descriptions-item label="Jenkins服务器">{{ currentEnv.jenkins_server_name || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="任务名称" :span="2">{{ currentEnv.job_name || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="任务链接" :span="2" v-if="currentEnv.job_url">
                      <el-link :href="currentEnv.job_url" target="_blank" type="primary">
                        {{ currentEnv.job_url }}
                      </el-link>
                    </el-descriptions-item>
                  </el-descriptions>
                </div>

                <!-- 技术信息 -->
                <div class="info-section">
                  <h3 class="section-title">技术信息</h3>
                  <el-descriptions :column="2" border>
                    <el-descriptions-item label="编程语言">
                      <el-tag v-if="currentApplication.programming_lang" type="primary" size="small">
                        {{ currentApplication.programming_lang }}
                      </el-tag>
                      <span v-else>-</span>
                    </el-descriptions-item>
                    <el-descriptions-item label="健康检查">{{ currentApplication.health_api || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="启动命令">{{ currentApplication.start_command || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="停止命令">{{ currentApplication.stop_command || '-' }}</el-descriptions-item>
                  </el-descriptions>
                </div>

                <!-- 负责人信息 -->
                <div class="info-section">
                  <h3 class="section-title">负责人信息</h3>
                  <el-descriptions :column="1" border>
                    <el-descriptions-item label="研发负责人">{{ currentApplication.dev_owner_names?.join(', ') || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="测试负责人">{{ currentApplication.test_owner_names?.join(', ') || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="运维负责人">{{ currentApplication.ops_owner_names?.join(', ') || '-' }}</el-descriptions-item>
                  </el-descriptions>
                </div>

                <!-- 关联资源 -->
                <div class="info-section">
                  <h3 class="section-title">关联资源</h3>
                  <el-descriptions :column="1" border>
                    <el-descriptions-item label="关联域名">{{ formatArrayDisplay(currentApplication.domains) }}</el-descriptions-item>
                    <el-descriptions-item label="关联主机">{{ getHostNamesByIds(currentApplication.hosts) }}</el-descriptions-item>
                    <el-descriptions-item label="关联数据库">{{ getDatabaseNamesByIds(currentApplication.databases) }}</el-descriptions-item>
                  </el-descriptions>
                </div>

                <!-- 其他资源 -->
                <div class="info-section" v-if="currentApplication.other_res && hasOtherResources">
                  <h3 class="section-title">其他资源</h3>
                  <el-descriptions :column="2" border>
                    <el-descriptions-item label="Kafka" v-if="currentApplication.other_res.kafka">
                      {{ formatArrayDisplay(currentApplication.other_res.kafka) }}
                    </el-descriptions-item>
                    <el-descriptions-item label="RabbitMQ" v-if="currentApplication.other_res.rabbitmq">
                      {{ formatArrayDisplay(currentApplication.other_res.rabbitmq) }}
                    </el-descriptions-item>
                    <el-descriptions-item label="Kubernetes" v-if="currentApplication.other_res.zookeeper">
                      {{ formatArrayDisplay(currentApplication.other_res.zookeeper) }}
                    </el-descriptions-item>
                    <el-descriptions-item label="其他资源" v-if="currentApplication.other_res.other">
                      {{ formatArrayDisplay(currentApplication.other_res.other) }}
                    </el-descriptions-item>
                  </el-descriptions>
                </div>
              </el-tab-pane>

              <!-- 主机列表标签页 -->
              <el-tab-pane label="主机列表" name="host-list">
                <div v-if="!currentApplication?.hosts || currentApplication.hosts.length === 0" class="placeholder-content">
                  <el-icon class="placeholder-icon"><Monitor /></el-icon>
                  <p class="placeholder-text">暂无关联主机</p>
                  <p class="placeholder-desc">该应用当前未关联任何主机</p>
                </div>
                <div v-else v-loading="hostListLoading">
                  <el-table :data="hostListData" border stripe style="width: 100%">
                    <el-table-column prop="name" label="主机名称" min-width="150" />
                    <el-table-column label="IP地址" min-width="130">
                      <template #default="{ row }">
                        {{ row.privateIp || row.publicIp || row.sshIp || '-' }}
                      </template>
                    </el-table-column>
                    <el-table-column label="CPU" min-width="100">
                      <template #default="{ row }">
                        {{ row.cpu ? `${row.cpu} 核` : '-' }}
                      </template>
                    </el-table-column>
                    <el-table-column label="内存" min-width="100">
                      <template #default="{ row }">
                        {{ row.memory ? `${row.memory} GB` : '-' }}
                      </template>
                    </el-table-column>
                    <el-table-column label="磁盘" min-width="100">
                      <template #default="{ row }">
                        {{ row.disk ? `${row.disk} GB` : '-' }}
                      </template>
                    </el-table-column>
                    <el-table-column label="操作系统" min-width="120">
                      <template #default="{ row }">
                        {{ row.os || row.osType || '-' }}
                      </template>
                    </el-table-column>
                    <el-table-column prop="tag" label="标签" min-width="120" />
                    <el-table-column label="状态" min-width="80">
                      <template #default="{ row }">
                        <el-tag v-if="row.status === 1" type="success" size="small">在线</el-tag>
                        <el-tag v-else-if="row.status === 0" type="danger" size="small">离线</el-tag>
                        <el-tag v-else type="info" size="small">未知</el-tag>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </el-tab-pane>

              <!-- 发布详情标签页 -->
              <el-tab-pane label="发布详情" name="release-details">
                <div v-if="!currentApplication" class="placeholder-content">
                  <el-icon class="placeholder-icon"><Setting /></el-icon>
                  <p class="placeholder-text">请选择应用</p>
                  <p class="placeholder-desc">请在左侧业务树中选择应用查看发布历史</p>
                </div>
                <div v-else>

                  <!-- 发布历史表格 -->
                  <div v-loading="releaseHistoryLoading">
                    <el-table :data="releaseHistoryData" border stripe style="width: 100%">
                      <el-table-column label="上线时间" min-width="150">
                        <template #default="{ row }">
                          {{ row.created_at || '-' }}
                        </template>
                      </el-table-column>
                      <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
                      <el-table-column label="环境" min-width="80">
                        <template #default="{ row }">
                          <el-tag v-if="row.environment === 'prod'" type="danger" size="small">生产</el-tag>
                          <el-tag v-else-if="row.environment === 'test'" type="warning" size="small">测试</el-tag>
                          <el-tag v-else-if="row.environment === 'dev'" type="info" size="small">开发</el-tag>
                          <el-tag v-else size="small">{{ row.environment }}</el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column prop="applicant_name" label="申请人" min-width="100" />
                      <el-table-column label="状态" min-width="100">
                        <template #default="{ row }">
                          <el-tag :type="getStatusTagType(row.status)" size="small">
                            {{ row.status_text || '-' }}
                          </el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column prop="branch" label="分支" min-width="120" show-overflow-tooltip />
                      <el-table-column label="构建编号" min-width="100">
                        <template #default="{ row }">
                          #{{ row.build_number || '-' }}
                        </template>
                      </el-table-column>
                      <el-table-column label="耗时" min-width="100">
                        <template #default="{ row }">
                          {{ formatDuration(row.duration) }}
                        </template>
                      </el-table-column>
                      <el-table-column label="操作" min-width="100" fixed="right">
                        <template #default="{ row }">
                          <el-link v-if="row.jenkins_job_url" :href="row.jenkins_job_url" target="_blank" type="primary" :underline="false" style="font-size: 12px;">
                            查看日志
                          </el-link>
                          <span v-else style="color: #909399; font-size: 12px;">-</span>
                        </template>
                      </el-table-column>
                    </el-table>

                    <!-- 分页 -->
                    <div class="pagination-container">
                      <el-pagination
                        v-model:current-page="releaseHistoryPage"
                        v-model:page-size="releaseHistoryPageSize"
                        :page-sizes="[10, 20, 50, 100]"
                        :total="releaseHistoryTotal"
                        layout="total, sizes, prev, pager, next, jumper"
                        @size-change="handleReleasePageSizeChange"
                        @current-change="handleReleasePageChange"
                      />
                    </div>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
/* 🎨 现代化科技感设计风格 - 参考cmdbHost.vue */

.application-view {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

.main-card {
  background: transparent !important;
  backdrop-filter: none;
  border-radius: 16px;
  box-shadow: none !important;
  border: none !important;
}

.main-card :deep(.el-card__body) {
  padding: 5px;
  background: transparent !important;
  border: none !important;
}

.view-container {
  display: flex;
  gap: 16px;
  height: calc(100vh - 180px);
}

/* 左侧树面板 */
.tree-section {
  width: 220px;
  min-width: 200px;
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.1);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.tree-header {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  background: linear-gradient(135deg, rgba(103, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-bottom: 1px solid rgba(103, 126, 234, 0.1);
}

.header-icon {
  color: #409EFF;
  font-size: 14px;
}

.tree-title {
  font-size: 13px;
  font-weight: 600;
  color: #2c3e50;
}

.business-tree {
  flex: 1;
  overflow-y: auto;
  padding: 6px;
}

/* 树节点样式 - 紧凑模式 */
.business-tree :deep(.el-tree-node) {
  margin-bottom: 1px;
}

.business-tree :deep(.el-tree-node__content) {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 4px;
  margin: 0;
  padding: 2px 4px;
  height: 30px;
  line-height: 30px;
}

.business-tree :deep(.el-tree-node__content:hover) {
  background-color: rgba(103, 126, 234, 0.1);
  transform: translateX(4px);
}

.business-tree :deep(.el-tree-node.is-current > .el-tree-node__content) {
  background: linear-gradient(90deg, rgba(103, 126, 234, 0.2) 0%, rgba(103, 126, 234, 0.1) 100%);
  border-left: 3px solid #667eea;
}

/* 减少树节点的子节点间距 */
.business-tree :deep(.el-tree-node__children) {
  overflow: visible;
}

.business-tree :deep(.el-tree-node__children .el-tree-node) {
  padding-top: 0;
  padding-bottom: 0;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0;
  width: 100%;
  font-size: 13px;
}

/* 图片图标样式 */
.node-icon-img {
  width: 15px;
  height: 15px;
  flex-shrink: 0;
  transition: all 0.3s ease;
  object-fit: contain;
}

/* Element Plus 图标样式（仅用于环境节点） */
.node-icon {
  font-size: 15px;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.env-icon {
  color: #E6A23C;
}

.node-label {
  font-size: 13px;
  color: #606266;
  transition: color 0.3s ease;
  line-height: 1.3;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.business-tree :deep(.el-tree-node__content:hover .node-label) {
  color: #667eea;
}

.business-tree :deep(.el-tree-node__content:hover .node-icon) {
  transform: scale(1.1);
}

/* 右侧详情区域 */
.detail-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.action-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.1);
  margin-bottom: 12px;
}

.section-label {
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #909399;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.empty-icon {
  font-size: 64px;
  color: #DCDFE6;
  margin-bottom: 16px;
  opacity: 0.6;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.empty-text {
  font-size: 14px;
  color: #909399;
}

/* 详情内容 */
.detail-content {
  height: 570px;
  overflow-y: auto;
  padding: 12px;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

/* 标签页样式 */
.detail-tabs {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.detail-tabs :deep(.el-tabs__header) {
  margin-bottom: 12px;
  background: rgba(103, 126, 234, 0.05);
  padding: 6px 12px;
  border-radius: 6px;
}

.detail-tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.detail-tabs :deep(.el-tabs__item) {
  font-size: 13px;
  font-weight: 500;
  color: #606266;
  padding: 0 16px;
  transition: all 0.3s ease;
}

.detail-tabs :deep(.el-tabs__item:hover) {
  color: #667eea;
}

.detail-tabs :deep(.el-tabs__item.is-active) {
  color: #667eea;
  font-weight: 600;
}

.detail-tabs :deep(.el-tabs__active-bar) {
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  height: 2px;
}

.detail-tabs :deep(.el-tabs__content) {
  flex: 1;
  overflow-y: auto;
}

/* 占位内容样式 */
.placeholder-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
}

.placeholder-icon {
  font-size: 64px;
  color: #DCDFE6;
  margin-bottom: 16px;
  opacity: 0.6;
  animation: float 3s ease-in-out infinite;
}

.placeholder-text {
  font-size: 15px;
  font-weight: 600;
  color: #606266;
  margin-bottom: 10px;
}

.placeholder-desc {
  font-size: 13px;
  color: #909399;
  margin: 0;
}

.info-section {
  margin-bottom: 16px;
}

.info-section:last-child {
  margin-bottom: 0;
}

.section-title {
  margin: 0 0 10px 0;
  padding-bottom: 6px;
  font-size: 13px;
  font-weight: 600;
  color: #409eff;
  border-bottom: 2px solid rgba(64, 158, 255, 0.2);
}

/* el-descriptions 样式优化 */
.detail-content :deep(.el-descriptions) {
  background: rgba(255, 255, 255, 0.8);
}

.detail-content :deep(.el-descriptions__label) {
  font-weight: 600;
  color: #606266;
  background: rgba(248, 249, 250, 0.8);
  width: 100px;
  font-size: 13px;
}

.detail-content :deep(.el-descriptions__content) {
  color: #303133;
  font-size: 13px;
}

.detail-content :deep(.el-descriptions__body) {
  border-radius: 6px;
  overflow: hidden;
}

.detail-content :deep(.el-descriptions__cell) {
  padding: 8px 12px;
  transition: background-color 0.3s ease;
}

.detail-content :deep(.el-descriptions__cell:hover) {
  background-color: rgba(103, 126, 234, 0.05);
}

/* 按钮样式 */
.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.el-button:active {
  transform: translateY(0);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.el-button--primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

.el-button--primary:hover {
  background: linear-gradient(135deg, #7689ef 0%, #8659b5 100%);
}

/* 标签样式 */
.el-tag {
  font-weight: 500;
  border-radius: 4px;
  border: none;
  padding: 2px 8px;
  font-size: 12px;
  transition: all 0.3s ease;
}

.el-tag:hover {
  transform: scale(1.05);
}

.el-tag--success {
  background: linear-gradient(135deg, #67C23A 0%, #85ce61 100%);
  color: white;
}

.el-tag--primary {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  color: white;
}

/* 链接样式 */
.el-link {
  font-size: 13px;
  transition: all 0.3s ease;
}

.el-link:hover {
  transform: translateX(2px);
}

/* 发布历史样式 */
.release-filter {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 6px;
}

.filter-label {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
  margin-right: 12px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
  padding: 12px;
}

.pagination-container :deep(.el-pagination) {
  font-size: 13px;
}

.pagination-container :deep(.el-pagination__total) {
  font-size: 13px;
}

.pagination-container :deep(.el-select .el-input) {
  font-size: 13px;
}

/* 滚动条美化 */
.business-tree::-webkit-scrollbar,
.detail-content::-webkit-scrollbar {
  width: 6px;
}

.business-tree::-webkit-scrollbar-track,
.detail-content::-webkit-scrollbar-track {
  background: rgba(103, 126, 234, 0.05);
  border-radius: 3px;
}

.business-tree::-webkit-scrollbar-thumb,
.detail-content::-webkit-scrollbar-thumb {
  background: rgba(103, 126, 234, 0.3);
  border-radius: 3px;
  transition: background 0.3s ease;
}

.business-tree::-webkit-scrollbar-thumb:hover,
.detail-content::-webkit-scrollbar-thumb:hover {
  background: rgba(103, 126, 234, 0.5);
}

/* 加载遮罩样式 */
:deep(.el-loading-mask) {
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(4px);
}

:deep(.el-loading-spinner .el-loading-text) {
  color: #667eea;
  font-weight: 500;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .view-container {
    flex-direction: column;
    height: auto;
  }

  .tree-section {
    width: 100%;
    max-height: 400px;
  }
}

@media (max-width: 768px) {
  .application-view {
    padding: 12px;
  }

  .main-card {
    border-radius: 12px;
  }

  .section-label {
    font-size: 14px;
  }

  .detail-content :deep(.el-descriptions__label) {
    width: 100px;
  }

  .detail-content :deep(.el-descriptions__cell) {
    padding: 10px 12px;
  }
}

@media (max-width: 480px) {
  .application-view {
    padding: 8px;
  }

  .action-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}
</style>
