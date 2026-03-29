import service from '@/utils/request'

// ======================= 配置管理 =======================

// @Tags AIAssistant
// @Summary 获取 AI 助手配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /aiAssistant/getConfig [get]
export const getConfig = () => {
  return service({
    url: '/aiAssistant/getConfig',
    method: 'get'
  })
}

// @Tags AIAssistant
// @Summary 更新 AI 助手配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AIAssistantConfig true "更新 AI 助手配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /aiAssistant/updateConfig [put]
export const updateConfig = (data) => {
  return service({
    url: '/aiAssistant/updateConfig',
    method: 'put',
    data
  })
}

// ======================= 预置提示词管理 =======================

// @Tags AIAssistant
// @Summary 创建预置提示词
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AIAssistantPrompt true "创建预置提示词"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /aiAssistant/createPrompt [post]
export const createPrompt = (data) => {
  return service({
    url: '/aiAssistant/createPrompt',
    method: 'post',
    data
  })
}

// @Tags AIAssistant
// @Summary 删除预置提示词
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AIAssistantPrompt true "删除预置提示词"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /aiAssistant/deletePrompt [delete]
export const deletePrompt = (params) => {
  return service({
    url: '/aiAssistant/deletePrompt',
    method: 'delete',
    params
  })
}

// @Tags AIAssistant
// @Summary 批量删除预置提示词
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除预置提示词"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /aiAssistant/deletePromptByIds [delete]
export const deletePromptByIds = (params) => {
  return service({
    url: '/aiAssistant/deletePromptByIds',
    method: 'delete',
    params
  })
}

// @Tags AIAssistant
// @Summary 更新预置提示词
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AIAssistantPrompt true "更新预置提示词"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /aiAssistant/updatePrompt [put]
export const updatePrompt = (data) => {
  return service({
    url: '/aiAssistant/updatePrompt',
    method: 'put',
    data
  })
}

// @Tags AIAssistant
// @Summary 用 id 查询预置提示词
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.AIAssistantPrompt true "用 id 查询预置提示词"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /aiAssistant/findPrompt [get]
export const findPrompt = (params) => {
  return service({
    url: '/aiAssistant/findPrompt',
    method: 'get',
    params
  })
}

// @Tags AIAssistant
// @Summary 分页获取预置提示词列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取预置提示词列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /aiAssistant/getPromptList [get]
export const getPromptList = (params) => {
  return service({
    url: '/aiAssistant/getPromptList',
    method: 'get',
    params
  })
}

// @Tags AIAssistant
// @Summary 获取已启用的预置提示词
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /aiAssistant/getEnabledPrompts [get]
export const getEnabledPrompts = () => {
  return service({
    url: '/aiAssistant/getEnabledPrompts',
    method: 'get'
  })
}

// ======================= AI 聊天（SSE 流式） =======================

const joinPath = (base, path) => {
  const cleanBase = (base || '').replace(/\/+$/, '')
  if (!path) return cleanBase
  return `${cleanBase}/${path.replace(/^\/+/, '')}`
}

const parseSSEEvent = (eventRaw) => {
  const lines = eventRaw.split(/\r?\n/)
  let event = 'message'
  const dataParts = []

  for (const rawLine of lines) {
    if (!rawLine || rawLine.startsWith(':')) continue
    if (rawLine.startsWith('event:')) {
      event = rawLine.slice(6).trim() || 'message'
      continue
    }
    if (rawLine.startsWith('data:')) {
      let data = rawLine.slice(5)
      if (data.startsWith(' ')) data = data.slice(1)
      dataParts.push(data.replace(/\r$/, ''))
      continue
    }
    dataParts.push(rawLine.replace(/\r$/, ''))
  }

  return {
    event,
    data: dataParts.join('\n')
  }
}

const extractStreamError = (payload, fallback = 'SSE error') => {
  if (!payload || typeof payload !== 'object') return fallback
  if (typeof payload.message === 'string' && payload.message.trim()) return payload.message.trim()
  if (typeof payload.msg === 'string' && payload.msg.trim()) return payload.msg.trim()
  if (payload.error && typeof payload.error === 'object') {
    if (typeof payload.error.message === 'string' && payload.error.message.trim()) {
      return payload.error.message.trim()
    }
  }
  return fallback
}

const extractStreamDelta = (payload) => {
  if (typeof payload === 'string') return payload
  if (!payload || typeof payload !== 'object') return ''

  if (Array.isArray(payload.choices)) {
    const parts = []
    payload.choices.forEach((choice) => {
      const delta = choice?.delta?.content
      const message = choice?.message?.content
      const text = choice?.text
      if (typeof delta === 'string') parts.push(delta)
      if (typeof message === 'string') parts.push(message)
      if (typeof text === 'string') parts.push(text)
    })
    if (parts.length) return parts.join('')
  }

  if (typeof payload.content === 'string') return payload.content
  if (typeof payload.output_text === 'string') return payload.output_text
  if (typeof payload?.delta?.content === 'string') return payload.delta.content

  return ''
}

const isStreamDone = (payload) => {
  if (!payload || typeof payload !== 'object') return false
  if (payload.done === true || payload.is_end === true || payload.finish === true) return true
  if (Array.isArray(payload.choices)) {
    return payload.choices.some((choice) => {
      return typeof choice?.finish_reason === 'string' && choice.finish_reason.length > 0
    })
  }
  return false
}

/**
 * 使用 SSE 流式聊天
 * @param {Object} data - 聊天数据
 * @param {Function} onMessage - 接收消息回调 (content: string)
 * @param {Function} onError - 错误回调
 * @param {Function} onComplete - 完成回调
 * @returns {Object} 控制器 { abort: Function }
 */
export const chatStream = (data, onMessage, onError, onComplete) => {
  let controller = new AbortController()
  let isCompleted = false

  const completeOnce = () => {
    if (isCompleted) return
    isCompleted = true
    if (onComplete) onComplete()
  }

  const errorOnce = (error) => {
    if (isCompleted) return
    if (onError) onError(error instanceof Error ? error : new Error(String(error)))
  }

  ;(async () => {
    try {
      const { useUserStore } = await import('@/pinia/modules/user')
      const store = useUserStore()
      const token = store.token || ''
      const userId = store.userInfo?.ID || ''
      const url = joinPath(import.meta.env.VITE_BASE_API, '/aiAssistant/chat')

      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'text/event-stream, application/json',
          'x-token': token,
          'x-user-id': userId
        },
        body: JSON.stringify({
          ...data,
          stream: true
        }),
        signal: controller.signal
      })

      if (!response.ok) {
        const errorText = await response.text()
        let msg = `HTTP ${response.status}`
        try {
          const errorData = JSON.parse(errorText)
          msg = errorData?.msg || errorData?.message || msg
        } catch {
          if (errorText?.trim()) msg = errorText.trim()
        }
        throw new Error(msg)
      }

      if (!response.body) {
        throw new Error('SSE response body is empty')
      }

      const reader = response.body.getReader()
      const decoder = new TextDecoder('utf-8')
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) {
          if (buffer.trim()) {
            const event = parseSSEEvent(buffer)
            const finished = handleSSEParsedEvent(event, onMessage, errorOnce)
            if (finished) {
              completeOnce()
              return
            }
          }
          break
        }

        buffer += decoder.decode(value, { stream: true })
        const events = buffer.split(/\r?\n\r?\n/)
        buffer = events.pop() || ''

        for (const eventRaw of events) {
          if (!eventRaw.trim()) continue
          const event = parseSSEEvent(eventRaw)
          const finished = handleSSEParsedEvent(event, onMessage, errorOnce)
          if (finished) {
            completeOnce()
            return
          }
        }
      }

      completeOnce()
    } catch (error) {
      if (error.name === 'AbortError') {
        return
      }
      errorOnce(error)
    }
  })()

  return {
    abort: () => {
      if (controller) {
        controller.abort()
        controller = null
      }
    }
  }
}

const handleSSEParsedEvent = (event, onMessage, onError) => {
  const dataLine = (event.data || '').trim()
  const eventType = event.event || 'message'

  if (!dataLine) return false
  if (dataLine === '[DONE]' || eventType === 'done') return true

  let payload = dataLine
  if (dataLine.startsWith('{') || dataLine.startsWith('[')) {
    try {
      payload = JSON.parse(dataLine)
    } catch {
      payload = dataLine
    }
  }

  if (eventType === 'meta') return false
  if (eventType === 'error') {
    const msg = extractStreamError(payload, dataLine)
    if (onError) onError(new Error(msg))
    return true
  }

  const delta = extractStreamDelta(payload)
  if (typeof delta === 'string' && delta.length > 0 && onMessage) {
    onMessage(delta)
  }

  if (payload && typeof payload === 'object' && isStreamDone(payload)) {
    return true
  }

  return false
}

// 非流式聊天（备用）
export const chat = (data) => {
  return service({
    url: '/aiAssistant/chat',
    method: 'post',
    data,
    timeout: 60000
  })
}
