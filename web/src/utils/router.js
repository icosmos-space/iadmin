import { computed } from 'vue'
import { useRouterStore } from '@/pinia/modules/router'

export const useRouterViewKeepAlive = () => {
  const routerStore = useRouterStore()
  return computed(() => routerStore.keepAliveRouters)
}
