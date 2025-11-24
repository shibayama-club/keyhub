import { useState } from 'react';

type CreateKeyFormProps = {
  onSubmit: (data: { keyNumber: string }) => void;
  isSubmitting?: boolean;
};

export const CreateKeyForm = ({ onSubmit, isSubmitting }: CreateKeyFormProps) => {
  const [keyNumber, setKeyNumber] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ keyNumber });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label htmlFor="keyNumber" className="block text-sm font-medium text-gray-700">
          鍵番号 <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          id="keyNumber"
          value={keyNumber}
          onChange={(e) => setKeyNumber(e.target.value)}
          required
          className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none sm:text-sm"
          placeholder="例: K001"
        />
      </div>

      <div className="flex justify-end space-x-3">
        <button
          type="submit"
          disabled={isSubmitting}
          className="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-none"
        >
          {isSubmitting ? '作成中...' : '鍵を作成'}
        </button>
      </div>
    </form>
  );
};
