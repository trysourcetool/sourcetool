import arg from 'arg';
import { build, context } from 'esbuild';
import type { Plugin, PluginBuild, BuildOptions } from 'esbuild';
import * as glob from 'glob';
import path from 'path';
import fs from 'fs';
import { exec } from 'child_process';
import { promisify } from 'util';

const entryPoints = glob.sync('./src/index.ts');

const args = arg({
  '--watch': Boolean,
});

const isWatch = args['--watch'] || false;

const commonOptions: BuildOptions = {
  entryPoints,
  logLevel: 'info',
  platform: 'node',
};

const buildOptions = (format: 'esm' | 'cjs'): BuildOptions => {
  return {
    ...commonOptions,
    outbase: './src',
    outdir: format === 'esm' ? './dist/esm' : './dist/cjs',
    format,
    bundle: format === 'esm',
    outExtension: { '.js': format === 'esm' ? '.mjs' : '.js' },
  };
};

// Function to generate type definitions using tsc
const generateTypeDefinitions = async () => {
  const execPromise = promisify(exec);
  console.log('Generating type definitions...');

  try {
    // Run tsc to generate declaration files
    await execPromise(
      `npx tsc ${isWatch ? '--watch' : ''} --declaration --emitDeclarationOnly --project tsconfig.build.json`,
    );

    console.log('Type definitions generated successfully in ./dist/types');
  } catch (error) {
    console.error('Error generating type definitions:', error);
    throw error;
  }
};

if (isWatch) {
  const ctxs = await Promise.all([
    context(buildOptions('esm')),
    context(buildOptions('cjs')),
  ]);
  for (const ctx of ctxs) {
    ctx.watch();
  }

  await generateTypeDefinitions();

  // Note: Type definitions are not automatically regenerated in watch mode
  // You would need to implement a file watcher for .ts files to regenerate types
  console.log(
    'Watch mode active. Type definitions will not be automatically regenerated.',
  );
} else {
  await Promise.all([build(buildOptions('esm')), build(buildOptions('cjs'))]);
  await generateTypeDefinitions();
}
