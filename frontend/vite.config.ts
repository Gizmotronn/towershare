import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const pbUrl = process.env.VITE_PB_URL ?? 'http://127.0.0.1:8090'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': { target: pbUrl, changeOrigin: true },
      '/_': { target: pbUrl, changeOrigin: true },
    },
  },
})
