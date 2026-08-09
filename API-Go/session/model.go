package session

import (
	"time"

	"github.com/jesus-ariel/api-jeussairel/shared/uuid"
	"github.com/jesus-ariel/api-jeussairel/user"
	"gorm.io/gorm"
)

// RefreshToken representa la tabla session.refresh_token.
type RefreshToken struct {
	ID         string     `gorm:"column:id;type:uuid;primaryKey"`
	UserID     string     `gorm:"column:user_id;type:uuid;not null;index:ix_refresh_token_user_id_is_revoked,priority:1"`
	TokenHash  string     `gorm:"column:token_hash;type:text;not null;uniqueIndex:uq_refresh_token_token_hash"`
	DeviceHint *string    `gorm:"column:device_hint;type:varchar(200)"`
	IPAddress  *string    `gorm:"column:ip_address;type:varchar(45)"`
	ExpiresAt  time.Time  `gorm:"column:expires_at;type:timestamptz;not null;default:(now() + interval '7 days')"`
	IsRevoked  bool       `gorm:"column:is_revoked;type:boolean;not null;default:false;index:ix_refresh_token_user_id_is_revoked,priority:2"`
	RevokedAt  *time.Time `gorm:"column:revoked_at;type:timestamptz"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`

	User *user.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
}

// BeforeCreate asigna un UUID antes de insertar.
func (t *RefreshToken) BeforeCreate(_ *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New()
	}
	return nil
}

// TableName apunta al esquema y tabla exactos de PostgreSQL.
func (RefreshToken) TableName() string {
	return "session.refresh_token"
}

// PasswordResetRequest representa la tabla session.password_reset_request.
type PasswordResetRequest struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	UserID      string    `gorm:"column:user_id;type:uuid;not null;index:ix_password_reset_request_user_id"`
	TokenHash   string    `gorm:"column:token_hash;type:text;not null"`
	ExpiresAt   time.Time `gorm:"column:expires_at;type:timestamptz;not null;default:(now() + interval '1 hour')"`
	IsUsed      bool      `gorm:"column:is_used;type:boolean;not null;default:false"`
	RequestedAt time.Time `gorm:"column:requested_at;type:timestamptz;not null;default:now()"`
	IPAddress   *string   `gorm:"column:ip_address;type:varchar(45)"`

	User *user.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
}

// BeforeCreate asigna un UUID antes de insertar.
func (u *PasswordResetRequest) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New()
	}
	return nil
}

// TableName apunta al esquema y tabla exactos de PostgreSQL.
func (PasswordResetRequest) TableName() string {
	return "session.password_reset_request"
}

// EmailVerificationCode representa la tabla session.email_verification_code.
// Almacena el hash del código de 6 dígitos usado para activar una cuenta.
type EmailVerificationCode struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	UserID      string    `gorm:"column:user_id;type:uuid;not null;index:ix_email_verification_user_id"`
	CodeHash    string    `gorm:"column:code_hash;type:text;not null"`
	ExpiresAt   time.Time `gorm:"column:expires_at;type:timestamptz;not null;default:(now() + interval '15 minutes')"`
	IsUsed      bool      `gorm:"column:is_used;type:boolean;not null;default:false"`
	RequestedAt time.Time `gorm:"column:requested_at;type:timestamptz;not null;default:now()"`
	IPAddress   *string   `gorm:"column:ip_address;type:varchar(45)"`

	User *user.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
}

// BeforeCreate asigna un UUID antes de insertar.
func (v *EmailVerificationCode) BeforeCreate(_ *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.New()
	}
	return nil
}

// TableName apunta al esquema y tabla exactos de PostgreSQL.
func (EmailVerificationCode) TableName() string {
	return "session.email_verification_code"
}
