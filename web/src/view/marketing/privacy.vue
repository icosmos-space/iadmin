<template>
  <MarketingLayout>
    <section class="doc">
      <el-tag type="info" effect="plain" round>Privacy Policy</el-tag>
      <h1>隐私政策</h1>
      <p class="meta">最后更新：2026-03-29</p>
      <el-divider />

      <h3>1. 我们收集的信息</h3>
      <p>仅收集为提供服务所必需的信息，包括账号、操作日志和安全审计相关数据。</p>

      <h3>2. 信息使用方式</h3>
      <p>用于账号管理、服务稳定性保障、故障排查以及安全合规审计。</p>

      <h3>3. 信息保护</h3>
      <p>采用权限控制、日志审计、最小化存储等措施保护数据安全。</p>

      <h3>4. 联系我们</h3>
      <p>如需数据处理说明，请发送邮件至 support@icosmos.space。</p>
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
