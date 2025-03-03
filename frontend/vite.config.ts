import { defineConfig, loadEnv } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';
import { reactRouter } from '@react-router/dev/vite';

export default defineConfig(({ mode }) => {
  const envPrefix = ['VITE_', 'APP_ENV'];
  const env = loadEnv(mode, '.', envPrefix);
  return {
    server: {
      host: '127.0.0.1',
      port: 5173,
      allowedHosts: true,
    },
    plugins: [reactRouter(), tsconfigPaths()],
    esbuild: {
      drop: env['APP_ENV'] === 'production' ? ['console', 'debugger'] : [],
    },
  };
});
