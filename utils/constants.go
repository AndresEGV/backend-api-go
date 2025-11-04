package utils

// Mensajes de éxito
const (
	MsgTaskCreated = "Tarea creada exitosamente"
	MsgTaskUpdated = "Tarea actualizada exitosamente"
	MsgTaskDeleted = "Tarea eliminada exitosamente"
	MsgTasksFetched = "Tareas obtenidas exitosamente"
)

// Mensajes de error
const (
	ErrTaskNotFound = "Tarea no encontrada"
	ErrInvalidInput = "Datos de entrada inválidos"
	ErrDatabaseOperation = "Error en operación de base de datos"
	ErrInternalServer = "Error interno del servidor"
)

// Límites y validaciones
const (
	MaxTitleLength = 200
	MinTitleLength = 3
	MaxDescriptionLength = 1000
)
