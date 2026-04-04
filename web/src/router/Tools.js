import Tools from '@/views/Tools/Tools'
import Agent from '@/views/Tools/Agent'
import Knowledge from '@/views/Tools/Knowledge'
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
        path: '/ops/knowledge',
        component: Knowledge,
        meta: {sTitle: '运维工具', tTitle: '运维知识库'}
    },
]

export default routes