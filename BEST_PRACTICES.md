# Buenas Prácticas - API REST en Go

Este documento explica las buenas prácticas implementadas en este proyecto y por qué son importantes.

## 📋 Tabla de Contenidos

1. [Estructura del Proyecto](#estructura-del-proyecto)
2. [Configuración Centralizada](#configuración-centralizada)
3. [Manejo de Errores](#manejo-de-errores)
4. [Validaciones](#validaciones)
5. [Middleware](#middleware)
6. [Base de Datos](#base-de-datos)
7. [Respuestas HTTP Estandarizadas](#respuestas-http-estandarizadas)
8. [Seguridad](#seguridad)
9. [Testing](#testing)
10. [Deploy](#deploy)

---

## 1. Estructura del Proyecto

### ✅ Buena Práctica: Separación de responsabilidades

```
backend-api-go/
├── config/          # Configuración centralizada
├── database/        # Conexión y configuración de BD
├── handlers/        # Controladores (lógica de negocio)
├── middleware/      # Middleware personalizados
├── models/          # Modelos de datos
├── utils/           # Utilidades y helpers
└── main.go          # Punto de entrada
```

**Por qué:**
- ✅ Código más organizado y mantenible
- ✅ Fácil de escalar y agregar nuevas funcionalidades
- ✅ Cada paquete tiene una responsabilidad clara
- ✅ Facilita el testing unitario

**❌ Evitar:**
```go
// NO: Todo en un solo archivo main.go de 1000+ líneas
```

---

## 2. Configuración Centralizada

### ✅ Buena Práctica: Variables de entorno con `.env`

**Archivo:** `config/config.go`

```go
// Usar variables de entorno para configuración
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    CORS     CORSConfig
}
```

**Por qué:**
- ✅ Configuración diferente para dev/staging/prod
- ✅ No hardcodear valores sensibles en el código
- ✅ Fácil cambiar configuración sin recompilar
- ✅ Secretos fuera del código fuente

**❌ Evitar:**
```go
// NO: Hardcodear configuración
db, _ := gorm.Open(sqlite.Open("tasks.db"))
router.Run(":8080")
```

### Uso correcto:

```go
// ✅ SI: Usar configuración centralizada
config.Load()
db, _ := gorm.Open(sqlite.Open(config.AppConfig.Database.Name))
router.Run(":" + config.AppConfig.Server.Port)
```

---

## 3. Manejo de Errores

### ✅ Buena Práctica: Respuestas de error consistentes

**Archivo:** `utils/response.go`

```go
// Funciones estandarizadas para errores
ErrorResponse(c, statusCode, message)
ValidationErrorResponse(c, err)
NotFoundResponse(c, resource)
InternalServerErrorResponse(c, message)
```

**Por qué:**
- ✅ API consistente para el cliente
- ✅ Códigos HTTP apropiados
- ✅ Mensajes de error claros
- ✅ Estructura uniforme de respuestas

**Ejemplo:**
```go
// ✅ CORRECTO
if result.Error != nil {
    utils.NotFoundResponse(c, "Tarea")
    return
}

// ❌ INCORRECTO
if result.Error != nil {
    c.JSON(500, gin.H{"error": "algo salió mal"})
}
```

### Códigos HTTP apropiados:

| Código | Uso | Ejemplo |
|--------|-----|---------|
| 200 | Éxito | GET exitoso |
| 201 | Creado | POST exitoso |
| 400 | Error de validación | JSON inválido |
| 404 | No encontrado | Recurso inexistente |
| 500 | Error del servidor | Error de BD |

---

## 4. Validaciones

### ✅ Buena Práctica: Validaciones en el modelo

**Archivo:** `models/task.go`

```go
type CreateTaskInput struct {
    Title       string `json:"title" binding:"required,min=3,max=200"`
    Description string `json:"description" binding:"max=1000"`
}
```

**Validaciones disponibles:**
- `required` - Campo obligatorio
- `min=N` - Longitud/valor mínimo
- `max=N` - Longitud/valor máximo
- `email` - Formato de email
- `url` - Formato de URL
- `oneof=red green blue` - Uno de los valores

**Por qué:**
- ✅ Validación automática por Gin
- ✅ Código más limpio
- ✅ Errores claros para el cliente
- ✅ Seguridad contra datos inválidos

**Constantes para límites:**

**Archivo:** `utils/constants.go`

```go
const (
    MaxTitleLength = 200
    MinTitleLength = 3
    MaxDescriptionLength = 1000
)
```

---

## 5. Middleware

### ✅ Buena Práctica: Middleware reutilizables

#### Logger Personalizado
**Archivo:** `middleware/logger.go`

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()
        c.Next()
        latency := time.Since(startTime)
        // Log detallado...
    }
}
```

**Beneficios:**
- ✅ Debugging más fácil
- ✅ Monitoreo de performance
- ✅ Auditoría de requests

#### CORS Configurado
**Archivo:** `middleware/cors.go`

```go
func CORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins: config.AppConfig.CORS.AllowedOrigins,
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
        // ...
    })
}
```

**Beneficios:**
- ✅ Permite requests desde frontend
- ✅ Seguridad configurada
- ✅ Control de orígenes permitidos

### Otros middleware útiles:

```go
// Rate limiting
router.Use(ratelimit.Middleware())

// Autenticación JWT
router.Use(auth.JWTMiddleware())

// Compresión GZIP
router.Use(gzip.Gzip(gzip.DefaultCompression))
```

---

## 6. Base de Datos

### ✅ Buena Práctica: Abstracción de base de datos

**Archivo:** `database/database.go`

```go
// Soporta múltiples bases de datos
switch config.AppConfig.Database.Type {
case "sqlite":
    dialector = sqlite.Open(config.AppConfig.Database.Name)
case "postgres":
    dialector = postgres.Open(dsn)
}
```

**Por qué:**
- ✅ Fácil cambiar de SQLite a PostgreSQL
- ✅ Mismo código para diferentes ambientes
- ✅ Configuración sin cambiar código

### Migraciones automáticas:

```go
// AutoMigrate actualiza el schema automáticamente
DB.AutoMigrate(&models.Task{})
```

**Ventajas:**
- ✅ No escribir SQL manualmente
- ✅ Schema actualizado con el código
- ✅ Desarrollo más rápido

**⚠️ Producción:**
Para producción, considera usar migraciones versionadas con herramientas como:
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)

---

## 7. Respuestas HTTP Estandarizadas

### ✅ Estructura consistente

```go
type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Code    int         `json:"code"`
}
```

**Respuesta exitosa:**
```json
{
  "success": true,
  "message": "Tarea creada exitosamente",
  "data": { "id": 1, "title": "..." },
  "code": 201
}
```

**Respuesta de error:**
```json
{
  "success": false,
  "error": "Tarea no encontrada",
  "code": 404
}
```

**Beneficios:**
- ✅ Cliente sabe qué esperar
- ✅ Fácil manejar errores en frontend
- ✅ API profesional y consistente

---

## 8. Seguridad

### Prácticas implementadas:

#### 1. **Validación de entrada**
```go
// Validar automáticamente con binding tags
if err := c.ShouldBindJSON(&input); err != nil {
    utils.ValidationErrorResponse(c, err)
    return
}
```

#### 2. **CORS configurado**
```go
// Solo orígenes permitidos
AllowOrigins: config.AppConfig.CORS.AllowedOrigins
```

#### 3. **Soft Delete**
```go
// GORM usa soft delete por defecto
// Los registros no se eliminan realmente
DB.Delete(&task) // Marca DeletedAt, no elimina
```

### 🔒 Seguridad adicional recomendada:

#### Rate Limiting
```go
import "github.com/ulule/limiter/v3"

// Limitar requests por IP
store := memory.NewStore()
rate := limiter.Rate{
    Period: 1 * time.Hour,
    Limit:  100,
}
```

#### Autenticación JWT
```go
import "github.com/golang-jwt/jwt/v5"

// Proteger rutas con JWT
router.Use(jwtMiddleware())
```

#### Sanitización de inputs
```go
import "github.com/microcosm-cc/bluemonday"

// Limpiar HTML/XSS
p := bluemonday.StrictPolicy()
cleanTitle := p.Sanitize(input.Title)
```

#### HTTPS en producción
```go
// Usar TLS certificates
router.RunTLS(":443", "cert.pem", "key.pem")
```

---

## 9. Testing

### Estructura de tests:

```
backend-api-go/
├── handlers/
│   ├── task_handler.go
│   └── task_handler_test.go
├── models/
│   ├── task.go
│   └── task_test.go
```

### Ejemplo de test:

```go
package handlers_test

import (
    "testing"
    "net/http/httptest"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestGetAllTasks(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.GET("/tasks", handlers.GetAllTasks)

    // Request
    req := httptest.NewRequest("GET", "/tasks", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, 200, w.Code)
}
```

**Ejecutar tests:**
```bash
go test ./...                # Todos los tests
go test -v ./handlers        # Tests de handlers
go test -cover ./...         # Con cobertura
```

---

## 10. Deploy

### Opción 1: Binario compilado

```bash
# Compilar para Linux
GOOS=linux GOARCH=amd64 go build -o api-server

# Ejecutar
./api-server
```

### Opción 2: Docker

**Dockerfile:**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api-server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-server .
COPY .env.example .env
EXPOSE 8080
CMD ["./api-server"]
```

**docker-compose.yml:**
```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_TYPE=postgres
      - DB_HOST=db
    depends_on:
      - db

  db:
    image: postgres:15
    environment:
      POSTGRES_DB: tasks_db
      POSTGRES_PASSWORD: password
```

### Opción 3: Plataformas Cloud

#### Railway
```bash
# Instalar Railway CLI
npm i -g @railway/cli

# Deploy
railway login
railway init
railway up
```

#### Fly.io
```bash
# Instalar Fly CLI
curl -L https://fly.io/install.sh | sh

# Deploy
fly launch
fly deploy
```

---

## 📊 Checklist de Buenas Prácticas

### Código
- [x] Estructura de carpetas organizada
- [x] Separación de responsabilidades
- [x] Constantes para valores mágicos
- [x] Comentarios en funciones públicas
- [x] Nombres descriptivos de variables

### Configuración
- [x] Variables de entorno con `.env`
- [x] Configuración centralizada
- [x] Diferentes configs por ambiente
- [x] `.env.example` documentado

### API
- [x] Códigos HTTP apropiados
- [x] Respuestas estandarizadas
- [x] Validación de inputs
- [x] Manejo de errores consistente
- [x] Documentación de endpoints

### Base de Datos
- [x] ORM para queries
- [x] Migraciones automáticas
- [x] Soft deletes
- [x] Índices en campos frecuentes

### Seguridad
- [x] Validación de inputs
- [x] CORS configurado
- [ ] Rate limiting (TODO)
- [ ] Autenticación JWT (TODO)
- [ ] HTTPS en producción (TODO)

### Middleware
- [x] Logger personalizado
- [x] CORS
- [x] Recovery de panics
- [ ] Rate limiting (TODO)
- [ ] Compresión GZIP (TODO)

### Testing
- [ ] Tests unitarios (TODO)
- [ ] Tests de integración (TODO)
- [ ] Cobertura > 70% (TODO)

### Deploy
- [ ] Dockerfile (TODO)
- [ ] CI/CD pipeline (TODO)
- [ ] Monitoring (TODO)
- [ ] Logs centralizados (TODO)

---

## 🎯 Próximos Pasos

1. **Agregar autenticación**
   - JWT tokens
   - Roles y permisos
   - Refresh tokens

2. **Mejorar base de datos**
   - Migraciones versionadas
   - Seeders para datos de prueba
   - Índices optimizados

3. **Testing completo**
   - Tests unitarios
   - Tests de integración
   - Tests E2E

4. **Monitoreo**
   - Prometheus metrics
   - Grafana dashboards
   - Error tracking (Sentry)

5. **Documentación**
   - Swagger/OpenAPI
   - Postman collection
   - Ejemplos de cliente

---

## 📚 Recursos Adicionales

- [Effective Go](https://go.dev/doc/effective_go)
- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [The Twelve-Factor App](https://12factor.net/)
