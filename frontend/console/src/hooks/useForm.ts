import { useState, useEffect, useMemo } from 'react';
import { z } from 'zod';

export type UseFormReturn<T extends z.ZodObject<z.ZodRawShape>> = {
  state: z.infer<T>;
  errors: { [K in keyof z.infer<T>]: string[] };
  updateField: <K extends keyof z.infer<T>>(key: K, value: z.infer<T>[K]) => void;
  validateField: (key: keyof z.infer<T>) => boolean;
  validate: () => z.ZodSafeParseResult<z.infer<T>>;
  setState: React.Dispatch<React.SetStateAction<z.infer<T>>>;
};

export const useForm = <T extends z.ZodObject<z.ZodRawShape>>(
  schema: T,
  options: { revalidate?: boolean; initialValues?: Partial<z.infer<T>> } = {},
): UseFormReturn<T> => {
  type FormType = z.infer<T>;
  type FormErrorsType = { [K in keyof FormType]: string[] };

  const [state, setState] = useState<FormType>(
    options.initialValues ? { ...({} as FormType), ...options.initialValues } : ({} as FormType),
  );

  const memoizedInitialValues = useMemo(() => {
    return options.initialValues || {};
  }, [options.initialValues]);

  useEffect(() => {
    if (Object.keys(memoizedInitialValues).length > 0) {
      setState((prev) => {
        const next = { ...prev } as FormType;
        (Object.keys(memoizedInitialValues) as Array<keyof FormType>).forEach((key) => {
          const incoming = (memoizedInitialValues as Partial<FormType>)[key];
          if (prev[key] === undefined && incoming !== undefined) {
            next[key] = incoming as FormType[typeof key];
          }
        });
        return next;
      });
    }
  }, [memoizedInitialValues]);
  const [errors, setErrors] = useState<FormErrorsType>({} as FormErrorsType);
  const [fieldToValidate, setFieldToValidate] = useState<keyof FormType | null>(null);

  useEffect(() => {
    if (fieldToValidate && options.revalidate) {
      const result = (schema.shape[fieldToValidate as string] as z.ZodTypeAny).safeParse(state[fieldToValidate]);
      setErrors((prev) => ({
        ...prev,
        [fieldToValidate]: result.success ? [] : result.error.issues.map((issue) => issue.message),
      }));
      setFieldToValidate(null);
    }
  }, [state, fieldToValidate, options.revalidate, schema]);

  const updateField = (key: keyof FormType, value: FormType[keyof FormType]) => {
    setState((prev) => ({ ...prev, [key]: value }));
    if (options.revalidate) {
      setFieldToValidate(key);
    }
  };

  const validateField = (key: keyof FormType) => {
    const result = (schema.shape[key as string] as z.ZodTypeAny).safeParse(state[key]);
    setErrors((prev) => ({
      ...prev,
      [key]: result.success ? [] : result.error.issues.map((issue) => issue.message),
    }));
    return result.success;
  };

  const validate = () => {
    const result = schema.safeParse(state);
    if (!result.success) {
      const newErrors = {} as FormErrorsType;
      result.error.issues.forEach((issue) => {
        const field = issue.path[0] as keyof FormType;
        if (!newErrors[field]) {
          newErrors[field] = [];
        }
        newErrors[field].push(issue.message);
      });
      setErrors(newErrors);
    }
    return result;
  };

  return {
    state,
    errors,
    updateField,
    validateField,
    validate,
    setState,
  };
};
