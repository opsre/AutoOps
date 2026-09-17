import ecskey from '@/views/configcenter/ecs-key.vue'
import accountauth from '@/views/configcenter/accountauth.vue'
import KeyManage from '@/views/configcenter/KeyManage'
const routes = [
    {
        path: '/config/ecskey',
        component: ecskey,
        meta: {sTitle: '配置中心', tTitle: '主机凭据'}
    },
    {
        path: '/config/accountauth',
        component: accountauth,
        meta: {sTitle: '配置中心', tTitle: '通用凭据'}
    },
    {
        path: '/config/keymanage',
        component: KeyManage,
        meta: {sTitle: '配置中心', tTitle: '密钥管理'}
    }
]

export default routes
