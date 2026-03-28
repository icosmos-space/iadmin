// AI Assistant Plugin
// 导出 API
export * from './api/assistant'

// 导出页面视图
export { default as AIChatPage } from './view/chat.vue'
export { default as AIConfig } from './view/config.vue'
export { default as AIPrompt } from './view/prompt.vue'

// 导出组件（可复用）
export { default as AIChat } from './components/chat.vue'
