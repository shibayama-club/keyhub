import { TenantType } from '../../../gen/src/keyhub/console/v1/console_pb';

export const TENANT_TYPE_OPTIONS = [
  { value: TenantType.TEAM, label: 'チーム' },
  { value: TenantType.DEPARTMENT, label: '部署' },
  { value: TenantType.PROJECT, label: 'プロジェクト' },
  { value: TenantType.LABORATORY, label: '研究室' },
] as const;
