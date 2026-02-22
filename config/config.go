package config

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL string
    JWTSecret   string
    Port        string
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    return &Config{
        DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/invoice_db?sslmode=disable"),
        JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
        Port:        getEnv("PORT", "8080"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}