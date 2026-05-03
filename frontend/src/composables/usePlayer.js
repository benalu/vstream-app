import { ref, nextTick } from 'vue'
import Plyr from 'plyr'
import 'plyr/dist/plyr.css'
import { SUBTITLE_COLORS } from './useSubtitleSettings'

export function usePlayer({ subtitleSize, subtitleColor, increaseSubtitle, decreaseSubtitle, restoreProgress, saveProgress, clearProgress, checkSkipIntro, hasEpisodes, reportError }) {
  const plyrInstance = ref(null)
  const videoRef     = ref(null)
  let _menuObserver  = null

  const destroyPlayer = () => {
    _menuObserver?.disconnect()
    _menuObserver = null
    plyrInstance.value?.destroy()
    plyrInstance.value = null
  }

  const watchPlyrMenu = () => {
    _menuObserver?.disconnect()
    const container = plyrInstance.value?.elements?.container
    if (!container) return
    _menuObserver = new MutationObserver(() => {
      const captionsPanel = container.querySelector('[id$="-captions"]:not([hidden])')
      if (captionsPanel && !captionsPanel.querySelector('.plyr-subtitle-size-panel')) {
        injectSubtitleControls(captionsPanel)
      }
    })
    _menuObserver.observe(container, { childList: true, subtree: true, attributes: true, attributeFilter: ['hidden'] })
  }

  const syncSubtitleUI = (wrapper) => {
    const display = wrapper?.querySelector('.plyr-sub-size-display')
    if (display) display.textContent = `${subtitleSize.value}px`
    const slider = wrapper?.querySelector('.plyr-sub-size-slider')
    if (slider) slider.value = subtitleSize.value
  }

  const injectSubtitleControls = (captionsPanel) => {
    const wrapper = document.createElement('div')
    wrapper.className = 'plyr-subtitle-size-panel'
    wrapper.innerHTML = `
      <div class="plyr-sub-size-label">Ukuran Subtitle</div>
      <div class="plyr-sub-size-row">
        <button type="button" class="plyr-sub-size-btn" data-action="decrease">A−</button>
        <span class="plyr-sub-size-display">${subtitleSize.value}px</span>
        <button type="button" class="plyr-sub-size-btn" data-action="increase">A+</button>
      </div>
      <div class="plyr-sub-size-slider-wrap">
        <span class="plyr-sub-size-min">10px</span>
        <input type="range" class="plyr-sub-size-slider" min="10" max="150" step="2" value="${subtitleSize.value}" />
        <span class="plyr-sub-size-max">150px</span>
      </div>
      <div class="plyr-sub-size-label" style="margin-top:4px">Warna Subtitle</div>
      <div class="plyr-sub-color-row">
        ${SUBTITLE_COLORS.map(c => `
          <button type="button" class="plyr-sub-color-btn ${subtitleColor.value === c.value ? 'plyr-sub-color-btn--active' : ''}"
            data-color="${c.value}" title="${c.label}" style="background:${c.value}"></button>
        `).join('')}
      </div>
    `
    captionsPanel.appendChild(wrapper)

    wrapper.querySelectorAll('.plyr-sub-size-btn').forEach(b => {
      b.addEventListener('click', () => {
        if (b.dataset.action === 'increase') increaseSubtitle()
        else decreaseSubtitle()
        syncSubtitleUI(wrapper)
      })
    })

    wrapper.querySelector('.plyr-sub-size-slider').addEventListener('input', e => {
      subtitleSize.value = Number(e.target.value)
      syncSubtitleUI(wrapper)
    })

    wrapper.querySelectorAll('.plyr-sub-color-btn').forEach(b => {
      b.addEventListener('click', () => {
        subtitleColor.value = b.dataset.color
        wrapper.querySelectorAll('.plyr-sub-color-btn').forEach(x =>
          x.classList.toggle('plyr-sub-color-btn--active', x.dataset.color === b.dataset.color)
        )
      })
    })
  }

  const initPlyr = async (subtitles) => {
    await nextTick()
    if (!videoRef.value) return
    destroyPlayer()

    plyrInstance.value = new Plyr(videoRef.value, {
      controls: ['play-large','play','rewind','fast-forward','progress','current-time','duration','mute','volume','captions','settings','pip','fullscreen'],
      settings: ['captions', 'speed'],
      captions: { active: true, language: 'id', update: true },
      speed:    { selected: 1, options: [0.5, 0.75, 1, 1.25, 1.5, 2] },
      keyboard: { focused: true, global: true },
      tooltips: { controls: true, seek: true },
      i18n: {
        restart: 'Mulai Ulang', rewind: 'Mundur {seektime}s', play: 'Putar',
        pause: 'Jeda', fastForward: 'Maju {seektime}s', seek: 'Cari',
        seekLabel: '{currentTime} dari {duration}', played: 'Diputar',
        buffered: 'Buffered', currentTime: 'Waktu saat ini', duration: 'Durasi',
        volume: 'Volume', mute: 'Bisukan', unmute: 'Aktifkan suara',
        enableCaptions: 'Aktifkan subtitle', disableCaptions: 'Nonaktifkan subtitle',
        enterFullscreen: 'Layar penuh', exitFullscreen: 'Keluar layar penuh',
        captions: 'Subtitle', settings: 'Pengaturan', pip: 'PiP',
        menuBack: 'Kembali', speed: 'Kecepatan', normal: 'Normal',
      },
    })

    plyrInstance.value.on('ready', () => {
      if (plyrInstance.value?.captions) plyrInstance.value.captions.active = true
      restoreProgress()
      watchPlyrMenu()
    })

    plyrInstance.value.on('timeupdate', () => {
      saveProgress()
      checkSkipIntro(hasEpisodes.value)
    })

    plyrInstance.value.on('ended', () => { clearProgress() })
    plyrInstance.value.on('error',  () => { reportError('load_error') })
  }

  return { videoRef, plyrInstance, initPlyr, destroyPlayer }
}