import { ref, computed } from 'vue'

export function useWatchData() {
  const info      = ref(null)
  const seasons   = ref([])
  const subtitles = ref([])
  const loading   = ref(true)
  const error     = ref(null)

  const hasEpisodes = computed(() =>
    info.value?.type === 'series' ||
    (info.value?.type === 'anime' && info.value?.has_episodes)
  )

  const currentSeason = (activeSeason) =>
    seasons.value.find(s => s.season_num === activeSeason.value)

  const fetchWatch = async (type, tmdbId) => {
    loading.value = true
    error.value   = null
    try {
      const base = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api'
      const res  = await fetch(`${base}/watch/${type}/${tmdbId}`)
      const data = await res.json()
      if (!data.success) throw new Error(data.error)
      info.value      = data.data.info
      seasons.value   = data.data.seasons ?? []
      subtitles.value = data.data.subtitles ?? []
    } catch (e) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  return { info, seasons, subtitles, loading, error, hasEpisodes, currentSeason, fetchWatch }
}