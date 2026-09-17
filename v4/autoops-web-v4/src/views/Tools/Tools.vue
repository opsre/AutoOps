<template>
  <div class="service-market">
    <el-card shadow="hover" class="market-card">
      <template #header>
        <div class="card-header">
          <span class="title">运维工具箱 - 服务市场</span>
          <el-tabs v-model="innerTab" class="inner-tabs">
            <el-tab-pane label="服务市场" name="market" />
            <el-tab-pane label="部署管理" name="deploy" />
          </el-tabs>
        </div>
      </template>

      <!-- 市场视图 -->
      <template v-if="innerTab === 'market'">
        <!-- 分类筛选 -->
        <div class="category-filter">
          <el-radio-group v-model="selectedCategory" @change="handleCategoryChange" size="small">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.id"
            >
              {{ cat.name }}
            </el-radio-button>
          </el-radio-group>
        </div>

        <!-- 服务列表 -->
        <div class="service-grid" v-loading="loading">
          <el-card
            v-for="service in displayServices"
            :key="service.id"
            class="service-card"
            shadow="hover"
            @click="openDeployDialog(service)"
          >
            <div class="service-icon">
              <img
                :src="getServiceIcon(service.id)"
                :alt="service.name"
                class="service-svg-icon"
              />
            </div>
            <h3>{{ service.name }}</h3>
            <p class="description">{{ service.description }}</p>
            <div class="versions">
              <el-tag
                v-for="version in service.versions.slice(0, 2)"
                :key="version.id"
                :type="version.recommended ? 'success' : ''"
                size="small"
              >
                {{ version.name }}
              </el-tag>
            </div>
          </el-card>
        </div>
      </template>

      <!-- 部署管理视图（嵌入模式） -->
      <template v-else>
        <div class="embedded-deploy-wrapper">
          <DeployManage :embedded="true" />
        </div>
      </template>
    </el-card>

    <!-- 部署对话框 -->
    <DeployDialog
      v-model="deployDialogVisible"
      :service="selectedService"
      @deployed="handleDeployed"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getServicesList } from '@/api/tool'
import { ElMessage } from 'element-plus'
import DeployDialog from './docker_tool/DeployDialog.vue'
import DeployManage from './docker_tool/DeployManage.vue'

const services = ref([])
const categories = ref([])
const selectedCategory = ref('')
const deployDialogVisible = ref(false)
const selectedService = ref(null)
const loading = ref(false)
const innerTab = ref('market')

// 获取服务图标 - 使用动态 require
const getServiceIcon = (serviceId) => {
  try {
    // 文件名映射（处理大小写）
    const fileNameMap = {
      'mysql': 'mysql',
      'redis': 'redis',
      'postgresql': 'PostgreSQL',
      'jenkins': 'Jenkins',
      'gitlab': 'gitlab',
      'grafana': 'grafana',
      'elasticsearch': 'Elasticsearch',
      'loki': 'loki',
      'prometheus': 'Prometheus',
      'elk': 'ELK',
      'n9e': 'n9e',
      'jumpserver': 'jumpserver',
      'nodejs': 'nodejs',
      'java': 'java',
      'golang': 'golang',
      'mongodb': 'mongodb',
      'fluentd': 'fluentd'
    }

    // 转换为小写查找
    const key = serviceId.toLowerCase()
    const fileName = fileNameMap[key]

    if (fileName) {
      // 使用 require 动态加载图标
      return require(`@/assets/image/${fileName}.svg`)
    }
  } catch (error) {
    console.warn(`加载图标失败: ${serviceId}`, error)
  }

  // 返回默认图标
  return require('@/assets/image/云主机服务器.svg')
}

// 过滤服务
const filteredServices = computed(() => {
  if (!selectedCategory.value) return services.value
  return services.value.filter(s => s.category === selectedCategory.value)
})

// 显示服务（全部分类显示所有，其他分类显示前10个）
const displayServices = computed(() => {
  if (!selectedCategory.value) {
    return filteredServices.value
  }
  return filteredServices.value.slice(0, 10)
})

// 加载服务列表
const loadServices = async () => {
  loading.value = true
  try {
    const res = await getServicesList()
    if (res.data?.code === 200) {
      services.value = res.data.data?.services || []
      categories.value = res.data.data?.categories || []
    } else {
      throw new Error(res.data?.message || '获取服务列表失败')
    }
  } catch (error) {
    console.error('加载服务列表失败:', error)
    ElMessage.error(`加载服务列表失败: ${error.message}`)
  } finally {
    loading.value = false
  }
}

// 打开部署对话框
const openDeployDialog = (service) => {
  selectedService.value = service
  deployDialogVisible.value = true
}

// 部署完成回调
const handleDeployed = (deployId) => {
  ElMessage.success('部署任务已创建')
}

const handleCategoryChange = () => {
  // 分类变化处理
}

onMounted(() => {
  loadServices()
})
</script>

<style scoped>
.service-market {
  padding: 5px;
  height: 80vh;
  background: #eef1f6;
}

.market-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  height: calc(80vh - 10px);
  display: flex;
  flex-direction: column;
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

.category-filter {
  margin-bottom: 20px;
  padding: 20px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

/* 内部切换（参考 k8s-storage PVC/PV 样式） */
.inner-tabs {
  margin: 0;
}

.inner-tabs :deep(.el-tabs__header) {
  margin: 0;
  border: none;
}

.inner-tabs :deep(.el-tabs__item) {
  font-weight: 500;
  color: #606266;
}

.inner-tabs :deep(.el-tabs__item.is-active) {
  color: #409EFF;
  font-weight: 600;
}

.embedded-deploy-wrapper {
  padding: 10px 20px 0;
}

.service-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
  padding: 20px 0;
  overflow-y: auto;
  max-height: calc(80vh - 250px);
}

.service-card {
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: center;
  border-radius: 12px;
}

.service-card :deep(.el-card__body) {
  padding: 16px !important;
}

.service-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(103, 126, 234, 0.3);
}

.service-icon {
  margin-bottom: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 48px;
}

.service-svg-icon {
  width: 48px;
  height: 48px;
  object-fit: contain;
}

.service-card h3 {
  margin: 8px 0;
  font-size: 15px;
  font-weight: 600;
  color: #2c3e50;
}

.description {
  color: #999;
  font-size: 13px;
  margin: 8px 0;
  min-height: 36px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.versions {
  display: flex;
  gap: 4px;
  justify-content: center;
  flex-wrap: wrap;
  margin-top: 8px;
}

.versions .el-tag {
  font-size: 12px;
  padding: 2px 6px;
}

.el-radio-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

:deep(.el-radio-button__inner) {
  border-radius: 8px !important;
  transition: all 0.3s ease;
}

:deep(.el-radio-button__inner:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style>
