package database

import (
	"fmt"
	"log"

	"github.com/AndresEGV/backend-api-go/config"
	"github.com/AndresEGV/backend-api-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB es la instancia de la base de datos que usaremos en toda la aplicación
var DB *gorm.DB

// Connect inicializa la conexión a la base de datos usando la configuración
func Connect() {
	var err error
	var dialector gorm.Dialector

	// Configurar el logger de GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silenciar logs SQL en producción
	}

	// Seleccionar el driver según la configuración
	switch config.AppConfig.Database.Type {
	case "sqlite":
		dialector = sqlite.Open(config.AppConfig.Database.Name)
		log.Printf("📁 Conectando a SQLite: %s", config.AppConfig.Database.Name)

	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.AppConfig.Database.Host,
			config.AppConfig.Database.User,
			config.AppConfig.Database.Password,
			config.AppConfig.Database.Name,
			config.AppConfig.Database.Port,
			config.AppConfig.Database.SSLMode,
		)
		dialector = postgres.Open(dsn)
		log.Printf("🐘 Conectando a PostgreSQL: %s@%s:%s/%s",
			config.AppConfig.Database.User,
			config.AppConfig.Database.Host,
			config.AppConfig.Database.Port,
			config.AppConfig.Database.Name,
		)

	default:
		log.Fatalf("❌ Tipo de base de datos no soportado: %s", config.AppConfig.Database.Type)
	}

	// Conectar a la base de datos
	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		log.Fatal("❌ Error al conectar a la base de datos:", err)
	}

	log.Println("✅ Conexión a la base de datos establecida")

	// AutoMigrate crea/actualiza las tablas basándose en los modelos
	// Esto es el "ORM en acción": convierte tus structs de Go en tablas SQL
	err = DB.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatal("❌ Error al migrar modelos:", err)
	}

	log.Println("✅ Migraciones ejecutadas correctamente")
}
