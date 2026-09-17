<template>
  <el-select
    :model-value="modelValue"
    placeholder="请选择集群"
    :size="size"
    :style="{ width, marginRight }"
    :class="['cluster-selector', theme ? `theme-${theme}` : '']"
    @change="handleChange"
    :loading="!clusterList.length"
  >
    <template #prefix v-if="theme === 'cluster'">
      <el-icon class="selector-prefix-icon cluster-prefix-icon"><Connection /></el-icon>
    </template>
    <el-option
      v-for="cluster in clusterList"
      :key="cluster.id"
      :label="`${cluster.name} ${getClusterStatusText(cluster.status)}`"
      :value="cluster.id"
      :disabled="cluster.status !== 2"
    >
      <div class="cluster-option">
        <span class="cluster-name">{{ cluster.name }}</span>
        <el-tag
          :type="getClusterStatusTag(cluster.status)"
          size="small"
          class="cluster-status-tag"
        >
          {{ getClusterStatusText(cluster.status) }}
        </el-tag>
      </div>
    </el-option>
  </el-select>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Connection } from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import {
  getStoredClusterId,
  saveStoredClusterId,
  resolvePreferredClusterId
} from '@/utils/k8sClusterSelection'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  width: {
    type: String,
    default: '200px'
  },
  size: {
    type: String,
    default: 'small'
  },
  marginRight: {
    type: String,
    default: '12px'
  },
  theme: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const clusterList = ref([])

// 获取集群列表
const fetchClusterList = async () => {
  try {
    const response = await k8sApi.getClusterList()
    const responseData = response.data || response

    if (responseData.code === 200 || responseData.success) {
      const clusters = responseData.data?.list || responseData.data || []
      clusterList.value = clusters.map(cluster => ({
        id: cluster.id,
        name: cluster.name,
        status: cluster.status
      }))

      if (clusterList.value.length > 0 && !props.modelValue) {
        const selectedId = resolvePreferredClusterId(clusterList.value, getStoredClusterId())

        if (selectedId) {
          saveStoredClusterId(selectedId)
        }

        emit('update:modelValue', selectedId)
        emit('change', selectedId)
      }

      console.log('集群列表加载成功:', clusterList.value)
    } else {
      ElMessage.error(responseData.message || '获取集群列表失败')
    }
  } catch (error) {
    console.error('获取集群列表失败:', error)
    ElMessage.warning('无法获取集群列表，请检查后端服务')
  }
}

const handleChange = (value) => {
  saveStoredClusterId(value)
  emit('update:modelValue', value)
  emit('change', value)
}

// 集群状态相关方法
const getClusterStatusTag = (status) => {
  const tagMap = {
    1: 'info',
    2: 'success',
    3: 'warning',
    0: 'danger'
  }
  return tagMap[status] || 'info'
}

const getClusterStatusText = (status) => {
  const textMap = {
    0: '离线',
    1: '连接中',
    2: '在线',
    3: '异常'
  }
  return textMap[status] || '未知'
}

onMounted(() => {
  fetchClusterList()
})
</script>

<style scoped>
/* 集群选择样式 */
.cluster-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}
.cluster-name {
  font-weight: 500;
  color: #6b9bd1;
}
.cluster-status-tag {
  margin-left: 8px;
}

.selector-prefix-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-right: 6px;
  font-size: 15px;
}

.cluster-prefix-icon {
  color: #67c23a;
}

.cluster-selector.theme-cluster :deep(.el-input__wrapper) {
  background: rgba(103, 194, 58, 0.08);
  border: 1px solid rgba(103, 194, 58, 0.28);
  box-shadow: none;
}

.cluster-selector.theme-cluster :deep(.el-input__wrapper:hover) {
  border-color: #67c23a;
}

.cluster-selector.theme-cluster :deep(.el-input__wrapper.is-focus) {
  border-color: #67c23a;
  box-shadow: 0 0 0 2px rgba(103, 194, 58, 0.18);
}

.cluster-selector.theme-cluster :deep(.el-select__selected-item),
.cluster-selector.theme-cluster :deep(.el-input__inner) {
  color: #4f8f1f;
  font-weight: 600;
}

.cluster-selector.theme-cluster :deep(.el-select__caret) {
  color: #67c23a;
}
</style>
