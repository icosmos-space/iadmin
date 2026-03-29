<template>
  <div class="mk-page" @scroll="handleScroll">
    <div class="mk-header-shell" :class="{ 'is-scrolled': scrolled }">
      <header class="mk-header-inner">
        <RouterLink class="mk-brand" to="/landing">
          <span class="mk-brand-dot" />
          <span>iAdmin Orbit</span>
        </RouterLink>

        <nav class="mk-nav">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="mk-nav-link"
            active-class="is-active"
          >
            {{ item.label }}
          </RouterLink>
        </nav>

        <el-space v-if="!isLogin" class="mk-actions" alignment="center">
          <RouterLink to="/init">
            <el-button plain>初始化</el-button>
          </RouterLink>
          <RouterLink to="/login">
            <el-button type="primary">登录后台</el-button>
          </RouterLink>
        </el-space>

        <div v-else class="mk-actions mk-account-wrap">
          <el-dropdown trigger="click">
            <button class="mk-account-btn" type="button">
              <span class="mk-account-name">{{ accountName }}</span>
              <el-icon><arrow-down /></el-icon>
            </button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item disabled>
                  当前角色：{{ roleName }}
                </el-dropdown-item>
                <el-dropdown-item @click="goConsole">进入控制台</el-dropdown-item>
                <el-dropdown-item @click="goProfile">个人信息</el-dropdown-item>
                <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>
    </div>

    <main class="mk-main">
      <slot />
    </main>

    <footer class="mk-footer">
      <div class="mk-footer-inner">
        <div class="mk-footer-left">
          <p class="mk-footer-copy">© {{ currentYear }} iAdmin. All rights reserved.</p>
          <p class="mk-footer-desc">A management platform powered by Golang + Vue</p>
        </div>
        <el-space alignment="center">
          <RouterLink class="mk-footer-link" to="/privacy">Privacy</RouterLink>
          <el-divider direction="vertical" />
          <RouterLink class="mk-footer-link" to="/terms">Terms</RouterLink>
        </el-space>
      </div>
    </footer>
  </div>
</template>

<script setup>
  import { computed, onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useUserStore } from '@/pinia/modules/user'

  const scrolled = ref(false)
  const currentYear = new Date().getFullYear()
  const router = useRouter()
  const userStore = useUserStore()

  const navItems = [
    { label: '首页', to: '/landing' },
    { label: '功能', to: '/features' },
    { label: '价格', to: '/pricing' },
    { label: '联系', to: '/contact' }
  ]

  const isLogin = computed(() => !!userStore.token)
  const accountName = computed(() => userStore.userInfo?.nickName || '账户')
  const roleName = computed(() => userStore.userInfo?.authority?.authorityName || '已登录')

  const handleScroll = (event) => {
    const top = event?.target?.scrollTop || 0
    scrolled.value = top > 28
  }

  const goProfile = () => {
    router.push({ name: 'person' })
  }

  const goConsole = async () => {
    let target = userStore.userInfo?.authority?.defaultRouter
    if (!target) {
      try {
        await userStore.GetUserInfo()
      } catch (_) {
        // ignore and use fallback below
      }
      target = userStore.userInfo?.authority?.defaultRouter
    }
    if (target) {
      router.push({ name: target })
    } else {
      router.push({ name: 'Login' })
    }
  }

  const logout = () => {
    userStore.LoginOut()
  }

  onMounted(async () => {
    if (isLogin.value && !userStore.userInfo?.nickName) {
      try {
        await userStore.GetUserInfo()
      } catch (_) {
        // ignore
      }
    }
  })
</script>

<style scoped lang="scss">
  .mk-page {
    position: relative;
    width: 100%;
    height: 100vh;
    overflow: auto;
    background: #f6f9ff;
    color: #1f3a68;
    font-family: 'Segoe UI', 'Microsoft YaHei', 'PingFang SC', sans-serif;
  }

  .mk-header-shell {
    position: fixed;
    inset: 0 0 auto;
    z-index: 40;
    transition: all 0.28s ease;
  }

  .mk-header-inner {
    width: 100%;
    display: grid;
    grid-template-columns: auto 1fr auto;
    align-items: center;
    gap: 18px;
    padding: 18px 28px;
    border-bottom: 1px solid rgba(59, 130, 246, 0.14);
    background: rgba(255, 255, 255, 0.92);
    transition: all 0.28s ease;
  }

  .mk-header-shell.is-scrolled .mk-header-inner {
    max-width: 1240px;
    margin: 10px auto 0;
    padding: 12px 18px;
    border: 1px solid rgba(59, 130, 246, 0.2);
    border-radius: 16px;
    background: rgba(255, 255, 255, 0.95);
    box-shadow: 0 10px 30px rgba(59, 130, 246, 0.12);
  }

  .mk-brand {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    color: #1e3a8a;
    text-decoration: none;
    font-size: 19px;
    font-weight: 700;
  }

  .mk-brand-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background: #2563eb;
    box-shadow: 0 0 0 6px rgba(37, 99, 235, 0.12);
  }

  .mk-nav {
    justify-self: center;
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }

  .mk-nav-link {
    text-decoration: none;
    color: #35568c;
    padding: 8px 14px;
    border-radius: 999px;
    font-size: 14px;
    transition: all 0.2s ease;
  }

  .mk-nav-link:hover,
  .mk-nav-link.is-active {
    color: #1e40af;
    background: rgba(59, 130, 246, 0.14);
  }

  .mk-actions {
    justify-self: end;
  }

  .mk-account-wrap {
    display: flex;
    align-items: center;
  }

  .mk-account-btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    border-radius: 10px;
    border: 1px solid rgba(37, 99, 235, 0.2);
    background: #ffffff;
    color: #1f3a68;
    cursor: pointer;
  }

  .mk-account-btn:hover {
    border-color: rgba(37, 99, 235, 0.35);
  }

  .mk-account-name {
    font-weight: 600;
    max-width: 140px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .mk-main {
    margin: 0 auto;
    max-width: 1240px;
    padding: 122px 24px 36px;
  }

  .mk-footer {
    margin-top: 12px;
    border-top: 1px solid rgba(59, 130, 246, 0.18);
    background: rgba(255, 255, 255, 0.6);
  }

  .mk-footer-inner {
    margin: 0 auto;
    max-width: 1240px;
    padding: 20px 24px 26px;
    display: flex;
    justify-content: space-between;
    gap: 10px;
    align-items: center;
  }

  .mk-footer-left {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .mk-footer-copy {
    margin: 0;
    font-size: 13px;
    color: #35568c;
  }

  .mk-footer-desc {
    margin: 0;
    font-size: 12px;
    color: #5f7fb3;
  }

  .mk-footer-link {
    text-decoration: none;
    color: #2563eb;
    font-size: 13px;
    font-weight: 600;
  }

  .mk-footer-link:hover {
    color: #1d4ed8;
  }

  @media (max-width: 980px) {
    .mk-header-inner {
      grid-template-columns: 1fr;
      padding: 14px 16px;
      gap: 12px;
    }

    .mk-header-shell.is-scrolled .mk-header-inner {
      margin: 8px 12px 0;
      max-width: none;
      padding: 12px 12px;
    }

    .mk-nav,
    .mk-actions {
      justify-self: stretch;
    }

    .mk-actions {
      justify-content: space-between;
    }

    .mk-account-wrap {
      justify-content: flex-end;
    }

    .mk-main {
      padding: 160px 16px 30px;
    }

    .mk-footer-inner {
      padding: 18px 16px 22px;
      flex-direction: column;
      align-items: flex-start;
    }
  }
</style>
