import { createClient } from '@connectrpc/connect';
import { create } from 'zustand';
import { AuthService } from '../../../gen/src/keyhub/app/v1/app_pb';
import { transport } from './connect';

const STORAGE_USER_KEY = 'app_user';

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
  checkAuth: () => Promise<void>;
  clearAuth: () => void;
}

const authClient = createClient(AuthService, transport);

export const useAuthStore = create<AuthState>((set, get) => ({
  isAuthenticated: false,
  user: null,
  isLoading: true,

  setUser: (user: User) => {
    localStorage.setItem(STORAGE_USER_KEY, JSON.stringify(user));

    set({
      isAuthenticated: true,
      user,
      isLoading: false,
    });
  },

  // アプリ起動時の認証確認。localStorageを信用せず、必ずサーバーの
  // セッション（Cookie）をGetMeで検証する。検証に通った場合のみ認証済みとし、
  // 失敗した場合はlocalStorageごとクリアしてログイン画面へ送る。
  checkAuth: async () => {
    set({ isLoading: true });

    try {
      const { user } = await authClient.getMe({});

      if (!user) {
        get().clearAuth();
        return;
      }

      get().setUser({
        id: user.id,
        email: user.email,
        name: user.name,
        picture: user.icon,
      });
    } catch {
      get().clearAuth();
    }
  },

  clearAuth: () => {
    localStorage.removeItem(STORAGE_USER_KEY);

    set({
      isAuthenticated: false,
      user: null,
      isLoading: false,
    });
  },
}));
