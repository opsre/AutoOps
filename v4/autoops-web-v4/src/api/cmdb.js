import request from "@/utils/request"

export default {
    // 获取主机监控全量数据
    getHostMonitorData(params) {
        console.log('请求监控数据参数:', JSON.stringify(params, null, 2))
        const { id, duration, start, end } = params
        let url = `monitor/hosts/${id}/all-metrics`
        const queryParams = {}
        
        if (duration) {
            queryParams.duration = duration
        } else if (start && end) {
            queryParams.start = start
            queryParams.end = end
        }
        
        return request({
            url,
            method: 'get',
            params: queryParams
        })
    },

    // 批量获取主机简要监控数据
    getHostsMonitorData(ids) {
        console.log('请求批量监控数据，主机IDs:', ids)
        return request({
            url: 'monitor/hosts',
            method: 'get',
            params: { ids }
        })
    },

    // 获取主机TOP进程监控数据
    getHostTopProcesses(hostId) {
        return request({
            url: `monitor/hosts/${hostId}/top-processes`,
            method: 'get'
        })
    },

    // 获取主机TCP端口监听状态
    getHostTcpPorts(hostId) {
        return request({
            url: `monitor/hosts/${hostId}/ports`,
            method: 'get'
        })
    },    // 获取主机列表
    getHostList(params) {
        return request({
            url: 'cmdb/host/list',
            method: 'get',
            params
        })
    },
    getAllCmdbGroups() {
        return request({
            url: 'cmdb/grouplist',
            method: 'get'
        })
    },
    getCmdbGroupByName(name) {
        return request({
            url: 'cmdb/groupbyname',
            method: 'get',
            params: {
                name: name
            }
        })
    },
    createCmdbGroup(data) {
        return request({
            url: 'cmdb/groupadd',
            method: 'post',
            data: data
        })
    },
    updateCmdbGroup(data) {
        return request({
            url: 'cmdb/groupupdate',
            method: 'put',
            data: data
        })
    },
    deleteCmdbGroup(id) {
        const data = {
            id
        }
        return request({
            url: 'cmdb/groupdelete',
            method: 'delete',
            data: data
        })
    },

    // 主机管理
    createCmdbHost(data) {
        return request({
            url: 'cmdb/hostcreate',
            method: 'post',
            data: data
        })
    },
    updateCmdbHost(data) {
        return request({
            url: 'cmdb/hostupdate',
            method: 'put',
            data: data
        })
    },
    deleteCmdbHost(id) {
        return request({
            url: 'cmdb/hostdelete',
            method: 'delete',
            data: { id }
        })
    },
    // 批量删除主机
    batchDeleteCmdbHost(hostIds) {
        return request({
            url: 'cmdb/hostbatchdelete',
            method: 'delete',
            data: { hostIds }
        })
    },
    getCmdbHostList(params) {
        // 确保参数名与后端API一致
        const queryParams = {
            page: params.page || 1,
            pageSize: params.pageSize || 10,
            sortBy: 'update_time',        // 按更新时间排序
            sortOrder: 'desc'    // 降序排列（最新的在前面）
        }

        // 添加搜索条件，确保参数名与后端API一致
        if (params.name) queryParams.name = params.name
        if (params.ip) queryParams.ip = params.ip
        if (params.status) queryParams.status = params.status
        if (params.groupId) queryParams.groupId = params.groupId

        console.log('最终API查询参数:', JSON.stringify(queryParams, null, 2))
        return request({
            url: 'cmdb/hostlist',
            method: 'get',
            params: queryParams
        })
    },
    getCmdbHostById(id) {
        return request({
            url: 'cmdb/hostinfo',
            method: 'get',
            params: { id }
        })
    },
    getCmdbHostsByGroupId(groupId) {
        return request({
            url: 'cmdb/hostgroup',
            method: 'get',
            params: { groupId }
        })
    },
    GetCmdbHostsByHostNameLike(hostName, params) {
        return request({
            url: '/cmdb/hostbyname',
            method: 'get',
            params: {
                name: hostName,
                sortBy: 'id',
                sortOrder: 'desc',
                ...params
            }
        })
    },
    GetCmdbHostsByIP(ip, params) {
        return request({
            url: 'cmdb/hostbyip',
            method: 'get',
            params: {
                ip,
                sortBy: 'id',
                sortOrder: 'desc',
                ...params
            }
        })
    },
    GetCmdbHostsByStatus(status, params) {
        return request({
            url: 'cmdb/hostbystatus',
            method: 'get',
            params: {
                status,
                sortBy: 'id',
                sortOrder: 'desc',
                ...params
            }
        })
    },
    getCmdbHostsByGroupId(groupId, params) {
        return request({
            url: 'cmdb/hostgroup',
            method: 'get',
            params: {
                groupId,
                sortBy: 'id',
                sortOrder: 'desc',
                ...params
            }
        })
    },
    hostcloudcreatealiyun(data) {
        return request({
            url: 'cmdb/hostcloudcreatealiyun',
            method: 'post',
            data: data
        })
    },
    hostcloudcreatetencent(data) {
        return request({
            url: 'cmdb/hostcloudcreatetencent',
            method: 'post',
            data: data
        })
    },
    hostcloudcreatehuawei(data) {
        return request({
            url: 'cmdb/hostcloudcreatehuawei',
            method: 'post',
            data: data
        })
    },
    hostcloudcreateaws(data) {
        return request({
            url: 'cmdb/hostcloudcreateaws',
            method: 'post',
            data: data
        })
    },
    // 获取云主机统计信息
    getHostCloudStatistics() {
        return request({
            url: 'cmdb/hostcloudstatistics',
            method: 'get'
        })
    },

    // 获取带主机列表的分组数据
    getGroupListWithHosts() {
        return request({
            url: 'cmdb/grouplistwithhosts',
            method: 'get'
        })
    },

    // WebSocket连接
    getHostSSHWebSocketUrl(hostId) {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        return `${protocol}//${window.location.host}/api/v1/cmdb/hostssh/connect/${hostId}`
    },

    // 文件上传到主机
    uploadFileToHost(hostId, formData, config = {}) {
        return request({
            url: `cmdb/hostssh/upload/${hostId}`,
            method: 'post',
            data: formData,
            headers: {
                'Content-Type': 'multipart/form-data'
            },
            ...config
        })
    },

    // 查询主机文件上传任务状态
    getUploadFileTaskStatus(taskId, config = {}) {
        return request({
            url: `cmdb/hostssh/upload/task/${taskId}`,
            method: 'get',
            ...config
        })
    },

    // 从主机下载文件
    downloadFileFromHost(hostId, remotePath, config = {}) {
        return request({
            url: `cmdb/hostssh/download/${hostId}`,
            method: 'get',
            params: {
                remotePath: remotePath
            },
            responseType: 'blob',
            timeout: 600000, // 10分钟超时，支持大文件下载
            ...config
        })
    },

    // 列出主机远程目录文件
    listRemoteFiles(hostId, remotePath = '/') {
        return request({
            url: `cmdb/hostssh/listfiles/${hostId}`,
            method: 'get',
            params: {
                remotePath: remotePath
            }
        })
    },

    // SSH执行命令
    executeHostCommand(hostId, command) {
        return request({
            url: `cmdb/hostssh/command/${hostId}`,
            method: 'get',
            params: {
                command: command
            }
        })
    },

    // 数据库管理
    createDatabase(data) {
        return request({
            url: 'cmdb/database',
            method: 'post',
            data: data
        })
    },
    updateDatabase(data) {
        return request({
            url: 'cmdb/database',
            method: 'put',
            data: data
        })
    },
    deleteDatabase(id) {
        return request({
            url: 'cmdb/database',
            method: 'delete',
            data: { id }
        })
    },
    getDatabase(id) {
        return request({
            url: 'cmdb/database/info',
            method: 'get',
            params: { id }
        })
    },
    listDatabases(params = {}) {
        return request({
            url: 'cmdb/databaselist',
            method: 'get',
            params: params
        })
    },
    getDatabasesByName(name, params = {}) {
        return request({
            url: 'cmdb/database/byname',
            method: 'get',
            params: {
                name,
                ...params
            }
        })
    },
    getDatabasesByType(type, params = {}) {
        return request({
            url: 'cmdb/database/bytype',
            method: 'get',
            params: {
                type,
                ...params
            }
        })
    },

    // SQL查询
    executeSelectSQL(data) {
        console.log('executeSelectSQL请求数据:', JSON.stringify(data, null, 2))
        return request({
            url: 'cmdb/sql/select',
            method: 'post',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },
    // SQL插入
    executeInsertSQL(data) {
        console.log('executeInsertSQL请求数据:', JSON.stringify(data, null, 2))
        return request({
            url: 'cmdb/sql',
            method: 'post',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },
    // SQL更新
    executeUpdateSQL(data) {
        console.log('executeUpdateSQL请求数据:', JSON.stringify(data, null, 2))
        return request({
            url: 'cmdb/sql',
            method: 'put',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },
    executeDeleteSQL(data) {
        console.log('executeDeleteSQL请求数据:', JSON.stringify(data, null, 2))
        return request({
            url: 'cmdb/sql',
            method: 'delete',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },
    executeRawSQL(data) {
        console.log('executeRawSQL请求数据:', JSON.stringify(data, null, 2))
        return request({
            url: 'cmdb/sql/execute',
            method: 'post',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },
    // 获取数据库列表
    executeDatabase(data) {
        console.log('executeDatabase请求数据:', JSON.stringify(data, null, 2))
        return request({
            url: 'cmdb/sql/databaselist',
            method: 'post',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },

    // SQL日志管理
    GetCmdbSqlLogList(params) {
        return request({
            url: 'cmdb/sqlLog/list',
            method: 'get',
            params: params
        })
    },
    DeleteCmdbSqlLogById(id) {
        return request({
            url: 'cmdb/sqlLog/delete',
            method: 'delete',
            data: { id }
        })
    },
    CleanCmdbSqlLog() {
        return request({
            url: 'cmdb/sqlLog/clean',
            method: 'delete'
        })
    },

    // 从Excel导入主机
    ImportHostsFromExcel(data, config = {}) {
        return request({
            url: 'cmdb/hostimport',
            method: 'post',
            data: data,
            headers: {
                'Content-Type': 'multipart/form-data'
            },
            ...config
        })
    },

    // 下载主机导入模板
    DownloadHostTemplate() {
        return request({
            url: 'cmdb/hosttemplate',
            method: 'get',
            responseType: 'blob'
        })
    },

    // 同步主机配置信息
    syncHostConfig(id) {
        return request({
            url: 'cmdb/hostsync',
            method: 'post',
            data: {
                id: id
            },
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },

    // Agent管理
    // 部署agent到指定主机 (支持单个或批量)
    deployAgent(hostIds, version = '1.0.0') {
        // 确保hostIds是数组格式
        const ids = Array.isArray(hostIds) ? hostIds : [hostIds]
        return request({
            url: 'monitor/agent/deploy',
            method: 'post',
            data: {
                hostIds: ids,
                version: version
            }
        })
    },
    
    // 卸载指定主机的agent (支持单个或批量)
    uninstallAgent(hostIds) {
        // 确保hostIds是数组格式
        const ids = Array.isArray(hostIds) ? hostIds : [hostIds]
        return request({
            url: 'monitor/agent/uninstall',
            method: 'delete',
            data: {
                hostIds: ids
            }
        })
    },
    
    // 根据主机id获取agent状态
    getAgentStatus(id) {
        return request({
            url: `monitor/agent/status/${id}`,
            method: 'get'
        })
    },
    
    // 重启agent
    restartAgent(id) {
        return request({
            url: `monitor/agent/restart/${id}`,
            method: 'post'
        })
    },
    
    // 获取agent列表
    getAgentList(params) {
        return request({
            url: 'monitor/agent/list',
            method: 'get',
            params
        })
    },

    // 删除agent数据 (用于离线服务器)
    deleteAgent(id) {
        return request({
            url: `monitor/agent/delete/${id}`,
            method: 'delete'
        })
    },

    // 批量更新主机分组
    batchUpdateHostGroup(hostIds, groupId) {
        return request({
            url: 'cmdb/hostbatchupdategroup',
            method: 'put',
            data: {
                hostIds: hostIds,
                groupId: groupId
            }
        })
    },

    // SNMP设备管理
    // 获取SNMP设备列表
    getSnmpDevices(config) {
        return request({
            url: 'cmdb/snmp-devices',
            method: 'get',
            ...config
        })
    },

    // 根据ID获取SNMP设备详情
    getSnmpDeviceById(id) {
        return request({
            url: `cmdb/snmp-devices/${id}`,
            method: 'get'
        })
    },

    // 创建SNMP设备
    createSnmpDevice(data) {
        return request({
            url: 'cmdb/snmp-devices',
            method: 'post',
            data: data
        })
    },

    // 更新SNMP设备
    updateSnmpDevice(id, data) {
        return request({
            url: `cmdb/snmp-devices/${id}`,
            method: 'put',
            data: data
        })
    },

    // 删除SNMP设备
    deleteSnmpDevice(id) {
        return request({
            url: `cmdb/snmp-devices/${id}`,
            method: 'delete'
        })
    },

    // 立即采集SNMP设备数据
    collectSnmpDevice(id) {
        return request({
            url: `cmdb/snmp-devices/${id}/collect`,
            method: 'post'
        })
    },

    // 启用SNMP设备监控
    enableSnmpDeviceMonitor(id) {
        return request({
            url: `cmdb/snmp-devices/${id}/enable`,
            method: 'post'
        })
    },

    // 停用SNMP设备监控
    disableSnmpDeviceMonitor(id) {
        return request({
            url: `cmdb/snmp-devices/${id}/disable`,
            method: 'post'
        })
    },

    // SNMP设备SSH终端连接
    getSnmpDeviceSSHWebSocketUrl(id) {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const token = localStorage.getItem('token') || sessionStorage.getItem('token') || ''
        return `${protocol}//${window.location.host}/api/v1/cmdb/snmp-devices/${id}/ssh/connect?token=${token}`
    },

    // SNMP设备SSH执行命令
    executeSnmpDeviceCommand(id, command) {
        return request({
            url: `cmdb/snmp-devices/${id}/ssh/command`,
            method: 'get',
            params: {
                command: command
            }
        })
    },

    // Windows Exporter 管理
    // 查询 Windows Exporter 安装状态
    getWindowsExporterStatus(deviceId) {
        return request({
            url: `cmdb/windows-exporter/status/${deviceId}`,
            method: 'get'
        })
    },

    // 安装 Windows Exporter
    installWindowsExporter(data) {
        return request({
            url: 'cmdb/windows-exporter/install',
            method: 'post',
            data: data
        })
    },

    // 卸载 Windows Exporter
    uninstallWindowsExporter(deviceId) {
        return request({
            url: `cmdb/windows-exporter/uninstall/${deviceId}`,
            method: 'delete'
        })
    },

    // 获取SNMP设备监控数据
    getSnmpDeviceMonitorData(id) {
        return request({
            url: `cmdb/snmp-devices/${id}/monitor`,
            method: 'get'
        })
    },

    // 获取 Windows 设备监控数据（Prometheus windows_exporter 指标聚合）
    getWindowsDeviceMonitor(id) {
        return request({
            url: `cmdb/windows-devices/${id}/monitor`,
            method: 'get'
        })
    },

    // 批量获取SNMP设备资源使用情况
    getSnmpDevicesBatchResource(deviceIds) {
        return request({
            url: 'cmdb/snmp-devices/batch-resource',
            method: 'get',
            params: {
                deviceIds: deviceIds
            }
        })
    },

    // 批量获取Windows设备实时CPU和内存数据
    getWindowsDevicesBatchRealtime(deviceIds) {
        return request({
            url: 'cmdb/windows-devices/batch-realtime',
            method: 'get',
            params: {
                ids: deviceIds
            }
        })
    },

    // JumpServer主机同步
    syncHostFromJumpServer(data) {
        return request({
            url: 'cmdb/hostjumpserversync',
            method: 'post',
            data: data
        })
    }
}
