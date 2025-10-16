package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort      string
	DatabaseURL     string
	JWTSecret       string
	LogLevel        string
	CertFilePath    string
	KeyFilePath     string
	DBHost          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBPort          string
	DevelopmentMode bool
}

var AppConfig Config

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		CertFilePath:    getEnv("CERT_FILE_PATH", "cert.pem"),
		KeyFilePath:     getEnv("KEY_FILE_PATH", "key.pem"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "data_mapping"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DevelopmentMode: getEnv("DEVELOPMENT_MODE", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
