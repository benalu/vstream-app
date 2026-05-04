import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/static/subtitles': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
    build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
        pure_funcs: ['console.log', 'console.error', 'console.warn'],
      },
    },
    rollupOptions: {
      output: {
        manualChunks(id) {
          // Vue core — needed everywhere
          if (id.includes('node_modules/vue') || id.includes('node_modules/vue-router')) {
            return 'vue-core'
          }
          // Plyr is heavy (~300KB) and only used on the Watch page
          // It's already lazy-imported in usePlayer.js via import('plyr')
          // so Rollup should already split it — this just makes it explicit
          if (id.includes('node_modules/plyr')) {
            return 'plyr'
          }
          // Lucide icons — only needed in admin/UI
          if (id.includes('node_modules/lucide')) {
            return 'lucide'
          }
          // Axios
          if (id.includes('node_modules/axios')) {
            return 'axios'
          }
          // Fontsource — split per font to avoid loading unused weights
          if (id.includes('@fontsource')) {
            return 'fonts'
          }
        },
      },
    },
    cssCodeSplit: true,
  },
})