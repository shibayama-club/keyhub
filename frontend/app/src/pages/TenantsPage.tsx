import { Header } from '../components/Header';
import { TenantList } from '../components/TenantList';
import { useNavigate } from 'react-router-dom';

export const TenantsPage = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-gray-50">
      <Header showBackButton backPath="/home" backLabel="ホーム" />

      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Header Section */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">マイテナント</h1>
              <p className="mt-2 text-sm text-gray-600">参加しているテナントの一覧です</p>
            </div>
            <button
              onClick={() => navigate('/join-tenant')}
              className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            >
              テナントに参加
            </button>
          </div>
        </div>

        {/* Tenant List */}
        <TenantList />
      </main>
    </div>
  );
};
