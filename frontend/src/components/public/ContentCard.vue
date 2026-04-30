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
       <!-- Rating badge -->
      <span v-if="item.rating" class="cc__rating">
        <svg width="10" height="10" viewBox="0 0 24 24" fill="#f59e0b">
          <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
        </svg>
        {{ item.rating }}
      </span>
    </div>

    <!-- Info -->
    <div class="cc__info">
      <p class="cc__title">{{ item.title }}</p>
      <span v-if="item.year" class="cc__year">{{ item.year }}</span>
    </div>
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

/* Rating */
.cc__rating {
  position: absolute;
  top: 8px; right: 8px;
  display: flex;
  align-items: center;
  gap: 3px;
  font-size: 10px; font-weight: 700;
  color: #f59e0b;
  background: rgba(0,0,0,0.65);
  backdrop-filter: blur(4px);
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid rgba(245,158,11,0.25);
}

/* Title */
.cc__info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.cc__title {
  font-size: 12.5px; font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}
.cc__year {
  font-size: 11px; font-weight: 500;
  color: #64748b;
  flex-shrink: 0;
}
</style>