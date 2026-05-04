import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'

const app = createApp(App)
app.use(router)
app.mount('#app')

// Load fonts non-blocking setelah app mount
if (typeof requestIdleCallback !== 'undefined') {
  requestIdleCallback(() => {
    import(/* @vite-ignore */ '@fontsource/bebas-neue')
    import(/* @vite-ignore */ '@fontsource/dm-sans/400.css')
    import(/* @vite-ignore */ '@fontsource/dm-sans/500.css')
    import(/* @vite-ignore */ '@fontsource/dm-sans/600.css')
    import(/* @vite-ignore */ '@fontsource/dm-sans/700.css')
  })
} else {
  setTimeout(() => {
    import(/* @vite-ignore */ '@fontsource/bebas-neue')
    import(/* @vite-ignore */ '@fontsource/dm-sans/400.css')
    import(/* @vite-ignore */ '@fontsource/dm-sans/500.css')
    import(/* @vite-ignore */ '@fontsource/dm-sans/600.css')
    import(/* @vite-ignore */ '@fontsource/dm-sans/700.css')
  }, 0)
}