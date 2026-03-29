import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/landing',
    name: 'Landing',
    meta: {
      title: '首页宣传',
      public: true
    },
    component: () => import('@/view/marketing/landing.vue')
  },
  {
    path: '/features',
    name: 'Features',
    meta: {
      title: '功能介绍',
      public: true
    },
    component: () => import('@/view/marketing/features.vue')
  },
  {
    path: '/pricing',
    name: 'Pricing',
    meta: {
      title: '价格方案',
      public: true
    },
    component: () => import('@/view/marketing/pricing.vue')
  },
  {
    path: '/contact',
    name: 'Contact',
    meta: {
      title: '联系咨询',
      public: true
    },
    component: () => import('@/view/marketing/contact.vue')
  },
  {
    path: '/privacy',
    name: 'Privacy',
    meta: {
      title: '隐私政策',
      public: true
    },
    component: () => import('@/view/marketing/privacy.vue')
  },
  {
    path: '/terms',
    name: 'Terms',
    meta: {
      title: '服务条款',
      public: true
    },
    component: () => import('@/view/marketing/terms.vue')
  },
  {
    path: '/init',
    name: 'Init',
    meta: {
      public: true
    },
    component: () => import('@/view/init/index.vue')
  },
  {
    path: '/login',
    name: 'Login',
    meta: {
      public: true
    },
    component: () => import('@/view/login/index.vue')
  },
  {
    path: '/scanUpload',
    name: 'ScanUpload',
    meta: {
      title: '扫码上传',
      client: true
    },
    component: () => import('@/view/example/upload/scanUpload.vue')
  },
  {
    path: '/:catchAll(.*)',
    meta: {
      closeTab: true
    },
    component: () => import('@/view/error/index.vue')
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
