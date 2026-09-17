<script setup>
import { ref } from 'vue'
import { Panel, VueFlow, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import ProcessNode from '../K8s/test/ProcessNode.vue'
import { useRunProcess } from '../K8s/test/useRunProcess'
import { useLayout } from '../K8s/test/useLayout'

// 定义CICD流程节点
const nodes = ref([
  {
    id: '1',
    type: 'process',
    position: { x: 0, y: 0 },
    data: { 
      status: null,
      label: '拉取代码',
      logs: '从Git仓库拉取最新代码...'
    }
  },
  {
    id: '2',
    type: 'process',
    position: { x: 0, y: 0 },
    data: { 
      status: null,
      label: '编译项目',
      logs: '执行编译命令，生成构建产物...'
    }
  },
  {
    id: '3',
    type: 'process',
    position: { x: 0, y: 0 },
    data: { 
      status: null,
      label: '切换流量',
      logs: '将流量切换到新版本...'
    }
  },
  {
    id: '4',
    type: 'process',
    position: { x: 0, y: 0 },
    data: { 
      status: null,
      label: '部署服务',
      logs: '部署新版本服务到生产环境...'
    }
  },
  {
    id: '5',
    type: 'process',
    position: { x: 0, y: 0 },
    data: { 
      status: null,
      label: '流量上线',
      logs: '确认新版本服务正常运行...'
    }
  }
])

// 定义流程连接
const edges = ref([
  { id: 'e1-2', source: '1', target: '2' },
  { id: 'e2-3', source: '2', target: '3' },
  { id: 'e3-4', source: '3', target: '4' },
  { id: 'e4-5', source: '4', target: '5' }
])

const selectedNode = ref(null)
const showLogs = ref(false)

const { graph, layout } = useLayout()
const { run, stop, reset, isRunning } = useRunProcess({ graph, cancelOnError: true })
const { fitView } = useVueFlow()

// 布局函数
async function layoutGraph(direction) {
  await stop()
  reset(nodes.value)
  nodes.value = layout(nodes.value, edges.value, direction)
  nextTick(() => {
    fitView()
  })
}

// 节点点击事件
function onNodeClick(node) {
  selectedNode.value = node
  showLogs.value = true
}
</script>

<template>
  <div class="cicd-flow">
    <VueFlow
      v-model:nodes="nodes"
      v-model:edges="edges"
      :default-edge-options="{ type: 'animation', animated: true }"
      @nodes-initialized="layoutGraph('LR')"
      @node-click="onNodeClick"
    >
      <template #node-process="props">
        <ProcessNode 
          :data="props.data" 
          :source-position="props.sourcePosition" 
          :target-position="props.targetPosition" 
        />
      </template>

      <Background />

      <Panel class="process-panel" position="top-right">
        <div class="layout-panel">
          <button v-if="isRunning" class="stop-btn" title="停止" @click="stop">
            <span class="spinner" />
            停止
          </button>
          <button v-else title="开始" @click="run(nodes)">
            开始
          </button>

          <button title="水平布局" @click="layoutGraph('LR')">
            水平
          </button>

          <button title="垂直布局" @click="layoutGraph('TB')">
            垂直
          </button>
        </div>
      </Panel>
    </VueFlow>

    <!-- 日志弹窗 -->
    <div v-if="showLogs" class="log-modal" @click.self="showLogs = false">
      <div class="log-content">
        <h3>{{ selectedNode?.data?.label }}</h3>
        <div class="log-text">
          {{ selectedNode?.data?.logs }}
        </div>
        <button class="close-btn" @click="showLogs = false">关闭</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cicd-flow {
  height: 100%;
  width: 100%;
  position: relative;
}

.process-panel {
  background-color: #2d3748;
  padding: 10px;
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
}

.process-panel button {
  margin: 0 5px;
  padding: 5px 10px;
  background-color: #4a5568;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.process-panel button:hover {
  background-color: #2563eb;
}

.log-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.log-content {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  max-width: 600px;
  width: 80%;
}

.log-text {
  margin: 15px 0;
  padding: 10px;
  background-color: #f5f5f5;
  border-radius: 4px;
  max-height: 300px;
  overflow-y: auto;
}

.close-btn {
  padding: 5px 15px;
  background-color: #2563eb;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
</style>
</template>
