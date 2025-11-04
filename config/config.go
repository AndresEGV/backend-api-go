package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config almacena toda la configuración de la aplicación
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
}

// ServerConfig configuración del servidor HTTP
type ServerConfig struct {
	Port    string
	GinMode string
}

// DatabaseConfig configuración de la base de datos
type DatabaseConfig struct {
	Type     string // sqlite, postgres, mysql
	Name     string
	Host     string
	Port     string
	User     string
	Password string
	SSLMode  string
}

// CORSConfig configuración de CORS
type CORSConfig struct {
	AllowedOrigins []string
}

var AppConfig *Config

// Load carga la configuración desde variables de entorno
func Load() {
	// Cargar archivo .env si existe (en desarrollo)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No se encontró archivo .env, usando variables de entorno del sistema")
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "sqlite"),
			Name:     getEnv("DB_NAME", "tasks.db"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ORIGINS", []string{"*"}),
		},
	}

	log.Println("✅ Configuración cargada correctamente")
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsSlice obtiene una variable de entorno como slice separado por comas
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
