package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/middleware"
	"github.com/jesus-ariel/api-jeussairel/shared/response"
)

// Handler expone los endpoints de autenticación a través de HTTP.
type Handler struct {
	service *Service
}

// NewHandler crea el handler de autenticación.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register crea una cuenta de forma pública (sin autenticación).
// @Summary Registrar cuenta
// @Description Crea un usuario nuevo (USER, INSTRUCTOR o LEARNER) pendiente de verificación: se envía un código de 6 dígitos por email. No requiere token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Datos del nuevo usuario"
// @Success 201 {object} response.Body "Cuenta creada (pendiente de verificar email)"
// @Failure 400 {object} response.Body "Datos inválidos"
// @Failure 409 {object} response.Body "Email ya registrado"
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	u, devCode, err := h.service.Register(req, c.ClientIP())
	if err != nil {
		response.AppError(c, err)
		return
	}

	data := gin.H{"user": u}
	if devCode != "" {
		data["dev_code"] = devCode
		data["message"] = "en desarrollo: usa este código para verificar tu cuenta"
	}

	response.Success(c, http.StatusCreated, "cuenta creada, revisa tu email para verificar la cuenta", data)
}

// VerifyEmail confirma la cuenta con el código de verificación.
// @Summary Verificar email
// @Description Activa la cuenta usando el código de 6 dígitos enviado por correo durante el registro.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body VerifyEmailRequest true "Email y código de verificación"
// @Success 200 {object} response.Body "Cuenta verificada"
// @Failure 400 {object} response.Body "Código inválido o expirado"
// @Router /auth/verify-email [post]
func (h *Handler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	if err := h.service.VerifyEmail(req, c.ClientIP()); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusOK, "cuenta verificada, ya puedes iniciar sesión")
}

// Login autentica al usuario y devuelve el par de tokens.
// @Summary Iniciar sesión
// @Description Autentica al usuario con email y contraseña y devuelve el par de tokens (access + refresh).
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Credenciales de acceso"
// @Success 200 {object} response.Body{data=auth.TokenResponse} "Inicio de sesión exitoso"
// @Failure 400 {object} response.Body "Datos inválidos"
// @Failure 401 {object} response.Body "Credenciales inválidas"
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	tokens, err := h.service.Login(req, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "inicio de sesión exitoso", tokens)
}

// Refresh renueva el par de tokens con un refresh token válido.
// @Summary Renovar tokens
// @Description Renueva el access token (y el refresh token) usando un refresh token válido.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RefreshRequest true "Refresh token"
// @Success 200 {object} response.Body{data=auth.TokenResponse} "Tokens renovados"
// @Failure 400 {object} response.Body "Datos inválidos"
// @Failure 401 {object} response.Body "Refresh token inválido o expirado"
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	tokens, err := h.service.Refresh(req, c.ClientIP())
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "tokens renovados", tokens)
}

// Logout revoca el refresh token.
// @Summary Cerrar sesión
// @Description Revoca el refresh token y cierra la sesión activa.
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body LogoutRequest true "Refresh token"
// @Success 200 {object} response.Body "Sesión cerrada"
// @Failure 400 {object} response.Body "Datos inválidos"
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	if err := h.service.Logout(req); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusOK, "sesión cerrada")
}

// ForgotPassword genera una solicitud de reseteo de contraseña.
// @Summary Solicitar recuperación de contraseña
// @Description Genera una solicitud de reseteo. Si el email existe, se envía un código de 6 dígitos. En desarrollo sin SMTP, devuelve el código en la respuesta.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body ForgotPasswordRequest true "Email del usuario"
// @Success 200 {object} response.Body "Si el email existe, se envió la instrucción"
// @Failure 400 {object} response.Body "Datos inválidos"
// @Router /auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	devCode, err := h.service.ForgotPassword(req, c.ClientIP())
	if err != nil {
		response.AppError(c, err)
		return
	}

	if devCode != "" {
		response.Success(c, http.StatusOK, "en desarrollo: usa este código para restablecer tu contraseña", gin.H{"code": devCode})
		return
	}

	response.SuccessNoData(c, http.StatusOK, "si el email existe, se envió la instrucción")
}

// ResetPassword cambia la contraseña con el token de recuperación.
// @Summary Restablecer contraseña
// @Description Cambia la contraseña usando el token de recuperación generado en forgot-password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body ResetPasswordRequest true "Código y nueva contraseña"
// @Success 200 {object} response.Body "Contraseña actualizada"
// @Failure 400 {object} response.Body "Datos inválidos"
// @Failure 401 {object} response.Body "Código inválido o expirado"
// @Router /auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "datos inválidos en la petición")
		return
	}

	if err := h.service.ResetPassword(req); err != nil {
		response.AppError(c, err)
		return
	}

	response.SuccessNoData(c, http.StatusOK, "contraseña actualizada")
}

// Me devuelve la sesión del usuario autenticado (requiere token).
// @Summary Perfil y permisos del usuario autenticado
// @Description Devuelve los datos del usuario logueado, sus roles activos y los features/scope que tiene permitidos usar en la UI.
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Body{data=auth.MeResponse} "Sesión del usuario"
// @Failure 401 {object} response.Body "Token inválido o expirado"
// @Router /auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)

	me, err := h.service.Me(userID.(string))
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "sesión obtenida", me)
}