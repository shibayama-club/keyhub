import { useParams, useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import * as Sentry from '@sentry/react';
import { Navbar } from '../components/Navbar';
import { AssignRoomForm } from '../components/AssignRoomForm';
import { useMutationAssignRoomToTenant, useQueryGetTenantById, queryClient } from '../libs/query';
import { timestampFromDate } from '@bufbuild/protobuf/wkt';
import { getTenantTypeLabel } from '../utils/tenant';

export const AssignRoomPage = () => {
  const { tenantId } = useParams<{ tenantId: string }>();
  const navigate = useNavigate();
  const { mutateAsync: assignRoom, isPending } = useMutationAssignRoomToTenant();
  const { data: tenantData, isLoading: isTenantLoading, isError: isTenantError } = useQueryGetTenantById(tenantId || '');

  if (!tenantId) {
    navigate('/tenants');
    return null;
  }

  const handleAssignRoom = async (data: { roomId: string; expiresAt?: Date }) => {
    try {
      await assignRoom({
        tenantId,
        roomId: data.roomId,
        expiresAt: data.expiresAt ? timestampFromDate(data.expiresAt) : undefined,
      });

      // TanStack Queryのキャッシュを無効化
      await queryClient.invalidateQueries();

      toast.success('Roomをテナントに割り当てました');
    } catch (error) {
      // Sentryでエラーをキャプチャ
      Sentry.captureException(error);
      toast.error('Roomの割り当てに失敗しました');
      console.error('Room assignment error:', error);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <button
            onClick={() => navigate('/tenants')}
            className="mb-4 inline-flex items-center text-sm text-gray-600 hover:text-gray-900"
          >
            <svg className="mr-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            テナント一覧に戻る
          </button>
          <h1 className="text-3xl font-bold text-gray-900">Roomの割り当て</h1>
          <p className="mt-2 text-sm text-gray-600">このテナントにRoomを割り当てます。</p>
        </div>

        <div className="rounded-lg border-2 border-gray-200 bg-white p-8 shadow-md">
          <div className="mb-6">
            <h2 className="text-lg font-semibold text-gray-900">テナント情報</h2>
            {isTenantLoading ? (
              <div className="mt-2 rounded-md bg-gray-50 p-4">
                <div className="flex items-center">
                  <div className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-solid border-indigo-600 border-r-transparent"></div>
                  <p className="ml-2 text-sm text-gray-600">読み込み中...</p>
                </div>
              </div>
            ) : isTenantError ? (
              <div className="mt-2 rounded-md bg-red-50 p-4">
                <p className="text-sm text-red-600">テナント情報の取得に失敗しました</p>
              </div>
            ) : tenantData?.tenant ? (
              <div className="mt-2 space-y-3 rounded-md bg-gray-50 p-4">
                <div className="flex items-center space-x-2">
                  <p className="text-base font-medium text-gray-900">{tenantData.tenant.name}</p>
                  <span className="inline-flex items-center rounded-full bg-indigo-100 px-2.5 py-0.5 text-xs font-medium text-indigo-800">
                    {getTenantTypeLabel(tenantData.tenant.tenantType)}
                  </span>
                </div>
                {tenantData.tenant.description && (
                  <p className="text-sm text-gray-600">{tenantData.tenant.description}</p>
                )}
                <p className="text-xs text-gray-500">
                  <span className="font-medium">ID:</span> {tenantId}
                </p>
              </div>
            ) : (
              <div className="mt-2 rounded-md bg-gray-50 p-4">
                <p className="text-sm text-gray-600">
                  <span className="font-medium">Tenant ID:</span> {tenantId}
                </p>
              </div>
            )}
          </div>

          <div className="border-t border-gray-200 pt-6">
            <h2 className="mb-4 text-lg font-semibold text-gray-900">Room割り当て</h2>
            <AssignRoomForm onSubmit={handleAssignRoom} isSubmitting={isPending} />
          </div>
        </div>

        <div className="mt-6 rounded-lg border border-blue-200 bg-blue-50 p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-blue-800">Room IDについて</h3>
              <div className="mt-2 text-sm text-blue-700">
                <p>
                  Room
                  IDは、Roomsページから作成したRoomのIDです。作成時にレスポンスで返されるUUID形式のIDを入力してください。
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
