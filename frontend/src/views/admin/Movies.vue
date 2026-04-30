<script setup>
import { ref, onMounted } from 'vue'
import api from '@/lib/api'  // ← ganti axios
import { Plus, Film, Search, Loader2, PlayCircle, Star, Calendar, Trash2, Edit3 } from 'lucide-vue-next'

const movies = ref([])
const loading = ref(false)
const showModal = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const isSearchingTMDB = ref(false)
const previewData = ref(null)

const form = ref({ tmdb_id: '', url1: '', url2: '', url3: '' })

const fetchMovies = async () => {
  loading.value = true
  try {
    const res = await api.get('/admin/movies')
    movies.value = res.data.data
  } catch (err) { console.error(err) }
  finally { loading.value = false }
}

const openEditModal = (movie) => {
  isEdit.value = true
  currentId.value = movie.id
  form.value = { tmdb_id: movie.tmdb_id, url1: movie.url1, url2: movie.url2, url3: movie.url3 }
  showModal.value = true
}

const submitMovie = async () => {
  loading.value = true
  try {
    if (isEdit.value) {
      await api.put(`/admin/movies/${currentId.value}`, form.value)
    } else {
      await api.post('/admin/movies', form.value)
    }
    showModal.value = false
    resetForm()
    await fetchMovies()
  } catch (err) {
    alert(err.response?.data?.error || "Terjadi kesalahan")
  } finally {
    loading.value = false
  }
}

const deleteMovie = async (id) => {
  if (!confirm("Apakah Anda yakin ingin menghapus film ini?")) return
  try {
    await api.delete(`/admin/movies/${id}`)
    await fetchMovies()
  } catch (err) {
    alert("Gagal menghapus film")
  }
}

const searchTMDB = async () => {
  if (!form.value.tmdb_id) return
  isSearchingTMDB.value = true
  previewData.value = null
  try {
    const res = await api.get(`/admin/tmdb/${form.value.tmdb_id}`)
    previewData.value = res.data.data
  } catch (err) {
    alert("Film tidak ditemukan di TMDB. Cek kembali ID-nya.")
  } finally {
    isSearchingTMDB.value = false
  }
}

const resetForm = () => {
  isEdit.value = false
  currentId.value = null
  form.value = { tmdb_id: '', url1: '', url2: '', url3: '' }
  previewData.value = null
}

onMounted(fetchMovies)
</script>

<template>
  <div>
    <!-- Header -->
    <div class="vs-page-header">
      <div>
        <h1 class="vs-page-title">
          <Film :size="22" />
          Manajemen Film
        </h1>
        <p class="vs-page-subtitle">Kelola database film untuk platform V-STREAM</p>
      </div>
      <button class="vs-btn vs-btn--primary" @click="() => { resetForm(); showModal = true }">
        <Plus :size="15" />
        Tambah Film
      </button>
    </div>

    <!-- Table card -->
    <div class="vs-card">

      <!-- Toolbar -->
      <div class="mv-toolbar">
        <div class="mv-search-wrap">
          <Search :size="14" class="mv-search-icon" />
          <input type="text" placeholder="Cari film..." class="vs-input mv-search-input" />
        </div>
        <div class="mv-toolbar-right">
          <span class="mv-count">
            {{ movies.length }} film
          </span>
        </div>
      </div>

      <!-- Loading skeleton -->
      <div v-if="loading && movies.length === 0" class="mv-loading">
        <div v-for="i in 5" :key="i" class="mv-skeleton-row">
          <div class="mv-skeleton-poster" />
          <div class="mv-skeleton-body">
            <div class="mv-skeleton-line mv-skeleton-line--title" />
            <div class="mv-skeleton-line mv-skeleton-line--sub" />
          </div>
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
              <th class="mv-th mv-th--main">Film</th>
              <th class="mv-th mv-th--hide-sm">Genre</th>
              <th class="mv-th mv-th--hide-sm">Rating</th>
              <th class="mv-th mv-th--hide-md">Tahun</th>
              <th class="mv-th mv-th--right">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="movie in movies" :key="movie.id" class="mv-tr">
              <td class="mv-td mv-td--film">
                <div class="mv-film-cell">
                  <div class="mv-poster-wrap">
                    <img :src="movie.poster" class="mv-poster" :alt="movie.title" />
                  </div>
                  <div>
                    <p class="mv-film-title">{{ movie.title }}</p>
                    <p class="mv-film-meta">
                      <PlayCircle :size="11" />
                      {{ movie.duration }}
                    </p>
                  </div>
                </div>
              </td>
              <td class="mv-td mv-td--hide-sm">
                <span class="vs-badge vs-badge--slate">{{ movie.genre }}</span>
              </td>
              <td class="mv-td mv-td--hide-sm">
                <span class="vs-badge vs-badge--yellow">
                  <Star :size="10" style="fill:currentColor" />
                  {{ movie.rating }}
                </span>
              </td>
              <td class="mv-td mv-td--year mv-td--hide-md">
                <span class="mv-year">
                  <Calendar :size="11" />
                  {{ movie.year }}
                </span>
              </td>
              <td class="mv-td mv-td--actions">
                <div class="mv-actions">
                  <button class="mv-action-btn mv-action-btn--edit" @click="openEditModal(movie)" title="Edit">
                    <Edit3 :size="14" />
                  </button>
                  <button class="mv-action-btn mv-action-btn--delete" @click="deleteMovie(movie.id)" title="Hapus">
                    <Trash2 :size="14" />
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="movies.length === 0">
              <td colspan="5" class="mv-empty">
                <Film :size="32" />
                <p>Belum ada data film</p>
                <p class="mv-empty-sub">Tambahkan film pertama Anda</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

    </div>

    <!-- Modal -->
    <Teleport to="body">
      <Transition name="mv-modal">
        <div v-if="showModal" class="mv-backdrop" @click.self="showModal = false">
          <div class="mv-modal">

            <!-- Modal header -->
            <div class="mv-modal-header">
              <div class="mv-modal-title-wrap">
                <div class="mv-modal-icon">
                  <Film :size="16" />
                </div>
                <h2 class="mv-modal-title">{{ isEdit ? 'Edit Film' : 'Tambah Film Baru' }}</h2>
              </div>
              <button class="mv-modal-close" @click="showModal = false">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 6L6 18M6 6l12 12"/>
                </svg>
              </button>
            </div>

            <!-- Modal body -->
            <form @submit.prevent="submitMovie" class="mv-modal-body">

              <div class="mv-field">
                <label class="mv-label">TMDB ID <span class="mv-label-required">*</span></label>
                <div style="display: flex; gap: 8px;">
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
                    style="flex-shrink: 0;"
                  >
                    <Loader2 v-if="isSearchingTMDB" :size="14" class="mv-spin" />
                    <Search v-else :size="14" />
                    Cek
                  </button>
                </div>
                <p v-if="!isEdit" class="mv-field-hint">Cari ID film lalu klik Cek untuk melihat preview.</p>
              </div>

              <div v-if="previewData && !isEdit" class="mv-preview-box">
                <img 
                  :src="'https://image.tmdb.org/t/p/w200' + previewData.poster_path" 
                  class="mv-preview-poster" 
                  alt="Poster"
                />
                <div class="mv-preview-info">
                  <h4 class="mv-preview-title">{{ previewData.title }}</h4>
                  <div class="mv-preview-meta">
                    <Star :size="11" style="color: #f59e0b; fill: #f59e0b;" /> {{ previewData.vote_average.toFixed(1) }}
                    <span style="opacity: 0.5;">•</span>
                    <Calendar :size="11" /> {{ previewData.release_date ? previewData.release_date.substring(0,4) : '-' }}
                  </div>
                  <p class="mv-preview-desc">{{ previewData.overview }}</p>
                </div>
              </div>

              <div class="mv-field">
                <label class="mv-label">Video URL — Server 1 <span class="mv-label-required">*</span></label>
                <input v-model="form.url1" placeholder="https://..." required class="vs-input" />
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
                <button type="button" class="vs-btn vs-btn--ghost" @click="showModal = false">
                  Batal
                </button>
                <button type="submit" class="vs-btn vs-btn--primary" :disabled="loading">
                  <Loader2 v-if="loading" :size="14" class="mv-spin" />
                  <span>{{ loading ? 'Menyimpan...' : (isEdit ? 'Simpan Perubahan' : 'Tambah Film') }}</span>
                </button>
              </div>

            </form>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
/* ── Toolbar ────────────────────────────────────────────── */
.mv-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  border-bottom: 1px solid var(--border);
  flex-wrap: wrap;
}
.mv-search-wrap {
  position: relative;
  flex: 1;
  min-width: 200px;
  max-width: 320px;
}
.mv-search-icon {
  position: absolute;
  left: 12px; top: 50%;
  transform: translateY(-50%);
  color: var(--muted2);
  pointer-events: none;
}
.mv-search-input { padding-left: 36px !important; }

.mv-toolbar-right { display: flex; align-items: center; gap: 8px; margin-left: auto; }
.mv-count { font-size: 12px; color: var(--muted); }

/* ── Table ──────────────────────────────────────────────── */
.mv-table-wrap { overflow-x: auto; }
.mv-table { width: 100%; border-collapse: collapse; }

.mv-th {
  padding: 10px 16px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--muted2);
  text-align: left;
  border-bottom: 1px solid var(--border);
  white-space: nowrap;
}
.mv-th--right { text-align: right; }

.mv-tr { transition: background 0.12s; }
.mv-tr:hover { background: rgba(255,255,255,0.02); }
.mv-tr:not(:last-child) .mv-td { border-bottom: 1px solid var(--border2); }

.mv-td { padding: 12px 16px; vertical-align: middle; }
.mv-td--film { min-width: 220px; }

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
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: var(--muted);
  margin-top: 3px;
}

.mv-year { display: flex; align-items: center; gap: 5px; font-size: 12px; color: var(--muted); }

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
.mv-action-btn--edit:hover  { color: var(--accent-hi); border-color: rgba(99,102,241,0.4); background: var(--accent-lo); }
.mv-action-btn--delete:hover{ color: var(--danger); border-color: rgba(239,68,68,0.3); background: rgba(239,68,68,0.08); }

.mv-td--hide-sm { }
.mv-td--hide-md { }
.mv-th--hide-sm { }
.mv-th--hide-md { }

@media (max-width: 700px) {
  .mv-td--hide-sm, .mv-th--hide-sm { display: none; }
}
@media (max-width: 900px) {
  .mv-td--hide-md, .mv-th--hide-md { display: none; }
}

/* ── Empty ──────────────────────────────────────────────── */
.mv-empty {
  text-align: center;
  padding: 60px 20px;
  color: var(--muted2);
  font-size: 14px;
  font-weight: 500;
}
.mv-empty svg { margin: 0 auto 12px; display: block; }
.mv-empty-sub { font-size: 12px; color: var(--muted2); opacity: 0.6; margin-top: 4px; font-weight: 400; }

/* ── Skeleton ───────────────────────────────────────────── */
.mv-skeleton-row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border2);
}
.mv-skeleton-poster {
  width: 36px; height: 52px;
  border-radius: 6px;
  flex-shrink: 0;
  background: linear-gradient(90deg, var(--surface2) 25%, var(--border) 50%, var(--surface2) 75%);
  background-size: 200%;
  animation: shimmer 1.4s infinite;
}
.mv-skeleton-body { flex: 1; }
.mv-skeleton-line {
  border-radius: 4px;
  background: linear-gradient(90deg, var(--surface2) 25%, var(--border) 50%, var(--surface2) 75%);
  background-size: 200%;
  animation: shimmer 1.4s infinite;
}
.mv-skeleton-line--title { width: 50%; height: 13px; margin-bottom: 8px; }
.mv-skeleton-line--sub   { width: 30%; height: 11px; }
.mv-skeleton-line--pill  { width: 64px; height: 22px; border-radius: 20px; flex-shrink: 0; }
.mv-skeleton-actions     { width: 64px; height: 30px; border-radius: 7px;
  background: linear-gradient(90deg, var(--surface2) 25%, var(--border) 50%, var(--surface2) 75%);
  background-size: 200%;
  animation: shimmer 1.4s infinite;
  flex-shrink: 0;
}
@keyframes shimmer { to { background-position: -200% 0; } }

/* ── Modal ──────────────────────────────────────────────── */
.mv-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.65);
  backdrop-filter: blur(6px);
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}
.mv-modal {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 16px;
  width: 100%;
  max-width: 460px;
  box-shadow: 0 25px 60px rgba(0,0,0,0.5), 0 0 0 1px rgba(255,255,255,0.04);
  overflow: hidden;
}
.mv-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 20px;
  border-bottom: 1px solid var(--border);
}
.mv-modal-title-wrap { display: flex; align-items: center; gap: 10px; }
.mv-modal-icon {
  width: 32px; height: 32px;
  border-radius: 8px;
  background: var(--accent-lo);
  color: var(--accent-hi);
  display: grid; place-items: center;
}
.mv-modal-title { font-size: 14px; font-weight: 600; color: #fff; }
.mv-modal-close {
  width: 28px; height: 28px;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: transparent;
  color: var(--muted);
  display: grid; place-items: center;
  cursor: pointer;
  transition: all 0.15s;
}
.mv-modal-close:hover { color: #fff; border-color: var(--border2); background: var(--border2); }

.mv-modal-body { padding: 20px; display: flex; flex-direction: column; gap: 16px; }

.mv-field { display: flex; flex-direction: column; gap: 6px; }
.mv-field-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.mv-label { font-size: 12px; font-weight: 500; color: var(--muted); }
.mv-label-required { color: var(--danger); }
.mv-field-hint { font-size: 11px; color: var(--muted2); }

.mv-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding-top: 4px;
  border-top: 1px solid var(--border);
  padding-top: 16px;
}

.mv-spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* ── Modal transitions ──────────────────────────────────── */
.mv-modal-enter-active { animation: modal-in 0.22s cubic-bezier(0.34,1.56,0.64,1); }
.mv-modal-leave-active { animation: modal-out 0.18s ease; }
@keyframes modal-in  { from { opacity:0; transform:scale(0.94) translateY(10px); } }
@keyframes modal-out { to   { opacity:0; transform:scale(0.96) translateY(6px); } }

/* ── TMDB Preview Box ───────────────────────────────────── */
.mv-preview-box {
  display: flex;
  gap: 12px;
  background: var(--surface2);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 12px;
  margin-bottom: 4px;
}
.mv-preview-poster {
  width: 50px;
  height: 75px;
  border-radius: 4px;
  object-fit: cover;
  background: var(--border);
  flex-shrink: 0;
}
.mv-preview-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-width: 0;
}
.mv-preview-title {
  font-size: 13px;
  font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.mv-preview-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: var(--muted);
}
.mv-preview-desc {
  font-size: 11px;
  color: var(--muted2);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>