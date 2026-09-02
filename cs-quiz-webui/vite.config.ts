import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'node:path'

const basePath = process.env.VITE_BASE_PATH || '/'
const baseNoSlash = basePath.replace(/\/$/, '') || ''
const apiProxyTarget = process.env.VITE_API_PROXY || 'http://127.0.0.1:8090'
const hmrClientPort = Number(process.env.VITE_HMR_CLIENT_PORT || 5173)

export default defineConfig({
  base: basePath,
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: true,
    port: 5173,
    allowedHosts: true,
    hmr: { clientPort: hmrClientPort },
    proxy: {
      [`${baseNoSlash}/api`]: {
        target: apiProxyTarget,
        changeOrigin: true,
        rewrite: (p) => p.replace(new RegExp(`^${baseNoSlash}/api`), '/api'),
      },
      [`${baseNoSlash}/health`]: {
        target: apiProxyTarget,
        changeOrigin: true,
        rewrite: (p) => p.replace(new RegExp(`^${baseNoSlash}/health`), '/health'),
      },
      [`${baseNoSlash}/uploads`]: {
        target: apiProxyTarget,
        changeOrigin: true,
        rewrite: (p) => p.replace(new RegExp(`^${baseNoSlash}/uploads`), '/uploads'),
      },
    },
  },
})
