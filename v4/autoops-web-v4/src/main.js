import { createApp } from 'vue'
import App from './App.vue'
import router from './router/router'
import store from './store'
import api from './api'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import handleTree from '@/utils/common'
import AuthorityDirective from "@/permission/index";
import storage from '@/utils/storage'
import { setAppContext } from '@/utils/appContext'
import { openYamlDialog } from '@/views/K8s/services/yamlDialog'

// 全局样式
import './assets/css/global.css'

// 统一导入 el-icon 图标
import * as ElIconModules from '@element-plus/icons-vue'

// 创建 Vue 应用实例 ✅
const app = createApp(App)

// 统一注册 el-icon 图标组件 ✅
for (let iconName in ElIconModules) {
    app.component(iconName, ElIconModules[iconName])
}

// 挂载全局工具方法
app.config.globalProperties.$storage = storage
app.config.globalProperties.$api = api
app.config.globalProperties.$handleTree = handleTree
app.config.globalProperties.$openYamlDialog = openYamlDialog

// 全局错误处理
app.config.errorHandler = (err, _, info) => {
    console.error('全局错误处理:', err, info)
    // 如果是TOKEN相关错误，显示友好提示
    if (err.message && err.message.includes('TOKEN') || err.message && err.message.includes('401')) {
        console.log('检测到TOKEN过期错误，准备跳转登录页')
    }
}

// 处理Promise未捕获的错误
window.addEventListener('unhandledrejection', event => {
    console.error('未处理的Promise错误:', event.reason)
    // 防止错误冒泡到控制台
    event.preventDefault()
})

// 使用插件和模块
app.use(router)
app.use(store)
app.use(ElementPlus)
app.use(handleTree)
AuthorityDirective.install(app)

// 挂载根组件
setAppContext(app._context)
app.mount('#app')
