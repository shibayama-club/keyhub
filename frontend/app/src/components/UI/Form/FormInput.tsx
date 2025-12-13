import { type ReactNode } from 'react';
import { useController, type Control, type FieldValues, type Path } from 'react-hook-form';
import { Input } from './Input';

type FormInputProps<T extends FieldValues> = {
  name: Path<T>;
  control: Control<T>;
  label?: string;
  hint?: string;
  required?: boolean;
  leftIcon?: ReactNode;
  rightIcon?: ReactNode;
  type?: string;
  placeholder?: string;
  disabled?: boolean;
  className?: string;
};

export const FormInput = <T extends FieldValues>({
  name,
  control,
  label,
  hint,
  required,
  leftIcon,
  rightIcon,
  type = 'text',
  placeholder,
  disabled,
  className,
}: FormInputProps<T>) => {
  const {
    field,
    fieldState: { error },
  } = useController({ name, control });

  return (
    <Input
      {...field}
      type={type}
      label={label}
      error={error?.message}
      hint={hint}
      required={required}
      leftIcon={leftIcon}
      rightIcon={rightIcon}
      placeholder={placeholder}
      disabled={disabled}
      className={className}
    />
  );
};
