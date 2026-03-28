<template>
  <div v-if="visible" class="gva-ai-chat-root">
    <transition name="gva-ai-chat-panel">
      <section
        v-show="panelVisible"
        class="gva-ai-chat-panel"
        :style="panelStyle"
        @mousedown.stop
      >
        <el-card
          shadow="always"
          class="gva-ai-chat-panel-card"
          @dragenter.prevent="onDragEnter"
          @dragover.prevent="onDragOver"
          @dragleave.prevent="onDragLeave"
          @drop.prevent="onDrop"
        >
          <template #header>
            <div class="gva-ai-chat-header">
              <div class="gva-ai-chat-head-left">
                <el-avatar :size="34" class="gva-ai-chat-head-avatar">
                  <el-icon>
                    <Service />
                  </el-icon>
                </el-avatar>
                <div class="gva-ai-chat-head-text">
                  <div class="gva-ai-chat-title">{{ title }}</div>
                  <div class="gva-ai-chat-subtitle">{{ subtitle }}</div>
                  <div class="gva-ai-chat-mode-tags">
                    <el-tag size="small" type="primary" effect="plain">
                      {{ chatModeLabel }}
                    </el-tag>
                    <el-tag size="small" type="success" effect="plain">
                      {{ contextModeLabel }}
                    </el-tag>
                    <el-tag size="small" type="warning" effect="plain">
                      图片拖入支持
                    </el-tag>
                  </div>
                </div>
              </div>
              <div class="gva-ai-chat-actions">
                <el-tooltip content="清空会话" placement="top">
                  <el-button
                    :icon="Delete"
                    circle
                    text
                    type="primary"
                    :disabled="!messages.length"
                    @click="clearHistory"
                  />
                </el-tooltip>
                <el-tooltip content="关闭窗口" placement="top">
                  <el-button
                    :icon="CloseBold"
                    circle
                    text
                    type="primary"
                    @click="panelVisible = false"
                  />
                </el-tooltip>
              </div>
            </div>
          </template>

          <el-scrollbar ref="messageBoxRef" class="gva-ai-chat-content">
            <div v-if="!messages.length" class="gva-ai-chat-empty-state">
              <el-empty description="开始和 AI 对话吧" :image-size="72" />
              <div class="gva-ai-chat-presets">
                <el-button
                  v-for="prompt in presetPrompts"
                  :key="prompt"
                  class="gva-ai-chat-preset"
                  type="primary"
                  plain
                  @click="usePreset(prompt)"
                >
                  {{ prompt }}
                </el-button>
              </div>
            </div>

            <div
              v-for="item in messages"
              :key="item.id"
              class="gva-ai-chat-item"
              :class="`is-${item.role}`"
            >
              <div class="gva-ai-chat-item-meta">
                <el-avatar
                  :size="22"
                  class="gva-ai-chat-msg-avatar"
                  :class="item.role === 'user' ? 'is-user' : 'is-assistant'"
                >
                  {{ item.role === 'user' ? '你' : 'AI' }}
                </el-avatar>
                <span class="gva-ai-chat-time">{{ formatTime(item.timestamp) }}</span>
              </div>

              <div
                v-if="item.content"
                class="gva-ai-chat-bubble"
                :class="
                  item.role === 'user' ? 'is-user-bubble' : 'is-assistant-bubble'
                "
              >
                <div
                  v-if="item.role === 'assistant'"
                  class="gva-ai-chat-markdown"
                  v-html="renderAssistantContent(item.content)"
                ></div>
                <div v-else>{{ item.content }}</div>
              </div>

              <div
                v-if="item.images && item.images.length"
                class="gva-ai-chat-image-grid"
                :class="
                  item.role === 'user' ? 'is-user-images' : 'is-assistant-images'
                "
              >
                <div
                  v-for="(image, index) in item.images"
                  :key="`${item.id}-img-${index}`"
                  class="gva-ai-chat-image-item"
                >
                  <el-image
                    class="gva-ai-chat-image"
                    :src="image.url"
                    :preview-src-list="item.images.map((i) => i.url)"
                    :initial-index="index"
                    fit="cover"
                    preview-teleported
                  />
                  <div class="gva-ai-chat-image-name" :title="image.name || 'image'">
                    {{ image.name || 'image' }}
                  </div>
                </div>
              </div>

              <div
                v-if="item.files && item.files.length"
                class="gva-ai-chat-file-grid"
                :class="
                  item.role === 'user' ? 'is-user-files' : 'is-assistant-files'
                "
              >
                <div
                  v-for="(file, index) in item.files"
                  :key="`${item.id}-file-${index}`"
                  class="gva-ai-chat-file-item"
                >
                  <el-icon class="gva-ai-chat-msg-file-icon">
                    <Document />
                  </el-icon>
                  <div class="gva-ai-chat-file-meta">
                    <div class="gva-ai-chat-file-name" :title="file.name || 'file'">
                      {{ file.name || 'file' }}
                    </div>
                    <div class="gva-ai-chat-file-size">
                      {{ file.size ? formatFileSize(file.size) : file.type || getFileExt(file.name) }}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <el-alert
              v-if="isRequesting"
              class="gva-ai-chat-thinking"
              title="AI 正在思考，请稍候..."
              type="info"
              :closable="false"
              show-icon
            />
          </el-scrollbar>

          <footer class="gva-ai-chat-footer">
            <div v-if="pendingFiles.length" class="gva-ai-chat-upload-list">
              <div
                v-for="file in pendingFiles"
                :key="file.id"
                class="gva-ai-chat-upload-item"
              >
                <el-image
                  v-if="file.isImage"
                  class="gva-ai-chat-upload-image"
                  :src="file.url || file.base64"
                  fit="cover"
                  :preview-src-list="pendingImagePreviewList"
                  preview-teleported
                />
                <div v-else class="gva-ai-chat-file-placeholder">
                  <el-icon class="gva-ai-chat-file-icon">
                    <Document />
                  </el-icon>
                  <span>{{ getFileExt(file.name) }}</span>
                </div>
                <div class="gva-ai-chat-upload-meta">
                  <div class="gva-ai-chat-upload-name" :title="file.name">
                    {{ file.name }}
                  </div>
                  <div class="gva-ai-chat-upload-size">
                    {{ formatFileSize(file.size) }}
                  </div>
                </div>
                <el-button
                  class="gva-ai-chat-upload-remove"
                  :icon="Close"
                  circle
                  text
                  type="danger"
                  @click="removePendingFile(file.id)"
                />
              </div>
            </div>

            <el-input
              v-model="inputValue"
              class="gva-ai-chat-input"
              type="textarea"
              :autosize="{ minRows: 2, maxRows: 4 }"
              :maxlength="maxInputLength"
              show-word-limit
              resize="none"
              placeholder="输入问题，按 Enter 发送，Shift + Enter 换行。可拖拽图片到窗口"
              @keydown.enter.exact.prevent="handleSend"
            />

            <div class="gva-ai-chat-footer-actions">
              <div class="gva-ai-chat-footer-left">
                <input
                  ref="fileInputRef"
                  type="file"
                  multiple
                  class="gva-ai-chat-file-input"
                  @change="onFileInputChange"
                />
                <el-button size="small" :icon="Plus" @click="openFilePicker">
                  添加文件
                </el-button>
                <el-button
                  size="small"
                  :icon="Delete"
                  :disabled="!pendingFiles.length"
                  @click="clearPendingFiles"
                >
                  清空
                </el-button>
              </div>
              <div class="gva-ai-chat-footer-right">
                <el-button
                  :icon="VideoPause"
                  :disabled="!isRequesting"
                  @click="abortCurrentRequest"
                >
                  停止
                </el-button>
                <el-button
                  type="primary"
                  :icon="Promotion"
                  :disabled="sendDisabled"
                  @click="handleSend"
                >
                  发送
                </el-button>
              </div>
            </div>
          </footer>

          <transition name="gva-ai-chat-drop-mask">
            <div v-if="isDragOver" class="gva-ai-chat-drop-mask">
              <el-icon class="gva-ai-chat-drop-icon">
                <PictureFilled />
              </el-icon>
              <div class="gva-ai-chat-drop-text">释放鼠标以上传文件</div>
              <div class="gva-ai-chat-drop-subtext">
                支持任意文件，最多 8 个，单个 10MB；图片会直接预览
              </div>
            </div>
          </transition>
        </el-card>
      </section>
    </transition>

    <el-badge
      :is-dot="isRequesting"
      type="danger"
      class="gva-ai-chat-ball-wrap"
      :style="ballStyle"
    >
      <el-button
        class="gva-ai-chat-ball"
        circle
        type="primary"
        title="AI 助手"
        @pointerdown="onPointerDown"
        @click="togglePanel"
      >
        <el-icon class="gva-ai-chat-ball-icon">
          <ChatDotRound />
        </el-icon>
      </el-button>
      <div class="gva-ai-chat-ball-text">AI助手</div>
    </el-badge>
  </div>
</template>

<script setup>
  import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
  import { marked } from 'marked'
  import { ElMessage } from 'element-plus'
  import { useRoute } from 'vue-router'
  import {
    ChatDotRound,
    Close,
    CloseBold,
    Delete,
    Document,
    PictureFilled,
    Plus,
    Promotion,
    Service,
    VideoPause
  } from '@element-plus/icons-vue'
  import { useUserStore } from '@/pinia/modules/user'
  import { getBaseUrl } from '@/utils/format'

  defineOptions({
    name: 'AiSuspendedBallChat'
  })

  const props = defineProps({
    title: {
      type: String,
      default: 'AI 助手'
    },
    subtitle: {
      type: String,
      default: '悬浮球聊天插件'
    },
    maxHistoryCount: {
      type: Number,
      default: 60
    },
    maxInputLength: {
      type: Number,
      default: 2000
    },
    enableStreaming: {
      type: Boolean,
      default: true
    },
    enableContext: {
      type: Boolean,
      default: true
    },
    contextHistoryCount: {
      type: Number,
      default: 12
    },
    requestTimeout: {
      type: Number,
      default: 60000
    }
  })

  const BALL_SIZE = 58
  const MAX_FILE_COUNT = 8
  const MAX_FILE_SIZE = 10 * 1024 * 1024

  const userStore = useUserStore()
  const route = useRoute()

  const historyStorageKey = 'gva-ai-chat-history-v1'
  const ballStorageKey = 'gva-ai-chat-ball-position-v1'
  const presetPrompts = [
    '帮我总结今天的工作重点',
    '帮我排查一个前端报错',
    '给我一个 SQL 优化建议'
  ]

  const panelVisible = ref(false)
  const inputValue = ref('')
  const isRequesting = ref(false)
  const isDragOver = ref(false)
  const messageBoxRef = ref(null)
  const fileInputRef = ref(null)

  const messages = ref([])
  const pendingFiles = ref([])

  const ballPosition = ref(getDefaultBallPosition())
  const viewport = ref({
    width: window.innerWidth,
    height: window.innerHeight
  })

  const dragState = {
    dragging: false,
    moved: false,
    offsetX: 0,
    offsetY: 0
  }

  const dragDropState = {
    counter: 0
  }

  let suppressClick = false
  let abortController = null

  marked.setOptions({
    gfm: true,
    breaks: true
  })

  const visible = computed(() => {
    return !!userStore.token && route.name !== 'Login'
  })

  const sendDisabled = computed(() => {
    return (
      (!inputValue.value.trim() && pendingFiles.value.length === 0) ||
      isRequesting.value
    )
  })

  const pendingImagePreviewList = computed(() => {
    return pendingFiles.value
      .filter((item) => item.isImage)
      .map((item) => item.url || item.base64)
      .filter(Boolean)
  })

  const chatModeLabel = computed(() => {
    return props.enableStreaming ? '流式模式' : '标准模式'
  })

  const contextModeLabel = computed(() => {
    return props.enableContext ? '携带上下文' : '单轮问答'
  })

  const ballStyle = computed(() => {
    return {
      left: `${ballPosition.value.x}px`,
      top: `${ballPosition.value.y}px`
    }
  })

  const panelStyle = computed(() => {
    const panelWidth = viewport.value.width <= 768 ? viewport.value.width - 24 : 440
    const panelHeight = viewport.value.height <= 720 ? viewport.value.height - 120 : 620
    const offset = 14
    const placeLeft = ballPosition.value.x > viewport.value.width / 2

    let left = placeLeft
      ? ballPosition.value.x - panelWidth - offset
      : ballPosition.value.x + BALL_SIZE + offset
    let top = ballPosition.value.y - panelHeight + BALL_SIZE

    left = clamp(left, 12, viewport.value.width - panelWidth - 12)
    top = clamp(top, 64, viewport.value.height - panelHeight - 24)

    return {
      width: `${panelWidth}px`,
      height: `${panelHeight}px`,
      left: `${left}px`,
      top: `${top}px`
    }
  })

  watch(
    messages,
    () => {
      saveHistory()
      scrollToBottom()
    },
    { deep: true }
  )

  watch(visible, (val) => {
    if (!val) {
      panelVisible.value = false
      clearPendingFiles()
    }
  })

  onMounted(() => {
    restoreHistory()
    restoreBallPosition()
    window.addEventListener('resize', handleResize)
  })

  onBeforeUnmount(() => {
    window.removeEventListener('resize', handleResize)
    removePointerEvents()
    abortCurrentRequest()
    clearPendingFiles()
  })

  function getDefaultBallPosition() {
    return {
      x: Math.max(window.innerWidth - 86, 12),
      y: Math.max(window.innerHeight - 148, 80)
    }
  }

  function clamp(value, min, max) {
    return Math.min(Math.max(value, min), max)
  }

  function normalizeBallPosition(pos) {
    const maxX = Math.max(viewport.value.width - BALL_SIZE, 12)
    const maxY = Math.max(viewport.value.height - BALL_SIZE, 72)
    return {
      x: clamp(pos.x, 12, maxX),
      y: clamp(pos.y, 72, maxY)
    }
  }

  function handleResize() {
    viewport.value = {
      width: window.innerWidth,
      height: window.innerHeight
    }
    ballPosition.value = normalizeBallPosition(ballPosition.value)
    saveBallPosition()
  }

  function onPointerDown(event) {
    dragState.dragging = true
    dragState.moved = false
    dragState.offsetX = event.clientX - ballPosition.value.x
    dragState.offsetY = event.clientY - ballPosition.value.y
    window.addEventListener('pointermove', onPointerMove)
    window.addEventListener('pointerup', onPointerUp)
  }

  function onPointerMove(event) {
    if (!dragState.dragging) return
    dragState.moved = true
    ballPosition.value = normalizeBallPosition({
      x: event.clientX - dragState.offsetX,
      y: event.clientY - dragState.offsetY
    })
  }

  function onPointerUp() {
    if (!dragState.dragging) return
    dragState.dragging = false
    suppressClick = dragState.moved
    saveBallPosition()
    removePointerEvents()
    setTimeout(() => {
      suppressClick = false
    }, 120)
  }

  function removePointerEvents() {
    window.removeEventListener('pointermove', onPointerMove)
    window.removeEventListener('pointerup', onPointerUp)
  }

  function togglePanel() {
    if (suppressClick) return
    panelVisible.value = !panelVisible.value
    if (panelVisible.value) {
      scrollToBottom()
    }
  }

  function buildMessage(role, content, images = [], files = []) {
    return {
      id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
      role,
      content,
      images,
      files,
      timestamp: Date.now()
    }
  }

  function sanitizeImagesForStorage(images = []) {
    return images.filter((item) => {
      const url = item?.url || ''
      return /^https?:\/\//i.test(url)
    })
  }

  function sanitizeFilesForStorage(files = []) {
    if (!Array.isArray(files)) return []
    return files.map((item) => {
      return {
        name: item?.name || 'file',
        size: Number(item?.size) || 0,
        type: item?.type || '',
        isImage: false
      }
    })
  }

  function saveHistory() {
    const latest = messages.value.slice(-props.maxHistoryCount).map((item) => {
      return {
        ...item,
        images: sanitizeImagesForStorage(item.images),
        files: sanitizeFilesForStorage(item.files)
      }
    })
    localStorage.setItem(historyStorageKey, JSON.stringify(latest))
  }

  function restoreHistory() {
    try {
      const raw = localStorage.getItem(historyStorageKey)
      if (!raw) return
      const parsed = JSON.parse(raw)
      if (Array.isArray(parsed)) {
        messages.value = parsed.slice(-props.maxHistoryCount).map((item) => {
          return {
            ...item,
            images: Array.isArray(item.images) ? item.images : [],
            files: Array.isArray(item.files) ? item.files : []
          }
        })
      }
    } catch (err) {
      console.warn('Failed to restore ai chat history', err)
    }
  }

  function clearHistory() {
    messages.value = []
    localStorage.removeItem(historyStorageKey)
  }

  async function usePreset(prompt) {
    inputValue.value = prompt
    await nextTick()
    handleSend()
  }

  function saveBallPosition() {
    localStorage.setItem(ballStorageKey, JSON.stringify(ballPosition.value))
  }

  function restoreBallPosition() {
    try {
      const raw = localStorage.getItem(ballStorageKey)
      if (!raw) return
      const parsed = JSON.parse(raw)
      if (
        parsed &&
        typeof parsed.x === 'number' &&
        typeof parsed.y === 'number'
      ) {
        ballPosition.value = normalizeBallPosition(parsed)
      }
    } catch (err) {
      console.warn('Failed to restore ai chat ball position', err)
    }
  }

  function formatTime(timestamp) {
    try {
      return new Date(timestamp).toLocaleTimeString([], {
        hour: '2-digit',
        minute: '2-digit'
      })
    } catch (err) {
      return ''
    }
  }

  function sanitizeMarkdown(raw) {
    return String(raw || '').replaceAll(/<script[\s\S]*?>[\s\S]*?<\/script>/gi, '')
  }

  function sanitizeRenderedHtml(raw) {
    return String(raw || '')
      .replaceAll(/\son\w+="[^"]*"/gi, '')
      .replaceAll(/\son\w+='[^']*'/gi, '')
      .replaceAll(
        /\s(href|src)=("|')\s*javascript:[\s\S]*?\2/gi,
        ' $1="#"'
      )
      .replaceAll(
        /\s(href|src)=("|')\s*data:text\/html[\s\S]*?\2/gi,
        ' $1="#"'
      )
  }

  function renderAssistantContent(content) {
    const safeMarkdown = sanitizeMarkdown(content)
    const rendered = marked.parse(safeMarkdown)
    return sanitizeRenderedHtml(rendered)
  }

  function scrollToBottom() {
    nextTick(() => {
      const scrollbar = messageBoxRef.value
      if (!scrollbar || !scrollbar.wrapRef) return
      scrollbar.setScrollTop(scrollbar.wrapRef.scrollHeight)
    })
  }

  function resolveApiUrl() {
    const configured = import.meta.env.VITE_AI_CHAT_API || '/ai/chat'
    if (/^https?:\/\//i.test(configured)) {
      return configured
    }
    const prefix = getBaseUrl() || ''
    const path = configured.startsWith('/') ? configured : `/${configured}`
    return `${prefix}${path}`
  }

  function abortCurrentRequest() {
    if (!abortController) return
    abortController.abort()
    abortController = null
    isRequesting.value = false
  }

  function hasTransferFiles(dataTransfer) {
    if (!dataTransfer || !dataTransfer.items) return false
    return Array.from(dataTransfer.items).some((item) => item.kind === 'file')
  }

  function onDragEnter(event) {
    if (!hasTransferFiles(event.dataTransfer)) return
    dragDropState.counter += 1
    isDragOver.value = true
  }

  function onDragOver(event) {
    if (!hasTransferFiles(event.dataTransfer)) return
    event.dataTransfer.dropEffect = 'copy'
    isDragOver.value = true
  }

  function onDragLeave(event) {
    if (!hasTransferFiles(event.dataTransfer)) return
    dragDropState.counter -= 1
    if (dragDropState.counter <= 0) {
      dragDropState.counter = 0
      isDragOver.value = false
    }
  }

  async function onDrop(event) {
    dragDropState.counter = 0
    isDragOver.value = false
    const files = Array.from(event.dataTransfer?.files || [])
    await appendPendingFiles(files)
  }

  function openFilePicker() {
    fileInputRef.value?.click()
  }

  async function onFileInputChange(event) {
    const files = Array.from(event.target.files || [])
    await appendPendingFiles(files)
    event.target.value = ''
  }

  async function appendPendingFiles(files) {
    const validFiles = files.filter((file) => file instanceof File)
    if (!validFiles.length) return

    const remain = MAX_FILE_COUNT - pendingFiles.value.length
    if (remain <= 0) {
      ElMessage.warning(`最多上传 ${MAX_FILE_COUNT} 个文件`)
      return
    }

    const accepted = validFiles.slice(0, remain)
    for (const file of accepted) {
      if (file.size > MAX_FILE_SIZE) {
        ElMessage.warning(`${file.name} 超过 10MB，已跳过`)
        continue
      }
      try {
        const base64 = await fileToDataURL(file)
        const isImage = file.type.startsWith('image/')
        const url = isImage ? URL.createObjectURL(file) : ''
        pendingFiles.value.push({
          id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
          name: file.name,
          size: file.size,
          type: file.type,
          base64,
          isImage,
          url
        })
      } catch (err) {
        ElMessage.error(`读取文件失败: ${file.name}`)
      }
    }

    if (validFiles.length > remain) {
      ElMessage.warning(`最多上传 ${MAX_FILE_COUNT} 个文件，多余文件已忽略`)
    }
  }

  function fileToDataURL(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(String(reader.result || ''))
      reader.onerror = () => reject(new Error('read file failed'))
      reader.readAsDataURL(file)
    })
  }

  function removePendingFile(id) {
    const index = pendingFiles.value.findIndex((item) => item.id === id)
    if (index === -1) return
    const target = pendingFiles.value[index]
    if (target.url?.startsWith('blob:')) {
      URL.revokeObjectURL(target.url)
    }
    pendingFiles.value.splice(index, 1)
  }

  function clearPendingFiles() {
    pendingFiles.value.forEach((item) => {
      if (item.url?.startsWith('blob:')) {
        URL.revokeObjectURL(item.url)
      }
    })
    pendingFiles.value = []
  }

  function getFileExt(name = '') {
    const ext = String(name).split('.').pop() || ''
    return ext ? ext.toUpperCase() : 'FILE'
  }

  function formatFileSize(size = 0) {
    if (size >= 1024 * 1024) {
      return `${(size / (1024 * 1024)).toFixed(1)} MB`
    }
    if (size >= 1024) {
      return `${Math.round(size / 1024)} KB`
    }
    return `${size} B`
  }

  function normalizeImageObject(item) {
    if (typeof item === 'string') {
      return { url: item, name: 'image' }
    }
    if (!item || typeof item !== 'object') return null

    const url = item.url || item.src || item.link || item.base64 || ''
    if (!url) return null

    return {
      url,
      name: item.name || item.fileName || 'image'
    }
  }

  function normalizeImageList(raw) {
    if (!Array.isArray(raw)) return []
    const seen = new Set()
    const images = []
    raw.forEach((item) => {
      const normalized = normalizeImageObject(item)
      if (!normalized) return
      if (seen.has(normalized.url)) return
      seen.add(normalized.url)
      images.push(normalized)
    })
    return images
  }

  function mergeImageList(origin, incoming) {
    return normalizeImageList([...(origin || []), ...(incoming || [])])
  }

  function normalizeFileObject(item) {
    if (typeof item === 'string') {
      return {
        name: item,
        size: 0,
        type: '',
        isImage: false
      }
    }

    if (!item || typeof item !== 'object') return null

    const name =
      item.name || item.fileName || item.filename || item.title || item.url || 'file'
    const type = item.type || item.mime || ''
    const size = Number(item.size) || 0
    const isImage = Boolean(item.isImage || String(type).startsWith('image/'))

    return {
      name,
      size,
      type,
      isImage
    }
  }

  function normalizeFileList(raw) {
    if (!Array.isArray(raw)) return []
    const seen = new Set()
    const files = []

    raw.forEach((item) => {
      const normalized = normalizeFileObject(item)
      if (!normalized) return
      const key = `${normalized.name}|${normalized.size}|${normalized.type}`
      if (seen.has(key)) return
      seen.add(key)
      files.push(normalized)
    })

    return files
  }

  function mergeFileList(origin, incoming) {
    return normalizeFileList([...(origin || []), ...(incoming || [])])
  }

  function extractContentText(raw) {
    if (typeof raw === 'string') return raw
    if (!raw) return ''

    if (Array.isArray(raw)) {
      return raw
        .map((item) => {
          if (typeof item === 'string') return item
          if (!item || typeof item !== 'object') return ''
          if (typeof item.text === 'string') return item.text
          if (typeof item.content === 'string') return item.content
          if (typeof item.delta === 'string') return item.delta
          if (typeof item.value === 'string') return item.value
          if (typeof item?.text?.value === 'string') return item.text.value
          return ''
        })
        .join('')
    }

    if (typeof raw === 'object') {
      if (typeof raw.text === 'string') return raw.text
      if (typeof raw.content === 'string') return raw.content
      if (typeof raw.delta === 'string') return raw.delta
      if (typeof raw.value === 'string') return raw.value
      if (typeof raw?.text?.value === 'string') return raw.text.value
    }

    return ''
  }

  function extractImagesFromContentNodes(raw) {
    if (!Array.isArray(raw)) return []
    return normalizeImageList(
      raw
        .map((item) => {
          if (!item || typeof item !== 'object') return null
          if (item.type === 'image_url') {
            return item.image_url?.url || item.url || null
          }
          if (item.type === 'input_image') {
            return item.image_url || item.url || null
          }
          if (item.type === 'image') {
            return item.url || item.src || null
          }
          return null
        })
        .filter(Boolean)
    )
  }

  function mergeStreamText(current = '', incoming = '') {
    if (!incoming) return current
    if (!current) return incoming
    if (incoming.startsWith(current)) return incoming
    if (current.endsWith(incoming)) return current
    return `${current}${incoming}`
  }

  function extractResponsePayload(raw) {
    if (typeof raw === 'string') {
      return {
        content: raw,
        images: [],
        files: []
      }
    }

    if (!raw || typeof raw !== 'object') {
      return {
        content: '',
        images: [],
        files: []
      }
    }

    const contentCandidates = [
      raw.result?.answer,
      raw.result?.content,
      raw.result?.text,
      raw.data?.answer,
      raw.data?.result,
      raw.data?.content,
      raw.answer,
      raw.result,
      raw.content,
      raw.message,
      raw.text,
      raw.delta,
      raw.output_text,
      raw.response?.output_text,
      raw.choices?.[0]?.delta?.content,
      raw.choices?.[0]?.delta,
      raw.choices?.[0]?.message?.content
    ]

    let content = ''
    for (const candidate of contentCandidates) {
      const text = extractContentText(candidate)
      if (text) {
        content = text
        break
      }
    }

    const imageCandidates = [
      raw.images,
      raw.result?.images,
      raw.data?.images,
      raw.result?.imageUrls,
      raw.result?.image_urls,
      raw.data?.imageUrls,
      raw.data?.image_urls,
      raw.output?.images,
      raw.data?.output?.images,
      raw.result?.output?.images,
      extractImagesFromContentNodes(raw.choices?.[0]?.message?.content),
      extractImagesFromContentNodes(raw.choices?.[0]?.delta?.content)
    ]

    let images = []
    for (const candidate of imageCandidates) {
      const normalized = normalizeImageList(candidate)
      if (normalized.length) {
        images = normalized
        break
      }
    }

    const fileCandidates = [
      raw.files,
      raw.attachments,
      raw.result?.files,
      raw.result?.attachments,
      raw.data?.files,
      raw.data?.attachments
    ]

    let files = []
    for (const candidate of fileCandidates) {
      const normalized = normalizeFileList(candidate).filter((item) => !item.isImage)
      if (normalized.length) {
        files = normalized
        break
      }
    }

    return {
      content,
      images,
      files
    }
  }

  async function handleSend() {
    const questionRaw = inputValue.value.trim()
    const hasOnlyFile = !questionRaw && pendingFiles.value.length > 0
    const question = questionRaw || (hasOnlyFile ? '请分析这些文件' : '')
    if (!question || isRequesting.value) return

    panelVisible.value = true

    const outgoingFiles = pendingFiles.value.map((item) => {
      return {
        name: item.name,
        type: item.type,
        size: item.size,
        base64: item.base64,
        isImage: item.isImage
      }
    })

    const outgoingImages = outgoingFiles.filter((item) => item.isImage)

    const userDisplayImages = pendingFiles.value
      .filter((item) => item.isImage)
      .map((item) => {
        return {
          name: item.name,
          url: item.base64
        }
      })

    const userDisplayFiles = pendingFiles.value
      .filter((item) => !item.isImage)
      .map((item) => {
        return {
          name: item.name,
          size: item.size,
          type: item.type,
          isImage: false
        }
      })

    inputValue.value = ''
    clearPendingFiles()

    messages.value.push(
      buildMessage('user', question, userDisplayImages, userDisplayFiles)
    )
    const assistantMessage = buildMessage('assistant', '', [], [])
    messages.value.push(assistantMessage)

    const recentHistory = props.enableContext
      ? messages.value
          .filter((item) => item.role === 'user' || item.role === 'assistant')
          .map((item) => ({
            role: item.role,
            content: item.content
          }))
          .slice(-props.contextHistoryCount)
      : []

    const payload = {
      message: question,
      question,
      prompt: question,
      appName: 'iadmin-web',
      domainName:
        userStore.userInfo?.userName || userStore.userInfo?.nickName || 'user',
      stream: props.enableStreaming,
      messages: recentHistory,
      images: outgoingImages,
      imageList: outgoingImages,
      attachments: outgoingFiles,
      files: outgoingFiles,
      hasImage: outgoingImages.length > 0,
      hasFile: outgoingFiles.length > 0
    }

    isRequesting.value = true
    abortController = new AbortController()

    try {
      const headers = {
        'Content-Type': 'application/json',
        Accept: props.enableStreaming
          ? 'text/event-stream, application/json'
          : 'application/json',
        'x-token': userStore.token || '',
        'x-user-id': userStore.userInfo?.ID || ''
      }

      const requestPromise = fetch(resolveApiUrl(), {
        method: 'POST',
        headers,
        body: JSON.stringify(payload),
        signal: abortController.signal
      })

      const response = await withTimeout(requestPromise, props.requestTimeout)

      if (!response.ok) {
        const errorText = await safeReadText(response)
        throw new Error(errorText || `请求失败(${response.status})`)
      }

      const contentType = (response.headers.get('content-type') || '').toLowerCase()
      if (
        props.enableStreaming &&
        response.body &&
        contentType.includes('text/event-stream')
      ) {
        await consumeStreamResponse(response, assistantMessage)
      } else {
        const data = await parseResponseBody(response)
        const parsed = extractResponsePayload(data)
        if (!parsed.content && !parsed.images.length && !parsed.files.length) {
          throw new Error('接口返回为空，请检查后端响应结构')
        }
        assistantMessage.content = parsed.content
        assistantMessage.images = parsed.images
        assistantMessage.files = parsed.files
      }
    } catch (err) {
      const aborted = err?.name === 'AbortError'
      if (aborted) {
        assistantMessage.content = assistantMessage.content || '已停止本次请求。'
      } else {
        assistantMessage.content = `请求失败：${err.message || '未知错误'}`
      }
    } finally {
      isRequesting.value = false
      abortController = null
      scrollToBottom()
    }
  }

  async function consumeStreamResponse(response, assistantMessage) {
    const reader = response.body.getReader()
    const decoder = new TextDecoder('utf-8')
    let buffer = ''
    let hasContent = false
    let ended = false

    const processEvent = (eventRaw) => {
      const lines = eventRaw.split(/\r?\n/)
      const dataParts = []

      for (const rawLine of lines) {
        if (!rawLine) continue
        if (rawLine.startsWith(':')) continue

        if (rawLine.startsWith('data:')) {
          let data = rawLine.slice(5)
          if (data.startsWith(' ')) data = data.slice(1)
          dataParts.push(data.replace(/\r$/, ''))
          continue
        }

        dataParts.push(rawLine.replace(/\r$/, ''))
      }

      const dataLine = dataParts.join('\n')
      if (!dataLine) return false
      if (dataLine === '[DONE]') return true

      let parsed = dataLine
      const trimStartLine = dataLine.trimStart()
      if (trimStartLine.startsWith('{') || trimStartLine.startsWith('[')) {
        try {
          parsed = JSON.parse(trimStartLine)
        } catch (err) {
          parsed = dataLine
        }
      }

      const streamPayload = extractResponsePayload(parsed)
      if (streamPayload.content) {
        assistantMessage.content = mergeStreamText(
          assistantMessage.content,
          streamPayload.content
        )
        hasContent = true
      }

      if (streamPayload.images.length) {
        assistantMessage.images = mergeImageList(
          assistantMessage.images,
          streamPayload.images
        )
        hasContent = true
      }

      if (streamPayload.files.length) {
        assistantMessage.files = mergeFileList(
          assistantMessage.files,
          streamPayload.files
        )
        hasContent = true
      }

      return isStreamEnd(parsed)
    }

    while (!ended) {
      const { value, done } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const events = buffer.split(/\r?\n\r?\n/)
      buffer = events.pop() || ''

      for (const eventRaw of events) {
        if (!eventRaw.trim()) continue
        ended = processEvent(eventRaw)
        if (ended) break
      }

      scrollToBottom()
    }

    if (!ended && buffer.trim()) {
      ended = processEvent(buffer)
      if (ended) {
        scrollToBottom()
      }
    }

    if (!hasContent) {
      throw new Error('stream response has no usable content')
    }
  }

  function isStreamEnd(data) {
    if (!data || typeof data !== 'object') return false
    return data.is_end === true || data.done === true || data.finish === true
  }

  async function parseResponseBody(response) {
    const contentType = (response.headers.get('content-type') || '').toLowerCase()
    if (contentType.includes('application/json')) {
      return response.json()
    }
    return response.text()
  }

  async function safeReadText(response) {
    try {
      const text = await response.text()
      if (!text) return ''
      if (text.startsWith('{')) {
        const parsed = JSON.parse(text)
        return parsed.msg || parsed.message || text
      }
      return text
    } catch (err) {
      return ''
    }
  }

  async function withTimeout(promise, timeout) {
    let timer = null
    try {
      return await Promise.race([
        promise,
        new Promise((_, reject) => {
          timer = setTimeout(() => {
            reject(new Error(`请求超时(${timeout}ms)`))
          }, timeout)
        })
      ])
    } finally {
      if (timer) clearTimeout(timer)
    }
  }
</script>

<style scoped lang="scss">
  .gva-ai-chat-root {
    position: fixed;
    inset: 0;
    pointer-events: none;
    z-index: 2999;
    --gva-ai-radius-panel: calc(var(--el-border-radius-base) + 12px);
    --gva-ai-radius-bubble: calc(var(--el-border-radius-base) + 10px);
    --gva-ai-radius-tail: calc(var(--el-border-radius-base) + 2px);
    --gva-ai-radius-card: calc(var(--el-border-radius-base) + 6px);
    --gva-ai-layer-0: var(--el-bg-color-page);
    --gva-ai-layer-1: var(--el-bg-color);
    --gva-ai-layer-2: var(--el-bg-color-overlay);
    --gva-ai-shadow-panel: 0 24px 60px
      color-mix(in srgb, var(--el-text-color-primary) 20%, transparent);
    --gva-ai-shadow-soft: 0 8px 18px
      color-mix(in srgb, var(--el-text-color-primary) 12%, transparent);
  }

  :global(html.dark) .gva-ai-chat-root {
    --gva-ai-layer-0: color-mix(in srgb, var(--el-bg-color-page) 86%, #000 14%);
    --gva-ai-layer-1: color-mix(in srgb, var(--el-bg-color) 90%, #000 10%);
    --gva-ai-layer-2: color-mix(in srgb, var(--el-bg-color-overlay) 86%, #000 14%);
    --gva-ai-shadow-panel: 0 20px 50px rgb(0 0 0 / 45%);
    --gva-ai-shadow-soft: 0 8px 18px rgb(0 0 0 / 30%);
  }

  .gva-ai-chat-panel {
    position: fixed;
    pointer-events: auto;
  }

  .gva-ai-chat-panel-card {
    position: relative;
    height: 100%;
    border-radius: var(--gva-ai-radius-panel);
    border: 1px solid var(--el-border-color-light);
    box-shadow: var(--gva-ai-shadow-panel);
    overflow: hidden;
  }

  .gva-ai-chat-panel-card :deep(.el-card__header) {
    padding: 0;
    border-bottom: 0;
  }

  .gva-ai-chat-panel-card :deep(.el-card__body) {
    height: calc(100% - 98px);
    padding: 0;
    display: flex;
    flex-direction: column;
    background: var(--gva-ai-layer-0);
  }

  .gva-ai-chat-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 10px;
    padding: 14px 14px 12px;
    background:
      radial-gradient(circle at 8% 10%, rgb(255 255 255 / 55%), transparent 40%),
      linear-gradient(
        130deg,
        color-mix(in srgb, var(--el-color-primary-light-9) 85%, #fff 15%),
        var(--el-color-primary-light-8) 45%,
        color-mix(in srgb, var(--el-color-primary-light-9) 70%, var(--el-color-success-light-9) 30%)
      );
    border-bottom: 1px solid var(--el-border-color-lighter);
  }

  :global(html.dark) .gva-ai-chat-header {
    background:
      radial-gradient(circle at 10% 10%, rgb(255 255 255 / 8%), transparent 36%),
      linear-gradient(
        130deg,
        color-mix(in srgb, var(--el-color-primary) 24%, #111827 76%),
        color-mix(in srgb, var(--el-color-primary-light-3) 16%, #0f172a 84%)
      );
    border-bottom-color: var(--el-border-color);
  }

  .gva-ai-chat-head-left {
    min-width: 0;
    display: flex;
    align-items: flex-start;
    gap: 10px;
  }

  .gva-ai-chat-head-avatar {
    flex-shrink: 0;
    color: #fff;
    background: linear-gradient(
      140deg,
      var(--el-color-primary),
      var(--el-color-success)
    );
    box-shadow: 0 6px 16px rgb(37 99 235 / 28%);
  }

  .gva-ai-chat-head-text {
    min-width: 0;
  }

  .gva-ai-chat-title {
    color: var(--el-text-color-primary);
    font-size: 14px;
    font-weight: 700;
    line-height: 1.2;
  }

  .gva-ai-chat-subtitle {
    margin-top: 2px;
    color: var(--el-text-color-secondary);
    font-size: 12px;
    line-height: 1.2;
  }

  .gva-ai-chat-mode-tags {
    margin-top: 8px;
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: wrap;
  }

  .gva-ai-chat-actions {
    display: flex;
    align-items: center;
    gap: 4px;
    flex-shrink: 0;
  }

  .gva-ai-chat-content {
    flex: 1;
    padding: 12px;
  }

  .gva-ai-chat-content :deep(.el-scrollbar__wrap) {
    padding-right: 2px;
  }

  .gva-ai-chat-empty-state {
    min-height: 100%;
    padding: 18px 8px 8px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }

  .gva-ai-chat-empty-state :deep(.el-empty) {
    padding: 0;
  }

  .gva-ai-chat-item {
    margin-bottom: 12px;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .gva-ai-chat-item.is-user {
    align-items: flex-end;
  }

  .gva-ai-chat-item.is-assistant {
    align-items: flex-start;
  }

  .gva-ai-chat-item-meta {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 0 2px;
  }

  .gva-ai-chat-msg-avatar.is-user {
    color: #fff;
    background: var(--el-color-primary);
  }

  .gva-ai-chat-msg-avatar.is-assistant {
    color: #fff;
    background: var(--el-color-success);
  }

  .gva-ai-chat-time {
    color: var(--el-text-color-secondary);
    font-size: 11px;
  }

  .gva-ai-chat-bubble {
    max-width: calc(100% - 26px);
    padding: 10px 12px;
    white-space: pre-wrap;
    word-break: break-word;
    font-size: 13px;
    line-height: 1.6;
    border-width: 1px;
    border-style: solid;
    border-radius: var(--gva-ai-radius-bubble);
  }

  .is-user-bubble {
    color: #fff;
    background: linear-gradient(
      135deg,
      var(--el-color-primary),
      var(--el-color-primary-dark-2)
    );
    border-color: color-mix(in srgb, var(--el-color-primary) 80%, #000 20%);
    border-radius: var(--gva-ai-radius-bubble) var(--gva-ai-radius-bubble)
      var(--gva-ai-radius-tail) var(--gva-ai-radius-bubble);
    box-shadow: var(--gva-ai-shadow-soft);
  }

  .is-assistant-bubble {
    color: var(--el-text-color-primary);
    background: var(--gva-ai-layer-2);
    border-color: var(--el-border-color-lighter);
    border-radius: var(--gva-ai-radius-bubble) var(--gva-ai-radius-bubble)
      var(--gva-ai-radius-bubble) var(--gva-ai-radius-tail);
    box-shadow: 0 2px 10px rgb(15 23 42 / 6%);
  }

  :global(html.dark) .is-user-bubble {
    background: linear-gradient(
      135deg,
      color-mix(in srgb, var(--el-color-primary) 86%, #0f172a 14%),
      color-mix(in srgb, var(--el-color-primary-dark-2) 75%, #0b1120 25%)
    );
    border-color: color-mix(in srgb, var(--el-color-primary-light-3) 36%, #000 64%);
    box-shadow: 0 10px 20px color-mix(in srgb, var(--el-color-primary) 32%, #000 68%);
  }

  :global(html.dark) .is-assistant-bubble {
    background: var(--gva-ai-layer-2);
    border-color: var(--el-border-color);
    box-shadow: 0 2px 10px rgb(0 0 0 / 26%);
  }

  .gva-ai-chat-markdown :deep(p) {
    margin: 0 0 0.45em;
  }

  .gva-ai-chat-markdown :deep(p:last-child) {
    margin-bottom: 0;
  }

  .gva-ai-chat-markdown :deep(pre) {
    margin: 8px 0;
    padding: 8px 10px;
    overflow: auto;
    border-radius: var(--gva-ai-radius-card);
    background: var(--el-fill-color-light);
    border: 1px solid var(--el-border-color-lighter);
  }

  .gva-ai-chat-markdown :deep(code) {
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas,
      'Liberation Mono', 'Courier New', monospace;
  }

  .gva-ai-chat-markdown :deep(a) {
    color: var(--el-color-primary);
    text-decoration: none;
  }

  .gva-ai-chat-markdown :deep(a:hover) {
    text-decoration: underline;
  }

  .gva-ai-chat-markdown :deep(img) {
    max-width: 100%;
    border-radius: var(--gva-ai-radius-card);
    border: 1px solid var(--el-border-color-lighter);
    display: block;
    margin-top: 6px;
  }

  .gva-ai-chat-image-grid {
    max-width: calc(100% - 26px);
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 8px;
    width: fit-content;
  }

  .is-user-images {
    align-self: flex-end;
    justify-content: end;
    justify-items: end;
  }

  .is-assistant-images {
    align-self: flex-start;
    justify-content: start;
    justify-items: start;
  }

  .gva-ai-chat-image-item {
    width: 142px;
    background: var(--gva-ai-layer-2);
    border-radius: var(--gva-ai-radius-card);
    border: 1px solid var(--el-border-color-lighter);
    overflow: hidden;
  }

  .gva-ai-chat-image {
    width: 100%;
    height: 92px;
  }

  .gva-ai-chat-image-name {
    padding: 6px 8px;
    font-size: 11px;
    color: var(--el-text-color-secondary);
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
  }

  .gva-ai-chat-file-grid {
    max-width: calc(100% - 26px);
    display: flex;
    flex-direction: column;
    gap: 6px;
    width: min(260px, 100%);
  }

  .is-user-files {
    align-self: flex-end;
  }

  .is-assistant-files {
    align-self: flex-start;
  }

  .gva-ai-chat-file-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 10px;
    border-radius: var(--gva-ai-radius-card);
    border: 1px solid var(--el-border-color-lighter);
    background: var(--gva-ai-layer-2);
  }

  .gva-ai-chat-msg-file-icon {
    flex-shrink: 0;
    font-size: 16px;
    color: var(--el-color-primary);
  }

  .gva-ai-chat-file-meta {
    min-width: 0;
    flex: 1;
  }

  .gva-ai-chat-file-name {
    font-size: 12px;
    color: var(--el-text-color-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .gva-ai-chat-file-size {
    margin-top: 2px;
    font-size: 11px;
    color: var(--el-text-color-secondary);
  }

  .gva-ai-chat-thinking {
    margin-top: 10px;
  }

  .gva-ai-chat-presets {
    width: 100%;
    padding: 2px 20px 8px;
    display: flex;
    justify-content: center;
    gap: 8px;
    flex-wrap: wrap;
  }

  .gva-ai-chat-preset {
    border-radius: 999px;
    border-color: color-mix(in srgb, var(--el-color-primary) 20%, var(--el-border-color));
    background: color-mix(in srgb, var(--el-color-primary-light-9) 88%, var(--el-bg-color) 12%);
    color: var(--el-color-primary);
    box-shadow: 0 6px 14px color-mix(in srgb, var(--el-color-primary) 14%, transparent);
    transition:
      transform 0.2s ease,
      box-shadow 0.2s ease;
  }

  .gva-ai-chat-preset:hover {
    transform: translateY(-1px);
    box-shadow: 0 10px 16px color-mix(in srgb, var(--el-color-primary) 22%, transparent);
  }

  :global(html.dark) .gva-ai-chat-preset {
    background: color-mix(in srgb, var(--el-color-primary) 16%, #111827 84%);
    border-color: color-mix(in srgb, var(--el-color-primary-light-3) 26%, #1f2937 74%);
    color: var(--el-color-primary-light-3);
  }

  .gva-ai-chat-footer {
    border-top: 1px solid var(--el-border-color-lighter);
    background: var(--gva-ai-layer-1);
    padding: 12px;
  }

  .gva-ai-chat-upload-list {
    margin-bottom: 10px;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(110px, 1fr));
    gap: 8px;
  }

  .gva-ai-chat-upload-item {
    position: relative;
    border-radius: var(--gva-ai-radius-card);
    overflow: hidden;
    border: 1px solid var(--el-border-color-lighter);
    background: var(--el-fill-color-blank);
  }

  .gva-ai-chat-upload-image {
    width: 100%;
    height: 80px;
  }

  .gva-ai-chat-file-placeholder {
    height: 80px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 6px;
    color: var(--el-color-primary);
    background: color-mix(in srgb, var(--el-color-primary-light-9) 84%, var(--el-bg-color) 16%);
  }

  .gva-ai-chat-file-icon {
    font-size: 22px;
  }

  .gva-ai-chat-upload-meta {
    padding: 4px 6px 6px;
  }

  .gva-ai-chat-upload-name {
    font-size: 11px;
    color: var(--el-text-color-secondary);
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
  }

  .gva-ai-chat-upload-size {
    margin-top: 2px;
    font-size: 11px;
    color: var(--el-text-color-placeholder);
  }

  .gva-ai-chat-upload-remove {
    position: absolute;
    top: 2px;
    right: 2px;
    background: rgb(255 255 255 / 90%);
  }

  :global(html.dark) .gva-ai-chat-upload-remove {
    background: rgb(15 23 42 / 85%);
  }

  .gva-ai-chat-input :deep(.el-textarea__inner) {
    border-radius: var(--gva-ai-radius-card);
  }

  .gva-ai-chat-footer-actions {
    margin-top: 10px;
    display: flex;
    justify-content: space-between;
    gap: 8px;
  }

  .gva-ai-chat-footer-left,
  .gva-ai-chat-footer-right {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }

  .gva-ai-chat-footer-left :deep(.el-button) {
    min-height: 28px;
    padding: 6px 10px;
  }

  .gva-ai-chat-file-input {
    display: none;
  }

  .gva-ai-chat-drop-mask {
    position: absolute;
    inset: 0;
    z-index: 5;
    background: rgb(59 130 246 / 16%);
    backdrop-filter: blur(2px);
    border: 2px dashed var(--el-color-primary);
    border-radius: var(--gva-ai-radius-panel);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 6px;
    color: var(--el-color-primary-dark-2);
    pointer-events: none;
  }

  :global(html.dark) .gva-ai-chat-drop-mask {
    background: rgb(59 130 246 / 24%);
    color: var(--el-color-primary-light-3);
  }

  .gva-ai-chat-drop-icon {
    font-size: 34px;
  }

  .gva-ai-chat-drop-text {
    font-size: 15px;
    font-weight: 700;
  }

  .gva-ai-chat-drop-subtext {
    font-size: 12px;
    color: var(--el-text-color-secondary);
  }

  .gva-ai-chat-ball-wrap {
    position: fixed;
    pointer-events: auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 5px;
  }

  .gva-ai-chat-ball {
    width: 58px;
    height: 58px;
    background: linear-gradient(
      140deg,
      var(--el-color-primary),
      color-mix(in srgb, var(--el-color-primary-dark-2) 88%, #000 12%)
    );
    box-shadow:
      0 14px 26px rgb(37 99 235 / 30%),
      0 4px 10px rgb(0 0 0 / 20%);
    border: 2px solid var(--gva-ai-layer-1);
    transition:
      transform 0.2s ease,
      box-shadow 0.2s ease;
  }

  :global(html.dark) .gva-ai-chat-ball {
    box-shadow:
      0 14px 24px color-mix(in srgb, var(--el-color-primary) 30%, #000 70%),
      0 6px 14px rgb(0 0 0 / 34%);
    border-color: color-mix(in srgb, var(--el-color-primary-light-8) 20%, #0f172a 80%);
  }

  .gva-ai-chat-ball :deep(.el-icon) {
    font-size: 22px;
  }

  .gva-ai-chat-ball-icon {
    display: inline-flex;
    transform-origin: 50% 58%;
  }

  .gva-ai-chat-ball-text {
    padding: 0 6px;
    line-height: 1;
    font-size: 12px;
    font-weight: 600;
    color: var(--el-text-color-secondary);
    text-shadow: 0 1px 2px rgb(0 0 0 / 16%);
    user-select: none;
    pointer-events: none;
    white-space: nowrap;
    transition:
      transform 0.2s ease,
      color 0.2s ease;
  }

  .gva-ai-chat-ball-wrap:hover .gva-ai-chat-ball {
    transform: translateY(-1px);
  }

  .gva-ai-chat-ball-wrap:hover .gva-ai-chat-ball-icon {
    animation: gva-ai-chat-ball-wiggle 0.5s ease;
  }

  .gva-ai-chat-ball-wrap:hover .gva-ai-chat-ball-text {
    transform: scale(1.12);
    color: var(--el-color-primary);
  }

  :global(html.dark) .gva-ai-chat-ball-text {
    color: var(--el-text-color-regular);
  }

  :global(html.dark) .gva-ai-chat-ball-wrap:hover .gva-ai-chat-ball-text {
    color: var(--el-color-primary-light-3);
  }

  @keyframes gva-ai-chat-ball-wiggle {
    0% {
      transform: rotate(0deg) scale(1);
    }
    30% {
      transform: rotate(-12deg) scale(1.08);
    }
    60% {
      transform: rotate(10deg) scale(1.06);
    }
    100% {
      transform: rotate(0deg) scale(1);
    }
  }

  .gva-ai-chat-panel-enter-active,
  .gva-ai-chat-panel-leave-active {
    transition: all 0.2s ease;
  }

  .gva-ai-chat-panel-enter-from,
  .gva-ai-chat-panel-leave-to {
    opacity: 0;
    transform: translateY(8px) scale(0.98);
  }

  .gva-ai-chat-drop-mask-enter-active,
  .gva-ai-chat-drop-mask-leave-active {
    transition: opacity 0.2s ease;
  }

  .gva-ai-chat-drop-mask-enter-from,
  .gva-ai-chat-drop-mask-leave-to {
    opacity: 0;
  }
</style>

