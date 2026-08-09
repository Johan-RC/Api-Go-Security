package audit

import (
	"time"

	"github.com/jesus-ariel/api-jeussairel/shared/uuid"
	"github.com/jesus-ariel/api-jeussairel/user"
	"gorm.io/gorm"
)

// AuditLogin representa la tabla identity_audit.audit_login.
type AuditLogin struct {
	ID             string     `gorm:"column:id;type:uuid;primaryKey"`
	UserID         *string    `gorm:"column:user_id;type:uuid;index:ix_audit_login_user_id"`
	EmailAttempted string     `gorm:"column:email_attempted;type:varchar(255);not null;index:ix_audit_login_email_attempted_attempted_at,priority:1"`
	Outcome        string     `gorm:"column:outcome;type:varchar(20);not null"`
	IPAddress      *string    `gorm:"column:ip_address;type:varchar(45)"`
	UserAgent      *string    `gorm:"column:user_agent;type:text"`
	AttemptedAt    time.Time  `gorm:"column:attempted_at;type:timestamptz;not null;default:now();index:ix_audit_login_email_attempted_attempted_at,priority:2"`

	User *user.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:SET NULL"`
}

// BeforeCreate asigna un UUID antes de insertar.
func (a *AuditLogin) BeforeCreate(_ *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New()
	}
	return nil
}

// TableName apunta al esquema y tabla exactos de PostgreSQL.
func (AuditLogin) TableName() string {
	return "identity_audit.audit_login"
}

// Restricciones adicionales de identity_audit.audit_login:
//   - CHECK outcome IN ('SUCCESS','INVALID_PASSWORD','USER_NOT_FOUND',
//     'ACCOUNT_LOCKED','TOKEN_EXPIRED')
