import request from '@/utils/request'

export function CreateTemplate(data) {
  return request({
    url: 'template/add',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'application/json'
    }
  })
}

export function GetAllTemplates(params) {
  return request({
    url: 'template/list',
    method: 'get',
    params,
    headers: {
      'Accept': 'application/json'
    }
  })
}

export function UpdateTemplate(data) {
  return request({
    url: `template/update?id=${data.id}`,
    method: 'put',
    data: {
      name: data.name,
      type: data.type,
      content: data.content,
      remark: data.remark
    }
  })
}

export function DeleteTemplate(params) {
  return request({
    url: 'template/delete',
    method: 'delete',
    params
  })
}
// 获取模板
export function GetTemplateByID(params) {
  return request({
    url: 'template/info/' + params.id,
    method: 'get'
  })
}

export function GetTemplateContent(params) {
  return request({
    url: 'template/content/' + params.id,
    method: 'get',
    headers: {
      'Accept': 'text/plain'
    }
  })
}

export function GetTemplatesByName(params) {
  console.log('查询模板名称参数:', params)
  return request({
    url: `template/query/name?name=${encodeURIComponent(params.name || '')}`,
    method: 'get',
    params: {
      pageNum: params.pageNum,
      pageSize: params.pageSize
    },
    headers: {
      'Accept': 'application/json',
      'Authorization': localStorage.getItem('token') || ''
    }
  })
}

export function GetTemplatesByType(params) {
  console.log('查询模板类型参数:', params)
  return request({
    url: `template/query/type?type=${params.type}`,
    method: 'get',
    params: {
      pageNum: params.pageNum,
      pageSize: params.pageSize
    },
    headers: {
      'Accept': 'application/json',
      'Authorization': localStorage.getItem('token') || ''
    }
  })
}

// 任务管理API
export function CreateTask(data) {
  return request({
    url: 'task/add',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'application/json'
    }
  })
}

export function GetTaskByID(params) {
  return request({
    url: 'task/get',
    method: 'get',
    params
  })
}

export function UpdateTask(data) {
  return request({
    url: 'task/update',
    method: 'put',
    data
  })
}

export function DeleteTask(data) {
  return request({
    url: 'task/delete',
    method: 'delete',
    data,
    headers: {
      'Content-Type': 'application/json'
    }
  })
}

export function ListTasks(params) {
  return request({
    url: 'task/list',
    method: 'get',
    params
  })
}

export function GetTasksByName(params) {
  return request({
    url: 'task/query/name',
    method: 'get',
    params
  })
}

export function GetTasksByType(params) {
  return request({
    url: 'task/query/type',
    method: 'get',
    params
  })
}

export function GetTasksByStatus(params) {
  return request({
    url: 'task/query/status',
    method: 'get',
    params
  })
}

export function GetNextExecutionTime(params) {
  return request({
    url: 'task/next-execution',
    method: 'get',
    params
  })
}

// 新增获取任务模板接口
export function GetTaskTemplates(params) {
  console.log('获取任务模板参数:', params)
  return request({
    url: 'task/templates',
    method: 'get',
    params: {
      id: params.id
    },
    headers: {
      'Accept': 'application/json',
      'Authorization': localStorage.getItem('token') || ''
    }
  }).catch(error => {
    console.error('获取任务模板失败:', error)
    throw error
  })
}

export function GetTaskJobLog(params) {
  return request({
    url: 'taskjob/log',
    method: 'get',
    params: {
      taskId: params.taskId,
      templateId: params.templateId
    },
    headers: {
      'Accept': 'application/json'
    }
  })
}

export function StartJob(data) {
  return request({
    url: `taskjob/start?taskId=${data.taskId}`,
    method: 'post',
    headers: {
      'accept': 'application/json'
    }
  })
}

export function StopJob(params) {
  return request({
    url: 'taskjob/stop',
    method: 'post',
    params: {
      taskId: params.taskId,
      templateId: params.templateId
    },
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      'Accept': 'application/json'
    }
  })
}

// 获取Ansible配置列表
export function GetAnsibleConfigList(params) {
  return request({
    url: '/config/ansible',
    method: 'get',
    params
  })
}

// 创建Ansible配置
export function CreateAnsibleConfig(data) {
  return request({
    url: '/config/ansible',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'application/json'
    }
  })
}

// 更新Ansible配置
/**
 * 
 * body参数示例：
 {
  "content": "string",
  "name": "string",
  "remark": "string",
  "type": 0
}
 */
export function UpdateAnsibleConfig(data) {
  return request({
    url: `/config/ansible/${data.id}`,
    method: 'put',
    data,
    headers: {
      'Content-Type': 'application/json'
    }
  })
}

// 获取Ansible配置详情
export function GetAnsibleConfigById(id) {
  return request({
    url: `/config/ansible/${id}`,
    method: 'get'
  })
}

// 删除Ansible配置
export function DeleteAnsibleConfig(id) {
  return request({
    url: `/config/ansible/${id}`,
    method: 'delete'
  })
}

// Ansible任务管理API
export function GetAnsibleTaskList(params) {
  return request({
    url: 'task/ansiblelist',
    method: 'get',
    params
  })
}

export function CreateAnsibleTask(data) {
  return request({
    url: 'task/ansible',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export function UpdateAnsibleTask(data) {
  return request({
    url: `task/ansible/${data.id}`,
    method: 'put',
    data,
    headers: {
      'Content-Type': 'application/json'
    }
  })
}

export function GetAnsibleTaskById(id) {
  return request({
    url: `task/ansible/${id}`,
    method: 'get'
  })
}

export function StartAnsibleTask(id) {
  return request({
    url: `task/ansible/${id}/start`,
    method: 'post'
  })
}

/**
 * 建立Ansible任务实时日志的SSE连接
 * @param {string|number} id - 任务ID
 * @param {string|number} workId - 执行ID
 * @returns {EventSource} 返回EventSource实例，调用者需自行添加 onmessage 监听和 close 处理
 */
export function GetAnsibleTaskLog(id, workId, retryCount = 0) {
  const maxRetries = 1 // 减少重试次数，避免过长等待
  // 根据后端处理时间调整：首次30秒，重试3分钟
  const timeout = retryCount === 0 ? 30000 : 180000

  return request({
    url: `/task/ansible/${id}/log/${workId}`,
    method: 'get',
    timeout: timeout,
    params: {
      // 添加时间戳避免缓存
      t: Date.now(),
      // 请求实时日志，即使任务在运行中
      realtime: true,
      // 包含缓冲区日志
      includeBuffer: true
    }
  }).then(response => {
    // 记录实际的API调用和响应
    console.log('🔍 日志API调用详情:', {
      taskId: id,
      workId: workId,
      requestUrl: `/api/v1/task/ansible/${id}/log/${workId}`,
      timestamp: new Date().toISOString(),
      responseSize: response.data ? JSON.stringify(response.data).length : 0,
      responsePreview: response.data ? JSON.stringify(response.data).substring(0, 100) + '...' : 'null',
      retryCount: retryCount,
      timeout: timeout / 1000 + 's'
    })
    return response
  }).catch(error => {
    console.error('❌ 获取历史日志失败:', error)

    // 如果是超时错误且未达到最大重试次数，进行重试


    // 抛出带有更多上下文信息的错误
    const timeoutInfo = retryCount === 0 ? '30秒' : '3分钟'
    const enhancedError = {
      ...error,
      contextMessage: `后端处理超过${timeoutInfo}，可能正在执行长时间操作`
    }
    throw enhancedError
  })
}

// 获取Ansible任务日志历史记录列表
export function GetAnsibleTaskHistory(params) {
  return request({
    url: `/api/v1/task/ansible/${params.id}/history`,
    method: 'get',
    params: {
      page: params.page,
      limit: params.pageSize,
    }
  })
}

// 根据任务ID、WORKID和HistoryID获取历史任务日志
export function GetAnsibleTaskLogByHistory(params) {
  return request({
    url: `/api/v1/task/ansible/history/detail/task/${params.id}/work/${params.workId}/history/${params.historyId}/log`,
    method: 'get',
    params: {
      t: Date.now(),
      // 获取历史日志时不需要实时参数
      realtime: false,
      includeBuffer: false
    }
  })
}

// 获取任务的历史执行详情，包含每个主机的执行日志
export function GetAnsibleHistoryDetail(params) {
  return request({
    url: `/api/v1/task/ansible/${params.id}/history/${params.historyId}`,
    method: 'get',
    headers: {
      'Accept': 'application/json'
    }
  })
}

// 尝试直接获取日志文件内容的备选方法
export function GetAnsibleTaskLogDirect(id, workId) {
  console.log('🔧 尝试直接日志文件访问方法')
  return request({
    url: `/task/ansible/${id}/log/${workId}/direct`,
    method: 'get',
    timeout: 10000, // 10秒超时
    params: {
      t: Date.now(),
      // 强制读取最新日志
      tail: true,
      // 获取最后1000行
      lines: 1000,
      // 绕过缓存
      nocache: true
    }
  }).catch(error => {
    console.warn('直接日志访问失败，使用标准方法:', error.message)
    throw error
  })
}

// 获取有效token的工具函数
function getValidToken() {
  // 尝试多种存储位置获取token
  const storageKeys = ['token', 'access_token', 'jwt_token', 'authToken']
  let token = null

  // 优先从localStorage获取
  for (const key of storageKeys) {
    token = localStorage.getItem(key)
    if (token && token !== 'null' && token !== 'undefined') {
      break
    }
  }

  // 如果localStorage没有，尝试sessionStorage
  if (!token || token === 'null' || token === 'undefined') {
    for (const key of storageKeys) {
      token = sessionStorage.getItem(key)
      if (token && token !== 'null' && token !== 'undefined') {
        break
      }
    }
  }

  // 如果token是JSON对象，尝试解析
  if (token && typeof token === 'string' && token.startsWith('{')) {
    try {
      const tokenObj = JSON.parse(token)
      token = tokenObj.token || tokenObj.access_token || tokenObj.value || tokenObj.jwt
    } catch (e) {
      console.warn('无法解析token JSON:', e)
    }
  }

  // 确保token不是字符串'null'
  if (token === 'null' || token === 'undefined' || !token) {
    token = null
  }

  console.log('Token获取结果:', {
    hasToken: !!token,
    tokenPreview: token ? `${token.substring(0, 10)}...` : 'null',
    tokenLength: token ? token.length : 0
  })

  return token
}

// SSE实时日志流接口
export function GetAnsibleTaskLogStream(id, workId) {
  const token = getValidToken()
  // 使用当前页面的协议和主机，支持Docker部署
  const protocol = window.location.protocol
  const host = window.location.host
  const baseURL = `${protocol}//${host}`
  const url = `${baseURL}/api/v1/task/ansible/${id}/log/${workId}`

  if (!token) {
    console.error('警告: 未找到有效的认证token')
  }

  console.log('构造SSE URL:', {
    hasValidToken: !!token,
    tokenPreview: token ? `${token.substring(0, 10)}...` : 'null',
    baseURL,
    id,
    workId,
    finalUrl: `${url}?token=${encodeURIComponent(token || '')}`
  })

  return {
    url: `${url}?token=${encodeURIComponent(token || '')}`
  }
}

// 获取Ansible任务详情和Works列表
export function GetAnsibleTaskDetail(id) {
  return request({
    url: `task/ansible/${id}`,
    method: 'get',
    headers: {
      'Accept': 'application/json'
    }
  })
}

// 启动Ansible任务
export function StartAnsibleTaskFlow(id) {
  return request({
    url: `task/ansible/${id}/start`,
    method: 'post',
    headers: {
      'Accept': 'application/json'
    }
  })
}

export function DeleteAnsibleTask(id) {
  return request({
    url: `task/ansible/${id}`,
    method: 'delete',
    headers: {
      'Accept': 'application/json'
    }
  })
}

// 根据名称模糊查询Ansible任务
export function GetAnsibleTasksByName(params) {
  console.log('查询Ansible任务名称参数:', params)
  return request({
    url: `task/ansible/query/name?name=${encodeURIComponent(params.name || '')}`,
    method: 'get',
    params: {
      page: params.page,
      pageSize: params.pageSize
    },
    headers: {
      'Accept': 'application/json',
      'Authorization': localStorage.getItem('token') || ''
    }
  }).catch(error => {
    console.error('按名称查询Ansible任务失败:', error)
    throw error
  })
}

// 根据类型查询Ansible任务
export function GetAnsibleTasksByType(params) {
  console.log('查询Ansible任务类型参数:', params)
  return request({
    url: `task/ansible/query/type?type=${params.type}`,
    method: 'get',
    params: {
      page: params.page,
      pageSize: params.pageSize
    },
    headers: {
      'Accept': 'application/json',
      'Authorization': localStorage.getItem('token') || ''
    }
  }).catch(error => {
    console.error('按类型查询Ansible任务失败:', error)
    throw error
  })
}

// 暂停定时任务
export function PauseScheduledTask(taskId) {
  return request({
    url: 'task/monitor/scheduled/pause',
    method: 'post',
    params: {
      task_id: taskId
    },
    headers: {
      'Accept': 'application/json'
    }
  })
}

// 恢复定时任务
export function ResumeScheduledTask(taskId) {
  return request({
    url: 'task/monitor/scheduled/resume',
    method: 'post',
    params: {
      task_id: taskId
    },
    headers: {
      'Accept': 'application/json'
    }
  })
}
