import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';
import { reactRouter } from '@react-router/dev/vite';

export default defineConfig(({ mode }) => {
  console.log({ mode });
  return {
    server: {
      host: '127.0.0.1',
      port: 5173,
      allowedHosts: true,
    },
    plugins: [reactRouter(), tsconfigPaths()],
    esbuild: {
      drop: mode === 'production' ? ['console', 'debugger'] : [],
    },
  };
});
