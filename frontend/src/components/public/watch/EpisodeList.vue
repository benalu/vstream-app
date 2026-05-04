<script setup>
/**
 * EpisodeList — season tabs + episode grid untuk series & anime.
 *
 * Props:
 *   seasons        Array   — [{ season_num, episodes: [{ id, ep_num, title, url1, url2 }] }]
 *   activeEpisode  Object  — episode yang sedang aktif
 *   activeSeason   Number  — season yang sedang ditampilkan (v-model)
 *
 * Emits:
 *   update:activeSeason — saat season tab diklik
 *   select-episode      — saat episode card diklik, payload: episode object
 */
defineProps({
  seasons: {
    type: Array,
    required: true,
  },
  activeEpisode: {
    type: Object,
    default: null,
  },
  activeSeason: {
    type: Number,
    required: true,
  },
})

const emit = defineEmits(['update:activeSeason', 'select-episode'])
</script>

<template>
  <div v-if="seasons.length > 0" class="el-wrap">

    <!-- Header: judul + season tabs -->
    <div class="el-header">
      <div class="el-header-left">
        <div class="el-bar" aria-hidden="true" />
        <h2 class="el-title">Episode</h2>
      </div>

      <div v-if="seasons.length > 1" class="el-season-tabs" role="tablist">
        <button
          v-for="s in seasons"
          :key="s.season_num"
          class="el-season-tab"
          :class="{ 'el-season-tab--active': activeSeason === s.season_num }"
          role="tab"
          :aria-selected="activeSeason === s.season_num"
          @click="emit('update:activeSeason', s.season_num)"
        >
          Season {{ s.season_num }}
        </button>
      </div>
    </div>

    <!-- Episode grid -->
    <div
      class="el-grid"
      role="list"
      :aria-label="`Episode Season ${activeSeason}`"
    >
      <button
        v-for="ep in seasons.find(s => s.season_num === activeSeason)?.episodes ?? []"
        :key="ep.id"
        class="el-ep-card"
        :class="{ 'el-ep-card--active': activeEpisode?.id === ep.id }"
        role="listitem"
        :aria-current="activeEpisode?.id === ep.id ? 'true' : undefined"
        @click="emit('select-episode', ep)"
      >
        <!-- Number / play icon -->
        <div class="el-ep-num">
          <svg
            v-if="activeEpisode?.id === ep.id"
            width="14" height="14"
            viewBox="0 0 24 24"
            fill="currentColor"
            aria-hidden="true"
          >
            <polygon points="5 3 19 12 5 21 5 3"/>
          </svg>
          <span v-else>{{ ep.ep_num }}</span>
        </div>

        <!-- Info -->
        <div class="el-ep-info">
          <span class="el-ep-label">Episode {{ ep.ep_num }}</span>
          <span v-if="ep.title" class="el-ep-title">{{ ep.title }}</span>
        </div>
      </button>
    </div>

  </div>
</template>

<style scoped>
.el-wrap {
  margin-top: 8px;
}

/* ── Header ──────────────────────────────────────────────────── */
.el-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  gap: 16px;
  flex-wrap: wrap;
}

.el-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.el-bar {
  width: 4px; height: 20px;
  border-radius: 2px;
  background: #6366f1;
  box-shadow: 0 0 10px rgba(99, 102, 241, 0.6);
  flex-shrink: 0;
}

.el-title {
  font-size: 17px;
  font-weight: 700;
  color: #fff;
}

/* ── Season tabs ─────────────────────────────────────────────── */
.el-season-tabs {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.el-season-tab {
  padding: 5px 14px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: transparent;
  color: #64748b;
  font-size: 12.5px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s;
}

.el-season-tab:hover { color: #fff; }

.el-season-tab--active {
  background: rgba(99, 102, 241, 0.15);
  border-color: rgba(99, 102, 241, 0.4);
  color: #818cf8;
}

/* ── Episode grid ────────────────────────────────────────────── */
.el-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 10px;
}

.el-ep-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(255, 255, 255, 0.02);
  color: #64748b;
  cursor: pointer;
  text-align: left;
  font-family: inherit;
  transition: all 0.15s;
  width: 100%;
}

.el-ep-card:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.12);
  color: #fff;
  transform: translateY(-2px);
}

.el-ep-card--active {
  background: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.3);
  color: #fff;
}

/* Number badge */
.el-ep-num {
  width: 30px; height: 30px;
  border-radius: 7px;
  background: rgba(255, 255, 255, 0.06);
  display: grid; place-items: center;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
  color: #94a3b8;
}

.el-ep-card--active .el-ep-num {
  background: rgba(99, 102, 241, 0.25);
  color: #818cf8;
}

/* Info */
.el-ep-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.el-ep-label {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
}

.el-ep-card--active .el-ep-label { color: #818cf8; }

.el-ep-title {
  font-size: 12px;
  font-weight: 500;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* ── Responsive ──────────────────────────────────────────────── */
@media (max-width: 640px) {
  .el-grid {
    grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
  }
}
</style>