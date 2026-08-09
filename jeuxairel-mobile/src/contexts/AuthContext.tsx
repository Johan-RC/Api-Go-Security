import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { saveTokens, clearTokens, loadTokens } from '@/api/tokenStore';
import { login as apiLogin, logout as apiLogout, me as apiMe } from '@/api/endpoints';
import { setOnSessionExpired } from '@/api/client';
import { decodeJwt } from '@/utils/jwt';
import type { ActorType, MePermission, RoleResponse } from '@/types/models';

export interface AuthUser {
  id: string;
  email: string;
  actorType: string;
}

// Jerarquía de acciones del RBAC (igual que el middleware en Go):
// quien tiene WRITE también puede READ; quien tiene APPROVE puede todo lo inferior.
const ACTION_RANK: Record<string, number> = {
  READ: 1,
  WRITE: 2,
  DELETE: 3,
  PUBLISH: 4,
  APPROVE: 5,
};

interface AuthContextValue {
  user: AuthUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  permissions: MePermission[];
  roles: RoleResponse[];
  /** True si el usuario puede hacer `action` sobre el feature `code`. */
  can: (code: string, action?: string) => boolean;
  hasRole: (name: string) => boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [permissions, setPermissions] = useState<MePermission[]>([]);
  const [roles, setRoles] = useState<RoleResponse[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Carga permisos/roles reales desde /auth/me y ajusta el usuario del JWT.
  const refreshProfile = useCallback(async () => {
    try {
      const me = await apiMe();
      setPermissions(Array.isArray(me.permissions) ? me.permissions : []);
      setRoles(Array.isArray(me.roles) ? me.roles : []);
      setUser({
        id: me.user.id,
        email: me.user.email,
        actorType: me.user.actor_type as ActorType,
      });
    } catch {
      // Si /me falla, se mantiene la sesión con los datos del JWT.
      setPermissions([]);
      setRoles([]);
    }
  }, []);

  // Al abrir la app, restaura la sesión desde el almacenamiento.
  useEffect(() => {
    (async () => {
      try {
        const tokens = await loadTokens();
        if (tokens.accessToken) {
          const payload = decodeJwt(tokens.accessToken);
          if (payload?.user_id) {
            setUser({
              id: payload.user_id as string,
              email: (payload.email as string) ?? '',
              actorType: (payload.actor_type as ActorType) ?? 'USER',
            });
            await refreshProfile();
          }
        }
      } finally {
        setIsLoading(false);
      }
    })();
  }, [refreshProfile]);

  // Si el backend no permite renovar el token, la sesión termina.
  useEffect(() => {
    setOnSessionExpired(() => {
      setUser(null);
    });
    return () => setOnSessionExpired(null);
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const tokens = await apiLogin(email, password);
    await saveTokens(tokens.access_token, tokens.refresh_token);

    const payload = decodeJwt(tokens.access_token);
    setUser({
      id: payload?.user_id ?? '',
      email: payload?.email ?? email,
      actorType: (payload?.actor_type as ActorType) ?? 'USER',
    });
    await refreshProfile();
  }, [refreshProfile]);

  const logout = useCallback(async () => {
    const { refreshToken } = await loadTokens();
    try {
      if (refreshToken) await apiLogout(refreshToken);
    } catch {
      // Aún si el servidor falla, la sesión local se cierra.
    }
    await clearTokens();
    setUser(null);
    setPermissions([]);
    setRoles([]);
  }, []);

  const can = useCallback(
    (code: string, action?: string): boolean => {
      const list = permissions ?? [];
      if (!action) {
        return list.some((p) => p.code === code);
      }
      const required = ACTION_RANK[action];
      if (!required) return false;
      return list.some((p) => p.code === code && ACTION_RANK[p.action] >= required);
    },
    [permissions],
  );

  const hasRole = useCallback(
    (name: string): boolean => (roles ?? []).some((r) => r.name === name),
    [roles],
  );

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      isAuthenticated: !!user,
      isLoading,
      permissions,
      roles,
      can,
      hasRole,
      login,
      logout,
    }),
    [user, isLoading, permissions, roles, can, hasRole, login, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth debe usarse dentro de <AuthProvider>');
  return ctx;
}