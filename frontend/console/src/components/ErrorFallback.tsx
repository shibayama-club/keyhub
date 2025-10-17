import { useEffect } from 'react';
import { useRouteError } from 'react-router-dom';
import * as Sentry from '@sentry/react';

export default function ErrorFallback() {
  const error = useRouteError();

  useEffect(() => {
    // Sentryにエラーを送信
    Sentry.captureException(error, {
      tags: {
        location: 'ErrorFallback',
        type: 'router_error',
      },
    });

    // コンソールにもログ出力
    console.error('Router error:', error);
  }, [error]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md rounded-lg bg-white p-8 shadow-lg">
        <h2 className="mb-4 text-2xl font-bold text-red-600">エラーが発生しました</h2>
        <p className="mb-6 text-gray-600">
          予期しないエラーが発生しました。ページを再読み込みしてください。
        </p>
        <button
          onClick={() => window.location.reload()}
          className="w-full rounded-md bg-indigo-600 px-4 py-2 text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
        >
          ページを再読み込み
        </button>
      </div>
    </div>
  );
}
