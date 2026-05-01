import { ref, computed } from 'vue'
import axios from 'axios'

// Public API instance — tidak butuh credentials
const publicApi = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api',
})

/**
 * useContent — fetches and manages public content data.
 *
 * Usage:
 *   const { featured, top10, movies, series, anime, loading, error, fetchAll } = useContent()
 *   await fetchAll()
 *
 * Individual fetchers are exposed too for targeted refreshes:
 *   await fetchFeatured()
 *   await fetchTop10()
 *   await fetchListing('movie')  // or 'series' | 'anime'
 */
export function useContent() {
  // ── State ────────────────────────────────────────────────────
  const featured     = ref(null)
  const heroSlides   = ref([])
  const top10        = ref([])
  const movies       = ref([])
  const series       = ref([])
  const anime        = ref([])

  // Per-request loading flags
  const loadingFeatured    = ref(false)
  const loadingHeroSlides  = ref(false)
  const loadingTop10       = ref(false)
  const loadingMovies      = ref(false)
  const loadingSeries      = ref(false)
  const loadingAnime       = ref(false)

  // Errors — null means no error
  const errorFeatured     = ref(null)
  const errorHeroSlides   = ref(null)
  const errorTop10        = ref(null)
  const errorMovies       = ref(null)
  const errorSeries       = ref(null)
  const errorAnime        = ref(null)

  // Convenience: any fetch in-flight
  const loading = computed(() =>
    loadingFeatured.value ||
    loadingTop10.value    ||
    loadingMovies.value   ||
    loadingSeries.value   ||
    loadingAnime.value
  )

  // First error found (for simple error display)
  const error = computed(() =>
    errorFeatured.value ||
    errorTop10.value    ||
    errorMovies.value   ||
    errorSeries.value   ||
    errorAnime.value    ||
    null
  )

  // ── Fetchers ─────────────────────────────────────────────────

  const fetchFeatured = async () => {
    loadingFeatured.value = true
    errorFeatured.value = null
    try {
      const res = await publicApi.get('/featured')
      featured.value = res.data.data  // null if 204 / no content
    } catch (err) {
      errorFeatured.value = err.message
      console.error('[useContent] fetchFeatured:', err)
    } finally {
      loadingFeatured.value = false
    }
  }
  const fetchHeroSlides = async () => {
    loadingHeroSlides.value = true
    errorHeroSlides.value = null
    try {
      const res = await publicApi.get('/hero-slides')
      heroSlides.value = res.data.data ?? []
      // featured tetap pakai slide pertama untuk backward compat
      featured.value = heroSlides.value[0]?.movie ?? null
    } catch (err) {
      errorHeroSlides.value = err.message
      console.error('[useContent] fetchHeroSlides:', err)
    } finally {
      loadingHeroSlides.value = false
    }
  }

  const fetchTop10 = async () => {
    loadingTop10.value = true
    errorTop10.value = null
    try {
      const res = await publicApi.get('/top10')
      top10.value = res.data.data ?? []
    } catch (err) {
      errorTop10.value = err.message
      console.error('[useContent] fetchTop10:', err)
    } finally {
      loadingTop10.value = false
    }
  }

  /**
   * fetchListing — fetch a typed listing.
   * @param {'movie'|'series'|'anime'} type
   * @param {{ limit?: number, offset?: number }} options
   */
  const fetchListing = async (type, { limit = 20, offset = 0 } = {}) => {
    const ENDPOINT   = { movie: '/movies', series: '/series', anime: '/anime' }
    const LOAD_FLAG  = { movie: loadingMovies, series: loadingSeries, anime: loadingAnime }
    const ERROR_FLAG = { movie: errorMovies,   series: errorSeries,   anime: errorAnime   }
    const DATA_REF   = { movie: movies,        series: series,        anime: anime        }

    if (!ENDPOINT[type]) {
      console.warn('[useContent] unknown type:', type)
      return
    }

    LOAD_FLAG[type].value = true
    ERROR_FLAG[type].value = null

    try {
      const res = await publicApi.get(ENDPOINT[type], { params: { limit, offset } })
      DATA_REF[type].value = res.data.data ?? []
    } catch (err) {
      ERROR_FLAG[type].value = err.message
      console.error(`[useContent] fetchListing(${type}):`, err)
    } finally {
      LOAD_FLAG[type].value = false
    }
  }

  /**
   * fetchAll — kicks off featured + top10 + all three listings in parallel.
   * Use this on Home page mount.
   */
  const fetchAll = () => Promise.all([
    fetchHeroSlides(),
    fetchTop10(),
    fetchListing('movie'),
    fetchListing('series'),
    fetchListing('anime'),
  ])

  /**
   * recentlyAdded — cross-type list sorted by insertion order (already desc from API).
   * Merges movies + series + anime and takes the first `n`.
   */
  const recentlyAdded = computed(() => {
    return [...movies.value, ...series.value, ...anime.value]
      .slice(0, 12)
  })

  // ── Expose ───────────────────────────────────────────────────
  return {
    // data
    featured,
    heroSlides,
    top10,
    movies,
    series,
    anime,
    recentlyAdded,

    // loading flags
    loading,
    loadingFeatured,
    loadingTop10,
    loadingMovies,
    loadingSeries,
    loadingAnime,

    // errors
    error,
    errorFeatured,
    errorTop10,
    errorMovies,
    errorSeries,
    errorAnime,

    // actions
    fetchAll,
    fetchFeatured,
    fetchHeroSlides,
    fetchTop10,
    fetchListing,
  }
}