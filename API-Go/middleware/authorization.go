package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/shared/response"
	"gorm.io/gorm"
)

// Jerarquía de niveles de acción. Un permiso con nivel superior cubre
// los inferiores (READ < WRITE < DELETE < PUBLISH < APPROVE).
var actionRank = map[string]int{
	"READ":    1,
	"WRITE":   2,
	"DELETE":  3,
	"PUBLISH": 4,
	"APPROVE": 5,
}

// RequirePermission es el middleware de autorización RBAC. Debe ejecutarse
// después de RequireAuth: comprueba que el usuario autenticado tiene el
// feature requerido con un nivel de acción suficiente, ya sea por un
// user_scope_override vigente o por sus roles activos.
func RequirePermission(db *gorm.DB, featureCode, action string) gin.HandlerFunc {
	required := actionRank[action]
	if required <= 0 {
		panic("middleware: acción de permiso inválida: " + action)
	}

	return func(c *gin.Context) {
		userID, _ := c.Get(ContextUserID)

		// 1. Resolver el feature requerido y su nivel base.
		var featureID string
		if err := db.Raw(
			`SELECT id FROM rbac_catalog.feature
			 WHERE code = ? AND is_active = true LIMIT 1`,
			featureCode,
		).Row().Scan(&featureID); err != nil {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "permiso requerido")
			c.Abort()
			return
		}

		// 2. Overrides de scope vigentes para este feature.
		//    Un override implícito ('denegar') bloquea; un override de
		//    concesión vigente lo permite.
		var suOutcome string
		if err := db.Raw(
			`SELECT CASE
				WHEN COUNT(*) FILTER (WHERE is_allowed = false) > 0 THEN 'denied'
				WHEN COUNT(*) FILTER (WHERE is_allowed = true) > 0 THEN 'granted'
				ELSE 'none'
			 END
			 FROM rbac.user_scope_override
			 WHERE user_id = $1 AND feature_id = $2
			   AND (expires_at IS NULL OR expires_at > now())`,
			userID, featureID,
		).Row().Scan(&suOutcome); err == nil {
			switch suOutcome {
			case "denied":
				response.Error(c, http.StatusForbidden, "FORBIDDEN", "permiso denegado")
				c.Abort()
				return
			case "granted":
				c.Next()
				return
			}
		}

		// 3. Buscar un rol activo del usuario que incluya el feature con
		//    el nivel de acción requerido.
		var featureAction string
		if err := db.Raw(
			`SELECT f.action_level
			 FROM rbac.user_role AS ur
			 JOIN rbac.role_feature AS rf ON rf.role_id = ur.role_id
			 JOIN rbac_catalog.feature AS f ON f.id = rf.feature_id
			 WHERE ur.user_id = $1 AND f.id = $2
			   AND (ur.expires_at IS NULL OR ur.expires_at > now())
			 ORDER BY (CASE f.action_level
				WHEN 'APPROVE' THEN 5 WHEN 'PUBLISH' THEN 4
				WHEN 'DELETE' THEN 3 WHEN 'WRITE' THEN 2 ELSE 1 END) DESC
			 LIMIT 1`,
			userID, featureID,
		).Row().Scan(&featureAction); err != nil {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "permiso requerido")
			c.Abort()
			return
		}

		if actionRank[featureAction] >= required {
			c.Next()
			return
		}

		response.Error(c, http.StatusForbidden, "FORBIDDEN", "permiso requerido")
		c.Abort()
	}
}