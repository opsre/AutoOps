<template>
  <div class="ansible-flow-container" style="margin-top: -5px !important; top: -5px !important;">
    <div class="flow-header">
      <el-button v-if="!historyMode" v-authority="['task:ansible:jobstart']" @click="startAnsibleTask" type="primary" size="medium" style="height: 40px; padding: 0 20px;" :loading="starting">
        <el-icon size="20"><VideoPlay /></el-icon>
        <span style="margin-left: 8px;">启动任务</span>
      </el-button>
      <el-button @click="handleRefresh" type="info" size="medium" style="height: 40px; padding: 0 20px;">
        <el-icon size="20"><Refresh /></el-icon>
        <span style="margin-left: 8px;">刷新状态</span>
      </el-button>
    </div>
    <svg class="connector-canvas" ref="connectorCanvas"></svg>
    <div 
      v-for="(step, index) in orderedSteps" 
      :key="index" 
      class="flow-step"
      :style="{'--step-order': step.order}"
    >
      <svg class="connector-defs">
        <defs>
          <marker :id="'arrowhead-' + index" markerWidth="6" markerHeight="4" 
                  refX="5" refY="2" orient="auto">
            <path d="M0,0 L0,4 L6,2 Z" :fill="step.status === 'completed' ? 'rgba(76,175,80,0.7)' : 
                  step.status === 'active' ? 'rgba(33,150,243,0.7)' : 'rgba(128,128,128,0.7)'"/>
          </marker>
          <marker :id="'flow-arrow-' + index" markerWidth="6" markerHeight="4"
                  refX="5" refY="2" orient="auto">
            <path d="M0,0 L0,4 L6,2 Z" :fill="step.status === 'completed' ? 'rgba(76,175,80,0.7)' : 
                  step.status === 'active' ? 'rgba(33,150,243,0.7)' : 'rgba(128,128,128,0.7)'"/>
          </marker>
        </defs>
      </svg>
      <div :class="['flow-card', getStepClass(step)]">
        <div class="status-text" :class="'status-' + step.status">
          {{
            step.status === 'completed' ? '完成' :
            step.status === 'active' ? '运行中' :
            step.status === 'error' ? '异常' : '等待'
          }}
        </div>
        <h3 class="card-title">
          <svg class="title-icon" viewBox="0 0 1024 1024" width="20" height="20">
            <path fill="#E6522C" d="M512 1024C229.23 1024 0 794.77 0 512S229.23 0 512 0s512 229.23 512 512-229.23 512-512 512zm0-938.67C276.36 85.33 85.33 276.36 85.33 512S276.36 938.67 512 938.67 938.67 747.64 938.67 512 747.64 85.33 512 85.33z"/>
            <path fill="#E6522C" d="M341.33 426.67a42.67 42.67 0 100-85.34 42.67 42.67 0 000 85.34zM682.67 426.67a42.67 42.67 0 100-85.34 42.67 42.67 0 000 85.34z"/>
            <path fill="#E6522C" d="M512 768c-117.82 0-213.33-95.51-213.33-213.33h426.66C725.33 672.49 629.82 768 512 768z"/>
          </svg>
          {{ step.title }}
        </h3>
        <p class="card-content">任务名称: {{ step.entry_file_name }}</p>
        <div class="card-duration" v-if="step.duration > 0">
          执行时长: {{ formatDuration(step.duration) }}秒
        </div>
        <div class="card-buttons">
          <el-tooltip effect="dark" content="查看日志" placement="top">
            <button @click.stop="onLogClick(step)" class="btn-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M21 10H7M21 6H3M21 14H3M21 18H7" stroke="#722ed1" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </button>
          </el-tooltip>
        </div>
      </div>
    </div>
    <div style="position: fixed; z-index: 9999;">
      <ansible-log-dialog 
        ref="ansibleLogDialog"
        :history-mode="historyMode"
        :history-id="historyId"
      ></ansible-log-dialog>
    </div>
  </div>
</template>

<script>
import { GetAnsibleTaskLog, StartAnsibleTaskFlow, GetAnsibleTaskDetail, GetAnsibleHistoryDetail } from '@/api/task'
import { VideoPlay, Refresh } from '@element-plus/icons-vue'
import AnsibleLogDialog from '@/views/task/Job/AnsibleLogDialog.vue'

export default {
  components: {
    AnsibleLogDialog
  },
  data() {
    return {
      starting: false
    }
  },
  props: {
    steps: {
      type: Array,
      required: true,
      default: () => [],
      validator: (value) => {
        return Array.isArray(value) && 
          value.every(step => typeof step === 'object' && step !== null);
      }
    },
    taskId: {
      type: [String, Number],
      required: true,
      default: null,
      validator: (value) => {
        return value !== null && value !== undefined && value !== ''
      }
    },
    historyMode: {
      type: Boolean,
      default: false
    },
    historyId: {
      type: [String, Number],
      default: null
    }
  },
  watch: {
    steps: {
      handler(newVal) {
        console.log('Ansible Steps prop changed - raw:', JSON.parse(JSON.stringify(newVal)));
        console.log('Current orderedSteps:', JSON.parse(JSON.stringify(this.orderedSteps)));
        this.$nextTick(() => {
          console.log('After nextTick - orderedSteps:', JSON.parse(JSON.stringify(this.orderedSteps)));
          this.drawConnectors();
        });
      },
      deep: true,
      immediate: true
    }
  },
  computed: {
    orderedSteps() {
      if (!this.steps || !Array.isArray(this.steps)) {
        console.warn('Invalid Ansible steps data:', this.steps);
        return [];
      }
      
      console.log('原始Ansible steps数据:', JSON.parse(JSON.stringify(this.steps)));
      
      const result = this.steps.map((step, index) => {
        // Handle both integer (1-4) and string status from backend
        const statusValue = step.status !== undefined ? step.status : 1;
        let status = 'pending';
        
        if (typeof statusValue === 'number') {
           status = ['pending', 'active', 'completed', 'error'][statusValue - 1] || 'pending';
        } else if (typeof statusValue === 'string') {
           const lower = statusValue.toLowerCase();
           if (lower === 'success' || lower === 'completed') status = 'completed';
           else if (lower === 'running' || lower === 'active') status = 'active';
           else if (lower === 'failed' || lower === 'error') status = 'error';
           else status = lower; // fallback
        }
        
        console.log(`Ansible步骤 ${step.work_id || index} 状态映射:`, 
                   `原始: ${step.status}, 最终: ${status}`);
        
        return {
          ...step,
          title: step.entry_file_name || `Work ${index + 1}`,
          content: step.entry_file_name || '',
          order: index + 1,
          status: status,
          _status: status,
          // 确保保留原始status
          original_status: step.status
        };
      }).sort((a, b) => a.order - b.order);
      
      console.log('Ansible orderedSteps计算结果:', JSON.parse(JSON.stringify(result)));
      return result;
    }
  },
  methods: {
    getStepClass(step) {
      return `status-${step.status}`;
    },
    async startAnsibleTask() {
      if (this.historyMode) {
        this.$message.warning('历史记录模式下无法启动任务');
        return;
      }
      this.starting = true;
      try {
        console.log('Starting Ansible task with taskId:', this.taskId);
        
        if (!this.taskId) {
          throw new Error('Ansible任务ID未定义');
        }
        
        console.log('Using Ansible task ID:', this.taskId);
        
        await StartAnsibleTaskFlow(this.taskId);
        this.$message.success('Ansible任务已启动');
        this.$emit('refresh');
      } catch (error) {
        console.error('启动Ansible任务失败:', error);
        this.$message.error('启动Ansible任务失败: ' + (error.message || '未知错误'));
      } finally {
        this.starting = false;
      }
    },

    async handleRefresh() {
      try {
        console.log('开始刷新Ansible任务状态...');
        
        let updatedSteps = [];

        if (this.historyMode && this.historyId) {
           console.log('刷新历史记录模式...', 'taskId:', this.taskId, 'historyId:', this.historyId);
           const historyResponse = await GetAnsibleHistoryDetail({
               id: this.taskId,
               historyId: this.historyId
           });
           
           if (!historyResponse?.data?.data) {
               throw new Error('无效的历史详情响应');
           }
           
           const historyData = historyResponse.data.data;
           const works = historyData.WorkHistories || [];
           
           updatedSteps = works.map((work, index) => ({
               task_id: this.taskId,
               work_id: work.WorkID,
               entry_file_name: work.HostName, // Note: in history detail, HostName might be used for entry_file_name
               status: work.Status,
               duration: work.Duration,
               order: index + 1
           }));

        } else {
            console.log('使用taskId:', this.taskId);
            
            const response = await GetAnsibleTaskDetail(this.taskId);
            console.log('Ansible API完整响应:', response);
            
            if (!response?.data?.data?.task_info) {
              throw new Error('无效的API响应');
            }
            
            // 深度拷贝API数据
            const apiData = JSON.parse(JSON.stringify(response.data.data.task_info));
            console.log('原始Ansible API数据:', apiData);
            
            updatedSteps = (apiData.Works || []).map((work, index) => {
              // 保留原始status
              const statusValue = work.status ?? 1;
              
              return {
                task_id: this.taskId,
                work_id: work.workid,
                entry_file_name: work.EntryFileName || `Work ${index + 1}`,
                status: statusValue, // 保留原始值
                duration: work.Duration || 0,
                order: index + 1
              };
            });
        }
        
        console.log('更新后的Ansible步骤数据:', JSON.parse(JSON.stringify(updatedSteps)));
        
        // 确保emit的数据结构与父组件期望的一致
        const emitData = updatedSteps.map(step => ({
          ...step,
          // 确保包含所有必要字段
          task_id: step.task_id,
          work_id: step.work_id,
          entry_file_name: step.entry_file_name,
          status: step.status,
          duration: step.duration
        }));
        
        console.log('准备emit的Ansible数据:', JSON.parse(JSON.stringify(emitData)));
        this.$emit('update:steps', emitData);
        this.$emit('update-steps', emitData);
        
        // 直接emit更新父组件数据
        this.$emit('update:steps', [...updatedSteps]);
        this.$emit('update-steps', [...updatedSteps]);
        this.$emit('refresh', {
          steps: [...updatedSteps],
          taskId: this.taskId
        });
        
        // 检查更新后的状态
        this.$nextTick(() => {
          console.log('更新后的Ansible orderedSteps:', this.orderedSteps);
          console.log('当前Ansible卡片状态:', 
            this.orderedSteps.map(s => `${s.entry_file_name}: ${s.status}`).join(', '));
        });
        
        this.$message.success('Ansible任务状态已刷新');
      } catch (error) {
        console.error('刷新Ansible任务状态失败:', error);
        this.$message.error('刷新Ansible任务状态失败: ' + (error.message || '未知错误'));
      }
    },

    async onLogClick(step) {
      try {
        console.log('Ansible Step data:', step);
        console.log('详细的step对象结构:', JSON.stringify(step, null, 2));
        console.log('step.task_id:', step.task_id, 'typeof:', typeof step.task_id);
        console.log('step.work_id:', step.work_id, 'typeof:', typeof step.work_id);
        console.log('step.workid:', step.workid, 'typeof:', typeof step.workid);
        
        if (!step.task_id || (!step.work_id && !step.workid)) {
          console.error('缺少任务ID或WorkID，当前step:', step);
          console.error('可用的ID字段:', {
            task_id: step.task_id,
            work_id: step.work_id,
            workid: step.workid,
            id: step.id
          });
          throw new Error('缺少任务ID或WorkID，请确保步骤数据包含task_id和work_id');
        }
        
        // 立即显示对话框，提升用户体验
        console.log('AnsibleLogDialog ref:', this.$refs.ansibleLogDialog);
        if (this.$refs.ansibleLogDialog) {
          console.log('Calling show method with loading state');
          
          const isRunningTask = step.status === 'active' || step.status === 2;
          
          console.log('任务状态检查:', {
            stepStatus: step.status,
            statusType: typeof step.status,
            isRunningTask: isRunningTask,
            stepData: step
          });
          
          // 由于SSE不可用，优先使用定时刷新模式
          console.log('SSE不可用，优先尝试定时刷新模式');
          
          const workId = step.work_id || step.workid;
          console.log('使用的workId:', workId, '来源:', step.work_id ? 'work_id' : 'workid');
          
          this.$refs.ansibleLogDialog.show({
            taskId: step.task_id,
            workId: workId,
            fileName: step.entry_file_name
          });
          
          console.log('WebSocket日志对话框已显示');
        } else {
          console.error('AnsibleLogDialog ref not found');
        }
      } catch (error) {
        console.error('获取Ansible任务日志失败:', error);
        
        // 如果是超时错误，给出更友好的提示
        const isTimeoutError = error.message && error.message.includes('timeout');
        const isRunningTask = step.status === 'active' || step.status === 2;
        
        if (isTimeoutError && isRunningTask) {
          console.log('任务运行中的超时，切换到实时模式');
          this.$message.warning('网络响应超时，将自动切换到实时日志模式');
        } else if (isTimeoutError) {
          console.log('日志获取超时，可能是网络波动');
          this.$message.warning('日志获取超时，请稍后重试');
        } else {
          // WebSocket模式下错误处理由WebSocket事件处理器自动处理
          console.log('WebSocket模式下的错误将由连接处理器自动显示');
          this.$message.error(`获取Ansible任务日志失败: ${error.message || '未知错误'}`);
        }
      }
    },

    formatDuration(seconds) {
      if (!seconds || seconds === 0) return '0';
      return seconds.toString();
    },

    drawConnectors() {
      const svg = this.$refs.connectorCanvas;
      if (!svg || !this.$el) return;
      
      if (this.drawTimeout) clearTimeout(this.drawTimeout);
      
      this.drawTimeout = setTimeout(() => {
        // 按任务顺序获取卡片
        const cards = {};
        const steps = Array.from(this.$el.querySelectorAll('.flow-step'));
        
        steps.forEach(step => {
          const card = step.querySelector('.flow-card');
          if (!card) return;
          const rect = card.getBoundingClientRect();
          if (rect.width > 0 && rect.height > 0) {
            const order = parseInt(step.style.getPropertyValue('--step-order'));
            cards[order] = {
              step,
              rect
            };
          }
        });
        
        const containerRect = this.$el.getBoundingClientRect();
        svg.innerHTML = '';
        const fragment = document.createDocumentFragment();
        
        // 按任务顺序连接卡片
        const sortedCards = Object.values(cards).sort((a, b) => {
          const aOrder = parseInt(a.step.style.getPropertyValue('--step-order'));
          const bOrder = parseInt(b.step.style.getPropertyValue('--step-order'));
          return aOrder - bOrder;
        });
        
        // 连接所有任务
        for (let i = 0; i < sortedCards.length; i++) {
          const current = sortedCards[i];
          const currentOrder = parseInt(current.step.style.getPropertyValue('--step-order'));
          
          // 处理普通水平连接
          if (i < sortedCards.length - 1) {
            const next = sortedCards[i+1];
            const nextOrder = parseInt(next.step.style.getPropertyValue('--step-order'));
            if (nextOrder === currentOrder + 1) {
              const startX = current.rect.right - containerRect.left;
              const startY = current.rect.top + current.rect.height / 2 - containerRect.top;
              const endX = next.rect.left - containerRect.left;
              const endY = next.rect.top + next.rect.height / 2 - containerRect.top;
              
              const path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
              path.setAttribute('d', `M${startX},${startY} L${endX},${endY}`);
              path.setAttribute('class', 'flow-path');
              path.setAttribute('stroke', current.step.status === 'completed' ? 'rgba(76,175,80,0.5)' : 
                  current.step.status === 'active' ? 'rgba(33,150,243,0.5)' : 'rgba(128,128,128,0.5)');
              path.setAttribute('stroke-width', '2');
              path.setAttribute('stroke-dasharray', '8,4');
              path.setAttribute('marker-end', `url(#arrowhead-${currentOrder})`);
              path.setAttribute('fill', 'none');
              
              // 添加流动箭头
              const arrow = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
              arrow.setAttribute('cx', '0');
              arrow.setAttribute('cy', '0');
              arrow.setAttribute('r', '3');
              arrow.setAttribute('fill', current.step.status === 'completed' ? 'rgba(76,175,80,0.7)' : 
                  current.step.status === 'active' ? 'rgba(33,150,243,0.7)' : 'rgba(128,128,128,0.7)');
              arrow.setAttribute('class', 'flow-arrow');
              
              const animate = document.createElementNS('http://www.w3.org/2000/svg', 'animateMotion');
              animate.setAttribute('path', `M${startX},${startY} L${endX},${endY}`);
              animate.setAttribute('dur', '3s');
              animate.setAttribute('repeatCount', 'indefinite');
              
              arrow.appendChild(animate);
              fragment.appendChild(path);
              fragment.appendChild(arrow);
            }
          }
        }
        
        svg.appendChild(fragment);
      }, 150);
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.drawConnectors();
      this.observer = new ResizeObserver(() => {
        this.drawConnectors();
      });
      this.observer.observe(this.$el);
      
      this.mutationObserver = new MutationObserver((mutations) => {
        mutations.forEach(() => this.drawConnectors());
      });
      this.mutationObserver.observe(this.$el, {
        childList: true,
        subtree: true,
        attributes: true
      });
    });
  },
  
  updated() {
    this.$nextTick(() => this.drawConnectors());
  },
  
  beforeUnmount() {
    console.log('Ansible Component destroying...');
    this.observer?.disconnect();
    this.mutationObserver?.disconnect();
    if (this.drawTimeout) clearTimeout(this.drawTimeout);
    console.log('Ansible Component destroyed');
  }
}
</script>

<style scoped>
div.ansible-flow-container {
  display: flex !important;
  flex-wrap: wrap !important;
  gap: 3px !important;
  padding: 0 !important;
  margin: -5px 0 0 0 !important;
  position: relative !important;
  top: -5px !important;
  left: 0 !important;
}

.flow-header {
  width: 100%;
  padding: 10px;
  display: flex;
  gap: 10px;
}

.connector-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 10;
}

.flow-path {
  stroke: rgba(128,128,128,0.5);
  stroke-width: 2;
  stroke-dasharray: 10;
  stroke-dashoffset: 10;
  animation: flow 3s linear infinite;
}

@keyframes flow {
  to {
    stroke-dashoffset: 0;
  }
}

.flow-arrow {
  fill: rgba(128,128,128,0.7);
  animation: arrowFlow 3s linear infinite;
  z-index: 12;
  pointer-events: none;
  transform-box: fill-box;
  transform-origin: center;
  position: relative;
  top: 2px;
}

@keyframes arrowFlow {
  0% {
    opacity: 0;
    transform: translateX(-20px);
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0;
    transform: translateX(20px);
  }
}

.flow-step {
  width: calc(33.33% - 14px);
  order: var(--step-order);
}

.flow-card {
  position: relative;
  padding: 10px;
  width: 240px;
  min-height: 120px;
  border-radius: 8px;
  background: linear-gradient(to bottom right, #E6522C 0%, #FFE0B2 100%) !important;
  box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.4);
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  color: #333;
  border: 1px solid rgba(255, 255, 255, 0.3);
  transform: translateZ(0);
  will-change: transform, box-shadow;
  overflow: hidden;
  z-index: 5;
}

.flow-card.status-completed {
  background: linear-gradient(to bottom right, #4CAF50 0%, #C8E6C9 100%) !important;
}

.flow-card.status-active {
  background: linear-gradient(to bottom right, #2196F3 0%, #BBDEFB 100%) !important;
}

.flow-card.status-pending {
  background: linear-gradient(to bottom right, #FF9800 0%, #FFE0B2 100%) !important;
}

.flow-card.status-error {
  background: linear-gradient(to bottom right, #F44336 0%, #FFCDD2 100%) !important;
}

.flow-card:hover {
  transform: translateY(-8px) scale(1.03);
  box-shadow: 0 8px 24px 0 rgba(0, 0, 0, 0.5);
  z-index: 10;
  filter: brightness(1.05);
}

.status-text {
  position: absolute;
  top: 10px;
  right: 10px;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
}

.status-completed {
  background: #f0f9eb;
  color: #67c23a;
}

.status-active {
  background: #ecf5ff;
  color: #409eff;
}

.status-pending {
  background: #fdf6ec;
  color: #e6a23c;
}

.status-error {
  background: #fef0f0;
  color: #f56c6c;
}

.card-title {
  display: flex;
  align-items: center;
  margin: 0 0 10px;
  font-size: 16px;
  color: #1A237E;
  font-weight: 500;
}

.title-icon {
  margin-right: 8px;
}

.card-content {
  margin: 0 0 10px;
  font-size: 14px;
  color: #333;
}

.card-duration {
  margin: 5px 0;
  font-size: 12px;
  color: #666;
}

.card-buttons {
  position: absolute;
  bottom: 10px;
  right: 10px;
  display: flex;
  background: transparent;
  padding: 0;
  border-radius: 4px;
}

.btn-icon {
  width: 28px;
  height: 28px;
  padding: 4px;
  border-radius: 4px;
  background: rgba(255,255,255,0.9);
  border: none;
  cursor: pointer;
  transition: all 0.2s;
  margin-left: 4px;
}

.btn-icon:hover {
  background: #e4e7ed;
}
</style>