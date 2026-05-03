<script setup>
defineProps({
  title:    { type: String, required: true },
  overview: { type: String, default: '' },
  genre:    { type: String, default: '' },
})
</script>

<template>
  <div class="wp-info">
    <h1 class="wp-title">{{ title }}</h1>
    <p v-if="overview" class="wp-overview">{{ overview }}</p>
    <div v-if="genre" class="wp-genres">
      <span v-for="g in genre.split(',')" :key="g" class="wp-genre-pill">{{ g.trim() }}</span>
    </div>
  </div>
</template>