<template>
  <div
    id="app"
    class="bg-gray-50 text-slate-700 !dark:text-slate-500 dark:bg-slate-800"
  >
    <el-config-provider :locale="zhCn" :size="appStore.config.global_size">
      <router-view />
      <Application />
      <AiSuspendedBallChat v-if="showAiChat" />
    </el-config-provider>
  </div>
</template>

<script setup>
  import { computed } from 'vue'
  import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
  import Application from '@/components/application/index.vue'
  import AiSuspendedBallChat from '@/components/aiSuspendedBallChat/index.vue'
  import { useAppStore } from '@/pinia'
  import { useRoute } from 'vue-router'

  const appStore = useAppStore()
  const route = useRoute()

  const showAiChat = computed(() => {
    return appStore.config.enable_ai_chat && !route.meta?.public
  })

  defineOptions({
    name: 'App'
  })
</script>
<style lang="scss">
  // 引入初始化样式
  #app {
    height: 100vh;
    overflow: hidden;
    font-weight: 400 !important;
  }

  .el-button {
    font-weight: 400 !important;
  }

  .gva-body-h {
    min-height: calc(100% - 3rem);
  }

  .gva-container {
    height: calc(100% - 2.5rem);
  }

  .gva-container2 {
    height: calc(100% - 4.5rem);
  }
</style>
