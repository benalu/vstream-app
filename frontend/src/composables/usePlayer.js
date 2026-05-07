import { ref, onBeforeUnmount } from 'vue'

export function usePlayer({
  subtitleSize,
  subtitleColor,
  increaseSubtitle,
  decreaseSubtitle,
  setSubtitleColor,
  SUBTITLE_COLORS,
  onReady,
  onTimeUpdate,
  onEnded,
  onError,
} = {}) {
  const plyrInstance = ref(null)
  let _clickOutsideHandler = null

  // ── Plyr config ────────────────────────────────────────────
  const PLYR_OPTIONS = {
    controls: `
      <div class="plyr__controls">
        <div class="plyr__controls__row plyr__controls__row--top">
          <div class="plyr__progress">
            <input data-plyr="seek" type="range" min="0" max="100" step="0.01" value="0" autocomplete="off" role="slider" aria-label="Seek" aria-valuemin="0" aria-valuemax="100" aria-valuenow="0">
            <progress class="plyr__progress__buffer" min="0" max="100" value="0">% buffered</progress>
            <span role="tooltip" class="plyr__tooltip">00:00</span>
          </div>
          <div class="plyr__time-display">
            <div class="plyr__time plyr__time--current" aria-label="Current time" role="timer">00:00</div>
            <span class="plyr__time-sep">/</span>
            <div class="plyr__time plyr__time--duration" aria-label="Duration">00:00</div>
          </div>
        </div>
        <div class="plyr__controls__row plyr__controls__row--bottom">
          <button type="button" class="plyr__control" aria-label="Play, {title}" data-plyr="play">
            <svg class="icon--pressed" role="presentation"><use xlink:href="#plyr-pause"></use></svg>
            <svg class="icon--not-pressed" role="presentation"><use xlink:href="#plyr-play"></use></svg>
          </button>
          <button type="button" class="plyr__control" aria-label="Mundur 10 detik" data-plyr="rewind">
            <svg role="presentation"><use xlink:href="#plyr-rewind"></use></svg>
          </button>
          <button type="button" class="plyr__control" aria-label="Maju 10 detik" data-plyr="fast-forward">
            <svg role="presentation"><use xlink:href="#plyr-fast-forward"></use></svg>
          </button>
          <div class="plyr__volume">
            <button type="button" class="plyr__control" data-plyr="mute">
              <svg class="icon--pressed" role="presentation"><use xlink:href="#plyr-muted"></use></svg>
              <svg class="icon--not-pressed" role="presentation"><use xlink:href="#plyr-volume"></use></svg>
            </button>
            <input data-plyr="volume" type="range" min="0" max="1" step="0.05" value="1" autocomplete="off" aria-label="Volume">
          </div>
          <div style="flex:1"></div>
          <button type="button" class="plyr__control plyr-custom-settings-btn" aria-label="Pengaturan">
            <svg role="presentation"><use xlink:href="#plyr-settings"></use></svg>
          </button>
          <button type="button" class="plyr__control" data-plyr="captions">
            <svg class="icon--pressed" role="presentation"><use xlink:href="#plyr-captions-on"></use></svg>
            <svg class="icon--not-pressed" role="presentation"><use xlink:href="#plyr-captions-off"></use></svg>
          </button>
          <button type="button" class="plyr__control" data-plyr="fullscreen">
            <svg class="icon--pressed" role="presentation"><use xlink:href="#plyr-exit-fullscreen"></use></svg>
            <svg class="icon--not-pressed" role="presentation"><use xlink:href="#plyr-enter-fullscreen"></use></svg>
          </button>
        </div>
      </div>
    `,
    captions:  { active: true, language: 'id', update: true },
    speed:     { selected: 1, options: [0.5, 0.75, 1, 1.25, 1.5, 2] },
    keyboard:  { focused: true, global: true },
    tooltips:  { controls: true, seek: true },
    i18n: {
      restart:         'Mulai Ulang',
      rewind:          'Mundur {seektime}s',
      play:            'Putar',
      pause:           'Jeda',
      fastForward:     'Maju {seektime}s',
      seek:            'Cari',
      seekLabel:       '{currentTime} dari {duration}',
      played:          'Diputar',
      buffered:        'Buffered',
      currentTime:     'Waktu saat ini',
      duration:        'Durasi',
      volume:          'Volume',
      mute:            'Bisukan',
      unmute:          'Aktifkan suara',
      enableCaptions:  'Aktifkan subtitle',
      disableCaptions: 'Nonaktifkan subtitle',
      enterFullscreen: 'Layar penuh',
      exitFullscreen:  'Keluar layar penuh',
      captions:        'Subtitle',
      settings:        'Pengaturan',
      pip:             'PiP',
      menuBack:        'Kembali',
      speed:           'Kecepatan',
      normal:          'Normal',
    },
  }

  // ── Init ───────────────────────────────────────────────────
  const initPlyr = async (videoEl) => {
    destroyPlyr()

    const [{ default: Plyr }] = await Promise.all([
      import('plyr'),
      import('plyr/dist/plyr.css'),
    ])

    plyrInstance.value = new Plyr(videoEl, PLYR_OPTIONS)

    plyrInstance.value.on('ready', () => {
      if (plyrInstance.value?.captions) {
        plyrInstance.value.captions.active = true
      }
      _buildCustomSettingsPanel()
      onReady?.()
    })

    plyrInstance.value.on('timeupdate', () => onTimeUpdate?.())
    plyrInstance.value.on('ended',      () => onEnded?.())
    plyrInstance.value.on('error',      () => onError?.())
  }

  // ── Destroy ────────────────────────────────────────────────
  const destroyPlyr = () => {
    if (_clickOutsideHandler) {
      document.removeEventListener('click', _clickOutsideHandler)
      _clickOutsideHandler = null
    }
    plyrInstance.value?.destroy()
    plyrInstance.value = null
  }

  // ── Build custom settings panel ────────────────────────────
  const _buildCustomSettingsPanel = () => {
  const container = plyrInstance.value?.elements?.container
  if (!container) return

  container.querySelector('.plyr-custom-menu')?.remove()

  const panel = document.createElement('div')
  panel.className = 'plyr-custom-menu'
  panel.hidden = true
  panel.innerHTML = `
    <div class="plyr-custom-menu__inner">

      <!-- Tab Header -->
      <div class="plyr-tabs">
        <button type="button" class="plyr-tab-btn" data-tab="speed">Kecepatan</button>
        <button type="button" class="plyr-tab-btn plyr-tab-btn--active" data-tab="subtitle">Subtitle</button>
        <button type="button" class="plyr-tab-close" aria-label="Tutup">✕</button>
      </div>

      <!-- Tab: Speed -->
      <div class="plyr-tab-panel" data-panel="speed" hidden>
        <div class="plyr-menu-list">
          ${[0.5, 0.75, 1, 1.25, 1.5, 2].map(s => `
            <button type="button" class="plyr-menu-item plyr-speed-btn ${s === 1 ? 'plyr-menu-item--active' : ''}" data-speed="${s}">
              <span>${s === 1 ? 'Normal' : s + 'x'}</span>
              <svg class="plyr-menu-check" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
            </button>
          `).join('')}
        </div>
      </div>

      <!-- Tab: Subtitle -->
      <div class="plyr-tab-panel" data-panel="subtitle">
        <div class="plyr-menu-list">

          <!-- Font Size -->
          <div class="plyr-menu-slider-row">
            <div class="plyr-menu-slider-header">
              <span class="plyr-menu-slider-icon">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="4 7 4 4 20 4 20 7"/><line x1="9" y1="20" x2="15" y2="20"/><line x1="12" y1="4" x2="12" y2="20"/>
                </svg>
              </span>
              <span class="plyr-menu-slider-label">Font Size</span>
              <span class="plyr-menu-slider-value plyr-sub-size-display">${subtitleSize.value}px</span>
            </div>
            <input type="range" class="plyr-menu-slider plyr-sub-size-slider" min="10" max="150" step="2" value="${subtitleSize.value}" />
          </div>

          <!-- Background Blur -->
          <div class="plyr-menu-slider-row">
            <div class="plyr-menu-slider-header">
              <span class="plyr-menu-slider-icon">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="3"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
                </svg>
              </span>
              <span class="plyr-menu-slider-label">Background Blur</span>
              <span class="plyr-menu-slider-value plyr-sub-blur-display">0%</span>
            </div>
            <input type="range" class="plyr-menu-slider plyr-sub-blur-slider" min="0" max="100" step="1" value="0" />
          </div>

          <!-- Color -->
          <div class="plyr-menu-color-section">
            <div class="plyr-menu-slider-header" style="margin-bottom:10px">
              <span class="plyr-menu-slider-icon">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="3"/>
                </svg>
              </span>
              <span class="plyr-menu-slider-label">Color</span>
            </div>
            <div class="plyr-sub-color-grid">
              ${SUBTITLE_COLORS.map(c => `
                <button
                  type="button"
                  class="plyr-sub-color-btn ${subtitleColor.value === c.value ? 'plyr-sub-color-btn--active' : ''}"
                  data-color="${c.value}"
                  title="${c.label}"
                  style="background:${c.value}; ${c.value === '#ffffff' ? 'outline: 1px solid #444' : ''}"
                ></button>
              `).join('')}
            </div>
          </div>

        </div>
      </div>

    </div>
  `

  container.appendChild(panel)

  // ── Tab switching ──
  const tabBtns   = panel.querySelectorAll('.plyr-tab-btn')
  const tabPanels = panel.querySelectorAll('.plyr-tab-panel')

  tabBtns.forEach(btn => {
    btn.addEventListener('click', () => {
      const target = btn.dataset.tab
      tabBtns.forEach(b => b.classList.toggle('plyr-tab-btn--active', b.dataset.tab === target))
      tabPanels.forEach(p => { p.hidden = p.dataset.panel !== target })
    })
  })

  // ── Close button ──
  panel.querySelector('.plyr-tab-close')?.addEventListener('click', () => {
    panel.hidden = true
  })

  // ── Toggle panel ──
  const settingsBtn = container.querySelector('.plyr-custom-settings-btn')
  settingsBtn?.addEventListener('click', (e) => {
    e.stopPropagation()
    panel.hidden = !panel.hidden
  })

  // ── Click outside ──
  if (_clickOutsideHandler) {
    document.removeEventListener('click', _clickOutsideHandler)
  }
  _clickOutsideHandler = (e) => {
    if (!panel.hidden && !panel.contains(e.target) && e.target !== settingsBtn) {
      panel.hidden = true
    }
  }
  document.addEventListener('click', _clickOutsideHandler)

  // ── Speed ──
  panel.querySelectorAll('.plyr-speed-btn').forEach(b => {
    b.addEventListener('click', () => {
      const speed = parseFloat(b.dataset.speed)
      if (plyrInstance.value) plyrInstance.value.speed = speed
      panel.querySelectorAll('.plyr-speed-btn').forEach(x =>
        x.classList.toggle('plyr-menu-item--active', x.dataset.speed === b.dataset.speed)
      )
    })
  })

  // ── Subtitle size buttons ──
  panel.querySelectorAll('.plyr-sub-size-btn').forEach(b => {
    b.addEventListener('click', () => {
      if (b.dataset.action === 'increase') increaseSubtitle?.()
      else decreaseSubtitle?.()
      _syncSubtitleUI(panel)
    })
  })

  // ── Subtitle slider ──
  panel.querySelector('.plyr-sub-size-slider')?.addEventListener('input', (e) => {
    subtitleSize.value = Number(e.target.value)
    _syncSubtitleUI(panel)
  })

   // ── Subtitle blur ──
  panel.querySelector('.plyr-sub-blur-slider')?.addEventListener('input', (e) => {
    const val = Number(e.target.value)
    const display = panel.querySelector('.plyr-sub-blur-display')
    if (display) display.textContent = `${val}%`

    // val 0–100 → opacity background 0–0.9
    const bgOpacity = (val / 100 * 0.9).toFixed(2)
    const container = plyrInstance.value?.elements?.container
    
    // Inject atau update style tag khusus untuk caption background
    let styleTag = container?.querySelector('#plyr-caption-blur-style')
    if (!styleTag) {
      styleTag = document.createElement('style')
      styleTag.id = 'plyr-caption-blur-style'
      container?.appendChild(styleTag)
    }
    styleTag.textContent = val > 0
      ? `.plyr__caption { background: rgba(0,0,0,${bgOpacity}) !important; }`
      : `.plyr__caption { background: rgba(0,0,0,0.35) !important; }`
  })

  // ── Subtitle color ──
  panel.querySelectorAll('.plyr-sub-color-btn').forEach(b => {
    b.addEventListener('click', () => {
      setSubtitleColor?.(b.dataset.color)
      panel.querySelectorAll('.plyr-sub-color-btn').forEach(x =>
        x.classList.toggle('plyr-sub-color-btn--active', x.dataset.color === b.dataset.color)
      )
    })
  })
}

  // ── Sync subtitle UI state ─────────────────────────────────
  const _syncSubtitleUI = (wrapper) => {
    const display = wrapper?.querySelector('.plyr-sub-size-display')
    if (display) display.textContent = `${subtitleSize.value}px`

    const slider = wrapper?.querySelector('.plyr-sub-size-slider')
    if (slider) slider.value = subtitleSize.value
  }

  // ── Auto-cleanup on unmount ────────────────────────────────
  onBeforeUnmount(destroyPlyr)

  return {
    plyrInstance,
    initPlyr,
    destroyPlyr,
  }
}