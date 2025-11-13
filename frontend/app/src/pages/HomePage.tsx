import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { useAuthStore } from '../lib/auth';
import { useMutationLogout } from '../lib/query';

export const HomePage = () => {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const clearAuth = useAuthStore((state) => state.clearAuth);

  const { mutate: logout, isPending: isLoggingOut } = useMutationLogout();

  const handleLogout = () => {
    logout(
      {},
      {
        onSuccess: () => {
          clearAuth();
          toast.success('Successfully signed out');
          navigate('/login', { replace: true });
        },
        onError: (error) => {
          console.error('Logout error:', error);
          clearAuth();
          toast.error('Signed out (with errors)');
          navigate('/login', { replace: true });
        },
      },
    );
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold text-gray-900">KeyHub App</h1>
            <button
              onClick={handleLogout}
              disabled={isLoggingOut}
              className="rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-gray-300 transition-colors ring-inset hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {isLoggingOut ? 'Signing out...' : 'Sign out'}
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Welcome Section */}
        <div className="mb-8 rounded-lg bg-white p-6 shadow">
          <div className="flex items-center gap-4">
            {user?.picture && <img src={user.picture} alt={user.name} className="h-16 w-16 rounded-full" />}
            <div>
              <h2 className="text-xl font-semibold text-gray-900">Welcome, {user?.name || 'User'}!</h2>
              <p className="text-sm text-gray-600">{user?.email}</p>
            </div>
          </div>
        </div>

        {/* Dashboard Cards */}
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div className="rounded-lg bg-white p-6 shadow">
            <h3 className="text-lg font-semibold text-gray-900">Profile</h3>
            <p className="mt-2 text-sm text-gray-600">View and edit your profile information</p>
            <button className="mt-4 text-sm font-semibold text-indigo-600 hover:text-indigo-500">View Profile →</button>
          </div>

          <div className="rounded-lg bg-white p-6 shadow">
            <h3 className="text-lg font-semibold text-gray-900">Settings</h3>
            <p className="mt-2 text-sm text-gray-600">Manage your account settings and preferences</p>
            <button className="mt-4 text-sm font-semibold text-indigo-600 hover:text-indigo-500">
              Go to Settings →
            </button>
          </div>

          <div className="rounded-lg bg-white p-6 shadow">
            <h3 className="text-lg font-semibold text-gray-900">Help</h3>
            <p className="mt-2 text-sm text-gray-600">Get help and support for your account</p>
            <button className="mt-4 text-sm font-semibold text-indigo-600 hover:text-indigo-500">Get Help →</button>
          </div>
        </div>

        {/* Info Section */}
        <div className="mt-8 rounded-lg bg-indigo-50 p-6">
          <h3 className="text-lg font-semibold text-indigo-900">Welcome to KeyHub App</h3>
          <p className="mt-2 text-sm text-indigo-700">
            You have successfully authenticated with Google OAuth. This is your home dashboard where you can manage your
            account and access various features.
          </p>
        </div>
      </main>
    </div>
  );
};
