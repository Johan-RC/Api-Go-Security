package role

import "time"

// CreateRoleRequest es la entrada para crear un rol.
type CreateRoleRequest struct {
	Name        string  `json:"name" binding:"required,max=50"`
	DisplayName string  `json:"display_name" binding:"required,max=100"`
	Description *string `json:"description"`
}

// UpdateRoleRequest permite actualizar los campos editables de un rol.
type UpdateRoleRequest struct {
	DisplayName *string `json:"display_name"`
	Description *string `json:"description"`
}

// AssignFeaturesRequest asigna funcionalidades a un rol.
type AssignFeaturesRequest struct {
	Features []FeatureAssignment `json:"features" binding:"required"`
}

// FeatureAssignment relaciona una funcionalidad con un ámbito.
type FeatureAssignment struct {
	FeatureID string `json:"feature_id" binding:"required"`
	ScopeType string `json:"scope_type" binding:"required,oneof=GLOBAL TRAINING_CENTER AREA OWN_FICHAS OWN_SCHEDULE OWN_PROFILE OWN_FICHA_AS_LEARNER"`
}

// CreateScopeOverrideRequest crea un permiso especial de scope para un usuario.
type CreateScopeOverrideRequest struct {
	UserID    string     `json:"user_id" binding:"required"`
	FeatureID string     `json:"feature_id" binding:"required"`
	ScopeType string     `json:"scope_type" binding:"required,oneof=GLOBAL TRAINING_CENTER AREA OWN_FICHAS OWN_SCHEDULE OWN_PROFILE OWN_FICHA_AS_LEARNER"`
	IsAllowed *bool      `json:"is_allowed"`
	Reason    string     `json:"reason" binding:"required"`
	GrantedBy string     `json:"granted_by" binding:"required"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// AssignRoleRequest asigna un rol a un usuario.
type AssignRoleRequest struct {
	RoleID string `json:"role_id" binding:"required"`
}

// RoleResponse es la salida de un rol.
type RoleResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	DisplayName  string  `json:"display_name"`
	Description  *string `json:"description"`
	IsSystemRole bool    `json:"is_system_role"`
}

// ListResponse es la salida paginada de roles.
type ListResponse struct {
	Items []RoleResponse `json:"items"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
}

// ListFeaturesResponse es la salida de las funcionalidades asignadas a un rol.
type ListFeaturesResponse struct {
	Items []FeatureAssignment `json:"items"`
}

// UserRoleResponse es la salida de un rol asignado a un usuario.
type UserRoleResponse struct {
	ID         string `json:"id"`
	RoleID     string `json:"role_id"`
	RoleName   string `json:"role_name"`
	DisplayName string `json:"display_name"`
	ScopeType string `json:"scope_type"`
	AssignedBy string `json:"assigned_by"`
}

// ScopeOverrideResponse es la salida de un override de scope.
type ScopeOverrideResponse struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	FeatureID  string  `json:"feature_id"`
	ScopeType  string  `json:"scope_type"`
	IsAllowed  bool    `json:"is_allowed"`
	Reason     string  `json:"reason"`
	GrantedBy  string  `json:"granted_by"`
	ExpiresAt  *time.Time `json:"expires_at"`
}