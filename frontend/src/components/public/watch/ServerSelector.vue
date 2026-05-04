<script setup>
/**
 * ServerSelector — tombol pilih server stream.
 *
 * Props:
 *   servers      Array   — [{ n: Number, url: String }]
 *   modelValue   Number  — server aktif saat ini (v-model)
 *
 * Emits:
 *   update:modelValue — saat server dipilih
 */
defineProps({
  servers: {
    type: Array,
    required: true,
    // [{ n: 1, url: '...' }, { n: 2, url: '...' }]
  },
  modelValue: {
    type: Number,
    required: true,
  },
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <div v-if="servers.length > 0" class="ss-wrap">
    <span class="ss-label">Server:</span>
    <div class="ss-buttons">
      <button
        v-for="s in servers"
        :key="s.n"
        class="ss-btn"
        :class="{ 'ss-btn--active': modelValue === s.n }"
        @click="emit('update:modelValue', s.n)"
        :aria-label="`Server ${s.n}`"
        :aria-pressed="modelValue === s.n"
      >
        {{ s.n }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.ss-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.ss-label {
  font-size: 12px;
  font-weight: 500;
  color: #64748b;
  white-space: nowrap;
}

.ss-buttons {
  display: flex;
  gap: 4px;
}

.ss-btn {
  width: 32px; height: 32px;
  border-radius: 7px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
  color: #64748b;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s;
  display: grid; place-items: center;
}

.ss-btn:hover {
  color: #fff;
  border-color: rgba(255, 255, 255, 0.15);
  background: rgba(255, 255, 255, 0.08);
}

.ss-btn--active {
  background: #6366f1;
  border-color: transparent;
  color: #fff;
  box-shadow: 0 0 14px rgba(99, 102, 241, 0.4);
}
</style>