package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jesus-ariel/api-jeussairel/audit"
	"github.com/jesus-ariel/api-jeussairel/config"
	"github.com/jesus-ariel/api-jeussairel/role"
	"github.com/jesus-ariel/api-jeussairel/session"
	"github.com/jesus-ariel/api-jeussairel/shared/apperror"
	"github.com/jesus-ariel/api-jeussairel/shared/bcrypt"
	"github.com/jesus-ariel/api-jeussairel/shared/jwt"
	"github.com/jesus-ariel/api-jeussairel/shared/validation"
	"github.com/jesus-ariel/api-jeussairel/user"
	"gorm.io/gorm"
)

// Códigos de error de dominio del módulo auth.
var (
	ErrInvalidCredentials = apperror.New("INVALID_CREDENTIALS", "credenciales inválidas", http.StatusUnauthorized)
	ErrAccountLocked      = apperror.New("ACCOUNT_LOCKED", "cuenta bloqueada temporalmente", http.StatusUnauthorized)
	ErrInactiveUser       = apperror.New("USER_INACTIVE", "usuario inactivo", http.StatusForbidden)
	ErrInvalidToken       = apperror.New("INVALID_TOKEN", "token inválido", http.StatusUnauthorized)
	ErrTokenExpired       = apperror.New("TOKEN_EXPIRED", "token expirado", http.StatusUnauthorized)
	ErrTokenRevoked       = apperror.New("TOKEN_REVOKED", "token revocado", http.StatusUnauthorized)
	ErrWeakPassword       = apperror.New("WEAK_PASSWORD", "la contraseña no cumple la política mínima", http.StatusBadRequest)
	ErrEmailTaken         = apperror.New("EMAIL_TAKEN", "ya existe un usuario con ese email", http.StatusConflict)
	ErrInternal           = apperror.New("INTERNAL", "error interno del servidor", http.StatusInternalServerError)
	ErrEmailNotSent       = apperror.New("EMAIL_NOT_SENT", "no se pudo enviar el código", http.StatusInternalServerError)
	ErrInvalidCode        = apperror.New("INVALID_CODE", "el código de verificación no es válido", http.StatusBadRequest)
	ErrCodeExpired        = apperror.New("CODE_EXPIRED", "el código de verificación caducó", http.StatusBadRequest)
	ErrInvalidName        = apperror.New("INVALID_NAME", "el nombre solo puede contener letras y espacios", http.StatusBadRequest)
)

// Constantes de intentos y resultados de auditoría.
const (
	MaxFailedAttempts = 5
	LockDuration      = 15 * time.Minute
)

// Service contiene la lógica de negocio de autenticación.
type Service struct {
	repo *Repository
	cfg  *config.Config
}

// NewService crea el servicio de autenticación.
func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

// Register crea una cuenta desde el registro público (sin autenticación).
// La cuenta queda inactiva hasta que confirme un código de verificación por email.
// En desarrollo sin SMTP, devuelve el código en el segundo valor.
func (s *Service) Register(req RegisterRequest, ipAddress string) (*user.UserResponse, string, error) {
	if !validation.ValidateEmail(req.Email) {
		return nil, "", ErrInvalidCredentials
	}
	if !validation.ValidatePassword(req.Password) {
		return nil, "", ErrWeakPassword
	}
	if !validation.ValidateName(req.FirstName, 100) || !validation.ValidateName(req.LastName, 100) {
		return nil, "", ErrInvalidName
	}

	exists, err := s.repo.Users.ExistsByEmail(req.Email)
	if err != nil {
		return nil, "", ErrInternal
	}
	if exists {
		return nil, "", ErrEmailTaken
	}

	passwordHash, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return nil, "", ErrInternal
	}

	u := &user.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		ActorType:    req.ActorType,
		IsActive:     false,
	}

	if err := s.repo.Users.Create(u); err != nil {
		return nil, "", ErrInternal
	}

	code, err := generateResetCode()
	if err != nil {
		return nil, "", ErrInternal
	}

	verification := &session.EmailVerificationCode{
		UserID:    u.ID,
		CodeHash:  hashToken(code),
		ExpiresAt: time.Now().Add(15 * time.Minute),
		IPAddress: &ipAddress,
	}
	if err := s.repo.Sessions.CreateEmailVerification(verification); err != nil {
		return nil, "", ErrInternal
	}

	smtpCfg := s.cfg.SMTP()
	if smtpCfg.Enabled() {
		body := fmt.Sprintf(
			"Hola %s,\n\nBienvenido a Jeux Airel.\n\nTu código de verificación de cuenta es:\n\n%s\n\nCaduca en 15 minutos.",
			u.FirstName, code,
		)
		if err := smtpCfg.Send(u.Email, "Código de verificación - Jeux Airel", body); err != nil {
			slog.Error("error enviando email de verificación", "user_id", u.ID, "error", err)
			return nil, "", ErrEmailNotSent
		}
		return toUserResponse(u), "", nil
	}

	slog.Info("verificación de email generada (SMTP no configurado)", "user_id", u.ID, "code", code)
	return toUserResponse(u), code, nil
}

// VerifyEmail confirma el código de verificación y activa la cuenta.
func (s *Service) VerifyEmail(req VerifyEmailRequest, ipAddress string) error {
	u, err := s.repo.Users.FindByEmail(req.Email)
	if err != nil {
		return ErrInvalidCode
	}
	if u.IsActive {
		return nil
	}

	hash := hashToken(req.Code)
	verification, err := s.repo.Sessions.FindEmailVerificationByHash(hash)
	if err != nil || verification.UserID != u.ID || verification.IsUsed {
		return ErrInvalidCode
	}
	if verification.ExpiresAt.Before(time.Now()) {
		return ErrCodeExpired
	}

	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		sessRepo := session.NewRepository(tx)
		userRepo := user.NewRepository(tx)

		if err := userRepo.SetActive(u.ID, true); err != nil {
			return err
		}
		return sessRepo.MarkEmailVerificationUsed(verification.ID)
	})
}

func toUserResponse(u *user.User) *user.UserResponse {
	return &user.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		ActorType: u.ActorType,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}

// Me devuelve la información del usuario autenticado junto con sus roles
// y los permisos efectivos que alimentan la navegación y acciones de la UI.
func (s *Service) Me(userID string) (*MeResponse, error) {
	u, err := s.repo.Users.FindByID(userID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	roles, err := s.repo.ListRolesByUser(userID)
	if err != nil {
		return nil, ErrInternal
	}

	permissions, err := s.repo.ListPermissionsByUser(userID)
	if err != nil {
		return nil, ErrInternal
	}

	roleResponses := make([]role.RoleResponse, 0, len(roles))
	for _, r := range roles {
		roleResponses = append(roleResponses, role.RoleResponse{
			ID:           r.ID,
			Name:         r.Name,
			DisplayName:  r.DisplayName,
			Description:  r.Description,
			IsSystemRole: r.IsSystemRole,
		})
	}

	// Garantiza que roles y permisos siempre sean arreglos en JSON (nunca null),
	// para que el front no reciba `permissions: null` y puedan renderizar vacío.
	if permissions == nil {
		permissions = []Permission{}
	}
	if roleResponses == nil {
		roleResponses = []role.RoleResponse{}
	}

	return &MeResponse{
		User:        toUserResponse(u),
		Roles:       roleResponses,
		Permissions: permissions,
	}, nil
}

// Login autentica a un usuario y emite el par de tokens.
func (s *Service) Login(req LoginRequest, ipAddress, userAgent string) (*TokenResponse, error) {
	u, err := s.repo.Users.FindByEmail(req.Email)
	if err != nil {
		s.recordAudit(nil, req.Email, "USER_NOT_FOUND", ipAddress, userAgent)
		return nil, ErrInvalidCredentials
	}

	// Cuenta bloqueada.
	if u.LockedUntil != nil && u.LockedUntil.After(time.Now()) {
		s.recordAudit(&u.ID, req.Email, "ACCOUNT_LOCKED", ipAddress, userAgent)
		return nil, ErrAccountLocked
	}

	// Usuario inactivo.
	if !u.IsActive {
		s.recordAudit(&u.ID, req.Email, "ACCOUNT_LOCKED", ipAddress, userAgent)
		return nil, ErrInactiveUser
	}

	// Contraseña incorrecta.
	if !bcrypt.ComparePassword(req.Password, u.PasswordHash) {
		_ = s.repo.Users.IncrementFailedAttempts(u.ID)

		// Bloqueo tras demasiados intentos fallidos.
		if u.FailedAttempts+1 >= MaxFailedAttempts {
			lockedUntil := time.Now().Add(LockDuration)
			_ = s.repo.Users.SetLocked(u.ID, &lockedUntil)
			// Revocar sesiones vigentes para impedir refrescos durante el bloqueo.
			_ = s.repo.Sessions.RevokeAllByUser(u.ID)
			s.recordAudit(&u.ID, req.Email, "ACCOUNT_LOCKED", ipAddress, userAgent)
			return nil, ErrAccountLocked
		}

		s.recordAudit(&u.ID, req.Email, "INVALID_PASSWORD", ipAddress, userAgent)
		return nil, ErrInvalidCredentials
	}

	// Login exitoso: se resetean los contadores.
	now := time.Now()
	_ = s.repo.Users.ResetFailedAttempts(u.ID)
	_ = s.repo.Users.SetLocked(u.ID, nil)
	_ = s.repo.Users.UpdateLastLogin(u.ID, now)

	tokens, err := s.issueTokens(u, ipAddress)
	if err != nil {
		return nil, err
	}

	s.recordAudit(&u.ID, req.Email, "SUCCESS", ipAddress, userAgent)
	return tokens, nil
}

// Refresh valida un refresh token y emite un nuevo par de tokens (rotación).
func (s *Service) Refresh(req RefreshRequest, ipAddress string) (*TokenResponse, error) {
	hash := hashToken(req.RefreshToken)

	rt, err := s.repo.Sessions.FindRefreshTokenByHash(hash)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if rt.IsRevoked {
		return nil, ErrTokenRevoked
	}
	if rt.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	u, err := s.repo.Users.FindByID(rt.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if u.LockedUntil != nil && u.LockedUntil.After(time.Now()) {
		return nil, ErrAccountLocked
	}
	if !u.IsActive {
		return nil, ErrInactiveUser
	}

	// El token antiguo queda revocado (rotación).
	_ = s.repo.Sessions.RevokeRefreshToken(rt.ID)

	return s.issueTokens(u, ipAddress)
}

// Logout revoca el refresh token del usuario.
func (s *Service) Logout(req LogoutRequest) error {
	hash := hashToken(req.RefreshToken)

	rt, err := s.repo.Sessions.FindRefreshTokenByHash(hash)
	if err != nil {
		// Un token inexistente ya es "cerrado"; no se expone error.
		return nil
	}

	return s.repo.Sessions.RevokeRefreshToken(rt.ID)
}

// ForgotPassword genera una solicitud de reseteo de contraseña.
// Por seguridad, la respuesta es la misma exista o no el usuario.
// En desarrollo, sin SMTP configurado, devuelve el código en el segundo valor.
func (s *Service) ForgotPassword(req ForgotPasswordRequest, ipAddress string) (string, error) {
	u, err := s.repo.Users.FindByEmail(req.Email)
	if err != nil {
		// No se confirma la existencia del email.
		return "", nil
	}

	code, err := generateResetCode()
	if err != nil {
		return "", err
	}

	// Invalida solicitudes previas pendientes del usuario y crea la nueva.
	if err := s.repo.Sessions.InvalidatePasswordResetsByUser(u.ID); err != nil {
		return "", err
	}

	reset := &session.PasswordResetRequest{
		UserID:    u.ID,
		TokenHash: hashToken(code),
		ExpiresAt: time.Now().Add(15 * time.Minute),
		IPAddress: &ipAddress,
	}

	if err := s.repo.Sessions.CreatePasswordReset(reset); err != nil {
		return "", err
	}

	// Envío del código: por correo si SMTP está configurado; si no, se loguea (desarrollo).
	smtpCfg := s.cfg.SMTP()
	if smtpCfg.Enabled() {
		body := fmt.Sprintf(
			"Hola %s,\n\nRecibiste este correo porque solicitaste restablecer tu contraseña.\n\nTu código de recuperación es:\n\n%s\n\nCaduca en 15 minutos. Si no la solicitaste, ignora este mensaje.",
			u.FirstName, code,
		)
		if err := smtpCfg.Send(u.Email, "Código de recuperación de contraseña - Jeux Airel", body); err != nil {
			slog.Error("error enviando email de recuperación", "user_id", u.ID, "error", err)
			return "", ErrEmailNotSent
		}
		return "", nil
	}

	slog.Info("password reset generado (SMTP no configurado)", "user_id", u.ID, "code", code)
	return code, nil
}

// ResetPassword valida el código de reseteo y cambia la contraseña.
func (s *Service) ResetPassword(req ResetPasswordRequest) error {
	if !validation.ValidatePassword(req.NewPassword) {
		return ErrWeakPassword
	}

	hash := hashToken(req.Code)

	newHash, err := bcrypt.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Transacción: validar, actualizar contraseña, marcar el reset como usado
	// y revocar sesiones vigentes de forma atómica.
	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		// Rechequeo dentro de la transacción para evitar condiciones de carrera
		// entre handler y confirmación del código.
		var xt *session.PasswordResetRequest
		if err := tx.Where("token_hash = ?", hash).First(&xt).Error; err != nil {
			return ErrInvalidToken
		}
		if xt.IsUsed {
			return ErrInvalidToken
		}
		if xt.ExpiresAt.Before(time.Now()) {
			return ErrTokenExpired
		}

		sessRepo := session.NewRepository(tx)
		userRepo := user.NewRepository(tx)

		if err := sessRepo.RevokeAllByUser(xt.UserID); err != nil {
			return err
		}
		if err := userRepo.UpdatePassword(xt.UserID, newHash); err != nil {
			return err
		}
		return sessRepo.MarkPasswordResetUsed(xt.ID)
	})
}

// issueTokens genera los tokens de acceso y refresco y persiste el refresh.
func (s *Service) issueTokens(u *user.User, ipAddress string) (*TokenResponse, error) {
	accessToken, err := jwt.GenerateAccessToken(
		u.ID, u.Email, u.ActorType,
		s.cfg.JWTSecret, s.cfg.JWTAccessTTL,
	)
	if err != nil {
		return nil, err
	}

	rawRefresh, err := randomToken()
	if err != nil {
		return nil, err
	}

	refresh := &session.RefreshToken{
		UserID:    u.ID,
		TokenHash: hashToken(rawRefresh),
		IPAddress: &ipAddress,
		ExpiresAt: time.Now().Add(s.cfg.JWTRefreshTTL),
	}
	if err := s.repo.Sessions.CreateRefreshToken(refresh); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}

// recordAudit registra un evento de inicio de sesión en la auditoría.
func (s *Service) recordAudit(userID *string, email, outcome, ipAddress, userAgent string) {
	entry := &audit.AuditLogin{
		UserID:         userID,
		EmailAttempted: email,
		Outcome:        outcome,
		IPAddress:      &ipAddress,
		UserAgent:      &userAgent,
	}
	if err := s.repo.Audit.Create(entry); err != nil {
		slog.Warn("no se pudo registrar auditoría de login", "error", err)
	}
}

// hashToken calcula el SHA-256 de un token en formato hexadecimal.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// randomToken genera un token aleatorio seguro de 32 bytes en hex.
func randomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// generateResetCode genera un código numérico de 6 dígitos.
func generateResetCode() (string, error) {
	var b [3]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	return fmt.Sprintf("%06d", n%1000000), nil
}
