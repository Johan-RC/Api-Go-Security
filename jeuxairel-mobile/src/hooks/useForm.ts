import { useCallback, useState } from 'react';

export type Validator<T = string> = (value: T) => string | undefined;

export type FieldName = string;

export interface UseFormOptions<TValues extends object> {
  initial: TValues;
  validators: Partial<Record<keyof TValues, Validator<TValues[keyof TValues]>>>;
}

/**
 * Mini sistema de formularios con validación contextual:
 * el error se muestra al salir del campo (onBlur) y se limpia al cambiar
 * el valor (onChange). `validate()` recalcula todos los errores al enviar.
 */
export function useForm<TValues extends object>({ initial, validators }: UseFormOptions<TValues>) {
  const [values, setValues] = useState<TValues>(initial);
  const [errors, setErrors] = useState<Partial<Record<keyof TValues, string>>>({});
  const [touched, setTouched] = useState<Partial<Record<keyof TValues, boolean>>>({});

  const validateField = useCallback(
    (field: keyof TValues, value: TValues[keyof TValues]) => {
      const validator = validators[field];
      if (!validator) return undefined as string | undefined;
      return validator(value);
    },
    [validators],
  );

  const setField = useCallback(
    <K extends keyof TValues>(field: K, value: TValues[K]) => {
      setValues((prev) => ({ ...prev, [field]: value }));
      if (touched[field]) {
        setErrors((prev) => ({ ...prev, [field]: validateField(field, value) }));
      }
    },
    [touched, validateField],
  );

  const blurField = useCallback(
    <K extends keyof TValues>(field: K) => {
      setTouched((prev) => ({ ...prev, [field]: true }));
      setErrors((prev) => ({ ...prev, [field]: validateField(field, values[field]) }));
    },
    [values, validateField],
  );

  const validate = useCallback((): boolean => {
    const next: Partial<Record<keyof TValues, string>> = {};
    for (const field of Object.keys(validators) as (keyof TValues)[]) {
      next[field] = validateField(field, values[field]);
    }
    setErrors(next);
    setTouched(
      Object.fromEntries((Object.keys(validators) as (keyof TValues)[]).map((f) => [f, true])) as Partial<
        Record<keyof TValues, boolean>
      >,
    );
    return Object.values(next).every((e) => e === undefined);
  }, [validators, values, validateField]);

  const reset = useCallback(() => {
    setValues(initial);
    setErrors({});
    setTouched({});
  }, [initial]);

  return { values, errors, touched, setField, blurField, validate, reset };
}