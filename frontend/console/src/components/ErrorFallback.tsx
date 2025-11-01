import { useEffect } from 'react';
import { useRouteError } from 'react-router-dom';
import { logError } from '../libs/sentry';

export const ErrorFallback = () => {
  const error = useRouteError();

  useEffect(() => {
    // Sentryにエラーを送信（fingerprintingとコンテキスト付き）
    logError(error, {
      location: 'ErrorFallback',
      type: 'router_error',
      path: window.location.pathname,
    });
  }, [error]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md rounded-lg bg-white p-8 shadow-lg">
        <h2 className="mb-4 text-2xl font-bold text-red-600">エラーが発生しました</h2>
        <p className="mb-6 text-gray-600">予期しないエラーが発生しました。ページを再読み込みしてください。</p>
        <button
          onClick={() => window.location.reload()}
          className="w-full rounded-md bg-indigo-600 px-4 py-2 text-white hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-none"
        >
          ページを再読み込み
        </button>
      </div>
    </div>
  );
};
