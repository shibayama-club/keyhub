import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { useAuthStore } from '../lib/auth';
import { useQueryGetMe } from '../lib/query';

export const CallbackPage = () => {
  const navigate = useNavigate();
  const setUser = useAuthStore((state) => state.setUser);
  const [isProcessing, setIsProcessing] = useState(true);
  const hasRun = useRef(false);

  const { refetch: fetchMe } = useQueryGetMe();

  useEffect(() => {
    if (hasRun.current) return;
    hasRun.current = true;

    const handleCallback = async () => {
      await new Promise((resolve) => setTimeout(resolve, 100));

      try {
        const { data, error: fetchError } = await fetchMe();

        if (fetchError) {
          toast.error('認証に失敗しました。もう一度お試しください。');
          navigate('/login', { replace: true });
          return;
        }

        if (data?.user) {
          setUser({
            id: data.user.id,
            email: data.user.email,
            name: data.user.name,
            picture: data.user.icon,
          });

          toast.success('ログインしました！');
          navigate('/home', { replace: true });
        } else {
          toast.error('ユーザー情報が見つかりませんでした');
          navigate('/login', { replace: true });
        }
      } catch (error) {
        toast.error('ログイン中にエラーが発生しました');
        navigate('/login', { replace: true });
      } finally {
        setIsProcessing(false);
      }
    };

    handleCallback();
  }, []);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gray-50">
      <div className="text-center">
        {isProcessing ? (
          <>
            <div className="mb-4">
              <div className="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-indigo-200 border-t-indigo-600"></div>
            </div>
            <h2 className="text-xl font-semibold text-gray-900">Signing you in...</h2>
            <p className="mt-2 text-sm text-gray-600">Please wait a moment</p>
          </>
        ) : (
          <>
            <h2 className="text-xl font-semibold text-gray-900">Processing...</h2>
          </>
        )}
      </div>
    </div>
  );
};
