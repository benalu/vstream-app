// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import api from '@/lib/api'

// Lazy load semua view agar bundle lebih kecil
const Login       = () => import('@/views/admin/Login.vue')
const AdminLayout = () => import('@/layouts/AdminLayout.vue')
const Dashboard   = () => import('@/views/admin/Dashboard.vue')
const Movies      = () => import('@/views/admin/Movies.vue')

// Cache hasil /auth/me agar tidak hit backend di setiap navigasi.
// null  = belum pernah dicek
// true  = sesi valid
// false = tidak login
let sessionCache = null

// Dipanggil oleh api.js interceptor saat 401, agar guard re-fetch di navigasi berikutnya
export const clearSessionCache = () => {
  sessionCache = null
}

const routes = [
  // ── Redirect root ke dashboard ───────────────────────────
  {
    path: '/',
    redirect: '/admin',
  },

  // ── Login (tidak butuh auth) ─────────────────────────────
  {
    path: '/admin/login',
    name: 'Login',
    component: Login,
    meta: { requiresGuest: true },
  },

  // ── Admin area (butuh auth) ──────────────────────────────
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: Dashboard,
      },
      {
        path: 'movies',
        name: 'Movies',
        component: Movies,
      },
    ],
  },

  // ── 404 fallback ─────────────────────────────────────────
  {
    path: '/:pathMatch(.*)*',
    redirect: '/admin',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// ── Navigation Guard ─────────────────────────────────────────
// Cek validitas session. Hasil di-cache agar tidak hit /auth/me
// di setiap navigasi — cache di-reset oleh api.js saat menerima 401.
router.beforeEach(async (to) => {
  const needsAuth  = to.matched.some((r) => r.meta.requiresAuth)
  const needsGuest = to.matched.some((r) => r.meta.requiresGuest)

  if (!needsAuth && !needsGuest) return true

  // Gunakan cache jika sudah ada hasil sebelumnya
  let isLoggedIn = sessionCache
  if (isLoggedIn === null) {
    try {
      await api.get('/auth/me')
      isLoggedIn = true
    } catch {
      isLoggedIn = false
    }
    sessionCache = isLoggedIn
  }

  if (needsAuth && !isLoggedIn) return { name: 'Login' }
  if (needsGuest && isLoggedIn) return { name: 'Dashboard' }

  return true
})

export default router