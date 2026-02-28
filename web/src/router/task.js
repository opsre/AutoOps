import Template from '@/views/task/TaskTemplate.vue'
import Job from '@/views/task/TaskJob.vue'
import Ansible from '@/views/task/TaskAnsible.vue'
import AnsibleHistory from '@/views/task/AnsibleTaskHistory.vue'
import AnsibleConfig from '@/views/task/TaskConfig.vue'

const routes = [
    {
        path: '/task/template',
        component: Template,
        meta: { sTitle: '任务中心', tTitle: '任务模版' }
    },
    {
        path: '/task/job',
        component: Job,
        meta: { sTitle: '任务中心', tTitle: '任务作业' }
    },
    {
        path: '/task/ansible',
        component: Ansible,
        meta: { sTitle: '任务中心', tTitle: 'Ansible任务' }
    },
    {
        path: '/task/config',
        component: AnsibleConfig,
        meta: { sTitle: '任务中心', tTitle: '配置管理' }
    },
    {
        path: '/task/ansible/history',
        name: 'AnsibleTaskHistory',
        component: AnsibleHistory,
        meta: { sTitle: '任务中心', tTitle: '执行历史', hidden: true }
    }
]

export default routes
