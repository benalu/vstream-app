<script setup>
import { onMounted, onBeforeUnmount, watch, watchEffect, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useWatchData }        from '@/composables/useWatchData'
import { useSubtitleSettings } from '@/composables/useSubtitleSettings'
import { useProgress }         from '@/composables/useProgress'
import { useSkipIntro }        from '@/composables/useSkipIntro'
import { usePlayer }           from '@/composables/usePlayer'

import WatchPlayer   from '@/components/public/watch/WatchPlayer.vue'
import WatchMeta     from '@/components/public/watch/WatchMeta.vue'
import WatchInfo     from '@/components/public/watch/WatchInfo.vue'
import WatchEpisodes from '@/components/public/watch/WatchEpisodes.vue'

const route  = useRoute()
const router = useRouter()

// ── Data ──────────────────────────────────────────────────────
const { info, seasons, subtitles, loading, error, hasEpisodes, currentSeason, fetchWatch } = useWatchData()

// ── UI State ──────────────────────────────────────────────────
const activeServer  = ref(1)
const activeSeason  = ref(1)
const activeEpisode = ref(null)
const hasError      = ref(false)

// ── Subtitle ──────────────────────────────────────────────────
const { subtitleSize, subtitleColor, subtitleStyle, increaseSubtitle, decreaseSubtitle } = useSubtitleSettings()

// ── Player ────────────────────────────────────────────────────
const { videoRef, plyrInstance, initPlyr, destroyPlayer } = usePlayer({
  subtitleSize, subtitleColor, increaseSubtitle, decreaseSubtitle,
  restoreProgress: () => progress.restoreProgress(),
  saveProgress:    () => progress.saveProgress(),
  clearProgress:   () => progress.clearProgress(),
  checkSkipIntro:  (hasEp) => skipIntro.checkSkipIntro(hasEp),
  hasEpisodes,
  reportError,
})

// ── Progress ──────────────────────────────────────────────────
const progress = useProgress(info, plyrInstance)
const { showResumePrompt, resumeFormatted, resume, startOver } = progress

// ── Skip Intro ────────────────────────────────────────────────
const skipIntro = useSkipIntro(plyrInstance)
const { showSkipIntro } = skipIntro

// ── Servers & URL ─────────────────────────────────────────────
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
  return servers.find(s => s.n === activeServer.value)?.url ?? servers[0]?.url ?? ''
})

const isDirectVideo = computed(() => {
  if (!currentUrl.value) return false
  return /\.(mp4|mkv|webm|ogg|avi|mov|m3u8)$/i.test(currentUrl.value.split('?')[0])
})

// ── Report Error ──────────────────────────────────────────────
const reportError = async (errorType = 'load_error') => {
  hasError.value = true
  if (!info.value?.id) return
  try {
    const base = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api'
    await fetch(`${base}/logs/playback-error`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ movie_id: info.value.id, server: `url${activeServer.value}`, error_type: errorType }),
    })
  } catch { /* silent */ }
}

// ── Fetch ─────────────────────────────────────────────────────
const doFetch = async () => {
  showResumePrompt.value = false
  destroyPlayer()
  await fetchWatch(route.params.type, route.params.tmdbId)
  if (seasons.value.length > 0 && seasons.value[0].episodes?.length > 0) {
    activeSeason.value  = seasons.value[0].season_num
    activeEpisode.value = seasons.value[0].episodes[0]
  }
}

onMounted(doFetch)
watch(() => route.params, doFetch)

// ── Init player saat URL ready ────────────────────────────────
watchEffect(async () => {
  if (!videoRef.value || !currentUrl.value || !isDirectVideo.value) return
  await initPlyr(subtitles.value)
})

// ── Server change ─────────────────────────────────────────────
watch(activeServer, () => {
  destroyPlayer()
  hasError.value         = false
  showResumePrompt.value = false
  showSkipIntro.value    = false
})

// ── Cleanup ───────────────────────────────────────────────────
onBeforeUnmount(() => {
  progress.saveProgress()
  destroyPlayer()
})

// ── Episode ───────────────────────────────────────────────────
const selectEpisode = (ep) => {
  progress.saveProgress()
  activeEpisode.value    = ep
  activeServer.value     = 1
  hasError.value         = false
  showResumePrompt.value = false
  showSkipIntro.value    = false
  document.querySelector('.wp-player-wrap')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}
</script>

<template>
  <div class="wp-root">

    <div v-if="loading" class="wp-loading">
      <div class="wp-spinner" />
      <p>Memuat konten...</p>
    </div>

    <div v-else-if="error" class="wp-error">
      <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <circle cx="12" cy="12" r="10"/>
        <line x1="12" y1="8" x2="12" y2="12"/>
        <line x1="12" y1="16" x2="12.01" y2="16"/>
      </svg>
      <p>{{ error }}</p>
      <button class="wp-back-btn" @click="router.back()">← Kembali</button>
    </div>

    <div v-else-if="info" class="wp-page">

      <div class="wp-breadcrumb">
        <button class="wp-bc-btn" @click="router.push('/')">Home</button>
        <span class="wp-bc-sep">›</span>
        <span class="wp-bc-current">{{ info.title }}</span>
        <template v-if="hasEpisodes && activeEpisode">
          <span class="wp-bc-sep">›</span>
          <span class="wp-bc-current">S{{ activeSeason }}E{{ activeEpisode.ep_num }}</span>
        </template>
      </div>

      <WatchPlayer
        :ref="el => videoRef = el?.$el ?? el"
        :current-url="currentUrl"
        :is-direct-video="isDirectVideo"
        :subtitles="subtitles"
        :subtitle-style="subtitleStyle"
        :show-skip-intro="showSkipIntro"
        :has-episodes="hasEpisodes"
        :show-resume-prompt="showResumePrompt"
        :resume-formatted="resumeFormatted"
        @skipIntro="skipIntro.skipIntro()"
        @resume="resume"
        @startOver="startOver"
      />

      <WatchMeta
        :info="info"
        :has-episodes="hasEpisodes"
        :active-episode="activeEpisode"
        :active-season="activeSeason"
        :available-servers="availableServers"
        :active-server="activeServer"
        :has-error="hasError"
        @update:activeServer="activeServer = $event"
        @reportError="reportError"
      />

      <WatchInfo :title="info.title" :overview="info.overview" :genre="info.genre" />

      <WatchEpisodes
        v-if="hasEpisodes"
        :seasons="seasons"
        :current-season="currentSeason(activeSeason)"
        :active-season="activeSeason"
        :active-episode="activeEpisode"
        @update:activeSeason="activeSeason = $event"
        @selectEpisode="selectEpisode"
      />

    </div>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Bebas+Neue&family=DM+Sans:wght@300;400;500;600;700&display=swap');

/* ── Base ────────────────────────────────────────────────────── */
.wp-root {
  min-height: 100vh;
  background: #07070e;
  color: #e2e8f0;
  font-family: 'DM Sans', system-ui, sans-serif;
}

/* ── Loading & Error ─────────────────────────────────────────── */
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

/* ── Page layout ─────────────────────────────────────────────── */
.wp-page {
  padding-top: 70px;
  max-width: 1100px;
  margin: 0 auto;
  padding-left: 32px;
  padding-right: 32px;
  padding-bottom: 60px;
}

/* ── Breadcrumb ──────────────────────────────────────────────── */
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
.wp-bc-sep { color: #334155; }
.wp-bc-current {
  color: #64748b;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  max-width: 280px;
}

/* ── Player wrap ─────────────────────────────────────────────── */
.wp-player-wrap {
  width: 100%;
  aspect-ratio: 16/9;
  background: #000;
  border-radius: 14px;
  overflow: hidden;
  box-shadow: 0 24px 80px rgba(0,0,0,0.7);
  margin-bottom: 0;
  position: relative;
  --subtitle-size: 16px;
}

/* ── Plyr overrides ──────────────────────────────────────────── */
:deep(.plyr) {
  width: 100%;
  height: 100%;
  border-radius: 14px;
  overflow: hidden;
  --plyr-color-main: #6366f1;
  --plyr-font-family: 'DM Sans', system-ui, sans-serif;
  --plyr-font-size-base: 13px;
  --plyr-control-radius: 6px;
  --plyr-video-background: #000;
  --plyr-range-fill-background: #6366f1;
}
:deep(.plyr--video .plyr__controls) {
  background: linear-gradient(transparent, rgba(0,0,0,0.7));
  padding: 20px 14px 12px;
}
:deep(.plyr__caption) {
  font-size: var(--subtitle-size) !important;
  color: var(--subtitle-color, #ffffff) !important;
  font-family: 'DM Sans', system-ui, sans-serif;
  font-weight: 500;
  text-shadow: 0 1px 4px rgba(0,0,0,0.9);
  background: rgba(0,0,0,0.5);
  border-radius: 4px;
  padding: 2px 8px;
  transition: font-size 0.15s ease, color 0.15s ease;
}

/* ── No URL placeholder ──────────────────────────────────────── */
.wp-no-url {
  width: 100%; height: 100%;
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  gap: 12px; color: #334155; font-size: 13px;
}

/* ── Skip Intro ──────────────────────────────────────────────── */
.wp-skip-intro {
  position: absolute;
  bottom: 80px;
  right: 20px;
  z-index: 10;
  padding: 8px 18px;
  border-radius: 4px;
  border: 2px solid rgba(255,255,255,0.8);
  background: rgba(0,0,0,0.65);
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  backdrop-filter: blur(4px);
  transition: all 0.2s;
}
.wp-skip-intro:hover {
  background: rgba(255,255,255,0.15);
  border-color: #fff;
}

/* ── Resume Prompt ───────────────────────────────────────────── */
.wp-resume-prompt {
  position: absolute;
  bottom: 80px;
  left: 20px;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  border-radius: 8px;
  background: rgba(0,0,0,0.75);
  border: 1px solid rgba(255,255,255,0.12);
  backdrop-filter: blur(8px);
  flex-wrap: wrap;
}
.wp-resume-text {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #fff;
  font-size: 13px;
  font-weight: 500;
}
.wp-resume-btn {
  padding: 4px 12px;
  border-radius: 5px;
  border: 1px solid rgba(255,255,255,0.2);
  background: transparent;
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s;
}
.wp-resume-btn:hover { background: rgba(255,255,255,0.1); }
.wp-resume-btn--primary {
  background: #6366f1;
  border-color: transparent;
}
.wp-resume-btn--primary:hover { background: #818cf8; }

/* ── Transitions ─────────────────────────────────────────────── */
.wp-fade-slide-enter-active,
.wp-fade-slide-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.wp-fade-slide-enter-from,
.wp-fade-slide-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

/* ── Report bar ──────────────────────────────────────────────── */
.wp-report-bar {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  padding: 6px 0;
}
.wp-report-label { font-size: 11.5px; color: #475569; }
.wp-report-btn {
  font-size: 11.5px; font-weight: 600; color: #ef4444;
  background: rgba(239,68,68,0.08);
  border: 1px solid rgba(239,68,68,0.2);
  border-radius: 6px; padding: 4px 10px;
  cursor: pointer; font-family: inherit; transition: all 0.15s;
}
.wp-report-btn:hover { background: rgba(239,68,68,0.15); }

/* ── Custom subtitle size panel di dalam Plyr settings ───────── */
:deep(.plyr__menu__container) {
  min-width: 200px;
}

:deep(.plyr-subtitle-size-panel) {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px 14px 16px;
}

:deep(.plyr-sub-size-row) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  background: rgba(0,0,0,0.25);
  border-radius: 8px;
  padding: 6px 8px;
}

:deep(.plyr-sub-size-btn) {
  background: rgba(99,102,241,0.2);
  border: 1px solid rgba(99,102,241,0.35);
  border-radius: 6px;
  color: #a5b4fc;
  font-size: 12px;
  font-weight: 700;
  padding: 5px 14px;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s;
  flex-shrink: 0;
  /* HAPUS width: 100% — supaya tidak full width lagi */
}

:deep(.plyr-sub-size-btn:hover) {
  background: rgba(99,102,241,0.4);
  border-color: rgba(99,102,241,0.7);
  color: #fff;
}

:deep(.plyr-sub-size-display) {
  font-size: 15px;
  font-weight: 700;
  color: #fff;
  min-width: 42px;
  text-align: center;
  letter-spacing: -0.01em;
}

:deep(.plyr-sub-size-slider-wrap) {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 0 2px;
}

:deep(.plyr-sub-size-min),
:deep(.plyr-sub-size-max) {
  font-size: 10px;
  color: #64748b;
  white-space: nowrap;
  flex-shrink: 0;
}

:deep(.plyr-sub-size-slider) {
  flex: 1;
  accent-color: #6366f1;
  cursor: pointer;
  height: 4px;
  min-width: 0; /* fix overflow di container sempit */
}

:deep(.plyr__menu__container) {
  min-width: 200px;
  background: #1a1a2e !important;
  border: 1px solid rgba(99,102,241,0.2);
}

/* Paksa semua teks dalam menu jadi terang */
:deep(.plyr__menu__container *) {
  color: #c7c8ff;
}

:deep(.plyr-sub-size-label) {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  padding: 8px 14px 4px;
  border-top: 1px solid rgba(255,255,255,0.08);
  margin-top: 4px;
}

:deep(.plyr-sub-color-row) {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 14px 12px;
}

:deep(.plyr-sub-color-btn) {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  padding: 0;
  transition: all 0.15s;
  flex-shrink: 0;
}

:deep(.plyr-sub-color-btn:hover) {
  transform: scale(1.2);
}

:deep(.plyr-sub-color-btn--active) {
  border-color: #fff;
  box-shadow: 0 0 0 2px rgba(99,102,241,0.6);
  transform: scale(1.15);
}

/* ── Meta bar ────────────────────────────────────────────────── */
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
.wp-servers { display: flex; align-items: center; gap: 6px; }
.wp-servers-label { font-size: 11.5px; color: #64748b; font-weight: 500; }
.wp-server-btn {
  width: 32px; height: 32px; border-radius: 7px;
  border: 1px solid rgba(255,255,255,0.08);
  background: rgba(255,255,255,0.04); color: #64748b;
  font-size: 12px; font-weight: 700;
  cursor: pointer; font-family: inherit; transition: all 0.15s;
}
.wp-server-btn:hover { color: #fff; border-color: rgba(255,255,255,0.15); }
.wp-server-btn--active {
  background: #6366f1; border-color: transparent; color: #fff;
  box-shadow: 0 0 14px rgba(99,102,241,0.4);
}

/* ── Info ────────────────────────────────────────────────────── */
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

/* ── Episodes ────────────────────────────────────────────────── */
.wp-episodes { margin-top: 8px; }
.wp-ep-header {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 20px; gap: 16px; flex-wrap: wrap;
}
.wp-ep-header-left { display: flex; align-items: center; gap: 10px; }
.wp-ep-bar {
  width: 4px; height: 20px; border-radius: 2px;
  background: #6366f1; box-shadow: 0 0 10px rgba(99,102,241,0.6); flex-shrink: 0;
}
.wp-ep-title { font-size: 17px; font-weight: 700; color: #fff; }
.wp-season-tabs { display: flex; gap: 6px; flex-wrap: wrap; }
.wp-season-tab {
  padding: 5px 14px; border-radius: 20px;
  border: 1px solid rgba(255,255,255,0.08);
  background: transparent; color: #64748b;
  font-size: 12.5px; font-weight: 600;
  cursor: pointer; font-family: inherit; transition: all 0.15s;
}
.wp-season-tab:hover { color: #fff; }
.wp-season-tab--active {
  background: rgba(99,102,241,0.15);
  border-color: rgba(99,102,241,0.4); color: #818cf8;
}
.wp-ep-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 10px;
}
.wp-ep-card {
  display: flex; align-items: center; gap: 10px;
  padding: 12px 14px; border-radius: 10px;
  border: 1px solid rgba(255,255,255,0.06);
  background: rgba(255,255,255,0.02); color: #64748b;
  cursor: pointer; text-align: left; font-family: inherit; transition: all 0.15s;
}
.wp-ep-card:hover {
  background: rgba(255,255,255,0.05);
  border-color: rgba(255,255,255,0.12); color: #fff;
  transform: translateY(-2px);
}
.wp-ep-card--active {
  background: rgba(99,102,241,0.1);
  border-color: rgba(99,102,241,0.3); color: #fff;
}
.wp-ep-card-num {
  width: 30px; height: 30px; border-radius: 7px;
  background: rgba(255,255,255,0.06);
  display: grid; place-items: center;
  font-size: 12px; font-weight: 700; flex-shrink: 0; color: #94a3b8;
}
.wp-ep-card--active .wp-ep-card-num { background: rgba(99,102,241,0.25); color: #818cf8; }
.wp-ep-card-info { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.wp-ep-card-label { font-size: 11px; font-weight: 600; color: #64748b; }
.wp-ep-card--active .wp-ep-card-label { color: #818cf8; }
.wp-ep-card-title {
  font-size: 12px; font-weight: 500; color: #fff;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}

/* ── Responsive ──────────────────────────────────────────────── */
@media (max-width: 640px) {
  .wp-page { padding-left: 16px; padding-right: 16px; }
  .wp-ep-grid { grid-template-columns: repeat(auto-fill, minmax(130px, 1fr)); }
  .wp-meta-bar { flex-direction: column; align-items: flex-start; }
  .wp-skip-intro { bottom: 70px; right: 12px; font-size: 12px; padding: 6px 14px; }
  .wp-resume-prompt { bottom: 70px; left: 12px; }
  .wp-sub-control { gap: 8px; }
  .wp-sub-slider { max-width: 100px; }
}
</style>