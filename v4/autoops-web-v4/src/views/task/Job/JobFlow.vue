<template>
<el-dialog 
  title="任务流程" 
  v-model="dialogVisible" 
  width="80%"
  top="5vh"
  custom-class="task-flow-dialog"
  @close="handleDialogClose"
>
    <div class="flow-container">
      <task-flow ref="taskFlowRef" :steps="steps" :task-id="currentTaskId" />
    </div>
    <template #footer>
      <el-button @click="dialogVisible = false">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { GetTaskTemplates } from '@/api/task'
import TaskFlow from './TaskFlowTemp.vue'

const dialogVisible = ref(false)
const taskFlowRef = ref(null)
const taskCount = ref(0)
const currentTaskId = ref(null)
const templatesData = ref([])

  const steps = computed(() => {
    console.log('计算steps - currentTaskId:', currentTaskId.value);
    console.log('计算steps - templatesData:', JSON.parse(JSON.stringify(templatesData.value)));
    
    // 如果没有任务ID，使用模拟数据
    if (!currentTaskId.value) {
      const mockSteps = [
        {
          task_count: 1,
          task_status: 1,
          template_id: 1,
          template_name: '默认任务',
          template_remark: '等待启动',
          id: 'task-1',
          hasDownConnector: false,
          showHorizontalConnector: false
        }
      ];
      console.log('返回模拟steps:', mockSteps);
      return mockSteps;
    }
    
    // 使用存储的API数据
    if (currentTaskId.value && templatesData.value.length > 0) {
      const apiSteps = templatesData.value.map((template, index) => {
        console.log(`处理步骤${index} - 原始task_status:`, template.task_status);
        return {
          ...template,
          id: `task-${index+1}`,
          hasDownConnector: (index + 1) % 3 === 0 && index < templatesData.value.length - 1,
          showHorizontalConnector: (index + 1) % 3 !== 0 && index < templatesData.value.length - 1
        };
      });
      console.log('返回API steps:', JSON.parse(JSON.stringify(apiSteps)));
      return apiSteps;
    }
    
    console.log('返回空steps数组');
    return [];
  })

    const showFlow = async (taskId) => {
      try {
        console.log('Received taskId:', taskId);
        if (!taskId) {
          throw new Error('任务ID未定义');
        }
        
        // 重置当前数据
        currentTaskId.value = taskId;
        templatesData.value = [];
        
        const response = await GetTaskTemplates({ id: taskId });
        console.log('获取模板响应:', response);
        
        if (!response?.data?.data) {
          throw new Error('无效的API响应数据');
        }
        
        // 深度拷贝API数据
        const apiData = JSON.parse(JSON.stringify(response.data.data));
        console.log('API返回的原始数据:', apiData);
        
        // 确保每个步骤都有task_status
        const processedData = apiData.map(step => ({
          ...step,
          task_status: step.task_status ?? 1 // 默认pending状态
        }));
        
        templatesData.value = processedData;
        console.log('处理后的模板列表:', templatesData.value);
        
        // 从API响应中获取实际的task_id
        const apiTaskId = templatesData.value[0]?.task_id;
        if (!apiTaskId) {
          throw new Error('API响应中未包含task_id');
        }
        
        // 更新currentTaskId为API返回的task_id
        currentTaskId.value = apiTaskId;
        taskCount.value = templatesData.value.length;
        
        console.log('更新后的templatesData:', templatesData.value);
        console.log('更新后的taskCount:', taskCount.value);
        
        dialogVisible.value = true;
        return apiTaskId;
      } catch (error) {
        console.error('获取模板信息失败:', error);
        templatesData.value = [];
        taskCount.value = 1;
        dialogVisible.value = true;
        throw error;
      }
    }
const handleDialogClose = () => {
  if (taskFlowRef.value?.stopAllPolling) {
    taskFlowRef.value.stopAllPolling()
  }
}

defineExpose({
  showFlow
})
</script>

<style scoped>
.flow-container {
  min-height: 400px;
  padding: 20px;
  background: #e8e8e8;
  border-radius: 6px;
}

:deep(.task-flow-dialog) {
  background: #e6e6e6;
}

:deep(.task-flow-dialog .el-dialog__header) {
  background: #e6e6e6;
  border-bottom: 1px solid #d9d9d9;
}

:deep(.task-flow-dialog .el-dialog__body) {
  padding: 20px;
  background: #e6e6e6;
}

:deep(.task-flow-dialog .el-dialog__footer) {
  background: #e6e6e6;
  border-top: 1px solid #d9d9d9;
}
</style>
