import { ref, computed } from 'vue'

export function useProgress(info, plyrInstance) {
  const showResumePrompt = ref(false)
  const resumeTime       = ref(0)

  const storageKey = computed(() =>
    info.value ? `vstream_progress_${info.value.id}` : null
  )

  const resumeFormatted = computed(() => {
    const m = Math.floor(resumeTime.value / 60)
    const s = String(Math.floor(resumeTime.value % 60)).padStart(2, '0')
    return `${m}:${s}`
  })

  const saveProgress = () => {
    if (!plyrInstance.value || !storageKey.value) return
    const t = plyrInstance.value.currentTime
    const d = plyrInstance.value.duration
    if (t < 5 || !d) return
    if (t / d > 0.95) { localStorage.removeItem(storageKey.value); return }
    localStorage.setItem(storageKey.value, JSON.stringify({
      time: t, duration: d, saved_at: Date.now(),
    }))
  }

  const restoreProgress = () => {
    if (!storageKey.value) return
    const raw = localStorage.getItem(storageKey.value)
    if (!raw) return
    try {
      const { time, duration, saved_at } = JSON.parse(raw)
      const expired  = Date.now() - saved_at > 30 * 24 * 60 * 60 * 1000
      const finished = duration && (time / duration) > 0.95
      if (expired || finished || time < 5) return
      resumeTime.value       = time
      showResumePrompt.value = true
    } catch { localStorage.removeItem(storageKey.value) }
  }

  const resume = () => {
    if (plyrInstance.value) {
      plyrInstance.value.currentTime = resumeTime.value
      plyrInstance.value.play()
    }
    showResumePrompt.value = false
  }

  const startOver    = () => { showResumePrompt.value = false; resumeTime.value = 0 }
  const clearProgress = () => { if (storageKey.value) localStorage.removeItem(storageKey.value) }

  return { showResumePrompt, resumeTime, resumeFormatted, saveProgress, restoreProgress, resume, startOver, clearProgress }
}