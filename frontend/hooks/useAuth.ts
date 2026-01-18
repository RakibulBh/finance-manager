import { apiRequest } from '@/lib/api';
import { useAuthStore } from '@/store/useAuthStore';
import { useMutation } from '@tanstack/react-query';

interface User {
  id: string;
  email: string;
  family_id: string;
}

interface LoginResponse {
  token: string;
  user: User;
}

interface LoginCredentials {
  email: string;
  password: string;
}

interface RegisterData {
  email: string;
  password: string;
  family_name: string;
}

export function useLogin() {
  const setAuth = useAuthStore((state) => state.setAuth);

  return useMutation({
    mutationFn: async (credentials: LoginCredentials) => {
      const data = await apiRequest<LoginResponse>("/auth/login", {
        method: "POST",
        body: JSON.stringify(credentials),
      });
      return data;
    },
    onSuccess: (data) => {
      setAuth(data.user, data.token);
    },
  });
}

export function useRegister() {
  const setAuth = useAuthStore((state) => state.setAuth);

  return useMutation({
    mutationFn: async (formData: RegisterData) => {
      const data = await apiRequest<LoginResponse>("/auth/register", {
        method: "POST",
        body: JSON.stringify(formData),
      });
      return data;
    },
    onSuccess: (data) => {
      setAuth(data.user, data.token);
    },
  });
}
