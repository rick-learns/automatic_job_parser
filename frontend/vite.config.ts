import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: '../jobsite_golang/public/latest',
    emptyOutDir: false, // CRITICAL: don't nuke jobs.json/csv
  },
  base: './',
  server: {
    proxy: {
      '/latest': { target: 'http://localhost:8080', changeOrigin: true },
    },
  },
})
