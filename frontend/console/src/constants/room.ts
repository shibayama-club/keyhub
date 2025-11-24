import { RoomType } from '../../../gen/src/keyhub/console/v1/common_pb';

export const ROOM_TYPE_OPTIONS = [
  { value: RoomType.CLASSROOM, label: '教室' },
  { value: RoomType.MEETING_ROOM, label: '会議室' },
  { value: RoomType.LABORATORY, label: '実験室' },
  { value: RoomType.OFFICE, label: 'オフィス' },
  { value: RoomType.WORKSHOP, label: '作業室' },
  { value: RoomType.STORAGE, label: '倉庫' },
] as const;
