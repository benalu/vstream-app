<script setup>
/**
 * ContentCard — reusable poster card.
 *
 * Props:
 *   item  Object — { id, title, poster, type, tmdb_id }
 *
 * Emits:
 *   click — when card / play button is clicked
 */
defineProps({
  item: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['click'])

const TYPE_LABEL = { movie: 'Movie', series: 'Series', anime: 'Anime' }
</script>

<template>
  <article
    class="cc"
    @click="emit('click', item)"
    :aria-label="item.title"
    role="button"
    tabindex="0"
    @keydown.enter="emit('click', item)"
  >
    <!-- Poster -->
    <div class="cc__poster-wrap">
      <img :src="item.poster" class="cc__poster" :alt="item.title" loading="lazy" />

      <!-- Play overlay -->
      <div class="cc__overlay">
        <button class="cc__play" tabindex="-1" aria-hidden="true">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 3 19 12 5 21 5 3"/>
          </svg>
        </button>
      </div>

      <!-- Type badge -->
      <span class="cc__type" :class="`cc__type--${item.type}`">
        {{ TYPE_LABEL[item.type] ?? item.type }}
      </span>
    </div>

    <!-- Info -->
    <p class="cc__title">{{ item.title }}</p>
  </article>
</template>

<style scoped>
.cc {
  cursor: pointer;
  transition: transform 0.2s;
  outline: none;
}
.cc:focus-visible .cc__poster-wrap {
  box-shadow: 0 0 0 2px #6366f1;
}
.cc:hover { transform: translateY(-4px); }

/* Poster */
.cc__poster-wrap {
  position: relative;
  width: 100%;
  aspect-ratio: 2/3;
  border-radius: 10px;
  overflow: hidden;
  background: #0d0d18;
  margin-bottom: 10px;
  box-shadow: 0 6px 20px rgba(0,0,0,0.4);
}
.cc__poster {
  width: 100%; height: 100%;
  object-fit: cover;
  transition: transform 0.4s;
}
.cc:hover .cc__poster { transform: scale(1.07); }

/* Play overlay */
.cc__overlay {
  position: absolute; inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}
.cc:hover .cc__overlay { opacity: 1; }
.cc__play {
  width: 42px; height: 42px;
  border-radius: 50%;
  border: none;
  background: rgba(255,255,255,0.92);
  color: #0a0a14;
  display: grid; place-items: center;
  cursor: pointer;
}

/* Type badge */
.cc__type {
  position: absolute;
  top: 8px; left: 8px;
  font-size: 9px; font-weight: 700;
  padding: 2px 7px; border-radius: 4px;
  text-transform: uppercase; letter-spacing: 0.06em;
  color: #fff;
}
.cc__type--movie  { background: rgba(99,102,241,0.85); }
.cc__type--series { background: rgba(6,182,212,0.85);  }
.cc__type--anime  { background: rgba(244,63,94,0.85);  }

/* Title */
.cc__title {
  font-size: 12.5px; font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>