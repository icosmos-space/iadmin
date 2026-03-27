<template>
  <div class="sns-auth-btn-group">
    <el-button
      v-for="item in providers"
      :key="item.provider"
      class="sns-btn"
      @click="handleClick(item.provider)"
    >
      {{ modeLabel }}{{ providerLabel(item.provider) }}
    </el-button>
  </div>
  <el-dialog v-model="wechatDialogVisible" title="请使用微信扫码继续" width="360px">
    <div class="wechat-qr-wrap">
      <img v-if="wechatQrUrl" :src="wechatQrUrl" alt="wechat auth qrcode" class="wechat-qr-img" />
      <div class="wechat-qr-tip">
        请使用微信扫一扫，完成{{ modeLabel }}授权
      </div>
      <el-link v-if="wechatAuthUrl" :href="wechatAuthUrl" target="_blank" type="primary">
        无法扫码时，点此打开授权链接
      </el-link>
    </div>
  </el-dialog>
</template>

<script setup>
  import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import {
    getEnabledSnsProviderList,
    getSnsBindURL,
    getSnsLoginURL,
    getSnsProviderList,
    telegramBind,
    telegramLogin
  } from '@/plugin/snsauth/api/auth'

  const props = defineProps({
    mode: {
      type: String,
      default: 'login'
    }
  })

  const emit = defineEmits(['success', 'error'])

  const allProviders = ref([])
  const wechatDialogVisible = ref(false)
  const wechatQrUrl = ref('')
  const wechatAuthUrl = ref('')
  const providers = computed(() =>
    allProviders.value.filter((item) => item.enabled)
  )
  const modeLabel = computed(() => (props.mode === 'bind' ? '绑定' : '登录'))

  const providerLabel = (provider) => {
    const map = {
      github: 'GitHub',
      feishu: '飞书',
      wechat: '微信',
      telegram: 'Telegram'
    }
    return map[provider] || provider
  }

  const loadProviders = async () => {
    const api = props.mode === 'bind' ? getSnsProviderList : getEnabledSnsProviderList
    const res = await api()
    if (res.code === 0) {
      allProviders.value = props.mode === 'bind' ? (res.data || []) : (res.data || []).map((i) => ({ ...i, enabled: true }))
    }
  }

  const handleClick = async (provider) => {
    if (provider === 'telegram') {
      const raw = window.prompt('请输入 Telegram Login Widget 返回的 JSON 数据')
      if (!raw) return
      let payload = null
      try {
        payload = JSON.parse(raw)
      } catch (e) {
        ElMessage.error('JSON格式错误')
        return
      }
      const req = props.mode === 'bind' ? telegramBind : telegramLogin
      const res = await req(payload)
      if (res.code !== 0) return
      if (props.mode === 'login') {
        emit('success', res.data)
      } else {
        ElMessage.success(res.data?.msg || '绑定成功')
        emit('success', res.data)
      }
      return
    }
    const req = props.mode === 'bind' ? getSnsBindURL : getSnsLoginURL
    const res = await req({ provider })
    if (res.code !== 0) return
    if (provider === 'wechat') {
      const isWechatBrowser = /MicroMessenger/i.test(navigator.userAgent)
      if (isWechatBrowser) {
        window.location.href = res.data.url
        return
      }
      wechatAuthUrl.value = res.data.url
      wechatQrUrl.value = `https://api.qrserver.com/v1/create-qr-code/?size=220x220&data=${encodeURIComponent(res.data.url)}`
      wechatDialogVisible.value = true
      return
    }
    window.open(res.data.url, '_blank', 'width=720,height=760')
  }

  const onMessage = (event) => {
    const data = event.data
    if (!data || data.type !== 'SNS_AUTH_RESULT') return
    if (!data.ok) {
      ElMessage.error(data.error || 'SNS操作失败')
      emit('error', data)
      return
    }
    ElMessage.success(data.data?.msg || '操作成功')
    emit('success', data.data)
  }

  onMounted(() => {
    loadProviders()
    window.addEventListener('message', onMessage)
  })

  onBeforeUnmount(() => {
    window.removeEventListener('message', onMessage)
  })
</script>

<style scoped>
  .sns-auth-btn-group {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .wechat-qr-wrap {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
  }

  .wechat-qr-img {
    width: 220px;
    height: 220px;
    border-radius: 8px;
  }

  .wechat-qr-tip {
    color: #606266;
    font-size: 13px;
  }
</style>
