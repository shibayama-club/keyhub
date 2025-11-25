import type { Key } from '../../../gen/src/keyhub/console/v1/common_pb';
import { getKeyStatusLabel, getKeyStatusBadgeColor } from '../utils/key';

export const KeyList = ({
  keys,
  isLoading,
  isError,
}: {
  keys: Key[];
  isLoading: boolean;
  isError: boolean;
}) => {
  if (isLoading) {
    return (
      <div className="px-4 py-12 text-center">
        <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-indigo-600 border-r-transparent"></div>
        <p className="mt-2 text-sm text-gray-600">読み込み中...</p>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="px-4 py-5 sm:px-6">
        <p className="text-red-600">鍵の取得に失敗しました。</p>
      </div>
    );
  }

  if (keys.length === 0) {
    return (
      <div className="px-4 py-5 sm:px-6">
        <p className="text-gray-600">鍵が見つかりません。最初の鍵を作成して始めましょう。</p>
      </div>
    );
  }

  return (
    <ul className="divide-y divide-gray-200">
      {keys.map((key) => (
        <li key={key.id} className="px-4 py-5 hover:bg-gray-50 sm:px-6">
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <div className="flex items-center space-x-3">
                <h4 className="text-base font-medium text-gray-900">{key.keyNumber}</h4>
                <span className={'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ' + getKeyStatusBadgeColor(key.status)}>
                  {getKeyStatusLabel(key.status)}
                </span>
              </div>
            </div>
          </div>
        </li>
      ))}
    </ul>
  );
};
