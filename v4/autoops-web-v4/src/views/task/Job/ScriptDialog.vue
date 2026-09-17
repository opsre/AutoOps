<template>
  <div v-if="visible" class="custom-dialog">
    <div class="dialog-mask"></div>
    <div class="dialog-container">
      <div class="dialog-header">
        <h3>{{ title }}</h3>
        <button @click="handleClose">×</button>
      </div>
      <div class="dialog-body" v-html="content"></div>
      <div class="dialog-footer">
        <button @click="handleClose">关闭</button>
      </div>
    </div>
  </div>
</template>

<script>
import { highlight } from '@/utils/highlight'

export default {
  name: 'ScriptDialog',
  props: {
    title: {
      type: String,
      default: '脚本内容'
    }
  },
  data() {
    return {
      visible: false,
      content: ''
    }
  },
  methods: {
    show(content) {
      let processedContent = content
      if (typeof processedContent !== 'string') {
        processedContent = processedContent.content || JSON.stringify(processedContent, null, 2)
      }
      processedContent = processedContent.replace(/\\n/g, '\n').replace(/\\r/g, '\r')
      const highlightedContent = highlight(processedContent, 'bash')
      this.content = `
        <div class="highlight-container">
          <pre style="white-space: pre-wrap; background: #121313e9; padding: 16px; border-radius: 6px;">
            ${highlightedContent}
          </pre>
        </div>
      `
      this.visible = true
    },
    handleClose() {
      this.visible = false
      this.$emit('close')
    }
  }
}
</script>

<style scoped>
.custom-dialog {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 9999;
}

.dialog-mask {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
}

.dialog-container {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 80%;
  max-width: 800px;
  background: white;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(29, 23, 23, 0.94);
}

.dialog-header {
  padding: 20px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dialog-body {
  padding: 20px;
  max-height: 70vh;
  overflow: auto;
}

.dialog-footer {
  padding: 10px 20px;
  border-top: 1px solid #eee;
  text-align: right;
}
</style>
