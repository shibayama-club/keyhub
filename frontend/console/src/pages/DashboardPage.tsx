import { Navbar } from '../components/Navbar';

export const DashboardPage = () => {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

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
};
