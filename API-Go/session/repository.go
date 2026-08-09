package session

import (
	"time"

	"gorm.io/gorm"
)

// Repository encapsula el acceso a datos de sesiones (solo GORM).
type Repository struct {
	db *gorm.DB
}

// NewRepository crea el repositorio de sesiones.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ============================== Refresh Tokens ==============================

// CreateRefreshToken guarda un refresh token.
func (r *Repository) CreateRefreshToken(t *RefreshToken) error {
	return r.db.Create(t).Error
}

// FindRefreshTokenByHash obtiene un refresh token por su hash.
func (r *Repository) FindRefreshTokenByHash(tokenHash string) (*RefreshToken, error) {
	var t RefreshToken
	if err := r.db.Where("token_hash = ?", tokenHash).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// ListActiveByUser devuelve los refresh tokens activos de un usuario.
func (r *Repository) ListActiveByUser(userID string) ([]RefreshToken, error) {
	var tokens []RefreshToken
	if err := r.db.
		Where("user_id = ? AND is_revoked = ? AND expires_at > ?", userID, false, time.Now()).
		Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

// RevokeRefreshToken marca un token como revocado.
// Devuelve gorm.ErrRecordNotFound si el token no existe.
func (r *Repository) RevokeRefreshToken(id string) error {
	res := r.db.Model(&RefreshToken{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// RevokeAllByUser revoca todos los tokens activos de un usuario.
func (r *Repository) RevokeAllByUser(userID string) error {
	return r.db.Model(&RefreshToken{}).Where("user_id = ? AND is_revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": time.Now(),
		}).Error
}

// DeleteExpiredRefreshTokens elimina tokens expirados o revocados.
func (r *Repository) DeleteExpiredRefreshTokens() error {
	return r.db.Delete(&RefreshToken{}, "is_revoked = ? OR expires_at <= ?", true, time.Now()).Error
}

// ============================== Password Reset ==============================

// CreatePasswordReset guarda una solicitud de reseteo de contraseña.
func (r *Repository) CreatePasswordReset(p *PasswordResetRequest) error {
	return r.db.Create(p).Error
}

// FindPasswordResetByHash obtiene una solicitud de reseteo por su hash.
func (r *Repository) FindPasswordResetByHash(tokenHash string) (*PasswordResetRequest, error) {
	var p PasswordResetRequest
	if err := r.db.Where("token_hash = ?", tokenHash).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// MarkPasswordResetUsed marca una solicitud como usada.
func (r *Repository) MarkPasswordResetUsed(id string) error {
	return r.db.Model(&PasswordResetRequest{}).Where("id = ?", id).
		Update("is_used", true).Error
}

// InvalidatePasswordResetsByUser invalida todas las solicitudes pendientes de un usuario.
func (r *Repository) InvalidatePasswordResetsByUser(userID string) error {
	return r.db.Model(&PasswordResetRequest{}).
		Where("user_id = ? AND is_used = ?", userID, false).
		Update("is_used", true).Error
}

// ============================== Email Verification ==============================

// CreateEmailVerification guarda un código de verificación de email.
func (r *Repository) CreateEmailVerification(v *EmailVerificationCode) error {
	return r.db.Create(v).Error
}

// InvalidateEmailVerificationsByUser marca como usadas las verificaciones
// pendientes de un usuario (se le manda un código nuevo).
func (r *Repository) InvalidateEmailVerificationsByUser(userID string) error {
	return r.db.Model(&EmailVerificationCode{}).
		Where("user_id = ? AND is_used = ?", userID, false).
		Update("is_used", true).Error
}

// FindEmailVerificationByHash obtiene una verificación por su hash.
func (r *Repository) FindEmailVerificationByHash(codeHash string) (*EmailVerificationCode, error) {
	var v EmailVerificationCode
	if err := r.db.Where("code_hash = ?", codeHash).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// MarkEmailVerificationUsed marca una verificación como usada.
func (r *Repository) MarkEmailVerificationUsed(id string) error {
	return r.db.Model(&EmailVerificationCode{}).Where("id = ?", id).
		Update("is_used", true).Error
}