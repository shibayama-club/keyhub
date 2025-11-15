import { useState } from 'react';
import { useQueryGetTenantByJoinCode } from '../lib/query';
import { TenantType } from '../../../gen/src/keyhub/app/v1/app_pb';

const TENANT_TYPE_LABELS: Record<TenantType, string> = {
  [TenantType.UNSPECIFIED]: '未指定',
  [TenantType.TEAM]: 'チーム',
  [TenantType.DEPARTMENT]: '部署',
  [TenantType.PROJECT]: 'プロジェクト',
  [TenantType.LABORATORY]: '研究室',
};

export function JoinTenantPage() {
  const [joinCode, setJoinCode] = useState('');
  const [searchCode, setSearchCode] = useState('');

  const { data, isLoading, error } = useQueryGetTenantByJoinCode(searchCode);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchCode(joinCode);
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <div className="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
        <h1 className="text-2xl font-bold text-center mb-8">テナントに参加</h1>

        <form onSubmit={handleSearch} className="space-y-4">
          <div>
            <label htmlFor="joinCode" className="block text-sm font-medium text-gray-700 mb-2">
              参加コード
            </label>
            <input
              id="joinCode"
              type="text"
              value={joinCode}
              onChange={(e) => setJoinCode(e.target.value)}
              placeholder="参加コードを入力してください"
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              required
            />
          </div>

          <button
            type="submit"
            disabled={isLoading || !joinCode}
            className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            {isLoading ? '検索中...' : 'テナント情報を表示'}
          </button>
        </form>

        {error && (
          <div className="mt-6 p-4 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-800 text-sm">
              {error.message || '参加コードが無効です'}
            </p>
          </div>
        )}

        {data && (
          <div className="mt-6 p-6 bg-blue-50 border border-blue-200 rounded-md space-y-4">
            <h2 className="text-lg font-semibold text-gray-900">テナント情報</h2>

            <div className="space-y-3">
              <div>
                <p className="text-xs text-gray-500">テナント名</p>
                <p className="text-base font-medium text-gray-900">{data.name}</p>
              </div>

              {data.description && (
                <div>
                  <p className="text-xs text-gray-500">説明</p>
                  <p className="text-base text-gray-700">{data.description}</p>
                </div>
              )}

              <div>
                <p className="text-xs text-gray-500">タイプ</p>
                <p className="text-base font-medium text-gray-900">
                  {TENANT_TYPE_LABELS[data.tenantType] || '未指定'}
                </p>
              </div>
            </div>

            <button
              type="button"
              className="w-full bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700 transition-colors mt-4"
            >
              このテナントに参加
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
