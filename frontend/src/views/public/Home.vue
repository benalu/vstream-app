<script setup>
import { ref, computed, onMounted } from 'vue'
import TheNavbar      from '@/components/public/TheNavbar.vue'
import HeroSection    from '@/components/public/HeroSection.vue'
import CategoryFilter from '@/components/public/CategoryFilter.vue'
import Top10Row       from '@/components/public/Top10Row.vue'
import ContentGrid    from '@/components/public/ContentGrid.vue'
import TheFooter      from '@/components/public/TheFooter.vue'
import { useContent } from '@/composables/useContent'

// ── Data ─────────────────────────────────────────────────────
const {
  featured,
  top10,
  recentlyAdded,
  loading,
  loadingFeatured,
  error,
  fetchAll,
} = useContent()

onMounted(fetchAll)

// ── Category filter ───────────────────────────────────────────
const CATEGORIES = [
  { id: 'all',    label: 'Semua'  },
  { id: 'movie',  label: 'Movies' },
  { id: 'series', label: 'Series' },
  { id: 'anime',  label: 'Anime'  },
]
const activeCategory = ref('all')

const filteredRecent = computed(() => {
  if (activeCategory.value === 'all') return recentlyAdded.value
  return recentlyAdded.value.filter(i => i.type === activeCategory.value)
})

const sectionTitle = computed(() => ({
  all: 'Baru Ditambahkan', movie: 'Movies', series: 'Series', anime: 'Anime',
}[activeCategory.value]))

// ── Handlers ──────────────────────────────────────────────────
const handlePlay      = (item) => console.log('play:', item)       // TODO: navigate /watch
const handleMore      = (item) => console.log('more:', item)       // TODO: navigate /detail
const handleCardClick = (item) => console.log('card:', item)
const handleMoreClick = ()     => console.log('lihat semua')
</script>

<template>
  <div class="home">
    <TheNavbar />

    <!-- Hero: loading -->
    <div v-if="loadingFeatured" class="home__hero-skeleton" />

    <!-- Hero: data nyata -->
    <HeroSection
      v-else-if="featured"
      :content="featured"
      @play="handlePlay"
      @more="handleMore"
    />

    <!-- Hero: tidak ada konten -->
    <div v-else class="home__hero-empty">
      <p>Belum ada konten. Tambahkan melalui
        <router-link to="/admin">Admin Dashboard</router-link>.
      </p>
    </div>

    <main class="home__main">
      <!-- Error banner -->
      <div v-if="error" class="home__error">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="8" x2="12" y2="12"/>
          <line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        Gagal memuat sebagian konten. Periksa koneksi dan coba lagi.
      </div>

      <!-- Category filter -->
      <CategoryFilter v-model="activeCategory" :categories="CATEGORIES" />

      <!-- TOP 10 -->
      <section class="home__section">
        <div class="home__section-header">
          <div class="home__section-bar" />
          <h2 class="home__section-title">TOP 10 Hari Ini</h2>
        </div>

        <div v-if="loading && top10.length === 0" class="home__top10-skeleton">
          <div v-for="i in 7" :key="i" class="home__skeleton-card" />
        </div>
        <Top10Row
          v-else-if="top10.length > 0"
          :items="top10"
          @click="handleCardClick"
        />
        <p v-else class="home__empty-text">Belum ada konten.</p>
      </section>

      <!-- Recently added / filtered -->
      <ContentGrid
        :title="sectionTitle"
        :items="filteredRecent"
        :loading="loading"
        @cardClick="handleCardClick"
        @moreClick="handleMoreClick"
      />
    </main>

    <TheFooter />
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Bebas+Neue&family=DM+Sans:wght@300;400;500;600;700&display=swap');

:root {
  --pub-bg:        #07070e;
  --pub-surface:   #0d0d18;
  --pub-border:    rgba(255,255,255,0.07);
  --pub-accent:    #6366f1;
  --pub-accent-hi: #818cf8;
  --pub-text:      #e2e8f0;
  --pub-muted:     #64748b;
}
</style>

<style scoped>
.home {
  min-height: 100vh;
  background: var(--pub-bg);
  color: var(--pub-text);
  font-family: 'DM Sans', system-ui, sans-serif;
}

/* Hero states */
.home__hero-skeleton {
  height: 100svh;
  min-height: 600px;
  max-height: 900px;
  background: linear-gradient(90deg, #0d0d18 25%, #12121f 50%, #0d0d18 75%);
  background-size: 200% 100%;
  animation: shimmer 1.6s infinite;
}
.home__hero-empty {
  height: 50vh;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: var(--pub-muted);
}
.home__hero-empty a { color: var(--pub-accent-hi); }

/* Error */
.home__error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 8px;
  background: rgba(239,68,68,0.08);
  border: 1px solid rgba(239,68,68,0.2);
  color: #f87171;
  font-size: 12.5px;
}

/* Main layout */
.home__main {
  padding: 40px 48px 60px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 52px;
}

/* Section header (TOP 10) */
.home__section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 22px;
}
.home__section-bar {
  width: 4px; height: 22px;
  border-radius: 2px;
  background: var(--pub-accent);
  box-shadow: 0 0 10px rgba(99,102,241,0.6);
  flex-shrink: 0;
}
.home__section-title {
  font-size: 18px; font-weight: 700;
  color: #fff; letter-spacing: -0.01em;
}

/* TOP 10 skeleton */
.home__top10-skeleton {
  display: flex;
  gap: 8px;
  overflow: hidden;
}
.home__skeleton-card {
  width: 190px; height: 190px;
  flex-shrink: 0;
  border-radius: 10px;
  background: linear-gradient(90deg, #0d0d18 25%, #12121f 50%, #0d0d18 75%);
  background-size: 200% 100%;
  animation: shimmer 1.6s infinite;
}

@keyframes shimmer { to { background-position: -200% 0; } }

.home__empty-text {
  color: var(--pub-muted);
  font-size: 13px;
  padding: 24px 0;
}

@media (max-width: 768px) {
  .home__main { padding: 32px 20px 48px; gap: 40px; }
}
</style>