<template>
  <el-dialog
    v-model="dialogVisible"
    title="注册 K8s 集群"
    width="600px"
    class="register-dialog"
    @closed="handleDialogClosed"
  >
    <el-form :model="form" label-width="120px" ref="formRef">
      <el-form-item 
        label="集群名称" 
        required
        prop="clusterName"
        :rules="[
          { required: true, message: '请输入集群名称', trigger: 'blur' },
          { min: 2, max: 50, message: '集群名称长度在 2 到 50 个字符', trigger: 'blur' }
        ]"
      >
        <el-input
          v-model="form.clusterName"
          placeholder="请输入集群名称"
          clearable
          :disabled="loading"
        />
      </el-form-item>
      
      <el-form-item 
        label="集群凭证" 
        required
        prop="credential"
        :rules="[
          { required: true, message: '请输入集群凭证', trigger: 'blur' }
        ]"
      >
        <CredentialEditor
          v-model="form.credential"
          height="200px"
          placeholder="请输入或上传集群凭证配置"
          :disabled="loading"
        />
      </el-form-item>
      
      <el-form-item label="描述信息">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="2"
          placeholder="请输入集群描述"
          :disabled="loading"
          maxlength="200"
          show-word-limit
        />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleCancel" :disabled="loading">取消</el-button>
        <el-button 
          type="primary" 
          @click="handleSubmit"
          :loading="loading"
        >
          注册
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import k8sApi from '@/api/k8s'
import CredentialEditor from '../components/CredentialEditor.vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'success', 'closed'])

const formRef = ref(null)
const loading = ref(false)

// 双向绑定 visible
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 表单数据
const form = reactive({
  clusterName: '',
  credential: '',
  description: ''
})

// 重置表单
const resetForm = () => {
  form.clusterName = ''
  form.credential = ''
  form.description = ''
  // 清除验证状态
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 提交注册
const handleSubmit = async () => {
  try {
    // 表单验证
    if (!formRef.value) {
      ElMessage.error('表单验证失败')
      return
    }
    
    const isValid = await formRef.value.validate().catch(() => false)
    if (!isValid) {
      return
    }
    
    loading.value = true
    
    const response = await k8sApi.registerExternalCluster({
      name: form.clusterName,
      description: form.description,
      clusterType: 2,  // 导入集群
      kubeconfig: form.credential
    })
    
    const responseData = response.data || response
    console.log('注册集群API响应:', responseData)

    if (responseData.code === 200 || responseData.success) {
      ElMessage.success('集群注册成功，正在同步集群信息...')
      dialogVisible.value = false
      resetForm()
      emit('success', responseData)
    } else {
      ElMessage.error(responseData.message || '注册集群失败')
    }
  } catch (error) {
    console.error('注册集群失败:', error)
    
    // 根据错误类型给出更友好的提示
    if (error.response?.status === 400) {
      ElMessage.error('请求参数错误，请检查集群信息')
    } else if (error.response?.status === 409) {
      ElMessage.error('集群名称已存在，请使用其他名称')
    } else if (error.code === 'ERR_NETWORK') {
      ElMessage.error('网络连接失败，请检查网络状态')
    } else {
      ElMessage.error('注册集群失败，请重试')
    }
  } finally {
    loading.value = false
  }
}

// 取消操作
const handleCancel = () => {
  dialogVisible.value = false
}

// 对话框关闭处理
const handleDialogClosed = () => {
  resetForm()
  emit('closed')
}

// 监听对话框打开，重置表单
watch(() => props.visible, (newVisible) => {
  if (newVisible) {
    resetForm()
  }
})
</script>

<style scoped>
/* 注册对话框样式 */
.register-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 3px solid #409EFF;
}

.register-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.register-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
  font-size: 16px;
}

.register-dialog :deep(.el-dialog__headerbtn) {
  top: 20px;
  right: 20px;
}

.register-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
  color: white;
  font-size: 18px;
}

.register-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

/* 表单样式 */
.register-dialog .el-form-item {
  margin-bottom: 24px;
}

.register-dialog .el-form-item__label {
  font-weight: 500;
  color: #606266;
}

.register-dialog .el-input :deep(.el-input__wrapper) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

.register-dialog .el-input :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.register-dialog .el-input :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
}

.register-dialog .el-textarea :deep(.el-textarea__inner) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

.register-dialog .el-textarea :deep(.el-textarea__inner):hover {
  border-color: #c0c4cc;
}

.register-dialog .el-textarea :deep(.el-textarea__inner):focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
}

/* 按钮样式 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.dialog-footer .el-button {
  border-radius: 8px;
  font-weight: 500;
  padding: 12px 24px;
  transition: all 0.3s ease;
}

.dialog-footer .el-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.dialog-footer .el-button--primary {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
}

.dialog-footer .el-button--primary:hover {
  background: linear-gradient(135deg, #5a6fd8, #6a4190);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .register-dialog :deep(.el-dialog) {
    width: 95% !important;
    margin: 5vh auto !important;
  }
  
  .register-dialog :deep(.el-dialog__body) {
    padding: 20px;
  }
  
  .dialog-footer {
    flex-direction: column;
  }
  
  .dialog-footer .el-button {
    width: 100%;
    margin: 0;
  }
}
</style>