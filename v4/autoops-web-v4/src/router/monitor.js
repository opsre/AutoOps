import Domain from '@/views/monitor/monitor_ssl.vue'
import Incident from '@/views/monitor/System_failure.vue'
import Alert from '@/views/monitor/Alert_Rule.vue'
import AlertPolicyRule from '@/views/monitor/alert/AlertPolicyRule.vue'
import API from '@/views/monitor/Api_Alert.vue'
import BusinessFunctionCatalog from '@/views/monitor/BusinessFunctionCatalog.vue'
const routes = [
    {
        path: '/monitor/domain',
        component: Domain,
        meta: {sTitle: '监控中心', tTitle: '域名监控'}
    },
    {
        path: '/monitor/incident',
        component: Incident,
        meta: {sTitle: '监控中心', tTitle: '故障管理'}
    },
    {
        path: '/monitor/alert',
        component: Alert,
        meta: {sTitle: '监控中心', tTitle: '告警规则'}
    },
    {
        path: '/monitor/alert-rule',
        component: AlertPolicyRule,
        meta: {sTitle: '监控中心', tTitle: '配置告警规则'}
    },
    {
        path: '/monitor/api',
        component: API,
        meta: {sTitle: '监控中心', tTitle: 'API监控'}
    },
    {
        path: '/monitor/business-function-catalog',
        component: BusinessFunctionCatalog,
        meta: {sTitle: '监控中心', tTitle: 'API浏览'}
    }
]

export default routes