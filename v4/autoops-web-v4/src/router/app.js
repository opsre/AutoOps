import App from '@/views/app/application.vue'
import QuickRelease from '@/views/app/app_quick_release.vue'
import QuickTemp from '@/views/app/app_quick_temp.vue'
import Upgrade   from '@/views/app/service_release.vue'
import ServiceReleaseDetail from '@/views/app/service_release_detail.vue'
import ShRelease from '@/views/app/sh_release.vue'
import ShReleaseDetail from '@/views/app/sh_release_detail.vue'
import AppView from '@/views/app/application-view.vue'
const routes = [
    {
        path: '/app/application',
        component: App,
        meta: {sTitle: '服务管理', tTitle: '应用列表'}
    },
    {
        path: '/app/view',
        component: AppView,
        meta: {sTitle: '服务管理', tTitle: '应用视图'}
    },
    {
        path: '/app/quick-release',
        component: QuickRelease,
        meta: {sTitle: '服务管理', tTitle: '快速发布'}
    },
    {
        path: '/app/quick-temp/:id',
        component: QuickTemp,
        meta: {sTitle: '服务管理', tTitle: '发布模板'}
    },
    {
        path: '/app/service_work',
        component: Upgrade,
        meta: {sTitle: '服务管理', tTitle: '服务上线'}
    },
    {
        path: '/app/service-release/:id',
        component: ServiceReleaseDetail,
        meta: {sTitle: '服务管理', tTitle: '工单详情'}
    },
    {
        path: '/app/sh-release',
        component: ShRelease,
        meta: {sTitle: '服务管理', tTitle: '脚本上线'}
    },
    {
        path: '/app/sh-release/:id',
        component: ShReleaseDetail,
        meta: {sTitle: '服务管理', tTitle: '脚本上线详情'}
    },
]

export default routes
