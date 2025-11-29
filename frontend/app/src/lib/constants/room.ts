/**
 * 部屋タイプのラベル定義
 */
export const ROOM_TYPE_LABELS: Record<number, string> = {
  0: '未指定',
  1: '教室',
  2: '会議室',
  3: '実験室',
  4: 'オフィス',
  5: '作業室',
  6: '倉庫',
} as const;

/**
 * 鍵ステータスのラベル定義
 */
export const KEY_STATUS_LABELS: Record<number, string> = {
  0: '未指定',
  1: '利用可能',
  2: '使用中',
  3: '紛失',
  4: '破損',
} as const;

/**
 * 鍵ステータスの色定義（テキスト用）
 */
export const KEY_STATUS_TEXT_COLORS: Record<number, string> = {
  0: 'text-gray-700',
  1: 'text-green-700',
  2: 'text-yellow-700',
  3: 'text-red-700',
  4: 'text-red-700',
} as const;
