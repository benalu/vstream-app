import { ref } from 'vue'

const INTRO_START = 30
const INTRO_END   = 90

export function useSkipIntro(plyrInstance) {
  const showSkipIntro = ref(false)

  const checkSkipIntro = (hasEpisodes) => {
    if (!hasEpisodes) return
    const t = plyrInstance.value?.currentTime ?? 0
    showSkipIntro.value = t >= INTRO_START && t <= INTRO_END
  }

  const skipIntro = () => {
    if (plyrInstance.value) plyrInstance.value.currentTime = INTRO_END
    showSkipIntro.value = false
  }

  return { showSkipIntro, checkSkipIntro, skipIntro }
}