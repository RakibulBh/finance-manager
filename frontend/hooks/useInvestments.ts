import { apiRequest } from '@/lib/api';
import { useQuery } from '@tanstack/react-query';
import { Account } from './useAccounts';

export function useInvestments() {
  return useQuery({
    queryKey: ['investments'],
    queryFn: async () => {
      const data = await apiRequest<{ accounts: Account[] }>("/accounts");
      return data.accounts?.filter(a => a.type === 'investment') || [];
    },
  });
}
