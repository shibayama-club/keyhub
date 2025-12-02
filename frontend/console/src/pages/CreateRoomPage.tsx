import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import * as Sentry from '@sentry/react';
import { Navbar } from '../components/Navbar';
import { CreateRoomForm } from '../components/CreateRoomForm';
import { useMutationCreateRoom, queryClient } from '../libs/query';
import { RoomType } from '../../../gen/src/keyhub/console/v1/common_pb';

export const CreateRoomPage = () => {
  const navigate = useNavigate();
  const { mutateAsync: createRoom, isPending } = useMutationCreateRoom();

  const handleSubmit = async (data: {
    name: string;
    buildingName: string;
    floorNumber: string;
    roomType: RoomType;
    description?: string;
  }) => {
    try {
      await createRoom({
        name: data.name,
        buildingName: data.buildingName,
        floorNumber: data.floorNumber,
        roomType: data.roomType,
        description: data.description || '',
      });

      // TanStack Queryのキャッシュを無効化
      await queryClient.invalidateQueries();

      toast.success('部屋を作成しました');
      // 成功後は部屋一覧に戻る
      navigate('/rooms', { replace: true });
    } catch (error) {
      // Sentryでエラーをキャプチャ
      Sentry.captureException(error);
      toast.error('部屋の作成に失敗しました');
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">新しい部屋を作成</h1>
          <p className="mt-2 text-sm text-gray-600">組織の新しい部屋を作成します。以下の詳細を入力してください。</p>
        </div>

        <div className="rounded-lg border-2 border-gray-200 bg-white p-8 shadow-md">
          <CreateRoomForm onSubmit={handleSubmit} isSubmitting={isPending} />
        </div>
      </div>
    </div>
  );
};
