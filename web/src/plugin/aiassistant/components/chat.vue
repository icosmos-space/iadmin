<template>
  <div class="ai-chat-container">
    <!-- 聊天消息列表 -->
    <div class="chat-messages" ref="messagesContainer">
      <div
        v-for="(message, index) in messages"
        :key="index"
        :class="['message-item', message.role]"
      >
        <div class="message-avatar">
          <el-icon v-if="message.role === 'user'"><User /></el-icon>
          <el-icon v-else><Cpu /></el-icon>
        </div>
        <div class="message-content">
          <div class="message-text" v-html="renderMarkdown(message.content)"></div>
          <div class="message-time">{{ formatTime(message.createdAt) }}</div>
        </div>
      </div>
      
      <!-- 加载中状态 -->
      <div v-if="loading" class="message-item assistant">
        <div class="message-avatar">
          <el-icon><Cpu /></el-icon>
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

    <!-- 快捷提示词 -->
    <div v-if="showPrompts && enabledPrompts.length > 0" class="quick-prompts">
      <el-tag
        v-for="prompt in enabledPrompts"
        :key="prompt.ID"
        class="prompt-tag"
        size="large"
        round
        @click="usePrompt(prompt)"
      >
        {{ prompt.title }}
      </el-tag>
    </div>

    <!-- 输入区域 -->
    <div class="chat-input-area">
      <el-input
        v-model="inputMessage"
        type="textarea"
        :rows="3"
        placeholder="输入消息，按 Enter 发送（Shift+Enter 换行）"
        :disabled="loading"
        @keydown.enter.exact.prevent="sendMessage"
        resize="none"
      />
      <div class="input-actions">
        <el-button
          v-if="hasSelection"
          type="info"
          size="small"
          @click="useSelectedText"
        >
          <el-icon><Document /></el-icon>
          使用选中内容
        </el-button>
        <el-button
          type="primary"
          :loading="loading"
          @click="sendMessage"
          :disabled="!canSend"
        >
          <el-icon><Promotion /></el-icon>
          发送
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { chatStream, getEnabledPrompts } from '@/plugin/aiassistant/api/assistant'
  import Markdown from 'markdown-it'

  defineOptions({
    name: 'AIChat'
  })

  const props = defineProps({
    showPrompts: {
      type: Boolean,
      default: true
    },
    initialPrompt: {
      type: String,
      default: ''
    }
  })

  const emit = defineEmits(['send', 'clear'])

  const md = new Markdown({
    html: false,
    breaks: true,
    linkify: true
  })

  const messages = ref([])
  const inputMessage = ref('')
  const loading = ref(false)
  const enabledPrompts = ref([])
  const messagesContainer = ref(null)
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

  // 使用快捷提示词
  const usePrompt = (prompt) => {
    inputMessage.value = prompt.content
    sendMessage()
  }

  // 使用选中文字
  const useSelectedText = () => {
    if (selectedText.value) {
      inputMessage.value = `请分析以下内容：${selectedText.value}\n\n`
    }
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
    emit('send', message)

    await nextTick()
    scrollToBottom()

    const messagesPayload = messages.value.map((m) => ({
      role: m.role,
      content: m.content
    }))

    const assistantMessage = {
      role: 'assistant',
      content: '',
      createdAt: new Date()
    }
    messages.value.push(assistantMessage)

    streamController.value = chatStream(
      { messages: messagesPayload },
      (delta) => {
        assistantMessage.content += delta
        scrollToBottom()
      },
      (error) => {
        ElMessage.error(error?.message || 'AI 响应失败')
        if (!assistantMessage.content) {
          assistantMessage.content = `抱歉，发生错误：${error?.message || '未知错误'}`
        }
        loading.value = false
        streamController.value = null
      },
      async () => {
        loading.value = false
        streamController.value = null
        await nextTick()
        scrollToBottom()
      }
    )
  }

  // 滚动到底部
  const scrollToBottom = () => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
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
  const clearMessages = () => {
    messages.value = []
    emit('clear')
  }

  // 暴露方法给父组件
  defineExpose({
    clearMessages,
    addMessage: (role, content) => {
      messages.value.push({ role, content, createdAt: new Date() })
    }
  })

  onMounted(() => {
    loadEnabledPrompts()
    document.addEventListener('selectionchange', handleSelectionChange)
    
    // 如果有初始提示
    if (props.initialPrompt) {
      inputMessage.value = props.initialPrompt
    }
  })

  onUnmounted(() => {
    document.removeEventListener('selectionchange', handleSelectionChange)
    if (streamController.value) {
      streamController.value.abort()
      streamController.value = null
    }
  })
</script>

<style scoped>
  .ai-chat-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: #f5f7fa;
  }

  .chat-messages {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
  }

  .message-item {
    display: flex;
    align-items: flex-start;
    margin-bottom: 20px;
    gap: 12px;
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
    font-size: 20px;
    flex-shrink: 0;
  }

  .message-item.user .message-avatar {
    background: #409eff;
    color: white;
  }

  .message-item.assistant .message-avatar {
    background: #67c23a;
    color: white;
  }

  .message-content {
    max-width: 70%;
    background: white;
    padding: 12px 16px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .message-item.user .message-content {
    background: #409eff;
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
    background: white;
    border-top: 1px solid #e4e7ed;
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
  }

  .chat-input-area {
    padding: 16px 20px;
    background: white;
    border-top: 1px solid #e4e7ed;
  }

  .input-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 12px;
  }
</style>
