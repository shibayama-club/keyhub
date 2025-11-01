import { useNavigate } from 'react-router-dom';
import { Navbar } from '../components/Navbar';

export const TenantsPage = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <div className="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          {/* Header */}
          <div className="mb-8 flex items-center justify-between">
            <div>
              <h2 className="text-2xl font-bold text-gray-900">テナント</h2>
              <p className="mt-1 text-sm text-gray-600">組織のテナントを管理します。</p>
            </div>
            <button
              onClick={() => navigate('/tenants/create')}
              className="inline-flex items-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-none"
            >
              新しいテナントを作成
            </button>
          </div>

          {/* Tenants List */}
          <div className="overflow-hidden bg-white shadow sm:rounded-lg">
            <div className="px-4 py-5 sm:px-6">
              <h3 className="text-lg leading-6 font-medium text-gray-900">既存のテナント</h3>
            </div>
            <div className="border-t border-gray-200 px-4 py-5 sm:px-6">
              <p className="text-gray-600">テナントが見つかりません。最初のテナントを作成して始めましょう。</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
