package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware personalizado para logging detallado
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tiempo de inicio
		startTime := time.Now()

		// Procesar la petición
		c.Next()

		// Tiempo de finalización
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Obtener información de la petición
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// Log detallado
		log.Printf("[%s] %s | Status: %d | Latency: %v | IP: %s",
			method,
			path,
			statusCode,
			latency,
			clientIP,
		)
	}
}
