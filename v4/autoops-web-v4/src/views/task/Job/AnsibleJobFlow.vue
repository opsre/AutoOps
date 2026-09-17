<template>
  <el-dialog 
    title="Ansible任务流程" 
    v-model="dialogVisible" 
    width="80%"
    top="5vh"
    custom-class="ansible-flow-dialog"
    @close="handleDialogClose"
  >
    <div class="flow-container">
      <ansible-flow 
        ref="ansibleFlowRef" 
        :steps="steps" 
        :task-id="currentTaskId"
        @refresh="handleRefresh"
      />
    </div>
    <template #footer>
      <el-button @click="dialogVisible = false">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { GetAnsibleTaskDetail } from '@/api/task'
import AnsibleFlow from './AnsibleFlowTemp.vue'

const dialogVisible = ref(false)
const ansibleFlowRef = ref(null)
const currentTaskId = ref(null)
const taskData = ref(null)
const worksData = ref([])

const steps = computed(() => {
  console.log('计算Ansible steps - currentTaskId:', currentTaskId.value);
  console.log('计算Ansible steps - worksData:', JSON.parse(JSON.stringify(worksData.value)));
  
  // 如果没有任务ID，使用模拟数据
  if (!currentTaskId.value) {
    const mockSteps = [
      {
        work_id: 1,
        entry_file_name: '默认任务',
        status: 1,
        duration: 0,
        id: 'ansible-work-1',
        hasDownConnector: false,
        showHorizontalConnector: false
      }
    ];
    console.log('返回模拟Ansible steps:', mockSteps);
    return mockSteps;
  }
  
  // 使用存储的Works数据
  if (currentTaskId.value && worksData.value.length > 0) {
    const ansibleSteps = worksData.value.map((work, index) => {
      console.log(`处理步骤${index} - 原始status:`, work.status);
      return {
        ...work,
        task_id: currentTaskId.value,
        work_id: work.workid,
        entry_file_name: work.EntryFileName,
        status: work.status,
        duration: work.Duration,
        id: `ansible-work-${index+1}`,
        hasDownConnector: (index + 1) % 3 === 0 && index < worksData.value.length - 1,
        showHorizontalConnector: (index + 1) % 3 !== 0 && index < worksData.value.length - 1
      };
    });
    console.log('返回Ansible API steps:', JSON.parse(JSON.stringify(ansibleSteps)));
    return ansibleSteps;
  }
  
  console.log('返回空Ansible steps数组');
  return [];
})

const showFlow = async (taskId) => {
  try {
    console.log('Received Ansible taskId:', taskId);
    if (!taskId) {
      throw new Error('Ansible任务ID未定义');
    }
    
    // 重置当前数据
    currentTaskId.value = taskId;
    taskData.value = null;
    worksData.value = [];
    
    const response = await GetAnsibleTaskDetail(taskId);
    console.log('获取Ansible任务详情响应:', response);
    
    if (!response?.data?.data?.task_info) {
      throw new Error('无效的API响应数据');
    }
    
    // 深度拷贝API数据
    const apiData = JSON.parse(JSON.stringify(response.data.data));
    console.log('Ansible API返回的原始数据:', apiData);
    
    taskData.value = apiData.task_info;
    worksData.value = apiData.task_info.Works || [];
    
    console.log('处理后的任务数据:', taskData.value);
    console.log('处理后的Works列表:', worksData.value);
    
    dialogVisible.value = true;
    return currentTaskId.value;
  } catch (error) {
    console.error('获取Ansible任务信息失败:', error);
    worksData.value = [];
    dialogVisible.value = true;
    throw error;
  }
}

const handleRefresh = async () => {
  if (currentTaskId.value) {
    await showFlow(currentTaskId.value);
  }
}

const handleDialogClose = () => {
  if (ansibleFlowRef.value?.stopAllPolling) {
    ansibleFlowRef.value.stopAllPolling()
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

:deep(.ansible-flow-dialog) {
  background: #e6e6e6;
}

:deep(.ansible-flow-dialog .el-dialog__header) {
  background: #e6e6e6;
  border-bottom: 1px solid #d9d9d9;
}

:deep(.ansible-flow-dialog .el-dialog__body) {
  padding: 20px;
  background: #e6e6e6;
}

:deep(.ansible-flow-dialog .el-dialog__footer) {
  background: #e6e6e6;
  border-top: 1px solid #d9d9d9;
}
</style>