package user

import (
	"errors"
	"net/http"

	"github.com/jesus-ariel/api-jeussairel/shared/apperror"
	"github.com/jesus-ariel/api-jeussairel/shared/bcrypt"
	"github.com/jesus-ariel/api-jeussairel/shared/validation"
	"gorm.io/gorm"
)

// Errores de dominio del módulo user.
var (
	ErrNotFound     = apperror.New("USER_NOT_FOUND", "usuario no encontrado", http.StatusNotFound)
	ErrEmailTaken   = apperror.New("EMAIL_TAKEN", "ya existe un usuario con ese email", http.StatusConflict)
	ErrInvalidEmail = apperror.New("INVALID_EMAIL", "email inválido", http.StatusBadRequest)
	ErrWeakPassword = apperror.New("WEAK_PASSWORD", "la contraseña no cumple la política mínima", http.StatusBadRequest)
	ErrInvalidActor = apperror.New("INVALID_ACTOR_TYPE", "tipo de actor inválido", http.StatusBadRequest)
	ErrInternal     = apperror.New("INTERNAL", "error interno del servidor", http.StatusInternalServerError)
)

// Valid actor types.
var validActorTypes = map[string]bool{
	"USER": true, "INSTRUCTOR": true, "LEARNER": true,
}

// Service contiene la lógica de negocio de usuarios.
type Service struct {
	repo *Repository
}

// NewService crea el servicio de usuarios.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create registra un nuevo usuario con su contraseña hasheada.
func (s *Service) Create(req CreateUserRequest) (*UserResponse, error) {
	if !validation.ValidateEmail(req.Email) {
		return nil, ErrInvalidEmail
	}
	if !validation.ValidatePassword(req.Password) {
		return nil, ErrWeakPassword
	}
	if !validActorTypes[req.ActorType] {
		return nil, ErrInvalidActor
	}

	exists, err := s.repo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, ErrInternal
	}
	if exists {
		return nil, ErrEmailTaken
	}

	passwordHash, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return nil, ErrInternal
	}

	u := &User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		ActorType:    req.ActorType,
		ActorID:      req.ActorID,
		IsActive:     true,
	}

	if err := s.repo.Create(u); err != nil {
		return nil, ErrInternal
	}

	return toResponse(u), nil
}

// GetById obtiene un usuario por su ID.
func (s *Service) GetById(id string) (*UserResponse, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return toResponse(u), nil
}

// List devuelve usuarios paginados.
func (s *Service) List(page, pageSize int) ([]UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	users, err := s.repo.List(offset, pageSize)
	if err != nil {
		return nil, 0, ErrInternal
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, ErrInternal
	}

	result := make([]UserResponse, 0, len(users))
	for i := range users {
		result = append(result, *toResponse(&users[i]))
	}
	return result, total, nil
}

// Update modifica los datos de perfil de un usuario.
func (s *Service) Update(id string, req UpdateUserRequest) (*UserResponse, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	if req.FirstName != nil {
		u.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		u.LastName = *req.LastName
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}

	if err := s.repo.Update(u); err != nil {
		return nil, ErrInternal
	}
	return toResponse(u), nil
}

// Delete elimina un usuario por su ID.
func (s *Service) Delete(id string) error {
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}
	if err := s.repo.Delete(id); err != nil {
		return ErrInternal
	}
	return nil
}

// toResponse convierte un User al DTO sin datos sensibles.
func toResponse(u *User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		ActorType: u.ActorType,
		IsActive:  u.IsActive,
		LastLogin: u.LastLoginAt,
		CreatedAt: u.CreatedAt,
	}
}
