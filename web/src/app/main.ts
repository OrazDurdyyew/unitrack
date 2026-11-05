import '@fontsource-variable/inter';
import '@/shared/ui/tokens.css';
import './styles/reset.css';
import './styles/base.css';

import { createPinia } from 'pinia';
import { createApp } from 'vue';

import App from './App.vue';
import { i18n } from './providers/i18n';
import { queryOptions, queryPlugin } from './providers/query';
import { router } from './router';

createApp(App)
  .use(createPinia())
  .use(i18n)
  .use(router)
  .use(queryPlugin, queryOptions)
  .mount('#app');
