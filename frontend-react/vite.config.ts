import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
// Using URL-based resolution to avoid depending on Node type definitions
const r = (p: string) => new URL(p, import.meta.url).pathname;

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@services': r('./src/services'),
      '@components': r('./src/components'),
      '@pages': r('./src/pages'),
      '@ui': r('./src/components/ui'),
      '@shared': r('./src/shared')
    }
  }
});
