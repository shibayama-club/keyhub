import * as Sentry from '@sentry/react';
import { Code, ConnectError } from '@connectrpc/connect';
import { SENTRY_DSN } from './env';

/**
 * Sentryが有効かどうか
 */
const isSentryEnabled = !!SENTRY_DSN;

/**
 * Organization コンテキストをSentryに設定
 * Console認証ではユーザーではなく組織単位でログインするため、
 * 組織情報をタグとして設定する
 *
 * ログイン時に呼び出す
 */
export function setSentryOrganization(organizationId: string) {
  if (!isSentryEnabled) return;

  // 組織IDをタグとして設定（検索・フィルタリング用）
  Sentry.setTag('organization_id', organizationId);

  // 組織IDをコンテキストとしても設定（詳細情報用）
  Sentry.setContext('organization', {
    id: organizationId,
    loginTime: new Date().toISOString(),
  });
}

/**
 * Organization コンテキストをクリア
 * ログアウト時に呼び出す
 */
export function clearSentryOrganization() {
  if (!isSentryEnabled) return;

  Sentry.setTag('organization_id', null);
  Sentry.setContext('organization', null);
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
