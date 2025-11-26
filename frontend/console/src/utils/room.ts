import { ROOM_TYPE_OPTIONS } from '../constants/room';

export const getRoomTypeLabel = (roomType: number): string => {
  const option = ROOM_TYPE_OPTIONS.find((opt) => opt.value === roomType);
  return option?.label || '不明';
};
