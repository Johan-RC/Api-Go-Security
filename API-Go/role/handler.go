package role

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/shared/response"
)

// Handler expone los endpoints de roles y RBAC a través de HTTP.
type Handler struct {
	service *Service
}

// NewHandler crea el handler de roles.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// assigneeID recupera el user_id del contexto (inyectado por RequireAuth)
// sin panics si la clave o el tipo no existen.
func assigneeID(c *gin.Context) string {
	v, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	id, ok := v.(string)
	if !ok {
		return ""
	}
	return id
}

// Create registra un nuevo rol.
// @Summary Crear rol
// @Description Registra un nuevo rol en el sistema.
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateRoleRequest true "Datos del rol"
// @Success 201 {object} response.Body{data=role.RoleResponse} "Rol creado"
// @Failure 400 {object} response.Body "Datos inválidos"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_ROLE_MANAGE:write)"
// @Router /roles [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	rl, err := h.service.Create(req)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "rol creado", rl)
}

// GetById devuelve un rol por su ID.
// @Summary Obtener rol por ID
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "UUID del rol"
// @Accept json
// @Produce json
// @Success 200 {object} response.Body{data=role.RoleResponse} "Rol encontrado"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_ROLE_VIEW:read)"
// @Failure 404 {object} response.Body "Rol no encontrado"
// @Router /roles/{id} [get]
func (h *Handler) GetById(c *gin.Context) {
	rl, err := h.service.GetById(c.Param("id"))
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "rol encontrado", rl)
}

// List devuelve roles paginados.
// @Summary Listar roles
// @Description Devuelve una lista paginada de roles.
// @Tags Roles
// @Security BearerAuth
// @Param page query int false "Número de página (default 1)"
// @Param page_size query int false "Tamaño de página 1-100 (default 20)"
// @Success 200 {object} response.Body{data=role.ListResponse} "Roles obtenidos"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_ROLE_VIEW:read)"
// @Router /roles [get]
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	roles, total, err := h.service.List(page, pageSize)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "roles obtenidos", gin.H{
		"items": roles,
		"total": total,
		"page":  page,
	})
}

// Update modifica un rol.
// @Summary Actualizar rol
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "UUID del rol"
// @Param body body UpdateRoleRequest true "Campos a actualizar"
// @Success 200 {object} response.Body{data=role.RoleResponse} "Rol actualizado"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_ROLE_MANAGE:write)"
// @Failure 404 {object} response.Body "Rol no encontrado"
// @Router /roles/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	rl, err := h.service.Update(c.Param("id"), req)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "rol actualizado", rl)
}

// Delete elimina un rol.
// @Summary Eliminar rol
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "UUID del rol"
// @Success 200 {object} response.Body "Rol eliminado"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_ROLE_MANAGE:write)"
// @Failure 404 {object} response.Body "Rol no encontrado"
// @Router /roles/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusOK, "rol eliminado")
}

// AssignFeatures reasigna las funcionalidades de un rol.
// @Summary Asignar funcionalidades a un rol
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "UUID del rol"
// @Param body body AssignFeaturesRequest true "Lista de funcionalidades con ámbito"
// @Success 200 {object} response.Body "Funcionalidades asignadas"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_ROLE_MANAGE:write)"
// @Failure 404 {object} response.Body "Rol no encontrado"
// @Router /roles/{id}/features [put]
func (h *Handler) AssignFeatures(c *gin.Context) {
	var req AssignFeaturesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	if err := h.service.AssignFeatures(c.Param("id"), req); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusOK, "funcionalidades asignadas")
}

// ListFeatures devuelve las funcionalidades de un rol.
// @Summary Listar funcionalidades de un rol
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "UUID del rol"
// @Success 200 {object} response.Body{data=role.ListFeaturesResponse} "Funcionalidades del rol"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_ROLE_VIEW:read)"
// @Failure 404 {object} response.Body "Rol no encontrado"
// @Router /roles/{id}/features [get]
func (h *Handler) ListFeatures(c *gin.Context) {
	items, err := h.service.ListFeatures(c.Param("id"))
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "funcionalidades del rol", gin.H{"items": items})
}

// AssignRoleToUser asigna un rol a un usuario.
// @Summary Asignar rol a un usuario
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "UUID del usuario"
// @Param body body AssignRoleRequest true "ID del rol"
// @Success 201 {object} response.Body "Rol asignado al usuario"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_ROLE_ASSIGN:write)"
// @Failure 404 {object} response.Body "Usuario o rol no encontrado"
// @Router /users/{id}/roles [post]
func (h *Handler) AssignRoleToUser(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	if err := h.service.AssignRoleToUser(c.Param("id"), req.RoleID, assigneeID(c)); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusCreated, "rol asignado al usuario")
}

// ListRolesByUser devuelve los roles de un usuario.
// @Summary Listar roles de un usuario
// @Security BearerAuth
// @Param id path string true "UUID del usuario"
// @Success 200 {object} response.Body{data=role.UserRoleResponse} "Roles del usuario"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_ROLE_VIEW:read)"
// @Failure 404 {object} response.Body "Usuario no encontrado"
// @Router /users/{id}/roles [get]
func (h *Handler) ListRolesByUser(c *gin.Context) {
	items, err := h.service.ListRolesByUser(c.Param("id"))
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "roles del usuario", gin.H{"items": items})
}

// RemoveRoleFromUser quita un rol de un usuario.
// @Summary Quitar rol de un usuario
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "UUID del usuario"
// @Param userRoleId path string true "UUID de la asignación usuario-rol"
// @Success 200 {object} response.Body "Rol removido del usuario"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_ROLE_ASSIGN:write)"
// @Failure 404 {object} response.Body "Asignación no encontrada"
// @Router /users/{id}/roles/{userRoleId} [delete]
func (h *Handler) RemoveRoleFromUser(c *gin.Context) {
	if err := h.service.RemoveRoleFromUser(c.Param("userRoleId")); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusOK, "rol removido del usuario")
}

// CreateScopeOverride crea un permiso especial de scope.
// @Summary Crear override de scope
// @Security BearerAuth
// @Param body body CreateScopeOverrideRequest true "Detalle del override"
// @Success 201 {object} response.Body{data=role.ScopeOverrideResponse} "Override creado"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_SCOPE_MANAGE:write)"
// @Router /scope-overrides [post]
func (h *Handler) CreateScopeOverride(c *gin.Context) {
	var req CreateScopeOverrideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	item, err := h.service.CreateScopeOverride(req)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "override creado", item)
}

// ListScopeOverrides devuelve los overrides de un usuario.
// @Summary Listar overrides de scope de un usuario
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "UUID del usuario"
// @Success 200 {object} response.Body{data=role.ScopeOverrideResponse} "Overrides obtenidos"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_SCOPE_MANAGE:write)"
// @Failure 404 {object} response.Body "Usuario no encontrado"
// @Router /scope-overrides/users/{id} [get]
func (h *Handler) ListScopeOverrides(c *gin.Context) {
	items, err := h.service.ListScopeOverrides(c.Param("id"))
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "overrides obtenidos", gin.H{"items": items})
}

// RemoveScopeOverride elimina un override.
// @Summary Eliminar override de scope
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "UUID del override"
// @Success 200 {object} response.Body "Override eliminado"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_SCOPE_MANAGE:write)"
// @Failure 404 {object} response.Body "Override no encontrado"
// @Router /scope-overrides/{id} [delete]
func (h *Handler) RemoveScopeOverride(c *gin.Context) {
	if err := h.service.RemoveScopeOverride(c.Param("id")); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusOK, "override eliminado")
}