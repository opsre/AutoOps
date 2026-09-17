<template>
  <el-dialog
    v-model="visible"
    :title="`部署 ${service?.name}`"
    width="700px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="120px"
      v-loading="loading"
    >
      <!-- 版本选择 -->
      <el-form-item label="选择版本" prop="version">
        <el-select v-model="formData.version" placeholder="请选择版本" style="width: 100%">
          <el-option
            v-for="ver in service?.versions"
            :key="ver.id"
            :label="ver.name"
            :value="ver.id"
          >
            <span>{{ ver.name }}</span>
            <el-tag
              v-if="ver.recommended"
              type="success"
              size="small"
              style="margin-left: 10px"
            >
              推荐
            </el-tag>
          </el-option>
        </el-select>
      </el-form-item>

      <!-- 主机选择 -->
      <el-form-item label="目标主机" prop="hostId">
        <el-select
          v-model="formData.hostId"
          placeholder="请选择主机"
          filterable
          style="width: 100%"
        >
          <el-option
            v-for="host in hosts"
            :key="host.id"
            :label="`${host.hostName} (${host.sshIp || '-'})`"
            :value="host.id"
            :disabled="host.status !== 3"
          >
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <span>{{ host.hostName }}</span>
              <span style="color: #8492a6; font-size: 13px">
                {{ host.sshIp }}
              </span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>

      <!-- 安装目录 -->
      <el-form-item label="安装目录" prop="installDir">
        <el-input
          v-model="formData.installDir"
          placeholder="/opt/data/mysql"
        />
      </el-form-item>

      <!-- 环境变量 -->
      <el-form-item
        v-for="envVar in service?.env_vars"
        :key="envVar.name"
        :label="envVar.description"
        :prop="`envVars.${envVar.name}`"
        :rules="envVar.required ? [{ required: true, message: '不能为空' }] : []"
      >
        <el-input
          v-model="formData.envVars[envVar.name]"
          :placeholder="envVar.default"
          :type="envVar.name.includes('PASSWORD') ? 'password' : 'text'"
          show-password
        />
        <div class="env-tip">默认值: {{ envVar.default || '无' }}</div>
      </el-form-item>

      <!-- 自动启动 -->
      <el-form-item label="自动启动">
        <el-switch v-model="formData.autoStart" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button
        type="primary"
        @click="handleDeploy"
        :loading="deploying"
      >
        {{ deploying ? '部署中...' : '开始部署' }}
      </el-button>
    </template>

    <!-- 部署进度对话框 -->
    <DeployProgress
      v-model="progressVisible"
      :deploy-id="deployId"
      @close="handleProgressClose"
    />
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'
import { createDeploy } from '@/api/tool'
import cmdbAPI from '@/api/cmdb'
import { ElMessage } from 'element-plus'
import DeployProgress from './DeployProgress.vue'

const props = defineProps({
  modelValue: Boolean,
  service: Object
})

const emit = defineEmits(['update:modelValue', 'deployed'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const formRef = ref()
const deploying = ref(false)
const loading = ref(false)
const progressVisible = ref(false)
const deployId = ref(null)
const hosts = ref([])

const formData = reactive({
  version: '',
  hostId: null,
  installDir: '',
  envVars: {},
  autoStart: true
})

const rules = {
  version: [{ required: true, message: '请选择版本', trigger: 'change' }],
  hostId: [{ required: true, message: '请选择主机', trigger: 'change' }],
  installDir: [{ required: true, message: '请输入安装目录', trigger: 'blur' }]
}

// 加载主机列表（参考 Agent 页面，使用 Agent 列表以保证可部署性）
const loadHosts = async () => {
  loading.value = true
  try {
    // 拉取 Agent 列表，包含 hostId/hostName/sshIp/status
    const response = await cmdbAPI.getAgentList({ page: 1, pageSize: 1000 })
    if (response.data?.code === 200) {
      const list = response.data.data?.list || []
      // 统一映射为下拉所需字段：id(=hostId)、hostName、sshIp、status
      hosts.value = list.map(item => ({
        id: item.hostId,
        hostName: item.hostName,
        sshIp: item.sshIp,
        status: item.status
      }))
    } else {
      hosts.value = []
    }
  } catch (error) {
    console.error('加载主机列表失败:', error)
    ElMessage.error('加载主机列表失败')
  } finally {
    loading.value = false
  }
}

// 监听服务变化，初始化表单
watch(() => props.service, (newVal) => {
  if (newVal) {
    // 设置默认版本（推荐版本）
    const recommended = newVal.versions?.find(v => v.recommended)
    formData.version = recommended?.id || newVal.versions?.[0]?.id || ''

    // 设置默认安装目录
    formData.installDir = `/opt/data/${newVal.id}`

    // 初始化环境变量
    formData.envVars = {}
    newVal.env_vars?.forEach(env => {
      formData.envVars[env.name] = env.default
    })

    // 加载主机列表
    if (props.modelValue) {
      loadHosts()
    }
  }
}, { immediate: true })

// 监听对话框打开
watch(() => props.modelValue, (newVal) => {
  if (newVal && hosts.value.length === 0) {
    loadHosts()
  }
})

// 开始部署
const handleDeploy = async () => {
  try {
    await formRef.value.validate()

    deploying.value = true
    const res = await createDeploy({
      serviceId: props.service.id,
      version: formData.version,
      hostId: formData.hostId,
      installDir: formData.installDir,
      envVars: formData.envVars,
      autoStart: formData.autoStart
    })

    deploying.value = false

    if (res.data?.code === 200) {
      ElMessage.success('部署任务已创建')
      deployId.value = res.data.data?.deployId
      progressVisible.value = true
      visible.value = false
      emit('deployed', res.data.data?.deployId)
    } else {
      ElMessage.error(res.data?.message || '部署失败')
    }
  } catch (error) {
    deploying.value = false
    if (error !== 'cancel') {
      console.error('部署失败:', error)
      ElMessage.error(`部署失败: ${error.message || '未知错误'}`)
    }
  }
}

const handleClose = () => {
  formRef.value?.resetFields()
}

const handleProgressClose = () => {
  progressVisible.value = false
}
</script>

<style scoped>
.env-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

:deep(.el-dialog) {
  border-radius: 12px;
}

:deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  padding: 20px;
  border-radius: 12px 12px 0 0;
}

:deep(.el-dialog__title) {
  color: #fff;
  font-weight: 600;
}

:deep(.el-form-item__label) {
  color: #606266;
  font-weight: 500;
}

.el-input :deep(.el-input__wrapper) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.el-input :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
}

.el-select :deep(.el-input__wrapper) {
  border-radius: 8px;
}

.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
