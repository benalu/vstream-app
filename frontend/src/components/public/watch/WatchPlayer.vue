<script setup>
import { computed } from 'vue'

const props = defineProps({
  videoRef:        { type: Object },
  currentUrl:      { type: String, required: true },
  isDirectVideo:   { type: Boolean, required: true },
  subtitles:       { type: Array,   default: () => [] },
  subtitleStyle:   { type: Object,  default: () => ({}) },
  showSkipIntro:   { type: Boolean, default: false },
  hasEpisodes:     { type: Boolean, default: false },
  showResumePrompt:{ type: Boolean, default: false },
  resumeFormatted: { type: String,  default: '' },
})

const emit = defineEmits(['skipIntro', 'resume', 'startOver'])
</script>

<template>
  <div class="wp-player-wrap" :style="subtitleStyle">

    <template v-if="currentUrl && isDirectVideo">
      <video ref="videoRef" :key="currentUrl" :src="currentUrl" class="wp-video" preload="metadata" playsinline>
        <track v-for="sub in subtitles" :key="sub.lang" :src="sub.url" :srclang="sub.lang"
          :label="sub.label" kind="subtitles" :default="sub.lang === 'id'" />
      </video>
    </template>

    <iframe v-else-if="currentUrl && !isDirectVideo" :key="'iframe-' + currentUrl"
      :src="currentUrl" class="wp-iframe" allowfullscreen allow="autoplay; fullscreen"
      frameborder="0" sandbox="allow-scripts allow-same-origin allow-presentation allow-forms" />

    <div v-else class="wp-no-url">
      <svg width="52" height="52" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2">
        <polygon points="5 3 19 12 5 21 5 3" opacity=".3"/>
      </svg>
      <p>URL stream belum tersedia</p>
    </div>

    <Transition name="wp-fade-slide">
      <button v-if="showSkipIntro && hasEpisodes" class="wp-skip-intro" @click="emit('skipIntro')">
        Skip Intro →
      </button>
    </Transition>

    <Transition name="wp-fade-slide">
      <div v-if="showResumePrompt" class="wp-resume-prompt">
        <span class="wp-resume-text">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="currentColor"><polygon points="5 3 19 12 5 21 5 3"/></svg>
          Lanjut dari {{ resumeFormatted }}?
        </span>
        <button class="wp-resume-btn wp-resume-btn--primary" @click="emit('resume')">Lanjutkan</button>
        <button class="wp-resume-btn" @click="emit('startOver')">Dari Awal</button>
      </div>
    </Transition>

  </div>
</template>