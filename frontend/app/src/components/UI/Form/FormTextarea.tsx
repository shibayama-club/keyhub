import { useController, type Control, type FieldValues, type Path } from 'react-hook-form';
import { Textarea } from './Textarea';

type FormTextareaProps<T extends FieldValues> = {
  name: Path<T>;
  control: Control<T>;
  label?: string;
  hint?: string;
  required?: boolean;
  placeholder?: string;
  disabled?: boolean;
  rows?: number;
  className?: string;
};

export const FormTextarea = <T extends FieldValues>({
  name,
  control,
  label,
  hint,
  required,
  placeholder,
  disabled,
  rows = 4,
  className,
}: FormTextareaProps<T>) => {
  const {
    field,
    fieldState: { error },
  } = useController({ name, control });

  return (
    <Textarea
      {...field}
      label={label}
      error={error?.message}
      hint={hint}
      required={required}
      placeholder={placeholder}
      disabled={disabled}
      rows={rows}
      className={className}
    />
  );
};
