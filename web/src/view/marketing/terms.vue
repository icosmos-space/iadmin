<template>
  <MarketingLayout>
    <section class="doc">
      <el-tag type="info" effect="plain" round>Terms of Service</el-tag>
      <h1>服务条款</h1>
      <p class="meta">最后更新：2026-03-29</p>
      <el-divider />

      <h3>1. 使用范围</h3>
      <p>你可以在遵守法律法规和本条款的前提下使用 iAdmin 及相关服务。</p>

      <h3>2. 账号责任</h3>
      <p>账号使用方需对其账号操作行为负责，并妥善保管凭据和访问权限。</p>

      <h3>3. 服务变更</h3>
      <p>我们可能根据产品演进调整服务内容，并提前在页面或文档中说明。</p>

      <h3>4. 免责说明</h3>
      <p>对因不可抗力或第三方原因导致的服务中断，我们将在合理范围内协助恢复。</p>
    </section>
  </MarketingLayout>
</template>

<script setup>
  import MarketingLayout from '@/view/marketing/layout.vue'
</script>

<style scoped lang="scss">
  .doc {
    padding: 16px 8px;
    border-top: 1px solid rgba(37, 99, 235, 0.13);
    color: #2a4b7e;
  }

  h1 {
    margin: 12px 0 4px;
    font-size: clamp(34px, 4.2vw, 50px);
    color: #123f8d;
  }

  .meta {
    margin: 0;
    color: #6a8bc0;
    font-size: 13px;
  }

  h3 {
    margin: 14px 0 8px;
    font-size: 19px;
    color: #1f3a68;
  }

  p {
    margin: 0;
    color: #5f7fb3;
    line-height: 1.75;
  }
</style>
