/**
 * テナントタイプのラベル定義
 */
export const TENANT_TYPE_LABELS: Record<number, string> = {
  0: '未指定',
  1: 'チーム',
  2: '部門',
  3: 'プロジェクト',
  4: 'ラボ',
} as const;
