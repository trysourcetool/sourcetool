import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';
import react from '@vitejs/plugin-react';
import { TanStackRouterVite } from '@tanstack/router-plugin/vite';

export default defineConfig(({ mode }) => {
  console.log({ mode });
  return {
    server: {
      host: '127.0.0.1',
      port: 5173,
      allowedHosts: true,
    },
    plugins: [
      TanStackRouterVite({
        target: 'react',
        autoCodeSplitting: true,
        virtualRouteConfig: './app/routes.ts',
        generatedRouteTree: './app/routeTree.gen.ts',
        routesDirectory: './app/routes/',
      }),
      react(),
      tsconfigPaths(),
    ],
    esbuild: {
      drop: mode === 'production' ? ['console', 'debugger'] : [],
    },
  };
});
