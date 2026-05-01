<script setup>
import { ref, onMounted } from 'vue'
import api from '@/lib/api'
import { Plus, Trash2, GripVertical, Loader2, ImagePlay } from 'lucide-vue-next'

const slides = ref([])
const allContent = ref([])
const loading = ref(false)
const showModal = ref(false)
const selectedMovieId = ref(0)
const selectedOrder = ref(1)

const fetchSlides = async () => {
  loading.value = true
  try {
    const res = await api.get('/admin/hero-slides')
    slides.value = res.data.data ?? []
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const fetchAllContent = async () => {
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
  } catch (err) {
    console.error(err)
  }
}

const addSlide = async () => {
  console.log('addSlide called, selectedMovieId:', selectedMovieId.value)
  if (!selectedMovieId.value) {
    console.log('BLOCKED: selectedMovieId is falsy')
    return
  }
  try {
    const payload = {
      movie_id: selectedMovieId.value,
      order: selectedOrder.value,
    }
    console.log('Posting payload:', payload)
    const res = await api.post('/admin/hero-slides', payload)
    console.log('Response:', res.data)
    showModal.value = false
    await fetchSlides()
  } catch (err) {
    console.error('Full error:', err)
    console.error('Response data:', err.response?.data)
    alert(err.response?.data?.error || 'Gagal menambah slide')
  }
}

const onSelectChange = (event) => {
  const idx = event.target.selectedIndex - 1 // minus 1 karena ada option disabled di index 0
  const content = availableContent()
  if (idx >= 0 && content[idx]) {
    selectedMovieId.value = content[idx].id
    console.log('Selected id:', selectedMovieId.value)
  }
}
const deleteSlide = async (id) => {
  if (!confirm('Hapus slide ini?')) return
  try {
    await api.delete(`/admin/hero-slides/${id}`)
    await fetchSlides()
  } catch {
    alert('Gagal menghapus slide')
  }
}

const updateOrder = async (slide, newOrder) => {
  try {
    await api.put(`/admin/hero-slides/${slide.id}`, { order: Number(newOrder) })
    await fetchSlides()
  } catch {
    alert('Gagal update urutan')
  }
}

const availableContent = () => {
  const usedIds = slides.value.map(s => Number(s.movie?.id)).filter(Boolean)
  return allContent.value.filter(c => !usedIds.includes(Number(c.id)))
}

const selectedPreview = () => {
  return allContent.value.find(c => c.id === selectedMovieId.value) ?? null
}

const openAddModal = () => {
  selectedMovieId.value = 0
  selectedOrder.value = slides.value.length + 1
  showModal.value = true
}

onMounted(() => {
  fetchSlides()
  fetchAllContent()
})
</script>

<template>
  <div>
    <!-- Header -->
    <div class="vs-page-header">
      <div>
        <h1 class="vs-page-title">
          <ImagePlay :size="22" />
          Hero Slider
        </h1>
        <p class="vs-page-subtitle">Kelola konten yang tampil di hero section — maksimal 5 slide</p>
      </div>
      <button
        class="vs-btn vs-btn--primary"
        type="button"
        :disabled="slides.length >= 5"
        @click="openAddModal"
        >
        <Plus :size="15" />
        Tambah Slide
      </button>
    </div>

    <!-- Card -->
    <div class="vs-card">

      <!-- Loading -->
      <div v-if="loading" class="hs-loading">
        <Loader2 :size="20" class="hs-spin" />
        <p>Memuat slides...</p>
      </div>

      <!-- Empty -->
      <div v-else-if="slides.length === 0" class="hs-empty">
        <ImagePlay :size="36" />
        <p>Belum ada slide hero</p>
        <p class="hs-empty-sub">Klik "Tambah Slide" untuk menambahkan konten ke hero section</p>
      </div>

      <!-- Slide list -->
      <div v-else class="hs-list">
        <div
          v-for="slide in slides"
          :key="slide.id"
          class="hs-item"
        >
          <!-- Drag handle + order -->
          <div class="hs-order">
            <GripVertical :size="14" class="hs-grip" />
            <input
              type="number"
              class="hs-order-input vs-input"
              :value="slide.order"
              min="1"
              max="5"
              @change="updateOrder(slide, $event.target.value)"
            />
          </div>

          <!-- Backdrop preview -->
          <div class="hs-backdrop-wrap">
            <img
              v-if="slide.movie?.backdrop"
              :src="slide.movie.backdrop"
              class="hs-backdrop"
              :alt="slide.movie?.title"
            />
            <div v-else class="hs-backdrop-empty" />
          </div>

          <!-- Poster -->
          <img
            v-if="slide.movie?.poster"
            :src="slide.movie.poster"
            class="hs-poster"
            :alt="slide.movie?.title"
          />

          <!-- Info -->
          <div class="hs-info">
            <p class="hs-title">{{ slide.movie?.title }}</p>
            <div class="hs-meta">
              <span
                class="hs-type-badge"
                :class="`hs-type-badge--${slide.movie?.type}`"
              >
                {{ slide.movie?.type === 'movie' ? 'Movie' : slide.movie?.type === 'series' ? 'Series' : 'Anime' }}
              </span>
              <span class="hs-year">{{ slide.movie?.year }}</span>
              <span class="hs-rating">★ {{ slide.movie?.rating }}</span>
            </div>
            <p class="hs-genre">{{ slide.movie?.genre }}</p>
          </div>

          <!-- Delete -->
          <button class="mv-action-btn mv-action-btn--delete hs-del-btn" @click="deleteSlide(slide.id)">
            <Trash2 :size="14" />
          </button>
        </div>
      </div>

      <!-- Footer -->
      <div class="hs-footer">
        <span class="hs-count">{{ slides.length }} / 5 slide terpakai</span>
        <div class="hs-dots-preview">
          <span
            v-for="i in 5"
            :key="i"
            class="hs-dot-preview"
            :class="{ 'hs-dot-preview--active': i <= slides.length }"
          />
        </div>
      </div>
    </div>

    <!-- Modal tambah slide -->
    <Teleport to="body">
      <Transition name="mv-modal">
        <div v-if="showModal" class="mv-backdrop" @click.self="showModal = false">
          <div class="mv-modal">
            <div class="mv-modal-header">
              <div class="mv-modal-title-wrap">
                <div class="mv-modal-icon">
                  <ImagePlay :size="16" />
                </div>
                <h2 class="mv-modal-title">Tambah Slide Hero</h2>
              </div>
              <button class="mv-modal-close" @click="showModal = false">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 6L6 18M6 6l12 12"/>
                </svg>
              </button>
            </div>

            <div class="mv-modal-body">
              <!-- Pilih konten -->
              <div class="mv-field">
                <label class="mv-label">Pilih Konten <span class="mv-label-required">*</span></label>
                <select
                  class="vs-input hs-select"
                  @change="onSelectChange"
                >
                  <option disabled selected>— Pilih film / series / anime —</option>
                  <option
                    v-for="c in availableContent()"
                    :key="c.id"
                  >
                    [{{ c.type.toUpperCase() }}] {{ c.title }} ({{ c.year }})
                  </option>
                </select>
                <p v-if="availableContent().length === 0" class="mv-field-hint" style="color:var(--warning)">
                  Semua konten sudah dipakai sebagai slide.
                </p>
              </div>

              <!-- Preview konten terpilih -->
              <div v-if="selectedMovieId && selectedPreview()" class="mv-preview-box">
                <img :src="selectedPreview().poster" class="mv-preview-poster" :alt="selectedPreview().title" />
                <div class="mv-preview-info">
                  <h4 class="mv-preview-title">{{ selectedPreview().title }}</h4>
                  <div class="mv-preview-meta">
                    ★ {{ selectedPreview().rating }} · {{ selectedPreview().year }} · {{ selectedPreview().genre }}
                  </div>
                </div>
              </div>

              <!-- Order -->
              <div class="mv-field">
                <label class="mv-label">Urutan Tampil</label>
                <input
                  v-model.number="selectedOrder"
                  type="number"
                  min="1"
                  max="5"
                  class="vs-input"
                  style="max-width: 100px"
                />
              </div>

              <div class="mv-modal-footer">
                <button class="vs-btn vs-btn--ghost" @click="showModal = false">Batal</button>
                <button
                  class="vs-btn vs-btn--primary"
                  :disabled="!selectedMovieId"
                  @click="addSlide"
                >
                  Tambah Slide
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
.hs-loading {
  display: flex; flex-direction: column; align-items: center;
  gap: 10px; padding: 48px; color: var(--muted); font-size: 13px;
}
.hs-spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.hs-empty {
  display: flex; flex-direction: column; align-items: center;
  gap: 8px; padding: 60px 20px;
  color: var(--muted2); font-size: 14px; font-weight: 500; text-align: center;
}
.hs-empty svg { opacity: 0.4; }
.hs-empty-sub { font-size: 12px; opacity: 0.6; font-weight: 400; margin-top: 2px; }

/* Slide list */
.hs-list { display: flex; flex-direction: column; }

.hs-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 18px;
  border-bottom: 1px solid var(--border2);
  transition: background 0.15s;
}
.hs-item:last-child { border-bottom: none; }
.hs-item:hover { background: rgba(255,255,255,0.02); }

.hs-order {
  display: flex; align-items: center; gap: 6px; flex-shrink: 0;
}
.hs-grip { color: var(--muted2); cursor: grab; }
.hs-order-input {
  width: 48px !important;
  padding: 5px 8px !important;
  text-align: center;
  font-size: 13px !important;
}

/* Backdrop thumbnail */
.hs-backdrop-wrap {
  width: 120px; height: 68px;
  border-radius: 6px; overflow: hidden;
  flex-shrink: 0; background: var(--surface2);
}
.hs-backdrop {
  width: 100%; height: 100%;
  object-fit: cover;
}
.hs-backdrop-empty {
  width: 100%; height: 100%;
  background: var(--surface2);
}

/* Poster */
.hs-poster {
  width: 36px; height: 52px;
  border-radius: 5px; object-fit: cover;
  flex-shrink: 0;
}

/* Info */
.hs-info { flex: 1; min-width: 0; }
.hs-title { font-size: 13.5px; font-weight: 600; color: #fff; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.hs-meta { display: flex; align-items: center; gap: 8px; margin-top: 4px; flex-wrap: wrap; }
.hs-type-badge {
  font-size: 10px; font-weight: 700;
  padding: 2px 7px; border-radius: 4px;
  text-transform: uppercase; letter-spacing: 0.06em; color: #fff;
}
.hs-type-badge--movie  { background: rgba(99,102,241,0.85); }
.hs-type-badge--series { background: rgba(6,182,212,0.85); }
.hs-type-badge--anime  { background: rgba(244,63,94,0.85); }
.hs-year, .hs-rating { font-size: 11.5px; color: var(--muted); }
.hs-rating { color: #f59e0b; }
.hs-genre { font-size: 11px; color: var(--muted2); margin-top: 3px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.hs-del-btn { flex-shrink: 0; }

/* Footer */
.hs-footer {
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 18px;
  border-top: 1px solid var(--border);
  background: var(--surface2);
}
.hs-count { font-size: 12px; color: var(--muted); }
.hs-dots-preview { display: flex; gap: 5px; align-items: center; }
.hs-dot-preview {
  width: 6px; height: 6px; border-radius: 50%;
  background: var(--border);
  transition: all 0.2s;
}
.hs-dot-preview--active {
  background: var(--accent-hi);
  width: 18px; border-radius: 3px;
}

/* Select */
.hs-select { cursor: pointer; }
option { background: #14141f; }

/* Modal preview */
.hs-modal-preview { margin: -4px 0; }

/* Shared modal styles dari movies.vue */
.mv-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.65);
  backdrop-filter: blur(6px);
  z-index: 100;
  display: flex; align-items: center; justify-content: center;
  padding: 16px;
}
.mv-modal {
  background: var(--surface); border: 1px solid var(--border);
  border-radius: 16px; width: 100%; max-width: 460px;
  box-shadow: 0 25px 60px rgba(0,0,0,0.5);
  overflow: hidden; max-height: 90vh;
  display: flex; flex-direction: column;
}
.mv-modal-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 18px 20px; border-bottom: 1px solid var(--border); flex-shrink: 0;
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
.mv-modal-close:hover { color: #fff; background: var(--border2); }
.mv-modal-body {
  padding: 20px; display: flex; flex-direction: column; gap: 16px; overflow-y: auto;
}
.mv-field { display: flex; flex-direction: column; gap: 6px; }
.mv-label { font-size: 12px; font-weight: 500; color: var(--muted); }
.mv-label-required { color: var(--danger); }
.mv-field-hint { font-size: 11px; color: var(--muted2); }
.mv-modal-footer {
  display: flex; justify-content: flex-end; gap: 8px;
  padding-top: 16px; border-top: 1px solid var(--border);
}
.mv-preview-box {
  display: flex; gap: 12px;
  background: var(--surface2); border: 1px solid var(--border);
  border-radius: var(--radius-sm); padding: 12px;
}
.mv-preview-poster { width: 50px; height: 75px; border-radius: 4px; object-fit: cover; flex-shrink: 0; }
.mv-preview-info { display: flex; flex-direction: column; gap: 4px; flex: 1; min-width: 0; }
.mv-preview-title { font-size: 13px; font-weight: 600; color: #fff; }
.mv-preview-meta { font-size: 11px; color: var(--muted); }

/* Action button (reuse dari movies) */
.mv-action-btn {
  width: 30px; height: 30px; border-radius: 7px;
  border: 1px solid var(--border); background: transparent;
  cursor: pointer; display: grid; place-items: center;
  transition: all 0.15s; color: var(--muted);
}
.mv-action-btn--delete:hover { color: var(--danger); border-color: rgba(239,68,68,0.3); background: rgba(239,68,68,0.08); }

.mv-modal-enter-active { animation: modal-in 0.22s cubic-bezier(0.34,1.56,0.64,1); }
.mv-modal-leave-active { animation: modal-out 0.18s ease; }
@keyframes modal-in  { from { opacity:0; transform:scale(0.94) translateY(10px); } }
@keyframes modal-out { to   { opacity:0; transform:scale(0.96) translateY(6px); } }
</style>