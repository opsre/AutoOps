import Tools from '@/views/Tools/Tools'
import Agent from '@/views/Tools/Agent'
import AI from '@/views/Tools/AI'
import Wiki from '@/views/Tools/FeiShu_Wiki.vue'
const routes = [
    {
        path: '/ops/tools',
        component: Tools,
        meta: {sTitle: '运维工具', tTitle: '工具列表'}
    },
    {
        path: '/ops/agent',
        component: Agent,
        meta: {sTitle: '运维工具', tTitle: 'agent列表'}
    },
        {
        path: '/ops/ai',
        component: AI,
        meta: {sTitle: '运维工具', tTitle: 'AI助手'}
    },
        {
        path: '/ops/wiki',
        component: Wiki,
        meta: {sTitle: '运维工具', tTitle: '飞书文档'}
    },
]

export default routes