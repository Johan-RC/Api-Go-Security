/** Formatea una fecha ISO a un texto amigable en español. */
export function formatDate(value: string | null | undefined): string {
  if (!value) return '—';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '—';
  return date.toLocaleDateString('es-ES', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  });
}

export function formatDateTime(value: string | null | undefined): string {
  if (!value) return '—';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '—';
  return date.toLocaleString('es-ES', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

/** Capitaliza la primera letra y baja el resto. */
export function capitalize(value: string): string {
  if (!value) return value;
  return value.charAt(0).toUpperCase() + value.slice(1).toLowerCase();
}

/** Reemplaza guiones bajos por espacios y capitaliza. */
export function humanize(value: string): string {
  return capitalize(value.replace(/_/g, ' '));
}

/** Obtiene las iniciales de un nombre (máx 2). */
export function initials(firstName: string, lastName?: string): string {
  const a = firstName?.charAt(0) ?? '';
  const b = lastName?.charAt(0) ?? '';
  return `${a}${b}`.toUpperCase();
}