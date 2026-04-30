<script setup>
import { useScrolled } from '@/composables/useScrolled'

const { scrolled } = useScrolled(60)
</script>

<template>
  <nav :class="['vs-nav', { 'vs-nav--scrolled': scrolled }]">
    <!-- Brand -->
    <div class="vs-nav__brand">
      <div class="vs-nav__logo">
        <svg width="18" height="18" viewBox="0 0 20 20" fill="none">
          <path d="M4 4h5v5H4zM11 4h5v5h-5zM4 11h5v5H4zM11 11h5v5h-5z" fill="currentColor" opacity=".35"/>
          <path d="M7 2L2 7l5 5 5-5z" fill="currentColor"/>
        </svg>
      </div>
      <span class="vs-nav__name">V-STREAM</span>
    </div>

    <!-- Links -->
    <div class="vs-nav__links">
      <a href="#" class="vs-nav__link vs-nav__link--active">Home</a>
      <a href="#" class="vs-nav__link">Movies</a>
      <a href="#" class="vs-nav__link">Series</a>
      <a href="#" class="vs-nav__link">Anime</a>
    </div>

    <!-- Right actions -->
    <div class="vs-nav__right">
      <button class="vs-nav__search-btn" aria-label="Cari">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
        </svg>
      </button>
      <router-link to="/admin" class="vs-nav__admin-btn">Admin</router-link>
    </div>
  </nav>
</template>

<style scoped>
.vs-nav {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  gap: 40px;
  padding: 20px 48px;
  transition: background 0.4s, padding 0.3s, backdrop-filter 0.4s;
}
.vs-nav--scrolled {
  background: rgba(7, 7, 14, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  padding: 14px 48px;
  border-bottom: 1px solid rgba(255,255,255,0.07);
}

/* Brand */
.vs-nav__brand {
  display: flex;
  align-items: center;
  gap: 9px;
  flex-shrink: 0;
  text-decoration: none;
}
.vs-nav__logo {
  width: 34px; height: 34px;
  border-radius: 8px;
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  display: grid; place-items: center;
  color: #fff;
}
.vs-nav__name {
  font-family: 'Bebas Neue', sans-serif;
  font-size: 22px;
  letter-spacing: 0.12em;
  color: #fff;
}

/* Links */
.vs-nav__links {
  display: flex;
  gap: 28px;
  flex: 1;
}
.vs-nav__link {
  font-size: 13.5px;
  font-weight: 500;
  color: #64748b;
  text-decoration: none;
  position: relative;
  transition: color 0.15s;
}
.vs-nav__link:hover { color: #fff; }
.vs-nav__link--active { color: #fff; }
.vs-nav__link--active::after {
  content: '';
  position: absolute;
  bottom: -4px; left: 0; right: 0;
  height: 2px;
  background: #6366f1;
  border-radius: 2px;
}

/* Right */
.vs-nav__right {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}
.vs-nav__search-btn {
  width: 36px; height: 36px;
  border-radius: 8px;
  border: 1px solid rgba(255,255,255,0.07);
  background: rgba(255,255,255,0.04);
  color: #64748b;
  display: grid; place-items: center;
  cursor: pointer;
  transition: all 0.15s;
}
.vs-nav__search-btn:hover { color: #fff; border-color: rgba(255,255,255,0.15); }

.vs-nav__admin-btn {
  padding: 8px 18px;
  border-radius: 8px;
  background: #6366f1;
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  text-decoration: none;
  box-shadow: 0 0 20px rgba(99,102,241,0.35);
  transition: background 0.15s, transform 0.15s;
}
.vs-nav__admin-btn:hover { background: #818cf8; transform: translateY(-1px); }

@media (max-width: 768px) {
  .vs-nav { padding: 16px 20px; }
  .vs-nav--scrolled { padding: 12px 20px; }
  .vs-nav__links { display: none; }
}
</style>