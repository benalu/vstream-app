<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/lib/api'
import { Plus, Search, Clapperboard, Film, Tv, Sparkles } from 'lucide-vue-next'
import MovieTable    from '@/components/admin/MovieTable.vue'
import MovieFormModal from '@/components/admin/MovieFormModal.vue'
import SeasonModal   from '@/components/admin/SeasonModal.vue'
import SubtitleModal  from '@/components/admin/SubtitleModal.vue'

// ── State ────────────────────────────────────────────────────
const allContent   = ref([])
const loading      = ref(false)
const activeTab    = ref('movie')
const searchQuery  = ref('')
const expandedRows = ref(new Set())
const seasons      = ref({})
const seasonsLoading = ref({})

const showFormModal   = ref(false)
const showSeasonModal = ref(false)
const isEdit          = ref(false)
const currentId       = ref(null)
const currentSeriesId = ref(null)

const formModalRef   = ref(null)
const seasonModalRef = ref(null)

const showSubtitleModal = ref(false)
const subtitleModalRef  = ref(null)

// ── Computed ─────────────────────────────────────────────────
const tabs = [
  { id: 'movie',  label: 'Movies', color: '#6366f1' },
  { id: 'series', label: 'Series', color: '#06b6d4' },
  { id: 'anime',  label: 'Anime',  color: '#f43f5e' },
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

// ── API ───────────────────────────────────────────────────────
const fetchAll = async () => {
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

const getEndpoint = (type) =>
  ({ movie: '/admin/movies', series: '/admin/series', anime: '/admin/anime' }[type] ?? '/admin/movies')

const handleSearchTMDB = async (tmdbId, type) => {
  try {
    const res = await api.get(`/admin/tmdb/${tmdbId}?type=${type}`)
    formModalRef.value?.setPreview(res.data.data)
  } catch {
    alert('Film tidak ditemukan di TMDB. Cek kembali ID-nya.')
  }
}

const handleFormSubmit = async (formData) => {
  loading.value = true
  try {
    const endpoint = getEndpoint(formData.type)
    if (isEdit.value) {
      await api.put(`${endpoint}/${currentId.value}`, formData)
    } else {
      await api.post(endpoint, formData)
    }
    showFormModal.value = false
    await fetchAll()
  } catch (err) {
    alert(err.response?.data?.error || 'Terjadi kesalahan')
  } finally {
    loading.value = false
  }
}

const handleEdit = (item) => {
  isEdit.value = true
  currentId.value = item.id
  formModalRef.value?.setForm({
    tmdb_id: item.tmdb_id, type: item.type,
    has_episodes: item.has_episodes || false,
    url1: item.url1, url2: item.url2, url3: item.url3,
  })
  showFormModal.value = true
}

const handleDelete = async (item) => {
  if (!confirm('Apakah Anda yakin ingin menghapus konten ini?')) return
  try {
    await api.delete(`${getEndpoint(item.type)}/${item.id}`)
    await fetchAll()
  } catch { alert('Gagal menghapus konten') }
}

const openAddModal = () => {
  isEdit.value = false
  currentId.value = null
  formModalRef.value?.reset()
  showFormModal.value = true
}

const handleOpenSubtitle = (item) => {
  showSubtitleModal.value = true
  // nextTick agar modal sudah mounted sebelum open dipanggil
  import('vue').then(({ nextTick }) => {
    nextTick(() => subtitleModalRef.value?.open(item.id, item.title))
  })
}

// ── Expand / Seasons ─────────────────────────────────────────
const handleToggleExpand = async (item) => {
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
      const res = await api.get(`${getEndpoint(item.type)}/${item.id}/seasons`)
      seasons.value[item.id] = res.data.data
    } catch {
      seasons.value[item.id] = []
    } finally {
      seasonsLoading.value[item.id] = false
    }
  }
}

const handleOpenAddSeason = (seriesId) => {
  currentSeriesId.value = seriesId
  const existing = seasons.value[seriesId] ?? []
  seasonModalRef.value?.reset(existing.length + 1)
  showSeasonModal.value = true
}

const handleSeasonSubmit = async (formData) => {
  try {
    const item = allContent.value.find(x => x.id === currentSeriesId.value)
    if (!item) return
    const endpoint = getEndpoint(item.type)
    await api.post(`${endpoint}/${currentSeriesId.value}/seasons`, formData)
    const res = await api.get(`${endpoint}/${currentSeriesId.value}/seasons`)
    seasons.value[currentSeriesId.value] = res.data.data
    showSeasonModal.value = false
  } catch (err) {
    alert(err.response?.data?.error || 'Gagal menambah season')
  }
}

onMounted(fetchAll)
</script>

<template>
  <div>
    <!-- Header -->
    <div class="vs-page-header">
      <div>
        <h1 class="vs-page-title"><Clapperboard :size="22" />Manajemen Konten</h1>
        <p class="vs-page-subtitle">Kelola Movies, Series & Anime untuk platform V-STREAM</p>
      </div>
      <button class="vs-btn vs-btn--primary" @click="openAddModal">
        <Plus :size="15" />Tambah Konten
      </button>
    </div>

    <!-- Card -->
    <div class="vs-card">
      <!-- Toolbar -->
      <div class="mv-toolbar">
        <div class="mv-tabs">
          <button
            v-for="tab in tabs" :key="tab.id"
            class="mv-tab"
            :class="{ 'mv-tab--active': activeTab === tab.id }"
            :style="activeTab === tab.id ? `--tab-color: ${tab.color}` : ''"
            @click="activeTab = tab.id; searchQuery = ''"
          >
            <Film     v-if="tab.id === 'movie'"  :size="13" />
            <Tv       v-else-if="tab.id === 'series'" :size="13" />
            <Sparkles v-else :size="13" />
            {{ tab.label }}
            <span class="mv-tab-count">{{ tabCount(tab.id) }}</span>
          </button>
        </div>
        <div class="mv-toolbar-right">
          <div class="mv-search-wrap">
            <Search :size="13" class="mv-search-icon" />
            <input v-model="searchQuery" type="text" placeholder="Cari judul..." class="vs-input mv-search-input" />
          </div>
          <span class="mv-count">{{ filteredContent.length }} konten</span>
        </div>
      </div>

      <!-- Table -->
      <MovieTable
        :items="filteredContent"
        :active-tab="activeTab"
        :seasons="seasons"
        :seasons-loading="seasonsLoading"
        :expanded-rows="expandedRows"
        :loading="loading"
        @edit="handleEdit"
        @delete="handleDelete"
        @toggle-expand="handleToggleExpand"
        @open-add-season="handleOpenAddSeason"
        @open-subtitle="handleOpenSubtitle"
      />
    </div>

    <!-- Modals -->
    <MovieFormModal
      ref="formModalRef"
      :show="showFormModal"
      :is-edit="isEdit"
      :active-tab="activeTab"
      :loading="loading"
      @close="showFormModal = false"
      @submit="handleFormSubmit"
      @search-tmdb="handleSearchTMDB"
    />

    <SeasonModal
      ref="seasonModalRef"
      :show="showSeasonModal"
      @close="showSeasonModal = false"
      @submit="handleSeasonSubmit"
    />

    <SubtitleModal
      ref="subtitleModalRef"
      :show="showSubtitleModal"
      @close="showSubtitleModal = false"
    />
  </div>
</template>

<style scoped>
.mv-toolbar {
  display: flex; align-items: center; gap: 12px;
  padding: 14px 18px; border-bottom: 1px solid var(--border); flex-wrap: wrap;
}
.mv-tabs {
  display: flex; gap: 4px;
  background: var(--surface2); border: 1px solid var(--border);
  border-radius: 10px; padding: 3px;
}
.mv-tab {
  display: flex; align-items: center; gap: 6px;
  padding: 6px 14px; border-radius: 7px;
  font-size: 12.5px; font-weight: 600; border: none;
  background: transparent; color: var(--muted);
  cursor: pointer; transition: all 0.15s; font-family: var(--font);
}
.mv-tab:hover { color: var(--text); }
.mv-tab--active {
  background: color-mix(in srgb, var(--tab-color) 15%, transparent);
  color: var(--tab-color);
  box-shadow: 0 0 0 1px color-mix(in srgb, var(--tab-color) 30%, transparent);
}
.mv-tab-count {
  font-size: 10px; font-weight: 700;
  padding: 1px 5px; border-radius: 20px;
  background: var(--border); color: var(--muted);
  min-width: 18px; text-align: center;
}
.mv-tab--active .mv-tab-count {
  background: color-mix(in srgb, var(--tab-color) 20%, transparent);
  color: var(--tab-color);
}
.mv-toolbar-right { display: flex; align-items: center; gap: 10px; margin-left: auto; }
.mv-search-wrap { position: relative; width: 220px; }
.mv-search-icon { position: absolute; left: 11px; top: 50%; transform: translateY(-50%); color: var(--muted2); pointer-events: none; }
.mv-search-input { padding-left: 34px !important; }
.mv-count { font-size: 12px; color: var(--muted); white-space: nowrap; }
</style>