<template>
  <el-dialog
    :model-value="visible"
    :title="`修改标签 - ${namespaceName}`"
    width="600px"
    class="ns-labels-dialog"
    :close-on-click-modal="false"
    @update:model-value="handleClose"
    @open="initLabels"
  >
    <div class="labels-content">
      <div class="labels-tip">
        <el-alert type="warning" :closable="false" show-icon>
          <template #title>标签将全量替换，提交后旧标签全部覆盖</template>
        </el-alert>
      </div>

      <div class="labels-list">
        <div v-for="(item, index) in labelItems" :key="index" class="label-row">
          <el-input v-model="item.key" placeholder="key" size="small" class="label-key" />
          <span class="label-sep">=</span>
          <el-input v-model="item.value" placeholder="value" size="small" class="label-value" />
          <el-button :icon="Delete" type="danger" size="small" circle @click="removeLabel(index)" />
        </div>

        <el-button :icon="Plus" type="primary" size="small" plain @click="addLabel">添加标签</el-button>
      </div>

      <div class="action-buttons">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="saveLabels" :loading="saving">保存</el-button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Delete, Plus } from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'

const props = defineProps({
  visible: { type: Boolean, default: false },
  clusterId: { type: [String, Number], required: true },
  namespaceName: { type: String, default: '' },
  currentLabels: { type: Object, default: () => ({}) }
})
const emit = defineEmits(['update:visible', 'saved'])

const labelItems = ref([])
const saving = ref(false)

const systemPrefixes = ['kubernetes.io/', 'beta.kubernetes.io/', 'app.kubernetes.io/']
const isSystemLabel = (key) => systemPrefixes.some(p => key.startsWith(p))

const initLabels = () => {
  // 过滤系统标签，只展示用户标签
  labelItems.value = Object.entries(props.currentLabels || {})
    .filter(([key]) => !isSystemLabel(key))
    .map(([key, value]) => ({ key, value }))
}

watch(() => props.visible, (val) => { if (val) initLabels() })

const addLabel = () => { labelItems.value.push({ key: '', value: '' }) }
const removeLabel = (index) => { labelItems.value.splice(index, 1) }

const saveLabels = async () => {
  // 验证 key 不为空
  for (const item of labelItems.value) {
    if (!item.key.trim()) { ElMessage.warning('标签 key 不能为空'); return }
  }
  // 构建标签对象（用户标签 + 保留系统标签）
  const userLabels = Object.fromEntries(labelItems.value.map(i => [i.key.trim(), i.value.trim()]))
  const systemLabels = Object.fromEntries(
    Object.entries(props.currentLabels || {}).filter(([key]) => isSystemLabel(key))
  )
  const labels = { ...systemLabels, ...userLabels }

  saving.value = true
  try {
    const res = await k8sApi.updateNamespaceLabels(props.clusterId, props.namespaceName, labels)
    const data = res.data || res
    if (data.code === 200 || data.success) {
      ElMessage.success('标签更新成功')
      emit('saved')
      handleClose()
    } else {
      ElMessage.error(data.message || '更新失败')
    }
  } catch {
    ElMessage.error('更新标签失败')
  } finally {
    saving.value = false
  }
}

const handleClose = () => emit('update:visible', false)
</script>

<style scoped>
.labels-content { display: flex; flex-direction: column; gap: 16px; }
.labels-tip { margin-bottom: 4px; }
.labels-list { display: flex; flex-direction: column; gap: 10px; }
.label-row { display: flex; align-items: center; gap: 8px; }
.label-key { flex: 1; }
.label-value { flex: 2; }
.label-sep { color: #909399; font-weight: 700; flex-shrink: 0; }
.action-buttons { display: flex; justify-content: flex-end; gap: 12px; padding-top: 16px; border-top: 1px solid #e5e7eb; }
.ns-labels-dialog :deep(.el-dialog) { border-radius: 16px; }
.ns-labels-dialog :deep(.el-dialog__header) { background: linear-gradient(135deg, #667eea, #764ba2); border-top-left-radius: 16px; border-top-right-radius: 16px; padding: 20px 24px; }
.ns-labels-dialog :deep(.el-dialog__title) { color: white; font-weight: 600; }
.ns-labels-dialog :deep(.el-dialog__body) { padding: 20px; }
</style>
