<script setup>
/**
 * Top10Row — horizontal scrollable TOP 10 list.
 *
 * Props:
 *   items  Array — list of content objects { id, title, poster, type, tmdb_id }
 *
 * Emits:
 *   click  — { item }  when a card is clicked
 */
defineProps({
  items: {
    type: Array,
    required: true,
  },
})

const emit = defineEmits(['click'])
</script>

<template>
  <div class="t10" role="list">
    <article
      v-for="(item, i) in items"
      :key="item.id"
      class="t10__card"
      role="listitem"
      :aria-label="`#${i + 1} ${item.title}`"
      @click="emit('click', item)"
    >
      <!-- Big rank number -->
      <span class="t10__rank" aria-hidden="true">
        {{ String(i + 1).padStart(2, '0') }}
      </span>

      <!-- Poster -->
      <div class="t10__poster-wrap">
        <img :src="item.poster" class="t10__poster" :alt="item.title" loading="lazy" />

        <!-- Hover overlay -->
        <div class="t10__overlay">
          <button class="t10__play-btn" tabindex="-1" aria-hidden="true">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor">
              <polygon points="5 3 19 12 5 21 5 3"/>
            </svg>
          </button>
        </div>

        <!-- Type badge -->
        <span class="t10__type" :class="`t10__type--${item.type}`">
          {{ item.type === 'movie' ? 'Movie' : item.type === 'series' ? 'Series' : 'Anime' }}
        </span>
      </div>
    </article>
  </div>
</template>

<style scoped>
.t10 {
  display: flex;
  gap: 0;
  overflow-x: auto;
  padding-bottom: 12px;
  scrollbar-width: thin;
  scrollbar-color: rgba(255,255,255,0.08) transparent;
}
.t10::-webkit-scrollbar { height: 4px; }
.t10::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.08); border-radius: 4px; }

/* Card */
.t10__card {
  position: relative;
  display: flex;
  align-items: flex-end;
  flex-shrink: 0;
  cursor: pointer;
  transition: transform 0.2s;
}
.t10__card:hover { transform: scale(1.04) translateY(-4px); z-index: 2; }

/* Rank */
.t10__rank {
  font-family: 'Bebas Neue', sans-serif;
  font-size: clamp(80px, 9vw, 130px);
  color: transparent;
  -webkit-text-stroke: 2px rgba(255,255,255,0.12);
  line-height: 1;
  margin-right: -28px;
  z-index: 1;
  flex-shrink: 0;
  align-self: flex-end;
  padding-bottom: 6px;
  user-select: none;
}

/* Poster */
.t10__poster-wrap {
  position: relative;
  width: 160px; height: 235px;
  border-radius: 10px;
  overflow: hidden;
  flex-shrink: 0;
  box-shadow: 0 8px 24px rgba(0,0,0,0.5);
}
.t10__poster {
  width: 100%; height: 100%;
  object-fit: cover;
  transition: transform 0.4s;
}
.t10__card:hover .t10__poster { transform: scale(1.08); }

/* Play overlay */
.t10__overlay {
  position: absolute; inset: 0;
  background: rgba(0,0,0,0.55);
  display: flex; align-items: center; justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}
.t10__card:hover .t10__overlay { opacity: 1; }
.t10__play-btn {
  width: 46px; height: 46px;
  border-radius: 50%;
  border: none;
  background: rgba(255,255,255,0.9);
  color: #0a0a14;
  display: grid; place-items: center;
  cursor: pointer;
  transform: scale(0.8);
  transition: transform 0.2s;
}
.t10__card:hover .t10__play-btn { transform: scale(1); }

/* Type badge */
.t10__type {
  position: absolute;
  top: 8px; left: 8px;
  font-size: 9px; font-weight: 700;
  padding: 2px 7px; border-radius: 4px;
  text-transform: uppercase; letter-spacing: 0.06em;
  color: #fff;
}
.t10__type--movie  { background: rgba(99,102,241,0.85); }
.t10__type--series { background: rgba(6,182,212,0.85);  }
.t10__type--anime  { background: rgba(244,63,94,0.85);  }
</style>