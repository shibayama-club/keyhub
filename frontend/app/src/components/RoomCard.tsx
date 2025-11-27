import type { Room } from '@keyhub/gen/keyhub/app/v1/common_pb.ts';
import { ROOM_TYPE_LABELS, KEY_STATUS_LABELS, KEY_STATUS_TEXT_COLORS } from '../lib/constants/room';

interface RoomCardProps {
  room: Room;
}

export const RoomCard = ({ room }: RoomCardProps) => {
  return (
    <div className="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <h3 className="text-lg font-semibold text-gray-900">{room.name}</h3>
          <p className="mt-1 text-sm text-gray-600">{room.description || '説明なし'}</p>
        </div>
        <span className="rounded-full bg-blue-100 px-3 py-1 text-xs font-medium text-blue-800">
          {ROOM_TYPE_LABELS[room.roomType] || '未指定'}
        </span>
      </div>

      <div className="mt-4 flex items-center gap-4 text-sm text-gray-500">
        <div className="flex items-center gap-1">
          <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
            />
          </svg>
          <span>{room.buildingName}</span>
        </div>
        <div className="flex items-center gap-1">
          <span>{room.floorNumber}階</span>
        </div>
      </div>

      {room.keys.length > 0 && (
        <div className="mt-4">
          <p className="mb-2 text-xs font-medium text-gray-500">鍵情報:</p>
          <div className="flex flex-wrap gap-2">
            {room.keys.map((key) => (
              <span
                key={key.id}
                className="inline-flex items-center space-x-1 rounded-md bg-gray-100 px-2 py-1 text-xs"
              >
                <span className="font-medium text-gray-900">{key.keyNumber}</span>
                <span className="text-gray-400">-</span>
                <span className={KEY_STATUS_TEXT_COLORS[key.status] || 'text-gray-700'}>
                  {KEY_STATUS_LABELS[key.status] || '不明'}
                </span>
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};
