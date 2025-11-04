package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/AndresEGV/backend-api-go/config"
	"github.com/AndresEGV/backend-api-go/database"
	"github.com/AndresEGV/backend-api-go/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// APIResponse estructura estándar
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code"`
}

// Health Check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{
		"status":   "ok",
		"database": "connected",
	})
}

// Welcome
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]interface{}{
		"message": "¡API de Tareas en Go con Chi Router! 🚀",
		"version": "2.0.0-chi",
		"features": map[string]string{
			"Framework":     "Chi (minimalista)",
			"Dependencies":  "1 pequeña",
			"Database":      "GORM + SQLite/PostgreSQL",
			"Compatibility": "100% stdlib compatible",
		},
	})
}

// Get All Tasks
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task

	if result := database.DB.Find(&tasks); result.Error != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, APIResponse{
			Success: false,
			Error:   "Error al obtener tareas",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	render.JSON(w, r, APIResponse{
		Success: true,
		Message: "Tareas obtenidas exitosamente",
		Data: map[string]interface{}{
			"tasks": tasks,
			"count": len(tasks),
		},
		Code: http.StatusOK,
	})
}

// Get Task by ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Chi extrae el parámetro de la URL
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, APIResponse{
			Success: false,
			Error:   "ID inválido",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var task models.Task
	if result := database.DB.First(&task, id); result.Error != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, APIResponse{
			Success: false,
			Error:   "Tarea no encontrada",
			Code:    http.StatusNotFound,
		})
		return
	}

	render.JSON(w, r, APIResponse{
		Success: true,
		Data:    task,
		Code:    http.StatusOK,
	})
}

// Create Task
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, APIResponse{
			Success: false,
			Error:   "JSON inválido: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validación
	if input.Title == "" || len(input.Title) < 3 || len(input.Title) > 200 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, APIResponse{
			Success: false,
			Error:   "Title debe tener entre 3 y 200 caracteres",
			Code:    http.StatusBadRequest,
		})
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
	}

	if result := database.DB.Create(&task); result.Error != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, APIResponse{
			Success: false,
			Error:   "Error al crear tarea",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, APIResponse{
		Success: true,
		Message: "Tarea creada exitosamente",
		Data:    task,
		Code:    http.StatusCreated,
	})
}

// Update Task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, APIResponse{
			Success: false,
			Error:   "ID inválido",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var task models.Task
	if result := database.DB.First(&task, id); result.Error != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, APIResponse{
			Success: false,
			Error:   "Tarea no encontrada",
			Code:    http.StatusNotFound,
		})
		return
	}

	var input models.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, APIResponse{
			Success: false,
			Error:   "JSON inválido",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = input.Description
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	database.DB.Save(&task)

	render.JSON(w, r, APIResponse{
		Success: true,
		Message: "Tarea actualizada exitosamente",
		Data:    task,
		Code:    http.StatusOK,
	})
}

// Delete Task
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, APIResponse{
			Success: false,
			Error:   "ID inválido",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var task models.Task
	if result := database.DB.First(&task, id); result.Error != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, APIResponse{
			Success: false,
			Error:   "Tarea no encontrada",
			Code:    http.StatusNotFound,
		})
		return
	}

	database.DB.Delete(&task)

	render.JSON(w, r, APIResponse{
		Success: true,
		Message: "Tarea eliminada exitosamente",
		Code:    http.StatusOK,
	})
}

func main() {
	// 1. Cargar configuración
	log.Println("🔧 Cargando configuración...")
	config.Load()

	// 2. Conectar a base de datos
	log.Println("🔌 Conectando a la base de datos...")
	database.Connect()

	// 3. Crear router Chi
	r := chi.NewRouter()

	// 4. Middleware
	r.Use(middleware.Logger)        // Logger de Chi
	r.Use(middleware.Recoverer)     // Recovery de panics
	r.Use(middleware.RequestID)     // Request ID único
	r.Use(middleware.RealIP)        // IP real del cliente
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// 5. Rutas
	r.Get("/health", healthHandler)
	r.Get("/", welcomeHandler)

	// Rutas de tareas - Con path parameters estilo REST
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", getTasksHandler)        // GET /tasks
		r.Post("/", createTaskHandler)     // POST /tasks
		r.Get("/{id}", getTaskHandler)     // GET /tasks/1
		r.Put("/{id}", updateTaskHandler)  // PUT /tasks/1
		r.Delete("/{id}", deleteTaskHandler) // DELETE /tasks/1
	})

	// 6. Iniciar servidor
	port := config.AppConfig.Server.Port
	addr := ":" + port

	log.Printf("🚀 Servidor corriendo en http://localhost:%s", port)
	log.Printf("📚 Documentación: http://localhost:%s/", port)
	log.Printf("🏥 Health check: http://localhost:%s/health", port)
	log.Println("✨ Usando Chi Router (minimalista + stdlib compatible)")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("❌ Error al iniciar el servidor:", err)
	}
}
