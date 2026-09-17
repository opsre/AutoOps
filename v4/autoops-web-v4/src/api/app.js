import request from '@/utils/request'

// �(�API
const appApi = {
  // ��(
  createApplication(data) {
    return request({
      url: '/apps',
      method: 'post',
      data
    })
  },

  // �֔(h
  getApplicationList(params) {
    return request({
      url: '/apps',
      method: 'get',
      params
    })
  },

  // �֔(��
  getApplicationDetail(id) {
    return request({
      url: `/apps/${id}`,
      method: 'get'
    })
  },

  // ���(
  updateApplication(id, data) {
    return request({
      url: `/apps/${id}`,
      method: 'put',
      data
    })
  },

  //  d�(
  deleteApplication(id) {
    return request({
      url: `/apps/${id}`,
      method: 'delete'
    })
  },

  // �֔(�Jenkins��Mn
  getAppJenkinsEnvs(id) {
    return request({
      url: `/apps/${id}/jenkins-envs`,
      method: 'get'
    })
  },

  // ���(�Jenkins��Mn
  addAppJenkinsEnv(id, data) {
    return request({
      url: `/apps/${id}/jenkins-envs`,
      method: 'post',
      data
    })
  },

  // ���(�Jenkins��Mn
  updateAppJenkinsEnv(id, envId, data) {
    return request({
      url: `/apps/${id}/jenkins-envs/${envId}`,
      method: 'put',
      data
    })
  },

  //  d�(�Jenkins��Mn
  deleteAppJenkinsEnv(id, envId) {
    return request({
      url: `/apps/${id}/jenkins-envs/${envId}`,
      method: 'delete'
    })
  },

  // 获取Jenkins服务器列表
  getJenkinsServers() {
    return request({
      url: '/apps/jenkins-servers',
      method: 'get'
    })
  },

  // 搜索Jenkins任务
  searchJenkinsJobs(serverId, searchKey = '') {
    return request({
      url: `/jenkins/${serverId}/jobs/search`,
      method: 'get',
      params: {
        search: searchKey
      }
    })
  },

  // 验证Jenkins任务是否存在
  validateJenkinsJob(data) {
    return request({
      url: '/apps/jenkins-job/validate',
      method: 'post',
      data
    })
  },

  // 获取Jenkins任务参数
  getJenkinsJobParameters(serverId, jobName) {
    return request({
      url: `/jenkins/${serverId}/jobs/${jobName}/parameters`,
      method: 'get'
    })
  },

  // 启动Jenkins任务(带参数)
  startJenkinsJob(serverId, jobName, data) {
    return request({
      url: `/jenkins/${serverId}/jobs/${jobName}/start`,
      method: 'post',
      data
    })
  },

  // 快速发布相关API

  // 获取可发布应用列表
  getDeployableApplications(params) {
    return request({
      url: '/apps/deployment/applications',
      method: 'get',
      params
    })
  },

  // 创建快速发布流程
  createQuickDeployment(data) {
    return request({
      url: '/apps/deployment/quick',
      method: 'post',
      data
    })
  },

  // 执行快速发布流程
  executeQuickDeployment(data) {
    return request({
      url: '/apps/deployment/execute',
      method: 'post',
      data
    })
  },

  // 获取发布详情
  getDeploymentDetail(id) {
    return request({
      url: `/apps/deployment/${id}`,
      method: 'get'
    })
  },

  // 获取发布列表
  getDeploymentList(params) {
    return request({
      url: '/apps/deployment/list',
      method: 'get',
      params
    })
  },

  // 业务组筛选服务树
  getServiceTree(params) {
    return request({
      url: '/apps/service-tree',
      method: 'get',
      params
    })
  },

  // 获取业务组选项数据（级联选择器用）
  getBusinessGroupOptions() {
    return request({
      url: '/apps/business-group-options',
      method: 'get'
    })
  },

  // 获取应用环境配置
  getEnvironment(params) {
    return request({
      url: '/apps/environment',
      method: 'get',
      params
    })
  },

  // === 新增的4个核心API ===

  // 1. 启动Jenkins任务(总任务，多任务，依次串行执行)
  // 执行发布任务已存在: executeQuickDeployment

  // 2. 实时获取每个任务的日志
  getTaskLog(taskId, start = 0) {
    return request({
      url: `/apps/deployment/tasks/${taskId}/log`,
      method: 'get',
      params: { start }
    })
  },

  // 3. 实时获取每个任务的构建状态
  getTaskStatus(taskId) {
    return request({
      url: `/apps/deployment/tasks/${taskId}/status`,
      method: 'get'
    })
  },

  // 4. 停止任务
  stopTask(taskId) {
    return request({
      url: `/apps/deployment/tasks/${taskId}/stop`,
      method: 'post'
    })
  },

  // 通过Jenkins API停止任务（备用方法）
  stopJenkinsTask(serverId, jobName, buildNumber) {
    return request({
      url: `/jenkins/${serverId}/jobs/${jobName}/builds/${buildNumber}/stop`,
      method: 'post'
    })
  },

  // 批量获取任务状态（优化性能）
  getBatchTaskStatus(taskIds) {
    return request({
      url: '/apps/deployment/tasks/batch-status',
      method: 'post',
      data: { task_ids: taskIds }
    })
  },

  // 删除发布记录
  deleteDeployment(id) {
    return request({
      url: `/apps/deployment/${id}`,
      method: 'delete'
    })
  },

  // ========== 服务上线工单 API ==========

  // 创建服务上线工单
  createServiceRelease(data) {
    return request({
      url: '/apps/service-release',
      method: 'post',
      data
    })
  },

  // 获取服务上线工单列表
  getServiceReleaseList(params) {
    return request({
      url: '/apps/service-release/list',
      method: 'get',
      params
    })
  },

  // 获取服务上线工单详情
  getServiceReleaseDetail(id) {
    return request({
      url: `/apps/service-release/${id}`,
      method: 'get'
    })
  },

  // 审批服务上线工单
  approveServiceRelease(data) {
    return request({
      url: '/apps/service-release/approve',
      method: 'post',
      data
    })
  },

  // 运维发布
  deployServiceRelease(data) {
    return request({
      url: '/apps/service-release/deploy',
      method: 'post',
      data
    })
  },

  // 回归测试
  regressionTestServiceRelease(data) {
    return request({
      url: '/apps/service-release/regression-test',
      method: 'post',
      data
    })
  },

  // 获取服务发布的服务日志
  getServiceReleaseLog(serviceId, start = 0) {
    return request({
      url: `/apps/service-release/services/${serviceId}/log`,
      method: 'get',
      params: { start }
    })
  },

  // 检查用户权限
  checkServiceReleasePermission(releaseId, operationType) {
    return request({
      url: `/apps/service-release/${releaseId}/permission`,
      method: 'get',
      params: { operation_type: operationType }
    })
  },

  // 删除服务上线工单
  deleteServiceRelease(id) {
    return request({
      url: `/apps/service-release/${id}`,
      method: 'delete'
    })
  },

  // ========== 脚本上线工单 API ==========

  // 创建脚本上线工单
  createShRelease(data) {
    return request({
      url: '/apps/sh-release',
      method: 'post',
      data
    })
  },

  // 获取脚本上线工单列表
  getShReleaseList(params) {
    return request({
      url: '/apps/sh-release/list',
      method: 'get',
      params
    })
  },

  // 获取脚本上线工单详情
  getShReleaseDetail(id) {
    return request({
      url: `/apps/sh-release/${id}`,
      method: 'get'
    })
  },

  // 审核脚本上线工单
  approveShRelease(data) {
    return request({
      url: '/apps/sh-release/approve',
      method: 'post',
      data
    })
  },

  // 执行Jenkins任务
  executeShRelease(data) {
    return request({
      url: '/apps/sh-release/pull-code',
      method: 'post',
      data
    })
  },

  // 执行脚本（通过SSH）
  executeShReleaseScript(data) {
    return request({
      url: '/apps/sh-release/execute-script',
      method: 'post',
      data
    })
  },

  // 获取Jenkins执行日志
  getShReleaseLog(id, start = 0) {
    return request({
      url: `/apps/sh-release/${id}/jenkins-log`,
      method: 'get',
      params: { start }
    })
  },

  // 获取脚本执行输出
  getShReleaseScriptOutput(id) {
    return request({
      url: `/apps/sh-release/${id}/script-output`,
      method: 'get'
    })
  },

  // 检查脚本上线用户权限
  checkShReleasePermission(releaseId, operationType) {
    return request({
      url: `/apps/sh-release/${releaseId}/permission`,
      method: 'get',
      params: { operation_type: operationType }
    })
  },

  // 列出脚本发布工单的服务器文件
  listShReleaseFiles(id) {
    return request({
      url: `/apps/sh-release/${id}/files`,
      method: 'get'
    })
  },

  // 下载脚本发布工单的服务器文件
  downloadShReleaseFile(id, filePath) {
    return request({
      url: `/apps/sh-release/${id}/download`,
      method: 'get',
      params: {
        file_path: filePath
      },
      responseType: 'blob'
    })
  },

  // 删除脚本上线工单
  deleteShRelease(id) {
    return request({
      url: `/apps/sh-release/${id}`,
      method: 'delete'
    })
  },

  // ========== 应用视图树 API ==========

  // 获取应用视图树（后端已构建好的树结构）
  getApplicationViewTree(params) {
    return request({
      url: '/apps/view-tree',
      method: 'get',
      params
    })
  },

  // 获取应用完整详情（用于右侧详情面板）
  getApplicationDetail(id) {
    return request({
      url: `/apps/${id}/detail`,
      method: 'get'
    })
  },

  // 获取应用发布历史
  getApplicationReleaseHistory(params) {
    return request({
      url: '/apps/release-history',
      method: 'get',
      params
    })
  }
}

export default appApi