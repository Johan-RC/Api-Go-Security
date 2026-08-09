package user

import (
	"time"

	"github.com/jesus-ariel/api-jeussairel/shared/uuid"
	"gorm.io/gorm"
)

// User representa la tabla identity.user.
type User struct {
	ID             string     `gorm:"column:id;type:uuid;primaryKey"`
	Email          string     `gorm:"column:email;type:varchar(255);not null;uniqueIndex:uq_user_email"`
	PasswordHash   string     `gorm:"column:password_hash;type:text;not null"`
	FirstName      string     `gorm:"column:first_name;type:varchar(100);not null"`
	LastName       string     `gorm:"column:last_name;type:varchar(100);not null"`
	ActorType      string     `gorm:"column:actor_type;type:varchar(20);not null"`
	ActorID        *string    `gorm:"column:actor_id;type:uuid;index:ix_user_actor_id"`
	IsActive       bool       `gorm:"column:is_active;type:boolean;not null"`
	LastLoginAt    *time.Time `gorm:"column:last_login_at;type:timestamptz"`
	FailedAttempts int16      `gorm:"column:failed_attempts;type:smallint;not null;default:0"`
	LockedUntil    *time.Time `gorm:"column:locked_until;type:timestamptz"`
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

// BeforeCreate asigna un UUID antes de insertar.
func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New()
	}
	return nil
}

// TableName apunta al esquema y tabla exactos de PostgreSQL.
func (User) TableName() string {
	return "identity.user"
}

// Restricciones e índices adicionales de la tabla identity.user:
//   - CHECK actor_type IN ('USER', 'INSTRUCTOR', 'LEARNER')
//   - UNIQUE (email) -> uq_user_email
//   - UNIQUE INDEX lower(email) -> uq_user_email_lower (no expresable como tag GORM)
//   - INDEX (last_name, first_name) -> ix_user_last_name_first_name
//   - PARTIAL INDEX (is_active) WHERE is_active = TRUE -> ix_user_is_active
