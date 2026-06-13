import { create } from 'zustand';

interface AuthState {
  token: string | null;
  refreshToken: string | null;
  user: any | null;
  setToken: (token: string, refreshToken?: string) => void;
  setUser: (user: any) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('token'),
  refreshToken: localStorage.getItem('refresh_token'),
  user: null,
  setToken: (token, refreshToken) => {
    localStorage.setItem('token', token);
    if (refreshToken) {
      localStorage.setItem('refresh_token', refreshToken);
      set({ token, refreshToken });
    } else {
      set({ token });
    }
  },
  setUser: (user) => set({ user }),
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('refresh_token');
    set({ token: null, refreshToken: null, user: null });
  },
}));
