import { fileURLToPath, URL } from 'node:url';

import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  server: {
    port: 5173,
    // В разработке API поднят отдельно (task backend:dev); в продакшене /api проксирует nginx.
    proxy: { '/api': 'http://localhost:8080' },
  },
});
