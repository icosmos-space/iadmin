<template>
  <MarketingLayout>
    <section class="top">
      <el-tag type="primary" effect="plain" round>Contact</el-tag>
      <h1>告诉我们目标，我们协助你完成落地</h1>
      <p>可通过邮箱或 GitHub Issue 联系，我们通常在 1 个工作日内响应。</p>
    </section>

    <section class="contact-wrap">
      <el-row :gutter="14">
        <el-col :xs="24" :md="12">
          <article class="block">
            <h3>联系通道</h3>
            <el-space direction="vertical" alignment="start" fill>
              <el-link href="mailto:support@icosmos.space" type="primary">support@icosmos.space</el-link>
              <el-link href="https://github.com/icosmos-space/iadmin" target="_blank" type="primary">
                github.com/icosmos-space/iadmin
              </el-link>
              <p>服务时间：工作日 10:00 - 18:00 (UTC+8)</p>
            </el-space>
          </article>
        </el-col>

        <el-col :xs="24" :md="12">
          <article class="block">
            <h3>建议沟通信息</h3>
            <ul>
              <li>团队规模和角色分布</li>
              <li>预计部署方式（公有云/私有化）</li>
              <li>核心场景与性能目标</li>
              <li>需要接入的现有系统</li>
            </ul>
            <RouterLink to="/login">
              <el-button type="primary" plain>立即体验后台</el-button>
            </RouterLink>
          </article>
        </el-col>
      </el-row>
    </section>
  </MarketingLayout>
</template>

<script setup>
  import MarketingLayout from '@/view/marketing/layout.vue'
</script>

<style scoped lang="scss">
  .top {
    padding: 10px 6px 18px;
  }

  .top h1 {
    margin: 12px 0 0;
    font-size: clamp(34px, 4.2vw, 50px);
    color: #123f8d;
    line-height: 1.1;
  }

  .top p {
    margin: 14px 0 0;
    max-width: 780px;
    color: #4f74ad;
    line-height: 1.75;
  }

  .contact-wrap {
    padding: 18px 6px;
    border-top: 1px solid rgba(37, 99, 235, 0.13);
  }

  .block {
    height: 100%;
    padding: 18px;
    border-radius: 14px;
    background: rgba(255, 255, 255, 0.66);
    border: 1px solid rgba(37, 99, 235, 0.12);
  }

  .block h3 {
    margin: 0 0 12px;
    font-size: 22px;
    color: #1f3a68;
  }

  p,
  li {
    color: #5f7fb3;
    line-height: 1.6;
  }

  ul {
    margin: 0 0 18px;
    padding-left: 18px;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
</style>
