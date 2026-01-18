import { apiRequest } from '@/lib/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export interface Account {
  id: string;
  name: string;
  type: string;
  subtype: string;
  balance: number;
  currency: string;
  institution_name?: string;
}

export function useAccounts() {
  return useQuery({
    queryKey: ['accounts'],
    queryFn: async () => {
      const data = await apiRequest<{ accounts: Account[] }>("/accounts");
      return data.accounts || [];
    },
  });
}

export function useNetWorth() {
  return useQuery({
    queryKey: ['net-worth'],
    queryFn: async () => {
      const data = await apiRequest<{ net_worth: number }>("/accounts/net-worth");
      return data.net_worth || 0;
    },
  });
}

export function useCreateAccount() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (account: Partial<Account>) => {
      return apiRequest("/accounts", {
        method: "POST",
        body: JSON.stringify(account),
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] });
      queryClient.invalidateQueries({ queryKey: ['net-worth'] });
    },
  });
}
