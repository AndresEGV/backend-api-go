# Comparación: Gin vs Chi vs Stdlib

## 🎯 Mismo Endpoint, 3 Enfoques

### Crear una Tarea (POST /tasks)

#### Con Gin (Actual)
```go
func CreateTask(c *gin.Context) {
    var input models.CreateTaskInput

    // Validación automática
    if err := c.ShouldBindJSON(&input); err != nil {
        utils.ValidationErrorResponse(c, err)
        return
    }

    task := models.Task{
        Title:       input.Title,
        Description: input.Description,
    }

    database.DB.Create(&task)

    // Helper de Gin
    utils.SuccessResponse(c, 201, "Creada", task)
}
```

**Líneas de código:** ~15
**Dependencias:** ~40 paquetes
**Abstracción:** Alta (oculta `http.ResponseWriter`)

---

#### Con Chi (Recomendado)
```go
func CreateTask(w http.ResponseWriter, r *http.Request) {
    var input models.CreateTaskInput

    // Decodificar JSON manualmente
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        render.Status(r, 400)
        render.JSON(w, r, APIResponse{Error: err.Error()})
        return
    }

    // Validación manual (o usar biblioteca)
    if len(input.Title) < 3 {
        render.Status(r, 400)
        render.JSON(w, r, APIResponse{Error: "Title muy corto"})
        return
    }

    task := models.Task{
        Title:       input.Title,
        Description: input.Description,
    }

    database.DB.Create(&task)

    // Helper de Chi
    render.Status(r, 201)
    render.JSON(w, r, APIResponse{
        Success: true,
        Data:    task,
    })
}
```

**Líneas de código:** ~25
**Dependencias:** 1 paquete pequeño
**Abstracción:** Media (usa stdlib + helpers)

---

#### Con Stdlib Puro
```go
func CreateTask(w http.ResponseWriter, r *http.Request) {
    var input models.CreateTaskInput

    // Decodificar JSON
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        json.NewEncoder(w).Encode(map[string]string{
            "error": err.Error(),
        })
        return
    }

    // Validación manual
    if len(input.Title) < 3 || len(input.Title) > 200 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Title debe tener entre 3 y 200 caracteres",
        })
        return
    }

    task := models.Task{
        Title:       input.Title,
        Description: input.Description,
    }

    database.DB.Create(&task)

    // Respuesta manual
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(201)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data":    task,
    })
}
```

**Líneas de código:** ~35
**Dependencias:** 0 externas
**Abstracción:** Ninguna (HTTP puro)

---

## 📊 Comparación Completa

### 1. Routing (Definir Rutas)

#### Gin
```go
router := gin.Default()
router.GET("/tasks/:id", handlers.GetTask)
router.POST("/tasks", handlers.CreateTask)
```

#### Chi
```go
r := chi.NewRouter()
r.Get("/tasks/{id}", handlers.GetTask)
r.Post("/tasks", handlers.CreateTask)
```

#### Stdlib (Go 1.22+)
```go
mux := http.NewServeMux()
mux.HandleFunc("GET /tasks/{id}", handlers.GetTask)
mux.HandleFunc("POST /tasks", handlers.CreateTask)
```

**Ganador:** Empate (todos simples en Go 1.22+)

---

### 2. Extraer Parámetros de URL

#### Gin
```go
id := c.Param("id")  // /tasks/:id
```

#### Chi
```go
id := chi.URLParam(r, "id")  // /tasks/{id}
```

#### Stdlib (Go 1.22+)
```go
id := r.PathValue("id")  // /tasks/{id}
```

**Ganador:** Gin (más corto), pero todos son fáciles

---

### 3. JSON Response

#### Gin
```go
c.JSON(200, gin.H{
    "message": "OK",
    "data":    data,
})
```

#### Chi
```go
render.JSON(w, r, map[string]interface{}{
    "message": "OK",
    "data":    data,
})
```

#### Stdlib
```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(200)
json.NewEncoder(w).Encode(map[string]interface{}{
    "message": "OK",
    "data":    data,
})
```

**Ganador:** Gin (más conciso)

---

### 4. Middleware

#### Gin
```go
router.Use(gin.Logger())
router.Use(gin.Recovery())
router.Use(customMiddleware())
```

#### Chi
```go
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(customMiddleware)
```

#### Stdlib
```go
// Tienes que envolver manualmente
mux.HandleFunc("/", loggingMiddleware(recoveryMiddleware(handler)))
```

**Ganador:** Chi y Gin (stdlib es más verboso)

---

### 5. Grupos de Rutas

#### Gin
```go
tasks := router.Group("/tasks")
{
    tasks.GET("", handlers.GetAll)
    tasks.POST("", handlers.Create)
}
```

#### Chi
```go
r.Route("/tasks", func(r chi.Router) {
    r.Get("/", handlers.GetAll)
    r.Post("/", handlers.Create)
})
```

#### Stdlib
```go
// No hay concepto de grupos
// Tienes que repetir el prefijo
mux.HandleFunc("GET /tasks", handlers.GetAll)
mux.HandleFunc("POST /tasks", handlers.Create)
```

**Ganador:** Chi (más elegante)

---

## 🏆 Tabla Resumen

| Característica | Gin | Chi | Stdlib |
|----------------|-----|-----|--------|
| **Líneas de código** | Menos | Medio | Más |
| **Dependencias** | 40+ | 1 | 0 |
| **Rendimiento** | Bueno | Excelente | Excelente |
| **Curva aprendizaje** | Baja | Media | Alta |
| **Flexibilidad** | Media | Alta | Total |
| **Estabilidad** | Media | Alta | Total |
| **Comunidad** | Grande | Mediana | Oficial |
| **Compatibilidad stdlib** | Baja | 100% | 100% |
| **Validaciones** | ✅ Auto | ❌ Manual | ❌ Manual |
| **JSON helpers** | ✅ Si | ✅ Si | ❌ No |
| **Routing avanzado** | ✅ Si | ✅ Si | ⚠️ Básico |
| **Middleware** | ✅ Muchos | ✅ Muchos | ⚠️ Manual |

---

## 🎯 Cuándo Usar Cada Uno

### Usa **Gin** SI:
- ✅ Necesitas **rapidez extrema** de desarrollo
- ✅ Tu equipo viene de **Node.js/Express**
- ✅ Quieres **validaciones automáticas**
- ✅ Es un **prototipo** o MVP
- ✅ El proyecto durará **< 2 años**

### Usa **Chi** SI:
- ✅ Quieres **balance** entre stdlib y framework
- ✅ Valoras **compatibilidad stdlib** al 100%
- ✅ Necesitas **middleware** sin peso extra
- ✅ Proyecto **mediano/grande** de producción
- ✅ Tu equipo **conoce Go bien**

### Usa **Stdlib** SI:
- ✅ Valoras **cero dependencias**
- ✅ Necesitas **máximo rendimiento**
- ✅ El proyecto durará **10+ años**
- ✅ Trabajas en **empresa grande** (Google, Uber)
- ✅ Prefieres **control total**

---

## 💡 Mi Recomendación Final

### Para este proyecto (Aprendiendo Go):

**Orden recomendado de aprendizaje:**

1. **Primero:** Usa **Stdlib** puro → Entiendes HTTP real
2. **Segundo:** Prueba **Chi** → Ves el valor de herramientas ligeras
3. **Tercero:** Compara con **Gin** → Entiendes los trade-offs

### Para proyectos reales 2025:

- **Startup/MVP:** Gin (velocidad > todo)
- **Empresa mediana:** Chi (balance perfecto)
- **Empresa grande:** Stdlib + helpers propios (control total)

---

## 📈 Tendencia 2025

```
2018: 90% Gin/Echo/Fiber
2020: 70% Gin/Echo, 20% Chi, 10% Stdlib
2022: 50% Gin/Echo, 30% Chi, 20% Stdlib
2024: 30% Gin/Echo, 40% Chi, 30% Stdlib
2025: 20% Gin/Echo, 40% Chi, 40% Stdlib ← Estamos aquí
```

**Conclusión:** La comunidad Go está volviendo a las raíces (stdlib).

---

## 🚀 Migración

### De Gin a Chi

```bash
# Instalar Chi
go get github.com/go-chi/chi/v5
go get github.com/go-chi/render

# Cambiar imports
- "github.com/gin-gonic/gin"
+ "github.com/go-chi/chi/v5"
+ "github.com/go-chi/render"

# Cambiar función signatures
- func Handler(c *gin.Context)
+ func Handler(w http.ResponseWriter, r *http.Request)
```

### De Chi a Stdlib

```bash
# Remover Chi
go mod tidy

# Ya sabes usar stdlib!
```

---

## 🎓 Ejercicio Práctico

**Desafío:** Implementa el endpoint `GET /tasks/{id}` en los 3 estilos.

Tienes los ejemplos en:
- `main.go` (Gin - actual)
- `stdlib-version/main-chi.go` (Chi)
- `stdlib-version/main-stdlib.go` (Stdlib puro)

**Compara:**
1. Líneas de código
2. Claridad
3. Qué prefieres y por qué

---

## 📚 Recursos

- [Stdlib net/http](https://pkg.go.dev/net/http)
- [Chi Router](https://github.com/go-chi/chi)
- [Gin Framework](https://gin-gonic.com/)
- [Go 1.22 Release Notes](https://go.dev/doc/go1.22) (routing mejorado)
