# Ejemplos de Uso de la API

Este documento contiene ejemplos prácticos de cómo usar la API de Tareas.

## Usando curl (Terminal)

### 1. Ver información de la API
```bash
curl http://localhost:8080/
```

**Respuesta:**
```json
{
  "message": "¡Bienvenido a la API de Tareas en Go! 🚀",
  "version": "1.0.0",
  "endpoints": {
    "GET /tasks": "Obtener todas las tareas",
    "GET /tasks/:id": "Obtener una tarea por ID",
    "POST /tasks": "Crear una nueva tarea",
    "PUT /tasks/:id": "Actualizar una tarea",
    "DELETE /tasks/:id": "Eliminar una tarea"
  }
}
```

### 2. Crear tareas
```bash
# Tarea 1
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Aprender Go",
    "description": "Completar el tutorial de API REST en Go"
  }'

# Tarea 2
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Practicar GORM",
    "description": "Hacer ejercicios con el ORM"
  }'

# Tarea 3
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Crear mi propia API"
  }'
```

**Respuesta (ejemplo):**
```json
{
  "message": "Tarea creada exitosamente",
  "data": {
    "ID": 1,
    "CreatedAt": "2024-11-04T10:30:00Z",
    "UpdatedAt": "2024-11-04T10:30:00Z",
    "DeletedAt": null,
    "title": "Aprender Go",
    "description": "Completar el tutorial de API REST en Go",
    "completed": false
  }
}
```

### 3. Obtener todas las tareas
```bash
curl http://localhost:8080/tasks
```

**Respuesta:**
```json
{
  "count": 3,
  "data": [
    {
      "ID": 1,
      "title": "Aprender Go",
      "description": "Completar el tutorial de API REST en Go",
      "completed": false
    },
    {
      "ID": 2,
      "title": "Practicar GORM",
      "description": "Hacer ejercicios con el ORM",
      "completed": false
    },
    {
      "ID": 3,
      "title": "Crear mi propia API",
      "description": "",
      "completed": false
    }
  ]
}
```

### 4. Obtener una tarea específica
```bash
curl http://localhost:8080/tasks/1
```

**Respuesta:**
```json
{
  "data": {
    "ID": 1,
    "title": "Aprender Go",
    "description": "Completar el tutorial de API REST en Go",
    "completed": false
  }
}
```

### 5. Actualizar una tarea
```bash
# Marcar como completada
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "completed": true
  }'

# Actualizar título y descripción
curl -X PUT http://localhost:8080/tasks/2 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Practicar GORM - Actualizado",
    "description": "Completé los ejercicios básicos",
    "completed": true
  }'
```

**Respuesta:**
```json
{
  "message": "Tarea actualizada exitosamente",
  "data": {
    "ID": 1,
    "title": "Aprender Go",
    "description": "Completar el tutorial de API REST en Go",
    "completed": true
  }
}
```

### 6. Eliminar una tarea
```bash
curl -X DELETE http://localhost:8080/tasks/3
```

**Respuesta:**
```json
{
  "message": "Tarea eliminada exitosamente"
}
```

---

## Usando JavaScript (Fetch API)

### Crear una tarea
```javascript
fetch('http://localhost:8080/tasks', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: 'Aprender Go',
    description: 'Completar el tutorial de API REST'
  })
})
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

### Obtener todas las tareas
```javascript
fetch('http://localhost:8080/tasks')
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

### Actualizar una tarea
```javascript
fetch('http://localhost:8080/tasks/1', {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    completed: true
  })
})
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

### Eliminar una tarea
```javascript
fetch('http://localhost:8080/tasks/1', {
  method: 'DELETE',
})
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

---

## Usando Python (requests)

Primero instala requests:
```bash
pip install requests
```

### Script completo
```python
import requests
import json

BASE_URL = "http://localhost:8080"

# Crear una tarea
def create_task(title, description=""):
    response = requests.post(
        f"{BASE_URL}/tasks",
        json={"title": title, "description": description}
    )
    print("Crear:", response.json())
    return response.json()

# Obtener todas las tareas
def get_tasks():
    response = requests.get(f"{BASE_URL}/tasks")
    print("Todas las tareas:", response.json())
    return response.json()

# Obtener una tarea
def get_task(task_id):
    response = requests.get(f"{BASE_URL}/tasks/{task_id}")
    print(f"Tarea {task_id}:", response.json())
    return response.json()

# Actualizar una tarea
def update_task(task_id, **kwargs):
    response = requests.put(
        f"{BASE_URL}/tasks/{task_id}",
        json=kwargs
    )
    print("Actualizar:", response.json())
    return response.json()

# Eliminar una tarea
def delete_task(task_id):
    response = requests.delete(f"{BASE_URL}/tasks/{task_id}")
    print("Eliminar:", response.json())
    return response.json()

# Ejemplo de uso
if __name__ == "__main__":
    # Crear tareas
    task1 = create_task("Aprender Go", "Completar tutorial")
    task2 = create_task("Practicar GORM")

    # Ver todas
    get_tasks()

    # Ver una específica
    get_task(1)

    # Actualizar
    update_task(1, completed=True)

    # Eliminar
    delete_task(2)

    # Ver el resultado final
    get_tasks()
```

---

## Probando con Postman

### 1. Importar colección

Crea un archivo `postman_collection.json` con:

```json
{
  "info": {
    "name": "API Tareas Go",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Obtener todas las tareas",
      "request": {
        "method": "GET",
        "url": "http://localhost:8080/tasks"
      }
    },
    {
      "name": "Crear tarea",
      "request": {
        "method": "POST",
        "url": "http://localhost:8080/tasks",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"title\": \"Nueva tarea\",\n  \"description\": \"Descripción de la tarea\"\n}"
        }
      }
    },
    {
      "name": "Actualizar tarea",
      "request": {
        "method": "PUT",
        "url": "http://localhost:8080/tasks/1",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"completed\": true\n}"
        }
      }
    },
    {
      "name": "Eliminar tarea",
      "request": {
        "method": "DELETE",
        "url": "http://localhost:8080/tasks/1"
      }
    }
  ]
}
```

### 2. Importar en Postman
- Abre Postman
- Click en "Import"
- Selecciona el archivo `postman_collection.json`
- ¡Listo! Puedes probar todos los endpoints

---

## Casos de error

### 1. Crear tarea sin título (campo requerido)
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Solo descripción"
  }'
```

**Respuesta:**
```json
{
  "error": "Key: 'CreateTaskInput.Title' Error:Field validation for 'Title' failed on the 'required' tag"
}
```

### 2. Obtener tarea que no existe
```bash
curl http://localhost:8080/tasks/999
```

**Respuesta:**
```json
{
  "error": "Tarea no encontrada"
}
```

### 3. Actualizar tarea que no existe
```bash
curl -X PUT http://localhost:8080/tasks/999 \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```

**Respuesta:**
```json
{
  "error": "Tarea no encontrada"
}
```

---

## Tips para desarrollo

### Ver logs del servidor
Los logs te ayudarán a entender qué está pasando:

```
[GIN] 2024/11/04 - 10:30:45 | 200 |     125.5µs |       127.0.0.1 | GET      "/tasks"
[GIN] 2024/11/04 - 10:31:12 | 201 |    1.234ms |       127.0.0.1 | POST     "/tasks"
[GIN] 2024/11/04 - 10:31:45 | 404 |      89.2µs |       127.0.0.1 | GET      "/tasks/999"
```

### Usar herramientas de desarrollo
- **httpie**: Alternativa más amigable a curl
  ```bash
  # Instalar
  brew install httpie  # macOS
  apt install httpie   # Linux

  # Usar
  http GET localhost:8080/tasks
  http POST localhost:8080/tasks title="Nueva tarea"
  ```

- **Thunder Client**: Extensión de VS Code (alternativa a Postman)

### Reiniciar la base de datos
Si quieres empezar de cero, simplemente elimina el archivo:
```bash
rm tasks.db
# Al reiniciar el servidor, se creará una nueva base de datos vacía
```
