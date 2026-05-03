<script setup>
defineProps({
  seasons:       { type: Array,  required: true },
  currentSeason: { type: Object, default: null },
  activeSeason:  { type: Number, required: true },
  activeEpisode: { type: Object, default: null },
})

const emit = defineEmits(['update:activeSeason', 'selectEpisode'])
</script>

<template>
  <div v-if="seasons.length > 0" class="wp-episodes">
    <div class="wp-ep-header">
      <div class="wp-ep-header-left">
        <div class="wp-ep-bar" />
        <h2 class="wp-ep-title">Episode</h2>
      </div>
      <div v-if="seasons.length > 1" class="wp-season-tabs">
        <button v-for="s in seasons" :key="s.season_num" class="wp-season-tab"
          :class="{ 'wp-season-tab--active': activeSeason === s.season_num }"
          @click="emit('update:activeSeason', s.season_num)">
          Season {{ s.season_num }}
        </button>
      </div>
    </div>

    <div class="wp-ep-grid">
      <button v-for="ep in currentSeason?.episodes ?? []" :key="ep.id"
        class="wp-ep-card" :class="{ 'wp-ep-card--active': activeEpisode?.id === ep.id }"
        @click="emit('selectEpisode', ep)">
        <div class="wp-ep-card-num">
          <svg v-if="activeEpisode?.id === ep.id" width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
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
</template>