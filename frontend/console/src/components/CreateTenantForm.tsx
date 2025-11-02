import { useMemo } from 'react';
import { useForm } from '../hooks/useForm';
import { useFormField } from '../hooks/useFormField';
import { TenantType } from '../../../gen/src/keyhub/console/v1/console_pb';
import { tenantSchema, type TenantFormData } from '../libs/utils/schema';
import { TENANT_TYPE_OPTIONS } from '../constants/tenant';

type CreateTenantFormProps = {
  onSubmit: (data: TenantFormData) => void;
  isSubmitting?: boolean;
};

export const CreateTenantForm = ({ onSubmit, isSubmitting = false }: CreateTenantFormProps) => {
  // initialValuesをメモ化して無限ループを防ぐ
  const initialValues = useMemo(
    () => ({
      name: '',
      description: '',
      tenantType: TenantType.TEAM,
    }),
    [],
  );

  const form = useForm(tenantSchema, {
    revalidate: true,
    initialValues,
  });

  const nameField = useFormField(form, 'name');
  const descriptionField = useFormField(form, 'description');
  const tenantTypeField = useFormField(form, 'tenantType', {
    transform: (value) => Number(value) as TenantType,
  });

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
          テナント名 <span className="text-red-500">*</span>
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

      <div>
        <label htmlFor="tenantType" className="mb-2 block text-sm font-medium text-gray-700">
          テナントタイプ <span className="text-red-500">*</span>
        </label>
        <select
          id="tenantType"
          value={tenantTypeField.value ?? TenantType.TEAM}
          onChange={tenantTypeField.onChange}
          onBlur={tenantTypeField.onBlur}
          className="mt-1 block w-full rounded-md border-gray-300 px-4 py-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          disabled={isSubmitting}
        >
          {TENANT_TYPE_OPTIONS.map((type) => (
            <option key={type.value} value={type.value}>
              {type.label}
            </option>
          ))}
        </select>
        {tenantTypeField.error.length > 0 && <p className="mt-2 text-sm text-red-600">{tenantTypeField.error[0]}</p>}
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
