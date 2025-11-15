import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import * as Sentry from '@sentry/react';
import { Navbar } from '../components/Navbar';
import { CreateTenantForm } from '../components/CreateTenantForm';
import { useMutationCreateTenant, queryClient } from '../libs/query';
import { TenantType } from '../../../gen/src/keyhub/console/v1/console_pb';
import { timestampFromDate } from '@bufbuild/protobuf/wkt';

export const CreateTenantPage = () => {
  const navigate = useNavigate();
  const { mutateAsync: createTenant, isPending } = useMutationCreateTenant();

  const handleSubmit = async (data: {
    name: string;
    description?: string;
    tenantType: TenantType;
    joinCode: string;
    joinCodeExpiry?: Date;
    joinCodeMaxUse?: number;
  }) => {
    try {
      await createTenant({
        name: data.name,
        description: data.description || '',
        tenantType: data.tenantType,
        joinCode: data.joinCode,
        joinCodeExpiry: data.joinCodeExpiry ? timestampFromDate(data.joinCodeExpiry) : undefined,
        joinCodeMaxUse: data.joinCodeMaxUse ?? 0,
      });

      // TanStack Queryのキャッシュを無効化してテナント一覧を最新に保つ
      await queryClient.invalidateQueries();

      toast.success('テナントを作成しました');
      // 成功後はテナント一覧に戻る
      navigate('/tenants', { replace: true });
    } catch (error) {
      // Sentryでエラーをキャプチャ
      Sentry.captureException(error);
      toast.error('テナントの作成に失敗しました');
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">新しいテナントを作成</h1>
          <p className="mt-2 text-sm text-gray-600">組織の新しいテナントを作成します。以下の詳細を入力してください。</p>
        </div>

        <div className="rounded-lg border-2 border-gray-200 bg-white p-8 shadow-md">
          <CreateTenantForm onSubmit={handleSubmit} isSubmitting={isPending} />
        </div>
      </div>
    </div>
  );
};
