<script setup>
import { ref, watch } from 'vue'
import { Film, Search, Loader2, Star, Calendar, Tv, Sparkles } from 'lucide-vue-next'

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

const form = ref({ tmdb_id: '', type: 'movie', has_episodes: false, url1: '', url2: '', url3: '' })
const previewData     = ref(null)
const isSearchingTMDB = ref(false)
const tmdbIdError     = ref('')

// Sync type ke activeTab saat modal dibuka
watch(() => props.show, (val) => {
  if (val && !props.isEdit) form.value.type = props.activeTab
})

const validateTmdbId = () => {
  const val = form.value.tmdb_id
  if (!val)              { tmdbIdError.value = ''; return }
  if (!/^\d+$/.test(val))  { tmdbIdError.value = 'TMDB ID hanya boleh berisi angka'; return }
  if (val.length > 10)     { tmdbIdError.value = 'TMDB ID maksimal 10 digit'; return }
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

const handleSubmit = () => {
  validateTmdbId()
  if (tmdbIdError.value) return
  emit('submit', { ...form.value })
}

// Expose untuk parent bisa set form (saat edit)
const setForm = (data) => { form.value = { ...data } }
const setPreview = (data) => { previewData.value = data }
const reset = () => {
  form.value = { tmdb_id: '', type: props.activeTab, has_episodes: false, url1: '', url2: '', url3: '' }
  previewData.value = null
  tmdbIdError.value = ''
}

defineExpose({ setForm, setPreview, reset })
</script>

<template>
  <Teleport to="body">
    <Transition name="mv-modal">
      <div v-if="show" class="mv-backdrop" @click.self="emit('close')">
        <div class="mv-modal">
          <div class="mv-modal-header">
            <div class="mv-modal-title-wrap">
              <div class="mv-modal-icon"><Film :size="16" /></div>
              <h2 class="mv-modal-title">{{ isEdit ? 'Edit Konten' : 'Tambah Konten Baru' }}</h2>
            </div>
            <button class="mv-modal-close" @click="emit('close')">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <div class="mv-modal-body">
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
                  <Film     v-if="tab.id === 'movie'"  :size="13" />
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
                    required
                    inputmode="numeric"
                    pattern="[0-9]{1,10}"
                    maxlength="10"
                    class="vs-input"
                    :class="{ 'vs-input--error': tmdbIdError }"
                    :disabled="isEdit"
                    @input="validateTmdbId"
                  />
                  <p v-if="tmdbIdError" class="mv-field-error">{{ tmdbIdError }}</p>
                </div>
                <button
                  type="button"
                  class="vs-btn vs-btn--ghost"
                  @click="searchTMDB"
                  :disabled="!form.tmdb_id || !!tmdbIdError || isSearchingTMDB || isEdit"
                  style="flex-shrink:0; align-self:flex-start"
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
                class="mv-preview-poster" alt="Poster"
              />
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
              <input
                v-model="form.url1"
                placeholder="https://..."
                :required="form.type === 'movie' || (form.type === 'anime' && !form.has_episodes)"
                class="vs-input"
              />
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
              <button type="button" class="vs-btn vs-btn--ghost" @click="emit('close')">Batal</button>
              <button type="button" class="vs-btn vs-btn--primary" :disabled="loading" @click="handleSubmit">
                <Loader2 v-if="loading" :size="14" class="mv-spin" />
                {{ loading ? 'Menyimpan...' : (isEdit ? 'Simpan Perubahan' : 'Tambah Konten') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.mv-backdrop {
  position: fixed; inset: 0; background: rgba(0,0,0,0.65);
  backdrop-filter: blur(6px); z-index: 100;
  display: flex; align-items: center; justify-content: center; padding: 16px;
}
.mv-modal {
  background: var(--surface); border: 1px solid var(--border);
  border-radius: 16px; width: 100%; max-width: 460px;
  box-shadow: 0 25px 60px rgba(0,0,0,0.5); overflow: hidden;
  max-height: 90vh; display: flex; flex-direction: column;
}
.mv-modal-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 18px 20px; border-bottom: 1px solid var(--border); flex-shrink: 0;
}
.mv-modal-title-wrap { display: flex; align-items: center; gap: 10px; }
.mv-modal-icon {
  width: 32px; height: 32px; border-radius: 8px;
  background: var(--accent-lo); color: var(--accent-hi); display: grid; place-items: center;
}
.mv-modal-title { font-size: 14px; font-weight: 600; color: #fff; }
.mv-modal-close {
  width: 28px; height: 28px; border-radius: 6px;
  border: 1px solid var(--border); background: transparent;
  color: var(--muted); display: grid; place-items: center; cursor: pointer; transition: all 0.15s;
}
.mv-modal-close:hover { color: #fff; background: var(--border2); }
.mv-modal-body { padding: 20px; display: flex; flex-direction: column; gap: 16px; overflow-y: auto; }
.mv-field { display: flex; flex-direction: column; gap: 6px; }
.mv-field-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.mv-label { font-size: 12px; font-weight: 500; color: var(--muted); }
.mv-label-required { color: var(--danger); }
.mv-field-hint  { font-size: 11px; color: var(--muted2); }
.mv-field-error { font-size: 11px; color: var(--danger); margin-top: 2px; }
.mv-modal-footer { display: flex; justify-content: flex-end; gap: 8px; padding-top: 16px; border-top: 1px solid var(--border); }

.mv-type-selector { display: flex; gap: 6px; }
.mv-type-opt {
  flex: 1; display: flex; align-items: center; justify-content: center; gap: 6px;
  padding: 8px 12px; border-radius: 8px; border: 1px solid var(--border);
  background: var(--surface2); color: var(--muted);
  font-size: 12.5px; font-weight: 600; cursor: pointer; transition: all 0.15s; font-family: var(--font);
}
.mv-type-opt:hover { color: var(--text); }
.mv-type-opt--active {
  background: color-mix(in srgb, var(--opt-color) 12%, transparent);
  color: var(--opt-color); border-color: color-mix(in srgb, var(--opt-color) 35%, transparent);
}
.vs-input--error { border-color: rgba(239,68,68,0.6) !important; box-shadow: 0 0 0 3px rgba(239,68,68,0.08) !important; }

.mv-preview-box { display: flex; gap: 12px; background: var(--surface2); border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 12px; }
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