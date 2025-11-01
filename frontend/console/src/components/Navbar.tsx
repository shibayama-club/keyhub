import { useNavigate, useLocation } from 'react-router-dom';
import toast from 'react-hot-toast';
import { useMutationLogout } from '../libs/query';
import { useAuthStore } from '../libs/auth';

export const Navbar = () => {
  const navigate = useNavigate();
  const location = useLocation();
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
      },
    );
  };

  const isActive = (path: string) => location.pathname === path;

  return (
    <nav className="bg-white shadow">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 justify-between">
          <div className="flex">
            <div className="flex flex-shrink-0 items-center">
              <h1 className="text-xl font-semibold">KeyHub Console</h1>
            </div>
            <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
              <a
                href="/"
                onClick={(e) => {
                  e.preventDefault();
                  navigate('/');
                }}
                className={`inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium ${
                  isActive('/')
                    ? 'border-indigo-500 text-gray-900'
                    : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                }`}
              >
                Dashboard
              </a>
              <a
                href="/tenants"
                onClick={(e) => {
                  e.preventDefault();
                  navigate('/tenants');
                }}
                className={`inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium ${
                  isActive('/tenants')
                    ? 'border-indigo-500 text-gray-900'
                    : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                }`}
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
  );
};
