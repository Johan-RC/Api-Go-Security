package role

import (
	"time"

	"github.com/jesus-ariel/api-jeussairel/feature"
	"github.com/jesus-ariel/api-jeussairel/shared/uuid"
	"github.com/jesus-ariel/api-jeussairel/user"
	"gorm.io/gorm"
)

// Role representa la tabla rbac.role.
type Role struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey"`
	Name         string    `gorm:"column:name;type:varchar(50);not null;uniqueIndex:uq_role_name"`
	DisplayName  string    `gorm:"column:display_name;type:varchar(100);not null"`
	Description  *string   `gorm:"column:description;type:text"`
	IsSystemRole bool      `gorm:"column:is_system_role;type:boolean;not null;default:false"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()"`

	RoleFeatures []RoleFeature `gorm:"foreignKey:RoleID"`
	UserRoles    []UserRole    `gorm:"foreignKey:RoleID"`
}

// BeforeCreate asigna un UUID antes de insertar.
func (r *Role) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New()
	}
	return nil
}

// TableName apunta al esquema y tabla exactos de PostgreSQL.
func (Role) TableName() string {
	return "rbac.role"
}

// RoleFeature representa la tabla rbac.role_feature.
type RoleFeature struct {
	ID        string `gorm:"column:id;type:uuid;primaryKey"`
	RoleID    string `gorm:"column:role_id;type:uuid;not null;uniqueIndex:uq_role_feature_role_id_feature_id,priority:1"`
	FeatureID string `gorm:"column:feature_id;type:uuid;not null;uniqueIndex:uq_role_feature_role_id_feature_id,priority:2;index:ix_role_feature_feature_id"`
	ScopeType string `gorm:"column:scope_type;type:varchar(30);not null"`

	Role    *Role            `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	Feature *feature.Feature `gorm:"foreignKey:FeatureID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

// BeforeCreate asigna un UUID antes de insertar.
func (r *RoleFeature) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New()
	}
	return nil
}

// TableName apunta al esquema y tabla exactos de PostgreSQL.
func (RoleFeature) TableName() string {
	return "rbac.role_feature"
}

// Restricciones adicionales de rbac.role_feature:
//   - CHECK scope_type IN ('GLOBAL','TRAINING_CENTER','AREA','OWN_FICHAS',
//     'OWN_SCHEDULE','OWN_PROFILE','OWN_FICHA_AS_LEARNER')

// UserRole representa la tabla rbac.user_role.
type UserRole struct {
	ID              string     `gorm:"column:id;type:uuid;primaryKey"`
	UserID          string     `gorm:"column:user_id;type:uuid;not null;uniqueIndex:uq_user_role_user_id_role_id_training_center_id,priority:1;index:ix_user_role_user_id"`
	RoleID          string     `gorm:"column:role_id;type:uuid;not null;uniqueIndex:uq_user_role_user_id_role_id_training_center_id,priority:2;index:ix_user_role_role_id"`
	TrainingCenterID *string   `gorm:"column:training_center_id;type:uuid;uniqueIndex:uq_user_role_user_id_role_id_training_center_id,priority:3"`
	AssignedBy      string     `gorm:"column:assigned_by;type:uuid;not null;index:ix_user_role_assigned_by"`
	AssignedAt      time.Time  `gorm:"column:assigned_at;type:timestamptz;not null;default:now()"`
	ExpiresAt       *time.Time `gorm:"column:expires_at;type:timestamptz"`

	User           *user.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	Role           *Role      `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	AssignedByUser *user.User `gorm:"foreignKey:AssignedBy;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

// BeforeCreate asigna un UUID antes de insertar.
func (u *UserRole) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New()
	}
	return nil
}

// TableName apunta al esquema y tabla exactos de PostgreSQL.
func (UserRole) TableName() string {
	return "rbac.user_role"
}

// UserScopeOverride representa la tabla rbac.user_scope_override.
type UserScopeOverride struct {
	ID        string     `gorm:"column:id;type:uuid;primaryKey"`
	UserID    string     `gorm:"column:user_id;type:uuid;not null;index:ix_user_scope_override_user_id_feature_id,priority:1"`
	FeatureID string     `gorm:"column:feature_id;type:uuid;not null;index:ix_user_scope_override_user_id_feature_id,priority:2;index:ix_user_scope_override_feature_id"`
	ScopeType string     `gorm:"column:scope_type;type:varchar(30);not null"`
	IsAllowed bool       `gorm:"column:is_allowed;type:boolean;not null;default:true"`
	Reason    string     `gorm:"column:reason;type:text;not null"`
	GrantedBy string     `gorm:"column:granted_by;type:uuid;not null;index:ix_user_scope_override_granted_by"`
	ExpiresAt *time.Time `gorm:"column:expires_at;type:timestamptz"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`

	User           *user.User       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	Feature        *feature.Feature `gorm:"foreignKey:FeatureID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	GrantedByUser  *user.User       `gorm:"foreignKey:GrantedBy;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

// BeforeCreate asigna un UUID antes de insertar.
func (u *UserScopeOverride) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New()
	}
	return nil
}

// TableName apunta al esquema y tabla exactos de PostgreSQL.
func (UserScopeOverride) TableName() string {
	return "rbac.user_scope_override"
}

// Restricciones adicionales de rbac.user_scope_override:
//   - CHECK scope_type IN ('GLOBAL','TRAINING_CENTER','AREA','OWN_FICHAS',
//     'OWN_SCHEDULE','OWN_PROFILE','OWN_FICHA_AS_LEARNER')
