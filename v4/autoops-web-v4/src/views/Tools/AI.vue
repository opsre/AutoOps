<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  ArrowUp,
  ChatDotRound,
  Connection,
  Cpu,
  DataAnalysis,
  Delete,
  DocumentCopy,
  MagicStick,
  Monitor,
  RefreshRight
} from '@element-plus/icons-vue'
import aiApi from '@/api/ai'
import MarkdownTable from '@/components/MarkdownTable.vue'
import aiAssistantAvatar from '@/assets/image/AI小助手.svg'
import userAvatar from '@/assets/image/123.png'

const AI_SESSION_STORAGE_KEY = 'ops_ai_session_id'

const selectedModelId = ref(null)
const enabledModels = ref([])
const chatMessages = ref([])
const userInput = ref('')
const isLoading = ref(false)
const copiedIndex = ref(null)
const aiSessionId = ref(null)
const chatMessagesRef = ref(null)

const quickPrompts = [
  '你好，你是谁？',
  'ip/name 帮我分析当前服务是否有压力',
  'ip/name 帮我检查服务/CPU/内存/磁盘',
  'ip/name 帮我在这台服务器上安装nginx服务',
  'ip/name 帮我在这台服务器上执行命令 df -h',
  '生成一份 Kubernetes 故障排查 SOP'
]

const overviewMetrics = computed(() => [
  {
    label: '在线模型',
    value: enabledModels.value.length || '0',
    icon: Cpu
  },
  {
    label: '上下文轮次',
    value: chatMessages.value.filter(item => item.role === 'assistant').length || '0',
    icon: Connection
  }
])

const currentModel = computed(() => {
  return enabledModels.value.find(model => model.id === selectedModelId.value) || null
})

const hasMessages = computed(() => chatMessages.value.length > 0)

function isMarkdownTable(text) {
  if (!text || typeof text !== 'string') return false
  return /\|\s*[^|]+\s*\|/.test(text) && /-\s*-\s*-/.test(text)
}

function extractTableTitle(text) {
  if (!text || typeof text !== 'string') return ''
  const firstLine = text.split(/\r?\n/)[0]?.trim()
  if (firstLine && !firstLine.startsWith('|')) return firstLine
  return ''
}

function parseListItems(content) {
  if (!content) return []
  return content
    .split('\n')
    .map(item => item.trim())
    .filter(item => item.length > 0)
    .map(item => item.replace(/^[-*]\s*/, '').replace(/^\d+\.\s*/, ''))
}

async function loadEnabledModels() {
  try {
    const res = await aiApi.getEnabledAiModels()
    if (res.data && (res.data.code === 0 || res.data.code === 200)) {
      enabledModels.value = res.data.data || []

      if (enabledModels.value.length > 0 && !selectedModelId.value) {
        selectedModelId.value = enabledModels.value[0].id
      }

      if (enabledModels.value.length === 0) {
        ElMessage.warning('暂无可用的AI模型，请先在系统配置中添加并启用模型')
      }
      return
    }

    ElMessage.error(res.data?.message || res.data?.msg || '获取AI模型失败')
  } catch (error) {
    console.error('获取AI模型失败:', error)
    ElMessage.error('获取AI模型失败')
  }
}

function handleModelChange() {
  if (!currentModel.value) return
  ElMessage.success(`已切换至 ${currentModel.value.name}`)
}

function handleKeydown(event) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    sendMessage()
  }
}

function usePrompt(prompt) {
  userInput.value = prompt
}

function resetConversation() {
  chatMessages.value = []
  userInput.value = ''
  aiSessionId.value = null
  copiedIndex.value = null
  sessionStorage.removeItem(AI_SESSION_STORAGE_KEY)
}

async function sendMessage() {
  if (!userInput.value.trim()) return

  if (!selectedModelId.value) {
    ElMessage.warning('请先选择AI模型')
    return
  }

  const userMessage = {
    role: 'user',
    content: userInput.value.trim()
  }

  chatMessages.value.push(userMessage)
  const input = userInput.value
  userInput.value = ''
  isLoading.value = true

  try {
    const payload = {
      modelId: selectedModelId.value,
      messages: [userMessage]
    }

    if (aiSessionId.value) {
      payload.sessionId = aiSessionId.value
    }

    const res = await aiApi.aiChat(payload)

    if (res.data && (res.data.code === 0 || res.data.code === 200)) {
      const responseData = res.data.data || {}
      const assistantMessage = {
        role: 'assistant',
        content: responseData.text || '抱歉，我没有理解您的问题',
        contentBlocks: Array.isArray(responseData.content) ? responseData.content : null
      }

      chatMessages.value.push(assistantMessage)

      if (responseData.sessionId) {
        aiSessionId.value = responseData.sessionId
        sessionStorage.setItem(AI_SESSION_STORAGE_KEY, responseData.sessionId)
      }
      return
    }

    ElMessage.error(res.data?.message || res.data?.msg || 'AI对话失败')
    userInput.value = input
  } catch (error) {
    console.error('AI对话失败:', error)
    ElMessage.error(`AI对话失败: ${error.message || '未知错误'}`)
    userInput.value = input
  } finally {
    isLoading.value = false
  }
}

async function copyCode(code, index) {
  try {
    await navigator.clipboard.writeText(code)
    copiedIndex.value = index
    ElMessage.success('代码已复制到剪贴板')
    setTimeout(() => {
      copiedIndex.value = null
    }, 2000)
  } catch (error) {
    console.error('复制失败:', error)
    copyCodeFallback(code, index)
  }
}

function copyCodeFallback(code, index) {
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
      copiedIndex.value = index
      ElMessage.success('代码已复制到剪贴板')
      setTimeout(() => {
        copiedIndex.value = null
      }, 2000)
    } else {
      ElMessage.error('复制失败，请手动复制')
    }
  } catch (error) {
    console.error('降级复制也失败:', error)
    ElMessage.error('复制失败，请手动复制')
  }
}

function scrollToBottom() {
  if (!chatMessagesRef.value) return
  chatMessagesRef.value.scrollTop = chatMessagesRef.value.scrollHeight
}

watch(
  () => [chatMessages.value.length, isLoading.value],
  async () => {
    await nextTick()
    scrollToBottom()
  }
)

onMounted(() => {
  const savedSessionId = sessionStorage.getItem(AI_SESSION_STORAGE_KEY)
  if (savedSessionId) {
    aiSessionId.value = savedSessionId
  }
  loadEnabledModels()
})
</script>

<template>
  <div class="ai-workbench">
    <div class="ai-shell">
      <section class="hero-panel">
        <div class="hero-copy">
          <div class="hero-badge">
            <el-icon><MagicStick /></el-icon>
            <span>智能协作中枢</span>
          </div>
          <h1 class="hero-title">面向运维与交付场景的智能协作空间</h1>
          <p class="hero-desc">
            统一承载问答、方案生成、故障分析与代码输出，视觉上保持克制、紧凑和更偏系统化的科技质感。
          </p>
          <div class="hero-actions">
            <el-button type="primary" @click="resetConversation">
              <el-icon><RefreshRight /></el-icon>
              新建会话
            </el-button>
            <el-button plain @click="loadEnabledModels">
              <el-icon><Cpu /></el-icon>
              刷新模型
            </el-button>
          </div>
        </div>

        <div class="hero-metrics">
          <div v-for="item in overviewMetrics" :key="item.label" class="metric-card">
            <div class="metric-icon">
              <el-icon><component :is="item.icon" /></el-icon>
            </div>
            <div class="metric-value">{{ item.value }}</div>
            <div class="metric-label">{{ item.label }}</div>
          </div>
        </div>
      </section>

      <section class="workspace-grid">
        <aside class="sidebar-panel">
          <div class="panel-card model-panel">
            <div class="panel-head">
              <div>
                <div class="panel-kicker">Model Matrix</div>
                <h3>模型配置</h3>
              </div>
              <span class="panel-dot"></span>
            </div>

            <el-select
              v-model="selectedModelId"
              placeholder="选择模型"
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

            <div class="model-status-card">
              <div class="status-line">
                <span>当前模型</span>
                <strong>{{ currentModel?.name || '未选择' }}</strong>
              </div>
              <div class="status-line">
                <span>接入状态</span>
                <em>{{ enabledModels.length > 0 ? 'Available' : 'Unavailable' }}</em>
              </div>
            </div>
          </div>

          <div class="panel-card prompt-panel">
            <div class="panel-head">
              <div>
                <div class="panel-kicker">Quick Start</div>
                <h3>快捷指令</h3>
              </div>
            </div>

            <button
              v-for="prompt in quickPrompts"
              :key="prompt"
              type="button"
              class="prompt-chip"
              @click="usePrompt(prompt)"
            >
              <el-icon><ChatDotRound /></el-icon>
              <span>{{ prompt }}</span>
            </button>
          </div>

          <div class="panel-card capability-panel">
            <div class="panel-head">
              <div>
                <div class="panel-kicker">Capability</div>
                <h3>能力侧写</h3>
              </div>
            </div>

            <div class="capability-item">
              <el-icon><Monitor /></el-icon>
              <div>
                <strong>场景问答</strong>
                <p>适合运维诊断、发布策略、脚本生成与知识整理。</p>
              </div>
            </div>
            <div class="capability-item">
              <el-icon><Connection /></el-icon>
              <div>
                <strong>上下文续聊</strong>
                <p>会话 ID 持续复用，适合多轮追问与逐步收敛方案。</p>
              </div>
            </div>
          </div>
        </aside>

        <main class="chat-panel panel-card">
          <div class="chat-head">
            <div>
              <div class="panel-kicker">Conversation</div>
              <h2>智能对话中心</h2>
            </div>
            <div class="chat-head-right">
              <el-button
                v-if="hasMessages"
                plain
                size="small"
                class="clear-btn"
                @click="resetConversation"
              >
                <el-icon><Delete /></el-icon>
                清除对话
              </el-button>
              <div class="chat-state">
                <span class="state-dot" :class="{ active: !isLoading }"></span>
                <span>{{ isLoading ? '模型处理中' : '随时可提问' }}</span>
              </div>
            </div>
          </div>

          <div ref="chatMessagesRef" class="chat-messages">
            <div v-if="!hasMessages" class="empty-state">
              <div class="empty-orb"></div>
              <div class="empty-content">
                <h3>开始一段更有设计感的 AI 协作</h3>
                <p>你可以直接提问，也可以从左侧快捷指令开始。支持代码块、列表、标题和 Markdown 表格展示。</p>
              </div>
            </div>

            <template v-else>
              <div
                v-for="(message, index) in chatMessages"
                :key="index"
                :class="['message-item', message.role === 'user' ? 'user-message' : 'assistant-message']"
              >
                <div class="message-avatar">
                  <img
                    :src="message.role === 'user' ? userAvatar : aiAssistantAvatar"
                    :alt="message.role === 'user' ? '用户头像' : 'AI助手头像'"
                    class="avatar-image"
                  >
                </div>

                <div class="message-body">
                  <div class="message-meta">
                    <span>{{ message.role === 'user' ? '你' : currentModel?.name || 'AI 助手' }}</span>
                    <em>{{ message.role === 'user' ? 'Prompt' : 'Response' }}</em>
                  </div>

                  <div class="message-text">
                    <template v-if="message.role === 'user'">
                      {{ message.content }}
                    </template>

                    <template v-else-if="message.contentBlocks && message.contentBlocks.length > 0">
                      <div
                        v-for="(block, blockIndex) in message.contentBlocks"
                        :key="blockIndex"
                        :class="['content-block', `block-${block.type}`]"
                      >
                        <div v-if="block.type === 'text' && isMarkdownTable(block.content)" class="markdown-table-block">
                          <component
                            :is="`h4`"
                            v-if="extractTableTitle(block.content)"
                            class="heading-block"
                          >
                            {{ extractTableTitle(block.content) }}
                          </component>
                          <MarkdownTable :source="block" :maxHeight="360" />
                        </div>

                        <div v-else-if="block.type === 'code'" class="code-block">
                          <div class="code-header">
                            <span class="code-language">{{ block.language || 'text' }}</span>
                            <el-button type="primary" link @click="copyCode(block.content, `${index}-${blockIndex}`)">
                              <el-icon><DocumentCopy /></el-icon>
                              {{ copiedIndex === `${index}-${blockIndex}` ? '已复制' : '复制代码' }}
                            </el-button>
                          </div>
                          <pre><code>{{ block.content }}</code></pre>
                        </div>

                        <component
                          :is="`h${block.level || 3}`"
                          v-else-if="block.type === 'heading'"
                          class="heading-block"
                        >
                          {{ block.content }}
                        </component>

                        <ul v-else-if="block.type === 'list'" class="list-block">
                          <li v-for="(item, itemIndex) in parseListItems(block.content)" :key="itemIndex">
                            {{ item }}
                          </li>
                        </ul>

                        <div v-else class="text-block markdown-text">{{ block.content }}</div>
                      </div>
                    </template>

                    <template v-else>
                      <div v-if="isMarkdownTable(message.content)" class="markdown-table-block">
                        <component
                          :is="`h4`"
                          v-if="extractTableTitle(message.content)"
                          class="heading-block"
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
                  <img
                    :src="aiAssistantAvatar"
                    alt="AI助手头像"
                    class="avatar-image"
                  >
                </div>
                <div class="message-body">
                  <div class="message-meta">
                    <span>{{ currentModel?.name || 'AI 助手' }}</span>
                    <em>Thinking</em>
                  </div>
                  <div class="message-text typing-indicator">
                    <span></span><span></span><span></span>
                  </div>
                </div>
              </div>
            </template>
          </div>

          <div class="composer-panel">
            <div class="composer-topline">
              <span>输入问题后按 Enter 发送，Shift + Enter 换行</span>
              <span>{{ currentModel?.name || '请先选择模型' }}</span>
            </div>
            <div class="composer-box">
              <el-input
                v-model="userInput"
                type="textarea"
                :rows="4"
                :disabled="!selectedModelId"
                placeholder="请输入您的问题，例如：帮我设计一套发布回滚方案，并列出执行步骤。"
                class="chat-textarea"
                @keydown="handleKeydown"
              />
              <div class="composer-actions">
                <div class="composer-hint">
                  <el-icon><DataAnalysis /></el-icon>
                  <span>支持结构化回复、代码输出、表格展示</span>
                </div>
                <div class="composer-buttons">
                  <el-button
                    v-if="hasMessages"
                    class="clear-screen-btn"
                    :disabled="isLoading"
                    @click="resetConversation"
                    v-authority="['tools:ai:clear-screen']"
                  >
                    <el-icon><Delete /></el-icon>
                    清屏
                  </el-button>
                  <el-button
                    type="primary"
                    class="send-button"
                    :loading="isLoading"
                    :disabled="!userInput.trim() || !selectedModelId"
                    @click="sendMessage"
                    v-authority="['tools:ai:post']"
                  >
                    <el-icon><ArrowUp /></el-icon>
                    发送
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </main>
      </section>
    </div>
  </div>
</template>

<style scoped>
.ai-workbench {
  height: calc(100vh - 84px + 250px);
  min-height: 0;
  padding: 6px;
  box-sizing: border-box;
  overflow: hidden;
  background:
    radial-gradient(ellipse at 15% 0%, rgba(102, 126, 234, 0.07), transparent 40%),
    radial-gradient(ellipse at 85% 0%, rgba(118, 75, 162, 0.05), transparent 38%),
    linear-gradient(180deg, #eef1f6 0%, #e8ecf3 100%);
}

.ai-shell {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  gap: 10px;
}

.hero-panel,
.panel-card {
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(102, 126, 234, 0.12);
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.08), 0 1px 4px rgba(0, 0, 0, 0.04);
  backdrop-filter: blur(12px);
}

.hero-panel {
  position: relative;
  overflow: hidden;
  display: grid;
  grid-template-columns: minmax(0, 1.45fr) minmax(220px, 280px);
  gap: 12px;
  padding: 16px 20px;
}

/* 顶部品牌色渐变边框线 */
.hero-panel::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #667eea, #764ba2, #667eea);
  border-radius: 16px 16px 0 0;
  pointer-events: none;
}

/* 右上角装饰光晕 */
.hero-panel::after {
  content: '';
  position: absolute;
  top: -40px;
  right: -40px;
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(118, 75, 162, 0.06), transparent 70%);
  pointer-events: none;
}

.hero-copy,
.hero-metrics,
.workspace-grid,
.sidebar-panel,
.chat-panel {
  position: relative;
  z-index: 1;
}

.hero-badge,
.panel-kicker,
.chat-state,
.composer-topline span,
.message-meta em {
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 3px 10px 3px 8px;
  border-radius: 999px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.07));
  border: 1px solid rgba(102, 126, 234, 0.2);
  color: #667eea;
  font-size: 10px;
  font-weight: 700;
}

.hero-title {
  margin: 8px 0 5px;
  max-width: 720px;
  font-size: 20px;
  line-height: 1.25;
  font-weight: 700;
  background: linear-gradient(135deg, #2c3e50 20%, #667eea 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-desc {
  margin: 0;
  max-width: 680px;
  color: #606266;
  font-size: 12px;
  line-height: 1.5;
}

.hero-actions {
  display: flex;
  gap: 10px;
  margin-top: 14px;
}

.hero-actions :deep(.el-button--primary) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-color: transparent;
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.28);
}

.hero-actions :deep(.el-button--default) {
  color: #667eea;
  border-color: rgba(102, 126, 234, 0.25);
  background: rgba(255, 255, 255, 0.9);
}

.hero-metrics {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
  align-self: center;
}

.metric-card {
  display: grid;
  grid-template-columns: 36px 1fr;
  grid-template-areas:
    'icon value'
    'icon label';
  column-gap: 10px;
  row-gap: 2px;
  align-items: center;
  padding: 12px 14px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.06), rgba(118, 75, 162, 0.03));
  border: 1px solid rgba(102, 126, 234, 0.12);
  transition: all 0.2s ease;
}

.metric-card:hover {
  border-color: rgba(102, 126, 234, 0.25);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.1);
  transform: translateY(-1px);
}

.metric-icon {
  grid-area: icon;
  width: 34px;
  height: 34px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.metric-value {
  grid-area: value;
  font-size: 22px;
  line-height: 1;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea, #764ba2);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.metric-label {
  grid-area: label;
  color: #909399;
  font-size: 11px;
  font-weight: 500;
}

.workspace-grid {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  grid-template-rows: minmax(0, 1fr);
  flex: 1;
  min-height: 0;
  gap: 14px;
  align-items: stretch;
  overflow: visible;
}

.sidebar-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  gap: 12px;
  overflow-y: auto;
  padding-right: 0;
}

.panel-card {
  padding: 16px;
}

.panel-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding-bottom: 12px;
  margin-bottom: 12px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.08);
}

.panel-kicker {
  color: #667eea;
  font-size: 10px;
  font-weight: 700;
}

.panel-head h3,
.chat-head h2 {
  margin: 4px 0 0;
  color: #2c3e50;
  font-size: 15px;
  font-weight: 700;
}

.panel-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: linear-gradient(135deg, #67C23A, #22c55e);
  box-shadow: 0 0 0 4px rgba(103, 194, 58, 0.15);
  flex-shrink: 0;
  margin-top: 5px;
}

.model-select {
  width: 100%;
}

.model-status-card {
  margin-top: 12px;
  padding: 12px 14px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.06), rgba(118, 75, 162, 0.03));
  border: 1px solid rgba(102, 126, 234, 0.12);
}

.status-line {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  color: #909399;
  font-size: 12px;
}

.status-line + .status-line {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid rgba(102, 126, 234, 0.08);
}

.status-line strong,
.status-line em {
  color: #2c3e50;
  font-style: normal;
  font-weight: 600;
}

.prompt-chip {
  width: 100%;
  border: 1px solid rgba(102, 126, 234, 0.12);
  background: rgba(249, 250, 255, 0.9);
  border-radius: 10px;
  padding: 9px 12px;
  display: flex;
  align-items: flex-start;
  gap: 10px;
  text-align: left;
  color: #606266;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 12px;
  line-height: 1.45;
}

.prompt-chip + .prompt-chip {
  margin-top: 7px;
}

.prompt-chip:hover {
  transform: translateY(-2px);
  border-color: rgba(102, 126, 234, 0.3);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.1);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.06), rgba(118, 75, 162, 0.04));
  color: #667eea;
}

.prompt-chip .el-icon {
  margin-top: 2px;
  flex-shrink: 0;
  color: #667eea;
}

.capability-item {
  display: flex;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid rgba(102, 126, 234, 0.07);
}

.capability-item:last-child {
  padding-bottom: 0;
  border-bottom: none;
}

.capability-item .el-icon {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.07));
  color: #667eea;
}

.capability-item strong {
  display: block;
  margin-bottom: 5px;
  color: #2c3e50;
  font-size: 13px;
}

.capability-item p {
  margin: 0;
  color: #909399;
  font-size: 12px;
  line-height: 1.55;
}

.capability-panel {
  padding-bottom: 84px;
}

.chat-panel {
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.chat-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

.chat-head-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.clear-btn {
  color: #f56c6c;
  border-color: rgba(245, 108, 108, 0.25);
  background: rgba(255, 245, 245, 0.9);
  font-size: 12px;
  height: 28px;
  padding: 0 10px;
}

.clear-btn:hover {
  color: #fff;
  background: #f56c6c;
  border-color: #f56c6c;
}

.chat-state {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #909399;
  font-size: 11px;
  font-weight: 600;
}

.state-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #E6A23C;
  box-shadow: 0 0 0 4px rgba(230, 162, 60, 0.15);
}

.state-dot.active {
  background: #67C23A;
  box-shadow: 0 0 0 4px rgba(103, 194, 58, 0.15);
  animation: pulseGreen 2s ease-in-out infinite;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding: 14px 4px 8px;
}

.chat-messages::-webkit-scrollbar {
  width: 5px;
}

.chat-messages::-webkit-scrollbar-track {
  background: transparent;
}

.chat-messages::-webkit-scrollbar-thumb {
  background: rgba(102, 126, 234, 0.18);
  border-radius: 999px;
}

.chat-messages::-webkit-scrollbar-thumb:hover {
  background: rgba(102, 126, 234, 0.32);
}

.empty-state {
  min-height: 280px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 20px;
  text-align: center;
}

.empty-orb {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background:
    radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.9), transparent 28%),
    linear-gradient(135deg, #667eea, #764ba2);
  box-shadow:
    0 20px 40px rgba(102, 126, 234, 0.25),
    inset 0 0 20px rgba(255, 255, 255, 0.2);
  animation: floatOrb 4s ease-in-out infinite;
}

.empty-content h3 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #2c3e50, #667eea);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.empty-content p {
  margin: 0;
  max-width: 480px;
  color: #909399;
  font-size: 13px;
  line-height: 1.7;
}

.message-item {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  animation: messageSlideIn 0.3s ease;
}

.message-avatar {
  width: 36px;
  height: 36px;
  flex-shrink: 0;
  border-radius: 10px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  border: 1px solid rgba(102, 126, 234, 0.15);
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.08);
}

.user-message {
  flex-direction: row-reverse;
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: inherit;
}

.message-body {
  flex: 1;
  min-width: 0;
  max-width: calc(100% - 46px);
}

.message-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 5px;
}

.message-meta span {
  color: #2c3e50;
  font-size: 12px;
  font-weight: 700;
}

.message-meta em {
  color: #c0c4cc;
  font-size: 9px;
  font-style: normal;
  font-weight: 600;
}

.message-text {
  padding: 10px 14px;
  border-radius: 12px;
  background: #fafbff;
  border: 1px solid rgba(102, 126, 234, 0.1);
  color: #2c3e50;
  font-size: 13px;
  line-height: 1.65;
  word-break: break-word;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.04);
}

.assistant-message .message-text {
  border-top-left-radius: 4px;
  background: rgba(255, 255, 255, 0.98);
  border-color: rgba(102, 126, 234, 0.1);
}

.user-message .message-text {
  border-top-right-radius: 4px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-color: transparent;
  color: rgba(255, 255, 255, 0.96);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.25);
}

.user-message .message-meta span,
.user-message .message-meta em {
  color: #909399;
}

.assistant-message .message-meta span {
  color: #2c3e50;
}

.content-block + .content-block {
  margin-top: 10px;
}

.text-block {
  margin: 0;
}

.markdown-text {
  white-space: pre-wrap;
  line-height: 1.75;
  color: inherit;
}

.heading-block {
  margin: 0 0 8px;
  color: #5667d8;
  font-weight: 700;
}

.list-block {
  margin: 0;
  padding-left: 20px;
}

.list-block li + li {
  margin-top: 6px;
}

.code-block {
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(102, 126, 234, 0.1);
  background: linear-gradient(180deg, #1a1e2e, #13151f);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

.code-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 14px;
  background: rgba(102, 126, 234, 0.1);
  border-bottom: 1px solid rgba(102, 126, 234, 0.08);
}

.code-language {
  color: #667eea;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.code-block pre {
  margin: 0;
  padding: 14px;
  overflow-x: auto;
}

.code-block pre::-webkit-scrollbar {
  height: 4px;
}

.code-block pre::-webkit-scrollbar-thumb {
  background: rgba(102, 126, 234, 0.25);
  border-radius: 999px;
}

.code-block code {
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', 'Courier New', monospace;
  color: #c8d3f5;
  font-size: 12px;
  line-height: 1.6;
  white-space: pre;
}

.typing-indicator {
  display: flex;
  align-items: center;
  gap: 5px;
}

.typing-indicator span {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  animation: typingBounce 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

.composer-panel {
  margin-top: 6px;
  padding-top: 10px;
  border-top: 1px solid rgba(102, 126, 234, 0.1);
}

.composer-topline {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
}

.composer-topline span {
  color: #c0c4cc;
  font-size: 10px;
  font-weight: 600;
}

.composer-box {
  border-radius: 12px;
  border: 1px solid rgba(102, 126, 234, 0.15);
  background: rgba(255, 255, 255, 0.95);
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(102, 126, 234, 0.06);
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.composer-box:focus-within {
  border-color: rgba(102, 126, 234, 0.35);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.08), 0 2px 12px rgba(102, 126, 234, 0.08);
}

.composer-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 6px 10px 8px;
  border-top: 1px solid rgba(102, 126, 234, 0.08);
  background: rgba(249, 250, 255, 0.6);
}

.composer-hint {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #c0c4cc;
  font-size: 11px;
}

.composer-hint .el-icon {
  color: rgba(102, 126, 234, 0.5);
}

.send-button {
  min-width: 80px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.28);
}

.composer-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}

.clear-screen-btn {
  height: 32px;
  border-radius: 8px;
  border-color: rgba(245, 108, 108, 0.25);
  color: #f56c6c;
  background: rgba(255, 245, 245, 0.9);
  font-size: 13px;
  padding: 0 12px;
}

.clear-screen-btn:hover:not(:disabled) {
  color: #fff;
  background: #f56c6c;
  border-color: #f56c6c;
}

.send-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 18px rgba(102, 126, 234, 0.4);
}

.chat-textarea :deep(.el-textarea__inner) {
  min-height: 72px !important;
  padding: 10px 14px 6px;
  border: none;
  border-radius: 0;
  box-shadow: none;
  resize: none;
  background: transparent;
  color: #2c3e50;
  font-size: 13px;
  line-height: 1.5;
}

.chat-textarea :deep(.el-textarea__inner)::placeholder {
  color: #c0c4cc;
}

.model-select :deep(.el-input__wrapper),
.chat-textarea :deep(.el-textarea__inner),
:deep(.el-button) {
  transition: all 0.2s ease;
}

.model-select :deep(.el-input__wrapper) {
  border-radius: 10px;
  min-height: 40px;
  box-shadow: none;
  background: rgba(249, 250, 255, 0.95);
  border: 1px solid rgba(102, 126, 234, 0.15);
}

.model-select :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  border-color: rgba(102, 126, 234, 0.35);
}

:deep(.el-button) {
  border-radius: 8px;
  font-weight: 600;
}

:deep(.el-button--default) {
  border-color: rgba(102, 126, 234, 0.15);
}

:deep(.el-button:hover) {
  transform: translateY(-1px);
}

@keyframes messageSlideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes typingBounce {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.5;
  }
  30% {
    transform: translateY(-7px);
    opacity: 1;
  }
}

@keyframes floatOrb {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-8px);
  }
}

@keyframes pulseGreen {
  0%, 100% { box-shadow: 0 0 0 4px rgba(103, 194, 58, 0.15); }
  50% { box-shadow: 0 0 0 6px rgba(103, 194, 58, 0.25); }
}

@media (max-width: 1280px) {
  .hero-panel,
  .workspace-grid {
    grid-template-columns: 1fr;
  }

  .hero-metrics {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .ai-workbench {
    padding: 0;
    height: auto;
    min-height: calc(100vh - 84px + 250px);
    overflow: visible;
  }

  .hero-panel,
  .panel-card {
    border-radius: 12px;
  }

  .hero-panel {
    padding: 16px;
  }

  .hero-title {
    font-size: 18px;
  }

  .hero-metrics {
    grid-template-columns: 1fr;
  }

  .workspace-grid {
    gap: 10px;
    min-height: auto;
    overflow: visible;
  }

  .sidebar-panel {
    overflow: visible;
    padding-right: 0;
  }

  .chat-head,
  .composer-topline,
  .composer-actions {
    flex-direction: column;
    align-items: flex-start;
  }

  .send-button {
    width: 100%;
  }
}
</style>
