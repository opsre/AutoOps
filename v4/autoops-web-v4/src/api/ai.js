/**
 * AI模型管理API
 *
 * @author xiaoRui
 */

import request from "@/utils/request"

export default {
  // ============== AI模型管理 ==============

  // 创建AI模型配置
  createAiModel(data) {
    return request({
      url: '/api/v1/ai/model/add',
      method: 'post',
      data
    })
  },

  // 获取模型详情
  getAiModelInfo(id) {
    return request({
      url: '/api/v1/ai/model/info',
      method: 'get',
      params: { id }
    })
  },

  // 更新模型配置
  updateAiModel(data) {
    return request({
      url: '/api/v1/ai/model/update',
      method: 'put',
      data
    })
  },

  // 删除模型
  deleteAiModel(id) {
    return request({
      url: '/api/v1/ai/model/delete',
      method: 'delete',
      params: { id }
    })
  },

  // 获取模型列表(支持分页、筛选)
  getAiModelList(params) {
    return request({
      url: '/api/v1/ai/model/list',
      method: 'get',
      params
    })
  },

  // 更新模型状态(启用/禁用)
  updateAiModelStatus(id, status) {
    return request({
      url: '/api/v1/ai/model/updateStatus',
      method: 'put',
      data: { id, status }
    })
  },

  // 获取所有启用的模型
  getEnabledAiModels() {
    return request({
      url: '/api/v1/ai/model/enabled',
      method: 'get'
    })
  },

  // ============== AI对话 ==============

  // AI对话接口(支持阿里百炼等多种模型)
  // 注意：AI对话可能需要较长时间，设置180秒超时
  aiChat(data) {
    return request({
      url: '/api/v1/ai/chat',
      method: 'post',
      data,
      timeout: 180000 // 180秒超时（3分钟），适用于内容较长的AI回复
    })
  }
}
