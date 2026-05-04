import { ref, computed, watch } from 'vue'

/**
 * useSubtitle — manage subtitle size & color, persisted to localStorage.
 *
 * Usage:
 *   const { subtitleSize, subtitleColor, subtitleStyle,
 *           increaseSubtitle, decreaseSubtitle,
 *           SUBTITLE_COLORS } = useSubtitle()
 */
export function useSubtitle() {
  // ── Constants ──────────────────────────────────────────────
  const SUBTITLE_COLORS = [
    { label: 'Putih',  value: '#ffffff' },
    { label: 'Kuning', value: '#facc15' },
    { label: 'Hijau',  value: '#4ade80' },
    { label: 'Cyan',   value: '#22d3ee' },
    { label: 'Merah',  value: '#f87171' },
  ]

  const SIZE_MIN     = 10
  const SIZE_MAX     = 150
  const SIZE_STEP    = 2
  const SIZE_DEFAULT = 25
  const COLOR_DEFAULT = '#ffffff'

  // ── State (restored from localStorage) ────────────────────
  const subtitleSize = ref(
    Number(localStorage.getItem('vstream_sub_size')) || SIZE_DEFAULT
  )

  const subtitleColor = ref(
    localStorage.getItem('vstream_sub_color') || COLOR_DEFAULT
  )

  // ── Persist on change ──────────────────────────────────────
  watch(subtitleSize, (val) => localStorage.setItem('vstream_sub_size', val))
  watch(subtitleColor, (val) => localStorage.setItem('vstream_sub_color', val))

  // ── CSS vars for the player wrapper ───────────────────────
  const subtitleStyle = computed(() => ({
    '--subtitle-size':  `${subtitleSize.value}px`,
    '--subtitle-color': subtitleColor.value,
  }))

  // ── Actions ────────────────────────────────────────────────
  const increaseSubtitle = () => {
    if (subtitleSize.value < SIZE_MAX) subtitleSize.value += SIZE_STEP
  }

  const decreaseSubtitle = () => {
    if (subtitleSize.value > SIZE_MIN) subtitleSize.value -= SIZE_STEP
  }

  const setSubtitleColor = (color) => {
    subtitleColor.value = color
  }

  return {
    subtitleSize,
    subtitleColor,
    subtitleStyle,
    increaseSubtitle,
    decreaseSubtitle,
    setSubtitleColor,
    SUBTITLE_COLORS,
    SIZE_MIN,
    SIZE_MAX,
    SIZE_STEP,
  }
}