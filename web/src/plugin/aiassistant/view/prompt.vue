<template>
  <div class="ai-prompt-container">
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="关键字" prop="keyword">
          <el-input
            v-model="searchInfo.keyword"
            placeholder="搜索标题或内容"
            clearable
          />
        </el-form-item>

        <el-form-item label="状态" prop="enabled">
          <el-select v-model="searchInfo.enabled" placeholder="全部" clearable style="width: 120px">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            查询
          </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">
          新增
        </el-button>
        <el-button
          icon="delete"
          style="margin-left: 10px"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >
          删除
        </el-button>
      </div>

      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />

        <el-table-column align="left" label="ID" prop="ID" width="80" />

        <el-table-column align="left" label="标题" prop="title" min-width="150" />

        <el-table-column align="left" label="内容" prop="content" min-width="300" show-overflow-tooltip />

        <el-table-column align="left" label="排序" prop="sort" width="80" />

        <el-table-column align="left" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.enabled ? 'success' : 'info'">
              {{ scope.row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column align="left" label="操作" fixed="right" min-width="200">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              @click="updatePromptFunc(scope.row)"
            >
              变更
            </el-button>
            <el-button
              type="danger"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
            >
              删除
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

    <!-- 新增/编辑弹窗 -->
    <el-drawer
      v-model="dialogFormVisible"
      destroy-on-close
      size="600"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '添加' : '修改' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog"> 确 定 </el-button>
            <el-button @click="closeDialog"> 取 消 </el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="elFormRef"
        :model="formData"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="标题" prop="title">
          <el-input
            v-model="formData.title"
            placeholder="请输入标题"
            clearable
          />
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <el-input
            v-model="formData.content"
            type="textarea"
            :rows="6"
            placeholder="请输入提示词内容"
            resize="vertical"
          />
        </el-form-item>

        <el-form-item label="排序" prop="sort">
          <el-input-number
            v-model="formData.sort"
            :min="0"
            :max="999"
            style="width: 100%"
          />
          <span style="margin-left: 12px; color: #909399; font-size: 12px;">
            数字越小越靠前
          </span>
        </el-form-item>

        <el-form-item label="状态" prop="enabled">
          <el-switch
            v-model="formData.enabled"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createPrompt,
    deletePrompt,
    deletePromptByIds,
    updatePrompt,
    findPrompt,
    getPromptList
  } from '@/plugin/aiassistant/api/assistant'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref, reactive } from 'vue'

  defineOptions({
    name: 'AIAssistantPrompt'
  })

  // 表单数据
  const formData = ref({
    title: '',
    content: '',
    sort: 0,
    enabled: true
  })

  // 验证规则
  const rules = reactive({
    title: [
      { required: true, message: '请输入标题', trigger: 'blur' }
    ],
    content: [
      { required: true, message: '请输入内容', trigger: 'blur' }
    ]
  })

  const elFormRef = ref()
  const elSearchFormRef = ref()

  // =========== 表格控制部分 ===========
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const searchInfo = ref({})

  // 重置
  const onReset = () => {
    searchInfo.value = {}
    getTableData()
  }

  // 搜索
  const onSubmit = () => {
    elSearchFormRef.value?.validate(async (valid) => {
      if (!valid) return
      page.value = 1
      getTableData()
    })
  }

  // 分页
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // 查询
  const getTableData = async () => {
    const table = await getPromptList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  getTableData()

  // 多选数据
  const multipleSelection = ref([])
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }

  // 删除行
  const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      deletePromptFunc(row)
    })
  }

  // 多选删除
  const onDelete = async () => {
    ElMessageBox.confirm('确定要删除选中的数据吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map((item) => {
          IDs.push(item.ID)
        })
      const res = await deletePromptByIds({ IDs })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === IDs.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
  }

  // 行为控制标记
  const type = ref('')
  const dialogFormVisible = ref(false)

  // 更新行
  const updatePromptFunc = async (row) => {
    const res = await findPrompt({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
      formData.value = {
        ID: res.data.ID,
        title: res.data.title,
        content: res.data.content,
        sort: res.data.sort,
        enabled: res.data.enabled
      }
      dialogFormVisible.value = true
    }
  }

  // 删除行
  const deletePromptFunc = async (row) => {
    const res = await deletePrompt({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  }

  // 打开弹窗
  const openDialog = () => {
    type.value = 'create'
    formData.value = {
      title: '',
      content: '',
      sort: 0,
      enabled: true
    }
    dialogFormVisible.value = true
  }

  // 关闭弹窗
  const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
      title: '',
      content: '',
      sort: 0,
      enabled: true
    }
  }

  // 弹窗确定
  const enterDialog = async () => {
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return
      let res
      switch (type.value) {
        case 'create':
          res = await createPrompt(formData.value)
          break
        case 'update':
          res = await updatePrompt(formData.value)
          break
        default:
          res = await createPrompt(formData.value)
          break
      }
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '操作成功'
        })
        closeDialog()
        getTableData()
      }
    })
  }
</script>

<style scoped>
  .ai-prompt-container {
    padding: 20px;
  }
</style>
