package models

import "gorm.io/gorm"

// Task representa una tarea en nuestra aplicación
type Task struct {
	// gorm.Model incluye automáticamente: ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model

	// Título de la tarea (requerido)
	Title string `json:"title" binding:"required"`

	// Descripción de la tarea (opcional)
	Description string `json:"description"`

	// Estado: si la tarea está completada o no
	Completed bool `json:"completed" gorm:"default:false"`
}

// CreateTaskInput define los datos necesarios para crear una tarea
type CreateTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// UpdateTaskInput define los datos para actualizar una tarea
type UpdateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"` // Usamos puntero para permitir false como valor válido
}
