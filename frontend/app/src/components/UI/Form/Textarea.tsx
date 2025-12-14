import { forwardRef, type TextareaHTMLAttributes } from 'react';

type TextareaProps = {
  label?: string;
  error?: string;
  hint?: string;
  required?: boolean;
} & TextareaHTMLAttributes<HTMLTextAreaElement>;

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ label, error, hint, required, className = '', id, ...props }, ref) => {
    const textareaId = id || props.name;
    const hasError = !!error;

    const baseStyles =
      'block w-full rounded-md border px-4 py-2 text-sm transition-colors focus:outline-none focus:ring-2 disabled:cursor-not-allowed disabled:bg-gray-50 disabled:opacity-50';

    const textareaStyles = hasError
      ? 'border-red-300 text-red-900 placeholder-red-300 focus:border-red-500 focus:ring-red-500'
      : 'border-gray-300 focus:border-orange-500 focus:ring-orange-500';

    return (
      <div className={className}>
        {label && (
          <label htmlFor={textareaId} className="mb-2 block text-sm font-medium text-gray-700">
            {label}
            {required && <span className="ml-1 text-red-500">*</span>}
          </label>
        )}
        <textarea
          ref={ref}
          id={textareaId}
          className={`${baseStyles} ${textareaStyles}`}
          aria-invalid={hasError}
          aria-describedby={hasError ? `${textareaId}-error` : hint ? `${textareaId}-hint` : undefined}
          {...props}
        />
        {error && (
          <p id={`${textareaId}-error`} className="mt-2 text-sm text-red-600">
            {error}
          </p>
        )}
        {hint && !error && (
          <p id={`${textareaId}-hint`} className="mt-2 text-sm text-gray-500">
            {hint}
          </p>
        )}
      </div>
    );
  },
);

Textarea.displayName = 'Textarea';
