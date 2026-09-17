<template>
  <el-dialog
    v-model="isVisible"
    :title="title"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      label-width="120px"
      :rules="formRules"
    >
      <el-form-item label="配置名称" prop="name">
        <el-input
          v-model="form.name"
          placeholder="请输入配置名称"
          clearable
        />
      </el-form-item>

      <el-form-item label="AccessKey" prop="accessKey">
        <el-input
          v-model="form.accessKey"
          placeholder="请输入AccessKey"
          clearable
          show-password
        />
      </el-form-item>

      <el-form-item label="AccessSecret" prop="accessSecret">
        <el-input
          v-model="form.accessSecret"
          placeholder="请输入AccessSecret"
          clearable
          show-password
        />
      </el-form-item>

      <el-form-item label="邮箱" prop="email">
        <el-input
          v-model="form.email"
          placeholder="请输入邮箱地址"
          type="email"
          clearable
        />
      </el-form-item>

      <el-form-item label="电话" prop="phone">
        <el-input
          v-model="form.phone"
          placeholder="请输入电话号码"
          clearable
        />
      </el-form-item>

      <el-form-item label="用户名" prop="username">
        <el-input
          v-model="form.username"
          placeholder="请输入用户名"
          clearable
        />
      </el-form-item>

      <el-form-item label="状态" prop="status">
        <el-select v-model="form.status" placeholder="请选择状态">
          <el-option label="启用" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
      </el-form-item>
    </el-form>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          确定
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
  title: {
    type: String,
    default: '添加阿里云配置'
  },
  form: {
    type: Object,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'submit'])

const isVisible = ref(false)
const formRef = ref(null)

// 监听 visible prop
watch(() => props.visible, (newVal) => {
  isVisible.value = newVal
})

// 监听 isVisible 变化，更新 visible prop
watch(isVisible, (newVal) => {
  emit('update:visible', newVal)
})

// 表单验证规则
const formRules = {
  name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
  accessKey: [{ required: true, message: '请输入AccessKey', trigger: 'blur' }],
  accessSecret: [{ required: true, message: '请输入AccessSecret', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  phone: [{ required: true, message: '请输入电话号码', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

// 关闭对话框
const handleClose = () => {
  isVisible.value = false
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
</style>
