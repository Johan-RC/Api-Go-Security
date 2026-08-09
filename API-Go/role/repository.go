package role

import (
	"gorm.io/gorm"
)

// Repository encapsula el acceso a datos de roles y asignaciones RBAC (solo GORM).
type Repository struct {
	db *gorm.DB
}

// NewRepository crea el repositorio de roles.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ============================== Roles ==============================

// Create inserta un nuevo rol.
func (r *Repository) CreateRole(rl *Role) error {
	return r.db.Create(rl).Error
}

// UpdateRole actualiza un rol completo por su ID.
func (r *Repository) UpdateRole(rl *Role) error {
	return r.db.Save(rl).Error
}

// DeleteRole elimina un rol por su ID.
func (r *Repository) DeleteRole(id string) error {
	return r.db.Delete(&Role{}, "id = ?", id).Error
}

// FindRoleByID obtiene un rol por su ID.
func (r *Repository) FindRoleByID(id string) (*Role, error) {
	var rl Role
	if err := r.db.Where("id = ?", id).First(&rl).Error; err != nil {
		return nil, err
	}
	return &rl, nil
}

// FindRoleByName obtiene un rol por su nombre.
func (r *Repository) FindRoleByName(name string) (*Role, error) {
	var rl Role
	if err := r.db.Where("name = ?", name).First(&rl).Error; err != nil {
		return nil, err
	}
	return &rl, nil
}

// ListRoles obtiene roles paginados.
func (r *Repository) ListRoles(offset, limit int) ([]Role, error) {
	var roles []Role
	if err := r.db.
		Order("name ASC").
		Offset(offset).Limit(limit).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// CountRoles devuelve el total de roles.
func (r *Repository) CountRoles() (int64, error) {
	var count int64
	if err := r.db.Model(&Role{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsRoleByName indica si ya existe un rol con ese nombre.
func (r *Repository) ExistsRoleByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&Role{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// --- RoleFeature (funcionalidades asignadas a un rol) ---

// AssignFeatures crea las asignaciones de features a un rol.
func (r *Repository) AssignFeatures(items []RoleFeature) error {
	return r.db.Create(&items).Error
}

// ListFeaturesByRole obtiene las features asignadas a un rol.
func (r *Repository) ListFeaturesByRole(roleID string) ([]RoleFeature, error) {
	var items []RoleFeature
	if err := r.db.Where("role_id = ?", roleID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// RemoveFeature elimina una asignación feature-rol por su ID.
func (r *Repository) RemoveFeature(id string) error {
	return r.db.Delete(&RoleFeature{}, "id = ?", id).Error
}

// RemoveFeaturesByRole elimina todas las features de un rol.
func (r *Repository) RemoveFeaturesByRole(roleID string) error {
	return r.db.Delete(&RoleFeature{}, "role_id = ?", roleID).Error
}

// --- UserRole ------------------------ asignaciones de roles a usuarios ---

// AssignUserRole asigna un rol a un usuario.
func (r *Repository) AssignUserRole(item *UserRole) error {
	return r.db.Create(item).Error
}

// ListUserRolesByUser devuelve los roles asignados a un usuario.
func (r *Repository) ListUserRolesByUser(userID string) ([]UserRole, error) {
	var items []UserRole
	if err := r.db.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// RemoveUserRole elimina una asignación de rol por su ID.
func (r *Repository) RemoveUserRole(id string) error {
	return r.db.Delete(&UserRole{}, "id = ?", id).Error
}

// RemoveUserRolesByUser elimina todos los roles de un usuario.
func (r *Repository) RemoveUserRolesByUser(userID string) error {
	return r.db.Delete(&UserRole{}, "user_id = ?", userID).Error
}

// --- UserScopeOverride --------------------- permisos especiales por usuario ---

// CreateScopeOverride crea un override de scope.
func (r *Repository) CreateScopeOverride(item *UserScopeOverride) error {
	return r.db.Create(item).Error
}

// ListScopeOverridesByUser devuelve los overrides de un usuario.
func (r *Repository) ListScopeOverridesByUser(userID string) ([]UserScopeOverride, error) {
	var items []UserScopeOverride
	if err := r.db.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// RemoveScopeOverride elimina un override por su ID.
func (r *Repository) RemoveScopeOverride(id string) error {
	return r.db.Delete(&UserScopeOverride{}, "id = ?", id).Error
}

// RemoveScopeOverridesByUser elimina todos los overrides de un usuario.
func (r *Repository) RemoveScopeOverridesByUser(userID string) error {
	return r.db.Delete(&UserScopeOverride{}, "user_id = ?", userID).Error
}