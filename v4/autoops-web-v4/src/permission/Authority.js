/**
 * 权限指令 - Vue 3 版本
 * 修改为：按钮可见但禁用 + 点击提示权限不足
 *
 * @author xiaoRui
 */
import { checkAuthority } from "@/utils/authority"
import storage from "@/utils/storage"
import { ElMessage } from 'element-plus'

export default {
    mounted(el, binding) {
        const { value } = binding
        const permissions = storage.getItem('permissionList') || []
        const hasPermission = checkAuthority(value, permissions)

        if (!hasPermission) {
            // 如果没有权限，禁用按钮而不是移除
            el.disabled = true

            // 添加视觉反馈样式
            el.style.opacity = '0.5'
            el.style.cursor = 'not-allowed'
            el.style.pointerEvents = 'none'

            // 给按钮添加禁用状态的类名
            el.classList.add('no-permission')

            // 创建一个包装容器来捕获点击事件
            const wrapper = document.createElement('div')
            wrapper.style.position = 'relative'
            wrapper.style.display = 'inline-block'

            // 将原按钮包装起来
            if (el.parentNode) {
                el.parentNode.insertBefore(wrapper, el)
                wrapper.appendChild(el)
            }

            // 在包装容器上添加点击事件监听
            wrapper.addEventListener('click', (e) => {
                e.preventDefault()
                e.stopPropagation()
                ElMessage.warning('您没有该操作权限，请联系管理员')
            })

            // 恢复按钮的 pointer-events 以便包装容器能捕获事件
            el.style.pointerEvents = 'auto'

        } else {
            // 有权限时添加 code 属性用于调试
            el.setAttribute('code', value)
        }
    }
}
