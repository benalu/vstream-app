import { ref, onMounted, onUnmounted } from 'vue'

/**
 * Tracks whether the page has been scrolled past a threshold.
 * @param {number} threshold - scroll distance in px (default 60)
 */
export function useScrolled(threshold = 60) {
  const scrolled = ref(false)

  const onScroll = () => {
    scrolled.value = window.scrollY > threshold
  }

  onMounted(() => window.addEventListener('scroll', onScroll, { passive: true }))
  onUnmounted(() => window.removeEventListener('scroll', onScroll))

  return { scrolled }
}