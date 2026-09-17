import request from '@/utils/request'

export const licenseApi = {
  getMachineCode() {
    return request({
      url: '/license/machine-code',
      method: 'get'
    })
  },

  getStatus() {
    return request({
      url: '/license/status',
      method: 'get'
    })
  },

  activate(data) {
    return request({
      url: '/license/activate',
      method: 'post',
      data
    })
  }
}
