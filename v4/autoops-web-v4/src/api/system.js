import request from "@/utils/request"

export default {
    // post岗位
    queryPostList(params) {
        return request({
            url: 'post/list',
            method: 'get',
            params: params
        })
    },
    batchDeleteSysPost(ids) {
        const data = {
            ids
        }
        return request({
            url: 'post/batch/delete',
            method: 'delete',
            data: data
        })
    },
    deleteSysPost(id) {
        const data = {
            id
        }
        return request({
            url: 'post/delete',
            method: 'delete',
            data: data
        })
    },
    querySysPostVoList() {
        return request({
            url: 'post/vo/list',
            method: 'get'
        })
    },
    addPost(data) {
        return request({
            url: 'post/add',
            method: 'post',
            data: data
        })
    },
    postInfo(id) {
        return request({
            url: 'post/info',
            method: 'get',
            params: { id }
        })
    },
    updatePost(data) {
        return request({
            url: 'post/update',
            method: 'put',
            data: data
        })
    },
    updatePostStatus(id, postStatus) {
        return request({
            url: 'post/updateStatus',
            method: 'put',
            data: {
                id,
                postStatus
            }
        })
    },
    // dept部门
    queryDeptList(params) {
        return request({
            url: 'dept/list',
            method: 'get',
            params: params
        })
    },
    querySysDeptVoList() {
        return request({
            url: 'dept/vo/list',
            method: 'get'
        })
    },
    addDept(data) {
        return request({
            url: 'dept/add',
            method: 'post',
            data: data
        })
    },
    deleteDept(id) {
        return request({
            url: 'dept/delete',
            method: 'delete',
            data: { id }
        })
    },
    deptInfo(id) {
        return request({
            url: 'dept/info',
            method: 'get',
            params: { id }
        })
    },
    deptUsers(id) {
        return request({
            url: 'dept/users',
            method: 'get',
            params: { deptId: id }
        })
    },
    deptUpdate(data) {
        return request({
            url: 'dept/update',
            method: 'put',
            data: data
        })
    },
    // menu菜单
    queryMenuList(params) {
        return request({
            url: 'menu/list',
            method: 'get',
            params: params
        })
    },
    querySysMenuVoList() {
        return request({
            url: 'menu/vo/list',
            method: 'get'
        })
    },
    addMenu(data) {
        return request({
            url: 'menu/add',
            method: 'post',
            data: data
        })
    },
    menuInfo(id) {
        return request({
            url: 'menu/info',
            method: 'get',
            params: { id }
        })
    },
    menuUpdate(data) {
        return request({
            url: 'menu/update',
            method: 'put',
            data: data
        })
    },
    menuDelete(id) {
        return request({
            url: 'menu/delete',
            method: 'delete',
            data: { id }
        })
    },
    // role角色
    queryRoleList(params) {
        return request({
            url: 'role/list',
            method: 'get',
            params: params
        })
    },
    querySysRoleVoList() {
        return request({
            url: 'role/vo/list',
            method: 'get'
        })
    },
    addRole(data) {
        return request({
            url: 'role/add',
            method: 'post',
            data: data
        })
    },
    roleInfo(id) {
        return request({
            url: 'role/info',
            method: 'get',
            params: { id }
        })
    },
    roleUpdate(data) {
        return request({
            url: 'role/update',
            method: 'put',
            data: data
        })
    },
    deleteRole(id) {
        return request({
            url: 'role/delete',
            method: 'delete',
            data: { id }
        })
    },
    updateRoleStatus(id, status) {
        return request({
            url: 'role/updateStatus',
            method: 'put',
            data: {
                id,
                status
            }
        })
    },
    QueryRoleMenuIdList(id) {
        return request({
            url: 'role/vo/idList',
            method: 'get',
            params: { id }
        })
    },
    AssignPermissions(id, menuIds) {
        return request({
            url: 'role/assignPermissions',
            method: 'put',
            data: {
                id,
                menuIds
            }
        })
    },
    // admin用户
    queryAdminList(params) {
        return request({
            url: 'admin/list',
            method: 'get',
            params: params
        })
    },
    updateAdminStatus(id, status) {
        return request({
            url: 'admin/updateStatus',
            method: 'put',
            data: {
                id,
                status
            }
        })
    },
    addAdmin(data) {
        return request({
            url: 'admin/add',
            method: 'post',
            data: data
        })
    },
    adminInfo(id) {
        return request({
            url: 'admin/info',
            method: 'get',
            params: { id }
        })
    },
    adminUpdate(data) {
        return request({
            url: 'admin/update',
            method: 'put',
            data: data
        })
    },
    resetPassword(id, password) {
        return request({
            url: 'admin/updatePassword',
            method: 'put',
            data: {
                id,
                password
            }
        })
    },
    deleteAdmin(id) {
        return request({
            url: 'admin/delete',
            method: 'delete',
            data: { id }
        })
    },
    adminUpdatePersonal(data) {
        return request({
            url: 'admin/updatePersonal',
            method: 'put',
            data: data
        })
    },
    adminUpdatePersonalPassword(data) {
        return request({
            url: 'admin/updatePersonalPassword',
            method: 'put',
            data: data
        })
    },
    // sysLoginInfo登录日志
    querySysLoginInfoList(params) {
        return request({
            url: 'sysLoginInfo/list',
            method: 'get',
            params: params
        })
    },
    batchDeleteSysLoginInfo(ids) {
        return request({
            url: 'sysLoginInfo/batch/delete',
            method: 'delete',
            data: { ids }
        })
    },
    cleanSysLoginInfo() {
        return request({
            url: 'sysLoginInfo/clean',
            method: 'delete'
        })
    },
    deleteSysLoginInfo(id) {
        return request({
            url: 'sysLoginInfo/delete',
            method: 'delete',
            data: { id }
        })
    },
    // SysOperationLog操作日志
    querySysOperationLogList(params) {
        return request({
            url: 'sysOperationLog/list',
            method: 'get',
            params: params
        })
    },
    batchDeleteSysOperationLog(ids) {
        return request({
            url: 'sysOperationLog/batch/delete',
            method: 'delete',
            data: { ids }
        })
    },
    cleanSysOperationLog() {
        return request({
            url: 'sysOperationLog/clean',
            method: 'delete'
        })
    },
    deleteSysOperationLog(id) {
        return request({
            url: 'sysOperationLog/delete',
            method: 'delete',
            data: { id }
        })
    },
    // 文件上传
    getUploadUrl() {
        return '/api/v1/upload'
    },
    // 系统配置
    // LDAP配置接口
    saveLdapConfig(data) {
        return request({
            url: 'ldap/config',
            method: 'post',
            data: data
        })
    },
    getLdapConfig() {
        return request({
            url: 'ldap/config',
            method: 'get'
        })
    },
    testLdapConnection() {
        return request({
            url: 'ldap/test',
            method: 'post'
        })
    },
    // LDAP登录
    ldapLogin(data) {
        return request({
            url: 'ldap/login',
            method: 'post',
            data: data
        })
    },
    // 通用系统配置接口
    addConfig(data) {
        return request({
            url: 'system/config/add',
            method: 'post',
            data: data
        })
    },
    getConfigInfo(id) {
        return request({
            url: 'system/config/info',
            method: 'get',
            params: { id }
        })
    },
    getConfigByKey(configKey) {
        return request({
            url: 'system/config/key',
            method: 'get',
            params: { configKey }
        })
    },
    updateConfig(data) {
        return request({
            url: 'system/config/update',
            method: 'put',
            data: data
        })
    },
    deleteConfig(id) {
        return request({
            url: 'system/config/delete',
            method: 'delete',
            data: { id }
        })
    },
    getConfigList(params) {
        return request({
            url: 'system/config/list',
            method: 'get',
            params: params
        })
    },
    updateConfigStatus(id, status) {
        return request({
            url: 'system/config/updateStatus',
            method: 'put',
            data: {
                id,
                status
            }
        })
    },
    // 系统动态日志
    getActivityLogs(params) {
        return request({
            url: 'activity-log/latest',
            method: 'get',
            params: params
        })
    },

    // ============ 堡垒机用户组管理 ============
    // 获取用户组列表
    getUserGroupList(params) {
        return request({
            url: '/api/v1/bastion/group/list',
            method: 'get',
            params
        })
    },
    // 创建用户组
    createUserGroup(data) {
        return request({
            url: '/api/v1/bastion/group/create',
            method: 'post',
            data
        })
    },
    // 更新用户组
    updateUserGroup(data) {
        return request({
            url: '/api/v1/bastion/group/update',
            method: 'put',
            data
        })
    },
    // 删除用户组
    deleteUserGroup(data) {
        return request({
            url: '/api/v1/bastion/group/delete',
            method: 'delete',
            data
        })
    },
    // 获取用户组详情
    getUserGroupDetail(id) {
        return request({
            url: '/api/v1/bastion/group/info',
            method: 'get',
            params: { id }
        })
    },
    // 添加用户组成员
    addGroupMembers(data) {
        return request({
            url: '/api/v1/bastion/group/member/add',
            method: 'post',
            data
        })
    },
    // 移除用户组成员
    removeGroupMembers(data) {
        return request({
            url: '/api/v1/bastion/group/member/remove',
            method: 'post',
            data
        })
    },
    // 获取用户列表（用于选择成员）
    getUserList(params) {
        return request({
            url: 'admin/list',
            method: 'get',
            params
        })
    },

    // ============ 堡垒机主机授权管理 ============
    // 授权主机给用户
    authorizeHostsToUser(data) {
        return request({
            url: '/api/v1/bastion/auth/user/authorize',
            method: 'post',
            data
        })
    },
    // 授权主机给用户组
    authorizeHostsToGroup(data) {
        return request({
            url: '/api/v1/bastion/auth/group/authorize',
            method: 'post',
            data
        })
    },
    // 取消用户的主机授权
    revokeHostsFromUser(data) {
        return request({
            url: '/api/v1/bastion/auth/user/revoke',
            method: 'post',
            data
        })
    },
    // 取消用户组的主机授权
    revokeHostsFromGroup(data) {
        return request({
            url: '/api/v1/bastion/auth/group/revoke',
            method: 'post',
            data
        })
    },
    // 获取授权主机列表
    getAuthorizedHosts(params) {
        return request({
            url: '/api/v1/bastion/auth/hosts',
            method: 'get',
            params
        })
    },
    // 获取用户有权限的主机列表
    getUserAuthorizedHosts(params) {
        return request({
            url: '/api/v1/bastion/user/authorized-hosts',
            method: 'get',
            params
        })
    },
    // 检查用户是否有主机权限
    checkUserHostPermission(userId, hostId) {
        return request({
            url: '/api/v1/bastion/user/check-permission',
            method: 'get',
            params: { userId, hostId }
        })
    }
}
