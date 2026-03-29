<template>
  <MarketingLayout>
    <section class="top">
      <el-tag type="warning" effect="plain" round>Pricing</el-tag>
      <h1>按团队阶段选择投入，按业务增长平滑升级</h1>
      <p>示意套餐如下，可根据部署规模、响应等级和上线节奏进行定制。</p>
    </section>

    <section class="plans">
      <el-row :gutter="14">
        <el-col v-for="plan in plans" :key="plan.name" :xs="24" :sm="12" :md="8">
          <article class="plan" :class="{ featured: plan.featured }">
            <el-space alignment="center">
              <h3>{{ plan.name }}</h3>
              <el-tag v-if="plan.featured" type="success" round>推荐</el-tag>
            </el-space>
            <p class="price">{{ plan.price }}</p>
            <p class="desc">{{ plan.desc }}</p>
            <ul>
              <li v-for="item in plan.items" :key="item">{{ item }}</li>
            </ul>
            <RouterLink to="/contact">
              <el-button :type="plan.featured ? 'primary' : undefined" plain>
                {{ plan.featured ? '预约演示' : '咨询详情' }}
              </el-button>
            </RouterLink>
          </article>
        </el-col>
      </el-row>
    </section>
  </MarketingLayout>
</template>

<script setup>
  import MarketingLayout from '@/view/marketing/layout.vue'

  const plans = [
    {
      name: 'Starter',
      price: '￥0 / 月',
      desc: '适合内部试用与小团队验证',
      items: ['基础权限管理', '登录与菜单体系', '社区支持']
    },
    {
      name: 'Growth',
      price: '￥2,999 / 月',
      desc: '适合快速增长中的业务团队',
      featured: true,
      items: ['插件能力扩展', '日志审计与告警', '工作日技术支持']
    },
    {
      name: 'Enterprise',
      price: '定制报价',
      desc: '适合高并发和合规要求场景',
      items: ['专属架构方案', '私有化部署支持', 'SLA 与响应保障']
    }
  ]
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

  .plans {
    padding: 18px 6px;
    border-top: 1px solid rgba(37, 99, 235, 0.13);
  }

  .plan {
    height: 100%;
    padding: 18px;
    border-radius: 14px;
    background: rgba(255, 255, 255, 0.66);
    border: 1px solid rgba(37, 99, 235, 0.12);
  }

  .plan.featured {
    border-color: rgba(37, 99, 235, 0.3);
    box-shadow: 0 10px 24px rgba(37, 99, 235, 0.14);
  }

  .plan h3 {
    margin: 0;
    color: #1f3a68;
    font-size: 24px;
  }

  .price {
    margin: 10px 0 0;
    color: #1d4ed8;
    font-size: 34px;
    font-weight: 700;
  }

  .desc {
    margin: 8px 0 0;
    color: #4f74ad;
    line-height: 1.65;
  }

  ul {
    list-style: none;
    margin: 14px 0 18px;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  li {
    color: #5f7fb3;
    line-height: 1.5;
  }
</style>
