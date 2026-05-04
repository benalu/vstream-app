<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import VideoPlayer   from '@/components/public/watch/VideoPlayer.vue'
import EpisodeList   from '@/components/public/watch/EpisodeList.vue'
import ReportBar     from '@/components/public/watch/ReportBar.vue'
import ServerSelector from '@/components/public/watch/ServerSelector.vue'

const route  = useRoute()
const router = useRouter()

// ── State ──────────────────────────────────────────────────────
const info       = ref(null)
const seasons    = ref([])
const subtitles  = ref([])
const loading    = ref(true)
const error      = ref(null)

const activeServer  = ref(1)
const activeSeason  = ref(1)
const activeEpisode = ref(null)
const hasError      = ref(false)

// ── Fetch ──────────────────────────────────────────────────────
const fetchWatch = async () => {
  loading.value  = true
  error.value    = null
  hasError.value = false

  try {
    const { type, tmdbId } = route.params
    const base = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api'
    const res  = await fetch(`${base}/watch/${type}/${tmdbId}`)
    const data = await res.json()
    if (!data.success) throw new Error(data.error)

    info.value      = data.data.info
    seasons.value   = data.data.seasons ?? []
    subtitles.value = data.data.subtitles ?? []

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
  const found = availableServers.value.find(s => s.n === activeServer.value)
  return found?.url ?? availableServers.value[0]?.url ?? ''
})

const storageKey = computed(() =>
  info.value ? `vstream_progress_${info.value.id}` : null
)

const TYPE_LABEL = { movie: 'Movie', series: 'Series', anime: 'Anime' }

// ── Report Error ───────────────────────────────────────────────
const reportError = async (errorType = 'load_error') => {
  hasError.value = true
  if (!info.value?.id) return
  try {
    const base = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api'
    await fetch(`${base}/logs/playback-error`, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        movie_id:   info.value.id,
        server:     `url${activeServer.value}`,
        error_type: errorType,
      }),
    })
  } catch { /* silent */ }
}

// ── Episode ────────────────────────────────────────────────────
const playerRef = ref(null)

const selectEpisode = (ep) => {
  activeEpisode.value = ep
  activeServer.value  = 1
  hasError.value      = false
  document.querySelector('.wp-player-wrap')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

// Reset error saat server bergan
watch(activeServer, () => { hasError.value = false })
</script>

<template>
  <div class="wp-root">

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

      <!-- Breadcrumb -->
      <div class="wp-breadcrumb">
        <button class="wp-bc-btn" @click="router.push('/')">Home</button>
        <span class="wp-bc-sep">›</span>
        <span class="wp-bc-current">{{ info.title }}</span>
        <template v-if="hasEpisodes && activeEpisode">
          <span class="wp-bc-sep">›</span>
          <span class="wp-bc-current">S{{ activeSeason }}E{{ activeEpisode.ep_num }}</span>
        </template>
      </div>

      <!-- Player -->
      <div class="wp-player-wrap">
        <VideoPlayer
          ref="playerRef"
          :src="currentUrl"
          :subtitles="subtitles"
          :has-episodes="hasEpisodes"
          :storage-key="storageKey"
          @error="reportError('load_error')"
        />
      </div>

      <!-- Report Bar -->
      <ReportBar
        :show="hasError"
        :active-server="activeServer"
        @report="reportError('manual')"
      />

      <!-- Meta bar -->
      <div class="wp-meta-bar">
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

        <ServerSelector
          v-model="activeServer"
          :servers="availableServers"
        />
      </div>

      <!-- Info -->
      <div class="wp-info">
        <h1 class="wp-title">{{ info.title }}</h1>
        <p v-if="info.overview" class="wp-overview">{{ info.overview }}</p>
        <div v-if="info.genre" class="wp-genres">
          <span v-for="g in info.genre.split(',')" :key="g" class="wp-genre-pill">
            {{ g.trim() }}
          </span>
        </div>
      </div>

      <!-- Episodes -->
      <EpisodeList
        v-if="hasEpisodes"
        :seasons="seasons"
        :active-episode="activeEpisode"
        v-model:active-season="activeSeason"
        @select-episode="selectEpisode"
      />

    </div>
  </div>
</template>

<style scoped>


.wp-root {
  min-height: 100vh;
  background: #07070e;
  color: #e2e8f0;
  font-family: 'DM Sans', system-ui, sans-serif;
}

.wp-loading, .wp-error {
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  min-height: 80vh; gap: 16px;
  color: #64748b; font-size: 14px;
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

.wp-page {
  padding-top: 70px;
  max-width: 1100px;
  margin: 0 auto;
  padding-left: 32px;
  padding-right: 32px;
  padding-bottom: 60px;
}

.wp-breadcrumb {
  display: flex; align-items: center; gap: 8px;
  padding: 20px 0 14px; font-size: 12.5px;
}
.wp-bc-btn {
  background: none; border: none; color: #6366f1;
  cursor: pointer; font-size: 12.5px; font-family: inherit; padding: 0;
  transition: color 0.15s;
}
.wp-bc-btn:hover { color: #818cf8; }
.wp-bc-sep  { color: #334155; }
.wp-bc-current {
  color: #64748b;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  max-width: 280px;
}

.wp-player-wrap {
  width: 100%;
  border-radius: 14px;
  overflow: hidden;
  margin-bottom: 0;
}

.wp-meta-bar {
  display: flex; align-items: center; justify-content: space-between;
  gap: 16px; padding: 14px 0;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  margin-bottom: 20px; flex-wrap: wrap;
}
.wp-meta-left { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }

.wp-type-badge {
  font-size: 9px; font-weight: 700;
  padding: 2px 8px; border-radius: 4px;
  text-transform: uppercase; letter-spacing: 0.07em; color: #fff;
}
.wp-type-badge--movie  { background: rgba(99,102,241,0.85); }
.wp-type-badge--series { background: rgba(6,182,212,0.85);  }
.wp-type-badge--anime  { background: rgba(244,63,94,0.85);  }

.wp-year, .wp-duration { font-size: 12.5px; color: #64748b; font-weight: 500; }
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

.wp-info { margin-bottom: 40px; }
.wp-title {
  font-family: 'Bebas Neue', sans-serif;
  font-size: clamp(32px, 5vw, 56px);
  letter-spacing: 0.04em; color: #fff;
  line-height: 1; margin-bottom: 12px;
}
.wp-overview {
  font-size: 13.5px; line-height: 1.7; color: #94a3b8;
  margin-bottom: 14px; max-width: 720px;
}
.wp-genres { display: flex; gap: 8px; flex-wrap: wrap; }
.wp-genre-pill {
  font-size: 11.5px; font-weight: 500; color: #94a3b8;
  padding: 3px 10px; border-radius: 20px;
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.08);
}

@media (max-width: 640px) {
  .wp-page { padding-left: 16px; padding-right: 32px; }
  .wp-meta-bar { flex-direction: column; align-items: flex-start; }
}
</style>