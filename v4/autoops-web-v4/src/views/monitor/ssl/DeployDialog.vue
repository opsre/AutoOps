<template>
  <el-dialog
    v-model="isVisible"
    title="部署证书"
    width="700px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      label-width="120px"
      :rules="formRules"
    >
      <el-form-item label="证书">
        <span v-if="cert">{{ cert.domain }} (ID: {{ cert.id }})</span>
      </el-form-item>

      <el-form-item label="目标主机" prop="hostId">
        <el-select
          v-model="form.hostId"
          placeholder="请选择目标主机"
          filterable
        >
          <el-option
            v-for="host in hostOptions"
            :key="host.value"
            :label="host.label"
            :value="host.value"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="部署路径" prop="deployPath">
        <el-input
          v-model="form.deployPath"
          placeholder="/etc/nginx/ssl"
        />
      </el-form-item>

      <!-- 部署日志显示 -->
      <el-form-item v-if="logs.length > 0" label="部署日志">
        <el-alert
          v-for="(log, index) in logs"
          :key="index"
          :title="log"
          type="info"
          closable
          class="log-item"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          确定部署
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  cert: {
    type: Object,
    default: null
  },
  form: {
    type: Object,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  hostOptions: {
    type: Array,
    default: () => []
  },
  logs: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:visible', 'submit', 'close'])

const isVisible = ref(false)
const formRef = ref(null)

// 监听 visible prop
watch(() => props.visible, (newVal) => {
  isVisible.value = newVal
})

// 监听 isVisible 变化，更新 visible prop
watch(isVisible, (newVal) => {
  if (!newVal) {
    emit('update:visible', false)
  }
})

// 表单验证规则
const formRules = {
  hostId: [{ required: true, message: '请选择目标主机', trigger: 'change' }],
  deployPath: [{ required: true, message: '请输入部署路径', trigger: 'blur' }]
}

// 关闭对话框
const handleClose = () => {
  isVisible.value = false
  emit('close')
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate((valid) => {
    if (valid) {
      emit('submit')
    }
  }).catch(error => {
    console.error('表单验证失败:', error)
  })
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.log-item {
  margin-bottom: 10px;
  word-break: break-all;
}
</style>
