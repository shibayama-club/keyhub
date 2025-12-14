import { useController, type Control, type FieldValues, type Path } from 'react-hook-form';
import { Select } from './Select';

type SelectOption = {
  value: string | number;
  label: string;
};

type FormSelectProps<T extends FieldValues> = {
  name: Path<T>;
  control: Control<T>;
  options: SelectOption[];
  label?: string;
  hint?: string;
  required?: boolean;
  placeholder?: string;
  disabled?: boolean;
  className?: string;
};

export const FormSelect = <T extends FieldValues>({
  name,
  control,
  options,
  label,
  hint,
  required,
  placeholder,
  disabled,
  className,
}: FormSelectProps<T>) => {
  const {
    field,
    fieldState: { error },
  } = useController({ name, control });

  return (
    <Select
      {...field}
      options={options}
      label={label}
      error={error?.message}
      hint={hint}
      required={required}
      placeholder={placeholder}
      disabled={disabled}
      className={className}
    />
  );
};
