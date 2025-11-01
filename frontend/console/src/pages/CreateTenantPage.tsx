import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import * as Sentry from '@sentry/react';
import { Navbar } from '../components/Navbar';
import { CreateTenantForm } from '../components/CreateTenantForm';
import { useMutationCreateTenant } from '../libs/query';
import { TenantType } from '../../../gen/src/keyhub/console/v1/console_pb';

export const CreateTenantPage = () => {
  const navigate = useNavigate();
  const { mutate: createTenant, isPending } = useMutationCreateTenant();

  const handleSubmit = (data: { name: string; description?: string; tenantType: TenantType }) => {
    createTenant(
      {
        name: data.name,
        description: data.description || '',
        tenantType: data.tenantType,
      },
      {
        onSuccess: () => {
          toast.success('テナントを作成しました');
          navigate('/tenants');
        },
        onError: (error) => {
          Sentry.captureException(error);
          toast.error('テナントの作成に失敗しました');
        },
      },
    );
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
