import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { saveTokens, clearTokens, loadTokens } from './tokenStore';

// URL base de la API. Por defecto localhost:8080 (API en Go).
// Emulador Android: usa http://10.0.2.2:8080/api/v1
// Puedes sobreescribirla creando un archivo .env:
//   EXPO_PUBLIC_API_URL=http://tu-ip:8080/api/v1
export const API_BASE_URL =
  process.env.EXPO_PUBLIC_API_URL ?? 'http://localhost:8080/api/v1';

export const client = axios.create({
  baseURL: API_BASE_URL,
  headers: { 'Content-Type': 'application/json' },
  timeout: 15_000,
});

// Callback que se ejecuta cuando la sesión expira (no se pudo renovar el token).
let onSessionExpired: (() => void) | null = null;
export function setOnSessionExpired(cb: (() => void) | null) {
  onSessionExpired = cb;
}

// Evita disparar varios refrescos al mismo tiempo (single-flight).
let refreshInFlight: Promise<string | null> | null = null;

async function refreshTokens(): Promise<string | null> {
  if (refreshInFlight) return refreshInFlight;

  refreshInFlight = (async () => {
    const { refreshToken } = await loadTokens();
    if (!refreshToken) return null;

    try {
      const { data } = await axios.post<{
        data: { access_token: string; refresh_token: string };
      }>(`${API_BASE_URL}/auth/refresh`, {
        refresh_token: refreshToken,
      });
      await saveTokens(data.data.access_token, data.data.refresh_token);
      return data.data.access_token;
    } catch {
      await clearTokens();
      onSessionExpired?.();
      return null;
    }
  })().finally(() => {
    refreshInFlight = null;
  });

  return refreshInFlight;
}

type RequestConfig = InternalAxiosRequestConfig & { _retry?: boolean };

// Adjunta el access token a cada petición.
client.interceptors.request.use(async (config) => {
  const { accessToken } = await loadTokens();
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`;
  }
  return config;
});

// Si el access token expiró (401), intenta renovarlo una vez y reenvía.
client.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const original = error.config as RequestConfig | undefined;
    if (!original || original._retry || error.response?.status !== 401) {
      return Promise.reject(error);
    }

    const newToken = await refreshTokens();
    if (newToken) {
      original._retry = true;
      original.headers.Authorization = `Bearer ${newToken}`;
      return client.request(original);
    }

    return Promise.reject(error);
  },
);