/// <reference types="vite/client" />

// Типы для импорта .vue-файлов из .ts (vue-tsc выводит точные типы сам; это для ESLint).
declare module '*.vue' {
  import type { DefineComponent } from 'vue';

  const component: DefineComponent<object, object, unknown>;
  export default component;
}
