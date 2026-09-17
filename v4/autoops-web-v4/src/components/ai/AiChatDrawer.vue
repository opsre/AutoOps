<template>
  <!-- AI助手抽屉 -->
  <el-drawer
    v-model="visible"
    title=""
    direction="rtl"
    size="500px"
    :close-on-click-modal="false"
    class="ai-drawer"
    @close="handleClose"
  >
    <!-- 自定义头部 -->
    <template #header>
      <div class="ai-drawer-header">
        <div class="ai-header-content">
          <div class="ai-title-group">
            <h3 class="ai-title">AI智能助手</h3>
            <span class="ai-subtitle">Powered by Advanced AI</span>
          </div>
        </div>
        <div class="ai-status-indicator">
          <span class="status-dot"></span>
          <span class="status-text">在线</span>
        </div>
      </div>
    </template>
    <div class="ai-chat-container">
      <!-- 对话消息列表 -->
      <div class="chat-messages" ref="chatMessages">
        <div
          v-for="(message, index) in chatMessages"
          :key="index"
          :class="['message-item', message.role === 'user' ? 'user-message' : 'assistant-message']"
        >
          <div class="message-avatar">
            <el-icon v-if="message.role === 'user'" size="24"><User /></el-icon>
            <img v-else src="../../assets/image/智能ai.svg" alt="AI" class="avatar-ai-icon" />
          </div>
          <div class="message-content">
            <!-- 用户消息：简单文本显示 -->
            <div v-if="message.role === 'user'" class="message-text">{{ message.content }}</div>

            <!-- AI消息：支持结构化内容 -->
            <div v-else class="message-text">
              <!-- 如果有结构化内容，渲染结构化内容 -->
              <template v-if="message.contentBlocks && message.contentBlocks.length > 0">
                <div
                  v-for="(block, blockIndex) in message.contentBlocks"
                  :key="blockIndex"
                  :class="['content-block', `block-${block.type}`]"
                >
                  <!-- Markdown表格（自动识别） -->
                  <div v-if="block.type === 'text' && isMarkdownTable(block.content)" class="markdown-table-block">
                    <component
                      :is="`h4`"
                      class="heading-block"
                      v-if="extractTableTitle(block.content)"
                    >
                      {{ extractTableTitle(block.content) }}
                    </component>
                    <MarkdownTable :source="block" :maxHeight="360" />
                  </div>

                  <!-- 代码块 -->
                  <div v-else-if="block.type === 'code'" class="code-block">
                    <div class="code-header">
                      <span class="code-language">{{ block.language || 'text' }}</span>
                      <el-button
                        type="primary"
                        size="small"
                        link
                        @click="copyCode(block.content, blockIndex)"
                      >
                        <el-icon><DocumentCopy /></el-icon>
                        {{ copiedIndex === blockIndex ? '已复制' : '复制代码' }}
                      </el-button>
                    </div>
                    <pre><code>{{ block.content }}</code></pre>
                  </div>

                  <!-- 标题 -->
                  <component
                    v-else-if="block.type === 'heading'"
                    :is="`h${block.level || 3}`"
                    class="heading-block"
                  >
                    {{ block.content }}
                  </component>

                  <!-- 列表 -->
                  <ul v-else-if="block.type === 'list'" class="list-block">
                    <li v-for="(item, itemIndex) in parseListItems(block.content)" :key="itemIndex">
                      {{ item }}
                    </li>
                  </ul>

                  <!-- 普通文本（兜底） -->
                  <div v-else-if="block.type === 'text'" class="text-block markdown-text">{{ block.content }}</div>

                  <!-- 其他类型：降级为文本显示 -->
                  <div class="text-block markdown-text">{{ block.content }}</div>
                </div>
              </template>

              <!-- 降级方案：如果没有结构化内容，也尝试识别Markdown表格 -->
              <template v-else>
                <div v-if="isMarkdownTable(message.content)" class="markdown-table-block">
                  <component
                    :is="`h4`"
                    class="heading-block"
                    v-if="extractTableTitle(message.content)"
                  >
                    {{ extractTableTitle(message.content) }}
                  </component>
                  <MarkdownTable :source="message.content" :maxHeight="360" />
                </div>
                <div v-else class="text-block markdown-text">{{ message.content }}</div>
              </template>
            </div>
          </div>
        </div>
        <div v-if="isLoading" class="message-item assistant-message">
          <div class="message-avatar">
            <img src="../../assets/image/智能ai.svg" alt="AI" class="avatar-ai-icon" />
          </div>
          <div class="message-content">
            <div class="message-text typing-indicator">
              <span></span><span></span><span></span>
            </div>
          </div>
        </div>
      </div>

      <!-- 输入框 -->
      <div class="chat-input-container">
        <div class="input-wrapper">
          <el-input
            v-model="userInput"
            type="textarea"
            :rows="3"
            placeholder="请输入您的问题... (Enter发送 / Shift+Enter换行)"
            @keydown="handleKeydown"
            :disabled="!selectedModelId"
            class="chat-textarea"
          />
          <div class="input-divider"></div>
          <div class="input-footer">
            <el-select
              v-model="selectedModelId"
              placeholder="选择模型"
              size="small"
              class="model-select"
              @change="handleModelChange"
            >
              <el-option
                v-for="model in enabledModels"
                :key="model.id"
                :label="model.name"
                :value="model.id"
              />
            </el-select>
            <el-button
              type="primary"
              :loading="isLoading"
              :disabled="!userInput.trim() || !selectedModelId"
              @click="sendMessage"
              class="send-button"
              circle
              icon="ArrowUp"
            >
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </el-drawer>
</template>

<script>
import { User, DocumentCopy, ArrowUp } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import aiApi from '@/api/ai'
import MarkdownTable from '../MarkdownTable.vue'

const AI_DRAWER_SESSION_STORAGE_KEY = 'ops_ai_drawer_session_id'

export default {
  name: 'AiChatDrawer',
  components: { User, DocumentCopy, ArrowUp, MarkdownTable },
  props: {
    modelValue: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:modelValue'],
  data() {
    return {
      selectedModelId: null,
      enabledModels: [],
      chatMessages: [],
      userInput: '',
      isLoading: false,
      copiedIndex: null,
      aiSessionId: null
    }
  },
  computed: {
    visible: {
      get() {
        return this.modelValue
      },
      set(val) {
        this.$emit('update:modelValue', val)
      }
    }
  },
  watch: {
    modelValue(val) {
      if (val && this.enabledModels.length === 0) {
        this.loadEnabledModels()
      }
    }
  },
  methods: {
    resetConversation() {
      this.chatMessages = []
      this.userInput = ''
      this.isLoading = false
      this.copiedIndex = null
      this.aiSessionId = null
      sessionStorage.removeItem(AI_DRAWER_SESSION_STORAGE_KEY)
    },

    // 简单判断文本是否为紧凑型Markdown表格（含表头与分隔线）
    isMarkdownTable(text) {
      if (!text || typeof text !== 'string') return false
      return /\|\s*[^\|]+\s*\|/.test(text) && /-\s*-\s*-/.test(text)
    },

    // 提取表格标题：位于第一行，且不以"|"开头
    extractTableTitle(text) {
      if (!text || typeof text !== 'string') return ''
      const firstLine = text.split(/\r?\n/)[0]?.trim()
      if (firstLine && !firstLine.startsWith('|')) return firstLine
      return ''
    },

    // 加载启用的AI模型
    async loadEnabledModels() {
      try {
        const res = await aiApi.getEnabledAiModels()
        if (res.data && (res.data.code === 0 || res.data.code === 200)) {
          this.enabledModels = res.data.data || []

          if (this.enabledModels.length > 0 && !this.selectedModelId) {
            this.selectedModelId = this.enabledModels[0].id
          }

          if (this.enabledModels.length === 0) {
            ElMessage.warning('暂无可用的AI模型，请先在系统配置中添加并启用模型')
          }
        } else {
          ElMessage.error(res.data?.message || res.data?.msg || '获取AI模型失败')
        }
      } catch (error) {
        console.error('获取AI模型失败:', error)
        ElMessage.error('获取AI模型失败')
      }
    },

    // 切换模型
    handleModelChange() {
      console.log('切换模型:', this.selectedModelId)
    },

    // 处理键盘事件
    handleKeydown(event) {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault()
        this.sendMessage()
      }
    },

    // 发送消息
    async sendMessage() {
      if (!this.userInput.trim()) {
        return
      }

      if (!this.selectedModelId) {
        ElMessage.warning('请先选择AI模型')
        return
      }

      const userMessage = {
        role: 'user',
        content: this.userInput.trim()
      }

      this.chatMessages.push(userMessage)
      const input = this.userInput
      this.userInput = ''
      this.isLoading = true

      try {
        const payload = {
          modelId: this.selectedModelId,
          messages: [userMessage]
        }
        if (this.aiSessionId) payload.sessionId = this.aiSessionId

        const res = await aiApi.aiChat(payload)

        if (res.data && (res.data.code === 0 || res.data.code === 200)) {
          const responseData = res.data.data || {}

          const assistantMessage = {
            role: 'assistant',
            content: responseData.text || '抱歉，我没有理解您的问题',
            contentBlocks: Array.isArray(responseData.content) ? responseData.content : null
          }

          this.chatMessages.push(assistantMessage)

          if (responseData.sessionId) {
            this.aiSessionId = responseData.sessionId
            sessionStorage.setItem(AI_DRAWER_SESSION_STORAGE_KEY, responseData.sessionId)
          }

          this.$nextTick(() => {
            this.scrollToBottom()
          })
        } else {
          ElMessage.error(res.data?.message || res.data?.msg || 'AI对话失败')
          this.userInput = input
        }
      } catch (error) {
        console.error('AI对话失败:', error)
        ElMessage.error('AI对话失败: ' + (error.message || '未知错误'))
        this.userInput = input
      } finally {
        this.isLoading = false
      }
    },

    // 复制代码到剪贴板
    async copyCode(code, index) {
      try {
        await navigator.clipboard.writeText(code)
        this.copiedIndex = index
        ElMessage.success('代码已复制到剪贴板')
        setTimeout(() => {
          this.copiedIndex = null
        }, 2000)
      } catch (error) {
        console.error('复制失败:', error)
        this.copyCodeFallback(code, index)
      }
    },

    // 复制代码降级方案
    copyCodeFallback(code, index) {
      try {
        const textarea = document.createElement('textarea')
        textarea.value = code
        textarea.style.position = 'fixed'
        textarea.style.opacity = '0'
        document.body.appendChild(textarea)
        textarea.select()

        const successful = document.execCommand('copy')
        document.body.removeChild(textarea)

        if (successful) {
          this.copiedIndex = index
          ElMessage.success('代码已复制到剪贴板')
          setTimeout(() => {
            this.copiedIndex = null
          }, 2000)
        } else {
          ElMessage.error('复制失败，请手动复制')
        }
      } catch (error) {
        console.error('降级复制也失败:', error)
        ElMessage.error('复制失败，请手动复制')
      }
    },

    // 解析列表项
    parseListItems(content) {
      if (!content) return []
      return content
        .split('\n')
        .map(item => item.trim())
        .filter(item => item.length > 0)
        .map(item => {
          return item.replace(/^[-*]\s*/, '').replace(/^\d+\.\s*/, '')
        })
    },

    // 滚动到底部
    scrollToBottom() {
      if (this.$refs.chatMessages) {
        this.$refs.chatMessages.scrollTop = this.$refs.chatMessages.scrollHeight
      }
    },

    // 关闭抽屉
    handleClose() {
      this.$emit('update:modelValue', false)
    }
  },
  mounted() {
    const savedSessionId = sessionStorage.getItem(AI_DRAWER_SESSION_STORAGE_KEY)
    if (savedSessionId) {
      this.aiSessionId = savedSessionId
    }
  }
}
</script>

<style lang="less">
// AI聊天容器样式
.ai-chat-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 180px);
  overflow: hidden;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px 16px 0 16px;

  // 自定义滚动条样式
  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background: rgba(102, 126, 234, 0.05);
    border-radius: 3px;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(102, 126, 234, 0.3);
    border-radius: 3px;
    transition: all 0.3s ease;

    &:hover {
      background: rgba(102, 126, 234, 0.5);
    }
  }

  .message-item {
    display: flex;
    gap: 10px;
    margin-bottom: 16px;
    animation: messageSlideIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;

    .message-avatar {
      flex-shrink: 0;
      width: 38px;
      height: 38px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
      position: relative;

      .avatar-ai-icon {
        width: 32px;
        height: 32px;
        object-fit: contain;
      }

      &::before {
        content: '';
        position: absolute;
        inset: -2px;
        border-radius: 50%;
        padding: 2px;
        background: linear-gradient(135deg, #667eea, #764ba2);
        -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
        -webkit-mask-composite: xor;
        mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
        mask-composite: exclude;
        opacity: 0.5;
      }
    }

    .message-content {
      flex: 1;
      min-width: 0;

      .message-text {
        background: linear-gradient(135deg, rgba(26, 31, 46, 0.8) 0%, rgba(26, 31, 46, 0.6) 100%);
        border: 1px solid rgba(102, 126, 234, 0.2);
        padding: 10px 14px;
        border-radius: 12px;
        line-height: 1.5;
        word-wrap: break-word;
        font-size: 13px;
        color: rgba(255, 255, 255, 0.9);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
        position: relative;
        backdrop-filter: blur(10px);

        // 内容块通用样式
        .content-block {
          margin-bottom: 12px;

          &:last-child {
            margin-bottom: 0;
          }
        }

        // 文本块样式
        .text-block {
          margin: 0 0 8px 0;
          line-height: 1.6;
          font-size: 13px;
          color: rgba(255, 255, 255, 0.9);

          &:last-child {
            margin-bottom: 0;
          }
        }

        .markdown-text {
          white-space: pre-wrap;
        }

        // 代码块样式
        .code-block {
          background: linear-gradient(135deg, #1a1d2e 0%, #16192a 100%);
          border: 1px solid rgba(102, 126, 234, 0.3);
          border-radius: 12px;
          overflow: hidden;
          margin: 12px 0;
          box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);

          .code-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px 14px;
            background: rgba(102, 126, 234, 0.1);
            border-bottom: 1px solid rgba(102, 126, 234, 0.2);

            .code-language {
              font-size: 11px;
              color: #a78bfa;
              font-weight: 600;
              text-transform: uppercase;
              letter-spacing: 0.5px;
            }

            .el-button {
              font-size: 11px;
              padding: 6px 10px;
              height: auto;
              background: rgba(102, 126, 234, 0.2);
              border: 1px solid rgba(102, 126, 234, 0.3);
              transition: all 0.3s ease;

              &:hover {
                background: rgba(102, 126, 234, 0.3);
                transform: translateY(-1px);
                box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
              }

              .el-icon {
                margin-right: 4px;
              }
            }
          }

          pre {
            margin: 0;
            padding: 14px;
            background: transparent;
            overflow-x: auto;

            &::-webkit-scrollbar {
              height: 6px;
            }

            &::-webkit-scrollbar-track {
              background: rgba(102, 126, 234, 0.05);
            }

            &::-webkit-scrollbar-thumb {
              background: rgba(102, 126, 234, 0.3);
              border-radius: 3px;
            }

            code {
              font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', 'Monaco', 'Courier New', monospace;
              font-size: 12px;
              line-height: 1.6;
              color: #e2e8f0;
              display: block;
              white-space: pre;
            }
          }
        }

        // 标题样式
        .heading-block {
          margin: 12px 0 8px 0;
          font-weight: 600;
          background: linear-gradient(135deg, #667eea 0%, #a78bfa 100%);
          background-clip: text;
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          line-height: 1.4;

          &:first-child {
            margin-top: 0;
          }
        }

        h1.heading-block {
          font-size: 18px;
        }

        h2.heading-block {
          font-size: 16px;
        }

        h3.heading-block {
          font-size: 14px;
        }

        h4.heading-block,
        h5.heading-block,
        h6.heading-block {
          font-size: 13px;
        }

        // 列表样式
        .list-block {
          margin: 8px 0;
          padding-left: 20px;
          font-size: 13px;

          li {
            margin-bottom: 6px;
            line-height: 1.6;
            color: rgba(255, 255, 255, 0.85);
            position: relative;

            &::marker {
              color: #a78bfa;
            }

            &:last-child {
              margin-bottom: 0;
            }
          }
        }
      }
    }

    &.user-message {
      flex-direction: row-reverse;

      .message-avatar {
        background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
        box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);

        &::before {
          background: linear-gradient(135deg, #3b82f6, #2563eb);
        }
      }

      .message-content .message-text {
        background: linear-gradient(135deg, rgba(191, 219, 254, 0.95) 0%, rgba(219, 234, 254, 0.95) 100%);
        border: 1px solid rgba(59, 130, 246, 0.3);
        color: #1e293b;
        text-align: right;
        box-shadow: 0 2px 8px rgba(59, 130, 246, 0.15);
      }
    }

    &.assistant-message {
      .message-avatar {
        background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
        box-shadow: 0 4px 12px rgba(251, 191, 36, 0.4);

        &::before {
          background: linear-gradient(135deg, #fbbf24, #f59e0b);
        }
      }
    }
  }
}

@keyframes messageSlideIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.typing-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 12px 0;

  span {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: linear-gradient(135deg, #667eea 0%, #a78bfa 100%);
    box-shadow: 0 0 8px rgba(102, 126, 234, 0.5);
    animation: typingBounce 1.4s infinite ease-in-out;

    &:nth-child(2) {
      animation-delay: 0.2s;
    }

    &:nth-child(3) {
      animation-delay: 0.4s;
    }
  }
}

@keyframes typingBounce {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.7;
  }
  30% {
    transform: translateY(-12px);
    opacity: 1;
  }
}

.chat-input-container {
  padding: 16px;
  position: relative;
  flex-shrink: 0;

  .input-wrapper {
    background: transparent;
    border: 1px solid rgba(102, 126, 234, 0.25);
    border-radius: 16px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
    overflow: visible;
    transition: all 0.3s ease;
    display: flex;
    flex-direction: column;
    backdrop-filter: blur(20px);

    &:hover {
      border-color: rgba(102, 126, 234, 0.4);
      box-shadow: 0 6px 28px rgba(102, 126, 234, 0.15);
    }

    &:focus-within {
      border-color: #667eea;
      box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.2), 0 8px 32px rgba(102, 126, 234, 0.25);
    }
  }

  .chat-textarea {
    :deep(.el-textarea__inner) {
      background: linear-gradient(135deg, rgba(26, 31, 46, 0.6) 0%, rgba(30, 35, 50, 0.6) 100%);
      border: none;
      color: rgba(255, 255, 255, 0.95);
      padding: 14px 18px;
      font-size: 14px;
      line-height: 1.6;
      resize: none;
      min-height: 72px;
      max-height: 120px;
      box-shadow: none !important;
      flex-grow: 1;
      border-radius: 16px 16px 0 0;

      &::placeholder {
        color: rgba(167, 139, 250, 0.45);
      }

      &:focus {
        outline: none;
        box-shadow: none;
        background: linear-gradient(135deg, rgba(26, 31, 46, 0.8) 0%, rgba(30, 35, 50, 0.8) 100%);
      }
    }
  }

  .input-divider {
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent 0%,
      rgba(102, 126, 234, 0.15) 10%,
      rgba(102, 126, 234, 0.3) 50%,
      rgba(102, 126, 234, 0.15) 90%,
      transparent 100%
    );
    margin: 0 12px;
    position: relative;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 1px;
      background: linear-gradient(
        90deg,
        transparent 0%,
        rgba(167, 139, 250, 0.2) 50%,
        transparent 100%
      );
      animation: dividerShimmer 3s ease-in-out infinite;
    }
  }

  .input-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 14px 12px 14px;
    gap: 12px;
    background: rgba(150, 160, 180, 0.08);
    border-radius: 0 0 16px 16px;

    .model-select {
      width: 140px;
      flex-shrink: 0;

      :deep(.el-input__wrapper) {
        background: rgba(102, 126, 234, 0.12);
        border: 1px solid rgba(102, 126, 234, 0.25);
        box-shadow: none;
        border-radius: 8px;
        transition: all 0.3s ease;
        height: 32px;

        &:hover {
          background: rgba(102, 126, 234, 0.18);
          border-color: rgba(102, 126, 234, 0.4);
        }

        &.is-focus {
          background: rgba(102, 126, 234, 0.2);
          border-color: #667eea;
          box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.15);
        }

        .el-input__inner {
          color: rgba(255, 255, 255, 0.95);
          font-size: 13px;
          font-weight: 500;
        }

        .el-input__suffix {
          .el-icon {
            color: rgba(167, 139, 250, 0.7);
          }
        }
      }
    }

    .send-button {
      width: 36px;
      height: 36px;
      padding: 0;
      flex-shrink: 0;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border: none;
      box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      position: relative;
      overflow: hidden;

      &::before {
        content: '';
        position: absolute;
        inset: 0;
        background: linear-gradient(135deg, rgba(255, 255, 255, 0.2) 0%, transparent 100%);
        opacity: 0;
        transition: opacity 0.3s ease;
      }

      &:hover:not(.is-disabled) {
        transform: translateY(-3px) scale(1.05);
        box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);

        &::before {
          opacity: 1;
        }
      }

      &:active:not(.is-disabled) {
        transform: translateY(-1px) scale(1.02);
        box-shadow: 0 4px 12px rgba(102, 126, 234, 0.5);
      }

      &.is-disabled {
        background: rgba(102, 126, 234, 0.25);
        box-shadow: none;
        cursor: not-allowed;
        opacity: 0.4;
        transform: none;
      }

      :deep(.el-icon) {
        font-size: 18px;
        font-weight: bold;
      }
    }
  }
}

@keyframes dividerShimmer {
  0%, 100% {
    opacity: 0.3;
    transform: translateX(-100%);
  }
  50% {
    opacity: 1;
    transform: translateX(100%);
  }
}

// ========== AI助手抽屉样式 ==========
.ai-drawer {
  :deep(.el-drawer) {
    background: linear-gradient(145deg, #3a3a38 0%, #2e2e2c 50%, #343432 100%) !important;
    box-shadow: -10px 0 50px rgba(0, 0, 0, 0.5) !important;

    // 添加背景纹理
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background:
        radial-gradient(circle at 20% 30%, rgba(180, 160, 100, 0.08) 0%, transparent 50%),
        radial-gradient(circle at 80% 70%, rgba(140, 130, 100, 0.05) 0%, transparent 50%);
      pointer-events: none;
    }

    // 网格背景
    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background-image:
        linear-gradient(rgba(180, 160, 100, 0.04) 1px, transparent 1px),
        linear-gradient(90deg, rgba(180, 160, 100, 0.04) 1px, transparent 1px);
      background-size: 50px 50px;
      pointer-events: none;
      opacity: 0.3;
    }
  }

  :deep(.el-drawer__header) {
    margin-bottom: 0;
    padding: 0;
    border-bottom: 1px solid rgba(180, 160, 100, 0.2) !important;
    background: rgba(50, 50, 48, 0.95) !important;
    backdrop-filter: blur(10px);
  }

  :deep(.el-drawer__body) {
    padding: 20px;
    position: relative;
    z-index: 1;
  }
}

// AI抽屉头部样式
.ai-drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px 8px 24px;

  .ai-header-content {
    display: flex;
    align-items: center;

    .ai-title-group {
      display: flex;
      flex-direction: column;
      gap: 4px;

      .ai-title {
        margin: 0;
        font-size: 18px;
        font-weight: 600;
        background: linear-gradient(135deg, #667eea 0%, #a78bfa 50%, #c084fc 100%);
        background-clip: text;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        letter-spacing: 0.5px;
      }

      .ai-subtitle {
        font-size: 11px;
        color: rgba(167, 139, 250, 0.6);
        font-weight: 400;
        letter-spacing: 1px;
        text-transform: uppercase;
      }
    }
  }

  .ai-status-indicator {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    background: rgba(102, 126, 234, 0.1);
    border: 1px solid rgba(102, 126, 234, 0.2);
    border-radius: 20px;

    .status-dot {
      width: 6px;
      height: 6px;
      border-radius: 50%;
      background: #10b981;
      box-shadow: 0 0 8px rgba(16, 185, 129, 0.8);
      animation: statusPulse 2s ease-in-out infinite;
    }

    .status-text {
      font-size: 12px;
      color: rgba(16, 185, 129, 0.9);
      font-weight: 500;
    }
  }
}

@keyframes aiIconFloat {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-5px) rotate(2deg);
  }
}

@keyframes statusPulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.6;
    transform: scale(0.9);
  }
}
</style>
