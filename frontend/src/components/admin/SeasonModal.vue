<script setup>
import { ref } from 'vue'
import { Tv, Plus, Trash2 } from 'lucide-vue-next'

defineProps({
  show: { type: Boolean, required: true },
})

const emit = defineEmits(['close', 'submit'])

const form = ref({
  season_num: 1,
  episodes: [{ ep_num: 1, title: '', url1: '', url2: '' }],
})

const addEpisode = () => {
  const last = form.value.episodes.at(-1)
  form.value.episodes.push({ ep_num: last.ep_num + 1, title: '', url1: '', url2: '' })
}

const removeEpisode = (i) => {
  if (form.value.episodes.length > 1) form.value.episodes.splice(i, 1)
}

const reset = (nextSeasonNum = 1) => {
  form.value = {
    season_num: nextSeasonNum,
    episodes: [{ ep_num: 1, title: '', url1: '', url2: '' }],
  }
}

defineExpose({ reset })
</script>

<template>
  <Teleport to="body">
    <Transition name="mv-modal">
      <div v-if="show" class="mv-backdrop" @click.self="emit('close')">
        <div class="mv-modal mv-modal--wide">
          <div class="mv-modal-header">
            <div class="mv-modal-title-wrap">
              <div class="mv-modal-icon" style="background:rgba(6,182,212,0.12);color:#22d3ee">
                <Tv :size="16" />
              </div>
              <h2 class="mv-modal-title">Tambah Season</h2>
            </div>
            <button class="mv-modal-close" @click="emit('close')">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <div class="mv-modal-body">
            <div class="mv-field">
              <label class="mv-label">Nomor Season</label>
              <input v-model.number="form.season_num" type="number" min="1" class="vs-input" style="max-width:120px" />
            </div>

            <div class="mv-field">
              <div class="mv-ep-list-header">
                <label class="mv-label">Daftar Episode</label>
                <button type="button" class="vs-btn vs-btn--ghost mv-btn-sm" @click="addEpisode">
                  <Plus :size="12" /> Tambah Episode
                </button>
              </div>
              <div class="mv-ep-list">
                <div v-for="(ep, i) in form.episodes" :key="i" class="mv-ep-item">
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
                    :disabled="form.episodes.length === 1"
                  >
                    <Trash2 :size="13" />
                  </button>
                </div>
              </div>
            </div>

            <div class="mv-modal-footer">
              <button type="button" class="vs-btn vs-btn--ghost" @click="emit('close')">Batal</button>
              <button type="button" class="vs-btn vs-btn--primary" @click="emit('submit', { ...form.value })">
                Simpan Season
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.mv-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,0.65); backdrop-filter: blur(6px); z-index: 100; display: flex; align-items: center; justify-content: center; padding: 16px; }
.mv-modal { background: var(--surface); border: 1px solid var(--border); border-radius: 16px; width: 100%; max-width: 460px; box-shadow: 0 25px 60px rgba(0,0,0,0.5); overflow: hidden; max-height: 90vh; display: flex; flex-direction: column; }
.mv-modal--wide { max-width: 600px; }
.mv-modal-header { display: flex; align-items: center; justify-content: space-between; padding: 18px 20px; border-bottom: 1px solid var(--border); flex-shrink: 0; }
.mv-modal-title-wrap { display: flex; align-items: center; gap: 10px; }
.mv-modal-icon { width: 32px; height: 32px; border-radius: 8px; display: grid; place-items: center; }
.mv-modal-title { font-size: 14px; font-weight: 600; color: #fff; }
.mv-modal-close { width: 28px; height: 28px; border-radius: 6px; border: 1px solid var(--border); background: transparent; color: var(--muted); display: grid; place-items: center; cursor: pointer; transition: all 0.15s; }
.mv-modal-close:hover { color: #fff; background: var(--border2); }
.mv-modal-body { padding: 20px; display: flex; flex-direction: column; gap: 16px; overflow-y: auto; }
.mv-field { display: flex; flex-direction: column; gap: 6px; }
.mv-label { font-size: 12px; font-weight: 500; color: var(--muted); }
.mv-modal-footer { display: flex; justify-content: flex-end; gap: 8px; padding-top: 16px; border-top: 1px solid var(--border); }

.mv-ep-list-header { display: flex; align-items: center; justify-content: space-between; }
.mv-ep-list { display: flex; flex-direction: column; gap: 10px; margin-top: 8px; }
.mv-ep-item { display: flex; align-items: flex-start; gap: 10px; background: var(--surface2); border: 1px solid var(--border); border-radius: 8px; padding: 10px 12px; }
.mv-ep-num-badge { width: 26px; height: 26px; border-radius: 6px; background: var(--border); color: var(--muted); display: grid; place-items: center; font-size: 11px; font-weight: 700; flex-shrink: 0; }
.mv-ep-inputs { flex: 1; display: flex; flex-direction: column; gap: 7px; }
.mv-ep-title-input { font-weight: 500 !important; }
.mv-ep-urls { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }

.mv-action-btn { width: 30px; height: 30px; border-radius: 7px; border: 1px solid var(--border); background: transparent; cursor: pointer; display: grid; place-items: center; transition: all 0.15s; color: var(--muted); }
.mv-action-btn:disabled { opacity: 0.3; cursor: not-allowed; }
.mv-action-btn--delete:hover { color: var(--danger); border-color: rgba(239,68,68,0.3); background: rgba(239,68,68,0.08); }
.mv-btn-sm { padding: 5px 10px !important; font-size: 11.5px !important; }

.mv-modal-enter-active { animation: modal-in 0.22s cubic-bezier(0.34,1.56,0.64,1); }
.mv-modal-leave-active { animation: modal-out 0.18s ease; }
@keyframes modal-in  { from { opacity:0; transform:scale(0.94) translateY(10px); } }
@keyframes modal-out { to   { opacity:0; transform:scale(0.96) translateY(6px); } }
</style>