import { create } from 'zustand';
import { setSentryOrganization, clearSentryOrganization } from './sentry';

// Storage keys constants
const STORAGE_TOKEN_KEY = 'console_token';
const STORAGE_EXPIRES_KEY = 'console_expires_at';

interface AuthState {
  isAuthenticated: boolean;
  token: string | null;
  expiresAt: number | null;
  setAuthData: (token: string, expiresInSeconds: number, organizationId: string) => void;
  checkAuth: () => void;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  isAuthenticated: false,
  token: null,
  expiresAt: null,

  setAuthData: (token: string, expiresInSeconds: number, organizationId: string) => {
    const expiresAt = Date.now() + expiresInSeconds * 1000;

    localStorage.setItem(STORAGE_TOKEN_KEY, token);
    localStorage.setItem(STORAGE_EXPIRES_KEY, expiresAt.toString());

    // Sentryに組織コンテキストを設定
    setSentryOrganization(organizationId);

    set({
      isAuthenticated: true,
      token,
      expiresAt,
    });
  },

  checkAuth: () => {
    const token = localStorage.getItem(STORAGE_TOKEN_KEY);
    const expiresAt = localStorage.getItem(STORAGE_EXPIRES_KEY);

    if (!token || !expiresAt) {
      get().clearAuth();
      return;
    }

    const expiresAtNum = parseInt(expiresAt, 10);
    if (isNaN(expiresAtNum) || Date.now() >= expiresAtNum) {
      get().clearAuth();
      return;
    }

    set({
      isAuthenticated: true,
      token,
      expiresAt: expiresAtNum,
    });
  },

  clearAuth: () => {
    localStorage.removeItem(STORAGE_TOKEN_KEY);
    localStorage.removeItem(STORAGE_EXPIRES_KEY);

    // Sentryの組織コンテキストをクリア
    clearSentryOrganization();

    set({
      isAuthenticated: false,
      token: null,
      expiresAt: null,
    });
  },
}));
