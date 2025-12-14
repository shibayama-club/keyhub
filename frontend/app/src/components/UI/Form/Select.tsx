import { forwardRef, type SelectHTMLAttributes } from 'react';

type SelectOption = {
  value: string | number;
  label: string;
};

type SelectProps = {
  label?: string;
  error?: string;
  hint?: string;
  required?: boolean;
  options: SelectOption[];
  placeholder?: string;
} & Omit<SelectHTMLAttributes<HTMLSelectElement>, 'children'>;

export const Select = forwardRef<HTMLSelectElement, SelectProps>(
  ({ label, error, hint, required, options, placeholder, className = '', id, ...props }, ref) => {
    const selectId = id || props.name;
    const hasError = !!error;

    const baseStyles =
      'block w-full rounded-md border px-4 py-2 text-sm transition-colors focus:outline-none focus:ring-2 disabled:cursor-not-allowed disabled:bg-gray-50 disabled:opacity-50';

    const selectStyles = hasError
      ? 'border-red-300 text-red-900 focus:border-red-500 focus:ring-red-500'
      : 'border-gray-300 focus:border-orange-500 focus:ring-orange-500';

    return (
      <div className={className}>
        {label && (
          <label htmlFor={selectId} className="mb-2 block text-sm font-medium text-gray-700">
            {label}
            {required && <span className="ml-1 text-red-500">*</span>}
          </label>
        )}
        <select
          ref={ref}
          id={selectId}
          className={`${baseStyles} ${selectStyles}`}
          aria-invalid={hasError}
          aria-describedby={hasError ? `${selectId}-error` : hint ? `${selectId}-hint` : undefined}
          {...props}
        >
          {placeholder && (
            <option value="" disabled>
              {placeholder}
            </option>
          )}
          {options.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
        {error && (
          <p id={`${selectId}-error`} className="mt-2 text-sm text-red-600">
            {error}
          </p>
        )}
        {hint && !error && (
          <p id={`${selectId}-hint`} className="mt-2 text-sm text-gray-500">
            {hint}
          </p>
        )}
      </div>
    );
  },
);

Select.displayName = 'Select';
