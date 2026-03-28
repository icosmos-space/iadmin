<template>
  <div class="state-page">
    <el-row :gutter="16">
      <el-col :span="12">
        <el-card class="section-card" shadow="never">
          <template #header>
            <div class="section-title">
              <el-icon><Cpu /></el-icon>
              <span>CPU 使用情况</span>
              <el-tooltip content="基于当前采样周期的 CPU 使用率统计" placement="top">
                <el-icon class="tip-icon"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <div class="top-metrics">
            <div class="metric-box">
              <div class="metric-label">核心数</div>
              <el-progress type="dashboard" :percentage="corePercent" :color="colors">
                <template #default>
                  <div class="metric-value">{{ state.cpu?.cores || 0 }}</div>
                </template>
              </el-progress>
            </div>
            <div class="metric-box">
              <div class="metric-label">使用率</div>
              <el-progress type="dashboard" :percentage="avgCpuUsage" :color="colors" />
              <div class="metric-sub">空闲率 {{ 100 - avgCpuUsage }}%</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="section-card" shadow="never">
          <template #header>
            <div class="section-title">
              <el-icon><Coin /></el-icon>
              <span>内存使用情况</span>
              <el-tooltip content="展示系统内存占用和剩余容量" placement="top">
                <el-icon class="tip-icon"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <div class="top-metrics">
            <div class="metric-box">
              <div class="metric-label">系统内存</div>
              <el-progress type="dashboard" :percentage="state.ram?.usedPercent || 0" :color="colors" />
              <div class="metric-sub">
                {{ formatGB(state.ram?.usedMb) }}GB / {{ formatGB(state.ram?.totalMb) }}GB
              </div>
            </div>
            <div class="metric-box">
              <div class="metric-label">剩余内存</div>
              <el-progress type="dashboard" :percentage="memoryFreePercent" :color="colors" />
              <div class="metric-sub">{{ formatGB(memoryFreeMB) }}GB 可用</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="section-card mt-16" shadow="never">
      <template #header>
        <div class="section-title">
          <el-icon><Monitor /></el-icon>
          <span>服务器基础信息</span>
          <el-tooltip content="服务运行环境和进程基础信息" placement="top">
            <el-icon class="tip-icon"><QuestionFilled /></el-icon>
          </el-tooltip>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="系统平台">{{ state.os?.goos || '-' }}</el-descriptions-item>
        <el-descriptions-item label="CPU 核数">{{ state.os?.numCpu || '-' }}</el-descriptions-item>
        <el-descriptions-item label="Go 版本">{{ state.os?.goVersion || '-' }}</el-descriptions-item>
        <el-descriptions-item label="编译器">{{ state.os?.compiler || '-' }}</el-descriptions-item>
        <el-descriptions-item label="Goroutine">{{ state.os?.numGoroutine || '-' }}</el-descriptions-item>
        <el-descriptions-item label="磁盘数量">{{ (state.disk || []).length }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card class="section-card mt-16" shadow="never">
      <template #header>
        <div class="section-title">
          <el-icon><Timer /></el-icon>
          <span>运行时信息</span>
          <el-tooltip content="当前采样时刻的运行状态信息" placement="top">
            <el-icon class="tip-icon"><QuestionFilled /></el-icon>
          </el-tooltip>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="运行时">{{ state.os?.goVersion || '-' }}</el-descriptions-item>
        <el-descriptions-item label="架构">{{ state.os?.goos || '-' }}</el-descriptions-item>
        <el-descriptions-item label="当前采样 CPU">{{ avgCpuUsage }}%</el-descriptions-item>
        <el-descriptions-item label="当前采样内存">{{ state.ram?.usedPercent || 0 }}%</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card class="section-card mt-16" shadow="never">
      <template #header>
        <div class="section-title">
          <el-icon><FolderOpened /></el-icon>
          <span>磁盘使用情况</span>
          <el-tooltip content="根据配置的挂载点统计磁盘容量占用" placement="top">
            <el-icon class="tip-icon"><QuestionFilled /></el-icon>
          </el-tooltip>
        </div>
      </template>
      <el-table :data="state.disk || []" border size="small">
        <el-table-column prop="mountPoint" label="挂载路径" min-width="120" />
        <el-table-column label="文件系统" min-width="140">
          <template #default>
            未提供
          </template>
        </el-table-column>
        <el-table-column label="总容量" min-width="100">
          <template #default="{ row }">{{ row.totalGb }}GB</template>
        </el-table-column>
        <el-table-column label="已用容量" min-width="100">
          <template #default="{ row }">{{ row.usedGb }}GB</template>
        </el-table-column>
        <el-table-column label="使用率" min-width="220">
          <template #default="{ row }">
            <el-progress :percentage="row.usedPercent" :color="colors" />
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
  import { getSystemState } from '@/api/system'
  import { computed, onUnmounted, ref } from 'vue'
  import {
    Coin,
    Cpu,
    FolderOpened,
    Monitor,
    QuestionFilled,
    Timer
  } from '@element-plus/icons-vue'

  defineOptions({
    name: 'State'
  })

  const timer = ref(null)
  const state = ref({})
  const colors = ref([
    { color: '#67c23a', percentage: 40 },
    { color: '#e6a23c', percentage: 70 },
    { color: '#f56c6c', percentage: 90 }
  ])

  const avgCpuUsage = computed(() => {
    const cpus = state.value?.cpu?.cpus || []
    if (!cpus.length) return 0
    const total = cpus.reduce((sum, item) => sum + item, 0)
    return Number((total / cpus.length).toFixed(0))
  })

  const corePercent = computed(() => {
    const cores = state.value?.cpu?.cores || 0
    return Math.min(100, cores * 10)
  })
  const memoryFreeMB = computed(() => {
    const total = state.value?.ram?.totalMb || 0
    const used = state.value?.ram?.usedMb || 0
    return Math.max(total - used, 0)
  })
  const memoryFreePercent = computed(() => {
    const total = state.value?.ram?.totalMb || 0
    if (!total) return 0
    return Number(((memoryFreeMB.value / total) * 100).toFixed(0))
  })

  const formatGB = (mb) => {
    if (!mb && mb !== 0) return '-'
    return (mb / 1024).toFixed(2)
  }

  const reload = async () => {
    const { data } = await getSystemState()
    state.value = data.server || {}
  }

  reload()
  timer.value = setInterval(reload, 1000 * 10)

  onUnmounted(() => {
    clearInterval(timer.value)
    timer.value = null
  })
</script>

<style scoped>
  .state-page {
    padding: 8px;
  }

  .section-card {
    border-radius: 8px;
  }

  .section-title {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-weight: 600;
    color: var(--el-text-color-primary);
  }

  .tip-icon {
    color: var(--el-color-info);
    cursor: pointer;
    margin-left: 2px;
  }

  .top-metrics {
    display: flex;
    gap: 24px;
    justify-content: space-between;
  }

  .metric-box {
    flex: 1;
    text-align: center;
  }

  .metric-label {
    margin-bottom: 10px;
    color: var(--el-text-color-secondary);
    font-size: 13px;
  }

  .metric-value {
    font-size: 22px;
    font-weight: 700;
    line-height: 1;
  }

  .metric-sub {
    margin-top: 8px;
    color: var(--el-text-color-secondary);
    font-size: 12px;
  }

  .mt-16 {
    margin-top: 16px;
  }
</style>
