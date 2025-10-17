import { createConnectTransport } from '@connectrpc/connect-web';
import { createClient } from '@connectrpc/connect';
import { ConsoleAuthService } from '../../../gen/src/keyhub/console/v1/console_pb';

const STORAGE_TOKEN_KEY = 'console_token';

const getBaseUrl = (): string => {
  // 本番環境時にセットが必須
  if (import.meta.env.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL;
  }
  // 開発環境の場合、Viteプロキシが使われる
  if (import.meta.env.DEV) {
    return '';
  }

  return 'http://localhost:8081';
};

export const transport = createConnectTransport({
  baseUrl: getBaseUrl(),
  interceptors: [
    // 認証インターセプター
    (next) => async (req) => {
      const token = localStorage.getItem(STORAGE_TOKEN_KEY);
      if (token) {
        req.header.set('Authorization', `Bearer ${token}`);
      }
      return await next(req);
    },
    (next) => async (req) => {
      try {
        return await next(req);
      } catch (error) {
        console.error('Connect RPC request failed:', {
          procedure: req.url,
          error,
        });
        throw error;
      }
    },
  ],
});

export const consoleAuthClient = createClient(ConsoleAuthService, transport);
