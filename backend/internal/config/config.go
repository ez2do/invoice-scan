package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Host         string
	GeminiAPIKey string
	CORSOrigin   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Port:         getEnv("PORT", "3001"),
		Host:         getEnv("HOST", "localhost"),
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
		CORSOrigin:   getEnv("CORS_ORIGIN", "http://localhost:5173"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
