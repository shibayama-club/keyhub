import { z } from 'zod';
import type { UseFormReturn } from './useForm';

export type UseFormFieldProps<T> = {
  value: T;
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => void;
  onBlur: () => void;
  error: string[];
};

export const useFormField = <T extends z.ZodObject<z.ZodRawShape>, F extends keyof z.infer<T>>(
  form: UseFormReturn<T>,
  field: F,
  options?: {
    transform?: (value: string) => z.infer<T>[F];
  },
): UseFormFieldProps<z.infer<T>[F]> => {
  return {
    value: form.state[field],
    onChange: (e) => {
      const value = options?.transform ? options.transform(e.target.value) : (e.target.value as z.infer<T>[F]);
      form.updateField(field, value);
    },
    onBlur: () => {
      form.validateField(field);
    },
    error: form.errors[field] || [],
  };
};
