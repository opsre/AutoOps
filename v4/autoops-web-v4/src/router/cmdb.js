import Host from '@/views/cmdb/cmdbHost.vue'
import Switch from '@/views/cmdb/Cmdb_Switch.vue'
import Db from '@/views/cmdb/cmdbDB.vue'
import SSH from '@/views/cmdb/Host/SSH.vue'
import HostDetail from '@/views/cmdb/Host/HostDetail.vue'
import NetworkSSH from '@/views/cmdb/Network/NetworkSSH.vue'
import NetworkRDP from '@/views/cmdb/Network/NetworkRDP.vue'
import NetworkMonitor from '@/views/cmdb/Network/NetworkMonitor.vue'
import WindowsMonitor from '@/views/cmdb/Network/WindowsMonitor.vue'
import DBTerminal from '@/views/cmdb/DB/cmdbDBTerminal.vue'
import RedisTerminal from '@/views/cmdb/DB/cmdbRedisTerminal.vue'
import ESTerminal from '@/views/cmdb/DB/cmdbESTerminal.vue'
import MongoTerminal from '@/views/cmdb/DB/cmdbMongoTerminal.vue'
const routes = [
    {
        path: '/cmdb/ecs',
        component: Host,
        meta: {sTitle: '资产管理', tTitle: '主机管理'}
    },
    {
        path: '/cmdb/switch',
        component: Switch,
        meta: {sTitle: '资产管理', tTitle: '网络设备'}
    },
    {
        path: '/cmdb/switch/:id',
        component: NetworkMonitor,
        meta: {sTitle: '资产管理', tTitle: '交换机监控'}
    },
    {
        path: '/cmdb/windows/:id',
        component: WindowsMonitor,
        meta: {sTitle: '资产管理', tTitle: 'Windows监控'}
    },
    {
        path: '/cmdb/db',
        component: Db,
        meta: {sTitle: '资产管理', tTitle: '数据管理'}
    }, 
    {
        path: '/cmdb/ssh',
        component: SSH,
        meta: {sTitle: '资产管理', tTitle: '终端登录'}
    },
    {
        path: '/cmdb/host/detail/:id',
        component: HostDetail,
        meta: {sTitle: '资产管理', tTitle: '主机详情'}
    },
    {
        path: '/cmdb/dbterminal',
        component: DBTerminal,
        meta: {sTitle: '数据管理', tTitle: '数据库终端'}
    },
    {
        path: '/cmdb/redis-terminal',
        component: RedisTerminal,
        meta: {sTitle: '数据管理', tTitle: 'Redis终端'}
    },
    {
        path: '/cmdb/es-terminal',
        component: ESTerminal,
        meta: {sTitle: '数据管理', tTitle: 'ES终端'}
    },
    {
        path: '/cmdb/mongo-terminal',
        component: MongoTerminal,
        meta: {sTitle: '数据管理', tTitle: 'MongoDB终端'}
    },
    {
        path: '/cmdb/network/ssh-terminal',
        component: NetworkSSH,
        meta: {sTitle: '资产管理', tTitle: '网络设备终端'}
    },
    {
        path: '/cmdb/network/rdp-terminal',
        component: NetworkRDP,
        meta: {sTitle: '资产管理', tTitle: 'RDP远程桌面'}
    }

]

export default routes
