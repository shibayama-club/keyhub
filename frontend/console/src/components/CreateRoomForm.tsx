import { useMemo } from 'react';
import { useForm } from '../hooks/useForm';
import { useFormField } from '../hooks/useFormField';
import { RoomType } from '../../../gen/src/keyhub/console/v1/common_pb';
import { roomSchema, type RoomFormData } from '../libs/utils/schema';
import { ROOM_TYPE_OPTIONS } from '../constants/room';

type CreateRoomFormProps = {
  onSubmit: (data: RoomFormData) => void;
  isSubmitting?: boolean;
};

export const CreateRoomForm = ({ onSubmit, isSubmitting = false }: CreateRoomFormProps) => {
  const initialValues = useMemo(
    () => ({
      name: '',
      buildingName: '',
      floorNumber: '',
      roomType: RoomType.CLASSROOM,
      description: '',
    }),
    [],
  );

  const form = useForm(roomSchema, {
    revalidate: true,
    initialValues,
  });

  const nameField = useFormField(form, 'name');
  const buildingNameField = useFormField(form, 'buildingName');
  const floorNumberField = useFormField(form, 'floorNumber');
  const roomTypeField = useFormField(form, 'roomType', {
    transform: (value) => Number(value) as RoomType,
  });
  const descriptionField = useFormField(form, 'description');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const result = form.validate();
    if (result.success) {
      onSubmit(result.data);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-8">
      <div>
        <label htmlFor="name" className="mb-2 block text-sm font-medium text-gray-700">
          部屋名 <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          id="name"
          value={nameField.value || ''}
          onChange={nameField.onChange}
          onBlur={nameField.onBlur}
          className="mt-1 block w-full rounded-md border-gray-300 px-4 py-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          disabled={isSubmitting}
        />
        {nameField.error.length > 0 && <p className="mt-2 text-sm text-red-600">{nameField.error[0]}</p>}
      </div>

      <div>
        <label htmlFor="buildingName" className="mb-2 block text-sm font-medium text-gray-700">
          建物名 <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          id="buildingName"
          value={buildingNameField.value || ''}
          onChange={buildingNameField.onChange}
          onBlur={buildingNameField.onBlur}
          className="mt-1 block w-full rounded-md border-gray-300 px-4 py-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          disabled={isSubmitting}
        />
        {buildingNameField.error.length > 0 && (
          <p className="mt-2 text-sm text-red-600">{buildingNameField.error[0]}</p>
        )}
      </div>

      <div>
        <label htmlFor="floorNumber" className="mb-2 block text-sm font-medium text-gray-700">
          階数 <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          id="floorNumber"
          value={floorNumberField.value || ''}
          onChange={floorNumberField.onChange}
          onBlur={floorNumberField.onBlur}
          placeholder="例: 3F, B1F"
          className="mt-1 block w-full rounded-md border-gray-300 px-4 py-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          disabled={isSubmitting}
        />
        {floorNumberField.error.length > 0 && (
          <p className="mt-2 text-sm text-red-600">{floorNumberField.error[0]}</p>
        )}
      </div>

      <div>
        <label htmlFor="roomType" className="mb-2 block text-sm font-medium text-gray-700">
          部屋タイプ <span className="text-red-500">*</span>
        </label>
        <select
          id="roomType"
          value={roomTypeField.value ?? RoomType.CLASSROOM}
          onChange={roomTypeField.onChange}
          onBlur={roomTypeField.onBlur}
          className="mt-1 block w-full rounded-md border-gray-300 px-4 py-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          disabled={isSubmitting}
        >
          {ROOM_TYPE_OPTIONS.map((type) => (
            <option key={type.value} value={type.value}>
              {type.label}
            </option>
          ))}
        </select>
        {roomTypeField.error.length > 0 && <p className="mt-2 text-sm text-red-600">{roomTypeField.error[0]}</p>}
      </div>

      <div>
        <label htmlFor="description" className="mb-2 block text-sm font-medium text-gray-700">
          説明
        </label>
        <textarea
          id="description"
          rows={4}
          value={descriptionField.value || ''}
          onChange={descriptionField.onChange}
          onBlur={descriptionField.onBlur}
          className="mt-1 block w-full rounded-md border-gray-300 px-4 py-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          disabled={isSubmitting}
        />
        {descriptionField.error.length > 0 && <p className="mt-2 text-sm text-red-600">{descriptionField.error[0]}</p>}
      </div>

      <div className="flex justify-end space-x-3 pt-4">
        <button
          type="submit"
          disabled={isSubmitting}
          className="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 px-6 py-3 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
        >
          {isSubmitting ? '作成中...' : '作成'}
        </button>
      </div>
    </form>
  );
};
