import { createRouter, createWebHistory } from 'vue-router'
import api from '@/lib/api'

const Home         = () => import('@/views/public/Home.vue')
const Watch        = () => import('@/views/public/Watch.vue')
const Login        = () => import('@/views/admin/Login.vue')
const PublicLayout = () => import('@/layouts/PublicLayout.vue')
const AdminLayout  = () => import('@/layouts/AdminLayout.vue')
const Dashboard    = () => import('@/views/admin/Dashboard.vue')
const Movies       = () => import('@/views/admin/Movies.vue')
const HeroSlides   = () => import('@/views/admin/HeroSlides.vue')
const NotFound     = () => import('@/views/NotFound.vue')

let sessionCache = null
export const clearSessionCache = () => { sessionCache = null }

const routes = [
  {
    path: '/',
    component: PublicLayout,
    children: [
      { path: '',                    name: 'Home',  component: Home  },
      { path: 'watch/:type/:tmdbId', name: 'Watch', component: Watch },
    ],
  },
  {
    path: '/admin/login',
    name: 'Login',
    component: Login,
    meta: { requiresGuest: true },
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '',            name: 'Dashboard', component: Dashboard  },
      { path: 'movies',      name: 'Movies',    component: Movies     },
      { path: 'hero-slides', name: 'Hero',      component: HeroSlides },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to) => {
  const needsAuth  = to.matched.some((r) => r.meta.requiresAuth)
  const needsGuest = to.matched.some((r) => r.meta.requiresGuest)

  if (!needsAuth && !needsGuest) return true

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

  if (needsAuth  && !isLoggedIn) return { name: 'Login'     }
  if (needsGuest &&  isLoggedIn) return { name: 'Dashboard' }

  return true
})

export default router