import globals from 'globals';
import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import pluginVitest from 'eslint-plugin-vitest';
import * as eslintPluginImport from 'eslint-plugin-import';
import pluginUnusedImports from 'eslint-plugin-unused-imports';
import pluginReact from 'eslint-plugin-react';
import pluginReactHooks from 'eslint-plugin-react-hooks';
import eslintConfigPrettier from 'eslint-config-prettier';
export default tseslint.config(
  eslint.configs.recommended,
  tseslint.configs.recommended,
  {
    name: 'ignore-global-rules',
    ignores: [
      '**/node_modules/**',
      '**/build',
      '**/dist',
      '**/.react-router',
      '**/.docusaurus',
    ],
  },
  {
    name: 'global-configs',
    languageOptions: {
      ecmaVersion: 2022,
      sourceType: 'module',
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
    plugins: {
      vitest: pluginVitest,
      react: pluginReact,
      'react-hooks': pluginReactHooks,
      import: eslintPluginImport,
      'unused-imports': pluginUnusedImports,
    },
  },
  {
    name: 'global-rules',
    rules: {
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/naming-convention': [
        'error',
        {
          selector: ['import', 'variable'],
          format: ['camelCase', 'PascalCase', 'UPPER_CASE'],
          leadingUnderscore: 'allow',
        },
      ],
      'no-unused-vars': 'error',
      'unused-imports/no-unused-imports': 'error',
      'unused-imports/no-unused-vars': [
        'error',
        {
          vars: 'all',
          varsIgnorePattern: '^_',
          args: 'after-used',
          argsIgnorePattern: '^_',
        },
      ],
      curly: 'error',
      eqeqeq: 'error',
      'no-throw-literal': 'warn',
      semi: 'error',
    },
  },
  {
    name: 'tsx-rules',
    files: ['**/*.tsx'],
    rules: {
      ...pluginReactHooks.configs.recommended.rules,
    },
  },
  {
    name: 'eslint-config-prettier',
    ...eslintConfigPrettier,
  },
);
