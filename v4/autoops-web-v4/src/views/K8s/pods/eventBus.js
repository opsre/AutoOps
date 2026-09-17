import { reactive } from 'vue'

// 事件总线实现
class EventBus {
  constructor() {
    this.events = {}
  }

  // 订阅事件
  on(event, callback) {
    if (!this.events[event]) {
      this.events[event] = []
    }
    this.events[event].push(callback)

    // 返回取消订阅函数
    return () => this.off(event, callback)
  }

  // 取消订阅
  off(event, callback) {
    if (!this.events[event]) return

    if (callback) {
      const index = this.events[event].indexOf(callback)
      if (index > -1) {
        this.events[event].splice(index, 1)
      }
    } else {
      delete this.events[event]
    }
  }

  // 触发事件
  emit(event, ...args) {
    if (!this.events[event]) return

    this.events[event].forEach(callback => {
      try {
        callback(...args)
      } catch (error) {
        console.error(`Event handler error for "${event}":`, error)
      }
    })
  }

  // 一次性订阅
  once(event, callback) {
    const onceCallback = (...args) => {
      callback(...args)
      this.off(event, onceCallback)
    }
    return this.on(event, onceCallback)
  }

  // 清空所有事件
  clear() {
    this.events = {}
  }
}

// 创建全局事件总线实例
export const k8sEventBus = new EventBus()

// K8s 相关事件定义
export const K8S_EVENTS = {
  // 集群相关
  CLUSTER_CHANGED: 'cluster:changed',
  CLUSTER_STATUS_UPDATED: 'cluster:status_updated',

  // 命名空间相关
  NAMESPACE_CHANGED: 'namespace:changed',
  NAMESPACE_LIST_UPDATED: 'namespace:list_updated',

  // 工作负载相关
  WORKLOAD_CREATED: 'workload:created',
  WORKLOAD_UPDATED: 'workload:updated',
  WORKLOAD_DELETED: 'workload:deleted',
  WORKLOAD_SCALED: 'workload:scaled',
  WORKLOAD_RESTARTED: 'workload:restarted',
  WORKLOAD_LIST_REFRESH: 'workload:list_refresh',

  // Pod 相关
  POD_CREATED: 'pod:created',
  POD_UPDATED: 'pod:updated',
  POD_DELETED: 'pod:deleted',
  POD_STATUS_CHANGED: 'pod:status_changed',
  POD_LOGS_UPDATED: 'pod:logs_updated',

  // 对话框相关
  DIALOG_OPEN: 'dialog:open',
  DIALOG_CLOSE: 'dialog:close',

  // 通用操作
  REFRESH_DATA: 'data:refresh',
  LOADING_START: 'loading:start',
  LOADING_END: 'loading:end',
  ERROR_OCCURRED: 'error:occurred',
  SUCCESS_OCCURRED: 'success:occurred'
}

// 事件数据类型定义
export const createEventData = {
  // 集群变化事件数据
  clusterChanged: (clusterId, clusterInfo) => ({
    clusterId,
    clusterInfo,
    timestamp: Date.now()
  }),

  // 命名空间变化事件数据
  namespaceChanged: (namespace, clusterId) => ({
    namespace,
    clusterId,
    timestamp: Date.now()
  }),

  // 工作负载操作事件数据
  workloadOperation: (operation, workload, result = null) => ({
    operation, // 'create', 'update', 'delete', 'scale', 'restart'
    workload,
    result,
    timestamp: Date.now()
  }),

  // Pod 操作事件数据
  podOperation: (operation, pod, result = null) => ({
    operation, // 'create', 'update', 'delete'
    pod,
    result,
    timestamp: Date.now()
  }),

  // 对话框操作事件数据
  dialogOperation: (action, dialogType, data = null) => ({
    action, // 'open', 'close'
    dialogType, // 'workload-detail', 'pod-list', 'logs', 'yaml', 'scale' etc.
    data,
    timestamp: Date.now()
  }),

  // 错误事件数据
  error: (error, context = '') => ({
    error: error.message || error,
    context,
    timestamp: Date.now(),
    stack: error.stack
  }),

  // 成功事件数据
  success: (message, context = '') => ({
    message,
    context,
    timestamp: Date.now()
  })
}

// 事件总线 Composable
export function useEventBus() {
  return {
    // 订阅事件
    on: k8sEventBus.on.bind(k8sEventBus),

    // 取消订阅
    off: k8sEventBus.off.bind(k8sEventBus),

    // 触发事件
    emit: k8sEventBus.emit.bind(k8sEventBus),

    // 一次性订阅
    once: k8sEventBus.once.bind(k8sEventBus),

    // 清空事件
    clear: k8sEventBus.clear.bind(k8sEventBus),

    // 事件常量
    events: K8S_EVENTS,

    // 事件数据创建工具
    createEventData
  }
}

// 全局状态管理（简单版 store）
export const k8sGlobalState = reactive({
  // 当前选中的集群
  selectedCluster: {
    id: '',
    name: '',
    status: 0
  },

  // 当前选中的命名空间
  selectedNamespace: 'default',

  // 工作负载列表
  workloadList: [],

  // 加载状态
  loading: {
    cluster: false,
    namespace: false,
    workload: false,
    pods: false
  },

  // 对话框状态
  dialogs: {
    workloadDetail: { visible: false, data: null },
    podList: { visible: false, data: null },
    logs: { visible: false, data: null },
    yaml: { visible: false, data: null },
    scale: { visible: false, data: null }
  },

  // 错误信息
  errors: [],

  // 最后更新时间
  lastUpdate: null
})

// 全局状态管理 Composable
export function useK8sGlobalState() {
  const updateCluster = (cluster) => {
    k8sGlobalState.selectedCluster = cluster
    k8sEventBus.emit(K8S_EVENTS.CLUSTER_CHANGED, createEventData.clusterChanged(cluster.id, cluster))
  }

  const updateNamespace = (namespace) => {
    k8sGlobalState.selectedNamespace = namespace
    k8sEventBus.emit(K8S_EVENTS.NAMESPACE_CHANGED, createEventData.namespaceChanged(namespace, k8sGlobalState.selectedCluster.id))
  }

  const updateWorkloadList = (workloads) => {
    k8sGlobalState.workloadList = workloads
    k8sGlobalState.lastUpdate = Date.now()
  }

  const setLoading = (key, value) => {
    k8sGlobalState.loading[key] = value
    if (value) {
      k8sEventBus.emit(K8S_EVENTS.LOADING_START, { key })
    } else {
      k8sEventBus.emit(K8S_EVENTS.LOADING_END, { key })
    }
  }

  const openDialog = (type, data = null) => {
    if (k8sGlobalState.dialogs[type]) {
      k8sGlobalState.dialogs[type] = { visible: true, data }
      k8sEventBus.emit(K8S_EVENTS.DIALOG_OPEN, createEventData.dialogOperation('open', type, data))
    }
  }

  const closeDialog = (type) => {
    if (k8sGlobalState.dialogs[type]) {
      k8sGlobalState.dialogs[type] = { visible: false, data: null }
      k8sEventBus.emit(K8S_EVENTS.DIALOG_CLOSE, createEventData.dialogOperation('close', type))
    }
  }

  const addError = (error, context = '') => {
    const errorData = createEventData.error(error, context)
    k8sGlobalState.errors.push(errorData)
    k8sEventBus.emit(K8S_EVENTS.ERROR_OCCURRED, errorData)
  }

  const clearErrors = () => {
    k8sGlobalState.errors = []
  }

  const refreshData = () => {
    k8sEventBus.emit(K8S_EVENTS.REFRESH_DATA)
  }

  return {
    // 状态
    state: k8sGlobalState,

    // 方法
    updateCluster,
    updateNamespace,
    updateWorkloadList,
    setLoading,
    openDialog,
    closeDialog,
    addError,
    clearErrors,
    refreshData
  }
}