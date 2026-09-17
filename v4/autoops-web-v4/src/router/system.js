import Personal from '@/views/system/Personal.vue'
import Admin from '@/views/system/Admin.vue'
import Role from '@/views/system/Role.vue'
import Dept from '@/views/system/user/Dept.vue'
import Post from '@/views/system/HostGroup.vue'
import Menu from '@/views/system/Menu.vue'
import Config from '@/views/system/sysconfig.vue'
import LoginLog from '@/views/audit/LoginLog.vue'
import Operator from '@/views/audit/Operator.vue'
import DBLog from '@/views/audit/DBLog.vue'
import Recording from '@/views/audit/recording.vue'
import HostAuthorization from '@/views/system/HostAuthorization.vue'
const routes = [
    {
        path: '/system/personal',
        component: Personal,
        meta: {sTitle: '个人中心', tTitle: '个人信息'}
    },
    {
        path: '/system/admin',
        component: Admin,
        meta: {sTitle: '基础管理', tTitle: '用户信息'}
    },
    {
        path: '/system/role',
        component: Role,
        meta: {sTitle: '基础管理', tTitle: '角色信息'}
    },
    {
        path: '/system/menu',
        component: Menu,
        meta: {sTitle: '基础管理', tTitle: '菜单信息'}
    },
    {
        path: '/system/dept',
        component: Dept,
        meta: {sTitle: '基础管理', tTitle: '部门信息'}
    },
    {
        path: '/system/post',
        component: Post,
        meta: {sTitle: '基础管理', tTitle: '用户组管理'}
    },
    {
        path: '/bastion/authorization',
        name: 'BastionAuthorization',
        component: HostAuthorization,
        meta: {sTitle: '基础管理', tTitle: '主机授权'}
    },
    {
        path: '/system/config',
        component: Config,
        meta: {sTitle: '基础管理', tTitle: '系统配置'}
    },
    {
        path: '/monitor/loginlog',
        component: LoginLog,
        meta: {sTitle: '日志管理', tTitle: '登录日志'}
    },
    {
        path: '/monitor/operator',
        component: Operator,
        meta: {sTitle: '日志管理', tTitle: '操作日志'}
    },
    {
        path: '/monitor/dblog',
        component: DBLog,
        meta: {sTitle: '日志管理', tTitle: '数据日志'}
    },
    {
        path: '/monitor/recording',
        component: Recording,
        meta: {sTitle: '日志管理', tTitle: '会话录制'}
    }
]

export default routes
