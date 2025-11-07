import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  // Public routes
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { requiresAuth: false, hideLayout: true }
  },
  {
    path: '/signup',
    name: 'Signup',
    component: () => import('@/views/SignupView.vue'),
    meta: { requiresAuth: false, hideLayout: true }
  },

  // Protected routes
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/DashboardView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/accounts',
    name: 'Accounts',
    component: () => import('@/views/AccountsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/account-types',
    name: 'AccountTypes',
    component: () => import('@/views/AccountTypesView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/transactions',
    name: 'Transactions',
    component: () => import('@/views/TransactionsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/budgets',
    name: 'Budgets',
    component: () => import('@/views/BudgetsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/credit-cards',
    name: 'CreditCards',
    component: () => import('@/views/CreditCardsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/goals',
    name: 'Goals',
    component: () => import('@/views/GoalsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/reports',
    name: 'Reports',
    component: () => import('@/views/ReportsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/bills',
    name: 'Bills',
    component: () => import('@/views/BillsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/scheduled-transactions',
    name: 'ScheduledTransactions',
    component: () => import('@/views/ScheduledTransactionsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/SettingsView.vue'),
    meta: { requiresAuth: true }
  },
  // 404 catch-all route
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/dashboard'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Error handling for lazy-loaded components
router.onError((error) => {
  console.error('Router error:', error)
  if (error.message.includes('Failed to fetch dynamically imported module') ||
      error.message.includes('Importing a module script failed')) {
    // Chunk loading error - reload the page to get fresh chunks
    console.warn('Chunk loading failed, reloading page...')
    window.location.reload()
  }
})

// Track if auth has been initialized
let authInitialized = false

// Navigation guard
router.beforeEach(async (to, from, next) => {
  try {
    const authStore = useAuthStore()

    // Initialize auth once on first navigation
    if (!authInitialized) {
      await authStore.initializeAuth()
      authInitialized = true
    }

    const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false)
    const isAuthPage = to.path === '/login' || to.path === '/signup'

    if (requiresAuth && !authStore.isAuthenticated) {
      // Redirect to login if trying to access protected route
      next({ name: 'Login', query: { redirect: to.fullPath } })
    } else if (isAuthPage && authStore.isAuthenticated) {
      // Redirect to dashboard if already logged in and trying to access auth pages
      next({ name: 'Dashboard' })
    } else {
      next()
    }
  } catch (error) {
    console.error('Navigation guard error:', error)
    // Allow navigation to continue even if auth check fails
    next()
  }
})

export default router
