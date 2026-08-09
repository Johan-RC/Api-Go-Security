export type ActorType = 'USER' | 'INSTRUCTOR' | 'LEARNER';
export type ScopeType =
  | 'GLOBAL'
  | 'TRAINING_CENTER'
  | 'AREA'
  | 'OWN_FICHAS'
  | 'OWN_SCHEDULE'
  | 'OWN_PROFILE'
  | 'OWN_FICHA_AS_LEARNER';
export type ActionLevel = 'READ' | 'WRITE' | 'DELETE' | 'PUBLISH' | 'APPROVE';

/** Permiso efectivo del usuario autenticado (viene de /auth/me). */
export interface MePermission {
  code: string;
  action: string;
  scope: string;
}

/** Respuesta de /auth/me: usuario + roles + permisos efectivos. */
export interface MeResponse {
  user: UserResponse;
  roles: RoleResponse[];
  permissions: MePermission[];
}

export interface UserResponse {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  actor_type: ActorType;
  is_active: boolean;
  last_login_at: string | null;
  created_at: string;
}

export interface CreateUserRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  actor_type: ActorType;
  actor_id?: string | null;
}

export interface CreateModuleRequest {
  code: string;
  name: string;
  description?: string;
  display_order: number;
  icon_key?: string;
}

export interface CreateFeatureRequest {
  module_id: string;
  code: string;
  name: string;
  description?: string;
  action_level: ActionLevel;
}

export interface UpdateUserRequest {
  first_name?: string;
  last_name?: string;
  is_active?: boolean;
}

export interface RoleResponse {
  id: string;
  name: string;
  display_name: string;
  description?: string | null;
  is_system_role: boolean;
}

export interface CreateRoleRequest {
  name: string;
  display_name: string;
  description?: string;
}

export interface AssignRoleRequest {
  role_id: string;
}

export interface FeatureAssignment {
  feature_id: string;
  scope_type: ScopeType;
}

export interface AssignFeaturesRequest {
  features: FeatureAssignment[];
}

export interface UserRoleResponse {
  id: string;
  role_id: string;
  role_name: string;
  display_name: string;
  scope_type: ScopeType;
  assigned_by: string;
}

export interface ScopeOverrideResponse {
  id: string;
  user_id: string;
  feature_id: string;
  scope_type: ScopeType;
  is_allowed: boolean;
  reason: string;
  granted_by: string;
  expires_at: string | null;
}

export interface CreateScopeOverrideRequest {
  user_id: string;
  feature_id: string;
  scope_type: ScopeType;
  is_allowed?: boolean;
  reason: string;
  granted_by: string;
  expires_at?: string | null;
}

export interface ModuleResponse {
  id: string;
  code: string;
  name: string;
  display_order: number;
  is_active: boolean;
}

export interface FeatureResponse {
  id: string;
  module_id: string;
  code: string;
  name: string;
  action_level: ActionLevel;
  is_active: boolean;
}

export interface SessionResponse {
  id: string;
  user_id: string;
  device_hint?: string | null;
  ip_address?: string | null;
  expires_at: string;
  is_revoked: boolean;
  created_at: string;
}

export interface AuditLoginResponse {
  id: string;
  user_id: string | null;
  email_attempted: string;
  outcome: AuditOutcome;
  ip_address?: string | null;
  user_agent?: string | null;
  attempted_at: string;
}

export type AuditOutcome =
  | 'SUCCESS'
  | 'INVALID_PASSWORD'
  | 'USER_NOT_FOUND'
  | 'ACCOUNT_LOCKED'
  | 'TOKEN_EXPIRED';