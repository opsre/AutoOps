/**
 * 监控管理 API（SSL证书监控 + 故障管理）
 *
 * @author xiaoRui
 */

import request from "@/utils/request"

const BASE_URL = '/api/v1/monitor'
const CMDB_BASE_URL = '/api/v1/cmdb'

// SSL证书域名监控 API
export const domainApi = {
  // 创建域名监控
  create(data) {
    return request({
      url: `${BASE_URL}/domain`,
      method: 'post',
      data
    })
  },

  // 获取域名监控列表
  getList(params) {
    return request({
      url: `${BASE_URL}/domain/list`,
      method: 'get',
      params
    })
  },

  // 获取域名监控详情
  getDetail(id) {
    return request({
      url: `${BASE_URL}/domain/${id}`,
      method: 'get'
    })
  },

  // 更新域名监控
  update(data) {
    return request({
      url: `${BASE_URL}/domain`,
      method: 'put',
      data
    })
  },

  // 删除域名监控
  delete(id) {
    return request({
      url: `${BASE_URL}/domain/${id}`,
      method: 'delete'
    })
  },

  // 检查单个域名
  check(id) {
    return request({
      url: `${BASE_URL}/domain/check/${id}`,
      method: 'post'
    })
  },

  // 检查所有域名
  checkAll() {
    return request({
      url: `${BASE_URL}/domain/check-all`,
      method: 'post'
    })
  },

  // 获取统计信息
  getStatistics() {
    return request({
      url: `${BASE_URL}/domain/statistics`,
      method: 'get'
    })
  },

  // 获取即将过期的SSL证书域名
  getSslExpiring(days) {
    return request({
      url: `${BASE_URL}/domain/ssl-expiring`,
      method: 'get',
      params: days ? { days } : {}
    })
  },

  // 获取定时巡检配置
  getScheduleConfig() {
    return request({
      url: `${BASE_URL}/domain/schedule/config`,
      method: 'get'
    })
  },

  // 更新定时巡检配置
  updateScheduleConfig(data) {
    return request({
      url: `${BASE_URL}/domain/schedule/config`,
      method: 'post',
      data
    })
  },

  // 暂停定时巡检
  pauseSchedule() {
    return request({
      url: `${BASE_URL}/domain/schedule/pause`,
      method: 'post'
    })
  },

  // 恢复定时巡检
  resumeSchedule() {
    return request({
      url: `${BASE_URL}/domain/schedule/resume`,
      method: 'post'
    })
  },

  // 获取当前巡检状态
  getScheduleStatus() {
    return request({
      url: `${BASE_URL}/domain/schedule/status`,
      method: 'get'
    })
  },

  // 获取Cron表达式下次执行时间
  getNextExecutionTime(cron) {
    return request({
      url: '/api/v1/task/next-execution',
      method: 'get',
      params: { cron }
    })
  },

  // 批量上传域名文件
  batchUpload(file) {
    const formData = new FormData()
    formData.append('file', file)
    return request({
      url: `${BASE_URL}/domain/batch-upload`,
      method: 'post',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 从阿里云DNS导入域名
  importAliyunDns(data) {
    return request({
      url: `${BASE_URL}/domain/import-aliyun-dns`,
      method: 'post',
      data
    })
  },

  // 获取通知机器人列表
  getNotifyRobots() {
    return request({
      url: `${BASE_URL}/robot/list`,
      method: 'get',
      params: { pageSize: 100 }
    })
  }
}

// 故障管理 API
export const incidentApi = {
  // 创建故障记录
  create(data) {
    return request({
      url: `${BASE_URL}/incident`,
      method: 'post',
      data
    })
  },

  // 获取故障记录列表
  getList(params) {
    return request({
      url: `${BASE_URL}/incident/list`,
      method: 'get',
      params
    })
  },

  // 获取故障记录详情
  getDetail(id) {
    return request({
      url: `${BASE_URL}/incident/${id}`,
      method: 'get'
    })
  },

  // 更新故障记录
  update(data) {
    return request({
      url: `${BASE_URL}/incident`,
      method: 'put',
      data
    })
  },

  // 删除故障记录
  delete(id) {
    return request({
      url: `${BASE_URL}/incident/${id}`,
      method: 'delete'
    })
  },

  // 批量删除故障记录
  batchDelete(ids) {
    return request({
      url: `${BASE_URL}/incident/batch-delete`,
      method: 'post',
      data: ids
    })
  },

  // 更新故障状态
  updateStatus(id, status) {
    return request({
      url: `${BASE_URL}/incident/${id}/status`,
      method: 'put',
      params: { status }
    })
  },

  // 获取故障统计信息
  getStatistics() {
    return request({
      url: `${BASE_URL}/incident/statistics`,
      method: 'get'
    })
  }
}

// SSL证书管理 API
export const sslCertApi = {
  // 申请SSL证书
  apply(data) {
    return request({
      url: `${BASE_URL}/ssl-cert/apply`,
      method: 'post',
      data
    })
  },

  // 获取SSL证书列表
  getList(params) {
    return request({
      url: `${BASE_URL}/ssl-cert/list`,
      method: 'get',
      params
    })
  },

  // 查询证书状态
  getState(id) {
    return request({
      url: `${BASE_URL}/ssl-cert/${id}/state`,
      method: 'get'
    })
  },

  // 下载证书
  download(id) {
    return request({
      url: `${BASE_URL}/ssl-cert/${id}/download`,
      method: 'get',
      responseType: 'blob' // 重要：指定响应类型为blob，用于下载文件
    })
  },

  // 删除证书
  delete(id) {
    return request({
      url: `${BASE_URL}/ssl-cert/${id}`,
      method: 'delete'
    })
  },

  // 添加DNS验证记录
  addDnsValidation(id) {
    return request({
      url: `${BASE_URL}/ssl-cert/${id}/dns-validation`,
      method: 'post'
    })
  },

  // 检查DNS验证状态
  checkValidation(id) {
    return request({
      url: `${BASE_URL}/ssl-cert/${id}/check-validation`,
      method: 'get'
    })
  },

  // 部署证书到主机
  deploy(data) {
    return request({
      url: `${BASE_URL}/ssl-cert/deploy`,
      method: 'post',
      data
    })
  },

  // 获取部署日志
  getDeployLogs(params) {
    return request({
      url: `${BASE_URL}/ssl-cert/deploy-logs`,
      method: 'get',
      params
    })
  },

  // 上传证书
  upload(formData) {
    return request({
      url: `${BASE_URL}/ssl-cert/upload`,
      method: 'post',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}

// 阿里云配置管理 API
export const aliyunConfigApi = {
  // 创建阿里云配置
  create(data) {
    return request({
      url: `${BASE_URL}/aliyun-config`,
      method: 'post',
      data
    })
  },

  // 获取阿里云配置列表
  getList(params) {
    return request({
      url: `${BASE_URL}/aliyun-config/list`,
      method: 'get',
      params
    })
  },

  // 获取阿里云配置详情
  getDetail(id) {
    return request({
      url: `${BASE_URL}/aliyun-config/${id}`,
      method: 'get'
    })
  },

  // 更新阿里云配置
  update(data) {
    return request({
      url: `${BASE_URL}/aliyun-config`,
      method: 'put',
      data
    })
  },

  // 删除阿里云配置
  delete(id) {
    return request({
      url: `${BASE_URL}/aliyun-config/${id}`,
      method: 'delete'
    })
  }
}

// 主机监控 API
export const hostMonitorApi = {
  // 获取 TOP 使用率主机
  getTopUsageHosts(params) {
    return request({
      url: `${BASE_URL}/hosts/top-usage`,
      method: 'get',
      params
    })
  },

  // 获取主机完整监控数据
  getCompleteMonitorData(hostId) {
    return request({
      url: `${BASE_URL}/hosts/${hostId}/complete`,
      method: 'get'
    })
  },

  // 获取主机Top进程
  getTopProcesses(hostId) {
    return request({
      url: `${BASE_URL}/hosts/${hostId}/top-processes`,
      method: 'get'
    })
  }
}

// 机器人管理 API
export const robotApi = {
  // 创建机器人
  create(data) {
    return request({
      url: `${BASE_URL}/robot`,
      method: 'post',
      data
    })
  },

  // 更新机器人
  update(data) {
    return request({
      url: `${BASE_URL}/robot`,
      method: 'put',
      data
    })
  },

  // 获取机器人详情
  getDetail(id) {
    return request({
      url: `${BASE_URL}/robot/${id}`,
      method: 'get'
    })
  },

  // 删除机器人
  delete(id) {
    return request({
      url: `${BASE_URL}/robot/${id}`,
      method: 'delete'
    })
  },

  // 获取机器人列表
  getList(params) {
    return request({
      url: `${BASE_URL}/robot/list`,
      method: 'get',
      params
    })
  },

  // 测试发送消息
  test(id, data) {
    return request({
      url: `${BASE_URL}/robot/test/${id}`,
      method: 'post',
      data: data || {}
    })
  }
}

// 用户统计 API
export const userStatsApi = {
  // 获取用户登录统计
  getUserStats() {
    return request({
      url: `${BASE_URL}/user-stats/login`,
      method: 'get'
    })
  }
}

// 业务组管理 API
export const businessGroupApi = {
  // 获取业务组列表（获取顶层组）
  getGroupList() {
    return request({
      url: `${CMDB_BASE_URL}/grouplist`,
      method: 'get'
    })
  },

  // 获取业务组树形统计（包含API数量、健康状态等）
  getGroupTreeStats() {
    return request({
      url: `${BASE_URL}/groups-tree-api-stats`,
      method: 'get'
    })
  }
}

// API 业务功能管理 API
export const apiBusinessFunctionApi = {
  // 获取业务功能列表
  getList(params) {
    return request({
      url: `${BASE_URL}/business-functions`,
      method: 'get',
      params
    })
  },

  // 创建业务功能
  create(data) {
    return request({
      url: `${BASE_URL}/api-business-functions`,
      method: 'post',
      data
    })
  },

  // 更新业务功能
  update(id, data) {
    return request({
      url: `${BASE_URL}/api-business-functions/${id}`,
      method: 'put',
      data
    })
  },

  // 删除业务功能
  delete(id) {
    return request({
      url: `${BASE_URL}/api-business-functions/${id}`,
      method: 'delete'
    })
  },

  // 获取业务功能详情
  getDetail(id) {
    return request({
      url: `${BASE_URL}/api-business-functions/${id}`,
      method: 'get'
    })
  },

  // 启用/禁用业务功能
  updateStatus(id, data) {
    return request({
      url: `${BASE_URL}/api-business-functions/${id}/status`,
      method: 'put',
      data
    })
  }
}

// API 监控管理 API
export const apiMonitorApi = {
  // 创建API监控
  create(data) {
    return request({
      url: `${BASE_URL}/api-monitors`,
      method: 'post',
      data
    })
  },

  // 简化版创建API监控 - 只需4个必填参数
  createSimple(data) {
    return request({
      url: `${BASE_URL}/api-monitors/simple`,
      method: 'post',
      data
    })
  },

  // 更新API监控
  update(id, data) {
    return request({
      url: `${BASE_URL}/api-monitors/${id}`,
      method: 'put',
      data
    })
  },

  // 删除API监控
  delete(id) {
    return request({
      url: `${BASE_URL}/api-monitors/${id}`,
      method: 'delete'
    })
  },

  // 获取监控详情
  getDetail(id) {
    return request({
      url: `${BASE_URL}/api-monitors/${id}`,
      method: 'get'
    })
  },

  // 获取监控列表
  getList(params) {
    return request({
      url: `${BASE_URL}/api-monitors`,
      method: 'get',
      params
    })
  },

  // 启用/禁用监控
  updateStatus(id, status) {
    return request({
      url: `${BASE_URL}/api-monitors/${id}/status`,
      method: 'put',
      params: { status }
    })
  },

  // 获取总体统计
  getStats() {
    return request({
      url: `${BASE_URL}/api-monitors/stats`,
      method: 'get'
    })
  },

  // 按业务组统计
  getGroupStats() {
    return request({
      url: `${BASE_URL}/api-monitors/group-stats`,
      method: 'get'
    })
  },

  // 按业务功能统计
  getBusinessFunctionStats() {
    return request({
      url: `${BASE_URL}/api-monitors/business-function-stats`,
      method: 'get'
    })
  },

  // 测试API接口
  test(data) {
    return request({
      url: `${BASE_URL}/api-monitors/test`,
      method: 'post',
      data
    })
  },

  // 获取指定业务功能的所有API
  getByBusinessFunction(businessFunctionId) {
    return request({
      url: `${BASE_URL}/api-monitors-by-function`,
      method: 'get',
      params: { businessFunctionId }
    })
  },

  // ==================== 分层监控相关接口 ====================

  // 获取顶级监控列表（支持分页/过滤）
  getTopLevel(params) {
    return request({
      url: `${BASE_URL}/api-monitors/top-level`,
      method: 'get',
      params
    })
  },

  // 获取监控的树形结构
  getHierarchy(id, depth) {
    return request({
      url: `${BASE_URL}/api-monitors/${id}/hierarchy`,
      method: 'get',
      params: { depth }
    })
  },

  // 为父监控添加子监控
  addChild(parentId, childMonitorId) {
    return request({
      url: `${BASE_URL}/api-monitors/${parentId}/add-child`,
      method: 'post',
      data: { childMonitorId }
    })
  },

  // 移除子监控与父监控的关系
  removeParent(childId) {
    return request({
      url: `${BASE_URL}/api-monitors/${childId}/parent`,
      method: 'delete'
    })
  },

  // 删除监控（支持处理子监控）
  deleteWithStrategy(id, strategy) {
    return request({
      url: `${BASE_URL}/api-monitors/${id}`,
      method: 'delete',
      params: { strategy }
    })
  }
}

// 获取业务功能列表 API
export const apiBusinessFunctionListApi = {
  // 获取指定业务组的业务功能列表
  getByGroup(groupId) {
    return request({
      url: `${BASE_URL}/business-functions-by-group`,
      method: 'get',
      params: { groupId }
    })
  }
}

export default {
  domain: domainApi,
  incident: incidentApi,
  sslCert: sslCertApi,
  aliyunConfig: aliyunConfigApi,
  hostMonitor: hostMonitorApi,
  robot: robotApi,
  getUserStats: userStatsApi.getUserStats,
  businessGroup: businessGroupApi,
  apiBusinessFunction: apiBusinessFunctionApi,
  apiBusinessFunctionList: apiBusinessFunctionListApi,
  apiMonitor: apiMonitorApi
}
