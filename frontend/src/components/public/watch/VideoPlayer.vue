<script setup>
/**
 * VideoPlayer — player wrapper yang menggabungkan Plyr (direct video)
 * atau iframe (embed), skip intro, resume prompt, dan subtitle tracks.
 *
 * Props:
 *   src          String  — URL stream saat ini
 *   subtitles    Array   — [{ lang, label, url }]
 *   hasEpisodes  Boolean — true jika series/anime (aktifkan skip intro)
 *   storageKey   String  — localStorage key untuk progress (null = disable)
 *
 * Emits:
 *   error        — saat video gagal load
 *   timeupdate   — setiap timeupdate dari Plyr (untuk logic di parent)
 *   ended        — saat video selesai
 */
import { ref, computed, watch, nextTick, watchEffect, onBeforeUnmount } from 'vue'
import { usePlayer }   from '@/composables/usePlayer'
import { useSubtitle } from '@/composables/useSubtitle'
import { useProgress } from '@/composables/useProgress'

const props = defineProps({
  src:         { type: String,  default: '' },
  subtitles:   { type: Array,   default: () => [] },
  hasEpisodes: { type: Boolean, default: false },
  storageKey:  { type: String,  default: null },
})

const emit = defineEmits(['error', 'timeupdate', 'ended'])

// ── Refs ───────────────────────────────────────────────────────
const videoRef  = ref(null)
const hasError  = ref(false)

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

// ── Player ─────────────────────────────────────────────────────
const storageKeyRef = computed(() => props.storageKey)

const {
  showResumePrompt,
  resumeFormatted,
  saveProgress,
  restoreProgress,
  resume,
  startOver,
  clearProgress,
} = useProgress(storageKeyRef, ref(null)) // plyrInstance injected below

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

// Patch useProgress agar pakai plyrInstance yang baru dibuat
// (workaround karena plyrInstance dibuat setelah useProgress dipanggil)
const _origSave    = saveProgress
const _origRestore = restoreProgress

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

// Reset state on server / episode change (when src changes)
watch(() => props.src, () => {
  hasError.value         = false
  showSkipIntro.value    = false
})

onBeforeUnmount(() => {
  saveProgress()
  destroyPlyr()
})

// Expose untuk parent (misal scroll to player)
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
        <span class="vp-resume-text">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 3 19 12 5 21 5 3"/>
          </svg>
          Lanjut dari {{ resumeFormatted }}?
        </span>
        <button class="vp-resume-btn vp-resume-btn--primary" @click="resume">
          Lanjutkan
        </button>
        <button class="vp-resume-btn" @click="startOver">
          Dari Awal
        </button>
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
  bottom: 80px; left: 20px;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.75);
  border: 1px solid rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(8px);
  flex-wrap: wrap;
  max-width: min(400px, calc(100% - 40px));
}
.vp-resume-text {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #fff;
  font-size: 13px;
  font-weight: 500;
}
.vp-resume-btn {
  padding: 4px 12px;
  border-radius: 5px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: transparent;
  color: #fff;
  font-size: 12px;
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
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
  padding: 20px 14px 12px;
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


/* ── Plyr settings menu overrides ────────────────────────────── */
:deep(.plyr__menu__container) {
  min-width: 200px;
  background: #1a1a2e !important;
  border: 1px solid rgba(99, 102, 241, 0.2);
}
:deep(.plyr__menu__container *) { color: #c7c8ff; }

:deep(.plyr-subtitle-size-panel) {
  display: flex; flex-direction: column; gap: 12px;
  padding: 12px 14px 16px;
}
:deep(.plyr-sub-size-label) {
  font-size: 11px; font-weight: 600;
  color: #64748b;
  text-transform: uppercase; letter-spacing: 0.06em;
  padding: 8px 14px 4px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  margin-top: 4px;
}
:deep(.plyr-sub-size-row) {
  display: flex; align-items: center;
  justify-content: space-between; gap: 8px;
  background: rgba(0, 0, 0, 0.25);
  border-radius: 8px; padding: 6px 8px;
}
:deep(.plyr-sub-size-btn) {
  background: rgba(99, 102, 241, 0.2);
  border: 1px solid rgba(99, 102, 241, 0.35);
  border-radius: 6px; color: #a5b4fc;
  font-size: 12px; font-weight: 700;
  padding: 5px 14px; cursor: pointer;
  font-family: inherit; transition: all 0.15s; flex-shrink: 0;
}
:deep(.plyr-sub-size-btn:hover) {
  background: rgba(99, 102, 241, 0.4);
  border-color: rgba(99, 102, 241, 0.7); color: #fff;
}
:deep(.plyr-sub-size-display) {
  font-size: 15px; font-weight: 700; color: #fff;
  min-width: 42px; text-align: center; letter-spacing: -0.01em;
}
:deep(.plyr-sub-size-slider-wrap) {
  display: flex; align-items: center;
  gap: 8px; width: 100%; padding: 0 2px;
}
:deep(.plyr-sub-size-min),
:deep(.plyr-sub-size-max) {
  font-size: 10px; color: #64748b; white-space: nowrap; flex-shrink: 0;
}
:deep(.plyr-sub-size-slider) {
  flex: 1; accent-color: #6366f1; cursor: pointer; height: 4px; min-width: 0;
}
:deep(.plyr-sub-color-row) {
  display: flex; align-items: center; gap: 8px;
  padding: 4px 14px 12px;
}
:deep(.plyr-sub-color-btn) {
  width: 22px; height: 22px; border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer; padding: 0; transition: all 0.15s; flex-shrink: 0;
}
:deep(.plyr-sub-color-btn:hover) { transform: scale(1.2); }
:deep(.plyr-sub-color-btn--active) {
  border-color: #fff;
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.6);
  transform: scale(1.15);
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

/* ── Responsive ──────────────────────────────────────────────── */
@media (max-width: 640px) {
  .vp-resume-prompt {
    bottom: 56px;
    left: 10px;
    right: 10px;
    max-width: unset;
    padding: 8px 12px;
    gap: 8px;
    border-radius: 6px;
    backdrop-filter: none;
    background: rgba(0, 0, 0, 0.88);
  }
  .vp-resume-text {
    font-size: 12px;
    flex: 1;
    min-width: 0;
  }
  .vp-resume-btn {
    padding: 4px 10px;
    font-size: 11px;
  }
}

@media (max-width: 768px) {
  .vp-skip-intro {
    backdrop-filter: none;
    background: rgba(0, 0, 0, 0.85);
  }
}
</style>