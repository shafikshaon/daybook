import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  base: '/',
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        api: 'modern-compiler',
        silenceDeprecations: ['legacy-js-api', 'import', 'global-builtin', 'color-functions']
      }
    }
  },
  build: {
    // Generate proper hashed filenames for cache busting
    rollupOptions: {
      output: {
        // Add timestamp to chunk names for better cache busting
        entryFileNames: `assets/[name]-[hash].js`,
        chunkFileNames: `assets/[name]-[hash].js`,
        assetFileNames: `assets/[name]-[hash].[ext]`
      }
    },
    // Generate source maps for debugging
    sourcemap: false,
    // Ensure all assets get proper hashing
    assetsInlineLimit: 0
  },
  server: {
    port: 3000,
    host: true,
    strictPort: true,
    hmr: {
      overlay: true
    },
    watch: {
      usePolling: false
    }
  }
})
