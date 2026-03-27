<template>
  <div class="sns-bind-panel bg-white dark:bg-slate-800 rounded-xl p-6 profile-card">
    <h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
      <el-icon class="text-blue-500"><Link /></el-icon>
      SNS 账号绑定
    </h2>
    <SnsAuthBtnGroup mode="bind" @success="onBindSuccess" />
    <div class="mt-4">
      <el-table :data="bindings" size="small">
        <el-table-column label="平台" width="140">
          <template #default="scope">{{ providerLabel(scope.row.provider) }}</template>
        </el-table-column>
        <el-table-column label="OpenID" prop="openID" min-width="220" />
        <el-table-column label="绑定时间" min-width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button type="danger" link @click="doUnbind(scope.row.provider)">
              解绑
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import SnsAuthBtnGroup from '@/plugin/snsauth/components/sns-auth-btn-group.vue'
  import { getMySnsBindings, unbindSns } from '@/plugin/snsauth/api/auth'

  const bindings = ref([])

  const providerLabel = (provider) => {
    const map = {
      github: 'GitHub',
      feishu: '飞书',
      wechat: '微信',
      telegram: 'Telegram'
    }
    return map[provider] || provider
  }

  const formatDate = (value) => {
    if (!value) return '-'
    const date = new Date(value)
    if (Number.isNaN(date.getTime())) return '-'
    return date.toLocaleString()
  }

  const loadBindings = async () => {
    const res = await getMySnsBindings()
    if (res.code === 0) {
      bindings.value = res.data || []
    }
  }

  const onBindSuccess = () => {
    loadBindings()
  }

  const doUnbind = async (provider) => {
    await ElMessageBox.confirm('确定解绑该SNS账号吗？', '提示', { type: 'warning' })
    const res = await unbindSns({ provider })
    if (res.code === 0) {
      ElMessage.success('解绑成功')
      loadBindings()
    }
  }

  loadBindings()
</script>
