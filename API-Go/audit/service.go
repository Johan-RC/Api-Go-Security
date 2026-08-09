package audit

import (
	"net/http"

	"github.com/jesus-ariel/api-jeussairel/shared/apperror"
)

// Errores de dominio del módulo audit.
var (
	ErrInternal = apperror.New("INTERNAL", "error interno del servidor", http.StatusInternalServerError)
)

// Service contiene la lógica de negocio de auditoría.
type Service struct {
	repo *Repository
}

// NewService crea el servicio de auditoría.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ListByUser devuelve la auditoría de login de un usuario.
func (s *Service) ListByUser(userID string) ([]AuditLogin, error) {
	entries, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, ErrInternal
	}
	return entries, nil
}

// List devuelve auditoría paginada.
func (s *Service) List(page, pageSize int) ([]AuditLogin, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	entries, err := s.repo.List(offset, pageSize)
	if err != nil {
		return nil, 0, ErrInternal
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, ErrInternal
	}
	return entries, total, nil
}