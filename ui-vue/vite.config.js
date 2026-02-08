import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    proxy: {
      '/ui': {
        target: 'http://localhost:32241',
        changeOrigin: true,
      },
      '/api': {
        target: 'http://localhost:32241',
        changeOrigin: true,
      },
      '/setup': {
        target: 'http://localhost:32241',
        changeOrigin: true,
      },
      '/management': {
        target: 'http://localhost:32241',
        changeOrigin: true,
      }
    }
  }
})

