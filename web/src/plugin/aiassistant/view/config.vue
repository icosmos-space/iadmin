<template>
  <div class="ai-config-container">
    <div class="app-card">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>AI 助手配置</span>
            <el-button type="primary" @click="saveConfig" :loading="saving">
              <el-icon><Check /></el-icon>
              保存配置
            </el-button>
          </div>
        </template>

        <el-form
          ref="configFormRef"
          :model="formData"
          :rules="rules"
          label-width="120px"
        >
          <el-form-item label="启用状态" prop="enabled">
            <el-switch
              v-model="formData.enabled"
              active-text="启用"
              inactive-text="禁用"
            />
            <div class="form-tip">
              禁用后将无法使用 AI 聊天功能
            </div>
          </el-form-item>

          <el-divider />

          <el-form-item label="API 地址" prop="baseURL">
            <el-input
              v-model="formData.baseURL"
              placeholder="https://api.openai.com"
              clearable
            />
            <div class="form-tip">
              AI 服务的 Base URL，例如：https://api.openai.com
            </div>
          </el-form-item>

          <el-form-item label="API 令牌" prop="token">
            <el-input
              v-model="formData.token"
              type="password"
              placeholder="sk-..."
              clearable
              show-password
            />
            <div class="form-tip">
              AI 服务的访问令牌，请妥善保管
            </div>
          </el-form-item>

          <el-divider />

          <el-form-item label="模型名称" prop="model">
            <el-input
              v-model="formData.model"
              placeholder="gpt-4o-mini"
              clearable
            >
              <template #append>
                <el-select
                  v-model="formData.model"
                  placeholder="选择模型"
                  style="width: 180px"
                >
                  <el-option label="GPT-4o Mini" value="gpt-4o-mini" />
                  <el-option label="GPT-4o" value="gpt-4o" />
                  <el-option label="GPT-4 Turbo" value="gpt-4-turbo" />
                  <el-option label="GPT-3.5 Turbo" value="gpt-3.5-turbo" />
                  <el-option label="Claude 3.5 Sonnet" value="claude-3-5-sonnet" />
                  <el-option label="Claude 3 Opus" value="claude-3-opus" />
                </el-select>
              </template>
            </el-input>
            <div class="form-tip">
              使用的 AI 模型名称
            </div>
          </el-form-item>

          <el-form-item label="聊天路径" prop="chatPath">
            <el-input
              v-model="formData.chatPath"
              placeholder="/v1/chat/completions"
              clearable
            />
            <div class="form-tip">
              聊天接口路径，通常为 /v1/chat/completions
            </div>
          </el-form-item>

          <el-divider />

          <el-form-item label="温度" prop="temperature">
            <el-slider
              v-model="formData.temperature"
              :min="0"
              :max="2"
              :step="0.1"
              :marks="{ 0: '严谨', 0.7: '平衡', 2: '发散' }"
              show-input
            />
            <div class="form-tip">
              控制输出的随机性，值越高越发散，越低越严谨
            </div>
          </el-form-item>

          <el-form-item label="超时时间" prop="timeoutSec">
            <el-input-number
              v-model="formData.timeoutSec"
              :min="10"
              :max="300"
              :step="10"
              style="width: 200px"
            />
            <span style="margin-left: 12px">秒</span>
            <div class="form-tip">
              请求超时时间，建议 60 秒以上
            </div>
          </el-form-item>

          <el-divider />

          <el-form-item label="测试连接">
            <el-button type="success" @click="testConnection" :loading="testing">
              <el-icon><Connection /></el-icon>
              测试连接
            </el-button>
            <span v-if="testResult" :class="['test-result', testResult.success ? 'success' : 'error']">
              {{ testResult.message }}
            </span>
          </el-form-item>
        </el-form>
      </el-card>
    </div>

    <!-- 配置说明 -->
    <div class="app-card">
      <el-card>
        <template #header>
          <span>配置说明</span>
        </template>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="支持的 API 格式">
            兼容 OpenAI API 格式的服务，包括 OpenAI、Azure OpenAI、本地部署的 Ollama、LocalAI 等
          </el-descriptions-item>
          <el-descriptions-item label="推荐模型">
            日常使用推荐 gpt-4o-mini，性价比高；复杂任务可使用 gpt-4o
          </el-descriptions-item>
          <el-descriptions-item label="安全提示">
            API 令牌会加密存储，仅服务器端使用，不会暴露给前端
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>
  </div>
</template>

<script setup>
  import { ref, reactive, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getConfig, updateConfig, chatStream } from '@/plugin/aiassistant/api/assistant'

  defineOptions({
    name: 'AIAssistantConfig'
  })

  const configFormRef = ref()
  const saving = ref(false)
  const testing = ref(false)
  const testResult = ref(null)

  const formData = ref({
    enabled: false,
    baseURL: '',
    token: '',
    model: 'gpt-4o-mini',
    chatPath: '/v1/chat/completions',
    temperature: 0.7,
    timeoutSec: 60
  })

  const rules = reactive({
    baseURL: [
      { required: true, message: '请输入 API 地址', trigger: 'blur' },
      { pattern: '^https?://', message: '请输入有效的 URL 地址', trigger: 'blur' }
    ],
    token: [
      { required: true, message: '请输入 API 令牌', trigger: 'blur' }
    ],
    model: [
      { required: true, message: '请输入模型名称', trigger: 'blur' }
    ]
  })

  // 加载配置
  const loadConfig = async () => {
    const res = await getConfig()
    if (res.code === 0) {
      formData.value = {
        enabled: res.data.enabled || false,
        baseURL: res.data.baseURL || '',
        token: res.data.token || '',
        model: res.data.model || 'gpt-4o-mini',
        chatPath: res.data.chatPath || '/v1/chat/completions',
        temperature: res.data.temperature || 0.7,
        timeoutSec: res.data.timeoutSec || 60
      }
    }
  }

  // 保存配置
  const saveConfig = async () => {
    if (!configFormRef.value) return
    
    await configFormRef.value.validate(async (valid) => {
      if (!valid) return

      saving.value = true
      try {
        const res = await updateConfig(formData.value)
        if (res.code === 0) {
          ElMessage.success('配置保存成功')
        } else {
          ElMessage.error(res.msg || '保存失败')
        }
      } catch (error) {
        ElMessage.error(error.message || '保存失败')
      } finally {
        saving.value = false
      }
    })
  }

  // 测试连接
  const testConnection = async () => {
    if (!formData.value.baseURL || !formData.value.token) {
      ElMessage.warning('请先填写 API 地址和令牌')
      return
    }

    testing.value = true
    testResult.value = null

    try {
      const res = await chatStream({
        messages: [{ role: 'user', content: 'Hello' }],
        stream: false,
        model: formData.value.model
      })

      if (res.code === 0) {
        testResult.value = {
          success: true,
          message: '连接测试成功！'
        }
        ElMessage.success('连接测试成功')
      } else {
        testResult.value = {
          success: false,
          message: res.msg || '连接失败'
        }
        ElMessage.error(testResult.value.message)
      }
    } catch (error) {
      testResult.value = {
        success: false,
        message: error.message || '连接失败'
      }
      ElMessage.error(testResult.value.message)
    } finally {
      testing.value = false
    }
  }

  onMounted(() => {
    loadConfig()
  })
</script>

<style scoped>
  .ai-config-container {
    padding: 20px;
  }

  .app-card {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .form-tip {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
    line-height: 1.5;
  }

  .test-result {
    margin-left: 12px;
    font-size: 14px;
  }

  .test-result.success {
    color: #67c23a;
  }

  .test-result.error {
    color: #f56c6c;
  }

  .el-divider {
    margin: 16px 0;
  }
</style>
