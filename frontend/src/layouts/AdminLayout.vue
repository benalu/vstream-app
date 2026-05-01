<template>
  <div class="vs-root">
    <div v-if="sidebarOpen" class="vs-overlay" @click="sidebarOpen = false" />

    <aside :class="['vs-sidebar', { 'vs-sidebar--open': sidebarOpen }]">
      <div class="vs-brand">
        <div class="vs-brand__icon">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path d="M4 4h5v5H4zM11 4h5v5h-5zM4 11h5v5H4zM11 11h5v5h-5z" fill="currentColor" opacity=".3"/>
            <path d="M7 2L2 7l5 5 5-5z" fill="currentColor"/>
          </svg>
        </div>
        <div class="vs-brand__text">
          <span class="vs-brand__name">V-STREAM</span>
          <span class="vs-brand__tag">PRO</span>
        </div>
      </div>

      <p class="vs-nav-label">Menu Utama</p>

      <nav class="vs-nav">
        <router-link to="/admin" class="vs-nav-item" exact-active-class="vs-nav-item--active">
          <span class="vs-nav-item__icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/>
              <rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/>
            </svg>
          </span>
          <span>Dashboard</span>
          <span class="vs-nav-item__dot" />
        </router-link>

        <router-link to="/admin/movies" class="vs-nav-item" active-class="vs-nav-item--active">
          <span class="vs-nav-item__icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="2" width="20" height="20" rx="2.18"/>
              <line x1="7" y1="2" x2="7" y2="22"/><line x1="17" y1="2" x2="17" y2="22"/>
              <line x1="2" y1="12" x2="22" y2="12"/><line x1="2" y1="7" x2="7" y2="7"/>
              <line x1="2" y1="17" x2="7" y2="17"/><line x1="17" y1="17" x2="22" y2="17"/>
              <line x1="17" y1="7" x2="22" y2="7"/>
            </svg>
          </span>
          <span>Kelola Film</span>
          <span class="vs-nav-item__dot" />
        </router-link>
        <router-link to="/admin/hero-slides" class="vs-nav-item" active-class="vs-nav-item--active">
          <span class="vs-nav-item__icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="7" width="20" height="15" rx="2"/>
              <path d="M16 3l-4 4-4-4"/>
            </svg>
          </span>
          <span>Hero Slider</span>
          <span class="vs-nav-item__dot" />
      </router-link>
      </nav>

      <!-- Footer: user + tombol logout -->
      <div class="vs-sidebar-footer">
        <div class="vs-user">
          <div class="vs-user__avatar">A</div>
          <div class="vs-user__info">
            <span class="vs-user__name">Administrator</span>
            <span class="vs-user__role">Super Admin</span>
          </div>
        </div>
        <button class="vs-logout-btn" @click="logout" :disabled="loggingOut" title="Logout">
          <!-- Icon logout -->
          <svg v-if="!loggingOut" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
            <polyline points="16 17 21 12 16 7"/>
            <line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
          <!-- Spinner saat proses logout -->
          <svg v-else class="vs-spin" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
          </svg>
        </button>
      </div>
    </aside>

    <div class="vs-main-wrapper">
      <header class="vs-topbar">
        <button class="vs-menu-btn" @click="sidebarOpen = !sidebarOpen">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/>
          </svg>
        </button>
        <div class="vs-topbar__right">
          <div class="vs-status-badge">
            <span class="vs-status-dot" />
            Server Normal
          </div>
        </div>
      </header>

      <main class="vs-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/lib/api'

const router = useRouter()
const sidebarOpen = ref(false)
const loggingOut = ref(false)

const logout = async () => {
  loggingOut.value = true
  try {
    await api.post('/auth/logout')
  } finally {
    loggingOut.value = false
    router.push('/admin/login')
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:wght@300;400;500;600;700&display=swap');

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

:root {
  --bg:        #09090f;
  --surface:   #0f0f1a;
  --surface2:  #14141f;
  --border:    rgba(255,255,255,0.07);
  --border2:   rgba(255,255,255,0.04);
  --accent:    #6366f1;
  --accent-lo: rgba(99,102,241,0.12);
  --accent-hi: #818cf8;
  --text:      #e2e8f0;
  --muted:     #64748b;
  --muted2:    #475569;
  --success:   #10b981;
  --warning:   #f59e0b;
  --danger:    #ef4444;
  --sidebar-w: 240px;
  --topbar-h:  56px;
  --radius:    12px;
  --radius-sm: 8px;
  --font: 'DM Sans', system-ui, sans-serif;
}

body { background: var(--bg); color: var(--text); font-family: var(--font); }

.vs-root { display: flex; min-height: 100vh; background: var(--bg); }
.vs-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.6); backdrop-filter: blur(4px); z-index: 40; }

.vs-sidebar {
  width: var(--sidebar-w); height: 100vh; position: fixed; top: 0; left: 0;
  display: flex; flex-direction: column;
  background: var(--surface); border-right: 1px solid var(--border);
  z-index: 50; transition: transform 0.3s cubic-bezier(0.4,0,0.2,1);
}

.vs-brand { display: flex; align-items: center; gap: 10px; padding: 20px 18px 16px; border-bottom: 1px solid var(--border); }
.vs-brand__icon { width: 34px; height: 34px; border-radius: 8px; background: var(--accent-lo); color: var(--accent-hi); display: grid; place-items: center; flex-shrink: 0; }
.vs-brand__name { font-weight: 700; font-size: 14px; letter-spacing: 0.06em; color: #fff; }
.vs-brand__tag { font-size: 9px; font-weight: 600; letter-spacing: 0.08em; color: var(--accent-hi); background: var(--accent-lo); border: 1px solid rgba(99,102,241,0.25); padding: 1px 5px; border-radius: 4px; }

.vs-nav-label { font-size: 10px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: var(--muted); padding: 20px 18px 8px; }
.vs-nav { display: flex; flex-direction: column; gap: 2px; padding: 0 10px; }

.vs-nav-item { display: flex; align-items: center; gap: 10px; padding: 10px 12px; border-radius: var(--radius-sm); font-size: 13.5px; font-weight: 500; color: var(--muted); text-decoration: none; transition: all 0.15s; }
.vs-nav-item__icon { width: 30px; height: 30px; display: grid; place-items: center; border-radius: 6px; flex-shrink: 0; transition: all 0.15s; }
.vs-nav-item__dot { width: 5px; height: 5px; border-radius: 50%; background: var(--accent); margin-left: auto; opacity: 0; transform: scale(0); transition: all 0.2s; }
.vs-nav-item:hover { color: var(--text); background: var(--border2); }
.vs-nav-item:hover .vs-nav-item__icon { background: var(--border); }
.vs-nav-item--active { color: #fff; background: var(--accent-lo); }
.vs-nav-item--active .vs-nav-item__icon { background: var(--accent-lo); color: var(--accent-hi); }
.vs-nav-item--active .vs-nav-item__dot { opacity: 1; transform: scale(1); }

.vs-sidebar-footer {
  margin-top: auto; padding: 14px 10px;
  border-top: 1px solid var(--border);
  display: flex; align-items: center; gap: 6px;
}
.vs-user { display: flex; align-items: center; gap: 10px; padding: 8px 10px; border-radius: var(--radius-sm); flex: 1; min-width: 0; }
.vs-user__avatar { width: 30px; height: 30px; border-radius: 6px; background: var(--accent-lo); color: var(--accent-hi); font-size: 12px; font-weight: 700; display: grid; place-items: center; flex-shrink: 0; }
.vs-user__name { display: block; font-size: 12.5px; font-weight: 600; color: var(--text); }
.vs-user__role { display: block; font-size: 11px; color: var(--muted); }

.vs-logout-btn {
  width: 32px; height: 32px; flex-shrink: 0;
  border-radius: var(--radius-sm); border: 1px solid var(--border);
  background: transparent; color: var(--muted);
  display: grid; place-items: center; cursor: pointer; transition: all 0.15s;
}
.vs-logout-btn:hover:not(:disabled) { color: var(--danger); border-color: rgba(239,68,68,0.3); background: rgba(239,68,68,0.06); }
.vs-logout-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.vs-spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.vs-main-wrapper { flex: 1; margin-left: var(--sidebar-w); display: flex; flex-direction: column; min-height: 100vh; }

.vs-topbar { height: var(--topbar-h); display: flex; align-items: center; justify-content: space-between; padding: 0 24px; border-bottom: 1px solid var(--border); background: var(--surface); position: sticky; top: 0; z-index: 30; }
.vs-menu-btn { display: none; background: none; border: none; color: var(--muted); cursor: pointer; padding: 6px; border-radius: var(--radius-sm); transition: all 0.15s; }
.vs-menu-btn:hover { color: var(--text); background: var(--border); }
.vs-topbar__right { display: flex; align-items: center; gap: 12px; margin-left: auto; }

.vs-status-badge { display: flex; align-items: center; gap: 6px; font-size: 12px; color: var(--muted); padding: 5px 10px; border-radius: 20px; background: var(--surface2); border: 1px solid var(--border); }
.vs-status-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--success); box-shadow: 0 0 6px var(--success); animation: pulse-dot 2s ease infinite; }
@keyframes pulse-dot { 0%, 100% { opacity:1; } 50% { opacity:0.5; } }

.vs-content { flex: 1; padding: 28px; max-width: 1400px; width: 100%; }

@media (max-width: 768px) {
  .vs-sidebar { transform: translateX(-100%); }
  .vs-sidebar--open { transform: translateX(0); }
  .vs-main-wrapper { margin-left: 0; }
  .vs-menu-btn { display: flex; }
  .vs-content { padding: 20px 16px; }
}

/* Shared design tokens */
.vs-page-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; margin-bottom: 28px; flex-wrap: wrap; }
.vs-page-title { font-size: 22px; font-weight: 700; color: #fff; display: flex; align-items: center; gap: 10px; letter-spacing: -0.02em; }
.vs-page-title svg { color: var(--accent-hi); }
.vs-page-subtitle { font-size: 13px; color: var(--muted); margin-top: 4px; }

.vs-btn { display: inline-flex; align-items: center; gap: 7px; padding: 9px 18px; border-radius: var(--radius-sm); font-size: 13px; font-weight: 600; cursor: pointer; border: none; transition: all 0.15s; font-family: var(--font); text-decoration: none; white-space: nowrap; }
.vs-btn--primary { background: var(--accent); color: #fff; box-shadow: 0 0 20px rgba(99,102,241,0.3); }
.vs-btn--primary:hover { background: var(--accent-hi); box-shadow: 0 0 28px rgba(99,102,241,0.45); transform: translateY(-1px); }
.vs-btn--ghost { background: var(--surface2); color: var(--muted); border: 1px solid var(--border); }
.vs-btn--ghost:hover { color: var(--text); }

.vs-card { background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius); overflow: hidden; }

.vs-input { width: 100%; background: var(--surface2); border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 9px 14px; font-size: 13px; color: var(--text); font-family: var(--font); transition: all 0.15s; outline: none; }
.vs-input:focus { border-color: rgba(99,102,241,0.5); box-shadow: 0 0 0 3px rgba(99,102,241,0.08); }
.vs-input::placeholder { color: var(--muted2); }

.vs-badge { display: inline-flex; align-items: center; gap: 4px; padding: 3px 8px; border-radius: 20px; font-size: 11px; font-weight: 600; }
.vs-badge--yellow { background: rgba(245,158,11,0.12); color: #f59e0b; }
.vs-badge--green  { background: rgba(16,185,129,0.12); color: #10b981; }
.vs-badge--slate  { background: var(--border); color: var(--muted); }

::-webkit-scrollbar { width: 5px; height: 5px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: var(--border); border-radius: 10px; }
</style>