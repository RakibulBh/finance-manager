import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

interface User {
  id: string;
  email: string;
  family_id: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  setAuth: (user: User, token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      setAuth: (user, token) => {
        set({ user, token, isAuthenticated: true });
        // Also keep in localStorage for immediate sync (lib/api.ts uses it)
        localStorage.setItem("token", token);
        localStorage.setItem("user", JSON.stringify(user));
      },
      logout: () => {
        set({ user: null, token: null, isAuthenticated: false });
        localStorage.removeItem("token");
        localStorage.removeItem("user");
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage),
    }
  )
);
