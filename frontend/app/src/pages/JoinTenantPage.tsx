import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import * as Sentry from '@sentry/react';
import { Header } from '../components/Header';
import { useQueryGetTenantByJoinCode, useMutationJoinTenant, queryClient } from '../lib/query';
import { getMyTenants } from '../../../gen/src/keyhub/app/v1/app-TenantService_connectquery';
import { TENANT_TYPE_LABELS } from '../lib/constants/tenant';

export function JoinTenantPage() {
  const [joinCode, setJoinCode] = useState('');
  const [searchCode, setSearchCode] = useState('');
  const navigate = useNavigate();

  const { data, isLoading, error } = useQueryGetTenantByJoinCode(searchCode);
  const joinMutation = useMutationJoinTenant();

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setSearchCode(joinCode);
    } catch (error) {
      Sentry.captureException(error);
      toast.error('テナント情報の取得に失敗しました');
    }
  };

  const handleJoin = async () => {
    try {
      await joinMutation.mutateAsync({ joinCode: searchCode });
      await queryClient.invalidateQueries({ queryKey: [getMyTenants] });
      toast.success('テナントに参加しました');
      navigate('/tenants');
    } catch (error) {
      Sentry.captureException(error);
      toast.error('テナントへの参加に失敗しました');
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header showBackButton backPath="/home" backLabel="ホーム" />
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mx-auto w-full max-w-md rounded-lg bg-white p-8 shadow-lg">
          <h1 className="mb-8 text-center text-2xl font-bold">テナントに参加</h1>

          <form onSubmit={handleSearch} className="space-y-4">
            <div>
              <label htmlFor="joinCode" className="mb-2 block text-sm font-medium text-gray-700">
                参加コード
              </label>
              <input
                id="joinCode"
                type="text"
                value={joinCode}
                onChange={(e) => setJoinCode(e.target.value)}
                placeholder="参加コードを入力してください"
                className="w-full rounded-md border border-gray-300 px-4 py-2 focus:border-blue-500 focus:ring-2 focus:ring-blue-500"
                required
              />
            </div>

            <button
              type="submit"
              disabled={isLoading || !joinCode}
              className="w-full rounded-md bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-gray-400"
            >
              {isLoading ? '検索中...' : 'テナント情報を表示'}
            </button>
          </form>

          {error && (
            <div className="mt-6 rounded-md border border-red-200 bg-red-50 p-4">
              <p className="text-sm font-semibold text-red-800">エラー</p>
              <p className="mt-1 text-sm text-red-700">
                {error.message ||
                  'テナント情報の取得に失敗しました。参加コードが正しいか、有効期限が切れていないかご確認ください。'}
              </p>
            </div>
          )}

          {data && (
            <div className="mt-6 space-y-4 rounded-md border border-blue-200 bg-blue-50 p-6">
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
                onClick={handleJoin}
                disabled={joinMutation.isPending}
                className="mt-4 w-full rounded-md bg-green-600 px-4 py-2 text-white transition-colors hover:bg-green-700 disabled:cursor-not-allowed disabled:bg-gray-400"
              >
                {joinMutation.isPending ? '参加中...' : 'このテナントに参加'}
              </button>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
