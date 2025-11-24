import { useState } from 'react';

type AssignRoomFormProps = {
  onSubmit: (data: { roomId: string; expiresAt?: Date }) => void;
  isSubmitting?: boolean;
};

export const AssignRoomForm = ({ onSubmit, isSubmitting = false }: AssignRoomFormProps) => {
  const [roomId, setRoomId] = useState('');
  const [expiresAt, setExpiresAt] = useState('');
  const [errors, setErrors] = useState<{ roomId?: string; expiresAt?: string }>({});

  const validateForm = () => {
    const newErrors: { roomId?: string; expiresAt?: string } = {};

    if (!roomId.trim()) {
      newErrors.roomId = 'Room IDを入力してください';
    } else if (!/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(roomId.trim())) {
      newErrors.roomId = '有効なUUID形式のRoom IDを入力してください';
    }

    if (expiresAt && isNaN(new Date(expiresAt).getTime())) {
      newErrors.expiresAt = '有効な日時を入力してください';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    const data: { roomId: string; expiresAt?: Date } = {
      roomId: roomId.trim(),
    };

    if (expiresAt) {
      data.expiresAt = new Date(expiresAt);
    }

    onSubmit(data);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label htmlFor="roomId" className="mb-2 block text-sm font-medium text-gray-700">
          Room ID <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          id="roomId"
          value={roomId}
          onChange={(e) => {
            setRoomId(e.target.value);
            if (errors.roomId) setErrors({ ...errors, roomId: undefined });
          }}
          placeholder="例: 550e8400-e29b-41d4-a716-446655440000"
          className="mt-1 block w-full rounded-md border-gray-300 px-4 py-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          disabled={isSubmitting}
        />
        {errors.roomId && <p className="mt-2 text-sm text-red-600">{errors.roomId}</p>}
        <p className="mt-1 text-xs text-gray-500">作成済みのRoom IDを入力してください</p>
      </div>

      <div>
        <label htmlFor="expiresAt" className="mb-2 block text-sm font-medium text-gray-700">
          有効期限
        </label>
        <input
          type="datetime-local"
          id="expiresAt"
          value={expiresAt}
          onChange={(e) => {
            setExpiresAt(e.target.value);
            if (errors.expiresAt) setErrors({ ...errors, expiresAt: undefined });
          }}
          className="mt-1 block w-full rounded-md border-gray-300 px-4 py-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          disabled={isSubmitting}
        />
        {errors.expiresAt && <p className="mt-2 text-sm text-red-600">{errors.expiresAt}</p>}
        <p className="mt-1 text-xs text-gray-500">指定しない場合は無期限で割り当てられます</p>
      </div>

      <div className="flex justify-end space-x-3 pt-4">
        <button
          type="submit"
          disabled={isSubmitting}
          className="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 px-6 py-3 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
        >
          {isSubmitting ? '割り当て中...' : '割り当て'}
        </button>
      </div>
    </form>
  );
};
