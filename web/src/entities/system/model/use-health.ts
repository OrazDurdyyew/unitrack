import { useQuery } from '@tanstack/vue-query';

import { getHealth } from '../api/health';

export const healthQueryKey = ['system', 'health'] as const;

export function useHealthQuery() {
  return useQuery({
    queryKey: healthQueryKey,
    queryFn: getHealth,
    refetchInterval: 60_000,
  });
}
