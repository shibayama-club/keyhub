// 環境変数の型定義と取得

/**
 * Sentry DSN
 * エラートラッキングのためのSentry設定
 */
export const SENTRY_DSN = import.meta.env.VITE_SENTRY_DSN as string | undefined;

/**
 * 環境名
 * development, staging, production など
 */
export const APP_ENVIRONMENT = import.meta.env.MODE as string;

/**
 * 本番環境かどうか
 */
export const IS_PRODUCTION = APP_ENVIRONMENT === 'production';

/**
 * API Base URL
 * バックエンドAPIのベースURL
 */
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL as string | undefined;
