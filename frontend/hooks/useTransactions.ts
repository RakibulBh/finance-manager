import { apiRequest } from '@/lib/api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export interface Transaction {
  id: string;
  date: string;
  name: string;
  merchant_name?: string;
  category_name?: string;
  account_name: string;
  amount: number;
  account_id: string;
  category_id?: string;
}

export function useTransactions() {
  return useQuery({
    queryKey: ['transactions'],
    queryFn: async () => {
      const data = await apiRequest<Transaction[]>("/transactions");
      return data || [];
    },
  });
}

export function useCreateTransaction() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (transaction: Partial<Transaction>) => {
      return apiRequest("/transactions", {
        method: "POST",
        body: JSON.stringify(transaction),
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] });
      queryClient.invalidateQueries({ queryKey: ['accounts'] });
      queryClient.invalidateQueries({ queryKey: ['net-worth'] });
    },
  });
}
