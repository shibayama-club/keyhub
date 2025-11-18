import type { Tenant } from '@keyhub/gen/keyhub/app/v1/app_pb.ts';
import { TENANT_TYPE_LABELS } from '../lib/constants/tenant';
import { formatTimestampToJapaneseDate } from '../utils/date';

interface TenantCardProps {
  tenant: Tenant;
}

export const TenantCard = ({ tenant }: TenantCardProps) => {
  const formattedDate = formatTimestampToJapaneseDate(tenant.createdAt);

  return (
    <div className="rounded-lg border border-gray-200 bg-white p-6 shadow-sm transition-shadow hover:shadow-md">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <h3 className="text-lg font-semibold text-gray-900">{tenant.name}</h3>
          <p className="mt-1 text-sm text-gray-600">{tenant.description || '説明なし'}</p>
        </div>
        <span className="rounded-full bg-indigo-100 px-3 py-1 text-xs font-medium text-indigo-800">
          {TENANT_TYPE_LABELS[tenant.tenantType] || '未指定'}
        </span>
      </div>

      <div className="mt-4 flex items-center gap-4 text-sm text-gray-500">
        <div className="flex items-center gap-1">
          <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
            />
          </svg>
          <span>{tenant.memberCount}人</span>
        </div>
        <div className="flex items-center gap-1">
          <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
            />
          </svg>
          <span>{formattedDate}</span>
        </div>
      </div>
    </div>
  );
};
