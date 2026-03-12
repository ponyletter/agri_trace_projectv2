import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/dashboard' },
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/login/LoginView.vue'),
      meta: { requiresAuth: false },
    },
    {
      // 公开溯源查询（消费者入口）
      path: '/trace',
      name: 'TracePublic',
      component: () => import('../views/trace/TracePublicView.vue'),
      meta: { requiresAuth: false },
    },
    {
      // 兼容旧路由 /trace/:code
      path: '/trace/:code',
      redirect: (to) => ({ path: '/trace', query: { code: to.params.code } }),
    },
    {
      // 主布局（含侧边栏）
      path: '/',
      component: () => import('../views/layout/LayoutView.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('../views/DashboardView.vue'),
          meta: { requiresAuth: true, title: '数据大屏' },
        },
        {
          path: 'batches',
          name: 'BatchList',
          component: () => import('../views/batch/BatchListView.vue'),
          meta: { requiresAuth: true, title: '批次管理' },
        },
        {
          path: 'trace-records',
          name: 'TraceRecords',
          component: () => import('../views/trace/TraceRecordsView.vue'),
          meta: { requiresAuth: true, title: '溯源记录' },
        },
        {
          path: 'planting',
          name: 'Planting',
          component: () => import('../views/planting/PlantingView.vue'),
          meta: { requiresAuth: true, title: '种植管理' },
        },
        {
          path: 'processing',
          name: 'Processing',
          component: () => import('../views/processing/ProcessingView.vue'),
          meta: { requiresAuth: true, title: '加工管理' },
        },
        {
          path: 'inspection',
          name: 'Inspection',
          component: () => import('../views/inspection/InspectionView.vue'),
          meta: { requiresAuth: true, title: '质检管理' },
        },
        {
          path: 'logistics',
          name: 'Logistics',
          component: () => import('../views/logistics/LogisticsView.vue'),
          meta: { requiresAuth: true, title: '物流追踪' },
        },
        {
          path: 'certificate',
          name: 'Certificate',
          component: () => import('../views/certificate/CertificateView.vue'),
          meta: { requiresAuth: true, title: '电子合格证' },
        },
        {
          path: 'blockchain',
          name: 'Blockchain',
          component: () => import('../views/blockchain/BlockchainView.vue'),
          meta: { requiresAuth: true, title: '区块链浏览器' },
        },
        {
          path: 'admin/users',
          name: 'AdminUsers',
          component: () => import('../views/admin/UserManageView.vue'),
          meta: { requiresAuth: true, role: 'admin', title: '用户管理' },
        },
      ],
    },
  ],
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
