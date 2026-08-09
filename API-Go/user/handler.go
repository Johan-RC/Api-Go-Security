package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/shared/response"
)

// Handler expone los endpoints de usuarios a través de HTTP.
type Handler struct {
	service *Service
}

// NewHandler crea el handler de usuarios.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create registra un nuevo usuario.
// @Summary Crear usuario
// @Description Registra un nuevo usuario en el sistema.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateUserRequest true "Datos del usuario"
// @Success 201 {object} response.Body{data=user.UserResponse} "Usuario creado"
// @Failure 400 {object} response.Body "Datos inválidos"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_USER_MANAGE:write)"
// @Router /users [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	u, err := h.service.Create(req)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "usuario creado", u)
}

// GetById devuelve un usuario por su ID.
// @Summary Obtener usuario por ID
// @Description Devuelve un usuario por su UUID.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID del usuario"
// @Success 200 {object} response.Body{data=user.UserResponse} "Usuario encontrado"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_USER_VIEW:read)"
// @Failure 404 {object} response.Body "Usuario no encontrado"
// @Router /users/{id} [get]
func (h *Handler) GetById(c *gin.Context) {
	u, err := h.service.GetById(c.Param("id"))
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "usuario encontrado", u)
}

// List devuelve usuarios paginados.
// @Summary Listar usuarios
// @Description Devuelve una lista paginada de usuarios.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página (default 1)"
// @Param page_size query int false "Tamaño de página 1-100 (default 20)"
// @Success 200 {object} response.Body{data=user.ListResponse} "Usuarios obtenidos"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_USER_VIEW:read)"
// @Router /users [get]
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := h.service.List(page, pageSize)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "usuarios obtenidos", gin.H{
		"items": users,
		"total": total,
		"page":  page,
	})
}

// Update modifica un usuario.
// @Summary Actualizar usuario
// @Description Actualiza los campos editables de un usuario.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID del usuario"
// @Param body body UpdateUserRequest true "Campos a actualizar"
// @Success 200 {object} response.Body{data=user.UserResponse} "Usuario actualizado"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_USER_MANAGE:write)"
// @Failure 404 {object} response.Body "Usuario no encontrado"
// @Router /users/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	u, err := h.service.Update(c.Param("id"), req)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "usuario actualizado", u)
}

// Delete elimina un usuario.
// @Summary Eliminar usuario
// @Description Elimina (desactiva) un usuario por su UUID.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID del usuario"
// @Success 200 {object} response.Body "Usuario eliminado"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_USER_MANAGE:write)"
// @Failure 404 {object} response.Body "Usuario no encontrado"
// @Router /users/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusOK, "usuario eliminado")
}