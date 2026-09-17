import request from "@/utils/request"

export default {
    // K8s集群管理
    
    // 创建K8s集群
    createCluster(data) {
        return request({
            url: '/api/v1/k8s/cluster',
            method: 'post',
            data: data
        })
    },

    // 获取K8s集群列表
    getClusterList(params) {
        return request({
            url: '/api/v1/k8s/cluster',
            method: 'get',
            params: params || {}
        })
    },

    // 获取K8s集群详情
    getClusterById(id) {
        return request({
            url: `/api/v1/k8s/cluster/${id}`,
            method: 'get'
        })
    },

    // 获取K8s集群完整详情信息（仪表板）
    getClusterDetail(id) {
        return request({
            url: `/api/v1/k8s/cluster/${id}/detail`,
            method: 'get'
        })
    },

    // 更新K8s集群
    updateCluster(id, data) {
        return request({
            url: `/api/v1/k8s/cluster/${id}`,
            method: 'put',
            data: data
        })
    },

    // 删除K8s集群
    deleteCluster(id) {
        return request({
            url: `/api/v1/k8s/cluster/${id}`,
            method: 'delete'
        })
    },

    // 获取K8s集群状态
    getClusterStatus(id) {
        return request({
            url: `/api/v1/k8s/cluster/${id}/status`,
            method: 'get'
        })
    },

    // K8s节点管理
    
    // 获取集群节点列表
    getClusterNodes(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes`,
            method: 'get',
            params: params
        })
    },

    // 根据节点名称搜索
    searchNodesByName(clusterId, name) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/search`,
            method: 'get',
            params: { name }
        })
    },

    // 根据节点角色查询
    searchNodesByRole(clusterId, role) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/byRole`,
            method: 'get',
            params: { role }
        })
    },

    // 根据节点状态查询
    searchNodesByStatus(clusterId, status) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/byStatus`,
            method: 'get',
            params: { status }
        })
    },

    // 获取节点详细信息
    getNodeDetail(clusterId, nodeName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/${nodeName}`,
            method: 'get'
        })
    },

    // 获取节点资源分配详情
    getNodeResources(clusterId, nodeName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/${nodeName}/resources`,
            method: 'get'
        })
    },

    // 添加节点污点
    addNodeTaint(clusterId, nodeName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/${nodeName}/taints`,
            method: 'post',
            data: data
        })
    },

    // 移除节点污点
    removeNodeTaint(clusterId, nodeName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/${nodeName}/taints`,
            method: 'delete',
            data: data
        })
    },

    // 添加节点标签
    addNodeLabel(clusterId, nodeName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/${nodeName}/labels`,
            method: 'post',
            data: data
        })
    },

    // 移除节点标签
    removeNodeLabel(clusterId, nodeName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/${nodeName}/labels`,
            method: 'delete',
            data: data
        })
    },

    // 封锁/解封节点
    cordonNode(clusterId, nodeName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/${nodeName}/cordon`,
            method: 'post',
            data: data
        })
    },

    // 获取节点YAML
    getNodeYaml(clusterId, nodeName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/${nodeName}/yaml`,
            method: 'get'
        })
    },

    // 更新节点YAML
    updateNodeYaml(clusterId, nodeName, yamlContent) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/${nodeName}/yaml`,
            method: 'put',
            data: { yaml: yamlContent }
        })
    },

    // 驱逐节点
    drainNode(clusterId, nodeName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/${nodeName}/drain`,
            method: 'post',
            data: data
        })
    },

    // 添加节点到集群
    addNodeToCluster(clusterId, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes`,
            method: 'post',
            data: data
        })
    },

    // 从集群移除节点
    removeNodeFromCluster(clusterId, nodeId) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/nodes/${nodeId}`,
            method: 'delete'
        })
    },

    // K8s任务管理
    
    // 创建K8s部署任务
    createK8sTask(data) {
        return request({
            url: '/api/v1/task/k8s',
            method: 'post',
            data: data
        })
    },

    // 启动Ansible任务
    startAnsibleTask(taskId) {
        return request({
            url: `/api/v1/task/ansible/${taskId}/start`,
            method: 'post'
        })
    },

    // 获取任务状态
    getTaskStatus(taskId) {
        return request({
            url: `/api/v1/task/${taskId}/status`,
            method: 'get'
        })
    },

    // 获取任务日志
    getTaskLogs(taskId, params) {
        return request({
            url: `/api/v1/task/${taskId}/logs`,
            method: 'get',
            params: params
        })
    },

    // K8s应用管理
    // 部署应用到K8s集群
    deployApplication(clusterId, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/applications`,
            method: 'post',
            data: data
        })
    },

    // 获取集群中的应用列表
    getApplicationList(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/applications`,
            method: 'get',
            params: params
        })
    },

    // 删除应用
    deleteApplication(clusterId, appId) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/applications/${appId}`,
            method: 'delete'
        })
    },

    // K8s资源管理
    
    // 获取集群资源使用情况
    getClusterResources(clusterId) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/resources`,
            method: 'get'
        })
    },

    // 获取集群事件
    getClusterEvents(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/events`,
            method: 'get',
            params: params
        })
    },

    // 获取命名空间事件
    getNamespaceEvents(clusterId, namespaceName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/events`,
            method: 'get',
            params: params
        })
    },

    // 执行kubectl命令
    executeKubectl(clusterId, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/kubectl`,
            method: 'post',
            data: data
        })
    },

    // K8s配置管理
    
    // 更新集群配置
    updateClusterConfig(clusterId, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/config`,
            method: 'put',
            data: data
        })
    },

    // 获取集群kubeconfig
    getKubeConfig(clusterId) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/kubeconfig`,
            method: 'get'
        })
    },

    // 注册外部K8s集群
    registerExternalCluster(data) {
        return request({
            url: '/api/v1/k8s/cluster',
            method: 'post',
            data: data
        })
    },

    // 同步集群信息
    syncCluster(clusterId) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/sync`,
            method: 'post'
        })
    },

    // K8s命名空间管理

    // 命名空间名称模糊查询
    searchNamespaces(clusterId, keyword) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/search`,
            method: 'get',
            params: { keyword }
        })
    },

    // 获取命名空间YAML
    getNamespaceYaml(clusterId, namespaceName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/yaml`,
            method: 'get'
        })
    },

    // 更新命名空间YAML
    updateNamespaceYaml(clusterId, namespaceName, yamlContent) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/yaml`,
            method: 'put',
            data: { yaml: yamlContent }
        })
    },

    // 更新命名空间标签（全量替换）
    updateNamespaceLabels(clusterId, namespaceName, labels) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/labels`,
            method: 'patch',
            data: { labels }
        })
    },

    // 获取集群命名空间列表
    getNamespaces(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces`,
            method: 'get',
            params: params
        })
    },

    // 创建命名空间
    createNamespace(clusterId, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces`,
            method: 'post',
            data: data
        })
    },

    // 获取命名空间详情
    getNamespaceDetail(clusterId, namespaceName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}`,
            method: 'get'
        })
    },

    // 删除命名空间
    deleteNamespace(clusterId, namespaceName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}`,
            method: 'delete'
        })
    },

    // ResourceQuota管理
    
    // 获取ResourceQuota列表
    getResourceQuotas(clusterId, namespaceName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/resourcequotas`,
            method: 'get'
        })
    },

    // 创建ResourceQuota
    createResourceQuota(clusterId, namespaceName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/resourcequotas`,
            method: 'post',
            data: data
        })
    },

    // 更新ResourceQuota
    updateResourceQuota(clusterId, namespaceName, quotaName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/resourcequotas/${quotaName}`,
            method: 'put',
            data: data
        })
    },

    // 删除ResourceQuota
    deleteResourceQuota(clusterId, namespaceName, quotaName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/resourcequotas/${quotaName}`,
            method: 'delete'
        })
    },

    // LimitRange管理
    
    // 获取LimitRange列表
    getLimitRanges(clusterId, namespaceName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/limitranges`,
            method: 'get'
        })
    },

    // 创建LimitRange
    createLimitRange(clusterId, namespaceName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/limitranges`,
            method: 'post',
            data: data
        })
    },

    // 更新LimitRange
    updateLimitRange(clusterId, namespaceName, limitRangeName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/limitranges/${limitRangeName}`,
            method: 'put',
            data: data
        })
    },

    // 删除LimitRange
    deleteLimitRange(clusterId, namespaceName, limitRangeName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/limitranges/${limitRangeName}`,
            method: 'delete'
        })
    },

    // K8s工作负载管理
    
    // 获取工作负载列表
    getWorkloadList(clusterId, namespaceName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/workloads`,
            method: 'get',
            params: params
        })
    },

    // 获取无控制器 Pod 列表
    getStandalonePodList(clusterId, namespaceName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/standalone`,
            method: 'get',
            params: params
        })
    },

    // 获取集群级工作负载列表（全部命名空间）
    getClusterWorkloadList(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/workloads`,
            method: 'get',
            params: params
        })
    },

    // 按名称模糊搜索集群级工作负载
    searchClusterWorkloadList(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/workloads/search`,
            method: 'get',
            params: params
        })
    },

    // 获取集群级异常工作负载列表
    getClusterAbnormalWorkloadList(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/workloads/abnormal`,
            method: 'get',
            params: params
        })
    },

    // 按名称模糊搜索工作负载
    searchWorkloadList(clusterId, namespaceName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/workloads/search`,
            method: 'get',
            params: params
        })
    },

    // 获取异常工作负载列表
    getAbnormalWorkloadList(clusterId, namespaceName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/workloads/abnormal`,
            method: 'get',
            params: params
        })
    },

    // 获取工作负载详情
    getWorkloadDetail(clusterId, namespaceName, workloadType, workloadName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/workloads/${workloadType}/${workloadName}`,
            method: 'get'
        })
    },

    // 删除工作负载
    deleteWorkload(clusterId, namespaceName, workloadType, workloadName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/workloads/${workloadType}/${workloadName}`,
            method: 'delete'
        })
    },

    // Deployment管理
    
    // 创建Deployment
    createDeployment(clusterId, namespaceName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/deployments`,
            method: 'post',
            data: data
        })
    },

    // 更新Deployment
    updateDeployment(clusterId, namespaceName, deploymentName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/deployments/${deploymentName}`,
            method: 'put',
            data: data
        })
    },

    // 删除Deployment
    deleteDeployment(clusterId, namespaceName, deploymentName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/deployments/${deploymentName}`,
            method: 'delete'
        })
    },

    // 扩缩容Deployment
    scaleDeployment(clusterId, namespaceName, deploymentName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/deployments/${deploymentName}/scale`,
            method: 'post',
            data: data
        })
    },

    // 重启Deployment
    restartDeployment(clusterId, namespaceName, deploymentName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/deployments/${deploymentName}/restart`,
            method: 'post'
        })
    },

    // 获取Deployment历史版本
    getDeploymentHistory(clusterId, namespaceName, deploymentName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/deployments/${deploymentName}/history`,
            method: 'get'
        })
    },

    // 回滚Deployment到指定版本
    rollbackDeployment(clusterId, namespaceName, deploymentName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/deployments/${deploymentName}/rollback`,
            method: 'post',
            data: data
        })
    },

    // K8s Pod管理
    
    // 获取Pod列表
    getPodList(clusterId, namespaceName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods`,
            method: 'get',
            params: params
        })
    },


    // 删除Pod
    deletePod(clusterId, namespaceName, podName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/${podName}`,
            method: 'delete'
        })
    },

    // 获取Pod日志
    getPodLogs(clusterId, namespaceName, podName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/${podName}/logs`,
            method: 'get',
            params: params
        })
    },

    // 上传文件到Pod
    uploadPodFile(clusterId, namespaceName, podName, formData, config = {}) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/${podName}/files/upload`,
            method: 'post',
            data: formData,
            headers: {
                'Content-Type': 'multipart/form-data'
            },
            ...config
        })
    },

    // 列出Pod文件列表
    listPodFiles(clusterId, namespaceName, podName, params, config = {}) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/${podName}/files/list`,
            method: 'get',
            params: params,
            ...config
        })
    },

    // 从Pod下载文件
    downloadPodFile(clusterId, namespaceName, podName, params, config = {}) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/${podName}/files/download`,
            method: 'get',
            params: params,
            responseType: 'blob',
            timeout: 600000,
            ...config
        })
    },

    // 获取Pod YAML
    getPodYaml(clusterId, namespaceName, podName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/${podName}/yaml`,
            method: 'get'
        })
    },

    // 获取Pod事件
    getPodEvents(clusterId, namespaceName, podName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/${podName}/events`,
            method: 'get'
        })
    },

    // 获取工作负载YAML
    getWorkloadYaml(clusterId, namespaceName, workloadType, workloadName) {
        return request({
            url: `/k8s/cluster/${clusterId}/namespaces/${namespaceName}/workload-yaml/${workloadType}/${workloadName}`,
            method: 'get'
        })
    },

    // 获取工作负载下的Pod列表
    getWorkloadPods(clusterId, namespaceName, workloadType, workloadName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/workloads/${workloadType}/${workloadName}/pods`,
            method: 'get'
        })
    },

    // 更新Pod YAML
    updatePodYaml(clusterId, namespaceName, podName, yamlContent) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/${podName}/yaml`,
            method: 'put',
            data: { yaml: yamlContent }
        })
    },

    // 更新工作负载YAML
    updateWorkloadYaml(clusterId, namespaceName, workloadType, workloadName, yamlContent, dryRun = false) {
        return request({
            url: `/k8s/cluster/${clusterId}/namespaces/${namespaceName}/workload-yaml`,
            method: 'put',
            data: {
                workloadType: workloadType,
                workloadName: workloadName,
                yamlContent: yamlContent,
                dryRun: dryRun
            }
        })
    },

    // 通过kubectl apply更新资源
    applyYaml(clusterId, yamlContent) {
        return request({
            url: `/k8s/cluster/${clusterId}/apply`,
            method: 'post',
            data: { yamlContent: yamlContent }
        })
    },

    // 获取Pod容器列表
    getPodContainers(clusterId, namespaceName, podName) {
        return request({
            url: `/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/${podName}/containers`,
            method: 'get'
        })
    },

    // 创建Pod临时调试容器
    createPodDebugContainer(clusterId, namespaceName, podName, data = {}) {
        return request({
            url: `/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/${podName}/debug`,
            method: 'post',
            data
        })
    },

    // K8s监控管理
    
    // 获取Pod监控指标
    getPodMetrics(clusterId, namespaceName, podName) {
        return request({
            url: `/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/${podName}/metrics`,
            method: 'get'
        })
    },

    // 获取节点监控指标
    getNodeMetrics(clusterId, nodeName) {
        return request({
            url: `/k8s/cluster/${clusterId}/nodes/${nodeName}/metrics`,
            method: 'get'
        })
    },

    // 获取命名空间监控指标
    getNamespaceMetrics(clusterId, namespaceName) {
        return request({
            url: `/k8s/cluster/${clusterId}/namespaces/${namespaceName}/metrics`,
            method: 'get'
        })
    },

    // K8s Pod创建管理

    // 通过YAML创建Pod/工作负载
    createPodFromYaml(clusterId, namespaceName, data) {
        return request({
            url: `/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pods/yaml`,
            method: 'post',
            data: { yamlContent: data.yamlContent }
        })
    },

    // 校验YAML格式
    validateYaml(clusterId, yamlContent) {
        return request({
            url: `/k8s/cluster/${clusterId}/yaml/validate`,
            method: 'post',
            data: { yamlContent: yamlContent }
        })
    },

    // K8s Service 管理

    // 获取Service列表
    getServiceList(clusterId, namespaceName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/services`,
            method: 'get',
            params: params
        })
    },

    // 获取Service详情
    getServiceDetail(clusterId, namespaceName, serviceName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/services/${serviceName}`,
            method: 'get'
        })
    },

    // 创建Service
    createService(clusterId, namespaceName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/services`,
            method: 'post',
            data: data
        })
    },

    // 更新Service
    updateService(clusterId, namespaceName, serviceName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/services/${serviceName}`,
            method: 'put',
            data: data
        })
    },

    // 删除Service
    deleteService(clusterId, namespaceName, serviceName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/services/${serviceName}`,
            method: 'delete'
        })
    },

    // 获取Service YAML
    getServiceYaml(clusterId, namespaceName, serviceName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/services/${serviceName}/yaml`,
            method: 'get'
        })
    },

    // 更新Service YAML
    updateServiceYaml(clusterId, namespaceName, serviceName, yamlContent) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/services/${serviceName}/yaml`,
            method: 'put',
            data: { yaml: yamlContent }
        })
    },

    // 获取Service事件
    getServiceEvents(clusterId, namespaceName, serviceName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/services/${serviceName}/events`,
            method: 'get'
        })
    },

    // 获取Service监控指标
    getServiceMetrics(clusterId, namespaceName, serviceName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/services/${serviceName}/metrics`,
            method: 'get'
        })
    },

    // K8s Ingress 管理

    // 获取Ingress列表
    getIngressList(clusterId, namespaceName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/ingresses`,
            method: 'get',
            params: params
        })
    },

    // 获取Ingress详情
    getIngressDetail(clusterId, namespaceName, ingressName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/ingresses/${ingressName}`,
            method: 'get'
        })
    },

    // 创建Ingress
    createIngress(clusterId, namespaceName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/ingresses`,
            method: 'post',
            data: data
        })
    },

    // 更新Ingress
    updateIngress(clusterId, namespaceName, ingressName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/ingresses/${ingressName}`,
            method: 'put',
            data: data
        })
    },

    // 删除Ingress
    deleteIngress(clusterId, namespaceName, ingressName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/ingresses/${ingressName}`,
            method: 'delete'
        })
    },

    // 获取Ingress YAML
    getIngressYaml(clusterId, namespaceName, ingressName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/ingresses/${ingressName}/yaml`,
            method: 'get'
        })
    },

    // 更新Ingress YAML
    updateIngressYaml(clusterId, namespaceName, ingressName, yamlContent) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/ingresses/${ingressName}/yaml`,
            method: 'put',
            data: { yaml: yamlContent }
        })
    },

    // 获取Ingress事件
    getIngressEvents(clusterId, namespaceName, ingressName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/ingresses/${ingressName}/events`,
            method: 'get'
        })
    },

    // 获取Ingress监控指标
    getIngressMetrics(clusterId, namespaceName, ingressName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/ingresses/${ingressName}/metrics`,
            method: 'get'
        })
    },

    // 测试Ingress后端服务
    testIngressBackend(clusterId, namespaceName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/ingresses/test-backend`,
            method: 'post',
            data: data
        })
    },

    // K8s 存储管理

    // PVC 管理
    // 获取PVC列表
    getPVCList(clusterId, namespaceName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pvcs`,
            method: 'get',
            params: params || {}
        })
    },

    // 获取PVC详情
    getPVCDetail(clusterId, namespaceName, pvcName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pvcs/${pvcName}`,
            method: 'get'
        })
    },

    // 创建PVC
    createPVC(clusterId, namespaceName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pvcs`,
            method: 'post',
            data: data
        })
    },

    // 更新PVC
    updatePVC(clusterId, namespaceName, pvcName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pvcs/${pvcName}`,
            method: 'put',
            data: data
        })
    },

    // 删除PVC
    deletePVC(clusterId, namespaceName, pvcName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pvcs/${pvcName}`,
            method: 'delete'
        })
    },

    // 获取PVC YAML
    getPVCYaml(clusterId, namespaceName, pvcName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pvcs/${pvcName}/yaml`,
            method: 'get'
        })
    },

    // 更新PVC YAML
    updatePVCYaml(clusterId, namespaceName, pvcName, yamlContent) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/pvcs/${pvcName}/yaml`,
            method: 'put',
            data: { yaml: yamlContent }
        })
    },

    // PV 管理
    // 获取PV列表
    getPVList(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/pvs`,
            method: 'get',
            params: params || {}
        })
    },

    // 获取PV详情
    getPVDetail(clusterId, pvName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/pvs/${pvName}`,
            method: 'get'
        })
    },

    // 创建PV
    createPV(clusterId, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/pvs`,
            method: 'post',
            data: data
        })
    },

    // 删除PV
    deletePV(clusterId, pvName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/pvs/${pvName}`,
            method: 'delete'
        })
    },

    // 获取PV YAML
    getPVYaml(clusterId, pvName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/pvs/${pvName}/yaml`,
            method: 'get'
        })
    },

    // 更新PV YAML
    updatePVYaml(clusterId, pvName, yamlContent) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/pvs/${pvName}/yaml`,
            method: 'put',
            data: { yaml: yamlContent }
        })
    },

    // StorageClass 管理
    // 获取StorageClass列表
    getStorageClassList(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/storageclasses`,
            method: 'get',
            params: params || {}
        })
    },

    // 获取StorageClass详情
    getStorageClassDetail(clusterId, storageClassName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/storageclasses/${storageClassName}`,
            method: 'get'
        })
    },

    // 创建StorageClass
    createStorageClass(clusterId, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/storageclasses`,
            method: 'post',
            data: data
        })
    },

    // 删除StorageClass
    deleteStorageClass(clusterId, storageClassName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/storageclasses/${storageClassName}`,
            method: 'delete'
        })
    },

    // 获取StorageClass YAML
    getStorageClassYaml(clusterId, storageClassName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/storageclasses/${storageClassName}/yaml`,
            method: 'get'
        })
    },

    // 更新StorageClass YAML
    updateStorageClassYaml(clusterId, storageClassName, yamlContent) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/storageclasses/${storageClassName}/yaml`,
            method: 'put',
            data: { yaml: yamlContent }
        })
    },

    // K8s 配置管理

    // ConfigMap 管理
    // 获取ConfigMap列表
    getConfigMaps(clusterId, namespaceName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/configmaps`,
            method: 'get',
            params: params || {}
        })
    },

    // 获取ConfigMap详情
    getConfigMapDetail(clusterId, namespaceName, configMapName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/configmaps/${configMapName}`,
            method: 'get'
        })
    },

    // 创建ConfigMap
    createConfigMap(clusterId, namespaceName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/configmaps`,
            method: 'post',
            data: data
        })
    },

    // 更新ConfigMap
    updateConfigMap(clusterId, namespaceName, configMapName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/configmaps/${configMapName}`,
            method: 'put',
            data: data
        })
    },

    // 删除ConfigMap
    deleteConfigMap(clusterId, namespaceName, configMapName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/configmaps/${configMapName}`,
            method: 'delete'
        })
    },

    // 获取ConfigMap YAML
    getConfigMapYaml(clusterId, namespaceName, configMapName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/configmaps/${configMapName}/yaml`,
            method: 'get'
        })
    },

    // 更新ConfigMap YAML
    updateConfigMapYaml(clusterId, namespaceName, configMapName, yamlContent) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/configmaps/${configMapName}/yaml`,
            method: 'put',
            data: { yaml: yamlContent }
        })
    },

    // Secret 管理
    // 获取Secret列表
    getSecrets(clusterId, namespaceName, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/secrets`,
            method: 'get',
            params: params || {}
        })
    },

    // 获取Secret详情
    getSecretDetail(clusterId, namespaceName, secretName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/secrets/${secretName}`,
            method: 'get'
        })
    },

    // 创建Secret
    createSecret(clusterId, namespaceName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/secrets`,
            method: 'post',
            data: data
        })
    },

    // 更新Secret
    updateSecret(clusterId, namespaceName, secretName, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/secrets/${secretName}`,
            method: 'put',
            data: data
        })
    },

    // 删除Secret
    deleteSecret(clusterId, namespaceName, secretName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/secrets/${secretName}`,
            method: 'delete'
        })
    },

    // 获取Secret YAML
    getSecretYaml(clusterId, namespaceName, secretName) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/secrets/${secretName}/yaml`,
            method: 'get'
        })
    },

    // 更新Secret YAML
    updateSecretYaml(clusterId, namespaceName, secretName, yamlContent) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/namespaces/${namespaceName}/secrets/${secretName}/yaml`,
            method: 'put',
            data: { yaml: yamlContent }
        })
    },

    // K8s Arthas Java诊断管理

    // 检查Arthas安装状态
    checkArthasInstall(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/check`,
            method: 'get',
            params: params
        })
    },

    // 安装Arthas
    installArthas(clusterId, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/install`,
            method: 'post',
            data: data
        })
    },

    // 列表Java进程
    listJavaProcesses(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/java-processes`,
            method: 'get',
            params: params
        })
    },

    // 获取JVM仪表板
    getArthasDashboard(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/dashboard`,
            method: 'get',
            params: params,
            timeout: 60000 // 60秒
        })
    },

    // 获取线程信息
    getArthasThreadInfo(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/thread`,
            method: 'get',
            params: params,
            timeout: 60000 // 60秒
        })
    },

    // 获取线程堆栈
    getArthasThreadStack(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/thread/stack`,
            method: 'get',
            params: params,
            timeout: 60000 // 60秒
        })
    },

    // 获取JVM信息
    getArthasJvmInfo(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/jvm`,
            method: 'get',
            params: params,
            timeout: 60000 // 60秒
        })
    },

    // 获取内存信息
    getArthasMemoryInfo(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/memory`,
            method: 'get',
            params: params,
            timeout: 60000 // 60秒
        })
    },

    // 获取系统环境变量
    getArthasSysenv(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/sysenv`,
            method: 'get',
            params: params,
            timeout: 60000 // 60秒
        })
    },

    // 获取系统属性
    getArthasSysprop(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/sysprop`,
            method: 'get',
            params: params,
            timeout: 60000 // 60秒
        })
    },

    // 获取性能计数器
    getArthasPerfcounter(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/perfcounter`,
            method: 'get',
            params: params,
            timeout: 60000 // 60秒
        })
    },

    // 反编译类
    getArthasJad(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/jad`,
            method: 'get',
            params: params,
            timeout: 60000 // 60秒
        })
    },

    // 搜索类
    searchArthasClass(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/sc`,
            method: 'get',
            params: params,
            timeout: 60000 // 60秒
        })
    },

    // 搜索方法
    searchArthasMethod(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/sm`,
            method: 'get',
            params: params,
            timeout: 60000 // 60秒
        })
    },

    // 启动火焰图采样（立即返回 taskId）
    startArthasProfiler(clusterId, params) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/profiler/start`,
            method: 'post',
            params: params,
            timeout: 30000
        })
    },

    // 查询火焰图任务状态
    getArthasProfilerTask(clusterId, taskId) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/profiler/task/${taskId}`,
            method: 'get',
            timeout: 15000
        })
    },

    // 执行自定义命令
    executeArthasCommand(clusterId, data) {
        return request({
            url: `/api/v1/k8s/cluster/${clusterId}/arthas/command`,
            method: 'post',
            data: data,
            timeout: 60000 // 60秒
        })
    },

    // 下载 arthas-boot.jar
    downloadArthas() {
        return request({
            url: `/api/v1/k8s/cluster/arthas/download`,
            method: 'post'
        })
    },

    // 上传 arthas-boot.jar
    uploadArthas(formData) {
        return request({
            url: `/api/v1/k8s/cluster/arthas/upload`,
            method: 'post',
            data: formData,
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        })
    }
}
