package feature

import (
	"time"

	"github.com/jesus-ariel/api-jeussairel/shared/uuid"
	"gorm.io/gorm"
)

// Module representa la tabla rbac_catalog.module.
type Module struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey"`
	Code         string    `gorm:"column:code;type:varchar(30);not null;uniqueIndex:uq_module_code"`
	Name         string    `gorm:"column:name;type:varchar(100);not null"`
	Description  *string   `gorm:"column:description;type:text"`
	DisplayOrder int16     `gorm:"column:display_order;type:smallint;not null"`
	IconKey      *string   `gorm:"column:icon_key;type:varchar(50)"`
	IsActive     bool      `gorm:"column:is_active;type:boolean;not null;default:true"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`

	Features []Feature `gorm:"foreignKey:ModuleID"`
}

// BeforeCreate asigna un UUID antes de insertar.
func (m *Module) BeforeCreate(_ *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New()
	}
	return nil
}

// TableName apunta al esquema y tabla exactos de PostgreSQL.
func (Module) TableName() string {
	return "rbac_catalog.module"
}

// Feature representa la tabla rbac_catalog.feature.
type Feature struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	ModuleID    string    `gorm:"column:module_id;type:uuid;not null;index:ix_feature_module_id"`
	Code        string    `gorm:"column:code;type:varchar(60);not null;uniqueIndex:uq_feature_code"`
	Name        string    `gorm:"column:name;type:varchar(120);not null"`
	Description *string   `gorm:"column:description;type:text"`
	ActionLevel string    `gorm:"column:action_level;type:varchar(20);not null"`
	IsActive    bool      `gorm:"column:is_active;type:boolean;not null;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`

	Module *Module `gorm:"foreignKey:ModuleID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

// BeforeCreate asigna un UUID antes de insertar.
func (f *Feature) BeforeCreate(_ *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New()
	}
	return nil
}

// TableName apunta al esquema y tabla exactos de PostgreSQL.
func (Feature) TableName() string {
	return "rbac_catalog.feature"
}

// Restricciones adicionales de rbac_catalog.feature:
//   - CHECK action_level IN ('READ', 'WRITE', 'DELETE', 'PUBLISH', 'APPROVE')