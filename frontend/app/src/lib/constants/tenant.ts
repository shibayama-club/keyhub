import { TenantType } from '../../../../gen/src/keyhub/app/v1/app_pb';

/**
 * テナントタイプのラベル定義
 */
export const TENANT_TYPE_LABELS: Record<TenantType, string> = {
  [TenantType.UNSPECIFIED]: '未指定',
  [TenantType.TEAM]: 'チーム',
  [TenantType.DEPARTMENT]: '部署',
  [TenantType.PROJECT]: 'プロジェクト',
  [TenantType.LABORATORY]: '研究室',
} as const;
