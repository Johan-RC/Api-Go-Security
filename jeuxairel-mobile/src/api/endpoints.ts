import { client } from './client';
import type { ApiEnvelope, TokenPair, ListEnvelope } from '@/types/api';
import type {
  AuditLoginResponse,
  CreateFeatureRequest,
  CreateModuleRequest,
  CreateRoleRequest,
  CreateScopeOverrideRequest,
  CreateUserRequest,
  FeatureAssignment,
  FeatureResponse,
  MeResponse,
  ModuleResponse,
  RoleResponse,
  ScopeOverrideResponse,
  SessionResponse,
  UpdateUserRequest,
  UserResponse,
  UserRoleResponse,
} from '@/types/models';

async function unwrap<T>(p: Promise<{ data: ApiEnvelope<T> }>): Promise<T> {
  const res = await p;
  return res.data.data as T;
}

// Los endpoints que devuelven listas no paginadas del backend usan el
// formato {"items": [...]}. Este helper extrae el arreglo.
async function unwrapItems<T>(p: Promise<{ data: ApiEnvelope<{ items: T[] }> }>): Promise<T[]> {
  const res = await p;
  return (res.data.data as { items: T[] }).items;
}

async function unwrapVoid(p: Promise<{ data: ApiEnvelope<undefined> }>): Promise<void> {
  await p;
}

// ---------------- Auth ----------------

export function login(email: string, password: string): Promise<TokenPair> {
  return unwrap(client.post('/auth/login', { email, password }));
}

export function logout(refreshToken: string): Promise<void> {
  return unwrapVoid(client.post('/auth/logout', { refresh_token: refreshToken }));
}

export async function register(req: CreateUserRequest): Promise<{ user: UserResponse; dev_code?: string }> {
  const res = await client.post<ApiEnvelope<{ user: UserResponse; dev_code?: string }>>('/auth/register', req);
  return res.data.data as { user: UserResponse; dev_code?: string };
}

export async function forgotPassword(email: string): Promise<{ code?: string }> {
  const res = await client.post<ApiEnvelope<{ code?: string } | undefined>>('/auth/forgot-password', { email });
  const data = res.data.data;
  return data && typeof data === 'object' ? { code: (data as { code?: string }).code } : {};
}

export function verifyEmail(email: string, code: string): Promise<void> {
  return unwrapVoid(client.post('/auth/verify-email', { email, code }));
}

export function resetPassword(code: string, newPassword: string): Promise<void> {
  return unwrapVoid(client.post('/auth/reset-password', { code, new_password: newPassword }));
}

export function me(): Promise<MeResponse> {
  return unwrap(client.get('/auth/me'));
}

// ---------------- Users ----------------

export function listUsers(page = 1, pageSize = 20): Promise<ListEnvelope<UserResponse>> {
  return unwrap(client.get('/users', { params: { page, page_size: pageSize } }));
}

export function getUser(id: string): Promise<UserResponse> {
  return unwrap(client.get(`/users/${id}`));
}

export function createUser(req: CreateUserRequest): Promise<UserResponse> {
  return unwrap(client.post('/users', req));
}

export function updateUser(id: string, req: UpdateUserRequest): Promise<UserResponse> {
  return unwrap(client.put(`/users/${id}`, req));
}

export function deleteUser(id: string): Promise<void> {
  return unwrapVoid(client.delete(`/users/${id}`));
}

// ---------------- Roles ----------------

export function listRoles(page = 1, pageSize = 20): Promise<ListEnvelope<RoleResponse>> {
  return unwrap(client.get('/roles', { params: { page, page_size: pageSize } }));
}

export function createRole(req: CreateRoleRequest): Promise<RoleResponse> {
  return unwrap(client.post('/roles', req));
}

export function getRole(id: string): Promise<RoleResponse> {
  return unwrap(client.get(`/roles/${id}`));
}

export function updateRole(id: string, req: { display_name?: string; description?: string }): Promise<RoleResponse> {
  return unwrap(client.put(`/roles/${id}`, req));
}

export function deleteRole(id: string): Promise<void> {
  return unwrapVoid(client.delete(`/roles/${id}`));
}

export function listRoleFeatures(roleId: string): Promise<FeatureAssignment[]> {
  return unwrapItems(client.get(`/roles/${roleId}/features`));
}

export function assignRoleFeatures(roleId: string, features: FeatureAssignment[]): Promise<void> {
  return unwrapVoid(client.put(`/roles/${roleId}/features`, { features }));
}

// ---------------- Asignación rol <-> usuario ----------------

export function assignRoleToUser(userId: string, roleId: string): Promise<void> {
  return unwrapVoid(client.post(`/users/${userId}/roles`, { role_id: roleId }));
}

export function listUserRoles(userId: string): Promise<UserRoleResponse[]> {
  return unwrapItems(client.get(`/users/${userId}/roles`));
}

export function removeUserRole(userId: string, userRoleId: string): Promise<void> {
  return unwrapVoid(client.delete(`/users/${userId}/roles/${userRoleId}`));
}

// ---------------- Scope overrides ----------------

export function createScopeOverride(req: CreateScopeOverrideRequest): Promise<ScopeOverrideResponse> {
  return unwrap(client.post('/scope-overrides', req));
}

export function listScopeOverrides(userId: string): Promise<ScopeOverrideResponse[]> {
  return unwrapItems(client.get(`/scope-overrides/users/${userId}`));
}

export function removeScopeOverride(id: string): Promise<void> {
  return unwrapVoid(client.delete(`/scope-overrides/${id}`));
}

// ---------------- Módulos y funcionalidades ----------------

export function listModules(page = 1, pageSize = 20): Promise<ListEnvelope<ModuleResponse>> {
  return unwrap(client.get('/modules', { params: { page, page_size: pageSize } }));
}

export function createModule(req: CreateModuleRequest): Promise<ModuleResponse> {
  return unwrap(client.post('/modules', req));
}

export function deleteModule(id: string): Promise<void> {
  return unwrapVoid(client.delete(`/modules/${id}`));
}

export function listFeaturesByModule(moduleId: string): Promise<FeatureResponse[]> {
  return unwrapItems(client.get(`/modules/${moduleId}/features`));
}

export function listFeatures(page = 1, pageSize = 100): Promise<ListEnvelope<FeatureResponse>> {
  return unwrap(client.get('/features', { params: { page, page_size: pageSize } }));
}

export function createFeature(req: CreateFeatureRequest): Promise<FeatureResponse> {
  return unwrap(client.post('/features', req));
}

export function deleteFeature(id: string): Promise<void> {
  return unwrapVoid(client.delete(`/features/${id}`));
}

// ---------------- Sesiones ----------------

export function listActiveSessions(): Promise<SessionResponse[]> {
  return unwrapItems(client.get('/sessions/active'));
}

export function revokeSession(id: string): Promise<void> {
  return unwrapVoid(client.delete(`/sessions/${id}`));
}

// ---------------- Auditoría ----------------

export function listAuditLogins(page = 1, pageSize = 20): Promise<ListEnvelope<AuditLoginResponse>> {
  return unwrap(client.get('/audit/logins', { params: { page, page_size: pageSize } }));
}

export function listAuditByUser(userId: string): Promise<AuditLoginResponse[]> {
  return unwrapItems(client.get(`/audit/logins/users/${userId}`));
}