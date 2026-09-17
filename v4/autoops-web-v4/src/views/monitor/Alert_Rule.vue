<script setup>
import { ref, onMounted, reactive } from 'vue'
import {
  AlertCircle,
  SuccessFilled,
  Clock,
  Calendar
} from '@element-plus/icons-vue'
import { recordApi } from '@/api/alert'
import AlertTemplate from './alert/AlertTemplate.vue'
import AlertPolicy from './alert/AlertPolicy.vue'
import AlertHistory from './alert/AlertHistory.vue'

// 当前活动标签
const activeTab = ref('policy')

// 统计数据
const statistics = reactive({
  firing: 0,
  resolved: 0,
  acknowledged: 0,
  today: 0
})

const statsLoading = ref(false)

// 获取统计数据
const fetchStatistics = async () => {
  statsLoading.value = true
  try {
    const { data: res } = await recordApi.getStats()
    if (res.code === 200) {
      const data = res.data
      statistics.firing = data.firing || 0
      statistics.resolved = data.resolved || 0
      statistics.acknowledged = data.acknowledged || 0
      statistics.today = data.today || 0
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    statsLoading.value = false
  }
}

// 刷新统计方法（传递给子组件）
const handleRefreshStats = () => {
  fetchStatistics()
}

onMounted(() => {
  fetchStatistics()
})
</script>

<template>
  <div class="alert-rule-container">
    <!-- 统计卡片 -->
    <div class="statistics-cards">
      <!-- Firing告警 -->
      <div class="stat-card stat-firing">
        <div class="stat-icon">
          <el-icon :size="32"><AlertCircle /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">🔥 Firing告警</div>
          <div class="stat-value">{{ statistics.firing }}</div>
        </div>
      </div>

      <!-- 已处理 -->
      <div class="stat-card stat-resolved">
        <div class="stat-icon">
          <el-icon :size="32"><SuccessFilled /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">✅ 已处理</div>
          <div class="stat-value">{{ statistics.resolved }}</div>
        </div>
      </div>

      <!-- 待认领 -->
      <div class="stat-card stat-acknowledged">
        <div class="stat-icon">
          <el-icon :size="32"><Clock /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">⏳ 待认领</div>
          <div class="stat-value">{{ statistics.acknowledged }}</div>
        </div>
      </div>

      <!-- 今日告警 -->
      <div class="stat-card stat-today">
        <div class="stat-icon">
          <el-icon :size="32"><Calendar /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">📊 今日告警</div>
          <div class="stat-value">{{ statistics.today }}</div>
        </div>
      </div>
    </div>

    <!-- 告警管理卡片 -->
    <el-card shadow="hover" class="alert-card">
      <template #header>
        <div class="card-header">
          <span class="title">告警管理</span>
        </div>
      </template>

      <!-- 标签页 -->
      <el-tabs v-model="activeTab" class="alert-tabs">
        <!-- 告警策略配置标签页 -->
        <el-tab-pane name="policy">
          <template #label>
            <span class="tab-label">
              <el-icon>⚙️</el-icon>
              告警策略配置
            </span>
          </template>
          <AlertPolicy @refresh-stats="handleRefreshStats" />
        </el-tab-pane>
        <!-- 历史告警记录标签页 -->
        <el-tab-pane name="history">
          <template #label>
            <span class="tab-label">
              <el-icon>📜</el-icon>
              历史告警记录
            </span>
          </template>
          <AlertHistory @refresh-stats="handleRefreshStats" />
        </el-tab-pane>
        
        <!-- 告警规则模板标签页 -->
        <el-tab-pane name="template">
          <template #label>
            <span class="tab-label">
              <el-icon>📋</el-icon>
              告警规则模板
            </span>
          </template>
          <AlertTemplate @refresh-stats="handleRefreshStats" />
        </el-tab-pane>




      </el-tabs>
    </el-card>
  </div>
</template>

<style scoped>
.alert-rule-container {
  padding: 5px;
  min-height: 100vh;
  background: #eef1f6;
}

/* 统计卡片区域 */
.statistics-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 12px;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 6px;
  font-weight: 500;
}

.stat-value {
  font-size: 26px;
  font-weight: bold;
  line-height: 1;
}

/* 不同类型统计卡片颜色 */
.stat-card.stat-firing .stat-icon {
  background: linear-gradient(135deg, #F56C6C 0%, #F78989 100%);
  color: white;
}

.stat-card.stat-firing .stat-value {
  color: #F56C6C;
}

.stat-card.stat-resolved .stat-icon {
  background: linear-gradient(135deg, #67C23A 0%, #85CE61 100%);
  color: white;
}

.stat-card.stat-resolved .stat-value {
  color: #67C23A;
}

.stat-card.stat-acknowledged .stat-icon {
  background: linear-gradient(135deg, #E6A23C 0%, #EEBE77 100%);
  color: white;
}

.stat-card.stat-acknowledged .stat-value {
  color: #E6A23C;
}

.stat-card.stat-today .stat-icon {
  background: linear-gradient(135deg, #409EFF 0%, #66B3FF 100%);
  color: white;
}

.stat-card.stat-today .stat-value {
  color: #409EFF;
}

/* 主卡片 */
.alert-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  min-height: calc(100vh - 100px);
}

.alert-card :deep(.el-card__body) {
  min-height: 500px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  color: transparent;
}

/* 标签页样式 */
.alert-tabs {
  margin-top: 0;
}

.alert-tabs :deep(.el-tabs__header) {
  margin-bottom: 20px;
}

.alert-tabs :deep(.el-tabs__item) {
  font-weight: 500;
  color: #606266;
  padding: 0 20px;
}

.alert-tabs :deep(.el-tabs__item.is-active) {
  color: #409EFF;
  font-weight: 600;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
}

.tab-label .el-icon {
  font-size: 16px;
}

/* 响应式调整 */
@media (max-width: 1400px) {
  .statistics-cards {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }
}

@media (max-width: 768px) {
  .alert-rule-container {
    padding: 10px;
  }

  .statistics-cards {
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
    margin-bottom: 15px;
  }

  .stat-card {
    padding: 12px;
    gap: 10px;
  }

  .stat-icon {
    width: 42px;
    height: 42px;
  }

  .stat-icon :deep(.el-icon) {
    font-size: 22px !important;
  }

  .stat-label {
    font-size: 12px;
  }

  .stat-value {
    font-size: 20px;
  }
}
</style>
