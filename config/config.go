package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config almacena las variables de entorno de la aplicación
type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	AppEnv     string
}

// LoadConfig lee las configuraciones de entorno de forma segura
func LoadConfig() *Config {
	// Intentamos cargar el archivo .env si existe (entorno de desarrollo)
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró un archivo .env activo, se usarán variables del sistema.")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "ebooks_db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		AppEnv:     getEnv("APP_ENV", "production"),
	}
}

// Función auxiliar para leer variables con un valor por defecto
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
