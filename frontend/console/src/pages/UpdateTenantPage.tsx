import { useParams, useNavigate } from 'react-router-dom';
import { useQueryGetTenantById, useMutationUpdateTenant } from '../libs/query';

import { Navbar } from '../components/Navbar';
import { TenantType } from '../../../gen/src/keyhub/console/v1/console_pb';
import { timestampFromDate } from '@bufbuild/protobuf/wkt';
import toast from 'react-hot-toast';
import * as Sentry from '@sentry/react';
import { queryClient } from '../libs/query';
import type { TenantFormData } from '../libs/utils/schema';
import { UpdateTenantForm } from '../components/UpdateTenantForm';

export const UpdateTenantPage = () => {
  const { tenantId } = useParams<{ tenantId: string }>();
  const navigate = useNavigate();
  const { mutateAsync: updateTenant, isPending } = useMutationUpdateTenant();
  const {
    data: tenantData,
    isLoading: isTenantLoading,
    isError: isTenantError,
  } = useQueryGetTenantById(tenantId || '');

  if (!tenantId) {
    navigate('/tenants');
    return null;
  }

  const handleUpdateTenant = async (data: TenantFormData) => {
    try {
      await updateTenant({
        id: tenantId,
        name: data.name,
        description: data.description,
        tenantType: data.tenantType as TenantType,
        joinCode: data.joinCode,
        joinCodeExpiry: data.joinCodeExpiry ? timestampFromDate(data.joinCodeExpiry) : undefined,
        joinCodeMaxUse: data.joinCodeMaxUse,
      });
      await queryClient.invalidateQueries();
      toast.success('テナント情報を更新しました');
    } catch (error) {
      Sentry.captureException(error);
      toast.error('テナント情報の更新に失敗しました');
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <button
            onClick={() => navigate('/tenants')}
            className="mb-4 inline-flex items-center text-sm text-gray-600 hover:text-gray-900"
          >
            <svg className="mr-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            テナント一覧に戻る
          </button>
          <h1 className="text-3xl font-bold text-gray-900">テナント情報の更新</h1>
          <p className="mt-2 text-sm text-gray-600">このテナントの情報を更新します。</p>
        </div>

        <div className="rounded-lg border-2 border-gray-200 bg-white p-8 shadow-md">
          <div className="mb-6">
            {isTenantLoading ? (
              <div className="mt-2 rounded-md bg-gray-50 p-4">
                <div className="flex items-center">
                  <p className="ml-2 text-sm text-gray-600">読み込み中...</p>
                </div>
              </div>
            ) : isTenantError ? (
              <div className="mt-2 rounded-md bg-red-50 p-4">
                <p className="text-sm text-red-600">テナント情報の取得に失敗しました</p>
              </div>
            ) : null}
          </div>

          <div className="border-t border-gray-200 pt-6">
            <h2 className="mb-4 text-lg font-semibold text-gray-900">テナント情報の更新</h2>
            {tenantData ? (
              <UpdateTenantForm onSubmit={handleUpdateTenant} isSubmitting={isPending} tenantData={tenantData} />
            ) : (
              <div className="text-sm text-gray-500">テナント情報を読み込み中...</div>
            )}
            <div>{tenantData?.tenant?.name}</div>
          </div>
        </div>
      </div>
    </div>
  );
};
