import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import ViteFonts from 'unplugin-fonts/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    ViteFonts({
      fontsource: {
        families: [
          {
            name: 'Roboto',
            weights: [100, 300, 400, 500, 700, 900],
            styles: ['normal', 'italic'],
          },
        ],
      },
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },

  // --- START CONFIGURATION FOR GIN INTEGRATION ---
  build: {
    outDir: '../dist', 
    emptyOutDir: true, // Clean the directory before building
  },

  server: {
    // Crucial for development: Proxy API calls to the Gin server (running on 8080)
    proxy: {
      '/api': {
        target: 'http://localhost:8080', 
        changeOrigin: true,
      }
    }
  }
  // --- END CONFIGURATION FOR GIN INTEGRATION ---
})

