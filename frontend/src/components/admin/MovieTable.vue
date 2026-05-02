<script setup>
import {
  Film, Clock, Captions, Star, Calendar, Trash2, Edit3,
  ChevronDown, ChevronRight, Plus, Tv, Sparkles,
  Loader2, PlayCircle
} from 'lucide-vue-next'

defineProps({
  items:         { type: Array,   required: true },
  activeTab:     { type: String,  required: true },
  seasons:       { type: Object,  default: () => ({}) },
  seasonsLoading:{ type: Object,  default: () => ({}) },
  expandedRows:  { type: Object,  required: true }, // Set
  loading:       { type: Boolean, default: false },
})

const emit = defineEmits([
  'edit',
  'delete',
  'toggle-expand',
  'open-add-season',
  'open-subtitle',
])

const hasExpandable = (activeTab, items) =>
  activeTab === 'series' ||
  (activeTab === 'anime' && items.some(i => i.has_episodes))
</script>

<template>
  <div class="mv-table-wrap">
    <!-- Skeleton -->
    <div v-if="loading && items.length === 0" class="mv-loading">
      <div v-for="i in 4" :key="i" class="mv-skeleton-row">
        <div class="mv-skeleton-poster" />
        <div class="mv-skeleton-body">
          <div class="mv-skeleton-line mv-skeleton-line--title" />
          <div class="mv-skeleton-line mv-skeleton-line--sub" />
        </div>
        <div class="mv-skeleton-line mv-skeleton-line--pill" />
        <div class="mv-skeleton-line mv-skeleton-line--pill" />
        <div class="mv-skeleton-line mv-skeleton-line--pill" />
        <div class="mv-skeleton-actions" />
      </div>
    </div>

    <!-- Table -->
    <table v-else class="mv-table">
      <thead>
        <tr>
          <th v-if="hasExpandable(activeTab, items)" class="mv-th mv-th--expand" />
          <th class="mv-th mv-th--main">Judul</th>
          <th class="mv-th mv-th--type">Kategori</th>
          <th class="mv-th mv-th--hide-sm">Genre</th>
          <th class="mv-th mv-th--hide-sm">Rating</th>
          <th class="mv-th mv-th--hide-md">Tahun</th>
          <th class="mv-th mv-th--right">Aksi</th>
        </tr>
      </thead>
      <tbody>
        <template v-for="item in items" :key="item.id">
          <!-- Main row -->
          <tr class="mv-tr" :class="{ 'mv-tr--expanded': expandedRows.has(item.id) }">
            <td
              v-if="activeTab === 'series' || (activeTab === 'anime' && item.has_episodes)"
              class="mv-td mv-td--expand"
              @click="emit('toggle-expand', item)"
            >
              <button class="mv-expand-btn">
                <ChevronDown v-if="expandedRows.has(item.id)" :size="14" />
                <ChevronRight v-else :size="14" />
              </button>
            </td>

            <td class="mv-td mv-td--film">
              <div class="mv-film-cell">
                <div class="mv-poster-wrap">
                  <img :src="item.poster" class="mv-poster" :alt="item.title" />
                </div>
                <div>
                  <p class="mv-film-title">{{ item.title }}</p>
                  <p class="mv-film-meta">
                    <Clock :size="11" />{{ item.duration }}
                  </p>
                </div>
              </div>
            </td>

            <td class="mv-td mv-td--type">
              <span class="mv-type-badge"
                :class="{
                  'mv-type-badge--movie':  item.type === 'movie',
                  'mv-type-badge--series': item.type === 'series',
                  'mv-type-badge--anime':  item.type === 'anime',
                }"
              >
                <Film     v-if="item.type === 'movie'"  :size="10" />
                <Tv       v-else-if="item.type === 'series'" :size="10" />
                <Sparkles v-else :size="10" />
                {{ item.type === 'movie' ? 'Movie' : item.type === 'series' ? 'Series' : 'Anime' }}
              </span>
            </td>

            <td class="mv-td mv-td--hide-sm">
              <span class="vs-badge vs-badge--slate mv-genre-badge">{{ item.genre }}</span>
            </td>
            <td class="mv-td mv-td--hide-sm">
              <span class="vs-badge vs-badge--yellow">
                <Star :size="10" style="fill:currentColor" />{{ item.rating }}
              </span>
            </td>
            <td class="mv-td mv-td--year mv-td--hide-md">
              <span class="mv-year"><Calendar :size="11" />{{ item.year }}</span>
            </td>
            <td class="mv-td mv-td--actions">
              <div class="mv-actions">
                <button
                  v-if="activeTab === 'series' || (activeTab === 'anime' && item.has_episodes)"
                  class="mv-action-btn mv-action-btn--season"
                  title="Tambah Season"
                  @click="emit('toggle-expand', item); emit('open-add-season', item.id)"
                >
                  <Plus :size="13" />
                </button>
                <button
                  class="mv-action-btn mv-action-btn--subtitle"
                  title="Kelola Subtitle"
                  @click="emit('open-subtitle', item)"
                >
                  <Captions :size="14" />
                </button>
                <button class="mv-action-btn mv-action-btn--edit"   @click="emit('edit', item)"   title="Edit">
                  <Edit3 :size="14" />
                </button>
                <button class="mv-action-btn mv-action-btn--delete" @click="emit('delete', item)" title="Hapus">
                  <Trash2 :size="14" />
                </button>
              </div>
            </td>
          </tr>

          <!-- Expanded seasons -->
          <tr
            v-if="(activeTab === 'series' || (activeTab === 'anime' && item.has_episodes)) && expandedRows.has(item.id)"
            class="mv-tr-expanded"
          >
            <td :colspan="hasExpandable(activeTab, items) ? 7 : 6" class="mv-td-expanded">
              <div class="mv-seasons-wrap">
                <div v-if="seasonsLoading[item.id]" class="mv-seasons-empty">
                  <Loader2 :size="20" class="mv-spin" />
                  <p>Memuat season...</p>
                </div>

                <div v-else-if="!seasons[item.id] || seasons[item.id].length === 0" class="mv-seasons-empty">
                  <Tv :size="20" />
                  <p>Belum ada season.</p>
                  <button class="vs-btn vs-btn--ghost mv-btn-sm" @click="emit('open-add-season', item.id)">
                    <Plus :size="12" /> Tambah Season
                  </button>
                </div>

                <div v-else>
                  <div v-for="season in seasons[item.id]" :key="season.season_num" class="mv-season-block">
                    <div class="mv-season-header">
                      <span class="mv-season-label">Season {{ season.season_num }}</span>
                      <span class="mv-season-epcount">{{ season.episodes.length }} episode</span>
                    </div>
                    <table class="mv-ep-table">
                      <thead>
                        <tr>
                          <th class="mv-ep-th">Ep</th>
                          <th class="mv-ep-th">Judul Episode</th>
                          <th class="mv-ep-th">Server 1</th>
                          <th class="mv-ep-th">Server 2</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="ep in season.episodes" :key="ep.ep_num" class="mv-ep-tr">
                          <td class="mv-ep-td mv-ep-num">{{ ep.ep_num }}</td>
                          <td class="mv-ep-td">{{ ep.title || '—' }}</td>
                          <td class="mv-ep-td">
                            <span v-if="ep.url1" class="mv-ep-url-badge">
                              <PlayCircle :size="10" /> Server 1
                            </span>
                            <span v-else class="mv-ep-url-empty">—</span>
                          </td>
                          <td class="mv-ep-td">
                            <span v-if="ep.url2" class="mv-ep-url-badge">
                              <PlayCircle :size="10" /> Server 2
                            </span>
                            <span v-else class="mv-ep-url-empty">—</span>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                  <button class="vs-btn vs-btn--ghost mv-btn-sm mv-add-season-btn" @click="emit('open-add-season', item.id)">
                    <Plus :size="12" /> Tambah Season
                  </button>
                </div>
              </div>
            </td>
          </tr>
        </template>

        <!-- Empty -->
        <tr v-if="items.length === 0 && !loading">
          <td :colspan="hasExpandable(activeTab, items) ? 7 : 6" class="mv-empty">
            <Film     v-if="activeTab === 'movie'"  :size="32" />
            <Tv       v-else-if="activeTab === 'series'" :size="32" />
            <Sparkles v-else :size="32" />
            <p>Belum ada {{ activeTab === 'movie' ? 'movie' : activeTab === 'series' ? 'series' : 'anime' }}</p>
            <p class="mv-empty-sub">Klik "Tambah Konten" untuk menambahkan</p>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.mv-table-wrap { overflow-x: auto; }
.mv-table { width: 100%; border-collapse: collapse; }

.mv-th {
  padding: 10px 14px;
  font-size: 10.5px; font-weight: 600;
  letter-spacing: 0.07em; text-transform: uppercase;
  color: var(--muted2); text-align: left;
  border-bottom: 1px solid var(--border); white-space: nowrap;
}
.mv-th--expand { width: 36px; padding: 10px 6px 10px 14px; }
.mv-th--right  { text-align: right; }
.mv-th--type   { width: 110px; }

.mv-tr { transition: background 0.12s; }
.mv-tr:hover { background: rgba(255,255,255,0.02); }
.mv-tr:not(:last-child) .mv-td { border-bottom: 1px solid var(--border2); }
.mv-tr--expanded .mv-td { border-bottom: none !important; }

.mv-td { padding: 11px 14px; vertical-align: middle; }
.mv-td--expand { padding: 11px 6px 11px 14px; width: 36px; cursor: pointer; }
.mv-td--film   { min-width: 200px; }

.mv-expand-btn {
  width: 24px; height: 24px; border-radius: 5px;
  border: 1px solid var(--border); background: var(--surface2);
  color: var(--muted); display: grid; place-items: center;
  cursor: pointer; transition: all 0.15s;
}
.mv-expand-btn:hover { color: var(--text); }

.mv-film-cell { display: flex; align-items: center; gap: 12px; }
.mv-poster-wrap { width: 36px; height: 52px; border-radius: 6px; overflow: hidden; flex-shrink: 0; background: var(--surface2); }
.mv-poster { width: 100%; height: 100%; object-fit: cover; }
.mv-film-title { font-size: 13px; font-weight: 500; color: #fff; }
.mv-film-meta  { display: flex; align-items: center; gap: 4px; font-size: 11px; color: var(--muted); margin-top: 3px; }

.mv-type-badge {
  display: inline-flex; align-items: center; gap: 4px;
  padding: 3px 9px; border-radius: 20px;
  font-size: 11px; font-weight: 600; white-space: nowrap;
}
.mv-type-badge--movie  { background: rgba(99,102,241,0.12); color: #818cf8; border: 1px solid rgba(99,102,241,0.2); }
.mv-type-badge--series { background: rgba(6,182,212,0.10);  color: #22d3ee; border: 1px solid rgba(6,182,212,0.2); }
.mv-type-badge--anime  { background: rgba(244,63,94,0.10);  color: #fb7185; border: 1px solid rgba(244,63,94,0.2); }

.mv-genre-badge { max-width: 200px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; display: inline-block; }
.mv-year { display: flex; align-items: center; gap: 5px; font-size: 12px; color: var(--muted); }

.mv-actions { display: flex; align-items: center; justify-content: flex-end; gap: 4px; }
.mv-action-btn {
  width: 30px; height: 30px; border-radius: 7px;
  border: 1px solid var(--border); background: transparent;
  cursor: pointer; display: grid; place-items: center;
  transition: all 0.15s; color: var(--muted);
}
.mv-action-btn--subtitle:hover {
  color: #22d3ee;
  border-color: rgba(6,182,212,0.3);
  background: rgba(6,182,212,0.08);
}
.mv-action-btn:disabled { opacity: 0.3; cursor: not-allowed; }
.mv-action-btn--edit:hover   { color: var(--accent-hi); border-color: rgba(99,102,241,0.4); background: var(--accent-lo); }
.mv-action-btn--delete:hover { color: var(--danger); border-color: rgba(239,68,68,0.3); background: rgba(239,68,68,0.08); }
.mv-action-btn--season:hover { color: #22d3ee; border-color: rgba(6,182,212,0.3); background: rgba(6,182,212,0.08); }

@media (max-width: 700px) { .mv-td--hide-sm, .mv-th--hide-sm { display: none; } }
@media (max-width: 900px) { .mv-td--hide-md, .mv-th--hide-md { display: none; } }

.mv-tr-expanded { background: var(--surface2); }
.mv-td-expanded { padding: 0 14px 14px 52px; border-bottom: 1px solid var(--border) !important; }
.mv-seasons-wrap { padding-top: 12px; }
.mv-seasons-empty {
  display: flex; flex-direction: column; align-items: center;
  gap: 8px; padding: 24px; color: var(--muted2); font-size: 13px; text-align: center;
}
.mv-seasons-empty svg { opacity: 0.5; }

.mv-season-block { margin-bottom: 16px; }
.mv-season-header { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.mv-season-label  { font-size: 12px; font-weight: 700; color: #22d3ee; letter-spacing: 0.05em; text-transform: uppercase; }
.mv-season-epcount { font-size: 11px; color: var(--muted); background: var(--border); padding: 1px 7px; border-radius: 20px; }

.mv-ep-table { width: 100%; border-collapse: collapse; }
.mv-ep-th { font-size: 10px; font-weight: 600; letter-spacing: 0.07em; text-transform: uppercase; color: var(--muted2); padding: 6px 10px; text-align: left; border-bottom: 1px solid var(--border2); }
.mv-ep-tr:hover { background: rgba(255,255,255,0.02); }
.mv-ep-tr:not(:last-child) .mv-ep-td { border-bottom: 1px solid var(--border2); }
.mv-ep-td { padding: 7px 10px; font-size: 12px; color: var(--muted); vertical-align: middle; }
.mv-ep-num { width: 40px; font-weight: 700; color: var(--muted2); font-size: 11px; }
.mv-ep-url-badge {
  display: inline-flex; align-items: center; gap: 4px;
  font-size: 10.5px; color: var(--success);
  background: rgba(16,185,129,0.1); border: 1px solid rgba(16,185,129,0.2);
  padding: 2px 7px; border-radius: 4px; font-weight: 600;
}
.mv-ep-url-empty { color: var(--muted2); opacity: 0.5; }
.mv-add-season-btn { margin-top: 10px; }
.mv-btn-sm { padding: 5px 10px !important; font-size: 11.5px !important; }

.mv-empty { text-align: center; padding: 60px 20px; color: var(--muted2); font-size: 14px; font-weight: 500; }
.mv-empty svg { margin: 0 auto 12px; display: block; opacity: 0.5; }
.mv-empty-sub { font-size: 12px; opacity: 0.6; margin-top: 4px; font-weight: 400; }

.mv-loading { padding: 0; }
.mv-skeleton-row { display: flex; align-items: center; gap: 14px; padding: 14px 16px; border-bottom: 1px solid var(--border2); }
.mv-skeleton-poster { width: 36px; height: 52px; border-radius: 6px; flex-shrink: 0; background: linear-gradient(90deg, var(--surface2) 25%, var(--border) 50%, var(--surface2) 75%); background-size: 200%; animation: shimmer 1.4s infinite; }
.mv-skeleton-body { flex: 1; }
.mv-skeleton-line { border-radius: 4px; background: linear-gradient(90deg, var(--surface2) 25%, var(--border) 50%, var(--surface2) 75%); background-size: 200%; animation: shimmer 1.4s infinite; }
.mv-skeleton-line--title { width: 50%; height: 13px; margin-bottom: 8px; }
.mv-skeleton-line--sub   { width: 30%; height: 11px; }
.mv-skeleton-line--pill  { width: 64px; height: 22px; border-radius: 20px; flex-shrink: 0; }
.mv-skeleton-actions     { width: 64px; height: 30px; border-radius: 7px; flex-shrink: 0; background: linear-gradient(90deg, var(--surface2) 25%, var(--border) 50%, var(--surface2) 75%); background-size: 200%; animation: shimmer 1.4s infinite; }
@keyframes shimmer { to { background-position: -200% 0; } }

.mv-spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>