import { ref, computed, watch } from 'vue'

export const SUBTITLE_COLORS = [
  { label: 'Putih',  value: '#ffffff' },
  { label: 'Kuning', value: '#facc15' },
  { label: 'Hijau',  value: '#4ade80' },
  { label: 'Cyan',   value: '#22d3ee' },
  { label: 'Merah',  value: '#f87171' },
]

export function useSubtitleSettings() {
  const subtitleSize = ref(
    Number(localStorage.getItem('vstream_sub_size')) || 25
  )
  const subtitleColor = ref(
    localStorage.getItem('vstream_sub_color') || '#ffffff'
  )

  watch(subtitleSize, v => localStorage.setItem('vstream_sub_size', v))
  watch(subtitleColor, v => localStorage.setItem('vstream_sub_color', v))

  const subtitleStyle = computed(() => ({
    '--subtitle-size':  `${subtitleSize.value}px`,
    '--subtitle-color': subtitleColor.value,
  }))

  const increaseSubtitle = () => { if (subtitleSize.value < 150) subtitleSize.value += 2 }
  const decreaseSubtitle = () => { if (subtitleSize.value > 10)  subtitleSize.value -= 2 }

  return { subtitleSize, subtitleColor, subtitleStyle, increaseSubtitle, decreaseSubtitle }
}