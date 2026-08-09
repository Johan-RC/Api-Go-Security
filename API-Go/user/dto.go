package user

import (
	"time"
)

// CreateUserRequest es la entrada para crear un usuario.
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	ActorType string `json:"actor_type" binding:"required,oneof=USER INSTRUCTOR LEARNER"`
	ActorID   *string `json:"actor_id"`
}

// UpdateUserRequest permite actualizar datos de perfil.
type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	IsActive  *bool   `json:"is_active"`
}

// UserResponse es la salida de un usuario sin datos sensibles.
type UserResponse struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	ActorType  string     `json:"actor_type"`
	IsActive   bool       `json:"is_active"`
	LastLogin  *time.Time `json:"last_login_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ListResponse es la salida paginada de listados.
type ListResponse struct {
	Items []UserResponse `json:"items"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
}
