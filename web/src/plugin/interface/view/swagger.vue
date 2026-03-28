<template>
  <div class="swagger-plugin-page">
    <el-alert
      class="mb-12"
      :closable="false"
      title="接口管理（Swagger）"
      type="info"
      description="当前页面用于查看后端 Swagger 文档。仅在后端以 dev 标签构建时可用。"
      show-icon
    />
    <iframe class="swagger-frame" :src="swaggerUrl" />
  </div>
</template>

<script setup>
  import { computed } from 'vue'

  defineOptions({
    name: 'ApiSwaggerManager'
  })

  const swaggerUrl = computed(() => {
    const origin = window.location.origin
    const base = import.meta.env.VITE_BASE_API || ''
    const normalizedBase = base.startsWith('/') ? base : `/${base}`
    return `${origin}${normalizedBase}/swagger/index.html`
  })
</script>

<style scoped>
  .swagger-plugin-page {
    width: 100%;
    height: calc(100vh - 140px);
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .mb-12 {
    margin-bottom: 12px;
  }

  .swagger-frame {
    flex: 1;
    width: 100%;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 6px;
    background-color: #fff;
  }
</style>
