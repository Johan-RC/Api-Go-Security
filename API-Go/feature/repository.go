package feature

import (
	"gorm.io/gorm"
)

// Repository encapsula el acceso a datos de módulos y funcionalidades (solo GORM).
type Repository struct {
	db *gorm.DB
}

// NewRepository crea el repositorio de features.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ============================== Modules ==============================

// CreateModule inserta un nuevo módulo.
func (r *Repository) CreateModule(m *Module) error {
	return r.db.Create(m).Error
}

// UpdateModule actualiza un módulo completo por su ID.
func (r *Repository) UpdateModule(m *Module) error {
	return r.db.Save(m).Error
}

// DeleteModule elimina un módulo por su ID.
func (r *Repository) DeleteModule(id string) error {
	return r.db.Delete(&Module{}, "id = ?", id).Error
}

// FindModuleByID obtiene un módulo por su ID.
func (r *Repository) FindModuleByID(id string) (*Module, error) {
	var m Module
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// FindModuleByCode obtiene un módulo por su código.
func (r *Repository) FindModuleByCode(code string) (*Module, error) {
	var m Module
	if err := r.db.Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// ListModules obtiene módulos paginados (los inactivos incluidos).
func (r *Repository) ListModules(offset, limit int) ([]Module, error) {
	var modules []Module
	if err := r.db.
		Order("display_order ASC").
		Offset(offset).Limit(limit).
		Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

// CountModules devuelve el total de módulos.
func (r *Repository) CountModules() (int64, error) {
	var count int64
	if err := r.db.Model(&Module{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsModuleByCode indica si ya existe un módulo con ese código.
func (r *Repository) ExistsModuleByCode(code string) (bool, error) {
	var count int64
	err := r.db.Model(&Module{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

// ============================== Features ==============================

// CreateFeature inserta una nueva funcionalidad.
func (r *Repository) CreateFeature(f *Feature) error {
	return r.db.Create(f).Error
}

// UpdateFeature actualiza una funcionalidad completa por su ID.
func (r *Repository) UpdateFeature(f *Feature) error {
	return r.db.Save(f).Error
}

// DeleteFeature elimina una funcionalidad por su ID.
func (r *Repository) DeleteFeature(id string) error {
	return r.db.Delete(&Feature{}, "id = ?", id).Error
}

// FindFeatureByID obtiene una funcionalidad por su ID.
func (r *Repository) FindFeatureByID(id string) (*Feature, error) {
	var f Feature
	if err := r.db.Where("id = ?", id).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

// FindFeatureByCode obtiene una funcionalidad por su código.
func (r *Repository) FindFeatureByCode(code string) (*Feature, error) {
	var f Feature
	if err := r.db.Where("code = ?", code).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

// ListFeatures obtiene funcionalidades paginadas.
func (r *Repository) ListFeatures(offset, limit int) ([]Feature, error) {
	var features []Feature
	if err := r.db.
		Order("name ASC").
		Offset(offset).Limit(limit).
		Find(&features).Error; err != nil {
		return nil, err
	}
	return features, nil
}

// ListFeaturesByModule obtiene las funcionalidades de un módulo.
func (r *Repository) ListFeaturesByModule(moduleID string) ([]Feature, error) {
	var features []Feature
	if err := r.db.Where("module_id = ?", moduleID).Find(&features).Error; err != nil {
		return nil, err
	}
	return features, nil
}

// CountFeatures devuelve el total de funcionalidades.
func (r *Repository) CountFeatures() (int64, error) {
	var count int64
	if err := r.db.Model(&Feature{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsFeatureByCode indica si ya existe una funcionalidad con ese código.
func (r *Repository) ExistsFeatureByCode(code string) (bool, error) {
	var count int64
	err := r.db.Model(&Feature{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}