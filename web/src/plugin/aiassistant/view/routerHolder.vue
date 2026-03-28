<template>
  <router-view v-slot="{ Component }">
    <transition mode="out-in" name="el-fade-in-linear">
      <keep-alive :include="routerViewKeepAlive">
        <component :is="Component" />
      </keep-alive>
    </transition>
  </router-view>
</template>

<script setup>
  import { computed } from 'vue'
  import { useRoute } from 'vue-router'
  import { useRouterViewKeepAlive } from '@/utils/router'

  defineOptions({
    name: 'AIAssistantHolder'
  })

  const route = useRoute()
  const routerViewKeepAlive = useRouterViewKeepAlive(route)
</script>

<style scoped>
</style>
