export interface ApiEnvelope<T = unknown> {
  success: boolean;
  message: string;
  data?: T;
  error?: { code: string };
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
}

export interface ListEnvelope<T> {
  items: T[];
  total: number;
  page: number;
}