<template>
  <el-dialog
    :model-value="visible"
    :title="`更新Pod配置 - ${workload.name || ''}`"
    width="600px"
    class="pod-config-dialog"
    @update:model-value="handleClose"
  >
    <el-form :model="configForm" label-width="120px">
      <el-form-item label="资源限制">
        <div class="resource-config">
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="resource-section">
                <h5>CPU</h5>
                <el-form-item label="Request:" size="small">
                  <el-input
                    v-model="configForm.cpuRequest"
                    placeholder="如: 100m"
                    style="width: 100%"
                  />
                </el-form-item>
                <el-form-item label="Limit:" size="small">
                  <el-input
                    v-model="configForm.cpuLimit"
                    placeholder="如: 500m"
                    style="width: 100%"
                  />
                </el-form-item>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="resource-section">
                <h5>内存</h5>
                <el-form-item label="Request:" size="small">
                  <el-input
                    v-model="configForm.memoryRequest"
                    placeholder="如: 128Mi"
                    style="width: 100%"
                  />
                </el-form-item>
                <el-form-item label="Limit:" size="small">
                  <el-input
                    v-model="configForm.memoryLimit"
                    placeholder="如: 512Mi"
                    style="width: 100%"
                  />
                </el-form-item>
              </div>
            </el-col>
          </el-row>
        </div>
      </el-form-item>
      <el-form-item label="环境变量">
        <div class="env-vars-section">
          <div v-for="(envVar, index) in configForm.envVars" :key="index" class="env-var-item">
            <el-row :gutter="10">
              <el-col :span="10">
                <el-input
                  v-model="envVar.name"
                  placeholder="变量名"
                  size="small"
                />
              </el-col>
              <el-col :span="10">
                <el-input
                  v-model="envVar.value"
                  placeholder="变量值"
                  size="small"
                />
              </el-col>
              <el-col :span="4">
                <el-button
                  type="danger"
                  size="small"
                  @click="removeEnvVar(index)"
                >
                  删除
                </el-button>
              </el-col>
            </el-row>
          </div>
          <el-button type="text" size="small" @click="addEnvVar">+ 添加环境变量</el-button>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">更新配置</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch, defineProps, defineEmits } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  workload: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:visible', 'close', 'submit'])

const submitting = ref(false)

// 配置表单数据
const configForm = reactive({
  cpuRequest: '',
  cpuLimit: '',
  memoryRequest: '',
  memoryLimit: '',
  envVars: []
})

// 监听 workload 变化，初始化表单数据
watch(() => props.workload, (newWorkload) => {
  if (newWorkload && Object.keys(newWorkload).length > 0) {
    // 初始化资源配置
    configForm.cpuRequest = newWorkload.resources?.requests?.cpu || ''
    configForm.cpuLimit = newWorkload.resources?.limits?.cpu || ''
    configForm.memoryRequest = newWorkload.resources?.requests?.memory || ''
    configForm.memoryLimit = newWorkload.resources?.limits?.memory || ''

    // 初始化环境变量
    configForm.envVars = newWorkload.envVars || []
  }
}, { immediate: true, deep: true })

const handleClose = () => {
  emit('update:visible', false)
  emit('close')
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    // 构建更新数据，按照正确的API格式
    const updateData = {
      template: {
        labels: {
          app: props.workload.labels?.app || props.workload.name
        },
        containers: [{
          name: props.workload.containerName || props.workload.name || 'nginx',
          image: props.workload.image || 'nginx:latest', // 必需字段
          resources: {
            requests: {
              cpu: configForm.cpuRequest || '100m',
              memory: configForm.memoryRequest || '128Mi'
            },
            limits: {
              cpu: configForm.cpuLimit || '500m',
              memory: configForm.memoryLimit || '512Mi'
            }
          },
          env: configForm.envVars.filter(env => env.name && env.value).map(env => ({
            name: env.name,
            value: env.value
          }))
        }]
      }
    }

    // 如果有端口配置，添加到容器中
    if (props.workload.ports && props.workload.ports.length > 0) {
      updateData.template.containers[0].ports = props.workload.ports
    }

    emit('submit', updateData)
  } finally {
    submitting.value = false
  }
}

// 添加环境变量
const addEnvVar = () => {
  configForm.envVars.push({ name: '', value: '' })
}

// 删除环境变量
const removeEnvVar = (index) => {
  configForm.envVars.splice(index, 1)
}
</script>

<style scoped>
.resource-config {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 16px;
  background-color: #fafafa;
}

.resource-section h5 {
  margin: 0 0 12px 0;
  color: #303133;
  font-weight: 600;
}

.env-vars-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.env-var-item {
  padding: 8px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  background-color: #f5f7fa;
}

.dialog-footer {
  text-align: right;
}

/* 对话框样式 */
.pod-config-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.pod-config-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.pod-config-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.pod-config-dialog :deep(.el-form-item__label) {
  font-weight: 500;
}
</style>