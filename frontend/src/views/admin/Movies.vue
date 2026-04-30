<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/lib/api'
import {
  Plus, Film, Search, Loader2, PlayCircle, Star,
  Calendar, Trash2, Edit3, ChevronDown, ChevronRight,
  Tv, Clapperboard, Sparkles, Clock
} from 'lucide-vue-next'

// ── State ────────────────────────────────────────────────────
const allContent  = ref([])
const loading     = ref(false)
const showModal   = ref(false)
const isEdit      = ref(false)
const currentId   = ref(null)
const isSearchingTMDB = ref(false)
const previewData = ref(null)
const searchQuery = ref('')
const activeTab   = ref('movie') // 'movie', 'series', 'anime'
const expandedRows = ref(new Set()) // set of content IDs yang di-expand

const form = ref({ tmdb_id: '', type: 'movie', has_episodes: false, url1: '', url2: '', url3: '' })

// Season/episode management
const showSeasonModal  = ref(false)
const currentSeriesId  = ref(null)
const seasons          = ref({}) // { [contentId]: [{season_num, episodes:[]}] }
const seasonsLoading   = ref({})
const seasonForm       = ref({ season_num: 1, episodes: [{ ep_num: 1, title: '', url1: '', url2: '' }] })

// ── Computed ─────────────────────────────────────────────────
const tabs = [
  { id: 'movie', label: 'Movies',  icon: 'movie',  color: '#6366f1' },
  { id: 'series', label: 'Series',  icon: 'series', color: '#06b6d4' },
  { id: 'anime',   label: 'Anime',   icon: 'anime',  color: '#f43f5e' },
]

const filteredContent = computed(() => {
  let list = allContent.value.filter(m => m.type === activeTab.value)
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(m =>
      m.title.toLowerCase().includes(q) ||
      m.genre?.toLowerCase().includes(q)
    )
  }
  return list
})

const tabCount = (typeId) => allContent.value.filter(m => m.type === typeId).length
const hasExpandable = computed(() =>
  activeTab.value === 'series' ||
  (activeTab.value === 'anime' && filteredContent.value.some(i => i.has_episodes))
)

// ── API ───────────────────────────────────────────────────────
const fetchMovies = async () => {
  loading.value = true
  try {
    const [movRes, serRes, anmRes] = await Promise.all([
      api.get('/admin/movies'),
      api.get('/admin/series'),
      api.get('/admin/anime'),
    ])
    allContent.value = [
      ...(movRes.data.data || []),
      ...(serRes.data.data || []),
      ...(anmRes.data.data || []),
    ]
  } catch (err) { console.error(err) }
  finally { loading.value = false }
}

const searchTMDB = async () => {
  if (!form.value.tmdb_id) return
  isSearchingTMDB.value = true
  previewData.value = null
  try {
    const res = await api.get(`/admin/tmdb/${form.value.tmdb_id}?type=${form.value.type}`)
    previewData.value = res.data.data
  } catch {
    alert('Film tidak ditemukan di TMDB. Cek kembali ID-nya.')
  } finally {
    isSearchingTMDB.value = false
  }
}

const getEndpoint = (type) => {
  const endpoints = { movie: '/admin/movies', series: '/admin/series', anime: '/admin/anime' }
  return endpoints[type] || '/admin/movies'
}

const submitMovie = async () => {
  loading.value = true
  try {
    const endpoint = getEndpoint(form.value.type)
    if (isEdit.value) {
      await api.put(`${endpoint}/${currentId.value}`, form.value)
    } else {
      await api.post(endpoint, form.value)
    }
    showModal.value = false
    resetForm()
    await fetchMovies()
  } catch (err) {
    alert(err.response?.data?.error || 'Terjadi kesalahan')
  } finally {
    loading.value = false
  }
}

const deleteMovie = async (item) => {
  if (!confirm('Apakah Anda yakin ingin menghapus konten ini?')) return
  try {
    const endpoint = getEndpoint(item.type)
    await api.delete(`${endpoint}/${item.id}`)
    await fetchMovies()
  } catch {
    alert('Gagal menghapus konten')
  }
}

const openEditModal = (item) => {
  isEdit.value = true
  currentId.value = item.id
  form.value = {
    tmdb_id: item.tmdb_id,
    type: item.type,
    has_episodes: item.has_episodes || false,
    url1: item.url1,
    url2: item.url2,
    url3: item.url3
  }
  showModal.value = true
}

const resetForm = () => {
  isEdit.value = false
  currentId.value = null
  form.value = { tmdb_id: '', type: activeTab.value, has_episodes: false, url1: '', url2: '', url3: '' }
  previewData.value = null
}

// ── Expand row (series/anime) ────────────────────────────────
const toggleExpand = async (item) => {
  if (expandedRows.value.has(item.id)) {
    expandedRows.value.delete(item.id)
    expandedRows.value = new Set(expandedRows.value)
    return
  }
  expandedRows.value.add(item.id)
  expandedRows.value = new Set(expandedRows.value)

  if (!seasons.value[item.id]) {
    seasonsLoading.value[item.id] = true
    try {
      const endpoint = getEndpoint(item.type)
      const res = await api.get(`${endpoint}/${item.id}/seasons`)
      seasons.value[item.id] = res.data.data
    } catch {
      seasons.value[item.id] = []
    } finally {
      seasonsLoading.value[item.id] = false
    }
  }
}

const openAddSeasonModal = (seriesId) => {
  currentSeriesId.value = seriesId
  seasonForm.value = { season_num: 1, episodes: [{ ep_num: 1, title: '', url1: '', url2: '' }] }
  showSeasonModal.value = true
}

const addEpisode = () => {
  const last = seasonForm.value.episodes.at(-1)
  seasonForm.value.episodes.push({ ep_num: last.ep_num + 1, title: '', url1: '', url2: '' })
}

const removeEpisode = (i) => {
  if (seasonForm.value.episodes.length > 1)
    seasonForm.value.episodes.splice(i, 1)
}

const submitSeason = async () => {
  try {
    const item = allContent.value.find(x => x.id === currentSeriesId.value)
    if (!item) return
    
    const endpoint = getEndpoint(item.type)
    await api.post(`${endpoint}/${currentSeriesId.value}/seasons`, seasonForm.value)
    
    // Refresh seasons data
    const res = await api.get(`${endpoint}/${currentSeriesId.value}/seasons`)
    seasons.value[currentSeriesId.value] = res.data.data
    
    showSeasonModal.value = false
  } catch (err) {
    alert(err.response?.data?.error || 'Gagal menambah season')
  }
}

onMounted(fetchMovies)
</script>

<template>
  <div>
    <!-- ── Header ──────────────────────────────────────────── -->
    <div class="vs-page-header">
      <div>
        <h1 class="vs-page-title">
          <Clapperboard :size="22" />
          Manajemen Konten
        </h1>
        <p class="vs-page-subtitle">Kelola Movies, Series & Anime untuk platform V-STREAM</p>
      </div>
      <button class="vs-btn vs-btn--primary" @click="() => { resetForm(); showModal = true }">
        <Plus :size="15" />
        Tambah Konten
      </button>
    </div>

    <!-- ── Card ───────────────────────────────────────────── -->
    <div class="vs-card">

      <!-- Toolbar -->
      <div class="mv-toolbar">
        <!-- Tabs -->
        <div class="mv-tabs">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            class="mv-tab"
            :class="{ 'mv-tab--active': activeTab === tab.id }"
            :style="activeTab === tab.id ? `--tab-color: ${tab.color}` : ''"
            @click="activeTab = tab.id; searchQuery = ''"
          >
            <span class="mv-tab-icon">
              <Film v-if="tab.id === 'movie'" :size="13" />
              <Tv  v-else-if="tab.id === 'series'" :size="13" />
              <Sparkles v-else :size="13" />
            </span>
            {{ tab.label }}
            <span class="mv-tab-count">{{ tabCount(tab.id) }}</span>
          </button>
        </div>

        <!-- Search + count -->
        <div class="mv-toolbar-right">
          <div class="mv-search-wrap">
            <Search :size="13" class="mv-search-icon" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari judul..."
              class="vs-input mv-search-input"
            />
          </div>
          <span class="mv-count">{{ filteredContent.length }} konten</span>
        </div>
      </div>

      <!-- Skeleton -->
      <div v-if="loading && allContent.length === 0" class="mv-loading">
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
      <div v-else class="mv-table-wrap">
        <table class="mv-table">
          <thead>
            <tr>
              <th class="mv-th mv-th--expand" v-if="hasExpandable"></th>
              <th class="mv-th mv-th--main">Judul</th>
              <th class="mv-th mv-th--type">Kategori</th>
              <th class="mv-th mv-th--hide-sm">Genre</th>
              <th class="mv-th mv-th--hide-sm">Rating</th>
              <th class="mv-th mv-th--hide-md">Tahun</th>
              <th class="mv-th mv-th--right">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="item in filteredContent" :key="item.id">
              <!-- Main row -->
              <tr
                class="mv-tr"
                :class="{ 'mv-tr--expanded': expandedRows.has(item.id) }"
              >
                <!-- Expand toggle (series/anime only) -->
                <td v-if="activeTab === 'series' || (activeTab === 'anime' && item.has_episodes)"
                  class="mv-td mv-td--expand" @click="toggleExpand(item)">
                  <button class="mv-expand-btn">
                    <ChevronDown v-if="expandedRows.has(item.id)" :size="14" />
                    <ChevronRight v-else :size="14" />
                  </button>
                </td>

                <!-- Film cell -->
                <td class="mv-td mv-td--film">
                  <div class="mv-film-cell">
                    <div class="mv-poster-wrap">
                      <img :src="item.poster" class="mv-poster" :alt="item.title" />
                    </div>
                    <div>
                      <p class="mv-film-title">{{ item.title }}</p>
                      <p class="mv-film-meta">
                        <Clock :size="11" />
                        {{ item.duration }}
                      </p>
                    </div>
                  </div>
                </td>

                <!-- Type badge -->
                <td class="mv-td mv-td--type">
                  <span
                    class="mv-type-badge"
                    :class="{
                      'mv-type-badge--movie':  item.type === 'movie',
                      'mv-type-badge--series': item.type === 'series',
                      'mv-type-badge--anime':  item.type === 'anime',
                    }"
                  >
                    <Film      v-if="item.type === 'movie'" :size="10" />
                    <Tv        v-else-if="item.type === 'series'" :size="10" />
                    <Sparkles  v-else :size="10" />
                    {{ item.type === 'movie' ? 'Movie' : item.type === 'series' ? 'Series' : 'Anime' }}
                  </span>
                </td>

                <td class="mv-td mv-td--hide-sm">
                  <span class="vs-badge vs-badge--slate mv-genre-badge">{{ item.genre }}</span>
                </td>
                <td class="mv-td mv-td--hide-sm">
                  <span class="vs-badge vs-badge--yellow">
                    <Star :size="10" style="fill:currentColor" />
                    {{ item.rating }}
                  </span>
                </td>
                <td class="mv-td mv-td--year mv-td--hide-md">
                  <span class="mv-year">
                    <Calendar :size="11" />
                    {{ item.year }}
                  </span>
                </td>
                <td class="mv-td mv-td--actions">
                  <div class="mv-actions">
                    <!-- Add Season (series/anime only) -->
                    <button
                      v-if="activeTab === 'series' || (activeTab === 'anime' && item.has_episodes)"
                      class="mv-action-btn mv-action-btn--season"
                      title="Tambah Season"
                      @click="toggleExpand(item); openAddSeasonModal(item.id)"
                    >
                      <Plus :size="13" />
                    </button>
                    <button class="mv-action-btn mv-action-btn--edit" @click="openEditModal(item)" title="Edit">
                      <Edit3 :size="14" />
                    </button>
                    <button class="mv-action-btn mv-action-btn--delete" @click="deleteMovie(item)" title="Hapus">
                      <Trash2 :size="14" />
                    </button>
                  </div>
                </td>
              </tr>

              <!-- Expanded seasons row (series/anime) -->
              <tr v-if="(activeTab === 'series' || (activeTab === 'anime' && item.has_episodes)) && expandedRows.has(item.id)" class="mv-tr-expanded">
                <td :colspan="hasExpandable ? 7 : 6" class="mv-td-expanded">
                  <div class="mv-seasons-wrap">
                    <div v-if="seasonsLoading[item.id]" class="mv-seasons-empty">
                      <Loader2 :size="20" class="mv-spin" />
                      <p>Memuat season...</p>
                    </div>
                    <div v-if="!seasons[item.id] || seasons[item.id].length === 0" class="mv-seasons-empty">
                      <Tv :size="20" />
                      <p>Belum ada season.</p>
                      <button class="vs-btn vs-btn--ghost mv-btn-sm" @click="openAddSeasonModal(item.id)">
                        <Plus :size="12" /> Tambah Season
                      </button>
                    </div>

                    <div v-else>
                      <div
                        v-for="season in seasons[item.id]"
                        :key="season.season_num"
                        class="mv-season-block"
                      >
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
                              <td class="mv-ep-td mv-ep-url">
                                <span v-if="ep.url1" class="mv-ep-url-badge">
                                  <PlayCircle :size="10" /> Server 1
                                </span>
                                <span v-else class="mv-ep-url-empty">—</span>
                              </td>
                              <td class="mv-ep-td mv-ep-url">
                                <span v-if="ep.url2" class="mv-ep-url-badge">
                                  <PlayCircle :size="10" /> Server 2
                                </span>
                                <span v-else class="mv-ep-url-empty">—</span>
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>

                      <button class="vs-btn vs-btn--ghost mv-btn-sm mv-add-season-btn" @click="openAddSeasonModal(item.id)">
                        <Plus :size="12" /> Tambah Season
                      </button>
                    </div>

                  </div>
                </td>
              </tr>
            </template>

            <!-- Empty state -->
            <tr v-if="filteredContent.length === 0 && !loading">
              <td :colspan="activeTab !== 'movie' ? 7 : 6" class="mv-empty">
                <Film v-if="activeTab === 'movie'" :size="32" />
                <Tv   v-else-if="activeTab === 'series'" :size="32" />
                <Sparkles v-else :size="32" />
                <p>Belum ada {{ activeTab === 'movie' ? 'movie' : activeTab === 'series' ? 'series' : 'anime' }}</p>
                <p class="mv-empty-sub">Klik "Tambah Konten" untuk menambahkan</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ── Add/Edit Modal ────────────────────────────────────── -->
    <Teleport to="body">
      <Transition name="mv-modal">
        <div v-if="showModal" class="mv-backdrop" @click.self="showModal = false">
          <div class="mv-modal">

            <div class="mv-modal-header">
              <div class="mv-modal-title-wrap">
                <div class="mv-modal-icon">
                  <Film :size="16" />
                </div>
                <h2 class="mv-modal-title">{{ isEdit ? 'Edit Konten' : 'Tambah Konten Baru' }}</h2>
              </div>
              <button class="mv-modal-close" @click="showModal = false">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 6L6 18M6 6l12 12"/>
                </svg>
              </button>
            </div>

            <form @submit.prevent="submitMovie" class="mv-modal-body">

              <!-- Type selector -->
              <div class="mv-field">
                <label class="mv-label">Kategori <span class="mv-label-required">*</span></label>
                <div class="mv-type-selector">
                  <button
                    v-for="tab in tabs" :key="tab.id"
                    type="button"
                    class="mv-type-opt"
                    :class="{ 'mv-type-opt--active': form.type === tab.id }"
                    :style="form.type === tab.id ? `--opt-color: ${tab.color}` : ''"
                    @click="form.type = tab.id"
                  >
                    <Film     v-if="tab.id === 'movie'" :size="13" />
                    <Tv       v-else-if="tab.id === 'series'" :size="13" />
                    <Sparkles v-else :size="13" />
                    {{ tab.label }}
                  </button>
                </div>
              </div>

              <!-- Has Episodes (anime only) -->
              <div v-if="form.type === 'anime'" class="mv-field">
                <label class="mv-label">
                  <input type="checkbox" v-model="form.has_episodes" style="margin-right:6px" />
                  Anime dengan Episode/Season
                </label>
                <p class="mv-field-hint">Jika dipilih, URL diisi di level episode bukan konten</p>
              </div>

              <!-- TMDB ID -->
              <div class="mv-field">
                <label class="mv-label">TMDB ID <span class="mv-label-required">*</span></label>
                <div style="display:flex; gap:8px;">
                  <input
                    v-model="form.tmdb_id"
                    placeholder="Contoh: 550"
                    required
                    class="vs-input"
                    :disabled="isEdit"
                  />
                  <button
                    type="button"
                    class="vs-btn vs-btn--ghost"
                    @click="searchTMDB"
                    :disabled="!form.tmdb_id || isSearchingTMDB || isEdit"
                    style="flex-shrink:0"
                  >
                    <Loader2 v-if="isSearchingTMDB" :size="14" class="mv-spin" />
                    <Search  v-else :size="14" />
                    Cek
                  </button>
                </div>
                <p v-if="!isEdit" class="mv-field-hint">Masukkan ID TMDB lalu klik Cek untuk preview.</p>
              </div>

              <!-- Preview -->
              <div v-if="previewData && !isEdit" class="mv-preview-box">
                <img
                  :src="'https://image.tmdb.org/t/p/w200' + previewData.poster_path"
                  class="mv-preview-poster"
                  alt="Poster"
                />
                <div class="mv-preview-info">
                  <h4 class="mv-preview-title">{{ previewData.title }}</h4>
                  <div class="mv-preview-meta">
                    <Star :size="11" style="color:#f59e0b;fill:#f59e0b" />
                    {{ previewData.vote_average?.toFixed(1) }}
                    <span style="opacity:.4">•</span>
                    <Calendar :size="11" />
                    {{ previewData.release_date?.substring(0,4) || '-' }}
                  </div>
                  <p class="mv-preview-desc">{{ previewData.overview }}</p>
                </div>
              </div>

              <!-- URLs -->
              <div class="mv-field">
                <label class="mv-label">
                  Video URL — Server 1
                  <span v-if="form.type === 'movie' || (form.type === 'anime' && !form.has_episodes)" class="mv-label-required">*</span>
                </label>
                <input
                  v-model="form.url1"
                  placeholder="https://..."
                  :required="form.type === 'movie' || (form.type === 'anime' && !form.has_episodes)"
                  class="vs-input"
                />
              </div>

              <!-- url2 dan url3 selalu opsional, tampilkan sebagai field-row -->
              <div class="mv-field-row">
                <div class="mv-field">
                  <label class="mv-label">Server 2</label>
                  <input v-model="form.url2" placeholder="https://..." class="vs-input" />
                </div>
                <div class="mv-field">
                  <label class="mv-label">Server 3</label>
                  <input v-model="form.url3" placeholder="https://..." class="vs-input" />
                </div>
              </div>

              <div class="mv-field-row">
                <div class="mv-field">
                  <label class="mv-label">Server 2</label>
                  <input v-model="form.url2" placeholder="https://..." class="vs-input" />
                </div>
                <div class="mv-field">
                  <label class="mv-label">Server 3</label>
                  <input v-model="form.url3" placeholder="https://..." class="vs-input" />
                </div>
              </div>

              <div class="mv-modal-footer">
                <button type="button" class="vs-btn vs-btn--ghost" @click="showModal = false">Batal</button>
                <button type="submit" class="vs-btn vs-btn--primary" :disabled="loading">
                  <Loader2 v-if="loading" :size="14" class="mv-spin" />
                  <span>{{ loading ? 'Menyimpan...' : (isEdit ? 'Simpan Perubahan' : 'Tambah Konten') }}</span>
                </button>
              </div>

            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- ── Season Modal ──────────────────────────────────────── -->
    <Teleport to="body">
      <Transition name="mv-modal">
        <div v-if="showSeasonModal" class="mv-backdrop" @click.self="showSeasonModal = false">
          <div class="mv-modal mv-modal--wide">

            <div class="mv-modal-header">
              <div class="mv-modal-title-wrap">
                <div class="mv-modal-icon" style="background:rgba(6,182,212,0.12);color:#22d3ee">
                  <Tv :size="16" />
                </div>
                <h2 class="mv-modal-title">Tambah Season</h2>
              </div>
              <button class="mv-modal-close" @click="showSeasonModal = false">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 6L6 18M6 6l12 12"/>
                </svg>
              </button>
            </div>

            <div class="mv-modal-body">
              <!-- Season number -->
              <div class="mv-field">
                <label class="mv-label">Nomor Season</label>
                <input v-model.number="seasonForm.season_num" type="number" min="1" class="vs-input" style="max-width:120px" />
              </div>

              <!-- Episodes -->
              <div class="mv-field">
                <div class="mv-ep-list-header">
                  <label class="mv-label">Daftar Episode</label>
                  <button type="button" class="vs-btn vs-btn--ghost mv-btn-sm" @click="addEpisode">
                    <Plus :size="12" /> Tambah Episode
                  </button>
                </div>

                <div class="mv-ep-list">
                  <div
                    v-for="(ep, i) in seasonForm.episodes"
                    :key="i"
                    class="mv-ep-item"
                  >
                    <div class="mv-ep-num-badge">{{ ep.ep_num }}</div>
                    <div class="mv-ep-inputs">
                      <input v-model="ep.title" :placeholder="`Judul Episode ${ep.ep_num}`" class="vs-input mv-ep-title-input" />
                      <div class="mv-ep-urls">
                        <input v-model="ep.url1" placeholder="URL Server 1" class="vs-input" />
                        <input v-model="ep.url2" placeholder="URL Server 2" class="vs-input" />
                      </div>
                    </div>
                    <button
                      type="button"
                      class="mv-action-btn mv-action-btn--delete"
                      @click="removeEpisode(i)"
                      :disabled="seasonForm.episodes.length === 1"
                    >
                      <Trash2 :size="13" />
                    </button>
                  </div>
                </div>
              </div>

              <div class="mv-modal-footer">
                <button type="button" class="vs-btn vs-btn--ghost" @click="showSeasonModal = false">Batal</button>
                <button type="button" class="vs-btn vs-btn--primary" @click="submitSeason">
                  Simpan Season
                </button>
              </div>
            </div>

          </div>
        </div>
      </Transition>
    </Teleport>

  </div>
</template>

<style scoped>
/* ── Tabs ───────────────────────────────────────────────── */
.mv-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  border-bottom: 1px solid var(--border);
  flex-wrap: wrap;
}

.mv-tabs {
  display: flex;
  gap: 4px;
  background: var(--surface2);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 3px;
}

.mv-tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 7px;
  font-size: 12.5px;
  font-weight: 600;
  border: none;
  background: transparent;
  color: var(--muted);
  cursor: pointer;
  transition: all 0.15s;
  font-family: var(--font);
}
.mv-tab:hover { color: var(--text); }
.mv-tab--active {
  background: color-mix(in srgb, var(--tab-color) 15%, transparent);
  color: var(--tab-color);
  box-shadow: 0 0 0 1px color-mix(in srgb, var(--tab-color) 30%, transparent);
}
.mv-tab-icon { display: grid; place-items: center; }
.mv-tab-count {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 20px;
  background: var(--border);
  color: var(--muted);
  min-width: 18px;
  text-align: center;
}
.mv-tab--active .mv-tab-count {
  background: color-mix(in srgb, var(--tab-color) 20%, transparent);
  color: var(--tab-color);
}

/* ── Search ─────────────────────────────────────────────── */
.mv-toolbar-right {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-left: auto;
}
.mv-search-wrap {
  position: relative;
  width: 220px;
}
.mv-search-icon {
  position: absolute;
  left: 11px; top: 50%;
  transform: translateY(-50%);
  color: var(--muted2);
  pointer-events: none;
}
.mv-search-input { padding-left: 34px !important; }
.mv-count { font-size: 12px; color: var(--muted); white-space: nowrap; }

/* ── Table ──────────────────────────────────────────────── */
.mv-table-wrap { overflow-x: auto; }
.mv-table { width: 100%; border-collapse: collapse; }

.mv-th {
  padding: 10px 14px;
  font-size: 10.5px;
  font-weight: 600;
  letter-spacing: 0.07em;
  text-transform: uppercase;
  color: var(--muted2);
  text-align: left;
  border-bottom: 1px solid var(--border);
  white-space: nowrap;
}
.mv-th--expand { width: 36px; padding: 10px 6px 10px 14px; }
.mv-th--right  { text-align: right; }
.mv-th--type   { width: 110px; }

.mv-tr { transition: background 0.12s; }
.mv-tr:hover { background: rgba(255,255,255,0.02); }
.mv-tr:not(:last-child) .mv-td { border-bottom: 1px solid var(--border2); }
.mv-tr--expanded .mv-td { border-bottom: none !important; }

.mv-td { padding: 11px 14px; vertical-align: middle; }
.mv-td--expand {
  padding: 11px 6px 11px 14px;
  width: 36px;
  cursor: pointer;
}
.mv-td--film { min-width: 200px; }

.mv-expand-btn {
  width: 24px; height: 24px;
  border-radius: 5px;
  border: 1px solid var(--border);
  background: var(--surface2);
  color: var(--muted);
  display: grid; place-items: center;
  cursor: pointer;
  transition: all 0.15s;
}
.mv-expand-btn:hover { color: var(--text); border-color: var(--border2); }

/* Film cell */
.mv-film-cell { display: flex; align-items: center; gap: 12px; }
.mv-poster-wrap {
  width: 36px; height: 52px;
  border-radius: 6px;
  overflow: hidden;
  flex-shrink: 0;
  background: var(--surface2);
}
.mv-poster { width: 100%; height: 100%; object-fit: cover; }
.mv-film-title { font-size: 13px; font-weight: 500; color: #fff; }
.mv-film-meta {
  display: flex; align-items: center; gap: 4px;
  font-size: 11px; color: var(--muted); margin-top: 3px;
}

/* Type badge */
.mv-type-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 9px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
  white-space: nowrap;
}
.mv-type-badge--movie  { background: rgba(99,102,241,0.12); color: #818cf8; border: 1px solid rgba(99,102,241,0.2); }
.mv-type-badge--series { background: rgba(6,182,212,0.10);  color: #22d3ee; border: 1px solid rgba(6,182,212,0.2); }
.mv-type-badge--anime  { background: rgba(244,63,94,0.10);  color: #fb7185; border: 1px solid rgba(244,63,94,0.2); }

.mv-genre-badge {
  max-width: 200px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: inline-block;
}

.mv-year { display: flex; align-items: center; gap: 5px; font-size: 12px; color: var(--muted); }

/* Actions */
.mv-actions { display: flex; align-items: center; justify-content: flex-end; gap: 4px; }
.mv-action-btn {
  width: 30px; height: 30px;
  border-radius: 7px;
  border: 1px solid var(--border);
  background: transparent;
  cursor: pointer;
  display: grid; place-items: center;
  transition: all 0.15s;
  color: var(--muted);
}
.mv-action-btn:disabled { opacity: 0.3; cursor: not-allowed; }
.mv-action-btn--edit:hover   { color: var(--accent-hi); border-color: rgba(99,102,241,0.4); background: var(--accent-lo); }
.mv-action-btn--delete:hover { color: var(--danger); border-color: rgba(239,68,68,0.3); background: rgba(239,68,68,0.08); }
.mv-action-btn--season:hover { color: #22d3ee; border-color: rgba(6,182,212,0.3); background: rgba(6,182,212,0.08); }

@media (max-width: 700px) { .mv-td--hide-sm, .mv-th--hide-sm { display: none; } }
@media (max-width: 900px) { .mv-td--hide-md, .mv-th--hide-md { display: none; } }

/* ── Expanded Season Row ────────────────────────────────── */
.mv-tr-expanded { background: var(--surface2); }
.mv-td-expanded { padding: 0 14px 14px 52px; border-bottom: 1px solid var(--border) !important; }

.mv-seasons-wrap { padding-top: 12px; }
.mv-seasons-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px;
  color: var(--muted2);
  font-size: 13px;
  text-align: center;
}
.mv-seasons-empty svg { color: var(--muted2); opacity: 0.5; }

.mv-season-block { margin-bottom: 16px; }
.mv-season-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.mv-season-label {
  font-size: 12px;
  font-weight: 700;
  color: #22d3ee;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}
.mv-season-epcount {
  font-size: 11px;
  color: var(--muted);
  background: var(--border);
  padding: 1px 7px;
  border-radius: 20px;
}

.mv-ep-table { width: 100%; border-collapse: collapse; }
.mv-ep-th {
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.07em;
  text-transform: uppercase;
  color: var(--muted2);
  padding: 6px 10px;
  text-align: left;
  border-bottom: 1px solid var(--border2);
}
.mv-ep-tr { transition: background 0.1s; }
.mv-ep-tr:hover { background: rgba(255,255,255,0.02); }
.mv-ep-tr:not(:last-child) .mv-ep-td { border-bottom: 1px solid var(--border2); }
.mv-ep-td { padding: 7px 10px; font-size: 12px; color: var(--muted); vertical-align: middle; }
.mv-ep-num {
  width: 40px;
  font-weight: 700;
  color: var(--muted2);
  font-size: 11px;
}
.mv-ep-url-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 10.5px;
  color: var(--success);
  background: rgba(16,185,129,0.1);
  border: 1px solid rgba(16,185,129,0.2);
  padding: 2px 7px;
  border-radius: 4px;
  font-weight: 600;
}
.mv-ep-url-empty { color: var(--muted2); opacity: 0.5; }

.mv-add-season-btn { margin-top: 10px; }
.mv-btn-sm { padding: 5px 10px !important; font-size: 11.5px !important; }

/* ── Empty ──────────────────────────────────────────────── */
.mv-empty {
  text-align: center;
  padding: 60px 20px;
  color: var(--muted2);
  font-size: 14px;
  font-weight: 500;
}
.mv-empty svg { margin: 0 auto 12px; display: block; opacity: 0.5; }
.mv-empty-sub { font-size: 12px; opacity: 0.6; margin-top: 4px; font-weight: 400; }

/* ── Skeleton ───────────────────────────────────────────── */
.mv-loading { padding: 0; }
.mv-skeleton-row {
  display: flex; align-items: center; gap: 14px;
  padding: 14px 16px; border-bottom: 1px solid var(--border2);
}
.mv-skeleton-poster {
  width: 36px; height: 52px; border-radius: 6px; flex-shrink: 0;
  background: linear-gradient(90deg, var(--surface2) 25%, var(--border) 50%, var(--surface2) 75%);
  background-size: 200%; animation: shimmer 1.4s infinite;
}
.mv-skeleton-body { flex: 1; }
.mv-skeleton-line {
  border-radius: 4px;
  background: linear-gradient(90deg, var(--surface2) 25%, var(--border) 50%, var(--surface2) 75%);
  background-size: 200%; animation: shimmer 1.4s infinite;
}
.mv-skeleton-line--title { width: 50%; height: 13px; margin-bottom: 8px; }
.mv-skeleton-line--sub   { width: 30%; height: 11px; }
.mv-skeleton-line--pill  { width: 64px; height: 22px; border-radius: 20px; flex-shrink: 0; }
.mv-skeleton-actions {
  width: 64px; height: 30px; border-radius: 7px; flex-shrink: 0;
  background: linear-gradient(90deg, var(--surface2) 25%, var(--border) 50%, var(--surface2) 75%);
  background-size: 200%; animation: shimmer 1.4s infinite;
}
@keyframes shimmer { to { background-position: -200% 0; } }

/* ── Modal ──────────────────────────────────────────────── */
.mv-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.65);
  backdrop-filter: blur(6px);
  z-index: 100;
  display: flex; align-items: center; justify-content: center;
  padding: 16px;
}
.mv-modal {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 16px;
  width: 100%; max-width: 460px;
  box-shadow: 0 25px 60px rgba(0,0,0,0.5), 0 0 0 1px rgba(255,255,255,0.04);
  overflow: hidden;
  max-height: 90vh;
  display: flex; flex-direction: column;
}
.mv-modal--wide { max-width: 600px; }

.mv-modal-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 18px 20px; border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.mv-modal-title-wrap { display: flex; align-items: center; gap: 10px; }
.mv-modal-icon {
  width: 32px; height: 32px; border-radius: 8px;
  background: var(--accent-lo); color: var(--accent-hi);
  display: grid; place-items: center;
}
.mv-modal-title { font-size: 14px; font-weight: 600; color: #fff; }
.mv-modal-close {
  width: 28px; height: 28px; border-radius: 6px;
  border: 1px solid var(--border); background: transparent;
  color: var(--muted); display: grid; place-items: center;
  cursor: pointer; transition: all 0.15s;
}
.mv-modal-close:hover { color: #fff; border-color: var(--border2); background: var(--border2); }

.mv-modal-body {
  padding: 20px; display: flex; flex-direction: column; gap: 16px;
  overflow-y: auto;
}
.mv-field { display: flex; flex-direction: column; gap: 6px; }
.mv-field-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.mv-label { font-size: 12px; font-weight: 500; color: var(--muted); }
.mv-label-required { color: var(--danger); }
.mv-field-hint { font-size: 11px; color: var(--muted2); }

.mv-modal-footer {
  display: flex; justify-content: flex-end; gap: 8px;
  padding-top: 16px; border-top: 1px solid var(--border);
}

/* Type selector */
.mv-type-selector { display: flex; gap: 6px; }
.mv-type-opt {
  flex: 1; display: flex; align-items: center; justify-content: center; gap: 6px;
  padding: 8px 12px; border-radius: 8px;
  border: 1px solid var(--border);
  background: var(--surface2); color: var(--muted);
  font-size: 12.5px; font-weight: 600; cursor: pointer;
  transition: all 0.15s; font-family: var(--font);
}
.mv-type-opt:hover { color: var(--text); }
.mv-type-opt--active {
  background: color-mix(in srgb, var(--opt-color) 12%, transparent);
  color: var(--opt-color);
  border-color: color-mix(in srgb, var(--opt-color) 35%, transparent);
}

/* Season modal episode list */
.mv-ep-list-header { display: flex; align-items: center; justify-content: space-between; }
.mv-ep-list { display: flex; flex-direction: column; gap: 10px; margin-top: 8px; }
.mv-ep-item {
  display: flex; align-items: flex-start; gap: 10px;
  background: var(--surface2); border: 1px solid var(--border);
  border-radius: 8px; padding: 10px 12px;
}
.mv-ep-num-badge {
  width: 26px; height: 26px; border-radius: 6px;
  background: var(--border); color: var(--muted);
  display: grid; place-items: center;
  font-size: 11px; font-weight: 700; flex-shrink: 0;
}
.mv-ep-inputs { flex: 1; display: flex; flex-direction: column; gap: 7px; }
.mv-ep-title-input { font-weight: 500 !important; }
.mv-ep-urls { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }

/* Preview box */
.mv-preview-box {
  display: flex; gap: 12px;
  background: var(--surface2); border: 1px solid var(--border);
  border-radius: var(--radius-sm); padding: 12px; margin-bottom: 4px;
}
.mv-preview-poster { width: 50px; height: 75px; border-radius: 4px; object-fit: cover; flex-shrink: 0; }
.mv-preview-info { display: flex; flex-direction: column; gap: 4px; flex: 1; min-width: 0; }
.mv-preview-title { font-size: 13px; font-weight: 600; color: #fff; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.mv-preview-meta { display: flex; align-items: center; gap: 6px; font-size: 11px; color: var(--muted); }
.mv-preview-desc { font-size: 11px; color: var(--muted2); line-height: 1.4; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }

.mv-spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.mv-modal-enter-active { animation: modal-in 0.22s cubic-bezier(0.34,1.56,0.64,1); }
.mv-modal-leave-active { animation: modal-out 0.18s ease; }
@keyframes modal-in  { from { opacity:0; transform:scale(0.94) translateY(10px); } }
@keyframes modal-out { to   { opacity:0; transform:scale(0.96) translateY(6px); } }
</style>