package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/shared/jwt"
	"github.com/jesus-ariel/api-jeussairel/shared/response"
)

// Context constants inyectadas en el gin.Context.
const (
	ContextUserID    = "user_id"
	ContextEmail     = "email"
	ContextActorType = "actor_type"
)

// RequireAuth valida el token JWT del header Authorization y deja
// los claims en el contexto para los handlers.
func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, http.StatusUnauthorized, "MISSING_TOKEN", "token de autorización requerido")
			c.Abort()
			return
		}

		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			response.Error(c, http.StatusUnauthorized, "INVALID_TOKEN", "formato de token inválido")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(token, secret)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				response.Error(c, http.StatusUnauthorized, "TOKEN_EXPIRED", "token expirado")
			} else {
				response.Error(c, http.StatusUnauthorized, "INVALID_TOKEN", "token inválido")
			}
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextActorType, claims.ActorType)
		c.Next()
	}
}