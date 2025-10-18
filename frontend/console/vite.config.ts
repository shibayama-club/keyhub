import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  plugins: [react(), tailwindcss()],
  base: '/console/',
  server: {
    port: 5174,
    proxy: {
      '/keyhub.console': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
    },
  },
});
