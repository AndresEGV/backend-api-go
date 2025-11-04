package main

import (
	"fmt"
	"log"

	"github.com/AndresEGV/backend-api-go/config"
	"github.com/AndresEGV/backend-api-go/database"
	"github.com/AndresEGV/backend-api-go/handlers"
	"github.com/AndresEGV/backend-api-go/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Cargar configuración
	log.Println("🔧 Cargando configuración...")
	config.Load()

	// 2. Configurar modo de Gin (debug/release)
	gin.SetMode(config.AppConfig.Server.GinMode)

	// 3. Conectar a la base de datos
	log.Println("🔌 Conectando a la base de datos...")
	database.Connect()

	// 4. Crear el router de Gin sin middleware por defecto
	router := gin.New()

	// 5. Aplicar middleware globales
	router.Use(gin.Recovery())        // Recuperación de panics
	router.Use(middleware.Logger())   // Logger personalizado
	router.Use(middleware.CORS())     // CORS configurado

	// 6. Ruta de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"database": "connected",
		})
	})

	// 7. Ruta de bienvenida
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "¡Bienvenido a la API de Tareas en Go! 🚀",
			"version": "2.0.0",
			"features": gin.H{
				"ORM":         "GORM",
				"Database":    config.AppConfig.Database.Type,
				"CORS":        "Habilitado",
				"Validations": "Automáticas",
			},
			"endpoints": gin.H{
				"GET /health":       "Health check",
				"GET /tasks":        "Obtener todas las tareas",
				"GET /tasks/:id":    "Obtener una tarea por ID",
				"POST /tasks":       "Crear una nueva tarea",
				"PUT /tasks/:id":    "Actualizar una tarea",
				"DELETE /tasks/:id": "Eliminar una tarea",
			},
		})
	})

	// 8. Agrupar las rutas de tareas bajo /tasks
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("", handlers.GetAllTasks)       // GET /tasks
		taskRoutes.GET("/:id", handlers.GetTaskByID)   // GET /tasks/:id
		taskRoutes.POST("", handlers.CreateTask)       // POST /tasks
		taskRoutes.PUT("/:id", handlers.UpdateTask)    // PUT /tasks/:id
		taskRoutes.DELETE("/:id", handlers.DeleteTask) // DELETE /tasks/:id
	}

	// 9. Iniciar el servidor
	port := config.AppConfig.Server.Port
	serverAddr := fmt.Sprintf(":%s", port)

	log.Printf("🚀 Servidor corriendo en http://localhost:%s", port)
	log.Printf("📚 Documentación: http://localhost:%s/", port)
	log.Printf("🏥 Health check: http://localhost:%s/health", port)

	if err := router.Run(serverAddr); err != nil {
		log.Fatal("❌ Error al iniciar el servidor:", err)
	}
}
