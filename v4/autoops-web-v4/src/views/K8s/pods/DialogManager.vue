<template>
  <!-- 工作负载详情对话框 -->
  <WorkloadDetailDialog
    v-model="dialogs.workloadDetail.visible"
    :workload="dialogs.workloadDetail.data"
    @close="closeDialog('workloadDetail')"
  />

  <!-- Pod 列表对话框 -->
  <PodListDialog
    v-model="dialogs.podList.visible"
    :workload-name="dialogs.podList.data?.workloadName || ''"
    :pod-list="dialogs.podList.data?.pods || []"
    :loading="dialogs.podList.data?.loading || false"
    @close="closeDialog('podList')"
    @refresh="handlePodListRefresh"
    @view-detail="handleViewPodDetail"
    @view-logs="handleViewPodLogs"
    @view-yaml="handleViewPodYaml"
    @delete-pod="handleDeletePod"
  />

  <!-- 日志查看对话框 -->
  <LogViewerDialog
    v-model="dialogs.logs.visible"
    :pod-name="dialogs.logs.data?.podName || ''"
    :namespace="dialogs.logs.data?.namespace || ''"
    :containers="dialogs.logs.data?.containers || []"
    :logs="dialogs.logs.data?.logs || ''"
    :loading="dialogs.logs.data?.loading || false"
    @close="closeDialog('logs')"
    @refresh-logs="handleRefreshLogs"
    @start-follow="handleStartFollowLogs"
    @stop-follow="handleStopFollowLogs"
  />

  <!-- YAML 查看对话框 -->
  <YamlViewerDialog
    v-model="dialogs.yaml.visible"
    :title="dialogs.yaml.data?.title || 'YAML查看'"
    :yaml-content="dialogs.yaml.data?.yamlContent || ''"
    :resource-type="dialogs.yaml.data?.resourceType || ''"
    :resource-name="dialogs.yaml.data?.resourceName || ''"
    :namespace="dialogs.yaml.data?.namespace || ''"
    :loading="dialogs.yaml.data?.loading || false"
    :editable="dialogs.yaml.data?.editable || false"
    @close="closeDialog('yaml')"
    @save="handleSaveYaml"
    @validate="handleValidateYaml"
  />

  <!-- 扩缩容对话框 -->
  <ScaleDialog
    v-model="dialogs.scale.visible"
    :workload="dialogs.scale.data?.workload || {}"
    :scaling="dialogs.scale.data?.scaling || false"
    @close="closeDialog('scale')"
    @scale="handleScale"
  />
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'
import WorkloadDetailDialog from './WorkloadDetailDialog.vue'
import PodListDialog from './PodListDialog.vue'
import LogViewerDialog from './LogViewerDialog.vue'
import YamlViewerDialog from './YamlViewerDialog.vue'
import ScaleDialog from './ScaleDialog.vue'
import { useK8sGlobalState, useEventBus, K8S_EVENTS } from './eventBus.js'

const props = defineProps({
  clusterId: {
    type: String,
    required: true
  },
  namespace: {
    type: String,
    required: true
  }
})

const emit = defineEmits([
  'refresh-workloads',
  'refresh-pod-list',
  'get-pod-logs',
  'get-yaml',
  'save-yaml',
  'validate-yaml',
  'scale-workload',
  'delete-pod'
])

const { state: globalState, closeDialog } = useK8sGlobalState()
const { emit: emitEvent, events } = useEventBus()

// 获取对话框状态
const dialogs = globalState.dialogs

// Pod 列表相关处理
const handlePodListRefresh = () => {
  const workload = dialogs.podList.data?.workload
  if (workload) {
    emit('refresh-pod-list', workload)
  }
}

const handleViewPodDetail = (pod) => {
  // 这里可以打开 Pod 详情对话框
  // 目前先触发事件
  emitEvent(events.POD_STATUS_CHANGED, { pod, action: 'view-detail' })
}

const handleViewPodLogs = (pod) => {
  emit('get-pod-logs', {
    clusterId: props.clusterId,
    namespace: props.namespace,
    podName: pod.name,
    containers: pod.containers || []
  })
}

const handleViewPodYaml = (pod) => {
  emit('get-yaml', {
    clusterId: props.clusterId,
    namespace: props.namespace,
    resourceType: 'pod',
    resourceName: pod.name
  })
}

const handleDeletePod = (pod) => {
  emit('delete-pod', {
    clusterId: props.clusterId,
    namespace: props.namespace,
    podName: pod.name,
    workload: dialogs.podList.data?.workload
  })
}

// 日志相关处理
const handleRefreshLogs = (params) => {
  const logData = dialogs.logs.data
  if (logData) {
    emit('get-pod-logs', {
      clusterId: props.clusterId,
      namespace: props.namespace,
      podName: logData.podName,
      container: params.container,
      lines: params.lines
    })
  }
}

const handleStartFollowLogs = () => {
  // 开始跟踪日志
  emitEvent(events.POD_LOGS_UPDATED, { action: 'start-follow' })
}

const handleStopFollowLogs = () => {
  // 停止跟踪日志
  emitEvent(events.POD_LOGS_UPDATED, { action: 'stop-follow' })
}

// YAML 相关处理
const handleSaveYaml = (yamlData) => {
  emit('save-yaml', {
    clusterId: props.clusterId,
    ...yamlData
  })
}

const handleValidateYaml = (yamlContent) => {
  emit('validate-yaml', {
    clusterId: props.clusterId,
    yamlContent
  })
}

// 扩缩容处理
const handleScale = (scaleData) => {
  emit('scale-workload', {
    clusterId: props.clusterId,
    namespace: props.namespace,
    ...scaleData
  })
}
</script>

<style scoped>
/* 这个组件主要用于管理对话框，不需要特殊样式 */
</style>