<script setup>
import { ref, watch } from 'vue'
import { Captions, Upload, Trash2, Loader2, CheckCircle2 } from 'lucide-vue-next'
import api from '@/lib/api'

const props = defineProps({
  show: { type: Boolean, required: true },
})

const emit = defineEmits(['close'])

// ── State ─────────────────────────────────────────────────────
const movieId    = ref(null)
const movieTitle = ref('')
const existingSubs = ref([])  // subtitle yang sudah ada
const loading    = ref(false)
const loadingSubs = ref(false)
const deletingId = ref(null)
const successMsg = ref('')

const form = ref({
  lang:  'id',
  label: 'Indonesia',
  file:  null,
})
const fileInputRef = ref(null)
const dragover     = ref(false)
const fileError    = ref('')

const LANG_OPTIONS = [
  { lang: 'id', label: 'Indonesia' },
  { lang: 'en', label: 'English'   },
]

// ── Watchers ──────────────────────────────────────────────────
watch(() => form.value.lang, (lang) => {
  form.value.label = LANG_OPTIONS.find(o => o.lang === lang)?.label ?? lang
})

// ── Expose untuk parent ───────────────────────────────────────
const open = async (id, title) => {
  movieId.value    = id
  movieTitle.value = title
  form.value = { lang: 'id', label: 'Indonesia', file: null }
  fileError.value  = ''
  successMsg.value = ''
  if (fileInputRef.value) fileInputRef.value.value = ''
  await fetchExisting()
}

const fetchExisting = async () => {
  if (!movieId.value) return
  loadingSubs.value = true
  try {
    const res = await api.get(`/admin/content/${movieId.value}/subtitles`)
    existingSubs.value = res.data.data ?? []
  } catch {
    existingSubs.value = []
  } finally {
    loadingSubs.value = false
  }
}

defineExpose({ open })

// ── File handling ─────────────────────────────────────────────
const validateFile = (file) => {
  if (!file) return 'File wajib dipilih'
  const ext = file.name.split('.').pop().toLowerCase()
  if (!['vtt', 'srt'].includes(ext)) return 'Hanya file .vtt atau .srt yang diizinkan'
  if (file.size > 5 * 1024 * 1024) return 'Ukuran file maksimal 5MB'
  return ''
}

const onFileChange = (e) => {
  const file = e.target.files?.[0] ?? null
  fileError.value = validateFile(file)
  form.value.file = fileError.value ? null : file
}

const onDrop = (e) => {
  dragover.value = false
  const file = e.dataTransfer.files?.[0] ?? null
  fileError.value = validateFile(file)
  form.value.file = fileError.value ? null : file
  // Sync ke input element juga (opsional, untuk tampilan nama file)
  if (!fileError.value && fileInputRef.value) {
    const dt = new DataTransfer()
    dt.items.add(file)
    fileInputRef.value.files = dt.files
  }
}

// ── Submit ────────────────────────────────────────────────────
const handleSubmit = async () => {
  fileError.value = validateFile(form.value.file)
  if (fileError.value) return

  loading.value = true
  successMsg.value = ''
  try {
    const fd = new FormData()
    fd.append('lang',  form.value.lang)
    fd.append('label', form.value.label)
    fd.append('file',  form.value.file)

    await api.post(`/admin/content/${movieId.value}/subtitles`, fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })

    successMsg.value = `Subtitle ${form.value.label} berhasil diupload`
    form.value.file = null
    if (fileInputRef.value) fileInputRef.value.value = ''
    await fetchExisting()
  } catch (err) {
    fileError.value = err.response?.data?.error ?? 'Gagal mengupload subtitle'
  } finally {
    loading.value = false
  }
}

// ── Delete ────────────────────────────────────────────────────
const handleDelete = async (sub) => {
  if (!confirm(`Hapus subtitle ${sub.label}?`)) return
  deletingId.value = sub.id
  try {
    await api.delete(`/admin/subtitles/${sub.id}`)
    await fetchExisting()
  } catch {
    alert('Gagal menghapus subtitle')
  } finally {
    deletingId.value = null
  }
}

// ── Close ─────────────────────────────────────────────────────
const handleClose = () => {
  movieId.value    = null
  movieTitle.value = ''
  existingSubs.value = []
  successMsg.value = ''
  fileError.value  = ''
  emit('close')
}
</script>

<template>
  <Teleport to="body">
    <Transition name="mv-modal">
      <div v-if="show" class="mv-backdrop" @click.self="handleClose">
        <div class="mv-modal">

          <!-- Header -->
          <div class="mv-modal-header">
            <div class="mv-modal-title-wrap">
              <div class="mv-modal-icon">
                <Captions :size="16" />
              </div>
              <div>
                <h2 class="mv-modal-title">Kelola Subtitle</h2>
                <p class="sub-movie-title">{{ movieTitle }}</p>
              </div>
            </div>
            <button class="mv-modal-close" @click="handleClose">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <div class="mv-modal-body">

            <!-- Existing subtitles -->
            <div class="sub-section">
              <p class="sub-section-label">Subtitle Tersedia</p>

              <div v-if="loadingSubs" class="sub-loading">
                <Loader2 :size="16" class="mv-spin" />
                <span>Memuat...</span>
              </div>

              <div v-else-if="existingSubs.length === 0" class="sub-empty">
                Belum ada subtitle untuk konten ini
              </div>

              <div v-else class="sub-list">
                <div v-for="sub in existingSubs" :key="sub.id" class="sub-item">
                  <div class="sub-item-left">
                    <span class="sub-lang-badge">{{ sub.lang.toUpperCase() }}</span>
                    <div>
                      <p class="sub-item-label">{{ sub.label }}</p>
                      <p class="sub-item-file">{{ sub.filename }}</p>
                    </div>
                  </div>
                  <button
                    class="mv-action-btn mv-action-btn--delete"
                    @click="handleDelete(sub)"
                    :disabled="deletingId === sub.id"
                    title="Hapus subtitle"
                  >
                    <Loader2 v-if="deletingId === sub.id" :size="13" class="mv-spin" />
                    <Trash2  v-else :size="13" />
                  </button>
                </div>
              </div>
            </div>

            <div class="sub-divider" />

            <!-- Upload form -->
            <div class="sub-section">
              <p class="sub-section-label">Upload Subtitle Baru</p>

              <!-- Success message -->
              <Transition name="fade">
                <div v-if="successMsg" class="sub-success">
                  <CheckCircle2 :size="14" />
                  {{ successMsg }}
                </div>
              </Transition>

              <!-- Lang selector -->
              <div class="mv-field">
                <label class="mv-label">Bahasa <span class="mv-label-required">*</span></label>
                <div class="sub-lang-selector">
                  <button
                    v-for="opt in LANG_OPTIONS"
                    :key="opt.lang"
                    type="button"
                    class="sub-lang-opt"
                    :class="{ 'sub-lang-opt--active': form.lang === opt.lang }"
                    @click="form.lang = opt.lang"
                  >
                    {{ opt.label }}
                    <span class="sub-lang-code">{{ opt.lang }}</span>
                  </button>
                </div>
              </div>

              <!-- Drop zone -->
              <div class="mv-field">
                <label class="mv-label">File Subtitle <span class="mv-label-required">*</span></label>
                <div
                  class="sub-dropzone"
                  :class="{
                    'sub-dropzone--dragover': dragover,
                    'sub-dropzone--error':    !!fileError,
                    'sub-dropzone--ok':       form.file && !fileError,
                  }"
                  @dragover.prevent="dragover = true"
                  @dragleave="dragover = false"
                  @drop.prevent="onDrop"
                  @click="fileInputRef?.click()"
                >
                  <input
                    ref="fileInputRef"
                    type="file"
                    accept=".vtt,.srt"
                    class="sub-file-input"
                    @change="onFileChange"
                  />

                  <div v-if="form.file && !fileError" class="sub-dropzone-content">
                    <CheckCircle2 :size="22" class="sub-icon-ok" />
                    <p class="sub-dropzone-filename">{{ form.file.name }}</p>
                    <p class="sub-dropzone-hint">Klik untuk ganti file</p>
                  </div>
                  <div v-else class="sub-dropzone-content">
                    <Upload :size="22" class="sub-icon-upload" />
                    <p class="sub-dropzone-text">
                      Drag & drop atau <span class="sub-dropzone-link">pilih file</span>
                    </p>
                    <p class="sub-dropzone-hint">.vtt atau .srt · maks 5MB</p>
                  </div>
                </div>
                <p v-if="fileError" class="mv-field-error">{{ fileError }}</p>
              </div>
            </div>

            <!-- Footer -->
            <div class="mv-modal-footer">
              <button type="button" class="vs-btn vs-btn--ghost" @click="handleClose">Tutup</button>
              <button
                type="button"
                class="vs-btn vs-btn--primary"
                :disabled="!form.file || !!fileError || loading"
                @click="handleSubmit"
              >
                <Loader2 v-if="loading" :size="14" class="mv-spin" />
                {{ loading ? 'Mengupload...' : 'Upload Subtitle' }}
              </button>
            </div>

          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
/* ── Modal shell (sama dengan modal lain) ─────────────────── */
.mv-backdrop {
  position: fixed; inset: 0; background: rgba(0,0,0,0.65);
  backdrop-filter: blur(6px); z-index: 100;
  display: flex; align-items: center; justify-content: center; padding: 16px;
}
.mv-modal {
  background: var(--surface); border: 1px solid var(--border);
  border-radius: 16px; width: 100%; max-width: 460px;
  box-shadow: 0 25px 60px rgba(0,0,0,0.5);
  overflow: hidden; max-height: 90vh; display: flex; flex-direction: column;
}
.mv-modal-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 18px 20px; border-bottom: 1px solid var(--border); flex-shrink: 0;
}
.mv-modal-title-wrap { display: flex; align-items: center; gap: 10px; }
.mv-modal-icon {
  width: 32px; height: 32px; border-radius: 8px;
  background: rgba(6,182,212,0.12); color: #22d3ee;
  display: grid; place-items: center; flex-shrink: 0;
}
.mv-modal-title  { font-size: 14px; font-weight: 600; color: #fff; line-height: 1.2; }
.sub-movie-title { font-size: 11px; color: var(--muted); margin-top: 2px; }
.mv-modal-close {
  width: 28px; height: 28px; border-radius: 6px;
  border: 1px solid var(--border); background: transparent;
  color: var(--muted); display: grid; place-items: center;
  cursor: pointer; transition: all 0.15s; flex-shrink: 0;
}
.mv-modal-close:hover { color: #fff; background: var(--border2); }
.mv-modal-body {
  padding: 20px; display: flex; flex-direction: column;
  gap: 16px; overflow-y: auto;
}
.mv-field       { display: flex; flex-direction: column; gap: 6px; }
.mv-label       { font-size: 12px; font-weight: 500; color: var(--muted); }
.mv-label-required { color: var(--danger); }
.mv-field-error { font-size: 11px; color: var(--danger); margin-top: 2px; }
.mv-modal-footer {
  display: flex; justify-content: flex-end; gap: 8px;
  padding-top: 16px; border-top: 1px solid var(--border);
}
.mv-action-btn {
  width: 30px; height: 30px; border-radius: 7px;
  border: 1px solid var(--border); background: transparent;
  cursor: pointer; display: grid; place-items: center;
  transition: all 0.15s; color: var(--muted); flex-shrink: 0;
}
.mv-action-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.mv-action-btn--delete:hover { color: var(--danger); border-color: rgba(239,68,68,0.3); background: rgba(239,68,68,0.08); }

/* ── Sections ─────────────────────────────────────────────── */
.sub-section       { display: flex; flex-direction: column; gap: 10px; }
.sub-section-label { font-size: 10.5px; font-weight: 600; letter-spacing: 0.07em; text-transform: uppercase; color: var(--muted2); }
.sub-divider       { height: 1px; background: var(--border); }

/* ── Existing list ────────────────────────────────────────── */
.sub-loading { display: flex; align-items: center; gap: 8px; font-size: 12px; color: var(--muted); padding: 8px 0; }
.sub-empty   { font-size: 12.5px; color: var(--muted2); padding: 10px 0; }
.sub-list    { display: flex; flex-direction: column; gap: 6px; }
.sub-item    {
  display: flex; align-items: center; justify-content: space-between;
  background: var(--surface2); border: 1px solid var(--border);
  border-radius: 8px; padding: 10px 12px; gap: 10px;
}
.sub-item-left  { display: flex; align-items: center; gap: 10px; min-width: 0; }
.sub-lang-badge {
  font-size: 10px; font-weight: 700; letter-spacing: 0.08em;
  color: #22d3ee; background: rgba(6,182,212,0.12);
  border: 1px solid rgba(6,182,212,0.25);
  padding: 2px 7px; border-radius: 4px; flex-shrink: 0;
}
.sub-item-label { font-size: 12.5px; font-weight: 500; color: #fff; }
.sub-item-file  { font-size: 10.5px; color: var(--muted2); margin-top: 1px; }

/* ── Success ──────────────────────────────────────────────── */
.sub-success {
  display: flex; align-items: center; gap: 8px;
  padding: 9px 12px; border-radius: 8px;
  background: rgba(16,185,129,0.08); border: 1px solid rgba(16,185,129,0.2);
  color: #34d399; font-size: 12.5px;
}

/* ── Lang selector ────────────────────────────────────────── */
.sub-lang-selector { display: flex; gap: 6px; }
.sub-lang-opt {
  flex: 1; display: flex; align-items: center; justify-content: center; gap: 6px;
  padding: 8px 12px; border-radius: 8px;
  border: 1px solid var(--border); background: var(--surface2);
  color: var(--muted); font-size: 12.5px; font-weight: 600;
  cursor: pointer; transition: all 0.15s; font-family: var(--font);
}
.sub-lang-opt:hover { color: var(--text); }
.sub-lang-opt--active {
  background: rgba(6,182,212,0.1); color: #22d3ee;
  border-color: rgba(6,182,212,0.35);
}
.sub-lang-code {
  font-size: 9px; font-weight: 700; letter-spacing: 0.07em;
  color: inherit; opacity: 0.6; text-transform: uppercase;
}

/* ── Drop zone ────────────────────────────────────────────── */
.sub-dropzone {
  border: 1.5px dashed rgba(255,255,255,0.1);
  border-radius: 10px; background: var(--surface2);
  padding: 24px 20px; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  transition: all 0.2s; position: relative;
}
.sub-dropzone:hover          { border-color: rgba(6,182,212,0.35); background: rgba(6,182,212,0.04); }
.sub-dropzone--dragover      { border-color: #22d3ee; background: rgba(6,182,212,0.08); }
.sub-dropzone--error         { border-color: rgba(239,68,68,0.4); }
.sub-dropzone--ok            { border-color: rgba(16,185,129,0.4); border-style: solid; background: rgba(16,185,129,0.04); }
.sub-file-input              { position: absolute; inset: 0; opacity: 0; cursor: pointer; width: 100%; height: 100%; }
.sub-dropzone-content        { display: flex; flex-direction: column; align-items: center; gap: 6px; pointer-events: none; }
.sub-icon-upload             { color: var(--muted2); }
.sub-icon-ok                 { color: var(--success); }
.sub-dropzone-text           { font-size: 13px; font-weight: 500; color: var(--muted); }
.sub-dropzone-link           { color: #22d3ee; }
.sub-dropzone-filename       { font-size: 12.5px; font-weight: 600; color: #fff; }
.sub-dropzone-hint           { font-size: 11px; color: var(--muted2); }

/* ── Animations ───────────────────────────────────────────── */
.mv-spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.mv-modal-enter-active { animation: modal-in 0.22s cubic-bezier(0.34,1.56,0.64,1); }
.mv-modal-leave-active { animation: modal-out 0.18s ease; }
@keyframes modal-in  { from { opacity:0; transform:scale(0.94) translateY(10px); } }
@keyframes modal-out { to   { opacity:0; transform:scale(0.96) translateY(6px); } }
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to       { opacity: 0; }
</style>