<template>
  <div class="flow-container" style="margin-top: -5px !important; top: -5px !important;">
    <div class="flow-header">
      <el-button
        @click="startAllSteps"
        v-authority="['task:job:jobstart']"
        type="primary"
        size="small"
        style="height: 32px; padding: 0 16px;"
        :loading="starting"
        :disabled="hasStarted || starting"
        :class="{ 'started-button': hasStarted }"
      >
        <el-icon size="16"><VideoPlay /></el-icon>
        <span style="margin-left: 6px;">{{ hasStarted ? '已启动' : '启动脚本' }}</span>
      </el-button>
      <el-button
        @click="handleRefresh"
        type="warning"
        size="small"
        style="height: 32px; padding: 0 16px;"
        :loading="refreshing"
      >
        <el-icon size="16"><Refresh /></el-icon>
        <span style="margin-left: 6px;">刷新状态</span>
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
<path fill="#9C27B0" d="M512 1024C229.23 1024 0 794.77 0 512S229.23 0 512 0s512 229.23 512 512-229.23 512-512 512zm0-938.67C276.36 85.33 85.33 276.36 85.33 512S276.36 938.67 512 938.67 938.67 747.64 938.67 512 747.64 85.33 512 85.33z"/>
<path fill="#9C27B0" d="M341.33 426.67a42.67 42.67 0 100-85.34 42.67 42.67 0 000 85.34zM682.67 426.67a42.67 42.67 0 100-85.34 42.67 42.67 0 000 85.34z"/>
<path fill="#9C27B0" d="M512 768c-117.82 0-213.33-95.51-213.33-213.33h426.66C725.33 672.49 629.82 768 512 768z"/>
          </svg>
          {{ step.title }}
        </h3>
        <p class="card-content">{{ step.content }}</p>
        <div class="card-buttons">
          <el-tooltip effect="dark" content="脚本" placement="top">
            <button @click.stop="onDetailClick(step)" class="btn-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M14 2H6C4.9 2 4 2.9 4 4V20C4 21.1 4.9 22 6 22H18C19.1 22 20 21.1 20 20V8L14 2Z" fill="#1890ff"/>
                <path d="M14 2V8H20" fill="#1890ff" fill-opacity="0.5"/>
                <path d="M16 13H8V11H16V13ZM16 17H8V15H16V17ZM13 9V3.5L18.5 9H13Z" fill="white"/>
              </svg>
            </button>
          </el-tooltip>
          <el-tooltip effect="dark" content="日志" placement="top">
            <button @click.stop="onLogClick(step)" class="btn-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M21 10H7M21 6H3M21 14H3M21 18H7" stroke="#722ed1" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </button>
          </el-tooltip>
          <el-tooltip effect="dark" content="停止" placement="top">
            <button v-authority="['task:job:jobstop']" @click.stop="onStopClick(step)" class="btn-icon" :loading="stopping">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <circle cx="12" cy="12" r="8" fill="#f5222d"/>
              </svg>
            </button>
          </el-tooltip>
        </div>
      </div>
    </div>
    <div style="position: fixed; z-index: 9999;">
      <script-dialog ref="scriptDialog"></script-dialog>
      <log-dialog ref="logDialog"></log-dialog>
    </div>
  </div>
</template>

<script>
import { GetTemplateContent, GetTaskJobLog, StartJob, StopJob, GetTaskTemplates } from '@/api/task'
import { VideoPlay, Refresh } from '@element-plus/icons-vue'
import ScriptDialog from '@/views/task/Job/ScriptDialog.vue'
import LogDialog from '@/views/task/Job/LogDialog.vue'

export default {
    components: {
      ScriptDialog,
      LogDialog
    },
    data() {
      return {
        starting: false,
        stopping: false,
        hasStarted: false, // 追踪是否已经启动过
        refreshing: false,  // 追踪刷新状态
        forceUpdateKey: 0,  // 强制更新key
        localSteps: []      // 本地步骤数据副本
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
    }
  },
  watch: {
    steps: {
      handler(newVal, oldVal) {
        console.log('Steps prop changed - raw:', JSON.parse(JSON.stringify(newVal)));
        console.log('Current orderedSteps:', JSON.parse(JSON.stringify(this.orderedSteps)));

        // 更新本地数据副本
        if (newVal && Array.isArray(newVal)) {
          this.localSteps = [...newVal];
          this.forceUpdateKey += 1;
        }

        // 检查是否有任务在运行，如果有则设置已启动状态
        this.checkTaskStatus(newVal);

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
      // 使用forceUpdateKey来强制重新计算
      const updateKey = this.forceUpdateKey;

      // 优先使用localSteps，如果为空则使用props.steps
      const currentSteps = this.localSteps.length > 0 ? this.localSteps : this.steps;

      if (!currentSteps || !Array.isArray(currentSteps)) {
        console.warn('Invalid steps data:', currentSteps);
        return [];
      }

      console.log('原始steps数据 (forceUpdateKey=' + updateKey + '):', JSON.parse(JSON.stringify(currentSteps)));

      const result = currentSteps.map((step, index) => {
        // 直接使用API返回的task_status，不进行任何转换
        const taskStatus = step.task_status !== undefined ? step.task_status : 1;
        const status = ['pending', 'active', 'completed', 'error'][taskStatus - 1] || 'pending';

        console.log(`步骤 ${step.template_id || index} 状态映射:`,
                   `原始: ${step.task_status}, 最终: ${status}`);

        return {
          ...step,
          title: step.template_name || '',
          content: step.template_remark || '',
          order: step.order || index + 1,
          status: status,
          _status: status,
          // 确保保留原始task_status
          task_status: step.task_status
        };
      }).sort((a, b) => a.order - b.order);

      console.log('orderedSteps计算结果:', JSON.parse(JSON.stringify(result)));
      return result;
    }
  },
  methods: {
    // 检查任务状态，判断是否已经启动
    checkTaskStatus(steps) {
      if (!steps || !Array.isArray(steps)) return;

      // 如果有任务状态不是 pending（等待），说明已经启动过
      const hasStartedTask = steps.some(step => {
        const taskStatus = step.task_status ?? 1;
        return taskStatus > 1; // 1=pending, 2=active, 3=completed, 4=error
      });

      this.hasStarted = hasStartedTask;
      console.log('Task status check:', { hasStartedTask, hasStarted: this.hasStarted });
    },

    getStepClass(step) {
      return `status-${step.status}`;
    },
    async startAllSteps() {
      this.starting = true;
      try {
        console.log('Starting task with props taskId:', this.taskId);
        console.log('Steps data:', this.steps);
        
        if (!this.taskId && (!this.steps || this.steps.length === 0)) {
          throw new Error('任务ID未定义');
        }
        
        // 优先使用steps中的task_id
        const apiTaskId = this.steps[0]?.task_id || this.taskId;
        if (!apiTaskId) {
          throw new Error('无法获取有效的任务ID');
        }
        
        console.log('Using task ID:', apiTaskId);
        
        // 确保参数格式正确
        const params = new URLSearchParams();
        params.append('taskId', apiTaskId);
        
        await StartJob(
          { taskId: apiTaskId },
          { params: params }
        );

        // 设置已启动状态
        this.hasStarted = true;

        this.$message.success('任务已启动');
        this.$emit('refresh');
      } catch (error) {
        console.error('启动任务失败:', error);
        this.$message.error('启动任务失败: ' + (error.message || '未知错误'));
      } finally {
        this.starting = false;
      }
    },

    async handleRefresh() {
      this.refreshing = true;
      try {
        console.log('开始刷新任务状态...');
        const apiTaskId = this.steps[0]?.task_id || this.taskId;
        console.log('使用taskId:', apiTaskId);

        const response = await GetTaskTemplates({ id: apiTaskId });
        console.log('API完整响应:', response);

        if (!response?.data?.data) {
          throw new Error('无效的API响应');
        }

        // 深度拷贝API数据
        const apiData = JSON.parse(JSON.stringify(response.data.data));
        console.log('原始API数据:', apiData);

        const updatedSteps = apiData.map((step, index) => {
          // 保留原始task_status
          const taskStatus = step.task_status ?? 1;
          const status = ['pending', 'active', 'completed', 'error'][taskStatus - 1] || 'pending';

          console.log(`步骤 ${step.template_id || index} 状态映射:`,
                     `原始: ${step.task_status}, 最终: ${status}`);

          return {
            ...step,
            task_id: step.task_id || this.taskId,
            template_id: step.template_id || index + 1,
            template_name: step.template_name || `步骤 ${index + 1}`,
            template_remark: step.template_remark || '',
            task_status: taskStatus, // 保留原始值
            status: status,
            order: step.order || index + 1
          };
        });

        console.log('更新后的步骤数据:', JSON.parse(JSON.stringify(updatedSteps)));

        // 更新本地数据副本
        this.localSteps = [...updatedSteps];

        // 强制更新computed属性
        this.forceUpdateKey += 1;

        // 强制更新父组件数据
        this.$emit('refresh-steps', updatedSteps);
        this.$emit('update:steps', updatedSteps);

        // 强制重新渲染
        this.$nextTick(() => {
          this.$forceUpdate();
          console.log('强制更新后的orderedSteps:', this.orderedSteps);
          console.log('当前卡片状态:',
            this.orderedSteps.map(s => `${s.template_name}: ${s.status}`).join(', '));

          // 验证状态映射
          updatedSteps.forEach((step, index) => {
            console.log(`验证步骤 ${index + 1}:`, {
              template_name: step.template_name,
              task_status: step.task_status,
              mapped_status: ['pending', 'active', 'completed', 'error'][step.task_status - 1],
              final_status: step.status
            });
          });

          // 重新绘制连接线
          this.drawConnectors();
        });

        this.$message.success('任务状态已刷新');
      } catch (error) {
        console.error('刷新任务状态失败:', error);
        this.$message.error('刷新任务状态失败: ' + (error.message || '未知错误'));
      } finally {
        this.refreshing = false;
      }
    },
    onDeployClick(step) {
      this.$emit('deploy', step);
    },
    async onDetailClick(step) {
      try {
        const response = await GetTemplateContent({ id: step.template_id });
        this.$emit('detail', {
          ...step,
          content: response.data
        });
        
        console.log('ScriptDialog ref:', this.$refs.scriptDialog);
        console.log('ScriptDialog component:', this.$refs.scriptDialog?.$el);
        if (this.$refs.scriptDialog) {
          console.log('Calling show method');
          this.$refs.scriptDialog.show(response.data);
          console.log('Show method called');
        } else {
          console.error('ScriptDialog ref not found');
        }
      } catch (error) {
        console.error('获取模板内容失败:', error);
        this.$message.error('获取模板内容失败');
      }
    },
    async onStopClick(step) {
      this.stopping = true;
      try {
        await StopJob({
          taskId: step.task_id,
          templateId: step.template_id
        });
        this.$message.success('任务已停止');
        this.$emit('refresh');
      } catch (error) {
        console.error('停止任务失败:', error);
        this.$message.error('停止任务失败: ' + (error.message || '未知错误'));
      } finally {
        this.stopping = false;
      }
    },

    async onLogClick(step) {
      try {
        console.log('Step data:', step);
        if (!step.task_id) {
          console.error('缺少任务ID，当前step:', step);
          throw new Error('缺少任务ID，请确保步骤数据包含task_id');
        }
        
        console.log('调用日志API参数:', { 
          taskId: step.task_id, 
          templateId: step.template_id 
        });
        const response = await GetTaskJobLog({
          taskId: step.task_id,
          templateId: step.template_id
        });
        
        console.log('LogDialog ref:', this.$refs.logDialog);
        console.log('LogDialog component:', this.$refs.logDialog?.$el);
        if (this.$refs.logDialog) {
          console.log('Calling show method');
          this.$refs.logDialog.show(response.data);
          console.log('Show method called');
          this.$emit('log', {
            ...step,
            logs: response.data
          });
        } else {
          console.error('LogDialog ref not found');
        }
      } catch (error) {
        console.error('获取任务日志失败:', error);
        this.$message.error(`获取任务日志失败: ${error.message || '未知错误'}`);
      }
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
        
        // 特殊处理3→4和9→10的垂直连接
        const verticalConnections = {
          3: 4,
          9: 10
        };
        
        // 连接所有任务
        for (let i = 0; i < sortedCards.length; i++) {
          const current = sortedCards[i];
          const currentOrder = parseInt(current.step.style.getPropertyValue('--step-order'));
          
          // 处理特殊垂直连接
          if (verticalConnections[currentOrder]) {
            const next = sortedCards.find(card => 
              parseInt(card.step.style.getPropertyValue('--step-order')) === verticalConnections[currentOrder]
            );
            if (next) {
              const startX = current.rect.left + current.rect.width/4 - containerRect.left;
              const startY = current.rect.bottom - containerRect.top;
              const endX = next.rect.right - next.rect.width/4 - containerRect.left;
              const endY = next.rect.top - containerRect.top;
              
              const path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
              // 创建贝塞尔曲线路径
              const controlX1 = startX + (endX - startX) * 0.25;
              const controlY1 = startY;
              const controlX2 = endX - (endX - startX) * 0.25;
              const controlY2 = endY;
              path.setAttribute('d', `M${startX},${startY} C${controlX1},${controlY1} ${controlX2},${controlY2} ${endX},${endY}`);
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
              arrow.setAttribute('class', `flow-arrow status-${current.step.status}`);
              arrow.setAttribute('class', 'flow-arrow');
              
              // 创建动画元素
              const animate = document.createElementNS('http://www.w3.org/2000/svg', 'animateMotion');
              animate.setAttribute('path', `M${startX},${startY} L${endX},${endY}`);
              animate.setAttribute('dur', '3s');
              animate.setAttribute('repeatCount', 'indefinite');
              
              arrow.appendChild(animate);
              fragment.appendChild(path);
              fragment.appendChild(arrow);
              continue;
            }
          }
          
          // 处理普通水平连接
          if (i < sortedCards.length - 1) {
            const next = sortedCards[i+1];
            const nextOrder = parseInt(next.step.style.getPropertyValue('--step-order'));
            if (nextOrder === currentOrder + 1 && !verticalConnections[currentOrder]) {
              // 特殊处理6-7连接
              let startX, startY, endX, endY, path;
              if (currentOrder === 6) {
                startX = current.rect.left + current.rect.width/4 - containerRect.left;
                startY = current.rect.bottom - containerRect.top;
                endX = next.rect.right - next.rect.width/4 - containerRect.left;
                endY = next.rect.top - containerRect.top;
                
                path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
                const controlX1 = startX + (endX - startX) * 0.25;
                const controlY1 = startY;
                const controlX2 = endX - (endX - startX) * 0.25;
                const controlY2 = endY;
                path.setAttribute('d', `M${startX},${startY} C${controlX1},${controlY1} ${controlX2},${controlY2} ${endX},${endY}`);
                path.setAttribute('class', 'flow-path');
                path.setAttribute('stroke', current.step.status === 'completed' ? 'rgba(76,175,80,0.5)' : 
                  current.step.status === 'active' ? 'rgba(33,150,243,0.5)' : 'rgba(128,128,128,0.5)');
                path.setAttribute('stroke-width', '2');
                path.setAttribute('stroke-dasharray', '8,4');
                path.setAttribute('marker-end', `url(#arrowhead-${currentOrder})`);
                path.setAttribute('fill', 'none');
              } else {
                startX = current.rect.right - containerRect.left;
                startY = current.rect.top + current.rect.height / 2 - containerRect.top;
                endX = next.rect.left - containerRect.left;
                endY = next.rect.top + next.rect.height / 2 - containerRect.top;
                
                path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
                path.setAttribute('d', `M${startX},${startY} L${endX},${endY}`);
                path.setAttribute('class', 'flow-path');
                path.setAttribute('stroke', current.step.status === 'completed' ? 'rgba(76,175,80,0.5)' : 
                    current.step.status === 'active' ? 'rgba(33,150,243,0.5)' : 'rgba(128,128,128,0.5)');
                path.setAttribute('stroke-width', '2');
                path.setAttribute('stroke-dasharray', '8,4');
                path.setAttribute('marker-end', `url(#arrowhead-${currentOrder})`);
                path.setAttribute('fill', 'none');
              }
              
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
      console.log('Component destroying...');
      this.observer?.disconnect();
      this.mutationObserver?.disconnect();
      if (this.drawTimeout) clearTimeout(this.drawTimeout);
      console.log('Component destroyed');
    }
}
</script>

<style scoped>
div.flow-container {
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

/* 更具体的嵌套选择器 */
.task-flow-wrapper > div.flow-container,
.task-flow-container > div.flow-container,
#taskFlow > div.flow-container {
  margin-top: -7.5px !important;
  top: -7.5px !important;
}

/* 如果仍然不生效，可以尝试在模板中添加内联样式 */
/* <div class="flow-container" style="margin-top: -7.5px !important; top: -7.5px !important;"> */

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
}

@keyframes arrowFlow {
  0% {
    opacity: 0;
    transform: translateX(-10px);
  }
  30% {
    opacity: 1;
  }
  70% {
    opacity: 1;
  }
  100% {
    opacity: 0;
    transform: translateX(10px);
  }
}

.flow-arrow {
  z-index: 12;
  animation: arrowFlow 3s linear infinite;
  pointer-events: none;
  transform-box: fill-box;
  transform-origin: center;
  position: relative;
  top: 2px;
}

@keyframes flow {
  to {
    stroke-dashoffset: 0;
  }
}

.flow-arrow {
  animation: arrowFlow 1s linear infinite;
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
  min-height: 110px;
  border-radius: 8px;
  background: linear-gradient(to bottom right, #1E88E5 0%, #B3E5FC 100%) !important;
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

.flow-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(30, 136, 229, 0.8) 0%, transparent 50%);
  z-index: -1;
}

.flow-card:hover {
  transform: translateY(-8px) scale(1.03);
  box-shadow: 0 8px 24px 0 rgba(0, 0, 0, 0.5);
  z-index: 10;
  filter: brightness(1.05);
}

.flow-card.status-completed:hover {
  background: linear-gradient(to bottom right, #388E3C 0%, #A5D6A7 100%) !important;
}

.flow-card.status-active:hover {
  background: linear-gradient(to bottom right, #1976D2 0%, #90CAF9 100%) !important;
}

.flow-card.status-pending:hover {
  background: linear-gradient(to bottom right, #F57C00 0%, #FFCC80 100%) !important;
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
  margin: 0;
  font-size: 14px;
  color: #333;
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

/* 启动按钮已启动状态样式 */
.started-button {
  background-color: #909399 !important;
  border-color: #909399 !important;
  cursor: not-allowed !important;
}

.started-button:hover {
  background-color: #909399 !important;
  border-color: #909399 !important;
}

</style>
