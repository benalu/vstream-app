<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  slides: {
    type: Array,
    required: true,
    // [{ id, order, movie: { title, rating, year, genres, overview, backdrop, tmdb_id, type, logos } }]
  },
})

const emit = defineEmits(['play', 'more'])

const activeIndex = ref(0)
const backdropLoaded = ref(false)
let timer = null

const current = computed(() => props.slides[activeIndex.value]?.movie ?? null)

const heroLogo = computed(() => {
  try {
    const raw = current.value?.logos
    if (!raw) return null
    const logos = typeof raw === 'string' ? JSON.parse(raw) : raw
    return logos.find(l => l.aspect_ratio >= 3 && l.aspect_ratio <= 8) ?? logos[0] ?? null
  } catch {
    return null
  }
})

const goTo = (i) => {
  backdropLoaded.value = false
  activeIndex.value = i
  resetTimer()
}

const resetTimer = () => {
  clearInterval(timer)
  if (props.slides.length > 1) {
    timer = setInterval(() => {
      backdropLoaded.value = false
      activeIndex.value = (activeIndex.value + 1) % props.slides.length
    }, 7000)
  }
}

onMounted(resetTimer)
onUnmounted(() => clearInterval(timer))
</script>

<template>
  <section class="hero" v-if="current">
    <!-- Backdrop -->
    <div class="hero__backdrop-wrap">
      <img
        :key="current.tmdb_id"
        :src="current.backdrop"
        class="hero__backdrop"
        :class="{ 'hero__backdrop--loaded': backdropLoaded }"
        @load="backdropLoaded = true"
        alt=""
        loading="eager"
      />
      <div class="hero__grad-bottom" />
      <div class="hero__grad-left" />
      <div class="hero__grad-right" />
      <div class="hero__noise" />
    </div>

    <!-- Content -->
    <div class="hero__content" :key="current.tmdb_id">
      <div class="hero__meta">
        <span class="hero__type-badge" :class="`hero__type-badge--${current.type}`">
          {{ current.type === 'movie' ? 'Movie' : current.type === 'series' ? 'Series' : 'Anime' }}
        </span>
        <span class="hero__dot">·</span>
        <div v-if="current.rating" class="hero__rating">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="#f59e0b">
            <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
          </svg>
          <span class="hero__rating-val">{{ current.rating }}</span>
        </div>
        <span class="hero__dot">·</span>
        <span class="hero__year">{{ current.year }}</span>
      </div>

      <img v-if="heroLogo" :src="heroLogo.file_path" :alt="current.title" class="hero__logo" />
      <h1 v-else class="hero__title">{{ current.title }}</h1>

      <p class="hero__overview">{{ current.overview }}</p>

      <div class="hero__ctas">
        <button class="hero__btn hero__btn--play" @click="emit('play', current)">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 3 19 12 5 21 5 3"/>
          </svg>
          Tonton Sekarang
        </button>
        <button class="hero__btn hero__btn--info" @click="emit('more', current)">
          <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          Selengkapnya
        </button>
      </div>
    </div>

    <!-- Dots -->
    <div v-if="slides.length > 1" class="hero__dots">
      <button
        v-for="(_, i) in slides"
        :key="i"
        class="hero__dot-btn"
        :class="{ 'hero__dot-btn--active': i === activeIndex }"
        @click="goTo(i)"
        :aria-label="`Slide ${i + 1}`"
      />
    </div>
  </section>
</template>

<style scoped>
.hero {
  position: relative;
  height: 100svh;
  min-height: 600px;
  max-height: 900px;
  display: flex;
  align-items: flex-end;
  overflow: hidden;
}

.hero__backdrop-wrap { position: absolute; inset: 0; }
.hero__backdrop {
  width: 100%; height: 100%;
  object-fit: cover;
  object-position: center 20%;
  opacity: 0;
  transition: opacity 0.8s ease;
  animation: hero-zoom 18s ease-in-out infinite alternate;
}
.hero__backdrop--loaded { opacity: 1; }
@keyframes hero-zoom {
  from { transform: scale(1.04); }
  to   { transform: scale(1.0); }
}

.hero__grad-bottom {
  position: absolute; bottom: 0; left: 0; right: 0; height: 75%;
  background: linear-gradient(to top, #07070e 0%, rgba(7,7,14,0.7) 50%, transparent 100%);
}
.hero__grad-left {
  position: absolute; inset: 0;
  background: linear-gradient(to right, rgba(7,7,14,0.85) 0%, rgba(7,7,14,0.3) 45%, transparent 70%);
}
.hero__grad-right {
  position: absolute; inset: 0;
  background: linear-gradient(to left, rgba(7,7,14,0.4) 0%, transparent 50%);
}
.hero__noise {
  position: absolute; inset: 0;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='0.03'/%3E%3C/svg%3E");
  opacity: 0.5; pointer-events: none;
}

.hero__content {
  position: relative; z-index: 2;
  padding: 0 48px 72px;
  max-width: 620px;
  animation: hero-in 0.7s cubic-bezier(0.22, 1, 0.36, 1) both;
}
@keyframes hero-in {
  from { opacity: 0; transform: translateY(24px); }
}

.hero__meta {
  display: flex; align-items: center; gap: 8px;
  margin-bottom: 14px; flex-wrap: wrap;
}
.hero__year { font-size: 13px; font-weight: 600; color: #64748b; }
.hero__dot  { color: #64748b; font-size: 12px; }

.hero__title {
  font-family: 'Bebas Neue', sans-serif;
  font-size: clamp(52px, 8vw, 96px);
  letter-spacing: 0.04em; color: #fff;
  line-height: 0.9; margin-bottom: 14px;
  text-shadow: 0 4px 40px rgba(0,0,0,0.5);
}
.hero__logo {
  max-width: min(380px, 60vw); max-height: 120px;
  width: auto; height: auto;
  object-fit: contain; object-position: left center;
  margin-bottom: 14px;
  filter: drop-shadow(0 4px 20px rgba(0,0,0,0.6));
  display: block;
}

.hero__type-badge {
  font-size: 10px; font-weight: 700;
  padding: 3px 9px; border-radius: 4px;
  text-transform: uppercase; letter-spacing: 0.07em; color: #fff;
}
.hero__type-badge--movie  { background: rgba(99,102,241,0.85); }
.hero__type-badge--series { background: rgba(6,182,212,0.85); }
.hero__type-badge--anime  { background: rgba(244,63,94,0.85); }

.hero__rating { display: flex; align-items: center; gap: 4px; }
.hero__rating-val { font-size: 13px; font-weight: 700; color: #f59e0b; }

.hero__overview {
  font-size: 14.5px; line-height: 1.65;
  color: #94a3b8; margin-bottom: 28px; max-width: 500px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.hero__ctas { display: flex; gap: 12px; flex-wrap: wrap; }
.hero__btn {
  display: inline-flex; align-items: center; gap: 9px;
  padding: 12px 26px; border-radius: 40px;
  font-size: 14px; font-weight: 600;
  font-family: 'DM Sans', sans-serif;
  cursor: pointer; border: none; transition: all 0.2s; white-space: nowrap;
}
.hero__btn--play {
  background: #fff; color: #0a0a14;
  box-shadow: 0 8px 30px rgba(255,255,255,0.15);
}
.hero__btn--play:hover {
  background: #e2e8f0; transform: translateY(-2px);
  box-shadow: 0 12px 36px rgba(255,255,255,0.2);
}
.hero__btn--info {
  background: rgba(255,255,255,0.12); color: #fff;
  border: 1px solid rgba(255,255,255,0.18); backdrop-filter: blur(10px);
}
.hero__btn--info:hover { background: rgba(255,255,255,0.2); transform: translateY(-2px); }

/* Dots */
.hero__dots {
  position: absolute;
  bottom: 28px; right: 48px;
  z-index: 10;
  display: flex; gap: 6px; align-items: center;
}
.hero__dot-btn {
  width: 6px; height: 6px;
  border-radius: 50%; border: none;
  background: rgba(255,255,255,0.3);
  cursor: pointer; padding: 0;
  transition: all 0.3s;
}
.hero__dot-btn--active {
  background: #fff;
  width: 20px;
  border-radius: 3px;
}

@media (max-width: 768px) {
  .hero__content { padding: 0 20px 60px; }
  .hero__title   { font-size: 52px; }
  .hero__dots    { right: 20px; bottom: 20px; }
}
</style>