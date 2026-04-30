<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import TheNavbar from '@/components/public/TheNavbar.vue'

const route  = useRoute()
const router = useRouter()

// ── State ─────────────────────────────────────────────────────
const info       = ref(null)
const seasons    = ref([])
const loading    = ref(true)
const error      = ref(null)

const activeServer  = ref(1)
const activeSeason  = ref(1)
const activeEpisode = ref(null)

// ── Fetch ──────────────────────────────────────────────────────
const fetchWatch = async () => {
  loading.value = true
  error.value   = null
  try {
    const { type, tmdbId } = route.params
    const base = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api'
    const res  = await fetch(`${base}/watch/${type}/${tmdbId}`)
    const data = await res.json()
    if (!data.success) throw new Error(data.error)

    info.value    = data.data.info
    seasons.value = data.data.seasons ?? []

    // Default: pilih episode pertama season pertama
    if (seasons.value.length > 0 && seasons.value[0].episodes?.length > 0) {
      activeSeason.value  = seasons.value[0].season_num
      activeEpisode.value = seasons.value[0].episodes[0]
    }
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

onMounted(fetchWatch)
watch(() => route.params, fetchWatch)

// ── Computed ───────────────────────────────────────────────────
const hasEpisodes = computed(() =>
  info.value?.type === 'series' ||
  (info.value?.type === 'anime' && info.value?.has_episodes)
)

const currentSeason = computed(() =>
  seasons.value.find(s => s.season_num === activeSeason.value)
)

const availableServers = computed(() => {
  if (hasEpisodes.value && activeEpisode.value) {
    return [
      { n: 1, url: activeEpisode.value.url1 },
      { n: 2, url: activeEpisode.value.url2 },
    ].filter(s => s.url)
  }
  if (!info.value) return []
  return [
    { n: 1, url: info.value.url1 },
    { n: 2, url: info.value.url2 },
    { n: 3, url: info.value.url3 },
  ].filter(s => s.url)
})

const currentUrl = computed(() => {
  const servers = availableServers.value
  const found   = servers.find(s => s.n === activeServer.value)
  return found?.url ?? servers[0]?.url ?? ''
})

const selectEpisode = (ep) => {
  activeEpisode.value = ep
  activeServer.value  = 1
  // Scroll ke player
  document.querySelector('.wp-player-wrap')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

const TYPE_LABEL = { movie: 'Movie', series: 'Series', anime: 'Anime' }
</script>

<template>
  <div class="wp-root">
    <TheNavbar />

    <!-- Loading -->
    <div v-if="loading" class="wp-loading">
      <div class="wp-spinner" />
      <p>Memuat konten...</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="wp-error">
      <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <circle cx="12" cy="12" r="10"/>
        <line x1="12" y1="8" x2="12" y2="12"/>
        <line x1="12" y1="16" x2="12.01" y2="16"/>
      </svg>
      <p>{{ error }}</p>
      <button class="wp-back-btn" @click="router.back()">← Kembali</button>
    </div>

    <!-- Main -->
    <div v-else-if="info" class="wp-page">

      <!-- ── Breadcrumb ─────────────────────────────────── -->
      <div class="wp-breadcrumb">
        <button class="wp-bc-btn" @click="router.push('/')">Home</button>
        <span class="wp-bc-sep">›</span>
        <span class="wp-bc-current">{{ info.title }}</span>
        <span v-if="hasEpisodes && activeEpisode" class="wp-bc-sep">›</span>
        <span v-if="hasEpisodes && activeEpisode" class="wp-bc-current">
          S{{ activeSeason }}E{{ activeEpisode.ep_num }}
        </span>
      </div>

      <!-- ── Player wrap ────────────────────────────────── -->
      <div class="wp-player-wrap">
        <iframe
          v-if="currentUrl"
          :key="currentUrl"
          :src="currentUrl"
          class="wp-iframe"
          allowfullscreen
          allow="autoplay; fullscreen"
          frameborder="0"
          sandbox="allow-scripts allow-same-origin allow-presentation allow-forms"
        />
        <div v-else class="wp-no-url">
          <svg width="52" height="52" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2">
            <polygon points="5 3 19 12 5 21 5 3" opacity=".3"/>
          </svg>
          <p>URL stream belum tersedia</p>
        </div>
      </div>

      <!-- ── Server selector ───────────────────────────── -->
      <div class="wp-meta-bar">
        <!-- Kiri: info singkat -->
        <div class="wp-meta-left">
          <span class="wp-type-badge" :class="`wp-type-badge--${info.type}`">
            {{ TYPE_LABEL[info.type] }}
          </span>
          <span class="wp-year">{{ info.year }}</span>
          <span v-if="info.rating" class="wp-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#f59e0b">
              <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
            </svg>
            {{ info.rating }}
          </span>
          <span v-if="info.duration" class="wp-duration">{{ info.duration }}</span>
          <span v-if="hasEpisodes && activeEpisode" class="wp-ep-now-playing">
            S{{ activeSeason }} E{{ activeEpisode.ep_num }}
            <span v-if="activeEpisode.title"> · {{ activeEpisode.title }}</span>
          </span>
        </div>

        <!-- Kanan: server buttons -->
        <div v-if="availableServers.length > 0" class="wp-servers">
          <span class="wp-servers-label">Server:</span>
          <button
            v-for="s in availableServers"
            :key="s.n"
            class="wp-server-btn"
            :class="{ 'wp-server-btn--active': activeServer === s.n }"
            @click="activeServer = s.n"
          >
            {{ s.n }}
          </button>
        </div>
      </div>

      <!-- ── Info section ───────────────────────────────── -->
      <div class="wp-info">
        <h1 class="wp-title">{{ info.title }}</h1>
        <p v-if="info.overview" class="wp-overview">{{ info.overview }}</p>
        <div v-if="info.genre" class="wp-genres">
          <span v-for="g in info.genre.split(',')" :key="g" class="wp-genre-pill">
            {{ g.trim() }}
          </span>
        </div>
      </div>

      <!-- ── Episode list (series/anime) ───────────────── -->
      <div v-if="hasEpisodes && seasons.length > 0" class="wp-episodes">

        <!-- Section header + season tabs -->
        <div class="wp-ep-header">
          <div class="wp-ep-header-left">
            <div class="wp-ep-bar" />
            <h2 class="wp-ep-title">Episode</h2>
          </div>
          <div v-if="seasons.length > 1" class="wp-season-tabs">
            <button
              v-for="s in seasons"
              :key="s.season_num"
              class="wp-season-tab"
              :class="{ 'wp-season-tab--active': activeSeason === s.season_num }"
              @click="activeSeason = s.season_num"
            >
              Season {{ s.season_num }}
            </button>
          </div>
        </div>

        <!-- Episode grid -->
        <div class="wp-ep-grid">
          <button
            v-for="ep in currentSeason?.episodes ?? []"
            :key="ep.id"
            class="wp-ep-card"
            :class="{ 'wp-ep-card--active': activeEpisode?.id === ep.id }"
            @click="selectEpisode(ep)"
          >
            <!-- Play indicator -->
            <div class="wp-ep-card-num">
              <svg v-if="activeEpisode?.id === ep.id"
                width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                <polygon points="5 3 19 12 5 21 5 3"/>
              </svg>
              <span v-else>{{ ep.ep_num }}</span>
            </div>
            <div class="wp-ep-card-info">
              <span class="wp-ep-card-label">Episode {{ ep.ep_num }}</span>
              <span v-if="ep.title" class="wp-ep-card-title">{{ ep.title }}</span>
            </div>
          </button>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Bebas+Neue&family=DM+Sans:wght@300;400;500;600;700&display=swap');

.wp-root {
  min-height: 100vh;
  background: #07070e;
  color: #e2e8f0;
  font-family: 'DM Sans', system-ui, sans-serif;
}

/* ── Loading / Error ─────────────────────────────────── */
.wp-loading, .wp-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 80vh;
  gap: 16px;
  color: #64748b;
  font-size: 14px;
}
.wp-spinner {
  width: 36px; height: 36px;
  border: 3px solid rgba(99,102,241,0.2);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
.wp-back-btn {
  padding: 8px 18px; border-radius: 8px;
  border: 1px solid rgba(255,255,255,0.1);
  background: transparent; color: #94a3b8;
  font-size: 13px; cursor: pointer;
  transition: all 0.15s; font-family: inherit;
}
.wp-back-btn:hover { color: #fff; border-color: rgba(255,255,255,0.2); }

/* ── Page wrapper ────────────────────────────────────── */
.wp-page {
  padding-top: 70px;
  max-width: 1100px;
  margin: 0 auto;
  padding-left: 32px;
  padding-right: 32px;
  padding-bottom: 60px;
}

/* ── Breadcrumb ──────────────────────────────────────── */
.wp-breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px 0 14px;
  font-size: 12.5px;
}
.wp-bc-btn {
  background: none; border: none;
  color: #6366f1; cursor: pointer;
  font-size: 12.5px; font-family: inherit; padding: 0;
  transition: color 0.15s;
}
.wp-bc-btn:hover { color: #818cf8; }
.wp-bc-sep { color: #334155; }
.wp-bc-current {
  color: #64748b;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  max-width: 280px;
}

/* ── Player ──────────────────────────────────────────── */
.wp-player-wrap {
  width: 100%;
  aspect-ratio: 16/9;
  background: #000;
  border-radius: 14px;
  overflow: hidden;
  box-shadow: 0 24px 80px rgba(0,0,0,0.7);
  margin-bottom: 0;
}
.wp-iframe { width: 100%; height: 100%; border: none; display: block; }
.wp-no-url {
  width: 100%; height: 100%;
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  gap: 12px; color: #334155; font-size: 13px;
}

/* ── Meta bar (below player) ─────────────────────────── */
.wp-meta-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 0;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  margin-bottom: 20px;
  flex-wrap: wrap;
}
.wp-meta-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.wp-type-badge {
  font-size: 9px; font-weight: 700;
  padding: 2px 8px; border-radius: 4px;
  text-transform: uppercase; letter-spacing: 0.07em; color: #fff;
}
.wp-type-badge--movie  { background: rgba(99,102,241,0.85); }
.wp-type-badge--series { background: rgba(6,182,212,0.85);  }
.wp-type-badge--anime  { background: rgba(244,63,94,0.85);  }
.wp-year, .wp-duration {
  font-size: 12.5px; color: #64748b; font-weight: 500;
}
.wp-rating {
  display: flex; align-items: center; gap: 4px;
  font-size: 12.5px; font-weight: 700; color: #f59e0b;
}
.wp-ep-now-playing {
  font-size: 12px; color: #818cf8; font-weight: 600;
  background: rgba(99,102,241,0.1);
  border: 1px solid rgba(99,102,241,0.2);
  padding: 2px 9px; border-radius: 20px;
}

/* Server buttons */
.wp-servers {
  display: flex; align-items: center; gap: 6px;
}
.wp-servers-label {
  font-size: 11.5px; color: #64748b; font-weight: 500;
}
.wp-server-btn {
  width: 32px; height: 32px;
  border-radius: 7px;
  border: 1px solid rgba(255,255,255,0.08);
  background: rgba(255,255,255,0.04);
  color: #64748b;
  font-size: 12px; font-weight: 700;
  cursor: pointer; font-family: inherit;
  transition: all 0.15s;
}
.wp-server-btn:hover { color: #fff; border-color: rgba(255,255,255,0.15); }
.wp-server-btn--active {
  background: #6366f1;
  border-color: transparent; color: #fff;
  box-shadow: 0 0 14px rgba(99,102,241,0.4);
}

/* ── Info ────────────────────────────────────────────── */
.wp-info { margin-bottom: 40px; }
.wp-title {
  font-family: 'Bebas Neue', sans-serif;
  font-size: clamp(32px, 5vw, 56px);
  letter-spacing: 0.04em; color: #fff;
  line-height: 1; margin-bottom: 12px;
}
.wp-overview {
  font-size: 13.5px; line-height: 1.7;
  color: #94a3b8; margin-bottom: 14px;
  max-width: 720px;
}
.wp-genres { display: flex; gap: 8px; flex-wrap: wrap; }
.wp-genre-pill {
  font-size: 11.5px; font-weight: 500; color: #94a3b8;
  padding: 3px 10px; border-radius: 20px;
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.08);
}

/* ── Episodes ────────────────────────────────────────── */
.wp-episodes { margin-top: 8px; }

.wp-ep-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  gap: 16px;
  flex-wrap: wrap;
}
.wp-ep-header-left {
  display: flex; align-items: center; gap: 10px;
}
.wp-ep-bar {
  width: 4px; height: 20px; border-radius: 2px;
  background: #6366f1;
  box-shadow: 0 0 10px rgba(99,102,241,0.6);
  flex-shrink: 0;
}
.wp-ep-title {
  font-size: 17px; font-weight: 700; color: #fff;
}

/* Season tabs */
.wp-season-tabs { display: flex; gap: 6px; flex-wrap: wrap; }
.wp-season-tab {
  padding: 5px 14px; border-radius: 20px;
  border: 1px solid rgba(255,255,255,0.08);
  background: transparent; color: #64748b;
  font-size: 12.5px; font-weight: 600;
  cursor: pointer; font-family: inherit;
  transition: all 0.15s;
}
.wp-season-tab:hover { color: #fff; }
.wp-season-tab--active {
  background: rgba(99,102,241,0.15);
  border-color: rgba(99,102,241,0.4);
  color: #818cf8;
}

/* Episode grid */
.wp-ep-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 10px;
}
.wp-ep-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid rgba(255,255,255,0.06);
  background: rgba(255,255,255,0.02);
  color: #64748b;
  cursor: pointer; text-align: left;
  font-family: inherit;
  transition: all 0.15s;
}
.wp-ep-card:hover {
  background: rgba(255,255,255,0.05);
  border-color: rgba(255,255,255,0.12);
  color: #fff;
  transform: translateY(-2px);
}
.wp-ep-card--active {
  background: rgba(99,102,241,0.1);
  border-color: rgba(99,102,241,0.3);
  color: #fff;
}
.wp-ep-card-num {
  width: 30px; height: 30px;
  border-radius: 7px;
  background: rgba(255,255,255,0.06);
  display: grid; place-items: center;
  font-size: 12px; font-weight: 700;
  flex-shrink: 0; color: #94a3b8;
}
.wp-ep-card--active .wp-ep-card-num {
  background: rgba(99,102,241,0.25);
  color: #818cf8;
}
.wp-ep-card-info {
  display: flex; flex-direction: column; gap: 2px;
  min-width: 0;
}
.wp-ep-card-label {
  font-size: 11px; font-weight: 600; color: #64748b;
}
.wp-ep-card--active .wp-ep-card-label { color: #818cf8; }
.wp-ep-card-title {
  font-size: 12px; font-weight: 500; color: #fff;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}

/* ── Responsive ──────────────────────────────────────── */
@media (max-width: 640px) {
  .wp-page { padding-left: 16px; padding-right: 16px; }
  .wp-ep-grid { grid-template-columns: repeat(auto-fill, minmax(130px, 1fr)); }
  .wp-meta-bar { flex-direction: column; align-items: flex-start; }
}
</style>