import { z } from 'zod';
import type { UseFormReturn } from './useForm';

export type UseFormFieldProps<T> = {
  value: T;
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
  onBlur: () => void;
  error: string[];
};

export const useFormField = <T extends z.ZodObject<z.ZodRawShape>, F extends keyof z.infer<T>>(
  form: UseFormReturn<T>,
  field: F,
): UseFormFieldProps<z.infer<T>[F]> => {
  return {
    value: form.state[field],
    onChange: (e) => {
      form.updateField(field, e.target.value as z.infer<T>[F]);
    },
    onBlur: () => {
      form.validateField(field);
    },
    error: form.errors[field] || [],
  };
};
