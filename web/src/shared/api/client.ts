import createClient from 'openapi-fetch';

import { env } from '@/shared/config';

import type { paths } from './schema';

// Единственный HTTP-клиент приложения. Типы — из api/openapi.yaml (schema.d.ts генерируется `task gen`).
export const api = createClient<paths>({ baseUrl: env.apiBaseUrl });
