import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import * as Sentry from '@sentry/react';
import { Code, ConnectError } from '@connectrpc/connect';
import { useMutationLoginWithOrgId } from '../libs/query';
import { useAuthStore } from '../libs/auth';

// UUID validation regex
const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

export const LoginPage = () => {
  const navigate = useNavigate();
  const { mutate: loginMutation, isPending } = useMutationLoginWithOrgId();
  const setAuthData = useAuthStore((state) => state.setAuthData);

  const [organizationId, setOrganizationId] = useState('');
  const [organizationKey, setOrganizationKey] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validation
    if (!organizationId.trim() || !organizationKey.trim()) {
      toast.error('組織IDとキーを入力してください');
      return;
    }

    if (!UUID_REGEX.test(organizationId.trim())) {
      toast.error('組織IDは有効なUUID形式である必要があります');
      return;
    }

    const orgId = organizationId.trim();
    const orgKey = organizationKey.trim();

    loginMutation(
      {
        organizationId: orgId,
        organizationKey: orgKey,
      },
      {
        onSuccess: (data) => {
          // 認証データを保存
          setAuthData(data.sessionToken, Number(data.expiresIn), orgId);
          toast.success('ログインしました');
          navigate('/dashboard');
        },
        onError: (error) => {
          // Handle specific error types
          Sentry.captureException(error);
          if (error instanceof ConnectError) {
            if (error.code === Code.Unauthenticated) {
              toast.error('組織IDまたはキーが無効です');
            } else {
              toast.error(`ログインに失敗しました: ${error.message}`);
            }
          } else {
            toast.error('予期しないエラーが発生しました');
          }
          console.error('Login error:', error);
        },
      },
    );
  };

  // For development, show default credentials hint
  const fillDefaultCredentials = () => {
    setOrganizationId('550e8400-e29b-41d4-a716-446655440000');
    setOrganizationKey('org_key_example_12345');
  };

  return (
    <div className="flex min-h-screen flex-col justify-center bg-gray-50 py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">KeyHub Console</h2>
        <p className="mt-2 text-center text-sm text-gray-600">Sign in with your Organization credentials</p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white px-4 py-8 shadow sm:rounded-lg sm:px-10">
          <form className="space-y-6" onSubmit={handleSubmit}>
            <div>
              <label htmlFor="organizationId" className="block text-sm font-medium text-gray-700">
                Organization ID
              </label>
              <div className="mt-1">
                <input
                  id="organizationId"
                  name="organizationId"
                  type="text"
                  autoComplete="off"
                  required
                  value={organizationId}
                  onChange={(e) => setOrganizationId(e.target.value)}
                  className="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none sm:text-sm"
                  placeholder="Enter Organization ID (UUID format)"
                />
              </div>
            </div>

            <div>
              <label htmlFor="organizationKey" className="block text-sm font-medium text-gray-700">
                Organization Key
              </label>
              <div className="mt-1">
                <input
                  id="organizationKey"
                  name="organizationKey"
                  type="password"
                  autoComplete="off"
                  required
                  value={organizationKey}
                  onChange={(e) => setOrganizationKey(e.target.value)}
                  className="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none sm:text-sm"
                  placeholder="Enter Organization Key"
                />
              </div>
            </div>

            <div>
              <button
                type="submit"
                disabled={isPending}
                className="flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
              >
                {isPending ? 'Signing in...' : 'Sign in'}
              </button>
            </div>
          </form>

          {/* Development helper - remove in production */}
          {import.meta.env.DEV && (
            <div className="mt-6 border-t border-gray-200 pt-6">
              <button onClick={fillDefaultCredentials} className="w-full text-sm text-gray-500 hover:text-gray-700">
                Use default development credentials
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
