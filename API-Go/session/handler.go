package session

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/shared/response"
)

// Handler expone los endpoints de sesiones a través de HTTP.
type Handler struct {
	service *Service
}

// NewHandler crea el handler de sesiones.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListActiveByUser devuelve las sesiones activas del usuario autenticado.
// @Summary Listar sesiones activas
// @Description Devuelve las sesiones (refresh tokens) activas del usuario autenticado.
// @Tags Sessions
// @Security BearerAuth
// @Success 200 {object} response.Body "Sesiones activas"
// @Failure 401 {object} response.Body "Token inválido o expirado"
// @Router /sessions/active [get]
func (h *Handler) ListActiveByUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := userID.(string)

	items, err := h.service.ListActiveByUser(id)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "sesiones activas", gin.H{"items": items})
}

// RevokeById revoca una sesión por su ID.
// @Summary Revocar sesión
// @Description Revoca una sesión del usuario autenticado por su ID.
// @Tags Sessions
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID de la sesión"
// @Success 200 {object} response.Body "Sesión revocada"
// @Failure 401 {object} response.Body "Token inválido o expirado"
// @Failure 404 {object} response.Body "Sesión no encontrada"
// @Router /sessions/{id} [delete]
func (h *Handler) RevokeById(c *gin.Context) {
	if err := h.service.RevokeById(c.Param("id")); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusOK, "sesión revocada")
}