import { createI18n } from 'vue-i18n';

import { defaultLocale, messages, type MessageSchema } from '@/shared/i18n';

export const i18n = createI18n<[MessageSchema], 'ru'>({
  locale: defaultLocale,
  fallbackLocale: defaultLocale,
  messages,
});
