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

      <div class="chat-body">
        <el-scrollbar ref="scrollbarRef" class="chat-scrollbar">
          <div class="messages-container">
            <div v-if="messages.length === 0" class="welcome-message">
              <el-empty description="开始和 AI 助手对话吧">
                <el-button type="primary" @click="focusInput">
                  开始提问
                </el-button>
              </el-empty>
            </div>

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
                  v-if="message.content"
                  class="message-text"
                  v-html="renderMarkdown(message.content)"
                ></div>

                <div
                  v-if="message.attachments && message.attachments.length > 0"
                  class="message-attachments"
                >
                  <div
                    v-for="(attachment, attachmentIndex) in message.attachments"
                    :key="`${attachment.url}-${attachmentIndex}`"
                    class="attachment-item"
                  >
                    <template v-if="attachment.isImage">
                      <a
                        :href="attachment.url"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="attachment-image-wrap"
                      >
                        <img
                          class="attachment-image"
                          :src="attachment.previewUrl || attachment.url"
                          :alt="attachment.name"
                        />
                      </a>
                    </template>

                    <div class="attachment-meta">
                      <a
                        class="attachment-link"
                        :href="attachment.url"
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        {{ attachment.name || '附件' }}
                      </a>
                      <span v-if="attachment.size" class="attachment-size">
                        {{ formatFileSize(attachment.size) }}
                      </span>
                    </div>
                  </div>
                </div>

                <div class="message-time">
                  {{ formatTime(message.createdAt) }}
                </div>
              </div>
            </div>

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

      <div class="chat-footer">
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

        <input
          ref="fileInputRef"
          class="hidden-file-input"
          type="file"
          multiple
          @change="handleFilePick"
        />

        <div
          class="chat-input-wrapper"
          :class="{ active: isDragOver }"
          @dragenter.prevent.stop="handleDragEnter"
          @dragover.prevent.stop="handleDragOver"
          @dragleave.prevent.stop="handleDragLeave"
          @drop.prevent.stop="handleDrop"
        >
          <div v-if="pendingUploads.length > 0" class="input-upload-list">
            <div
              v-for="item in pendingUploads"
              :key="item.id"
              class="input-upload-item"
              :class="{
                uploading: item.uploading,
                error: !!item.error
              }"
            >
              <img
                v-if="item.isImage"
                class="input-upload-preview"
                :src="item.previewUrl"
                :alt="item.name"
              />
              <div v-else class="input-upload-file-icon">
                <el-icon><Document /></el-icon>
              </div>
              <div class="input-upload-main">
                <a
                  v-if="item.url && !item.error"
                  :href="item.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="input-upload-name"
                >
                  {{ item.name }}
                </a>
                <span v-else class="input-upload-name">{{ item.name }}</span>
                <div class="input-upload-meta">
                  <span>{{ formatFileSize(item.size) }}</span>
                  <span v-if="item.uploading">上传中...</span>
                  <span v-else-if="item.error" class="error-text">{{ item.error }}</span>
                  <span v-else>已上传</span>
                </div>
              </div>
              <el-button
                text
                type="danger"
                size="small"
                class="input-upload-remove"
                @click.stop="removePendingUpload(item.id)"
              >
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
          </div>

          <el-input
            ref="inputRef"
            v-model="inputMessage"
            type="textarea"
            :rows="3"
            placeholder="输入问题，按 Enter 发送（Shift+Enter 换行），也可直接拖入文件/图片"
            :disabled="loading"
            @keydown.enter.exact.prevent="sendMessage"
            resize="none"
            class="chat-input"
          />
        </div>

        <div class="input-actions">
          <div class="left-actions">
            <el-button
              type="info"
              size="small"
              plain
              @click="triggerFilePicker"
              :disabled="loading"
            >
              <el-icon><Upload /></el-icon>
              选择文件
            </el-button>
            <span class="upload-count" v-if="readyUploadCount > 0">
              已上传 {{ readyUploadCount }} 个附件
            </span>
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
  const fileInputRef = ref()
  const scrollbarRef = ref()
  const messages = ref([])
  const inputMessage = ref('')
  const loading = ref(false)
  const enabledPrompts = ref([])
  const hasSelection = ref(false)
  const selectedText = ref('')
  const streamController = ref(null)
  const pendingUploads = ref([])
  const isDragOver = ref(false)
  const dragDepth = ref(0)

  const uploadInProgressCount = computed(() => {
    return pendingUploads.value.filter((item) => item.uploading).length
  })

  const readyUploadCount = computed(() => {
    return pendingUploads.value.filter((item) => !item.uploading && !item.error && item.url).length
  })

  const canSend = computed(() => {
    const hasText = inputMessage.value.trim().length > 0
    return (hasText || readyUploadCount.value > 0) && !loading.value
  })

  const loadEnabledPrompts = async () => {
    const res = await getEnabledPrompts()
    if (res.code === 0) {
      enabledPrompts.value = res.data || []
    }
  }

  const handleSelectionChange = () => {
    const selection = window.getSelection()
    const text = selection?.toString?.().trim() || ''
    hasSelection.value = !!text
    selectedText.value = text
  }

  const focusInput = () => {
    inputRef.value?.focus()
  }

  const usePrompt = (content) => {
    inputMessage.value = content
    sendMessage()
  }

  const useSelectedText = () => {
    if (selectedText.value) {
      inputMessage.value = `请分析以下内容：${selectedText.value}\n\n`
      focusInput()
    }
  }

  const stopGenerating = () => {
    if (streamController.value) {
      streamController.value.abort()
      streamController.value = null
    }
    loading.value = false
    ElMessage.info('已停止生成')
  }

  const createUploadId = () => {
    return `${Date.now()}-${Math.random().toString(36).slice(2, 10)}`
  }

  const isImageFile = (file) => {
    if (!file) return false
    if (file.type && file.type.startsWith('image/')) return true
    return /\.(png|jpe?g|gif|webp|bmp|svg)$/i.test(file.name || '')
  }

  const buildMockUploadUrl = (file, id) => {
    const safeName = encodeURIComponent((file.name || 'file').replace(/\s+/g, '_'))
    return `https://mock-upload.local/aiassistant/${Date.now()}-${id}-${safeName}`
  }

  const mockUploadFile = (file, id) => {
    const delay = 300 + Math.floor(Math.random() * 800)
    return new Promise((resolve) => {
      window.setTimeout(() => {
        resolve(buildMockUploadUrl(file, id))
      }, delay)
    })
  }

  const releasePreview = (item) => {
    if (item?.previewUrl && item.previewUrl.startsWith('blob:')) {
      URL.revokeObjectURL(item.previewUrl)
    }
  }

  const releaseMessagePreviews = () => {
    messages.value.forEach((message) => {
      ;(message.attachments || []).forEach((attachment) => {
        releasePreview(attachment)
      })
    })
  }

  const removePendingUpload = (id) => {
    const index = pendingUploads.value.findIndex((item) => item.id === id)
    if (index < 0) return
    const [removed] = pendingUploads.value.splice(index, 1)
    releasePreview(removed)
  }

  const enqueueUpload = async (file) => {
    const id = createUploadId()
    const image = isImageFile(file)
    const record = {
      id,
      name: file.name || '未命名文件',
      type: file.type || 'application/octet-stream',
      size: file.size || 0,
      isImage: image,
      previewUrl: image ? URL.createObjectURL(file) : '',
      url: '',
      uploading: true,
      error: ''
    }

    pendingUploads.value.push(record)

    try {
      const uploadedUrl = await mockUploadFile(file, id)
      record.url = uploadedUrl
      record.uploading = false
    } catch (error) {
      record.uploading = false
      record.error = error?.message || '上传失败'
    }
    return record
  }

  const handleFiles = async (fileList) => {
    const files = Array.from(fileList || [])
    if (files.length === 0) return

    const uploadResults = await Promise.all(files.map((file) => enqueueUpload(file)))
    const successCount = uploadResults.filter((item) => !item.error).length
    if (successCount > 0) {
      ElMessage.success(`已完成 ${successCount} 个文件的模拟上传`)
    }
  }

  const triggerFilePicker = () => {
    fileInputRef.value?.click()
  }

  const handleFilePick = async (event) => {
    const files = event?.target?.files
    await handleFiles(files)
    if (event?.target) {
      event.target.value = ''
    }
  }

  const handleDragEnter = () => {
    dragDepth.value += 1
    isDragOver.value = true
  }

  const handleDragOver = (event) => {
    event.dataTransfer.dropEffect = 'copy'
    isDragOver.value = true
  }

  const handleDragLeave = () => {
    dragDepth.value -= 1
    if (dragDepth.value <= 0) {
      dragDepth.value = 0
      isDragOver.value = false
    }
  }

  const handleDrop = async (event) => {
    dragDepth.value = 0
    isDragOver.value = false
    await handleFiles(event.dataTransfer?.files)
  }

  const sendMessage = async () => {
    if (loading.value) return

    const plainText = inputMessage.value.trim()

    const outgoingAttachments = pendingUploads.value
      .filter((item) => !item.uploading && !item.error && item.url)
      .map((item) => ({
        id: item.id,
        name: item.name,
        type: item.type,
        size: item.size,
        url: item.url,
        isImage: item.isImage,
        previewUrl: item.previewUrl
      }))

    if (!plainText && outgoingAttachments.length === 0) return
    if (uploadInProgressCount.value > 0) {
      ElMessage.info('部分文件仍在上传，本次先发送已上传的附件')
    }

    const content = plainText || '请结合这些附件 URL 进行分析。'

    messages.value.push({
      role: 'user',
      content,
      attachments: outgoingAttachments,
      createdAt: new Date()
    })

    const outgoingIDs = new Set(outgoingAttachments.map((item) => item.id))

    inputMessage.value = ''
    pendingUploads.value = pendingUploads.value.filter((item) => !outgoingIDs.has(item.id))
    loading.value = true

    await nextTick()
    scrollToBottom()

    try {
      const messagesPayload = messages.value.map((m) => ({
        role: m.role,
        content: m.content,
        attachments: (m.attachments || []).map((item) => ({
          name: item.name,
          url: item.url,
          type: item.type,
          size: item.size
        }))
      }))

      const assistantMessage = {
        role: 'assistant',
        content: '',
        attachments: [],
        createdAt: new Date()
      }
      messages.value.push(assistantMessage)

      streamController.value = chatStream(
        {
          messages: messagesPayload,
          attachments: outgoingAttachments.map((item) => ({
            name: item.name,
            url: item.url,
            type: item.type,
            size: item.size
          }))
        },
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
        () => {
          loading.value = false
          streamController.value = null
        }
      )
    } catch (error) {
      ElMessage.error(error?.message || 'AI 响应失败')
      loading.value = false
      messages.value.push({
        role: 'assistant',
        content: `抱歉，发生错误：${error?.message || '未知错误'}`,
        attachments: [],
        createdAt: new Date()
      })
    }
  }

  const scrollToBottom = () => {
    const container = scrollbarRef.value?.wrapRef
    if (container) {
      container.scrollTop = container.scrollHeight
    }
  }

  const renderMarkdown = (text) => {
    if (!text) return ''
    return md.render(text)
  }

  const formatTime = (date) => {
    if (!date) return ''
    const d = new Date(date)
    return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }

  const formatFileSize = (size) => {
    const value = Number(size) || 0
    if (value < 1024) return `${value}B`
    if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)}KB`
    if (value < 1024 * 1024 * 1024) return `${(value / 1024 / 1024).toFixed(1)}MB`
    return `${(value / 1024 / 1024 / 1024).toFixed(1)}GB`
  }

  const clearChat = () => {
    if (messages.value.length === 0) return
    releaseMessagePreviews()
    releasePendingPreviews()
    messages.value = []
    ElMessage.success('对话已清空')
  }

  const releasePendingPreviews = () => {
    pendingUploads.value.forEach((item) => releasePreview(item))
    pendingUploads.value = []
  }

  onMounted(() => {
    loadEnabledPrompts()
    document.addEventListener('selectionchange', handleSelectionChange)
    focusInput()
  })

  onUnmounted(() => {
    document.removeEventListener('selectionchange', handleSelectionChange)
    if (streamController.value) {
      streamController.value.abort()
      streamController.value = null
    }
    releasePendingPreviews()
    releaseMessagePreviews()
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
    overflow: hidden;
  }

  .chat-card :deep(.el-card__body) {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
    padding: 0;
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
    background: #f8fafc;
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
    margin-bottom: 20px;
    gap: 12px;
    animation: fadeIn 0.24s ease;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(8px);
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
    max-width: 74%;
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

  .message-attachments {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-top: 10px;
  }

  .attachment-item {
    border: 1px solid rgba(64, 158, 255, 0.28);
    border-radius: 8px;
    padding: 8px;
    background: rgba(64, 158, 255, 0.05);
  }

  .message-item.user .attachment-item {
    border-color: rgba(255, 255, 255, 0.4);
    background: rgba(255, 255, 255, 0.14);
  }

  .attachment-image-wrap {
    display: block;
    width: 160px;
    max-width: 100%;
  }

  .attachment-image {
    width: 100%;
    border-radius: 6px;
    display: block;
  }

  .attachment-meta {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 6px;
    flex-wrap: wrap;
  }

  .attachment-link {
    color: #1f6feb;
    text-decoration: none;
    font-weight: 500;
  }

  .attachment-link:hover {
    text-decoration: underline;
  }

  .message-item.user .attachment-link {
    color: #ffffff;
  }

  .attachment-size {
    color: #64748b;
    font-size: 12px;
  }

  .message-item.user .attachment-size {
    color: rgba(255, 255, 255, 0.8);
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

  .chat-footer {
    position: sticky;
    bottom: 0;
    z-index: 3;
    border-top: 1px solid #e4e7ed;
    background: #ffffff;
    padding: 12px 16px;
    box-shadow: 0 -4px 16px rgba(15, 23, 42, 0.06);
  }

  .quick-prompts {
    padding: 10px 12px;
    background: #f8fafc;
    border-radius: 8px;
    margin-bottom: 12px;
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
    transition: all 0.25s;
  }

  .prompt-tag:hover {
    transform: translateY(-1px);
  }

  .hidden-file-input {
    display: none;
  }

  .chat-input-wrapper {
    border: 1px solid #dcdfe6;
    border-radius: 10px;
    background: #ffffff;
    margin-bottom: 10px;
    transition: all 0.2s;
    overflow: hidden;
  }

  .chat-input-wrapper.active {
    border-color: #409eff;
    box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.14);
    background: #f5f9ff;
  }

  .input-upload-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    padding: 10px;
    border-bottom: 1px dashed #e5e7eb;
    background: #f8fafc;
  }

  .input-upload-item {
    max-width: 220px;
    display: flex;
    align-items: center;
    gap: 8px;
    border: 1px solid #e5e7eb;
    background: #ffffff;
    border-radius: 8px;
    padding: 6px 8px;
    min-width: 0;
  }

  .input-upload-item.uploading {
    border-color: #93c5fd;
    background: #eff6ff;
  }

  .input-upload-item.error {
    border-color: #fca5a5;
    background: #fff1f2;
  }

  .input-upload-preview,
  .input-upload-file-icon {
    width: 34px;
    height: 34px;
    border-radius: 6px;
    flex-shrink: 0;
  }

  .input-upload-preview {
    object-fit: cover;
    display: block;
  }

  .input-upload-file-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    background: #eef2ff;
    color: #4f46e5;
    font-size: 16px;
  }

  .input-upload-main {
    min-width: 0;
    flex: 1;
  }

  .input-upload-name {
    display: block;
    color: #0f172a;
    font-size: 12px;
    text-decoration: none;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .input-upload-name:hover {
    text-decoration: underline;
  }

  .input-upload-meta {
    margin-top: 2px;
    font-size: 11px;
    color: #64748b;
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
  }

  .input-upload-remove {
    padding: 0;
    min-height: auto;
    flex-shrink: 0;
  }

  .error-text {
    color: #ef4444;
  }

  .chat-input {
    margin-bottom: 0;
  }

  .chat-input :deep(.el-textarea__inner) {
    border: 0;
    box-shadow: none;
    resize: none;
    min-height: 84px;
    border-radius: 0;
    background: transparent;
  }

  .chat-input :deep(.el-textarea__inner):focus {
    box-shadow: none;
  }

  .input-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
  }

  .left-actions,
  .right-actions {
    display: flex;
    gap: 8px;
    align-items: center;
    flex-wrap: wrap;
  }

  .upload-count {
    color: #475569;
    font-size: 12px;
  }

  @media (max-width: 768px) {
    .ai-chat-page {
      padding: 12px;
      height: calc(100vh - 100px);
    }

    .message-content {
      max-width: 88%;
    }

    .chat-footer {
      padding: 10px 12px;
    }

    .attachment-image-wrap {
      width: 120px;
    }
  }
</style>
