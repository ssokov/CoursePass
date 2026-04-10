import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/v1': {
        target: 'http://localhost:8075',
        changeOrigin: true,
      },
      '/media': {
        target: 'http://localhost:8075',
        changeOrigin: true,
      },
    },
  },
})
