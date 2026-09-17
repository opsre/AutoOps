import Template from '@/views/task/TaskTemplate.vue'
import Job from '@/views/task/TaskJob.vue'
import Ansible from '@/views/task/TaskAnsible.vue'

const routes = [
    {
        path: '/task/template',
        component: Template,
        meta: {sTitle: '任务中心', tTitle: '任务模版'}
    },
    {
        path: '/task/job',
        component: Job,
        meta: {sTitle: '任务中心', tTitle: '任务作业'}
    },
    {
        path: '/task/ansible',
        component: Ansible,
        meta: {sTitle: '任务中心', tTitle: 'Ansible任务'}
    }
]

export default routes
