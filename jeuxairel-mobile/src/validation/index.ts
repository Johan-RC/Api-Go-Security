export type FieldValue = unknown;

type ValidatorFn = (value: string) => string | undefined;

/** Ejecuta varias reglas en orden y devuelve el primer error. */
export function seq(...rules: ValidatorFn[]): ValidatorFn {
  return (value: string) => {
    for (const rule of rules) {
      const err = rule(value);
      if (err) return err;
    }
    return undefined;
  };
}

const EMAIL_REGEX =
  /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,63}$/;

const PASS_UPPER = /[A-Z]/;
const PASS_LOWER = /[a-z]/;
const PASS_DIGIT = /[0-9]/;
const PASS_SPECIAL = /[^A-Za-z0-9]/;

// Partes locales que son sistemas/listas de distribución, no bandejas personales.
const RESERVED_LOCAL_PARTS = new Set([
  'abuse', 'bounce', 'contact', 'devnull', 'info', 'list',
  'mailer-daemon', 'no-reply', 'noreply', 'postmaster', 'root', 'sales',
  'security', 'support', 'team', 'webmaster',
]);

export function isEmail(value: string): boolean {
  const email = value.trim().toLowerCase();
  if (email.length > 254 || !email.includes('@')) return false;
  if (!EMAIL_REGEX.test(email)) return false;
  const local = email.split('@')[0];
  return !RESERVED_LOCAL_PARTS.has(local);
}

export function isRequired(value: FieldValue): boolean {
  if (value === null || value === undefined) return false;
  if (typeof value === 'string') return value.trim().length > 0;
  if (Array.isArray(value)) return value.length > 0;
  return true;
}

export function minLength(value: string, min: number): boolean {
  return value.trim().length >= min;
}

export function maxLength(value: string, max: number): boolean {
  return value.trim().length <= max;
}

export function exactLength(value: string, len: number): boolean {
  return value.trim().length === len;
}

export function isOneOf(value: string, allowed: readonly string[]): boolean {
  return allowed.includes(value);
}

/** Solo letras, espacios, tildes, ñ/Ñ, apóstrofe y guion. */
export function isName(value: string): boolean {
  const name = value.trim();
  if (name.length === 0) return false;
  return /^[\p{L}\s'’-]+$/u.test(name);
}

export interface PasswordCheck {
  length: boolean;
  upper: boolean;
  lower: boolean;
  digit: boolean;
  special: boolean;
}

/** Devuelve qué criterios cumple la contraseña (para UI de política). */
export function checkPassword(value: string): PasswordCheck {
  const password = value.trim();
  return {
    length: password.length >= 8,
    upper: PASS_UPPER.test(password),
    lower: PASS_LOWER.test(password),
    digit: PASS_DIGIT.test(password),
    special: PASS_SPECIAL.test(password),
  };
}

/** Valida la política completa de contraseña. */
export function isStrongPassword(value: string): boolean {
  const c = checkPassword(value);
  return c.length && c.upper && c.lower && c.digit && c.special;
}

/** Mensajes en español por tipo de validación. */
export const messages = {
  required: (label: string) => `${label} es obligatorio.`,
  email: 'Ingresa un correo electrónico válido.',
  emailReserved: 'No se permiten correos de sistema (ej. no-reply@…).',
  min: (label: string, min: number) => `${label} debe tener al menos ${min} caracteres.`,
  max: (label: string, max: number) => `${label} no puede superar ${max} caracteres.`,
  name: 'Solo se permiten letras y espacios.',
  weak: 'Debe tener 8+ caracteres, mayúscula, minúscula, número y símbolo.',
  oneOf: (label: string) => `${label} no es un valor permitido.`,
  match: (label: string) => `${label} no coincide.`,
  code: 'El código debe tener exactamente 6 dígitos.',
};

// ---- Fábricas de reglas (devuelven validadores para useForm) ----

export const ruleRequired = (label: string): ValidatorFn => (value) =>
  isRequired(value) ? undefined : messages.required(label);

export const ruleEmail: ValidatorFn = (value) => {
  if (!isRequired(value)) return messages.required('El correo');
  return isEmail(value) ? undefined : messages.email;
};

export const ruleName = (label: string): ValidatorFn => (value) => {
  if (!isRequired(value)) return messages.required(label);
  return isName(value) ? undefined : messages.name;
};

export const ruleMinLength = (label: string, min: number): ValidatorFn => (value) =>
  minLength(value, min) ? undefined : messages.min(label, min);

export const ruleMaxLength = (label: string, max: number): ValidatorFn => (value) =>
  maxLength(value, max) ? undefined : messages.max(label, max);

export const ruleStrongPassword: ValidatorFn = (value) => {
  if (!isRequired(value)) return messages.required('La contraseña');
  return isStrongPassword(value) ? undefined : messages.weak;
};

export const ruleConfirm = (label: string, otherValue: () => string): ValidatorFn => (value) =>
  isRequired(value) ? (value === otherValue() ? undefined : messages.match(label)) : messages.required(label);

export const ruleSixDigitCode: ValidatorFn = (value) => {
  if (!isRequired(value)) return messages.required('El código');
  return /^\d{6}$/.test(value.trim()) ? undefined : messages.code;
};