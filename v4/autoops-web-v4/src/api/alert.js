/**
 * 告警管理 API
 *
 * @author xiaoRui
 */

import request from "@/utils/request"

const BASE_URL = '/api/v1/monitor/alert'

// ==================== 告警规则模板 API ====================
export const templateApi = {
  // 获取模板列表
  getList(params) {
    return request({
      url: `${BASE_URL}/templates`,
      method: 'get',
      params
    })
  },

  // 创建模板
  create(data) {
    return request({
      url: `${BASE_URL}/templates`,
      method: 'post',
      data
    })
  },

  // 更新模板
  update(data) {
    return request({
      url: `${BASE_URL}/templates/${data.id}`,
      method: 'put',
      data
    })
  },

  // 删除模板
  delete(id) {
    return request({
      url: `${BASE_URL}/templates/${id}`,
      method: 'delete'
    })
  },

  // 批量删除模板
  batchDelete(ids) {
    return request({
      url: `${BASE_URL}/templates/batch-delete`,
      method: 'post',
      data: ids
    })
  },

  // 更新模板状态（启用/禁用）
  updateStatus(id, status) {
    return request({
      url: `${BASE_URL}/templates/${id}/status`,
      method: 'put',
      params: { status }
    })
  },

  // 获取模板分类列表
  getCategories() {
    return request({
      url: `${BASE_URL}/templates/categories`,
      method: 'get'
    })
  },

  // 获取模板详情
  getDetail(id) {
    return request({
      url: `${BASE_URL}/templates/${id}`,
      method: 'get'
    })
  }
}

// ==================== 告警策略 API ====================
export const policyApi = {
  // 获取策略列表
  getList(params) {
    return request({
      url: `${BASE_URL}/policies`,
      method: 'get',
      params
    })
  },

  // 创建策略
  create(data) {
    return request({
      url: `${BASE_URL}/policies`,
      method: 'post',
      data
    })
  },

  // 更新策略
  update(data) {
    return request({
      url: `${BASE_URL}/policies/${data.id}`,
      method: 'put',
      data
    })
  },

  // 删除策略
  delete(id) {
    return request({
      url: `${BASE_URL}/policies/${id}`,
      method: 'delete'
    })
  },

  // 批量删除策略
  batchDelete(ids) {
    return request({
      url: `${BASE_URL}/policies/batch-delete`,
      method: 'post',
      data: ids
    })
  },

  // 更新策略状态（启用/禁用）
  updateStatus(id, status) {
    return request({
      url: `${BASE_URL}/policies/${id}/status`,
      method: 'put',
      params: { status }
    })
  },

  // PromQL 预览
  preview(data) {
    return request({
      url: `${BASE_URL}/policies/preview`,
      method: 'post',
      data
    })
  },

  // 获取策略详情
  getDetail(id) {
    return request({
      url: `${BASE_URL}/policies/${id}`,
      method: 'get'
    })
  }
}

// ==================== 告警记录 API ====================
export const recordApi = {
  // 获取告警记录列表
  getList(params) {
    return request({
      url: `${BASE_URL}/records`,
      method: 'get',
      params
    })
  },

  // 获取告警记录详情
  getDetail(id) {
    return request({
      url: `${BASE_URL}/records/${id}`,
      method: 'get'
    })
  },

  // 认领告警
  acknowledge(id, data) {
    return request({
      url: `${BASE_URL}/records/${id}/acknowledge`,
      method: 'post',
      data
    })
  },

  // 统一处理告警
  process(id, data) {
    return request({
      url: `${BASE_URL}/records/${id}/process`,
      method: 'put',
      data
    })
  },

  // 处理告警
  handle(id, data) {
    return request({
      url: `${BASE_URL}/records/${id}/handle`,
      method: 'post',
      data
    })
  },

  // 获取告警统计
  getStats() {
    return request({
      url: `${BASE_URL}/records/stats`,
      method: 'get'
    })
  },

  // 删除单条告警记录
  delete(id) {
    return request({
      url: `${BASE_URL}/records/${id}`,
      method: 'delete'
    })
  },

  // 批量删除告警记录
  batchDelete(ids) {
    return request({
      url: `${BASE_URL}/records/batch-delete`,
      method: 'post',
      data: { ids }
    })
  },

  // 清空所有告警记录
  clearAll() {
    return request({
      url: `${BASE_URL}/records/clear-all`,
      method: 'delete'
    })
  }
}

// ==================== 告警策略规则 API ====================
export const policyRuleApi = {
  // 获取策略下的规则列表
  getList(params) {
    return request({
      url: `${BASE_URL}/policy-rules`,
      method: 'get',
      params
    })
  },

  // 从模板创建策略规则
  create(data) {
    return request({
      url: `${BASE_URL}/policy-rules`,
      method: 'post',
      data
    })
  },

  // 更新策略规则
  update(data) {
    return request({
      url: `${BASE_URL}/policy-rules/${data.id}`,
      method: 'put',
      data
    })
  },

  // 删除策略规则
  delete(id) {
    return request({
      url: `${BASE_URL}/policy-rules/${id}`,
      method: 'delete'
    })
  },

  // 批量删除规则
  batchDelete(ids) {
    return request({
      url: `${BASE_URL}/policy-rules/batch-delete`,
      method: 'post',
      data: ids
    })
  },

  // 获取规则详情
  getDetail(id) {
    return request({
      url: `${BASE_URL}/policy-rules/${id}`,
      method: 'get'
    })
  },

  // 更新策略规则状态（启用/禁用）
  updateStatus(id, status) {
    return request({
      url: `${BASE_URL}/policy-rules/${id}/status`,
      method: 'put',
      params: { status }
    })
  },

  // 预览规则的实时数据
  preview(id) {
    return request({
      url: `${BASE_URL}/policy-rules/${id}/preview`,
      method: 'get'
    })
  }
}
export const ruleApi = {
  // 获取规则列表
  getList(params) {
    return request({
      url: `${BASE_URL}/rules`,
      method: 'get',
      params
    })
  },

  // 创建规则
  create(data) {
    return request({
      url: `${BASE_URL}/rules`,
      method: 'post',
      data
    })
  },

  // 更新规则
  update(data) {
    return request({
      url: `${BASE_URL}/rules/${data.id}`,
      method: 'put',
      data
    })
  },

  // 删除规则
  delete(id) {
    return request({
      url: `${BASE_URL}/rules/${id}`,
      method: 'delete'
    })
  },

  // 批量删除规则
  batchDelete(ids) {
    return request({
      url: `${BASE_URL}/rules/batch-delete`,
      method: 'post',
      data: ids
    })
  },

  // 更新规则状态
  updateStatus(id, status) {
    return request({
      url: `${BASE_URL}/rules/${id}/status`,
      method: 'put',
      params: { status }
    })
  }
}

// ==================== 默认导出 ====================
export default {
  template: templateApi,
  policy: policyApi,
  policyRule: policyRuleApi,
  record: recordApi,
  rule: ruleApi
}
