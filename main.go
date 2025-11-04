package main

import (
	"log"

	"github.com/AndresEGV/backend-api-go/database"
	"github.com/AndresEGV/backend-api-go/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Conectar a la base de datos
	database.Connect()

	// Crear el router de Gin
	router := gin.Default()

	// Ruta de bienvenida
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "¡Bienvenido a la API de Tareas en Go! 🚀",
			"version": "1.0.0",
			"endpoints": gin.H{
				"GET /tasks":        "Obtener todas las tareas",
				"GET /tasks/:id":    "Obtener una tarea por ID",
				"POST /tasks":       "Crear una nueva tarea",
				"PUT /tasks/:id":    "Actualizar una tarea",
				"DELETE /tasks/:id": "Eliminar una tarea",
			},
		})
	})

	// Agrupar las rutas de tareas bajo /tasks
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("", handlers.GetAllTasks)       // GET /tasks
		taskRoutes.GET("/:id", handlers.GetTaskByID)   // GET /tasks/:id
		taskRoutes.POST("", handlers.CreateTask)       // POST /tasks
		taskRoutes.PUT("/:id", handlers.UpdateTask)    // PUT /tasks/:id
		taskRoutes.DELETE("/:id", handlers.DeleteTask) // DELETE /tasks/:id
	}

	// Iniciar el servidor en el puerto 8080
	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
