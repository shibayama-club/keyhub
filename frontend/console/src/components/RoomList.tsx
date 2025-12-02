import { useNavigate } from 'react-router-dom';
import type { Room } from '../../../gen/src/keyhub/console/v1/common_pb';
import { getRoomTypeLabel } from '../utils/room';
import { KeyBadgeList } from './KeyBadgeList';

export const RoomList = ({ rooms, isLoading, isError }: { rooms: Room[]; isLoading: boolean; isError: boolean }) => {
  const navigate = useNavigate();
  if (isLoading) {
    return (
      <div className="px-4 py-12 text-center">
        <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-indigo-600 border-r-transparent"></div>
        <p className="mt-2 text-sm text-gray-600">読み込み中...</p>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="px-4 py-5 sm:px-6">
        <p className="text-red-600">Roomの取得に失敗しました。</p>
      </div>
    );
  }

  if (rooms.length === 0) {
    return (
      <div className="px-4 py-5 sm:px-6">
        <p className="text-gray-600">Roomが見つかりません。最初のRoomを作成して始めましょう。</p>
      </div>
    );
  }

  return (
    <ul className="divide-y divide-gray-200">
      {rooms.map((room) => (
        <li key={room.id} className="px-4 py-5 hover:bg-gray-50 sm:px-6">
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <div className="flex items-center space-x-3">
                <h4 className="text-base font-medium text-gray-900">{room.name}</h4>
                <span className="inline-flex items-center rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-medium text-blue-800">
                  {getRoomTypeLabel(room.roomType)}
                </span>
              </div>
              <div className="mt-1 flex items-center space-x-3 text-sm text-gray-600">
                <span>{room.buildingName}</span>
                <span>•</span>
                <span>{room.floorNumber}階</span>
              </div>
              {room.description && <p className="mt-1 text-sm text-gray-600">{room.description}</p>}
              <KeyBadgeList keys={room.keys} />
            </div>
            <div className="flex space-x-2">
              <button
                onClick={() => navigate(`/rooms/${room.id}/keys`)}
                className="inline-flex items-center rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-none"
              >
                <svg className="mr-1.5 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
                鍵一覧
              </button>
              <button
                onClick={() => navigate(`/rooms/${room.id}/keys/create`)}
                className="inline-flex items-center rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-none"
              >
                <svg className="mr-1.5 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                  />
                </svg>
                鍵を作成
              </button>
            </div>
          </div>
        </li>
      ))}
    </ul>
  );
};
