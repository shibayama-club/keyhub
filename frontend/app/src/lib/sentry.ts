import { useEffect } from 'react';
import * as Sentry from '@sentry/react';
import { useLocation, useNavigationType, createRoutesFromChildren, matchRoutes } from 'react-router-dom';
import { Code, ConnectError } from '@connectrpc/connect';
import { SENTRY_DSN, APP_ENVIRONMENT, IS_PRODUCTION } from './env';

export const initSentry = () => {
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
        Sentry.replayIntegration(),
      ],
      // パフォーマンス監視のサンプリング率
      // 本番: 50%、開発/ステージング: 100%
      tracesSampleRate: IS_PRODUCTION ? 0.5 : 1.0,
      // Replay設定
      replaysSessionSampleRate: 0.1,
      replaysOnErrorSampleRate: 1.0,
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
};

/**
 * Sentryが有効かどうか
 */
const isSentryEnabled = !!SENTRY_DSN;

/**
 * User コンテキストをSentryに設定
 *
 * ログイン時に呼び出す
 */
export function setSentryUser(userId: string, email: string, name: string) {
  if (!isSentryEnabled) return;

  Sentry.setUser({
    id: userId,
    email,
    username: name,
  });

  // ログイン時刻をコンテキストとして設定
  Sentry.setContext('user', {
    id: userId,
    email,
    name,
    loginTime: new Date().toISOString(),
  });
}

/**
 * User コンテキストをクリア
 * ログアウト時に呼び出す
 */
export function clearSentryUser() {
  if (!isSentryEnabled) return;

  Sentry.setUser(null);
  Sentry.setContext('user', null);
}

/**
 * スキップすべきエラータイプ
 * これらのエラーはSentryに送信しない
 */
const SKIP_ERROR_TYPES = ['ValidationError', 'PermissionError', 'AuthenticationError'];

/**
 * スキップすべきConnect RPCエラーコード
 * これらはビジネスロジックの一部であり、システムエラーではない
 */
const SKIP_RPC_CODES = [
  Code.Unauthenticated, // 認証されていない（ログインページへリダイレクト）
  Code.PermissionDenied, // 権限がない（正常なビジネスロジック）
  Code.InvalidArgument, // バリデーションエラー（ユーザー入力ミス）
  Code.NotFound, // リソースが見つからない（404相当）
  Code.AlreadyExists, // すでに存在する（409相当）
];

/**
 * エラーをSentryに送信するかどうかを判定
 */
function shouldCaptureError(error: unknown): boolean {
  // Sentryが無効な場合はスキップ
  if (!isSentryEnabled) return false;

  // Connect RPCエラーの場合、ビジネスエラーはスキップ
  if (error instanceof ConnectError) {
    if (SKIP_RPC_CODES.includes(error.code)) {
      return false;
    }
  }

  // エラーオブジェクトの名前をチェック
  if (error instanceof Error && SKIP_ERROR_TYPES.includes(error.name)) {
    return false;
  }

  return true;
}

/**
 * エラーをログに記録し、Sentryに送信
 *
 * @example
 * ```typescript
 * try {
 *   await fetchData();
 * } catch (error) {
 *   logError(error, { context: 'データ取得' });
 * }
 * ```
 */
export function logError(error: unknown, context?: { [key: string]: unknown }) {
  // コンソールには常に出力
  console.error('Error:', error, context);

  // Sentryへの送信判定
  if (!shouldCaptureError(error)) {
    return;
  }

  // Connect RPCエラーの場合
  if (error instanceof ConnectError) {
    const endpoint = context?.endpoint as string | undefined;
    const method = context?.method as string | undefined;

    Sentry.withScope((scope) => {
      // 同じエンドポイント・ステータスコードのエラーをグルーピング
      if (endpoint) {
        scope.setFingerprint(['{{ default }}', endpoint, String(error.code), method || 'unknown']);
      }

      // コンテキスト情報を追加
      scope.setContext('connectRPC', {
        code: error.code,
        message: error.message,
        endpoint,
        method,
        metadata: error.metadata,
      });

      // タグを設定
      scope.setTag('error_type', 'connect_rpc');
      scope.setTag('rpc_code', error.code);

      Sentry.captureException(error);
    });
    return;
  }

  // その他のエラー
  Sentry.withScope((scope) => {
    if (context) {
      scope.setContext('additional', context);
    }

    Sentry.captureException(error);
  });
}

/**
 * 手動でメッセージをSentryに送信
 *
 * @example
 * ```typescript
 * logMessage('想定外の状態が発生', 'warning', { userId: '123' });
 * ```
 */
export function logMessage(
  message: string,
  level: 'info' | 'warning' | 'error' = 'info',
  context?: { [key: string]: unknown },
) {
  if (!isSentryEnabled) {
    console.log(`[${level}]`, message, context);
    return;
  }

  Sentry.withScope((scope) => {
    scope.setLevel(level);

    if (context) {
      scope.setContext('message_context', context);
    }

    Sentry.captureMessage(message);
  });
}

// 後方互換性のためのエイリアス
export const setUser = setSentryUser;
export const clearUser = clearSentryUser;
