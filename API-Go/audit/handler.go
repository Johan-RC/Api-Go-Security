package audit

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/shared/response"
)

// Handler expone los endpoints de auditoría a través de HTTP.
type Handler struct {
	service *Service
}

// NewHandler crea el handler de auditoría.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListByUser devuelve la auditoría de login de un usuario.
// @Summary Auditoría de login de un usuario
// @Description Devuelve el historial de inicios de sesión de un usuario.
// @Tags Audit
// @Security BearerAuth
// @Param userId path string true "UUID del usuario"
// @Success 200 {object} response.Body "Auditoría del usuario"
// @Failure 403 {object} response.Body "Acceso denegado (AUDIT_LOG_VIEW:read)"
// @Failure 404 {object} response.Body "Usuario no encontrado"
// @Router /audit/logins/users/{userId} [get]
func (h *Handler) ListByUser(c *gin.Context) {
	entries, err := h.service.ListByUser(c.Param("userId"))
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "auditoría del usuario", gin.H{"items": entries})
}

// List devuelve la auditoría paginada.
// @Summary Listar auditoría de logins
// @Description Devuelve la auditoría de inicios de sesión paginada.
// @Tags Audit
// @Security BearerAuth
// @Param page query int false "Número de página (default 1)"
// @Param page_size query int false "Tamaño de página (1-100, default 20)"
// @Success 200 {object} response.Body "Auditoría obtenida"
// @Failure 403 {object} response.Body "Acceso denegado (AUDIT_LOG_VIEW:read)"
// @Router /audit/logins [get]
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	entries, total, err := h.service.List(page, pageSize)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "auditoría obtenida", gin.H{
		"items": entries,
		"total": total,
		"page":  page,
	})
}