import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { useAuthStore } from '../lib/auth';
import { useMutationLogout } from '../lib/query';

interface HeaderProps {
  showBackButton?: boolean;
  backPath?: string;
  backLabel?: string;
}

export const Header = ({ showBackButton = false, backPath = '/home', backLabel = 'ホーム' }: HeaderProps) => {
  const navigate = useNavigate();
  const clearAuth = useAuthStore((state) => state.clearAuth);
  const { mutate: logout, isPending: isLoggingOut } = useMutationLogout();

  const handleLogout = () => {
    logout(
      {},
      {
        onSuccess: () => {
          clearAuth();
          toast.success('ログアウトしました');
          navigate('/login', { replace: true });
        },
        onError: () => {
          clearAuth();
          toast.error('ログアウトしました（エラーあり）');
          navigate('/login', { replace: true });
        },
      },
    );
  };

  const handleBack = () => {
    navigate(backPath);
  };

  return (
    <header className="bg-white shadow">
      <div className="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <h1 className="text-2xl font-bold text-gray-900">KeyHub App</h1>
            {showBackButton && (
              <button
                onClick={handleBack}
                className="flex items-center gap-1 text-sm text-gray-600 hover:text-gray-900"
              >
                ← {backLabel}
              </button>
            )}
          </div>
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
  );
};
