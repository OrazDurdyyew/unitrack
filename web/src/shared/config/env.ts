import { z } from 'zod';

// Настройки сборки. Проверяются при старте, чтобы ошибка конфигурации была видна сразу.
const schema = z.object({
  apiBaseUrl: z.string().min(1).default('/api/v1'),
});

export const env = schema.parse({
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL as string | undefined,
});
