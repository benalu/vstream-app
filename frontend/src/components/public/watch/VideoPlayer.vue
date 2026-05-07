<script setup>
import { ref, computed, watch, nextTick, watchEffect, onBeforeUnmount } from 'vue'
import { usePlayer }   from '@/composables/usePlayer'
import { useSubtitle } from '@/composables/useSubtitle'
import { useProgress } from '@/composables/useProgress'

const props = defineProps({
  src:         { type: String,  default: '' },
  subtitles:   { type: Array,   default: () => [] },
  hasEpisodes: { type: Boolean, default: false },
  storageKey:  { type: String,  default: null },
  poster:      { type: String,  default: '' },
})

const emit = defineEmits(['error', 'timeupdate', 'ended'])

// ── Refs ───────────────────────────────────────────────────────
const videoRef     = ref(null)
const hasError     = ref(false)
const hasStarted   = ref(false)

// ── Subtitle ───────────────────────────────────────────────────
const {
  subtitleSize,
  subtitleColor,
  subtitleStyle,
  increaseSubtitle,
  decreaseSubtitle,
  setSubtitleColor,
  SUBTITLE_COLORS,
} = useSubtitle()

// ── Skip Intro ─────────────────────────────────────────────────
const showSkipIntro = ref(false)
const INTRO_START   = 30
const INTRO_END     = 90

const skipIntro = () => {
  if (plyrInstance.value) plyrInstance.value.currentTime = INTRO_END
  showSkipIntro.value = false
}

// ── Progress ───────────────────────────────────────────────────
const storageKeyRef = computed(() => props.storageKey)

const {
  showResumePrompt,
  resumeFormatted,
  saveProgress,
  restoreProgress,
  resume,
  startOver,
  clearProgress,
} = useProgress(storageKeyRef, ref(null))

// ── Player ─────────────────────────────────────────────────────
const { plyrInstance, initPlyr, destroyPlyr } = usePlayer({
  subtitleSize,
  subtitleColor,
  increaseSubtitle,
  decreaseSubtitle,
  setSubtitleColor,
  SUBTITLE_COLORS,

  onReady: () => {
    restoreProgress()
  },

  onTimeUpdate: () => {
    if (!hasStarted.value) hasStarted.value = true
    saveProgress()
    if (props.hasEpisodes) {
      const t = plyrInstance.value?.currentTime ?? 0
      showSkipIntro.value = t >= INTRO_START && t <= INTRO_END
    }
    emit('timeupdate')
  },

  onEnded: () => {
    showSkipIntro.value = false
    clearProgress()
    emit('ended')
  },

  onError: () => {
    hasError.value = true
    emit('error')
  },


})

// ── Computed ───────────────────────────────────────────────────
const isDirectVideo = computed(() => {
  if (!props.src) return false
  const url = props.src.toLowerCase().split('?')[0]
  return /\.(mp4|mkv|webm|ogg|avi|mov|m3u8)$/.test(url)
})

// ── Init / destroy Plyr on src change ─────────────────────────
watchEffect(async () => {
  const url = props.src
  const ok  = isDirectVideo.value
  if (!url || !ok || !videoRef.value) return
  await nextTick()
  initPlyr(videoRef.value)
})

// Reset state on src change
watch(() => props.src, () => {
  hasError.value      = false
  hasStarted.value    = false
  showSkipIntro.value = false
})

onBeforeUnmount(() => {
  saveProgress()
  destroyPlyr()
})

const posterSrcset = computed(() => {
  const url = props.poster
  if (!url) return ''
  // Replace w1280 → ukuran lain untuk srcset
  const w300  = url.replace('/w1280/', '/w300/')
  const w500  = url.replace('/w1280/', '/w500/')
  const w780  = url.replace('/w1280/', '/w780/')
  return `${w300} 300w, ${w500} 500w, ${w780} 780w, ${url} 1280w`
})

defineExpose({ plyrInstance })
</script>

<template>
  <div class="vp-wrap" :style="subtitleStyle">

    <!-- ── Direct video: Plyr ─────────────────────────────── -->
    <template v-if="src && isDirectVideo">
      <video
        ref="videoRef"
        :key="src"
        :src="src"
        class="vp-video"
        preload="metadata"
        playsinline
      >
        <track
          v-for="sub in subtitles"
          :key="sub.lang"
          :src="sub.url"
          :srclang="sub.lang"
          :label="sub.label"
          kind="subtitles"
          :default="sub.lang === 'id'"
        />
      </video>
    </template>

    <!-- ── Iframe embed ───────────────────────────────────── -->
    <iframe
      v-else-if="src && !isDirectVideo"
      :key="'iframe-' + src"
      :src="src"
      class="vp-iframe"
      allowfullscreen
      allow="autoplay; fullscreen"
      frameborder="0"
      sandbox="allow-scripts allow-same-origin allow-presentation allow-forms"
    />

    <!-- ── No URL placeholder ────────────────────────────── -->
    <div v-else class="vp-no-url">
      <svg width="52" height="52" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2">
        <polygon points="5 3 19 12 5 21 5 3" opacity=".3"/>
      </svg>
      <p>URL stream belum tersedia</p>
    </div>

    <!-- ── Poster Overlay ─────────────────────────────────── -->
    <Transition name="vp-fade">
      <div v-if="poster && !hasStarted" class="vp-poster">
        <img
          :src="poster"
          :srcset="posterSrcset"
          sizes="(max-width: 375px) 375px,
                (max-width: 640px) 640px,
                (max-width: 1024px) 780px,
                1280px"
          class="vp-poster__img"
          fetchpriority="high"
          alt=""
          draggable="false"
        />
        <div class="vp-poster__gradient" />
        <button class="vp-poster__play" aria-label="Putar film" @click="plyrInstance?.play()">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 3 19 12 5 21 5 3"/>
          </svg>
        </button>
      </div>
    </Transition>

    <!-- ── Skip Intro ─────────────────────────────────────── -->
    <Transition name="vp-fade-slide">
      <button
        v-if="showSkipIntro && hasEpisodes"
        class="vp-skip-intro"
        @click="skipIntro"
      >
        Skip Intro →
      </button>
    </Transition>

    <!-- ── Resume Prompt ──────────────────────────────────── -->
    <Transition name="vp-fade-slide">
      <div v-if="showResumePrompt" class="vp-resume-prompt">
        <svg width="11" height="11" viewBox="0 0 24 24" fill="currentColor" class="vp-resume-icon">
          <polygon points="5 3 19 12 5 21 5 3"/>
        </svg>
        <span class="vp-resume-text">Lanjut dari {{ resumeFormatted }}?</span>
        <div class="vp-resume-actions">
          <button class="vp-resume-btn vp-resume-btn--primary" @click="resume">Lanjutkan</button>
          <button class="vp-resume-btn" @click="startOver">Dari Awal</button>
        </div>
      </div>
    </Transition>

  </div>
</template>

<style scoped>
/* ── Wrapper ─────────────────────────────────────────────────── */
.vp-wrap {
  width: 100%;
  aspect-ratio: 16 / 9;
  background: #000;
  border-radius: 14px;
  overflow: hidden;
  box-shadow: 0 24px 80px rgba(0, 0, 0, 0.7);
  position: relative;
  --subtitle-size: 16px;
  --subtitle-color: #ffffff;
}

/* ── Video / Iframe ──────────────────────────────────────────── */
.vp-video,
.vp-iframe {
  width: 100%;
  height: 100%;
  display: block;
  border: none;
}

/* ── No URL ──────────────────────────────────────────────────── */
.vp-no-url {
  width: 100%; height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #334155;
  font-size: 13px;
}

/* ── Skip Intro ──────────────────────────────────────────────── */
.vp-skip-intro {
  position: absolute;
  bottom: 80px; right: 20px;
  z-index: 10;
  padding: 8px 18px;
  border-radius: 4px;
  border: 2px solid rgba(255, 255, 255, 0.8);
  background: rgba(0, 0, 0, 0.65);
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  backdrop-filter: blur(4px);
  transition: all 0.2s;
}
.vp-skip-intro:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: #fff;
}

/* ── Resume Prompt ───────────────────────────────────────────── */
.vp-resume-prompt {
  position: absolute;
  bottom: 56px;
  left: 12px;
  transform: none;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 10px;
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.82);
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(8px);
  white-space: nowrap;
  max-width: calc(100% - 24px);
}

.vp-resume-icon {
  color: rgba(255, 255, 255, 0.5);
  flex-shrink: 0;
}

.vp-resume-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.8);
  font-weight: 500;
  flex-shrink: 0;
}

.vp-resume-actions {
  display: flex;
  gap: 5px;
  flex-shrink: 0;
}

.vp-resume-btn {
  padding: 3px 10px;
  border-radius: 5px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  background: transparent;
  color: #fff;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s;
  white-space: nowrap;
}
.vp-resume-btn:hover { background: rgba(255, 255, 255, 0.1); }
.vp-resume-btn--primary {
  background: #6366f1;
  border-color: transparent;
}
.vp-resume-btn--primary:hover { background: #818cf8; }

/* Mobile */
@media (max-width: 640px) {
  .vp-resume-prompt {
    bottom: 18px;
    left: 8px;
    gap: 6px;
    padding: 6px 8px;
  }
  .vp-resume-text {
    font-size: 11px;
  }
  .vp-resume-btn {
    padding: 3px 8px;
    font-size: 10px;
  }
}

/* ── Plyr overrides ──────────────────────────────────────────── */
:deep(.plyr) {
  width: 100%; height: 100%;
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
  display: flex;
  flex-direction: column;
  padding: 0 14px 12px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
  gap: 0;
}

:deep(.plyr__controls__row--top) {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 0 2px;
}

:deep(.plyr__controls__row--top .plyr__progress) {
  flex: 1;
}

:deep(.plyr__controls__row--bottom) {
  display: flex;
  align-items: center;
  gap: 4px;
  width: 100%;
}

:deep(.plyr__controls__row--top .plyr__time--current) {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.85);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
  flex-shrink: 0;
}

:deep(.plyr__caption) {
  font-size: var(--subtitle-size) !important;
  color: var(--subtitle-color, #ffffff) !important;
  font-family: 'DM Sans', system-ui, sans-serif;
  font-weight: 500;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.9);
  background: rgba(0, 0, 0, 0.5);
  border-radius: 4px;
  padding: 2px 8px;
  transition: font-size 0.15s ease, color 0.15s ease;
}

:deep(.plyr__time-display) {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
  white-space: nowrap;
}

:deep(.plyr__time-sep) {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

:deep(.plyr__time--duration) {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  font-variant-numeric: tabular-nums;
}

/* ── Custom Settings Menu ─────────────────────────────────────── */
:deep(.plyr-custom-menu) {
  position: absolute;
  bottom: 60px;
  right: 14px;
  z-index: 20;
  background: rgba(12, 12, 22, 0.97);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 14px;
  width: 280px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(16px);
  overflow: hidden;
}

:deep(.plyr-custom-menu__inner) {
  display: flex;
  flex-direction: column;
}

/* ── Tabs ── */
:deep(.plyr-tabs) {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 10px 10px 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.07);
}

:deep(.plyr-tab-btn) {
  flex: 1;
  padding: 5px 8px;
  border-radius: 20px;
  border: 1px solid transparent;
  background: transparent;
  color: #64748b;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s;
  text-align: center;
}

:deep(.plyr-tab-btn:hover) {
  color: #cbd5e1;
}

:deep(.plyr-tab-btn--active) {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.12);
  color: #fff;
}

:deep(.plyr-tab-close) {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  border: none;
  background: rgba(255, 255, 255, 0.07);
  color: #64748b;
  font-size: 11px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
  flex-shrink: 0;
  font-family: inherit;
}

:deep(.plyr-tab-close:hover) {
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
}

/* ── Tab panels ── */
:deep(.plyr-tab-panel) {
  padding: 6px 0 8px;
}

/* ── Menu list (speed items) ── */
:deep(.plyr-menu-list) {
  display: flex;
  flex-direction: column;
  padding: 0 8px;
}

:deep(.plyr-menu-item) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 9px 10px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: #94a3b8;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s;
  text-align: left;
  width: 100%;
}

:deep(.plyr-menu-item:hover) {
  background: rgba(255, 255, 255, 0.06);
  color: #fff;
}

:deep(.plyr-menu-item--active) {
  color: #fff;
  background: rgba(99, 102, 241, 0.15);
}

:deep(.plyr-menu-check) {
  width: 14px;
  height: 14px;
  color: #6366f1;
  opacity: 0;
  flex-shrink: 0;
  transition: opacity 0.15s;
}

:deep(.plyr-menu-item--active .plyr-menu-check) {
  opacity: 1;
}

/* ── Slider rows ── */
:deep(.plyr-menu-slider-row) {
  padding: 8px 10px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 10px;
  margin: 4px 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

:deep(.plyr-menu-slider-header) {
  display: flex;
  align-items: center;
  gap: 7px;
}

:deep(.plyr-menu-slider-icon) {
  color: #6366f1;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

:deep(.plyr-menu-slider-label) {
  font-size: 12px;
  font-weight: 600;
  color: #e2e8f0;
  flex: 1;
}

:deep(.plyr-menu-slider-value) {
  font-size: 12px;
  font-weight: 700;
  color: #6366f1;
  min-width: 36px;
  text-align: right;
  flex-shrink: 0;
}

:deep(.plyr-menu-slider) {
  width: 100%;
  accent-color: #6366f1;
  cursor: pointer;
  height: 4px;
  border-radius: 2px;
}

/* ── Color section ── */
:deep(.plyr-menu-color-section) {
  padding: 8px 10px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 10px;
  margin: 4px 8px 8px;
}

:deep(.plyr-sub-color-grid) {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 7px;
}

:deep(.plyr-sub-color-btn) {
  width: 100%;
  aspect-ratio: 1;
  border-radius: 8px;
  border: 2px solid transparent;
  cursor: pointer;
  padding: 0;
  transition: all 0.15s;
}

:deep(.plyr-sub-color-btn:hover) {
  transform: scale(1.1);
  border-color: rgba(255, 255, 255, 0.3);
}

:deep(.plyr-sub-color-btn--active) {
  border-color: #fff;
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.7);
  transform: scale(1.05);
}

/* ── Play Button (poster overlay) ───────────────────────────── */
.vp-poster__play {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 68px;
  height: 68px;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.9);
  background: rgba(0, 0, 0, 0.45);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  backdrop-filter: blur(6px);
  transition: all 0.2s ease;
  pointer-events: all;   /* override pointer-events: none dari .vp-poster */
  padding-left: 4px;     /* optik agar icon play tampak centered */
}

.vp-poster__play:hover {
  background: rgba(99, 102, 241, 0.75);
  border-color: #fff;
  transform: translate(-50%, -50%) scale(1.08);
}

.vp-poster__play:active {
  transform: translate(-50%, -50%) scale(0.96);
}

/* Mobile — sedikit lebih kecil */
@media (max-width: 640px) {
  .vp-poster__play {
    width: 54px;
    height: 54px;
  }
  .vp-poster__play svg {
    width: 22px;
    height: 22px;
  }
}

/* ── Transitions ─────────────────────────────────────────────── */
.vp-fade-slide-enter-active,
.vp-fade-slide-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.vp-fade-slide-enter-from,
.vp-fade-slide-leave-to {
  opacity: 0; transform: translateY(8px);
}

/* ── Poster Overlay ──────────────────────────────────────────── */
.vp-poster {
  position: absolute;
  inset: 0;
  z-index: 5;
  background-size: cover;
  background-position: center top;
  pointer-events: none;   /* biar controls tetap bisa diklik */
}

.vp-poster__gradient {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to bottom,
    rgba(0,0,0,0.15) 0%,
    rgba(0,0,0,0.05) 30%,
    rgba(0,0,0,0.55) 70%,
    rgba(0,0,0,0.88) 100%
  );
}

.vp-poster__img {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center top;
}

/* ── Fade transition ─────────────────────────────────────────── */
.vp-fade-enter-active,
.vp-fade-leave-active {
  transition: opacity 0.4s ease;
}
.vp-fade-enter-from,
.vp-fade-leave-to {
  opacity: 0;
}

/* ── Mobile ──────────────────────────────────────────────────── */
@media (max-width: 640px) {

}

/* ── Responsive ──────────────────────────────────────────────── */
@media (max-width: 640px) {
  /* Settings panel — compact, tidak menutupi player */
  :deep(.plyr-custom-menu) {
    position: absolute;
    width: 220px;          /* lebar fixed, tidak full width */
    right: 8px;
    left: auto;            /* reset left */
    bottom: 52px;          /* tepat di atas control bar */
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.08);
    background: rgba(10, 10, 20, 0.99);
    backdrop-filter: none;
    max-height: 200px;     /* batasi tinggi agar muat di player */
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
  }

  /* Tabs lebih compact */
  :deep(.plyr-tabs) {
    padding: 8px 8px 6px;
    gap: 3px;
  }

  :deep(.plyr-tab-btn) {
    padding: 4px 6px;
    font-size: 11px;
  }

  :deep(.plyr-tab-close) {
    width: 22px;
    height: 22px;
    font-size: 11px;
  }

  /* Slider rows lebih compact */
  :deep(.plyr-menu-slider-row) {
    margin: 3px 6px;
    padding: 6px 8px;
    gap: 6px;
  }

  :deep(.plyr-menu-slider-label) {
    font-size: 11px;
  }

  :deep(.plyr-menu-slider-value) {
    font-size: 11px;
    min-width: 28px;
  }

  :deep(.plyr-menu-slider-icon svg) {
    width: 11px;
    height: 11px;
  }

  /* Color section compact */
  :deep(.plyr-menu-color-section) {
    margin: 3px 6px 6px;
    padding: 6px 8px;
  }

  :deep(.plyr-sub-color-grid) {
    grid-template-columns: repeat(6, 1fr);  /* tetap 6 kolom tapi lebih kecil */
    gap: 5px;
  }

  /* Speed items */
  :deep(.plyr-menu-item) {
    padding: 8px 8px;
    font-size: 12px;
  }

  :deep(.plyr-menu-list) {
    padding: 0 4px;
  }
}

@media (max-width: 768px) {
  .vp-skip-intro {
    backdrop-filter: none;
    background: rgba(0, 0, 0, 0.85);
  }
}
</style>