<script setup>
/**
 * CategoryFilter — pill buttons for All / Movies / Series / Anime
 *
 * Props:
 *   modelValue  String  — currently active category id
 *   categories  Array   — [{ id, label }]
 *
 * Emits:
 *   update:modelValue  — when a category is clicked (v-model compatible)
 */
defineProps({
  modelValue: {
    type: String,
    required: true,
  },
  categories: {
    type: Array,
    default: () => [
      { id: 'all',    label: 'Semua'  },
      { id: 'movie',  label: 'Movies' },
      { id: 'series', label: 'Series' },
      { id: 'anime',  label: 'Anime'  },
    ],
  },
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <div class="cat-filter" role="tablist" aria-label="Filter kategori">
    <button
      v-for="cat in categories"
      :key="cat.id"
      class="cat-filter__btn"
      :class="{ 'cat-filter__btn--active': modelValue === cat.id }"
      role="tab"
      :aria-selected="modelValue === cat.id"
      @click="emit('update:modelValue', cat.id)"
    >
      {{ cat.label }}
    </button>
  </div>
</template>

<style scoped>
.cat-filter {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.cat-filter__btn {
  padding: 7px 18px;
  border-radius: 20px;
  border: 1px solid rgba(255,255,255,0.07);
  background: transparent;
  color: #64748b;
  font-size: 13px;
  font-weight: 500;
  font-family: 'DM Sans', sans-serif;
  cursor: pointer;
  transition: all 0.15s;
}
.cat-filter__btn:hover {
  color: #fff;
  border-color: rgba(255,255,255,0.15);
}
.cat-filter__btn--active {
  background: #6366f1;
  border-color: transparent;
  color: #fff;
  box-shadow: 0 0 18px rgba(99,102,241,0.3);
}
</style>