import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { useMutationLogout } from '../libs/query';
import { useAuthStore } from '../libs/auth';

export default function DashboardPage() {
  const navigate = useNavigate();
  const { mutate: logoutMutation } = useMutationLogout();
  const clearAuth = useAuthStore((state) => state.clearAuth);

  const handleLogout = () => {
    logoutMutation(
      {},
      {
        onSuccess: () => {
          clearAuth();
          toast.success('Logged out successfully');
          navigate('/login');
        },
        onError: (error) => {
          console.error('Logout error:', error);
          // エラーが発生してもローカルの認証情報はクリア
          clearAuth();
          toast.error('Error logging out');
          navigate('/login');
        },
      }
    );
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex h-16 justify-between">
            <div className="flex">
              <div className="flex flex-shrink-0 items-center">
                <h1 className="text-xl font-semibold">KeyHub Console</h1>
              </div>
              <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
                <a
                  href="#"
                  className="inline-flex items-center border-b-2 border-indigo-500 px-1 pt-1 text-sm font-medium text-gray-900"
                >
                  Dashboard
                </a>
                <a
                  href="#"
                  className="inline-flex items-center border-b-2 border-transparent px-1 pt-1 text-sm font-medium text-gray-500 hover:border-gray-300 hover:text-gray-700"
                >
                  Tenants
                </a>
                <a
                  href="#"
                  className="inline-flex items-center border-b-2 border-transparent px-1 pt-1 text-sm font-medium text-gray-500 hover:border-gray-300 hover:text-gray-700"
                >
                  Users
                </a>
              </div>
            </div>
            <div className="hidden sm:ml-6 sm:flex sm:items-center">
              <button
                onClick={handleLogout}
                className="inline-flex items-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-none"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <div className="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {/* Tenant Summary Card */}
            <div className="overflow-hidden rounded-lg bg-white shadow">
              <div className="px-4 py-5 sm:p-6">
                <dt className="truncate text-sm font-medium text-gray-500">Total Tenants</dt>
                <dd className="mt-1 text-3xl font-semibold text-gray-900">0</dd>
              </div>
            </div>

            {/* Users Summary Card */}
            <div className="overflow-hidden rounded-lg bg-white shadow">
              <div className="px-4 py-5 sm:p-6">
                <dt className="truncate text-sm font-medium text-gray-500">Total Users</dt>
                <dd className="mt-1 text-3xl font-semibold text-gray-900">0</dd>
              </div>
            </div>

            {/* Active Sessions Card */}
            <div className="overflow-hidden rounded-lg bg-white shadow">
              <div className="px-4 py-5 sm:p-6">
                <dt className="truncate text-sm font-medium text-gray-500">Active Sessions</dt>
                <dd className="mt-1 text-3xl font-semibold text-gray-900">0</dd>
              </div>
            </div>
          </div>

          {/* Main Content Area */}
          <div className="mt-8">
            <div className="overflow-hidden bg-white shadow sm:rounded-lg">
              <div className="px-4 py-5 sm:px-6">
                <h3 className="text-lg leading-6 font-medium text-gray-900">Console Dashboard</h3>
                <p className="mt-1 max-w-2xl text-sm text-gray-500">Manage your organization's tenants and users.</p>
              </div>
              <div className="border-t border-gray-200 px-4 py-5 sm:px-6">
                <p className="text-gray-600">
                  Welcome to KeyHub Console. Use the navigation above to manage tenants and users.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
