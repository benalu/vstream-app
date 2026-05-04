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
          drop_console: true,   // hapus console.log di production
          drop_debugger: true,
          pure_funcs: ['console.log', 'console.error', 'console.warn'],
        },
      },
      rollupOptions: {
        output: {
          manualChunks: {
            'vue-core': ['vue', 'vue-router'],
            'plyr':     ['plyr'],
            'lucide':   ['lucide-vue-next'],
            'axios':    ['axios'],
          },
        },
      },
      cssCodeSplit: true,
    },
})