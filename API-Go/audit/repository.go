package audit

import (
	"time"

	"gorm.io/gorm"
)

// Repository encapsula el acceso a datos de auditoría (solo GORM).
type Repository struct {
	db *gorm.DB
}

// NewRepository crea el repositorio de auditoría.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create registra un evento de inicio de sesión (o intento).
func (r *Repository) Create(entry *AuditLogin) error {
	return r.db.Create(entry).Error
}

// FindByID obtiene un evento de auditoría por su ID.
func (r *Repository) FindByID(id string) (*AuditLogin, error) {
	var e AuditLogin
	if err := r.db.Where("id = ?", id).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

// ListByUser devuelve los eventos de un usuario, del más reciente al más antiguo.
func (r *Repository) ListByUser(userID string) ([]AuditLogin, error) {
	var entries []AuditLogin
	if err := r.db.
		Where("user_id = ?", userID).
		Order("attempted_at DESC").
		Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

// List obtiene eventos de auditoría paginados y ordenados por fecha.
func (r *Repository) List(offset, limit int) ([]AuditLogin, error) {
	var entries []AuditLogin
	if err := r.db.
		Order("attempted_at DESC").
		Offset(offset).Limit(limit).
		Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

// Count devuelve el total de eventos de auditoría.
func (r *Repository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&AuditLogin{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByOutcome cuenta los eventos de un resultado específico desde una fecha.
func (r *Repository) CountByOutcome(outcome string, since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&AuditLogin{}).
		Where("outcome = ? AND attempted_at >= ?", outcome, since).
		Count(&count).Error
	return count, err
}