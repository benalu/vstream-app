<script setup>
const TYPE_LABEL = { movie: 'Movie', series: 'Series', anime: 'Anime' }

defineProps({
  info:             { type: Object,  required: true },
  hasEpisodes:      { type: Boolean, default: false },
  activeEpisode:    { type: Object,  default: null },
  activeSeason:     { type: Number,  default: 1 },
  availableServers: { type: Array,   default: () => [] },
  activeServer:     { type: Number,  default: 1 },
  hasError:         { type: Boolean, default: false },
})

const emit = defineEmits(['update:activeServer', 'reportError'])
</script>

<template>
  <div>
    <div v-if="hasError" class="wp-report-bar">
      <span class="wp-report-label">Video tidak bisa diputar?</span>
      <button class="wp-report-btn" @click="emit('reportError', 'manual')">
        Laporkan Server {{ activeServer }}
      </button>
    </div>

    <div class="wp-meta-bar">
      <div class="wp-meta-left">
        <span class="wp-type-badge" :class="`wp-type-badge--${info.type}`">{{ TYPE_LABEL[info.type] }}</span>
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
      <div v-if="availableServers.length > 0" class="wp-servers">
        <span class="wp-servers-label">Server:</span>
        <button v-for="s in availableServers" :key="s.n" class="wp-server-btn"
          :class="{ 'wp-server-btn--active': activeServer === s.n }"
          @click="emit('update:activeServer', s.n)">
          {{ s.n }}
        </button>
      </div>
    </div>
  </div>
</template>