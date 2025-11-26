import { useNavigate, useParams } from 'react-router-dom';
import toast from 'react-hot-toast';
import * as Sentry from '@sentry/react';
import { Navbar } from '../components/Navbar';
import { CreateKeyForm } from '../components/CreateKeyForm';
import { useMutationCreateKey, queryClient } from '../libs/query';

export const CreateKeyPage = () => {
  const navigate = useNavigate();
  const { roomId } = useParams<{ roomId: string }>();
  const { mutateAsync: createKey, isPending } = useMutationCreateKey();

  if (!roomId) {
    navigate('/rooms');
    return null;
  }

  const handleSubmit = async (data: { keyNumber: string }) => {
    try {
      await createKey({
        roomId: roomId,
        keyNumber: data.keyNumber,
      });

      // TanStack Queryのキャッシュを無効化
      await queryClient.invalidateQueries();

      toast.success('鍵を作成しました');
      // 成功後はRoom一覧に戻る
      navigate('/rooms', { replace: true });
    } catch (error) {
      // Sentryでエラーをキャプチャ
      Sentry.captureException(error);
      toast.error('鍵の作成に失敗しました');
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">新しい鍵を作成</h1>
          <p className="mt-2 text-sm text-gray-600">このRoomの新しい鍵を作成します。鍵番号を入力してください。</p>
        </div>

        <div className="rounded-lg border-2 border-gray-200 bg-white p-8 shadow-md">
          <CreateKeyForm onSubmit={handleSubmit} isSubmitting={isPending} />
        </div>
      </div>
    </div>
  );
};
