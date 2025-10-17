import { StrictMode, useEffect } from 'react';
import { createRoot } from 'react-dom/client';
import * as Sentry from '@sentry/react';
import { useLocation, useNavigationType, createRoutesFromChildren, matchRoutes } from 'react-router-dom';
import './index.css';
import App from './App.tsx';
import { SENTRY_DSN, APP_ENVIRONMENT, IS_PRODUCTION } from './libs/env';

// Sentryの初期化
if (SENTRY_DSN) {
  Sentry.init({
    dsn: SENTRY_DSN,
    environment: APP_ENVIRONMENT,
    normalizeDepth: 10,
    integrations: [
      Sentry.reactRouterV7BrowserTracingIntegration({
        useEffect,
        useLocation,
        useNavigationType,
        createRoutesFromChildren,
        matchRoutes,
      }),
    ],
    // パフォーマンス監視のサンプリング率
    // 本番: 50%、開発/ステージング: 100%
    tracesSampleRate: IS_PRODUCTION ? 0.5 : 1.0,
    // 無視するエラーパターン
    ignoreErrors: [
      // ネットワークエラー（Sentryで追跡する必要なし）
      'Network Error',
      'NetworkError',
      'Failed to fetch',
      // キャンセルされたリクエスト
      'AbortError',
      'Request aborted',
      // React Suspenseの内部エラー（無害）
      'Minified React error',
      // 一般的なブラウザ拡張機能のエラー
      'Non-Error promise rejection captured',
    ],
  });
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
