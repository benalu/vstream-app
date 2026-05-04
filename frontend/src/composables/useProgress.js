import { ref, computed } from 'vue'

/**
 * useProgress — save & restore video playback progress via localStorage.
 *
 * Usage:
 *   const { showResumePrompt, resumeTime, resumeFormatted,
 *           saveProgress, restoreProgress, resume, startOver } = useProgress(storageKey, plyrInstance)
 *
 * @param {import('vue').ComputedRef<string|null>} storageKey  — reactive key (e.g. `vstream_progress_${id}`)
 * @param {import('vue').Ref}                      player      — ref to Plyr instance
 */
export function useProgress(storageKey, player) {
  const showResumePrompt = ref(false)
  const resumeTime       = ref(0)

  // ── Helpers ────────────────────────────────────────────────
  const EXPIRY_MS = 30 * 24 * 60 * 60 * 1000 // 30 days

  const resumeFormatted = computed(() => {
    const m = Math.floor(resumeTime.value / 60)
    const s = String(Math.floor(resumeTime.value % 60)).padStart(2, '0')
    return `${m}:${s}`
  })

  // ── Save ───────────────────────────────────────────────────
  const saveProgress = () => {
    if (!player.value || !storageKey.value) return

    const t = player.value.currentTime
    const d = player.value.duration
    if (t < 5 || !d) return

    // If >95% watched, clear saved progress
    if (t / d > 0.95) {
      localStorage.removeItem(storageKey.value)
      return
    }

    localStorage.setItem(
      storageKey.value,
      JSON.stringify({ time: t, duration: d, saved_at: Date.now() })
    )
  }

  // ── Restore ────────────────────────────────────────────────
  const restoreProgress = () => {
    if (!storageKey.value) return

    const raw = localStorage.getItem(storageKey.value)
    if (!raw) return

    try {
      const { time, duration, saved_at } = JSON.parse(raw)
      const expired  = Date.now() - saved_at > EXPIRY_MS
      const finished = duration && time / duration > 0.95
      if (expired || finished || time < 5) return

      resumeTime.value       = time
      showResumePrompt.value = true
    } catch {
      localStorage.removeItem(storageKey.value)
    }
  }

  // ── Resume / Start Over ────────────────────────────────────
  const resume = () => {
    if (player.value) {
      player.value.currentTime = resumeTime.value
      player.value.play()
    }
    showResumePrompt.value = false
  }

  const startOver = () => {
    showResumePrompt.value = false
    resumeTime.value       = 0
  }

  // ── Clear (on ended) ──────────────────────────────────────
  const clearProgress = () => {
    if (storageKey.value) localStorage.removeItem(storageKey.value)
    showResumePrompt.value = false
  }

  return {
    showResumePrompt,
    resumeTime,
    resumeFormatted,
    saveProgress,
    restoreProgress,
    resume,
    startOver,
    clearProgress,
  }
}