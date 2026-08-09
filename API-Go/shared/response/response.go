package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/shared/apperror"
)

// Body es la envoltura estándar de todas las respuestas del API.
type Body struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success devuelve una respuesta exitosa con datos.
func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Body{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessNoData devuelve una respuesta exitosa sin datos.
func SuccessNoData(c *gin.Context, status int, message string) {
	c.JSON(status, Body{
		Success: true,
		Message: message,
	})
}

// Error devuelve una respuesta de error con código de dominio.
func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, Body{
		Success: false,
		Message: message,
		Error: map[string]string{
			"code": code,
		},
	})
}

// AppError traduce un error de dominio (*apperror.Error) a su respuesta HTTP.
// Si el error no es de dominio, responde 500.
func AppError(c *gin.Context, err error) {
	e, ok := apperror.AsError(err)
	if !ok {
		c.JSON(http.StatusInternalServerError, Body{
			Success: false,
			Message: "error interno del servidor",
			Error: map[string]string{
				"code": "INTERNAL",
			},
		})
		return
	}

	c.JSON(e.Status, Body{
		Success: false,
		Message: e.Message,
		Error: map[string]string{
			"code": e.Code,
		},
	})
}