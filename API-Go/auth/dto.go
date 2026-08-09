package auth

import (
	"github.com/jesus-ariel/api-jeussairel/role"
	"github.com/jesus-ariel/api-jeussairel/user"
)

// RegisterRequest es la entrada del registro público (sin autenticación).
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	ActorType string `json:"actor_type" binding:"required,oneof=USER INSTRUCTOR LEARNER"`
}

// LoginRequest es la entrada del endpoint de login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest es la entrada del endpoint de refresh token.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest es la entrada del endpoint de logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// TokenResponse es la salida de autenticación exitosa.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ForgotPasswordRequest solicita la recuperación de contraseña.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest cambia la contraseña con el código de recuperación.
type ResetPasswordRequest struct {
	Code        string `json:"code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// VerifyEmailRequest confirma la cuenta con el código de verificación.
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// Permission es un feature al que tiene acceso el usuario logueado.
type Permission struct {
	Code    string `json:"code"`
	Action  string `json:"action"`
	Scope   string `json:"scope"`
}

// MeResponse es la info de sesión del usuario autenticado.
type MeResponse struct {
	User        *user.UserResponse `json:"user"`
	Roles       []role.RoleResponse `json:"roles"`
	Permissions []Permission       `json:"permissions"`
}
