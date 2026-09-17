<template>
  <el-dialog
    v-model="dialogVisible"
    title="启动任务"
    width="500px"
    :before-close="handleClose"
  >
    <div class="start-task-content">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="任务名称">
          <el-input :value="taskName" disabled />
        </el-form-item>
        
        <el-form-item label="确认启动" prop="confirm">
          <el-checkbox v-model="form.confirm">
            我确认要启动此Ansible任务
          </el-checkbox>
        </el-form-item>
        
        <div class="task-info" v-if="taskInfo">
          <p><strong>任务ID:</strong> {{ taskInfo.id }}</p>
          <p><strong>描述:</strong> {{ taskInfo.description || '无描述' }}</p>
          <p><strong>预计执行时间:</strong> {{ taskInfo.estimatedTime || '未知' }}</p>
        </div>
      </el-form>
    </div>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button 
          type="primary" 
          @click="handleConfirm"
          :loading="loading"
          :disabled="!form.confirm"
        >
          {{ loading ? '启动中...' : '确认启动' }}
        </el-button>
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
  taskName: {
    type: String,
    default: ''
  },
  taskInfo: {
    type: Object,
    default: () => ({})
  }
})

// Emits
const emit = defineEmits(['update:visible', 'confirm', 'cancel'])

// 响应式数据
const dialogVisible = ref(false)
const loading = ref(false)
const formRef = ref()

const form = reactive({
  confirm: false
})

const rules = {
  confirm: [
    { 
      validator: (rule, value, callback) => {
        if (!value) {
          callback(new Error('请确认启动任务'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ]
}

// 监听visible变化
watch(() => props.visible, (newVal) => {
  dialogVisible.value = newVal
  if (newVal) {
    // 重置表单
    form.confirm = false
  }
})

watch(dialogVisible, (newVal) => {
  emit('update:visible', newVal)
})

// 处理确认
const handleConfirm = async () => {
  try {
    await formRef.value?.validate()
    
    loading.value = true
    
    // 发送确认事件
    emit('confirm', {
      taskInfo: props.taskInfo,
      confirmed: true
    })
    
    // 短暂延迟以显示加载状态
    setTimeout(() => {
      loading.value = false
      dialogVisible.value = false
      ElMessage.success('任务启动成功')
    }, 1000)
    
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

// 处理关闭
const handleClose = () => {
  dialogVisible.value = false
  emit('cancel')
}
</script>

<style scoped>
.start-task-content {
  padding: 20px 0;
}

.task-info {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 6px;
  margin-top: 15px;
  border-left: 4px solid #409eff;
}

.task-info p {
  margin: 5px 0;
  font-size: 14px;
  color: #606266;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>