<template>
  <div class="gva-table-box">
    <el-table :data="tableData" border>
      <el-table-column label="平台" prop="provider" width="140" />
      <el-table-column label="启用" width="100">
        <template #default="scope">
          <el-switch v-model="scope.row.enabled" />
        </template>
      </el-table-column>
      <el-table-column label="Client ID" min-width="180">
        <template #default="scope">
          <el-input v-model="scope.row.clientID" placeholder="Client ID" />
        </template>
      </el-table-column>
      <el-table-column label="Client Secret" min-width="180">
        <template #default="scope">
          <el-input v-model="scope.row.clientSecret" placeholder="Client Secret" />
        </template>
      </el-table-column>
      <el-table-column label="Redirect URL" min-width="220">
        <template #default="scope">
          <el-input v-model="scope.row.redirectURL" placeholder="回调地址" />
        </template>
      </el-table-column>
      <el-table-column label="Scopes" min-width="180">
        <template #default="scope">
          <el-input v-model="scope.row.scopes" placeholder="scope1 scope2" />
        </template>
      </el-table-column>
      <el-table-column label="Auth URL" min-width="220">
        <template #default="scope">
          <el-input v-model="scope.row.authURL" placeholder="授权地址" />
        </template>
      </el-table-column>
      <el-table-column label="Token URL" min-width="220">
        <template #default="scope">
          <el-input v-model="scope.row.tokenURL" placeholder="Token地址" />
        </template>
      </el-table-column>
      <el-table-column label="UserInfo URL" min-width="220">
        <template #default="scope">
          <el-input v-model="scope.row.userInfoURL" placeholder="用户信息地址" />
        </template>
      </el-table-column>
      <el-table-column label="操作" fixed="right" width="100">
        <template #default="scope">
          <el-button type="primary" link @click="save(scope.row)">保存</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
  import { ElMessage } from 'element-plus'
  import { ref } from 'vue'
  import { getSnsProviderList, updateSnsProviderConfig } from '@/plugin/snsauth/api/auth'

  defineOptions({
    name: 'SnsAuthConfig'
  })

  const tableData = ref([])

  const load = async () => {
    const res = await getSnsProviderList()
    if (res.code === 0) {
      tableData.value = (res.data || []).map((item) => ({
        provider: item.provider,
        enabled: item.enabled,
        clientID: item.clientID || '',
        clientSecret: item.clientSecret || '',
        redirectURL: item.redirectURL || '',
        scopes: item.scopes || '',
        authURL: item.authURL || '',
        tokenURL: item.tokenURL || '',
        userInfoURL: item.userInfoURL || ''
      }))
    }
  }

  const save = async (row) => {
    const res = await updateSnsProviderConfig(row)
    if (res.code === 0) {
      ElMessage.success('保存成功')
      load()
    }
  }

  load()
</script>
