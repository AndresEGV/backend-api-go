package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse estructura estándar para respuestas de la API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code"`
}

// SuccessResponse respuesta exitosa
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Code:    statusCode,
	})
}

// ErrorResponse respuesta de error
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
	})
}

// ValidationErrorResponse respuesta de error de validación
func ValidationErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Error:   "Error de validación: " + err.Error(),
		Code:    http.StatusBadRequest,
	})
}

// NotFoundResponse respuesta de recurso no encontrado
func NotFoundResponse(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Error:   resource + " no encontrado",
		Code:    http.StatusNotFound,
	})
}

// InternalServerErrorResponse respuesta de error interno del servidor
func InternalServerErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Error:   "Error interno del servidor: " + message,
		Code:    http.StatusInternalServerError,
	})
}
