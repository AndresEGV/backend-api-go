package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AndresEGV/backend-api-go/config"
	"github.com/AndresEGV/backend-api-go/database"
	"github.com/AndresEGV/backend-api-go/models"
)

// APIResponse estructura estándar para respuestas
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code"`
}

// respondJSON helper para respuestas JSON
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// Middleware: Logger
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("[%s] %s | Latency: %v | IP: %s",
			r.Method, r.URL.Path, time.Since(start), r.RemoteAddr)
	}
}

// Middleware: CORS
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Handler: Health Check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":   "ok",
		"database": "connected",
	})
}

// Handler: Welcome
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "¡API de Tareas en Go con Standard Library! 🚀",
		"version": "2.0.0-stdlib",
		"features": map[string]string{
			"Framework":   "None (stdlib only)",
			"Dependencies": "0 external",
			"Database":    "GORM + SQLite/PostgreSQL",
		},
		"endpoints": map[string]string{
			"GET /health":       "Health check",
			"GET /tasks":        "Obtener todas las tareas",
			"GET /tasks?id=1":   "Obtener una tarea",
			"POST /tasks":       "Crear nueva tarea",
			"PUT /tasks?id=1":   "Actualizar tarea",
			"DELETE /tasks?id=1": "Eliminar tarea",
		},
	})
}

// Handler: Get All Tasks
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task

	if result := database.DB.Find(&tasks); result.Error != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Error al obtener tareas",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Tareas obtenidas exitosamente",
		Data: map[string]interface{}{
			"tasks": tasks,
			"count": len(tasks),
		},
		Code: http.StatusOK,
	})
}

// Handler: Get Task by ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "ID es requerido",
			Code:    http.StatusBadRequest,
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "ID inválido",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var task models.Task
	if result := database.DB.First(&task, id); result.Error != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Tarea no encontrada",
			Code:    http.StatusNotFound,
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    task,
		Code:    http.StatusOK,
	})
}

// Handler: Create Task
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "JSON inválido: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validación manual
	if input.Title == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Title es requerido",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if len(input.Title) < 3 || len(input.Title) > 200 {
		respondJSON(w, http.StatusBadRequest, APIResponse{
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
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Error al crear tarea",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "Tarea creada exitosamente",
		Data:    task,
		Code:    http.StatusCreated,
	})
}

// Handler: Update Task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "ID es requerido",
			Code:    http.StatusBadRequest,
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "ID inválido",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var task models.Task
	if result := database.DB.First(&task, id); result.Error != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Tarea no encontrada",
			Code:    http.StatusNotFound,
		})
		return
	}

	var input models.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "JSON inválido: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Actualizar campos
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

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Tarea actualizada exitosamente",
		Data:    task,
		Code:    http.StatusOK,
	})
}

// Handler: Delete Task
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "ID es requerido",
			Code:    http.StatusBadRequest,
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "ID inválido",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var task models.Task
	if result := database.DB.First(&task, id); result.Error != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Tarea no encontrada",
			Code:    http.StatusNotFound,
		})
		return
	}

	database.DB.Delete(&task)

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Tarea eliminada exitosamente",
		Code:    http.StatusOK,
	})
}

// Aplicar middleware a un handler
func withMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return corsMiddleware(loggingMiddleware(h))
}

func main() {
	// 1. Cargar configuración
	log.Println("🔧 Cargando configuración...")
	config.Load()

	// 2. Conectar a base de datos
	log.Println("🔌 Conectando a la base de datos...")
	database.Connect()

	// 3. Crear router (ServeMux)
	mux := http.NewServeMux()

	// 4. Registrar rutas
	mux.HandleFunc("/health", withMiddleware(healthHandler))
	mux.HandleFunc("/", withMiddleware(welcomeHandler))

	// Rutas de tareas
	mux.HandleFunc("/tasks", withMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				getTaskHandler(w, r)
			} else {
				getTasksHandler(w, r)
			}
		case http.MethodPost:
			createTaskHandler(w, r)
		case http.MethodPut:
			updateTaskHandler(w, r)
		case http.MethodDelete:
			deleteTaskHandler(w, r)
		default:
			respondJSON(w, http.StatusMethodNotAllowed, APIResponse{
				Success: false,
				Error:   "Método no permitido",
				Code:    http.StatusMethodNotAllowed,
			})
		}
	}))

	// 5. Iniciar servidor
	port := config.AppConfig.Server.Port
	addr := ":" + port

	log.Printf("🚀 Servidor corriendo en http://localhost:%s", port)
	log.Printf("📚 Documentación: http://localhost:%s/", port)
	log.Printf("🏥 Health check: http://localhost:%s/health", port)
	log.Println("✨ Usando SOLO Standard Library (0 dependencias externas)")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("❌ Error al iniciar el servidor:", err)
	}
}
