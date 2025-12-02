import { KEY_STATUS_OPTIONS } from '../constants/key';

export const getKeyStatusLabel = (keyStatus: number): string => {
  const option = KEY_STATUS_OPTIONS.find((opt) => opt.value === keyStatus);
  return option?.label || '不明';
};

export const getKeyStatusTextColor = (status: number): string => {
  if (status === 1) return 'text-green-700';
  if (status === 2) return 'text-yellow-700';
  if (status === 3) return 'text-red-700';
  return 'text-gray-700';
};

export const getKeyStatusBadgeColor = (status: number): string => {
  if (status === 1) return 'bg-green-100 text-green-800';
  if (status === 2) return 'bg-yellow-100 text-yellow-800';
  if (status === 3) return 'bg-red-100 text-red-800';
  return 'bg-gray-100 text-gray-800';
};
