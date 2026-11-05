import { api } from '@/shared/api';

export interface HealthInfo {
  version: string;
}

// Состояние сервера. Ошибка сети или не-2xx ответ → исключение, чтобы TanStack Query перешёл в isError.
export async function getHealth(): Promise<HealthInfo> {
  const { data, response } = await api.GET('/health');
  if (data === undefined) {
    throw new Error(`Сервер недоступен: HTTP ${String(response.status)}`);
  }
  return { version: data.version };
}
