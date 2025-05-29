package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    DatabaseURL string
}

func Load() *Config {
    _ = godotenv.Load()

    dbUser := getEnv("DB_USER", "root")
    dbPass := getEnv("DB_PASS", "root")
    dbHost := getEnv("DB_HOST", "localhost")
    dbPort := getEnv("DB_PORT", "3306")
    dbName := getEnv("DB_NAME", "ordermatching")

    dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        dbUser, dbPass, dbHost, dbPort, dbName)

    return &Config{
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: dbURL,
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
