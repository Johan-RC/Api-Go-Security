package feature

import (
	"errors"
	"net/http"

	"github.com/jesus-ariel/api-jeussairel/shared/apperror"
	"gorm.io/gorm"
)

// Errores de dominio del módulo feature.
var (
	ErrModuleNotFound   = apperror.New("MODULE_NOT_FOUND", "módulo no encontrado", http.StatusNotFound)
	ErrModuleCodeTaken  = apperror.New("MODULE_CODE_TAKEN", "ya existe un módulo con ese código", http.StatusConflict)
	ErrFeatureNotFound  = apperror.New("FEATURE_NOT_FOUND", "funcionalidad no encontrada", http.StatusNotFound)
	ErrFeatureCodeTaken = apperror.New("FEATURE_CODE_TAKEN", "ya existe una funcionalidad con ese código", http.StatusConflict)
	ErrInternal         = apperror.New("INTERNAL", "error interno del servidor", http.StatusInternalServerError)
)

// Service contiene la lógica de negocio de módulos y funcionalidades.
type Service struct {
	repo *Repository
}

// NewService crea el servicio de features.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ============================== Modules ==============================

// CreateModule registra un nuevo módulo.
func (s *Service) CreateModule(req CreateModuleRequest) (*ModuleResponse, error) {
	exists, err := s.repo.ExistsModuleByCode(req.Code)
	if err != nil {
		return nil, ErrInternal
	}
	if exists {
		return nil, ErrModuleCodeTaken
	}

	m := &Module{
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		DisplayOrder: req.DisplayOrder,
		IconKey:      req.IconKey,
		IsActive:     true,
	}

	if err := s.repo.CreateModule(m); err != nil {
		return nil, ErrInternal
	}
	return toModuleResponse(m), nil
}

// GetModuleById obtiene un módulo por su ID.
func (s *Service) GetModuleById(id string) (*ModuleResponse, error) {
	m, err := s.repo.FindModuleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModuleNotFound
		}
		return nil, ErrInternal
	}
	return toModuleResponse(m), nil
}

// ListModules devuelve módulos paginados.
func (s *Service) ListModules(page, pageSize int) ([]ModuleResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	modules, err := s.repo.ListModules(offset, pageSize)
	if err != nil {
		return nil, 0, ErrInternal
	}
	total, err := s.repo.CountModules()
	if err != nil {
		return nil, 0, ErrInternal
	}

	result := make([]ModuleResponse, 0, len(modules))
	for i := range modules {
		result = append(result, *toModuleResponse(&modules[i]))
	}
	return result, total, nil
}

// UpdateModule actualiza un módulo existente.
func (s *Service) UpdateModule(id string, req UpdateModuleRequest) (*ModuleResponse, error) {
	m, err := s.repo.FindModuleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModuleNotFound
		}
		return nil, ErrInternal
	}

	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	if req.DisplayOrder != nil {
		m.DisplayOrder = *req.DisplayOrder
	}
	if req.IconKey != nil {
		m.IconKey = req.IconKey
	}
	if req.IsActive != nil {
		m.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateModule(m); err != nil {
		return nil, ErrInternal
	}
	return toModuleResponse(m), nil
}

// DeleteModule elimina un módulo por su ID.
func (s *Service) DeleteModule(id string) error {
	if _, err := s.repo.FindModuleByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrModuleNotFound
		}
		return ErrInternal
	}
	if err := s.repo.DeleteModule(id); err != nil {
		return ErrInternal
	}
	return nil
}

// ============================== Features ==============================

// CreateFeature registra una nueva funcionalidad.
func (s *Service) CreateFeature(req CreateFeatureRequest) (*FeatureResponse, error) {
	if _, err := s.repo.FindModuleByID(req.ModuleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModuleNotFound
		}
		return nil, ErrInternal
	}

	exists, err := s.repo.ExistsFeatureByCode(req.Code)
	if err != nil {
		return nil, ErrInternal
	}
	if exists {
		return nil, ErrFeatureCodeTaken
	}

	f := &Feature{
		ModuleID:    req.ModuleID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		ActionLevel: req.ActionLevel,
		IsActive:    true,
	}

	if err := s.repo.CreateFeature(f); err != nil {
		return nil, ErrInternal
	}
	return toFeatureResponse(f), nil
}

// GetFeatureById obtiene una funcionalidad por su ID.
func (s *Service) GetFeatureById(id string) (*FeatureResponse, error) {
	f, err := s.repo.FindFeatureByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFeatureNotFound
		}
		return nil, ErrInternal
	}
	return toFeatureResponse(f), nil
}

// ListFeatures devuelve funcionalidades paginadas.
func (s *Service) ListFeatures(page, pageSize int) ([]FeatureResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	features, err := s.repo.ListFeatures(offset, pageSize)
	if err != nil {
		return nil, 0, ErrInternal
	}
	total, err := s.repo.CountFeatures()
	if err != nil {
		return nil, 0, ErrInternal
	}

	result := make([]FeatureResponse, 0, len(features))
	for i := range features {
		result = append(result, *toFeatureResponse(&features[i]))
	}
	return result, total, nil
}

// ListFeaturesByModule devuelve las funcionalidades de un módulo.
func (s *Service) ListFeaturesByModule(moduleID string) ([]FeatureResponse, error) {
	if _, err := s.repo.FindModuleByID(moduleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModuleNotFound
		}
		return nil, ErrInternal
	}

	features, err := s.repo.ListFeaturesByModule(moduleID)
	if err != nil {
		return nil, ErrInternal
	}

	result := make([]FeatureResponse, 0, len(features))
	for i := range features {
		result = append(result, *toFeatureResponse(&features[i]))
	}
	return result, nil
}

// UpdateFeature actualiza una funcionalidad existente.
func (s *Service) UpdateFeature(id string, req UpdateFeatureRequest) (*FeatureResponse, error) {
	f, err := s.repo.FindFeatureByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFeatureNotFound
		}
		return nil, ErrInternal
	}

	if req.Name != nil {
		f.Name = *req.Name
	}
	if req.Description != nil {
		f.Description = req.Description
	}
	if req.ActionLevel != nil {
		f.ActionLevel = *req.ActionLevel
	}
	if req.IsActive != nil {
		f.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateFeature(f); err != nil {
		return nil, ErrInternal
	}
	return toFeatureResponse(f), nil
}

// DeleteFeature elimina una funcionalidad por su ID.
func (s *Service) DeleteFeature(id string) error {
	if _, err := s.repo.FindFeatureByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFeatureNotFound
		}
		return ErrInternal
	}
	if err := s.repo.DeleteFeature(id); err != nil {
		return ErrInternal
	}
	return nil
}

// toModuleResponse convierte un Module al DTO.
func toModuleResponse(m *Module) *ModuleResponse {
	return &ModuleResponse{
		ID:           m.ID,
		Code:         m.Code,
		Name:         m.Name,
		DisplayOrder: m.DisplayOrder,
		IsActive:     m.IsActive,
	}
}

// toFeatureResponse convierte un Feature al DTO.
func toFeatureResponse(f *Feature) *FeatureResponse {
	return &FeatureResponse{
		ID:          f.ID,
		ModuleID:    f.ModuleID,
		Code:        f.Code,
		Name:        f.Name,
		ActionLevel: f.ActionLevel,
		IsActive:    f.IsActive,
	}
}