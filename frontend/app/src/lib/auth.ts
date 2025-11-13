import { create } from 'zustand';

const STORAGE_USER_KEY = 'app_user';
const STORAGE_CHECKED_KEY = 'app_auth_checked';

interface User {
  id: string;
  email: string;
  name: string;
  picture?: string;
}

interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  isLoading: boolean;
  setUser: (user: User) => void;
  checkAuth: () => void;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  isAuthenticated: false,
  user: null,
  isLoading: true,

  setUser: (user: User) => {
    localStorage.setItem(STORAGE_USER_KEY, JSON.stringify(user));
    localStorage.setItem(STORAGE_CHECKED_KEY, 'true');

    set({
      isAuthenticated: true,
      user,
      isLoading: false,
    });
  },

  checkAuth: () => {
    const userStr = localStorage.getItem(STORAGE_USER_KEY);
    const checked = localStorage.getItem(STORAGE_CHECKED_KEY);

    if (!userStr || !checked) {
      set({
        isAuthenticated: false,
        user: null,
        isLoading: false,
      });
      return;
    }

    try {
      const user = JSON.parse(userStr) as User;
      set({
        isAuthenticated: true,
        user,
        isLoading: false,
      });
    } catch (error) {
      get().clearAuth();
    }
  },

  clearAuth: () => {
    localStorage.removeItem(STORAGE_USER_KEY);
    localStorage.removeItem(STORAGE_CHECKED_KEY);

    set({
      isAuthenticated: false,
      user: null,
      isLoading: false,
    });
  },
}));
