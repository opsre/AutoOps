/**
 * 业务数据的提交
 *
 *  @author xiaoRui
 */

import storage from "@/utils/storage"

export default {
    saveSysAdmin(state, sysAdmin) {
        state.sysAdmin = sysAdmin
        storage.setItem('sysAdmin', sysAdmin)
    },
    saveToken(state, token) {
        state.token = token
        storage.setItem('token', token)
    },
    saveLeftMenuList(state, leftMenuList) {
        console.log("【mutations】准备保存 leftMenuList：", leftMenuList)
        state.leftMenuList = leftMenuList
        storage.setItem('leftMenuList', leftMenuList)
    },
    savePermissionList(state, permissionList) {
        console.log("【mutations】准备保存 permissionList：", permissionList)
        state.permissionList = permissionList  // ✅ 修复
        storage.setItem('permissionList', permissionList)
    },
    saveActivePath(state, activePath) {
        state.activePath = activePath  // ✅ 修复
        storage.setItem('activePath', activePath)
    }
}
