package feature

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/shared/response"
)

// Handler expone los endpoints de módulos y funcionalidades a través de HTTP.
type Handler struct {
	service *Service
}

// NewHandler crea el handler de features.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateModule registra un módulo.
// @Summary Crear módulo
// @Tags Modules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateModuleRequest true "Datos del módulo"
// @Success 201 {object} response.Body{data=feature.ModuleResponse} "Módulo creado"
// @Failure 400 {object} response.Body "Datos inválidos"
// @Failure 403 {object} response.Body "Sin permisos (IDENTITY_ROLE_MANAGE:write)"
// @Router /modules [post]
func (h *Handler) CreateModule(c *gin.Context) {
	var req CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	m, err := h.service.CreateModule(req)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "módulo creado", m)
}

// GetModuleById devuelve un módulo.
// @Summary Obtener módulo por ID
// @Tags Modules
// @Security BearerAuth
// @Param id path string true "UUID del módulo"
// @Success 200 {object} response.Body{data=feature.ModuleResponse} "Módulo encontrado"
// @Failure 403 {object} response.Body "Acceso denegado (IDENTITY_ROLE_MANAGE:write)"
// @Failure 404 {object} response.Body "Módulo no encontrado"
// @Router /modules/{id} [get]
func (h *Handler) GetModuleById(c *gin.Context) {
	m, err := h.service.GetModuleById(c.Param("id"))
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "módulo encontrado", m)
}

// ListModules devuelve módulos paginados.
// @Summary Listar módulos
// @Tags API
// @Security BearerAuth
// @Param page query int false "Número de página (default 1)"
// @Param page_size query int false "Tamaño de página (1-100, default 20)"
// @Success 200 {object} response.Body{data=feature.ListModulesResponse} "Módulos obtenidos"
// @Failure 403 {object} response.Body "Acceso denegado (IDENTITY_ROLE_MANAGE:write)"
// @Router /modules [get]
func (h *Handler) ListModules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	modules, total, err := h.service.ListModules(page, pageSize)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "módulos obtenidos", gin.H{
		"items": modules,
		"total": total,
		"page":  page,
	})
}

// UpdateModule actualiza un módulo.
// @Summary Actualizar módulo
// @Security BearerAuth
// @Param id path string true "UUID del módulo"
// @Param body body UpdateModuleRequest true "Campos a actualizar"
// @Success 200 {object} response.Body{data=feature.ModuleResponse} "Módulo actualizado"
// @Failure 403 {object} response.Body "Acceso denegado (IDENTITY_ROLE_MANAGE:write)"
// @Failure 404 {object} response.Body "Módulo no encontrado"
// @Router /modules/{id} [put]
func (h *Handler) UpdateModule(c *gin.Context) {
	var req UpdateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	m, err := h.service.UpdateModule(c.Param("id"), req)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "módulo actualizado", m)
}

// DeleteModule elimina un módulo.
// @Summary Eliminar módulo
// @Tags API
// @Security BearerAuth
// @Param id path string true "UUID del módulo"
// @Success 200 {object} response.Body "Módulo eliminado"
// @Failure 403 {object} response.Body "Acceso denegado (IDENTITY_ROLE_MANAGE:write)"
// @Failure 404 {object} response.Body "Módulo no encontrado"
// @Router /modules/{id} [delete]
func (h *Handler) DeleteModule(c *gin.Context) {
	if err := h.service.DeleteModule(c.Param("id")); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusOK, "módulo eliminado")
}

// ListFeaturesByModule devuelve las funcionalidades de un módulo.
// @Summary Listar funcionalidades de un módulo
// @Tags API
// @Security BearerAuth
// @Param id path string true "UUID del módulo"
// @Success 200 {object} response.Body{data=feature.ListFeaturesResponse} "Funcionalidades del módulo"
// @Failure 403 {object} response.Body "Acceso denegado (IDENTITY_ROLE_MANAGE:write)"
// @Failure 404 {object} response.Body "Módulo no encontrado"
// @Router /modules/{id}/features [get]
func (h *Handler) ListFeaturesByModule(c *gin.Context) {
	features, err := h.service.ListFeaturesByModule(c.Param("id"))
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "funcionalidades del módulo", gin.H{"items": features})
}

// CreateFeature CRUD funcionalidades.
// @Summary Crear funcionalidad
// @Tags API
// @Security BearerAuth
// @Param body body CreateFeatureRequest true "Datos de la funcionalidad"
// @Success 201 {object} response.Body{data=feature.FeatureResponse} "Funcionalidad creada"
// @Failure 400 {object} response.Body "Datos inválidos"
// @Failure 403 {object} response.Body "Acceso denegado (IDENTITY_ROLE_MANAGE:write)"
// @Router /features [post]
func (h *Handler) CreateFeature(c *gin.Context) {
	var req CreateFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	f, err := h.service.CreateFeature(req)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "funcionalidad creada", f)
}

// GetFeatureById devuelve una funcionalidad.
// @Summary Obtener funcionalidad por ID
// @Tags API
// @Security BearerAuth
// @Param id path string true "UUID de la funcionalidad"
// @Success 200 {object} response.Body{data=feature.FeatureResponse} "Funcionalidad encontrada"
// @Failure 403 {object} response.Body "Acceso denegado (IDENTITY_ROLE_MANAGE:write)"
// @Failure 404 {object} response.Body "Funcionalidad no encontrada"
// @Router /features/{id} [get]
func (h *Handler) GetFeatureById(c *gin.Context) {
	f, err := h.service.GetFeatureById(c.Param("id"))
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "funcionalidad encontrada", f)
}

// ListFeatures devuelve funcionalidades paginadas.
// @Summary Listar funcionalidades
// @Tags API
// @Security BearerAuth
// @Param page query int false "Número de página (default 1)"
// @Param page_size query int false "Tamaño de página (1-100, default 20)"
// @Success 200 {object} response.Body{data=feature.ListFeaturesResponse} "Funcionalidades obtenidas"
// @Failure 403 {object} response.Body "Acceso denegado (IDENTITY_ROLE_MANAGE:write)"
// @Router /features [get]
func (h *Handler) ListFeatures(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	features, total, err := h.service.ListFeatures(page, pageSize)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "funcionalidades obtenidas", gin.H{
		"items": features,
		"total": total,
		"page":  page,
	})
}

// UpdateFeature actualiza una funcionalidad.
// @Summary Actualizar funcionalidad
// @Tags API
// @Security BearerAuth
// @Param id path string true "UUID de la funcionalidad"
// @Param body body UpdateFeatureRequest true "Campos a actualizar"
// @Success 200 {object} response.Body{data=feature.FeatureResponse} "Funcionalidad actualizada"
// @Failure 403 {object} response.Body "Acceso denegado (IDENTITY_ROLE_MANAGE:write)"
// @Failure 404 {object} response.Body "Funcionalidad no encontrada"
// @Router /features/{id} [put]
func (h *Handler) UpdateFeature(c *gin.Context) {
	var req UpdateFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	f, err := h.service.UpdateFeature(c.Param("id"), req)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "funcionalidad actualizada", f)
}

// DeleteFeature elimina una funcionalidad.
// @Summary Eliminar funcionalidad
// @Tags API
// @Security BearerAuth
// @Param id path string true "UUID de la funcionalidad"
// @Success 200 {object} response.Body "Funcionalidad eliminada"
// @Failure 403 {object} response.Body "Acceso denegado (IDENTITY_ROLE_MANAGE:write)"
// @Failure 404 {object} response.Body "Funcionalidad no encontrada"
// @Router /features/{id} [delete]
func (h *Handler) DeleteFeature(c *gin.Context) {
	if err := h.service.DeleteFeature(c.Param("id")); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusOK, "funcionalidad eliminada")
}