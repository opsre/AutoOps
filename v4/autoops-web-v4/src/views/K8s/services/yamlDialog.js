import { defineComponent, h, reactive, render } from 'vue'
import ResourceYamlDialog from '../components/ResourceYamlDialog.vue'
import { getAppContext } from '@/utils/appContext'

export function openYamlDialog(options) {
  const context = getAppContext()

  if (!context) {
    throw new Error('YAML Dialog service 未初始化')
  }

  const container = document.createElement('div')
  document.body.appendChild(container)

  const state = reactive({
    visible: true
  })

  const cleanup = () => {
    render(null, container)
    if (container.parentNode) {
      container.parentNode.removeChild(container)
    }
  }

  const close = () => {
    state.visible = false
  }

  const DialogHost = defineComponent({
    name: 'YamlDialogHost',
    setup() {
      return () => h(ResourceYamlDialog, {
        ...options,
        visible: state.visible,
        'onUpdate:visible': (value) => {
          state.visible = value
          if (!value) {
            cleanup()
          }
        },
        onSaved: async () => {
          if (typeof options.onSaved === 'function') {
            await options.onSaved()
          }
        }
      })
    }
  })

  const vnode = h(DialogHost)
  vnode.appContext = context
  render(vnode, container)

  return {
    close
  }
}
