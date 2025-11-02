import { TENANT_TYPE_OPTIONS } from '../constants/tenant';

export const getTenantTypeLabel = (tenantType: number): string => {
  const option = TENANT_TYPE_OPTIONS.find((opt) => opt.value === tenantType);
  return option?.label || '不明';
};
