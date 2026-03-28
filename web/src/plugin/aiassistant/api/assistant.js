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
  let isConnected = false

  // 立即执行异步操作
  ;(async () => {
    try {
      const { useUserStore } = await import('@/pinia/modules/user')
      const store = useUserStore()
      const token = store.token
      const userId = store.userInfo.ID

      const baseURL = import.meta.env.VITE_BASE_API
      const url = `${baseURL}/aiAssistant/chat`

      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
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
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.msg || `HTTP ${response.status}`)
      }

      isConnected = true
      const reader = response.body.getReader()
      const decoder = new TextDecoder('utf-8')
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })

        // 处理 SSE 格式
        const lines = buffer.split('\n')
        buffer = lines.pop() || '' // 保留不完整行

        for (const line of lines) {
          const trimmedLine = line.trim()
          if (!trimmedLine || trimmedLine.startsWith(':')) continue

          // 解析 SSE 事件
          if (trimmedLine.startsWith('event: ')) {
            const eventType = trimmedLine.slice(7)
            // 读取下一行 data
            const nextLine = lines.shift()
            if (nextLine && nextLine.trim().startsWith('data: ')) {
              const dataStr = nextLine.trim().slice(6)
              try {
                const parsed = JSON.parse(dataStr)
                
                switch (eventType) {
                  case 'meta':
                    // 连接建立，无需处理
                    break
                  case 'error':
                    if (onError) onError(new Error(parsed.msg || 'SSE error'))
                    return
                  case 'done':
                    if (onComplete) onComplete()
                    return
                }
              } catch (e) {
                console.warn('SSE 事件解析失败:', e)
              }
            }
            continue
          }

          // 解析 data 行（AI 原始内容）
          if (trimmedLine.startsWith('data: ')) {
            const dataStr = trimmedLine.slice(6)
            if (dataStr === '[DONE]') continue

            try {
              const parsed = JSON.parse(dataStr)
              const delta = parsed.choices?.[0]?.delta?.content
              if (delta) {
                onMessage(delta)
              }
            } catch (e) {
              console.warn('SSE data 解析失败:', e)
            }
          }
        }
      }

      if (onComplete) onComplete()
    } catch (error) {
      if (error.name === 'AbortError') {
        console.log('SSE 请求已取消')
        return
      }
      if (onError) onError(error)
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

// 非流式聊天（备用）
export const chat = (data) => {
  return service({
    url: '/aiAssistant/chat',
    method: 'post',
    data,
    timeout: 60000
  })
}
