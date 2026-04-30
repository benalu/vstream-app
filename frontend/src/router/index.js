// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import api from '@/lib/api'

// Lazy load semua view agar bundle lebih kecil
const Login       = () => import('@/views/admin/Login.vue')
const AdminLayout = () => import('@/layouts/AdminLayout.vue')
const Dashboard   = () => import('@/views/admin/Dashboard.vue')
const Movies      = () => import('@/views/admin/Movies.vue')

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
    meta: { requiresGuest: true }, // redirect ke dashboard jika sudah login
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
// Cek validitas session ke backend setiap navigasi ke route yang dilindungi
router.beforeEach(async (to) => {
  const needsAuth  = to.matched.some((r) => r.meta.requiresAuth)
  const needsGuest = to.matched.some((r) => r.meta.requiresGuest)

  if (!needsAuth && !needsGuest) return true // route bebas

  // Panggil /auth/me — sukses = sesi valid, 401 = belum/tidak login
  let isLoggedIn = false
  try {
    await api.get('/auth/me')
    isLoggedIn = true
  } catch {
    isLoggedIn = false
  }

  if (needsAuth && !isLoggedIn) {
    // Belum login → ke halaman login
    return { name: 'Login' }
  }

  if (needsGuest && isLoggedIn) {
    // Sudah login → tidak perlu ke halaman login lagi
    return { name: 'Dashboard' }
  }

  return true
})

export default router