package handlers

import (
	"net/http"

	"github.com/AndresEGV/backend-api-go/database"
	"github.com/AndresEGV/backend-api-go/models"
	"github.com/AndresEGV/backend-api-go/utils"
	"github.com/gin-gonic/gin"
)

// GetAllTasks obtiene todas las tareas
// GET /tasks
func GetAllTasks(c *gin.Context) {
	var tasks []models.Task

	// Buscar todas las tareas en la base de datos
	if result := database.DB.Find(&tasks); result.Error != nil {
		utils.InternalServerErrorResponse(c, utils.ErrDatabaseOperation)
		return
	}

	// Retornar las tareas en formato JSON
	utils.SuccessResponse(c, http.StatusOK, utils.MsgTasksFetched, gin.H{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// GetTaskByID obtiene una tarea por su ID
// GET /tasks/:id
func GetTaskByID(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	// Buscar la tarea por ID
	if result := database.DB.First(&task, id); result.Error != nil {
		utils.NotFoundResponse(c, "Tarea")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", task)
}

// CreateTask crea una nueva tarea
// POST /tasks
func CreateTask(c *gin.Context) {
	var input models.CreateTaskInput

	// Validar el JSON recibido
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Crear la tarea
	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
	}

	// Guardar en la base de datos
	if result := database.DB.Create(&task); result.Error != nil {
		utils.InternalServerErrorResponse(c, utils.ErrDatabaseOperation)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, utils.MsgTaskCreated, task)
}

// UpdateTask actualiza una tarea existente
// PUT /tasks/:id
func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	// Buscar la tarea
	if result := database.DB.First(&task, id); result.Error != nil {
		utils.NotFoundResponse(c, "Tarea")
		return
	}

	// Validar el input
	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Actualizar los campos solo si se proporcionaron
	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = input.Description
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	// Guardar cambios
	if result := database.DB.Save(&task); result.Error != nil {
		utils.InternalServerErrorResponse(c, utils.ErrDatabaseOperation)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.MsgTaskUpdated, task)
}

// DeleteTask elimina una tarea (soft delete)
// DELETE /tasks/:id
func DeleteTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	// Buscar la tarea
	if result := database.DB.First(&task, id); result.Error != nil {
		utils.NotFoundResponse(c, "Tarea")
		return
	}

	// Eliminar la tarea (soft delete por defecto con GORM)
	if result := database.DB.Delete(&task); result.Error != nil {
		utils.InternalServerErrorResponse(c, utils.ErrDatabaseOperation)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.MsgTaskDeleted, nil)
}
