<script setup>
/**
 * ReportBar — muncul saat video gagal diputar.
 *
 * Props:
 *   show          Boolean — tampilkan bar
 *   activeServer  Number  — server yang sedang aktif (untuk label)
 *
 * Emits:
 *   report — saat user klik tombol lapor
 */
defineProps({
  show:         { type: Boolean, required: true },
  activeServer: { type: Number,  default: 1 },
})

const emit = defineEmits(['report'])
</script>

<template>
  <Transition name="rb-fade">
    <div v-if="show" class="rb-bar">
      <div class="rb-left">
        <div class="rb-icon">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
        </div>
        <span class="rb-label">Video tidak bisa diputar?</span>
      </div>
      <button class="rb-btn" @click="emit('report')">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"/>
          <line x1="4" y1="22" x2="4" y2="15"/>
        </svg>
        Laporkan Server {{ activeServer }}
      </button>
    </div>
  </Transition>
</template>

<style scoped>
.rb-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 14px;
  border-radius: 8px;
  background: rgba(239, 68, 68, 0.06);
  border: 1px solid rgba(239, 68, 68, 0.15);
  margin-top: 8px;
}

.rb-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rb-icon {
  width: 24px; height: 24px;
  border-radius: 6px;
  background: rgba(239, 68, 68, 0.1);
  color: #f87171;
  display: grid; place-items: center;
  flex-shrink: 0;
}

.rb-label {
  font-size: 12.5px;
  color: #94a3b8;
}

.rb-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: 6px;
  border: 1px solid rgba(239, 68, 68, 0.25);
  background: rgba(239, 68, 68, 0.08);
  color: #f87171;
  font-size: 12px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.15s;
  white-space: nowrap;
  flex-shrink: 0;
}

.rb-btn:hover {
  background: rgba(239, 68, 68, 0.15);
  border-color: rgba(239, 68, 68, 0.4);
}

/* Transition */
.rb-fade-enter-active { animation: rb-in 0.2s ease; }
.rb-fade-leave-active { animation: rb-in 0.15s ease reverse; }
@keyframes rb-in {
  from { opacity: 0; transform: translateY(-4px); }
}
</style>