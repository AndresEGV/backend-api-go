package database

import (
	"log"

	"github.com/AndresEGV/backend-api-go/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB es la instancia de la base de datos que usaremos en toda la aplicación
var DB *gorm.DB

// Connect inicializa la conexión a la base de datos
func Connect() {
	var err error

	// Conectar a SQLite (crea el archivo tasks.db)
	DB, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	log.Println("✅ Conexión a la base de datos establecida")

	// AutoMigrate crea/actualiza las tablas basándose en los modelos
	// Esto es el "ORM en acción": convierte tus structs de Go en tablas SQL
	err = DB.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatal("Error al migrar modelos:", err)
	}

	log.Println("✅ Migraciones ejecutadas correctamente")
}
