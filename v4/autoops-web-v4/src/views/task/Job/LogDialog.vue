<template>
  <div v-if="visible" class="custom-dialog log-dialog">
    <div class="dialog-mask"></div>
    <div class="dialog-container">
      <div class="dialog-header">
        <h3 style="display: block;">{{ title }}</h3>
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
  name: 'LogDialog',
  props: {
    title: {
      type: String,
      default: '任务日志'
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
        processedContent = processedContent.logs || JSON.stringify(processedContent, null, 2)
      }
      processedContent = processedContent.replace(/\\n/g, '\n').replace(/\\r/g, '\r')
      const highlightedContent = highlight(processedContent, 'bash')
      this.content = `
        <div class="highlight-container">
          <pre style="white-space: pre-wrap; background: #000; color: #fff; padding: 16px; border-radius: 6px; font-family: 'Courier New', monospace;">
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
.log-dialog .dialog-container {
  background: #fff;
  color: #fff;
  width: 80%;
  max-width: 1200px;
  min-width: 800px;
}

.log-dialog .dialog-header {
  border-bottom: 1px solid #333;
}

.log-dialog .dialog-header h3 {
  color: #000;
  margin: 0;
  font-size: 18px;
}

.log-dialog .dialog-header button {
  color: #000;
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
}

.log-dialog .dialog-footer {
  border-top: 1px solid #333;
}

.log-dialog .dialog-footer button {
  background-color: #0e639c;
  color: #fff;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}

.log-dialog .dialog-footer button:hover {
  background-color: #1177bb;
}

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
  background: white;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
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
