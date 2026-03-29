<template>
  <MarketingLayout>
    <section class="top">
      <el-tag type="primary" effect="plain" round>Feature Stack</el-tag>
      <h1>为高效率团队准备的能力栈</h1>
      <p>权限、路由、代码生成、日志审计与插件扩展，覆盖后台交付全流程。</p>
    </section>

    <section class="list">
      <el-row :gutter="14">
        <el-col v-for="item in featureList" :key="item.title" :xs="24" :sm="12" :md="8">
          <article class="item">
            <el-tag size="small" type="info" effect="plain">模块能力</el-tag>
            <h3>{{ item.title }}</h3>
            <p>{{ item.desc }}</p>
          </article>
        </el-col>
      </el-row>
    </section>
  </MarketingLayout>
</template>

<script setup>
  import MarketingLayout from '@/view/marketing/layout.vue'

  const featureList = [
    { title: '角色与菜单权限', desc: '可视化分配菜单与按钮权限，支持细粒度访问控制。' },
    { title: '动态路由加载', desc: '权限驱动路由注入，兼顾安全、性能与可维护性。' },
    { title: '系统参数中心', desc: '把分散配置统一管理，减少多环境差异造成的问题。' },
    { title: '审计日志体系', desc: '登录日志、操作日志、错误记录可追踪可回放。' },
    { title: '多插件架构', desc: '业务能力按模块扩展，避免核心系统过度耦合。' },
    { title: 'AI 助手接入', desc: '后台内直接使用 AI 能力提升运营与研发效率。' }
  ]
</script>

<style scoped lang="scss">
  .top {
    padding: 10px 6px 18px;
  }

  .top h1 {
    margin: 12px 0 0;
    font-size: clamp(36px, 4.5vw, 52px);
    color: #123f8d;
    line-height: 1.08;
  }

  .top p {
    margin: 14px 0 0;
    max-width: 780px;
    color: #4f74ad;
    line-height: 1.75;
  }

  .list {
    padding: 18px 6px;
    border-top: 1px solid rgba(37, 99, 235, 0.13);
  }

  .item {
    height: 100%;
    padding: 18px;
    border-radius: 14px;
    background: rgba(255, 255, 255, 0.65);
    border: 1px solid rgba(37, 99, 235, 0.12);
  }

  .item h3 {
    margin: 12px 0 0;
    color: #1f3a68;
    font-size: 22px;
  }

  .item p {
    margin: 10px 0 0;
    color: #5f7fb3;
    line-height: 1.65;
  }
</style>
