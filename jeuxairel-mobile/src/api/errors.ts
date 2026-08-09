import axios from 'axios';
import type { ApiEnvelope } from '@/types/api';

export const commonMessages: Record<string, string> = {
  NETWORK: 'No se pudo conectar con el servidor.',
  INVALID_CREDENTIALS: 'Correo o contraseña incorrectos.',
  ACCOUNT_LOCKED: 'Cuenta bloqueada temporalmente por varios intentos fallidos.',
  USER_INACTIVE: 'Tu cuenta está inactiva. Contacta al administrador.',
  INVALID_TOKEN: 'El token es inválido.',
  TOKEN_EXPIRED: 'Tu sesión expiró. Inicia sesión de nuevo.',
  TOKEN_REVOKED: 'Tu sesión fue revocada. Inicia sesión de nuevo.',
  WEAK_PASSWORD: 'La contraseña no cumple la política mínima de seguridad.',
  VALIDATION_ERROR: 'Algunos campos enviados no son válidos.',
  FORBIDDEN: 'No tienes permisos para esta acción.',
  NOT_FOUND: 'El recurso solicitado no existe.',
  SESSION_NOT_FOUND: 'La sesión no existe o ya fue revocada.',
  INTERNAL: 'Ocurrió un error en el servidor. Intenta más tarde.',
};

/** Convierte cualquier error de axios en un mensaje legible para el usuario. */
export function getErrorMessage(err: unknown): string {
  if (axios.isAxiosError(err)) {
    if (!err.response) return 'No se pudo conectar con el servidor. Revisa tu conexión.';

    const envelope = err.response.data as ApiEnvelope | undefined;
    const code = envelope?.error?.code;

    if (code && commonMessages[code]) return commonMessages[code];
    if (err.response.status === 401) {
      return 'Tu sesión expiró o las credenciales son inválidas.';
    }
    if (err.response.status === 403) return 'No tienes permisos para hacer esta acción.';
    if (err.response.status === 404) return 'El recurso solicitado no existe.';
    if (err.response.status >= 500) return 'Ocurrió un error en el servidor. Intenta más tarde.';

    if (envelope?.message) return envelope.message;
    return 'Ocurrió un error inesperado.';
  }

  if (err instanceof Error && err.message) return err.message;
  return 'Ocurrió un error inesperado.';
}