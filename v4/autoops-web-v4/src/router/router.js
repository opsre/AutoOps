import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import Home from '@/views/Home.vue'
import Dashboard from '@/views/dashboard/index.vue'
import systemRoutes from './system'
import cmdbRoutes from './cmdb'
import k8sRoutes from './k8s'
import configRoutes from './config'
import task from './task'
import tools from './Tools'
import app from './app'
import monitorRoutes from './monitor'

// 路由集合
const routes = [
    {path: '/',      redirect: '/login'},
    {path: '/login', component: Login},
    {
        path: '/home',
        component: Home,
        redirect: '/dashboard',
        children: [
            {
                path: '/dashboard',
                component: Dashboard,
                meta: {tTitle: '仪表盘'}
            },
            // 以下路由已注释，如需使用请先创建对应组件并导入
            // {
            //     path: '/dashboard/assets',
            //     component: AssetsManagement,
            //     meta: {tTitle: '资产管理'}
            // },
            // {
            //     path: '/dashboard/services',
            //     component: ServiceManagement,
            //     meta: {tTitle: '服务管理'}
            // },
            // {
            //     path: '/dashboard/monitoring',
            //     component: MonitoringAlerts,
            //     meta: {tTitle: '监控告警'}
            // },
            // {
            //     path: '/dashboard/users',
            //     component: UserStatistics,
            //     meta: {tTitle: '用户统计'}
            // },
            ...systemRoutes,
            ...cmdbRoutes,
            ...k8sRoutes,
            ...configRoutes,
            ...monitorRoutes,
            ...task,
            ...tools,
            ...app,
        ]
    }
]

// 创建 router 实例 ✅ 唯一的一次
const router = createRouter({
    history: createWebHistory(),
    routes // 使用上面定义的 routes 数组
})

// 引入 storage 工具类
import storage from '@/utils/storage'

// 挂载导航守卫
router.beforeEach((to, from, next) => {
    const token = storage.getItem('token')
    const sysAdmin = storage.getItem('sysAdmin')

    // 如果已经登录，且访问的是登录页，则跳转到首页
    if (to.path === '/login' && token && sysAdmin) {
        return next('/dashboard')
    }

    // 如果访问的是非登录页，检查是否登录
    if (to.path !== '/login') {
        if (!token || !sysAdmin) {
            // 未登录，跳转到登录页
            return next('/login')
        }
    }

    // 其他情况正常放行
    next()
})

export default router
