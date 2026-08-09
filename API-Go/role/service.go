package role

import (
	"errors"
	"net/http"

	"github.com/jesus-ariel/api-jeussairel/feature"
	"github.com/jesus-ariel/api-jeussairel/shared/apperror"
	"github.com/jesus-ariel/api-jeussairel/user"
	"gorm.io/gorm"
)

// Errores de dominio del módulo role.
var (
	ErrNotFound        = apperror.New("ROLE_NOT_FOUND", "rol no encontrado", http.StatusNotFound)
	ErrNameTaken       = apperror.New("ROLE_NAME_TAKEN", "ya existe un rol con ese nombre", http.StatusConflict)
	ErrFeatureNotFound = apperror.New("FEATURE_NOT_FOUND", "funcionalidad no encontrada", http.StatusNotFound)
	ErrUserNotFound    = apperror.New("USER_NOT_FOUND", "usuario no encontrado", http.StatusNotFound)
	ErrRoleAssigned    = apperror.New("ROLE_ALREADY_ASSIGNED", "el rol ya está asignado al usuario", http.StatusConflict)
	ErrInternal        = apperror.New("INTERNAL", "error interno del servidor", http.StatusInternalServerError)
)

// Valid scope types.
var validScopes = map[string]bool{
	"GLOBAL": true, "TRAINING_CENTER": true, "AREA": true,
	"OWN_FICHAS": true, "OWN_SCHEDULE": true, "OWN_PROFILE": true,
	"OWN_FICHA_AS_LEARNER": true,
}

// Service contiene la lógica de negocio de roles y RBAC.
type Service struct {
	repo     *Repository
	features *feature.Repository
	users    *user.Repository
}

// NewService crea el servicio de roles.
func NewService(repo *Repository, features *feature.Repository, users *user.Repository) *Service {
	return &Service{repo: repo, features: features, users: users}
}

// ============================== CRUD Roles ==============================

// Create registra un nuevo rol.
func (s *Service) Create(req CreateRoleRequest) (*RoleResponse, error) {
	exists, err := s.repo.ExistsRoleByName(req.Name)
	if err != nil {
		return nil, ErrInternal
	}
	if exists {
		return nil, ErrNameTaken
	}

	rl := &Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
	}

	if err := s.repo.CreateRole(rl); err != nil {
		return nil, ErrInternal
	}
	return toRoleResponse(rl), nil
}

// GetById obtiene un rol por su ID.
func (s *Service) GetById(id string) (*RoleResponse, error) {
	rl, err := s.repo.FindRoleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return toRoleResponse(rl), nil
}

// List devuelve roles paginados.
func (s *Service) List(page, pageSize int) ([]RoleResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	roles, err := s.repo.ListRoles(offset, pageSize)
	if err != nil {
		return nil, 0, ErrInternal
	}
	total, err := s.repo.CountRoles()
	if err != nil {
		return nil, 0, ErrInternal
	}

	result := make([]RoleResponse, 0, len(roles))
	for i := range roles {
		result = append(result, *toRoleResponse(&roles[i]))
	}
	return result, total, nil
}

// Update modifica un rol existente.
func (s *Service) Update(id string, req UpdateRoleRequest) (*RoleResponse, error) {
	rl, err := s.repo.FindRoleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	if req.DisplayName != nil {
		rl.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		rl.Description = req.Description
	}

	if err := s.repo.UpdateRole(rl); err != nil {
		return nil, ErrInternal
	}
	return toRoleResponse(rl), nil
}

// Delete elimina un rol por su ID.
func (s *Service) Delete(id string) error {
	if _, err := s.repo.FindRoleByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}
	if err := s.repo.DeleteRole(id); err != nil {
		return ErrInternal
	}
	return nil
}

// ============================== Asignación de features ==============================

// AssignFeatures reemplaza las funcionalidades asignadas a un rol.
func (s *Service) AssignFeatures(roleID string, req AssignFeaturesRequest) error {
	if _, err := s.repo.FindRoleByID(roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}

	items := make([]RoleFeature, 0, len(req.Features))
	for _, f := range req.Features {
		if !validScopes[f.ScopeType] {
			return apperror.New("INVALID_SCOPE_TYPE", "tipo de ámbito inválido", http.StatusBadRequest)
		}
		if _, err := s.features.FindFeatureByID(f.FeatureID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrFeatureNotFound
			}
			return ErrInternal
		}
		items = append(items, RoleFeature{
			RoleID:    roleID,
			FeatureID: f.FeatureID,
			ScopeType: f.ScopeType,
		})
	}

	// Reasignación: se eliminan las actuales y se insertan las nuevas.
	if err := s.repo.RemoveFeaturesByRole(roleID); err != nil {
		return ErrInternal
	}
	if len(items) == 0 {
		return nil
	}
	if err := s.repo.AssignFeatures(items); err != nil {
		return ErrInternal
	}
	return nil
}

// ListFeatures devuelve las funcionalidades asignadas a un rol.
func (s *Service) ListFeatures(roleID string) ([]RoleFeature, error) {
	if _, err := s.repo.FindRoleByID(roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return s.repo.ListFeaturesByRole(roleID)
}

// RemoveFeature quita una funcionalidad de un rol.
func (s *Service) RemoveFeature(roleFeatureID string) error {
	if err := s.repo.RemoveFeature(roleFeatureID); err != nil {
		return ErrInternal
	}
	return nil
}

// ============================== Asignación de roles a usuarios ==============================

// AssignRoleToUser asigna un rol a un usuario.
func (s *Service) AssignRoleToUser(userID, roleID, assignedBy string) error {
	if _, err := s.users.FindByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return ErrInternal
	}
	if _, err := s.repo.FindRoleByID(roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}

	existing, err := s.repo.ListUserRolesByUser(userID)
	if err != nil {
		return ErrInternal
	}
	for _, ur := range existing {
		if ur.RoleID == roleID {
			return ErrRoleAssigned
		}
	}

	item := &UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: assignedBy,
	}

	if err := s.repo.AssignUserRole(item); err != nil {
		return ErrInternal
	}
	return nil
}

// ListRolesByUser devuelve los roles asignados a un usuario.
func (s *Service) ListRolesByUser(userID string) ([]UserRole, error) {
	if _, err := s.users.FindByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}
	return s.repo.ListUserRolesByUser(userID)
}

// RemoveRoleFromUser quita un rol de un usuario por el ID de la asignación.
func (s *Service) RemoveRoleFromUser(userRoleID string) error {
	if err := s.repo.RemoveUserRole(userRoleID); err != nil {
		return ErrInternal
	}
	return nil
}

// ============================== Overrides de scope ==============================

// CreateScopeOverride crea un permiso especial para un usuario.
func (s *Service) CreateScopeOverride(req CreateScopeOverrideRequest) (*UserScopeOverride, error) {
	if !validScopes[req.ScopeType] {
		return nil, apperror.New("INVALID_SCOPE_TYPE", "tipo de ámbito inválido", http.StatusBadRequest)
	}
	if _, err := s.features.FindFeatureByID(req.FeatureID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFeatureNotFound
		}
		return nil, ErrInternal
	}

	item := &UserScopeOverride{
		UserID:    req.UserID,
		FeatureID: req.FeatureID,
		ScopeType: req.ScopeType,
		IsAllowed: true,
		Reason:    req.Reason,
		GrantedBy: req.GrantedBy,
		ExpiresAt: req.ExpiresAt,
	}
	if req.IsAllowed != nil {
		item.IsAllowed = *req.IsAllowed
	}

	if err := s.repo.CreateScopeOverride(item); err != nil {
		return nil, ErrInternal
	}
	return item, nil
}

// ListScopeOverrides devuelve los overrides de un usuario.
func (s *Service) ListScopeOverrides(userID string) ([]UserScopeOverride, error) {
	return s.repo.ListScopeOverridesByUser(userID)
}

// RemoveScopeOverride elimina un override por su ID.
func (s *Service) RemoveScopeOverride(id string) error {
	if err := s.repo.RemoveScopeOverride(id); err != nil {
		return ErrInternal
	}
	return nil
}

// toRoleResponse convierte un Role al DTO.
func toRoleResponse(rl *Role) *RoleResponse {
	return &RoleResponse{
		ID:           rl.ID,
		Name:         rl.Name,
		DisplayName:  rl.DisplayName,
		Description:  rl.Description,
		IsSystemRole: rl.IsSystemRole,
	}
}