# API REST en Go - Tutorial Completo

Este es un tutorial paso a paso para aprender a crear una API REST en Go usando Gin y GORM.

## ¿Qué tecnologías usamos?

### 1. **Go (Golang)**
Lenguaje de programación compilado, eficiente y con concurrencia nativa.

### 2. **Gin Framework**
Framework web rápido y minimalista para Go. Maneja las rutas HTTP, validación JSON, y respuestas.

### 3. **GORM (ORM)**
**Object-Relational Mapping** - Te permite trabajar con bases de datos usando structs de Go en lugar de SQL directo.

**Ventajas del ORM:**
- ✅ No necesitas escribir SQL manualmente
- ✅ Automáticamente crea/actualiza tablas (migrations)
- ✅ Previene inyección SQL
- ✅ Código más limpio y mantenible

**Ejemplo sin ORM (SQL puro):**
```go
db.Exec("INSERT INTO tasks (title, description) VALUES (?, ?)", title, desc)
```

**Ejemplo con ORM (GORM):**
```go
db.Create(&Task{Title: title, Description: desc})
```

### 4. **SQLite**
Base de datos ligera en un solo archivo. Perfecta para desarrollo y aplicaciones pequeñas.

---

## Estructura del Proyecto

```
backend-api-go/
├── main.go              # Punto de entrada, configura rutas y servidor
├── models/              # Modelos de datos (structs)
│   └── task.go
├── database/            # Configuración de la base de datos
│   └── database.go
├── handlers/            # Controladores (lógica de negocio)
│   └── task_handler.go
├── tasks.db             # Base de datos SQLite (se crea automáticamente)
└── go.mod               # Dependencias del proyecto
```

---

## Conceptos Clave en Go

### 1. **Structs**
Los structs son como "clases" en otros lenguajes. Definen la estructura de tus datos.

```go
type Task struct {
    gorm.Model              // Añade ID, CreatedAt, UpdatedAt, DeletedAt
    Title       string      // Campo título
    Description string      // Campo descripción
    Completed   bool        // Estado de completado
}
```

### 2. **Tags (Etiquetas)**
Los tags entre `` definen comportamiento especial:

```go
Title string `json:"title" binding:"required"`
//            ↑                ↑
//         Nombre en JSON    Validación
```

- `json:"title"` → El campo se llamará "title" en JSON
- `binding:"required"` → El campo es obligatorio
- `gorm:"default:false"` → Valor por defecto en la BD

### 3. **Punteros**
En Go, `*` indica un puntero. Útil para diferenciar "sin valor" de "valor false":

```go
Completed *bool  // Puede ser: true, false, o nil (sin valor)
Completed bool   // Solo puede ser: true o false
```

---

## Cómo funciona el flujo de la API

```
Cliente (Postman/Browser)
    ↓
Petición HTTP (GET, POST, PUT, DELETE)
    ↓
Router (Gin) → identifica la ruta
    ↓
Handler (función controladora)
    ↓
Base de datos (GORM + SQLite)
    ↓
Respuesta JSON al cliente
```

---

## Endpoints de la API

### 1. **GET /** - Información de la API
```bash
curl http://localhost:8080/
```

### 2. **GET /tasks** - Obtener todas las tareas
```bash
curl http://localhost:8080/tasks
```

### 3. **GET /tasks/:id** - Obtener una tarea específica
```bash
curl http://localhost:8080/tasks/1
```

### 4. **POST /tasks** - Crear una nueva tarea
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Aprender Go",
    "description": "Completar el tutorial de API REST"
  }'
```

### 5. **PUT /tasks/:id** - Actualizar una tarea
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Aprender Go - Actualizado",
    "completed": true
  }'
```

### 6. **DELETE /tasks/:id** - Eliminar una tarea
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

---

## Cómo ejecutar el proyecto

### 1. **Instalar dependencias**
```bash
go mod download
```

### 2. **Ejecutar el servidor**
```bash
go run main.go
```

Verás:
```
✅ Conexión a la base de datos establecida
✅ Migraciones ejecutadas correctamente
🚀 Servidor corriendo en http://localhost:8080
```

### 3. **Probar la API**

Puedes usar:
- **curl** (desde la terminal)
- **Postman** (aplicación GUI)
- **Thunder Client** (extensión de VS Code)
- **Navegador** (solo para GET)

---

## Explicación de cada archivo

### **models/task.go**
Define la estructura de datos de una tarea. Es como crear una tabla en SQL pero con código Go.

**Conceptos:**
- `gorm.Model` → Añade campos automáticos: ID, CreatedAt, UpdatedAt, DeletedAt
- `CreateTaskInput` → Define qué datos se necesitan para crear una tarea
- `UpdateTaskInput` → Define qué datos se pueden actualizar

### **database/database.go**
Configura la conexión a la base de datos.

**Conceptos:**
- `gorm.Open()` → Conecta a la base de datos
- `AutoMigrate()` → Crea/actualiza las tablas automáticamente basándose en los modelos

### **handlers/task_handler.go**
Contiene la lógica de negocio (CRUD).

**Funciones:**
- `GetAllTasks()` → SELECT * FROM tasks
- `GetTaskByID()` → SELECT * FROM tasks WHERE id = ?
- `CreateTask()` → INSERT INTO tasks
- `UpdateTask()` → UPDATE tasks SET ...
- `DeleteTask()` → DELETE FROM tasks WHERE id = ?

### **main.go**
Punto de entrada de la aplicación.

**Funciones:**
- Conecta a la base de datos
- Define las rutas HTTP
- Inicia el servidor en el puerto 8080

---

## Conceptos importantes de GORM

### Operaciones CRUD con GORM

```go
// CREATE (Crear)
db.Create(&task)

// READ (Leer)
db.Find(&tasks)           // Obtener todos
db.First(&task, id)       // Obtener por ID

// UPDATE (Actualizar)
db.Save(&task)            // Guardar cambios

// DELETE (Eliminar)
db.Delete(&task)          // Soft delete (no elimina realmente)
db.Unscoped().Delete(&task) // Hard delete (elimina permanentemente)
```

### Auto Migrations
GORM crea/actualiza las tablas automáticamente:

```go
db.AutoMigrate(&models.Task{})
```

Esto crea una tabla como:
```sql
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE
);
```

---

## Próximos pasos para aprender más

1. **Agregar validaciones personalizadas**
   - Validar longitud de título
   - Validar formato de email, etc.

2. **Agregar middleware**
   - Autenticación con JWT
   - Logging de peticiones
   - CORS para frontend

3. **Usar PostgreSQL o MySQL**
   - Cambiar de SQLite a base de datos más robusta
   ```go
   import "gorm.io/driver/postgres"
   dsn := "host=localhost user=gorm password=gorm dbname=gorm"
   db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
   ```

4. **Agregar relaciones**
   - Tareas con usuarios (1 usuario → N tareas)
   - Tareas con etiquetas (N tareas ↔ N etiquetas)

5. **Testing**
   - Escribir tests unitarios
   - Tests de integración

6. **Deploy**
   - Dockerizar la aplicación
   - Desplegar en Railway, Fly.io, o AWS

---

## Recursos adicionales

- [Documentación oficial de Go](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Go by Example](https://gobyexample.com/)

---

## ¿Preguntas frecuentes?

**P: ¿Por qué usar un ORM?**
R: Facilita el trabajo con bases de datos, previene errores SQL, y hace el código más mantenible.

**P: ¿Es Go difícil de aprender?**
R: Go es uno de los lenguajes más simples. Tiene pocas palabras clave y una sintaxis limpia.

**P: ¿Cuándo usar Go vs otros lenguajes?**
R: Go es excelente para APIs, microservicios, herramientas CLI y sistemas distribuidos.

**P: ¿Necesito saber SQL para usar GORM?**
R: No es estrictamente necesario, pero entender SQL te ayudará a usar GORM mejor.
