<script setup>
import ContentCard from './ContentCard.vue'

/**
 * ContentGrid — a titled section with a responsive grid of ContentCards.
 *
 * Props:
 *   title    String   — section title (e.g. "Baru Ditambahkan")
 *   items    Array    — list of content objects
 *   showMore Boolean  — whether to show "Lihat Semua" link (default true)
 *   moreHref String   — href for the "Lihat Semua" link
 *
 * Emits:
 *   cardClick  — item when a ContentCard is clicked
 *   moreClick  — when "Lihat Semua" is clicked
 */
defineProps({
  title: {
    type: String,
    required: true,
  },
  items: {
    type: Array,
    required: true,
  },
  showMore: {
    type: Boolean,
    default: true,
  },
  moreHref: {
    type: String,
    default: '#',
  },
})

const emit = defineEmits(['cardClick', 'moreClick'])
</script>

<template>
  <section class="cg">
    <!-- Section header -->
    <div class="cg__header">
      <div class="cg__bar" aria-hidden="true" />
      <h2 class="cg__title">{{ title }}</h2>
      <a
        v-if="showMore"
        :href="moreHref"
        class="cg__more"
        @click.prevent="emit('moreClick')"
      >
        Lihat Semua →
      </a>
    </div>

    <!-- Grid -->
    <div class="cg__grid">
      <ContentCard
        v-for="item in items"
        :key="item.id"
        :item="item"
        @click="emit('cardClick', $event)"
      />
    </div>

    <!-- Empty state -->
    <div v-if="items.length === 0" class="cg__empty">
      <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <rect x="2" y="2" width="20" height="20" rx="2.18"/>
        <line x1="7" y1="2" x2="7" y2="22"/><line x1="17" y1="2" x2="17" y2="22"/>
        <line x1="2" y1="12" x2="22" y2="12"/>
      </svg>
      <p>Belum ada konten</p>
    </div>
  </section>
</template>

<style scoped>
.cg { /* no extra margin — parent controls spacing via gap */ }

/* Header */
.cg__header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 22px;
}
.cg__bar {
  width: 4px; height: 22px;
  border-radius: 2px;
  background: #6366f1;
  box-shadow: 0 0 10px rgba(99,102,241,0.6);
  flex-shrink: 0;
}
.cg__title {
  font-size: 18px; font-weight: 700;
  color: #fff;
  letter-spacing: -0.01em;
}
.cg__more {
  font-size: 12.5px;
  color: #818cf8;
  text-decoration: none;
  margin-left: auto;
  opacity: 0.8;
  transition: opacity 0.15s;
}
.cg__more:hover { opacity: 1; }

/* Grid */
.cg__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(148px, 1fr));
  gap: 16px;
}

/* Empty */
.cg__empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 48px 20px;
  color: #475569;
  font-size: 13px;
}
</style>