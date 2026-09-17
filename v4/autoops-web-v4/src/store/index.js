/**
 * vuex状态管理（Vue 3 版本）
 *
 * @author xiaoRui
 */

import { createStore } from 'vuex'
import mutations from './mutations'
import storage from '@/utils/storage'

// 创建并导出 Vuex Store
export default createStore({
    state: {
        sysAdmin: storage.getItem("sysAdmin") || "",
        token: storage.getItem("token") || "",
        leftMenuList: storage.getItem("leftMenuList") || "",
        permissionList: storage.getItem("permissionList") || "",
        activePath: storage.getItem("activePath") || ""
    },
    mutations
})
