package handlers

import (
	"net/http"

	"github.com/AndresEGV/backend-api-go/database"
	"github.com/AndresEGV/backend-api-go/models"
	"github.com/gin-gonic/gin"
)

// GetAllTasks obtiene todas las tareas
// GET /tasks
func GetAllTasks(c *gin.Context) {
	var tasks []models.Task

	// Buscar todas las tareas en la base de datos
	result := database.DB.Find(&tasks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al obtener las tareas",
		})
		return
	}

	// Retornar las tareas en formato JSON
	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"count": len(tasks),
	})
}

// GetTaskByID obtiene una tarea por su ID
// GET /tasks/:id
func GetTaskByID(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	// Buscar la tarea por ID
	result := database.DB.First(&task, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tarea no encontrada",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": task,
	})
}

// CreateTask crea una nueva tarea
// POST /tasks
func CreateTask(c *gin.Context) {
	var input models.CreateTaskInput

	// Validar el JSON recibido
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Crear la tarea
	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
	}

	// Guardar en la base de datos
	result := database.DB.Create(&task)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al crear la tarea",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tarea creada exitosamente",
		"data":    task,
	})
}

// UpdateTask actualiza una tarea existente
// PUT /tasks/:id
func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	// Buscar la tarea
	result := database.DB.First(&task, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tarea no encontrada",
		})
		return
	}

	// Validar el input
	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Actualizar los campos
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
	database.DB.Save(&task)

	c.JSON(http.StatusOK, gin.H{
		"message": "Tarea actualizada exitosamente",
		"data":    task,
	})
}

// DeleteTask elimina una tarea
// DELETE /tasks/:id
func DeleteTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	// Buscar la tarea
	result := database.DB.First(&task, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tarea no encontrada",
		})
		return
	}

	// Eliminar la tarea
	database.DB.Delete(&task)

	c.JSON(http.StatusOK, gin.H{
		"message": "Tarea eliminada exitosamente",
	})
}
