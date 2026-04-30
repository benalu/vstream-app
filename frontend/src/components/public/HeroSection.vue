<script setup>
import { ref } from 'vue'

const props = defineProps({
  /** Featured content object */
  content: {
    type: Object,
    required: true,
    // Expected shape:
    // { title, rating, year, genres: [], overview, backdrop, tmdb_id, type }
  },
})

const emit = defineEmits(['play', 'more'])

const backdropLoaded = ref(false)
</script>

<template>
  <section class="hero">
    <!-- Backdrop -->
    <div class="hero__backdrop-wrap">
      <img
        :src="content.backdrop"
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
    <div class="hero__content">
      <!-- Meta row -->
      <div class="hero__meta">
        <span class="hero__year">{{ content.year }}</span>
        <span class="hero__dot">·</span>
        <span
          v-for="g in content.genres"
          :key="g"
          class="hero__genre-pill"
        >{{ g }}</span>
      </div>

      <!-- Title -->
      <h1 class="hero__title">{{ content.title }}</h1>

      <!-- Rating -->
      <div class="hero__rating">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="#f59e0b">
          <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
        </svg>
        <span class="hero__rating-val">{{ content.rating }}</span>
        <span class="hero__rating-sep">·</span>
        <span class="hero__rating-label">IMDb</span>
      </div>

      <!-- Overview -->
      <p class="hero__overview">{{ content.overview }}</p>

      <!-- CTA -->
      <div class="hero__ctas">
        <button class="hero__btn hero__btn--play" @click="emit('play', content)">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 3 19 12 5 21 5 3"/>
          </svg>
          Tonton Sekarang
        </button>
        <button class="hero__btn hero__btn--info" @click="emit('more', content)">
          <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          Selengkapnya
        </button>
      </div>
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

/* ── Backdrop ───────────────────────────────────────────────── */
.hero__backdrop-wrap {
  position: absolute;
  inset: 0;
}
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
  to   { transform: scale(1.0);  }
}

/* Gradient layers */
.hero__grad-bottom {
  position: absolute;
  bottom: 0; left: 0; right: 0;
  height: 75%;
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
  opacity: 0.5;
  pointer-events: none;
}

/* ── Content ────────────────────────────────────────────────── */
.hero__content {
  position: relative;
  z-index: 2;
  padding: 0 48px 72px;
  max-width: 620px;
  animation: hero-in 0.9s cubic-bezier(0.22, 1, 0.36, 1) both;
}
@keyframes hero-in {
  from { opacity: 0; transform: translateY(30px); }
}

.hero__meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
  flex-wrap: wrap;
}
.hero__year {
  font-size: 13px; font-weight: 600;
  color: #64748b;
}
.hero__dot { color: #64748b; font-size: 12px; }
.hero__genre-pill {
  font-size: 12px; font-weight: 500;
  color: #a0aec0;
  padding: 3px 10px;
  border-radius: 20px;
  background: rgba(255,255,255,0.08);
  border: 1px solid rgba(255,255,255,0.1);
}

.hero__title {
  font-family: 'Bebas Neue', sans-serif;
  font-size: clamp(52px, 8vw, 96px);
  letter-spacing: 0.04em;
  color: #fff;
  line-height: 0.9;
  margin-bottom: 14px;
  text-shadow: 0 4px 40px rgba(0,0,0,0.5);
}

.hero__rating {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 16px;
}
.hero__rating-val { font-size: 14px; font-weight: 700; color: #f59e0b; }
.hero__rating-sep { color: #64748b; font-size: 12px; }
.hero__rating-label { font-size: 12px; color: #64748b; }

.hero__overview {
  font-size: 14.5px;
  line-height: 1.65;
  color: #94a3b8;
  margin-bottom: 28px;
  max-width: 500px;
}

/* CTA buttons */
.hero__ctas { display: flex; gap: 12px; flex-wrap: wrap; }
.hero__btn {
  display: inline-flex;
  align-items: center;
  gap: 9px;
  padding: 12px 26px;
  border-radius: 40px;
  font-size: 14px; font-weight: 600;
  font-family: 'DM Sans', sans-serif;
  cursor: pointer; border: none;
  transition: all 0.2s;
  white-space: nowrap;
}
.hero__btn--play {
  background: #fff; color: #0a0a14;
  box-shadow: 0 8px 30px rgba(255,255,255,0.15);
}
.hero__btn--play:hover {
  background: #e2e8f0;
  transform: translateY(-2px);
  box-shadow: 0 12px 36px rgba(255,255,255,0.2);
}
.hero__btn--info {
  background: rgba(255,255,255,0.12); color: #fff;
  border: 1px solid rgba(255,255,255,0.18);
  backdrop-filter: blur(10px);
}
.hero__btn--info:hover {
  background: rgba(255,255,255,0.2);
  transform: translateY(-2px);
}

@media (max-width: 768px) {
  .hero__content { padding: 0 20px 60px; }
  .hero__title   { font-size: 52px; }
}
</style>