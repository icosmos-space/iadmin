<template>
  <MarketingLayout>
    <section class="hero">
      <el-row :gutter="26" align="middle">
        <el-col :xs="24" :md="15">
          <el-space wrap class="hero-tags">
            <el-tag type="primary" effect="dark" round>Enterprise Console</el-tag>
            <el-tag effect="plain" round>Golang + Vue + AI</el-tag>
          </el-space>

          <h1>企业后台的一体化蓝图</h1>
          <p class="hero-desc">
            用统一权限、动态路由、插件体系与 AI 能力，支撑从 0 到 1 到规模化的后台交付。
            少做重复搭建，把时间留给核心业务。
          </p>

          <el-space wrap class="hero-actions">
            <RouterLink to="/features"><el-button type="primary" size="large">查看能力</el-button></RouterLink>
            <RouterLink to="/pricing"><el-button size="large" plain>查看方案</el-button></RouterLink>
            <RouterLink to="/contact"><el-button size="large" text bg>预约演示</el-button></RouterLink>
          </el-space>

          <div class="stats">
            <article v-for="item in heroStats" :key="item.label" class="stat-item">
              <el-statistic :value="item.value" :suffix="item.suffix" />
              <p>{{ item.label }}</p>
            </article>
          </div>
        </el-col>

        <el-col :xs="24" :md="9">
          <div class="hero-side">
            <h3>典型交付节奏</h3>
            <el-timeline>
              <el-timeline-item v-for="item in timeline" :key="item.title" :timestamp="item.time" type="primary">
                <p class="step-title">{{ item.title }}</p>
                <p class="step-desc">{{ item.desc }}</p>
              </el-timeline-item>
            </el-timeline>
          </div>
          <p class="hero-side-note">从“搭框架”升级到“造业务”</p>
        </el-col>
      </el-row>
    </section>

    <section class="section">
      <div class="section-head">
        <el-tag type="primary" effect="plain" round>Capabilities</el-tag>
        <h2>核心能力矩阵</h2>
      </div>

      <el-row :gutter="14" class="cap-grid">
        <el-col v-for="item in capabilities" :key="item.title" :xs="24" :sm="12" :md="8">
          <article class="cap-item">
            <el-icon class="cap-icon"><component :is="item.icon" /></el-icon>
            <h3>{{ item.title }}</h3>
            <p>{{ item.desc }}</p>
          </article>
        </el-col>
      </el-row>
    </section>

    <section class="section">
      <div class="section-head">
        <el-tag type="success" effect="plain" round>Use Cases</el-tag>
        <h2>典型应用场景</h2>
      </div>

      <el-row :gutter="14" class="case-grid">
        <el-col v-for="item in showcases" :key="item.title" :xs="24" :md="8">
          <article class="case-item">
            <h3>{{ item.title }}</h3>
            <p>{{ item.desc }}</p>
            <el-link :underline="false" type="primary" href="#">{{ item.cta }}</el-link>
          </article>
        </el-col>
      </el-row>
    </section>

    <section class="section trust">
      <div class="section-head">
        <el-tag type="info" effect="plain" round>Trusted By Teams</el-tag>
        <h2>多行业团队实践</h2>
      </div>
      <el-space wrap class="brand-list">
        <el-tag v-for="brand in brands" :key="brand" round effect="plain" type="info">{{ brand }}</el-tag>
      </el-space>
    </section>

    <section class="cta">
      <h2>把后台工程变成可复制的组织能力</h2>
      <p>现在就开始体验 iAdmin，或者联系我们获取面向你业务场景的落地方案。</p>
      <el-space wrap>
        <RouterLink to="/login"><el-button type="primary" size="large">立即体验</el-button></RouterLink>
        <RouterLink to="/contact"><el-button size="large" plain>联系我们</el-button></RouterLink>
      </el-space>
    </section>
  </MarketingLayout>
</template>

<script setup>
  import {
    Cpu,
    DataAnalysis,
    Connection,
    MagicStick,
    Setting,
    Promotion
  } from '@element-plus/icons-vue'

  import MarketingLayout from '@/view/marketing/layout.vue'

  const heroStats = [
    { value: 3, suffix: 'x', label: '交付效率提升' },
    { value: 99.9, suffix: '%', label: '可用性目标' },
    { value: 1, suffix: ' Day', label: '新业务落地' },
    { value: 50, suffix: '+', label: '可扩展模块' }
  ]

  const timeline = [
    { time: 'T+0', title: '初始化项目', desc: '开箱即用能力就绪，快速完成基础配置。' },
    { time: 'T+1', title: '接入业务模型', desc: '自动生成 CRUD 与权限路由，减少重复编码。' },
    { time: 'T+3', title: '发布上线', desc: '通过日志审计与监控面板保障稳定交付。' }
  ]

  const capabilities = [
    { title: '统一权限治理', desc: '菜单、接口、按钮、数据权限集中管控。', icon: Connection },
    { title: '智能代码生成', desc: '从数据结构到页面接口，一键生成业务基础模块。', icon: MagicStick },
    { title: 'AI 助手协作', desc: '用 AI 提升配置效率、排障速度和运营效率。', icon: Cpu },
    { title: '插件化架构', desc: '按需扩展功能，保持核心系统稳定轻量。', icon: Promotion },
    { title: '可观测与审计', desc: '登录、操作、错误日志全链路可追踪。', icon: DataAnalysis },
    { title: '工程化最佳实践', desc: '标准化分层与结构，提升团队协作上限。', icon: Setting }
  ]

  const showcases = [
    { title: '企业中后台', desc: '组织、用户、权限、流程等核心管理场景。', cta: '查看解决方案' },
    { title: '运营活动系统', desc: '活动配置、资源管理、数据统计一体化支撑。', cta: '查看能力清单' },
    { title: '多租户 SaaS', desc: '策略路由与模块扩展支持多业务线独立演进。', cta: '咨询架构建议' }
  ]

  const brands = ['FinTech Team', 'Retail Cloud', 'Gov Platform', 'SaaS Startup', 'AI Ops Group', 'E-Commerce Lab']
</script>

<style scoped lang="scss">
  .hero {
    padding: 14px 6px 24px;
  }

  .hero h1 {
    margin: 14px 0 0;
    font-size: clamp(38px, 6vw, 62px);
    line-height: 1.05;
    letter-spacing: -0.03em;
    color: #113f8c;
  }

  .hero-desc {
    margin: 16px 0 0;
    max-width: 92%;
    color: #35568c;
    line-height: 1.8;
    font-size: 16px;
  }

  .hero-actions {
    margin-top: 22px;
  }

  .stats {
    margin-top: 24px;
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 12px;
  }

  .stat-item {
    padding: 14px;
    border-radius: 14px;
    border: 1px solid rgba(37, 99, 235, 0.16);
    background: rgba(255, 255, 255, 0.65);
  }

  :deep(.stat-item .el-statistic__head),
  :deep(.stat-item .el-statistic__content) {
    color: #1d4ed8;
  }

  :deep(.stat-item .el-statistic__number) {
    font-size: 24px;
  }

  .stat-item p {
    margin: 8px 0 0;
    color: #5f7fb3;
    font-size: 13px;
  }

  .hero-side {
    padding: 16px 16px 4px;
    border-radius: 16px;
    border: 1px solid rgba(37, 99, 235, 0.15);
    background: rgba(255, 255, 255, 0.66);
  }

  .hero-side h3 {
    margin: 0;
    color: #1d4ed8;
    font-size: 20px;
  }

  .step-title {
    margin: 0;
    color: #1f3a68;
    font-weight: 600;
  }

  .step-desc {
    margin: 4px 0 0;
    color: #5f7fb3;
    font-size: 13px;
    line-height: 1.5;
  }

  .hero-side-note {
    margin: 12px 4px 0;
    color: #2563eb;
    font-weight: 600;
  }

  .section {
    padding: 18px 6px;
    border-top: 1px solid rgba(37, 99, 235, 0.13);
  }

  .section-head h2 {
    margin: 12px 0 0;
    font-size: clamp(30px, 4vw, 42px);
    color: #123f8d;
  }

  .cap-grid,
  .case-grid {
    margin-top: 14px;
  }

  .cap-item,
  .case-item {
    height: 100%;
    padding: 18px;
    border-radius: 14px;
    background: rgba(255, 255, 255, 0.64);
    border: 1px solid rgba(37, 99, 235, 0.12);
  }

  .cap-icon {
    width: 42px;
    height: 42px;
    border-radius: 10px;
    background: rgba(59, 130, 246, 0.14);
    color: #2563eb;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 21px;
  }

  .cap-item h3,
  .case-item h3 {
    margin: 14px 0 0;
    color: #1f3a68;
    font-size: 22px;
  }

  .cap-item p,
  .case-item p {
    margin: 10px 0 0;
    color: #5f7fb3;
    line-height: 1.65;
  }

  .trust {
    padding-bottom: 12px;
  }

  .brand-list {
    margin-top: 14px;
    row-gap: 10px;
  }

  .brand-list :deep(.el-tag) {
    border-color: rgba(37, 99, 235, 0.2);
    color: #1e40af;
    background: rgba(255, 255, 255, 0.72);
  }

  .cta {
    margin-top: 8px;
    padding: 28px 10px 8px;
    text-align: center;
    border-top: 1px solid rgba(37, 99, 235, 0.13);
  }

  .cta h2 {
    margin: 0;
    font-size: clamp(34px, 4.5vw, 50px);
    color: #11408f;
    line-height: 1.12;
  }

  .cta p {
    margin: 14px auto 22px;
    max-width: 760px;
    color: #4f74ad;
    line-height: 1.75;
  }

  @media (max-width: 980px) {
    .hero-desc {
      max-width: 100%;
    }

    .stats {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
  }
</style>
