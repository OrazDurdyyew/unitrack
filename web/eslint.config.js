// ESLint: TypeScript strict с проверкой типов, правила Vue, границы слоёв FSD.
// Правила проекта описаны в web/CLAUDE.md.
import js from '@eslint/js';
import boundaries from 'eslint-plugin-boundaries';
import prettier from 'eslint-config-prettier';
import pluginVue from 'eslint-plugin-vue';
import tseslint from 'typescript-eslint';
import vueParser from 'vue-eslint-parser';

// Слои FSD сверху вниз. Слой импортирует только нижележащие; слайсы одного слоя изолированы;
// у слайсов (pages, widgets, features, entities) публичный API — только index.ts; shared открыт целиком.
const SLICED = ['pages', 'widgets', 'features', 'entities'];
const LAYERS = ['app', ...SLICED, 'shared'];

const fsdPolicies = LAYERS.map((layer, i) => {
  const lower = LAYERS.slice(i + 1);
  const targets = [
    ...lower
      .filter((l) => SLICED.includes(l))
      .map((type) => ({ element: { type, fileInternalPath: 'index.ts' } })),
    { element: { type: 'shared' } },
  ];
  return { from: { element: { type: layer } }, allow: { to: targets } };
});

export default tseslint.config(
  { ignores: ['dist/**', 'node_modules/**', 'coverage/**', 'src/shared/api/schema.d.ts'] },

  js.configs.recommended,
  ...tseslint.configs.strictTypeChecked,
  ...tseslint.configs.stylisticTypeChecked,
  ...pluginVue.configs['flat/recommended'],

  {
    files: ['**/*.{ts,vue}'],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: tseslint.parser,
        projectService: true,
        extraFileExtensions: ['.vue'],
        tsconfigRootDir: import.meta.dirname,
      },
    },
    plugins: { boundaries },
    settings: {
      'import/resolver': { typescript: { project: './tsconfig.app.json' } },
      'boundaries/include': ['src/**/*'],
      'boundaries/elements': [
        { type: 'app', pattern: 'src/app' },
        { type: 'pages', pattern: 'src/pages/*', capture: ['slice'] },
        { type: 'widgets', pattern: 'src/widgets/*', capture: ['slice'] },
        { type: 'features', pattern: 'src/features/*', capture: ['slice'] },
        { type: 'entities', pattern: 'src/entities/*', capture: ['slice'] },
        { type: 'shared', pattern: 'src/shared/*', capture: ['segment'] },
      ],
    },
    rules: {
      // Границы слоёв (только локальные импорты; внешние пакеты не ограничиваются)
      'boundaries/dependencies': [
        'error',
        {
          default: 'disallow',
          policies: fsdPolicies,
          message:
            'Нарушение FSD: {{from.element.type}} → {{to.element.type}} ({{to.element.fileInternalPath}}). Слой импортирует только нижележащие слои, слайсы — через index.ts',
        },
      ],
      'boundaries/no-unknown-files': 'error',
      'boundaries/no-unknown-dependencies': 'error',

      // TypeScript
      '@typescript-eslint/no-explicit-any': 'error',
      '@typescript-eslint/consistent-type-imports': ['error', { prefer: 'type-imports' }],
      '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],

      // Vue: только script setup + TS, строки интерфейса — через i18n
      'vue/block-lang': ['error', { script: { lang: 'ts' } }],
      'vue/component-api-style': ['error', ['script-setup']],
      'vue/no-bare-strings-in-template': 'error',
      'vue/multi-word-component-names': 'off',
      'vue/define-macros-order': [
        'error',
        { order: ['defineProps', 'defineEmits', 'defineModel'] },
      ],
    },
  },

  // Файлы конфигурации и тесты Node — без проверки типов проекта
  {
    files: ['*.config.{js,ts}', 'eslint.config.js'],
    ...tseslint.configs.disableTypeChecked,
  },

  prettier,
);
