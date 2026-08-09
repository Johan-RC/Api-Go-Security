package session

import (
	"errors"
	"net/http"

	"github.com/jesus-ariel/api-jeussairel/shared/apperror"
	"gorm.io/gorm"
)

// Errores de dominio del módulo session.
var (
	ErrNotFound = apperror.New("SESSION_NOT_FOUND", "sesión no encontrada", http.StatusNotFound)
	ErrInternal = apperror.New("INTERNAL", "error interno del servidor", http.StatusInternalServerError)
)

// Service contiene la lógica de negocio de sesiones.
type Service struct {
	repo *Repository
}

// NewService crea el servicio de sesiones.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ListActiveByUser devuelve las sesiones activas de un usuario.
func (s *Service) ListActiveByUser(userID string) ([]RefreshToken, error) {
	tokens, err := s.repo.ListActiveByUser(userID)
	if err != nil {
		return nil, ErrInternal
	}
	return tokens, nil
}

// RevokeById revoca una sesión por su ID.
func (s *Service) RevokeById(id string) error {
	if err := s.repo.RevokeRefreshToken(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}
	return nil
}