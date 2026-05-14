<script setup>
import { ref, watch, computed } from 'vue'
import {
  Film, Search, Loader2, Star, Calendar, Tv, Sparkles,
  RefreshCw, Link2, Info, Image, ChevronRight
} from 'lucide-vue-next'

const props = defineProps({
  show:      { type: Boolean, required: true },
  isEdit:    { type: Boolean, default: false },
  activeTab: { type: String,  default: 'movie' },
  loading:   { type: Boolean, default: false },
})

const emit = defineEmits(['close', 'submit', 'search-tmdb'])

const tabs = [
  { id: 'movie',  label: 'Movies', color: '#6366f1' },
  { id: 'series', label: 'Series', color: '#06b6d4' },
  { id: 'anime',  label: 'Anime',  color: '#f43f5e' },
]

// Tab di dalam modal saat mode edit
const editTab = ref('meta') // 'meta' | 'links'

const form = ref({
  tmdb_id: '', type: 'movie', has_episodes: false,
  url1: '', url2: '', url3: '',
  // metadata fields (edit mode)
  title: '', overview: '', poster: '', backdrop: '', genre: '', year: '', duration: '', rating: '',
})

const previewData     = ref(null)
const isSearchingTMDB = ref(false)
const isRefreshing    = ref(false)
const tmdbIdError     = ref('')

// Sync type ke activeTab saat modal dibuka (add mode)
watch(() => props.show, (val) => {
  if (val && !props.isEdit) {
    form.value.type = props.activeTab
    editTab.value = 'meta'
  }
})

const validateTmdbId = () => {
  const val = form.value.tmdb_id
  if (!val)               { tmdbIdError.value = ''; return }
  if (!/^\d+$/.test(val)) { tmdbIdError.value = 'TMDB ID hanya boleh berisi angka'; return }
  if (val.length > 10)    { tmdbIdError.value = 'TMDB ID maksimal 10 digit'; return }
  tmdbIdError.value = ''
}

const searchTMDB = async () => {
  if (!form.value.tmdb_id || tmdbIdError.value) return
  isSearchingTMDB.value = true
  previewData.value = null
  try {
    const result = await emit('search-tmdb', form.value.tmdb_id, form.value.type)
    previewData.value = result
  } finally {
    isSearchingTMDB.value = false
  }
}

// Refresh metadata dari TMDB (edit mode)
const refreshFromTMDB = async () => {
  isRefreshing.value = true
  try {
    emit('submit', { ...form.value, refresh_from_tmdb: true })
  } finally {
    isRefreshing.value = false
  }
}

const handleSubmit = () => {
  validateTmdbId()
  if (tmdbIdError.value) return
  emit('submit', { ...form.value })
}

// Poster preview URL — pakai custom kalau diisi, fallback ke TMDB
const posterPreviewUrl = computed(() => {
  if (form.value.poster) return form.value.poster
  return null
})

// Expose untuk parent bisa set form (saat edit)
const setForm = (data) => {
  form.value = {
    tmdb_id:     data.tmdb_id      || '',
    type:        data.type         || 'movie',
    has_episodes:data.has_episodes || false,
    url1:        data.url1         || '',
    url2:        data.url2         || '',
    url3:        data.url3         || '',
    title:       data.title        || '',
    overview:    data.overview     || '',
    poster:      data.poster       || '',
    backdrop:    data.backdrop     || '',
    genre:       data.genre        || '',
    year:        data.year         || '',
    duration:    data.duration     || '',
    rating:      data.rating       || '',
  }
  editTab.value = 'meta'
}

const setPreview = (data) => { previewData.value = data }
const reset = () => {
  form.value = {
    tmdb_id: '', type: props.activeTab, has_episodes: false,
    url1: '', url2: '', url3: '',
    title: '', overview: '', poster: '', backdrop: '', genre: '', year: '', duration: '', rating: '',
  }
  previewData.value = null
  tmdbIdError.value = ''
  editTab.value = 'meta'
}

defineExpose({ setForm, setPreview, reset })
</script>

<template>
  <Teleport to="body">
    <Transition name="mv-modal">
      <div v-if="show" class="mv-backdrop" @click.self="emit('close')">
        <div class="mv-modal" :class="{ 'mv-modal--wide': isEdit }">

          <!-- Header -->
          <div class="mv-modal-header">
            <div class="mv-modal-title-wrap">
              <div class="mv-modal-icon">
                <Film :size="16" />
              </div>
              <div>
                <h2 class="mv-modal-title">{{ isEdit ? 'Edit Konten' : 'Tambah Konten Baru' }}</h2>
                <p v-if="isEdit" class="mv-modal-subtitle">TMDB ID: {{ form.tmdb_id }} · {{ form.type }}</p>
              </div>
            </div>
            <button class="mv-modal-close" @click="emit('close')">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <!-- Edit mode: inner tab bar -->
          <div v-if="isEdit" class="mv-edit-tabs">
            <button
              class="mv-edit-tab"
              :class="{ 'mv-edit-tab--active': editTab === 'meta' }"
              @click="editTab = 'meta'"
            >
              <Info :size="13" />Metadata
            </button>
            <button
              class="mv-edit-tab"
              :class="{ 'mv-edit-tab--active': editTab === 'media' }"
              @click="editTab = 'media'"
            >
              <Image :size="13" />Gambar
            </button>
            <button
              class="mv-edit-tab"
              :class="{ 'mv-edit-tab--active': editTab === 'links' }"
              @click="editTab = 'links'"
            >
              <Link2 :size="13" />Video Links
            </button>
          </div>

          <div class="mv-modal-body">

            <!-- ══════════════════ ADD MODE ══════════════════ -->
            <template v-if="!isEdit">
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
                    <Film     v-if="tab.id === 'movie'"   :size="13" />
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
                  <div style="flex:1">
                    <input
                      v-model="form.tmdb_id"
                      placeholder="Contoh: 550"
                      required inputmode="numeric" pattern="[0-9]{1,10}" maxlength="10"
                      class="vs-input" :class="{ 'vs-input--error': tmdbIdError }"
                      @input="validateTmdbId"
                    />
                    <p v-if="tmdbIdError" class="mv-field-error">{{ tmdbIdError }}</p>
                  </div>
                  <button
                    type="button" class="vs-btn vs-btn--ghost"
                    @click="searchTMDB"
                    :disabled="!form.tmdb_id || !!tmdbIdError || isSearchingTMDB"
                    style="flex-shrink:0; align-self:flex-start"
                  >
                    <Loader2 v-if="isSearchingTMDB" :size="14" class="mv-spin" />
                    <Search  v-else :size="14" />
                    Cek
                  </button>
                </div>
                <p class="mv-field-hint">Masukkan ID TMDB lalu klik Cek untuk preview.</p>
              </div>

              <!-- Preview -->
              <div v-if="previewData" class="mv-preview-box">
                <img :src="'https://image.tmdb.org/t/p/w200' + previewData.poster_path" class="mv-preview-poster" alt="Poster" />
                <div class="mv-preview-info">
                  <h4 class="mv-preview-title">{{ previewData.title || previewData.name }}</h4>
                  <div class="mv-preview-meta">
                    <Star :size="11" style="color:#f59e0b;fill:#f59e0b" />
                    {{ previewData.vote_average?.toFixed(1) }}
                    <span style="opacity:.4">•</span>
                    <Calendar :size="11" />
                    {{ (previewData.release_date || previewData.first_air_date)?.substring(0,4) || '-' }}
                  </div>
                  <p class="mv-preview-desc">{{ previewData.overview }}</p>
                </div>
              </div>

              <!-- URL fields -->
              <div class="mv-field">
                <label class="mv-label">
                  Video URL — Server 1
                  <span v-if="form.type === 'movie' || (form.type === 'anime' && !form.has_episodes)" class="mv-label-required">*</span>
                </label>
                <input v-model="form.url1" placeholder="https://..." class="vs-input" />
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
            </template>

            <!-- ══════════════════ EDIT MODE ══════════════════ -->
            <template v-else>

              <!-- ── Tab: Metadata ── -->
              <template v-if="editTab === 'meta'">
                <div class="mv-refresh-banner">
                  <div class="mv-refresh-info">
                    <RefreshCw :size="13" />
                    <span>Sync ulang semua metadata dari TMDB (judul, genre, rating, dll)</span>
                  </div>
                  <button
                    type="button"
                    class="vs-btn vs-btn--ghost mv-refresh-btn"
                    @click="refreshFromTMDB"
                    :disabled="loading || isRefreshing"
                  >
                    <Loader2 v-if="isRefreshing || loading" :size="13" class="mv-spin" />
                    <RefreshCw v-else :size="13" />
                    Refresh TMDB
                  </button>
                </div>

                <div class="mv-field">
                  <label class="mv-label">Judul <span class="mv-label-required">*</span></label>
                  <input v-model="form.title" placeholder="Judul konten" class="vs-input" />
                </div>

                <div class="mv-field-row">
                  <div class="mv-field">
                    <label class="mv-label">Tahun</label>
                    <input v-model="form.year" placeholder="2024" maxlength="4" class="vs-input" />
                  </div>
                  <div class="mv-field">
                    <label class="mv-label">Durasi</label>
                    <input v-model="form.duration" placeholder="120 min" class="vs-input" />
                  </div>
                  <div class="mv-field">
                    <label class="mv-label">Rating</label>
                    <input v-model="form.rating" placeholder="7.5" maxlength="4" class="vs-input" />
                  </div>
                </div>

                <div class="mv-field">
                  <label class="mv-label">Genre</label>
                  <input v-model="form.genre" placeholder="Action, Drama, Sci-Fi" class="vs-input" />
                  <p class="mv-field-hint">Pisahkan dengan koma</p>
                </div>

                <div class="mv-field">
                  <label class="mv-label">Overview / Sinopsis</label>
                  <textarea
                    v-model="form.overview"
                    placeholder="Deskripsi singkat..."
                    class="vs-input mv-textarea"
                    rows="4"
                  />
                </div>
              </template>

              <!-- ── Tab: Gambar ── -->
              <template v-if="editTab === 'media'">
                <div class="mv-field">
                  <label class="mv-label">URL Poster</label>
                  <input v-model="form.poster" placeholder="https://image.tmdb.org/t/p/w500/..." class="vs-input" />
                  <p class="mv-field-hint">Gunakan URL TMDB (w500) atau URL gambar eksternal</p>
                </div>

                <!-- Poster preview -->
                <div v-if="form.poster" class="mv-img-preview-wrap">
                  <div class="mv-img-preview-label">Preview Poster</div>
                  <div class="mv-poster-preview-box">
                    <img :src="form.poster" alt="Poster preview" class="mv-poster-preview-img"
                      @error="$event.target.style.display='none'" />
                  </div>
                </div>

                <div class="mv-field">
                  <label class="mv-label">URL Backdrop</label>
                  <input v-model="form.backdrop" placeholder="https://image.tmdb.org/t/p/w1280/..." class="vs-input" />
                  <p class="mv-field-hint">Gambar lebar untuk hero section (w1280 disarankan)</p>
                </div>

                <!-- Backdrop preview -->
                <div v-if="form.backdrop" class="mv-img-preview-wrap">
                  <div class="mv-img-preview-label">Preview Backdrop</div>
                  <div class="mv-backdrop-preview-box">
                    <img :src="form.backdrop" alt="Backdrop preview" class="mv-backdrop-preview-img"
                      @error="$event.target.style.display='none'" />
                  </div>
                </div>
              </template>

              <!-- ── Tab: Video Links ── -->
              <template v-if="editTab === 'links'">
                <!-- has_episodes toggle (anime only) -->
                <div v-if="form.type === 'anime'" class="mv-field">
                  <label class="mv-label">
                    <input type="checkbox" v-model="form.has_episodes" style="margin-right:6px" />
                    Anime dengan Episode/Season
                  </label>
                  <p class="mv-field-hint">Jika aktif, URL dikelola di level episode</p>
                </div>

                <template v-if="form.type !== 'anime' || !form.has_episodes">
                  <div class="mv-field">
                    <label class="mv-label">
                      Video URL — Server 1
                      <span class="mv-label-required">*</span>
                    </label>
                    <input v-model="form.url1" placeholder="https://..." class="vs-input" />
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
                </template>

                <div v-else class="mv-episodes-note">
                  <Tv :size="16" />
                  <p>URL untuk anime ini dikelola di level episode/season.<br>Klik tombol <strong>expand ▶</strong> di tabel untuk menambah season.</p>
                </div>
              </template>

            </template>
          </div>

          <!-- Footer -->
          <div class="mv-modal-footer">
            <button type="button" class="vs-btn vs-btn--ghost" @click="emit('close')">Batal</button>

            <!-- Edit mode: tombol berbeda per tab -->
            <template v-if="isEdit">
              <button
                v-if="editTab !== 'links' || (form.type !== 'anime' || !form.has_episodes)"
                type="button"
                class="vs-btn vs-btn--primary"
                :disabled="loading"
                @click="handleSubmit"
              >
                <Loader2 v-if="loading" :size="14" class="mv-spin" />
                {{ loading ? 'Menyimpan...' : 'Simpan Perubahan' }}
              </button>
              <!-- Navigasi antar tab dengan tombol Next -->
              <button
                v-if="editTab === 'meta'"
                type="button" class="vs-btn vs-btn--ghost mv-next-btn"
                @click="editTab = 'media'"
              >
                Gambar <ChevronRight :size="14" />
              </button>
              <button
                v-if="editTab === 'media'"
                type="button" class="vs-btn vs-btn--ghost mv-next-btn"
                @click="editTab = 'links'"
              >
                Video Links <ChevronRight :size="14" />
              </button>
            </template>

            <!-- Add mode: satu tombol submit -->
            <button
              v-else
              type="button"
              class="vs-btn vs-btn--primary"
              :disabled="loading"
              @click="handleSubmit"
            >
              <Loader2 v-if="loading" :size="14" class="mv-spin" />
              {{ loading ? 'Menyimpan...' : 'Tambah Konten' }}
            </button>
          </div>

        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.mv-backdrop {
  position: fixed; inset: 0; background: rgba(0,0,0,0.7);
  backdrop-filter: blur(6px); z-index: 100;
  display: flex; align-items: center; justify-content: center; padding: 16px;
}
.mv-modal {
  background: var(--surface); border: 1px solid var(--border);
  border-radius: 16px; width: 100%; max-width: 480px;
  box-shadow: 0 25px 60px rgba(0,0,0,0.5); overflow: hidden;
  max-height: 90vh; display: flex; flex-direction: column;
}
.mv-modal--wide { max-width: 560px; }

/* Header */
.mv-modal-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 16px 20px; border-bottom: 1px solid var(--border); flex-shrink: 0;
}
.mv-modal-title-wrap { display: flex; align-items: center; gap: 10px; }
.mv-modal-icon {
  width: 32px; height: 32px; border-radius: 8px;
  background: var(--accent-lo); color: var(--accent-hi); display: grid; place-items: center; flex-shrink: 0;
}
.mv-modal-title    { font-size: 14px; font-weight: 600; color: #fff; }
.mv-modal-subtitle { font-size: 11px; color: var(--muted); margin-top: 1px; }
.mv-modal-close {
  width: 28px; height: 28px; border-radius: 6px;
  border: 1px solid var(--border); background: transparent;
  color: var(--muted); display: grid; place-items: center; cursor: pointer; transition: all 0.15s;
}
.mv-modal-close:hover { color: #fff; background: var(--border2); }

/* Edit inner tabs */
.mv-edit-tabs {
  display: flex; gap: 2px; padding: 10px 16px 0;
  border-bottom: 1px solid var(--border); flex-shrink: 0; background: var(--surface);
}
.mv-edit-tab {
  display: flex; align-items: center; gap: 6px;
  padding: 7px 14px; border-radius: 8px 8px 0 0;
  font-size: 12px; font-weight: 600; border: none;
  background: transparent; color: var(--muted);
  cursor: pointer; transition: all 0.15s; font-family: var(--font);
  border-bottom: 2px solid transparent; margin-bottom: -1px;
}
.mv-edit-tab:hover { color: var(--text); }
.mv-edit-tab--active {
  color: var(--accent-hi);
  border-bottom-color: var(--accent-hi);
  background: var(--accent-lo);
}

/* Body */
.mv-modal-body { padding: 20px; display: flex; flex-direction: column; gap: 14px; overflow-y: auto; }

/* Footer */
.mv-modal-footer {
  display: flex; justify-content: flex-end; gap: 8px;
  padding: 14px 20px; border-top: 1px solid var(--border); flex-shrink: 0;
}
.mv-next-btn { gap: 4px; }

/* Fields */
.mv-field { display: flex; flex-direction: column; gap: 6px; }
.mv-field-row { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 10px; }
.mv-field-row:has(.mv-field:only-child) { grid-template-columns: 1fr 1fr; }
.mv-label          { font-size: 12px; font-weight: 500; color: var(--muted); }
.mv-label-required { color: var(--danger); }
.mv-field-hint     { font-size: 11px; color: var(--muted2); }
.mv-field-error    { font-size: 11px; color: var(--danger); margin-top: 2px; }
.vs-input--error   { border-color: rgba(239,68,68,0.6) !important; box-shadow: 0 0 0 3px rgba(239,68,68,0.08) !important; }

.mv-textarea {
  resize: vertical; min-height: 90px; line-height: 1.5;
  font-family: var(--font); font-size: 13px;
}

/* Type selector (add mode) */
.mv-type-selector { display: flex; gap: 6px; }
.mv-type-opt {
  flex: 1; display: flex; align-items: center; justify-content: center; gap: 6px;
  padding: 8px 10px; border-radius: 8px; border: 1px solid var(--border);
  background: var(--surface2); color: var(--muted);
  font-size: 12.5px; font-weight: 600; cursor: pointer; transition: all 0.15s; font-family: var(--font);
}
.mv-type-opt:hover { color: var(--text); }
.mv-type-opt--active {
  background: color-mix(in srgb, var(--opt-color) 12%, transparent);
  color: var(--opt-color); border-color: color-mix(in srgb, var(--opt-color) 35%, transparent);
}

/* Refresh banner */
.mv-refresh-banner {
  display: flex; align-items: center; justify-content: space-between; gap: 10px;
  background: rgba(99,102,241,0.08); border: 1px solid rgba(99,102,241,0.2);
  border-radius: var(--radius-sm); padding: 10px 12px;
}
.mv-refresh-info {
  display: flex; align-items: center; gap: 8px;
  font-size: 11.5px; color: var(--muted);
}
.mv-refresh-btn { padding: 6px 12px !important; font-size: 12px !important; }

/* Preview (add mode) */
.mv-preview-box {
  display: flex; gap: 12px; background: var(--surface2);
  border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 12px;
}
.mv-preview-poster { width: 50px; height: 75px; border-radius: 4px; object-fit: cover; flex-shrink: 0; }
.mv-preview-info   { display: flex; flex-direction: column; gap: 4px; flex: 1; min-width: 0; }
.mv-preview-title  { font-size: 13px; font-weight: 600; color: #fff; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.mv-preview-meta   { display: flex; align-items: center; gap: 6px; font-size: 11px; color: var(--muted); }
.mv-preview-desc   {
  font-size: 11px; color: var(--muted2); line-height: 1.4;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
}

/* Image previews (edit → gambar tab) */
.mv-img-preview-wrap  { display: flex; flex-direction: column; gap: 6px; }
.mv-img-preview-label { font-size: 11px; color: var(--muted2); }
.mv-poster-preview-box {
  width: 80px; height: 120px; border-radius: 8px; overflow: hidden;
  background: var(--surface2); border: 1px solid var(--border);
}
.mv-poster-preview-img   { width: 100%; height: 100%; object-fit: cover; }
.mv-backdrop-preview-box {
  width: 100%; height: 120px; border-radius: 8px; overflow: hidden;
  background: var(--surface2); border: 1px solid var(--border);
}
.mv-backdrop-preview-img { width: 100%; height: 100%; object-fit: cover; }

/* Episodes note */
.mv-episodes-note {
  display: flex; align-items: flex-start; gap: 10px;
  background: var(--surface2); border: 1px solid var(--border);
  border-radius: var(--radius-sm); padding: 14px;
  font-size: 12.5px; color: var(--muted); line-height: 1.6;
}
.mv-episodes-note svg { flex-shrink: 0; margin-top: 1px; color: #22d3ee; }
.mv-episodes-note strong { color: var(--text); }

/* Animations */
.mv-spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.mv-modal-enter-active { animation: modal-in 0.22s cubic-bezier(0.34,1.56,0.64,1); }
.mv-modal-leave-active { animation: modal-out 0.18s ease; }
@keyframes modal-in  { from { opacity:0; transform:scale(0.94) translateY(10px); } }
@keyframes modal-out { to   { opacity:0; transform:scale(0.96) translateY(6px); } }
</style>