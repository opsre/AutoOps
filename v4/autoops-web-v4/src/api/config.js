import request from '@/utils/request'

export default {
    // 账号认证管理API
    listAccountAuth(params) {
        return request({
            url: 'config/accountauth/list',
            method: 'get',
            params
        })
    },
    getAccountAuthByAlias(alias) {
        return request({
            url: 'config/accountauth/alias',
            method: 'get',
            params: { alias }
        })
    },
    getAccountAuthByType(type) {
        return request({
            url: 'config/accountauth/type',
            method: 'get',
            params: { type }
        })
    },
    // 账号认证管理
    createAccountAuth(data) {
        return request({
            url: 'config/accountauth',
            method: 'post',
            data
        })
    },
    updateAccountAuth(data) {
        return request({
            url: 'config/accountauth',
            method: 'put',
            data
        })
    },
    deleteAccountAuth(id) {
        return request({
            url: 'config/accountauth',
            method: 'delete',
            params: { id: Number(id) }
        })
    },
    decryptPassword(data, config) {
        return request({
            url: 'config/accountauth/decrypt',
            method: 'post',
            ...config,
            data
        })
    },
    // 测试账号连接
    testAccountAuth(id) {
        return request({
            url: 'config/accountauth/test',
            method: 'post',
            params: { id }
        })
    },

    // ECS凭据API
    getEcsAuthList(params) {
        return request({
            url: 'config/ecsauthlist',
            method: 'get',
            params
        })
    },
    getEcsAuthByName(name) {
        return request({
            url: 'config/ecsauthinfo',
            method: 'get',
            params: { name }
        })
    },
    createEcsAuth(data) {
        return request({
            url: 'config/ecsauthadd',
            method: 'post',
            data: data
        })
    },
    updateEcsAuth(data) {
        return request({
            url: 'config/ecsauthupdate',
            method: 'put',
            data: data
        })
    },
    deleteEcsAuth(id) {
        return request({
            url: 'config/ecsauthdelete',
            method: 'delete',
            data: { id }
        })
    },

    // 密钥管理API
    getKeyManageList(params) {
        return request({
            url: 'config/keymanage/list',
            method: 'get',
            params
        })
    },
    getKeyManageById(id) {
        return request({
            url: 'config/keymanage',
            method: 'get',
            params: { id }
        })
    },
    getKeyManageByType(keyType) {
        return request({
            url: 'config/keymanage/type',
            method: 'get',
            params: { keyType }
        })
    },
    createKeyManage(data) {
        return request({
            url: 'config/keymanage',
            method: 'post',
            data
        })
    },
    updateKeyManage(data) {
        return request({
            url: 'config/keymanage',
            method: 'put',
            data
        })
    },
    deleteKeyManage(id) {
        return request({
            url: 'config/keymanage',
            method: 'delete',
            params: { id }
        })
    },
    decryptKeys(id) {
        return request({
            url: 'config/keymanage/decrypt',
            method: 'post',
            params: { id }
        })
    },
    syncCloudHosts(data) {
        return request({
            url: 'config/keymanage/sync',
            method: 'post',
            data
        })
    },
    syncAliyunHosts(data) {
        return request({
            url: 'config/keymanage/sync/aliyun',
            method: 'post',
            data
        })
    },
    syncTencentHosts(data) {
        return request({
            url: 'config/keymanage/sync/tencent',
            method: 'post',
            data
        })
    },

    // 定时同步配置API
    getSyncScheduleList(params) {
        return request({
            url: 'config/sync-schedule/list',
            method: 'get',
            params
        })
    },
    createSyncSchedule(data) {
        return request({
            url: 'config/sync-schedule',
            method: 'post',
            data
        })
    },
    updateSyncSchedule(data) {
        return request({
            url: 'config/sync-schedule',
            method: 'put',
            data
        })
    },
    deleteSyncSchedule(id) {
        return request({
            url: 'config/sync-schedule',
            method: 'delete',
            params: { id }
        })
    },
    toggleSyncScheduleStatus(data) {
        return request({
            url: 'config/sync-schedule/toggle-status',
            method: 'post',
            data
        })
    },

    // 手动触发同步任务
    triggerSyncSchedule(id) {
        return request({
            url: 'config/sync-schedule/trigger',
            method: 'post',
            params: { id }
        })
    },

    // 获取任务下次执行时间
    getNextExecutionTime(cron) {
        return request({
            url: 'task/next-execution',
            method: 'get',
            params: { cron }
        })
    },

    // 获取同步日志
    getSyncScheduleLog(id) {
        return request({
            url: 'config/sync-schedule/log',
            method: 'get',
            params: { id }
        })
    }
}
