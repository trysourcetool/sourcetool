import typescriptEslint from '@typescript-eslint/eslint-plugin';
import pluginVitest from 'eslint-plugin-vitest';
import tsParser from '@typescript-eslint/parser';

export default [
  {
    files: ['app/**/*.ts', 'app/**/*.tsx'],
  },
  {
    plugins: {
      '@typescript-eslint': typescriptEslint,
      vitest: pluginVitest,
    },

    languageOptions: {
      parser: tsParser,
      ecmaVersion: 2022,
      sourceType: 'module',
    },

    rules: {
      '@typescript-eslint/naming-convention': [
        'warn',
        {
          selector: 'import',
          format: ['camelCase', 'PascalCase'],
        },
      ],
      '@typescript-eslint/no-unused-vars': 'warn',
      curly: 'warn',
      eqeqeq: 'warn',
      'no-throw-literal': 'warn',
      semi: 'warn',
    },
    ignores: ['node_modules/**', '.docusaurus/**'],
  },
];
