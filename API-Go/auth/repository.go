package auth

import (
	"github.com/jesus-ariel/api-jeussairel/audit"
	"github.com/jesus-ariel/api-jeussairel/role"
	"github.com/jesus-ariel/api-jeussairel/session"
	"github.com/jesus-ariel/api-jeussairel/user"
	"gorm.io/gorm"
)

// Repository encapsula el acceso a datos del módulo auth.
// Delega en los repositorios de los módulos dueños de cada tabla (user,
// role, session, audit) y ejecuta SQL directo para consultas de permisos
// efectivos del usuario autenticado.
type Repository struct {
	db      *gorm.DB
	Users   *user.Repository
	Roles   *role.Repository
	Sessions *session.Repository
	Audit   *audit.Repository
}

// NewRepository crea el repositorio de autenticación inyectando los repos
// de los módulos que componen el flujo de login.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db:       db,
		Users:    user.NewRepository(db),
		Roles:    role.NewRepository(db),
		Sessions: session.NewRepository(db),
		Audit:    audit.NewRepository(db),
	}
}

// ListRolesByUser devuelve los roles activos del usuario autenticado.
func (r *Repository) ListRolesByUser(userID string) ([]role.Role, error) {
	var roles []role.Role
	if err := r.db.Raw(
		`SELECT rl.*
		 FROM rbac.user_role AS ur
		 JOIN rbac.role AS rl ON rl.id = ur.role_id
		 WHERE ur.user_id = ? AND (ur.expires_at IS NULL OR ur.expires_at > now())`,
		userID,
	).Scan(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// ListPermissionsByUser devuelve los features con su nivel de acción y el
// scope al que tiene acceso el usuario (vía roles activos o overrides).
// Replica la lógica del middleware de autorización para alimentar la UI.
func (r *Repository) ListPermissionsByUser(userID string) ([]Permission, error) {
	rows, err := r.db.Raw(
		`WITH resolved AS (
			SELECT DISTINCT f.code, f.action_level, rf.scope_type
			FROM rbac.user_role AS ur
			JOIN rbac.role_feature AS rf ON rf.role_id = ur.role_id
			JOIN rbac_catalog.feature AS f ON f.id = rf.feature_id
			WHERE ur.user_id = ? AND (ur.expires_at IS NULL OR ur.expires_at > now())
			UNION
			SELECT DISTINCT f.code, f.action_level, so.scope_type
			FROM rbac.user_scope_override AS so
			JOIN rbac_catalog.feature AS f ON f.id = so.feature_id
			WHERE so.user_id = ? AND so.is_allowed = true
			  AND (so.expires_at IS NULL OR so.expires_at > now())
		)
		SELECT code, MAX(action_level) AS action_level, MAX(scope_type) AS scope_type
		FROM resolved
		GROUP BY code`,
		userID, userID,
	).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.Code, &p.Action, &p.Scope); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, rows.Err()
}

// DB devuelve la conexión subyacente para poder ejecutar transacciones
// que abarquen varios repositorios.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
