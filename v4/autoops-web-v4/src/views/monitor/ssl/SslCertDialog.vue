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
      <el-form-item label="域名" prop="domain">
        <el-input
          v-model="form.domain"
          placeholder="请输入域名，支持通配符如 *.example.com"
          clearable
        />
      </el-form-item>

      <el-form-item label="证书类型" prop="certType">
        <el-select v-model="form.certType" placeholder="请选择证书类型">
          <el-option label="ACME (免费)" :value="1" />
          <el-option label="阿里云CAS" :value="2" />
        </el-select>
      </el-form-item>

      <!-- 阿里云配置 - 两种方式都需要 -->
      <el-form-item label="阿里云配置" prop="aliyunConfigId">
        <el-select
          v-model="form.aliyunConfigId"
          placeholder="请选择阿里云配置"
          clearable
        >
          <el-option
            v-for="config in aliyunConfigs"
            :key="config.id"
            :label="config.name"
            :value="config.id"
          />
        </el-select>
        <div style="margin-top: 4px; color: #909399; font-size: 12px;">
          {{ isACME ? '用于DNS验证（需要阿里云DNS解析权限）' : '用于阿里云CAS证书申请' }}
        </div>
      </el-form-item>

      <!-- ACME方式的额外选项 -->
      <template v-if="isACME">
        <el-form-item label="CA提供商" prop="caProvider">
          <el-select v-model="form.caProvider" placeholder="请选择CA提供商">
            <el-option label="Let's Encrypt" value="LetsEncrypt" />
            <el-option label="ZeroSSL" value="ZeroSSL" />
            <el-option label="Google Trust Services" value="Google" />
            <el-option label="BuyPass" value="BuyPass" />
          </el-select>
        </el-form-item>

        <el-form-item label="加密算法" prop="algorithm">
          <el-select v-model="form.algorithm" placeholder="请选择加密算法">
            <el-option label="RSA 2048" value="RSA2048" />
            <el-option label="RSA 4096" value="RSA4096" />
            <el-option label="ECDSA 256" value="EC256" />
          </el-select>
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="form.email"
            placeholder="请输入邮箱地址"
            type="email"
            clearable
          />
          <div style="margin-top: 4px; color: #909399; font-size: 12px;">
            用于接收证书相关通知
          </div>
        </el-form-item>
      </template>
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
import { ref, computed, watch } from 'vue'
import { sslCertApi, aliyunConfigApi } from '@/api/monitor'
import { ElMessage } from 'element-plus'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '申请SSL证书'
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
const aliyunConfigs = ref([])

// 监听 visible prop
watch(() => props.visible, (newVal) => {
  isVisible.value = newVal
  if (newVal) {
    fetchAliyunConfigs()
  }
})

// 监听 isVisible 变化，更新 visible prop
watch(isVisible, (newVal) => {
  emit('update:visible', newVal)
})

// 判断是否ACME方式
const isACME = computed(() => {
  return props.form.certType === 1
})

// 表单验证规则
const formRules = {
  domain: [
    { required: true, message: '请输入域名', trigger: 'blur' },
    {
      pattern: /^(\*\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:[0-9]{1,5})?$/,
      message: '请输入有效的域名格式（支持通配符 *.example.com）',
      trigger: 'blur'
    }
  ],
  certType: [{ required: true, message: '请选择证书类型', trigger: 'change' }],
  aliyunConfigId: [
    { required: true, message: '请选择阿里云配置', trigger: 'change' }
  ],
  email: [
    {
      validator: (rule, value, callback) => {
        if (isACME.value && !value) {
          callback(new Error('ACME方式需要填写邮箱'))
        } else if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
          callback(new Error('请输入有效的邮箱地址'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 获取阿里云配置列表
const fetchAliyunConfigs = async () => {
  try {
    const { data: res } = await aliyunConfigApi.getList({
      pageNum: 1,
      pageSize: 100,
      status: 1
    })
    if (res.code === 200) {
      aliyunConfigs.value = Array.isArray(res.data) ? res.data : res.data?.list || []
    }
  } catch (error) {
    console.error('获取阿里云配置失败:', error)
    ElMessage.error('获取阿里云配置失败')
  }
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
