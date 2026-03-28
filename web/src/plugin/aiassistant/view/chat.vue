<template>
  <div class="ai-chat-page">
    <el-card class="chat-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon :size="20"><ChatDotRound /></el-icon>
            <span class="title">AI 智能助手</span>
          </div>
          <div class="header-actions">
            <el-button
              type="info"
              size="small"
              @click="clearChat"
              :disabled="messages.length === 0"
            >
              <el-icon><Delete /></el-icon>
              清空对话
            </el-button>
            <el-button
              v-if="loading"
              type="danger"
              size="small"
              @click="stopGenerating"
            >
              <el-icon><VideoPause /></el-icon>
              停止生成
            </el-button>
          </div>
        </div>
      </template>

      <!-- 聊天区域 -->
      <div class="chat-body">
        <el-scrollbar ref="scrollbarRef" class="chat-scrollbar">
          <div class="messages-container">
            <!-- 欢迎消息 -->
            <div v-if="messages.length === 0" class="welcome-message">
              <el-empty description="开始和 AI 助手对话吧">
                <el-button type="primary" @click="focusInput">
                  开始提问
                </el-button>
              </el-empty>
            </div>

            <!-- 消息列表 -->
            <div
              v-for="(message, index) in messages"
              :key="index"
              :class="['message-item', message.role]"
            >
              <div class="message-avatar">
                <el-icon v-if="message.role === 'user'" :size="20"><User /></el-icon>
                <el-icon v-else :size="20"><Cpu /></el-icon>
              </div>
              <div class="message-content">
                <div
                  class="message-text"
                  v-html="renderMarkdown(message.content)"
                ></div>
                <div class="message-time">
                  {{ formatTime(message.createdAt) }}
                </div>
              </div>
            </div>

            <!-- 加载中 -->
            <div v-if="loading" class="message-item assistant">
              <div class="message-avatar">
                <el-icon :size="20"><Cpu /></el-icon>
              </div>
              <div class="message-content">
                <div class="typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            </div>
          </div>
        </el-scrollbar>
      </div>

      <!-- 快捷提示词 -->
      <div v-if="enabledPrompts.length > 0" class="quick-prompts">
        <div class="prompts-label">
          <el-icon><Star /></el-icon>
          快捷提问：
        </div>
        <div class="prompts-list">
          <el-tag
            v-for="prompt in enabledPrompts"
            :key="prompt.ID"
            class="prompt-tag"
            size="large"
            round
            effect="plain"
            @click="usePrompt(prompt.content)"
          >
            {{ prompt.title }}
          </el-tag>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="chat-footer">
        <el-input
          ref="inputRef"
          v-model="inputMessage"
          type="textarea"
          :rows="3"
          placeholder="输入问题，按 Enter 发送（Shift+Enter 换行）"
          :disabled="loading"
          @keydown.enter.exact="sendMessage"
          resize="none"
          class="chat-input"
        />
        <div class="input-actions">
          <div class="left-actions">
            <el-button
              v-if="hasSelection"
              type="info"
              size="small"
              plain
              @click="useSelectedText"
            >
              <el-icon><Document /></el-icon>
              使用选中内容 ({{ selectedText.length }}字)
            </el-button>
          </div>
          <div class="right-actions">
            <el-button
              type="primary"
              :loading="loading"
              :disabled="!canSend"
              @click="sendMessage"
            >
              <el-icon><Promotion /></el-icon>
              发送
            </el-button>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
  import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getEnabledPrompts, chatStream } from '@/plugin/aiassistant/api/assistant'
  import Markdown from 'markdown-it'

  defineOptions({
    name: 'AIChatPage'
  })

  const md = new Markdown({
    html: false,
    breaks: true,
    linkify: true
  })

  const inputRef = ref()
  const scrollbarRef = ref()
  const messages = ref([])
  const inputMessage = ref('')
  const loading = ref(false)
  const enabledPrompts = ref([])
  const hasSelection = ref(false)
  const selectedText = ref('')
  const streamController = ref(null)

  const canSend = computed(() => {
    return inputMessage.value.trim() && !loading.value
  })

  // 加载已启用的提示词
  const loadEnabledPrompts = async () => {
    const res = await getEnabledPrompts()
    if (res.code === 0) {
      enabledPrompts.value = res.data || []
    }
  }

  // 监听页面选中文字
  const handleSelectionChange = () => {
    const selection = window.getSelection()
    const text = selection.toString().trim()
    hasSelection.value = !!text
    selectedText.value = text
  }

  // 聚焦输入框
  const focusInput = () => {
    inputRef.value?.focus()
  }

  // 使用快捷提示词
  const usePrompt = (content) => {
    inputMessage.value = content
    sendMessage()
  }

  // 使用选中文字
  const useSelectedText = () => {
    if (selectedText.value) {
      inputMessage.value = `请分析以下内容：${selectedText.value}\n\n`
      focusInput()
    }
  }

  // 停止生成
  const stopGenerating = () => {
    if (streamController.value) {
      streamController.value.abort()
      streamController.value = null
    }
    loading.value = false
    ElMessage.info('已停止生成')
  }

  // 发送消息
  const sendMessage = async () => {
    const message = inputMessage.value.trim()
    if (!message || loading.value) return

    // 添加用户消息
    messages.value.push({
      role: 'user',
      content: message,
      createdAt: new Date()
    })

    inputMessage.value = ''
    loading.value = true

    await nextTick()
    scrollToBottom()

    try {
      // 构建消息历史
      const messagesPayload = messages.value.map((m) => ({
        role: m.role,
        content: m.content
      }))

      // 添加助手消息占位
      const assistantMessage = {
        role: 'assistant',
        content: '',
        createdAt: new Date()
      }
      messages.value.push(assistantMessage)

      // 使用 SSE 流式聊天
      streamController.value = chatStream(
        { messages: messagesPayload },
        // onMessage: 接收流式内容
        (delta) => {
          assistantMessage.content += delta
          scrollToBottom()
        },
        // onError: 错误处理
        (error) => {
          ElMessage.error(error.message || 'AI 响应失败')
          if (assistantMessage.content === '') {
            assistantMessage.content = `抱歉，发生错误：${error.message}`
          }
        },
        // onComplete: 完成回调
        () => {
          loading.value = false
          streamController.value = null
        }
      )
    } catch (error) {
      ElMessage.error(error.message || 'AI 响应失败')
      loading.value = false
      // 添加错误消息
      messages.value.push({
        role: 'assistant',
        content: `抱歉，发生错误：${error.message}`,
        createdAt: new Date()
      })
    }
  }

  // 滚动到底部
  const scrollToBottom = () => {
    const container = scrollbarRef.value?.wrapRef
    if (container) {
      container.scrollTop = container.scrollHeight
    }
  }

  // 渲染 Markdown
  const renderMarkdown = (text) => {
    if (!text) return ''
    return md.render(text)
  }

  // 格式化时间
  const formatTime = (date) => {
    if (!date) return ''
    const d = new Date(date)
    return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }

  // 清空对话
  const clearChat = () => {
    if (messages.value.length === 0) return
    messages.value = []
    ElMessage.success('对话已清空')
  }

  onMounted(() => {
    loadEnabledPrompts()
    document.addEventListener('selectionchange', handleSelectionChange)
    focusInput()
  })

  onUnmounted(() => {
    document.removeEventListener('selectionchange', handleSelectionChange)
    // 清理未完成的流
    if (streamController.value) {
      streamController.value.abort()
    }
  })
</script>

<style scoped>
  .ai-chat-page {
    padding: 20px;
    height: calc(100vh - 120px);
  }

  .chat-card {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .header-left .title {
    font-size: 16px;
    font-weight: 600;
  }

  .chat-body {
    flex: 1;
    min-height: 0;
    margin-bottom: 16px;
  }

  .chat-scrollbar {
    height: 100%;
  }

  .messages-container {
    padding: 20px;
  }

  .welcome-message {
    display: flex;
    justify-content: center;
    padding: 60px 0;
  }

  .message-item {
    display: flex;
    align-items: flex-start;
    margin-bottom: 24px;
    gap: 12px;
    animation: fadeIn 0.3s ease;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .message-item.user {
    flex-direction: row-reverse;
  }

  .message-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .message-item.user .message-avatar {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  .message-item.assistant .message-avatar {
    background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
    color: white;
  }

  .message-content {
    max-width: 70%;
    background: white;
    padding: 12px 16px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }

  .message-item.user .message-content {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  .message-text {
    line-height: 1.6;
    word-break: break-word;
  }

  .message-text :deep(p) {
    margin: 0.5em 0;
  }

  .message-text :deep(p:first-child) {
    margin-top: 0;
  }

  .message-text :deep(p:last-child) {
    margin-bottom: 0;
  }

  .message-text :deep(code) {
    background: rgba(0, 0, 0, 0.06);
    padding: 2px 6px;
    border-radius: 4px;
    font-family: 'Consolas', 'Monaco', monospace;
    font-size: 0.9em;
  }

  .message-item.user .message-text :deep(code) {
    background: rgba(255, 255, 255, 0.2);
  }

  .message-text :deep(pre) {
    background: #282c34;
    color: #abb2bf;
    padding: 12px;
    border-radius: 6px;
    overflow-x: auto;
    margin: 8px 0;
  }

  .message-text :deep(pre code) {
    background: transparent;
    padding: 0;
  }

  .message-time {
    font-size: 12px;
    opacity: 0.6;
    margin-top: 6px;
    text-align: right;
  }

  .typing-indicator {
    display: flex;
    gap: 4px;
    padding: 8px 0;
  }

  .typing-indicator span {
    width: 8px;
    height: 8px;
    background: #909399;
    border-radius: 50%;
    animation: typing 1.4s infinite;
  }

  .typing-indicator span:nth-child(2) {
    animation-delay: 0.2s;
  }

  .typing-indicator span:nth-child(3) {
    animation-delay: 0.4s;
  }

  @keyframes typing {
    0%, 100% {
      transform: translateY(0);
      opacity: 0.4;
    }
    50% {
      transform: translateY(-4px);
      opacity: 1;
    }
  }

  .quick-prompts {
    padding: 12px 20px;
    background: #f5f7fa;
    border-radius: 8px;
    margin-bottom: 16px;
  }

  .prompts-label {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 13px;
    color: #606266;
    margin-bottom: 8px;
  }

  .prompts-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .prompt-tag {
    cursor: pointer;
    transition: all 0.3s;
  }

  .prompt-tag:hover {
    transform: translateY(-2px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  }

  .chat-footer {
    border-top: 1px solid #e4e7ed;
    padding-top: 16px;
  }

  .chat-input {
    margin-bottom: 12px;
  }

  .chat-input :deep(.el-textarea__inner) {
    resize: none;
  }

  .input-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .left-actions,
  .right-actions {
    display: flex;
    gap: 8px;
    align-items: center;
  }
</style>
