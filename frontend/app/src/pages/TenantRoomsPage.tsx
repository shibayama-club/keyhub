import { useNavigate, useParams } from 'react-router-dom';
import { Header } from '../components/Header';
import { RoomCard } from '../components/RoomCard';
import { useQueryGetRoomsByTenant } from '../lib/query';

export const TenantRoomsPage = () => {
  const navigate = useNavigate();
  const { tenantId } = useParams<{ tenantId: string }>();
  const { data, isLoading, error } = useQueryGetRoomsByTenant(tenantId || '');

  if (!tenantId) {
    navigate('/tenants');
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header showBackButton backPath="/tenants" backLabel="テナント一覧" />

      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">部屋一覧</h1>
          <p className="mt-2 text-sm text-gray-600">このテナントに割り当てられた部屋の一覧です</p>
        </div>

        {isLoading && (
          <div className="flex items-center justify-center py-12">
            <div className="h-8 w-8 animate-spin rounded-full border-4 border-indigo-600 border-t-transparent"></div>
          </div>
        )}

        {error && (
          <div className="rounded-lg border border-red-200 bg-red-50 p-4">
            <div className="flex">
              <div className="flex-shrink-0">
                <svg className="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              </div>
              <div className="ml-3">
                <h3 className="text-sm font-medium text-red-800">エラーが発生しました</h3>
                <div className="mt-2 text-sm text-red-700">
                  <p>部屋情報の取得に失敗しました。時間をおいて再度お試しください。</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {!isLoading && !error && (!data?.rooms || data.rooms.length === 0) && (
          <div className="rounded-lg border-2 border-dashed border-gray-300 p-12 text-center">
            <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
              />
            </svg>
            <h3 className="mt-2 text-sm font-semibold text-gray-900">割り当てられた部屋がありません</h3>
            <p className="mt-1 text-sm text-gray-500">管理者に部屋の割り当てを依頼してください</p>
          </div>
        )}

        {!isLoading && !error && data?.rooms && data.rooms.length > 0 && (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {data.rooms.map((room) => (
              <RoomCard key={room.id} room={room} />
            ))}
          </div>
        )}
      </main>
    </div>
  );
};
