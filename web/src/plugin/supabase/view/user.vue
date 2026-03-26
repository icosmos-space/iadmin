<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="关键词">
          <el-input
            v-model="searchInfo.keyword"
            clearable
            placeholder="邮箱/手机号/用户ID"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            查询
          </el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table :data="tableData" style="width: 100%">
        <el-table-column label="User ID" prop="id" min-width="280" />
        <el-table-column label="邮箱" prop="email" min-width="220" />
        <el-table-column label="手机号" prop="phone" min-width="160" />
        <el-table-column label="邮箱验证" min-width="120">
          <template #default="scope">
            {{ scope.row.email_confirmed_at ? '是' : '否' }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="最后登录" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.last_sign_in_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="140">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              @click="openResetPwdDialog(scope.row)"
            >
              修改密码
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" title="修改 Supabase 用户密码" width="520px">
      <el-form ref="elFormRef" :model="formData" :rules="rule" label-width="100px">
        <el-form-item label="用户ID">
          <el-input v-model="formData.userID" disabled />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="formData.newPassword"
            show-password
            placeholder="请输入新密码（至少6位）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeDialog">取消</el-button>
        <el-button type="primary" @click="submitResetPwd">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import {
    getSupabaseUserList,
    updateSupabaseUserPassword
  } from '@/plugin/supabase/api/user'
  import { ElMessage } from 'element-plus'
  import { reactive, ref } from 'vue'

  defineOptions({
    name: 'SupabaseUser'
  })

  const formatDate = (value) => {
    if (!value) return '-'
    const date = new Date(value)
    if (Number.isNaN(date.getTime())) return '-'
    return date.toLocaleString()
  }

  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const tableData = ref([])
  const searchInfo = ref({
    keyword: ''
  })

  const getTableData = async () => {
    const res = await getSupabaseUserList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
      page.value = res.data.page || 1
      pageSize.value = res.data.pageSize || 10
    }
  }

  const onSubmit = () => {
    page.value = 1
    getTableData()
  }

  const onReset = () => {
    searchInfo.value = { keyword: '' }
    page.value = 1
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const dialogVisible = ref(false)
  const elFormRef = ref()
  const formData = ref({
    userID: '',
    newPassword: ''
  })
  const rule = reactive({
    newPassword: [
      { required: true, message: '请输入新密码', trigger: 'blur' },
      { min: 6, message: '密码至少 6 位', trigger: 'blur' }
    ]
  })

  const openResetPwdDialog = (row) => {
    formData.value = {
      userID: row.id,
      newPassword: ''
    }
    dialogVisible.value = true
  }

  const closeDialog = () => {
    dialogVisible.value = false
    formData.value = {
      userID: '',
      newPassword: ''
    }
  }

  const submitResetPwd = () => {
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return
      const res = await updateSupabaseUserPassword(formData.value)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '密码修改成功'
        })
        closeDialog()
      }
    })
  }

  getTableData()
</script>
