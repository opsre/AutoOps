<template>
  <!-- 创建命名空间对话框 -->
  <el-dialog
    :model-value="visible"
    title="创建命名空间"
    width="500px"
    class="create-namespace-dialog"
    @update:model-value="handleClose"
    @close="handleClose"
  >
    <el-form :model="createForm" label-width="100px" ref="formRef">
      <el-form-item 
        label="名称" 
        required
        :rules="[
          { required: true, message: '请输入命名空间名称', trigger: 'blur' },
          { pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/, message: '命名空间名称只能包含小写字母、数字和连字符，且必须以字母数字字符开头和结尾', trigger: 'blur' }
        ]"
      >
        <el-input
          v-model="createForm.name"
          placeholder="请输入命名空间名称"
        />
      </el-form-item>
      
      <el-form-item label="描述">
        <el-input
          v-model="createForm.description"
          placeholder="请输入描述信息"
          type="textarea"
          :rows="2"
        />
      </el-form-item>

      <el-form-item label="标签">
        <el-input
          v-model="createForm.customLabels"
          placeholder="格式：key=value，每行一个"
          type="textarea"
          :rows="3"
        />
        <div class="form-hint">
        </div>
      </el-form-item>
    </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">创建</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'

// Props
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  selectedClusterId: {
    type: [String, Number],
    required: true
  }
})

// Emits
const emit = defineEmits(['update:visible', 'success'])

// Form ref
const formRef = ref(null)

// 提交状态
const submitting = ref(false)

// 创建命名空间表单
const createForm = reactive({
  name: '',
  description: '',
  customLabels: ''
})

// 监听 visible 变化，重置表单
watch(() => props.visible, (newVal) => {
  if (newVal) {
    resetForm()
  }
})

// 重置表单
const resetForm = () => {
  createForm.name = ''
  createForm.description = ''
  createForm.customLabels = ''
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 验证标签键值格式
const validateLabel = (key, value) => {
  // 标签键验证：最多253字符，格式为[前缀/]名称
  const keyRegex = /^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\/)?[a-zA-Z0-9]([-_.a-zA-Z0-9]*[a-zA-Z0-9])?$/
  if (key.length > 253 || !keyRegex.test(key)) {
    return `标签键 "${key}" 格式不正确。最多253字符，格式：[前缀/]名称`
  }
  
  // 标签值验证：最多63字符，只能包含字母数字、连字符、下划线、点号
  const valueRegex = /^[a-zA-Z0-9]([-_.a-zA-Z0-9]*[a-zA-Z0-9])?$/
  if (value && (value.length > 63 || !valueRegex.test(value))) {
    return `标签值 "${value}" 格式不正确。最多63字符，只能包含字母数字、连字符、下划线、点号`
  }
  
  return null
}

// 解析键值对字符串
const parseKeyValuePairs = (str) => {
  const result = {}
  if (!str.trim()) return result
  
  const pairs = str.split('\n').filter(line => line.trim())
  for (const pair of pairs) {
    const [key, ...valueParts] = pair.split('=')
    if (key && key.trim()) {
      result[key.trim()] = valueParts.join('=').trim() || ''
    }
  }
  return result
}

// 处理取消
const handleCancel = () => {
  emit('update:visible', false)
}

// 处理关闭
const handleClose = () => {
  emit('update:visible', false)
}

// 处理提交
const handleSubmit = async () => {
  try {
    // 表单验证
    if (!formRef.value) return
    const valid = await formRef.value.validate()
    if (!valid) return

    if (!createForm.name) {
      ElMessage.warning('请输入命名空间名称')
      return
    }

    submitting.value = true
    
    const data = {
      name: createForm.name,
      labels: {
        'app.kubernetes.io/managed-by': 'k8s-platform'
      },
      annotations: {}
    }
    
    // 解析自定义标签
    if (createForm.customLabels) {
      const customLabels = parseKeyValuePairs(createForm.customLabels)
      for (const [key, value] of Object.entries(customLabels)) {
        const error = validateLabel(key, value)
        if (error) {
          ElMessage.warning(error)
          return
        }
        data.labels[key] = value
      }
    }
    
    // 描述作为注释
    if (createForm.description) {
      data.annotations['description'] = createForm.description
    }
    
    // 添加创建者注释
    data.annotations['kubernetes.io/created-by'] = 'k8s-platform'
    data.annotations['created-at'] = new Date().toISOString()
    
    // 发射成功事件，由父组件处理API调用
    emit('success', data)
    
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    submitting.value = false
  }
}

// 暴露方法供父组件调用
defineExpose({
  resetForm,
  setSubmitting: (loading) => {
    submitting.value = loading
  }
})
</script>

<style scoped>
/* 对话框样式 */
.create-namespace-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.create-namespace-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.create-namespace-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.create-namespace-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

.dialog-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

/* 表单提示样式 */
.form-hint {
  margin-top: 4px;
  padding: 4px 8px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 4px;
  border-left: 3px solid rgba(103, 126, 234, 0.3);
}

.form-hint .el-text {
  line-height: 1.4;
}

.create-namespace-dialog :deep(.el-divider__text) {
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 600;
}

.create-namespace-dialog :deep(.el-divider) {
  margin: 20px 0 16px 0;
}

.create-namespace-dialog :deep(.el-textarea__inner) {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
}

/* 通用元素样式 */
.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.el-input :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.el-input :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}
</style>